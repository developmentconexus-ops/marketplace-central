package domain

import "time"

// OrderShipment is the order_shipments (migration 0088, IC-03) persistence fact — a
// module-owned shape, deliberately NOT a reuse of connectorsdomain.ShipmentDetail. Same
// precedent as MarketplaceOrder not reusing connectorsdomain.OrderDetail (order.go): connectors
// types stay behind the ports boundary, this package owns what gets persisted.
//
// No raw column by design (0088, P7 r01 B-7): a shipment payload carries delivery PII. Every
// optional field is a nil pointer (ADR-17) — provider absence is never a fabricated zero/blank.
type OrderShipment struct {
	ProviderCode       string
	ProviderShipmentID string
	// ProviderOrderID is a soft link (0088 carries no FK to orders_marketplace_orders — see
	// migration comment); IngestOrder is what keeps the two rows correlated in practice, one
	// order fetch away from one shipment fetch, same transaction.
	ProviderOrderID string

	Status    string
	Substatus *string

	LogisticType   *string
	TrackingMethod *string
	TrackingNumber *string

	SLAStatus  *string
	SLALimitAt *time.Time

	// CostGross/CostSeller are the parsed decimal amounts (order_shipments.cost_gross/
	// cost_seller are numeric(14,2)); connectorsdomain.ShipmentDetail carries these as decimal
	// strings (its own doc comment explains why) — IngestOrder parses them once at the ingest
	// boundary so the repository writes a plain numeric parameter, same as every other money
	// column in this module (nullableFloat8).
	CostGross  *float64
	CostSeller *float64
	Currency   *string

	ReceiverName *string
	DestCity     *string
	DestState    *string
	DestZip      *string

	// SourceTime is the provider's own timestamp for this shipment fact (date_created,
	// verbatim). FetchedAt is NOT NULL on 0088 and is this ingest's own clock — it is what
	// upsertOrderShipment's replay/freshness guard is built on, since SourceTime can be nil.
	SourceTime *time.Time
	FetchedAt  time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}
