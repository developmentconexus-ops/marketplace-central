# D8 — Golden Flows

> **Status:** OPEN / ACTIVE — NOT YET DERIVED  
> **Program:** Architecture Rebaseline / Technical System Design  
> **Opened:** 2026-08-22  
> **Accepted prerequisites:** D0–D7 — ACCEPTED / CLOSED  
> **Method:** DevelopmentConexus Engineering Method v1.0.0

## Purpose

D8 defines the smallest set of representative **golden flows** needed to prove that the accepted D0–D7 Product, ownership, API, frontend and runtime authorities compose coherently end to end before the adversarial D9 review.

D8 owns flow selection, choreography and proof expectations only. It does **not** reopen accepted Product operations, Permissions, Principal kinds, semantic owners, provider contracts, frontend authority or D7 runtime mechanisms by convenience, and it does not begin Product implementation.

## Accepted baseline

```text
D0–D7                  ACCEPTED / CLOSED
Product operations     99
ordinary Permissions   30
Principal kinds        H / A / S only
stable origin          https://conexus.fun
active runtime         NONE
D9                      BLOCKED
Product implementation BLOCKED UNTIL D9
```

## Target

> Derive a finite, high-leverage set of golden flows whose combined success, failure, ambiguity, authorization, Organization-isolation and recovery cases are sufficient to expose contradictions between accepted stage authorities without turning D8 into an exhaustive operation-by-operation test catalog.

## Boundary

- Prefer representative flows that cross multiple accepted boundaries and can falsify architecture composition.
- Reuse the canonical Product OAD and accepted D6/D7 authorities rather than creating parallel route, DTO, workflow or runtime truth.
- Include negative/ambiguous/recovery cases only where they test a material accepted invariant.
- Do not select implementation details that D0–D7 deliberately left manifest/configuration-level.
- Do not begin D9 or Product implementation.

## Exact next action

**Derive and adjudicate the smallest D8 golden-flow set and its falsifiable acceptance matrix from current repository authority. Start from this owner and switch only to the exact prior owner/OAD needed for each candidate flow.**
