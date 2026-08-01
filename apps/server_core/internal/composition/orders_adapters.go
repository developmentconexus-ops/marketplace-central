package composition

import (
	"context"
	"fmt"
	"time"

	mercadolivreconnector "marketplace-central/apps/server_core/internal/modules/connectors/adapters/mercado_livre"
	connectorsdomain "marketplace-central/apps/server_core/internal/modules/connectors/domain"
	integrationsapp "marketplace-central/apps/server_core/internal/modules/integrations/application"
	internalreadapp "marketplace-central/apps/server_core/internal/modules/internal_read/application"
	internalreaddomain "marketplace-central/apps/server_core/internal/modules/internal_read/domain"
	internalreadports "marketplace-central/apps/server_core/internal/modules/internal_read/ports"
	ordersports "marketplace-central/apps/server_core/internal/modules/orders/ports"
)

// orders_adapters.go wires the M-08 orders enrichment ports (CostReader,
// ShipmentReader) to concrete composition-owned adapters. It mirrors
// market_adapters.go end-to-end (same IC-02 CostAsOf policy, same
// ProviderAccountRef resolution via accountRefForInstallation) and lives in
// composition for the same reason: modules/connectors already imports
// modules/orders indirectly through this composition root, so an
// orders-side adapter over connectors would risk an import cycle if placed
// inside modules/orders.

// --- ports.CostReader --------------------------------------------------------

// ordersCostReaderAdapter backs ports.CostReader with the same IC-02
// internal_read.GetCostAsOf pattern as marketCostReaderAdapter
// (market_adapters.go:40-78), but returns the internal_read CostAsOf
// straight through since EnrichService reads cost.Amount directly (no Money
// mapping needed here, unlike market's). When no internal_read source is
// configured (available == false) or the source is unavailable, it honestly
// returns a zero-value CostAsOf (nil Amount) rather than a fabricated cost.
type ordersCostReaderAdapter struct {
	service   internalreadapp.Service
	available bool
}

var _ ordersports.CostReader = (*ordersCostReaderAdapter)(nil)

func newOrdersCostReaderAdapter(service internalreadapp.Service, available bool) *ordersCostReaderAdapter {
	return &ordersCostReaderAdapter{service: service, available: available}
}

func (r *ordersCostReaderAdapter) GetCostAsOf(ctx context.Context, productID int, effectiveAt time.Time) (internalreaddomain.CostAsOf, error) {
	if !r.available {
		return internalreaddomain.CostAsOf{}, nil
	}
	cost, err := r.service.GetCostAsOf(ctx, internalreadports.CostAsOfInput{
		ProductID: productID,
		Policy: internalreaddomain.CostAsOfPolicy{
			CompanyID:   1,
			EffectiveAt: effectiveAt,
			Basis:       internalreaddomain.CostBasisCUSSEMICM,
		},
		Freshness: internalreaddomain.FreshnessPolicy{},
	})
	if err != nil {
		if internalreaddomain.IsReadErrorCode(err, internalreaddomain.ReadErrorSourceUnavailable) {
			return internalreaddomain.CostAsOf{}, nil
		}
		return internalreaddomain.CostAsOf{}, err
	}
	return cost, nil
}

// --- ports.ShipmentReader ----------------------------------------------------
//
// ordersShipmentReaderAdapter (the mercado_livre-backed ports.ShipmentReader
// adapter, formerly here) was retired by F-03 (read-path switch): the
// interactive GET /orders/{id} enrich path now reads order_shipments
// (migration 0088, populated by F-02's IngestOrder) via
// orders/adapters/postgres.ShipmentReader instead, wired directly in
// root.go — no composition adapter needed since the Postgres reader already
// implements ports.ShipmentReader.

// --- ports.BuyerFiscalReader -------------------------------------------------

// ordersBuyerFiscalReaderAdapter backs ports.BuyerFiscalReader with the
// mercado_livre connectors capability adapter, resolving installationID to a
// connectorsdomain.ProviderAccountRef via the same accountRefForInstallation
// helper as ordersShipmentReaderAdapter. The connectors adapter runs the whole
// documented two-step billing-info flow and already degrades a buyer without
// billing data to an honest-absence value; EnrichService degrades any error
// here to a nil block plus a warn-once (enrich_service.go resolveBuyerFiscal),
// so the provider error passes straight up unmapped.
type ordersBuyerFiscalReaderAdapter struct {
	capabilities  *mercadolivreconnector.CapabilityAdapter
	installations *integrationsapp.InstallationService
	tenantID      string
}

var _ ordersports.BuyerFiscalReader = (*ordersBuyerFiscalReaderAdapter)(nil)

func newOrdersBuyerFiscalReaderAdapter(
	capabilities *mercadolivreconnector.CapabilityAdapter,
	installations *integrationsapp.InstallationService,
	tenantID string,
) *ordersBuyerFiscalReaderAdapter {
	return &ordersBuyerFiscalReaderAdapter{
		capabilities:  capabilities,
		installations: installations,
		tenantID:      tenantID,
	}
}

func (r *ordersBuyerFiscalReaderAdapter) GetBuyerFiscal(ctx context.Context, installationID, providerOrderID string) (connectorsdomain.BuyerFiscalInfo, error) {
	inst, found, err := r.installations.Get(ctx, installationID)
	if err != nil {
		return connectorsdomain.BuyerFiscalInfo{}, err
	}
	if !found {
		return connectorsdomain.BuyerFiscalInfo{}, fmt.Errorf("orders: installation %q not found", installationID)
	}
	ref := accountRefForInstallation(r.tenantID, inst)
	return r.capabilities.GetBuyerFiscalInfo(ctx, ref, providerOrderID)
}
