package composition

import (
	"context"
	"errors"
	"testing"
	"time"

	integrationsdomain "marketplace-central/apps/server_core/internal/modules/integrations/domain"
)

type fakeInstallationLister struct {
	installations []integrationsdomain.Installation
	err           error
}

func (l *fakeInstallationLister) List(context.Context) ([]integrationsdomain.Installation, error) {
	return l.installations, l.err
}

// TestResolveOrdersSchedulersOnlyMercadoLivre proves the fan-out filters to
// Mercado Livre installations only. The shopee installation is the negative
// control: without the filter, a fan-out that just returned "one scheduler
// per installation" would pass equally.
func TestResolveOrdersSchedulersOnlyMercadoLivre(t *testing.T) {
	lister := &fakeInstallationLister{installations: []integrationsdomain.Installation{
		{InstallationID: "inst-ml-1", TenantID: "t1", ProviderCode: "mercado_livre"},
		{InstallationID: "inst-shopee-1", TenantID: "t1", ProviderCode: "shopee"},
		{InstallationID: "inst-ml-2", TenantID: "t1", ProviderCode: "mercado_livre"},
	}}

	group := NewOrdersSchedulers(nil, time.Minute, 50, 0, lister, nil)
	got := resolveOrdersSchedulers(context.Background(), group)
	if len(got) != 2 {
		t.Fatalf("schedulers: quero 2 (so' mercado_livre), recebi %d", len(got))
	}
}

// TestResolveOrdersSchedulersReturnsNilOnListError proves the fail-open boot
// posture: a List error must not crash composition and must yield zero
// schedulers rather than a partial/guessed list.
func TestResolveOrdersSchedulersReturnsNilOnListError(t *testing.T) {
	group := NewOrdersSchedulers(nil, time.Minute, 50, 0, &fakeInstallationLister{err: errors.New("banco fora")}, nil)
	got := resolveOrdersSchedulers(context.Background(), group)
	if got != nil {
		t.Fatalf("falha ao listar nao pode virar schedulers fantasma; recebi %d", len(got))
	}
}
