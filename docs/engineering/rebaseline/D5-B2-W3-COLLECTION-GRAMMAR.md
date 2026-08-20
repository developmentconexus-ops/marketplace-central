# D5-B2 — W3 Collections / Pagination / Filter / Search / Cursor Grammar

> **Status:** ACCEPTED / CANONICAL — Whole-W3 operator-ratified  
> **Parent Wire Contract:** `D5-B2-WIRE-CONTRACT.md`  
> **Schema authority:** `D5-B2-W2-SCHEMA-GRAMMAR.md`  
> **Operation inventory:** `D5-B2-OPERATION-ADMISSION-MATRIX.md`  
> **Parent authorities:** accepted D0→D4 + D4-R1 + D5-B1 + Decision Reconciliation Baseline + ratified D5-B2 Whole-Matrix + canonical W1/W2  
> **Method:** DevelopmentConexus Engineering Method v1.0.0  
> **Whole-W3 final ratification incorporated:** 2026-08-19

## 1. Purpose and authority

W3 is the **single canonical Product API collection/query authority** for the already-admitted List/Search Q operations.

It consolidates the formerly staged W3-A/B/C work plus the operator-ratified Whole-W3 lead/Fable/GPT adversarial review. Git history preserves those staging/review snapshots; they are not parallel active authorities.

> **Uniform continuation mechanics may be shared; collection population, owner meaning, coverage, filters, search semantics and ordering remain operation-specific. A cursor carries continuation, never business selection authority. Pagination state never becomes knowledge/completeness authority.**

W3 does not derive contracts from legacy OpenAPI/routes/controllers/frontend tables and does not choose D7 database/index/cache/provider-paging realization.

---

# 2. Canonical collection core

## 2.1 Named collection responses; no universal page/result wrapper

Every List/Search operation owns a response schema named for that owner/operation meaning.

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

## 2.2 ListItem semantic-subset law

A List operation is not required to return the full point-Get representation.

Where a collection member has an admitted point resource/Q, its ListItem is a **semantic subset of that same owner meaning**:

- it may omit detail, history, evidence or large nested sections;
- it may carry the identity/current owner fields needed to scan, select and navigate;
- it may **not** introduce a list-only derived business conclusion;
- when a field is reused, its schema/name/meaning remains the same except for qualifiers already made unambiguous by the operation/path scope.

For collection-only owner values without a point Get, the item remains the owner-native W2 value schema and likewise cannot invent a second business conclusion.

Do not introduce generic `Summary<T>`, `View<T>`, `ListResource<T>`, `fields`, `select`, `expand` or projection DSL machinery. A naturally small resource may reuse its point representation when honest.

## 2.3 Shared request grammar

All admitted List/Search operations are pagination-capable using only the shared continuation mechanics:

```text
limit?
cursor?
```

`limit` is a requested maximum, never a promise that exactly that many members are returned.

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

`pagination-capable` does not mean a small collection must issue a cursor in normal use.

## 2.4 `next_cursor`

When another traversal page exists, the response may contain:

```text
next_cursor: opaque non-empty string
```

When no subsequent traversal page exists, `next_cursor` is omitted. `null` does not mean end-of-pagination.

A page may contain fewer members than `limit`, including zero members, and still carry `next_cursor`; item count alone never determines exhaustion.

## 2.5 Cursor opacity and authority fence

A Product cursor is an opaque continuation token. Clients cannot depend on its encoding or inspect it for ordering/source/provider state.

D7 may internally encode/reference proportionately:

- ordering continuation state;
- provider continuation state;
- server-side continuation identity;
- semantic-query fingerprint;
- acquisition/coverage boundary;
- bounded seen-set/search execution state.

None becomes Product API business meaning.

Raw provider paging tokens, scroll IDs, database keys or transparent base64-encoded implementation state must not leak as Product contract.

A cursor never proves Principal identity, Membership, Permission, current source access or business authorization. Every continuation request is authenticated/authorized normally.

## 2.6 Explicit semantic query on every continuation request

A cursor carries **continuation state only**; it never replaces the explicit Product query.

Every continuation request must repeat:

- every operation-required subject/search parameter;
- the same effective optional narrowing/filter semantics selected on the first request;
- the same Organization/operation scope.

An optional filter omitted on the first request remains semantically omitted unless that operation explicitly defines an equivalent default.

Only `limit` may vary between continuation requests because it changes requested response cardinality, not the selected population.

The cursor may contain a query fingerprint for verification. The fingerprint verifies the explicit query; it does not become the only carrier of that query.

Processing precedence:

1. decode/basic request contract and required query fields;
2. ordinary authentication/access/privacy checks;
3. cursor validation against the well-formed explicit semantic query.

Therefore:

- a missing/invalid operation-required subject/search/filter field is ordinary W2 `422 validation-error`;
- a well-formed request whose semantic query does not match the supplied cursor is W2/W3 `400 invalid-cursor`.

## 2.7 Pagination state is not coverage/completeness

Absence of `next_cursor` means only:

> no subsequent page exists in this Product traversal.

It does **not** prove:

- universal source/provider completeness;
- all-time population completeness;
- cancellation-inclusive Sales completeness;
- complete market universe;
- complete Product master universe;
- owner knowledge beyond that collection's explicit coverage/evidence contract.

A response may therefore have no `next_cursor` while still declaring partial/unknown/unavailable coverage where that owner exposes collection coverage.

```text
cursor exhaustion
!= deduplication
!= source enumeration completeness
!= owner knowledge completeness
```

## 2.8 No baseline total count

W3 admits no universal `total`, `total_count` or estimated-count field.

A later real consumer may justify an operation-specific count only when the counted universe, filters, freshness/coverage and authoritative semantics are explicit and honest.

Database COUNT convenience is not Product-contract evidence.

## 2.9 No caller-selectable sort baseline

Baseline caller-selectable sort count is **0**.

There is no generic:

```text
sort
order_by
sort_by
field,-field
```

Each admitted collection has one owner-defined deterministic default order below. A second caller-selectable ordering requires a real supported consumer and a stable semantic contract.

## 2.10 No universal snapshot-isolation promise

A cursor traversal does not imply one immutable collection snapshot across pages unless an operation-specific contract explicitly says so.

W3 guarantees honest continuation semantics, not a global transaction snapshot.

The one ComparableOffer-specific evaluation/acquisition basis in §5.5 is an owner-local correctness requirement for identity-less ordered evidence; it is not a shared Snapshot/Traversal abstraction.

## 2.11 Provider pagination remains D4-local mechanism

```text
Product cursor
    ↓
owner/application collection semantics
    ↓
D4 adapter
    ↓
provider token / offset / scroll / search-after / provider-specific paging
```

Provider pagination changes do not automatically alter Product cursor grammar. D7 owns persistence/signing/index/cache/seen-set realization.

---

# 3. Query / filter / search grammar

## 3.1 Governing query invariant

> **Expose the smallest typed query vocabulary that answers the Product question, not a language capable of asking future questions.**

Admitted filter fields are ordinary typed Product query parameters or typed source-qualified references whose semantics are closed per operation.

Shared composition is deliberately small:

```text
different admitted filter fields → AND
same filter field                → one value baseline
```

No baseline:

- OR / NOT / IN;
- wildcard expression grammar;
- field traversal;
- function/expression language;
- OData/SQL/GraphQL-like filter DSL;
- generic range language;
- provider/database JSON path;
- provider status taxonomy as Product filters.

Identity point lookups do not become list filters where a point Get already exists. Correlation/navigation filters exist only for real bounded 1:N/cross-owner read paths.

## 3.2 Source qualification and Organization safety

Every external identity filter remains source-qualified.

A Product/Sale/Listing/Shipment reference never becomes a bare native ID query parameter. Exact OpenAPI encoding may flatten a typed reference into several query parameters, but the semantic reference remains qualified and unambiguous.

Every Organization-owned secondary reference resolves inside the path Organization. Query filtering is never a cross-Organization discovery mechanism.

## 3.3 Time windows

Time windows exist only for operations named in the matrix below and use the owner/source event meaning named by that operation.

Intervals use half-open semantics:

```text
[from, before)
```

Generic `from`, `to` or cross-API `created_at` filters are not admitted by convention.

## 3.4 SearchSourceProductsForMarketplace

`SearchSourceProductsForMarketplace` is the only baseline Search operation.

Every request, including continuation requests, requires:

- Marketplace Installation context;
- SourceInstance context;
- non-empty opaque user `query`;
- optional `limit`;
- optional `cursor`.

Omitting/blanking `query` never turns Search into a source-Product list/master API.

The Search contract returns source-qualified evidence only and uses bounded matching over legitimate source identification/display evidence supported by the sanctioned source contract.

Where the sanctioned source contract establishes them, exact native Product key / exact legitimate identifier matches rank ahead of textual matches.

Textual matching/ranking must be deterministic enough for one traversal, but W3 freezes **no universal**:

- tokenizer/token boundary algorithm;
- case-folding/case-insensitive algorithm;
- diacritics/accent-folding behavior;
- locale/collation behavior;
- stemming;
- typo/edit-distance/fuzzy engine;
- vector/embedding/AI semantic search.

Provider query syntax never crosses the Product API. No public relevance score is admitted.

Stable SourceProductRef is the Product search-result member identity/tie-breaker. Within one traversal the same SourceProductRef is returned at most once.

If a real SourceInstance cannot satisfy the materially required Search through sanctioned reads/projections without introducing a new MPC Product-search mirror/index/data authority, **STOP / targeted re-adjudication**. Do not silently build the mirror to satisfy W3 wording.

---

# 4. Canonical 26-operation List/Search matrix

The ratified operation inventory contains **26 admitted List/Search Q operations**. W3 introduces zero additional List/Search operations by symmetry.

## 4.1 Identity / Access / Portfolio / Readiness

| Operation | Enumerable population / subject | Baseline narrowing | Default order | Collection coverage |
|---|---|---|---|---|
| `ListOrganizationMembers` | current Organization Membership relations visible under `access.read` | none | `principal_id ASC` | no generic field |
| `ListAccessRoles` | Product-defined AccessRole definitions | none | `role_key ASC` | no generic field |
| `ListMarketplaceInstallations` | Organization MarketplaceInstallation resources | none | `marketplace_installation_id ASC` | no generic field |
| `ListSellingEntities` | Organization Portfolio SellingEntity registry | none | `selling_entity_id ASC` | no generic field |
| `SearchSourceProductsForMarketplace` | bounded source search results for explicit Installation + SourceInstance + query; **not** a Product universe | required Installation + SourceInstance + non-empty query | source-capability-backed search ranking → `SourceProductRef` | explicit owner/source coverage/provenance |

The four administrative Lists intentionally do not gain text/lifecycle/member filters before a real Product consumer requires them.

## 4.2 Offering

| Operation | Enumerable population | Baseline filters | Default order | Collection coverage |
|---|---|---|---|---|
| `ListMarketplaceListings` | currently known external Listing subjects within the Installation namespace | none; Installation path qualifies source namespace | `native_listing_key ASC` | explicit Listing acquisition coverage/provenance |
| `ListListingIntents` | MPC ListingIntent resources, including current/historical lifecycle states | `marketplace_installation_id?`, `lifecycle?` | `created_at DESC` → `listing_intent_id DESC` | no generic field |
| `ListPriceIntents` | MPC PriceIntent resources | `marketplace_installation_id?` | `created_at DESC` → `price_intent_id DESC` | no generic field |

`ListingIntent.lifecycle` means the MPC Intent lifecycle, never provider Listing status. No baseline dispatchability/provider-effect/convergence/creator filters.

## 4.3 Availability

| Operation | Enumerable population | Baseline filters | Default order | Collection coverage |
|---|---|---|---|---|
| `ListSellableAvailability` | current Offering target universe defined in §4.3.1 | `marketplace_installation_id?` | fixed target-kind order → target identity | no new generic coverage field; existing-Listing universe bounded as below |
| `ListInventorySources` | MPC InventorySource resources | none | `inventory_source_id ASC` | no generic field |

### 4.3.1 SellableAvailability current population

`ListSellableAvailability` enumerates the **currently known Offering target universe for which SellableAvailability is a current meaning**, never persistence rows that happen to be configured/materialized.

It contains:

1. **pre-creation members** — ListingIntents whose authoring target is `new_listing`, whose Intent lifecycle is `draft` or `submitted`, and for which no provider Listing has yet been established as the current external Listing target;
2. **post-creation/current external members** — currently known existing marketplace Listing subjects.

Consequences:

- `discarded` ListingIntents are excluded from the current Availability population;
- an incomplete `draft` remains included: Availability may honestly be unknown/unavailable; omission must not hide the target merely because knowledge is weak;
- a submitted intent remains included while genuinely pre-creation because dispatch may still depend on Availability and other prerequisites;
- once the provider Listing is established, the current member is addressed by `existing_listing`; the pre-creation member leaves the current population rather than remaining a second current authority;
- completeness of the existing-Listing portion is bounded by Listing acquisition coverage from `ListMarketplaceListings`; no duplicate collection-level coverage field is introduced.

Default total order:

```text
1. pre_creation_listing_intent → listing_intent_id ASC
2. existing_listing            → source-qualified Listing ref ASC
```

The closed target-kind order is deterministic grouping only, not business priority.

No baseline quantity/convergence/control/provider-stock filters.

## 4.4 Market Intelligence

| Operation | Enumerable population / subject | Baseline filters/subject | Default order | Collection coverage |
|---|---|---|---|---|
| `ListCompetitivePositions` | **currently known existing Listing subjects only** | `marketplace_installation_id?` | source-qualified Listing ref ASC | no new field; universe completeness bounded by Listing acquisition coverage |
| `ListComparableOffers` | ComparableOffer evidence for one required MarketSubject and one owner-local evaluation/acquisition basis | required typed MarketSubject | `delivered_price ASC` → stable provider member identity when real, else acquisition-local discriminator | explicit Market coverage/provenance |

The point `GetCompetitivePosition` subject union remains wider than the List population and may address `source_product_marketplace_context` for pre-listing analysis.

Pre-listing Product 1.0 flow remains:

```text
SearchSourceProductsForMarketplace
→ selected source Product ref
→ point GetCompetitivePosition / Economics point Q or EvaluatePriceScenario
```

There is no hidden enumeration of every source Product × marketplace context and no enumeration of whichever keyed Qs happen to be cached/materialized.

`ListCompetitivePositions` traversal exhaustion therefore never proves complete marketplace portfolio coverage; its enumerable universe derives from currently known Listing knowledge and is bounded by Listing acquisition coverage without adding a duplicate coverage field.

## 4.5 Commercial Economics

| Operation | Enumerable population | Baseline filters | Default order | Collection coverage |
|---|---|---|---|---|
| `ListExpectedEconomics` | **currently known existing Listing subjects only** | `marketplace_installation_id?` | source-qualified Listing ref ASC | no new field; universe completeness bounded by Listing acquisition coverage |
| `ListSaleEconomics` | known SaleEconomics meanings keyed by source-qualified Sales | `marketplace_installation_id?`, `sale_occurred_from?`, `sale_occurred_before?` | `sale_occurred_at DESC` → source-qualified Sale ref | explicit where underlying Sales/economic universe is incomplete |
| `ListEconomicAttributions` | persistent EconomicAttribution resources | `attribution_state?` | `economic_attribution_id ASC` | no generic field |

The point ExpectedEconomics subject union remains wider than the List population and may address a source Product + marketplace context explicitly.

No subject-grouping order is added to EconomicAttribution: the reconciliation consumer already narrows by `attribution_state`; opaque ID order provides one stable total order without fabricating priority.

## 4.6 Governance

| Operation | Enumerable population | Baseline filters | Default order | Collection coverage |
|---|---|---|---|---|
| `ListAuthorizationDecisions` | immutable AuthorizationDecision occurrences | `decided_from?`, `decided_before?` | `decided_at DESC` → AuthorizationDecision ID | no generic field |
| `ListAuthorizationDelegations` | AuthorizationDelegation resources | `delegate_principal_id?` | `authorization_delegation_id ASC` | no generic field |

Decision time filtering is an immutable audit/history question. Delegation-by-Principal is a standing-authority management query; no generic target/scope/action expression filters.

## 4.7 Marketplace Sales

| Operation | Enumerable population | Baseline filters | Default order | Collection coverage |
|---|---|---|---|---|
| `ListMarketplaceSales` | currently known external Sale subjects within the Installation namespace | `sale_occurred_from?`, `sale_occurred_before?`, `selling_entity_id?` | `sale_occurred_at DESC` → native Sale key | explicit Sales acquisition coverage/provenance |

SellingEntity narrowing is legitimate Sales-owned transaction attribution. Provider order/status/substatus filters are not Product semantics.

## 4.8 Business-System Materialization

| Operation | Enumerable population | Baseline filters | Default order | Collection coverage |
|---|---|---|---|---|
| `ListBusinessOrderIntents` | MPC BusinessOrderIntent resources | typed source-qualified `sale?`, `source_instance_id?` | `created_at DESC` → BusinessOrderIntent ID | no generic field |
| `ListInvoicingIntents` | MPC InvoicingIntent resources | `business_order_intent_id?` | `created_at DESC` → InvoicingIntent ID | no generic field |

These are bounded correlation/navigation reads under the zero-P baseline; they do not create a cross-owner operational projection or client-commanded workflow.

## 4.9 Fulfillment

| Operation | Enumerable population / containment | Baseline filters | Default order | Collection coverage |
|---|---|---|---|---|
| `ListFulfillmentStates` (wire resource `FulfillmentExecution`) | MPC FulfillmentExecution resources | typed source-qualified `sale?`, `fulfillment_node_id?` | `created_at DESC` → `fulfillment_execution_id` | no generic field |
| `ListFulfillmentNodes` | MPC FulfillmentNode resources | none | `fulfillment_node_id ASC` | no generic field |
| `ListFulfillmentArtifacts` | artifacts contained/scoped by one FulfillmentExecution | parent FulfillmentExecution scope | `recorded_at ASC` → `artifact_key` | no generic field |
| `ListShipments` | currently known source-qualified Shipment subjects | typed source-qualified `sale?` | Installation-qualified native Shipment key ASC | explicit Shipment acquisition coverage/provenance |

`FulfillmentState` remains an operation-semantic label over the one `FulfillmentExecution` wire resource. W3 does not invent an unproven universal shipment timestamp merely for sorting.

## 4.10 Post-Sale

| Operation | Enumerable population | Baseline filters | Default order | Collection coverage |
|---|---|---|---|---|
| `ListPostSaleResolutions` | MPC PostSaleResolution resources | typed source-qualified `sale?`, `lifecycle?` | `created_at DESC` → PostSaleResolution ID | no generic field |

Lifecycle is PostSale-owned open/closed meaning, never provider claim/return status. No refund/return/claim boolean filter family by symmetry.

## 4.11 Operational Work

`ListWork` enumerates MPC Work obligations.

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

This is oldest-obligation-first without inventing prioritization/scheduling authority. A real source deadline remains owner/Work meaning; a second ordering requires consumer proof.

## 4.12 Matrix completeness result

```text
admitted List/Search Qs                         26
W3 collection homes                             26
additional List/Search operations by symmetry    0
caller-selectable sort operations                 0
```

---

# 5. Collection stability / cursor safety

## 5.1 Cursor Problem Details

W2 is the single Product Problem Details catalog. W3 defines applicability for its two cursor-specific types:

### `invalid-cursor` — HTTP 400

Use when a supplied cursor cannot validly continue a **well-formed** request, including:

- malformed/unparseable token;
- integrity/signature/tamper failure;
- unknown token;
- operation mismatch;
- Organization mismatch after normal privacy/access checks;
- semantic query/filter/search mismatch.

### `cursor-expired` — HTTP 400

Use when the cursor was legitimately issued but that exact continuation can no longer be resumed honestly, e.g. required server-side continuation state retired or a provider continuation definitively expired and cannot be translated without changing Product semantics.

Recovery:

```text
discard expired cursor
→ repeat the semantic query without cursor
→ begin a fresh traversal
```

Never silently restart at page one or continue from a best-effort nearby position.

`410 Gone` is not used merely because a cursor expired: the collection target still exists and the cursor is not a Product resource. `409` is not baseline because ordinary population evolution is not a resource conflict. Do not normalize these cursor cases into `422`; W2 deliberately distinguishes required-query validation from cursor continuation validity.

No `cursor-stale`, `cursor-gone`, `cursor-conflict` or provider cursor taxonomy.

## 5.2 Cursor lifetime

A Product cursor is an **ephemeral continuation token, not a persistent bookmark**.

W3 freezes no public minimum TTL. D7 chooses practical retention.

A future real consumer that must suspend/resume traversal for a materially long/offline period is the reopen trigger for a stronger lifetime/bookmark contract.

Cursor lifetime is separate from source freshness and collection coverage.

## 5.3 Population mutation between pages

W3 does not promise one immutable collection snapshot across a traversal.

Changes between calls may affect members not yet across the continuation boundary:

- a newly created member whose order position is before the boundary may not appear until a new traversal;
- a member deleted before being reached may never appear;
- a member that stops matching before being reached may disappear;
- a member that starts matching may or may not appear depending on resulting order position.

These are not pagination defects under the no-snapshot contract.

A client that needs the current universe starts a new traversal.

## 5.4 Stable-identity at-most-once guarantee

For collection members with stable MPC or source-qualified identity/key:

> **The same semantic member identity MUST NOT be returned more than once in one Product traversal.**

This applies to MPC resources and source-qualified Listing/Sale/Shipment/SourceProductRef/keyed-Q subjects with stable identity.

Provider duplicate pages/reordered duplicate delivery do not leak directly into Product collection semantics.

D7 may use keyset boundaries, provider-stable continuation, bounded seen-set/traversal state or another mechanism that proves the property. D7 may not silently weaken it.

At-most-once identity does **not** imply snapshot, no-omission or completeness.

A target re-keying such as SellableAvailability moving from `pre_creation_listing_intent` to `existing_listing` is current-population mutation: the two target identities are not the same semantic member identity for this guarantee.

## 5.5 ComparableOffer owner-local resumable basis

`ComparableOffer` deliberately may have no stable external/MPC identity while its ordering key (`delivered_price`) can change.

Therefore one `ListComparableOffers` cursor chain is bound to **one Market Intelligence evaluation/acquisition basis**:

- the response/provenance identifies that Market owner evaluation sufficiently for honesty;
- `delivered_price ASC` is evaluated inside that basis;
- a stable provider-native member identity is used where real;
- otherwise an acquisition-local discriminator may exist only as traversal mechanism;
- D7 may retain the existing Market evaluation result only long enough to resume the chain;
- if the same basis cannot be resumed, return `cursor-expired` rather than re-querying current market and pretending it is continuation;
- a request without cursor starts a fresh current evaluation/traversal.

`evaluation/acquisition basis` is owner-local description, **not** a shared named `EvaluationBasis`, `TraversalBasis`, Snapshot, session or Product resource type.

## 5.6 Non-identifiable evidence members

No evidence receives a fabricated canonical MPC identity merely to make deduplication/pagination convenient.

Where D4 exposes a stable qualified external identity, it may support traversal identity. Where it does not, an acquisition-local discriminator:

- is technical only;
- need not be publicly resolvable later;
- does not prove two value-equal offers/evidence values are the same external member.

If sufficient member continuity cannot be preserved, cursor expiration is honest; invented identity is not.

## 5.7 Updates do not turn pagination into a change feed

If a member has already been returned and later changes a non-identity/non-ordering field, it is not re-emitted merely because its representation changed.

Current state is obtained by point Get, a new traversal or a later independently admitted event/change-feed contract. Pagination is not a subscription mechanism.

Ordering therefore prefers immutable create/occurrence time or stable identity/tuple over mutable `updated_at` when that serves the consumer.

## 5.8 Provider continuation failure

If a Product cursor internally depends on provider continuation state and the provider continuation expires/changes:

1. D4/D7 may reconstruct it transparently only when the **same Product continuation semantics** remain provable;
2. otherwise return `cursor-expired`;
3. never silently restart, skip approximately or return a partial page while pretending the old continuation was satisfied.

A temporary provider/source outage is **not** automatically cursor expiration. It is availability/coverage/service behavior unless the continuation itself is definitively lost.

## 5.9 No public traversal/snapshot resource

No baseline:

```text
traversal_id
snapshot_id
search_session_id
```

The opaque cursor is the sufficient Product continuation seam. D7 may internally persist traversal state without creating a Product authority.

---

# 6. `limit` contract

W3 freezes semantic validation but deliberately does not invent numeric scale assumptions.

Binding behavior:

```text
limit omitted
  → use documented Product default

limit > 0
  → requested maximum number of members in the response

server returned members
  → MUST be <= limit

limit = 0 or negative
  → validation error

limit > documented maximum
  → validation error; no silent clamp
```

Small collections obey the same contract. `limit=1` never permits returning more than one member merely because the collection is small.

## 6.1 Exact numeric default/max — DEFER SAFELY

Exact numeric default/max remain intentionally deferred to the final OpenAPI contract/tooling sub-batch after concrete ListItem wire representations are serialized/measured.

Binding fences already fixed:

- finite positive default;
- finite positive maximum;
- one shared pair preferred when evidence supports it;
- operation-specific exception only for measured payload/PII/provider/consumer constraint;
- no silent clamp;
- later OpenAPI work may fill numbers but may not redesign W3 pagination semantics.

This is a bounded `DEFER SAFELY`, not permission for unbounded responses.

---

# 7. Retry semantics

List/Search Product operations are GET reads and do not use `Idempotency-Key`.

A same-request retry may succeed, fail normally or discover that the continuation expired. If expired, the client begins a new traversal explicitly.

Read retry never turns an expired cursor into a silent restart or consequential intake concept.

---

# 8. Canonical W3 negative controls

Later OpenAPI/runtime proof must reject at least:

1. raw provider cursor/scroll token exposed directly;
2. cursor accepted under another Organization/operation/query;
3. continuation request omitting operation-required semantic subject/search fields and relying on cursor as hidden query authority;
4. materially changed optional filters accepted with an old cursor;
5. cursor treated as authorization;
6. cursor exhaustion interpreted as source/provider/market/Product completeness;
7. `len(items) < limit` interpreted as traversal exhaustion;
8. generic total/estimated count added because storage can count;
9. offset/page-number/previous-cursor baseline by familiarity;
10. arbitrary caller sort expression;
11. universal snapshot/traversal/session resource;
12. universal `Page<T>`/`metadata` wrapper;
13. generic `filter`, `where`, OData/SQL/GraphQL-like query language;
14. provider status, Sankhya TOP/NUNOTA/CODEMP or provider JSON path as Product filters;
15. blank Source Product query becoming an unbounded Product list;
16. a Product mirror/index created silently only to honor an unproven search algorithm promise;
17. universal tokenizer/case/diacritics/locale/collation/fuzzy/vector search promise without source evidence;
18. bare native Product/Sale/Listing/Shipment identity without namespace qualification;
19. provider Listing lifecycle substituted for ListingIntent lifecycle;
20. Work priority/tags/free-text filters without accepted meaning;
21. generic created-time range applied to every collection;
22. mutable `updated_at` history ordering where immutable ordering suffices;
23. ListCompetitivePositions/ListExpectedEconomics enumerating all source Products or cache/materialization rows;
24. ListSellableAvailability enumerating only persisted/configured availability rows and silently dropping unknown targets;
25. full ListingIntent historical dispatch basis repeated in every ListItem by schema symmetry;
26. list-only derived business conclusion absent from the point owner meaning;
27. generic `fields/select/expand` projection language;
28. same stable Sale/Listing/Work/SourceProduct identity returned twice in one traversal;
29. provider duplicate member passed through unchanged;
30. ComparableOffer receiving fake canonical ID solely for pagination;
31. equal-looking ComparableOffers collapsed without identity evidence;
32. ComparableOffer page 2 re-querying a new market basis while pretending to continue page 1;
33. expired cursor silently restarting;
34. `410 Gone` promoting cursor to resource or `409` treating ordinary population mutation as conflict;
35. transient provider outage mislabeled cursor expiration;
36. provider continuation loss approximated by nearby skip/restart while returning success;
37. at-most-once identity interpreted as no-omission/snapshot/completeness;
38. server returning more members than caller `limit`;
39. small collection ignoring caller `limit`;
40. exact arbitrary default/max numbers introduced without concrete payload/consumer evidence.

---

# 9. Reopen triggers

Reopen only the smallest W3 decision when material evidence proves one of these conditions:

- a supported client needs a second semantic sort/filter/count that cannot be served by the baseline vocabulary;
- a real pre-listing portfolio consumer requires an enumerable universe beyond existing Listings; any extension must prove a bounded owner-known universe (for example established correspondence), never “all Products” by default;
- a SourceInstance cannot satisfy the admitted source Product Search through sanctioned reads/projections without a materially new Product-search data authority;
- a real client requires durable/offline cursor bookmark lifetime;
- D7 proves the at-most-once guarantee has disproportionate unavoidable cost for a concrete collection; D7 must surface this for W3 adjudication rather than weakening the contract silently;
- a source/provider proves a stronger or different collection identity/order requirement materially necessary for correctness.

Current Whole-W3 review proves **no D0→D5-B1/W1 semantic parent reopen**. W2 was amended only to preserve one canonical Product Problem Details catalog.

---

# 10. Whole-W3 method disposition

```text
D0→D5-B1 / ratified B2 semantic authority     CURRENT STRUCTURE CONFIRMED
Wire W1                                        ACCEPTED / CANONICAL
Wire W2                                        ACCEPTED / CANONICAL
W3 collection/query grammar                    ACCEPTED / CANONICAL
Whole-W3 lead review                           RESTRUCTURE W3-LOCAL
Fable independent Whole-W3 review              COMPLETE
GPT final adjudication                         CONVERGED
Operator final Whole-W3 ratification           COMPLETE
Parent-stage reopen                            NONE
```

**Global Maximum:** explicit owner/query meaning + minimal shared continuation mechanics + honest source-bound coverage/stability, without Product/PIM mirror, generic Page/query/search/sort/projection/snapshot/traversal framework or fabricated evidence identity.

W3 is closed/canonical. Remaining Wire Contract obligations continue only in router order. Implementation remains blocked until D9.
