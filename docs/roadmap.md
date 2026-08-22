# Marketplace Central — Roadmap

<!-- program-status-authority -->

> **Role:** sole mutable current-stage/status/allowed-work/next-action authority. Detailed semantics live in routed owning documents.

## Current checkpoint

| Field | Current value |
| --- | --- |
| Product | **Marketplace Operations Control Plane + Commercial Intelligence** |
| Current stage | **D8 — Golden Flows — CLOSEOUT RATIFIED / INTEGRATION PENDING** |
| Accepted baseline | **D0–D8 ACCEPTED / CLOSED on candidate branch** |
| D8 authority | [D8 Golden Flows](engineering/rebaseline/D8-GOLDEN-FLOWS.md) + [D8-R1 Proof Closure & Implementation-Readiness Coherence](engineering/rebaseline/D8-R1-PROOF-CLOSURE-COHERENCE.md) — **ACCEPTED / CLOSED / OPERATOR-RATIFIED 2026-08-22** |
| Live-probe execution | [D8 Controlled Live Probe Protocol](engineering/rebaseline/D8-LIVE-PROBE-PROTOCOL.md) — **EXECUTED 2026-08-22**; outcomes in [D8 Live Probe Evidence](engineering/rebaseline/D8-LIVE-PROBE-EVIDENCE.md) |
| Canonical Product OAD | `contracts/api/product/openapi.yaml` |
| Product surface | **99 Product operations · 30 ordinary Permissions · Principal kinds H / A / S only** |
| Stable origin | `https://conexus.fun` |
| Active runtime baseline | **NONE** |
| Exact next action | **Evaluate and operator-ratify the D6-R2 frontend-realization method while PR #59 remains the sole D8 integration vehicle. Do not merge without explicit merge authorization; do not open/stack D6-R2 work until D8 is integrated into `main`.** |
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
| D8 — Golden Flows | **ACCEPTED / CLOSED — OPERATOR-RATIFIED / PENDING INTEGRATION INTO `main`** |
| D6-R2 — Complete Frontend Realization Closure | **BLOCKED UNTIL D8 INTEGRATION** |
| Pre-D9 Implementation Readiness Contract | **BLOCKED UNTIL D6-R2** |
| D9 — Adversarial Architecture Review | BLOCKED |
| Implementation | BLOCKED UNTIL D9 |

## Accepted D8 closure

D8 confirmed the smallest representative falsifiable set:

```text
GF-01  Publication & Marketplace Convergence
GF-02  Sale → Business System → Fiscal → Fulfillment → Outcome
GF-03  Performance Evidence Honesty
SR-01  PITR / Timeline Continuity Recovery
```

This remains **3 business golden flows + 1 systemic recovery falsifier**, not an exhaustive 99-operation test catalog and not Product implementation. Independent Fable review, GPT adjudication and controlled real probes closed without a D0–D7 semantic reopen.

Final probe disposition:

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

P5 is a capability narrowing, not a contradiction: the sanctioned Sankhya contact reference survives fiscal progression, but it is not a full street/fiscal destination override. Partner master corruption or duplicate-Party creation remains forbidden. P2 travels forward as an explicit proof obligation and may not silently disappear.

## Approved pre-D9 realization sequence

D8 established that **architecture closure is not by itself implementation-readiness closure**. After D8 is integrated, the accepted sequence is:

```text
D8 integrated
→ D6-R2 — Complete Frontend Realization Closure
→ Pre-D9 Implementation Readiness Contract
→ D9 — Adversarial Architecture Review
→ Product implementation only after accepted D9
```

D6-R2 must complete material route/screen/state/action/API/Permission/wireframe realization against current accepted Product + D7/D8 authority without reopening accepted frontend technologies/topology by preference. The Pre-D9 Implementation Readiness Contract then projects accepted architecture into per-slice outcomes, the complete D6 dependency/import law, negative controls, executable acceptance and explicit exit states so coding agents do not decide material behavior during implementation.

If D6-R2 or the readiness contract legitimately changes an accepted operation, Permission, semantic owner or another material invariant through the smallest owning authority, affected D8 flows/controls must be revalidated before D9 opens. D9 never reviews golden flows derived from superseded authority.

Until PR #59 is explicitly authorized and merged, keep D6-R2 blocked and do not stack its branch on this candidate. After merge, revalidate `main`, branches, PRs and CI first. Product implementation remains blocked until accepted D9.

One coherent gate lands before the next. For task-specific reading, return to [`index.md`](index.md).
