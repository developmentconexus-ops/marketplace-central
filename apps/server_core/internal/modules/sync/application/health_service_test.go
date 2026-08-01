package application_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/modules/sync/application"
	"marketplace-central/apps/server_core/internal/modules/sync/domain"
)

type fakeHealthReader struct {
	rows []application.HealthEntityRow
	err  error
}

var _ application.HealthReader = (*fakeHealthReader)(nil)

func (f *fakeHealthReader) ReadAll(context.Context, string) ([]application.HealthEntityRow, error) {
	return f.rows, f.err
}

type fakeWebhookStatsReader struct {
	stats application.WebhookStats
	err   error
}

var _ application.WebhookStatsReader = (*fakeWebhookStatsReader)(nil)

func (f *fakeWebhookStatsReader) WebhookStats(context.Context, string) (application.WebhookStats, error) {
	return f.stats, f.err
}

// TestHealthServiceLastSuccessAtIsGreatestNotFullSyncAlone is the single most
// important test in F-01 (audited finding F-r04-1): an entity that runs
// incremental syncs every few minutes has last_full_sync_at frozen at its
// one-time initial backfill. An implementation that reports last_full_sync_at
// alone would show this healthy, freshly-synced entity as stale/dead — this
// test fails if that regresses.
func TestHealthServiceLastSuccessAtIsGreatestNotFullSyncAlone(t *testing.T) {
	oldFull := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	recentIncremental := time.Date(2026, 7, 31, 11, 0, 0, 0, time.UTC)

	reader := &fakeHealthReader{rows: []application.HealthEntityRow{
		{
			Entity:            "products",
			LastFullSyncAt:    &oldFull,
			LastIncrementalAt: &recentIncremental,
		},
	}}
	svc := application.NewHealthService(reader, "tenant-1")

	result, err := svc.GetHealth(context.Background())
	if err != nil {
		t.Fatalf("GetHealth: %v", err)
	}
	if len(result.Entities) != 1 {
		t.Fatalf("entities = %d, want 1", len(result.Entities))
	}
	got := result.Entities[0].LastSuccessAt
	if got == nil || !got.Equal(recentIncremental) {
		t.Fatalf("last_success_at = %v, want %v (GREATEST must pick the recent incremental, not the stale full sync)", got, recentIncremental)
	}
}

func TestHealthServiceLastSuccessAtNilSafety(t *testing.T) {
	full := time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)
	inc := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name string
		row  application.HealthEntityRow
		want *time.Time
	}{
		{name: "both nil", row: application.HealthEntityRow{Entity: "orders"}, want: nil},
		{name: "only full set", row: application.HealthEntityRow{Entity: "orders", LastFullSyncAt: &full}, want: &full},
		{name: "only incremental set", row: application.HealthEntityRow{Entity: "orders", LastIncrementalAt: &inc}, want: &inc},
		{name: "full is later", row: application.HealthEntityRow{Entity: "orders", LastFullSyncAt: &full, LastIncrementalAt: &inc}, want: &full},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			reader := &fakeHealthReader{rows: []application.HealthEntityRow{tc.row}}
			svc := application.NewHealthService(reader, "tenant-1")
			result, err := svc.GetHealth(context.Background())
			if err != nil {
				t.Fatalf("GetHealth: %v", err)
			}
			got := result.Entities[0].LastSuccessAt
			switch {
			case tc.want == nil && got != nil:
				t.Fatalf("last_success_at = %v, want nil", got)
			case tc.want != nil && (got == nil || !got.Equal(*tc.want)):
				t.Fatalf("last_success_at = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestHealthServiceEmptyTenantReturnsEmptySliceNotNil(t *testing.T) {
	reader := &fakeHealthReader{rows: []application.HealthEntityRow{}}
	svc := application.NewHealthService(reader, "tenant-empty")

	result, err := svc.GetHealth(context.Background())
	if err != nil {
		t.Fatalf("GetHealth: %v", err)
	}
	if result.Entities == nil {
		t.Fatal("entities is nil, want a non-nil empty slice (JSON must encode [], not null)")
	}
	if len(result.Entities) != 0 {
		t.Fatalf("entities = %d, want 0", len(result.Entities))
	}
}

func TestHealthServiceCursorPhaseExtraction(t *testing.T) {
	tests := []struct {
		name   string
		cursor json.RawMessage
		want   *string
	}{
		{name: "nil cursor", cursor: nil, want: nil},
		{name: "valid phase", cursor: json.RawMessage(`{"phase":"incremental","other":"stuff"}`), want: strPtr("incremental")},
		{name: "no phase key", cursor: json.RawMessage(`{"other":"stuff"}`), want: nil},
		{name: "malformed non-JSON", cursor: json.RawMessage(`not json at all {{{`), want: nil},
		{name: "phase null", cursor: json.RawMessage(`{"phase":null}`), want: nil},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			reader := &fakeHealthReader{rows: []application.HealthEntityRow{
				{Entity: "products", Cursor: tc.cursor},
			}}
			svc := application.NewHealthService(reader, "tenant-1")
			result, err := svc.GetHealth(context.Background())
			if err != nil {
				t.Fatalf("GetHealth: %v", err)
			}
			got := result.Entities[0].Phase
			switch {
			case tc.want == nil && got != nil:
				t.Fatalf("phase = %v, want nil", *got)
			case tc.want != nil && (got == nil || *got != *tc.want):
				t.Fatalf("phase = %v, want %v", got, *tc.want)
			}
		})
	}
}

func TestHealthServiceReaderErrorPropagates(t *testing.T) {
	wantErr := errors.New("boom")
	reader := &fakeHealthReader{err: wantErr}
	svc := application.NewHealthService(reader, "tenant-1")

	_, err := svc.GetHealth(context.Background())
	if !errors.Is(err, wantErr) {
		t.Fatalf("err = %v, want %v", err, wantErr)
	}
}

func TestHealthServiceDefaultWebhookIsCanonicalZero(t *testing.T) {
	reader := &fakeHealthReader{rows: []application.HealthEntityRow{}}
	svc := application.NewHealthService(reader, "tenant-1")

	result, err := svc.GetHealth(context.Background())
	if err != nil {
		t.Fatalf("GetHealth: %v", err)
	}
	want := application.WebhookStats{}
	if result.Webhook != want {
		t.Fatalf("webhook = %+v, want canonical zero %+v", result.Webhook, want)
	}
}

// TestHealthServiceWithWebhookStatsReaderMutatesInPlace proves the by-reference
// injection semantics IC-05 requires: calling WithWebhookStatsReader on an
// ALREADY-CONSTRUCTED service must change what THAT SAME service answers on
// its next call — not require the caller to re-wire a new copy into whatever
// already holds a reference to the service (the transport handler, in
// production). A value-receiver builder that returned an updated copy would
// fail this test: `svc` itself would keep answering the default forever.
func TestHealthServiceWithWebhookStatsReaderMutatesInPlace(t *testing.T) {
	reader := &fakeHealthReader{rows: []application.HealthEntityRow{}}
	svc := application.NewHealthService(reader, "tenant-1")

	before, err := svc.GetHealth(context.Background())
	if err != nil {
		t.Fatalf("GetHealth (before): %v", err)
	}
	if before.Webhook != (application.WebhookStats{}) {
		t.Fatalf("webhook before injection = %+v, want canonical zero", before.Webhook)
	}

	notified := time.Date(2026, 8, 1, 9, 0, 0, 0, time.UTC)
	fake := &fakeWebhookStatsReader{stats: application.WebhookStats{
		LastNotificationAt: &notified,
		Pending:            7,
		Dropped24h:         2,
	}}

	// The setter is called on the SAME svc value already in use above — this
	// is the exact shape a future milestone's root.go region will exercise
	// against an already-`.Register(mux)`'d handler.
	if got := svc.WithWebhookStatsReader(fake); got != svc {
		t.Fatalf("WithWebhookStatsReader returned %p, want the same instance %p", got, svc)
	}

	after, err := svc.GetHealth(context.Background())
	if err != nil {
		t.Fatalf("GetHealth (after): %v", err)
	}
	if after.Webhook.LastNotificationAt == nil || !after.Webhook.LastNotificationAt.Equal(notified) {
		t.Fatalf("webhook.last_notification_at = %v, want %v", after.Webhook.LastNotificationAt, notified)
	}
	if after.Webhook.Pending != 7 || after.Webhook.Dropped24h != 2 {
		t.Fatalf("webhook = %+v, want pending=7 dropped_24h=2", after.Webhook)
	}
}

func TestHealthServiceLastErrorPassesThroughVerbatim(t *testing.T) {
	syncErr := &domain.SyncError{Message: "provider timeout", At: time.Date(2026, 7, 30, 8, 0, 0, 0, time.UTC)}
	reader := &fakeHealthReader{rows: []application.HealthEntityRow{
		{Entity: "listings", LastError: syncErr, ConsecutiveFailures: 3},
	}}
	svc := application.NewHealthService(reader, "tenant-1")

	result, err := svc.GetHealth(context.Background())
	if err != nil {
		t.Fatalf("GetHealth: %v", err)
	}
	got := result.Entities[0].LastError
	if got == nil || got.Message != "provider timeout" || !got.At.Equal(syncErr.At) {
		t.Fatalf("last_error = %+v, want %+v", got, syncErr)
	}
	if result.Entities[0].ConsecutiveFailures != 3 {
		t.Fatalf("consecutive_failures = %d, want 3", result.Entities[0].ConsecutiveFailures)
	}
}

func strPtr(s string) *string { return &s }
