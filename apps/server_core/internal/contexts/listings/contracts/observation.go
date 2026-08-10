package contracts

import (
	"fmt"

	"marketplace-central/apps/server_core/internal/kernel/exact"
	"marketplace-central/apps/server_core/internal/kernel/fact"
	"marketplace-central/apps/server_core/internal/kernel/provenance"
)

// VariationObservation is one variation as the channel reported it. SellerSKU
// and GTIN travel here because the linking context (next leg) anchors on them;
// dropping them now would force a second pull later.
type VariationObservation struct {
	VariationID       string
	Price             fact.Fact[exact.Money]
	AvailableQuantity fact.Fact[int]
	SellerSKU         fact.Fact[string]
	GTIN              fact.Fact[string]
}

// ListingObservation is what a channel adapter hands to Listings: one listing
// as the channel saw it at one moment. Every field the channel may omit is a
// Fact, never a zero value (protocolo §4.1).
//
// RawPayload is the channel's own bytes for this observation, opaque to this
// context: stored for audit and reconciliation (§15.3 "payloads reais
// gravados"), never parsed past the adapter.
type ListingObservation struct {
	Key               SourceListingKey
	Title             fact.Fact[string]
	Status            fact.Fact[string]
	ListingType       fact.Fact[string]
	Price             fact.Fact[exact.Money]
	AvailableQuantity fact.Fact[int]
	SellerSKU         fact.Fact[string]
	GTIN              fact.Fact[string]
	Variations        []VariationObservation
	RawPayload        []byte
	Evidence          provenance.Evidence
}

// Validate rejects an observation that cannot be recorded. It deliberately
// accepts every fact as Unknown: a channel that said nothing is a fact about
// the channel.
func (o ListingObservation) Validate() error {
	if o.Key.IsZero() {
		return fmt.Errorf("%w: key", ErrBlank)
	}
	if o.Evidence.IsZero() {
		return fmt.Errorf("%w: evidence", ErrBlank)
	}
	for i, v := range o.Variations {
		if v.VariationID == "" {
			return fmt.Errorf("%w: variation id at index %d", ErrBlank, i)
		}
	}
	return nil
}

// Disposition is what ingesting an observation did — same closed set the
// catalog leg ratified (contexts/catalog/contracts/observation.go:39-49).
type Disposition string

const (
	DispositionCreated    Disposition = "created"
	DispositionChanged    Disposition = "changed"
	DispositionIdempotent Disposition = "idempotent"
)

// IngestResult reports what happened to one observation.
type IngestResult struct {
	Disposition Disposition
	Version     int
}
