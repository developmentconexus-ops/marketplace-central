# ASSENTO B — GPT-5.6 Sol medium — CHIP-VINC-NEUTRO rodada 8 — **VERBATIM**

`colado no RETORNO, antes de qualquer análise` · `tip julgado: bfc1d9bb`
`entrada: code-diff-r8.patch + GATE-P6-r8.md + DRIVE-EVIDENCE.md + lane-r8/` · cego ao assento A

---

SEAT-VERDICT: BLOCKED

| Severity | Criterion | file:line | Defect | Concrete failure scenario |
|---|---|---|---|---|
| BLOCKING | V8-2 | apps/web/src/pages/vinculos/wireFixtures.ts:450 | The docstring says the guard does not check the `side` of an absence on a supplied anchor, but line 294 rejects `side` on every non-INCOMPARABLE reason. The declared hole is therefore false. | A producible `70/MEDIA/CONFIRM/exact_sku/seller_sku` candidate with `seller_sku/FOR`, supplied-anchor `ean/UNAVAILABLE/side=erp`, and the required `marca/UNAVAILABLE` is documented as passing this unchecked dimension, but `wireCandidate` throws that `ean` carries a side on UNAVAILABLE. |
| REPORT | V8-4 | apps/web/src/pages/vinculos/wireFixtures.guard.test.ts:59 | `GO_SEAM` names only the vocabulary file and Mercado Livre adapter; it cannot discover a newly added capability adapter. | Add `connectors/adapters/shopee/capability_adapter.go` declaring `IdentityAnchors`. The guard remains green while `DECLARED_PROVIDER_CAPABILITIES` still omits `shopee`, so `wireCandidate({provider_code:"shopee"})` wrongly rejects a now-producible row. |

V8-1 PASS — DRIVE-EVIDENCE contains three live 70% rows with exactly `seller_sku/FOR`, `ean/INCOMPARABLE`, and `marca/UNAVAILABLE`, discharging "Rows of exactly that shape are on the screen today."

V8-3 PASS — `V-signals-round8-arms-RED.md` names `MUST-PASS — title FOR and title AGAINST` red under the historical wide rule and green under the current duplicate-absence rule; the restored lane reports 550/550 passing.

V8-4 current-total PASS — DRIVE-EVIDENCE states `mercado_livre` is the only capability declaration today; the expiration weakness is reported above.

V8-5 PASS — `V-signals-round8-arms-GREEN.log` reports `67 passed` and `550 passed`; no third-party test assertion failed.

Review limited to the four authorized inputs; no repository tree, chip evidence pack, git, browser, or mutations were used.
The drive witnesses nine rows from 3 of 16 producible sites; I could not verify reachability of the other thirteen.
The TypeScript lane has 12 errors outside `pages/vinculos`; their claimed pre-existence was not independently reproducible under the read-only boundary.
