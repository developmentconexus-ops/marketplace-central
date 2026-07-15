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

type listingRowScanner interface {
	Scan(...any) error
}

const listingLinkState = "COALESCE(NULLIF(pl.state, 'none'), 'unresolved')"

func listingProjectionSQL() string {
	return fmt.Sprintf(`SELECT l.installation_id,l.provider,l.provider_listing_id,l.variation_id,l.title,l.listing_type_code,l.status,
		l.price_amount::text,l.price_currency,l.published_quantity,l.sync_state,l.sync_error,l.quality_score,l.sales_30d,l.fetched_at,
		%s,CASE WHEN pl.state='resolved' THEN pl.internal_product_id::text END,NULLIF(pls.seller_sku,'')
	FROM listings l
	LEFT JOIN product_links pl ON pl.tenant_id=l.tenant_id AND pl.installation_id=l.installation_id AND pl.provider_item_id=l.provider_listing_id AND pl.provider_variation_id=CASE WHEN l.variation_id='-' THEN '' ELSE l.variation_id END
	LEFT JOIN product_link_listing_snapshots pls ON pls.tenant_id=l.tenant_id AND pls.installation_id=l.installation_id AND pls.provider_item_id=l.provider_listing_id AND pls.provider_variation_id=CASE WHEN l.variation_id='-' THEN '' ELSE l.variation_id END`, listingLinkState)
}

func scanListingReadModel(row listingRowScanner) (domain.ListingReadModel, error) {
	var m domain.ListingReadModel
	var variation string
	var typeCode *domain.ListingTypeCode
	var priceAmount *string
	var currency *domain.PriceCurrency
	var syncJSON []byte
	var linkState domain.LinkState
	if err := row.Scan(&m.InstallationID, &m.Provider, &m.ProviderListingID, &variation, &m.Title, &typeCode, &m.Status, &priceAmount, &currency, &m.PublishedQuantity, &m.SyncState, &syncJSON, &m.QualityScore, &m.Sales30D, &m.FetchedAt, &linkState, &m.Link.ProductID, &m.Link.SellerSKU); err != nil {
		return domain.ListingReadModel{}, err
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
			return domain.ListingReadModel{}, fmt.Errorf("decode listing sync error: %w", err)
		}
		m.SyncError = &se
	}
	return m, nil
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
	if q.Filter.LinkState != "" {
		add(listingLinkState+" = $%d", q.Filter.LinkState)
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
		where = append(where, listingLinkState+" = 'unresolved'")
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
	sql := fmt.Sprintf(`%s WHERE %s ORDER BY l.title,l.provider_listing_id,l.variation_id LIMIT $%d`, listingProjectionSQL(), strings.Join(where, " AND "), limitArg)
	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return ports.ListingRowPage{}, fmt.Errorf("list listing rows: %w", err)
	}
	defer rows.Close()
	items := make([]domain.ListingReadModel, 0, q.Limit+1)
	for rows.Next() {
		m, err := scanListingReadModel(rows)
		if err != nil {
			return ports.ListingRowPage{}, fmt.Errorf("scan listing row: %w", err)
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

func (r *Repository) ListListingGroupRows(ctx context.Context, q ports.ListingGroupQuery) (ports.ListingGroupRowPage, error) {
	args, where := listingGroupWhere(r.tenantID, q)
	if !q.Cursor.IsFirstPage() {
		args = append(args, q.Cursor.NullLast, q.Cursor.ProductTitle, q.Cursor.ProductID)
		n := len(args)
		where = append(where, fmt.Sprintf("(CASE WHEN pl.state='resolved' AND pl.internal_product_id IS NOT NULL THEN false ELSE true END, CASE WHEN pl.state='resolved' AND pl.internal_product_id IS NOT NULL THEN pl.internal_product_name ELSE '' END, CASE WHEN pl.state='resolved' AND pl.internal_product_id IS NOT NULL THEN pl.internal_product_id::text ELSE '' END) > ($%d,$%d,$%d)", n-2, n-1, n))
	}
	args = append(args, q.Limit+1)
	keySQL := fmt.Sprintf(`SELECT DISTINCT
		CASE WHEN pl.state='resolved' AND pl.internal_product_id IS NOT NULL THEN false ELSE true END AS null_last,
		CASE WHEN pl.state='resolved' AND pl.internal_product_id IS NOT NULL THEN pl.internal_product_name ELSE '' END AS product_title,
		CASE WHEN pl.state='resolved' AND pl.internal_product_id IS NOT NULL THEN pl.internal_product_id::text ELSE '' END AS product_id
	FROM listings l
	LEFT JOIN product_links pl ON pl.tenant_id=l.tenant_id AND pl.installation_id=l.installation_id AND pl.provider_item_id=l.provider_listing_id AND pl.provider_variation_id=CASE WHEN l.variation_id='-' THEN '' ELSE l.variation_id END
	LEFT JOIN product_link_listing_snapshots pls ON pls.tenant_id=l.tenant_id AND pls.installation_id=l.installation_id AND pls.provider_item_id=l.provider_listing_id AND pls.provider_variation_id=CASE WHEN l.variation_id='-' THEN '' ELSE l.variation_id END
	WHERE %s ORDER BY null_last,product_title,product_id LIMIT $%d`, strings.Join(where, " AND "), len(args))
	rows, err := r.pool.Query(ctx, keySQL, args...)
	if err != nil {
		return ports.ListingGroupRowPage{}, fmt.Errorf("list listing group keys: %w", err)
	}
	keys := make([]ports.GroupCursor, 0, q.Limit+1)
	for rows.Next() {
		var k ports.GroupCursor
		if err := rows.Scan(&k.NullLast, &k.ProductTitle, &k.ProductID); err != nil {
			rows.Close()
			return ports.ListingGroupRowPage{}, fmt.Errorf("scan listing group key: %w", err)
		}
		keys = append(keys, k)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return ports.ListingGroupRowPage{}, fmt.Errorf("iterate listing group keys: %w", err)
	}
	rows.Close()
	page := ports.ListingGroupRowPage{}
	if len(keys) > q.Limit {
		keys = keys[:q.Limit]
		c := keys[len(keys)-1]
		page.NextCursor = &c
	}
	if len(keys) == 0 {
		page.Groups = []domain.ListingGroup{}
		return page, nil
	}

	childArgs, childWhere := listingGroupWhere(r.tenantID, q)
	groups := make([]domain.ListingGroup, len(keys))
	keyIndex := make(map[string]int, len(keys))
	selectors := make([]string, 0, len(keys))
	for i, k := range keys {
		if k.NullLast {
			title := "sem produto"
			groups[i] = domain.ListingGroup{ProductTitle: &title, Listings: []domain.ListingReadModel{}}
			selectors = append(selectors, "NOT (pl.state='resolved' AND pl.internal_product_id IS NOT NULL)")
			keyIndex["null"] = i
			continue
		}
		id, title := k.ProductID, k.ProductTitle
		groups[i] = domain.ListingGroup{ProductID: &id, ProductTitle: &title, Listings: []domain.ListingReadModel{}}
		childArgs = append(childArgs, k.ProductTitle, k.ProductID)
		n := len(childArgs)
		selectors = append(selectors, fmt.Sprintf("(pl.state='resolved' AND pl.internal_product_id IS NOT NULL AND pl.internal_product_name=$%d AND pl.internal_product_id::text=$%d)", n-1, n))
		keyIndex[k.ProductTitle+"\x00"+k.ProductID] = i
	}
	childWhere = append(childWhere, "("+strings.Join(selectors, " OR ")+")")
	projection := strings.Replace(listingProjectionSQL(), "SELECT ", `SELECT CASE WHEN pl.state='resolved' AND pl.internal_product_id IS NOT NULL THEN pl.internal_product_name END, CASE WHEN pl.state='resolved' AND pl.internal_product_id IS NOT NULL THEN pl.internal_product_id::text END, `, 1)
	childSQL := fmt.Sprintf(`%s WHERE %s ORDER BY l.title,l.provider_listing_id,l.variation_id`, projection, strings.Join(childWhere, " AND "))
	rows, err = r.pool.Query(ctx, childSQL, childArgs...)
	if err != nil {
		return ports.ListingGroupRowPage{}, fmt.Errorf("list listing group children: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var title, id *string
		wrapped := &prefixedListingScanner{row: rows, prefix: []any{&title, &id}}
		item, err := scanListingReadModel(wrapped)
		if err != nil {
			return ports.ListingGroupRowPage{}, fmt.Errorf("scan listing group child: %w", err)
		}
		key := "null"
		if id != nil && title != nil {
			key = *title + "\x00" + *id
		}
		groups[keyIndex[key]].Listings = append(groups[keyIndex[key]].Listings, item)
	}
	if err := rows.Err(); err != nil {
		return ports.ListingGroupRowPage{}, fmt.Errorf("iterate listing group children: %w", err)
	}
	page.Groups = groups
	return page, nil
}

type prefixedListingScanner struct {
	row    pgx.Rows
	prefix []any
}

func (s *prefixedListingScanner) Scan(dest ...any) error {
	return s.row.Scan(append(s.prefix, dest...)...)
}

func listingGroupWhere(tenantID string, q ports.ListingGroupQuery) ([]any, []string) {
	args := []any{tenantID, q.InstallationID}
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
	if q.Filter.LinkState != "" {
		add(listingLinkState+" = $%d", q.Filter.LinkState)
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
		where = append(where, listingLinkState+" = 'unresolved'")
	}
	if q.Q != "" {
		escaped := strings.NewReplacer("\\", "\\\\", "%", "\\%", "_", "\\_").Replace(q.Q)
		args = append(args, escaped)
		n := len(args)
		where = append(where, fmt.Sprintf("(l.title ILIKE '%%%%' || $%d || '%%%%' ESCAPE E'\\\\' OR l.provider_listing_id ILIKE '%%%%' || $%d || '%%%%' ESCAPE E'\\\\' OR pls.seller_sku ILIKE '%%%%' || $%d || '%%%%' ESCAPE E'\\\\')", n, n, n))
	}
	return args, where
}
func (r *Repository) GetListingRow(ctx context.Context, key domain.ListingKey) (domain.ListingReadModel, bool, error) {
	row := r.pool.QueryRow(ctx, fmt.Sprintf(`%s
	WHERE l.tenant_id=$1 AND l.installation_id=$2 AND l.provider_listing_id=$3 AND l.variation_id=$4`, listingProjectionSQL()), r.tenantID, key.InstallationID, key.ProviderListingID, key.VariationID)
	model, err := scanListingReadModel(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ListingReadModel{}, false, nil
		}
		return domain.ListingReadModel{}, false, fmt.Errorf("get listing row: %w", err)
	}
	return model, true, nil
}
func (*Repository) GetListingsSummary(context.Context, ports.SummaryQuery) (ports.ListingSummaryRow, error) {
	return ports.ListingSummaryRow{}, errors.New("listings read: GetListingsSummary not implemented (Slice 3/4/5)")
}
func (r *Repository) ListListingTimeline(ctx context.Context, key domain.ListingKey, limit int) ([]domain.TimelineEvent, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT at,kind,message_pt
		FROM listing_sync_events
		WHERE tenant_id=$1 AND installation_id=$2 AND provider_listing_id=$3 AND variation_id=$4
		ORDER BY at DESC,event_id DESC
		LIMIT $5`, r.tenantID, key.InstallationID, key.ProviderListingID, key.VariationID, limit)
	if err != nil {
		return nil, fmt.Errorf("list listing timeline: %w", err)
	}
	defer rows.Close()
	events := make([]domain.TimelineEvent, 0)
	for rows.Next() {
		var event domain.TimelineEvent
		if err := rows.Scan(&event.At, &event.Kind, &event.MessagePT); err != nil {
			return nil, fmt.Errorf("scan listing timeline: %w", err)
		}
		events = append(events, event)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate listing timeline: %w", err)
	}
	return events, nil
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
