package application

import (
	"context"
	"time"
)

// WebhookStats is the webhook block of the /sync/health response (F-01,
// M-09). A future milestone (M-08) injects a real WebhookStatsReader against
// its own webhook bookkeeping; today every tenant answers the canonical
// zero/unknown state.
type WebhookStats struct {
	LastNotificationAt *time.Time
	Pending            int
	Dropped24h         int
}

// WebhookStatsReader is the pluggable seam for the webhook block. F-01/M-09
// ships only the default implementation below; M-08 wires a real reader.
type WebhookStatsReader interface {
	WebhookStats(ctx context.Context, tenantID string) (WebhookStats, error)
}

// defaultWebhookStatsReader is the M-09 seam default: the canonical initial
// state, never an error, never a fabricated timestamp.
type defaultWebhookStatsReader struct{}

func (defaultWebhookStatsReader) WebhookStats(context.Context, string) (WebhookStats, error) {
	return WebhookStats{}, nil
}
