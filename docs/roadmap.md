# Marketplace Central — Roadmap

<!-- program-status-authority -->

> **Role:** sole mutable current-program stage/status/allowed-work/next-action authority. Detailed Product and architecture semantics remain in their owning documents.

## Current checkpoint

| Field | Current value |
| --- | --- |
| Product | **Marketplace Operations Control Plane + Commercial Intelligence** |
| Current stage | **D6 — Frontend — OPEN / ACTIVE** |
| Accepted through | **D0–D5**, including the canonical Product OpenAPI Description and executable proof |
| Exact next action | **D6-B1 — operator-review/adjudicate the 95-operation interaction map + ten-state low-fidelity HTML wireframe proof; only if approved proceed to frontend topology/dependency adjudication** |
| Canonical Product OAD | `contracts/api/product/openapi.yaml` |
| Product surface | **95 Product operations · 29 ordinary Permissions · Principal kinds H / A / S only** |
| Stable origin | `https://conexus.fun` |
| Product Problem namespace | `https://conexus.fun/marketplace-central/problems/product/{slug}` |
| Technical Problem namespace | `https://conexus.fun/marketplace-central/problems/technical/{surface}/{slug}` when independently contracted |
| Ngrok | preview/tunnel only; never canonical `Problem.type` authority |
| Active runtime baseline | **NONE** |
| Implementation | **BLOCKED UNTIL D9** |

The admitted Product surface is contract scope, not an instruction to build 95 runtime handlers or 95 frontend screens.

## Stage progression

| Stage | Status |
| --- | --- |
| D0 — Product / System Definition | ACCEPTED / CLOSED |
| D1 — Domains / Boundaries | ACCEPTED / CLOSED |
| D2 — Identity / Tenant / Data Ownership | ACCEPTED / CLOSED |
| D3 — Communication / Events | ACCEPTED / CLOSED |
| D4 — External Integrations | ACCEPTED / CLOSED |
| D4-R1 — Publication Input / Listing Authoring | ACCEPTED / CANONICAL |
| D5 — API | ACCEPTED / CLOSED |
| D6 — Frontend | OPEN / ACTIVE |
| D7 — Runtime / Jobs / Transactions | BLOCKED |
| D8 — Golden Flows | BLOCKED |
| D9 — Adversarial Architecture Review | BLOCKED |
| Implementation | BLOCKED UNTIL D9 |

## D5 closure

D5 is closed on the canonical Product OAD landed in `main` after executable proof, independent Fable challenge and review adjudication.

Closure preserves:

- exactly 95 Product operations and 29 ordinary Permissions;
- Principal kinds H / A / S only;
- Product Problem identity under the stable `https://conexus.fun` origin;
- Product/Technical Ingress separation and no provider/OAuth ingress in the Product OAD;
- generated TypeScript and Go projections as temporary derived proof artifacts, never a second authority;
- Organization privacy-preserving `404`, required idempotency carriers, strong validators and accepted collection grammar;
- zero population of the retired legacy runtime/OpenAPI/manual SDK;
- no Product runtime, router, database or deployment selection.

D5 proof demonstrates source/wire/generator compatibility only. Runtime schema rejection, concrete router compatibility for canonical `:verb` paths, supported Go runtime floor, persistence and transaction mechanics remain D7 obligations.

## D6 opening

The operator explicitly authorized opening D6 on 2026-08-21 from the accepted D5 checkpoint.

D6 begins with [D6-B1 — Frontend Interaction & Authority Model](engineering/rebaseline/D6-FRONTEND.md). B1 is **OPEN / ACTIVE**, not accepted merely because D6 was authorized to begin.

Current D6 boundary:

- React remains the accepted Product API client technology; business policy is not duplicated in React;
- TanStack Query remains the accepted server-state client unless material evidence explicitly reopens that decision;
- canonical Product OAD remains the wire authority; D6 does not create a screen-shaped/BFF second Product contract;
- server state, URL/navigation state, form draft state and ephemeral UI state remain distinct unless evidence proves another class is necessary;
- target screens/interactions must map to accepted semantic owners and explicit Product operations/queries/capabilities plus ordinary Permission where applicable;
- known-empty / unknown / unavailable / partial / materially stale states and accepted / pending / rejected / ambiguous consequential outcomes remain honest where reachable;
- route/button visibility is usability, never authorization;
- consequential retry UX may not weaken idempotency, concurrency or no-blind-replay semantics;
- provider/OAuth/technical ingress is not a Product frontend shortcut;
- frontend dependency/package topology is selected only after a concrete D6 property requires it and evidence supports the smallest fit.

D6 opening does **not** authorize Product implementation.

## D6-B1 current result

The operator has approved the global App Shell / information architecture:

> **Organization-global + task-oriented primary navigation + explicit contextual Marketplace Installation where the Product contract admits that dimension + low-frequency Settings separated from routine operation + read-only Overview composition.**

The current D6-B1 proof package includes [D6-B1 Frontend Interaction Map](engineering/rebaseline/D6-B1-INTERACTION-MAP.md), which derives ten user flows, 32 screen/route states and explicit interaction homes for all 95 admitted Product operations.

Three bounded frontend-coverage findings have been adjudicated:

1. **Available channels / authorization:** no generic Integration/Channel Catalog Product operation and no fake Product `ConnectMarketplace` capability. Product 1.0 remains Mercado Livre-only for current connectability; D5 Technical Ingress remains the provider authorization ceremony.
2. **Presentation identity:** [D2-R1 Presentation Identity](engineering/rebaseline/D2-R1-PRESENTATION-IDENTITY.md) plus the bounded D5 wire correction provides non-authoritative human-readable labels while Organization/Principal/AccessRole IDs/keys remain canonical authority. The canonical OAD proof remains green at 95 operations / 29 Permissions / H-A-S.
3. **SourceInstance discovery:** `SearchSourceProductsForMarketplace` treats `source_instance_id` as an optional operation-local narrowing filter. Omission means bounded search across current Organization-scoped Readiness-admitted/configured SourceInstances and never ambient/default-source selection; every result remains source-qualified. Point readiness/requirements reads still require the exact SourceInstance selected from the search result. No Source Registry/Product operation is added, and the canonical OAD proof remains green.

The low-fidelity proof now exists at `qualification/d6-wireframes/index.html`. It contains ten navigable grayscale/static wireframe states and passed the internal interaction-map negative-control challenge without exposing a new Product/API prerequisite. It is design evidence only, not Product/runtime/browser proof.

D6-B1 is **not yet ratified**. The exact current gate is operator review of the interaction map and wireframes. Frontend topology/dependency selection begins only if that interaction model is approved.

## D6 blocked/deferred boundary

While D6 is active:

- do not begin D7–D9;
- do not select server runtime/router/database/transaction/worker/deployment mechanics;
- do not implement Product code;
- do not restore retired runtime, frontend, OpenAPI or manual SDK authority;
- do not reopen D0–D5 for preference, topology convenience or screen design;
- reopen an accepted earlier decision only if executable/material evidence exposes a contradiction, and reopen only the smallest implicated authority.

## Progression law

One coherent stage/gate lands before the next begins. A material contradiction reopens only the smallest implicated decision. Product implementation becomes eligible only after accepted D9.

For task-specific reading, return to [`index.md`](index.md).