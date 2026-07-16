package ports

import "context"

type LinkageSummary struct {
	PendingLinks int64
	MissingGTIN  int64
}

type SummaryReader interface {
	GetLinkageSummary(ctx context.Context, installationID string) (LinkageSummary, error)
}
