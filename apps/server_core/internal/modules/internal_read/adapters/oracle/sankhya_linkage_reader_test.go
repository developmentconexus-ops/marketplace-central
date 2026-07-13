package oracle

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"strings"
	"sync/atomic"
	"testing"

	"marketplace-central/apps/server_core/internal/modules/internal_read/domain"
	"marketplace-central/apps/server_core/internal/modules/internal_read/ports"
)

var sankhyaDriverSequence uint64

type scriptedSankhyaDB struct {
	queries []string
	respond func(string, []driver.NamedValue) ([]string, [][]driver.Value, error)
}

type scriptedSankhyaDriver struct{ script *scriptedSankhyaDB }
type scriptedSankhyaConn struct{ script *scriptedSankhyaDB }
type scriptedSankhyaRows struct {
	columns []string
	rows    [][]driver.Value
	index   int
}

func (d scriptedSankhyaDriver) Open(string) (driver.Conn, error) {
	return &scriptedSankhyaConn{script: d.script}, nil
}

func (c *scriptedSankhyaConn) Prepare(string) (driver.Stmt, error) {
	return nil, errors.New("prepare is not supported")
}

func (c *scriptedSankhyaConn) Close() error { return nil }
func (c *scriptedSankhyaConn) Begin() (driver.Tx, error) {
	return nil, errors.New("transactions are not supported")
}
func (c *scriptedSankhyaConn) Ping(context.Context) error { return nil }

func (c *scriptedSankhyaConn) QueryContext(_ context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	c.script.queries = append(c.script.queries, query)
	columns, rows, err := c.script.respond(query, args)
	if err != nil {
		return nil, err
	}
	return &scriptedSankhyaRows{columns: columns, rows: rows}, nil
}

func (r *scriptedSankhyaRows) Columns() []string { return r.columns }
func (r *scriptedSankhyaRows) Close() error      { return nil }
func (r *scriptedSankhyaRows) Next(values []driver.Value) error {
	if r.index >= len(r.rows) {
		return io.EOF
	}
	copy(values, r.rows[r.index])
	r.index++
	return nil
}

func openScriptedSankhyaDB(t *testing.T, script *scriptedSankhyaDB) *sql.DB {
	t.Helper()
	name := fmt.Sprintf("sankhya-linkage-%d", atomic.AddUint64(&sankhyaDriverSequence, 1))
	sql.Register(name, scriptedSankhyaDriver{script: script})
	db, err := sql.Open(name, "")
	if err != nil {
		t.Fatalf("sql.Open() error = %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	return db
}

type countingSankhyaQueryer struct {
	calls int
}

func (q *countingSankhyaQueryer) PingContext(context.Context) error {
	q.calls++
	return nil
}

func (q *countingSankhyaQueryer) QueryContext(context.Context, string, ...any) (*sql.Rows, error) {
	q.calls++
	return nil, nil
}

func (q *countingSankhyaQueryer) QueryRowContext(context.Context, string, ...any) *sql.Row {
	q.calls++
	return &sql.Row{}
}

func validSankhyaLinkageConfig() SankhyaLinkageConfig {
	return SankhyaLinkageConfig{
		Schema:                  "METALPRD",
		HeaderFieldName:         "AD_MPC_ORDER_KEY",
		ConfigurationRevision:   "rev-1",
		ExpectedOriginTOP:       313,
		ExpectedDestinationTOP:  306,
		UniquenessAttestationID: "attestation-1",
		CandidateLimit:          20,
		CandidateLineLimit:      100,
		LineageLimit:            100,
	}
}

func TestSankhyaLinkageInvalidIdentifierFailsBeforeDatabaseAccess(t *testing.T) {
	db := &countingSankhyaQueryer{}
	config := validSankhyaLinkageConfig()
	config.HeaderFieldName = `AD_BAD" OR 1=1`
	reader := NewSankhyaLinkageReader(db, config)

	_, err := reader.FindCandidates(context.Background(), ports.SankhyaCandidateInput{ExternalOrderKey: "123456"})
	if !domain.IsReadErrorCode(err, domain.ReadErrorConfigurationInvalid) {
		t.Fatalf("error = %v, want configuration_invalid", err)
	}
	if db.calls != 0 {
		t.Fatalf("database calls = %d, want 0", db.calls)
	}
}

func TestSankhyaLinkageMissingAttestationFailsBeforeDatabaseAccess(t *testing.T) {
	db := &countingSankhyaQueryer{}
	config := validSankhyaLinkageConfig()
	config.UniquenessAttestationID = ""
	err := NewSankhyaLinkageReader(db, config).ValidateConfiguration(context.Background())
	if !domain.IsReadErrorCode(err, domain.ReadErrorUniquenessUnproved) {
		t.Fatalf("error = %v, want uniqueness_unproved", err)
	}
	if db.calls != 0 {
		t.Fatalf("database calls = %d, want 0", db.calls)
	}
}

func TestSankhyaLinkageUnavailableDatabaseIsTyped(t *testing.T) {
	err := NewSankhyaLinkageReader(nil, validSankhyaLinkageConfig()).ValidateConfiguration(context.Background())
	if !domain.IsReadErrorCode(err, domain.ReadErrorSourceUnavailable) {
		t.Fatalf("error = %v, want source_unavailable", err)
	}
}

func TestSankhyaLinkageLineAndLineageLimitsAreExplicitAndBounded(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*SankhyaLinkageConfig)
	}{
		{name: "missing candidate line limit", mutate: func(config *SankhyaLinkageConfig) { config.CandidateLineLimit = 0 }},
		{name: "excess candidate line limit", mutate: func(config *SankhyaLinkageConfig) { config.CandidateLineLimit = maxCandidateLineLimit + 1 }},
		{name: "missing lineage limit", mutate: func(config *SankhyaLinkageConfig) { config.LineageLimit = 0 }},
		{name: "excess lineage limit", mutate: func(config *SankhyaLinkageConfig) { config.LineageLimit = maxLineageLimit + 1 }},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			config := validSankhyaLinkageConfig()
			test.mutate(&config)
			if err := validateSankhyaLinkageConfig(config); !domain.IsReadErrorCode(err, domain.ReadErrorConfigurationInvalid) {
				t.Fatalf("validateSankhyaLinkageConfig() error = %v, want configuration_invalid", err)
			}
		})
	}
}

func TestSankhyaLinkageMissingMetadataStopsBeforeCandidateQuery(t *testing.T) {
	script := &scriptedSankhyaDB{respond: func(query string, _ []driver.NamedValue) ([]string, [][]driver.Value, error) {
		if strings.Contains(query, "ALL_TAB_COLUMNS") {
			return []string{"DATA_TYPE", "NULLABLE"}, nil, nil
		}
		return nil, nil, fmt.Errorf("unexpected query after missing metadata: %s", query)
	}}
	reader := NewSankhyaLinkageReader(openScriptedSankhyaDB(t, script), validSankhyaLinkageConfig())

	_, err := reader.FindCandidates(context.Background(), ports.SankhyaCandidateInput{ExternalOrderKey: "123456"})
	if !domain.IsReadErrorCode(err, domain.ReadErrorMetadataMismatch) {
		t.Fatalf("error = %v, want field_metadata_mismatch", err)
	}
	if len(script.queries) != 1 || !strings.Contains(script.queries[0], "ALL_TAB_COLUMNS") {
		t.Fatalf("queries = %#v, want metadata only", script.queries)
	}
}

func TestSankhyaLinkageDuplicateProbeStopsBeforeCandidateQuery(t *testing.T) {
	script := &scriptedSankhyaDB{respond: func(query string, _ []driver.NamedValue) ([]string, [][]driver.Value, error) {
		switch {
		case strings.Contains(query, "ALL_TAB_COLUMNS"):
			return []string{"DATA_TYPE", "NULLABLE"}, [][]driver.Value{{"CLOB", "Y"}}, nil
		case strings.Contains(query, "other.ROWID <> c.ROWID"):
			return []string{"DUPLICATE_GROUPS"}, [][]driver.Value{{int64(1)}}, nil
		default:
			return nil, nil, fmt.Errorf("unexpected query after duplicate probe: %s", query)
		}
	}}
	reader := NewSankhyaLinkageReader(openScriptedSankhyaDB(t, script), validSankhyaLinkageConfig())

	_, err := reader.FindCandidates(context.Background(), ports.SankhyaCandidateInput{ExternalOrderKey: "123456"})
	if !domain.IsReadErrorCode(err, domain.ReadErrorUniquenessUnproved) {
		t.Fatalf("error = %v, want uniqueness_unproved", err)
	}
	if len(script.queries) != 2 || !strings.Contains(script.queries[1], "DBMS_LOB.COMPARE") {
		t.Fatalf("queries = %#v, want metadata then duplicate probe", script.queries)
	}
}

func TestSankhyaLinkageExternalOrderKeyMustBeDigitsOnly(t *testing.T) {
	for _, key := range []string{"", " 123", "123 ", "ml:v1:123", "123-456"} {
		if sankhyaExternalOrderKey.MatchString(key) {
			t.Fatalf("external key %q was accepted", key)
		}
	}
	if !sankhyaExternalOrderKey.MatchString("123456") {
		t.Fatal("digits-only external key was rejected")
	}
}

func TestSankhyaLinkageReaderPreservesOneToManyNullableDescendants(t *testing.T) {
	script := &scriptedSankhyaDB{respond: func(query string, _ []driver.NamedValue) ([]string, [][]driver.Value, error) {
		switch {
		case strings.Contains(query, "ALL_TAB_COLUMNS"):
			return []string{"DATA_TYPE", "NULLABLE"}, [][]driver.Value{{"CLOB", "Y"}}, nil
		case strings.Contains(query, "other.ROWID <> c.ROWID"):
			return []string{"DUPLICATE_GROUPS"}, [][]driver.Value{{int64(0)}}, nil
		case strings.Contains(query, "TGFVAR v"):
			return []string{"NUNOTA", "SEQUENCIA", "CODTIPOPER", "QTDATENDIDA"}, [][]driver.Value{
				{int64(30601), int64(1), int64(306), float64(1)},
				{int64(30602), int64(3), int64(306), nil},
			}, nil
		default:
			return nil, nil, fmt.Errorf("unexpected query: %s", query)
		}
	}}
	reader := NewSankhyaLinkageReader(openScriptedSankhyaDB(t, script), validSankhyaLinkageConfig())

	got, err := reader.ListDescendants(context.Background(), ports.SankhyaDescendantInput{
		Origin: domain.InternalDocumentLine{DocumentID: 31301, LineNumber: 7},
	})
	if err != nil {
		t.Fatalf("ListDescendants() error = %v", err)
	}
	if got.State != domain.SankhyaLineagePartial || len(got.Descendants) != 2 || got.Descendants[1].AttendedQuantity != nil {
		t.Fatalf("ListDescendants() = %+v, want two partial descendants with nullable quantity", got)
	}
}

func TestSankhyaLinkageCandidateLineOverflowFailsAmbiguous(t *testing.T) {
	script := &scriptedSankhyaDB{respond: func(query string, _ []driver.NamedValue) ([]string, [][]driver.Value, error) {
		switch {
		case strings.Contains(query, "ALL_TAB_COLUMNS"):
			return []string{"DATA_TYPE", "NULLABLE"}, [][]driver.Value{{"CLOB", "Y"}}, nil
		case strings.Contains(query, "other.ROWID <> c.ROWID"):
			return []string{"DUPLICATE_GROUPS"}, [][]driver.Value{{int64(0)}}, nil
		case strings.Contains(query, "SELECT c.NUNOTA"):
			return []string{"NUNOTA", "NUMNOTA", "CODTIPOPER", "DTNEG", "VLRNOTA"}, [][]driver.Value{{int64(31301), nil, int64(313), nil, nil}}, nil
		case strings.Contains(query, "TGFITE i"):
			return []string{"NUNOTA", "SEQUENCIA", "CODPROD", "QTDNEG", "VLRTOT"}, [][]driver.Value{
				{int64(31301), int64(1), nil, nil, nil},
				{int64(31301), int64(2), nil, nil, nil},
			}, nil
		default:
			return nil, nil, fmt.Errorf("unexpected query: %s", query)
		}
	}}
	config := validSankhyaLinkageConfig()
	config.CandidateLineLimit = 1
	reader := NewSankhyaLinkageReader(openScriptedSankhyaDB(t, script), config)

	_, err := reader.FindCandidates(context.Background(), ports.SankhyaCandidateInput{ExternalOrderKey: "123456"})
	if !domain.IsReadErrorCode(err, domain.ReadErrorCandidateAmbiguous) {
		t.Fatalf("error = %v, want candidate_ambiguous", err)
	}
}

func TestSankhyaLinkageDescendantOverflowFailsConflict(t *testing.T) {
	script := &scriptedSankhyaDB{respond: func(query string, _ []driver.NamedValue) ([]string, [][]driver.Value, error) {
		switch {
		case strings.Contains(query, "ALL_TAB_COLUMNS"):
			return []string{"DATA_TYPE", "NULLABLE"}, [][]driver.Value{{"CLOB", "Y"}}, nil
		case strings.Contains(query, "other.ROWID <> c.ROWID"):
			return []string{"DUPLICATE_GROUPS"}, [][]driver.Value{{int64(0)}}, nil
		case strings.Contains(query, "TGFVAR v"):
			return []string{"NUNOTA", "SEQUENCIA", "CODTIPOPER", "QTDATENDIDA"}, [][]driver.Value{
				{int64(30601), int64(1), int64(306), float64(1)},
				{int64(30602), int64(1), int64(306), float64(1)},
			}, nil
		default:
			return nil, nil, fmt.Errorf("unexpected query: %s", query)
		}
	}}
	config := validSankhyaLinkageConfig()
	config.LineageLimit = 1
	reader := NewSankhyaLinkageReader(openScriptedSankhyaDB(t, script), config)

	got, err := reader.ListDescendants(context.Background(), ports.SankhyaDescendantInput{
		Origin: domain.InternalDocumentLine{DocumentID: 31301, LineNumber: 7},
	})
	if !domain.IsReadErrorCode(err, domain.ReadErrorLineageConflict) || got.State != domain.SankhyaLineageConflict {
		t.Fatalf("ListDescendants() = %+v, %v; want lineage conflict", got, err)
	}
}

func TestSankhyaLinkageMetadataQueryUsesOnlyBinds(t *testing.T) {
	config := validSankhyaLinkageConfig()
	query, args := buildSankhyaMetadataQuery(config)
	if strings.Contains(query, config.Schema) || strings.Contains(query, config.HeaderFieldName) {
		t.Fatalf("metadata values were interpolated: %s", query)
	}
	for _, predicate := range []string{"c.OWNER = :1", "c.TABLE_NAME = :2", "c.COLUMN_NAME = :3", "FETCH FIRST 2 ROWS ONLY"} {
		if !strings.Contains(query, predicate) {
			t.Fatalf("metadata query missing %q: %s", predicate, query)
		}
	}
	if len(args) != 3 || args[0] != config.Schema || args[1] != "TGFCAB" || args[2] != config.HeaderFieldName {
		t.Fatalf("metadata args = %#v", args)
	}
}

func TestSankhyaLinkageCandidateQueryQuotesIdentifierAndBindsValues(t *testing.T) {
	config := validSankhyaLinkageConfig()
	externalKey := "123456"
	query, args, err := buildSankhyaCandidateQuery(config, externalKey)
	if err != nil {
		t.Fatalf("buildSankhyaCandidateQuery() error = %v", err)
	}

	for _, fragment := range []string{`FROM "METALPRD".TGFCAB`, `DBMS_LOB.COMPARE(c."AD_MPC_ORDER_KEY", TO_CLOB(:1)) = 0`, "c.CODTIPOPER = :2", "FETCH FIRST :3 ROWS ONLY"} {
		if !strings.Contains(query, fragment) {
			t.Fatalf("candidate query missing %q: %s", fragment, query)
		}
	}
	if strings.Contains(query, externalKey) || len(args) != 3 || args[0] != externalKey || args[1] != int64(313) || args[2] != 21 {
		t.Fatalf("candidate query/args = %s / %#v", query, args)
	}
}

func TestSankhyaLinkageDuplicateProbeUsesValidatedQuotedIdentifier(t *testing.T) {
	config := validSankhyaLinkageConfig()
	query, args, err := buildSankhyaDuplicateQuery(config)
	if err != nil {
		t.Fatalf("buildSankhyaDuplicateQuery() error = %v", err)
	}
	if len(args) != 0 || !strings.Contains(query, `c."AD_MPC_ORDER_KEY" IS NOT NULL`) || !strings.Contains(query, `other.ROWID <> c.ROWID`) || !strings.Contains(query, `DBMS_LOB.COMPARE(other."AD_MPC_ORDER_KEY", c."AD_MPC_ORDER_KEY") = 0`) {
		t.Fatalf("duplicate query/args = %s / %#v", query, args)
	}
}

func TestSankhyaLinkageLineAndDescendantQueriesUseExactIdentity(t *testing.T) {
	config := validSankhyaLinkageConfig()
	lineQuery, lineArgs, err := buildSankhyaCandidateLinesQuery(config, 31301)
	if err != nil {
		t.Fatalf("buildSankhyaCandidateLinesQuery() error = %v", err)
	}
	if !strings.Contains(lineQuery, "i.NUNOTA = :1") || !strings.Contains(lineQuery, "FETCH FIRST :2 ROWS ONLY") || len(lineArgs) != 2 || lineArgs[0] != int64(31301) || lineArgs[1] != 101 {
		t.Fatalf("line query/args = %s / %#v", lineQuery, lineArgs)
	}

	origin := domain.InternalDocumentLine{DocumentID: 31301, LineNumber: 7}
	descendantQuery, descendantArgs, err := buildSankhyaDescendantQuery(config, origin)
	if err != nil {
		t.Fatalf("buildSankhyaDescendantQuery() error = %v", err)
	}
	for _, fragment := range []string{"TGFVAR v", "v.NUNOTAORIG = :1", "v.SEQUENCIAORIG = :2", "c.CODTIPOPER = :3", "i.NUNOTA = v.NUNOTA", "i.SEQUENCIA = v.SEQUENCIA", "FETCH FIRST :4 ROWS ONLY"} {
		if !strings.Contains(descendantQuery, fragment) {
			t.Fatalf("descendant query missing %q: %s", fragment, descendantQuery)
		}
	}
	if len(descendantArgs) != 4 || descendantArgs[0] != int64(31301) || descendantArgs[1] != 7 || descendantArgs[2] != int64(306) || descendantArgs[3] != 101 {
		t.Fatalf("descendant args = %#v", descendantArgs)
	}
}

func TestSankhyaLinkageBuildersReturnTypedErrorsWithoutPanic(t *testing.T) {
	config := validSankhyaLinkageConfig()
	config.Schema = `BAD"SCHEMA`
	builders := []func() error{
		func() error { _, _, err := buildSankhyaDuplicateQuery(config); return err },
		func() error { _, _, err := buildSankhyaCandidateQuery(config, "123456"); return err },
		func() error { _, _, err := buildSankhyaCandidateLinesQuery(config, 1); return err },
		func() error {
			_, _, err := buildSankhyaDescendantQuery(config, domain.InternalDocumentLine{DocumentID: 1, LineNumber: 1})
			return err
		},
	}
	for index, build := range builders {
		if err := build(); !domain.IsReadErrorCode(err, domain.ReadErrorConfigurationInvalid) {
			t.Fatalf("builder %d error = %v, want configuration_invalid", index, err)
		}
	}
}

func TestSankhyaLinkageChangedProductionFilesContainNoPanic(t *testing.T) {
	files := []string{
		"sankhya_linkage_reader.go",
		"../../domain/sankhya_linkage.go",
		"../../ports/sankhya_linkage.go",
		"../../application/sankhya_linkage_service.go",
	}
	for _, path := range files {
		fileSet := token.NewFileSet()
		parsed, err := parser.ParseFile(fileSet, path, nil, 0)
		if err != nil {
			t.Fatalf("parse %s: %v", path, err)
		}
		ast.Inspect(parsed, func(node ast.Node) bool {
			call, ok := node.(*ast.CallExpr)
			if !ok {
				return true
			}
			identifier, ok := call.Fun.(*ast.Ident)
			if ok && identifier.Name == "panic" {
				t.Errorf("production panic found in %s at %s", path, fileSet.Position(call.Pos()))
			}
			return true
		})
	}
}

func TestSankhyaLinkageIdentifierAndMetadataCompatibilityAreStrict(t *testing.T) {
	if _, err := quoteSankhyaIdentifier("AD_VALID_FIELD"); err != nil {
		t.Fatalf("valid identifier rejected: %v", err)
	}
	for _, invalid := range []string{"ad_lower", " AD_SPACE", "AD.DOTTED", strings.Repeat("A", 31)} {
		if _, err := quoteSankhyaIdentifier(invalid); err == nil {
			t.Fatalf("invalid identifier %q accepted", invalid)
		}
	}
	if !compatibleSankhyaHeaderField("CLOB", "Y") {
		t.Fatal("nullable CLOB should be compatible")
	}
	if compatibleSankhyaHeaderField("CLOB", "N") || compatibleSankhyaHeaderField("VARCHAR2", "Y") {
		t.Fatal("incompatible metadata accepted")
	}
}
