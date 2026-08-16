# Marketplace Central — Architecture Rebaseline

> **Role:** sole current-program status/router after `AGENTS.md`  
> **Current phase:** **D2 — IDENTITY / TENANT / DATA OWNERSHIP — GLOBAL COHERENCE / CLOSURE REVIEW**  
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
D2 — Identity / Tenant / Data Ownership — GLOBAL COHERENCE / CLOSURE REVIEW
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

### D2 — B1+B2 CONSOLIDATED; CLOSURE REVIEW IN PROGRESS

`D2-IDENTITY-TENANT-DATA-OWNERSHIP.md` is the current D2 authority. B1+B2 are operator-approved, independently challenged and consolidated; no B3 is required.

Current locked D2 direction includes:

- canonical identity follows semantic authority; MPC-owned IDs are opaque/stable/non-reusable;
- MPC-owned Organization, Marketplace Installation, Selling Entity, Inventory Source, Fulfillment Node and human/automation/system Principal semantics;
- Organization as canonical tenant/isolation root with no duplicate target `Tenant` identity;
- source-qualified external Product, Listing/Variation, Sale/Order, Shipment and native financial identities;
- minimal SourceInstance identity while D4 owns concrete source contracts/capabilities/config/credentials;
- MPC-owned Business/Invoicing Intents, Post-Sale Resolution, Economic Attribution/Reconciliation, Operational Work and Authorization Decision/Grant semantics;
- external OIDC human AuthN with MPC-owned Principal/binding/Membership/Role/Permission ordinary-access state and a strict fence against business-action authority leaking into identity/access;
- one D1 semantic write authority per canonical business meaning; typed references/snapshots/projections do not become current write authority;
- exact Money/rate/material-quantity semantics, bounded `Fact<T>` scope, provenance and distinct material time meanings;
- clean target persistence baseline with **no legacy-data migration or archival requirement** for the pre-rebaseline MPC database;
- old ADR structures do not carry forward by inheritance; D2 has adjudicated ADR-011/012/022/028/031 and defined safe rehoming gates/new ADR-036+ transition.

The remaining work is only the final D2 Global Coherence + YAGNI / Overengineering / Future-Cost review and any bounded corrections it finds.

## 4. Engineering method and repo lifecycle

Engineering reasoning follows the **DevelopmentConexus Engineering Method** identified in `AGENTS.md`; the local file in this authority path is only the consumed context copy.

This router defines the Marketplace Central D0–D9 status/lifecycle and allowed work. It is repo-specific specialization, not a second organizational engineering method. Conflicts inside the organizational method's scope must be surfaced, never silently reinterpreted.

## 5. What is prohibited now

Until D2 is accepted as a whole:

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

**Complete the final D2 Global Coherence + YAGNI / Overengineering / Future-Cost review.**

The review must check the complete consolidated D2 set for duplicate/missing authority, contradictory assumptions, circular ownership, hidden generic frameworks, lost seams, unnecessary future-proofing and later-stage mechanism leakage.

If bounded corrections resolve all findings without reopening D0/D1 or creating a new material D2 batch, present D2 as a closure candidate for explicit operator ratification.

Do not advance D3 or product implementation before D2 is accepted as a whole and this router advances.

## 7. Fresh-session success test

A fresh session should conclude that:

- D0 and D1 are **CLOSED / ACCEPTED**;
- D1 defines exactly 12 semantic business boundaries but no runtime topology;
- D2 B1+B2 are consolidated and **GLOBAL COHERENCE / CLOSURE REVIEW** is the current activity;
- no B3 is planned unless Global Coherence finds a genuinely new material D2 problem;
- current modules/contexts/schema/legacy ADRs remain evidence, not target authority by inheritance;
- implementation remains blocked until D9;
- the exact next action is final D2 Global Coherence + YAGNI / Overengineering / Future-Cost review;
- D3–D9 design is not yet authorized.

If it cannot, the authority path is incomplete or contradictory.