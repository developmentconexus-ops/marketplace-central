// Package ports declares the divergences core interfaces (IC-02). M-05/M-06
// producers and the M-03/M-04/M-07 read seams depend on these interfaces,
// never on the Postgres implementation in adapters/postgres directly.
package ports

import (
	"context"

	"marketplace-central/apps/server_core/internal/modules/divergences/domain"
)

// DivergenceRecorder evaluates one (entity, kind) observation and upserts
// the one-open-row ledger (IC-02, ADR-10): a divergent observation opens or
// refreshes the open row (detected_at immutable); a convergent observation
// auto-resolves an open row if one exists.
type DivergenceRecorder interface {
	Evaluate(ctx context.Context, input domain.EvaluateInput) (domain.EvaluateOutcome, error)
}

// DivergenceReader lists the currently-open divergences for one entity,
// newest first by detected_at (IC-02 Operations).
type DivergenceReader interface {
	ListOpenByEntity(ctx context.Context, tenantID, provider, entityType, entityID string) ([]domain.OpenDivergence, error)
}
