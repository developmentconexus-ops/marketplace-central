package domain

import "time"

type OrderRealizationState string

const (
	OrderRealizationRealized    OrderRealizationState = "realized"
	OrderRealizationNotRealized OrderRealizationState = "not_realized"
	OrderRealizationUnknown     OrderRealizationState = "unknown"
)

type OrderLinkQuality string

const (
	OrderLinkResolved   OrderLinkQuality = "resolved"
	OrderLinkRejected   OrderLinkQuality = "rejected"
	OrderLinkConflict   OrderLinkQuality = "conflict"
	OrderLinkUnresolved OrderLinkQuality = "unresolved"
	OrderLinkMissing    OrderLinkQuality = "missing"
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
	ProviderItemID      string
	ProviderVariationID string
	Quantity            int
	UnitPrice           *float64
	SaleFeeAmount       *float64
	LinkQuality         OrderLinkQuality
	InternalProductID   *int
}
