package domain

import "testing"

func TestNewDefaultCalcProfileIsSimplesFourPercentWithExplicitOrigem(t *testing.T) {
	profile := NewDefaultCalcProfile()

	if profile.Regime != RegimeSimples {
		t.Fatalf("Regime = %q, want %q", profile.Regime, RegimeSimples)
	}
	if profile.AliquotaPct != "4" {
		t.Fatalf("AliquotaPct = %q, want 4", profile.AliquotaPct)
	}
	if profile.Origem != OrigemDefault {
		t.Fatalf("Origem = %q, want %q (explicit, not implied)", profile.Origem, OrigemDefault)
	}
	if profile.LimiarVerdePct != "18" {
		t.Fatalf("LimiarVerdePct = %q, want 18", profile.LimiarVerdePct)
	}
	if profile.LimiarAmareloPct != "10" {
		t.Fatalf("LimiarAmareloPct = %q, want 10", profile.LimiarAmareloPct)
	}
	if profile.TarifaFull != nil {
		t.Fatalf("TarifaFull = %+v, want nil (nullable, unset by default)", profile.TarifaFull)
	}
}

func TestDefaultAliquotaPctPerRegime(t *testing.T) {
	if got := DefaultAliquotaPct(RegimeSimples); got != "4" {
		t.Fatalf("DefaultAliquotaPct(SIMPLES) = %q, want 4", got)
	}
	if got := DefaultAliquotaPct(RegimePresumido); got != "9.25" {
		t.Fatalf("DefaultAliquotaPct(PRESUMIDO) = %q, want 9.25", got)
	}
}

func TestCalcProfileTarifaFullIsNullableMoney(t *testing.T) {
	known := &Money{Amount: "12.50", Currency: "BRL"}
	profile := CalcProfile{
		Regime:      RegimePresumido,
		AliquotaPct: DefaultAliquotaPct(RegimePresumido),
		TarifaFull:  known,
	}
	if profile.TarifaFull.IsUnknown() {
		t.Fatalf("a set TarifaFull must not report unknown")
	}

	profile.TarifaFull = nil
	if !profile.TarifaFull.IsUnknown() {
		t.Fatalf("nil TarifaFull must report unknown (ADR-17), never a zero value")
	}
}
