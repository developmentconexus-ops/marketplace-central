# Marketplace Central — Architecture Rebaseline

> **Role:** sole current-program status/router after `AGENTS.md`  
> **Current phase:** **D3 — COMMUNICATION / EVENTS — OPEN / IN PROGRESS; D3-B1 ACCEPTED, D3-B2 is the exact next batch**  
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
D3 — Communication / Events — OPEN / IN PROGRESS
  ├─ B1 Communication Topology & Edge Matrix — ACCEPTED
  └─ B2 Communication Contract & Failure Semantics — NEXT
  ↓
D4 — External Integrations — NOT OPEN
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

### D3 — OPEN / B1 ACCEPTED

`D3-COMMUNICATION-EVENTS.md` is the active D3 authority.

Accepted D3-B1 direction includes:

- a **semantic hybrid Q/C/E/P** communication model rather than all-sync or all-event-driven ideology;
- Q for current owner meaning needed to decide now;
- C for asking an owner to perform/accept owner-owned work;
- E only for already-committed producer-owned facts with real independent consumer reactions;
- P for read-only multi-authority composition;
- progression events wake consumers, but current owner revalidation remains required when currentness is material;
- evidence-consuming edges preserve/recover each material occurrence from the smallest sufficient durable authority rather than relying on latest mutable state or universal event sourcing;
- consequential event propagation is detectable/recoverable rather than silently lossy;
- accepted feed-forward edges remain Q-only by baseline unless a real autonomous reaction is proven;
- Sales fan-out to Materialization/Fulfillment/Economics is event-based; Sales→Post-Sale events are limited to post-sale-relevant facts;
- Materialization⇄Fulfillment remains a semantic business cycle with separate write authorities and committed checkpoints, not bilateral mutation;
- Governance uses C/Q from action owner and E/Q back for later committed decisions without becoming executor;
- Operational Work consumes source-owned actionable-condition facts while retaining Work-obligation representation/lifecycle authority;
- Post-Sale coordinates through C/E/Q without absorbing consequence-owner semantics;
- ordinary identity/access correctness remains Q-based; revocation cannot depend solely on eventual events;
- projections are rebuildable from owner state/evidence and never make event transport the sole history/rebuild authority;
- provider webhook/poll/callback evidence is not itself a D3 domain event;
- no cross-owner distributed transaction is required;
- shared external-effect safety mechanics verify proofs but do not own business disposition/policy/authorization;
- D3 B1 partially adjudicated legacy ADR-018/019/024/026 while leaving B2/D7 residues explicit;
- no generic Event Bus/Command Bus/Workflow engine/event sourcing/universal CQRS/broker/outbox/runtime topology/microservice split is chosen by B1.

## 4. Engineering method and repo lifecycle

Engineering reasoning follows the **DevelopmentConexus Engineering Method** identified in `AGENTS.md`; the local file in this authority path is only the consumed context copy.

This router defines the Marketplace Central D0–D9 status/lifecycle and allowed work. It is repo-specific specialization, not a second organizational engineering method. Conflicts inside the organizational method's scope must be surfaced, never silently reinterpreted.

## 5. What is prohibited now

While D3 remains open:

- do not begin D4–D9 target design prematurely;
- do not implement product architecture/features;
- do not let legacy IDs/tables/ADRs imply target identity, persistence ownership, communication architecture or target contracts;
- do not silently alter accepted D0/D1/D2 or accepted D3-B1 authority while choosing B2 communication contracts;
- do not create a semantic dependency outside D1 by hiding it in an event, API, queue, projection or database;
- do not choose provider/ERP transport contracts or credentials;
- do not choose HTTP/frontend/runtime topology;
- do not choose workers, brokers, outbox/transaction implementation, locks, RLS enforcement or deployment topology;
- do not make event transport a business/historical authority or require universal event sourcing;
- do not treat `AI-DIALOG.md`, review candidates or reviewer/chat summaries as target authority.

Existing code/module/context/schema shape remains current-state evidence only.

## 6. Exact next action

**Open D3-B2 — Communication Contract & Failure Semantics.**

B1 is canonical and must not be reopened merely to choose technology.

B2 must define, only to the semantic depth required before D4/D7:

- Organization scope;
- actor/Principal attribution where material;
- domain-local Intent references;
- event occurrence identity/producer ownership;
- causation/correlation;
- provenance/material time;
- duplicate/idempotency semantics;
- ordering/late-delivery/anti-regression assumptions;
- missed-delivery recovery and replay/reconciliation expectations;
- progression current-reread vs evidence occurrence recovery;
- capability outcomes/ambiguity;
- projection rebuild completeness/fail-honesty;
- multi-target granular outcomes where material;
- remaining D3 semantic residues of ADR-019/024.

B2 must **not** select broker, outbox table, queue, worker, scheduler, lock, transaction or deployment topology; those remain D7. Concrete provider/source acquisition contracts remain D4.

After B2 is operator-approved, independently challenged and consolidated, D3 proceeds to final Global Coherence + YAGNI / Overengineering / Future-Cost review before whole-stage ratification.

If D3 discovers a genuinely necessary semantic dependency not allowed by D1, reopen only the implicated D1 decision rather than hiding the dependency in mechanism.

Do not advance D4 or product implementation before D3 is accepted as a whole.

## 7. Fresh-session success test

A fresh session should conclude that:

- D0, D1 and D2 are **CLOSED / ACCEPTED**;
- D1 defines exactly 12 semantic business boundaries but no runtime topology;
- D2 fixes canonical/external identities, tenant/isolation semantics, persistent ownership and shared value/knowledge/time semantics;
- D3 is **OPEN / IN PROGRESS**;
- **D3-B1 is ACCEPTED / CANONICAL** and defines the semantic Q/C/E/P topology/edge matrix;
- **D3-B2 is NEXT** and owns communication contract/failure semantics;
- current modules/contexts/schema/legacy ADRs remain evidence, not target authority by inheritance;
- implementation remains blocked until D9;
- D4 is not yet open.

If it cannot, the authority path is incomplete or contradictory.
