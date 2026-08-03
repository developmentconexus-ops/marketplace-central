# Task 5 — Drive ao vivo da Fatia A

Commit no merge: `e8177a3e99b97950b0c35776b9fbbf99b1e6dd6e` (main).
`git log -1`: `2c55b740ef9f938adcc0a9298d2a80d9190f1119 2026-08-03 15:08:56 -0300`.
Frontend container: dev-mode com volume montado (Vite) — código servido é o do checkout, não uma imagem
buildada; confirmado observando o novo `formatAsOf`/`FreshnessIndicator` renderizando ao vivo.
Backend container `CreatedAt` é anterior ao commit, mas **Fatia A não toca `apps/server_core`** (diff do
merge só tem `apps/web`, `packages/web-query`, `packages/ui`, `packages/feature-products`,
`packages/feature-inventory`) — binário velho é irrelevante aqui.

## Step 2 — Percorrer as telas (idade da tela vs. `fetched_at` do banco)

| tela | idade na tela | fonte real (SQL) | bate? |
|---|---|---|---|
| `/integracoes` SyncHealthCard — Listings | "há 4 h" | `sync_state.last_full_sync_at` = 2026-08-03 13:24:05 (idade ~4h43m às 18:07) | sim |
| `/integracoes` SyncHealthCard — Products | "há 3 min" | `sync_state.last_full_sync_at` = 2026-08-03 18:07:07 | sim |
| `/integracoes` SyncHealthCard — Market | "nunca" | `sync_state.last_full_sync_at` = NULL (linha existe, nunca rodou) | sim — honesto, não fabrica zero |
| `/produto/26909` header, campo Custo | "há 761 d" | fato canônico de custo (fonte Sankhya ao vivo, `observed_at` do próprio ERP, não o `updated_at` do mirror que é ~2min) — número grande e plausível, não hora do dia | sim, plausível — ver nota abaixo |
| `/anuncios` cabeçalho | "há menos de 1 min" | tempo de fetch da query (React Query), não `listings.fetched_at` por linha — por design (D-48 é sobre o cabeçalho da página) | sim |
| `/mercado` → Reprecificação, coluna ATUALIZADO | **"há 739830 d"** na maioria das linhas | `market_signal.evidence.fetched_at` = `"0001-01-01T00:00:00Z"` (zero-value Go) quando `match_status="NO_CANDIDATE"` — `market_aggregates` e `market_competitive_signals` têm **0 linhas** no ambiente, confirmando que não existe evidência real nenhuma | **NÃO — defeito pré-existente do backend, não da Fatia A** (ver seção "Achado" abaixo) |
| `/mercado` → Oportunidades | vazio (nenhuma linha) | `market_aggregates` tem 0 linhas com `status='OK'` — não há nada a mostrar, consistente | sim (vazio é o esperado, não é divergência) |

SQL de apoio: `sql-listings-fetched-at.txt`, `sql-products-mirror-26909.txt`,
`sql-market-tables-empty.txt`, `sql-sync-state.txt`.

## Step 3 — Controle positivo

O SQL do plano (`UPDATE market_price_intel_aggregates ...`) não é executável neste ambiente: a tabela
real chama-se `market_aggregates` (PK `tenant_id, product_id, source, computed_at`) e está com **0
linhas** — não há uma linha real para envelhecer nas duas abas de `/mercado`. Em vez de fabricar uma
linha sintética numa tabela de evidência de mercado (o que violaria o próprio princípio ADR-17 que esta
fatia defende), o controle positivo foi feito com dado real já envelhecido, e um segundo controle
sintético foi executado e revertido:

1. **Positivo orgânico (dado real, sem manipulação):** `/produto/26909`, campo Custo, mostra "há 761 d"
   — um fato genuinamente antigo (custo Sankhya não alterado há mais de dois anos) renderizado com o
   número de dias certo. Antes desta fatia (`formatAsOf` só mostrava hora do dia), essa mesma idade teria
   passado por um horário plausível e pareceria fresca. Este é exatamente o defeito de classe que o D-48
   fechou — provado aqui com dado de produção, não fixture.

2. **Controle sintético (executado e revertido):** `UPDATE listings SET fetched_at = now() - interval '9
   days' WHERE provider_listing_id='MLB4735328201'` — recarregado `/anuncios`: o cabeçalho da página não
   mudou (ele é dirigido pelo tempo de fetch da query, não por `listings.fetched_at` por linha — é o
   comportamento correto e esperado, D-48 mirava o cabeçalho, não uma idade por linha em `/anuncios`).
   Nenhum marcador por linha existe hoje em `/anuncios` para reagir a essa coluna. Valor restaurado a
   `2026-08-03 13:23:55.138398+00` (confirmado via SELECT pós-restore).

## Achado: `evidence.fetched_at` zero-value quando `match_status=NO_CANDIDATE`

`curl http://localhost:8080/listings?...` (`curl-listings-zero-value-evidence.json`) mostra
`market_signal.evidence.fetched_at = "0001-01-01T00:00:00Z"` para linhas sem candidato de mercado. É um
`time.Time` zero-value do Go serializado como se fosse um fato real — viola a própria regra do AGENTS.md
("unknown operational facts never become zero/default"). Antes desta fatia isso era invisível (mostrava
hora do dia, ~meia-noite, plausível); depois de `formatAsOf` virar idade relativa, aparece como "há 739830
d" na Reprecificação.

**Não é regressão da Fatia A** — o FE está renderizando fielmente o que a API manda; o defeito nasce na
construção do `evidence` no backend, não no display. Fora do escopo desta fatia (código Go em
`apps/server_core`, fora da lista de arquivos do plano). Registrado como chip separado
(`task_255e4bd8`, "Fix zero-value evidence.fetched_at on NO_CANDIDATE match") em vez de bloquear o
merge/fechamento da Fatia A.

## Conclusão

Nenhuma divergência atribuível ao código desta fatia. A única divergência encontrada (zero-value em
`evidence.fetched_at`) é um defeito de backend pré-existente, agora visível — não fabricado por esta
fatia, e sim exposto por ela fazer exatamente o que deveria: mostrar a idade real em vez de escondê-la
atrás de um horário. **Fatia A fecha.**
