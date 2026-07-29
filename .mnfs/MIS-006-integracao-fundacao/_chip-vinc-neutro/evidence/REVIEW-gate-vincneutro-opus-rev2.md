# GATE ROUND 2 (Opus side) — CHIP-VINC-NEUTRO — VERDICT

Reviewer: cold Opus subagent, read-only by construction (no Bash, no Edit/Write). Independent of
this session, which is the implementer (HARNESS-CORE §4 item 3, implementer ≠ reviewer).
Prompt: evidence/PROMPT-gate-vincneutro-rev3.md (the same frozen input the GPT side got).
Frozen input: commit `394c83c`.

## PROVENANCE — read this before trusting the text below

The task's on-disk output artifact
(`…/tasks/adfed0fb58dc6c987.output`) is **0 bytes**. The verdict reached this session only through
the completion notification. What follows is transcribed from that notification, not copied from a
file, and I am recording that rather than presenting a transcription as a captured artifact — this
chip's own rule is that streaming is not persisting and an empty artifact means the pass did not
happen.

What makes the round auditable anyway: every finding below was verified INDEPENDENTLY against the
named source before being acted on, and those verifications are written up with file:line in
evidence/V-round4-must-fail.txt. The findings stand on that verification, not on the transcription.

## Two limits of this round, stated up front

1. **It reviewed `394c83c`.** That is the commit BEFORE the `providerDisplayName` injectivity fix
   the GPT side (rev3) forced. Its `V10: PASS` — and specifically its explicit refutation of
   attack 2, "the split `/[_\-\s]+/` only collapses separators, so no registered `provider_code`
   typesets into another provider's name" — is **wrong**, and the GPT side caught what this side
   missed. Both `amazon_marketplace` and `amazon-marketplace` are registrable (registry.go
   dedupes by exact string equality) and both typeset to `Amazon Marketplace`. This side's V10 PASS
   must NOT be read as corroborating that function.
2. **Three criteria came back NOT-PROVEN for want of tools, not for want of substance** — V5, V8,
   V11 — because the seat had no Bash (so no `tsc`, no `vitest`, no `git diff`) and was not given
   EVIDENCE.md. Those are honest NOT-PROVENs about THIS SEAT's reach, not defects found.

## VERDICT: REFUTED

    V1: PASS — `INCOMPARABLE` has its own glyph `?` (QueueRow.tsx:211) and its own real token pair
    `bg-info-soft`/`text-info` (QueueRow.tsx:48, defined at apps/web/src/index.css:27-28 + 74-75,
    pre-existing and used by packages/ui/src/Badge.tsx:5), distinct from FOR/accent, AGAINST/warn
    and UNAVAILABLE/surface-2; the drawer's former third map is gone — it now imports the same
    Record (VinculoDrawer.tsx:6,121) — and the tokens are DOM-asserted at QueueTab.test.tsx:213-218.
    V2: PASS — an all-INCOMPARABLE row renders 2 real `motivo-chip` nodes and no bare "+N"
    (QueueTab.test.tsx:196-201, assertions on the rendered DOM), because the ranking is a total sort
    over `Record<Direction, number>` (QueueRow.tsx:320-326) so `shown` is empty only when `reasons` is.
    V3: PASS (artifact-only — I have no Bash and cannot re-run vitest) — evidence/V3-must-fail.txt:290-291
    shows `1 failed | 9 skipped` with the failure landing on the `getAllByTestId("motivo-chip")` line
    (:282) under the reverted three-direction enumeration.
    V4: PASS — `side` is read from the field, never parsed from `detail` (QueueRow.tsx:169-172), rides
    inline in the compact chip (QueueRow.tsx:229) and is asserted as `? SKU (falta no anúncio)` /
    `? EAN (falta no ERP)` / `? EAN (falta nos dois lados)` (QueueTab.test.tsx:206-207,245); the
    side-less INCOMPARABLE path stays silent (QueueTab.test.tsx:275-277) and matches
    generation_service.go:711.
    V5: NOT-PROVEN — no Bash in this seat, so `tsc -p apps/web --noEmit` could not be run; the pack
    artifact evidence/L0-tsc-after-opus-fixes.txt:1-16 lists exactly 12 errors, none under
    `pages/vinculos`, but that is the chip's own transcript, not an independent run.
    V6: PASS — `decidingAnchors` reads the same pair the trail files from: ACCEPT→`CODPROD + EAN`,
    CONFIRM→by `match_input`, all else empty (QueueRow.tsx:101-122) mirrors `decisionRuleForCandidate`
    verbatim (resolution_service.go:812-835); the `title FOR` negative case is DOM-proven with
    `td[5] === "—"` (QueueTab.test.tsx:542-548), and `td[5]` is still the Identificado por cell by
    construction (QueueRow.tsx:360,373,382,389,409,424). Attack 1 refuted by construction, not by
    comment: the `Record<ProductLinkMatchStatus, …>` annotation on the map literal (QueueRow.tsx:109-112)
    is untouched, so a sixth SDK status fails TS2741 on the literal; `if (!rule)` (QueueRow.tsx:132) is
    a use-site truthiness check that TS permits and that cannot relax an object-literal completeness
    check (tsconfig.base.json:7 `strict`, no `noUncheckedIndexedAccess`).
    V7: PASS — predicate is the resolving audit entry's `actor.actor_type === "system"`
    (useVinculosResolved.ts:49-51), badge rendered only when true (ResolvidosTab.tsx:53-61), three DOM
    cases green: system→present, operator→absent, pre-M-05 empty audit→absent with no `undefined`
    (ResolvidosTab.test.tsx:96,117,132-135).
    V8: NOT-PROVEN — I was barred from reading EVIDENCE.md, so I cannot confirm the pack records it; the
    substance is independently true (zero `rule_matched` hits outside the chip's own comments in
    contracts//sdk-runtime; migrations/0082_product_link_decisions.sql:54 is verbatim
    `CHECK (actor <> 'system' OR rule_matched = 'concordant_codprod_ean')`; the predicate used is
    documented at useVinculosResolved.ts:29-50).
    V9: PASS — `refforn` was KEPT, with the reason written down (QueueRow.tsx:189-201: D-A does not
    migrate already-persisted reasons), and `anchorShortLabel` still returns unknown anchors verbatim
    (QueueRow.tsx:215), exercised live by the `refforn` motivo at QueueTab.test.tsx:136,158.
    V10: PASS — structural labels neutral (QueueTab.tsx:137-141; ResolvidosTab.tsx:104) while the
    provider survives as DATA (QueueRow.tsx:384 and VinculoDrawer.tsx:182 via `providerDisplayName`),
    and no raw slug reaches the DOM (QueueTab.test.tsx:337,360; VinculosDesign.golden.test.tsx:127).
    Attack 2 refuted: the split `/[_\-\s]+/` (QueueRow.tsx:75) only collapses separators, so no
    registered `provider_code` typesets into another provider's name; the only collisions possible are
    separator/duplication variants of one slug, and the `anchorShortLabel` asymmetry is mandated by V9
    itself, not a contradiction.
    V11: NOT-PROVEN — the vitest half is artifact-only but internally consistent (511 tests/62 files
    green at evidence/L1-vitest-baseline.txt:82-83 → 522/63 green at
    evidence/L1-vitest-after-opus-fixes.txt:83-84, golden 3 green at :51, no green→red); the scope half
    `git diff --name-only … HEAD` CANNOT be run in this seat (no Bash), so zero-paths-outside-
    `pages/vinculos` and zero touches of `VinculosPage.tsx`/`ImportacaoSection.tsx` are unverified —
    filesystem-only observations: the `vitest.chip.config.ts` leftover is gone (only
    apps/web/vitest.config.ts exists) and the `bg-info-soft` tokens needed no `index.css` edit.

    FINDINGS (most severe first)
    - VinculosDesign.golden.test.tsx:42-79 and :195-208 — both "made honest" fixtures are still NOT
      wire-producible: the reasons array omits the `marca` motivo the generator appends to EVERY
      candidate of EVERY provider — a real `mercado_livre` exact_ean CONFIRM (60/MEDIA, ean FOR +
      seller_sku INCOMPARABLE/provider) is emitted with a THIRD motivo, because
      identity_anchor_adapter.go:28-35 always declares all four `KnownIdentityAnchors`
      (connectors/ports/marketplace_capability.go:40-45), generation_service.go:159-160 hard-fails
      rather than degrading that list, and `classifyProviderIdentityAnchor` always emits for `marca` —
      UNAVAILABLE at generation_service.go:704 when not supplied, side-less INCOMPARABLE at :711 when
      supplied (no `marca` case in `identityAnchorValues`, :749-750; locked by
      generation_service_test.go:521-525) → the production Motivo cell shows 2 chips PLUS a "+1"
      toggle, while the golden locks a 2-chip cell with no toggle (:143-146), so the EXEMPLO-IO row
      still gates a layout the backend never produces and the "fixture is now producible" claim is
      false.
    - QueueRow.tsx:230, :243, :324 — the drift hardening was applied to `match_status` only;
      `direction` is still an unchecked `Record` index, and the fix made that path newly REACHABLE
      (the old literal enumeration dropped unknown directions from the compact cell, the new total
      sort keeps them) → wire reason `{anchor:"ean", direction:"PARTIAL", detail:"x"}` renders the
      chip text `undefined EAN`, a className ending in `undefined`, and `directionRank[...]` undefined
      makes the comparator NaN so the cell's order is unspecified — the literal string "undefined" on
      screen is the same failure V7 case 3 forbids.
    - VinculosDesign.golden.test.tsx:154-171 and :209-217 — the NO_CANDIDATE fixtures inherit
      `state: "exact_ean"` from `base()` while carrying `match_status: "NO_CANDIDATE"`,
      `match_input: "none"`, `reasons: []`, which `applyUnresolvedScore` cannot produce
      (generation_service.go:620-628 always emits at least the two absence reasons, plus `marca`, and
      the state on that path is `unresolved`) — same non-producible-fixture class the fix claims to
      have closed, in the same file the fix edited, left untouched.
    - VinculosDesign.golden.test.tsx:134 vs :224-227 — after the swap the EXEMPLO-IO test asserts no
      confidence-band pill at all (only `60%`), so the ALTA/accent token path is covered only by the
      negative off-theme sweep, which still passes if the pill loses its token entirely; whether the
      pre-swap golden asserted that pill positively is NOT-PROVEN here (requires git).

## DISPOSITION

All four findings verified independently against the named sources and ACCEPTED. Fixes and the
verbatim verification trail: evidence/V-round4-must-fail.txt, sections A, C and D.

Finding #2 was extended, not merely applied: a grep-by-STRING for the same shape across the
write-set found a fourth unguarded site the reviewer did not name — `bandLabels`/`bandClasses`
indexed by `confidence_band` at the CONFIANÇA cell — which renders a pill whose class attribute
ends in the literal `undefined`. Hardened in the same pass.

Finding #4 is the one the fix could not fully close on its own terms: whether the PRE-swap golden
asserted the ALTA pill positively is answerable with git, and the answer is that it did not — the
pre-swap EXEMPLO-IO asserted `92%` and no band pill either. So the swap did not LOSE that coverage;
it never existed. The finding is still correct that the coverage SHOULD exist, and it now does, on
both bands.
