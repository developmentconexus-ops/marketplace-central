package transport

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"marketplace-central/apps/server_core/internal/modules/orders/application"
	"marketplace-central/apps/server_core/internal/modules/orders/domain"
	"marketplace-central/apps/server_core/internal/platform/httpx"
)

const maxAssistedSankhyaRequestBytes = 64 << 10

type AssistedSankhyaLinkageApplication interface {
	GetCurrent(context.Context, application.GetAssistedSankhyaLinkageInput) (domain.SankhyaLinkage, error)
	ListCandidates(context.Context, application.ListAssistedSankhyaCandidatesInput) (domain.AssistedSankhyaCandidateResult, error)
	Confirm(context.Context, application.ConfirmAssistedSankhyaLinkageInput) (domain.AssistedSankhyaConfirmationResult, error)
}

type SankhyaLinkageHandler struct {
	service  AssistedSankhyaLinkageApplication
	tenantID string
}

func NewSankhyaLinkageHandler(service AssistedSankhyaLinkageApplication, tenantID string) SankhyaLinkageHandler {
	return SankhyaLinkageHandler{service: service, tenantID: strings.TrimSpace(tenantID)}
}

func (h SankhyaLinkageHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /orders/{provider_order_id}/sankhya-linkage", h.handleCurrent)
	mux.HandleFunc("GET /orders/{provider_order_id}/sankhya-linkage/candidates", h.handleCandidates)
	mux.HandleFunc("POST /orders/{provider_order_id}/sankhya-linkage/confirm", h.handleConfirm)
}

func (h SankhyaLinkageHandler) handleCurrent(w http.ResponseWriter, r *http.Request) {
	scope, ok := h.requestScope(w, r)
	if !ok {
		return
	}
	if h.service == nil {
		writeSankhyaLinkageError(w, &domain.AssistedSankhyaReadError{Kind: domain.AssistedSankhyaReadConfigurationInvalid})
		return
	}
	linkage, err := h.service.GetCurrent(r.Context(), application.GetAssistedSankhyaLinkageInput{
		TenantID: h.tenantID, InstallationID: scope.installationID, ProviderOrderID: scope.providerOrderID,
	})
	if err != nil {
		writeSankhyaLinkageError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, newSankhyaLinkageResponse(linkage))
}

func (h SankhyaLinkageHandler) handleCandidates(w http.ResponseWriter, r *http.Request) {
	scope, ok := h.requestScope(w, r)
	if !ok {
		return
	}
	if h.service == nil {
		writeSankhyaLinkageError(w, &domain.AssistedSankhyaReadError{Kind: domain.AssistedSankhyaReadConfigurationInvalid})
		return
	}
	result, err := h.service.ListCandidates(r.Context(), application.ListAssistedSankhyaCandidatesInput{
		TenantID: h.tenantID, InstallationID: scope.installationID, ProviderOrderID: scope.providerOrderID,
	})
	if err != nil {
		writeSankhyaLinkageError(w, err)
		return
	}
	response := sankhyaCandidateListResponse{InstallationID: result.Scope.InstallationID, ProviderOrderID: result.Scope.ProviderOrderID}
	for _, candidate := range result.Candidates {
		response.Candidates = append(response.Candidates, newSankhyaCandidateResponse(candidate))
	}
	if response.Candidates == nil {
		response.Candidates = []sankhyaCandidateResponse{}
	}
	httpx.WriteJSON(w, http.StatusOK, response)
}

func (h SankhyaLinkageHandler) handleConfirm(w http.ResponseWriter, r *http.Request) {
	scope, ok := h.requestScope(w, r)
	if !ok {
		return
	}
	if h.service == nil {
		writeSankhyaLinkageError(w, &domain.AssistedSankhyaReadError{Kind: domain.AssistedSankhyaReadConfigurationInvalid})
		return
	}
	var request sankhyaConfirmRequest
	if err := decodeStrictJSON(w, r, &request); err != nil {
		writeSankhyaLinkageInvalid(w)
		return
	}
	sourceAt, err := time.Parse(time.RFC3339Nano, request.SourceAt)
	if err != nil || sourceAt.IsZero() {
		writeSankhyaLinkageInvalid(w)
		return
	}
	selections := make([]domain.AssistedSankhyaLineSelection, 0, len(request.Selections))
	for _, selection := range request.Selections {
		selections = append(selections, domain.AssistedSankhyaLineSelection{
			MPCLineID: domain.MPCLineID(selection.MPCLineID),
			CandidateLine: domain.InternalDocumentLineIdentity{
				DocumentID: selection.InternalDocumentID,
				LineNumber: selection.InternalLineNumber,
			},
		})
	}
	result, err := h.service.Confirm(r.Context(), application.ConfirmAssistedSankhyaLinkageInput{
		TenantID: h.tenantID, InstallationID: scope.installationID, ProviderOrderID: scope.providerOrderID,
		SelectedDocumentID: request.SelectedDocumentID, Selections: selections,
		ActorID: request.ActorID, Reason: request.Reason, SourceAt: sourceAt, IdempotencyKey: request.IdempotencyKey,
	})
	if err != nil {
		writeSankhyaLinkageError(w, err)
		return
	}
	response := sankhyaConfirmationResponse{sankhyaLinkageResponse: newSankhyaLinkageResponse(result.Linkage)}
	for _, lineage := range result.Lines {
		response.Lineage = append(response.Lineage, newSankhyaLineageResponse(lineage))
	}
	if response.Lineage == nil {
		response.Lineage = []sankhyaLineageResponse{}
	}
	httpx.WriteJSON(w, http.StatusOK, response)
}

type sankhyaRequestScope struct {
	installationID  string
	providerOrderID string
}

func (h SankhyaLinkageHandler) requestScope(w http.ResponseWriter, r *http.Request) (sankhyaRequestScope, bool) {
	query := r.URL.Query()
	installationValues, present := query["installation_id"]
	installationID := ""
	if present && len(installationValues) == 1 {
		installationID = installationValues[0]
	}
	providerOrderID := r.PathValue("provider_order_id")
	if h.tenantID == "" || len(query) != 1 || !validSankhyaScopeValue(installationID) || !validSankhyaScopeValue(providerOrderID) {
		writeSankhyaLinkageInvalid(w)
		return sankhyaRequestScope{}, false
	}
	return sankhyaRequestScope{installationID: installationID, providerOrderID: providerOrderID}, true
}

func validSankhyaScopeValue(value string) bool {
	return value != "" && value == strings.TrimSpace(value) && len(value) <= 200 && !strings.ContainsAny(value, "\r\n\t")
}

func decodeStrictJSON(w http.ResponseWriter, r *http.Request, target any) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxAssistedSankhyaRequestBytes)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return err
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return errors.New("request body must contain exactly one JSON value")
	}
	return nil
}

func writeSankhyaLinkageInvalid(w http.ResponseWriter) {
	writeOrdersError(w, http.StatusBadRequest, "ORDERS_SANKHYA_LINKAGE_INVALID_REQUEST", "invalid assisted Sankhya linkage request")
}

func writeSankhyaLinkageError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	code := "ORDERS_SANKHYA_LINKAGE_INTERNAL_ERROR"
	message := "assisted Sankhya linkage request failed"
	var readErr *domain.AssistedSankhyaReadError
	switch {
	case errors.Is(err, domain.ErrInvalidAssistedSankhyaLinkage), errors.Is(err, domain.ErrInvalidSankhyaLinkage):
		status, code, message = http.StatusBadRequest, "ORDERS_SANKHYA_LINKAGE_INVALID_REQUEST", "invalid assisted Sankhya linkage request"
	case errors.Is(err, domain.ErrAssistedSankhyaOrderNotFound), errors.Is(err, domain.ErrSankhyaLinkageNotFound):
		status, code, message = http.StatusNotFound, "ORDERS_SANKHYA_LINKAGE_NOT_FOUND", "assisted Sankhya linkage not found"
	case errors.Is(err, domain.ErrSankhyaLinkageConflict):
		status, code, message = http.StatusConflict, "ORDERS_SANKHYA_LINKAGE_CONFLICT", "assisted Sankhya linkage conflict"
	case errors.As(err, &readErr) && readErr.Kind == domain.AssistedSankhyaReadCandidateConflict:
		status, code, message = http.StatusConflict, "ORDERS_SANKHYA_LINKAGE_CONFLICT", "assisted Sankhya linkage conflict"
	case errors.As(err, &readErr):
		status, code, message = http.StatusServiceUnavailable, "ORDERS_SANKHYA_LINKAGE_UNAVAILABLE", "assisted Sankhya linkage unavailable"
	}
	writeOrdersError(w, status, code, message)
}

type sankhyaConfirmRequest struct {
	SelectedDocumentID int64                  `json:"selected_document_id"`
	Selections         []sankhyaLineSelection `json:"selections"`
	ActorID            string                 `json:"actor_id"`
	Reason             string                 `json:"reason"`
	SourceAt           string                 `json:"source_at"`
	IdempotencyKey     string                 `json:"idempotency_key"`
}

type sankhyaLineSelection struct {
	MPCLineID          string `json:"mpc_line_id"`
	InternalDocumentID int64  `json:"internal_document_id"`
	InternalLineNumber int    `json:"internal_line_number"`
}

type sankhyaCandidateListResponse struct {
	InstallationID  string                     `json:"installation_id"`
	ProviderOrderID string                     `json:"provider_order_id"`
	Candidates      []sankhyaCandidateResponse `json:"candidates"`
}

type sankhyaCandidateResponse struct {
	InternalDocumentID int64                  `json:"internal_document_id"`
	DocumentNumber     *int64                 `json:"document_number"`
	OperationCode      int64                  `json:"operation_code"`
	NegotiatedAt       *time.Time             `json:"negotiated_at"`
	Amount             *float64               `json:"amount"`
	Lines              []sankhyaCandidateLine `json:"lines"`
}

type sankhyaCandidateLine struct {
	InternalDocumentID int64    `json:"internal_document_id"`
	InternalLineNumber int      `json:"internal_line_number"`
	ProductID          *int64   `json:"product_id"`
	Quantity           *float64 `json:"quantity"`
	Amount             *float64 `json:"amount"`
}

type sankhyaLinkageResponse struct {
	InstallationID     string               `json:"installation_id"`
	ProviderOrderID    string               `json:"provider_order_id"`
	InternalDocumentID int64                `json:"internal_document_id"`
	Lines              []sankhyaMappingLine `json:"lines"`
	Audit              sankhyaAuditResponse `json:"audit"`
}

type sankhyaMappingLine struct {
	MPCLineID          string `json:"mpc_line_id"`
	InternalDocumentID int64  `json:"internal_document_id"`
	InternalLineNumber int    `json:"internal_line_number"`
}

type sankhyaAuditResponse struct {
	EventID               string    `json:"event_id"`
	EventType             string    `json:"event_type"`
	ActorType             string    `json:"actor_type"`
	ActorID               string    `json:"actor_id"`
	Reason                string    `json:"reason"`
	SourceAt              time.Time `json:"source_at"`
	RecordedAt            time.Time `json:"recorded_at"`
	ConfigurationRevision string    `json:"configuration_revision"`
	EvidenceState         string    `json:"evidence_state"`
	EvidenceReference     string    `json:"evidence_reference"`
}

type sankhyaConfirmationResponse struct {
	sankhyaLinkageResponse
	Lineage []sankhyaLineageResponse `json:"lineage"`
}

type sankhyaLineageResponse struct {
	MPCLineID          string                             `json:"mpc_line_id"`
	InternalDocumentID int64                              `json:"internal_document_id"`
	InternalLineNumber int                                `json:"internal_line_number"`
	State              domain.AssistedSankhyaLineageState `json:"state"`
	Descendants        []sankhyaDescendantResponse        `json:"descendants"`
}

type sankhyaDescendantResponse struct {
	InternalDocumentID int64    `json:"internal_document_id"`
	InternalLineNumber int      `json:"internal_line_number"`
	AttendedQuantity   *float64 `json:"attended_quantity"`
}

func newSankhyaCandidateResponse(candidate domain.AssistedSankhyaCandidate) sankhyaCandidateResponse {
	result := sankhyaCandidateResponse{
		InternalDocumentID: candidate.Header.DocumentID, DocumentNumber: candidate.DocumentNumber,
		OperationCode: candidate.OperationCode, NegotiatedAt: candidate.NegotiatedAt, Amount: candidate.Amount,
	}
	for _, line := range candidate.Lines {
		result.Lines = append(result.Lines, sankhyaCandidateLine{
			InternalDocumentID: line.Identity.DocumentID, InternalLineNumber: line.Identity.LineNumber,
			ProductID: line.ProductID, Quantity: line.Quantity, Amount: line.Amount,
		})
	}
	if result.Lines == nil {
		result.Lines = []sankhyaCandidateLine{}
	}
	return result
}

func newSankhyaLinkageResponse(linkage domain.SankhyaLinkage) sankhyaLinkageResponse {
	result := sankhyaLinkageResponse{
		InstallationID: linkage.Scope.InstallationID, ProviderOrderID: linkage.Scope.ProviderOrderID,
		InternalDocumentID: linkage.Header.DocumentID,
		Audit: sankhyaAuditResponse{
			EventID: linkage.Audit.EventID, EventType: string(linkage.Audit.EventType),
			ActorType: linkage.Audit.ActorType, ActorID: linkage.Audit.ActorID, Reason: linkage.Audit.Reason,
			SourceAt: linkage.Audit.SourceAt, RecordedAt: linkage.Audit.RecordedAt,
			ConfigurationRevision: linkage.Audit.ConfigurationRevision,
			EvidenceState:         string(linkage.Audit.EvidenceState), EvidenceReference: linkage.Audit.EvidenceReference,
		},
	}
	for _, line := range linkage.Lines {
		result.Lines = append(result.Lines, sankhyaMappingLine{
			MPCLineID: string(line.MPCLineID), InternalDocumentID: line.Origin.DocumentID, InternalLineNumber: line.Origin.LineNumber,
		})
	}
	if result.Lines == nil {
		result.Lines = []sankhyaMappingLine{}
	}
	return result
}

func newSankhyaLineageResponse(lineage domain.AssistedSankhyaLineage) sankhyaLineageResponse {
	result := sankhyaLineageResponse{
		MPCLineID: string(lineage.MPCLineID), InternalDocumentID: lineage.Origin.DocumentID,
		InternalLineNumber: lineage.Origin.LineNumber, State: lineage.State,
	}
	for _, descendant := range lineage.Descendants {
		result.Descendants = append(result.Descendants, sankhyaDescendantResponse{
			InternalDocumentID: descendant.Identity.DocumentID, InternalLineNumber: descendant.Identity.LineNumber,
			AttendedQuantity: descendant.AttendedQuantity,
		})
	}
	if result.Descendants == nil {
		result.Descendants = []sankhyaDescendantResponse{}
	}
	return result
}
