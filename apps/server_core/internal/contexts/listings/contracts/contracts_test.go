package contracts_test

import (
	"strings"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/contexts/listings/contracts"
	"marketplace-central/apps/server_core/internal/kernel/channel"
	"marketplace-central/apps/server_core/internal/kernel/exact"
	"marketplace-central/apps/server_core/internal/kernel/fact"
	"marketplace-central/apps/server_core/internal/kernel/provenance"
	"marketplace-central/apps/server_core/internal/kernel/tenant"
)

func testEvidence(t *testing.T) provenance.Evidence {
	t.Helper()
	e, err := provenance.NewEvidence("mercado_livre", "item", "MLB123", time.Date(2026, 8, 10, 12, 0, 0, 0, time.UTC), "hash-1")
	if err != nil {
		t.Fatalf("evidence: %v", err)
	}
	return e
}

func testKey(t *testing.T) contracts.SourceListingKey {
	t.Helper()
	tid, err := tenant.Parse("tenant_default")
	if err != nil {
		t.Fatalf("tenant: %v", err)
	}
	code, err := channel.ParseCode("mercado_livre")
	if err != nil {
		t.Fatalf("channel: %v", err)
	}
	account, err := channel.NewAccountRef(code, "179571326")
	if err != nil {
		t.Fatalf("account: %v", err)
	}
	key, err := contracts.NewSourceListingKey(tid, account, "MLB123")
	if err != nil {
		t.Fatalf("key: %v", err)
	}
	return key
}

func TestNewSourceListingKeyRejectsBlankListingID(t *testing.T) {
	tid, _ := tenant.Parse("tenant_default")
	code, _ := channel.ParseCode("mercado_livre")
	account, _ := channel.NewAccountRef(code, "179571326")
	if _, err := contracts.NewSourceListingKey(tid, account, "  "); err == nil {
		t.Fatal("blank listing id accepted")
	}
}

func TestValidateRejectsZeroKeyAndZeroEvidence(t *testing.T) {
	title, _ := fact.NewUnknown[string]("ml omitted title", testEvidence(t))
	obs := contracts.ListingObservation{State: contracts.ListingState{Title: title}}
	if err := obs.Validate(); err == nil || !strings.Contains(err.Error(), "key") {
		t.Fatalf("zero key accepted or wrong error: %v", err)
	}
	obs.Key = testKey(t)
	if err := obs.Validate(); err == nil || !strings.Contains(err.Error(), "evidence") {
		t.Fatalf("zero evidence accepted or wrong error: %v", err)
	}
}

func TestValidateRejectsVariationWithBlankID(t *testing.T) {
	obs := contracts.ListingObservation{
		Key:      testKey(t),
		Evidence: testEvidence(t),
		State:    contracts.ListingState{Variations: []contracts.VariationObservation{{VariationID: ""}}},
	}
	if err := obs.Validate(); err == nil || !strings.Contains(err.Error(), "variation") {
		t.Fatalf("blank variation id accepted or wrong error: %v", err)
	}
}

func TestValidateAcceptsAllFactsUnknown(t *testing.T) {
	// Um anúncio de que o ML só devolveu o id é um fato sobre o ML; Listings
	// grava Unknown, nunca recusa nem inventa (protocolo §4.1). But "Unknown"
	// here means a deliberate fact.NewUnknown with a reason, not a Go zero
	// value: the zero value is ALSO Unknown-shaped but carries no reason, and
	// that is the malformed case Validate now rejects (0099's CHECK
	// constraints reject it at the database with an opaque error).
	e := testEvidence(t)
	unknownString := func(reason string) fact.Fact[string] {
		f, err := fact.NewUnknown[string](reason, e)
		if err != nil {
			t.Fatalf("fact.NewUnknown: %v", err)
		}
		return f
	}
	unknownInt := func(reason string) fact.Fact[int] {
		f, err := fact.NewUnknown[int](reason, e)
		if err != nil {
			t.Fatalf("fact.NewUnknown: %v", err)
		}
		return f
	}
	unknownMoney := func(reason string) fact.Fact[exact.Money] {
		f, err := fact.NewUnknown[exact.Money](reason, e)
		if err != nil {
			t.Fatalf("fact.NewUnknown: %v", err)
		}
		return f
	}
	obs := contracts.ListingObservation{
		Key:      testKey(t),
		Evidence: e,
		State: contracts.ListingState{
			Title:             unknownString("ml omitted title"),
			Status:            unknownString("ml omitted status"),
			ListingType:       unknownString("ml omitted listing_type_id"),
			Price:             unknownMoney("ml omitted price or currency"),
			AvailableQuantity: unknownInt("ml omitted available_quantity"),
			SellerSKU:         unknownString("ml omitted seller_sku"),
			GTIN:              unknownString("ml omitted gtin attribute"),
		},
	}
	if err := obs.Validate(); err != nil {
		t.Fatalf("all-unknown observation with reasons rejected: %v", err)
	}
}

func TestValidateRejectsZeroValueFactAsMalformed(t *testing.T) {
	// A zero-value fact.Fact[T] is Unknown-shaped but has no reason: that is
	// not the same thing as a deliberate fact.NewUnknown(reason, evidence),
	// and letting it through means the database CHECK constraint (0099) is
	// the first place it is caught, with an opaque error instead of this one.
	obs := contracts.ListingObservation{Key: testKey(t), Evidence: testEvidence(t)}
	if err := obs.Validate(); err == nil {
		t.Fatal("observation with zero-value (reasonless) facts accepted")
	}
}

func TestValidateRejectsWhitespaceOnlyVariationID(t *testing.T) {
	obs := contracts.ListingObservation{
		Key:      testKey(t),
		Evidence: testEvidence(t),
		State:    contracts.ListingState{Variations: []contracts.VariationObservation{{VariationID: "   "}}},
	}
	if err := obs.Validate(); err == nil || !strings.Contains(err.Error(), "variation") {
		t.Fatalf("whitespace-only variation id accepted or wrong error: %v", err)
	}
}

// testState builds a fully-known ListingState so each Equal test only has to
// mutate the one thing it means to test.
func testState(t *testing.T, e provenance.Evidence) contracts.ListingState {
	t.Helper()
	title, err := fact.NewKnown("Titulo", e)
	if err != nil {
		t.Fatalf("fact.NewKnown title: %v", err)
	}
	status, err := fact.NewKnown("active", e)
	if err != nil {
		t.Fatalf("fact.NewKnown status: %v", err)
	}
	listingType, err := fact.NewKnown("gold_special", e)
	if err != nil {
		t.Fatalf("fact.NewKnown listing_type: %v", err)
	}
	currency, err := exact.ParseCurrency("BRL")
	if err != nil {
		t.Fatalf("exact.ParseCurrency: %v", err)
	}
	amount, err := exact.ParseMoney("199.90", currency)
	if err != nil {
		t.Fatalf("exact.ParseMoney: %v", err)
	}
	price, err := fact.NewKnown(amount, e)
	if err != nil {
		t.Fatalf("fact.NewKnown price: %v", err)
	}
	qty, err := fact.NewKnown(10, e)
	if err != nil {
		t.Fatalf("fact.NewKnown qty: %v", err)
	}
	sku, err := fact.NewKnown("SKU-1", e)
	if err != nil {
		t.Fatalf("fact.NewKnown sku: %v", err)
	}
	gtin, err := fact.NewKnown("7890000000001", e)
	if err != nil {
		t.Fatalf("fact.NewKnown gtin: %v", err)
	}
	return contracts.ListingState{
		Title: title, Status: status, ListingType: listingType,
		Price: price, AvailableQuantity: qty, SellerSKU: sku, GTIN: gtin,
		Variations: []contracts.VariationObservation{testVariation(t, e, "V1", "199.90")},
	}
}

func testVariation(t *testing.T, e provenance.Evidence, id, amount string) contracts.VariationObservation {
	t.Helper()
	currency, err := exact.ParseCurrency("BRL")
	if err != nil {
		t.Fatalf("exact.ParseCurrency: %v", err)
	}
	money, err := exact.ParseMoney(amount, currency)
	if err != nil {
		t.Fatalf("exact.ParseMoney: %v", err)
	}
	price, err := fact.NewKnown(money, e)
	if err != nil {
		t.Fatalf("fact.NewKnown variation price: %v", err)
	}
	qty, err := fact.NewKnown(5, e)
	if err != nil {
		t.Fatalf("fact.NewKnown variation qty: %v", err)
	}
	sku, err := fact.NewKnown("SKU-"+id, e)
	if err != nil {
		t.Fatalf("fact.NewKnown variation sku: %v", err)
	}
	gtin, err := fact.NewKnown("7890000001"+id, e)
	if err != nil {
		t.Fatalf("fact.NewKnown variation gtin: %v", err)
	}
	return contracts.VariationObservation{VariationID: id, Price: price, AvailableQuantity: qty, SellerSKU: sku, GTIN: gtin}
}

func TestListingStateEqualSameFactsAreEqual(t *testing.T) {
	e := testEvidence(t)
	a := testState(t, e)
	b := testState(t, e)
	if !a.Equal(b) {
		t.Fatal("identical states compared unequal")
	}
}

func TestListingStateEqualDiffersInState(t *testing.T) {
	e := testEvidence(t)
	a := testState(t, e)
	b := testState(t, e)
	unknownTitle, err := fact.NewUnknown[string]("ml omitted title", e)
	if err != nil {
		t.Fatalf("fact.NewUnknown: %v", err)
	}
	b.Title = unknownTitle
	if a.Equal(b) {
		t.Fatal("states differing only in knowledge state compared equal")
	}
}

func TestListingStateEqualDiffersInReason(t *testing.T) {
	e := testEvidence(t)
	a := testState(t, e)
	b := testState(t, e)
	reasonA, err := fact.NewEstimated("Titulo", "reason a", e)
	if err != nil {
		t.Fatalf("fact.NewEstimated: %v", err)
	}
	reasonB, err := fact.NewEstimated("Titulo", "reason b", e)
	if err != nil {
		t.Fatalf("fact.NewEstimated: %v", err)
	}
	a.Title, b.Title = reasonA, reasonB
	if a.Equal(b) {
		t.Fatal("states differing only in reason compared equal")
	}
}

// TestListingStateEqualDiffersInPriceAmountOnly is the regression test for the
// bug fact.EqualFunc exists to prevent: Money holds a Decimal with a pointer
// inside, so comparing it with == would compare pointer identity, and two
// facts stating different amounts could read as equal simply because they
// were built independently. If Price ever regresses to fact.Equal (which uses
// ==), this test starts failing to fail -- it must FAIL for unequal amounts.
func TestListingStateEqualDiffersInPriceAmountOnly(t *testing.T) {
	e := testEvidence(t)
	a := testState(t, e)
	b := testState(t, e)
	currency, err := exact.ParseCurrency("BRL")
	if err != nil {
		t.Fatalf("exact.ParseCurrency: %v", err)
	}
	otherAmount, err := exact.ParseMoney("299.90", currency)
	if err != nil {
		t.Fatalf("exact.ParseMoney: %v", err)
	}
	otherPrice, err := fact.NewKnown(otherAmount, e)
	if err != nil {
		t.Fatalf("fact.NewKnown: %v", err)
	}
	b.Price = otherPrice
	if a.Equal(b) {
		t.Fatal("states differing only in price amount compared equal")
	}
}

func TestListingStateEqualVariationsDifferInCount(t *testing.T) {
	e := testEvidence(t)
	a := testState(t, e)
	b := testState(t, e)
	b.Variations = append(b.Variations, testVariation(t, e, "V2", "50.00"))
	if a.Equal(b) {
		t.Fatal("states differing in variation count compared equal")
	}
}

// TestListingStateEqualIgnoresVariationOrder pins the one property that keeps
// this comparison usable against the database. Nothing stores a variation's
// position: listings.listing_variations is keyed by variation_id and the
// repository reads it back ORDER BY variation_id, while the channel hands the
// adapter whatever order it likes. Comparing by position would therefore
// disagree with the stored row on every poll of any listing whose channel
// order is not the sorted order, mint a new version, and never converge —
// the same shape of defect as the payload-hash fold, only inverted: a change
// reported where nothing changed.
func TestListingStateEqualIgnoresVariationOrder(t *testing.T) {
	e := testEvidence(t)
	a := testState(t, e)
	a.Variations = []contracts.VariationObservation{
		testVariation(t, e, "V1", "199.90"),
		testVariation(t, e, "V2", "50.00"),
	}
	b := testState(t, e)
	b.Variations = []contracts.VariationObservation{
		testVariation(t, e, "V2", "50.00"),
		testVariation(t, e, "V1", "199.90"),
	}
	if !a.Equal(b) {
		t.Fatal("the same variations in a different order compared unequal — every poll would mint a version")
	}
}

// TestListingStateEqualSeesAChangeInsideAReorderedVariation closes the hole the
// order-insensitive comparison could open: matching by id must still compare
// every fact of the matched variation, not merely confirm the ids line up.
func TestListingStateEqualSeesAChangeInsideAReorderedVariation(t *testing.T) {
	e := testEvidence(t)
	a := testState(t, e)
	a.Variations = []contracts.VariationObservation{
		testVariation(t, e, "V1", "199.90"),
		testVariation(t, e, "V2", "50.00"),
	}
	b := testState(t, e)
	b.Variations = []contracts.VariationObservation{
		testVariation(t, e, "V2", "50.00"),
		testVariation(t, e, "V1", "149.90"),
	}
	if a.Equal(b) {
		t.Fatal("a reordered variation whose price moved compared equal")
	}
}

func TestListingStateEqualVariationsDifferInOneFact(t *testing.T) {
	e := testEvidence(t)
	a := testState(t, e)
	b := testState(t, e)
	otherSKU, err := fact.NewKnown("SKU-DIFFERENT", e)
	if err != nil {
		t.Fatalf("fact.NewKnown: %v", err)
	}
	b.Variations[0].SellerSKU = otherSKU
	if a.Equal(b) {
		t.Fatal("states differing in one variation fact compared equal")
	}
}
