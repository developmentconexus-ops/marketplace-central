package internalread

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/modules/erp_import/adapters/xlsx"
	"marketplace-central/apps/server_core/internal/modules/erp_import/application"
	erpdomain "marketplace-central/apps/server_core/internal/modules/erp_import/domain"
	erpports "marketplace-central/apps/server_core/internal/modules/erp_import/ports"
	readdomain "marketplace-central/apps/server_core/internal/modules/internal_read/domain"
	readports "marketplace-central/apps/server_core/internal/modules/internal_read/ports"
)

const (
	contractTenant = "tenant-contract"
	contractImport = "2026-07-17T12:00:00Z"
	contractFetch  = "2026-07-17T13:00:00Z"
)

var contractClock = func() time.Time {
	t, _ := time.Parse(time.RFC3339, contractImport)
	return t
}()

var contractFetchedAt = func() time.Time {
	t, _ := time.Parse(time.RFC3339, contractFetch)
	return t
}()

func TestXLSXReaderSubstitutabilityContract(t *testing.T) {
	fixtures := contractFixturesDir(t)
	writeContractFixtures(t, fixtures)

	exampleData, err := os.ReadFile(filepath.Join(fixtures, "example-erp.xlsx"))
	if err != nil {
		t.Fatal(err)
	}
	exampleRepo, exampleReport := importContractFixture(t, exampleData, "123e4567-e89b-42d3-a456-426614174001")
	if exampleReport.Status != erpdomain.ImportStatusCompleted || exampleReport.AcceptedCount < 50 || exampleReport.WarningCount != 0 {
		t.Fatalf("example report = %+v, want COMPLETED with at least 50 accepted rows", exampleReport)
	}
	if !regexp.MustCompile(`^#[0-9]{3,}-E$`).MatchString(string(exampleReport.Protocol)) {
		t.Fatalf("example protocol = %q, want #NNN-E shape", exampleReport.Protocol)
	}

	xlsxReader := NewReader(exampleRepo, contractTenant, WithClock(func() time.Time { return contractFetchedAt }))
	oracleReader := newOracleShapedReader()
	ctx := context.Background()
	stockPolicy := readdomain.DefaultSellableStockPolicy()
	costPolicy := readdomain.CostAsOfPolicy{CompanyID: 1, Basis: readdomain.CostBasisCUSSEMICM, EffectiveAt: contractClock}
	taxPolicy := readdomain.TaxPolicy{EffectiveAt: contractClock}

	xlsxCandidate, err := xlsxReader.FindProductsForLinking(ctx, readports.FindProductsInput{ProductID: intPointer(1001)})
	if err != nil {
		t.Fatal(err)
	}
	oracleCandidate, err := oracleReader.FindProductsForLinking(ctx, readports.FindProductsInput{ProductID: intPointer(1001)})
	if err != nil {
		t.Fatal(err)
	}
	assertCandidateShapeParity(t, xlsxCandidate, oracleCandidate)
	collisionCandidate, err := xlsxReader.FindProductsForLinking(ctx, readports.FindProductsInput{ProductID: intPointer(1002)})
	if err != nil || len(collisionCandidate) != 1 || collisionCandidate[0].EAN == nil || *collisionCandidate[0].EAN != *xlsxCandidate[0].EAN || !readdomain.HasQualityFlag(collisionCandidate[0].QualityFlags, readdomain.QualityEANCollision) || collisionCandidate[0].ReferenceCode == nil || *collisionCandidate[0].ReferenceCode != "REF-1002" {
		t.Fatalf("EAN collision pair projection = %+v, err=%v", collisionCandidate, err)
	}

	xlsxStock, err := xlsxReader.GetSellableStock(ctx, readports.SellableStockInput{ProductID: 1001, Policy: stockPolicy})
	if err != nil {
		t.Fatal(err)
	}
	oracleStock, err := oracleReader.GetSellableStock(ctx, readports.SellableStockInput{ProductID: 1001, Policy: stockPolicy})
	if err != nil {
		t.Fatal(err)
	}
	assertStockShapeParity(t, xlsxStock, oracleStock)
	if xlsxStock.Quantity == nil || *xlsxStock.Quantity != 8 {
		t.Fatalf("reserved-present stock = %#v, want 8", xlsxStock.Quantity)
	}

	reservedAbsent, err := xlsxReader.GetSellableStock(ctx, readports.SellableStockInput{ProductID: 1002, Policy: stockPolicy})
	if err != nil {
		t.Fatal(err)
	}
	if reservedAbsent.Quantity != nil || !readdomain.HasQualityFlag(reservedAbsent.QualityFlags, readdomain.QualityMissingStock) {
		t.Fatalf("reserved-absent stock = %+v, want unknown with missing_stock", reservedAbsent)
	}
	if row := acceptedRowByID(exampleRepo.snapshot.AcceptedRows, "1002"); row == nil || row.StockPhysical != "11" || row.StockReserved != nil {
		t.Fatalf("reserved-absent persisted row = %+v, want physical 11 and nil reserved", row)
	}
	oracleReservedAbsent, err := oracleReader.GetSellableStock(ctx, readports.SellableStockInput{ProductID: 1002, Policy: stockPolicy})
	if err != nil {
		t.Fatal(err)
	}
	assertStockShapeParity(t, reservedAbsent, oracleReservedAbsent)

	xlsxCost, err := xlsxReader.GetCostAsOf(ctx, readports.CostAsOfInput{ProductID: 1001, Policy: costPolicy})
	if err != nil {
		t.Fatal(err)
	}
	oracleCost, err := oracleReader.GetCostAsOf(ctx, readports.CostAsOfInput{ProductID: 1001, Policy: costPolicy})
	if err != nil {
		t.Fatal(err)
	}
	assertCostShapeParity(t, xlsxCost, oracleCost)

	xlsxTax, err := xlsxReader.GetTaxInputs(ctx, readports.TaxInput{ProductID: 1001, Policy: taxPolicy})
	if err != nil {
		t.Fatal(err)
	}
	oracleTax, err := oracleReader.GetTaxInputs(ctx, readports.TaxInput{ProductID: 1001, Policy: taxPolicy})
	if err != nil {
		t.Fatal(err)
	}
	assertTaxShapeParity(t, xlsxTax, oracleTax)

	// Current price + sales history: the xlsx snapshot is a strict SUBSET of the Oracle source. It
	// carries neither, so it returns a TYPED unsupported outcome (ADR-17: never a zero price / empty
	// history passed off as real data). The Oracle source DOES serve both (oracle/reader.go queries
	// VW_PRECO_TABELA and the sales views by default). The substitutability contract pinned here is
	// exactly that divergence — xlsx signals typed-unsupported where Oracle returns data — NOT identical
	// outcomes. Asserting "both unsupported" would be theater (it would only test a fabricated fake).
	t.Run("current price divergence", func(t *testing.T) {
		if _, err := xlsxReader.GetCurrentPrice(ctx, readports.CurrentPriceInput{ProductID: 1001}); err == nil || err.Error() != "xlsx snapshot does not contain current price" || !readdomain.IsReadErrorCode(err, readdomain.ReadErrorUnsupportedQuery) {
			t.Fatalf("xlsx current price = %v, want typed unsupported", err)
		}
		if price, err := oracleReader.GetCurrentPrice(ctx, readports.CurrentPriceInput{ProductID: 1001}); err != nil || price.Amount == nil {
			t.Fatalf("oracle current price = %+v err=%v, want served data (source supports it)", price, err)
		}
	})
	t.Run("sales history divergence", func(t *testing.T) {
		if _, err := xlsxReader.GetSalesHistory(ctx, readports.SalesHistoryInput{ProductID: intPointer(1001)}); err == nil || err.Error() != "xlsx snapshot does not contain sales history" || !readdomain.IsReadErrorCode(err, readdomain.ReadErrorUnsupportedQuery) {
			t.Fatalf("xlsx sales history = %v, want typed unsupported", err)
		}
		if history, err := oracleReader.GetSalesHistory(ctx, readports.SalesHistoryInput{ProductID: intPointer(1001)}); err != nil || len(history.Entries) == 0 {
			t.Fatalf("oracle sales history = %+v err=%v, want served data (source supports it)", history, err)
		}
	})

	xlsxPage, err := xlsxReader.ListCatalogProductFacts(ctx, readports.Cursor{InternalProductID: 1000}, 1)
	if err != nil {
		t.Fatal(err)
	}
	oraclePage, err := oracleReader.ListCatalogProductFacts(ctx, readports.Cursor{InternalProductID: 1000}, 1)
	if err != nil {
		t.Fatal(err)
	}
	assertCatalogPageShapeParity(t, xlsxPage, oraclePage)

	xlsxSearch, err := xlsxReader.SearchCatalogProductFacts(ctx, "Example Product 1001", 1)
	if err != nil {
		t.Fatal(err)
	}
	oracleSearch, err := oracleReader.SearchCatalogProductFacts(ctx, "Example Product 1001", 1)
	if err != nil {
		t.Fatal(err)
	}
	assertCatalogPageShapeParity(t, xlsxSearch, oracleSearch)

	identityData, err := os.ReadFile(filepath.Join(fixtures, "identity-rejections.xlsx"))
	if err != nil {
		t.Fatal(err)
	}
	identityRepo, identityReport := importContractFixture(t, identityData, "123e4567-e89b-42d3-a456-426614174002")
	if identityReport.Status != erpdomain.ImportStatusCompleted || identityReport.AcceptedCount != 3 || identityReport.RejectedCount != 1 || identityReport.WarningCount != 2 {
		t.Fatalf("identity report = %+v, want 3 accepted, 1 rejected, 2 warnings", identityReport)
	}
	assertIssueKind(t, identityReport.Issues, erpdomain.CodeEmptyDescrprod, erpdomain.Rejection)
	assertIssueKind(t, identityReport.Issues, erpdomain.CodeInvalidEAN, erpdomain.Warning)
	assertIssueKind(t, identityReport.Issues, erpdomain.CodeInvalidNCM, erpdomain.Warning)

	identityReader := NewReader(identityRepo, contractTenant, WithClock(func() time.Time { return contractFetchedAt }))
	invalidEAN, err := identityReader.FindProductsForLinking(ctx, readports.FindProductsInput{ProductID: intPointer(2002)})
	if err != nil || len(invalidEAN) != 1 || invalidEAN[0].EAN != nil || invalidEAN[0].ReferenceCode == nil || *invalidEAN[0].ReferenceCode != "REF-RETAIN-2002" || !readdomain.HasQualityFlag(invalidEAN[0].QualityFlags, readdomain.QualityInvalidEAN) {
		t.Fatalf("invalid EAN projection = %+v, err=%v", invalidEAN, err)
	}
	malformedNCM, err := identityReader.FindProductsForLinking(ctx, readports.FindProductsInput{ProductID: intPointer(2003)})
	if err != nil || len(malformedNCM) != 1 || malformedNCM[0].NCM != nil {
		t.Fatalf("malformed NCM projection = %+v, err=%v", malformedNCM, err)
	}

	allRejected := newContractWorkbook([][]string{
		{"3001", "", "12,00", "1", "", validGTIN(301), "REF-3001", "Synthetic", "12345678"},
	})
	_, allRejectedReport := importContractFixture(t, allRejected, "123e4567-e89b-42d3-a456-426614174003")
	if allRejectedReport.Status != erpdomain.ImportStatusRejected || allRejectedReport.AcceptedCount != 0 || allRejectedReport.RejectedCount != 1 {
		t.Fatalf("all-rejected report = %+v", allRejectedReport)
	}

	noSnapshot := NewReader(&contractRepo{}, contractTenant)
	_, err = noSnapshot.GetSellableStock(ctx, readports.SellableStockInput{ProductID: 1001})
	if !errors.Is(err, ErrNoErpSnapshot) || !readdomain.IsReadErrorCode(err, readdomain.ReadErrorSourceUnavailable) {
		t.Fatalf("unavailable snapshot error = %v", err)
	}
}

func importContractFixture(t *testing.T, data []byte, id erpdomain.ImportID) (*contractRepo, erpdomain.ImportReport) {
	t.Helper()
	repo := &contractRepo{nextIDs: []erpdomain.ImportID{id}}
	service := application.NewImportService(xlsx.NewParser(), repo, contractTenant, application.WithClock(func() time.Time { return contractClock }), application.WithIDGenerator(seqIDs(repo.nextIDs)))
	report, err := service.RunImport(context.Background(), bytes.NewReader(data))
	if err != nil {
		t.Fatalf("RunImport() error = %v", err)
	}
	if !repo.persistCalled {
		t.Fatal("RunImport did not persist a snapshot")
	}
	if len(repo.snapshot.AcceptedRows) != report.AcceptedCount {
		t.Fatalf("persisted accepted rows = %d, report.AcceptedCount = %d (persistence out of sync with returned summary)", len(repo.snapshot.AcceptedRows), report.AcceptedCount)
	}
	return repo, report
}

func seqIDs(ids []erpdomain.ImportID) func() (erpdomain.ImportID, error) {
	index := 0
	return func() (erpdomain.ImportID, error) {
		if index >= len(ids) {
			return "", errors.New("deterministic ID sequence exhausted")
		}
		id := ids[index]
		index++
		return id, nil
	}
}

type contractRepo struct {
	snapshot      erpdomain.ImportSnapshot
	persistCalled bool
	nextIDs       []erpdomain.ImportID
}

func (r *contractRepo) PersistSnapshotAtomically(_ context.Context, _ string, snapshot erpdomain.ImportSnapshot) error {
	r.snapshot = snapshot
	r.persistCalled = true
	return nil
}

func (r *contractRepo) FindByFileSHA256(context.Context, string, erpdomain.FileSHA256) (*erpdomain.ImportReport, error) {
	return nil, nil
}

func (r *contractRepo) ListImports(context.Context, string) ([]erpdomain.ImportReport, error) {
	return nil, nil
}

func (r *contractRepo) GetImport(_ context.Context, _ string, _ erpdomain.ImportID) (erpdomain.ImportReport, error) {
	rejectedRows := map[int]struct{}{}
	warnings := 0
	for _, issue := range r.snapshot.Issues {
		if issue.Kind == erpdomain.Rejection {
			rejectedRows[issue.Row] = struct{}{}
		}
		if issue.Kind == erpdomain.Warning {
			warnings++
		}
	}
	return erpdomain.ImportReport{ID: r.snapshot.ID, Protocol: r.snapshot.Protocol, FileSHA256: r.snapshot.FileSHA256, Source: r.snapshot.Source, ImportedAt: r.snapshot.ImportedAt, Status: r.snapshot.Status, AcceptedCount: len(r.snapshot.AcceptedRows), RejectedCount: len(rejectedRows), WarningCount: warnings, Issues: r.snapshot.Issues}, nil
}

func (r *contractRepo) LatestCompletedSnapshot(context.Context, string, erpdomain.ImportSource) (erpdomain.ImportSnapshot, error) {
	if r.snapshot.Status != erpdomain.ImportStatusCompleted {
		return erpdomain.ImportSnapshot{}, erpports.ErrImportNotFound
	}
	return r.snapshot, nil
}

func (r *contractRepo) SyncLatestCompletedSnapshot(context.Context, string, erpdomain.ImportSource) (int, error) {
	return 0, nil
}
func (r *contractRepo) MirrorRows(context.Context, string, erpdomain.ImportSource) ([]erpdomain.MirrorProduct, error) {
	return nil, nil
}
func (r *contractRepo) MirrorProductByCode(context.Context, string, erpdomain.ImportSource, string) (erpdomain.MirrorProduct, bool, error) {
	return erpdomain.MirrorProduct{}, false, nil
}
func (r *contractRepo) MirrorCatalogPage(context.Context, string, erpdomain.ImportSource, string, int64, int) ([]erpdomain.MirrorProduct, error) {
	return nil, nil
}
func (r *contractRepo) MirrorEANCollisionCounts(context.Context, string, erpdomain.ImportSource) (map[string]int, error) {
	return nil, nil
}

var _ erpports.ImportRepository = (*contractRepo)(nil)

type oracleShapedReader struct {
	candidates map[int]readdomain.ProductCandidate
	stocks     map[int]readdomain.SellableStock
	costs      map[int]readdomain.CostAsOf
	taxes      map[int]readdomain.TaxInputs
	facts      map[int64]readports.CatalogProductFact
}

func newOracleShapedReader() *oracleShapedReader {
	ean := "7890000000017"
	flags := []readdomain.QualityFlag{readdomain.QualityComplete, readdomain.QualityEANCollision}
	return &oracleShapedReader{
		candidates: map[int]readdomain.ProductCandidate{1001: {
			InternalProductID: internalIDPointer(1001), ProductID: 1001, Name: "Example Product 1001", EAN: stringPointer(ean), ReferenceCode: stringPointer("REF-1001"), NCM: stringPointer("12345678"), BrandName: stringPointer("Synthetic Brand"), IsActive: true, Source: readdomain.SourceMetadata{System: "oracle", FetchedAt: contractFetchedAt}, QualityFlags: flags,
		}},
		stocks: map[int]readdomain.SellableStock{
			1001: {ProductID: 1001, Quantity: floatPointer(8), Policy: readdomain.DefaultSellableStockPolicy(), Source: readdomain.SourceMetadata{System: "oracle", FetchedAt: contractFetchedAt, ObservedAt: timePointer(contractClock)}},
			1002: {ProductID: 1002, Quantity: nil, Policy: readdomain.DefaultSellableStockPolicy(), Source: readdomain.SourceMetadata{System: "oracle", FetchedAt: contractFetchedAt, ObservedAt: timePointer(contractClock)}, QualityFlags: []readdomain.QualityFlag{readdomain.QualityMissingStock}},
		},
		costs: map[int]readdomain.CostAsOf{1001: {ProductID: 1001, CompanyID: 1, Basis: readdomain.CostBasisCUSSEMICM, EffectiveAt: contractClock, Amount: floatPointer(13.34), AmountScope: readdomain.CostAmountScopePerUnit, Source: readdomain.SourceMetadata{System: "oracle", FetchedAt: contractFetchedAt, ObservedAt: timePointer(contractClock)}}},
		taxes: map[int]readdomain.TaxInputs{1001: {ProductID: 1001, EffectiveAt: contractClock, NCM: stringPointer("12345678"), Source: readdomain.SourceMetadata{System: "oracle", FetchedAt: contractFetchedAt, ObservedAt: timePointer(contractClock)}}},
		facts: map[int64]readports.CatalogProductFact{1001: oracleFact(1001, ean, flags, 8), 1002: oracleFact(1002, ean, flags, 0)},
	}
}

func (r *oracleShapedReader) FindProductsForLinking(_ context.Context, input readports.FindProductsInput) ([]readdomain.ProductCandidate, error) {
	if input.ProductID != nil {
		candidate, ok := r.candidates[*input.ProductID]
		if !ok {
			return nil, errors.New("oracle-shaped product not found")
		}
		return []readdomain.ProductCandidate{candidate}, nil
	}
	return nil, nil
}
func (r *oracleShapedReader) GetSellableStock(_ context.Context, input readports.SellableStockInput) (readdomain.SellableStock, error) {
	result := r.stocks[input.ProductID]
	result.Policy = input.Policy
	return result, nil
}
func (*oracleShapedReader) GetCurrentPrice(_ context.Context, input readports.CurrentPriceInput) (readdomain.CurrentPrice, error) {
	// Oracle serves current price by default (oracle/reader.go queries VW_PRECO_TABELA). Model that
	// so the divergence assertion compares against real Oracle behavior, not a fabricated unsupported.
	return readdomain.CurrentPrice{ProductID: input.ProductID, Amount: floatPointer(29.90), Source: readdomain.SourceMetadata{System: "oracle", FetchedAt: contractFetchedAt}}, nil
}
func (r *oracleShapedReader) GetCostAsOf(_ context.Context, input readports.CostAsOfInput) (readdomain.CostAsOf, error) {
	result := r.costs[input.ProductID]
	result.CompanyID = input.Policy.CompanyID
	result.Basis = input.Policy.Basis
	result.EffectiveAt = input.Policy.EffectiveAt
	return result, nil
}
func (*oracleShapedReader) GetSalesHistory(_ context.Context, _ readports.SalesHistoryInput) (readdomain.SalesHistory, error) {
	// Oracle serves sales history by default (oracle/reader.go:281). Model served data so the
	// divergence assertion reflects the real Oracle contract, not a copy of the xlsx unsupported error.
	return readdomain.SalesHistory{Entries: []readdomain.SalesHistoryEntry{{}}, Source: readdomain.SourceMetadata{System: "oracle", FetchedAt: contractFetchedAt}}, nil
}
func (r *oracleShapedReader) GetTaxInputs(_ context.Context, input readports.TaxInput) (readdomain.TaxInputs, error) {
	result := r.taxes[input.ProductID]
	result.EffectiveAt = input.Policy.EffectiveAt
	result.IncidenceCode = input.Policy.IncidenceCode
	result.SourceIdentity = input.Policy.Source
	return result, nil
}
func (r *oracleShapedReader) ListCatalogProductFacts(_ context.Context, cursor readports.Cursor, limit int) (readports.CatalogFactPage, error) {
	if cursor.InternalProductID >= 1001 || limit <= 0 {
		return readports.CatalogFactPage{Items: []readports.CatalogProductFact{}}, nil
	}
	page := readports.CatalogFactPage{Items: []readports.CatalogProductFact{r.facts[1001]}, AsOf: contractClock}
	if limit == 1 {
		page.NextCursor = &readports.Cursor{InternalProductID: 1001}
	}
	return page, nil
}
func (r *oracleShapedReader) SearchCatalogProductFacts(_ context.Context, query string, _ int) (readports.CatalogFactPage, error) {
	if !strings.Contains(strings.ToLower("Example Product 1001"), strings.ToLower(query)) {
		return readports.CatalogFactPage{Items: []readports.CatalogProductFact{}, AsOf: contractClock}, nil
	}
	return readports.CatalogFactPage{Items: []readports.CatalogProductFact{r.facts[1001]}, AsOf: contractClock}, nil
}

var _ readports.Reader = (*oracleShapedReader)(nil)
var _ readports.CatalogPageReader = (*oracleShapedReader)(nil)

func oracleFact(id int64, ean string, flags []readdomain.QualityFlag, quantity float64) readports.CatalogProductFact {
	fact := readports.CatalogProductFact{InternalProductID: id, Reference: stringPointer("REF-1001"), ManufacturerReference: stringPointer("REF-1001"), Description: stringPointer("Example Product 1001"), EAN: stringPointer(ean), BrandName: stringPointer("Synthetic Brand"), NCM: stringPointer("12345678"), Active: true, SellableStock: readports.CatalogQuantityFact{Quantity: floatPointer(quantity)}, CurrentPrice: readports.CatalogMoneyFact{Currency: "BRL", Quality: []string{string(readdomain.QualityMissingPrice)}}, Cost: readports.CatalogMoneyFact{Amount: stringPointer("13.34"), Currency: "BRL"}}
	for _, flag := range flags {
		fact.QualityFlags = append(fact.QualityFlags, string(flag))
	}
	if id == 1002 {
		fact.Reference, fact.ManufacturerReference, fact.Description = stringPointer("REF-1002"), stringPointer("REF-1002"), stringPointer("Example Product 1002")
		fact.SellableStock.Quantity = nil
		fact.SellableStock.Quality = []string{string(readdomain.QualityMissingStock)}
	}
	return fact
}

func assertCandidateShapeParity(t *testing.T, left, right []readdomain.ProductCandidate) {
	t.Helper()
	if len(left) != len(right) {
		t.Fatalf("candidate lengths = %d/%d", len(left), len(right))
	}
	for i := range left {
		left[i].Source.System = ""
		right[i].Source.System = ""
		if !reflect.DeepEqual(left[i], right[i]) {
			t.Fatalf("candidate shape mismatch:\n xlsx=%+v\n oracle=%+v", left[i], right[i])
		}
	}
}
func assertStockShapeParity(t *testing.T, left, right readdomain.SellableStock) {
	t.Helper()
	left.Source.System, right.Source.System = "", ""
	if !reflect.DeepEqual(left, right) {
		t.Fatalf("stock shape mismatch:\n xlsx=%+v\n oracle=%+v", left, right)
	}
}
func assertCostShapeParity(t *testing.T, left, right readdomain.CostAsOf) {
	t.Helper()
	left.Source.System, right.Source.System = "", ""
	if !reflect.DeepEqual(left, right) {
		t.Fatalf("cost shape mismatch:\n xlsx=%+v\n oracle=%+v", left, right)
	}
}
func assertTaxShapeParity(t *testing.T, left, right readdomain.TaxInputs) {
	t.Helper()
	left.Source.System, right.Source.System = "", ""
	if !reflect.DeepEqual(left, right) {
		t.Fatalf("tax shape mismatch:\n xlsx=%+v\n oracle=%+v", left, right)
	}
}
func assertCatalogPageShapeParity(t *testing.T, left, right readports.CatalogFactPage) {
	t.Helper()
	if !reflect.DeepEqual(left, right) {
		t.Fatalf("catalog page shape mismatch:\n xlsx=%+v\n oracle=%+v", left, right)
	}
}
func assertIssueKind(t *testing.T, issues []erpdomain.Issue, code erpdomain.IssueCode, kind erpdomain.IssueKind) {
	t.Helper()
	for _, issue := range issues {
		if issue.Code == code {
			if issue.Kind != kind {
				t.Fatalf("issue %q kind = %q, want %q", code, issue.Kind, kind)
			}
			return
		}
	}
	t.Fatalf("issue %q not found in %#v", code, issues)
}

func acceptedRowByID(rows []erpdomain.NormalizedRow, id string) *erpdomain.NormalizedRow {
	for i := range rows {
		if rows[i].Codprod == id {
			return &rows[i]
		}
	}
	return nil
}

func contractFixturesDir(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	dir := filepath.Dir(file)
	for {
		if _, err := os.Stat(filepath.Join(dir, ".mnfs")); err == nil {
			return filepath.Join(dir, ".mnfs", "MIS-004-mvp-demo", "M-01-erp-xlsx-identity", "fixtures")
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("could not locate repository .mnfs directory")
		}
		dir = parent
	}
}

func writeContractFixtures(t *testing.T, dir string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	exampleRows := [][]string{{"CODPROD", "DESCRPROD", "CUSTO", "ESTOQUE_FISICO", "ESTOQUE_RESERVADO", "EAN", "REFFORN", "MARCA", "NCM"}}
	for id := 1001; id <= 1055; id++ {
		reserved := ""
		if id == 1001 {
			reserved = "2"
		}
		ean := validGTIN(id - 1000)
		if id == 1002 {
			ean = validGTIN(1)
		}
		exampleRows = append(exampleRows, []string{fmt.Sprintf("%d", id), fmt.Sprintf("Example Product %d", id), fmt.Sprintf("%d,34", id-988), fmt.Sprintf("%d", id-991), reserved, ean, fmt.Sprintf("REF-%d", id), "Synthetic Brand", "12345678"})
	}
	writeFixture(t, filepath.Join(dir, "example-erp.xlsx"), newContractWorkbook(exampleRows[1:]))

	rejectionRows := [][]string{
		{"2001", "", "10,00", "5", "1", validGTIN(201), "REF-EMPTY-2001", "Synthetic Brand", "12345678"},
		{"2002", "Invalid EAN Product", "11,00", "6", "1", "7894900011518", "REF-RETAIN-2002", "Synthetic Brand", "12345678"},
		{"2003", "Malformed NCM Product", "12,00", "7", "2", validGTIN(203), "REF-NCM-2003", "Synthetic Brand", "12AB"},
		{"2004", "Valid Identity Product", "13,00", "8", "3", validGTIN(204), "REF-VALID-2004", "Synthetic Brand", "12345678"},
	}
	writeFixture(t, filepath.Join(dir, "identity-rejections.xlsx"), newContractWorkbook(rejectionRows))
}

func writeFixture(t *testing.T, path string, data []byte) {
	t.Helper()
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatalf("write fixture %s: %v", path, err)
	}
}

func newContractWorkbook(rows [][]string) []byte {
	rows = append([][]string{{"CODPROD", "DESCRPROD", "CUSTO", "ESTOQUE_FISICO", "ESTOQUE_RESERVADO", "EAN", "REFFORN", "MARCA", "NCM"}}, rows...)
	var output bytes.Buffer
	archive := zip.NewWriter(&output)
	writeZipEntry(&archive, "[Content_Types].xml", `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types"><Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/><Default Extension="xml" ContentType="application/xml"/><Override PartName="/xl/workbook.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.sheet.main+xml"/><Override PartName="/xl/worksheets/sheet1.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.worksheet+xml"/></Types>`)
	writeZipEntry(&archive, "_rels/.rels", `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"><Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="xl/workbook.xml"/></Relationships>`)
	writeZipEntry(&archive, "xl/_rels/workbook.xml.rels", `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"><Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/worksheet" Target="worksheets/sheet1.xml"/></Relationships>`)
	writeZipEntry(&archive, "xl/workbook.xml", `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><workbook xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships"><sheets><sheet name="ERP" sheetId="1" r:id="rId1"/></sheets></workbook>`)
	writeZipEntry(&archive, "xl/worksheets/sheet1.xml", contractSheetXML(rows))
	if err := archive.Close(); err != nil {
		panic(err)
	}
	return output.Bytes()
}

func writeZipEntry(archive **zip.Writer, name, content string) {
	file, err := (*archive).Create(name)
	if err != nil {
		panic(err)
	}
	if _, err := io.WriteString(file, content); err != nil {
		panic(err)
	}
}

func contractSheetXML(rows [][]string) string {
	var sheet strings.Builder
	sheet.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?><worksheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"><sheetData>`)
	for rowIndex, row := range rows {
		sheet.WriteString(`<row r="` + itoaContract(rowIndex+1) + `">`)
		for columnIndex, value := range row {
			if value == "" {
				continue
			}
			sheet.WriteString(`<c r="` + cellReferenceContract(columnIndex, rowIndex) + `" t="inlineStr"><is><t>`)
			escapeContract(&sheet, value)
			sheet.WriteString(`</t></is></c>`)
		}
		sheet.WriteString(`</row>`)
	}
	sheet.WriteString(`</sheetData></worksheet>`)
	return sheet.String()
}

func cellReferenceContract(column, row int) string {
	return string(rune('A'+column)) + itoaContract(row+1)
}
func escapeContract(builder *strings.Builder, value string) {
	_ = xml.EscapeText(builder, []byte(value))
}
func itoaContract(value int) string { return fmt.Sprintf("%d", value) }

func validGTIN(seed int) string {
	body := fmt.Sprintf("789%09d", seed)
	sum := 0
	for i := len(body) - 1; i >= 0; i-- {
		weight := 1
		if (len(body)-1-i)%2 == 0 {
			weight = 3
		}
		sum += int(body[i]-'0') * weight
	}
	return body + fmt.Sprintf("%d", (10-sum%10)%10)
}

func intPointer(value int) *int { return &value }
func internalIDPointer(value int) *readdomain.InternalProductID {
	id := readdomain.InternalProductID(value)
	return &id
}
func stringPointer(value string) *string     { return &value }
func floatPointer(value float64) *float64    { return &value }
func timePointer(value time.Time) *time.Time { return &value }
