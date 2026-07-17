package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"marketplace-central/apps/server_core/internal/modules/mutations/application"
	"marketplace-central/apps/server_core/internal/modules/mutations/domain"
	"marketplace-central/apps/server_core/internal/modules/mutations/ports"
)

func (r *Repository) CloneRetry(ctx context.Context, sourceID string, createdAt time.Time) (ports.Protocol, bool, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return ports.Protocol{}, false, fmt.Errorf("begin retry mutation protocol: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()
	var source ports.Protocol
	if err := tx.QueryRow(ctx, `SELECT protocol_id,installation_id,type,state,actor,intent,selection,totals,source_as_of,retried_from,created_at,previewed_at,approved_at,finished_at FROM mutation_protocols WHERE tenant_id=$1 AND protocol_id=$2 FOR UPDATE`, r.tenantID, sourceID).Scan(&source.ProtocolID, &source.InstallationID, &source.Type, &source.State, &source.Actor, &source.Intent, &source.Selection, &source.Totals, &source.SourceAsOf, &source.RetriedFrom, &source.CreatedAt, &source.PreviewedAt, &source.ApprovedAt, &source.FinishedAt); err != nil {
		return ports.Protocol{}, false, fmt.Errorf("lock retry source protocol: %w", err)
	}
	rows, err := tx.Query(ctx, `SELECT listing_id FROM mutation_items WHERE tenant_id=$1 AND protocol_id=$2 AND state='failed' AND failure->>'code' IN ('provider_rate_limited','provider_unavailable','stale_source') ORDER BY seq`, r.tenantID, sourceID)
	if err != nil {
		return ports.Protocol{}, false, fmt.Errorf("read retryable mutation failures: %w", err)
	}
	listingIDs := make([]string, 0)
	for rows.Next() {
		var listingID string
		if err := rows.Scan(&listingID); err != nil {
			rows.Close()
			return ports.Protocol{}, false, fmt.Errorf("scan retryable mutation failure: %w", err)
		}
		listingIDs = append(listingIDs, listingID)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return ports.Protocol{}, false, fmt.Errorf("iterate retryable mutation failures: %w", err)
	}
	rows.Close()
	if len(listingIDs) == 0 {
		return ports.Protocol{}, false, nil
	}
	selection, err := json.Marshal(struct {
		Mode       string   `json:"mode"`
		ListingIDs []string `json:"listing_ids"`
	}{Mode: "explicit", ListingIDs: listingIDs})
	if err != nil {
		return ports.Protocol{}, false, fmt.Errorf("marshal retry mutation selection: %w", err)
	}
	id, err := r.allocateProtocolID(ctx, tx)
	if err != nil {
		return ports.Protocol{}, false, err
	}
	totals := []byte(`{"items":0,"previewed":0,"applied":0,"failed":0,"skipped":0}`)
	createdAt = createdAt.UTC()
	if _, err := tx.Exec(ctx, `INSERT INTO mutation_protocols(tenant_id,protocol_id,installation_id,type,state,actor,intent,selection,totals,source_as_of,retried_from,created_at) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,NULL,$10,$11)`, r.tenantID, id, source.InstallationID, source.Type, domain.ProtocolStateDraft, source.Actor, source.Intent, selection, totals, sourceID, createdAt); err != nil {
		return ports.Protocol{}, false, fmt.Errorf("insert retry protocol: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return ports.Protocol{}, false, fmt.Errorf("commit retry protocol: %w", err)
	}
	return ports.Protocol{ProtocolID: id, InstallationID: source.InstallationID, Type: source.Type, State: domain.ProtocolStateDraft, Actor: source.Actor, Intent: append([]byte(nil), source.Intent...), Selection: selection, Totals: totals, RetriedFrom: &sourceID, CreatedAt: createdAt}, true, nil
}

func (r *Repository) ListProtocols(ctx context.Context, q application.ProtocolQuery) (application.ProtocolPage, error) {
	args := []any{r.tenantID, q.InstallationID}
	where := []string{"tenant_id=$1", "installation_id=$2"}
	if q.State != "" {
		args = append(args, q.State)
		where = append(where, fmt.Sprintf("state=$%d", len(args)))
	}
	if q.Type != "" {
		args = append(args, q.Type)
		where = append(where, fmt.Sprintf("type=$%d", len(args)))
	}
	if !q.Cursor.CreatedAt.IsZero() {
		args = append(args, q.Cursor.CreatedAt.UTC(), q.Cursor.ProtocolID)
		where = append(where, fmt.Sprintf("(created_at,protocol_id)<($%d,$%d)", len(args)-1, len(args)))
	}
	args = append(args, q.Limit+1)
	rows, err := r.pool.Query(ctx, fmt.Sprintf(`SELECT protocol_id,installation_id,type,state,actor,intent,selection,totals,source_as_of,retried_from,created_at,previewed_at,approved_at,finished_at FROM mutation_protocols WHERE %s ORDER BY created_at DESC,protocol_id DESC LIMIT $%d`, strings.Join(where, " AND "), len(args)), args...)
	if err != nil {
		return application.ProtocolPage{}, fmt.Errorf("list mutation protocols: %w", err)
	}
	defer rows.Close()
	items := make([]ports.Protocol, 0, q.Limit+1)
	for rows.Next() {
		item, err := scanProtocol(rows)
		if err != nil {
			return application.ProtocolPage{}, fmt.Errorf("scan mutation protocol: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return application.ProtocolPage{}, fmt.Errorf("iterate mutation protocols: %w", err)
	}
	page := application.ProtocolPage{Items: items}
	if len(items) > q.Limit {
		page.Items = items[:q.Limit]
		last := page.Items[len(page.Items)-1]
		page.NextCursor = &application.ProtocolCursor{CreatedAt: last.CreatedAt, ProtocolID: last.ProtocolID}
	}
	return page, nil
}

func (r *Repository) ListItems(ctx context.Context, q application.ItemQuery) (application.ItemPage, error) {
	args := []any{r.tenantID, q.ProtocolID, q.Cursor.Seq}
	where := []string{"tenant_id=$1", "protocol_id=$2", "seq>$3"}
	if q.State != "" {
		args = append(args, q.State)
		where = append(where, fmt.Sprintf("state=$%d", len(args)))
	}
	args = append(args, q.Limit+1)
	rows, err := r.pool.Query(ctx, fmt.Sprintf(`SELECT seq,item_id,listing_id,idempotency_key,before,after,state,failure,applied_at FROM mutation_items WHERE %s ORDER BY seq ASC LIMIT $%d`, strings.Join(where, " AND "), len(args)), args...)
	if err != nil {
		return application.ItemPage{}, fmt.Errorf("list mutation items: %w", err)
	}
	defer rows.Close()
	items := make([]ports.MutationItem, 0, q.Limit+1)
	for rows.Next() {
		var item ports.MutationItem
		if err := rows.Scan(&item.Seq, &item.ItemID, &item.ListingID, &item.IdempotencyKey, &item.Before, &item.After, &item.State, &item.Failure, &item.AppliedAt); err != nil {
			return application.ItemPage{}, fmt.Errorf("scan mutation item: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return application.ItemPage{}, fmt.Errorf("iterate mutation items: %w", err)
	}
	page := application.ItemPage{Items: items}
	if len(items) > q.Limit {
		page.Items = items[:q.Limit]
		page.NextCursor = &application.ItemCursor{Seq: page.Items[len(page.Items)-1].Seq}
	}
	return page, nil
}
