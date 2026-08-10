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
	obs := contracts.ListingObservation{Title: title}
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
		Key:        testKey(t),
		Evidence:   testEvidence(t),
		Variations: []contracts.VariationObservation{{VariationID: ""}},
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
		Key:               testKey(t),
		Evidence:          e,
		Title:             unknownString("ml omitted title"),
		Status:            unknownString("ml omitted status"),
		ListingType:       unknownString("ml omitted listing_type_id"),
		Price:             unknownMoney("ml omitted price or currency"),
		AvailableQuantity: unknownInt("ml omitted available_quantity"),
		SellerSKU:         unknownString("ml omitted seller_sku"),
		GTIN:              unknownString("ml omitted gtin attribute"),
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
		Key:        testKey(t),
		Evidence:   testEvidence(t),
		Variations: []contracts.VariationObservation{{VariationID: "   "}},
	}
	if err := obs.Validate(); err == nil || !strings.Contains(err.Error(), "variation") {
		t.Fatalf("whitespace-only variation id accepted or wrong error: %v", err)
	}
}
