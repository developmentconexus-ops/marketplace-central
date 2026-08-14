# Marketplace Central Architecture

> **Status:** stable product-level constraints during Architecture Rebaseline  
> **Detailed target architecture:** intentionally under D0–D9 design; see `docs/engineering/rebaseline/README.md`  
> **Last updated:** 2026-08-13

## Purpose

This file contains only architecture constraints that remain stable while the detailed system is re-adjudicated.

It is deliberately **not** a catalog of current modules, tables, routes, frontend packages or processes. Those are current-state evidence in D0 and target-design questions in D1–D9.

A new structural decision becomes durable only after the relevant D-stage is accepted and, when appropriate, an ADR records it.

## Product North Star

Marketplace Central is an internal operations and intelligence system for marketplace commerce, initially Mercado Livre, backed by real Sankhya/Oracle operational facts.

It must support trustworthy flows for:

- internal product identity and source observations;
- marketplace listings/variations and channel observations;
- product↔listing linkage with explicit evidence;
- stock/cost/tax/price semantics;
- safe external actions with preview/policy/audit/reconciliation;
- marketplace orders and realized profitability;
- operational views whose freshness/completeness is honest.

## Stable platform constraints

1. **Independent monorepo.** Marketplace Central remains its own repository/application. Future integration into a broader product does not justify coupling current domain boundaries to another repository.
2. **Go backend is canonical business execution.** `apps/server_core` is the server-side application. Business policy is not duplicated in React.
3. **React frontend is a client, not a second domain authority.** Server state is managed through TanStack Query (ADR-021) unless a material finding explicitly reopens that decision.
4. **PostgreSQL stores MPC-owned canonical state.** External systems are sources/dependencies, not alternate writable application stores.
5. **Sankhya/Oracle is external to MPC.** Oracle reads are behind MPC-owned adapter boundaries (ADR-006); the current canonical driver/runtime is godror/OCI (ADR-007).
6. **Mercado Livre first.** Other marketplace providers are deferred until the Mercado Livre operating loop is coherent and the adapter protocol is proven (ADR-005).
7. **Marketplace provider boundary.** Provider integrations enter through `internal/adapters/marketplace/<vendor>` and implement ports owned by consuming business contexts (ADR-033). Provider wire DTO/protocol knowledge stays inside that vendor tree.
8. **Honest absence.** Unknown facts do not become plausible zero/default values. `internal/kernel/fact` is an accepted primitive for uncertainty where semantically appropriate (ADR-034); D2 decides its correct scope rather than forcing it onto every value.
9. **Exactness where the domain requires it.** Money/tax/cost/pricing values must not lose correctness through floating-point convenience. D2 decides the exact shared/domain representation.
10. **Tenant-ready data isolation is a real invariant.** The exact tenant identity/runtime/RLS model is under D2/D7, but tenant isolation may not depend solely on developers remembering predicates.
11. **External writes are controlled.** A provider write has explicit authority/policy, duplicate protection, auditability and reconciliation. An ambiguous outcome is not blindly retried (ADR-029).
12. **Provider PII is minimized.** Raw external PII is not retained merely because a payload contains it (ADR-025).
13. **Partial observations are honest.** Absence from a partial provider pull does not prove closure/deletion (ADR-027).
14. **No compatibility tax without a consumer.** There are no production users requiring current route/schema/package compatibility; hard cutover is allowed under ADR-035.
15. **Git history is history.** Active source/document trees do not keep `old/` copies or parallel legacy roadmaps.

## Architecture Rebaseline authority

ADR-035 governs detailed target-design authority during D0–D9.

The following are **not currently frozen target decisions** even if old code/ADRs once specified them:

- exact context/module set;
- whether `costing`, `tax`, `profitability`, `intelligence`, `account`, `changecontrol` or other candidates are separate contexts;
- exact database schemas/tables/FKs/migration strategy;
- whether the database is migrated or reset to a clean baseline;
- exact sync/event/projection graph;
- exact scheduler/process/worker topology;
- exact HTTP path namespace and operation set;
- generated server/client technology choice;
- frontend feature/package topology;
- exact transaction/outbox implementation;
- old `connectors`, `integrations`, `marketplaces`, `mutations`, `sync`, `internal_read`, `dashboard` or other legacy module responsibilities;
- old manual SDK or proxy-table synchronization mechanisms.

A future session must not infer target architecture from those existing artifacts before their D-stage is accepted.

## Target reasoning shape

The rebaseline is testing — not blindly accepting — a top-level shape with:

```text
apps/server_core/internal/
  contexts/       business authorities
  adapters/       external-system translations
  kernel/         tiny shared value semantics only
  platform/       technical runtime mechanisms without business policy
  composition/    final assembly only
  views/          rebuildable read projections when justified
```

The exact context set and allowed edges are D1 outcomes. Current `internal/modules` and the two existing contexts are D0 evidence.

## Communication principles under evaluation

D3 will make the full matrix, but these meanings are stable:

- **synchronous capability/query:** caller needs the answer to complete its current decision;
- **event:** another component reacts to a fact that has already committed;
- **projection/view:** multiple authorities are combined for reading without becoming a new business authority.

Cross-context SQL and importing another context’s private implementation are not accepted as unnamed communication mechanisms.

## External-integration principles

D4 maps each Mercado Livre/Sankhya capability in detail. Until then:

- consumer context owns the port;
- adapter owns provider/driver protocol;
- provider DTOs/errors/auth/pagination stay at the provider boundary;
- current provider/reference behavior must be verified against current official behavior when materially unstable;
- live integration claims require real-dependency evidence, not only mocks.

## API and frontend

The current OpenAPI/SDK/routes are **current runtime contract evidence**, not an obligation to preserve every operation.

D5 may delete/redesign legacy operations because there is no external-client compatibility requirement. It will also decide the generator/validation authority.

D6 maps every target screen to explicit API/query/mutation ownership and decides whether current `packages/feature-*` organization is worth retaining.

## Runtime and persistence

D7 decides serving/worker/scheduler/outbox/transaction topology. Do not preserve a dedicated executable or poller because it exists today.

D2 classifies persistent state and decides the clean-baseline vs migration strategy. Historical migrations do not automatically define the target model.

## Proof bar

A structural rule should, where reasonable, fail at the strongest available boundary:

- illegal private import → compile failure;
- invalid value combination → type/constructor/schema failure;
- foreign/unowned write → structurally unavailable or mechanically blocked;
- contract drift → generation/validation red;
- RLS bypass → boot/integration failure;
- custom guard → negative fixture proves it fires.

A green artifact that did not execute the relevant subject is no proof.

## Current stage

Read `docs/engineering/rebaseline/README.md` for the sole current status and exact next action.