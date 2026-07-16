package domain

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type ProviderCode string
type CapabilityStatus string
type StockScope string
type StockWriteResultStatus string
type ProviderOperationState string
type ErrorCode string

const (
	CapabilityStatusSupported   CapabilityStatus = "supported"
	CapabilityStatusUnsupported CapabilityStatus = "unsupported"
	CapabilityStatusDegraded    CapabilityStatus = "degraded"
	CapabilityStatusBlocked     CapabilityStatus = "blocked"

	StockScopeItem      StockScope = "item"
	StockScopeVariation StockScope = "variation"

	StockWriteResultApplied          StockWriteResultStatus = "applied"
	StockWriteResultRejected         StockWriteResultStatus = "rejected"
	StockWriteResultTransientFailure StockWriteResultStatus = "transient_failure"
	StockWriteResultUnsupportedShape StockWriteResultStatus = "unsupported_shape"

	ProviderOperationStatePending   ProviderOperationState = "pending"
	ProviderOperationStateRunning   ProviderOperationState = "running"
	ProviderOperationStateSucceeded ProviderOperationState = "succeeded"
	ProviderOperationStateFailed    ProviderOperationState = "failed"
	ProviderOperationStateSkipped   ProviderOperationState = "skipped"
)

const (
	ErrCodeProviderRateLimited      ErrorCode = "CONNECTORS_PROVIDER_RATE_LIMITED"
	ErrCodeProviderAuth             ErrorCode = "CONNECTORS_PROVIDER_AUTH"
	ErrCodeProviderValidation       ErrorCode = "CONNECTORS_PROVIDER_VALIDATION"
	ErrCodeProviderUnsupportedShape ErrorCode = "CONNECTORS_PROVIDER_UNSUPPORTED_SHAPE"
	ErrCodeProviderTransient        ErrorCode = "CONNECTORS_PROVIDER_TRANSIENT"
	ErrCodeProviderPayloadInvalid   ErrorCode = "CONNECTORS_PROVIDER_PAYLOAD_INVALID"
	ErrCodeProviderInvalidReference ErrorCode = "CONNECTORS_PROVIDER_INVALID_REFERENCE"
)

type CapabilityError struct {
	Code    ErrorCode
	Message string
	Cause   error
}

func (e *CapabilityError) Error() string {
	if e == nil {
		return ""
	}
	if strings.TrimSpace(e.Message) == "" {
		return string(e.Code)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *CapabilityError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Cause
}

func NewCapabilityError(code ErrorCode, message string) error {
	return &CapabilityError{
		Code:    code,
		Message: strings.TrimSpace(message),
	}
}

func ErrorCodeOf(err error) ErrorCode {
	var target *CapabilityError
	if errors.As(err, &target) && target != nil {
		return target.Code
	}
	return ""
}

type ProviderAccountRef struct {
	TenantID          string
	InstallationID    string
	ProviderCode      string
	ProviderAccountID string
}

type ProviderListingRef struct {
	AccountRef          ProviderAccountRef
	ProviderItemID      string
	ProviderVariationID string
}

type ProviderOrderRef struct {
	AccountRef      ProviderAccountRef
	ProviderOrderID string
}

type ListListingsInput struct {
	AccountRef ProviderAccountRef
	Cursor     string
	Status     string
	Limit      int
}

type ListOrdersInput struct {
	AccountRef   ProviderAccountRef
	Cursor       string
	Status       string
	UpdatedAfter *time.Time
	Limit        int
}

type ListingSnapshot struct {
	ProviderCode        string                     `json:"provider_code"`
	ProviderItemID      string                     `json:"provider_item_id"`
	ProviderVariationID string                     `json:"provider_variation_id,omitempty"`
	ProviderStatus      string                     `json:"provider_status,omitempty"`
	PriceAmount         *string                    `json:"price_amount,omitempty"`
	PriceCurrency       string                     `json:"price_currency,omitempty"`
	ListingTypeCode     string                     `json:"listing_type_code,omitempty"`
	SellerSKU           string                     `json:"seller_sku,omitempty"`
	EAN                 string                     `json:"ean,omitempty"`
	Title               string                     `json:"title,omitempty"`
	AvailableQuantity   *int                       `json:"available_quantity,omitempty"`
	SourceUpdatedAt     *time.Time                 `json:"source_updated_at,omitempty"`
	FetchedAt           time.Time                  `json:"fetched_at"`
	RawProviderRef      any                        `json:"raw_provider_ref,omitempty"`
	Variations          []ListingVariationSnapshot `json:"variations,omitempty"`
}

type AccountSnapshot struct {
	ProviderCode        string    `json:"provider_code"`
	ProviderAccountID   string    `json:"provider_account_id"`
	ProviderAccountName string    `json:"provider_account_name"`
	SiteID              string    `json:"site_id"`
	Status              string    `json:"status"`
	FetchedAt           time.Time `json:"fetched_at"`
	RawProviderRef      any       `json:"raw_provider_ref,omitempty"`
}

type ListingVariationSnapshot struct {
	ProviderVariationID string `json:"provider_variation_id"`
	SellerSKU           string `json:"seller_sku,omitempty"`
	EAN                 string `json:"ean,omitempty"`
	AvailableQuantity   *int   `json:"available_quantity,omitempty"`
}

type StockSnapshot struct {
	ProviderCode        string     `json:"provider_code"`
	ProviderItemID      string     `json:"provider_item_id"`
	ProviderVariationID string     `json:"provider_variation_id,omitempty"`
	AvailableQuantity   int        `json:"available_quantity"`
	ProviderStatus      string     `json:"provider_status,omitempty"`
	SellerSKU           string     `json:"seller_sku,omitempty"`
	EAN                 string     `json:"ean,omitempty"`
	Title               string     `json:"title,omitempty"`
	SourceUpdatedAt     *time.Time `json:"source_updated_at,omitempty"`
	FetchedAt           time.Time  `json:"fetched_at"`
	RawProviderRef      any        `json:"raw_provider_ref,omitempty"`
	Scope               StockScope `json:"scope"`
}

type StockWriteRequest struct {
	TenantID            string
	InstallationID      string
	ProviderCode        string
	ProviderAccountID   string
	ProviderItemID      string
	ProviderVariationID string
	IdempotencyKey      string
	RequestedQuantity   int
	Reason              string
}

type StockWriteResult struct {
	IDempotencyKey      string                 `json:"idempotency_key"`
	Result              StockWriteResultStatus `json:"result"`
	ProviderCode        string                 `json:"provider_code"`
	ProviderItemID      string                 `json:"provider_item_id"`
	ProviderVariationID string                 `json:"provider_variation_id,omitempty"`
	Message             string                 `json:"message,omitempty"`
	ProviderResponseRef any                    `json:"provider_response_ref,omitempty"`
}

type OrderSnapshot struct {
	ProviderCode         string                 `json:"provider_code"`
	ProviderOrderID      string                 `json:"provider_order_id"`
	ProviderStatus       string                 `json:"provider_status,omitempty"`
	ProviderStatusDetail string                 `json:"provider_status_detail,omitempty"`
	ProviderCreatedAt    *time.Time             `json:"provider_created_at,omitempty"`
	ProviderClosedAt     *time.Time             `json:"provider_closed_at,omitempty"`
	ProviderUpdatedAt    *time.Time             `json:"provider_updated_at,omitempty"`
	FetchedAt            time.Time              `json:"fetched_at"`
	RawProviderRef       any                    `json:"raw_provider_ref,omitempty"`
	Items                []OrderItemSnapshot    `json:"items"`
	SaleFeeAmount        *float64               `json:"sale_fee_amount,omitempty"`
	Payments             []OrderPaymentSnapshot `json:"payments"`
	ShippingID           string                 `json:"shipping_id,omitempty"`
	CancellationDetail   string                 `json:"cancellation_detail,omitempty"`
	Tags                 []string               `json:"tags,omitempty"`
}

type OrderItemSnapshot struct {
	ProviderItemID      string   `json:"provider_item_id"`
	ProviderVariationID string   `json:"provider_variation_id,omitempty"`
	SellerSKU           string   `json:"seller_sku,omitempty"`
	EAN                 string   `json:"ean,omitempty"`
	Title               string   `json:"title,omitempty"`
	Quantity            int      `json:"quantity"`
	UnitPrice           *float64 `json:"unit_price,omitempty"`
	SaleFeeAmount       *float64 `json:"sale_fee_amount,omitempty"`
}

type OrderPaymentSnapshot struct {
	PaymentID         string   `json:"payment_id"`
	Status            string   `json:"status,omitempty"`
	Amount            *float64 `json:"amount,omitempty"`
	TransactionAmount *float64 `json:"transaction_amount,omitempty"`
	TotalPaidAmount   *float64 `json:"total_paid_amount,omitempty"`
}

type FeeQuoteInput struct {
	AccountRef    ProviderAccountRef
	SiteID        string
	CategoryID    string
	ListingTypeID string
	PriceAmount   float64
	CurrencyID    string
}

type FeeQuoteSnapshot struct {
	ProviderCode      string     `json:"provider_code"`
	SiteID            string     `json:"site_id"`
	CategoryID        string     `json:"category_id"`
	ListingTypeID     string     `json:"listing_type_id"`
	PriceAmount       float64    `json:"price_amount"`
	CurrencyID        string     `json:"currency_id"`
	CommissionPercent *float64   `json:"commission_percent,omitempty"`
	FixedFeeAmount    *float64   `json:"fixed_fee_amount,omitempty"`
	SourceUpdatedAt   *time.Time `json:"source_updated_at,omitempty"`
	FetchedAt         time.Time  `json:"fetched_at"`
	RawProviderRef    any        `json:"raw_provider_ref,omitempty"`
}
