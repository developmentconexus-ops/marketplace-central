package domain

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestListingReadModelNullableFieldsRemainNull(t *testing.T) {
	encoded, err := json.Marshal(ListingReadModel{})
	if err != nil {
		t.Fatal(err)
	}
	var got map[string]any
	if err := json.Unmarshal(encoded, &got); err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"listing_type", "price", "published_quantity", "sync_error", "pending_issue", "cost", "below_margin_worst_case", "icms_worst_case_by_uf", "fetched_at"} {
		value, exists := got[key]
		if !exists || value != nil {
			t.Errorf("%s = %#v, exists=%v; want explicit null", key, value, exists)
		}
	}
	// quality_score/sales_30d left the HTTP read contract in Task 8 (ADR-C3).
	for _, key := range []string{"quality_score", "sales_30d"} {
		if _, exists := got[key]; exists {
			t.Errorf("%s must not be on the wire (dropped, no producer)", key)
		}
	}
}

func TestListingTypeLabelsAreRatifiedD5Values(t *testing.T) {
	want := map[ListingTypeCode]string{
		"gold_pro": "Premium", "gold_special": "Clássico", "gold_premium": "Ouro Premium (legado)",
		"gold": "Ouro (legado)", "silver": "Prata (legado)", "bronze": "Bronze (legado)", "free": "Grátis",
	}
	for code, label := range want {
		got, ok := ListingTypeForCode(code)
		if !ok || got.Code != code || got.Label != label {
			t.Errorf("ListingTypeForCode(%q) = %#v, %v; want label %q", code, got, ok, label)
		}
	}
	if _, ok := ListingTypeForCode("unrecognized"); ok {
		t.Fatal("unrecognized modality must remain unknown")
	}
}

func TestParseListingIDMalformedUsesTypedNotFoundPath(t *testing.T) {
	for _, value := range []string{"", "installation~listing", "installation~~-", "installation~listing~-~extra", " installation~listing~-"} {
		_, err := ParseListingID(value)
		if !errors.Is(err, ErrListingNotFound) {
			t.Errorf("ParseListingID(%q) error = %v; want listing_not_found", value, err)
		}
		var typed *ListingNotFoundError
		if !errors.As(err, &typed) {
			t.Errorf("ParseListingID(%q) error type = %T; want *ListingNotFoundError", value, err)
		}
	}
}

func TestBelowMarginWorstCaseDefaultsToUnknown(t *testing.T) {
	// Slice 1 holds the field only; the worst-case tri-state formula (and its full
	// null cost/price/policy/source matrix) lands in the compute slice. Contract
	// here: an unpopulated read model leaves below_margin_worst_case nil (unknown),
	// never a zero-value false (ADR-17).
	if model := (ListingReadModel{}); model.BelowMarginWorstCase != nil {
		t.Fatalf("zero-value below_margin_worst_case = %v; want nil (unknown)", *model.BelowMarginWorstCase)
	}
}
