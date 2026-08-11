// Package port carries what Listings asks of a channel source. The cursor is
// opaque on purpose: ML pages by scroll_id, another channel may page by
// timestamp, and neither shape belongs in this contract (the legacy port that
// typed one source's row id into the cursor is the measured counter-example,
// contexts/catalog/port/feed.go:14-18).
package port

import (
	"context"

	"marketplace-central/apps/server_core/internal/contexts/listings/contracts"
	"marketplace-central/apps/server_core/internal/kernel/tenant"
)

// Cursor is the source's own position marker, unreadable by Listings.
type Cursor struct{ token string }

func NewCursor(token string) Cursor { return Cursor{token: token} }
func (c Cursor) Token() string      { return c.token }
func (c Cursor) IsStart() bool      { return c.token == "" }

// Page is one batch of observations plus where to continue. Done is explicit:
// an empty page mid-feed is legal and must not stop the walk.
type Page struct {
	Observations []contracts.ListingObservation
	Next         Cursor
	Done         bool
}

// ListingFeed is a source of listing observations. Listings asks; the adapter
// decides how to page and how to authenticate.
type ListingFeed interface {
	NextPage(ctx context.Context, t tenant.ID, after Cursor, limit int) (Page, error)
}
