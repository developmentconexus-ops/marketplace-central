# Marketplace Central — Architecture Rebaseline

> **Role:** sole current-program status/router after `AGENTS.md`  
> **Current phase:** **D5 — API — OPEN / ACTIVE; D5-B1 ACCEPTED / CANONICAL; D5-B2 OPEN / ACTIVE**  
> **D5-B2 current state:** **Operation Matrix + Whole-Matrix RATIFIED; Wire W1 + W2 ACCEPTED / CANONICAL; W3-A + W3-B ACCEPTED IN-STAGE; W3-C Cursor Validity / Population Change / Deduplication / Limits / Problem Grammar = NEXT**  
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
18. `docs/engineering/rebaseline/D5-B2-W3-COLLECTION-GRAMMAR.md` — current W3 authority/design home
19. `docs/engineering/rebaseline/EVIDENCE-REGISTER.md`
20. code/OpenAPI/schemas/tests/runtime only as current-state evidence when needed

This file alone owns **program status, allowed/blocked work and exact next action**. Detailed semantics remain in the accepted artifacts above.

Former staged W2-C/D/E artifacts and the Whole-W2 review candidate are absent from the active tree; Git history is the archive. `AI-DIALOG.md` is protocol-only and is not architecture authority.

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
        W3-C Cursor Validity / Population Change /
             Deduplication / Limits / Problem Grammar   NEXT
        Whole-W3 Global Coherence                        AFTER W3-C
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

The accepted operation-specific matrix is authoritative in `D5-B2-W3-COLLECTION-GRAMMAR.md`, including bounded filters/order for ListingIntents, PriceIntents, Availability, Market, Economics, Governance, Sales, Materialization, Fulfillment, Post-Sale and Work.

---

## 4. Prohibited now

While W3-C is next:

- do not begin D6–D9 target design or implementation;
- do not reopen accepted D0→D5-B1/B2/W1/W2/W3-A/B for naming/style or implementation convenience;
- do not derive collection behavior from legacy routes/OpenAPI/controllers/frontend tables;
- do not introduce generic query/filter DSL, arbitrary caller sort, generic priority/tag/search fields or provider-native query vocabulary;
- do not make cursor exhaustion imply knowledge/source completeness;
- do not invent universal total counts, snapshot isolation or stable-traversal guarantees stronger than owner/source evidence;
- do not add bulk operations by symmetry;
- do not choose D7 cursor storage/signing/database mechanism during W3;
- do not choose OpenAPI minor/generator before W3 coherence.

---

## 5. Exact next action

**Derive W3-C — Cursor Validity / Population Change / Deduplication / Limits / Problem Grammar.**

W3-C must decide and adversarially stress-test:

1. malformed cursor versus valid-but-query-mismatched versus expired/invalidated continuation;
2. exact Problem Details/status grammar for those failure classes;
3. whether Product cursor has a stable public lifetime promise or only implementation-bounded validity;
4. insert/update/delete behavior between page requests without claiming universal snapshot isolation;
5. duplicate/omission guarantees a Product client may rely on, per collection class;
6. source-backed/provider traversal behavior when the provider population or paging token changes/expires mid-traversal;
7. whether any collection needs a bounded traversal/snapshot identifier beyond the cursor; reject by default;
8. default `limit`, accepted range and maximum — shared where honest, operation-specific only when a real scale/PII/payload constraint requires it;
9. whether very small definition/configuration collections should ignore caller `limit` by returning all values or still obey the requested maximum consistently;
10. restart/recovery semantics after stale/expired cursor: explicit new traversal, never silent first-page fallback;
11. deduplication/stability for tie-breakers that are not Product identities (e.g. ComparableOffer evidence continuation);
12. final W3 negative controls and reopen triggers.

After W3-C, run a **Whole-W3 Global Coherence Review** before accepting W3 as a whole. Do not begin remaining Wire Contract obligations before W3 coherence.

After accepted W3, continue router-ordered Wire obligations: exact Permission→operation/client-class mapping; technical non-Product ingress classification; final Problem/media consistency as needed; and the single machine-readable OpenAPI authority/tooling decision.

Implementation remains blocked until D9.

---

## 6. Fresh-session success test

A fresh session must conclude:

- D0→D4/D4-R1 + D5-B1 are accepted/canonical;
- D5-B2 Operation Matrix + Whole-Matrix are ratified;
- W1 and W2 are canonical;
- `D5-B2-W3-COLLECTION-GRAMMAR.md` is the active W3 authority/design home with **W3-A + W3-B accepted in-stage**;
- W3-A defines owner-named forward opaque cursor traversal without total/count/snapshot/completeness fiction;
- W3-B defines the bounded per-operation typed filter/search/order matrix with zero caller-selectable sort baseline;
- **W3-C Cursor Validity / Population Change / Deduplication / Limits / Problem Grammar is the exact next action**;
- implementation remains blocked until D9.

If not, the active authority tree is inconsistent.
