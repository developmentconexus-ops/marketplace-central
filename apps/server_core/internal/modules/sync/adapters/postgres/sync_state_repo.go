// Package postgres implements the sync_state persistence seam against PostgreSQL.
// Every statement is scoped to the repository's tenant_id (AC-01); NULL columns
// scan to nil/absent Go values (ADR-17); the cursor-append is a single atomic
// UPDATE (M01-C11).
package postgres

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"marketplace-central/apps/server_core/internal/modules/sync/application"
	"marketplace-central/apps/server_core/internal/modules/sync/domain"
)

// SyncStateRepository satisfies the scheduler's StateStore seam.
var _ application.StateStore = (*SyncStateRepository)(nil)

// SyncStateRepository reads/writes sync_state for a single tenant.
type SyncStateRepository struct {
	pool     *pgxpool.Pool
	tenantID string
}

func NewSyncStateRepository(pool *pgxpool.Pool, tenantID string) *SyncStateRepository {
	return &SyncStateRepository{pool: pool, tenantID: tenantID}
}

// Read returns the state for (installation, entity). A missing row is found=false
// with no error and no fabricated cursor — an honest first run.
func (r *SyncStateRepository) Read(ctx context.Context, installationID string, entity domain.Entity) (domain.SyncState, bool, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT tenant_id, installation_id, entity,
		       cursor, schedule, last_full_sync_at, last_incremental_at,
		       last_error, consecutive_failures
		FROM sync_state
		WHERE tenant_id = $1
		  AND installation_id = $2
		  AND entity = $3
	`, r.tenantID, installationID, string(entity))

	state, err := scanSyncState(row)
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.SyncState{}, false, nil
		}
		return domain.SyncState{}, false, err
	}
	return state, true, nil
}

// RecordSuccess upserts the cursor + the relevant timestamp, clears last_error,
// and resets consecutive_failures to 0. A nil cursor is stored as SQL NULL, never
// a fabricated `{}` (ADR-17).
func (r *SyncStateRepository) RecordSuccess(ctx context.Context, installationID string, entity domain.Entity, cursor json.RawMessage, at time.Time, incremental bool) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO sync_state (
			tenant_id, installation_id, entity,
			cursor, last_full_sync_at, last_incremental_at,
			last_error, consecutive_failures
		) VALUES ($1, $2, $3, $4, $5, $6, NULL, 0)
		ON CONFLICT (tenant_id, installation_id, entity) DO UPDATE SET
			cursor = EXCLUDED.cursor,
			last_full_sync_at = COALESCE(EXCLUDED.last_full_sync_at, sync_state.last_full_sync_at),
			last_incremental_at = COALESCE(EXCLUDED.last_incremental_at, sync_state.last_incremental_at),
			last_error = NULL,
			consecutive_failures = 0
	`, r.tenantID, installationID, string(entity),
		jsonbArg(cursor), timestampArg(!incremental, at), timestampArg(incremental, at))
	return err
}

// RecordFailure upserts last_error and increments consecutive_failures by 1 in a
// single atomic statement (read-modify-write in the DB, never in the app — no
// lost update). It leaves the success timestamps untouched.
func (r *SyncStateRepository) RecordFailure(ctx context.Context, installationID string, entity domain.Entity, syncErr domain.SyncError) error {
	payload, err := json.Marshal(syncErr)
	if err != nil {
		return err
	}
	_, err = r.pool.Exec(ctx, `
		INSERT INTO sync_state (
			tenant_id, installation_id, entity,
			last_error, consecutive_failures
		) VALUES ($1, $2, $3, $4, 1)
		ON CONFLICT (tenant_id, installation_id, entity) DO UPDATE SET
			last_error = EXCLUDED.last_error,
			consecutive_failures = sync_state.consecutive_failures + 1
	`, r.tenantID, installationID, string(entity), payload)
	return err
}

// AppendPendingCodigo atomically appends a product code to the entity cursor's
// {"pending":[...]} array in a single UPDATE — no read-modify-write in the app —
// so two concurrent producers (e.g. an M-03 import hook and an M-07 enqueue) never
// lose an update (M01-C11). This operation is owned by M-01 (owner of sync_state);
// consumers call it, they do not reimplement it. The row is created on first
// append.
func (r *SyncStateRepository) AppendPendingCodigo(ctx context.Context, installationID string, entity domain.Entity, codigo string) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO sync_state (
			tenant_id, installation_id, entity, cursor, consecutive_failures
		) VALUES (
			$1, $2, $3, jsonb_build_object('pending', jsonb_build_array($4::text)), 0
		)
		ON CONFLICT (tenant_id, installation_id, entity) DO UPDATE SET
			cursor = jsonb_set(
				COALESCE(sync_state.cursor, '{}'::jsonb),
				'{pending}',
				COALESCE(sync_state.cursor -> 'pending', '[]'::jsonb) || to_jsonb($4::text),
				true
			)
	`, r.tenantID, installationID, string(entity), codigo)
	return err
}

// jsonbArg maps a nil/empty RawMessage to a SQL NULL param (ADR-17) and a
// non-empty one to its bytes.
func jsonbArg(raw json.RawMessage) any {
	if len(raw) == 0 {
		return nil
	}
	return []byte(raw)
}

// timestampArg returns at (UTC) when set, otherwise a SQL NULL param.
func timestampArg(set bool, at time.Time) any {
	if set {
		return at.UTC()
	}
	return nil
}

func scanSyncState(row pgx.Row) (domain.SyncState, error) {
	var (
		st          domain.SyncState
		entity      string
		cursor      []byte
		schedule    []byte
		lastFull    pgtype.Timestamptz
		lastInc     pgtype.Timestamptz
		lastErrRaw  []byte
		consecutive int32
	)
	if err := row.Scan(
		&st.TenantID, &st.InstallationID, &entity,
		&cursor, &schedule, &lastFull, &lastInc, &lastErrRaw, &consecutive,
	); err != nil {
		return domain.SyncState{}, err
	}

	st.Entity = domain.Entity(entity)
	if len(cursor) > 0 {
		st.Cursor = json.RawMessage(cursor)
	}
	if len(schedule) > 0 {
		st.Schedule = json.RawMessage(schedule)
	}
	if lastFull.Valid {
		t := lastFull.Time.UTC()
		st.LastFullSyncAt = &t
	}
	if lastInc.Valid {
		t := lastInc.Time.UTC()
		st.LastIncrementalAt = &t
	}
	if len(lastErrRaw) > 0 {
		var se domain.SyncError
		if err := json.Unmarshal(lastErrRaw, &se); err == nil {
			st.LastError = &se
		}
	}
	st.ConsecutiveFailures = int(consecutive)
	return st, nil
}
