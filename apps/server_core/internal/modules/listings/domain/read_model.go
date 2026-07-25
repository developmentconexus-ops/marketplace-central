package domain

import (
	"errors"
	"strings"
	"time"
)

type LinkState string

const (
	LinkStateUnresolved LinkState = "unresolved"
	LinkStateConflict   LinkState = "conflict"
	LinkStateResolved   LinkState = "resolved"
	LinkStateRejected   LinkState = "rejected"
)

type ListingException string

const (
	ListingExceptionSyncError    ListingException = "sync_error"
	ListingExceptionStale        ListingException = "stale"
	ListingExceptionUnlinked     ListingException = "unlinked"
	ListingExceptionBelowMargin  ListingException = "below_margin"
	ListingExceptionSemVinculo   ListingException = "sem_vinculo"
	ListingExceptionAbaixoCusto  ListingException = "abaixo_custo"
	ListingExceptionSemEvidencia ListingException = "sem_evidencia"
)

type ListingFilter struct {
	Status          ListingStatus
	SyncState       ListingSyncState
	LinkState       LinkState
	Exception       ListingException
	HasException    *bool
	ListingTypeCode ListingTypeCode
	ProductID       string
}

type ListingID struct {
	InstallationID    string
	ProviderListingID string
	VariationID       string
}

var ErrListingNotFound = errors.New("listing_not_found")

type ListingNotFoundError struct{}

func (*ListingNotFoundError) Error() string { return ErrListingNotFound.Error() }
func (*ListingNotFoundError) Unwrap() error { return ErrListingNotFound }

func ParseListingID(value string) (ListingID, error) {
	if value == "" || value != strings.TrimSpace(value) {
		return ListingID{}, &ListingNotFoundError{}
	}
	parts := strings.Split(value, "~")
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
		return ListingID{}, &ListingNotFoundError{}
	}
	return ListingID{InstallationID: parts[0], ProviderListingID: parts[1], VariationID: parts[2]}, nil
}

func (id ListingID) String() string {
	return id.InstallationID + "~" + id.ProviderListingID + "~" + id.VariationID
}

type ListingType struct {
	Code  ListingTypeCode `json:"code"`
	Label string          `json:"label"`
}

var listingTypeLabels = map[ListingTypeCode]string{
	"gold_pro": "Premium", "gold_special": "Clássico", "gold_premium": "Ouro Premium (legado)",
	"gold": "Ouro (legado)", "silver": "Prata (legado)", "bronze": "Bronze (legado)", "free": "Grátis",
}

func ListingTypeForCode(code ListingTypeCode) (ListingType, bool) {
	label, ok := listingTypeLabels[code]
	return ListingType{Code: code, Label: label}, ok
}

type Money struct {
	Amount   string        `json:"amount"`
	Currency PriceCurrency `json:"currency"`
}

type ListingLink struct {
	State     LinkState `json:"state"`
	ProductID *string   `json:"product_id"`
	SellerSKU *string   `json:"seller_sku"`
}

type ReadSyncError struct {
	Code            string  `json:"code"`
	MessagePT       string  `json:"message_pt"`
	MessageProvider *string `json:"message_provider"`
}

type PendingIssueKind string

const (
	PendingIssueSyncError         PendingIssueKind = "sync_error"
	PendingIssueStale             PendingIssueKind = "stale"
	PendingIssueUnlinked          PendingIssueKind = "unlinked"
	PendingIssueBelowMargin       PendingIssueKind = "below_margin"
	PendingIssueAttributeRequired PendingIssueKind = "attribute_required"
)

type PendingIssue struct {
	Kind      PendingIssueKind `json:"kind"`
	MessagePT string           `json:"message_pt"`
}

type ICMWorstCaseByUF struct {
	DestinationUF    string  `json:"destination_uf"`
	WorstCaseICMSPct *string `json:"worst_case_icms_pct"`
	PriceNetBasis    *string `json:"price_net_basis"`
	BelowMarginAtUF  *bool   `json:"below_margin_at_uf"`
}

type ListingReadModel struct {
	ListingID            string              `json:"listing_id"`
	InstallationID       string              `json:"installation_id"`
	Provider             string              `json:"provider"`
	ProviderListingID    string              `json:"provider_listing_id"`
	// A provider listing with variations produces one row per variation, all
	// carrying the same provider_listing_id. Without the variation the screen
	// shows N identical-looking rows the operator cannot tell apart. Absent
	// (NoVariationID) stays null instead of a fabricated "-".
	VariationID          *string             `json:"variation_id"`
	Title                string              `json:"title"`
	ListingType          *ListingType        `json:"listing_type"`
	Status               ListingStatus       `json:"status"`
	Link                 ListingLink         `json:"link"`
	Price                *Money              `json:"price"`
	PublishedQuantity    *int                `json:"published_quantity"`
	SyncState            ListingSyncState    `json:"sync_state"`
	SyncError            *ReadSyncError      `json:"sync_error"`
	QualityScore         *int                `json:"quality_score"`
	PendingIssue         *PendingIssue       `json:"pending_issue"`
	Sales30D             *int                `json:"sales_30d"`
	Cost                 *Money              `json:"cost"`
	BelowMarginWorstCase *bool               `json:"below_margin_worst_case"`
	ICMSWorstCaseByUF    *[]ICMWorstCaseByUF `json:"icms_worst_case_by_uf"`
	FetchedAt            *time.Time          `json:"fetched_at"`
	MarketSignal         *MarketSignal       `json:"market_signal"`
	SignalStatus         SignalStatus        `json:"signal_status"`
}

type GroupState string

const (
	GroupStateOK        GroupState = "ok"
	GroupStateAttention GroupState = "attention"
	GroupStateError     GroupState = "error"
)

type ListingGroup struct {
	ProductID    *string            `json:"product_id"`
	ProductTitle *string            `json:"product_title"`
	ListingCount int                `json:"listing_count"`
	GroupState   GroupState         `json:"group_state"`
	Listings     []ListingReadModel `json:"listings"`
}

type TimelineEvent struct {
	At        time.Time `json:"at"`
	Kind      string    `json:"kind"`
	MessagePT string    `json:"message_pt"`
}
