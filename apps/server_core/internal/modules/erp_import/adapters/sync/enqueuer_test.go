package erpsync

import (
	"context"
	"errors"
	"reflect"
	"testing"

	integrationsdomain "marketplace-central/apps/server_core/internal/modules/integrations/domain"
	syncdomain "marketplace-central/apps/server_core/internal/modules/sync/domain"
)

type fakeInstallationLister struct {
	installations []integrationsdomain.Installation
	err           error
}

func (f *fakeInstallationLister) List(context.Context) ([]integrationsdomain.Installation, error) {
	return f.installations, f.err
}

type appendCall struct {
	installationID string
	entity         syncdomain.Entity
	codigo         string
}

type fakeMarketCursorAppender struct {
	calls   []appendCall
	failAt  int
	failErr error
}

func (f *fakeMarketCursorAppender) AppendPendingCodigo(_ context.Context, installationID string, entity syncdomain.Entity, codigo string) error {
	f.calls = append(f.calls, appendCall{installationID: installationID, entity: entity, codigo: codigo})
	if f.failAt > 0 && len(f.calls) == f.failAt {
		return f.failErr
	}
	return nil
}

func TestMarketEnqueuerEnqueuesEveryCodeForEveryInstallation(t *testing.T) {
	lister := &fakeInstallationLister{installations: []integrationsdomain.Installation{
		{InstallationID: "inst-a"},
		{InstallationID: "inst-b"},
	}}
	appender := &fakeMarketCursorAppender{}

	installationIDs, err := NewMarketEnqueuer(lister, appender).EnqueueMarketProducts(context.Background(), []string{"100", "200", "300"})
	if err != nil {
		t.Fatalf("EnqueueMarketProducts() error = %v", err)
	}

	wantIDs := []string{"inst-a", "inst-b"}
	if !reflect.DeepEqual(installationIDs, wantIDs) {
		t.Fatalf("installationIDs = %#v, want %#v", installationIDs, wantIDs)
	}

	wantCalls := []appendCall{
		{installationID: "inst-a", entity: syncdomain.EntityMarket, codigo: "100"},
		{installationID: "inst-a", entity: syncdomain.EntityMarket, codigo: "200"},
		{installationID: "inst-a", entity: syncdomain.EntityMarket, codigo: "300"},
		{installationID: "inst-b", entity: syncdomain.EntityMarket, codigo: "100"},
		{installationID: "inst-b", entity: syncdomain.EntityMarket, codigo: "200"},
		{installationID: "inst-b", entity: syncdomain.EntityMarket, codigo: "300"},
	}
	if !reflect.DeepEqual(appender.calls, wantCalls) {
		t.Fatalf("append calls = %#v, want %#v", appender.calls, wantCalls)
	}
}

func TestMarketEnqueuerPropagatesAppendError(t *testing.T) {
	wantErr := errors.New("append failed")
	lister := &fakeInstallationLister{installations: []integrationsdomain.Installation{{InstallationID: "inst-a"}}}
	appender := &fakeMarketCursorAppender{failAt: 1, failErr: wantErr}

	_, err := NewMarketEnqueuer(lister, appender).EnqueueMarketProducts(context.Background(), []string{"100"})
	if !errors.Is(err, wantErr) {
		t.Fatalf("EnqueueMarketProducts() error = %v, want wrapped %v", err, wantErr)
	}
}

func TestMarketEnqueuerEmptyCodesDoesNotAppend(t *testing.T) {
	lister := &fakeInstallationLister{installations: []integrationsdomain.Installation{{InstallationID: "inst-a"}}}
	appender := &fakeMarketCursorAppender{}

	installationIDs, err := NewMarketEnqueuer(lister, appender).EnqueueMarketProducts(context.Background(), nil)
	if err != nil {
		t.Fatalf("EnqueueMarketProducts() error = %v", err)
	}
	if len(installationIDs) != 0 {
		t.Fatalf("installationIDs = %#v, want empty", installationIDs)
	}
	if len(appender.calls) != 0 {
		t.Fatalf("append calls = %#v, want none", appender.calls)
	}
}

func TestMarketEnqueuerPropagatesListError(t *testing.T) {
	wantErr := errors.New("list failed")
	lister := &fakeInstallationLister{err: wantErr}
	appender := &fakeMarketCursorAppender{}

	_, err := NewMarketEnqueuer(lister, appender).EnqueueMarketProducts(context.Background(), []string{"100"})
	if !errors.Is(err, wantErr) {
		t.Fatalf("EnqueueMarketProducts() error = %v, want wrapped %v", err, wantErr)
	}
	if len(appender.calls) != 0 {
		t.Fatalf("append calls = %#v, want none", appender.calls)
	}
}
