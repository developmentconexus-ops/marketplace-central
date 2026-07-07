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
	return []domain.ProductCandidate{}, nil
}

func (r *Reader) GetSellableStock(_ context.Context, input ports.SellableStockInput) (domain.SellableStock, error) {
	if stock, ok := r.fixtures.Stocks[input.Codprod]; ok {
		return stock, nil
	}
	return domain.SellableStock{
		Codprod:         input.Codprod,
		Quantity:        0,
		Scope:           domain.DefaultSellableStockScope(),
		QualityFlags:    []domain.QualityFlag{domain.QualityMissingStock},
		SourceFetchedAt: time.Time{},
	}, nil
}

func (r *Reader) GetCurrentPrice(_ context.Context, input ports.CurrentPriceInput) (domain.CurrentPrice, error) {
	if price, ok := r.fixtures.Prices[input.Codprod]; ok {
		return price, nil
	}

	codtab := 0
	if input.Codtab != nil {
		codtab = *input.Codtab
	}
	codlocal := 10101
	if input.Codlocal != nil {
		codlocal = *input.Codlocal
	}

	return domain.CurrentPrice{
		Codprod:      input.Codprod,
		Codtab:       codtab,
		Codlocal:     codlocal,
		Price:        nil,
		QualityFlags: []domain.QualityFlag{domain.QualityMissingCost},
	}, nil
}

func (r *Reader) GetCostAsOf(_ context.Context, input ports.CostAsOfInput) (domain.CostAsOf, error) {
	key := fmt.Sprintf("%d:%d:%s", input.Codprod, input.Codemp, input.SaleDate)
	if cost, ok := r.fixtures.Costs[key]; ok {
		return cost, nil
	}
	return domain.CostAsOf{
		Codprod:      input.Codprod,
		Codemp:       input.Codemp,
		SaleDate:     input.SaleDate,
		CUSSEMICM:    nil,
		QualityFlags: []domain.QualityFlag{domain.QualityMissingCost},
	}, nil
}

func (r *Reader) GetSalesHistory(_ context.Context, input ports.SalesHistoryInput) (domain.SalesHistory, error) {
	productKey := "product:nil"
	if input.Codprod != nil {
		productKey = fmt.Sprintf("product:%d", *input.Codprod)
	}
	groupKey := "group:nil"
	if input.CodgrupoProd != nil {
		groupKey = fmt.Sprintf("group:%d", *input.CodgrupoProd)
	}
	key := fmt.Sprintf("%s|%s|%s|%s", productKey, groupKey, input.StartDate, input.EndDate)
	if sales, ok := r.fixtures.Sales[key]; ok {
		return sales, nil
	}
	return domain.SalesHistory{}, nil
}

func (r *Reader) GetTaxInputs(_ context.Context, input ports.TaxInput) (domain.TaxInputs, error) {
	key := fmt.Sprintf("%d:%s", input.Codprod, input.SaleDate)
	if taxes, ok := r.fixtures.Taxes[key]; ok {
		return taxes, nil
	}
	return domain.TaxInputs{
		Codprod:      input.Codprod,
		SaleDate:     input.SaleDate,
		QualityFlags: []domain.QualityFlag{domain.QualityMissingTax},
	}, nil
}
