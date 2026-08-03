package ports

import (
	"context"
	"time"

	"marketplace-central/apps/server_core/internal/modules/orders/domain"
)

// ListOrdersInput é a janela de enumeração. Ela existe como struct — e não como
// mais dois parâmetros posicionais — porque a assinatura estreita
// (installationID, limit) já custou uma feature: o adapter Mercado Livre sabe
// traduzir UpdatedAfter para order.date_last_updated.from desde sempre, e três
// camadas intermediárias jogavam o valor fora antes de chegar nele.
type ListOrdersInput struct {
	InstallationID string
	Limit          int
	Offset         int
	// UpdatedAfter nil significa "sem janela" (varredura), não "desde a época
	// zero". Os dois produzem consultas diferentes no provider.
	UpdatedAfter *time.Time
}

type OrderSource interface {
	ListOrders(ctx context.Context, input ListOrdersInput) ([]domain.OrderIngestionSnapshot, error)
}
