package ports

import (
	"errors"
	"strings"
	"time"

	"marketplace-central/apps/server_core/internal/modules/integrations/domain"
)

var ErrInvalidRunListQuery = errors.New("invalid_run_list_query")

type RunListQuery struct {
	InstallationID string
	Cursor         RunCursor
	Limit          int
	Module         string
	Status         domain.OperationRunStatus
}

func (q RunListQuery) Normalize() (RunListQuery, error) {
	q.InstallationID = strings.TrimSpace(q.InstallationID)
	q.Module = strings.TrimSpace(q.Module)
	if q.InstallationID == "" || q.Limit < 0 {
		return RunListQuery{}, ErrInvalidRunListQuery
	}
	if q.Limit == 0 {
		q.Limit = 20
	}
	if q.Limit > 100 {
		q.Limit = 100
	}
	return q, nil
}

type RunReadModel struct {
	OperationRunID      string
	InstallationID      string
	Module              string
	Status              domain.OperationRunStatus
	ResultCode          string
	FailureCode         string
	TranslatedErrorCode string
	AttemptCount        int
	DurationMs          int64
	StartedAt           *time.Time
	FinishedAt          *time.Time
}

type RunPage struct {
	Items      []RunReadModel
	NextCursor *RunCursor
}

type LatestRunByModule struct {
	OperationType      string
	LatestAttemptedAt  *time.Time
	LatestSuccessfulAt *time.Time
}
