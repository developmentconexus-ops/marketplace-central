# Phase 2 — Synthesis, Arm B

Written against `SYNTHESIS-BRIEF.md`, `PHASE-0.md` (budget section as amended 2026-08-07), and the ten
lane reports. Read-only pass; no gate, suite, or harness command was executed. Claims not made by any
lane are marked **[Arm B]** with the check that would confirm them.

Two numbers govern everything below and both come from PHASE-0's established facts, not from me:
**nothing blocks a merge to `main` today and nothing ever has**, and **every instrument that exists is
red at HEAD**. The second is why the first was never fixed. That pairing is the spine of this
synthesis.

---

## 0. What the two constraint amendments changed

GitHub Actions is available and billing is resolved, so level 3 is reachable for any rule expressible
as a command that exits non-zero. The remote is real: `origin/main` is `7df7d011` (2026-08-04), local
`main` is **98 ahead, 0 behind, fast-forward possible, nothing diverged**. And the repository will be
made **public on a free personal account**, which makes rulesets and enforced required status checks
available at zero cost.

**A correction to a lane finding, stated loudly.** `lanes/cicd.md` F2 records `origin` as
unresolvable. That was an authentication artifact, not a fact about the repo: two GitHub identities
sit in the keyring, `gh` asks as `leandrotcawork` while git authenticates as
`developmentconexus-ops`. The `legacy` remote (`leandrotcawork/marketplace-central`, default branch
`master`, last push 2026-07-20) is superseded and unreachable under current git credentials — that
part of F2 stands. **Everything in this document is written against a reachable `origin`.** The
cicd lane's conclusion that nothing blocks a merge is unaffected; only its explanation was wrong.

What remains true, and is still the first-order problem:

- Nothing blocks a merge to `main` today and nothing ever has (PHASE-0).
- Every instrument that exists is red at HEAD. That is why the first fact was never fixed.
- **No change has ever travelled through the remote.** 98 commits sit local, there are no pull
  requests, and `docs/HARNESS-PROFILE.md:279` forbids the act that would create one — *"Push: NEVER
  without explicit operator permission"* — with the landing path defined as a local merge into an
  integrated `main` followed by a post-merge ladder (`:77`). **[Arm B]** No lane examined the
  harness's landing path, because no lane was pointed at it; confirm by reading
  `HARNESS-PROFILE.md` §merge and by `git log --format=%D` showing chip tips merged into local `main`
  with no intervening remote ref.

So the remaining defect is not reachability. It is that **a required check would sit beside a road no
change travels.** `release-images.yml` is the proof already in the repo: it exists, it is well-formed,
and it has one recorded run, which failed (`lanes/cicd.md` F1, F3). Adding a workflow without changing
the landing path inherits exactly that history.

Two further consequences shape §4 throughout rather than as footnotes. **Public repositories change
what a workflow may do** — fork pull requests execute, logs are world-readable, and
`pull_request_target` with a writable token is a known privilege-escalation path — so every workflow
below states its trigger, its `permissions:` block, and whether an untrusted fork can reach it. And
**the 98 commits stay local until the gates exist**, so the sequence in §5 lands the gate topology
before the first push, not after.

---

## 1. Axes

Eight axes. Each finding appears once. Where a lane filed a finding under a symptom that belongs to a
different cause, the reattribution is stated.

### A — The verifier is not on the path a change takes

**Cause.** Two facts compose into one defect: no instrument is invoked by an event, and no change
travels through a place where an event could fire. Fixing either alone is strictly more expensive.
Wire workflows first and you get a gate with zero executions. Push 96 unvalidated commits first and
you get a permanently red default branch that every later required check inherits, which is the one
outcome that guarantees the checks get turned off.

**Findings.** `lanes/cicd.md` F1 (no `pull_request` workflow), F2 (no branch protection — its
"unreachable origin" explanation is corrected in §0; the absence of protection is real), F3 (the one
remote run — history, not constraint), F4 (no git hooks, no husky, no lint-staged), F5
(no `.claude/settings.json` hooks), F10 (all 11 `scripts/tests/*.tests.ps1` invoked by nothing;
`Invoke-Pester` → 0 hits), F11 (6 of 11 non-green cold), F12 (`knowledge-routes.json:12` cites a
test file that does not exist), F13 (`Policy.psm1:447` — `-BaseSha` is opt-in with no default, so a
governance run without it silently skips every diff-scoped rule), F14 (`scripts/release.sh` bypasses
CI by construction), F15 (dead warn-only `wiki-lint` on a `legacy` branch, 3 runs 3 failures).
`lanes/testing.md` T-1 (`internal/**` — 353 test files, ~1800 `func Test*` — reached by no npm or
harness command), T-3 (9 FE test files / 111 cases structurally unreachable from root), T-6 (the
integration lane prints no RUN/PASS/SKIP counts, so an all-skipped run is byte-identical to a green
one), T-7, T-8, T-10 (`apps/web/vitest.config.ts:14` pins a package by exact literal filename), T-12
(orphaned root `vitest.config.ts`). `lanes/frontend.md` F2 (`tsc --noEmit` wired to nothing, fails
with 12 real errors, 3 in production code, under `tsconfig.base.json:6` `"strict": true`).
PHASE-0 facts 1–3, and the D-51/D-52 open decisions.

The two aggregate red numbers — `lanes/cicd.md` F6 (58 governance violations at tip with zero diff)
and F7 (44 archscan findings: kernel 0, contexts 1, adapters 1, composition 42) — belong here **as
ratchet inputs only**. The violations they count are attributed to the axes that own their causes
(37 `RCFG_*` and 20 `GOV_MODULE_*` → axis D). A number is not a finding; it is a starting line.

**Why it is one axis.** Every item is the same shape: a correct instrument disconnected from
execution. Five lanes found it independently — cicd in the pipeline, testing in the test census,
frontend in `tsc`, layering in the scanners, persistence in the tenant-predicate check that does not
exist. When five lanes converge on one shape from five directions, that shape is a cause.

### B — The published contract is four hand-written copies with no arbiter

**Cause.** The same shape is authored four times (Go DTO → OpenAPI → SDK → FE consumer) and the only
rule connecting them is `Policy.psm1:460`, a literal XOR that requires the OpenAPI file and the SDK
file to change in the *same commit* and never checks that they *agree*.

**Findings.** `lanes/duplication.md` DUP-1 (239 OpenAPI schemas vs 172 SDK interfaces / 2595 lines,
with the traced `ListingReadModel` divergences: `listing_type` and `price` are Go pointers that can
emit `null` while OpenAPI marks them required and the SDK types them non-null; `market_signal` and
`signal_status` are always serialized by Go but optional in both). `lanes/delivery.md` F3 (12
independent `map*Error` functions string-prefix-matching `err.Error()`; zero `errors.Is`/`errors.As`
in transport), F6 (7 hand-copied page-envelope structs), F7 (3 divergent envelope shapes including
`market/collection_handler.go:64` emitting the Portuguese JSON key `"decisões"` — the only
non-English wire field in the product), F12 (12 distinct error codes for "malformed JSON body"), F13
(2 hand-maintained string-literal copies of the error envelope, forced by the `apierror→httpx` import
cycle), F14 (10 OpenAPI+router operations with zero SDK method, which `GOV_FRONTEND_FETCH` renders
unreachable from the web app under the repo's own rule), F15 (5 SDK query-builders plus ad hoc
concatenation). `lanes/frontend.md` F1 (nothing links OpenAPI to the SDK but a code comment at
`index.ts:1747`; no zod/io-ts/ajv/joi/yup anywhere), F7 (33 of 41 `<ErrorState>` sites discard the
backend's typed error code), F13 (the SDK is 354 lines wider than fact 5's headline).

`lanes/delivery.md` F4 also lands here: Go 1.22 method-prefixed routes (51 of 98 registrations) get
stdlib's plain-text 405 from `net/http/server.go:2707`, bypassing `apierror.Write`. That is a
*fourth, accidental* producer of the envelope — a contract violation the contract never named.

**Why it is one axis.** One mechanism closes all of it: generate from one authored artifact, and make
disagreement a compile or test failure rather than a review question. The approved-but-uninstalled
`oapi-codegen` + `openapi-typescript` (fact 7) is exactly this mechanism, already decided and never
installed.

### C — Identity is a boot-time constant, so every downstream enforcement enforces a constant

**Cause.** There is no per-request principal and no per-request tenant. `pgdb/tenant.go:5-9` returns
a config fallback or the literal `"tenant_default"`. Everything built on top of identity — RLS, the
tenant predicate in 246 query sites, CORS, rate limiting, audit — is therefore enforcing the same
value for every caller, which is the same as enforcing nothing.

**Findings.** `lanes/security.md` SEC-1 (zero HTTP authentication on the entire API; no login page;
no `Authorization` sender in the FE), SEC-2 (buyer full name, CPF/CNPJ and full address served by
that unauthenticated surface at `orders/transport/http_handler.go:608-618`), SEC-3
(`pgdb.LoadConfig()` fails **open** on `MC_DEFAULT_TENANT_ID` while `MC_DATABASE_URL` and
`MPC_ENCRYPTION_KEY` fail closed in the same function), SEC-4 (RLS's two preconditions both fail: a
superuser+`bypassrls` DSN, and the GUC set only on the new-context write path, not the legacy read
path that actually serves `/catalog/*`), SEC-5 (1 of 85 migrations enables RLS; 51 of 52
tenant-bearing migrations have no DB backstop), SEC-6 (`Access-Control-Allow-Origin: *` with
`Authorization` in the allowed headers — inert today, a live hole the hour auth ships), SEC-8 (a
second, dead, fail-open `DefaultTenantID` helper). `lanes/delivery.md` F1, F2, F10 (`MaxBytesReader`
in 2 of 20 transport packages), F11 (CORS `*` on all ~104 routes). `lanes/persistence.md` P-1 (RLS on
4 of 64 live tables), P-2 (tenant scoping is 100% hand-convention across 246 raw query sites; 285/285
sampled statements correct; **zero** automated checks), P-3 (`set_config('app.tenant_id', …)` at
exactly one site: `internal/contexts/catalog/internal/postgres/repository.go:49`), P-4
(`pgdb/config.go:23-24` still defaults `tenant_default` for `cmd/server`; D-39's guard covered only
`cmd/catalogingest`).

SEC-2 is the axis's severity anchor and it is not a code-quality finding. The unmasked fiscal block is
*correct* — a Brazilian nota fiscal legally requires it, and the DTO comment cites C06/LGPD. The
defect is entirely in who can read it, which is: anyone who can route a packet to the process.

**Sequencing fork, stated explicitly.** PHASE-0 marks network reachability `unverified`. If the
service is reachable from anything but localhost, C is not axis four — it is a live incident and it
preempts the entire program. Confirm before scheduling.

### D — Rules and their exception ledgers are hand-maintained against the tree being abandoned

**Cause.** Every governance rule in this repo is a hand-written walk over a hardcoded path, paired
with a hand-written exception list that nothing reconciles against the code. Both halves decay, and
the decay is invisible because the ledger is the only thing anyone reads.

**Findings.** `lanes/layering.md` L-01 (`modules.json` `temporary_exceptions` is a second,
disconnected ledger; `Policy.psm1:98-102` never opens a `.go` file; 3 of 5 declared exception paths
do not exist on disk, the other 2 sit inside the 234 uncounted violations, and all 5 carry
`removal_owner: M-10`, a milestone deferred since 2026-07-14 under a superseded mission), L-02 (the
same ~50-line boundary scanner hand-written twice in `channelfees` and `divergences`, zero reuse of
the superset `arch.ScanVendorTokens`, 19 of 21 modules uncovered), L-03 (disputed — see §7), L-04
(`TestModuleBoundaryADR023` matches only the literal `.../internal/modules/` prefix, so
modules↔contexts coupling is structurally invisible), L-05 (of the 234 violations, 146 originate in
`adapters/`, 42 in `application/`, 9 in `ports/` — ports are typed by the target module's `domain`,
so the boundary is itself the leak), L-06 (49 directed edges across 17 of 21 modules; 3 bidirectional
pairs; `connectors` + `internal_read` + `erp_import` absorb 132 of 234), L-08 (no detector reaches
`apps/web/src` or `packages/sdk-runtime`; 77 FE vendor-token occurrences are data, not a rule
violation). `lanes/cicd.md` F8 (`Policy.psm1:302` hardcodes the walk root to `internal/modules`, so
`internal/contexts` — the tree ADR-023 says is the future — has **zero** governance coverage: 10
non-test files ungoverned against 462 governed), F9 (a negative fixture already asserts a
`GOV_CONTEXT_UNREGISTERED` code that does not exist in `Policy.psm1`).

The registry half of `lanes/cicd.md` F6 is reattributed here: 37 of the 58 violations are `RCFG_*`
(20 `RCFG_UNAPPROVED_READER`, 10 `RCFG_READER_MISSING`, 7 `RCFG_UNDECLARED_READ`) — a reader registry
that has drifted from the code by exactly L-01's mechanism, in a second file. The 20
`GOV_MODULE_DEPENDENCY` + `GOV_MODULE_LAYER` violations are the same story in a third.
`lanes/persistence.md` P-8 (2 migration-prefix collisions) is the fourth instance; see §7 for why its
single exception entry is the clearest proof the mechanism has no expiry.

**Why it is one axis.** Four registries, one failure mode. A per-registry fix is four fixes plus a
fifth registry next quarter. The mechanism fix is one rule: *an exception entry that does not
correspond to a live violation is itself a violation.*

### E — The idiom for an unhandled error is `_ =`, and nothing makes a discarded failure findable

**Cause.** Failure has no durable channel. `slog` is used in 36 files but **no handler is ever
configured** (`lanes/observability.md` OBS-9), so output is Go's default text on stderr, from a
process whose stderr is not collected. In that world, discarding an error and logging one are
equivalent, and the code has consistently chosen the cheaper.

**Findings.** `lanes/observability.md` OBS-1 (`/healthz` returns 200 unconditionally with no
dependency check, and is the Docker `HEALTHCHECK` at `docker-compose.yml:46`;
`deploy/docker-compose.prod.yml` has no backend healthcheck at all), OBS-2 (5 of 6 ticker loops have
no `recover()`, in a single process that also serves HTTP — one panic takes the API down; the fix is
already written in-repo at `sync/scheduler.go:201-208`), OBS-3 (`FeeSyncScheduler` imports no logger
at all), OBS-4 (`runJob` silently returns on cursor-read error and discards both `RecordFailure` and
`RecordSuccess` write errors), OBS-5 (13 of 15 ML adapter files never call `slog`), OBS-7 (zero
metrics, zero tracing, no `/metrics`), OBS-8 (no correlation ID), OBS-11
(`cmd/catalogingest` failure surfaces only as stderr plus exit code), plus the 10-site
silent-failure census. `lanes/security.md` SEC-7 (raw `recover()` value logged verbatim).
`lanes/persistence.md` P-7 (`stale_since` written twice, read zero times).
`lanes/goidiom.md`'s 9 discarded `json.Unmarshal` returns.

One census detail decides the fix: **zero `slog.Debug` calls exist in the repo.** The failure mode is
therefore never "logged at a level nobody reads." It is "not logged at all." That means a log-level
policy is worthless here and the actual fix is small: configure one JSON handler, wrap the six loops,
and make the discard sites emit.

OBS-6 (the `/sync/health` row collapse on `installation_id`, invisible until a second installation
exists) also sits here: it is a failure that cannot be *seen*, which is the same cause one layer up.

**Explicitly excluded from this axis:** the ML token-refresh failure path. The observability lane
re-verified it as durable and tested (`recordRefreshFailure`, `degradeAfterRefreshFailure`,
`integrations_refresh_failure_test.go`). It is on the do-not-touch list. Do not re-open it.

### F — Money has three representations and none of them is compulsory

**Cause.** An exact-decimal type exists, is correct, and is optional. `lanes/duplication.md` DUP-3 is
the whole axis in one line: **`internal/kernel/exact` has zero production callers** — all four
importers are `_test.go` — under a doc comment reading "There is no constructor from float64 anywhere
in this package, and that is the point." The point was made and nothing enforces it.

**Findings.** `lanes/goidiom.md` G-1 (373 float64/float32 occurrences across 78 files, concentrated
in money fields in `orders`, `profitability`, `marketplaces`, while two exact engines exist).
`lanes/duplication.md` DUP-2 (`Money` defined 5×, 4 byte-identical `{Amount string; Currency string}`;
canonical at `kernel/exact/money.go:45`), DUP-3, DUP-4 (12 lossy `strconv.ParseFloat` sites,
including `orders/adapters/pricingtax/reader.go:309` round-tripping an already-exact
`FormatRatHalfUp(r,2)` string back through float64), DUP-5 (the margin formula implemented twice with
a **unit mismatch** — `orders/domain/order_decomposition.go:118-163` uses `*float64` and a raw
fraction; `pricing/domain/decompose.go:141-286` uses `big.Rat` with ×100 baked in; zero shared code),
DUP-9 (15 files hand-construct `big.Rat`). `lanes/persistence.md` P-5 (money scanned as
`pgtype.Float8` at 25 sites in `orders` and `profitability`, against the documented `::text`
discipline used only by `pricing/matrix_reader.go`). `lanes/delivery.md` F8 (two money wire dialects:
exact string and bare float64). `lanes/frontend.md` F14 (client-side DIFAL aggregation in JS
doubles).

DUP-5 is the finding that makes this an axis rather than a cleanup. Two margin formulas with
different units, no shared code, both live: the failure is not precision loss, it is **a number that
is wrong by 100×** and looks plausible in both places. That is the machine-author failure mode
exactly — the wrong noun used correctly everywhere.

The one guarded crossing, `pricingtax/reader.go:266` `lineTotal`, is the model to generalize, not to
remove.

### G — The guard and its proof are written by the same agent in the same change

**Cause.** Nothing in this repo distinguishes a guard that has been proven to fire from one that has
only been observed to pass. Green is the default state of an assertion that cannot fail.

**Findings.** `lanes/testing.md` T-2 (`TestModuleBoundaryADR023`: one test function, **no `testdata/`
directory at all**, no positive control, no must-fail proof — while both its siblings have one;
`internal/arch/scan_test.go` runs parallel `testdata/{violations,clean}` and `archguard_test.go`
carries 3 must-fail fixtures naming the offending symbol), T-4 (two boundary invariants have
same-commit-only evidence: catalog RLS, where `7e3dcc47` adds both `0098_catalog_app_role.sql` and its
only test; and the catalogingest tenant guard, where `47a76837` adds `requireTenantConfigured` and
`TestRequireTenantConfigured` together), T-5 (all 15 ML adapter test files use inline hand-authored
`httptest` JSON, zero golden or live-captured fixtures, and the known real ML quirks are encoded in
none of them), T-9 (`router_registration_test.go:172-179` asserts "not 404" only), T-11 (6 `t.Skip`
Oracle-live sites plus a `//go:build !cgo` double-skip). `lanes/cicd.md` F9 belongs to D as a symptom
of registry drift but is *evidence* for G: a fixture asserting a reason code that does not exist has
never been run against a failing input, because it would have failed.

**[Arm B]** I verified T-2's central claim directly: `find apps/server_core/internal/composition -type
d -name testdata` returns empty while `internal/arch/testdata/{clean,violations,adapters}` exists.
This is why I side with testing over layering in §7.

**Why it is its own axis and not part of A.** A is about instruments not running. G is about
instruments that run and cannot fail. They need opposite fixes — A needs wiring, G needs a negative
fixture requirement — and G must land *before* the bulk of A's ratchets, or the ratchets certify
themselves.

### H — The browser has no failure containment, so any render error is total

**Cause.** There is no boundary between a component's failure and the application's. Every other item
here is an amplifier of that one fact.

**Findings.** `lanes/frontend.md` F4 (zero React error boundaries anywhere in the repo), F3 (two live
routed components render `<ErrorState detail="…" />` without the required `onRetry` —
`MutationPreviewModal.tsx:210` and `MutationResultSummary.tsx:22` — a `TS2741` that becomes a runtime
`TypeError` on click), F8 (`ClassificationsPage.tsx` combines every risk: 7× `catch (err: any)`, the
only non-test `any` in the repo; no react-query at all; and a test file no wired command has ever
run), F9 (three incompatible data-fetching patterns), F10 (the D-11 `query.status`-switch fix applied
to exactly 1 file, 9 still on the vulnerable ternary chain).

F3 + F4 together are the entire axis in one sentence: a type error that `tsc` already catches, in code
`tsc` is not run against, in an app where the consequence of that error is a white screen instead of a
broken panel. Small, cheap, and it changes the blast radius of every future frontend defect.

---

## 2. Kill the noise

Findings that are real, measured, and do not earn a slot. For each: cost of fixing versus cost of
leaving it.

**Leave forever — the current state is correct or the fix is a net negative.**

- `lanes/duplication.md` DUP-10 — 83 `strings.TrimSpace(x) == ""` checks. This is idiomatic Go. A
  helper saves nothing and adds an indirection. Leave forever.
- `lanes/duplication.md` DUP-11 — `Assert-True` defined 6× in PowerShell test files. Test-local
  helpers duplicated across independent scripts is the correct trade; sharing them couples the
  scripts. Leave forever.
- `lanes/goidiom.md` G-7 — 27 of 35 adapter interfaces correctly unexported and local. This is a
  *compliance* measurement reported as a finding. Nothing to do.
- `lanes/goidiom.md` G-8 — 200 `map[string]any` occurrences, 4 in `domain/` as named passthrough
  fields, **0 domain-layer key access**. The measured shape is already the safe one. Leave.
- `lanes/goidiom.md` G-5 — PHASE-0 fact 11 already declares this deliberate foundation. Not a finding.
- `lanes/frontend.md` F15 — a11y baseline (136 `aria-*`, 59 `role=`, 0 bare `<img>`, 0 div/span
  `onClick`). Reported as fine; keep it that way, do nothing.
- `lanes/persistence.md` P-9 — 14 `fmt.Sprintf` WHERE-builders, **all sampled safe**. A guard here
  would be written against zero measured instances of the defect. Building a detector for a defect
  you cannot demonstrate is how you get a detector nobody can prove fires (axis G). Leave; revisit if
  one unsafe instance ever appears.

**Fix the instance, never build the mechanism.**

- `lanes/delivery.md` F5 — `parsePositiveInt` copied 5×, one behaviourally diverged
  (`product_links/.../http_handler.go:514-524` uses `fmt.Sscanf`, which accepts `"20xyz"` → 20;
  `orders/.../http_handler.go:751-757` uses `strconv.Atoi`, which rejects it). Fix the one file. Cost:
  minutes. A shared helper for a 5-line parser is not worth a new import edge across 5 modules.
- `lanes/goidiom.md` G-2 — clock injection inconsistency. Only one real bug in it:
  `classifications/application/service.go:46,74` drops `.UTC()`. One line. The other 6 inline
  `time.Now()` files are not worth an injection campaign. `auth_flow_service.go:812` bypassing its own
  injected clock is a second one-liner; take it.
- `lanes/goidiom.md` G-3 (2 `fmt.Errorf` with no verbs) and G-6 (1 `%v` that should be `%w`, out of
  12). Three-line total. `go vet` already passes, so these are style, not behaviour. Take them
  opportunistically; do not schedule them.
- `lanes/frontend.md` F11 — dead `pages/DashboardPage.tsx`, 160 lines, zero importers. Delete it. Cost:
  one commit. Cost of leaving: it holds 2 live SDK calls that will show up in every future contract
  reconciliation as false positives.
- `lanes/frontend.md` F12 — 3 empty workspace shells plus a non-functional `apps/landing`. Delete or
  leave; either is fine. They cost nothing but they widen every workspace-glob decision.

**Real, but the cost of leaving is genuinely low.**

- `lanes/duplication.md` DUP-7 — 19 module-local `write*Error` wrappers, 3 byte-identical, **all
  bottoming out in `apierror.Write`**. The sink is already unified (CHIP-ERROR-UNIFY landed). Nineteen
  thin wrappers over one correct sink is duplication with no divergence surface. Leave.
- `lanes/duplication.md` DUP-8 — 2 near-duplicate `policy_reader.go` adapters. Two is not a pattern.
  Leave until three.
- `lanes/duplication.md` DUP-6 — 82 hand-written `for rows.Next()` loops, 0 uses of pgx v5's
  `CollectRows`/`RowToStructByName`. **Mis-sized upward; see §7.** Leave forever, deliberately.
- `lanes/observability.md` OBS-10 — two logging mechanisms. Collapses to nothing once E configures one
  handler; not separately schedulable.
- `lanes/observability.md` OBS-12 — webhook stats permanently canonical-zero because no receiver
  exists. This is ADR-017/034 working correctly: an unknown fact rendered as unknown. Not a defect.
- `lanes/cicd.md` F12 — a knowledge-route citing a nonexistent test file. One-line correction, or
  delete the route. Not a program item.
- `lanes/cicd.md` F15 — the dead `wiki-lint` workflow on a `legacy` branch. Delete the branch or
  ignore it. It cannot run.
- `lanes/persistence.md` P-8 — 2 migration-prefix collisions. The runner sorts by **full filename**,
  so both collisions are deterministic and harmless today. The *exception entry* is the finding (§7,
  axis D); the collisions themselves are noise.

**One I am declining to kill, against the temptation.** `lanes/persistence.md` P-6 — 65 hand-typed
`CHECK (… IN (…))` constraints, **0 `CREATE TYPE`**, at least 2 vocabularies duplicated verbatim
across files. This looks like style and is not: it is the same shape as DUP-5 and B's four-copy
contract, one layer down, and a vocabulary that drifts between two `CHECK` lists produces rows that
are valid in one query and invalid in another. It rides along with axis D's registry work at near-zero
marginal cost. Not scheduled separately; not killed.

---

## 3. Consolidated do-not-touch list

Merged from every lane's "what is actually fine" section. These are load-bearing and correct. A change
program that damages one of these has gone backwards.

**Verification and guards that actually work**
- `internal/arch/scan_test.go` and its `testdata/{clean,violations,adapters}` fixture pattern — the
  reference shape for every guard in §4.
- `internal/platform/archguard` and its 5 must-fail fixtures naming the offending symbol.
- The exact 8-panic / 8-declaration match between the code and `invariants.json`; 0 undeclared, 0 in
  `cmd/`.
- `go vet` and `go build` clean at HEAD.
- The ephemeral `mpc_test_<hex>` database pattern in the hermetic lane.
- `Policy.psm1`'s rule logic where it has coverage, including its cross-branch-contamination exclusion
  list.
- `GOV_FRONTEND_FETCH` — a working level-3-shaped rule with a real regex and a real scope.

**Contract and error handling**
- `apierror.Write` as the single envelope sink (192 sites). CHIP-ERROR-UNIFY landed; do not
  re-litigate the sink. B's work is upstream of it, not on it.
- `errors.Is` / `errors.As` discipline in the non-transport layers: 120 + 37 against 42 sentinels.
- 355 `%w` against 12 `%v`.
- 141 of 280 interfaces in `ports/`, **0 in `domain/`**.
- The `route_deadline` middleware — no bypass path.
- ADR-017/034 honest-unknown formatters.

**Money and persistence**
- `pricing/domain.Money` and `kernel/exact` as types. F makes them compulsory; it does not redesign
  them.
- `lineTotal`'s guarded float64→decimal crossing at `pricingtax/reader.go:266` — the model to
  generalize.
- The migration runner's per-migration transaction and its idempotency.
- All 116 `timestamptz` columns.
- The ERP→Postgres `mergeSnapshotTx` write path.
- The two correctly non-tenant global tables.
- The now-wired `icms_matrix_mirror` scheduler.
- The DIFAL override mechanism — real, wired, audited. The duplication lane explicitly retired the
  stale "3 copies / 4 unused overrides" number; the three ICMS/tariff tables are three different
  facts.

**Security and integration**
- The ML token-refresh failure path: `recordRefreshFailure`, `degradeAfterRefreshFailure`,
  `integrations_refresh_failure_test.go`. **Fixed. Do not re-open.**
- AES-256-GCM credential encryption, fail-closed on `MPC_ENCRYPTION_KEY`.
- The Oracle config loader — the reference shape for fixing SEC-3.
- PII redaction on the raw-capture path and its tests; ML secret redaction tested in both directions.
- `MPC_PROVIDER_WRITES_ENABLED`, `MPC_ML_CATALOG_OFFERS_ENABLED`,
  `MPC_ASSISTED_SANKHYA_LINKAGE_ENABLED` — all fail closed. These are real level-2 controls and §4
  leans on them.
- Oracle boot degradation to explicit `Unavailable` readers.
- The Sankhya/Oracle read path in `internal/adapters/erp/sankhyaoracle` — the approved adapter molde.

**Runtime and frontend**
- `sync.Scheduler.safeInvoke` (`scheduler.go:201-208`) and `apierror.Recover`.
- `/sync/health` plus `SyncHealthCard.tsx` and `SyncHealthCard.test.tsx`.
- `tsconfig.base.json` `"strict": true`.
- Zero `fetch()` bypass in the FE; `any` nearly absent outside `ClassificationsPage.tsx`.
- `sdk-runtime` type-checks clean in isolation.
- `contexts/catalog/module.go`'s `New()` shape.
- `cache_composed_test.go`, `catalog_rls_role_test.go`.

**Process**
- `scripts/release.sh` and the Docker dev stack as *local* tools. F14 calls the script a CI bypass;
  that is true and it stays, because §4 does not put release on the required path.

---

## 4. Target gate topology

### 4.0 First: how work reaches the remote, and in what order

**No gate design is meaningful until this is answered.** A required check is a verifier the author does
not control at merge time. Today the merge does not happen at a place a verifier can stand.

The remote is reachable and the visibility question is decided, so this is now a **sequencing**
problem with one hard constraint: the operator wants the first remote execution to be already gated,
not a red run against 98 commits whose failures this audit has already catalogued. That constraint
plus the mechanics of rulesets fixes the order completely.

**P1 — Fix the identity split first.** Two accounts sit in the keyring: `gh` asks as
`leandrotcawork`, git authenticates as `developmentconexus-ops`. This is why the cicd lane concluded
the remote was gone. Until `gh` and `git` agree, every `gh`-mediated step below (creating the ruleset,
reading check status, opening PRs) reports on the wrong repository or fails, and — worse — reports
plausibly. Operator action, minutes. Delete or re-scope the `legacy` remote at the same time; it is
superseded and unreachable, and leaving a second `origin`-shaped ref in the config is how the split
survives.

**P2 — Build the gates locally, on a branch off local `main`.** All of axis A, on the local
checkout, against the 98 commits as they stand. Nothing is pushed.

**P3 — Baseline locally, and commit the baselines.** Run every instrument against HEAD, commit the
resulting counts as ratchet baseline files. Without this, the first remote run is red on 58 governance
violations, 44 archscan findings, 12 `tsc` errors and 22 unformatted files, and a permanently red
required check is switched off within a week. The baseline commit is what makes the first green run
possible; it is not bookkeeping.

**P4 — Push once, while the branch is still unprotected.** Fast-forward `main` (98 + gate commits) in
a single push. This is deliberate and it is the last ungated write to `main` in the program's history,
by construction: a ruleset that blocks direct pushes to `main` would block this very push, so the
ruleset is enabled *after* it, not before. The pushed tree already contains the workflow, so the push
itself triggers the first run — which is the operator's requirement, satisfied by the gate being *in*
the pushed tree rather than by protection being on before it.

**P5 — Confirm green on the remote before anything else.** A local green and a runner green are
different claims; the runner has a clean checkout, LF line endings, no `.gomodcache`, and no
warm caches. Expect one or two environment defects here and budget for them.

**P6 — Flip to public, gated on the secrets/PII sweep.** Going public exposes the entire history, not
the current tree, and the sweep now running separately is what says whether that history is safe. **Do
not enable any workflow that consumes a repository secret until the sweep returns**, and see §4.3:
this program's recommendation is that the repository never holds an ERP or provider credential at all,
which makes the gate topology independent of the sweep's verdict. What *does* depend on it is the flip
itself, and therefore P7.

**P7 — Enable the ruleset, with an empty bypass list.** Required status checks on `main`, block direct
pushes, require a pull request. Details and the bypass-list argument in §4.2.

**P8 — Change the landing path in doctrine.** `docs/HARNESS-PROFILE.md:279` forbids pushing without
explicit operator permission and `:77` defines the post-merge ladder on an already-integrated local
`main`. Chips must land as `branch → push → PR → required checks → merge`, and the doctrine document
must say so, or the checks sit beside the road exactly as `release-images.yml` does. **This is the
largest non-code item in the program and no lane costed it**, because no lane was pointed at the
harness. **[Arm B]** Confirm by reading `HARNESS-PROFILE.md` §merge and observing that no chip tip in
`git log` has ever been a remote ref. Standing operator authorization to push *chip branches* — never
`main` — is the minimum change.

P1 through P5 are prerequisites to everything. P6–P8 can trail axis G by a few days without harm, but
P8 must land before the second axis ships, or the gates built in P2 go back to gating nothing.

### 4.1 What fires where

| # | Rule | Level | Where it fires | Blocks / annotates | Negative fixture (ships in the same change) |
|---|---|---|---|---|---|
| G1 | Money has no float64 constructor on the path from DB to wire | **1** unrepresentable | Go type system: `exact.Money` with unexported field, no `FromFloat`, `MarshalJSON` emitting a string | Blocks — non-compiling | A `_test.go` in a `//go:build fixture_fail` file constructing `Money` from a float64; a meta-test asserts it does not compile |
| G2 | SDK types are generated, never authored | **1** | `openapi-typescript` output committed; the file is `// Code generated` and `.gitattributes`-marked | Blocks — a hand edit is overwritten and the diff test fails | A committed hand-edit to a generated line; the regen-diff job must fail on it |
| G3 | Tenant comes from the request, not from config | **1** | Repository constructors take a `TenantID` value type with no zero value and no default; `pgdb/tenant.go`'s fallback deleted | Blocks — non-compiling at all 70 `root.go` sites | A fixture constructing a repository with `TenantID{}`; must not compile |
| G4 | `MC_DEFAULT_TENANT_ID` fails closed like every other required var | **2** boot-fatal | `pgdb.LoadConfig()` — the same function that already fails closed on `MC_DATABASE_URL` | Blocks boot | `TestLoadConfigMissingTenant` asserting a non-nil error; the test must fail if the fallback is restored |
| G5 | The app's DB role cannot bypass RLS | **2** | Startup assertion querying `rolbypassrls` and `rolsuper` for `current_user`; refuses to serve | Blocks boot | An integration test booting against the superuser DSN and asserting boot failure |
| G6 | Provider writes / ML catalog / assisted linkage stay off unless explicitly enabled | **2** | Already implemented and correct — `MPC_PROVIDER_WRITES_ENABLED` et al | Blocks the action | Already exists. Do not touch |
| G7 | One command is the gate | **3** red build | `npm run gate` — identical invocation locally and in `ci.yml` | Blocks the PR | The gate's own smoke test: a fixture branch with one known violation of each wired rule; `gate` must exit non-zero |
| G8 | Ratchets shrink only | **3** | `ci.yml` job comparing measured counts to committed baselines | Blocks on increase; annotates the current number | A fixture adding one violation of each ratcheted class; the job must fail |
| G9 | Generated artifacts match their source | **3** | `ci.yml`: regenerate SDK + server types, `git diff --exit-code` | Blocks | A committed drift between `.openapi.yaml` and the generated SDK; must fail |
| G10 | Route ⇄ spec ⇄ SDK are the same set | **3** | A Go test enumerating registered routes, an OpenAPI parse, and an SDK method extract; three-way set equality **joined on `operationId`, not on path template** | Blocks | Delete one SDK method; the test names the missing operation. Also: a route registered but absent from the spec |
| G11 | Every guard has an input that makes it fail | **3** meta | A test that walks the guard inventory and requires each to have a registered failing fixture that is actually executed | Blocks the addition of an unproven guard | A guard registered with no fixture; the meta-check must fail. **This is the axis-G mechanism** |
| G12 | `internal/contexts` is governed | **3** | `Policy.psm1:302` walk root becomes the pair `(kind, id)` over both trees; `GOV_CONTEXT_UNREGISTERED` implemented | Blocks | The fixture at `lanes/cicd.md` F9 already exists and currently asserts a nonexistent code — implementing the code makes an existing dead fixture live |
| G13 | An exception entry with no live violation is a violation | **3** | Governance run: for each `temporary_exceptions` entry, assert the violation it excuses still occurs | Blocks | An exception for a path that does not exist — 3 of the 5 in `modules.json` qualify today |
| G14 | Every tenant-bearing query carries a tenant predicate | **3** | A Go AST checker over the 246 raw query sites, allow-listing the 2 global tables | Blocks | A query against a tenant-bearing table with no `tenant_id` predicate |
| G15 | `tsc --noEmit` is clean | **3** | `ci.yml`, root and every workspace | Blocks | The `MutationPreviewModal.tsx:210` missing-`onRetry` case *is* the fixture — fix it last, after the gate is wired, so the gate is proven red first |
| G16 | Every test file is reached by a command | **3** | A census check: enumerate `*_test.go` and `*.test.tsx`, compare against the union of what `gate` executes | Blocks | A new test file in an unreached directory; census must fail. Kills T-1, T-3, T-7, T-8, T-10, T-12 as a class |
| G17 | The integration lane reports counts | **3** | `-v` plus a parsed RUN/PASS/SKIP summary; zero-ran is a failure | Blocks | A run with every test skipped must exit non-zero (T-6) |
| G18 | `gofmt -l` is empty | **3** | `ci.yml` on `ubuntu-latest` | Blocks | Any unformatted file. Precondition: one `gofmt -w` commit over the 22 genuinely unformatted files |
| G19 | Transport errors are typed | **3** | A checker banning `strings.Contains(err.Error(), …)` and `strings.HasPrefix(err.Error(), …)` in `*/transport/**` | Blocks | A transport handler string-matching an error; must fail |
| G20 | `/healthz` reflects dependencies | **4** runtime assertion | The handler actually checks DB and Oracle reader state; `docker-compose.yml:46` already consumes it | Blocks container promotion | A test booting with a dead DB asserting non-200 |
| G21 | No background loop can kill the process | **4** | All 6 ticker loops wrapped in `sync.Scheduler.safeInvoke`'s existing pattern | Contains — panic logged, loop survives, HTTP stays up | A job that panics; the test asserts the scheduler survives and the failure was recorded |
| G22 | Every request has a correlation ID and one JSON log handler exists | **4** | Middleware plus `slog.SetDefault` with a JSON handler in `cmd/server` | Annotates — but makes every other failure findable | A request asserting the ID appears in the emitted line and in the response header |
| G23 | A render error breaks a panel, not the app | **4** | React `ErrorBoundary` at the route level and around each data panel | Contains | A component that throws; the test asserts siblings still render |

### 4.1a Workflow inventory — trigger, permissions, fork reachability

On a public repository these three columns are part of the gate's correctness, not deployment detail.
Anyone may fork and open a pull request; every run's logs are world-readable.

| Workflow | Trigger | `permissions:` | Reachable by an untrusted fork? |
|---|---|---|---|
| `ci.yml` — G7 (`npm run gate`), G8–G11, G15–G19 | `pull_request` (all branches) + `push` on `main` | `contents: read` only. Nothing else. | **Yes, by design.** Safe because a `pull_request` run against a fork gets a read-only `GITHUB_TOKEN` and **no secrets**, and the workflow needs neither |
| `gate-integrity.yml` — C1 mixed-change ban | `pull_request` | `contents: read`, `pull-requests: read` | Yes. Reads the changed-file list via the API with a read token; no write path |
| `release-images.yml` — existing | `push` on tags + `workflow_dispatch` | `contents: read`, `packages: write` | **No, and it must stay that way.** Tag pushes and manual dispatch are not fork-reachable. Registry credentials live here and nowhere else |

Three rules bind every workflow this program adds:

1. **`pull_request_target` is banned outright.** It executes in the base repository's context with a
   writable token and access to secrets, while checking out code the fork author controls. There is no
   requirement in this topology that needs it. A ban is cheaper than a review policy, and it is
   enforceable: add `pull_request_target` to the C1 mixed-change job's forbidden-pattern list so a
   workflow introducing it fails its own gate. Negative fixture: a workflow file containing the
   trigger; must fail.
2. **Default `permissions: contents: read` at the workflow level**, with any elevation declared
   per-job. GitHub's repository default for new public repos is already read-only, but stating it in
   the file makes the intent survive a settings change nobody remembers making.
3. **Never `echo` an environment.** Logs are world-readable. This is the same rule the audit itself
   ran under, now with a public consequence: no bare `printenv`, no `docker inspect`, no `set -x` over
   a block that touches configuration. The governance and archscan output is fine — it prints paths
   and code that are public anyway once the flip lands.

Fork-PR workflow approval should be set to **"Require approval for all external contributors"**, not
the default first-time-contributor setting. It costs one click per genuine outside PR — of which this
repository will have approximately zero — and it removes the entire class of "fork opens a PR to burn
minutes or probe the workflow surface."

### 4.2 Protecting the gate from its own author

The obvious hole, stated plainly: the author of nearly all code here is an agent with write access to
`.github/`. A required check the author can edit is not a separation of powers.

With rulesets available at zero cost on a public repo, the primary control is the enforced required
check itself — and its strength turns entirely on one setting.

**R1 — Ruleset on `main` with an empty bypass list.** Require a pull request, require the `ci.yml` and
`gate-integrity.yml` checks to pass, block direct pushes and force-pushes. **Do not add the repository
owner to the bypass list.** This is the whole control. GitHub applies rulesets to admins unless they
are explicitly exempted, so an empty bypass list means the operator is bound by the same checks as
anyone else, and the only route around a red check is to go and edit the ruleset — a deliberate,
separate, logged act rather than a merge button clicked at speed. That is what separation of powers
reduces to in a one-person system: not a second reviewer, but a second *deliberate step* on a
different surface. Negative fixture: attempt a direct push of a trivial commit to `main`; it must be
rejected.

**CODEOWNERS buys less here than it looks like it does, and I would rather say so than repeat the
suggestion.** `developmentconexus-ops` is a personal account, so there are no teams — CODEOWNERS can
name individual users only, and the only user is the owner. "Require review from Code Owners" then
means: an outside fork's PR needs the owner's review (already true, since only the owner can merge),
and the owner's own PRs — which includes everything the agents produce under the owner's credentials —
get **nothing**, because GitHub does not permit self-approval and the requirement is either unmeetable
or admin-bypassed. It is a two-party control asked to work in a one-party system. Add the file if you
like, for the day a second person appears; do not count it as protection today.

Two checks that do work, both level 3, both usable by one person:

**C1 — Mixed-change ban.** A required job reads the PR's changed-file list and fails if the diff
touches both a gate path (`.github/workflows/**`, `scripts/harness/Policy.psm1`,
`contracts/governance/**`, the ratchet baseline files, `.githooks/**`) and any file under `apps/` or
`packages/`. It also fails on any workflow introducing `pull_request_target` or elevating
`permissions:` beyond `contents: read` without a matching entry in a declared allow-list. Cost: ~40
lines. Negative fixtures, both required: a PR touching one workflow file and one handler; and a
workflow file containing `pull_request_target`. This forces gate weakening to be a *separate, visible,
single-purpose PR* — which, combined with R1's empty bypass list, means weakening the gate is a PR
that must itself pass the gate.

**C2 — Signed commits on gate paths (optional, defense in depth).** A job asserting that every commit
touching a gate path carries a valid signature from the operator's key. An agent with write access to
`.github/` cannot sign as the operator. This is *not* the primary control — the amendment is right
that degrading to credential custody would be a retreat, and R1 makes it unnecessary. But it closes
C1's residual two-PR path (weaken the gate in one PR, land the code in the next) at the cost of one
job and a signing key, and it is the only control in this document that an agent cannot satisfy at
all. Take it if the signing setup is already in place; skip it otherwise and revisit.

### 4.3 What CI cannot run, and why that is acceptable

Per the amendment, stated explicitly rather than proposed and quietly dropped:

- **The Oracle/Sankhya live lane** needs ERP credentials, `cgo`, and an Oracle client. Those
  credentials are not going into a runner. Stays local. **Level 5.**
- **The Docker dev-stack live-drive / browser lane** needs the composed stack and real provider
  credentials. Stays local. **Level 5.**
- **`harness:provider-write`** performs live ML writes. Never in CI. Stays local. **Level 5.**

Level 5 is not a control, so the honest accounting is: those three surfaces are *not* gated by CI —
they are gated by **G6**, the fail-closed environment variables that already exist and already work
(`MPC_PROVIDER_WRITES_ENABLED`, the Oracle vars, `MPC_ML_CATALOG_OFFERS_ENABLED`). That is a level-2
boot-fatal control, and it is the reason this gap is tolerable rather than alarming. The residual risk
is regression in ERP *read* mapping, which is real and unmitigated; the mitigation available in CI is
G7-executed **golden fixtures captured from live responses** (axis G, T-5), not a live connection.

**On a public repository this stops being a limitation and becomes a requirement.** The recommendation
is stronger than "do not run these lanes in CI": **do not add ERP, Sankhya, Oracle, Mercado Livre or
any provider credential as a repository secret at all.** Reasons, in order of weight:

- Every secret in a public repository is one misconfigured trigger away from disclosure, and the
  misconfiguration in question — `pull_request_target` — is the single most common workflow mistake in
  the ecosystem. C1 bans it, but a control that has to hold forever against an agent that edits
  `.github/` is a worse bet than not having the secret.
- Logs are world-readable. GitHub masks registered secret *values*, but it does not mask a connection
  string assembled at runtime, a stack trace containing a DSN, or an ERP row echoed by a failing
  assertion. The audit's own PII constraint applies to the runner.
- There is no upside. The only lanes that would consume them are the three declared level 5 above.

This has a useful side effect on §4.0 P6: because the gate topology consumes **no secrets**, the
required checks can be built, pushed and proven green *before* the secrets/PII sweep returns its
verdict. The sweep gates the public flip; it does not gate the gates. The one recommendation that does
depend on the sweep is R1, since rulesets on a free personal account require the repository to be
public — so a `BLOCK` verdict from the sweep postpones enforcement, not construction. **[Arm B]** The
existing `release-images.yml` already holds registry credentials; whether those are repository secrets
that survive the flip is a question for the sweep, and it is the one place where §4.1a's "no secrets"
column has an exception.

One narrower point: `TestModuleBoundaryADR023` and the arch scanners are pure static analysis over the
source tree. They need nothing but Go and a checkout. **They belong in CI today** — the fact that
they are currently invoked only by `bash scripts/arch-gate.sh` on a developer's machine is the entire
content of T-1, not a constraint.

### 4.4 Minutes

Cost is not availability, so: the whole non-integration gate is Linux-only and cheap.

- `ubuntu-latest`, 1× multiplier. Measured local times from the lanes: `arch-gate.sh` ~78s,
  `harness:unit` ~82s, `go test ./internal/...` ~68s cached. A cold CI run without caches lands
  around **5–7 minutes per PR** for build + vet + gofmt + go test + archscan + governance + tsc +
  vitest.
- The integration lane (~259s locally) runs against a `services: postgres:` container — free, adds
  ~5 minutes. Run it on `pull_request` only, not on every push to a branch.
- **`windows-latest` is 2× and is not needed.** The PowerShell harness runs under `pwsh` on Linux,
  and the CRLF problem below does not exist on a Linux checkout. No matrix.
- **No matrix builds.** One Go version, one Node version, the ones in `go.mod` and `package.json`.
- Estimated ~12 minutes per PR. **On a public repository, GitHub-hosted standard runners are
  unmetered**, so the bill is zero regardless of volume.

Unmetered is not a licence to be careless, for two reasons that survive the price being zero. Wall
time is still paid by the operator: a 25-minute gate is a gate people learn to merge around, and the
12-minute figure above is already at the edge of tolerable for a solo workflow — cache Go modules and
`node_modules`, and keep the integration lane off branch pushes. And on a public repo, run volume is
partly controlled by strangers; the fork-approval setting in §4.1a is what keeps that from becoming
someone else's decision.

### 4.5 D-52 / CRLF, re-sized

`ubuntu-latest` checks out LF. The 635 pure-CRLF `gofmt` failures **cannot occur there**; only the 22
genuinely unformatted files remain. So:

- **The gate question is settled and trivial:** one `gofmt -w` commit over 22 files, then G18 blocks at
  zero. No ratchet, no exceptions, no `.gitattributes` dependency.
- **The `.gitattributes` renormalize decision drops off the critical path entirely.** It is now a local
  developer-experience question: without it, a developer running `gofmt -l` locally on Windows
  (`core.autocrlf=true`, `.gitattributes` containing only `*.sh text eol=lf`) still sees 635 files and
  the local instrument stays useless. I would still do it — one `* text=auto` line plus one mechanical
  renormalize commit — because a local instrument nobody can read is how the gate stopped being run in
  the first place. But it is a convenience item, not a decision the operator needs to adjudicate. **This
  is a re-sizing of PHASE-0's open decision D-52, downward, from blocking to optional.**

### 4.6 Where CodeRabbit sits

The operator's goal names CodeRabbit explicitly, so: it is an **annotator**, permanently. Its output is
prose, and prose cannot be a required check without making every PR's fate a judgement call — which is
level 5 with a robot wearing it. Its correct position is upstream of the gate, as a **rule
generator**: every CodeRabbit finding the operator accepts becomes a level-1/2/3 rule with a negative
fixture under G11. That converts a review tool into gate material, which is the only way it makes bad
changes *hard to land* rather than merely *commented on*.

---

## 5. The sequence

Sizes are working days for one operator plus agents.

**Prerequisites, interleaved with axis A rather than preceding it.** §4.0's P1–P8 are not a separate
phase: P1 (identity split) is an hour and comes first; P2 and P3 *are* axis A's work, done locally;
P4/P5 (the single push and the first green remote run) fall at the end of A and are its acceptance
criterion; P6/P7 (public flip, ruleset) follow the sweep and add ~1 day; P8 (the doctrine change) is
half a day and must land before axis B ships. **The 98 commits stay local through P1–P3.** The first
thing the remote ever sees is a tree that already contains the gate.

| # | Axis | Days | Depends on | What it makes cheaper |
|---|---|---|---|---|
| 1 | **A** — verifier not on the path (incl. P1–P5) | 7–10 | — | Everything. Every later axis's fix is verifiable the day it lands instead of the day someone remembers to check |
| — | P6–P8: sweep verdict, public flip, ruleset, doctrine | 1–2 | A, and the secrets/PII sweep | Converts A's checks from annotations into blocks. Until P7 the checks run and report; they do not stop a merge |
| 2 | **G** — guard and proof same author | 1–2 | A (needs `gate` to run the meta-check) | Every guard in axes B–F. Without G11, each later axis ships guards that certify themselves |
| 3 | **B** — contract has four copies | 5–8 | A, G | C (auth needs one place to add a header), H (typed errors are what `<ErrorState>` needs), and every future endpoint |
| 4 | **C** — identity is a constant | 7–11 | A, G, B | Nothing downstream; it is the terminal risk item. **Jumps to position 1 if the service is network-reachable** |
| 5 | **F** — money is optional | 4–6 | A, G | Runs in parallel with C; no dependency between them |
| 6 | **E** — no durable failure channel | 3–5 | A | Makes C's and F's production incidents diagnosable. Could move earlier at low cost |
| 7 | **D** — hand-maintained ledgers | 3–4 | A, G | Cheapest after A, because A's ratchets already measure what D must shrink |
| 8 | **H** — no failure containment | 1–2 | A (for `tsc`) | Slot anywhere after A; it is filler-sized and independent |

**Total: 33–48 days including prerequisites.**

**One ordering constraint that is not a preference.** Between P4 (the push) and P7 (the ruleset) there
is a window in which the checks run and report but do not block, because rulesets on a free personal
account require the repository to be public and the flip waits on the secrets/PII sweep. That window
is a level-5 interval by definition and should be kept short and *named* — during it, the operator is
the enforcement. If the sweep returns `BLOCK` and the flip is postponed by weeks, the honest response
is a local `pre-push` hook running the same `npm run gate` as an interim, explicitly labelled as a
stopgap that the operator can skip with `--no-verify` and therefore is not a control. Do not let the
interim become the design.

**Why A is first, justified against the inverted evidence standard.** Not because the code is currently
organized that way — it is organized the opposite way, with instruments and no wiring. The argument is
a design principle with a named failure mode, plus measured scale. The failure mode: *a verification
program whose first deliverable is code produces N fixes and zero evidence that fix N+1 will not undo
fix 3.* The measured scale: 91% of Go test files and 15% of frontend test cases are currently reached
by no command (`lanes/testing.md` census), which means 91% of the repo's existing verification
investment is already paid for and yielding nothing. A returns that investment before spending new
money. And the operator's own sentence — *"it gets so much harder to send bad PRs"* — is a statement
about a mechanism, not about code quality; A is the only axis that is entirely that mechanism.

**Why G is second and not last.** G is cheap and it is the only axis that changes the *quality of every
subsequent axis's output*. Ship A's ratchets without G11 and the ratchets are guards written by the
agent whose work they measure — the exact machine-author failure mode the brief names. One to two days
spent here changes what the following 25 days produce.

**Why C is fourth and not first, with a stated fork.** C is the highest-severity axis in the report:
an unauthenticated API serving CPF/CNPJ and full addresses. It is fourth **only** because PHASE-0
records network reachability as `unverified`, and because C's fix is the most invasive in the program
(a type-level tenant threaded through 70 construction sites, plus RLS, plus a role change, plus a login
surface) — doing it before A means doing it without a verifier, which is how a 7-day auth change
becomes a 20-day auth change. **If the service is reachable from outside localhost, this ordering is
void and C.1 starts today.** That is the one thing in this document I would want confirmed before the
operator reads any further.

**Why D is last among the real axes.** Its findings are the least dangerous — a stale exception entry
harms nobody directly — and its fix is nearly free once A's ratchets exist, because the ratchet
already computes the number D needs. Doing D early means building measurement twice.

---

## 6. Inversion tests

One line per structural conclusion, per the brief's form. Tests 10 and 11 replace an earlier pair
written when required checks looked unreachable; they are revisions, not additions.

1. **A verifier must sit on the path the change takes.** *Survives an opposite CI vendor, an opposite
   branching model, and an opposite repo layout, because a verifier positioned off the route a change
   travels has an execution count of zero in every topology.*
2. **Baseline before enabling, never after.** *Survives an opposite set of rules and an opposite
   codebase, because a permanently-red required check is disabled by its owner in every organization,
   regardless of what it is red about.*
3. **One authored artifact, everything else generated.** *Survives an opposite direction of generation
   (spec-from-code instead of code-from-spec) and an opposite serialization format, because two
   independently maintained descriptions of one wire shape diverge under maintenance in every language.*
4. **Identity must be per-request.** *Survives an opposite storage engine, an opposite auth protocol,
   and an opposite tenancy model, because an enforcement predicate whose input is constant partitions
   nothing in any system.*
5. **An exception entry with no live violation is a violation.** *Survives an opposite rule set and an
   opposite exception format, because a hand-maintained list that nothing reconciles against reality
   drifts in one direction only, in every registry ever built.*
6. **A guard ships with an input that makes it fail.** *Survives an opposite language and an opposite
   test framework, because an assertion never observed to fail is indistinguishable from an assertion
   that cannot fail, under every possible implementation of both.*
7. **The exact type must be compulsory, not available.** *Survives an opposite numeric library and an
   opposite currency, because an optional correctness affordance is used at the rate developers
   remember it, which is not 100% in any codebase — measured here as 0% (`kernel/exact` has zero
   production callers).*
8. **Failure needs a durable channel before it needs a policy.** *Survives an opposite logging library
   and an opposite deployment target, because a log line written to a stream nobody collects and a
   discarded error are the same artifact in every runtime.*
9. **A render error must break a panel, not an application.** *Survives an opposite frontend framework
   and an opposite rendering model, because unbounded fault propagation converts any single component
   defect into a total outage in every architecture that lacks a boundary.*
10. **Separation of powers in a one-person team is a second deliberate act, not a second person.**
    *Survives an opposite hosting provider and an opposite review culture, because self-approval is
    unavailable in every system that has approvals, so the only enforceable boundary a lone author can
    build is one that requires them to change a different surface on purpose.*
11. **A gate must not hold a credential it does not need.** *Survives an opposite CI vendor and an
    opposite repository visibility, because a secret that exists can be disclosed by a
    misconfiguration and a secret that does not exist cannot, in every pipeline ever built.*
12. **A control that exists and does not fire is absent.** *Survives every opposite in this document,
    because it is a statement about effects and the entire inadmissible-evidence category is statements
    about artifacts.*

---

## 7. Where I disagree with the lanes

**0. `lanes/cicd.md` F2's explanation is wrong: `origin` is reachable, and the lane's own tooling
misled it.** Recorded first because it was established by operator measurement during this synthesis,
not by me, and because it is the clearest instance in the audit of a correct measurement producing a
wrong fact.

F2 reports that `origin` cannot be resolved and infers that branch protection is unreachable and that
the remote is stale. Measured: `origin/main` is `7df7d011` (2026-08-04), local `main` is 98 ahead,
0 behind, fast-forward possible, nothing diverged. The cause is an identity split — `gh` asks as
`leandrotcawork`, git authenticates as `developmentconexus-ops` — so `gh` was answering truthfully
about a repository the operator does not own while git was talking to the one they do. The `legacy`
remote genuinely is superseded and unreachable, which is what made the wrong reading plausible.

The lane's *conclusion* — no branch protection exists, nothing blocks a merge — is unaffected and
correct. What was wrong was the reason, and the reason is what a remediation is built from: a program
premised on an unreachable remote designs local hooks and credential custody, which is precisely the
retreat the second amendment forbids. **This is the strongest argument in the audit for §4.0's P1.**
An identity split does not produce an error; it produces a confident answer about the wrong object,
which is the machine-author failure mode arriving through the tooling instead of through the prose.

**1. `lanes/layering.md` L-03 is factually wrong about the mechanism, and the wrong remediation
follows from it.**

L-03 states that `internal/composition` "is never passed to `arch.ScanVendorTokens`" and is therefore
scanned by nothing, citing 43 vendor-token occurrences across 5 files as unmonitored.

`scripts/arch-gate.sh:30` reads:

```bash
for root in internal/kernel internal/contexts internal/adapters internal/composition; do
  if (cd "$SERVER" && go run ./internal/arch/cmd/archscan -root "$root"); then
```

`internal/composition` is in the list. It **is** scanned — and `lanes/cicd.md` F7 measured the result:
42 findings from that exact root, the dominant term in the 44 total. Two lanes measured the same tree
and one of them concluded it was unmeasured.

The *verdict* survives, but only under the brief's own rule — a control that does not fire is absent,
and `arch-gate.sh` is invoked by no command (T-1). The **remediation changes completely**: L-03 as
written implies building or extending a detector. The real work is wiring a script that already
exists and then shrinking the 42 it already reports. That is a difference of days, and it is the
difference between axis A and axis D.

**2. `lanes/duplication.md` DUP-6 is mis-sized upward, and the correct disposition is the opposite of
the one implied.**

DUP-6 reports 82 hand-written `for rows.Next()` loops with zero uses of pgx v5's `CollectRows` /
`RowToStructByName`, framed as duplication to be removed.

`RowToStructByName` maps columns to struct fields by **name, at runtime, via reflection**. The
hand-written loop's `rows.Scan(&a, &b, &c)` is checked at compile time for arity and for type. Adopting
`RowToStructByName` therefore trades a compile-time check for a runtime one — in a codebase whose
named failure mode is *fluent, internally consistent wrongness: the wrong noun used correctly
everywhere*. A renamed column or a reordered select is caught today by the compiler and would be
caught tomorrow by a query that returns zeros in production.

This finding should move to §2 and stay there permanently. 82 verbose-but-checked loops is the correct
state. The genuine duplication cost — the boilerplate — is real and is worth less than the check it
would cost. **This is a disagreement about direction, not about size.**

**3. The `GOV_MIGRATION_PREFIX` reconciliation exposes a mechanism neither lane named. [Arm B]**

`lanes/cicd.md` F6 counts **1** `GOV_MIGRATION_PREFIX` violation. `lanes/persistence.md` P-8 counts
**2** prefix collisions (`0021_integration_operation_run_evidence.sql` /
`0021_integrations_provider_auth_strategy_shopee_partner.sql`, and
`0093_orders_status_details_nullable.sql` / `0093_sync_state_market_queue_entity_split.sql`). Do not
average: both are correct. `contracts/governance/invariants.json:34-38` carries a
`migration-prefix-0021-duplicate` `temporary_exception` with `exception_mode: exact-path-set`, and
`Policy.psm1:478-484` matches the sorted path set against it. One collision is excused; the other is
reported.

The finding neither lane made is in the dates. The `0093` pair was committed in `2263a12f`
(2026-08-03). The `migration-prefix-0021-duplicate` exception was introduced in `5fdb88ed`
(2026-07-11, *"feat(governance): add executable repository contracts"*). **The exception predates the
second collision by three weeks.** It was written to excuse a specific historical pair, it has no
expiry, no owner, and nothing has ever asked whether it is still needed — and in the interval a second
instance of the very defect it excuses landed with no friction.

That is `lanes/layering.md` L-01's decay pattern reproduced in a second, independent registry, which
is what promoted it from "one stale entry" to the naming of axis D and to gate G13. Confirm with
`git log --diff-filter=A -- apps/server_core/migrations/0093_*` and
`git log -S migration-prefix-0021-duplicate -- contracts/governance/invariants.json`.

**4. `lanes/testing.md` T-2 is right and `lanes/layering.md`'s "actually fine" verdict on
`TestModuleBoundaryADR023` is wrong.**

Layering lists the ADR-023 detector among the things that work. Testing lists it as having no positive
control, no `testdata/`, and no must-fail proof. Layering's evidence for "fine" is the detector's
*output* — it produces sensible findings — which is reasoning from the artifact, the exact category
the brief rules inadmissible.

I verified the disputed fact directly: `find apps/server_core/internal/composition -type d -name
testdata` returns **empty**, while `internal/arch/testdata/{clean,violations,adapters}` exists for the
sibling scanner. A detector with no failing input has never been shown to fail, and this one has a
known blind spot (L-04: it matches only the literal `.../internal/modules/` prefix, so modules↔contexts
coupling is invisible to it). Side with testing, decisively. It belongs in axis G, not on the
do-not-touch list.

**5. Not a disagreement — a reconciliation `lanes/frontend.md` left open.**

The frontend lane reports the SDK's interface count as both 172 and 220 without resolving it.
`grep -c "^export interface"` returns **172**; `grep -cE "^export (interface|type)"` returns **220**.
The difference is exactly the 48 exported type aliases. PHASE-0 established fact 5 uses the interface
count and **stands unchallenged** — I raise this only so the number does not get re-litigated as a
contradiction of an established fact when it is a difference in what was counted.
