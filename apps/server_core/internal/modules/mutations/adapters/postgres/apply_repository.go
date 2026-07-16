package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"marketplace-central/apps/server_core/internal/modules/mutations/domain"
	"marketplace-central/apps/server_core/internal/modules/mutations/ports"
)

var ErrImmutableItem = errors.New("mutation item outcome is immutable")

func (r *Repository) ReplaceItems(ctx context.Context, protocolID string, inputs []ports.ReplaceItemInput) ([]ports.MutationItem, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin replace mutation items: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()
	var state domain.ProtocolState
	if err := tx.QueryRow(ctx, `SELECT state FROM mutation_protocols WHERE tenant_id=$1 AND protocol_id=$2 FOR UPDATE`, r.tenantID, protocolID).Scan(&state); err != nil {
		return nil, fmt.Errorf("lock mutation protocol for item replacement: %w", err)
	}
	if state != domain.ProtocolStateDraft && state != domain.ProtocolStatePreviewed {
		return nil, fmt.Errorf("replace mutation items: %w", &domain.InvalidStateTransitionError{Code: domain.InvalidStateCode, From: state, To: domain.ProtocolStatePreviewed})
	}
	if _, err := tx.Exec(ctx, `DELETE FROM mutation_items WHERE tenant_id=$1 AND protocol_id=$2`, r.tenantID, protocolID); err != nil {
		return nil, fmt.Errorf("delete mutation items: %w", err)
	}
	items := make([]ports.MutationItem, 0, len(inputs))
	for i, input := range inputs {
		if len(input.After) == 0 {
			return nil, fmt.Errorf("mutation item %d after is required", i+1)
		}
		seq := i + 1
		itemID := fmt.Sprintf("%s-%03d", protocolID, seq)
		key := protocolID + ":" + input.ListingID
		_, err := tx.Exec(ctx, `INSERT INTO mutation_items (tenant_id,protocol_id,seq,item_id,listing_id,idempotency_key,before,after,state) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`, r.tenantID, protocolID, seq, itemID, input.ListingID, key, nullableJSON(input.Before), input.After, domain.ItemStatePreviewed)
		if err != nil {
			return nil, fmt.Errorf("insert mutation item %d: %w", seq, err)
		}
		items = append(items, ports.MutationItem{Seq: seq, ItemID: itemID, ListingID: input.ListingID, IdempotencyKey: key, Before: input.Before, After: input.After, State: domain.ItemStatePreviewed})
	}
	if _, err := tx.Exec(ctx, `UPDATE mutation_protocols SET state='previewed', previewed_at=COALESCE(previewed_at,now()), totals=jsonb_build_object('items',$3::int,'previewed',$3::int,'applied',0,'failed',0,'skipped',0) WHERE tenant_id=$1 AND protocol_id=$2`, r.tenantID, protocolID, len(items)); err != nil {
		return nil, fmt.Errorf("mark mutation protocol previewed: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit mutation item replacement: %w", err)
	}
	return items, nil
}

func (r *Repository) ApproveItems(ctx context.Context, protocolID string, approvedAt time.Time) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	var state domain.ProtocolState
	if err := tx.QueryRow(ctx, `SELECT state FROM mutation_protocols WHERE tenant_id=$1 AND protocol_id=$2 FOR UPDATE`, r.tenantID, protocolID).Scan(&state); err != nil {
		return fmt.Errorf("lock mutation protocol for approval: %w", err)
	}
	if err := domain.TransitionProtocolState(state, domain.ProtocolStateApproved); err != nil {
		return err
	}
	if _, err = tx.Exec(ctx, `UPDATE mutation_items SET state='approved' WHERE tenant_id=$1 AND protocol_id=$2 AND state='previewed'`, r.tenantID, protocolID); err != nil {
		return err
	}
	if _, err = tx.Exec(ctx, `UPDATE mutation_protocols SET state='approved',approved_at=$3 WHERE tenant_id=$1 AND protocol_id=$2`, r.tenantID, protocolID, approvedAt.UTC()); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

type protocolClaim struct {
	tx       pgx.Tx
	pool     *pgxpool.Pool
	tenantID string
	protocol ports.Protocol
	done     bool
}

func (c *protocolClaim) Protocol() ports.Protocol { return c.protocol }

func (r *Repository) ClaimProtocol(ctx context.Context, installationID string) (ports.ProtocolClaim, bool, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, false, err
	}
	c := &protocolClaim{tx: tx, pool: r.pool, tenantID: r.tenantID}
	err = scanProtocolInto(&c.protocol, tx.QueryRow(ctx, `SELECT protocol_id,installation_id,type,state,actor,intent,selection,totals,source_as_of,retried_from,created_at,previewed_at,approved_at,finished_at FROM mutation_protocols WHERE tenant_id=$1 AND installation_id=$2 AND state IN ('approved','applying') ORDER BY created_at FOR UPDATE SKIP LOCKED LIMIT 1`, r.tenantID, installationID))
	if errors.Is(err, pgx.ErrNoRows) {
		_ = tx.Rollback(ctx)
		return nil, false, nil
	}
	if err != nil {
		_ = tx.Rollback(ctx)
		return nil, false, err
	}
	if c.protocol.State == domain.ProtocolStateApproved {
		if err := domain.TransitionProtocolState(c.protocol.State, domain.ProtocolStateApplying); err != nil {
			_ = tx.Rollback(ctx)
			return nil, false, err
		}
		if _, err = tx.Exec(ctx, `UPDATE mutation_protocols SET state='applying' WHERE tenant_id=$1 AND protocol_id=$2`, r.tenantID, c.protocol.ProtocolID); err != nil {
			_ = tx.Rollback(ctx)
			return nil, false, err
		}
		c.protocol.State = domain.ProtocolStateApplying
	}
	return c, true, nil
}

func (c *protocolClaim) FetchPendingItems(ctx context.Context) ([]ports.MutationItem, error) {
	rows, err := c.tx.Query(ctx, `SELECT seq,item_id,listing_id,idempotency_key,before,after,state,failure,applied_at FROM mutation_items WHERE tenant_id=$1 AND protocol_id=$2 AND state NOT IN ('applied','failed','skipped') ORDER BY seq LIMIT 20`, c.tenantID, c.protocol.ProtocolID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []ports.MutationItem
	for rows.Next() {
		var i ports.MutationItem
		if err := rows.Scan(&i.Seq, &i.ItemID, &i.ListingID, &i.IdempotencyKey, &i.Before, &i.After, &i.State, &i.Failure, &i.AppliedAt); err != nil {
			return nil, err
		}
		out = append(out, i)
	}
	return out, rows.Err()
}
func (c *protocolClaim) WriteItemOutcome(ctx context.Context, itemID string, outcome ports.ItemOutcome) error {
	if outcome.State != domain.ItemStateApplied && outcome.State != domain.ItemStateFailed && outcome.State != domain.ItemStateSkipped {
		return domain.ErrItemOutcomeInvalid
	}
	var appliedAt any
	if outcome.State == domain.ItemStateApplied {
		appliedAt = outcome.AppliedAt.UTC()
	}
	tag, err := c.pool.Exec(ctx, `UPDATE mutation_items SET state=$4,failure=$5,applied_at=$6 WHERE tenant_id=$1 AND protocol_id=$2 AND item_id=$3 AND state NOT IN ('applied','failed','skipped')`, c.tenantID, c.protocol.ProtocolID, itemID, outcome.State, nullableJSON(outcome.Failure), appliedAt)
	if err != nil {
		return err
	}
	if tag.RowsAffected() != 1 {
		return ErrImmutableItem
	}
	return nil
}
func (c *protocolClaim) AppliedIdempotencyKeys(ctx context.Context, keys []string) ([]string, error) {
	rows, err := c.pool.Query(ctx, `SELECT idempotency_key FROM mutation_items WHERE tenant_id=$1 AND protocol_id=$2 AND idempotency_key=ANY($3) AND state='applied' ORDER BY seq`, c.tenantID, c.protocol.ProtocolID, keys)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []string
	for rows.Next() {
		var k string
		if err := rows.Scan(&k); err != nil {
			return nil, err
		}
		out = append(out, k)
	}
	return out, rows.Err()
}
func (c *protocolClaim) ItemStateCounts(ctx context.Context) (map[domain.ItemState]int, error) {
	rows, err := c.pool.Query(ctx, `SELECT state,count(*) FROM mutation_items WHERE tenant_id=$1 AND protocol_id=$2 GROUP BY state`, c.tenantID, c.protocol.ProtocolID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	counts := make(map[domain.ItemState]int)
	for rows.Next() {
		var state domain.ItemState
		var count int
		if err := rows.Scan(&state, &count); err != nil {
			return nil, err
		}
		counts[state] = count
	}
	return counts, rows.Err()
}
func (c *protocolClaim) Finish(ctx context.Context, state domain.ProtocolState, finishedAt time.Time) error {
	if err := domain.TransitionProtocolState(c.protocol.State, state); err != nil {
		return err
	}
	tag, err := c.tx.Exec(ctx, `UPDATE mutation_protocols SET state=$3,finished_at=$4 WHERE tenant_id=$1 AND protocol_id=$2 AND state='applying'`, c.tenantID, c.protocol.ProtocolID, state, finishedAt.UTC())
	if err != nil {
		return err
	}
	if tag.RowsAffected() != 1 {
		return domain.ErrInvalidStateTransition
	}
	c.protocol.State = state
	c.protocol.FinishedAt = &finishedAt
	return nil
}
func (c *protocolClaim) Commit(ctx context.Context) error {
	if c.done {
		return nil
	}
	c.done = true
	return c.tx.Commit(ctx)
}
func (c *protocolClaim) Rollback(ctx context.Context) error {
	if c.done {
		return nil
	}
	c.done = true
	return c.tx.Rollback(ctx)
}
func scanProtocolInto(p *ports.Protocol, row protocolScanner) error {
	return row.Scan(&p.ProtocolID, &p.InstallationID, &p.Type, &p.State, &p.Actor, &p.Intent, &p.Selection, &p.Totals, &p.SourceAsOf, &p.RetriedFrom, &p.CreatedAt, &p.PreviewedAt, &p.ApprovedAt, &p.FinishedAt)
}
func nullableJSON(value json.RawMessage) any {
	if len(value) == 0 {
		return nil
	}
	return value
}
