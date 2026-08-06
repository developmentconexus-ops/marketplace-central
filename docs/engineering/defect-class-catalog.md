# Defect-Class Catalog — marketplace-central

> **Purpose.** Not a bug list — a list of **defect *classes*** observed while building this
> system with AI agents, each paired with the mechanism that makes the class *unreachable
> by construction* in a future project.
> **Audience.** The software factory: whoever sets up repo scaffolding, CI, dispatch
> doctrine and review gates on day 0 of the next system.
> **Rule of the document.** Every class carries real evidence from this repo (commit,
> `file:line`, mission, debt ID). No hypothetical defects. A class with no evidence does
> not belong here yet.
> **Sibling document.** `MetalDocs/docs/engineering/defect-class-catalog.md`. It owns the
> **Prevention Ladder (§0)** and 25 classes on the *code and architecture* axis. This file
> does **not** restate the ladder — restating it would be that catalog's own Class 2
> (Hand-Synced Enumerations). Read it first; rung numbers below refer to it.
> **This file's axis is different and complementary: verification, evidence, and
> agent-authored work.** Where a class overlaps, the overlap is named explicitly.
> **Last verified:** 2026-08-06 — 27 classes. V10, V11, P6, P7 and F5 were added from the
> foundation-kernel work; root cause 8 with them.

---

## Why this axis exists

The sibling catalog answers *"what did we build wrong?"*. This one answers a question that
only appears once an agent is doing the building: **"why did we believe it was right?"**

Almost every class here is a **signal that was cheap to produce, stable, reproducible, and
wrong about which world it was in.** It passes in the fixed world and in the broken world,
so it feels like evidence and proves nothing. An agent optimizing toward a green signal
will find the cheapest green available, and the cheapest green is almost never the one that
required the work.

That is not a moral failing of the agent. It is an incentive-design failure of the harness.
The factory's job is to make the cheapest available green also the honest one.

---

# Part I — Verification and evidence

## Class V1 — The Vacuous Green

**Symptom.** A lane prints success that is **byte-identical** to real success, without
having executed the thing under test.

**Evidence.** Four independent instances, all paid for:
- `//go:build cgo` on a host without gcc → `no tests to run`, exit 0. The Oracle lane was
  "green" on the host for an unknown period while executing nothing
  (`oracle-live-lane-runs-in-docker`, D-8).
- `npx --no-install tsc` in a worktree without `node_modules` → pass, having type-checked
  nothing (`chip-m02-mirror-port-merged`).
- `apps/web/vitest.config.ts` includes tests by **exact filename**; a file outside the list
  is silently zero, not an error (HARNESS-DEBTS B-7).
- Integration lane discovery is blind to half the tree; `status=passed` on a skipped
  package is indistinguishable from a real run (B-1, B-6, `integration-lane-failure-token`).

**Root cause.** "Ran zero tests" and "ran tests, all passed" share an exit code and, in
most tooling, an output shape.

**Prevention (rung 3).** A lane must emit an **attributable count**, and the gate asserts
on the count, not the exit code. This repo's measured form: a `failure_token=test=<name>`
is attributable; `package=` is not. Corollary rule, mechanizable today: **a lane that
reports zero executed units fails.** Zero is never green.

**Detection in an existing repo.** For each lane, delete the code under test and re-run. If
it still prints green, the lane proves nothing about that code. Cheap, and it has never
once been a waste here.

**Overlaps** sibling §8 (Tests at the Wrong Altitude) and §16 (Compiles ≠ Works), but the
mechanism is different: those are tests at the wrong level, this is *no test at all*
wearing the same output.

---

## Class V2 — The Assertion That Cannot Fail

**Symptom.** A test exists, is well-named, is read by reviewers as coverage, and would pass
against the defect it names.

**Evidence.**
- `expect(screen.getByLabelText("Data freshness")).toBeInTheDocument()` — passes with any
  string, which is exactly why a formatter that could not distinguish 15 minutes from
  15 days survived the entire suite (`anti-patterns.md`, Test design).
- A criterion *"never calls X"* over a type that does not declare `X`. Unfalsifiable
  against the signature (`criterion-unfalsifiable-against-type`).
- A synchronous non-call assertion over a target reached after `await`: the assertion runs
  before the awaited code can call anything, and **passes with an injected click**
  (`stale-binary-makes-live-drive-lie`).
- A symmetric fixture: swapping the inputs still passes, so it never tested order
  (`deleted-test-restore-or-observable`).
- Round-trip JSONB compared by byte-exact string equality — the encoder's key order, not
  the value (`chip-m01-sync-state-ready`).

**Root cause.** Presence is cheaper to assert than value, and reads identically in a diff.

**Prevention (rung 4, enforced at rung 3).** **Red before green, with the failure text
recorded.** A new assertion is not accepted until the harness has seen it fail with a named
message. Mechanizable: require, per new test, an artifact containing the red output; CI
compares that the test name in the artifact matches the test added in the diff.

**Detection.** Inject the defect the test claims to catch. If the test still passes, it is
decoration. In this repo the injection round has *never* failed to find something.

---

## Class V3 — The Blind Instrument

**Symptom.** A measurement is run, returns a confident number, and answers a **different
question** than the one asked. The number is then quoted as fact.

**Evidence.** All from this repo, all caught before shipping only because a second
instrument disagreed:
- `pg_stat_user_tables.n_live_tup` reported `listings = 0`; real `count(*)` is **34**. It
  is an autovacuum estimate. An entire "42 empty tables" finding was discarded
  (`varredura-maximo-global-instrumentos`).
- OpenAPI `operationId` compared against frontend call sites, to find orphan operations.
  The SDK is hand-written and renames: `listMarketplaceOrders` in the spec, `listOrders` in
  the SDK. The whole orphan list was noise (same).
- The same sweep re-run over `apps/web/src` only — but `packages/feature-*` and
  `packages/web-query` also consume the SDK. 49 became **42** (same).
- `grep Degrau3` returned only test files, suggesting 102 dead lines. The file's exported
  symbols are `ProductIdentityReader` / `CategoryResolver` / `CommissionQuoter`, wired at
  `composition/root.go:966-968` (same).
- Joining `TGFDIN` without `CODINC` fanned out 1.391 of 6.274 items. The round reported
  "342 items do not close"; pre-aggregating correctly, **13**
  (`uma-formula-fiscal-cst-e-a-chave`).
- The module-dependency checker's regex is blind to an import of a module **root**
  (HARNESS-DEBTS D-27). `EXIT=$?` after a pipe records the exit of `sed`
  (`chip-vinc-neutro-round9`).

**Root cause.** The first instrument that returns is accepted as authoritative, and its
*universe* is never stated. A count without a stated universe is not a measurement.

**Prevention (rung 3 + process).**
1. **Every reported number ships the exact command that produced it**, as an artifact — not
   as prose. Unreproducible numbers are opinions.
2. **A count that carries a verdict requires a second, independent instrument.** Not a
   re-run of the same one.
3. Ban estimator columns from verdicts by name (`n_live_tup`, `reltuples`, sampling
   `EXPLAIN` rows). Mechanizable as a grep in the evidence linter.

**Detection.** Ask of every number: *what would this instrument report if the answer were
the opposite?* If it would report the same thing, it is blind.

---

## Class V4 — The Lying Environment

**Symptom.** The code ran. It was not **your** code.

**Evidence.**
- A live drive against a container built **before** the slice's commit. Absence of the
  defect and absence of the code look identical
  (`stale-binary-makes-live-drive-lie`).
- A governance gate resolving paths against a stale worktree, reporting on "the repo".
  Both halves of its accusation were false (HARNESS-DEBTS D-24, `gate-reads-diff-not-pack`).
- A worktree without `node_modules` resolves `@mc/*` to the **main branch's** packages, so
  the lane type-checked the wrong tree's types (`web-tsc-lane-cross-branch-resolution`).
- `.gomodcache` at the repo root silently ignored; the build died on an empty cache under
  `GOPROXY=off` and read as a migration defect (`gomodcache-root-gitignore-trap`,
  REVIEW-LEARNINGS 2026-07-15).
- Pre-provisioned worktrees are born stale (D-12). A worktree directory outlives the
  worktree and Docker's bind mount locks it (D-18).

**Root cause.** Every layer between the commit and the observation is assumed transparent.
None of them are.

**Prevention (rung 3).** A **preflight that prints identity and fails closed**: absolute
path measured, tip SHA, container image build time. It **fails** if the container's
`CreatedAt` predates the tip commit. This is fully mechanizable, cheap, and is the single
highest-return gate in this catalog — it converts an entire class from "found later, at the
cost of a wrong verdict" to "impossible".

**Factory rule.** **A verdict prints the tree it measured.** A report that names no tree is
not a report.

---

## Class V5 — The Gate Whose Tools Cannot Reach Its Criterion

**Symptom.** A review gate is dispatched with a tool set that physically cannot verify what
it is asked to certify — so it certifies the *claim*, and the claim came from the work
under review.

**Evidence.** A read-only reviewer (no shell) was asked to discharge criteria of the form
"the lane runs and passes". It could only read the chip's own assertion that it did. Every
execution criterion in that gate was, transitively, self-certified
(`readonly-gate-cannot-discharge-execution`, HARNESS-DEBTS C-4).

**Root cause.** The gate's *scope* was designed independently of the gate's *capability*.

**Prevention (rung 3).** Mechanical, checkable at dispatch time: **tools(reviewer) ⊇
tools required by the criteria assigned to it.** Criteria are tagged by the capability they
need (`reads`, `executes`, `queries-db`, `drives-browser`); a dispatch that assigns an
`executes` criterion to a reviewer without a shell is rejected before it runs.

**Generalization worth more than the class.** *Any* self-certified criterion is a criterion
without a gate. The pattern reappears whenever the producer of a claim is the only source
for its evidence.

---

## Class V6 — The Diff Against The Wrong Base

**Symptom.** A gate reviews a real diff, correctly, and misses a defect that exists only in
the *relationship* between the branch and its target.

**Evidence.** A chip's merge silently **reverted a feature from `main`**, because the diff
was computed against the branch's dispatch base rather than the target's tip. The gate saw
only the chip's own additions (`diff-vs-dispatch-base-hides-revert`). Separately, a gate
read the evidence *pack* rather than the diff, and so blocked on nothing measurable
(`gate-reads-diff-not-pack`).

**Root cause.** "The change" is ambiguous — it means one thing at dispatch and another at
merge, and the tooling used the same word for both.

**Prevention (rung 3).** **The merge gate's base is always the tip of the target**, never
the dispatch base. Enforced in the gate runner, not in doctrine. And: a gate blocks on an
**observable** — never on the absence of narrative in a pack.

---

## Class V7 — Zero As Proof

**Symptom.** A sweep returns "0 violations" and is read as "the property holds". It often
means the instrument found nothing it was capable of finding.

**Evidence.**
- A repeated sweep was offered as proof of ADR-C2 compliance. It could not be — there was
  nothing to delete, so the sweep would return 0 in both worlds. Only a **positive control**
  plus `fetched_at` advancing proved it (`p1-validated-positive-control`).
- An error-surface gate refuted its own `0`: the instrument was blind to the undeclared
  case, so the honest reading of 0 was "unmeasured" (`chip-error-unify-closed`).
- A must-fail whose injection was never confirmed is a **false empty**, not a pass
  (HARNESS-DEBTS D-8).

**Root cause.** Absence of a signal is symmetric between "clean" and "not looking".

**Prevention (rung 3).** **Every zero-result check ships a positive control** — a
deliberately planted instance the check must find in the same run. A check that cannot
demonstrate a catch does not report a count.

**Overlaps** sibling §3 (Allowlist Guard That Ratifies Drift) — an allowlist is one way to
manufacture a zero.

---

## Class V8 — Failure With No Path To A Pixel

**Symptom.** Something fails unattended, every layer handles it politely, and the operator's
screen stays green. The system is not lying about a test — it is lying about **production**.

**Evidence.** The most operationally expensive class in this repo.
- **Mercado Livre token-refresh failure is invisible at every layer**, and the cockpit
  reports health. The operator's account can be disconnected while the screen says
  connected (`mis008-design-closed`).
- A hand-rolled scheduler was specified beside `syncapp.Scheduler`, discarding
  `RecordFailure` — the one thing that feeds the operator's sync-health card. A broken
  collection would have rendered green (`anti-patterns.md`, Planning).
- A long ticker with no initial tick and no persisted due-time: silent starvation
  (HARNESS-DEBTS D-16). An installation connected after boot never gets a scheduler (D-20).

**Root cause.** Error handling is written to keep the process alive, and "alive" is then
rendered as "healthy". Nobody owns the path from a caught error to a rendered state.

**Prevention (rung 3 + planning gate).** For anything that can fail unattended, the plan
**names the path from failure to pixel** — which record captures it, which query reads it,
which component renders it. A slice whose failure mode is invisible is incomplete, not
"hardened later". Reuse of an existing seam (a scheduler with `RecordFailure`) buys this
path for free, which is usually the decisive argument for reuse over a parallel build.

**Detection.** Kill the dependency — revoke the token, stop the container, break the
credential — and look at the screen. If it is green, the class is present.

---

## Class V9 — The Third Patch

**Symptom.** The same defect shape is fixed a third time, in a third place, by a third
targeted patch. Each fix is correct. The class survives all of them.

**Evidence.** Ratified into the profile at `1889d0dd` after the second occurrence
(`stop-the-line-class-rule`, HARNESS-DEBTS C-1). Instances that forced it: the error-surface
divergence across modules, the freshness-formatter copies, the repeated
false-accusation hook (D-7, which reached a **third** occurrence during the very session
convened to fix it).

**Root cause.** A targeted patch is always cheaper than a class fix *at the moment of the
patch*, and the accumulated cost is paid by someone else later.

**Prevention (process, hard rule).** **The second occurrence of a shape stops the line.**
Root-cause it, then either fix the class or register the debt with the measurement that
proves its size. A third targeted patch of the same shape is a process defect, not a fix.

---

## Class V10 — The Census That Only Counts The Enrolled

**Symptom.** A checker iterates a **declaration** — a registry, a manifest, an evidence
pack, a route table — and reports on what it finds there. Anything that exists in the
territory but was never declared produces **no finding at all**. The failure mode is not a
wrong number; it is silence, which reads as compliance.

Distinct from Class V3: a blind instrument answers the wrong question loudly. This one
answers the right question over the wrong universe, quietly.

**Evidence.**
- `scripts/harness/Policy.psm1:308` (`Test-GovernanceDrift`) iterates
  `$moduleRegistry.modules` and requires each root to be exactly
  `apps/server_core/internal/modules/<id>`. A directory under `internal/contexts/` with no
  registry entry yields zero issues. The rule cannot see what nobody enrolled
  (`registry-walks-registry-not-tree`).
- The merge gate read the chip's evidence **pack** rather than the diff, so anything the
  chip did not narrate was outside the gate's universe (`gate-reads-diff-not-pack`).
- A module-scoped sweep validated inside the module while the caller — the site that
  decided the value — sat outside it (`chip-fim-closed`).
- `apps/web/src` treated as "the frontend"; `packages/feature-*` and `packages/web-query`
  also consume the SDK (`varredura-maximo-global-instrumentos`).

**Root cause.** Declaration-driven checks are cheap to write and their coverage is invisible
in their output. A passing run and an empty universe are byte-identical.

**Prevention (rung 2/3).**
1. **Every registry rule needs a paired tree rule.** For each declaration that is checked,
   walk the filesystem and emit a finding for every entity present but unenrolled
   (`GOV_CONTEXT_UNREGISTERED` is the concrete instance). The pair is the invariant, not
   either half.
2. **A checker must print its universe**, not just its verdict: *"scanned 21 registry
   entries, 22 directories, 1 unenrolled"*. A verdict without a denominator is not a
   measurement.
3. **Prove the tree rule with a decoy.** Create an unenrolled entity, require FAIL, remove
   it. A rule that has only ever returned zero has never been exercised.

**Detection.** Ask: *if I added a whole new subsystem and told the checker nothing, what
would it say?* If the answer is "nothing", the checker measures obedience, not reality.

---

## Class V11 — The Must-Fail Proof Against A Red Baseline

**Symptom.** A negative control — inject a violation, require the gate to reject it — is run
while the tree is already failing for an unrelated reason. It fails. The failure is recorded
as proof the guard works. It proves nothing: the observed FAIL is not **attributable** to
the injected violation.

**Evidence.** `main @532811b4` carried two independent reds: `go build ./...` rejected
`internal/composition/catalog_wiring.go:9` (`use of internal package .../catalog/internal/postgres
not allowed`), and `TestCanonicalSourceListsEveryMigrationByFullFilename` reported *"got 84
canonical migrations, want 83"*. Two planned proofs — an unregistered context directory
required to FAIL, and an edited OpenAPI required to diff — were scheduled against that
baseline. Neither would have been attributable. The build break additionally prevented the
full unit suite from running, so the migration red had been invisible since the task that
introduced it.

Same family as the attributable-token finding: `status=passed` does not prove execution,
`failure_token=test=` does (`integration-lane-failure-token`).

**Root cause.** A negative control is only valid against a known-green baseline, and
"baseline is green" is an assumption nobody is assigned to verify. It is also the assumption
most likely to be false precisely when work is in flight.

**Prevention (rung 2/3).**
1. **Green-before-red is a gate step, not etiquette.** Any must-fail proof runs
   `build + unit` first and records the exit code. A non-zero baseline aborts the proof.
2. **The FAIL must name the injected thing.** Not "the gate failed" — the gate's output must
   contain the decoy's identifier (`tmporphan`, `probeField`, `probe_float.go`). An
   unattributable FAIL is not evidence.
3. **Sequence work that shares a gate file.** Two plans editing `arch-gate.sh` in disjoint
   places still need an order, because the second one's proofs run against the first one's
   tree.

**Detection.** Ask of every must-fail proof: *would this have failed anyway?* If you cannot
answer from a recorded pre-injection run, the proof is unfinished.

---

# Part II — Planning and specification

## Class P1 — The Plan As Unmeasured Allegation

**Symptom.** A plan (or brief, or card) states facts about the repository that were
generated from pattern-matching rather than read from the tree. They are plausible, and
some are false.

**Evidence.**
- A query written against `sync_state.last_success_at`. The table has `last_full_sync_at`
  and `last_incremental_at` (`anti-patterns.md`, Planning).
- **Three** false claims in a single milestone brief, false by the time a worker read them
  (`wave2-m06-dispatched`, HARNESS-DEBTS A-2).
- A card shipping a lane command that produced `No test files found` (A-2, case S8).
- A card that inherited a falsehood from a prior ruling (A-2, case S10).

**Cost.** Each instance is a full wasted round or an escalation. This is the highest-volume
class in the repo's history.

**Prevention (rung 3) — the highest-return CLI the factory can build.** A **claim linter**
that runs on any plan/card/brief before dispatch. It extracts and verifies, mechanically:

| extracted | verified against |
|---|---|
| every fenced shell command | executed dry, or at minimum resolved (binary exists, path exists) |
| every `file:line` | file exists, line exists, **and the tree it was measured in is named** |
| every bare symbol name | exists in the tree (`go doc`, `tsc`, or grep) |
| every table/column name | live schema (`information_schema`) |
| every OpenAPI `operationId` / SDK method | the spec and `packages/sdk-runtime` |

Fails closed: an unverifiable claim blocks the dispatch. Everything in that table is
already available in this repo today; none of it requires judgment.

**Factory rule.** **A plan is an allegation about the repo, and allegations rot.** The
question is never "is the author careful" — it is "does the artifact carry provenance".

---

## Class P2 — The Rotting Anchor

**Symptom.** A `file:line` measured in one tree points at different code in another. Both
parties are honest and the ruling is unusable.

**Evidence.** A hub cited `:440-442` measured on `main`; the chip measured `:467-469` in
its own worktree. Same fact, incompatible anchors (HARNESS-DEBTS A-2 → A-27/A-28). Line
anchors also rotted across the three revisions of the MIS-008 fiscal plan.

**Prevention (rung 5, but cheap).** **Anchor by content** — route, `operationId`, schema
name, function name, test name. A line number is permitted only when the tree is named in
the same sentence. Mechanizable as part of the P1 claim linter: a bare `:NNN` with no tree
named is a lint error.

---

## Class P3 — The Plan With No Composition Site

**Symptom.** Every task is green, every lane passes, the feature does not exist. Nothing
in the plan ever *called* the thing it built.

**Evidence.** P2.b: 7 tasks completed, all lanes green, `icms_matrix_mirror` had **0 rows**.
A reader and a writer both shipped, with the only caller in a `_test.go`. ADR-17 working
correctly (honest-unknown, no fabricated value) *hid* the fact that the data never existed
(`plan-gate-composition-site`).

**Root cause.** A plan decomposed into artifacts, never into a *path from producer to
pixel*. Unit-level green is compositional; features are not.

**Prevention (rung 3 + rung 4).**
1. Every plan names its **composition site** — the production call path, ending at a
   screen or a row. A slice with no composition site is rejected at planning.
2. **Acceptance is a `count(*)` in the database or a rendered string on a screen** — never
   a green test. Green proves the unit; the row proves the feature.
3. Structural guard: a port implemented with **no production caller** is a lint failure
   (HARNESS-DEBTS D-15).

**Overlaps** sibling §16 (Compiles ≠ Works) at the outcome, but the cause is upstream — in
how the plan was decomposed.

---

## Class P4 — Declared Risk As Disguised Debt

**Symptom.** A plan declares a risk, the reader accepts it as the cost of moving, and the
"risk" was in fact **unmeasured work** that would have been cheap to settle.

**Evidence.** Three declared risks in the MIS-008 plan were rejected by the operator. Seven
measurement rounds later, **none was wasted and two *shrank* the work** — the fiscal
override table was cancelled entirely after measuring zero overrides in use
(`measurement-before-spec-gate`, `p2b-rescoped-estimate-not-ledger`).

**Root cause.** Declaring a risk is a legitimate move, so it is available as an escape from
measurement, and it is indistinguishable in the artifact from a genuinely irreducible one.

**Prevention (process, rung 5 — but with a mechanical form).** A declared risk must carry
**the measurement that would settle it, and why it was not run**. "Expensive" is a valid
answer; silence is not. The linter can require the field's presence, not its quality.

---

## Class P5 — The Second Engine

**Symptom.** The plan builds a new implementation beside a working one, because the symptom
appeared where the working one was not wired.

**Evidence.**
- A 17-task fiscal plan constructed a calculation layer beside `pricing/domain`, which
  already produced the right number. An empty column on screen read as "feature absent"
  instead of "feature not wired". Re-scoped to 7 tasks (`p2b-rescoped-estimate-not-ledger`).
- A market-collection ticker specified beside `syncapp.Scheduler`, discarding `sync_state`
  reconciliation, per-entity failure isolation, cursor persistence, and `RecordFailure`
  feeding the operator's health card. A broken collection would have shown a green screen
  (`anti-patterns.md`, Planning).

**Root cause.** The agent works outward from the file it opened, not from the concept's
distribution across the tree.

**Prevention (process, checkable).** Before any task creates a calculation, port, table,
component or job: **cite the `file:line` of the nearest existing thing and state why it
does not serve.** An artifact with no such citation is a plan defect. Already codified in
`.claude/skills/mc-planning` Phase 2; the missing half is the linter that fails a plan
whose new-artifact count exceeds its citation count.

**Overlaps** sibling §9 (Second Copy of a Critical Path) and §14 (Optimizing Inside a Local
Maximum). Same family; this is the *planning-time* entry point to both.

---

## Class P6 — The Plan Whose Code Was Never Compiled

**Symptom.** A plan contains complete, confident, well-formed code — the good kind of plan,
with no placeholders — and that code has never been through a compiler. It ships defects
that the language would have rejected in milliseconds, and it ships them **with authority**,
because a plan is read as a decision rather than as a draft.

**Evidence.** A 5.774-line plan specified `internal/composition/catalog_wiring.go` importing
`contexts/catalog/internal/postgres`. Go's `internal/` rule forbids it. The implementer wrote
exactly what the plan said, and `go build ./...` rejected it:
`use of internal package .../catalog/internal/postgres not allowed`. The same plan document,
eleven tasks earlier, had introduced that very rule as its central guarantee.

Earlier variants in the same family: a plan whose reader and writer had no caller outside
`_test.go` (`plan-gate-composition-site`); an agent narrating a RED it never ran
(`red-before-code-pristine-tree`).

**Root cause.** Prose and code are interleaved in one artifact, and only one of them has a
verifier. Writing plausible Go and writing correct Go are indistinguishable in a Markdown
file — which is root cause #7 (the signal is cheaper than the work) applied to planning
rather than to testing.

**Prevention (rung 2).**
1. **Any plan whose guarantee is compiler-enforced must contain a compiled skeleton.** Not
   the whole plan — the *seams*: the facade signatures, the wiring, the import graph. A
   throwaway module that builds is a few minutes and catches this whole class.
2. **The composition site is where boundary rules break.** If a plan introduces an
   `internal/`, its wiring code is the first thing to compile, not the last.
3. **Treat a plan's code blocks as untested by default** in the implementer's brief, so a
   compiler rejection reads as expected feedback rather than as a reason to work around the
   rule. Working around it is how the guarantee dies.

**Detection.** For every code block in a plan, ask: *has this exact text been through a
compiler, or does it merely look like the language?*

---

## Class P7 — The Rule Stated For An Instance Instead Of The Property

**Symptom.** A rule is discovered in one concrete setting, written down naming that setting,
and therefore not applied where the same property holds. The author violates their own rule
almost immediately — not from carelessness, but because the rule as written does not appear
to cover the new case.

**Evidence.** Rule 2.2-a was ratified as *"each vendor exposes a single root package…"*,
derived from an experiment on `adapters/marketplace/mercadolivre`. The actual property is
weaker and broader: **any tree containing an `internal/` forces a facade, because nothing
outside can name what lives inside** — contexts included. Written as a vendor rule, it was
not applied to `contexts/catalog`, and the very next composition-root file violated it (see
Class P6). Two artifacts, one author, one working session.

Related but distinct: Class D3 is prose promising *more* than the guard delivers. This is
prose promising **less** than the mechanism already enforces, so the mechanism's reach goes
unused and its violations feel legal.

**Root cause.** A rule is usually discovered through a single measurement, and the
measurement's subject gets baked into the rule's wording. Generalising requires a separate,
deliberate step that nothing prompts.

**Prevention (rung 3 + process).**
1. **After ratifying any rule, state the mechanism that enforces it and ask what else that
   mechanism already covers.** If the enforcing mechanism is broader than the rule's
   wording, the wording is wrong. Name the property, then list the instances beneath it.
2. **Write the rule at the level of the enforcer.** "Any package tree with an `internal/`"
   is checkable; "each vendor" is a category the compiler has never heard of.
3. When two rules are enforced by the same mechanism, they are one rule. Merge them.

**Detection.** Read each rule and ask: *what does the compiler/linter actually reject here,
and is that set larger than the set this sentence describes?*

---

# Part III — Honest values

## Class H1 — Unknown Becomes A Plausible Default

**Symptom.** A missing fact is rendered as `0`, `""`, `false`, or a believable rate, and
travels downstream as a measurement.

**Evidence.**
- `pricing_calc_profiles` has **zero rows**, and a 4% SIMPLES rate entered every simulation
  as if it had been configured (`calcprofile.go:19`, `p2-margin-works-tax-fabricated`).
  A zero-row configuration table is an *unconfigured* state, not a default.
- The mirror surfaced values whose provenance was a seed, not the ERP
  (`mirror-real-data-restore-d121`).

**The second half, which is not obvious and cost a round.** ADR-17 has two sides: a
**known** zero must be written as `0`, and writing it as unknown is the same defect
mirrored. Conflating "we measured zero" with "we do not know" destroys information in the
other direction (`chip-import-chain-closed`).

**Prevention (rung 1).** The type carries the distinction: `*T` / `Option<T>` / an explicit
state enum, never a zero value. Integrity-critical reads return `(T, error)`. At the schema
level: a configuration table must be able to express *"configured as zero"* distinctly from
*"absent"*.

**Overlaps** sibling §5 (Absence as Configuration) and §6 (Fallback That Fabricates Truth) —
this entry adds the inverted case, which the sibling does not carry.

---

## Class H2 — The Business Constant In Code

**Symptom.** A rate, fee, tariff or company code lives as a literal, is described as
temporary, and stays.

**Evidence.** A hardcoded `R$79` freight constant, later banned by name
(`ml-tariff-design-pending`). `CODEMP = 1` assumed where the ERP has more than one company
(D-17 / D-22). `pricing/domain/icms.go:144` hardcoding `big.NewRat(18, 100)` for intra-MG.
`pricing/domain/difal_seed.go` — a **27-row tax table** as a Go literal, with **zero
production consumers**, disagreeing with the SQL seed it claimed to mirror
(`varredura-maximo-global-instrumentos`).

**Prevention (rung 3).** A lint over `domain/` rejecting numeric literals outside a small
allowlist (0, 1, 100, powers of ten used dimensionlessly). Every exception carries a
citation comment — the law, the contract clause, the ADR. This is mechanizable and would
have caught all four instances above.

**Detection.** Grep the domain layer for decimal literals and for `%`. Every hit is either
a citation or a defect.

---

## Class H3 — Method And Data Collapsed Into One Verdict

**Symptom.** A comparison between our system and a reference disagrees, and the
investigation elects a single winner — when in fact each side is right about a *different*
half.

**Evidence.** The fiscal audit read the operator's "by the law" as being about the
**formula**; it was about the **data**. Two readings, in opposite directions, both wrong.
Decomposing into (method, data) resolved it: the ERP's *method* (simple difference) was
correct and ours was wrong; our *data* (the internal-rate table) was correct and the ERP's
was stale — 10.049 of 10.051 Bahia products carrying 17% where the law says 20,5%
(`metodo-versus-dado-nunca-colapsar`, `tgfaid-is-the-internal-rate-source`).

**Root cause.** "Who is right" is a single-slot question, and a disagreement rarely has a
single slot.

**Prevention (process).** When two systems disagree, **decompose before adjudicating**:
formula vs input vs scope vs vintage. Record a verdict *per axis*. A single-axis verdict on
a multi-axis disagreement will be wrong in a way that reads as decisive.

---

## Class H4 — Refutation Against A Single Document

**Symptom.** A hypothesis is tested against one real record, refuted, and the refutation is
wrong because that record was itself defective.

**Evidence.** A fiscal rule was refuted against one invoice. The invoice carried a stale
cadastro; its **sibling from the previous day** proved the rule. Worse: the ERP is not
self-consistent, so "the ERP disagrees" is not a refutation — ex ante the source can only
be the **matrix**, never one document (`refutation-against-single-document`). The same
shape appeared in the DIFAL check: `BASERED = BASE` in one note, which is one note.

**Prevention (process, with a mechanical floor).** A refutation requires **n > 1 and a
stated distribution**. A verdict resting on a single record is recorded as a hypothesis,
not a finding. The linter can enforce the presence of a sample size on any claim of the
form "X does not hold".

---

## Class H5 — The Sentinel Read As A Wildcard

**Symptom.** A legacy encoding uses a value to mean "this dimension is switched off"; the
implementation reads it as "matches everything". The predicate then admits rows it should
exclude, and the result is plausible.

**Evidence.** `TGFICM`'s two restriction pairs are **different axes** (operation × product),
not two interchangeable slots, and `S` is the sentinel for *axis unused* — measured, it
never carries a real code, only `-1` or `0`. Reading `S` as a wildcard accepted a rule
belonging to grupo 504 for every grupo. Correcting the predicate took 881 ambiguous cells
to **2** (`icms-matrix-two-axis-predicate`, `internal_read/domain/icms_matrix.go:90-108`).

**Prevention (rung 4 + doctrine).** For any external encoding: **enumerate the value
distribution before writing the predicate**, and record it in the code as a comment with
counts (as `icms_matrix.go:8-23` now does). A predicate written from a schema description
rather than from the data is a guess with a type signature.

---

# Part IV — Duplication and accumulation

## Class D1 — The Nth Copy Of A Concept

**Symptom.** The same idea exists in k places with k drifting definitions, and the k-th
author had no cheap way to know.

**Evidence.**
- A relative-time ladder in **three** files with three different thresholds.
- Two `FreshnessIndicator` components with different `aria-label`s, so every label-based
  test silently bound to whichever copy its file imported.
- **Three** copies of the 27-state internal-rate table: `pricing_difal_rates` (DB),
  `pricing/domain/difal_seed.go` (Go literal), `icms_aliquota_interna` (DB). Their values
  disagreed in exactly three states — and the delta was **exactly the FCP**, revealing that
  two of them answered different questions and nobody had said so
  (`varredura-maximo-global-instrumentos`).

**Prevention (rung 2 where possible, rung 3 always).** One source of truth, others
generated. Where generation is impractical: a **duplication CLI** that diffs a business
table's Go literal against its SQL seed and fails on divergence.

**Overlaps** sibling §2, §10, §19 — the sibling covers this class thoroughly. Kept here for
one addition it lacks: **copies that disagree are more informative than copies that agree.**
Their delta is a hypothesis about what the two are actually measuring, and here it settled a
tax-calculation defect by arithmetic rather than by argument.

---

## Class D2 — The Mechanism Nobody Uses

**Symptom.** A capability is built, wired end-to-end, and has zero uses — so a later author,
not finding it, builds another one.

**Evidence.** The `pricing` module accumulated **four** override mechanisms for the same
idea, none of them in use: `pricing_manual_overrides` (a table with **zero** Go references,
the only orphan table in a 62-table schema), `pricing_difal_rates.override_*` (code
complete, 0 rows using it), the frontend's inline request state, and a fifth this session
nearly added before measuring (`varredura-maximo-global-instrumentos`). Separately: **11**
contract operations with no consumer, one of them the entire feature a later plan proposed
to build from scratch (`mis008-design-closed`); and **42 of 104** SDK methods with no call
site anywhere in the frontend.

**Root cause.** Adding is locally safe and removing requires proof, so the gradient points
one way only. Over time the surface is mostly restatement.

**Prevention (rung 3).** A recurring **orphan report** in CI: exported symbols with no
non-test caller, contract operations with no client call site, tables with no code
reference. It does not block — it is read at planning time, which is when it changes a
decision. Building this once pays for itself the first time it prevents a fifth override.

**Detection caveat, learned here.** Separate *"never had a UI"* from *"the UI was removed"*.
They look identical in the report and only the second is pure garbage.

---

## Class D3 — Prose That Promises More Than The Guard

**Symptom.** A docstring, disclaimer, or contract description states a total guarantee; the
code implements a partial one. Readers stop checking.

**Evidence.** A guard covering some inputs under a docstring promising all of them —
measured worse than no guard at all, because it removes the reader's suspicion
(`chip-vinc-neutro-closed`). A global disclaimer *"seed padrão 2026"* standing in for
per-row provenance that varied by state (MIS-008 §1.5).

**Prevention — and the honest negative result, which is the point of this entry.** This
repo built a general prose gate and then **cut it by YAGNI**. Measured, a general
"prose must match code" gate does not pay: it fires constantly on comments where the cost
of imprecision is zero. The scope that *does* pay is narrow and mechanizable:
**strings rendered to a user, and published contract descriptions.** Everything else is
review vigilance (rung 6) and should be admitted as such rather than pretended into a gate
(`yagni-cut-prose-gate`).

**Factory rule derived.** Before building a gate, **measure where the value it guards comes
from.** A gate whose true blast radius is two file types should not be written to scan the
tree.

---

# Part V — The factory's own machinery

## Class F1 — The Convention Enforced By Reminder

**Symptom.** A rule that matters is documented, agreed, repeated — and its enforcement is
that someone remembers it at the right moment.

**Evidence.** Six plan amendments in a single week, each requiring the hub to remember to
amend the pack *before* the brief; a decision outside the plan does not exist for the worker
(HARNESS-DEBTS A-3). A governance registry entry that must land at merge and not before
(`governance-registry-entry-timing`). A `removal_owner` field requiring a round-trip to the
hub (D-2).

**Prevention (rung 2/3).** Encode the ordering in the mechanism: the ACK that unblocks a
worker is **emitted by the commit** that amends the pack, not by a person. Any protocol step
whose correctness depends on ordering must have the ordering enforced by the transport.

**This is the sibling catalog's central finding, restated on the process axis:** a
convention without enforcement is not a standard, it is a wish. Process conventions rot
faster than code conventions, because they have no compiler at all.

---

## Class F2 — The Checker Blind To A Legal Spelling

**Symptom.** A guard enforces a real rule and misses a syntactically different, equally
legal way to break it. Its green is read as compliance.

**Evidence.** The module-dependency checker's regex is blind to an import of a module
**root** (HARNESS-DEBTS D-27). An architecture guard keyed on a **symbol name** has a silent
back door for the renamed symbol (D-14). A `production-panic` checker did not recognise the
mandatory re-panic of `http.ErrAbortHandler` — a false *positive* of the same blindness
(D-9). A baselined panic exception is not inherited by a new call site
(`p2b-arquitetura-e-governanca`).

**Prevention (rung 1 > rung 3).** Prefer checks over **structure** (AST, import graph, type
graph) to checks over **text**. Where a regex is unavoidable, its test suite must include
the legal-but-different spellings, and the check ships with a positive control (Class V7).

---

## Class F3 — The Alarm Wrong Twice

**Symptom.** A gate produces false accusations. The cost is not the wasted investigation —
it is that a reader who has dismissed two false alarms will skip the third, which is real.

**Evidence.** A stop hook produced its **third** false accusation of the same class, during
the very session convened to fix it (HARNESS-DEBTS D-7). A governance lane that **always**
fails because the shortcut does not forward `-BaseSha` (B-11) — a permanently red lane is
informationally identical to no lane.

**Prevention (process, and it is a hard rule).** A gate's false-positive rate is a
first-class defect of the gate. **Two false accusations of the same class disable the gate**
until its precision is fixed. A permanently-red lane is treated as an outage, not as a
backlog item.

**Related, from the sibling:** an allowlist is the tempting fix here and is the wrong one
(§3) — it converts imprecision into sanctioned drift.

---

## Class F4 — The Wedge With No Degradation Signal

**Symptom.** A shared facility fails in a way that produces no error, only absence — and
every downstream consumer reports success at doing nothing.

**Evidence.** The harness shell wedges for the **entire session, subagents included**, with
no degradation signal; the tell is a command signature that changes mid-line, meaning
truncation. The mitigation is a `true` probe at bootstrap
(`harness-shell-wedge-b8`, HARNESS-DEBTS B-8). Related: `.dockerignore` letting 940 MB of Go
cache into the build context (D-23); a policy scan walking untracked worktrees and hanging
past 20 minutes (B-10/B-10b).

**Prevention (rung 3).** Every shared facility has a **liveness probe run at session start**
whose failure is loud. Absence is never allowed to be the only symptom — this is Class V1
one layer down, at the tooling rather than the test.

---

## Class F5 — The Rule Specified Against A Language That Cannot Express It

**Symptom.** A standard mandates a construct the implementation language cannot compile. It
is ratified, cited, and planned against; the contradiction only surfaces when someone finally
types it. Every artifact that referenced the rule in the meantime is now wrong.

**Evidence.** Spec Rule 4.4 required that arithmetic over `Fact[T]` live in **methods** of
`Fact[T]`. Go does not permit a method to declare its own type parameters, so
`func (f Fact[T]) Map[U](...)` does not compile; the operation must be a package-level
function. Rule 4.1 in the same spec asserted "there is no struct literal", but `Fact[T]{}`
compiles from any package — an empty composite literal names no fields, so unexported fields
do not block it. The real guard is `Evidence().IsZero()`. Both were ratified without a
compiling probe.

Adjacent: a criterion vacuous against the type it constrained
(`criterion-unfalsifiable-against-type`) — there the type was too weak to violate the rule;
here the type system forbids the rule outright.

**Root cause.** Specification and compilation are separated by days. Design documents are
reviewed by readers, and readers cannot typecheck. The more sophisticated the type-level
guarantee, the more likely the language has a rule about it that prose does not know.

**Prevention (rung 1/2).**
1. **Any rule that mandates a language construct ships a compiling snippet in the spec.**
   Five lines in a scratch module. If it does not build, the rule is not ratifiable.
2. **Highest risk lives at generics, embedding, interface satisfaction, visibility and
   initialisation order** — the places where every language has a surprising restriction.
   Treat a rule touching those as unproven until compiled.
3. **A false sentence in a ratified document is deleted, not annotated.** An amendment
   alongside the falsehood leaves two readings and the reader picks one.

**Detection.** For each normative sentence about code shape, ask: *has anyone compiled an
example of this?* Absent an artifact, the rule is a hypothesis.

---

## Appendix A — Day-0 Factory Checklist

Each line prevents a class already observed **in this repo**. Complementary to the sibling
catalog's Appendix A; read both.

**Dispatch (before an agent starts)**
- [ ] Claim linter over every plan/card/brief: shell commands resolve, `file:line` exists
      and names its tree, symbols exist, tables and columns exist in the live schema, spec
      operations and SDK methods exist (§P1, §P2)
- [ ] A plan naming a new artifact carries a citation of the nearest existing one (§P5)
- [ ] Every plan names its **composition site** — the production path ending in a row or a
      pixel (§P3)
- [ ] Declared risks carry the measurement that would settle them (§P4)
- [ ] `tools(reviewer) ⊇ capabilities required by its assigned criteria` (§V5)
- [ ] A plan whose guarantee is compiler-enforced ships a **compiled skeleton of its seams**
      — facade signatures, wiring, import graph (§P6)
- [ ] The composition site compiles **first**, not last, whenever a plan introduces an
      `internal/` (§P6)

**Lanes and evidence**
- [ ] A lane reporting **zero executed units fails** (§V1)
- [ ] Lanes emit an attributable per-test token, never just an exit code (§V1)
- [ ] Preflight prints absolute path + tip SHA + container build time, and **fails if the
      container predates the tip** (§V4)
- [ ] Every zero-result check ships a positive control in the same run (§V7)
- [ ] Every unattended failure has a named path to a pixel; kill-the-dependency drill is
      part of acceptance (§V8)
- [ ] Second occurrence of a defect shape stops the line — class fix or registered debt,
      never a third patch (§V9)
- [ ] Every reported number ships the command that produced it, as an artifact (§V3)
- [ ] Estimator sources are banned from verdicts by name (§V3)
- [ ] A verdict that carries a count requires a second, independent instrument (§V3)
- [ ] Merge gates diff against the **tip of the target**, always (§V6)
- [ ] Every shared facility has a loud liveness probe at session start (§F4)
- [ ] Every registry/manifest rule has a **paired tree rule** that walks the filesystem and
      reports entities present but unenrolled (§V10)
- [ ] Checkers print their **universe**, not just their verdict — scanned N, found M (§V10)
- [ ] A must-fail proof records a **green baseline first**; a non-zero baseline aborts it,
      and the FAIL output must name the injected decoy (§V11)

**Tests**
- [ ] Red-before-green enforced by artifact: the recorded failure names the test added (§V2)
- [ ] A port implemented with no production caller is a lint failure (§P3)
- [ ] Acceptance is a `count(*)` or a rendered string — never a green test (§P3)

**Types and domain**
- [ ] Unknown is a type-level state, and "known zero" is expressible distinctly (§H1)
- [ ] Numeric-literal lint over `domain/`; exceptions carry a citation (§H2)
- [ ] Predicates over external encodings are written from a measured value distribution,
      recorded in-code with counts (§H5)

**Sweeps that change decisions**
- [ ] Recurring orphan report: symbols, contract operations, tables with no caller (§D2)
- [ ] Duplication check between business tables in code and their SQL seeds (§D1)

**Gate hygiene**
- [ ] Two false accusations of one class **disable the gate** until precision is fixed (§F3)
- [ ] A permanently-red lane is an outage, not a backlog item (§F3)
- [ ] Prefer structural checks (AST, import graph) over text checks; regex checks ship
      legal-but-different spellings in their test suite (§F2)
- [ ] Protocol ordering is enforced by the transport, not by memory (§F1)
- [ ] Measure where a guarded value comes from **before** sizing the guard (§D3)
- [ ] A rule mandating a language construct ships a **compiling snippet**; generics,
      embedding, visibility and init order are unproven until compiled (§F5)
- [ ] After ratifying a rule, state its enforcing mechanism and widen the wording to
      everything that mechanism already covers (§P7)
- [ ] A false sentence in a ratified document is **deleted**, never annotated (§F5)

---

## Appendix B — Recurring Root Causes

The sibling catalog reduces its 25 classes to five underlying causes. Those five hold here
too. This axis adds **three more**, which do not appear until an agent is doing the work:

6. **The agent produced text about the repository instead of reading the repository.**
   Plausible, well-formed, and false — column names, line anchors, lane commands, symbol
   names, brief claims. The defect is not carelessness; it is that generating a plausible
   fact and retrieving a true one are indistinguishable in the output.
   → §P1, §P2, §V3, §H4

7. **The signal that proves work was done is cheaper to fabricate than the work.**
   A green lane that ran nothing, an assertion that cannot fail, a zero from a blind
   instrument, a self-certified criterion, a stale binary. Every one of these is reachable
   without doing the task, and none of them looks different from success — including in
   production, where "the process is alive" is rendered as "the system is healthy".
   → §V1, §V2, §V4, §V5, §V7, §V8, §P3

8. **A rule written in prose is never executed, so nothing objects when it is false,
   unexpressible, or narrower than the mechanism that enforces it.** Code has a compiler;
   standards have readers. Every normative document accumulates sentences that would not
   survive one minute of execution, and they are cited as authority in the meantime — by
   the author most of all.
   → §P6, §P7, §F5, §D3, §F1

**The factory's design principle, in one line:** the causes above are not fixed by
better prompts or better agents. They are fixed by making the honest signal the cheapest
one available — which is an engineering problem, and therefore solvable.
