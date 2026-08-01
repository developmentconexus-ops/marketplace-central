package postgres

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	connectorsdomain "marketplace-central/apps/server_core/internal/modules/connectors/domain"
	"marketplace-central/apps/server_core/internal/modules/orders/ports"
)

// shipmentReaderProviderCode matches the literal used everywhere else this module writes/reads
// order_shipments (order_repo.go upsertOrderShipment via ordersdomain.OrderShipment.ProviderCode,
// ingest_repo_test.go fixtures) — order_shipments carries no CHECK on `provider` (0088: the
// vocabulary belongs to the provider, not this repo), so the literal is repeated here rather than
// introduced as a new shared constant no other call site uses yet.
const shipmentReaderProviderCode = "mercado_livre"

// ShipmentReader backs ports.ShipmentReader (F-03 read-path switch) by reading the order_shipments
// table (migration 0088, populated by F-02's IngestOrder) instead of a live Mercado Livre call.
// installationID is accepted only to satisfy the port signature — order_shipments carries no
// installation_id column (0088's PK is tenant_id/provider/provider_shipment_id), so tenant_id +
// shipmentID alone is the query's scoping, same as every other tenant-scoped query in this
// package.
type ShipmentReader struct {
	pool     *pgxpool.Pool
	tenantID string
}

var _ ports.ShipmentReader = (*ShipmentReader)(nil)

func NewShipmentReader(pool *pgxpool.Pool, tenantID string) *ShipmentReader {
	return &ShipmentReader{pool: pool, tenantID: tenantID}
}

// GetShipment reads the persisted shipment fact for shipmentID. A row that does not exist yet
// (order ingested before F-02, or the order carries no shipment) is honest-unknown, never an
// error: it returns ports.ErrShipmentNotFound, which EnrichService.fetchShipment recognizes and
// degrades to a nil ShipmentEnrichment WITHOUT a warn log (unlike a real read failure, which still
// warns) — a not-yet-ingested order is an expected, common case, not an operational fault.
func (r *ShipmentReader) GetShipment(ctx context.Context, installationID, shipmentID string) (connectorsdomain.ShipmentInfo, error) {
	if r.pool == nil {
		return connectorsdomain.ShipmentInfo{}, ports.ErrShipmentNotFound
	}
	row := r.pool.QueryRow(ctx, `
		SELECT provider_shipment_id, status, substatus, sla_limit_at,
		       dest_state, dest_city, dest_zip, receiver_name,
		       cost_gross::text, cost_seller::text, currency, fetched_at
		FROM order_shipments
		WHERE tenant_id = $1 AND provider = $2 AND provider_shipment_id = $3
	`, r.tenantID, shipmentReaderProviderCode, strings.TrimSpace(shipmentID))
	info, err := scanShipmentRow(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return connectorsdomain.ShipmentInfo{}, ports.ErrShipmentNotFound
	}
	return info, err
}

// scanShipmentRow maps one order_shipments row onto connectorsdomain.ShipmentInfo. Column order
// mirrors the SELECT in GetShipment exactly. Split out as a pure function (generic Scan-capable
// interface, same pattern as scanReadModel/scanOrder in order_repo.go) so the mapping — including
// the Money.Amount formatting the F-03 brief calls out as load-bearing for retorno_liquido — is
// unit-testable without a database (see shipment_reader_test.go).
//
// Delayed/CarrierName/TrackingURL/Costs.ReceiverCost are ALWAYS nil: migration 0088 has no
// delayed/carrier/tracking-url/receiver-cost columns (see the column-mapping table in
// F-03-read-path-switch/feature.md and validation.md for the accepted, out-of-scope UX gap this
// produces on the rastreio/frete_real DTOs).
func scanShipmentRow(scanner interface{ Scan(dest ...any) error }) (connectorsdomain.ShipmentInfo, error) {
	var (
		providerShipmentID string
		status             string
		substatus          pgtype.Text
		slaLimitAt         pgtype.Timestamptz
		destState          pgtype.Text
		destCity           pgtype.Text
		destZip            pgtype.Text
		receiverName       pgtype.Text
		costGrossText      pgtype.Text
		costSellerText     pgtype.Text
		currency           pgtype.Text
		fetchedAt          pgtype.Timestamptz
	)
	if err := scanner.Scan(
		&providerShipmentID, &status, &substatus, &slaLimitAt,
		&destState, &destCity, &destZip, &receiverName,
		&costGrossText, &costSellerText, &currency, &fetchedAt,
	); err != nil {
		return connectorsdomain.ShipmentInfo{}, err
	}

	info := connectorsdomain.ShipmentInfo{
		ID:              providerShipmentID,
		Status:          status,
		Substatus:       textValue(substatus),
		SLADue:          scanTime(slaLimitAt),
		DestinationUF:   textPtr(destState),
		DestinationCity: textPtr(destCity),
		DestinationZip:  textPtr(destZip),
		ReceiverName:    textPtr(receiverName),
		FetchedAt:       fetchedAt.Time.UTC(),
	}

	// cost_gross/cost_seller are cast to ::text in the SELECT so Postgres formats the
	// numeric(14,2) value directly (e.g. "42.50") — no float64 round-trip, so no risk of the
	// scientific-notation/imprecision failure mode domain.ValidateMoney's decimal regex rejects.
	// A cost value present without a valid currency is honest-unknown (ADR-17): Money always
	// needs both, never a fabricated blank currency.
	currencyValue := textValue(currency)
	var costs connectorsdomain.ShipmentCosts
	hasCosts := false
	if costGrossText.Valid && currencyValue != "" {
		costs.GrossAmount = &connectorsdomain.Money{Amount: costGrossText.String, Currency: currencyValue}
		hasCosts = true
	}
	if costSellerText.Valid && currencyValue != "" {
		costs.SenderCost = &connectorsdomain.Money{Amount: costSellerText.String, Currency: currencyValue}
		hasCosts = true
	}
	if hasCosts {
		info.Costs = &costs
	}
	return info, nil
}

func textPtr(v pgtype.Text) *string {
	if !v.Valid {
		return nil
	}
	s := v.String
	return &s
}

func textValue(v pgtype.Text) string {
	if !v.Valid {
		return ""
	}
	return v.String
}
