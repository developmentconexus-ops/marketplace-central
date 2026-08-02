package postgres

import (
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// fakeRow feeds a fixed 11-column sequence into scanReadModel without a
// database, proving the shipping_id + tags scan mapping (ADR-17: absent stays
// absent, never fabricated). Column order mirrors the ListOrders/GetOrder
// SELECT lists in order_repo.go: provider_order_id, provider_code,
// provider_status, provider_status_detail, provider_created_at,
// provider_closed_at, provider_updated_at, nf_state, total, shipping_id,
// tags_json.
type fakeRow struct {
	shippingID           pgtype.Text
	tagsJSON             []byte
	providerStatusDetail pgtype.Text
}

func (f fakeRow) Scan(dest ...any) error {
	*dest[0].(*string) = "ord-1"
	*dest[1].(*string) = "mercado_livre"
	*dest[2].(*string) = "paid"
	*dest[3].(*pgtype.Text) = f.providerStatusDetail
	*dest[4].(*pgtype.Timestamptz) = pgtype.Timestamptz{Time: time.Now(), Valid: true}
	*dest[5].(*pgtype.Timestamptz) = pgtype.Timestamptz{}
	*dest[6].(*pgtype.Timestamptz) = pgtype.Timestamptz{}
	*dest[7].(*pgtype.Text) = pgtype.Text{}
	*dest[8].(*pgtype.Float8) = pgtype.Float8{}
	*dest[9].(*pgtype.Text) = f.shippingID
	*dest[10].(*[]byte) = f.tagsJSON
	return nil
}

func TestScanReadModel_ShippingID(t *testing.T) {
	t.Run("shipping_id populated when present", func(t *testing.T) {
		model, err := scanReadModel(fakeRow{shippingID: pgtype.Text{String: "SHIP-123", Valid: true}})
		if err != nil {
			t.Fatalf("scanReadModel returned error: %v", err)
		}
		if model.ShippingID != "SHIP-123" {
			t.Fatalf("ShippingID = %q, want %q", model.ShippingID, "SHIP-123")
		}
	})

	t.Run("shipping_id absent (NULL) stays empty", func(t *testing.T) {
		model, err := scanReadModel(fakeRow{shippingID: pgtype.Text{}})
		if err != nil {
			t.Fatalf("scanReadModel returned error: %v", err)
		}
		if model.ShippingID != "" {
			t.Fatalf("ShippingID = %q, want empty", model.ShippingID)
		}
	})

	t.Run("tags populated from tags_json (delivered signal for bucketing)", func(t *testing.T) {
		model, err := scanReadModel(fakeRow{tagsJSON: []byte(`["paid","delivered"]`)})
		if err != nil {
			t.Fatalf("scanReadModel returned error: %v", err)
		}
		if len(model.Tags) != 2 || model.Tags[0] != "paid" || model.Tags[1] != "delivered" {
			t.Fatalf("Tags = %v, want [paid delivered]", model.Tags)
		}
	})

	t.Run("tags absent (NULL tags_json) stays nil", func(t *testing.T) {
		model, err := scanReadModel(fakeRow{tagsJSON: nil})
		if err != nil {
			t.Fatalf("scanReadModel returned error: %v", err)
		}
		if len(model.Tags) != 0 {
			t.Fatalf("Tags = %v, want empty", model.Tags)
		}
	})

	// provider_status_detail (0027: text NOT NULL DEFAULT '') can now hold a
	// genuine NULL — this must not panic/error the scan (cannot scan NULL
	// into *string), and must not be confused with a present-but-empty value.
	t.Run("provider_status_detail populated when present", func(t *testing.T) {
		model, err := scanReadModel(fakeRow{providerStatusDetail: pgtype.Text{String: "accredited", Valid: true}})
		if err != nil {
			t.Fatalf("scanReadModel returned error: %v", err)
		}
		if model.ProviderStatusDetail == nil || *model.ProviderStatusDetail != "accredited" {
			t.Fatalf("ProviderStatusDetail = %v, want %q", model.ProviderStatusDetail, "accredited")
		}
	})

	t.Run("provider_status_detail absent (NULL) does not error and stays nil", func(t *testing.T) {
		model, err := scanReadModel(fakeRow{providerStatusDetail: pgtype.Text{}})
		if err != nil {
			t.Fatalf("scanReadModel returned error: %v", err)
		}
		if model.ProviderStatusDetail != nil {
			t.Fatalf("ProviderStatusDetail = %v, want nil", model.ProviderStatusDetail)
		}
	})
}
