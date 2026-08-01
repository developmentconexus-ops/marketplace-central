// Fixture for TestFixture_AliasedSiteIsDetectedAndNamed
// (F-04-read-guard-allowlist, MIS-007). Reproduces the adversarial-review
// alias-bypass counterexample verbatim: the 4 real composition wiring calls,
// PLUS one INJECTED unallowlisted site reached not by passing
// mercadoLivreCapabilities directly, but by first assigning it to a
// local alias (mlAlias) and passing the alias to the new constructor.
//
// A detector that matches call arguments by literal identifier name alone
// (pre-hardening scanWiringFile) never flags newSneakyRefundAdapter here:
// mlAlias is not itself the string "mercadoLivreCapabilities", so the
// raw hit count stays at 4 (matching len(mlAllowlist)) and
// compareToAllowlist returns nil -- guard passes green -- despite a
// brand-new interactive site reaching the ML capability. This is exactly
// the one-line "extract variable" refactor the review demonstrated.
//
// This file lives under testdata/, which `go build`, `go vet`, and
// `go test` all skip when compiling real packages -- it is never part of
// any binary. It is only ever read as text by scanWiringFile, so it does
// not need to compile or type-check; the undefined identifiers below
// (mercadoLivreCapabilities, installationSvc, tenantID, feeReader) are
// intentional and harmless, same as the sibling fixtures.
package composition

func wireFixture() {
	ordersShipmentReader := newOrdersShipmentReaderAdapter(mercadoLivreCapabilities, installationSvc, tenantID)
	ordersBuyerFiscalReader := newOrdersBuyerFiscalReaderAdapter(mercadoLivreCapabilities, installationSvc, tenantID)
	categoryResolver := newPricingCategoryResolverAdapter(mercadoLivreCapabilities, installationSvc, tenantID)
	commissionQuoter := newPricingCommissionQuoterAdapter(feeReader, installationSvc, tenantID)

	// INJECTED: alias-bypass shape from the adversarial review -- a plain
	// local rename of mercadoLivreCapabilities into mlAlias, then passed to
	// a brand-new, unallowlisted constructor.
	mlAlias := mercadoLivreCapabilities
	sneakyRefundReader := newSneakyRefundAdapter(mlAlias, installationSvc, tenantID)

	_ = ordersShipmentReader
	_ = ordersBuyerFiscalReader
	_ = categoryResolver
	_ = commissionQuoter
	_ = sneakyRefundReader
}
