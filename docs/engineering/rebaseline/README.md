# Marketplace Central — Architecture Rebaseline

> **Role:** sole current-program status/router after `AGENTS.md`  
> **Current phase:** **D1 — DOMAINS / BOUNDARIES — CLOSED / ACCEPTED; D2 — IDENTITY / TENANT / DATA OWNERSHIP is the exact next stage**  
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
8. active D2 artifact once D2 is opened
9. `docs/engineering/rebaseline/EVIDENCE-REGISTER.md`
10. code, OpenAPI, schemas, tests and runtime only as current-state evidence when needed

This file alone answers **where the program is and what happens next**. Stable architecture belongs in `ARCHITECTURE.md`; accepted stage semantics belong in D-stage artifacts; Git history is the archive.

Do not reconstruct target authority from memory, legacy package shape, historical plans or stale docs.

## 2. Program state

```text
Documentary / governance cleanup — DONE
  ↓
D0 — Product / System Definition — CLOSED / ACCEPTED
  ↓
D1 — Domains / Boundaries — CLOSED / ACCEPTED
  ↓
D2 — Identity / Tenant / Data Ownership — NEXT, NOT YET OPENED
  ↓
D3 — Communication / Events
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

## 3. Accepted baseline

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

## 4. Engineering method and repo lifecycle

Engineering reasoning follows the **DevelopmentConexus Engineering Method** identified in `AGENTS.md`; the local file in this authority path is only the consumed context copy.

This router defines the Marketplace Central D0–D9 status/lifecycle and allowed work. That is **repo-specific specialization, not a second organizational engineering method**. If the lifecycle ever conflicts with the organizational method inside the method's scope, surface the conflict instead of reinterpreting either silently.

## 5. What is prohibited now

Until D2 is adjudicated:

- do not begin D3–D9 target design prematurely;
- do not implement product architecture/features;
- do not let legacy IDs/tables imply canonical identity or persistence ownership;
- do not choose events/outbox/sync-vs-async communication;
- do not choose provider/ERP transport contracts;
- do not choose HTTP/frontend/runtime topology;
- do not silently contradict accepted D0/D1 authority.

Existing code/module/context/schema shape remains current-state evidence only.

## 6. Exact next action

**Open D2 — Identity / Tenant / Data Ownership with the operator.**

D2 must determine canonical identities, tenant/isolation semantics, persistence/data ownership and exact shared/domain value/evidence/time representations from accepted D0/D1 authority. It must not infer target identity from legacy `CODPROD`, current tables or provider/ERP nouns merely because they exist.

Do not advance D3 or product implementation before D2 is accepted.

## 7. Fresh-session success test

A fresh session should conclude that:

- D0 and D1 are **CLOSED / ACCEPTED**;
- D1 defines exactly 12 semantic business boundaries but no runtime topology;
- current modules/contexts/schema remain evidence, not target authority;
- implementation remains blocked until D9;
- the exact next stage is **D2 — Identity / Tenant / Data Ownership**;
- D3–D9 design is not authorized yet.

If it cannot, the authority path is incomplete or contradictory.
