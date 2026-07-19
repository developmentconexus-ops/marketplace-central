# F-01 produto-detalhe — Dispatch Ledger

Chip CHIP-M06-PRODUTO · branch chip-m06-produto @ base 89de2fef · contingency lane §12 (codex quota).
impl-pack v1.0.0 (CORE §4). Every dispatch row written at dispatch time; completed with verdict/artifact.

| # | Role | Model | Phase | Prompt/artifact | Result |
|---|---|---|---|---|---|
| D1 | Investigator (codebase map) | sonnet | P2-pre | Agent a40993bf (map A–F) | DONE — map returned; GAPs: no per-domain SDK files, verdict_label null (M-07), no reservado/disponível split |
| D2 | Planner (batch plan) | cold Opus | P2 | Agent a39f60a4 → plan.md | DONE — 7 slices, all open_questions:[], verified vs source |
| D3 | Implementer S1 (scaffold+route+mktclient+tabcodec) | sonnet | P3 | Agent a9521d4b | DONE @8170084b — 11 produto + 18 AppRouter green, build green; chip re-ran produto vitest (11 green) |
| D4 | Implementer S2 (header widget) | sonnet | P3 | Agent (S2) | DONE @1d8f9d6e — 13 produto green, build green; chip re-ran (green) |
| D5 | Implementer S3 (veredicto box: verdict + sync collect + refetch, no polling) | sonnet | P3 | Agent aff1030c | DONE @42f3c7c — 18 produto green, build green; chip re-ran (18 green) + grep-confirmed NO polling construct |
| D6 | Implementer S4 (anúncios vinculados tab) | sonnet | P3 | Agent a467e2df | DONE @085de0c — 21 produto green, build green; chip re-ran (21 green) |
| D7 | Implementer S5 (estoque tab) | sonnet | P3 | Agent ad61efb6 | DONE @b0d313e — 27 produto green, build green; chip re-ran (27 green) |
| D8 | Implementer S6 (deep-link + F5 hardening test) | sonnet | P3 | Agent ac46c62e | DONE @07adec7 — no source change (behavior already correct from S1); 31 produto green |
| D9 | Implementer S7 (partial-failure isolation test) | sonnet | P3 | Agent af126c23 | DONE @bad1663 — no widget change (isolation+retry already correct); 34 produto green, build green |
| D10 | Adversarial feature reviewer (whole diff, refute-framed) | sonnet | P4 | Agent a5aea618 | HOLD — 2 blockers: B1 REFFORN dropped when EAN present (C01); B2 fabricated "0 vendedores" on null market_range (ADR-17). No polling / keys match / scope clean confirmed |
| D11 | Review-fix implementer (B1+B2+minor) | sonnet | P4 | Agent a26298cf | DONE @0931c5b — REFFORN always shown; honest "vendedores insuficientes"; Anúncios group fallback tightened; 38 produto green, build green; chip re-ran (38 green) + spot-checked both fixes |
| D12 | P6 gate A (cold, independent) | Opus | P6 | Agent abfba36a | GATE: SHIP @0931c5b — 0 blockers; C01–C05 all clean; no polling/keys/scope confirmed; flagged costFact unwired (minor) |
| D13 | P6 gate B (independent, refute) | sonnet | P6 | Agent a6fd555 | GATE: SHIP @0931c5b — 0 blockers; independently re-derived all checks; flagged costFact unwired (major, non-blocking) |
| D14 | costFact wiring fix | sonnet | P6-fix | Agent a29347c5 | DONE @dd67067 — shared-key getCatalogProduct query feeds costFact into VeredictoBox; stale doc refreshed; 39 produto green, build green; chip re-ran (39 green) |
| D15 | P6 re-gate delta (both gates, 0931c5b→dd67067) | Opus+sonnet | P6 | Agents abfba36a + a6fd555 | BOTH GATE: SHIP @dd67067 — costFact minor/major CLOSED; shared-key/no-blank/honest-fallback/scope/green all re-confirmed; 0 delta blockers. **P6-DUAL-GATE: AGREEMENT** |

## Changed-path reconciliation (per slice: declared / changed-undeclared / declared-but-unchanged)

- S1 @8170084b: declared = {routes/produto.tsx, pages/produto/{ProdutoPage,productQueryState,marketClient}.tsx/.ts + 2 tests, app/AppRouter.test.tsx(granted)}. changed-undeclared = none. declared-but-unchanged = none. AppRouter.test.tsx diff = single assertion+title+path (grant-scope clean).
- S2 @1d8f9d6e: declared = {pages/produto/ProdutoHeader.tsx/.test.tsx, ProdutoPage.tsx compose}. changed-undeclared = ProdutoPage.test.tsx (vi.mock ClientContext, mock-only, no assertion change — owned dir, justified). declared-but-unchanged = none.
- S3 @42f3c7c: declared = {pages/produto/VeredictoBox.tsx/.test.tsx, ProdutoPage.tsx mount}. changed-undeclared = none (ProdutoPage.test.tsx left untouched — market client self-contained via useProdutoMarketClient, no mock needed). declared-but-unchanged = ProdutoPage.test.tsx. costFact passed undefined (UnknownValue) — no 2nd catalog query; deferred to a later thread-through if wanted.
- S4 @085de0c: declared = {pages/produto/AnunciosVinculadosTab.tsx/.test.tsx, ProdutoPage.tsx mount}. changed-undeclared = none. declared-but-unchanged = none. Reused ui DataTable + injectable client/installationId props (VeredictoBox idiom).
- S5 @b0d313e: declared = {pages/produto/EstoqueTab.tsx/.test.tsx, ProdutoPage.tsx mount}. changed-undeclared = none. Reuses ProdutoHeader catalog key ["catalog","product",productId] (shared cache).
- S6 @07adec7: declared = {pages/produto/ProdutoPage.deeplink.test.tsx}. changed-undeclared = none. ProdutoPage.tsx/productQueryState.ts NOT touched (behavior already correct). Test harness wraps InstallationProvider (Anúncios panel needs useInstallation).
- S7 @bad1663: declared = {pages/produto/ProdutoPage.partialFailure.test.tsx}. changed-undeclared = none. No widget change (per-widget isolation + ErrorState refetch retry already correct S2–S5).

## Full changed-path reconciliation vs base 89de2fef (feature close)

git diff --name-only 89de2fef..HEAD (source, excl .mnfs evidence) — see chip verification below.
Expected ⊆ { apps/web/src/pages/produto/**, apps/web/src/routes/produto.tsx, apps/web/src/app/AppRouter.test.tsx (granted, 1 assertion) }.
