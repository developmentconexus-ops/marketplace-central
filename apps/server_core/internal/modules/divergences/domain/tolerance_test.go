package domain

import "testing"

// IC-02 Persistence Expectations: estoque tolerance is 0 — any difference
// diverges.
func TestDivergesEstoqueAnyDifferenceDiverges(t *testing.T) {
	diverges, err := Diverges(KindEstoque, "12", "9")
	if err != nil {
		t.Fatalf("Diverges: %v", err)
	}
	if !diverges {
		t.Fatal("expected estoque 12 vs 9 to diverge (tolerance 0)")
	}
}

func TestDivergesEstoqueExactMatchDoesNotDiverge(t *testing.T) {
	diverges, err := Diverges(KindEstoque, "12", "12")
	if err != nil {
		t.Fatalf("Diverges: %v", err)
	}
	if diverges {
		t.Fatal("expected estoque 12 vs 12 (exact match) to NOT diverge")
	}
}

// IC-02 Must Preserve: tariff tolerance is R$0.01 absolute per line; ≤0.01 is
// rounding, not divergence; >0.01 diverges.
func TestDivergesTarifaExactlyPointZeroOneDoesNotDiverge(t *testing.T) {
	diverges, err := Diverges(KindTarifa, "24.33", "24.34")
	if err != nil {
		t.Fatalf("Diverges: %v", err)
	}
	if diverges {
		t.Fatal("expected tarifa delta of exactly 0.01 to NOT diverge (rounding, not disagreement)")
	}
}

func TestDivergesTarifaAboveTolerance(t *testing.T) {
	diverges, err := Diverges(KindTarifa, "24.33", "24.35")
	if err != nil {
		t.Fatalf("Diverges: %v", err)
	}
	if !diverges {
		t.Fatal("expected tarifa delta of 0.02 to diverge")
	}
}

func TestDivergesTarifaDirectionDoesNotMatter(t *testing.T) {
	diverges, err := Diverges(KindTarifa, "24.35", "24.33")
	if err != nil {
		t.Fatalf("Diverges: %v", err)
	}
	if !diverges {
		t.Fatal("expected tarifa delta of 0.02 in the other direction to also diverge (absolute delta)")
	}
}

func TestDivergesRejectsUnknownKind(t *testing.T) {
	if _, err := Diverges("bogus", "1", "1"); err == nil {
		t.Fatal("expected an error for an unknown kind")
	}
}

func TestDivergesRejectsUnparsableValue(t *testing.T) {
	if _, err := Diverges(KindEstoque, "not-a-number", "1"); err == nil {
		t.Fatal("expected an error for an unparsable expected value")
	}
}
