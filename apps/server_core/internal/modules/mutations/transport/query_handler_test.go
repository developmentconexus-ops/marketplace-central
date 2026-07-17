package transport

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/modules/mutations/application"
	"marketplace-central/apps/server_core/internal/modules/mutations/domain"
	"marketplace-central/apps/server_core/internal/modules/mutations/ports"
)

type mutationReadFake struct {
	protocolPage application.ProtocolPage
	itemPage     application.ItemPage
	protocol     ports.Protocol
	protocolErr  error
	listErr      error
	itemsErr     error
	protocolQ    application.ProtocolQuery
	itemQ        application.ItemQuery
	getID        string
}

func (f *mutationReadFake) ListProtocols(_ context.Context, query application.ProtocolQuery) (application.ProtocolPage, error) {
	f.protocolQ = query
	return f.protocolPage, f.listErr
}

func (f *mutationReadFake) GetProtocol(_ context.Context, id string) (ports.Protocol, error) {
	f.getID = id
	return f.protocol, f.protocolErr
}

func (f *mutationReadFake) ListItems(_ context.Context, query application.ItemQuery) (application.ItemPage, error) {
	f.itemQ = query
	return f.itemPage, f.itemsErr
}

func TestQueryHandlerListRequiresInstallationAndValidatesIC03Filters(t *testing.T) {
	fake := &mutationReadFake{}
	h := NewQueryHandler(fake)
	tests := []struct {
		name string
		url  string
		code string
	}{
		{name: "missing installation", url: "/mutations", code: "installation_required"},
		{name: "invalid state", url: "/mutations?installation_id=inst-1&state=unknown", code: "invalid_filter"},
		{name: "invalid type", url: "/mutations?installation_id=inst-1&type=listing_create", code: "invalid_filter"},
		{name: "invalid limit", url: "/mutations?installation_id=inst-1&limit=201", code: "invalid_filter"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			h.HandleList(w, httptest.NewRequest(http.MethodGet, tc.url, nil))
			if w.Code != http.StatusBadRequest || !strings.Contains(w.Body.String(), `"code":"`+tc.code+`"`) {
				t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
			}
		})
	}

	created := time.Date(2026, 7, 16, 12, 0, 0, 0, time.UTC)
	fake.protocolPage = application.ProtocolPage{Items: []ports.Protocol{{ProtocolID: "MP-2"}}, NextCursor: &application.ProtocolCursor{CreatedAt: created, ProtocolID: "MP-2"}}
	w := httptest.NewRecorder()
	h.HandleList(w, httptest.NewRequest(http.MethodGet, "/mutations?installation_id=inst-1&state=applied&type=price_update&limit=20", nil))
	if w.Code != http.StatusOK || fake.protocolQ.InstallationID != "inst-1" || fake.protocolQ.State != domain.ProtocolStateApplied || fake.protocolQ.Type != domain.ProtocolTypePriceUpdate || fake.protocolQ.Limit != 20 {
		t.Fatalf("status=%d query=%+v body=%s", w.Code, fake.protocolQ, w.Body.String())
	}
}

func TestQueryHandlerProtocolCursorRoundTripsAndPreservesStableKeyset(t *testing.T) {
	created := time.Date(2026, 7, 16, 12, 0, 0, 123, time.UTC)
	fake := &mutationReadFake{protocolPage: application.ProtocolPage{NextCursor: &application.ProtocolCursor{CreatedAt: created, ProtocolID: "MP-000002"}}}
	h := NewQueryHandler(fake)
	w := httptest.NewRecorder()
	h.HandleList(w, httptest.NewRequest(http.MethodGet, "/mutations?installation_id=inst-1", nil))
	var response mutationPageEnvelope
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil || response.NextCursor == nil {
		t.Fatalf("status=%d body=%s err=%v", w.Code, w.Body.String(), err)
	}
	w = httptest.NewRecorder()
	h.HandleList(w, httptest.NewRequest(http.MethodGet, "/mutations?installation_id=inst-1&cursor="+*response.NextCursor, nil))
	if w.Code != http.StatusOK || !fake.protocolQ.Cursor.CreatedAt.Equal(created) || fake.protocolQ.Cursor.ProtocolID != "MP-000002" {
		t.Fatalf("status=%d query=%+v body=%s", w.Code, fake.protocolQ, w.Body.String())
	}
}

func TestQueryHandlerRejectsMalformedCursorAndDuplicateFilters(t *testing.T) {
	h := NewQueryHandler(&mutationReadFake{})
	for _, tc := range []struct{ url, code string }{
		{"/mutations?installation_id=inst-1&cursor=bad", "invalid_cursor"},
		{"/mutations?installation_id=inst-1&state=applied&state=failed", "invalid_filter"},
		{"/mutations?installation_id=inst-1&filter.state=applied", "invalid_filter"},
	} {
		w := httptest.NewRecorder()
		h.HandleList(w, httptest.NewRequest(http.MethodGet, tc.url, nil))
		if w.Code != http.StatusBadRequest || !strings.Contains(w.Body.String(), `"code":"`+tc.code+`"`) {
			t.Fatalf("url=%s status=%d body=%s", tc.url, w.Code, w.Body.String())
		}
	}
}

func TestQueryHandlerDetailAndItemsUsePathIDOnlyAndMapNotFound(t *testing.T) {
	fake := &mutationReadFake{protocolErr: application.ErrProtocolNotFound}
	h := NewQueryHandler(fake)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/mutations/MP-000001?installation_id=other-tenant", nil)
	req.SetPathValue("id", "MP-000001")
	h.HandleGet(w, req)
	if w.Code != http.StatusNotFound || fake.getID != "MP-000001" || !strings.Contains(w.Body.String(), `"code":"protocol_not_found"`) {
		t.Fatalf("status=%d id=%q body=%s", w.Code, fake.getID, w.Body.String())
	}

	fake.protocolErr = nil
	fake.itemPage = application.ItemPage{NextCursor: &application.ItemCursor{Seq: 4}}
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/mutations/MP-000001/items?state=failed&limit=10", nil)
	req.SetPathValue("id", "MP-000001")
	h.HandleItems(w, req)
	if w.Code != http.StatusOK || fake.itemQ.ProtocolID != "MP-000001" || fake.itemQ.State != domain.ItemStateFailed || fake.itemQ.Limit != 10 {
		t.Fatalf("status=%d query=%+v body=%s", w.Code, fake.itemQ, w.Body.String())
	}
}

func TestQueryHandlerDoesNotSerializeRequestHeaders(t *testing.T) {
	fake := &mutationReadFake{protocol: ports.Protocol{ProtocolID: "MP-000001", InstallationID: "inst-1", Actor: "operator"}}
	h := NewQueryHandler(fake)
	req := httptest.NewRequest(http.MethodGet, "/mutations/MP-000001", nil)
	req.SetPathValue("id", "MP-000001")
	req.Header.Set("Authorization", "Bearer secret-token")
	w := httptest.NewRecorder()
	h.HandleGet(w, req)
	if w.Code != http.StatusOK || strings.Contains(w.Body.String(), "secret-token") || strings.Contains(w.Body.String(), "Authorization") {
		t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
	}
}
