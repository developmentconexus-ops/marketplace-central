package domain

import "time"

type ShipmentInfo struct {
	ID            string         `json:"id"`
	Status        string         `json:"status"`
	Substatus     string         `json:"substatus"`
	SLADue        *time.Time     `json:"sla_due"`
	Delayed       *bool          `json:"delayed"`
	Costs         *ShipmentCosts `json:"costs"`
	DestinationUF *string        `json:"destination_uf"`
	FetchedAt     time.Time      `json:"fetched_at"`
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
