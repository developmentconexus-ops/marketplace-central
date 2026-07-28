## PERSISTED BY THE ORCHESTRATOR — read this header before the verdict

This artifact was written by CHIP-ANCHORS-3 (the orchestrating session), not by the seat that
produced it. The seat is `harness:gate-reviewer` (Opus), which is physically read-only — it has
Read/Grep/Glob and no Write — so it cannot persist its own verdict. It returns the verdict as its
final message and then exits.

- Seat: `harness:gate-reviewer`, Opus, round 3, task-id `a0a33c7e5a6f2beab`.
- Brief: `dispatches/p6-gate-brief-r3.md`.
- HEAD under review: `0264ba84`.
- Written contemporaneously: the verdict was persisted in the same turn the task notification
  arrived, before its content was analysed. This is NOT a transcript salvage (contrast
  `p6-opus-gate-r1.md`, which is).
- One transport change, and it is the only change: the notification channel HTML-escaped `&` and
  `>` inside code spans (`&amp;`, `&gt;`). Those are decoded back here. No other edit — no
  reordering, no trimming, no summarising.

---

VERDICT: REFUTED

---

# P6 gate — side A (Opus, `harness:gate-reviewer`), ROUND 3 — CHIP-ANCHORS-3

Scope read: `p6-delta-r3-to-r4.patch` (1 file, +24/-3), plus the MANDATORY SWEEP of both whole files, plus `EVIDENCE.md`, `dispatches/p6-sol-gate-r2.md`, `dispatches/p6-opus-gate-r1.md`, and `p6-input-r4.patch` (used only to establish authorship of swept text). Physically read-only: Read/Grep/Glob, no Bash, no git. No `.env*` read, no server, no environment dump.

---

## 1. MANDATORY SWEEP

**What I grepped** (case-insensitive, whole files, not the delta):

- Pattern A (English absolutes, superset of the brief's list): `\b(never|always|only|unreachable|cannot|impossible|no longer|every|all|any|must|guarantee|ensures?|None)\b`
- Pattern B (Portuguese absolutes, because half the comments and every `Detail` string are Portuguese): `\b(sempre|nunca|jamais|nenhum[a]?|todo[s]?|toda[s]?|impossível|inalcançável|única|único|só )\b`
- Pattern C (targeted, to confirm round-2 deletions and to find survivors of the same string): `unreachable for seller_sku|degrades like both of its siblings|is unreachable in`

**Raw hit counts:** `generation_service_test.go` = **39** matching lines (pattern A) + **5** (pattern B); `generation_service.go` = **25** matching lines (pattern A) + **4** (pattern B). Pattern C = **1** hit total (`generation_service_test.go:402`); the two strings round 2 blocked on return **0** hits — deleted, not annotated.

I read every hit against the code it describes. Hits that were fixture names, `t.Fatalf` message text, or `Detail` payloads asserted verbatim elsewhere are collapsed into one row. **Claim sites: 26.** Verdicts:

| # | file:line | The claim (exact string, abbreviated where long) | Verdict |
|---|---|---|---|
| 1 | `generation_service_test.go:222-224` | "A capped generation run leaves every listing beyond the cap without a candidate" | **NOT-VERIFIED** — true of a first run; on a re-run it depends on `ReplaceLinkCandidates` delete-scope, which lives in the postgres adapter, outside the two swept files. Rationale comment, not a correctness claim. |
| 2 | `generation_service_test.go:400-403` | "a `ProductCandidate{}` is **unreachable in production**, and an assertion over it pins nothing (B-01)" | **FALSE as written** — see BLOCKING 2. `generation_service.go:215`, `:340`, `:379` pass literally `internalreaddomain.ProductCandidate{}` on production paths. |
| 3 | `generation_service_test.go:418-422` | "with a product present, side=erp **cannot** arise for seller_sku" | **TRUE.** With a non-nil product carrying a canonical id, `identityAnchorValues` (`generation_service.go:744-759`) returns `productValue = strconv.Itoa(codprod)` ≠ `""`, so `missingMatchedAnchorReason:648` (`product == nil || productValue == ""`) cannot fire, and `classifyProviderIdentityAnchor:729-738` cannot reach the `SideERP` return either. |
| 4 | `generation_service_test.go:469-470` | "it went through findProducts, which filters out **any** candidate without a canonical id" | **TRUE** — `generation_service.go:277-283`. |
| 5 | `generation_service_test.go:475-476` | "which is **impossible**: CODPROD is part of the primary key of erp_import_products" | **NOT-VERIFIED** — a DB-schema claim; I did not open the migrations, and the swept files do not evidence it. Not load-bearing for the delta. |
| 6 | `generation_service_test.go:477` | "the outcome **no longer depends on refforn at all**" | **TRUE as scoped** (subject of the paragraph is "the ERP side of the seller_sku anchor"); note `newCandidate:402-405` still reads `product.ReferenceCode` into `InternalReferenceCode`, so "at all" is broader than the whole candidate. Nit, not a finding. |
| 7 | `generation_service_test.go:496-497` | "refforn is **no longer read on this side**" | **TRUE** — zero `ReferenceCode` inside `identityAnchorValues`. |
| 8 | `generation_service_test.go:505-511` (**delta**) | "**Every case in this table** runs with a product present … side=erp **cannot** arise for seller_sku **here**" | **TRUE** — see (a). |
| 9 | `generation_service_test.go:544` (**delta**) | "the **single call site** (buildExactCandidates) **always** passes a non-nil pointer" | **TRUE** — `buildConcordantCandidate` has exactly one caller, `generation_service.go:303`, passing `&product` where `product := skuMatches.Products[0]` (`:302`). |
| 10 | `generation_service_test.go:546-550` (**delta**) | "nil-checks the product like **both of its siblings** … **the siblings degrade into ABSENCE reasons**" | **PARTIAL** — see NON-BLOCKING 1. |
| 11 | `generation_service_test.go:642-643` | "so the honest side is erp — **the one direction no other seller_sku assertion in this file covers**" | **FALSE** — see BLOCKING 1. Refuted by `generation_service_test.go:1823-1826` in the same file. |
| 12 | `generation_service_test.go:646-648` | "classifyProviderIdentityAnchor **emits nothing** when the product is nil and the listing value is present" | **TRUE** — `generation_service.go:723-727` returns `("", "", "", false)`; `:686-688` then `continue`s. |
| 13 | `generation_service_test.go:925-929` | "the resolve loop **never** runs" on an empty batch / "this engine was **never** wired" | **TRUE** — `resolveIdentityAnchors:150-167` iterates snapshots; the guard is at `:76`, before any work. |
| 14 | `generation_service_test.go:948, 978, 1009, 1040, 1371, 1399, 1658` | `t.Fatalf` message text ("want no write at all", "want none", …) | **TRUE / N-A** — failure prose, matched by the assertion above each. |
| 15 | `generation_service_test.go:1235-1238` | fixture literals (`"LEGACY-ONLY"`, `sku:LEGACY-ONLY`) | **N-A** — fixture data, no claim. |
| 16 | `generation_service_test.go:1285-1287` | "marca must **ALWAYS** appear as UNAVAILABLE (ADR-17) regardless of band/status" | **TRUE as scoped** to the fixtures below it (all `mercado_livre`, where `marca.Supplied == false`). Strictly false for a provider that *does* declare `marca`: `classifyProviderIdentityAnchor:719-721` then yields INCOMPARABLE, which `TestProviderDeclaredUnmodelledAnchorIsIncomparableWithoutSide:678-689` itself pins. Pre-existing, not in this chip's diff. |
| 17 | `generation_service_test.go:1507-1508` | subtest names "listing side only" / "internal side only" | **N-A**. |
| 18 | `generation_service_test.go:1557` | "um humano diz sim — **never** the machine" | **TRUE** — `applySingleAnchorScore:539` gives CONFIRM, never ACCEPT. |
| 19 | `generation_service_test.go:1616-1622` | "the ERP product resolved by the EAN **always** carries a CODPROD"; "seller_sku was the outlier **only** because it read `refforn`" | **TRUE** — filtering at `:277-283` + `uniqueProducts:437-452`; the sibling `ean` case is pinned by `TestExactSKUWithUnmatchedListingEANKeepsSeededEANReason`. |
| 20 | `generation_service_test.go:1636-1637` | "title is ranking-only and **can never grant ACCEPT/REVIEW-grade confidence**" | **FALSE, pre-existing, not this chip's** — see NON-BLOCKING 3. `generation_service.go:559` sets `MatchStatusReview` for `StateTitleMatch`, and the test 17 lines below asserts exactly that (`:1654`). |
| 21 | `generation_service_test.go:1797-1802` | "the operator has to see WHY the matcher had nothing, **never an empty row**" | **TRUE** — `applyUnresolvedScore:634-638` always seeds two reasons. |
| 22 | `generation_service_test.go:1834-1837` | "marca **must appear** as UNAVAILABLE regardless of band/status" | Same as #16. **TRUE as scoped.** |
| 23 | `generation_service_test.go:1855-1858` | "kit/combo, cor, medida, voltagem are **all binding**" | **TRUE** — `detectHardNegative:822-844`, four checks, each returns on divergence. |
| 24 | `generation_service.go:85-89` | "An absent limit means the WHOLE installation, **never** a 20-row default page" | **TRUE for this method** — `:90` passes `input.Limit` raw. (`ListLinkCandidates:184-186` does default to 20, but that is a different entry point and the comment sits on the generator.) |
| 25 | `generation_service.go:121-127` | "Auto-approval runs **only after** the candidates are persisted"; "**Every** corroborated candidate is attempted" | **TRUE** — `:117` then `:124`; the loop at `:129-136` never breaks and joins errors. |
| 26 | `generation_service.go:224-235`, `:347-352`, `:471-477`, `:561`, `:607-610`, `:622-629`, `:748-754`, `:802-805`, `:857-861` | "**Only** the generator … may fill it"; "ACCEPT is reachable **only** from buildConcordantCandidate"; "**every** product in that set is offered for review"; "title ranks only, **never** accepts"; "nunca ACCEPT"; "**never** confirmation"; "**never** a silent empty row"; "A product that reaches this comparison **always** carries a CODPROD … a claim this side can **never** truthfully make"; "**never** collides with the 'M' grade size"; "**normal titles without measurements are never flagged**" | **TRUE**, except the last: `MatchStatusAccept` is assigned exactly once, at `:515`; `uniqueProducts` + Review for every collision product; `hardNegativeSizePattern:805` is case-sensitive so lowercase `m` cannot match. **The last one is FALSE** — `hardNegativeDimension` also keys on bare grade-size letters, so two titles with **no measurement at all** are flagged: `generation_service_test.go:1517-1520` asserts `detectHardNegative("Camisa M", "Camisa G") == true`. Pre-existing, not in this chip's diff (`p6-input-r4.patch` has zero hits for that comment). |

**Sweep result: 26 claim sites, 2 false and chip-authored (#2, #11), 1 partial and chip-authored (#10), 2 false and pre-existing (#20, #26-last), 2 NOT-VERIFIED (#1, #5), the rest true.**

---

## 2. (a) The replacement comment is true and is not a smaller universal — **PASS**

Exact string checked, `generation_service_test.go:505-511`:

> `// Every case in this table runs with a product present (product` / `// above always carries a canonical id). The ERP side of seller_sku` / `// is the CODPROD, and findProducts drops any candidate without one` / `// — so with a product present, side=erp cannot arise for` / `// seller_sku here. It arises on the nil-product (unresolved) path` / `// instead, pinned by` / `// TestUnresolvedListingSellerSKUIsIncomparableOnTheERPSide.`

I verified the **scope claim itself**, per the brief, rather than the chip's account of it:

1. *"Every case in this table runs with a product present."* The table has three cases (`:490`, `:498`, `:512`). The runner at `:518-529` builds **one** product literal for all three — `internalreaddomain.ProductCandidate{InternalProductID: canonicalIDPtr(101), …}` (`:523-526`) — and registers it under `"ean:EAN-B01"`, while every case's snapshot carries `EAN: "EAN-B01"` (`:521`). `canonicalIDPtr` (`:139-142`) returns a pointer to `InternalProductID(101)`, so `canonicalProductID` (`generation_service.go:464-469`) returns `ok`. All three cases therefore route `skuMatches=0 / eanMatches=1` → `buildExactCandidates:313` → `buildCandidatesFromProducts:386-393` → `applySingleAnchorScore(&candidate, StateExactEAN, …&product)`. **Product present in every case: TRUE.**
2. *"side=erp cannot arise for seller_sku here."* Two and only two producers of a `seller_sku` reason on this path: `missingMatchedAnchorReason` (`generation_service.go:556` → `:641-655`) — `productValue = "101"`, so the `SideERP` arm at `:648` is unreachable and the outcome is `provider` (case 3) or `UNAVAILABLE` (cases 1-2); and `classifyProviderIdentityAnchor` (`:712-739`) — with `comparison.product != nil` and `productValue != ""` the only reachable returns are `provider`, `both`, or no-emit. **TRUE.**
3. *"It arises on the nil-product (unresolved) path instead, pinned by TestUnresolvedListingSellerSKUIsIncomparableOnTheERPSide."* That test exists at `:649-676` and asserts `Direction == INCOMPARABLE && Side == SideERP && Detail == "seller_sku sem correspondência"` (`:667-675`). **TRUE.** (It is not the *only* such pin — see BLOCKING 1 — but this sentence does not claim exclusivity, so it is not falsified here.)

Nit, non-blocking: *"(product above always carries a canonical id)"* — the product literal is at `:523`, i.e. **below** the comment at `:505`. Positional word is wrong; the substance is not.

**(a) PASS.** The replacement is scoped, and the scope is real.

---

## 3. (b) The four new assertions pin what they claim and could fail — **PASS on the assertions, PARTIAL on the rewritten doc comment**

| Assertion | file:line | Reads the field it names? | Right constant? | Can it fail? |
|---|---|---|---|---|
| `ean` FOR reason | `:573-575` `findReason(candidate.Reasons, "ean", …DirectionFor)` | Yes — `candidate.Reasons`, produced at `generation_service.go:505` and preserved by `appendProviderDeclaredUnavailableReasons:662-670` (`declared` is nil here, so no promotion path can mask it) | Yes — `LinkCandidateReasonDirectionFor` | **Yes.** It fails the moment the scorer degrades into absence instead of corroboration — which is the exact behaviour the pack is being held to. |
| `Confidence == 95` | `:576-578` | Yes — `candidate.Confidence`, set at `generation_service.go:513` | Numeric 95, matching `:513` | **Yes**, and non-trivially: the branch is chosen by `detectHardNegative(snapshot.Title, product.Name)` at `:507` over a **zeroed** product, so it is contingent on `detectHardNegative:818` short-circuiting on the empty internal name. Change that and this assertion flips to 25. |
| `ConfidenceBand == …BandAlta` | `:579-581` | Yes — `candidate.ConfidenceBand`, `generation_service.go:514` | Domain constant, not a literal | Yes (same branch). |
| `MatchStatus == …StatusAccept` | `:582-584` | Yes — `candidate.MatchStatus`, `generation_service.go:515` | Domain constant | Yes (same branch). |

None is tautological: none re-asserts a value the fixture itself supplies. The fixture supplies only `snapshot` and two `productMatchResult{InputValue: …}` values; all four asserted fields are computed by the function under test.

Honest caveat I will not dress up: the last three read **one** `else` branch (`generation_service.go:512-516`), so they are three expressions of a single behavioural pin, not three independent ones. That is redundancy, not falsity, and it is what the pack claimed ("95 / ALTA / ACCEPT").

**Does the test still drive the nil-product path?** Yes. `:559-565` still calls `buildConcordantCandidate(…, newProviderIdentityAnchorComparison(snapshot, nil, nil), now)` and `:567-569` still asserts `candidate.InternalProductID == nil`. The four additions are appended after the existing assertions and change nothing about what is driven.

**The rewritten doc comment (`:543-550`)** — the three-function accuracy question the brief asks:

- *"pins the guard, not a production path: the single call site (buildExactCandidates) always passes a non-nil pointer"* — **TRUE** (one call site, `generation_service.go:303`).
- *"this scorer nil-checks the product like both of its siblings, so it no longer panics"* — **TRUE for `applySingleAnchorScore`** (`generation_service.go:524`, byte-identical guard). **Loose for `applyUnresolvedScore`**, which does not nil-check `comparison.product` at all — it hard-codes `nil` into `missingMatchedAnchorReason` (`:635-636`) and lets `:648` handle it. No deref, so no panic; but "nil-checks" is not what that function does.
- *"It does NOT degrade like them: the siblings degrade into ABSENCE reasons"* — **PARTIAL, see NON-BLOCKING 1.**

**(b) PASS** on the four assertions. The doc comment carries one overstated clause, recorded below.

---

## 4. (c) The pack's newest self-corrections are TRUE — **PASS, with one over-broad companion left standing**

| Row (`EVIDENCE.md`) | Claim | My verdict |
|---|---|---|
| `:700` A11/R5 — the test "fixa exatamente" 95/ALTA/ACCEPT was false; "o corpo só asseria `InternalProductID != nil` e um motivo `seller_sku`" | **TRUE.** The delta's pre-image (`p6-delta-r3-to-r4.patch:35-37,50`) shows the body ending after the `seller_sku` FOR check. The self-denunciation is accurate, and the repair is on the test side (`:573-584`), not on the sentence — which is the correct direction. |
| `:701` A2 (2ª vez) — "comentário e cobertura corrigidos em `54342331`" was false; a second site survived at `generation_service_test.go:505` | **TRUE.** The delta removes exactly `// side=erp is unreachable for seller_sku now, so the only honest` / `// INCOMPARABLE is the provider side: the ANÚNCIO carries no SKU.` at `:505-506`, and `grep "unreachable for seller_sku"` over the module now returns **0**. Note the *older* row at `:691` still asserts the superseded "corrigidos em `54342331`" — it is corrected by `:701` two rows down, so the table reads as a log, not as a live contradiction. Accepted. |
| `:702` gate — `dispatches/p6-opus-gate-r1.md` was cited and did not exist | **TRUE that it now exists, and the header is honest.** `p6-opus-gate-r1.md:3-22` opens `## SALVAGED ARTIFACT — read this header before the verdict`, states "This file did NOT exist when the round-2 gate ran", names the seat's incapacity (no Write), names the source jsonl and task-id, and says in terms: "it is recovery from a transcript, not a contemporaneous write, and it should be read with that provenance." It does not launder itself as contemporaneous. **PASS on honesty of the header.** Whether the extraction is verbatim is NOT-VERIFIED — verifying it means opening a session transcript, which is exactly the artefact that leaked a container password in this chip's own R8; I will not read it. |
| `:708-711` Nota de fase — R4/R5/R6/A11 rewritten mid-flight against `2bed7d9d` | **TRUE and correctly disclosed.** It states the code did not change and the pack text did, and it pre-empts a round-2 verdict quoting stale wording. This is the honest form. |

**Over-correction check** (the brief's specific worry — retracting something that was right):

- `:698` R5 retracts the chip's own earlier "o comentário do guard nil é falso e eu ia deletá-lo". The retraction is **correct**: `generation_service.go:490-491` "Both sibling scorers nil-check this pointer; an unconditional deref here made this the one site that panics" holds for `applySingleAnchorScore:524` and no scorer derefs a nil product. Deleting it would indeed have erased truth. No credit bought that was not earned.
- `:697` R4 retracts the overcount claim against real data. Out of my reach (no DB, no shell) — NOT-VERIFIED, and the pack says it was the hub's measurement, not the chip's.

**Working-tree honesty** (`EVIDENCE.md:527-531`): the pack says `git status` shows ` M generation_service.go` and calls it a stat-cache artefact under `autocrlf`, with byte counts, `git diff` empty, and an explicit "Registrado porque 'árvore limpa' seria alegação falsa". Judged as the brief asks — **on the honesty of the claim, not the git fact** — this is the correct form: it states the dirty marker plainly instead of asserting a clean tree. The byte counts themselves are NOT-VERIFIED (no shell).

**Side B characterisation** (`EVIDENCE.md:486-491`): I read `dispatches/p6-sol-gate-r2.md` before answering. The chip's framing is **fair**. Side B's (a) section quotes `generation_service_test.go:418` and grades it PASS on the merits (`p6-sol-gate-r2.md:5-16`) — which my own sweep agrees with (sweep row #3) — but its `## Findings` section is literally `None` (`:74-76`) and its `## What I could not verify` section (`:78-86`) lists **only execution-shaped items**; it does not declare a limit of coverage over the file. So a false sentence 87 lines below the site it graded went unreported, and the artefact reads as full coverage. "O veredito não estava errado sobre o que leu; estava errado sobre o que não leu" is an accurate description of that artefact. Grading it as a coverage failure rather than as corroboration is **fair**, and I reached that independently.

**(c) PASS.**

---

## 5. FINDINGS — most severe first

### BLOCKING 1 — a chip-authored false universal survives, and it is the entire justification for the test it heads

**`apps/server_core/internal/modules/product_links/application/generation_service_test.go:642-643`**, exact string:

> `// and no ERP counterpart at all, so the honest side is erp — the one direction`
> `// no other seller_sku assertion in this file covers.`

**Refuted by the same file, `generation_service_test.go:1823-1826`:**

```go
skuReason, ok := findReason(candidate.Reasons, "seller_sku", productlinksdomain.LinkCandidateReasonDirectionIncomparable)
if !ok || skuReason.Side != productlinksdomain.LinkCandidateReasonSideERP {
    t.Fatalf("reasons=%#v, want seller_sku INCOMPARABLE/erp sem correspondência", candidate.Reasons)
```

`TestCase8NoAnchorResolvedYieldsZeroConfidenceNoCandidateWithReasons` (`:1803-1832`) is a **seller_sku assertion in this file that covers side=erp**, and it covers it by the *same* mechanism the new test claims as unique: it drives `generateSingle` (`:1812`) → `GenerateLinkCandidates` → the production unresolved call at `generation_service.go:216`, with `ProviderCode: "mercado_livre"` (`:1807`) so `mercadoLivreIdentityAnchorReader` declares `seller_sku` `Supplied: true` (`:192`) — i.e. it *also* proves the seeded `side=erp` survives the declared-anchor pass, which `EVIDENCE.md:144-148` gives as the reason the new test had to be separate. Its fixture (`SellerSKU: "SKU-UNKNOWN"`, empty matcher, `:1808-1810`) is behaviourally the same listing as `SKU-SEM-PRODUTO` at `:654`.

**Authorship is this chip's.** `p6-input-r4.patch:741-753` shows the whole doc comment and test added by this chip (`+` lines). `TestCase8` is **not** in the chip's cumulative patch at all — `grep "MLB-FX8"` over `p6-input-r4.patch` = 0 hits, not even as context — so it predates the chip and was never removed. And its `side=erp` outcome is independent of B-01: on the unresolved path `product == nil`, so `missingMatchedAnchorReason:648` fires on the `product == nil` half regardless of whether `identityAnchorValues` reads `refforn` or CODPROD. It asserted `erp` before the chip and after it.

**The same falsehood is in the pack**, `EVIDENCE.md:134`:

> "Efeito líquido antes do conserto: a única asserção de `seller_sku`/`side=erp` da suíte tinha sido removida"

Not the only one. `TestCase8` at `:1823-1826` was never touched.

**Why blocking, not a nit:** this is the defect class this round exists to eliminate — total in wording, partial in code, refutable by one grep of the file the criterion exists to protect (R-24). It is not a stray adjective: it is the stated reason a new test was added and the stated measure of lost coverage in the pack. Correct would be a claim that is true — either the coverage was already there and the new test is corroboration (say so), or name what it adds that `TestCase8` does not.

### BLOCKING 2 — a second chip-authored universal, refutable by grepping the production file for the literal it declares unreachable

**`generation_service_test.go:400-403`**, exact string:

> `// Every product fixture below carries a canonical id, because`
> `// generation_service.go's findProducts drops candidates without one before`
> `// any candidate exists: a \`ProductCandidate{}\` is unreachable in`
> `// production, and an assertion over it pins nothing (B-01).`

`internalreaddomain.ProductCandidate{}` is passed **literally, on production paths**, three times:

- `generation_service.go:215` — `newCandidate(snapshot, …StateUnresolved, …MatchInputNone, "", internalreaddomain.ProductCandidate{}, now)`, the unresolved path that the pack itself (`EVIDENCE.md:131-132`) calls "produção, hoje";
- `generation_service.go:340` — conflict path with no products;
- `generation_service.go:379` — collision path with no products.

and it is **constructed** at `generation_service.go:497` and `:523` whenever `comparison.product == nil`.

Authorship: `p6-input-r4.patch:554-557` — added by this chip.

The charitable narrowing — "a *matcher-derived* candidate lacking a canonical id never reaches this comparison" — is true (`findProducts:277-283`, `uniqueProducts:437-452`) and is what the paragraph needs. But the sentence as committed is total, and the chip has already had two rounds' worth of exactly this shape blocked. It is the same standard the chip applied to itself at `EVIDENCE.md:136-137` ("total na redação, parcial no código"). Correct would be the narrowing, in the text, in the file — the way `:505-511` was repaired this round. I flag the charitable reading explicitly so the hub can adjudicate severity; I will not grade it PASS on a reading the string does not carry.

### NON-BLOCKING 1 — the delta's new contrast about the siblings is total in wording, partial in code

**`generation_service_test.go:547-548`** (added by this delta):

> `// longer panics. It does NOT degrade like them: the siblings degrade into`
> `// ABSENCE reasons, while this scorer degrades into CORROBORATION over a`

Under a nil product, `applySingleAnchorScore` does **not** degrade into absence reasons only. With `state == StateExactSKU` it still emits `{Anchor: "seller_sku", Direction: For, Detail: "seller_sku resolve exato para codprod"}` (`generation_service.go:544-547`) at `70 / MEDIA / CONFIRM` (`:539`) — corroboration asserted with no product — and only its *other*, product-reading reason falls to absence via `missingMatchedAnchorReason`. The difference from `buildConcordantCandidate` is of degree (95/ALTA/ACCEPT vs 70/MEDIA/CONFIRM), not of kind.

Two mitigations, which is why this is non-blocking and not blocking 3: (i) under the **reachable** reading — product present but its EAN/name empty — the sibling's degradation genuinely is absence-only, and that is the only degradation reachable in production, since the `default` arm (`:565-567`) routes nil to `applyUnresolvedScore` and the sole caller `buildCandidatesFromProducts:390` always passes `&product`; (ii) the clause is true of `applyUnresolvedScore` (`:634-638`, pure absence, confidence 0, NO_CANDIDATE). Correct would be naming which degradation is being compared.

### NON-BLOCKING 2 — "nil-checks … like both of its siblings" is loose for one of the two

`generation_service_test.go:546` and `generation_service.go:490`. `applySingleAnchorScore:524` does carry the identical `if comparison.product != nil`. `applyUnresolvedScore:630-639` does not check anything — it passes a hard-coded `nil` (`:635-636`) and relies on `missingMatchedAnchorReason:648`. The conclusion ("no panic") holds; the mechanism named does not, for one of the two.

### NON-BLOCKING 3 — pre-existing false universals found by the sweep, NOT this chip's

Reported because the sweep is over the whole files and silence would read as "all true". Neither is in `p6-input-r4.patch` (0 hits), so neither is chargeable to this delta:

- `generation_service_test.go:1636-1637`: *"title is ranking-only and can never grant ACCEPT/REVIEW-grade confidence"* — `generation_service.go:559` assigns `MatchStatusReview` for `StateTitleMatch`, and the test directly below asserts it (`:1654-1656`). Only the ACCEPT half is true (`:561` "nunca ACCEPT" is the accurate form).
- `generation_service.go:858-861`: *"normal titles without measurements are never flagged (… no false positives)"* — `hardNegativeDimension:876-878` also folds bare grade-size tokens into the signature, so two titles with no measurement at all are flagged; `generation_service_test.go:1517-1520` asserts `detectHardNegative("Camisa M", "Camisa G") == true` as intended behaviour.

---

## 6. Pack markers

| Marker | State |
|---|---|
| `P6-DUAL-GATE:` | Correctly **absent** — it is the hub's line, being produced now. I do not write it and did not find it fabricated anywhere in `_chip-anchors-3/`. |
| `LIVE-VERIFIED:` / `LIVE-WAIVED-BY-OPERATOR:` | **Absent from the whole pack.** The contract excludes L2 from this chip (`EVIDENCE.md:676-681` records it, sourced to my own round-1 seat), so I do not charge it to the chip — but the milestone does not close on this pack as it stands. Unchanged from round 1. |
| `EXEMPLO-IO` golden case asserted by a test | **Absent.** No hit anywhere in `_chip-anchors-3/`. Unchanged from round 1. |

---

## 7. What I could not verify, and why

Leaving this empty would itself be a false claim of coverage.

**Execution seat (the hub's, not mine — I have no Bash and no git):**
- That `TestConcordantCandidateDoesNotDerefNilProduct` actually **passes** with the four new assertions. I verified they are consistent with `generation_service.go:507-516` and that they are falsifiable; I did not run them. In particular `Confidence == 95` depends on `detectHardNegative("Produto Nil", "")` returning false via the `in == ""` short-circuit at `:818` — read, not executed.
- The four must-fail signatures quoted at `EVIDENCE.md:517-521` (`confidence=90`, `confidence_band="MEDIA"`, `match_status="CONFIRM"`, the missing `ean` FOR). Pasted output; not reproduced by me.
- `go build ./...`, `go vet ./...`, `-count=1/-count=10` ladders, the D-121 policy suite, `vitest`, and every integration/Postgres criterion (A5/A6/A7/A10).
- Every git fact: that HEAD is `0264ba84`, that `generation_service.go` is byte-identical to HEAD despite the ` M` marker, the 41615-byte counts, the empty `git diff` after must-fail round-trips, and the sha256 of both frozen patches. I read the patch files' **content**; I cannot hash them.

**Read-scope limits of this seat:**
- Whether `dispatches/p6-opus-gate-r1.md` is a **verbatim** extraction of the round-1 return. Confirming it means reading a session transcript (`8f6e3346-…jsonl`), and this chip's own R8 records a worker leaking a container password into a transcript. I judged the provenance **header** honest on its face and stopped there deliberately.
- `generation_service_test.go:475-476` ("CODPROD is part of the primary key of erp_import_products") — I did not open the migrations; nothing in the two swept files evidences it.
- `generation_service_test.go:222-224` (capped run leaves listings without a candidate) — depends on the delete-scope of `ReplaceLinkCandidates` in the postgres adapter, outside the two swept files.
- The real-DB measurements that rejected the R4 overcount (`EVIDENCE.md:572-579`): the hub's data, not reproducible here.
- Out of scope by brief and not examined: the AGAINST branch of A2-R1, G4/index sargability, B-02 (`apps/web`), B-08 (`platform/httpx`), `apps/web` tsc, L2/live drive, and the `ltrim(x,'0')` all-zero collision.

---

**VERDICT: REFUTED** — 2 blocking. (a) PASS, (b) PASS on the assertions, (c) PASS. The delta itself repairs what round 2 named. It fails on the class: the mandatory sweep found **two further chip-authored universals that a one-line grep refutes** — `generation_service_test.go:642-643` (contradicted by `:1823-1826` in the same file, and repeated in `EVIDENCE.md:134`) and `generation_service_test.go:400-403` (contradicted by `generation_service.go:215`, `:340`, `:379`). Fixing the two sites a brief names, for the third round running, is not sweeping the class.
