# DISPATCH LEDGER — CHIP-M01 (M-01 erp-xlsx-identity)

```yaml
milestone: M-01-erp-xlsx-identity
mission: MIS-004-mvp-demo
chip_session: CHIP-M01
accepted_base_sha: 59d0e62fdbf15db068542432ef5d5731b6fa9f83
base_subject: "merge(chip-sat): W1 CHIP-SAT ledger close-out docs"
worktree: .claude/worktrees/chip-m01-erp-xlsx-identity
branch: chip/m01-erp-xlsx-identity
hub_session: local_c5532363-d407-49ff-9bac-120c328b9dcf  # CONFIRMED "HUB: MIS-004 marketplace-central hub" (title verified via list_sessions 2026-07-17)
db_specialist: local_ec787804-f8e9-4981-9c12-7d3f45292294
impl_pack_version: v1.0.0
doctrine: docs/HARNESS-CORE.md + docs/HARNESS-PROFILE.md @ base (amendment log thru 2026-07-16; §11/§12 review bindings present)
```

## Boot findings (report to hub)
- No session titled "HUB: MIS-004" exists; assuming "Harness Hub" `local_c5532363` (right cwd, spawned wave-A siblings). CHIP-M02 `local_c6df14c8` + CHIP-M03 `local_1477e70a` running parallel.
- Governance HEAD `claude/confident-allen-9c7ab6` @ 3d158885 diverged from base (base NOT ancestor); doctrine content equivalent (base amendment log already carries the §11/§12 review hardening). Non-blocking.
- **xlsx lib ABSENT** (no excelize/tealeg/xuri in go.mod/go.work/go.sum) → F-02 parse slice needs a **dependency REQUEST** to hub. Blocks F-02 parse only; planning proceeds.
- Migrations top = `0044_market_references.sql`; `runner_test.go` count fixture = 41 (bump per migration).
- Modules present: catalog, internal_read (adapters: oracle/fake/cache + application/service.go), orders, product_links, market, etc. No erp_import (F-02 creates).
- gomodcache warmed (exit 0). codex-cli 0.144.4 present.

## Dep grants (binding)
- **DEP-GRANT-01** github.com/xuri/excelize/v2 (latest stable) — GRANTED by hub 2026-07-17 (hub row D-02 @ c6df5cc1). Conditions bound into F02-S3: (1) excelize imported in ONE file only (erp_import/adapters/xlsx/parser.go), behind ports.Parser; (2) go.mod+go.sum land in the SAME slice as the parse adapter (F02-S3), reviewed like any slice; (3) `go mod tidy` hermetic (GOMODCACHE local, GOCACHE absolute). F02-S3 now unblocked pending its F02-S1 dep.

## Environment
- gomodcache: `apps/server_core/.gomodcache` warmed 2026-07-17.
- node_modules: NOT yet installed (npm ci needed before web tsc/vitest lanes).
- Scratchpad (prompts/logs): session scratchpad dir.

## Dispatch rows
| # | Phase | Role | Model/Effort | Path | Prompt file | Output artifact | Status | Verdict |
|---|---|---|---|---|---|---|---|---|
| D01 | P2-evidence | Investigator | gpt-5.6-luna / medium | companion --wait | scratchpad/prompt-D01-evidence.md | scratchpad/D01-evidence.md | ABANDONED (hung 19min on pwsh grep; killed PID 24248) | — |
| D02 | P2-plan | Feature planner | gpt-5.6-sol / medium | OS-process | scratchpad/prompt-D02-planner.md | scratchpad/agent__D02-planner.last.md (+.log); PLAN.md | complete (exit 0) | **ACCEPTED** w/ 4 binding corrections (see PLAN §G); 14 slices F01=4·F02=7·F03=3 |
| — | REQUEST | DEP-GRANT-01 excelize/v2 | — | ccd send_message | — | hub local_c5532363 | **GRANTED** (hub ledger D-02 @ c6df5cc1) | proceed |
| D03 | impl F01-S1 | Implementer (standard) | gpt-5.6-luna / high | OS-process | scratchpad/prompt-D03-F01S1.md | scratchpad/agent__D03-F01S1.last.md (+.log); gtin.go + gtin_test.go | **MERGED @ 362107ec** | review PASS (sonnet, recomputed 4 checksums) |
| D04 | impl F02-S2 | Implementer (standard) | gpt-5.6-luna / high | OS-process | scratchpad/prompt-D04-F02S2.md | scratchpad/agent__D04-F02S2.last.md (+.log); 0045-0047 + runner_test.go | **MERGED @ c4fb75bd** | review PASS + 3 orch hardening (protocol {3,}, composite tenant FK ×2, UNIQUE tenant/protocol) |
| D05 | impl F01-S2 | Implementer (complex) | gpt-5.6-sol / low | OS-process | scratchpad/prompt-D05-F01S2.md | scratchpad/agent__D05-F01S2.last.md (+.log); internal_read domain+oracle identity remap | **MERGED @ aa9ef075** | deep review PASS (column-align, collision, REFFORN separation traced); catalog_page 2nd path deferred to F01-S3 |
| D06 | impl F02-S1 | Implementer (standard) | gpt-5.6-luna / high | OS-process | scratchpad/prompt-D06-F02S1.md | scratchpad/agent__D06-F02S1.last.md (+.log); erp_import domain+ports+codes | **MERGED @ f73df384** | review PASS (D06R, 2🟡: bare-decimal canon FIXED+tested, StockReserved silent-nil accepted-OOS) |
| D06R | review F02-S1 | Reviewer (independent) | claude sonnet | Agent bg | (inline) | (agent return) | done | PASS — 0 BLOCKER; enum/port/tenant/float/ADR-17/dedup all clean |
| D07 | impl F01-S3 | Implementer (complex) | gpt-5.6-sol / low | OS-process | scratchpad/prompt-D07-F01S3.md | scratchpad/agent__D07-F01S3.last.md (+.log); catalog-page identity + catalog DTOs (9 files) | **MERGED @ 6977a56f** | review PASS (D07R, 8/8 clean: Scan/column-order traced, no REFERENCIA leak, mirror-fidelity vs reader.go, null-not-empty, both CODPRODs) |
| D07R | review F01-S3 | Reviewer (independent) | claude sonnet | Agent bg | (inline) | (agent return) | done | PASS — 0 defect |
| D10 | impl F01-S4 | Implementer (standard) | gpt-5.6-luna / high | OS-process | scratchpad/prompt-D10-F01S4.md | scratchpad/agent__D10-F01S4.last.md (+.log); OpenAPI catalog schemas + manual SDK parity (atomic) | **MERGED @ a2c06985** | review 4🔴 (enum too narrow) FIXED pre-commit; F-01 COMPLETE |
| D10R | review F01-S4 | Reviewer (independent) | claude sonnet | Agent bg | (inline) | (agent return) | done | 4🔴: quality_flags enum `[invalid_ean,ean_collision]` under-wide vs unfiltered reader passthrough (real wire carries `complete` default + `missing_product`/`ambiguous_product`) → clients would reject valid responses. FIXED: shared QualityFlag component = full closed 11-value domain (spec⊇wire), $ref from both schemas + SDK union; std `npm run test` 60/60, tsc clean, legacy CatalogProduct + scope untouched. Verified in worktree. |
| D08 | impl F02-S3 | Implementer (standard) | gpt-5.6-luna / high | OS-process | scratchpad/prompt-D08-F02S3.md | scratchpad/agent__D08-F02S3.last.md (+.log); xlsx parse adapter + excelize v2.11.0 | **MERGED @ 55e84c36** | review PASS (D08R, 8/8: excelize isolated, pure mapper, accent-fold, typed FileError, ADR-17 nil, real fixtures) |
| D08R | review F02-S3 | Reviewer (independent) | claude sonnet | Agent bg | (inline) | (agent return) | done | PASS — 0 defect |
| D09 | impl F02-S4 | Implementer (complex) | gpt-5.6-sol / low | OS-process | scratchpad/prompt-D09-F02S4.md | scratchpad/agent__D09-F02S4.last.md (+.log); postgres import repo (+ports/errors.go) | **MERGED @ 52605864** | review PASS (D09R, 0🔴 4🟡 ALL FIXED: lock-key namespace, dedup unit test, ×2 ORDER BY id-DESC tiebreak) + 2 pre-review orch fixes |
| D09R | review F02-S4 | Reviewer (independent) | claude sonnet | Agent bg | (inline) | (agent return) | done | PASS — 0 BLOCKER; tenancy/atomicity/lock/status/columns/typed-errors/reads/skip-guard all clean; 4🟡 fixed pre-commit |
| D11 | impl F02-S5 | Implementer (complex) | gpt-5.6-sol / low | OS-process | scratchpad/prompt-D11-F02S5.md | scratchpad/agent__D11-F02S5.last.md (+.log); erp_import application ImportService+QueryService (4 files) | **MERGED @ 758a2a25** | worker BLOCKED-honest on wrong hardcoded SHA fixture (else 7/7 pass); orch fixed test to compute hash in-test; D11R PASS + 1🟡 fixed |
| D11R | review F02-S5 | Reviewer (independent) | claude sonnet | Agent bg | (inline) | (agent return) | done | PASS — 0🔴 1🟡: newUUIDv4 (prod default) uncovered (all tests override via WithIDGenerator) → FIXED, added TestNewUUIDv4IsCanonicalV4 (RFC-4122 v4 regex + uniqueness). Tenancy/ordering/FileError-short-circuit/zero-valid-persist/protocol/UUID-nibbles/%w-typed-errors/no-slop/test-honesty/boundary all clean. |
| D12 | impl F02-S6 | Implementer (standard) | gpt-5.6-luna / high | OS-process | scratchpad/prompt-D12-F02S6.md | scratchpad/agent__D12-F02S6.last.md (+.log); erp_import HTTP transport (multipart POST + list/detail GET) | **MERGED @ a568dbe7** | 2 files, tenant-free, full error matrix, RFC3339-UTC list, []-not-null issues; D12R 4🟡 fixed pre-commit |
| D12R | review F02-S6 | Reviewer (independent) | claude sonnet | Agent bg | (inline) | (agent return) ada2f74974980863f | done | 0🔴 4🟡 ALL FIXED pre-commit: (1)(2) malformed/missing-file 400 carried fabricated `detail` (contract reserves detail for parser INVALID_FILE only) → dropped; (3) oversize MaxBytesReader path untested → added streamed >25MiB test (MultiReader, no big alloc); (4) imported_at UTC never asserted despite BRT fixture → added TestHandlerImportedAtIsRFC3339UTC. Tenant-free/duplicate-original/all-rejected-201/no-leak all clean. |
| D13 | impl F03-S1 | Implementer (complex) | gpt-5.6-sol / low | OS-process | scratchpad/prompt-D13-F03S1.md | scratchpad/agent__D13-F03S1.last.md (+.log); erp_import→internal_read xlsx Reader adapter (reader.go + unit + integration tests) | **MERGED @ e314c74a** | 3 files only; gofmt/vet/build(both tags)/unit pass; integration skips clean no-DSN; D13R 1🔴 FIXED pre-commit (see D13R) |
| D13R | review F03-S1 | Reviewer (independent) | claude sonnet | Agent bg | (inline) | (agent return) a73a974ff76485516 | done | 1🔴 reader.go:212 candidate() hard-failed non-numeric CODPROD via fmt.Errorf("parse codprod")—aborts whole FindProductsForLinking + leaks offending value, inconsistent with catalogPage silent continue. FIXED: parse moved inline into loop, `continue` on parseErr/id≤0 (mirrors catalogPage); candidate() takes pre-parsed id; added TestFindProductsSkipsNonNumericCodprod. Re-verified green, committed e314c74a. Dual-match wrap, ERPProductNotFoundError-not-fabrication, reserved-"0", as-of<= boundary, repo-only-LatestCompletedSnapshot, skip-guard all clean. |
| D14 | impl F02-S7 | Implementer (standard) | gpt-5.6-luna / high | OS-process | scratchpad/prompt-D14-F02S7.md | scratchpad/agent__D14-F02S7.last.md (+.log); OpenAPI /erp/imports* paths + ErpImport* schemas + sdk-runtime/src/erpImport.ts + .test.ts | **MERGED @ 894ece73** | mirrors merged F02-S6 wire (flat lowercase errors, file_sha256/source, status[COMPLETED,REJECTED], issue-code enum). D14R 1🔴 FIXED pre-commit. Orch also widened vitest include src/**/*.test.ts (worker's "62/62" came from non-prescribed run; prescribed npm test only ran index.test.ts—new suite was orphaned). F-02 COMPLETE. |
| D15 | impl F03-S2 | Implementer (standard) | gpt-5.6-luna / high | OS-process | scratchpad/prompt-D15-F03S2.md | scratchpad/agent__D15-F03S2.last.md (+.log); composition source-selection + erp_import wiring + governance | **MERGED @ 45481943** | 3 files (root.go/root_test.go/modules.json). erpSource empty→oracle, oracle/xlsx trim+case-fold, else startup error; handler self-classifies routes → registerBatchRoutes untouched (no double-reg); oracle path preserved; xlsx leaves oracle-only readers unavailable (stock batch/profitability internal/sankhya/poolStats/oracleDB); repo tenant-free ctor. Orch re-verified warmed-cache: gofmt/vet clean, build composition+cmd/server OK (-buildvcs=false), TestERPSource(7)+RejectsInvalidERPSource+RegistersERPImportRoutes+WiresWithoutRepositoryTenant + all existing PASS. Cleaned stray apps/server_core/apps GOCACHE doubling. |
| D15R | review F03-S2 | Reviewer (independent) | claude sonnet | Agent bg ab213bb02827d6d1a | (inline) | (agent return) | done | **0🔴 0🟡 — clean, committed 45481943.** All 8 binary checks PASS: oracle path byte-preserved (source==oracle when empty; sankhya/profitability/oracleDB.Close paths intact), xlsx cache→timing order correct + oracle-only readers left nil (no fabrication/nil-deref), erpSource hard-fail propagates via `return nil, err` (no warn-continue), write path always-registered + repo tenant-free ctor + no double-reg (handler self-classifies), tenant via cfg.DefaultTenantID only, exact sorted governance entry (dashboard→erp_import→integrations), scope 3 files/additive imports/no dep change. |
| D16 | impl F03-S3 | Implementer (standard) | gpt-5.6-luna / high | OS-process | scratchpad/prompt-D16-F03S3.md | scratchpad/agent__D16-F03S3.last.md (+.log); xlsx/oracle substitutability contract test + 2 fixtures (example-erp.xlsx ≥50 valid, identity-rejections.xlsx C02) | **MERGED @ 4da666fb** (retry bbqqymv3b; bln2le4iu died on teardown, 0 files) | FINAL SLICE — F-03 COMPLETE. 3 files: internalread/source_contract_test.go + fixtures/{example-erp(55 rows),identity-rejections(4 rows)}.xlsx. Fixtures hand-rolled archive/zip (NO excelize — grep clean). Hermetic (no tag/DB/server). Drives merged ImportService.RunImport (real parser+ValidateRows, det clock/IDs); independent oracle-shaped fake; C01/C02/C03/C04 pinned + all-rejected→REJECTED + unavailable dual-sentinel. Orch re-verified: gofmt/vet/build(-buildvcs=false)/test both tags GREEN, fixture SHA byte-stable. D16R 2🔴 FIXED pre-commit (see D16R). |
| D16R | review F03-S3 | Reviewer (independent) | claude sonnet | Agent bg a0035cd9ce128211f | (inline) | (agent return) | done | **2🔴 FIXED pre-commit.** Both = test theater on the unsupported loop: oracle-shaped fake hardcoded the *xlsx* error string for GetCurrentPrice/GetSalesHistory, so parity compared xlsx-vs-copy-of-itself. Orch verified against real oracle/reader.go:153-192,281 — Oracle SERVES price (VW_PRECO_TABELA) + sales by default, unsupported only on edge cases; the "both unsupported" premise (my prompt) was factually wrong. FIXED: fake now models Oracle serving data (Amount/Entries populated, nil err); loop rewritten as DIVERGENCE assertions — xlsx typed-unsupported (ADR-17, never zero-as-data) vs Oracle returns data. All else CLEAN: real RunImport (no parallel validator), C04 nil≠0, C03 ean-nil+collision+REFFORN, C02 reject-vs-warn kinds, GTIN 7894900011518 genuinely fails checksum, no excelize, hermetic+deterministic, scope 3 files. Re-verified GREEN. |
| D14R | review F02-S7 | Reviewer (independent) | claude sonnet | Agent bg a7c7e582da9b8bbaf | (inline) | (agent return) | done | 1🔴 1❓. 🔴 GET /erp/imports (list) declared only 200 but handler emits 500 internal_error on ListImports fail (http_handler.go:83)—spec narrower than wire. Verified BROADER: 500 reachable on POST(118/136)+list(83)+detail(100); internal_error∈enum but no path declared 500. FIXED: explicit "500"→ErpImportError on all 3 ops (erp block uses explicit codes, not default→ErrorResponse; flat-error consistent) + regression test (3× "500", no ErrorResponse). ❓ vitest glob = intentional orch fix, coverage maintained. Contract-truth (issue-code/status enum, flat errors, no filename, additive+F01-untouched, index.ts-untouched, test-honesty) all PASS. |

## Binding names / decisions (downstream slices MUST honor)
- Snapshot dedup column = **`file_sha256`** on `erp_import_protocols` (IC-02 `file_hash` reconciled). F02-S4 repo + F02-S5 service use `file_sha256`.
- issues.code CHECK enum (committed 0047) is the CLOSED set of valid codes; F02-S1 domain code constants are a spelled-identical subset.
- protocol DB CHECK = `^#[0-9]{3,}-E$` (≥3 digits, zero-padded, unbounded) — service `#NNN-E` generator must match.
- Migrations are UP-ONLY plain SQL (no down sections) — repo loader idiom (matches 0044).

## Field findings (report to hub)
- codex companion (luna/medium) spawns a FRESH pwsh per grep on Windows → pathologically slow, hang-prone for multi-section investigations. D01 froze 19min on one grep. Mitigation: self-gather via native Grep/Read (fast), reserve codex for OS-process reasoning passes with pre-seeded evidence (no repo tree-crawl).
- F-01 SEMANTIC TENSION (contract vs code): oracle reader ALREADY maps REFERENCIA→manufacturer-ref (ReferenceCode) + hardcodes EAN=nil (reader.go:69-72, catalog_page.go:128). IC-01 requires REFERENCIA→GTIN-validate→ean-or-null, refforn EXCLUSIVELY from REFFORN, +ncm. Truth-order clean (IC-01>code) → F-01 re-maps + flips existing tests (reader_test.go:69-78). Not a BLOCK; scoped into planner.
- F02-S2 REVIEW GAP (schema over-addition): 0045 shipped `filename_sanitized TEXT NOT NULL` — NOT in IC-02 §Persistence, no domain field backs it, nothing reads it; forced D09 to fabricate `""` (ADR-17 violation). Corrected in D09's commit (dropped column; INSERT realigned). IC-02 §"Must Not Decide In Feature Execution" locks columns to contract — the migration author invented one. warning_count (also absent from the contract list) was KEPT: report-backed (ImportReport.WarningCount, GET report exposes warnings). Lesson for hub: migration reviews must diff the CREATE TABLE column set against the interface-contract persistence list, not just constraints.
- F02-S4 correctness (caught pre-review): repo derived rejected_count from rejection-ISSUE count, but ValidateRow emits ≥1 rejection issue per bad field (multi-error row → over-count). Fixed to distinct-rejected-ROW count so accepted+rejected == total rows.
- F01-S4 enum-under-wideness (D10R 4🔴, fixed pre-commit): worker mirrored the planner's "enum exactly [invalid_ean, ean_collision]" literally, but the committed F01-S3 handler (reader.go product() + catalogProductFactResponse) passes reader `QualityFlags` through UNFILTERED — real wire carries `complete` (catalog_page.go:170 default) plus `missing_product`/`ambiguous_product` (oracle reader.go:105/111 on no-match/multi-match). A closed enum narrower than the wire makes conformant clients reject valid 200s. Lesson for hub: an OpenAPI enum on a passthrough field must be derived from the emitter's reachable value set, not from a contract's "intended" subset — or the handler must filter to the declared subset. Fix chose superset (full QualityFlag domain as a shared component) to avoid reopening merged handler logic. If IC-01 truly wants product `quality_flags` scoped to identity-only, that is a follow-up handler-filter decision for the hub (would also let the enum re-narrow).

## Milestone close — P5 VERIFY / P6 DUAL GATE / P7 QA / P8 CLOSE

### P5 aggregate ladder (re-run from clean state @ HEAD db91f385, 2026-07-17)
| Lane | Command | Result |
|---|---|---|
| gofmt | `gofmt -l erp_import internal_read composition` | milestone-authored (erp_import/composition/internalread) CLEAN; pre-existing MIS-002 internal_read files flagged **CRLF-only** (see FINDING F-ENV-M01) |
| vet | `go vet ./erp_import/... ./internal_read/... ./composition/...` | clean |
| build | `go build -buildvcs=false ./cmd/server` | exit 0 (`-buildvcs=false` = worktree VCS-stamp workaround, not a correctness flag) |
| go test | `go test ./erp_import/... ./catalog/... ./internal_read/... ./composition/...` | ALL ok (hermetic, non-integration) |
| integration compile | `go test -tags=integration -run=xxx ./erp_import/... ./internal_read/...` | compiles + skips clean (no DSN) |
| SDK tsc | `tsc --noEmit` | exit 0 |
| SDK vitest | `vitest run --config vitest.config.ts` | **63/63 passed** |

### P6 DUAL GATE (Opus + Sol-medium on fixed-SHA diff 59d0e62f..HEAD)
| Gate | Reviewer | Verdict | Findings |
|---|---|---|---|
| GATE-A | COLD Opus subagent | **PASS** | 0🔴; 1🟡 = listingCostReader over nil oracleDB in xlsx mode → confirm honest-degrade (QA/hub watch) |
| GATE-B | GPT-5.6 Sol / medium (OS-process, scratchpad/prompt-GATE-B-sol.md) | **BLOCK** (4🔴) | R1,R2,R3,R4 (below) |

**Gates DISAGREED (A=PASS, B=BLOCK). Adjudication per finding (dual gate requires concurrence):**
- **R3** (persist tautology): contractRepo asserted `persistedSnapshot` truthiness = compared value to planted copy of itself. **REAL test theater → FIXED.** contractRepo now has explicit `persistCalled bool`; assertion pins the real persist call. Committed **db91f385**.
- **R4** (lock-key mismatch): integration contention test held `hashtextextended($1,0)` but production locks `hashtextextended('erp_import:'||$1,0)` (import_repository.go:29) → no contention created → ErrImportInProgress assertion PASSED VACUOUSLY. **REAL correctness bug → FIXED.** Test now holds the production key (import_repository_test.go:78). Committed **db91f385**.
- **R1** (duplicate error casing): validation-contract.md:37 says re-POST → `409 DUPLICATE_FILE` (UPPERCASE, IC-02 prose), but the ratified wire error family is flat-lowercase (`duplicate_file`, merged F02-S6/S7/D14R + OpenAPI + SDK). **CONTRACT-vs-WIRE conflict — NOT chip-fixable → ESCALATE to hub.** Chip recommendation: lowercase `duplicate_file` STANDS (consistent with the whole ratified erp error family); amend contract prose `DUPLICATE_FILE`→`duplicate_file`. Reversing to uppercase = a wire break across handler+OpenAPI+SDK and breaks family consistency.
- **R2** (Oracle sellable-stock `NVL(est.ESTOQUE,0)-NVL(est.RESERVADO,0)` — ADR-17 "unknown≠zero"): reader.go:467 + catalog_page.go:94. **PRE-EXISTING, UNTOUCHED by M-01** (verified: no +/- hunk in the milestone diff touches these NVL lines; this is MIS-002 code). Out of milestone scope → **DOWNGRADE to hub FINDING** (tracked, not a chip blocker).

**Reconciliation result:** 2 of 4 Gate-B blockers were REAL and are FIXED (R3+R4 @ db91f385, re-verified GREEN). R1 = binding-contract ruling owed by hub. R2 = pre-existing-code finding for hub triage. Gate A's 1🟡 forwarded as a QA watch-item. Dual gate WORKED: it surfaced 2 real bugs Gate A missed + 1 legit contract question.

### P7 QA — DELEGATED to hub (REQUEST)
Chip **structurally cannot** perform the P7 live-drive: C01–C06 require the full dev stack (Go server + Postgres + Oracle creds via `.env`) which is hub-owned; chip never boots a server, binds ports, or loads `.env` (profile §7 + HARNESS-CORE §2). → **REQUEST hub P7 live-drive** vs validation-contract.md C01–C06 (HTTP/API import + DB inspect). No `validation-result.md` authored in-chip (only QA passes a milestone; QA runs hub-side).

### P8 CLOSE artifacts
- Per-feature evidence: F-01/F-02/F-03 `validation.md` (slice→commit→gate map).
- CLOSED event → hub `local_c5532363` (this ledger + feature validation.md are the evidence pointers).

### FINDING F-ENV-M01 (report to hub, ratify into profile §3 false-alarm list)
`gofmt -l` over pre-existing MIS-002 `internal_read/**` .go files reports them "unformatted" — but `gofmt -d` shows the **only** delta is `^M` (CR) stripping (Windows `autocrlf` checked them out CRLF; gofmt normalizes to LF). ZERO token/format changes. Milestone-authored files (erp_import/**, internalread adapter, composition) are LF and gofmt-clean. This is a checkout line-ending artifact, NOT a code defect — a full-tree `gofmt -l` on this Windows worktree is noisy by construction. Mitigation: scope gofmt to milestone-authored dirs, or `git config core.autocrlf false` + re-checkout before a full-tree gofmt gate.

### Hub decisions/requests carried in CLOSED (summary)
| Ref | Type | Ask |
|---|---|---|
| R1 | **DECISION (blocks acceptance)** | Rule duplicate error casing: chip recommends lowercase `duplicate_file` stands + amend contract prose |
| P7 | **REQUEST** | Hub run P7 QA live-drive C01–C06 on local stack (chip cannot boot stack) |
| BARREL-01 | **REQUEST** | Add erpImport re-export to sdk-runtime/src/index.ts (shared-seam, hub-owned barrel) |
| R2 | FINDING | Pre-existing Oracle sellable-stock NVL(...,0) ADR-17 review (untouched by M-01) |
| GATE-A-🟡 | FINDING | listingCostReader over nil oracleDB in xlsx mode — confirm honest-degrade-not-panic |
| GATE-B-Y1 | FINDING | Reader port exposes sellable stock but not physical stock independently when reserved absent — C04 wording "físico segue consultável como físico" (contract:82) may want a port-level physical read; snapshot retains it, port doesn't surface it |
| F-ENV-M01 | FINDING | Windows CRLF gofmt -l false-positive on pre-existing files — ratify into profile §3 |
