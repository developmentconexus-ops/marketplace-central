# Feature Spec

```yaml
id: F-04
type: feature-spec
status: spec_ready
owner: Feature Implementer
parent: F-04
created: 2026-07-09
updated: 2026-07-09
validation_level: QA-0
lifecycle_scope: feature
```

## Feature ID

F-04-stock-seguro-dashboard

## Problem

M-05 already has the inventory policy, risk classifier, and manual action safety rules, but operators still cannot see linked listing stock risk or apply a confirmed audited correction from the product. The missing slice is the vertical contract: inventory API, OpenAPI, SDK runtime surface, action persistence, and an operator-facing dashboard that renders only backend-calculated truth.

## Requirements

- Expose a stock risk list API under `inventory` using persisted data only.
- The stock risk list must not trigger synchronous provider probes to render rows.
- Read risk evidence from existing listing snapshots, product link truth, internal sellable stock, and the inventory policy model.
- Expose a manual stock action API that reuses the inventory action service and persists audit evidence.
- Persist stock action audit rows in Postgres with before quantity, requested quantity, operator, trigger, policy id, source timestamps, provider response summary, idempotency key, and failure/blocking reason.
- Add SDK runtime types and client methods for the inventory risk/action endpoints.
- Add an operator route at `/inventory/stock-seguro`.
- The dashboard must support filtering by risk state, link state, and actionability.
- Row detail must show internal stock, Mercado Livre stock, recommendation, policy id, freshness/source timestamps, and blocking reason when present.
- Manual apply must be explicitly confirmatory and must refetch the affected row list after completion.
- React must not calculate stock recommendation or safety math.

## Non-Goals

- No automatic stock correction, scheduler, or bulk apply.
- No runtime capability registration that exposes `stock_write` as generally available in this slice.
- No synchronous provider stock probe as the dashboard source of truth.
- No local browser storage as inventory truth.
- No fake-only validation claim for live inventory writes.

## Design

Create the first `inventory` vertical from transport to UI:

1. `inventory/application` gets a risk listing service that composes:
   - persisted product link listing snapshots,
   - persisted product link workflow truth,
   - `internal_read` sellable stock reads mapped from the inventory policy,
   - inventory domain classification.
2. `inventory/application` gets a manual action facade that:
   - resolves installation/provider account context from `integrations`,
   - maps inventory-owned stock write requests to the Mercado Livre stock writer adapter,
   - persists action audit rows through an inventory repository,
   - returns an operator-facing action result.
3. `inventory/adapters/postgres` persists manual action rows and reads them back by `stock_action_id`.
4. `inventory/transport` exposes:
   - `GET /inventory/stock-risks`
   - `POST /inventory/stock-actions/manual-apply`
5. `packages/sdk-runtime` adds typed inventory clients.
6. `packages/feature-inventory` renders the dashboard and action panel.
7. `apps/web` registers the route and sidebar entry.

The dashboard route is `/inventory/stock-seguro`. Filters are query-param driven so operators can share or refresh the same cockpit state.

## API Shape

- `GET /inventory/stock-risks?installation_id=...&state=...&link_state=...&actionability=...&limit=...`
  - returns `{ items: StockRiskListItem[] }`
- `POST /inventory/stock-actions/manual-apply`
  - accepts listing identity, operator approval metadata, optional reason, and `stock_action_id`
  - returns `{ action: StockActionRecord, risk: StockRiskListItem }`

`StockRiskListItem` carries:

- listing identity and provider title/reference
- current link state and resolved internal product reference when present
- provider quantity and provider source timestamps
- internal quantity and internal source timestamps
- recommended quantity and policy id
- risk state plus blocking reason
- `actionable` boolean computed by backend
- last action summary when an action already exists for that listing/action id

## Edge Cases

- Unresolved, conflict, and rejected links stay visible as blocked rows.
- Missing provider quantity is rendered as unsupported, not zero.
- Missing internal stock is rendered as stale, not zero.
- If Oracle/internal read is unavailable, the API returns a structured inventory error rather than optimistic quantities.
- If the provider write is rejected or fails, the UI shows the persisted failed result and the refetched row state.
- Repeating the same `stock_action_id` returns the stored action instead of producing a duplicate write.

## Acceptance Criteria

- M-05-C03: Manual action API persists complete audit evidence in Postgres.
- M-05-C04: Frontend tests cover loading, error, empty, healthy, oversell, undersell, stale, unresolved, conflict, and action result states.
- M-05-C04: Browser validation shows the dashboard route at `/inventory/stock-seguro` with readable desktop and mobile layout.
- M-05 done-mean: stock risk rows render from persisted read sources, not synchronous provider probes.

## Handoff

- Current status: `spec_ready`
- Next owner: Feature Implementer
- Next action: implement API, SDK, UI, and browser validation
- Required files/evidence: `plan.md`, tests, `validation.md`
- Blockers or open decisions: live write validation still requires explicit operator approval during QA
