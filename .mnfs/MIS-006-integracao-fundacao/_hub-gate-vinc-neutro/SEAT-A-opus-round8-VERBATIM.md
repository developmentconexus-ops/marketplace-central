# ASSENTO A — Opus frio, adversarial — CHIP-VINC-NEUTRO rodada 8 — **VERBATIM**

`colado no RETORNO, antes de qualquer análise` (§11, emenda de `0bda36bb`: a rodada 7 perdeu o
texto cru deste assento porque a janela entre retornar e colar atravessou uma compactação).

`tip julgado: bfc1d9bb` · `entrada: code-diff-r8.patch + GATE-P6-r8.md + DRIVE-EVIDENCE.md + lane-r8/`

---

SEAT-VERDICT: BLOCKED

| Severity | Criterion | file:line | Defect | Concrete failure scenario |
|---|---|---|---|---|
| BLOCKING | outside V8-1..5, raised under the RULER (wrong data on screen + claim the code contradicts) | `apps/web/src/pages/vinculos/QueueRow.tsx` new-file ~119-137 = `code-diff-r8.patch:104-122` (`providerDisplayName` / `restored = painted.toLowerCase()...`) | The round-trip is **not** an injectivity check across CASE. It lower-cases `painted` before comparing, so a lower-case code is typeset while a code differing only by case falls through verbatim — and both land on the same painted string. The docstring one screen above asserts the opposite: "the transform is applied only where it is INJECTIVE" and "Anything else — hyphens, **mixed case**, embedded spaces … renders verbatim" (`patch:86-91`). `amazon-marketplace` also does not render verbatim; it renders `Amazon-marketplace`. | Two providers registered at once — `provider_code = "amazon"` and `provider_code = "Amazon"` (legal simultaneously by the diff's own cited fact, `registry.go:100-114` dedupes by exact string equality, `patch:80-83`). Row A: `typesetSlug("amazon")="Amazon"`, `restored="amazon"===code` → prints **"Amazon"**. Row B: `typesetSlug("Amazon")="Amazon"`, `restored="amazon"!=="Amazon"` → prints **"Amazon"** verbatim. Canal column shows one name for two marketplaces — exactly "two providers wearing one name is wrong information" (`patch:83-84`). The guard test `does not let two provider codes collapse onto one name through whitespace` (`patch:1399-1421`) never feeds a case variant, so it stays green. |
| REPORT | V8-2 (second direction) | `wireFixtures.ts:~450` = `patch:2971-2972` vs `patch:2965` and body `patch:2815-2821` | Declared hole slightly wider than the real one: "does NOT check … the direction/side of an absence on a SUPPLIED anchor" — but `{anchor:"ean", direction:"UNAVAILABLE", side:"erp"}` (ean IS supplied) **throws** at `patch:2815`. Not blocking: the same docstring's list 1 states that exact rule ("no non-INCOMPARABLE reason carries a `side`"), so the two lists together are accurate. | `wireCandidate({reasons:[…{anchor:"ean",direction:"UNAVAILABLE",side:"erp"}…]})` throws `ean carries a side on a UNAVAILABLE reason` where the second list says it would pass. |
| REPORT | V8-4 | `wireFixtures.ts:57` = `patch:2578`; guard `patch:2256-2269` | "whole set of declarations that exist" is TRUE today (hub-measured: `DRIVE-EVIDENCE.md:41`), but the guard only reads `mercado_livre`'s adapter — it cannot detect a *second* adapter landing. Expiration trigger IS named in-tree (`patch:2996-2999`), which is the form V8-4 asks for. | A new `connectors/adapters/shopee/capability_adapter.go` declaring `IdentityAnchors` leaves all 8 guard tests green while `DECLARED_PROVIDER_CAPABILITIES` is stale and its shopee fixtures stay `driftCandidate`. |

Criteria discharged:
- **V8-1 PASS.** "Rows of exactly that shape are on the screen today" (`patch:390`) is witnessed: `DRIVE-EVIDENCE.md:24-26` show three live rows `[(seller_sku,FOR,) (ean,INCOMPARABLE,both/provider) (marca,UNAVAILABLE,)]`, and `:85` records old `[FOR, UNAVAILABLE]` → new `[FOR, INCOMPARABLE]`, muda **SIM**. Existential claim, existential witness.
- **V8-2 list 1 PASS.** All three bullets enforced: capability + per-unsupplied-anchor sentence `patch:2791-2799, 2854-2872`; vocabulary/two-absences/signal+absence/side `patch:2806, 2815, 2831, 2839`; tuple AND signal set `patch:2891-2919` (`sameSignals`, multiplicity-sensitive).
- **V8-3 PASS.** Pair present and arm-named: RED `× MUST-PASS — \`title\` FOR and \`title\` AGAINST … stays accepted` under the hoisted (wide) dup rule, `lane-r8/V-signals-round8-arms-RED.md:62`; GREEN `-t "MUST-"` → `wireFixtures.guard.test.ts (8 tests | 3 skipped)`, `6 passed | 0 failed`, `lane-r8/V-mustfail-round8.log:66,77`.
- **V8-5 PASS.** `550 passed (550)` / `67 passed (67)` (`V-signals-round8-arms-GREEN.log:88-89`); tsc 12 errors, every one outside `pages/vinculos` (`V-tsc-round8.log:1-16`), and the diff touches only `apps/web/src/pages/vinculos/*`.

Limits of this review (read-only, no Bash, no browser):
1. I could not run anything: RED/GREEN, md5 restore of `wireFixtures.ts`, and the mutation transcript in `arms-RED.md` are chip-authored and self-certified — the arms-RED file is a transcription, not a captured `.log`.
2. I could not verify V8-4's "single adapter" by string in the tree (read-set restricted to the four inputs); I took `DRIVE-EVIDENCE.md:41` as the hub's measurement.
3. Reachability beyond the nine witnessed rows is NOT-EVIDENCED by construction: 13 of 16 `PRODUCIBLE_SITES` are unwitnessed (`DRIVE-EVIDENCE.md:103`), and I did not treat the drive as wire coverage.
4. `PRODUCIBLE_SITES` line cites into `generation_service.go` are unverifiable from here — the Go source is outside my read-set, so every `:NNN` in that table is NOT-EVIDENCED.
5. The blocking row is outside V8-1..V8-5; the hub decides scope, but I will not soften it — it is a named observable producing wrong information on screen.
