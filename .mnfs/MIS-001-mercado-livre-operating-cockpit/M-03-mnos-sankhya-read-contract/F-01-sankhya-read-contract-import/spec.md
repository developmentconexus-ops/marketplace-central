# Feature Spec

```yaml
id: F-01
type: feature-spec
status: spec_ready
owner: Feature Implementer
parent: F-01
created: 2026-07-07
updated: 2026-07-07
validation_level: QA-0
lifecycle_scope: feature
```

## Feature ID

F-01-sankhya-read-contract-import

## MNOS Evidence

- `C:\Users\leandro.theodoro\Documents\MNOS\semantic\views\vw_estoque_saldo.sql`
- `C:\Users\leandro.theodoro\Documents\MNOS\semantic\views\vw_preco_tabela.sql`
- `C:\Users\leandro.theodoro\Documents\MNOS\semantic\views\vw_fat_venda_item.sql`
- `C:\Users\leandro.theodoro\Documents\MNOS\semantic\views\vw_imposto_item.sql`
- `C:\Users\leandro.theodoro\Documents\MNOS\semantic\governance\tgfcus.yml`
- `C:\Users\leandro.theodoro\Documents\MNOS\semantic\governance\tgfpro.yml`

## Problem

Marketplace Central needs an MPC-owned read contract for internal product, stock, price, cost, sales, and tax facts before `product_links`, `inventory`, and `profitability` can depend on Sankhya/MNOS semantics safely. Without this seam, future features will drift into ad hoc SQL, silent zero defaults, or Sankhya-coupled write paths.

## Requirements

- Requirement: Import IC-002 into MPC-owned domain and port types for product linking, sellable stock, current price, cost-as-of, sales history, and tax inputs.
  - Acceptance evidence: `apps/server_core/internal/modules/internal_read/domain/*.go` and `apps/server_core/internal/modules/internal_read/ports/reader.go` compile under focused Go tests.
- Requirement: Preserve the default stock semantics exactly as `SUM(ESTOQUE - RESERVADO)` with `CODEMP IN (1,2)` and `CODLOCAL=10101`.
  - Acceptance evidence: contract tests assert the default `StockScope` values and scope code.
- Requirement: Keep `CODLOCAL=10108` excluded from the default sellable stock scope.
  - Acceptance evidence: only location `10101` appears in the default stock scope contract.
- Requirement: Model missing values as explicit quality flags, never as zero defaults.
  - Acceptance evidence: contract tests assert the required quality flags and pointer-based nullable cost/price/tax fields.
- Requirement: Name `CUSSEMICM` as the initial cost basis for margin inputs.
  - Acceptance evidence: `domain.CostAsOf` exposes `CUSSEMICM *float64`.
- Requirement: Keep the seam read-only and do not introduce Sankhya write paths or ERP mirroring.
  - Acceptance evidence: changed paths are limited to feature artifacts, domain contract files, and the read port.

## Non-Goals

- Implement the real Oracle/Sankhya adapter.
- Add business logic for link resolution, stock reconciliation, or profitability math.
- Introduce HTTP routes, persistence, or MPC snapshot tables.
- Mirror MNOS/Sankhya tables inside MPC.

## Design

This feature adds a new `internal_read` module surface with pure domain structs and a single application-facing `Reader` port. The domain layer defines the quality flags, stock scope, and typed read models for product candidates, sellable stock, current price, cost-as-of, tax inputs, and sales history. Nullable source values use pointers so missing data remains explicit and can later feed quality rules instead of becoming zero.

The stock contract encodes the mission defaults directly in the domain surface: companies `1` and `2`, location `10101`, formula `SUM(ESTOQUE - RESERVADO)`, and a `revenda` scope code. Showroom stock stays outside the default simply by not appearing in the location list. The port layer exposes the six IC-002 operations with strongly typed inputs and outputs, giving future adapters and modules one seam to depend on without reaching for module-local SQL.

## Edge Cases

- Product linking may find no candidate: return an empty candidate set or candidates flagged `missing_product`, not synthetic placeholder products.
- Product linking may find more than one exact candidate: consumers must receive `ambiguous_product`, not an arbitrary winner.
- Stock may be negative when reservations exceed stock: the contract must allow negative sellable quantities.
- Cost, price, or tax values may be missing: fields stay `nil` and quality flags carry the missing state.
- Source freshness matters for future consumers: stock/price/tax/sales models keep `SourceFetchedAt`.

## Acceptance Criteria

- Criterion: The internal read contract preserves the default sellable stock rule with mission company and location scope.
  - Traces to milestone criterion ID: M-03-C01
  - Proven by (verification command or QA step): `cd apps/server_core; $env:GOCACHE=(Resolve-Path .gocache); go test ./internal/modules/internal_read/...`
- Criterion: The contract exposes explicit missing-value quality flags and `CUSSEMICM` as the nullable cost basis without zero-filling.
  - Traces to milestone criterion ID: M-03-C02
  - Proven by (verification command or QA step): `cd apps/server_core; $env:GOCACHE=(Resolve-Path .gocache); go test ./internal/modules/internal_read/...`

## Handoff

- Current status: `spec_ready`
- Next owner: Feature Implementer
- Next action: Write `plan.md` and implement the scoped feature.
- Required files/evidence: feature brief, spec, milestone contract, validation expectations
- Blockers or open decisions: None.
