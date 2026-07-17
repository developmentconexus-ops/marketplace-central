package transport

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"marketplace-central/apps/server_core/internal/modules/mutations/application"
	"marketplace-central/apps/server_core/internal/modules/mutations/domain"
	"marketplace-central/apps/server_core/internal/modules/mutations/ports"
)

type commandServiceFake struct {
	createResult ports.Protocol
	createErr    error
	preview      application.PreviewResult
	previewErr   error
	createCalls  int
	createInput  application.CreateInput
}

func (f *commandServiceFake) Create(_ context.Context, input application.CreateInput) (ports.Protocol, error) {
	f.createCalls++
	f.createInput = input
	return f.createResult, f.createErr
}

func (f *commandServiceFake) Preview(_ context.Context, _ string) (application.PreviewResult, error) {
	return f.preview, f.previewErr
}

type approvalServiceFake struct {
	approveResult ports.Protocol
	approveErr    error
	cancelResult  ports.Protocol
	cancelErr     error
	approveBody   json.RawMessage
}

func (f *approvalServiceFake) Approve(_ context.Context, _ string, body json.RawMessage) (ports.Protocol, error) {
	f.approveBody = append(f.approveBody[:0], body...)
	return f.approveResult, f.approveErr
}

func (f *approvalServiceFake) Cancel(_ context.Context, _ string) (ports.Protocol, error) {
	return f.cancelResult, f.cancelErr
}

type retryServiceFake struct {
	result ports.Protocol
	err    error
}

func (f *retryServiceFake) Retry(context.Context, string) (ports.Protocol, error) {
	return f.result, f.err
}

func TestCommandHandlersMapStatusAndCode(t *testing.T) {
	validCreate := `{"installation_id":"inst-1","type":"price_update","selection":{"mode":"explicit","listing_ids":["inst-1~MLB1~-"]},"intent":{"new_price":{"amount":"10.00","currency":"BRL"}},"actor":"operator"}`
	tests := []struct {
		name       string
		path       string
		body       string
		commands   *commandServiceFake
		approval   *approvalServiceFake
		retry      *retryServiceFake
		wantStatus int
		wantCode   string
	}{
		{name: "invalid json", path: "/mutations", body: `{`, commands: &commandServiceFake{}, wantStatus: http.StatusBadRequest, wantCode: "invalid_body"},
		{name: "unknown field", path: "/mutations", body: validCreate[:len(validCreate)-1] + `,"extra":true}`, commands: &commandServiceFake{}, wantStatus: http.StatusBadRequest, wantCode: "invalid_body"},
		{name: "trailing content", path: "/mutations", body: validCreate + `{}`, commands: &commandServiceFake{}, wantStatus: http.StatusBadRequest, wantCode: "invalid_body"},
		{name: "invalid intent", path: "/mutations", body: strings.Replace(validCreate, `"intent":{"new_price":{"amount":"10.00","currency":"BRL"}}`, `"intent":[]`, 1), commands: &commandServiceFake{}, wantStatus: http.StatusBadRequest, wantCode: "invalid_intent"},
		{name: "actor gate", path: "/mutations", body: validCreate, commands: &commandServiceFake{createErr: &application.GateError{Code: application.FailureCodeActorRequired, MessagePT: "ator é obrigatório"}}, wantStatus: http.StatusUnprocessableEntity, wantCode: "actor_required"},
		{name: "disabled listing create", path: "/mutations", body: strings.Replace(validCreate, `"type":"price_update"`, `"type":"listing_create"`, 1), commands: &commandServiceFake{}, wantStatus: http.StatusUnprocessableEntity, wantCode: "type_not_enabled"},
		{name: "empty selection", path: "/mutations/MP-1/preview", commands: &commandServiceFake{previewErr: &application.GateError{Code: application.FailureCodeEmptySelection, MessagePT: "vazia"}}, wantStatus: http.StatusUnprocessableEntity, wantCode: "empty_selection"},
		{name: "large selection", path: "/mutations/MP-1/preview", commands: &commandServiceFake{previewErr: &application.GateError{Code: application.FailureCodeSelectionTooLarge, MessagePT: "grande"}}, wantStatus: http.StatusUnprocessableEntity, wantCode: "selection_too_large"},
		{name: "execute gate", path: "/mutations/MP-1/approve", body: `{}`, approval: &approvalServiceFake{approveErr: &application.GateError{Code: application.FailureCodeExecuteRequired, MessagePT: "execute"}}, wantStatus: http.StatusUnprocessableEntity, wantCode: "execute_required"},
		{name: "stale preview", path: "/mutations/MP-1/approve", body: `{"execute":true}`, approval: &approvalServiceFake{approveErr: &application.GateError{Code: application.FailureCodePreviewStale, MessagePT: "stale"}}, wantStatus: http.StatusConflict, wantCode: "preview_stale"},
		{name: "illegal state", path: "/mutations/MP-1/cancel", approval: &approvalServiceFake{cancelErr: &domain.InvalidStateTransitionError{Code: domain.InvalidStateCode, From: domain.ProtocolStateApproved, To: domain.ProtocolStateCancelled}}, wantStatus: http.StatusConflict, wantCode: "invalid_state"},
		{name: "unknown protocol", path: "/mutations/MP-missing/preview", commands: &commandServiceFake{previewErr: application.ErrProtocolNotFound}, wantStatus: http.StatusNotFound, wantCode: "protocol_not_found"},
		{name: "nothing to retry", path: "/mutations/MP-1/retry", retry: &retryServiceFake{err: &application.GateError{Code: application.FailureCodeNothingToRetry, MessagePT: "nada"}}, wantStatus: http.StatusConflict, wantCode: "nothing_to_retry"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.commands == nil {
				tc.commands = &commandServiceFake{}
			}
			if tc.approval == nil {
				tc.approval = &approvalServiceFake{}
			}
			if tc.retry == nil {
				tc.retry = &retryServiceFake{}
			}
			h := NewHandler(tc.commands, tc.approval, tc.retry)
			req := httptest.NewRequest(http.MethodPost, tc.path, strings.NewReader(tc.body))
			w := httptest.NewRecorder()
			mux := http.NewServeMux()
			h.Register(mux)
			mux.ServeHTTP(w, req)

			if w.Code != tc.wantStatus {
				t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
			}
			var envelope struct {
				Error struct {
					Code    string         `json:"code"`
					Message string         `json:"message"`
					Details map[string]any `json:"details"`
				} `json:"error"`
			}
			if err := json.Unmarshal(w.Body.Bytes(), &envelope); err != nil {
				t.Fatal(err)
			}
			if envelope.Error.Code != tc.wantCode || envelope.Error.Details == nil {
				t.Fatalf("error=%+v body=%s", envelope.Error, w.Body.String())
			}
		})
	}
}

func TestCommandHandlerUsesBodyActorAndDoesNotPersistListingCreate(t *testing.T) {
	commands := &commandServiceFake{createResult: ports.Protocol{ProtocolID: "MP-1"}}
	h := NewHandler(commands, &approvalServiceFake{}, &retryServiceFake{})
	req := httptest.NewRequest(http.MethodPost, "/mutations", strings.NewReader(`{"installation_id":"inst-1","type":"price_update","selection":{"mode":"explicit","listing_ids":["inst-1~MLB1~-"]},"intent":{"new_price":{}},"actor":"body-actor"}`))
	req.Header.Set("X-Actor", "header-actor")
	w := httptest.NewRecorder()
	h.HandleCreate(w, req)
	if w.Code != http.StatusCreated || commands.createInput.Actor != "body-actor" {
		t.Fatalf("status=%d actor=%q body=%s", w.Code, commands.createInput.Actor, w.Body.String())
	}

	commands.createCalls = 0
	req = httptest.NewRequest(http.MethodPost, "/mutations", strings.NewReader(`{"installation_id":"inst-1","type":"listing_create","selection":{"mode":"explicit","listing_ids":["inst-1~MLB1~-"]},"intent":{"product_id":1,"listing_type_code":"gold","price":10,"category_id":"cat","attributes":[]},"actor":"body-actor"}`))
	w = httptest.NewRecorder()
	h.HandleCreate(w, req)
	if w.Code != http.StatusUnprocessableEntity || commands.createCalls != 0 {
		t.Fatalf("status=%d create_calls=%d body=%s", w.Code, commands.createCalls, w.Body.String())
	}
}

func TestCommandHandlerRegistersFiveSeparatePostRoutes(t *testing.T) {
	commands := &commandServiceFake{}
	approval := &approvalServiceFake{}
	retry := &retryServiceFake{}
	h := NewHandler(commands, approval, retry)
	mux := http.NewServeMux()
	h.Register(mux)

	for _, path := range []string{"/mutations/MP-1/preview", "/mutations/MP-1/approve", "/mutations/MP-1/cancel", "/mutations/MP-1/retry"} {
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, httptest.NewRequest(http.MethodPost, path, strings.NewReader(`{}`)))
		if w.Code == http.StatusNotFound {
			t.Fatalf("route %s was not registered", path)
		}
	}
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/mutations/MP-1/preview/approve", strings.NewReader(`{}`)))
	if w.Code != http.StatusNotFound {
		t.Fatalf("combined route status=%d body=%s", w.Code, w.Body.String())
	}
}

func TestCommandHandlerRejectsOptionalCommandBodyWithUnknownFields(t *testing.T) {
	h := NewHandler(&commandServiceFake{}, &approvalServiceFake{}, &retryServiceFake{})
	mux := http.NewServeMux()
	h.Register(mux)
	for _, path := range []string{"/mutations/MP-1/preview", "/mutations/MP-1/cancel", "/mutations/MP-1/retry"} {
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, httptest.NewRequest(http.MethodPost, path, strings.NewReader(`{"unknown":true}`)))
		if w.Code != http.StatusBadRequest {
			t.Fatalf("path=%s status=%d body=%s", path, w.Code, w.Body.String())
		}
	}
}

func TestCommandErrorMappingDoesNotCollapseUnknownErrors(t *testing.T) {
	h := NewHandler(&commandServiceFake{previewErr: errors.New("database unavailable")}, &approvalServiceFake{}, &retryServiceFake{})
	w := httptest.NewRecorder()
	h.HandlePreview(w, httptest.NewRequest(http.MethodPost, "/mutations/MP-1/preview", nil))
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
	}
}
