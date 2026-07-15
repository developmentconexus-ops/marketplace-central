# F-03-estoque-workspace

```yaml
id: F-03
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-04
created: 2026-07-14
updated: 2026-07-14
validation_level: QA-0
lifecycle_scope: feature
```

## Mission

MIS-003. Binding contracts: IC-05 (route `/estoque`, staleTime class `stock` 45s / invalidation namespace `inventory`, states), IC-03 (`stock_correct` intent), R-01 screen 1h, R-02 (current stock-seguro page facts). ADR-17.

## Milestone

M-04 read-workspaces. Depends on M-03 (stock_correct envelope) + F-02 pattern precedent.

## Brief

Rebuild Estoque (wireframe 1h) at `/estoque`, replacing `/inventory/stock-seguro`: table with físico/reservado/disponível/seguro per product (existing inventory endpoints), divergence detection column (ConflictTag "divergente: ERP=n" when ERP vs published mismatch per 1h), stock-seguro rule editing (existing endpoints, parity with legacy page), bulk "Corrigir estoque" via M-03 modal (`stock_correct` intents), freshness (inventory 45s + refresh). ⚑ "estoque: ERP" provenance banner.

EARS:
- While ERP disponível differs from provider published stock, when the row renders, ConflictTag shall show "divergente: ERP={n}" with both values in detail.
- While the operator selects rows and confirms "Corrigir estoque", when the envelope applies, corrections shall target ERP-disponível-derived values (source_timestamp gate carries the inventory fetch time) and `invalidateAfterMutation('stock_correct')` shall fire on terminal.
- While any quantity is unknown (null), when rendered, the cell shall show UnknownValue — never 0 (unknown stock is not zero stock: overselling guard).
- While legacy `/inventory/stock-seguro?…` is visited, redirect shall land `/estoque?…`.

## Inputs

- R-01 §1h inventory, R-02 stock-seguro page facts (parity checklist), existing inventory endpoints + stock-seguro rule endpoints (sdk-runtime), IC-03 stock_correct schema, M-03 modal, IC-05 keys/components.

## Expected Output

- `/estoque` page; legacy page deleted; stock-seguro rule editing preserved (parity).
- Divergence computed from fields existing endpoints already return (spec.md names exact fields from R-02; if publish-side stock unavailable pre-M-01-refresh, column renders UnknownValue with hint — no fake comparison).
- Component tests: divergence render, unknown-not-zero, modal handoff (intent carries source_timestamp), redirect.

## Constraints

- Writes only via M-03 modal/envelope; legacy StockActionService UI paths replaced (fold landed in M-03 F-02).
- No new server endpoints; reads from existing inventory + listings modules.
- pt-BR glossary: físico/reservado/disponível/seguro exact.

## Negative Scenarios

- Inventory API error → page ErrorState with retry.
- Correction on listing with unresolved link → item fails `link_unresolved` in protocolo; UI shows failureCopy string (proof reuses M-03 surface).
- Zero selected rows → "Corrigir estoque" disabled.

## State Model

Queries: inventory list (namespace `inventory`, staleTime class `stock` 45s), `listListingsByProduct` (IC-02) for per-row published-stock join (listings ns — `getListingsSummary` carries only counters, not per-row quantities). Selection component-local, cleared on installation switch. Modal owns write lifecycle (M-03 pattern).

## Validation Expectations

- Vitest output: divergence, unknown-not-zero, modal-handoff tests green.
- Browser proof: 1h screenshot with ConflictTag rows; redirect proof; stock-seguro rule edit parity proof.
- Protocolo proof: stub-adapter stock_correct flow from this page reaching terminal with invalidation.

## Execution Artifact Rules

`spec.md`, `plan.md`, `validation.md` created during feature execution.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer (after F-02 accepted).
- Next action: compile context pack; read R-01 §1h + R-02 stock section + IC-03/IC-05.
- Required files/evidence: `validation.md` in this folder.
- Blockers or open decisions: exact divergence source fields fixed in spec.md from R-02 (bounded).
