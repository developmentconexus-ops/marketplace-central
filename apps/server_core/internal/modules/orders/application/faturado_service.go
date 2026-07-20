package application

import (
	"context"
	"errors"
	"strings"

	"marketplace-central/apps/server_core/internal/modules/orders/ports"
)

// ErrFaturadoInvalidInput signals a missing installationID/providerOrderID
// on MarkFaturado — a client error, never a store call.
var ErrFaturadoInvalidInput = errors.New("ORDERS_FATURADO_INVALID_REQUEST")

// FaturadoService wraps ports.OrderFaturadoStore for the "Foi faturado"
// operator action (goal B). Our-DB only; never a Mercado Livre write.
type FaturadoService struct {
	store ports.OrderFaturadoStore
}

func NewFaturadoService(store ports.OrderFaturadoStore) FaturadoService {
	return FaturadoService{store: store}
}

func (s FaturadoService) MarkFaturado(ctx context.Context, installationID, providerOrderID string) error {
	installationID = strings.TrimSpace(installationID)
	providerOrderID = strings.TrimSpace(providerOrderID)
	if installationID == "" || providerOrderID == "" {
		return ErrFaturadoInvalidInput
	}
	if s.store == nil {
		return errors.New("ORDERS_FATURADO_STORE_NOT_CONFIGURED")
	}
	return s.store.MarkOrderFaturado(ctx, installationID, providerOrderID)
}
