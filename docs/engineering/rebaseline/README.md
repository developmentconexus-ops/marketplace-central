# Marketplace Central — Architecture Rebaseline

> **Role:** sole current-program status/router after `AGENTS.md`  
> **Current phase:** **D2 — IDENTITY / TENANT / DATA OWNERSHIP — CLOSURE CANDIDATE / AWAITING OPERATOR RATIFICATION**  
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
9. `docs/engineering/rebaseline/EVIDENCE-REGISTER.md`
10. code, OpenAPI, schemas, tests and runtime only as current-state evidence when needed

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
D2 — Identity / Tenant / Data Ownership — CLOSURE CANDIDATE / AWAITING OPERATOR RATIFICATION
  ↓
D3 — Communication / Events — BLOCKED BY D2
  ↓
D4 — External Integrations
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

## 3. Accepted/current baseline

### D0 — CLOSED

Product 1.0 is **Marketplace Operations + Commercial Intelligence**. MPC is the marketplace operations control plane: external systems retain authority for facts/processes inherently theirs while MPC owns the cross-system marketplace operating semantics needed to observe, understand, decide, execute, verify and reconcile.

D0 authority and non-goals are defined only in `D0-PRODUCT-SYSTEM-DEFINITION.md`.

### D1 — CLOSED

`D1-DOMAINS-BOUNDARIES.md` is the accepted D1 authority. It defines 12 business boundaries, explicit ownership/non-ownership, cross-cutting non-domain treatment, legal semantic authority edges, forbidden boundary violations, legacy semantic disposition, D2–D7 defers and reopen triggers.

D1 closure included operator adjudication, independent Fable review and final Global Coherence + YAGNI / Overengineering / Future-Cost review. The 12 D1 business boundaries do **not** imply 12 services, databases, processes or deployments.

### D2 — CLOSURE CANDIDATE

`D2-IDENTITY-TENANT-DATA-OWNERSHIP.md` is the current D2 authority. B1+B2 are operator-approved and independently challenged. Final Global Coherence + YAGNI / Overengineering / Future-Cost review completed with only bounded consolidation corrections and no D0/D1 reopen, no B3 and no material contradiction.

Current locked D2 direction includes:

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
- old ADR structures do not carry forward by inheritance; D2 has adjudicated ADR-011/012/022/028/031 and defines safe rehoming gates/new ADR-036+ transition.

## 4. Engineering method and repo lifecycle

Engineering reasoning follows the **DevelopmentConexus Engineering Method** identified in `AGENTS.md`; the local file in this authority path is only the consumed context copy.

This router defines the Marketplace Central D0–D9 status/lifecycle and allowed work. It is repo-specific specialization, not a second organizational engineering method. Conflicts inside the organizational method's scope must be surfaced, never silently reinterpreted.

## 5. What is prohibited now

Until D2 is explicitly ratified as a whole:

- do not begin D3–D9 target design prematurely;
- do not implement product architecture/features;
- do not let legacy IDs/tables/ADRs imply target identity, persistence ownership or architecture;
- do not choose events/outbox/sync-vs-async communication;
- do not choose provider/ERP transport contracts or credentials;
- do not choose HTTP/frontend/runtime topology;
- do not silently contradict accepted D0/D1 or locked D2 authority;
- do not treat `AI-DIALOG.md` or reviewer/chat summaries as target authority.

Existing code/module/context/schema shape remains current-state evidence only.

## 6. Exact next action

**Operator ratifies or amends D2 as a whole.**

The complete D2 closure candidate is `docs/engineering/rebaseline/D2-IDENTITY-TENANT-DATA-OWNERSHIP.md`, including its final Global Coherence + YAGNI / Overengineering / Future-Cost review.

If the operator approves the corrected whole:

1. mark D2 `CLOSED / ACCEPTED`;
2. advance this router so **D3 — Communication / Events** is exact next stage;
3. do not begin implementation; implementation remains blocked until D9.

If the operator finds a material issue, reopen only the implicated D2 decision rather than re-running the whole stage.

## 7. Fresh-session success test

A fresh session should conclude that:

- D0 and D1 are **CLOSED / ACCEPTED**;
- D1 defines exactly 12 semantic business boundaries but no runtime topology;
- D2 is a **CLOSURE CANDIDATE / AWAITING OPERATOR RATIFICATION**;
- B1+B2 are consolidated and independently challenged; no B3 is planned;
- final D2 Global Coherence completed without a material contradiction;
- current modules/contexts/schema/legacy ADRs remain evidence, not target authority by inheritance;
- implementation remains blocked until D9;
- the exact next action is explicit operator ratification of D2 as a whole;
- D3 remains blocked until that ratification.

If it cannot, the authority path is incomplete or contradictory.