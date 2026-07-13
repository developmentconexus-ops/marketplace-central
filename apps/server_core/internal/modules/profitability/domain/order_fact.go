package domain

import "time"

type OrderRealizationState string

const (
	OrderRealizationRealized    OrderRealizationState = "realized"
	OrderRealizationNotRealized OrderRealizationState = "not_realized"
	OrderRealizationUnknown     OrderRealizationState = "unknown"
)

type OrderLinkQuality string
type OrderLineReconciliationState string

const (
	OrderLinkResolved   OrderLinkQuality = "resolved"
	OrderLinkRejected   OrderLinkQuality = "rejected"
	OrderLinkConflict   OrderLinkQuality = "conflict"
	OrderLinkUnresolved OrderLinkQuality = "unresolved"
	OrderLinkMissing    OrderLinkQuality = "missing"

	OrderLineStable           OrderLineReconciliationState = "stable"
	OrderLineLegacyUnresolved OrderLineReconciliationState = "legacy_unresolved"
	OrderLineAmbiguous        OrderLineReconciliationState = "ambiguous"
)

type OrderFact struct {
	InstallationID    string
	ProviderOrderID   string
	RealizationState  OrderRealizationState
	ProviderCreatedAt *time.Time
	ProviderClosedAt  *time.Time
	ProviderUpdatedAt *time.Time
	FetchedAt         time.Time
	Items             []OrderItemFact
}

type OrderItemFact struct {
	MPCLineID           string
	ReconciliationState OrderLineReconciliationState
	ProviderItemID      string
	ProviderVariationID string
	Quantity            int
	UnitPrice           *float64
	SaleFeeAmount       *float64
	LinkQuality         OrderLinkQuality
	InternalProductID   *int
}
