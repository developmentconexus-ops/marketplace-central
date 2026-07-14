# F-01 catalog page port — validation evidence

The commands below were run from `apps/server_core`. Windows Go requires an
absolute `GOCACHE` path; it points to the workspace-local `.gocache` directory
for each transcript.

## Build

Command:

```text
$env:GOCACHE=(Resolve-Path .gocache).Path; go build ./...
```

Transcript:

```text
Exit code: 0
```

## Scoped catalog page tests

Command:

```text
$env:GOCACHE=(Resolve-Path .gocache).Path; go test ./internal/modules/internal_read/... ./internal/modules/catalog/... -run 'CatalogPage|Nullable|Quality|Ambiguous|Cursor' -v
```

Transcript:

```text
=== RUN   TestFakeReaderMissingStockStaysNilWithQualityFlag
--- PASS: TestFakeReaderMissingStockStaysNilWithQualityFlag (0.00s)
PASS
ok  marketplace-central/apps/server_core/internal/modules/internal_read/adapters/fake
=== RUN   TestCatalogPageUsesOneQueryForEveryListSize
=== RUN   TestCatalogPageUsesOneQueryForEveryListSize/limit-1
=== RUN   TestCatalogPageUsesOneQueryForEveryListSize/limit-50
=== RUN   TestCatalogPageUsesOneQueryForEveryListSize/limit-100
--- PASS: TestCatalogPageUsesOneQueryForEveryListSize (0.00s)
    --- PASS: TestCatalogPageUsesOneQueryForEveryListSize/limit-1 (0.00s)
    --- PASS: TestCatalogPageUsesOneQueryForEveryListSize/limit-50 (0.00s)
    --- PASS: TestCatalogPageUsesOneQueryForEveryListSize/limit-100 (0.00s)
=== RUN   TestCatalogPageCursorChainIsGaplessAndNonOverlapping
--- PASS: TestCatalogPageCursorChainIsGaplessAndNonOverlapping (0.00s)
=== RUN   TestCatalogPageMapsNullableFactsAndAmbiguousPrice
--- PASS: TestCatalogPageMapsNullableFactsAndAmbiguousPrice (0.00s)
=== RUN   TestCatalogPageInvalidCursorDoesNotTouchOracle
--- PASS: TestCatalogPageInvalidCursorDoesNotTouchOracle (0.00s)
=== RUN   TestCatalogSearchUsesOneBoundedQueryAndNoCursor
--- PASS: TestCatalogSearchUsesOneBoundedQueryAndNoCursor (0.00s)
=== RUN   TestSankhyaLinkageReaderPreservesOneToManyNullableDescendants
--- PASS: TestSankhyaLinkageReaderPreservesOneToManyNullableDescendants (0.00s)
=== RUN   TestSankhyaLinkageCandidateLineOverflowFailsAmbiguous
--- PASS: TestSankhyaLinkageCandidateLineOverflowFailsAmbiguous (0.00s)
PASS
ok  marketplace-central/apps/server_core/internal/modules/internal_read/adapters/oracle
testing: warning: no tests to run
PASS
ok  marketplace-central/apps/server_core/internal/modules/internal_read/application [no tests to run]
=== RUN   TestRequiredQualityFlagsRemainExplicit
--- PASS: TestRequiredQualityFlagsRemainExplicit (0.00s)
=== RUN   TestMissingCostStaysNilWithQualityFlag
--- PASS: TestMissingCostStaysNilWithQualityFlag (0.00s)
=== RUN   TestQualityFlagsAreStable
--- PASS: TestQualityFlagsAreStable (0.00s)
PASS
ok  marketplace-central/apps/server_core/internal/modules/internal_read/domain
testing: warning: no tests to run
PASS
ok  marketplace-central/apps/server_core/internal/modules/internal_read/observability [no tests to run]
=== RUN   TestCatalogCursorRoundTripsLastProductID
--- PASS: TestCatalogCursorRoundTripsLastProductID (0.00s)
=== RUN   TestCatalogCursorRejectsInvalidValues
--- PASS: TestCatalogCursorRejectsInvalidValues (0.00s)
PASS
ok  marketplace-central/apps/server_core/internal/modules/internal_read/ports
=== RUN   TestFactProjectsQualityFlagsWithoutInventingCurrentData
=== RUN   TestFactProjectsQualityFlagsWithoutInventingCurrentData/complete_known_zero
=== RUN   TestFactProjectsQualityFlagsWithoutInventingCurrentData/stale_keeps_prior_value
=== RUN   TestFactProjectsQualityFlagsWithoutInventingCurrentData/missing_discards_value
=== RUN   TestFactProjectsQualityFlagsWithoutInventingCurrentData/conflict_has_no_canonical_value
--- PASS: TestFactProjectsQualityFlagsWithoutInventingCurrentData (0.00s)
    --- PASS: TestFactProjectsQualityFlagsWithoutInventingCurrentData/complete_known_zero (0.00s)
    --- PASS: TestFactProjectsQualityFlagsWithoutInventingCurrentData/stale_keeps_prior_value (0.00s)
    --- PASS: TestFactProjectsQualityFlagsWithoutInventingCurrentData/missing_discards_value (0.00s)
    --- PASS: TestFactProjectsQualityFlagsWithoutInventingCurrentData/conflict_has_no_canonical_value (0.00s)
PASS
ok  marketplace-central/apps/server_core/internal/modules/catalog/adapters/internalread
testing: warning: no tests to run
PASS
ok  marketplace-central/apps/server_core/internal/modules/catalog/adapters/postgres [no tests to run]
testing: warning: no tests to run
PASS
ok  marketplace-central/apps/server_core/internal/modules/catalog/application [no tests to run]
testing: warning: no tests to run
PASS
ok  marketplace-central/apps/server_core/internal/modules/catalog/domain [no tests to run]
?   marketplace-central/apps/server_core/internal/modules/catalog/events [no test files]
?   marketplace-central/apps/server_core/internal/modules/catalog/ports [no test files]
?   marketplace-central/apps/server_core/internal/modules/catalog/readmodel [no test files]
?   marketplace-central/apps/server_core/internal/modules/catalog/transport [no test files]
```

## Full test suite

Command:

```text
$env:GOCACHE=(Resolve-Path .gocache).Path; go test ./...
```

Transcript:

```text
?    marketplace-central/apps/server_core/cmd/migrate [no test files]
?    marketplace-central/apps/server_core/cmd/server [no test files]
?    marketplace-central/apps/server_core/cmd/testdb [no test files]
ok   marketplace-central/apps/server_core/internal/composition
ok   marketplace-central/apps/server_core/internal/modules/catalog/adapters/internalread
ok   marketplace-central/apps/server_core/internal/modules/catalog/adapters/postgres
ok   marketplace-central/apps/server_core/internal/modules/catalog/application
ok   marketplace-central/apps/server_core/internal/modules/catalog/domain
?    marketplace-central/apps/server_core/internal/modules/catalog/events [no test files]
?    marketplace-central/apps/server_core/internal/modules/catalog/ports [no test files]
?    marketplace-central/apps/server_core/internal/modules/catalog/readmodel [no test files]
?    marketplace-central/apps/server_core/internal/modules/catalog/transport [no test files]
ok   marketplace-central/apps/server_core/internal/modules/classifications/adapters/postgres
?    marketplace-central/apps/server_core/internal/modules/classifications/application [no test files]
?    marketplace-central/apps/server_core/internal/modules/classifications/domain [no test files]
?    marketplace-central/apps/server_core/internal/modules/classifications/ports [no test files]
?    marketplace-central/apps/server_core/internal/modules/classifications/transport [no test files]
?    marketplace-central/apps/server_core/internal/modules/connectors/adapters/magalu [no test files]
ok   marketplace-central/apps/server_core/internal/modules/connectors/adapters/melhorenvio
ok   marketplace-central/apps/server_core/internal/modules/connectors/adapters/mercado_livre
?    marketplace-central/apps/server_core/internal/modules/connectors/adapters/shopee [no test files]
ok   marketplace-central/apps/server_core/internal/modules/connectors/application
?    marketplace-central/apps/server_core/internal/modules/connectors/domain [no test files]
?    marketplace-central/apps/server_core/internal/modules/connectors/events [no test files]
?    marketplace-central/apps/server_core/internal/modules/connectors/ports [no test files]
?    marketplace-central/apps/server_core/internal/modules/connectors/readmodel [no test files]
?    marketplace-central/apps/server_core/internal/modules/connectors/transport [no test files]
ok   marketplace-central/apps/server_core/internal/modules/integrations/adapters/amazon
ok   marketplace-central/apps/server_core/internal/modules/integrations/adapters/crypto
ok   marketplace-central/apps/server_core/internal/modules/integrations/adapters/feesync
ok   marketplace-central/apps/server_core/internal/modules/integrations/adapters/leroymerlin
ok   marketplace-central/apps/server_core/internal/modules/integrations/adapters/madeiramadeira
ok   marketplace-central/apps/server_core/internal/modules/integrations/adapters/magalu
ok   marketplace-central/apps/server_core/internal/modules/integrations/adapters/mercadolivre
ok   marketplace-central/apps/server_core/internal/modules/integrations/adapters/postgres
ok   marketplace-central/apps/server_core/internal/modules/integrations/adapters/providers
ok   marketplace-central/apps/server_core/internal/modules/integrations/adapters/shopee
ok   marketplace-central/apps/server_core/internal/modules/integrations/application
ok   marketplace-central/apps/server_core/internal/modules/integrations/background
ok   marketplace-central/apps/server_core/internal/modules/integrations/domain
?    marketplace-central/apps/server_core/internal/modules/integrations/ports [no test files]
ok   marketplace-central/apps/server_core/internal/modules/integrations/transport
ok   marketplace-central/apps/server_core/internal/modules/internal_read/adapters/fake
ok   marketplace-central/apps/server_core/internal/modules/internal_read/adapters/oracle
ok   marketplace-central/apps/server_core/internal/modules/internal_read/application
ok   marketplace-central/apps/server_core/internal/modules/internal_read/domain
ok   marketplace-central/apps/server_core/internal/modules/internal_read/observability
ok   marketplace-central/apps/server_core/internal/modules/internal_read/ports
?    marketplace-central/apps/server_core/internal/modules/inventory/adapters/connectors [no test files]
?    marketplace-central/apps/server_core/internal/modules/inventory/adapters/integrations [no test files]
ok   marketplace-central/apps/server_core/internal/modules/inventory/adapters/internalread
?    marketplace-central/apps/server_core/internal/modules/inventory/adapters/postgres [no test files]
?    marketplace-central/apps/server_core/internal/modules/inventory/adapters/productlinks [no test files]
ok   marketplace-central/apps/server_core/internal/modules/inventory/application
ok   marketplace-central/apps/server_core/internal/modules/inventory/domain
?    marketplace-central/apps/server_core/internal/modules/inventory/ports [no test files]
ok   marketplace-central/apps/server_core/internal/modules/inventory/transport
?    marketplace-central/apps/server_core/internal/modules/marketplaces/adapters/postgres [no test files]
ok   marketplace-central/apps/server_core/internal/modules/marketplaces/application
?    marketplace-central/apps/server_core/internal/modules/marketplaces/domain [no test files]
?    marketplace-central/apps/server_core/internal/modules/marketplaces/events [no test files]
?    marketplace-central/apps/server_core/internal/modules/marketplaces/ports [no test files]
?    marketplace-central/apps/server_core/internal/modules/marketplaces/readmodel [no test files]
ok   marketplace-central/apps/server_core/internal/modules/marketplaces/registry
ok   marketplace-central/apps/server_core/internal/modules/marketplaces/transport
ok   marketplace-central/apps/server_core/internal/modules/orders/adapters/integrations
ok   marketplace-central/apps/server_core/internal/modules/orders/adapters/internalread
?    marketplace-central/apps/server_core/internal/modules/orders/adapters/postgres [no test files]
ok   marketplace-central/apps/server_core/internal/modules/orders/adapters/productlinks
ok   marketplace-central/apps/server_core/internal/modules/orders/application
ok   marketplace-central/apps/server_core/internal/modules/orders/domain
?    marketplace-central/apps/server_core/internal/modules/orders/ports [no test files]
ok   marketplace-central/apps/server_core/internal/modules/orders/transport
ok   marketplace-central/apps/server_core/internal/modules/pricing/adapters/catalog
?    marketplace-central/apps/server_core/internal/modules/pricing/adapters/feeschedule [no test files]
ok   marketplace-central/apps/server_core/internal/modules/pricing/adapters/marketplace
ok   marketplace-central/apps/server_core/internal/modules/pricing/adapters/postgres
ok   marketplace-central/apps/server_core/internal/modules/pricing/application
?    marketplace-central/apps/server_core/internal/modules/pricing/domain [no test files]
?    marketplace-central/apps/server_core/internal/modules/pricing/events [no test files]
?    marketplace-central/apps/server_core/internal/modules/pricing/ports [no test files]
?    marketplace-central/apps/server_core/internal/modules/pricing/readmodel [no test files]
ok   marketplace-central/apps/server_core/internal/modules/pricing/transport
?    marketplace-central/apps/server_core/internal/modules/product_links/adapters/postgres [no test files]
ok   marketplace-central/apps/server_core/internal/modules/product_links/application
ok   marketplace-central/apps/server_core/internal/modules/product_links/domain
?    marketplace-central/apps/server_core/internal/modules/product_links/ports [no test files]
ok   marketplace-central/apps/server_core/internal/modules/product_links/transport
?    marketplace-central/apps/server_core/internal/modules/profitability/adapters/internalread [no test files]
ok   marketplace-central/apps/server_core/internal/modules/profitability/adapters/orders
?    marketplace-central/apps/server_core/internal/modules/profitability/adapters/postgres [no test files]
ok   marketplace-central/apps/server_core/internal/modules/profitability/application
?    marketplace-central/apps/server_core/internal/modules/profitability/domain [no test files]
?    marketplace-central/apps/server_core/internal/modules/profitability/ports [no test files]
ok   marketplace-central/apps/server_core/internal/modules/profitability/transport
?    marketplace-central/apps/server_core/internal/platform/config [no test files]
ok   marketplace-central/apps/server_core/internal/platform/httpx
?    marketplace-central/apps/server_core/internal/platform/logging [no test files]
ok   marketplace-central/apps/server_core/internal/platform/migrate
ok   marketplace-central/apps/server_core/internal/platform/pgdb
ok   marketplace-central/apps/server_core/internal/testsupport/postgres
ok   marketplace-central/apps/server_core/migrations
ok   marketplace-central/apps/server_core/tests/integration
ok   marketplace-central/apps/server_core/tests/unit
```
