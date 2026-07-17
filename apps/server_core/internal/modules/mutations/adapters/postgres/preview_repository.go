package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"marketplace-central/apps/server_core/internal/modules/mutations/domain"
	"marketplace-central/apps/server_core/internal/modules/mutations/ports"
)

func (r *Repository) ReplacePreview(ctx context.Context, protocolID string, inputs []ports.ReplaceItemInput, sourceAsOf, previewedAt time.Time) ([]ports.MutationItem, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, fmt.Errorf("begin replace mutation preview: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var state domain.ProtocolState
	if err := tx.QueryRow(ctx, `SELECT state FROM mutation_protocols WHERE tenant_id=$1 AND protocol_id=$2 FOR UPDATE`, r.tenantID, protocolID).Scan(&state); err != nil {
		return nil, fmt.Errorf("lock mutation protocol for preview: %w", err)
	}
	if err := domain.TransitionProtocolState(state, domain.ProtocolStatePreviewed); err != nil {
		return nil, err
	}
	if _, err := tx.Exec(ctx, `DELETE FROM mutation_items WHERE tenant_id=$1 AND protocol_id=$2`, r.tenantID, protocolID); err != nil {
		return nil, fmt.Errorf("delete old mutation preview: %w", err)
	}

	items := make([]ports.MutationItem, 0, len(inputs))
	for i, input := range inputs {
		seq := i + 1
		item := ports.MutationItem{
			Seq: seq, ItemID: fmt.Sprintf("%s-%03d", protocolID, seq), ListingID: input.ListingID,
			IdempotencyKey: protocolID + ":" + input.ListingID, Before: input.Before, After: input.After,
			State: domain.ItemStatePreviewed,
		}
		if len(item.After) == 0 {
			return nil, fmt.Errorf("mutation item %d after is required", seq)
		}
		if _, err := tx.Exec(ctx, `INSERT INTO mutation_items (tenant_id,protocol_id,seq,item_id,listing_id,idempotency_key,before,after,state) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`, r.tenantID, protocolID, item.Seq, item.ItemID, item.ListingID, item.IdempotencyKey, nullableJSON(item.Before), item.After, item.State); err != nil {
			return nil, fmt.Errorf("insert preview item %d: %w", seq, err)
		}
		items = append(items, item)
	}
	if _, err := tx.Exec(ctx, `UPDATE mutation_protocols SET state='previewed',source_as_of=$3,previewed_at=$4,totals=jsonb_build_object('items',$5::int,'previewed',$5::int,'applied',0,'failed',0,'skipped',0) WHERE tenant_id=$1 AND protocol_id=$2`, r.tenantID, protocolID, sourceAsOf.UTC(), previewedAt.UTC(), len(items)); err != nil {
		return nil, fmt.Errorf("update mutation preview: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit mutation preview: %w", err)
	}
	return items, nil
}
