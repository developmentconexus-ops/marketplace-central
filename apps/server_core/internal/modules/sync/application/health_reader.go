package application

import (
	"context"
	"encoding/json"
	"time"

	"marketplace-central/apps/server_core/internal/modules/sync/domain"
)

// HealthEntityRow is one sync_state row as read for the /sync/health
// aggregate (F-01, M-09). It is deliberately NOT keyed by installation_id —
// the aggregate reports one entry per row found, across every installation,
// including the "erp" sentinel where the products entity lives.
type HealthEntityRow struct {
	Entity              string
	Cursor              json.RawMessage
	LastFullSyncAt      *time.Time
	LastIncrementalAt   *time.Time
	LastError           *domain.SyncError
	ConsecutiveFailures int
}

// HealthReader lists every sync_state row for a tenant, ordered by entity
// ascending, with no installation_id filter. A tenant with no rows returns an
// empty (non-nil) slice and a nil error — an honest empty state, not a 404.
type HealthReader interface {
	ReadAll(ctx context.Context, tenantID string) ([]HealthEntityRow, error)
}
