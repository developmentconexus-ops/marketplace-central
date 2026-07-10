# M-06 F-01 Correction Review Package

## Requirements

- `.superpowers/sdd/m06-f01-correction-brief.md`

## Implementer Report

- `.superpowers/sdd/m06-f01-correction-report.md`

## Current-State Files To Review

- `apps/server_core/internal/modules/orders/application/import_service.go`
- `apps/server_core/internal/modules/orders/application/import_service_test.go`
- `apps/server_core/internal/modules/orders/adapters/postgres/order_repo.go`
- `apps/server_core/internal/modules/orders/adapters/postgres/order_repo_test.go`

The orders module is untracked in a shared dirty main worktree, so a meaningful commit-range diff does not exist yet. Review the current files above as the complete correction state and compare them to the brief and report. The relevant changes are the atomic conditional upsert, child replacement guarded by the winning write, safe provider-reference normalization, cancelled-order coverage, and focused real-PostgreSQL tests.

## Binding Constraints

- Newer and equal provider timestamps replace the snapshot.
- Older snapshots are skipped.
- Unknown incoming freshness cannot erase known stored freshness.
- Unknown stored freshness may be replaced.
- Raw provider data is reduced to an operationally safe reference; buyer PII is never persisted or returned.
- Tenant scope and transactional child replacement must remain intact.
- Real-DB evidence must come only from the PostgreSQL command recorded in the report; do not infer live Mercado Livre validation from it.
