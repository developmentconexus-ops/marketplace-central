package internalread

import (
	internalreaddomain "marketplace-central/apps/server_core/internal/modules/internal_read/domain"
	inventorydomain "marketplace-central/apps/server_core/internal/modules/inventory/domain"
)

func MapStockPolicy(policy inventorydomain.StockPolicy) internalreaddomain.SellableStockPolicy {
	return internalreaddomain.SellableStockPolicy{
		CompanyIDs:          append([]int(nil), policy.Scope.CompanyIDs...),
		LocationIDs:         append([]int(nil), policy.Scope.LocationIDs...),
		ExcludedLocationIDs: append([]int(nil), policy.Scope.ExcludedLocationIDs...),
		Formula:             policy.Scope.Formula,
		Scope:               mapStockScope(policy.Scope.Scope),
	}
}

func MapFreshnessPolicy(policy inventorydomain.FreshnessPolicy) internalreaddomain.FreshnessPolicy {
	return internalreaddomain.FreshnessPolicy{MaxAge: policy.MaxAge}
}

func mapStockScope(scope inventorydomain.StockScope) internalreaddomain.StockScope {
	switch scope {
	case inventorydomain.StockScopeCustom:
		return internalreaddomain.SellableStockScopeCustom
	default:
		return internalreaddomain.SellableStockScopeResale
	}
}
