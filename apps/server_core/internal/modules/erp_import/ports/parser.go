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

// Parser reads a workbook into normalized rows. The second return value carries
// file-level issues that are not tied to a single data row (e.g. a worksheet
// skipped for lacking a product header) — they are folded into the import
// report alongside the per-row validation issues.
type Parser interface {
	Parse(ctx context.Context, source io.Reader) ([]domain.NormalizedRow, []domain.Issue, error)
	ParseLenient(ctx context.Context, source io.Reader) ([]domain.NormalizedRow, []domain.Issue, error)
}
