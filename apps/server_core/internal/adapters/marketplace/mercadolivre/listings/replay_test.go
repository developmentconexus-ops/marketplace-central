package listings_test

import (
	"encoding/json"
	"strings"
	"testing"

	"marketplace-central/apps/server_core/internal/adapters/marketplace/mercadolivre"
	"marketplace-central/apps/server_core/internal/contexts/listings/contracts"
)

// storedFrom is what the listings context hands back to this adapter for a
// reprocess: the bytes it stored, the hash it stored beside them, and when the
// channel was observed. Deliberately built from a live observation, so a case
// cannot drift from the shape the feed actually writes.
func storedFrom(obs contracts.ListingObservation) contracts.StoredObservation {
	return contracts.StoredObservation{
		Key:         obs.Key,
		RawPayload:  obs.RawPayload,
		PayloadHash: obs.Evidence.PayloadHash(),
		ObservedAt:  obs.Evidence.ObservedAt(),
	}
}

// TestMapStoredReproducesTheFactsTheLiveFeedProduced is the property the whole
// reprocess rests on: mapping stored bytes offline must produce exactly the
// facts a live fetch of those same bytes produces. If the two mappers could
// disagree, a reprocess would not heal rows, it would rewrite them with a
// second, quieter answer.
func TestMapStoredReproducesTheFactsTheLiveFeedProduced(t *testing.T) {
	live := observeOne(t, map[string]any{
		"title": "Produto MLB1", "status": "active",
		"price": json.Number("199.90"), "currency_id": "BRL",
		"listing_type_id": "gold_special", "available_quantity": 5,
		"seller_custom_field": "SCF-MLB1",
		"attributes": []map[string]any{
			{"id": "GTIN", "value_name": "7891234567890"},
			{"id": "SELLER_SKU", "value_name": "SKU-MLB1"},
		},
		"variations": []map[string]any{{
			"id": 111, "price": json.Number("199.90"), "available_quantity": 3,
			"attributes": []map[string]any{{"id": "SELLER_SKU", "value_name": "VSKU-MLB1"}},
		}},
	})

	replayed, err := mercadolivre.NewListingMapper().MapStored(storedFrom(live))
	if err != nil {
		t.Fatalf("map stored: %v", err)
	}
	if !replayed.State.Equal(live.State) {
		t.Fatalf("offline mapping disagrees with the live mapping:\n live = %+v\n replay = %+v", live.State, replayed.State)
	}
	if err := replayed.Validate(); err != nil {
		t.Fatalf("replayed observation does not validate: %v", err)
	}
}

// TestMapStoredKeepsTheOriginalEvidence pins that a reprocess does not
// fabricate a fresh observation. Nothing was observed: the channel was never
// called. Stamping now() here would claim ML confirmed these facts today, and
// the row's observed_at is the only thing that says how old the truth is.
func TestMapStoredKeepsTheOriginalEvidence(t *testing.T) {
	live := observeOne(t, map[string]any{"title": "Produto MLB1"})
	stored := storedFrom(live)

	replayed, err := mercadolivre.NewListingMapper().MapStored(stored)
	if err != nil {
		t.Fatalf("map stored: %v", err)
	}
	if got := replayed.Evidence.ObservedAt(); !got.Equal(stored.ObservedAt) {
		t.Fatalf("observed_at = %s, want the stored %s", got, stored.ObservedAt)
	}
	if got := replayed.Evidence.PayloadHash(); got != stored.PayloadHash {
		t.Fatalf("payload_hash = %s, want the stored %s", got, stored.PayloadHash)
	}
	if string(replayed.RawPayload) != string(stored.RawPayload) {
		t.Fatal("replayed observation did not carry the stored bytes through unchanged")
	}
}

// TestMapStoredRefusesBytesThatDoNotHashToTheStoredHash: the hash and the bytes
// are stored in two columns of the same row, and the whole reason the payload
// column is bytea is that sha256(payload) must reproduce payload_hash. A
// reprocess reads both, so it is the cheapest place in the system to notice
// that they stopped agreeing — and mapping bytes under a hash that does not
// describe them would write facts whose evidence points at a different payload.
func TestMapStoredRefusesBytesThatDoNotHashToTheStoredHash(t *testing.T) {
	stored := storedFrom(observeOne(t, map[string]any{"title": "Produto MLB1"}))
	stored.PayloadHash = strings.Repeat("0", 64)

	if _, err := mercadolivre.NewListingMapper().MapStored(stored); err == nil {
		t.Fatal("a payload that does not hash to its stored hash was mapped anyway")
	} else if !strings.Contains(err.Error(), "hash") {
		t.Fatalf("error must name the hash mismatch, got: %v", err)
	}
}

// TestMapStoredRefusesAPayloadThatNamesADifferentListing closes the other half
// of the same seam: the key comes from the row, the id comes from the bytes,
// and if they disagree the facts would land on the wrong listing.
func TestMapStoredRefusesAPayloadThatNamesADifferentListing(t *testing.T) {
	live := observeOne(t, map[string]any{"title": "Produto MLB1"})
	stored := storedFrom(live)
	// Same bytes and hash as live, but filed under another listing's key.
	otherKey, err := contracts.NewSourceListingKey(live.Key.Tenant(), live.Key.Account(), "MLB999")
	if err != nil {
		t.Fatalf("key: %v", err)
	}
	stored.Key = otherKey

	if _, err := mercadolivre.NewListingMapper().MapStored(stored); err == nil {
		t.Fatal("a payload was mapped onto a listing id it does not name")
	} else if !strings.Contains(err.Error(), "MLB999") || !strings.Contains(err.Error(), "MLB1") {
		t.Fatalf("error must name both ids, got: %v", err)
	}
}

// TestMapStoredRefusesUnparseableBytes: a stored payload that is not JSON is a
// data-integrity failure, not an empty listing. Mapping it to a row of Unknown
// facts would erase real ones under a reason that reads like the channel's.
func TestMapStoredRefusesUnparseableBytes(t *testing.T) {
	stored := storedFrom(observeOne(t, map[string]any{"title": "Produto MLB1"}))
	stored.RawPayload = []byte("{not json")

	if _, err := mercadolivre.NewListingMapper().MapStored(stored); err == nil {
		t.Fatal("unparseable stored bytes were mapped instead of refused")
	}
}
