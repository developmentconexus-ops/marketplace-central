// Fixture for TestFixture_FifthSiteIsDetectedAndNamed
// (F-04-read-guard-allowlist, MIS-007). Mirrors the 4 real composition
// wiring calls in ../../../composition/root.go, PLUS one INJECTED 5th
// site (newOrdersRefundReaderAdapter) simulating a future change that wires
// a new interactive path to the Mercado Livre capability without updating
// mlAllowlist.
//
// This file lives under testdata/, which `go build`, `go vet`, and
// `go test` all skip when compiling real packages -- it is never part of
// any binary. It is only ever read as text by scanWiringFile, so it does
// not need to compile or type-check; the undefined identifiers below
// (mercadoLivreCapabilities, installationSvc, tenantID, feeReader) are
// intentional and harmless.
package composition

func wireFixture() {
	ordersShipmentReader := newOrdersShipmentReaderAdapter(mercadoLivreCapabilities, installationSvc, tenantID)
	ordersBuyerFiscalReader := newOrdersBuyerFiscalReaderAdapter(mercadoLivreCapabilities, installationSvc, tenantID)
	categoryResolver := newPricingCategoryResolverAdapter(mercadoLivreCapabilities, installationSvc, tenantID)
	commissionQuoter := newPricingCommissionQuoterAdapter(feeReader, installationSvc, tenantID)

	// INJECTED: a new interactive refund-lookup adapter reaching ML without
	// any corresponding mlAllowlist entry.
	ordersRefundReader := newOrdersRefundReaderAdapter(mercadoLivreCapabilities, installationSvc, tenantID)

	_ = ordersShipmentReader
	_ = ordersBuyerFiscalReader
	_ = categoryResolver
	_ = commissionQuoter
	_ = ordersRefundReader
}
