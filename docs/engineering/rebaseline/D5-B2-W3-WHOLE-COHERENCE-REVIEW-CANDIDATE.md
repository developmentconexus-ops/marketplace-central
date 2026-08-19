# D5-B2 — Whole-W3 Global Coherence Review Candidate

> **Status:** NON-AUTHORITATIVE LEAD REVIEW CANDIDATE  
> **Review subject:** accepted-in-stage W3-A/B/C as one Product collection/query system  
> **Authority:** none — review evidence only until operator ratification and canonical filing  
> **Method:** DevelopmentConexus Engineering Method v1.0.0  
> **Reviewed:** 2026-08-19

## 1. Review outcome

**Method outcome:** `RESTRUCTURE NOW — W3-LOCAL`.

The parent architecture survives. No current finding requires a D0→D5-B1/W1/W2 semantic reopen.

The combined W3 package is globally sound in its main direction — owner-named collections, forward opaque continuation, operation-specific typed filters, zero caller-selectable sort, no universal snapshot/total/query language — but six W3-local contradictions/gaps must be corrected before W3 can be accepted as a whole.

## 2. What was challenged and passed

- admitted List/Search operation coverage: **26/26 have a W3 home**;
- list-by-symmetry operations introduced after the ratified operation matrix: **0**;
- universal `Page<T>` / `PagedResult<T>` / generic metadata wrapper: **rejected / PASS**;
- offset/page-number pagination baseline: **rejected / PASS**;
- caller-selectable sort baseline: **0 / PASS**;
- generic filter/OData/SQL/GraphQL-like expression language: **rejected / PASS**;
- raw provider/database cursor leakage: **rejected / PASS**;
- cursor as authorization: **rejected / PASS**;
- traversal exhaustion => source/provider completeness: **rejected / PASS**;
- universal total count: **rejected / PASS**;
- universal snapshot/traversal resource: **rejected / PASS**;
- D7 cursor/storage implementation leaking into W3 authority: **not required / PASS**;
- stable-identity at-most-once member guarantee: **survives**, because W3-C already permits internal keyset/traversal state and explicit `cursor-expired` rather than fabricating public snapshot/identity when continuity cannot be preserved;
- exact numeric limit default/max deferral: **bounded / PASS**, because positive finite semantics + no-silent-clamp + later measurement owner are already fixed.

## 3. W3-G1 — continuation request query-carrier contradiction — REVISE

### Evidence / contradiction

W3-A says a cursor is bound to the same semantic query and says the client **may repeat** query fields. W3-B simultaneously makes some semantic selectors mandatory, especially:

```text
SearchSourceProductsForMarketplace
  marketplace installation
  source instance
  non-empty query
```

and `ListComparableOffers` requires a MarketSubject.

If continuation may omit those selectors because the cursor silently carries them, the cursor becomes a hidden query/selection authority. If those selectors remain mandatory, “may repeat” is false/ambiguous.

### Root cause

Continuation state and semantic query identity were not separated strongly enough.

### Corrected invariant / direction

> **A Product cursor carries continuation state only; it never replaces the explicit semantic query. Every continuation request repeats the operation-required subject/search parameters and the same effective optional narrowing semantics. Only `limit` may vary.**

Rules:

- path Organization remains explicit;
- required operation subject/query fields remain required on every page request;
- any optional filter selected on page 1 must have the same effective value on continuation;
- optional filter omitted on page 1 remains semantically omitted unless the contract explicitly defines an equivalent default;
- cursor verifies the repeated semantic query; it does not become the only place that query exists;
- `limit` is excluded from semantic-query equality and may change as already accepted.

### Alternatives

1. cursor reconstructs the query and fields become optional — convenient, but hides business selection inside opaque continuation;
2. both styles allowed — more states and ambiguous precedence;
3. **repeat explicit query + cursor only for continuation — selected Global Maximum.**

### Parent reopen

None. W3-local clarification only.

---

## 4. W3-G2 — undefined enumerable population for keyed Q lists — REVISE

### Evidence / contradiction

W2 deliberately allows:

```text
MarketSubject
  existing_listing
  source_product_marketplace_context

EconomicsSubject
  existing_listing
  source_product_marketplace_context
```

so `GetCompetitivePosition` / `GetExpectedEconomics` can reason before a listing exists.

W3-B then admits:

- `ListCompetitivePositions` over a generic stable MarketSubject tuple;
- `ListExpectedEconomics` over a generic stable EconomicsSubject tuple.

But MPC explicitly has **no Product master/list universe**. Interpreting those Lists as “all possible source Product × marketplace contexts” would recreate the very Product-list/master surface D0/D2/D4-R1 rejected.

### Root cause

Point-addressable keyed Q subject space was mistaken for an enumerable collection population.

### Corrected invariant / direction

Baseline List populations are narrower than point-Q subject unions:

- `ListCompetitivePositions` enumerates CompetitivePosition meanings for **currently known existing marketplace Listing subjects** only;
- `ListExpectedEconomics` enumerates ExpectedEconomics meanings for **currently known existing marketplace Listing subjects** only;
- pre-listing source-product marketplace contexts remain reachable by explicit `SearchSourceProductsForMarketplace` → selected Product ref → point `GetCompetitivePosition` / `GetExpectedEconomics` or `EvaluatePriceScenario` as admitted;
- no hidden “all source Products” enumeration is introduced.

Their default ordering therefore becomes stable source-qualified Listing reference ordering, optionally narrowed by the already-admitted Marketplace Installation filter.

### Alternatives

1. enumerate all possible source Product contexts — rejected: resurrects Product master/list universe;
2. enumerate whichever keyed Qs happen to be cached/materialized — rejected: persistence state becomes Product collection semantics;
3. **list existing marketplace-offer/listing subjects; point-query pre-listing contexts — selected.**

### Parent reopen

None. W3 narrows collection population without changing W2 point-Q subject capability.

---

## 5. W3-G3 — deterministic ordering/tie-breaker underspecification — REVISE

### Evidence

W3-B requires one deterministic ordering + tie-breaker, but several rows remain placeholders rather than executable wire meaning:

- `ListSellableAvailability` → “stable target tuple”;
- `ListCompetitivePositions` → “stable MarketSubject tuple”;
- `ListExpectedEconomics` → “stable EconomicsSubject tuple”;
- `ListEconomicAttributions` → “stable attribution-subject tuple”;
- `ListShipments` → “stable source-qualified Shipment order”.

At least one is insufficiently unique: multiple EconomicAttribution records may concern the same semantic subject, while W2 exposes them as mutable persistent resources with point resolution semantics.

### Corrected direction

Use the smallest explicit stable keys already present in accepted meaning:

- `ListSellableAvailability`: fixed closed target-kind order, then target identity:
  1. `pre_creation_listing_intent` → ListingIntentId;
  2. `existing_listing` → source-qualified Listing ref;
- `ListCompetitivePositions`: after W3-G2, source-qualified Listing ref;
- `ListExpectedEconomics`: after W3-G2, source-qualified Listing ref;
- `ListEconomicAttributions`: stable `economic_attribution_id ASC` as final ordering key rather than non-unique subject-only tuple;
- `ListShipments`: Installation-qualified native Shipment key ASC within the already source-qualified namespace.

Other accepted time + ID/ref orderings remain unchanged.

Opaque IDs/native keys are used only as deterministic stable order/tie-breaker, not business priority/ranking.

### Parent reopen

None.

---

## 6. W3-G4 — ComparableOffer price ordering needs one resumable evaluation basis — REVISE

### Evidence / contradiction

`ComparableOffer` deliberately has no synthetic MPC ID when provider evidence has no stable member identity.

W3-B orders comparable offers by mutable `delivered_price ASC` plus an acquisition-local discriminator. If each page independently re-queries a changing market population, price movement plus missing identity can reorder members across the boundary and make deduplication impossible without inventing identity.

### Corrected invariant / direction

> **One `ListComparableOffers` cursor chain is bound to one Market Intelligence evaluation/acquisition basis.**

- response/provenance identifies the owner evaluation basis sufficiently for honesty;
- D7 may retain that basis internally; no public Snapshot/Traversal resource is created;
- `delivered_price ASC` is evaluated inside that basis;
- a stable provider-native member identity may be used where real;
- otherwise an acquisition-local discriminator remains technical only;
- if the same basis cannot be resumed, return `cursor-expired` rather than re-querying current market and pretending it is continuation;
- starting without cursor requests a fresh current evaluation/traversal.

This is a collection-specific stability guarantee, not a universal W3 snapshot promise.

### Parent reopen

None; it applies existing Market evidence/provenance authority.

---

## 7. W3-G5 — Source Product search matching is stronger than proven D4 capability — REVISE

### Evidence

W3-B freezes “case-insensitive lexical token match” over source Product name/display evidence.

D4 proves sanctioned bounded Product reads and provider-local criteria/paging quirks, but does not establish a cross-SourceInstance tokenizer/case-insensitive lexical-search contract. Forcing that exact algorithm could silently require a new MPC Product search mirror/index solely to honor W3 wording.

### Root cause

A useful Product Search operation was over-specified into an implementation/search-engine guarantee not required by the real consumer.

### Corrected direction

`SearchSourceProductsForMarketplace` keeps:

- required SourceInstance + Marketplace Installation + non-empty opaque user `query`;
- source-qualified results only;
- bounded matching over legitimate source identification/display evidence supported by the sanctioned source contract;
- exact native Product key / exact legitimate identifier matches, when established, rank ahead of textual matches;
- textual matching/ranking must remain deterministic for one traversal basis but does **not** promise a universal tokenizer, stemming, fuzzy, vector or case algorithm before source capability proves it;
- provider query syntax never crosses the Product API;
- if the sanctioned source cannot satisfy a materially required search without an MPC Product mirror/index that changes architecture, STOP/re-adjudicate rather than silently creating the mirror.

No public relevance score.

### Parent reopen trigger

Only if a real SourceInstance cannot support the admitted Product search through sanctioned reads/projections without creating a materially new Product-search data authority. Current evidence does not prove that trigger.

---

## 8. W3-G6 — cursor Problem Details must have one canonical catalog — REVISE

### Evidence / contradiction

Canonical W2 currently declares the small Product Problem Details family list and does not include:

- `invalid-cursor`;
- `cursor-expired`.

W3-C justifiably introduces both because client recovery differs: malformed/mismatched cursor is not the same case as a legitimately issued continuation that can no longer resume.

Leaving them only in W3 would create two active catalogs of Product problem types.

### Corrected direction

Upon Whole-W3 final ratification/canonical consolidation:

- amend canonical W2 Problem Details catalog to include `invalid-cursor` and `cursor-expired`;
- W3 remains authority for when these two collection-specific types apply;
- both keep HTTP 400 under the ratified W3-C direction;
- do not create `cursor-stale`, `cursor-conflict`, `cursor-gone`, provider-specific cursor errors or another taxonomy;
- ordinary auth/access/service/coverage/business meanings remain in their existing problem/semantic homes.

This is one catalog with later W3 additions, not W2 vs W3 competing authorities.

### Parent reopen

None. W2 artifact amendment is a D5 wire-consistency correction, not a semantic parent reopen.

---

## 9. Strong adversarial challenge that survived

### At-most-once stable identity

Challenge: the guarantee may force stateful seen-sets for mutable/source-backed traversals.

Disposition: **KEEP**.

Reason:

- immutable-order MPC collections can usually realize it cheaply with keyset-like continuation;
- source/provider collections may rely on provider-stable continuation or bounded internal traversal state;
- W3-C already allows explicit cursor expiration when identical continuation cannot be maintained;
- the guarantee materially prevents provider duplicate paging from becoming Product duplicate membership and does not require a public snapshot resource.

Reopen in D7 only if proof shows the guarantee has disproportionate unavoidable cost for a concrete collection; D7 may not silently weaken the Product contract.

### ListItem projection law

Challenge: could become generic “view” framework.

Disposition: **KEEP WITH FENCE**.

- projection stays owner-specific/read-only;
- point identity/correlation remains the same;
- no `fields/select/expand` language;
- final OpenAPI may instantiate named list-item schemas but may not invent a second business conclusion/authority;
- exact numeric limits may be finalized only after those wire shapes are concrete.

### Limit numeric deferral

Disposition: **KEEP / DEFER SAFELY**.

The deferral has a named later owner (final OpenAPI contract), finite positive bounds, no-silent-clamp law and a measured-payload prerequisite. It cannot silently become unbounded.

## 10. Global Maximum after corrections

```text
explicit semantic query on every continuation request
+ opaque continuation-only cursor
+ owner-defined enumerable populations
+ deterministic stable ordering/tie-breakers
+ operation-specific evaluation basis only where identity/order requires it
+ source-capability-backed Search semantics
+ one Product Problem Details catalog
```

while preserving:

```text
no generic Page<T>
no query/filter/sort DSL
no universal snapshot
no total_count baseline
no Product/PIM mirror
no fabricated ComparableOffer identity
no D7 mechanism authority
```

## 11. Reopen classification

- D0 product boundary: **NO REOPEN**
- D1 semantic owners: **NO REOPEN**
- D2 identity: **NO REOPEN**
- D3 communication: **NO REOPEN**
- D4 integration semantics: **NO REOPEN**; only existing Search capability proof/reopen trigger retained
- D4-R1: **NO REOPEN**
- D5-B1/W1: **NO REOPEN**
- W2: **artifact amendment required only for canonical cursor problem catalog after final ratification**
- W3-A/B/C: **targeted W3-local corrections required**

## 12. Operator decision requested

Ratify/revise W3-G1…W3-G6 as the lead Whole-W3 correction direction.

Until operator decision:

- W3-A/B/C accepted artifacts remain current in-stage authority;
- this file is review evidence only;
- do not silently apply findings to W3/W2;
- do not begin later Wire Contract obligations.
