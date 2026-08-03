package background

import (
	"context"
	"log/slog"
	"time"

	"marketplace-central/apps/server_core/internal/modules/integrations/application"
	"marketplace-central/apps/server_core/internal/modules/integrations/domain"
)

type expiringSessionLister interface {
	ListExpiringSessions(ctx context.Context, expiresWithin time.Duration) ([]domain.AuthSession, error)
}

type credentialRefresher interface {
	RefreshCredential(ctx context.Context, input application.RefreshCredentialInput) (application.AuthStatus, error)
}

type RefreshTicker struct {
	sessions      expiringSessionLister
	flow          credentialRefresher
	logger        *slog.Logger
	interval      time.Duration
	expiresWithin time.Duration
	stop          chan struct{}
}

// NewRefreshTicker aceita logger nil e cai em slog.Default(), mesmo contrato de
// mutations/background/poller.go:26 — chamador de teste não precisa montar um
// logger só para exercitar o loop.
func NewRefreshTicker(sessions expiringSessionLister, flow credentialRefresher, interval time.Duration, logger *slog.Logger) *RefreshTicker {
	if logger == nil {
		logger = slog.Default()
	}
	return &RefreshTicker{
		sessions:      sessions,
		flow:          flow,
		logger:        logger,
		interval:      interval,
		expiresWithin: 10 * time.Minute,
		stop:          make(chan struct{}),
	}
}

func (t *RefreshTicker) RunOnce(ctx context.Context) error {
	sessions, err := t.sessions.ListExpiringSessions(ctx, t.expiresWithin)
	if err != nil {
		// Falhar em LISTAR é falha do lote: não há o que iterar. Sobe.
		return err
	}

	for _, session := range sessions {
		// Falha de UM item não é falha do lote. Abortar aqui fazia a primeira
		// conta com token quebrado impedir todas as seguintes de tentarem —
		// invisível com uma conta só, fatal com duas. Mesmo padrão de
		// mutations/background/poller.go:70-72.
		if _, err := t.flow.RefreshCredential(ctx, application.RefreshCredentialInput{
			InstallationID: session.InstallationID,
		}); err != nil {
			t.logger.Error("integrations refresh ticker item failed",
				"installation_id", session.InstallationID,
				"err", err,
			)
		}
	}
	return nil
}

func (t *RefreshTicker) Start(ctx context.Context) {
	ticker := time.NewTicker(t.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-t.stop:
			return
		case <-ticker.C:
			// O erro era descartado com `_ =`: uma listagem falhando em todo
			// tick não deixava rastro nenhum (F-A1).
			if err := t.RunOnce(ctx); err != nil {
				t.logger.Error("integrations refresh ticker pass failed", "err", err)
			}
		}
	}
}

func (t *RefreshTicker) Stop() {
	select {
	case <-t.stop:
	default:
		close(t.stop)
	}
}
