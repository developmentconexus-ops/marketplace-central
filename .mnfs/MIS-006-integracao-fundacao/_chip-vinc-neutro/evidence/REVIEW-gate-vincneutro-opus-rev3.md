# ROUND 5 DELTA GATE — OPUS SIDE — verdict

**PROVENANCE, declared before the content.** Seat: `harness:gate-reviewer` subagent on **opus**,
tool-set **Read / Grep / Glob — no Bash, no Write**. A seat with no Write cannot author a file, so
this verdict is **TRANSCRIBED from the task result, not captured from disk**. That is the expected
form of this seat, not an anomaly (hub ruling on ledger row 6); the residual risk is **omission** in
transcription, not fabrication, and it is inauditable from this artifact by construction. Mitigation
per §11: the text below was pasted verbatim as the FIRST act after the completion notification
arrived, before any analysis of it.

Frozen prompt: `evidence/PROMPT-gate-round5-delta.md`. Range: `394c83c..3915f33b`.

---

```
VERDICT: REFUTED
```

```
V1: PASS — QueueRow.tsx:70 `INCOMPARABLE: "bg-info-soft text-info"` + glyph "?" at :254 are real design tokens (packages/ui/src/Badge.tsx:5) and are none of FOR/AGAINST/UNAVAILABLE's pairs; the drawer's former second copy is gone — VinculoDrawer.tsx:6,121 import `directionClass` from QueueRow.
V2: PASS — QueueRow.tsx:399-405 ranks (never filters) over `directionRank`, so `shown` is empty only when `reasons` is; DOM-level proof at QueueTab.test.tsx:207-213, and the new test at QueueTab.test.tsx:248-290 is field-for-field what `applyUnresolvedScore` + `appendProviderDeclaredUnavailableReasons` emit for mercado_livre (0/BAIXA/NO_CANDIDATE, seller_sku+ean INCOMPARABLE/erp via generation_service.go:625-626,638-640, title suppressed by :719-721, `marca` UNAVAILABLE via :704 — order and details match).
V3: PASS — the recorded must-fail signature is verifiable by construction: with the reverted enumeration, reasons [INC,INC,UNAVAILABLE] yields `shown = [marca]`, i.e. exactly the "got 1 … the one that survived is – Marca" recorded at V-fixture-producibility-sweep.md:103-109; the execution itself is self-certified, not reproducible in this seat.
V4: PASS — `reasonSideLabel` reads the field, never `detail` (QueueRow.tsx:212-215), rides inline in the compact chip (:298), and is asserted on producible fixtures (QueueTab.test.tsx:280-281 "falta no ERP"; VinculosDesign.golden.test.tsx:161 "? SKU (falta no anúncio)"). The side-less path is real (generation_service.go:709-711) — but the fixture that carries it is not (finding 1).
V5: NOT-PROVEN — no execution tool in this seat.
V6: PASS — `confirmDecidingAnchors`/`statusDecidingAnchors` are keyed by the full unions (QueueRow.tsx:144,152-165), source is match_status+match_input; the negative case is asserted positionally, QueueTab.test.tsx:651-658 (`td[5].textContent === "—"` for the title-only REVIEW row).
V7: PASS — three DOM cases, ResolvidosTab.test.tsx:80 (badge present, actor_type "system"), :101 (operator → absent), :121 (no resolving audit entry → absent, no "undefined").
V8: PASS — EVIDENCE.md:275-292 records all four items: predicate not implemented as briefed, `rule_matched` grep empty in `contracts/`+`packages/sdk-runtime/src/`, the `0082_product_link_decisions.sql:53-54` CHECK quoted, and the predicate actually used.
V9: PASS — `refforn` KEPT with the reason stated at QueueRow.tsx:231-237,243; unknown anchors still fall through verbatim (QueueRow.tsx:284) and that fallback is exercised on the DOM at QueueTab.test.tsx:108 ("✓ SKU idêntico").
V10: PASS — structural headers neutral (QueueTab.tsx:137-141) while the provider survives as DATA (QueueRow.tsx:463), and zero raw slug on the rendered surface (QueueTab.test.tsx:416,468-469). Injectivity holds byte-for-byte, but not as rendered (finding 4).
V11: NOT-PROVEN — no execution tool in this seat.
```

```
FINDINGS (most severe first)
- apps/web/src/pages/vinculos/QueueTab.test.tsx:177-185 and :334-342 (annotations) — THIRD INSTANCE of the recurring class, inside the write-set: both sweep annotations are FALSE, not merely optimistic. `resolveIdentityAnchors` (generation_service.go:149-169) ABORTS generation unless the declaration resolves, and identity_anchor_adapter.go:28-35 always returns all four `KnownIdentityAnchors`, so EVERY candidate carries exactly one `marca` reason under EVERY declaration. Under the very capability the annotation invokes (`marca` supplied), `cand_inc` would have THREE reasons (marca INCOMPARABLE via :709-711), the cell renders 2 chips + "+1", and the test's own assertion at :213 (`queryByRole(/Mostrar todos os/)` absent) FAILS; under today's ML declaration it also has three (marca UNAVAILABLE) and fails identically. There is no declaration under which the fixture as written occurs. Same for `cand_noside`: a candidate with `internal_product_id: 666` always carries a seeded FOR/AGAINST (generation_service.go:494-495,535,545,551,579,592,606) plus `marca`, so a single-reason array is unreachable, and with two reasons the FOR chip sorts first and `const [chip]` at :353 makes the "? Marca" assertion at :354 fail. The disposition at V-fixture-producibility-sweep.md:82-90 ("both remain valid under a capability declaration, which is a real future provider") is therefore a dodge, not an honest degrade.
- apps/web/src/pages/vinculos/VinculoDrawer.tsx:17-21,78 — the hardening MISSED a site: `bandClasses[candidate.confidence_band]` is an unguarded `Record` index whose value is interpolated into `className` at :30. Round-4 disposition #3 (EVIDENCE.md:796) names this exact defect ("whose pill renders with a class attribute ending in the literal `undefined`") and claims "both hardened" — but the grep it cites was scoped to `QueueRow.tsx`, and the drawer keeps its own copy of the map. Scenario: API ships `confidence_band: "CRITICA"` before SDK regen → the queue row degrades honestly (QueueRow.tsx:55-57, asserted at QueueTab.test.tsx:766-769), the operator opens the drawer via `?candidate=` → the pill renders `class="… font-medium undefined"`, unstyled, with the literal string `undefined` in the DOM. This is the same duplicate-map mechanism the file's own comment at VinculoDrawer.tsx:23-26 says it removed for `direction`.
- .mnfs/.../evidence/V-fixture-producibility-sweep.md:9 — the exhaustiveness claim is false by enumeration: "the other four test files under `pages/vinculos/` construct no candidates". `VinculosPage.test.tsx:45-58` constructs a full `ProductLinkCandidateItem` — a 29th fixture the sweep never saw — and it is non-producible twice over by the sweep's own criteria: `confidence: 95` + `confidence_band: "ALTA"` + `match_status: "REVIEW"` (ALTA is emitted only at generation_service.go:503-505 and only for ACCEPT, the exact reasoning the author used to rewrite golden `base()`), and `reasons: []` (impossible, per the `marca` invariant above — the identical shape as round-4 finding #5). The 28-fixture count across the two swept files is correct; the claim about the remaining files is not.
- apps/web/src/pages/vinculos/QueueRow.tsx:104-122 — the round-trip proves injectivity on the STRING, not on the rendered cell, because the cell (`<span className="text-xs text-muted">`, QueueRow.tsx:462) collapses whitespace. `typesetSlug` splits on `_` only, so `amazon__marketplace` → `"Amazon  Marketplace"` and `restored === providerCode` PASSES; `_amazon` → `" Amazon"` and `amazon_` → `"Amazon "` also pass. `buildDefinitions` (registry.go:100-115) only `TrimSpace`s and dedupes by exact equality, so `amazon_marketplace` and `amazon__marketplace` are both registrable simultaneously — and both render as "Amazon Marketplace" to the operator, which is the precise failure the docblock at :83-98 claims the round-trip closes ("two providers wearing one name is wrong information"). The `amazon-marketplace` case the test at :443-465 covers is genuinely fixed; the whitespace family is not.
```
