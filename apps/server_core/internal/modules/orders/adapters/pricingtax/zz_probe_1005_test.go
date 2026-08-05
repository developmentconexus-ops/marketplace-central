package pricingtax

import (
	"math/big"
	"strconv"
	"testing"
)

// TestProbe1005ExactValue is the measurement lineTotal's doc comment cites.
//
// The review of this slice caught a false claim in that comment: it said the
// "decimal-exact answer" for lineTotal(1.005, 1) was "1.01". It is not.
// float64(1.005) sits BELOW the half-up boundary, so "1.00" is the faithful
// answer for that double, and the real defect is the dropped third decimal —
// not a rounding error. Kept as a probe so the number in the comment stays
// checkable instead of being taken on trust. Excluded by -skip 'TestProbe'.
func TestProbe1005ExactValue(t *testing.T) {
	const f = 1.005
	exact := new(big.Rat).SetFloat64(f)
	t.Logf("float64(1.005) exact = %s", exact.FloatString(40))
	t.Logf("above the 1.005 decimal? %v", exact.Cmp(big.NewRat(1005, 1000)) > 0)
	t.Logf("shortest repr = %s", strconv.FormatFloat(f, 'g', -1, 64))
}
