# Marketplace Central — Roadmap

<!-- program-status-authority -->

> **Role:** sole mutable current-stage/status/allowed-work/next-action authority. Detailed semantics live in routed owning documents.

## Current checkpoint

| Field | Current value |
| --- | --- |
| Product | **Marketplace Operations Control Plane + Commercial Intelligence** |
| Current stage | **D6-R2 — Complete Frontend Realization Closure — OPEN / ACTIVE** |
| Accepted baseline | **D0–D8 ACCEPTED / CLOSED / INTEGRATED; D5-R2 + D8-R2 ACCEPTED** |
| Integrated checkpoint | **PR #61 squash-integrated at `b54d17bfe6d794645d198a9160f4a2a1c63647e8`; D6-R2 remained OPEN** |
| Current acceptance increment | **PR #62 — DevelopmentConexus methodology adoption — OPEN / CANDIDATE; repository-governance only** |
| Methodology target pin | **`developmentconexus-ops/conexus-methodology@9c7210d1504bef01c0d134a6c3ae8627deebb535`** |
| D6-R2 | [Closure](engineering/rebaseline/D6-R2-COMPLETE-FRONTEND-REALIZATION-CLOSURE.md) · [authority route](engineering/rebaseline/D6-R2-AUTHORITY-ROUTE.md) — **WHOLE-REPO FABLE R-1 + R-2 CLOSED; B10 PAUSED** |
| Canonical Product OAD | `contracts/api/product/openapi.yaml` |
| Product surface | **106 Product operations · 31 ordinary Permissions · Principal kinds H / A / S only** |
| Decision carrier | `CreateAuthorizationDecision` = `Idempotency-Key` + typed body `etag`/`outcome`; revision validation **422/409** |
| OAD source hygiene | **49 exact frozen historical-proof definitions; 0 new orphans; content digest guarded** |
| Stable origin | `https://conexus.fun` |
| Active runtime baseline | **NONE** |
| Exact next action | **Complete PR #62 as the bounded methodology-adoption candidate, prove aggregate conformance, then independently challenge the exact candidate under the pinned adversarial-review method. Do not resume B10, merge PR #62, begin Pre-D9/D9, or implement Product code first.** |
| Pre-D9 readiness | **BLOCKED UNTIL D6-R2 ACCEPTED / CLOSED** |
| D9 | **BLOCKED** |
| Implementation | **BLOCKED UNTIL D9** |

## Stage progression

| Stage / gate | Status |
| --- | --- |
| D0 — Product / System Definition | ACCEPTED / CLOSED; NOTIF-01 D0-R ACCEPTED |
| D1 — Domains / Boundaries | ACCEPTED / CLOSED; D1-R2 PASS |
| D2 — Identity / Tenant / Data Ownership | ACCEPTED / CLOSED; D2-R6 ACCEPTED |
| D3 — Communication / Events | ACCEPTED / CLOSED; D3-R3 ACCEPTED |
| D4 — External Integrations | ACCEPTED / CLOSED |
| D4-R1 — Publication Input / Listing Authoring | ACCEPTED / CANONICAL |
| D5 — API | ACCEPTED / CLOSED; **D5-R6 106/31 PROVED; D5-R7 W1 REPAIR ACCEPTED** |
| D5-R2 — Operational Read Projection Repair | **ACCEPTED / CANONICAL** |
| D6 — Frontend | **ACCEPTED / CLOSED baseline; B00/B01/B00-R2/B11/B12/B110 LOCKED; B110 REVALIDATED / UNAFFECTED** |
| D7 — Runtime / Jobs / Transactions | **ACCEPTED / CLOSED baseline; D7-R ACCEPTED** |
| D8 — Golden Flows | **ACCEPTED / CLOSED — OPERATOR-RATIFIED / INTEGRATED; D8-R ACCEPTED** |
| D8-R2 — GF-02 Operational Read Revalidation | **ACCEPTED / PASS** |
| D6-R2 — Complete Frontend Realization Closure | **OPEN / ACTIVE — PR #61 INTEGRATED; METHODOLOGY ADOPTION PR #62 OPEN; B10 PAUSED** |
| Pre-D9 readiness | **BLOCKED** |
| D9 — Adversarial Architecture Review | **BLOCKED** |
| Implementation | **BLOCKED UNTIL D9** |

## Current result

- PR #61 is integrated. Its Product/frontend/architecture authority remains intact; methodology adoption does not reinterpret or mutate it.
- PR #62 is a separate repository-governance acceptance increment. Its only purpose is to pin and route the accepted DevelopmentConexus methodology, remove duplicate active local methodology authority, preserve genuine MPC specialization, and strengthen aggregate repository conformance.
- AuthorizationRequest Global Maximum, Product **106/31/H-A-S**, D7-R, D8-R, OAD source-hygiene guards, and all operator LOCKs remain unchanged.
- B10 remains paused until PR #62 is accepted/integrated and B10 is then boundedly revalidated against the accepted pinned frontend method.

```text
PR #61 integrated → methodology-adoption PR #62 → B10 continuation inside D6-R2 → Pre-D9 readiness → D9 → Product implementation only after accepted D9
```

One coherent acceptance increment lands before the next. Return to [`index.md`](index.md) for task routing.
