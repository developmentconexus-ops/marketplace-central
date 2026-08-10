// Package contracts is Listings' published vocabulary: what a channel adapter
// hands in, and what ingesting it did. Nothing here names a vendor type.
package contracts

import (
	"errors"
	"fmt"
	"strings"

	"marketplace-central/apps/server_core/internal/kernel/channel"
	"marketplace-central/apps/server_core/internal/kernel/tenant"
)

// ErrBlank marks a required field that arrived empty.
var ErrBlank = errors.New("listings contracts: required field is blank")

// SourceListingKey is the identity of one listing as one channel account
// publishes it. Listings does not mint ids: the provider's listing id IS the
// identity, scoped by tenant and channel account.
type SourceListingKey struct {
	tenant    tenant.ID
	account   channel.AccountRef
	listingID string
}

// NewSourceListingKey validates every component; a key with a blank part would
// silently collide rows across accounts.
func NewSourceListingKey(t tenant.ID, account channel.AccountRef, listingID string) (SourceListingKey, error) {
	if t.IsZero() {
		return SourceListingKey{}, fmt.Errorf("%w: tenant", ErrBlank)
	}
	if account.Channel().IsZero() {
		return SourceListingKey{}, fmt.Errorf("%w: channel account", ErrBlank)
	}
	listingID = strings.TrimSpace(listingID)
	if listingID == "" {
		return SourceListingKey{}, fmt.Errorf("%w: listing id", ErrBlank)
	}
	return SourceListingKey{tenant: t, account: account, listingID: listingID}, nil
}

func (k SourceListingKey) Tenant() tenant.ID           { return k.tenant }
func (k SourceListingKey) Account() channel.AccountRef { return k.account }
func (k SourceListingKey) ListingID() string           { return k.listingID }
func (k SourceListingKey) IsZero() bool                { return k.listingID == "" }
