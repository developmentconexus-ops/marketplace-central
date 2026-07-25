package pricingtax

import (
	"context"
	"testing"

	pricingdomain "marketplace-central/apps/server_core/internal/modules/pricing/domain"
	pricingports "marketplace-central/apps/server_core/internal/modules/pricing/ports"
)

type fakeSource struct {
	profile   pricingdomain.CalcProfile
	profErr   error
	rate      pricingdomain.DifalRate
	rateErr   error
	askedUF   string
	askedTnt  string
	rateCalls int
}

func (f *fakeSource) GetProfile(_ context.Context, tenantID string) (pricingdomain.CalcProfile, error) {
	f.askedTnt = tenantID
	return f.profile, f.profErr
}

func (f *fakeSource) RateForUF(_ context.Context, tenantID, uf string) (pricingdomain.DifalRate, error) {
	f.rateCalls++
	f.askedTnt = tenantID
	f.askedUF = uf
	return f.rate, f.rateErr
}

// TestImpostoUsesProfileAliquotaAndRoundsToCentavos pins the default SIMPLES
// profile applied to a real order total: 4% of 129,90 is 5,196, which must
// reach the operator as 5,20 — the same centavo the Simulador shows.
func TestImpostoUsesProfileAliquotaAndRoundsToCentavos(t *testing.T) {
	source := &fakeSource{profile: pricingdomain.NewDefaultCalcProfile()}
	taxes, err := NewReader(source, "tenant-x").TaxesForOrder(context.Background(), 129.90, "SP")
	if err != nil {
		t.Fatalf("TaxesForOrder: %v", err)
	}
	if taxes.Imposto == nil || *taxes.Imposto != 5.20 {
		t.Fatalf("Imposto = %v, want 5.20", taxes.Imposto)
	}
	if source.askedTnt != "tenant-x" {
		t.Fatalf("tenant = %q, want tenant-x", source.askedTnt)
	}
}

// TestDifalOffIsExplicitZeroNotUnknown: a tenant with DIFAL switched off owes
// nothing, and that zero must be reported so the margem can close.
func TestDifalOffIsExplicitZeroNotUnknown(t *testing.T) {
	source := &fakeSource{profile: pricingdomain.NewDefaultCalcProfile()}
	taxes, err := NewReader(source, "tenant-x").TaxesForOrder(context.Background(), 100, "SP")
	if err != nil {
		t.Fatalf("TaxesForOrder: %v", err)
	}
	if taxes.Difal == nil || *taxes.Difal != 0 {
		t.Fatalf("Difal = %v, want explicit 0", taxes.Difal)
	}
	if source.rateCalls != 0 {
		t.Fatalf("rate lookups = %d, want 0 when DIFAL is off", source.rateCalls)
	}
}

// TestDifalOnUsesDestinoEfetivo checks the enabled path takes the shipment's
// destination state through the pricing DIFAL table: efetivo 18−12 = 6% of 200.
func TestDifalOnUsesDestinoEfetivo(t *testing.T) {
	profile := pricingdomain.NewDefaultCalcProfile()
	profile.DifalEnabled = true
	source := &fakeSource{
		profile: profile,
		rate: pricingdomain.DifalRate{
			UF:               "SP",
			InternaPct:       "18",
			InterestadualPct: "12",
			EfetivoPct:       "6",
		},
	}
	taxes, err := NewReader(source, "tenant-x").TaxesForOrder(context.Background(), 200, "SP")
	if err != nil {
		t.Fatalf("TaxesForOrder: %v", err)
	}
	if taxes.Difal == nil || *taxes.Difal != 12 {
		t.Fatalf("Difal = %v, want 12", taxes.Difal)
	}
	if source.askedUF != "SP" {
		t.Fatalf("asked UF = %q, want SP", source.askedUF)
	}
}

// TestDifalOnWithoutDestinoStaysUnknown: DIFAL is owed but we cannot say how
// much, so it must stay nil — never a convenient 0 that inflates the margem.
func TestDifalOnWithoutDestinoStaysUnknown(t *testing.T) {
	profile := pricingdomain.NewDefaultCalcProfile()
	profile.DifalEnabled = true
	source := &fakeSource{profile: profile}
	taxes, err := NewReader(source, "tenant-x").TaxesForOrder(context.Background(), 200, "")
	if err != nil {
		t.Fatalf("TaxesForOrder: %v", err)
	}
	if taxes.Difal != nil {
		t.Fatalf("Difal = %v, want nil", *taxes.Difal)
	}
	if taxes.Imposto == nil {
		t.Fatal("Imposto = nil, want the aliquota to still resolve")
	}
}

// TestDifalUFWithoutRateRowStaysUnknown: a state absent from the rate table is
// an honest gap, not a zero.
func TestDifalUFWithoutRateRowStaysUnknown(t *testing.T) {
	profile := pricingdomain.NewDefaultCalcProfile()
	profile.DifalEnabled = true
	source := &fakeSource{profile: profile, rateErr: pricingports.ErrDifalUFNotFound}
	taxes, err := NewReader(source, "tenant-x").TaxesForOrder(context.Background(), 200, "ZZ")
	if err != nil {
		t.Fatalf("TaxesForOrder: %v", err)
	}
	if taxes.Difal != nil {
		t.Fatalf("Difal = %v, want nil", *taxes.Difal)
	}
}
