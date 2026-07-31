package transport

import (
	"context"
	"errors"
	"io"
	"net/http"
	"time"

	"marketplace-central/apps/server_core/internal/modules/erp_import/domain"
	"marketplace-central/apps/server_core/internal/modules/erp_import/ports"
	"marketplace-central/apps/server_core/internal/platform/apierror"
	"marketplace-central/apps/server_core/internal/platform/httpx"
)

const maxUploadBytes = 25 << 20

type ImportRunner interface {
	RunImport(ctx context.Context, source io.Reader) (domain.ImportReport, error)
	RunImportLenient(ctx context.Context, source io.Reader) (domain.ImportReport, error)
}

type ImportQuerier interface {
	ListImports(ctx context.Context) ([]domain.ImportReport, error)
	GetImport(ctx context.Context, id domain.ImportID) (domain.ImportReport, error)
	GetImportChain(ctx context.Context, id domain.ImportID) (domain.ImportChain, error)
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
	registerInteractiveRoute(mux, "/erp/imports/{id}/chain", h.handleGetImportChain)
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
		apierror.Write(w, http.StatusBadRequest, "invalid_file", importErrorMessage("invalid_file", ""), nil)
		return
	}
	file, _, err := r.FormFile("file")
	if err != nil {
		apierror.Write(w, http.StatusBadRequest, "invalid_file", importErrorMessage("invalid_file", ""), nil)
		return
	}
	defer file.Close()

	run := h.importer.RunImport
	if isLenientImportRequest(r) {
		run = h.importer.RunImportLenient
	}
	report, err := run(r.Context(), file)
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

// isLenientImportRequest opts into the client-catalog (lenient) import path
// when the multipart form declares source=catalogo_cliente or mode=lenient.
// Default (no field, or any other value) stays on the strict path.
func isLenientImportRequest(r *http.Request) bool {
	if r.MultipartForm == nil {
		return false
	}
	return r.FormValue("source") == "catalogo_cliente" || r.FormValue("mode") == "lenient"
}

func (h Handler) handleListImports(w http.ResponseWriter, r *http.Request) {
	reports, err := h.queries.ListImports(r.Context())
	if err != nil {
		apierror.Write(w, http.StatusInternalServerError, "internal_error", importErrorMessage("internal_error", ""), nil)
		return
	}
	items := make([]importSummaryResponse, 0, len(reports))
	for _, report := range reports {
		items = append(items, newImportSummaryResponse(report))
	}
	httpx.WriteJSON(w, http.StatusOK, importListResponse{Items: items})
}

// readImportID validates the {id} path value before it can reach a query.
// Both {id} routes share this one site so the two cannot drift apart: an
// invalid uuid used to travel into the `uuid` column and come back as a 500.
func readImportID(w http.ResponseWriter, r *http.Request) (domain.ImportID, bool) {
	importID := domain.ImportID(r.PathValue("id"))
	if !importID.IsValid() {
		apierror.Write(w, http.StatusBadRequest, "invalid_import_id", importErrorMessage("invalid_import_id", ""), nil)
		return "", false
	}
	return importID, true
}

func (h Handler) handleGetImport(w http.ResponseWriter, r *http.Request) {
	importID, ok := readImportID(w, r)
	if !ok {
		return
	}
	report, err := h.queries.GetImport(r.Context(), importID)
	if err != nil {
		if errors.Is(err, ports.ErrImportNotFound) {
			apierror.Write(w, http.StatusNotFound, "import_not_found", importErrorMessage("import_not_found", ""), nil)
			return
		}
		apierror.Write(w, http.StatusInternalServerError, "internal_error", importErrorMessage("internal_error", ""), nil)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, newImportDetailResponse(report))
}

func (h Handler) handleGetImportChain(w http.ResponseWriter, r *http.Request) {
	importID, ok := readImportID(w, r)
	if !ok {
		return
	}
	chain, err := h.queries.GetImportChain(r.Context(), importID)
	if err != nil {
		if errors.Is(err, ports.ErrImportNotFound) {
			apierror.Write(w, http.StatusNotFound, "import_not_found", importErrorMessage("import_not_found", ""), nil)
			return
		}
		apierror.Write(w, http.StatusInternalServerError, "internal_error", importErrorMessage("internal_error", ""), nil)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, importChainResponse{
		Protocol:     string(chain.Protocol),
		Importados:   chain.Importados,
		Vinculados:   chain.Vinculados,
		Enfileirados: chain.Enfileirados,
		QueueReadAt:  chain.QueueReadAt.UTC().Format(time.RFC3339),
	})
}

func writeImportError(w http.ResponseWriter, err error) {
	var fileErr *ports.FileError
	if errors.As(err, &fileErr) {
		switch fileErr.Code {
		case domain.CodeInvalidFile:
			apierror.Write(w, http.StatusBadRequest, "invalid_file", importErrorMessage("invalid_file", fileErr.Detail), nil)
		case domain.CodeMissingRequiredColumn:
			apierror.Write(w, http.StatusUnprocessableEntity, "missing_required_column", "coluna obrigatória ausente", map[string]any{"column": fileErr.Column})
		default:
			apierror.Write(w, http.StatusInternalServerError, "internal_error", importErrorMessage("internal_error", ""), nil)
		}
		return
	}

	var duplicateErr *ports.DuplicateFileError
	if errors.As(err, &duplicateErr) {
		apierror.Write(w, http.StatusConflict, "duplicate_file", "arquivo já importado", map[string]any{
			"import_id": string(duplicateErr.ExistingID),
			"protocol":  string(duplicateErr.ExistingProtocol),
		})
		return
	}
	if errors.Is(err, ports.ErrImportInProgress) {
		apierror.Write(w, http.StatusConflict, "import_in_progress", importErrorMessage("import_in_progress", ""), nil)
		return
	}
	apierror.Write(w, http.StatusInternalServerError, "internal_error", importErrorMessage("internal_error", ""), nil)
}

// importErrorMessage answers the human-facing text for writeError's callers.
// Most of them wrote no message at all under the flat body (writeError only
// set "detail" when the caller passed one), so the code-keyed defaults below
// are new operator-facing text, in pt-BR to match the module's existing
// erp_source/duplicate-file strings. A caller that already had real text
// (fileErr.Detail) keeps it verbatim.
func importErrorMessage(code, message string) string {
	if message != "" {
		return message
	}
	switch code {
	case "invalid_file":
		return "arquivo inválido"
	case "internal_error":
		return "erro interno do servidor"
	case "invalid_import_id":
		return "id de importação inválido"
	case "import_not_found":
		return "importação não encontrada"
	case "import_in_progress":
		return "importação em andamento"
	}
	return code
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

type importChainResponse struct {
	Protocol     string `json:"protocol"`
	Importados   int    `json:"importados"`
	Vinculados   int    `json:"vinculados"`
	Enfileirados int    `json:"enfileirados"`
	QueueReadAt  string `json:"queue_read_at"`
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
