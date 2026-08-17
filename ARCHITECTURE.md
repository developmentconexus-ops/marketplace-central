# Marketplace Central Architecture

> **Status:** stable product-level constraints during Architecture Rebaseline  
> **Detailed target architecture:** intentionally under D0–D9 design; see `docs/engineering/rebaseline/README.md`  
> **Last updated:** 2026-08-17

## Purpose

This file contains only architecture constraints that remain stable while the detailed system is re-adjudicated.

It is deliberately **not** a catalog of current modules, tables, routes, frontend packages or processes. Those are current-state evidence in D0 and target-design questions in D1–D9.

A new structural decision becomes durable only after the relevant D-stage is accepted and, when appropriate, an ADR records it.

## Product North Star

Marketplace Central is an internal operations and intelligence system for marketplace commerce, initially Mercado Livre, backed by real Sankhya business-system operational facts.

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
5. **Sankhya is external to MPC; its transport is not business authority.** Business-system access stays behind MPC-owned consumer ports/adapters. D4-B1 makes provider-sanctioned Sankhya API Gateway the default target transport for new MPC↔Sankhya contracts. Direct database access is not a silent fallback: it requires explicit current provider/customer entitlement/support evidence plus a targeted D4-B3 finding that the sanctioned surface cannot meet required correctness, coverage and operational viability. Client-specific direct-DB entitlement/compliance remains unknown until proven. If a direct-DB exception survives B3, its driver/runtime is a later D7 decision rather than an inherited godror default.
6. **Mercado Livre first.** Other marketplace providers are deferred until the Mercado Livre operating loop is coherent and the adapter protocol is proven (ADR-005).
7. **Marketplace provider boundary.** Provider integrations enter through vendor adapters and implement ports owned by consuming business contexts (ADR-033). Provider wire DTO/protocol knowledge stays inside the vendor boundary; exact target package layout remains later realization detail unless a stage explicitly freezes it.
8. **Honest absence.** Unknown facts do not become plausible zero/default values. `internal/kernel/fact` is an accepted primitive for uncertainty where semantically appropriate (ADR-034); D2 decides its correct scope rather than forcing it onto every value.
9. **Exactness where the domain requires it.** Money/tax/cost/pricing values must not lose correctness through floating-point convenience. D2 owns the exact shared/domain representation.
10. **Tenant-ready data isolation is a real invariant.** The exact tenant runtime/RLS model is under D7, but tenant isolation may not depend solely on developers remembering predicates.
11. **External writes are controlled.** A provider write has explicit authority/policy, duplicate protection, auditability and reconciliation. An ambiguous outcome is not blindly retried (ADR-029).
12. **Provider PII is minimized.** Raw external PII is not retained merely because a payload contains it (ADR-025).
13. **Partial observations are honest.** Absence from a partial provider/source pull does not prove closure/deletion (ADR-027).
14. **No compatibility tax without a consumer.** There are no production users requiring current route/schema/package compatibility; hard cutover is allowed under ADR-035.
15. **Git history is history.** Active source/document trees do not keep `old/` copies or parallel legacy roadmaps.

## Architecture Rebaseline authority

ADR-035 governs detailed target-design authority during D0–D9.

The following are **not currently frozen target decisions** even if old code/ADRs once specified them:

- exact context/module set beyond accepted D1 semantic authorities;
- exact database schemas/tables/FKs and physical persistence realization;
- exact scheduler/process/worker topology;
- exact HTTP path namespace and operation set;
- generated server/client technology choice;
- frontend feature/package topology;
- exact transaction/outbox implementation;
- legacy `connectors`, `integrations`, `marketplaces`, `mutations`, `sync`, `internal_read`, `dashboard` or other module structures;
- old manual SDK or proxy-table synchronization mechanisms;
- direct-Oracle/godror as a default Sankhya transport/runtime.

A future session must not infer target architecture from those existing artifacts before their responsible D-stage accepts the relevant meaning.

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

The exact context set and allowed semantic edges are D1 outcomes. Runtime/package realization remains subject to later stages.

## Communication principles

D3 is accepted and canonical:

- **Q:** current producer-owned meaning required for the consumer's current decision;
- **C:** caller asks the owner to accept/perform owner-owned work;
- **E:** already-committed producer-owned fact with a real independent consumer reaction;
- **P:** multiple authorities composed for reading without becoming write authority.

Communication may duplicate, arrive late/out of order, fail or replay without changing business truth. Cross-context SQL/private implementation access is not unnamed communication. D7 chooses transport/runtime realization without moving authority.

## External-integration principles

D4-B1 is accepted and canonical. External integrations obey these constraints:

- consumer context owns semantic meaning/port;
- adapter owns provider/business-system protocol, DTOs, auth and pagination;
- Marketplace Installation/SourceInstance qualify the correct external namespace without becoming credentials or provider IDs;
- namespace mismatch fails closed where authoritative source/provider markers allow it to be detected;
- provider notification/callback is acquisition evidence; authoritative reread establishes current provider state where material;
- point/enumeration/delta/notification coverage is operation-scoped; incomplete/unavailable never becomes plausible absence;
- Integration Support and Provider Effective Capability/Requirement do not become Effective Business Capability;
- external-effect contracts distinguish acceptance/ambiguity from convergence and name an authoritative reread/reconciliation surface;
- current unstable provider/reference behavior must be verified against current official behavior for the concrete decision that depends on it;
- live integration claims require real-dependency evidence, not only mocks;
- no speculative universal provider/integration framework is introduced.

D4-B2/B3/B4 decide the concrete Mercado Livre, Sankhya and market/economics/settlement surfaces still open.

## API and frontend

The current OpenAPI/SDK/routes are **current runtime contract evidence**, not an obligation to preserve every operation.

D5 may delete/redesign legacy operations because there is no external-client compatibility requirement. It will also decide generator/validation authority.

D6 maps every target screen to explicit API/query/mutation ownership and decides frontend feature/package topology.

## Runtime and persistence

D7 decides serving/worker/scheduler/outbox/transaction/cursor/secret/deployment topology and any driver/runtime needed by an explicitly admitted external transport. Do not preserve a dedicated executable, poller or database driver because it exists today.

D2 classifies persistent state and the clean target baseline. Historical migrations do not automatically define the target model.

## Proof bar

A structural rule should, where reasonable, fail at the strongest available boundary:

- illegal private import → compile failure;
- invalid value combination → type/constructor/schema failure;
- foreign/unowned write → structurally unavailable or mechanically blocked;
- contract drift → generation/validation red;
- RLS/isolation bypass → boot/integration failure;
- custom guard → negative fixture proves it fires;
- external namespace mismatch → fail closed before attribution/effect where the source exposes authoritative qualification;
- partial acquisition → cannot pass as complete in contract/integration proof.

A green artifact that did not execute the relevant subject is no proof.

## Current stage

Read `docs/engineering/rebaseline/README.md` for the sole current status and exact next action.
