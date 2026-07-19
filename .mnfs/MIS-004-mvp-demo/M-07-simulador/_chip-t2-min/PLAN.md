# CHIP-T2-MIN — degrau 3 (COTAÇÃO) — Implementation Plan

BaseSha (fork): `40f7f4deaaece5be198789cf27cad28241a63add` (main tip @ LADDER GREEN)
Branch: `chip/t2-min-degrau3`
Lane: Claude-only (Codex quota exhausted). Plan = cold Opus planner (verified against code). P6 = Opus frio + sonnet adversarial, AGREEMENT.

## Verified facts (cold planner + own re-verification)
- IC-04 freeze pins ONLY Decomposition/DifalForUFResult/DecomposeInput (calc_ports_contract_test.go). Tariff types NOT frozen → additive changes safe.
- calcInputDTO already has `product_id *int` → resolution SERVER-SIDE, no codprod, no request change.
- OpenAPI PricingTarifaComponent already has `fonte:string`(no enum)+`data:string nullable` (both required); SDK mirrors → **ZERO contract change**.
- resolveTariff currently passes ONLY `{Modalidade}` (calc_service.go ~405); manual-override check comes after → gap to fill + manual gate.
- toTarifaDTO hardcodes `Data:nil` → must map through.
- D-21 composition-owned adapter pattern (market_adapters.go / orders_adapters.go) bridges pricing↔connectors w/o import cycle. `accountRefForTenant` (market_adapters.go:228) picks tenant first Connected ML installation. Template: `newOrdersShipmentReaderAdapter(mercadoLivreCapabilities, installationSvc, cfg.DefaultTenantID)` root.go:509.
- Probe (integrations/application/catalog_match_probe.go) records provider_operations audit EVERY call → degrau 3 reuses RAW connectors readers (mercadoLivreCapabilities.ReadCatalogMatch; marketplaceCapabilities.FeeQuoteReader("mercado_livre").ReadFeeQuote), NEVER the probe wrapper.
- listing_type: Classico→gold_special, Premium|Full→gold_pro (probe consts).

## Architecture (all additive)
New pricing ports (neutral): ProductIdentityReader, CategoryResolver, CommissionQuoter.
New pricing adapters: adapters/catalog/identity_reader.go (ProductIdentityReader over catalog svc);
  adapters/tarifflive/resolver.go (degrau-3 assembler, silent-miss); adapters/tariffcomposite/resolver.go (ports.TariffResolver, degrau3→degrau4).
Domain: FonteCotacao="COTACAO"; ComponentResolution.Data *string (additive).
TariffRequest grows: ProductID *int, PriceBasis *string (pre-authorized by port comment).
Composition ML adapters: composition/pricing_adapters.go (CategoryResolver+CommissionQuoter over raw readers + accountRefForTenant).

Call chain: transport → calc_service.resolveTariff(+productID/priceBasis, gate ProductID on reqComissao=="") →
  composite.Resolve → degrau4(base config) then degrau3.TryCommission(live, silent) →
  ProductIdentityReader → CategoryResolver(ReadCatalogMatch EAN+Query D-83) → CommissionQuoter(ReadFeeQuote gold_special/gold_pro) → comissao{COTACAO,degrau3,data}.
Miss (ProductID nil / no id / no EAN+titulo / no category / no commission / ANY error) → silent degrau-4 PADRAO. Never error, never fabricated (ADR-17).
Frete stays degrau 4 in MIN (no frete-by-dims).

## TDD slices (failing test first → minimal green → commit)
- [x] S1 domain+DTO: FonteCotacao const; ComponentResolution.Data; toTarifaDTO maps Data.
- [x] S2 ports: 3 new port files + TariffRequest.ProductID/PriceBasis. IC-04 freeze stays green.
- [x] S3 ProductIdentityReader adapter (catalog svc; int→string; unknown→ok=false).
- [ ] S4 tarifflive assembler (silent-miss matrix; D-83 EAN+titulo together; happy=COTACAO/3/data).
- [ ] S5 tariffcomposite (chain; frete untouched; ProductID nil → degrau3 skipped).
- [ ] S6 calc_service threading (productID+priceBasis; manual-commission gate → ProductID nil).
- [ ] S7 composition wiring (pricing_adapters.go raw readers + accountRefForTenant + modalidade→listing_type; root.go wire). Full build/vet/test + IC-04 green.

## Gates (close)
GOCACHE=C:/Users/leandro.theodoro/Documents/marketplace-central/.gocache absolute, no GOFLAGS.
go build/vet/test pricing+composition green; IC-04 freeze green; tsc sdk 0 (no SDK change, confirm); contract tests via mocks (never live ML); ZERO ML writes.
Evidence → this dir. P6 dual-gate AGREEMENT. CLOSED to hub w/ tip SHA.

## Risks (planner)
xlsx products have no price → Solve (no PriceBasis) w/ 0 catalog+no buybox price → silent degrau4 (honest, MIN-acceptable). No cotação cache in MIN (T3). Predita not surfaced in MIN wire (deferred).
