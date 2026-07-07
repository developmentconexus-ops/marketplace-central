# Task 1 Report

Status: DONE

Feature: F-01 Contract Import And Domain Surface

Summary:
- Added the MPC-owned `internal_read` domain surface for quality flags, product candidates, default sellable stock scope, price, cost, tax, and sales reads.
- Added the read-only `ports.Reader` contract for all six IC-002 operations.
- Preserved exact stock defaults: `SUM(ESTOQUE - RESERVADO)`, `CODEMP IN (1,2)`, `CODLOCAL=10101`, with `10108` excluded from the default scope.
- Preserved `CUSSEMICM` as nullable initial cost basis; missing cost is represented by `nil` plus `missing_cost`, not zero.
- Created F-01 `validation.md` with quick-validation evidence.

Validation:
- RED: `cd apps/server_core; $env:GOCACHE=(Resolve-Path .gocache).Path; go test ./internal/modules/internal_read/...` failed before implementation on missing `Reader`, `DefaultSellableStockScope`, scope constants, and quality flags.
- GREEN: `cd apps/server_core; $env:GOCACHE=(Resolve-Path .gocache).Path; go test ./internal/modules/internal_read/...` passed for `domain` and `ports`.

Changed paths:
- `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-03-mnos-sankhya-read-contract/F-01-sankhya-read-contract-import/spec.md`
- `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-03-mnos-sankhya-read-contract/F-01-sankhya-read-contract-import/plan.md`
- `.mnfs/MIS-001-mercado-livre-operating-cockpit/M-03-mnos-sankhya-read-contract/F-01-sankhya-read-contract-import/validation.md`
- `apps/server_core/internal/modules/internal_read/domain/*.go`
- `apps/server_core/internal/modules/internal_read/domain/contract_test.go`
- `apps/server_core/internal/modules/internal_read/ports/reader.go`
- `apps/server_core/internal/modules/internal_read/ports/reader_test.go`

Concerns:
- None for Task 1 scope.

Fix cycle after task review:
- Updated `CurrentPriceInput` so `Codtab` and `Codlocal` are optional pointers, preserving IC-002 default lookup semantics.
- Updated `SalesHistoryInput` so the seam supports product or group queries via `Codprod` or `CodgrupoProd`.
- Re-ran focused validation after the contract-shape fix.

Fix validation:
- `cd apps/server_core; $env:GOCACHE=(Join-Path (Get-Location) '.gocache'); go test ./internal/modules/internal_read/...`
- Result: pass

---

Follow-up fix:
- Narrowed `CurrentPriceInput` so `Codtab` and `Codlocal` are optional pointers, allowing implicit defaults for `CODTAB=0` and `CODLOCAL=10101`.
- Expanded `SalesHistoryInput` with `CodgrupoProd` alongside `Codprod`, so the read port can express product-or-group sales windows per IC-002.
- Verified with `cd apps/server_core; $env:GOCACHE=(Resolve-Path .gocache); go test ./internal/modules/internal_read/...`.
