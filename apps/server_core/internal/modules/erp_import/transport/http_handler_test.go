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
	items   []domain.ImportReport
	item    domain.ImportReport
	listErr error
	getErr  error
}

func (f *fakeImportQuerier) ListImports(context.Context) ([]domain.ImportReport, error) {
	return f.items, f.listErr
}

func (f *fakeImportQuerier) GetImport(context.Context, domain.ImportID) (domain.ImportReport, error) {
	return f.item, f.getErr
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
			var body map[string]any
			decodeJSON(t, response, &body)
			if body["error"] != tc.wantError {
				t.Fatalf("error = %v, want %q", body["error"], tc.wantError)
			}
			if tc.wantColumn != "" && body["column"] != tc.wantColumn {
				t.Fatalf("column = %v, want %q", body["column"], tc.wantColumn)
			}
			if tc.name == "duplicate" {
				if body["import_id"] != string(duplicate.ExistingID) || body["protocol"] != string(duplicate.ExistingProtocol) {
					t.Fatalf("duplicate body = %v, want original import", body)
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
			var body map[string]any
			decodeJSON(t, response, &body)
			if body["error"] != "invalid_file" {
				t.Fatalf("error = %v, want invalid_file", body["error"])
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
	response := performRequest(t, &fakeImportRunner{}, &fakeImportQuerier{item: report}, httptest.NewRequest(http.MethodGet, "/erp/imports/imp-1", nil))
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
	emptyResponse := performRequest(t, &fakeImportRunner{}, &fakeImportQuerier{item: emptyReport}, httptest.NewRequest(http.MethodGet, "/erp/imports/empty", nil))
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
	response := performRequest(t, &fakeImportRunner{}, &fakeImportQuerier{getErr: ports.ErrImportNotFound}, httptest.NewRequest(http.MethodGet, "/erp/imports/missing", nil))
	if response.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", response.Code)
	}
	var body map[string]any
	decodeJSON(t, response, &body)
	if body["error"] != "import_not_found" {
		t.Fatalf("error = %v, want import_not_found", body["error"])
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
	var decoded map[string]any
	decodeJSON(t, response, &decoded)
	if decoded["error"] != "invalid_file" {
		t.Fatalf("error = %v, want invalid_file", decoded["error"])
	}
}

func TestHandlerImportedAtIsRFC3339UTC(t *testing.T) {
	// Fixture ImportedAt is 12:34:56 in BRT (-03:00); the wire must carry its UTC instant.
	report := importReport("imp-1", "proto-1", domain.ImportStatusCompleted)
	response := performRequest(t, &fakeImportRunner{}, &fakeImportQuerier{item: report}, httptest.NewRequest(http.MethodGet, "/erp/imports/imp-1", nil))
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
