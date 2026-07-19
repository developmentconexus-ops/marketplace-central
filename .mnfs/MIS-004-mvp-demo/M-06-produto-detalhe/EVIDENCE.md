# M-06 produto-detalhe — Evidence Pack

Chip: CHIP-M06-PRODUTO · branch `chip-m06-produto` · base `89de2fef` · HEAD `dd67067`
Contingency lane §12 (codex quota until 2026-07-25) — all roles ran on Claude (cold Opus planner/gate + sonnet implementers/reviewers).

## Merge-gate lines (harness 0.4.0)

- `P6-DUAL-GATE: AGREEMENT`
- `LIVE-VERIFIED: 2026-07-19 hub P7 live-drive on clean docker dev-stack (frontend :5174 + backend :8080 on M-06 worktree mount). /catalogo/produtos/412 (CHAVE COMBINADA, cost 42.1, no market evidence): Veredicto tab honest no-evidence ("sem candidato de mercado — veredicto não avaliável"; Mínimo/Mediana/Vendedores all "—"; real Custo ERP 42.1; evidence source erp_cost; categorical verdict deferred "— M-07" per D-98); Estoque tab FÍSICO=12 real + RESERVADO/DISPONÍVEL "— DESCONHECIDO" (ADR-17, no fabricated zeros); Anúncios vinculados honest empty ("Nenhum registro encontrado" + vincular action). URL tab-state (?tab=estoque|anuncios) deep-links correctly. Theme paper rgb(251,250,247)+Instrument Sans light; a11y tablist/tab/tabpanel/note roles; zero console errors. Populated market path (mediana/vendedores) is data-gated (fixture ERP products all missing_price) — covered by VeredictoBox.test.tsx (182 lines), to be confirmed with real demo product at rehearsal C01. NOTE: reaching the page required an ephemeral hub dev-stack fix — vite proxy key "/catalog" prefix-matches the FE route "/catalogo/*" (dev-only 404 on full-load deep-links); patched to "/catalog/" for the drive, reverted at teardown. Pre-existing infra bug (not in M-06 diff), flagged for profile.>`

## Outcome

`/catalogo/produtos/:productId` real page composed 100% client-side from existing SDK reads (M-01/M-02/M-04/M-05). ZERO new backend / OpenAPI / migration / SDK / governance. Route swap only inside owned surface.

## Slice ledger (P3)

| Slice | Commit | Content |
|---|---|---|
| S1 | 8170084b | route swap + ProdutoPage shell + tablist + local market-client hook + tab URL codec |
| S2 | 1d8f9d6e | ProdutoHeader (identity, honest nulls) |
| S3 | 42f3c7c | VeredictoBox (price-evidence verdict + synchronous collect + refetch, NO polling) |
| S4 | 085de0c | AnunciosVinculadosTab (per-listing position vs mercado) |
| S5 | b0d313e | EstoqueTab (físico honesto, reservado/disponível DESCONHECIDO) |
| S6 | 07adec7 | deep-link + F5 restore hardening test (no source change) |
| S7 | bad1663 | partial-failure isolation test (no widget change) |
| P4-fix | 0931c5b | REFFORN always shown (C01) + honest "vendedores insuficientes" (ADR-17) + Anúncios group-fallback tighten |
| costFact | dd67067 | real Custo ERP wired into VeredictoBox via shared catalog cache key |

## Contract acceptance (M-06-C01..C05)

- **C01** ProdutoHeader: CODPROD/name/EAN/REFFORN/marca/NCM/custo/completude; every null → UnknownValue "—", never 0/fabricated; REFFORN present with or without EAN.
- **C02** VeredictoBox: single synchronous `useMutation(collectMarketPriceIntel)` + invalidate BOTH `["market","verdict",productId]` and `["listings","by-product"]`; button disabled in-flight; **no polling** (no setInterval/refetchInterval, one collect per click); `verdict_label` always null → UnknownValue hint "veredicto de margem — M-07" (hub RULING ratified honest framing, no scope expansion); no R$0/green for no-evidence states; INSUFFICIENT_MARKET exact copy `N vendedores — mínimo 5`, honest no-count when market_range null.
- **C03** AnunciosVinculadosTab: per-listing rank/total, price_to_win, delta_pct, freshness; null/SEM_VINCULO → "—"; empty → EmptyState + /vinculos CTA; never blank/fabricated/cross-product rows.
- **C04** EstoqueTab: físico = stock_quantity value+observed_at; reservado + disponível always UnknownValue DESCONHECIDO (never 0/derived); null value or missing_stock → "importar planilha" (never 0); Concorrência/Pedidos/Histórico absent.
- **C05** ownership clean: changed source ⊆ { apps/web/src/pages/produto/**, apps/web/src/routes/produto.tsx, apps/web/src/app/AppRouter.test.tsx (granted single assertion) }. Verified via `git diff --name-only 89de2fef..HEAD`.

## Verification

- Unit/integration: `npm run test --workspace apps/web -- src/pages/produto/` → 8 files, 39 tests green.
- Shell: `npm run test --workspace apps/web -- src/app/AppRouter.test.tsx` → 18 green.
- Build: `npm run build --workspace apps/web` → green (vite).
- tsc baseline main = 10 pre-existing type-only (D-97) — no new produto type error.

## Reviews

- **P4 adversarial** (sonnet, Agent a5aea618): HOLD → 2 blockers (REFFORN drop, fabricated "0 vendedores") fixed @0931c5b.
- **P6 dual gate @0931c5b**: cold Opus (abfba36a) GATE: SHIP + independent refutation sonnet (a6fd555) GATE: SHIP — AGREEMENT.
- **P6 re-gate @dd67067** (delta after costFact fix): both GATE: SHIP — AGREEMENT holds.

## Hub rulings referenced

- FINDING-1 (verdict_label null): hub VERIFIED, honest framing RATIFIED, categorical margin engine is design-first POST-demo (operator). No M-07 dependency pre-demo.
- FINDING-2 (Estoque KPIs físico + reservado/disponível "—"): hub confirmed matches C04.
- Local per-domain market client inside pages/produto (MarketComparison pattern replica): hub confirmed IN-SCOPE.

Full dispatch ledger: `F-01-produto-detalhe-page/dispatch-ledger.md`.
