package domain

import "testing"

func TestUnknownOrderProfitability(t *testing.T) {
	got := UnknownOrderProfitability()

	if got.RetornoLiquido != nil {
		t.Fatalf("RetornoLiquido = %v, want nil", *got.RetornoLiquido)
	}
	if got.MargemPct != nil {
		t.Fatalf("MargemPct = %v, want nil", *got.MargemPct)
	}

	d := got.Decomposition
	for name, ptr := range map[string]*float64{
		"Comissao":    d.Comissao,
		"TaxaFixa":    d.TaxaFixa,
		"Frete":       d.Frete,
		"Imposto":     d.Imposto,
		"Difal":       d.Difal,
		"TarifaFull":  d.TarifaFull,
		"Custo":       d.Custo,
		"MargemValor": d.MargemValor,
		"MargemPct":   d.MargemPct,
	} {
		if ptr != nil {
			t.Fatalf("Decomposition.%s = %v, want nil", name, *ptr)
		}
	}

	wantComponents := []string{"comissao", "taxa_fixa", "frete", "imposto", "difal", "tarifa_full", "custo"}
	if len(d.ComponentesDesconhecidos) != len(wantComponents) {
		t.Fatalf("ComponentesDesconhecidos = %v, want %v", d.ComponentesDesconhecidos, wantComponents)
	}
	for i, want := range wantComponents {
		if d.ComponentesDesconhecidos[i] != want {
			t.Fatalf("ComponentesDesconhecidos[%d] = %q, want %q (full: %v)", i, d.ComponentesDesconhecidos[i], want, d.ComponentesDesconhecidos)
		}
	}

	f := got.Difal
	if f.Amount != nil {
		t.Fatalf("Difal.Amount = %v, want nil", *f.Amount)
	}
	if f.UFRoute != nil {
		t.Fatalf("Difal.UFRoute = %v, want nil", *f.UFRoute)
	}
	if f.DueDate != nil {
		t.Fatalf("Difal.DueDate = %v, want nil", *f.DueDate)
	}
	if f.Paid != nil {
		t.Fatalf("Difal.Paid = %v, want nil", *f.Paid)
	}
}

func TestUnknownOrderProfitabilityComponentsAreIndependentSlices(t *testing.T) {
	a := UnknownOrderProfitability()
	b := UnknownOrderProfitability()
	a.Decomposition.ComponentesDesconhecidos[0] = "mutated"
	if b.Decomposition.ComponentesDesconhecidos[0] != "comissao" {
		t.Fatalf("mutating one call's slice affected another: %v", b.Decomposition.ComponentesDesconhecidos)
	}
}

// TestBuildProfitabilityGoldenOrder2000017276984774 pins the EXEMPLO order's
// real comissão (22.95) surfaced verbatim and difal staying honest-unknown
// ("—") because no pricing engine is wired yet (ADR-17).
func TestBuildProfitabilityGoldenOrder2000017276984774(t *testing.T) {
	total := 229.50
	comissao := 22.95
	got := BuildProfitability(ProfitabilityInputs{
		Total:    &total,
		Comissao: &comissao,
	})
	if got.Decomposition.Comissao == nil || *got.Decomposition.Comissao != 22.95 {
		t.Fatalf("Decomposition.Comissao = %v, want 22.95", got.Decomposition.Comissao)
	}
	if got.Difal.Amount != nil {
		t.Fatalf("Difal.Amount = %v, want nil (honest unknown)", *got.Difal.Amount)
	}
	found := false
	for _, name := range got.Decomposition.ComponentesDesconhecidos {
		if name == "difal" {
			found = true
		}
	}
	if !found {
		t.Fatalf("ComponentesDesconhecidos = %v, want to contain difal", got.Decomposition.ComponentesDesconhecidos)
	}
}

// TestBuildProfitabilityAllKnownDerivesMargem verifies margem/retorno only
// derive when every source input is known, and MargemPct is a fraction (not
// pre-multiplied by 100 — the FE formatPercent does that).
func TestBuildProfitabilityAllKnownDerivesMargem(t *testing.T) {
	total, comissao, custo, frete, imposto := 100.0, 10.0, 40.0, 5.0, 3.0
	// DIFAL switched off on the tenant profile is an explicit 0, not an unknown.
	difal := 0.0
	got := BuildProfitability(ProfitabilityInputs{
		Total:    &total,
		Comissao: &comissao,
		Custo:    &custo,
		Frete:    &frete,
		Imposto:  &imposto,
		Difal:    &difal,
	})
	if got.Decomposition.MargemValor == nil || *got.Decomposition.MargemValor != 42 {
		t.Fatalf("MargemValor = %v, want 42", got.Decomposition.MargemValor)
	}
	if got.Decomposition.MargemPct == nil || *got.Decomposition.MargemPct != 0.42 {
		t.Fatalf("MargemPct = %v, want 0.42", got.Decomposition.MargemPct)
	}
	if got.RetornoLiquido == nil || *got.RetornoLiquido != 42 {
		t.Fatalf("RetornoLiquido = %v, want 42", got.RetornoLiquido)
	}
	if got.MargemPct == nil || *got.MargemPct != 0.42 {
		t.Fatalf("MargemPct (top-level) = %v, want 0.42", got.MargemPct)
	}
	wantUnknown := []string{"taxa_fixa", "tarifa_full"}
	if len(got.Decomposition.ComponentesDesconhecidos) != len(wantUnknown) {
		t.Fatalf("ComponentesDesconhecidos = %v, want %v", got.Decomposition.ComponentesDesconhecidos, wantUnknown)
	}
	for i, name := range wantUnknown {
		if got.Decomposition.ComponentesDesconhecidos[i] != name {
			t.Fatalf("ComponentesDesconhecidos[%d] = %q, want %q (full: %v)", i, got.Decomposition.ComponentesDesconhecidos[i], name, got.Decomposition.ComponentesDesconhecidos)
		}
	}
}

// TestBuildProfitabilityAnyUnknownInputSuppressesMargem asserts a single
// unknown input (Imposto here) suppresses margem/retorno entirely — never a
// partial/fabricated margin (ADR-17).
func TestBuildProfitabilityAnyUnknownInputSuppressesMargem(t *testing.T) {
	total, comissao, custo, frete := 100.0, 10.0, 40.0, 5.0
	got := BuildProfitability(ProfitabilityInputs{
		Total:    &total,
		Comissao: &comissao,
		Custo:    &custo,
		Frete:    &frete,
		Imposto:  nil,
	})
	if got.Decomposition.MargemValor != nil {
		t.Fatalf("MargemValor = %v, want nil", *got.Decomposition.MargemValor)
	}
	if got.Decomposition.MargemPct != nil {
		t.Fatalf("MargemPct = %v, want nil", *got.Decomposition.MargemPct)
	}
	if got.RetornoLiquido != nil {
		t.Fatalf("RetornoLiquido = %v, want nil", *got.RetornoLiquido)
	}
	if got.MargemPct != nil {
		t.Fatalf("MargemPct (top-level) = %v, want nil", *got.MargemPct)
	}
	found := false
	for _, name := range got.Decomposition.ComponentesDesconhecidos {
		if name == "imposto" {
			found = true
		}
	}
	if !found {
		t.Fatalf("ComponentesDesconhecidos = %v, want to contain imposto", got.Decomposition.ComponentesDesconhecidos)
	}
}

// TestBuildProfitabilityUnknownDifalSuppressesMargem pins that a DIFAL we could
// not source (destino state unknown while the tenant HAS DIFAL switched on) is
// not quietly treated as zero: the whole margem stays "—" (ADR-17).
func TestBuildProfitabilityUnknownDifalSuppressesMargem(t *testing.T) {
	total, comissao, custo, frete, imposto := 100.0, 10.0, 40.0, 5.0, 3.0
	got := BuildProfitability(ProfitabilityInputs{
		Total:    &total,
		Comissao: &comissao,
		Custo:    &custo,
		Frete:    &frete,
		Imposto:  &imposto,
	})
	if got.RetornoLiquido != nil {
		t.Fatalf("RetornoLiquido = %v, want nil when difal is unknown", *got.RetornoLiquido)
	}
	found := false
	for _, name := range got.Decomposition.ComponentesDesconhecidos {
		if name == "difal" {
			found = true
		}
	}
	if !found {
		t.Fatalf("ComponentesDesconhecidos = %v, want to contain difal", got.Decomposition.ComponentesDesconhecidos)
	}
}
