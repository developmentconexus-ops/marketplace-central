# Marketplace Central — Architecture Rebaseline

> **Role:** sole current-program status/router after `AGENTS.md`  
> **Current phase:** **D2 — IDENTITY / TENANT / DATA OWNERSHIP — OPEN / IN PROGRESS**  
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

Do not reconstruct target authority from memory, legacy package shape, historical plans, `AI-DIALOG.md` or stale docs.

## 2. Program state

```text
Documentary / governance cleanup — DONE
  ↓
D0 — Product / System Definition — CLOSED / ACCEPTED
  ↓
D1 — Domains / Boundaries — CLOSED / ACCEPTED
  ↓
D2 — Identity / Tenant / Data Ownership — OPEN / IN PROGRESS
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

D1 closure included:

- iterative operator adjudication of D1.1–D1.18;
- independent adversarial Fable review;
- adjudication of all six Fable amendments without adding/removing/merging/splitting a boundary;
- final Global Coherence + YAGNI / Overengineering / Future-Cost review;
- the D1.G1 correction assigning provider-requirement closure for claimed fulfillment paths to Fulfillment Lifecycle while provider-native truth/protocol remains external/D4-owned;
- explicit operator approval of the corrected whole.

The 12 D1 business boundaries do **not** imply 12 services, databases, processes or deployments.

### D2 — OPEN / IN PROGRESS

`D2-IDENTITY-TENANT-DATA-OWNERSHIP.md` is the current D2 authority. Operator-approved decisions already recorded there are binding within D2; the stage is not yet accepted as a whole.

Current locked D2 direction includes:

- canonical identity follows semantic authority;
- MPC-owned Organization, Marketplace Installation, Selling Entity, Inventory Source, Fulfillment Node and Principal identities;
- Organization as the canonical tenant/isolation root with no duplicate target `Tenant` identity;
- source-qualified external Product, Listing/Variation, Sale/Order, Shipment and native financial identities;
- minimal Source Instance identity in D2 while D4 owns source contracts/capabilities/config/credentials;
- MPC-owned Business Order Intent, Invoicing Intent, Post-Sale Resolution, Economic Attribution/Reconciliation semantics, Operational Work and Authorization Decision;
- external OIDC AuthN with MPC-owned Principal/Membership/ordinary Role/Permission state and a strict fence against business-action authority leaking into the identity/access substrate;
- one D1 semantic write authority per canonical business meaning; projections/snapshots/references do not become current write authority;
- exact Money/rate semantics, bounded `Fact<T>` scope, provenance and distinct material time meanings;
- clean target persistence baseline with **no legacy-data migration or archival requirement** for the pre-rebaseline MPC database, per operator ruling.

D2 remains open only for the minimum remaining representation/persistence closure and final stage coherence work named by the D2 artifact.

## 4. Engineering method and repo lifecycle

Engineering reasoning follows the **DevelopmentConexus Engineering Method** identified in `AGENTS.md`; the local file in this authority path is only the consumed context copy.

This router defines the Marketplace Central D0–D9 status/lifecycle and allowed work. That is **repo-specific specialization, not a second organizational engineering method**. If the lifecycle ever conflicts with the organizational method inside the method's scope, surface the conflict instead of reinterpreting either silently.

## 5. What is prohibited now

Until D2 is closed as a whole:

- do not begin D3–D9 target design prematurely;
- do not implement product architecture/features;
- do not let legacy IDs/tables imply canonical identity or persistence ownership;
- do not choose events/outbox/sync-vs-async communication;
- do not choose provider/ERP transport contracts or credentials;
- do not choose HTTP/frontend/runtime topology;
- do not silently contradict accepted D0/D1 or locked D2 authority;
- do not treat `AI-DIALOG.md` or reviewer/chat summaries as target authority.

Existing code/module/context/schema shape remains current-state evidence only.

## 6. Exact next action

**Continue D2 with Batch B2 — representation/persistence closure and Global Coherence preparation.**

B2 must identify only the remaining D2 decisions necessary to make identity/value/persistent-state semantics implementation-ready, explicitly defer D3–D7 mechanisms, and prepare the final D2 Global Coherence + YAGNI / Overengineering / Future-Cost review.

Do not advance D3 or product implementation before D2 is accepted as a whole and this router advances.

## 7. Fresh-session success test

A fresh session should conclude that:

- D0 and D1 are **CLOSED / ACCEPTED**;
- D1 defines exactly 12 semantic business boundaries but no runtime topology;
- D2 is **OPEN / IN PROGRESS** and `D2-IDENTITY-TENANT-DATA-OWNERSHIP.md` is its current authority;
- locked D2 decisions include the clean-baseline/no-legacy-data-migration operator ruling;
- current modules/contexts/schema remain evidence, not target authority;
- implementation remains blocked until D9;
- the exact next action is **D2 Batch B2 — representation/persistence closure and Global Coherence preparation**;
- D3–D9 design is not yet authorized.

If it cannot, the authority path is incomplete or contradictory.
