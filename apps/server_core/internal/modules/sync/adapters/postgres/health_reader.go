package postgres

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"marketplace-central/apps/server_core/internal/modules/sync/application"
	"marketplace-central/apps/server_core/internal/modules/sync/domain"
)

// HealthReader satisfies the /sync/health aggregate's read port (F-01, M-09).
var _ application.HealthReader = (*HealthReader)(nil)

// HealthReader lists every sync_state row for a tenant. Unlike
// SyncStateRepository (scoped to one (installation, entity) key), this reads
// the whole tenant with no installation_id filter — the "erp" sentinel row
// where the products entity lives must not be filtered out (F-01 spec).
type HealthReader struct {
	pool     *pgxpool.Pool
	tenantID string
}

func NewHealthReader(pool *pgxpool.Pool, tenantID string) *HealthReader {
	return &HealthReader{pool: pool, tenantID: tenantID}
}

// ReadAll returns one HealthEntityRow per sync_state row for the tenant,
// ordered by entity ascending, across every installation_id. A tenant with no
// rows returns an empty (non-nil) slice and a nil error.
func (r *HealthReader) ReadAll(ctx context.Context, tenantID string) ([]application.HealthEntityRow, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT entity, cursor, last_full_sync_at, last_incremental_at,
		       last_error, consecutive_failures
		FROM sync_state
		WHERE tenant_id = $1
		ORDER BY entity ASC
	`, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]application.HealthEntityRow, 0)
	for rows.Next() {
		var (
			entity      string
			cursor      []byte
			lastFull    pgtype.Timestamptz
			lastInc     pgtype.Timestamptz
			lastErrRaw  []byte
			consecutive int32
		)
		if err := rows.Scan(&entity, &cursor, &lastFull, &lastInc, &lastErrRaw, &consecutive); err != nil {
			return nil, err
		}

		out := application.HealthEntityRow{Entity: entity, ConsecutiveFailures: int(consecutive)}
		if len(cursor) > 0 {
			out.Cursor = json.RawMessage(cursor)
		}
		if lastFull.Valid {
			t := lastFull.Time.UTC()
			out.LastFullSyncAt = &t
		}
		if lastInc.Valid {
			t := lastInc.Time.UTC()
			out.LastIncrementalAt = &t
		}
		if len(lastErrRaw) > 0 {
			var se domain.SyncError
			if err := json.Unmarshal(lastErrRaw, &se); err == nil {
				out.LastError = &se
			}
		}
		result = append(result, out)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}
