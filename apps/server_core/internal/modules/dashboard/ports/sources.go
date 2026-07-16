package ports

import (
	"context"
	"time"

	integrationsdomain "marketplace-central/apps/server_core/internal/modules/integrations/domain"
	integrationports "marketplace-central/apps/server_core/internal/modules/integrations/ports"
	listingsports "marketplace-central/apps/server_core/internal/modules/listings/ports"
	ordersports "marketplace-central/apps/server_core/internal/modules/orders/ports"
	linkageports "marketplace-central/apps/server_core/internal/modules/product_links/ports"
)

type InstallationSource interface {
	Get(context.Context, string) (integrationsdomain.Installation, bool, error)
}

type ListingsSource interface {
	Summary(context.Context, listingsports.SummaryQuery) (listingsports.ListingSummaryRow, error)
}

type LinkageSource interface {
	Summary(context.Context, string) (linkageports.LinkageSummary, error)
}

type OrdersSource interface {
	Summary(context.Context, string, time.Time) (ordersports.OrderSummary, error)
}

type SyncSource interface {
	LatestRunsByModule(context.Context, string) ([]integrationports.LatestRunByModule, error)
}
