# GATE ROUND 3 (GPT side) — CHIP-VINC-NEUTRO — VERDICT, VERBATIM

Reviewer: `gpt-5.6-sol` / effort `medium`, OS-process dispatch, cold and read-only.
Prompt: evidence/PROMPT-gate-vincneutro-rev3.md (frozen input: commit 394c83c).
Raw stream: scratchpad/agent__gate-rev3.log · final message: scratchpad/agent__gate-rev3.last.md
Copied here verbatim — the scratchpad is not durable, and streaming is not persisting.

```
VERDICT: REFUTED
V1: PASS — QueueRow.tsx:44-49 and QueueRow.tsx:207-212 give INCOMPARABLE its own `?` glyph and `bg-info-soft text-info` tokens.
V2: PASS — QueueRow.tsx:320-335 ranks and slices every typed direction without filtering all-INCOMPARABLE rows.
V3: PASS — EVIDENCE.md:73-87 records the reverted implementation failing with zero chips and a lone `+2`.
V4: PASS — QueueRow.tsx:153-171 maps the wire `side` field to distinct inline operator guidance without inventing absent values.
V5: PASS — EVIDENCE.md:742 records 15→12 TypeScript errors with zero remaining under `pages/vinculos/`.
V6: PASS — QueueRow.tsx:109-133 retains exhaustive status/input Records, and QueueTab.test.tsx:498-548 proves a title-FOR REVIEW does not become a deciding anchor.
V7: PASS — ResolvidosTab.test.tsx:80-136 covers system, operator, and audit-less legacy rows with DOM assertions.
V8: PASS — EVIDENCE.md:272-297 records the stale predicate, empty wire grep, database constraint, and actual `actor_type === "system"` predicate.
V9: PASS — QueueRow.tsx:185-215 preserves historical `refforn` presentation and returns unknown anchors verbatim.
V10: FAIL — QueueRow.tsx:74-78 destructively normalizes distinct unmapped provider codes into identical display strings, erasing provider identity.
V11: PASS — the bcab8269→HEAD app diff contains only allowed `apps/web/src/pages/vinculos/` paths, while EVIDENCE.md:742-744 records 63 files and 522 tests passing.
FINDINGS (most severe first)
- apps/web/src/pages/vinculos/QueueRow.tsx:74 — the fallback is lossy, contradicting its claim that every wire character remains legible — `provider_code="amazon_marketplace"` and `"amazon-marketplace"` are distinct codes accepted by registry.go:104-111 but both render as `Amazon Marketplace`, so the operator cannot identify which provider owns the listing.
- apps/web/src/pages/vinculos/QueueTab.test.tsx:340 — the unmapped-provider test proves only one hard-coded input/output pair and can pass against a fabricated constant fallback — an implementation returning `Shopee` for every unknown code makes `provider_code="shopee"` pass while `provider_code="amazon_marketplace"` incorrectly renders `Shopee`.```
