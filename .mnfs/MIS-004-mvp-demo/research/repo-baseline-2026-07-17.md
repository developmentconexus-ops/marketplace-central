# Research Note

```yaml
id: R-01
type: research
status: draft
owner: Codebase Investigator
parent: MIS-004
created: 2026-07-17
updated: 2026-07-17
validation_level: QA-0
lifecycle_scope: support
```

## Topic

Baseline do repo para o replan MIS-004/MIS-005: estado de merge W1, módulos backend, superfície OpenAPI, adapter ML, ERP read path, FE, tenancy, migrações, test lanes.

## Sources Checked

- Source: árvore do repo em main @ `cd74b4011abe09597ea677473bb31e6202cc56e8` (investigador dedicado, 2026-07-17).
- Why it matters: MIS-004 estende módulos existentes; decomposição e matriz de colisão dependem destes fatos.

## Findings

### Git / estado W1

- main @ `cd74b4011abe09597ea677473bb31e6202cc56e8`; branch do checkout hub = `main`.
- Chips W1 NÃO mergeados: `chip/m-02-frontend-platform-anuncios` @ 10f6a1f5, `chip/m-03-mutation-envelope-writes` @ 0923b31c, `chip/sat-m05f01-m06f02` @ 8996d788 — todos divergindo de `a49168e6`, worktrees irmãs vivas.
- SAT parcialmente aterrissado em main (F-4 fixtures @ 2c095f3b); registro governance do módulo mutations foi commitado e REVERTIDO (66892000 → f869c4a6) — entra de novo só via merge do chip M-03.

### Módulos backend (contracts/governance/modules.json — 14)

- `catalog` (identidade CODPROD, enrichment) · `classifications` · `connectors` (adapters ML/MelhorEnvio/fees) · `dashboard` (/dashboard/summary) · `integrations` (OAuth, credenciais AES em `integration_credentials`, `MPC_ENCRYPTION_KEY`) · `internal_read` (boundary ERP/Oracle, zero OpenAPI) · `inventory` · `listings` (read model ML + refresh) · `market` (contract-only: /market/observations, /market/references) · `marketplaces` · `orders` (import + Sankhya linkage) · `pricing` (/pricing/simulations + batch) · `product_links` (imports/generations/candidates/workflows/resolutions) · `profitability`.
- Composition root: `apps/server_core/internal/composition/root.go` `NewRootRuntime()` (198-566).
- 4 temporary_exceptions no governance, todas owner M-10 (connectors→integrations fee edges).

### OpenAPI (contracts/api/marketplace-central.openapi.yaml)

- Grupos existentes: /catalog/products*, /listings*, /classifications*, /integrations/*, /product-links/*, /inventory/*, /orders* (+/{id}/sankhya-linkage), /profitability/*, /marketplaces/*, /admin/fee-schedules/*, /pricing/simulations(+/batch), /connectors/melhor-envio/*, /dashboard/summary, /market/observations|references, /sync/runs (GET-only).
- Gaps: sem /mutations* (chip M-03), sem import ERP.

### sdk-runtime

- `packages/sdk-runtime` é client TS ESCRITO À MÃO (sem codegen; scripts só build/test — package.json:7-10). Profile chama de "generated" — mismatch doc/realidade. Regra prática: OpenAPI + edição manual do sdk-runtime no MESMO commit.

### Migrações

- `apps/server_core/migrations/*.sql`, prefixo 4 dígitos, ordenação lexicográfica por filename completo; topo = `0044_market_references.sql`; 39 arquivos (gaps 0002, 0038-0042); prefixo 0021 duplicado 2×.
- Fixture: `internal/platform/migrate/runner_test.go:25-27` hardcoda 39 — toda migração nova bumpa o literal.

### ERP read path (internal_read)

- Port `ports/reader.go:48-55`: FindProductsForLinking, GetSellableStock, GetCurrentPrice, GetCostAsOf (basis `cussemicm`), GetSalesHistory, GetTaxInputs. Ports extras: stock_batch, batch_reader, sankhya_linkage, catalog_page. Domain inclui icms_ceiling.
- SQL: TGFPRO.CODPROD/REFERENCIA/REFFORN/CODGRUPOPROD/DESCRPROD (`adapters/oracle/reader.go:69,153,220,374,389-426`). Match hoje por CODPROD ou REFFORN; REFERENCIA tratada como "governed manufacturer/reference value", NÃO como EAN — contradiz research de identidade (REFERENCIA=EAN no Sankhya). Correção de contrato necessária (research §4).
- Oracle indisponível → composition root injeta `Unavailable*` readers que erram sempre (root.go:327-362) — nunca zero-default.

### Adapter Mercado Livre (connectors/adapters/mercado_livre/capability_adapter.go)

- Já chama: GET /users/me, GET /users/{id}/items/search, GET /items/{id}, PUT /items/{id} (stock, gated: idempotency key), GET /orders/search, GET /orders/{id}, GET /sites/{id}/listing_prices.
- FALTAM p/ MIS-004: /items/{id}/sale_price, /items/{id}/price_to_win?version=v2, /products/search, /products/{id}, /products/{id}/items (flag), /shipments/{id}(+costs/delays), /users/{id}/shipping_options/free.
- OAuth: `integrations/adapters/mercadolivre/auth_adapter.go`; tokens cifrados via CredentialResolver.

### M-03 mutation envelope

- NÃO está em main. Design em `.mnfs/MIS-003-.../mission.md`: ADR-13 protocolo table + in-process poller; lifecycle 8 estados (IC-03); caminho UI → /mutations → poller → capability adapter. Implementação no chip não mergeado.

### Frontend (apps/web)

- Rotas atuais (AppRouter.tsx:54-72): /, /products, /classifications, /marketplaces, /integrations, /product-links, /inventory/stock-seguro, /orders, /simulator. Sem /listings em main (chip M-02).
- Shell: Layout.tsx sidebar escura hardcoded, labels EM INGLÊS, lucide-react. SEM mecanismo de tema (nenhum ThemeProvider/data-theme).
- Query state: @tanstack/react-query via packages/web-query (staleTimes por domínio, cache-busting fetch), SDK client via ClientContext (VITE_API_BASE_URL, dev localhost:8080).

### Tenancy/auth

- SEM middleware de auth/tenant. CORS `*` dev-only (httpx/router.go:16-29). Tenant fixo `MC_DEFAULT_TENANT_ID` → "tenant_default" (pgdb/config.go:14-30), threaded à mão nos repos. Single-tenant/no-auth hoje; multi-tenant real = greenfield (MIS-005).

### Test lanes (contracts/governance/execution-lanes.json)

- 7 lanes: unit, integration (postgres efêmero), dev-invariance, live-oracle, live-provider-read, browser, provider-write (7 gates: actor/idempotency/execute/resolved-link/policy/source-timestamp/before-after-audit).
- Go: `//go:build integration` central (tests/integration, 9 arquivos) + colocado (~16); scripts npm harness:* → scripts/harness.ps1.

## Recommendation

Planejar MIS-004 como EXTENSÃO dos módulos existentes (market, connectors, product_links, pricing, orders, internal_read, dashboard) — não criar módulos paralelos sem razão. Declarar base = main pós-merge W1 e a dependência explicitamente. Registrar mismatch sdk-runtime "generated" vs hand-written.

## Impact On Mission

Muda decomposição de milestones (extensão vs greenfield), matriz de colisão (adapter ML = seam único connectors; shell FE = seam M-02-chip), blocos de migração (0045+), e adiciona milestone de auth/tenancy ao MIS-005.

## Handoff

- Current status: completo para planning; worktrees dos chips NÃO inspecionadas (fora de escopo baseline-main).
- Next owner: Mission Strategist (P3).
- Next action: candidato de escopo P3.
- Required files/evidence: este arquivo; modules.json; openapi.yaml.
- Blockers or open decisions: nenhum bloqueante; merge W1 é pré-condição de execução, não de planning.
