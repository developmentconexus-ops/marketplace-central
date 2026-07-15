# Interface Contract

```yaml
id: IC-05
type: interface-contract
status: planned
owner: Mission Strategist
parent: MIS-003
created: 2026-07-14
updated: 2026-07-14
validation_level: QA-0
lifecycle_scope: support
```

## Boundary

Frontend platform seam: AppRouter route map + redirects, installation/empresa context, web-query key/staleTime/invalidation contract, shared state-vocabulary components, pt-BR glossary. One writer (M-02) creates the seam; every later milestone consumes it read-only.

## Why This Contract Exists

Six new workspaces would each reinvent installation propagation (today: per-page `?installation=` + independent installation fetch), invalidation, and state copy. Route and query-key collisions across parallel feature workers are the classic drift seam.

## Resources Or Entities

### Route map (client, react-router-dom 7)

| Route | Screen (wireframe) | Owner milestone | Notes |
| --- | --- | --- | --- |
| `/` | Visão geral (1e) | M-05 | until M-05 lands, `/` keeps current Dashboard |
| `/anuncios` | Anúncios unificado (2a) | M-02 | params below |
| `/catalogo` | Catálogo (1f) | M-04 | replaces `/products` |
| `/catalogo/produtos/:productId` | Detalhe do produto (2b) | M-04 | `productId` = CODPROD string; deep link survives reload |
| `/vinculos` | Vínculos & Importação (1g) | M-04 | replaces `/product-links` |
| `/estoque` | Estoque (1h) | M-04 | replaces `/inventory/stock-seguro` |
| `/precos` | Preços & Simulador (1i/2d minus @mercado) | M-04 | replaces `/simulator` |
| `/pedidos` | Pedidos read (1j) | M-05 | replaces `/orders` |
| `/integracoes` | Integrações & Sync (1k) | M-05 | replaces `/integrations` |
| `/protocolos/:protocolId` | protocolo detail (fila de sync) | M-03 | `protocolId` = `MP-000042` form |
| `/classifications`, `/marketplaces` | config surfaces (not in wireframe nav) | unchanged | kept off the main sidebar; linked from `/integracoes`/`/precos` as secondary |

Legacy redirects (301-style `<Navigate replace>`): `/products→/catalogo`, `/product-links→/vinculos`, `/inventory/stock-seguro→/estoque`, `/orders→/pedidos`, `/integrations→/integracoes`, `/simulator→/precos`. Query strings preserved on redirect.

Sidebar (M-02, wireframe deck-2 order, Mercado omitted this mission): Visão geral · Catálogo · Anúncios · Vínculos & Import. · Estoque · Preços & Simulador · Pedidos · Integrações & Sync.

### Context contract

- `InstallationContext` (new, `apps/web/src/app/InstallationContext.tsx`): provides `{installationId, setInstallationId, installations[], status}`. Single fetch of `listIntegrationInstallations`; persists selection to URL `?installation=` on every route (read on mount, written on change, survives reload and navigation); falls back to first installation. Top-bar pill `ML: <nome> ▾` renders from it. Empresa pill renders but is static single-option this mission (P1b-4).
- Pages MUST consume `useInstallation()`; direct `listIntegrationInstallations` calls in pages are forbidden after M-02.

### Query contract (packages/web-query additions)

- `QUERY_STALE_TIME` gains: `listings: 45_000`, `mutations: 5_000`, `orders: 120_000`, `sync: 30_000`, `market: 300_000`. Existing catalog/stock/pricecost unchanged (mirror L2 classes).
- `queryKeyNamespaces` gains: `listings: ["listings"]`, `mutations: ["mutations"]`, `orders: ["orders"]`, `sync: ["sync"]`, `market: ["market"]`.
- Key builders (exact): `listingsQueryKeys.page(installationId, filters)`, `.byProduct(installationId, filters)`, `.detail(listingId)`, `.summary(installationId)`; `mutationsQueryKeys.list(installationId, filters)`, `.detail(protocolId)`, `.items(protocolId)`; `ordersQueryKeys.list(installationId, filters)`; `syncQueryKeys.runs(installationId, filters)`.
- New helper `invalidateAfterMutation(queryClient, type)` — THE invalidation crosswalk (single implementation, no per-page inline invalidation for envelope writes):

| Mutation type (IC-03) | Server L2 invalidation | Client namespaces invalidated |
| --- | --- | --- |
| `price_update` | — (listing read model is Postgres; refreshed by resync) | `listings`, `mutations` |
| `stock_correct` | `stock` class (server does it on apply) | `listings`, `inventory`, `mutations` |
| `link_apply` | — | `listings`, `linkage`, `catalog`, `mutations` |
| `listing_pause` / `listing_resync` / `listing_edit` | — | `listings`, `mutations` |
| product enrichment edit (IC-01 dormant row ACTIVATED, M-04) | `catalog` class | `catalog` |

- All server state via TanStack Query through these builders. Direct `useEffect` fetch is forbidden in new/rebuilt pages; remaining legacy direct-fetch pages (Dashboard until M-05, Classifications, Marketplaces until their rebuild) carry migration briefs (see mission `## Migration Briefs`).

### State vocabulary (shared components, `packages/ui`)

Exactly six states, fixed pt-BR copy:

| State | Component | Copy |
| --- | --- | --- |
| loading | `<LoadingState />` | "Carregando…" (skeleton variant available) |
| error | `<ErrorState onRetry>` | "Erro ao carregar. {detalhe}" + button "Tentar novamente" |
| empty | `<EmptyState>` | "Nenhum registro encontrado." (+ contextual hint slot) |
| stale | `<FreshnessIndicator asOf>` + refresh affordance | "dados de HH:MM:SS" (existing) |
| unknown | `<UnknownValue hint>` | renders "—" with tooltip hint, e.g. "sem custo no ERP → não simulado" |
| conflict | `<ConflictTag>` | "divergente" (amber) + detail, e.g. "divergente: ERP=35" |

Tag palette (wireframe): ok=green, warn=amber, err=red, mut=gray, blu=blue (sincronizando/processando/fulfillment).

### Glossary (binding pt-BR copy)

produto (mestre, ERP, SKU=CODPROD) · anúncio (listing) · vínculo · estoque físico/reservado/disponível/seguro · modalidade (Clássico/Premium/Full) · política de preço · margem · pedido · NF · sincronização/sync · completude (cadastro) vs qualidade (anúncio) · protocolo (recibo de mutação, `MP-nnnnnn`) · "⚑" badges de sistema-mestre: `⚑ catálogo/estoque/NF: ERP` · `⚑ preço/anúncio: HUB`. Provider raw strings only behind "▸ técnico". New UI copy is pt-BR ONLY; English strings in rebuilt pages are defects.

## Operations

Client-side contract only; no server operations. (Server contracts: IC-02/03/04 + existing OpenAPI.)

## Enums And Statuses

UI state enums above; sync_state labels fixed in IC-02.

## Error Cases

SDK `MarketplaceCentralClientError {status, error:{code,message}}` is the only error shape pages handle; `<ErrorState>` renders `error.message` when present, generic copy otherwise. Failure codes from IC-03 map to fixed pt-BR strings in ONE module (`packages/web-query/src/failureCopy.ts`).

## Persistence Expectations

Context selection persists in URL only (no localStorage this mission — reload/deep-link fidelity is the requirement).

## Canonical Examples

Deep link that must survive reload: `/anuncios?installation=inst_1&tab=pendencia&filter.exception=sync_error` → same tab, filter, and installation after F5.
Redirect example: `/inventory/stock-seguro?installation=inst_1` → `/estoque?installation=inst_1`.

## Error Matrix

| Case | Status | Code | Notes |
| --- | --- | --- | --- |
| no installations exist | — | — | context `status="empty"`; pages render `<EmptyState>` with hint "Conecte uma conta em Integrações" |
| installation in URL unknown | — | — | context falls back to first installation and REWRITES the URL param (no silent mismatch) |

## Database Shape

None (client contract).

## Seed Data

Component tests use fixed installation fixture `inst_test`.

## Timestamp And ID Semantics

`formatAsOf` (existing) is the only timestamp formatter for freshness copy.

## Compatibility Rules

- New namespaces/key builders extend; existing keys never reshaped (in-flight caches).
- Route additions extend the table above; renames require redirect rows.

## Route Namespace

- Client: table above. Reserved prefixes per milestone owner as listed — no other milestone mounts routes.
- Server: IC-02/03/04 own the `/listings`, `/mutations`, `/market` prefixes; M-05 F-01 owns the new `/dashboard` and `/sync` prefixes (shapes pinned in the M-05 F-01 brief; OpenAPI + sdk-runtime updated in the same commit).

## Transport And Integration

- Vite dev proxy list must include: `/listings`, `/mutations`, `/market`, `/orders`, `/profitability` (fixes R-02 defect), plus `/dashboard` and `/sync`. Writer sequence (single writer per commit): the first five rows are added by M-02 F-02; the `/dashboard` and `/sync` rows are added by M-05 F-01 when those prefixes go live — sequential writers, never concurrent edits to `apps/web/vite.config.ts`.
- `createRefreshableFetch`/`withNoCache` remain the forced-refresh mechanism.

## Must Preserve

- Deep links survive reload (route + installation + filter state in URL).
- One writer per seam: AppRouter/Layout/InstallationContext/web-query files owned by M-02; later milestones add routes ONLY via the route table rows assigned to them.
- IC-01 dormant row: any product-edit surface invalidates server class `catalog` + client `['catalog']`.
- TanStack Query exclusive for server state in rebuilt pages.

## Must Not Decide In Feature Execution

- Route paths, redirect map, query-key shapes, staleTimes, state copy, glossary terms, invalidation crosswalk.

## Validation Impact

M-02 criteria bind to: redirect proofs, reload-survival proofs, context single-fetch proof, state-component render proofs, crosswalk unit test.
