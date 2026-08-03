# Task 9 — Drive ao vivo da Fatia B (job periódico de mercado)

## Setup
- Merge `fa3-idade-honesta` -> `main` (91a8680f) + fix entity-split `market`/`market_queue` (afb6b54a).
- Backend rebuildado, container StartedAt=2026-08-03T20:18:08Z, "server starting"=20:18:24Z.
- Envelhecida via mutação real do app (botão "coletar agora", produto 26909) e SQL de idade em `computed_at`
  (campo real que `StaleAggregateProductIDs` filtra — não `fetched_at`, que é o nome usado no texto do plano).

## Passo 1-2: ciclo aparece no log
```
backend-1 | 2026/08/03 20:48:31 INFO market collection job: ciclo concluído colhidos=1 falhas=0 teto=50
```
Detectado às 20:48:37Z, dentro da janela esperada (boot 20:18:24Z + 30min = 20:48:24Z).

## Passo 3: nova linha append-only em market_aggregates
```
product_id | source            | status               | computed_at              | age
26909      | ml_catalog_offers | INSUFFICIENT_MARKET  | 2026-08-03 20:48:28 UTC  | 31s
26909      | ml_catalog_offers | INSUFFICIENT_MARKET  | 2026-08-03 16:13:38 UTC  | 4h35m  (linha antiga preservada)
```
2 linhas para o mesmo produto/fonte -> confirma append-only, não overwrite.

## sync_state pós-ciclo (entity separado, prova da migração 0093)
```
installation_id                                          | entity       | last_full_sync_at        | consecutive_failures
inst-mercado_livre-...a0e0                                | listings     | 2026-08-03 13:24:05 UTC  | 0
market                                                    | market       | 2026-08-03 20:48:31 UTC  | 0   <- job periódico
inst-mercado_livre-...a0e0                                | market_queue |  (NULL, fila nunca RecordSuccess) | 0
erp                                                        | products     | 2026-08-03 20:48:51 UTC  | 0
```
`entity=market` (job) e `entity=market_queue` (fila do erp_import) agora são linhas distintas — sem mais
colisão de key React / identidade ambígua no /sync/health.

## Passo 4: tela /mercado (produto 26909) reflete o ciclo
`GET /catalogo/produtos/26909` — texto extraído da página mostra a evidência atualizada:
```
Evidência
erp_cost           2026-08-03T20:49:11.374252624Z
ml_catalog_offers  2026-08-03T20:48:28.043679Z
```
Timestamp bate exatamente com o `computed_at` da nova linha (20:48:28). Confirma que a tela lê o dado
recém-colhido, não um cache. (Screenshot não capturado — bug conhecido do browser pane não compositando
frames nesta sessão; prova por texto extraído é suficiente e consistente com os dados do banco acima.)

## Passo 5: controle negativo — outros schedulers não quebraram
```
erp/products: last_full_sync_at avançou 19:40:32 -> 20:48:51 (seu próprio ciclo, cadência independente)
             consecutive_failures = 0
listings:     consecutive_failures = 0 (inalterado, sem falha introduzida pelo job de mercado)
```
Nenhum dos outros dois schedulers teve `consecutive_failures` incrementado nem parou de rodar por causa
do scheduler de mercado rodando na mesma tenant/pool.

## Resultado
Fatia B (Tasks 6-8) confirmada ao vivo: job periódico roda, respeita teto (50), grava append-only,
sync_state agora distingue fila de job graças ao fix de entity-split (afb6b54a), tela reflete dado novo,
sem regressão nos outros dois schedulers.
