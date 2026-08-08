package mercadolivre

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"marketplace-central/apps/server_core/internal/modules/connectors/domain"
)

// F-01-resilience-decorator (MIS-007 / M-01).
//
// This file wraps the adapter's single HTTP choke point (doRawOnce, the former
// body of doRawWithHeaders) with two independent concerns:
//
//  1. A per-installation token-bucket rate limiter (resilienceDecorator.bucketFor)
//     that every outbound call funnels through, regardless of retry policy.
//  2. An opt-in retry-with-backoff loop for HTTP 429 (resilienceDecorator.doRetryable),
//     bounded by a max-attempts / max-total-wait budget.
//
// Scope note (read before extending): doRawWithHeaders — the shared choke point
// consumed by doJSON/doJSONWithHeaders and, transitively, by every reader and
// writer file in this package — applies retry policy by HTTP method, not by
// call site. GET requests (every reader: pricing_reader.go,
// catalog_identity_reader.go, catalog_match_reader.go, catalog_offers_reader.go,
// buyer_fiscal_reader.go, shipping_reader.go, items_multiget_reader.go,
// prices_reader.go) retry on HTTP 429 through this decorator's doRetryable.
// Non-GET requests (price_writer.go, listing_writer.go, and the stock-quantity
// PUT in capability_adapter.go) are routed to doRawWithHeadersNoRetry and never
// auto-retry, since idempotent replay safety isn't assumed for writes yet
// (ADR-029: "opt-out no-retry para writes"). Two previously-existing regression
// tests locked single-attempt behavior on HTTP 429 through GET call sites
// under the old (call-site-opt-in) design — pricing_reader_test.go's
// TestPriceToWinErrorMappingWithoutRetry and
// items_multiget_reader_test.go's TestGetItemsMultigetBatch429Propagates —
// and have been updated in this same slice to assert retry-then-exhaust
// (calls > 1, same sentinel error) instead of single-attempt. Writer tests
// (price_writer_test.go's rate_limit subtest, capability_adapter's stock-write
// test) are unaffected and remain calls == 1, since PUT never enters the retry
// branch. The rate limiter applies universally regardless of method or retry
// outcome (every attempt, retried or not, consumes a token).
const (
	defaultRateLimitPerMinute = 900
	defaultMaxRetryAttempts   = 5
	defaultMaxTotalRetryWait  = 30 * time.Second
	defaultRetryBaseDelay     = 250 * time.Millisecond
	defaultJitterMax          = 250 * time.Millisecond
)

// RateLimitExhaustedError is the Cause carried by the domain.CapabilityError
// returned when the resilience decorator's retry budget is exhausted while the
// provider keeps answering HTTP 429. It is exported so callers (and tests) can
// recover the attempt count / last observed Retry-After via errors.As without
// parsing the error message.
type RateLimitExhaustedError struct {
	Attempts       int
	LastRetryAfter time.Duration
}

func (e *RateLimitExhaustedError) Error() string {
	return fmt.Sprintf("resilience decorator: rate limit retry budget exhausted after %d attempts (last retry-after %s)", e.Attempts, e.LastRetryAfter)
}

// ctxSleep sleeps for d, returning ctx.Err() immediately if ctx is cancelled
// (or already done) before or during the wait instead of holding the caller
// (and any rate-limit token it queued for) for the full duration.
func ctxSleep(ctx context.Context, d time.Duration) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if d <= 0 {
		return nil
	}
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-timer.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// tokenBucket is a minimum-interval limiter (a token bucket with capacity 1):
// each call to Take reserves the next available slot under mutex, then sleeps
// outside the lock (ctx-aware) so concurrent waiters queue without serializing
// on each other's sleep. This guarantees combined throughput never exceeds
// ratePerMinute for the installation the bucket belongs to.
type tokenBucket struct {
	mu          sync.Mutex
	interval    time.Duration
	nextAllowed time.Time
	now         func() time.Time
	sleep       func(context.Context, time.Duration) error
}

func newTokenBucket(ratePerMinute int, now func() time.Time, sleep func(context.Context, time.Duration) error) *tokenBucket {
	if ratePerMinute <= 0 {
		ratePerMinute = defaultRateLimitPerMinute
	}
	return &tokenBucket{
		interval: time.Minute / time.Duration(ratePerMinute),
		now:      now,
		sleep:    sleep,
	}
}

func (b *tokenBucket) Take(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	b.mu.Lock()
	now := b.now()
	base := b.nextAllowed
	if base.Before(now) {
		base = now
	}
	wait := base.Sub(now)
	b.nextAllowed = base.Add(b.interval)
	b.mu.Unlock()

	if wait <= 0 {
		return nil
	}
	return b.sleep(ctx, wait)
}

// resilienceDecorator owns the per-installation rate limiter map and the
// retry/backoff budget for the adapter instance's lifetime.
type resilienceDecorator struct {
	mu                 sync.Mutex
	buckets            map[string]*tokenBucket
	rateLimitPerMinute int
	maxAttempts        int
	maxTotalWait       time.Duration
	baseDelay          time.Duration
	now                func() time.Time
	sleep              func(context.Context, time.Duration) error
	jitter             func(base time.Duration) time.Duration
}

func newResilienceDecorator(rateLimitPerMinute int, maxAttempts int, maxTotalWait time.Duration, baseDelay time.Duration) *resilienceDecorator {
	if rateLimitPerMinute <= 0 {
		rateLimitPerMinute = defaultRateLimitPerMinute
	}
	if maxAttempts <= 0 {
		maxAttempts = defaultMaxRetryAttempts
	}
	if maxTotalWait <= 0 {
		maxTotalWait = defaultMaxTotalRetryWait
	}
	if baseDelay <= 0 {
		baseDelay = defaultRetryBaseDelay
	}
	return &resilienceDecorator{
		buckets:            make(map[string]*tokenBucket),
		rateLimitPerMinute: rateLimitPerMinute,
		maxAttempts:        maxAttempts,
		maxTotalWait:       maxTotalWait,
		baseDelay:          baseDelay,
		now:                time.Now,
		sleep:              ctxSleep,
		jitter:             defaultJitter,
	}
}

// defaultJitter returns extra positive delay proportional to base (capped at
// 25% of it) rather than a fixed absolute amount — a fixed jitter large
// enough to matter for multi-second Retry-After waits would otherwise dwarf a
// millisecond-scale exponential-backoff base (and, worse, could tip a
// caller's configured max-total-wait budget into premature exhaustion). Wait
// is always base + jitter, never less than base.
func defaultJitter(base time.Duration) time.Duration {
	if base <= 0 {
		return 0
	}
	max := base / 4
	if max <= 0 {
		return 0
	}
	return time.Duration(rand.Int63n(int64(max) + 1))
}

func (d *resilienceDecorator) bucketFor(installationID string) *tokenBucket {
	d.mu.Lock()
	defer d.mu.Unlock()
	b, ok := d.buckets[installationID]
	if !ok {
		b = newTokenBucket(d.rateLimitPerMinute, d.now, d.sleep)
		d.buckets[installationID] = b
	}
	return b
}

// acquireToken blocks (queues) until a rate-limit token is available for the
// given installation, or returns ctx.Err() immediately if ctx is cancelled
// while waiting.
func (d *resilienceDecorator) acquireToken(ctx context.Context, installationID string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return d.bucketFor(installationID).Take(ctx)
}

// parseRetryAfterSeconds parses the provider's Retry-After header, which ML
// sends as an integer count of seconds. A missing or non-numeric header is
// treated as absent (never panics) so the caller falls back to exponential
// backoff.
func parseRetryAfterSeconds(header string) (int, bool) {
	header = strings.TrimSpace(header)
	if header == "" {
		return 0, false
	}
	seconds, err := strconv.Atoi(header)
	if err != nil || seconds < 0 {
		return 0, false
	}
	return seconds, true
}

// computeWait returns how long to wait before the next attempt: at least
// Retry-After (plus jitter) when the header is present and numeric, otherwise
// exponential backoff from baseDelay (plus jitter).
func (d *resilienceDecorator) computeWait(retryAfterHeader string, attempt int) time.Duration {
	if seconds, ok := parseRetryAfterSeconds(retryAfterHeader); ok {
		base := time.Duration(seconds) * time.Second
		return base + d.jitter(base)
	}
	shift := attempt - 1
	if shift < 0 {
		shift = 0
	}
	if shift > 20 { // guard against overflow on pathological attempt counts
		shift = 20
	}
	backoff := d.baseDelay * time.Duration(int64(1)<<uint(shift))
	return backoff + d.jitter(backoff)
}

// rawCall is the shape of a single HTTP round trip through doRawOnce, closed
// over the request-specific arguments so doRetryable can invoke it repeatedly.
type rawCall func(context.Context) (*http.Response, []byte, error)

// doRetryable acquires a rate-limit token and performs call; on HTTP 429 it
// waits (Retry-After if present, else exponential backoff+jitter) and retries
// the same request until it succeeds, the provider stops answering 429, or the
// retry budget (max attempts / max total wait) is exhausted. Exhaustion
// returns a *domain.CapabilityError{Code: ErrCodeProviderRateLimited} whose
// Cause is a *RateLimitExhaustedError naming the attempt count and the last
// observed Retry-After, inspectable via errors.As. Context cancellation during
// any wait (token acquire or backoff) returns ctx.Err() immediately.
func (d *resilienceDecorator) doRetryable(ctx context.Context, installationID string, call rawCall) (*http.Response, []byte, error) {
	var totalWait time.Duration
	var lastRetryAfter time.Duration
	attempts := 0
	for {
		attempts++
		if err := d.acquireToken(ctx, installationID); err != nil {
			return nil, nil, err
		}
		resp, rawBody, err := call(ctx)
		if err != nil {
			return nil, nil, err
		}
		if resp.StatusCode != http.StatusTooManyRequests {
			return resp, rawBody, nil
		}

		wait := d.computeWait(resp.Header.Get("Retry-After"), attempts)
		lastRetryAfter = wait
		if attempts >= d.maxAttempts || totalWait+wait > d.maxTotalWait {
			return nil, nil, &domain.CapabilityError{
				Code:    domain.ErrCodeProviderRateLimited,
				Message: fmt.Sprintf("provider rate limit retry budget exhausted after %d attempts (last retry-after %s)", attempts, lastRetryAfter),
				Cause:   &RateLimitExhaustedError{Attempts: attempts, LastRetryAfter: lastRetryAfter},
			}
		}
		if err := d.sleep(ctx, wait); err != nil {
			return nil, nil, err
		}
		totalWait += wait
	}
}
