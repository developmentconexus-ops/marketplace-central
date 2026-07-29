# CHIP-VENDAVEL — dispatch ledger

Chip: CHIP-VENDAVEL (MIS-006, interim post-mission chip).
Worktree: `.claude/worktrees/chip-vendavel` · branch `worktree-chip-vendavel`.
BASE-SHA: `554788d576d04b21719f4a17e4702dd0f0aff4e1`.
Hub session: `local_99feb041-a5b3-4161-b6dc-bd38e65b6156`.
Contract: `.mnfs/MIS-006-integracao-fundacao/_pos-chip-vendavel/validation-contract.md` (VC-1..VC-7).

Ownership declared (write-set, single track — no concurrent chip in flight):
`apps/server_core/migrations/` block 0083–0084 only · `internal_read/adapters/oracle/{sync.go,catalog_page.go}` ·
`products_mirror` writers/readers (additive columns `usoprod`, `ad_ecommerce`) · xlsx parser optional columns ·
module `tenant_config` (3 boolean columns + handler) · `MirrorMatcher` path ·
FE `apps/web/src/pages/integracoes/IntegracoesPage.tsx`, `packages/feature-products` catalog page,
`packages/web-query` + `packages/sdk-runtime` domain blocks.
Pre-authorised additive grants: composition-root wiring region (`internal/composition/root.go`),
SDK domain block, migration-count fixtures.

## Operator ratifications received in this chip session

`RATIFIED-BY-OPERATOR:` 2026-07-29 — stock LOCATION rule, amending the context pack's
`only_em_estoque` (which named only `CODEMP IN (1,2)`). Operator, verbatim:

> "Uma coisa CODLOCAL é legal 10101 é estoque revenda, 10108 é show room não é pra contar, e
> 10102 é outlet esses que vendem"

Asked whether it binds the filter only or also the stock NUMBER already rendered on the catalog
screen (which reads only `10101` today), the operator chose "filtro e número, os dois". So
sellable stock = `CODEMP IN (1, 2) AND CODLOCAL IN (10101, 10102)` — show room `10108` excluded —
and the same definition serves the catalog's displayed stock and the mirror's `estoque_total`.

Open against this ratification (routed as `REQUEST db-consult` to the hub, 2026-07-29): whether
available = `ESTOQUE - RESERVADO` (what both existing queries already do) or physical `ESTOQUE`
(operator wrote "e temos os reserváveis também"); whether `CODPARC = 0` remains the own-stock
predicate; the live count under the amended rule (the pack's 3.822 was measured with no CODLOCAL
filter, so it is an upper bound); and whether any other CODLOCAL sells. Not blocking — the chip
implements `CODEMP IN (1,2) AND CODLOCAL IN (10101,10102)` with available = `ESTOQUE - RESERVADO`
and leaves the VC-2 number to the hub's live measurement post-merge.

Rows are written AT DISPATCH TIME (core §1 ledger rule) and completed with the result.

| # | phase | role | model / effort | path | prompt | log / output | result | SHA |
|---|---|---|---|---|---|---|---|---|
| 1 | pre-P2 | Investigator / repo map | Claude subagent (Explore, sonnet) | Agent tool, sync | inline brief: map tenant_config / products_mirror writers / sync Q1+Q4 / catalog_page / xlsx parser / MirrorMatcher / OpenAPI+SDK / FE cards + hooks + tests | returned to session context; facts quoted verbatim into the P2 planner prompt (`scratchpad/prompt-p2-planner.md`) | COMPLETED — 13-section `file:line` map | base 554788d5 |
| 2 | P2 | Feature planner (batch, whole chip) | `gpt-5.6-sol` / `medium` | OS-process codex, `--sandbox read-only`, stdin closed, cwd = chip worktree | `scratchpad/prompt-p2-planner.md` | `scratchpad/agent__p2-planner.log` · `scratchpad/agent__p2-planner.last.md` | COMPLETED — 14 slice cards (S0..S13), write-set DAG, contract-satisfiability table, verification map, seam-closure checklist, must-fail design. Copied verbatim to `BATCH-PLAN.md`; chip amendments appended there (A-1..A-6) because the plan predates the closed A7 ruling | base 554788d5 |
| 3 | P2 | `REQUEST db-consult` → hub → MNOS Oracle specialist | — (cross-session event) | hub `local_99feb041` | 4 questions: RESERVADO semantics, `CODPARC = 0`, live count under the amended rule, other selling CODLOCAL | hub reply quoted verbatim into `ADDENDUM-01-codlocal.md` | ANSWERED — A7 closed: subtract RESERVADO (354 products zero out by reservation alone), keep `CODPARC = 0` (defensive), **VC-2 = 2.923 distinct CODPROD**, named whitelist of 8 locations (only 10101 + 10102 sell). Provenance: db-consult MNOS 2026-07-29 via hub, read-only aggregation | main @ aee1a222 |
| 4 | P2 | `REQUEST DR-1` → hub | — (cross-session event) | hub `local_99feb041` | mirror-catalog parity for `xlsx`/`catalogo_cliente` tenants: apply the predicate + live counts on the upload path (~8 files / ~215 lines), or issue a named+dated deferral with a non-numeric UI state | hub ruling quoted into `BATCH-PLAN.md` §A-7 | **GRANTED** — parity approved, no deferral; the pack is corrected to THREE binding read paths (Oracle · mirror-serving-catalog · MirrorMatcher). Conditions: SQL filter before `LIMIT+1` with a >1-page fixture · measure absent-vs-zero on the xlsx writer · must-fail at the CALLER (`routing.Reader`, source `xlsx`). Mirror counter counts against the mirror. VC-2's 2.923 is the LIVE ORACLE number only; the mirror path discharges by agreement within its own population | main @ 66f32125 |
| 5 | P3 | Implement worker — slice S1-SCHEMA | `gpt-5.6-luna` / `high` | OS-process codex, `--sandbox workspace-write`, stdin closed, cwd = chip worktree | `scratchpad/prompt-s1-schema.md` | `scratchpad/agent__s1-schema.log` · `scratchpad/agent__s1-schema.last.md` | **ACCEPTED** at P4 — see the review below | `23b9d428` |
| 6 | P3 | Implement worker — slice S2-CONFIG-DB | `gpt-5.6-luna` / `high` | OS-process codex, `--sandbox workspace-write`, stdin closed, cwd = chip worktree | `scratchpad/prompt-s2-config-db.md` | `scratchpad/agent__s2-config-db.log` · `scratchpad/agent__s2-config-db.last.md` | **ACCEPTED** at P4, committed by the orchestrator (the worker could not) | `29e7cb69` |

## P4 adversarial review — S2-CONFIG-DB (`29e7cb69`)

**ACCEPTED**, with the same treatment S1 got: the worker's report is a claim, not evidence.

1. **`go build ./...` refuted a SECOND time.** The worker again reported the build "blocked by
   VCS stamping" with `-buildvcs=false` passing. Measured from `apps/server_core` with absolute
   caches: **EXIT=0**, unqualified. Two workers, two identical false alarms, both from the codex
   seat and neither reproducible at the orchestrator's. Recorded as an environment signature of
   that seat, NOT as an L0 result — and no `-buildvcs=false` is buried in any lane to silence it,
   because that trades an alarm for blindness.
2. **The commit failure was real, and it is a finding.** `git commit` died on
   `.git/worktrees/chip-vendavel/index.lock: permission denied`. No lock file exists now. In a
   git worktree the real git directory lives under the MAIN checkout
   (`.git/worktrees/<name>/`), outside the sandbox's writable root — so a worker under
   `--sandbox workspace-write` can be denied the index while every file write succeeds. S1's
   worker did commit, so it is not deterministic; the orchestrator commits when it happens.
   Worth reporting to the hub for the profile.
3. **The RED is weaker than S1's and needed the same correction.**
   `repository_test.go:79:9: got.SellableAssortment undefined` is a COMPILE error — it proves the
   field did not exist, not that any assertion discriminates a value. Two mutations, both killed:
   - `Set`'s `ON CONFLICT DO UPDATE` given `only_ecommerce = false` (a source switch quietly
     resetting the operator's rule) → `repository_test.go:186: after source switch OnlyEcommerce
     = false, want true`.
   - `SetSellableAssortment`'s `WHERE tenant_id = $1` replaced by `WHERE $1::text IS NOT NULL`
     (tenant scoping dropped, the classic leak) → killed in THREE independent places: the first
     tenant re-read after the second write (`:160`, `:163`), the source-switch preservation
     (`:180`, `:183`), and the missing-tenant sentinel returning `<nil>` instead of
     `ErrUnknownActiveSource` (`:190`).
   A first attempt at that second mutation (`WHERE $1 IS NOT NULL`, no cast) is NOT counted: it
   died on `SQLSTATE 42P18 could not determine data type of parameter $1`, which is the mutation
   failing to compile rather than the test catching a leak. A mutation that kills itself proves
   nothing.
4. **Counts measured, not read off a tail:** 13 PASS, **0 SKIP**, 0 FAIL across two packages.
   The skip count matters here — these tests call `SkipWithoutTarget` and a fully skipped run
   prints `ok` too.
5. **Diff reviewed for scope:** `Get` reads the three columns in the SAME query, still
   `WHERE tenant_id = $1`, no second round trip and no default-if-missing branch; `Set` is
   untouched; `SellableAssortment` holds plain booleans, which is correct because the columns are
   `NOT NULL DEFAULT` — the value is always known, so a pointer would invent an unknown that
   cannot occur.

## P4 adversarial review — S1-SCHEMA (`23b9d428`), by the orchestrator

**ACCEPTED.** Four checks, three of them against the worker's own report:

1. **The worker's build failure is REFUTED by measurement.** It reported `go build ./...`
   failing with `error obtaining VCS status: exit status 128`. Re-run from `apps/server_core`
   with absolute `GOCACHE`/`GOMODCACHE`: `go build ./...` → **EXIT=0**, and
   `go build -buildvcs=false ./...` → EXIT=0. Environment false alarm at the worker's seat, not
   an L0 failure. A worker's report is never the evidence.
2. **The RED is real but weak on its own** — `sellable_assortment_test.go:9: open
   0083_sellable_assortment_config.sql: file does not exist` names the test and fails, but a
   file-absent RED proves only that the file was missing, NOT that the assertions discriminate
   the ratified VALUES. So the orchestrator ran a mutation: flipped `only_ecommerce` to
   `DEFAULT true` and made `usoprod` `NOT NULL` in the two migrations, re-ran, and both mutants
   were killed by name:
   `missing exact declaration "…only_ecommerce boolean not null default false"` and
   `products_mirror.usoprod must have no DEFAULT and no NOT NULL`. Tree restored, re-run green,
   `git status --porcelain migrations/` empty. The defaults and the nullability are now pinned
   by an instrument that has been shown to go red.
3. **Independent confirmation of the migration count:** the orchestrator's own
   `go run ./cmd/testdb migrate` against a fresh database printed `applied 71 migration(s)` —
   the same 71 both `runner_test.go` assertions were bumped to, measured from the DB side
   instead of the fixture side. Both assertion sites were bumped (`:25` and `:64`).
4. **Scope and hygiene:** no `.gomodcache` at the repo root (only `apps/server_core/.gomodcache`);
   `writer.go` gained the two struct fields and NOT the `upsertSQL` columns, which is correct —
   S4 owns that seam and has its own test.

Reservation recorded, not blocking: `mirror_test.go` asserts through `reflect.FieldByName`. That
is only justified because the test had to compile before the fields existed (a direct reference
would have made the RED a compile error rather than a named failing test). Later slices assert
these fields directly; if any of them keeps using reflection instead, that is test theater and
gets rejected there.

## Test-database provisioning (orchestrator glue, not a slice)

The hermetic lane is the chip's own instrument (profile §4) — distinct from the hub-owned dev
stack (`:8080`/`:5174`), which the chip never touches. Provisioned once, so no worker starts a
container or invents a DSN:

- `npm run harness:pg:up` → session container `mpc-pg-session-3eee515d`, port `55864`,
  `status=ready` (one per checkout, hashed from the path — hub and chip never collide).
- Slice database `mpc_test_deadbeefdeadbeefdeadbeefdeadbeef`. The name is not decorative:
  `testsupport/postgres/target.go:19` rejects anything that does not match
  `^mpc_test_[0-9a-f]{32}$` with `HPG_TARGET_INVALID` — a friendly name like `mpc_chip_vendavel`
  is refused, which is how the first attempt failed.
- `go run ./cmd/testdb migrate` → `applied 71 migration(s)`.
- Workers receive `scratchpad/testdb-env.ps1` (dot-source only). It derives the URL from
  `scripts/.runs/pg-session.json` at run time, so the container password is never written into
  a prompt, a log, or the evidence pack.

`SkipWithoutTarget` makes a skipped run and a green run byte-identical at the tail (profile §11,
vacuous green). Every DB slice brief therefore demands RUN/PASS/**SKIP** counts as a result.

---

## S2B-RENAME — `only_ecommerce` → `only_ecommerce_eligible`

| field | value |
|---|---|
| dispatched | worker `gpt-5.6-luna` / `high`, OS-process, `agent__s2b-rename.log` + `.last.md` |
| base | `2e943795` |
| commit | `f76c1325` (orchestrator committed; see the build claim below) |
| evidence | `evidence/S2B-red.txt`, `evidence/S2B-green.txt` (worker), `evidence/S2B-orchestrator.txt` |

Hub ruling. The name asserted the opposite of the clause it controls: `only_ecommerce` reads as
"only the published ones", while the ratified clause is `NVL(AD_ECOMMERCE,'X') <> 'N'`, which
removes only what the ERP explicitly marked as outside e-commerce and keeps everything undecided.
Measured live, the strict reading would cut 2.923 → 442. The rename lands INSIDE `0083`, which has
not reached main — no `ALTER RENAME`, no new file.

### P4 — what the orchestrator measured, not what the worker reported

The RED discriminates: the migration test failed naming the OLD declaration string verbatim, so
the assertion is bound to the exact declaration rather than to the column's existence.

The slice database was dropped, recreated and re-migrated after the rename — mandatory, because a
stale schema would make the worker's own package read as a code defect. `applied 71 migration(s)`,
the SAME count as before the rename: independent confirmation, by measurement rather than by
assertion, that the count fixture is untouched.

`go test ./internal/modules/tenant_config/... ./migrations -count=1 -v` → **33 RUN / 33 PASS /
0 SKIP / 0 FAIL**. `SKIP=0` carries the weight: without `MPC_TEST_DATABASE_URL` these tests skip
silently, and a fully skipped run is byte-identical at the tail to a fully green one (profile §11).

Sweep at the CALLER, anchored, asserting ABSENCE — `git grep -nE 'only_ecommerce($|[^_])|OnlyEcommerce($|[^E])'`:
zero in `apps/server_core`, `apps/web`, `packages/`, `contracts/`; zero `_eligible_eligible` (the
prefix trap did not fire); new-name count **19**, identical to the pre-image count of the old name
taken from `git HEAD` BEFORE the worker ran — one-to-one, nothing lost or doubled. Five survivors
are deliberate: `BATCH-PLAN:1066` names the rename and must quote the dead name, and
`DISPATCH-LEDGER:70,71,103,106` are verbatim history of the S1/S2 runs. Rewriting those would
falsify the trail.

### The rename destapou prosa falsa, corrigida in the same commit

`BATCH-PLAN:60` still carried `m.ad_ecommerce = 'S'` — the formula the DR-3 revision revoked — and
`CONTEXT-PACK:23` carried both the dead name and the revoked semantics. A worker reads the plan,
not the migration: that clause would have entered S7 as the strict reading and cut the assortment
to 442. The S3 brief carried the same falsehood and was corrected BEFORE dispatch.

### Finding — third false build claim from the same seat

The worker reported `go build ./...` failing with `error obtaining VCS status: exit status 128`
and **withheld its commit on that basis**. Measured in the same worktree, `cd apps/server_core`,
absolute caches: `EXIT=0`. Third identical occurrence (S1, S2, S2B) — a signature of the codex
seat, not a repo fact. No `-buildvcs=false` was buried in the lane: that trades an alarm for a
blindness (hub, binding). The real cost this time was a green slice left uncommitted, so the S3
brief now states that this observation must not withhold a commit.

### Finding — the sweep instrument is partial, cause NOT established

The Grep tool rooted at the worktree's ABSOLUTE path returned 5 matches and missed all 19 code
occurrences; `git grep` from inside the worktree returns all 19. The `.gitignore '.claude/*'`
hypothesis was tested and REFUTED (`rg` from inside the worktree finds them). Recorded as an
instrument constraint, not a diagnosis. Operationally: rename sweeps in this chip use `git grep`
from inside the worktree. A partial sweep is worse than none — it has the shape of a complete one,
and a gate seat using the wrong tool would report "zero occurrences" with full confidence.

---

## S3-XLSX — optional `USOPROD` / `AD_ECOMMERCE` on the upload path

| field | value |
|---|---|
| dispatched | worker `gpt-5.6-luna` / `high`, OS-process, `agent__s3-xlsx.log` + `.last.md` |
| base | `a94694c2` |
| commit | `8f825a7c` |
| evidence | `evidence/S3-red.txt`, `evidence/S3-green.txt` (worker), `evidence/S3-orchestrator.txt` |

Production diff is correct: both fields flow through `optionalCell`, the two required lists
(`Parse`, `ParseLenient`) are untouched, and `usoprod`/`ad_ecommerce` reach the stage table, the
`CopyFrom` column list, the insert and the `ON CONFLICT DO UPDATE SET`. `nullableStringValue` is a
justified thin wrapper — the existing `nullableString` takes a value, not a pointer.

### P4 — three mutations, three kills

The slice is entirely honest-unknown, so presence assertions would be worthless. Injected and
killed: (M1) parser emitting `""` for an absent column → caught at `parser_test.go:81`;
(M2) `nullableStringValue` persisting `""` instead of SQL NULL → caught at
`mirror_repository_integration_test.go:37`, which matters because the round-trip is what proves
SQL NULL rather than a Go nil that a nil-preserving scan would yield anyway; (M3) adding
`USOPROD` to the required list → four PRE-EXISTING tests fail with `MISSING_REQUIRED_COLUMN`, so
"legacy workbooks keep importing" is guarded rather than merely true by inspection. All reverted,
tree clean against `8f825a7c`, `erp_import/...` green.

### Finding — a pre-existing clock-skew flake, NOT caused by S3

The worker reported `TestGetImportChainCountsCurrentQueueAcrossInstallations` failing as
"unrelated timing ... reproduced in isolation". The second half is FALSE and was refuted: it
PASSES in isolation and PASSES the full 28-test package; it fails only in the `erp_import/...`
tree run, and non-deterministically (run 1 FAIL, run 2 PASS at the same tip).

Root cause: `query_repository.go` selects `statement_timestamp() AS queue_read_at` — the
POSTGRES clock, inside the container — while the test brackets it with host-side `time.Now()`
(`chain_query_repository_integration_test.go:74-85`). It asserts a container clock reading falls
inside a host clock window of ~10ms. The observed miss was ~77µs, and the precision tell is in
the output: 6 decimals (PG microsecond) against 7 (Windows 100ns).

The test is NOT in `8f825a7c` — S3 neither wrote nor touched it. Latent defect, surfaced by added
package load. Out of the S3 write-set, so reported to the hub and not fixed here. It will bite
the P5 ladder and any gate seat that reads a single red tree run as attributable.

### Plan rot corrected in the same commit (hub standing order)

Formula revocation is a write-set change. Sweeping the pack for every predicate DR-2/DR-3/A-14
touched found the S6 slice card still carrying ALL THREE dead clauses in its `done_criteria` —
`(p.USOPROD IS NULL OR = 'R')`, `(stock.sellable_qty IS NULL OR > 0)` and
`(p.AD_ECOMMERCE ... = 'S')` — superseded by amendments ~750 lines further down the same file.
S6's worker reads the card, not the amendment. Corrected to the binding live forms with the
supersession recorded inline and a pointer to the A-14 asymmetry table. `BATCH-PLAN:1086` also
printed the revoked formula under the heading "Original reasoning (the principle, unchanged)" —
the principle is unchanged, the formula is not; struck through and marked.
