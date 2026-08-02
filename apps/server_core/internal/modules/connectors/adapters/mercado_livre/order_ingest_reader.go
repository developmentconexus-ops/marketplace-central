package mercadolivre

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"

	"marketplace-central/apps/server_core/internal/modules/connectors/domain"
)

// mlIngestOrderResponse is the FULL GET /orders/{id} decode target for the M-03 ingest reader
// — a richer sibling of mlOrderResponse (capability_adapter.go:1092, owned by M-01/closed).
// It adds `pack_id` (verified live, research/external-ml-api-facts.md fact #17,
// "pack_id agrupa pedidos (carrinho)" — top-level, nullable). Named distinctly (mlIngest*) to
// avoid redeclaring the existing mlOrder*/mlPayment private types in this same package.
//
// PackID is decoded as `any` (not *string): like order.id and shipping.id, ML emits this as a
// bare JSON number, not a string — a *string target would fail json.Unmarshal outright on any
// real payload (caught by TestGetOrderDetailDecodesAllTargetFields against a realistic
// fixture). Normalized to string via normalizeAnyID in mapOrderDetail, same as every other
// provider id in this file.
type mlIngestOrderResponse struct {
	ID              any                   `json:"id"`
	Status          string                `json:"status"`
	StatusDetail    string                `json:"status_detail"`
	DateCreated     string                `json:"date_created"`
	DateClosed      string                `json:"date_closed"`
	LastUpdated     string                `json:"last_updated"`
	DateLastUpdated string                `json:"date_last_updated"`
	PackID          any                   `json:"pack_id"`
	OrderItems      []mlIngestOrderItem   `json:"order_items"`
	Payments        []mlIngestPayment     `json:"payments"`
	Shipping        mlIngestOrderShipping `json:"shipping"`
	Tags            []string              `json:"tags"`
	Buyer           *mlIngestOrderBuyer   `json:"buyer"`
	// CancelDetail is the REAL source of the cancellation reason. status_detail came back
	// null in 7/7 measured cancelled orders while cancel_detail was present — it is not
	// duplicate data, it is the data that was missing (rawkeys/P2). Only requested_by and
	// code are ever consumed (see cancellationDetail below); description is ML's own UI
	// copy and changes without notice, so it is never read into our value (ADR-C4).
	CancelDetail *mlIngestOrderCancelDetail `json:"cancel_detail"`
	// CurrencyID feeds orders.currency (Task 5 connects the domain/repo column — this DTO
	// change only stops the value from being silently discarded by encoding/json).
	CurrencyID string `json:"currency_id"`
}

type mlIngestOrderCancelDetail struct {
	Group       string `json:"group"`
	Code        string `json:"code"`
	Description string `json:"description"`
	RequestedBy string `json:"requested_by"`
}

type mlIngestOrderShipping struct {
	ID any `json:"id"`
}

type mlIngestOrderBuyer struct {
	Nickname string `json:"nickname"`
}

type mlIngestOrderItem struct {
	Item      mlIngestOrderItemDetails `json:"item"`
	Quantity  int                      `json:"quantity"`
	UnitPrice *float64                 `json:"unit_price"`
	// SaleFee is the provider's `sale_fee` — PER UNIT (fact T2, live-verified 2026-07-20;
	// the doc does not declare this, it was established by measurement). Mapped onto
	// domain.OrderDetailItem.SaleFeeUnit — never renamed/summed into a total here.
	SaleFee *float64 `json:"sale_fee"`
}

type mlIngestOrderItemDetails struct {
	ID                  string        `json:"id"`
	VariationID         any           `json:"variation_id"`
	SellerSKU           string        `json:"seller_sku"`
	Title               string        `json:"title"`
	VariationAttributes []mlAttribute `json:"variation_attributes"`
}

type mlIngestPayment struct {
	ID                any      `json:"id"`
	Status            string   `json:"status"`
	TotalPaidAmount   *float64 `json:"total_paid_amount"`
	TransactionAmount *float64 `json:"transaction_amount"`
}

// GetOrderDetail resolves the full /orders/{id} payload for ingest (MIS-007/M-03 F-01,
// IC-03): the typed domain.OrderDetail (every column the 0089 extension + the existing
// baseline from import_service.go normalizeOrders need) plus a raw capture with
// buyer.billing_info stripped (ADR-03) — both from a SINGLE GET (EARS: "1 GET"), via the
// shared M-01 resilience choke point (doRawWithHeaders, capability_adapter.go:751).
//
// A 403 (third-party order — a documented live fact, feature.md Negative Scenarios) maps to
// domain.ErrUnauthorized via mapPricingReaderError: a typed, non-retryable error. It is never
// silently degraded — the caller must mark the ingest attempt and move on.
func (a *CapabilityAdapter) GetOrderDetail(ctx context.Context, accountRef domain.ProviderAccountRef, providerOrderID string) (domain.OrderDetail, error) {
	accountRef, err := normalizeAccountRef(accountRef)
	if err != nil {
		return domain.OrderDetail{}, err
	}
	providerOrderID = strings.TrimSpace(providerOrderID)
	if providerOrderID == "" {
		return domain.OrderDetail{}, domain.NewCapabilityError(domain.ErrCodeProviderInvalidReference, "provider order id is empty")
	}
	token, err := a.accessToken(ctx, accountRef)
	if err != nil {
		return domain.OrderDetail{}, mapPricingReaderError(err)
	}

	var order mlIngestOrderResponse
	path := "/orders/" + url.PathEscape(providerOrderID)
	rawBody, err := a.doJSONRaw(ctx, accountRef, token, http.MethodGet, path, &order)
	if err != nil {
		return domain.OrderDetail{}, mapPricingReaderError(err)
	}

	rawOrder, err := redactBillingInfo(rawBody)
	if err != nil {
		// Cannot safely guarantee the billing_info redaction succeeded — drop the raw
		// capture entirely rather than risk ADR-03 leaking a PII/fiscal fragment. The
		// typed DTO below is the load-bearing data; raw is a supplementary capture only.
		rawOrder = nil
	}

	return mapOrderDetail(order, rawOrder, a.now()), nil
}

// doJSONRaw performs a single GET through the shared resilience choke point (doRaw ->
// doRawWithHeaders -> doRawWithIdempotency -> doRawOnce, capability_adapter.go:732 onward)
// and returns BOTH the decoded target and the raw response body, so one provider round trip
// can serve a typed DTO and a raw-payload capture. Status-code mapping mirrors
// doJSONWithHeaders (capability_adapter.go:667) exactly; duplicated here (not added there)
// because that file is closed (M-01) and doJSONWithHeaders does not expose raw bytes to its
// caller.
func (a *CapabilityAdapter) doJSONRaw(ctx context.Context, accountRef domain.ProviderAccountRef, token, method, path string, target any) ([]byte, error) {
	resp, rawBody, err := a.doRaw(ctx, accountRef, token, method, path, nil)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, domain.NewCapabilityError(domain.ErrCodeProviderAuth, "provider credentials were rejected")
	}
	if resp.StatusCode == http.StatusNotFound {
		return nil, domain.NewCapabilityError(domain.ErrCodeProviderInvalidReference, "provider resource was not found")
	}
	if resp.StatusCode >= 500 {
		return nil, domain.NewCapabilityError(domain.ErrCodeProviderTransient, "provider temporarily unavailable")
	}
	if resp.StatusCode >= 400 {
		return nil, domain.NewCapabilityError(domain.ErrCodeProviderValidation, providerDiag(method, path, resp.StatusCode, rawBody))
	}
	if target != nil {
		if err := json.Unmarshal(rawBody, target); err != nil {
			return nil, domain.NewCapabilityError(domain.ErrCodeProviderPayloadInvalid, providerDiag(method, path, resp.StatusCode, rawBody))
		}
	}
	return rawBody, nil
}

// redactBillingInfo returns the order payload with `buyer.billing_info` removed (ADR-03: raw
// order capture may be held in memory MINUS billing_info). It operates on
// map[string]json.RawMessage at the top level and inside `buyer` only, so every OTHER field
// (order_items, payments, large numeric ids, …) is re-emitted byte-for-byte from its original
// RawMessage — no float64 round-trip, no precision loss. If `buyer` cannot be safely
// re-decoded as an object, it is dropped entirely rather than risked: the invariant this
// function exists to guarantee (no `billing_info` substring survives) is non-negotiable.
func redactBillingInfo(raw []byte) (json.RawMessage, error) {
	var top map[string]json.RawMessage
	if err := json.Unmarshal(raw, &top); err != nil {
		return nil, err
	}
	if buyerRaw, ok := top["buyer"]; ok && len(buyerRaw) > 0 {
		var buyer map[string]json.RawMessage
		if err := json.Unmarshal(buyerRaw, &buyer); err != nil {
			delete(top, "buyer")
		} else {
			delete(buyer, "billing_info")
			if scrubbed, err := json.Marshal(buyer); err == nil {
				top["buyer"] = scrubbed
			} else {
				delete(top, "buyer")
			}
		}
	}
	return json.Marshal(top)
}

// cancellationDetail derives the canonical "<requested_by>:<code>" form from the order's
// cancel_detail block (e.g. "buyer:cancel_purchase_by_buyer") — requested_by answers "who
// cancelled" (buyer/meli/seller), which is the operational question. Returns "" when the
// order carries no cancel_detail (the common case: cancel_detail only exists on cancelled
// orders) or when both fields are blank; never falls back to status_detail — that fallback
// was the P2 bug, not a legitimate substitute.
func cancellationDetail(cancel *mlIngestOrderCancelDetail) string {
	if cancel == nil {
		return ""
	}
	requestedBy := strings.TrimSpace(cancel.RequestedBy)
	code := strings.TrimSpace(cancel.Code)
	if requestedBy == "" && code == "" {
		return ""
	}
	return requestedBy + ":" + code
}

func mapOrderDetail(order mlIngestOrderResponse, rawOrder json.RawMessage, fetchedAt time.Time) domain.OrderDetail {
	detail := domain.OrderDetail{
		ProviderCode:         "mercado_livre",
		ProviderOrderID:      normalizeAnyID(order.ID),
		ProviderStatus:       strings.TrimSpace(order.Status),
		ProviderStatusDetail: strings.TrimSpace(order.StatusDetail),
		ProviderCreatedAt:    parseTimePtr(order.DateCreated),
		ProviderClosedAt:     parseTimePtr(order.DateClosed),
		ProviderUpdatedAt:    parseTimePtr(firstNonEmpty(order.LastUpdated, order.DateLastUpdated)),
		DateLastUpdatedML:    parseTimePtr(order.DateLastUpdated),
		FetchedAt:            fetchedAt,
		CancellationDetail:   cancellationDetail(order.CancelDetail),
		Tags:                 trimStrings(order.Tags),
		RawOrder:             rawOrder,
	}
	if packID := normalizeAnyID(order.PackID); packID != "" {
		detail.PackID = &packID
	}
	if shippingID := normalizeAnyID(order.Shipping.ID); shippingID != "" {
		detail.ShippingID = &shippingID
	}
	if order.Buyer != nil {
		detail.BuyerNickname = trimmedStringPtr(order.Buyer.Nickname)
	}
	for _, item := range order.OrderItems {
		detail.Items = append(detail.Items, domain.OrderDetailItem{
			ProviderItemID:      strings.TrimSpace(item.Item.ID),
			ProviderVariationID: normalizeAnyID(item.Item.VariationID),
			SellerSKU:           strings.TrimSpace(item.Item.SellerSKU),
			EAN:                 attributeValue(item.Item.VariationAttributes, "GTIN", "EAN"),
			Title:               strings.TrimSpace(item.Item.Title),
			Quantity:            item.Quantity,
			UnitPrice:           item.UnitPrice,
			SaleFeeUnit:         item.SaleFee,
		})
	}
	for _, payment := range order.Payments {
		detail.Payments = append(detail.Payments, domain.OrderDetailPayment{
			PaymentID:         normalizeAnyID(payment.ID),
			Status:            strings.TrimSpace(payment.Status),
			TransactionAmount: payment.TransactionAmount,
			TotalPaidAmount:   payment.TotalPaidAmount,
		})
	}
	return detail
}
