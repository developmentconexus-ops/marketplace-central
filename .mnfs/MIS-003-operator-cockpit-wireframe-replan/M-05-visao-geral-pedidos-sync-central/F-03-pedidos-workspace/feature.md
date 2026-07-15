# F-03-pedidos-workspace

```yaml
id: F-03
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-05
created: 2026-07-14
updated: 2026-07-14
validation_level: QA-0
lifecycle_scope: feature
```

## Mission

MIS-003. Binding contracts: IC-05 (route `/pedidos`, orders ns 120s), R-01 screen 1j, F-01 orders API. ADR-17. Non-Scope: faturar/NF emission.

## Milestone

M-05. Depends on F-01 + F-02 (route seam order).

## Brief

Build Pedidos read workspace (wireframe 1j) at `/pedidos`: orders table (pedido id, comprador, itens resumo, total, status tag, fulfillment tag blu, NF state, data), filters per 1j (status, fulfillment, período, busca), detail drawer from `getOrder` (itens, valores, timeline, "▸ técnico"), NF state display-only (emitida/pendente/null→"—"). Read-only: zero mutation controls.

EARS:
- While orders load, when a row renders, status/fulfillment tags shall use the fixed tag palette (fulfillment=blu) with pt-BR labels fixed in spec.md from R-01 §1j.
- While NF state is unknown, when rendered, the cell shall show "—" with hint "NF: ERP" (provenance) — never "pendente" by default.
- While filters/search change, when state updates, URL shall encode it; reload restores view.
- While `/orders?…` legacy URL is visited, redirect shall land `/pedidos?…`.

## Inputs

- R-01 §1j inventory, F-01 orders endpoints (sdk-runtime), IC-05 `ordersQueryKeys` + components.

## Expected Output

- `/pedidos` page, TanStack Query only; detail drawer; URL round-trip.
- Component tests: label map, unknown-NF render, round-trip, read-only audit (no mutation calls importable).

## Constraints

- Read-only; no faturar/emitir/cancel controls even disabled (Non-Scope, not "coming soon").
- Currency formatting via existing shared formatter; pt-BR.

## Negative Scenarios

- Unknown order id in drawer → 404 copy "Pedido não encontrado.".
- Malformed period filter from stale URL → 400 handled: ErrorState with clear-filter retry.
- Empty result → EmptyState "Nenhum registro encontrado.".

## Validation Expectations

- Vitest output: listed tests green.
- Browser proof: 1j screenshot; F5 fidelity; redirect proof.
- Grep proof: no mutation sdk methods imported by the page.

## Execution Artifact Rules

`spec.md`, `plan.md`, `validation.md` created during feature execution.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer (after F-02 accepted).
- Next action: compile context pack; read R-01 §1j + F-01 output + IC-05.
- Required files/evidence: `validation.md` in this folder.
- Blockers or open decisions: none.
