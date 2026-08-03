package background

import (
	"context"
	"errors"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/modules/integrations/application"
	"marketplace-central/apps/server_core/internal/modules/integrations/domain"
)

type refreshSessionStore struct {
	items   []domain.AuthSession
	listErr error
}

func (s refreshSessionStore) ListExpiringSessions(context.Context, time.Duration) ([]domain.AuthSession, error) {
	if s.listErr != nil {
		return nil, s.listErr
	}
	return append([]domain.AuthSession(nil), s.items...), nil
}

type refreshFlow struct {
	inputs  []application.RefreshCredentialInput
	failFor map[string]error
}

func (s *refreshFlow) RefreshCredential(ctx context.Context, input application.RefreshCredentialInput) (application.AuthStatus, error) {
	s.inputs = append(s.inputs, input)
	if err, ok := s.failFor[input.InstallationID]; ok {
		return application.AuthStatus{}, err
	}
	return application.AuthStatus{InstallationID: input.InstallationID, Status: domain.InstallationStatusConnected}, nil
}

func TestRefreshTickerUsesListExpiringSessions(t *testing.T) {
	t.Parallel()

	flow := &refreshFlow{}
	job := NewRefreshTicker(refreshSessionStore{items: []domain.AuthSession{
		{InstallationID: "installation_expiring_1"},
		{InstallationID: "installation_expiring_2"},
	}}, flow, time.Minute, nil)

	if err := job.RunOnce(context.Background()); err != nil {
		t.Fatalf("RunOnce() error = %v", err)
	}
	if len(flow.inputs) != 2 {
		t.Fatalf("refresh count = %d, want 2", len(flow.inputs))
	}
	if flow.inputs[0].InstallationID != "installation_expiring_1" || flow.inputs[1].InstallationID != "installation_expiring_2" {
		t.Fatalf("refresh inputs = %#v, want expiring-session installation IDs", flow.inputs)
	}
}

func TestRunOnceContinuesAfterItemFailure(t *testing.T) {
	t.Parallel()

	sessions := refreshSessionStore{
		items: []domain.AuthSession{
			{InstallationID: "inst-a"},
			{InstallationID: "inst-b"},
			{InstallationID: "inst-c"},
		},
	}
	flow := &refreshFlow{failFor: map[string]error{"inst-a": errors.New("boom")}}

	ticker := NewRefreshTicker(sessions, flow, time.Minute, nil)
	if err := ticker.RunOnce(context.Background()); err != nil {
		t.Fatalf("RunOnce err = %v, want nil (falha de item não é falha do lote)", err)
	}

	if len(flow.inputs) != 3 {
		t.Fatalf("chamadas = %#v, want 3 (a falha de inst-a matou o resto do lote)", flow.inputs)
	}
}

func TestRunOnceReturnsErrorWhenListingFails(t *testing.T) {
	t.Parallel()

	// Controle negativo: falhar em LISTAR é falha do lote inteiro e continua
	// subindo. Sem este teste, "RunOnce sempre retorna nil" passaria vacuoso.
	sessions := refreshSessionStore{listErr: errors.New("db down")}
	ticker := NewRefreshTicker(sessions, &refreshFlow{}, time.Minute, nil)

	if err := ticker.RunOnce(context.Background()); err == nil {
		t.Fatal("RunOnce err = nil, want erro de listagem")
	}
}
