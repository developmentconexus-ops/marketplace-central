# M-02 Price Intelligence Core — Batch Feature Plan

Plan basis: accepted base `59d0e62fdbf15db068542432ef5d5731b6fa9f83`.

Global immutable boundaries:

- `apps/server_core/internal/modules/connectors/adapters/mercado_livre/listing_writer.go` and `price_writer.go` are forbidden to edit. Their `PUT /items/{id}` path through `doRawWithIdempotency` remains unchanged (seam map lines 77–101).
- All new provider operations use the existing `CapabilityAdapter.accessToken` + `doJSON` GET path and the existing `CredentialResolver.ResolveAccessToken` wiring (seam map lines 22–61, 142–178).
- `/sites/MLB/search`, scraping, provider DTO leakage, ML PUT/POST, silent pagination truncation, and integrity-read fallbacks are forbidden.
- Decimal values use the repository’s exact decimal-string convention; no decimal or metrics dependency is added (seam map lines 307–342).
- Every PostgreSQL repository is tenant-bound through `Repository.tenantID`, and every query contains `tenant_id = $1` (seam map lines 196–229, 325–332).
- PostgreSQL timestamp round-trip fixtures use `UTC().Truncate(time.Microsecond)` (seam map lines 344–348).
- The accepted absolute Go cache is `C:\Users\leandro.theodoro\Documents\marketplace-central\.claude\worktrees\chip-m02-price-intel-core\.gocache`.

## 1. Slice cards

### F-01 — ML adapter read ports

#### F-01-S1

- **id:** F-01-S1
- **goal:** Publish the seven IC-06 read interfaces, normalized domain shapes, typed errors, and capability-service accessors without exposing provider DTOs.
- **validation_kind:** unit-contract
- **commands:** `$env:GOCACHE='C:\Users\leandro.theodoro\Documents\marketplace-central\.claude\worktrees\chip-m02-price-intel-core\.gocache'; go test ./apps/server_core/internal/modules/connectors/domain ./apps/server_core/internal/modules/connectors/ports ./apps/server_core/internal/modules/connectors/application -count=1; go build ./...; go vet ./...`
- **expected_artifacts:** `market_read_test.go` transcripts covering nullable fields, exact enums, `errors.Is` typed errors, capability lookup, and compile-time interface assertions.
- **write_set:** `apps/server_core/internal/modules/connectors/domain/market_read.go`; `apps/server_core/internal/modules/connectors/domain/market_read_test.go`; `apps/server_core/internal/modules/connectors/ports/market_read.go`; `apps/server_core/internal/modules/connectors/application/marketplace_capability_service.go`; `apps/server_core/internal/modules/connectors/application/marketplace_capability_service_test.go`; `apps/server_core/internal/modules/connectors/adapters/mercado_livre/capability_adapter.go`
- **failing_test_first:** Add compile-time assertions for all seven ports plus table tests proving nil optional values stay nil and the five typed provider errors are distinguishable.
- **done_criteria:** Seven additive read capabilities are retrievable from `ProviderCapabilitySet`; account context remains `ProviderAccountRef`; normalized public types contain no `ml*` DTO; existing write capability fields and methods remain byte-functionally unchanged.
- **complexity:** standard→Luna-high
- **open_questions:** Concrete IC-06 types are incomplete: the contracts do not define the concrete `Attrs` shape, shipment `Costs` fields, free-shipping `item` request shape, or where `FetchedAt` sits in `SearchCatalogByEAN`’s list result.

#### F-01-S2

- **id:** F-01-S2
- **goal:** Implement `GetOwnItemPricing` and `GetPriceToWin` as additive authenticated GET operations with exact nullable mappings.
- **validation_kind:** unit-contract
- **commands:** `$env:GOCACHE='C:\Users\leandro.theodoro\Documents\marketplace-central\.claude\worktrees\chip-m02-price-intel-core\.gocache'; go test ./apps/server_core/internal/modules/connectors/adapters/mercado_livre -run 'Test(GetOwnItemPricing|GetPriceToWin|MarketReadErrorMapping)' -count=1; go build ./...; go vet ./...`
- **expected_artifacts:** Mock-server request transcripts for `/items/{id}/sale_price` and `/items/{id}/price_to_win?version=v2`, including 401/403/404/429/5xx/timeout mappings.
- **write_set:** `apps/server_core/internal/modules/connectors/adapters/mercado_livre/pricing_reader.go`; `apps/server_core/internal/modules/connectors/adapters/mercado_livre/pricing_reader_test.go`
- **failing_test_first:** Add request/response fixtures proving `TargetPrice`, `RegularAmount`, and `Position` remain nil when absent and `price_to_win` is not mapped to catalog minimum.
- **done_criteria:** Both routes use `accessToken` and `doJSON` with `http.MethodGet`; `FetchedAt` comes from the injected clock; no provider DTO crosses the adapter package.
- **complexity:** standard→Luna-high
- **open_questions:** []

#### F-01-S3

- **id:** F-01-S3
- **goal:** Implement EAN catalog search and catalog-product detail reads, preserving provider order and nullable buy-box data.
- **validation_kind:** unit-contract
- **commands:** `$env:GOCACHE='C:\Users\leandro.theodoro\Documents\marketplace-central\.claude\worktrees\chip-m02-price-intel-core\.gocache'; go test ./apps/server_core/internal/modules/connectors/adapters/mercado_livre -run 'Test(SearchCatalogByEAN|GetCatalogProduct)' -count=1; go build ./...; go vet ./...`
- **expected_artifacts:** Mock transcripts for `/products/search?product_identifier=...` and `/products/{id}`, including `buy_box_winner: null`, provider-order preservation, URL escaping, and `/sites/MLB/search` zero-call assertion.
- **write_set:** `apps/server_core/internal/modules/connectors/adapters/mercado_livre/catalog_reader.go`; `apps/server_core/internal/modules/connectors/adapters/mercado_livre/catalog_reader_test.go`
- **failing_test_first:** Add a null-buy-box fixture expecting `BuyBoxWinner == nil` and a multi-candidate fixture expecting exact provider order.
- **done_criteria:** Only approved official routes are called; null winner never becomes zero money; raw catalog response types remain private to `mercado_livre`.
- **complexity:** standard→Luna-high
- **open_questions:** []

#### F-01-S4

- **id:** F-01-S4
- **goal:** Implement the default-OFF catalog-offers flag, complete pagination, typed unavailability, and structured per-call telemetry.
- **validation_kind:** unit-contract
- **commands:** `$env:GOCACHE='C:\Users\leandro.theodoro\Documents\marketplace-central\.claude\worktrees\chip-m02-price-intel-core\.gocache'; go test ./apps/server_core/internal/modules/connectors/adapters/mercado_livre -run 'TestListCatalogOffers' -count=1; go build ./...; go vet ./...`
- **expected_artifacts:** Three-page success transcript; page-two failure transcript returning no partial slice; flag-OFF zero-HTTP transcript; captured `slog` records with `route`, `status`, `page`, `duration_ms`, and `count`.
- **write_set:** `apps/server_core/internal/modules/connectors/adapters/mercado_livre/catalog_offers_reader.go`; `apps/server_core/internal/modules/connectors/adapters/mercado_livre/catalog_offers_reader_test.go`; `apps/server_core/internal/modules/connectors/adapters/mercado_livre/capability_adapter.go`
- **failing_test_first:** Add a three-page fixture plus a second-page 503 fixture that rejects all accumulated results and matches `ErrCatalogOffersUnavailable`.
- **done_criteria:** Default construction and absent/false `MC_ML_CATALOG_OFFERS_ENABLED` perform no network call; ON exhausts pagination; any failed page returns typed unavailability and zero partial data; logs match the existing `slog` timing idiom (seam map lines 307–323).
- **complexity:** complex→Sol-low
- **open_questions:** []

#### F-01-S5

- **id:** F-01-S5
- **goal:** Implement shipment and free-shipping reads and add the real-installation live read-only contract suite for all seven ports.
- **validation_kind:** live-provider
- **commands:** `$env:GOCACHE='C:\Users\leandro.theodoro\Documents\marketplace-central\.claude\worktrees\chip-m02-price-intel-core\.gocache'; go test ./apps/server_core/internal/modules/connectors/adapters/mercado_livre -run 'Test(GetShipmentInfo|GetFreeShippingCost)' -count=1; npm run harness:live -- -Target provider -PreflightOnly; docker compose exec -T backend go test -tags=live ./apps/server_core/internal/modules/connectors/adapters/mercado_livre -run TestLiveCapabilityAdapterReadPorts -count=1 -v`
- **expected_artifacts:** Unit transcripts for shipment nullables and free-shipping unknown cost; `scripts/.runs/<run-id>/summary.txt`; live log listing only GET routes, recent `FetchedAt`, pagination telemetry, and redacted identifiers.
- **write_set:** `apps/server_core/internal/modules/connectors/adapters/mercado_livre/logistics_reader.go`; `apps/server_core/internal/modules/connectors/adapters/mercado_livre/logistics_reader_test.go`; `apps/server_core/internal/modules/connectors/adapters/mercado_livre/capability_adapter_live_test.go`
- **failing_test_first:** Add build-tagged live interface assertions and unit fixtures where SLA/cost/destination fields are absent and must remain nil.
- **done_criteria:** Shipment and shipping-option DTOs remain adapter-private; live suite obtains credentials only through the composed `CredentialResolver`; all observed provider methods are GET; sale-price and price-to-win are non-null for the selected active own listing.
- **complexity:** standard→Luna-high
- **open_questions:** []

### F-02 — Market persistence

#### F-02-S1

- **id:** F-02-S1
- **goal:** Define IC-01/IC-03 evidence entities and domain invariants without altering legacy observations or references.
- **validation_kind:** unit-contract
- **commands:** `$env:GOCACHE='C:\Users\leandro.theodoro\Documents\marketplace-central\.claude\worktrees\chip-m02-price-intel-core\.gocache'; go test ./apps/server_core/internal/modules/market/domain -run 'Test(MarketPriceSnapshot|ValidatedOffer|MatchDecision|CompetitiveSignal|MarketAggregate)' -count=1; go build ./...; go vet ./...`
- **expected_artifacts:** Domain table tests for exact uppercase enums, positive known prices, mandatory FAILED reason, nullable unknowns, aggregate counts, and three distinct verdict-state enums.
- **write_set:** `apps/server_core/internal/modules/market/domain/evidence.go`; `apps/server_core/internal/modules/market/domain/evidence_test.go`
- **failing_test_first:** Add constructors rejecting unknown-as-zero, FAILED without reason, non-BRL aggregate offers, and invalid enum combinations.
- **done_criteria:** New types coexist with existing `MarketObservation`/`MarketReference`; price-to-win, own price, catalog aggregate, and verdict states cannot be conflated by type.
- **complexity:** standard→Luna-high
- **open_questions:** []

#### F-02-S2

- **id:** F-02-S2
- **goal:** Add the complete tenant-scoped evidence schema in migrations 0050–0053 and set the exact migration fixture count to 45.
- **validation_kind:** migration-lane
- **commands:** `$env:GOCACHE='C:\Users\leandro.theodoro\Documents\marketplace-central\.claude\worktrees\chip-m02-price-intel-core\.gocache'; go test ./apps/server_core/internal/platform/migrate -count=1; npm run harness:pg:up; npm run harness:integration`
- **expected_artifacts:** Integration log showing four pending migrations on first pass and zero on second; schema assertions for five tables, tenant keys, uniqueness constraints, checks, and indexes; `runner_test.go` expected count `45`.
- **write_set:** `apps/server_core/migrations/0050_market_price_snapshots.sql`; `apps/server_core/migrations/0051_market_validated_offers.sql`; `apps/server_core/migrations/0052_market_match_decisions.sql`; `apps/server_core/migrations/0053_market_signals_aggregates.sql`; `apps/server_core/internal/platform/migrate/runner_test.go`; `apps/server_core/internal/modules/market/adapters/postgres/evidence_schema_integration_test.go`
- **failing_test_first:** Change the fixture/schema test to expect 45 migrations and the five absent tables before adding SQL.
- **done_criteria:** `0050` owns snapshots, `0051` offers, `0052` match decisions, and `0053` competitive signals plus aggregates; no ALTER touches `market_observations` or `market_references`; final count is `41 + 4 = 45`; `0054` is not created.
- **complexity:** complex→Sol-low
- **open_questions:** []

#### F-02-S3

- **id:** F-02-S3
- **goal:** Implement append/idempotency, latest-valid, and series repositories for snapshots and validated offers.
- **validation_kind:** hermetic-integration
- **commands:** `$env:GOCACHE='C:\Users\leandro.theodoro\Documents\marketplace-central\.claude\worktrees\chip-m02-price-intel-core\.gocache'; go test ./apps/server_core/internal/modules/market/adapters/postgres -run 'Test(Snapshot|ValidatedOffer)' -count=1; npm run harness:integration`
- **expected_artifacts:** SQL transcript for VALID → FAILED → latest VALID unchanged; duplicate-insert row count `1`; cross-tenant isolation; expired response without deletion.
- **write_set:** `apps/server_core/internal/modules/market/ports/evidence_store.go`; `apps/server_core/internal/modules/market/adapters/postgres/snapshot_repository.go`; `apps/server_core/internal/modules/market/adapters/postgres/snapshot_repository_integration_test.go`; `apps/server_core/internal/modules/market/adapters/postgres/offer_repository.go`; `apps/server_core/internal/modules/market/adapters/postgres/offer_repository_integration_test.go`
- **failing_test_first:** Seed VALID, append FAILED, assert the original `fetched_at` and price remain the latest valid while both rows remain visible in series.
- **done_criteria:** Idempotency key is `(tenant, scope, ref_id, source, fetched_at)`; latest-valid explicitly filters `status = VALID`; expiry is projected as EXPIRED without update/delete; every query binds tenant first.
- **complexity:** complex→Sol-low
- **open_questions:** []

#### F-02-S4

- **id:** F-02-S4
- **goal:** Implement repositories for match decisions, competitive signals, and aggregates with ordered batch reads.
- **validation_kind:** hermetic-integration
- **commands:** `$env:GOCACHE='C:\Users\leandro.theodoro\Documents\marketplace-central\.claude\worktrees\chip-m02-price-intel-core\.gocache'; go test ./apps/server_core/internal/modules/market/adapters/postgres -run 'Test(MatchDecision|CompetitiveSignal|MarketAggregate)' -count=1; npm run harness:integration`
- **expected_artifacts:** Integration transcripts proving anchors/contradictions round-trip, exact decimal strings, input-order reads, idempotency, and cross-tenant invisibility.
- **write_set:** `apps/server_core/internal/modules/market/adapters/postgres/match_decision_repository.go`; `apps/server_core/internal/modules/market/adapters/postgres/match_decision_repository_integration_test.go`; `apps/server_core/internal/modules/market/adapters/postgres/competitive_signal_repository.go`; `apps/server_core/internal/modules/market/adapters/postgres/competitive_signal_repository_integration_test.go`; `apps/server_core/internal/modules/market/adapters/postgres/aggregate_repository.go`; `apps/server_core/internal/modules/market/adapters/postgres/aggregate_repository_integration_test.go`
- **failing_test_first:** Add ordered-batch and tenant-A/tenant-B tests before repository queries.
- **done_criteria:** Stored decisions retain exact positive anchors and contradictions; aggregate rows retain source/fetched/computed times and counts; missing rows are not synthesized by persistence.
- **complexity:** complex→Sol-low
- **open_questions:** []

### F-03 — Identity resolver

#### F-03-S1

- **id:** F-03-S1
- **goal:** Implement the pure deterministic IC-01 resolver and the complete hard-negative fixture matrix.
- **validation_kind:** unit-contract
- **commands:** `$env:GOCACHE='C:\Users\leandro.theodoro\Documents\marketplace-central\.claude\worktrees\chip-m02-price-intel-core\.gocache'; go test ./apps/server_core/internal/modules/market/resolver -count=3; go build ./...; go vet ./...`
- **expected_artifacts:** At least 12 named fixtures, including Doka/Menegotti, Doka/VW, kit, quantity, color, measure/capacity, voltage, title-only, missing EAN, EAN collision, tie, and no candidate; three-run deterministic transcript.
- **write_set:** `apps/server_core/internal/modules/market/resolver/resolver.go`; `apps/server_core/internal/modules/market/resolver/normalization.go`; `apps/server_core/internal/modules/market/resolver/resolver_test.go`; `apps/server_core/internal/modules/market/resolver/fixtures_test.go`; `apps/server_core/internal/modules/market/resolver/testdata/identity_cases.json`
- **failing_test_first:** Add the frozen fixture table with exact expected status, confidence, band, anchors, contradictions, and selected candidate before resolver code.
- **done_criteria:** Two independent anchors are required for ACCEPT; hard negative always REJECTs; non-hard contradiction caps at REVIEW; missing/colliding EAN caps at REVIEW; equal ACCEPT-eligible candidates become REVIEW; zero candidates becomes NO_CANDIDATE.
- **complexity:** complex→Sol-low
- **open_questions:** IC-01 fixes confidence bands but does not supply anchor weights, exact numeric score calculation, normalization vocabulary, candidate tie-breaking order, or the exact hard-negative comparison table needed by the required exact-confidence fixtures.

#### F-03-S2

- **id:** F-03-S2
- **goal:** Persist resolver decisions through the F-02 decision store without adding endpoints or network behavior.
- **validation_kind:** hermetic-integration
- **commands:** `$env:GOCACHE='C:\Users\leandro.theodoro\Documents\marketplace-central\.claude\worktrees\chip-m02-price-intel-core\.gocache'; go test ./apps/server_core/internal/modules/market/application -run TestIdentityResolutionService -count=1; npm run harness:integration`
- **expected_artifacts:** Service test proving one persisted deterministic decision with inspectable anchors/contradictions and an integration row transcript.
- **write_set:** `apps/server_core/internal/modules/market/application/identity_resolution_service.go`; `apps/server_core/internal/modules/market/application/identity_resolution_service_test.go`
- **failing_test_first:** Add a fake decision-store test requiring the exact resolver result and a repository-backed round-trip for serialized evidence.
- **done_criteria:** Same input produces the same persisted decision; no HTTP, ML call, new migration, or product-links dependency is introduced.
- **complexity:** standard→Luna-high
- **open_questions:** []

### F-04 — Collection, verdict, reads, and API

#### F-04-S1

- **id:** F-04-S1
- **goal:** Implement validated-offer aggregation with seller deduplication and honest sample states.
- **validation_kind:** unit-contract
- **commands:** `$env:GOCACHE='C:\Users\leandro.theodoro\Documents\marketplace-central\.claude\worktrees\chip-m02-price-intel-core\.gocache'; go test ./apps/server_core/internal/modules/market/application -run TestMarketAggregator -count=1; go build ./...; go vet ./...`
- **expected_artifacts:** Aggregation table covering BRL/new/ACCEPT filtering, lowest-per-seller dedupe, median/min, exact `n_offers`/`n_sellers`, zero evidence, and 4-vs-5 seller boundary.
- **write_set:** `apps/server_core/internal/modules/market/application/aggregator.go`; `apps/server_core/internal/modules/market/application/aggregator_test.go`
- **failing_test_first:** Add fixtures where duplicate seller offers differ in price, invalid condition/currency/match rows are excluded, four sellers are insufficient, and five are OK.
- **done_criteria:** Provider order has no effect; no empty-set zero aggregate exists; `<5` valid sellers is exactly INSUFFICIENT_MARKET.
- **complexity:** complex→Sol-low
- **open_questions:** []

#### F-04-S2

- **id:** F-04-S2
- **goal:** Implement the verdict engine with three separate enums and complete unknown-cost/fee/freight/tax handling.
- **validation_kind:** unit-contract
- **commands:** `$env:GOCACHE='C:\Users\leandro.theodoro\Documents\marketplace-central\.claude\worktrees\chip-m02-price-intel-core\.gocache'; go test ./apps/server_core/internal/modules/market/application -run TestVerdictEngine -count=1; go build ./...; go vet ./...`
- **expected_artifacts:** Decision-table transcript for every match/evidence/blocking state, all four labels, and null verdict label with market range under SEM_CUSTO.
- **write_set:** `apps/server_core/internal/modules/market/domain/verdict.go`; `apps/server_core/internal/modules/market/domain/verdict_test.go`; `apps/server_core/internal/modules/market/application/verdict_engine.go`; `apps/server_core/internal/modules/market/application/verdict_engine_test.go`
- **failing_test_first:** Add the complete label/blocking decision table, including unavailable cost, fee, freight, and tax inputs, before implementation.
- **done_criteria:** `match_status`, `price_evidence_status`, and `blocking_state` remain separate; no unknown input becomes zero; a label is emitted only when every required profitability input is known.
- **complexity:** complex→Sol-low
- **open_questions:** No frozen contract defines the fee/freight/tax policies and source-time inputs or the numeric threshold table mapping market range/margin to `saudavel`, `viavel_preco_mercado`, `apertado`, and `nao_vale`. Implementing a table would re-decide frozen semantics or duplicate M-07 pricing policy.

#### F-04-S3

- **id:** F-04-S3
- **goal:** Build market-owned adapters for canonical product identity, linked listings, installation/account resolution, cost, fees, freight, and tax inputs.
- **validation_kind:** unit-contract
- **commands:** `$env:GOCACHE='C:\Users\leandro.theodoro\Documents\marketplace-central\.claude\worktrees\chip-m02-price-intel-core\.gocache'; go test ./apps/server_core/internal/modules/market/adapters/... -run 'Test(ProductIdentity|LinkedListings|CollectionAccount|VerdictInputs)' -count=1; go build ./...; go vet ./...`
- **expected_artifacts:** Adapter tests proving CODPROD-only identity, exact `ProviderAccountRef`, active Mercado Livre installation selection, nullable source facts, and no HTTP self-call.
- **write_set:** `apps/server_core/internal/modules/market/ports/collection_dependencies.go`; `apps/server_core/internal/modules/market/adapters/internalread/product_identity.go`; `apps/server_core/internal/modules/market/adapters/internalread/product_identity_test.go`; `apps/server_core/internal/modules/market/adapters/internalread/verdict_inputs.go`; `apps/server_core/internal/modules/market/adapters/internalread/verdict_inputs_test.go`; `apps/server_core/internal/modules/market/adapters/listings/linked_listing_reader.go`; `apps/server_core/internal/modules/market/adapters/listings/linked_listing_reader_test.go`; `apps/server_core/internal/modules/market/adapters/integrations/account_resolver.go`; `apps/server_core/internal/modules/market/adapters/integrations/account_resolver_test.go`
- **failing_test_first:** Add tests for nonexistent CODPROD, missing EAN, unresolved listing, no connected ML installation, multiple connected installations, and nil cost/tax facts.
- **done_criteria:** Existing published ports are adapted at the market boundary; no sibling module is edited; zero installations is unavailable and multiple installations is explicit ambiguity, never first-row fallback.
- **complexity:** standard→Luna-high
- **open_questions:** []

#### F-04-S4

- **id:** F-04-S4
- **goal:** Implement bounded synchronous collection, per-CODPROD concurrency exclusion, partial-failure accounting, and ADR-17-safe writes.
- **validation_kind:** hermetic-integration
- **commands:** `$env:GOCACHE='C:\Users\leandro.theodoro\Documents\marketplace-central\.claude\worktrees\chip-m02-price-intel-core\.gocache'; go test ./apps/server_core/internal/modules/market/application -run 'Test(CollectionService|CollectionGuard)' -count=1; npm run harness:integration`
- **expected_artifacts:** Transcripts for COMPLETED, PARTIAL, 404 product, 409 concurrent collection, missing EAN zero-catalog-call, flag-OFF no evidence, rate limit without retry storm, and VALID-before/after equality under provider failure.
- **write_set:** `apps/server_core/internal/modules/market/application/collection_guard.go`; `apps/server_core/internal/modules/market/application/collection_guard_test.go`; `apps/server_core/internal/modules/market/application/collection_orchestrator.go`; `apps/server_core/internal/modules/market/application/collection_orchestrator_test.go`; `apps/server_core/internal/modules/market/application/collection_orchestrator_integration_test.go`
- **failing_test_first:** Seed a VALID snapshot, force the second provider call to fail, and require PARTIAL plus an appended FAILED attempt while latest-valid remains byte-equivalent.
- **done_criteria:** POST work is synchronous and bounded; no job/poller/scheduler table exists; the same CODPROD conflicts while different CODPRODs may proceed; `ErrCatalogOffersUnavailable` becomes NO_PRICE_EVIDENCE; failures never overwrite valid evidence.
- **complexity:** complex→Sol-low
- **open_questions:** []

#### F-04-S5

- **id:** F-04-S5
- **goal:** Publish the frozen internal `market.EvidenceReader` and ordered batch read service for signals, aggregates, and verdicts.
- **validation_kind:** hermetic-integration
- **commands:** `$env:GOCACHE='C:\Users\leandro.theodoro\Documents\marketplace-central\.claude\worktrees\chip-m02-price-intel-core\.gocache'; go test ./apps/server_core/internal/modules/market/ports ./apps/server_core/internal/modules/market/application ./apps/server_core/internal/modules/market/adapters/postgres -run 'TestEvidenceReader' -count=1; npm run harness:integration`
- **expected_artifacts:** Compile-time consumer contract and repository-backed comparison showing identical domain objects for internal and HTTP serializers.
- **write_set:** `apps/server_core/internal/modules/market/ports/evidence_reader.go`; `apps/server_core/internal/modules/market/ports/evidence_reader_test.go`; `apps/server_core/internal/modules/market/application/evidence_read_service.go`; `apps/server_core/internal/modules/market/application/evidence_read_service_test.go`; `apps/server_core/internal/modules/market/adapters/postgres/verdict_repository.go`; `apps/server_core/internal/modules/market/adapters/postgres/verdict_repository_integration_test.go`
- **failing_test_first:** Add a compile-time fixture for `Signals(ctx, listingIDs)`, `Aggregates(ctx, codprods)`, and `Verdicts(ctx, codprods)` plus order-preservation tests.
- **done_criteria:** The three methods use slices and preserve input order; outputs carry source, fetched time, counts where aggregate, and match status; no implementation performs an HTTP self-call.
- **complexity:** complex→Sol-low
- **open_questions:** []

#### F-04-S6

- **id:** F-04-S6
- **goal:** Add the four new `/market/*` handlers and exact status/error behavior while leaving observations/references untouched.
- **validation_kind:** unit-contract
- **commands:** `$env:GOCACHE='C:\Users\leandro.theodoro\Documents\marketplace-central\.claude\worktrees\chip-m02-price-intel-core\.gocache'; go test ./apps/server_core/internal/modules/market/transport -run 'Test(MarketCollections|MarketSignals|MarketAggregates|MarketVerdicts)' -count=1; go build ./...; go vet ./...`
- **expected_artifacts:** HTTP transcripts for synchronous 200, 404 PRODUCT_NOT_FOUND, 409 COLLECTION_IN_PROGRESS, ordered GET arrays, empty arrays, exact `codprod` parameter, and legacy route regression.
- **write_set:** `apps/server_core/internal/modules/market/transport/intelligence_query.go`; `apps/server_core/internal/modules/market/transport/intelligence_query_test.go`; `apps/server_core/internal/modules/market/transport/intelligence_handler.go`; `apps/server_core/internal/modules/market/transport/intelligence_handler_test.go`; `apps/server_core/internal/modules/market/transport/http_handler.go`; `apps/server_core/internal/modules/market/transport/http_handler_test.go`
- **failing_test_first:** Register expected routes against a recording mux and add golden JSON/status tests before handler methods.
- **done_criteria:** POST returns 200, never 202; GET filters use `listing_ids` or exact `codprod`; observations/references registrations and response behavior remain unchanged.
- **complexity:** standard→Luna-high
- **open_questions:** []

#### F-04-S7

- **id:** F-04-S7
- **goal:** Atomically publish additive OpenAPI and `sdk-runtime/src/market.ts` contracts for all four routes.
- **validation_kind:** openapi-parity
- **commands:** `npm run build --workspace @marketplace-central/sdk-runtime; npm run test --workspace @marketplace-central/sdk-runtime -- --run; npm exec -- tsc --noEmit -p apps/web/tsconfig.json; npm run harness:governance -- -BaseSha 59d0e62fdbf15db068542432ef5d5731b6fa9f83`
- **expected_artifacts:** OpenAPI/SDK/handler parity test; SDK request transcripts; governance receipt; single atomic commit containing the OpenAPI section and `market.ts`.
- **write_set:** `contracts/api/marketplace-central.openapi.yaml`; `packages/sdk-runtime/src/market.ts`; `packages/sdk-runtime/src/market.test.ts`
- **failing_test_first:** Add SDK golden requests and a parity test reading the four OpenAPI path blocks before defining SDK types/client functions.
- **done_criteria:** Four additive paths and their schemas match handler JSON exactly; existing observation/reference SDK symbols remain available; OpenAPI and `market.ts` land in one commit under the hub-held contract lock.
- **complexity:** standard→Luna-high
- **open_questions:** IC-03 does not define the exact JSON members/types for collection `decisões`, `contagens`, `causas`, Verdict `inputs_used`, or the served market-price range. OpenAPI cannot be frozen without those shapes.

### Shared seams

#### SHARED-S1

- **id:** SHARED-S1
- **goal:** Add the one-line SDK barrel export under hub ownership in the same contract-lock transaction as F-04-S7.
- **validation_kind:** openapi-parity
- **commands:** `npm run build --workspace @marketplace-central/sdk-runtime; npm run test --workspace @marketplace-central/sdk-runtime -- --run; npm exec -- tsc --noEmit -p apps/web/tsconfig.json`
- **expected_artifacts:** `market.test.ts` proves consumers can import new market symbols from the package barrel; diff shows one additive export line in `index.ts`.
- **write_set:** `packages/sdk-runtime/src/index.ts`
- **failing_test_first:** F-04-S7 first imports one new symbol through `./index`, which fails until the hub adds the export.
- **done_criteria:** Exactly one additive barrel export is added; existing monolithic observation/reference client and types are not moved or renamed.
- **complexity:** standard→Luna-high
- **open_questions:** []

#### SHARED-S2

- **id:** SHARED-S2
- **goal:** Wire the real market collection/read dependencies and route registration in the composition root without a stub or alternate token path.
- **validation_kind:** hermetic-integration
- **commands:** `$env:GOCACHE='C:\Users\leandro.theodoro\Documents\marketplace-central\.claude\worktrees\chip-m02-price-intel-core\.gocache'; go test ./apps/server_core/internal/composition -run 'TestMarketIntelligenceRegistration|TestRootRuntime' -count=1; go build ./...; go vet ./...; npm run harness:integration`
- **expected_artifacts:** Root test proving all six legacy/new market routes register, real repositories are non-nil, the existing credential resolver feeds the same CapabilityAdapter, and catalog-offers defaults OFF.
- **write_set:** `apps/server_core/internal/composition/root.go`; `apps/server_core/internal/composition/root_test.go`
- **failing_test_first:** Extend the root route-registration test to require collections/signals/aggregates/verdicts and real dependencies before adding the composition block.
- **done_criteria:** Root adds only M-02 imports/configuration/registration; `MC_ML_CATALOG_OFFERS_ENABLED` is parsed default-false and passed to the existing adapter; no permanent fake/nil dependency ships.
- **complexity:** standard→Luna-high
- **open_questions:** []

### NEEDS INVESTIGATION

These questions must be resolved and folded into the frozen contracts or an explicit hub ruling before the affected cards are dispatched:

1. **F-01-S1 — IC-06 concrete Go shapes:** What are the exact field names/types for catalog `Attrs`, shipment `Costs`, free-shipping `item`, and search-result `FetchedAt`? The answer must preserve decimal strings and require no dependency.
2. **F-03-S1 — resolver decision table:** What exact normalization, anchor weights, score formula, hard-negative comparisons, and tie-break rules produce the required numeric confidence and confidence band?
3. **F-04-S2 — verdict policy:** Which fee, freight, tax, and cost policies/source times are mandatory, and what exact threshold table maps them plus market range to the four verdict labels without duplicating M-07?
4. **F-04-S7 — public shapes:** What exact JSON schema is frozen for collection decisions/counts/causes, Verdict `inputs_used`, and the market-price range returned under SEM_CUSTO?

No implementation slice may invent these answers.

## 2. Per-feature write-set and write-DAG

### F-01 write-DAG

Union:

- `apps/server_core/internal/modules/connectors/domain/market_read*`
- `apps/server_core/internal/modules/connectors/ports/market_read.go`
- `apps/server_core/internal/modules/connectors/application/marketplace_capability_service*`
- `apps/server_core/internal/modules/connectors/adapters/mercado_livre/capability_adapter.go`
- `pricing_reader*`, `catalog_reader*`, `catalog_offers_reader*`, `logistics_reader*`, and `capability_adapter_live_test.go` under the same adapter directory

Edges:

```text
F-01-S1 → F-01-S2 → F-01-S3 → F-01-S4 → F-01-S5
```

The sequence protects the shared `CapabilityAdapter` seam. `listing_writer.go` and `price_writer.go` are outside the union and forbidden.

### F-02 write-DAG

Union:

- `apps/server_core/internal/modules/market/domain/evidence*`
- `apps/server_core/internal/modules/market/ports/evidence_store.go`
- new evidence repository files/tests under `market/adapters/postgres`
- migrations `0050`–`0053`
- `apps/server_core/internal/platform/migrate/runner_test.go`

Edges:

```text
F-02-S1 → F-02-S2 → F-02-S3 → F-02-S4
```

`F-02-S3` and `F-02-S4` share `Repository` state and run serially. Migration `0054` is deliberately unused; final fixture count is exactly `45`.

### F-03 write-DAG

Union:

- `apps/server_core/internal/modules/market/resolver/**`
- `apps/server_core/internal/modules/market/application/identity_resolution_service*`

Edges:

```text
(F-01-S3 + F-02-S4 + investigation resolution) → F-03-S1 → F-03-S2
```

F-03 reads F-01 shapes and F-02 decision storage but edits neither seam.

### F-04 write-DAG

Union:

- `apps/server_core/internal/modules/market/domain/verdict*`
- `apps/server_core/internal/modules/market/application/{aggregator,verdict_engine,collection_guard,collection_orchestrator,evidence_read_service}*`
- `apps/server_core/internal/modules/market/ports/{collection_dependencies,evidence_reader}*`
- market-owned adapters under `market/adapters/{internalread,listings,integrations}`
- `market/adapters/postgres/verdict_repository*`
- market transport intelligence files plus additive edits to existing handler tests/registration
- OpenAPI `/market/*` section
- `packages/sdk-runtime/src/market.ts` and `market.test.ts`

Edges:

```text
F-02-S4 → F-04-S1
investigation resolution → F-04-S2
F-01-S5 + F-03-S2 → F-04-S3
F-04-S1 + F-04-S2 + F-04-S3 → F-04-S4
F-02-S4 + F-04-S2 → F-04-S5
F-04-S4 + F-04-S5 → F-04-S6
F-04-S6 + public-shape investigation → F-04-S7
F-04-S7 ↔ SHARED-S1  (same contract-lock transaction)
F-04-S6 + SHARED-S2 → live QA
```

### Shared-seam write-DAG

```text
F-04-S7 test-first → SHARED-S1 barrel export
F-01-S5 + F-02-S4 + F-03-S2 + F-04-S6 → SHARED-S2 root wiring
SHARED-S1 + SHARED-S2 → L0/L1/live-provider close
```

Only the hub writes `index.ts` and `root.go`. The F-02 owner writes the migration fixture under its pre-allocated grant.

## 3. Contract-satisfiability check

Current OpenAPI state at the accepted base:

| Claimed path | Current contract state | Collision result |
|---|---|---|
| `POST /market/collections` | Absent | New; satisfiable |
| `GET /market/signals` | Absent | New; satisfiable |
| `GET /market/aggregates` | Absent | New; satisfiable |
| `GET /market/verdicts` | Absent | New; satisfiable |
| `/market/observations` | Present at current OpenAPI line 2353 | Existing SAT endpoint; untouched |
| `/market/references` | Present at current OpenAPI line 2399 | Existing SAT endpoint; untouched |

The current `/market` path space contains only `/market/observations` and `/market/references`. All four claimed M-02 paths are genuinely additive.

Sibling-track comparison:

- M-01 claims `/erp/imports*`; no overlap with `/market/*`.
- M-03 claims no OpenAPI paths.
- No M-01 or M-03 claim occupies the four proposed paths.
- `/marketplaces/*` is a separate prefix and is not a collision.

SDK state:

- `packages/sdk-runtime/src/market.ts` does not exist.
- `packages/sdk-runtime/src/index.ts` already contains legacy `MarketObservation`, `MarketReference`, and their client methods.
- New `market.ts` is satisfiable if it preserves those symbols and avoids duplicate barrel-export names. SHARED-S1 is therefore an additive one-line export, not a refactor.
- The unresolved public response shapes listed under NEEDS INVESTIGATION are a planning defect; path availability alone does not make F-04-S7 dispatch-ready.

## 4. Pre-identified additive contract-locks

| Seam | File/section | Why | Owning slice / action |
|---|---|---|---|
| OpenAPI lock | `contracts/api/marketplace-central.openapi.yaml`, additive `/market/collections`, `/market/signals`, `/market/aggregates`, `/market/verdicts` paths and schemas | Shared contract artifact; must not collide with sibling work | Hub grants lock to F-04-S7 |
| SDK contract file | `packages/sdk-runtime/src/market.ts` | Manual client/types must land with OpenAPI under ADR-12 | F-04-S7, same commit as OpenAPI |
| SDK barrel | `packages/sdk-runtime/src/index.ts`, one export line | Hub-owned shared barrel; current file already hosts legacy market symbols | SHARED-S1, hub writer |
| Composition root | `apps/server_core/internal/composition/root.go`, existing connectors construction around seam-map lines 290–298 and market registration around lines 280–288 | Must wire the real credential, repository, cost, listing, and installation seams; no stub is acceptable | SHARED-S2, hub writer |
| Migration allocation | `apps/server_core/migrations/0050*`–`0053*` | M-02’s exclusive allocated migration block | F-02-S2 |
| Migration fixture | `apps/server_core/internal/platform/migrate/runner_test.go:25` | Current exact assertion is `len(want) != 41` (seam map lines 262–278); four migrations make the target `45` | F-02-S2; hub adjudicates merge conflicts |
| Migration 0054 | Reserved, no file | F-02’s four files already carry all five required tables; no F-04 column is currently justified | Must remain unused unless a reviewed schema defect proves necessity |
| Governance registry | `contracts/governance/modules.json` | No new module is created; `market` and `connectors` already exist | No touch; a change would be out of scope |

## 5. Per-criterion verification map

### Milestone C01–C07

| Criterion | Verification command / QA step | Test and implementation files |
|---|---|---|
| C01 synchronous collection summary | F-04-S4/S6 tests; live `Invoke-RestMethod -Method Post -Uri http://localhost:8080/market/collections -ContentType application/json -Body (@{codprod=$env:MC_QA_CODPROD}\|ConvertTo-Json -Compress)` | `collection_orchestrator*`, `collection_guard*`, `intelligence_handler*`, OpenAPI and `market.ts` |
| C02 IC-03 reads/evidence/enums | F-04-S1/S2/S5/S6 tests; live GET `/market/aggregates?codprod=...` and `/market/verdicts?codprod=...` | `aggregator*`, `verdict_engine*`, `evidence_read_service*`, `intelligence_handler*`, aggregate/verdict repositories |
| C03 Go EvidenceReader | `$env:GOCACHE='<absolute>'; go test ./apps/server_core/internal/modules/market/... -run TestEvidenceReader -count=1`; `rg -n 'http.*market/' apps/server_core/internal/modules/market` must show no self-call | `ports/evidence_reader*`, `application/evidence_read_service*`, `adapters/postgres/verdict_repository*` |
| C04 deterministic hard negatives | `$env:GOCACHE='<absolute>'; go test ./apps/server_core/internal/modules/market/resolver -count=3` | `resolver.go`, `normalization.go`, `fixtures_test.go`, `testdata/identity_cases.json` |
| C05 ADR-17 negative | F-02-S3 repository test plus F-04-S4 integration test with provider failure after seeded VALID | `snapshot_repository_integration_test.go`, `collection_orchestrator_integration_test.go` |
| C06 ML flag/pagination/telemetry/read-only | `npm run harness:live -- -Target provider -PreflightOnly`; live build-tag suite; `docker compose logs --no-color backend \| Select-String 'ml_catalog_offers'` | F-01 adapter files/tests and `capability_adapter_live_test.go`; SHARED-S2 root wiring |
| C07 migrations/seams | `Get-ChildItem apps/server_core/migrations/005[0-4]*`; migrate unit test; `npm run harness:integration`; governance against full base SHA | migrations `0050`–`0053`, `runner_test.go`, declared write sets, governance receipt |

### F-01 acceptance coverage

| Feature criterion | Verification | Files |
|---|---|---|
| Seven normalized read ports with FetchedAt | F-01-S1 compile/shape tests; F-01-S5 live suite | `domain/market_read*`, `ports/market_read.go`, reader files |
| Flag OFF returns typed unavailable without HTTP | F-01-S4 flag-OFF mock hit count `0` | `catalog_offers_reader_test.go` |
| 429 maps to `ErrRateLimited` without retry | F-01-S2/S4 error table and exact request count `1` | pricing/catalog-offers tests |
| Null buy-box remains nil | F-01-S3 null fixture | `catalog_reader_test.go` |
| 404/unauthorized/5xx/timeout typed mappings | F-01-S2 error table | `pricing_reader_test.go`, `domain/market_read_test.go` |
| Three-page completeness; failed middle page returns no partial result | F-01-S4 pagination tests | `catalog_offers_reader_test.go` |
| Structured route/status/page/duration/count telemetry | F-01-S4 captured slog test and C06 live logs | `catalog_offers_reader.go/test` |
| Provider DTO confinement | package compile/API review plus governance | private DTOs in adapter reader files; public types only in connectors domain |
| Existing token resolver reused | SHARED-S2 root test and live suite | `root.go/root_test.go`, `capability_adapter_live_test.go` |
| Read-only; no forbidden route | unit request method/path assertions and C06 live log | all F-01 adapter tests |
| Own price + price-to-win live | F-01-S5 live suite | `capability_adapter_live_test.go` |

### F-02 acceptance coverage

| Feature criterion | Verification | Files |
|---|---|---|
| Five tenant-scoped evidence tables | F-02-S2 schema integration test | migrations `0050`–`0053`, `evidence_schema_integration_test.go` |
| Idempotent snapshot writes | F-02-S3 duplicate insert test | `snapshot_repository*` |
| Latest VALID survives FAILED | F-02-S3 and F-04-S4 negative tests | snapshot repository and collection integration tests |
| FAILED reason mandatory | F-02-S1 constructor tests plus SQL check | `domain/evidence*`, migration `0050` |
| Unknown amount cannot become zero | F-02-S1 domain rejection | `domain/evidence_test.go` |
| Expiry is projected, not deleted | F-02-S3 read test and before/after row count | `snapshot_repository*` |
| Every query tenant-scoped | Cross-tenant tests in F-02-S3/S4 | all new PostgreSQL repository integration tests |
| Match anchors/contradictions persist | F-02-S4 round-trip | `match_decision_repository*` |
| Ordered signal/aggregate reads | F-02-S4 batch tests | signal and aggregate repositories |
| Migration fixture equals real count | migrate test expects 45 | `runner_test.go` |

### F-03 acceptance coverage

| Feature criterion | Verification | Files |
|---|---|---|
| EAN + brand ACCEPT with both anchors | Resolver fixture | `testdata/identity_cases.json`, `fixtures_test.go` |
| Hard negative with equal EAN REJECT | Doka and category-specific hard-negative fixtures | resolver fixture files |
| Non-hard contradiction caps REVIEW | Resolver fixture | resolver fixture files |
| Missing EAN caps REVIEW | Resolver fixture | resolver fixture files |
| Zero candidates is NO_CANDIDATE | Resolver fixture | resolver fixture files |
| Empty attributes are unavailable, not support/contradiction | Resolver fixture | resolver fixture files |
| Equal ACCEPT candidates become REVIEW | Resolver tie fixture | resolver fixture files |
| EAN collision caps REVIEW | Doka collision fixtures | resolver fixture files |
| Title-only never ACCEPTs | Resolver fixture | resolver fixture files |
| At least 12 exact cases | Test asserts fixture count and names | `fixtures_test.go` |
| Determinism across three runs | `go test ... -count=3` plus same-input loop | `resolver_test.go` |
| Persist exact evidence | F-03-S2 integration round-trip | `identity_resolution_service*`, F-02 match repository |

### F-04 acceptance coverage

| Feature criterion | Verification | Files |
|---|---|---|
| POST 200 synchronous summary; no job/polling | F-04-S4/S6 tests; OpenAPI parity; schema grep for no job table | orchestrator, handler, OpenAPI, migrations |
| COMPLETED/PARTIAL with named causes | Orchestrator decision table and HTTP golden JSON | `collection_orchestrator_test.go`, `intelligence_handler_test.go` |
| Missing EAN: NO_CANDIDATE/NO_PRICE_EVIDENCE and zero catalog calls | Fake-port call-count test | `collection_orchestrator_test.go` |
| Flag OFF: explicit NO_PRICE_EVIDENCE | F-01 typed error → F-04 orchestration test | `collection_orchestrator_test.go` |
| Provider failure appends FAILED and preserves VALID | F-04-S4 integration test | `collection_orchestrator_integration_test.go` |
| Rate limit does not retry-storm | Provider fake returns rate-limit; call count `1` | `collection_orchestrator_test.go` |
| 409 same-CODPROD collection conflict | Concurrent guard test and HTTP test | `collection_guard_test.go`, `intelligence_handler_test.go` |
| 404 nonexistent CODPROD | Identity adapter + HTTP test | `product_identity_test.go`, `intelligence_handler_test.go` |
| Aggregate BRL/new/ACCEPT/dedupe/<5 semantics | F-04-S1 table | `aggregator_test.go` |
| SEM_CUSTO serves range with null label | F-04-S2 decision table after investigation | `verdict_engine_test.go`, `domain/verdict_test.go` |
| Never-collected product returns NO_PRICE_EVIDENCE with null collection time | EvidenceReader/handler test | `evidence_read_service_test.go`, `intelligence_handler_test.go` |
| GET signals/aggregates/verdicts preserve input order | F-04-S5/S6 batch tests | EvidenceReader and handler tests |
| Every price read carries source/fetched/counts/match | Domain and parity tests | `domain/verdict*`, evidence read service, OpenAPI, `market.ts` |
| EvidenceReader matches HTTP shapes | Shared serializer/domain-object comparison | `evidence_read_service_test.go`, `intelligence_handler_test.go` |
| OpenAPI and SDK atomic parity | F-04-S7 commands and commit inspection | OpenAPI, `market.ts`, `market.test.ts` |
| Live collection, verdict, and price-to-win | Hub-provisioned dev stack POST/GET transcript | SHARED-S2 wiring, F-04 handler/orchestration, validation-result evidence |

### Full ladder before close

```powershell
$env:GOCACHE='C:\Users\leandro.theodoro\Documents\marketplace-central\.claude\worktrees\chip-m02-price-intel-core\.gocache'
go build ./...
go vet ./...
npm exec -- tsc --noEmit -p apps/web/tsconfig.json
npm run build --workspace @marketplace-central/sdk-runtime
npm run harness:governance -- -BaseSha 59d0e62fdbf15db068542432ef5d5731b6fa9f83

go test ./... -count=1
npm run test --workspace @marketplace-central/sdk-runtime -- --run
npm run harness:pg:up
npm run harness:integration

npm run harness:live -- -Target provider -PreflightOnly
```

The governance command runs from a clean detached verification worktree. The integration lane is serialized by the hub against sibling tracks with divergent migration fingerprints. C01/C02/C06 then run against the hub-provisioned dev stack; mocks cannot satisfy them.

## 6. Recommended slice dispatch order

1. Resolve all four NEEDS INVESTIGATION questions and freeze the resulting contract addenda.
2. Start F-01 and F-02 concurrently; their modules and files are disjoint.

   - F-01 serial chain: `S1 → S2 → S3 → S4 → S5`.
   - F-02 serial chain: `S1 → S2 → S3 → S4`.
   - F-01 and F-02 integration/live lanes may overlap only when the hub confirms their runtime resources do not collide.

3. After F-01-S3 shapes and F-02-S4 decision persistence are reviewed green, dispatch `F-03-S1 → F-03-S2`.
4. After F-01-S5, F-02-S4, and F-03-S2 are reviewed green:

   - F-04-S1 and F-04-S3 may overlap because their files are disjoint.
   - F-04-S2 may overlap with them only after its verdict-policy investigation is resolved.
   - Serialize `F-04-S4` after S1/S2/S3.
   - F-04-S5 may begin after F-02-S4 and the verdict domain from S2 are green.
   - Serialize `F-04-S6` after S4/S5.
   - Dispatch F-04-S7 only after public-shape investigation and handler JSON are frozen.

5. Under the hub’s contract lock, land `F-04-S7 + SHARED-S1` as one atomic OpenAPI/SDK commit.
6. Dispatch SHARED-S2 after all concrete constructors and handlers exist.
7. Run L0, L1, integration, dual gate, then hub-provisioned C01/C02/C06 live drive. Any fix restarts the ladder at L0.

**DISPATCH-READY: no**