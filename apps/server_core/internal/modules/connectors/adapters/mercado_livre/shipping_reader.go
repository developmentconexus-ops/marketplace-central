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

type mlShipmentResponse struct {
	ID          flexString          `json:"id"`
	SiteID      string              `json:"site_id"`
	Status      string              `json:"status"`
	Substatus   *string             `json:"substatus"`
	LeadTime    *mlShipmentLeadTime `json:"lead_time"`
	Delayed     *bool               `json:"delayed"`
	Destination *mlDestination      `json:"destination"`
}

// mlSiteCurrency resolves the ISO currency for an ML site_id. The new-format
// /shipments/{id}/costs payload carries NO currency field (ML docs:
// gerenciamento-de-envios — the cost body is gross_amount + receiver + senders
// only), so the amount currency is resolved from the shipment's site. Sites and
// their currencies are ML-documented (developers site list): MLB→BRL, MLA→ARS,
// MLM→MXN, MLC→CLP, MCO→COP, MLU→UYU, MPE→PEN. This BR connector defaults to BRL
// when site_id is absent/unknown rather than emitting an empty currency, which
// domain.ValidateMoney rejects and which would sink the read.
func mlSiteCurrency(siteID string) string {
	switch strings.ToUpper(strings.TrimSpace(siteID)) {
	case "MLA":
		return "ARS"
	case "MLM":
		return "MXN"
	case "MLC":
		return "CLP"
	case "MCO":
		return "COP"
	case "MLU":
		return "UYU"
	case "MPE":
		return "PEN"
	default:
		return "BRL"
	}
}

// flexString decodes a JSON value that ML may send as either a string or a bare
// number. The x-format-new shipment JSON carries a NUMERIC id (ML docs:
// status-de-pedidos-rastreamento), so a plain `string` field would fail to
// unmarshal and sink the whole read. Accepting both keeps decoding tolerant of
// that shape drift (ADR-17: shape drift degrades the field, never the read).
type flexString string

func (s *flexString) UnmarshalJSON(data []byte) error {
	trimmed := strings.TrimSpace(string(data))
	if trimmed == "null" || trimmed == "" {
		*s = ""
		return nil
	}
	if trimmed[0] == '"' {
		var str string
		if err := json.Unmarshal(data, &str); err != nil {
			return err
		}
		*s = flexString(str)
		return nil
	}
	// Bare number (or any other scalar token): keep its literal text.
	*s = flexString(trimmed)
	return nil
}

type mlShipmentLeadTime struct {
	EstimatedDeliveryLimit *mlEstimatedDeliveryLimit `json:"estimated_delivery_limit"`
}

type mlEstimatedDeliveryLimit struct {
	Date *string `json:"date"`
}

// mlDestination mirrors the documented new-format GET /shipments/{id}
// `destination` object (ML docs: gerenciamento-de-envios). Rounds 2–3's decode
// targeted a `receiver_address` field that does NOT exist in the new-format
// schema — the buyer delivery address lives under `destination.shipping_address`.
// Only the fields we actually surface are decoded (state, city name, zip,
// receiver_name); the many other documented sub-fields (agency, geolocation,
// scoring, municipality, …) are intentionally ignored (YAGNI).
type mlDestination struct {
	ReceiverName    *string            `json:"receiver_name"`
	ShippingAddress *mlShippingAddress `json:"shipping_address"`
}

type mlShippingAddress struct {
	State   *mlNamedRef `json:"state"`
	City    *mlNamedRef `json:"city"`
	ZipCode *string     `json:"zip_code"`
}

// mlNamedRef is the documented `{id, name}` shape ML uses for state and city.
type mlNamedRef struct {
	ID   *string `json:"id"`
	Name *string `json:"name"`
}

// mlShipmentCarrierResponse is the dedicated GET /shipments/{id}/carrier
// sub-resource (ML docs: gerenciamento-de-envios) — carrier name + tracking URL
// only. Fetched with x-format-new:true; a 404/decode failure degrades to a nil
// carrier (honest absence — ML returns 404 as a normal "none" condition on
// sibling sub-resources like /delays), never sinking the shipment read.
type mlShipmentCarrierResponse struct {
	URL  *string `json:"url"`
	Name *string `json:"name"`
}

type mlShipmentCostsResponse struct {
	GrossAmount *json.Number            `json:"gross_amount"`
	Receiver    *mlShipmentCostReceiver `json:"receiver"`
	Senders     []mlShipmentCostSender  `json:"senders"`
}

type mlShipmentCostReceiver struct {
	Cost *json.Number `json:"cost"`
}

type mlShipmentCostSender struct {
	Cost *json.Number `json:"cost"`
}

type mlFreeShippingResponse struct {
	Cost     *json.Number        `json:"cost"`
	Coverage *mlShippingCoverage `json:"coverage"`
	Currency string              `json:"currency_id"`
}

type mlShippingCoverage struct {
	AllCountry *mlShippingCoverageCountry `json:"all_country"`
}

type mlShippingCoverageCountry struct {
	ListCost *json.Number `json:"list_cost"`
}

func (a *CapabilityAdapter) getShipmentInfo(ctx context.Context, accountRef domain.ProviderAccountRef, token, shipmentID string) (domain.ShipmentInfo, error) {
	var shipment mlShipmentResponse
	path := "/shipments/" + url.PathEscape(shipmentID)
	// The PRIMARY shipment GET also REQUIRES `x-format-new: true`; without it ML
	// returns the legacy shipment shape and the decode sinks the whole read (ML
	// docs: gerenciamento-de-envios — "ao fazer uma requisição GET, é necessário
	// enviar o header 'x-format-new: true'"). A primary decode failure stays a
	// hard error (honest): unlike costs, there is no partial shipment to degrade to.
	if err := a.doJSONWithHeaders(ctx, accountRef, token, http.MethodGet, path, nil, map[string]string{"x-format-new": "true"}, &shipment); err != nil {
		return domain.ShipmentInfo{}, mapPricingReaderError(err)
	}

	// Status/substatus/destino are in hand — build the result now, and fill the
	// carrier sub-resource independent of costs (carrier survives a costs failure).
	result := mapShipmentInfo(shipment, a.now())
	a.fillShipmentCarrier(ctx, accountRef, token, path, &result)

	var costsResponse mlShipmentCostsResponse
	costsPath := path + "/costs"
	// The ML shipments /costs endpoint REQUIRES `x-format-new: true`; without it
	// ML returns the legacy cost shape which fails to decode into
	// mlShipmentCostsResponse (ML docs: mercado-envios-2 "Get Shipment Details
	// with New Format").
	if err := a.doJSONWithHeaders(ctx, accountRef, token, http.MethodGet, costsPath, nil, map[string]string{"x-format-new": "true"}, &costsResponse); err != nil {
		// A missing shipment (404) OR a costs payload we cannot decode degrades
		// to a shipment WITHOUT costs — the shipment STATUS/destino/carrier from
		// the calls above still stand, so bucket/rastreio survive a costs failure.
		// Only these two non-fatal codes degrade; auth/transient/rate-limit/
		// validation still sink the whole read (never silently zero costs).
		if code := domain.ErrorCodeOf(err); code == domain.ErrCodeProviderInvalidReference || code == domain.ErrCodeProviderPayloadInvalid {
			return result, nil
		}
		return domain.ShipmentInfo{}, mapPricingReaderError(err)
	}

	// A costs MAPPING failure (e.g. an amount that fails domain.ValidateMoney)
	// degrades to the costs-less shipment exactly like a costs DECODE failure — it
	// must never sink the whole read (round-1 invariant: costs failure degrades,
	// status survives). Currency is resolved from the shipment site (new-format
	// /costs carries none — see mlSiteCurrency); prefer the shipment's own site_id,
	// fall back to the account-configured site before the BRL default so a non-BR
	// tenant is not silently relabelled BRL (unknown ≠ default).
	costs, err := mapShipmentCosts(costsResponse, mlSiteCurrency(firstNonEmpty(shipment.SiteID, a.siteID)))
	if err != nil {
		return result, nil
	}
	result.Costs = costs
	return result, nil
}

// fillShipmentCarrier fetches GET /shipments/{id}/carrier and fills
// CarrierName/TrackingURL on result. Any degradable failure (404 or a payload
// we cannot decode) leaves them nil — ML returns 404 as a normal "none"
// condition on shipment sub-resources (docs: /delays), so this is honest
// absence, never a WARN. Non-degradable failures (auth/transient/rate-limit)
// also leave carrier nil rather than sinking the whole shipment read, since the
// status/costs already in hand are the load-bearing facts.
func (a *CapabilityAdapter) fillShipmentCarrier(ctx context.Context, accountRef domain.ProviderAccountRef, token, shipmentPath string, result *domain.ShipmentInfo) {
	var carrier mlShipmentCarrierResponse
	if err := a.doJSONWithHeaders(ctx, accountRef, token, http.MethodGet, shipmentPath+"/carrier", nil, map[string]string{"x-format-new": "true"}, &carrier); err != nil {
		return
	}
	result.CarrierName = trimmedPtr(carrier.Name)
	result.TrackingURL = trimmedPtr(carrier.URL)
}

func (a *CapabilityAdapter) getFreeShippingCost(ctx context.Context, accountRef domain.ProviderAccountRef, token string, query domain.FreeShippingQuery) (domain.FreeShippingCost, error) {
	values := url.Values{}
	values.Set("item_id", strings.TrimSpace(query.ItemID))
	path := "/users/" + url.PathEscape(accountRef.ProviderAccountID) + "/shipping_options/free?" + values.Encode()

	var response mlFreeShippingResponse
	if err := a.doJSON(ctx, accountRef, token, http.MethodGet, path, nil, &response); err != nil {
		return domain.FreeShippingCost{}, mapPricingReaderError(err)
	}

	value := response.Cost
	if response.Coverage != nil && response.Coverage.AllCountry != nil && response.Coverage.AllCountry.ListCost != nil {
		value = response.Coverage.AllCountry.ListCost
	}
	cost, err := providerMoney("free_shipping_cost", value, response.Currency)
	if err != nil {
		return domain.FreeShippingCost{}, err
	}
	return domain.FreeShippingCost{Cost: cost, FetchedAt: a.now()}, nil
}

func mapShipmentInfo(shipment mlShipmentResponse, fetchedAt time.Time) domain.ShipmentInfo {
	var slaDue *time.Time
	if shipment.LeadTime != nil && shipment.LeadTime.EstimatedDeliveryLimit != nil && shipment.LeadTime.EstimatedDeliveryLimit.Date != nil {
		slaDue = parseTimePtr(*shipment.LeadTime.EstimatedDeliveryLimit.Date)
	}

	var destinationUF, destinationCity, destinationZip, receiverName *string
	if shipment.Destination != nil {
		receiverName = trimmedPtr(shipment.Destination.ReceiverName)
		if addr := shipment.Destination.ShippingAddress; addr != nil {
			if addr.State != nil && addr.State.ID != nil {
				// state.id is the ML state code, e.g. "BR-SP" → bare UF "SP".
				value := strings.TrimPrefix(strings.TrimSpace(*addr.State.ID), "BR-")
				if value != "" {
					destinationUF = &value
				}
			}
			if addr.City != nil {
				destinationCity = trimmedPtr(addr.City.Name)
			}
			destinationZip = trimmedPtr(addr.ZipCode)
		}
	}

	var substatus string
	if shipment.Substatus != nil {
		substatus = strings.TrimSpace(*shipment.Substatus)
	}

	return domain.ShipmentInfo{
		ID:              strings.TrimSpace(string(shipment.ID)),
		Status:          strings.TrimSpace(shipment.Status),
		Substatus:       substatus,
		SLADue:          slaDue,
		Delayed:         shipment.Delayed,
		DestinationUF:   destinationUF,
		DestinationCity: destinationCity,
		DestinationZip:  destinationZip,
		ReceiverName:    receiverName,
		FetchedAt:       fetchedAt,
	}
}

// trimmedPtr returns a pointer to the trimmed string, or nil when the input is
// nil or trims to empty. This is how a masked/obfuscated provider value (ML
// obfuscates buyer address fields until payment is confirmed) degrades to
// honest absence instead of a fabricated blank (ADR-17) — the caller never
// warns or errors on a masked value.
func trimmedPtr(s *string) *string {
	if s == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*s)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func mapShipmentCosts(response mlShipmentCostsResponse, currency string) (*domain.ShipmentCosts, error) {
	grossAmount, err := providerMoney("shipment_gross_amount", response.GrossAmount, currency)
	if err != nil {
		return nil, err
	}

	var receiverCost *domain.Money
	if response.Receiver != nil {
		receiverCost, err = providerMoney("shipment_receiver_cost", response.Receiver.Cost, currency)
		if err != nil {
			return nil, err
		}
	}

	var senderCost *domain.Money
	if len(response.Senders) > 0 {
		senderCost, err = providerMoney("shipment_sender_cost", response.Senders[0].Cost, currency)
		if err != nil {
			return nil, err
		}
	}

	return &domain.ShipmentCosts{
		GrossAmount:  grossAmount,
		ReceiverCost: receiverCost,
		SenderCost:   senderCost,
	}, nil
}
