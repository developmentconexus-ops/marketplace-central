package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	ordersdomain "marketplace-central/apps/server_core/internal/modules/orders/domain"
	"marketplace-central/apps/server_core/internal/modules/orders/ports"
)

var _ ports.SankhyaLinkageRepository = (*SankhyaLinkageRepository)(nil)

type SankhyaLinkageRepository struct {
	pool     *pgxpool.Pool
	tenantID string
}

func NewSankhyaLinkageRepository(pool *pgxpool.Pool, tenantID string) *SankhyaLinkageRepository {
	return &SankhyaLinkageRepository{pool: pool, tenantID: strings.TrimSpace(tenantID)}
}

func (r *SankhyaLinkageRepository) LoadCurrent(ctx context.Context, scope ordersdomain.LinkageScope) (ordersdomain.SankhyaLinkage, bool, error) {
	if err := r.validateScope(scope); err != nil {
		return ordersdomain.SankhyaLinkage{}, false, err
	}
	if r.pool == nil {
		return ordersdomain.SankhyaLinkage{}, false, nil
	}
	return loadLinkageByProviderOrder(ctx, r.pool, scope)
}

func (r *SankhyaLinkageRepository) AppendConfirmation(ctx context.Context, linkage ordersdomain.SankhyaLinkage) (ordersdomain.SankhyaLinkage, error) {
	if err := linkage.Validate(); err != nil {
		return ordersdomain.SankhyaLinkage{}, err
	}
	if err := r.validateScope(linkage.Scope); err != nil {
		return ordersdomain.SankhyaLinkage{}, err
	}
	if r.pool == nil {
		return ordersdomain.SankhyaLinkage{}, fmt.Errorf("%w: repository unavailable", ordersdomain.ErrSankhyaLinkageNotFound)
	}

	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return ordersdomain.SankhyaLinkage{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	tag, err := tx.Exec(ctx, `
		INSERT INTO orders_sankhya_linkage_events (
			tenant_id, installation_id, event_id, event_type, provider_order_id, external_order_key,
			internal_document_id, actor_type, actor_id, reason, idempotency_key,
			source_at, recorded_at, previous_event_id, configuration_revision,
			evidence_state, evidence_reference
		) VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8, $9, $10, $11,
			$12, $13, $14, $15,
			$16, $17
		)
		ON CONFLICT DO NOTHING
	`, linkage.Scope.TenantID, linkage.Scope.InstallationID, linkage.Audit.EventID, linkage.Audit.EventType, linkage.Scope.ProviderOrderID, linkage.ExternalOrderKey,
		linkage.Header.DocumentID, linkage.Audit.ActorType, linkage.Audit.ActorID, linkage.Audit.Reason, linkage.Audit.IdempotencyKey,
		linkage.Audit.SourceAt.UTC(), linkage.Audit.RecordedAt.UTC(), nullableText(linkage.Audit.PreviousEventID), linkage.Audit.ConfigurationRevision,
		linkage.Audit.EvidenceState, linkage.Audit.EvidenceReference)
	if err != nil {
		return ordersdomain.SankhyaLinkage{}, classifyLinkagePersistenceError(err)
	}
	if tag.RowsAffected() == 0 {
		stored, found, err := loadLinkageByIdempotencyKey(ctx, tx, linkage.Scope.TenantID, linkage.Scope.InstallationID, linkage.Audit.IdempotencyKey)
		if err != nil {
			return ordersdomain.SankhyaLinkage{}, err
		}
		if !found {
			stored, found, err = loadLinkageByProviderOrder(ctx, tx, linkage.Scope)
			if err != nil {
				return ordersdomain.SankhyaLinkage{}, err
			}
		}
		if found && stored.SemanticallyEqual(linkage) {
			if err := tx.Commit(ctx); err != nil {
				return ordersdomain.SankhyaLinkage{}, err
			}
			return stored, nil
		}
		return ordersdomain.SankhyaLinkage{}, ordersdomain.ErrSankhyaLinkageConflict
	}

	for _, line := range linkage.Lines {
		var state ordersdomain.LineReconciliationState
		err := tx.QueryRow(ctx, `
			SELECT reconciliation_state
			FROM orders_marketplace_order_items
			WHERE tenant_id = $1 AND installation_id = $2
			  AND provider_order_id = $3 AND mpc_line_id = $4
			FOR KEY SHARE
		`, linkage.Scope.TenantID, linkage.Scope.InstallationID, linkage.Scope.ProviderOrderID, line.MPCLineID).Scan(&state)
		if errors.Is(err, pgx.ErrNoRows) {
			return ordersdomain.SankhyaLinkage{}, ordersdomain.ErrSankhyaLinkageNotFound
		}
		if err != nil {
			return ordersdomain.SankhyaLinkage{}, err
		}
		if state != ordersdomain.LineReconciliationStable {
			return ordersdomain.SankhyaLinkage{}, ordersdomain.ErrSankhyaLinkageConflict
		}

		if _, err := tx.Exec(ctx, `
			INSERT INTO orders_sankhya_linkage_lines (
				tenant_id, installation_id, event_id, provider_order_id, mpc_line_id,
				internal_document_id, internal_line_number
			) VALUES ($1, $2, $3, $4, $5, $6, $7)
		`, linkage.Scope.TenantID, linkage.Scope.InstallationID, linkage.Audit.EventID, linkage.Scope.ProviderOrderID, line.MPCLineID,
			line.Origin.DocumentID, line.Origin.LineNumber); err != nil {
			return ordersdomain.SankhyaLinkage{}, classifyLinkagePersistenceError(err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return ordersdomain.SankhyaLinkage{}, classifyLinkagePersistenceError(err)
	}
	return linkage, nil
}

func (r *SankhyaLinkageRepository) validateScope(scope ordersdomain.LinkageScope) error {
	if strings.TrimSpace(scope.TenantID) == "" || scope.TenantID != r.tenantID || strings.TrimSpace(scope.InstallationID) == "" || strings.TrimSpace(scope.ProviderOrderID) == "" {
		return fmt.Errorf("%w: linkage scope", ordersdomain.ErrInvalidSankhyaLinkage)
	}
	return nil
}

type linkageQuerier interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

func loadLinkageByProviderOrder(ctx context.Context, querier linkageQuerier, scope ordersdomain.LinkageScope) (ordersdomain.SankhyaLinkage, bool, error) {
	return loadLinkage(ctx, querier, `
		WHERE tenant_id = $1 AND installation_id = $2 AND provider_order_id = $3
	`, scope.TenantID, scope.InstallationID, scope.ProviderOrderID)
}

func loadLinkageByIdempotencyKey(ctx context.Context, querier linkageQuerier, tenantID, installationID, idempotencyKey string) (ordersdomain.SankhyaLinkage, bool, error) {
	return loadLinkage(ctx, querier, `
		WHERE tenant_id = $1 AND installation_id = $2 AND idempotency_key = $3
	`, tenantID, installationID, idempotencyKey)
}

func loadLinkage(ctx context.Context, querier linkageQuerier, predicate string, args ...any) (ordersdomain.SankhyaLinkage, bool, error) {
	var linkage ordersdomain.SankhyaLinkage
	var sourceAt time.Time
	var recordedAt time.Time
	var previousEventID pgtype.Text
	err := querier.QueryRow(ctx, `
		SELECT tenant_id, installation_id, provider_order_id, external_order_key, internal_document_id,
		       event_id, event_type, actor_type, actor_id, reason, idempotency_key,
		       source_at, recorded_at, previous_event_id, configuration_revision,
		       evidence_state, evidence_reference
		FROM orders_sankhya_linkage_events
	`+predicate, args...).Scan(
		&linkage.Scope.TenantID, &linkage.Scope.InstallationID, &linkage.Scope.ProviderOrderID, &linkage.ExternalOrderKey, &linkage.Header.DocumentID,
		&linkage.Audit.EventID, &linkage.Audit.EventType, &linkage.Audit.ActorType, &linkage.Audit.ActorID, &linkage.Audit.Reason, &linkage.Audit.IdempotencyKey,
		&sourceAt, &recordedAt, &previousEventID, &linkage.Audit.ConfigurationRevision,
		&linkage.Audit.EvidenceState, &linkage.Audit.EvidenceReference,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return ordersdomain.SankhyaLinkage{}, false, nil
	}
	if err != nil {
		return ordersdomain.SankhyaLinkage{}, false, err
	}
	linkage.Audit.SourceAt = sourceAt.UTC()
	linkage.Audit.RecordedAt = recordedAt.UTC()
	if previousEventID.Valid {
		linkage.Audit.PreviousEventID = previousEventID.String
	}

	rows, err := querier.Query(ctx, `
		SELECT mpc_line_id, internal_document_id, internal_line_number
		FROM orders_sankhya_linkage_lines
		WHERE tenant_id = $1 AND installation_id = $2 AND event_id = $3
		ORDER BY mpc_line_id ASC
	`, linkage.Scope.TenantID, linkage.Scope.InstallationID, linkage.Audit.EventID)
	if err != nil {
		return ordersdomain.SankhyaLinkage{}, false, err
	}
	defer rows.Close()
	for rows.Next() {
		var line ordersdomain.SankhyaLineMapping
		if err := rows.Scan(&line.MPCLineID, &line.Origin.DocumentID, &line.Origin.LineNumber); err != nil {
			return ordersdomain.SankhyaLinkage{}, false, err
		}
		linkage.Lines = append(linkage.Lines, line)
	}
	if err := rows.Err(); err != nil {
		return ordersdomain.SankhyaLinkage{}, false, err
	}
	return linkage, true, nil
}

func classifyLinkagePersistenceError(err error) error {
	var postgresError *pgconn.PgError
	if !errors.As(err, &postgresError) {
		return err
	}
	switch postgresError.Code {
	case "23503":
		return ordersdomain.ErrSankhyaLinkageNotFound
	case "23505", "23514":
		return ordersdomain.ErrSankhyaLinkageConflict
	default:
		return err
	}
}

func nullableText(value string) pgtype.Text {
	if value == "" {
		return pgtype.Text{}
	}
	return pgtype.Text{String: value, Valid: true}
}
