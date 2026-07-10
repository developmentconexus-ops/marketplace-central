package domain

import (
	"errors"
	"strings"
	"time"
)

type RuntimeCapabilityCode string
type RuntimeCapabilityState string

const (
	RuntimeCapabilityAccountProbe   RuntimeCapabilityCode = "account_probe"
	RuntimeCapabilityListingRead    RuntimeCapabilityCode = "listing_read"
	RuntimeCapabilityOrderRead      RuntimeCapabilityCode = "order_read"
	RuntimeCapabilityFeeQuoteRead   RuntimeCapabilityCode = "fee_quote_read"
	RuntimeCapabilityStockRead      RuntimeCapabilityCode = "stock_read"
	RuntimeCapabilityStockWrite     RuntimeCapabilityCode = "stock_write"
	RuntimeCapabilityMessageReply   RuntimeCapabilityCode = "message_reply"
	RuntimeCapabilityShipmentRead   RuntimeCapabilityCode = "shipment_read"
	RuntimeCapabilityWebhookReceive RuntimeCapabilityCode = "webhook_receive"

	RuntimeCapabilityStateAvailable     RuntimeCapabilityState = "available"
	RuntimeCapabilityStateUnavailable   RuntimeCapabilityState = "unavailable"
	RuntimeCapabilityStateNeedsAuth     RuntimeCapabilityState = "needs_auth"
	RuntimeCapabilityStateDegraded      RuntimeCapabilityState = "degraded"
	RuntimeCapabilityStateNotConfigured RuntimeCapabilityState = "not_configured"
)

type RuntimeCapability struct {
	Code              RuntimeCapabilityCode  `json:"code"`
	State             RuntimeCapabilityState `json:"state"`
	Executable        bool                   `json:"executable"`
	LiveValidated     bool                   `json:"live_validated"`
	LocalValidated    bool                   `json:"local_validated"`
	UnavailableReason string                 `json:"unavailable_reason,omitempty"`
	LastValidatedAt   *time.Time             `json:"last_validated_at,omitempty"`
}

func (c RuntimeCapability) Available() bool {
	return c.Code != "" && c.State == RuntimeCapabilityStateAvailable && c.Executable
}

func (c RuntimeCapability) Runnable() bool {
	if c.Code == "" || !c.Executable {
		return false
	}

	switch c.State {
	case RuntimeCapabilityStateAvailable, RuntimeCapabilityStateDegraded:
		return true
	default:
		return false
	}
}

func ValidateRuntimeCapability(capability RuntimeCapability) error {
	if strings.TrimSpace(string(capability.Code)) == "" {
		return errors.New("INTEGRATIONS_CAPABILITY_INVALID")
	}
	if capability.State == "" {
		return errors.New("INTEGRATIONS_CAPABILITY_INVALID")
	}

	return nil
}
