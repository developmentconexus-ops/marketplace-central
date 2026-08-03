package application

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	syncapp "marketplace-central/apps/server_core/internal/modules/sync/application"
)

// staleAggregateReader is the consumer-side port for enumerating market
// evidence that has passed its TTL. *postgres.Repository (via
// ports.MarketAggregateStore, Task 6) satisfies this structurally.
type staleAggregateReader interface {
	StaleAggregateProductIDs(ctx context.Context, olderThan time.Duration, limit int) ([]string, error)
}

// productCollector is the consumer-side port for renewing one product's market
// evidence. *CollectionPipelineService satisfies this structurally.
type productCollector interface {
	Collect(ctx context.Context, codprod string) (CollectionSummary, error)
}

// NewCollectionJob returns the sync.JobFunc that renews market evidence that has
// passed its TTL. It is a JobFunc, not a scheduler: the syncapp.Scheduler that
// already exists (sync/application/scheduler.go) is what ticks, reads/writes the
// cursor, and reconciles sync_state — registering this job there (Task 8, not this
// one) is what makes a collection failure reach SyncHealthCard via RecordFailure.
//
// The job does NOT sweep the catalog — sweeping the ~2,923-product sellable
// catalog at up to 6 ML calls per product (CollectionPipelineService.Collect)
// would cost ~17,500 calls and starve every other sync job sharing the 900/min
// ML rate-limit bucket for ~20 minutes. Instead it only renews evidence that
// was already collected and has now gone stale (via StaleAggregateProductIDs),
// bounded by ceiling per cycle. Discovering new products stays a manual
// operator action, unrelated to this job.
func NewCollectionJob(
	reader staleAggregateReader,
	collector productCollector,
	ttl time.Duration,
	ceiling int,
	logger *slog.Logger,
) syncapp.JobFunc {
	return func(ctx context.Context, cursor json.RawMessage) (json.RawMessage, error) {
		ids, err := reader.StaleAggregateProductIDs(ctx, ttl, ceiling)
		if err != nil {
			logger.Error("market collection job: enumeração falhou", "err", err)
			return cursor, err
		}
		if len(ids) == 0 {
			return cursor, nil
		}

		var failures int
		for _, id := range ids {
			if ctx.Err() != nil {
				return cursor, ctx.Err()
			}
			if _, err := collector.Collect(ctx, id); err != nil {
				failures++
				logger.Error("market collection job: produto falhou", "codprod", id, "err", err)
			}
		}
		logger.Info("market collection job: ciclo concluído",
			"colhidos", len(ids)-failures, "falhas", failures, "teto", ceiling)

		// A per-product failure does not fail the cycle: one bad codprod would
		// otherwise mark the whole entity broken. What failed stays stale and is
		// picked up again by the next enumeration on its own.
		//
		// Cursor deliberately unchanged: this job's position is derived from
		// computed_at in the database, not from a cursor carried between cycles —
		// writing one here would be a second, competing source of truth.
		return cursor, nil
	}
}
