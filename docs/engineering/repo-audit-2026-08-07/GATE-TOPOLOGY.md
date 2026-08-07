# Phase 3 — The target gate topology (merged)

Inputs: `synthesis-arm-a.md` §4.0–4.7, `synthesis-arm-b.md` §4.0–4.6. Neither arm saw the other.
Companion to `RECONCILIATION.md`, which merged the *axes*. This file merges the *gates*.

Why this file is the deliverable: the operator's success condition is a mechanism — *"this way it
gets so much harder to send bad PRs."* The ten axes fix code. This fixes what judges code, and it
outlives every one of them.

Binding rules carried from `SYNTHESIS-BRIEF.md`:

- Firing hierarchy: **1** unrepresentable → **2** boot-fatal → **3** red build → **4** runtime
  assertion → **5** discipline. **Level 5 is not a control.**
- **A control that exists and does not fire is absent.** Size it as absent.
- Every guard ships with an input that makes it fail, **in the same change**.
- The author of nearly all code here is a machine whose failure mode is *fluent, internally
  consistent wrongness — the wrong noun used correctly everywhere, including inside the guard
  written to catch it.*

---

## 0. Where the two arms diverged, and how each resolves

Seven divergences. Convergence is not recorded here; it is folded into the register in §2.

| # | Arm A | Arm B | Resolution |
|---|---|---|---|
| **G-1** Workflow shape | One `verify.yml`, two jobs (`verify-fast`, `verify-full`) | Two workflows: `ci.yml` (product gate) + `gate-integrity.yml` (meta gate) | **Both, composed.** B's split is right for a reason neither arm stated: two *separately named* required checks mean deleting either one shows up as "expected, not received" on a named check. Fold them into one file and a single workflow deletion removes both required checks in one edit. A's fast/full job split is orthogonal and correct — it lives **inside** `ci.yml`. |
| **G-2** CODEOWNERS | "Annotation, not enforcement, on a personal account" | Buys nothing today, with the mechanism: no teams on a personal account, and GitHub forbids self-approval, so the requirement is either unmeetable or admin-bypassed | **B.** Same verdict, but B names *why*, which is what stops it being re-proposed in six months. It is a two-party control asked to work in a one-party system. Add the file for the day a second person exists; **do not count it in the topology.** |
| **G-3** Signed commits on gate paths | Not proposed | C2, optional defense in depth — an agent with write access to `.github/` cannot sign as the operator | **Deferred, not scheduled.** B's residual threat is real (weaken the gate in PR 1, land the code in PR 2), but R1's empty bypass list plus C1 already force that weakening to be a single-purpose PR that must itself pass the gate. Take C2 only if a signing key is already configured. **Recorded as open, §8.** |
| **G-4** Integration lane in CI | CI-able: `Invoke-Integration` (`scripts/harness.ps1:68-88`) needs only Docker + `pwsh`, provisions its own ephemeral Postgres, embeds its migrations, needs **no ERP** | Also budgets it in CI, but against a `services: postgres:` container | **A's mechanism.** B's `services:` variant creates a *second* provisioning path, so the CI invocation would stop being byte-identical to the local one — which contradicts G7 ("one command is the gate"). Use the lane's own provisioning. **Requires one measured run on a runner to confirm; it is the only unverified claim in this file.** |
| **G-5** Push versus flip order | Flip public (step 4), *then* fast-forward `main` (step 5) | Push at P4 while still private and unprotected, confirm green at P5, flip at P6 | **B, and it is not close.** B's order runs the first remote execution — the one both arms expect to surface environment defects — with **logs that are not world-readable**. A's order publishes the history first and then debugs CI in the open. No cost to B's order: pushing to a private repo discloses nothing, so it does not need the sweep. The sweep gates the **flip**, not the push. |
| **G-6** D-52 / CRLF | "This is where D-52 dissolves" | Re-sizes D-52 **downward from blocking to optional**; the 635 CRLF `gofmt` failures cannot occur on a Linux checkout, leaving 22 real ones | **Converged; B's sizing adopted.** `.gitattributes` renormalize drops off the critical path and becomes local developer experience. Still worth doing — a local instrument nobody can read is how the gate stopped being run — but the operator does not need to adjudicate it. |
| **G-7** What gets ratcheted | Three shrink-only ratchets: archscan 44, governance 58, ADR-023 234 | Same three, plus treats `gofmt` (22) and `tsc` (12) as pay-down-now rather than ratchet | **Both, non-conflicting.** Ratchet what cannot be paid down in one sitting; pay down what can. 22 `gofmt` files = one commit. 12 `tsc` errors = one change. 44/58/234 = ratchets. |

---

## 1. The landing path — P1 to P8

No gate design means anything until this is settled: **a required check is a verifier the author does
not control at merge time, and today the merge does not happen anywhere a verifier can stand.** Every
one of the 98 commits reached local `main` without passing through a place a check could run.

| # | Step | Who | Notes |
|---|---|---|---|
| **P1** | ~~Fix the account split~~ **PARTIALLY CLOSED 2026-08-07 — see the defect below** | operator | `gh` asked as `leandrotcawork` while git authenticated as `developmentconexus-ops` — which is why the `cicd` lane concluded the remote was gone, **and reported that plausibly**. The `legacy` remote is removed. Both keyring tokens carry `repo` + `workflow`, so pushing `.github/workflows/` is in scope. |

> **P1 note — gh account state is shared mutable state across concurrent sessions.**
>
> **Root-caused by the operator, 2026-08-07: a concurrent session switched the account back.
> `gh auth switch` persists correctly; there is no gh defect.** An earlier revision of this file
> recorded one, on three observations of the account reverting between shell invocations. That
> inference was wrong and is corrected here rather than deleted, because the observations were real
> and the operational consequence survives the corrected cause.
>
> **What survives:** `gh`'s active account is **global to the machine, not to a session.** Any
> concurrent session can change which account — and therefore which of the two same-named private
> repositories — a `gh` command acts on, between one invocation and the next. The failure is not an
> error; it is a *wrong answer delivered confidently*: `gh repo list` did not fail, it listed the
> other account's repositories. Same class as the `cicd` lane's F2 concluding the remote was gone.
>
> **Standing practice, and it is cheap:** switch, verify identity, and perform the mutating operation
> **in one shell invocation**, with a hard abort between the verification and the mutation. The
> eleven issues were filed that way. The abort is what kept the first failed attempt at **zero
> created** rather than eleven filed into the wrong repository. Always pin `-R <owner>/<repo>` on any
> `gh` command that writes.
>
> **Two private repositories share the name:** `developmentconexus-ops/marketplace-central` (the
> `origin`, the live one) and `leandrotcawork/marketplace-central`.
>
> **Residual, surfaces at P4:** `credential.helper = manager` is a store separate from gh's keyring,
> so `git push` may present a different identity than `gh` does. `gh auth setup-git` makes git
> delegate to gh and removes the possibility of disagreement — **available, not contraindicated.** An
> earlier revision advised against it on the strength of the retracted defect above; that advice is
> withdrawn.
| **P2** | Build the gates locally | agents | All of §2, on a branch off local `main`, against the 98 commits as they stand. Nothing pushed. |
| **P3** | Baseline locally and commit the baselines | agents | Without this the first remote run is red on 58 governance violations, 44 archscan findings, 12 `tsc` errors and 22 unformatted files — **and a permanently red required check is switched off within a week.** The baseline commit is what makes the first green run possible. It is not bookkeeping. |
| **P3.5** | **I1-edge — close the unauthenticated PII route** | agents, hours | Carried from `RECONCILIATION.md` Divergence 1, unchanged and binding. Must land before P6. |
| **P4** | Push once, while unprotected | **needs explicit operator permission** | Fast-forward `main` — 98 commits plus the gate commits — in one push. **This is the last ungated write to `main` in the program's history, by construction:** a ruleset blocking direct pushes to `main` would block this very push, so the ruleset comes after it. The pushed tree already contains the workflows, so the push itself triggers the first run. |
| **P5** | Confirm green on the remote | agents | A local green and a runner green are different claims. The runner has a clean checkout, LF line endings, no `.gomodcache`, no warm caches. **Budget for one or two environment defects here.** Still private, so these logs are not public. |
| **P6** | Flip to public | **operator** | Gated on `PRE-PUBLIC-SWEEP.md` (verdict: `SAFE AFTER LISTED REMEDIATIONS`) **and** on P3.5. Going public exposes the entire history, not the current tree. **Operator decision 2026-08-07: the B-1/B-2 credential is not rotated now, and does not need to be before P4.** The push is not the disclosure event — the flip is. B-1's working-tree copy was scrubbed at `1ec2d081`; B-2 remains in a historical blob that is already an ancestor of `origin/main`, so rotation stays the only remediation that reaches the pushed copy. **The open question that decides whether rotation is needed at all: is that PostgreSQL role valid anywhere other than the operator's local machine?** Local-only puts it with N-1/N-2 in the accepted column and the sweep over-rated it; valid on any reachable host and the flip hands out a working credential. Operator answers before P6, not before P4. **ANSWERED 2026-08-07: local-only. B-1/B-2 closed, no rotation, P6 unblocked** — see the closure box in `PRE-PUBLIC-SWEEP.md` §2. **Operator directed the flip on 2026-08-07.** |
| **P7** | Enable the ruleset | **operator** | Required checks (`verify-fast`, `verify-full`, `gate-integrity`), require a pull request, linear history, block force-push and deletion, **empty bypass list**. §3. **Measured 2026-08-07: none of this exists on a personal account with GitHub Free while the repo is private** — both the rulesets and branch-protection REST endpoints return `Upgrade to GitHub Pro or make this repository public to enable this feature` (403). P6 is therefore not merely *before* P7 in order; it is what makes P7 possible at all. See `GATE-DESIGN.md` §0. |
| **P8** | Change the landing path in doctrine | agents + operator | `docs/HARNESS-PROFILE.md:279` forbids pushing without explicit permission; `:77` defines a post-merge ladder over an already-integrated local `main`. Chips must land as `branch → push → PR → required checks → merge`, and the doctrine must **say so**, or the checks sit beside the road exactly as `release-images.yml` does today. **This is the largest non-code item in the program and no lane costed it, because no lane was pointed at the harness.** Minimum change: standing authorization to push *chip branches* — never `main`. |

P1–P5 are prerequisites to everything. **P8 must land before the second axis ships, or the gates
built in P2 go back to gating nothing.**

---

## 2. The gate register

Levels are climbed as high as each rule allows. Where a rule could be expressed as a command that
exits non-zero and is nonetheless proposed at level 5, that is a defect in the recommendation, not a
constraint — GitHub Actions is available.

### 2.1 Level 1 — unrepresentable

Strongest and cheapest to maintain, because **nothing has to fire**.

| # | Makes unrepresentable | Axis | Negative fixture |
|---|---|---|---|
| **L1-a** | Money constructed from, or extracted into, a binary float. `exact.Money` with an unexported field, no `FromFloat`, no `Float64()`, `MarshalJSON` emitting a string. The `strconv.ParseFloat` sites and the 373 `float64` occurrences stop compiling rather than being flagged. | M1 | A `_test.go` under a `//go:build fixture_fail` tag attempting `Money(1.5)` and `m.Float64()`; a meta-check asserts `go build` returns non-zero. |
| **L1-b** | A hand-edited SDK type diverging from OpenAPI. The hand copy stops existing rather than being compared. `openapi-typescript` output committed, marked `// Code generated`. | C1 | A committed hand-edit to a generated line; the regen-diff job must go red. |
| **L1-c** | A repository constructed without a request-derived tenant. `TenantID` value type with no usable zero value; the `pgdb/config.go:23-24` `tenant_default` fallback deleted. **Fails to compile at all 70 `root.go` sites.** | I1 | A fixture constructing a repository with `TenantID{}`; must not compile. |
| **L1-d** | A value outside a vocabulary. Postgres enum types replacing the 65 `CHECK (… IN (…))` constraints — which also kills the drift between copies. | M1 | A migration test inserting an invalid value; must error. |
| **L1-e** | An error code with no HTTP status mapping. Typed error registry with an exhaustive switch. | R1 | Add a code with no mapping; the build must fail. |
| **L1-f** | A backend route prefix present in the Go router but absent from `deploy/Caddyfile` or `apps/web/vite.config.ts`. Both tables emitted from one authority instead of hand-maintained under a comment that instructs a human to keep them in sync. | C1 | CI regenerates and runs `git diff --exit-code`. **Interim, before generation exists:** a `verify-fast` step parsing all three and failing on disagreement — comparison is the fallback, not the target. |

`tsconfig.base.json` `"strict": true` already sits at this level, and it is the reason the frontend
lane's F3 defect is *findable at all*. The checker just never runs. That is V2, not a level-1 gap.

### 2.2 Level 2 — boot-fatal

Refuse to start rather than run wrong. The repo already does this correctly for `MC_DATABASE_URL`,
`MPC_ENCRYPTION_KEY` and the Oracle loader — **those are the template, not aspiration.**

| # | Fires at | Blocks | Negative fixture |
|---|---|---|---|
| **L2-a** | `pgdb.LoadConfig()` | Boot with no explicit tenant source. Same function that already fails closed on `MC_DATABASE_URL`. Deletes the dead `pgdb.DefaultTenantID` helper (SEC-8). | `TestLoadConfigMissingTenant` asserting a non-nil error. **The test must fail if the fallback is restored** — that is what makes it a fixture rather than a mirror. |
| **L2-b** | boot | The connecting DB role being `BYPASSRLS`, superuser, or the table owner. Queries `rolbypassrls` / `rolsuper` for `current_user`; refuses to serve. | Boot against the owner DSN; must refuse. **This fixture is the whole point** — it is the assertion that would have caught D-44 the day RLS was written. |
| **L2-c** | boot | A route in the PII set composed without the identity middleware. `composition/root.go:994` currently reads `CORSMiddleware(apierror.Recover(mux))` and contains no identity check. | Compose a PII route without the middleware; boot must fail. **This is the structural form of P3.5** — the edge control buys hours; this makes the exposure impossible rather than merely closed once. |
| **L2-d** | boot | A registered route absent from the auth allowlist. *Applicable only after auth exists* (I1-structural). | Register a route with no allowlist entry; boot must fail. |
| **L2-e** | boot | Silent degradation of logging. `slog.SetDefault` with a JSON handler in `cmd/server`; fatal if the destination is unwritable. | Point at an unwritable destination; must exit rather than fall back to text. |
| **L2-f** | the action | Provider writes / ML catalog offers / assisted linkage without explicit enablement. **Already implemented and already correct** — `MPC_PROVIDER_WRITES_ENABLED`, `MPC_ML_CATALOG_OFFERS_ENABLED`, the Oracle vars. | Already exists. **Do not touch.** This is also the entire reason §5's CI gap is tolerable. |

### 2.3 Level 3 — red build (the primary enforcement layer)

`ci.yml`, `pull_request` on all branches plus `push` on `main`, `ubuntu-latest`. **The same command
runs locally.** A gate the operator cannot reproduce on their own machine is its own defect class.

**Job `verify-fast`** — every push to a PR branch, target **under 3 minutes**:

| Step | Closes | Notes |
|---|---|---|
| `gofmt -l` over the whole Go module, not just `internal` | V2 | On a Linux runner the 635 CRLF artifacts cannot occur; only the 22 real ones remain. Pay those down in one commit, then this blocks at zero. **D-52 dissolves here.** |
| `go vet ./...`, `go build ./...` | baseline | Clean today. Keep it that way. |
| **`golangci-lint run`** | V2 | §2.3a. Has never run here; baseline and ratchet shrink-only on the first pass. |
| **`prettier --check`** over TS/TSX/JSON/MD/YAML | V2 | §2.3a. Blocks at zero, no ratchet — one mechanical commit precedes it. |
| **`eslint`** (type-aware, flat config) | V2, F1 | §2.3a. Has never run here; baseline and ratchet. |
| **PR title is a Conventional Commit** | V1 | §2.3a. Linear history + squash means the title becomes the commit subject. Fixture: a PR titled `wip`. |
| `tsc --noEmit`, root and every workspace | V2, F1 | 12 errors today, 3 in production code. **Wire the gate first and watch it go red on the `MutationPreviewModal.tsx:210` missing-`onRetry` case, then fix all 12** — proving it red before proving it green is the only order that certifies anything. |
| archscan over **all four roots including `internal/composition`** | B1 | 44 findings today. `scripts/arch-gate.sh:30` already scans that root, so this is **wiring a script, not building a detector** (RECONCILIATION B-2). Shrink-only ratchet. |
| governance validate + drift, `-BaseSha` from `git merge-base origin/main HEAD` | V2, B1 | 58 violations at HEAD with a zero diff. Baselined day one, shrink-only. |
| ADR-023 module-boundary detector | B1 | 234 violations. Ratchet. **Pure static analysis over the source tree — needs nothing but Go and a checkout.** That it is invoked only by `bash scripts/arch-gate.sh` on a developer's machine is the *content* of the finding, not a constraint on it. |
| Count assertion on every step | V3 | Each step prints an attributable count; **a step reporting zero executed units fails.** This is the anti-vacuity rule, and it is what separates a green run from a skipped one. |
| `pull_request_target` grep over `.github/` | §3 | Fails on the string. |

**Job `verify-full`** — on PR and on merge to `main`, target **under 12 minutes**:

| Step | Closes | Notes |
|---|---|---|
| `go test ./...` — the whole module | V2 | Brings 369 currently-dark test files into the light. **Expect the first run to be red. That is the point.** |
| vitest across all workspaces, glob-based, no filename pins | V2 | Removes the `apps/web/vitest.config.ts:14` name-pin and the orphaned root config. |
| Pester over `scripts/tests/*.tests.ps1` | V2 | ~250 assertions, 6/11 red cold. `pwsh` is present on `ubuntu-latest`. |
| **Integration lane** | V2, V3 | `Invoke-Integration` (`scripts/harness.ps1:68-88`): Docker + `pwsh`, both on `ubuntu-latest`; provisions its own ephemeral Postgres, embeds its migrations, needs **no ERP and no dev stack**. ~259s locally — the dominant term. **Must emit RUN/PASS/SKIP/FAIL counts before it is trusted; a zero-ran run must exit non-zero.** Pulled-vs-green is byte-identical without this. |
| `git diff --exit-code` on regenerated SDK + server types | C1 | After C1 lands. |
| Route ⇄ spec ⇄ SDK three-way set equality | C1 | **Joined on `operationId`, not on path template.** Fixture: delete one SDK method — the test must name the missing operation; and a route registered but absent from the spec. |
| Tenant-predicate AST checker over the 246 raw query sites | I1 | Allow-lists the 2 genuinely global tables. Fixture: a query against a tenant-bearing table with no `tenant_id` predicate. |
| Typed-error checker over `*/transport/**` | R1 | Bans `strings.Contains(err.Error(), …)` and `strings.HasPrefix(err.Error(), …)`. Fixture: a handler string-matching an error. |
| **Test-file census** | V2 | Enumerate every `*_test.go` and `*.test.tsx`; compare against the union of what the gate actually executes. **Kills T-1, T-3, T-7, T-8, T-10, T-12 as a class** rather than one at a time. Fixture: a new test file in an unreached directory. |
| **Guard-fixture meta-check** | V3 | Walk the guard inventory; every guard must have a registered failing fixture **that is actually executed**. Fixture: register a guard with no fixture — the meta-check must fail. **This is the mechanism of axis V3 and the answer to "evidence certifies itself."** |
| Exception-liveness check | B1 | For each `temporary_exceptions` entry in `modules.json`, assert the violation it excuses still occurs. **3 of the 5 entries fail this today**, including `migration-prefix-0021-duplicate` (2026-07-11), which predates by three weeks the 2026-08-03 collision it does not cover. |
| `internal/contexts` governed | B1 | `Policy.psm1:302` walk root becomes the pair `(kind, id)` over both trees; implement `GOV_CONTEXT_UNREGISTERED`. **The fixture at `lanes/cicd.md` F9 already exists and currently asserts a code that does not exist** — implementing the code turns a dead fixture live. |
| Ratchet comparison | V1 | Measured counts versus committed baselines. Blocks on increase, annotates the current number. Fixture: a PR adding one violation of each ratcheted class. |

### 2.3a The lint standard — settled here, because it was not settled anywhere else

**Measured 2026-08-07: this repository has no linter.** No `.golangci.yml`. No ESLint, Biome,
Prettier or dprint config. No SQL or migration linter. No commitlint, no hooks — `core.hooksPath` is
unset and `.githooks/` does not exist. **Not one `lint` script in any `package.json` across the
workspace.** `gofmt` and `go vet` are formatting and a sliver of vetting; archscan and governance are
architecture. None of that is lint, and until this section existed the topology quietly implied
otherwise.

The choice principle: **every linter here earns its place by catching a failure mode this audit
actually found, or by compensating for a decision this audit actually made.** No linter is included
because it is popular.

**Go — `golangci-lint`, blocking, ratcheted.**

| Linter | Why this one |
|---|---|
| `errcheck` | Unchecked errors. The single largest silent-failure class in Go and nothing here looks for it. |
| `staticcheck` | Broad correctness. Supersedes most hand-rolled checks. |
| `unused` | The audit found dead code that keeps a retired concept alive — `pgdb.DefaultTenantID` (SEC-8). This finds the rest. |
| `ineffassign` | Assignments that never reach a read. |
| `rowserrcheck`, `sqlclosecheck` | **These two are the compensating control for the do-not-touch ruling on the 82 hand-written `rows.Next()` loops.** Keeping the loops keeps a compile-time type check and keeps the risk of an unchecked `rows.Err()` or an unclosed `rows`. Rejecting `RowToStructByName` without adding these would be trading a real risk for nothing. |
| `bodyclose` | Leaked HTTP response bodies. This codebase is mostly provider HTTP clients. |
| `errorlint` | `errors.Is` / `errors.As` instead of `==` and type assertions. **The mechanical half of issue #7** — #7 bans string-matching on `err.Error()`; this catches the comparison forms it does not reach. |
| `exhaustive` | Non-exhaustive switches over enumerated types. **Directly supports #6's typed error registry and #8's vocabularies** — an unmapped code should not be reachable by omission. |
| `noctx` | HTTP requests built without a context. |

**Explicitly excluded, as a decision rather than an oversight:** `lll`, `gocyclo`, `funlen`,
`gocognit`, `dupl` and every other metric linter. They generate volume rather than defects, and a
ratchet on a number that means nothing is worse than no ratchet — it teaches the team that ratchets
are noise. Also excluded: `depguard`, because archscan and governance already own import boundaries
and two instruments over one rule is how they come to disagree.

**TypeScript — ESLint flat config with `typescript-eslint`, type-aware, blocking, ratcheted.**

| Rule | Why this one |
|---|---|
| `@typescript-eslint/no-floating-promises` | An unawaited promise is a silent failure with no stack. Requires type information, which is why `tsc` alone does not catch it. |
| `@typescript-eslint/no-misused-promises` | An async function passed where a void callback is expected — including React handlers. |
| `@typescript-eslint/await-thenable` | Awaiting a non-promise: usually a refactor that lost a call. |
| `react-hooks/rules-of-hooks`, `react-hooks/exhaustive-deps` | Stale-closure bugs. Not findable by review at speed. |
| `@typescript-eslint/no-unused-vars` | Same reason as Go's `unused`. |

**Strictly correctness rules. No stylistic rules in ESLint** — formatting belongs to Prettier, and
keeping the two apart is what stops the config becoming a preference argument.

**Formatting — `gofmt` (Go) and Prettier (TS/TSX/JSON/MD/YAML). Blocking at zero, never ratcheted.**
Formatting admits no legitimate exception, so a ratchet would only record how long the drift has been
tolerated.

**PR title convention — one job, blocking.** The repository already writes Conventional Commits
consistently (`docs(audit):`, `fix(catalogingest):`); nothing enforces it. With linear history in the
ruleset, a squash merge makes **the PR title the commit subject**, so validating the title is
sufficient and validating every commit in the branch is not. Negative fixture: a PR titled `wip`.

**Deliberately not adopted, each with its reason recorded so it is a decision:**

- **No coverage threshold.** With 91% of the suite currently unreachable, a coverage number would be
  a fabricated signal — and thresholds are satisfied by tests that assert nothing, which is this
  program's named failure mode (V3). Revisit after #2 and #3, never before.
- **No Dependabot version PRs.** `AGENTS.md` requires operator ACK for any dependency change; a bot
  opening version bumps manufactures exactly the churn that rule exists to prevent. **Enable
  Dependabot *alerts*** — notification only, no PRs, no conflict with the ACK rule.
- **No SQL linter yet.** The migration rules that matter here (prefix collisions, tenant predicates)
  are already owned by governance and by #5's AST checker. A general SQL linter would overlap both.
- **`.coderabbit.yaml` is configuration, not a gate.** Path filters excluding generated output, so
  CodeRabbit does not review a file whose authority is a schema. Non-blocking, per §6.

**Rollout, and this is the part that decides whether it survives.** `golangci-lint` and ESLint have
never run here, so the first run produces hundreds of findings. **Baseline and ratchet shrink-only**,
exactly as archscan (44), governance (58) and ADR-023 (234) are handled. Paying them down as a
precondition would stall #2 behind unrelated work, and running them non-blocking would make them
annotation — which is level 5 with a config file on. Formatting is the single exception: `gofmt` and
Prettier go to zero immediately, because that is one mechanical commit each.

### 2.4 Level 4 — runtime assertion

| # | Fires | Blocks / contains | Negative fixture |
|---|---|---|---|
| **L4-a** | `/healthz` | Container promotion. The handler actually checks DB and Oracle-reader state; `docker-compose.yml:46` already consumes it. | Boot with a dead DB; assert non-200. |
| **L4-b** | All 6 ticker loops | Contains — a panicking job is logged, the loop survives, HTTP stays up. Wrap in the existing `sync.Scheduler.safeInvoke` pattern. | A job that panics; assert the scheduler survives **and the failure was recorded.** |
| **L4-c** | Every request | Annotates — correlation ID in middleware, in the emitted log line, and in the response header. Makes every *other* failure findable. | A request asserting the ID appears in both places. |
| **L4-d** | React render | Contains — `ErrorBoundary` at the route level and around each data panel. A render error breaks a panel, not the app. | A component that throws; assert siblings still render. |

---

## 3. Protecting the gate from its own author

Stated plainly: **the author of nearly all code here is an agent with write access to `.github/`. A
required check the author can edit is not a separation of powers.**

**R1 — Ruleset on `main` with an empty bypass list.** This is the whole control. GitHub applies
rulesets to admins *unless they are explicitly exempted*, so an empty bypass list means the operator
is bound by the same checks as anyone else, and the only route around a red check is to go and edit
the ruleset — a **deliberate, separate, logged act on a different surface**, rather than a merge
button clicked at speed. That is what separation of powers reduces to in a one-person system: not a
second reviewer, but a second *deliberate step*.
*Negative fixture:* attempt a direct push of a trivial commit to `main`; it must be rejected.

**R2 — A required check that never reports must block as "expected, not received."** Verify this
setting deliberately, once. It is what makes *deleting the workflow* a losing move rather than a
winning one. Without it, the cheapest way past every gate in this document is `rm .github/workflows/`.
*Negative fixture — this is the topology's own fixture, and it must actually be run:* open a PR that
deletes `ci.yml`, and confirm the merge is blocked.

**R3 — Mixed-change ban** (`gate-integrity.yml`, level 3, ~40 lines). Fails if one PR's changed-file
list touches **both** a gate path and any file under `apps/` or `packages/`. The gate paths, measured
against the tree rather than assumed:

```
.github/workflows/**
scripts/harness/Policy.psm1
scripts/arch-gate.sh
contracts/governance/**
.golangci.yml
eslint.config.js
.prettierrc*
<the ratchet baseline files, once §1 P3 creates them>
```

*Correction recorded:* an earlier draft of this list included `.githooks/**`. **That directory does
not exist and `core.hooksPath` is unset** — it was written from assumption rather than measurement,
which is the exact error class this audit exists to catch. Removed. If hooks are ever added they
belong on the list, and they are level 5 regardless, so they never substitute for a job here.

Also fails on any workflow introducing `pull_request_target`, or elevating
`permissions:` beyond `contents: read` without a declared allow-list entry.
Combined with R1, this means **weakening the gate is a separate, single-purpose, visible PR that
must itself pass the gate.**
*Negative fixtures, both required:* a PR touching one workflow file and one handler; and a workflow
file containing `pull_request_target`.

**Not controls, and recorded as such so they are not re-proposed:**

- **CODEOWNERS.** `developmentconexus-ops` is a personal account, so there are no teams — CODEOWNERS
  can name individual users only, and the only user is the owner. "Require review from Code Owners"
  then means an outside fork's PR needs the owner's review (already true), while the owner's own PRs
  — which includes everything the agents produce under the owner's credentials — get **nothing**,
  because GitHub does not permit self-approval, so the requirement is either unmeetable or
  admin-bypassed. Add the file for the day a second person appears. **Do not count it today.**
- **Any rule whose enforcement is "remember to."** Level 5. Not a control.

---

## 4. Workflow inventory — trigger, permissions, fork reachability

On a public repository these three columns are **part of the gate's correctness**, not deployment
detail. Anyone may fork and open a pull request; every run's logs are world-readable.

| Workflow | Trigger | `permissions:` | Fork-reachable? | Ruling |
|---|---|---|---|---|
| `ci.yml` — `verify-fast` + `verify-full` | `pull_request` (all branches) + `push` on `main` | `contents: read` **only** | **Yes, by design** | Safe. A `pull_request` run from a fork gets a read-only `GITHUB_TOKEN` and **no secrets**, and neither job needs any — the integration lane provisions its own Postgres in-runner. |
| `gate-integrity.yml` — R3 | `pull_request` | `contents: read`, `pull-requests: read` | Yes | Safe. Reads the changed-file list with a read token; no write path. |
| `release-images.yml` — existing | `push` on tags + `workflow_dispatch` | `contents: read`, `packages: write` | **No, and it must stay that way** | Correct as written. Tag pushes and manual dispatch are not fork-reachable. Registry credentials live here and nowhere else. Its write scope is exactly why it must never gain a `pull_request` trigger. |

Three rules bind every workflow this program adds:

1. **`pull_request_target` is banned outright.** It executes in the base repository's context with a
   writable token and access to secrets, while checking out code the fork author controls — the
   textbook privilege-escalation path on a public repo. **No gate in this document needs it.** A ban
   is cheaper than a review policy and it is enforceable: it is in R3's forbidden-pattern list, so a
   workflow introducing it fails its own gate.
2. **`permissions:` declared explicitly at the top of every workflow**, never inherited. GitHub's
   default for new public repos is already read-only; stating it in the file makes the intent survive
   a settings change nobody remembers making.
3. **Never `echo` an environment.** Logs are world-readable. Same rule this audit itself ran under,
   now with a public consequence: no bare `printenv`, no `docker inspect`, no `set -x` over a block
   that touches configuration. Governance and archscan output is fine — it prints paths and code that
   are public anyway once the flip lands.

**Fork-PR approval must be set to "Require approval for all external contributors"**, not the default
first-time-contributor setting. One click per genuine outside PR — of which this repository will have
approximately zero — and it removes the whole class of "a fork opens a PR to burn minutes or probe
the workflow surface."

---

## 5. What CI cannot run, and why that is acceptable

Stated explicitly rather than proposed and quietly dropped:

- **The Oracle/Sankhya live lane** — needs ERP credentials, `cgo`, an Oracle client. Stays local. **Level 5.**
- **The Docker dev-stack live-drive / browser lane** — needs the composed stack and real provider
  credentials. Stays local. **Level 5.**
- **`harness:provider-write`** — performs live ML writes. Never in CI. Stays local. **Level 5.**

Level 5 is not a control, so the honest accounting is: **those three surfaces are not gated by CI.**
They are gated by **L2-f**, the fail-closed environment variables that already exist and already
work. That is a level-2 boot-fatal control, and it is why this gap is tolerable rather than alarming.

**Residual risk, unmitigated and named:** regression in ERP *read* mapping. The mitigation available
in CI is **golden fixtures captured from live responses**, not a live connection.

**On a public repository this stops being a limitation and becomes a requirement.** The
recommendation is stronger than "do not run these lanes in CI":

> **Do not add ERP, Sankhya, Oracle, Mercado Livre or any provider credential as a repository secret
> at all.**

Reasons in order of weight:

1. Every secret in a public repository is one misconfigured trigger away from disclosure, and the
   misconfiguration in question — `pull_request_target` — is the single most common workflow mistake
   in the ecosystem. R3 bans it, but **a control that must hold forever against an agent that edits
   `.github/` is a worse bet than not having the secret.**
2. Logs are world-readable. GitHub masks registered secret *values*; it does not mask a connection
   string assembled at runtime, a stack trace containing a DSN, or an ERP row echoed by a failing
   assertion. The audit's own PII constraint applies to the runner.
3. **There is no upside.** The only lanes that would consume them are the three declared level 5 above.

Useful consequence for P6: because the gate topology consumes **no secrets**, the required checks can
be built, pushed and proven green *before* the sweep returns its verdict. **The sweep gates the flip;
it does not gate the gates.** The one exception is `release-images.yml`, which already holds registry
credentials — whether those survive the flip is a question for the sweep.

---

## 6. Where CodeRabbit sits

The operator's goal names CodeRabbit explicitly, so, plainly: **it is an annotator, permanently.**

Its output is prose, and prose cannot be a required check without making every PR's fate a judgement
call — which is level 5 with a robot wearing it.

Its correct position is **upstream of the gate, as a rule generator**: every CodeRabbit finding the
operator accepts becomes a level-1/2/3 rule with a negative fixture, registered in the guard
inventory so the §2.3 meta-check executes it. That converts a review tool into gate material, which
is the only way it makes bad changes **hard to land** rather than merely **commented on**.

---

## 7. Cost

Cost is not availability, so: the whole non-integration gate is Linux-only and cheap.

- `ubuntu-latest`, 1× multiplier. Measured local times: `arch-gate.sh` ~78s, `harness:unit` ~82s,
  `go test ./internal/...` ~68s cached. A cold CI run without caches lands around **5–7 minutes** for
  build + vet + gofmt + go test + archscan + governance + tsc + vitest.
- Integration lane ~259s locally; run it on `pull_request` only, **not** on every branch push.
- **`windows-latest` is 2× and is not needed.** The PowerShell harness runs under `pwsh` on Linux and
  the CRLF problem does not exist on a Linux checkout. **No matrix.**
- One Go version, one Node version — the ones in `go.mod` and `package.json`.
- **~12 minutes per PR. On a public repository, GitHub-hosted standard runners are unmetered**, so
  the bill is zero regardless of volume.

Unmetered is not licence to be careless, for two reasons that survive the price being zero. **Wall
time is still paid by the operator** — a 25-minute gate is a gate people learn to merge around, and
12 minutes is already at the edge of tolerable for a solo workflow; cache Go modules and
`node_modules`, keep the integration lane off branch pushes. And on a public repo, **run volume is
partly controlled by strangers** — the fork-approval setting in §4 is what keeps that from becoming
someone else's decision.

---

## 8. Open, and not closable by this document

- **Integration lane on a runner.** §0 G-4. The claim that `Invoke-Integration` needs only Docker +
  `pwsh` is derived from `scripts/harness.ps1:68-88`, not measured on a runner. **One run confirms
  or refutes it.** It is the only unverified claim in this file.
- **Signed commits on gate paths (C2).** §0 G-3. Deferred. Revisit if a signing key gets configured;
  it is the only control here an agent cannot satisfy at all.
- **P4 requires explicit operator permission at the time.** Not pre-granted by this document.
- **P6–P7 are operator actions.** Visibility flip, ruleset configuration. (P1 closed 2026-08-07.)
- **Is the B-1/B-2 PostgreSQL role valid off the operator's machine?** Decides whether rotation is
  required at all. Gates P6. See the P6 row in §1.
- **L2-c and L2-d are deferred — authentication is not being built yet.** Operator decision
  2026-08-07: make the platform work first; there is one human user and no second principal.
  **Consequence, and it is the one to watch:** issue #1's edge control was written as a stopgap that
  L2-c would retire, and it is now the *only* control on the PII route, indefinitely — a Caddy config
  and an ngrok scope, level 3 at best, with no type and no boot condition behind it. That is fine
  while the repository is private. It sharpens at P6, which publishes `deploy/Caddyfile`,
  `apps/web/vite.config.ts` and `orders/transport/http_handler.go:608-618` — the route, its predicate,
  and the fields it returns — with a deny rule as the only thing between them. **#1's negative fixture
  is therefore not optional:** a request through the composed stack asserting non-200 without
  credentials, wired into `verify-full`. Without it, "the door is closed" is a claim about a config
  file nobody re-reads.
- **`.gitattributes` renormalize** — optional, local developer experience, no longer blocking. D-52
  is re-sized downward.

---

## 9. Construction order

The gates are built in the order that makes each next one cheaper, not in level order.

1. **R2 verified first** (`expected, not received`) — every later gate's strength is conditional on it.
2. `ci.yml` skeleton with the count assertion, wired to nothing. Proves the harness before the rules.
3. `verify-fast` steps, each with its fixture, each proven **red before green**.
4. Ratchet baselines committed (P3).
5. `verify-full`, integration lane last (it is the dominant term and the unverified one).
6. `gate-integrity.yml` (R3) — it must exist before the first PR that could touch a gate path.
7. R1 ruleset at P7, after the push.
8. The level-1 and level-2 gates arrive with their axes (M1, C1, I1, R1) — they are not a separate
   phase, and each one **removes** level-3 checks rather than adding to them. That is the direction
   the topology is supposed to move.
