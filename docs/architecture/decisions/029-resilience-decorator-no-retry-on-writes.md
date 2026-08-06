# ADR-029: Resilience decorator retries GET, opt-out no-retry for writes

**Date:** 2026-08-05
**Status:** accepted
**Reconstructed:** this decision governed the Mercado Livre adapter from MIS-007's
M-01 slice onward but was only ever cited by its local two-digit label `ADR-02`, a
number that collides with three unrelated decisions from other missions and with a
pre-existing three-digit document about Postgres schema topology. It is reconstructed
here from the 3 live-code citations harvested at
`docs/architecture/decisions/_citations/adr-02-twodigit-citations.md`, Assertion A1.
The renumbering to ADR-029 is recorded in
`docs/architecture/decisions/_citations/RENUMBERING-REGISTRY.md:88`.

## Context

The Mercado Livre `CapabilityAdapter` funnels every outbound HTTP call through one
choke point. The provider rate-limits with HTTP 429 and a `Retry-After` header, and
the adapter needs a uniform way to back off and retry rather than making every reader
and writer file reimplement its own retry loop. But a retry is only safe to apply
automatically when replaying the request does not risk a duplicate side effect. A
`GET` is safe to repeat: it reads, nothing changes if it runs twice. A `PUT` or `POST`
against Mercado Livre — updating price, stock quantity, or a listing — is not safe to
repeat without an idempotency key binding the retry to the original attempt, and the
adapter does not have that binding for writes yet.

## Decision

**The resilience decorator retries GET requests automatically on HTTP 429 (honoring
`Retry-After`, else exponential backoff with jitter, bounded by a max-attempt and
max-total-wait budget). Non-GET requests — writes — opt out of that retry loop and
make exactly one attempt.**

**§1 — The opt-out is by HTTP method, not by call site.** `doRawWithHeaders` is the
shared entry point used by every reader and writer in the package. It branches on
method: non-GET is routed to `doRawWithHeadersNoRetry`, GET is routed through the
resilience decorator's `doRetryable`.
> `apps/server_core/internal/modules/connectors/adapters/mercado_livre/capability_adapter.go:753-759`
> — "if method != http.MethodGet { return a.doRawWithHeadersNoRetry(...) } return
> a.resilience.doRetryable(...)"

**§2 — The reason is stated as a rule, not a preference.** Writes do not retry because
idempotent replay safety is not assumed for them yet.
> `apps/server_core/internal/modules/connectors/adapters/mercado_livre/capability_adapter.go:747-749`
> — "non-GET requests (writes) never retry automatically, because idempotent replay
> safety isn't assumed for them yet (ADR-02: \"opt-out no-retry para writes\" — a
> future milestone may add idempotency-keyed write retry, this one does not)."

**§3 — `doRawWithHeadersNoRetry` still consumes a rate-limit token.** The opt-out is
from the retry loop only, not from the shared per-installation token bucket: every
attempt, retried or not, GET or write, consumes a token before it reaches the
provider.
> `apps/server_core/internal/modules/connectors/adapters/mercado_livre/resilience_decorator.go:36`
> — "Non-GET requests are routed to doRawWithHeadersNoRetry and never auto-retry,
> since idempotent replay safety isn't assumed for writes yet (ADR-02: \"opt-out
> no-retry para writes\")." (same paragraph states the rate limiter "applies
> universally regardless of method or retry outcome")

At the time this rule shipped, the concrete write call sites were `price_writer.go`,
`listing_writer.go`, and the stock-quantity `PUT` in `capability_adapter.go`
(`UpdateAvailableQuantity`), all routed through `doRawWithHeadersNoRetry`.

Two pre-existing regression tests were updated in the same slice to match the new
GET-retries behavior: `pricing_reader_test.go`'s
`TestPriceToWinErrorMappingWithoutRetry` and `items_multiget_reader_test.go`'s
`TestGetItemsMultigetBatch429Propagates`, both of which previously asserted
single-attempt behavior on HTTP 429 for GET call sites (under the prior,
call-site-opt-in design) and now assert retry-then-exhaust — more than one call
attempt, the same sentinel error surfaced once the retry budget is spent. Writer
tests (`price_writer_test.go`'s rate-limit subtest, `capability_adapter`'s
stock-write test) were left asserting exactly one call, since a `PUT` never enters
the retry branch.
> `apps/server_core/internal/modules/connectors/adapters/mercado_livre/resilience_decorator.go:36-46`

## Rationale

A retry that replays a write can double-apply it: two price updates, two stock
decrements, two listing mutations from what the caller believed was one call. A
retry that replays a read costs nothing beyond the extra round trip. Splitting the
policy by HTTP method, at the single choke point every call funnels through, means
the split cannot be bypassed by a call site forgetting to opt in — the earlier
design (call-site opt-in) is exactly the version this rule replaced.

## Consequences

- A `429` on a write surfaces to the caller after exactly one attempt; retrying a
  write, if ever added, requires an idempotency key binding first (the code comment
  names this as a possible future milestone, not committed work).
- A `429` on a GET is invisible to the caller as long as the retry budget is not
  exhausted; once exhausted, `RateLimitExhaustedError` (attempts + last
  `Retry-After`) is surfaced.
- If a POST read is ever added to this package, the GET/non-GET branch stops being
  exact and the routing in `doRawWithHeaders` needs revisiting — the code comment
  flags this explicitly.

## Alternatives Considered

**Retry every method uniformly.** Rejected: doing so risks double-applying writes
during a provider rate-limit episode, and Mercado Livre's write endpoints (price,
stock, listing) were not being called with an idempotency key that would make a
replay detectable and safe to collapse.

**Call-site opt-in to retry (the prior design).** Superseded: retry behavior lived
per file rather than at the one choke point, so a new writer could retry by omission
if it forgot to opt out. Moving the branch into the shared `doRawWithHeaders` removes
that failure mode.

## Unverified claims

None. All three live-code anchors in the harvest were read and confirmed to state
the clauses above; no clause beyond what the code comments assert was added.
