package mercadolivre

import (
	"encoding/json"
	"strings"
)

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
