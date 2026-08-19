# Marketplace Central — Architecture Rebaseline

> **Role:** sole current-program status/router after `AGENTS.md`  
> **Current phase:** **D5 — API — OPEN / ACTIVE; D5-B1 ACCEPTED / CANONICAL; D5-B2 OPEN / ACTIVE**  
> **D5-B2 current state:** **Wire W1 + W2-A/B/C/D/E ACCEPTED IN-STAGE; Whole-W2 lead review = RESTRUCTURE NOW / W2-local; operator ratification of lead correction direction = NEXT**  
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

This file alone owns **program status, allowed/blocked work and exact next action**. Detailed semantics remain in the accepted artifacts above.

`docs/engineering/rebaseline/D5-B2-W2-WHOLE-COHERENCE-REVIEW-CANDIDATE.md` is **NON-AUTHORITATIVE review input** and is deliberately excluded from the authority path. `AI-DIALOG.md`, that candidate, Git history, legacy routes/OpenAPI and current code are never target/status authority by inheritance.

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
      Whole-W2 lead Global Coherence                     RESTRUCTURE NOW / COMPLETE
      operator ratification of lead direction            NEXT
      Fable coherent W2 review                           BLOCKED UNTIL RATIFICATION
D6 — Frontend                                             BLOCKED BY D5
D7 — Runtime / Jobs / Transactions                        BLOCKED
D8 — Golden Flows                                         BLOCKED
D9 — Adversarial Architecture Review                      BLOCKED
Implementation                                            BLOCKED UNTIL D9
```

## 3. Accepted W2 baseline still in force

Until review findings are independently challenged, adjudicated and operator-ratified into canonical artifacts, W2-A/B/C/D/E remain the accepted in-stage authority.

Load-bearing baseline remains:

- Organization-scoped semantic Product API; opaque MPC IDs and explicit external source qualification;
- no Product/PIM, generic Integration/Action/Operation/Workflow/Rules/Task/Finance/AI authority;
- exact decimal Money; explicit knowledge/freshness/coverage; no universal Fact/Evidence/Result/ExternalRef/property bags;
- ListingIntent, PriceIntent and Availability remain distinct through initial publication;
- Readiness owns requirements/source candidates, Offering owns draft dispatchability, D4 owns provider protocol;
- Market coverage != evidence sufficiency; Economics never fabricates profitability from missing evidence;
- Governance authorizes but never executes; Sales remains external; BusinessOrder/Invoicing are owner reactions;
- physical Fulfillment checkpoints remain owner facts; Post-Sale consequence tracks remain distinct; Work never owns source truth;
- direct concurrency uses strong `ETag` + `If-Match`; referenced-resource revision uses a distinct referenced precondition; idempotency is duplicate-intake safety, not retry permission;
- Problem Details remains separate from valid business outcomes;
- implementation remains blocked until D9.

## 4. Whole-W2 lead findings — non-authoritative until ratified/converged

The lead candidate currently proposes seven bounded findings:

1. **G1 — historical publication basis:** add immutable ListingIntent dispatch/attempt basis sufficient to explain exact values/provenance/requirement revision/media/owner correlations and provider result without a new generic Operation/PIM resource;
2. **G2 — missing admitted schema homes:** add bounded `MarketplaceListing`, `FulfillmentNode` and `EconomicPerformanceSummary` schemas; W2-E's current complete-coverage claim is otherwise false;
3. **G3 — media precondition contradiction:** replace media collection POST + parent ETag misuse with a ListingIntent-bound media capability (candidate direction `POST {listing-intent}:create-media`) so `If-Match` conditions the actual selected resource; binary content participates in semantic idempotency equivalence;
4. **G4 — publication value completeness:** current provider evidence requires bounded `number_unit` support and explicit N/A meaning without a generic UoM/property-bag system;
5. **G5 — draft-dependent conditional requirements:** current provider conditional-required validation can depend on concrete listing/joint owner inputs; preserve Readiness requirement authority and Offering dispatchability while adding a bounded D4 draft-dependent validation seam rather than provider rule DSL/client payload authority;
6. **G6 — Fulfillment identity hardening:** if `FulfillmentExecutionId` survives, it is the one concrete D2 durable Fulfillment lifecycle/intent identity, not a second ID justified by speculative split fulfillment;
7. **G7 — scoped-key historical stability:** durable references such as `sale_line_key` cannot later be recycled/rebound; transient publication candidate/option keys must resolve into historical material values/provenance before consequential history depends on them.

Current lead conclusion: **no parent-stage reopen is proven**. G5 carries a targeted D1/D3/D4-R1 reopen trigger only if real implementation proves conditional requirement evaluation cannot fit accepted owner-preserving technical composition.

## 5. What is prohibited now

While operator ratification is pending:

- do not mutate W2-A/B/C/D/E to incorporate the lead candidate;
- do not invoke Fable yet;
- do not begin collection/pagination/Permission/OpenAPI/tooling sub-batches;
- do not begin D6–D9 or implementation;
- do not treat reviewer/lead severity as requirement authority;
- do not re-open parents merely because current Mercado Libre protocol is complex; first test whether the accepted D4 mechanism/authority seam is sufficient;
- do not weaken historical explainability, source qualification, knowledge honesty, owner separation, idempotency, concurrency or client-class safety for API convenience.

## 6. Exact next action

**Operator reviews and ratifies/revises the Whole-W2 lead correction direction G1–G7.**

If ratified:

1. retain W2-A/B/C/D/E as accepted authority while independent review is open;
2. use `D5-B2-W2-WHOLE-COHERENCE-REVIEW-CANDIDATE.md` as the bounded non-authoritative review target;
3. open `AI-DIALOG.md` for one coherent independent Fable Whole-W2 review using the canonical Standard Fable workflow;
4. Fable reconstructs authority independently, attacks G1–G7 and searches for additional material contradictions;
5. GPT adjudicates reviewer findings technically;
6. Round 2 only if a real material contradiction survives;
7. after convergence, operator ratifies the final package;
8. only then consolidate corrections into canonical W2 artifacts, delete the candidate/reset AI-DIALOG and advance to the next Wire Contract sub-batch.

Implementation remains blocked until D9.

## 7. Fresh-session success test

A fresh session must conclude:

- D0→D4/D4-R1 + D5-B1 accepted/canonical;
- D5-B2 Operation Matrix + Whole-Matrix ratified;
- W1 + W2-A/B/C/D/E accepted in-stage;
- Whole-W2 lead review found W2-local corrections and is non-authoritative evidence;
- operator ratification of G1–G7 is exact next action;
- Fable is blocked until that ratification;
- implementation remains blocked until D9.

If not, the active authority tree is inconsistent.
