# P5 Prerequisites — MIS-007 ml-sync feature-brief facts

Status: complete. All 12 items located and verified by direct file read on main (tip dd89d4b3 at investigation time).
Convention: `path:line` line numbers are 1-based as of this read. All identifiers/values verbatim.

## 1. DeriveOrderBucket

- File: `apps/server_core/internal/modules/orders/domain/order_bucket.go`
- Signature (`:48`):
  `func DeriveOrderBucket(providerStatus, shipmentStatus string, tags []string, faturado bool) OrderBucket`
- Type + every literal it can return (`:9-17`):
  ```go
  type OrderBucket string

  const (
      BucketNovo      OrderBucket = "novo"
      BucketFaturar   OrderBucket = "faturar"
      BucketEnviar    OrderBucket = "enviar"
      BucketEnviado   OrderBucket = "enviado"
      BucketCancelado OrderBucket = "cancelado"
  )
  ```
- Inputs consumed (decision order, `:48-84`):
  1. `providerStatus` lowercased/trimmed; `"cancelled"` or `"invalid"` → `BucketCancelado` (`:53-55`)
  2. order tag `"delivered"` (case-insensitive, via `hasTag`, `:58-60`) → `BucketEnviado`
  3. live `shipmentStatus` `"shipped"` or `"delivered"` (`:66-69`) → `BucketEnviado`
  4. `providerStatus` `"paid"`/`"confirmed"`: `faturado==true` → `BucketEnviar`, else `BucketFaturar` (`:73-78`)
  5. `providerStatus` `"delivered"`/`"shipped"` → `BucketEnviado` (`:79-80`); `default` → `BucketNovo` (`:81-82`)
- No dates are consumed. `faturado` is passed as `e.Order.FaturadoAt != nil` at the transport call site.
- Call sites: `orders/transport/http_handler.go:549` (`Bucket: domain.DeriveOrderBucket(e.Order.Status, shipmentStatusOf(e.Shipment), e.Order.Tags, e.Order.FaturadoAt != nil)`) and `orders/adapters/postgres/order_repo.go:378` (`DeriveOrderBucket(providerStatus, "", tags, faturadoAt.Valid)` — summary counts path, no shipment status).
- Exhaustive truth table test: `orders/domain/order_bucket_test.go:8` (`TestDeriveOrderBucket`).

## 2. Buyer fiscal

### Adapter (mercadolivre)
- File: `apps/server_core/internal/modules/connectors/adapters/mercado_livre/buyer_fiscal_reader.go`
- Entry: `func (a *CapabilityAdapter) GetBuyerFiscalInfo(ctx context.Context, accountRef domain.ProviderAccountRef, providerOrderID string) (domain.BuyerFiscalInfo, error)` (`:59`)
- Two-step flow: `GET /orders/{id}` → `buyer.billing_info.id` → `GET /orders/billing-info/{SITE_ID}/{billing_info_id}` (`:71-107`). 404/undecodable payload degrades to empty `BuyerFiscalInfo{FetchedAt: a.now()}` (honest absence), never an error (`:101-103`).
- Provider decode structs (`:18-52`): `mlBillingInfoResponse{Buyer}` → `mlBillingInfoBuyer{BillingInfo}` → `mlBillingInfoData{Name, LastName, Identification, Address}` → `mlBillingIdentification{Type, Number}` (Type is OPAQUE, never mapped to CPF/CNPJ enum) and `mlBillingAddress{StreetName, StreetNumber, CityName, State{Code,Name}, ZipCode, CountryID}`.

### Returned domain DTO
- File: `apps/server_core/internal/modules/connectors/domain/buyer_fiscal.go`
- `type BuyerFiscalInfo struct` (`:16-22`):
  ```go
  Name      *string             `json:"name"`
  DocType   *string             `json:"doc_type"`
  DocNumber *string             `json:"doc_number"`
  Address   *BuyerFiscalAddress `json:"address"`
  FetchedAt time.Time           `json:"fetched_at"`
  ```
- `type BuyerFiscalAddress struct` (`:27-35`):
  ```go
  StreetName   *string `json:"street_name"`
  StreetNumber *string `json:"street_number"`
  City         *string `json:"city"`
  StateCode    *string `json:"state_code"`
  StateName    *string `json:"state_name"`
  ZipCode      *string `json:"zip_code"`
  CountryID    *string `json:"country_id"`
  ```
- `func (b BuyerFiscalInfo) HasData() bool` (`:40-42`).
- Orders-local consumer port: `apps/server_core/internal/modules/orders/ports/buyer_fiscal_reader.go:19-21` — `GetBuyerFiscal(ctx context.Context, installationID, providerOrderID string) (connectorsdomain.BuyerFiscalInfo, error)`.

### Handler mapping into response
- File: `apps/server_core/internal/modules/orders/transport/http_handler.go`
- Route: `mux.HandleFunc("/orders/{provider_order_id}", h.handleGet)` (`:89`).
- Response field on `enrichedOrderDTO` (`:514`): `CompradorFiscal *compradorFiscalDTO \`json:"comprador_fiscal,omitempty"\``
- Assigned at `:569`: `dto.CompradorFiscal = mapCompradorFiscal(e.BuyerFiscal)`; mapper `mapCompradorFiscal` `:578-588`, `mapCompradorFiscalEndereco` `:593-606`.
- `type compradorFiscalDTO struct` (`:462-467`): `Nome *string json:"nome,omitempty"`, `DocTipo *string json:"doc_tipo,omitempty"`, `DocNumero *string json:"doc_numero,omitempty"`, `Endereco *compradorFiscalEnderecoDTO json:"endereco,omitempty"`.
- `type compradorFiscalEnderecoDTO struct` (`:446-454`): `logradouro`, `numero`, `cidade`, `uf_codigo`, `uf_nome`, `cep`, `pais` (all `*string`, `omitempty`).
- Enrichment production: `orders/application/enrich_service.go:194-198` (`EnrichOne`, BuyerFiscal enrichment seam; degrade semantics no doc comment `:188-193` — range medido, sweep residual F-r08-3).

### FE renderer
- Component: `CompradorFiscalSection` in `apps/web/src/pages/pedidos/PedidoDrawer.tsx:355-368`, mounted in `DrawerBody` at `:375`. Section title string: `"Comprador · fiscal (ERP)"`.
- Fields displayed: `cf.nome` (label "Nome/Razão"), `cf.doc_tipo` + `cf.doc_numero` joined with a space (label "Documento", mono), and `endereco` composed by `formatEndereco` (`:341-347`) from `logradouro`, `numero` (joined ", "), `cidade`/`uf_codigo` (joined "/"), `cep`, `pais` (joined " · ").
- NOT displayed: `uf_nome`, `fetched_at`, `street state_name`.
- SDK type: `comprador_fiscal?: OrderCompradorFiscal` on OrderRead, `packages/sdk-runtime/src/index.ts:784`.

## 3. GET /listings

### Handler + response
- File: `apps/server_core/internal/modules/listings/transport/http_handler.go`
- Registration (`:137-148`): `mux.HandleFunc("GET /listings", h.HandleList)` (`:144`); also `GET /listings/by-product` (`:145`), `GET /listings/summary` (`:146`), `GET /listings/{id}` (`:147`). All registered `httpx.InteractiveRouteClass` (`:139-141`).
- `HandleList` body: `:168-226`. Page envelope (`:150-155`):
  ```go
  type listingPageEnvelope struct {
      Items      any     `json:"items"`
      NextCursor *string `json:"next_cursor"`
      PageSize   int     `json:"page_size"`
      AsOf       any     `json:"as_of"`
  }
  ```
- Per-listing item = `domain.ListingReadModel`, `apps/server_core/internal/modules/listings/domain/read_model.go:122-149`. JSON field names verbatim: `listing_id`, `installation_id`, `provider`, `provider_listing_id`, `variation_id`, `title`, `listing_type`, `status`, `link`, `price`, `published_quantity`, `sync_state`, `sync_error`, `quality_score`, `pending_issue`, `sales_30d`, `cost`, `below_margin_worst_case`, `icms_worst_case_by_uf`, `fetched_at`, `market_signal`, `signal_status`.
- Pagination/query params (`transport/query.go:28-67`): `installation_id` (required, `http_handler.go:175-179`), `cursor` (opaque, `ports.DecodeListingCursor`, `:181-192`), `q`, `limit` (default 50, valid 1..200, `query.go:29,47-53`), `filter.<key>` with allowed keys `domain.FilterKeys = []string{"status", "sync_state", "link_state", "exception", "has_exception", "listing_type_code", "product_id"}` (`domain/filter.go:9`).
- Sort order (`adapters/postgres/repository.go:128`): `ORDER BY l.title,l.provider_listing_id,l.variation_id LIMIT $n`. by-product groups: `ORDER BY null_last,product_title,product_id` (`:170`); group children `ORDER BY l.title,l.provider_listing_id,l.variation_id` (`:221`); timeline `ORDER BY at DESC,event_id DESC` (`:363`).

### FE /anuncios
- Route: `apps/web/src/routes/anuncios.tsx`; page: `apps/web/src/pages/AnunciosPage.tsx`; table: `apps/web/src/pages/AnunciosTable.tsx`; queries: `apps/web/src/pages/anunciosQueries.ts`; URL state: `apps/web/src/pages/anunciosQueryState.ts`; detail panel: `apps/web/src/pages/ListingDetailPanel.tsx`; summary strip: `apps/web/src/pages/ListingsSummary.tsx`.
- Current columns (`AnunciosTable.tsx:277-294`): selection checkbox column, then `MLB`, `TÍTULO`, `PRODUTO`, `PREÇO`, `EST.`, `SYNC`, `QUAL.`, `PENDÊNCIA`. Supports agrupar-por-produto group-header rows (`renderGroupHeader`, `:223`) and collapse state (`:148-171`).

## 4. IntegracoesPage

- File: `apps/web/src/pages/integracoes/IntegracoesPage.tsx` (+ `IntegracoesPage.test.tsx`, `useErpImportUpload.ts` in same dir).
- Page component `IntegracoesPage` (`:558-574`), h1 = "Configuração da plataforma". Section order:
  1. `ActiveSourceCard` (`:308-379`) — radio "Fonte ativa": options `sankhya` / `xlsx` / `catalogo_cliente` (`ACTIVE_SOURCE_OPTIONS`, `:291-295`).
  2. `SellableAssortmentCard` (`:403-487`) — checkboxes `only_revenda` / `only_em_estoque` / `only_ecommerce_eligible` (`:381-401`).
  3. `UploadCard` (`:142-283`) — .xlsx import, source radio `catalogo_cliente` | `xlsx` (`SOURCE_OPTIONS`, `:33-44`), renders `ResultSummary` (`:99-140`).
  4. `ProviderConnectCard` (`:496-556`) — Mercado Livre OAuth connect.
  5. `ImportacaoSection` — imported from `../importacoes/ImportacaoSection` (`:19`).
- Data fetching:
  - Hooks from `@marketplace-central/web-query` (`:10-16`): `useActiveSourceQuery`, `useCatalogAssortmentCountsQuery`, `useSellableAssortmentQuery`, `useSetActiveSourceMutation`, `useSetSellableAssortmentMutation`.
  - `useErpImportDetail` from `../vinculos/useErpImports` (`:20`, used at `:100`); `useErpImportUpload` (`:150`).
  - Direct SDK client calls (via `useClient()`): `client.listIntegrationInstallations()` (`:508`), `client.createIntegrationInstallation({installation_id, provider_code: "mercado_livre", family: "marketplace", display_name})` (`:520-525`), `client.startIntegrationAuthorization(installationId)` (`:527`).
  - Fail-closed source-unset detection: `isApiError(error) && error.status === 400 && hasCode(error, "unknown_erp_source")` (`:304-306`).

## 5. Orders import path

- File: `apps/server_core/internal/modules/orders/application/import_service.go`
- Input: `type ImportOrdersInput struct { InstallationID string; Limit int }` (`:14-17`); `Limit <= 0` defaults to `20` (`:54-57`).
- Constructor: `func NewImportService(cfg ImportServiceConfig) *ImportService` (`:33`); config = `{Source ports.OrderSource; Links ports.LinkReader; Store ports.OrderStore; Now func() time.Time}` (`:26-31`).
- Entry: `func (s *ImportService) Import(ctx context.Context, input ImportOrdersInput) (domain.ImportResult, error)` (`:46-82`). Flow: `s.source.ListOrders(ctx, installationID, limit)` (`:59`) → `s.links.ResolveLinks(ctx, installationID, identities)` (`:65`) → `normalizeOrders(...)` (`:70`) → `s.store.UpsertOrders(ctx, orders)` (`:71`).
- What it persists today (`normalizeOrders`, `:108-162`): `domain.MarketplaceOrder` header (`InstallationID, ProviderCode, ProviderOrderID, ProviderStatus, ProviderStatusDetail, ProviderCreatedAt, ProviderClosedAt, ProviderUpdatedAt, FetchedAt, ShippingID, BuyerNickname, CancellationDetail, Tags, RawProviderRef, CreatedAt, UpdatedAt`) + `Payments []MarketplaceOrderPayment{ProviderPaymentID, ProviderStatus, TransactionAmount, TotalPaidAmount}` + `Items []MarketplaceOrderItem{ProviderItemID, ProviderVariationID, SellerSKU, Title, Quantity, UnitPrice, SaleFeeAmount, LinkQuality, InternalProductID}`. `RawProviderRef` = `"/orders/"+PathEscape(id)` only for `mercado_livre` (`:164-173`).
- Store port: `UpsertOrders(ctx context.Context, orders []domain.MarketplaceOrder) (importedCount int, skippedCount int, err error)` — `orders/ports/order_store.go:32`.
- Repository upsert: `func (r *OrderRepository) UpsertOrders(ctx context.Context, orders []ordersdomain.MarketplaceOrder) (int, int, error)` — `orders/adapters/postgres/order_repo.go:467-514`. Per order: `r.upsertOrder` (header wins only when `provider_updated_at` bumped; `INSERT INTO orders_marketplace_orders` at `:603`); on win → `r.replaceItems` (`:501`, `INSERT INTO orders_marketplace_order_items` at `:682`) + `r.replacePayments` (`:504`, `INSERT INTO orders_marketplace_order_payments` at `:772`); on lose → `r.backfillMissingLines` (`:492`, defined `:519` — writes lines only if the stored order has none).
- HTTP entry: `mux.HandleFunc("/orders/import", h.handleImport)` — `orders/transport/http_handler.go:85`; `/orders/import` is registered `BatchRouteClass` in `registerBatchRoutes` (root.go:263).

## 6. Pricing tariff chain

### Live resolver (degrau 3)
- File: `apps/server_core/internal/modules/pricing/adapters/tarifflive/resolver.go` (import alias in root.go: `pricingtarifflive "marketplace-central/apps/server_core/internal/modules/pricing/adapters/tarifflive"`, root.go:101).
- Constructor (`:35`): `func NewResolver(identities ports.ProductIdentityReader, categories ports.CategoryResolver, commissions ports.CommissionQuoter) *Resolver`
- Method (`:43`): `func (r *Resolver) ResolveCommission(ctx context.Context, req ports.TariffRequest) (domain.ComponentResolution, bool, error)` — chains identity → category → commission; every miss returns `found=false`. Constants: `listingTypeClassico = "gold_special"`, `listingTypePremium = "gold_pro"` (`:20-21`), `priceBasisCurrency = "BRL"` (`:26`). Successful quote → `domain.ComponentResolution{Valor, Fonte: domain.FonteCotacao, Degrau: 3, Data}` (`:79-84`).
- NOTE: `tarifflive.Resolver` does NOT implement `ports.TariffResolver`; it satisfies the composite's private `commissionResolver` seam.

### Composite
- File: `apps/server_core/internal/modules/pricing/adapters/tariffcomposite/resolver.go`
- `type Resolver struct { base ports.TariffResolver; live commissionResolver }` (`:23-26`); `var _ ports.TariffResolver = (*Resolver)(nil)` (`:28`); `func NewResolver(base ports.TariffResolver, live commissionResolver) *Resolver` (`:32`); `func (r *Resolver) Resolve(ctx context.Context, req ports.TariffRequest) (domain.TariffResolution, error)` (`:41`) — degrau-4 error propagates; degrau-3 miss OR error silently keeps degrau-4 commission.

### Port
- File: `apps/server_core/internal/modules/pricing/ports/tariff_resolver.go`
- `type TariffRequest struct { Modalidade domain.Modalidade; ProductID *int; PriceBasis *string }` (`:14-24`).
- `type TariffResolver interface { Resolve(ctx context.Context, req TariffRequest) (domain.TariffResolution, error) }` (`:30-32`).

### root.go wiring (verbatim)
- `apps/server_core/internal/composition/root.go:837-856`:
  ```go
  calcTariffResolver := pricingtariffdefaults.NewResolver(calcRepo, cfg.DefaultTenantID, "")     // :837
  var tariffResolver pricingports.TariffResolver = calcTariffResolver                            // :844
  if feeReader, ferr := marketplaceCapabilities.FeeQuoteReader(mercadoLivreProviderCode); ferr == nil {  // :845
      identityReader := pricingcatalog.NewProductIdentityReader(catalogSvc)                      // :846
      categoryResolver := newPricingCategoryResolverAdapter(mercadoLivreCapabilities, installationSvc, cfg.DefaultTenantID)  // :847
      commissionQuoter := newPricingCommissionQuoterAdapter(feeReader, installationSvc, cfg.DefaultTenantID)                 // :848
      liveResolver := pricingtarifflive.NewResolver(identityReader, categoryResolver, commissionQuoter)                      // :849
      tariffResolver = pricingtariffcomposite.NewResolver(calcTariffResolver, liveResolver)      // :850
  }
  batchOrch.WithTariffResolver(tariffResolver)                                                   // :852
  calcSvc := pricingapp.NewCalcService(calcRepo, calcCost, calcProducts, cfg.DefaultTenantID).
      WithTariffStore(calcRepo).
      WithTariffResolver(tariffResolver)                                                         // :853-855
  ```
  Whole block gated by `if internalReadAvailable {` (`:829`); `calcRepo := pricingpostgres.NewCalcRepository(pool)` (`:830`).

### pricing_tariff_defaults read path
- File: `apps/server_core/internal/modules/pricing/adapters/postgres/calc_repository.go`
- `func (r *CalcRepository) GetTariffDefaults(ctx context.Context, tenantID, installationID string) (domain.TariffDefaults, error)` (`:239-269`). Materialize-on-first-read: `INSERT INTO pricing_tariff_defaults (tenant_id, installation_id) VALUES ($1, $2) ON CONFLICT (tenant_id, installation_id) DO NOTHING` (`:240-244`; DB column DEFAULTs 13.00/16.00/sem_dados per comment `:236`), then `SELECT comissao_classico_pct::text, comissao_premium_pct::text, frete_estimativa_amount::text, frete_policy FROM pricing_tariff_defaults WHERE tenant_id = $1 AND installation_id = $2` (`:254-259`).
- `func (r *CalcRepository) UpsertTariffDefaults(ctx context.Context, tenantID, installationID string, in domain.TariffDefaults) (domain.TariffDefaults, error)` (`:274-310`), `ON CONFLICT (tenant_id, installation_id) DO UPDATE` (`:290`).

### PutTariffDefaults endpoint
- File: `apps/server_core/internal/modules/pricing/transport/calc_handler.go`
- Routes: `mux.HandleFunc("GET /pricing/tariff-defaults", h.handleGetTariffDefaults)` (`:37`), `mux.HandleFunc("PUT /pricing/tariff-defaults", h.handlePutTariffDefaults)` (`:38`).
- `handlePutTariffDefaults` (`:461-478`) → `h.calc.PutTariffDefaults(r.Context(), tariffInstallationID(r), domain.TariffDefaults{...})`; `tariffInstallationID` reads optional `installation_id` query param, `""` = single-installation sentinel (`:448-450`).
- DTO (`:430-435`): `comissao_classico_pct` (string), `comissao_premium_pct` (string), `frete_estimativa_amount` (*string), `frete_policy` (string).
- Service: `func (s CalcService) PutTariffDefaults(ctx context.Context, installationID string, in domain.TariffDefaults) (domain.TariffDefaults, error)` — `pricing/application/calc_service.go:129`.
- SDK: there is NO client method for /pricing/tariff-defaults in `packages/sdk-runtime/src/index.ts`; only the interface `PricingTariffDefaults` exists (`index.ts:1600`).

## 7. root.go wiring anchors (current exact ranges)

File: `apps/server_core/internal/composition/root.go`

- registerBatchRoutes: `:259-272` — patterns `/profitability/margin-inputs/import`, `/profitability/profit-snapshots/calculate`, `/orders/import`, `/product-links/listing-snapshots/imports`, `/product-links/link-candidates/generations`, `/pricing/simulations/batch`, `/admin/fee-schedules/sync`, `/admin/fee-schedules/seed` → `mux.RegisterRouteClass(pattern, httpx.BatchRouteClass)`. Called at `:282` inside `NewRootRuntime` (`:274`).
- Orders wiring incl. enrich: `:576-601`. Constructors: `orderspostgres.NewOrderRepository` (`:576`), `ordersproductlinks.NewLinkReader` (`:580`), `ordersapp.NewImportService(ordersapp.ImportServiceConfig{Source: ordersintegrations.NewOrderSource(providerOperationSvc), Links: ordersLinkReader, Store: ordersRepo})` (`:581-585`), `ordersapp.NewListService` (`:586`), `orderspostgres.NewOrderReadRepository` (`:587`), `ordersapp.NewReadService` (`:588`), `ordersapp.NewSummaryServiceWithBuckets` (`:589`), `newOrdersCostReaderAdapter` (`:590`), `newOrdersShipmentReaderAdapter` (`:591`), `newOrdersBuyerFiscalReaderAdapter` (`:592`), `orderspricingtax.NewReader(pricingpostgres.NewCalcRepository(pool), cfg.DefaultTenantID)` (`:596`), `ordersapp.NewEnrichServiceWithReaders(ordersCostReader, ordersShipmentReader, nil, ordersBuyerFiscalReader, slog.Default()).WithLinkRefresh(ordersLinkReader).WithTaxes(ordersTaxReader)` (`:597-599`), `ordersapp.NewFaturadoService` (`:600`), `orderstransport.NewHandlerWithSummary(...).WithFaturador(ordersFaturadoSvc).Register(mux)` (`:601`).
- Background schedulers: `:661-678`. `integrationsbg.NewRefreshTicker(authSessionRepo, authFlowSvc, 5*time.Minute)` (`:661`), `integrationsbg.NewStateCleanup(oauthStateRepo, time.Hour)` (`:662`), `integrationsbg.NewFeeSyncScheduler(installationSvc, providerSvc, feeSyncSvc, 15*time.Minute)` (`:663`). Products sync (`:668-678`): `if activeSourceLookup != nil { go synccomposition.NewProductsScheduler(pool, cfg.DefaultTenantID, 15*time.Minute, activeSourceLookup, productSyncAdapters, synccomposition.WithLinkCandidateRefresh(productlinkscomposition.NewLinkCandidateRefresher(installationSvc, productLinkGenerationSvc))).Start(context.Background()) }`.
- Listings read wiring: `:701-714`. `listingspostgres.NewRepository(pool, cfg.DefaultTenantID)` (`:701`), `listingsinternalread.NewCostReader(internalreadoracle.NewBatchReader(oracleDB, oracleBatchSemaphore))` (`:702`), `listingsmarketplaces.NewPolicyReader(marketSvc)` (`:703`), `listingsintegrations.NewInstallationReader(installationSvc)` (`:704`), `newListingsEvidenceAdapter(marketEvidenceSvc)` (`:705`), `listingsapp.NewReadServiceWithEvidence(listingRepo, listingCostReader, listingPolicyReader, listingInstallationReader, time.Now, listingEvidenceReader)` (`:706-713`), `listingstransport.NewReadHandler(listingSvc).Register(mux)` (`:714`).
- Listings refresh/ingestion wiring: `:728-736`. `listingsconnectors.NewSourceWithObserver(marketplaceCapabilities, productLinkImportSvc)` (`:728`), `listingsapp.NewIngestion(listingIngestionSource, listingRepo, 100, time.Now)` (`:729`), `listingsintegrations.NewGateway(installationSvc, operationSvc)` (`:730`), `listingsapp.NewRefreshService(listingRefreshGateway, listingIngestion, func(task func()) { go task() }, time.Now, func(err error) { slog.Error(...) }, func(err error) string { return string(connectorsdomain.ErrorCodeOf(err)) })` (`:731-735`), `listingstransport.NewRefreshHandler(listingRefreshSvc).Register(mux)` (`:736`).
- Pricing: block `:811-858` (repo/service `:811-812`, batch orchestrator `:820-822`, calc + tariff chain `:828-857` — verbatim in item 6), `pricingHandler.Register(mux)` (`:858`).

## 8. Listings repository / ingestion / source

### ApplyCompletedPull
- File: `apps/server_core/internal/modules/listings/adapters/postgres/repository.go`
- Signature (`:383`): `func (r *Repository) ApplyCompletedPull(ctx context.Context, installationID string, rows []domain.Listing, completedAt time.Time) error`
- Mass-closure UPDATE verbatim (`:390-394`):
  ```sql
  UPDATE listings SET status = 'closed', updated_at = $3
  WHERE tenant_id = $1 AND installation_id = $2
  RETURNING provider_listing_id, variation_id
  ```
  (closes EVERY listing of the installation first; rows in the pull are then re-upserted at `:426-440` with `ON CONFLICT (tenant_id, installation_id, provider_listing_id, variation_id) DO UPDATE`; keys previously present but absent from the pull get a `"closed"` / `"Anúncio ausente no provedor — fechado"` timeline event, `:453-460`; all inside one tx, commit `:461`.)
- Timeline event insert helper `insertEvent` (`:467-480`) → `INSERT INTO listing_sync_events`.

### Ingestion.Pull
- File: `apps/server_core/internal/modules/listings/application/ingestion.go`
- `type Ingestion struct { source ports.PageSource; store ports.CompletedPullStore; pageSize int; now func() time.Time }` (`:23-28`); `const maxIngestionPages = 10_000` (`:21`).
- `func NewIngestion(source ports.PageSource, store ports.CompletedPullStore, pageSize int, now func() time.Time) *Ingestion` (`:30`).
- `func (i *Ingestion) Pull(ctx context.Context, account ports.InstallationAccount) error` (`:34`).

### AbsorbProviderSnapshots
- Declaration: `apps/server_core/internal/modules/product_links/application/import_service.go:84`:
  `func (s *ImportService) AbsorbProviderSnapshots(ctx context.Context, installationID string, listings []connectorsdomain.ListingSnapshot) error`
- Consumer interface + caller: `apps/server_core/internal/modules/listings/adapters/connectors/source.go` — `type SnapshotObserver interface { AbsorbProviderSnapshots(ctx context.Context, installationID string, snapshots []connectorsdomain.ListingSnapshot) error }` (`:17-19`); called from `(Source).observe` at `:44` (`return s.observer.AbsorbProviderSnapshots(ctx, account.InstallationID, snapshots)`); `observe` invoked from `ReadPage` at `:67` (scan-page path) and `:81` (offset path). Observer failure FAILS the read (`:36-39`). `NewSourceWithObserver(capabilities, observer)` (`:32-34`); scan capability seam `scanPageReader.ListListingsScanPage` (`:50-52`).

## 9. sdk-runtime client

- File: `packages/sdk-runtime/src/index.ts`. Client object literal: `return {` at `:2113` through `};` at `:2446` (factory function closes `:2447`).
- Listings methods: `listListings` (`:2130`), `listListingsByProduct` (`:2132`), `getListing` (`:2134`), `getCategoryAttributes` (`:2136`), `getListingsSummary` (`:2140`), `refreshListings` (`:2146`, POST `/listings/refresh`).
- Orders methods: `importMarketplaceOrders` (`:2281`, POST `/orders/import`), `listMarketplaceOrders` (`:2284`, @deprecated), `listOrders` (`:2288`), `getOrder` (`:2290`), `markOrderFaturado` (`:2296`), `getOrderSummary` (`:2300`), `getAssistedSankhyaLinkage` (`:2304`), `listAssistedSankhyaLinkageCandidates` (`:2308`), `confirmAssistedSankhyaLinkage` (`:2312`).
- Pricing methods: `listPricingSimulations` (`:2352`), `runPricingSimulation` (`:2359`), `runBatchSimulation` (`:2361`), `getPricingProfile` (`:2363`), `putPricingProfile` (`:2364`), `listPricingDifal` (`:2366`), `putPricingDifalOverride` (`:2367`), `listPricingScenarios` (`:2369`), `createPricingScenario` (`:2370`), `deletePricingScenario` (`:2372`), `pricingDecompose` (`:2374`), `pricingSolveTarget` (`:2376`). No tariff-defaults method (see item 6).
- Sync: `listSyncRuns` (`:2144`, GET `/sync/runs`). Types `SyncRun`/`SyncRunPage`/`SyncRunListOptions`/`SyncRunStatus` at `:2461-2490`.
- Integrations: `listIntegrationProviders` (`:2176`), `listIntegrationInstallations` (`:2177`), `listIntegrationOperationRuns` (`:2178`), `startIntegrationAuthorization` (`:2180`), `startIntegrationReauthorization` (`:2182`), `submitIntegrationCredentials` (`:2184`), `disconnectIntegrationInstallation` (`:2186`), `getIntegrationAuthStatus` (`:2188`), `startIntegrationFeeSync` (`:2190`), `probeIntegrationAccount` (`:2192`), `probeIntegrationListings` (`:2194`), `probeIntegrationOrders` (`:2198`), `probeIntegrationFeeQuote` (`:2202`), `probeIntegrationStock` (`:2209`), `createIntegrationInstallation` (`:2357`).
- Config (active source / sortimento): `getActiveSource` (`:2163`), `setActiveSource` (`:2164`), `getSellableAssortment` (`:2166`), `setSellableAssortment` (`:2168`).

## 10. sync_state scheduler + repo

- File: `apps/server_core/internal/modules/sync/application/scheduler.go`
- `JobFunc` verbatim (`:46`, doc `:36-45`):
  ```go
  type JobFunc func(ctx context.Context, cursor json.RawMessage) (json.RawMessage, error)
  ```
  Doc contract: cursor is nil on first run; returned cursor is AUTHORITATIVE and replaces stored cursor verbatim; returning nil ERASES the persisted cursor (SQL NULL) — a no-progress cycle MUST re-return the received cursor.
- `func (s *Scheduler) RegisterJob(entity domain.Entity, fn JobFunc) error` (`:85`) — rejects invalid entity (`domain.ErrUnknownEntity`), nil fn, and duplicate entity registration.
- `func NewScheduler(store StateStore, installationID string, interval time.Duration, now func() time.Time) *Scheduler` (`:68`); `Start` (`:105`), `RunOnce` (`:124`).
- Nil-cursor / honest-read branch, `runJob` (`:141-161`) verbatim core:
  ```go
  state, _, err := s.store.Read(ctx, s.installationID, j.entity)
  if err != nil {
      return
  }
  cursor := state.Cursor
  ```
  (comment `:142-145`: a read ERROR skips the entity this cycle — it is NOT a first run; `found=false` with no error IS a first run → nil cursor.)
- Incremental flag call site (`:160`) verbatim:
  ```go
  _ = s.store.RecordSuccess(ctx, s.installationID, j.entity, next, s.now().UTC(), false)
  ```
  (hardcoded `incremental=false` — every scheduler success today records a full sync.)
- Failure path (`:154-157`): `RecordFailure(..., domain.SyncError{Message: jobErr.Error(), At: s.now().UTC()})`; panics converted per-job by `safeInvoke` (`:168-176`).
- StateStore port (`:22-34`) and repo `apps/server_core/internal/modules/sync/adapters/postgres/sync_state_repo.go`:
  - `func (r *SyncStateRepository) Read(ctx context.Context, installationID string, entity domain.Entity) (domain.SyncState, bool, error)` (`:35`)
  - `func (r *SyncStateRepository) RecordSuccess(ctx context.Context, installationID string, entity domain.Entity, cursor json.RawMessage, at time.Time, incremental bool) error` (`:62`)
  - `func (r *SyncStateRepository) RecordFailure(ctx context.Context, installationID string, entity domain.Entity, syncErr domain.SyncError) error` (`:86`)

## 11. Migrations

- Dir: `apps/server_core/migrations/`. Highest migration number: **0085**.
- Files with prefix >= 0080:
  - `0080_unaccent_extension.sql`
  - `0081_drop_legacy_marketplaces_and_orphan_tables.sql`
  - `0082_product_link_decisions.sql`
  - `0083_sellable_assortment_config.sql`
  - `0084_products_mirror_sellable_fields.sql`
  - `0085_erp_import_products_sellable_fields.sql`
- (Numbering has gaps: no 0002, 0040-0042, 0048-0049, 0054, 0059-0064, 0077; 0021 appears twice: `0021_integration_operation_run_evidence.sql` and `0021_integrations_provider_auth_strategy_shopee_partner.sql`.)

## 12. GET /sync/runs

- Handler: `apps/server_core/internal/modules/integrations/transport/run_read_handler.go` — route `mux.HandleFunc("GET /sync/runs", h.HandleList)` (`:34`), registered `httpx.InteractiveRouteClass` (`:32`); `HandleList` `:58-106`.
- Response shape verbatim (`:37-56`):
  ```go
  type runPageEnvelope struct {
      Items      []runReadItem `json:"items"`
      NextCursor *string       `json:"next_cursor"`
      PageSize   int           `json:"page_size"`
      AsOf       time.Time     `json:"as_of"`
  }
  type runReadItem struct {
      OperationRunID      string     `json:"operation_run_id"`
      InstallationID      string     `json:"installation_id"`
      Module              string     `json:"module"`
      Status              string     `json:"status"`
      ResultCode          string     `json:"result_code"`
      FailureCode         string     `json:"failure_code"`
      TranslatedErrorCode string     `json:"translated_error_code"`
      AttemptCount        int        `json:"attempt_count"`
      DurationMs          int64      `json:"duration_ms"`
      StartedAt           *time.Time `json:"started_at"`
      FinishedAt          *time.Time `json:"finished_at"`
  }
  ```
- Query params (per `ParseRunQuery` + tests `run_read_handler_test.go`, integration test `tests/integration/aggregate_sync_read_test.go:163-197`): `installation_id` (required), `cursor`, `limit`, `filter.module`, `filter.status`. Error codes: `invalid_cursor`, `invalid_filter`, installation-required, `internal_error` (`:65-82, 108-122`).
- SDK statuses: `SyncRunStatus = "queued" | "running" | "succeeded" | "failed" | "cancelled"` (`packages/sdk-runtime/src/index.ts:2461`).
