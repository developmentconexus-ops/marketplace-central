package application

import (
	"context"
	"time"

	"marketplace-central/apps/server_core/internal/modules/inventory/domain"
	"marketplace-central/apps/server_core/internal/modules/inventory/ports"
)

type ListStockRisksInput struct {
	InstallationID string
	State          string
	LinkState      string
	Actionability  string
	Limit          int
}

type StockRiskService struct {
	snapshots ports.ListingSnapshotReader
	links     ports.ListingLinkReader
	stocks    ports.InternalStockReader
}

func NewStockRiskService(snapshots ports.ListingSnapshotReader, links ports.ListingLinkReader, stocks ports.InternalStockReader) StockRiskService {
	return StockRiskService{snapshots: snapshots, links: links, stocks: stocks}
}

func (s StockRiskService) List(ctx context.Context, input ListStockRisksInput) ([]domain.StockRiskListItem, error) {
	limit := input.Limit
	if limit <= 0 {
		limit = 50
	}
	snapshots, err := s.snapshots.ListListingSnapshots(ctx, input.InstallationID, limit)
	if err != nil {
		return nil, err
	}
	links, err := s.links.ListLinks(ctx, input.InstallationID, limit)
	if err != nil {
		return nil, err
	}
	linksByKey := map[string]domain.LinkRecord{}
	for _, link := range links {
		linksByKey[keyOf(link.Identity)] = link
	}
	items := make([]domain.StockRiskListItem, 0, len(snapshots))
	for _, snapshot := range snapshots {
		link := linksByKey[keyOf(snapshot.Identity)]
		row, err := s.classifySnapshot(ctx, snapshot, link)
		if err != nil {
			return nil, err
		}
		if matchesRiskFilters(row, input) {
			items = append(items, row)
		}
	}
	return items, nil
}

func (s StockRiskService) classifySnapshot(ctx context.Context, snapshot domain.ListingSnapshot, link domain.LinkRecord) (domain.StockRiskListItem, error) {
	policy := domain.DefaultStockPolicy()
	input := domain.RiskInput{
		Policy: policy,
		Link: domain.LinkEvidence{
			State:             link.State,
			InternalProductID: link.InternalProductID,
		},
		ProviderStock: domain.ProviderStockEvidence{
			AvailableQuantity: snapshot.AvailableQuantity,
			Source:            domain.SourceEvidence{ObservedAt: snapshot.ObservedAt},
		},
		Product: domain.ProductEvidence{},
	}
	if input.Link.State == "" {
		input.Link.State = domain.LinkStateNone
	}
	if link.InternalProductID != nil {
		internalStock, product, err := s.stocks.GetSellableStock(ctx, *link.InternalProductID, policy)
		if err != nil {
			return domain.StockRiskListItem{}, err
		}
		input.InternalStock = internalStock
		input.Product = product
	}
	input.Now = time.Now().UTC()
	row := domain.ClassifyStockRisk(input)
	actionable := row.IsActionable()
	actionability := domain.StockRiskActionabilityBlocked
	if actionable {
		actionability = domain.StockRiskActionabilityActionable
	}
	return domain.StockRiskListItem{
		Identity:              snapshot.Identity,
		ProviderCode:          snapshot.ProviderCode,
		ProviderStatus:        snapshot.ProviderStatus,
		SellerSKU:             snapshot.SellerSKU,
		EAN:                   snapshot.EAN,
		Title:                 snapshot.Title,
		LinkState:             input.Link.State,
		InternalProductID:     link.InternalProductID,
		InternalProductName:   link.InternalProductName,
		InternalReferenceCode: link.InternalReferenceCode,
		State:                 row.State,
		Actionability:         actionability,
		Actionable:            actionable,
		InternalQuantity:      row.InternalQuantity,
		ProviderQuantity:      row.ProviderQuantity,
		RecommendedQuantity:   row.RecommendedQuantity,
		PolicyID:              row.PolicyID,
		InternalObservedAt:    row.InternalObservedAt,
		ProviderObservedAt:    row.ProviderObservedAt,
		BlockingReason:        row.BlockingReason,
	}, nil
}

func matchesRiskFilters(item domain.StockRiskListItem, input ListStockRisksInput) bool {
	if input.State != "" && input.State != string(item.State) {
		return false
	}
	if input.LinkState != "" && input.LinkState != string(item.LinkState) {
		return false
	}
	if input.Actionability != "" && input.Actionability != string(item.Actionability) {
		return false
	}
	return true
}

func keyOf(identity domain.ListingIdentity) string {
	return identity.InstallationID + "|" + identity.ProviderItemID + "|" + identity.ProviderVariationID
}
