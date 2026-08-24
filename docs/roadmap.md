# Marketplace Central — Roadmap

<!-- program-status-authority -->

> **Role:** sole mutable current-stage/status/allowed-work/next-action authority. Detailed semantics live in routed owning documents.

## Current checkpoint

| Field | Current value |
| --- | --- |
| Product | **Marketplace Operations Control Plane + Commercial Intelligence** |
| Current stage | **D6-R2 — Complete Frontend Realization Closure — OPEN / ACTIVE** |
| Accepted baseline | **D0–D8 ACCEPTED / CLOSED / INTEGRATED; D5-R2 + D8-R2 ACCEPTED** |
| Current acceptance increment | **PR #61 accumulated Product/frontend/architecture checkpoint — PROVED / READY FOR INTEGRATION; D6-R2 remains OPEN after integration** |
| D6-R2 | [Closure](engineering/rebaseline/D6-R2-COMPLETE-FRONTEND-REALIZATION-CLOSURE.md) · [authority route](engineering/rebaseline/D6-R2-AUTHORITY-ROUTE.md) — **WHOLE-REPO FABLE R-1 + R-2 CLOSED; B10 PAUSED** |
| Canonical Product OAD | `contracts/api/product/openapi.yaml` |
| Product surface | **106 Product operations · 31 ordinary Permissions · Principal kinds H / A / S only** |
| Decision carrier | `CreateAuthorizationDecision` = `Idempotency-Key` + typed body `etag`/`outcome`; revision validation **422/409** |
| OAD source hygiene | **49 exact frozen historical-proof definitions; 0 new orphans; content digest guarded** |
| Stable origin | `https://conexus.fun` |
| Active runtime baseline | **NONE** |
| Exact next action | **Await separate explicit operator authorization to squash-merge PR #61. After merge: revalidate `main`, then open a new independent methodology-adoption acceptance increment pinned to `developmentconexus-ops/conexus-methodology@9c7210d1504bef01c0d134a6c3ae8627deebb535`. Do not adopt methodology in PR #61, resume B10, or implement Product code first.** |
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
| D6-R2 — Complete Frontend Realization Closure | **OPEN / ACTIVE — PR #61 CHECKPOINT PROVED / READY FOR INTEGRATION; B10 PAUSED** |
| Pre-D9 readiness | **BLOCKED** |
| D9 — Adversarial Architecture Review | **BLOCKED** |
| Implementation | **BLOCKED UNTIL D9** |

## Current result

- PR #61 remains an accumulated legacy checkpoint, not the future PR-granularity model. Its Product/frontend/architecture content is proved and ready for integration; merging it does **not** close D6-R2.
- AuthorizationRequest Global Maximum remains accepted; D5-R7 changes only the custom-method revision carrier to canonical W1 typed-request semantics.
- Fable R-2 remains closed with **16 pathItems + 33 schemas = 49** exact frozen historical-proof definitions and **0 new orphans**; historical **95/29 + 99/30** and current **106/31** proofs remain binding.
- The accepted DevelopmentConexus methodology upgrade is a separate repository-governance acceptance increment. It must start from revalidated post-#61 `main`; it does not reinterpret or mutate this Product checkpoint.
- B10 stays paused until that separate methodology adoption is integrated and B10 is boundedly revalidated against it.

```text
PR #61 integration → methodology-adoption increment → B10 continuation inside D6-R2 → Pre-D9 readiness → D9 → Product implementation only after accepted D9
```

One coherent acceptance increment lands before the next. Return to [`index.md`](index.md) for task routing.
