# Feature Spec

```yaml
id: F-01
type: feature-spec
status: spec_ready
owner: Feature Implementer
parent: F-01
created: 2026-07-08
updated: 2026-07-08
validation_level: QA-0
lifecycle_scope: feature
```

## Feature ID

F-01-listing-snapshot-import

## Problem

`product_links` cannot generate safe candidates until Mercado Livre listings exist as tenant-scoped, normalized, re-importable internal records instead of transient probe output.

## Requirements

- Create a dedicated `product_links` persistence path for normalized listing snapshots.
- Reuse the existing installation-aware Mercado Livre listing read surface instead of talking to provider adapters directly.
- Persist one row per linkable unit: item-level row when no variations exist, variation-level row when variations exist.
- Keep snapshots tenant-scoped and idempotent by provider item + variation identity.
- Expose a manual import trigger so live validation can happen before any scheduler exists.

## Non-Goals

- Candidate generation or confidence scoring.
- Manual approval/rejection workflow.
- Any Mercado Livre write action.

## Acceptance Criteria

- Tests cover listing import without variations and with variations.
- Re-import of the same listing set does not duplicate persisted rows.
- OpenAPI and `sdk-runtime` expose the manual import contract.
- Live validation proves the import works against a connected Mercado Livre installation.

## Handoff

- Current status: `spec_ready`
- Next owner: Feature Implementer
- Next action: execute the import slice and record real validation evidence
- Required files/evidence: plan, validation, API contract update, persisted snapshot evidence
- Blockers or open decisions: scheduler intentionally deferred
