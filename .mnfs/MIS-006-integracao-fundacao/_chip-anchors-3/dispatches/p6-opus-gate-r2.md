# P6 gate — side A (Claude/Opus, `harness:gate-reviewer`), ROUND 2 — CHIP-ANCHORS-3

Persisted verbatim by the chip session. The seat is physically read-only (no Write/Bash), so it
cannot author its own artifact; it returned this as its final message and the chip wrote it to disk
unedited. Only transport-level HTML entities (`&amp;`, `&gt;`) were decoded back to `&` and `>`.

- AGENT: `harness:gate-reviewer`, model Opus, read-only (Read/Grep/Glob only)
- TASK-ID: `ab0f52d206c131de6`
- HEAD under review: `2bed7d9d`
- SCOPE: delta `p6-delta-r1-to-r3.patch` (4 files, +135/-7); execution criteria explicitly OUT
- BRIEF: `dispatches/p6-gate-brief-r2.md`
- Tokens 139392 / 27 tool uses / 420964 ms

---

VERDICT: REFUTED

## Targets

| Target | Verdict | Proof (file:line, exact string checked) |
|---|---|---|
| (a) false universal stayed deleted, did not become a different one | **FAIL** | The named site is repaired and TRUE — `apps/server_core/internal/modules/product_links/application/generation_service_test.go:418-422` `"The ERP side of seller_sku is the CODPROD, and findProducts drops any candidate without one — so with a product present, side=erp cannot arise for seller_sku."` holds: `generation_service.go:279` `if _, ok := canonicalProductID(product); !ok || ... { continue }` and `applySingleAnchorScore` is reached only from `generation_service.go:390` `newProviderIdentityAnchorComparison(snapshot, identityAnchors, &product)` over findProducts-filtered slices, so `identityAnchorValues:755-758` always yields a non-empty `productValue` and `missingMatchedAnchorReason:648` cannot fire. **But the identical false universal survives 83 lines below**, at `generation_service_test.go:505-506` — see BLOCKING 1 — and a second sentence the pack itself grades false survives at `:541` — see BLOCKING 2. Both are `+` lines of this chip (`p6-input-r3.patch:652`, `:688`) and neither is in the delta. |
| (b) restored coverage covers what it claims | **PASS** | End-to-end: `generation_service_test.go:176` `generateSingle(t, snapshot, &stubProductMatcher{results: map[string][]internalreaddomain.ProductCandidate{}}, now)` → `:1310-1317` `NewGenerationService(...)` + `svc.GenerateLinkCandidates(...)`; empty matcher ⇒ `generation_service.go:215-217` `applyUnresolvedScore(&unresolved, newProviderIdentityAnchorComparison(snapshot, identityAnchors, nil))` ⇒ `:635` `missingMatchedAnchorReason("seller_sku", "seller_sku sem correspondência", comparison.listing, nil)` ⇒ `:648` `case product == nil || productValue == "":` → `SideERP`. The nil comes from production code, not a fixture. Survival: declared anchors are real (`:189-197` `{Anchor: "seller_sku", Supplied: true}`), `appendProviderDeclaredUnavailableReasons:686` calls `classifyProviderIdentityAnchor`, which at `:723-727` returns `"", "", "", false` when `comparison.product == nil` and `listingValue != ""` — the promotion **does not fire**, it does not merely agree. The test discriminates: had it fired with `INCOMPARABLE`, `:700-703` retains `Detail` but overwrites `Side`, and the `reason.Side` assertion at `:188` would catch it. Separate-test-over-table-case reason **holds**: the table runner at `:457` is `applySingleAnchorScore(&candidate, tc.state, newProviderIdentityAnchorComparison(tc.snapshot, nil, tc.product))` — declarations hardcoded `nil`, internal scorer, so a table case could not show survival. |
| (c) the six self-corrections are each TRUE | **FAIL (4 of 6 true, 2 defective)** | Per-row below. |

### (c) row by row

| # | Row | Verdict | Proof |
|---|---|---|---|
| 1 | A2 — "side=erp INALCANÇÁVEL" → reachable via unresolved path | **FAIL (incomplete)** | Truth-half verified (`generation_service.go:216`, `:648`). But `EVIDENCE.md:509` "Comentário e cobertura corrigidos em `54342331`" is false as stated: `generation_service_test.go:505-506` still asserts the retracted claim. BLOCKING 1. |
| 2 | A3 — "dois discriminam", not one | **TRUE** | `a3-mustfail-raw.txt:18` `--- FAIL: TestSellerSKUAnchorReadsCanonicalCodprodNotSupplierReference/listing_has_no_seller_sku`; `:19` `--- PASS: .../listing_has_a_seller_sku,_ERP_product_has_a_refforn`. |
| 3 | A7 — skew diagnosis retracted, REPORT kept | **TRUE, and the retraction is correct** | `impl-integration-a5-a7-a10.md:285-292` measures host `16:09:04.608380400` vs `statement_timestamp()` `16:09:05.208959+00` ≈ +600 ms constant. A constant +600 ms against a host-side `[before, after]` bracket fails every run; the observed result is 5 PASS / 4 FAIL (`:271-276`). The diagnosis does not explain the data; grading it "não provado" while keeping REPORT is honest. Residue: the "~1 ms" magnitude in `EVIDENCE.md:203` appears in no artifact I could find (the log records only "queue_read_at outside [before,after] window"); the retraction survives without it. |
| 4 | A12 — wrong raw artifact, recaptured | **TRUE** | `ladder-l0-l1-raw.txt:1-3` opens straight on the package list, and `:3` is `?   .../cmd/testdb  [no test fi107` — truncated; no `go build`, `go vet` or `count=10` anywhere in it. `ladder-l0-l1-raw-r2.txt:4-15` carries all four with the command echoed above each. |
| 5 | A11 — "degrada como as duas irmãs" → degrades to corroboration | **FAIL (over-claims)** | Substance TRUE: `generation_service.go:503-506` seeds `seller_sku` FOR + `ean` FOR unconditionally, `:513-515` sets `95` / `Alta` / `Accept`, `autoApprovals:239` selects on `Accept`. But the false sentence was left in the tree and the "test pinning it" does not pin it. BLOCKING 2. |
| 6 | ledger — slice reviewer's CONFIRMED counts less | **TRUE** | `review-adversarial-r1.md:19` reads `verified by hand-trace` and grades PASS the very subtest framing that `a3-mustfail-raw.txt:18` refutes. |

### Also in the delta

| Item | Verdict | Proof |
|---|---|---|
| `queued_products` canonical predicate + raw counted key, both CTEs | **PASS on reading** | `query_repository.go:97-98` and `:127-128` now use one identity test; both CTEs count `SELECT DISTINCT products.codprod` drawn from `import_products`, the same rows `importados` counts with `count(*)` (`:133`). `migrations/0046_create_erp_import_products.sql:13` `PRIMARY KEY (tenant_id, protocol_id, codprod)` makes codprod unique inside that CTE, so `count(*) == count(DISTINCT codprod)` and neither joined counter can exceed `importados`. The comment's "the counted key stays the raw codprod, which is what importados counts" is accurate. Fixture arithmetic checks out by reading: importados 3; resolved `ltrim('101','0')='101'=ltrim('00101','0')` → 1; queued `{'101','00102','OUTSIDE'}` → `{'101','102','OUTSIDE'}` matches `'00101'` and `'102'` → DISTINCT 2; under the raw predicate both matches vanish → `Enfileirados:0`, the reported must-fail. Execution is the other seat's. |
| Two-sentence addition above the nil guard | **PASS** | `generation_service.go:491-496` is accurate about the function, traced above; the pre-existing true sentence is preserved. |

---

## Findings

**BLOCKING 1 — the retracted universal is still asserted in the same file.**
`apps/server_core/internal/modules/product_links/application/generation_service_test.go:505-506`:
```go
		// side=erp is unreachable for seller_sku now, so the only honest
		// INCOMPARABLE is the provider side: the ANÚNCIO carries no SKU.
```
This is round 1's blocking claim, unqualified, at a second site. It is refuted by `generation_service.go:648` reached from `:216`, `:341` and `:380` (`applyUnresolvedScore(..., nil)`), and refuted by this chip's own new test at `generation_service_test.go:169-196`, which asserts `Side == SideERP` for `seller_sku`. It is a `+` line of this chip (`p6-input-r3.patch:652`), so it is not pre-existing debt, and the delta did not touch it. R-25 requires the false half DELETED or NARROWED; here it is neither. Correct would be the same narrowing applied at `:418-422`, or deletion. The pack's `EVIDENCE.md:118` "a única asserção de `seller_sku`/`side=erp` da suíte tinha sido removida" and `:509` "Comentário e cobertura corrigidos" both read as total and are partial in code — the exact R-24 shape the contract names.

**BLOCKING 2 — a sentence the pack grades FALSE was left in the code, and a test is credited with an assertion it does not make.**
`generation_service_test.go:541`: `// that this scorer now degrades like both of its siblings instead of panicking.` — `EVIDENCE.md:283-284` says of that claim: «**não é verdade** — `applyUnresolvedScore` degrada para motivos de ausência, este degrada para corroboração». The delta put the true description into `generation_service.go:491-496` and left the contradicting sentence standing, so two comments in the module now contradict each other and the pack says the surviving one is wrong.
Compounding: `EVIDENCE.md:282-283` and `:471` claim `TestConcordantCandidateDoesNotDerefNilProduct` "fixa exatamente isso" / "com teste fixando o comportamento" (95 / ALTA / ACCEPT). The test asserts only `candidate.InternalProductID != nil` (`:558`) and the presence of a `seller_sku` FOR reason (`:561`) — nothing about `Confidence`, `ConfidenceBand` or `MatchStatus`. The chip's own worker says so: `impl-r4-queued-canonicalization.md:429-431` "there is no test that constructs a nil-product comparison and observes 95 / ALTA / ACCEPT". R5 is graded as a reported-and-held hazard when nothing holds it. Correct would be: delete/narrow `:541`, and state R5 as unheld.

**BLOCKING 3 — the pack does not describe the HEAD it is submitted against, and its R4 reports as open a defect the delta closed.**
`EVIDENCE.md:7` `head: 54342331`; the commit table `:24-29` has no row for `2bed7d9d`. `EVIDENCE.md:460-461`: "Além disso `queued_products` **continua juntando cru**, então dois contadores lado a lado na mesma tela usam duas regras de identidade para a mesma coluna" — false at HEAD, refuted by `query_repository.go:127-128`. The new integration test `TestGetImportChainCountsLeadingZeroCodprodInQueue` and the R5 comment sentence appear nowhere in the pack. Also `:464` "para o dono do G4: `ltrim(...)` dos dois lados torna o predicado não-sargable" now understates — it is two CTEs, not one. Correct would be a pack whose criteria and REPORTs are written against `2bed7d9d`.

**BLOCKING 4 — a ledger artifact citation that does not resolve; same class as round 1's blocking 3, and absent from the corrections table.**
`EVIDENCE.md:47` cites `dispatches/p6-opus-gate-r1.md` for the round-1 side-A REFUTED. No such file exists under `_chip-anchors-3/dispatches/` (my glob of `.mnfs/MIS-006-integracao-fundacao/_chip-anchors-3/dispatches/p6-*` returns only `p6-gate-brief.md`, `p6-sol-gate-r1.md`, `p6-sol-gate-r1.stderr.log`, `p6-gate-brief-r2.md`). The only `p6-opus-gate-r1.md` in the repo is `.mnfs/MIS-006-integracao-fundacao/_hub-gate-anchors-2/p6-opus-gate-r1.md`, whose line 1 reads `# P6 gate — Claude/Opus side, round 1 (CHIP-ANCHORS-2, post-merge)` — a different chip. `EVIDENCE.md` attributes at least six substantive claims to "o lado A" (`:85-86` upstream canonicity, `:203-205` the A7 arithmetic, `:244-247` the A9 double-check, `:279-287` the corroboration finding, `:378-382` anterioridade, `:457` R4) that no reader can retrieve. Stated honestly: this may be a hub-side persistence failure rather than the chip's authorship — but at HEAD the citation is unresolvable, and "unwritten = didn't happen" applies to it as it applied to `ladder-l0-l1-raw.txt`.

**NON-BLOCKING 1 — mis-attributed quotation inside a self-correction.** `EVIDENCE.md:283` presents «Meu comentário diz que ele "degrada como as duas irmãs"» in the A11 context, i.e. sourced to `generation_service.go`. That string lives in `generation_service_test.go:541`; the old `generation_service.go` text (`impl-r4-queued-canonicalization.md:364-366`) never contained it.

**NON-BLOCKING 2 — `ltrim(x,'0')` collides all-zero codprods, now in two CTEs.** `'0'`, `'00'`, `'000'` all canonicalize to `''`. The worker declares it untested (`impl-r4-queued-canonicalization.md:441-446`); the delta doubles its blast radius and `EVIDENCE.md` does not carry it. It should be a named REPORT, not only a worker-report line.

**NON-BLOCKING 3 — `EXEMPLO-IO` marker absent from the pack.** No occurrence anywhere under `_chip-anchors-3/`. A golden case does exist as a test (`generation_service_test.go:1327` `TestGoldenToalheiroDimensionUnitEquivalenceYieldsConfirm`) but is not labelled as the golden case. Recorded, not relitigated — `LIVE-VERIFIED:` / L2 are excluded by contract and by this brief.

---

## What I could not verify, and why

- **Every execution-shaped claim.** No shell, no git, by design this round. Not verified by me: the `--- PASS: TestGetImportChainCountsLeadingZeroCodprodInQueue (0.03s)` run and its `Enfileirados:0` must-fail (`impl-r4-queued-canonicalization.md:220-221`, `:257`); the A3 round-trip; the 9 A7 runs; `go build` / `go vet` / `-count=10` 20/20 / 107 packages in `ladder-l0-l1-raw-r2.txt`; `vitest` 4/4; the D-121 suite count. I read the SQL and Go and the arithmetic is consistent with the pasted output — that is corroboration, not proof. Execution seat.
- **Git facts.** That HEAD is `2bed7d9d`; that the delta is exactly commits `54342331` + `2bed7d9d`; that `git diff HEAD` was empty after each must-fail restore; that `055d1705` changed behaviour + OpenAPI + SDK atomically; that `p6-delta-r1-to-r3.patch` and `p6-input-r3.patch` match the tree. I compared the delta patch against the working-tree files and they agree, but I cannot distinguish committed from uncommitted state. Execution seat.
- **Hashes.** `p6-input-r2.patch` sha256 `d06f9881…` (`EVIDENCE.md:347-348`) — read, not hashed. Execution seat.
- **Whether the round-1 side-A verdict exists anywhere outside this worktree.** I searched `.mnfs` and the repo; see BLOCKING 4. If the hub holds it elsewhere, that resolves the citation but not the ledger path.
- **The "~1 ms" failure magnitude** underpinning the A7 retraction arithmetic (`EVIDENCE.md:203`). Not present in any artifact I read. The retraction's conclusion stands on the 9/9-vs-4/9 half alone, so I did not treat it as a finding.
- **Whether any real producer writes an unpadded CODPROD into `sync_state.cursor->'pending'`.** The worker states `grep` finds no non-test caller of `AppendPendingCodigo` (`impl-r4-queued-canonicalization.md:433-440`); I did not repeat that sweep. The fix is defensible as an invariant either way — two counters that disagreed on identity now do not.
- **Round-1 coverage of unchanged hunks**, per the delta-round rule. I did not re-review them; my (a)/(c) findings concern text this chip added and this chip claims to have corrected, not unchanged hunks.
- **Out of scope by brief, not graded:** the `AGAINST` branch of A2-R1, G4, B-02, B-08, `apps/web` tsc, L2 / `LIVE-VERIFIED:`.
