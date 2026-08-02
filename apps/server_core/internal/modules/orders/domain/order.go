package domain

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"sort"
	"strings"
	"time"
)

var (
	ErrInvalidLineIdentity       = errors.New("invalid order line identity")
	ErrOrderLineIdentityConflict = errors.New("order line identity conflict")
)

type MPCLineID string

type LineReconciliationState string

const (
	LineReconciliationStable           LineReconciliationState = "stable"
	LineReconciliationLegacyUnresolved LineReconciliationState = "legacy_unresolved"
	LineReconciliationAmbiguous        LineReconciliationState = "ambiguous"
)

type LinkQuality string

const (
	LinkQualityResolved   LinkQuality = "resolved"
	LinkQualityRejected   LinkQuality = "rejected"
	LinkQualityConflict   LinkQuality = "conflict"
	LinkQualityUnresolved LinkQuality = "unresolved"
	LinkQualityMissing    LinkQuality = "missing"
)

type ListingIdentity struct {
	InstallationID      string `json:"installation_id"`
	ProviderItemID      string `json:"provider_item_id"`
	ProviderVariationID string `json:"provider_variation_id,omitempty"`
}

type MarketplaceOrder struct {
	InstallationID       string                    `json:"installation_id"`
	ProviderCode         string                    `json:"provider_code"`
	ProviderOrderID      string                    `json:"provider_order_id"`
	ProviderStatus       string                    `json:"provider_status,omitempty"`
	// ProviderStatusDetail e CancellationDetail são *string, não string: o
	// provider distingue "não mandou o campo" de "mandou vazio", e string
	// colapsa os dois em "". Foi esse colapso que fez cancellation_detail
	// aparecer preenchido em 38/38 pedidos — inclusive nos 7 cancelados, onde
	// o motivo real do cancelamento é justamente o que faltava.
	ProviderStatusDetail *string                   `json:"provider_status_detail,omitempty"`
	ProviderCreatedAt    *time.Time                `json:"provider_created_at,omitempty"`
	ProviderClosedAt     *time.Time                `json:"provider_closed_at,omitempty"`
	ProviderUpdatedAt    *time.Time                `json:"provider_updated_at,omitempty"`
	FetchedAt            time.Time                 `json:"fetched_at"`
	ShippingID           string                    `json:"shipping_id,omitempty"`
	CancellationDetail   *string                   `json:"cancellation_detail,omitempty"`
	Tags                 []string                  `json:"tags,omitempty"`
	BuyerNickname        string                    `json:"buyer_nickname,omitempty"`
	RawProviderRef       string                    `json:"raw_provider_ref,omitempty"`
	Items                []MarketplaceOrderItem    `json:"items"`
	Payments             []MarketplaceOrderPayment `json:"payments"`
	CreatedAt            time.Time                 `json:"created_at"`
	UpdatedAt            time.Time                 `json:"updated_at"`

	// PackID/ProviderShipmentID/Bucket/DateLastUpdatedML/Buyer* below are the migration 0089
	// extension (IC-03/F-02). ProviderShipmentID is a soft link to order_shipments (0088 carries
	// no FK); nil when the order has no shipping.id. Bucket is domain.DeriveOrderBucket's OUTPUT,
	// computed once at ingest time and persisted — never re-derived in SQL (order_repo.go
	// GetOrderBucketCounts already documents why the derivation stays in Go). DateLastUpdatedML
	// is the provider's own date_last_updated, verbatim — the IC-06 incremental-sync watermark,
	// never defaulted from ProviderUpdatedAt's last_updated fallback.
	//
	// The 9 Buyer* fields are the flat buyer-fiscal columns 0089 adds (BuyerFiscalReader's
	// nested BuyerFiscalInfo/BuyerFiscalAddress shape, flattened for storage). StateCode only —
	// 0089 deliberately excludes state_name/uf_nome and fetched_at (see migration comment).
	//
	// net_amount/margin_pct/decomposition (0089) are M-06 scope — deliberately absent from this
	// struct's write path; IngestOrder must never populate them.
	PackID             *string     `json:"pack_id,omitempty"`
	ProviderShipmentID *string     `json:"provider_shipment_id,omitempty"`
	Bucket             OrderBucket `json:"bucket,omitempty"`
	DateLastUpdatedML  *time.Time  `json:"date_last_updated_ml,omitempty"`

	BuyerName             *string `json:"buyer_name,omitempty"`
	BuyerDocType          *string `json:"buyer_doc_type,omitempty"`
	BuyerDocNumber        *string `json:"buyer_doc_number,omitempty"`
	BuyerAddressStreet    *string `json:"buyer_address_street,omitempty"`
	BuyerAddressNumber    *string `json:"buyer_address_number,omitempty"`
	BuyerAddressCity      *string `json:"buyer_address_city,omitempty"`
	BuyerAddressStateCode *string `json:"buyer_address_state_code,omitempty"`
	BuyerAddressZip       *string `json:"buyer_address_zip,omitempty"`
	BuyerAddressCountry   *string `json:"buyer_address_country,omitempty"`
}

type MarketplaceOrderItem struct {
	MPCLineID           MPCLineID               `json:"-"`
	ReconciliationState LineReconciliationState `json:"-"`
	ProviderItemID      string                  `json:"provider_item_id"`
	ProviderVariationID string                  `json:"provider_variation_id,omitempty"`
	SellerSKU           string                  `json:"seller_sku,omitempty"`
	Title               string                  `json:"title,omitempty"`
	Quantity            int                     `json:"quantity"`
	UnitPrice           *float64                `json:"unit_price,omitempty"`
	// SaleFeeAmount is PER-UNIT, verbatim from the provider's sale_fee (OrderDetailItem.SaleFeeUnit) —
	// never multiplied by Quantity. A consumer needing a line/order total must multiply by
	// Quantity itself.
	SaleFeeAmount     *float64    `json:"sale_fee_amount,omitempty"`
	LinkQuality       LinkQuality `json:"link_quality"`
	InternalProductID *int        `json:"internal_product_id,omitempty"`
}

type MarketplaceOrderPayment struct {
	ProviderPaymentID string   `json:"provider_payment_id"`
	ProviderStatus    string   `json:"provider_status,omitempty"`
	TransactionAmount *float64 `json:"transaction_amount,omitempty"`
	TotalPaidAmount   *float64 `json:"total_paid_amount,omitempty"`
}

type ImportResult struct {
	InstallationID string             `json:"installation_id"`
	ImportedCount  int                `json:"imported_count"`
	SkippedCount   int                `json:"skipped_count"`
	Items          []MarketplaceOrder `json:"items"`
}

type ItemLink struct {
	Identity          ListingIdentity
	Quality           LinkQuality
	InternalProductID *int
}

// NewMPCLineID creates an opaque, MPC-owned identity. Provider attributes and
// positional line numbers are deliberately absent from the value.
func NewMPCLineID() (MPCLineID, error) {
	var value [16]byte
	if _, err := rand.Read(value[:]); err != nil {
		return "", err
	}
	return MPCLineID("mpl_" + hex.EncodeToString(value[:])), nil
}

func (id MPCLineID) Valid() bool {
	value := string(id)
	if len(value) != 36 || !strings.HasPrefix(value, "mpl_") {
		return false
	}
	_, err := hex.DecodeString(value[4:])
	return err == nil
}

func (state LineReconciliationState) Valid() bool {
	switch state {
	case LineReconciliationStable, LineReconciliationLegacyUnresolved, LineReconciliationAmbiguous:
		return true
	default:
		return false
	}
}

type LineIDAllocator func() (MPCLineID, error)

type providerLineIdentity struct {
	itemID      string
	variationID string
	sellerSKU   string
}

// ReconcileOrderLineIdentities carries stable identities across a refresh.
// Only provider item, variation, and seller SKU identify a group. Quantity,
// price, and incoming position are mutable snapshot evidence. Duplicate groups
// retain the old ID multiset but remain ambiguous, so assigning that set in
// deterministic storage order does not claim an individual match.
func ReconcileOrderLineIdentities(existing, incoming []MarketplaceOrderItem, allocate LineIDAllocator) ([]MarketplaceOrderItem, error) {
	if allocate == nil {
		return nil, ErrInvalidLineIdentity
	}

	result := append([]MarketplaceOrderItem(nil), incoming...)
	existingGroups := make(map[providerLineIdentity][]MarketplaceOrderItem)
	for _, item := range existing {
		key := lineIdentityKey(item)
		existingGroups[key] = append(existingGroups[key], item)
	}
	incomingGroups := make(map[providerLineIdentity][]int)
	for index, item := range result {
		key := lineIdentityKey(item)
		incomingGroups[key] = append(incomingGroups[key], index)
	}

	for key, indexes := range incomingGroups {
		previous := append([]MarketplaceOrderItem(nil), existingGroups[key]...)
		sort.Slice(previous, func(left, right int) bool {
			return previous[left].MPCLineID < previous[right].MPCLineID
		})

		state := LineReconciliationStable
		if len(previous) > 1 || len(indexes) > 1 || containsReconciliationState(previous, LineReconciliationAmbiguous) {
			state = LineReconciliationAmbiguous
		} else if containsReconciliationState(previous, LineReconciliationLegacyUnresolved) {
			state = LineReconciliationLegacyUnresolved
		}

		for position, index := range indexes {
			var id MPCLineID
			if position < len(previous) && previous[position].MPCLineID.Valid() {
				id = previous[position].MPCLineID
			} else {
				var err error
				id, err = allocate()
				if err != nil || !id.Valid() {
					return nil, ErrInvalidLineIdentity
				}
			}
			result[index].MPCLineID = id
			result[index].ReconciliationState = state
		}
	}
	return result, nil
}

func lineIdentityKey(item MarketplaceOrderItem) providerLineIdentity {
	return providerLineIdentity{
		itemID:      item.ProviderItemID,
		variationID: item.ProviderVariationID,
		sellerSKU:   item.SellerSKU,
	}
}

func containsReconciliationState(items []MarketplaceOrderItem, state LineReconciliationState) bool {
	for _, item := range items {
		if item.ReconciliationState == state {
			return true
		}
	}
	return false
}
