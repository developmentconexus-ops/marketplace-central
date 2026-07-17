package ports

import "context"

type PolicyReader interface {
	HasPricingPolicy(context.Context, string) (bool, error)
}
