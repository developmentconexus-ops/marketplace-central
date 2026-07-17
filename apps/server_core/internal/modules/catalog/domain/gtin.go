package domain

import "strings"

// IsValidGTIN reports whether value is a syntactically valid GTIN.
func IsValidGTIN(value string) bool {
	value = strings.TrimSpace(value)
	switch len(value) {
	case 8, 12, 13, 14:
	default:
		return false
	}

	for i := 0; i < len(value); i++ {
		if value[i] < '0' || value[i] > '9' {
			return false
		}
	}

	sum := 0
	weight := 3
	for i := len(value) - 2; i >= 0; i-- {
		sum += int(value[i]-'0') * weight
		if weight == 3 {
			weight = 1
		} else {
			weight = 3
		}
	}

	checkDigit := (10 - (sum % 10)) % 10
	return int(value[len(value)-1]-'0') == checkDigit
}
