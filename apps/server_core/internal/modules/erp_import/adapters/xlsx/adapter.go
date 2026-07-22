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
}

func NewXlsxAdapter(reader readports.Reader, syncer snapshotSyncer, tenantID string) *XlsxAdapter {
	return &XlsxAdapter{
		Reader:   reader,
		syncer:   syncer,
		tenantID: tenantID,
	}
}

func (a *XlsxAdapter) Sync(ctx context.Context) (readports.SyncResult, error) {
	processed, err := a.syncer.SyncLatestCompletedSnapshot(ctx, a.tenantID, erpdomain.SourceXLSX)
	if err != nil {
		return readports.SyncResult{}, fmt.Errorf("xlsx sync: %w", err)
	}
	return readports.SyncResult{Processed: processed}, nil
}

func (a *XlsxAdapter) Kind() sourcekind.SourceKind {
	return sourcekind.UploadSnapshot
}

var _ readports.ProductSourceAdapter = (*XlsxAdapter)(nil)
