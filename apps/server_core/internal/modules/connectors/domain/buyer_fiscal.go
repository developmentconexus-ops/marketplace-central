package domain

import "time"

// BuyerFiscalInfo is the marketplace-neutral buyer fiscal identity used to register a sale in
// the ERP: legal name, a fiscal document (opaque type + number) and a billing address. No
// provider payload shape (ML "billing_info"/"identification"/"cust_id") leaks past the adapter.
//
// Every field is a pointer — nil is honest absence (ADR-17). ML obfuscates buyer fiscal data
// until it is available and returns 404 for a buyer without billing data; a masked/blank
// provider value degrades to nil, never a fabricated blank.
//
// DocType is the document type string EXACTLY as the provider returns it (e.g. the ML
// identification.type). It is opaque: the decoder never assumes a literal ("CPF"/"CNPJ"/…) nor
// maps it to an enum — the ERP/UI render it as-is (ADR-17).
type BuyerFiscalInfo struct {
	Name      *string             `json:"name"`
	DocType   *string             `json:"doc_type"`
	DocNumber *string             `json:"doc_number"`
	Address   *BuyerFiscalAddress `json:"address"`
	FetchedAt time.Time           `json:"fetched_at"`
}

// BuyerFiscalAddress is the neutral billing address. StateCode/StateName are the provider's
// own state code + label rendered verbatim (no UF derivation — that transform belongs to the
// shipment destination fact, not the fiscal address). Every field is honest-absence nil.
type BuyerFiscalAddress struct {
	StreetName   *string `json:"street_name"`
	StreetNumber *string `json:"street_number"`
	City         *string `json:"city"`
	StateCode    *string `json:"state_code"`
	StateName    *string `json:"state_name"`
	ZipCode      *string `json:"zip_code"`
	CountryID    *string `json:"country_id"`
}

// HasData reports whether any identity field is present. A BuyerFiscalInfo with no name, no
// document and no address is honest absence (a buyer without billing data, or a 404) — the
// caller treats it as nil rather than surfacing an empty fiscal block.
func (b BuyerFiscalInfo) HasData() bool {
	return b.Name != nil || b.DocType != nil || b.DocNumber != nil || b.Address != nil
}
