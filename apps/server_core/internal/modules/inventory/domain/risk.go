package domain

import "time"

type LinkState string

const (
	LinkStateNone       LinkState = "none"
	LinkStateUnresolved LinkState = "unresolved"
	LinkStateConflict   LinkState = "conflict"
	LinkStateResolved   LinkState = "resolved"
	LinkStateRejected   LinkState = "rejected"
)

type StockRiskState string

const (
	StockRiskHealthy     StockRiskState = "healthy"
	StockRiskOversell    StockRiskState = "oversell"
	StockRiskUndersell   StockRiskState = "undersell"
	StockRiskStale       StockRiskState = "stale"
	StockRiskUnresolved  StockRiskState = "unresolved"
	StockRiskConflict    StockRiskState = "conflict"
	StockRiskIneligible  StockRiskState = "ineligible"
	StockRiskUnsupported StockRiskState = "unsupported"
)

type LinkEvidence struct {
	State             LinkState
	InternalProductID *int
}

type InternalStockEvidence struct {
	Quantity *int
	Source   SourceEvidence
}

type ProviderStockEvidence struct {
	AvailableQuantity *int
	Source            SourceEvidence
}

type RiskInput struct {
	Policy        StockPolicy
	Link          LinkEvidence
	InternalStock InternalStockEvidence
	ProviderStock ProviderStockEvidence
	Product       ProductEvidence
	Now           time.Time
}

type BlockingReason struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type StockRiskRow struct {
	State               StockRiskState
	InternalQuantity    *int
	ProviderQuantity    *int
	RecommendedQuantity *int
	PolicyID            string
	InternalObservedAt  *time.Time
	ProviderObservedAt  *time.Time
	BlockingReason      BlockingReason
}

func ClassifyStockRisk(input RiskInput) StockRiskRow {
	row := StockRiskRow{
		InternalQuantity:   input.InternalStock.Quantity,
		ProviderQuantity:   input.ProviderStock.AvailableQuantity,
		PolicyID:           input.Policy.PolicyID,
		InternalObservedAt: input.InternalStock.Source.ObservedAt,
		ProviderObservedAt: input.ProviderStock.Source.ObservedAt,
	}

	switch input.Link.State {
	case LinkStateResolved:
		if input.Link.InternalProductID == nil {
			return block(row, StockRiskUnresolved, "missing_internal_product_link", "resolved link has no internal product")
		}
	case LinkStateRejected:
		return block(row, StockRiskUnresolved, "rejected_link", "product link was rejected")
	case LinkStateConflict:
		return block(row, StockRiskConflict, "conflict_link", "product link is in conflict")
	default:
		return block(row, StockRiskUnresolved, "unresolved_link", "product link is not resolved")
	}

	if input.InternalStock.Quantity == nil {
		return block(row, StockRiskStale, "missing_internal_stock", "internal sellable stock is missing")
	}
	if freshness := input.Policy.Freshness.Evaluate(input.InternalStock.Source, input.Now); freshness.State == SourceFreshnessStale {
		return block(row, StockRiskStale, "stale_internal_source", freshness.Reason)
	}
	if input.ProviderStock.AvailableQuantity == nil {
		return block(row, StockRiskUnsupported, "missing_provider_quantity", "provider stock quantity is missing")
	}
	if freshness := input.Policy.Freshness.Evaluate(input.ProviderStock.Source, input.Now); freshness.State == SourceFreshnessStale {
		return block(row, StockRiskStale, "stale_provider_source", freshness.Reason)
	}
	if eligibility := input.Policy.Eligibility.Evaluate(input.Product); eligibility.State == ProductEligibilityIneligible {
		return block(row, StockRiskIneligible, "ineligible_product", eligibility.Reason)
	}

	recommended := input.Policy.RecommendedProviderQuantity(*input.InternalStock.Quantity)
	row.RecommendedQuantity = &recommended

	switch {
	case *input.ProviderStock.AvailableQuantity > recommended:
		row.State = StockRiskOversell
	case *input.ProviderStock.AvailableQuantity < recommended:
		row.State = StockRiskUndersell
	default:
		row.State = StockRiskHealthy
	}
	return row
}

func block(row StockRiskRow, state StockRiskState, code string, message string) StockRiskRow {
	row.State = state
	row.RecommendedQuantity = nil
	row.BlockingReason = BlockingReason{Code: code, Message: message}
	return row
}
