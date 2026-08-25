# Marketplace Central — Roadmap

<!-- program-status-authority -->

> **Role:** sole mutable current-stage/status/allowed-work/next-action authority. Detailed semantics live in routed owning documents.

## Current checkpoint

| Field | Current value |
| --- | --- |
| Product | **Marketplace Operations Control Plane + Commercial Intelligence** |
| Current stage | **D6-R2 — Complete Frontend Realization Closure — OPEN / ACTIVE** |
| Accepted baseline | **D0–D8 ACCEPTED / CLOSED / INTEGRATED; D5-R2 + D8-R2 ACCEPTED** |
| Integrated checkpoints | **PR #61 → `b54d17bfe6d794645d198a9160f4a2a1c63647e8`; PR #62 methodology → `689dab34b0a756cbd7c790a6c3ae8627deebb535`; PR #65 proportional CI → `a2aeac19816c90ee30bf373cef0448d52a486c7e`; PR #68 PublicationRequirements wire → `ed3d164b0574b7950c2c7467d150c89576bba1ec`** |
| Method profile | **`developmentconexus-ops/conexus-methodology@9c7210d1504bef01c0d134a6c3ae8627deebb535` → `METHOD + FRONTEND-METHOD`** |
| Current acceptance increment | **PR #64 / B10 — Preparação — PR #68 bounded rebaseline complete; P8 GREEN CANDIDATE / OPERATOR WALKTHROUGH NEXT / NOT LOCKED** |
| B10 structural proof | **Draft quick CI #721 PASS at `9a10c09316901cd438861844ba2bb7e48210be0d`; B10 falsifiers 10/10; bootstrap 19054/20480** |
| Resolved upstream finding | **PublicationRequirements machine-readable wire lacked accepted provider/context/source distinctions; RESOLVED by integrated PR #68** |
| LOCK impact | **B00 / B01 / B00-R2 / B11 / B12 / B110 UNAFFECTED** |
| Canonical Product OAD | `contracts/api/product/openapi.yaml` |
| Product surface | **106 Product operations · 31 ordinary Permissions · Principal kinds H / A / S only** |
| Stable origin | `https://conexus.fun` |
| Active runtime baseline | **NONE** |
| Exact next action | **Operator operates the browser B10 P8 candidate and returns REVISE / UPSTREAM FINDING / explicit LOCK, plus A01 disposition `ACCEPT_FOR_LOCK_WITH_LATER_PROBE` or `BLOCK_LOCK`. Do not begin P9 before explicit operator LOCK.** |
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
| D5 — API | ACCEPTED / CLOSED; D5-R6 106/31 PROVED; D5-R7 W1 REPAIR ACCEPTED |
| D5-R2 — Operational Read Projection Repair | **ACCEPTED / CANONICAL** |
| D6 — Frontend | **ACCEPTED / CLOSED baseline; B00/B01/B00-R2/B11/B12/B110 LOCKED** |
| D7 — Runtime / Jobs / Transactions | **ACCEPTED / CLOSED baseline; D7-R ACCEPTED** |
| D8 — Golden Flows | **ACCEPTED / CLOSED — OPERATOR-RATIFIED / INTEGRATED; D8-R ACCEPTED** |
| D8-R2 — GF-02 Operational Read Revalidation | **ACCEPTED / PASS** |
| D6-R2 — Complete Frontend Realization Closure | **OPEN / ACTIVE — B10 P8 GREEN CANDIDATE / OPERATOR WALKTHROUGH NEXT / NOT LOCKED** |
| Pre-D9 readiness | **BLOCKED** |
| D9 — Adversarial Architecture Review | **BLOCKED** |
| Implementation | **BLOCKED UNTIL D9** |

## Current result

- PR #68 is integrated and `main` ruleset-required/full CI passed at `ed3d164b0574b7950c2c7467d150c89576bba1ec`.
- B10 remains the same search/list → exact-subject detail structure; P7 layout hypotheses remain not triggered.
- B10 now projects `PublicationRequirements` with independent `requirement_class` and `applicability`, all six source-evidence states, bounded `value_spec`, `not_applicable_allowed`, opaque source-candidate identity and separate source-media candidates.
- `missing source` remains source truth, not `publication impossible`; actual `FOLLOW_SOURCE` / `EXPLICIT_OVERRIDE` authoring remains downstream in ListingIntent/Offering.
- Draft quick proof is GREEN with B10 negative controls **10/10**; the PR remains Draft because the operator has not LOCKED B10.
- Existing LOCKED frontend blocks are unaffected. Product **106/31/H-A-S**, runtime NONE and Pre-D9/D9/implementation blocks remain unchanged.
- A01 remains **PENDING OPERATOR**: `ACCEPT_FOR_LOCK_WITH_LATER_PROBE` or `BLOCK_LOCK`.

```text
PR #68 INTEGRATED / MAIN GREEN
→ B10 P8 GREEN CANDIDATE
→ operator walkthrough + A01 disposition
→ explicit LOCK only if accepted
→ P9 only after LOCK
```

One coherent acceptance increment lands before the next. Return to [`index.md`](index.md) for task routing.
