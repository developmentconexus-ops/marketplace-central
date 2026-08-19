# D5-B2 — W3-C Cursor Validity / Population Change / Deduplication / Limits / Problem Grammar

> **Status:** ACCEPTED IN-STAGE / OPERATOR-RATIFIED  
> **Parent W3:** `D5-B2-W3-COLLECTION-GRAMMAR.md`  
> **Parent Wire Contract:** `D5-B2-WIRE-CONTRACT.md`  
> **Schema authority:** `D5-B2-W2-SCHEMA-GRAMMAR.md`  
> **Operation inventory:** `D5-B2-OPERATION-ADMISSION-MATRIX.md`  
> **Parent authorities:** accepted D0→D4 + D4-R1 + D5-B1 + Decision Reconciliation Baseline + ratified D5-B2 Whole-Matrix + canonical W1/W2 + accepted W3-A/W3-B  
> **Method:** DevelopmentConexus Engineering Method v1.0.0  
> **Accepted:** 2026-08-19

## 1. Purpose

W3-C closes the Product API collection guarantees that W3-A/B intentionally left open: cursor validity, population mutation between pages, duplicate/member stability, list-item projections, limit behavior and cursor Problem Details.

W3-C does not choose D7 cursor persistence/signing/index/cache/provider-token realization and does not introduce a traversal/snapshot business resource.

> **A Product collection is an ordered traversal of owner meaning, not a universal transactional snapshot. Stable members are not repeated within one traversal; population change may affect not-yet-traversed members; cursor failure is explicit; pagination never fabricates source completeness or durable bookmark semantics.**

---

## 2. List item projection law

A List operation is not required to return the full point-Get representation.

For owner meanings whose full representation contains material history, evidence or large nested detail, the collection may expose a bounded owner-specific list-item projection sufficient to:

- identify the member;
- scan/select it;
- navigate to the point resource/Q;
- understand the material owner state needed for the collection use.

Example:

```text
GetListingIntent
  → full current detail
  → append-only dispatch/effect basis
  → richer evidence/history

ListListingIntents
  → ListingIntentListItem
  → identity / target / lifecycle / bounded scan fields
```

This is a read projection of the same owner meaning, never a second business authority.

Do not introduce generic `Summary<T>`, `ListResource<T>`, `fields`, `select`, `expand` or projection DSL machinery. A naturally small resource may reuse its point representation when honest.

---

## 3. Cursor problem classes

W3 admits only two cursor-specific Problem Details types.

### 3.1 `invalid-cursor` — HTTP 400

Use when the supplied cursor cannot validly continue the request, including:

- malformed/unparseable token;
- integrity/signature/tamper failure;
- unknown token;
- cursor bound to another operation;
- cursor bound to another Organization after normal access/privacy checks;
- material semantic query/filter/search mismatch.

Cursor invalidity is request invalidity, not ordinary access authority and not a business outcome.

### 3.2 `cursor-expired` — HTTP 400

Use when the cursor was legitimately issued but that exact continuation can no longer be resumed honestly, e.g.:

- required server-side continuation state retired;
- provider continuation window/token definitively expired;
- cursor generation intentionally retired and cannot be safely translated.

The Product contract does not expose those internal causes as a taxonomy.

Recovery is explicit:

```text
discard expired cursor
→ repeat the original query without cursor
→ begin a new traversal
```

Never silently restart at page one or continue from a best-effort nearby position.

`410 Gone` is not used merely because a cursor expired: the Product collection target resource still exists and the cursor is not promoted to a Product resource. `409` is likewise not the baseline because normal population evolution is not itself a conflict.

---

## 4. Cursor lifetime

A Product cursor is an **ephemeral continuation token, not a persistent bookmark**.

W3 freezes no public minimum TTL such as one hour/day/week. D7 chooses practical storage/retention and may allow ordinary client traversal to resume for a reasonable operational window.

A future real consumer that must suspend and resume traversals for a materially long/offline period is the reopen trigger for a stronger lifetime contract.

Cursor lifetime remains separate from source/provider data freshness and collection coverage.

---

## 5. Population mutation between pages

W3 does not promise one immutable collection snapshot across a traversal.

For a stable ordered traversal, changes between calls may affect members that have not yet crossed the continuation boundary.

Examples:

- a newly created member whose order position is before the current boundary may not appear until a new traversal;
- a member deleted before being reached may never appear;
- a member that stops matching a filter before being reached may disappear;
- a member that starts matching during traversal may or may not appear depending on its resulting position relative to the continuation boundary.

These are not Product pagination defects when the no-snapshot contract is preserved.

A client that needs the current universe starts a new traversal.

---

## 6. Stable-identity deduplication guarantee

For collection members with stable MPC or source-qualified identity/key:

> **the same semantic member identity MUST NOT be returned more than once in one Product traversal.**

This includes proportionately:

- MPC-owned resources such as ListingIntent, PriceIntent, Work, AuthorizationDecision, InventorySource and FulfillmentExecution;
- external/source-qualified Listing, Sale, Shipment and SourceProductRef.

Provider duplicate pages/reordered duplicate delivery do not leak directly into Product collection semantics.

D7 chooses the mechanism: keyset boundaries, provider-stable continuation, seen-set/traversal state or another implementation that proves the property.

At-most-once member identity does **not** imply snapshot or completeness.

---

## 7. Updates do not turn pagination into a change feed

If a member has already been returned and later changes a non-identity/non-ordering field, it is not re-emitted merely because its representation changed.

Example:

```text
Work W1 returned as open
→ W1 later becomes held
→ same traversal does not return W1 again for that update
```

Current state is obtained by point `Get`, a new traversal, or a later independently admitted event/change-feed contract. Pagination is not a subscription mechanism.

W3-B ordering therefore deliberately prefers immutable creation/occurrence time or stable semantic tuple over mutable `updated_at` when that satisfies the consumer.

---

## 8. Non-identifiable evidence members

ComparableOffer and similar evidence do not receive a fabricated MPC identity merely to make deduplication/pagination convenient.

Where D4 exposes a stable qualified external member identity, it may support continuation/deduplication under that evidence contract.

Where no stable business/external identity exists, D7 may use an opaque acquisition-local member discriminator solely for traversal mechanics. It:

- is not a canonical Product ID;
- need not be publicly resolvable later;
- does not prove two value-equal offers are the same external offer.

If the source cannot preserve sufficient member continuity to resume honestly, the cursor may expire rather than invent identity or silently approximate continuation.

---

## 9. SearchSourceProducts traversal stability

`SearchSourceProductsForMarketplace` uses stable SourceProductRef as the member identity.

Within one search traversal the same SourceProductRef is returned at most once even if source ranking/data changes between provider pages.

Search ranking may evolve as source evidence changes; W3 does not promise that every Product that could have matched at some point during the traversal will appear. A fresh search starts a fresh current-universe traversal.

---

## 10. Provider continuation failure

If a Product cursor internally depends on provider continuation state and the provider continuation expires/changes:

1. D4/D7 may reconstruct continuation transparently only when the **same Product continuation semantics** can still be proven;
2. otherwise return `400 cursor-expired`;
3. never silently restart enumeration, skip approximately, or return a partial page while pretending the old continuation was satisfied.

A temporary provider/source outage is **not** automatically cursor expiration. It is availability/coverage/service behavior. The cursor remains conceptually valid unless the continuation itself is definitively lost.

---

## 11. No public traversal/snapshot resource

W3 admits no baseline:

```text
traversal_id
snapshot_id
search_session_id
```

The opaque cursor is the sufficient Product continuation seam.

D7 may internally persist traversal state, seen sets, provider tokens, ordering high-water marks or search execution state without creating a Product business/resource authority.

---

## 12. `limit` contract

W3 freezes semantic validation but deliberately does not invent numeric scale assumptions.

Binding behavior:

```text
limit omitted
  → use documented Product default

limit > 0
  → requested maximum number of members in that response

server returned members
  → MUST be <= limit

limit = 0 or negative
  → validation error

limit > documented maximum
  → validation error; no silent clamp
```

Small collections obey the same contract. `limit=1` never permits returning more than one member merely because the collection is small.

### 12.1 Exact numeric default/max — DEFER SAFELY

Exact values such as `default=50` / `max=200` are intentionally deferred to the final OpenAPI contract/tooling sub-batch after concrete list-item wire representations are serialized/measured.

Fences already fixed:

- finite positive default;
- finite positive maximum;
- one shared pair preferred when evidence supports it;
- operation-specific exceptions only for measured payload/PII/provider/consumer constraints;
- no silent clamp;
- later OpenAPI work may fill the numbers but may not redesign W3 pagination semantics.

This is a bounded `DEFER SAFELY`, not permission for unbounded collection responses.

---

## 13. Retry semantics

List/Search Product operations are GET reads and do not use `Idempotency-Key`.

A same-request retry may succeed, fail normally, or discover that its continuation expired. If expired, the client starts a new traversal explicitly.

Read retry never converts an expired cursor into a silent restart or a consequential intake concept.

---

## 14. W3-C proof / negative controls

Later contract/runtime proof must falsify at least:

1. same stable Sale/Listing/Work/SourceProduct identity returned twice in one traversal;
2. provider duplicate member passed directly through to Product client;
3. mutable `updated_at` ordering causing a previously returned history member to reappear when immutable ordering was sufficient;
4. expired cursor silently restarting at the first page;
5. `410 Gone` treating the cursor itself as a Product resource by convenience;
6. temporary provider outage being mislabeled `cursor-expired`;
7. public `snapshot_id`/`traversal_id` created only to simplify implementation;
8. exact arbitrary `50/200` defaults introduced without concrete payload/consumer evidence;
9. server returning more than caller `limit`;
10. small administrative collection ignoring caller `limit`;
11. full ListingIntent historical dispatch basis repeated in every list row by schema symmetry;
12. generic `fields/select/expand` projection language;
13. ComparableOffer receiving a fake canonical ID for deduplication;
14. equal-looking ComparableOffers collapsed as duplicates without identity evidence;
15. a newly inserted member before the continuation boundary being treated as proof pagination failed;
16. at-most-once member guarantee being interpreted as source/population completeness;
17. provider continuation failure approximated by skipping/restarting while returning success;
18. invalid/query-mismatched cursor returning an empty page instead of an explicit problem.

---

## 15. W3-C outcome

**Outcome:** `CURRENT PARENT STRUCTURE CONFIRMED` with one W3-local crystallization and one bounded deferral.

### W3-local crystallization

> **List representation may be an owner-specific bounded projection rather than the full point-resource representation when full detail/history/evidence would create unnecessary collection payload. This projection never becomes a second authority or generic projection framework.**

### Bounded deferral

> **Exact numeric `limit` default/max is `DEFER SAFELY` to final OpenAPI contract closure after list-item representations are concrete; positive finite semantics and no-silent-clamp behavior are already binding.**

No D0→D5-B1/W1/W2/W3-A/W3-B parent reopen is required.

---

## 16. Next action

Run the **Whole-W3 Global Coherence Review** over W3-A + W3-B + W3-C as one collection/query system.

The review must challenge duplicate/missing collection authority, query/cursor binding, ordering stability, provider-source coverage, deduplication versus fabricated identity, list-item projection leakage, Problem Details, limit deferral and any hidden generic query/snapshot/traversal abstraction.

Do not begin later Wire Contract obligations until Whole-W3 is coherent.

Implementation remains blocked until D9.
