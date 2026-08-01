# MIS-007 research — lado READ + contratos (fatos de codebase)

```yaml
id: MIS-007-research-codebase-read-side
type: research-evidence
author: codebase-investigator (P2)
date: 2026-07-31
base: main (dd89d4b3)
method: leitura estática de arquivos; nenhum comando executado; nenhum código alterado
```

Convenção: caminhos relativos à raiz do repo. Todo fato tem `path:line`. Interpretação separada em "Leitura" quando necessário.

---

## 1. Sítios que chamam ML no READ (alvo da Onda 0)

Todo tráfego ML passa por um único client HTTP: `apps/server_core/internal/modules/connectors/adapters/mercado_livre/capability_adapter.go:22` — `defaultBaseURL = "https://api.mercadolibre.com"`.

**Contagem: o "4" do design CONFIRMA como 4 ADAPTERS read-time fiados na composition root para rotas interativas de tela.** Em termos de rotas HTTP interativas são 4: `GET /orders`, `GET /orders/{provider_order_id}`, `POST /pricing/decompose`, `POST /pricing/solve`. Há OUTROS chamadores ML em request-time que NÃO são read-de-tela (lista completa em 1.5) — o recorte "4" só fecha sob a definição "handler interativo que uma tela chama para LER".

### 1.1 Sítio A — GET /orders (lista de /pedidos) → shipment por pedido

- `apps/server_core/internal/composition/root.go:591` — `ordersShipmentReader := newOrdersShipmentReaderAdapter(mercadoLivreCapabilities, installationSvc, cfg.DefaultTenantID)`
- `apps/server_core/internal/composition/root.go:597-599` — o reader entra no `ordersEnrichSvc` (`NewEnrichServiceWithReaders(...)`), que o handler usa em list E detail.
- `apps/server_core/internal/composition/orders_adapters.go:98-108` — `GetShipment` resolve installation e chama `a.capabilities.GetShipmentInfo(ctx, ref, shipmentID)` (ML vivo).
- `apps/server_core/internal/modules/orders/application/enrich_service.go:150-162` — `Enrich` (caminho de LISTA): um `errgroup` com `g.SetLimit(shipmentConcurrency)` dispara `fetchShipment` POR PEDIDO da página.
- `apps/server_core/internal/modules/orders/application/enrich_service.go:21` — `const shipmentConcurrency = 8`.
- `apps/server_core/internal/modules/orders/application/enrich_service.go:287-291` — `fetchShipment` → `s.shipment.GetShipment(ctx, installationID, order.ShippingID)`.
- Cada `GetShipmentInfo` = **3 GETs ML sequenciais**:
  - `apps/server_core/internal/modules/connectors/adapters/mercado_livre/shipping_reader.go:149-155` — `GET /shipments/{id}` (header `x-format-new: true` obrigatório);
  - `shipping_reader.go:162` + `:204-211` — `GET /shipments/{id}/carrier` (transportadora/tracking; falha degrada p/ nil);
  - `shipping_reader.go:164-180` — `GET /shipments/{id}/costs` (frete real; 404/decode degrada, resto derruba a leitura).
- Handler: `apps/server_core/internal/modules/orders/transport/http_handler.go:86` (`mux.HandleFunc("/orders", h.handleList)`) → `:164-203` `handleReadList` chama `h.enricher.Enrich(r.Context(), installationID, page.Items)` (linha 194) DENTRO do request.
- O que a tela mostra com isso: `bucket` (fila/kanban — `http_handler.go:549`, `DeriveOrderBucket` consome status do shipment vivo), `sla` (`:556`), `destino_uf/cep`, `destinatario`, `frete_real` (`:557-560`), `rastreio` status/substatus/transportadora/url (`:561-567`).

### 1.2 Sítio B — GET /orders/{id} (drawer de /pedidos) → buyer fiscal 2-step

- `apps/server_core/internal/composition/root.go:592` — `ordersBuyerFiscalReader := newOrdersBuyerFiscalReaderAdapter(...)`.
- `apps/server_core/internal/composition/orders_adapters.go:140-150` — `GetBuyerFiscal` → `a.capabilities.GetBuyerFiscalInfo(ctx, ref, providerOrderID)`.
- `apps/server_core/internal/modules/orders/application/enrich_service.go:194-198` — `EnrichOne` (caminho de DETALHE) = `Enrich` de 1 pedido (logo, os 3 GETs de shipment) + `resolveBuyerFiscal` (`:343-347`).
- Comentário load-bearing: `enrich_service.go:143-149` — a lista NÃO resolve buyer fiscal porque são "+2 sequential ML calls per order" e isso "would blow the request deadline across a full list page (FINDING-M08-LIST-TIMEOUT)".
- Adapter ML: `apps/server_core/internal/modules/connectors/adapters/mercado_livre/buyer_fiscal_reader.go:59-94` — fluxo documentado em 2 passos, termina em `GET /orders/billing-info/{site}/{billing_info_id}` (`:94`).
- Handler: `orders/transport/http_handler.go:89` (`/orders/{provider_order_id}`) → `:237-239` `EnrichOne` no request. Tela mostra: `comprador_fiscal` (nome, doc, endereço NF — DTO `:462-467`).

### 1.3 Sítios C e D — POST /pricing/decompose e /pricing/solve (tela /precos) → tarifa viva (categoria + comissão)

- `apps/server_core/internal/composition/root.go:845-851` — quando `FeeQuoteReader("mercado_livre")` existe: `categoryResolver` (`:847`, sítio C) + `commissionQuoter` (`:848`, sítio D) compõem `pricingtarifflive.NewResolver` sobre o resolver de defaults (degrau 4), via `pricingtariffcomposite`.
- `apps/server_core/internal/composition/pricing_adapters.go:58` — `newPricingCategoryResolverAdapter(catalogMatch connectorsports.CatalogMatchReader, ...)`; ML: `catalog_match_reader.go:78` `GET /products/search`, `:123` `GET /sites/{site}/domain_discovery/search`, `:98-` `GET /products/{id}` (buy box).
- `apps/server_core/internal/composition/pricing_adapters.go:123` — `newPricingCommissionQuoterAdapter(fees connectorsports.FeeQuoteReader, ...)`; ML: `capability_adapter.go:571` `GET /sites/{site}/listing_prices?...` (`ReadFeeQuote`, `:546`).
- Uso no serviço: `apps/server_core/internal/modules/pricing/adapters/tarifflive/resolver.go:43-69` — `ResolveCommission` resolve categoria e então `QuoteCommission` (2+ chamadas ML por decompose que precisa do degrau 3).
- Rotas: `apps/server_core/internal/modules/pricing/transport/calc_handler.go:35-36` — `POST /pricing/decompose`, `POST /pricing/solve`. Ambas fora de `registerBatchRoutes` → classe interativa 15s.
- Também compartilha o resolver: `root.go:852` — `batchOrch.WithTariffResolver(tariffResolver)` → `POST /pricing/simulations/batch` (classe batch, `root.go:266`).
- Tela: `apps/web/src/pages/precos/PricingPage.tsx:228-230` — `useQuery` → `client.pricingDecompose(...)` no load/edição; `apps/web/src/pages/precos/PricingMatrix.tsx:117` — `pricingDecompose` POR LINHA da matriz (multiplica o custo vivo).

### 1.4 Confirmação do fechamento do recorte

- `GET /listings` (tela /anuncios) NÃO chama ML: `root.go:701-714` — `listingSvc` = repo Postgres (`listingspostgres.NewRepository`) + custo Oracle + policy + installations + evidence (Postgres do market). Nenhum tipo `connectors` no read de listings.
- `GET /orders` lê o pedido base do Postgres (`orderspostgres.NewOrderReadRepository`, `root.go:587`) — o ML entra só pelo enrich.

### 1.5 Outros chamadores ML em request-time (NÃO contam no "4", mas existem)

- `POST /market/collections` — `apps/server_core/internal/modules/market/transport/http_handler.go:56`; pipeline usa `marketPriceIntelReader` (`root.go:684`), ML vivo (`market_adapters.go:320-338` `SearchCatalogByEAN` etc.). Ação da tela /mercado; escopo MIS-008 (mission.md Non-Scope).
- Probes `GET|POST /integrations/installations/{id}/probes/{account|listings|orders|fee-quote|stock|catalog-match}` — `apps/server_core/internal/modules/integrations/transport/auth_handler.go:262-382`; ML vivo via `providerOperationSvc` (`root.go:400-407`). NENHUMA tela do FE chama probes hoje (grep `probeIntegration` em `apps/web/src` = 0 usos fora do SDK).
- `POST /orders/import` — classe batch (`root.go:263`); lê pedidos do ML e grava no PG. Disparado pelo botão "Atualizar" de /pedidos (`apps/web/src/pages/pedidos/PedidosPage.tsx:163-166`).
- `POST /listings/refresh` — `apps/server_core/internal/modules/listings/transport/http_handler.go:30`; responde 202 e roda a ingestão ML em goroutine (`root.go:728-736`, `func(task func()) { go task() }`). Disparado pelo `ListingsRefreshControl` de /anuncios.
- `POST /product-links/listing-snapshots/imports` e `/link-candidates/generations` — classe batch (`root.go:264-265`); fonte = `providerOperationSvc` (`root.go:422-425`) e âncoras de identidade ML (`root.go:389`).
- Background (não-request): `mutationsbg.Poller` (`root.go:766-777`), `FeeSyncScheduler` 15min (`root.go:663`), `RefreshTicker` de credenciais (`root.go:661`).

---

## 2. Telas

### 2.1 /pedidos

- Componente: `apps/web/src/pages/pedidos/PedidosPage.tsx:122` (`PedidosPage`); views `FilaView`/`KanbanView`/`ListaView` + `PedidoDrawer`.
- Endpoint: `PedidosPage.tsx:41-64` — `fetchAllOrders` acumula `client.listOrders(...)` (GET /orders) seguindo `next_cursor` até esgotar (cap `MAX_ORDER_PAGES = 20`, `:32`). UMA query traz o dataset inteiro; KPI+tabs derivam dele (`:34-40`).
- Drawer: `getOrder` → `GET /orders/{id}` (`packages/sdk-runtime/src/index.ts:2290-2293`).
- Botão Atualizar: `importMarketplaceOrders` → `POST /orders/import` (`PedidosPage.tsx:163-166`; SDK `index.ts:2281-2282`).
- Colunas da Lista: `apps/web/src/pages/pedidos/PedidosTable.tsx:95-146` — **Pedido · Data · Comprador · Itens · Valor · Retorno · SLA · DIFAL · Ação**.
- **Onde a lista gasta tempo (a ~10.8s conhecida):** `GET /orders` com enrich inline — por página, até 8 shipments em paralelo × 3 GETs ML cada (`enrich_service.go:150-162` + `shipping_reader.go:147-194`), repetido por página do `fetchAllOrders`. A rota é classe interativa (15s — `route_deadline.go:26`), então a lista vive a <5s da parede. O número 10.8s é medição registrada no fechamento do M-08 (memória `m08-pedidos-closed`, "lista ~10.8s = chip futuro bounded-parallel"); não re-medido nesta pesquisa.

### 2.2 /anuncios

- Componente: `apps/web/src/pages/AnunciosPage.tsx` + `AnunciosTable.tsx`.
- Endpoints: `apps/web/src/pages/AnunciosPage.tsx:163` — `client.listListings` ou `client.listListingsByProduct` (agrupado); `:173` — summary via `anunciosSummaryQuery` → `getListingsSummary` (`anunciosQueries.ts:23-29`). Client da página: `anunciosQueries.ts:5-8` — `"listListings" | "getListingsSummary" | "refreshListings" | "listIntegrationOperationRuns"`.
- Refresh manual: `apps/web/src/pages/ListingsRefreshControl.tsx:89` — polling de `listIntegrationOperationRuns(installationId)` (GET /integrations/installations/{id}/operations) após `refreshListings` (POST /listings/refresh, 202+409 `refresh_in_progress`).
- Colunas: `apps/web/src/pages/AnunciosTable.tsx:23-25` — "Ratified 9-column set: checkbox · MLB · TÍTULO · PRODUTO · PREÇO · EST. · SYNC · QUAL. · PENDÊNCIA" (`<th>` em `:287-294`).

### 2.3 /precos

- Componente: `apps/web/src/pages/precos/PricingPage.tsx` + `PricingMatrix.tsx`.
- Endpoints (`PricingPage.tsx:85-141, 173-176, 228-230`): `getPricingProfile` (GET /pricing/profile), `listListings({link_state:"resolved"})`, `catalogProductFactsByIds`, `listPricingDifal` (GET /pricing/difal), `putPricingProfile`, `putPricingDifalOverride`, `listListingsByProduct`, `pricingDecompose` (POST /pricing/decompose). Matrix: `PricingMatrix.tsx:117` decompose por linha. Solve: `pricingSolveTarget` no client (testes `PricingPage.test.tsx:97`).

### 2.4 /integracoes

- Componente: `apps/web/src/pages/integracoes/IntegracoesPage.tsx:558-574` — seções: `ActiveSourceCard` (GET/PUT /config/active-source), `SellableAssortmentCard` (GET/PUT /config/sellable-assortment + counts), `UploadCard` (POST /erp/imports multipart), `ProviderConnectCard` (`:496-556` — listIntegrationInstallations, createIntegrationInstallation, startIntegrationAuthorization → redirect OAuth ML), `ImportacaoSection` (histórico de imports ERP).
- **NÃO existe hoje seção de saúde de sync por entidade nem status de webhook** nessa tela (Q4 da missão é lacuna).

---

## 3. OpenAPI + SDK

- Spec: `contracts/api/marketplace-central.openapi.yaml` (OpenAPI 3.1.0, `info.version: 2026-07-13` — linha 1-4). Rotas relevantes já lá: `/listings:647`, `/listings/refresh:919`, `/orders:1848`, `/orders/{provider_order_id}:1952`, `/pricing/decompose:2718`, `/pricing/solve:2747`, `/sync/runs:3172`, probes `1357-1466`. **Não existe `/webhooks/*` no spec.**
- SDK: `packages/sdk-runtime/src/index.ts` — pacote `@marketplace-central/sdk-runtime` v0.1.0 (consumido por `apps/web/package.json:19`). É ESCRITO À MÃO (não gerado): client factory `createMarketplaceCentralClient` retorna um objeto literal de métodos (`index.ts:2113-2330+`), um método por rota, ex.: `listOrders` (`:2288-2289`), `refreshListings` (`:2146-2147`), `listSyncRuns` (`:2144-2145`).
- Convenção de rota nova: adicionar path no YAML + método no client do SDK + tipos; regra vinculante: `docs/HARNESS-PROFILE.md:248` — "OpenAPI + `sdk-runtime` land in the same commit"; `:193` — contract lock (um dono por vez, hub-owned `:202`).
- Exemplo recente: commit `bce555a8` "fix(contracts): /marketplaces/fee-schedules 400 declares the error envelope it already returns" (git log do bootstrap); e o comentário-âncora do DTO de summary: `orders/transport/http_handler.go:290-296` — "Keys are locked byte-identical to the OpenAPI OrderSummary schema and the SDK OrderSummary interface".

---

## 4. Envelope de erro (pós CHIP-ERROR-UNIFY)

- Produtor único: `apps/server_core/internal/platform/apierror/apierror.go:1-3` — "Package apierror is the single producer of the backend's HTTP error envelope."
- API: `apierror.go:25` — `func Write(w http.ResponseWriter, status int, code, message string, details map[string]any)`; shape no fio `{"error":{"code","message","details"}}` (`:12-20`); `details` nunca ausente (`:26-28`).
- Recover central: `apps/server_core/internal/platform/apierror/recover.go:20` — `func Recover(next http.Handler) http.Handler`; aplicado no root: `root.go:863` — `httpx.CORSMiddleware(apierror.Recover(mux))`.
- Exemplo canônico de handler: `orders/transport/http_handler.go:100` — `apierror.Write(w, http.StatusMethodNotAllowed, "ORDERS_METHOD_NOT_ALLOWED", "method not allowed", nil)`; mapeamento de erros de leitura `:689-704` (`writeOrderReadError`).
- Deadline também fala o envelope: `platform/httpx/route_deadline.go:129` — `{"error":{"code":"deadline_exceeded","message":"tempo limite excedido","details":{}}}` (literal para evitar ciclo de import, `:124-128`).
- FE/SDK: `packages/sdk-runtime/src/index.ts:1807-1821` — `MarketplaceCentralClientError` (Error real com `.status/.code/.details`); `:1958` `parseApiError` (parse único; payload malformado vira `internal_error` com `details.raw`); `:1897-1909` `hasCode<C>()` — checa código E valida os campos de `details` que o spec garante (`REQUIRED_DETAIL_FIELDS`, `:1867-1873`); códigos novos entram nas unions por domínio (`:1727-1796`, "sourced verbatim from the `code` enum ... in contracts/api/...yaml"). Uso no FE: `apps/web/src/pages/AnunciosPage.tsx:135`, `integracoes/useErpImportUpload.ts:32-48`, `mutations/MutationPreviewModal.tsx:125`.
- 409 fora das unions tem guard dedicado: `index.ts:1920-1933` — `isRefreshInProgressError`.

---

## 5. Transport / rotas / classes de deadline

- Composition root: `apps/server_core/internal/composition/root.go:274` — `NewRootRuntime` cria `httpx.NewRouteClassMux()` (`:281`), registra classes batch (`:282`) e então cada módulo se registra via `XXXtransport.NewHandler(...).Register(mux)` (ex.: `:416-418`, `:601`, `:714`, `:736`, `:858`).
- Classes: `apps/server_core/internal/platform/httpx/route_deadline.go:23-28` — `interactiveDeadline = 15s`, `batchDeadline = 120s`; default = interativa (`:44-45`); classe declarada pelo PATH sem método (`:47-53`).
- Lista batch atual: `root.go:259-272` — `registerBatchRoutes`: `/profitability/margin-inputs/import`, `/profitability/profit-snapshots/calculate`, `/orders/import`, `/product-links/listing-snapshots/imports`, `/product-links/link-candidates/generations`, `/pricing/simulations/batch`, `/admin/fee-schedules/sync`, `/admin/fee-schedules/seed`.
- Leitura (para a missão): `POST /webhooks/{provider}` entraria como rota nova registrada no root com `mux.HandleFunc` (padrão method-aware do ServeMux, ex. `"POST /market/collections"` em `market/transport/http_handler.go:56`); precisa decidir classe (o transport fino de 200-imediato cabe na interativa). Endpoints de refresh em lote seguem o padrão `registerBatchRoutes` (`root.go:260` — adicionar o path à lista ANTES do registro do handler).
- Não existe hoje NENHUMA rota `/webhooks` no backend: grep `webhook` em `internal/` só acha strings de capability (`integrations/domain/runtime_capability.go:21`, `marketplaces/registry/mercado_livre.go:30`).

---

## 6. /integracoes hoje — saúde/status de sync no backend

- `GET /sync/runs` — `apps/server_core/internal/modules/integrations/transport/run_read_handler.go:34`; classe interativa explícita (`:32`). Shape: `runPageEnvelope{items, next_cursor, page_size, as_of}` (`:37-42`); item: `operation_run_id, installation_id, module, status, result_code, failure_code, translated_error_code, attempt_count, duration_ms, started_at, finished_at` (`:44-56`). SDK: `listSyncRuns` (`index.ts:2144-2145`) — **sem uso no FE hoje** (grep em `apps/web/src` e `packages/feature-*` = 0).
- `GET /integrations/installations/{id}/operations` — `auth_handler.go:246-261`; usado pelo `ListingsRefreshControl` de /anuncios (não pela tela /integracoes).
- `GET /integrations/installations/{id}/auth/status` — `auth_handler.go:209-226`.
- sync_state (0075/MIS-006): scheduler de produtos ERP — `root.go:668-678` (`synccomposition.NewProductsScheduler`, 15min) + `syncpg.NewSyncStateRepository` (`root.go:557`). Entidades `listings`/`orders` no sync_state: ainda não existem (escopo da missão).

---

## 7. Padrão de teste

- Lane hermética Go: `apps/server_core/tests/integration/` — 10 arquivos com `//go:build integration` (grep): `market_contract_test.go`, `listings_read_test.go`, `listings_read_perf_test.go`, `listings_refresh_test.go`, `aggregate_sync_read_test.go`, `product_link_decisions_test.go`, `integrations_fee_sync_test.go`, `integrations_credential_repo_test.go`, `marketplaces_repository_test.go`, `phase1_smoke_test.go`, `migrate_runner_test.go`. Regra da lane (memória `integration-lane-not-superset-of-unit`): só `./tests/integration` + módulos com a tag nas 5 primeiras linhas; `GOCACHE=.gocache` (AGENTS.md).
- Guard anti-scraping/anti-live no contrato de market: `tests/integration/market_contract_test.go:533` — `forbidden := []string{"api.mercadolibre.com/sites", "search?", "colly", "goquery", "chromedp"}`.
- Fixtures com paginação: `tests/integration/listings_read_test.go` — `seedListing`/`seedContractRows` (`:103`, `:158-171`) + 12 usos de cursor no arquivo; lição CHIP-MERCADO (memória `chip-mercado-closed`): truncamento de página-1 é invisível ao live-drive quando o DB tem <50 rows — "prove paginação com fixture >1 página" (também R-3 da mission.md:140).
- FE: vitest — `apps/web/package.json:10` (`vitest run --config vitest.config.ts`); testes por página (`PedidosPage.test.tsx`, `AnunciosPage.test.tsx`, `PricingPage.test.tsx` etc.) mockando o client via `useClient` (`PricingMatrix.test.tsx:143`). Worktree: `npm ci` na raiz (profile §3; junction banida).

---

## 8. Fee estático 16%/22% — estado REAL diverge do design

- O seeder flat 16% Clássico / 22% Premium foi REMOVIDO na main: `apps/server_core/internal/modules/integrations/adapters/providers/registry_test.go:90-103` — "No marketplace seeds fee schedules any more. The ML seeder wrote a flat 16% Clássico / 22% Premium ... An empty schedule is honest"; `TestBuildFeeSyncersSeedsNothing` exige `len(BuildFeeSyncers()) == 0`.
- Migração que apagou as linhas legadas: `apps/server_core/migrations/0081_drop_legacy_marketplaces_and_orphan_tables.sql:15-16` — "not one flat 16%/22% per listing type ... contradicted pricing_tariff_defaults (13%/16%)".
- **A referência do design `fee_sync.go:29` está OBSOLETA**: nenhum `fee_sync.go` na árvore carrega 16%/22% (arquivos existentes: `connectors/application/fee_sync_service.go`, `integrations/application/fee_sync_service.go`, `integrations/background/fee_sync_scheduler.go`). O outcome "seed estático morre" já aconteceu para o SEED; o que resta de fee sem proveniência é o degrau-4 de `pricing_tariff_defaults` (config do tenant) + o degrau-3 vivo (sítios C/D acima).

---

## Lacunas (o que NÃO existe hoje)

- `POST /webhooks/{provider}`: zero rota, zero handler, zero entrada no OpenAPI (§5).
- `notifications_inbox`, `channel_fees`, `order_shipments`, `divergences`, `listing_variations` formalizada: nenhuma migração/tabela (migrações vão até 0085; nenhuma cita esses nomes).
- Decomposição de margem PERSISTIDA: hoje é calculada NO READ (`enrich_service.go:373-`) e nunca gravada; custo não é congelado.
- Backfill retomável de listings/orders: não existe; `POST /orders/import` e `POST /listings/refresh` são pulls limitados sem cursor persistido de backfill.
- `sync_state` para entidades `listings`/`orders`: scheduler atual cobre só produto ERP (`root.go:668-678`); `erpsync.NewMarketEnqueuer` (`root.go:557`) enfileira, mas não há scheduler ML de listings/orders.
- Backoff exponencial/jitter/Retry-After no client ML: `capability_adapter.go` `doRawWithHeaders` (`:712-`) não tem retry/backoff (pré-requisito das ondas).
- Regra DTO `Raw json.RawMessage`: os DTOs ML atuais (`mlItemResponse`, `mlOrderResponse`, `mlShipmentResponse`) não persistem raw.
- Saúde de sync por entidade em /integracoes: a tela não tem a seção (§2.4); `GET /sync/runs` existe mas nenhum FE o consome (§6).
- Badge/filtro de divergência em /anuncios: colunas atuais são as 9 ratificadas (§2.2), sem divergência.
- Colunas E3 (sold_quantity, category_id, permalink, thumbnail, estoque ML etc.): ausentes do read model de listings (9 colunas atuais; repo PG `listingspostgres`).
- Medição atual do 10.8s: valor vem de memória de missão (M-08), não re-medido aqui — re-medir no live-drive da Onda 0.
