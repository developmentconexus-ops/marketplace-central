package xlsx

import (
	"context"
	"errors"
	"testing"

	erpdomain "marketplace-central/apps/server_core/internal/modules/erp_import/domain"
	readdomain "marketplace-central/apps/server_core/internal/modules/internal_read/domain"
	readports "marketplace-central/apps/server_core/internal/modules/internal_read/ports"
	"marketplace-central/apps/server_core/internal/modules/sourcekind"
)

var _ readports.ProductSourceAdapter = (*XlsxAdapter)(nil)

type recordingReader struct {
	readports.Reader
	currentPriceCalled bool
}

func (r *recordingReader) GetCurrentPrice(context.Context, readports.CurrentPriceInput) (readdomain.CurrentPrice, error) {
	r.currentPriceCalled = true
	return readdomain.CurrentPrice{}, nil
}

type fakeSnapshotSyncer struct {
	tenant string
	source erpdomain.ImportSource
	result int
	err    error

	calls int
}

func (f *fakeSnapshotSyncer) SyncLatestCompletedSnapshot(_ context.Context, tenantID string, source erpdomain.ImportSource) (int, error) {
	f.calls++
	f.tenant = tenantID
	f.source = source
	if f.err != nil {
		return 0, f.err
	}
	return f.result, nil
}

func TestXlsxAdapterSatisfiesProductSourceAdapter(t *testing.T) {
	reader := &recordingReader{}
	adapter := NewXlsxAdapter(reader, &fakeSnapshotSyncer{}, "tenant-a")

	if _, err := adapter.GetCurrentPrice(context.Background(), readports.CurrentPriceInput{}); err != nil {
		t.Fatal(err)
	}
	if !reader.currentPriceCalled {
		t.Fatal("embedded reader method did not delegate to the supplied reader")
	}
}

func TestXlsxAdapterSyncForwardsSourceAndMapsProcessed(t *testing.T) {
	syncer := &fakeSnapshotSyncer{result: 7}
	adapter := NewXlsxAdapter(&recordingReader{}, syncer, "tenant-a")

	got, err := adapter.Sync(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if got.Processed != 7 || got.Errors != 0 {
		t.Fatalf("SyncResult = %+v, want Processed=7 and Errors=0", got)
	}
	if syncer.tenant != "tenant-a" || syncer.source != erpdomain.SourceXLSX {
		t.Fatalf("sync call = tenant %q, source %q", syncer.tenant, syncer.source)
	}
}

func TestXlsxAdapterSyncWrapsErrorWithoutPartialResult(t *testing.T) {
	wantErr := errors.New("merge failed")
	syncer := &fakeSnapshotSyncer{err: wantErr}
	adapter := NewXlsxAdapter(&recordingReader{}, syncer, "tenant-a")

	got, err := adapter.Sync(context.Background())
	if err == nil || !errors.Is(err, wantErr) {
		t.Fatalf("Sync error = %v, want wrapped merge error", err)
	}
	if got != (readports.SyncResult{}) {
		t.Fatalf("SyncResult = %+v, want zero result on error", got)
	}
}

func TestXlsxAdapterSyncDelegatesToSnapshotSyncer(t *testing.T) {
	syncer := &fakeSnapshotSyncer{result: 2}
	adapter := NewXlsxAdapter(&recordingReader{}, syncer, "tenant-a")

	got, err := adapter.Sync(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if got.Processed != 2 || got.Errors != 0 {
		t.Fatalf("SyncResult = %+v, want Processed=2 and Errors=0", got)
	}
	if syncer.calls != 1 || syncer.tenant != "tenant-a" || syncer.source != erpdomain.SourceXLSX {
		t.Fatalf("sync calls=%d tenant=%q source=%q", syncer.calls, syncer.tenant, syncer.source)
	}
}

func TestXlsxAdapterKind(t *testing.T) {
	adapter := NewXlsxAdapter(&recordingReader{}, &fakeSnapshotSyncer{}, "tenant-a")

	if got := adapter.Kind(); got != sourcekind.UploadSnapshot {
		t.Fatalf("Kind() = %q, want %q", got, sourcekind.UploadSnapshot)
	}
}
