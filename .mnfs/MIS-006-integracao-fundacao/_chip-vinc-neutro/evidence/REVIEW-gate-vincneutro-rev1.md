# Gate review round 1 — CHIP-VINC-NEUTRO

Reviewer: gpt-5.6-sol / effort medium, OS-process (codex-cli 0.144.4, explicit binary), --sandbox read-only.
Prompt: see PROMPT-gate-vincneutro-rev1.md. Reviewed commit fa6ca3a2 against base bcab8269.
Verbatim final message follows.

```
VERDICT: REFUTED
V1: PASS — `INCOMPARABLE` renders glyph `?` and its exclusive `bg-info-soft text-info` pair in apps/web/src/pages/vinculos/QueueRow.tsx:42.
V2: PASS — the all-`INCOMPARABLE` test asserts two rendered `motivo-chip` DOM elements in apps/web/src/pages/vinculos/QueueTab.test.tsx:177.
V3: PASS — reverting the ranking produced zero motivo chips and a lone `+2`, failing with “Unable to find…motivo-chip” in .mnfs/MIS-006-integracao-fundacao/_chip-vinc-neutro/evidence/V3-must-fail.txt:22.
V4: PASS — `side` is read directly and rendered inline while absent `side` remains unreported in apps/web/src/pages/vinculos/QueueRow.tsx:96.
V5: PASS — the three `/vinculos` errors at baseline disappear while the twelve declared errors remain in .mnfs/MIS-006-integracao-fundacao/_chip-vinc-neutro/evidence/L0-tsc-after.txt:1.
V6: FAIL — `match_status` and `match_input` are present in packages/sdk-runtime/src/index.ts:1080 and already determine the decision anchors exactly as implemented by apps/server_core/internal/modules/product_links/application/resolution_service.go:812, yet apps/web/src/pages/vinculos/QueueTab.tsx:137 retains `Canal` and `GTIN` instead of rendering `Identificado por`.
V7: PASS — the resolving audit entry is selected and tested in the DOM for system, operator, and missing-audit cases in apps/web/src/pages/vinculos/ResolvidosTab.test.tsx:81.
V8: PASS — the corrected predicate and both brief defects are recorded in apps/web/src/pages/vinculos/useVinculosResolved.ts:32 and independently corroborated by apps/server_core/migrations/0082_product_link_decisions.sql:54 and apps/server_core/internal/modules/product_links/application/resolution_service.go:280.
V9: PASS — historical `refforn` is deliberately retained with verbatim fallback preserved in apps/web/src/pages/vinculos/QueueRow.tsx:112.
V10: PASS — neutral headers retain “Mercado Livre” as provider data without rendering its slug in apps/web/src/pages/vinculos/QueueTab.test.tsx:309, including drawer display-name conversion at apps/web/src/pages/vinculos/VinculoDrawer.tsx:182.
V11: NOT-PROVEN — the write-set excludes both parallel-owned files and the golden gate was strengthened at apps/web/src/pages/vinculos/VinculosDesign.golden.test.tsx:104, but only the after-run count is recorded in .mnfs/MIS-006-integracao-fundacao/_chip-vinc-neutro/evidence/L1-vitest-after.txt:98.

FINDINGS (most severe first)
- apps/web/src/pages/vinculos/QueueTab.tsx:137 — the frozen `Identificado por` contract was skipped despite a correct derivation already existing — `match_status=ACCEPT` plus `match_input=seller_sku` → backend records `concordant_codprod_ean`, but the screen shows separate `Canal`/`GTIN` columns instead of `CODPROD + EAN`.```
