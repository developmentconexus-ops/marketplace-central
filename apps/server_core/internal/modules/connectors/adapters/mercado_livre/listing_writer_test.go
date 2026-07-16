package mercadolivre

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"marketplace-central/apps/server_core/internal/modules/connectors/domain"
)

func TestUpdateListingSendsCanonicalPauseAndEdit(t *testing.T) {
	tests := []struct {
		name string
		req  domain.ListingWriteRequest
		want map[string]any
	}{
		{"pause", listingRequest(domain.ListingWritePause, nil), map[string]any{"status": "paused"}},
		{"edit", listingRequest(domain.ListingWriteEdit, []domain.ListingAttribute{{ID: "BRAND", ValueName: "Acme"}}), map[string]any{"attributes": []any{map[string]any{"id": "BRAND", "value_name": "Acme"}}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var calls int
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls++
				if r.Method != http.MethodPut || r.URL.Path != "/items/MLB-1" {
					t.Errorf("request = %s %s", r.Method, r.URL.Path)
				}
				if got := r.Header.Get("X-Idempotency-Key"); got != "MP-1:MLB-1" {
					t.Errorf("idempotency = %q", got)
				}
				var got map[string]any
				if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
					t.Fatal(err)
				}
				if encoded, _ := json.Marshal(got); string(encoded) != mustJSON(tt.want) {
					t.Errorf("body = %s, want %s", encoded, mustJSON(tt.want))
				}
				status := "active"
				if tt.req.Action == domain.ListingWritePause {
					status = "paused"
				}
				_, _ = io.WriteString(w, `{"id":"MLB-1","status":"`+status+`"}`)
			}))
			defer server.Close()
			result, err := testListingAdapter(server).UpdateListing(context.Background(), tt.req)
			if err != nil {
				t.Fatalf("UpdateListing() error = %v", err)
			}
			if calls != 1 || result.Result != domain.WriteResultApplied || result.Action != tt.req.Action {
				t.Fatalf("calls=%d result=%+v", calls, result)
			}
		})
	}
}

func TestUpdateListingMapsRemoteStateHonestly(t *testing.T) {
	for _, status := range []string{"paused", "closed"} {
		t.Run(status, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				_, _ = io.WriteString(w, `{"id":"MLB-1","status":"`+status+`"}`)
			}))
			defer server.Close()
			result, err := testListingAdapter(server).UpdateListing(context.Background(), listingRequest(domain.ListingWriteEdit, []domain.ListingAttribute{{ID: "BRAND", ValueName: "Acme"}}))
			if err != nil {
				t.Fatal(err)
			}
			if result.Result != domain.WriteResultRejected || !strings.Contains(result.Message, status) {
				t.Fatalf("result = %+v", result)
			}
		})
	}
}

func TestUpdateListingRejectsMalformedEditBeforeNetwork(t *testing.T) {
	var calls int
	server := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { calls++ }))
	defer server.Close()
	for _, attrs := range [][]domain.ListingAttribute{nil, {{ID: "", ValueName: "Acme"}}, {{ID: "BRAND", ValueName: ""}}, {{ID: "BRAND", ValueName: "A"}, {ID: "BRAND", ValueName: "B"}}} {
		_, err := testListingAdapter(server).UpdateListing(context.Background(), listingRequest(domain.ListingWriteEdit, attrs))
		if domain.ErrorCodeOf(err) != domain.ErrCodeProviderValidation {
			t.Errorf("attrs=%+v err=%v", attrs, err)
		}
	}
	if calls != 0 {
		t.Fatalf("provider calls = %d", calls)
	}
}

func TestUpdateListingMapsAndSanitizesProviderFailures(t *testing.T) {
	for _, tt := range []struct {
		status int
		code   domain.ErrorCode
	}{{401, domain.ErrCodeProviderAuth}, {403, domain.ErrCodeProviderAuth}, {429, domain.ErrCodeProviderRateLimited}, {400, domain.ErrCodeProviderValidation}, {502, domain.ErrCodeProviderTransient}, {302, domain.ErrorCode("CONNECTORS_INTERNAL")}} {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(tt.status)
			_, _ = io.WriteString(w, `{"message":"bad token=secret-value"}`)
		}))
		_, err := testListingAdapter(server).UpdateListing(context.Background(), listingRequest(domain.ListingWritePause, nil))
		server.Close()
		if domain.ErrorCodeOf(err) != tt.code || strings.Contains(err.Error(), "secret-value") || strings.Contains(err.Error(), `{"message"`) {
			t.Fatalf("status=%d err=%v", tt.status, err)
		}
	}
}

func TestUpdateListingMapsTimeoutAsUnavailable(t *testing.T) {
	adapter := NewCapabilityAdapter(CapabilityAdapterConfig{BaseURL: "http://provider.invalid", HTTPClient: &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) { return nil, context.DeadlineExceeded })}, AccessTokenResolver: func(context.Context, domain.ProviderAccountRef) (string, error) { return "token", nil }})
	_, err := adapter.UpdateListing(context.Background(), listingRequest(domain.ListingWritePause, nil))
	if domain.ErrorCodeOf(err) != domain.ErrCodeProviderTransient || strings.Contains(err.Error(), "token") {
		t.Fatalf("err = %v", err)
	}
}

func listingRequest(action domain.ListingWriteAction, attrs []domain.ListingAttribute) domain.ListingWriteRequest {
	return domain.ListingWriteRequest{TenantID: "tenant-1", InstallationID: "installation-1", ListingID: "MLB-1", IdempotencyKey: "MP-1:MLB-1", Action: action, Attributes: attrs}
}
func testListingAdapter(server *httptest.Server) *CapabilityAdapter {
	return NewCapabilityAdapter(CapabilityAdapterConfig{BaseURL: server.URL, HTTPClient: server.Client(), AccessTokenResolver: func(context.Context, domain.ProviderAccountRef) (string, error) { return "token", nil }})
}
func mustJSON(v any) string { b, _ := json.Marshal(v); return string(b) }
