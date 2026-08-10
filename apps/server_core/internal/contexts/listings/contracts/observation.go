package contracts

import (
	"fmt"
	"strings"

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
// the channel. But that only holds for a DELIBERATE unknown — one built
// through fact.NewUnknown/NewNotApplicable, which forces a reason. A
// zero-value fact.Fact[T] is Unknown-shaped too, with no reason, and that is
// not a fact about the channel: it is a value nobody built, and it fails at
// the database instead of here (0099_listings_context.sql's CHECK
// constraints reject a null reason on an unknown title/price). Validate
// closes that gap by rejecting any fact whose state is Unknown or
// NotApplicable and whose reason is blank.
func (o ListingObservation) Validate() error {
	if o.Key.IsZero() {
		return fmt.Errorf("%w: key", ErrBlank)
	}
	if o.Evidence.IsZero() {
		return fmt.Errorf("%w: evidence", ErrBlank)
	}
	for i, v := range o.Variations {
		if strings.TrimSpace(v.VariationID) == "" {
			return fmt.Errorf("%w: variation id at index %d", ErrBlank, i)
		}
	}
	if err := requireFactReason(o.Title, "title"); err != nil {
		return err
	}
	if err := requireFactReason(o.Status, "status"); err != nil {
		return err
	}
	if err := requireFactReason(o.ListingType, "listing_type"); err != nil {
		return err
	}
	if err := requireFactReason(o.Price, "price"); err != nil {
		return err
	}
	if err := requireFactReason(o.AvailableQuantity, "available_quantity"); err != nil {
		return err
	}
	if err := requireFactReason(o.SellerSKU, "seller_sku"); err != nil {
		return err
	}
	if err := requireFactReason(o.GTIN, "gtin"); err != nil {
		return err
	}
	for i, v := range o.Variations {
		if err := requireFactReason(v.Price, fmt.Sprintf("variation[%d] price", i)); err != nil {
			return err
		}
		if err := requireFactReason(v.AvailableQuantity, fmt.Sprintf("variation[%d] available_quantity", i)); err != nil {
			return err
		}
		if err := requireFactReason(v.SellerSKU, fmt.Sprintf("variation[%d] seller_sku", i)); err != nil {
			return err
		}
		if err := requireFactReason(v.GTIN, fmt.Sprintf("variation[%d] gtin", i)); err != nil {
			return err
		}
	}
	return nil
}

// requireFactReason rejects a fact that is Unknown or NotApplicable with a
// blank reason — the zero-value shape that fact.NewUnknown/NewNotApplicable
// cannot produce (both require a non-blank reason) but a struct literal can.
//
// Written as a condition rather than a switch on purpose: a switch over a
// state enum is a claim that every state has been considered, and this rule is
// not about the state space — Known and Estimated are governed elsewhere (they
// must carry a value, which their constructors enforce). A two-arm switch here
// would read as an exhaustiveness bug to anyone adding a fifth state, and the
// exhaustive linter reads it that way too.
func requireFactReason[T any](f fact.Fact[T], field string) error {
	state := f.State()
	needsReason := state == fact.Unknown || state == fact.NotApplicable
	if needsReason && strings.TrimSpace(f.Reason()) == "" {
		return fmt.Errorf("%w: %s fact is %s with no reason", ErrBlank, field, state)
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
