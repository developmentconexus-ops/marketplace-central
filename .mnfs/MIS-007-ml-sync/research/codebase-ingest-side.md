# MIS-007 ML-sync — Evidência do lado INGEST (codebase-investigator, 2026-07-31)

Base: main @ dd89d4b3. Backend Go em `apps/server_core`. Todos os paths abaixo relativos à raiz do repo, file:line exatos.

## 1. Adapter ML atual

- `apps/server_core/internal/modules/connectors/adapters/mercado_livre/capability_adapter.go:1` — pacote real chama-se `mercadolivre` mas vive no diretório `mercado_livre`; é O adapter ML de capacidades (readers+writers).
- `capability_adapter.go:22-23` — `defaultBaseURL = "https://api.mercadolibre.com"`, `defaultSiteID = "MLB"`.
- `capability_adapter.go:62` — HTTP client default `&http.Client{Timeout: 15 * time.Second}`; sem transporte custom.
- `capability_adapter.go:79-92` — `ProviderCapabilitySet` declara `ProviderCode: "mercado_livre"` com Listings/FeeQuotes/StockReads/Orders/PriceWrites/StockWrites/ListingWrites e âncoras `seller_sku|EAN|title`.
- Endpoints já chamados:
  - `capability_adapter.go:208` — `GET /users/me` (probe de conta).
  - `capability_adapter.go:241` — `GET /users/{id}/items/search?limit&offset&include_filters=false` (paginação offset).
  - `capability_adapter.go:277-287` — `GET /users/{id}/items/search?search_type=scan&scroll_id=...` (`ListListingsScanPage`; comentário :260-265 documenta o cap de 1000 itens do offset).
  - `capability_adapter.go:628` — `GET /items/{id}` (leitura de item, um por vez).
  - `capability_adapter.go:502` — `GET /orders/search?seller&limit&offset&order.status&order.date_last_updated.from`.
  - `capability_adapter.go:540` — `GET /orders/{id}`.
  - `capability_adapter.go:571` — `GET /sites/{site}/listing_prices?price&listing_type_id&category_id` (fee quote).
  - `capability_adapter.go:444` — `PUT /items/{id}` (stock write).
  - `capability_adapter.go:142-152` + comentário :639-641 — `GetShipmentInfo` via shipping_reader; endpoint /shipments/costs exige header `x-format-new: true`.
- Paginação: scan/scroll_id EXISTE (`:266-304`) e o cursor volta pro chamador; offset é fallback. **Hidratação é N+1**: cada item da página de search vira um `ReadListing` → `GET /items/{id}` individual (`:246-255` e `:292-301`). Não há multiget `/items?ids=`.
- `capability_adapter.go:507-520` — ListOrders também é N+1: cada resultado do search dispara `ReadOrder`.
- Headers: `capability_adapter.go:717-731` — `Authorization: Bearer`, `Accept: application/json`, `Content-Type` (se body), `X-Tenant-ID`, `X-Installation-ID`, `X-Idempotency-Key` (writes), mais mapa por request (x-format-new).
- DTOs: tipados, structs privados `mlItemResponse`/`mlOrderResponse`/`mlVariation`/`mlAttribute`/`mlPayment`/`mlListingPriceResponse` (`capability_adapter.go:975-1083`). **Raw JSON não é guardado**: body lido com `io.LimitReader(resp.Body, 1<<20)` (`:747`) e descartado após decode; `RawProviderRef` guarda só o PATH (ex.: `"/items/"+id`, `:362`), e o diagnóstico de erro é clipado a 512 runes (`providerDiag`, `:681-688`).
- 429/retry/backoff: `capability_adapter.go:654-655` — `StatusTooManyRequests` → `ErrCodeProviderRateLimited`, retorno imediato, **uma tentativa, zero retry/backoff**; idem `:462-465` e `:578-583`. Não existe loop de retry em nenhum caminho do adapter.
- Autenticação: `capability_adapter.go:26` — `AccessTokenResolver func(ctx, ProviderAccountRef) (string, error)` injetado; `apps/server_core/internal/composition/root.go:370-378` liga ao `integrationsapp.NewCredentialResolver(credentialSvc, encryptionSvc).ResolveAccessToken`.
- Token refresh: `apps/server_core/internal/modules/integrations/adapters/mercadolivre/auth_adapter.go:223-228` — `RefreshToken` com `grant_type=refresh_token` contra `https://api.mercadolibre.com/oauth/token` (URL em `:73`); ticker de refresh `root.go:661` (`NewRefreshTicker(..., 5*time.Minute)`); backoff EXISTE só aqui: `integrations/domain/refresh_policy.go:18-27` (BackoffBase 30s, BackoffMax 15m) — é para refresh de credencial OAuth, não para chamadas de API.
- Rate limiting de chamadas ML: inexistente. O único semáforo do sistema é do Oracle (`root.go:436`, `oraclebatch.NewSemaphore`).

## 2. Scheduler + sync_state

- `apps/server_core/migrations/0075_sync_sync_state.sql:19-30` — schema exato de `sync_state`: `tenant_id text, installation_id text, entity text, cursor jsonb NULL, schedule jsonb NULL, last_full_sync_at timestamptz NULL, last_incremental_at timestamptz NULL, last_error jsonb NULL, consecutive_failures int NOT NULL DEFAULT 0, PK(tenant_id, installation_id, entity)`.
- `0075_sync_sync_state.sql:12-14` — `entity` é enum SEMÂNTICO (`products|listings|orders|market|tariffs`) validado na aplicação (`domain.Entity.Valid()`), sem CHECK no DB.
- Assinatura de um job novo: `apps/server_core/internal/modules/sync/application/scheduler.go:46` — `type JobFunc func(ctx context.Context, cursor json.RawMessage) (json.RawMessage, error)`; cursor retornado é AUTORITATIVO, nil APAGA o cursor persistido (`:42-45`).
- Registro: `scheduler.go:85-101` — `RegisterJob(entity, fn)`; um job por entity, duplicata rejeitada, entity inválida fail-closed.
- Loop: `scheduler.go:105-119` `Start` (ticker no intervalo injetado, cadence-agnostic), `:124-139` `RunOnce` (falha isolada por entity), `:141-161` `runJob` — erro de leitura de cursor = pula o ciclo (não fabrica nil); `:160` **sempre grava `incremental=false`** hoje.
- StateStore seam: `scheduler.go:22-34` — `Read/RecordSuccess/RecordFailure`, implementado por `sync/adapters/postgres/sync_state_repo.go` (upserts em `:67-78` sucesso, `:95-102` falha com `consecutive_failures+1` atômico; helper `AppendPendingCodigo` `:117-125` faz append em `cursor->pending` via `jsonb_set` num único UPDATE).
- Composition: `apps/server_core/internal/modules/sync/composition/scheduler.go:11` — `InstallationScopeERP = "erp"` (installation_id p/ entidades ERP). `products_job.go:41-56` — `NewProductsScheduler(pool, tenantID, interval, lookup, adapters, opts...)` registra `domain.EntityProducts` com corpo real.
- Bootstrap: `root.go:672-677` — `go synccomposition.NewProductsScheduler(pool, cfg.DefaultTenantID, 15*time.Minute, activeSourceLookup, productSyncAdapters, WithLinkCandidateRefresh(...)).Start(...)` — **products é a ÚNICA entity registrada**; listings/orders/market/tariffs não têm job.
- Interface que o corpo do job usa por baixo: `apps/server_core/internal/modules/internal_read/ports/adapter.go:36-40` — `ProductSourceAdapter interface { Reader; Sync(ctx) (SyncResult, error); Kind() sourcekind.SourceKind }`.
- Outros tickers no boot: `root.go:661` refresh de token 5min, `:662` cleanup OAuth 1h, `:663` fee-sync scheduler 15min (ver §4).

## 3. Migrações

- Diretório: `apps/server_core/migrations/` (arquivos `NNNN_snake_case.sql`, embutidos via `canonical.Source()` — `apps/server_core/cmd/migrate/main.go:31`).
- Aplicação: `apps/server_core/internal/platform/migrate/runner.go:33-40` — `migrate.Run` cria `schema_migrations(filename PK, applied_at)` e aplica em ordem por filename os que faltam; disparado pelo binário `cmd/migrate` (o `cmd/server/main.go` NÃO roda migrações no boot — `main.go:15-55` só monta pool + composition).
- Última migração: **0085**. As 5 últimas: `0081_drop_legacy_marketplaces_and_orphan_tables.sql`, `0082_product_link_decisions.sql`, `0083_sellable_assortment_config.sql`, `0084_products_mirror_sellable_fields.sql`, `0085_erp_import_products_sellable_fields.sql`.
- Anomalias de numeração: `0002` não existe; `0021` está DUPLICADO (`0021_integration_operation_run_evidence.sql` e `0021_integrations_provider_auth_strategy_shopee_partner.sql` — ordenação por filename resolve, mas o próximo autor deve pinçar 0086+); buracos 0040-0042, 0048-0049, 0054, 0059-0064, 0077.
- `migrations/0001_foundation.sql:1` cria `platform_migrations` — tabela legada; o runner usa `schema_migrations`.
- Há testes de migração no próprio diretório (`migrations/listings_test.go`, `migrations/product_link_decisions_test.go`): `:25` é lista de SUBSTRINGS obrigatórias (`required` slice); o assert por regex mora em `createTableBody` (`listings_test.go:101`) — os dois idiomas, medidos (P5 r06 F-r06-5).

## 4. Fees estáticos

- **Seed estático vivo (13/16)**: `apps/server_core/migrations/0068_pricing_tariff_defaults.sql:13-14` — `comissao_classico_pct DEFAULT 13.00`, `comissao_premium_pct DEFAULT 16.00` em `pricing_tariff_defaults`. Materializado on-first-read por `INSERT ... ON CONFLICT DO NOTHING` em `apps/server_core/internal/modules/pricing/adapters/postgres/calc_repository.go:240-246`; lido em `:254-259`.
- Consumo: `apps/server_core/internal/modules/pricing/application/calc_service.go:118-122` — `CalcService.GetTariffDefaults` (tela /precos, simulador); `PutTariffDefaults` `:129` permite edição (upsert `calc_repository.go:286-296`).
- Baseline hardcoded adicional: `apps/server_core/internal/modules/integrations/adapters/mercadolivre/auth_adapter.go:47-48` — metadata do provider com `"baseline_commission_percent": 0.16`, `"baseline_fixed_fee_amount": 0.0`.
- Sobre "16 e 22": não há seed de produção com 22. Os pares 0.16/0.22 aparecem SÓ em testes (`apps/server_core/internal/modules/marketplaces/application/fee_schedule_service_test.go:178-179`); o `22.5` de `pricing/domain/difal_seed.go:34` é alíquota DIFAL do PI, não comissão ML.
- Seeder ML de fee schedule foi REMOVIDO: `apps/server_core/internal/modules/marketplaces/registry/mercado_livre.go:46-48` — `SeedFees` é no-op com comentário "ML fees are read per listing from the Fees API (sale_fee, per unit), never seeded as a flat rate"; `FeeSource: "api_sync"` (`:20`).
- Fee-sync runtime: `root.go:337-338` — `BuildFeeSyncers()` monta os syncers registrados via `RegisterFeeSyncerFactory` (`integrations/adapters/providers/registry.go:38-42,84-98`), **mas nenhum código de produção chama RegisterFeeSyncerFactory** (só testes) → executor sem syncer p/ ML; `feesync/marketplace_executor.go:62-66` retorna `INTEGRATIONS_FEE_SYNC_UNSUPPORTED`. O scheduler `root.go:663` tica a cada 15min contra isso.
- Startup seed genérico (skeleton): `apps/server_core/internal/modules/connectors/application/fee_sync_service.go:22-41` — `SeedAll` idempotente sobre `marketplace_fee_schedules` (tabela em `migrations/0011_marketplace_fee_schedules.sql:1`).
- sale_fee POR PEDIDO (fonte real de fee hoje): `capability_adapter.go:818-837` — soma `order_items[].sale_fee` no snapshot; persiste em `orders_marketplace_order_items.sale_fee_amount` (`0027:32`); consumido por profitability (`apps/server_core/internal/modules/profitability/application/service.go:262-269` vira `InputKindSaleFee`; `:1014-1015` seta flag `MissingSaleFee`; `:1037` contribui na margem).

## 5. listings hoje

- Schema: `apps/server_core/migrations/0036_listings.sql:2-31` — `listings(tenant_id, installation_id, provider, provider_listing_id, variation_id, title, listing_type_code, status, price_amount, price_currency, published_quantity, sync_state, sync_error jsonb, quality_score, sales_30d, fetched_at, created_at, updated_at)`, PK `(tenant_id, installation_id, provider_listing_id, variation_id)`; CHECKs de status/sync_state/quality/par-de-preço.
- `0036_listings.sql:36-48` — tabela irmã `listing_sync_events` (timeline por listing).
- Alteração: `0037_listings_status.sql:3-4` — status ampliado p/ `under_review|inactive|payment_required|not_yet_active`.
- **`listing_variations` como tabela: NÃO EXISTE.** Variação é linha da própria `listings` via `variation_id` no PK (sentinela `'-'` p/ sem-variação — comentário `0036:1`); no domínio do connector existe `domain.ListingVariationSnapshot` populado em `capability_adapter.go:786-794`; identidade por variação (SKU/EAN) persiste em `product_link_listing_snapshots` (`migrations/0022_product_links_listing_snapshots.sql:1-17`, PK com `provider_variation_id`).
- Quem escreve `listings`: único escritor de produção é `apps/server_core/internal/modules/listings/adapters/postgres/repository.go:383-465` `ApplyCompletedPull` — **fecha TUDO (`UPDATE ... SET status='closed'` sem filtro além de tenant/installation, `:390-394`) e re-upserta o que o pull devolveu (`:426-440`, ON CONFLICT no PK)**. É o MASS-CLOSURE do audit D-120: pull parcial fecha o catálogo inteiro.
- Pipeline: `listings/application/ingestion.go:34-76` — `Ingestion.Pull` pagina via `ports.PageSource` até página curta (teto 10.000 páginas `:21`), acumula TUDO em memória e só então aplica; `listings/adapters/connectors/source.go:54-89` — `Source.ReadPage` usa `ListListingsScanPage` quando disponível e alimenta o observer do product_links (`:17-19`, `AbsorbProviderSnapshots` — resolve o chicken-egg de EAN/SKU antes do achatamento).
- Disparo: manual, não agendado — `root.go:729-736` monta `NewIngestion(source, listingRepo, 100, time.Now)` + `NewRefreshService` exposto por `listingstransport.NewRefreshHandler` (HTTP). Nenhum job `EntityListings` no scheduler.

## 6. orders hoje

- Schema: `migrations/0027_orders_marketplace_orders.sql:1-19` — `orders_marketplace_orders` PK `(tenant_id, installation_id, provider_order_id)` com `provider_status`, timestamps do provider, `fetched_at`, `shipping_id text` (`:12`), `tags_json`, `raw_provider_ref jsonb NULL`; `:21-38` `orders_marketplace_order_items` (com `sale_fee_amount`, `link_quality`, `internal_product_id`); `:40-52` `orders_marketplace_order_payments`.
- Alterações: `0033_orders_sankhya_linkage.sql:2-18` — `mpc_line_id` (imutável por trigger `:20-34`) + `reconciliation_state` nos items; `0074_orders_faturado_at.sql:7` — `faturado_at timestamptz`; `0079_orders_buyer_nickname.sql:6` — `buyer_nickname`.
- Como entram: **import manual via HTTP, não é read-through vivo nem agendado.** `orders/transport/http_handler.go:85` registra `POST /orders/import` → `:96-121` chama `ImportService.Import`; `orders/application/import_service.go:46-82` — `source.ListOrders` (live contra ML naquele momento, limit default 20 `:56`), resolve links, `store.UpsertOrders`. A LISTA (`GET /orders`) lê só do Postgres: `orders/application/list_service.go:52-57` (`ReadService.List` → `OrderReadStore`).
- Source: `orders/adapters/integrations/order_source.go:22-70` mapeia `connectorsdomain.OrderSnapshot` → `OrderIngestionSnapshot` (inclui `SaleFeeAmount`, `ShippingID`, `BuyerNickname`).
- Upsert: `orders/adapters/postgres/order_repo.go:614` — `ON CONFLICT (tenant_id, installation_id, provider_order_id) DO UPDATE`; items `:691` — `ON CONFLICT (..., mpc_line_id) DO UPDATE`.
- **pack_id: não existe no schema nem no módulo orders.** Única ocorrência é a sonda `apps/server_core/cmd/mlprobe/main.go:358-405` (T3-pack explorando `/packs/{id}`).
- Shipment: só `shipping_id` na row (`0027:12`) + leitura viva opcional `GetShipmentInfo` (`capability_adapter.go:142-152`, porta `orders/ports/shipment_reader.go`); nenhuma tabela de shipments.

## 7. Vínculo + mirror

- Tabela de vínculo: `product_links` — `migrations/0025_product_link_workflows.sql:1-15`: PK `(tenant_id, installation_id, provider_item_id, provider_variation_id)`, `state text`, `internal_product_id integer` (CODPROD), `internal_product_name`, `internal_reference_code`; auditoria em `product_link_audit_entries` (`:20-39`).
- Decisões (D-121): `migrations/0082_product_link_decisions.sql:25-63` — append-only, `link_id` GENERATED da chave natural (`:35-37`), `rule_matched IN (exact_codprod_unique, exact_ean_unique, concordant_codprod_ean, manual)` (`:50`), CHECK `actor='system' ⇒ rule='concordant_codprod_ean'` (`:53-54`).
- Candidatos: `0023_product_link_candidates.sql`, confidence `0065`, batches `0066`, audit batch_id `0067`.
- Mirror: `migrations/0076_products_mirror_active_source.sql:18-39` — `products_mirror(tenant_id, source, codigo_produto, descricao, referencia, ean, marca, grupo_*, ncm, custo NUMERIC NULL, preco_venda NUMERIC NULL, estoque_total NUMERIC NULL, protocol_id, absent_in_last_snapshot, stale_since, updated_at)`; `0078:19-31` re-keyed PK `(tenant_id, source, codigo_produto)`; `0078:69-73` adiciona `estoque_fisico`/`estoque_reservado`; `0084:5-9` adiciona `usoprod`, `ad_ecommerce`. Stock por local: `products_mirror_stock_locations` (`0076:44-51`, + source no PK via `0078:49-61`). Config: `active_source` (`0076:57-67`) + flags de sortimento `only_revenda/only_em_estoque/only_ecommerce_eligible` (`0083:6-13`).
- Disponível vendável — computado em DOIS lugares:
  - Lado Sankhya (produtor): `apps/server_core/internal/modules/internal_read/adapters/oracle/sync.go:257-265` — `sankhyaSellableLocationCodes = "10101, 10102"` e Q4 `SUM(NVL(ESTOQUE,0) - NVL(RESERVADO,0))` com `CODPARC=0 AND CODEMP IN (1,2) AND CODLOCAL IN (10101,10102)`; USOPROD/AD_ECOMMERCE vêm da Q1 (`:104-111`); CODEMP=1 pinado `:17-20`.
  - Lado leitura (telas): `apps/server_core/internal/modules/erp_import/adapters/postgres/mirror_query_repository.go:133-145` — `mirrorAssortmentPredicate`: `usoprod IS NULL OR usoprod='R'`, `estoque_total IS NULL OR estoque_total > 0`, `ad_ecommerce IS NULL OR <> 'N'` (NULL passa = honest-unknown incluído); aplicado no catálogo `:147-160` e nas contagens `:180-192`.

## 8. Padrão de upsert/idempotência (ERP→mirror)

- Exemplo canônico: `apps/server_core/internal/modules/internal_read/adapters/mirror/writer.go:74-95` — `INSERT INTO products_mirror (...) VALUES ($1,'sankhya',...) ON CONFLICT (tenant_id, source, codigo_produto) DO UPDATE SET <todas as colunas de fato> , absent_in_last_snapshot=false, stale_since=NULL, updated_at=now()`.
- Keep-absent (ADR-04): `writer.go:104-112` (doc comment em `:97-103` — range medido, P5
  r07 F-r07-4) — `keepAbsentSQL` marca `absent_in_last_snapshot=true` + `stale_since=COALESCE(...)` SÓ nas rows da própria source; nunca DELETE.
- Outros upserts por chave natural no mesmo padrão: `sync/adapters/postgres/sync_state_repo.go:67-78` (sync_state, COALESCE nos timestamps p/ não regredir); `orders/adapters/postgres/order_repo.go:614` e `:691`; `listings/adapters/postgres/repository.go:432`; `pricing/adapters/postgres/calc_repository.go:290` (tariff defaults DO UPDATE); `product_links/adapters/postgres/listing_snapshot_repo.go:33` (snapshots de identidade).

## Lacunas (o que NÃO existe hoje)

- **Retry/backoff para 429 da API ML: inexistente** — 429 vira `ErrCodeProviderRateLimited` e a chamada morre na primeira tentativa (`capability_adapter.go:654-655`). O único backoff do sistema é do refresh de credencial OAuth (`refresh_policy.go:18-27`).
- **Rate limiter para chamadas ML: inexistente** (semáforo existe só para Oracle, `root.go:436`).
- **Multiget `/items?ids=`: inexistente** — hidratação de listings e orders é N+1 (1 GET por item/pedido). (D-120: multiget de shipments nem existe na API.)
- **Persistência do payload raw da ML: inexistente** — body limitado a 1MB, decodificado e descartado; `raw_provider_ref` guarda só o path do recurso.
- **Webhook/notifications receiver ML: inexistente** — únicas rotas de callback são OAuth (`integrations/transport/auth_handler.go:48`, melhor-envio); `Webhooks: CapabilitySupported` em `registry/mercado_livre.go:30` é declaração de capacidade do provider, sem handler.
- **Jobs de sync para `listings`/`orders`/`market`/`tariffs`: inexistentes** — só `EntityProducts` registrado (`root.go:672-677`); entidades previstas no enum da 0075 estão órfãs.
- **Uso de `last_incremental_at`: inexistente** — o scheduler grava sempre `incremental=false` (`scheduler.go:160`); `schedule jsonb` da 0075 nunca é lido pelo loop (intervalo vem hardcoded da composition: 15min).
- **Tabela `listing_variations`: inexistente** (variação achatada no PK de `listings`).
- **pack_id / tabela de shipments: inexistentes** (pack_id só no mlprobe).
- **Orders agendado/read-through vivo: inexistente** — só `POST /orders/import` manual, limit default 20.
- **FeeScheduleSyncer de produção para ML: inexistente** — `RegisterFeeSyncerFactory` sem chamador fora de testes; o fee-sync scheduler de 15min executa e retorna UNSUPPORTED para ML; fee real entra por pedido (`sale_fee`) e por quote on-demand (`/sites/MLB/listing_prices`).
- **Listing pull incremental: inexistente** — todo pull é full-catalog em memória + MASS-CLOSURE (`repository.go:390-394`), sem `date_last_updated` filter para items (só orders search tem `order.date_last_updated.from`, `capability_adapter.go:494`).
