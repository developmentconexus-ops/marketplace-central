# F-01-sankhya-read-contract-import

```yaml
id: F-01
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-03
created: 2026-07-06
updated: 2026-07-06
validation_level: QA-0
lifecycle_scope: feature
```

## Mission

MIS-001 Mercado Livre Operating Cockpit.

## Milestone

M-03 MNOS/Sankhya read contract.

## Brief

Translate IC-002 and MNOS source evidence into MPC-owned docs/types/ports for product, stock, price, cost, sales, and tax inputs.

## Inputs

- IC-002.
- MNOS files: `VW_ESTOQUE_SALDO`, `VW_PRECO_TABELA`, `TGFCUS`, `VW_FAT_VENDA_ITEM`, `VW_IMPOSTO_ITEM`, `TGFPRO`.

## Expected Output

- Contract doc and typed interfaces define required fields and quality flags.
- Feature specs cite exact MNOS files and rules.

## Constraints

- Do not copy MNOS implementation wholesale.
- Do not add direct ad hoc SQL inside business modules.

## Validation Expectations

- Contract includes stock formula, company/location defaults, `CUSSEMICM`, missing-value quality states, and no-write rule.

## Execution Artifact Rules

`spec.md`, `plan.md`, and `validation.md` are created during feature execution, not mission planning.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer.
- Next action: Create `spec.md`.
- Required files/evidence: F-01/validation.md.
- Blockers or open decisions: None.
