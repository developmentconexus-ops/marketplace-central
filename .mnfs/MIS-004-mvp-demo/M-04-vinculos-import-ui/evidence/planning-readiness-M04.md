# M-04 Planning-Readiness Memo (CHIP-M04)

Synthesis of D01 (baseline) + D02 (gap-closer). Feeds the P2 batch planner. Two findings need
hub adjudication BEFORE the plan is dispatch-ready (§4.1: non-empty open_questions ⇒ not ready).
Status: **BLOCKED on codex quota + 2 contract/seam rulings.** Evidence: `evidence/d01-*.md`, `evidence/d02-*.md`.

---

## Baseline reality (load-bearing)

- **product_links today** (`apps/server_core/internal/modules/product_links/`): candidates carry
  a SINGLE `match_input` anchor (seller_sku|ean|title|manual|none) + `match_value`; NO
  confidence/band/reasons. Generation ladder = exact SKU / exact EAN / conflict / title / unresolved
  (`generation_service.go:120-145`). Resolution flows: approve-candidate / reject-listing /
  manual-resolve → persists `product_links` + `product_link_audit_entries` (already state/actor/reason).
  Routes are `/product-links/link-resolutions/*` (NOT the brief's `/resolutions/*` — planner aligns naming).
- **Migrations**: top = 0047; runner_test.go count = **44** (`runner_test.go:25,64`). My block **0065-0069**
  is clear (non-contiguous numbering already exists — 0024 skipped). Each new migration bumps the fixture.
- **Audit precedent**: `product_link_audit_entries` (0025) + `erp_import_protocols/issues` (0045/0047) —
  batch audit table mirrors these (tenant_id, actor=`operator` literal MVP, per-item results, batch_id).

## Finding A — SDK seam reality (REQUEST: additive index.ts lock; the `productLinks.ts` grant is moot)

- sdk-runtime is **hand-authored, NOT generated** (no codegen; `tsc --noEmit` only). ALL domain client
  methods live INLINE in `packages/sdk-runtime/src/index.ts` (1804 lines) inside one object literal in
  `createMarketplaceCentralClient` (index.ts:1322+). The **7 existing product-links methods are already
  there** (index.ts:1508-1545); their types at index.ts:850-949.
- `erpImport.ts` (types-only + one barrel line) is a ONE-OFF exception that required an explicit
  hub-owned BARREL grant (HUB-LEDGER D-13). The default/idiomatic path for a new domain = add types +
  inline methods in index.ts, same commit as OpenAPI.
- **Implication**: the collision-matrix grant `sdk-runtime/src/productLinks.ts` (new file) does NOT match
  repo idiom — splitting only the NEW batch/undo methods into a separate file while the existing 7 sit in
  index.ts is an idiom-mismatch (AI-slop). Correct move = add batch-preview/batch/undo methods + enriched
  candidate types (confidence/band/reasons) INLINE in index.ts alongside the existing product-links block.
- **ASK**: grant a TEMPORARY ADDITIVE-ONLY lock on `index.ts` for the product-links method+type block
  (new methods + enriched ProductLinkCandidateItem type; no edits to other domains' lines), released at
  CLOSED, diff called out in CLOSED payload — same mechanism as BARREL-01/D-13. If hub prefers the
  productLinks.ts-new-file path anyway (to keep index.ts hub-exclusive), say so and I'll follow it, but
  flag the idiom split in evidence.

## Finding B — IC-01 anchor model not satisfiable from provider data (ESCALATION: IC-01 owner ruling)

- IC-01 auto-ACCEPT requires **2 independent anchors** from {EAN, marca, refforn}. But:
  - INTERNAL product side HAS all fields — `ProductCandidate` exposes BrandName, ReferenceCode(refforn),
    NCM, EAN (`internal_read/domain/internal_product.go:3-21`).
  - PROVIDER (listing snapshot) side has ONLY seller_sku + ean + title (migration 0022 +
    `listing_snapshot.go`) — **no brand, no manufacturer_reference**. So marca/refforn CANNOT be
    cross-matched; they're computable only internal-side.
  - EAN caveat: internal "EAN" is TGFPRO.REFERENCIA-derived; `reader.go:431-443` explicitly marks
    EAN-only linking UNPROVED and returns unsupported for EAN-only queries.
- **Consequence**: the only cross-side anchors that exist are **seller_sku (→codprod)** and **EAN (weak)**.
  IC-01's marca/refforn anchor pairs are structurally uncomputable in the demo data. Title = ranking only
  (IC-01), never an anchor.
- IC-01 says "Must Not Decide In Feature Execution: lista de âncoras, thresholds" — so I CANNOT narrow the
  anchor model unilaterally. **ASK (IC-01 owner = Mission Strategist/hub)**: ratify the feasible model:
  - Anchors available: `seller_sku` (exact→codprod) + `ean` (corroboration, flagged unproved).
    `reasons[]` reports marca/refforn as `UNAVAILABLE` (honest, ADR-17 — motivo sempre visível).
  - Bands: seller_sku+ean agree on same product ⇒ ALTA (auto-ACCEPT proxy for "2 anchors");
    single anchor ⇒ MEDIA; title-only ⇒ BAIXA; SKU/EAN conflict or hard-negative title contradiction
    (kit/combo/cor/medida/voltagem, IC-01) ⇒ cap BAIXA + reason AGAINST; EAN absent ⇒ max MEDIA/REVIEW.
  - Confirm whether `seller_sku` exact counts as an ACCEPT-grade anchor (IC-01 says seller_sku resolves
    ONLY to codprod — arguably the strongest identity claim, but IC-01's ACCEPT text is anchored on
    EAN/marca/refforn, not seller_sku).
- **Dependency**: M-02 F-03 "shared collision fixtures" do NOT exist on base 28b8447c (M-02 in-flight,
  own branch). F-01 must author its OWN IC-01-derived fixtures. The anchor-model ruling above must
  propagate to M-02 for the "verdade única" invariant to hold.

## Finding C — FE has no shared tabs/modal/toast (no ruling needed; page-local per precedent)

- No shared Tab/Modal/Toast in packages/ui. Precedents: AnunciosPage local tabs (`role="tab"`),
  `pages/mutations/MutationPreviewModal.tsx` local `role="dialog"`, `pages/mutations/` sub-component
  folder. DetailPanel (not DetailDrawer — unused) is the real drawer, wrapped page-local
  (`ListingDetailPanel.tsx`). NO toast → success/error via inline banner (ErrorState style) or a
  page-local toast component. F-02 builds tabs/modal/(toast?) as page-local components in
  `pages/vinculos/**` — within ownership, no packages/ui change.
- web-query: `linkageQueryKeys.workflows` exists; `invalidateAfterMutation` already maps
  `link_apply → [listings, linkage, catalog, mutations]` (invalidation.ts:26) — F-02 REUSES it.
  `QUERY_STALE_TIME` has NO `linkage` key; new queue/resolved keys ⇒ either page-local keys OR a small
  additive edit to `packages/web-query/src/index.ts` (shared seam — REQUEST if we touch it; page-local
  keys avoid the REQUEST). Planner decides; default = page-local to avoid the seam.

## Candidate slice skeleton (INPUT to the planner, NOT a ratified plan)

F-01 (backend, after Finding A+B rulings):
- S1 (complex): migration 0065 ALTER product_link_candidates + candidates confidence/band/reasons;
  generation computes confidence per ratified anchor model; backfill next-generation only. Fixture bump.
  Fixtures ≥8 IC-01 cases (band+reasons exact). [C-none-direct; feeds C03 presentation]
- S2 (standard): migration 0066 batch audit table; `POST link-resolutions/batch-preview` dry-run itemized
  (OK|FAILED+cause, nothing persisted). [C01]
- S3 (complex): `POST link-resolutions/batch` apply-per-item reusing approve-candidate logic, partial
  failure itemized, batch audit row. [C02]
- S4 (standard): `POST link-resolutions/{id}/undo` (404 missing / 409 ALREADY_UNDONE / 409 SUPERSEDED),
  audited reversal. [C02]
- S5 (standard): OpenAPI /product-links/* additive + SDK methods/types inline (Finding A seam). [C05]

F-02 (frontend, after F-01 OpenAPI committed):
- S6: routes/vinculos.tsx content + pages/vinculos/ scaffold; Fila/Resolvidos tabs, KPIs. [C03]
- S7: queue rows (produto+candidato+banda%+anchor chips), drawer (?candidate= deep-link), individual
  approve/reject via existing manual-resolve/reject-listing. [C03]
- S8: bulk select → preview modal (dry-run) → apply → applied/failed feedback → refetch (reuse
  link_apply invalidation); NO_CANDIDATE honest state. [C03]
- S9: Importação section reads GET /erp/imports (protocol #NNN-E, status, counts, rejection report). [C04]

Per-criterion verification map + write-sets + contract-satisfiability = planner output once A+B ruled.

## Open questions (blockers to dispatch-ready plan)
1. codex quota exhausted til Jul-25 (post-demo) → execution lane ruling (BLOCKED sent).
2. Finding A: index.ts additive lock vs productLinks.ts-new-file.
3. Finding B: feasible anchor model ratification (IC-01 owner).
4. Who runs P2 planner (Sol unavailable) — cold Opus subagent? (part of lane ruling).
