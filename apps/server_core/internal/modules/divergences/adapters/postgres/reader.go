package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"marketplace-central/apps/server_core/internal/modules/divergences/domain"
	"marketplace-central/apps/server_core/internal/modules/divergences/ports"
)

// Reader is the Postgres-backed ports.DivergenceReader.
type Reader struct {
	pool *pgxpool.Pool
}

// NewReader builds a Reader over pool.
func NewReader(pool *pgxpool.Pool) *Reader {
	return &Reader{pool: pool}
}

var _ ports.DivergenceReader = (*Reader)(nil)

const listOpenByEntitySQL = `
SELECT kind, expected_value, observed_value, expected_source, observed_source,
	expected_observed_at, observed_observed_at, detected_at
FROM divergences
WHERE tenant_id = $1 AND provider = $2 AND entity_type = $3 AND entity_id = $4
  AND resolved_at IS NULL
ORDER BY detected_at DESC`

// ListOpenByEntity returns every currently-open (resolved_at IS NULL)
// divergence for one entity, newest first by detected_at (IC-02 Operations).
// An empty slice (never nil-vs-empty ambiguity for the DTO layer) means no
// open divergence.
func (r *Reader) ListOpenByEntity(ctx context.Context, tenantID, provider, entityType, entityID string) ([]domain.OpenDivergence, error) {
	rows, err := r.pool.Query(ctx, listOpenByEntitySQL, tenantID, provider, entityType, entityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	open := []domain.OpenDivergence{}
	for rows.Next() {
		var d domain.OpenDivergence
		if err := rows.Scan(&d.Kind, &d.ExpectedValue, &d.ObservedValue, &d.ExpectedSource, &d.ObservedSource,
			&d.ExpectedObservedAt, &d.ObservedObservedAt, &d.DetectedAt); err != nil {
			return nil, err
		}
		open = append(open, d)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return open, nil
}
