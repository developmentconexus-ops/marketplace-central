package transport

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"marketplace-central/apps/server_core/internal/modules/product_links/application"
	"marketplace-central/apps/server_core/internal/modules/product_links/domain"
	"marketplace-central/apps/server_core/internal/platform/httpx"
)

type ListingSnapshotImporter interface {
	ImportListingSnapshots(ctx context.Context, input application.ImportListingSnapshotsInput) (domain.ListingSnapshotImportResult, error)
}

type LinkCandidateGenerator interface {
	GenerateLinkCandidates(ctx context.Context, input application.GenerateLinkCandidatesInput) (domain.LinkCandidateGenerationResult, error)
}

type LinkCandidateReader interface {
	ListLinkCandidates(ctx context.Context, input application.ListLinkCandidatesInput) ([]domain.LinkCandidate, error)
}

type LinkWorkflowReader interface {
	ListLinkWorkflows(ctx context.Context, input application.ListLinkWorkflowsInput) ([]domain.ProductLinkWorkflowItem, error)
}

type LinkWorkflowResolver interface {
	ApproveCandidate(ctx context.Context, input application.ApproveCandidateInput) (domain.ProductLinkResolutionResult, error)
	RejectListing(ctx context.Context, input application.RejectListingInput) (domain.ProductLinkResolutionResult, error)
	ManualResolve(ctx context.Context, input application.ManualResolveInput) (domain.ProductLinkResolutionResult, error)
}

// BatchResolver is the batch dry-run (S2) + batch-apply (S3) surface for
// product-links resolutions. Kept as its own narrow interface — additive to
// LinkWorkflowResolver — so wiring it does not require touching NewHandler's
// existing signature or its callers. A single injected *application.BatchService
// satisfies both methods, so one root.go wiring call covers preview+apply.
type BatchResolver interface {
	PreviewBatch(ctx context.Context, input application.PreviewBatchInput) (application.PreviewBatchResult, error)
	ApplyBatch(ctx context.Context, input application.ApplyBatchInput) (application.ApplyBatchResult, error)
}

type Handler struct {
	importer         ListingSnapshotImporter
	candidateMaker   LinkCandidateGenerator
	candidateReader  LinkCandidateReader
	workflowReader   LinkWorkflowReader
	workflowResolver LinkWorkflowResolver
	batchResolver    BatchResolver
}

func NewHandler(importer ListingSnapshotImporter, candidateMaker LinkCandidateGenerator, candidateReader LinkCandidateReader, workflowReader LinkWorkflowReader, workflowResolver LinkWorkflowResolver) Handler {
	return Handler{
		importer:         importer,
		candidateMaker:   candidateMaker,
		candidateReader:  candidateReader,
		workflowReader:   workflowReader,
		workflowResolver: workflowResolver,
	}
}

// NewHandlerWithBatchPreview is an additive constructor (same precedent as
// orders.NewHandlerWithReader) that also wires the batch surface (S2 preview
// + S3 apply, via the single BatchResolver), without changing NewHandler's
// existing signature/callers.
func NewHandlerWithBatchPreview(importer ListingSnapshotImporter, candidateMaker LinkCandidateGenerator, candidateReader LinkCandidateReader, workflowReader LinkWorkflowReader, workflowResolver LinkWorkflowResolver, batchResolver BatchResolver) Handler {
	h := NewHandler(importer, candidateMaker, candidateReader, workflowReader, workflowResolver)
	h.batchResolver = batchResolver
	return h
}

func (h Handler) Register(mux httpx.RouteRegistrar) {
	mux.HandleFunc("/product-links/listing-snapshots/imports", h.handleListingSnapshotImports)
	mux.HandleFunc("/product-links/link-candidates/generations", h.handleLinkCandidateGenerations)
	mux.HandleFunc("/product-links/link-candidates", h.handleLinkCandidates)
	mux.HandleFunc("/product-links/link-workflows", h.handleLinkWorkflows)
	mux.HandleFunc("/product-links/link-resolutions/approve-candidate", h.handleApproveCandidate)
	mux.HandleFunc("/product-links/link-resolutions/reject-listing", h.handleRejectListing)
	mux.HandleFunc("/product-links/link-resolutions/manual-resolve", h.handleManualResolve)
	mux.HandleFunc("/product-links/link-resolutions/batch-preview", h.handleBatchPreview)
	mux.HandleFunc("/product-links/link-resolutions/batch", h.handleBatchApply)
}

type apiError struct {
	Code    string         `json:"code"`
	Message string         `json:"message"`
	Details map[string]any `json:"details"`
}

type apiErrorResponse struct {
	Error apiError `json:"error"`
}

func writeProductLinksError(w http.ResponseWriter, status int, code, message string) {
	httpx.WriteJSON(w, status, apiErrorResponse{
		Error: apiError{
			Code:    code,
			Message: message,
			Details: map[string]any{},
		},
	})
}

func mapProductLinksError(err error) (int, string, string) {
	if err == nil {
		return http.StatusInternalServerError, "PRODUCT_LINKS_INTERNAL_ERROR", "internal error"
	}
	msg := err.Error()
	if msg == "PRODUCT_LINKS_INTERNAL_READ_UNAVAILABLE" {
		return http.StatusServiceUnavailable, msg, msg
	}
	if strings.HasPrefix(msg, "PRODUCT_LINKS_") || strings.HasPrefix(msg, "INTEGRATIONS_") {
		if strings.HasSuffix(msg, "_NOT_FOUND") {
			return http.StatusNotFound, msg, msg
		}
		return http.StatusBadRequest, msg, msg
	}
	return http.StatusInternalServerError, "PRODUCT_LINKS_INTERNAL_ERROR", "internal error"
}

func (h Handler) handleListingSnapshotImports(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		slog.Info("product_links.listing_snapshot_imports", "action", "reject_method", "result", "405", "duration_ms", time.Since(start).Milliseconds())
		writeProductLinksError(w, http.StatusMethodNotAllowed, "PRODUCT_LINKS_METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	var req struct {
		InstallationID string `json:"installation_id"`
		Limit          int    `json:"limit"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Info("product_links.listing_snapshot_imports", "action", "decode", "result", "400", "duration_ms", time.Since(start).Milliseconds())
		writeProductLinksError(w, http.StatusBadRequest, "PRODUCT_LINKS_INVALID_REQUEST", "malformed request body")
		return
	}

	limit := req.Limit
	if limit <= 0 {
		limit = 20
	}

	result, err := h.importer.ImportListingSnapshots(r.Context(), application.ImportListingSnapshotsInput{
		InstallationID: req.InstallationID,
		Limit:          limit,
	})
	if err != nil {
		status, code, message := mapProductLinksError(err)
		slog.Error("product_links.listing_snapshot_imports", "action", "import", "result", status, "error", err.Error(), "duration_ms", time.Since(start).Milliseconds())
		writeProductLinksError(w, status, code, message)
		return
	}

	slog.Info("product_links.listing_snapshot_imports", "action", "import", "result", "200", "installation_id", result.InstallationID, "count", result.ImportedCount, "duration_ms", time.Since(start).Milliseconds())
	httpx.WriteJSON(w, http.StatusOK, result)
}

func (h Handler) handleLinkCandidateGenerations(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		slog.Info("product_links.link_candidate_generations", "action", "reject_method", "result", "405", "duration_ms", time.Since(start).Milliseconds())
		writeProductLinksError(w, http.StatusMethodNotAllowed, "PRODUCT_LINKS_METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	var req struct {
		InstallationID string `json:"installation_id"`
		Limit          int    `json:"limit"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Info("product_links.link_candidate_generations", "action", "decode", "result", "400", "duration_ms", time.Since(start).Milliseconds())
		writeProductLinksError(w, http.StatusBadRequest, "PRODUCT_LINKS_INVALID_REQUEST", "malformed request body")
		return
	}

	result, err := h.candidateMaker.GenerateLinkCandidates(r.Context(), application.GenerateLinkCandidatesInput{
		InstallationID: req.InstallationID,
		Limit:          req.Limit,
	})
	if err != nil {
		status, code, message := mapProductLinksError(err)
		slog.Error("product_links.link_candidate_generations", "action", "generate", "result", status, "error", err.Error(), "duration_ms", time.Since(start).Milliseconds())
		writeProductLinksError(w, status, code, message)
		return
	}

	slog.Info("product_links.link_candidate_generations", "action", "generate", "result", "200", "installation_id", result.InstallationID, "count", result.GeneratedCount, "duration_ms", time.Since(start).Milliseconds())
	httpx.WriteJSON(w, http.StatusOK, result)
}

func (h Handler) handleLinkCandidates(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		slog.Info("product_links.link_candidates", "action", "reject_method", "result", "405", "duration_ms", time.Since(start).Milliseconds())
		writeProductLinksError(w, http.StatusMethodNotAllowed, "PRODUCT_LINKS_METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	items, err := h.candidateReader.ListLinkCandidates(r.Context(), application.ListLinkCandidatesInput{
		InstallationID: r.URL.Query().Get("installation_id"),
		Limit:          parsePositiveInt(r.URL.Query().Get("limit"), 20),
	})
	if err != nil {
		status, code, message := mapProductLinksError(err)
		slog.Error("product_links.link_candidates", "action", "list", "result", status, "error", err.Error(), "duration_ms", time.Since(start).Milliseconds())
		writeProductLinksError(w, status, code, message)
		return
	}

	slog.Info("product_links.link_candidates", "action", "list", "result", "200", "count", len(items), "duration_ms", time.Since(start).Milliseconds())
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"items": items})
}

func (h Handler) handleLinkWorkflows(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		slog.Info("product_links.link_workflows", "action", "reject_method", "result", "405", "duration_ms", time.Since(start).Milliseconds())
		writeProductLinksError(w, http.StatusMethodNotAllowed, "PRODUCT_LINKS_METHOD_NOT_ALLOWED", "method not allowed")
		return
	}
	items, err := h.workflowReader.ListLinkWorkflows(r.Context(), application.ListLinkWorkflowsInput{
		InstallationID: r.URL.Query().Get("installation_id"),
		Limit:          parsePositiveInt(r.URL.Query().Get("limit"), 20),
	})
	if err != nil {
		status, code, message := mapProductLinksError(err)
		slog.Error("product_links.link_workflows", "action", "list", "result", status, "error", err.Error(), "duration_ms", time.Since(start).Milliseconds())
		writeProductLinksError(w, status, code, message)
		return
	}
	slog.Info("product_links.link_workflows", "action", "list", "result", "200", "count", len(items), "duration_ms", time.Since(start).Milliseconds())
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"items": items})
}

func (h Handler) handleApproveCandidate(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		slog.Info("product_links.link_resolutions", "action", "reject_method", "result", "405", "duration_ms", time.Since(start).Milliseconds())
		writeProductLinksError(w, http.StatusMethodNotAllowed, "PRODUCT_LINKS_METHOD_NOT_ALLOWED", "method not allowed")
		return
	}
	var req struct {
		CandidateID string               `json:"candidate_id"`
		Reason      string               `json:"reason"`
		Actor       domain.ActorMetadata `json:"actor"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeProductLinksError(w, http.StatusBadRequest, "PRODUCT_LINKS_INVALID_REQUEST", "malformed request body")
		return
	}
	result, err := h.workflowResolver.ApproveCandidate(r.Context(), application.ApproveCandidateInput{
		CandidateID: req.CandidateID,
		Actor:       req.Actor,
		Reason:      req.Reason,
	})
	if err != nil {
		status, code, message := mapProductLinksError(err)
		slog.Error("product_links.link_resolutions", "action", "approve_candidate", "result", status, "error", err.Error(), "duration_ms", time.Since(start).Milliseconds())
		writeProductLinksError(w, status, code, message)
		return
	}
	slog.Info("product_links.link_resolutions", "action", "approve_candidate", "result", "200", "candidate_id", req.CandidateID, "duration_ms", time.Since(start).Milliseconds())
	httpx.WriteJSON(w, http.StatusOK, result)
}

func (h Handler) handleRejectListing(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		writeProductLinksError(w, http.StatusMethodNotAllowed, "PRODUCT_LINKS_METHOD_NOT_ALLOWED", "method not allowed")
		return
	}
	var req struct {
		InstallationID      string               `json:"installation_id"`
		ProviderCode        string               `json:"provider_code"`
		ProviderItemID      string               `json:"provider_item_id"`
		ProviderVariationID string               `json:"provider_variation_id"`
		Reason              string               `json:"reason"`
		Actor               domain.ActorMetadata `json:"actor"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeProductLinksError(w, http.StatusBadRequest, "PRODUCT_LINKS_INVALID_REQUEST", "malformed request body")
		return
	}
	result, err := h.workflowResolver.RejectListing(r.Context(), application.RejectListingInput{
		InstallationID:      req.InstallationID,
		ProviderCode:        req.ProviderCode,
		ProviderItemID:      req.ProviderItemID,
		ProviderVariationID: req.ProviderVariationID,
		Actor:               req.Actor,
		Reason:              req.Reason,
	})
	if err != nil {
		status, code, message := mapProductLinksError(err)
		slog.Error("product_links.link_resolutions", "action", "reject_listing", "result", status, "error", err.Error(), "duration_ms", time.Since(start).Milliseconds())
		writeProductLinksError(w, status, code, message)
		return
	}
	slog.Info("product_links.link_resolutions", "action", "reject_listing", "result", "200", "installation_id", req.InstallationID, "duration_ms", time.Since(start).Milliseconds())
	httpx.WriteJSON(w, http.StatusOK, result)
}

func (h Handler) handleManualResolve(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		writeProductLinksError(w, http.StatusMethodNotAllowed, "PRODUCT_LINKS_METHOD_NOT_ALLOWED", "method not allowed")
		return
	}
	var req struct {
		InstallationID        string               `json:"installation_id"`
		ProviderCode          string               `json:"provider_code"`
		ProviderItemID        string               `json:"provider_item_id"`
		ProviderVariationID   string               `json:"provider_variation_id"`
		InternalProductID     int                  `json:"internal_product_id"`
		InternalProductName   string               `json:"internal_product_name"`
		InternalReferenceCode string               `json:"internal_reference_code"`
		Reason                string               `json:"reason"`
		Actor                 domain.ActorMetadata `json:"actor"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeProductLinksError(w, http.StatusBadRequest, "PRODUCT_LINKS_INVALID_REQUEST", "malformed request body")
		return
	}
	if err := domain.ValidateInternalProductID(req.InternalProductID); err != nil {
		writeProductLinksError(w, http.StatusBadRequest, "invalid_identity", "internal_product_id must be a positive integer")
		return
	}
	result, err := h.workflowResolver.ManualResolve(r.Context(), application.ManualResolveInput{
		InstallationID:        req.InstallationID,
		ProviderCode:          req.ProviderCode,
		ProviderItemID:        req.ProviderItemID,
		ProviderVariationID:   req.ProviderVariationID,
		InternalProductID:     req.InternalProductID,
		InternalProductName:   req.InternalProductName,
		InternalReferenceCode: req.InternalReferenceCode,
		Actor:                 req.Actor,
		Reason:                req.Reason,
	})
	if err != nil {
		status, code, message := mapProductLinksError(err)
		slog.Error("product_links.link_resolutions", "action", "manual_resolve", "result", status, "error", err.Error(), "duration_ms", time.Since(start).Milliseconds())
		writeProductLinksError(w, status, code, message)
		return
	}
	slog.Info("product_links.link_resolutions", "action", "manual_resolve", "result", "200", "installation_id", req.InstallationID, "duration_ms", time.Since(start).Milliseconds())
	httpx.WriteJSON(w, http.StatusOK, result)
}

func (h Handler) handleBatchPreview(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		slog.Info("product_links.link_resolutions", "action", "batch_preview_method", "result", "405", "duration_ms", time.Since(start).Milliseconds())
		writeProductLinksError(w, http.StatusMethodNotAllowed, "PRODUCT_LINKS_METHOD_NOT_ALLOWED", "method not allowed")
		return
	}
	if h.batchResolver == nil {
		slog.Error("product_links.link_resolutions", "action", "batch_preview", "result", "500", "error", "batch resolver not configured", "duration_ms", time.Since(start).Milliseconds())
		writeProductLinksError(w, http.StatusInternalServerError, "PRODUCT_LINKS_INTERNAL_ERROR", "internal error")
		return
	}
	var req struct {
		Approvals []struct {
			CandidateID string `json:"candidate_id"`
		} `json:"approvals"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Info("product_links.link_resolutions", "action", "batch_preview_decode", "result", "400", "duration_ms", time.Since(start).Milliseconds())
		writeProductLinksError(w, http.StatusBadRequest, "PRODUCT_LINKS_INVALID_REQUEST", "malformed request body")
		return
	}
	approvals := make([]application.BatchApprovalInput, 0, len(req.Approvals))
	for _, item := range req.Approvals {
		approvals = append(approvals, application.BatchApprovalInput{CandidateID: item.CandidateID})
	}
	result, err := h.batchResolver.PreviewBatch(r.Context(), application.PreviewBatchInput{Approvals: approvals})
	if err != nil {
		if errors.Is(err, application.ErrBatchApprovalsRequired) {
			slog.Info("product_links.link_resolutions", "action", "batch_preview", "result", "422", "duration_ms", time.Since(start).Milliseconds())
			writeProductLinksError(w, http.StatusUnprocessableEntity, "PRODUCT_LINKS_BATCH_APPROVALS_REQUIRED", "at least one approval is required")
			return
		}
		status, code, message := mapProductLinksError(err)
		slog.Error("product_links.link_resolutions", "action", "batch_preview", "result", status, "error", err.Error(), "duration_ms", time.Since(start).Milliseconds())
		writeProductLinksError(w, status, code, message)
		return
	}
	slog.Info("product_links.link_resolutions", "action", "batch_preview", "result", "200", "count", len(result.Items), "duration_ms", time.Since(start).Milliseconds())
	httpx.WriteJSON(w, http.StatusOK, result)
}

func (h Handler) handleBatchApply(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		slog.Info("product_links.link_resolutions", "action", "reject_method", "result", "405", "duration_ms", time.Since(start).Milliseconds())
		writeProductLinksError(w, http.StatusMethodNotAllowed, "PRODUCT_LINKS_METHOD_NOT_ALLOWED", "method not allowed")
		return
	}
	if h.batchResolver == nil {
		slog.Error("product_links.link_resolutions", "action", "batch_apply", "result", "500", "error", "batch resolver not configured", "duration_ms", time.Since(start).Milliseconds())
		writeProductLinksError(w, http.StatusInternalServerError, "PRODUCT_LINKS_INTERNAL_ERROR", "internal error")
		return
	}
	var req struct {
		Approvals []struct {
			CandidateID string `json:"candidate_id"`
		} `json:"approvals"`
		Actor domain.ActorMetadata `json:"actor"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Info("product_links.link_resolutions", "action", "batch_apply_decode", "result", "400", "duration_ms", time.Since(start).Milliseconds())
		writeProductLinksError(w, http.StatusBadRequest, "PRODUCT_LINKS_INVALID_REQUEST", "malformed request body")
		return
	}
	approvals := make([]application.ApplyApprovalInput, 0, len(req.Approvals))
	for _, item := range req.Approvals {
		approvals = append(approvals, application.ApplyApprovalInput{CandidateID: item.CandidateID})
	}
	result, err := h.batchResolver.ApplyBatch(r.Context(), application.ApplyBatchInput{Approvals: approvals, Actor: req.Actor})
	if err != nil {
		if errors.Is(err, application.ErrBatchApprovalsRequired) {
			slog.Info("product_links.link_resolutions", "action", "batch_apply", "result", "422", "duration_ms", time.Since(start).Milliseconds())
			writeProductLinksError(w, http.StatusUnprocessableEntity, "PRODUCT_LINKS_BATCH_APPROVALS_REQUIRED", "at least one approval is required")
			return
		}
		status, code, message := mapProductLinksError(err)
		slog.Error("product_links.link_resolutions", "action", "batch_apply", "result", status, "error", err.Error(), "duration_ms", time.Since(start).Milliseconds())
		writeProductLinksError(w, status, code, message)
		return
	}
	slog.Info("product_links.link_resolutions", "action", "batch_apply", "result", "200", "batch_id", result.BatchID, "applied", len(result.Applied), "failed", len(result.Failed), "duration_ms", time.Since(start).Milliseconds())
	httpx.WriteJSON(w, http.StatusOK, result)
}

func parsePositiveInt(raw string, fallback int) int {
	value := strings.TrimSpace(raw)
	if value == "" {
		return fallback
	}
	var parsed int
	if _, err := fmt.Sscanf(value, "%d", &parsed); err != nil || parsed <= 0 {
		return fallback
	}
	return parsed
}
