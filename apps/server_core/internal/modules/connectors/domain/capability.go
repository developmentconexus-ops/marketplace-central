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
	ProviderCode        string
	ProviderItemID      string
	ProviderVariationID string
	ProviderStatus      string
	SellerSKU           string
	EAN                 string
	Title               string
	AvailableQuantity   *int
	SourceUpdatedAt     *time.Time
	FetchedAt           time.Time
	RawProviderRef      any
	Variations          []ListingVariationSnapshot
}

type ListingVariationSnapshot struct {
	ProviderVariationID string
	SellerSKU           string
	EAN                 string
	AvailableQuantity   *int
}

type StockSnapshot struct {
	ProviderCode        string
	ProviderItemID      string
	ProviderVariationID string
	AvailableQuantity   int
	ProviderStatus      string
	SellerSKU           string
	EAN                 string
	Title               string
	SourceUpdatedAt     *time.Time
	FetchedAt           time.Time
	RawProviderRef      any
	Scope               StockScope
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
	IDempotencyKey      string
	Result              StockWriteResultStatus
	ProviderCode        string
	ProviderItemID      string
	ProviderVariationID string
	Message             string
	ProviderResponseRef any
}

type OrderSnapshot struct {
	ProviderCode       string
	ProviderOrderID    string
	ProviderStatus     string
	SourceUpdatedAt    *time.Time
	FetchedAt          time.Time
	RawProviderRef     any
	Items              []OrderItemSnapshot
	SaleFeeAmount      *float64
	Payments           []OrderPaymentSnapshot
	ShippingID         string
	CancellationDetail string
}

type OrderItemSnapshot struct {
	ProviderItemID      string
	ProviderVariationID string
	SellerSKU           string
	EAN                 string
	Title               string
	Quantity            int
	UnitPrice           *float64
}

type OrderPaymentSnapshot struct {
	PaymentID string
	Status    string
	Amount    *float64
}
