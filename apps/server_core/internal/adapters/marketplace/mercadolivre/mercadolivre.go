// Package mercadolivre is the vendor root: New() is the only importable
// surface (§2). The wire lives under internal/api and cannot leave.
package mercadolivre

import (
	"context"
	"time"

	"marketplace-central/apps/server_core/internal/adapters/marketplace/mercadolivre/internal/api"
	mllistings "marketplace-central/apps/server_core/internal/adapters/marketplace/mercadolivre/listings"
	listingsport "marketplace-central/apps/server_core/internal/contexts/listings/port"
	"marketplace-central/apps/server_core/internal/kernel/channel"
)

// Config carries what the composition root decides: where ML is, which
// account, and how tokens are obtained. Token is a plain func type, NOT
// api.TokenSource: an internal/ type on the façade would not compile for
// callers outside the vendor tree (§2.2-a). Go assigns this func value to the
// named api.TokenSource implicitly — identical underlying type.
type Config struct {
	BaseURL   string
	UserID    string
	Channel   string
	AccountID string
	Token     func(ctx context.Context) (string, error)
}

type Bundle struct {
	ListingFeed listingsport.ListingFeed
}

// NewListingMapper builds the offline half of this adapter: the translation
// from Mercado Livre item bytes to listing facts, with no client behind it.
//
// It is a second constructor on this root rather than a field of Bundle
// because Bundle's whole content is credential-bound — New() takes a token
// source and fails without one — and a reprocess must run with no credential
// at all. Handing it a Config it does not use, or a token func that exists
// only to be refused, would put a live-call shape around an operation whose
// point is that it makes none.
func NewListingMapper() listingsport.ListingMapper { return mllistings.NewMapper() }

func New(cfg Config) (Bundle, error) {
	client, err := api.NewClient(cfg.BaseURL, cfg.UserID, cfg.Token)
	if err != nil {
		return Bundle{}, err
	}
	code, err := channel.ParseCode(cfg.Channel)
	if err != nil {
		return Bundle{}, err
	}
	account, err := channel.NewAccountRef(code, cfg.AccountID)
	if err != nil {
		return Bundle{}, err
	}
	return Bundle{ListingFeed: mllistings.NewFeed(client, account, time.Now)}, nil
}
