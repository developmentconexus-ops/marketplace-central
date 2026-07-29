package transport

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/modules/tenant_config"
)

type fakeStore struct {
	getCfg                   tenant_config.Config
	getErr                   error
	setErr                   error
	setSellableAssortmentErr error

	setCalledWith tenant_config.Config
	setCalled     bool

	setSellableAssortmentCalledWith tenant_config.SellableAssortment
	setSellableAssortmentCalled     bool
	getCalls                        int

	// postSetCfg, when non-zero-valued (postSetSet true), is returned by Get
	// on the call AFTER Set succeeds — simulating a re-Get of the persisted
	// row (derived source_kind included).
	postSetCfg tenant_config.Config
	postSetSet bool
	// postSetErr covers the row vanishing between a successful write and read-back.
	postSetErr error
}

func (f *fakeStore) Get(context.Context, string) (tenant_config.Config, error) {
	f.getCalls++
	if (f.setCalled || f.setSellableAssortmentCalled) && f.postSetSet {
		return f.postSetCfg, f.postSetErr
	}
	return f.getCfg, f.getErr
}

func (f *fakeStore) Set(_ context.Context, cfg tenant_config.Config) error {
	f.setCalled = true
	f.setCalledWith = cfg
	return f.setErr
}

func (f *fakeStore) SetSellableAssortment(_ context.Context, _ string, a tenant_config.SellableAssortment) error {
	f.setSellableAssortmentCalled = true
	f.setSellableAssortmentCalledWith = a
	return f.setSellableAssortmentErr
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

func TestHandlerSellableAssortmentPutThenGetReturnsStoredValues(t *testing.T) {
	store := &fakeStore{
		postSetSet: true,
		postSetCfg: tenant_config.Config{
			TenantID: "tenant-1",
			SellableAssortment: tenant_config.SellableAssortment{
				OnlyRevenda:           false,
				OnlyEmEstoque:         true,
				OnlyEcommerceEligible: false,
			},
		},
	}
	req := httptest.NewRequest(http.MethodPut, "/config/sellable-assortment", strings.NewReader(
		`{"only_revenda":true,"only_em_estoque":false,"only_ecommerce_eligible":true}`))
	rec := httptest.NewRecorder()
	newMux(store).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}
	wantRequest := tenant_config.SellableAssortment{OnlyRevenda: true, OnlyEcommerceEligible: true}
	if !store.setSellableAssortmentCalled || store.setSellableAssortmentCalledWith != wantRequest {
		t.Fatalf("SetSellableAssortment called with %+v, want %+v", store.setSellableAssortmentCalledWith, wantRequest)
	}
	var body sellableAssortmentResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	want := sellableAssortmentResponse{OnlyRevenda: false, OnlyEmEstoque: true, OnlyEcommerceEligible: false}
	if body != want {
		t.Fatalf("body = %+v, want stored values %+v", body, store.postSetCfg.SellableAssortment)
	}
}

func TestHandlerSellableAssortmentFieldMapping(t *testing.T) {
	tests := []struct {
		name               string
		requestBody        string
		assortment         tenant_config.SellableAssortment
		wantResponseFields map[string]bool
	}{
		{
			name:        "only_revenda",
			requestBody: `{"only_revenda":true,"only_em_estoque":false,"only_ecommerce_eligible":false}`,
			assortment: tenant_config.SellableAssortment{
				OnlyRevenda: true,
			},
			wantResponseFields: map[string]bool{
				"only_revenda": true, "only_em_estoque": false, "only_ecommerce_eligible": false,
			},
		},
		{
			name:        "only_em_estoque",
			requestBody: `{"only_revenda":false,"only_em_estoque":true,"only_ecommerce_eligible":false}`,
			assortment: tenant_config.SellableAssortment{
				OnlyEmEstoque: true,
			},
			wantResponseFields: map[string]bool{
				"only_revenda": false, "only_em_estoque": true, "only_ecommerce_eligible": false,
			},
		},
		{
			name:        "only_ecommerce_eligible",
			requestBody: `{"only_revenda":false,"only_em_estoque":false,"only_ecommerce_eligible":true}`,
			assortment: tenant_config.SellableAssortment{
				OnlyEcommerceEligible: true,
			},
			wantResponseFields: map[string]bool{
				"only_revenda": false, "only_em_estoque": false, "only_ecommerce_eligible": true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name+"/PUT", func(t *testing.T) {
			store := &fakeStore{postSetSet: true, postSetCfg: tenant_config.Config{
				SellableAssortment: tt.assortment,
			}}
			req := httptest.NewRequest(http.MethodPut, "/config/sellable-assortment", strings.NewReader(tt.requestBody))
			rec := httptest.NewRecorder()
			newMux(store).ServeHTTP(rec, req)

			if rec.Code != http.StatusOK {
				t.Fatalf("status = %d, want %d, body=%s", rec.Code, http.StatusOK, rec.Body.String())
			}
			if !store.setSellableAssortmentCalled || store.setSellableAssortmentCalledWith != tt.assortment {
				t.Fatalf("SetSellableAssortment called with %+v, want %+v", store.setSellableAssortmentCalledWith, tt.assortment)
			}
			var body map[string]bool
			if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
				t.Fatalf("decode response: %v", err)
			}
			if !reflect.DeepEqual(body, tt.wantResponseFields) {
				t.Fatalf("response fields = %+v, want %+v", body, tt.wantResponseFields)
			}
		})

		t.Run(tt.name+"/GET", func(t *testing.T) {
			store := &fakeStore{getCfg: tenant_config.Config{SellableAssortment: tt.assortment}}
			req := httptest.NewRequest(http.MethodGet, "/config/sellable-assortment", nil)
			rec := httptest.NewRecorder()
			newMux(store).ServeHTTP(rec, req)

			if rec.Code != http.StatusOK {
				t.Fatalf("status = %d, want %d, body=%s", rec.Code, http.StatusOK, rec.Body.String())
			}
			var body map[string]bool
			if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
				t.Fatalf("decode response: %v", err)
			}
			if !reflect.DeepEqual(body, tt.wantResponseFields) {
				t.Fatalf("response fields = %+v, want %+v", body, tt.wantResponseFields)
			}
		})
	}
}

func TestHandlerGetSellableAssortment(t *testing.T) {
	store := &fakeStore{getCfg: tenant_config.Config{SellableAssortment: tenant_config.SellableAssortment{
		OnlyRevenda: true, OnlyEmEstoque: false, OnlyEcommerceEligible: true,
	}}}
	req := httptest.NewRequest(http.MethodGet, "/config/sellable-assortment", nil)
	rec := httptest.NewRecorder()
	newMux(store).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}
	var body sellableAssortmentResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	want := sellableAssortmentResponse{OnlyRevenda: true, OnlyEcommerceEligible: true}
	if body != want {
		t.Fatalf("body = %+v, want %+v", body, want)
	}
}

func TestHandlerGetSellableAssortmentUnknownSourceReturnsBadRequest(t *testing.T) {
	store := &fakeStore{getErr: tenant_config.ErrUnknownActiveSource}
	req := httptest.NewRequest(http.MethodGet, "/config/sellable-assortment", nil)
	rec := httptest.NewRecorder()
	newMux(store).ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d, body=%s", rec.Code, http.StatusBadRequest, rec.Body.String())
	}
	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if !reflect.DeepEqual(body, map[string]string{"error": "unknown_erp_source"}) {
		t.Fatalf("body = %#v, want %#v", body, map[string]string{"error": "unknown_erp_source"})
	}
}

func TestHandlerPutSellableAssortmentUnknownSourceReturnsBadRequest(t *testing.T) {
	store := &fakeStore{setSellableAssortmentErr: tenant_config.ErrUnknownActiveSource}
	req := httptest.NewRequest(http.MethodPut, "/config/sellable-assortment", strings.NewReader(
		`{"only_revenda":true,"only_em_estoque":false,"only_ecommerce_eligible":true}`))
	rec := httptest.NewRecorder()
	newMux(store).ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d, body=%s", rec.Code, http.StatusBadRequest, rec.Body.String())
	}
	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if !reflect.DeepEqual(body, map[string]string{"error": "unknown_erp_source"}) {
		t.Fatalf("body = %#v, want %#v", body, map[string]string{"error": "unknown_erp_source"})
	}
	if store.getCalls != 0 {
		t.Fatalf("Get calls = %d, want 0 after SetSellableAssortment error", store.getCalls)
	}
}

func TestHandlerPutSellableAssortmentUnknownSourceOnReadBackReturnsBadRequest(t *testing.T) {
	store := &fakeStore{
		postSetSet: true,
		postSetErr: tenant_config.ErrUnknownActiveSource,
	}
	req := httptest.NewRequest(http.MethodPut, "/config/sellable-assortment", strings.NewReader(
		`{"only_revenda":true,"only_em_estoque":false,"only_ecommerce_eligible":true}`))
	rec := httptest.NewRecorder()
	newMux(store).ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d for unknown_erp_source, body=%s", rec.Code, http.StatusBadRequest, rec.Body.String())
	}
	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if !reflect.DeepEqual(body, map[string]string{"error": "unknown_erp_source"}) {
		t.Fatalf("body = %#v, want %#v", body, map[string]string{"error": "unknown_erp_source"})
	}
	if !store.setSellableAssortmentCalled {
		t.Fatalf("SetSellableAssortment must be called before read-back")
	}
}

func TestHandlerGetSellableAssortmentStoreFailureReturnsInternalError(t *testing.T) {
	errBoom := errors.New("boom")
	store := &fakeStore{getErr: errBoom}
	req := httptest.NewRequest(http.MethodGet, "/config/sellable-assortment", nil)
	rec := httptest.NewRecorder()
	newMux(store).ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d, body=%s", rec.Code, http.StatusInternalServerError, rec.Body.String())
	}
	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if !reflect.DeepEqual(body, map[string]string{"error": "internal_error"}) {
		t.Fatalf("body = %#v, want %#v", body, map[string]string{"error": "internal_error"})
	}
}

func TestHandlerPutSellableAssortmentStoreFailureReturnsInternalError(t *testing.T) {
	errBoom := errors.New("boom")
	store := &fakeStore{setSellableAssortmentErr: errBoom}
	req := httptest.NewRequest(http.MethodPut, "/config/sellable-assortment", strings.NewReader(
		`{"only_revenda":true,"only_em_estoque":false,"only_ecommerce_eligible":true}`))
	rec := httptest.NewRecorder()
	newMux(store).ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d, body=%s", rec.Code, http.StatusInternalServerError, rec.Body.String())
	}
	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if !reflect.DeepEqual(body, map[string]string{"error": "internal_error"}) {
		t.Fatalf("body = %#v, want %#v", body, map[string]string{"error": "internal_error"})
	}
}

func TestHandlerPutSellableAssortmentReadBackFailureReturnsInternalError(t *testing.T) {
	errBoom := errors.New("boom")
	store := &fakeStore{postSetSet: true, postSetErr: errBoom}
	req := httptest.NewRequest(http.MethodPut, "/config/sellable-assortment", strings.NewReader(
		`{"only_revenda":true,"only_em_estoque":false,"only_ecommerce_eligible":true}`))
	rec := httptest.NewRecorder()
	newMux(store).ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d, body=%s", rec.Code, http.StatusInternalServerError, rec.Body.String())
	}
	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if !reflect.DeepEqual(body, map[string]string{"error": "internal_error"}) {
		t.Fatalf("body = %#v, want %#v", body, map[string]string{"error": "internal_error"})
	}
	if !store.setSellableAssortmentCalled {
		t.Fatalf("SetSellableAssortment must be called before read-back")
	}
}

func TestHandlerPutSellableAssortmentMissingBooleanReturnsInvalidBody(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{name: "only_revenda", body: `{"only_em_estoque":false,"only_ecommerce_eligible":true}`},
		{name: "only_em_estoque", body: `{"only_revenda":true,"only_ecommerce_eligible":true}`},
		{name: "only_ecommerce_eligible", body: `{"only_revenda":true,"only_em_estoque":false}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &fakeStore{}
			req := httptest.NewRequest(http.MethodPut, "/config/sellable-assortment", strings.NewReader(tt.body))
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
			if store.setSellableAssortmentCalled {
				t.Fatalf("SetSellableAssortment must not be called for missing %s", tt.name)
			}
		})
	}
}

func TestHandlerPutSellableAssortmentMalformedBody(t *testing.T) {
	store := &fakeStore{}
	req := httptest.NewRequest(http.MethodPut, "/config/sellable-assortment", strings.NewReader(`{not json`))
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
	if store.setSellableAssortmentCalled {
		t.Fatalf("SetSellableAssortment must not be called on malformed body")
	}
}

func TestHandlerPutSellableAssortmentTrailingGarbageReturnsInvalidBody(t *testing.T) {
	store := &fakeStore{}
	req := httptest.NewRequest(http.MethodPut, "/config/sellable-assortment", strings.NewReader(
		`{"only_revenda":true,"only_em_estoque":false,"only_ecommerce_eligible":true}{"only_revenda":false}`))
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
	if store.setSellableAssortmentCalled {
		t.Fatalf("SetSellableAssortment must not be called on trailing garbage")
	}
}
