# Marketplace Central — Architecture Rebaseline

> **Role:** sole current-program status/router after `AGENTS.md`  
> **Current phase:** **D3 — COMMUNICATION / EVENTS — CLOSURE CANDIDATE / AWAITING OPERATOR RATIFICATION**  
> **Implementation:** BLOCKED until D9 is accepted  
> **Last updated:** 2026-08-16

## 1. Authority path

A fresh session reads, in order:

1. `AGENTS.md`
2. this file
3. `docs/engineering/standards/root-cause-global-maximum-method.md`
4. `ARCHITECTURE.md`
5. `docs/architecture/decisions/README.md`
6. `docs/engineering/rebaseline/D0-PRODUCT-SYSTEM-DEFINITION.md`
7. `docs/engineering/rebaseline/D1-DOMAINS-BOUNDARIES.md`
8. `docs/engineering/rebaseline/D2-IDENTITY-TENANT-DATA-OWNERSHIP.md`
9. `docs/engineering/rebaseline/D3-COMMUNICATION-EVENTS.md`
10. `docs/engineering/rebaseline/EVIDENCE-REGISTER.md`
11. code, OpenAPI, schemas, tests and runtime only as current-state evidence when needed

This file alone answers **where the program is and what happens next**. Stable architecture belongs in `ARCHITECTURE.md`; accepted/current stage semantics belong in D-stage artifacts; Git history is the archive.

Do not reconstruct target authority from memory, legacy package shape, historical plans, `AI-DIALOG.md`, review candidates or stale docs.

## 2. Program state

```text
Documentary / governance cleanup — DONE
  ↓
D0 — Product / System Definition — CLOSED / ACCEPTED
  ↓
D1 — Domains / Boundaries — CLOSED / ACCEPTED
  ↓
D2 — Identity / Tenant / Data Ownership — CLOSED / ACCEPTED
  ↓
D3 — Communication / Events — CLOSURE CANDIDATE / AWAITING OPERATOR RATIFICATION
  ├─ B1 Communication Topology & Edge Matrix — ACCEPTED / CANONICAL
  ├─ B2 Communication Contract & Failure Semantics — ACCEPTED / CANONICAL
  └─ Final Global Coherence + YAGNI / Overengineering / Future-Cost — COMPLETED / PASS
  ↓
D4 — External Integrations — BLOCKED UNTIL D3 RATIFICATION
  ↓
D5 — API
  ↓
D6 — Frontend
  ↓
D7 — Runtime / Jobs / Transactions
  ↓
D8 — Golden Flows
  ↓
D9 — Adversarial Architecture Review
  ↓
Implementation DAG / Plan
  ↓
Implementation
```

Product implementation remains blocked until D9 is accepted.

## 3. Accepted baseline

### D0 — CLOSED

Product 1.0 is **Marketplace Operations + Commercial Intelligence**. MPC is the marketplace operations control plane: external systems retain authority for facts/processes inherently theirs while MPC owns the cross-system marketplace operating semantics needed to observe, understand, decide, execute, verify and reconcile.

D0 authority and non-goals are defined only in `D0-PRODUCT-SYSTEM-DEFINITION.md`.

### D1 — CLOSED

`D1-DOMAINS-BOUNDARIES.md` is the accepted D1 authority. It defines 12 business boundaries, explicit ownership/non-ownership, cross-cutting non-domain treatment, legal semantic authority edges, forbidden boundary violations, legacy semantic disposition, D2–D7 defers and reopen triggers.

D1 closure included operator adjudication, independent Fable review and final Global Coherence + YAGNI / Overengineering / Future-Cost review. The 12 D1 business boundaries do **not** imply 12 services, databases, processes or deployments.

### D2 — CLOSED

`D2-IDENTITY-TENANT-DATA-OWNERSHIP.md` is the accepted D2 authority. B1+B2 were operator-approved, independently challenged, consolidated and explicitly ratified as a whole after the final Global Coherence + YAGNI / Overengineering / Future-Cost review.

Accepted D2 direction includes:

- canonical identity follows semantic authority; MPC-owned IDs are opaque/stable/non-reusable;
- MPC-owned Organization, Marketplace Installation, Selling Entity, Inventory Source, Fulfillment Node and human/automation/system Principal semantics;
- Organization as canonical tenant/isolation root with no duplicate target `Tenant` identity;
- source-qualified external Product, Listing/Variation, Sale/Order, Shipment and native financial identities;
- minimal SourceInstance identity while D4 owns concrete source contracts/capabilities/config/credentials;
- stable domain-local identity for material Business Intents that participate in authorization/external effects/convergence/history;
- MPC-owned Post-Sale Resolution, Economic Attribution/Reconciliation, Operational Work and Authorization Decision/Grant semantics;
- external OIDC human AuthN with MPC-owned Principal/binding/Membership/AccessRole/Permission ordinary-access state and a strict fence against business-action authority leaking into identity/access;
- one D1 semantic write authority per canonical business meaning; typed references/snapshots/projections do not become current write authority;
- exact Money/rate/material-quantity semantics, bounded `Fact<T>` scope, provenance and distinct material time meanings;
- clean target persistence baseline with **no legacy-data migration or archival requirement** for the pre-rebaseline MPC database;
- old ADR structures do not carry forward by inheritance; D2 adjudicated ADR-011/012/022/028/031 and defined safe rehoming gates/new ADR-036+ transition.

### D3 — CLOSURE CANDIDATE

`D3-COMMUNICATION-EVENTS.md` is the complete current D3 authority.

Accepted D3 direction includes:

- semantic hybrid **Q/C/E/P** communication rather than all-sync or all-event-driven ideology;
- Q for current owner meaning needed to decide now;
- C for asking an owner to perform/accept owner-owned work;
- E only for already-committed producer-owned facts with real independent consumer reactions;
- P for read-only multi-authority composition;
- progression events wake consumers while current owner revalidation remains required when currentness is material;
- evidence-consuming edges preserve/recover material occurrences from the smallest sufficient durable authority rather than latest mutable state or universal event sourcing;
- consequential event propagation is recoverable rather than silently lossy;
- Organization scope remains explicit across communication/recovery and is never inferred from Installation/SourceInstance/provider IDs;
- duplicate delivery is safe through consumer semantic idempotency; transport dedupe is not business correctness;
- arrival order is not business order; no global sequence/order/version is assumed;
- evidence occurrences use bounded owner/source-defined stable discrimination only where same-vs-distinct correctness requires it; no universal EventID is created;
- missed-reaction conclusions stay with the domain whose progression/convergence is missing; automatic recovery does not automatically create Work;
- capability acceptance distinguishes accepted/rejected/pending/ambiguous where applicable and ambiguous retry reconciles by stable Organization-scoped semantic anchor rather than generic CommandID;
- Q preserves known/known-empty/unknown/unavailable semantics and material freshness provenance;
- replay/redelivery cannot rewrite history or blindly repeat external effects;
- projections remain rebuildable read state and never make transport logs historical authority;
- multi-target communication preserves intended/authorized/attempted scope and granular member outcomes;
- provider webhook/poll/callback evidence is not itself a D3 domain event;
- no cross-owner distributed transaction is required;
- shared external-effect safety mechanics verify proofs but do not own business disposition/policy/authorization;
- incompatible communication-contract cutover preserves still-required recoverable reactions without requiring permanent multi-version support;
- ADR-019 and ADR-024 are fully rehomed and historical; ADR-018/026 retain only D7 residue;
- no generic Event Bus/Command Bus/Workflow engine/event sourcing/universal CQRS/exactly-once/global ordering/schema-registry/broker/outbox/runtime topology/microservice split is chosen by D3.

D3-B1 and B2 were independently challenged and operator-ratified as batches. The final D3 Global Coherence + YAGNI / Overengineering / Future-Cost review completed with **CURRENT STRUCTURE CONFIRMED**, no material correction, no B3 and no D0/D1/D2 reopen.

D3 is **not closed as a whole until the operator explicitly ratifies this closure candidate**.

## 4. Engineering method and repo lifecycle

Engineering reasoning follows the **DevelopmentConexus Engineering Method** identified in `AGENTS.md`; the local file in this authority path is only the consumed context copy.

This router defines the Marketplace Central D0–D9 status/lifecycle and allowed work. It is repo-specific specialization, not a second organizational engineering method. Conflicts inside the organizational method's scope must be surfaced, never silently reinterpreted.

## 5. What is prohibited now

Until D3 is explicitly ratified as a whole:

- do not begin D4–D9 target design prematurely;
- do not implement product architecture/features;
- do not silently alter accepted D0/D1/D2 or canonical D3-B1/B2 authority;
- do not create semantic dependencies outside D1 by hiding them in events, APIs, queues, projections or databases;
- do not choose provider/ERP transport contracts/credentials yet;
- do not choose HTTP/frontend/runtime topology prematurely;
- do not make event transport a business/historical authority or require universal event sourcing;
- do not treat `AI-DIALOG.md`, deleted review candidates or reviewer/chat summaries as target authority.

Existing code/module/context/schema shape remains current-state evidence only.

## 6. Exact next action

**Operator ratifies or amends D3 as a whole.**

The complete closure candidate is `docs/engineering/rebaseline/D3-COMMUNICATION-EVENTS.md`, including accepted B1+B2 and the final Global Coherence + YAGNI / Overengineering / Future-Cost review.

If the operator approves the whole:

1. mark D3 `CLOSED / ACCEPTED`;
2. advance this router so **D4 — External Integrations** is the exact next stage;
3. do not begin implementation; implementation remains blocked until D9.

If the operator finds a material issue, reopen only the implicated D3 decision rather than re-running the whole stage.

## 7. Fresh-session success test

A fresh session should conclude that:

- D0, D1 and D2 are **CLOSED / ACCEPTED**;
- D1 defines exactly 12 semantic business boundaries but no runtime topology;
- D2 fixes canonical/external identities, tenant/isolation semantics, persistent ownership and shared value/knowledge/time semantics;
- D3 is a **CLOSURE CANDIDATE / AWAITING OPERATOR RATIFICATION**;
- D3-B1 and D3-B2 are **ACCEPTED / CANONICAL**;
- no B3 is required and final D3 Global Coherence completed with no material contradiction;
- current modules/contexts/schema/legacy ADRs remain evidence, not target authority by inheritance;
- implementation remains blocked until D9;
- D4 remains blocked until whole-stage D3 ratification;
- the exact next action is explicit operator ratification of D3 as a whole.

If it cannot, the authority path is incomplete or contradictory.
