package listings

import (
	"encoding/json"
	"fmt"
	"strings"

	"marketplace-central/apps/server_core/internal/adapters/marketplace/mercadolivre/internal/api"
	"marketplace-central/apps/server_core/internal/contexts/listings/contracts"
	"marketplace-central/apps/server_core/internal/contexts/listings/port"
	"marketplace-central/apps/server_core/internal/kernel/provenance"
)

// Mapper derives listing facts from bytes Listings already recorded. It holds
// no client, no account and no token on purpose: re-deriving facts from stored
// bytes must not be able to call Mercado Livre, and a type that cannot be
// constructed with a credential cannot spend one by accident.
type Mapper struct{}

func NewMapper() Mapper { return Mapper{} }

// MapStored runs the stored payload through the same translation NextPage runs
// on a fresh one (mapItemState), under the evidence the row already carries.
//
// The two integrity checks are here rather than in the context because only
// this adapter can perform them: the context stores bytes it cannot read, so
// "these bytes are the item this row claims" is a question only the code that
// speaks the channel's wire format can answer.
func (Mapper) MapStored(so contracts.StoredObservation) (contracts.ListingObservation, error) {
	if so.Key.IsZero() {
		return contracts.ListingObservation{}, fmt.Errorf("mercadolivre listings: stored observation has no key")
	}
	// Integrity before interpretation: bytes that do not hash to the hash
	// stored beside them are not evidence of anything, and mapping them would
	// write facts whose provenance points at a payload nobody has.
	if got := payloadHash(so.RawPayload); got != strings.TrimSpace(so.PayloadHash) {
		return contracts.ListingObservation{}, fmt.Errorf(
			"mercadolivre listings: stored payload for %s hashes to %s but its row stores payload_hash %s",
			so.Key.ListingID(), got, so.PayloadHash)
	}
	var item api.Item
	if err := json.Unmarshal(so.RawPayload, &item); err != nil {
		return contracts.ListingObservation{}, fmt.Errorf("mercadolivre listings: decode stored payload for %s: %w", so.Key.ListingID(), err)
	}
	item.Raw = so.RawPayload
	if id := strings.TrimSpace(item.ID); id != so.Key.ListingID() {
		return contracts.ListingObservation{}, fmt.Errorf(
			"mercadolivre listings: payload stored under listing %s names item %q", so.Key.ListingID(), id)
	}
	// so.ObservedAt, never now(): a reprocess observes nothing. The channel was
	// read once, and that moment is the only honest answer to how old these
	// facts are.
	evidence, err := provenance.NewEvidence(systemName, "item", so.Key.ListingID(), so.ObservedAt, so.PayloadHash)
	if err != nil {
		return contracts.ListingObservation{}, fmt.Errorf("mercadolivre listings: rebuild evidence for %s: %w", so.Key.ListingID(), err)
	}
	state, err := mapItemState(item, evidence)
	if err != nil {
		return contracts.ListingObservation{}, fmt.Errorf("mercadolivre listings: map stored %s: %w", so.Key.ListingID(), err)
	}
	return contracts.ListingObservation{Key: so.Key, State: state, RawPayload: so.RawPayload, Evidence: evidence}, nil
}

var _ port.ListingMapper = Mapper{}
