# Marketplace Central — Architecture Rebaseline

> **Role:** sole current-program status/router after `AGENTS.md`  
> **Current phase:** **D5 — API — OPEN / ACTIVE; D5-B1 ACCEPTED / CANONICAL; D5-B2 OPEN / ACTIVE**  
> **D5-B2 current state:** **Operation Matrix + Whole-Matrix RATIFIED; Wire W1 + W2 ACCEPTED / CANONICAL; W3 Collections / Pagination / Filter / Search / Cursor Grammar = NEXT**  
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
18. `docs/engineering/rebaseline/EVIDENCE-REGISTER.md`
19. code/OpenAPI/schemas/tests/runtime only as current-state evidence when needed

This file alone owns **program status, allowed/blocked work and exact next action**. Detailed semantics remain in the accepted artifacts above.

The former staged W2-C/D/E files and Whole-W2 review candidate were removed after final ratification; Git history is the archive. `AI-DIALOG.md` is again protocol-only and is not architecture authority.

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
         Cursor Grammar                                  NEXT
      Remaining wire obligations                         BLOCKED BY W3 SEQUENCE
D6 — Frontend                                             BLOCKED BY D5
D7 — Runtime / Jobs / Transactions                        BLOCKED
D8 — Golden Flows                                         BLOCKED
D9 — Adversarial Architecture Review                      BLOCKED
Implementation                                            BLOCKED UNTIL D9
```

Whole-W2 lead review, Fable Round 1, GPT adjudication, focused Fable Round 2 and final GPT adjudication converged and were operator-ratified on 2026-08-19. Their active meaning is now incorporated in canonical W1/W2; review dialogue itself is not an authority layer.

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

---

## 4. Prohibited now

While W3 is next:

- do not begin D6–D9 target design or implementation;
- do not reopen accepted D0→D5-B1/B2/W1/W2 for naming/style or implementation convenience;
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

**Derive W3 — Collections / Pagination / Filter / Search / Cursor Grammar from the admitted List/Search Q operations.**

W3 must decide and adversarially stress-test, as one coherent package:

1. canonical collection response shape without a universal business Result/metadata envelope;
2. opaque cursor semantics and how a client obtains the next page;
3. deterministic owner-meaning ordering and necessary tie-breakers without exposing database ordering as API meaning;
4. cursor continuation versus source/provider coverage/completeness — page exhaustion must never fabricate universal completeness;
5. page-size/limit defaults and maximums only to the extent a real Product/tooling need requires them;
6. list filtering as **operation-specific typed filters**, not a generic query DSL;
7. search semantics separately from list/filter semantics where the admitted operation is genuinely search (`SearchSourceProductsForMarketplace` in particular);
8. whether any admitted collection genuinely needs caller-selectable sort; reject arbitrary sort by default;
9. total-count semantics — include only where a real consumer and honest authoritative universe make it reliable; otherwise omit;
10. cursor invalid/expired/stale behavior and Problem Details without promising one universal snapshot transaction across owners/providers;
11. deduplication/stability expectations when external source populations change between pages;
12. source-qualified external collection identities and same-Organization reference safety;
13. provider-specific paging tokens remaining adapter-local unless an opaque Product cursor legitimately encapsulates them;
14. negative controls for unbounded page size, offset-by-default, raw provider cursor leakage, generic filter/sort maps and list-by-symmetry operations not admitted in B2.

W3 must not choose D7 database/index/cache implementation. It defines Product wire semantics only.

After W3, continue the remaining Wire Contract obligations in router order: exact Permission→operation/client-class mapping, technical non-Product ingress classification, final Problem/media consistency as needed, and the single machine-readable OpenAPI authority/tooling decision. Do not advance them before W3 is coherent.

Implementation remains blocked until D9.

---

## 6. Fresh-session success test

A fresh session must conclude:

- D0→D4/D4-R1 + D5-B1 are accepted/canonical;
- D5-B2 Operation Matrix + Whole-Matrix are ratified;
- `D5-B2-WIRE-CONTRACT.md` is canonical W1 authority;
- `D5-B2-W2-SCHEMA-GRAMMAR.md` is the **single consolidated canonical W2 authority**;
- former W2-C/D/E staging artifacts and Whole-W2 review candidate are absent from the active tree and preserved only in Git history;
- `AI-DIALOG.md` has no active Whole-W2 review cycle;
- W1 custom-method preconditions use typed request ETag rather than base-resource `If-Match` on a distinct URI;
- W2 includes the ratified historical publication, missing-schema, provider-requirement, Fulfillment identity/target, key-lifetime, media, idempotency and problem corrections;
- **W3 Collections / Pagination / Filter / Search / Cursor Grammar is the exact next action**;
- implementation remains blocked until D9.

If not, the active authority tree is inconsistent.
