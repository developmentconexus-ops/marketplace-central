package pricingtax

import (
	"context"
	"math"
	"os"
	"strings"
	"testing"

	ordersports "marketplace-central/apps/server_core/internal/modules/orders/ports"
	pricingports "marketplace-central/apps/server_core/internal/modules/pricing/ports"
)

// unitPriceMatrix/unitPriceProducts is the fully-resolvable single-item setup
// these tests vary ONLY the UnitPrice on, so any nil in the result is
// attributable to the price and nothing else. Same MG fixture as
// TestTaxesForItemsUnlinkedItemNilsWholeOrder (migrations/0094_icms_matrix.sql:98).
func unitPriceReader() *Reader {
	matrix := &fakeMatrix{
		cells:     map[string]pricingports.MatrixCell{"MG|MG|1": {Found: true, CodTrib: intp(0), Ambiguo: false}},
		aliquotas: map[string]pricingports.AliquotaInterna{"MG": {AliquotaPct: "18", FcpEmbutidoPct: "0"}},
	}
	products := &fakeProducts{
		facts: map[string]pricingports.ProductFiscalFacts{
			"201": {Found: true, GrupoICMS: intp(1), Origprod: intp(0)},
		},
	}
	return NewReader(matrix, products, "tenant-1")
}

func unitPriceTaxes(t *testing.T, price float64, qty int) ordersports.OrderTaxes {
	t.Helper()
	got, err := unitPriceReader().TaxesForItems(context.Background(),
		[]ordersports.TaxItem{{InternalProductID: intp(201), UnitPrice: &price, Quantity: qty}}, "MG")
	if err != nil {
		t.Fatalf("TaxesForItems error: %v", err)
	}
	return got
}

// TestTaxesForItemsRejectsUnreadableUnitPrice is the contract this slice adds.
//
// UnitPrice crosses into the money path as float64 (ports/tax_reader.go).
// lineTotal used to round it to 2dp SILENTLY, so a value that is not the
// faithful image of any 2-decimal amount became a number nobody could tell
// apart from a real one: measured, lineTotal(1.005, 1) == "1.00".
//
// NOT because "1.00" is the wrong rounding — float64(1.005) is exactly
// 1.00499999999999989341858963598497211933 (measured, zz_probe_1005_test.go),
// below the half-up boundary, so "1.00" is the faithful answer for that
// double. The defect is that the literal carried a third decimal the money
// path does not have, and the conversion dropped it without saying so. Which
// 2dp amount was meant is unknowable, and unknown takes the unknown channel.
//
// That was only ever safe by accident of schema — unit_price is
// numeric(14,2) (pinned by TestUnitPriceColumnIsTwoDecimalNumeric below), so
// the DB read path cannot produce a third decimal. But nothing in the code
// SAID so, and application/ingest_service.go populates UnitPrice straight
// from the provider payload in memory, before that column ever rounds it.
//
// A price we cannot read faithfully is not a price. It takes the SAME
// named-unknown channel an unlinked item already uses (resolveItem's nil
// cell), which the all-or-nothing aggregation turns into a whole-order nil —
// never a silently rounded number.
func TestTaxesForItemsRejectsUnreadableUnitPrice(t *testing.T) {
	got := unitPriceTaxes(t, 1.005, 1)

	if got.ICMSSaida != nil || got.Difal != nil || got.PisCofins != nil || got.RestituicaoST != nil {
		t.Fatalf("got %s, want all four nil — UnitPrice 1.005 is not a readable 2dp amount and must nil the WHOLE order, never round to 1.00 in silence", dumpTaxes(got))
	}
}

// TestTaxesForItemsRejectsNonFiniteUnitPrice covers the second half of the
// same silence. Measured on the pre-guard body (zz_probe_linetotal_test.go):
// lineTotal(NaN, 1), lineTotal(+Inf, 1) and lineTotal(-Inf, 1) ALL returned
// "0.00" — a fabricated zero standing in for a value that does not exist,
// which is the ADR-17 failure mode exactly. Not-a-number is unknown, and
// unknown takes the unknown channel.
func TestTaxesForItemsRejectsNonFiniteUnitPrice(t *testing.T) {
	for name, price := range map[string]float64{
		"NaN":  math.NaN(),
		"+Inf": math.Inf(1),
		"-Inf": math.Inf(-1),
	} {
		got := unitPriceTaxes(t, price, 1)
		if got.ICMSSaida != nil || got.Difal != nil || got.PisCofins != nil || got.RestituicaoST != nil {
			t.Errorf("UnitPrice %s: got %s, want all four nil — never a fabricated 0.00", name, dumpTaxes(got))
		}
	}
}

// TestTaxesForItemsAcceptsFaithfulTwoDecimalPrices is the control that keeps
// the guard from being a blanket ban, and it is the one that actually decides
// the shape of the predicate.
//
// The naive test — "is this float exactly representable as n/100" — REJECTS
// almost every real price: float64's nearest value to 0.10 is
// 0.1000000000000000055511151231257827, so 0.10 x 100 is not an integer.
// A guard written that way would nil practically every order.
//
// The correct question is the other direction: is this float the faithful
// IMAGE of some 2dp decimal? Round it to 2dp, take the decimal back to
// float64, and compare. 0.10 survives (it is the nearest double to "0.10");
// 1.005 does not (its nearest double rounds to "1.00", which is a DIFFERENT
// double). Every value below is a legitimate numeric(14,2) amount and must
// produce numbers.
func TestTaxesForItemsAcceptsFaithfulTwoDecimalPrices(t *testing.T) {
	for _, tc := range []struct {
		price float64
		qty   int
	}{
		{0.10, 1}, {0.07, 3}, {19.90, 2}, {1000.00, 1}, {1234.56, 7}, {0.01, 999},
	} {
		got := unitPriceTaxes(t, tc.price, tc.qty)
		if got.ICMSSaida == nil || got.Difal == nil || got.PisCofins == nil {
			t.Errorf("UnitPrice %v x %d: got %s, want numbers — this IS a valid numeric(14,2) amount and must not trip the guard", tc.price, tc.qty, dumpTaxes(got))
		}
	}
}

// TestTaxesForItemsUnitPriceQuantityUsesDecimalNotFloat pins the second half:
// once the price is known to be a faithful 2dp amount, the line total is
// derived from that DECIMAL, not from the float. 0.07 x 3 is exactly 0.21;
// carried through float64 the product is 0.21000000000000002. Both round to
// "0.21" here, so this asserts the RESULT rather than the mechanism — but it
// is the case that drifts first if the derivation ever goes back through the
// float.
func TestTaxesForItemsUnitPriceQuantityUsesDecimalNotFloat(t *testing.T) {
	got := unitPriceTaxes(t, 0.07, 3)
	// preço = 0.21, aliquota interna MG = 18% ⇒ ICMS = round2(0.21 x 0.18) = 0.04.
	if !floatEq(got.ICMSSaida, "0.04") {
		t.Fatalf("ICMSSaida = %s, want 0.04 (linha 0.07 x 3 = 0.21, x 18%%)", dumpTaxes(got))
	}
}

// TestUnitPriceColumnIsTwoDecimalNumeric pins the schema fact the read path's
// safety rests on. Before this slice that fact was load-bearing and
// unasserted: widen the column (or add a path that skips it) and the silent
// rounding above goes live with nothing failing. Now the premise fails
// LOUDLY instead of lying.
func TestUnitPriceColumnIsTwoDecimalNumeric(t *testing.T) {
	const migration = "../../../../../migrations/0027_orders_marketplace_orders.sql"
	b, err := os.ReadFile(migration)
	if err != nil {
		t.Fatalf("read %s: %v", migration, err)
	}
	if !strings.Contains(string(b), "unit_price numeric(14,2)") {
		t.Fatalf("%s no longer declares `unit_price numeric(14,2)` — the read path's 2dp premise moved; revisit lineTotal's guard before changing this test", migration)
	}
}
