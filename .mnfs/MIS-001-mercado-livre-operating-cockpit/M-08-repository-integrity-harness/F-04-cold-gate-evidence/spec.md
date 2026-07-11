# F-04 Deterministic Cold Gate and Evidence Manifest — Specification

```yaml
id: F-04
type: feature-spec
status: spec_ready
owner: Feature Implementer
parent: M-08
created: 2026-07-11
updated: 2026-07-11
validation_level: QA-3
lifecycle_scope: feature
```

## Problem

The accepted harness has hermetic unit and ephemeral-PostgreSQL lanes, but no
single command proves a clean candidate implementation SHA can provision its dependencies,
exercise every deterministic workspace, and leave a target-labelled redacted
outcome. Existing summaries omit command order, tool identity, aggregate
classification, and reproducibility across cold runs.

## Requirements

### Clean isolated source and provisioning

- A versioned local/future-CI command requires a candidate 40-hex SHA equal to
  clean caller `HEAD`, then creates a detached local clone under
  ignored `scripts/.runs/<run-id>/snapshot`. A local clone is chosen over
  `git archive` because governance drift needs Git metadata. It must not create
  worktree metadata or write the primary checkout. It clones committed Git
  objects only, checks out that exact candidate detached, then verifies clone
  `HEAD`, clean status, and absence of mutable source-worktree content.
  Before cloning, the parent resolves the canonical expected checkout root and
  requires its direct child `.git` to be a normal directory. Gitfiles,
  reparse-point gitdirs, external/ancestor gitdirs, and any non-exact gitdir
  fail closed as `GITDIR_UNTRUSTED`.
- All snapshot `node_modules`, npm cache, `GOCACHE`, and `GOMODCACHE` paths live
  under that run. Caller dependencies/caches are neither read nor changed.
- A registry-defined `cold-provision` lane receives only safe tool variables,
  no `.env`, proxy, registry credential, Oracle, OAuth, provider, tenant, or
  database value. Its evidence class is `provisioning`, never fake/integration/live.
- Provision in frozen order: Go modules from `go.sum`, npm workspaces from
  `package-lock.json`, then `docker pull postgres:16-bookworm`. Record the
  resolved PostgreSQL image identity without credentials or absolute paths.
  If a private credential is required, the image identity is unresolved or
  changes unexpectedly between the comparison runs, or provisioning would
  mutate the primary checkout/caller caches, stop blocked.

### Deterministic validation and aggregation

- After provisioning, network-disabled validation consumes the accepted
  Environment/Execution/Policy/Context seams and the F-03 integration lane.
  The ordered inventory covers governance contract/drift, hermetic unit tests,
  all npm workspaces through root test/build aliases, and F-03 ephemeral
  PostgreSQL with its existing `--pull=never` behavior.
- An active workspace may not be skipped silently. A missing alias, lockfile,
  tool, dependency, image, workspace, or expected command is a stable nonzero
  reason.
- Continue independent later stages after a failure when safe, preserve every
  completed record, and classify the aggregate nonzero if any required stage
  fails, blocks, times out, leaks resources, or cannot be evidenced.

### Evidence contract

- `Evidence.psm1` writes an ignored, redacted JSONL trace and a JSON outcome
  manifest validated by a committed schema. Raw subprocess logs remain ignored
  and are not copied into MNFS.
- The manifest records schema version, run ID, candidate SHA, source branch,
  dirty flag, later acceptance-link field (unset during candidate execution),
  sorted tool versions, exact ordered command inventory, stage, target label,
  evidence class, duration, exit code/reason, repository-relative artifact
  paths, PostgreSQL image identity, and aggregate classification.
- Evidence accepts only declared target/evidence-class pairs. Provisioning can
  establish readiness but cannot prove a fake, integration, live, browser, or
  production-like outcome. No live Oracle/provider/browser command or provider
  write is part of this gate.
- Redaction rejects secret-like values, buyer PII, database URLs, credentials,
  and absolute paths before persistence. Committed validation stores only safe
  summaries and repository-relative ignored artifact paths.
- Two successful runs from the same candidate SHA may differ in run ID and
  duration, but their projected ordered inventory, target/evidence classes,
  image identity, per-command exit classifications, and aggregate
  classification must be identical.

## Non-Goals

- No rewrite of F-03, `Execution.psm1`, `Environment.psm1`, `Policy.psm1`, or
  `Context.psm1` absent an architecture contradiction.
- No live Oracle/provider/browser execution, credentials, provider read/write,
  dev database mutation, Docker dev-stack change, dependency/lockfile update,
  application code, OpenAPI/SDK, or UI work.
- F-09 remains the final consumer for `M-08-C12`; F-04 supplies cold-gate
  evidence but does not run or accept the complete harness eval corpus.

## Edge Cases

- Wrong/dirty SHA fails before snapshot or provisioning and still emits a safe
  blocked outcome.
- Provision failure does not run dependent validation but does preserve safe
  preflight/provision records; unrelated safe evidence may continue.
- Test failure followed by cleanup failure records both and keeps the primary
  failure visible.
- Unknown target/evidence pair, malformed event, unsafe artifact path, redaction
  failure, or missing workspace fails the outcome rather than dropping data.

## Acceptance Criteria

### F04-AC01 — Cold gate is reproducible

- Traces to `M-08-C05`.
- Observable required verbatim: **both runs report identical command inventory
  and exit classifications; each produces a run ID, SHA, branch, dirty flag,
  tool versions, target types, exit codes, and redacted artifact paths.**
- Proven by the schema/manifest contract suite and two real cold runs from one
  clean fixed candidate implementation SHA equal to `HEAD`. F-04 acceptance
  occurs only after these runs and independent review; F-09 later reruns the
  cold gate from the accepted M-08 SHA.

### F04-AC02 — Cold dependencies and caller state are isolated

- Traces to `M-08-C05`.
- Two fresh run snapshots provision Go/npm dependencies and the PostgreSQL image
  through the secret-free provisioning lane; caller SHA/status, `node_modules`,
  and cache inventories are unchanged.
- Proven by snapshot/provision negative fixtures and the two-run integration
  comparison.

### F04-AC03 — Accepted deterministic lanes remain authoritative

- Traces to `M-08-C02`, `M-08-C03`, and `M-08-C09`.
- The fixed inventory consumes governance, hermetic unit, npm test/build, and
  F-03 ephemeral PostgreSQL lanes without ambient runtime configuration,
  persistent dev mutation, or missing active workspace.
- Proven by cold-gate contract tests and both real outcome manifests.

### F04-AC04 — Evidence is complete, redacted, and honest

- Traces to `M-08-C04` and `M-08-C05`.
- Schema-valid trace/outcome records contain every required field and only
  repository-relative artifact paths; injected secrets, PII, URLs, absolute
  paths, unknown targets, and fake-to-live promotion fail closed.
- Proven by adversarial evidence tests and redaction scan of both manifests.

### F04-AC05 — Aggregate failure preserves safe evidence

- Traces to `M-08-C05`.
- A failed/timed-out lane makes the aggregate nonzero with a stable reason while
  retaining completed independent records and F-03 cleanup evidence.
- Proven by aggregator failure fixtures and final regression execution.

### F04-AC06 — Ephemeral PostgreSQL remains isolated

- Traces to `M-08-C03`.
- The cold inventory consumes F-03 migrations/tests with `--pull=never`, reports
  `32/0`, preserves dev invariance, and leaves zero labelled resources.
- Proven by both real cold manifests and post-run resource inventory.

### F04-AC07 — Governance remains current in the snapshot

- Traces to `M-08-C09`.
- Governance schema/drift executes from the detached candidate-SHA clone and
  records a passing target-labelled result without stale context.
- Proven by both real cold manifests and current-context validation.

## Stop Conditions

- Provisioning needs private credentials, an undeclared proxy/runtime key, or a
  write outside ignored run state/global Docker image provisioning.
- The PostgreSQL tag cannot resolve to a stable recorded identity, F-03 attempts
  a pull during validation, or a Docker/database resource survives.
- The caller checkout/cache changes, a required workspace is absent, evidence
  cannot be redacted, or accepted seam behavior must be rewritten.
- A competing shared-seam writer or architecture/contract contradiction appears.

## Handoff

- Current status: `spec_ready`.
- Next owner: Fresh F-04 Build Feature Implementer.
- Next action: Execute the split TDD plan beginning with schema/manifest RED.
- Required evidence: RED/GREEN contracts, two real cold manifests, caller-state
  invariance, zero resource inventory, and fixed-commit independent review.
- Blockers: None during planning; real provisioning prerequisites are rechecked
  by the build session without requesting Oracle/provider credentials.
