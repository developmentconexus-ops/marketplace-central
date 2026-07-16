package ports

import (
	"context"
	"time"
)

const ListingsRefreshOperationType = "listings_refresh"

type RefreshRun struct {
	OperationRunID      string
	InstallationID      string
	Status              string
	TranslatedErrorCode string
	StartedAt           *time.Time
	CompletedAt         *time.Time
}

type RefreshIntegrationGateway interface {
	InstallationAccount(context.Context, string) (InstallationAccount, bool, error)
	BeginRefresh(context.Context, RefreshRun) (RefreshRun, bool, error)
	RecordRefresh(context.Context, RefreshRun) error
}
