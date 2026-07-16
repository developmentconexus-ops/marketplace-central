package ports

import "context"

type CostFact struct {
	Amount *float64
}

type ICMSCeiling struct {
	Percent *float64
}

type CostReader interface {
	GetCostFactsByIDs(context.Context, []int64) (map[int64]*CostFact, error)
	GetICMSCeilingByOrigin(context.Context, int64) (map[int64]*ICMSCeiling, error)
}

type PricingPolicy struct{ MinMarginPercent float64 }

type PolicyReader interface {
	GetPricingPolicyForInstallation(context.Context, string) (PricingPolicy, bool, error)
}

type InstallationReader interface {
	InstallationExists(context.Context, string) (bool, error)
}
