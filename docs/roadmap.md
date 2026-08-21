# Marketplace Central — Roadmap

<!-- program-status-authority -->

> **Role:** sole mutable current-program stage/status/allowed-work/next-action authority. Detailed Product and architecture semantics remain in their owning documents.

## Current checkpoint

| Field | Current value |
| --- | --- |
| Product | **Marketplace Operations Control Plane + Commercial Intelligence** |
| Current stage | **D5 — API — ACCEPTED / CLOSED** |
| Accepted through | **D0–D5**, including the canonical Product OpenAPI Description and executable proof |
| Exact next action | **D6 — Frontend — NEXT / NOT STARTED; await explicit operator authorization to open D6** |
| Canonical Product OAD | `contracts/api/product/openapi.yaml` |
| Product surface | **95 Product operations · 29 ordinary Permissions · Principal kinds H / A / S only** |
| Stable origin | `https://conexus.fun` |
| Product Problem namespace | `https://conexus.fun/marketplace-central/problems/product/{slug}` |
| Technical Problem namespace | `https://conexus.fun/marketplace-central/problems/technical/{surface}/{slug}` when independently contracted |
| Ngrok | preview/tunnel only; never canonical `Problem.type` authority |
| Active runtime baseline | **NONE** |
| Implementation | **BLOCKED UNTIL D9** |

The admitted Product surface is contract scope, not an instruction to build 95 runtime handlers at once.

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
| D6 — Frontend | NEXT / NOT STARTED |
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

## Post-D5 boundary

D6 is the next stage but has **not started**. Opening D6 requires explicit operator authorization and its own bounded authority pack.

Until then:

- do not begin D6 work by implication from D5 closure;
- do not begin D7–D9 or Product implementation;
- do not select runtime/router/database/deployment mechanics;
- do not restore retired legacy runtime/OpenAPI/manual SDK authority;
- reopen D5 only if executable evidence exposes a material contradiction, and then reopen only the smallest implicated authority.

## Progression law

One coherent stage/gate lands before the next begins. A material contradiction reopens only the smallest implicated decision. Product implementation becomes eligible only after accepted D9.

For task-specific reading, return to [`index.md`](index.md).
