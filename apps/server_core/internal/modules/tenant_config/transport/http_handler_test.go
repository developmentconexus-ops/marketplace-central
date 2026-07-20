package transport

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/modules/tenant_config"
)

type fakeStore struct {
	getCfg tenant_config.Config
	getErr error
	setErr error

	setCalledWith tenant_config.Config
	setCalled     bool

	// postSetCfg, when non-zero-valued (postSetSet true), is returned by Get
	// on the call AFTER Set succeeds — simulating a re-Get of the persisted
	// row (derived source_kind included).
	postSetCfg tenant_config.Config
	postSetSet bool
}

func (f *fakeStore) Get(context.Context, string) (tenant_config.Config, error) {
	if f.setCalled && f.postSetSet {
		return f.postSetCfg, nil
	}
	return f.getCfg, f.getErr
}

func (f *fakeStore) Set(_ context.Context, cfg tenant_config.Config) error {
	f.setCalled = true
	f.setCalledWith = cfg
	return f.setErr
}

func newMux(store Store) *http.ServeMux {
	mux := http.NewServeMux()
	NewHandler(store, "tenant-1").Register(mux)
	return mux
}

func TestHandlerGetActiveSource(t *testing.T) {
	setAt := time.Date(2026, 7, 20, 12, 0, 0, 0, time.UTC)
	store := &fakeStore{getCfg: tenant_config.Config{
		TenantID: "tenant-1",
		Source:   tenant_config.SourceXLSX,
		Kind:     "upload_snapshot",
		SetAt:    setAt,
		SetBy:    "operator-1",
	}}
	req := httptest.NewRequest(http.MethodGet, "/config/active-source", nil)
	rec := httptest.NewRecorder()
	newMux(store).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}
	var body activeSourceResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.ActiveSource != "xlsx" || body.SourceKind != "upload_snapshot" {
		t.Fatalf("body = %+v, want active_source=xlsx source_kind=upload_snapshot", body)
	}
}

func TestHandlerGetActiveSourceUnknown(t *testing.T) {
	store := &fakeStore{getErr: tenant_config.ErrUnknownActiveSource}
	req := httptest.NewRequest(http.MethodGet, "/config/active-source", nil)
	rec := httptest.NewRecorder()
	newMux(store).ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d, body=%s", rec.Code, http.StatusBadRequest, rec.Body.String())
	}
	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body["error"] != "unknown_erp_source" {
		t.Fatalf("error = %q, want unknown_erp_source", body["error"])
	}
}

func TestHandlerPutActiveSource(t *testing.T) {
	setAt := time.Date(2026, 7, 20, 13, 0, 0, 0, time.UTC)
	store := &fakeStore{
		postSetSet: true,
		postSetCfg: tenant_config.Config{
			TenantID: "tenant-1",
			Source:   tenant_config.SourceSankhya,
			Kind:     "live_read_through",
			SetAt:    setAt,
		},
	}
	req := httptest.NewRequest(http.MethodPut, "/config/active-source", strings.NewReader(`{"active_source":"sankhya"}`))
	rec := httptest.NewRecorder()
	newMux(store).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}
	if !store.setCalled || store.setCalledWith.Source != tenant_config.SourceSankhya {
		t.Fatalf("Set called with %+v, setCalled=%v", store.setCalledWith, store.setCalled)
	}
	var body activeSourceResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.ActiveSource != "sankhya" || body.SourceKind != "live_read_through" {
		t.Fatalf("body = %+v, want active_source=sankhya source_kind=live_read_through", body)
	}
}

func TestHandlerPutActiveSourceInvalid(t *testing.T) {
	store := &fakeStore{setErr: tenant_config.ErrInvalidActiveSource}
	req := httptest.NewRequest(http.MethodPut, "/config/active-source", strings.NewReader(`{"active_source":"garbage"}`))
	rec := httptest.NewRecorder()
	newMux(store).ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d, body=%s", rec.Code, http.StatusBadRequest, rec.Body.String())
	}
	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body["error"] != "invalid_active_source" {
		t.Fatalf("error = %q, want invalid_active_source", body["error"])
	}
}

func TestHandlerPutActiveSourceMalformedBody(t *testing.T) {
	store := &fakeStore{}
	req := httptest.NewRequest(http.MethodPut, "/config/active-source", strings.NewReader(`{not json`))
	rec := httptest.NewRecorder()
	newMux(store).ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d, body=%s", rec.Code, http.StatusBadRequest, rec.Body.String())
	}
	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body["error"] != "invalid_body" {
		t.Fatalf("error = %q, want invalid_body", body["error"])
	}
	if store.setCalled {
		t.Fatalf("Set must not be called on a malformed body")
	}
}

func TestHandlerPutActiveSourceIgnoresClientSourceKind(t *testing.T) {
	store := &fakeStore{
		postSetSet: true,
		postSetCfg: tenant_config.Config{
			TenantID: "tenant-1",
			Source:   tenant_config.SourceXLSX,
			Kind:     "upload_snapshot",
			SetAt:    time.Now().UTC(),
		},
	}
	// A malicious/careless client sends source_kind=live_read_through for an
	// xlsx source — the handler must derive-only and never forward it.
	req := httptest.NewRequest(http.MethodPut, "/config/active-source", strings.NewReader(
		`{"active_source":"xlsx","source_kind":"live_read_through"}`))
	rec := httptest.NewRecorder()
	newMux(store).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}
	if store.setCalledWith.Kind != "" {
		t.Fatalf("Set called with Kind=%q, want empty (server derives, client value must be ignored)", store.setCalledWith.Kind)
	}
	var body activeSourceResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.SourceKind != "upload_snapshot" {
		t.Fatalf("source_kind = %q, want upload_snapshot (server-derived, not the client-supplied live_read_through)", body.SourceKind)
	}
}
