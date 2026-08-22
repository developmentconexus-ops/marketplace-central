# Marketplace Central — Roadmap

<!-- program-status-authority -->

> **Role:** sole mutable current-stage/status/allowed-work/next-action authority. Detailed semantics live in routed owning documents.

## Current checkpoint

| Field | Current value |
| --- | --- |
| Product | **Marketplace Operations Control Plane + Commercial Intelligence** |
| Current stage | **D6-R2 — Complete Frontend Realization Closure — OPEN / ACTIVE** |
| Accepted baseline | **D0–D8 ACCEPTED / CLOSED / INTEGRATED** |
| D6-R2 authority | [D6-R2 Complete Frontend Realization Closure](engineering/rebaseline/D6-R2-COMPLETE-FRONTEND-REALIZATION-CLOSURE.md) — **P0–P3 DERIVED; P8 NOT STARTED; no D6-R2 UX block is operator-LOCKED** |
| Execution method | [Frontend Product Experience Planning Method v2.1](development/frontend-product-experience-planning-method.md) — reusable methodology; **not** stage/status authority |
| Canonical Product OAD | `contracts/api/product/openapi.yaml` |
| Product surface | **99 Product operations · 30 ordinary Permissions · Principal kinds H / A / S only** |
| Stable origin | `https://conexus.fun` |
| Active runtime baseline | **NONE** |
| Exact next action | **Execute P4 bounded IA revalidation from the derived N01–N16 / UF01–UF16 / P3 coverage matrix. Preserve accepted D6 IA/routes by default; reopen only on a material falsifier. Do not draw P8 wireframes yet.** |
| Pre-D9 Implementation Readiness Contract | **BLOCKED UNTIL D6-R2 ACCEPTED / CLOSED** |
| D9 | **BLOCKED** |
| Implementation | **BLOCKED UNTIL D9** |

## Stage progression

| Stage / gate | Status |
| --- | --- |
| D0 — Product / System Definition | ACCEPTED / CLOSED |
| D1 — Domains / Boundaries | ACCEPTED / CLOSED |
| D2 — Identity / Tenant / Data Ownership | ACCEPTED / CLOSED |
| D3 — Communication / Events | ACCEPTED / CLOSED |
| D4 — External Integrations | ACCEPTED / CLOSED |
| D4-R1 — Publication Input / Listing Authoring | ACCEPTED / CANONICAL |
| D5 — API | ACCEPTED / CLOSED |
| D6 — Frontend | **ACCEPTED / CLOSED** |
| D7 — Runtime / Jobs / Transactions | **ACCEPTED / CLOSED** |
| D8 — Golden Flows | **ACCEPTED / CLOSED — OPERATOR-RATIFIED / INTEGRATED** |
| D6-R2 — Complete Frontend Realization Closure | **OPEN / ACTIVE — P0–P3 DERIVED** |
| Pre-D9 Implementation Readiness Contract | **BLOCKED UNTIL D6-R2** |
| D9 — Adversarial Architecture Review | **BLOCKED** |
| Implementation | **BLOCKED UNTIL D9** |

## Accepted D8 result carried forward

Owning authority remains [D8 Golden Flows](engineering/rebaseline/D8-GOLDEN-FLOWS.md) + [D8-R1 Proof Closure & Implementation-Readiness Coherence](engineering/rebaseline/D8-R1-PROOF-CLOSURE-COHERENCE.md). The executed probe evidence remains [D8 Controlled Live Probe Protocol](engineering/rebaseline/D8-LIVE-PROBE-PROTOCOL.md) + [D8 Live Probe Evidence](engineering/rebaseline/D8-LIVE-PROBE-EVIDENCE.md).

D8 remains accepted authority and is not reopened by D6-R2 preference.

```text
GF-01  Publication & Marketplace Convergence
GF-02  Sale → Business System → Fiscal → Fulfillment → Outcome
GF-03  Performance Evidence Honesty
SR-01  PITR / Timeline Continuity Recovery
```

Final live-probe disposition carried forward:

```text
P1  ML Price/Availability                 PASS_CONVERGED
P2  ML fiscal/invoice/label               OPERATOR_RATIFIED_REDEFER
                                           first real open ML Sale / beta-flagged implementation drive
P3  Sankhya 313→306                       PASS_CONVERGED
P4  native Party create/update            NOT_TRIGGERED
P5  alternate destination/contact         CAPABILITY_NOT_PROVEN for full override
                                           contact reference converged; full destination stays external-required/unsupported
P6  additional fiscal branch/component    NOT_TRIGGERED
```

P2 remains an explicit future proof obligation. P5 is a capability narrowing: the sanctioned contact reference is not a full alternate street/fiscal destination override. D6-R2 may not silently present either as already-proven capability.

## D6-R2 execution contract

D6-R2 uses Frontend Product Experience Planning Method v2.1 with Marketplace Central tailoring:

```text
accepted D6 IA/routes/React/TanStack/topology remain accepted
P0–P3 global foundation first
P4 revalidates accepted IA rather than reopening it by preference
P6/P7 activate only on real ambiguity
P8 renders structural wireframes block by block; only operator may set LOCKED
P9 binds exact owner + operationId + Permission + H/A/S + identity + state + effect semantics
P14 closes frontend realization readiness
```

Current D6-R2 P0–P3 result:

- four accepted D0 actor contexts retained; no invented personas;
- 16 outcome-oriented user needs and 16 complete end-to-end flows derived without starting from screens;
- accepted D6-B1 coverage remains 99/99 Product operations across 14 owner/families;
- no new Product operation, Permission or semantic owner is required;
- no D0–D8 reopen is justified;
- task frequency/density, real device/work-floor distribution, detailed terminology comprehension, Overview priority and any future bulk-interaction need remain explicit P12 evidence assumptions rather than hidden design decisions;
- P8 is **NOT STARTED** and no D6-R2 block is `LOCKED`.

## Approved pre-D9 sequence

```text
D8 integrated
→ D6-R2 — Complete Frontend Realization Closure
→ Pre-D9 Implementation Readiness Contract
→ D9 — Adversarial Architecture Review
→ Product implementation only after accepted D9
```

If D6-R2 or the readiness contract legitimately changes an accepted operation, Permission, semantic owner or another material invariant through the smallest owning authority, affected D8 flows/controls must be revalidated before D9 opens. D9 never reviews golden flows derived from superseded authority.

One coherent gate lands before the next. Product implementation remains blocked. For task-specific reading, return to [`index.md`](index.md).
