# Phase 0 — What "good" means, in the operator's words

Method: `C:\Users\leandro.theodoro\Documents\MetalDocs\docs\engineering\repo-audit-playbook.md`.
That file is the single source of the method and is **not** vendored into this repository — a
second copy would be exactly the hand-sync defect the audit exists to find. Cite it by path.

Date: 2026-08-07. Repo state at authoring: `main` @ `1473e863`, working tree clean.

---

## The goal, in the operator's words

> I have run something really deep into MetalDocs and I am changing the way I code there to move to
> something more professional towards issues, PRs, PR review, CodeRabbit mechanical full validation
> and so much more. For that I had to identify every error in my code, my platform, to improve it and
> to create this full validation. I want to run it here as well so we move on the same path, this way
> it gets so much harder to send bad PRs.

Calibration, chosen explicitly by the operator: **solid professional level.** Not Google-tier.
That means: clear dependencies, clear modules, clear consumable surfaces, no hand-maintained
redundancy, following rules that have existed for decades, optimised for maintainability and future
scaling.

The operative sentence is the last one — **"this way it gets so much harder to send bad PRs."** The
program's success condition is a *mechanism* that makes a bad change hard to land, not a cleaner
codebase. A remediation that fixes code and leaves the judging to discipline has failed this goal
even if every finding is closed.

Paste this section verbatim into every brief. Do not paraphrase it.

---

## Hard constraints

**Budget: GitHub Actions is AVAILABLE.** *Amended 2026-08-07 by the operator, superseding the
zero-spend constraint this section originally carried.* The billing failure recorded below
(`gh run view 29712873135`, the only remote execution on record, which died before any build step)
is **resolved**. CI minutes are spendable.

Consequences, binding on Phase 2 and Phase 3:
- **Firing level 3 — "red build" — is reachable for every rule.** Recommendations must climb to it
  wherever the rule can be expressed as a command that exits non-zero. Proposing level 5 (discipline)
  for something a workflow could block is now a defect in the recommendation, not a constraint.
- **A required check is a verifier the author does not control at merge time** — the separation of
  powers the playbook's Phase 3 asks for. But the author here is an agent with write access to
  `.github/`, so the workflow file itself must be protected (CODEOWNERS on `.github/`, or a
  mixed-change ban that refuses a PR touching both a gate and the code it gates).
- **Linux runners dissolve the CRLF gate defect (D-52).** `ubuntu-latest` checks out LF; the 635
  pure-CRLF gofmt failures cannot occur there. Only the 22 genuinely unformatted files remain. The
  `.gitattributes`-renormalize-versus-selective-gofmt decision shrinks to a local-developer-experience
  question, not a gate question.
- **What CI still cannot run:** anything requiring the Oracle/Sankhya ERP or the local Docker dev
  stack. Those lanes stay local and are therefore still level 5. Say so explicitly rather than
  proposing a workflow that would need credentials this program will not put in a runner.

**Money still needs pricing, and cost is not the same as availability.** Label anything that consumes
minutes at a scale worth noticing (matrix builds, `windows-latest`, full integration runs on every
push). "Available" does not mean "free per run".

**The binding constraint is now delivery, not budget.** See "`origin` is stale" below: a required
check cannot gate work that is never pushed. Actions gates pull requests; today there are none, and
`main` is 96 commits ahead of a remote `gh` cannot even resolve. Any Phase 3 topology must state how
work reaches the remote before it states what blocks it there.

**Team: solo operator plus agents.** There is no second human reviewer. Consequences the lanes must
reason with, not around:
- Human PR review is **not** an available control. Anything whose acceptance is "a reviewer will
  catch it" is a level-5 discipline rule and must be recorded as such.
- Credential custody **is** available: the operator can be the only party holding a merge credential.
- The author of nearly all code is a machine. Its failure mode is fluent, internally consistent
  wrongness — the wrong noun used correctly everywhere, including inside the guard written to catch
  it. Post-hoc review comments do not correct that; single-source compiled vocabularies and an
  in-loop verify manifest do.

**Remote topology — corrected 2026-08-07.** The original text here claimed `origin` was unresolvable
and therefore of unknown visibility. That was an **authentication artifact, not a fact about the
remote.** Two GitHub accounts sit in the keyring:

```
$ gh auth status
  ✓ Logged in to github.com account leandrotcawork          (Active account: true)
  ✓ Logged in to github.com account developmentconexus-ops  (Active account: false)
```

`gh` asks as `leandrotcawork`; git authenticates as `developmentconexus-ops`. So
`gh repo view developmentconexus-ops/marketplace-central` 404s while `git ls-remote origin`
succeeds — and symmetrically, `gh repo view leandrotcawork/marketplace-central` succeeds while
`git ls-remote legacy` returns "Repository not found". Any future probe of either remote must state
which account it asked as.

**`origin` is real, reachable, and not diverged:**

```
$ git rev-list --left-right --count origin/main...main   → 0   98
$ git log -1 --format='%h %ad %s' --date=short origin/main
  7df7d011 2026-08-04 docs(onda-0): corrige F-15 — "nunca" da market_queue é projeto, não avaria
$ git merge-base --is-ancestor origin/main main          → fast-forward possible
```

`main` is 98 commits ahead and unpushed, but nothing has diverged. Claims about "what CI runs" must
still distinguish what is *configured in the repo* from what has ever *executed on the remote* — the
remote has not seen this work — but the remote itself is not the obstacle.

**`legacy` is superseded.** `leandrotcawork/marketplace-central`: PRIVATE, default branch `master`,
last push 2026-07-20, unreachable under the credentials git currently presents.

**`developmentconexus-ops` is a personal account, not an organization** (`gh api
orgs/developmentconexus-ops` → 404; it appears in `gh auth status` as an account). This is
load-bearing for Phase 3: **branch protection and rulesets on private repositories are a paid
feature** — GitHub Pro for a personal account, Team for an organization — while free accounts get
them on public repositories only. Repository visibility for `origin` is still **`unverified`**,
because confirming it requires asking as `developmentconexus-ops`.

Consequence to carry into the gate topology: if this repository stays private on a free personal
account, GitHub Actions will happily *run* every check and **cannot *require* any of them**. The
"required check" separation of powers then degrades to credential custody — the operator holding the
only merge credential — which is a real control but a level-5 one.

**Operator decision, 2026-08-07: the repository will be made PUBLIC on the free personal account.**
This resolves the branch above. Rulesets and required status checks are available at no cost, so
Phase 3 may assume **enforced** required checks — a verifier the author does not control at merge
time — rather than degrading to credential custody. Design the topology for that.

Two conditions attach, and both are hard:

1. **The flip is gated on a pre-public sweep.** Going public exposes the entire commit history, not
   the current tree; a credential committed once and deleted later is public forever. Real customer
   data is a live concern — the dev stack ran against ~38 real Mercado Livre orders and ~34 real
   listings, and Brazilian CPF/CNPJ in a committed fixture would be an LGPD exposure that no later
   commit can undo. `docs/engineering/repo-audit-2026-08-07/PRE-PUBLIC-SWEEP.md` carries the verdict.
   **A secret found in history is remediated by rotation, not by rewrite** — forks and caches outlive
   a `filter-repo`.
2. **Public changes what a gate must do.** Workflows on a public repository run for fork pull
   requests, log output is world-readable, and `pull_request_target` plus a writable token is a
   well-known escalation path. Any workflow this program proposes must state its trigger, its
   `permissions:` block, and whether it can be reached by an untrusted fork.

**Push status: the 98 local commits stay local for now.** The operator's sequencing is gates first,
then push, so that the first execution on the remote is already gated rather than a red run against
98 commits of history whose failures this audit has already catalogued.

**Network reachability of the unauthenticated API: MEASURED 2026-08-07, and the answer is YES.**
This resolves what earlier drafts left open, and it re-ranks the security axis. The API has **two
independent public paths**, both by design:

*Production* — `deploy/Caddyfile` terminates TLS on `{$MPC_DOMAIN}` and proxies to `backend:8080`:

```
@orders_api { path /orders /orders/*
              not header Accept *text/html* }
reverse_proxy @orders_api backend:8080
```

A request carrying `Accept: application/json` reaches `orders/transport/http_handler.go:608-618`,
which serialises buyer full name, CPF/CNPJ and address. Nothing in the chain checks identity —
`composition/root.go:994` composes the handler as `CORSMiddleware(apierror.Recover(mux))`.

*Development* — `docker-compose.yml` profile `oauth` runs
`ngrok http --url="$callback_host" frontend:5174`, publishing the frontend to the internet so
Mercado Livre's registered OAuth callback can reach it. `apps/web/vite.config.ts` proxies `/orders`,
`/pricing`, `/marketplaces/*`, `/profitability` and the rest to `backend:8080`. While that profile is
up, the same unauthenticated PII surface is live on a public URL.

Whether a production host is currently running is still **`unverified`** (`MPC_DOMAIN` is external),
but the dev path is one `--profile oauth` away and exists because the integration requires it. Do not
size the security axis as "localhost only".

*Related, same defect class:* the Caddyfile's own comment says it **mirrors the dev proxy table in
`apps/web/vite.config.ts` — keep both in sync when a new backend route prefix appears.** Two
hand-maintained route tables that must agree, with the Go router as a third source of truth. A route
added to the router and missed in either table fails in exactly one environment.

**What cannot break: `unverified`, operator has not enumerated it.** Treat the following as
candidates for the do-not-touch list and report evidence rather than assuming: the Docker dev stack,
the Sankhya/Oracle read path, the Mercado Livre integration. Lanes fill this in via their
"What is actually fine" section; the operator ratifies the consolidated list in Phase 2.

**Oracle/Sankhya is READ-ONLY.** No part of this audit or any remediation it produces writes to the
ERP. This is not negotiable and not subject to a cost/benefit argument.

---

## Repo shape (measured 2026-08-07 @ `1473e863`)

```
$ git ls-files '*.go' | grep -v '_test.go$' | xargs wc -l | tail -1     → 24039
$ git ls-files '*_test.go'                  | xargs wc -l | tail -1     → 23620
$ git ls-files '*.ts' '*.tsx'               | xargs wc -l | tail -1     → 35423
$ git ls-files '*.sql'                      | xargs wc -l | tail -1     →  3128
$ git ls-files '*.ps1' '*.psm1' '*.sh'      | xargs wc -l | tail -1     →  5805
$ ls apps/server_core/internal/modules | wc -l                          →    21
$ git ls-files | wc -l                                                  →  3121
```

~92k LOC tracked. Go backend in `apps/server_core`, React/TS frontend in `apps/web`, a hand-written
TypeScript SDK, PowerShell verification harness in `scripts/`, Postgres migrations in
`apps/server_core/migrations/`, Oracle/Sankhya as a read-only upstream mirrored into Postgres.

---

## Established facts — do NOT spend lane budget re-deriving these

Each is measured and reproducible. A lane that rediscovers one of these has wasted its budget; a lane
that *contradicts* one with better evidence should say so loudly, because that is a finding.

**Verification lanes and gates**

1. `scripts/arch-gate.sh` does **not** exit early when a step fails. `set -euo pipefail` is on, but
   every fallible command sits inside an `if` / `if !`, and under `-e` a status tested by `if` never
   triggers it. All five steps always run; `fail` aggregates and the exit decision happens at the
   bottom. (Debt D-51. A prior brief asserted the opposite and was corrected by measurement.)
2. `arch-gate.sh` is **FAIL** at HEAD, from three causes:
   - gofmt reports 639 files. **635 are pure CRLF artifacts** — `core.autocrlf=true`, git blobs are
     LF, `.gitattributes` pins only `*.sh`. This step can never pass on a Windows checkout. (D-52)
   - **22 files are genuinely unformatted**, independent of line endings.
   - `TestModuleBoundaryADR023` reports **234 violations**, all originating in and targeting
     `apps/server_core/internal/modules/`.
3. `npm run harness:unit` runs `go test ./tests/unit/... -count=1` plus the FE vitest suite. It does
   **not** run `./internal/...`. Those run only in `arch-gate.sh` step 4, which is not wired to any
   npm script (`grep -n "arch-gate" package.json` → 0 hits). (D-51)
4. The integration lane provisions an ephemeral `mpc_test_<hex>` database per run, so it never
   accumulates state between runs.

**Contract seam**

5. The TypeScript SDK is 100% hand-written: 2.595 lines, 172 interfaces. The same shape is
   transcribed four times — domain → DTO → OpenAPI → SDK.
6. Governance rule `GOV_API_SDK_SPLIT` requires only that those change in the **same commit**. It
   never checks that they **agree**.
7. `oapi-codegen` and `openapi-typescript` were **approved by the operator on 2026-08-07** as
   dependencies. They are not yet installed. A lane may cite them as available; none may assume they
   are in use.

**Architecture**

8. 32 ADRs are live. ADR-023 is the module protocol. ADR-033 (2026-08-07) amends `ARCHITECTURE.md`
   frozen decision 7: marketplace integrations enter through `adapters/marketplace/<vendor>`
   implementing ports owned by the consuming context, not through the `connectors` module. ADR-034
   supersedes ADR-017 — `kernel/fact.Fact` replaces the known-zero convention.
9. ADR-023's prose at `docs/architecture/decisions/023-module-protocol.md:82,300,306` still asserts
   "35 violations" three measurements out of date. The header is corrected to 234. (D-55)
10. `apps/server_core/internal/modules/` (21 modules) is the **legacy** tree being decomposed into
    `apps/server_core/internal/contexts/`. Both exist simultaneously. This is a known migration, not
    a finding — but *how far it has got* and *what crosses between them* is squarely in scope.

**Known-orphaned code (already declared as debt — report scale, do not re-discover)**

11. `internal/kernel/fact.Map`, `Combine2`, and `internal/kernel/provenance.Derived` have zero
    non-test callers. (Foundation for a later slice.)
12. Postgres role `mpc_app` (migration `0098`) is `NOLOGIN`, and the application DSN still connects
    as the table owner, which bypasses RLS. Catalog RLS is therefore true of the test suite and false
    of the running application. (D-44)
13. The catalog read side — `contexts/catalog/module.go:46` `Reader()`, `port.Reader`,
    `SummaryReader`, `Summary.DescriptionState` — has zero non-test callers. (D-53)

**The debt ledger**

14. `.mnfs/HARNESS-DEBTS.md` holds D-1 … D-55, each with `file:line` and a measurement. Read it
    before reporting anything as new. A finding already in the ledger should be cited by its number
    and **re-measured**, not re-narrated — the ledger records what was true when written.

---

## Lane roster

Ten concurrent read-only discovery lanes, one dimension each:

| Lane | Question it answers |
|---|---|
| `duplication` | What logic is implemented more than once? |
| `layering` | What dependency crosses a boundary this repo's own rules forbid? |
| `cicd` | Which gates exist, which actually fire, and what blocks a merge? |
| `delivery` | How many dialects solve the same HTTP request-boundary concern? |
| `persistence` | How is data-access correctness maintained — by machinery or by hand? |
| `testing` | What does the suite actually prove, and when does it prove it? |
| `observability` | Can you tell what the system is doing in production? |
| `security` | Which controls are real, which are contingent on configuration? |
| `goidiom` | What would a competent Go reviewer flag? |
| `frontend` | The same questions, other side of the wire (React/TS + the SDK) |

Lane reports land in `docs/engineering/repo-audit-2026-08-07/lanes/<LANE>.md`.

---

## Open operator decisions (not blocking Phase 1)

- **D-2 / D-50** — is `removal_owner=HARNESS-D-N` ratified permanent practice for governance
  exceptions with no open milestone, or must the four existing uses get real owners? Four exceptions
  in `contracts/governance/invariants.json` rest on an unratified precedent.
- **The CRLF gate (D-52)** — `.gitattributes` renormalize (one line, one large mechanical commit) vs
  running gofmt only where blobs and worktree agree. The `cicd` lane should size both before this is
  decided.
