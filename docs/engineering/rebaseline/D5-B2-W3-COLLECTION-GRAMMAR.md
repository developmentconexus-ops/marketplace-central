# D5-B2 — W3 Collections / Pagination / Filter / Search / Cursor Grammar

> **Status:** OPEN / ACTIVE — W3-A Collection + Cursor Core + W3-B Per-Operation Filter / Search / Ordering Matrix **ACCEPTED IN-STAGE / OPERATOR-RATIFIED**; W3-C Cursor Validity / Population Change / Deduplication / Limits / Problem Grammar next  
> **Parent Wire Contract:** `D5-B2-WIRE-CONTRACT.md`  
> **Schema authority:** `D5-B2-W2-SCHEMA-GRAMMAR.md`  
> **Operation inventory:** `D5-B2-OPERATION-ADMISSION-MATRIX.md`  
> **Parent authorities:** accepted D0→D4 + D4-R1 + D5-B1 + Decision Reconciliation Baseline + ratified D5-B2 Whole-Matrix + canonical W1/W2  
> **Method:** DevelopmentConexus Engineering Method v1.0.0  
> **W3-A accepted:** 2026-08-19  
> **W3-B accepted:** 2026-08-19

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

# 3. W3-B — Per-Operation Filter / Search / Ordering Matrix — ACCEPTED IN-STAGE

## 3.1 Governing query invariant

> **A collection exposes only typed selection parameters that correspond to a real Product question. Filters never form a general expression language. Every collection has one deterministic owner-defined default order; caller-selected order is absent unless a real consumer proves another semantic ordering is required. Search is a distinct semantic capability, not shorthand for arbitrary filtering.**

Shared query mechanics are deliberately small:

```text
different admitted filter fields → AND
same filter field                → one value baseline
```

There is no baseline `OR`, `NOT`, `IN`, wildcard, field traversal, function/expression grammar, generic range language, generic query DSL or provider/database field path.

Identity lookups do not become list filters when a point `Get` already exists. Correlation/navigation filters are admitted only where a real 1:N/cross-owner read path requires selecting a bounded population.

## 3.2 Typed filters, not generic filter language

Admitted filters are ordinary typed Product query parameters or a typed source-qualified reference encoding in query parameters. Their semantics are closed per operation.

Reject:

```text
filter=...
where=...
filters[field]=...
order_by=...
provider_status=...
field=payload.foo
```

A filter may never expose provider status taxonomy, database column, Sankhya TOP/NUNOTA/CODEMP, provider JSON path or frontend-table column merely because it is available internally.

Time windows are admitted only for operations named below. They use the owner/source event meaning named by that operation and half-open semantics:

```text
[from, before)
```

Generic `from` / `to` / `created_at` range parameters are not a cross-API convention.

## 3.3 Search is distinct from List

`SearchSourceProductsForMarketplace` is the only baseline Search operation.

It requires:

- Marketplace Installation context;
- SourceInstance context;
- non-empty `query`;
- optional `limit`;
- optional `cursor`.

Omitting/blanking `query` never turns Search into a source-Product list/master API.

Bounded matching semantics:

1. exact native Product key / legitimate source identifier match where available;
2. exact source SKU/reference/GTIN match where that evidence is legitimately exposed;
3. case-insensitive lexical token match over bounded source Product display/name evidence;
4. stable SourceProductRef tie-breaker.

No baseline typo correction, edit-distance/fuzzy engine, stemming, embedding/vector semantic search, AI search, provider-specific query syntax or public relevance score.

Search ranking is search-result ordering, not a durable business fact.

## 3.4 Ordering law

Every collection has one server/owner-defined deterministic ordering plus a stable tie-breaker where required.

Baseline caller-selectable sort count: **0**.

Do not expose database row order, provider-native incidental order, `updated_at` churn or a fabricated identity merely to stabilize pagination.

For durable history/Intent collections, prefer immutable occurrence/create time plus stable identity over mutable update time when that satisfies the real consumer.

For small product-defined administrative collections, stable semantic/key ordering is sufficient.

## 3.5 Identity / Portfolio / Readiness matrix

| Operation | Baseline narrowing/search | Default order |
|---|---|---|
| `ListOrganizationMembers` | none | `principal_id ASC` |
| `ListAccessRoles` | none | `role_key ASC` |
| `ListMarketplaceInstallations` | none | `marketplace_installation_id ASC` |
| `ListSellingEntities` | none | `selling_entity_id ASC` |
| `SearchSourceProductsForMarketplace` | required Marketplace Installation + SourceInstance + non-empty `query` | search ranking from §3.3 → stable SourceProductRef |

The four ordinary administrative lists intentionally do not gain member-name search, role filtering, lifecycle filtering or text search before a real Product consumer requires it.

## 3.6 Offering matrix

| Operation | Baseline filters | Default order |
|---|---|---|
| `ListMarketplaceListings` | none; Installation path already qualifies the source namespace | `native_listing_key ASC` |
| `ListListingIntents` | `marketplace_installation_id?`, `lifecycle?` | `created_at DESC` → `listing_intent_id DESC` |
| `ListPriceIntents` | `marketplace_installation_id?` | `created_at DESC` → `price_intent_id DESC` |

`ListingIntent.lifecycle` is the MPC Intent lifecycle, not provider Listing status.

Do not add baseline filters for dispatchability, provider effect, convergence, creator, arbitrary requirement values or provider state. Operational attention remains owner/Work meaning rather than a parallel query engine.

PriceIntent does not gain baseline pending/convergence filters merely because those axes exist in the response.

## 3.7 Availability matrix

| Operation | Baseline filters | Default order |
|---|---|---|
| `ListSellableAvailability` | `marketplace_installation_id?` | stable target tuple (`target.kind` + ListingIntentId or source-qualified Listing ref) |
| `ListInventorySources` | none | `inventory_source_id ASC` |

No baseline quantity/convergence/control/provider-stock filters. A Work-producing Availability exception remains source/Work semantics, not a second exception-query authority.

## 3.8 Market Intelligence matrix

| Operation | Baseline filters/subject | Default order |
|---|---|---|
| `ListCompetitivePositions` | `marketplace_installation_id?` | stable MarketSubject tuple |
| `ListComparableOffers` | required typed MarketSubject | `delivered_price ASC` → stable evidence continuation discriminator |

`ListComparableOffers` is a collection for one bounded MarketSubject, not a global ComparableOffer database.

If the provider exposes no stable comparable-offer identity, MPC does not mint one solely for pagination. D7 may use an opaque acquisition-local continuation discriminator without promoting it to Product identity.

## 3.9 Commercial Economics matrix

| Operation | Baseline filters | Default order |
|---|---|---|
| `ListExpectedEconomics` | `marketplace_installation_id?` | stable EconomicsSubject tuple |
| `ListSaleEconomics` | `marketplace_installation_id?`, `sale_occurred_from?`, `sale_occurred_before?` | `sale_occurred_at DESC` → source-qualified Sale ref |
| `ListEconomicAttributions` | `attribution_state?` | stable attribution-subject tuple |

Sale time filters use the Sale occurrence/business time owned by the Sales/Economics interpretation, not generic record creation time.

`attribution_state` is admitted because human reconciliation needs to select exact/partial/ambiguous/unresolved Economics state. This does not duplicate Work: Economics answers attribution state; Work answers responsibility/actionability.

## 3.10 Governance matrix

| Operation | Baseline filters | Default order |
|---|---|---|
| `ListAuthorizationDecisions` | `decided_from?`, `decided_before?` | `decided_at DESC` → AuthorizationDecision ID |
| `ListAuthorizationDelegations` | `delegate_principal_id?` | `authorization_delegation_id ASC` |

AuthorizationDecision time filtering is an immutable audit/history question. Delegation-by-Principal is a legitimate standing-authority management query.

No generic target-kind/scope-expression/action-expression filter language is admitted.

## 3.11 Marketplace Sales matrix

| Operation | Baseline filters | Default order |
|---|---|---|
| `ListMarketplaceSales` | `sale_occurred_from?`, `sale_occurred_before?`, `selling_entity_id?` | `sale_occurred_at DESC` → native Sale key |

The Installation path qualifies the external Sale namespace.

Selling Entity is a Sales-owned transaction attribution and therefore legitimate narrowing. Provider order/status/substatus filters are not Product semantics.

## 3.12 Materialization matrix

| Operation | Baseline filters | Default order |
|---|---|---|
| `ListBusinessOrderIntents` | typed source-qualified `sale?`, `source_instance_id?` | `created_at DESC` → BusinessOrderIntent ID |
| `ListInvoicingIntents` | `business_order_intent_id?` | `created_at DESC` → InvoicingIntent ID |

These filters are bounded correlation/navigation reads under the zero-P baseline. They let a client follow accepted owner relationships without creating a cross-owner operational projection or client-commanded workflow.

## 3.13 Fulfillment matrix

| Operation | Baseline filters/containment | Default order |
|---|---|---|
| `ListFulfillmentStates` (wire resource `FulfillmentExecution`) | typed source-qualified `sale?`, `fulfillment_node_id?` | `created_at DESC` → `fulfillment_execution_id` |
| `ListFulfillmentNodes` | none | `fulfillment_node_id ASC` |
| `ListFulfillmentArtifacts` | contained/scoped by one FulfillmentExecution | `recorded_at ASC` → `artifact_key` |
| `ListShipments` | typed source-qualified `sale?` | stable source-qualified Shipment order |

W3-B does not invent a universal `shipment_created_at` ordering when selected lanes have not proven one stable cross-provider temporal meaning. A later D4/consumer finding may justify a stronger temporal ordering for that operation only.

`FulfillmentState` remains the operation-semantic label for reads of the one `FulfillmentExecution` resource; no second resource is created.

## 3.14 Post-Sale matrix

| Operation | Baseline filters | Default order |
|---|---|---|
| `ListPostSaleResolutions` | typed source-qualified `sale?`, `lifecycle?` | `created_at DESC` → PostSaleResolution ID |

Sale correlation is essential because one Sale may have 0..N scoped Resolutions. Lifecycle is PostSale-owned `open/closed` meaning, never provider claim/return status.

No baseline refund/return/claim-status boolean filter family is admitted.

## 3.15 Operational Work matrix

`ListWork` has the broadest legitimate Product filter set because assignment/responsibility is Work's accepted business purpose.

Admitted filters:

- `lifecycle?`;
- `responsibility_role_key?`;
- `assigned_principal_id?`;
- `assignment_state? = assigned | unassigned`;
- `origin_kind?`.

Different supplied fields combine by AND.

Rejected baseline filters:

- generic priority/severity/score;
- arbitrary tags;
- free-text search;
- provider status/error field;
- arbitrary source-domain property traversal.

Default order:

```text
created_at ASC
→ work_id ASC
```

This is oldest-obligation-first without inventing prioritization/scheduling authority. A real deadline/source obligation remains visible in owner/Work meaning; caller-selectable deadline/priority ordering requires later consumer proof.

## 3.16 Source qualification and Organization safety in filters

Every external identity filter remains source-qualified.

A Sale/Listing/Shipment/Product reference never becomes a bare native ID query parameter. Exact OpenAPI encoding may flatten a typed reference into several query parameters, but the semantic filter is one qualified reference and must remain unambiguous.

Every Organization-owned secondary reference is resolved inside the path Organization. Query filtering is never a cross-Organization discovery mechanism.

## 3.17 Collection coverage exposure

Explicit collection-level coverage/provenance is required only where the selected population can be incomplete because of external/source acquisition semantics:

- `SearchSourceProductsForMarketplace`;
- `ListMarketplaceListings`;
- `ListComparableOffers`;
- `ListMarketplaceSales`;
- `ListShipments`;
- `ListSaleEconomics` when the underlying Sales/economic evidence universe is not sufficiently complete.

Coverage semantics remain owner-specific and may include bounded scope, partial/unknown/unavailable and material provenance/freshness.

Do not add a generic coverage field to MPC-owned administrative/state collections merely for schema symmetry.

Search exhaustion does not prove Product-master completeness. Marketplace Sales traversal exhaustion does not prove cancellation-inclusive/all-time completeness. ComparableOffer exhaustion does not prove universal market completeness.

## 3.18 No caller-selectable sort baseline

After the complete admitted List/Search sweep, caller-selectable sort remains **zero** in Product 1.0 baseline.

Current collection questions are sufficiently served by:

- search ranking;
- competitive delivered-price ordering;
- immutable occurrence/history time;
- stable resource/reference identity;
- oldest Work obligation first.

Reopen only when a real D6/other supported consumer requires a second stable semantic ordering for one operation.

## 3.19 W3-B negative controls

Later contract proof must reject at least:

1. generic `filter`, `where`, OData/SQL/GraphQL-style expression language;
2. arbitrary `order_by` / sort fields;
3. provider status/substatus filters becoming Product query vocabulary;
4. Sankhya TOP/NUNOTA/CODEMP/native columns becoming filters;
5. provider payload/JSON-path field traversal;
6. blank SearchSourceProducts query becoming an unbounded Product list;
7. bare native Product/Sale/Listing/Shipment IDs without namespace qualification;
8. provider Listing lifecycle enum substituted for ListingIntent lifecycle;
9. Work generic priority/tags/free-text filters without accepted meaning;
10. filters added solely because a React table has that column;
11. generic created-time range applied to every collection by convention;
12. ComparableOffer receiving a synthetic MPC ID merely for a stable tie-breaker;
13. materially different query filters accepted with an old cursor;
14. mutable `updated_at` used for history ordering where immutable create/occurrence time satisfies the consumer;
15. caller-selectable sort added only for API symmetry.

## 3.20 W3-B outcome

**Outcome:** `CURRENT PARENT STRUCTURE CONFIRMED`.

> **Expose the smallest typed query vocabulary that answers each Product question. Ordinary lists, correlation/navigation lists and the one true source-Product search remain distinct; no generic filter/search/sort language is admitted.**

No D0→D5-B1/W1/W2/W3-A reopen is required.

---

# 4. Exact next W3 work

**W3-C — Cursor Validity / Population Change / Deduplication / Limits / Problem Grammar.**

W3-C must close the mechanics/semantics that W3-A/B intentionally leave open:

1. malformed versus query-mismatched versus expired/invalidated cursor behavior and exact Problem Details/status grammar;
2. whether Product cursors have a semantic lifetime promise or only implementation-bounded validity;
3. insert/update/delete behavior between pages without claiming global snapshot isolation;
4. duplicate/omission expectations and whether any collection requires stronger stable-traversal semantics;
5. tie-breaker/continuation behavior for source-backed collections whose external population changes during traversal;
6. `limit` validation/default/maximum policy without inventing one magical scale assumption for all collections;
7. client behavior on stale/invalid cursor — explicit restart/new traversal versus silent continuation;
8. interaction between external collection coverage and continuation when provider enumeration itself expires/changes;
9. whether any operation needs a bounded traversal identifier/snapshot marker in addition to the cursor; reject by default;
10. final W3 Problem Details types and negative controls.

After W3-C, run a Whole-W3 Global Coherence review before accepting W3 as a whole. Do not begin later Wire Contract obligations until W3 is coherent.

Implementation remains blocked until D9.
