# Phase 4 — Issue set

Eleven issues for ten axes (I1 splits into an edge carve-out and the structural body — they land at
different times behind different blockers). One issue per **axis**, never per finding: child findings
live inside as sized evidence.

**Not yet filed.** `gh` authenticates as `leandrotcawork` while git authenticates as
`developmentconexus-ops` — every `gh`-mediated call today reports on the wrong repository, *and
reports plausibly*. File these after **P1** (`GATE-TOPOLOGY.md` §1), with `gh issue create --title …
--body-file …`. Filing them is an outward-facing action; it needs the operator's go-ahead at the time.

Titles name the **cause**, so someone who has never seen this codebase understands what is broken.

Sources: `RECONCILIATION.md` (axes, sequence), `GATE-TOPOLOGY.md` (gates), `lanes/*.md` (measurements).

---

## Sequence at a glance

| Order | Issue | Size | Independent? |
|---|---|---|---|
| 0 | **#1 I1-edge** — a PII endpoint is reachable with no identity check | hours | yes — blocks the public flip |
| 1 | **#2 V2** — the verifier is scripts to remember, not one product | 4–6d | no — the multiplier |
| 2 | **#3 V3** — evidence certifies itself | 3–4d | after #2 |
| 3 | **#4 V1** — nothing is required to pass through the delivery path | 1d | after #3 |
| 4 | **#5 I1-structural** — identity is a boot-time constant | 5–8d | after #4 |
| 5 | **#6 C1** — the published contract is six hand copies with no arbiter | 4–6d | after #4 |
| 6 | **#7 R1** — every request boundary is re-invented per module | 5–7d | after #6 |
| 7 | **#8 M1** — there is no single compiled vocabulary for a value | 7–10d | after #4 |
| 8 | **#9 B1** — the boundary instruments are anchored to the tree being abandoned | 4–5d | after #4 |
| 9 | **#10 O1** — a failure leaves no durable trace | 4–6d | after #4 |
| 10 | **#11 F1** — the browser has no failure containment | 1–2d | **yes** — any time after #2 |

**Genuinely parallel:** #1 runs now, beside everything. #11 runs beside any of #5–#10. #9 and #10
have no dependency on each other. Everything else is serialised by the reason stated in its
Sequencing section — and never by "it would be tidier."

~34–49 agent-days. Both synthesis arms landed within three days of this independently.

---

# #1 — I1-edge: a PII-bearing endpoint is reachable with no identity check

**Root cause.** `/orders` serialises buyer name, tax ID and address, and the middleware chain that
wraps it contains no identity check — so the only thing currently limiting access is that nobody
knows the route. Publishing the repository publishes the map.

**Evidence (sized).**
- `orders/transport/http_handler.go:608-618` — the handler serialises buyer name, CPF/CNPJ, address.
- `composition/root.go:994` — the chain is `CORSMiddleware(apierror.Recover(mux))`. No identity
  middleware anywhere in it.
- `deploy/Caddyfile` — the exact path predicate routing `/orders` with a non-HTML `Accept` header to
  `backend:8080`. This file goes public with the repository.
- `docker-compose.yml` profile `oauth` — an ngrok tunnel publishing the frontend.
- `apps/web/vite.config.ts` — the full proxy table.

**Deliverable.** The route is not reachable without an identity check. At the edge: a deny rule in
`deploy/Caddyfile` for the `@orders_api` matcher, plus the ngrok tunnel scoped so it does not
publish the API surface. This is a stopgap by design — #5 makes it structural.

**Explicitly not in scope.** Authentication. Sessions. Principals. Per-request tenancy. This issue
closes a door; it does not decide who has keys. Any change to `composition/root.go` beyond what the
deny rule requires belongs in #5.

**Acceptance.** Level 3 for the edge rule — a check that fetches the route through the composed
stack and asserts non-200 without credentials. Not higher, because a Caddy config is not
representable in the Go type system. **Arm A withdrew its own first proposal here** — a `127.0.0.1`
bind — after determining it would read as done and be inert. Do not re-propose it.

**The negative fixture is mandatory, and that changed on 2026-08-07.** As written, this issue was a
stopgap that `GATE-TOPOLOGY.md` L2-c would retire once #5 shipped. The operator has deferred
authentication — one human user, platform first — so L2-c is not coming, and **this is now the only
control on the PII route, indefinitely.** A Caddy config and an ngrok scope, with no type and no boot
condition behind them. Contained while the repository is private; sharper at P6, which publishes the
route, its predicate and the fields it returns with a deny rule as the only thing between them.
So the fixture is not a nicety: **a request through the composed stack asserting non-200 without
credentials, wired into `verify-full`.** Without it, "the door is closed" is a claim about a config
file nobody re-reads.

**Sequencing.** Depends on nothing. **Blocks P6, the public flip.** Hours.
*Correction recorded:* Arm B sequenced the flip before this and Arm A after; resolved in favour of A
(`RECONCILIATION.md` Divergence 1). The exposure is already live — the flip does not create it, it
removes the only thing currently limiting who knows about it.

---

# #2 — V2: the verifier is a set of scripts to remember, not one product with one entry point

**Root cause.** Checks exist and are good; nothing invokes them as a set, so which ones ran on any
given change is a function of what someone remembered.

**Evidence (sized).**
- 369 `*_test.go` files are not reached by any command anyone runs.
- `apps/web/vitest.config.ts:14` pins test discovery to filenames, so new test files are invisible.
- An orphaned root vitest config discovers nothing.
- `tsc --noEmit` is wired to no command. 12 errors sit at HEAD, 3 in production code — including
  `MutationPreviewModal.tsx:210`, whose missing `onRetry` has been live long enough to prove the
  checker has never run.
- `scripts/arch-gate.sh` is referenced by zero pipelines. 44 archscan findings at HEAD.
- Governance: 58 violations at HEAD with a zero diff.
- ADR-023 module boundary detector: 234 violations. Pure static analysis — needs only Go and a checkout.
- Pester: ~250 assertions across `scripts/tests/*.tests.ps1`, 6 of 11 files red from cold.
- `gofmt`: 22 genuinely unformatted files (the other 635 are CRLF artifacts that cannot occur on a
  Linux checkout).

**Deliverable.** One command — `npm run gate` — that runs every instrument above, invoked
**identically** locally and in `ci.yml`. Plus `ci.yml` itself with the `verify-fast` / `verify-full`
job split of `GATE-TOPOLOGY.md` §2.3, and the ratchet baselines of P3 committed.

**Explicitly not in scope.** Fixing the violations the instruments find. 44 + 58 + 234 get
**baselined**, not paid down — a permanently red required check is switched off within a week, and
that failure mode is more expensive than the violations. The two exceptions are cheap and bounded:
22 `gofmt` files (one commit) and 12 `tsc` errors (one change) get paid to zero, because a ratchet on
a number that small is more machinery than the fix.

**Acceptance.** Level 3. Not higher — "the checks ran" is not expressible in a type system and is not
a boot condition. It goes as high as it can go, which is a red build.

**Sequencing.** Depends on nothing. **Makes every later axis cheaper**, which is why it is first:
every later fix lands as *mechanism* if its gate fires before handoff, and as *discipline* if it does
not — and discipline regresses. This argument survives inversion (§below) because it does not depend
on what the code currently looks like. 4–6 days.
*Inversion:* survives an opposite codebase entirely. A verifier that nothing invokes has no effect
regardless of what it would have found.

---

# #3 — V3: evidence certifies itself

**Root cause.** Guards are written together with the tests that bless them, so a green result is
consistent with the guard never firing, the test never running, and the assertion naming a condition
the engine does not implement.

**Evidence (sized).**
- The catalog RLS guard and its only test were authored in the same commit (`0098` / `7e3dcc47`).
  Same for the tenant guard (`47a76837`). Neither has an input proven to make it fail.
- The ADR-023 detector has **no `testdata/`** anywhere under `internal/composition`. Nothing proves
  it can report a violation.
- The integration lane's all-skipped run is **byte-identical** to a green one. There is no count.
- A negative fixture at `lanes/cicd.md` F9 asserts the error code `GOV_CONTEXT_UNREGISTERED`, which
  the governance engine does not implement — so the fixture passes by never matching.
- `set -euo pipefail` does not fire on a status tested by `if`, so `scripts/arch-gate.sh` never exits
  early on a failing instrument (D-51).

**Deliverable.** Three mechanisms: (a) a **guard inventory** where every guard is registered with a
failing fixture, and a meta-check that walks the inventory and requires each fixture to be
*executed*, not merely present; (b) **attributable counts** on every gate step, where a step
reporting zero executed units fails; (c) `failure_token=test=` on the integration lane, with
RUN/PASS/SKIP/FAIL parsed and a zero-ran run exiting non-zero.

**Explicitly not in scope.** Auditing the correctness of every existing guard's logic. This issue
proves each guard **can fail**, not that it fails on the right things. A guard that is wrong but
demonstrably alive is a smaller problem than one that is right and inert, and it is findable later.

**Acceptance.** Level 3 meta. Not higher: a guard's liveness is a property of a test run, not of a
type or a boot.

**Sequencing.** Depends on **#2** — there is no inventory until there is one command that runs
things. **Must precede #4**, or the first required check certifies itself, which is the exact defect
this issue names. 3–4 days.
*Inversion:* survives an opposite implementation. A guard with no input proving it can fail is
indistinguishable from an absent guard in any codebase.

---

# #4 — V1: nothing is required to pass through the delivery path

**Root cause.** Work reaches `main` by local merge. No commit in this repository's history has ever
passed a check it could not skip, because there is nowhere in the path a check could stand.

**Evidence (sized).**
- 98 commits sit on local `main`, ahead of `origin/main` (`7df7d011`) by `0 98`. Fast-forward is
  possible; none of them has ever been through a PR.
- Zero required status checks. Zero rulesets. `release-images.yml` is the only workflow, and it
  triggers on tags and dispatch — it sits beside the delivery path, not on it.
- `docs/HARNESS-PROFILE.md:77` defines the post-merge ladder over an **already-integrated local
  `main`** — the doctrine describes a path with no gate position in it.
- `docs/HARNESS-PROFILE.md:279` forbids pushing without explicit operator permission, so chip
  branches never reach the remote where a check could see them.

**Deliverable.** `GATE-TOPOLOGY.md` §1 P4–P8 executed: push once while unprotected, confirm green on
the runner, flip to public, enable the ruleset with an **empty bypass list**, and **amend
`HARNESS-PROFILE.md` so the landing path is `branch → push → PR → required checks → merge`.** Plus
`gate-integrity.yml` (the mixed-change ban) and verification that a missing required check blocks as
*expected, not received*.

**Explicitly not in scope.** CODEOWNERS as a control. It is a two-party mechanism in a one-party
system: no teams on a personal account, and GitHub forbids self-approval, so the requirement is
either unmeetable or admin-bypassed. Add the file for the day a second person exists; do not count
it. Also not in scope: signed commits (deferred, `GATE-TOPOLOGY.md` §8).

**Acceptance.** Level 3, enforced by platform. Its own negative fixture is mandatory and must
actually be run: **open a PR deleting `ci.yml` and confirm the merge is blocked.** Without that, the
cheapest route past every gate in the program is `rm .github/workflows/`.

**Sequencing.** Depends on **#3**. The doctrine amendment (P8) is the largest non-code item in the
program and **no lane costed it, because no lane was pointed at the harness** — it must land before
#6 ships or the gates built in #2 go back to gating nothing. 1 day plus the doctrine change.
*Inversion:* survives an opposite implementation. A merge the author can perform without passing
anything is ungated in every architecture.

---

# #5 — I1-structural: identity is a boot-time constant, so every downstream enforcement enforces a constant

**Root cause.** The tenant is read from configuration at startup, so row-level security, tenant
predicates and per-tenant scoping are all enforcing the same constant for every request — and the
enforcement layers below it cannot be correct no matter how carefully they are written.

**Evidence (sized).**
- `pgdb/config.go:23-24` — `tenant_default` fallback. A missing tenant becomes a *value*.
- 70 construction sites in `composition/root.go` receive it.
- 246 raw query sites; the tenant predicate is applied by convention.
- The application's DB role can bypass RLS, which makes the policies inert against the connecting
  role (fact 12 / D-44).
- `pgdb.DefaultTenantID` — a dead helper that keeps the concept alive (SEC-8).
- No identity middleware exists in the composed chain (`composition/root.go:994`).

**Deliverable.** `TenantID` as a value type with no usable zero value, resolved **per request**;
the `tenant_default` fallback deleted; `MC_DEFAULT_TENANT_ID` failing closed in `pgdb.LoadConfig()`
exactly as `MC_DATABASE_URL` already does; a boot assertion that the connecting role is neither
`BYPASSRLS` nor superuser nor table owner; an AST checker over the 246 query sites allow-listing the
2 genuinely global tables; and a boot assertion that no route in the PII set is composed without the
identity middleware (`GATE-TOPOLOGY.md` L2-c) — which retires #1's stopgap.

**Explicitly not in scope, and this is the load-bearing exclusion.** **Authentication itself is a
product decision, not an audit output.** Who are the principals, what is a session, is there ever
more than one human user — both synthesis arms declined to size it, and correctly. This issue makes
identity *representable and required*; it does not choose an auth scheme.

**Operator decision 2026-08-07: authentication is deferred. One human user; make the platform work
first.** This does **not** stall the issue, because most of it never needed auth. `TenantID` with no
usable zero value, deleting the `tenant_default` fallback, `MC_DEFAULT_TENANT_ID` failing closed, the
RLS-bypass boot assertion and the 246-site query checker are all about a tenant being **explicit and
fail-closed rather than invented** — and one tenant is a perfectly good number of tenants for that.
The single piece that genuinely waits is the L2-c boot assertion (no PII route composed without
identity middleware), because there is no middleware to assert on. **Estimate drops from "5–8d plus
unsized auth" to 5–8d.** The deferred half moves to #1, which is now permanent until auth exists.

**Acceptance.** Level 1 for the tenant type (`TenantID{}` must not compile at all 70 sites), level 2
for the four boot assertions, level 3 for the query checker. Three levels because three different
things are being prevented — the highest each one admits.

**Sequencing.** Depends on **#4** (so the boot assertions land as mechanism) and on the auth decision
for its second half. 5–8 days plus unsized auth.
*Inversion:* survives an opposite schema entirely. **Row-level security is ineffective for a
bypassing connection in every topology**, and a tenant read once at boot is a constant regardless of
how the rest of the system is arranged.

---

# #6 — C1: the published contract is six hand copies with no arbiter

**Root cause.** One truth — the shape of a request and the set of routes — is transcribed by hand
into six places, and nothing compares any two of them. `GOV_API_SDK_SPLIT` requires only that they
change in the *same commit*, never that they *agree*.

**Evidence (sized).**
- The SDK is 100% hand-written: 2,595 lines, 172 interfaces.
- The same shape lives in `domain` → DTO → OpenAPI → SDK. Four copies for shape.
- Two more for routing: `deploy/Caddyfile` and `apps/web/vite.config.ts` maintain duplicate route
  tables, under a comment instructing a human to keep them in sync, with the Go router as a third
  source of truth.
- 3 OpenAPI operations are genuinely unreachable from the SDK.
- 6 of the 13 Fatia A findings were this one cause wearing different clothes.

**Deliverable.** `oapi-codegen` for server types and `openapi-typescript` for the SDK — both approved
2026-08-07 — with `git diff --exit-code` on the regenerated tree in CI; a three-way route ⇄ spec ⇄
SDK set-equality test **joined on `operationId`, not on path template**; and both routing tables
emitted from one authority.

**Explicitly not in scope.** Redesigning the API. This issue changes *how many copies exist*, not
what they say. Any operation whose shape is wrong stays wrong and gets its own issue — mixing the
two makes the regeneration diff unreviewable, which is how a generated-code migration silently
changes behaviour.
*Correction recorded:* the `delivery` lane's F14 claimed 10 unreachable operations. Four are
reachable at `market.ts:179,182,185,190`, and three are OAuth browser redirects where an SDK method
would itself be a defect. **The real gap is 3.** The refuting fact sat in `lanes/frontend.md:39` the
whole time and neither lane checked the other.

**Acceptance.** Level 1 — the hand copy stops existing rather than being compared. The three-way
equality test is level 3 and covers the residue that generation cannot reach (routes registered in Go).

**Sequencing.** Depends on **#4**. **Blocks #7** — a request boundary is much cheaper to provide once
when the types on both sides are generated. 4–6 days.
*Inversion:* survives an opposite API design. Several hand transcriptions of one truth with nothing
comparing them drift in any shape of contract.

---

# #7 — R1: every request boundary is re-invented per module instead of being provided once

**Root cause.** Each module writes its own decode, validate, error-map and respond path, so the
boundary's behaviour is a per-module accident rather than a property of the system.

**Evidence (sized).**
- Error mapping is by string matching on `err.Error()` in transport handlers — `strings.Contains` and
  `strings.HasPrefix` against message text.
- No typed error registry; no exhaustive mapping from code to HTTP status.
- Response shapes vary per module; the same failure surfaces differently depending on which handler
  it passed through.

**Deliverable.** One provided request boundary: a typed error registry with an exhaustive switch
(so a code with no HTTP mapping **fails to build**), one decode/validate/respond path used by every
transport package, and a checker banning `strings.Contains(err.Error(), …)` and
`strings.HasPrefix(err.Error(), …)` under `*/transport/**`.

**Explicitly not in scope.** React error boundaries. A render error boundary is not a request
boundary — the conflation is verbal, and merging them produces one issue whose "done" is
unfalsifiable. That work is **#11**.
*Correction recorded:* Arm A folded browser containment into this axis; Arm B split it out. Resolved
in favour of B (`RECONCILIATION.md` Divergence 4).

**Acceptance.** Level 1 for the error registry (an unmapped code does not compile), level 3 for the
string-matching ban.

**Sequencing.** Depends on **#6** — providing the boundary once is mechanical when the types on both
sides are generated, and hand-editing hundreds of call sites when they are not. 5–7 days.
*Inversion:* survives an opposite module layout. A boundary re-implemented per module has per-module
behaviour under any decomposition.

---

# #8 — M1: there is no single compiled vocabulary for a value

**Root cause.** Money and enumerated values have no compulsory representation, so every layer picks
one — `float64`, `string`, a `CHECK` constraint, a TypeScript union — and nothing forces agreement.

**Evidence (sized).**
- 373 `float64` occurrences on paths that carry money.
- `strconv.ParseFloat` at the ingest boundary.
- 65 `CHECK (… IN (…))` constraints encoding vocabularies that also exist in Go and in TypeScript,
  with nothing comparing the three.
- ~78 files in the conversion blast radius.
- **`lineTotal(1.005, 1)` returns `"1.00"`** — a real defect, currently unreachable only because
  `unit_price` is `numeric(14,2)` at the database. **No line of code asserts that premise**, and the
  in-memory ingest path bypasses it.

**Deliverable.** An `exact.Money` type with an unexported field, no `FromFloat`, no `Float64()`, and
`MarshalJSON` emitting a string; Postgres enum types replacing the 65 `CHECK` constraints; the
contract types (`OrderDecomposicao`, `PricingDecomposition`) carrying the exact representation rather
than `double`/`string` by accident.

**Explicitly not in scope.** Changing any monetary *calculation*. This is a representation change and
must be provably behaviour-preserving; a formula correction landing inside a 78-file conversion is
undetectable in review. Also not in scope: replacing the 82 hand-written `rows.Next()` loops with
`RowToStructByName` — see the do-not-touch list below.

**Acceptance.** Level 1. Constructing money from, or extracting it into, a binary float stops
compiling rather than being flagged. Negative fixture in a `//go:build fixture_fail` file with a
meta-check asserting `go build` returns non-zero.

**Sequencing.** Depends on **#4**, and **must not run while 91% of the test suite is unreachable** —
both arms placed it late for this reason and no other. 7–10 days.
*Inversion:* survives an opposite type layout. A value with several representations and no arbiter
drifts between them in any architecture; that is what "no single vocabulary" means.

---

# #9 — B1: the boundary instruments are prefix-anchored to the tree being abandoned

**Root cause.** The governance and architecture scanners walk `internal/modules/<id>` by path prefix,
so a module that lives anywhere else is not *failed* — it is **not seen**, which reports as zero
findings.

**Evidence (sized).**
- `scripts/harness/Policy.psm1:308` cements `internal/modules/<id>` as the walk root.
- A directory with no registry entry yields **zero findings, not a failure**.
- `id: "catalog"` already exists, so the registry key must become the pair `(kind, id)`.
- `internal/contexts` is ungoverned; `GOV_CONTEXT_UNREGISTERED` is asserted by an existing fixture
  but not implemented by the engine.
- 5 `temporary_exceptions` entries in `modules.json`; **3 no longer correspond to a live violation.**
  `migration-prefix-0021-duplicate` (2026-07-11) predates by three weeks the 2026-08-03 collision it
  does not cover.
- 44 archscan findings across `contexts` / `adapters` / `composition` at HEAD.

**Deliverable.** The walk root becomes the pair `(kind, id)` over both trees; `GOV_CONTEXT_UNREGISTERED`
implemented, which turns an existing dead fixture live; `scripts/arch-gate.sh` wired into CI over
**all four roots**; and an exception-liveness check asserting that every `temporary_exceptions` entry
still excuses a violation that actually occurs.

**Explicitly not in scope.** Moving modules. This issue makes the instruments see both trees; it does
not decide which tree anything belongs in.
*Correction recorded:* the `layering` lane's L-03 claimed `internal/composition` is unreachable by
every instrument. **`scripts/arch-gate.sh:30` does scan it**, and the `cicd` lane measured 42–44
findings from that root. The verdict survives only under *a control that does not fire is absent* —
but the remediation changes from **build a detector** to **wire a script**, which is materially
cheaper. Arm B caught the contradiction; neither lane checked the other.

**Acceptance.** Level 3. A registry walk is not expressible in a type system and is not a boot
condition.

**Sequencing.** Depends on **#4**. Independent of #10. 4–5 days.

**Open operator decision blocking part of this — D-2 / D-50, stated concretely.** Every governance
exception must name a `removal_owner`. Two of the three schemas
([`modules.schema.json:49`](contracts/governance/schemas/modules.schema.json:49),
[`runtime-config.schema.json:67`](contracts/governance/schemas/runtime-config.schema.json:67)) force
that value to be `M-NN` or `M-NN/F-NN` — **a milestone. Something scheduled. Something that ends.**
[`invariants.schema.json:54`](contracts/governance/schemas/invariants.schema.json:54) alone widens
the pattern with a fourth alternative, `HARNESS-D-[0-9]{1,3}`, pointing at a row in
`.mnfs/HARNESS-DEBTS.md`. Four exceptions use it: D-9, D-40, D-41, D-42.

A milestone is a commitment; a debt row is a list. So `HARNESS-D-N` produces exceptions that read as
temporary and are permanent — [`invariants.json:62`](contracts/governance/invariants.json:62) says
removal *"= the constructor ceasing to exist (HARNESS-DEBTS D-40)"*, which describes what removal
would look like rather than naming anyone who will do it. **That is the same shape as the three stale
exceptions this issue already found**, and the liveness check shipped here cannot distinguish
"permanent by design" from "abandoned."

The ruling needed, per exception:

- **Permanent by design** — clearest candidate is D-9
  ([`invariants.json:30`](contracts/governance/invariants.json:30)), re-raising
  `http.ErrAbortHandler`, whose own reason argues the panic is *correct* and that converting it would
  be the defect. Then `HARNESS-D-N` is the wrong slot entirely: this is **a rule the checker does not
  know yet**, and it belongs in the checker as a sanctioned idiom, not excused per-site forever.
- **Deferred work** — then each gets a real `M-NN` and the fourth schema alternative is deleted,
  closing the escape hatch.

Either ruling makes the liveness check unambiguous. Neither is urgent; it can wait until this issue
starts.
*Inversion:* survives an opposite directory layout. **An instrument anchored to a path prefix reports
zero for everything outside that prefix in any tree shape** — and zero findings is indistinguishable
from compliance.

---

# #10 — O1: a failure leaves no durable trace

**Root cause.** When something goes wrong there is no record that ties it to a request, so every
diagnosis starts from reconstruction rather than from evidence — and a background failure may leave
nothing at all.

**Evidence (sized).**
- No correlation ID on any request.
- No `slog.SetDefault` with a JSON handler in `cmd/server`; log output is unstructured.
- 6 ticker loops where a panicking job can take the process down.
- `/healthz` does not check its dependencies, while `docker-compose.yml:46` already consumes it for
  container promotion — so a container with a dead database is promoted as healthy.
- A failed ML token refresh is invisible at every layer and the screen stays green.

**Deliverable.** A correlation ID in middleware, in every emitted log line, and in the response
header; `slog` with a JSON handler configured at boot and **fatal if the destination is unwritable**
rather than silently degrading to text; all 6 ticker loops wrapped in the existing
`sync.Scheduler.safeInvoke` pattern so a panicking job is logged and the loop survives; `/healthz`
actually checking DB and Oracle-reader state.

**Explicitly not in scope.** Metrics, tracing, dashboards, alerting. This issue makes a failure
*findable after the fact*. Deciding what to watch continuously is a separate decision with a separate
cost, and bundling it is how an observability issue becomes a platform project.

**Acceptance.** Level 2 for the log-handler boot condition, level 4 for the correlation ID, the
scheduler containment and `/healthz`. Not level 3 — none of these is a property of the source tree;
they are properties of a running process, which is exactly why level 4 is the ceiling here.

**Sequencing.** Depends on **#4**. Independent of #9. 4–6 days.
*Inversion:* survives an opposite implementation. A failure with no durable trace is undiagnosable
regardless of how the code that failed is organised.

---

# #11 — F1: the browser has no failure containment

**Root cause.** A render error anywhere takes down the whole SPA, because there is no boundary
between a panel and the application.

**Evidence (sized).**
- No React `ErrorBoundary` at the route level or around any data panel.
- `MutationPreviewModal.tsx:210` renders `ErrorState` without the required `onRetry`.

**Deliverable.** An `ErrorBoundary` at the route level and around each data panel, so a render error
breaks a panel and leaves its siblings alive. Negative fixture: a component that throws, with the
test asserting siblings still render.

**Explicitly not in scope.** Request-boundary work — that is **#7**. Also not in scope: redesigning
`ErrorState`.
*Correction recorded, and this one was also stated wrongly to the operator in conversation and has
been corrected there.* The `frontend` lane's F3 called the missing `onRetry` a crash that downs the
SPA. **It is not.** `packages/ui/src/ErrorState.tsx` passes `onClick={onRetry}` with no invocation
guard, and React ignores an undefined handler — the result is a **permanently dead retry button**,
not a throw. Re-rated accordingly. Its value as live proof that `tsc` has never run is undiminished,
and it is the negative fixture for the `tsc` gate in #2: wire the gate, watch it go red on this
exact line, *then* fix it.

**Acceptance.** Level 4 — containment, not prevention. A render error cannot be made unrepresentable;
it can be stopped from propagating.

**Sequencing.** Depends only on **#2**. **Genuinely independent** — schedulable beside any of #5–#10.
1–2 days.
*Inversion:* survives an opposite component tree. An uncontained render error propagates to the root
in any React application.

---

## Consolidated do-not-touch list

Merged from every lane's "what is actually fine" section. These are **decisions**, not omissions —
recorded so they are not re-proposed in six months.

| Item | Ruling |
|---|---|
| The 82 hand-written `rows.Next()` loops | **Keep them, forever.** Replacing them with `RowToStructByName` trades a compile-time arity and type check for name-based runtime reflection — in a codebase whose named failure mode is *the wrong noun used correctly everywhere.* Rejecting this is consistent with the program's own goal. |
| `MPC_PROVIDER_WRITES_ENABLED`, `MPC_ML_CATALOG_OFFERS_ENABLED`, the Oracle fail-closed vars | **Already correct. Do not touch.** They are the level-2 template every new boot gate should copy, and they are the entire reason the CI gap in `GATE-TOPOLOGY.md` §5 is tolerable rather than alarming. |
| `release-images.yml` as written | Correct. Tag + dispatch triggers are not fork-reachable, and its `packages: write` scope is precisely why it must never gain a `pull_request` trigger. |
| `tsconfig.base.json` `"strict": true` | Already at level 1. The setting is not the problem; nothing invoking the checker is (#2). |
| CODEOWNERS | Not a control on a personal account. Add the file for a future second person if you like; never count it. |
| `.gitattributes` renormalize | Optional. Local developer experience only — the CI checkout is LF and the 635 CRLF `gofmt` failures cannot occur there. D-52 is **re-sized downward from blocking to optional**. |
| ERP / Sankhya / Oracle / Mercado Livre credentials as repository secrets | **Never add them.** No lane in this topology consumes them; every one that would is level 5 by design. |

---

## Corrections recorded by this audit

Written down rather than silently applied, because the reasoning is the durable artifact.

| Claim | Correction |
|---|---|
| `frontend` F3: missing `onRetry` crashes the SPA | **Wrong.** React ignores an undefined handler. It is a dead button. Verified against `packages/ui/src/ErrorState.tsx`. Also stated wrongly to the operator in conversation and corrected there. |
| `delivery` F14: 10 OpenAPI operations unreachable from the SDK | **Overcounted.** 4 are reachable (`market.ts:179,182,185,190`); 3 are OAuth browser redirects where an SDK method would be a defect. **Real gap: 3.** |
| `layering` L-03: `internal/composition` unreachable by every instrument | **Overstated.** `scripts/arch-gate.sh:30` scans it. Verdict survives only under *inert = absent*; remediation is cheaper than claimed. |
| `cicd` F2: the remote is gone | **Wrong**, and it reported plausibly. It is an auth split-brain: `gh` asks as `leandrotcawork`, git authenticates as `developmentconexus-ops`. `origin/main` is `7df7d011` and fast-forward is possible. Both synthesis arms refuted it independently. |
| PHASE-0: three causes of `arch-gate.sh` failure | **Four.** `archscan` also fails at HEAD with 44 findings across `contexts` / `adapters` / `composition`. Additive, not contradictory. |
| `duplication` DUP-6: replace the `rows.Next()` loops | **Rejected and promoted to the do-not-touch list.** See above. |
