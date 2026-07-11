# F-04 Deterministic Cold Gate and Evidence Manifest — Plan

```yaml
id: F-04
type: feature-plan
status: planned
owner: Feature Implementer
parent: M-08
created: 2026-07-11
updated: 2026-07-11
validation_level: QA-3
lifecycle_scope: feature
split_decision: split
```

## Scope and Ownership

F-04 is an L3 serial change to the harness/evidence and execution-lane
governance seams. The build session may extend the registry with one
non-validation provisioning lane, but must consume F-03 and the accepted
Execution/Environment/Policy modules without rewriting them. Context remains
consume-only except for the minimal `{candidate_sha}` template/aggregate-lane
contract extension required by the staged post-GREEN amendment. Split is
required because the plan has five TDD/real-evidence phases and an independent
fixed-commit review; dispatch a fresh build session.

## Machine Work Contract

```json
{
  "schema_version":"1.0",
  "feature_id":"F-04",
  "required_sources":[
    "scripts/harness.ps1",
    "scripts/harness/Execution.psm1",
    "scripts/harness/Environment.psm1",
    "scripts/harness/Policy.psm1",
    "scripts/harness/Context.psm1",
    "scripts/harness/Postgres.psm1",
    "contracts/governance/execution-lanes.json",
    "package-lock.json",
    "go.work",
    "apps/server_core/go.mod",
    ".gitignore"
  ],
  "allowed_paths":[
    "scripts/harness.ps1",
    "scripts/harness/Evidence.psm1",
    "scripts/harness/Context.psm1",
    "scripts/tests/cold-gate-evidence.tests.ps1",
    "scripts/tests/cold-gate-snapshot.tests.ps1",
    "scripts/tests/cold-gate.integration.tests.ps1",
    "scripts/tests/governance-contracts.tests.ps1",
    "contracts/governance/execution-lanes.json",
    "contracts/governance/schemas/execution-lanes.schema.json",
    "contracts/governance/schemas/feature-work-contract.schema.json",
    "contracts/governance/schemas/context-pack.schema.json",
    "contracts/governance/schemas/harness-outcome.schema.json",
    "package.json",
    ".mnfs/MIS-001-mercado-livre-operating-cockpit/M-08-repository-integrity-harness/F-04-cold-gate-evidence/**"
  ],
  "forbidden_paths":[
    "scripts/harness/Execution.psm1",
    "scripts/harness/Environment.psm1",
    "scripts/harness/Policy.psm1",
    "scripts/harness/Postgres.psm1",
    "package-lock.json",
    "apps/server_core/**",
    "apps/web/**",
    "packages/**",
    "contracts/api/**",
    "docker-compose.yml",
    "docker/**",
    ".env"
  ],
  "side_effects":{"allowed":["repository-write","database-write"],"forbidden":["provider-write"]},
  "commands":[
    {"id":"evidence-contract","command_template":"pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/tests/cold-gate-evidence.tests.ps1","lane_id":"unit","expected_exit_code":0},
    {"id":"snapshot-contract","command_template":"pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/tests/cold-gate-snapshot.tests.ps1","lane_id":"unit","expected_exit_code":0},
    {"id":"cold-regression","command_template":"pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/tests/cold-gate-evidence.tests.ps1; pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/tests/cold-gate-snapshot.tests.ps1","lane_id":"unit","expected_exit_code":0}
    ,{"id":"cold-real-1","command_template":"npm run harness:cold -- --CandidateSha {candidate_sha}","lane_id":"unit","expected_exit_code":0}
    ,{"id":"cold-real-2","command_template":"npm run harness:cold -- --CandidateSha {candidate_sha}","lane_id":"unit","expected_exit_code":0}
  ],
  "criteria":[
    {"id":"F04-AC01","milestone_criterion_id":"M-08-C05","command_ids":["evidence-contract","cold-regression","cold-real-1","cold-real-2"]},
    {"id":"F04-AC02","milestone_criterion_id":"M-08-C05","command_ids":["snapshot-contract","cold-real-1","cold-real-2"]},
    {"id":"F04-AC03","milestone_criterion_id":"M-08-C02","command_ids":["snapshot-contract","cold-regression","cold-real-1"]},
    {"id":"F04-AC04","milestone_criterion_id":"M-08-C04","command_ids":["evidence-contract","cold-regression","cold-real-1","cold-real-2"]},
    {"id":"F04-AC05","milestone_criterion_id":"M-08-C05","command_ids":["evidence-contract","cold-regression"]},
    {"id":"F04-AC06","milestone_criterion_id":"M-08-C03","command_ids":["snapshot-contract","cold-regression","cold-real-1","cold-real-2"]},
    {"id":"F04-AC07","milestone_criterion_id":"M-08-C09","command_ids":["snapshot-contract","cold-regression","cold-real-1","cold-real-2"]}
  ],
  "stop_conditions":[
    {"code":"private-provisioning","condition":"Cold dependency or image provisioning requires credentials, proxy secrets, or undeclared runtime configuration."},
    {"code":"caller-mutation","condition":"The primary checkout, caller node_modules, Go/npm caches, or lockfiles would change."},
    {"code":"image-identity","condition":"postgres:16-bookworm cannot resolve to a stable recorded identity or validation attempts a pull."},
    {"code":"evidence-unsafe","condition":"A secret, PII, URL, absolute path, unknown target, or evidence-class promotion reaches persisted evidence."},
    {"code":"scope-conflict","condition":"An accepted consume-only seam beyond the declared Context template extension needs rewriting or a competing shared-seam writer appears."}
  ],
  "retry_budget":{"max_correction_attempts":1},
  "handoff_fields":["commit-sha","changed-paths","red-evidence","green-evidence","cold-manifests","inventory-comparison","target-classes","caller-invariance","resource-inventory","remaining-risks"]
}
```

The current schema uses coarse `repository-write`/`database-write` tokens. Their
exact authorized scope is only: ignored run/snapshot/npm/Go cache writes, public
dependency-registry network reads, global Docker image-cache write for
`postgres:16-bookworm`, and generated `mpc_test_*` ephemeral PostgreSQL writes.
Forbidden regardless of coarse vocabulary: dev/live database mutation,
Oracle/provider/browser access, provider write, private registry credentials,
primary checkout writes, caller dependency/cache writes, and unowned Docker
resources. Phase 3 extends governance vocabulary atomically so the amended
contract expresses these exact scopes directly before Phase 4.

Snapshot trust correction: source validation must resolve the canonical expected
checkout root and prove that its direct child `.git` is a normal directory.
Gitfiles, reparse-point gitdirs, external/ancestor gitdirs, and non-exact paths
fail closed as `GITDIR_UNTRUSTED`; the child clone receives exactly one
structured `-c safe.directory=<canonical-exact-gitdir>` override.

## Phase 1 — Schema and Manifest RED

1. Add adversarial tests for the outcome schema and JSONL writer: required
   identity/tool/ordered-command fields, stable projection, target/evidence
   compatibility, relative paths, secret/PII/URL/absolute-path rejection, and
   aggregate failure preservation.
2. Extend governance tests first for the exact `cold-provision` lane:
   no inheritance/runtime keys, live dependency-registry network, disabled DB,
   isolated cache-write side effect, `external-dependency-registry` target, and
   `provisioning` evidence class. RED must be missing schema/implementation, not
   syntax or external availability.

Commit RED tests only: `test(harness): define cold evidence contract`.

## Phase 2 — Cold Snapshot and Provision RED

3. Add fake executable/repository fixtures proving clean candidate-SHA preflight,
   ignored detached clone/run layout, unchanged caller state/caches, secret-free
   Environment-derived provisioning, frozen Go/npm/image order, recorded image
   identity, no validation on failed prerequisites, and stable reason codes.

Run only fake/unit suites. No Docker/container/network command is allowed in
this phase. Commit: `test(harness): define cold snapshot contract`.

## Phase 3 — Aggregator GREEN

4. Add `Evidence.psm1`, the outcome schema, registry/schema lane, and the
   `harness.ps1 -Command cold` orchestration. Generate a fresh run environment,
   add only run-local npm/Go cache variables, provision Go/npm/image in order,
   then execute the fixed governance/unit/test/build/F-03 inventory from the
   detached snapshot. Keep provisioning records distinct from validation.
5. Continue safe independent commands after failures, aggregate all results,
   validate/redact before persistence, and add a thin `harness:cold` npm alias.
   Do not change lockfiles or accepted consume-only modules.
6. Atomically add `cold-provision` subprocess classification and distinct
   top-level `cold-gate` aggregate target/classification to execution-lane,
   work-contract, context-pack, and outcome schemas. Subprocess records retain
   `cold-provision`, `fake`, and `ephemeral-postgres`; the mixed top-level run is
   never labelled fake, provisioning, integration, or live. Amend this plan's
   machine contract with two load-bearing commands:

   - `cold-real-1`: `npm run harness:cold -- --CandidateSha {candidate_sha}`
   - `cold-real-2`: `npm run harness:cold -- --CandidateSha {candidate_sha}`

   Map both commands to F04-AC01/02/03/04/06/07 and at least one to F04-AC05.
   Extend Context only enough to bind `{candidate_sha}` to the clean fixed
   candidate `HEAD`. Recompile and current-validate the amended pack before any
   Phase 4 execution; failure blocks real execution.

Run Phase 1/2 suites plus current governance, environment, execution, alias,
and F-03 contract regressions. Commit: `feat(harness): add deterministic cold gate`.

## Phase 4 — Two Real Cold Runs

7. From one clean fixed candidate implementation SHA equal to `HEAD`, inventory
   caller status, dependency/cache paths,
   Docker resources, and current image identity; run the versioned cold command
   twice. This is the only phase authorized for package-registry/image network
   and ephemeral PostgreSQL execution.
8. Require both manifests to have identical ordered command inventory, target/
   evidence pairs, resolved image identity, exit classifications, and aggregate
   classification. Require distinct run IDs, relative artifact paths, unchanged
   caller inventory, F-03 migration `32/0`, and zero run-labelled resources.
   Record safe paths/summaries only; do not fetch live credentials.

## Phase 5 — Independent Final Review and Handoff

9. Run full planned regression on the fixed candidate implementation commit. Dispatch
   independent SPEC/SAFETY and QUALITY review against that same commit.
10. Consolidate findings into at most one correction batch, rerun the complete
   regression and both cold runs once, then create `validation.md`. A new
   material architecture/scope finding blocks or replans instead of opening a
   review loop. Only then may the Milestone Orchestrator accept F-04 and link
   the accepted SHA in evidence. F-09 later reruns the gate on the accepted M-08
   SHA for `M-08-C12`.

## Verification Mapping

- `evidence-contract` proves F04-AC01, F04-AC04, and F04-AC05.
- `snapshot-contract` proves F04-AC02 and the M-08-C02 portion of F04-AC03.
- The Phase 3 atomic contract amendment adds `cold-real-1` and `cold-real-2`
  using exact command `npm run harness:cold -- --CandidateSha {candidate_sha}`;
  these become the load-bearing proofs for F04-AC01 through F04-AC04, F04-AC06,
  and F04-AC07 before Phase 4. The current pre-build contract remains valid
  until the registry/Context vocabulary exists.
- `cold-regression` proves failure aggregation and contract stability for
  F04-AC05. Independent review checks every criterion on the fixed commit.

## Rollback/Risk Notes

- Mutable `postgres:16-bookworm` can resolve differently; record the resolved
  identity and block a same-SHA comparison mismatch rather than normalizing it.
- A local detached clone reads exact committed candidate objects, never mutable
  source-worktree files, and preserves governance semantics without registering
  a worktree or touching the caller checkout.
- Network success proves dependency provisioning only; it is never validation
  or live-provider/Oracle evidence.
- Remove only new ignored run state when cleanup ownership is exact. Never reset,
  stash, clean, restore, or overwrite user state.

## Handoff

- Current status: `planned`.
- Next owner: Fresh F-04 Build Feature Implementer.
- Next action: Start Phase 1 RED; do not execute Docker/network before Phase 4.
- Split decision: `split` — fresh build context required by plan size and L3
  evidence/review scope.
- Review policy: independent final SPEC/SAFETY and QUALITY reviews on one fixed
  commit, one consolidated correction batch maximum.
- Blockers: None at plan time; provisioning stop conditions are owner-reserved.
