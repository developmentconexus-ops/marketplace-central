package application

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	ordersdomain "marketplace-central/apps/server_core/internal/modules/orders/domain"
	"marketplace-central/apps/server_core/internal/modules/orders/ports"
)

type stubOrderSource struct {
	items []ordersdomain.OrderIngestionSnapshot
	err   error
}

func (s stubOrderSource) ListOrders(context.Context, ports.ListOrdersInput) ([]ordersdomain.OrderIngestionSnapshot, error) {
	return s.items, s.err
}

// recordingSource is a package-local ports.OrderSource fake used to assert that Import passes
// Offset/UpdatedAfter through to the enumeration port unchanged — legitimate because OrderSource
// is an application-domain port and what's under test is parameter pass-through, not provider
// integration.
type recordingSource struct{ got ports.ListOrdersInput }

func (s *recordingSource) ListOrders(_ context.Context, in ports.ListOrdersInput) ([]ordersdomain.OrderIngestionSnapshot, error) {
	s.got = in
	return nil, nil
}

// stubLinkReader is a package-wide test double for ports.LinkReader, kept in this file (as
// before the F-02 rewrite) since ingest_service_test.go also needs it.
type stubLinkReader struct {
	links map[string]ordersdomain.ItemLink
	err   error
}

func (s stubLinkReader) ResolveLinks(context.Context, string, []ordersdomain.ListingIdentity) (map[string]ordersdomain.ItemLink, error) {
	return s.links, s.err
}

// stubOrderIngestor fakes ports.OrderIngestor: errFor names which provider_order_id should fail
// (and with what error), recording every call so tests can assert exactly which ids Import
// enumerated and delegated.
type stubOrderIngestor struct {
	errFor map[string]error
	calls  []string
}

func (s *stubOrderIngestor) IngestOrder(_ context.Context, _ string, providerOrderID string) error {
	s.calls = append(s.calls, providerOrderID)
	if s.errFor != nil {
		if err, ok := s.errFor[providerOrderID]; ok {
			return err
		}
	}
	return nil
}

func snapshotWithID(providerOrderID string) ordersdomain.OrderIngestionSnapshot {
	return ordersdomain.OrderIngestionSnapshot{
		ProviderCode:    "mercado_livre",
		ProviderOrderID: providerOrderID,
		FetchedAt:       time.Date(2026, 7, 9, 12, 0, 0, 0, time.UTC),
	}
}

// floatPtr is a package-wide test helper (also used by enrich_service_test.go /
// enrich_link_refresh_test.go — F-03's fixtures) kept here since this file owned it before the
// F-02 rewrite.
func floatPtr(v float64) *float64 { return &v }
func intPtr(v int) *int           { return &v }

func TestImportServiceIngestsEachEnumeratedOrder(t *testing.T) {
	source := stubOrderSource{items: []ordersdomain.OrderIngestionSnapshot{
		snapshotWithID("2001"), snapshotWithID("2002"), snapshotWithID("2003"),
	}}
	ingestor := &stubOrderIngestor{}
	service := NewImportService(ImportServiceConfig{Source: source, Ingestor: ingestor})

	result, err := service.Import(context.Background(), ImportOrdersInput{InstallationID: "inst-1", Limit: 3})
	if err != nil {
		t.Fatalf("Import() error = %v", err)
	}
	if result.ImportedCount != 3 || result.SkippedCount != 0 {
		t.Fatalf("counts = imported %d skipped %d, want 3/0", result.ImportedCount, result.SkippedCount)
	}
	want := []string{"2001", "2002", "2003"}
	if len(ingestor.calls) != len(want) {
		t.Fatalf("ingestor calls = %v, want %v", ingestor.calls, want)
	}
	for i, id := range want {
		if ingestor.calls[i] != id {
			t.Fatalf("ingestor calls = %v, want %v", ingestor.calls, want)
		}
	}
}

// TestImportServiceCountsOrderUnavailableAsSkipWithoutFailingRun is the brief's named scenario
// (e): a 403/404 third-party order classified by IngestOrder as domain.ErrOrderUnavailable must
// be counted as a skip, and the run must still succeed and continue ingesting the rest.
func TestImportServiceCountsOrderUnavailableAsSkipWithoutFailingRun(t *testing.T) {
	source := stubOrderSource{items: []ordersdomain.OrderIngestionSnapshot{
		snapshotWithID("2001"), snapshotWithID("2002"),
	}}
	ingestor := &stubOrderIngestor{errFor: map[string]error{
		"2001": fmt.Errorf("%w: installation=inst-1 provider_order_id=2001: %w", ordersdomain.ErrOrderUnavailable, errors.New("provider unauthorized")),
	}}
	service := NewImportService(ImportServiceConfig{Source: source, Ingestor: ingestor})

	result, err := service.Import(context.Background(), ImportOrdersInput{InstallationID: "inst-1", Limit: 2})
	if err != nil {
		t.Fatalf("Import() error = %v, want run to succeed despite one unavailable order", err)
	}
	if result.ImportedCount != 1 || result.SkippedCount != 1 {
		t.Fatalf("counts = imported %d skipped %d, want 1/1", result.ImportedCount, result.SkippedCount)
	}
	if len(ingestor.calls) != 2 {
		t.Fatalf("ingestor calls = %v, want both orders attempted", ingestor.calls)
	}
}

// TestImportServiceCountsAnyIngestFailureAsSkip: IC-06 "erro por item no ingest: registra, segue
// o batch" is not scoped to the 403/404 case alone — any per-order failure (a transient shipment
// read, a DB blip) must not abort the whole run either.
func TestImportServiceCountsAnyIngestFailureAsSkip(t *testing.T) {
	source := stubOrderSource{items: []ordersdomain.OrderIngestionSnapshot{
		snapshotWithID("3001"), snapshotWithID("3002"),
	}}
	ingestor := &stubOrderIngestor{errFor: map[string]error{
		"3001": errors.New("provider temporarily unavailable"),
	}}
	service := NewImportService(ImportServiceConfig{Source: source, Ingestor: ingestor})

	result, err := service.Import(context.Background(), ImportOrdersInput{InstallationID: "inst-1", Limit: 2})
	if err != nil {
		t.Fatalf("Import() error = %v, want run to succeed despite one transient failure", err)
	}
	if result.ImportedCount != 1 || result.SkippedCount != 1 {
		t.Fatalf("counts = imported %d skipped %d, want 1/1", result.ImportedCount, result.SkippedCount)
	}
}

func TestImportServiceSkipsBlankProviderOrderID(t *testing.T) {
	source := stubOrderSource{items: []ordersdomain.OrderIngestionSnapshot{
		snapshotWithID(""), snapshotWithID("4001"),
	}}
	ingestor := &stubOrderIngestor{}
	service := NewImportService(ImportServiceConfig{Source: source, Ingestor: ingestor})

	result, err := service.Import(context.Background(), ImportOrdersInput{InstallationID: "inst-1", Limit: 2})
	if err != nil {
		t.Fatalf("Import() error = %v", err)
	}
	if result.ImportedCount != 1 || result.SkippedCount != 0 {
		t.Fatalf("counts = imported %d skipped %d, want 1/0 (blank id neither imported nor skipped)", result.ImportedCount, result.SkippedCount)
	}
	if len(ingestor.calls) != 1 || ingestor.calls[0] != "4001" {
		t.Fatalf("ingestor calls = %v, want only the non-blank id", ingestor.calls)
	}
}

func TestImportServiceFailsRunWhenSourceListOrdersErrors(t *testing.T) {
	source := stubOrderSource{err: errors.New("provider list failed")}
	ingestor := &stubOrderIngestor{}
	service := NewImportService(ImportServiceConfig{Source: source, Ingestor: ingestor})

	if _, err := service.Import(context.Background(), ImportOrdersInput{InstallationID: "inst-1", Limit: 2}); err == nil {
		t.Fatal("Import() error = nil, want enumeration failure to fail the run")
	}
	if len(ingestor.calls) != 0 {
		t.Fatalf("ingestor calls = %v, want none when enumeration itself fails", ingestor.calls)
	}
}

func TestImportServiceRequiresInstallationID(t *testing.T) {
	service := NewImportService(ImportServiceConfig{Source: stubOrderSource{}, Ingestor: &stubOrderIngestor{}})
	if _, err := service.Import(context.Background(), ImportOrdersInput{}); err == nil {
		t.Fatal("Import() error = nil, want missing installation id to fail")
	}
}

// TestImportPassesWindowAndOffsetToSource is F-00 Task 1's pass-through criterion: the narrow
// (installationID, limit) signature already cost a feature once (three intermediate layers threw
// UpdatedAfter away before it reached the Mercado Livre adapter, which has always known how to
// translate it). Import must forward Offset and UpdatedAfter to ports.OrderSource unchanged.
func TestImportPassesWindowAndOffsetToSource(t *testing.T) {
	updatedAfter := time.Date(2026, 8, 1, 12, 0, 0, 0, time.UTC)
	src := &recordingSource{}
	service := NewImportService(ImportServiceConfig{Source: src, Ingestor: &stubOrderIngestor{}})

	_, err := service.Import(context.Background(), ImportOrdersInput{
		InstallationID: "inst-1",
		Limit:          50,
		Offset:         100,
		UpdatedAfter:   &updatedAfter,
	})
	if err != nil {
		t.Fatalf("Import() error = %v", err)
	}
	if src.got.InstallationID != "inst-1" || src.got.Limit != 50 {
		t.Fatalf("installation_id/limit not passed through: %+v", src.got)
	}
	if src.got.Offset != 100 {
		t.Fatalf("offset: quero 100, recebi %d", src.got.Offset)
	}
	if src.got.UpdatedAfter == nil || !src.got.UpdatedAfter.Equal(updatedAfter) {
		t.Fatalf("updated_after: quero %s, recebi %v", updatedAfter, src.got.UpdatedAfter)
	}
}
