# F-10 Quick Validation

```yaml
feature_status: quick_validation_passed
evidence_target: fake/current-checkout
base_sha: 75299ea7fcf708ae10aac95efaaaaf7bcbbfc74e
```

## Criterion evidence

| Criterion | Evidence type | Result |
| --- | --- | --- |
| F10-AC01 / M-08-C17 | ran | Active-surface search returned no matches outside historical MNFS/design artifacts; current governance contract passed. |
| F10-AC02 / M-08-C15 | ran | `harness-impact.tests.ps1` proved exactly `first|second` registered IDs run in order and an unregistered ID blocks. |
| F10-AC03 / M-08-C15 | ran | The fixture proved stale-base, unknown-target, and out-of-scope packs block with `CTX_BASE_SHA_MISMATCH`, `CTX_SCHEMA_INVALID`, and `HIMPACT_PATH_OUTSIDE_SCOPE`. |
| F10-AC04 / M-08-C15 | ran | Outcome fixture proved schema-valid base/commit SHAs, relative changed paths, risk, ordered records, and aggregate outcome. |

## Commands run

| Command | Target | Evidence type | Result |
| --- | --- | --- | --- |
| `pwsh scripts/tests/harness-impact.tests.ps1` | fake | ran | Pass |
| `pwsh scripts/tests/governance-contracts.tests.ps1` | fake | ran | Pass |
| `pwsh scripts/tests/context-compiler.tests.ps1` | fake | ran | Pass |
| `pwsh scripts/tests/harness-aliases.tests.ps1` | fake | ran | Pass |
| `pwsh scripts/tests/harness-execution.tests.ps1` | fake | ran | Pass |
| `pwsh scripts/tests/harness-environment.tests.ps1` | fake | ran | Pass |
| `pwsh scripts/harness.ps1 -Command governance -BaseSha 75299ea7fcf708ae10aac95efaaaaf7bcbbfc74e` | fake | ran | Pass |
| `pwsh scripts/harness.ps1 -Command impact -ContextPath scripts/.runs/f10-impact-context.json` | fake/current-checkout | ran | Pass; outcome `scripts/.runs/61b60174b10644a78f86419f2f4650ac/outcome.json`. |

## Same-session fixup

The first current-checkout impact run recorded `HIMPACT_COMMAND_FAILED` because
the fixed-argv child environment omitted its allowed unit `PATH`. The impact
module now uses the existing `New-HarnessChildEnvironment` unit lane and the
full impact plan was rerun. A second failure showed that the alias test is a
separate dispatcher validation rather than a required impact-gate command; it
was removed from the F-10 pack and run separately above. No network, Docker,
browser, database, Oracle, or provider target was used.

## Changed paths

- `package.json`
- `contracts/governance/**`
- `scripts/harness.ps1`, `scripts/harness/{Context,Evidence,Impact}.psm1`
- `scripts/tests/**`, including retirement of the three cold-gate tests
- F-10 `spec.md`, `plan.md`, and this validation artifact

## Risks and handoff

The gate reports `dependency-graph` as the declared seam because `package.json`
changes. Evidence is fake/current-checkout only and does not assert any live
integration.

## Milestone acceptance review

- Decision: Accepted.
- Reviewer: Milestone Orchestrator.
- Review basis: one L2 fixed-commit review found a missing unknown-target
  fixture; the single consolidated correction added it and the focused impact
  suite passed.
- Scope: spec, plan, changed paths, governance/context contracts, gate outcome,
  and quick-validation evidence match F-10. No product or live-target scope
  entered the feature.
- Next owner: Milestone Orchestrator for serial F-05 dispatch.
