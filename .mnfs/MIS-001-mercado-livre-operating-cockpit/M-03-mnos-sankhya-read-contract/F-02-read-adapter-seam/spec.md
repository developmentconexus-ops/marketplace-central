# Feature Spec

```yaml
id: F-02
type: feature-spec
status: spec_ready
owner: Feature Implementer
parent: F-02
created: 2026-07-07
updated: 2026-07-07
validation_level: QA-0
lifecycle_scope: feature
```

## Feature ID

F-02-read-adapter-seam

## Problem

F-01 defined the contract surface, but MPC still needs an application-facing service, a deterministic fake adapter for tests, and a real Oracle seam that stays read-only and secret-safe.

## Requirements

- Add a thin application service over `ports.Reader`.
- Add a fake adapter that future module tests can consume without live Oracle access.
- Add a minimal Oracle seam that is read-only and returns structured unavailability rather than inventing writes.
- Keep secrets out of config/bootstrap errors.

## Non-Goals

- Live Oracle query implementation.
- HTTP routes.
- Product-link, inventory, or profitability business workflows.

## Acceptance Criteria

- `go test ./internal/modules/internal_read/...` passes with service and adapters compiling.
- Secret-safe config test passes without leaking secret values.
- No route registration or write-path logic is introduced.
