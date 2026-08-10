// Package postgres is Listings' own writer. Every statement runs inside a
// transaction that first pins app.tenant_id, because RLS is FORCEd and a
// query without the setting sees an empty world (catalog repository.go:49 is
// the ratified precedent).
package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"marketplace-central/apps/server_core/internal/contexts/listings/contracts"
	"marketplace-central/apps/server_core/internal/contexts/listings/internal/application"
	"marketplace-central/apps/server_core/internal/kernel/exact"
	"marketplace-central/apps/server_core/internal/kernel/fact"
)

type Repository struct{ pool *pgxpool.Pool }

func NewRepository(pool *pgxpool.Pool) *Repository { return &Repository{pool: pool} }

// factColumns flattens a Fact[string] into its state/value/reason triple.
func factColumns(f fact.Fact[string]) (state string, value, reason *string) {
	state = f.State().String()
	if v, ok := f.Value(); ok {
		value = &v
	}
	if r := f.Reason(); r != "" {
		reason = &r
	}
	return state, value, reason
}

func intFactColumns(f fact.Fact[int]) (state string, value *int, reason *string) {
	state = f.State().String()
	if v, ok := f.Value(); ok {
		value = &v
	}
	if r := f.Reason(); r != "" {
		reason = &r
	}
	return state, value, reason
}

func moneyFactColumns(f fact.Fact[exact.Money]) (state string, amount, currency, reason *string) {
	state = f.State().String()
	if v, ok := f.Value(); ok {
		// StringFixed, not String: Decimal has no String(), and persistence
		// renders at the column's scale.
		a := v.Amount().StringFixed(2)
		c := v.Currency().String()
		amount, currency = &a, &c
	}
	if r := f.Reason(); r != "" {
		reason = &r
	}
	return state, amount, currency, reason
}

func (r *Repository) withTenantTx(ctx context.Context, tenantID string, fn func(pgx.Tx) error) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("listings postgres: begin: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()
	if _, err := tx.Exec(ctx, "SELECT set_config('app.tenant_id', $1, true)", tenantID); err != nil {
		return fmt.Errorf("listings postgres: set tenant: %w", err)
	}
	if err := fn(tx); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *Repository) Current(ctx context.Context, key contracts.SourceListingKey) (application.CurrentListing, bool, error) {
	var current application.CurrentListing
	found := false
	err := r.withTenantTx(ctx, key.Tenant().String(), func(tx pgx.Tx) error {
		row := tx.QueryRow(ctx,
			`SELECT version, last_payload_hash FROM listings.listings
			 WHERE tenant_id=$1 AND channel=$2 AND account_external_id=$3 AND listing_id=$4`,
			key.Tenant().String(), key.Account().Channel().String(), key.Account().External(), key.ListingID())
		// errors.Is, not ==: a wrapped ErrNoRows read as a real failure would
		// turn a first-ingest into an error, and read as anything else would
		// turn a real failure into "no such listing" — a fabricated absence.
		switch err := row.Scan(&current.Version, &current.PayloadHash); {
		case err == nil:
			found = true
			return nil
		case errors.Is(err, pgx.ErrNoRows):
			return nil
		default:
			return fmt.Errorf("listings postgres: read current: %w", err)
		}
	})
	return current, found, err
}

func (r *Repository) SaveVersion(ctx context.Context, o contracts.ListingObservation, version int) error {
	k := o.Key
	return r.withTenantTx(ctx, k.Tenant().String(), func(tx pgx.Tx) error {
		titleS, titleV, titleR := factColumns(o.Title)
		statusS, statusV, statusR := factColumns(o.Status)
		typeS, typeV, typeR := factColumns(o.ListingType)
		priceS, priceA, priceC, priceR := moneyFactColumns(o.Price)
		qtyS, qtyV, qtyR := intFactColumns(o.AvailableQuantity)
		skuS, skuV, skuR := factColumns(o.SellerSKU)
		gtinS, gtinV, gtinR := factColumns(o.GTIN)
		if _, err := tx.Exec(ctx,
			`INSERT INTO listings.listings (
			   tenant_id, channel, account_external_id, listing_id, version,
			   title_state, title_value, title_reason,
			   status_state, status_value, status_reason,
			   listing_type_state, listing_type_value, listing_type_reason,
			   price_state, price_amount, price_currency, price_reason,
			   available_qty_state, available_qty_value, available_qty_reason,
			   seller_sku_state, seller_sku_value, seller_sku_reason,
			   gtin_state, gtin_value, gtin_reason,
			   last_payload_hash)
			 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26,$27,$28)
			 ON CONFLICT (tenant_id, channel, account_external_id, listing_id) DO UPDATE SET
			   version=EXCLUDED.version,
			   title_state=EXCLUDED.title_state, title_value=EXCLUDED.title_value, title_reason=EXCLUDED.title_reason,
			   status_state=EXCLUDED.status_state, status_value=EXCLUDED.status_value, status_reason=EXCLUDED.status_reason,
			   listing_type_state=EXCLUDED.listing_type_state, listing_type_value=EXCLUDED.listing_type_value, listing_type_reason=EXCLUDED.listing_type_reason,
			   price_state=EXCLUDED.price_state, price_amount=EXCLUDED.price_amount, price_currency=EXCLUDED.price_currency, price_reason=EXCLUDED.price_reason,
			   available_qty_state=EXCLUDED.available_qty_state, available_qty_value=EXCLUDED.available_qty_value, available_qty_reason=EXCLUDED.available_qty_reason,
			   seller_sku_state=EXCLUDED.seller_sku_state, seller_sku_value=EXCLUDED.seller_sku_value, seller_sku_reason=EXCLUDED.seller_sku_reason,
			   gtin_state=EXCLUDED.gtin_state, gtin_value=EXCLUDED.gtin_value, gtin_reason=EXCLUDED.gtin_reason,
			   last_payload_hash=EXCLUDED.last_payload_hash,
			   updated_at=now()`,
			k.Tenant().String(), k.Account().Channel().String(), k.Account().External(), k.ListingID(), version,
			titleS, titleV, titleR, statusS, statusV, statusR, typeS, typeV, typeR,
			priceS, priceA, priceC, priceR, qtyS, qtyV, qtyR, skuS, skuV, skuR, gtinS, gtinV, gtinR,
			o.Evidence.PayloadHash()); err != nil {
			return fmt.Errorf("listings postgres: upsert listing: %w", err)
		}
		if _, err := tx.Exec(ctx,
			`DELETE FROM listings.listing_variations
			 WHERE tenant_id=$1 AND channel=$2 AND account_external_id=$3 AND listing_id=$4`,
			k.Tenant().String(), k.Account().Channel().String(), k.Account().External(), k.ListingID()); err != nil {
			return fmt.Errorf("listings postgres: clear variations: %w", err)
		}
		for _, v := range o.Variations {
			vPriceS, vPriceA, vPriceC, vPriceR := moneyFactColumns(v.Price)
			vQtyS, vQtyV, vQtyR := intFactColumns(v.AvailableQuantity)
			vSkuS, vSkuV, vSkuR := factColumns(v.SellerSKU)
			vGtinS, vGtinV, vGtinR := factColumns(v.GTIN)
			if _, err := tx.Exec(ctx,
				`INSERT INTO listings.listing_variations (
				   tenant_id, channel, account_external_id, listing_id, variation_id,
				   price_state, price_amount, price_currency, price_reason,
				   available_qty_state, available_qty_value, available_qty_reason,
				   seller_sku_state, seller_sku_value, seller_sku_reason,
				   gtin_state, gtin_value, gtin_reason)
				 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18)`,
				k.Tenant().String(), k.Account().Channel().String(), k.Account().External(), k.ListingID(), v.VariationID,
				vPriceS, vPriceA, vPriceC, vPriceR, vQtyS, vQtyV, vQtyR, vSkuS, vSkuV, vSkuR, vGtinS, vGtinV, vGtinR); err != nil {
				return fmt.Errorf("listings postgres: insert variation %s: %w", v.VariationID, err)
			}
		}
		if _, err := tx.Exec(ctx,
			`INSERT INTO listings.source_observations (
			   tenant_id, channel, account_external_id, listing_id, payload_hash, payload, observed_at)
			 VALUES ($1,$2,$3,$4,$5,$6,$7)
			 ON CONFLICT (tenant_id, channel, account_external_id, listing_id, payload_hash) DO NOTHING`,
			k.Tenant().String(), k.Account().Channel().String(), k.Account().External(), k.ListingID(),
			o.Evidence.PayloadHash(), string(o.RawPayload), o.Evidence.ObservedAt()); err != nil {
			// string(o.RawPayload), not the raw []byte: pgx encodes a plain
			// []byte as bytea, which Postgres refuses for a jsonb column.
			// Passing it as text lets pgx pick the jsonb codec.
			return fmt.Errorf("listings postgres: record observation: %w", err)
		}
		return nil
	})
}

var _ application.Store = (*Repository)(nil)
