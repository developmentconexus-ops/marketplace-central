# Gate review round 2 — CHIP-VINC-NEUTRO (the V6 fix)

Reviewer: gpt-5.6-sol / effort medium, OS-process (codex-cli 0.144.4, explicit binary), --sandbox read-only.
Prompt: see PROMPT-gate-vincneutro-rev2.md. Reviewed commit 7a343fea (base bcab8269).
Verbatim final message follows.

```
VERDICT: REFUTED
V1: PASS — `INCOMPARABLE` has glyph `?` and distinct `bg-info-soft text-info` tokens in apps/web/src/pages/vinculos/QueueRow.tsx:44.
V2: PASS — the all-`INCOMPARABLE` DOM test requires visible motivo chips in apps/web/src/pages/vinculos/QueueTab.test.tsx:177.
V3: PASS — the recorded regression produces zero chips and one failing test in .mnfs/MIS-006-integracao-fundacao/_chip-vinc-neutro/evidence/V3-must-fail.txt:290.
V4: PASS — `side` is read only for `INCOMPARABLE` and rendered through the chip label in apps/web/src/pages/vinculos/QueueRow.tsx:145.
V5: PASS — the recorded post-change compilation retains only the twelve declared out-of-scope baseline errors in .mnfs/MIS-006-integracao-fundacao/_chip-vinc-neutro/EVIDENCE.md:111.
V6: FAIL — The status/input derivation at apps/web/src/pages/vinculos/QueueRow.tsx:85 faithfully mirrors resolution_service.go:812, but retaining `Canal` at QueueTab.tsx:138 overrides the frozen requirement that the new column replace both former slots at DECISOES-D122-anchors-telas.md:136.
V7: PASS — system, operator, and missing-audit DOM cases are covered at apps/web/src/pages/vinculos/ResolvidosTab.test.tsx:80.
V8: PASS — the unavailable `rule_matched`, forbidden system pairing, and actual predicate are recorded at .mnfs/MIS-006-integracao-fundacao/_chip-vinc-neutro/EVIDENCE.md:245.
V9: PASS — historical `refforn` remains mapped while unknown anchors fall through verbatim in apps/web/src/pages/vinculos/QueueRow.tsx:165.
V10: PASS — neutral headers retain provider identity as the display name and reject the raw slug in apps/web/src/pages/vinculos/QueueTab.test.tsx:309.
V11: PASS — persisted counts rise from 511 to 520 passing tests at .mnfs/MIS-006-integracao-fundacao/_chip-vinc-neutro/EVIDENCE.md:327, and the commit diff contains only permitted `/vinculos` paths.
FINDINGS (most severe first)
- apps/web/src/pages/vinculos/QueueTab.tsx:138 — the chip keeps the former `SKU ML` slot as `Canal`, overriding the frozen two-column replacement instead of relocating provider identity — any candidate with `provider_code=mercado_livre` → nine-column table containing both `Canal` and `Identificado por`, rather than `Identificado por` replacing the former `SKU ML` and `GTIN` slots.```
