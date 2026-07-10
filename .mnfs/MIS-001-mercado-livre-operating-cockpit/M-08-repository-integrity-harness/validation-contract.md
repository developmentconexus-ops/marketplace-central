# Milestone Validation Contract

```yaml
id: M-08
type: milestone-validation-contract
status: planned
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

A clean accepted baseline and isolated execution lanes make subsequent
milestone worktrees reproducible without contaminating dev/live state.

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
- Expected: exit 0 with `.env`, PostgreSQL, Oracle, provider network, and migration execution disabled.
- Actual:
- Artifact: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-08-repository-integrity-harness/validation-result.md`
Blocking failure: Unit execution reads live configuration or external state.
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

## Retry Policy

- correction_attempts:
- max_correction_attempts: 2
- last_validation_result:

## Handoff

- Current status: Pending execution.
- Next owner: QA Validator after milestone implementation.
- Next action: Validate M-08-C01 through M-08-C07 against current-state evidence.
- Required files/evidence: F-*/validation.md and `validation-result.md`.
- Blockers or open decisions: None during planning.

