package application

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/modules/erp_import/domain"
	"marketplace-central/apps/server_core/internal/modules/erp_import/ports"
)

type fakeParser struct {
	rows       []domain.NormalizedRow
	fileIssues []domain.Issue
	err        error
	input      []byte
}

func (p *fakeParser) Parse(_ context.Context, source io.Reader) ([]domain.NormalizedRow, []domain.Issue, error) {
	p.input, _ = io.ReadAll(source)
	return p.rows, p.fileIssues, p.err
}

func (p *fakeParser) ParseLenient(_ context.Context, source io.Reader) ([]domain.NormalizedRow, []domain.Issue, error) {
	p.input, _ = io.ReadAll(source)
	return p.rows, p.fileIssues, p.err
}

type fakeImportRepository struct {
	list              []domain.ImportReport
	getReport         domain.ImportReport
	persistErr        error
	getErr            error
	persistCalls      int
	persistTenant     string
	listTenant        string
	getTenant         string
	getID             domain.ImportID
	persistedSnapshot domain.ImportSnapshot
}

func (r *fakeImportRepository) PersistSnapshotAtomically(_ context.Context, tenantID string, snapshot domain.ImportSnapshot) error {
	r.persistCalls++
	r.persistTenant = tenantID
	r.persistedSnapshot = snapshot
	return r.persistErr
}

func (r *fakeImportRepository) FindByFileSHA256(context.Context, string, domain.FileSHA256) (*domain.ImportReport, error) {
	return nil, nil
}

func (r *fakeImportRepository) ListImports(_ context.Context, tenantID string) ([]domain.ImportReport, error) {
	r.listTenant = tenantID
	return r.list, nil
}

func (r *fakeImportRepository) GetImport(_ context.Context, tenantID string, id domain.ImportID) (domain.ImportReport, error) {
	r.getTenant = tenantID
	r.getID = id
	return r.getReport, r.getErr
}

func (r *fakeImportRepository) LatestCompletedSnapshot(context.Context, string, domain.ImportSource) (domain.ImportSnapshot, error) {
	return domain.ImportSnapshot{}, nil
}

func TestImportServiceHappyPath(t *testing.T) {
	clock := time.Date(2026, 7, 17, 13, 30, 0, 123, time.FixedZone("BRT", -3*60*60))
	id := domain.ImportID("123e4567-e89b-42d3-a456-426614174000")
	row := validImportRow("P-1")
	repo := &fakeImportRepository{getReport: domain.ImportReport{ID: id, AcceptedCount: 1}}
	parser := &fakeParser{rows: []domain.NormalizedRow{row}}
	service := NewImportService(parser, repo, "tenant-a", WithClock(func() time.Time { return clock }), WithIDGenerator(func() (domain.ImportID, error) { return id, nil }))

	payload := "xlsx bytes"
	report, err := service.RunImport(context.Background(), strings.NewReader(payload))
	if err != nil {
		t.Fatalf("RunImport() error = %v", err)
	}
	if !reflect.DeepEqual(report, repo.getReport) {
		t.Fatalf("report = %#v, want authoritative repo report %#v", report, repo.getReport)
	}
	digest := sha256.Sum256([]byte(payload))
	wantHash := domain.FileSHA256(hex.EncodeToString(digest[:]))
	snapshot := repo.persistedSnapshot
	if snapshot.ID != id || snapshot.Protocol != "#001-E" || snapshot.FileSHA256 != wantHash {
		t.Fatalf("snapshot identity = %#v", snapshot)
	}
	if snapshot.Status != domain.ImportStatusCompleted || snapshot.Source != domain.SourceXLSX {
		t.Fatalf("snapshot status/source = %q/%q", snapshot.Status, snapshot.Source)
	}
	if !snapshot.ImportedAt.Equal(clock.UTC()) || snapshot.ImportedAt.Location() != time.UTC {
		t.Fatalf("ImportedAt = %v, want %v in UTC", snapshot.ImportedAt, clock.UTC())
	}
	if repo.persistTenant != "tenant-a" || repo.listTenant != "tenant-a" || repo.getTenant != "tenant-a" {
		t.Fatalf("tenant forwarding = persist %q list %q get %q", repo.persistTenant, repo.listTenant, repo.getTenant)
	}
	if repo.getID != id || string(parser.input) != "xlsx bytes" {
		t.Fatalf("GetImport id/input = %q/%q", repo.getID, parser.input)
	}
}

func TestImportServiceValidationOutcomes(t *testing.T) {
	invalidEAN := "123"
	tests := []struct {
		name       string
		rows       []domain.NormalizedRow
		wantStatus domain.ImportStatus
		wantRows   int
		wantIssues int
	}{
		{name: "mixed rejections and warnings", rows: []domain.NormalizedRow{validImportRow("P-1"), {Codprod: "", Descrprod: "bad", Custo: "10", StockPhysical: "1"}, func() domain.NormalizedRow { row := validImportRow("P-2"); row.EAN = &invalidEAN; return row }()}, wantStatus: domain.ImportStatusCompleted, wantRows: 2, wantIssues: 2},
		{name: "all rejected", rows: []domain.NormalizedRow{{Codprod: "", Descrprod: "bad", Custo: "0", StockPhysical: "-1"}}, wantStatus: domain.ImportStatusRejected, wantRows: 0, wantIssues: 3},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := &fakeImportRepository{}
			service := deterministicImportService(&fakeParser{rows: test.rows}, repo)
			if _, err := service.RunImport(context.Background(), strings.NewReader("data")); err != nil {
				t.Fatalf("RunImport() error = %v", err)
			}
			if repo.persistCalls != 1 || repo.persistedSnapshot.Status != test.wantStatus {
				t.Fatalf("persist calls/status = %d/%q", repo.persistCalls, repo.persistedSnapshot.Status)
			}
			if len(repo.persistedSnapshot.AcceptedRows) != test.wantRows || len(repo.persistedSnapshot.Issues) != test.wantIssues {
				t.Fatalf("accepted/issues = %d/%d", len(repo.persistedSnapshot.AcceptedRows), len(repo.persistedSnapshot.Issues))
			}
			if test.wantRows > 0 && repo.persistedSnapshot.AcceptedRows[0].Codprod != "P-1" {
				t.Fatalf("accepted rows not forwarded: %#v", repo.persistedSnapshot.AcceptedRows)
			}
		})
	}
}

func TestImportServiceStructuralParserErrorsDoNotPersist(t *testing.T) {
	for _, code := range []domain.IssueCode{domain.CodeInvalidFile, domain.CodeMissingRequiredColumn} {
		t.Run(string(code), func(t *testing.T) {
			parserErr := &ports.FileError{Code: code, Detail: "structural failure"}
			repo := &fakeImportRepository{}
			_, err := deterministicImportService(&fakeParser{err: parserErr}, repo).RunImport(context.Background(), strings.NewReader("data"))
			var got *ports.FileError
			if !errors.As(err, &got) || got != parserErr {
				t.Fatalf("error = %v, want original FileError", err)
			}
			if repo.persistCalls != 0 {
				t.Fatalf("persist calls = %d, want 0", repo.persistCalls)
			}
		})
	}
}

func TestImportServiceConflictErrorsRemainTyped(t *testing.T) {
	t.Run("duplicate", func(t *testing.T) {
		want := &ports.DuplicateFileError{ExistingID: "existing-id", ExistingProtocol: "#007-E"}
		repo := &fakeImportRepository{persistErr: want}
		_, err := deterministicImportService(&fakeParser{rows: []domain.NormalizedRow{validImportRow("P-1")}}, repo).RunImport(context.Background(), strings.NewReader("data"))
		var got *ports.DuplicateFileError
		if !errors.As(err, &got) || got.ExistingID != want.ExistingID || got.ExistingProtocol != want.ExistingProtocol {
			t.Fatalf("error = %#v, want duplicate fields %#v", got, want)
		}
	})
	t.Run("in progress", func(t *testing.T) {
		repo := &fakeImportRepository{persistErr: ports.ErrImportInProgress}
		_, err := deterministicImportService(&fakeParser{rows: []domain.NormalizedRow{validImportRow("P-1")}}, repo).RunImport(context.Background(), strings.NewReader("data"))
		if !errors.Is(err, ports.ErrImportInProgress) {
			t.Fatalf("error = %v, want ErrImportInProgress", err)
		}
	})
}

func TestImportServiceProtocolGeneration(t *testing.T) {
	pattern := regexp.MustCompile(`^#[0-9]{3,}-E$`)
	for _, test := range []struct {
		name string
		list []domain.ImportReport
		want domain.Protocol
	}{
		{name: "first", want: "#001-E"},
		{name: "monotonic", list: []domain.ImportReport{{Protocol: "#001-E"}, {Protocol: "ignored"}}, want: "#002-E"},
		{name: "past padding", list: []domain.ImportReport{{Protocol: "#999-E"}}, want: "#1000-E"},
	} {
		t.Run(test.name, func(t *testing.T) {
			repo := &fakeImportRepository{list: test.list}
			_, err := deterministicImportService(&fakeParser{rows: []domain.NormalizedRow{validImportRow("P-1")}}, repo).RunImport(context.Background(), strings.NewReader("data"))
			if err != nil {
				t.Fatalf("RunImport() error = %v", err)
			}
			got := repo.persistedSnapshot.Protocol
			if got != test.want || !pattern.MatchString(string(got)) {
				t.Fatalf("protocol = %q, want %q and regex-valid", got, test.want)
			}
		})
	}
}

func TestImportServiceSHA256Determinism(t *testing.T) {
	hashes := make([]domain.FileSHA256, 0, 3)
	for _, input := range []string{"same", "same", "different"} {
		repo := &fakeImportRepository{}
		_, err := deterministicImportService(&fakeParser{rows: []domain.NormalizedRow{validImportRow("P-1")}}, repo).RunImport(context.Background(), strings.NewReader(input))
		if err != nil {
			t.Fatalf("RunImport(%q) error = %v", input, err)
		}
		hashes = append(hashes, repo.persistedSnapshot.FileSHA256)
	}
	if hashes[0] != hashes[1] || hashes[0] == hashes[2] {
		t.Fatalf("hashes = %v, want identical inputs equal and different input unequal", hashes)
	}
}

func TestNewUUIDv4IsCanonicalV4(t *testing.T) {
	v4 := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)
	seen := make(map[domain.ImportID]struct{}, 32)
	for i := 0; i < 32; i++ {
		id, err := newUUIDv4()
		if err != nil {
			t.Fatalf("newUUIDv4() error = %v", err)
		}
		if !v4.MatchString(string(id)) {
			t.Fatalf("newUUIDv4() = %q, not a canonical RFC-4122 v4 UUID", id)
		}
		if _, dup := seen[id]; dup {
			t.Fatalf("newUUIDv4() produced duplicate %q", id)
		}
		seen[id] = struct{}{}
	}
}

func deterministicImportService(parser ports.Parser, repo ports.ImportRepository) *ImportService {
	return NewImportService(parser, repo, "tenant-test", WithClock(func() time.Time { return time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC) }), WithIDGenerator(func() (domain.ImportID, error) { return "123e4567-e89b-42d3-a456-426614174000", nil }))
}

func validImportRow(codprod string) domain.NormalizedRow {
	return domain.NormalizedRow{Codprod: codprod, Descrprod: "Product", Custo: "10.00", StockPhysical: "2"}
}
