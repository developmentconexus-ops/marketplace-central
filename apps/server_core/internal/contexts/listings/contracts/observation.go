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

// ListingState is what we assert about a listing: the seven facts a channel
// reports plus its variations. It is deliberately separate from
// ListingObservation, which adds how we know it (RawPayload, Evidence) — a
// split that exists so the fold in application/ingest.go can ask "did what we
// store change?" without evidence, which differs on every poll by
// construction, drowning out the answer.
type ListingState struct {
	Title             fact.Fact[string]
	Status            fact.Fact[string]
	ListingType       fact.Fact[string]
	Price             fact.Fact[exact.Money]
	AvailableQuantity fact.Fact[int]
	SellerSKU         fact.Fact[string]
	GTIN              fact.Fact[string]
	Variations        []VariationObservation
}

// Validate rejects a state that cannot be recorded. It deliberately accepts
// every fact as Unknown: a channel that said nothing is a fact about the
// channel. But that only holds for a DELIBERATE unknown — one built through
// fact.NewUnknown/NewNotApplicable, which forces a reason. A zero-value
// fact.Fact[T] is Unknown-shaped too, with no reason, and that is not a fact
// about the channel: it is a value nobody built, and it fails at the database
// instead of here (0099_listings_context.sql's CHECK constraints reject a
// null reason on an unknown title/price). Validate closes that gap by
// rejecting any fact whose state is Unknown or NotApplicable and whose reason
// is blank.
func (s ListingState) Validate() error {
	for i, v := range s.Variations {
		if strings.TrimSpace(v.VariationID) == "" {
			return fmt.Errorf("%w: variation id at index %d", ErrBlank, i)
		}
	}
	if err := requireFactReason(s.Title, "title"); err != nil {
		return err
	}
	if err := requireFactReason(s.Status, "status"); err != nil {
		return err
	}
	if err := requireFactReason(s.ListingType, "listing_type"); err != nil {
		return err
	}
	if err := requireFactReason(s.Price, "price"); err != nil {
		return err
	}
	if err := requireFactReason(s.AvailableQuantity, "available_quantity"); err != nil {
		return err
	}
	if err := requireFactReason(s.SellerSKU, "seller_sku"); err != nil {
		return err
	}
	if err := requireFactReason(s.GTIN, "gtin"); err != nil {
		return err
	}
	for i, v := range s.Variations {
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

// Equal reports whether two states assert the same facts. Evidence never
// enters this comparison — ListingState does not carry any — and that is
// deliberately what makes this the right predicate for the ingest fold: two
// observations of the same true listing carry different evidence on every
// poll (a later timestamp, a fresh payload hash) even when nothing about the
// listing changed, and a corrected mapper can turn Unknown into Known for
// bytes that never moved at all. Comparing state instead of evidence is what
// lets that corrected mapper heal a row whose payload hash is identical to
// the one already stored — the property that motivated this method: a live
// ingest re-run after a mapper fix must write every row whose FACTS changed,
// not only the rows whose bytes changed.
func (s ListingState) Equal(other ListingState) bool {
	if !fact.Equal(s.Title, other.Title) {
		return false
	}
	if !fact.Equal(s.Status, other.Status) {
		return false
	}
	if !fact.Equal(s.ListingType, other.ListingType) {
		return false
	}
	// Price uses EqualFunc, not Equal: exact.Money holds a Decimal with a
	// pointer inside, so Fact[Money]'s == would compare pointer identity, not
	// amount, and two facts stating the same price could read as different.
	if !fact.EqualFunc(s.Price, other.Price, moneyEqual) {
		return false
	}
	if !fact.Equal(s.AvailableQuantity, other.AvailableQuantity) {
		return false
	}
	if !fact.Equal(s.SellerSKU, other.SellerSKU) {
		return false
	}
	if !fact.Equal(s.GTIN, other.GTIN) {
		return false
	}
	// Variations match by id, never by position. Nothing stores a variation's
	// position — listings.listing_variations is keyed by variation_id and the
	// repository reads it back sorted by that id, while the channel sends
	// whatever order it likes. Comparing by position would make the stored row
	// and a fresh observation disagree on every poll of any listing whose
	// channel order is not the sorted one, minting a version each time and
	// never converging.
	if len(s.Variations) != len(other.Variations) {
		return false
	}
	byID := make(map[string]VariationObservation, len(other.Variations))
	for _, o := range other.Variations {
		byID[o.VariationID] = o
	}
	for _, v := range s.Variations {
		o, ok := byID[v.VariationID]
		if !ok {
			return false
		}
		// Consume the match so a repeated id on one side cannot pair twice and
		// report equality against a variation the other side never carried.
		delete(byID, v.VariationID)
		if !fact.EqualFunc(v.Price, o.Price, moneyEqual) {
			return false
		}
		if !fact.Equal(v.AvailableQuantity, o.AvailableQuantity) {
			return false
		}
		if !fact.Equal(v.SellerSKU, o.SellerSKU) {
			return false
		}
		if !fact.Equal(v.GTIN, o.GTIN) {
			return false
		}
	}
	return true
}

// moneyEqual compares Money by amount and currency, the two accessors Money
// exposes on purpose instead of exported fields. Amount().Cmp, not ==, because
// Decimal's exact value lives behind a *big.Rat pointer.
func moneyEqual(a, b exact.Money) bool {
	return a.Amount().Cmp(b.Amount()) == 0 && a.Currency() == b.Currency()
}

// ListingObservation is what a channel adapter hands to Listings: one listing
// as the channel saw it at one moment, plus how we know it.
//
// RawPayload is the channel's own bytes for this observation, opaque to this
// context: stored for audit and reconciliation (§15.3 "payloads reais
// gravados"), never parsed past the adapter.
type ListingObservation struct {
	Key        SourceListingKey
	State      ListingState
	RawPayload []byte
	Evidence   provenance.Evidence
}

// Validate rejects an observation that cannot be recorded.
func (o ListingObservation) Validate() error {
	if o.Key.IsZero() {
		return fmt.Errorf("%w: key", ErrBlank)
	}
	if o.Evidence.IsZero() {
		return fmt.Errorf("%w: evidence", ErrBlank)
	}
	return o.State.Validate()
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
