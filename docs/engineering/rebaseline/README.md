# Marketplace Central — Architecture Rebaseline

> **Role:** sole current-program status/router after `AGENTS.md`  
> **Current phase:** **D5 — API — OPEN / ACTIVE; D5-B1 ACCEPTED / CANONICAL; D5-B2 OPEN / ACTIVE**  
> **D5-B2 current state:** **Wire W1 + W2-A/B/C/D/E ACCEPTED IN-STAGE; Whole-W2 lead review = RESTRUCTURE NOW / W2-local; operator ratified G1–G7 review direction; independent Fable Whole-W2 review = NEXT**  
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

`docs/engineering/rebaseline/D5-B2-W2-WHOLE-COHERENCE-REVIEW-CANDIDATE.md` and `AI-DIALOG.md` are **NON-AUTHORITATIVE review input** and deliberately excluded from the authority path. The operator ratified G1–G7 only as the lead direction to be independently challenged; that ratification does not modify W2-A/B/C/D/E by implication.

Older “next” wording in child artifacts or the review candidate is a stage snapshot; this router is current status authority.

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
      operator lead-direction ratification               COMPLETE
      independent Fable Whole-W2 review                  NEXT
      GPT adjudication                                   AFTER FABLE
      operator final converged ratification               AFTER ADJUDICATION
D6 — Frontend                                             BLOCKED BY D5
D7 — Runtime / Jobs / Transactions                        BLOCKED
D8 — Golden Flows                                         BLOCKED
D9 — Adversarial Architecture Review                      BLOCKED
Implementation                                            BLOCKED UNTIL D9
```

## 3. Accepted W2 baseline still in force

While independent review is open, W2-A/B/C/D/E remain the accepted in-stage authority. Review findings are evidence only.

Load-bearing baseline:

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

## 4. Operator-ratified lead review direction — still non-authoritative

The independent reviewer must challenge, not inherit, these seven lead findings:

1. **G1 — historical publication basis:** add immutable ListingIntent dispatch/attempt basis sufficient to explain exact resolved values/provenance/requirement revision/media/PriceIntent/Availability/authorization/scope/result without a new generic Operation/PIM resource;
2. **G2 — missing admitted schema homes:** add bounded `MarketplaceListing`, `FulfillmentNode` and `EconomicPerformanceSummary` schemas; current W2-E complete-coverage claim is otherwise false;
3. **G3 — media precondition contradiction:** challenge the candidate correction from media collection POST + parent ETag misuse to a ListingIntent-bound multipart `:create-media` capability; binary content must participate in semantic idempotency equivalence;
4. **G4 — publication value completeness:** challenge bounded `number_unit` + explicit `not_applicable` support without a generic UoM/property-bag system;
5. **G5 — draft-dependent conditional requirements:** challenge a bounded D4 technical evaluation seam that may compose current ListingIntent + PriceIntent + Availability only to evaluate provider conditional requirements, while Readiness keeps requirement authority and Offering keeps dispatchability;
6. **G6 — Fulfillment identity hardening:** challenge whether `FulfillmentExecutionId` is genuinely the one current durable Fulfillment lifecycle identity or should be removed as speculative/synthetic before OpenAPI;
7. **G7 — scoped-key historical stability:** durable `sale_line_key` references cannot be recycled/rebound; transient publication candidate/option keys must not be the sole basis for historical explainability.

Lead conclusion only: **no parent-stage reopen currently proven**. G5 retains a targeted D1/D3/D4-R1 reopen trigger if real evidence shows the owner-preserving technical seam cannot express the provider conditional requirement contract.

## 5. Review-cycle rules / prohibited work

While the independent Whole-W2 review is open:

- do not mutate W2-A/B/C/D/E to incorporate G1–G7;
- do not start collection/pagination/Permission/OpenAPI/tooling sub-batches;
- do not begin D6–D9 or implementation;
- do not treat lead/Fable finding severity as requirement authority;
- Fable must reconstruct authority independently and attack G1–G7 plus search for additional material contradictions;
- Fable may modify **only `AI-DIALOG.md`** under the operator-authorized review scope; no other repository file may be changed by the reviewer;
- do not reopen parents merely because the selected provider protocol is complex; first test whether the accepted D4 mechanism/authority seam is sufficient;
- do not weaken historical explainability, source qualification, knowledge honesty, owner separation, idempotency, concurrency, client-class safety or YAGNI for API convenience.

## 6. Exact next action

**Run one coherent independent Fable D5-B2 Whole-W2 Global Coherence review.**

Review inputs:

- current authority path above;
- accepted W2-A/B/C/D/E;
- `D5-B2-W2-WHOLE-COHERENCE-REVIEW-CANDIDATE.md` only as bounded non-authoritative lead evidence;
- active `AI-DIALOG.md` review cycle.

Fable must:

1. revalidate branch HEAD;
2. reconstruct repository authority independently from `AGENTS.md` + this router;
3. follow the canonical Standard Fable review workflow;
4. attack every G1–G7 finding and seek additional material contradictions;
5. use current primary external standards/provider evidence only when material;
6. return material findings only in `AI-DIALOG.md`;
7. modify no other file;
8. commit + push the `AI-DIALOG.md` review to the same branch;
9. end `HANDOFF → GPT`.

After Fable finishes, GPT revalidates HEAD, adjudicates every material finding against authority/evidence, and requests a focused Round 2 only if a real material contradiction survives. No canonical W2 correction lands before final convergence + operator ratification.

Implementation remains blocked until D9.

## 7. Fresh-session success test

A fresh session must conclude:

- D0→D4/D4-R1 + D5-B1 accepted/canonical;
- D5-B2 Operation Matrix + Whole-Matrix ratified;
- W1 + W2-A/B/C/D/E accepted in-stage;
- Whole-W2 lead review is non-authoritative evidence with G1–G7 direction operator-ratified only for independent challenge;
- `AI-DIALOG.md` has an active Whole-W2 Fable review cycle;
- **independent Fable Whole-W2 review is the exact next action**;
- canonical W2 artifacts remain unchanged until review adjudication + final operator ratification;
- implementation remains blocked until D9.

If not, the active authority tree is inconsistent.
