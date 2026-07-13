package domain

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"
)

func TestInternalProductIDRequiresPositiveCODPROD(t *testing.T) {
	got, err := NewInternalProductID(1001)
	if err != nil || got != 1001 {
		t.Fatalf("NewInternalProductID() = %d, %v", got, err)
	}
	if _, err := NewInternalProductID(0); !errors.Is(err, ErrInvalidInternalProductID) {
		t.Fatalf("zero error = %v", err)
	}
	if _, err := NewInternalProductID(-1); !errors.Is(err, ErrInvalidInternalProductID) {
		t.Fatalf("negative error = %v", err)
	}
}

func TestNumericSourceFactDistinguishesUnknownFromKnownZero(t *testing.T) {
	zero := 0.0
	observedAt := time.Date(2026, 7, 13, 12, 0, 0, 0, time.UTC)
	current, err := NewNumericSourceFact("sankhya", &zero, SourceFactQualityCurrent, &observedAt, nil)
	if err != nil {
		t.Fatal(err)
	}
	reason := "cost was not returned by source"
	unknown, err := NewNumericSourceFact("sankhya", nil, SourceFactQualityUnknown, nil, &reason)
	if err != nil {
		t.Fatal(err)
	}
	body, err := json.Marshal(struct {
		Current NumericSourceFact `json:"current"`
		Unknown NumericSourceFact `json:"unknown"`
	}{current, unknown})
	if err != nil {
		t.Fatal(err)
	}
	jsonText := string(body)
	if !strings.Contains(jsonText, `"value":0`) || !strings.Contains(jsonText, `"value":null`) {
		t.Fatalf("expected distinct zero and null: %s", jsonText)
	}
	if _, err := NewNumericSourceFact("sankhya", nil, SourceFactQualityUnknown, nil, nil); !errors.Is(err, ErrInvalidSourceFact) {
		t.Fatalf("missing unknown reason error = %v", err)
	}
	if _, err := NewNumericSourceFact("sankhya", nil, SourceFactQualityCurrent, &observedAt, nil); !errors.Is(err, ErrInvalidSourceFact) {
		t.Fatalf("missing current value error = %v", err)
	}
}

func TestCanonicalProductKeepsCODPRODSeparateFromIdentifiers(t *testing.T) {
	id, _ := NewInternalProductID(1001)
	zero := 0.0
	observedAt := time.Date(2026, 7, 13, 12, 0, 0, 0, time.UTC)
	fact, _ := NewNumericSourceFact("sankhya", &zero, SourceFactQualityCurrent, &observedAt, nil)
	ean, reference, sellerSKU := "7890000000000", "REF-1001", "SELLER-1001"
	product := CanonicalProduct{InternalProductID: id, Name: "Produto", EAN: &ean, ManufacturerReference: &reference, SellerSKU: &sellerSKU, CostAmount: fact, PriceAmount: fact, StockQuantity: fact}
	body, err := json.Marshal(product)
	if err != nil {
		t.Fatal(err)
	}
	for _, fragment := range []string{`"internal_product_id":1001`, `"ean":"7890000000000"`, `"manufacturer_reference":"REF-1001"`, `"seller_sku":"SELLER-1001"`} {
		if !strings.Contains(string(body), fragment) {
			t.Fatalf("missing %s in %s", fragment, body)
		}
	}
}
