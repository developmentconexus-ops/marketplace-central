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
