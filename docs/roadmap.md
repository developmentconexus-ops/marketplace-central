# Marketplace Central — Roadmap

<!-- program-status-authority -->

> **Role:** sole mutable current-program stage/status/allowed-work/next-action authority. Detailed semantics remain in their owning documents.

## Current checkpoint

| Field | Current value |
| --- | --- |
| Product | **Marketplace Operations Control Plane + Commercial Intelligence** |
| Current stage | **D6 — Frontend — OPEN / ACTIVE** |
| Accepted baseline | **D0–D5 ACCEPTED / CLOSED** |
| D6-discovered bounded amendment | **D6-R1 Marketplace Performance Intelligence — OPERATOR-APPROVED / FABLE ACCEPT WITH BOUNDED FIXES / FIXES ADJUDICATED IN CURRENT CANDIDATE** |
| D6-B1 interaction proof | **OPERATOR-RATIFIED — corrected Portuguese 99-operation interaction map + low-fidelity wireframes accepted** |
| Exact next action | **D6-B2 — adjudicate the smallest frontend feature/package topology and exact dependency needs from current official evidence; do not select D7 server/runtime/database/deployment mechanics** |
| Canonical Product OAD candidate | `contracts/api/product/openapi.yaml` |
| Product surface in current D6 candidate | **99 Product operations · 30 ordinary Permissions · Principal kinds H / A / S only** |
| Stable origin | `https://conexus.fun` |
| Active runtime baseline | **NONE** |
| Implementation | **BLOCKED UNTIL D9** |

The Product operation count is contract scope, not a handler/screen count.

## Stage progression

| Stage | Status |
| --- | --- |
| D0 — Product / System Definition | ACCEPTED / CLOSED; bounded D6-R1 amendment approved in current candidate |
| D1 — Domains / Boundaries | ACCEPTED / CLOSED; bounded D6-R1 amendment approved in current candidate |
| D2 — Identity / Tenant / Data Ownership | ACCEPTED / CLOSED; D2-R1 + D2-R2 bounded amendments approved in current candidate |
| D3 — Communication / Events | ACCEPTED / CLOSED; bounded D6-R1 Q/P amendment approved in current candidate |
| D4 — External Integrations | ACCEPTED / CLOSED; bounded Mercado Livre Performance evidence amendment approved in current candidate |
| D4-R1 — Publication Input / Listing Authoring | ACCEPTED / CANONICAL |
| D5 — API | ACCEPTED / CLOSED; bounded D6-R1 99/30 wire amendment proved in current candidate |
| D6 — Frontend | OPEN / ACTIVE; D6-B1 OPERATOR-RATIFIED; D6-B2 NEXT / ACTIVE DECISION |
| D7 — Runtime / Jobs / Transactions | BLOCKED |
| D8 — Golden Flows | BLOCKED |
| D9 — Adversarial Architecture Review | BLOCKED |
| Implementation | BLOCKED UNTIL D9 |

## D5 baseline and D6-R1 amendment

`main @ 9d2c81e175bc39ac388c9d8924ddad21f2a86480` remains the accepted D5 closeout baseline at 95 operations / 29 ordinary Permissions. D6 exposed a material strategic-analysis gap and the operator approved the smallest bounded parent repair recorded in [D6-R1 Marketplace Performance Intelligence](engineering/rebaseline/D6-R1-MARKETPLACE-PERFORMANCE-INTELLIGENCE.md).

The current D6 candidate adds only:

- one read/derive boundary, **Marketplace Performance Intelligence**;
- one ordinary Permission, `performance.read`;
- four H/A/S Product Qs: marketplace summary, Listing-performance list, Listing-performance point read and Retail Media performance list;
- D2-R2 historical Performance-evidence custody so provider retention does not become MPC historical retention;
- bounded Mercado Livre Visits + current Product Ads evidence semantics and a Technical Non-Product advertiser-binding ceremony under `portfolio.manage`;
- no Performance mutation, Ads management, generic Analytics/Metric API, Product P operation, AI/MCP authority or D7 mechanism.

Executable proof on the repaired OAD establishes:

```text
accepted D5 baseline non-regression  PASS  95/95 ops · 29/29 Permissions · 26/26 List/Search · 12/12 negative controls
current D6-R1 Product candidate      PASS  99/99 ops · 30/30 Permissions · 28/28 List/Search · 7/7 Performance negative controls
Performance knowledge wire          PASS  available requires measure · unavailable forbids measure · 2/2 negative controls
Principal kinds                      H / A / S
TypeScript + Go projections          deterministic / compilable
legacy runtime population            0
runtime schema enforcement           NOT_CLAIMED_D7
router selection                     NONE_D7
```

Independent Fable review concluded **ACCEPT WITH BOUNDED FIXES** and found no reason to reconstruct or reopen the 13th boundary, D2-R2 custody, the four-Q surface, `performance.read` or the proof approach. GPT adjudicated F-1…F-8 as bounded corrections only: frontend proof clarity, wire knowledge precision, baseline-fixture labeling, count-independent D2-R1 wording, honest repository negative-control reporting, routing freshness, and a D7 measure-by-scope obligation. No new Product capability or mechanism was admitted by the review.

The amendment is bounded; it is not permission to generally reopen D0–D5.

## D6-B1 ratified interaction authority

[D6 Frontend](engineering/rebaseline/D6-FRONTEND.md) remains the active frontend authority. Binding shell decisions:

- Organization is the global workspace/isolation context;
- primary navigation is task-oriented, not D1 package taxonomy;
- Marketplace Installation is explicit contextual navigation state and never ambient authority;
- Settings is low-frequency grouping, not a business owner;
- read-only screen composition never becomes write authority;
- route/button visibility is usability, never authorization;
- server state, URL/navigation state, form draft and ephemeral UI state remain distinct;
- unknown/partial/unavailable/stale and accepted/pending/rejected/ambiguous remain honest;
- Product/Technical Ingress separation remains binding.

D6-R1 refines the strategic group:

```text
ESTRATÉGIA E INTELIGÊNCIA
  Performance
    Resumo
    Publicações
    Mídia
  Mercado
  Economia
```

The corrected [interaction map](engineering/rebaseline/D6-B1-INTERACTION-MAP.md) maps all 99 Product operations and the Portuguese [low-fidelity wireframe proof](../qualification/d6-wireframes/index.html) exercises the approved shell plus the strategic Performance workspace. Following Fable review, the wireframe no longer implies an uncontracted `signals[]` surface, does not show a delta from partial evidence, and explicitly demonstrates both `insufficient_evidence` and `not_comparable` comparison refusal states. The operator ratified this corrected D6-B1 proof on 2026-08-21.

Performance requires an exact Marketplace Installation. Future multi-marketplace views may compose per-Installation results but may not add provider metrics as if their measurement bases were automatically equivalent.

## D6-B2 decision boundary

D6-B2 may decide only the smallest frontend realization contract required to implement the ratified interaction model later, including feature/package boundaries and concrete frontend-only dependency needs when current evidence proves them necessary.

D6-B2 must preserve:

- canonical generated Product client/types rather than hand-written API DTO duplication;
- TanStack Query as server-state owner unless material evidence reopens it;
- URL/navigation, form draft and ephemeral UI as distinct state classes;
- task-oriented UX while package boundaries remain semantic/feature-oriented and independently understandable;
- no generic action/workflow/client-domain framework;
- no microfrontends, offline-first, websocket/event-stream, universal design-system platform, generic analytics layer or AI-specific frontend architecture without a proven consumer.

## D6 blocked/deferred boundary

While D6 is active:

- do not begin D7–D9;
- do not select server runtime/router/database/transaction/worker/deployment mechanics;
- do not implement Product code;
- do not restore retired runtime/frontend/manual SDK authority;
- do not add Ads management, autonomous commercial optimization or AI/MCP authority;
- reopen an accepted earlier decision only for a material falsifier and only at the smallest responsible authority.

D7 inherits one explicit Performance realization obligation from Fable review: prove real Mercado Livre measure availability per exact Retail Media scope and report unsupported/unavailable honestly; do not assume campaign/item/catalog/family measurement symmetry.

## Progression law

One coherent stage/gate lands before the next begins. Product implementation becomes eligible only after accepted D9. PR #54 remains Draft and is not authorized for merge without explicit operator authorization.

For task-specific reading, return to [`index.md`](index.md).