package domain

import "time"

type ShipmentInfo struct {
	ID        string         `json:"id"`
	Status    string         `json:"status"`
	Substatus string         `json:"substatus"`
	SLADue    *time.Time     `json:"sla_due"`
	Delayed   *bool          `json:"delayed"`
	Costs     *ShipmentCosts `json:"costs"`
	// Destination/carrier facts are marketplace-neutral: no provider payload
	// shape (ML "destination"/"shipping_address"/"carrier") leaks past the
	// adapter. Every field is a pointer — nil is honest absence, and a
	// masked-but-present provider value that trims to empty degrades to nil, not
	// a fabricated blank (ADR-17; ML obfuscates buyer address until payment is
	// confirmed — masked ≠ error).
	DestinationUF   *string   `json:"destination_uf"`
	DestinationCity *string   `json:"destination_city"`
	DestinationZip  *string   `json:"destination_zip"`
	ReceiverName    *string   `json:"receiver_name"`
	CarrierName     *string   `json:"carrier_name"`
	TrackingURL     *string   `json:"tracking_url"`
	FetchedAt       time.Time `json:"fetched_at"`
}

type ShipmentCosts struct {
	GrossAmount  *Money `json:"gross_amount"`
	ReceiverCost *Money `json:"receiver_cost"`
	SenderCost   *Money `json:"sender_cost"`
}

type FreeShippingQuery struct {
	ItemID string `json:"item_id"`
}

type FreeShippingCost struct {
	Cost      *Money    `json:"cost"`
	FetchedAt time.Time `json:"fetched_at"`
}
