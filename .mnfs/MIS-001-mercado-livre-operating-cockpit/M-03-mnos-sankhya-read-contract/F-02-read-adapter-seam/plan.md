# Feature Plan

```yaml
id: F-02
type: feature-plan
status: planned
owner: Feature Implementer
parent: F-02
created: 2026-07-07
updated: 2026-07-07
validation_level: QA-0
lifecycle_scope: feature
```

## Steps

1. Add failing tests for service delegation and Oracle config secrecy.
2. Implement `application.Service` over `ports.Reader`.
3. Implement fake adapter fixtures and methods for all six operations.
4. Implement guarded Oracle config and read-only adapter seam.
5. Run focused package tests and record evidence.

## Files Expected To Change

- `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-03-mnos-sankhya-read-contract/F-02-read-adapter-seam/spec.md`
- `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-03-mnos-sankhya-read-contract/F-02-read-adapter-seam/plan.md`
- `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-03-mnos-sankhya-read-contract/F-02-read-adapter-seam/validation.md`
- `apps/server_core/internal/modules/internal_read/application/service.go`
- `apps/server_core/internal/modules/internal_read/application/service_test.go`
- `apps/server_core/internal/modules/internal_read/adapters/fake/reader.go`
- `apps/server_core/internal/modules/internal_read/adapters/oracle/config.go`
- `apps/server_core/internal/modules/internal_read/adapters/oracle/config_test.go`
- `apps/server_core/internal/modules/internal_read/adapters/oracle/reader.go`
