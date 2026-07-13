# F-08 Assisted Sankhya Linkage Discovery Plan

```yaml
id: F-08
type: feature-plan
status: planned
owner: Feature Implementer
parent: F-08
created: 2026-07-13
updated: 2026-07-13
validation_level: QA-1
lifecycle_scope: feature
split_decision: single
```

## Feature ID

F-08-assisted-sankhya-linkage-discovery

## Split Decision

`single` — although the documentation Feature produces more than four
artifacts and contains explicitly tolerated operational unknowns, all outputs
share one read-only contract decision and one owned directory; a fresh build
session would add no safer implementation boundary.

## Owner and Boundaries

- Owner: Feature Implementer for this Feature directory only.
- Domain/application boundary: the next worker will place linkage invariants
  in one orders-owned application service, behind repository and Sankhya read
  ports; this Feature only specifies that contract.
- Consumers: orders operator workflow and profitability's exact Oracle source
  identity boundary.
- Legacy decision: no partner/date/product/value matching and no positional
  line number as durable proof; candidate filters remain read-only hints.
- Explicit unknowns: deployed custom-field names, supported uniqueness/index
  mechanism, TOP copy behavior, sanctioned entry surface, live TOP lineage,
  and cancellation/devolution policy remain unknown unless bounded evidence
  proves them.

## Steps

1. Compile and validate the routed Feature context against the accepted base
   SHA, then read only its selectors.
2. Resolve the registered proof commands and inspect the approved live-Oracle
   lane without exposing runtime values.
3. Run only a safely supported, bounded SELECT discovery command; otherwise
   record why it could not run and retain operational facts as unknown.
4. Write `discovery.md`, `sankhya-admin-spec.md`, and
   `implementation-contract.md` with fail-closed semantics.
5. Run scoped content, traceability, diff, and registered focused validation;
   write `validation.md`, then commit only the owned Feature directory.

## Files Expected To Change

- `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-06-orders-margin-ml/F-08-assisted-sankhya-linkage-discovery/spec.md`
  - Reason: acceptance and design contract.
- `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-06-orders-margin-ml/F-08-assisted-sankhya-linkage-discovery/plan.md`
  - Reason: bounded execution and verification map.
- `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-06-orders-margin-ml/F-08-assisted-sankhya-linkage-discovery/context.json`
  - Reason: compiled and validated minimal context evidence.
- `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-06-orders-margin-ml/F-08-assisted-sankhya-linkage-discovery/discovery.md`
  - Reason: repository/runtime fact, inference, and unknown classification.
- `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-06-orders-margin-ml/F-08-assisted-sankhya-linkage-discovery/sankhya-admin-spec.md`
  - Reason: deployable fail-closed administrator contract.
- `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-06-orders-margin-ml/F-08-assisted-sankhya-linkage-discovery/implementation-contract.md`
  - Reason: exact next-worker boundaries and seams.
- `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-06-orders-margin-ml/F-08-assisted-sankhya-linkage-discovery/validation.md`
  - Reason: quick-validation evidence and limitations.

## Verification Commands

- Command: `pwsh -NoProfile -File scripts/harness.ps1 -Command context-compile -FeaturePath <feature-directory> -BaseSha 6fe6e2a056c0397d7c5ad45555581ab1175c7cef -AllowedPath <feature-directory>`; copy the values-safe generated pack to `context.json`; then run `pwsh -NoProfile -File scripts/harness.ps1 -Command context-validate -ContextPath <context.json> -RequireCurrentBase`.
  - Satisfies criterion ID: F08-AC01, F08-AC05.
  - Expected result: context compiles from named routes and validates against
    accepted base SHA without broad context expansion.
- Command ID `live-oracle-linkage-discovery`: run only an existing governed
  command definition that performs bounded aggregate/metadata SELECTs for the
  configured header field and TOP 313-to-306 lineage; if no such definition is
  registered in repository truth, record `could-not-run` and do not substitute
  ad hoc SQL or the generic smoke test.
  - Satisfies criterion ID: F08-AC01, F08-AC03, F08-AC05.
  - Expected result: bounded metadata/aggregate lineage observables with target
    `live-oracle`, or `could-not-run` with operational facts left unknown.
- Command ID `go-tax-provenance-focused`: set `GOCACHE=.gocache`; run
  `go test ./internal/modules/internal_read/adapters/oracle ./internal/modules/internal_read/adapters/fake` from `apps/server_core`.
  - Satisfies criterion ID: F08-AC03, F08-AC04.
  - Expected result: focused contract tests preserve exact identity and
    unknown-tax behavior.
- Command ID `git-diff-check`: run `git diff --check` plus scoped
  `git status --short`.
  - Satisfies criterion ID: F08-AC01, F08-AC02, F08-AC03, F08-AC04, F08-AC05.
  - Expected result: no whitespace errors and no Feature-authored changes
    outside the owned directory.
- Command: bounded artifact scan for credential-key assignments, common PII
  labels, write/DDL verbs presented as executed actions, and invented deployed
  field claims.
  - Satisfies criterion ID: F08-AC02, F08-AC05.
  - Expected result: no secret values, PII output, prohibited write evidence,
    or unproved live field name.

## QA Steps

- Step: Inspect the fact/inference/unknown tables across all three evidence
  artifacts.
  - Expected result: only repository/runtime observations are facts and all
    unproved operational details remain explicitly unknown.
- Step: Trace candidate selection through confirmation, idempotent persistence,
  exact 313 line identity, `TGFVAR`, and 306 descendants.
  - Expected result: candidate attributes never become proof; ambiguous and
    incomplete lineage remains fail-closed/missing.
- Step: Trace each F08 acceptance criterion to the milestone criteria and at
  least one exact command or QA step.
  - Expected result: complete acceptance mapping without a milestone verdict.

## Rollback/Risk Notes

Only Feature documentation is written. Rollback is omission of the commit by
the Milestone Orchestrator; no runtime or operational data needs reversal.
Live access, if unavailable or unsafe, is not a failure of the deployable spec:
facts remain unknown and the implementation contract gates activation. A
repository/runtime architecture or ownership conflict stops the Feature.

## Machine Work Contract

```json
{
  "schema_version": "1.0",
  "feature_id": "F-08",
  "required_sources": [
    ".mnfs/MIS-001-mercado-livre-operating-cockpit/M-06-orders-margin-ml/validation-contract.md",
    ".mnfs/MIS-001-mercado-livre-operating-cockpit/M-06-orders-margin-ml/research/order-sankhya-linkage-architecture-2026-07-12.md",
    ".mnfs/MIS-001-mercado-livre-operating-cockpit/M-06-orders-margin-ml/F-05-docker-live-oracle-runner/validation.md",
    ".mnfs/MIS-001-mercado-livre-operating-cockpit/M-06-orders-margin-ml/F-07-order-specific-tax-provenance/validation.md",
    "contracts/governance/execution-lanes.json"
  ],
  "knowledge_route_ids": ["portfolio-core", "orders-margin"],
  "allowed_paths": [
    ".mnfs/MIS-001-mercado-livre-operating-cockpit/M-06-orders-margin-ml/F-08-assisted-sankhya-linkage-discovery/**"
  ],
  "forbidden_paths": ["apps/**", "contracts/api/**", "packages/**", "docs/research/**", "contracts/governance/**"],
  "side_effects": {
    "allowed": ["repository-write", "isolated-cache-write"],
    "forbidden": ["database-mutation", "provider-write"]
  },
  "commands": [
    {"id": "live-oracle-linkage-discovery", "command_id": "live-oracle-linkage-discovery", "lane_id": "live-oracle", "expected_exit_code": 0},
    {"id": "go-tax-provenance-focused", "command_id": "go-tax-provenance-focused", "lane_id": "unit", "expected_exit_code": 0},
    {"id": "git-diff-check", "command_id": "git-diff-check", "lane_id": "unit", "expected_exit_code": 0}
  ],
  "criteria": [
    {"id": "F08-AC01", "milestone_criterion_id": "M-06-C02", "command_ids": ["live-oracle-linkage-discovery", "git-diff-check"]},
    {"id": "F08-AC02", "milestone_criterion_id": "M-06-C02", "command_ids": ["git-diff-check"]},
    {"id": "F08-AC03", "milestone_criterion_id": "M-06-C02", "command_ids": ["live-oracle-linkage-discovery", "go-tax-provenance-focused"]},
    {"id": "F08-AC04", "milestone_criterion_id": "M-06-C01", "command_ids": ["go-tax-provenance-focused", "git-diff-check"]},
    {"id": "F08-AC05", "milestone_criterion_id": "M-06-C02", "command_ids": ["live-oracle-linkage-discovery", "git-diff-check"]}
  ],
  "stop_conditions": [
    {"code": "write-attempt", "condition": "Any discovery command would perform DDL or an Oracle, provider, or application write."},
    {"code": "truth-conflict", "condition": "Repository and runtime truth conflict on architecture, contract, ownership, or verification."},
    {"code": "unsafe-output", "condition": "Any output would expose a secret, credential, buyer PII, or unbounded operational rows."}
  ],
  "retry_budget": {"max_correction_attempts": 1},
  "handoff_fields": ["status", "commit", "changed-paths", "evidence", "blockers", "next"]
}
```

## Handoff

- Current status: `planned`
- Next owner: Feature Implementer
- Next action: Compile context, perform bounded discovery if safely supported,
  write evidence artifacts, and quick-validate.
- Required files/evidence: spec, context, registered command resolution,
  discovery/admin/implementation artifacts, validation
- Blockers or open decisions: No blocking design decision; explicitly unknown
  operational facts are fail-closed activation prerequisites.
