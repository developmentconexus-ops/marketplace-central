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

	calls      int
	mergeState map[string]struct{}
	inserts    int
	deletes    int
}

func (f *fakeSnapshotSyncer) SyncLatestCompletedSnapshot(_ context.Context, tenantID string, source erpdomain.ImportSource) (int, error) {
	f.calls++
	f.tenant = tenantID
	f.source = source
	if f.err != nil {
		return 0, f.err
	}
	for _, product := range []string{"100", "200"} {
		if _, exists := f.mergeState[product]; !exists {
			f.mergeState[product] = struct{}{}
			f.inserts++
		}
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
	syncer := &fakeSnapshotSyncer{result: 7, mergeState: make(map[string]struct{})}
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
	syncer := &fakeSnapshotSyncer{err: wantErr, mergeState: make(map[string]struct{})}
	adapter := NewXlsxAdapter(&recordingReader{}, syncer, "tenant-a")

	got, err := adapter.Sync(context.Background())
	if err == nil || !errors.Is(err, wantErr) {
		t.Fatalf("Sync error = %v, want wrapped merge error", err)
	}
	if got != (readports.SyncResult{}) {
		t.Fatalf("SyncResult = %+v, want zero result on error", got)
	}
}

func TestXlsxAdapterSyncIsIdempotentThroughKeepAbsentMerge(t *testing.T) {
	syncer := &fakeSnapshotSyncer{result: 2, mergeState: make(map[string]struct{})}
	adapter := NewXlsxAdapter(&recordingReader{}, syncer, "tenant-a")

	first, err := adapter.Sync(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	stateAfterFirst := len(syncer.mergeState)
	second, err := adapter.Sync(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if second != first {
		t.Fatalf("second SyncResult = %+v, first = %+v", second, first)
	}
	if len(syncer.mergeState) != stateAfterFirst || len(syncer.mergeState) != 2 {
		t.Fatalf("merge state changed or duplicated: size after first=%d, final=%d", stateAfterFirst, len(syncer.mergeState))
	}
	if syncer.inserts != 2 || syncer.deletes != 0 || syncer.calls != 2 {
		t.Fatalf("merge calls = %d, inserts = %d, deletes = %d; want 2 calls, 2 initial inserts, no deletes", syncer.calls, syncer.inserts, syncer.deletes)
	}
}

func TestXlsxAdapterKind(t *testing.T) {
	adapter := NewXlsxAdapter(&recordingReader{}, &fakeSnapshotSyncer{mergeState: make(map[string]struct{})}, "tenant-a")

	if got := adapter.Kind(); got != sourcekind.UploadSnapshot {
		t.Fatalf("Kind() = %q, want %q", got, sourcekind.UploadSnapshot)
	}
}
