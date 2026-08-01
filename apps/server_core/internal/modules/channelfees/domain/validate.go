package domain

import (
	"encoding/json"
	"errors"
	"fmt"
)

// ErrCommissionDetailRequired is returned by ValidateFee when a layer 3
// fee_kind=commission row's detail is missing sale_fee_unit or quantity
// (IC-01 Canonical Examples "Rejeição canônica"). The mandate is scoped to
// commission only — layer 3 fee_kind=freight has no such decomposition and
// a NULL detail is accepted (IC-01 Persistence Expectations; auditoria P5
// r06 F-r06-1). Named so a caller can errors.Is against it instead of
// matching on error text.
var ErrCommissionDetailRequired = errors.New("channelfees: layer 3 commission fee requires detail.sale_fee_unit and detail.quantity")

// ValidateFee applies the one write-time mandate this port owns: layer 3
// commission rows must carry a decomposed detail (sale_fee_unit + quantity)
// so the realized-vs-estimated audit (IC-01 AuditRealizedVsEstimated) can
// read the per-unit fee back out of the row. Every other fee shape —
// including layer 3 freight and any layer 2/1 commission — passes through
// untouched; this function is deliberately narrow (IC-01 "Must Not Decide In
// Feature Execution" pins the rest of the shape upstream).
func ValidateFee(fee Fee) error {
	if fee.Layer != LayerOrder || fee.FeeKind != FeeKindCommission {
		return nil
	}
	if len(fee.Detail) == 0 {
		return fmt.Errorf("%w: detail is empty", ErrCommissionDetailRequired)
	}

	var decoded map[string]json.RawMessage
	if err := json.Unmarshal(fee.Detail, &decoded); err != nil {
		return fmt.Errorf("%w: detail is not a JSON object: %v", ErrCommissionDetailRequired, err)
	}

	var missing []string
	if _, ok := decoded["sale_fee_unit"]; !ok {
		missing = append(missing, "sale_fee_unit")
	}
	if _, ok := decoded["quantity"]; !ok {
		missing = append(missing, "quantity")
	}
	if len(missing) > 0 {
		return fmt.Errorf("%w: missing keys %v", ErrCommissionDetailRequired, missing)
	}
	return nil
}
