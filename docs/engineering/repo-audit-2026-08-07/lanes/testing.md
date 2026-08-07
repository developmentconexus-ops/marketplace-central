# Lane: testing

> I have run something really deep into MetalDocs and I am changing the way I code there to move to
> something more professional towards issues, PRs, PR review, CodeRabbit mechanical full validation
> and so much more. For that I had to identify every error in my code, my platform, to improve it and
> to create this full validation. I want to run it here as well so we move on the same path, this way
> it gets so much harder to send bad PRs.

Calibration: solid professional level, not Google-tier. Success condition: a mechanism that makes a
bad change hard to land — not a cleaner test suite. Method: `Documents\MetalDocs\docs\engineering\repo-audit-playbook.md`.

Repo state: `main` @ `1473e863`, working tree clean at lane start. Established facts 1-14 in
`docs/engineering/repo-audit-2026-08-07/PHASE-0.md` are treated as given and re-measured only where
noted.

## Findings

| ID | class | finding | evidence | scale |
|---|---|---|---|---|
| T-1 | gap | `internal/**` — the single largest Go test tree (353 files, ~1800 `func Test*`) including all 21 legacy modules — is reached by NO npm/harness-invocable lane. Only manual `bash scripts/arch-gate.sh` touches it, and that script has zero `npm`/CI wiring. | `git ls-files 'apps/server_core/internal/**_test.go' \| wc -l` → 353; `grep -c '^func Test' <files>` sums to 1800; `grep -n "arch-gate" package.json` → 0 hits (re-confirms D-51); `npm run harness:unit` output shows only `apps/server_core/tests/unit` for Go | 353 files / ~1800 funcs, 100% of `internal/` |
| T-2 | gap | `TestModuleBoundaryADR023` (the ADR-023 module-boundary detector) has exactly one test function: it scans the real repo and asserts violations. No fixture, no positive control (proof it returns 0 on a clean tree), no must-fail proof (an injected violation it can catch). Contrast with its siblings in `internal/arch` and `internal/platform/archguard`, both of which have exactly that. | `apps/server_core/internal/composition/module_boundary_arch_test.go` — 217 lines, 1 `func Test`; ran `go test ./internal/composition/...` → `FAIL … 234 violation(s)` | 1 detector, 0 fixture tests, 234 unverified-by-fixture violations |
| T-3 | gap | 9 FE test files (111 test cases) are structurally unreachable: `packages/sdk-runtime` (7 files/86 tests, own `vitest.config.ts` + own `test` script), `packages/feature-classifications` (1/15), `packages/feature-inventory` (1/10). Root `npm test`/`test:run` only delegate to the `web` workspace (`package.json:13-14`). All 3 packages pass when run directly. | `cat package.json` lines 13-14: `"test": "npm run test --workspace @marketplace-central/web"`; ran `npx vitest run` in each package dir → 86/15/10 passed | 9 files / 111 tests, 0 reachable from root |
| T-4 | hazard | Two security/data-integrity boundary invariants have their SOLE test evidence authored in the same commit that introduced the mechanism, with no independent/cold verification: (a) catalog RLS tenant isolation, (b) catalogingest tenant fail-closed guard (today's commit). | `git log --follow -- .../catalog_rls_role_test.go` → only commit `7e3dcc47` (same commit adds `migrations/0098_catalog_app_role.sql`); `git show --stat 47a76837` → adds `requireTenantConfigured` AND `TestRequireTenantConfigured` in cmd/catalogingest in one commit | 2 boundary invariants, 0 with independent evidence |
| T-5 | hazard | All 15 Mercado Livre HTTP adapter test files use `httptest.NewServer` with inline, hand-authored JSON fixtures written by the same author/change as the parser under test — no golden fixture captured from a real ML response, no schema check against ML's published contract. This is exactly the class AGENTS.md:14 warns "provider payloads remain at adapters" is meant to isolate, but nothing here proves the isolated fixture still matches the provider. | `grep -rln "httptest.NewServer" apps/server_core/internal/modules --include="*_test.go"` → 15 files; read `prices_reader_test.go:24-30` — literal JSON body written next to the assertions it feeds | 15 files, 0 golden/live-captured fixtures |
| T-6 | gap | `npm run harness:integration` prints no RUN/PASS/SKIP/FAIL counts on success — only `container=`, `resource_count=`, `status=passed`. A run where every test was silently skipped would print byte-identical output to a fully-green run. (Re-measurement of ledger D-26.) | ran `npm run harness:integration`; full stdout (16 lines) has no test-count line; `status=passed` is the only signal | 1 lane, 0 count signal on success |
| T-7 | gap | 13 PowerShell test files under `scripts/tests/` (`*.tests.ps1`, ~250+ `Assert-True`/`throw` checks) covering governance contracts, governance drift, postgres lifecycle, hermetic lanes, harness environment/execution/aliases — are invoked by NO automated lane. `scripts/harness/Policy.psm1:228` only *excludes* `scripts/tests` from governance scanning; nothing runs them. Manual `pwsh -File` per file is the only way any of these are ever exercised. | `grep -rln "scripts/tests" scripts/*.ps1 scripts/harness/*.psm1 package.json` → only the exclusion regex; `ls scripts/tests/*.ps1 \| wc -l` → 13 | 13 files, ~250 assertions, 0 automated invocations |
| T-8 | gap | `apps/server_core/migrations/*_test.go` (15 files, 26 `func Test`) and `cmd/catalogingest/main_test.go` (1 file, 3 funcs, includes the T-4(b) guard) are reached by no lane at all — not even `arch-gate.sh`'s `go test ./internal/...`, since `migrations/` and `cmd/` are outside `internal/`. Both pass when run directly. | `go test ./migrations/...` → `ok … 18.038s`; `go test ./cmd/...` → `ok … apps/server_core/cmd/catalogingest 2.561s` | 16 files / 29 funcs, 0 lanes |
| T-9 | idiom | `TestRouterRegistersAllFoundationEndpoints` (the only Go test proving all foundation HTTP routes exist) asserts presence only — "not 404" — never the actual response body/status for a correctly-handled request. A route registered for the wrong method would still pass. | `apps/server_core/tests/unit/router_registration_test.go:172-179` | 1 test, 10 routes checked by presence only |
| T-10 | drift | `apps/web/vitest.config.ts:14` includes `packages/feature-products` by an exact literal filename (`.../CatalogPage.test.tsx`), not a glob, while the two adjacent entries (`web-query`, `ui`) use globs. Today there is exactly one test file in `feature-products`, so this is currently latent, not active — but any second test file added to that package is silently never run. Re-confirmation of ledger memory item B-7. | `apps/web/vitest.config.ts:10-17`; `git ls-files 'packages/feature-products/**' \| grep test` → exactly 1 file | 1 package, 1 landmine, 0 files silently dropped today |
| T-11 | gap | 6 `t.Skip` sites across 3 files gate all live-Oracle validation behind `MPC_ORACLE_LIVE_TEST=1`, which no invocable lane ever sets, plus a `//go:build !cgo` fallback that unconditionally double-skips on a non-cgo host. Re-confirmation of ledger D-8. | `apps/server_core/internal/modules/internal_read/adapters/oracle/reader_live_nocgo_test.go:10-15`, `reader_live_test.go:1,20-22`, `sync_integration_test.go:36` | 3 files, 6 skip sites, 0 executions in any wired lane |
| T-12 | gap | The root `vitest.config.ts` exists (`include: ['**/*.test.ts','**/*.test.tsx']`) but no `npm` script ever passes it to `vitest`; the only root `test`/`test:run` scripts delegate straight into the `web` workspace's own config. | `package.json:13-14`; `vitest.config.ts` (root) | 1 orphaned config file |

## The five heaviest, with detail

**T-1 — the whole `internal/` tree is invisible to any invokable lane.** This is the largest finding
by volume: 353 test files, roughly 1800 top-level `func Test*` (a floor — many use `t.Run` subtests,
so the real assertion count is higher), spanning all 21 legacy modules plus `kernel`, `platform`,
`composition`, `arch`, `adapters`, `testsupport`. `npm run harness:unit` (the only unit-test entry
point most operators would ever type) runs `go test ./tests/unit/...` — 14 files, 69 funcs — and the
FE suite; it never touches `internal/`. The only thing that runs `go test ./internal/...` is
`scripts/arch-gate.sh` step 4, and `grep -n "arch-gate" package.json` returns 0 hits — nothing in
`npm` calls it. This matches and re-confirms ledger D-51 exactly (measured again today: `bash
scripts/arch-gate.sh` → `EXIT=1`, 78s total, `go test ./internal/...` alone → 68s cached, `FAIL` on
`internal/composition` only, all 14 other packages `ok`). Practically: a change that breaks any test
in `internal/modules/pricing`, `internal/modules/orders`, `internal/kernel/fact`, or any of the other
~1800 test functions is invisible to `npm run harness:unit` and to anyone who has not separately
learned to run `bash scripts/arch-gate.sh` by hand. This is the single biggest "what does the suite
prove and WHEN" gap in the repo — most of the suite's bulk proves nothing about the next `git commit`
unless the operator remembers a script `npm` doesn't expose.

**T-2 — the ADR-023 module-boundary detector has no self-test.** `TestModuleBoundaryADR023`
(`apps/server_core/internal/composition/module_boundary_arch_test.go`) is 217 lines and exactly one
`func Test`. It walks the real `internal/modules` tree and reports every cross-module import that
does not go through `ports`. Read in full: there is no `testdata/` directory, no fixture, no assertion
that the detector returns zero on a clean synthetic tree, and no assertion that it can catch an
injected violation independent of the real repo's current mess (currently 234 real violations,
unchanged since the last audit — `go test ./internal/composition/...` reproduces `FAIL … 234
violation(s)` today). Compare this to its two siblings that share the same job description
("detector that must fire and must be provably silent"):
`apps/server_core/internal/arch/scan_test.go` maintains parallel `testdata/violations/` and
`testdata/clean/` trees and asserts both directions for every rule (e.g.
`TestCrossContextInternalFiresOnFixture` / `TestCrossContextInternalIsSilentOnCleanFixture`).
`apps/server_core/internal/platform/archguard/archguard_test.go` goes further: it has three
*must-fail* fixtures (`testdata/five_sites`, `testdata/aliased_site`, `testdata/three_sites`) that
inject a violation and assert the guard both fails AND names the offending symbol/file — explicitly
written to defeat a "generic count mismatch" failure message. `TestModuleBoundaryADR023` has none of
this. Combined with T-1, the detector is both unproven-by-fixture and unreachable by any invocable
lane — a regression in `moduleAndLayer`/`importedModuleAndLayer` that silently stopped flagging a
whole class of import (say, broke the `ports` carve-out check) would not be caught by any test, and
even if it were, nothing would run it.

**T-3 — 111 FE test cases across 9 files are unreachable from the repo root.** `packages/sdk-runtime`
(the hand-written SDK named in established fact 5) has its own `vitest.config.ts` and its own `test`
script (`vitest run --config vitest.config.ts`) — 7 files, 84 `it`/`test` calls declared, 86 actually
executed (some are parametrized). `packages/feature-classifications` and `packages/feature-inventory`
each have a bare `"test": "vitest run"` script and one test file apiece (15 and 10 cases). None of
these three packages' `test` scripts is ever invoked: the root `package.json` `test` and `test:run`
scripts (lines 13-14) are hardcoded to `npm run test --workspace @marketplace-central/web`, and
`scripts/harness.ps1`'s `Invoke-Unit` (line 63) does the same. Ran all three directly to confirm they
are real, not decorative: `sdk-runtime` → `7 passed (7)` / `86 passed (86)`, `feature-classifications`
→ `1 passed (1)` / `15 passed (15)`, `feature-inventory` → `1 passed (1)` / `10 passed (10)`. This
means the entire SDK — the surface every FE page and both `feature-classifications`/
`feature-inventory` depend on, and the exact seam established fact 6 says `GOV_API_SDK_SPLIT` never
checks for *agreement* — has 86 tests whose only value today is what a human remembers to run by hand.

**T-4 — two boundary invariants are proven only by same-change tests.** The brief's rule: an
agent-authored test may never be the SOLE evidence for a security, boundary, or migration invariant.
Two concrete instances found by walking recent commits that introduced such invariants: (a) Catalog
RLS tenant isolation — `migrations/0098_catalog_app_role.sql` (creates the least-privilege `mpc_app`
role) and `tests/integration/catalog_rls_role_test.go` (the only test that proves RLS policies are
enforced, not decorative) were both added in commit `7e3dcc47` ("papel de menor privilegio torna o
RLS do contexto provavel") — `git log --follow` on the test file shows exactly one commit, the same
one. The test itself is well built (positive control, negative control, tenant-B isolation — see "What
is actually fine" below) but it is honest in its own comment that it cannot exercise the actual
production connection path: "the application connects as the table owner… every policy evaluated
against that role is skipped" — this is established fact 12 (D-44) restated inside the test's own
doc comment, and the test's SET LOCAL ROLE approach is explicitly a workaround, not a test of the real
DSN. So the one test proving RLS works was written alongside the RLS mechanism, by the same change,
and even that test cannot reach the path the running application actually uses. (b) Catalogingest
tenant fail-closed guard — commit `47a76837` (`fix(catalogingest): tenant ausente falha fechado`,
today, `2026-08-07`) adds `requireTenantConfigured` AND `TestRequireTenantConfigured` in the same
commit (`git show --stat 47a76837`). This guard exists specifically to stop the catalog ingest from
silently writing ~10k rows under a fabricated `tenant_default` (D-39) — a data-integrity invariant if
there ever was one — and its only test was written by the same hand, in the same commit, that wrote
the guard. It is additionally unreachable by any lane (T-8): `cmd/catalogingest` is outside both
`tests/unit` and `internal/`, so even this self-authored test is never automatically run.

**T-5 — all Mercado Livre adapter tests are self-consistent, not contract-verified.** 15 files under
`internal/modules/connectors/adapters/mercado_livre` use `httptest.NewServer` to fake the ML API.
Reading one representative file in full
(`internal/modules/connectors/adapters/mercado_livre/prices_reader_test.go:24-30`) shows the fixture
is a raw JSON string literal written directly above the assertions that decode it — the same author,
in the same file, wrote both the shape the fake server returns and the code that parses that shape.
There is no golden fixture captured from a real ML response, and no schema check against ML's
published OpenAPI. This is not hypothetical drift risk: session memory (`ml-catalog-offers-pricing-api.md`,
`ml-api-live-verdicts-d120.md`) records that real ML quirks — `/products/{id}/items` requiring LEAF
children, `sale_price` requiring `?context=channel_marketplace`, `/items` multiget not existing for a
third-party seller — were all discovered by live drive against the real API, not by any of these 15
files, and none of the 15 files encode those quirks today. A future ML contract change (a renamed
field, a restructured nested object) would leave all 15 files green while the adapter silently
misparses production responses. AGENTS.md:14 ("provider payloads remain at adapters") isolates the
blast radius of a provider quirk to the adapter layer; nothing in the test suite verifies the adapter
still matches the provider it was isolated from.

## Test census

| tree | files | test funcs / cases | executed by which lane | runtime (measured) |
|---|---|---|---|---|
| `apps/server_core/tests/unit` | 14 | 69 `func Test` | `npm run harness:unit` (Go step) | 5.056s (package `ok`, cached) |
| `apps/server_core/tests/integration` | 21 | 47 `func Test` | `npm run harness:integration` | 259s total (ephemeral Postgres provision + migrate + run); no per-test count printed on success (T-6) |
| `apps/server_core/internal/modules` | 320 | 1646 `func Test` | **none** (npm-wired); `bash scripts/arch-gate.sh` step 4 only | included in 68s (`go test ./internal/...`, cached) |
| `apps/server_core/internal/contexts` | 5 | 23 | same as above | same |
| `apps/server_core/internal/kernel` | 8 | 44 | same as above | same |
| `apps/server_core/internal/platform` | 9 | 28 | same as above | same |
| `apps/server_core/internal/composition` | 7 | 37 | same as above — `FAIL` (234 ADR-023 violations + `TestNoVendorTokenInKernel`) | same |
| `apps/server_core/internal/arch` | 2 | 14 | same as above | same |
| `apps/server_core/internal/adapters` | 1 | 6 | same as above | same |
| `apps/server_core/internal/testsupport` | 1 | 2 | same as above | same |
| `apps/server_core/migrations` | 15 | 26 | **none** | 18.038s (`ok`, run directly) |
| `apps/server_core/cmd/catalogingest` | 1 | 3 | **none** | 2.561s (`ok`, run directly) |
| `docs/research/probes` | 1 | 1 | **none** — outside the `apps/server_core` Go module tree entirely (no covering `go.mod`), and gated by 2 required env vars even if invoked ad hoc | not run (structurally excluded) |
| `scripts/tests/*.ps1` | 13 | ~250+ `Assert-True`/`throw` checks (flat scripts, not discrete funcs) | **none** — manual `pwsh -File` only | not measured (out of scope to run 13 files individually against live infra) |
| `apps/web/src` | 57 | (part of 605 total) | `npm run harness:unit` (web step) / `npm test` | included below |
| `packages/feature-products` | 1 | (part of 605) | web step, via exact-filename include (T-10) | included below |
| `packages/web-query` | 3 | (part of 605) | web step, via glob | included below |
| `packages/ui` | 11 | (part of 605) | web step, via glob | included below |
| **web workspace total (reachable)** | **72** | **605** | `npm run harness:unit` / `npm test` | ~77s (82s combined Go+web minus 5s Go) |
| `packages/sdk-runtime` | 7 | 86 | **none** | 2.67s (`7 passed`, run directly) |
| `packages/feature-classifications` | 1 | 15 | **none** | 68.57s (`1 passed`, run directly — mostly jsdom environment setup) |
| `packages/feature-inventory` | 1 | 10 | **none** | 7.70s (`1 passed`, run directly) |

Go total: 405 test files (`git ls-files '*_test.go' | wc -l`), of which 353 live under `internal/`
(unreached by any lane) and 16 live under `migrations/`+`cmd/` (unreached by any lane) — **369 of 405
Go test files (91%) are unreachable from any npm/harness-invocable command**, reachable only via
`tests/unit` (14) and `tests/integration` (21) plus manual `bash scripts/arch-gate.sh`.

FE total: 81 test files repo-wide, 72 reachable (605 cases), 9 unreachable (111 cases) — **11% of FE
test files, 15% of FE test cases, are unreachable from the repo root.**

## Vacuous or unexecuted tests

| file:line | why it proves nothing (or nothing today) |
|---|---|
| `apps/server_core/internal/composition/module_boundary_arch_test.go:89-217` | Sole test for the ADR-023 detector; no fixture, no positive control, no must-fail proof — see T-2. Currently always red (234 violations) and, per T-1, never runs in any invokable lane, so its permanent-red state is invisible to begin with. |
| `apps/server_core/tests/unit/router_registration_test.go:172-179` | Asserts "not 404" for 10 routes — presence, not behavior. A route wired to the wrong HTTP method, or one that 500s on a malformed but non-empty response, still passes. |
| `apps/server_core/internal/modules/internal_read/adapters/oracle/reader_live_nocgo_test.go:10-15` | `TestOracleLiveBaseline` on a non-cgo host: first `t.Skip` never fires false (env var unset in every wired lane), and even if it did, an unconditional second `t.Skip` follows — the function can never execute a real assertion in this build configuration. |
| `apps/server_core/internal/modules/internal_read/adapters/oracle/reader_live_test.go:19-22` (build-tag `cgo`) | Same `MPC_ORACLE_LIVE_TEST` gate; no wired lane sets it, so this file's tests never run their body even when the cgo build tag is satisfied. |
| `apps/server_core/internal/modules/internal_read/adapters/oracle/sync_integration_test.go:27,36` | Same pattern, third file. |
| `packages/sdk-runtime/src/*.test.ts` (7 files, 86 cases) | Pass when run, prove nothing about the current HEAD to anyone using `npm test`/`npm run harness:unit` — see T-3. |
| `packages/feature-classifications/src/ClassificationsPage.test.tsx` (15 cases) | Same — unreachable, see T-3. |
| `packages/feature-inventory/src/StockSeguroPage.test.tsx` (10 cases) | Same — unreachable, see T-3. |
| `apps/server_core/migrations/*_test.go` (15 files, 26 funcs) | Pass when run directly; zero lanes reach `./migrations/...`. |
| `apps/server_core/cmd/catalogingest/main_test.go` (3 funcs, incl. `TestRequireTenantConfigured`) | Pass when run directly; zero lanes reach `./cmd/...`. This is the sole evidence for the T-4(b) fail-closed tenant guard, and it never runs automatically. |
| `docs/research/probes/ml_official_durable_price_probe_test.go` | Outside the `apps/server_core` module tree (no covering `go.mod`); requires two file-path env vars even for manual invocation. Structurally excluded from `go build ./...`/`go test ./...` run from the module root. |
| `scripts/tests/*.ps1` (13 files) | No automated invocation exists anywhere in the repo; `Policy.psm1:228` only excludes the directory from governance scanning. |
| `apps/web/vitest.config.ts:10-17` include list | `feature-products` entry is an exact literal filename, not a glob (T-10) — not vacuous today (0 dropped files) but structurally the same defect class as an untested branch: it will silently stop covering a second test file the day one is added. |

No self-comparison assertions (`x == x` style) or symmetric-fixture guarantees were found in the
sampled files. Given lane budget, roughly 20 of 405 Go test files and 8 of 81 FE test files were read
in full; the classes above (unreached trees, presence-only assertions, same-commit-only invariant
evidence, hand-authored provider fixtures) are the ones that surfaced. A wider read would very likely
surface more individual instances of the same classes, not new classes.

## What is actually fine

- **`apps/server_core/internal/arch/scan_test.go`** — exemplary. Parallel `testdata/violations/` and
  `testdata/clean/` fixture trees; every rule has both a "fires on violation" and a "silent on clean"
  test; `TestVendorTokensIgnoresAdaptersAndOwnList` and `TestCrossContextInternalSeesOutsideContexts`
  are explicit positive-control tests documenting a real historical blind spot (the detector used to
  skip files outside `contexts/` and missed a real violation) and proving it's fixed. This is the
  model the ADR-023 detector (T-2) should be held to, not a place needing correction itself.
- **`apps/server_core/internal/platform/archguard/archguard_test.go`** — exemplary. Three separate
  must-fail fixtures (`five_sites`, `aliased_site`, `three_sites`) each proving the guard both catches
  an injected violation AND names the offending symbol/file, explicitly written against "the brief's
  warning against vacuous guards" (its own comment, line ~310). `resolveCapabilityAliases` and its
  fixture close a real adversarial-review bypass (variable-extraction alias). This is real security
  engineering discipline in test form — undermined only by T-1 (nobody's lane ever runs it).
- **`apps/server_core/tests/unit/cache_composed_test.go`** — a genuine composed-chain test (source →
  cache → timing → routing → application → HTTP), with call-count assertions
  (`oracle.listCalls.Load()`) proving cache hits/misses/bypasses actually happen, not just that a
  response arrives. Comments explicitly explain why the composition order matters and what a shortcut
  version used to hide.
- **`apps/server_core/tests/integration/catalog_rls_role_test.go`** — textbook control structure:
  positive control (uncontrolled read must see both tenants' rows, or the rest of the test proves
  nothing), tenant-A/tenant-B isolation in separate transactions, and an explicit negative/fail-closed
  control (no tenant scope → 0 rows, not all rows). The T-4(a) caveat (same-commit authorship, cannot
  reach the real production connection role) is about independent verification, not about this test's
  internal quality, which is high.
- **`apps/web/src/pages/integracoes/SyncHealthCard.test.tsx`** — mocks at the client-hook boundary
  (`useClient`), not by injecting rendered state directly (the D-11 pattern from session memory that
  used to make this exact area vacuous is not present here); the isolation test has an explicit
  anti-vacuity comment: "a partial mock would make every other card fail too… and the isolation
  assertion would be vacuous." This file appears to have been rebuilt since the memory note was
  written and is currently sound.
- **The ephemeral-Postgres integration lane** (established fact 4, re-confirmed: `container=ephemeral`
  or `session-reuse`, fresh `mpc_test_<hex>` database) genuinely prevents state accumulation between
  runs — corroborates the ledger's account of D-54 (the `count(*) = 3` absolute-assertion test is not
  a false-green *today* specifically because of this).
- `migrations/*_test.go` and `cmd/catalogingest/main_test.go`, though unreached by any lane (T-8), are
  real and passing when invoked — not decorative placeholders. The problem is purely "when," not
  "what."

## Unverified / needs judgment

- `scripts/tests/*.ps1` — did not execute all 13 files against live infrastructure (would require
  spinning up Postgres/Docker sessions repeatedly outside the harness's own orchestration, and several
  are named `*.integration.tests.ps1`, implying they expect state this lane should not provision ad
  hoc). One (`governance-contracts.tests.ps1`) was invoked directly and failed on a missing `rg` binary
  on `PATH` in this shell — `unverified` whether that is a real portability gap in the script or an
  artifact of this specific invocation environment (the harness's own child-process environment
  construction, `New-HarnessChildEnvironment`, was not used for this ad hoc run).
- Exact test-function counts for `internal/**` are a floor: `grep -c '^func Test'` counts top-level
  functions only. Table-driven tests using `t.Run` subtests execute more assertions than the function
  count shows; the true "how many independent checks" number is materially higher than 1800.
- `contracts/governance/modules.json`, `contracts/governance/schemas/modules.schema.json`,
  `scripts/harness/Policy.psm1`, `scripts/tests/governance-drift.tests.ps1` showed as modified in
  `git status` at measurement time, not by this lane (this lane made no edits outside this report).
  Consistent with other lanes running concurrently against the same working tree per the ten-lane
  roster; not investigated further as it is out of this lane's dimension.
- Did not exhaustively read all 405 Go / 81 FE test files; vacuity classes reported are what a
  budget-bounded sample surfaced, not a proof that no other instances exist.
- Whether the 234 ADR-023 violations reported by `TestModuleBoundaryADR023` today match the count in
  `docs/architecture/decisions/023-module-protocol.md`'s corrected header (234, per commit `82bd18ef`)
  was re-measured and does match — cited as corroboration, not as new work, per PHASE-0's instruction
  not to re-derive established facts.

## Commands run

```
$ cat .mnfs/HARNESS-DEBTS.md | grep -iE "test|mock|vacuous|gate|guard" -n | head -100
$ cat package.json | grep -A 30 '"scripts"'
$ cat scripts/arch-gate.sh
$ wc -l .mnfs/HARNESS-DEBTS.md && grep -n "^## D-\|^\*\*D-" .mnfs/HARNESS-DEBTS.md | tail -60
$ git ls-files 'apps/server_core/tests/unit/**_test.go' | wc -l            # 14
$ git ls-files 'apps/server_core/tests/integration/**_test.go' | wc -l     # 21
$ git ls-files 'apps/server_core/internal/**_test.go' | wc -l              # 353
$ git ls-files '*_test.go' | sed -E 's#...#\1#' | sort | uniq -c | sort -rn
$ for d in tests/unit tests/integration internal/modules internal/contexts internal/kernel \
    internal/platform internal/composition internal/arch internal/adapters internal/testsupport \
    migrations cmd; do
    git ls-files "$d/**_test.go" | xargs grep -hoE '^func (Test|Fuzz|Benchmark)[A-Za-z0-9_]*' | wc -l
  done
$ git ls-files 'apps/web/**' | grep -iE '\.test\.(ts|tsx)$' | wc -l         # 57
$ cat apps/web/vitest.config.ts
$ git ls-files 'packages/**' | grep -iE '\.test\.(ts|tsx)$' | sed -E 's#(packages/[^/]+)/.*#\1#' | sort | uniq -c
$ cat vitest.config.ts   (root)
$ cat packages/sdk-runtime/package.json
$ cat apps/web/package.json | grep -A5 '"scripts"'
$ grep -n '"test' package.json
$ npm run harness:unit                        # EXIT=0, DURATION=82s; Go ok 5.056s; Test Files 72 passed (72); Tests 605 passed (605)
$ npm run harness:integration                  # EXIT=0, DURATION=259s; container=session-reuse; status=passed; no test counts printed
$ cd apps/server_core && GOCACHE=$PWD/.gocache go test ./internal/...   # EXIT=1, DURATION=68s; FAIL internal/composition (234 violations); all else ok
$ bash scripts/arch-gate.sh                    # EXIT=1, DURATION=78s
$ grep -rn "t\.Skip\b" (git ls-files '*_test.go')   # 6 sites, 3 files, all Oracle live-gate
$ grep -rln "^//go:build integration" apps/server_core --include="*_test.go" | wc -l   # 74
$ grep -rl "toBeInTheDocument" apps/web/src --include="*.test.tsx" | wc -l  # 46
$ find packages -iname "vitest.config*"        # only packages/sdk-runtime/vitest.config.ts
$ cd packages/sdk-runtime && npx vitest run --config vitest.config.ts     # 7 passed (7) / 86 passed (86), 2.67s
$ cd packages/feature-classifications && npx vitest run                  # 1 passed (1) / 15 passed (15), 68.57s
$ cd packages/feature-inventory && npx vitest run                        # 1 passed (1) / 10 passed (10), 7.70s
$ cd apps/server_core && GOCACHE=$PWD/.gocache go test ./migrations/...  # ok, 18.038s
$ cd apps/server_core && GOCACHE=$PWD/.gocache go test ./cmd/...         # ok apps/server_core/cmd/catalogingest 2.561s; no test files elsewhere
$ grep -rln "httptest.NewServer" apps/server_core/internal/modules --include="*_test.go" | wc -l   # 15
$ git log --oneline -1 -- apps/server_core/migrations/0098_catalog_app_role.sql   # 7e3dcc47
$ git log --oneline --follow -- apps/server_core/tests/integration/catalog_rls_role_test.go  # only 7e3dcc47
$ git show --stat 47a76837
$ grep -n "GOV_API_SDK_SPLIT" -r scripts/harness/Policy.psm1 contracts/governance/*.json
$ pwsh -NoProfile -File scripts/tests/governance-contracts.tests.ps1      # real EXIT=1 (rg not on PATH in this shell)
$ ls scripts/tests/*.ps1 | wc -l                                          # 13
$ for f in scripts/tests/*.ps1; do grep -c "Assert-True\|Assert-False\|throw " "$f"; done
$ git status --porcelain --untracked-files=all
```
