# D6-R2 P5-R1 — B110 Approvals AuthorizationRequest Targeted Supersession

> **Status:** DERIVED / CURRENT FOR B110
> **Parent inventory:** [D6-R2 P5 Complete Screen / Material-Surface Inventory](D6-R2-P5-SCREEN-SURFACE-INVENTORY.md)
> **Identity authority:** [D2-R6 AuthorizationRequest Ratification](D6-R2-NOTIF-01-D2-R6-RATIFICATION.md)
> **Wire authority:** [D5-R6 AuthorizationRequest OAD Proof](D6-R2-NOTIF-01-D5-R6-AUTHORIZATION-REQUEST-OAD-WIRE-PROOF.md)
> **Boundary:** supersedes only the old R110/R111 Approvals disposition; no other P5 route meaning is reopened.

## 1. Why this bounded P5 repair exists

The old P5 model treated `/aprovacoes` primarily as AuthorizationDecision queue/history and its detail as target/revision decision context. D2-R6 + D5-R6 introduced a distinct canonical pre-decision `AuthorizationRequest` plus purpose-bounded actionable reads. Keeping the old R110/R111 wording would make frontend planning contradict the current 106/31 Product OAD.

This is a material user-flow correction, not preference-driven IA redesign.

## 2. Current Approvals surface inventory

### R110 — `/aprovacoes`

One existing global destination under `CONTROLE > Aprovações`, with two **local permission-independent lenses**:

```text
Para decidir
→ ListMyActionableAuthorizationRequests
→ governance.decide
→ exact current human actionability only

Histórico
→ ListAuthorizationDecisions
→ governance.read
→ immutable Decision history
```

`governance.decide` does not imply `governance.read`, and `governance.read` does not imply `governance.decide`.

### R111-A — `/aprovacoes/solicitacoes/:authorizationRequestId`

Material actionable-request detail:

```text
GetMyActionableAuthorizationRequest
→ one typed review basis
→ inline decision confirmation
→ CreateAuthorizationDecision
```

F13 `AUTHORIZATION_ACTION_REQUIRED` may deep-link here through `AuthorizationRequestRef`, but Notification possession remains awareness only and never substitutes for current server authorization.

### R111-B — `/aprovacoes/decisoes/:authorizationDecisionId`

Immutable historical Decision detail:

```text
GetAuthorizationDecision
→ outcome + deciding Principal + decided_at
→ governed target ref
→ immutable typed review basis
→ no decision controls
```

## 3. Structural consequence

```text
old Approvals content homes      2
current Approvals content homes  3
new top-level IA destinations    0
new Product operations from P5   0
screen-shaped API                0
```

The extra detail route reflects two genuinely different identities/lifecycles — pending `AuthorizationRequest` versus terminal `AuthorizationDecision` — rather than cosmetic page splitting.

## 4. Explicit non-goals

This repair does not admit:

- a second global `Histórico de aprovações` navigation destination;
- mixed pending/history authority in one undifferentiated list;
- approval search, totals, generic lifecycle filters, approver filters or bulk decisions;
- source-domain read Permission inference from `governance.decide`;
- modal-only consequential decision flow;
- frontend ownership of authorization validity.

## 5. Feed-forward

This targeted P5 supersession is the structural input to [B110 P8 candidate](D6-R2-P8-B110-APPROVALS-CANDIDATE.md). Final P9 remains blocked until the operator explicitly `LOCK`s the rendered B110 structure.
