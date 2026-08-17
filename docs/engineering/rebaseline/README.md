# Marketplace Central — Architecture Rebaseline

> **Role:** sole current-program status/router after `AGENTS.md`  
> **Current phase:** **D3 — COMMUNICATION / EVENTS — OPEN / IN PROGRESS; D3-B1 independent review is the exact next action**  
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

### D3 — OPEN

`D3-COMMUNICATION-EVENTS.md` is the active D3 stage artifact. It currently records only the D3 purpose/boundary, imported D0–D2 invariants, review protocol and current stage status.

**No D3 batch is canonical yet.** The operator approved the direction of D3-B1 for independent challenge; `D3-B1-REVIEW-CANDIDATE.md` is deliberately non-authoritative until reviewer convergence and explicit operator batch acceptance.

## 4. Engineering method and repo lifecycle

Engineering reasoning follows the **DevelopmentConexus Engineering Method** identified in `AGENTS.md`; the local file in this authority path is only the consumed context copy.

This router defines the Marketplace Central D0–D9 status/lifecycle and allowed work. It is repo-specific specialization, not a second organizational engineering method. Conflicts inside the organizational method's scope must be surfaced, never silently reinterpreted.

## 5. What is prohibited now

While D3 remains open:

- do not begin D4–D9 target design prematurely;
- do not implement product architecture/features;
- do not let legacy IDs/tables/ADRs imply target identity, persistence ownership, communication architecture or target contracts;
- do not silently alter accepted D0/D1/D2 authority while choosing D3 communication semantics;
- do not create a semantic dependency outside D1 by hiding it in an event, API, queue, projection or database;
- do not choose provider/ERP transport contracts or credentials;
- do not choose HTTP/frontend/runtime topology;
- do not choose workers, brokers, outbox/transaction implementation, locks, RLS enforcement or deployment topology;
- do not treat `AI-DIALOG.md`, `D3-B1-REVIEW-CANDIDATE.md` or reviewer/chat summaries as target authority.

Existing code/module/context/schema shape remains current-state evidence only.

## 6. Exact next action

**Independent review of D3-B1 — Communication Topology & Edge Matrix.**

The operator has approved B1 as a review candidate, not yet as canonical architecture.

Next sequence:

1. Fable is invoked independently by the operator and reads the current repository authority path itself.
2. Fable reviews `docs/engineering/rebaseline/D3-B1-REVIEW-CANDIDATE.md` and appends only material findings to `AI-DIALOG.md`.
3. GPT independently adjudicates every material finding against repository authority/evidence; reviewer findings remain evidence, not authority.
4. Material reviewer disagreement receives another GitHub round or operator adjudication; GPT does not simulate Fable.
5. The converged B1 batch returns to the operator for explicit acceptance.
6. Only after acceptance is B1 consolidated into `D3-COMMUNICATION-EVENTS.md`; then D3-B2 may proceed.

If D3 discovers a genuinely necessary semantic dependency not allowed by D1, reopen only the implicated D1 decision rather than hiding the dependency in an event, queue, API, projection or database.

Do not advance D4 or product implementation before D3 is accepted.

## 7. Fresh-session success test

A fresh session should conclude that:

- D0, D1 and D2 are **CLOSED / ACCEPTED**;
- D1 defines exactly 12 semantic business boundaries but no runtime topology;
- D2 fixes canonical/external identities, tenant/isolation semantics, persistent ownership and shared value/knowledge/time semantics;
- D3 is **OPEN / IN PROGRESS** and owns communication/event/projection semantics only inside D1/D2 authority;
- D3-B1 is currently an **independent-review candidate**, not canonical target authority;
- current modules/contexts/schema/legacy ADRs remain evidence, not target authority by inheritance;
- implementation remains blocked until D9;
- D4 is not yet open.

If it cannot, the authority path is incomplete or contradictory.