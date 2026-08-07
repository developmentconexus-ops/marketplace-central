# Issue #2 — the verification gate, designed against the field

Research pass, 2026-08-07. Four independent streams: an adversarial adviser, a field survey of
eight rigorous open-source repositories plus six solo-maintainer ones, a tool-and-platform
landscape sweep, and a feature plan from the GPT-5.6 Sol planner reading this repository.
Everything below that is stated as a fact about this repo or about GitHub was **measured on
2026-08-07**, and the command is given so it can be re-measured.

`GATE-TOPOLOGY.md` is the prior design. This document does not replace it. It corrects six of
its claims, reorders it, and adds three controls it does not contain.

---

## 0. The finding that reorders the program

**On a personal account, GitHub Free, private repository, required status checks do not exist.**

Measured directly, not inferred:

```
gh api repos/developmentconexus-ops/marketplace-central/rulesets
{"message":"Upgrade to GitHub Pro or make this repository public to enable this feature.","status":"403"}

gh api repos/developmentconexus-ops/marketplace-central/branches/main/protection
{"message":"Upgrade to GitHub Pro or make this repository public to enable this feature.","status":"403"}
```

`GATE-TOPOLOGY.md` §3 states that R1 — a ruleset on `main` with an empty bypass list — "is the
whole control", and it is right: it is the only control in the program that the machine author
neither wrote nor can maintain. Every level-3 check inherits its force from it. **It cannot be
created on this account today.**

Without it, everything issue #2 builds is an annotation. A check nobody is obliged to pass is
level 5 with a YAML file attached, which is the exact thing this program exists to stop being.

Three ways out, and this is the one decision the operator has to make before the rest is worth
building:

| | Cost | Unlocks |
|---|---|---|
| **Flip the repo public** (P6, already on the landing path) | Free | Rulesets, branch protection, required checks, CODEOWNERS enforcement, CodeQL, **unlimited Actions minutes**, secret-scanning push protection, and CodeRabbit Pro-tier free forever |
| **GitHub Pro** | ~US$4/mo (price not re-verified) | Rulesets and branch protection on the private repo. Actions minutes stay capped at 3,000/mo |
| **Neither** | — | The gate is advisory indefinitely |

Going public is already the plan (P6) and is gated on the pre-public sweep and the credential
rotation question. It is also, by a distance, the best value: one bit flip turns on every
platform control at once and removes the Actions-minutes budget from the design entirely.

**Consequence for sequencing:** issue #2 can and should be built now — the checks run on PRs
regardless. But P6-or-Pro must land before issue #4 makes anything required, and until then
nobody should describe the gate as enforcing anything.

---

## 1. Corrections to GATE-TOPOLOGY.md, each measured

| § | Claim | Measured | Command |
|---|---|---|---|
| §4 | `release-images.yml` triggers on "push on tags + workflow_dispatch", not fork-reachable | **Also `push: branches: [main]`.** Still not fork-reachable, so not a hole — but §4 calls these columns "part of the gate's correctness", and this row was written from assumption | `sed -n '1,12p' .github/workflows/release-images.yml` |
| §7 | The Node version is already pinned in `package.json` | **No `engines`, no `packageManager`, no `.node-version`, no `.nvmrc`.** Node is unpinned everywhere. Nothing makes the runner's Node match the operator's | `node -e "…"`; `ls .node-version .nvmrc` |
| §2.3 | 22 `gofmt` files | **26**, against LF index content | `git ls-files '*.go'` → `git show :$f \| gofmt -l` |
| §2.3 | "635 CRLF artifacts" on a Windows checkout | **678 of 925 tracked files.** `core.autocrlf=true`, so the whole tree is CRLF and `gofmt` rejects all of it | `git ls-files -z '*.go' \| xargs -0 -n 200 gofmt -l` |
| §2.3 | "Pester over `scripts/tests/*.tests.ps1`" | Most of the 14 files are **directly executable scripts**, not Pester suites; one imports Pester. The lane needs two adapters, not one | `scripts/tests/live-oracle-docker-runner.tests.ps1:3` |
| §2.3 | 369 dark Go test files of ~404 | **Confirmed.** 404 tracked `*_test.go` across 968 package dirs — the 2,721 first counted were `.gomodcache`, a 245 MB in-repo dependency cache | `find … -not -path './.gomodcache/*'` |

### 1a. The CRLF correction is not cosmetic — it breaks the plan's own first rule

§2.3 opens with: *"The same command runs locally. A gate the operator cannot reproduce on their
own machine is its own defect class."* The very first row of its table then resolves the CRLF
problem by observing that "on a Linux runner the artifacts cannot occur" — which accepts that
local and CI diverge, on check number one.

Fix, one line, extending what issue #1 already started in `.gitattributes`:

```
*.go text eol=lf
```

The working tree becomes LF on Windows too, `gofmt` agrees in both places, and the 678 collapse
to the 26 that are real. This must land before the `gofmt` check, or the operator's first
local run of the gate reports 678 problems that do not exist.

### 1b. Wall-clock: the fast job does not fit as specified

Measured warm, on this machine:

| | Warm |
|---|---|
| `go build ./...` | **1m 59s** |
| `go vet ./...` | **44s** (exit 0, zero diagnostics — the plan's "clean today" is confirmed) |

§2.3 lists `go build`, `go vet` and `golangci-lint` as three separate steps of a job budgeted at
under three minutes. Those are **three full compilations of the same module**, and Go compilation
alone already consumes the budget before `tsc`, ESLint, Prettier, archscan, governance and the
ADR-023 walk have run. On a cold runner it is worse.

`golangci-lint` runs `govet` as one of its linters and type-checks the whole module to do it. It
subsumes both. **One Go compilation pass in the fast job, not three.** A fast job that overruns
its budget is not a slow gate — it is a gate the operator switches off, which is how this program
actually dies.

---

## 2. What the field actually does

Fourteen repositories read as configuration files, not as documentation: kubernetes, grafana,
coder, prometheus, caddy, next.js, astro, vite, tRPC, and the six that matter most here because
they are one or two maintainers — litestream, pgx, goose, type-fest, zod, sqlc.

### 2a. The irreducible core — present and blocking essentially everywhere

1. **A formatter check implemented as a diff, not a flag.** `gofmt -l -s -w . && git diff --exit-code`
   (pgx), `pnpm format && git diff --exit-code` (vite). Universal.
2. **Typecheck as its own named job**, separate from tests. Universal.
3. **Tests run against something real.** pgx runs PG 14–18 + CockroachDB + PgBouncer; coder boots
   real Postgres. This is where the rigour lives in the small repos, not in linters.
4. **`git diff --exit-code` after regenerating anything generated.** 5 of 8 large repos. This is
   the whole drift mechanism — there is no cleverer version of it.
5. **One aggregator job**, so branch protection needs a single required check and the fan-out
   lives in YAML. coder's is the template:

   ```yaml
   required:
     needs: [changes, fmt, lint, gen, test-go-pg, test-js, test-e2e, sqlc-vet, check-build]
     if: always()
     # We allow skipped jobs to pass, but not failed or cancelled jobs.
   ```

**Not in the core, and worth recording because they get proposed:** coverage thresholds (zero of
fourteen gate on coverage), CODEOWNERS (absent in all six small repos), merge queues (one repo,
and personal accounts cannot have them at all), conventional-commit title checks (two).

### 2b. The ratchet is a big-monorepo instrument

The six solo/small repos have **zero** committed baselines, **zero** ratchets, **zero**
`--new-from-rev`. The most instructive data point is negative: `pressly/goose` has the ratchet
flags physically present in `lint.yaml` and **commented out** —

```yaml
    # args: --issues-exit-code=0
    # Optional: show only new issues if it's a pull request.
    # only-new-issues: true
```

— a solo maintainer looked at the ratchet and chose whole-tree blocking with a small linter set
instead. Where ratchets do exist they are large-repo artifacts: grafana's `eslint-suppressions.json`
is 3,006 lines and, per their own comment, could not be extended to all their code.

This repo has six violation counts (58 governance, 44 archscan, 234 ADR-023, plus unmeasured
golangci and ESLint first runs), so some baselining is unavoidable. But the field says: **keep the
baselined set as small as you can, and put nothing in it that one commit could clear.**

### 2c. No AI reviewer is a required check anywhere

Searched for, and not found in a single repository. The reasons are structural:

- **Copilot cannot.** GitHub's own docs: *"Copilot always leaves a 'Comment' review… its reviews
  do not count toward required approvals and will not block merging changes."*
- **`claude-code-action` has no blocking mode.** Two of the small rigorous repos run it
  (`type-fest`, `zod`) — both comment-only, both gated to non-fork PRs.
- **Cursor Bugbot admits it:** *"requiring the status alone does not block merges on findings
  because findings default to `neutral`."*
- **Django bans it outright:** *"Do not request automated AI reviews… These reviews do not replace
  human review and often generate noise that distracts maintainers."*

There is a platform mechanism underneath this that is worth knowing, because it makes a gate look
configured when it is not. GitHub's docs: *"Successful check statuses are `success`, `skipped`,
and `neutral`."* **`neutral` passes.** A vendor can truthfully advertise "posts a check run" and
still be incapable of blocking. You would not discover this from the branch-protection UI.

### 2d. What the field does about machine-authored code

Prose, everywhere, enforced by nothing. A survey of 1,000 popular repositories found 118 with AI
policies and **no CI-enforced disclosure gate in the sample**. Linux merged an `Assisted-by:`
trailer convention; nothing validates it; Kubernetes bans that same trailer.

Exactly two mechanical artifacts exist in the whole field, and both are worth copying in spirit:

1. **coder/coder compiles review judgements into blocking lint rules.** `scripts/rules.go` holds
   ~25 project-specific `ruleguard` matchers, all blocking, each one a review comment that became
   executable — *"Forgot to return early after writing to the http response writer"*, *"Avoid
   calling pubsub.Publish() inside database transactions"*. One of them addresses AI agents by
   name and requires a written justification to bypass. **Provenance is not checked; the failure
   mode characteristic of agent output is.**
2. **tRPC drift-checks its agent-facing docs.** `validate-skills.yml` validates `SKILL.md` files
   on every PR that touches them, and a release-time job opens a PR when they go stale. The
   machine-readable instructions are treated as a versioned artifact.

coder's policy also states the reason this program exists, better than the audit did:

> *"Traditionally, maintainers assumed that producing a pull request required more effort than
> reviewing it. With AI-assisted tools, the balance has shifted: generating code is often faster
> than reviewing it."*

---

## 3. The three controls the plan does not contain

### 3a. revert-red — the anti-slop control

**This is the highest-leverage item in the whole research pass, and it is absent from
GATE-TOPOLOGY.md.**

The program's named failure mode is a test that passes without proving anything. Every mechanism
in the plan for catching it — the guard-fixture registry, the guard inventory walk, the
meta-check — asks the author to register something honestly. The author is the machine whose
failure mode is fluent, internally consistent wrongness *including inside the guard written to
catch it*.

revert-red asks nothing of the author. In CI:

1. Take the diff. Apply **only its test-file hunks** onto the merge base.
2. Run the touched packages/projects. **The new tests must fail.**
3. Apply the whole diff. They must pass.

A test that asserts nothing passes on the base tree too, so the gate goes red. A test whose
assertion is vacuous, whose fixture is symmetric, whose `expect` sits after an `await` on a target
that has already resolved — all of them pass step 2 and are caught. It needs no registry, no
fixture bookkeeping, and no honesty.

This mechanises the loop the operator already ratified — *all the REDs first, against a pristine
tree* — converting it from discipline (level 5, not a control) into a red build (level 3).

Cost: one extra scoped test run. Needs an escape hatch for pure refactors and test-only PRs, by
label, visible to the operator.

The literature name for what it detects is **pseudo-tested code**: a method whose entire body can
be deleted with no test failing. Vera-Pérez et al. found them in *every* project studied,
including well-tested ones. The packaged tool for it (Descartes) is Java-only; there is no Go or
TypeScript port. revert-red is the diff-scoped approximation, and it is cheap.

### 3b. Credential asymmetry — the one real separation of powers available

GitHub refuses any push that creates or modifies `.github/workflows/**` unless the pushing
credential carries the `workflow` scope. So:

- **A fine-grained PAT for agent use: `contents: write`, `workflows: none`, no admin.** Agents
  push chip branches with it. An agent PR touching a workflow file is rejected **by GitHub at
  push time** — not by a check the agent maintains.
- The operator's own credential is the only thing that can alter the gate.

This closes a self-reference the plan leaves open: R3, the mixed-change ban, lives in a workflow
file the agent can edit. The guard guards the guard. With the PAT, the platform protects the
workflows and the workflows protect everything else.

**Named weakness, honestly:** this runs on one machine, and both current keyring tokens carry
`repo` + `workflow` (`GATE-TOPOLOGY.md:46`). Credential asymmetry is only as strong as keeping the
operator's token out of the agent-reachable store. It converts "agent silently edits the gate"
into "agent must exfiltrate a credential" — a different and louder class, not an impossibility.

Unlike rulesets, **this works today, on Free, on a private repo.**

### 3c. Deploy provenance — two lines against a class already paid for twice

Nothing in the topology connects gated-`main` → built image → running container. This repo's own
history records the consequence: a green live drive against a stale binary.

Embed the git SHA at build time (`-ldflags`), expose it on `/healthz` (the handler already
exists — L4-a), and have any live-drive or QA step assert SHA-under-test equals SHA-claimed.

---

## 4. Ratchet by identity, never by count

Two independent streams — the adversarial adviser and the Sol planner — arrived at this
separately, which is why it is stated as a finding rather than a preference.

§2.3's ratchet row reads: *"Measured counts versus committed baselines. Blocks on increase."*
Two defects:

1. **A count ratchets the number, not the set.** A PR that adds three archscan violations while a
   refactor incidentally removes three others passes at 44 = 44. A machine author does this
   constantly — not maliciously, just by moving code.
2. **Shrink-only with a static baseline accumulates headroom.** The count drops 44 → 30, the
   baseline stays 44, and there are now fourteen free violations to spend.

Replace with **exact identities**: `rule + normalised path + stable symbol fingerprint`, sorted
and committed. New identity → fail, naming it. Disappeared identity → fail as *stale baseline,
lower it in this PR*. This matches the mechanism this repository already ratified for governance
exceptions, where an exception names exact paths and is explicitly not a permission
(`contracts/governance/README.md:31`).

This exposes a live contradiction in the plan worth fixing before it is met: §3's R3 puts the
baseline files on the gate-path list, and the mixed-change ban forbids touching a gate path and
`apps/` in one PR — but every violation-paying-down PR must lower the baseline in the same PR. As
written, R3 either blocks legitimate pay-down or forces the headroom defect. Resolution:
**shrinking the baseline set is allowed in a mixed PR; growing or loosening it requires a
single-purpose PR.**

Where an off-the-shelf tool already does this, use it rather than building: `golangci-lint`'s
`--new-from-rev` lints only code added since a revision and needs no baseline file at all.
Semgrep has `--baseline-commit`. The bespoke identity form is only needed for the three in-house
engines (governance, archscan, ADR-023).

---

## 5. Where AI review sits, and what it costs

Ratified: **AI reviewers annotate, they never block.** Not a compromise — the field does not
contain a counter-example, and there are three separate reasons.

- **The `neutral` loophole** makes a "required" AI check cosmetic without saying so.
- **Non-determinism and vendor flakiness.** `claude-code-action` issue #1299 is a permanently red
  required check caused by the action reacting to its own bot's commit. A gate the operator
  re-rolls until it goes green is not a gate; it trains them to ignore red.
- **Prompt injection, which is the sharp edge here specifically.** The diff is *input* to the
  reviewer, and in this repo the diff's author is an AI. This is documented at RCE severity, not
  theory: CodeRabbit (Aug 2025 — crafted PR → RCE, leaked their GitHub App private key, write
  access to ~1M repositories), Copilot (CVE-2025-59145), Qodo Merge, GitLab Duo. An Oct 2025
  cross-lab study found twelve published defences all bypassed >90% of the time.

**Pricing, for one developer on a private repository:**

| Tool | Private cost | Free tier on private | Blocks? |
|---|---|---|---|
| **CodeRabbit** | $24/mo annual, $30 monthly | Yes, rate-limited — PR summarisation, reviews only via IDE/CLI | Capability yes, default no |
| **Codex** (`chatgpt-codex-connector`) | **$20/mo — already paid for** | No | Comment-only |
| **Claude hosted Code Review** | n/a | — | **Team/Enterprise plans only. Not available on Pro or Max** |
| **Greptile** | $30/seat | Yes — 50 credits/mo, 1 dev | Leaning no |
| **Copilot review** | $10/mo | — | **Confirmed no** |
| **Gemini Code Assist** | **dead** — consumer tier shut down 2026-07-17 | — | — |

The practical reading: **Codex already reviews PRs on the private repos and is already paid for.**
It produced a high-quality review on `mnfs` PR #25 — severity badge, exact `file:line`, cited
`AGENTS.md:L34` as authority, named a concrete failure mode. That is the annotator seat, filled,
at no additional cost.

CodeRabbit's full experience — the one worth admiring on `mnfs` — exists there **because that repo
is public**. Free plan on a private repo gives PR summarisation only. So CodeRabbit is not a
purchase decision; it is another thing that arrives free the moment P6 lands.

**The one thing the human actually needs, and the plan does not provide:** a non-blocking
**decision digest** per PR — behavioural changes, goldens moved, counts moved, contract surface
touched. The operator's job in a one-person system is *deciding*, and a raw diff is not a decision
interface. This is what turns the single human check into something real rather than a merge-button
reflex. Either CodeRabbit fills it, or fifty lines of script summarising the gate's own outputs do
most of it.

---

## 6. Adoption sequence

The program dies exactly one way: **a required check that is red, or slow, often enough that
editing the ruleset "just this once" becomes routine.** The plan knows the failure mode and its
order still risks it — §2.3 says of `go test ./...`, *"Expect the first run to be red. That is the
point."* That is correct for **certifying** a check and fatal for a **required** one.

The rule that prevents it, absent from the plan: **no check becomes required until it has been
proven red once, then run green on several consecutive real PRs. Red-prove → green-prove →
require, per check.**

1. **Step 0 — retire the unknowns.** The plan's only self-declared unverified claim (§8: the
   integration lane on a runner, which is the dominant wall-clock term) is scheduled **last**.
   Unknowns go first. Half a day on a throwaway repo, alongside verifying "expected — waiting for
   status" behaviour. If the integration lane cannot run on `ubuntu-latest`, `verify-full` changes
   shape and everything downstream re-plans.
2. **`.gitattributes` `*.go text eol=lf`**, then pay the 26 `gofmt` files and the 12 `tsc` errors
   to zero. Both are one mechanical commit; neither is baselined.
3. **One command, byte-identical locally and in CI** — `npm run gate -- -Lane fast|full`, over the
   existing `scripts/harness.ps1` dispatch skeleton, which is already the right shape. Every lane
   prints `discovered/run/pass/skip/fail`; **`run=0` exits non-zero.** Today `harness.ps1:49`
   records only target, status and run id — a fully skipped run is byte-identical to a green one.
4. **First required set = the always-green trivia:** compile, `gofmt`, `tsc`. Under three minutes,
   green from day one. The operator's first month with a gate must be boring.
5. **Identity baselines** for governance (58), archscan (44), ADR-023 (234). Green by construction.
6. **`verify-full` non-required for a week** while the 369 dark test files are triaged. Genuinely
   broken ones go on an owned skip list; skipped is not green.
7. **Linters, after one dependency ACK.** `golangci-lint` and ESLint have never run here, so the
   first-run counts are unknown and must be **measured, not invented**.
8. **revert-red.**
9. **Then issue #4** — required checks, the ruleset, the mixed-change ban, the landing-path
   doctrine change. All of which need P6 or Pro to exist at all.

Enforce the budget as a check: `timeout-minutes: 15`. A gate that exceeds its own budget is a
defect in the gate.

---

## 7. Not doing, with reasons

- **No blocking AI reviewer.** §5. It would train the author to persuade rather than to be correct.
- **No coverage threshold.** No repository in the field survey gates on coverage, and an
  assertion-free test produces full coverage — this program's named failure mode, rewarded.
- **No full-corpus mutation testing.** Go's tooling is weak (`gremlins` has no score gate, no
  baseline, no diff mode) and runs are hours. revert-red covers most of the same ground at a
  fraction of the cost. Revisit only if revert-red proves insufficient, and only diff-scoped.
- **No assertion-density or vacuous-test linting for Go.** No such linter exists — searched.
  Go has no standard assertion library, so a syntactic checker cannot distinguish "asserts
  nothing" from "asserts via a helper". It would be satisfied by `assert.NotNil` spam anyway.
  (TypeScript is the opposite case: `@vitest/eslint-plugin` ships `expect-expect`,
  `no-standalone-expect`, `no-conditional-expect` and `prefer-strict-boolean-matchers`, which map
  one-to-one onto defects in this repo's own history. Cheap, take it.)
- **No Prettier in issue #2.** Broad mechanical churn, no measured correctness defect closed.
- **No CODEOWNERS, no merge queue.** CODEOWNERS is vacuous with one possible reviewer — GitHub
  never permits self-approval. Merge queues are org-only; a personal account cannot have one even
  when public.
- **No new bespoke checker beyond what exists.** §2.3 proposes eight custom level-3 instruments,
  each new surface for the exact failure mode they are built to catch. A bespoke checker earns its
  place only when no off-the-shelf tool and no level-1 restructuring covers it, and every one must
  print population-versus-extraction counts — the rule this repository already ratified as *a
  sweep is only as wide as its pattern*.
- **No provider, ERP, Oracle or dev-stack credential in GitHub, ever.** Unchanged from §5 of the
  topology, and it hardens on a public repo.

## 8. Named and unowned

- **The fiscal golden corpus.** The operator-verified, centavo-exact datasets — 784/784 against
  real notes, 6,258/6,274 on the single-formula reconciliation — are the most expensive artifacts
  this project has produced, and **no CI lane can currently reach them.** §5 of the topology names
  golden fixtures as the mitigation for the ERP-read residual risk and then no issue owns it.
  They are perishable. Capture, sanitise, commit, replay against the adapters in `verify-full`.
  Guard the "update snapshots" reflex: a golden *change* goes through the growing-versus-shrinking
  rule; new goldens are free.
- **Windows Actions-minutes deduction rate.** Current GitHub docs no longer state whether Windows
  minutes deduct from the included quota at 2× or 1:1. This harness is PowerShell. Read the rate
  off Settings → Billing rather than trusting either reading. Moot if P6 lands — public
  repositories have unlimited minutes.
- **Node is unpinned.** §1. Pin it before claiming local/runner equivalence.
