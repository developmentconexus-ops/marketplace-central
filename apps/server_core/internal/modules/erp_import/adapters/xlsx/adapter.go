package xlsx

import (
	"context"
	"fmt"

	erpdomain "marketplace-central/apps/server_core/internal/modules/erp_import/domain"
	readports "marketplace-central/apps/server_core/internal/modules/internal_read/ports"
	"marketplace-central/apps/server_core/internal/modules/sourcekind"
)

type snapshotSyncer interface {
	SyncLatestCompletedSnapshot(ctx context.Context, tenantID string, source erpdomain.ImportSource) (int, error)
}

type XlsxAdapter struct {
	readports.Reader
	syncer   snapshotSyncer
	tenantID string
	source   erpdomain.ImportSource
}

// NewXlsxAdapter builds the adapter over the ERP spreadsheet source.
func NewXlsxAdapter(reader readports.Reader, syncer snapshotSyncer, tenantID string) *XlsxAdapter {
	return NewUploadAdapter(reader, syncer, tenantID, erpdomain.SourceXLSX)
}

// NewUploadAdapter builds the adapter for an explicit upload source. Both
// upload sources (the ERP spreadsheet and the lenient client catalog) share the
// exact same mechanism — mirror the latest completed protocol of THAT source —
// and differ only in which protocol history they read, so the source is a
// parameter rather than a second near-identical adapter type.
func NewUploadAdapter(reader readports.Reader, syncer snapshotSyncer, tenantID string, source erpdomain.ImportSource) *XlsxAdapter {
	return &XlsxAdapter{
		Reader:   reader,
		syncer:   syncer,
		tenantID: tenantID,
		source:   source,
	}
}

func (a *XlsxAdapter) Sync(ctx context.Context) (readports.SyncResult, error) {
	processed, err := a.syncer.SyncLatestCompletedSnapshot(ctx, a.tenantID, a.source)
	if err != nil {
		return readports.SyncResult{}, fmt.Errorf("%s sync: %w", a.source, err)
	}
	return readports.SyncResult{Processed: processed}, nil
}

func (a *XlsxAdapter) Kind() sourcekind.SourceKind {
	return sourcekind.UploadSnapshot
}

var _ readports.ProductSourceAdapter = (*XlsxAdapter)(nil)
