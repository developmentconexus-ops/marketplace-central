# Marketplace Central — Architecture Rebaseline

> **Role:** sole current-program status/router after `AGENTS.md`  
> **Current phase:** **D5 — API — OPEN / ACTIVE; D5-B1 ACCEPTED / CANONICAL; D5-B2 OPEN / ACTIVE**  
> **D5-B2 current state:** **Wire W1 + W2-A/B/C/D/E ACCEPTED IN-STAGE; Whole-W2 Fable Round 1 + GPT adjudication COMPLETE; focused Fable Round 2 on precondition grammar = NEXT**  
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
16. `docs/engineering/rebaseline/D5-B2-WIRE-CONTRACT.md`
17. `docs/engineering/rebaseline/D5-B2-W2-SCHEMA-GRAMMAR.md` — W2-A/B
18. `docs/engineering/rebaseline/D5-B2-W2-C-READINESS-MARKET-ECONOMICS.md` — W2-C
19. `docs/engineering/rebaseline/D5-B2-W2-D-OPERATIONAL-SCHEMAS.md` — W2-D
20. `docs/engineering/rebaseline/D5-B2-W2-E-TRANSVERSAL-CONSISTENCY.md` — W2-E
21. `docs/engineering/rebaseline/EVIDENCE-REGISTER.md`
22. code/OpenAPI/schemas/tests/runtime only as current-state evidence when needed

This file alone owns **program status, allowed/blocked work and exact next action**. Detailed semantics remain in accepted artifacts above.

`D5-B2-W2-WHOLE-COHERENCE-REVIEW-CANDIDATE.md` and `AI-DIALOG.md` remain **NON-AUTHORITATIVE review evidence** and are excluded from the authority path. Fable/GPT findings do not modify W2-A/B/C/D/E until final convergence + operator ratification + canonical consolidation.

Older “next” wording in child artifacts/review files is a stage snapshot; this router is current status authority.

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
      Whole-W2 lead Global Coherence                     COMPLETE / RESTRUCTURE W2-LOCAL
      operator lead-direction ratification               COMPLETE
      independent Fable Round 1                          COMPLETE / REVISE W2-LOCAL
      GPT Round-1 adjudication                           COMPLETE
      focused Fable Round 2 — precondition grammar      NEXT
      GPT final adjudication                             AFTER ROUND 2
      operator final converged ratification              AFTER ADJUDICATION
D6 — Frontend                                             BLOCKED BY D5
D7 — Runtime / Jobs / Transactions                        BLOCKED
D8 — Golden Flows                                         BLOCKED
D9 — Adversarial Architecture Review                      BLOCKED
Implementation                                            BLOCKED UNTIL D9
```

## 3. Whole-W2 Round-1 adjudication

### Converged W2-local correction direction

The following materially converged between lead review, independent Fable challenge and GPT adjudication:

1. **historical publication basis** — `GetListingIntent` must expose append-only historical dispatch/effect basis sufficient to explain exact attempt-time resolved Offering values/provenance, requirement revision, media, authorization/disposition and typed PriceIntent/Availability correlations without creating a `PublicationAttempt` Product resource or cross-owner current authority;
2. **missing schema homes** — add bounded `MarketplaceListing`, `FulfillmentNode` and `EconomicPerformanceSummary`; `FulfillmentState` is not a second resource and maps to the single `FulfillmentExecution` wire home;
3. **Fulfillment operating targets** — W2 must name a closed consumer-proven field inventory; if no concrete target meaning can be named from current authority/evidence, narrow/defer rather than expose a generic target/SLA map;
4. **media capability** — authored media is ListingIntent-owned capability; use a ListingIntent-bound multipart custom method, `200` + media descriptor, same media identity on exact idempotent replay, binary content included in semantic idempotency equivalence; exact precondition transport remains Round-2 scope;
5. **PublicationValue** — add bounded `number_unit = ExactDecimalString + requirement-scoped unit_key` and explicit `not_applicable` only when the current requirement permits N/A; no generic UoM/property bag;
6. **draft-dependent provider requirements** — preserve Readiness requirement definition, Offering dispatchability and separate Price/Availability ownership while D4 may technically compose exact current owner inputs for provider conditional validation; result is revision-anchored evidence and may be unknown/unavailable; selected User-Product lane concrete validation surface remains D4/D8 proof;
7. **FulfillmentExecutionId** — accepted as the one current durable Fulfillment lifecycle identity required by checkpoint/history/artifact/Work/Materialization correlation; no parallel FulfillmentIntent/Workflow ID for the same meaning and no speculative split-routing policy;
8. **sale_line_key** — once minted within a Sale it never rebinds; reinterpretation may retire and mint a new key. Transient publication candidate/option keys remain context-bound because historical dispatch basis snapshots material meaning.

No current item above proves a parent-stage reopen.

### New Round-1 material findings

**F-IND-1 — ProductChannelCorrespondence stale-state carrier.** Resolve/Clear correspondence require current correspondence revision, but the subject is keyed and has no canonical ID. The finding is accepted; final carrier depends on the focused precondition decision below.

**F-GPT-1 — custom-method `If-Match` contradiction.** This is the only material contradiction surviving Round 1.

RFC HTTP conditional semantics apply `If-Match` to the request's target resource. A resource-bound custom method URI such as:

```text
POST /listing-intents/{id}:submit
POST /listing-intents/{id}:create-media
POST /work/{id}:hold
POST /fulfillment-executions/{id}:record-conference
```

is not the same request target URI as the base resource GET that emitted `/.../{id}`'s ETag. W1/W2 therefore cannot simply reuse the base-resource ETag as an RFC `If-Match` header on `:verb` and claim standard conditional-request semantics.

The current Round-2 candidate is:

- reserve HTTP `If-Match` for requests whose request URI actually identifies the conditionally mutated resource (ordinary PATCH/PUT/DELETE or another genuine same-URI standard method);
- custom owner methods carry the acted-on resource's same opaque ETag in the typed request (`etag` field; multipart `etag` part) when current state is required;
- another referenced resource carries its ETag adjacent to that typed reference when exact revision is material;
- the ETag value remains one server-issued validator; request-field transport is not a second version authority;
- missing required request-field ETag → `422 validation-error`; stale request-field/referenced ETag → `409 resource-revision-conflict`; `428/412` remain reserved for actual HTTP conditional-header failures;
- idempotency fingerprint includes all material ETag/precondition values and exact replay is resolved before a now-stale revision is rechecked;
- custom requests may remain free of **business payload** while carrying technical freshness proof.

This is a W1/W2-local wire question, not a D5-B1 or parent semantic reopen.

## 4. What is prohibited now

While focused Round 2 is open:

- do not canonicalize any Whole-W2 review correction into W1/W2 artifacts;
- do not reopen G1/G2/G4/G5/G6/G7 unless the precondition challenge directly falsifies them;
- do not begin collection/pagination/Permission/OpenAPI/tooling sub-batches;
- do not begin D6–D9 or implementation;
- do not treat Fable/GPT review evidence as authority;
- do not use a colon custom-method URI as an implicit alias of the base resource merely to make `If-Match` fit;
- do not convert consequential owner capabilities into fake PATCH/status CRUD merely to reuse conditional headers;
- do not create a synthetic Correspondence ID merely to obtain concurrency;
- do not weaken idempotency-before-stale-replay or no-blind-external-retry safety.

## 5. Exact next action

**Run one focused Fable Round 2 on F-GPT-1 + F-IND-1, using media as the concrete adversarial example.**

Fable must challenge:

1. whether RFC `If-Match` on `/resource/{id}:verb` can legitimately validate the ETag from `GET /resource/{id}` without inventing alias resource semantics;
2. request `etag` field/part versus custom header versus redesign of individual capabilities as honest standard methods;
3. every admitted current-state-protected C operation, classifying its precondition carrier as:
   - same-resource standard method → HTTP `If-Match`;
   - owner custom method → request ETag candidate;
   - create/capability depending on another resource → typed referenced ETag;
   - keyed meaning such as ProductChannelCorrespondence → honest keyed standard update or custom method + request ETag;
4. idempotency/lost-response ordering under the revised carrier model;
5. whether `428/412` versus `422/409` is the smallest honest problem grammar.

Fable may modify **only `AI-DIALOG.md`**, commit + push to the same branch and end `HANDOFF → GPT`.

No parent reopen unless focused proof actually requires one.

Implementation remains blocked until D9.

## 6. Fresh-session success test

A fresh session must conclude:

- D0→D4/D4-R1 + D5-B1 accepted/canonical;
- D5-B2 Operation Matrix + Whole-Matrix ratified;
- W1 + W2-A/B/C/D/E remain accepted in-stage authority;
- Whole-W2 lead review and Fable Round 1 are non-authoritative evidence;
- GPT Round-1 adjudication converged G1/G2/G4/G5/G6/G7, accepted F-IND-1 as a real gap and found F-GPT-1 as the only surviving material contradiction;
- **focused Fable Round 2 on precondition grammar is the exact next action**;
- no canonical W2 corrections land until Round-2 convergence + operator final ratification;
- implementation remains blocked until D9.

If not, the active authority tree is inconsistent.
