//go:build integration

package integration

import (
	"context"
	"testing"

	marketplacespostgres "marketplace-central/apps/server_core/internal/modules/marketplaces/adapters/postgres"
	marketplacesapp "marketplace-central/apps/server_core/internal/modules/marketplaces/application"
	pricingpostgres "marketplace-central/apps/server_core/internal/modules/pricing/adapters/postgres"
	pricingapp "marketplace-central/apps/server_core/internal/modules/pricing/application"
	testpostgres "marketplace-central/apps/server_core/internal/testsupport/postgres"
)

// TestPhase1SmokeFlow exercises the full Phase 1 application service layer:
// create account → create policy → run simulation → list accounts → list policies → list simulations.
func TestPhase1SmokeFlow(t *testing.T) {
	pool, cfg := testpostgres.OpenPool(t, "tenant_harness_phase1")
	testpostgres.SeedMarketplaceDefinition(t, pool, testpostgres.MarketplaceDefinitionFixture{Code: "mercado_livre", DisplayName: "Mercado Livre", FeeSource: "api_sync"})
	testpostgres.CleanupMarketplaceAccounts(t, pool, cfg.DefaultTenantID, "smoke-acct-1")

	marketRepo := marketplacespostgres.NewRepository(pool, cfg.DefaultTenantID)
	marketSvc := marketplacesapp.NewService(marketRepo, cfg.DefaultTenantID)

	pricingRepo := pricingpostgres.NewRepository(pool, cfg.DefaultTenantID)
	pricingSvc := pricingapp.NewService(pricingRepo, cfg.DefaultTenantID)

	// Create account
	account, err := marketSvc.CreateAccount(context.Background(), marketplacesapp.CreateAccountInput{
		AccountID:       "smoke-acct-1",
		MarketplaceCode: "mercado_livre",
		ChannelCode:     "mercado_livre",
		DisplayName:     "Mercado Livre Smoke",
		ConnectionMode:  "api",
	})
	if err != nil {
		t.Fatalf("create account error: %v", err)
	}
	if account.AccountID == "" {
		t.Fatal("expected non-empty account ID")
	}

	// Create policy
	policy, err := marketSvc.CreatePolicy(context.Background(), marketplacesapp.CreatePolicyInput{
		PolicyID:           "smoke-policy-1",
		AccountID:          account.AccountID,
		CommissionPercent:  0.16,
		FixedFeeAmount:     5.0,
		DefaultShipping:    10.0,
		MinMarginPercent:   0.10,
		SLAQuestionMinutes: 60,
		SLADispatchHours:   24,
	})
	if err != nil {
		t.Fatalf("create policy error: %v", err)
	}
	if policy.PolicyID == "" {
		t.Fatal("expected non-empty policy ID")
	}

	// Run pricing simulation
	sim, err := pricingSvc.RunSimulation(context.Background(), pricingapp.RunSimulationInput{
		SimulationID:      "smoke-sim-1",
		ProductID:         "smoke-prod-1",
		AccountID:         account.AccountID,
		BasePriceAmount:   100.0,
		CostAmount:        60.0,
		CommissionPercent: policy.CommissionPercent,
		FixedFeeAmount:    policy.FixedFeeAmount,
		ShippingAmount:    policy.DefaultShipping,
		MinMarginPercent:  policy.MinMarginPercent,
	})
	if err != nil {
		t.Fatalf("run simulation error: %v", err)
	}
	if sim.Status == "" {
		t.Fatal("expected non-empty simulation status")
	}

	// List accounts
	accounts, err := marketSvc.ListAccounts(context.Background())
	if err != nil {
		t.Fatalf("list accounts error: %v", err)
	}
	if len(accounts) == 0 {
		t.Fatal("expected at least one account")
	}

	// List policies
	policies, err := marketSvc.ListPolicies(context.Background())
	if err != nil {
		t.Fatalf("list policies error: %v", err)
	}
	if len(policies) == 0 {
		t.Fatal("expected at least one policy")
	}

	// List simulations
	sims, err := pricingSvc.ListSimulations(context.Background())
	if err != nil {
		t.Fatalf("list simulations error: %v", err)
	}
	if len(sims) == 0 {
		t.Fatal("expected at least one simulation")
	}

	t.Logf("smoke flow OK: account=%s policy=%s sim=%s status=%s accounts=%d policies=%d sims=%d",
		account.AccountID, policy.PolicyID, sim.SimulationID, sim.Status,
		len(accounts), len(policies), len(sims))
}
