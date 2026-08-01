package listings

import (
	"context"
	"errors"
	"testing"

	listingsports "marketplace-central/apps/server_core/internal/modules/listings/ports"
	mutationsdomain "marketplace-central/apps/server_core/internal/modules/mutations/domain"
	mutationsports "marketplace-central/apps/server_core/internal/modules/mutations/ports"
)

type pullerStub struct {
	calls   int
	account listingsports.InstallationAccount
	err     error
}

func (s *pullerStub) Pull(_ context.Context, account listingsports.InstallationAccount) error {
	s.calls++
	s.account = account
	return s.err
}

func TestResyncWriterDelegatesToPullerWithAccount(t *testing.T) {
	account := listingsports.InstallationAccount{TenantID: "tenant-1", InstallationID: "installation-1", ProviderCode: "mercado_livre", ProviderAccountID: "seller-1"}
	puller := &pullerStub{}
	writer := NewResyncWriter(puller, account)

	outcome, err := writer.Apply(context.Background(), mutationsports.WriteItem{ProtocolType: mutationsdomain.ProtocolTypeListingResync, InstallationID: "installation-1", ListingID: "MLB-1"})
	if err != nil || outcome.Failure != nil {
		t.Fatalf("outcome=%+v err=%v", outcome, err)
	}
	if puller.calls != 1 {
		t.Fatalf("puller calls = %d, want 1", puller.calls)
	}
	if puller.account.TenantID != "tenant-1" || puller.account.InstallationID != "installation-1" || puller.account.ProviderAccountID != "seller-1" {
		t.Fatalf("puller.account = %+v", puller.account)
	}
}

func TestResyncWriterPropagatesPullerFailure(t *testing.T) {
	account := listingsports.InstallationAccount{TenantID: "tenant-1", InstallationID: "installation-1", ProviderCode: "mercado_livre"}
	puller := &pullerStub{err: errors.New("provider unavailable")}
	writer := NewResyncWriter(puller, account)

	_, err := writer.Apply(context.Background(), mutationsports.WriteItem{ProtocolType: mutationsdomain.ProtocolTypeListingResync, InstallationID: "installation-1", ListingID: "MLB-1"})
	if err == nil {
		t.Fatal("Apply() error = nil, want the puller failure")
	}
}

func TestResyncWriterRejectsCrossInstallationBeforePull(t *testing.T) {
	puller := &pullerStub{}
	writer := NewResyncWriter(puller, listingsports.InstallationAccount{TenantID: "tenant-1", InstallationID: "installation-1", ProviderCode: "mercado_livre"})
	_, err := writer.Apply(context.Background(), mutationsports.WriteItem{ProtocolType: mutationsdomain.ProtocolTypeListingResync, InstallationID: "other", ListingID: "MLB-1"})
	if err == nil || puller.calls != 0 {
		t.Fatalf("err=%v calls=%d", err, puller.calls)
	}
}

func TestResyncWriterNotConfigured(t *testing.T) {
	writer := NewResyncWriter(nil, listingsports.InstallationAccount{TenantID: "tenant-1", InstallationID: "installation-1"})
	_, err := writer.Apply(context.Background(), mutationsports.WriteItem{ProtocolType: mutationsdomain.ProtocolTypeListingResync, InstallationID: "installation-1", ListingID: "MLB-1"})
	if err == nil {
		t.Fatal("Apply() error = nil, want configuration error")
	}
}
