package ports

import (
	"context"
	"fmt"
	"io"

	"marketplace-central/apps/server_core/internal/modules/erp_import/domain"
)

type FileError struct {
	Code   domain.IssueCode
	Column string
	Detail string
}

func (e *FileError) Error() string {
	if e.Column == "" {
		return fmt.Sprintf("%s: %s", e.Code, e.Detail)
	}
	return fmt.Sprintf("%s: %s (%s)", e.Code, e.Detail, e.Column)
}

type Parser interface {
	Parse(ctx context.Context, source io.Reader) ([]domain.NormalizedRow, error)
}
