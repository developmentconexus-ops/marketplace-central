# Step A — Two-source snapshot resolution (hub-authored)

**Commit**: `1ee5f917` on `main` (parent `500ec1dd`). Hub-authored (not a chip); dual-gated per stop-hook.
**Date**: 2026-07-20 (D-113, demo day).

## Intent
The erp_import reader served whichever completed import had the newest `imported_at`.
Importing a lenient prospect catalog (`source=catalogo_cliente`, #004-E) therefore
hijacked the active dataset and displaced the real Sankhya ERP snapshot
(`source=xlsx`, #003-E). Step A makes `LatestCompletedSnapshot` source-aware,
defaulting to xlsx so the demo opens on real ERP data. Prospect remains reachable
via the future toggle layer (Steps B/C, deferred by operator D-113).

## Write-set (git show 1ee5f917 --name-only) — 8 files, all under erp_import/
- ports/repository.go — `LatestCompletedSnapshot(ctx, tenant, source domain.ImportSource)`
- adapters/postgres/query_repository.go — `WHERE tenant_id=$1 AND source=$2 AND status='COMPLETED'`
- adapters/postgres/import_repository.go — empty `snapshot.Source` → `SourceXLSX` (protects 0072 CHECK + pre-source fixtures)
- adapters/internalread/reader.go — `WithActiveSource`/`activeSourceFromContext` (default `SourceXLSX`); `snapshot()` threads it into all 5 reader methods
- adapters/internalread/reader_test.go, source_contract_test.go, application/import_service_test.go — 3 fakes updated to new signature
- adapters/postgres/query_repository_test.go — integration callers + regression `TestLatestCompletedSnapshotReadsNullCostAndStock`

## Gates
- `go build ./...` exit 0; `go test ./internal/modules/erp_import/...` all ok; `go vet -tags integration` (postgres pkg) clean.
- **Live verify**: backend rebuilt healthy. `GET /catalog/products` (default source) serves Sankhya #003-E (codprods 412+, real cost 42.1/296.51/74.25, quality=complete) — NOT the NULL-cost prospect. Both datasets persist (`GET /erp/imports`: #004-E catalogo_cliente + #003-E/#002-E/#001-E xlsx).

## P6 DUAL GATE
- **Cold gate** (harness:gate-reviewer, read-only): COLD-GATE: PASS. Criteria 1–6 evidenced with file:line; criterion 7 (scope confinement) NOT-EVIDENCED by the reviewer (no git tool in its sandbox) — **closed by hub**: `git show 1ee5f917 --name-only` = exactly 8 files, all under `erp_import/`, zero out-of-scope.
- **Adversarial refuter** (read-only skeptic): REFUTER-VERDICT: NO-REFUTATION. All 6 attack vectors CANNOT-REFUTE — (1) every implementer updated (10 grep hits, sole caller reader.go:69); (2) comma-ok assertion on unexported ctx key, no panic, deterministic xlsx default; (3) runImport sets Source before the status branch so both COMPLETED+REJECTED carry non-empty source → 0072 CHECK satisfied on all paths (empty→xlsx is a belt for fixtures); (4) positional args + tenant scoping intact, sibling queries untouched; (5) ADR-17 nullableString + honest-unknown flags preserved end-to-end; (6) single source in play in Step A → no cache cross-source poisoning (latent for Step B only).

Both gates agree PASS.

P6-DUAL-GATE: AGREEMENT
