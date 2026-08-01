package composition

import (
	"context"
	"fmt"

	mercadolivreconnector "marketplace-central/apps/server_core/internal/modules/connectors/adapters/mercado_livre"
	connectorsdomain "marketplace-central/apps/server_core/internal/modules/connectors/domain"
	integrationsapp "marketplace-central/apps/server_core/internal/modules/integrations/application"
	ordersports "marketplace-central/apps/server_core/internal/modules/orders/ports"
)

// orders_ingest_adapters.go wires the M-03 F-02 ingest ports (OrderDetailReader,
// ShipmentDetailReader) to concrete composition-owned adapters. Kept in a SEPARATE file from
// orders_adapters.go (F-02 feature.md: that file is a different feature's territory, not to be
// edited) even though the pattern — resolve installationID via installations.Get, build a
// connectorsdomain.ProviderAccountRef via accountRefForInstallation, delegate to the
// mercado_livre CapabilityAdapter — is identical to ordersShipmentReaderAdapter/
// ordersBuyerFiscalReaderAdapter there. Lives in composition for the same reason those do:
// modules/connectors already imports modules/orders indirectly through this composition root,
// so an orders-side adapter over connectors risks an import cycle if placed inside modules/orders.

// --- ports.OrderDetailReader --------------------------------------------------------

type ordersOrderDetailReaderAdapter struct {
	capabilities  *mercadolivreconnector.CapabilityAdapter
	installations *integrationsapp.InstallationService
	tenantID      string
}

var _ ordersports.OrderDetailReader = (*ordersOrderDetailReaderAdapter)(nil)

func newOrdersOrderDetailReaderAdapter(
	capabilities *mercadolivreconnector.CapabilityAdapter,
	installations *integrationsapp.InstallationService,
	tenantID string,
) *ordersOrderDetailReaderAdapter {
	return &ordersOrderDetailReaderAdapter{
		capabilities:  capabilities,
		installations: installations,
		tenantID:      tenantID,
	}
}

func (r *ordersOrderDetailReaderAdapter) GetOrderDetail(ctx context.Context, installationID, providerOrderID string) (connectorsdomain.OrderDetail, error) {
	inst, found, err := r.installations.Get(ctx, installationID)
	if err != nil {
		return connectorsdomain.OrderDetail{}, err
	}
	if !found {
		return connectorsdomain.OrderDetail{}, fmt.Errorf("orders: installation %q not found", installationID)
	}
	ref := accountRefForInstallation(r.tenantID, inst)
	return r.capabilities.GetOrderDetail(ctx, ref, providerOrderID)
}

// --- ports.ShipmentDetailReader ------------------------------------------------------

type ordersShipmentDetailReaderAdapter struct {
	capabilities  *mercadolivreconnector.CapabilityAdapter
	installations *integrationsapp.InstallationService
	tenantID      string
}

var _ ordersports.ShipmentDetailReader = (*ordersShipmentDetailReaderAdapter)(nil)

func newOrdersShipmentDetailReaderAdapter(
	capabilities *mercadolivreconnector.CapabilityAdapter,
	installations *integrationsapp.InstallationService,
	tenantID string,
) *ordersShipmentDetailReaderAdapter {
	return &ordersShipmentDetailReaderAdapter{
		capabilities:  capabilities,
		installations: installations,
		tenantID:      tenantID,
	}
}

func (r *ordersShipmentDetailReaderAdapter) GetShipmentDetail(ctx context.Context, installationID, shipmentID string) (connectorsdomain.ShipmentDetail, error) {
	inst, found, err := r.installations.Get(ctx, installationID)
	if err != nil {
		return connectorsdomain.ShipmentDetail{}, err
	}
	if !found {
		return connectorsdomain.ShipmentDetail{}, fmt.Errorf("orders: installation %q not found", installationID)
	}
	ref := accountRefForInstallation(r.tenantID, inst)
	return r.capabilities.GetShipmentDetail(ctx, ref, shipmentID)
}
