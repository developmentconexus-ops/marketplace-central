package pricingtax

import (
	"context"
	"fmt"
	"testing"

	ordersports "marketplace-central/apps/server_core/internal/modules/orders/ports"
	pricingports "marketplace-central/apps/server_core/internal/modules/pricing/ports"
)

func intp(i int) *int           { return &i }
func strp(s string) *string     { return &s }
func floatp(f float64) *float64 { return &f }

// fakeMatrix is a hand-rolled pricingports.TaxMatrixReader keyed by
// (ufOrigem, ufDestino, grupoICMS) for CellFor and by uf for
// AliquotaInternaFor — enough surface for this adapter's tests without
// pulling in the Postgres-backed MatrixReader (Task 4, a different layer).
type fakeMatrix struct {
	cells     map[string]pricingports.MatrixCell
	aliquotas map[string]pricingports.AliquotaInterna
}

func (f *fakeMatrix) CellFor(_ context.Context, _tenantID, ufOrigem, ufDestino string, grupoICMS int) (pricingports.MatrixCell, error) {
	key := fmt.Sprintf("%s|%s|%d", ufOrigem, ufDestino, grupoICMS)
	cell, ok := f.cells[key]
	if !ok {
		return pricingports.MatrixCell{Found: false}, nil
	}
	return cell, nil
}

func (f *fakeMatrix) AliquotaInternaFor(_ context.Context, uf string) (*pricingports.AliquotaInterna, error) {
	a, ok := f.aliquotas[uf]
	if !ok {
		return nil, nil
	}
	return &a, nil
}

// fakeProducts is a hand-rolled pricingports.ProductFiscalReader keyed by
// codigo_produto (already-resolved active source — that resolution is this
// task's ProductFiscalReader adapter's own job, a different layer from this
// one).
type fakeProducts struct {
	facts map[string]pricingports.ProductFiscalFacts
}

func (f *fakeProducts) FactsFor(_ context.Context, _tenantID, codigoProduto string) (pricingports.ProductFiscalFacts, error) {
	facts, ok := f.facts[codigoProduto]
	if !ok {
		return pricingports.ProductFiscalFacts{Found: false}, nil
	}
	return facts, nil
}

func floatEq(got *float64, want string) bool {
	if got == nil {
		return false
	}
	return fmt.Sprintf("%.2f", *got) == want
}

func dumpTaxes(t ordersports.OrderTaxes) string {
	deref := func(f *float64) string {
		if f == nil {
			return "<nil>"
		}
		return fmt.Sprintf("%.2f", *f)
	}
	return fmt.Sprintf("ICMSSaida=%s Difal=%s PisCofins=%s RestituicaoST=%s",
		deref(t.ICMSSaida), deref(t.Difal), deref(t.PisCofins), deref(t.RestituicaoST))
}

// TestTaxesForItemsSumsPerItemAcrossICMSGroups is the CRITICAL test (T5
// brief): two items in different ICMS grupos under the same order/destino —
// item A resolves as ST (codtrib=60, contributes an explicit "0.00"), item B
// resolves normally. The order's ICMSSaida must equal ONLY item B's
// contribution, proving the sum is built per-item, never a single rate
// applied to the order total. Fixtures mirror pricing/domain/icms_test.go
// cases C (ST, RJ) and F (origprod-importado, RJ) exactly, so the expected
// numbers are independently pinned there too.
func TestTaxesForItemsSumsPerItemAcrossICMSGroups(t *testing.T) {
	matrix := &fakeMatrix{
		cells: map[string]pricingports.MatrixCell{
			"MG|RJ|10": {Found: true, CodTrib: intp(60), Ambiguo: false}, // ST
			"MG|RJ|20": {Found: true, CodTrib: intp(0), Ambiguo: false},  // normal
		},
		// RJ real data (migrations/0094_icms_matrix.sql:104): aliquota=22.0,
		// fcp_embutido=2.0.
		aliquotas: map[string]pricingports.AliquotaInterna{"RJ": {AliquotaPct: "22", FcpEmbutidoPct: "2"}},
	}
	products := &fakeProducts{
		facts: map[string]pricingports.ProductFiscalFacts{
			"101": { // ST item — icms_test.go case C
				Found: true, GrupoICMS: intp(10), Origprod: nil,
				StRetidoEntrada: strp("40.00"), RestituicaoUnit: strp("15.00"),
			},
			"102": { // normal item — icms_test.go case F
				Found: true, GrupoICMS: intp(20), Origprod: intp(1),
			},
		},
	}
	r := NewReader(matrix, products, "tenant-1")

	items := []ordersports.TaxItem{
		{InternalProductID: intp(101), UnitPrice: floatp(800.00), Quantity: 1},
		{InternalProductID: intp(102), UnitPrice: floatp(1000.00), Quantity: 1},
	}

	got, err := r.TaxesForItems(context.Background(), items, "RJ")
	if err != nil {
		t.Fatalf("TaxesForItems error: %v", err)
	}

	// item A (ST) contributes an explicit "0.00"; item B contributes 40.00.
	// The order-level ICMSSaida must be exactly item B's value — proving the
	// zero from the ST item was summed in, not skipped, and the 40.00 was
	// never diluted or rated across both items.
	if !floatEq(got.ICMSSaida, "40.00") {
		t.Fatalf("ICMSSaida = %s, want 40.00 (only item B — item A's ST contribution is an explicit 0.00)", dumpTaxes(got))
	}
	// Task A2 recompute (gross-up replaced by diferença simples, D-41):
	// item 102's DIFAL is P×(aCusto−aInter) = 1000×(0.22−0.04) = 180.00
	// (was 230.77 under the old gross-up formula). Item 101 (ST) still
	// contributes 0.00 explicit. Sum = 0.00+180.00 = 180.00.
	if !floatEq(got.Difal, "180.00") {
		t.Fatalf("Difal = %s, want 180.00 (0.00 + 180.00)", dumpTaxes(got))
	}
	// item 102's PIS/COFINS base now uses aBase = aCusto−fcp = 0.22−0.02 =
	// 0.20: base = 1000×(1−0.20) = 800.00, PisCofins = 0.0925×800.00 = 74.00
	// (was 67.45 under the old gross-up `a`). Item 101 (ST branch,
	// untouched by this task) stays 70.30. Sum = 70.30+74.00 = 144.30.
	if !floatEq(got.PisCofins, "144.30") {
		t.Fatalf("PisCofins = %s, want 144.30 (70.30 + 74.00)", dumpTaxes(got))
	}
	if !floatEq(got.RestituicaoST, "15.00") {
		t.Fatalf("RestituicaoST = %s, want 15.00 (15.00 + 0.00)", dumpTaxes(got))
	}
}

// TestTaxesForItemsRejectsRateOnTotalBug is the mandated negative control:
// it hand-computes the exact bug this task exists to prevent — applying the
// FIRST item's effective rate to the ORDER TOTAL instead of computing each
// item's own ICMS and summing — and proves that naive number is observably
// different from (and NOT what the Reader returns for) the same fixtures as
// the critical test above. If a future change collapsed TaxesForItems back
// onto a single rate-on-total calculation, this test would catch it because
// the two numbers are provably different, not because of an incidental
// assertion.
func TestTaxesForItemsRejectsRateOnTotalBug(t *testing.T) {
	matrix := &fakeMatrix{
		cells: map[string]pricingports.MatrixCell{
			"MG|RJ|20": {Found: true, CodTrib: intp(0), Ambiguo: false},  // normal, first in this order
			"MG|RJ|10": {Found: true, CodTrib: intp(60), Ambiguo: false}, // ST, second
		},
		aliquotas: map[string]pricingports.AliquotaInterna{"RJ": {AliquotaPct: "22", FcpEmbutidoPct: "2"}},
	}
	products := &fakeProducts{
		facts: map[string]pricingports.ProductFiscalFacts{
			"102": {Found: true, GrupoICMS: intp(20), Origprod: intp(1)}, // a_inter=4% (origprod importado)
			"101": {
				Found: true, GrupoICMS: intp(10), Origprod: nil,
				StRetidoEntrada: strp("40.00"), RestituicaoUnit: strp("15.00"),
			},
		},
	}
	r := NewReader(matrix, products, "tenant-1")

	// item[0] is the NORMAL item (effective ICMS rate a_inter=4%), item[1] is
	// the ST item — the naive bug would take item[0]'s 4% and apply it to the
	// (800+1000=1800) order total.
	items := []ordersports.TaxItem{
		{InternalProductID: intp(102), UnitPrice: floatp(1000.00), Quantity: 1},
		{InternalProductID: intp(101), UnitPrice: floatp(800.00), Quantity: 1},
	}

	got, err := r.TaxesForItems(context.Background(), items, "RJ")
	if err != nil {
		t.Fatalf("TaxesForItems error: %v", err)
	}

	const correct = "40.00" // per-item: 40.00 (item B) + 0.00 (item A, ST)
	const naive = "72.00"   // BUG: 0.04 x (800.00+1000.00) = 72.00

	if naive == correct {
		t.Fatalf("negative control is degenerate: naive %s == correct %s, proves nothing", naive, correct)
	}
	if !floatEq(got.ICMSSaida, correct) {
		t.Fatalf("ICMSSaida = %s, want %s (correct per-item sum)", dumpTaxes(got), correct)
	}
	if floatEq(got.ICMSSaida, naive) {
		t.Fatalf("Reader collapsed onto the rate-on-total bug: ICMSSaida = %s (== naive %s)", dumpTaxes(got), naive)
	}
}

// TestTaxesForItemsUnlinkedItemNilsWholeOrder proves the all-or-nothing
// aggregation across items (never a partial sum over only the linked ones):
// one item is fully resolvable (MG interno, icms_test.go case A), the other
// has no InternalProductID at all. The whole order's four new components must
// come back nil — not the known item's contribution alone.
func TestTaxesForItemsUnlinkedItemNilsWholeOrder(t *testing.T) {
	matrix := &fakeMatrix{
		cells: map[string]pricingports.MatrixCell{
			"MG|MG|1": {Found: true, CodTrib: intp(0), Ambiguo: false},
		},
		// MG real data (migrations/0094_icms_matrix.sql:98): aliquota=18.0,
		// fcp_embutido=0. Task A2 (d): MG no longer bypasses this lookup —
		// without this entry item 201 would ALSO resolve to unknown
		// (D-37), defeating this test's premise that it is the ONE fully
		// resolvable item, contrasted against the unlinked second item.
		aliquotas: map[string]pricingports.AliquotaInterna{"MG": {AliquotaPct: "18", FcpEmbutidoPct: "0"}},
	}
	products := &fakeProducts{
		facts: map[string]pricingports.ProductFiscalFacts{
			"201": {
				Found: true, GrupoICMS: intp(1), Origprod: intp(0),
				RestituicaoUnit: strp("50.00"),
			},
		},
	}
	r := NewReader(matrix, products, "tenant-1")

	items := []ordersports.TaxItem{
		{InternalProductID: intp(201), UnitPrice: floatp(1000.00), Quantity: 1},
		{InternalProductID: nil, UnitPrice: floatp(200.00), Quantity: 1}, // never linked
	}

	got, err := r.TaxesForItems(context.Background(), items, "MG")
	if err != nil {
		t.Fatalf("TaxesForItems error: %v", err)
	}

	if got.ICMSSaida != nil || got.Difal != nil || got.PisCofins != nil || got.RestituicaoST != nil {
		t.Fatalf("got %s, want all four nil — one unlinked item must nil the WHOLE order, not just its own contribution", dumpTaxes(got))
	}
}

// TestTaxesForItemsScalesRestituicaoByQuantity proves S/R are per-unit
// products_mirror facts, scaled by Quantity before entering the ICMSCell:
// restituicao_unit="10.00" x Quantity=3 must yield RestituicaoST=30.00.
// GrupoICMS is left nil so ICMSSaida/Difal/PisCofins stay unknown (matrix
// cell never resolved) — isolating the assertion to the scaling behavior
// alone, since RestituicaoST is independent of the matrix cell resolution
// (pricing/domain/icms.go).
func TestTaxesForItemsScalesRestituicaoByQuantity(t *testing.T) {
	matrix := &fakeMatrix{}
	products := &fakeProducts{
		facts: map[string]pricingports.ProductFiscalFacts{
			"301": {Found: true, GrupoICMS: nil, RestituicaoUnit: strp("10.00")},
		},
	}
	r := NewReader(matrix, products, "tenant-1")

	items := []ordersports.TaxItem{
		{InternalProductID: intp(301), UnitPrice: floatp(100.00), Quantity: 3},
	}

	got, err := r.TaxesForItems(context.Background(), items, "SP")
	if err != nil {
		t.Fatalf("TaxesForItems error: %v", err)
	}

	if !floatEq(got.RestituicaoST, "30.00") {
		t.Fatalf("RestituicaoST = %s, want 30.00 (10.00 unit x Quantity=3)", dumpTaxes(got))
	}
}

func TestParseTaxComponentRejectsGarbageWithoutPanic(t *testing.T) {
	// Uma string de imposto que não parseia é erro de programação em
	// pricing/domain, mas em produção ela não pode derrubar o processo do
	// servidor inteiro: o pedido que não calcula é UM pedido, não o serviço.
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("panic em vez de erro: %v", r)
		}
	}()
	if _, err := parseTaxComponent("não é número", "ICMS saída"); err == nil {
		t.Fatal("err = nil, quero erro nomeando o componente")
	}
}
