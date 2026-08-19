# Marketplace Central — Architecture Rebaseline

> **Role:** sole current-program status/router after `AGENTS.md`  
> **Current phase:** **D5 — API — OPEN / ACTIVE; D5-B1 ACCEPTED / CANONICAL; D5-B2 OPEN / ACTIVE**  
> **D5-B2 current state:** **Operation Matrix + Whole-Matrix RATIFIED; Wire W1 + W2 ACCEPTED / CANONICAL; W3-A/B/C ACCEPTED IN-STAGE; Whole-W3 Global Coherence Review = NEXT**  
> **Implementation:** **BLOCKED until D9 is accepted**  
> **Last updated:** 2026-08-19

## 1. Authority path

A fresh session reads, in order:

1. `AGENTS.md`
2. this router
3. `docs/engineering/standards/root-cause-global-maximum-method.md`
4. `ARCHITECTURE.md`
5. `docs/engineering/rebaseline/DECISION-RECONCILIATION-BASELINE.md`
6. `docs/architecture/decisions/README.md`
7. `docs/engineering/rebaseline/D0-PRODUCT-SYSTEM-DEFINITION.md`
8. `docs/engineering/rebaseline/D1-DOMAINS-BOUNDARIES.md`
9. `docs/engineering/rebaseline/D2-IDENTITY-TENANT-DATA-OWNERSHIP.md`
10. `docs/engineering/rebaseline/D3-COMMUNICATION-EVENTS.md`
11. `docs/engineering/rebaseline/D4-EXTERNAL-INTEGRATIONS.md`
12. `docs/engineering/rebaseline/D4-R1-PUBLICATION-INPUT.md`
13. `docs/engineering/rebaseline/D5-API.md`
14. `docs/engineering/rebaseline/D5-B2-PRODUCT-OPERATION-SURFACE.md`
15. `docs/engineering/rebaseline/D5-B2-OPERATION-ADMISSION-MATRIX.md`
16. `docs/engineering/rebaseline/D5-B2-WIRE-CONTRACT.md` — canonical W1
17. `docs/engineering/rebaseline/D5-B2-W2-SCHEMA-GRAMMAR.md` — canonical W2
18. `docs/engineering/rebaseline/D5-B2-W3-COLLECTION-GRAMMAR.md` — W3-A/B accepted in-stage
19. `docs/engineering/rebaseline/D5-B2-W3-C-CURSOR-SAFETY.md` — W3-C accepted in-stage
20. `docs/engineering/rebaseline/EVIDENCE-REGISTER.md`
21. code/OpenAPI/schemas/tests/runtime only as current-state evidence when needed

This file alone owns **program status, allowed/blocked work and exact next action**. Detailed semantics remain in the accepted artifacts above.

W3-C is intentionally a bounded staging artifact while Whole-W3 is reviewed. If Whole-W3 converges and is operator-ratified, W3-A/B/C must be consolidated into one canonical W3 authority and the staging file removed; Git history is the archive.

`docs/engineering/rebaseline/cockpit.html` is a non-authoritative visual projection only. It never participates in this authority path.

Legacy/current code, routes, OpenAPI, SDK and frontend tables remain evidence only.

---

## 2. Program state

```text
D0 — Product / System Definition                         CLOSED / ACCEPTED
D1 — Domains / Boundaries                                CLOSED / ACCEPTED
D2 — Identity / Tenant / Data Ownership                  CLOSED / ACCEPTED
D3 — Communication / Events                              CLOSED / ACCEPTED
D4 — External Integrations                               CLOSED / ACCEPTED AS A WHOLE
D4-R1 — Publication Input & Listing Authoring            ACCEPTED / CANONICAL
Decision Reconciliation                                  ACCEPTED / CANONICAL
D5 — API                                                  OPEN / ACTIVE
  B1 Semantic API Model                                  ACCEPTED / CANONICAL
  B2 Product Operation / Resource Surface                 OPEN / ACTIVE
    B2-A Client/Auth                                     ACCEPTED IN-STAGE
    Operation Admission Matrix                           ACCEPTED / RATIFIED
    Whole-Matrix Global Coherence                        ACCEPTED / RATIFIED
    Wire Contract
      W1 Resource / Path / HTTP Grammar                  ACCEPTED / CANONICAL
      W2 Request / Response Schema Grammar               ACCEPTED / CANONICAL
      W3 Collections / Pagination / Filter / Search /
         Cursor Grammar                                  OPEN / ACTIVE
        W3-A Collection + Cursor Core                    ACCEPTED IN-STAGE
        W3-B Per-Operation Filter / Search / Ordering    ACCEPTED IN-STAGE
        W3-C Cursor Safety / Population / Limits         ACCEPTED IN-STAGE
        Whole-W3 Global Coherence                        NEXT
      Remaining wire obligations                         BLOCKED BY W3 SEQUENCE
D6 — Frontend                                             BLOCKED BY D5
D7 — Runtime / Jobs / Transactions                        BLOCKED
D8 — Golden Flows                                         BLOCKED
D9 — Adversarial Architecture Review                      BLOCKED
Implementation                                            BLOCKED UNTIL D9
```

Whole-W2 independent review is complete and incorporated into canonical W1/W2. W3 is now the only active Wire design surface.

---

## 3. Load-bearing current authority

### 3.1 W1/W2 laws that W3 may not weaken

- Product API is semantic-owner driven and Organization-scoped; protocol/provider topology is not Product ontology.
- external Product/Listing/Sale/Shipment identity remains source-qualified; no mirror IDs or bare native correlation keys.
- custom owner methods remain explicit capabilities; no generic Action/Command/Workflow surface.
- one opaque ETag revision authority per protected meaning; true same-resource standard mutation uses `If-Match`; custom/reference revision proofs use typed ETag request data.
- Idempotency-Key and revision/concurrency solve different failure classes.
- request/read schemas remain authority-separated and closed; `null` never means unknown/unavailable/partial/not-applicable.
- no universal Fact/Evidence/Result/Subject/Scope/Policy/Workflow wrappers.
- ListingIntent / PriceIntent / Availability remain distinct through publication.
- Market coverage is distinct from evidence sufficiency; Economics never fabricates conclusions from missing evidence.
- Work never becomes source-truth or generic resolution authority.
- implementation remains blocked until D9.

### 3.2 Accepted W3-A — Collection + Cursor Core

- every admitted List/Search operation is pagination-capable through shared `limit?` + opaque `cursor?` mechanics;
- responses are operation/owner-specific with semantic collection field + optional `next_cursor`; no universal Page/PagedResult/data/metadata wrapper;
- forward-only pagination; no baseline page number, offset, skip, before or previous cursor;
- `limit` is a requested maximum, not promised cardinality;
- fewer than `limit` items, including zero, does not prove exhaustion;
- only `next_cursor` indicates another traversal page;
- cursor is opaque, query/operation/Organization-bound, never authorization and never raw provider/database continuation state;
- material query changes invalidate continuation; `limit` may change;
- cursor exhaustion never proves source/provider/market/all-time completeness;
- no universal total count, caller sort or snapshot-isolation promise;
- invalid/expired/query-mismatched cursor fails explicitly;
- D4 owns provider paging protocol; D7 owns cursor persistence/signing/index/cache realization.

### 3.3 Accepted W3-B — Filter / Search / Ordering Matrix

- typed operation-specific query parameters only; no generic filter/query expression language;
- different admitted filter fields combine by AND; one value per field baseline; no OR/NOT/IN/functions/traversal DSL;
- identity point lookups stay Get operations; correlation filters exist only for real bounded navigation populations;
- `SearchSourceProductsForMarketplace` remains the only baseline Search and requires Marketplace Installation + SourceInstance + non-empty query;
- source Product search supports bounded exact identifier matching plus lexical name matching; no fuzzy/vector/AI search baseline and no public relevance score;
- caller-selectable sort baseline count = **0**;
- each collection has one owner-defined deterministic default order plus stable tie-breaker where required;
- external identity filters remain source-qualified and same-Organization references fail closed;
- collection-level coverage is exposed only for source-backed populations whose enumeration may be incomplete: source Product search, Marketplace Listings, Comparable Offers, Marketplace Sales, Shipments and Sale Economics where underlying coverage is not complete;
- no filter may expose provider status, Sankhya native columns/TOP/NUNOTA/CODEMP, provider JSON paths or frontend-table fields by convenience.

The accepted operation-specific matrix is authoritative in `D5-B2-W3-COLLECTION-GRAMMAR.md`.

### 3.4 Accepted W3-C — Cursor Safety / Population / Limits

- List operations may use bounded owner-specific list-item projections instead of full point-Get representations when full history/evidence would be excessive; no generic `fields/select/expand` projection DSL;
- `400 invalid-cursor` covers malformed/tampered/unknown/operation/Organization/query-mismatched cursors after normal access/privacy checks;
- `400 cursor-expired` covers legitimately issued continuation that can no longer be resumed honestly;
- cursor is an ephemeral continuation, not a durable bookmark; no public minimum TTL baseline;
- no silent cursor restart or best-effort nearby continuation;
- stable MPC/source-qualified member identity appears at most once per Product traversal;
- at-most-once member identity does not imply snapshot or completeness;
- inserts/deletes/filter changes between pages may affect not-yet-traversed population; a fresh current universe requires a new traversal;
- updates do not turn pagination into a change feed;
- ComparableOffer/non-identifiable evidence receives no fabricated canonical ID merely for deduplication;
- provider continuation expiry is reconstructed transparently only when identical Product continuation semantics can be proven; otherwise `cursor-expired`;
- transient provider outage is not cursor expiration by itself;
- no public traversal/snapshot/search-session resource baseline;
- `limit` must be a positive finite requested maximum and the server never returns more than it; zero/negative/above-max is validation error, with no silent clamp;
- exact numeric default/max is `DEFER SAFELY` to final OpenAPI closure after concrete list-item payloads are known, under already-fixed finite/no-silent-clamp fences.

---

## 4. Prohibited now

While Whole-W3 Global Coherence is next:

- do not begin remaining Wire Contract obligations, D6–D9 target design or implementation;
- do not reopen accepted D0→D5-B1/B2/W1/W2/W3-A/B/C for naming/style or implementation convenience;
- do not derive collection behavior from legacy routes/OpenAPI/controllers/frontend tables;
- do not introduce generic query/filter/sort/projection DSL, arbitrary caller sort, generic priority/tags/search fields or provider-native query vocabulary;
- do not make cursor exhaustion, at-most-once identity or traversal completion imply source/population completeness;
- do not invent universal total counts, snapshot isolation or durable traversal resources;
- do not fabricate identity for evidence values solely to stabilize pagination;
- do not select D7 cursor storage/signing/database mechanism during W3 review;
- do not choose exact OpenAPI minor/generator or numeric limit defaults before Whole-W3 coherence.

---

## 5. Exact next action

**Run the Whole-W3 Global Coherence Review over W3-A + W3-B + W3-C as one collection/query system.**

The review must adversarially stress-test at least:

1. every admitted List/Search operation has one honest collection home and no list-by-symmetry operation slipped in;
2. collection response/list-item projections do not create second read/business authority or generic projection language;
3. cursor semantic-query binding is coherent with every W3-B filter/search field and `limit` exception;
4. ordering + tie-breaker can support the at-most-once guarantee without forcing fake Product identity;
5. mutable populations do not accidentally imply snapshot, no-omission or completeness guarantees;
6. source/provider coverage remains independent from traversal exhaustion/deduplication;
7. provider paging/continuation failure stays D4/D7 mechanism while Product cursor semantics remain stable;
8. `invalid-cursor` / `cursor-expired` are sufficient and do not collide with access/validation/service/business problems;
9. exact numeric limit deferral is bounded enough that later OpenAPI work cannot silently redesign collection semantics;
10. no generic filter/sort/search/projection/traversal/snapshot framework emerges from the combined package;
11. Structural Inversion against legacy OpenAPI/frontend/database/provider paging still passes;
12. no D0→W2 parent reopen is actually required.

If no material contradiction survives, present the complete Whole-W3 package to the operator for final ratification before consolidating W3 into one canonical artifact.

If a material contradiction survives, return only the smallest implicated W3-local or parent scope to the Decision Loop; do not proceed to later Wire obligations.

Implementation remains blocked until D9.

---

## 6. Fresh-session success test

A fresh session must conclude:

- D0→D4/D4-R1 + D5-B1 are accepted/canonical;
- D5-B2 Operation Matrix + Whole-Matrix are ratified;
- W1 and W2 are canonical;
- W3-A/B are accepted in-stage in `D5-B2-W3-COLLECTION-GRAMMAR.md`;
- W3-C is accepted in-stage in `D5-B2-W3-C-CURSOR-SAFETY.md`;
- W3-A defines owner-named forward opaque cursor traversal without total/count/snapshot/completeness fiction;
- W3-B defines the bounded per-operation typed filter/search/order matrix with zero caller-selectable sort baseline;
- W3-C defines list-item projection, cursor failure/lifetime/population/deduplication/limit semantics without durable traversal resource or fake identity;
- **Whole-W3 Global Coherence Review is the exact next action**;
- implementation remains blocked until D9.

If not, the active authority tree is inconsistent.
