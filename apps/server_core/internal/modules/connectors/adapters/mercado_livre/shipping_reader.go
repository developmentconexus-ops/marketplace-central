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
	ID              string              `json:"id"`
	Status          string              `json:"status"`
	Substatus       *string             `json:"substatus"`
	LeadTime        *mlShipmentLeadTime `json:"lead_time"`
	Delayed         *bool               `json:"delayed"`
	ReceiverAddress *mlReceiverAddress  `json:"receiver_address"`
}

type mlShipmentLeadTime struct {
	EstimatedDeliveryLimit *mlEstimatedDeliveryLimit `json:"estimated_delivery_limit"`
}

type mlEstimatedDeliveryLimit struct {
	Date *string `json:"date"`
}

type mlReceiverAddress struct {
	State *mlShipmentState `json:"state"`
}

type mlShipmentState struct {
	ID *string `json:"id"`
}

type mlShipmentCostsResponse struct {
	GrossAmount *json.Number            `json:"gross_amount"`
	Receiver    *mlShipmentCostReceiver `json:"receiver"`
	Senders     []mlShipmentCostSender  `json:"senders"`
	Currency    string                  `json:"currency_id"`
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
	if err := a.doJSON(ctx, accountRef, token, http.MethodGet, path, nil, &shipment); err != nil {
		return domain.ShipmentInfo{}, mapPricingReaderError(err)
	}

	var costsResponse mlShipmentCostsResponse
	costsPath := path + "/costs"
	// The ML shipments /costs endpoint REQUIRES `x-format-new: true`; without it
	// ML returns the legacy cost shape which fails to decode into
	// mlShipmentCostsResponse (ML docs: mercado-envios-2 "Get Shipment Details
	// with New Format").
	if err := a.doJSONWithHeaders(ctx, accountRef, token, http.MethodGet, costsPath, nil, map[string]string{"x-format-new": "true"}, &costsResponse); err != nil {
		// A missing shipment (404) OR a costs payload we cannot decode degrades
		// to a shipment WITHOUT costs — the shipment STATUS from the first call
		// still stands, so the order's bucket/rastreio survive a costs failure.
		// Only these two non-fatal codes degrade; auth/transient/rate-limit/
		// validation still sink the whole read (never silently zero costs).
		if code := domain.ErrorCodeOf(err); code == domain.ErrCodeProviderInvalidReference || code == domain.ErrCodeProviderPayloadInvalid {
			return mapShipmentInfo(shipment, a.now()), nil
		}
		return domain.ShipmentInfo{}, mapPricingReaderError(err)
	}

	costs, err := mapShipmentCosts(costsResponse)
	if err != nil {
		return domain.ShipmentInfo{}, err
	}
	result := mapShipmentInfo(shipment, a.now())
	result.Costs = costs
	return result, nil
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

	var destinationUF *string
	if shipment.ReceiverAddress != nil && shipment.ReceiverAddress.State != nil && shipment.ReceiverAddress.State.ID != nil {
		value := strings.TrimPrefix(strings.TrimSpace(*shipment.ReceiverAddress.State.ID), "BR-")
		destinationUF = &value
	}

	var substatus string
	if shipment.Substatus != nil {
		substatus = strings.TrimSpace(*shipment.Substatus)
	}

	return domain.ShipmentInfo{
		ID:            strings.TrimSpace(shipment.ID),
		Status:        strings.TrimSpace(shipment.Status),
		Substatus:     substatus,
		SLADue:        slaDue,
		Delayed:       shipment.Delayed,
		DestinationUF: destinationUF,
		FetchedAt:     fetchedAt,
	}
}

func mapShipmentCosts(response mlShipmentCostsResponse) (*domain.ShipmentCosts, error) {
	grossAmount, err := providerMoney("shipment_gross_amount", response.GrossAmount, response.Currency)
	if err != nil {
		return nil, err
	}

	var receiverCost *domain.Money
	if response.Receiver != nil {
		receiverCost, err = providerMoney("shipment_receiver_cost", response.Receiver.Cost, response.Currency)
		if err != nil {
			return nil, err
		}
	}

	var senderCost *domain.Money
	if len(response.Senders) > 0 {
		senderCost, err = providerMoney("shipment_sender_cost", response.Senders[0].Cost, response.Currency)
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
