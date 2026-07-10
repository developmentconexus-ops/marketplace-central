package domain

import "time"

// OrderIngestionSnapshot is the sanitized order data contract consumed by orders.
type OrderIngestionSnapshot struct {
	ProviderCode         string                  `json:"provider_code"`
	ProviderOrderID      string                  `json:"provider_order_id"`
	ProviderStatus       string                  `json:"provider_status,omitempty"`
	ProviderStatusDetail string                  `json:"provider_status_detail,omitempty"`
	ProviderCreatedAt    *time.Time              `json:"provider_created_at,omitempty"`
	ProviderClosedAt     *time.Time              `json:"provider_closed_at,omitempty"`
	ProviderUpdatedAt    *time.Time              `json:"provider_updated_at,omitempty"`
	FetchedAt            time.Time               `json:"fetched_at"`
	Items                []OrderIngestionItem    `json:"items"`
	SaleFeeAmount        *float64                `json:"sale_fee_amount,omitempty"`
	Payments             []OrderIngestionPayment `json:"payments"`
	ShippingID           string                  `json:"shipping_id,omitempty"`
	CancellationDetail   string                  `json:"cancellation_detail,omitempty"`
	Tags                 []string                `json:"tags,omitempty"`
}

type OrderIngestionItem struct {
	ProviderItemID      string   `json:"provider_item_id"`
	ProviderVariationID string   `json:"provider_variation_id,omitempty"`
	SellerSKU           string   `json:"seller_sku,omitempty"`
	EAN                 string   `json:"ean,omitempty"`
	Title               string   `json:"title,omitempty"`
	Quantity            int      `json:"quantity"`
	UnitPrice           *float64 `json:"unit_price,omitempty"`
	SaleFeeAmount       *float64 `json:"sale_fee_amount,omitempty"`
}

type OrderIngestionPayment struct {
	PaymentID         string   `json:"payment_id"`
	Status            string   `json:"status,omitempty"`
	Amount            *float64 `json:"amount,omitempty"`
	TransactionAmount *float64 `json:"transaction_amount,omitempty"`
	TotalPaidAmount   *float64 `json:"total_paid_amount,omitempty"`
}
