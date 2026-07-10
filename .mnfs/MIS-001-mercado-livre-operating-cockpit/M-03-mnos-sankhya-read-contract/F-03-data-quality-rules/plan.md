# Feature Plan

```yaml
id: F-03
type: feature-plan
status: planned
owner: Feature Implementer
parent: F-03
created: 2026-07-07
updated: 2026-07-07
validation_level: QA-0
lifecycle_scope: feature
```

## Steps

1. Add or rewrite quality-state stability tests around the Oracle-first contract.
2. Extend fake-reader and Oracle-adapter tests so both respect nil-preserving and explicit-quality semantics.
3. Update module/wiki references that still imply fake-only validation is enough.
4. Run focused tests and document the exact boundary between local proof and real Oracle proof.

## Verification Commands

- Command: `cd apps/server_core; $env:GOCACHE=(Join-Path (Get-Location) '.gocache'); go test ./internal/modules/internal_read/...`
  - Satisfies criterion ID: M-03-C02
  - Expected result: Pass. Quality-state and nil-preserving tests remain green across the internal-read surface.
- Command: `rg -n "mock|fake|seam" .mnfs/MIS-001-mercado-livre-operating-cockpit/M-03-mnos-sankhya-read-contract`
  - Satisfies criterion ID: M-03-C03
  - Expected result: Pass with honest wording. Artifact language must distinguish local seam proof from live Oracle proof.

## Rollback/Risk Notes

- Risk: quality states drift between fake tests and Oracle adapter behavior.
- Recovery: keep the quality-state contract centralized in domain tests and reuse the same assertions across adapters where possible.

## Handoff

- Current status: `planned`
- Next owner: Feature Implementer
- Next action: implement and record fresh quality-state evidence
- Required files/evidence: spec, verification commands, QA notes, updated validation artifact
- Blockers or open decisions: none beyond the rewritten Oracle-first execution scope
