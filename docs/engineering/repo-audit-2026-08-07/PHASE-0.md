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

**Budget: zero-spend tooling only.** No paid CI minutes, no paid scanners, no paid required-checks.
This rules out platform-enforced branch protection if the repo is private, and pushes separation of
powers down to the mixed-change ban and shrink-only ratchets (playbook Phase 3, levels 4 and 5).
Every tooling recommendation must be free, or be labelled as costing money and priced.

**Team: solo operator plus agents.** There is no second human reviewer. Consequences the lanes must
reason with, not around:
- Human PR review is **not** an available control. Anything whose acceptance is "a reviewer will
  catch it" is a level-5 discipline rule and must be recorded as such.
- Credential custody **is** available: the operator can be the only party holding a merge credential.
- The author of nearly all code is a machine. Its failure mode is fluent, internally consistent
  wrongness — the wrong noun used correctly everywhere, including inside the guard written to catch
  it. Post-hoc review comments do not correct that; single-source compiled vocabularies and an
  in-loop verify manifest do.

**Repo visibility: `unverified`.** Two remotes exist — `origin`
(`developmentconexus-ops/marketplace-central`) and `legacy`
(`leandrotcawork/marketplace-central`). `gh repo view` cannot resolve `origin` from this machine's
credentials:

```
$ gh repo view --json visibility,nameWithOwner,defaultBranchRef
GraphQL: Could not resolve to a Repository with the name 'developmentconexus-ops/marketplace-central'. (repository)
```

Do not assert public or private. Where a recommendation depends on it, state both branches.

**`origin` is stale.** `main` is 96 commits ahead of `origin/main` and unpushed. Any claim about
"what CI runs" must distinguish what is *configured in the repo* from what has ever *executed on the
remote*, because the remote has not seen this work.

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
