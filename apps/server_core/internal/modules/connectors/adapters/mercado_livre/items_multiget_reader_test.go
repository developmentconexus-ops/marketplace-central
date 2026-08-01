package mercadolivre

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/modules/connectors/domain"
)

// fixedNow returns a stable, arbitrary "now" for tests that do not assert on
// FetchedAt — avoids importing time.Now() flakiness into every call site.
func fixedNow() time.Time {
	return time.Date(2026, 7, 31, 12, 0, 0, 0, time.UTC)
}

// TestGetItemsMultigetPartitionsInBatchesOf20 is the headline F-02 contract
// (fact #13, research/external-ml-api-facts.md): 45 ids must produce EXACTLY 3
// HTTP calls (20+20+5), each observable in the request's `ids` query string,
// and the returned DTOs must preserve the ORIGINAL global id order — not the
// per-batch order, the order the caller passed in across all batches.
func TestGetItemsMultigetPartitionsInBatchesOf20(t *testing.T) {
	t.Parallel()

	ids := make([]string, 45)
	for i := range ids {
		ids[i] = "MLB" + strconv.Itoa(1000+i)
	}

	var mu sync.Mutex
	var batchSizes []int
	callCount := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/items" {
			t.Fatalf("request = %s %s", r.Method, r.URL.Path)
		}
		requested := strings.Split(r.URL.Query().Get("ids"), ",")

		mu.Lock()
		callCount++
		batchSizes = append(batchSizes, len(requested))
		mu.Unlock()

		var sb strings.Builder
		sb.WriteByte('[')
		for i, id := range requested {
			if i > 0 {
				sb.WriteByte(',')
			}
			fmt.Fprintf(&sb, `{"code":200,"body":{"id":%q,"title":"item %s","status":"active"}}`, id, id)
		}
		sb.WriteByte(']')
		_, _ = io.WriteString(w, sb.String())
	}))
	defer server.Close()

	dtos, err := pricingTestAdapter(server.URL, fixedNow()).getItemsMultiget(context.Background(), pricingAccountRef(), "test-token", ids)
	if err != nil {
		t.Fatalf("getItemsMultiget() error = %v", err)
	}

	mu.Lock()
	gotCallCount := callCount
	gotBatchSizes := append([]int(nil), batchSizes...)
	mu.Unlock()

	if gotCallCount != 3 {
		t.Fatalf("HTTP calls = %d, want 3", gotCallCount)
	}
	if want := []int{20, 20, 5}; !equalInts(gotBatchSizes, want) {
		t.Fatalf("batch sizes = %v, want %v", gotBatchSizes, want)
	}
	if len(dtos) != 45 {
		t.Fatalf("dto count = %d, want 45", len(dtos))
	}
	for i, dto := range dtos {
		if dto.ProviderItemID != ids[i] {
			t.Fatalf("dto[%d].ProviderItemID = %q, want %q (global order must match input order)", i, dto.ProviderItemID, ids[i])
		}
		if dto.Err != nil {
			t.Fatalf("dto[%d].Err = %v, want nil", i, dto.Err)
		}
	}
}

// TestGetItemsMultigetSuccessfulDTOsHaveNonEmptyRaw covers the DTO Raw
// contract: every successfully-returned item must carry its own individual
// raw body bytes (Raw json.RawMessage), never empty on success.
func TestGetItemsMultigetSuccessfulDTOsHaveNonEmptyRaw(t *testing.T) {
	t.Parallel()

	ids := []string{"MLB1", "MLB2", "MLB3"}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, `[
			{"code":200,"body":{"id":"MLB1","title":"one"}},
			{"code":200,"body":{"id":"MLB2","title":"two"}},
			{"code":200,"body":{"id":"MLB3","title":"three"}}
		]`)
	}))
	defer server.Close()

	dtos, err := pricingTestAdapter(server.URL, fixedNow()).getItemsMultiget(context.Background(), pricingAccountRef(), "test-token", ids)
	if err != nil {
		t.Fatalf("getItemsMultiget() error = %v", err)
	}
	if len(dtos) != 3 {
		t.Fatalf("dto count = %d, want 3", len(dtos))
	}
	for i, dto := range dtos {
		if dto.Err != nil {
			t.Fatalf("dto[%d].Err = %v, want nil", i, dto.Err)
		}
		if len(dto.Raw) == 0 {
			t.Fatalf("dto[%d].Raw is empty, want non-empty on success", i)
		}
		if dto.RawTruncated {
			t.Fatalf("dto[%d].RawTruncated = true, want false (small body)", i)
		}
	}
}

// TestGetItemsMultigetPerItemErrorDoesNotFailBatch: a single id 404ing inside
// the multiget response must degrade to a per-item error in that DTO slot —
// the other 44 ids in the same batch must still decode cleanly and the
// function must not return a batch-level error.
func TestGetItemsMultigetPerItemErrorDoesNotFailBatch(t *testing.T) {
	t.Parallel()

	ids := make([]string, 45)
	for i := range ids {
		ids[i] = "MLB" + strconv.Itoa(2000+i)
	}
	const failingIndex = 17 // inside the first batch (0..19)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requested := strings.Split(r.URL.Query().Get("ids"), ",")
		var sb strings.Builder
		sb.WriteByte('[')
		for i, id := range requested {
			if i > 0 {
				sb.WriteByte(',')
			}
			globalIndex := indexOf(ids, id)
			if globalIndex == failingIndex {
				fmt.Fprintf(&sb, `{"code":404,"body":{"message":"item not found","error":"not_found"}}`)
				continue
			}
			fmt.Fprintf(&sb, `{"code":200,"body":{"id":%q,"title":"item %s"}}`, id, id)
		}
		sb.WriteByte(']')
		_, _ = io.WriteString(w, sb.String())
	}))
	defer server.Close()

	dtos, err := pricingTestAdapter(server.URL, fixedNow()).getItemsMultiget(context.Background(), pricingAccountRef(), "test-token", ids)
	if err != nil {
		t.Fatalf("getItemsMultiget() error = %v, want nil (per-item error must not fail the batch)", err)
	}
	if len(dtos) != 45 {
		t.Fatalf("dto count = %d, want 45", len(dtos))
	}

	validCount := 0
	for i, dto := range dtos {
		if i == failingIndex {
			if dto.Err == nil {
				t.Fatalf("dto[%d].Err = nil, want a per-item error", i)
			}
			if !errors.Is(dto.Err, domain.ErrNotFound) {
				t.Fatalf("dto[%d].Err = %v, want errors.Is(_, ErrNotFound)", i, dto.Err)
			}
			if dto.ProviderItemID != ids[failingIndex] {
				t.Fatalf("dto[%d].ProviderItemID = %q, want %q (id slot preserved even on error)", i, dto.ProviderItemID, ids[failingIndex])
			}
			continue
		}
		if dto.Err != nil {
			t.Fatalf("dto[%d].Err = %v, want nil", i, dto.Err)
		}
		validCount++
	}
	if validCount != 44 {
		t.Fatalf("valid dto count = %d, want 44", validCount)
	}
}

// TestGetItemsMultigetBodyNotJSONIsPerItemError: an item whose `body` is a bare
// JSON string (not the expected object shape) must degrade to a per-item
// error, while the batch and its sibling DTOs continue normally.
func TestGetItemsMultigetBodyNotJSONIsPerItemError(t *testing.T) {
	t.Parallel()

	ids := []string{"MLB1", "MLB2", "MLB3"}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, `[
			{"code":200,"body":{"id":"MLB1","title":"one"}},
			{"code":200,"body":"this-is-not-an-item-object"},
			{"code":200,"body":{"id":"MLB3","title":"three"}}
		]`)
	}))
	defer server.Close()

	dtos, err := pricingTestAdapter(server.URL, fixedNow()).getItemsMultiget(context.Background(), pricingAccountRef(), "test-token", ids)
	if err != nil {
		t.Fatalf("getItemsMultiget() error = %v, want nil (per-item decode failure must not fail the batch)", err)
	}
	if len(dtos) != 3 {
		t.Fatalf("dto count = %d, want 3", len(dtos))
	}
	if dtos[0].Err != nil || dtos[2].Err != nil {
		t.Fatalf("sibling DTOs must decode cleanly: dtos[0].Err=%v dtos[2].Err=%v", dtos[0].Err, dtos[2].Err)
	}
	if dtos[1].Err == nil {
		t.Fatal("dtos[1].Err = nil, want a per-item decode error")
	}
	if dtos[1].ProviderItemID != "MLB2" {
		t.Fatalf("dtos[1].ProviderItemID = %q, want MLB2 (id slot preserved even on decode failure)", dtos[1].ProviderItemID)
	}
}

// TestGetItemsMultigetOversizeItemBodyMarksRawTruncated: a per-item body
// larger than the 1MB per-item raw cap must degrade to a truncated,
// explicitly-marked Raw — never a silent partial prefix.
func TestGetItemsMultigetOversizeItemBodyMarksRawTruncated(t *testing.T) {
	t.Parallel()

	if itemMultigetRawCap != 256*1024 {
		t.Fatalf("itemMultigetRawCap = %d, want 256*1024 (documented per-item cap, see rationale comment on the const)", itemMultigetRawCap)
	}

	// Build a single-item body whose JSON encoding exceeds itemMultigetRawCap
	// (256KiB) via an oversized filler field alongside a normal id/title.
	filler := strings.Repeat("x", itemMultigetRawCap+1024)
	ids := []string{"MLB-BIG"}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `[{"code":200,"body":{"id":"MLB-BIG","title":"big","permalink":%q}}]`, filler)
	}))
	defer server.Close()

	dtos, err := pricingTestAdapter(server.URL, fixedNow()).getItemsMultiget(context.Background(), pricingAccountRef(), "test-token", ids)
	if err != nil {
		t.Fatalf("getItemsMultiget() error = %v", err)
	}
	if len(dtos) != 1 {
		t.Fatalf("dto count = %d, want 1", len(dtos))
	}
	if !dtos[0].RawTruncated {
		t.Fatal("RawTruncated = false, want true for an oversize per-item body")
	}
	if len(dtos[0].Raw) != itemMultigetRawCap {
		t.Fatalf("len(Raw) = %d, want exactly the %d-byte cap", len(dtos[0].Raw), itemMultigetRawCap)
	}
	if dtos[0].Err != nil {
		t.Fatalf("Err = %v, want nil (truncation is not a per-item error)", dtos[0].Err)
	}
}

// TestGetItemsMultigetBatch429Propagates: a whole-batch 429 from the choke
// point now retries (F-01's resilience decorator applies retry by HTTP
// method, not by call site — GET /items?ids= retries the same as every other
// GET reader; see the scope note atop resilience_decorator.go) and, once the
// tightened retry budget below is exhausted, still surfaces as
// ErrRateLimited. F-02 implements no retry of its own — the retry loop lives
// entirely in the shared choke point this reader's doJSON call funnels
// through.
func TestGetItemsMultigetBatch429Propagates(t *testing.T) {
	t.Parallel()

	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		calls++
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer server.Close()

	adapter := NewCapabilityAdapter(CapabilityAdapterConfig{
		BaseURL:    server.URL,
		HTTPClient: &http.Client{Timeout: time.Second},
		AccessTokenResolver: func(context.Context, domain.ProviderAccountRef) (string, error) {
			return "test-token", nil
		},
		Now:               fixedNow,
		MaxRetryAttempts:  2,
		RetryBaseDelay:    1 * time.Millisecond,
		MaxTotalRetryWait: 50 * time.Millisecond,
	})

	_, err := adapter.getItemsMultiget(context.Background(), pricingAccountRef(), "test-token", []string{"MLB1"})
	if !errors.Is(err, domain.ErrRateLimited) {
		t.Fatalf("error = %v, want errors.Is(_, ErrRateLimited)", err)
	}
	if calls != 2 {
		t.Fatalf("calls = %d, want 2 (MaxRetryAttempts, retry-then-exhaust)", calls)
	}
}

// TestGetItemsMultigetWholeBatchUndecodableIsBatchError: a response that is
// not a valid `[]{code,body}` array (e.g. truncated/garbled JSON) is a
// BATCH-level failure — distinct from a per-item error — and must propagate a
// clearly-named error rather than silently dropping ids.
func TestGetItemsMultigetWholeBatchUndecodableIsBatchError(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, `{"not":"an array"}`)
	}))
	defer server.Close()

	_, err := pricingTestAdapter(server.URL, fixedNow()).getItemsMultiget(context.Background(), pricingAccountRef(), "test-token", []string{"MLB1", "MLB2"})
	if err == nil {
		t.Fatal("error = nil, want a batch-level decode error")
	}
	if !errors.Is(err, ErrItemsMultigetBatchDecode) {
		t.Fatalf("error = %v, want errors.Is(_, ErrItemsMultigetBatchDecode)", err)
	}
}

// TestGetItemsMultigetEmptyIDsIsNoOp: an empty id slice makes zero HTTP calls
// and returns an empty result, never an error.
func TestGetItemsMultigetEmptyIDsIsNoOp(t *testing.T) {
	t.Parallel()

	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		calls++
	}))
	defer server.Close()

	dtos, err := pricingTestAdapter(server.URL, fixedNow()).getItemsMultiget(context.Background(), pricingAccountRef(), "test-token", nil)
	if err != nil {
		t.Fatalf("getItemsMultiget() error = %v, want nil", err)
	}
	if len(dtos) != 0 {
		t.Fatalf("dto count = %d, want 0", len(dtos))
	}
	if calls != 0 {
		t.Fatalf("calls = %d, want 0", calls)
	}
}

// TestGetItemsMultigetMapsE3Fields covers the IC-07 E3 field set the DTO must
// carry typed (in addition to Raw) per the feature brief: sold_quantity,
// category_id, condition, permalink, thumbnail, date_created, tags,
// catalog_product_id, shipping, available_quantity, and variations with
// price/available_quantity/sold_quantity/seller_sku/attributes.
func TestGetItemsMultigetMapsE3Fields(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, `[{"code":200,"body":{
			"id":"MLB1",
			"title":"Item E3",
			"status":"active",
			"sold_quantity":37,
			"category_id":"MLB1055",
			"condition":"new",
			"permalink":"https://produto.mercadolivre.com.br/MLB1",
			"thumbnail":"https://http2.mlstatic.com/thumb.jpg",
			"date_created":"2024-04-23T10:48:51Z",
			"tags":["good_quality_thumbnail","immediate_payment"],
			"catalog_product_id":"MLB99887766",
			"available_quantity":5,
			"shipping":{"mode":"me2","logistic_type":"fulfillment","free_shipping":true},
			"variations":[{"id":181,"price":79.9,"available_quantity":3,"sold_quantity":2,"seller_sku":"SKU-A","attributes":[{"id":"COLOR","value_name":"Azul"}]}]
		}}]`)
	}))
	defer server.Close()

	dtos, err := pricingTestAdapter(server.URL, fixedNow()).getItemsMultiget(context.Background(), pricingAccountRef(), "test-token", []string{"MLB1"})
	if err != nil {
		t.Fatalf("getItemsMultiget() error = %v", err)
	}
	if len(dtos) != 1 {
		t.Fatalf("dto count = %d, want 1", len(dtos))
	}
	dto := dtos[0]
	if dto.SoldQuantity == nil || *dto.SoldQuantity != 37 {
		t.Fatalf("SoldQuantity = %#v, want 37", dto.SoldQuantity)
	}
	if dto.CategoryID != "MLB1055" {
		t.Fatalf("CategoryID = %q, want MLB1055", dto.CategoryID)
	}
	if dto.Condition != "new" {
		t.Fatalf("Condition = %q, want new", dto.Condition)
	}
	if dto.Permalink != "https://produto.mercadolivre.com.br/MLB1" {
		t.Fatalf("Permalink = %q", dto.Permalink)
	}
	if dto.Thumbnail != "https://http2.mlstatic.com/thumb.jpg" {
		t.Fatalf("Thumbnail = %q", dto.Thumbnail)
	}
	if dto.DateCreated == nil {
		t.Fatal("DateCreated = nil, want parsed timestamp")
	}
	if len(dto.Tags) != 2 || dto.Tags[0] != "good_quality_thumbnail" {
		t.Fatalf("Tags = %#v", dto.Tags)
	}
	if dto.CatalogProductID != "MLB99887766" {
		t.Fatalf("CatalogProductID = %q", dto.CatalogProductID)
	}
	if dto.AvailableQuantity == nil || *dto.AvailableQuantity != 5 {
		t.Fatalf("AvailableQuantity = %#v, want 5", dto.AvailableQuantity)
	}
	if dto.Shipping == nil || dto.Shipping.Mode != "me2" || dto.Shipping.LogisticType != "fulfillment" || dto.Shipping.FreeShipping == nil || !*dto.Shipping.FreeShipping {
		t.Fatalf("Shipping = %#v", dto.Shipping)
	}
	if len(dto.Variations) != 1 {
		t.Fatalf("Variations count = %d, want 1", len(dto.Variations))
	}
	variation := dto.Variations[0]
	if variation.ProviderVariationID != "181" {
		t.Fatalf("Variation.ProviderVariationID = %q, want 181", variation.ProviderVariationID)
	}
	if variation.Price == nil || *variation.Price != "79.9" {
		t.Fatalf("Variation.Price = %#v, want 79.9", variation.Price)
	}
	if variation.AvailableQuantity == nil || *variation.AvailableQuantity != 3 {
		t.Fatalf("Variation.AvailableQuantity = %#v, want 3", variation.AvailableQuantity)
	}
	if variation.SoldQuantity == nil || *variation.SoldQuantity != 2 {
		t.Fatalf("Variation.SoldQuantity = %#v, want 2", variation.SoldQuantity)
	}
	if variation.SellerSKU != "SKU-A" {
		t.Fatalf("Variation.SellerSKU = %q, want SKU-A", variation.SellerSKU)
	}
	if len(variation.Attributes) != 1 || variation.Attributes[0].ID != "COLOR" || variation.Attributes[0].ValueName != "Azul" {
		t.Fatalf("Variation.Attributes = %#v", variation.Attributes)
	}
}

// TestGetItemsMultigetPublicEntryPointResolvesTokenAndAccount exercises the
// exported GetItemsMultiget wrapper (account normalization + token
// resolution), mirroring the Get*/normalizeAccountRef pattern of every other
// reader in this package (e.g. GetShipmentInfo, GetOwnItemPricing).
func TestGetItemsMultigetPublicEntryPointResolvesTokenAndAccount(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer test-token" {
			t.Fatalf("authorization = %q", got)
		}
		_, _ = io.WriteString(w, `[{"code":200,"body":{"id":"MLB1"}}]`)
	}))
	defer server.Close()

	dtos, err := pricingTestAdapter(server.URL, fixedNow()).GetItemsMultiget(context.Background(), pricingAccountRef(), []string{"MLB1"})
	if err != nil {
		t.Fatalf("GetItemsMultiget() error = %v", err)
	}
	if len(dtos) != 1 || dtos[0].ProviderItemID != "MLB1" {
		t.Fatalf("dtos = %#v", dtos)
	}
}

func indexOf(values []string, target string) int {
	for i, v := range values {
		if v == target {
			return i
		}
	}
	return -1
}

func equalInts(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
