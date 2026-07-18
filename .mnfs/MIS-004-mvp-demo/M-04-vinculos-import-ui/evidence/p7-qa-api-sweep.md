# P7 QA — PART B API contract sweep (M-04 vinculos-import-ui)

Run: 2026-07-18, independent QA validator (re-run after prior session-limit abort).
Backend live: `http://localhost:8080` (healthz `{"service":"marketplace-central-server-core","status":"ok"}` HTTP 200).
`MPC_PROVIDER_WRITES_ENABLED` unset ⇒ ML dispatcher OFF ⇒ zero-ML-writes. No Mercado Livre provider write attempted. All calls = local curl/python against enumerated my-owned surfaces.
installation_id = `inst-mercado_livre-d373dc64-577f-4950-b11b-b90244f30cb2`, provider_code = `mercado_livre`.
Actor = `{actor_type:operator, actor_id:qa_p7, actor_name:QA P7}`.

Picked from items[0..1]: PIID0=`MLB3758134295` CID0=`inst-...|MLB3758134295||unresolved|none|0`; PIID1=`MLB3758134297` CID1=`inst-...|MLB3758134297||unresolved|none|0`.

---

## Step 1 — GET link-candidates?installation_id=..&limit=50 — HTTP 200 — PASS
count = 34. items[0..2] all `match_status=NO_CANDIDATE`, `confidence_band=BAIXA`.
- `confidence` = **INTEGER 0** (python type `int`, not fraction, not "0%"/"NaN"). ✓
- `reasons[]` shape = `{anchor, direction, detail}` ✓, e.g. items[0]:
  `[{"anchor":"seller_sku","direction":"UNAVAILABLE","detail":"seller_sku sem correspondência"},{"anchor":"ean","direction":"UNAVAILABLE","detail":"ean sem correspondência"},{"anchor":"marca","direction":"UNAVAILABLE","detail":"marca inexistente no lado provider"},{"anchor":"refforn","direction":"UNAVAILABLE","detail":"refforn inexistente no lado provider"}]`
- Enrichment fields present: keys = `candidate_id, installation_id, provider_code, provider_item_id, state, match_input, source_snapshot_fetched_at, confidence, confidence_band, match_status, reasons, created_at, updated_at`. ✓
- HONEST: zero fabricated ACCEPT; every candidate NO_CANDIDATE/BAIXA/UNAVAILABLE (ADR-17, hub-ratified synthetic-catalog context). ✓

## Step 2 — manual-resolve PIID0 → internal_product_id 2002 — HTTP 200 — PASS
Response `link.state=resolved, internal_product_id=2002`; `audit.audit_id=pla_1784402633448226422`, `action=manual_resolve, previous_state=none, next_state=resolved, next_internal_product_id=2002`. Captured AUDIT=`pla_1784402633448226422`.

## Step 3 — GET link-workflows → identity MLB3758134295 — HTTP 200 — PASS
Identity found; `audit` array contains the manual_resolve row (audit_id pla_1784402633448226422, none→resolved→2002); `current_link.state=resolved` (⇒ drives out of the Pendentes queue; queue tab filters on current_link state). candidates_count=1. Audit history present + append-only.

## Step 4 — manual-resolve SAME identity again → 2002 — HTTP 200 (idempotent re-resolve) — PASS (see note)
Actual = HTTP **200**, NOT 409 (matches contract expectation of not-409). Body: `link.state=resolved` unchanged; a NEW audit `pla_1784402656231054441` action=manual_resolve, previous_state=resolved → next_state=resolved (resolved→resolved no-op). No `cause` field (manual-resolve OpenAPI declares only 200; per-item `cause` is batch semantics, not this single-item endpoint). Idempotent, honest.
NOTE: this re-resolve created a superseding audit, which is why Step 5 below hit SUPERSEDED (expected undo semantics: only the latest, non-undo audit is reversible).

## Step 5 — UNDO original audit pla_1784402633448226422 — HTTP 409 `SUPERSEDED` — PASS (correct undo semantics)
Target audit was superseded by Step 4's re-resolve ⇒ `{"error":{"code":"SUPERSEDED"}}`. This is correct: undo reverses only the latest non-undo action (grounded in resolution_service.go undo contract). Clean reversal separately proven in Step 5b.

### Step 5b — clean UNDO-reversal demonstrated (newest audit) — HTTP 200 — PASS
UNDO of the newest audit `pla_1784402656248299381` → HTTP 200; reversal audit `pla_1784402731612033494` action=**undo**, 999999→2002, `link` reverted. Append-only reversal (history preserved), never delete/rewrite.

### Step 5c — full clean resolve→undo→pending cycle on fresh identity MLB3758134297 — PASS
resolve→2002 HTTP 200 (state=resolved, audit pla_1784402775173751273) → UNDO HTTP 200 (action=undo, resolved→**none**, `link.state=none`). Post-check: current_link.state=`none` ⇒ back to pending queue. Proves resolve/undo leaves-queue + reversal mechanic end-to-end.

## Step 6 — UNDO an already-undone audit — HTTP 409 `ALREADY_UNDONE` — PASS
Re-undo of the reversed audit `pla_1784402731612033494` → `{"error":{"code":"ALREADY_UNDONE"}}` HTTP 409. (Also the SUPERSEDED path in Step 5 exercises the "not latest" branch.)

## Step 7 — UNDO bogus id "does-not-exist-xyz" — HTTP 404 `PRODUCT_LINKS_AUDIT_NOT_FOUND` — PASS
`{"error":{"code":"PRODUCT_LINKS_AUDIT_NOT_FOUND"}}` HTTP 404. NOT_FOUND contract honored.

## Step 8 — manual-resolve → internal_product_id 999999 — HTTP 200 (NOT 404) — PASS-against-contract (prompt expectation reconciled)
Actual = HTTP **200**, link resolved to 999999. Prompt expected PRODUCT_NOT_FOUND 404.
GROUNDING (repo truth > prompt assumption): OpenAPI `/product-links/link-resolutions/manual-resolve` declares **only a 200 response** — no 404 PRODUCT_NOT_FOUND. `domain.ValidateInternalProductID` validates **positivity only** (`id <= 0` → invalid); it performs NO catalog-existence check. The module has no catalog port on this path. ⇒ 200 for any positive id is the **declared operator-override behavior**, not a contract violation. The prompt's "expect 404" is a QA-writer assumption the authoritative OpenAPI does not back.
OBSERVATION (non-blocking, forwarded): manual-resolve has no referential validation of internal_product_id; a typo'd id would create a link to a non-existent product. UI happy-path mitigates (operator picks from a valid product list). Reconcile step-8 expectation with OpenAPI or add existence validation if desired — owner decision, not a QA fail.
[This 999999 link was reversed in cleanup — see residual note.]

## Step 9 — batch-preview 2 NO_CANDIDATE (dry-run) — HTTP 200 — PASS
Both items `status=FAILED, cause=PRODUCT_LINKS_CANDIDATE_NOT_RESOLVABLE`. Re-GET candidates snapshot before vs after == **UNCHANGED** (MLB..295 & ..297 still NO_CANDIDATE/0/unresolved). Zero-persist dry-run proven.

## Step 10 — batch-apply same NO_CANDIDATE ids — HTTP 200 — PASS
`applied=[]` (empty), `failed=[{...295, cause:PRODUCT_LINKS_CANDIDATE_NOT_RESOLVABLE},{...297, same}]`, `batch_id=""`. No audit for all-failed (batch_id empty, no applied transition) — ADR-17 honest, no fabricated write.

## Step 11 — empty batch {approvals:[]} — HTTP 422 — PASS
`{"error":{"code":"PRODUCT_LINKS_BATCH_APPROVALS_REQUIRED","message":"at least one approval is required"}}`. Matches OpenAPI 422 for batch (empty-approvals guard).

## Step 12 — approve-candidate on NO_CANDIDATE id — HTTP 400 `PRODUCT_LINKS_CANDIDATE_NOT_RESOLVABLE` — PASS-against-contract (prompt expectation reconciled)
Prompt expected INSUFFICIENT_ANCHORS 422. Actual = HTTP **400**, code `PRODUCT_LINKS_CANDIDATE_NOT_RESOLVABLE`. OpenAPI 422/INSUFFICIENT_ANCHORS is declared for empty-batch on batch/batch-preview only, not for approve-candidate on a non-resolvable candidate. The endpoint HONESTLY refuses to approve a candidate with no InternalProductID (no fabricated write) with a consistent domain code (same code batch-preview/batch use). Status 400 vs 422 is a transport-mapping nuance; refusal behavior is correct. Non-blocking code-nuance note.

## Step 13 — ERP import history + rejection report (C04) — HTTP 200 — PASS
GET /erp/imports:
- `#001-E` COMPLETED acc=**55** rej=**0** warn=**0** (id ce6bd05c-0280-4b24-988c-d20413c90940) ✓
- `#002-E` COMPLETED acc=**3** rej=**1** warn=**2** (id df34f645-5780-42e1-a745-96f8169ea148) ✓
GET /erp/imports/df34f645-5780-42e1-a745-96f8169ea148:
- protocol=#002-E status=COMPLETED acc=3 rej=1 warn=2
- rejected_rows[1]: `{row:1, code:EMPTY_DESCRPROD, column:DESCRPROD, detail:"descrprod is required"}` (row/code/detail present) ✓
- warnings[2]: `{row:2, code:INVALID_EAN, column:EAN, offending_value:7894900011518}`, `{row:3, code:INVALID_NCM, column:NCM, offending_value:12AB}` (itemized) ✓

---

## PART A cross-check (from p7-qa-live-drive.md, UI render live-driven)
PART A anchor claims corroborated 1:1 by PART B step 1 API data:
- PART A "sem confiança sem candidato" / "Sem candidato" ↔ API confidence=0 int + match_status=NO_CANDIDATE (no fake 0%/NaN). ✓
- PART A anchor list (seller_sku/ean/marca/refforn all UNAVAILABLE with honest detail) ↔ API reasons[] identical anchors/direction/detail. ✓
- PART A "34 candidates rendered" ↔ API count=34. ✓
NO_CANDIDATE honesty (ADR-17) live + API agree. No contradiction.

## Residual live-state (cleanup)
- MLB3758134297: reverted to `current_link.state=none` — CLEAN (back in pending queue).
- MLB3758134295: **residual — could NOT fully reverse to pending.** `current_link.state=resolved, internal_product_id=2002` (valid synthetic product, NOT phantom 999999 — the 999999 write was undone). Full unresolve-to-none is blocked because Step 4's idempotent re-resolve superseded the original none→resolved audit and there is no unresolve-to-pending API path (undo only reverses the latest non-undo audit). Benign: points at real synthetic product 2002. Audit history rows remain (append-only by design). RECOMMEND hub reset MLB3758134295's link (regen candidates / DB) if the demo needs it back in the Pendentes queue.
