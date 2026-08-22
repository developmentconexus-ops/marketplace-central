# Marketplace Central — Roadmap

<!-- program-status-authority -->

> **Role:** sole mutable current-stage/status/allowed-work/next-action authority. Detailed semantics live in routed owning documents.

## Current checkpoint

| Field | Current value |
| --- | --- |
| Product | **Marketplace Operations Control Plane + Commercial Intelligence** |
| Current stage | **D6 — Frontend — ACCEPTED / CLOSED** |
| Accepted baseline | **D0–D6 ACCEPTED / CLOSED** |
| D5-R1 auth correction | **OPERATOR-APPROVED / EXECUTABLE PROOF PASS** — [Human Browser Authentication Correction](engineering/rebaseline/D5-R1-HUMAN-BROWSER-AUTHENTICATION.md) |
| D6-R1 | **Marketplace Performance Intelligence — OPERATOR-APPROVED / FABLE REVIEWED / BOUNDED FIXES ADJUDICATED** |
| D6-B1 | **OPERATOR-RATIFIED — corrected Portuguese interaction map + wireframes** |
| D6-B2 | **OPERATOR-RATIFIED — frontend realization topology + dependency profile** |
| Final D6 challenge | **ACCEPT WITH BOUNDED FIXES — GPT ADJUDICATED / BOUNDED FIXES PROVED / NO ROUND 2 REQUIRED** |
| Cross-repo review | **Marketplace Central ↔ MetalDocs — OPERATOR-APPROVED / ADJUDICATED** |
| Canonical Product OAD | `contracts/api/product/openapi.yaml` |
| Product surface | **99 Product operations · 30 ordinary Permissions · Principal kinds H / A / S only** |
| Stable origin | `https://conexus.fun` |
| Active runtime baseline | **NONE** |
| Exact next action | **After PR #54 lands, revalidate merged `main`, then open D7 — Runtime / Jobs / Transactions on a fresh stage branch. D7 is not started by this closeout.** |
| Implementation | **BLOCKED UNTIL D9** |

## Stage progression

| Stage | Status |
| --- | --- |
| D0 — Product / System Definition | ACCEPTED / CLOSED |
| D1 — Domains / Boundaries | ACCEPTED / CLOSED |
| D2 — Identity / Tenant / Data Ownership | ACCEPTED / CLOSED |
| D3 — Communication / Events | ACCEPTED / CLOSED |
| D4 — External Integrations | ACCEPTED / CLOSED |
| D4-R1 — Publication Input / Listing Authoring | ACCEPTED / CANONICAL |
| D5 — API | ACCEPTED / CLOSED — bounded D5-R1 browser-auth carrier correction integrated and proved through D6 |
| D6 — Frontend | **ACCEPTED / CLOSED — operator-authorized after final independent challenge and adjudication** |
| D7 — Runtime / Jobs / Transactions | **NEXT / NOT STARTED** |
| D8 — Golden Flows | BLOCKED |
| D9 — Adversarial Architecture Review | BLOCKED |
| Implementation | BLOCKED UNTIL D9 |

## Accepted D6 authority

The accepted Product surface is now:

```text
99 Product operations
30 ordinary Permissions
Principal kinds H / A / S only
4 Marketplace Performance Qs under performance.read
stable origin https://conexus.fun
```

Human/machine authentication remains:

```text
H browser  -> server-side OIDC login -> Secure HttpOnly application session + CSRF on unsafe requests
A / S      -> Client Credentials -> audience-bound bearer
```

D6 frontend realization remains:

```text
React + TypeScript strict
TanStack Query              server-state authority
TanStack Router             typed/validated URL/navigation state
openapi-typescript          generated Product wire shapes
openapi-fetch               thin OAD-bound Product HTTP client
```

The detailed interaction, topology, URL, state, Product-client, retry/idempotency/concurrency and YAGNI laws remain owned by [D6 Frontend](engineering/rebaseline/D6-FRONTEND.md). The [interaction map](engineering/rebaseline/D6-B1-INTERACTION-MAP.md) and Portuguese [wireframes](../qualification/d6-wireframes/index.html) remain the ratified interaction proof.

## Final D6 proof

The final independent Fable challenge returned **ACCEPT WITH BOUNDED FIXES** and found no material reason to reopen D5-R1, D6-R1, D6-B1, D6-B2, the 99/30 Product surface, the authentication split, or the selected frontend topology/dependency profile.

GPT adjudicated every material finding. The accepted bounded corrections repaired the default gate entry point, D6-B1 status/count precision, and fixed-metric Traffic/Sales comparison-state precision without forcing false Retail Media metric uniformity. Repository-only reporting observations were not promoted into architecture work.

Fresh executable proof after those fixes establishes:

```text
accepted D5 baseline             95/95 operations · 29/29 Permissions · 12/12 controls
current D6 Product               99/99 operations · 30/30 Permissions · 28/28 List/Search
Performance controls             7/7 PASS
current auth profile             H session + CSRF · A/S bearer · 5/5 controls
current generated projections    TypeScript + Go PASS / deterministic / compilable
Performance knowledge            2/2 controls PASS
legacy runtime population        0
repository negative controls     1/1
gate                             PASS
```

No Critical finding survived adjudication and no material contradiction remains, so a second Fable round is not warranted.

## D7 boundary

D7 is **NEXT / NOT STARTED**. D6 does not select the server runtime/HTTP mux, database/RLS, transaction implementation, worker/queue realization, session persistence, CSRF bootstrap/rotation, Keycloak realm/deployment or production deployment topology.

D7 may independently evaluate the already pre-vetted candidates from accepted architecture, but none becomes selected merely because D6 mentioned it. Reopen D0–D6 only for a material falsifier at the smallest owning authority.

Do not begin D8–D9, restore retired runtime/manual SDK authority, add Ads-management/autonomous-optimization/AI business authority, or implement Product code before accepted D9.

One coherent gate lands before the next. For task-specific reading, return to [`index.md`](index.md).
