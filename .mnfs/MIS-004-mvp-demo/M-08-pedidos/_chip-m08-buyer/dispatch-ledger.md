# CHIP-M08-BUYER — Dispatch Ledger

- **Mission:** MIS-004-mvp-demo
- **Milestone:** M-08-pedidos
- **Chip:** CHIP-M08-BUYER (buyer fiscal identity on order-detail drawer for ERP registration)
- **Hub:** local_7c41fdd8-dd6f-487c-a7ea-1b511e214fa9
- **Branch:** chip/m08-buyer
- **BaseSha:** `0df200dce136ef43257f60252c2db908123c2e74`

## Commits (BaseSha..HEAD)

| SHA | Slice | Summary |
|-----|-------|---------|
| bb0f7df3 | 1 | feat(connectors): buyer fiscal reader (ML two-step billing-info flow) — domain BuyerFiscalInfo, connectors port, ML adapter reader + tests |
| e581a288 | 2 | feat(orders): read-time buyer fiscal enrichment — orders port, EnrichService.resolveBuyerFiscal (warn-once degrade), tests |
| 75c0bfc9 | 3 | feat(orders): surface comprador_fiscal on order detail DTO + wire adapter (composition bridge installationID→accountRef) |
| ccb68d67 | 4 | feat(contract): comprador_fiscal on OrderRead (OpenAPI + sdk-runtime, same commit) |
| c75141f3 | 5 | feat(web): render buyer fiscal identity + carrier/destino/frete on pedido drawer |
| 2f089604 | 6 | fix(chip): P6 gate — de-contaminate SDK (drop M-05 signal types) + gofmt var block |

## Gate results

- **Backend:** `go build ./apps/server_core/...` OK · `go vet` OK · `go test` (connectors ML + orders/…) all `ok`, no FAIL/panic.
- **Frontend:** vitest pedidos 23/23 · full web vitest 322 pass · vite build ✓ · SDK `tsc --noEmit` clean.
- **gofmt:** all seven touched Go files clean after 2f089604.
- **Governance lane:** `harness.ps1 governance -BaseSha <40hex>` → `status=passed` (only pre-existing baseline exceptions; no new drift).
- **P6 dual gate:** cold Opus + adversarial sonnet, independent read of fix.diff.
  - Round 1: both FAIL — (Opus) SDK carried M-05 listing-signal types with no OpenAPI source (scope contamination + OpenAPI↔SDK invariant break); (sonnet) capability_adapter.go var-block not gofmt-clean.
  - Fix: commit 2f089604 removes M-05 SDK types (keeps only comprador_fiscal), gofmt-realigns var block.
  - Round 2: re-verified against corrected diff (verdicts recorded in EVIDENCE.md).

## Binding constraints honored

- READ-only (both ML calls `http.MethodGet`); MPC_PROVIDER_WRITES_ENABLED unset; zero ML writes.
- billing-info 404/undecodable → nil silently at adapter; other errors propagate → application warns once. No WARN spam on the normal "none" path.
- LGPD: identification.number / doc_numero never logged (grep-verified); rendered by client only.
- ADR-17: honest-absence nil pointers, omitempty DTO, honest "—" UI; doc_tipo opaque (never mapped to CPF/CNPJ); invoice_type not decoded.
- Additive nullable; NO migration. OpenAPI + sdk-runtime landed in the same commit (ccb68d67).
- Hexagonal: provider structs die at adapter → neutral domain type; connectors port keyed by ProviderAccountRef, orders port keyed by installationID, composition bridges the two.
- Tenant scoping intact. No push. `git branch -d` only. No reset/revert/stash/clean/WSL. All writes to the worktree path.

## Out of scope for the chip (hub owns)

- Merge to main.
- Live P7 browser QA drive of the drawer.
