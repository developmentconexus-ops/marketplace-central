// Package exact carries the platform's numeric types. There is no constructor
// from float64 anywhere in this package, and that is the point: a binary float
// cannot represent 0.1, and a tax base built from one is wrong before any rule
// is applied. Arithmetic is exact rational arithmetic; rounding happens once,
// at presentation, and is half-to-even.
package exact

import (
	"errors"
	"fmt"
	"math/big"
	"strings"
)

var (
	// ErrNotANumber is returned when input is not a plain decimal literal.
	ErrNotANumber = errors.New("exact: not a decimal number")
	// ErrDivideByZero is returned by Div when the divisor is zero.
	ErrDivideByZero = errors.New("exact: divide by zero")
	// ErrNegativeScale is returned by StringFixed for a negative scale.
	ErrNegativeScale = errors.New("exact: negative scale")
)

// Decimal is an exact rational number. The zero value is zero and is usable.
type Decimal struct {
	// r is nil for the zero value. Every method must treat nil as zero.
	r *big.Rat
}

func (d Decimal) rat() *big.Rat {
	if d.r == nil {
		return new(big.Rat)
	}
	return d.r
}

// ParseDecimal accepts a plain decimal literal: optional sign, digits, an
// optional single dot, digits. It deliberately rejects exponent notation,
// thousands separators, comma decimal marks, NaN and Inf, because each of those
// is a sign the value came from somewhere that has not been read carefully.
func ParseDecimal(s string) (Decimal, error) {
	t := strings.TrimSpace(s)
	if t == "" {
		return Decimal{}, fmt.Errorf("%w: empty", ErrNotANumber)
	}
	body := t
	if body[0] == '+' || body[0] == '-' {
		body = body[1:]
	}
	if body == "" {
		return Decimal{}, fmt.Errorf("%w: %q", ErrNotANumber, s)
	}
	dots := 0
	digits := 0
	for _, c := range body {
		switch {
		case c == '.':
			dots++
			if dots > 1 {
				return Decimal{}, fmt.Errorf("%w: %q", ErrNotANumber, s)
			}
		case c >= '0' && c <= '9':
			digits++
		default:
			return Decimal{}, fmt.Errorf("%w: %q", ErrNotANumber, s)
		}
	}
	if digits == 0 {
		return Decimal{}, fmt.Errorf("%w: %q", ErrNotANumber, s)
	}
	r, ok := new(big.Rat).SetString(t)
	if !ok {
		return Decimal{}, fmt.Errorf("%w: %q", ErrNotANumber, s)
	}
	return Decimal{r: r}, nil
}

// MustParseDecimal is ParseDecimal for compile-time-known literals in tests and
// for constants. It panics on bad input, which is correct only because the
// input is a literal in the source and never data.
func MustParseDecimal(s string) Decimal {
	d, err := ParseDecimal(s)
	if err != nil {
		panic(err)
	}
	return d
}

// FromInt builds a Decimal from a whole number.
func FromInt(n int64) Decimal {
	return Decimal{r: new(big.Rat).SetInt64(n)}
}

// Add returns d + o.
func (d Decimal) Add(o Decimal) Decimal {
	return Decimal{r: new(big.Rat).Add(d.rat(), o.rat())}
}

// Sub returns d - o.
func (d Decimal) Sub(o Decimal) Decimal {
	return Decimal{r: new(big.Rat).Sub(d.rat(), o.rat())}
}

// Mul returns d * o.
func (d Decimal) Mul(o Decimal) Decimal {
	return Decimal{r: new(big.Rat).Mul(d.rat(), o.rat())}
}

// Div returns d / o, exactly, or ErrDivideByZero.
func (d Decimal) Div(o Decimal) (Decimal, error) {
	if o.IsZero() {
		return Decimal{}, ErrDivideByZero
	}
	return Decimal{r: new(big.Rat).Quo(d.rat(), o.rat())}, nil
}

// Neg returns -d.
func (d Decimal) Neg() Decimal {
	return Decimal{r: new(big.Rat).Neg(d.rat())}
}

// Cmp returns -1, 0 or 1 as d is less than, equal to, or greater than o.
func (d Decimal) Cmp(o Decimal) int { return d.rat().Cmp(o.rat()) }

// IsZero reports whether d is exactly zero.
func (d Decimal) IsZero() bool { return d.rat().Sign() == 0 }

// Sign returns -1, 0 or 1.
func (d Decimal) Sign() int { return d.rat().Sign() }

// Rat returns a copy of the underlying exact value. A copy, so no caller can
// mutate a Decimal that another value shares.
func (d Decimal) Rat() *big.Rat { return new(big.Rat).Set(d.rat()) }

// StringFixed renders d with exactly scale digits after the point, rounding
// half-to-even. Half-to-even and not half-up: half-up is biased away from zero,
// and a bias applied to every line of every order is a systematic error, not a
// rounding difference.
func (d Decimal) StringFixed(scale int) string {
	if scale < 0 {
		panic(ErrNegativeScale)
	}
	r := d.rat()
	neg := r.Sign() < 0
	abs := new(big.Rat).Abs(r)

	pow := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(scale)), nil)
	scaled := new(big.Rat).Mul(abs, new(big.Rat).SetInt(pow))

	quo, rem := new(big.Int).QuoRem(scaled.Num(), scaled.Denom(), new(big.Int))
	// Compare 2*rem against denom to find the half.
	twice := new(big.Int).Lsh(rem, 1)
	switch twice.Cmp(scaled.Denom()) {
	case 1:
		quo.Add(quo, big.NewInt(1))
	case 0:
		if quo.Bit(0) == 1 { // odd: round up to make it even
			quo.Add(quo, big.NewInt(1))
		}
	}

	digits := quo.String()
	var out string
	if scale == 0 {
		out = digits
	} else {
		for len(digits) <= scale {
			digits = "0" + digits
		}
		out = digits[:len(digits)-scale] + "." + digits[len(digits)-scale:]
	}
	if neg && strings.Trim(out, "0.") != "" {
		out = "-" + out
	}
	return out
}
