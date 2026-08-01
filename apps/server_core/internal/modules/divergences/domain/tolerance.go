package domain

import (
	"fmt"
	"math/big"
)

// Diverges reports whether observed differs from expected by more than
// Kind's tolerance (IC-02 Persistence Expectations): estoque tolerance is 0
// (any difference diverges); tarifa tolerance is R$0.01 absolute per line
// (≤0.01 is rounding, not divergence). Comparison uses exact rational
// arithmetic (math/big.Rat) rather than float64 — this codebase's house
// rule for money math (pricing/domain/decimal.go) applies equally here,
// since a boundary-exact 0.01 delta must never flip due to binary
// floating-point representation error.
func Diverges(kind, expected, observed string) (bool, error) {
	tolerance, err := toleranceFor(kind)
	if err != nil {
		return false, err
	}

	expectedRat, err := parseDecimal(expected)
	if err != nil {
		return false, fmt.Errorf("divergences: expected value: %w", err)
	}
	observedRat, err := parseDecimal(observed)
	if err != nil {
		return false, fmt.Errorf("divergences: observed value: %w", err)
	}

	delta := new(big.Rat).Sub(expectedRat, observedRat)
	if delta.Sign() < 0 {
		delta.Neg(delta)
	}
	return delta.Cmp(tolerance) > 0, nil
}

func toleranceFor(kind string) (*big.Rat, error) {
	switch kind {
	case KindEstoque:
		return big.NewRat(0, 1), nil
	case KindTarifa:
		return big.NewRat(1, 100), nil // R$0.01
	default:
		return nil, fmt.Errorf("divergences: unknown kind %q", kind)
	}
}

func parseDecimal(value string) (*big.Rat, error) {
	rat, ok := new(big.Rat).SetString(value)
	if !ok {
		return nil, fmt.Errorf("%q is not a decimal number", value)
	}
	return rat, nil
}
