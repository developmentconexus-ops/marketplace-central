package application_test

import (
	"context"
	"fmt"
	"testing"

	"marketplace-central/apps/server_core/internal/modules/pricing/application"
	pricingdomain "marketplace-central/apps/server_core/internal/modules/pricing/domain"
	pricingports "marketplace-central/apps/server_core/internal/modules/pricing/ports"
)

// stubTariffResolver returns a fixed Comissao.Valor (nil = unknown) from the
// shared ports.TariffResolver seam, mirroring the composite ML-category
// resolver used in production.
type stubTariffResolver struct {
	valor *string // Comissao.Valor to return; nil = unknown
	err   error
}

func (s *stubTariffResolver) Resolve(_ context.Context, _ pricingports.TariffRequest) (pricingdomain.TariffResolution, error) {
	if s.err != nil {
		return pricingdomain.TariffResolution{}, s.err
	}
	return pricingdomain.TariffResolution{Comissao: pricingdomain.ComponentResolution{Valor: s.valor}}, nil
}

type stubBatchProductProvider struct{}

func (s *stubBatchProductProvider) GetProductsForBatch(_ context.Context, _ []string) ([]pricingports.BatchProduct, error) {
	return []pricingports.BatchProduct{
		{ProductID: "prod-1", CategoryID: "electronics", CostAmount: 100, PriceAmount: 200},
	}, nil
}

type stubBatchPolicyProvider struct {
	override        *float64
	marketplaceCode string
}

func (s *stubBatchPolicyProvider) GetPoliciesForBatch(_ context.Context, _ []string) ([]pricingports.BatchPolicy, error) {
	return []pricingports.BatchPolicy{
		{
			PolicyID:           "pol-1",
			MarketplaceCode:    s.marketplaceCode,
			CommissionPercent:  0.10,
			CommissionOverride: s.override,
			FixedFeeAmount:     0,
			DefaultShipping:    0,
			ShippingProvider:   "fixed",
			MinMarginPercent:   0.05,
		},
	}, nil
}

type stubFreightQuoter struct{}

func (s *stubFreightQuoter) IsConnected(_ context.Context) (bool, error) { return false, nil }
func (s *stubFreightQuoter) QuoteFreight(_ context.Context, _ pricingports.FreightRequest) (map[string]pricingports.FreightResult, error) {
	return nil, nil
}

func ptrF64(f float64) *float64 { return &f }
func ptrStr(s string) *string   { return &s }

// runBatch wires a BatchOrchestrator with the given override/marketplaceCode
// and, when resolver is non-nil, attaches it via WithTariffResolver. It
// returns the whole first result item so tests can read CommissionAmount,
// CommissionSource and Status.
func runBatch(t *testing.T, override *float64, marketplaceCode string, resolver *stubTariffResolver) application.BatchSimulationItem {
	t.Helper()
	orch := application.NewBatchOrchestrator(
		&stubBatchProductProvider{},
		&stubBatchPolicyProvider{override: override, marketplaceCode: marketplaceCode},
		&stubFreightQuoter{},
		"tenant_default",
	)
	if resolver != nil {
		orch.WithTariffResolver(resolver)
	}
	result, err := orch.RunBatch(context.Background(), application.BatchRunRequest{
		ProductIDs: []string{"prod-1"},
		PolicyIDs:  []string{"pol-1"},
		Modalidade: pricingdomain.ModalidadeClassico,
	})
	if err != nil {
		t.Fatalf("RunBatch: %v", err)
	}
	if len(result.Items) == 0 {
		t.Fatal("expected at least one item")
	}
	return result.Items[0]
}

func TestBatchOrchestrator_CommissionOverrideTakesPriority(t *testing.T) {
	// override=0.05 wins over resolver="99" and policy=0.10
	// selling=200, commission=200*0.05=10
	item := runBatch(t, ptrF64(0.05), "mercado_livre", &stubTariffResolver{valor: ptrStr("99")})
	if item.CommissionAmount != 10 {
		t.Errorf("CommissionOverride priority: expected amount 10, got %v", item.CommissionAmount)
	}
	if item.CommissionSource != "override" {
		t.Errorf("CommissionOverride priority: expected source override, got %v", item.CommissionSource)
	}
}

func TestBatchOrchestrator_ResolverUsedWhenNoOverride(t *testing.T) {
	// No override; resolver="20.00" wins over policy=0.10
	// selling=200, commission=200*0.20=40
	item := runBatch(t, nil, "mercado_livre", &stubTariffResolver{valor: ptrStr("20.00")})
	if item.CommissionAmount != 40 {
		t.Errorf("resolver: expected amount 40, got %v", item.CommissionAmount)
	}
	if item.CommissionSource != "resolver" {
		t.Errorf("resolver: expected source resolver, got %v", item.CommissionSource)
	}
}

func TestBatchOrchestrator_PolicyRateUsedWhenNoResolver(t *testing.T) {
	// resolver nil (not wired); policy=0.10
	// selling=200, commission=200*0.10=20
	item := runBatch(t, nil, "mercado_livre", nil)
	if item.CommissionAmount != 20 {
		t.Errorf("policy fallback: expected amount 20, got %v", item.CommissionAmount)
	}
	if item.CommissionSource != "policy" {
		t.Errorf("policy fallback: expected source policy, got %v", item.CommissionSource)
	}
}

func TestBatchOrchestrator_UnresolvedCategoryIsHonestUnknown(t *testing.T) {
	// resolver wired but Valor=nil (unresolved category) → honest unknown,
	// NEVER the policy 0.10 default (ADR-17: unknown != zero/default).
	item := runBatch(t, nil, "mercado_livre", &stubTariffResolver{valor: nil})
	if item.CommissionSource != "unknown" {
		t.Errorf("unresolved category: expected source unknown, got %v", item.CommissionSource)
	}
	if item.Status != "critical" {
		t.Errorf("unresolved category: expected status critical, got %v", item.Status)
	}
	if item.CommissionAmount != 0 {
		t.Errorf("unresolved category: expected amount 0 (flagged sentinel, not fabricated 20), got %v", item.CommissionAmount)
	}
}

// TestBatchOrchestrator_ParityWithDecompose proves that the batch path and
// the single-product decompose path derive the SAME commission amount from
// the SAME shared resolver value, because both key commission off
// ports.TariffResolver.Resolve(...).Comissao.Valor. The input is chosen to
// be rounding-SENSITIVE — 78.90 @ 12% = 9.468 raw, which both paths must
// round half-up to 9.47; an unrounded batch (9.468) would fail here.
//
// This uses the pure-domain form (pricingdomain.Decompose) rather than a
// full CalcService stub: ports.CalcRepository is a large multi-method
// interface (DifalReader + profile/difal/scenario CRUD) that would need to
// be stubbed in full just to exercise GetProfile. Feeding the shared
// resolver's value into pricingdomain.Decompose directly (the same engine
// CalcService.Decompose delegates to once resolveTariff has produced the
// pct) still proves the commission math is identical off the shared seam.
func TestBatchOrchestrator_ParityWithDecompose(t *testing.T) {
	// Both inputs are rounding-SENSITIVE and cover the two ways float money
	// math diverges from the exact-decimal engine:
	//   78.90 @ 12% = 9.468   → half-up → 9.47 (non-boundary)
	//    6.10 @ 15% = 0.914999… (0.915 exact) → half-up → 0.92 (…x5 boundary
	//    that float round2 got wrong: round2(6.10*0.15)=0.91)
	// Both batch and decompose must agree to the cent off the shared percent.
	cases := []struct {
		price    float64
		pct      string
		expected string
	}{
		{price: 78.90, pct: "12.00", expected: "9.47"},
		{price: 6.10, pct: "15", expected: "0.92"},
	}
	for _, tc := range cases {
		shared := &stubTariffResolver{valor: ptrStr(tc.pct)}
		orch := application.NewBatchOrchestrator(
			&stubBatchProductProvider{},
			&stubBatchPolicyProvider{marketplaceCode: "mercado_livre"},
			&stubFreightQuoter{},
			"tenant_default",
		).WithTariffResolver(shared)
		result, err := orch.RunBatch(context.Background(), application.BatchRunRequest{
			ProductIDs:     []string{"prod-1"},
			PolicyIDs:      []string{"pol-1"},
			Modalidade:     pricingdomain.ModalidadeClassico,
			PriceOverrides: map[string]float64{"prod-1::pol-1": tc.price},
		})
		if err != nil {
			t.Fatalf("RunBatch(%.2f @ %s): %v", tc.price, tc.pct, err)
		}
		if len(result.Items) == 0 {
			t.Fatalf("RunBatch(%.2f @ %s): expected at least one item", tc.price, tc.pct)
		}
		batchItem := result.Items[0]
		batchCommissionStr := fmt.Sprintf("%.2f", batchItem.CommissionAmount)
		if batchCommissionStr != tc.expected {
			t.Fatalf("batch(%.2f @ %s): expected commission %s, got %s", tc.price, tc.pct, tc.expected, batchCommissionStr)
		}

		res, err := shared.Resolve(context.Background(), pricingports.TariffRequest{Modalidade: pricingdomain.ModalidadeClassico})
		if err != nil {
			t.Fatalf("shared resolver(%.2f @ %s): %v", tc.price, tc.pct, err)
		}
		decomp := pricingdomain.Decompose(pricingdomain.DecomposeInput{
			Preco:       fmt.Sprintf("%.2f", tc.price),
			ComissaoPct: *res.Comissao.Valor,
			AliquotaPct: "0",
			Modalidade:  pricingdomain.ModalidadeClassico,
			Custo:       &pricingdomain.Money{Amount: "100", Currency: "BRL"},
		})
		if decomp.Comissao != tc.expected {
			t.Fatalf("decompose(%.2f @ %s): expected Comissao %s, got %s", tc.price, tc.pct, tc.expected, decomp.Comissao)
		}
		if decomp.Comissao != batchCommissionStr {
			t.Fatalf("parity(%.2f @ %s): decompose commission %s != batch commission %s", tc.price, tc.pct, decomp.Comissao, batchCommissionStr)
		}
	}
}
