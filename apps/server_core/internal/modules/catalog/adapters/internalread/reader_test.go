package internalread

import (
	"reflect"
	"testing"
	"time"

	catalogdomain "marketplace-central/apps/server_core/internal/modules/catalog/domain"
	readdomain "marketplace-central/apps/server_core/internal/modules/internal_read/domain"
	readports "marketplace-central/apps/server_core/internal/modules/internal_read/ports"
)

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
