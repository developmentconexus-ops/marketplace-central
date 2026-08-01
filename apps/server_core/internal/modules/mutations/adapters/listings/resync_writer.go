package listings

import (
	"context"
	"errors"
	"strings"

	listingsports "marketplace-central/apps/server_core/internal/modules/listings/ports"
	mutationsdomain "marketplace-central/apps/server_core/internal/modules/mutations/domain"
	mutationsports "marketplace-central/apps/server_core/internal/modules/mutations/ports"
)

// BackfillPuller is the whole-catalog synchronous pull operation this writer
// delegates to — the SAME shape as listings/application's unexported
// refreshPuller interface RefreshService depends on. Root.go wires a single
// shared listingsapp.BackfillRunner instance to BOTH producers (the manual
// API refresh and this admin resync writer) through this interface, so they
// share identical hydrate+persist behavior end to end (ADR-04: one writer,
// no coexisting path — the old per-writer *listingsapp.Ingestion this writer
// used to construct fresh per Apply call is gone; F-04).
type BackfillPuller interface {
	Pull(ctx context.Context, account listingsports.InstallationAccount) error
}

// ResyncWriter applies an admin-triggered listing resync by delegating to a
// shared BackfillPuller, scoped to one installation account. It no longer
// owns any hydrate/persist logic itself (that lives in the puller) — its own
// job is strictly the guard clauses below plus the installation-scope check.
type ResyncWriter struct {
	puller  BackfillPuller
	account listingsports.InstallationAccount
}

func NewResyncWriter(puller BackfillPuller, account listingsports.InstallationAccount) *ResyncWriter {
	return &ResyncWriter{puller: puller, account: account}
}

func (w *ResyncWriter) Apply(ctx context.Context, item mutationsports.WriteItem) (mutationsports.WriteOutcome, error) {
	if w == nil || w.puller == nil {
		return mutationsports.WriteOutcome{}, errors.New("listing resync writer is not configured")
	}
	if item.ProtocolType != mutationsdomain.ProtocolTypeListingResync {
		return mutationsports.WriteOutcome{}, errors.New("listing resync writer received unsupported protocol type")
	}
	if strings.TrimSpace(w.account.TenantID) == "" || strings.TrimSpace(w.account.InstallationID) == "" || strings.TrimSpace(item.InstallationID) != strings.TrimSpace(w.account.InstallationID) {
		return mutationsports.WriteOutcome{}, errors.New("listing resync installation scope is invalid")
	}
	if err := w.puller.Pull(ctx, w.account); err != nil {
		return mutationsports.WriteOutcome{}, err
	}
	return mutationsports.WriteOutcome{}, nil
}
