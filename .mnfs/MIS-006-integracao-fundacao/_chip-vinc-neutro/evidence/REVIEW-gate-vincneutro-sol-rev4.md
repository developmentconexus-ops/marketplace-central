VERDICT: REFUTED
V1: PASS — QueueRow.tsx:66-70 and :250-254 map INCOMPARABLE to glyph `?` and its own `bg-info-soft text-info` tokens.
V2: PASS — QueueTab.test.tsx:248-289 uses the exact mercado_livre unresolved reason set and asserts two rendered actionable chips ahead of the permanent `marca` absence.
V3: PASS — evidence/V-fixture-producibility-sweep.md:101-116 records the producible V2 probe failing with one chip under the old three-direction enumeration.
V4: PASS — QueueRow.tsx:212-225 derives the displayed side only from `reason.side`, while QueueTab.test.tsx:215-219 asserts both provider and ERP side labels in the DOM.
V6: PASS — QueueRow.tsx:144-176 exhaustively derives deciding anchors from status plus match_input, and QueueTab.test.tsx:607-657 includes the required title-FOR negative case.
V7: PASS — useVinculosResolved.ts:49-50 tests the resolving audit actor for `system`, with system, operator, and missing-audit DOM cases at ResolvidosTab.test.tsx:80-137.
V8: PASS — EVIDENCE.md:272-297 records the stale predicate, empty wire grep, schema prohibition, and the actual `actor_type === "system"` predicate.
V9: PASS — QueueRow.tsx:232-243 declares why historical `refforn` remains mapped, while :283-284 returns unknown anchors verbatim.
V10: FAIL — QueueRow.tsx:104-121 accepts doubled and edge underscores even though HTML collapses the resulting whitespace, so distinct registered provider codes can still look identical.
FINDINGS (most severe first)
- apps/web/src/pages/vinculos/QueueRow.tsx:104 — the round-trip checks the pre-layout string, not its visible HTML representation — `amazon_marketplace` → `Amazon Marketplace` and `amazon__marketplace` → `Amazon  Marketplace` → both visibly render as “Amazon Marketplace”.
- apps/web/src/pages/vinculos/QueueTab.test.tsx:177 — the sweep falsely calls `cand_inc` producible under a provider declaring `marca` supplied — its two seeded INCOMPARABLE reasons → finalizer necessarily appends a third side-less `marca` INCOMPARABLE reason, contradicting the asserted two-chip/no-toggle row.
- apps/web/src/pages/vinculos/VinculoDrawer.tsx:17 — the claimed exhaustive union-index hardening missed the drawer confidence-band map at :78 — wire band `CRITICA` → visible label but class interpolation ending in literal `undefined`.
- apps/web/src/pages/vinculos/QueueTab.test.tsx:121 — the sweep classifies a conclusion-dependent impossible fixture as presentation isolation — current CONFIRM generation cannot emit removed `refforn`, yet the test calls this the “real CONFIRM shape” and requires four reasons, `+2`, and expanded `refforn` text.