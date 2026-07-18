package domain

import (
	"regexp"
	"strings"

	catalogdomain "marketplace-central/apps/server_core/internal/modules/catalog/domain"
)

var (
	decimalPattern = regexp.MustCompile(`^\+?(?:[0-9]+(?:[.,][0-9]+)?|[.,][0-9]+)$`)
	integerPattern = regexp.MustCompile(`^[0-9]+$`)
	ncmPattern     = regexp.MustCompile(`^[0-9]{8}$`)
)

type RowValidationResult struct {
	Row    *NormalizedRow
	Issues []Issue
}

// ValidateRow validates one normalized row. seenCodprod contains earlier rows
// in the same file and is read-only.
func ValidateRow(row NormalizedRow, seenCodprod map[string]struct{}) RowValidationResult {
	row.Codprod = strings.TrimSpace(row.Codprod)
	row.Descrprod = strings.TrimSpace(row.Descrprod)
	issues := make([]Issue, 0, 4)

	if row.Codprod == "" {
		issues = append(issues, issue("CODPROD", Rejection, CodeEmptyCodprod, "codprod is required", row.Codprod))
	} else if _, duplicate := seenCodprod[row.Codprod]; duplicate {
		issues = append(issues, issue("CODPROD", Rejection, CodeDuplicateCodprod, "codprod is duplicated in the file", row.Codprod))
	}

	if row.Descrprod == "" {
		issues = append(issues, issue("DESCRPROD", Rejection, CodeEmptyDescrprod, "descrprod is required", row.Descrprod))
	}

	if custo, valid := normalizeDecimal(row.Custo); valid {
		row.Custo = custo
	} else {
		value := string(row.Custo)
		issues = append(issues, issue("CUSTO", Rejection, CodeInvalidCusto, "custo must be a positive decimal", value))
	}

	if stock, valid := normalizeInteger(row.StockPhysical); valid {
		row.StockPhysical = stock
	} else {
		issues = append(issues, issue("ESTOQUE_FISICO", Rejection, CodeInvalidEstoque, "estoque_fisico must be a non-negative integer", row.StockPhysical))
	}

	row.StockReserved = normalizeOptionalInteger(row.StockReserved)
	row.Refforn = normalizeOptionalString(row.Refforn)
	row.Marca = normalizeOptionalString(row.Marca)

	if value := normalizeOptionalString(row.EAN); value != nil {
		if catalogdomain.IsValidGTIN(*value) {
			row.EAN = value
		} else {
			issues = append(issues, issue("EAN", Warning, CodeInvalidEAN, "ean is not a valid GTIN", *value))
			row.EAN = nil
		}
	} else {
		row.EAN = nil
	}

	if value := normalizeOptionalString(row.NCM); value != nil {
		if ncmPattern.MatchString(*value) {
			row.NCM = value
		} else {
			issues = append(issues, issue("NCM", Warning, CodeInvalidNCM, "ncm must contain exactly 8 digits", *value))
			row.NCM = nil
		}
	} else {
		row.NCM = nil
	}

	for _, current := range issues {
		if current.Kind == Rejection {
			return RowValidationResult{Issues: issues}
		}
	}
	return RowValidationResult{Row: &row, Issues: issues}
}

// ValidateRows validates a complete file snapshot and excludes rejected rows.
func ValidateRows(rows []NormalizedRow) (accepted []NormalizedRow, issues []Issue, rejected int) {
	seen := make(map[string]struct{}, len(rows))
	for index, row := range rows {
		codprod := strings.TrimSpace(row.Codprod)
		prior := make(map[string]struct{}, 1)
		if codprod != "" {
			if _, exists := seen[codprod]; exists {
				prior[codprod] = struct{}{}
			}
			seen[codprod] = struct{}{}
		}

		result := ValidateRow(row, prior)
		for issueIndex := range result.Issues {
			result.Issues[issueIndex].Row = index + 1
		}
		issues = append(issues, result.Issues...)
		if result.Row == nil {
			rejected++
			continue
		}
		accepted = append(accepted, *result.Row)
	}
	return accepted, issues, rejected
}

func normalizeDecimal(value Decimal) (Decimal, bool) {
	text := strings.TrimSpace(string(value))
	if !decimalPattern.MatchString(text) {
		return "", false
	}
	text = strings.TrimPrefix(text, "+")
	text = strings.ReplaceAll(text, ",", ".")
	if strings.Trim(text, "0.") == "" {
		return "", false
	}
	if strings.HasPrefix(text, ".") {
		text = "0" + text
	}
	return Decimal(text), true
}

func normalizeInteger(value string) (string, bool) {
	text := strings.TrimSpace(value)
	if !integerPattern.MatchString(text) {
		return "", false
	}
	text = strings.TrimLeft(text, "0")
	if text == "" {
		text = "0"
	}
	return text, true
}

func normalizeOptionalInteger(value *string) *string {
	if value == nil {
		return nil
	}
	normalized, valid := normalizeInteger(*value)
	if !valid {
		return nil
	}
	return &normalized
}

func normalizeOptionalString(value *string) *string {
	if value == nil {
		return nil
	}
	text := strings.TrimSpace(*value)
	if text == "" {
		return nil
	}
	return &text
}

func issue(column string, kind IssueKind, code IssueCode, detail, offending string) Issue {
	return Issue{
		Column:         &column,
		Kind:           kind,
		Code:           code,
		Detail:         detail,
		OffendingValue: &offending,
	}
}
