package ports

import (
	"context"
	"time"
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

// OrderRef é o resultado da enumeração no vocabulário de orders (F-00 Task 3).
// O que ImportService sempre consumiu do snapshot cheio era o id e, quando
// presente, a data de última atualização — o resto (itens, pagamentos) nunca
// tinha consumidor aqui; IngestOrder refaz a leitura completa pelo caminho de
// escrita único. ProviderUpdatedAt fica ponteiro pelo mesmo motivo do
// domain.OrderSearchHit de onde vem: ausente e "zero" são fatos diferentes.
type OrderRef struct {
	ProviderOrderID   string
	ProviderUpdatedAt *time.Time
}

type OrderSource interface {
	ListOrders(ctx context.Context, input ListOrdersInput) ([]OrderRef, error)
}
