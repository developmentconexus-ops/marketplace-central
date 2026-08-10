package port

import (
	"marketplace-central/apps/server_core/internal/contexts/listings/contracts"
)

// ListingMapper derives facts from a payload Listings already stored. It is
// the same translation ListingFeed performs on freshly fetched bytes, exposed
// on its own so it can run without the channel.
//
// Two capabilities come out of that separation, and both are why this is a
// port of its own rather than a method on ListingFeed:
//
//   - A mapper correction can be applied to everything already recorded. The
//     stored row is a function of the payload AND the mapping; when only the
//     mapping moves, re-fetching from the channel is neither necessary nor
//     sufficient — a listing that has since been deleted at the source would
//     never be re-fetched, and its stored facts would keep the old mapping's
//     answer forever.
//   - It needs no credential, no quota and no network. An implementation that
//     required a token to re-read bytes we already hold would make the
//     cheapest operation in this system depend on the most fragile one.
type ListingMapper interface {
	MapStored(o contracts.StoredObservation) (contracts.ListingObservation, error)
}
