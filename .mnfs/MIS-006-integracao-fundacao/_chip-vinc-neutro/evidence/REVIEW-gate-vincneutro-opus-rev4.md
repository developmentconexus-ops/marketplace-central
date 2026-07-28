> **PROVENANCE — transcribed, not captured.** This seat is `harness:gate-reviewer` (Read/Grep/Glob),
> which has no `Write` by construction, so its verdict cannot land on disk by its own hand. It is
> reproduced below from the task notification, in full and in the seat's own order, with nothing
> added, dropped, reordered or summarised.
>
> **One mechanical transform, declared:** the transport encoded `<` and `>` as `&lt;`/`&gt;` inside
> code spans. Those are decoded back to `<`/`>` here. That is the only edit; it restores what the
> seat wrote rather than changing it. Every other byte is as received.
>
> Seat: cold Opus subagent, read-only, background. Dispatch tip `2e9e9ce3`. Round 6, OPUS side.
> Blind to the GPT side.

---

REFUTED

---

## Findings (most severe first)

### F1 — MEDIUM-HIGH. `wireCandidate` certifies ACCEPT fixtures the sole ACCEPT producer cannot emit, and one of them is the fixture that gates the V6 golden assertion

`wireFixtures.ts:202` documents the contract: *"A candidate the backend can actually emit. Throws if it is not one."*

Measured against primary source:

- `LinkCandidateMatchStatusAccept` has exactly **one** non-test assignment — `generation_service.go:505` (grep across `apps/server_core`: 11 hits — 1 assignment, 3 comparisons at `resolution_service.go:231`/`:814` and `generation_service.go:239`, 7 in `_test.go`). The pack's REPORT-2 claim holds.
- That assignment sits in `buildConcordantCandidate`, whose candidate is built at `generation_service.go:491`:
  `newCandidate(snapshot, domain.LinkCandidateStateExactSKU, domain.LinkCandidateMatchInputSellerSKU, skuMatches.InputValue, product, now)`
  and `newCandidate` assigns `State: state, MatchInput: matchInput, MatchValue: strings.TrimSpace(matchValue)` (`:417-419`). The ACCEPT branch `:502-506` touches only Confidence/Band/MatchStatus.
- ⇒ **every** ACCEPT has `state="exact_sku"`, `match_input="seller_sku"`.

Two fixtures assert the opposite, both through `wireCandidate` (aliased at `QueueTab.test.tsx:30` — `const candidate = wireCandidate;`):

- `QueueTab.test.tsx:701-704` — `state: "exact_ean"`, `match_status: "ACCEPT"`, `match_input: "ean"`, `match_value: "7890000000001"`. **This is the row whose `identificado-por` is asserted `"CODPROD + EAN"` at `:746`** — the V6 positive case.
- `QueueTab.test.tsx:553-555` — `state: "exact_ean"`, `match_status: "ACCEPT"`, `match_input: "ean"`, under a comment at `:550-552` that cites `buildConcordantCandidate (:503-505)` by name.

`assertProducibleScore` encodes exactly this status⇒shape rule for NO_CANDIDATE (`wireFixtures.ts:164-169`) and none for ACCEPT. It is **not** inside the declared gap list at `wireFixtures.ts:60-64`, which names only detail wording, the FOR/AGAINST set vs `state`, and title suppression. (The two ACCEPT fixtures missing the `ean` FOR at `:61-68` and `:598-605` *are* inside that declared gap; I am not counting them.)

The pack's own REPORT-2 (`EVIDENCE.md:590-604`) states the invariant these violate — *"any second route to `ACCEPT` makes the column assert a corroboration that did not happen"* — while the fixture backing that column is itself a fabricated second route.

**What ships:** nothing visibly breaks. `decidingAnchors` ignores `match_input` on the ACCEPT branch (`QueueRow.tsx:175`). What breaks is the guarantee: this is the fifth occurrence of the class the mechanism was built to end, now inside the mechanism's own output.

### F2 — MEDIUM. V5's "12 baseline errors, one by one by path" lists 10; rows 9 and 10 are Dispatch-ledger rows pasted into the table

`EVIDENCE.md:127-141`. Row 9 reads `**Round 6 delta gate, GPT side** (code delta 4c000a04..7b5c18eb…` and row 10 `**Round 6 delta gate, OPUS side**…` — verbatim duplicates of `EVIDENCE.md:625` and `:626`.

The two real errors they displaced are present in the artifacts:
`L0-tsc-after.txt:13-14` and `L0-tsc-after-worktree-selfcontained.txt:13-14` —
`ProdutoPage.partialFailure.test.tsx(40,45)` and `(41,46)`, `TS2322: Type '"complete"' is not assignable to type 'CanonicalSourceFactQuality'`.

`validation-contract.md:46-48` requires *"os **12** baseline listados um a um, por caminho"*. Ten are listed. `EVIDENCE.md:1071` grades V5 **PASS**.

### F3 — MEDIUM. V10's `mercado_livre` sweep prints 4 where the tree has 12, and one of the 8 it misses is a value, not prose or a map key

`EVIDENCE.md:345-348`: *"returns **four hits and all four are prose or a map key**, none a rendered value: `QueueRow.tsx:51`, `QueueRow.tsx:56`, `QueueRow.tsx:309`, `VinculoDrawer.tsx:180`."*

Same anchor, non-test files under `apps/web/src/pages/vinculos/`, now:

- **POPULATION = 12** — `wireFixtures.ts` ×7 (`:47,:48,:73,:142,:145,:184,:221`), `QueueRow.tsx` ×4 (`:83,:86,:109,:474`), `VinculoDrawer.tsx` ×1 (`:177`)
- **pack's EXTRACTION = 4**, difference **8**, undeclared.

All four cited line numbers are also wrong for this tree. And `wireFixtures.ts:184` reads `provider_code: "mercado_livre",` — a **value assignment**, which the sentence's universal excludes. `wireFixtures.ts` landed in round 5 and is in the write-set; V10 was never re-measured after it.

The underlying property is likely still true (`QueueTab.test.tsx:486` and `VinculosDesign.golden.test.tsx:147` still assert the DOM; `wireFixtures.ts` is test-only) — but that is not what V10 says, and this is the chip's own R-25 rule (invoked at `EVIDENCE.md:713`) not applied to itself.

### F4 — LOW-MEDIUM. `incomparableSideLabels` is the one wire-indexed `Record` in `QueueRow.tsx` whose fall-through DROPS the wire value, and the round-4 sweep that closed the class did not name it

`QueueRow.tsx:212` `const incomparableSideLabels: Record<ProductLinkReasonSide, string>`, indexed raw at `:230`:

```ts
return reason.side ? incomparableSideLabels[reason.side] : undefined;
```

`ProductLinkReasonSide = "provider" | "erp" | "both"` and `side?` is optional (`packages/sdk-runtime/src/index.ts:1061`, `:1066`). A wire `side` outside the union is truthy, the lookup is `undefined`, and `compactChipLabel` (`:311-314`) renders the anchor alone. No crash, no `undefined` on screen — and no signal that the wire carried a side.

Every sibling resolves to the wire word: `bandLabel` `:60` `?? band`; `directionGlyph` `:292` `?? direction`; `anchorShortLabel` `:300` `?? anchor`; `directionRankOf` `:356` `?? UNKNOWN_DIRECTION_RANK`. The file's own doctrine at `:286-288` — *"The wire word is the only honest thing to show (ADR-17)"* — is not applied here. `EVIDENCE.md:835` claims the round-4 grep found *"a **fourth** site the reviewer did not name"* and that *"Both hardened; fall-through is the wire value verbatim."* This is a fifth site in the same file, unhardened. What ships: on `side` drift the operator silently loses the "where do I go fix it" datum that D-B/V4 exists to deliver.

### F5 — LOW. The generator line citations the throwing constructor prints point at a `}`, a blank line, a func signature, and two comment lines

`wireFixtures.ts:43-44` — *"UNAVAILABLE when unsupplied (:700-702), INCOMPARABLE when supplied (:706-708…)"*, repeated in the runtime failure message at `:147`.

Measured: `generation_service.go:700` is `}`, `:701` blank, `:702` the `func classifyProviderIdentityAnchor(` signature; the UNAVAILABLE return is `:703-704`. `:706` is blank, `:707-708` are the *"ListingSnapshot has no provider-side brand field"* comment; the INCOMPARABLE return is `:710-711`. `EVIDENCE.md:108-112` cites `:711`/`:715`/`:723`/`:726`/`:728` and **all five are exact** — so the code file and the pack disagree. The stated purpose is `wireFixtures.ts:25`: *"naming the generator line the fixture contradicts."*

### F6 — LOW. Two comment-level C5 hits

- `ResolvidosTab.tsx:25` still reads `{/* ANÚNCIO ML */}` over a column whose header is `Anúncio` (`:104`) with the F-05 neutrality comment beside it (`:103`). Comment-only, does not reach the operator.
- `VinculosPage.test.tsx:49-50`: *"Removing it is not doc-rot cleanup — **it is not** describing a seam whose removal is already decided…"*. `EVIDENCE.md:958-959` states the identical reason with the opposite polarity: *"the comment is removed because **it describes** a seam whose removal is already decided…"*. As written the in-code sentence negates its own stated reason.

**Verified TRUE and worth recording** (the pack's load-bearing Go claims that hold): `marketplace_capability.go:40-45` = `{seller_sku, ean, title, marca}` matches `KNOWN_IDENTITY_ANCHORS`; `mercado_livre/capability_adapter.go:90` = `{seller_sku, ean, title}` matches `MERCADO_LIVRE_SUPPLIED_ANCHORS`; `generation_service.go:704` emits `"provider não fornece a âncora %s"`; `identity_anchor_adapter.go:28-33` walks all four; `resolveIdentityAnchors` `:149-169` aborts on an unresolved declaration; `VinculosPage.tsx:159` does render `<ImportacaoSection />`; `ProductLinkBatchPreviewItem.status` is 2-member at `sdk-runtime/src/index.ts:1162`.

---

## SWEEP (mandatory section 1)

| # | class | anchor | POPULATION | EXTRACTION | verdict per site |
|---|---|---|---|---|---|
| C1 | union enumerated by string literal | `=== "`, `!== "`, `.filter(`, `switch (`, `? (` | **57 occ / 14 files** | **18** (subject's TS type is a union) | **2 live instances.** `BatchPreviewModal.tsx:44/:46/:94` — `status === "OK"`/`"FAILED"` read three ways that disagree under drift; already filed unfixed by the pack at `EVIDENCE.md:1011-1020`, union verified 2-member, **not a new finding**. `VinculosPage.tsx:39` `confidence_band === "ALTA"` KPI — a fourth "high" band would be undercounted silently; **outside this write-set** (`EVIDENCE.md:421` "Zero VinculosPage.tsx"), routed. The other 16 select a specific member on purpose and stay correct as the union grows. Residual 39 = JSX booleans, `=== undefined`, `.filter(` on arrays, prose. |
| C2 | check WIDER than the fact it names | `.toContain(`, `.includes(`, `.startsWith(`, `.match(` | **30 lines** | **9** (subject is a symbol/identifier) | **Zero unfixed.** The round-6 fix holds: `wireFixtures.guard.test.ts:139` asserts `GO_SEAM[key].extract.test(text)` with the same object the extractions run (`:152`, `:176`), and the `extract` regexes require `knownIdentityAnchors = ` — so `knownIdentityAnchorsXX` no longer satisfies the sentinel. `:192` `expect(declared).not.toContain("marca")` is array-element equality, exact. `QueueTab.test.tsx:284-288` are substring checks on `className`, which is the normal form for a class attribute — noted, not raised. |
| C3 | assertion that PASSES on empty extraction | `matchAll(`, `.match(`, `Array.from(`, `.length` | **35 `.length` occ / 11 files**; 5 `matchAll`/`match` | **8** (extractions without an inherent bound) | **Zero unguarded.** Five zero-extraction assertions in the guard, exactly as the pack claims: `:131` `text.length`, `:150` + `:174` `constants.size`, `:162` `anchors.length`, `:186` `declared.length`, all `.toBeGreaterThan(0)`. `QueueTab.test.tsx:269` has a positive lower bound. `VinculosDesign.golden.test.tsx:296-299` (`Array.from(querySelectorAll("[class]")) … expect(offenders).toEqual([])`) has **no** population bound, but is anchored by `:289` `findAllByTestId("queue-row")` and `:294` `getByText("ALTA")).toHaveClass(…)`, which prove the population non-empty — **PASS with the mitigation named**. |
| C4 | fixture not built through the throwing constructor | `candidate_id:`, `reasons: [`, `confidence_band:` | **57 lines / 59 occurrences** (see reconciliation) | **37** candidate-object sites | **Zero raw** — all 37 go through `wireCandidate` (direct, via the `candidate` alias `QueueTab.test.tsx:30`, or via `base()` `VinculosDesign.golden.test.tsx:49-50`) or `driftCandidate`. **But 2 of the 6 ACCEPT fixtures pass the constructor while contradicting the sole ACCEPT producer → F1.** The class "raw fixture" is closed; the class "fixture the backend cannot emit" is not. Residual 20 = batch preview/apply refs (a different union, different producer), 2 production sites, 1 interface field. |
| C5 | comment/pack sentence stating an unmeasured repo fact | 7-40 hex, paths, `main`, `deleted`, `no longer` | **5 lines** matched by my hex anchor, **2 false positives** (`match_value: "7890000000001"` at `VinculosDesign.golden.test.tsx:57` and `QueueTab.test.tsx:704` — 13 decimal digits are all in `[0-9a-f]`) → **real population 3** | **3 SHA claims + 4 prose claims** | **F2, F3, F5, F6 raised.** The 3 SHA claims (`VinculosPage.test.tsx:40-43`: merge-base `bcab8269`; `5441fe18` an ancestor of it; `45b887b3` deleted the component on `main`) are **NOT-EVIDENCED** from this seat — routed. The one half I could check is **TRUE**: `VinculosPage.tsx:159` renders `<ImportacaoSection />`. |
| C6 | `Record<…>` indexed by a value that can be outside the key type at runtime | `Record<`, `[` on a wire-typed value | **20 lines**, of which **5 are prose in comments** (`QueueRow.tsx:43,:184,:411`, `QueueTab.test.tsx:806`, `wireFixtures.ts:217`) → **15 real declarations** | **10** wire-indexed, all in `QueueRow.tsx` | **9 guarded, 1 not.** Guarded: `:23`/`:31` via `bandLabel`/`bandClass` (`:60`,`:64`), `:74` via `directionClass` (`:296`), `:108` via `if (mapped)` (`:121`), `:160` via `?? []` (`:192`), `:168` via `if (!rule)` (`:191`), `:254` via `?? anchor` (`:300`), `:266` via `?? direction` (`:292`), `:342` via `?? UNKNOWN_DIRECTION_RANK` (`:356`). **Unguarded-to-verbatim: `:212` `incomparableSideLabels` → F4.** `ImportacaoSection.tsx:7,12,21` are CHIP-IMPORT-CHAIN's per `validation-contract.md:119`. |

---

## COUNT RECONCILIATION (mandatory section 2)

Every sweep above prints POPULATION and EXTRACTION. Three declarations where they differ or where the instrument was wrong:

1. **C4, lines vs occurrences.** My anchor ran in content mode, so it returned **57 lines**. The true occurrence count is **59**: two lines carry two matches each — `QueueTab.test.tsx:656` and `BatchPreviewModal.test.tsx:46`, both `approvals: [{ candidate_id: "cand_1" }, { candidate_id: "cand_2" }]`. **These are precisely the two sites the pack's own §11 reconciliation was short by** (`EVIDENCE.md:984-987`), so that correction is independently confirmed here.
2. **C5, my own instrument's false positives.** My hex anchor `\b[0-9a-f]{7,40}\b` matched 5 lines; 2 are the EAN literal `7890000000001`. Real population 3. Printed rather than silently subtracted.
3. **F3 is a reconciliation failure in the pack**, not mine: POPULATION 12 vs the pack's stated EXTRACTION 4, difference 8, undeclared.

**Patterns that returned zero — none.** Every pattern I relied on is shown matching a known member:
- `/var knownIdentityAnchors = \[\]IdentityAnchor\{([^}]*)\}/` → matches `marketplace_capability.go:40`, block `:41-44`.
- `/IdentityAnchors:\s*\[\]ports\.IdentityAnchor\{([^}]*)\}/` → matches `capability_adapter.go:90`.
- `/(IdentityAnchor\w+)\s+IdentityAnchor\s*=\s*"([^"]+)"/g` → matches `marketplace_capability.go:25-28`.
- `match_status: "ACCEPT"` → 6 members in `pages/vinculos/`, all quoted in F1.

The only zero result I obtained is `apps/web/vitest.chip.config.ts` (`Glob vitest*.ts` in `apps/web` returns `vitest.config.ts` only). I did **not** demonstrate that pattern matching a known member elsewhere, so I mark it **unverified rather than clean** — see the routing note below.

---

## Section 3 — pack custody

**no shell — custody routed to the executor seat**

---

## Section 4 — the two adversarial questions

**Q1 — does the button depend on `match_status`, and does the test still prove what it proved?**

The button depends on `match_status` through exactly one branch: `QueueRow.tsx:449` `const noCandidate = candidate.match_status === "NO_CANDIDATE"`, and `:555` selects between `{Criar produto, Ignorar}` and `{Outro…, Vincular}`. `"Vincular"` is present for every status except `NO_CANDIDATE`. `REVIEW → ACCEPT` moves between two non-`NO_CANDIDATE` statuses, so label and presence are unchanged, and all four assertions (`:88`, `:89`, `:92`, `:94` — two tabs, `Vincular`, `MLB1`) are invariant under the change. **The test still proves exactly what it proved before. No finding.**

But the fixture does travel a different path, and none of the difference is asserted: `decidingAnchors` (`:183-193`) now takes `ACCEPT: () => ["CODPROD","EAN"]` (`:175`) instead of `REVIEW: () => []` (`:178`), so the `identificado-por` cell (`:519-530`) goes from `<UnknownValue />` to a rendered `CODPROD + EAN`; `data-match-status` (`:453`) changed; and `AnchorChips` now receives 3 reasons instead of `[]`, so the Motivo cell goes from `<UnknownValue hint="sem sinais de correspondência" />` (`:369-371`) to two chips plus a `+1` toggle. Net: no coverage lost, none gained, name/coverage mismatch unchanged. The observation that matters is that the risk did not move here — it sits in the sibling ACCEPT fixture at `QueueTab.test.tsx:702`, which **is** asserted (F1).

**Q2 — which test fails if the mirror diverges from the finalizer while the anchor vocabulary stays unchanged?**

**NONE.** That is a named gap, and it is the gap F1 fell into.

`wireFixtures.guard.test.ts` reads exactly two facts out of Go: `knownIdentityAnchors` (pattern `:58`, asserted `:169`) and `mercado_livre`'s `IdentityAnchors` (pattern `:64`, asserted `:188`). Everything else in the mirror is unguarded FE data:

- `MARCA_UNAVAILABLE_DETAIL` (`wireFixtures.ts:76`) vs `generation_service.go:704`'s `fmt.Sprintf("provider não fornece a âncora %s", …)` — rewrite the Go sentence and every fixture and every test stays green.
- `PRODUCIBLE_SCORES` (`:90-99`) vs the eight assignment sites — change `60` to `65` at `generation_service.go:539` and nothing on the FE moves.
- Rule 2 (one reason per anchor, `:120-123`), rule 3 (`side` only on INCOMPARABLE, `:125-130`), rule 4 (`marca` always present, `:133-141`), rule 7 (NO_CANDIDATE shape, `:164-169`) — four FE assertions about Go behaviour with no reader on the Go side.

The guard closes the **vocabulary** half of F-07 and leaves the **behaviour** half open. `wireFixtures.guard.test.ts:30-33` and `EVIDENCE.md:919-923` state the deployed-server limit honestly, but `EVIDENCE.md:863` records F-07 as *"CONFIRMED and CLOSED"*, which over-reads it: what is closed is the vocabulary, not the mirror. The ACCEPT shape is exactly a producer-behaviour fact with no reader — which is why F1 exists and why no lane went red on it.

---

## Named and routed to the executor seat (NOT-EVIDENCED here, no shell)

1. **The brief's own safety invariant**, `git diff --stat 7b5c18eb..HEAD -- apps/` → empty. I could not run it. Everything above is therefore a reading of the **working tree**, not of `7b5c18eb`. If that diff is non-empty, F1/F4 need re-siting.
2. **Pack custody** — `git ls-tree -r HEAD -- .mnfs/…/_chip-vinc-neutro/ | wc -l` and `git status --porcelain` on the same path.
3. **The three SHA claims at `VinculosPage.test.tsx:40-43`** — merge-base `bcab8269`; `5441fe18` an ancestor of it; `45b887b3` deleting `ImportacaoSection.tsx` **and** its test on `main`. The last one is a live merge-orphan hazard and only `git` can settle it.
4. **The harness handed me a start-of-session `git status` reading `M apps/web/src/pages/vinculos/QueueRow.tsx`, `M apps/web/src/pages/vinculos/QueueTab.tsx`, `?? apps/web/vitest.chip.config.ts`.** I am **not** reporting a dirty tree: `apps/web/vitest.chip.config.ts` does not exist in the tree I read, so that snapshot is stale relative to my reads. But I cannot prove the two `.tsx` files are committed, and both are code under review. Route it; do not take my silence as clearance.
5. **`EVIDENCE.md:138-139` cites the dispatch tip as `2e9e9ce3`** while the brief freezes the code delta at `7b5c18eb` and lists four commits ending there. Whether `2e9e9ce3` touches only `.mnfs/` is the same `git` question as item 1.

Lanes, the four must-fail arms, and the merge-of-proof are **not mine** per the brief and I make no claim about them.

---

## Verdict

**REFUTED.** F1 is sufficient on its own: the round exists to certify the mechanism, and the mechanism passes a fixture it documents itself as rejecting — in the file, and on the assertion, that gates V6. F2 leaves a contract criterion (V5's "12, one by one, by path") graded PASS on an enumeration of 10. F3 is an undeclared 4-vs-12 count in a criterion the chip's own R-25 rule says must be corrected at source.

Key paths:
- `C:\Users\leandro.theodoro\Documents\marketplace-central\.claude\worktrees\gifted-dhawan-5049f6\apps\web\src\pages\vinculos\QueueTab.test.tsx`
- `C:\Users\leandro.theodoro\Documents\marketplace-central\.claude\worktrees\gifted-dhawan-5049f6\apps\web\src\pages\vinculos\wireFixtures.ts`
- `C:\Users\leandro.theodoro\Documents\marketplace-central\.claude\worktrees\gifted-dhawan-5049f6\apps\web\src\pages\vinculos\QueueRow.tsx`
- `C:\Users\leandro.theodoro\Documents\marketplace-central\.claude\worktrees\gifted-dhawan-5049f6\apps\server_core\internal\modules\product_links\application\generation_service.go`
- `C:\Users\leandro.theodoro\Documents\marketplace-central\.claude\worktrees\gifted-dhawan-5049f6\.mnfs\MIS-006-integracao-fundacao\_chip-vinc-neutro\EVIDENCE.md`
