package transport

import (
	"context"
	"errors"
	"io"
	"net/http"
	"time"

	"marketplace-central/apps/server_core/internal/modules/erp_import/domain"
	"marketplace-central/apps/server_core/internal/modules/erp_import/ports"
	"marketplace-central/apps/server_core/internal/platform/httpx"
)

const maxUploadBytes = 25 << 20

type ImportRunner interface {
	RunImport(ctx context.Context, source io.Reader) (domain.ImportReport, error)
}

type ImportQuerier interface {
	ListImports(ctx context.Context) ([]domain.ImportReport, error)
	GetImport(ctx context.Context, id domain.ImportID) (domain.ImportReport, error)
}

type Handler struct {
	importer ImportRunner
	queries  ImportQuerier
}

func NewHandler(importer ImportRunner, queries ImportQuerier) Handler {
	return Handler{importer: importer, queries: queries}
}

func (h Handler) Register(mux httpx.RouteRegistrar) {
	if registrar, ok := mux.(routeClassRegistrar); ok {
		registrar.RegisterRouteClass("/erp/imports", httpx.BatchRouteClass)
	}
	mux.HandleFunc("POST /erp/imports", h.handlePostImport)
	mux.HandleFunc("GET /erp/imports", h.handleListImports)
	registerInteractiveRoute(mux, "/erp/imports/{id}", h.handleGetImport)
}

type routeClassRegistrar interface {
	RegisterRouteClass(string, httpx.RouteClass)
}

func registerInteractiveRoute(mux httpx.RouteRegistrar, pattern string, handler func(http.ResponseWriter, *http.Request)) {
	if registrar, ok := mux.(routeClassRegistrar); ok {
		registrar.RegisterRouteClass(pattern, httpx.InteractiveRouteClass)
	}
	mux.HandleFunc("GET "+pattern, handler)
}

func (h Handler) handlePostImport(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadBytes)
	if err := r.ParseMultipartForm(maxUploadBytes); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_file", "")
		return
	}
	file, _, err := r.FormFile("file")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_file", "")
		return
	}
	defer file.Close()

	report, err := h.importer.RunImport(r.Context(), file)
	if err != nil {
		writeImportError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, importCreatedResponse{
		ImportID: string(report.ID),
		Protocol: string(report.Protocol),
		Status:   string(report.Status),
	})
}

func (h Handler) handleListImports(w http.ResponseWriter, r *http.Request) {
	reports, err := h.queries.ListImports(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "")
		return
	}
	items := make([]importSummaryResponse, 0, len(reports))
	for _, report := range reports {
		items = append(items, newImportSummaryResponse(report))
	}
	httpx.WriteJSON(w, http.StatusOK, importListResponse{Items: items})
}

func (h Handler) handleGetImport(w http.ResponseWriter, r *http.Request) {
	report, err := h.queries.GetImport(r.Context(), domain.ImportID(r.PathValue("id")))
	if err != nil {
		if errors.Is(err, ports.ErrImportNotFound) {
			writeError(w, http.StatusNotFound, "import_not_found", "")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal_error", "")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, newImportDetailResponse(report))
}

func writeImportError(w http.ResponseWriter, err error) {
	var fileErr *ports.FileError
	if errors.As(err, &fileErr) {
		switch fileErr.Code {
		case domain.CodeInvalidFile:
			writeError(w, http.StatusBadRequest, "invalid_file", fileErr.Detail)
		case domain.CodeMissingRequiredColumn:
			httpx.WriteJSON(w, http.StatusUnprocessableEntity, map[string]string{
				"error":  "missing_required_column",
				"column": fileErr.Column,
			})
		default:
			writeError(w, http.StatusInternalServerError, "internal_error", "")
		}
		return
	}

	var duplicateErr *ports.DuplicateFileError
	if errors.As(err, &duplicateErr) {
		httpx.WriteJSON(w, http.StatusConflict, map[string]string{
			"error":     "duplicate_file",
			"import_id": string(duplicateErr.ExistingID),
			"protocol":  string(duplicateErr.ExistingProtocol),
		})
		return
	}
	if errors.Is(err, ports.ErrImportInProgress) {
		writeError(w, http.StatusConflict, "import_in_progress", "")
		return
	}
	writeError(w, http.StatusInternalServerError, "internal_error", "")
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	payload := map[string]string{"error": code}
	if message != "" {
		payload["detail"] = message
	}
	httpx.WriteJSON(w, status, payload)
}

type importCreatedResponse struct {
	ImportID string `json:"import_id"`
	Protocol string `json:"protocol"`
	Status   string `json:"status"`
}

type importListResponse struct {
	Items []importSummaryResponse `json:"items"`
}

type importSummaryResponse struct {
	ImportID      string `json:"import_id"`
	Protocol      string `json:"protocol"`
	FileSHA256    string `json:"file_sha256"`
	Source        string `json:"source"`
	ImportedAt    string `json:"imported_at"`
	Status        string `json:"status"`
	AcceptedCount int    `json:"accepted_count"`
	RejectedCount int    `json:"rejected_count"`
	WarningCount  int    `json:"warning_count"`
}

type importDetailResponse struct {
	importSummaryResponse
	RejectedRows []issueResponse `json:"rejected_rows"`
	Warnings     []issueResponse `json:"warnings"`
}

type issueResponse struct {
	Row            int     `json:"row"`
	Code           string  `json:"code"`
	Detail         string  `json:"detail"`
	Column         *string `json:"column,omitempty"`
	OffendingValue *string `json:"offending_value,omitempty"`
}

func newImportSummaryResponse(report domain.ImportReport) importSummaryResponse {
	return importSummaryResponse{
		ImportID:      string(report.ID),
		Protocol:      string(report.Protocol),
		FileSHA256:    string(report.FileSHA256),
		Source:        string(report.Source),
		ImportedAt:    report.ImportedAt.UTC().Format(time.RFC3339),
		Status:        string(report.Status),
		AcceptedCount: report.AcceptedCount,
		RejectedCount: report.RejectedCount,
		WarningCount:  report.WarningCount,
	}
}

func newImportDetailResponse(report domain.ImportReport) importDetailResponse {
	rejected := make([]issueResponse, 0)
	warnings := make([]issueResponse, 0)
	for _, issue := range report.Issues {
		response := issueResponse{
			Row:            issue.Row,
			Code:           string(issue.Code),
			Detail:         issue.Detail,
			Column:         issue.Column,
			OffendingValue: issue.OffendingValue,
		}
		switch issue.Kind {
		case domain.Rejection:
			rejected = append(rejected, response)
		case domain.Warning:
			warnings = append(warnings, response)
		}
	}
	return importDetailResponse{
		importSummaryResponse: newImportSummaryResponse(report),
		RejectedRows:          rejected,
		Warnings:              warnings,
	}
}
