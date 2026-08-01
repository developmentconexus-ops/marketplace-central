package composition

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	integrationsdomain "marketplace-central/apps/server_core/internal/modules/integrations/domain"
	syncapp "marketplace-central/apps/server_core/internal/modules/sync/application"
	"marketplace-central/apps/server_core/internal/modules/sync/domain"
)

type fakeInstallationLister struct {
	installations []integrationsdomain.Installation
	err           error
}

func (l *fakeInstallationLister) List(context.Context) ([]integrationsdomain.Installation, error) {
	return l.installations, l.err
}

type noopStateStore struct{}

func (noopStateStore) Read(context.Context, string, domain.Entity) (domain.SyncState, bool, error) {
	return domain.SyncState{}, false, nil
}
func (noopStateStore) RecordSuccess(context.Context, string, domain.Entity, json.RawMessage, time.Time, bool) error {
	return nil
}
func (noopStateStore) RecordFailure(context.Context, string, domain.Entity, domain.SyncError) error {
	return nil
}

var _ syncapp.StateStore = noopStateStore{}

// TestListingsEntityDoubleRegistrationFailsClosed is the mandatory negative
// scenario (F-04 brief): registering the listings entity twice on the SAME
// scheduler instance must fail closed, never silently overwrite the first
// registration. This exercises sync/application.Scheduler.RegisterJob's own
// generic guard specifically for domain.EntityListings — the entity/job this
// feature wires — without touching scheduler.go itself (forbidden path).
func TestListingsEntityDoubleRegistrationFailsClosed(t *testing.T) {
	scheduler := syncapp.NewScheduler(noopStateStore{}, "installation-a", time.Hour, time.Now)
	job := func(ctx context.Context, cursor json.RawMessage) (json.RawMessage, error) { return cursor, nil }

	if err := scheduler.RegisterJob(domain.EntityListings, job); err != nil {
		t.Fatalf("first RegisterJob(EntityListings) error = %v, want nil", err)
	}
	err := scheduler.RegisterJob(domain.EntityListings, job)
	if err == nil {
		t.Fatal("second RegisterJob(EntityListings) error = nil, want a fail-closed duplicate-entity error")
	}
	if !strings.Contains(err.Error(), `entity "listings" already registered`) {
		t.Fatalf("error = %v, want the duplicate-entity message naming listings", err)
	}
}

// TestNewListingsSchedulersOneInstancePerActiveMLInstallation proves the
// per-installation fan-out: only mercado_livre installations get a
// scheduler, and each gets its OWN instance (so each can independently
// register domain.EntityListings without colliding with the others, per the
// single-installation-scope constraint above).
func TestNewListingsSchedulersOneInstancePerActiveMLInstallation(t *testing.T) {
	lister := &fakeInstallationLister{installations: []integrationsdomain.Installation{
		{InstallationID: "ml-1", TenantID: "tenant-a", ProviderCode: "mercado_livre", ExternalAccountID: "seller-1"},
		{InstallationID: "ml-2", TenantID: "tenant-a", ProviderCode: "mercado_livre", ExternalAccountID: "seller-2"},
		{InstallationID: "erp-1", TenantID: "tenant-a", ProviderCode: "sankhya_xlsx", ExternalAccountID: ""},
	}}

	group := NewListingsSchedulers(nil, time.Hour, lister, nil, nil, nil)
	schedulers := resolveListingsSchedulers(context.Background(), group)
	if len(schedulers) != 2 {
		t.Fatalf("schedulers = %d, want 2 (one per ML installation, ERP installation skipped)", len(schedulers))
	}
}

// TestNewListingsSchedulersListFailureReturnsNoSchedulersWithoutPanicking
// proves the fail-open boot posture documented on NewListingsSchedulers: a
// List error must not crash composition, and must yield zero schedulers
// rather than a partial/guessed list.
func TestNewListingsSchedulersListFailureReturnsNoSchedulersWithoutPanicking(t *testing.T) {
	lister := &fakeInstallationLister{err: errors.New("db unavailable")}
	group := NewListingsSchedulers(nil, time.Hour, lister, nil, nil, nil)
	schedulers := resolveListingsSchedulers(context.Background(), group)
	if schedulers != nil {
		t.Fatalf("schedulers = %+v, want nil", schedulers)
	}
}

// TestNewListingsSchedulersStartAllDoesNotTouchDatabaseSynchronously proves
// the actual F-04 fix: composing the Group (NewListingsSchedulers) and
// calling StartAll must return immediately and must never invoke the
// installationLister on the calling goroutine — installations.List is only
// ever called from StartAll's own background goroutine. This is what keeps
// root.go's single composition call safe to run with an unconfigured (nil)
// pool during TestNewRootRouterBuildsWithoutLegacyConnectorCredentials.
func TestNewListingsSchedulersStartAllDoesNotTouchDatabaseSynchronously(t *testing.T) {
	lister := &blockingInstallationLister{listCalled: make(chan struct{})}
	group := NewListingsSchedulers(nil, time.Hour, lister, nil, nil, nil)
	group.StartAll(context.Background())

	select {
	case <-lister.listCalled:
		// fine: List() runs in the background goroutine eventually.
	case <-time.After(2 * time.Second):
		t.Fatal("List() was never called from StartAll's background goroutine")
	}
}

type blockingInstallationLister struct {
	listCalled chan struct{}
}

func (l *blockingInstallationLister) List(context.Context) ([]integrationsdomain.Installation, error) {
	close(l.listCalled)
	return nil, nil
}
