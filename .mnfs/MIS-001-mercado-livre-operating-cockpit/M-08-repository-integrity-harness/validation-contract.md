# Milestone Validation Contract

```yaml
id: M-08
type: milestone-validation-contract
status: in_progress
owner: Mission Strategist
parent: MIS-001
created: 2026-07-10
updated: 2026-07-11
validation_level: QA-3
lifecycle_scope: milestone
```

## Milestone ID

M-08

## QA Level

QA-3: repository integrity, autonomous development control, and real-target evidence safety.

## Required Outcome

A goal can run for a long duration through visible milestone orchestration,
bounded subagents, exact repository knowledge, proportional engineering gates,
and artifact-backed resume. The harness improves development speed and quality
without cold-cloning the repository or simulating a clean machine locally.

## Criteria

### Supersession Ledger

| Historical ID | Decision | Date | Active replacement |
| --- | --- | --- | --- |
| M-08-C05 | Excluded after the operator rejected local clean-machine simulation; the original definition is preserved below. | 2026-07-11 | M-08-C15 and M-08-C17 |
| M-08-C12 | Excluded because its original proof depended on the cold gate; its deterministic-regression intent continues. | 2026-07-11 | M-08-C16 and M-08-C17 |

The active required set is C01-C04, C06-C11, and C13-C17. Superseded IDs are
audited for historical honesty and absence from active execution; they are not
counted as passed requirements.

## Criterion: Baseline Integrity
ID: M-08-C01
Level: Milestone
Type: Engineering
Required: Yes
Status: Pending
Evidence:
- Command: `git status --short` plus baseline inventory verifier.
- Expected: empty status at accepted SHA; every original dirty path mapped to an intentional commit or explicit retained-state record.
- Actual:
- Artifact: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-08-repository-integrity-harness/validation-result.md`
Blocking failure: Any original path is lost, silently reverted, or unattributed.
Blocking failure observed: No
Owner: QA Validator

## Criterion: Safe Unit Lane
ID: M-08-C02
Level: Milestone
Type: Reliability
Required: Yes
Status: Pending
Evidence:
- Command: versioned harness unit command.
- Expected: exit 0 from a child process containing only safe tool keys and declared unit values; no PostgreSQL, Oracle, OAuth, provider, tenant, proxy, migration, or tunnel value reaches tests.
- Actual:
- Artifact: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-08-repository-integrity-harness/validation-result.md`
Blocking failure: Unit execution consumes live configuration or external state.
Blocking failure observed: No
Owner: QA Validator

## Criterion: Ephemeral PostgreSQL
ID: M-08-C03
Level: Milestone
Type: Reliability
Required: Yes
Status: Pending
Evidence:
- Command: versioned integration command followed by database/resource inventory.
- Expected: generated `mpc_test_*` database, first migration apply count `32`, second apply count `0`, cleanup in `finally`, and unchanged development digest.
- Actual:
- Artifact: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-08-repository-integrity-harness/validation-result.md`
Blocking failure: Integration can target or mutate persistent development/live data or leak a resource.
Blocking failure observed: No
Owner: QA Validator

## Criterion: Real-Target Evidence Safety
ID: M-08-C04
Level: Milestone
Type: Security
Required: Yes
Status: Pending
Evidence:
- Command: lane preflights plus selected real Oracle/provider/browser checks required by dogfood; no provider write.
- Expected: each result names exactly one of `fake`, `ephemeral-postgres`, `live-oracle`, `live-provider`, or `browser`; no secret/buyer PII is persisted; fake cannot satisfy a real criterion.
- Actual:
- Artifact: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-08-repository-integrity-harness/validation-result.md`
Blocking failure: Secret/PII disclosure, evidence-class inflation, or implicit provider write.
Blocking failure observed: No
Owner: QA Validator

### Historical immutable definition: M-08-C05

> Criterion: Cold Gate Reproducibility
> ID: M-08-C05
> Level: Milestone
> Type: QA
> Required: Yes
> Status: Pending
> Evidence:
> - Command: versioned cold-gate command executed twice from accepted SHA.
> - Expected: both runs report identical command inventory and exit classifications; each produces a run ID, SHA, branch, dirty flag, tool versions, target types, exit codes, and redacted artifact paths.
> - Actual:
> - Artifact: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-08-repository-integrity-harness/validation-result.md`
> Blocking failure: Gate result depends on hidden environment or omits a required workspace.
> Blocking failure observed: No
> Owner: QA Validator

## Criterion: Worktree and Writer Ownership
ID: M-08-C06
Level: Milestone
Type: Compatibility
Required: Yes
Status: Pending
Evidence:
- Command: native worktree task fixture plus lease/path checks.
- Expected: worktree starts from named accepted SHA, uses a separate branch/detached task state, reuses normal developer dependencies, and leaves primary checkout unchanged.
- Actual:
- Artifact: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-08-repository-integrity-harness/validation-result.md`
Blocking failure: Two writers share one checkout/seam or worktree operations mutate primary user work.
Blocking failure observed: No
Owner: QA Validator

## Criterion: Session Contract
ID: M-08-C07
Level: Milestone
Type: Documentation
Required: Yes
Status: Pending
Evidence:
- Command: context/dispatch schema validation plus fresh-task dry run.
- Expected: fresh task receives objective, exact read selectors, constraints, paths, seams, gates, evidence targets, commit rule, stop conditions, and compact return schema without broad discovery.
- Actual:
- Artifact: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-08-repository-integrity-harness/validation-result.md`
Blocking failure: Worker needs hidden transcript or repository-wide rediscovery to begin.
Blocking failure observed: No
Owner: QA Validator

## Criterion: Knowledge Authority
ID: M-08-C08
Level: Milestone
Type: Architecture
Required: Yes
Status: Pending
Evidence:
- Command: tracked-file and active-reference audit.
- Expected: `.brain` has zero tracked paths; architecture, execution, machine facts, wiki, contracts, and code each have one declared owner.
- Actual:
- Artifact: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-08-repository-integrity-harness/validation-result.md`
Blocking failure: Two writable sources own one fact or a current unique decision is lost.
Blocking failure observed: No
Owner: QA Validator

## Criterion: Governance, Knowledge Routes, and Context
ID: M-08-C09
Level: Milestone
Type: Engineering
Required: Yes
Status: Pending
Evidence:
- Command: governance schema/drift plus representative context compilation/current validation.
- Expected: routes resolve existing module/interface/consumer/test selectors;
  the measured initial set includes every harness-requested bootstrap and
  selector read; it records source hashes, risk, criteria, paths, seams, side
  effects, commands, and stop conditions; target is at most 2,000 estimated
  tokens, while a larger L2/L3 set must name an `overflow_reason` and only the
  additional required sources; source/base mutation exits nonzero.
- Actual:
- Artifact: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-08-repository-integrity-harness/validation-result.md`
Blocking failure: Pack is stale, incomplete, contradictory, has unjustified overflow, or loads unrelated module context.
Blocking failure observed: No
Owner: QA Validator

## Criterion: Allowed Paths and Shared-Seam Leases
ID: M-08-C10
Level: Milestone
Type: Reliability
Required: Yes
Status: Pending
Evidence:
- Command: two-writer seam fixture, out-of-scope change fixture, checkpoint/resume recovery.
- Expected: first writer obtains lease; conflicting writer exits nonzero before writing; undeclared changed path fails acceptance; stale recovery preserves all files and requires explicit disposition.
- Actual:
- Artifact: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-08-repository-integrity-harness/validation-result.md`
Blocking failure: A conflict or undeclared path is accepted, or recovery deletes/resets work.
Blocking failure observed: No
Owner: QA Validator

## Criterion: Native Goal and Codex Capability Contract
ID: M-08-C11
Level: Milestone
Type: Compatibility
Required: Yes
Status: Pending
Evidence:
- Command: fresh task invokes repo skill; supported task/subagent/read/steer/interrupt/worktree controls are exercised; unsupported controls use labelled fallback.
- Expected: goal reconciles to MNFS, context compiles, one milestone task and bounded subagent are observable, and no acceptance depends on experimental app-server methods.
- Actual:
- Artifact: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-08-repository-integrity-harness/validation-result.md`
Blocking failure: Workflow depends on hidden chat state or unsupported/experimental-only behavior without fallback.
Blocking failure observed: No
Owner: QA Validator

### Historical immutable definition: M-08-C12

> Criterion: Harness Regression Corpus
> ID: M-08-C12
> Level: Milestone
> Type: QA
> Required: Yes
> Status: Pending
> Evidence:
> - Command: versioned agent-harness eval suite followed by cold gate from a clean worktree.
> - Expected: all declared positive and negative cases receive the expected deterministic verdict; result records case ID, duration, target type, exit code, and artifact path; no case promotes fake evidence to live.
> - Actual:
> - Artifact: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-08-repository-integrity-harness/validation-result.md`
> Blocking failure: A known architecture, environment, evidence, path, seam, or provider-write regression is accepted.
> Blocking failure observed: No
> Owner: QA Validator

## Criterion: Goal Routing and Global-Maximum Planning
ID: M-08-C13
Level: Milestone
Type: Architecture
Required: Yes
Status: Pending
Evidence:
- Command: fixtures for one-module, cross-module, API/SDK, database, and real-provider goals.
- Expected: every fixture selects one owner, current interfaces/consumers, wrong-legacy decision, smallest durable abstraction, risk level, and exact proof commands; unknown cost/linkage remains explicit.
- Actual:
- Artifact: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-08-repository-integrity-harness/validation-result.md`
Blocking failure: Planner permits a local workaround, wrong module ownership, duplicated policy, provider-specific core abstraction, or unknown-to-zero.
Blocking failure observed: No
Owner: QA Validator

## Criterion: Long-Running Orchestration and Resume
ID: M-08-C14
Level: Milestone
Type: Reliability
Required: Yes
Status: Pending
Evidence:
- Command: visible milestone dogfood with bounded subagent, checkpoint, steering/interrupt path, and fresh continuation.
- Expected: checkpoint ID, base/current SHA, feature state, completed evidence,
  blocker, and next action persist; a native task/thread ID may appear only as
  correlation metadata; continuation states the same next action without the
  portfolio transcript.
- Actual:
- Artifact: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-08-repository-integrity-harness/validation-result.md`
Blocking failure: Work cannot resume from files or task control drifts from MNFS state.
Blocking failure observed: No
Owner: QA Validator

## Criterion: Risk-Adaptive Impact and QA Routing
ID: M-08-C15
Level: Milestone
Type: QA
Required: Yes
Status: Pending
Evidence:
- Command: L0-L3 routing fixtures plus current-checkout impact gate.
- Expected: L0 runs deterministic checks; L1 focused tests plus one review; L2 contract gate plus one fixed-commit review; L3 uses one serial writer followed by parallel independent SPEC/SAFETY and QUALITY reviews on the same fixed commit and the named real target; impact gate executes exactly registered commands selected by the pack.
- Actual:
- Artifact: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-08-repository-integrity-harness/validation-result.md`
Blocking failure: Low-risk work receives universal ceremony, high-risk work skips required evidence, or undeclared commands execute.
Blocking failure observed: No
Owner: QA Validator

## Criterion: Harness Eval and Efficiency Dogfood
ID: M-08-C16
Level: Milestone
Type: Performance
Required: Yes
Status: Pending
Evidence:
- Command: versioned eval corpus plus fresh-task dogfood and resume.
- Expected: every case emits pinned verdict/reason/target; the representative
  bootstrap-plus-selector read set meets the 2,000-token target or carries a
  justified L2/L3 overflow; unrelated module sources count `0`; every route
  miss is recorded and resolved or explicitly blocked; no more than one
  correction batch.
- Actual:
- Artifact: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-08-repository-integrity-harness/validation-result.md`
Blocking failure: Known regression is accepted, context is unrelated or unjustifiably large, or dogfood depends on hidden history.
Blocking failure observed: No
Owner: QA Validator

## Criterion: Pragmatic Cutover Integrity
ID: M-08-C17
Level: Milestone
Type: Engineering
Required: Yes
Status: Pending
Evidence:
- Command: governance/alias/boundary search after F-10.
- Expected: zero active cold command, npm alias, cold-provision lane, snapshot/provisioning implementation, cold-only test, or required cold criterion; historical F-04 artifacts remain.
- Actual:
- Artifact: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-08-repository-integrity-harness/validation-result.md`
Blocking failure: Active execution still contains cold-clone/clean-cache behavior or history was erased.
Blocking failure observed: No
Owner: QA Validator

## Evidence Requirements

- Feature execution evidence: `<feature-root>/validation.md`.
- Milestone QA rollup: this milestone's `validation-result.md`.
- Commands record target, exit, decisive output, and relative artifact path.
- Raw logs remain ignored and redacted; secrets and buyer PII never enter MNFS.
- Mock/fake/compile/preflight evidence never proves a real target.
- Superseded C05/C12 are audited for absence, not counted as passed requirements.

## Blocking Failures

- Lost/reverted user change or destructive recovery.
- Hidden transcript required to start or resume.
- Wrong module/interface ownership, unknown-to-zero, or local workaround accepted.
- Stale/over-budget/unrelated context accepted.
- Parallel writer conflict or undeclared path accepted.
- Secret/PII disclosure, evidence inflation, or implicit provider write.
- Universal cold gate remains active or product work is blocked on clean-machine simulation.
- Required real Oracle/provider/browser/database behavior is claimed from fake evidence.

## Retry Policy

- correction_attempts: 0
- max_correction_attempts: 2
- last_validation_result:

## Handoff

- Current status: Replanned and pending F-10, F-05, and F-09 execution.
- Next owner: Milestone Orchestrator, then independent QA Validator.
- Next action: Execute F-10 pragmatic cutover.
- Required files/evidence: accepted feature validations and `validation-result.md` covering required criteria.
- Blockers or open decisions: None.
