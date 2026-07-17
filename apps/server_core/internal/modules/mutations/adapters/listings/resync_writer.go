package listings

import (
	"context"
	"errors"
	"strings"
	"time"

	listingsapp "marketplace-central/apps/server_core/internal/modules/listings/application"
	listingsports "marketplace-central/apps/server_core/internal/modules/listings/ports"
	mutationsdomain "marketplace-central/apps/server_core/internal/modules/mutations/domain"
	mutationsports "marketplace-central/apps/server_core/internal/modules/mutations/ports"
)

type ResyncWriter struct {
	source   listingsports.PageSource
	store    listingsports.CompletedPullStore
	account  listingsports.InstallationAccount
	pageSize int
	now      func() time.Time
}

func NewResyncWriter(source listingsports.PageSource, store listingsports.CompletedPullStore, account listingsports.InstallationAccount, pageSize int, now func() time.Time) *ResyncWriter {
	return &ResyncWriter{source: source, store: store, account: account, pageSize: pageSize, now: now}
}

func (w *ResyncWriter) Apply(ctx context.Context, item mutationsports.WriteItem) (mutationsports.WriteOutcome, error) {
	if w == nil || w.source == nil || w.store == nil || w.now == nil || w.pageSize <= 0 {
		return mutationsports.WriteOutcome{}, errors.New("listing resync writer is not configured")
	}
	if item.ProtocolType != mutationsdomain.ProtocolTypeListingResync {
		return mutationsports.WriteOutcome{}, errors.New("listing resync writer received unsupported protocol type")
	}
	if strings.TrimSpace(w.account.TenantID) == "" || strings.TrimSpace(w.account.InstallationID) == "" || strings.TrimSpace(item.InstallationID) != strings.TrimSpace(w.account.InstallationID) {
		return mutationsports.WriteOutcome{}, errors.New("listing resync installation scope is invalid")
	}
	ingestion := listingsapp.NewIngestion(w.source, w.store, w.pageSize, w.now)
	if err := ingestion.Pull(ctx, w.account); err != nil {
		return mutationsports.WriteOutcome{}, err
	}
	return mutationsports.WriteOutcome{}, nil
}
