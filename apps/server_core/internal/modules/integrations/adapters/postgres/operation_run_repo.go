package postgres

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"marketplace-central/apps/server_core/internal/modules/integrations/domain"
	"marketplace-central/apps/server_core/internal/modules/integrations/ports"
)

var _ ports.OperationRunStore = (*OperationRunRepository)(nil)
var _ ports.OperationRunReadStore = (*OperationRunRepository)(nil)
var _ ports.OperationRunLatestReadStore = (*OperationRunRepository)(nil)

type OperationRunRepository struct {
	pool     *pgxpool.Pool
	tenantID string
}

func (r *OperationRunRepository) BeginExclusive(ctx context.Context, run domain.OperationRun) (domain.OperationRun, bool, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return domain.OperationRun{}, false, err
	}
	defer tx.Rollback(ctx)

	// hashtextextended produces one deterministic int8 key for the tenant-scoped
	// installation and operation type; the transaction lock releases on commit.
	lockKey := r.tenantID + "|" + run.InstallationID + "|" + run.OperationType
	if _, err = tx.Exec(ctx, `SELECT pg_advisory_xact_lock(hashtextextended($1, 0))`, lockKey); err != nil {
		return domain.OperationRun{}, false, err
	}
	row := tx.QueryRow(ctx, `
		SELECT operation_run_id, tenant_id, installation_id, operation_type, status,
			result_code, failure_code, translated_error_code, attempt_count, actor_type, actor_id,
			provider_evidence_json, duration_ms, started_at, completed_at, created_at, updated_at
		FROM integration_operation_runs
		WHERE tenant_id = $1 AND installation_id = $2 AND operation_type = $3
		  AND status IN ('queued', 'running')
		ORDER BY created_at DESC, operation_run_id DESC LIMIT 1
	`, r.tenantID, run.InstallationID, run.OperationType)
	active, _, scanErr := scanOperationRun(row)
	if scanErr == nil {
		if err = tx.Commit(ctx); err != nil {
			return domain.OperationRun{}, false, err
		}
		return active, true, nil
	}
	if scanErr != pgx.ErrNoRows {
		return domain.OperationRun{}, false, scanErr
	}
	_, err = tx.Exec(ctx, `
		INSERT INTO integration_operation_runs (
			tenant_id, operation_run_id, installation_id, operation_type, status,
			result_code, failure_code, translated_error_code, attempt_count, actor_type, actor_id,
			provider_evidence_json, duration_ms, started_at, completed_at, created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)
	`, r.tenantID, run.OperationRunID, run.InstallationID, run.OperationType, run.Status,
		run.ResultCode, run.FailureCode, run.TranslatedErrorCode, run.AttemptCount, run.ActorType, run.ActorID,
		marshalJSONObject(run.ProviderEvidence), run.DurationMs, timestamptzArg(run.StartedAt), timestamptzArg(run.CompletedAt), run.CreatedAt, run.UpdatedAt)
	if err != nil {
		return domain.OperationRun{}, false, err
	}
	if err = tx.Commit(ctx); err != nil {
		return domain.OperationRun{}, false, err
	}
	return run, false, nil
}

func NewOperationRunRepository(pool *pgxpool.Pool, tenantID string) *OperationRunRepository {
	return &OperationRunRepository{pool: pool, tenantID: tenantID}
}

func (r *OperationRunRepository) SaveOperationRun(ctx context.Context, run domain.OperationRun) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO integration_operation_runs (
			tenant_id, operation_run_id, installation_id, operation_type,
			status, result_code, failure_code, translated_error_code, attempt_count,
			actor_type, actor_id, provider_evidence_json, duration_ms, started_at, completed_at,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4,
			$5, $6, $7, $8, $9,
			$10, $11, $12, $13, $14, $15,
			$16, $17
		)
		ON CONFLICT (tenant_id, operation_run_id) DO UPDATE SET
			installation_id = EXCLUDED.installation_id,
			operation_type = EXCLUDED.operation_type,
			status = EXCLUDED.status,
			result_code = EXCLUDED.result_code,
			failure_code = EXCLUDED.failure_code,
			translated_error_code = EXCLUDED.translated_error_code,
			attempt_count = EXCLUDED.attempt_count,
			actor_type = EXCLUDED.actor_type,
			actor_id = EXCLUDED.actor_id,
			provider_evidence_json = EXCLUDED.provider_evidence_json,
			duration_ms = EXCLUDED.duration_ms,
			started_at = EXCLUDED.started_at,
			completed_at = EXCLUDED.completed_at,
			created_at = EXCLUDED.created_at,
			updated_at = EXCLUDED.updated_at
	`, r.tenantID, run.OperationRunID, run.InstallationID, run.OperationType, run.Status,
		run.ResultCode, run.FailureCode, run.TranslatedErrorCode, run.AttemptCount, run.ActorType, run.ActorID,
		marshalJSONObject(run.ProviderEvidence), run.DurationMs, timestamptzArg(run.StartedAt), timestamptzArg(run.CompletedAt), run.CreatedAt, run.UpdatedAt)
	return err
}

func (r *OperationRunRepository) ListByInstallation(ctx context.Context, installationID string) ([]domain.OperationRun, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT
			operation_run_id, tenant_id, installation_id, operation_type, status,
			result_code, failure_code, translated_error_code, attempt_count, actor_type, actor_id,
			provider_evidence_json, duration_ms,
			started_at, completed_at, created_at, updated_at
		FROM integration_operation_runs
		WHERE tenant_id = $1
		  AND installation_id = $2
		ORDER BY created_at DESC, operation_run_id DESC
	`, r.tenantID, installationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	runs := make([]domain.OperationRun, 0)
	for rows.Next() {
		run, _, err := scanOperationRun(rows)
		if err != nil {
			return nil, err
		}
		runs = append(runs, run)
	}
	return runs, rows.Err()
}

func (r *OperationRunRepository) LatestRunsByModule(ctx context.Context, installationID string) ([]ports.LatestRunByModule, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT operation_type,
			MAX(started_at) AS latest_attempted_at,
			MAX(completed_at) FILTER (WHERE status = 'succeeded') AS latest_successful_at
		FROM integration_operation_runs
		WHERE tenant_id = $1
		  AND installation_id = $2
		GROUP BY operation_type
		ORDER BY operation_type ASC
	`, r.tenantID, installationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	latestRuns := make([]ports.LatestRunByModule, 0)
	for rows.Next() {
		var item ports.LatestRunByModule
		var latestAttemptedAt pgtype.Timestamptz
		var latestSuccessfulAt pgtype.Timestamptz
		if err := rows.Scan(&item.OperationType, &latestAttemptedAt, &latestSuccessfulAt); err != nil {
			return nil, err
		}
		item.LatestAttemptedAt = scanTimestamptz(latestAttemptedAt)
		item.LatestSuccessfulAt = scanTimestamptz(latestSuccessfulAt)
		latestRuns = append(latestRuns, item)
	}
	return latestRuns, rows.Err()
}

func (r *OperationRunRepository) ListRuns(ctx context.Context, query ports.RunListQuery) (ports.RunPage, error) {
	query, err := query.Normalize()
	if err != nil {
		return ports.RunPage{}, err
	}
	rows, err := r.pool.Query(ctx, `
		SELECT operation_run_id, installation_id, operation_type, status,
			result_code, failure_code, translated_error_code, attempt_count, duration_ms,
			started_at, completed_at
		FROM integration_operation_runs
		WHERE tenant_id = $1
		  AND installation_id = $2
		  AND started_at >= now() - interval '90 days'
		  AND ($3 = '' OR operation_type = $3)
		  AND ($4 = '' OR status = $4)
		  AND ($5::timestamptz IS NULL OR (started_at, operation_run_id) < ($5, $6))
		ORDER BY started_at DESC, operation_run_id DESC
		LIMIT $7
	`, r.tenantID, query.InstallationID, query.Module, string(query.Status), nullableCursorTime(query.Cursor), query.Cursor.OperationRunID, query.Limit+1)
	if err != nil {
		return ports.RunPage{}, err
	}
	defer rows.Close()

	items := make([]ports.RunReadModel, 0, query.Limit+1)
	for rows.Next() {
		var item ports.RunReadModel
		var startedAt pgtype.Timestamptz
		var finishedAt pgtype.Timestamptz
		if err := rows.Scan(&item.OperationRunID, &item.InstallationID, &item.Module, &item.Status,
			&item.ResultCode, &item.FailureCode, &item.TranslatedErrorCode, &item.AttemptCount, &item.DurationMs,
			&startedAt, &finishedAt); err != nil {
			return ports.RunPage{}, err
		}
		item.StartedAt = scanTimestamptz(startedAt)
		item.FinishedAt = scanTimestamptz(finishedAt)
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return ports.RunPage{}, err
	}
	page := ports.RunPage{Items: items}
	if len(items) > query.Limit {
		page.Items = items[:query.Limit]
		last := page.Items[len(page.Items)-1]
		if last.StartedAt != nil {
			page.NextCursor = &ports.RunCursor{StartedAt: *last.StartedAt, OperationRunID: last.OperationRunID}
		}
	}
	return page, nil
}

func nullableCursorTime(cursor ports.RunCursor) any {
	if cursor.IsFirstPage() {
		return nil
	}
	return cursor.StartedAt
}

func scanOperationRun(scanner interface {
	Scan(dest ...any) error
}) (domain.OperationRun, bool, error) {
	var run domain.OperationRun
	var startedAt pgtype.Timestamptz
	var completedAt pgtype.Timestamptz
	var providerEvidenceRaw []byte

	err := scanner.Scan(
		&run.OperationRunID,
		&run.TenantID,
		&run.InstallationID,
		&run.OperationType,
		&run.Status,
		&run.ResultCode,
		&run.FailureCode,
		&run.TranslatedErrorCode,
		&run.AttemptCount,
		&run.ActorType,
		&run.ActorID,
		&providerEvidenceRaw,
		&run.DurationMs,
		&startedAt,
		&completedAt,
		&run.CreatedAt,
		&run.UpdatedAt,
	)
	if err != nil {
		return domain.OperationRun{}, false, err
	}

	run.StartedAt = scanTimestamptz(startedAt)
	run.CompletedAt = scanTimestamptz(completedAt)
	if len(providerEvidenceRaw) > 0 {
		_ = json.Unmarshal(providerEvidenceRaw, &run.ProviderEvidence)
	}

	return run, true, nil
}

func marshalJSONObject(value map[string]any) []byte {
	if len(value) == 0 {
		return []byte(`{}`)
	}
	raw, err := json.Marshal(value)
	if err != nil {
		return []byte(`{}`)
	}
	return raw
}
