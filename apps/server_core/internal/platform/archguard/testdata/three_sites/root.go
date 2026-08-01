// Fixture for TestFixture_ShrunkAllowlistStillPasses
// (F-04-read-guard-allowlist, MIS-007). Mirrors only 1 of the CURRENT 2-entry
// mlAllowlist's real composition wiring calls in
// ../../../composition/root.go (post F-03-read-path-switch: site A retired
// entirely, site B reclassified batch-only -- see archguard_test.go's
// mlAllowlist doc comment), simulating a future milestone retiring
// pricing.solve.commission_quoter / newPricingCommissionQuoterAdapter too.
//
// This file lives under testdata/ and is never built/vetted as part of any
// real package -- see the note in ../five_sites/root.go.
package composition

func wireFixture() {
	categoryResolver := newPricingCategoryResolverAdapter(mercadoLivreCapabilities, installationSvc, tenantID)

	_ = categoryResolver
}
