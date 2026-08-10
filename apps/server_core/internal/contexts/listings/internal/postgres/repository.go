// Package postgres is Listings' own writer. Every statement runs inside a
// transaction that first pins app.tenant_id, because RLS is FORCEd and a
// query without the setting sees an empty world (catalog repository.go:49 is
// the ratified precedent).
package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"marketplace-central/apps/server_core/internal/contexts/listings/contracts"
	"marketplace-central/apps/server_core/internal/contexts/listings/internal/application"
	"marketplace-central/apps/server_core/internal/kernel/exact"
	"marketplace-central/apps/server_core/internal/kernel/fact"
	"marketplace-central/apps/server_core/internal/kernel/provenance"
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

// factFromColumns is the inverse of factColumns: it rebuilds a Fact[string]
// from the state/value/reason triple a row stored, using the constructor that
// matches the stored state so a malformed combination (e.g. known with no
// value) fails loud here instead of silently producing a wrong fact.
func factFromColumns(state string, value, reason *string, e provenance.Evidence) (fact.Fact[string], error) {
	switch state {
	case fact.Known.String():
		if value == nil {
			return fact.Fact[string]{}, fmt.Errorf("listings postgres: known fact has no value")
		}
		return fact.NewKnown(*value, e)
	case fact.Estimated.String():
		if value == nil || reason == nil {
			return fact.Fact[string]{}, fmt.Errorf("listings postgres: estimated fact missing value or reason")
		}
		return fact.NewEstimated(*value, *reason, e)
	case fact.Unknown.String():
		if reason == nil {
			return fact.Fact[string]{}, fmt.Errorf("listings postgres: unknown fact has no reason")
		}
		return fact.NewUnknown[string](*reason, e)
	case fact.NotApplicable.String():
		if reason == nil {
			return fact.Fact[string]{}, fmt.Errorf("listings postgres: not_applicable fact has no reason")
		}
		return fact.NewNotApplicable[string](*reason, e)
	default:
		return fact.Fact[string]{}, fmt.Errorf("listings postgres: unrecognized fact state %q", state)
	}
}

// intFactFromColumns is factFromColumns for Fact[int].
func intFactFromColumns(state string, value *int, reason *string, e provenance.Evidence) (fact.Fact[int], error) {
	switch state {
	case fact.Known.String():
		if value == nil {
			return fact.Fact[int]{}, fmt.Errorf("listings postgres: known fact has no value")
		}
		return fact.NewKnown(*value, e)
	case fact.Estimated.String():
		if value == nil || reason == nil {
			return fact.Fact[int]{}, fmt.Errorf("listings postgres: estimated fact missing value or reason")
		}
		return fact.NewEstimated(*value, *reason, e)
	case fact.Unknown.String():
		if reason == nil {
			return fact.Fact[int]{}, fmt.Errorf("listings postgres: unknown fact has no reason")
		}
		return fact.NewUnknown[int](*reason, e)
	case fact.NotApplicable.String():
		if reason == nil {
			return fact.Fact[int]{}, fmt.Errorf("listings postgres: not_applicable fact has no reason")
		}
		return fact.NewNotApplicable[int](*reason, e)
	default:
		return fact.Fact[int]{}, fmt.Errorf("listings postgres: unrecognized fact state %q", state)
	}
}

// moneyFactFromColumns is factFromColumns for Fact[exact.Money]. amount and
// currency are parsed back through the same exact package that rendered them,
// so a round trip through StringFixed(2) reproduces an equal Money, not a
// merely similar one.
func moneyFactFromColumns(state string, amount, currency, reason *string, e provenance.Evidence) (fact.Fact[exact.Money], error) {
	parseMoney := func() (exact.Money, error) {
		if amount == nil || currency == nil {
			return exact.Money{}, fmt.Errorf("listings postgres: money fact missing amount or currency")
		}
		c, err := exact.ParseCurrency(*currency)
		if err != nil {
			return exact.Money{}, err
		}
		return exact.ParseMoney(*amount, c)
	}
	switch state {
	case fact.Known.String():
		m, err := parseMoney()
		if err != nil {
			return fact.Fact[exact.Money]{}, err
		}
		return fact.NewKnown(m, e)
	case fact.Estimated.String():
		if reason == nil {
			return fact.Fact[exact.Money]{}, fmt.Errorf("listings postgres: estimated fact missing reason")
		}
		m, err := parseMoney()
		if err != nil {
			return fact.Fact[exact.Money]{}, err
		}
		return fact.NewEstimated(m, *reason, e)
	case fact.Unknown.String():
		if reason == nil {
			return fact.Fact[exact.Money]{}, fmt.Errorf("listings postgres: unknown fact has no reason")
		}
		return fact.NewUnknown[exact.Money](*reason, e)
	case fact.NotApplicable.String():
		if reason == nil {
			return fact.Fact[exact.Money]{}, fmt.Errorf("listings postgres: not_applicable fact has no reason")
		}
		return fact.NewNotApplicable[exact.Money](*reason, e)
	default:
		return fact.Fact[exact.Money]{}, fmt.Errorf("listings postgres: unrecognized fact state %q", state)
	}
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

// Current reads back the full ListingState the fold compares against, not
// just the version and a hash: application.Service.Ingest calls
// State.Equal(o.State), so a fold that only had the hash to compare would
// have nothing to fold on any more.
//
// Rebuilding a fact.Fact[T] requires non-zero evidence, which this row does
// not carry directly — evidence lives in listings.source_observations,
// identified by last_payload_hash. The LEFT JOIN is deliberate: an INNER JOIN
// would make "no listing" and "listing exists but its observation row is
// gone" both read as zero rows, indistinguishable from each other, and the
// second case is not an absence to report quietly — it means the evidence
// this row's facts are supposedly grounded in cannot be found, which is a
// data-integrity failure this method must fail loud on rather than paper over
// with a fabricated timestamp.
func (r *Repository) Current(ctx context.Context, key contracts.SourceListingKey) (application.CurrentListing, bool, error) {
	var current application.CurrentListing
	found := false
	err := r.withTenantTx(ctx, key.Tenant().String(), func(tx pgx.Tx) error {
		var (
			lastPayloadHash             string
			titleS                      string
			titleVal, titleReason       *string
			statusS                     string
			statusVal, statusReason     *string
			typeS                       string
			typeVal, typeReason         *string
			priceS                      string
			priceA, priceC, priceReason *string
			qtyS                        string
			qtyVal                      *int
			qtyReason                   *string
			skuS                        string
			skuVal, skuReason           *string
			gtinS                       string
			gtinVal, gtinReason         *string
			observedAt                  *time.Time
		)
		row := tx.QueryRow(ctx,
			`SELECT l.version, l.last_payload_hash,
			        l.title_state, l.title_value, l.title_reason,
			        l.status_state, l.status_value, l.status_reason,
			        l.listing_type_state, l.listing_type_value, l.listing_type_reason,
			        l.price_state, l.price_amount, l.price_currency, l.price_reason,
			        l.available_qty_state, l.available_qty_value, l.available_qty_reason,
			        l.seller_sku_state, l.seller_sku_value, l.seller_sku_reason,
			        l.gtin_state, l.gtin_value, l.gtin_reason,
			        so.observed_at
			 FROM listings.listings l
			 LEFT JOIN listings.source_observations so
			   ON so.tenant_id = l.tenant_id AND so.channel = l.channel
			  AND so.account_external_id = l.account_external_id AND so.listing_id = l.listing_id
			  AND so.payload_hash = l.last_payload_hash
			 WHERE l.tenant_id=$1 AND l.channel=$2 AND l.account_external_id=$3 AND l.listing_id=$4`,
			key.Tenant().String(), key.Account().Channel().String(), key.Account().External(), key.ListingID())
		// errors.Is, not ==: a wrapped ErrNoRows read as a real failure would
		// turn a first-ingest into an error, and read as anything else would
		// turn a real failure into "no such listing" — a fabricated absence.
		switch err := row.Scan(
			&current.Version, &lastPayloadHash,
			&titleS, &titleVal, &titleReason,
			&statusS, &statusVal, &statusReason,
			&typeS, &typeVal, &typeReason,
			&priceS, &priceA, &priceC, &priceReason,
			&qtyS, &qtyVal, &qtyReason,
			&skuS, &skuVal, &skuReason,
			&gtinS, &gtinVal, &gtinReason,
			&observedAt,
		); {
		case err == nil:
			// fall through
		case errors.Is(err, pgx.ErrNoRows):
			return nil
		default:
			return fmt.Errorf("listings postgres: read current: %w", err)
		}
		if observedAt == nil {
			return fmt.Errorf(
				"listings postgres: listing %s/%s/%s has last_payload_hash=%s but no matching row in source_observations",
				key.Account().Channel().String(), key.Account().External(), key.ListingID(), lastPayloadHash)
		}
		evidence, err := provenance.NewEvidence(
			key.Account().Channel().String(), "item", key.ListingID(), *observedAt, lastPayloadHash)
		if err != nil {
			return fmt.Errorf("listings postgres: rebuild evidence: %w", err)
		}

		state := contracts.ListingState{}
		if state.Title, err = factFromColumns(titleS, titleVal, titleReason, evidence); err != nil {
			return fmt.Errorf("listings postgres: title: %w", err)
		}
		if state.Status, err = factFromColumns(statusS, statusVal, statusReason, evidence); err != nil {
			return fmt.Errorf("listings postgres: status: %w", err)
		}
		if state.ListingType, err = factFromColumns(typeS, typeVal, typeReason, evidence); err != nil {
			return fmt.Errorf("listings postgres: listing_type: %w", err)
		}
		if state.Price, err = moneyFactFromColumns(priceS, priceA, priceC, priceReason, evidence); err != nil {
			return fmt.Errorf("listings postgres: price: %w", err)
		}
		if state.AvailableQuantity, err = intFactFromColumns(qtyS, qtyVal, qtyReason, evidence); err != nil {
			return fmt.Errorf("listings postgres: available_quantity: %w", err)
		}
		if state.SellerSKU, err = factFromColumns(skuS, skuVal, skuReason, evidence); err != nil {
			return fmt.Errorf("listings postgres: seller_sku: %w", err)
		}
		if state.GTIN, err = factFromColumns(gtinS, gtinVal, gtinReason, evidence); err != nil {
			return fmt.Errorf("listings postgres: gtin: %w", err)
		}

		variations, err := readVariations(ctx, tx, key, evidence)
		if err != nil {
			return err
		}
		state.Variations = variations

		current.State = state
		found = true
		return nil
	})
	return current, found, err
}

// readVariations reads a listing's variations back in a deterministic order.
// SaveVersion writes them in the order o.State.Variations gave them but
// stores no explicit sequence column for that order — variation_id is the
// only ordering key this schema has — so ORDER BY variation_id is the best
// available deterministic read, not a guarantee that it matches the channel's
// original order when two ingests would have produced a genuine reordering.
func readVariations(ctx context.Context, tx pgx.Tx, key contracts.SourceListingKey, evidence provenance.Evidence) ([]contracts.VariationObservation, error) {
	rows, err := tx.Query(ctx,
		`SELECT variation_id,
		        price_state, price_amount, price_currency, price_reason,
		        available_qty_state, available_qty_value, available_qty_reason,
		        seller_sku_state, seller_sku_value, seller_sku_reason,
		        gtin_state, gtin_value, gtin_reason
		 FROM listings.listing_variations
		 WHERE tenant_id=$1 AND channel=$2 AND account_external_id=$3 AND listing_id=$4
		 ORDER BY variation_id`,
		key.Tenant().String(), key.Account().Channel().String(), key.Account().External(), key.ListingID())
	if err != nil {
		return nil, fmt.Errorf("listings postgres: read variations: %w", err)
	}
	defer rows.Close()

	var out []contracts.VariationObservation
	for rows.Next() {
		var (
			variationID                 string
			priceS                      string
			priceA, priceC, priceReason *string
			qtyS                        string
			qtyVal                      *int
			qtyReason                   *string
			skuS                        string
			skuVal, skuReason           *string
			gtinS                       string
			gtinVal, gtinReason         *string
		)
		if err := rows.Scan(
			&variationID,
			&priceS, &priceA, &priceC, &priceReason,
			&qtyS, &qtyVal, &qtyReason,
			&skuS, &skuVal, &skuReason,
			&gtinS, &gtinVal, &gtinReason,
		); err != nil {
			return nil, fmt.Errorf("listings postgres: scan variation: %w", err)
		}
		v := contracts.VariationObservation{VariationID: variationID}
		if v.Price, err = moneyFactFromColumns(priceS, priceA, priceC, priceReason, evidence); err != nil {
			return nil, fmt.Errorf("listings postgres: variation %s price: %w", variationID, err)
		}
		if v.AvailableQuantity, err = intFactFromColumns(qtyS, qtyVal, qtyReason, evidence); err != nil {
			return nil, fmt.Errorf("listings postgres: variation %s available_quantity: %w", variationID, err)
		}
		if v.SellerSKU, err = factFromColumns(skuS, skuVal, skuReason, evidence); err != nil {
			return nil, fmt.Errorf("listings postgres: variation %s seller_sku: %w", variationID, err)
		}
		if v.GTIN, err = factFromColumns(gtinS, gtinVal, gtinReason, evidence); err != nil {
			return nil, fmt.Errorf("listings postgres: variation %s gtin: %w", variationID, err)
		}
		out = append(out, v)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("listings postgres: read variations: %w", err)
	}
	return out, nil
}

func (r *Repository) SaveVersion(ctx context.Context, o contracts.ListingObservation, version int) error {
	k := o.Key
	return r.withTenantTx(ctx, k.Tenant().String(), func(tx pgx.Tx) error {
		titleS, titleV, titleR := factColumns(o.State.Title)
		statusS, statusV, statusR := factColumns(o.State.Status)
		typeS, typeV, typeR := factColumns(o.State.ListingType)
		priceS, priceA, priceC, priceR := moneyFactColumns(o.State.Price)
		qtyS, qtyV, qtyR := intFactColumns(o.State.AvailableQuantity)
		skuS, skuV, skuR := factColumns(o.State.SellerSKU)
		gtinS, gtinV, gtinR := factColumns(o.State.GTIN)
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
		for _, v := range o.State.Variations {
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
			o.Evidence.PayloadHash(), o.RawPayload, o.Evidence.ObservedAt()); err != nil {
			// o.RawPayload as []byte, not as a string: the column is bytea
			// (0099) precisely so the channel's own bytes survive unchanged and
			// still hash to the payload_hash stored beside them. pgx encodes a
			// plain []byte as bytea, so this is a copy, not a parse.
			return fmt.Errorf("listings postgres: record observation: %w", err)
		}
		return nil
	})
}

var _ application.Store = (*Repository)(nil)
