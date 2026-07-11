# Milestone Validation Contract

```yaml
id: M-08
type: milestone-validation-contract
status: in_progress
owner: Mission Strategist
parent: MIS-001
created: 2026-07-10
updated: 2026-07-10
validation_level: QA-3
lifecycle_scope: milestone
```

## Milestone ID

M-08

## QA Level

QA-3: repository integrity, deterministic integration, and operational harness.

## Required Outcome

A clean accepted baseline plus canonical knowledge, compiled context, isolated
execution, and risk-aware orchestration make subsequent work reproducible
without hidden history, duplicated truth, or dev/live contamination.

## Criteria

## Criterion: Baseline Integrity
ID: M-08-C01
Level: Milestone
Type: Engineering
Required: Yes
Status: Pending
Evidence:
- Command: `git status --short` plus baseline inventory verifier.
- Expected: empty status at accepted SHA; every original dirty path mapped to an intentional commit or an explicit retained-state record.
- Actual:
- Artifact: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-08-repository-integrity-harness/validation-result.md`
Blocking failure: Any original path is lost, silently reverted, or unattributed.
Blocking failure observed: No
Owner: QA Validator

## Criterion: Hermetic Unit Lane
ID: M-08-C02
Level: Milestone
Type: Reliability
Required: Yes
Status: Pending
Evidence:
- Command: versioned harness unit command.
- Expected: exit 0 from a child process whose environment contains only the declared unit allowlist; no configured PostgreSQL, Oracle, OAuth, provider, proxy, migration, live, tenant, or tunnel key reaches test subprocesses.
- Actual:
- Artifact: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-08-repository-integrity-harness/validation-result.md`
Blocking failure: Unit execution reads live configuration or external state.
Blocking failure observed: No
Owner: QA Validator

## Criterion: Knowledge Authority Cutover
ID: M-08-C08
Level: Milestone
Type: Architecture
Required: Yes
Status: Pending
Evidence:
- Command: tracked-file and active-reference audit after the authority migration.
- Expected: `git ls-files -- .brain` returns zero paths; active guidance resolves ADRs to `docs/architecture/decisions/`, execution status to `.mnfs`, and machine-owned facts to `contracts/governance`; historical evidence remains available through Git history.
- Actual:
- Artifact: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-08-repository-integrity-harness/validation-result.md`
Blocking failure: `.brain` remains active, a current unique decision is lost, or two writable sources own the same status/configuration/invariant fact.
Blocking failure observed: No
Owner: QA Validator

## Criterion: Governance and Context Contract
ID: M-08-C09
Level: Milestone
Type: Engineering
Required: Yes
Status: Pending
Evidence:
- Command: governance schema/drift suite plus context compilation from a pinned MNFS feature and base SHA.
- Expected: every registry validates; semantic drift exits 0; context pack records source hashes, risk, criterion mapping, allowed/forbidden paths, side effects, commands, evidence types, and stop conditions; base/source mutation makes validation exit non-zero.
- Actual:
- Artifact: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-08-repository-integrity-harness/validation-result.md`
Blocking failure: A pack starts from stale, incomplete, contradictory, or unowned context.
Blocking failure observed: No
Owner: QA Validator

## Criterion: Allowed Paths and Shared-Seam Leases
ID: M-08-C10
Level: Milestone
Type: Reliability
Required: Yes
Status: Pending
Evidence:
- Command: isolated lease/path fixture with two competing run IDs plus recovery check.
- Expected: first run acquires its declared seam; second run exits non-zero without writing; an out-of-scope changed path exits non-zero; stale recovery preserves worktree/files and requires an explicit disposition.
- Actual:
- Artifact: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-08-repository-integrity-harness/validation-result.md`
Blocking failure: Two writers can own one seam, an undeclared path is accepted, or recovery deletes/resets work.
Blocking failure observed: No
Owner: QA Validator

## Criterion: Goal and Codex Capability Contract
ID: M-08-C11
Level: Milestone
Type: Compatibility
Required: Yes
Status: Pending
Evidence:
- Command: fresh-task capability spike followed by repo-skill dry run against a pinned feature.
- Expected: installed host reports supported project agent, hook, model, reasoning, and permission surfaces; supported controls execute with recorded output; unsupported controls are advisory and the script fallback blocks the same unsafe case; repo skill reconciles active `/goal` into MNFS and emits a schema-valid, hash-current, criterion-complete context pack without hidden transcript.
- Actual:
- Artifact: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-08-repository-integrity-harness/validation-result.md`
Blocking failure: Completion depends on an unsupported Codex surface or hidden conversation state.
Blocking failure observed: No
Owner: QA Validator

## Criterion: Harness Regression Corpus
ID: M-08-C12
Level: Milestone
Type: QA
Required: Yes
Status: Pending
Evidence:
- Command: versioned agent-harness eval suite followed by cold gate from a clean worktree.
- Expected: all declared positive and negative cases receive the expected deterministic verdict; result records case ID, duration, target type, exit code, and artifact path; no case promotes fake evidence to live.
- Actual:
- Artifact: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-08-repository-integrity-harness/validation-result.md`
Blocking failure: A known architecture, environment, evidence, path, seam, or provider-write regression is accepted.
Blocking failure observed: No
Owner: QA Validator

## Criterion: Ephemeral PostgreSQL Isolation
ID: M-08-C03
Level: Milestone
Type: Reliability
Required: Yes
Status: Pending
Evidence:
- Command: versioned harness integration command followed by database inventory.
- Expected: generated database name starts `mpc_test_`; migrations apply then apply zero on second run; test resources are removed in `finally`; dev row counts are unchanged.
- Actual:
- Artifact: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-08-repository-integrity-harness/validation-result.md`
Blocking failure: Integration tests can target or mutate the persistent development database.
Blocking failure observed: No
Owner: QA Validator

## Criterion: Live Lane Safety
ID: M-08-C04
Level: Milestone
Type: Security
Required: Yes
Status: Pending
Evidence:
- Command: live preflight and redaction verifier without provider write.
- Expected: only allowlisted variable names are consumed; no secret value or buyer PII appears in stdout, logs, manifest, or committed evidence; provider write requires a separate explicit command and actor.
- Actual:
- Artifact: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-08-repository-integrity-harness/validation-result.md`
Blocking failure: Secret/PII disclosure or an implicit provider write.
Blocking failure observed: No
Owner: QA Validator

## Criterion: Cold Gate Reproducibility
ID: M-08-C05
Level: Milestone
Type: QA
Required: Yes
Status: Pending
Evidence:
- Command: versioned cold-gate command executed twice from accepted SHA.
- Expected: both runs report identical command inventory and exit classifications; each produces a run ID, SHA, branch, dirty flag, tool versions, target types, exit codes, and redacted artifact paths.
- Actual:
- Artifact: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-08-repository-integrity-harness/validation-result.md`
Blocking failure: Gate result depends on hidden environment or omits a required workspace.
Blocking failure observed: No
Owner: QA Validator

## Criterion: Worktree Isolation
ID: M-08-C06
Level: Milestone
Type: Compatibility
Required: Yes
Status: Pending
Evidence:
- Command: create one native Codex worktree from accepted SHA and run status/baseline checks.
- Expected: independent Compose project, ports/database namespace, clean git state, and passing baseline without modifying the primary checkout.
- Actual:
- Artifact: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-08-repository-integrity-harness/validation-result.md`
Blocking failure: Worktree shares mutable runtime identity or cannot reproduce baseline.
Blocking failure observed: No
Owner: QA Validator

## Criterion: Session Contract
ID: M-08-C07
Level: Milestone
Type: Documentation
Required: Yes
Status: Pending
Evidence:
- Command: context-pack/skill validation plus fresh-thread dry run.
- Expected: fresh thread receives objective, reads, constraints, owned/forbidden paths, gates, evidence types, commit rule, and compact handoff without repository-wide discovery.
- Actual:
- Artifact: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-08-repository-integrity-harness/validation-result.md`
Blocking failure: Worker requires hidden chat history or broad rediscovery to start.
Blocking failure observed: No
Owner: QA Validator

## Evidence Requirements

- Commands identify `fake`, `ephemeral-postgres`, `live-oracle`,
  `live-provider`, or `browser` target.
- Raw logs live under an ignored run-ID directory; MNFS records concise result
  and path only.
- No mock or compile-only output is accepted as live evidence.

## Blocking Failures

- Lost or reverted user change.
- Dev/live database mutation by integration tests.
- Secret or buyer PII exposure.
- Worktree created from an unaccepted baseline.
- Gate that excludes an active Go or npm workspace.
- Competing knowledge owners, stale context acceptance, or shared-seam conflict.
- Mandatory reliance on an unsupported Codex capability.

## Retry Policy

- correction_attempts:
- max_correction_attempts: 2
- last_validation_result:

## Handoff

- Current status: In progress; contract expanded after approved harness replan.
- Next owner: QA Validator after milestone implementation.
- Next action: Validate M-08-C01 through M-08-C12 against current-state evidence.
- Required files/evidence: F-*/validation.md and `validation-result.md`.
- Blockers or open decisions: None during planning.
