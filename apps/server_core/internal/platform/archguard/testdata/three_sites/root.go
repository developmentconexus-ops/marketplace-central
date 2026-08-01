// Fixture for TestFixture_ShrunkAllowlistStillPasses
// (F-04-read-guard-allowlist, MIS-007). Mirrors only 3 of the 4 real
// composition wiring calls in ../../../composition/root.go, simulating a
// future milestone (M-03/M-07 per AGENTS.md) retiring site D
// (pricing.solve.commission_quoter / newPricingCommissionQuoterAdapter).
//
// This file lives under testdata/ and is never built/vetted as part of any
// real package -- see the note in ../five_sites/root.go.
package composition

func wireFixture() {
	ordersShipmentReader := newOrdersShipmentReaderAdapter(mercadoLivreCapabilities, installationSvc, tenantID)
	ordersBuyerFiscalReader := newOrdersBuyerFiscalReaderAdapter(mercadoLivreCapabilities, installationSvc, tenantID)
	categoryResolver := newPricingCategoryResolverAdapter(mercadoLivreCapabilities, installationSvc, tenantID)

	_ = ordersShipmentReader
	_ = ordersBuyerFiscalReader
	_ = categoryResolver
}
