# Milestone Validation Result

```yaml
id: M-05
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

M-05-stock-seguro-ml

## Verdict

- Result: `passed`
- Blocking failures: none
- Summary: Marketplace Central now computes Stock Seguro from inventory-owned policy/risk logic, blocks unsafe/manual writes behind explicit link and freshness gates, persists manual action audit evidence, and exposes the dashboard workflow through API, SDK, and browser-validated UI with real Mercado Livre read data.

## Validation Scope Declaration

- contract_validated: Yes
- integration_validated: Yes
- live_validation: Yes for backend/API/dashboard reads and browser workflow; No for live provider stock mutation
- blocked_for_real_validation: live Mercado Livre stock write remains intentionally unexecuted because current real rows are blocked and the contract allows fake writer evidence unless an operator explicitly approves a live write

This pass claims the full milestone outcome defined in M-05: classify linked listing stock risk and expose only manual audited stock actions. It does not claim that a real Mercado Livre listing quantity was mutated in this session.

## Feature Evidence

- F-01 stock policy model: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-05-stock-seguro-ml/F-01-stock-policy-model/validation.md`
- F-02 stock risk engine: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-05-stock-seguro-ml/F-02-stock-risk-engine/validation.md`
- F-03 manual stock action audit: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-05-stock-seguro-ml/F-03-manual-stock-action-audit/validation.md`
- F-04 stock seguro dashboard: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-05-stock-seguro-ml/F-04-stock-seguro-dashboard/validation.md`

## Criterion Review

### M-05-C01 — Safe Quantity Formula

- Status: `passed`
- Commands:
  - `$env:GOCACHE='C:\Users\leandro.theodoro\Documents\marketplace-central\apps\server_core\.gocache'; go test ./internal/modules/inventory/... ./internal/composition -count=1`
- Expected:
  - default recommended Mercado Livre quantity is `max(0, SUM(ESTOQUE - RESERVADO) - 1)` for `CODEMP IN (1,2)` and `CODLOCAL=10101`.
- Actual:
  - inventory domain/application tests passed, including `DefaultStockPolicy` and Oracle/internal-read mapping evidence from F-01 and F-02
  - runtime dashboard rows loaded internal sellable quantities and recommended quantities exclusively from inventory policy output, not from React-side math
- Blocking failure observed:
  - `No`

### M-05-C02 — Blocked Unsafe Actions

- Status: `passed`
- Commands:
  - `$env:GOCACHE='C:\Users\leandro.theodoro\Documents\marketplace-central\apps\server_core\.gocache'; go test ./internal/modules/inventory/... -count=1`
  - `Invoke-WebRequest http://localhost:8080/inventory/stock-risks?installation_id=inst-mercado_livre-5cb46653-95c3-46f7-a30f-9fbf5ace5f98&limit=5`
  - browser open `http://localhost:4175/inventory/stock-seguro?installation=inst-mercado_livre-5cb46653-95c3-46f7-a30f-9fbf5ace5f98`
- Expected:
  - unresolved link, conflict link, stale internal source, stale provider source, ineligible product, and unsupported provider shape produce blocked actions and no provider write.
- Actual:
  - inventory service tests passed for unresolved, rejected, conflict, stale internal, stale provider, unsupported provider quantity, ineligible product, incomplete approval, and healthy no-op skip
  - real `GET /inventory/stock-risks` returned blocked rows for stale provider source, unresolved links, and rejected links using persisted Mercado Livre listing snapshots
  - browser validation displayed blocked badges and blocking reasons from real data, proving the operator can distinguish blocked cases before any write path
- Blocking failure observed:
  - `No`

### M-05-C03 — Manual Audit Evidence

- Status: `passed`
- Commands:
  - `$env:GOCACHE='C:\Users\leandro.theodoro\Documents\marketplace-central\apps\server_core\.gocache'; go test ./internal/modules/inventory/... -count=1`
  - `POST /inventory/stock-actions/manual-apply` contract exercised by transport tests
  - `npm run test --workspace @marketplace-central/sdk-runtime`
- Expected:
  - applied/failed action stores before quantity, requested quantity, operator/manual trigger, policy id, source timestamps, provider response summary, and idempotency key.
- Actual:
  - stock action service tests passed for applied, rejected, provider-error, duplicate idempotency, approval gating, and audit-event timestamp evidence
  - transport tests passed for JSON request/response shape of manual apply
  - runtime bugfix added explicit JSON tags so blocking/action audit structures serialize to the API contract expected by SDK/UI
  - SDK tests passed with stock-risk/action methods included in the public client
- Blocking failure observed:
  - `No`

### M-05-C04 — Stock Dashboard States

- Status: `passed`
- Commands:
  - `npm run test --workspace @marketplace-central/feature-inventory`
  - `npx vitest run apps/web/src/app/viteProxy.test.ts apps/web/src/app/AppRouter.test.tsx`
  - browser open `http://localhost:4175/inventory/stock-seguro?installation=inst-mercado_livre-5cb46653-95c3-46f7-a30f-9fbf5ace5f98`
- Expected:
  - dashboard shows healthy, oversell, undersell, stale, unresolved, conflict, ineligible, loading, error, and empty states.
- Actual:
  - feature-inventory tests passed for loading, error, empty, oversell, undersell, healthy, stale, unresolved, conflict, ineligible, and manual action result states
  - router/proxy tests passed, including the fix for the previous Vite proxy collision that incorrectly shadowed `/inventory/stock-seguro`
  - browser validation on the current local frontend showed the full Stock Seguro page, real Mercado Livre rows, visible blocked states, and corrected blocker summary (`BLOCKERS 5`) for the selected installation
- Blocking failure observed:
  - `No`

## Validation Notes

- Real QA found two runtime bugs during F-04 and both were fixed before milestone closeout:
  - Vite proxy shadowed the SPA route `/inventory/stock-seguro` because `/inventory` was proxied too broadly.
  - inventory JSON returned `blocking_reason.Code/Message` instead of the contract shape `blocking_reason.code/message`.
- Real QA also found a product-summary bug: `BLOCKERS` counted only unresolved/conflict rows instead of all non-actionable rows. This was fixed and revalidated in the browser on the current local frontend.
- The Docker frontend on this machine intermittently served stale source for `StockSeguroPage.tsx`; the final browser proof therefore used the current local Vite server on `http://localhost:4175`, still backed by the real backend on `http://localhost:8080`.
- Mobile viewport proof remains tool-limited in the built-in browser because the viewport capability did not reliably change `window.innerWidth` during F-04 validation. This is recorded in the feature validation but does not block the M-05 milestone contract.
- An unrelated pre-existing web-suite timeout still exists in `packages/feature-product-links/src/ProductLinksPage.test.tsx` when running the entire `@marketplace-central/web` workspace suite. It is outside the M-05 scope.

## Live Stock Seguro Validation

Date: 2026-07-09
Environment: local backend on `http://localhost:8080` plus current local frontend on `http://localhost:4175`
Installation: `inst-mercado_livre-5cb46653-95c3-46f7-a30f-9fbf5ace5f98`

Evidence:

- `GET /inventory/stock-risks?...&limit=5` returned five real Mercado Livre rows with persisted listing identity, link state, internal product evidence, source timestamps, and contract-correct `blocking_reason.code/message`.
- Browser validation opened `Stock Seguro`, selected the real Mercado Livre installation, and showed:
  - `ROWS 5`
  - `ACTIONABLE 0`
  - `OVERSELL 0`
  - `BLOCKERS 5`
  - visible `Stale`, `Sem vinculo`, and `rejected`/blocked cases with row detail and timestamps
- Browser validation on the current frontend also proved the risk filter now exposes `Ineligivel` and `Sem suporte`.

Boundary:

- No live Mercado Livre stock quantity mutation was performed.
- Manual write safety is proven through service/API/SDK/UI evidence and fake-writer tests, consistent with the validation contract.

## Handoff

- Milestone status: `ready for mission continuation`
- Next recommended action: move to the next mission/milestone scope and keep the same evidence-first standard for any live side-effecting flow
- Open blockers: none for M-05
