package application

import (
	"context"
	"errors"
	"strings"
	"time"

	"marketplace-central/apps/server_core/internal/modules/inventory/domain"
	"marketplace-central/apps/server_core/internal/modules/inventory/ports"
)

type ApplyManualActionInput struct {
	ActionID            string
	InstallationID      string
	ProviderItemID      string
	ProviderVariationID string
	RequestedQuantity   int
	Reason              string
	Approval            domain.ManualApproval
}

type ApplyManualActionResult struct {
	Action domain.StockAction       `json:"action"`
	Risk   domain.StockRiskListItem `json:"risk"`
}

type ManualActionFacade struct {
	risks         StockRiskService
	actions       StockActionService
	installations ports.InstallationReader
	snapshots     ports.ListingSnapshotStore
	now           func() time.Time
}

func NewManualActionFacade(risks StockRiskService, actions StockActionService, installations ports.InstallationReader, snapshots ports.ListingSnapshotStore) ManualActionFacade {
	return ManualActionFacade{
		risks:         risks,
		actions:       actions,
		installations: installations,
		snapshots:     snapshots,
		now:           func() time.Time { return time.Now().UTC() },
	}
}

func (f ManualActionFacade) Apply(ctx context.Context, input ApplyManualActionInput) (ApplyManualActionResult, error) {
	risks, err := f.risks.List(ctx, ListStockRisksInput{
		InstallationID: input.InstallationID,
		Limit:          100,
	})
	if err != nil {
		return ApplyManualActionResult{}, err
	}
	risk, found := findRisk(risks, input.InstallationID, input.ProviderItemID, input.ProviderVariationID)
	if !found {
		return ApplyManualActionResult{}, errors.New("INVENTORY_STOCK_RISK_NOT_FOUND")
	}
	ref, found, err := f.installations.GetInstallation(ctx, input.InstallationID)
	if err != nil {
		return ApplyManualActionResult{}, err
	}
	if !found {
		return ApplyManualActionResult{}, errors.New("INTEGRATIONS_INSTALLATION_NOT_FOUND")
	}
	ref.ProviderItemID = input.ProviderItemID
	ref.ProviderVariationID = input.ProviderVariationID
	now := f.now()
	action, err := f.actions.ApplyManual(ctx, ApplyManualStockActionInput{
		ActionID: input.ActionID,
		Risk: domain.RiskInput{
			Policy: domain.DefaultStockPolicy(),
			Link: domain.LinkEvidence{
				State:             risk.LinkState,
				InternalProductID: risk.InternalProductID,
			},
			InternalStock: domain.InternalStockEvidence{
				Quantity: risk.InternalQuantity,
				Source:   domain.SourceEvidence{ObservedAt: risk.InternalObservedAt},
			},
			ProviderStock: domain.ProviderStockEvidence{
				AvailableQuantity: risk.ProviderQuantity,
				Source:            domain.SourceEvidence{ObservedAt: risk.ProviderObservedAt},
			},
			Product: domain.ProductEvidence{
				ProductID: derefInt(risk.InternalProductID),
			},
			Now: now,
		},
		ProviderRef:       ref,
		RequestedQuantity: input.RequestedQuantity,
		Reason:            strings.TrimSpace(input.Reason),
		Approval:          input.Approval,
		Now:               now,
	})
	if err != nil {
		return ApplyManualActionResult{}, err
	}
	if action.State == domain.StockActionStateApplied {
		observedAt := now.UTC()
		if err := f.snapshots.SaveListingSnapshot(ctx, domain.ListingSnapshot{
			Identity:          risk.Identity,
			ProviderCode:      risk.ProviderCode,
			ProviderStatus:    risk.ProviderStatus,
			SellerSKU:         risk.SellerSKU,
			EAN:               risk.EAN,
			Title:             risk.Title,
			AvailableQuantity: &input.RequestedQuantity,
			ObservedAt:        &observedAt,
		}); err != nil {
			return ApplyManualActionResult{}, err
		}
		risk.ProviderQuantity = &input.RequestedQuantity
		risk.ProviderObservedAt = &observedAt
		risk.RecommendedQuantity = &input.RequestedQuantity
		risk.State = domain.StockRiskHealthy
		risk.Actionability = domain.StockRiskActionabilityBlocked
		risk.Actionable = false
		risk.BlockingReason = domain.BlockingReason{}
	}
	return ApplyManualActionResult{Action: action, Risk: risk}, nil
}

func findRisk(items []domain.StockRiskListItem, installationID, providerItemID, providerVariationID string) (domain.StockRiskListItem, bool) {
	for _, item := range items {
		if item.Identity.InstallationID == installationID &&
			item.Identity.ProviderItemID == providerItemID &&
			item.Identity.ProviderVariationID == providerVariationID {
			return item, true
		}
	}
	return domain.StockRiskListItem{}, false
}

func derefInt(value *int) int {
	if value == nil {
		return 0
	}
	return *value
}
