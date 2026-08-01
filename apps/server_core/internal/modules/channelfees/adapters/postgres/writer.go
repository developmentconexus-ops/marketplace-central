// Package postgres is the Postgres-backed implementation of the channel_fees
// core ports (IC-01). It follows the upsert/keep-absent idiom proven at
// internal_read/adapters/mirror/writer.go: INSERT ... ON CONFLICT DO UPDATE
// on the table's natural key, current-value semantics (no history).
package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"marketplace-central/apps/server_core/internal/modules/channelfees/domain"
	"marketplace-central/apps/server_core/internal/modules/channelfees/ports"
)

// Writer is the Postgres-backed ports.ChannelFeeWriter.
type Writer struct {
	pool *pgxpool.Pool
}

// NewWriter builds a Writer over pool.
func NewWriter(pool *pgxpool.Pool) *Writer {
	return &Writer{pool: pool}
}

var _ ports.ChannelFeeWriter = (*Writer)(nil)

const upsertFeeSQL = `
INSERT INTO channel_fees (
	tenant_id, provider, installation_id, layer, fee_kind, subject_type, subject_id,
	value_type, value, currency, detail, origem, coletado_em, source_time, updated_at
) VALUES (
	$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, now()
)
ON CONFLICT ON CONSTRAINT channel_fees_natural_key_key DO UPDATE SET
	value_type   = EXCLUDED.value_type,
	value        = EXCLUDED.value,
	currency     = EXCLUDED.currency,
	detail       = EXCLUDED.detail,
	origem       = EXCLUDED.origem,
	coletado_em  = EXCLUDED.coletado_em,
	source_time  = EXCLUDED.source_time,
	updated_at   = now()`

// UpsertFee validates the write-time mandate (domain.ValidateFee — layer 3
// commission requires a decomposed detail) and then upserts on the natural
// key: a 2nd UpsertFee for the same (tenant, provider, installation, subject,
// layer, fee_kind) collapses to the same row, coletado_em/source_time
// advanced (IC-01 Persistence Expectations — current-value semantics, no
// history in this migration).
func (w *Writer) UpsertFee(ctx context.Context, fee domain.Fee) error {
	if err := domain.ValidateFee(fee); err != nil {
		return err
	}
	if fee.TenantID == "" {
		return errors.New("channelfees: empty tenantID")
	}

	_, err := w.pool.Exec(ctx, upsertFeeSQL,
		fee.TenantID, fee.Provider, fee.InstallationID, fee.Layer, fee.FeeKind,
		fee.SubjectType, fee.SubjectID, fee.ValueType, fee.Value, fee.Currency,
		nullableJSON(fee.Detail), fee.Origem, fee.ColetadoEm, fee.SourceTime,
	)
	if err != nil {
		return fmt.Errorf("channelfees: upsert fee: %w", err)
	}
	return nil
}

// nullableJSON returns nil (SQL NULL) for an empty/absent detail, and the
// raw bytes otherwise — pgx encodes a nil json.RawMessage as SQL NULL, never
// as the JSON literal "null" (same helper shape as
// mutations/adapters/postgres/apply_repository.go:240).
func nullableJSON(value json.RawMessage) any {
	if len(value) == 0 {
		return nil
	}
	return value
}
