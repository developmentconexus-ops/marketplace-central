package fake

import (
	"context"
	"fmt"
	"time"

	"marketplace-central/apps/server_core/internal/modules/internal_read/domain"
	"marketplace-central/apps/server_core/internal/modules/internal_read/ports"
)

type Fixtures struct {
	Products map[string][]domain.ProductCandidate
	Stocks   map[int]domain.SellableStock
	Prices   map[int]domain.CurrentPrice
	Costs    map[string]domain.CostAsOf
	Sales    map[string]domain.SalesHistory
	Taxes    map[string]domain.TaxInputs
}

type Reader struct {
	fixtures Fixtures
}

func NewReader(fixtures Fixtures) *Reader {
	return &Reader{fixtures: fixtures}
}

func (r *Reader) FindProductsForLinking(_ context.Context, input ports.FindProductsInput) ([]domain.ProductCandidate, error) {
	if input.EAN != nil {
		if products, ok := r.fixtures.Products["ean:"+*input.EAN]; ok {
			return products, nil
		}
	}
	if input.SellerSKU != nil {
		if products, ok := r.fixtures.Products["sku:"+*input.SellerSKU]; ok {
			return products, nil
		}
	}
	if input.Title != nil {
		if products, ok := r.fixtures.Products["title:"+*input.Title]; ok {
			return products, nil
		}
	}
	return []domain.ProductCandidate{{
		Source:       domain.SourceMetadata{System: "fake"},
		QualityFlags: []domain.QualityFlag{domain.QualityMissingProduct},
	}}, nil
}

func (r *Reader) GetSellableStock(_ context.Context, input ports.SellableStockInput) (domain.SellableStock, error) {
	if stock, ok := r.fixtures.Stocks[input.ProductID]; ok {
		return stock, nil
	}
	return domain.SellableStock{
		ProductID: input.ProductID,
		Quantity:  nil,
		Policy:    fallbackStockPolicy(input.Policy),
		Source: domain.SourceMetadata{
			System:    "fake",
			FetchedAt: time.Time{},
		},
		QualityFlags: []domain.QualityFlag{domain.QualityMissingStock},
	}, nil
}

func (r *Reader) GetCurrentPrice(_ context.Context, input ports.CurrentPriceInput) (domain.CurrentPrice, error) {
	if price, ok := r.fixtures.Prices[input.ProductID]; ok {
		return price, nil
	}
	policy := fallbackPricePolicy(input.Policy)

	return domain.CurrentPrice{
		ProductID:    input.ProductID,
		PriceTableID: policy.PriceTableID,
		LocationID:   policy.LocationID,
		Amount:       nil,
		Source:       domain.SourceMetadata{System: "fake"},
		QualityFlags: []domain.QualityFlag{domain.QualityMissingPrice},
	}, nil
}

func (r *Reader) GetCostAsOf(_ context.Context, input ports.CostAsOfInput) (domain.CostAsOf, error) {
	policy := fallbackCostPolicy(input.Policy)
	key := fmt.Sprintf("%d:%d:%s", input.ProductID, policy.CompanyID, policy.EffectiveAt.Format(time.DateOnly))
	if cost, ok := r.fixtures.Costs[key]; ok {
		return cost, nil
	}
	return domain.CostAsOf{
		ProductID:   input.ProductID,
		CompanyID:   policy.CompanyID,
		Basis:       policy.Basis,
		EffectiveAt: policy.EffectiveAt,
		Amount:      nil,
		Source:      domain.SourceMetadata{System: "fake"},
		QualityFlags: []domain.QualityFlag{
			domain.QualityMissingCost,
		},
	}, nil
}

func (r *Reader) GetSalesHistory(_ context.Context, input ports.SalesHistoryInput) (domain.SalesHistory, error) {
	productKey := "product:nil"
	if input.ProductID != nil {
		productKey = fmt.Sprintf("product:%d", *input.ProductID)
	}
	groupKey := "group:nil"
	if input.ProductGroupID != nil {
		groupKey = fmt.Sprintf("group:%d", *input.ProductGroupID)
	}
	key := fmt.Sprintf("%s|%s|%s|%s", productKey, groupKey, input.Window.Start.Format(time.DateOnly), input.Window.End.Format(time.DateOnly))
	if sales, ok := r.fixtures.Sales[key]; ok {
		return sales, nil
	}
	return domain.SalesHistory{
		Source: domain.SourceMetadata{System: "fake"},
	}, nil
}

func (r *Reader) GetTaxInputs(_ context.Context, input ports.TaxInput) (domain.TaxInputs, error) {
	policy := fallbackTaxPolicy(input.Policy)
	if !policy.Source.Verified() {
		return domain.TaxInputs{ProductID: input.ProductID, EffectiveAt: policy.EffectiveAt, IncidenceCode: policy.IncidenceCode, Source: domain.SourceMetadata{System: "fake"}, QualityFlags: []domain.QualityFlag{domain.QualityMissingTax}}, nil
	}
	key := fmt.Sprintf("%d:%d:%d:%d", input.ProductID, policy.Source.DocumentID, policy.Source.LineNumber, policy.IncidenceCode)
	if taxes, ok := r.fixtures.Taxes[key]; ok {
		return taxes, nil
	}
	return domain.TaxInputs{
		ProductID:      input.ProductID,
		EffectiveAt:    policy.EffectiveAt,
		IncidenceCode:  policy.IncidenceCode,
		SourceIdentity: policy.Source,
		Source:         domain.SourceMetadata{System: "fake"},
		QualityFlags:   []domain.QualityFlag{domain.QualityMissingTax},
	}, nil
}

func fallbackStockPolicy(policy domain.SellableStockPolicy) domain.SellableStockPolicy {
	if len(policy.CompanyIDs) == 0 && len(policy.LocationIDs) == 0 && len(policy.ExcludedLocationIDs) == 0 && policy.Formula == "" && policy.Scope == "" {
		return domain.DefaultSellableStockPolicy()
	}
	return policy
}

func fallbackPricePolicy(policy domain.CurrentPricePolicy) domain.CurrentPricePolicy {
	if policy.PriceTableID == 0 && policy.LocationID == 0 && policy.EffectiveAt.IsZero() {
		return domain.DefaultCurrentPricePolicy(time.Time{})
	}
	return policy
}

func fallbackCostPolicy(policy domain.CostAsOfPolicy) domain.CostAsOfPolicy {
	if policy.Basis == "" {
		policy.Basis = domain.CostBasisCUSSEMICM
	}
	return policy
}

func fallbackTaxPolicy(policy domain.TaxPolicy) domain.TaxPolicy {
	return policy
}
