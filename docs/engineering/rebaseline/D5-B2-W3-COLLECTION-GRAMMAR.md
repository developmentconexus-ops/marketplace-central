# D5-B2 — W3 Collections / Pagination / Filter / Search / Cursor Grammar

> **Status:** OPEN / ACTIVE — W3-A Collection + Cursor Core **ACCEPTED IN-STAGE / OPERATOR-RATIFIED**; W3-B Per-Operation Filter / Search / Ordering Matrix next  
> **Parent Wire Contract:** `D5-B2-WIRE-CONTRACT.md`  
> **Schema authority:** `D5-B2-W2-SCHEMA-GRAMMAR.md`  
> **Operation inventory:** `D5-B2-OPERATION-ADMISSION-MATRIX.md`  
> **Parent authorities:** accepted D0→D4 + D4-R1 + D5-B1 + Decision Reconciliation Baseline + ratified D5-B2 Whole-Matrix + canonical W1/W2  
> **Method:** DevelopmentConexus Engineering Method v1.0.0  
> **W3-A accepted:** 2026-08-19

## 1. Purpose

W3 defines Product API collection traversal semantics only for the already-admitted List/Search Q operations.

It does not derive collection contracts from legacy routes/OpenAPI/frontend tables and does not choose D7 database/index/cache/provider-paging implementation.

> **Uniform continuation mechanics may be shared; business collection meaning, coverage, filters, search semantics and ordering remain operation/owner-specific. Pagination state never becomes knowledge/completeness authority.**

---

# 2. W3-A — Collection + Cursor Core — ACCEPTED IN-STAGE

## 2.1 Named collection responses; no universal page/result wrapper

Every List/Search operation owns a response schema named for that operation/owner meaning.

Example shapes:

```json
{
  "marketplace_listings": [],
  "next_cursor": "opaque"
}
```

```json
{
  "work": [],
  "next_cursor": "opaque"
}
```

Do not introduce `Page<T>`, `PagedResult<T>`, universal `data`, `metadata`, `Result<T>` or another generic business envelope.

Owner-specific coverage, evidence sufficiency, provenance or other collection semantics appear only where that collection legitimately owns them.

## 2.2 Baseline request grammar

All admitted List/Search operations are pagination-capable using only:

```text
limit?
cursor?
```

at the shared continuation-mechanism layer.

`limit` is a requested maximum, never a promise that exactly that many elements are returned.

Baseline does not admit:

```text
page
page_number
offset
skip
before
previous_cursor
```

Pagination is forward-only.

`pagination-capable` does not mean every small collection must actually issue a continuation cursor in normal use.

## 2.3 `next_cursor` grammar

When another traversal page exists, the response may contain:

```text
next_cursor: opaque non-empty string
```

When no subsequent traversal page exists, `next_cursor` is omitted. `null` does not mean end-of-pagination.

A page may contain fewer items than `limit`, including zero items, and still carry `next_cursor`; item count alone never determines exhaustion.

## 2.4 Cursor opacity

A Product cursor is an opaque continuation token. Clients cannot depend on its encoding or inspect it for ordering/source/provider state.

The implementation may later encode/reference proportionately:

- ordering continuation state;
- provider continuation state;
- server-side continuation identity;
- query fingerprint;
- acquisition/coverage boundary;
- other D7 mechanism state.

None becomes Product API meaning.

Raw provider paging tokens, scroll IDs, database keys or transparent base64-encoded implementation state must not leak as Product contract.

## 2.5 Cursor is not authorization

A previously issued cursor never proves:

- Principal identity;
- Organization Membership;
- Permission;
- current source/provider access;
- business authorization.

Every continuation request is authenticated/authorized normally. D7 may bind cursor integrity to Organization/operation/query as a mechanism, but the cursor never becomes access authority.

## 2.6 Cursor is bound to the semantic query

A cursor continues the same Product collection query for the same operation and Organization scope.

Changing operation, Organization or material collection-selection/filter/search parameters while supplying a cursor is invalid.

The client may repeat the query fields; the server verifies semantic equivalence against the cursor-bound traversal.

`limit` may change between requests because it changes only requested page cardinality, not the selected population.

Exact cursor integrity/signing/storage mechanics remain D7.

## 2.7 Pagination state is not coverage/completeness

Absence of `next_cursor` means only:

> no subsequent page exists in this Product collection traversal.

It does **not** prove:

- universal source/provider completeness;
- all-time population completeness;
- cancellation-inclusive Sales completeness;
- complete market universe;
- owner knowledge beyond that collection's explicit coverage contract.

An owner response may therefore have no `next_cursor` while still declaring `partial`, `unknown` or `unavailable` coverage/evidence state.

```text
cursor exhaustion != knowledge completeness
```

## 2.8 No baseline total count

W3-A admits no universal `total`, `total_count` or estimated-count field.

A later real consumer may justify an operation-specific count only when the counted universe, filters, freshness/coverage and authoritative semantics are explicit and honest.

Database COUNT convenience is not Product-contract evidence.

## 2.9 No caller-selectable sort baseline

W3-A admits no generic:

```text
sort
order_by
sort_by
field,-field
```

Each collection must have deterministic owner-meaning ordering plus required tie-breaker(s), decided operation-by-operation in W3-B.

Caller-selectable ordering is added only to an operation that proves a real consumer need and a stable semantic ordering contract.

## 2.10 No global snapshot-isolation promise

A cursor traversal does not imply one immutable snapshot across pages unless a specific collection later proves and documents that guarantee.

W3 guarantees honest continuation semantics, not a universal transaction snapshot.

When populations can change between calls, collection-specific stability/deduplication expectations belong to W3-B/W3-C as needed; D7 chooses implementation.

## 2.11 Provider pagination remains D4-local mechanism

Product traversal and provider traversal are distinct layers:

```text
Product cursor
    ↓
owner/application collection semantics
    ↓
D4 adapter
    ↓
provider token / offset / scroll / search-after / provider-specific paging
```

Provider pagination changes do not automatically alter Product cursor grammar.

## 2.12 Error behavior

Malformed, invalid, expired or query-mismatched cursor never silently becomes an empty/first page.

Exact Problem Details types/status codes are finalized within W3 once cursor lifetime/staleness classes are closed; the fail-explicit invariant is already binding.

## 2.13 W3-A negative controls

Later contract proof must falsify at least:

1. raw provider cursor/scroll token exposed directly;
2. cursor accepted under another Organization or operation;
3. cursor reused with materially changed filter/search parameters;
4. cursor treated as authorization;
5. missing `next_cursor` interpreted as universal source/provider completeness;
6. `len(items) < limit` interpreted as end-of-traversal;
7. `total_count` added merely because storage can count;
8. offset/page-number pagination made baseline by familiarity;
9. arbitrary caller sort expression admitted by symmetry;
10. invalid/expired cursor silently restarting or returning empty data;
11. one universal snapshot-isolation promise applied to provider/source-backed collections;
12. universal `Page<T>`/`metadata` wrapper becoming business schema authority.

## 2.14 W3-A outcome

**Outcome:** `CURRENT PARENT STRUCTURE CONFIRMED`.

> **Use owner-named collection schemas plus shared `limit` / opaque `cursor` / optional `next_cursor` continuation mechanics. Keep coverage, search, filtering and ordering owner-specific; cursor exhaustion never implies knowledge completeness.**

No D0→D5-B1/W1/W2 reopen is required.

---

# 3. Exact next W3 work

**W3-B — Per-Operation Filter / Search / Ordering Matrix.**

For every admitted List/Search Q, W3-B must decide from real consumer/owner meaning:

1. whether filtering is needed at all;
2. the smallest typed filter fields, with no generic query DSL;
3. whether operation is true search versus list/filter (`SearchSourceProductsForMarketplace` especially);
4. deterministic default semantic ordering and stable tie-breaker;
5. whether any caller-selectable sort is genuinely justified;
6. how source-qualified identity/same-Organization constraints participate in query selection;
7. collection-specific coverage/provenance exposure where material;
8. whether very small product-defined collections can simply return all values while remaining pagination-capable;
9. negative controls preventing provider/database field exposure or frontend-table-shaped filter sets.

After W3-B, close cursor invalid/expired/stale and population-change/deduplication semantics as required by the per-operation matrix before W3 can be accepted as a whole.

Implementation remains blocked until D9.
