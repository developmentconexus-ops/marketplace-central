package internalread

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	catalogdomain "marketplace-central/apps/server_core/internal/modules/catalog/domain"
	readdomain "marketplace-central/apps/server_core/internal/modules/internal_read/domain"
	readports "marketplace-central/apps/server_core/internal/modules/internal_read/ports"
)

// stubReadPortsReader is a hand-rolled readports.Reader whose per-fact returns
// (value + error) are set per test, used to drive product() enrichment paths.
type stubReadPortsReader struct {
	candidate readdomain.ProductCandidate
	stock     readdomain.SellableStock
	stockErr  error
	price     readdomain.CurrentPrice
	priceErr  error
	cost      readdomain.CostAsOf
	costErr   error
}

func (s stubReadPortsReader) FindProductsForLinking(context.Context, readports.FindProductsInput) ([]readdomain.ProductCandidate, error) {
	return []readdomain.ProductCandidate{s.candidate}, nil
}
func (s stubReadPortsReader) GetSellableStock(context.Context, readports.SellableStockInput) (readdomain.SellableStock, error) {
	return s.stock, s.stockErr
}
func (s stubReadPortsReader) GetCurrentPrice(context.Context, readports.CurrentPriceInput) (readdomain.CurrentPrice, error) {
	return s.price, s.priceErr
}
func (s stubReadPortsReader) GetCostAsOf(context.Context, readports.CostAsOfInput) (readdomain.CostAsOf, error) {
	return s.cost, s.costErr
}
func (s stubReadPortsReader) GetSalesHistory(context.Context, readports.SalesHistoryInput) (readdomain.SalesHistory, error) {
	return readdomain.SalesHistory{}, nil
}
func (s stubReadPortsReader) GetTaxInputs(context.Context, readports.TaxInput) (readdomain.TaxInputs, error) {
	return readdomain.TaxInputs{}, nil
}

func candidateForID(id int) readdomain.ProductCandidate {
	pid := readdomain.InternalProductID(id)
	return readdomain.ProductCandidate{InternalProductID: &pid, ProductID: id, Name: "Produto"}
}

// TestCanonicalProductDegradesUnavailableFactsToMissing is the regression guard
// for the xlsx 500: an existing product whose stock/price/cost reads return an
// honest "not available in this source" read error (unsupported_query — the xlsx
// GetCurrentPrice signal — or source_unavailable) must resolve to a
// CanonicalProduct with those facts marked missing and NO error, never abort.
func TestCanonicalProductDegradesUnavailableFactsToMissing(t *testing.T) {
	for _, code := range []readdomain.ReadErrorCode{readdomain.ReadErrorUnsupportedQuery, readdomain.ReadErrorSourceUnavailable} {
		t.Run(string(code), func(t *testing.T) {
			degradable := readdomain.NewReadError(code, "not available in this source mode", nil)
			r := &Reader{reader: stubReadPortsReader{
				candidate: candidateForID(90001),
				stockErr:  degradable,
				priceErr:  degradable,
				costErr:   degradable,
			}}
			got, err := r.GetCanonicalProduct(context.Background(), catalogdomain.InternalProductID(90001))
			if err != nil {
				t.Fatalf("expected degrade to missing, got error: %v", err)
			}
			if got.InternalProductID != catalogdomain.InternalProductID(90001) {
				t.Fatalf("product identity lost: %+v", got)
			}
			for name, f := range map[string]catalogdomain.NumericSourceFact{"stock": got.StockQuantity, "price": got.PriceAmount, "cost": got.CostAmount} {
				if f.Quality != catalogdomain.SourceFactQualityUnknown || f.Value != nil {
					t.Fatalf("%s fact should be missing, got quality=%s value=%v", name, f.Quality, f.Value)
				}
			}
		})
	}
}

// TestListProductsByIDsDegradesUnavailableFacts proves the same degrade holds
// through the by-ID list entrypoint pricing.Exists depends on.
func TestListProductsByIDsDegradesUnavailableFacts(t *testing.T) {
	degradable := readdomain.NewReadError(readdomain.ReadErrorUnsupportedQuery, "xlsx snapshot does not contain current price", nil)
	r := &Reader{reader: stubReadPortsReader{candidate: candidateForID(90001), priceErr: degradable}}
	products, err := r.ListProductsByIDs(context.Background(), []string{"90001"})
	if err != nil {
		t.Fatalf("expected degrade, got error: %v", err)
	}
	if len(products) != 1 {
		t.Fatalf("expected 1 product, got %d", len(products))
	}
}

// TestCanonicalProductPropagatesUnexpectedFactError is the over-degrade guard:
// a read error that is NOT an availability signal must still surface, never be
// swallowed into a fake-missing fact.
func TestCanonicalProductPropagatesUnexpectedFactError(t *testing.T) {
	boom := errors.New("db connection reset")
	cases := map[string]stubReadPortsReader{
		"stock": {candidate: candidateForID(90001), stockErr: boom},
		"price": {candidate: candidateForID(90001), priceErr: boom},
		"cost":  {candidate: candidateForID(90001), costErr: boom},
	}
	for name, stub := range cases {
		t.Run(name, func(t *testing.T) {
			r := &Reader{reader: stub}
			if _, err := r.GetCanonicalProduct(context.Background(), catalogdomain.InternalProductID(90001)); err == nil {
				t.Fatalf("expected unexpected %s error to propagate, got nil", name)
			}
		})
	}
}

func TestFactProjectsQualityFlagsWithoutInventingCurrentData(t *testing.T) {
	now := time.Date(2026, 7, 13, 12, 0, 0, 0, time.UTC)
	value := 0.0
	tests := []struct {
		name      string
		value     *float64
		flags     []readdomain.QualityFlag
		want      catalogdomain.SourceFactQuality
		wantValue bool
	}{
		{name: "complete known zero", value: &value, flags: []readdomain.QualityFlag{readdomain.QualityComplete}, want: catalogdomain.SourceFactQualityCurrent, wantValue: true},
		{name: "stale keeps prior value", value: &value, flags: []readdomain.QualityFlag{readdomain.QualityStaleSource}, want: catalogdomain.SourceFactQualityStale, wantValue: true},
		{name: "missing discards value", value: &value, flags: []readdomain.QualityFlag{readdomain.QualityMissingCost}, want: catalogdomain.SourceFactQualityUnknown},
		{name: "conflict has no canonical value", value: &value, flags: []readdomain.QualityFlag{readdomain.QualityAmbiguousProduct}, want: catalogdomain.SourceFactQualityConflict},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fact(tt.value, readdomain.SourceMetadata{System: "oracle", FetchedAt: now}, tt.flags, readdomain.QualityMissingCost, "cost unavailable")
			if err != nil {
				t.Fatal(err)
			}
			if got.Quality != tt.want || (got.Value != nil) != tt.wantValue {
				t.Fatalf("fact = %+v, want quality=%s value=%v", got, tt.want, tt.wantValue)
			}
		})
	}
}

func TestCatalogPageIdentityProjectsWithoutDroppingFlags(t *testing.T) {
	ean, reference, brand, ncm := "4006381333931", "MF-1", "Marca", "12345678"
	page := readports.CatalogFactPage{AsOf: time.Now().UTC(), Items: []readports.CatalogProductFact{{InternalProductID: 1, EAN: &ean, Reference: &reference, ManufacturerReference: &reference, BrandName: &brand, NCM: &ncm, QualityFlags: []string{"complete", "ean_collision"}}}}
	products, err := canonicalProductsFromPage(page)
	if err != nil {
		t.Fatal(err)
	}
	got := products[0]
	if got.EAN == nil || *got.EAN != ean || got.ManufacturerReference == nil || *got.ManufacturerReference != reference || got.BrandName == nil || *got.BrandName != brand || got.NCM == nil || *got.NCM != ncm || !reflect.DeepEqual(got.QualityFlags, []string{"complete", "ean_collision"}) {
		t.Fatalf("identity = %+v", got)
	}
}

func TestFactReturnsConstructorErrors(t *testing.T) {
	value := 1.0
	_, err := fact(&value, readdomain.SourceMetadata{System: "oracle"}, []readdomain.QualityFlag{readdomain.QualityComplete}, readdomain.QualityMissingCost, "cost unavailable")
	if err == nil {
		t.Fatal("expected current fact without observed time to return constructor error")
	}
}
