package application

import (
	"context"
	"errors"
	"strings"

	"marketplace-central/apps/server_core/internal/modules/product_links/ports"
)

var ErrSummaryInstallationRequired = errors.New("PRODUCT_LINKS_INSTALLATION_REQUIRED")
var ErrSummaryReaderNotConfigured = errors.New("PRODUCT_LINKS_SUMMARY_READER_NOT_CONFIGURED")

type SummaryService struct {
	reader ports.SummaryReader
}

func NewSummaryService(reader ports.SummaryReader) SummaryService {
	return SummaryService{reader: reader}
}

func (s SummaryService) Summary(ctx context.Context, installationID string) (ports.LinkageSummary, error) {
	if strings.TrimSpace(installationID) == "" {
		return ports.LinkageSummary{}, ErrSummaryInstallationRequired
	}
	if s.reader == nil {
		return ports.LinkageSummary{}, ErrSummaryReaderNotConfigured
	}
	return s.reader.GetLinkageSummary(ctx, strings.TrimSpace(installationID))
}
