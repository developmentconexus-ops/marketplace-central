# F-01-aggregate-sync-endpoints

```yaml
id: F-01
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-05
created: 2026-07-14
updated: 2026-07-14
validation_level: QA-0
lifecycle_scope: feature
```

## Mission

MIS-003. Binding contracts: IC-02 (summary reuse), R-03 (orders tables 0027/0033, `integration_operation_runs`), `GOV_API_SDK_SPLIT`. ADR-17 (unknown counters null).

## Milestone

M-05 visao-geral-pedidos-sync-central.

## Brief

Three small read APIs. **Dashboard summary** `GET /dashboard/summary?installation_id=`: composes existing module reads — listings summary (IC-02), pendências counts (sync errors, unresolved links, below_margin, sem GTIN), orders today/7d counts, last sync times per module — each counter nullable; a failed sub-read yields null for that counter plus a `degraded: [source]` list, never a 500 or a 0. **Orders read** `GET /orders` (cursor, filters: status, fulfillment, date range, q) + `GET /orders/{orderId}` over existing orders tables — canonical fields only, provider raw behind adapter. **Sync runs** `GET /sync/runs` (cursor, filters: module, status; sort `started_at DESC`, newest first) over `integration_operation_runs`, window fixed at last 90 days (`started_at >= now() - 90d`; older runs excluded from the listing — mission §Audit & history "sync central 90d view"). OpenAPI + sdk-runtime same commit. This feature owns the new `/dashboard` and `/sync` server prefixes and adds their two Vite dev-proxy rows to `apps/web/vite.config.ts` (per IC-05 writer sequence — M-02's five rows already merged).

EARS:
- While one sub-source fails, when summary is computed, the response shall carry null for that counter and name the source in `degraded[]` (200, partial-honest).
- While orders exist, when listed with `filter.status=`, the API shall return canonical order rows (id, provider_order_id, buyer nickname, total, currency, status, fulfillment, NF state nullable, created_at) cursor-paginated newest first.
- While a run is in progress, when runs are listed, its row shall show `status=running` with null finished_at.
- While date filters are malformed, when listing orders, the API shall return 400 `invalid_filter`.

## Inputs

- R-03: orders migrations 0027/0033 column facts, `integration_operation_runs` shape, listings summary op (IC-02), product_links pending counts source, existing cursor precedent (M-01 F-02).

## Expected Output

- Three endpoint groups in their owning modules (summary in a thin `dashboard` transport composing module application services — no cross-module SQL joins; orders in orders module; runs in integrations module).
- OpenAPI + sdk-runtime (`getDashboardSummary`, `listOrders`, `getOrder`, `listSyncRuns`) same commit; `/orders` proxy row exists from M-02.
- Integration tests: summary composition incl. degraded path (one source stubbed to fail), orders filters + cursor, runs listing.

## Constraints

- Read-only; no new tables, no writes, no order mutation endpoints (Non-Scope: faturar).
- Module boundaries: summary calls application services, never other modules' tables directly.
- Unknown counter ≠ 0 — JSON null mandatory; tests assert null not 0 on degraded.
- Tenant scoping on every query.

## Negative Scenarios

- Summary with unknown installation → 404 `installation_not_found`.
- `GET /orders/{unknown}` → 404 `order_not_found`.
- Malformed cursor → 400 `invalid_cursor`.
- All sub-sources fail → 200 with all-null counters + full degraded list (page still renders).

## Validation Expectations

- `go test` output: degraded-path test asserting `"pending_links": null` (not 0) + `degraded: ["linkage"]`; filter and cursor tests.
- Integration transcript: summary JSON with real seeded counts; orders page walk.
- Diff proof: openapi.yaml + sdk-runtime same commit.

## Execution Artifact Rules

`spec.md`, `plan.md`, `validation.md` created during feature execution.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer.
- Next action: compile context pack; read R-03 orders/runs sections + IC-02 summary + cursor precedent.
- Required files/evidence: `validation.md` in this folder.
- Blockers or open decisions: none.
