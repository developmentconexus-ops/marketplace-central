package application

import (
	"context"
	"time"

	"marketplace-central/apps/server_core/internal/modules/inventory/domain"
	"marketplace-central/apps/server_core/internal/modules/inventory/ports"
	"marketplace-central/apps/server_core/internal/modules/sourcekind"
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
	stocks    ports.InternalStockBatchReader
}

func NewStockRiskService(snapshots ports.ListingSnapshotReader, links ports.ListingLinkReader, stocks ports.InternalStockBatchReader) StockRiskService {
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

	stockIDs := make([]int64, 0, len(snapshots))
	seenStockIDs := make(map[int64]struct{}, len(snapshots))
	for _, snapshot := range snapshots {
		link := linksByKey[keyOf(snapshot.Identity)]
		if link.InternalProductID == nil {
			continue
		}
		id := int64(*link.InternalProductID)
		if _, seen := seenStockIDs[id]; seen {
			continue
		}
		seenStockIDs[id] = struct{}{}
		stockIDs = append(stockIDs, id)
	}
	stockFacts := map[int64]*ports.StockFact{}
	if len(stockIDs) > 0 {
		stockFacts, err = s.stocks.GetStockFactsByIDs(ctx, stockIDs)
		if err != nil {
			return nil, err
		}
	}

	items := make([]domain.StockRiskListItem, 0, len(snapshots))
	for _, snapshot := range snapshots {
		link := linksByKey[keyOf(snapshot.Identity)]
		row := s.classifySnapshot(snapshot, link, stockFacts)
		if matchesRiskFilters(row, input) {
			items = append(items, row)
		}
	}
	return items, nil
}

func (s StockRiskService) classifySnapshot(snapshot domain.ListingSnapshot, link domain.LinkRecord, stockFacts map[int64]*ports.StockFact) domain.StockRiskListItem {
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
		productID := int64(*link.InternalProductID)
		input.Product = domain.ProductEvidence{ProductID: *link.InternalProductID}
		if fact, ok := stockFacts[productID]; ok && fact != nil {
			input.InternalStock.Quantity = stockQuantityToInt(fact.Quantity)
			input.InternalStock.Source.ObservedAt = stockObservedAt(fact)
			// The freshness bar depends on how the source produces data: a live ERP
			// read is expected to be minutes old, an uploaded snapshot only as new
			// as the last import.
			input.InternalStock.Source.Kind = sourcekind.Of(fact.Source.System)
		}
	}
	input.Now = time.Now().UTC()
	row := domain.ClassifyStockRisk(input)
	actionable := row.IsActionable()
	actionability := domain.StockRiskActionabilityBlocked
	if actionable {
		actionability = domain.StockRiskActionabilityActionable
	}
	result := domain.StockRiskListItem{
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
		InternalSourceKind:    input.InternalStock.Source.Kind,
		ProviderObservedAt:    row.ProviderObservedAt,
		BlockingReason:        row.BlockingReason,
	}
	if link.InternalProductID != nil {
		fact, found := stockFacts[int64(*link.InternalProductID)]
		if !found || fact == nil {
			result.QualityFlags = []string{"missing_stock"}
		} else {
			result.QualityFlags = stockQualityFlags(fact)
		}
	}
	return result
}

func stockQuantityToInt(value *float64) *int {
	if value == nil {
		return nil
	}
	quantity := int(*value)
	return &quantity
}

func stockObservedAt(fact *ports.StockFact) *time.Time {
	if fact.Source.ObservedAt != nil {
		observedAt := fact.Source.ObservedAt.UTC()
		return &observedAt
	}
	if fact.Source.FetchedAt.IsZero() {
		return nil
	}
	fetchedAt := fact.Source.FetchedAt.UTC()
	return &fetchedAt
}

func stockQualityFlags(fact *ports.StockFact) []string {
	flags := make([]string, 0, len(fact.QualityFlags))
	for _, flag := range fact.QualityFlags {
		flags = append(flags, string(flag))
	}
	return flags
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
