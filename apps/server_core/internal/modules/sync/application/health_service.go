package application

import (
	"context"
	"encoding/json"
	"time"

	"marketplace-central/apps/server_core/internal/modules/sync/domain"
)

// HealthEntity is one entity's aggregated health, ready for the transport
// layer to encode (F-01, M-09).
type HealthEntity struct {
	Entity              string
	LastSuccessAt       *time.Time
	LastIncrementalAt   *time.Time
	ConsecutiveFailures int
	Phase               *string
	LastError           *domain.SyncError
}

// HealthResult is the full /sync/health aggregate.
type HealthResult struct {
	Entities []HealthEntity
	Webhook  WebhookStats
}

// HealthService assembles the /sync/health aggregate: per-entity sync_state
// rows plus the webhook stats block behind its own pluggable port.
//
// webhookReader is swapped THROUGH THE POINTER RECEIVER of WithWebhookStatsReader,
// not by a value-receiver builder that returns an updated copy (contrast the
// pricing transport's `WithCalc` idiom at composition/root.go). IC-05 requires
// a later milestone to call WithWebhookStatsReader on the service AFTER it is
// already constructed and `.Register(mux)`'d, and have the already-registered
// route observe the swap — a value-receiver copy would leave the live route
// serving the old reader forever. See health_service_test.go /
// health_handler_test.go for the must-fail this guards.
type HealthService struct {
	reader        HealthReader
	tenantID      string
	webhookReader WebhookStatsReader
}

func NewHealthService(reader HealthReader, tenantID string) *HealthService {
	return &HealthService{
		reader:        reader,
		tenantID:      tenantID,
		webhookReader: defaultWebhookStatsReader{},
	}
}

// WithWebhookStatsReader mutates s in place and returns s for chaining. Because
// it writes through the receiver instead of returning a new struct, a caller
// holding the SAME *HealthService that was already passed to a transport
// handler's constructor changes what that already-registered handler serves.
func (s *HealthService) WithWebhookStatsReader(reader WebhookStatsReader) *HealthService {
	s.webhookReader = reader
	return s
}

// GetHealth reads every sync_state row for the tenant plus the webhook stats
// block and assembles the aggregate. A reader error and a webhook-reader
// error both propagate — this endpoint is read-only and answers no partial
// success.
func (s *HealthService) GetHealth(ctx context.Context) (HealthResult, error) {
	rows, err := s.reader.ReadAll(ctx, s.tenantID)
	if err != nil {
		return HealthResult{}, err
	}

	entities := make([]HealthEntity, 0, len(rows))
	for _, row := range rows {
		entities = append(entities, HealthEntity{
			Entity:              row.Entity,
			LastSuccessAt:       greatestTime(row.LastFullSyncAt, row.LastIncrementalAt),
			LastIncrementalAt:   row.LastIncrementalAt,
			ConsecutiveFailures: row.ConsecutiveFailures,
			Phase:               extractPhase(row.Cursor),
			LastError:           row.LastError,
		})
	}

	webhook, err := s.webhookReader.WebhookStats(ctx, s.tenantID)
	if err != nil {
		return HealthResult{}, err
	}

	return HealthResult{Entities: entities, Webhook: webhook}, nil
}

// greatestTime is the nil-safe GREATEST(last_full_sync_at, last_incremental_at)
// (audited finding F-r04-1): an entity that runs frequent incremental syncs
// has last_full_sync_at frozen at its one-time initial backfill, so
// last_full_sync_at ALONE would report a healthy, freshly-synced entity as
// stale. NULL only when both inputs are NULL.
func greatestTime(fullSync, incremental *time.Time) *time.Time {
	switch {
	case fullSync == nil:
		return incremental
	case incremental == nil:
		return fullSync
	case incremental.After(*fullSync):
		return incremental
	default:
		return fullSync
	}
}

// extractPhase tolerantly reads cursor.phase (IC-06 shape) without
// validating or re-mapping its vocabulary — it passes the string through
// verbatim. A nil, unparseable, or key-less cursor yields nil, never an error
// or a panic.
func extractPhase(cursor json.RawMessage) *string {
	if len(cursor) == 0 {
		return nil
	}
	var payload struct {
		Phase *string `json:"phase"`
	}
	if err := json.Unmarshal(cursor, &payload); err != nil {
		return nil
	}
	return payload.Phase
}
