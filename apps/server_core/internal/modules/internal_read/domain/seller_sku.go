package domain

import "strings"

// maxCodprodDigits bounds a CODPROD at what an internal product id can hold
// (int32). A longer run of digits is not a CODPROD, it is another catalog's
// identifier that happens to be numeric.
const maxCodprodDigits = 9

// IsValidCodprod reports whether a provider-supplied seller_sku is
// syntactically a Sankhya CODPROD.
//
// D-121: on Mercado Livre the operator registers the ERP CODPROD in the
// listing's SKU field, so seller_sku is matched against CODPROD — not
// REFFORN, the manufacturer reference. A seller_sku that is NOT a CODPROD
// (legacy REFFORN like "ZP1704.1.", a free-text code, an empty string) must
// never become a match clause: the same discipline IsValidGTIN already
// applies to the EAN, so junk input never widens the query.
func IsValidCodprod(value string) bool {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" || len(trimmed) > maxCodprodDigits {
		return false
	}
	nonZero := false
	for _, r := range trimmed {
		if r < '0' || r > '9' {
			return false
		}
		if r != '0' {
			nonZero = true
		}
	}
	// "0"/"000" parse as digits but name no product — CODPROD is positive.
	return nonZero
}
