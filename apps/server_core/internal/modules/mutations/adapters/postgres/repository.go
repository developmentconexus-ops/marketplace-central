package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"marketplace-central/apps/server_core/internal/modules/mutations/domain"
	"marketplace-central/apps/server_core/internal/modules/mutations/ports"
)

var _ ports.ProtocolRepository = (*Repository)(nil)

type Repository struct {
	pool     *pgxpool.Pool
	tenantID string
}

type protocolScanner interface {
	Scan(...any) error
}

func NewRepository(pool *pgxpool.Pool, tenantID string) *Repository {
	return &Repository{pool: pool, tenantID: tenantID}
}

func (r *Repository) CreateProtocol(ctx context.Context, input ports.CreateProtocolInput) (ports.Protocol, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return ports.Protocol{}, fmt.Errorf("begin create mutation protocol: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	if _, err := tx.Exec(ctx, `SELECT pg_advisory_xact_lock(hashtextextended('mutation_protocols:' || $1, 0))`, r.tenantID); err != nil {
		return ports.Protocol{}, fmt.Errorf("lock mutation protocol tenant: %w", err)
	}
	var next int
	if err := tx.QueryRow(ctx, `
		SELECT COALESCE(MAX(substring(protocol_id FROM 4)::integer), 0) + 1
		FROM mutation_protocols
		WHERE tenant_id = $1`, r.tenantID).Scan(&next); err != nil {
		return ports.Protocol{}, fmt.Errorf("allocate mutation protocol id: %w", err)
	}
	if next > 999999 {
		return ports.Protocol{}, fmt.Errorf("mutation protocol id space exhausted for tenant")
	}
	protocolID := fmt.Sprintf("MP-%06d", next)
	totals := []byte(`{"items":0,"previewed":0,"applied":0,"failed":0,"skipped":0}`)
	createdAt := input.CreatedAt.UTC()

	_, err = tx.Exec(ctx, `
		INSERT INTO mutation_protocols
			(tenant_id, protocol_id, installation_id, type, state, actor, intent, selection, totals, source_as_of, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NULL, $10)`,
		r.tenantID, protocolID, input.InstallationID, input.Type, domain.ProtocolStateDraft,
		input.Actor, input.Intent, input.Selection, totals, createdAt)
	if err != nil {
		return ports.Protocol{}, fmt.Errorf("insert mutation protocol: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return ports.Protocol{}, fmt.Errorf("commit mutation protocol: %w", err)
	}
	return ports.Protocol{
		ProtocolID: protocolID, InstallationID: input.InstallationID, Type: input.Type,
		State: domain.ProtocolStateDraft, Actor: input.Actor, Intent: input.Intent,
		Selection: input.Selection, Totals: totals, CreatedAt: createdAt,
	}, nil
}

func (r *Repository) GetProtocol(ctx context.Context, protocolID string) (ports.Protocol, bool, error) {
	protocol, err := scanProtocol(r.pool.QueryRow(ctx, `
		SELECT protocol_id, installation_id, type, state, actor, intent, selection, totals,
			source_as_of, retried_from, created_at, previewed_at, approved_at, finished_at
		FROM mutation_protocols
		WHERE tenant_id = $1 AND protocol_id = $2`, r.tenantID, protocolID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ports.Protocol{}, false, nil
		}
		return ports.Protocol{}, false, fmt.Errorf("get mutation protocol: %w", err)
	}
	return protocol, true, nil
}

func scanProtocol(row protocolScanner) (ports.Protocol, error) {
	var protocol ports.Protocol
	if err := row.Scan(
		&protocol.ProtocolID, &protocol.InstallationID, &protocol.Type, &protocol.State,
		&protocol.Actor, &protocol.Intent, &protocol.Selection, &protocol.Totals,
		&protocol.SourceAsOf, &protocol.RetriedFrom, &protocol.CreatedAt,
		&protocol.PreviewedAt, &protocol.ApprovedAt, &protocol.FinishedAt,
	); err != nil {
		return ports.Protocol{}, err
	}
	return protocol, nil
}
