package domain

import (
	"errors"
	"strings"
	"time"
)

const NoVariationID = "-"

type ListingStatus string

const (
	ListingStatusActive          ListingStatus = "active"
	ListingStatusPaused          ListingStatus = "paused"
	ListingStatusClosed          ListingStatus = "closed"
	ListingStatusUnknown         ListingStatus = "unknown"
	ListingStatusUnderReview     ListingStatus = "under_review"
	ListingStatusInactive        ListingStatus = "inactive"
	ListingStatusPaymentRequired ListingStatus = "payment_required"
	ListingStatusNotYetActive    ListingStatus = "not_yet_active"
)

type ListingSyncState string

const (
	ListingSyncStateSynced     ListingSyncState = "synced"
	ListingSyncStateError      ListingSyncState = "error"
	ListingSyncStateStale      ListingSyncState = "stale"
	ListingSyncStateQueued     ListingSyncState = "queued"
	ListingSyncStateSyncing    ListingSyncState = "syncing"
	ListingSyncStatePausedSync ListingSyncState = "paused_sync"
)

type ListingTypeCode string
type PriceAmount string
type PriceCurrency string

type ListingSyncError struct {
	Code            string
	MessagePT       string
	MessageProvider string
}

type ListingKey struct {
	TenantID          string
	InstallationID    string
	ProviderListingID string
	VariationID       string
}

var (
	ErrTenantIDRequired          = errors.New("tenant_id is required")
	ErrInstallationIDRequired    = errors.New("installation_id is required")
	ErrProviderRequired          = errors.New("provider is required")
	ErrProviderListingIDRequired = errors.New("provider_listing_id is required")
	ErrTitleRequired             = errors.New("title is required")
	ErrListingStatusInvalid      = errors.New("listing status is invalid")
	ErrListingSyncStateInvalid   = errors.New("listing sync state is invalid")
)

func NewListingKey(tenantID, installationID, providerListingID, variationID string) (ListingKey, error) {
	tenantID = strings.TrimSpace(tenantID)
	installationID = strings.TrimSpace(installationID)
	providerListingID = strings.TrimSpace(providerListingID)
	if tenantID == "" {
		return ListingKey{}, ErrTenantIDRequired
	}
	if installationID == "" {
		return ListingKey{}, ErrInstallationIDRequired
	}
	if providerListingID == "" {
		return ListingKey{}, ErrProviderListingIDRequired
	}

	variationID = strings.TrimSpace(variationID)
	if variationID == "" {
		variationID = NoVariationID
	}
	return ListingKey{
		TenantID:          tenantID,
		InstallationID:    installationID,
		ProviderListingID: providerListingID,
		VariationID:       variationID,
	}, nil
}

type Listing struct {
	Key               ListingKey
	Provider          string
	Title             string
	ListingTypeCode   *ListingTypeCode
	Status            ListingStatus
	PriceAmount       *PriceAmount
	PriceCurrency     *PriceCurrency
	PublishedQuantity *int
	SyncState         ListingSyncState
	SyncError         *ListingSyncError
	QualityScore      *int
	Sales30D          *int
	FetchedAt         *time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time

	// E3 fields (0090, IC-07) — honest-unknown per ADR-17: nil means the
	// caller never populated the field from provider data, never a fabricated
	// zero/false/empty. Deliberately excludes commission_amount/commission_pct/
	// free_shipping_cost — those live in channel_fees (IC-01), not here.
	SoldQuantity      *int
	CategoryID        *string
	Condition         *string
	Permalink         *string
	Thumbnail         *string
	DateCreatedML     *time.Time
	Tags              []string
	CatalogProductID  *string
	ShippingMode      *string
	FreeShipping      *bool
	LogisticType      *string
	AvailableQuantity *int

	// Variations (0091) — optional child rows upserted alongside the parent
	// row, in the same transaction, by ApplyCompletedPull.
	Variations []ListingVariation
}

// ListingVariation is a child row of a listing (0091 listing_variations).
// Fields are pointer/nullable-friendly for the same honest-unknown reason as
// the Listing E3 fields above.
type ListingVariation struct {
	VariationID       string
	Price             *PriceAmount
	AvailableQuantity *int
	SoldQuantity      *int
	SellerSKU         *string
	Attributes        map[string]any
}

type ListingInput struct {
	TenantID          string
	InstallationID    string
	Provider          string
	ProviderListingID string
	VariationID       string
	Title             string
	ListingTypeCode   *ListingTypeCode
	Status            ListingStatus
	PriceAmount       *PriceAmount
	PriceCurrency     *PriceCurrency
	PublishedQuantity *int
	SyncState         ListingSyncState
	SyncError         *ListingSyncError
	QualityScore      *int
	Sales30D          *int
	FetchedAt         *time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time

	// E3 fields — see Listing above.
	SoldQuantity      *int
	CategoryID        *string
	Condition         *string
	Permalink         *string
	Thumbnail         *string
	DateCreatedML     *time.Time
	Tags              []string
	CatalogProductID  *string
	ShippingMode      *string
	FreeShipping      *bool
	LogisticType      *string
	AvailableQuantity *int
}

func NewListing(input ListingInput) (Listing, error) {
	key, err := NewListingKey(input.TenantID, input.InstallationID, input.ProviderListingID, input.VariationID)
	if err != nil {
		return Listing{}, err
	}
	provider := strings.TrimSpace(input.Provider)
	if provider == "" {
		return Listing{}, ErrProviderRequired
	}
	title := strings.TrimSpace(input.Title)
	if title == "" {
		return Listing{}, ErrTitleRequired
	}
	if !input.Status.IsValid() {
		return Listing{}, ErrListingStatusInvalid
	}
	if !input.SyncState.IsValid() {
		return Listing{}, ErrListingSyncStateInvalid
	}

	return Listing{
		Key:               key,
		Provider:          provider,
		Title:             title,
		ListingTypeCode:   input.ListingTypeCode,
		Status:            input.Status,
		PriceAmount:       input.PriceAmount,
		PriceCurrency:     input.PriceCurrency,
		PublishedQuantity: input.PublishedQuantity,
		SyncState:         input.SyncState,
		SyncError:         input.SyncError,
		QualityScore:      input.QualityScore,
		Sales30D:          input.Sales30D,
		FetchedAt:         input.FetchedAt,
		CreatedAt:         input.CreatedAt,
		UpdatedAt:         input.UpdatedAt,
		SoldQuantity:      input.SoldQuantity,
		CategoryID:        input.CategoryID,
		Condition:         input.Condition,
		Permalink:         input.Permalink,
		Thumbnail:         input.Thumbnail,
		DateCreatedML:     input.DateCreatedML,
		Tags:              input.Tags,
		CatalogProductID:  input.CatalogProductID,
		ShippingMode:      input.ShippingMode,
		FreeShipping:      input.FreeShipping,
		LogisticType:      input.LogisticType,
		AvailableQuantity: input.AvailableQuantity,
	}, nil
}

func (s ListingStatus) IsValid() bool {
	switch s {
	case ListingStatusActive, ListingStatusPaused, ListingStatusClosed, ListingStatusUnknown,
		ListingStatusUnderReview, ListingStatusInactive, ListingStatusPaymentRequired, ListingStatusNotYetActive:
		return true
	default:
		return false
	}
}

func (s ListingSyncState) IsValid() bool {
	switch s {
	case ListingSyncStateSynced, ListingSyncStateError, ListingSyncStateStale,
		ListingSyncStateQueued, ListingSyncStateSyncing, ListingSyncStatePausedSync:
		return true
	default:
		return false
	}
}
