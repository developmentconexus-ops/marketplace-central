package contracts

import "time"

// StoredObservation is one payload Listings already recorded, handed back out
// so the adapter that produced it can derive facts from it again.
//
// RawPayload stays opaque on the way out exactly as it is on the way in
// (observation.go): this context stores the channel's bytes and never parses
// them. That is what makes a reprocess possible at all — the mapping from
// bytes to facts lives in one place, the channel adapter, and re-running it
// over stored bytes is the same translation, not a second one.
//
// PayloadHash travels with the bytes rather than being recomputed by the
// reader, so the adapter can check that the two columns still describe each
// other. ObservedAt is the moment the channel was read, not the moment of the
// reprocess: nothing is observed when facts are re-derived offline, and a
// fresh timestamp here would claim the channel confirmed today what it said
// days ago.
type StoredObservation struct {
	Key         SourceListingKey
	RawPayload  []byte
	PayloadHash string
	ObservedAt  time.Time
}

// StoredCursor is a position in a walk over stored observations: the key of
// the last row a page returned. The zero value starts at the beginning, the
// same convention port.Cursor uses for a channel feed.
//
// It carries the whole identity — channel, account, listing — and not just the
// listing id, because listing ids are only unique within one channel account.
// A cursor of "the last listing id" would skip or repeat rows the moment a
// second account exists.
type StoredCursor struct {
	channel   string
	account   string
	listingID string
}

// StoredCursorAfter returns the cursor that resumes immediately after key.
func StoredCursorAfter(key SourceListingKey) StoredCursor {
	return StoredCursor{
		channel:   key.Account().Channel().String(),
		account:   key.Account().External(),
		listingID: key.ListingID(),
	}
}

func (c StoredCursor) IsStart() bool     { return c.listingID == "" }
func (c StoredCursor) Channel() string   { return c.channel }
func (c StoredCursor) Account() string   { return c.account }
func (c StoredCursor) ListingID() string { return c.listingID }

// StoredPage is one batch of stored observations plus where to continue. Done
// is explicit for the same reason port.Page's is: the reader decides when the
// walk is over, not the caller's arithmetic on the batch size.
type StoredPage struct {
	Observations []StoredObservation
	Next         StoredCursor
	Done         bool
}
