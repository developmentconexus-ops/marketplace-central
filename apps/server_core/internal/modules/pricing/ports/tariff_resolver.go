package ports

import (
	"context"

	"marketplace-central/apps/server_core/internal/modules/pricing/domain"
)

// TariffRequest carries the inputs needed to resolve one tariff. Degrau 4
// reads only Modalidade. Degrau 3 (live COTAÇÃO) additionally reads the
// per-product id and the price basis to quote against; both are OPTIONAL —
// when either is nil the composite resolver skips degrau 3 and degrau 4
// answers (ADR-17: a missing input is never coerced into a live quote).
type TariffRequest struct {
	Modalidade domain.Modalidade
	// ProductID identifies the product whose live category/commission to
	// quote. nil => degrau 3 skipped (e.g. a manual-commission override, or a
	// request with no resolvable product).
	ProductID *int
	// PriceBasis is the decimal price string the commission is quoted at
	// (the decompose target price). nil => degrau 3 falls back to the
	// product's catalog price, and skips if that is absent too.
	PriceBasis *string
}

// TariffResolver resolves the commission + shipping tariff components for a
// TariffRequest. Degrau 4 (adapters/tariffdefaults) implements this purely
// from tenant config defaults; a future composite resolver will chain
// degraus 1->2->3->4 behind the same port with no interface change.
type TariffResolver interface {
	Resolve(ctx context.Context, req TariffRequest) (domain.TariffResolution, error)
}
