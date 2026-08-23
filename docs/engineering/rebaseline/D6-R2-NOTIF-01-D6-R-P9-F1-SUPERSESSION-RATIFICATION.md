# NOTIF-01 D6-R P9-F1 — Query-Only Supersession & Targeted Architecture Reopen Ratification

> **Status:** ACCEPTED / OPERATOR-RATIFIED — QUERY-ONLY REMEDY SUPERSEDED
> **Operator ruling:** 2026-08-23
> **Supersedes only the remedy** proposed by [P9-F1 Actionable Governance Context Finding](D6-R2-NOTIF-01-D6-R-P9-F1-ACTIONABLE-GOVERNANCE-CONTEXT.md); the underlying P9 falsifier remains valid.
> **D1 revalidation:** [D1-R2 Authorization-Request Boundary Revalidation](D6-R2-NOTIF-01-D1-R2-AUTHORIZATION-REQUEST-BOUNDARY-REVALIDATION.md)
> **D2 targeted reopen:** [D2-R6 Authorization Request Identity & Lifecycle](D6-R2-NOTIF-01-D2-R6-AUTHORIZATION-REQUEST-IDENTITY-LIFECYCLE.md)
> **Canonical Product wire:** unchanged — 104 Product operations · 31 ordinary Permissions · H/A/S
> **Implementation:** BLOCKED UNTIL accepted D9

## 1. Operator ruling

The operator rejected freezing the P9-F1 query-only remedy as the target merely because it was the smallest correction inside the current wire structure.

The approved direction is:

```text
P9 falsifier
→ re-evaluate root cause under Engineering Method / Global Maximum
→ D1 boundary revalidation
→ targeted D2 identity/lineage reopen
→ only then D3/D5/D6 feed-forward
```

No OpenAPI change or operation-count prediction is authorized by this ratification.

## 2. What remains true from P9-F1

The factual defect remains:

```text
AUTHORIZATION_ACTION_REQUIRED
→ current human is told that authorization is required
→ current wire has no pre-decision canonical subject the human can open
→ CreateAuthorizationDecision expects exact current authorization target/revision
→ governance.decide does not imply governance.read or target-owner read
```

The original P9-F1 correctly proved that the frontend cannot truthfully close the actionable approval flow with the current 104/31 wire.

What is superseded is only this proposed remedy:

```text
add one purpose-bounded query over target + current revision
→ assume target + revision are sufficient identity for an authorization episode
```

## 3. New material evidence

D0 already requires authorization to remain valid only while materially governing context remains valid enough for the intended action. Authority, delegation, policy, evidence sufficiency/readiness or mandatory safety conditions may change independently from the business target's own revision.

Therefore one target revision can participate in materially distinct authorization episodes:

```text
PriceIntent PI-100 / revision 7
→ authorization episode A
→ episode A becomes invalidated by material governing-context drift

PriceIntent PI-100 / revision 7
→ authorization episode B
→ new human authorization required
```

The two episodes have the same target identity and may have the same target revision, but they are not the same authorization occurrence.

This falsifies `target + target revision` as complete authorization-episode identity.

## 4. Method classification

The Engineering Method says:

- do not patch the symptom before identifying the structural condition that made it possible;
- Global Maximum may change current structure when the current structure preserves the root cause;
- if correctness requires a materially new business identity not supported by D2, reopen the implicated D2 identity decision rather than hiding the gap in API/UI mechanics.

Outcome:

**RESTRUCTURE NOW — targeted pre-decision Governance identity repair.**

This is not permission for a generic workflow engine, approval DSL or authorization platform.

## 5. Scope of the approved direction

The operator approved these semantic statements:

1. the pre-decision authorization case is distinct from the later `AuthorizationDecision` occurrence;
2. the pre-decision case requires stable MPC-owned identity under `Controlled Action Governance`;
3. the canonical concept is named **`AuthorizationRequest`** unless the D2 derivation proves a semantic contradiction;
4. multiple distinct `AuthorizationRequest`s may legitimately reference the same business target/revision across separate authorization episodes;
5. `AuthorizationDecision` remains a distinct terminal decision occurrence and does not become `pending` merely to absorb the request lifecycle;
6. the action-owning domain retains Business Intent, business disposition and execution-time validity;
7. Governance retains only authorization-specific request/context/authority/decision meaning;
8. Notifications remain awareness only; F13 should eventually reference the exact pre-decision authorization case rather than store current decision data;
9. the canonical OAD remains 104/31 until D2→D3→D5 derivation proves the exact new wire;
10. operation count is a consequence, never a target.

## 6. Explicit non-goals

This ratification does not admit:

```text
workflow/BPM engine
arbitrary approval stages
quorum/voting
approval DSL
generic Request/Case platform
custom governed entity graph
new ordinary Permission
runtime/storage topology
Product implementation
```

## 7. Gate

```text
P8 B00-R2 / B11 / B12        OPERATOR-LOCKED
P9-F1 factual falsifier       VALID
P9-F1 query-only remedy       SUPERSEDED
D1-R2 boundary revalidation   CURRENT STRUCTURE CONFIRMED
D2-R6 identity/lifecycle      OPEN / NEXT FOR EXACT ADJUDICATION
canonical Product OAD         104/31 UNCHANGED
D7-R / D8-R                   BLOCKED
Product implementation        BLOCKED UNTIL D9
```

**Exact next action:** adjudicate the exact D2-R6 `AuthorizationRequest` identity/lifecycle contract. Do not modify D3/D5/OAD or resume P9 until D2-R6 closes.