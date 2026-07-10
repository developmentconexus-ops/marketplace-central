package internalread

import (
	"reflect"
	"testing"
	"time"

	internalreaddomain "marketplace-central/apps/server_core/internal/modules/internal_read/domain"
	inventorydomain "marketplace-central/apps/server_core/internal/modules/inventory/domain"
)

func TestMapStockPolicyPreservesInternalReadSellableScope(t *testing.T) {
	policy := inventorydomain.DefaultStockPolicy()

	got := MapStockPolicy(policy)

	if !reflect.DeepEqual(got.CompanyIDs, []int{1, 2}) {
		t.Fatalf("expected companies [1 2], got %v", got.CompanyIDs)
	}
	if !reflect.DeepEqual(got.LocationIDs, []int{10101}) {
		t.Fatalf("expected locations [10101], got %v", got.LocationIDs)
	}
	if !reflect.DeepEqual(got.ExcludedLocationIDs, []int{10108}) {
		t.Fatalf("expected excluded showroom [10108], got %v", got.ExcludedLocationIDs)
	}
	if got.Formula != "SUM(ESTOQUE - RESERVADO)" {
		t.Fatalf("expected formula SUM(ESTOQUE - RESERVADO), got %q", got.Formula)
	}
	if got.Scope != internalreaddomain.SellableStockScopeResale {
		t.Fatalf("expected internal read resale scope, got %q", got.Scope)
	}
}

func TestMapFreshnessPolicyPreservesMaxAge(t *testing.T) {
	got := MapFreshnessPolicy(inventorydomain.FreshnessPolicy{MaxAge: 45 * time.Minute})

	if got.MaxAge != 45*time.Minute {
		t.Fatalf("expected 45m max age, got %s", got.MaxAge)
	}
}

func TestMapStockPolicyPreservesCustomScope(t *testing.T) {
	policy := inventorydomain.DefaultStockPolicy()
	policy.Scope.Scope = inventorydomain.StockScopeCustom

	got := MapStockPolicy(policy)

	if got.Scope != internalreaddomain.SellableStockScopeCustom {
		t.Fatalf("expected custom scope, got %q", got.Scope)
	}
}
