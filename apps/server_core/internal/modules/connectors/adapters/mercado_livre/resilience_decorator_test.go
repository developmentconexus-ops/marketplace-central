package mercadolivre

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/modules/connectors/domain"
)

// F-01-resilience-decorator (MIS-007 / M-01) tests.
//
// Retry policy is method-based at the shared doJSON/doJSONWithHeaders choke
// point (doRawWithHeaders): every GET retries on HTTP 429 through
// doRetryable, non-GET (writes) never auto-retry — see the scope note atop
// resilience_decorator.go. The token-bucket rate limiter applies to every
// outbound call regardless of retry policy, so the concurrency test below
// exercises it through ProbeAccount (a plain GET, but one with no 429 in
// play) rather than ReadFeeQuote.

func resilienceAccountRef(installationID string) domain.ProviderAccountRef {
	return domain.ProviderAccountRef{
		TenantID:          "tenant-resilience",
		InstallationID:    installationID,
		ProviderCode:      "mercado_livre",
		ProviderAccountID: "seller-resilience",
	}
}

func resilienceFeeQuoteInput(installationID string) domain.FeeQuoteInput {
	return domain.FeeQuoteInput{
		AccountRef:    resilienceAccountRef(installationID),
		SiteID:        "MLB",
		ListingTypeID: "gold_special",
		PriceAmount:   100,
		CurrencyID:    "BRL",
	}
}

func writeListingPriceOK(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`[{"listing_type_id":"gold_special","sale_fee_details":{"percentage_fee":11,"fixed_fee":6}}]`))
}

// TestResilienceDecoratorHonorsRetryAfterHeader is the money assertion for
// requirement 1: the fake provider answers 429 with Retry-After: 2 on the
// first call and 200 on the second, and the test measures real elapsed time
// between the two attempts. An "eventually succeeds" assertion with no
// elapsed-time check would be a vacuous pass; this asserts elapsed >= 2s.
//
// Must-fail check performed manually during development: commenting out the
// `d.sleep(ctx, wait)` call in resilienceDecorator.doRetryable (so the retry
// fires immediately) makes this test fail with "elapsed < retry-after"
// (measured ~0s instead of >=2s). That mutation was not committed.
func TestResilienceDecoratorHonorsRetryAfterHeader(t *testing.T) {
	t.Parallel()

	var calls int32
	var firstCallAt, secondCallAt time.Time
	var mu sync.Mutex

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&calls, 1)
		mu.Lock()
		if n == 1 {
			firstCallAt = time.Now()
		} else {
			secondCallAt = time.Now()
		}
		mu.Unlock()
		if n == 1 {
			w.Header().Set("Retry-After", "2")
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		writeListingPriceOK(w)
	}))
	defer server.Close()

	adapter := NewCapabilityAdapter(CapabilityAdapterConfig{
		BaseURL:    server.URL,
		SiteID:     "MLB",
		HTTPClient: server.Client(),
		AccessTokenResolver: func(context.Context, domain.ProviderAccountRef) (string, error) {
			return "token-retry-after", nil
		},
	})

	t0 := time.Now()
	_, err := adapter.ReadFeeQuote(context.Background(), resilienceFeeQuoteInput("inst-retry-after"))
	elapsed := time.Since(t0)
	if err != nil {
		t.Fatalf("ReadFeeQuote() error = %v", err)
	}
	if atomic.LoadInt32(&calls) != 2 {
		t.Fatalf("provider calls = %d, want 2", calls)
	}
	if elapsed < 2*time.Second {
		t.Fatalf("elapsed = %s, want >= 2s (Retry-After was not honored)", elapsed)
	}
	mu.Lock()
	interCallGap := secondCallAt.Sub(firstCallAt)
	mu.Unlock()
	if interCallGap < 2*time.Second {
		t.Fatalf("inter-call gap = %s, want >= 2s", interCallGap)
	}
	t.Logf("measured total elapsed = %s, inter-call gap = %s", elapsed, interCallGap)
}

// TestResilienceDecoratorTokenBucketThrottlesConcurrentRequests fires N
// goroutines concurrently at the same installation and asserts on the
// OBSERVED request timestamps (not on the configured rate) that consecutive
// requests are spaced by at least the configured interval — i.e. the bucket
// actually queues callers instead of letting them burst through.
func TestResilienceDecoratorTokenBucketThrottlesConcurrentRequests(t *testing.T) {
	t.Parallel()

	const ratePerMinute = 600 // 1 request per 100ms
	const n = 5
	interval := time.Minute / time.Duration(ratePerMinute)

	var mu sync.Mutex
	var timestamps []time.Time

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		timestamps = append(timestamps, time.Now())
		mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":"acct-1","nickname":"seller","site_id":"MLB","status":{"site_status":"active"}}`))
	}))
	defer server.Close()

	adapter := NewCapabilityAdapter(CapabilityAdapterConfig{
		BaseURL:    server.URL,
		HTTPClient: server.Client(),
		AccessTokenResolver: func(context.Context, domain.ProviderAccountRef) (string, error) {
			return "token-bucket", nil
		},
		RateLimitPerMinute: ratePerMinute,
	})

	accountRef := resilienceAccountRef("inst-bucket")

	var wg sync.WaitGroup
	t0 := time.Now()
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := adapter.ProbeAccount(context.Background(), accountRef); err != nil {
				t.Errorf("ProbeAccount() error = %v", err)
			}
		}()
	}
	wg.Wait()
	totalElapsed := time.Since(t0)

	mu.Lock()
	defer mu.Unlock()
	if len(timestamps) != n {
		t.Fatalf("recorded timestamps = %d, want %d", len(timestamps), n)
	}
	sortTimes(timestamps)

	// What the bucket guarantees is release spacing, and therefore total
	// throughput -- NOT the spacing of server-side receipt timestamps. A
	// goroutine released at t=0 can be descheduled for 95ms and arrive 5ms
	// before the one released at t=100ms, so a per-pair gap assertion on
	// receipt times flakes under whole-suite load no matter the tolerance
	// (observed: a 5.4ms receipt gap from a correctly-queued release).
	// Total elapsed is robust against reordering: 5 queued requests take at
	// least (n-1)*interval end to end, a burst finishes in milliseconds.
	for i := 1; i < len(timestamps); i++ {
		gap := timestamps[i].Sub(timestamps[i-1])
		t.Logf("gap[%d] = %s (timestamps[%d]=%s timestamps[%d]=%s)", i, gap, i-1, timestamps[i-1].Format(time.RFC3339Nano), i, timestamps[i].Format(time.RFC3339Nano))
	}
	const tolerance = 15 * time.Millisecond
	minTotal := time.Duration(n-1) * interval
	if totalElapsed < minTotal-tolerance {
		t.Fatalf("total elapsed = %s, want >= ~%s for %d requests at %d/min (bucket allowed a burst)", totalElapsed, minTotal, n, ratePerMinute)
	}
	t.Logf("total elapsed for %d concurrent requests at %d/min = %s (min expected ~%s)", n, ratePerMinute, totalElapsed, minTotal)
}

func sortTimes(ts []time.Time) {
	for i := 1; i < len(ts); i++ {
		for j := i; j > 0 && ts[j].Before(ts[j-1]); j-- {
			ts[j], ts[j-1] = ts[j-1], ts[j]
		}
	}
}

// TestResilienceDecoratorBudgetExhaustedReturnsTypedError sets a small budget
// (via CapabilityAdapterConfig) so the test runs fast in real (not mocked)
// time while still exercising the real sleep path. The fake provider answers
// 429 with no Retry-After on every call, forcing exponential backoff on every
// attempt until the budget is exhausted.
func TestResilienceDecoratorBudgetExhaustedReturnsTypedError(t *testing.T) {
	t.Parallel()

	var calls int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&calls, 1)
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer server.Close()

	adapter := NewCapabilityAdapter(CapabilityAdapterConfig{
		BaseURL:    server.URL,
		SiteID:     "MLB",
		HTTPClient: server.Client(),
		AccessTokenResolver: func(context.Context, domain.ProviderAccountRef) (string, error) {
			return "token-budget", nil
		},
		MaxRetryAttempts:  3,
		MaxTotalRetryWait: 200 * time.Millisecond,
		RetryBaseDelay:    2 * time.Millisecond,
	})

	_, err := adapter.ReadFeeQuote(context.Background(), resilienceFeeQuoteInput("inst-budget"))
	if err == nil {
		t.Fatal("ReadFeeQuote() error = nil, want rate limit exhaustion error")
	}
	if got := domain.ErrorCodeOf(err); got != domain.ErrCodeProviderRateLimited {
		t.Fatalf("error code = %q, err = %v", got, err)
	}
	if !strings.Contains(err.Error(), "3 attempts") {
		t.Fatalf("error message = %q, want it to name the attempt count (3 attempts)", err.Error())
	}
	var exhausted *RateLimitExhaustedError
	if !errors.As(err, &exhausted) {
		t.Fatalf("errors.As(err, *RateLimitExhaustedError) = false, err = %v (%#v)", err, err)
	}
	if exhausted.Attempts != 3 {
		t.Fatalf("exhausted.Attempts = %d, want 3", exhausted.Attempts)
	}
	if exhausted.LastRetryAfter <= 0 {
		t.Fatalf("exhausted.LastRetryAfter = %s, want > 0", exhausted.LastRetryAfter)
	}
	if !strings.Contains(err.Error(), exhausted.LastRetryAfter.String()) {
		t.Fatalf("error message = %q, want it to contain the last retry-after %s", err.Error(), exhausted.LastRetryAfter)
	}
	if atomic.LoadInt32(&calls) != 3 {
		t.Fatalf("provider calls = %d, want exactly 3 (maxAttempts)", calls)
	}
}

// TestResilienceDecoratorBackoffWithoutRetryAfterHeader asserts that a 429
// with no (or non-numeric) Retry-After header does not retry instantly: the
// second attempt only fires after a real, positive backoff wait, and the
// decorator never panics parsing the malformed header.
func TestResilienceDecoratorBackoffWithoutRetryAfterHeader(t *testing.T) {
	t.Parallel()

	var calls int32
	var firstCallAt, secondCallAt time.Time
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&calls, 1)
		if n == 1 {
			firstCallAt = time.Now()
			// Non-numeric Retry-After must be treated as absent, never panic.
			w.Header().Set("Retry-After", "not-a-number")
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		secondCallAt = time.Now()
		writeListingPriceOK(w)
	}))
	defer server.Close()

	adapter := NewCapabilityAdapter(CapabilityAdapterConfig{
		BaseURL:    server.URL,
		SiteID:     "MLB",
		HTTPClient: server.Client(),
		AccessTokenResolver: func(context.Context, domain.ProviderAccountRef) (string, error) {
			return "token-backoff", nil
		},
		RetryBaseDelay: 20 * time.Millisecond,
	})

	quote, err := adapter.ReadFeeQuote(context.Background(), resilienceFeeQuoteInput("inst-backoff"))
	if err != nil {
		t.Fatalf("ReadFeeQuote() error = %v (must not panic or fail on malformed Retry-After)", err)
	}
	if quote.CommissionPercent == nil || *quote.CommissionPercent != 11 {
		t.Fatalf("quote = %#v, want successful decode after backoff+retry", quote)
	}
	if atomic.LoadInt32(&calls) != 2 {
		t.Fatalf("provider calls = %d, want 2", calls)
	}
	gap := secondCallAt.Sub(firstCallAt)
	if gap <= 0 {
		t.Fatalf("inter-attempt gap = %s, want > 0 (must not retry instantly)", gap)
	}
	if gap < 20*time.Millisecond {
		t.Fatalf("inter-attempt gap = %s, want >= configured base delay 20ms", gap)
	}
	t.Logf("measured backoff gap without Retry-After header = %s", gap)
}

// TestResilienceDecoratorStockWriteDoesNotRetry pins requirement 4: the PUT
// /items stock write must make exactly one attempt on 429 and surface the
// error immediately (fast — no wait), even though the same rate limiter
// applies to it. This exercises the same scenario as the pre-existing
// TestCapabilityAdapterUpdateAvailableQuantityMapsRateLimit but adds the
// explicit call-count and elapsed-time assertions this feature's contract
// requires.
func TestResilienceDecoratorStockWriteDoesNotRetry(t *testing.T) {
	t.Parallel()

	var calls int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			_ = json.NewEncoder(w).Encode(map[string]any{
				"id":                 "MLB-STOCK-NORETRY",
				"available_quantity": 1,
			})
		case http.MethodPut:
			atomic.AddInt32(&calls, 1)
			w.WriteHeader(http.StatusTooManyRequests)
		default:
			t.Fatalf("unexpected method %s", r.Method)
		}
	}))
	defer server.Close()

	adapter := NewCapabilityAdapter(CapabilityAdapterConfig{
		BaseURL:    server.URL,
		HTTPClient: server.Client(),
		AccessTokenResolver: func(context.Context, domain.ProviderAccountRef) (string, error) {
			return "token-stock-noretry", nil
		},
		// A base delay far above local-HTTP overhead: if the write path
		// incorrectly retried, elapsed necessarily crosses 2s; if it did not,
		// two local round trips stay well under it even on a loaded machine.
		RetryBaseDelay: 2 * time.Second,
	})

	t0 := time.Now()
	_, err := adapter.UpdateAvailableQuantity(context.Background(), stockWriteRequest("MLB-STOCK-NORETRY", "", "action-noretry", 5))
	elapsed := time.Since(t0)

	if got := domain.ErrorCodeOf(err); got != domain.ErrCodeProviderRateLimited {
		t.Fatalf("error code = %q, err = %v", got, err)
	}
	if atomic.LoadInt32(&calls) != 1 {
		t.Fatalf("provider PUT calls = %d, want exactly 1 (no retry on write)", calls)
	}
	if elapsed >= 2*time.Second {
		t.Fatalf("elapsed = %s, want under the 2s base delay (write must not enter the retry/backoff loop)", elapsed)
	}
	t.Logf("stock write 429 propagated after %s with %d attempt(s)", elapsed, calls)
}

// TestResilienceDecoratorContextCancellationDuringWait cancels the context
// shortly after the first 429 (which carries a long Retry-After) and asserts
// the call returns ctx.Err() promptly instead of holding the wait for the
// full Retry-After duration.
func TestResilienceDecoratorContextCancellationDuringWait(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 5s, comfortably inside the adapter's default 30s max-total-wait
		// budget even with jitter, so the decorator actually enters the sleep
		// instead of short-circuiting into budget exhaustion.
		w.Header().Set("Retry-After", "5")
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer server.Close()

	adapter := NewCapabilityAdapter(CapabilityAdapterConfig{
		BaseURL:    server.URL,
		SiteID:     "MLB",
		HTTPClient: server.Client(),
		AccessTokenResolver: func(context.Context, domain.ProviderAccountRef) (string, error) {
			return "token-cancel", nil
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	t0 := time.Now()
	_, err := adapter.ReadFeeQuote(ctx, resilienceFeeQuoteInput("inst-cancel"))
	elapsed := time.Since(t0)

	if !errors.Is(err, context.Canceled) {
		t.Fatalf("error = %v, want context.Canceled", err)
	}
	if elapsed >= 2*time.Second {
		t.Fatalf("elapsed = %s, want prompt cancellation (well under the 5s Retry-After)", elapsed)
	}
	t.Logf("context cancellation during 5s Retry-After wait returned after %s", elapsed)
}

// TestResilienceDecoratorRateLimitPerMinuteDefaultsWhenUnset documents the
// "never a compiled constant" requirement's default: an unset
// RateLimitPerMinute must not block a single request (it falls back to a
// generous default), and must remain configurable (proven by the token-bucket
// test above using a small explicit value).
func TestResilienceDecoratorRateLimitPerMinuteDefaultsWhenUnset(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":"acct-1","nickname":"seller","site_id":"MLB","status":{"site_status":"active"}}`))
	}))
	defer server.Close()

	adapter := NewCapabilityAdapter(CapabilityAdapterConfig{
		BaseURL:    server.URL,
		HTTPClient: server.Client(),
		AccessTokenResolver: func(context.Context, domain.ProviderAccountRef) (string, error) {
			return "token-default-rate", nil
		},
		// RateLimitPerMinute intentionally left zero.
	})

	t0 := time.Now()
	if _, err := adapter.ProbeAccount(context.Background(), resilienceAccountRef("inst-default-rate")); err != nil {
		t.Fatalf("ProbeAccount() error = %v", err)
	}
	// The bound refutes a blocking wait (a broken default would sleep for a
	// rate-limit interval, seconds), not scheduler noise: a single local HTTP
	// round trip has been measured past 200ms under whole-suite load.
	if elapsed := time.Since(t0); elapsed > 2*time.Second {
		t.Fatalf("elapsed = %s, want no blocking wait on the first call under the default rate limit", elapsed)
	}
}
