# D5-B2 — W3 Collections / Search / Cursor Grammar

> **Status:** ACCEPTED / CURRENT CONSOLIDATED AUTHORITY  
> **Parent:** `D5-B2-WIRE-CONTRACT.md`  
> **Schema authority:** `D5-B2-W2-SCHEMA-GRAMMAR.md`

## 1. Governing invariant

> **Expose the smallest typed collection/query vocabulary that answers a real Product question. Pagination carries continuation only; collection items remain subsets of owner meaning; traversal never creates knowledge/completeness authority.**

No universal Page/Result/View/projection/filter/sort DSL is admitted.

## 2. Named owner-specific collections

Each List/Search operation owns a named response shape:

```json
{
  "<owner_population>": [],
  "next_cursor": "opaque-if-more"
}
```

Owner-specific coverage/provenance may be added only where that collection legitimately owns it. No universal `data`, `metadata`, `total`, `total_count`, `page`, `offset` or estimated count exists.

## 3. ListItem semantic-subset law

A ListItem is a semantic subset of the same owner meaning needed to scan/select/navigate. It may omit detail/history/large evidence but cannot create a list-only business conclusion.

A point GET fan-out is not a substitute for a materially deficient collection item when the owner can directly provide the necessary current item meaning. Conversely, no generic `Summary<T>`, `View<T>`, `expand` or selection DSL is admitted.

When the admitted human consumer must scan/select/navigate a member and the owner can supply current presentation without a second business conclusion, the ListItem carries that owner-semantic presentation subset directly. Point-GET fan-out is not the baseline repair for a deficient collection item. This admits no generic `View<T>`, projection DSL, total count, alternate sort, or metadata envelope.

The currently proven collection set is exactly:

```text
SearchSourceProductsForMarketplace
ListMarketplaceListings
ListListingIntents
ListPriceIntents
ListSellableAvailability
ListCompetitivePositions
ListExpectedEconomics
ListMarketplaceListingPerformance
```

Any expansion beyond this set requires a separately proven human job and the smallest-owner reopen.

## 4. Shared continuation grammar

Admitted List/Search operations use only:

```text
limit?
cursor?
```

`limit` is a requested maximum, not an exact cardinality promise. `next_cursor` is present only when another traversal page exists. A page may contain fewer than `limit` members, including zero, and still have a continuation.

No baseline:

```text
page
page_number
offset
skip
before
previous_cursor
```

Traversal is forward-only.

## 5. Cursor opacity / explicit query

A cursor is opaque continuation state, never business selection/identity/authorization/completeness evidence. Provider/database tokens are not exposed as Product cursor semantics.

Every continuation request repeats the explicit semantic query/scope that selected the population. The cursor may carry a fingerprint for validation; it does not replace the query.

Processing preserves ordinary auth/access/privacy before treating the cursor as authorization. A validly shaped request whose explicit query disagrees with the cursor is invalid-cursor semantics, not a hidden query override.

Only `limit` may vary across continuation when the operation's semantics otherwise stay identical.

## 6. Pagination != completeness

No `next_cursor` means only no later page in that Product traversal. It does **not** prove provider/source/all-time/portfolio completeness or reconciliation closure.

```text
cursor exhaustion
!= source enumeration completeness
!= owner knowledge completeness
!= deduplication
```

Where coverage is material, the owning collection exposes an explicit coverage/provenance meaning independent of cursor state.

## 7. Sort / filters / query grammar

No generic caller-selectable sort baseline exists. Each collection has one owner-defined deterministic order; a second order requires a real consumer and stable semantics.

Filters are ordinary typed operation-specific parameters. Different admitted fields combine by AND; no generic OR/NOT/IN/expression language, SQL/OData/GraphQL field traversal or provider query syntax is admitted.

Time windows use operation-specific half-open `[from,before)` semantics where admitted; generic `from/to` is not a convention.

## 8. SearchSourceProductsForMarketplace

This remains the only Product Search operation and is **bounded marketplace-operating discovery**, not an MPC Product master/listing API.

Current request requires:

```text
marketplace_installation_id
query                 non-empty
limit?
cursor?
source_instance_id?   optional exact narrowing
```

When `source_instance_id` is omitted, Readiness searches the Organization-scoped SourceInstances currently admitted/configured for this search context. Omission never selects an ambient/default SourceInstance. Every hit remains explicitly source-qualified.

Matching may use legitimate source identification/display evidence supported by sanctioned source capability. Exact native/legitimate identifier matches may rank ahead of textual matches where source contract establishes them. No Product contract is frozen for tokenizer, stemming, fuzzy/vector search or relevance score.

Within one traversal the same `SourceProductRef` is returned at most once and remains a stable member identity/tie-breaker. If sanctioned source capability cannot satisfy materially required search without inventing an MPC Product-search mirror/index authority, STOP / targeted re-adjudication.

## 9. Important current collection populations

### Marketplace Listings

`ListMarketplaceListings` enumerates currently known external Listing subjects within exact Marketplace Installation scope, ordered by source-qualified Listing identity. Listing acquisition coverage is explicit and independent from cursor exhaustion.

### ListingIntents / PriceIntents

Enumerate MPC-owned intents including relevant lifecycle/history under their accepted filters/order. Provider Listing status is not ListingIntent lifecycle.

### SellableAvailability

Enumerates current Availability target universe, not persistence rows:

1. pre-creation ListingIntent targets still genuinely pre-creation/current;
2. currently known existing Listing targets.

Discarded intents are excluded; incomplete draft may remain present with unknown/unavailable Availability. When provider Listing is established, current identity shifts to existing Listing rather than duplicating one current target twice.

### CompetitivePositions / ExpectedEconomics

Enumerable collection baseline is currently known existing Listing subjects; pre-listing source-product analysis remains point-query/evaluation flow. Collection exhaustion never proves complete marketplace portfolio/source-product universe.

### Performance

Listing/media Performance collections remain exact Installation + explicit period/comparison scoped and use owner-local deterministic ordering/cursor. No generic metric/dimension/sort API exists.

## 10. Personal Notifications collections

### 10.1 Self Inbox

`ListMyNotifications` is exact-self H-only current-Membership awareness. Current admitted query controls are:

```text
archive_state?
read_state?
notification_kind[]?   non-empty unique list when supplied
limit?
cursor?
```

No recipient/principal filter, free-text search, generic status, source-type expression or unread total is admitted. Item order is owner-defined current Inbox order; Notification source truth remains elsewhere.

### 10.2 Notification routes

`ListNotificationRoutes` enumerates current Organization ORG_ROUTED configuration meanings. Route history/revision storage exists for temporal correctness under D2/D7 but is not automatically a public history collection.

### 10.3 Route recipient candidates

`ListNotificationRouteRecipientCandidates` uses only `limit?`/`cursor?` and returns bounded eligible Principal candidates needed for routing. It is not Organization-member/role/Permission browse by proxy.

## 11. Actionable AuthorizationRequest collection

`ListMyActionableAuthorizationRequests` is H-only `governance.decide` and enumerates only requests currently actionable by the exact Principal under Governance meaning.

Query is deliberately only:

```text
limit?
cursor?
```

No status, target-kind, requester, assignee, role, free-text or general Governance-history filter is admitted. List item is the W2 bounded scan/select subset; point detail returns the typed actionable review view.

This collection does not imply every pending AuthorizationRequest is enumerable by the caller; eligibility is part of the population definition, not a client filter.

## 12. Provider pagination fence

Product cursor may internally depend on provider token/offset/scroll/search-after/seen-set/checkpoint state, but that remains D4/D7 mechanism. Provider paging changes do not automatically change Product cursor grammar.

## 13. Reopen triggers

Reopen W3 only when a real consumer needs a new enumerable population/filter/order/search distinction that cannot be represented honestly by current operation semantics. Database convenience, screen-table affordance or provider query features alone are not admission evidence.
