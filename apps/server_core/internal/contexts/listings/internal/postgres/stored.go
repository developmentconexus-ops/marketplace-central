package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"marketplace-central/apps/server_core/internal/contexts/listings/contracts"
	"marketplace-central/apps/server_core/internal/kernel/channel"
	"marketplace-central/apps/server_core/internal/kernel/tenant"
)

// StoredObservations reads back the payload behind each listing's current
// version: the source_observations row whose payload_hash is the listing's
// last_payload_hash. Not every observation ever recorded — re-deriving facts
// from a payload the listing has already moved past would fold an older truth
// in as if it were the newest.
//
// The LEFT JOIN is the same decision Current makes and for the same reason: an
// INNER JOIN would drop a listing whose observation row is missing, and the
// walk would report "reprocessed 33" over 34 listings with nothing anywhere
// saying which one was skipped. Missing evidence is a data-integrity failure,
// so it fails loud, naming the listing.
func (r *Repository) StoredObservations(ctx context.Context, t tenant.ID, after contracts.StoredCursor, limit int) ([]contracts.StoredObservation, error) {
	var out []contracts.StoredObservation
	err := r.withTenantTx(ctx, t.String(), func(tx pgx.Tx) error {
		// The cursor is compared as a row value against the same three columns
		// the ORDER BY uses, so "strictly after" and "in order" are the same
		// definition. The zero cursor is three empty strings, which every real
		// key sorts after — the start needs no special case, and a special case
		// is exactly where a walk loses its first or last row.
		rows, err := tx.Query(ctx,
			`SELECT l.channel, l.account_external_id, l.listing_id,
			        l.last_payload_hash, so.payload, so.observed_at
			 FROM listings.listings l
			 LEFT JOIN listings.source_observations so
			   ON so.tenant_id = l.tenant_id AND so.channel = l.channel
			  AND so.account_external_id = l.account_external_id AND so.listing_id = l.listing_id
			  AND so.payload_hash = l.last_payload_hash
			 WHERE l.tenant_id = $1
			   AND (l.channel, l.account_external_id, l.listing_id) > ($2::text, $3::text, $4::text)
			 ORDER BY l.channel, l.account_external_id, l.listing_id
			 LIMIT $5`,
			t.String(), after.Channel(), after.Account(), after.ListingID(), limit)
		if err != nil {
			return fmt.Errorf("listings postgres: read stored observations: %w", err)
		}
		defer rows.Close()

		for rows.Next() {
			var (
				channelCode, accountID, listingID, payloadHash string
				payload                                        []byte
				observedAt                                     *time.Time
			)
			if err := rows.Scan(&channelCode, &accountID, &listingID, &payloadHash, &payload, &observedAt); err != nil {
				return fmt.Errorf("listings postgres: scan stored observation: %w", err)
			}
			if observedAt == nil {
				return fmt.Errorf(
					"listings postgres: listing %s/%s/%s has last_payload_hash=%s but no matching row in source_observations",
					channelCode, accountID, listingID, payloadHash)
			}
			key, err := storedKey(t, channelCode, accountID, listingID)
			if err != nil {
				return err
			}
			out = append(out, contracts.StoredObservation{
				Key: key,
				// Copied, not aliased: pgx may hand back a slice into the row
				// buffer it reuses on the next Next(), and these bytes must
				// still hash to payloadHash after the walk has moved on.
				RawPayload:  append([]byte(nil), payload...),
				PayloadHash: payloadHash,
				ObservedAt:  *observedAt,
			})
		}
		if err := rows.Err(); err != nil {
			return fmt.Errorf("listings postgres: read stored observations: %w", err)
		}
		return nil
	})
	return out, err
}

// storedKey rebuilds the context's key from the row's columns, through the
// same constructors that validated it on the way in — a row that no longer
// makes a valid key is reported, never carried.
func storedKey(t tenant.ID, channelCode, accountID, listingID string) (contracts.SourceListingKey, error) {
	code, err := channel.ParseCode(channelCode)
	if err != nil {
		return contracts.SourceListingKey{}, fmt.Errorf("listings postgres: listing %s channel: %w", listingID, err)
	}
	account, err := channel.NewAccountRef(code, accountID)
	if err != nil {
		return contracts.SourceListingKey{}, fmt.Errorf("listings postgres: listing %s account: %w", listingID, err)
	}
	key, err := contracts.NewSourceListingKey(t, account, listingID)
	if err != nil {
		return contracts.SourceListingKey{}, fmt.Errorf("listings postgres: listing %s key: %w", listingID, err)
	}
	return key, nil
}
