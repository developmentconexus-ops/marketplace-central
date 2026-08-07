# Lane: duplication

> I have run something really deep into MetalDocs and I am changing the way I code there to move to
> something more professional towards issues, PRs, PR review, CodeRabbit mechanical full validation
> and so much more. For that I had to identify every error in my code, my platform, to improve it and
> to create this full validation. I want to run it here as well so we move on the same path, this way
> it gets so much harder to send bad PRs.

Calibration: solid professional level, not Google-tier. Operative sentence: **"this way it gets so
much harder to send bad PRs"** — the target is a mechanism, not a cleanup.

## Findings

| ID | class | finding | evidence | scale |
|---|---|---|---|---|
| DUP-1 | duplication/drift | Four-copy contract seam: domain → DTO → OpenAPI → SDK, no compiler link, `GOV_API_SDK_SPLIT` checks same-commit only, never agreement | `scripts/harness/Policy.psm1:460` `if ($apiChanged -xor $sdkChanged)`; 239 OpenAPI schemas (`contracts/api/marketplace-central.openapi.yaml:3517-8574`, counted); 172 SDK interfaces / 2595 lines in one file (`packages/sdk-runtime/src/index.ts`, established fact 5, re-confirmed by `wc -l`); sampled `ListingReadModel` field-by-field, see heaviest #1 | 239 schemas, 172 interfaces, re-confirms D-48 |
| DUP-2 | duplication | `Money` struct independently defined 5 times; 4 of the 5 are byte-identical `{Amount string \`json:"amount"\`; Currency string \`json:"currency"\`}` | `apps/server_core/internal/modules/connectors/domain/money.go:10`, `apps/server_core/internal/modules/listings/domain/read_model.go:83` (Currency typed `PriceCurrency` alias of string), `apps/server_core/internal/modules/market/domain/market.go:13`, `apps/server_core/internal/modules/pricing/domain/decimal.go:12`; canonical 5th is `apps/server_core/internal/kernel/exact/money.go:45` (`{amount Decimal; currency Currency}`) | 5 struct defs, 4 identical shape |
| DUP-3 | gap | `internal/kernel/exact` (the purpose-built exact-decimal package: `Decimal`, `Money`, no float64 constructor by design) has **zero production callers** — only its own tests and `kernel/fact` tests import it | `grep -rln 'kernel/exact"' apps/server_core/internal --include=*.go` → 4 files, all `_test.go` | 0 production call sites for a canonical type that exists specifically to fix this class of bug |
| DUP-4 | duplication | Order/listing money values parsed with lossy `strconv.ParseFloat` instead of the exact-decimal path, at 12 call sites, several round-tripping a value that was already computed exactly via `big.Rat`/string back through a binary float | `apps/server_core/internal/modules/orders/adapters/pricingtax/reader.go:276,309` (line 309 parses `pricingdomain.FormatRatHalfUp(r, 2)` — an already-exact 2dp string — back into a float64); `apps/server_core/internal/modules/pricing/application/batch_orchestrator.go:286`; 9 more sites listed in commands below | 12 call sites, `grep -c` |
| DUP-5 | duplication | Profitability/margin formula (`total − comissão − frete − ICMS_saída − DIFAL − PIS/COFINS − custo + restituição_ST`, then `margem / total`) implemented independently twice, in two different numeric representations, with a different `×100` convention | `apps/server_core/internal/modules/orders/domain/order_decomposition.go:118-163` (`BuildProfitability`, `*float64` pointers, `pct := margem / *in.Total` — fraction, no ×100, per comment "FE formatPercent multiplies by 100") vs `apps/server_core/internal/modules/pricing/domain/decompose.go:141-286` (`Decompose`, `big.Rat`/string, `pct.Mul(pct, cem)` — ×100 baked in) | 2 independent implementations, 0 shared code between them |
| DUP-6 | duplication | pgx row-scan boilerplate: every repository hand-writes `for rows.Next() { rows.Scan(...) }`; pgx v5's built-in generic collectors are never used | `grep -rn "for rows.Next()" apps/server_core/internal --include=*.go \| grep -v _test.go` → 82; `grep -rn "defer rows.Close()"` → 75; `grep -n "jackc/pgx" apps/server_core/go.mod` → `v5.7.6` (has `pgx.CollectRows`/`pgx.RowToStructByName`); `grep -rln "pgx.CollectRows\|pgx.RowToStructByName" apps/server_core/internal --include=*.go` → 0 hits | 82 hand-written scan loops, 0 use of the stdlib-adjacent generic helper |
| DUP-7 | duplication | 19 module-local `write*Error`/`writeError` HTTP-error-envelope helper functions, each hand-copied per module instead of one shared helper; 3 are byte-identical | `grep -rn "^func write.*Error\|^func (.*) write.*Error" apps/server_core/internal/modules --include=*.go \| grep -v _test.go` → 19 hits across `dashboard`, `listings`, `integrations` (×2), `catalog` (×2), `market` (×3), `orders` (×3), `erp_import`, `mutations` (×3), `pricing`, `profitability`; the identical-signature subset `func write.*Error(w http.ResponseWriter, status int, code, message, key string)` at `dashboard/transport/http_handler.go:59`, `integrations/transport/run_read_handler.go:124`, `listings/transport/http_handler.go:326` is byte-for-byte the same 5-line body (`details := map[string]any{}; if key != "" {...}; apierror.Write(...)`) | 19 wrapper functions, 3 byte-identical (verified by diff-by-eye of extracted bodies) |
| DUP-8 | duplication | Per-consumer anti-corruption-layer adapter that shims `marketplaces.Service` into a module-local policy port is hand-written independently per consumer; 2 of the 3 do the *same* translation (`Policy → {MinMarginPercent}`) with different nil-handling and receiver conventions | `apps/server_core/internal/modules/listings/adapters/marketplaces/policy_reader.go` (21 lines, value receiver, no nil guard) vs `apps/server_core/internal/modules/mutations/adapters/marketplaces/policy_reader.go` (39 lines, pointer receiver, explicit `r == nil \|\| r.service == nil` guard, extra `HasPricingPolicy` method) — `diff` run, shown in commands; third instance `apps/server_core/internal/modules/pricing/adapters/marketplace/reader.go` translates to a richer `BatchPolicy` shape (legitimately different, not counted as duplicate) | 2 near-duplicate adapters (same filename, same source type, same translated field), 1 legitimately distinct |
| DUP-9 | idiom | 6 `big.Rat`-based decimal libraries worth of hand-rolled exact arithmetic outside `kernel/exact`: 15 files construct `big.Rat` directly (`big.NewRat`/`new(big.Rat)`) instead of going through a shared `Decimal` type | `grep -rln "big.NewRat\|new(big.Rat)" apps/server_core/internal --include=*.go \| grep -v _test.go` → 15 files, mostly `pricing/domain/*.go` | 15 files |
| DUP-10 | duplication | 83 hand-written `strings.TrimSpace(x) == ""` required-field checks scattered across transport handlers, no shared "require non-empty" validator | `grep -rn 'strings.TrimSpace(.*) == ""' apps/server_core/internal/modules --include=*.go \| grep -v _test.go` → 83 | 83 call sites (low severity — one-liner, but 83× the same one-liner) |
| DUP-11 | duplication (test-only, excluded from main count) | PowerShell `Assert-True` helper function redefined 6 times, only in `*.tests.ps1` files | `grep -rl "^function Assert-True"` → `scripts/tests/{governance-contracts,governance-drift,harness-environment,harness-execution,postgres-contract,postgres-lifecycle}.tests.ps1` | 6 test files — excluded from headline count per the test-file exclusion rule below, reported for completeness |

## The five heaviest, with detail

### 1. The four-copy contract seam — re-measured, and now field-level, not just line-count

Established fact 5 (SDK: 2595 lines / 172 interfaces, one file `packages/sdk-runtime/src/index.ts`)
and fact 6 (`GOV_API_SDK_SPLIT` same-commit-only) are re-confirmed:

```
$ wc -l packages/sdk-runtime/src/index.ts                    → 2595
$ grep -c "^export interface" packages/sdk-runtime/src/index.ts → 172
$ sed -n '3517,8574p' contracts/api/marketplace-central.openapi.yaml \
    | grep -c "^    [A-Za-z_][A-Za-z0-9_]*:$"                 → 239   (schema count)
```

`scripts/harness/Policy.psm1:460`: `if ($apiChanged -xor $sdkChanged) { … GOV_API_SDK_SPLIT … }` —
literally an XOR on "did this path change", never a shape comparison.

To move past the established fact into a real measurement, I traced one shape end-to-end:
`ListingReadModel` (Go domain `apps/server_core/internal/modules/listings/domain/read_model.go:122-152`,
OpenAPI `contracts/api/marketplace-central.openapi.yaml:3588-3653`, SDK
`packages/sdk-runtime/src/index.ts:380-405`). Two concrete, live divergences found, not just
theoretical risk:

- **`listing_type` / `price` nullability disagrees across all three.** Go: `ListingType *ListingType`
  and `Price *Money` — pointers, nil is a real, produced value. OpenAPI: `listing_type`/`price` are in
  the `required` array with no `nullable: true` on the ref (`:3590,3607-3617`) — non-nullable by
  omission. SDK: `listing_type: ListingType` and `price: ListingMoney` — non-null, no `| null`
  (`index.ts:387,390`). A consumer generated from the OpenAPI/SDK side would not defend against the
  null the Go side can actually emit.
- **`market_signal`/`signal_status` presence disagrees.** Go: plain (non-pointer) `SignalStatus`
  field and `*MarketSignal` pointer, both **always serialized** (no `omitempty`,
  `read_model.go:150-151`). OpenAPI: both properties exist but neither is in `required`
  (`:3648-3653`) — modeled as absent-or-present. SDK: both marked optional with `?`
  (`index.ts:403-404`) — matches OpenAPI, not Go. Two of three layers agree with each other and
  disagree with the actual producer.

This is a sample of one shape, chosen because it is read-heavy and screen-facing; it is not an
exhaustive audit of 239 schemas. It demonstrates the class D-48 already named (byte-divergent
transcription, unchecked by tooling) is not hypothetical — it is present in the first shape checked.

### 2. Money has 5 independent hand-written definitions; the one built to be exact is orphaned

```
$ grep -rn "type Money struct" apps/server_core/internal --include=*.go | grep -v _test.go
apps/server_core/internal/kernel/exact/money.go:45        {amount Decimal; currency Currency}   <- canonical, unexported fields, constructor-guarded
apps/server_core/internal/modules/connectors/domain/money.go:10   {Amount string; Currency string}
apps/server_core/internal/modules/listings/domain/read_model.go:83 {Amount string; Currency PriceCurrency}
apps/server_core/internal/modules/market/domain/market.go:13      {Amount string; Currency string}
apps/server_core/internal/modules/pricing/domain/decimal.go:12    {Amount string; Currency string}
```

4 of 5 are the same shape, hand-copied. None import `kernel/exact.Money`. And `kernel/exact` itself —
package doc at `apps/server_core/internal/kernel/exact/decimal.go:1-5`: *"There is no constructor from
float64 anywhere in this package, and that is the point"* — has zero production callers:

```
$ grep -rln 'kernel/exact"' apps/server_core/internal --include=*.go
apps/server_core/internal/kernel/exact/decimal_test.go
apps/server_core/internal/kernel/exact/money_test.go
apps/server_core/internal/kernel/fact/combine_test.go
apps/server_core/internal/kernel/fact/knowledge_test.go
```

All 4 hits are test files. This is the same shape of finding as established fact 11
(`kernel/fact.Map`/`Combine2`/`provenance.Derived` orphaned) but for a different package — worth
flagging as its own item since fact 11 does not list `exact.Money`/`exact.Decimal`.

Consequence, measured, not inferred: production code parses money through lossy `strconv.ParseFloat`
at 12 sites (`DUP-4` above), including one that round-trips an already-exact decimal string back
through a binary float (`orders/adapters/pricingtax/reader.go:309`:
`strconv.ParseFloat(pricingdomain.FormatRatHalfUp(r, 2), 64)`), and 15 files hand-construct
`big.Rat` directly instead of using the one type built to make that safe (`DUP-9`).

### 3. The margin formula is implemented twice, disagreeing on units

`orders/domain/order_decomposition.go:118-163` (`BuildProfitability`) and
`pricing/domain/decompose.go:141-286` (`Decompose`) both implement:

```
margem = total − comissão − frete − icms_saída − difal − pis_cofins − custo + restituição_st
margem_pct = margem / total
```

but:

| | orders | pricing |
|---|---|---|
| numeric type | `*float64` | `*big.Rat` → 2dp string |
| `margem_pct` | raw fraction (0.18), comment says FE multiplies by 100 | `pct.Mul(pct, cem)` — ×100 baked in |
| unknown handling | per-field nil checks, appends component name to `ComponentesDesconhecidos` | same idea, different code |

Zero shared code between the two functions — no shared helper, no common package. Two people (or two
generations of the same agent) solving the identical arithmetic problem, with a real unit mismatch
(fraction vs percent) that would silently corrupt a value if either caller ever called the wrong one
for the wrong context.

### 4. pgx row-scan boilerplate — 82 hand-written loops, 0 use of the driver's own generic collector

```
$ grep -rn "for rows.Next()" apps/server_core/internal --include=*.go | grep -v _test.go | wc -l
82
$ grep -rn "defer rows.Close()" apps/server_core/internal --include=*.go | grep -v _test.go | wc -l
75
$ grep -n "jackc/pgx" apps/server_core/go.mod
	github.com/jackc/pgx/v5 v5.7.6
$ grep -rln "pgx.CollectRows\|pgx.RowToStructByName" apps/server_core/internal --include=*.go
(no output)
```

pgx v5.7.6 ships `pgx.CollectRows`/`pgx.RowToStructByName` specifically to eliminate this loop. It is
imported and available (it's the driver already in `go.mod`) but never used — every one of the 82
call sites hand-rolls scan-loop, error-check, append.

### 5. 19 module-local HTTP-error-writer wrappers around one shared function, 3 of them byte-identical

```
$ grep -rn "^func write.*Error\|^func (.*) write.*Error" apps/server_core/internal/modules --include=*.go | grep -v _test.go | wc -l
19
```

All 19 ultimately call the shared `apierror.Write` (`apps/server_core/internal/platform/apierror/apierror.go:25`)
— so the *sink* is not duplicated, only the *wrapper*. But the wrapper itself is hand-copied:
`dashboard/transport/http_handler.go:59`, `integrations/transport/run_read_handler.go:124`, and
`listings/transport/http_handler.go:326` are the identical 5-line function
(`func write*Error(w http.ResponseWriter, status int, code, message, key string) { details := map[string]any{}; if key != "" { details["key"] = key }; apierror.Write(w, status, code, message, details) }`),
verbatim, in 3 different modules. This is the exact class the "error surface unification" backlog
item (memory: `error-surface-unification-backlog.md`; also `chip-error-unify-closed.md` — a related
but distinct duplication was already closed in that chip) still has open: the unification reached the
sink, not the call site.

## What is actually fine

- **`apierror.Write`** (`apps/server_core/internal/platform/apierror/apierror.go:25`) is a real,
  single, shared HTTP-error envelope writer that every one of the 19 module wrappers ultimately calls
  — the envelope shape itself is not duplicated, only the thin per-module dispatch to it.
- **No raw `fetch()` bypassing the SDK on the frontend.** Initial grep for `fetch(` hit 22 files, but
  every hit was `refetch(`/`Prefetch(` (react-query), not a literal `fetch()` call:
  `grep -rn "[^a-zA-Z]fetch(" apps/web/src --include=*.tsx --include=*.ts | grep -v test | grep -v "refetch(\|Prefetch(\|prefetch("` →
  0 hits. All 77 files using `useQuery`/`useSWR`/`@tanstack` go through the generated/hand-written
  SDK client, not ad hoc fetch.
- **FE query-hook layer is disciplined, not copy-pasted.** `apps/web/src/pages/vinculos/useVinculosQueue.ts`
  uses a shared `@marketplace-central/web-query` package (`invalidateAfterMutation`, `QUERY_STALE_TIME`)
  and page-local `vinculosQueryKeys` rather than hand-rolling cache invalidation per hook.
- **Only one live vendor integration (Mercado Livre).** `find apps/server_core/internal -iname "*mercado*"`
  vs no second vendor directory — the "repeated mapper/adapter code between vendors" candidate the
  brief asked me to check does not exist as a finding: there is nothing to duplicate against. VTEX is
  already dead per prior memory (`varredura-maximo-global-instrumentos.md`), re-confirmed by absence
  in this sweep.
- **The ICMS matrix is now a single producer.** `icms_matrix_mirror` had one reader (Oracle) and one
  writer (mirror) but no scheduler, so it sat at 0 rows (prior finding, `plan-gate-composition-site.md`).
  `apps/server_core/internal/composition/root.go:711-729` shows the scheduler now wired
  (`icmsMatrixScheduler`, `syncdomain.EntityICMSMatrix`, 24h cadence) — the gap is closed, not open.
  I did not re-query row counts live (no DB session in this lane); the code-level producer/consumer
  wiring is unambiguous and matches the comment's own account of the fix.
- **The DIFAL override mechanism is real, wired, and audited** — contradicts the prior memory note
  ("4 override mechanisms, zero in use", `varredura-maximo-global-instrumentos.md`). `DifalOverride`
  (`pricing/domain/difal.go:31`) has a persistence path (`calc_repository.go:145-168`,
  `UpsertDifalOverride`), an audit gate (`calc_service.go:35`,
  `overrideThresholdPP = big.NewRat(49, 1000)`), and a read path that reflects it
  (`calc_repository.go:319-fnbelow`, `scanDifalRate`). This is not one of 3-4 tariff-table copies
  either: the three ICMS/tariff-shaped tables that exist today
  (`pricing_tariff_defaults` 0068, `icms_aliquota_interna` 0094, `icms_matrix_mirror` 0094) are three
  different facts (channel-fee defaults, our own legal per-UF rate, mirrored ERP cell) with three
  different owners, not three copies of one table. **The prior "3 copies / 4 unused overrides" number
  is stale** — re-measurement in this sweep does not reproduce it. Flagging per the brief's own
  instruction ("a lane that contradicts one with better evidence should say so loudly").
- **Only 85 migrations, no obviously-duplicated table definitions found** beyond the tariff-table
  triage above; did not do a full pairwise column-diff of all 85 (out of budget), so this is a
  lighter-weight "fine" than the others — flagged as `unverified` at full depth.

## Unverified / needs judgment

- The four-copy seam's divergence rate is measured on **one sampled shape** (`ListingReadModel`), not
  all 239 schemas. Two real divergences were found in that one shape; extrapolating a repo-wide
  divergence percentage from N=1 would be fabrication. A generator/diff tool (already approved per
  established fact 7 — `oapi-codegen`/`openapi-typescript`, not yet installed) is the only honest way
  to get the full-repo number.
- `icms_matrix_mirror` row count: code-level wiring says the gap is closed; I did not run a live query
  against Postgres to confirm rows exist today. `unverified` at the data level, confirmed at the code
  level.
- DUP-8 (policy-reader adapter duplication): only 3 instances found repo-wide via
  `GetPricingPolicyForInstallation`/`ListPoliciesByIDs` grep. Did not exhaustively check every
  cross-module port adapter for the same shim pattern — plausible more exist under different method
  names; not counted without a name to grep for.
- Migrations (85 files, 3128 lines) were not fully pairwise-compared for duplicated DDL beyond the
  fiscal-table triage. `unverified`.
- DUP-10 (83 `TrimSpace == ""` checks): counted but not judged — could be idiomatic Go (no stdlib
  "require non-empty" helper exists to reuse) rather than a defect; reported as a size, not a verdict.

## Exclusions (what I did not count, and why)

- **`scripts/.runs/**`** — 7 full snapshot copies of `scripts/` plus `node_modules`, confirmed
  gitignored (`git ls-files scripts/.runs | wc -l` → 0; `git check-ignore -v` confirms
  `.gitignore:14:scripts/.runs/`) and untracked. These are harness run artifacts, not source; excluded
  from every PowerShell count in this report by filtering `grep -v '\.runs/'` after `git ls-files`.
- **`*.test.ts`/`*_test.go`/`*.tests.ps1`** — excluded from all clone/duplication counts per the lane
  brief. The one place this mattered concretely: `Assert-True` (PowerShell) is defined 6 times but
  only in `*.tests.ps1` files (DUP-11) — reported separately, not folded into the main duplication
  count, and not treated as a defect the way DUP-2/DUP-6/DUP-7 are.
- **Generated code**: none found to exclude — `oapi-codegen`/`openapi-typescript` are approved but not
  yet installed (established fact 7), so there is no generated SDK/DTO code in the tree today. The
  entire SDK is hand-written (established fact 5), which is why DUP-1 exists at all.

## Commands run

```
grep -n "D-3[0-9]\|D-4[0-9]\|D-5[0-9]\|duplicat\|copy\|copies\|four-copy\|contract seam" .mnfs/HARNESS-DEBTS.md
wc -l packages/sdk-runtime/src/index.ts
grep -c "^export interface" packages/sdk-runtime/src/index.ts
sed -n '3517,8574p' contracts/api/marketplace-central.openapi.yaml | grep -n "^    [A-Za-z_][A-Za-z0-9_]*:$" | wc -l
grep -rln "type Product struct\|type Listing struct\|type Order struct" apps/server_core/internal/modules/*/domain/
grep -rn "type Money struct" apps/server_core/internal --include=*.go | grep -v _test.go
grep -rln "exact\.Money" apps/server_core --include=*.go | grep -v _test.go
grep -rln 'kernel/exact"' apps/server_core/internal --include=*.go
grep -rn "strconv.ParseFloat" apps/server_core/internal --include=*.go | grep -v _test.go
grep -rln "big.NewRat\|new(big.Rat)" apps/server_core/internal --include=*.go | grep -v _test.go
grep -rn "^func BuildProfitability\|^func.*Decompose\|^func.*Profitability" apps/server_core/internal/modules/orders/domain/order_decomposition.go apps/server_core/internal/modules/pricing/domain/decompose.go
grep -rn "for rows.Next()" apps/server_core/internal --include=*.go | grep -v _test.go | wc -l
grep -rn "defer rows.Close()" apps/server_core/internal --include=*.go | grep -v _test.go | wc -l
grep -n "jackc/pgx" apps/server_core/go.mod
grep -rln "pgx.CollectRows\|pgx.RowToStructByName" apps/server_core/internal --include=*.go
grep -rn "^func write.*Error\|^func (.*) write.*Error" apps/server_core/internal/modules --include=*.go | grep -v _test.go
diff apps/server_core/internal/modules/listings/adapters/marketplaces/policy_reader.go apps/server_core/internal/modules/mutations/adapters/marketplaces/policy_reader.go
grep -rln "GetPricingPolicyForInstallation" apps/server_core/internal --include=*.go | grep -v _test.go
grep -rn "override\|Override" apps/server_core/internal/modules/pricing --include=*.go | grep -v _test.go
grep -rhoE "CREATE TABLE( IF NOT EXISTS)? [a-zA-Z_.\"]+" apps/server_core/migrations/*.sql | grep -iE "icms|tarif|aliquot"
sed -n '695,730p' apps/server_core/internal/composition/root.go
git ls-files scripts/.runs | wc -l
git check-ignore -v scripts/.runs/155085a389a44f3ca06339a018596c19/snapshot/run-server.ps1
git ls-files '*.ps1' '*.psm1' | grep -v '\.runs/' | xargs grep -ohE "^function [A-Za-z0-9_-]+" | sort | uniq -c | sort -rn
grep -rl "^function Assert-True" scripts -r
grep -rn 'strings.TrimSpace(.*) == ""' apps/server_core/internal/modules --include=*.go | grep -v _test.go | wc -l
grep -rn "[^a-zA-Z]fetch(" apps/web/src --include=*.tsx --include=*.ts | grep -v test | grep -v "refetch(\|Prefetch(\|prefetch("
grep -rl "useQuery\|useSWR\|react-query\|@tanstack" apps/web/src --include=*.tsx --include=*.ts | wc -l
grep -n "^components:\|^  schemas:" contracts/api/marketplace-central.openapi.yaml
```
