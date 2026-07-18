# PLAN-M04.md — Dispatch-Ready Batch Plan (M-04 vinculos-import-ui)

Mission MIS-004-mvp-demo · Milestone M-04 · Branch `chip/m04-vinculos-import-ui` @ base 28b8447c
Features: F-01 product-links-api-gaps (backend/SDK) → F-02 vinculos-screen (FE).
Authoritative contracts: IC-01 **with Amendment A2** (anchor/band/status), IC-05 (FE seams). Ratified locks: SDK-01 (inline index.ts), Finding B/A2 (anchor model), Finding C (FE page-local).
Author: cold-Opus P2 planner subagent (Claude-only lane D-23), 2026-07-18. Inputs: planning-readiness-M04.md, d01/d02 evidence, ic01-amended-A2.md, milestone/feature/validation-contract docs.

## 0. Ratified reconciliations (already closed — implementers do NOT reopen)

- **SDK path:** validation-contract C05 and both feature briefs name `sdk-runtime/src/productLinks.ts`. That grant is **RETIRED**. Per LOCK SDK-01 all new SDK types+methods go **ADDITIVELY INLINE** in `packages/sdk-runtime/src/index.ts` (enriched `ProductLinkCandidateItem` type near :866; new methods in the returned object literal adjacent to the existing product-links block :1508-1545). Zero edits to other domains' lines. No new per-domain file, no barrel touch.
- **Route naming:** feature briefs write `/product-links/resolutions/*`; the live module uses `/product-links/link-resolutions/*` (`http_handler.go:62-64`). **Aligned to the live prefix**: `/product-links/link-resolutions/batch-preview`, `/link-resolutions/batch`, `/link-resolutions/{id}/undo`, `/link-resolutions/batch/{batch_id}/undo`.
- **Migration block:** F-01 brief says 0065-0067, milestone says 0065-0069. Binding = **0065-0069**; this plan uses **0065, 0066, 0067** (0068-0069 reserved, no collision).
- **Undo scope:** F-01 EARS specifies single-resolution undo; milestone C02 additionally requires **batch undo (all items of a lote)**. S4 delivers BOTH endpoints. Undo is derivable from the audit chain (`previous_state`/`previous_internal_product_id` already persisted) + batch linkage — no schema for undo state itself.

## 1. DAG (dependency graph + parallelism)

```
F-01 (backend, serial spine):
  S1 ─┐
  S2 ─┼─────────► S5 (OpenAPI + SDK inline) ══F-01→F-02 GATE══╗
  S3 ─┤                                                        ║
  S4 ─┘                                                        ║
                                                               ▼
F-02 (frontend, after S5 committed & reviewed):
  S6 ──► S7 ──► S8
                 └─ S9 (parallel with S7/S8; only needs S6 scaffold)
```

- **S1** (candidate enrichment) is independent and can start immediately; it is the truth-source for S5's candidate schema.
- **S2, S3, S4** share the resolution/batch surface. Serialize **S3 → S4** (undo reverts what apply wrote; S4 reuses S3's batch linkage). **S2** (preview, read-only dry-run) is independent of S3/S4 and parallelizable with them, but both S2 and S3 share the same per-item eval helper — implement that helper in S2 and reuse in S3 (so **S2 before S3** is preferred, not strictly required).
- **S5** is the integration slice: it serializes AFTER S1-S4 land their Go handlers/domain shapes, because OpenAPI+SDK must mirror the final Go JSON shapes. S5 commit+review = the **F-01→F-02 gate**.
- **F-02**: **S6 → S7 → S8** serial (scaffold → rows/drawer → bulk). **S9** (Importação section) depends only on S6 scaffold and is parallelizable with S7/S8.
- Parallelizable set at start: {S1} ∥ {S2}. After S2: {S3}→{S4}. All F-02 slices blocked on S5.
- NOTE (single-worktree collision): all F-01 slices share `http_handler.go`, `workflow_store.go`, `link_candidate_repo.go`, `runner_test.go`; all F-02 slices share `pages/vinculos/**`. Implementers run in ONE chip worktree with no per-slice isolation ⇒ orchestrator SERIALIZES implementer dispatch even where the logical DAG allows parallelism, to avoid same-file races.

## 2. Migration ledger (block 0065-0069, all additive)

| Mig | Owner | Operation | Fixture bump |
|-----|-------|-----------|--------------|
| **0065** | S1 | `ALTER TABLE product_link_candidates ADD COLUMN confidence integer NOT NULL DEFAULT 0, ADD COLUMN confidence_band text NOT NULL DEFAULT '', ADD COLUMN match_status text NOT NULL DEFAULT '', ADD COLUMN reasons jsonb NOT NULL DEFAULT '[]'::jsonb` (additive; preserves `tenant_id` PK). Backfill = next generation only (not retroactive). | 44 → 45 |
| **0066** | S3 | `CREATE TABLE IF NOT EXISTS product_link_batches (tenant_id text NOT NULL, batch_id text NOT NULL, installation_id text NOT NULL, actor_type/actor_id/actor_name text, requested_count/applied_count/failed_count integer NOT NULL DEFAULT 0, status text NOT NULL, created_at timestamptz NOT NULL DEFAULT now(), PRIMARY KEY (tenant_id, batch_id))` — the module-owned batch audit table (C02 "tabela própria"). | 45 → 46 |
| **0067** | S3 | `ALTER TABLE product_link_audit_entries ADD COLUMN batch_id text NOT NULL DEFAULT ''` (links per-item audit rows to their batch; enables batch-undo grouping). | 46 → 47 |
| 0068-0069 | — | Reserved, unused. No file, no collision. | — |

Every new `.sql` increments the two hardcoded `44` literals in `apps/server_core/internal/platform/migrate/runner_test.go:25` and `:64`. Final count after F-01 = **47**. S4 (undo) needs **no migration** (undo state derived from the audit chain: latest-audit-for-identity + action=`undo` + `batch_id`).

## 3. Anchor/confidence model (IC-01 A2, deterministic — the truth for S1)

Cross-side anchors that exist provider-side: **`seller_sku`** (exact→codprod, ACCEPT-grade in M-04 scope) and **`ean`** (corroboration, `unproved` flag). `title` = ranking-only, never an anchor. `marca`/`refforn` exist only internal-side ⇒ **always emitted in `reasons[]` as `UNAVAILABLE`** (ADR-17, never silently dropped).

Deterministic band/status/score mapping (thresholds ≥85 ALTA / 50-84 MEDIA / <50 BAIXA, unchanged):

| Situation | confidence | confidence_band | match_status | reasons[] anchors |
|-----------|-----------|-----------------|--------------|-------------------|
| seller_sku + ean → same codprod, no hard-neg | 95 | ALTA | ACCEPT | seller_sku FOR, ean FOR, marca UNAVAILABLE, refforn UNAVAILABLE |
| seller_sku alone (ean absent or no ean-match) | 70 | MEDIA | REVIEW | seller_sku FOR, ean UNAVAILABLE (EAN ausente ⇒ máx REVIEW), marca/refforn UNAVAILABLE |
| ean alone (sku no match) | 60 | MEDIA | REVIEW | ean FOR (unproved), seller_sku UNAVAILABLE, marca/refforn UNAVAILABLE |
| title-only match | 35 | BAIXA | REVIEW | title FOR (ranking-only), seller_sku/ean/marca/refforn UNAVAILABLE |
| SKU/EAN conflict (sku→A, ean→B) | 20 | BAIXA | REJECT | per-candidate: FOR on its own anchor + AGAINST citing the conflicting anchor |
| hard-negative title (kit/combo/cor/medida/voltagem) even w/ EAN match | 25 | BAIXA | REJECT | ean FOR, title AGAINST (contradição vence EAN) |
| no product matched | 0 | BAIXA | NO_CANDIDATE | seller_sku/ean/marca/refforn UNAVAILABLE |

**EXACT enum vocabulary that must appear byte-identical in migration values, Go consts, OpenAPI enum, SDK type, and UI:** `match_status`: `ACCEPT|REVIEW|REJECT|NO_CANDIDATE`; `confidence_band`: `ALTA|MEDIA|BAIXA`; reason `direction`: `FOR|AGAINST|UNAVAILABLE`. Error codes: `INSUFFICIENT_ANCHORS` (422), `PRODUCT_NOT_FOUND` (404) from IC-01; plus F-01 batch/undo codes `ALREADY_RESOLVED` (409), `ALREADY_UNDONE` (409), `SUPERSEDED` (409), empty-batch (422).

## 4. Fixture plan for S1 (≥8 IC-01-A2 cases — anchor model made concrete)

Fixtures are F-01-owned (M-02 shared fixtures do NOT exist on base 28b8447c). Each = one listing snapshot (`seller_sku`, `ean`, `title`) + matcher stub returns → assert `confidence_band` + `match_status` + key `reasons[]`.

1. **CONCORDANT-ALTA** — sku `MLB-SKU-1`→codprod 100, ean `7891234567895`→codprod 100, titles agree. ⇒ `confidence≈95`, **ALTA**, **ACCEPT**; reasons: seller_sku FOR, ean FOR, marca UNAVAILABLE, refforn UNAVAILABLE. *(also the primary auto-ACCEPT proxy)*
2. **SKU-ALONE-MEDIA** — sku `MLB-SKU-2`→codprod 200; snapshot `ean=""`. ⇒ `confidence≈70`, **MEDIA**, **REVIEW**; reasons: seller_sku FOR, ean UNAVAILABLE detail "EAN ausente ⇒ máximo REVIEW", marca/refforn UNAVAILABLE. *(EAN-absent⇒REVIEW binding)*
3. **EAN-ALONE-MEDIA** — sku no match; ean `7890000000017`→codprod 300. ⇒ `confidence≈60`, **MEDIA**, **REVIEW**; reasons: ean FOR (unproved), seller_sku UNAVAILABLE, marca/refforn UNAVAILABLE.
4. **TITLE-ONLY-BAIXA** — sku/ean no match; title LIKE→codprod 400. ⇒ `confidence≈35`, **BAIXA**, **REVIEW**; reasons: title FOR (ranking-only), all anchors UNAVAILABLE.
5. **SKU-EAN-CONFLICT-REJECT** — sku `MLB-SKU-5`→codprod 500, ean `7899999999994`→codprod 501 (disagree). ⇒ two candidates, each **BAIXA/REJECT**; candidate-500: seller_sku FOR + ean AGAINST ("EAN aponta codprod 501"); candidate-501: ean FOR + seller_sku AGAINST.
6. **DOKA-HARD-NEGATIVE-REJECT** — ean `7896541230004`→codprod 600 (matches), but snapshot title `Andaime Doka Kit 12 peças` vs internal `Escora Menegotti unidade` (kit/combo + brand collision). ⇒ `confidence≈25`, **BAIXA**, **REJECT**; reasons: ean FOR, **title AGAINST** detail "hard-negative: kit/combo divergente (colisão Doka)". *(contradição vence EAN — mandatory Doka-style case)*
7. **VOLTAGE-HARD-NEGATIVE-REJECT** — sku→codprod 700 and ean→codprod 700 (would be ALTA) but title `Furadeira 220V` vs internal `Furadeira 110V`. ⇒ cap **BAIXA**, **REJECT**; reasons: seller_sku FOR, ean FOR, title AGAINST detail "voltagem divergente 220V≠110V".
8. **NO-CANDIDATE** — sku/ean/title all no match (unresolved). ⇒ `confidence=0`, **NO_CANDIDATE**; reasons: seller_sku/ean/marca/refforn all UNAVAILABLE (honest, non-empty).
9. **MARCA-REFFORN-UNAVAILABLE (explicit)** — any matched case; assert `reasons[]` ALWAYS contains `{anchor:"marca",direction:"UNAVAILABLE"}` and `{anchor:"refforn",direction:"UNAVAILABLE"}` regardless of band (ADR-17 motivo sempre visível). Bind to case 1's payload with a dedicated assertion.

---

## 5. Slice cards

### S1 — Candidate enrichment: confidence/band/status/reasons — F-01 — COMPLEX
**Why complex:** scoring state machine over the existing generation ladder + lexical **hard-negative detector** (kit/combo/cor/medida/voltagem) comparing snapshot title vs internal product name; deterministic band/status assignment must match A2 exactly.
- **write-set:**
  - `apps/server_core/migrations/0065_product_link_candidate_confidence.sql` (ALTER per ledger).
  - `apps/server_core/internal/platform/migrate/runner_test.go` (44→45 both literals).
  - `apps/server_core/internal/modules/product_links/domain/link_candidate.go` (add `Confidence int`, `ConfidenceBand`, `MatchStatus`, `Reasons []LinkCandidateReason` + new const enums `ALTA|MEDIA|BAIXA`, `ACCEPT|REVIEW|REJECT|NO_CANDIDATE`, `FOR|AGAINST|UNAVAILABLE`; new `LinkCandidateReason{Anchor,Direction,Detail}`).
  - `apps/server_core/internal/modules/product_links/application/generation_service.go` (compute confidence/band/status/reasons in `buildExactCandidates`/`buildConflictCandidates`/`buildCandidatesFromProducts`/`newCandidate`; add hard-negative helper).
  - `apps/server_core/internal/modules/product_links/application/generation_service_test.go` (**new** — the ≥8 fixtures of §4).
  - `apps/server_core/internal/modules/product_links/adapters/postgres/link_candidate_repo.go` (persist/scan new columns in `ReplaceLinkCandidates` + list/get; marshal `reasons` jsonb).
- **depends-on:** none. First to start.
- **contract-satisfiability:** IC-01 A2 anchor model (§3 table); emits `confidence_band ALTA|MEDIA|BAIXA` + `match_status ACCEPT|REVIEW|REJECT|NO_CANDIDATE` + `reasons[{anchor,direction,detail}]` with marca/refforn always UNAVAILABLE. `seller_sku`+`ean` concordant ⇒ ALTA/ACCEPT; single anchor ⇒ MEDIA/REVIEW; EAN-absent ⇒ max REVIEW; hard-neg/conflict ⇒ BAIXA/REJECT AGAINST.
- **verification map:** advances **C03** (fila presentation source) + is the data half of **C05**. Rungs **L0/L1/L2**. Evidence: `generation_service_test.go` 9 fixtures pass; `go test ./internal/modules/product_links/...` green; migration applies; runner_test 45.
- **additive-lock/notes:** 0065 additive ALTER; no barrel/shared-seam touch.

### S2 — Batch-preview dry-run (itemized, zero persist) — F-01 — STANDARD
- **write-set:**
  - `application/batch_service.go` (**new**) — `PreviewBatch(ctx, {Approvals:[{CandidateID}]}) → {Items:[{CandidateID, Status OK|FAILED, Cause?}]}`; per-item eval helper `evaluateApproval(candidateID)` reused by S3 (checks candidate exists, resolvable, not already ALREADY_RESOLVED). Reads only; **no store writes**.
  - `transport/http_handler.go` (register `POST /product-links/link-resolutions/batch-preview`; new `handleBatchPreview`; empty batch ⇒ 422; decode error ⇒ 400).
  - `application/batch_service_test.go` (**new**) — 3 approvals (1 invalid) ⇒ 2 OK + 1 FAILED, nothing persisted (assert store untouched).
- **depends-on:** S1 (candidate shape). Parallel-safe with S3/S4 but implement the shared `evaluateApproval` helper here.
- **contract-satisfiability:** F-01 preview 200 itemized, partial-fail not HTTP error; item states `OK|FAILED` + `cause`. Honors C01 "SELECT antes/depois idêntico".
- **verification map:** **C01**. Rungs L0/L1/L2. Evidence: unit test asserts zero repo mutation; dry-run response itemized. (Endpoint-level SELECT-before/after proof lands at QA via S5+live stack.)
- **additive-lock/notes:** no migration; additive route.

### S3 — Batch apply (per-item, partial-failure, batch audit) — F-01 — COMPLEX
**Why complex:** partial-failure semantics (one item's failure must not roll back others), reuse of `ApproveCandidate` transition logic per item, and a new batch audit aggregate row spanning N per-item transitions.
- **write-set:**
  - `apps/server_core/migrations/0066_product_link_batches.sql` (CREATE table) + `0067_product_link_audit_batch_id.sql` (ALTER audit ADD `batch_id`).
  - `runner_test.go` (45→47).
  - `application/batch_service.go` — `ApplyBatch(ctx,{Approvals,Actor}) → {BatchID, Applied:[...], Failed:[{CandidateID,Cause}]}`; each item calls the SAME resolution path as `ResolutionService.ApproveCandidate` (inject resolver or share `buildTransition`); item failure itemized (`ALREADY_RESOLVED`→FAILED, continue); writes one `product_link_batches` row (`operator` literal actor, counts, status) + tags each item's audit entry with `batch_id`.
  - `ports/workflow_store.go` (+ `InsertBatch`, and thread `batch_id` through `ApplyProductLinkTransition` — additive method/param).
  - `adapters/postgres/link_candidate_repo.go` (insert batch row; write `batch_id` on audit insert at :288-315).
  - `transport/http_handler.go` (register `POST /product-links/link-resolutions/batch`; `handleBatchApply`; empty ⇒ 422).
  - `application/batch_service_test.go` + repo integration test (2 applied + 1 failed; batch row + per-item audit inspectable; queue no longer returns the 2).
- **depends-on:** S1, S2 (eval helper). Serial before S4.
- **contract-satisfiability:** F-01 apply 200 with `applied`/`failed`; module-owned batch audit table (C02 "tabela própria"); actor `operator`; **zero `/mutations` envelope, zero ML write by construction** (batch calls only local resolver). `INSUFFICIENT_ANCHORS`/manual paths untouched.
- **verification map:** **C02** (apply + per-item audit + batch row). Rungs L0/L1/L2. Evidence: transcript 2 applied/1 failed; `SELECT product_link_batches`, `product_link_audit_entries WHERE batch_id=...`; connectors mock untouched (local-only proof).
- **additive-lock/notes:** 0066 CREATE + 0067 additive ALTER; workflow_store port extended additively (still inside `modules/product_links/**`).

### S4 — Undo (single resolution + batch) — F-01 — COMPLEX
**Why complex:** reversal must reconstruct prior link state from the audit chain; `ALREADY_UNDONE`/`SUPERSEDED` are ordering-sensitive over concurrent resolutions; batch-undo fans out over all `batch_id` items atomically-per-item.
- **write-set:**
  - `domain/product_link.go` (add `ProductLinkActionUndo ProductLinkAction = "undo"`).
  - `application/resolution_service.go` (add `UndoResolution(ctx, auditID)` + `UndoBatch(ctx, batchID)`): load target audit; if not latest-for-identity ⇒ `SUPERSEDED` 409; if latest action already `undo` ⇒ `ALREADY_UNDONE` 409; if none ⇒ `PRODUCT_LINKS_..._NOT_FOUND` 404; else write reversal transition (`state`←`previous_state`, product←`previous_internal_product_id`, action=`undo`). Batch-undo iterates all audit rows of `batch_id`.
  - `ports/workflow_store.go` (+ `GetAuditEntry(auditID)`, `LatestAuditForIdentity(identity)`, `ListAuditByBatch(batchID)` — additive).
  - `adapters/postgres/link_candidate_repo.go` (implement the three read queries, all `tenant_id`-scoped).
  - `transport/http_handler.go` (register `POST /product-links/link-resolutions/{id}/undo` + `POST /product-links/link-resolutions/batch/{batch_id}/undo`; path-param parse; map `SUPERSEDED`/`ALREADY_UNDONE`→409).
  - `application/resolution_service_test.go` (undo ⇒ candidate back in queue; second undo ⇒ 409 ALREADY_UNDONE; undo after newer resolution ⇒ 409 SUPERSEDED; batch-undo reverts ALL).
- **depends-on:** S3 (batch_id + apply). Serial after S3.
- **contract-satisfiability:** F-01 undo 404/409 ALREADY_UNDONE/409 SUPERSEDED; candidate returns to fila; history preserved (reversal is a new audit row, not a delete). Batch-undo reverts all items + audits reversal (C02).
- **verification map:** **C02** (undo path). Rungs L0/L1/L2. Evidence: transcripts above; audit trail shows reversal rows.
- **additive-lock/notes:** no migration (derives from audit chain); additive routes + additive port methods.

### S5 — OpenAPI additive + SDK inline (SDK-01) — F-01 — STANDARD — **F-01→F-02 GATE**
- **write-set:**
  - `contracts/api/marketplace-central.openapi.yaml`:
    - `ProductLinkCandidate` schema (:5057-5103) — ADD `confidence` (int 0-100), `confidence_band` enum `[ALTA,MEDIA,BAIXA]`, `match_status` enum `[ACCEPT,REVIEW,REJECT,NO_CANDIDATE]`, `reasons` array of `{anchor, direction enum [FOR,AGAINST,UNAVAILABLE], detail}` (all additive properties).
    - New paths under `/product-links/link-resolutions/*`: `batch-preview`, `batch`, `{id}/undo`, `batch/{batch_id}/undo` + request/response schemas (`BatchPreviewRequest`, `BatchPreviewResult`, `BatchApplyRequest`, `BatchApplyResult`, undo responses) inserted in the product-links region (near :1446-1496 paths, :5288+ schemas). Error responses `422/404/409`.
  - `packages/sdk-runtime/src/index.ts` (**SDK-01 additive block ONLY**): enrich `ProductLinkCandidateItem` type (~:866) with the 4 fields + `ProductLinkReason` type; add methods `previewProductLinkBatch`, `applyProductLinkBatch`, `undoProductLinkResolution`, `undoProductLinkBatch` inline in the returned object literal adjacent to :1508-1545 following the `postJson<T>` pattern. No edits to other domains' lines; **no barrel line**.
- **depends-on:** S1-S4 (mirrors final Go JSON shapes). **Its committed+reviewed state is the gate that unblocks all F-02 slices.**
- **contract-satisfiability:** IC-01 exact enum vocabulary surfaced in API+SDK unchanged; additive-only `/product-links/*`. IC-05 SDK-via-ClientContext contract preserved.
- **verification map:** **C05** (seams) + enables C01/C02/C03. Rungs **L0/L1/L3** (`tsc --noEmit` in sdk-runtime; OpenAPI lint/type-sync). Evidence: `tsc` clean; spec diff confined to `/product-links/*`; SDK diff confined to the product-links block.
- **additive-lock/notes:** SDK-01 inline lock consumed here; **flag the diff region in the CLOSED payload** (same mechanism as BARREL-01/D-13). No `packages/web-query` / `packages/ui` touch.

### S6 — /vinculos scaffold + tabs + KPIs — F-02 — STANDARD
- **write-set:**
  - `apps/web/src/routes/vinculos.tsx` (replace `WorkspacePlaceholder` with `<VinculosPage/>`; AppRouter untouched — IC-05).
  - `apps/web/src/pages/vinculos/VinculosPage.tsx` (**new**) — Fila/Resolvidos tabs (page-local `role="tab"` per AnunciosPage precedent), top KPIs (pendentes / alta confiança / resolvidos hoje), `LoadingState`/`ErrorState`/`EmptyState` from `@marketplace-central/ui`, `useClient`+`InstallationContext`.
  - `apps/web/src/pages/vinculos/vinculosQueryKeys.ts` (**new, page-local**) — `['product-links','queue',filters]`, `['product-links','resolved']` (Finding C: page-local, avoids web-query seam).
- **depends-on:** S5 committed+reviewed (gate); M-03 seams (present).
- **contract-satisfiability:** IC-05 `routes/vinculos.tsx` seam + `pages/vinculos/**` location + TanStack Query via web-query client; obligatory state components. Vínculos stays out of global nav (deep-link reachable).
- **verification map:** **C03** (scaffold). Rungs **L0/L3** (web build/tsc). Evidence: `/vinculos` renders tabs; placeholder-swap leaves AppRouter.tsx clean-diff.
- **additive-lock/notes:** page-local query keys (no web-query edit). If a shared `linkage` staleTime proves unavoidable → **REQUEST flag** (see open_questions gate — currently avoided).

### S7 — Queue rows + drawer + individual resolve — F-02 — STANDARD
- **write-set (all under `apps/web/src/pages/vinculos/`, new):** `QueueTab.tsx`, `QueueRow.tsx` (produto CODPROD+descrição+EAN/refforn, best candidate título+preço+banda%, anchor chips FOR green/AGAINST red/UNAVAILABLE gray — **motivo always rendered with %**), `VinculoDrawer.tsx` (field-by-field compare, all ranked candidates; wraps `DetailPanel` per ListingDetailPanel precedent; open state via `?candidate=` URL for deep-link/F5), `useVinculosQueue.ts` (queue query + approve via existing `approveProductLinkCandidate`, reject via `rejectProductLinkListing`, manual via `manualResolveProductLink`; invalidate queue+resolved+KPIs).
- **depends-on:** S6.
- **contract-satisfiability:** IC-01 presentation fidelity (motivo sempre visível, % nunca sozinho; band chips ALTA/MEDIA/BAIXA); `NO_CANDIDATE` honest state (not empty row); IC-05 drawer 300-380px, no inline table edit, SDK-only (no raw fetch). 409 `ALREADY_RESOLVED` ⇒ conflict banner + queue refetch.
- **verification map:** **C03**. Rungs L0/L3. Evidence: screenshot ≥3 rows distinct bands + anchor chips; drawer open; approve ⇒ 2xx + item leaves queue after refetch; deep-link `?candidate=X`+F5 reopens drawer.
- **additive-lock/notes:** page-local; reuses `invalidateAfterMutation('link_apply')` (already maps `linkage` — invalidation.ts:26).

### S8 — Bulk select → preview modal → apply → feedback — F-02 — COMPLEX
**Why complex:** two-phase dry-run→apply UX with partial-failure rendering, local selection state lifecycle, and post-apply refetch orchestration (no optimistic list mutation).
- **write-set (new under `pages/vinculos/`):** `BulkBar.tsx` (local multi-select state, cleared after enqueue), `BatchPreviewModal.tsx` (page-local `role="dialog"` per MutationPreviewModal precedent; renders dry-run items + predicted failures; "prosseguir só com válidos"), `BatchResultFeedback.tsx` (applied/failed + link to Resolvidos — page-local inline banner since **no toast primitive exists**), `useVinculosBatch.ts` (`previewProductLinkBatch` → confirm → `applyProductLinkBatch` → invalidate queue+resolved+KPIs → refetch).
- **depends-on:** S7.
- **contract-satisfiability:** C01/C02 consumed via UI — preview shows dry-run (nothing applied), apply shows applied/failed itemized, fila refetch post-2xx (no optimistic write per State model). IC-05 SDK-only.
- **verification map:** **C03** (preview→apply→Resolvidos flow) + live proof of **C01/C02**. Rungs L0/L3 + QA live-drive. Evidence: bulk 3 items ⇒ modal dry-run ⇒ apply ⇒ applied/failed feedback ⇒ fila reduced + items in Resolvidos, no manual F5.
- **additive-lock/notes:** page-local modal/toast (Finding C — no packages/ui change).

### S9 — Importação section (protocol read-only) — F-02 — STANDARD
- **write-set:** `apps/web/src/pages/vinculos/ImportacaoSection.tsx` (**new**) — reads `GET /erp/imports` (+ `/erp/imports/{id}` for rejection report) via SDK; shows protocolo `#NNN-E`, status `COMPLETED|REJECTED`, accepted/rejected counts, expandable per-row rejection report. `useErpImports.ts` (**new**, page-local). **No upload UI** (Non-Scope).
- **depends-on:** S6 scaffold only (parallel with S7/S8).
- **contract-satisfiability:** C04 — protocol + status + counts + rejection report from `/erp/imports`; ADR-17 honest states (rejections visible, never hidden).
- **verification map:** **C04**. Rungs L0/L3 + live-drive. Evidence: screenshot Importação section with last import protocol + expanded rejections.
- **additive-lock/notes:** consumes existing ERP SDK types (barrel `erpImport` types already exported); read-only; no seam touch.

## 6. Per-criterion verification rollup

| Criterion | Advanced by | Concrete proof |
|-----------|-------------|----------------|
| **C01** batch preview dry-run | S2 (+S5 wire, S8 UI) | `batch_service_test.go` zero-persist; live SELECT before/after identical; no `/mutations` protocol |
| **C02** batch apply + audit + undo | S3, S4 (+S8 UI) | 2 applied/1 failed transcript; `product_link_batches`+audit `batch_id` SELECT; undo reverts + reversal audit; connectors mock untouched |
| **C03** tela flow preview→apply→Resolvidos | S1, S6, S7, S8 | live-drive: bands+anchor chips, NO_CANDIDATE honest, apply refetch without F5, deep-link+F5 |
| **C04** import protocol visible | S9 | screenshot protocolo `#NNN-E`+status+counts+rejection rows |
| **C05** migrations + seams | S1, S3, S5 | `ls migrations | grep 006[5-9]` = 0065/0066/0067; runner_test 47; diff confined to owned paths; `tsc` clean |

## 7. Ownership self-audit (write-set ⊆ ALLOWED)

All write paths ∈ {`modules/product_links/**`, `apps/server_core/migrations/0065-0067`, `runner_test.go` fixture bump, OpenAPI `/product-links/*` additive, `index.ts` SDK-01 block, `apps/web/src/routes/vinculos.tsx`, `apps/web/src/pages/vinculos/**`, page-local query keys}. **No** touch to `modules/mutations/**`, `modules/market/**`, `connectors/**`, other domains' index.ts lines, SDK barrel, `packages/ui`, `packages/web-query`, `apps/web/src/app/**`, other routes, `.env*`, dev stack. Zero ML writes.

## 8. open_questions

**(empty — plan is dispatch-ready.)**

---

### Critical Files for Implementation
- `apps/server_core/internal/modules/product_links/application/generation_service.go` (S1 scoring/anchor model)
- `apps/server_core/internal/modules/product_links/application/resolution_service.go` (S3/S4 apply + undo reuse)
- `apps/server_core/internal/modules/product_links/transport/http_handler.go` (S2/S3/S4/S5 additive routes + error mapping)
- `contracts/api/marketplace-central.openapi.yaml` (S5 additive `/product-links/*` + enriched `ProductLinkCandidate`, region :1364-1496 / :5057-5330)
- `packages/sdk-runtime/src/index.ts` (S5 SDK-01 inline block: types ~:866, methods ~:1508-1545)
