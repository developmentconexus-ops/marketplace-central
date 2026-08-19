# Marketplace Central — Architecture Rebaseline

> **Role:** sole current-program status/router after `AGENTS.md`  
> **Current phase:** **D5 — API — OPEN / ACTIVE; D5-B1 ACCEPTED / CANONICAL; D5-B2 OPEN / ACTIVE**  
> **D5-B2 current state:** **Operation Matrix + Whole-Matrix RATIFIED; Wire W1 + W2-A/B/C/D/E ACCEPTED IN-STAGE; Whole-W2 Global Coherence lead review = NEXT**  
> **Implementation:** **BLOCKED until D9 is accepted**  
> **Last updated:** 2026-08-18

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
16. `docs/engineering/rebaseline/D5-B2-WIRE-CONTRACT.md`
17. `docs/engineering/rebaseline/D5-B2-W2-SCHEMA-GRAMMAR.md` — W2-A/B
18. `docs/engineering/rebaseline/D5-B2-W2-C-READINESS-MARKET-ECONOMICS.md` — W2-C
19. `docs/engineering/rebaseline/D5-B2-W2-D-OPERATIONAL-SCHEMAS.md` — W2-D
20. `docs/engineering/rebaseline/D5-B2-W2-E-TRANSVERSAL-CONSISTENCY.md` — W2-E
21. `docs/engineering/rebaseline/EVIDENCE-REGISTER.md`
22. code/OpenAPI/schemas/tests/runtime only as current-state evidence when needed

This file alone owns **program status, allowed/blocked work and exact next action**. Detailed semantics remain in the accepted artifacts above. `AI-DIALOG.md`, Git history, future review candidates, legacy routes/OpenAPI and current code are never target/status authority by inheritance.

Older “next” wording in child artifacts is a stage snapshot; this router is current status authority.

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
      W1 Resource / Path / HTTP Grammar                  ACCEPTED IN-STAGE
      W2 Request / Response Schema Grammar               OPEN / ACTIVE
        W2-A Core Schema Grammar                         ACCEPTED IN-STAGE
        W2-B ListingIntent / PriceIntent / Availability  ACCEPTED IN-STAGE
        W2-C Readiness / Market / Economics              ACCEPTED IN-STAGE
        W2-D Governance / Sales / Materialization /
             Fulfillment / Post-Sale / Work              ACCEPTED IN-STAGE
        W2-E Transversal / final consistency             ACCEPTED IN-STAGE
      Whole-W2 Global Coherence — lead                   NEXT
      W2 independent Fable review                        AFTER LEAD CONVERGENCE
D6 — Frontend                                             BLOCKED BY D5
D7 — Runtime / Jobs / Transactions                        BLOCKED
D8 — Golden Flows                                         BLOCKED
D9 — Adversarial Architecture Review                      BLOCKED
Implementation                                            BLOCKED UNTIL D9
```

## 3. Load-bearing Wire decisions

### W1

- Product API is Organization-scoped under `/organizations/{organization_id}/...`; `/access-context` is the bounded platform-scoped self-only Q.
- no `/v1` compatibility axis without a real entitled consumer;
- MPC IDs are opaque; external identities remain explicitly Marketplace Installation/SourceInstance-qualified;
- URI nesting means identity containment/namespace qualification, never workflow order;
- standard HTTP methods only when honest; owner-specific non-CRUD capability uses `POST {resource-uri}:verb`;
- no `PATCH status` escape hatch;
- direct stale-state protection uses strong opaque `ETag` + `If-Match`; required missing = 428, stale = 412.

### W2-A/B

- exact decimals serialize as strings; `Money = ExactDecimalString + currency`;
- typed source-qualified refs; no universal ExternalRef/entity graph;
- request objects are closed and do not accept server-owned effective authority/history;
- `null` never transports unknown/unavailable/partial;
- knowledge uses smallest owner-specific exclusive unions, not universal `Fact<T>`;
- no generic Result/Operation/Evidence/property-bag/rules DSL;
- typed owner-specific PATCH; omitted = unchanged;
- one sparse/declarative ListingIntent covers create/edit;
- dynamic publication uses Readiness requirement keys/revisions/candidates and only `FOLLOW_SOURCE | EXPLICIT_OVERRIDE`;
- PriceIntent always owns price, including initial publication; supersession is explicit;
- SellableAvailability separately exposes control, desired value, provider observation and convergence;
- client does not orchestrate price/availability correlations by PATCHing ListingIntent.

### W2-C

- source Product search remains source-qualified evidence, not Product mirror/master;
- ProductChannelReadiness and CompetitivePosition are contextual keyed Qs with no synthetic canonical ID;
- PublicationRequirements returns effective Readiness meaning, not a provider rules DSL;
- market coverage and evidence sufficiency remain distinct;
- ExpectedEconomics is components-first and does not fabricate profitability from incomplete evidence;
- scenario evaluation accepts hypothetical variables, not caller-authored authoritative evidence;
- SaleEconomics preserves L0/L1/L2 and R1/R2; realized economics is occurrence-based;
- Economic Attribution remains Economics-local exact/partial/ambiguous/unresolved meaning without a universal entity/reconciliation graph.

### W2-D

- AuthorizationDecision is immutable; AuthorizationDelegation has stable opaque identity only because update/revoke now require a stable wire target;
- Sale/Shipment remain source-qualified external identities; sale-line selector is Sale-scoped only;
- BusinessOrderIntent/InvoicingIntent remain owner-triggered tracking resources; native Sankhya results stay source-qualified and provider-local;
- PartyResolution/DestinationRealization are BusinessOrderIntent-contained meanings, not Customer/Address mastery;
- FulfillmentExecution is the Fulfillment-owned durable identity for physical checkpoint addressing, not a workflow/WMS domain;
- physical checkpoint occurrence is distinct from readiness conclusion; generic automation cannot self-declare trusted physical evidence;
- PostSaleResolution has scoped concurrent consequence tracks and no direct close/provider-action vocabulary;
- Work owns responsibility/lifecycle, not source truth; current Work closure-path audit passes and generic Work resolution remains deferred.

### W2-E

- shared inherit/override mechanics do not create cross-owner policy authority; effective owner configuration preserves provenance and provider deadlines remain external evidence;
- AccessContext/AccessRole/Membership and MarketplaceInstallation/SellingEntity schemas are covered without fake IDs or credential leakage;
- ListingIntent authored-media intake uses bounded multipart Product wire while storage/CDN/hash mechanics remain D7;
- owner outcome, external effect and convergence remain separate axes; no generic `failed`/OperationState;
- successful create = 201 + Location/resource; update/`:verb` = 200 + current resource/result; `202` is not a convergence flag;
- direct target stale state uses `If-Match`; exact referenced-resource revision uses a typed referenced-resource precondition and stale referenced state is 409, not misapplied 412;
- idempotency dedup is resolved before stale direct-resource precondition on exact replay; semantic fingerprint includes material preconditions and excludes raw JSON formatting/credentials;
- Problem Details uses a small stable RFC 9457 type catalog; valid business rejection/pending/unknown/unavailable/partial remains semantic meaning;
- all admitted Product owner families have a W2 schema home.

## 4. Prohibited now

While Whole-W2 lead review is open:

- do not begin collection/pagination/tooling/OpenAPI crystallization beyond what W1/W2 already accepted;
- do not begin D6–D9 target design or implementation;
- do not weaken accepted D0→D5-B1/B2/W1/W2-A/B/C/D/E for review convenience;
- do not derive target schema from legacy OpenAPI/provider DTO/database/frontend shape;
- do not introduce Product mirror/PIM, generic Result/Fact/Evidence/Operation/Subject/Scope/Policy/Workflow/property-bag abstractions;
- do not collapse knowledge, coverage, owner outcome, external effect and convergence into one status;
- do not expose provider/business-system native choreography as Product ontology;
- do not run Fable before the lead Whole-W2 review converges enough to prepare one bounded review candidate.

## 5. Exact next action

**Run the lead Whole-W2 Global Coherence Review across W2-A/B/C/D/E using the Method.**

Attack at least:

1. duplicate/missing authority or identity;
2. one meaning with two wire homes/IDs;
3. accidental universal wrappers emerging from repeated mechanics;
4. every admitted Product operation/family having a legitimate request/response home;
5. hidden client orchestration between owners;
6. provider/source ontology leakage;
7. knowledge/freshness/coverage/economic exactness honesty;
8. owner outcome vs external effect vs convergence separation;
9. direct `If-Match` vs referenced-resource precondition correctness;
10. idempotency replay/concurrency interaction;
11. API Problem Details vs valid business outcomes;
12. physical-fact client-class enforcement;
13. Work closure paths;
14. media Product semantics vs D7 mechanics;
15. policy default/inheritance/override without generic policy authority;
16. YAGNI/future retrofit and Structural Inversion against the legacy API.

Outcomes remain `RESTRUCTURE NOW`, `CURRENT STRUCTURE CONFIRMED` or `STOP / SPLIT PREREQUISITE`.

If a material contradiction is found, record the smallest W2/B2-local correction or targeted parent reopen and obtain operator ratification before canonical correction. Do not silently rewrite accepted W2 sections.

If the lead package converges, prepare one bounded **NON-AUTHORITATIVE W2 review candidate** and invoke Fable through the canonical Standard Fable review workflow. Fable output is evidence only until lead adjudication + operator ratification.

Implementation remains blocked until D9.

## 6. Fresh-session success test

A fresh session must conclude:

- D0→D4/D4-R1 + D5-B1 accepted/canonical;
- D5-B2 Operation Matrix + Whole-Matrix ratified;
- W1 + W2-A/B/C/D/E accepted in-stage;
- Whole-W2 lead coherence review is exact next action;
- independent Fable W2 review waits for a converged lead candidate;
- implementation remains blocked until D9.

If not, the active authority tree is inconsistent.
