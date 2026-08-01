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

// mlIngestShipmentResponse is the GET /shipments/{id} decode target for the M-03
// order_shipments ingest reader (IC-03, migration 0088). It reuses the already-verified
// private types from shipping_reader.go (flexString for the numeric-or-string id,
// mlDestination/mlShippingAddress/mlNamedRef for the destination block) rather than
// redefining them — same package, same confirmed shape, no drift risk. Named distinctly
// (mlIngest*) to avoid colliding with mlShipmentResponse (shipping_reader.go), which serves a
// DIFFERENT, still-LIVE call site (own-cost simulation) and must not be touched here.
type mlIngestShipmentResponse struct {
	ID     flexString `json:"id"`
	SiteID string     `json:"site_id"`
	Status string     `json:"status"`
	// Substatus/DateCreated/TrackingNumber/Destination all appear verbatim in an in-repo
	// captured fixture (shipping_reader_test.go TestGetShipmentInfoDecodesNewFormatNumericID
	// carries a real `date_created` value with numeric id 28264263908) or in
	// docs/design/handoff/API-MAP.md's doc-reviewed field list (`tracking_number` ✅ direct
	// API on GET /shipments/$ID) — both count as "verified against fixture/docs" per
	// feature.md, neither is a guess.
	Substatus      *string        `json:"substatus"`
	DateCreated    string         `json:"date_created"`
	TrackingNumber *string        `json:"tracking_number"`
	Destination    *mlDestination `json:"destination"`
	// LeadTime feeds ShipmentDetail.SLALimitAt (0088 sla_limit_at) — NOT from a separate
	// /sla call. Reuses mlShipmentLeadTime/mlEstimatedDeliveryLimit (shipping_reader.go:77-83),
	// the same types the existing, still-LIVE getShipmentInfo/mapShipmentInfo already decode
	// this exact field with (shipping_reader.go:19,236-237) and docs/design/handoff-2026-07/
	// API-MAP.md:18 documents as coming from GET /shipments/$ID's estimated_delivery_* data.
	LeadTime *mlShipmentLeadTime `json:"lead_time"`
}

// GetShipmentDetail resolves the full order_shipments persistence fact for one shipment
// (IC-03, migration 0088): primary GET /shipments/{id} (x-format-new:true, same requirement
// as the existing getShipmentInfo) plus GET /shipments/{id}/costs (same header) for
// cost_gross/cost_seller/currency. A 404 or 410 on the PRIMARY call is honest-absence —
// domain.ShipmentDetail{} (Found stays false), never an error (feature.md Negative
// Scenarios: "While ML responde 404/410 p/ shipment ... the retorno shall ser
// honest-absence"). A costs failure degrades the same way the existing getShipmentInfo
// degrades: 404/undecodable leaves cost fields nil but keeps the rest of the shipment;
// auth/transient/rate-limit/validation propagate as real errors, never silently zeroed.
//
// This reader intentionally does NOT call the separate GET /shipments/{id}/sla sub-resource
// IC-03 names as SLAStatus's source: that endpoint's existence is live-verified
// (research/external-ml-api-facts.md fact #15) but its response shape is not confirmed in
// this worktree (see domain.ShipmentDetail's SLAStatus doc comment) — calling it and mapping
// nothing would just spend a rate-limit token for zero benefit. SLAStatus stays nil; this is
// a tracked honest gap (validation.md), not a silent omission.
//
// SLALimitAt, however, does NOT come from that /sla sub-resource (IC-03 is wrong on this
// specific point — docs/ outranks .mnfs/research/ per this repo's truth order): it comes from
// `lead_time.estimated_delivery_limit.date` on this SAME primary GET /shipments/{id} call,
// confirmed live by the existing mapShipmentInfo path (shipping_reader.go:236-237) and
// docs/design/handoff-2026-07/API-MAP.md:18. See mapShipmentDetail below.
func (a *CapabilityAdapter) GetShipmentDetail(ctx context.Context, accountRef domain.ProviderAccountRef, shipmentID string) (domain.ShipmentDetail, error) {
	accountRef, err := normalizeAccountRef(accountRef)
	if err != nil {
		return domain.ShipmentDetail{}, err
	}
	shipmentID = strings.TrimSpace(shipmentID)
	if shipmentID == "" {
		return domain.ShipmentDetail{}, domain.NewCapabilityError(domain.ErrCodeProviderInvalidReference, "provider shipment id is empty")
	}
	token, err := a.accessToken(ctx, accountRef)
	if err != nil {
		return domain.ShipmentDetail{}, mapPricingReaderError(err)
	}
	return a.getShipmentDetail(ctx, accountRef, token, shipmentID)
}

func (a *CapabilityAdapter) getShipmentDetail(ctx context.Context, accountRef domain.ProviderAccountRef, token, shipmentID string) (domain.ShipmentDetail, error) {
	var shipment mlIngestShipmentResponse
	path := "/shipments/" + url.PathEscape(shipmentID)
	if err := a.doJSONAllowGone(ctx, accountRef, token, http.MethodGet, path, map[string]string{"x-format-new": "true"}, &shipment); err != nil {
		if domain.ErrorCodeOf(err) == domain.ErrCodeProviderInvalidReference {
			// 404 or 410 (doJSONAllowGone maps both to the same code): honest-absence,
			// never an error.
			return domain.ShipmentDetail{}, nil
		}
		return domain.ShipmentDetail{}, mapPricingReaderError(err)
	}

	result := mapShipmentDetail(shipment, a.now())
	result.Found = true

	var costsResponse mlShipmentCostsResponse
	costsPath := path + "/costs"
	if err := a.doJSONWithHeaders(ctx, accountRef, token, http.MethodGet, costsPath, nil, map[string]string{"x-format-new": "true"}, &costsResponse); err != nil {
		if code := domain.ErrorCodeOf(err); code == domain.ErrCodeProviderInvalidReference || code == domain.ErrCodeProviderPayloadInvalid {
			return result, nil
		}
		return domain.ShipmentDetail{}, mapPricingReaderError(err)
	}
	applyShipmentCosts(&result, costsResponse, mlSiteCurrency(firstNonEmpty(shipment.SiteID, a.siteID)))
	return result, nil
}

// doJSONAllowGone is doJSONWithHeaders (capability_adapter.go:667) plus a 410 branch: ML
// documents 410 Gone as an additional "resource no longer exists" status for shipments, and
// IC-03/feature.md require 404 AND 410 to both degrade to honest-absence. Kept local (not
// added to capability_adapter.go, which is closed/M-01-owned) since it is shipment-ingest
// specific — every other call site in this package only needs the 404 doJSONWithHeaders
// already handles.
func (a *CapabilityAdapter) doJSONAllowGone(ctx context.Context, accountRef domain.ProviderAccountRef, token, method, path string, headers map[string]string, target any) error {
	resp, rawBody, err := a.doRawWithHeaders(ctx, accountRef, token, method, path, nil, "", headers)
	if err != nil {
		return err
	}
	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return domain.NewCapabilityError(domain.ErrCodeProviderAuth, "provider credentials were rejected")
	}
	if resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusGone {
		return domain.NewCapabilityError(domain.ErrCodeProviderInvalidReference, "provider resource was not found")
	}
	if resp.StatusCode >= 500 {
		return domain.NewCapabilityError(domain.ErrCodeProviderTransient, "provider temporarily unavailable")
	}
	if resp.StatusCode >= 400 {
		return domain.NewCapabilityError(domain.ErrCodeProviderValidation, providerDiag(method, path, resp.StatusCode, rawBody))
	}
	if target == nil {
		return nil
	}
	if err := json.Unmarshal(rawBody, target); err != nil {
		return domain.NewCapabilityError(domain.ErrCodeProviderPayloadInvalid, providerDiag(method, path, resp.StatusCode, rawBody))
	}
	return nil
}

func mapShipmentDetail(shipment mlIngestShipmentResponse, fetchedAt time.Time) domain.ShipmentDetail {
	detail := domain.ShipmentDetail{
		ProviderShipmentID: strings.TrimSpace(string(shipment.ID)),
		Status:             strings.TrimSpace(shipment.Status),
		Substatus:          trimmedPtr(shipment.Substatus),
		TrackingNumber:     trimmedPtr(shipment.TrackingNumber),
		SourceTime:         parseTimePtr(shipment.DateCreated),
		FetchedAt:          fetchedAt,
	}
	// SLALimitAt: lead_time.estimated_delivery_limit.date on this SAME primary payload —
	// same field, same nil-guarded traversal, as the existing mapShipmentInfo
	// (shipping_reader.go:236-237). Absent at any level (no lead_time / no
	// estimated_delivery_limit / no date) stays nil, never fabricated.
	if shipment.LeadTime != nil && shipment.LeadTime.EstimatedDeliveryLimit != nil && shipment.LeadTime.EstimatedDeliveryLimit.Date != nil {
		detail.SLALimitAt = parseTimePtr(*shipment.LeadTime.EstimatedDeliveryLimit.Date)
	}
	if shipment.Destination == nil {
		return detail
	}
	detail.ReceiverName = trimmedPtr(shipment.Destination.ReceiverName)
	addr := shipment.Destination.ShippingAddress
	if addr == nil {
		return detail
	}
	if addr.State != nil && addr.State.ID != nil {
		// state.id is the ML state code, e.g. "BR-SP" -> bare UF "SP" (same trim as the
		// existing mapShipmentInfo, shipping_reader.go).
		value := strings.TrimPrefix(strings.TrimSpace(*addr.State.ID), "BR-")
		if value != "" {
			detail.DestState = &value
		}
	}
	if addr.City != nil {
		detail.DestCity = trimmedPtr(addr.City.Name)
	}
	detail.DestZip = trimmedPtr(addr.ZipCode)
	return detail
}

// applyShipmentCosts fills CostGross/CostSeller/Currency from the /costs response.
// CostSeller comes from `senders[0].cost` — the existing ShipmentCosts.SenderCost
// (shipping_reader.go) documents this as the cost charged to the seller, matching IC-03's
// `cost_seller`/decomposition `frete_seller`. Currency is set only when at least one amount
// was decodable — never a fabricated currency label with no associated value (ADR-17).
func applyShipmentCosts(detail *domain.ShipmentDetail, resp mlShipmentCostsResponse, currency string) {
	detail.CostGross = decimalString(resp.GrossAmount)
	if len(resp.Senders) > 0 {
		detail.CostSeller = decimalString(resp.Senders[0].Cost)
	}
	if detail.CostGross != nil || detail.CostSeller != nil {
		c := currency
		detail.Currency = &c
	}
}
