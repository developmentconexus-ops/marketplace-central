# Hub copy — round-6 COLD verdict (Opus, `harness:gate-reviewer`, area SEMANTICS)

**Provenance.** Received VERBATIM in the chip's BLOCKED event. **Hub's own copy**, filed per R-23.
The cold reviewer runs through the Agent tool, whose per-task output file came back EMPTY on all
twelve runs of this chip — so without a transcription this verdict would exist only in a session
transcript, exactly as the round-5 cold verdict did before it became unrecoverable. The chip holds
`_chip-anchors/verdict-r6-cold.md`; this is the second holder.

`gate-reviewer` has no Write **by construction** — that absence is what makes it cold — so the
transcription is necessarily someone else's act, and naming whose act it is, is what keeps it
honest. This copy was written by the hub from the event body.

**Honest limit** (R-23): two holders do not make either copy verifiable against the vanished
original. They remove unilateral control.

- Pack tip reviewed: `a23aee3a` · Code tip: `85b6c367`
- Verdict: **PASS** (cold half only; C5, C9, C10, C11, C12 NOT-CHECKED — no shell)

---

## 1. VERDICT

**VERDICT: PASS** — cold side, area SEMANTICS. No BLOCKING finding stands. One SHOULD-FIX (apparatus recognizer drops anchors silently), one NIT. The verdict covers only the criteria marked PASS below; C5, C9, C10, C11, C12 are NOT-CHECKED and the gate as a whole cannot close on my half alone.

## 2. Per-criterion table (C1–C12)

| # | Contract criterion | Verdict | Proof / what was missing |
|---|---|---|---|
| C1 | Provider **declares** anchors off capability registry; `product_links` consumes via own port + composition wiring | **PASS** | Declaration `connectors/adapters/mercado_livre/capability_adapter.go:90` (`IdentityAnchorSellerSKU, IdentityAnchorEAN, IdentityAnchorTitle`); vocabulary `connectors/ports/marketplace_capability.go:22-42`; service gate `connectors/application/marketplace_capability_service.go:129-155` (nil→`unsupported`, unknown/dup fail closed, returns clone); own port `product_links/ports/provider_identity_anchor_reader.go:3-8`; adapter + compile-time assert `product_links/adapters/connectors/identity_anchor_adapter.go:39`; wiring `composition/root.go:105,387,543` |
| C2 | `mandatoryUnavailableReasons()` gone; no anchor name hardcoded in the generator | **PASS** | Grep `mandatoryUnavailableReasons` over `apps/server_core/internal/modules/**/*.go` → 0 hits. Anchor names now come from the declaration loop `generation_service.go:640-654`. The pack's own disclosure that `"seller_sku"/"ean"/"title"` still appear as reason **seeds** in exactly five functions is accurate: `buildConflictCandidates:324-331`, `buildCollisionCandidates:360-361`, `buildConcordantCandidate:483`, `applySingleAnchorScore:519-537`, `applyUnresolvedScore:609`. Hub already ruled the criterion's prose was the defect; the code state matches the ruling |
| C3 | R2: zero provider branching inside `product_links` | **PASS** | Grep `mercado_livre\|ProviderCode ==\|ProviderCode !=\|switch .*Provider` over `product_links/` → 8 files, **all `_test.go`**. Zero production hits. Adapter carries no provider literal |
| C4 | `UNAVAILABLE` distinguishes "provider doesn't supply" from "supplies it, no value" via `detail` | **PASS** | `generation_service.go:646` `fmt.Sprintf("provider não fornece a âncora %s", anchor.Anchor)`; asserted with two distinct details **in one test** at `generation_service_test.go:136` — `refforn.Detail == "provider não fornece a âncora refforn"`, `ean.Detail == "sem EAN para corroborar o CODPROD"`, plus an explicit `if refforn.Detail == ean.Detail { t.Fatalf }`. Enum still three values: `domain/link_candidate.go:53-55` |
| C5 | R1: OpenAPI + `packages/sdk-runtime` byte-identical to BASE-SHA, `git diff --stat` pasted | **NOT-CHECKED** | Needs `git diff --stat BASE..85b6c367 -- contracts/ packages/sdk-runtime` at the frozen SHAs. No Bash. Derivation side owns it. I did confirm no wire-shape change in the Go types (`{Anchor, Direction, Detail}`, three directions) |
| C6 | D-A: `50cm` vs `500MM` not a contradiction, `50cm` vs `40cm` still is, + must-fail | **PASS** | Exact rationals `normalizeDimensionToken:779-807`, `big.NewRat(127,5)` at `:803-804`, suffix list `:788` with `"m"` last, fail-closed `:795-797`. Tests: golden `TestGoldenToalheiroDimensionUnitEquivalenceYieldsConfirm:800` (asserts CONFIRM/70/MEDIA and **no** `title AGAINST` — the match, not just the parser), `TestEquivalentDimensionUnitsDoNotRejectConcordantCandidate:831`, `TestDimensionCanonicalizationUsesExactMillimetres:855`, `TestDifferentCanonicalDimensionsStillRejectConcordantCandidate:926`, `TestNormalizeDimensionTokenFailsClosed:348`. Must-fail states observed failure modes, not "a proof was done" |
| C7 | D-A regression: hard negatives still block the corroborated path | **PASS** | `auto_link_policy_test.go:274` `TestHardNegativeKindsBlockConcordantSKUAndEAN` covers kit/combo, **cor** (`PUXADOR DHARMA AZUL` vs `PRETO`), voltagem on SKU+EAN-concordant input, asserting REJECT and no automatic link. Case tests `:1175` (kit, EAN-only), `:1206` (voltagem, SKU+EAN), `:1299` (dimension, EAN-only) match the cell's ROUND-4-corrected classification exactly |
| C8 | D-B: independent limits; 29-link fixture returns 29 + must-fail | **PASS** | Three separate `<= 0` defaults, no shared derivation: `resolution_service.go:370-373` (20), `:374-377` (2000), `:378-381` (10000). Tests: `resolution_service_test.go:361` asserts the three literals **and** `len(items) != 29 → Fatalf`; `:406` sweeps limits 5/20/500/5000; `:448` asserts 7/11/13 then 17/11/13. Load-bearing because the limiting stubs `:100`/`:117` record and truncate, unlike the base stub `generation_service_test.go:102` which discards the limit. Must-fail quotes real deltas (`link limit = 20000, want 2000`; `link limit = 100, want 2000`) |
| C9 | R3: no UPDATE/backfill of persisted motivos, no migration | **NOT-CHECKED** | Needs `git diff --name-only` at the frozen SHAs plus a migration-dir count. Adjacent observation only, not proof: the reason-assembly path `generation_service.go:615-657` is in-memory over the candidate being built and I saw no write-back |
| C10 | Non-scope: zero `apps/web/` file in the diff, `--name-only` pasted whole | **NOT-CHECKED** | Needs the actual `git diff --name-only`. Cannot infer file-set membership from reads |
| C11 | Ladder L0+L1 (absolute GOCACHE) + governance from clean detached worktree, 40-hex BaseSha | **NOT-CHECKED** | Needs the Go ladder and governance lane executed. Pack records L0a/L0b/L1/GOV rungs and a 53==53 differential at tip and BASE; verifying that is the derivation side's assignment |
| C12 | `root.go` grant additive-only in tokens, file gofmt-clean | **NOT-CHECKED** | Needs `git diff -w` and `gofmt -l`. What I can state from reading: the grant is exactly three lines — `root.go:105` (import), `:387` (`productLinkIdentityAnchorReader := productlinksconnectors.NewIdentityAnchorAdapter(marketplaceCapabilities)`), `:543` (`IdentityAnchors: productLinkIdentityAnchorReader,`). Whether anything was **removed** requires the diff |

Markers: `EXEMPLO-IO:` present at EVIDENCE.md:17-20 and asserted by a real test (`generation_service_test.go:800`). `LIVE-VERIFIED:` present at EVIDENCE.md:1496, backed by the hub artifact `hub-live-drive-u1-u3.md` with U1–U3 detail. `P6-DUAL-GATE:` is being produced by this round; not fabricated here.

## 3. Findings

**SHOULD-FIX — `cite-table.py` drops unresolvable anchors silently; the table's own scope sentence overclaims.**

Pack, `coordinates.txt` header, states the table is keyed by the anchors EVIDENCE.md names; EVIDENCE.md:1354 says "`coordinates.txt` resolves **every** anchor below to its file and line". Generator:

        locs = symbols.get(name)
        if not locs:
            continue

(`cite-table.py:283-284`), and the footer prints only `anchors resolved 49 / ambiguous 4` — there is **no unresolved count**.

Two distinct silent-drop paths exist, each with a live instance:

1. `TOPLEVEL` (`cite-table.py:~30`) matches only `func`/`type`/`var`/`const` at column 0, so a symbol declared inside a grouped `var (…)` is invisible. `hardNegativeDimensionPattern` is declared at `generation_service.go:684` inside such a block and is absent from `coordinates.txt`, while the prose leans on it in load-bearing sentences — EVIDENCE.md:451 "read `hardNegativeDimensionPattern`", EVIDENCE.md:2073 "a token the alternation `hardNegativeDimensionPattern` cannot emit".
2. `ANCHOR` requires the whole backticked span to be a bare identifier, so `` `ProviderCapabilitySet()` `` (EVIDENCE.md:276) is not even seen as an anchor, though `ProviderCapabilitySet` is a real top-level type at `marketplace_capability_service.go:10`.

This sits against R-17's requirement that the tool report what it cannot resolve. Not raised to BLOCKING because: every anchor named in the C1–C12 criteria cells **does** resolve, the two dropped names I traced turn out to be true claims (`hardNegativeDimensionPattern`'s alternation at `generation_service.go:684` is byte-identical to the pack's quotation; the "one production `ProviderCapabilitySet()` declaration site, the ML adapter" sweep is confirmed — production hits are exactly `capability_adapter.go:79-80` and the registry entry `root.go:385`, everything else is `_test.go`), and the apparatus is formally the GPT side's area. Flagging it because this is the third round in which a coverage artifact was narrower than its own claim, and only the hub can decide whether R-17 makes that blocking a fourth time.

**NIT — Rung LIVE, U2 cell inverts listing and ERP product, reproducing the inversion adjudication A3 already corrected.**

Pack, EVIDENCE.md:1515: "`MLB4735326915` "SOUL TOALHEIRO SIMPLES 500MM CR/POLIDO" went from BAIXA 25%…" — attributes the ERP description to the listing ID. The header and the golden test have it the other way: `generation_service_test.go:800-812` sets `ProviderItemID: "MLB4735326915"`, `Title: "Toalheiro Simples Soul Zen 50cm Cromado"`, and puts `SOUL TOALHEIRO SIMPLES 500MM CR/POLIDO` on the ERP side under `"sku:33698"`. The hub's own live artifact agrees with the code (`hub-live-drive-u1-u3.md:47`, "`500MM` (anúncio) e `50cm` (ERP)" — also inverted there). Cosmetic; no criterion depends on the orientation.

No BLOCKING semantic mismatch found: no test is cited as proving something it does not assert, no guard is described as covering a condition it does not cover, no function is credited with logic living in its caller. I sampled ~25 coordinate rows plus all 20 backticked `Test*` names and every row resolved to the stated file **and** line, with production symbols attributed to production files — the round-5 defect class is not present in what I checked. The four AMBIGUOUS rows (`GenerateLinkCandidates`, `InternalProductID`, `ListLinkWorkflows`, `canonicalProductID`) did not make any sentence unreadable and the header's reasoning against guessing is sound.

## 4. Honest limit of this review

- **Everything requiring execution.** No Bash: no `git diff`, no `gofmt -l`, no `go test`, no ladder, no governance lane, no re-run of the recorded mutations. C5, C9, C10, C11, C12 rest entirely on the derivation side.
- **Must-fail claims judged on form, not outcome.** I checked that each states a specific observed failure mode with a value, and that the quoted test names and assertion texts match the real test bodies at the cited lines. I did **not** and could not confirm those outputs were produced by the runs claimed.
- **Frozen-SHA identity unverified.** The brief's assertion that `85b6c367` and the pack tip `a23aee3a` have no differing file under `apps/`/`contracts/`/`packages/` was taken as given; I read the working tree, not either commit.
- **Coordinate table sampled, not exhaustively re-derived.** ~25 of 53 rows checked against code. A wrong row outside my sample would not have been caught.
- **`coordinates.txt` completeness not audited.** I found the silent-drop path by inspection of the generator plus two spot instances; I did not enumerate every backticked span in a 2187-line pack to bound how many anchors are missing.
- **Scope excluded by the brief:** `slices.SortFunc` stability, `refforn`, `package-lock.json`, and the superseded `cite-audit.*` round-4 artifacts.
- **U1–U3 not re-driven.** Hub-run criteria; I read the artifact, did not exercise the account.
