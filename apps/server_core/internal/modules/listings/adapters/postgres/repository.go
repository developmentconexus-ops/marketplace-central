package postgres

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"marketplace-central/apps/server_core/internal/modules/listings/domain"
	"marketplace-central/apps/server_core/internal/modules/listings/ports"
)

var _ ports.CompletedPullStore = (*Repository)(nil)
var _ ports.ListingReadRepository = (*Repository)(nil)

type Repository struct {
	pool     *pgxpool.Pool
	tenantID string
}

func NewRepository(pool *pgxpool.Pool, tenantID string) *Repository {
	return &Repository{pool: pool, tenantID: tenantID}
}

func (r *Repository) ListListingRows(ctx context.Context, q ports.ListingQuery) (ports.ListingRowPage, error) {
	args := []any{r.tenantID, q.InstallationID}
	where := []string{"l.tenant_id = $1", "l.installation_id = $2"}
	add := func(clause string, value any) {
		args = append(args, value)
		where = append(where, fmt.Sprintf(clause, len(args)))
	}
	if q.Filter.Status != "" {
		add("l.status = $%d", q.Filter.Status)
	}
	if q.Filter.SyncState != "" {
		add("l.sync_state = $%d", q.Filter.SyncState)
	}
	linkState := "COALESCE(NULLIF(pl.state, 'none'), 'unresolved')"
	if q.Filter.LinkState != "" {
		add(linkState+" = $%d", q.Filter.LinkState)
	}
	if q.Filter.ListingTypeCode != "" {
		add("l.listing_type_code = $%d", q.Filter.ListingTypeCode)
	}
	if q.Filter.ProductID != "" {
		add("pl.state = 'resolved' AND pl.internal_product_id = $%d::bigint", q.Filter.ProductID)
	}
	switch q.Filter.Exception {
	case domain.ListingExceptionSyncError:
		where = append(where, "(l.sync_state = 'error' OR l.sync_error IS NOT NULL)")
	case domain.ListingExceptionStale:
		where = append(where, "l.sync_state = 'stale'")
	case domain.ListingExceptionUnlinked:
		where = append(where, linkState+" = 'unresolved'")
	}
	if q.Q != "" {
		escaped := strings.NewReplacer("\\", "\\\\", "%", "\\%", "_", "\\_").Replace(q.Q)
		args = append(args, escaped)
		n := len(args)
		where = append(where, fmt.Sprintf("(l.title ILIKE '%%' || $%d || '%%' ESCAPE E'\\\\' OR l.provider_listing_id ILIKE '%%' || $%d || '%%' ESCAPE E'\\\\' OR pls.seller_sku ILIKE '%%' || $%d || '%%' ESCAPE E'\\\\')", n, n, n))
	}
	if !q.Cursor.IsFirstPage() {
		id, err := domain.ParseListingID(q.Cursor.ListingID)
		if err != nil {
			return ports.ListingRowPage{}, err
		}
		args = append(args, q.Cursor.LastTitle, id.ProviderListingID, id.VariationID)
		n := len(args)
		where = append(where, fmt.Sprintf("(l.title,l.provider_listing_id,l.variation_id) > ($%d,$%d,$%d)", n-2, n-1, n))
	}
	args = append(args, q.Limit+1)
	limitArg := len(args)
	sql := fmt.Sprintf(`SELECT l.installation_id,l.provider,l.provider_listing_id,l.variation_id,l.title,l.listing_type_code,l.status,
		l.price_amount::text,l.price_currency,l.published_quantity,l.sync_state,l.sync_error,l.quality_score,l.sales_30d,l.fetched_at,
		%s,CASE WHEN pl.state='resolved' THEN pl.internal_product_id::text END,NULLIF(pls.seller_sku,'')
	FROM listings l
	LEFT JOIN product_links pl ON pl.tenant_id=l.tenant_id AND pl.installation_id=l.installation_id AND pl.provider_item_id=l.provider_listing_id AND pl.provider_variation_id=CASE WHEN l.variation_id='-' THEN '' ELSE l.variation_id END
	LEFT JOIN product_link_listing_snapshots pls ON pls.tenant_id=l.tenant_id AND pls.installation_id=l.installation_id AND pls.provider_item_id=l.provider_listing_id AND pls.provider_variation_id=CASE WHEN l.variation_id='-' THEN '' ELSE l.variation_id END
	WHERE %s ORDER BY l.title,l.provider_listing_id,l.variation_id LIMIT $%d`, linkState, strings.Join(where, " AND "), limitArg)
	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return ports.ListingRowPage{}, fmt.Errorf("list listing rows: %w", err)
	}
	defer rows.Close()
	items := make([]domain.ListingReadModel, 0, q.Limit+1)
	for rows.Next() {
		var m domain.ListingReadModel
		var variation string
		var typeCode *domain.ListingTypeCode
		var priceAmount *string
		var currency *domain.PriceCurrency
		var syncJSON []byte
		var linkState domain.LinkState
		if err := rows.Scan(&m.InstallationID, &m.Provider, &m.ProviderListingID, &variation, &m.Title, &typeCode, &m.Status, &priceAmount, &currency, &m.PublishedQuantity, &m.SyncState, &syncJSON, &m.QualityScore, &m.Sales30D, &m.FetchedAt, &linkState, &m.Link.ProductID, &m.Link.SellerSKU); err != nil {
			return ports.ListingRowPage{}, fmt.Errorf("scan listing row: %w", err)
		}
		m.ListingID = domain.ListingID{InstallationID: m.InstallationID, ProviderListingID: m.ProviderListingID, VariationID: variation}.String()
		m.Link.State = linkState
		if typeCode != nil {
			if typ, ok := domain.ListingTypeForCode(*typeCode); ok {
				m.ListingType = &typ
			}
		}
		if priceAmount != nil && currency != nil {
			m.Price = &domain.Money{Amount: *priceAmount, Currency: *currency}
		}
		if len(syncJSON) > 0 {
			var se domain.ReadSyncError
			if err := json.Unmarshal(syncJSON, &se); err != nil {
				return ports.ListingRowPage{}, fmt.Errorf("decode listing sync error: %w", err)
			}
			m.SyncError = &se
		}
		items = append(items, m)
	}
	if err := rows.Err(); err != nil {
		return ports.ListingRowPage{}, fmt.Errorf("iterate listing rows: %w", err)
	}
	page := ports.ListingRowPage{Items: items}
	if len(items) > q.Limit {
		page.Items = items[:q.Limit]
		last := page.Items[len(page.Items)-1]
		c := ports.ListingCursor{LastTitle: last.Title, ListingID: last.ListingID}
		page.NextCursor = &c
	}
	return page, nil
}

func (*Repository) ListListingGroupRows(context.Context, ports.ListingGroupQuery) (ports.ListingGroupRowPage, error) {
	return ports.ListingGroupRowPage{}, errors.New("listings read: ListListingGroupRows not implemented (Slice 3/4/5)")
}
func (*Repository) GetListingRow(context.Context, domain.ListingKey) (domain.ListingReadModel, bool, error) {
	return domain.ListingReadModel{}, false, errors.New("listings read: GetListingRow not implemented (Slice 3/4/5)")
}
func (*Repository) GetListingsSummary(context.Context, ports.SummaryQuery) (ports.ListingSummaryRow, error) {
	return ports.ListingSummaryRow{}, errors.New("listings read: GetListingsSummary not implemented (Slice 3/4/5)")
}
func (*Repository) ListListingTimeline(context.Context, domain.ListingKey, int) ([]domain.TimelineEvent, error) {
	return nil, errors.New("listings read: ListListingTimeline not implemented (Slice 3/4/5)")
}

func (r *Repository) ApplyCompletedPull(ctx context.Context, installationID string, rows []domain.Listing, completedAt time.Time) error {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin listing completed pull: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	closedRows, err := tx.Query(ctx, `
		UPDATE listings SET status = 'closed', updated_at = $3
		WHERE tenant_id = $1 AND installation_id = $2
		RETURNING provider_listing_id, variation_id
	`, r.tenantID, installationID, completedAt.UTC())
	if err != nil {
		return fmt.Errorf("close listings: %w", err)
	}
	previous := make([]domain.ListingKey, 0)
	for closedRows.Next() {
		var key domain.ListingKey
		key.TenantID, key.InstallationID = r.tenantID, installationID
		if err := closedRows.Scan(&key.ProviderListingID, &key.VariationID); err != nil {
			closedRows.Close()
			return fmt.Errorf("scan closed listing key: %w", err)
		}
		previous = append(previous, key)
	}
	if err := closedRows.Err(); err != nil {
		closedRows.Close()
		return fmt.Errorf("iterate closed listing keys: %w", err)
	}
	closedRows.Close()

	returned := make(map[domain.ListingKey]struct{}, len(rows))
	for _, row := range rows {
		if row.Key.TenantID != r.tenantID || row.Key.InstallationID != installationID {
			return fmt.Errorf("listing key is outside repository tenant or installation")
		}
		syncError, err := json.Marshal(row.SyncError)
		if err != nil {
			return fmt.Errorf("marshal listing sync error: %w", err)
		}
		if row.SyncError == nil {
			syncError = nil
		}
		_, err = tx.Exec(ctx, `
			INSERT INTO listings (tenant_id, installation_id, provider, provider_listing_id, variation_id, title,
				listing_type_code, status, price_amount, price_currency, published_quantity, sync_state,
				sync_error, quality_score, sales_30d, fetched_at, created_at, updated_at)
			SELECT $1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18
			WHERE $1 = $19 AND $2 = $20
			ON CONFLICT (tenant_id, installation_id, provider_listing_id, variation_id) DO UPDATE SET
				provider=EXCLUDED.provider, title=EXCLUDED.title, listing_type_code=EXCLUDED.listing_type_code,
				status=EXCLUDED.status, price_amount=EXCLUDED.price_amount, price_currency=EXCLUDED.price_currency,
				published_quantity=EXCLUDED.published_quantity, sync_state=EXCLUDED.sync_state,
				sync_error=EXCLUDED.sync_error, quality_score=EXCLUDED.quality_score, sales_30d=EXCLUDED.sales_30d,
				fetched_at=EXCLUDED.fetched_at, updated_at=EXCLUDED.updated_at
		`, r.tenantID, installationID, row.Provider, row.Key.ProviderListingID, row.Key.VariationID, row.Title,
			row.ListingTypeCode, row.Status, row.PriceAmount, row.PriceCurrency, row.PublishedQuantity, row.SyncState,
			syncError, row.QualityScore, row.Sales30D, row.FetchedAt, row.CreatedAt.UTC(), completedAt.UTC(), r.tenantID, installationID)
		if err != nil {
			return fmt.Errorf("upsert listing %s/%s: %w", row.Key.ProviderListingID, row.Key.VariationID, err)
		}
		returned[row.Key] = struct{}{}
		kind, message := "synced", "Sincronizado"
		if row.Status == domain.ListingStatusPaused {
			kind, message = "paused", "Anúncio pausado no provedor"
		}
		if err := insertEvent(ctx, tx, row.Key, completedAt, kind, message); err != nil {
			return err
		}
	}
	for _, key := range previous {
		if _, ok := returned[key]; ok {
			continue
		}
		if err := insertEvent(ctx, tx, key, completedAt, "closed", "Anúncio ausente no provedor — fechado"); err != nil {
			return err
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit listing completed pull: %w", err)
	}
	return nil
}

func insertEvent(ctx context.Context, tx pgx.Tx, key domain.ListingKey, at time.Time, kind, message string) error {
	idBytes := make([]byte, 16)
	if _, err := rand.Read(idBytes); err != nil {
		return fmt.Errorf("create listing sync event id: %w", err)
	}
	_, err := tx.Exec(ctx, `
		INSERT INTO listing_sync_events (tenant_id, installation_id, provider_listing_id, variation_id, event_id, at, kind, message_pt)
		SELECT $1,$2,$3,$4,$5,$6,$7,$8 WHERE $1 = $9 AND $2 = $10
	`, key.TenantID, key.InstallationID, key.ProviderListingID, key.VariationID, hex.EncodeToString(idBytes), at.UTC(), kind, message, key.TenantID, key.InstallationID)
	if err != nil {
		return fmt.Errorf("insert listing sync event %s/%s: %w", key.ProviderListingID, key.VariationID, err)
	}
	return nil
}
