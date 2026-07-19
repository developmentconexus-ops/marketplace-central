# CHIP-FACT-FE — evidence pack

Mission: MIS-004-mvp-demo · Scope: migrate FE consumers off orphaned `CatalogProduct` SDK type onto ratified `CatalogProductFact`
Base SHA (fork point): `910d819688b5db64a870ccc26d464b14d32ffd0a`
Commit: `390b35bcd647f261c36a8df92033032346bd88f3`
Branch: `chip/fact-fe`
Migration grant: N/A (FE-only, zero backend/OpenAPI/SDK touch)

## Scope delivered
1. **feature-classifications** (ROUTED /classifications — live runtime defect): consumed orphaned
   `CatalogProduct` type → names/SKUs/IDs/prices/stock rendered blank at runtime. Migrated to
   `CatalogProductFact` via `listCatalogProductFacts` (cursor-drained pagination, 100/page).
2. **feature-simulator** (UNROUTED — superseded by in-app PricingPage/IC-04): `PricingSimulatorPage`
   proven fully dead (no source import anywhere; only inert `vi.mock` in AppRouter.test + tailwind
   `@source`). Deleted the 2 files; kept the package per hub bar on full-package deletion.

## Field map (CatalogProduct → CatalogProductFact)
| old | new |
|-----|-----|
| product_id (string) | internal_product_id (**number**) |
| name | description |
| sku | reference |
| cost_amount | cost.amount |
| price_amount | current_price.amount |
| stock_quantity | sellable_stock.quantity |

Dropped: `taxonomy_node_id` filter (backend `ListTaxonomyNodes` is a stub returning `[]`),
`ean` search (Fact.ean hardcoded NULL), `brand` column/search (brand_name projection DEFERRED
post-demo per hub ruling — NOT added, no placeholder).

## Number↔string membership seam (reconciliation table)
`Fact.internal_product_id` is `number`; `Classification.product_ids` is `string[]`; create/update
payloads take `string[]`. **Single stringification point** `factId(fact) = String(fact.internal_product_id)`.
Every membership/mutation path receives the already-stringified id — verified end-to-end:

| path | id source | stringified? |
|------|-----------|--------------|
| checkbox checked-state (`checkedIds.has(id)`) | `factId(p)` | ✓ |
| toggle (`handleToggleProduct(id)`) | `factId(p)` | ✓ |
| select-all-filtered (`filtered.map(factId)`) | `factId` | ✓ |
| clear-all (`[]`) | — | ✓ |
| create-from-first-check | `factId(p)` | ✓ |
| row `key` / `aria-label` | `factId(p)` | ✓ |
Edge `internal_product_id: 0` → `String(0)="0"` (truthy non-empty) — no falsy-id bug.

## ADR-17 honest-absent rendering
- `formatCurrency(amount: string|null)` → `null`/`NaN` render `"—"`, never `0`/blank.
- Stock: `p.sellable_stock.quantity ?? "—"` — uses `??` not `||`, so quantity `0` renders `0` (not `"—"`).
- Search: `p.description?.toLowerCase().includes(q) ?? false` (+ same for reference) — null-safe.

## Gate results (VERIFY)
- vitest **33/33 green** — feature-classifications 15/15, feature-simulator (post-delete package
  tests) + apps/web AppRouter.test 18/18 (vi.mock inert, still passes).
- `npx tsc --noEmit` (apps/web, chip-local paths → worktree packages): **0 CatalogProduct/Fact drift
  errors** (the 2 baseline drift errors GONE). 5 unrelated pre-existing errors remain, all OUT OF
  SCOPE (anunciosQueries ListingListOptions; MutationPreviewModal/MutationResultSummary onRetry) —
  owned by other seams. 449 jest-dom/import.meta.env fork-noise errors are the vite-client/jest-dom
  d.ts gap owned by task_232262e6, not this chip.
- `vite build` **green** (4.01s).

## Scope discipline
- Zero backend / OpenAPI / sdk-runtime edits. Zero apps/web/tsconfig.json touch (task_4238de79 owns).
- Zero apps/web/src/pages Anuncios* touch (M-05 owns).
- brand_name projection NOT added (hub deferral). taxonomy filter NOT stubbed-in.
- All chip-local verify artifacts (tsconfig.chip.json, vitest.chip.config.ts, vite.chip.config.ts,
  node_modules junction, dist) deleted pre-commit — none in the diff.

## Dual-gate review (Claude-only lane, codex quota DEAD til 2026-07-25)
- **Cold Opus reviewer** (model=opus, clean context): **NO BLOCKER**. Confirmed number→string seam
  consistent end-to-end, tests genuinely exercise it. 3 MINOR (formatCurrency empty-string edge;
  multi-page-drain test-coverage gap; simulator residue — all non-blocking).
- **Adversarial sonnet reviewer** (attack every vector): **clean against all** — id seam single
  stringification point, cursor loop terminates (`null` & `""` both fall out), stock-0 renders `0`,
  null-safe search, tests retain assertion strength. Same 2 low/cosmetic notes only.
- **Agreement: PASSED.**

## Residue flagged to HUB (cross-seam cleanup — hub-barred from this chip)
feature-simulator package kept (per hub bar on full-package deletion), so 2 dangling references
survive and want a hub-owned sweep:
- `apps/web/src/index.css:6` — `@source "../../../packages/feature-simulator/src";` (dir now empty;
  Tailwind v4 glob on missing path resolves to zero matches, does NOT error → build stays green).
- `apps/web/package.json` — `"@marketplace-central/feature-simulator": "0.1.0"` dep still listed.
Neither breaks the demo build. REQUEST: hub decide full feature-simulator package removal post-demo.

## Untested branch (informational, non-blocking)
`loadAllFacts` cursor-drain multi-page branch has no test (all mocks single-page `next_cursor:null`);
logic verified correct by both reviewers; no iteration cap. Demo dataset is single-page — branch
won't execute tomorrow. Post-demo hardening candidate.
