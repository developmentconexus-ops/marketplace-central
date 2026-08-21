# Marketplace Central — Roadmap

<!-- program-status-authority -->

> **Role:** sole mutable current-stage/status/allowed-work/next-action authority. Detailed semantics live in routed owning documents.

## Current checkpoint

| Field | Current value |
| --- | --- |
| Product | **Marketplace Operations Control Plane + Commercial Intelligence** |
| Current stage | **D6 — Frontend — OPEN / ACTIVE** |
| Accepted baseline | **D0–D5 ACCEPTED / CLOSED** |
| D6-R1 | **Marketplace Performance Intelligence — OPERATOR-APPROVED / FABLE REVIEWED / BOUNDED FIXES ADJUDICATED** |
| D6-B1 | **OPERATOR-RATIFIED — corrected Portuguese interaction map + wireframes** |
| Cross-repo review | **Marketplace Central ↔ MetalDocs — OPERATOR-APPROVED / ADJUDICATED** |
| D5-R1 auth correction | **OPERATOR-APPROVED / EXECUTABLE PROOF PASS** — [Human Browser Authentication Correction](engineering/rebaseline/D5-R1-HUMAN-BROWSER-AUTHENTICATION.md) |
| D6-B2 | **OPERATOR-RATIFIED — frontend realization topology + dependency profile** |
| Exact next action | **Run the isolated independent final D6 challenge against the exact candidate; adjudicate findings before any D6 closeout or D7 opening** |
| Canonical Product OAD | `contracts/api/product/openapi.yaml` |
| Product surface | **99 Product operations · 30 ordinary Permissions · Principal kinds H / A / S only** |
| Stable origin | `https://conexus.fun` |
| Active runtime baseline | **NONE** |
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
| D5 — API | ACCEPTED / CLOSED — bounded D5-R1 human-browser auth carrier correction proved in current D6 candidate |
| D6 — Frontend | OPEN / ACTIVE — B1 + B2 operator-ratified; final independent D6 challenge next |
| D7 — Runtime / Jobs / Transactions | BLOCKED |
| D8 — Golden Flows | BLOCKED |
| D9 — Adversarial Architecture Review | BLOCKED |
| Implementation | BLOCKED UNTIL D9 |

## Current accepted Product authority

`main @ 9d2c81e175bc39ac388c9d8924ddad21f2a86480` remains the accepted D5 closeout baseline at 95 operations / 29 Permissions. The bounded [D6-R1 Marketplace Performance Intelligence](engineering/rebaseline/D6-R1-MARKETPLACE-PERFORMANCE-INTELLIGENCE.md) repair produces the current 99/30 Product candidate with four Performance Qs and `performance.read`.

The operator-approved [D5-R1 Human Browser Authentication Correction](engineering/rebaseline/D5-R1-HUMAN-BROWSER-AUTHENTICATION.md) changes authentication carriers only:

```text
H browser  -> server-side OIDC login -> Secure HttpOnly application session + CSRF on unsafe requests
A / S      -> Client Credentials -> audience-bound bearer
```

It adds/removes no operation, Permission, Principal kind, owner or business semantic. Keycloak remains the preferred first OIDC provider candidate; exact session/CSRF/Keycloak/runtime realization remains D7.

Fresh executable proof establishes:

```text
accepted D5 baseline             95/95 operations · 29/29 Permissions · 12/12 controls
pre-auth D6-R1 surface           99/99 operations · 30/30 Permissions · 7/7 Performance controls
current auth profile             H session + CSRF · A/S bearer · 5/5 auth controls
current generated projections    TypeScript + Go PASS
Performance knowledge            2/2 controls PASS
legacy runtime population        0
```

## D6-B1 ratified interaction authority

[D6 Frontend](engineering/rebaseline/D6-FRONTEND.md) owns the detailed frontend laws. D6-B1 preserves explicit Organization/Installation context, task-oriented Portuguese navigation, Product/Technical-Ingress separation, read composition without write authority and honest knowledge/outcome states.

The [interaction map](engineering/rebaseline/D6-B1-INTERACTION-MAP.md) covers all 99 operations and the Portuguese [wireframes](../qualification/d6-wireframes/index.html) are operator-ratified.

## D6-B2 ratified frontend realization

D6-B2 selects the smallest frontend profile proven necessary by the ratified interactions and the Marketplace Central ↔ MetalDocs cross-repo adjudication:

```text
React + TypeScript strict
TanStack Query              server-state authority
TanStack Router             typed/validated URL/navigation state
openapi-typescript          generated Product wire shapes
openapi-fetch               thin OAD-bound Product HTTP client
```

Topology law:

```text
human lens/flow UI
  -> stateless owner/operation Product adapters
  -> one thin openapi-fetch transport
  -> generated OpenAPI shapes
```

UI follows stable human lenses/flows; API adapters follow Product ownership/families without becoming client-side domain authorities. Routes remain thin and never create a second data cache. Cross-owner composition stays in user-facing lenses; owner adapters do not import each other's internals.

The human browser consumes D5-R1 session + CSRF and never owns OIDC access/refresh tokens. No Redux/Zustand server mirror, generated generic query/workflow layer, Axios, universal form/schema/UI framework, microfrontend, SSR-by-fashion, offline-first or realtime architecture is admitted without a proven consumer.

Exact dependency versions and the import-boundary lint mechanism remain implementation details; the architectural dependency direction is binding and must become mechanically default-deny when a real frontend source tree exists.

## Blocked / deferred boundary

D6-B2 does **not** select D7 server runtime/HTTP mux/database/RLS/transactions/River/deployment/Keycloak realm topology.

D7 inherits pre-vetted candidates from the cross-repo review but none are selected by D6: modular-monolith class, pgx/pgxpool, structural tenant isolation/RLS proof, River-first durable-work falsification, OpenTelemetry/OTLP/slog, sqlc, tern and real-dependency test tooling.

Do not begin D7–D9, restore retired runtime/manual SDK authority, add Ads management/autonomous optimization/AI authority, or merge PR #54 without explicit operator authorization. Reopen accepted authority only for a material falsifier at the smallest owning stage.

## Final D6 challenge gate

The current D6 candidate now contains one coherent interaction + realization design. Before D6 closeout it must receive an independent challenge from the **exact candidate HEAD** using the isolated `review/*` protocol in `docs/development/engineering-rules.md`.

The review must attack at least:

- D5-R1 human-session / machine-bearer consumption;
- D6-R1 Performance authority and frontend use;
- D6-B1 99-operation interaction coverage;
- D6-B2 lens/flow + owner-adapter topology;
- TanStack Router / Query ownership separation;
- `openapi-typescript` + `openapi-fetch` wire discipline;
- retry/idempotency/concurrency preservation;
- permission/knowledge/identity leakage;
- D7 mechanism leakage and YAGNI.

Reviewer output is Evidence, not authority. Round 2 occurs only for a surviving material contradiction after GPT adjudication.

One coherent gate lands before the next. For task-specific reading, return to [`index.md`](index.md).
