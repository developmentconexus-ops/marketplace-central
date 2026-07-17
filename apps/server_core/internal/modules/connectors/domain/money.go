package domain

import (
	"fmt"
	"regexp"
	"strings"
)

// Money is a known monetary value. Unknown money is a nil *Money, never a zero Money.
type Money struct {
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
}

var decimalPattern = regexp.MustCompile(`^[+-]?(?:\d+(?:\.\d+)?|\.\d+)$`)

func ValidateMoney(field string, m *Money) error {
	if m == nil {
		return nil
	}
	if !decimalPattern.MatchString(m.Amount) {
		return fmt.Errorf("%s.amount must be a non-empty decimal string", field)
	}
	if strings.TrimSpace(m.Currency) == "" || m.Currency != strings.ToUpper(m.Currency) {
		return fmt.Errorf("%s.currency must be a non-empty upper-case string", field)
	}
	return nil
}
