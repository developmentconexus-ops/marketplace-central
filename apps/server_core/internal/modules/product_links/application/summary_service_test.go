package application

import (
	"context"
	"errors"
	"testing"

	"marketplace-central/apps/server_core/internal/modules/product_links/ports"
)

type fakeSummaryReader struct {
	summary        ports.LinkageSummary
	err            error
	installationID string
}

func (f *fakeSummaryReader) GetLinkageSummary(_ context.Context, installationID string) (ports.LinkageSummary, error) {
	f.installationID = installationID
	return f.summary, f.err
}

func TestSummaryServiceReturnsLinkageCounts(t *testing.T) {
	reader := &fakeSummaryReader{summary: ports.LinkageSummary{PendingLinks: 4, MissingGTIN: 3}}
	service := NewSummaryService(reader)

	got, err := service.Summary(context.Background(), "installation-1")
	if err != nil {
		t.Fatalf("Summary() error = %v", err)
	}
	if got != reader.summary {
		t.Fatalf("Summary() = %+v, want %+v", got, reader.summary)
	}
	if reader.installationID != "installation-1" {
		t.Fatalf("installation ID = %q, want installation-1", reader.installationID)
	}
}

func TestSummaryServicePropagatesReaderError(t *testing.T) {
	wantErr := errors.New("summary source unavailable")
	service := NewSummaryService(&fakeSummaryReader{err: wantErr})

	got, err := service.Summary(context.Background(), "installation-1")
	if !errors.Is(err, wantErr) {
		t.Fatalf("Summary() error = %v, want %v", err, wantErr)
	}
	if got != (ports.LinkageSummary{}) {
		t.Fatalf("Summary() = %+v, want zero result on error", got)
	}
}
