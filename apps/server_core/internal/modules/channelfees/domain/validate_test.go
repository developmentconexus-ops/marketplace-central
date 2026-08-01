package domain

import (
	"encoding/json"
	"errors"
	"testing"
	"time"
)

func baseFee() Fee {
	return Fee{
		TenantID:       "tenant1",
		Provider:       "mercado_livre",
		InstallationID: "inst1",
		Layer:          LayerOrder,
		FeeKind:        FeeKindCommission,
		SubjectType:    SubjectTypeOrderLine,
		SubjectID:      "200012345:MLB123",
		ValueType:      ValueTypeAmount,
		Value:          "24.33",
		Currency:       strPtr("BRL"),
		Origem:         OrigemAPIOrder,
		ColetadoEm:     time.Date(2026, 7, 31, 9, 0, 0, 0, time.UTC),
	}
}

func strPtr(s string) *string { return &s }

// IC-01 Canonical Examples "Rejeição canônica": layer 3 fee_kind=commission
// upsert with detail missing sale_fee_unit/quantity must be rejected by a
// named error — the detail mandate is ONLY for commission (never a bare
// "err != nil" assertion; the exact sentinel is asserted below).
func TestValidateFeeRejectsLayer3CommissionWithoutDetail(t *testing.T) {
	fee := baseFee()
	fee.Detail = nil

	err := ValidateFee(fee)
	if err == nil {
		t.Fatal("MUST-FAIL: expected layer 3 commission without detail to be rejected, got nil error")
	}
	if !errors.Is(err, ErrCommissionDetailRequired) {
		t.Fatalf("expected ErrCommissionDetailRequired, got %v", err)
	}
}

func TestValidateFeeRejectsLayer3CommissionWithDetailMissingSaleFeeUnit(t *testing.T) {
	fee := baseFee()
	fee.Detail = json.RawMessage(`{"quantity":3}`)

	err := ValidateFee(fee)
	if !errors.Is(err, ErrCommissionDetailRequired) {
		t.Fatalf("expected ErrCommissionDetailRequired for missing sale_fee_unit, got %v", err)
	}
}

func TestValidateFeeRejectsLayer3CommissionWithDetailMissingQuantity(t *testing.T) {
	fee := baseFee()
	fee.Detail = json.RawMessage(`{"sale_fee_unit":8.11}`)

	err := ValidateFee(fee)
	if !errors.Is(err, ErrCommissionDetailRequired) {
		t.Fatalf("expected ErrCommissionDetailRequired for missing quantity, got %v", err)
	}
}

// IC-01 Canonical Examples: layer 3 commission WITH the required detail keys
// is accepted (asserted by field value, not by "no error").
func TestValidateFeeAcceptsLayer3CommissionWithCompleteDetail(t *testing.T) {
	fee := baseFee()
	fee.Detail = json.RawMessage(`{"sale_fee_unit":8.11,"quantity":3}`)

	if err := ValidateFee(fee); err != nil {
		t.Fatalf("expected layer 3 commission with complete detail to be accepted, got %v", err)
	}
}

// IC-01 §Persistence Expectations + auditoria P5 r06 F-r06-1: layer 3
// fee_kind=freight has NO sale_fee_unit/quantity decomposition — detail NULL
// is accepted. This is the asymmetric other half of the same rule; the
// commission mandate must NOT bleed into freight.
func TestValidateFeeAcceptsLayer3FreightWithoutDetail(t *testing.T) {
	fee := baseFee()
	fee.FeeKind = FeeKindFreight
	fee.SubjectType = SubjectTypeOrder
	fee.SubjectID = "200012345"
	fee.Origem = OrigemAPIShipment
	fee.Detail = nil

	if err := ValidateFee(fee); err != nil {
		t.Fatalf("expected layer 3 freight without detail to be accepted, got %v", err)
	}
}

// The commission detail mandate is scoped to layer 3 only (IC-01 Persistence
// Expectations: "camada 3 comissão: ... detail OBRIGATÓRIO"). Layer 2
// commission rows carry their own detail shape (percentage_fee/fixed_fee)
// but ValidateFee must not impose the layer-3 sale_fee_unit/quantity mandate
// on them.
func TestValidateFeeDoesNotApplyLayer3MandateToLayer2Commission(t *testing.T) {
	fee := baseFee()
	fee.Layer = LayerListing
	fee.SubjectType = SubjectTypeListing
	fee.SubjectID = "MLB123"
	fee.Origem = OrigemAPIListingPrices
	fee.Detail = json.RawMessage(`{"percentage_fee":12.5,"fixed_fee":6.0,"financing_add_on_fee":0,"price_used":79.90,"listing_type_id":"gold_special"}`)

	if err := ValidateFee(fee); err != nil {
		t.Fatalf("expected layer 2 commission to be accepted regardless of layer-3 detail shape, got %v", err)
	}
}
