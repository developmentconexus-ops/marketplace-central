package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/modules/erp_import/domain"
	"marketplace-central/apps/server_core/internal/modules/erp_import/ports"
	"marketplace-central/apps/server_core/internal/platform/httpx"
)

// Both {id} routes validate the path value before the query, so request paths
// in these tests carry real UUIDs — the shape erp_import_protocols.id holds.
const (
	validImportID   = "11111111-1111-1111-1111-111111111111"
	missingImportID = "22222222-2222-2222-2222-222222222222"
	// Upper-case hex, accepted by ImportID.IsValid: a fixture whose bytes the
	// handler cannot normalise without the assertion noticing.
	mixedCaseImportID = "3A3a3A3a-BBBB-4ccc-8DDD-eeeeFFFF0000"
)

type fakeImportRunner struct {
	report        domain.ImportReport
	err           error
	lenientCalled bool
}

func (f *fakeImportRunner) RunImport(_ context.Context, _ io.Reader) (domain.ImportReport, error) {
	return f.report, f.err
}

func (f *fakeImportRunner) RunImportLenient(_ context.Context, _ io.Reader) (domain.ImportReport, error) {
	f.lenientCalled = true
	return f.report, f.err
}

type fakeImportQuerier struct {
	items    []domain.ImportReport
	item     domain.ImportReport
	chain    domain.ImportChain
	listErr  error
	getErr   error
	chainErr error
	getID    domain.ImportID
	chainID  domain.ImportID
}

func (f *fakeImportQuerier) ListImports(context.Context) ([]domain.ImportReport, error) {
	return f.items, f.listErr
}

func (f *fakeImportQuerier) GetImport(_ context.Context, id domain.ImportID) (domain.ImportReport, error) {
	f.getID = id
	return f.item, f.getErr
}

func (f *fakeImportQuerier) GetImportChain(_ context.Context, id domain.ImportID) (domain.ImportChain, error) {
	f.chainID = id
	return f.chain, f.chainErr
}

func TestHandlerPostImport(t *testing.T) {
	completed := importReport("imp-1", "proto-1", domain.ImportStatusCompleted)
	rejected := importReport("imp-2", "proto-2", domain.ImportStatusRejected)
	tests := []struct {
		name   string
		report domain.ImportReport
	}{
		{name: "completed", report: completed},
		{name: "all rejected", report: rejected},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			runner := &fakeImportRunner{report: tc.report}
			response := performRequest(t, runner, &fakeImportQuerier{}, multipartRequest(t, "payload.xlsx"))
			if response.Code != http.StatusCreated {
				t.Fatalf("status = %d, want %d", response.Code, http.StatusCreated)
			}

			var body struct {
				ImportID string `json:"import_id"`
				Protocol string `json:"protocol"`
				Status   string `json:"status"`
			}
			decodeJSON(t, response, &body)
			if body.ImportID != string(tc.report.ID) || body.Protocol != string(tc.report.Protocol) || body.Status != string(tc.report.Status) {
				t.Fatalf("body = %+v, want import_id=%q protocol=%q status=%q", body, tc.report.ID, tc.report.Protocol, tc.report.Status)
			}
		})
	}
}

// deadlineRecordingRunner reads the budget the upload actually got, from the
// context the handler hands to the importer.
type deadlineRecordingRunner struct {
	report      domain.ImportReport
	deadline    time.Time
	hasDeadline bool
}

func (r *deadlineRecordingRunner) RunImport(ctx context.Context, _ io.Reader) (domain.ImportReport, error) {
	r.deadline, r.hasDeadline = ctx.Deadline()
	return r.report, nil
}

func (r *deadlineRecordingRunner) RunImportLenient(ctx context.Context, source io.Reader) (domain.ImportReport, error) {
	return r.RunImport(ctx, source)
}

// B-08: Register declares /erp/imports as batch, so an xlsx upload has 120s.
// The route class is only honoured if the declaration and the ServeMux pattern
// agree on the key; when they did not, the upload ran on the 15s interactive
// default and a large planilha was cut off mid-import with a 504.
func TestHandlerPostImportRunsUnderTheDeclaredBatchDeadline(t *testing.T) {
	runner := &deadlineRecordingRunner{report: importReport("imp-batch", "proto-batch", domain.ImportStatusCompleted)}
	mux := httpx.NewRouteClassMux()
	NewHandler(runner, &fakeImportQuerier{}).Register(mux)

	response := httptest.NewRecorder()
	mux.ServeHTTP(response, multipartRequest(t, "payload.xlsx"))

	if response.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d; body=%q", response.Code, http.StatusCreated, response.Body.String())
	}
	if !runner.hasDeadline {
		t.Fatal("importer ran without a context deadline")
	}
	budget := time.Until(runner.deadline)
	if budget <= 15*time.Second {
		t.Fatalf("upload budget = %s, want the 120s batch budget; 15s is the interactive default", budget)
	}
}

func TestHandlerPostImportLenientTrigger(t *testing.T) {
	tests := []struct {
		name        string
		extraFields map[string]string
		wantLenient bool
	}{
		{name: "default strict", wantLenient: false},
		{name: "source catalogo_cliente", extraFields: map[string]string{"source": "catalogo_cliente"}, wantLenient: true},
		{name: "mode lenient", extraFields: map[string]string{"mode": "lenient"}, wantLenient: true},
		{name: "unrelated field", extraFields: map[string]string{"source": "erp"}, wantLenient: false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			runner := &fakeImportRunner{report: importReport("imp-1", "proto-1", domain.ImportStatusCompleted)}
			response := performRequest(t, runner, &fakeImportQuerier{}, multipartRequestWithFields(t, "payload.xlsx", tc.extraFields))
			if response.Code != http.StatusCreated {
				t.Fatalf("status = %d, want %d", response.Code, http.StatusCreated)
			}
			if runner.lenientCalled != tc.wantLenient {
				t.Fatalf("lenientCalled = %v, want %v", runner.lenientCalled, tc.wantLenient)
			}
		})
	}
}

func multipartRequestWithFields(t *testing.T, filename string, fields map[string]string) *http.Request {
	t.Helper()
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	for key, value := range fields {
		if err := writer.WriteField(key, value); err != nil {
			t.Fatal(err)
		}
	}
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := part.Write([]byte("xlsx payload")); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest(http.MethodPost, "/erp/imports", &body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	return request
}

func TestHandlerPostImportErrors(t *testing.T) {
	fileError := func(code domain.IssueCode) error {
		return &ports.FileError{Code: code, Column: "codprod", Detail: "bad workbook"}
	}
	duplicate := &ports.DuplicateFileError{ExistingID: "original-id", ExistingProtocol: "original-protocol"}
	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantError  string
		wantColumn string
	}{
		{name: "invalid file", err: fileError(domain.CodeInvalidFile), wantStatus: http.StatusBadRequest, wantError: "invalid_file"},
		{name: "missing required column", err: fileError(domain.CodeMissingRequiredColumn), wantStatus: http.StatusUnprocessableEntity, wantError: "missing_required_column", wantColumn: "codprod"},
		{name: "duplicate", err: duplicate, wantStatus: http.StatusConflict, wantError: "duplicate_file"},
		{name: "import in progress", err: ports.ErrImportInProgress, wantStatus: http.StatusConflict, wantError: "import_in_progress"},
		{name: "internal", err: errors.New("database details must not escape"), wantStatus: http.StatusInternalServerError, wantError: "internal_error"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			response := performRequest(t, &fakeImportRunner{err: tc.err}, &fakeImportQuerier{}, multipartRequest(t, "payload.xlsx"))
			if response.Code != tc.wantStatus {
				t.Fatalf("status = %d, want %d", response.Code, tc.wantStatus)
			}
			var body struct {
				Error struct {
					Code    string         `json:"code"`
					Details map[string]any `json:"details"`
				} `json:"error"`
			}
			decodeJSON(t, response, &body)
			if body.Error.Code != tc.wantError {
				t.Fatalf("error.code = %v, want %q", body.Error.Code, tc.wantError)
			}
			if tc.wantColumn != "" && body.Error.Details["column"] != tc.wantColumn {
				t.Fatalf("details.column = %v, want %q", body.Error.Details["column"], tc.wantColumn)
			}
			if tc.name == "duplicate" {
				if body.Error.Details["import_id"] != string(duplicate.ExistingID) || body.Error.Details["protocol"] != string(duplicate.ExistingProtocol) {
					t.Fatalf("duplicate details = %v, want original import", body.Error.Details)
				}
			}
		})
	}
}

func TestHandlerPostImportMalformedRequests(t *testing.T) {
	tests := []struct {
		name string
		body io.Reader
	}{
		{name: "missing file field", body: multipartFieldRequest(t)},
		{name: "non multipart", body: strings.NewReader("not multipart")},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			response := performRequest(t, &fakeImportRunner{}, &fakeImportQuerier{}, requestWithBody(t, tc.body, contentTypeFor(tc.name)))
			if response.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want %d", response.Code, http.StatusBadRequest)
			}
			var body struct {
				Error struct {
					Code string `json:"code"`
				} `json:"error"`
			}
			decodeJSON(t, response, &body)
			if body.Error.Code != "invalid_file" {
				t.Fatalf("error.code = %v, want invalid_file", body.Error.Code)
			}
		})
	}
}

func TestHandlerGetImportListPreservesServiceOrder(t *testing.T) {
	items := []domain.ImportReport{
		importReport("new", "p-new", domain.ImportStatusCompleted),
		importReport("old", "p-old", domain.ImportStatusRejected),
	}
	response := performRequest(t, &fakeImportRunner{}, &fakeImportQuerier{items: items}, httptest.NewRequest(http.MethodGet, "/erp/imports", nil))
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", response.Code)
	}
	var body struct {
		Items []struct {
			ImportID string `json:"import_id"`
			Protocol string `json:"protocol"`
		} `json:"items"`
	}
	decodeJSON(t, response, &body)
	if len(body.Items) != 2 || body.Items[0].ImportID != "new" || body.Items[1].ImportID != "old" {
		t.Fatalf("items = %+v, want service order new, old", body.Items)
	}
}

func TestHandlerGetImportDetailSplitsIssuesAndEmitsEmptyArrays(t *testing.T) {
	column := "codprod"
	report := importReport("imp-1", "proto-1", domain.ImportStatusRejected)
	report.Issues = []domain.Issue{
		{Row: 2, Kind: domain.Rejection, Code: domain.CodeEmptyCodprod, Detail: "missing product code"},
		{Row: 3, Kind: domain.Warning, Code: domain.CodeInvalidEAN, Detail: "invalid EAN", Column: &column},
	}
	response := performRequest(t, &fakeImportRunner{}, &fakeImportQuerier{item: report}, httptest.NewRequest(http.MethodGet, "/erp/imports/"+validImportID, nil))
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", response.Code)
	}
	var body struct {
		ImportID     string `json:"import_id"`
		RejectedRows []struct {
			Row    int    `json:"row"`
			Code   string `json:"code"`
			Detail string `json:"detail"`
		} `json:"rejected_rows"`
		Warnings []struct {
			Row    int    `json:"row"`
			Code   string `json:"code"`
			Detail string `json:"detail"`
			Column string `json:"column"`
		} `json:"warnings"`
	}
	decodeJSON(t, response, &body)
	if body.ImportID != "imp-1" || len(body.RejectedRows) != 1 || body.RejectedRows[0].Row != 2 || body.RejectedRows[0].Code != string(domain.CodeEmptyCodprod) || body.RejectedRows[0].Detail != "missing product code" {
		t.Fatalf("rejected_rows = %+v", body.RejectedRows)
	}
	if len(body.Warnings) != 1 || body.Warnings[0].Row != 3 || body.Warnings[0].Code != string(domain.CodeInvalidEAN) || body.Warnings[0].Column != column {
		t.Fatalf("warnings = %+v", body.Warnings)
	}

	emptyReport := importReport("empty", "p-empty", domain.ImportStatusCompleted)
	emptyResponse := performRequest(t, &fakeImportRunner{}, &fakeImportQuerier{item: emptyReport}, httptest.NewRequest(http.MethodGet, "/erp/imports/"+validImportID, nil))
	var emptyBody struct {
		RejectedRows []any `json:"rejected_rows"`
		Warnings     []any `json:"warnings"`
	}
	decodeJSON(t, emptyResponse, &emptyBody)
	if emptyBody.RejectedRows == nil || emptyBody.Warnings == nil {
		t.Fatalf("empty arrays must be []: %+v", emptyBody)
	}
}

func TestHandlerGetImportDetailUnknownID(t *testing.T) {
	response := performRequest(t, &fakeImportRunner{}, &fakeImportQuerier{getErr: ports.ErrImportNotFound}, httptest.NewRequest(http.MethodGet, "/erp/imports/"+missingImportID, nil))
	if response.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", response.Code)
	}
	var body struct {
		Error struct {
			Code string `json:"code"`
		} `json:"error"`
	}
	decodeJSON(t, response, &body)
	if body.Error.Code != "import_not_found" {
		t.Fatalf("error.code = %v, want import_not_found", body.Error.Code)
	}
}

func TestHandlerPostImportOversizeRejected(t *testing.T) {
	const boundary = "oversizeboundary"
	header := "--" + boundary + "\r\n" +
		"Content-Disposition: form-data; name=\"file\"; filename=\"big.xlsx\"\r\n" +
		"Content-Type: application/octet-stream\r\n\r\n"
	footer := "\r\n--" + boundary + "--\r\n"
	oversize := io.LimitReader(fillReader{}, int64(maxUploadBytes)+(1<<20))
	body := io.MultiReader(strings.NewReader(header), oversize, strings.NewReader(footer))
	request := requestWithBody(t, body, "multipart/form-data; boundary="+boundary)
	response := performRequest(t, &fakeImportRunner{}, &fakeImportQuerier{}, request)
	if response.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400 for oversize upload", response.Code)
	}
	var decoded struct {
		Error struct {
			Code string `json:"code"`
		} `json:"error"`
	}
	decodeJSON(t, response, &decoded)
	if decoded.Error.Code != "invalid_file" {
		t.Fatalf("error.code = %v, want invalid_file", decoded.Error.Code)
	}
}

func TestHandlerImportedAtIsRFC3339UTC(t *testing.T) {
	// Fixture ImportedAt is 12:34:56 in BRT (-03:00); the wire must carry its UTC instant.
	report := importReport("imp-1", "proto-1", domain.ImportStatusCompleted)
	response := performRequest(t, &fakeImportRunner{}, &fakeImportQuerier{item: report}, httptest.NewRequest(http.MethodGet, "/erp/imports/"+validImportID, nil))
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", response.Code)
	}
	var body struct {
		ImportedAt string `json:"imported_at"`
	}
	decodeJSON(t, response, &body)
	if body.ImportedAt != "2026-07-17T15:34:56Z" {
		t.Fatalf("imported_at = %q, want UTC RFC3339 2026-07-17T15:34:56Z", body.ImportedAt)
	}
}

func TestHandlerGetImportChainMapsServiceErrorsAndUTCResponse(t *testing.T) {
	querier := &fakeImportQuerier{chain: domain.ImportChain{
		Protocol:     "#001-E",
		Importados:   4,
		Vinculados:   2,
		Enfileirados: 3,
		QueueReadAt:  time.Date(2026, 7, 27, 12, 34, 56, 0, time.FixedZone("BRT", -3*60*60)),
	}}
	// mixedCaseImportID, not validImportID: the path value must reach the service
	// RAW. IsValid accepts upper-case hex, so an id whose case survives the round
	// trip is the only fixture where a normalising handler would be visible —
	// validImportID is all ones and would read the same either way.
	response := performRequest(t, &fakeImportRunner{}, querier, httptest.NewRequest(http.MethodGet, "/erp/imports/"+mixedCaseImportID+"/chain", nil))
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", response.Code, response.Body.String())
	}
	if querier.chainID != domain.ImportID(mixedCaseImportID) {
		t.Fatalf("service id = %q, want %q byte for byte", querier.chainID, mixedCaseImportID)
	}
	if got, want := response.Body.String(), "{\"protocol\":\"#001-E\",\"importados\":4,\"vinculados\":2,\"enfileirados\":3,\"queue_read_at\":\"2026-07-27T15:34:56Z\"}\n"; got != want {
		t.Fatalf("body = %q, want %q", got, want)
	}

	for _, tc := range []struct {
		name       string
		err        error
		wantStatus int
		wantBody   string
	}{
		{name: "not found", err: ports.ErrImportNotFound, wantStatus: http.StatusNotFound, wantBody: "{\"error\":{\"code\":\"import_not_found\",\"message\":\"importação não encontrada\",\"details\":{}}}\n"},
		{name: "internal", err: errors.New("database details must not escape"), wantStatus: http.StatusInternalServerError, wantBody: "{\"error\":{\"code\":\"internal_error\",\"message\":\"erro interno do servidor\",\"details\":{}}}\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := performRequest(t, &fakeImportRunner{}, &fakeImportQuerier{chainErr: tc.err}, httptest.NewRequest(http.MethodGet, "/erp/imports/"+missingImportID+"/chain", nil))
			if got.Code != tc.wantStatus || got.Body.String() != tc.wantBody {
				t.Fatalf("status/body = %d/%q, want %d/%q", got.Code, got.Body.String(), tc.wantStatus, tc.wantBody)
			}
		})
	}
}

func TestHandlerMalformedImportIDRejectedOnBothRoutesBeforeQuery(t *testing.T) {
	const malformedID = "not-a-uuid"
	tests := []struct {
		name string
		path string
	}{
		{name: "get import", path: "/erp/imports/" + malformedID},
		{name: "get import chain", path: "/erp/imports/" + malformedID + "/chain"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			querier := &fakeImportQuerier{}
			response := performRequest(t, &fakeImportRunner{}, querier, httptest.NewRequest(http.MethodGet, tc.path, nil))
			if response.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want 400", response.Code)
			}
			if got, want := response.Body.String(), "{\"error\":{\"code\":\"invalid_import_id\",\"message\":\"id de importação inválido\",\"details\":{}}}\n"; got != want {
				t.Fatalf("body = %q, want %q", got, want)
			}
			if querier.getID != "" {
				t.Fatalf("GetImport ran with id %q; the malformed path value must be rejected before any query", querier.getID)
			}
			if querier.chainID != "" {
				t.Fatalf("GetImportChain ran with id %q; the malformed path value must be rejected before any query", querier.chainID)
			}
		})
	}
}

type fillReader struct{}

func (fillReader) Read(p []byte) (int, error) {
	for i := range p {
		p[i] = 'a'
	}
	return len(p), nil
}

func importReport(id, protocol string, status domain.ImportStatus) domain.ImportReport {
	return domain.ImportReport{
		ID:         domain.ImportID(id),
		Protocol:   domain.Protocol(protocol),
		Source:     domain.SourceXLSX,
		ImportedAt: time.Date(2026, 7, 17, 12, 34, 56, 0, time.FixedZone("BRT", -3*60*60)),
		Status:     status,
	}
}

func multipartRequest(t *testing.T, filename string) *http.Request {
	t.Helper()
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := part.Write([]byte("xlsx payload")); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest(http.MethodPost, "/erp/imports", &body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	return request
}

func multipartFieldRequest(t *testing.T) io.Reader {
	t.Helper()
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	if err := writer.WriteField("other", "value"); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	return &body
}

func requestWithBody(t *testing.T, body io.Reader, contentType string) *http.Request {
	t.Helper()
	request := httptest.NewRequest(http.MethodPost, "/erp/imports", body)
	request.Header.Set("Content-Type", contentType)
	return request
}

func contentTypeFor(name string) string {
	if name == "missing file field" {
		return "multipart/form-data; boundary=unused"
	}
	return "text/plain"
}

func performRequest(t *testing.T, runner ImportRunner, querier ImportQuerier, request *http.Request) *httptest.ResponseRecorder {
	t.Helper()
	mux := http.NewServeMux()
	NewHandler(runner, querier).Register(mux)
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	return response
}

func decodeJSON(t *testing.T, response *httptest.ResponseRecorder, target any) {
	t.Helper()
	if err := json.NewDecoder(response.Body).Decode(target); err != nil {
		t.Fatalf("decode response JSON: %v; body=%q", err, response.Body.String())
	}
}
