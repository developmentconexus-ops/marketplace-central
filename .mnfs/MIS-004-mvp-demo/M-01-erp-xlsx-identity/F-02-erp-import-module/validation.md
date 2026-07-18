# F-02 erp-import-module — validation evidence

Feature = the ERP xlsx import path: POST xlsx → synchronous import → protocol `#NNN-E`, honest
per-row accept/reject report, dedup by file SHA-256, atomic persistence. Maps to validation-contract
**C01** (import), **C02** (rejection report), **C06** (migrations + seams).

Branch `chip/m01-erp-xlsx-identity` · base `59d0e62f`. Failing-test-first, independent review per slice.

| Slice | What | Dispatch | Commit | Review verdict |
|---|---|---|---|---|
| F02-S1 | erp_import domain + ports + issue-code constants | D06 (luna/high) | `f73df384` | PASS — enum/port/tenant/float/ADR-17/dedup clean (2🟡 fixed) |
| F02-S2 | migrations 0045–0047 (tenant-scoped, protocol/tenant UNIQUE, composite FKs) | D04 (luna/high) | `c4fb75bd` | PASS + 3 orch hardening (protocol `{3,}`, composite tenant FK ×2, UNIQUE tenant/protocol) |
| F02-S3 | xlsx parse adapter + excelize v2.11.0 (DEP-GRANT-01, ONE file) | D08 (luna/high) | `55e84c36` | PASS — 8/8; excelize isolated, pure mapper, typed FileError, ADR-17 nil |
| F02-S4 | postgres import repository (+ports/errors.go) | D09 (sol/low) | `52605864` | PASS — tenancy/atomicity/lock/status/typed-errors clean (4🟡 fixed) |
| F02-S5 | application ImportService + QueryService (RunImport, det clock/IDs) | D11 (sol/low) | `758a2a25` | PASS — 1🟡 (prod UUID default uncovered) fixed |
| F02-S6 | HTTP transport (multipart POST + list/detail GET, tenant-free handler) | D12 (luna/high) | `a568dbe7` | PASS — full error matrix, []-not-null issues, RFC3339-UTC (4🟡 fixed) |
| F02-S7 | OpenAPI /erp/imports* paths + ErpImport* schemas + sdk-runtime types (ADR-12) | D14 (luna/high) | `894ece73` | PASS — flat lowercase errors, status/issue-code enums mirror wire (1🔴 500-decl fixed) |

## Contract criteria satisfied
- **C01 import:** POST xlsx → 201 sync + protocol `#NNN-E`; re-POST same file → 409 dedup by
  `file_sha256`; no raw row payload leaked to logs (D12 no-leak review). Atomic persist + advisory-lock
  contention → ErrImportInProgress (postgres repo, R4-hardened).
- **C02 rejection report:** empty DESCRPROD → REJECTED `EMPTY_DESCRPROD`; invalid EAN → imported +
  `INVALID_EAN` warning + ean null; malformed NCM → imported + `INVALID_NCM` warning + ncm null; 100%
  rejected → REJECTED. Pinned in domain/validation_test.go + source_contract_test.go (C02 block).
- **C06 migrations + seams:** migrations 0045–0047 present, tenant-scoped, UP-only plain SQL (repo loader
  idiom); erp_import governance entry (F03-S2); composition wires the module (F03-S2).

## Notes / carried findings
- **DEP-GRANT-01** excelize/v2 imported in EXACTLY ONE file (xlsx/parser.go); grep-verified clean; go.mod+go.sum
  landed in the same slice (F02-S3).
- Error family = flat lowercase (`invalid_file`, `duplicate_file`, `import_in_progress`); issue codes = UPPERCASE
  (`EMPTY_DESCRPROD`, `INVALID_EAN`, `INVALID_NCM`). **R1 open:** validation-contract.md:37 prose says
  `DUPLICATE_FILE` (uppercase) — chip recommends lowercase stands; hub ruling owed.
- F02-S2 schema over-addition (`filename_sanitized`) caught + corrected in D09 (ledger §Field findings).

**Status: COMPLETE** (7/7 slices merged, all reviewed green). Aggregate ladder green @ db91f385.
