# F-02 Slice 8 (corrective) — read-path degrades on unavailable cost/ceiling reader

```yaml
origin: P7 live QA (C10) — hub ruling caminho (a), IN-SCOPE M-01 corrective slice
adr: ADR-17 (unknown operational facts -> null/unknown, never block the read)
review_sha_base: e2cde36 (delta gate will diff from here)
```

## Defect (live-QA-caught)

`ReadService.List`, `ByProduct`, and `Get` **hard-fail 503 `source_unavailable`** when the
**Oracle-backed** ICMS ceiling reader (`facts.GetICMSCeilingByOrigin`) or cost reader
(`facts.GetCostFactsByIDs`) is unavailable — while `Summary` **degrades** the same two errors to
`below_margin_worst_case: null` / `margin_unknown` and still returns rows. This couples the whole
listings read spine's availability to Oracle for an **optional** margin field, violating ADR-17
and diverging from Summary. Proven live: dev stack has no Oracle → `GET /listings/summary` 200 but
`GET /listings` + `/listings/by-product` 503 (blocks C10 secondary observation of MLB ids +
unknown% over 34 successfully-ingested rows).

## Defect loci (worktree `apps/server_core/internal/modules/listings/application/read_service.go`)

| Method | Hard-fail site | Reader | Should |
|--------|----------------|--------|--------|
| `List` | `:122-125` ceiling; `enrich` `:350-353` cost | Oracle | degrade → serve rows, `below_margin_worst_case:null`, cost `null`, count `margin_unknown` |
| `ByProduct` | `:157-159` ceiling; `enrichGroups`→`enrich` cost | Oracle | same |
| `Get` (detail) | `:278-281` ceiling; `:283-285` enrich cost | Oracle | same (single row) |
| `scan` `:298` / `scanGroups` `:177` | below-margin FILTER path | Oracle | see design note |

**Reference (the correct pattern already in the same file):** `Summary` `:59-67` (ceilingErr →
`nil,nil,return row`) and `:78-83` (cost err → `nil,nil,return row`). Policy + installation errors
stay HARD in Summary (Postgres-backed, available) — keep them hard in the fix too. Only the two
Oracle-backed readers (ceiling + cost) degrade.

## Required behavior

- Default `List` / `ByProduct` / `Get` (no below_margin filter): when ceiling OR cost reader is
  unavailable → **HTTP 200**, rows served, each row `cost: null`, `below_margin_worst_case: null`;
  summary-style counters where applicable count those as `margin_unknown`. Never 503 for this.
- Preserve exact behavior when the readers ARE available (no regression to C04/C05/C07/C09).
- Shapes unchanged (fields already nullable) → **no OpenAPI/SDK change** (behavior-only).

## Design note — below_margin FILTER path (`needsBelowMarginScan` true, `scan`/`scanGroups`)

A `filter.exception=below_margin_worst_case` query genuinely needs cost+ceiling. Plan decides the
honest degradation: preferred = return the scanned rows with `below_margin_worst_case:null` +
`margin_unknown` (do NOT silently drop or 503); acceptable = documented `service_degraded` 200
empty-with-flag. **No hard 503 on reader-unavailable anywhere.** Planner picks and justifies.

## Test-first (worker: Luna high)

Add failing tests BEFORE the fix, then make them pass:
- Unit (`read_service_test.go`): `fakeFacts` returning an error from `GetICMSCeilingByOrigin` AND
  from `GetCostFactsByIDs` → `List` and `ByProduct` and `Get` each return **nil error** with rows
  present, `BelowMarginWorstCase == nil`, cost nil (mirror the existing Summary degradation test).
- Integration if a lane fits: reader-unavailable → endpoint 200 + null margin.
- Regression: existing available-reader tests stay green.

## Scope guard

ONLY `internal/modules/listings/application/read_service.go` + its tests. Do NOT touch
`internal_read`/Oracle adapters, other modules, transport shapes, OpenAPI/SDK, migrations, or the
hub's `docker/dev/*.sh` / `.env`. No new dependencies.

## ADDENDUM (2026-07-16, FINAL) — hub ACK: wire real cost + policy readers

Hub verified the wiring diagnosis in-worktree and APPROVED full wiring; the transient deferral path
is VOID (hub deleted the board deferral-tracker task — "sem stub alcançável = melhor que deferral").
NO-STUB doctrine satisfied by REMOVING the stubs, not by deferring them. **Corrective slice wires both
readers real AND keeps fix (a) degrade (mandatory — Oracle can drop in prod).**

### Wiring loci (composition `internal/composition/root.go`) — chip-diagnosed + hub-verified

| Site | Now (violation) | Fix | Real, not a stub |
|------|-----------------|-----|------------------|
| `:494` cost reader | `NewBatchReader(nil, sem)` | `NewBatchReader(oracleDB, sem)` | `oracleDB` `:352` (`internalreadoracle.Database`), real when `internalReadAvailable`; identical to profitability `:474`. |
| `:495` policy reader | `NewPolicyReader(unavailablePolicyService{…})` | `NewPolicyReader(marketSvc)` | `marketSvc` `:491` (`marketplacesapp.NewService`), Postgres-backed (`marketplaces/application/service.go:109`→`adapters/postgres/repository.go:178`); satisfies `listingsmarketplaces.Service.GetPricingPolicyForInstallation` field-for-field. |

- **Keep** the `unavailableListingPolicyReader` wrapper (`root.go:115-123`) wrapping the REAL
  `NewPolicyReader(marketSvc)` — it translates runtime `SourceUnavailable`→graceful degrade (ADR-17).
- **Remove** the now-dead `unavailablePolicyService` type (`root.go:107-113`) — no reachable stub left
  (AI-slop checklist clean).

### Zero-value oracleDB safety (hub-flagged — reviewer MUST check)

When Oracle env absent/boot-fails, `internalReadAvailable=false` and `oracleDB` stays untyped-nil
(only assigned real `db` at `root.go:359` inside the available branch). `NewBatchReader(oracleDB, sem)`
is SAFE: `BatchReader.ensureBatchAvailable` (`batch_reader.go:204-212`) checks `r.db == nil` and
returns `ReadErrorSourceUnavailable` BEFORE any `r.db` call (guard hit at `:46`/`:104`) → **clean
degrade via fix (a), NO panic**. Untyped-nil interface → `r.db == nil` true (no typed-nil trap). No
extra `if internalReadAvailable` guard required; the reader self-guards. Reviewer confirms this branch.

### Scope

Corrective slice = `read_service.go` + `read_service_test.go` (fix a) **AND**
`internal/composition/root.go` (wiring; remove dead stub type). No adapter/OpenAPI/SDK/migration
change. Composition is a hub seam — hub explicitly authorized this touch. Still no `docker/dev/*.sh`
/ `.env` commits (hub cleans at close).

## Verify (chip, post-impl)

L0 `go build ./... && go vet` (listings + composition) + L1 `go test` listings unit + targeted
integration (`GOCACHE` absolute). Wiring proof = L0 build/vet + the C10 live re-drive (composition
wiring is out-of-band for unit tests; the live lane is its proof). Then independent slice review,
commit in worktree, REQUEST hub restart-backend, re-drive C10 (cost+policy now real → richer
evidence), dual-gate DELTA from e2cde36.
