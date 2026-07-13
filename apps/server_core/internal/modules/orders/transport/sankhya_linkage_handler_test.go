package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/modules/orders/application"
	"marketplace-central/apps/server_core/internal/modules/orders/domain"
)

type stubAssistedSankhyaApplication struct {
	getInput     application.GetAssistedSankhyaLinkageInput
	listInput    application.ListAssistedSankhyaCandidatesInput
	confirmInput application.ConfirmAssistedSankhyaLinkageInput
	current      domain.SankhyaLinkage
	candidates   domain.AssistedSankhyaCandidateResult
	confirmation domain.AssistedSankhyaConfirmationResult
	err          error
	confirmCalls int
}

func (s *stubAssistedSankhyaApplication) GetCurrent(_ context.Context, input application.GetAssistedSankhyaLinkageInput) (domain.SankhyaLinkage, error) {
	s.getInput = input
	return s.current, s.err
}

func (s *stubAssistedSankhyaApplication) ListCandidates(_ context.Context, input application.ListAssistedSankhyaCandidatesInput) (domain.AssistedSankhyaCandidateResult, error) {
	s.listInput = input
	return s.candidates, s.err
}

func (s *stubAssistedSankhyaApplication) Confirm(_ context.Context, input application.ConfirmAssistedSankhyaLinkageInput) (domain.AssistedSankhyaConfirmationResult, error) {
	s.confirmCalls++
	s.confirmInput = input
	return s.confirmation, s.err
}

func TestSankhyaLinkageHandlersUseDefaultTenantAndExactInstallationOrderScope(t *testing.T) {
	service := assistedSankhyaTransportFixture()
	mux := http.NewServeMux()
	NewSankhyaLinkageHandler(service, " tenant-default ").Register(mux)

	for _, path := range []string{
		"/orders/order-1/sankhya-linkage?installation_id=install-1",
		"/orders/order-1/sankhya-linkage/candidates?installation_id=install-1",
	} {
		recorder := httptest.NewRecorder()
		mux.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, path, nil))
		if recorder.Code != http.StatusOK {
			t.Fatalf("GET %s status = %d body=%s", path, recorder.Code, recorder.Body.String())
		}
	}
	if service.getInput.TenantID != "tenant-default" || service.getInput.InstallationID != "install-1" || service.getInput.ProviderOrderID != "order-1" {
		t.Fatalf("current input = %#v", service.getInput)
	}
	if service.listInput.TenantID != "tenant-default" || service.listInput.InstallationID != "install-1" || service.listInput.ProviderOrderID != "order-1" {
		t.Fatalf("candidate input = %#v", service.listInput)
	}
}

func TestSankhyaLinkageConfirmAcceptsOnlyOperatorSelectionFacts(t *testing.T) {
	service := assistedSankhyaTransportFixture()
	mux := http.NewServeMux()
	NewSankhyaLinkageHandler(service, "tenant-default").Register(mux)
	body := `{
  "selected_document_id":31301,
  "selections":[{"mpc_line_id":"mpl_aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","internal_document_id":31301,"internal_line_number":1}],
  "actor_id":"operator-7","reason":"exact selection","source_at":"2026-07-13T14:00:00-03:00","idempotency_key":"idem-1"
}`
	recorder := httptest.NewRecorder()
	mux.ServeHTTP(recorder, httptest.NewRequest(http.MethodPost, "/orders/order-1/sankhya-linkage/confirm?installation_id=install-1", bytes.NewBufferString(body)))
	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", recorder.Code, recorder.Body.String())
	}
	if service.confirmInput.TenantID != "tenant-default" || service.confirmInput.ProviderOrderID != "order-1" || service.confirmInput.SelectedDocumentID != 31301 || service.confirmInput.SourceAt.IsZero() {
		t.Fatalf("confirm input = %#v", service.confirmInput)
	}
	var payload map[string]any
	if err := json.NewDecoder(recorder.Body).Decode(&payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	audit := payload["audit"].(map[string]any)
	if audit["configuration_revision"] != "cfg-runtime-1" || audit["evidence_reference"] != "sankhya-linkage:v1:safe" || audit["actor_type"] != domain.AssistedSankhyaActorOperatorSuppliedUnverified {
		t.Fatalf("audit = %#v", audit)
	}
}

func TestSankhyaLinkageConfirmRejectsCallerControlledServerFactsAndExtraJSON(t *testing.T) {
	for _, field := range []string{"tenant_id", "event_id", "configuration_revision", "evidence_reference", "actor_type", "external_order_key"} {
		t.Run(field, func(t *testing.T) {
			service := assistedSankhyaTransportFixture()
			mux := http.NewServeMux()
			NewSankhyaLinkageHandler(service, "tenant-default").Register(mux)
			body := `{"selected_document_id":31301,"selections":[],"actor_id":"operator","reason":"reason","source_at":"2026-07-13T14:00:00Z","idempotency_key":"idem","` + field + `":"caller"}`
			recorder := httptest.NewRecorder()
			mux.ServeHTTP(recorder, httptest.NewRequest(http.MethodPost, "/orders/order-1/sankhya-linkage/confirm?installation_id=install-1", bytes.NewBufferString(body)))
			if recorder.Code != http.StatusBadRequest || service.confirmCalls != 0 {
				t.Fatalf("field %s status/calls = %d/%d body=%s", field, recorder.Code, service.confirmCalls, recorder.Body.String())
			}
		})
	}
}

func TestSankhyaLinkageTransportRejectsMalformedScopeTimeAndMultipleJSONValues(t *testing.T) {
	tests := []struct {
		method string
		path   string
		body   string
	}{
		{http.MethodGet, "/orders/order-1/sankhya-linkage", ""},
		{http.MethodGet, "/orders/order-1/sankhya-linkage?installation_id=install-1&tenant_id=caller", ""},
		{http.MethodPost, "/orders/order-1/sankhya-linkage/confirm?installation_id=install-1", `{"source_at":"not-a-time"}`},
		{http.MethodPost, "/orders/order-1/sankhya-linkage/confirm?installation_id=install-1", `{}` + `{}`},
	}
	for _, test := range tests {
		service := assistedSankhyaTransportFixture()
		mux := http.NewServeMux()
		NewSankhyaLinkageHandler(service, "tenant-default").Register(mux)
		recorder := httptest.NewRecorder()
		mux.ServeHTTP(recorder, httptest.NewRequest(test.method, test.path, bytes.NewBufferString(test.body)))
		if recorder.Code != http.StatusBadRequest {
			t.Fatalf("%s %s status=%d body=%s", test.method, test.path, recorder.Code, recorder.Body.String())
		}
	}
}

func TestSankhyaLinkageTransportMapsSafeStableErrors(t *testing.T) {
	tests := []struct {
		err    error
		status int
		code   string
	}{
		{domain.ErrAssistedSankhyaOrderNotFound, http.StatusNotFound, "ORDERS_SANKHYA_LINKAGE_NOT_FOUND"},
		{domain.ErrSankhyaLinkageConflict, http.StatusConflict, "ORDERS_SANKHYA_LINKAGE_CONFLICT"},
		{&domain.AssistedSankhyaReadError{Kind: domain.AssistedSankhyaReadConfigurationInvalid}, http.StatusServiceUnavailable, "ORDERS_SANKHYA_LINKAGE_UNAVAILABLE"},
		{errors.New("driver secret detail"), http.StatusInternalServerError, "ORDERS_SANKHYA_LINKAGE_INTERNAL_ERROR"},
	}
	for _, test := range tests {
		service := assistedSankhyaTransportFixture()
		service.err = test.err
		mux := http.NewServeMux()
		NewSankhyaLinkageHandler(service, "tenant-default").Register(mux)
		recorder := httptest.NewRecorder()
		mux.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/orders/order-1/sankhya-linkage?installation_id=install-1", nil))
		if recorder.Code != test.status || !bytes.Contains(recorder.Body.Bytes(), []byte(test.code)) || bytes.Contains(recorder.Body.Bytes(), []byte("driver secret")) {
			t.Fatalf("error %v status/body=%d/%s", test.err, recorder.Code, recorder.Body.String())
		}
	}
}

func TestSankhyaLinkageRoutesStayRegisteredWhenRuntimeUnavailable(t *testing.T) {
	mux := http.NewServeMux()
	NewSankhyaLinkageHandler(nil, "tenant-default").Register(mux)
	for _, path := range []string{
		"/orders/order-1/sankhya-linkage?installation_id=install-1",
		"/orders/order-1/sankhya-linkage/candidates?installation_id=install-1",
	} {
		recorder := httptest.NewRecorder()
		mux.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, path, nil))
		if recorder.Code != http.StatusServiceUnavailable {
			t.Fatalf("GET %s status=%d body=%s", path, recorder.Code, recorder.Body.String())
		}
	}
}

func assistedSankhyaTransportFixture() *stubAssistedSankhyaApplication {
	linkage := domain.SankhyaLinkage{
		Scope:            domain.LinkageScope{TenantID: "tenant-default", InstallationID: "install-1", ProviderOrderID: "order-1"},
		ExternalOrderKey: "ml:v1:install-1:order-1", Header: domain.InternalDocumentIdentity{DocumentID: 31301},
		Lines: []domain.SankhyaLineMapping{{MPCLineID: "mpl_aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", Origin: domain.InternalDocumentLineIdentity{DocumentID: 31301, LineNumber: 1}}},
		Audit: domain.LinkageAudit{
			EventID: "server-event-1", EventType: domain.LinkageEventConfirmed,
			ActorType: domain.AssistedSankhyaActorOperatorSuppliedUnverified, ActorID: "operator-7", Reason: "exact selection",
			IdempotencyKey: "idem-1", SourceAt: time.Date(2026, 7, 13, 17, 0, 0, 0, time.UTC), RecordedAt: time.Date(2026, 7, 13, 17, 1, 0, 0, time.UTC),
			ConfigurationRevision: "cfg-runtime-1", EvidenceState: domain.LinkageEvidenceExact, EvidenceReference: "sankhya-linkage:v1:safe",
		},
	}
	quantity := 2.0
	candidate := domain.AssistedSankhyaCandidate{
		Header: domain.InternalDocumentIdentity{DocumentID: 31301}, DocumentNumber: nil, OperationCode: 313, Amount: nil,
		Lines: []domain.AssistedSankhyaCandidateLine{{Identity: domain.InternalDocumentLineIdentity{DocumentID: 31301, LineNumber: 1}, Quantity: &quantity}},
	}
	lineage := domain.AssistedSankhyaLineage{MPCLineID: linkage.Lines[0].MPCLineID, Origin: linkage.Lines[0].Origin, State: domain.AssistedSankhyaLineageNone}
	return &stubAssistedSankhyaApplication{
		current:      linkage,
		candidates:   domain.AssistedSankhyaCandidateResult{Scope: linkage.Scope, ExternalOrderKey: linkage.ExternalOrderKey, Candidates: []domain.AssistedSankhyaCandidate{candidate}},
		confirmation: domain.AssistedSankhyaConfirmationResult{Linkage: linkage, Lines: []domain.AssistedSankhyaLineage{lineage}},
	}
}
