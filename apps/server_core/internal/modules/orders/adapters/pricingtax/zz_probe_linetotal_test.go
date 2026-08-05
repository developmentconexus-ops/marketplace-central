package pricingtax

import (
	"math"
	"math/big"
	"testing"

	pricingdomain "marketplace-central/apps/server_core/internal/modules/pricing/domain"
)

// TestProbeLegacyLineTotalFabricatesZero measures the PRE-GUARD lineTotal body
// verbatim (reader.go@b8c929e7:230-237) so the claim "NaN/±Inf became 0.00" is
// measured rather than narrated — the guard replaced that path, so the live
// function can no longer demonstrate it. Probe, not a contract: excluded by
// the lane's -skip 'TestProbe'; delete once nothing cites the legacy behavior.
func TestProbeLegacyLineTotalFabricatesZero(t *testing.T) {
	legacy := func(unitPrice float64, quantity int) string {
		v := new(big.Rat).SetFloat64(unitPrice)
		if v == nil { // unitPrice was NaN or ±Inf — not a real amount
			v = new(big.Rat)
		}
		v.Mul(v, big.NewRat(int64(quantity), 1))
		return pricingdomain.FormatRatHalfUp(v, 2)
	}

	for _, f := range []float64{math.NaN(), math.Inf(1), math.Inf(-1), 1.005, 0.10} {
		t.Logf("legacy lineTotal(%v, 1) = %q", f, legacy(f, 1))
	}
}
