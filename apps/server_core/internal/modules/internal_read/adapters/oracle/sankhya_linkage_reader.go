package oracle

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"marketplace-central/apps/server_core/internal/modules/internal_read/domain"
	"marketplace-central/apps/server_core/internal/modules/internal_read/ports"
)

const (
	sankhyaOriginTOP      int64 = 313
	sankhyaDestinationTOP int64 = 306
	maxCandidateLimit           = 50
	maxCandidateLineLimit       = 500
	maxLineageLimit             = 500
)

var oracleLinkageIdentifier = regexp.MustCompile(`^[A-Z][A-Z0-9_$#]{0,29}$`)

type SankhyaLinkageConfig struct {
	Schema                  string
	HeaderFieldName         string
	ConfigurationRevision   string
	ExpectedOriginTOP       int64
	ExpectedDestinationTOP  int64
	UniquenessAttestationID string
	CandidateLimit          int
	CandidateLineLimit      int
	LineageLimit            int
}

type SankhyaLinkageReader struct {
	db     queryer
	config SankhyaLinkageConfig
}

func NewSankhyaLinkageReader(db queryer, config SankhyaLinkageConfig) *SankhyaLinkageReader {
	return &SankhyaLinkageReader{db: db, config: config}
}

func (r *SankhyaLinkageReader) ValidateConfiguration(ctx context.Context) error {
	if r == nil {
		return domain.NewReadError(domain.ReadErrorConfigurationInvalid, "sankhya linkage reader is not configured", nil)
	}
	if err := validateSankhyaLinkageConfig(r.config); err != nil {
		return err
	}
	if r.db == nil {
		return domain.NewReadError(domain.ReadErrorSourceUnavailable, "sankhya linkage oracle reader is not configured", nil)
	}
	if err := r.db.PingContext(ctx); err != nil {
		return sankhyaOracleError("ping", err)
	}

	query, args := buildSankhyaMetadataQuery(r.config)
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return sankhyaOracleError("validate header field metadata", err)
	}
	defer rows.Close()

	rowCount := 0
	for rows.Next() {
		rowCount++
		var dataType string
		var charLength sql.NullInt64
		var charUsed sql.NullString
		if err := rows.Scan(&dataType, &charLength, &charUsed); err != nil {
			return sankhyaOracleError("scan header field metadata", err)
		}
		if !compatibleSankhyaHeaderField(dataType, charLength, charUsed) {
			return domain.NewReadError(domain.ReadErrorMetadataMismatch, "sankhya header field metadata is incompatible", nil)
		}
	}
	if err := rows.Err(); err != nil {
		return sankhyaOracleError("iterate header field metadata", err)
	}
	if rowCount != 1 {
		return domain.NewReadError(domain.ReadErrorMetadataMismatch, "sankhya header field metadata did not resolve exactly once", nil)
	}

	duplicateQuery, duplicateArgs, err := buildSankhyaDuplicateQuery(r.config)
	if err != nil {
		return err
	}
	var duplicateGroups int64
	if err := r.db.QueryRowContext(ctx, duplicateQuery, duplicateArgs...).Scan(&duplicateGroups); err != nil {
		return sankhyaOracleError("validate header field uniqueness", err)
	}
	if duplicateGroups != 0 {
		return domain.NewReadError(domain.ReadErrorUniquenessUnproved, "sankhya header field contains duplicate nonblank values", nil)
	}
	return nil
}

func (r *SankhyaLinkageReader) FindCandidates(ctx context.Context, input ports.SankhyaCandidateInput) ([]domain.SankhyaDocumentCandidate, error) {
	if err := r.ValidateConfiguration(ctx); err != nil {
		return nil, err
	}
	if input.ExternalOrderKey == "" || input.ExternalOrderKey != strings.TrimSpace(input.ExternalOrderKey) || !strings.HasPrefix(input.ExternalOrderKey, "ml:v1:") {
		return nil, domain.NewReadError(domain.ReadErrorConfigurationInvalid, "sankhya candidate lookup requires an exact account-scoped external key", nil)
	}

	query, args, err := buildSankhyaCandidateQuery(r.config, input.ExternalOrderKey)
	if err != nil {
		return nil, err
	}
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, sankhyaOracleError("read linkage candidates", err)
	}
	defer rows.Close()

	candidates := make([]domain.SankhyaDocumentCandidate, 0, r.config.CandidateLimit)
	for rows.Next() {
		var documentID int64
		var documentNumber sql.NullInt64
		var operationCode int64
		var negotiatedAt sql.NullTime
		var amount sql.NullFloat64
		if err := rows.Scan(&documentID, &documentNumber, &operationCode, &negotiatedAt, &amount); err != nil {
			return nil, sankhyaOracleError("scan linkage candidate", err)
		}
		if documentID <= 0 || operationCode != r.config.ExpectedOriginTOP {
			return nil, domain.NewReadError(domain.ReadErrorTOPMismatch, "sankhya candidate returned an invalid identity or origin TOP", nil)
		}
		candidates = append(candidates, domain.SankhyaDocumentCandidate{
			DocumentID:     documentID,
			DocumentNumber: sankhyaNullableInt64(documentNumber),
			OperationCode:  operationCode,
			NegotiatedAt:   nullableTime(negotiatedAt),
			Amount:         nullableFloat(amount),
		})
	}
	if err := rows.Err(); err != nil {
		return nil, sankhyaOracleError("iterate linkage candidates", err)
	}
	if len(candidates) > r.config.CandidateLimit {
		return nil, domain.NewReadError(domain.ReadErrorCandidateAmbiguous, "sankhya candidate lookup exceeded its configured bound", nil)
	}

	for index := range candidates {
		lines, err := r.readCandidateLines(ctx, candidates[index].DocumentID)
		if err != nil {
			return nil, err
		}
		candidates[index].Lines = lines
	}
	return candidates, nil
}

func (r *SankhyaLinkageReader) ListDescendants(ctx context.Context, input ports.SankhyaDescendantInput) (domain.SankhyaLineage, error) {
	if err := r.ValidateConfiguration(ctx); err != nil {
		return domain.SankhyaLineage{}, err
	}
	if !input.Origin.Valid() {
		return domain.SankhyaLineage{}, domain.NewReadError(domain.ReadErrorLineageConflict, "sankhya lineage requires an exact origin document line", nil)
	}

	query, args, err := buildSankhyaDescendantQuery(r.config, input.Origin)
	if err != nil {
		return domain.SankhyaLineage{}, err
	}
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return domain.SankhyaLineage{}, sankhyaOracleError("read linkage descendants", err)
	}
	defer rows.Close()

	descendants := make([]domain.SankhyaDescendant, 0)
	for rows.Next() {
		var documentID int64
		var lineNumber int
		var operationCode int64
		var attended sql.NullFloat64
		if err := rows.Scan(&documentID, &lineNumber, &operationCode, &attended); err != nil {
			return domain.SankhyaLineage{}, sankhyaOracleError("scan linkage descendant", err)
		}
		if operationCode != r.config.ExpectedDestinationTOP {
			return domain.SankhyaLineage{}, domain.NewReadError(domain.ReadErrorTOPMismatch, "sankhya descendant returned a non-destination TOP", nil)
		}
		descendants = append(descendants, domain.SankhyaDescendant{
			Identity:         domain.InternalDocumentLine{DocumentID: documentID, LineNumber: lineNumber},
			AttendedQuantity: nullableFloat(attended),
		})
	}
	if err := rows.Err(); err != nil {
		return domain.SankhyaLineage{}, sankhyaOracleError("iterate linkage descendants", err)
	}
	if len(descendants) > r.config.LineageLimit {
		return domain.SankhyaLineage{Origin: input.Origin, State: domain.SankhyaLineageConflict}, domain.NewReadError(domain.ReadErrorLineageConflict, "sankhya lineage exceeded its configured bound", nil)
	}
	return domain.NewSankhyaLineage(input.Origin, input.ExpectedOriginQuantity, descendants)
}

func (r *SankhyaLinkageReader) readCandidateLines(ctx context.Context, documentID int64) ([]domain.SankhyaLineEvidence, error) {
	query, args, err := buildSankhyaCandidateLinesQuery(r.config, documentID)
	if err != nil {
		return nil, err
	}
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, sankhyaOracleError("read candidate lines", err)
	}
	defer rows.Close()

	lines := make([]domain.SankhyaLineEvidence, 0)
	for rows.Next() {
		var returnedDocumentID int64
		var lineNumber int
		var productID sql.NullInt64
		var quantity sql.NullFloat64
		var amount sql.NullFloat64
		if err := rows.Scan(&returnedDocumentID, &lineNumber, &productID, &quantity, &amount); err != nil {
			return nil, sankhyaOracleError("scan candidate line", err)
		}
		identity := domain.InternalDocumentLine{DocumentID: returnedDocumentID, LineNumber: lineNumber}
		if returnedDocumentID != documentID || !identity.Valid() {
			return nil, domain.NewReadError(domain.ReadErrorCandidateAmbiguous, "sankhya candidate line identity is invalid", nil)
		}
		lines = append(lines, domain.SankhyaLineEvidence{
			Identity:  identity,
			ProductID: sankhyaNullableInt64(productID),
			Quantity:  nullableFloat(quantity),
			Amount:    nullableFloat(amount),
		})
	}
	if err := rows.Err(); err != nil {
		return nil, sankhyaOracleError("iterate candidate lines", err)
	}
	if len(lines) > r.config.CandidateLineLimit {
		return nil, domain.NewReadError(domain.ReadErrorCandidateAmbiguous, "sankhya candidate lines exceeded their configured bound", nil)
	}
	return lines, nil
}

func validateSankhyaLinkageConfig(config SankhyaLinkageConfig) error {
	if _, err := quoteSankhyaIdentifier(config.Schema); err != nil {
		return domain.NewReadError(domain.ReadErrorConfigurationInvalid, "sankhya linkage schema is invalid", err)
	}
	if _, err := quoteSankhyaIdentifier(config.HeaderFieldName); err != nil {
		return domain.NewReadError(domain.ReadErrorConfigurationInvalid, "sankhya linkage header field name is invalid", err)
	}
	if strings.TrimSpace(config.ConfigurationRevision) == "" || strings.TrimSpace(config.UniquenessAttestationID) == "" {
		return domain.NewReadError(domain.ReadErrorUniquenessUnproved, "sankhya linkage revision and uniqueness attestation are required", nil)
	}
	if config.ExpectedOriginTOP != sankhyaOriginTOP || config.ExpectedDestinationTOP != sankhyaDestinationTOP {
		return domain.NewReadError(domain.ReadErrorConfigurationInvalid, "sankhya linkage TOP configuration is invalid", nil)
	}
	if config.CandidateLimit <= 0 || config.CandidateLimit > maxCandidateLimit {
		return domain.NewReadError(domain.ReadErrorConfigurationInvalid, "sankhya linkage candidate limit is invalid", nil)
	}
	if config.CandidateLineLimit <= 0 || config.CandidateLineLimit > maxCandidateLineLimit {
		return domain.NewReadError(domain.ReadErrorConfigurationInvalid, "sankhya linkage candidate line limit is invalid", nil)
	}
	if config.LineageLimit <= 0 || config.LineageLimit > maxLineageLimit {
		return domain.NewReadError(domain.ReadErrorConfigurationInvalid, "sankhya linkage lineage limit is invalid", nil)
	}
	return nil
}

func compatibleSankhyaHeaderField(dataType string, charLength sql.NullInt64, charUsed sql.NullString) bool {
	typeName := strings.ToUpper(strings.TrimSpace(dataType))
	return (typeName == "VARCHAR2" || typeName == "NVARCHAR2") && charLength.Valid && charLength.Int64 >= 160 && charUsed.Valid && strings.EqualFold(charUsed.String, "C")
}

func quoteSankhyaIdentifier(identifier string) (string, error) {
	if identifier != strings.ToUpper(identifier) || !oracleLinkageIdentifier.MatchString(identifier) {
		return "", errors.New("identifier is not an uppercase Oracle-safe name")
	}
	return `"` + identifier + `"`, nil
}

func buildSankhyaMetadataQuery(config SankhyaLinkageConfig) (string, []any) {
	return `
SELECT c.DATA_TYPE, c.CHAR_LENGTH, c.CHAR_USED
FROM ALL_TAB_COLUMNS c
WHERE c.OWNER = :1
  AND c.TABLE_NAME = :2
  AND c.COLUMN_NAME = :3
FETCH FIRST 2 ROWS ONLY`, []any{config.Schema, "TGFCAB", config.HeaderFieldName}
}

func buildSankhyaDuplicateQuery(config SankhyaLinkageConfig) (string, []any, error) {
	schema, err := quoteSankhyaIdentifier(config.Schema)
	if err != nil {
		return "", nil, domain.NewReadError(domain.ReadErrorConfigurationInvalid, "sankhya linkage schema is invalid", err)
	}
	field, err := quoteSankhyaIdentifier(config.HeaderFieldName)
	if err != nil {
		return "", nil, domain.NewReadError(domain.ReadErrorConfigurationInvalid, "sankhya linkage header field name is invalid", err)
	}
	query := fmt.Sprintf(`
SELECT COUNT(*)
FROM (
  SELECT c.%s
  FROM %s.TGFCAB c
  WHERE c.%s IS NOT NULL
  GROUP BY c.%s
  HAVING COUNT(*) > 1
  FETCH FIRST 1 ROW ONLY
)`, field, schema, field, field)
	return query, nil, nil
}

func buildSankhyaCandidateQuery(config SankhyaLinkageConfig, externalKey string) (string, []any, error) {
	schema, err := quoteSankhyaIdentifier(config.Schema)
	if err != nil {
		return "", nil, domain.NewReadError(domain.ReadErrorConfigurationInvalid, "sankhya linkage schema is invalid", err)
	}
	field, err := quoteSankhyaIdentifier(config.HeaderFieldName)
	if err != nil {
		return "", nil, domain.NewReadError(domain.ReadErrorConfigurationInvalid, "sankhya linkage header field name is invalid", err)
	}
	query := fmt.Sprintf(`
SELECT c.NUNOTA, c.NUMNOTA, c.CODTIPOPER, c.DTNEG, c.VLRNOTA
FROM %s.TGFCAB c
WHERE c.%s = :1
  AND c.CODTIPOPER = :2
ORDER BY c.NUNOTA
FETCH FIRST :3 ROWS ONLY`, schema, field)
	return query, []any{externalKey, config.ExpectedOriginTOP, config.CandidateLimit + 1}, nil
}

func buildSankhyaCandidateLinesQuery(config SankhyaLinkageConfig, documentID int64) (string, []any, error) {
	schema, err := quoteSankhyaIdentifier(config.Schema)
	if err != nil {
		return "", nil, domain.NewReadError(domain.ReadErrorConfigurationInvalid, "sankhya linkage schema is invalid", err)
	}
	query := fmt.Sprintf(`
SELECT i.NUNOTA, i.SEQUENCIA, i.CODPROD, i.QTDNEG, i.VLRTOT
FROM %s.TGFITE i
WHERE i.NUNOTA = :1
ORDER BY i.SEQUENCIA
FETCH FIRST :2 ROWS ONLY`, schema)
	return query, []any{documentID, config.CandidateLineLimit + 1}, nil
}

func buildSankhyaDescendantQuery(config SankhyaLinkageConfig, origin domain.InternalDocumentLine) (string, []any, error) {
	schema, err := quoteSankhyaIdentifier(config.Schema)
	if err != nil {
		return "", nil, domain.NewReadError(domain.ReadErrorConfigurationInvalid, "sankhya linkage schema is invalid", err)
	}
	query := fmt.Sprintf(`
SELECT i.NUNOTA, i.SEQUENCIA, c.CODTIPOPER, v.QTDATENDIDA
FROM %s.TGFVAR v
JOIN %s.TGFCAB c ON c.NUNOTA = v.NUNOTA
JOIN %s.TGFITE i ON i.NUNOTA = v.NUNOTA AND i.SEQUENCIA = v.SEQUENCIA
WHERE v.NUNOTAORIG = :1
  AND v.SEQUENCIAORIG = :2
  AND c.CODTIPOPER = :3
ORDER BY i.NUNOTA, i.SEQUENCIA
FETCH FIRST :4 ROWS ONLY`, schema, schema, schema)
	return query, []any{origin.DocumentID, origin.LineNumber, config.ExpectedDestinationTOP, config.LineageLimit + 1}, nil
}

func sankhyaOracleError(operation string, err error) error {
	return domain.NewReadError(domain.ReadErrorSourceUnavailable, fmt.Sprintf("sankhya oracle %s failed", operation), err)
}

func sankhyaNullableInt64(value sql.NullInt64) *int64 {
	if !value.Valid {
		return nil
	}
	result := value.Int64
	return &result
}

var _ ports.SankhyaLinkageReader = (*SankhyaLinkageReader)(nil)
