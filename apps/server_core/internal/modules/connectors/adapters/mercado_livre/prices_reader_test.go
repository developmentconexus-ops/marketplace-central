package mercadolivre

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"marketplace-central/apps/server_core/internal/modules/connectors/domain"
)

// TestGetItemPricesMapsTypedFieldsAndRaw covers the F-02 ownership deliverable
// (feature brief point 5 / P5 F-9): GET /items/{id}/prices is the dedicated
// reader M-05 will consume for fee-layer pricing, since the multiget endpoint
// does not carry sale_price/commission data (see file header finding).
func TestGetItemPricesMapsTypedFieldsAndRaw(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/items/MLB-PRICE/prices" {
			t.Fatalf("request = %s %s", r.Method, r.URL.Path)
		}
		_, _ = io.WriteString(w, `{"id":"MLB-PRICE","prices":[
			{"type":"standard","amount":199.90,"regular_amount":249.90,"currency_id":"BRL"},
			{"type":"promotion","amount":179.90,"currency_id":"BRL"}
		]}`)
	}))
	defer server.Close()

	dto, err := pricingTestAdapter(server.URL, fixedNow()).GetItemPrices(context.Background(), pricingAccountRef(), "MLB-PRICE")
	if err != nil {
		t.Fatalf("GetItemPrices() error = %v", err)
	}
	if dto.ProviderItemID != "MLB-PRICE" {
		t.Fatalf("ProviderItemID = %q, want MLB-PRICE", dto.ProviderItemID)
	}
	if !dto.FetchedAt.Equal(fixedNow()) {
		t.Fatalf("FetchedAt = %v, want %v", dto.FetchedAt, fixedNow())
	}
	if len(dto.Raw) == 0 {
		t.Fatal("Raw is empty, want non-empty on success")
	}
	if dto.RawTruncated {
		t.Fatal("RawTruncated = true, want false (small body)")
	}
	if len(dto.Prices) != 2 {
		t.Fatalf("Prices count = %d, want 2", len(dto.Prices))
	}
	standard := dto.Prices[0]
	if standard.Type != "standard" || standard.Amount == nil || *standard.Amount != "199.90" || standard.RegularAmount == nil || *standard.RegularAmount != "249.90" || standard.Currency != "BRL" {
		t.Fatalf("Prices[0] = %+v", standard)
	}
	promotion := dto.Prices[1]
	if promotion.Type != "promotion" || promotion.Amount == nil || *promotion.Amount != "179.90" || promotion.RegularAmount != nil {
		t.Fatalf("Prices[1] = %+v", promotion)
	}
}

func TestGetItemPricesMapsProviderErrors(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		status int
		want   error
	}{
		{http.StatusUnauthorized, domain.ErrUnauthorized},
		{http.StatusNotFound, domain.ErrNotFound},
		{http.StatusTooManyRequests, domain.ErrRateLimited},
		{http.StatusInternalServerError, domain.ErrProviderUnavailable},
	} {
		t.Run(fmt.Sprint(test.status), func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(test.status) }))
			defer server.Close()

			_, err := pricingTestAdapter(server.URL, fixedNow()).GetItemPrices(context.Background(), pricingAccountRef(), "MLB-ERR")
			if !errors.Is(err, test.want) {
				t.Fatalf("GetItemPrices() error = %v, want errors.Is(_, %v)", err, test.want)
			}
		})
	}
}

func TestGetItemPricesUndecodableBodyFails(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, `not-json-at-all`)
	}))
	defer server.Close()

	_, err := pricingTestAdapter(server.URL, fixedNow()).GetItemPrices(context.Background(), pricingAccountRef(), "MLB-BAD")
	if err == nil {
		t.Fatal("error = nil, want a decode error")
	}
	if domain.ErrorCodeOf(err) != domain.ErrCodeProviderPayloadInvalid {
		t.Fatalf("error code = %q, want %q", domain.ErrorCodeOf(err), domain.ErrCodeProviderPayloadInvalid)
	}
}

// TestGetItemPricesOversizeBodyMarksRawTruncated mirrors the multiget
// per-item truncation contract (see items_multiget_reader.go's
// itemMultigetRawCap doc) for this single-item reader.
func TestGetItemPricesOversizeBodyMarksRawTruncated(t *testing.T) {
	t.Parallel()

	filler := strings.Repeat("y", itemMultigetRawCap+1024)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprintf(w, `{"id":"MLB-BIG","prices":[{"type":"standard","amount":10,"currency_id":"BRL"}],"_filler":%q}`, filler)
	}))
	defer server.Close()

	dto, err := pricingTestAdapter(server.URL, fixedNow()).GetItemPrices(context.Background(), pricingAccountRef(), "MLB-BIG")
	if err != nil {
		t.Fatalf("GetItemPrices() error = %v", err)
	}
	if !dto.RawTruncated {
		t.Fatal("RawTruncated = false, want true for an oversize body")
	}
	if len(dto.Raw) != itemMultigetRawCap {
		t.Fatalf("len(Raw) = %d, want exactly the %d-byte cap", len(dto.Raw), itemMultigetRawCap)
	}
	// The typed fields decode fine even though Raw is capped: capping is a
	// COPY taken after the typed decode, never a mutation of the bytes fed to
	// json.Unmarshal.
	if len(dto.Prices) != 1 || dto.Prices[0].Type != "standard" {
		t.Fatalf("Prices = %#v, want typed decode to survive raw truncation", dto.Prices)
	}
}

func TestGetItemPricesPublicEntryPointResolvesTokenAndAccount(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer test-token" {
			t.Fatalf("authorization = %q", got)
		}
		_, _ = io.WriteString(w, `{"id":"MLB-1","prices":[]}`)
	}))
	defer server.Close()

	dto, err := pricingTestAdapter(server.URL, fixedNow()).GetItemPrices(context.Background(), pricingAccountRef(), "MLB-1")
	if err != nil {
		t.Fatalf("GetItemPrices() error = %v", err)
	}
	if dto.ProviderItemID != "MLB-1" {
		t.Fatalf("ProviderItemID = %q, want MLB-1", dto.ProviderItemID)
	}
}
