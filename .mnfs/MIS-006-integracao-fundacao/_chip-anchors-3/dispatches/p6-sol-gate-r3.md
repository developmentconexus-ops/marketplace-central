## PERSISTED BY THE ORCHESTRATOR — read this header before the verdict

This artifact was written by CHIP-ANCHORS-3 (the orchestrating session), not by the seat that
produced it. The seat is GPT-5.6 Sol, effort medium, dispatched via `codex:codex-rescue` as the
round-3 side-B P6 gate reviewer.

**Why the seat did not write it.** The brief ordered the reviewer to write this exact path. It
tried, with `apply_patch` / `*** Add File:` on this path, and the runtime refused:

    patch rejected: writing is blocked by read-only sandbox; rejected by user approval settings

(codex rollout `019fa9da-fb61-7532-a3c4-8216d5fe04b3`, line 127, 2026-07-28T18:03:59.806Z). The
seat then reported in its final message: "I could not persist `p6-sol-gate-r3.md`: the workspace is
enforced read-only, and the write was rejected. The file remains absent." That report was true —
the file was absent until this recovery.

**What this body is.** It is the literal payload of that rejected write, not a summary and not a
reconstruction from prose: the `*** Add File:` hunk of the `apply_patch` call, with the apply-patch
`+` prefix stripped from each of its 208 lines. 16438 characters. Zero lines failed to carry the
`+` prefix, so nothing was interpolated. It is the file the seat authored and was prevented from
saving — closer to a contemporaneous write than the round-1 salvage in `p6-opus-gate-r1.md`, which
was recovered from a returned message. It is still a recovery, and it is labelled as one.

- Seat: GPT-5.6 Sol, effort medium, `codex:codex-rescue`, task-id `a8797556cdcca4ff1`, job `bzao5338y`.
- Brief: `dispatches/p6-gate-brief-r3.md`. HEAD under review: `0264ba84`.
- Extraction: `scratchpad/roll_extract.py` over the rollout jsonl; recovered 2026-07-28.

**FINDING for the hub (profile §11):** the read-only reviewer seat and the "write your verdict to
this path" instruction are in direct contradiction, on BOTH sides of the gate and for different
reasons — `harness:gate-reviewer` has no Write tool, and the codex read-only sandbox rejects the
patch. Round 1 lost side A's artifact to this. Round 3 would have lost side B's. The brief cannot
discharge artifact persistence by ordering the seat to do it; persistence is the orchestrator's
step, in the turn the verdict arrives.

---

VERDICT: REFUTED

## SWEEP

Grepped the whole of:

- `apps/server_core/internal/modules/product_links/application/generation_service_test.go`
- `apps/server_core/internal/modules/product_links/application/generation_service.go`

Pattern, case-insensitive and as word occurrences: `never|always|only|unreachable|cannot|impossible|no longer|every`.

Result: **50 occurrences total** — 27 in `generation_service_test.go`, 23 in
`generation_service.go`. Each occurrence:

### `generation_service_test.go` — 27 occurrences

1. `:222` `[every]` — **REFUTED**. `"A capped generation run leaves every listing beyond the cap without a candidate"` is broader than the code. `link_candidate_repo.go:38-48` deletes only identities supplied by the capped snapshot set, so an unprocessed listing can retain an existing candidate. Correct: the run does not generate or refresh candidates beyond the cap.
2. `:400` `[Every]` — **TRUE only for non-nil product fixtures**, but the following total claim is false. Every non-nil `product:` literal in the table has a positive canonical ID; the unresolved rows have `product == nil`.
3. `:402` `[unreachable]` — **REFUTED**. ``a `ProductCandidate{}` is unreachable in production`` is false: `generation_service.go:215` constructs the production unresolved candidate with `internalreaddomain.ProductCandidate{}`.
4. `:419` `[cannot]` — **TRUE in its stated scope**. A product that came through `findProducts` has a positive canonical ID (`generation_service.go:277-283`), so a present production product supplies the seller_sku ERP-side CODPROD.
5. `:475` `[impossible]` — **TRUE in the generator scope**. Resolved products reaching this path passed the canonical-ID filter.
6. `:477` `[no longer]` — **TRUE**. `identityAnchorValues` derives seller_sku's ERP value from `canonicalProductID`, not `ReferenceCode`.
7. `:497` `[no longer]` — **TRUE** for the same reason; changing `referenceCode` cannot change this comparison.
8. `:505` `[Every]` — **TRUE**. The table runner constructs `product` for every case at `:523-526` and the EAN matcher returns it at `:527-529`.
9. `:506` `[always]` — **TRUE**. That runner's product has `InternalProductID: canonicalIDPtr(101)`.
10. `:508` `[cannot]` — **TRUE in the stated table/product-present scope**. The canonical product makes seller_sku's ERP-side value non-empty.
11. `:544` `[always]` — **TRUE in the stated production scope**. The only production call at `generation_service.go:303` passes `&product`.
12. `:926` `[never]` — **TRUE**. With an empty snapshot slice, the loop at `generation_service.go:102` has no iteration.
13. `:928` `[never]` — **TRUE as the quoted erroneous interpretation** of a missing guard: an unconfigured engine would otherwise look like an empty result.
14. `:1235` `[ONLY]` in `"LEGACY-ONLY"` — **NON-CLAIM**, fixture data.
15. `:1235` `[only]` in `"Legacy only"` — **NON-CLAIM**, fixture data.
16. `:1238` `[ONLY]` in `"sku:LEGACY-ONLY"` — **NON-CLAIM**, fixture key.
17. `:1285` `[ALWAYS]` — **TRUE in the fixture scope**. The Mercado Livre declaration marks `marca` unsupplied and `assertProviderDeclaredUnavailableReasons` pins the resulting UNAVAILABLE reason.
18. `:1507` `[only]` — **NON-CLAIM**, table-case label `"listing side only"`.
19. `:1508` `[only]` — **NON-CLAIM**, table-case label `"internal side only"`.
20. `:1557` `[never]` — **TRUE**. The single-anchor branch is MEDIA/CONFIRM and `autoApprovals` selects only ACCEPT.
21. `:1616` `[always]` — **TRUE in the production path**. An EAN-resolved product passed the positive canonical-ID filter.
22. `:1620` `[only]` — **TRUE for the compared old/new behavior**: the seller_sku outlier came from reading `refforn`; the canonical CODPROD path aligns it with EAN.
23. `:1636` `[ONLY]` in `TITLE-ONLY-BAIXA` — **NON-CLAIM**, case label.
24. `:1637` `[only]` in `ranking-only` — **TRUE**: title is used for ranking/hard-negative scoring, not corroborating auto-acceptance.
25. `:1637` `[never]` — **REFUTED**. `"can never grant ACCEPT/REVIEW-grade confidence"` contradicts the title branch at `generation_service.go:559`, which sets `MatchStatusReview`, and this test's assertion at `:1654-1655`. Correct: title never grants ACCEPT; by itself it yields BAIXA/REVIEW.
26. `:1658` `[only]` in the failure message `ranking-only` — **NON-CLAIM**, assertion diagnostic.
27. `:1799` `[never]` — **TRUE**. `applyUnresolvedScore` seeds seller_sku and EAN reasons; the unresolved candidate is not a silent empty-reasons row.

### `generation_service.go` — 23 occurrences

1. `:85` `[never]` — **TRUE**. `GenerateLinkCandidates` forwards zero, and the repository omits `LIMIT` for `limit <= 0` (`listing_snapshot_repo.go:91-96`).
2. `:87` `[every]` — **REFUTED**. `"a capped run leaves every uncapped listing without a candidate"` ignores retained candidates. `link_candidate_repo.go:38-48` deletes only processed identities. Correct: uncapped listings are not regenerated/refreshed and may have no candidate or retain stale ones.
3. `:89` `[never]` — **TRUE for that invocation**. Listings excluded by the reader's cap never enter the generation loop.
4. `:121` `[only]` — **TRUE**. Persistence returns successfully at `:117-119` before the auto-approval loop begins.
5. `:125` `[Every]` — **TRUE**. The loop visits every element of `pendingApprovals` and joins errors instead of returning early.
6. `:225` `[Only]` — **TRUE as the local ownership rule**. The sole `autoApproval` construction is at `:243`, using the generator's match result.
7. `:231` `[only]` — **TRUE**. The only assignment of `LinkCandidateMatchStatusAccept` in this file is `buildConcordantCandidate:515`.
8. `:235` `[never]` — **TRUE**. The code takes one `len(skuMatches.Products)` reading here and does not perform a second matcher/database count.
9. `:348` `[every]` — **TRUE by canonical product identity**. The nested loops offer every unique product from both anchor match sets.
10. `:473` `[only]` — **REFUTED as worded**. `"seller_sku and ean are the only cross-side anchors available against provider data"` is broader than the implementation: title is also compared across provider listing title and ERP product name for ranking and hard negatives (`:507`, `:558-573`, `:813-846`). Correct: seller_sku and EAN are the only corroborating/auto-accept anchors; title is a cross-side ranking and veto signal.
11. `:474` `[only]` — **TRUE**. Title ranks or vetoes; it does not corroborate.
12. `:474` `[never]` — **TRUE**. Title-only is REVIEW and hard-negative title is REJECT, never ACCEPT.
13. `:561` `[only]` in `"ranking-only"` — **TRUE**, operator-facing detail consistent with the title scoring branch.
14. `:610` `[never]` — **TRUE**. `applyAmbiguousCorroborationScore` unconditionally sets REVIEW, not CONFIRM.
15. `:624` `[never]` — **TRUE**. The unresolved scorer seeds two absence reasons.
16. `:698` `[only]` — **REFUTED**. `"promote only its classification and retain that sentence"` is not total behavior: detail is retained only when the promoted direction is INCOMPARABLE (`:700-702`); UNAVAILABLE replaces the seeded detail with `"provider não fornece..."`. Correct: retain the seeded sentence only for INCOMPARABLE; otherwise replace the reason with the provider-declaration result.
17. `:748` `[never]` — **TRUE**. seller_sku's ERP counterpart is derived from canonical CODPROD, never `refforn`.
18. `:750` `[always]` — **TRUE in the stated generation/product-present scope**. `findProducts` filters missing/invalid canonical IDs before products reach candidate scoring.
19. `:753` `[never]` — **TRUE in that same scope**. A present resolved generator product cannot truthfully be described as lacking CODPROD; the separate unresolved nil-product path has no ERP product.
20. `:804` `[never]` — **TRUE for the stated lowercase/uppercase distinction**. The size regex is case-sensitive and does not classify lowercase `m` as grade `M`.
21. `:858` `[every]` — **TRUE relative to the function's declared regex vocabulary**. It collects all matches from the dimension and grade-size patterns.
22. `:859` `[only]` — **TRUE**. `detectHardNegative` requires both dimension parses to succeed and their normalized signatures to differ.
23. `:860` `[never]` — **TRUE for dimension flagging**. Titles without dimension tokens return `ok == false` and cannot trigger the dimension contradiction branch.

Sweep result: **REFUTED**. False universal claims remain at
`generation_service_test.go:222`, `:402`, `:1637` and
`generation_service.go:87`, `:473`, `:698`.

## (a) Replacement comment

**PASS.**

- `generation_service_test.go:505-511`, exact string checked:
  `"Every case in this table runs with a product present (product above always carries a canonical id) ... with a product present, side=erp cannot arise for seller_sku here."`
- The scope is real, not cosmetic. Every subtest constructs the same non-nil product at `:523-526`,
  with canonical ID 101, and exposes it through the EAN result at `:527-529`.
- `generation_service.go:277-283` filters invalid canonical products, and
  `identityAnchorValues:744-758` therefore produces a non-empty seller_sku ERP-side CODPROD.
- The comment correctly separates the production unresolved nil-product path, pinned at
  `generation_service_test.go:649-674`.

## (b) Four assertions and rewritten test comment

**PASS.**

- `generation_service_test.go:573-575`, exact string checked:
  `findReason(candidate.Reasons, "ean", ...DirectionFor)`. It reads the EAN reason and can fail if
  the scorer omits it or changes its direction.
- `:576-578`, exact string checked: `candidate.Confidence != 95`. It reads the named field and can
  fail independently; the fixture does not prepopulate candidate confidence.
- `:579-581`, exact string checked:
  `candidate.ConfidenceBand != ...LinkCandidateConfidenceBandAlta`. It reads the named field,
  compares to the correct domain constant, and can fail independently.
- `:582-584`, exact string checked:
  `candidate.MatchStatus != ...LinkCandidateMatchStatusAccept`. It reads the named field, compares
  to the correct domain constant, and can fail independently.
- None is tautological: `buildConcordantCandidate:507-516` computes all four outcomes.
- The test still drives the nil-product path directly at `generation_service_test.go:559-565`;
  the additions do not change the fixture or comparison.
- The comment at `:543-550` is accurate in substance. `buildConcordantCandidate`,
  `applySingleAnchorScore`, and the unresolved scoring/classification path are nil-safe; the latter
  two produce absence reasons, while the concordant nil path unconditionally seeds seller_sku/EAN
  FOR and reaches `95 / ALTA / ACCEPT` because the zero product name prevents a hard negative.

## (c) Newest EVIDENCE self-corrections

**PASS on reading; execution-shaped subclaims are execution seat.**

1. `EVIDENCE.md:700`, exact string checked:
   ``TestConcordantCandidateDoesNotDerefNilProduct "fixa exatamente" 95/ALTA/ACCEPT``.
   The retraction is true: the round-2 artifact describes the old test as asserting only nil ID and
   seller_sku FOR; the current test now adds EAN FOR, 95, ALTA and ACCEPT at
   `generation_service_test.go:573-584`.
2. `EVIDENCE.md:701`, exact string checked:
   `"comentário e cobertura corrigidos em 54342331"`.
   The second correction is true: `p6-opus-gate-r2.md:48-58` records the surviving second site at
   then-`:505` and the false degradation sentence; the current delta repairs those sites.
3. `EVIDENCE.md:702`, exact string checked:
   ``dispatches/p6-opus-gate-r1.md citado como artefato — Nunca existiu``.
   The file now exists. Its `SALVAGED ARTIFACT` header at `p6-opus-gate-r1.md:3-21` plainly says it
   did not exist during round 2, identifies transcript source/task/head/brief/recovery time, and
   distinguishes recovery from a contemporaneous write. That provenance is honest on its face.
4. `EVIDENCE.md:697-699` / `:708-711`, exact strings checked: the R4/R5/R6 rewrites.
   - R4's reading claim is coherent: both joined counters keep raw `products.codprod` as the
     distinct counted key while canonicalizing membership; the 10,584-row database measurements
     are **execution seat**.
   - R5 is true: the nil-check/safety claim is narrower than degradation parity, and the current
     code/test now states and pins the actual `95 / ALTA / ACCEPT` result.
   - R6 is true: `identity_anchor_adapter.go:23-34` computes
     `Supplied = (anchor in declaredSet)`; the known vocabulary includes `marca`, while Mercado
     Livre declares seller_sku/EAN/title only (`capability_adapter.go:79-91`), so the unsupplied
     `marca` branch is live.

Also judged:

- `EVIDENCE.md:527-531` does not claim a clean tree. It plainly reports
  `` M generation_service.go`` together with an empty content diff/stat-cache explanation.
- Characterizing round-2 side B's `Findings: None` as a coverage failure is fair.
  `dispatches/p6-sol-gate-r2.md:1-19` checks the brief-named `:418` site and stops; its findings
  section says none despite the false second site and false degradation sentence documented by
  side A. It is not corroboration of claims it never examined.

## Findings

### Blocking 1 — false zero-product unreachability claim

- `generation_service_test.go:400-403`
- Wrong:
  ``generation_service.go's findProducts drops candidates without one before any candidate exists: a `ProductCandidate{}` is unreachable in production``.
- Why: `generation_service.go:215` constructs the production unresolved candidate with exactly
  `internalreaddomain.ProductCandidate{}`.
- Correct: scope the statement to non-nil matched products returned by `findProducts`; do not claim
  the zero value is unreachable across candidate generation.

### Blocking 2 — capped-run comments overclaim absence

- `generation_service_test.go:222-224`
- `generation_service.go:85-89`
- Wrong: `"a capped run leaves every uncapped listing without a candidate"`.
- Why: `link_candidate_repo.go:38-48` deletes candidates only for processed identities. An
  unprocessed listing may retain a previous candidate.
- Correct: a capped run does not generate or refresh candidates beyond the cap and can therefore
  leave missing or stale coverage.

### Blocking 3 — title-only comment contradicts its own branch

- `generation_service_test.go:1636-1637`
- Wrong: `"title is ranking-only and can never grant ACCEPT/REVIEW-grade confidence"`.
- Why: `generation_service.go:558-564` sets title-only to BAIXA/REVIEW, and this test asserts REVIEW
  at `:1654-1655`.
- Correct: title-only never grants ACCEPT; it yields BAIXA/REVIEW.

### Blocking 4 — “only cross-side anchors” is broader than the code

- `generation_service.go:473-474`
- Wrong: `"seller_sku and ean are the only cross-side anchors available against provider data"`.
- Why: title is also compared across provider and ERP sides for ranking and hard-negative vetoes.
- Correct: seller_sku and EAN are the only corroborating/auto-accept anchors; title remains a
  cross-side ranking/veto signal.

### Blocking 5 — reason-promotion comment promises retention the code does not provide

- `generation_service.go:698-703`
- Wrong: `"promote only its classification and retain that sentence for the operator"`.
- Why: the seeded detail is retained only for INCOMPARABLE. On UNAVAILABLE, the whole declared
  reason including detail replaces it; EVIDENCE R6 correctly notes this branch is live for `marca`.
- Correct: state the conditional behavior: retain seeded detail for INCOMPARABLE, otherwise use the
  provider-declaration reason.

## What I could not verify, and why

- All test/build/vet/must-fail/runtime claims: **execution seat**.
- The 10,584-row database measurements and all real-Postgres fixture outcomes: **execution seat**.
- HEAD identity, commit boundaries/atomicity, status/stat-cache facts, restoration diffs, and git
  history: **execution seat**. I did not run git.
- The patch SHA-256 values were readable without execution and matched the brief:
  `4943e3c3f9c7e4972641b79c90e2345775a17016b35c68b102cb3f91cd37cbd8` and
  `085882a6085546666a2aa56adfb683ca3d5c551f73c00c3d406d62b6be2a7444`.
- Transcript authenticity beyond the salvaged artifact's internal provenance header: no independent
  transcript source was in this review scope.
- Out of scope and not graded: A2-R1/G4 AGAINST/sargability, B-02, B-08, `apps/web` tsc, L2/live
  drive/`LIVE-VERIFIED:`, and the `ltrim(x,'0')` all-zero collision.
