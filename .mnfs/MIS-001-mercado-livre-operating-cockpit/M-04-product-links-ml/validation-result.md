# Milestone Validation Result

```yaml
id: M-04
type: milestone-validation-result
status: passed
owner: QA Validator
parent: MIS-001
created: 2026-07-09
updated: 2026-07-09
validation_level: QA-0
lifecycle_scope: milestone
```

## Milestone

M-04-product-links-ml

## Verdict

- Result: `passed`
- Blocking failures: none
- Summary: Marketplace Central now imports real Mercado Livre listing identities, generates exact-first product-link candidates from live Oracle-backed internal product reads, and lets operators approve, reject, or manually resolve link truth with persisted audit evidence through API, SDK, UI, and live backend validation.

## Validation Scope Declaration

- contract_validated: Yes
- integration_validated: Yes
- live_validation: Yes, for listing import, candidate generation, and operator link-resolution persistence
- blocked_for_real_validation: downstream stock-write gating remains intentionally unvalidated here because stock actions belong to M-05

This pass covers the full M-04 scope for product-link truth. It does not claim that a stock action engine already consumes these states; that downstream enforcement belongs to M-05 Stock Seguro.

## Feature Evidence

- F-01 listing snapshot import: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-04-product-links-ml/F-01-listing-snapshot-import/validation.md`
- F-02 link candidate engine: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-04-product-links-ml/F-02-link-candidate-engine/validation.md`
- F-03 link resolution workflow: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-04-product-links-ml/F-03-link-resolution-workflow/validation.md`

## Criterion Review

### M-04-C01 — Exact Match Priority

- Status: `passed`
- Commands:
  - `$env:GOCACHE='C:\Users\leandro.theodoro\Documents\marketplace-central\.gocache'; go test ./apps/server_core/internal/modules/product_links/... ./apps/server_core/internal/composition -count=1`
  - `GET http://localhost:8080/product-links/link-workflows?installation_id=inst-mercado_livre-5cb46653-95c3-46f7-a30f-9fbf5ace5f98&limit=20`
- Expected:
  - exact EAN / seller-SKU candidate generation outranks heuristic states and remains visible as evidence before any operator decision.
- Actual:
  - focused `product_links` backend tests passed
  - live workflow output on the official Docker backend returned exact-EAN candidates for the imported Mercado Livre identities with preserved candidate metadata, internal product ids, and match values
  - operator approval consumed an exact-EAN candidate without mutating the underlying candidate evidence
- Blocking failure observed:
  - `No`

### M-04-C02 — Link State Is Explicit Enough To Block Downstream Writes

- Status: `passed`
- Commands:
  - `npm run test --workspace @marketplace-central/feature-product-links`
  - `POST http://localhost:8090/product-links/link-resolutions/approve-candidate`
  - `POST http://localhost:8090/product-links/link-resolutions/reject-listing`
  - `POST http://localhost:8090/product-links/link-resolutions/manual-resolve`
  - `GET http://localhost:8090/product-links/link-workflows?installation_id=inst-mercado_livre-5cb46653-95c3-46f7-a30f-9fbf5ace5f98&limit=20`
  - direct Postgres inspection of `product_links` and `product_link_audit_entries`
- Expected:
  - unresolved/conflict/rejected/resolved states stay explicit so downstream stock-action milestones can safely block unsafe write paths instead of inferring or guessing link truth.
- Actual:
  - product-links UI tests passed for unresolved, conflict, resolved, rejected, approve, reject, and manual resolve states
  - live approve/reject/manual-resolve calls persisted distinct `resolved` and `rejected` truth rows plus matching audit rows
  - workflow responses show untouched items still carrying only candidate evidence while acted-on items expose explicit `current_link` state and audit history
  - direct database proof showed:
    - `MLB4807275656` -> `resolved` with `internal_product_id=20307`
    - `MLB4834408384` -> `rejected`
    - `MLB4834419602` -> `resolved` with `internal_product_id=20312`
- Blocking failure observed:
  - `No`

## Validation Notes

- The first live workflow validation attempt failed because the Docker Postgres used by the running stack had not yet received `0025_product_link_workflows.sql`; after applying the migration, the live endpoints and persistence behaved correctly.
- Official Docker services were rebuilt and revalidated so the evidence does not depend only on the temporary host-side validation server.
- M-04 stops at link truth and audit. The first milestone that must prove a stock action engine actually refuses unsafe writes is M-05.

## Live Product-Link Validation

Date: 2026-07-09
Environment: local Docker backend/frontend plus active Mercado Livre development credentials
Installation: `inst-mercado_livre-5cb46653-95c3-46f7-a30f-9fbf5ace5f98`

Evidence:
- `GET /product-links/link-workflows` returned five live Mercado Livre listing identities with exact-EAN candidate evidence.
- `POST /product-links/link-resolutions/approve-candidate` resolved `MLB4807275656` to internal product `20307`.
- `POST /product-links/link-resolutions/reject-listing` persisted `MLB4834408384` as `rejected`.
- `POST /product-links/link-resolutions/manual-resolve` resolved `MLB4834419602` to internal product `20312`.
- Direct Postgres inspection confirmed persisted `product_links` rows and `product_link_audit_entries` rows with action, previous/next state, actor, and reason.
- After `docker compose up --build -d backend frontend`, the official `http://localhost:8080/product-links/link-workflows?...` endpoint returned the same persisted workflow truth from the running stack.

Boundary:
- No provider stock/price/listing writes were executed.
- Product-link live validation covers linking truth and audit only.
- Downstream write blocking remains the responsibility of M-05 consumers of these explicit states.

## Handoff

- Milestone status: `ready for mission continuation`
- Next recommended action: start M-05 Stock Seguro and consume `product_links` states as hard safety gates for manual stock actions
- Open blockers: none for M-04
