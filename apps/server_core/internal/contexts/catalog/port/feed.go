// Package port carries what other contexts may ASK Catalog, and what Catalog
// asks of a source. Both directions are questions, never tables.
package port

import (
	"context"

	"marketplace-central/apps/server_core/internal/contexts/catalog/contracts"
	"marketplace-central/apps/server_core/internal/kernel/tenant"
)

// Cursor is a source's position in its own feed, and Catalog cannot read it.
//
// The token is opaque on purpose. The legacy read port declared
// Cursor{InternalProductID int64} and thereby wrote one ERP's row-id shape into
// the contract: a source paging by string key, by timestamp, or by a provider's
// scroll id could not implement it without lying. Here the source encodes
// whatever it needs and Catalog hands it back untouched.
type Cursor struct{ token string }

// NewCursor builds a cursor from a source-defined token.
func NewCursor(token string) Cursor { return Cursor{token: token} }

// Token returns the source's own position marker.
func (c Cursor) Token() string { return c.token }

// IsStart reports whether this is the beginning of the feed. The zero cursor is
// the start, so a caller that forgets to seed one reads from the beginning
// rather than from an undefined place.
func (c Cursor) IsStart() bool { return c.token == "" }

// Page is one batch of observations plus where to continue.
type Page struct {
	Observations []contracts.ProductObservation
	Next         Cursor
	// Done means the source has no more rows. It is explicit rather than
	// inferred from an empty page, because an empty page in the middle of a
	// filtered feed is legal and must not stop the walk.
	Done bool
}

// ProductFeed is a source of product observations. Catalog asks; the adapter
// decides how to page.
type ProductFeed interface {
	NextPage(ctx context.Context, t tenant.ID, after Cursor, limit int) (Page, error)
}
