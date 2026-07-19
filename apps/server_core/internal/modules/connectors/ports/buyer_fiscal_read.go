package ports

import (
	"context"

	"marketplace-central/apps/server_core/internal/modules/connectors/domain"
)

// BuyerFiscalReader resolves a buyer's fiscal identity (name, document, billing address) for
// an order, via the provider's documented two-step flow (order -> billing_info id ->
// billing-info). A buyer without billing data degrades to an empty, honest-absence
// BuyerFiscalInfo (BuyerFiscalInfo.HasData() == false), never an error.
type BuyerFiscalReader interface {
	GetBuyerFiscalInfo(ctx context.Context, accountRef domain.ProviderAccountRef, providerOrderID string) (domain.BuyerFiscalInfo, error)
}
