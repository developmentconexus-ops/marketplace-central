package domain

import (
	"fmt"
	"math/big"
	"regexp"
)

// Money is a known monetary value. Unknown money is a nil *Money, never a
// zero-valued Money (ADR-17) — callers must check IsUnknown before reading
// Amount.
type Money struct {
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
}

// IsUnknown reports whether m represents an unknown/absent value. A nil
// *Money is unknown; a non-nil Money (including a zero amount) is known.
func (m *Money) IsUnknown() bool {
	return m == nil
}

var decimalStringPattern = regexp.MustCompile(`^[+-]?\d+\.\d+$|^[+-]?\d+$`)

// ParseRat parses a plain decimal string (e.g. "79.00") into an exact
// big.Rat. It rejects anything that is not a plain decimal literal —
// no scientific notation, no thousands separators, no surrounding
// whitespace, no trailing bare "." — so every value entering the pricing
// engine's decimal math is unambiguous.
func ParseRat(s string) (*big.Rat, error) {
	if !decimalStringPattern.MatchString(s) {
		return nil, fmt.Errorf("pricing: %q is not a decimal string", s)
	}
	value, ok := new(big.Rat).SetString(s)
	if !ok {
		return nil, fmt.Errorf("pricing: %q is not a decimal string", s)
	}
	return value, nil
}

// FormatRatHalfUp formats value as a fixed-point decimal string with the
// given number of places, rounding half away from zero at the boundary
// (…x5 rounds up). It is the single formatting routine every pricing money
// and percentage computation must route through (G1) so decimal handling
// never drifts to float64.
func FormatRatHalfUp(value *big.Rat, places int) string {
	scale := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(places)), nil)
	scaledNumerator := new(big.Int).Mul(value.Num(), scale)
	quotient, remainder := new(big.Int).QuoRem(scaledNumerator, value.Denom(), new(big.Int))

	remainder.Abs(remainder)
	denom := new(big.Int).Abs(value.Denom())
	if new(big.Int).Mul(remainder, big.NewInt(2)).Cmp(denom) >= 0 {
		if value.Sign() < 0 {
			quotient.Sub(quotient, big.NewInt(1))
		} else {
			quotient.Add(quotient, big.NewInt(1))
		}
	}

	negative := quotient.Sign() < 0
	digits := new(big.Int).Abs(quotient).String()
	for len(digits) <= places {
		digits = "0" + digits
	}
	result := digits[:len(digits)-places] + "." + digits[len(digits)-places:]
	if negative {
		result = "-" + result
	}
	return result
}
