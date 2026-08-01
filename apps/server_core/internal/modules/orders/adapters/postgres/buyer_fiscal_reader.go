package postgres

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	connectorsdomain "marketplace-central/apps/server_core/internal/modules/connectors/domain"
	"marketplace-central/apps/server_core/internal/modules/orders/ports"
)

// BuyerFiscalReader backs ports.BuyerFiscalReader (F-03 read-path switch) by reading the 9
// buyer_* columns on orders_marketplace_orders (migration 0089, populated by F-02's IngestOrder)
// instead of a live Mercado Livre call.
type BuyerFiscalReader struct {
	pool     *pgxpool.Pool
	tenantID string
}

var _ ports.BuyerFiscalReader = (*BuyerFiscalReader)(nil)

func NewBuyerFiscalReader(pool *pgxpool.Pool, tenantID string) *BuyerFiscalReader {
	return &BuyerFiscalReader{pool: pool, tenantID: tenantID}
}

// GetBuyerFiscal reads the persisted buyer-fiscal columns for one order, scoped by
// (tenant_id, installation_id, provider_order_id) — the same tenancy predicate idiom
// OrderRepository.GetFaturadoAt already uses for a single-order lookup. Two distinct
// "not present" cases per the port's doc comment (buyer_fiscal_reader.go) and the F-03 brief:
//
//   - The order row itself does not exist (bad/unknown providerOrderID): a REAL failure, not
//     absence — returned as the underlying pgx.ErrNoRows, same as GetFaturadoAt would surface it
//     had it not special-cased first-ingest. EnrichService.resolveBuyerFiscal already treats any
//     error here as a real failure (degrades to nil AND warns once).
//   - The order row exists but all 9 buyer_* columns are NULL (ML never had fiscal data for this
//     buyer, or the order predates F-02's ingest): honest absence — returns a BuyerFiscalInfo
//     whose HasData() is false, which resolveBuyerFiscal already treats as nil with no warn.
func (r *BuyerFiscalReader) GetBuyerFiscal(ctx context.Context, installationID, providerOrderID string) (connectorsdomain.BuyerFiscalInfo, error) {
	if r.pool == nil {
		return connectorsdomain.BuyerFiscalInfo{}, errors.New("orders buyer fiscal reader: database is not configured")
	}
	var (
		name, docType, docNumber                      pgtype.Text
		street, number, city, stateCode, zip, country pgtype.Text
	)
	err := r.pool.QueryRow(ctx, `
		SELECT buyer_name, buyer_doc_type, buyer_doc_number,
		       buyer_address_street, buyer_address_number, buyer_address_city,
		       buyer_address_state_code, buyer_address_zip, buyer_address_country
		FROM orders_marketplace_orders
		WHERE tenant_id = $1 AND installation_id = $2 AND provider_order_id = $3
	`, r.tenantID, strings.TrimSpace(installationID), strings.TrimSpace(providerOrderID)).Scan(
		&name, &docType, &docNumber,
		&street, &number, &city, &stateCode, &zip, &country,
	)
	if err != nil {
		return connectorsdomain.BuyerFiscalInfo{}, err
	}
	return buildBuyerFiscalInfo(name, docType, docNumber, street, number, city, stateCode, zip, country), nil
}

// buildBuyerFiscalInfo is split out as a pure function (same rationale as scanShipmentRow) so the
// honest-absence collapsing rule is unit-testable without a database: an address is included ONLY
// when at least one of its six fields is present, because BuyerFiscalInfo.HasData() checks
// `Address != nil` — a non-nil-but-all-empty address would flip HasData() to true for a buyer with
// zero fiscal data, turning honest absence into a fabricated empty block (ADR-17). StateName
// (uf_nome) is always nil: migration 0089 deliberately excludes it (see the migration's own
// comment) — VERIFIED safe because apps/web/src/pages/pedidos/PedidoDrawer.tsx's formatEndereco
// (around line 341-347) only reads uf_codigo, never uf_nome.
func buildBuyerFiscalInfo(name, docType, docNumber, street, number, city, stateCode, zip, country pgtype.Text) connectorsdomain.BuyerFiscalInfo {
	info := connectorsdomain.BuyerFiscalInfo{
		Name:      textPtr(name),
		DocType:   textPtr(docType),
		DocNumber: textPtr(docNumber),
	}
	if street.Valid || number.Valid || city.Valid || stateCode.Valid || zip.Valid || country.Valid {
		info.Address = &connectorsdomain.BuyerFiscalAddress{
			StreetName:   textPtr(street),
			StreetNumber: textPtr(number),
			City:         textPtr(city),
			StateCode:    textPtr(stateCode),
			StateName:    nil,
			ZipCode:      textPtr(zip),
			CountryID:    textPtr(country),
		}
	}
	return info
}
