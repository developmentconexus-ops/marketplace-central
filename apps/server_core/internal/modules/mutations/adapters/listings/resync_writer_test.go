package listings

import (
	"context"
	"errors"
	"sort"
	"testing"

	listingsports "marketplace-central/apps/server_core/internal/modules/listings/ports"
	mutationsdomain "marketplace-central/apps/server_core/internal/modules/mutations/domain"
	mutationsports "marketplace-central/apps/server_core/internal/modules/mutations/ports"
)

// ingestorStub is a fake ListingIngestor that records every provider listing
// id it was called with, so tests can assert exact per-item discrimination
// (not just a call count that could hide a shared/collapsed call).
type ingestorStub struct {
	calls []string
	err   error
}

func (s *ingestorStub) IngestListing(_ context.Context, providerListingID string) error {
	s.calls = append(s.calls, providerListingID)
	return s.err
}

func TestResyncWriterDelegatesToIngestorForTheNamedListing(t *testing.T) {
	account := listingsports.InstallationAccount{TenantID: "tenant-1", InstallationID: "installation-1", ProviderCode: "mercado_livre", ProviderAccountID: "seller-1"}
	ingestor := &ingestorStub{}
	writer := NewResyncWriter(ingestor, account)

	outcome, err := writer.Apply(context.Background(), mutationsports.WriteItem{ProtocolType: mutationsdomain.ProtocolTypeListingResync, InstallationID: "installation-1", ListingID: "installation-1~MLB-1~0"})
	if err != nil || outcome.Failure != nil {
		t.Fatalf("outcome=%+v err=%v", outcome, err)
	}
	if len(ingestor.calls) != 1 || ingestor.calls[0] != "MLB-1" {
		t.Fatalf("ingestor.calls = %+v, want exactly [MLB-1]", ingestor.calls)
	}
}

func TestResyncWriterPropagatesIngestorFailure(t *testing.T) {
	account := listingsports.InstallationAccount{TenantID: "tenant-1", InstallationID: "installation-1", ProviderCode: "mercado_livre"}
	ingestor := &ingestorStub{err: errors.New("provider unavailable")}
	writer := NewResyncWriter(ingestor, account)

	_, err := writer.Apply(context.Background(), mutationsports.WriteItem{ProtocolType: mutationsdomain.ProtocolTypeListingResync, InstallationID: "installation-1", ListingID: "installation-1~MLB-1~0"})
	if err == nil {
		t.Fatal("Apply() error = nil, want the ingestor failure")
	}
}

func TestResyncWriterRejectsCrossInstallationBeforeIngest(t *testing.T) {
	ingestor := &ingestorStub{}
	writer := NewResyncWriter(ingestor, listingsports.InstallationAccount{TenantID: "tenant-1", InstallationID: "installation-1", ProviderCode: "mercado_livre"})
	_, err := writer.Apply(context.Background(), mutationsports.WriteItem{ProtocolType: mutationsdomain.ProtocolTypeListingResync, InstallationID: "other", ListingID: "other~MLB-1~0"})
	if err == nil || len(ingestor.calls) != 0 {
		t.Fatalf("err=%v calls=%+v", err, ingestor.calls)
	}
}

func TestResyncWriterNotConfigured(t *testing.T) {
	writer := NewResyncWriter(nil, listingsports.InstallationAccount{TenantID: "tenant-1", InstallationID: "installation-1"})
	_, err := writer.Apply(context.Background(), mutationsports.WriteItem{ProtocolType: mutationsdomain.ProtocolTypeListingResync, InstallationID: "installation-1", ListingID: "installation-1~MLB-1~0"})
	if err == nil {
		t.Fatal("Apply() error = nil, want configuration error")
	}
}

func TestResyncWriterRejectsMalformedListingID(t *testing.T) {
	ingestor := &ingestorStub{}
	writer := NewResyncWriter(ingestor, listingsports.InstallationAccount{TenantID: "tenant-1", InstallationID: "installation-1"})
	_, err := writer.Apply(context.Background(), mutationsports.WriteItem{ProtocolType: mutationsdomain.ProtocolTypeListingResync, InstallationID: "installation-1", ListingID: "not-a-canonical-id"})
	if err == nil || len(ingestor.calls) != 0 {
		t.Fatalf("err=%v calls=%+v, want a parse error and zero ingest calls", err, ingestor.calls)
	}
}

// TestResyncWriterBatchCallsIngestorOncePerDistinctListing is the proof this
// bug fix exists for: a batch of N>1 resync WriteItems for distinct listings
// (the same shape mutations/application.Poller.Pass feeds one item at a time,
// in a loop, to WriterPort.Apply) must reach the per-item ListingIngestor
// exactly once per distinct provider listing id — never a single shared
// whole-catalog pull, and never a collapsed call that silently drops the
// other N-1 items.
func TestResyncWriterBatchCallsIngestorOncePerDistinctListing(t *testing.T) {
	account := listingsports.InstallationAccount{TenantID: "tenant-1", InstallationID: "installation-1", ProviderCode: "mercado_livre", ProviderAccountID: "seller-1"}
	ingestor := &ingestorStub{}
	writer := NewResyncWriter(ingestor, account)

	listingIDs := []string{
		"installation-1~MLB-1~0",
		"installation-1~MLB-2~0",
		"installation-1~MLB-3~0",
		"installation-1~MLB-4~0",
	}
	for _, listingID := range listingIDs {
		outcome, err := writer.Apply(context.Background(), mutationsports.WriteItem{ProtocolType: mutationsdomain.ProtocolTypeListingResync, InstallationID: "installation-1", ListingID: listingID})
		if err != nil || outcome.Failure != nil {
			t.Fatalf("Apply(%s) outcome=%+v err=%v", listingID, outcome, err)
		}
	}

	wantProviderIDs := []string{"MLB-1", "MLB-2", "MLB-3", "MLB-4"}
	if len(ingestor.calls) != len(wantProviderIDs) {
		t.Fatalf("ingestor.calls = %+v, want %d calls (one per distinct listing)", ingestor.calls, len(wantProviderIDs))
	}
	gotSorted := append([]string(nil), ingestor.calls...)
	sort.Strings(gotSorted)
	wantSorted := append([]string(nil), wantProviderIDs...)
	sort.Strings(wantSorted)
	for i := range wantSorted {
		if gotSorted[i] != wantSorted[i] {
			t.Fatalf("ingestor.calls (sorted) = %+v, want %+v", gotSorted, wantSorted)
		}
	}
}
