# Marketplace Central — Architecture Rebaseline

> **Role:** sole current-program status/router after `AGENTS.md`  
> **Current phase:** **D5 — API — OPEN / ACTIVE; D5-B1 ACCEPTED / CANONICAL; D5-B2 OPEN / ACTIVE**  
> **D5-B2 current state:** **Operation Matrix + Whole-Matrix RATIFIED; Wire W1 + W2 ACCEPTED / CANONICAL; W3-A Collection + Cursor Core ACCEPTED IN-STAGE; W3-B Per-Operation Filter / Search / Ordering Matrix = NEXT**  
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
17. `docs/engineering/rebaseline/D5-B2-W2-SCHEMA-GRAMMAR.md` — canonical consolidated W2
18. `docs/engineering/rebaseline/D5-B2-W3-COLLECTION-GRAMMAR.md` — W3-A accepted in-stage; current W3 design home
19. `docs/engineering/rebaseline/EVIDENCE-REGISTER.md`
20. code/OpenAPI/schemas/tests/runtime only as current-state evidence when needed

This file alone owns **program status, allowed/blocked work and exact next action**. Detailed semantics remain in the accepted artifacts above.

The former staged W2-C/D/E files and Whole-W2 review candidate were removed after final ratification; Git history is the archive. `AI-DIALOG.md` is protocol-only and is not architecture authority.

Legacy/current code, routes, OpenAPI and SDK remain evidence only until later D-stage target work replaces them.

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
        W3-B Per-Operation Filter / Search / Ordering    NEXT
      Remaining wire obligations                         BLOCKED BY W3 SEQUENCE
D6 — Frontend                                             BLOCKED BY D5
D7 — Runtime / Jobs / Transactions                        BLOCKED
D8 — Golden Flows                                         BLOCKED
D9 — Adversarial Architecture Review                      BLOCKED
Implementation                                            BLOCKED UNTIL D9
```

Whole-W2 lead review, Fable Round 1, GPT adjudication, focused Fable Round 2 and final GPT adjudication converged and were operator-ratified on 2026-08-19. Their active meaning is incorporated in canonical W1/W2; review dialogue itself is not an authority layer.

---

## 3. Load-bearing D5-B2 authority

### 3.1 Semantic/client authority

- Product API is semantic-owner driven, not CRUD-, screen- or provider-shaped;
- humans use OIDC Authorization Code + PKCE; confidential machines use Client Credentials/service-account semantics;
- MPC owns Principal, Membership, AccessRole/Permission and every business authority;
- `GET /access-context` is the only platform-scoped self-only Product Q baseline;
- Organization business API remains `/organizations/{organization_id}/...`;
- no Product/PIM, generic Integration/Action/Operation/Workflow/Rules/Task/Finance/AI owner;
- zero-P baseline remains until D6 consumer evidence proves a bounded composition need.

### 3.2 Operation authority

- ListingIntent, PriceIntent and Availability remain distinct through initial publication;
- BusinessOrderIntent/InvoicingIntent remain Materialization owner reactions, not direct Product commands;
- Governance authorizes but does not execute;
- Fulfillment physical facts require admitted human/proven physical-system authority and cannot be fabricated by ordinary automation;
- Post-Sale coordinates scoped consequences and does not gain provider-action vocabulary;
- Work owns responsibility/lifecycle, never source truth;
- every admitted consequential C operation has an explicit idempotency/current-state disposition.

### 3.3 Canonical W1

- no `/v1` compatibility axis without a real stable consumer;
- paths express Organization/canonical identity or source namespace, not D1 package names/workflow order;
- external Listing/Sale/Shipment keep Installation-qualified native identity; no mirror IDs;
- standard methods only when honest; owner capabilities use `POST {resource-or-keyed-subject-uri}:verb`;
- custom-method URI is not an implicit alias of the base resource for HTTP conditionals;
- one opaque ETag revision authority per protected meaning;
- true same-resource standard mutation uses `If-Match` (`428` missing / `412` false);
- current-state-protected custom method carries typed `etag` request data (`422` invalid/missing / `409 resource-revision-conflict` stale);
- exact revision of another resource carries ETag adjacent to the typed reference;
- ProductChannelCorrespondence remains keyed Readiness meaning with a correspondence-scoped ETag and Resolve/Clear capabilities; no synthetic Correspondence ID or forced PUT/DELETE;
- Idempotency-Key remains duplicate-intake safety, independent of revision proof;
- provider/business-system protocol ingress remains outside Product API roots.

### 3.4 Canonical W2

- opaque IDs/typed source refs; exact decimal Money; explicit temporal meanings;
- request/read schemas are authority-separated and closed;
- `null` never carries unknown/unavailable/partial/not-applicable;
- knowledge uses smallest owner-specific unions, not universal `Fact<T>`;
- no generic Result/Evidence/ExternalRef/Subject/Scope/Policy/Workflow wrappers;
- ListingIntent is sparse/declarative create/edit identity; historical dispatch/effect basis is append-only and preserves attempt-time resolved values/provenance without a PublicationAttempt resource;
- `PublicationValue` includes bounded `number_unit` and explicit requirement-permitted `not_applicable`, without generic UoM engine;
- MarketplaceListing is source-qualified actual-state evidence; Listing/Price convergence stay in their Intents;
- PriceIntent remains separate and explicit supersession carries prior PriceIntent revision proof;
- Availability separates control, desired meaning, provider observation and convergence;
- Readiness owns requirement definitions/candidates; draft-dependent provider requirement evaluation is D4 technical evidence feeding Offering dispatchability while Price/Availability stay separate;
- Market coverage remains distinct from evidence sufficiency;
- Economics is components-first, preserves L0/L1/L2 and R1/R2, and never fabricates profitability from missing evidence;
- EconomicPerformanceSummary is period/scope keyed with explicit coverage and no generic metrics map;
- AuthorizationDecision is immutable occurrence; AuthorizationDelegation has one justified stable ID;
- Sale remains external; `sale_line_key`, once minted, never rebinds;
- Party/Destination are bounded Materialization meanings, not Customer/Address masters;
- `FulfillmentExecutionId` is the one durable Fulfillment lifecycle identity; `FulfillmentState` is not another resource;
- FulfillmentNode is minimal MPC identity distinct from InventorySource/native warehouse;
- one baseline Fulfillment internal target only: `dispatch_handoff_lead_time_before_provider_deadline`;
- Post-Sale consequence tracks remain independent; Work generic resolution remains deferred and closure audit passes;
- media intake is ListingIntent-bound multipart `:create-media`, `200` descriptor, same media identity on exact replay, binary content in idempotency fingerprint;
- exact idempotent intake resolves before rechecking a revision advanced by the first successful call;
- canonical Problem Details includes unified `resource-revision-conflict` for stale typed revision proof;
- D4/D8 still must prove selected-lane N/A reread and User-Product conditional-requirement validation before claiming live convergence.

### 3.5 Accepted W3-A

- every admitted List/Search operation is pagination-capable through shared `limit?` + opaque `cursor?` mechanics;
- collection responses are operation/owner-specific and may expose optional `next_cursor`; no universal `Page<T>`, `PagedResult<T>`, `data/metadata` or generic Result wrapper;
- pagination is forward-only; no baseline page number/offset/skip/previous cursor;
- `limit` is a requested maximum, not guaranteed returned cardinality;
- fewer than `limit` items, including zero, does not prove exhaustion; `next_cursor` is the continuation authority;
- absence of `next_cursor` means no later page in that traversal only; it never proves source/provider/market/all-time knowledge completeness;
- cursor is opaque, bound to Organization + operation + semantic query, never authorization and never a raw provider/database paging token;
- material query changes invalidate continuation; `limit` may change between pages;
- no universal `total_count` baseline;
- no caller-selectable arbitrary sort baseline;
- no universal snapshot-isolation promise across pages;
- malformed/invalid/expired/query-mismatched cursor fails explicitly rather than silently restarting or returning an empty first page;
- provider pagination remains D4-local mechanism; D7 chooses cursor persistence/signing/index/cache implementation.

---

## 4. Prohibited now

While W3-B is next:

- do not begin D6–D9 target design or implementation;
- do not reopen accepted D0→D5-B1/B2/W1/W2/W3-A for naming/style or implementation convenience;
- do not reconstruct deleted staged W2 files or review candidate as parallel authority;
- do not derive collection contracts from legacy OpenAPI/routes/controllers/frontend tables;
- do not create generic query/filter DSL, arbitrary sort expressions or GraphQL-like query language;
- do not make cursor exhaustion imply universal source/provider completeness;
- do not invent total counts, snapshot isolation or stable ordering stronger than an owner/source can justify;
- do not add bulk operations by symmetry;
- do not introduce D7 paging cache/storage/database cursor mechanics during W3;
- do not choose exact OpenAPI minor/generator purely by recency.

---

## 5. Exact next action

**Derive W3-B — Per-Operation Filter / Search / Ordering Matrix from every admitted List/Search Q.**

W3-B must decide and adversarially stress-test:

1. whether each collection requires filtering at all;
2. the smallest operation-specific typed filter fields, rejecting generic filter maps/expressions;
3. true search semantics separately from ordinary list/filter semantics, especially `SearchSourceProductsForMarketplace`;
4. one deterministic owner-meaning default order per collection plus required stable tie-breaker(s), without exposing database ordering as API meaning;
5. whether any collection genuinely requires caller-selectable sorting; reject by default;
6. same-Organization/source-qualified identity constraints in query selection;
7. owner-specific coverage/provenance exposure where material;
8. whether small product-defined enumerations/definition lists should simply return all values while retaining the common pagination-capable contract;
9. negative controls against provider-native/database-field filters, frontend-table-shaped queries, fuzzy search where not semantically defined, and sort/filter parameters admitted only for symmetry.

After W3-B, close cursor invalid/expired/stale plus population-change/deduplication/stability semantics as required by the per-operation matrix before W3 can be accepted as a whole.

After W3 as a whole, continue remaining Wire Contract obligations in router order: exact Permission→operation/client-class mapping, technical non-Product ingress classification, final Problem/media consistency as needed, and the single machine-readable OpenAPI authority/tooling decision. Do not advance them before W3 is coherent.

Implementation remains blocked until D9.

---

## 6. Fresh-session success test

A fresh session must conclude:

- D0→D4/D4-R1 + D5-B1 are accepted/canonical;
- D5-B2 Operation Matrix + Whole-Matrix are ratified;
- `D5-B2-WIRE-CONTRACT.md` is canonical W1 authority;
- `D5-B2-W2-SCHEMA-GRAMMAR.md` is the single consolidated canonical W2 authority;
- `D5-B2-W3-COLLECTION-GRAMMAR.md` is the current W3 design home with W3-A accepted in-stage;
- former W2-C/D/E staging artifacts and Whole-W2 review candidate are absent from the active tree and preserved only in Git history;
- `AI-DIALOG.md` has no active Whole-W2 review cycle;
- W1 custom-method preconditions use typed request ETag rather than base-resource `If-Match` on a distinct URI;
- W2 includes the ratified historical publication, missing-schema, provider-requirement, Fulfillment identity/target, key-lifetime, media, idempotency and problem corrections;
- W3-A establishes named owner collections + forward opaque cursor semantics without total/count/snapshot/completeness fiction;
- **W3-B Per-Operation Filter / Search / Ordering Matrix is the exact next action**;
- implementation remains blocked until D9.

If not, the active authority tree is inconsistent.
