package exact

import (
	"errors"
	"fmt"
	"strings"
)

var (
	// ErrBadCurrency is returned for anything that is not a 3-letter code.
	ErrBadCurrency = errors.New("exact: currency must be a 3-letter code")
	// ErrCurrencyMismatch is returned when two Money values in different
	// currencies are combined.
	ErrCurrencyMismatch = errors.New("exact: currency mismatch")
)

// Currency is an ISO-4217-shaped code, upper-cased. This package validates the
// shape, not membership of the ISO list: the list changes and is not our truth.
type Currency struct {
	code string
}

// ParseCurrency builds a Currency from a three-letter code.
func ParseCurrency(s string) (Currency, error) {
	t := strings.ToUpper(strings.TrimSpace(s))
	if len(t) != 3 {
		return Currency{}, fmt.Errorf("%w: %q", ErrBadCurrency, s)
	}
	for _, c := range t {
		if c < 'A' || c > 'Z' {
			return Currency{}, fmt.Errorf("%w: %q", ErrBadCurrency, s)
		}
	}
	return Currency{code: t}, nil
}

// String returns the code, or the empty string for the zero value.
func (c Currency) String() string { return c.code }

// IsZero reports whether this is the zero value.
func (c Currency) IsZero() bool { return c.code == "" }

// Money is an amount in a currency. There is no exported field and no
// constructor from float64.
type Money struct {
	amount   Decimal
	currency Currency
}

// NewMoney pairs an amount with a currency, rejecting a zero currency. A Money
// without a currency is not an amount, it is a number that lost its meaning.
func NewMoney(amount Decimal, c Currency) (Money, error) {
	if c.IsZero() {
		return Money{}, fmt.Errorf("%w: currency is the zero value", ErrBadCurrency)
	}
	return Money{amount: amount, currency: c}, nil
}

// ParseMoney parses a decimal literal into Money.
func ParseMoney(s string, c Currency) (Money, error) {
	d, err := ParseDecimal(s)
	if err != nil {
		return Money{}, err
	}
	return NewMoney(d, c)
}

// Amount returns the exact amount.
func (m Money) Amount() Decimal { return m.amount }

// Currency returns the currency.
func (m Money) Currency() Currency { return m.currency }

// Add returns m + o, or ErrCurrencyMismatch.
func (m Money) Add(o Money) (Money, error) {
	if m.currency != o.currency {
		return Money{}, fmt.Errorf("%w: %s + %s", ErrCurrencyMismatch, m.currency, o.currency)
	}
	return Money{amount: m.amount.Add(o.amount), currency: m.currency}, nil
}

// Sub returns m - o, or ErrCurrencyMismatch.
func (m Money) Sub(o Money) (Money, error) {
	if m.currency != o.currency {
		return Money{}, fmt.Errorf("%w: %s - %s", ErrCurrencyMismatch, m.currency, o.currency)
	}
	return Money{amount: m.amount.Sub(o.amount), currency: m.currency}, nil
}

// MulDecimal scales an amount by a dimensionless factor — a commission rate, a
// tax rate, a quantity. The currency is unchanged because a rate has none.
func (m Money) MulDecimal(d Decimal) Money {
	return Money{amount: m.amount.Mul(d), currency: m.currency}
}

// String renders the amount at two decimal places with its currency. It is for
// humans and logs; persistence uses Amount().StringFixed at the column's scale.
func (m Money) String() string {
	if m.currency.IsZero() {
		return m.amount.StringFixed(2)
	}
	return m.currency.String() + " " + m.amount.StringFixed(2)
}
