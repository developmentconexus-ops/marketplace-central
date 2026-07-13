package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/modules/product_links/application"
	"marketplace-central/apps/server_core/internal/modules/product_links/domain"
)

type stubImporter struct {
	input  application.ImportListingSnapshotsInput
	result domain.ListingSnapshotImportResult
}

func (s *stubImporter) ImportListingSnapshots(_ context.Context, input application.ImportListingSnapshotsInput) (domain.ListingSnapshotImportResult, error) {
	s.input = input
	if s.result.InstallationID == "" {
		s.result = domain.ListingSnapshotImportResult{
			InstallationID: input.InstallationID,
			ImportedCount:  1,
			Items: []domain.ListingSnapshot{{
				InstallationID:    input.InstallationID,
				ProviderCode:      "mercado_livre",
				ProviderItemID:    "MLB123",
				ProviderStatus:    "active",
				Title:             "Produto teste",
				FetchedAt:         time.Date(2026, 7, 8, 19, 10, 0, 0, time.UTC),
				AvailableQuantity: ptrInt(5),
			}},
		}
	}
	return s.result, nil
}

type stubCandidateEngine struct {
	generateInput    application.GenerateLinkCandidatesInput
	listInput        application.ListLinkCandidatesInput
	workflowInput    application.ListLinkWorkflowsInput
	result           domain.LinkCandidateGenerationResult
	items            []domain.LinkCandidate
	workflowItems    []domain.ProductLinkWorkflowItem
	resolutionResult domain.ProductLinkResolutionResult
}

func (s *stubCandidateEngine) GenerateLinkCandidates(_ context.Context, input application.GenerateLinkCandidatesInput) (domain.LinkCandidateGenerationResult, error) {
	s.generateInput = input
	if s.result.InstallationID == "" {
		s.result = domain.LinkCandidateGenerationResult{
			InstallationID: input.InstallationID,
			GeneratedCount: 1,
			Items: []domain.LinkCandidate{{
				CandidateID:         "cand-1",
				InstallationID:      input.InstallationID,
				ProviderCode:        "mercado_livre",
				ProviderItemID:      "MLB123",
				InternalProductName: "Produto Interno",
				State:               domain.LinkCandidateStateExactEAN,
				MatchInput:          domain.LinkCandidateMatchInputEAN,
				MatchValue:          "789",
				CreatedAt:           time.Date(2026, 7, 8, 19, 10, 0, 0, time.UTC),
				UpdatedAt:           time.Date(2026, 7, 8, 19, 10, 0, 0, time.UTC),
			}},
		}
	}
	return s.result, nil
}

func (s *stubCandidateEngine) ListLinkCandidates(_ context.Context, input application.ListLinkCandidatesInput) ([]domain.LinkCandidate, error) {
	s.listInput = input
	if len(s.items) == 0 {
		s.items = []domain.LinkCandidate{{
			CandidateID:         "cand-1",
			InstallationID:      input.InstallationID,
			ProviderCode:        "mercado_livre",
			ProviderItemID:      "MLB123",
			InternalProductName: "Produto Interno",
			State:               domain.LinkCandidateStateExactEAN,
			MatchInput:          domain.LinkCandidateMatchInputEAN,
			MatchValue:          "789",
			CreatedAt:           time.Date(2026, 7, 8, 19, 10, 0, 0, time.UTC),
			UpdatedAt:           time.Date(2026, 7, 8, 19, 10, 0, 0, time.UTC),
		}}
	}
	return s.items, nil
}

func (s *stubCandidateEngine) ListLinkWorkflows(_ context.Context, input application.ListLinkWorkflowsInput) ([]domain.ProductLinkWorkflowItem, error) {
	s.workflowInput = input
	if len(s.workflowItems) == 0 {
		s.workflowItems = []domain.ProductLinkWorkflowItem{{
			Identity: domain.ListingIdentity{
				InstallationID: input.InstallationID,
				ProviderItemID: "MLB123",
			},
			Candidates: []domain.LinkCandidate{{
				CandidateID: "cand-1",
				State:       domain.LinkCandidateStateConflict,
			}},
		}}
	}
	return s.workflowItems, nil
}

func (s *stubCandidateEngine) ApproveCandidate(_ context.Context, input application.ApproveCandidateInput) (domain.ProductLinkResolutionResult, error) {
	if s.resolutionResult.Link.InstallationID == "" {
		productID := 10
		s.resolutionResult = domain.ProductLinkResolutionResult{
			Link: domain.ProductLink{
				InstallationID:    "inst-1",
				ProviderCode:      "mercado_livre",
				ProviderItemID:    "MLB123",
				State:             domain.ProductLinkStateResolved,
				SourceCandidateID: input.CandidateID,
				InternalProductID: &productID,
				CreatedAt:         time.Date(2026, 7, 8, 19, 10, 0, 0, time.UTC),
				UpdatedAt:         time.Date(2026, 7, 8, 19, 10, 0, 0, time.UTC),
			},
			Audit: domain.ProductLinkAuditEntry{
				AuditID:       "audit-1",
				Action:        domain.ProductLinkActionApproveCandidate,
				PreviousState: domain.ProductLinkStateNone,
				NextState:     domain.ProductLinkStateResolved,
				CreatedAt:     time.Date(2026, 7, 8, 19, 10, 0, 0, time.UTC),
			},
		}
	}
	return s.resolutionResult, nil
}

func (s *stubCandidateEngine) RejectListing(_ context.Context, _ application.RejectListingInput) (domain.ProductLinkResolutionResult, error) {
	return domain.ProductLinkResolutionResult{
		Link:  domain.ProductLink{InstallationID: "inst-1", ProviderCode: "mercado_livre", ProviderItemID: "MLB123", State: domain.ProductLinkStateRejected},
		Audit: domain.ProductLinkAuditEntry{AuditID: "audit-2", Action: domain.ProductLinkActionRejectListing, PreviousState: domain.ProductLinkStateResolved, NextState: domain.ProductLinkStateRejected},
	}, nil
}

func (s *stubCandidateEngine) ManualResolve(_ context.Context, input application.ManualResolveInput) (domain.ProductLinkResolutionResult, error) {
	productID := input.InternalProductID
	return domain.ProductLinkResolutionResult{
		Link:  domain.ProductLink{InstallationID: input.InstallationID, ProviderCode: input.ProviderCode, ProviderItemID: input.ProviderItemID, ProviderVariationID: input.ProviderVariationID, State: domain.ProductLinkStateResolved, InternalProductID: &productID},
		Audit: domain.ProductLinkAuditEntry{AuditID: "audit-3", Action: domain.ProductLinkActionManualResolve, PreviousState: domain.ProductLinkStateNone, NextState: domain.ProductLinkStateResolved},
	}, nil
}

func TestHandleListingSnapshotImportsDelegatesToImporter(t *testing.T) {
	t.Parallel()

	importer := &stubImporter{}
	handler := NewHandler(importer, &stubCandidateEngine{}, &stubCandidateEngine{}, &stubCandidateEngine{}, &stubCandidateEngine{})
	req := httptest.NewRequest(http.MethodPost, "/product-links/listing-snapshots/imports", bytes.NewReader([]byte(`{"installation_id":"inst-1","limit":3}`)))
	rr := httptest.NewRecorder()

	handler.handleListingSnapshotImports(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rr.Code, rr.Body.String())
	}
	if importer.input.InstallationID != "inst-1" || importer.input.Limit != 3 {
		t.Fatalf("input=%#v", importer.input)
	}
}

func TestHandleListingSnapshotImportsUsesSnakeCaseJSON(t *testing.T) {
	t.Parallel()

	handler := NewHandler(&stubImporter{}, &stubCandidateEngine{}, &stubCandidateEngine{}, &stubCandidateEngine{}, &stubCandidateEngine{})
	req := httptest.NewRequest(http.MethodPost, "/product-links/listing-snapshots/imports", bytes.NewReader([]byte(`{"installation_id":"inst-1"}`)))
	rr := httptest.NewRecorder()

	handler.handleListingSnapshotImports(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rr.Code, rr.Body.String())
	}
	var payload map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if _, ok := payload["imported_count"]; !ok {
		t.Fatalf("payload=%v", payload)
	}
	items, ok := payload["items"].([]any)
	if !ok || len(items) != 1 {
		t.Fatalf("payload=%v", payload)
	}
	item, ok := items[0].(map[string]any)
	if !ok {
		t.Fatalf("payload=%v", payload)
	}
	if _, ok := item["provider_item_id"]; !ok {
		t.Fatalf("payload=%v", item)
	}
	if _, ok := item["ProviderItemID"]; ok {
		t.Fatalf("payload=%v", item)
	}
}

func TestHandleListingSnapshotImportsRejectsWrongMethod(t *testing.T) {
	t.Parallel()

	handler := NewHandler(&stubImporter{}, &stubCandidateEngine{}, &stubCandidateEngine{}, &stubCandidateEngine{}, &stubCandidateEngine{})
	req := httptest.NewRequest(http.MethodGet, "/product-links/listing-snapshots/imports", nil)
	rr := httptest.NewRecorder()

	handler.handleListingSnapshotImports(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status=%d body=%s", rr.Code, rr.Body.String())
	}
}

func TestHandleLinkCandidateGenerationsDelegatesToService(t *testing.T) {
	t.Parallel()

	engine := &stubCandidateEngine{}
	handler := NewHandler(&stubImporter{}, engine, engine, engine, engine)
	req := httptest.NewRequest(http.MethodPost, "/product-links/link-candidates/generations", bytes.NewReader([]byte(`{"installation_id":"inst-1","limit":7}`)))
	rr := httptest.NewRecorder()

	handler.handleLinkCandidateGenerations(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rr.Code, rr.Body.String())
	}
	if engine.generateInput.InstallationID != "inst-1" || engine.generateInput.Limit != 7 {
		t.Fatalf("input=%#v", engine.generateInput)
	}
}

func TestHandleLinkCandidateGenerationsUsesSnakeCaseJSON(t *testing.T) {
	t.Parallel()

	engine := &stubCandidateEngine{}
	handler := NewHandler(&stubImporter{}, engine, engine, engine, engine)
	req := httptest.NewRequest(http.MethodPost, "/product-links/link-candidates/generations", bytes.NewReader([]byte(`{"installation_id":"inst-1"}`)))
	rr := httptest.NewRecorder()

	handler.handleLinkCandidateGenerations(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rr.Code, rr.Body.String())
	}
	var payload map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if _, ok := payload["generated_count"]; !ok {
		t.Fatalf("payload=%v", payload)
	}
}

func TestHandleLinkCandidatesReturnsItems(t *testing.T) {
	t.Parallel()

	engine := &stubCandidateEngine{}
	handler := NewHandler(&stubImporter{}, engine, engine, engine, engine)
	req := httptest.NewRequest(http.MethodGet, "/product-links/link-candidates?installation_id=inst-1&limit=4", nil)
	rr := httptest.NewRecorder()

	handler.handleLinkCandidates(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rr.Code, rr.Body.String())
	}
	if engine.listInput.InstallationID != "inst-1" || engine.listInput.Limit != 4 {
		t.Fatalf("input=%#v", engine.listInput)
	}
	var payload struct {
		Items []map[string]any `json:"items"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(payload.Items) != 1 {
		t.Fatalf("items=%v", payload.Items)
	}
	if _, ok := payload.Items[0]["candidate_id"]; !ok {
		t.Fatalf("payload=%v", payload.Items[0])
	}
}

func ptrInt(value int) *int {
	return &value
}

func TestHandleLinkWorkflowsReturnsItems(t *testing.T) {
	t.Parallel()

	engine := &stubCandidateEngine{}
	handler := NewHandler(&stubImporter{}, engine, engine, engine, engine)
	req := httptest.NewRequest(http.MethodGet, "/product-links/link-workflows?installation_id=inst-1&limit=4", nil)
	rr := httptest.NewRecorder()

	handler.handleLinkWorkflows(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rr.Code, rr.Body.String())
	}
	if engine.workflowInput.InstallationID != "inst-1" || engine.workflowInput.Limit != 4 {
		t.Fatalf("input=%#v", engine.workflowInput)
	}
}

func TestHandleApproveCandidateUsesSnakeCaseJSON(t *testing.T) {
	t.Parallel()

	engine := &stubCandidateEngine{}
	handler := NewHandler(&stubImporter{}, engine, engine, engine, engine)
	req := httptest.NewRequest(http.MethodPost, "/product-links/link-resolutions/approve-candidate", bytes.NewReader([]byte(`{"candidate_id":"cand-1","reason":"ok","actor":{"actor_type":"operator","actor_id":"u1"}}`)))
	rr := httptest.NewRecorder()

	handler.handleApproveCandidate(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rr.Code, rr.Body.String())
	}
	var payload map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if _, ok := payload["link"]; !ok {
		t.Fatalf("payload=%v", payload)
	}
	if _, ok := payload["audit"]; !ok {
		t.Fatalf("payload=%v", payload)
	}
}

func TestHandleManualResolveRejectsNonPositiveInternalProductID(t *testing.T) {
	t.Parallel()

	for _, id := range []int{-1, 0} {
		engine := &stubCandidateEngine{}
		handler := NewHandler(&stubImporter{}, engine, engine, engine, engine)
		body := fmt.Sprintf(`{"installation_id":"inst-1","provider_code":"mercado_livre","provider_item_id":"MLB1","internal_product_id":%d,"actor":{"actor_type":"operator","actor_id":"u1"}}`, id)
		req := httptest.NewRequest(http.MethodPost, "/product-links/link-resolutions/manual-resolve", bytes.NewReader([]byte(body)))
		rr := httptest.NewRecorder()

		handler.handleManualResolve(rr, req)

		if rr.Code != http.StatusBadRequest || !strings.Contains(rr.Body.String(), "invalid_identity") {
			t.Fatalf("id=%d status=%d body=%s", id, rr.Code, rr.Body.String())
		}
	}
}
