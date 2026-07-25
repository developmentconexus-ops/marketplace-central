package erpsync

import (
	"context"
	"fmt"

	erpapp "marketplace-central/apps/server_core/internal/modules/erp_import/application"
	integrationsdomain "marketplace-central/apps/server_core/internal/modules/integrations/domain"
	syncdomain "marketplace-central/apps/server_core/internal/modules/sync/domain"
)

type installationLister interface {
	List(ctx context.Context) ([]integrationsdomain.Installation, error)
}

type marketCursorAppender interface {
	AppendPendingCodigos(ctx context.Context, installationID string, entity syncdomain.Entity, codigos []string) error
}

type MarketEnqueuer struct {
	installations installationLister
	cursor        marketCursorAppender
}

func NewMarketEnqueuer(installations installationLister, cursor marketCursorAppender) *MarketEnqueuer {
	return &MarketEnqueuer{installations: installations, cursor: cursor}
}

func (e *MarketEnqueuer) EnqueueMarketProducts(ctx context.Context, productCodes []string) ([]string, error) {
	if len(productCodes) == 0 {
		return []string{}, nil
	}

	installations, err := e.installations.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list installations for market enqueue: %w", err)
	}

	installationIDs := make([]string, 0, len(installations))
	seen := make(map[string]struct{}, len(installations))
	for _, installation := range installations {
		// One statement per installation, not one per code: an ERP import carries a
		// couple of thousand codes and the per-code round trips dominated the upload.
		if err := e.cursor.AppendPendingCodigos(ctx, installation.InstallationID, syncdomain.EntityMarket, productCodes); err != nil {
			return nil, fmt.Errorf("append %d market products for installation %q: %w", len(productCodes), installation.InstallationID, err)
		}
		if _, ok := seen[installation.InstallationID]; !ok {
			seen[installation.InstallationID] = struct{}{}
			installationIDs = append(installationIDs, installation.InstallationID)
		}
	}

	return installationIDs, nil
}

var _ erpapp.SyncEnqueuer = (*MarketEnqueuer)(nil)
