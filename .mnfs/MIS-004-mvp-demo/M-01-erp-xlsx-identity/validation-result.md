# M-01 — Validation Result (P7 QA live-drive)

- **Verdict: PASS** — milestone closed.
- Date: 2026-07-18. Main @ 69764da (merge f19d4376 + hub fixes 332c592c/128e2e8b/8afd384a/69764da).
- QA: fresh independent persona (sonnet subagent), live-drive against dev stack (backend :8080, `MC_ERP_SOURCE=xlsx`, dev postgres :5435). Hub authored this result from the QA report; ladder evidence cited for C04–C06.

## §import (C01 — PASS)
- POST `/erp/imports` multipart `example-erp.xlsx` → **201 sync** `{"import_id":"ce6bd05c-0280-4b24-988c-d20413c90940","protocol":"#001-E","status":"COMPLETED"}`. Protocol matches `#NNN-E`.
- Re-POST same file → **409** `{"error":"duplicate_file", ...}` — lowercase wire code per D-13 ruling. DB: exactly one `erp_import_protocols` row; no duplicate snapshot.
- GET list (newest-first) + GET detail 200, counts `accepted:55 rejected:0 warning:0`.
- Logs: zero hits for row payloads (`Example Product`, custo values, refforn) — no raw payload leakage.
- Note (non-blocking): DB check-constraint enum stores `DUPLICATE_FILE` uppercase internally; wire is lowercase. Internal-only.

## §rejeicao (C02 — PASS)
- POST `identity-rejections.xlsx` → 201 `#002-E` COMPLETED (partial-issues import).
- Detail: `accepted:3 rejected:1 warning:2`; `rejected_rows=[{row:1, code:EMPTY_DESCRPROD, column:DESCRPROD, offending_value:""}]`; warnings `INVALID_EAN` (row 2, `7894900011518`) + `INVALID_NCM` (row 3, `12AB`).
- DB `erp_import_issues`: identical 3 rows (kind/code/offending_value preserved). `erp_import_products`: codprod 2002 `ean IS NULL` (refforn retained), codprod 2003 `ncm IS NULL` — real SQL NULL, never empty string.

## §identidade (C03 — PASS)
- `GET /catalog/products` + `/catalog/products/search`: `ean`/`ncm` are JSON `null` where invalid, never `""`; `quality_flags` correct (`invalid_ean` on 2002, `complete` on 2003/2004).
- Catalog serves latest-completed-snapshot only (by design, `LatestCompletedSnapshot`); first-import identity verified via `erp_import_products` rows for protocol ce6bd05c (codprod 1001: custo 13.34, ean 7890000000017, refforn REF-1001, marca "Synthetic Brand", ncm 12345678) — matches xlsx source exactly.

## §reader (C04 — PASS, ladder evidence)
- Hermetic reader-port integration test `TestReaderRealRepositoryCostReservedAsOfAndRejectedIgnored` green on session pg (run after hub fixture fix 128e2e8b): sellable=7 (9−2), cost-as-of 12.30 with observed_at, pre-snapshot as-of → `ErrNoErpSnapshot`, rejected snapshot ignored. Integration lane run 26f51f13 status=passed migrations_first=44.

## §oracle (C05 — PASS, ladder evidence)
- Oracle path intact: full `go build ./... && go vet ./...` green, full test sweep zero FAIL at 69764da. `erpSource()` default remains `"oracle"`; xlsx opt-in via `MC_ERP_SOURCE` (registered in runtime-config.json, approved reader root.go).
- Live-drive: zero `panic|goroutine|fatal error` in backend logs under xlsx mode (nil Oracle). Gate-A watch item confirmed: degrade, not panic.

## §seams (C06 — PASS, ladder evidence)
- Migrations 0045–0049 only; post-merge migrations_first=44→49 consistent; governance lane status=passed vs base 59d0e62f (clean detached worktree, full 40-hex BaseSha).
- SDK: tsc 0 errors, vitest 63/63; BARREL-01 (`export * from "./erpImport"`) applied per grant.
- Dev-stack seam (hub-owned): compose gained `MC_ERP_SOURCE: ${MC_ERP_SOURCE:-xlsx}` on backend for demo mode; overridable, code default untouched.

## Findings (non-blocking, carried forward)
- **F-QA-M01-1**: `GET /catalog/products/{id}` → 500 `CATALOG_INTERNAL_ERROR` in xlsx mode — `GetCurrentPrice` path hits nil Oracle (`oracle_code=unknown`), gracefully handled but blocks product-detail page. List/search unaffected. **Owner: M-06 produto-detalhe** (must serve detail from snapshot/internal_read in xlsx mode). Fold into M-06 context pack.
