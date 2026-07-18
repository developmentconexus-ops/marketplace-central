package postgres

import (
	"testing"

	"marketplace-central/apps/server_core/internal/modules/erp_import/domain"
)

// issueCounts must count distinct rejected ROWS (a single row can carry several
// rejection issues) and warning ISSUES. A regression to len(rejectionIssues)
// would fail the two-rejections-on-one-row case below.
func TestIssueCountsDistinctRejectedRows(t *testing.T) {
	issues := []domain.Issue{
		{Row: 1, Kind: domain.Rejection, Code: domain.CodeEmptyCodprod},
		{Row: 1, Kind: domain.Rejection, Code: domain.CodeInvalidCusto},
		{Row: 2, Kind: domain.Rejection, Code: domain.CodeEmptyDescrprod},
		{Row: 3, Kind: domain.Warning, Code: domain.CodeInvalidEAN},
		{Row: 3, Kind: domain.Warning, Code: domain.CodeInvalidNCM},
	}
	rejected, warnings := issueCounts(issues)
	if rejected != 2 {
		t.Fatalf("rejected rows = %d, want 2 (rows 1 and 2, row 1 counted once)", rejected)
	}
	if warnings != 2 {
		t.Fatalf("warning issues = %d, want 2", warnings)
	}
}
