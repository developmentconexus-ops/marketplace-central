# Marketplace Central — Architecture Rebaseline

> **Role:** sole current-program status/router after `AGENTS.md`  
> **Current phase:** **D5 — API — OPEN / ACTIVE; D5-B1 ACCEPTED / CANONICAL; D5-B2 OPEN / ACTIVE**  
> **D5-B2 current state:** **Wire W1 + W2-A/B/C/D/E ACCEPTED IN-STAGE; Whole-W2 Fable Rounds 1–2 + GPT final adjudication CONVERGED; operator final ratification of converged W2-local corrections = NEXT**  
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

This file alone owns **program status, allowed/blocked work and exact next action**. Detailed semantics remain in the accepted artifacts above.

`D5-B2-W2-WHOLE-COHERENCE-REVIEW-CANDIDATE.md` and `AI-DIALOG.md` are **NON-AUTHORITATIVE review evidence** and are deliberately excluded from the authority path. Fable/GPT findings do not modify W1/W2 until operator final ratification + canonical consolidation.

Older “next” wording in child/review artifacts is a stage snapshot; this router is current status authority.

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
      focused Fable Round 2                              COMPLETE
      GPT final adjudication                             COMPLETE / CONVERGED
      operator final converged ratification              NEXT
D6 — Frontend                                             BLOCKED BY D5
D7 — Runtime / Jobs / Transactions                        BLOCKED
D8 — Golden Flows                                         BLOCKED
D9 — Adversarial Architecture Review                      BLOCKED
Implementation                                            BLOCKED UNTIL D9
```

## 3. Converged Whole-W2 correction package — NON-AUTHORITATIVE UNTIL FINAL RATIFICATION

The review cycle converged on bounded W1/W2-local corrections only; **no parent-stage reopen is required**.

1. **Historical ListingIntent dispatch basis** — expose append-only attempt/effect history through existing `GetListingIntent`; preserve dispatch-time resolved values/knowledge/provenance, requirement revision, media, authorization/disposition and typed PriceIntent/Availability correlations without `PublicationAttempt` CRUD or cross-owner current authority.
2. **Missing schema homes** — add bounded `MarketplaceListing`, `FulfillmentNode`, `EconomicPerformanceSummary`; `FulfillmentState` is not a second resource and maps to `FulfillmentExecution`.
3. **Fulfillment operating target baseline** — one current evidence-backed field only: internal `dispatch_handoff_lead_time_before_provider_deadline`, with inherit/override/effective provenance; provider deadline remains external evidence; no generic SLA/target map.
4. **Media** — ListingIntent-bound multipart `:create-media`, `200` + ListingIntent-scoped descriptor, same media identity on exact idempotent replay, binary content in semantic fingerprint; no standalone media URI/Asset authority.
5. **PublicationValue** — add bounded `number_unit = ExactDecimalString + requirement-scoped unit_key` and explicit `not_applicable` only where Readiness says N/A is permitted; no UoM engine/property bag.
6. **Draft-dependent conditional requirements** — Readiness keeps requirement definition; Offering keeps dispatchability; D4 may technically compose exact current ListingIntent + PriceIntent + Availability inputs for provider conditional validation. Evaluation evidence is revision-anchored and may be unknown/unavailable. Selected User-Product lane concrete validation remains D4/D8 proof.
7. **Fulfillment identity** — `FulfillmentExecutionId` is the one durable Fulfillment lifecycle identity for checkpoint/history/artifact/Work/Materialization correlation; no parallel FulfillmentIntent/Workflow ID for the same meaning and no speculative split-routing policy.
8. **`sale_line_key` lifetime** — once minted within a Sale, it never rebinds; reinterpretation may retire + mint new. Transient publication candidate/option keys rely on historical dispatch snapshots instead of eternal resolvability.
9. **Revision-precondition grammar** — HTTP `If-Match` is used only on standard methods whose request URI is the actual conditionally protected resource. Owner custom methods carry the same opaque validator as typed `etag` request data; referenced resources carry typed adjacent `etag`. No custom conditional header and no implicit alias semantics for `:verb` URIs.
10. **ProductChannelCorrespondence** — keep Resolve/Clear as Readiness owner capabilities over the keyed Product+Marketplace subject; `GetProductChannelReadiness` exposes a correspondence-scoped `etag` distinct from `requirements_revision`; Resolve/Clear carry that typed `etag`. No synthetic Correspondence/Readiness ID and no forced PUT/DELETE resource fiction.
11. **Problem/idempotency hardening** — true HTTP conditionals: missing/false → `428/412`; typed revision proof missing/invalid/stale → `422 validation-error` / `409 resource-revision-conflict`; exact prior idempotent intake resolves before stale revision re-evaluation; all material revision proofs participate in semantic fingerprint.
12. **D4/D8 proof notes only** — selected provider lane must prove authoritative N/A reread semantics and concrete User-Product conditional-requirement validation before D8 claims convergence. These notes do not reopen D4 semantics now.

The final GPT adjudication deliberately revised one Fable Round-2 proposal: correspondence does **not** use PUT/DELETE + `If-Match`, because that model cannot honestly protect unresolved/conflicting correspondence revision without either an absent-target gap or a DELETE-that-does-not-delete resource fiction. The already-selected typed `etag` custom-method grammar is smaller and complete.

## 4. What is prohibited now

Until operator final ratification:

- do not modify W1/W2 canonical artifacts to incorporate the converged review package;
- do not delete the review candidate or reset `AI-DIALOG.md` yet;
- do not start collection/pagination/Permission/OpenAPI/tooling sub-batches;
- do not begin D6–D9 or implementation;
- do not treat Fable/GPT review text as authority before operator ratification;
- do not reopen D0–D4/D4-R1/D5-B1 merely because the selected provider protocol is complex;
- do not add generic Product/PIM/Asset/UoM/Rules/SLA/Workflow/Operation/Finance/Task abstractions;
- do not weaken source qualification, knowledge honesty, owner separation, idempotency, concurrency/preconditions, physical-fact trust or Work source-truth fences.

## 5. Exact next action

**Operator reviews and ratifies/revises the complete converged Whole-W2 correction package in §3.**

If ratified:

1. revalidate remote HEAD immediately before writes;
2. consolidate the package into canonical W1/W2 artifacts;
3. preserve the ratified Operation Admission Matrix semantic inventory/Permissions/client classes/safety tuples unless an exact canonicalization contradiction proves otherwise;
4. verify the exact post-write diff and resulting remote HEAD;
5. delete `D5-B2-W2-WHOLE-COHERENCE-REVIEW-CANDIDATE.md`;
6. reset `AI-DIALOG.md` to protocol-only state; Git history remains review archive;
7. update this router to the next Wire Contract sub-batch;
8. implementation remains blocked until D9.

Canonical consolidation scope includes at least:

- `D5-B2-WIRE-CONTRACT.md` — W1 precondition/custom-method law;
- `D5-B2-W2-SCHEMA-GRAMMAR.md` — historical ListingIntent basis, MarketplaceListing, PublicationValue, submit/media-related schema law;
- `D5-B2-W2-C-READINESS-MARKET-ECONOMICS.md` — EconomicPerformanceSummary, draft-dependent requirement evidence, correspondence-scoped validator/carrier;
- `D5-B2-W2-D-OPERATIONAL-SCHEMAS.md` — FulfillmentNode, FulfillmentExecution naming/identity, `sale_line_key`, typed-etag operational custom methods;
- `D5-B2-W2-E-TRANSVERSAL-CONSISTENCY.md` — operating-target field, media transport, two-class revision carrier grammar, response/status/idempotency/Problem Details changes.

No Round 3 is required by current evidence.

## 6. Fresh-session success test

A fresh session must conclude:

- D0→D4/D4-R1 + D5-B1 accepted/canonical;
- D5-B2 Operation Matrix + Whole-Matrix ratified;
- W1 + W2-A/B/C/D/E remain accepted in-stage authority and are not yet changed by review findings;
- Whole-W2 lead/Fable/GPT review is non-authoritative evidence but has **converged after focused Round 2 + GPT final adjudication**;
- no parent-stage reopen is required;
- **operator final ratification of the converged Whole-W2 correction package is the exact next action**;
- only after ratification may GPT canonicalize W1/W2, delete/reset review artifacts and advance the router;
- implementation remains blocked until D9.

If not, the active authority tree is inconsistent.
