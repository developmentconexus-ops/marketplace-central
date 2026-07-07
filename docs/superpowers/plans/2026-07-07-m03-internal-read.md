# M-03 Internal Read Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build the M-03 MNOS/Sankhya internal read seam so MPC owns typed product/stock/price/cost/tax/sales contracts, fake/test adapters, quality-flag semantics, and milestone evidence proving stock and missing-input rules.

**Architecture:** Add a new shared module at `apps/server_core/internal/modules/internal_read/` with `domain`, `ports`, `application`, and `adapters`. Keep M-03 read-only and fake-adapter-driven: Oracle remains a guarded seam, while tests and future consumers rely on stable domain contracts and explicit quality flags.

**Tech Stack:** Go 1.25.x, `pgx/v5`, `pgxpool.Pool`, PowerShell, MNFS artifacts under `.mnfs/`, Git.

## Global Constraints

- Preserve MNOS semantics from `.mnfs/MIS-001-mercado-livre-operating-cockpit/research/mnos-sankhya-read-interface-contract.md`.
- Default sellable stock remains `SUM(ESTOQUE - RESERVADO)` with `CODEMP IN (1,2)` and `CODLOCAL=10101`.
- `CODLOCAL=10108` showroom stock must not contribute to default sellable stock.
- Initial margin cost basis is `CUSSEMICM`, never `CUSVARIAVEL`.
- Missing product, stock, cost, tax, or ambiguous data must become explicit quality flags, never zero defaults.
- Sankhya access remains read-only.
- MPC must not mirror whole ERP tables; only MPC-owned state or auditable snapshots are allowed.
- Secrets and DSN values must not appear in logs, errors, or validation artifacts.
- Use `GOCACHE=.gocache` for Go test commands.
- Preserve unrelated user changes in the worktree.

---

### Task 1: F-01 Contract Import And Domain Surface

**Files:**
- Create: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-03-mnos-sankhya-read-contract/F-01-sankhya-read-contract-import/spec.md`
- Create: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-03-mnos-sankhya-read-contract/F-01-sankhya-read-contract-import/plan.md`
- Create: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-03-mnos-sankhya-read-contract/F-01-sankhya-read-contract-import/validation.md`
- Create: `apps/server_core/internal/modules/internal_read/domain/quality_flag.go`
- Create: `apps/server_core/internal/modules/internal_read/domain/internal_product.go`
- Create: `apps/server_core/internal/modules/internal_read/domain/internal_stock.go`
- Create: `apps/server_core/internal/modules/internal_read/domain/internal_price.go`
- Create: `apps/server_core/internal/modules/internal_read/domain/internal_cost.go`
- Create: `apps/server_core/internal/modules/internal_read/domain/internal_tax.go`
- Create: `apps/server_core/internal/modules/internal_read/domain/internal_sales.go`
- Create: `apps/server_core/internal/modules/internal_read/ports/reader.go`

**Interfaces:**
- Consumes: `.mnfs/MIS-001-mercado-livre-operating-cockpit/research/mnos-sankhya-read-interface-contract.md`
- Produces:
  - `type QualityFlag string`
  - `type StockScope struct { Companies []int; Locations []int; Formula string; ScopeCode string }`
  - `type ProductCandidate struct { Codprod int; Produto string; EAN string; Reference string; QualityFlags []QualityFlag }`
  - `type SellableStock struct { Codprod int; Quantity float64; Scope StockScope; QualityFlags []QualityFlag; SourceFetchedAt time.Time }`
  - `type CostAsOf struct { Codprod int; Codemp int; SaleDate string; CUSSEMICM *float64; QualityFlags []QualityFlag }`
  - `type Reader interface { FindProductsForLinking(...); GetSellableStock(...); GetCurrentPrice(...); GetCostAsOf(...); GetSalesHistory(...); GetTaxInputs(...) }`

- [ ] **Step 1: Write the feature spec**

```md
# F-01 Spec

## MNOS Evidence
- `C:\Users\leandro.theodoro\Documents\MNOS\semantic\views\vw_estoque_saldo.sql`
- `C:\Users\leandro.theodoro\Documents\MNOS\semantic\views\vw_preco_tabela.sql`
- `C:\Users\leandro.theodoro\Documents\MNOS\semantic\views\vw_fat_venda_item.sql`
- `C:\Users\leandro.theodoro\Documents\MNOS\semantic\views\vw_imposto_item.sql`
- `C:\Users\leandro.theodoro\Documents\MNOS\semantic\governance\tgfcus.yml`
- `C:\Users\leandro.theodoro\Documents\MNOS\semantic\governance\tgfpro.yml`

## Acceptance Criteria
- Domain contract names every required read operation from IC-002.
- Stock scope preserves `CODEMP IN (1,2)` + `CODLOCAL=10101`.
- Quality flags include `missing_product`, `missing_stock`, `missing_cost`, `missing_tax`, `ambiguous_product`, `stale_source`.
- Contract names `CUSSEMICM` as the initial cost basis.
```

- [ ] **Step 2: Write the feature plan**

```md
# F-01 Plan

1. Add domain value objects and quality flags.
2. Add read port signatures aligned with IC-002.
3. Record MNOS evidence in validation.
4. Run targeted compile/test command after file creation.
```

- [ ] **Step 3: Write the failing compile/test check**

```go
package ports

import "testing"

func TestReaderContractCompiles(t *testing.T) {
	var _ Reader
}
```

Run: `cd apps/server_core; $env:GOCACHE=(Resolve-Path .gocache); go test ./internal/modules/internal_read/...`
Expected: FAIL because package/files do not exist yet.

- [ ] **Step 4: Add minimal domain and port implementation**

```go
package domain

type QualityFlag string

const (
	QualityComplete         QualityFlag = "complete"
	QualityMissingProduct   QualityFlag = "missing_product"
	QualityMissingStock     QualityFlag = "missing_stock"
	QualityMissingCost      QualityFlag = "missing_cost"
	QualityMissingTax       QualityFlag = "missing_tax"
	QualityAmbiguousProduct QualityFlag = "ambiguous_product"
	QualityStaleSource      QualityFlag = "stale_source"
)
```

```go
package ports

import (
	"context"
	"time"

	"marketplace-central/apps/server_core/internal/modules/internal_read/domain"
)

type Reader interface {
	FindProductsForLinking(ctx context.Context, input FindProductsInput) ([]domain.ProductCandidate, error)
	GetSellableStock(ctx context.Context, input SellableStockInput) (domain.SellableStock, error)
	GetCurrentPrice(ctx context.Context, input CurrentPriceInput) (domain.CurrentPrice, error)
	GetCostAsOf(ctx context.Context, input CostAsOfInput) (domain.CostAsOf, error)
	GetSalesHistory(ctx context.Context, input SalesHistoryInput) (domain.SalesHistory, error)
	GetTaxInputs(ctx context.Context, input TaxInput) (domain.TaxInputs, error)
}

type SellableStockInput struct {
	Codprod   int
	Companies []int
	Locations []int
	Now       time.Time
}
```

- [ ] **Step 5: Run targeted tests**

Run: `cd apps/server_core; $env:GOCACHE=(Resolve-Path .gocache); go test ./internal/modules/internal_read/...`
Expected: PASS with only domain/ports package tests compiling.

- [ ] **Step 6: Update F-01 validation artifact**

```md
## Quick Validation
- Command: `go test ./internal/modules/internal_read/...`
- Evidence type: ran
- Result: pass
- Changed paths: domain/*.go, ports/reader.go
```

- [ ] **Step 7: Commit**

```bash
git add .mnfs/MIS-001-mercado-livre-operating-cockpit/M-03-mnos-sankhya-read-contract/F-01-sankhya-read-contract-import apps/server_core/internal/modules/internal_read/domain apps/server_core/internal/modules/internal_read/ports
git commit -m "feat(internal_read): add sankhya read contract types"
```

### Task 2: F-02 Application Service And Adapter Seams

**Files:**
- Create: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-03-mnos-sankhya-read-contract/F-02-read-adapter-seam/spec.md`
- Create: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-03-mnos-sankhya-read-contract/F-02-read-adapter-seam/plan.md`
- Create: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-03-mnos-sankhya-read-contract/F-02-read-adapter-seam/validation.md`
- Create: `apps/server_core/internal/modules/internal_read/application/service.go`
- Create: `apps/server_core/internal/modules/internal_read/application/service_test.go`
- Create: `apps/server_core/internal/modules/internal_read/adapters/fake/reader.go`
- Create: `apps/server_core/internal/modules/internal_read/adapters/oracle/reader.go`
- Create: `apps/server_core/internal/modules/internal_read/adapters/oracle/config.go`
- Modify: `apps/server_core/internal/composition/root.go`
- Check: `apps/server_core/internal/platform/msdb/pool.go`

**Interfaces:**
- Consumes:
  - `internal_read/ports.Reader`
  - `internal_read/domain.*`
- Produces:
  - `func NewService(reader ports.Reader) Service`
  - `func (s Service) GetSellableStock(ctx context.Context, input ports.SellableStockInput) (domain.SellableStock, error)`
  - `func NewReader(fixtures Fixtures) *Reader`
  - `func LoadConfigFromEnv(getenv func(string) string) (Config, error)`

- [ ] **Step 1: Write the failing service tests**

```go
func TestServiceDelegatesSellableStock(t *testing.T) {
	reader := fake.NewReader(fake.Fixtures{
		Stocks: map[int]domain.SellableStock{
			42664: {Codprod: 42664, Quantity: 3},
		},
	})
	svc := application.NewService(reader)

	got, err := svc.GetSellableStock(context.Background(), ports.SellableStockInput{Codprod: 42664})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Quantity != 3 {
		t.Fatalf("expected 3, got %v", got.Quantity)
	}
}
```

Run: `cd apps/server_core; $env:GOCACHE=(Resolve-Path .gocache); go test ./internal/modules/internal_read/application/...`
Expected: FAIL because service and fake adapter do not exist yet.

- [ ] **Step 2: Add the minimal service and fake adapter**

```go
package application

type Service struct {
	reader ports.Reader
}

func NewService(reader ports.Reader) Service {
	return Service{reader: reader}
}

func (s Service) GetSellableStock(ctx context.Context, input ports.SellableStockInput) (domain.SellableStock, error) {
	return s.reader.GetSellableStock(ctx, input)
}
```

```go
package fake

type Fixtures struct {
	Products map[string][]domain.ProductCandidate
	Stocks   map[int]domain.SellableStock
	Costs    map[string]domain.CostAsOf
	Taxes    map[string]domain.TaxInputs
}
```

- [ ] **Step 3: Add Oracle seam config and secret-safe tests**

```go
func TestLoadConfigFromEnvDoesNotLeakSecretValues(t *testing.T) {
	_, err := oracle.LoadConfigFromEnv(func(key string) string {
		switch key {
		case "SANKHYA_DSN":
			return ""
		case "SANKHYA_PASSWORD":
			return "super-secret"
		default:
			return ""
		}
	})
	if err == nil {
		t.Fatal("expected error")
	}
	if strings.Contains(err.Error(), "super-secret") {
		t.Fatal("error leaked secret value")
	}
}
```

Run: `cd apps/server_core; $env:GOCACHE=(Resolve-Path .gocache); go test ./internal/modules/internal_read/...`
Expected: FAIL until `LoadConfigFromEnv` names only missing keys and Oracle reader remains read-only.

- [ ] **Step 4: Implement the guarded Oracle seam**

```go
package oracle

type Config struct {
	DSN string
}

func LoadConfigFromEnv(getenv func(string) string) (Config, error) {
	dsn := strings.TrimSpace(getenv("SANKHYA_DSN"))
	if dsn == "" {
		return Config{}, fmt.Errorf("missing required env: SANKHYA_DSN")
	}
	return Config{DSN: dsn}, nil
}
```

```go
package oracle

type Reader struct{}

func NewReader(_ Config) *Reader { return &Reader{} }
```

- [ ] **Step 5: Wire composition only if a non-HTTP consumer needs it**

```go
// Keep root wiring compile-safe only.
// Do not register HTTP routes for internal_read in M-03.
```

Run: `cd apps/server_core; $env:GOCACHE=(Resolve-Path .gocache); go test ./internal/...`
Expected: PASS with no route changes and no secret leaks in test output.

- [ ] **Step 6: Update F-02 validation artifact**

```md
## Quick Validation
- Command: `go test ./internal/modules/internal_read/...`
- Command: `go test ./internal/...`
- Evidence type: ran
- Result: pass
- Notes: Oracle seam remains read-only and route-less.
```

- [ ] **Step 7: Commit**

```bash
git add .mnfs/MIS-001-mercado-livre-operating-cockpit/M-03-mnos-sankhya-read-contract/F-02-read-adapter-seam apps/server_core/internal/modules/internal_read/application apps/server_core/internal/modules/internal_read/adapters apps/server_core/internal/composition/root.go
git commit -m "feat(internal_read): add application and adapter seam"
```

### Task 3: F-03 Quality Rules And Contract Tests

**Files:**
- Create: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-03-mnos-sankhya-read-contract/F-03-data-quality-rules/spec.md`
- Create: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-03-mnos-sankhya-read-contract/F-03-data-quality-rules/plan.md`
- Create: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-03-mnos-sankhya-read-contract/F-03-data-quality-rules/validation.md`
- Create: `apps/server_core/internal/modules/internal_read/domain/quality_flag_test.go`
- Create: `apps/server_core/internal/modules/internal_read/adapters/fake/reader_test.go`
- Modify: `apps/server_core/internal/modules/internal_read/domain/internal_stock.go`
- Modify: `apps/server_core/internal/modules/internal_read/domain/internal_cost.go`
- Modify: `apps/server_core/internal/modules/internal_read/domain/internal_tax.go`
- Modify: `wiki/modules/product-links.md`
- Modify: `wiki/modules/inventory.md`
- Modify: `wiki/modules/profitability.md`

**Interfaces:**
- Consumes:
  - `domain.QualityFlag`
  - `fake.Fixtures`
- Produces:
  - tests proving `M-03-C01`
  - tests proving `M-03-C02`
  - documented quality semantics for downstream modules

- [ ] **Step 1: Write failing contract tests for stock semantics**

```go
func TestFakeReaderSellableStockExcludesShowroomByDefault(t *testing.T) {
	reader := fake.NewReader(fake.Fixtures{
		Stocks: map[int]domain.SellableStock{
			42664: {
				Codprod:  42664,
				Quantity: 3,
				Scope: domain.StockScope{
					Companies: []int{1, 2},
					Locations: []int{10101},
					Formula:   "SUM(ESTOQUE - RESERVADO)",
				},
			},
		},
	})

	got, err := reader.GetSellableStock(context.Background(), ports.SellableStockInput{Codprod: 42664})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Quantity != 3 {
		t.Fatalf("expected revenda stock only, got %v", got.Quantity)
	}
}
```

```go
func TestFakeReaderMissingCostRemainsFlagged(t *testing.T) {
	reader := fake.NewReader(fake.Fixtures{})
	got, err := reader.GetCostAsOf(context.Background(), ports.CostAsOfInput{Codprod: 42664, Codemp: 1, SaleDate: "2026-07-06"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.CUSSEMICM != nil {
		t.Fatalf("expected nil cost, got %v", *got.CUSSEMICM)
	}
	if !slices.Contains(got.QualityFlags, domain.QualityMissingCost) {
		t.Fatalf("expected missing_cost flag, got %v", got.QualityFlags)
	}
}
```

Run: `cd apps/server_core; $env:GOCACHE=(Resolve-Path .gocache); go test ./internal/modules/internal_read/...`
Expected: FAIL until fake adapter default behavior and quality semantics are implemented.

- [ ] **Step 2: Implement minimal quality-aware fake behavior**

```go
func (r *Reader) GetCostAsOf(ctx context.Context, input ports.CostAsOfInput) (domain.CostAsOf, error) {
	key := costKey(input.Codprod, input.Codemp, input.SaleDate)
	if cost, ok := r.fixtures.Costs[key]; ok {
		return cost, nil
	}
	return domain.CostAsOf{
		Codprod:      input.Codprod,
		Codemp:       input.Codemp,
		SaleDate:     input.SaleDate,
		CUSSEMICM:    nil,
		QualityFlags: []domain.QualityFlag{domain.QualityMissingCost},
	}, nil
}
```

```go
func (r *Reader) GetTaxInputs(ctx context.Context, input ports.TaxInput) (domain.TaxInputs, error) {
	return domain.TaxInputs{
		Codprod:      input.Codprod,
		QualityFlags: []domain.QualityFlag{domain.QualityMissingTax},
	}, nil
}
```

- [ ] **Step 3: Update module wiki semantics**

```md
- Inventory consumes `internal_read` sellable stock and treats `missing_stock` or `stale_source` as blocked states.
- Product Links consumes `internal_read` candidate search and treats `missing_product` or `ambiguous_product` as unresolved states.
- Profitability consumes `internal_read` cost/tax inputs and treats missing values as quality flags, never `0`.
```

- [ ] **Step 4: Run targeted tests and contract searches**

Run: `cd apps/server_core; $env:GOCACHE=(Resolve-Path .gocache); go test ./internal/modules/internal_read/...`
Expected: PASS

Run: `rg -n "CUSVARIAVEL|INSERT|UPDATE|DELETE" apps/server_core/internal/modules/internal_read`
Expected: no `CUSVARIAVEL`; no write-path SQL in Oracle adapter seam.

- [ ] **Step 5: Update F-03 validation artifact**

```md
## Quick Validation
- Command: `go test ./internal/modules/internal_read/...`
- Command: `rg -n "CUSVARIAVEL|INSERT|UPDATE|DELETE" apps/server_core/internal/modules/internal_read`
- Evidence type: ran
- Result: pass
```

- [ ] **Step 6: Commit**

```bash
git add .mnfs/MIS-001-mercado-livre-operating-cockpit/M-03-mnos-sankhya-read-contract/F-03-data-quality-rules apps/server_core/internal/modules/internal_read wiki/modules/product-links.md wiki/modules/inventory.md wiki/modules/profitability.md
git commit -m "feat(internal_read): enforce quality flag semantics"
```

### Task 4: Milestone Review, Acceptance, And Validation Result

**Files:**
- Modify: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-03-mnos-sankhya-read-contract/milestone.md`
- Create: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-03-mnos-sankhya-read-contract/validation-result.md`
- Modify: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-03-mnos-sankhya-read-contract/F-01-sankhya-read-contract-import/validation.md`
- Modify: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-03-mnos-sankhya-read-contract/F-02-read-adapter-seam/validation.md`
- Modify: `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-03-mnos-sankhya-read-contract/F-03-data-quality-rules/validation.md`

**Interfaces:**
- Consumes:
  - all feature `spec.md`, `plan.md`, `validation.md`
  - current `go test` and `rg` evidence
- Produces:
  - `validation-result.md` with explicit M-03-C01 and M-03-C02 verdicts
  - milestone status grounded in evidence

- [ ] **Step 1: Review each feature return against MNFS acceptance bar**

```md
Acceptance requires:
- `spec.md` present
- `plan.md` present
- changed paths listed
- `validation.md` present
- evidence type `ran` for load-bearing checks
```

- [ ] **Step 2: Run full milestone validation commands**

Run: `cd apps/server_core; $env:GOCACHE=(Resolve-Path .gocache); go test ./internal/modules/internal_read/...`
Expected: PASS

Run: `rg -n "CUSVARIAVEL|INSERT|UPDATE|DELETE" apps/server_core/internal/modules/internal_read`
Expected: no write-path evidence, no wrong cost basis.

Run: `rg -n "SANKHYA_DSN|SANKHYA_PASSWORD|password=" .mnfs apps/server_core/internal/modules/internal_read`
Expected: no secret values leaked in artifacts or module errors.

- [ ] **Step 3: Write milestone validation result**

```md
# M-03 Validation Result

## Verdict
Pass

## Criterion Results
- M-03-C01: Pass
  - Evidence: fake-adapter/unit tests prove `SUM(ESTOQUE - RESERVADO)`, `CODEMP IN (1,2)`, `CODLOCAL=10101`, showroom exclusion for `10108`.
- M-03-C02: Pass
  - Evidence: unit tests prove `missing_product`, `missing_cost`, `missing_tax`; no missing numeric converted to zero.

## Blocking Failure Check
- No Sankhya write path found.
- No `CUSVARIAVEL` fallback found.
- No secret leak found.
```

- [ ] **Step 4: Update milestone status only if validation-result proves pass**

```md
- status: passed only after `validation-result.md` says `Verdict: Pass`
```

- [ ] **Step 5: Commit**

```bash
git add .mnfs/MIS-001-mercado-livre-operating-cockpit/M-03-mnos-sankhya-read-contract
git commit -m "docs(m03): record validation result"
```

## Self-Review

- Spec coverage:
  - internal read module boundary: Task 1 and Task 2
  - fake/test seam: Task 2 and Task 3
  - quality flags and zero-default prevention: Task 3
  - M-03 feature artifacts and validation-result: Task 1 through Task 4
  - secret-safety and no-write checks: Task 2 through Task 4
- Placeholder scan:
  - no `TBD`, `TODO`, or deferred implementation markers remain
- Type consistency:
  - `domain.QualityFlag`, `ports.Reader`, `application.Service`, and fake/oracle adapter signatures are introduced in earlier tasks before later tasks depend on them

## Automatic Execution Choice

Proceed with **Subagent-Driven** execution automatically.

Reason:
- objective explicitly requires Milestone Orchestrator flow and feature implementation via subagents
- F-01, F-02, and F-03 have separable artifacts and review points
- main thread should stay focused on acceptance review and milestone validation
