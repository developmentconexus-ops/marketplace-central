# P3 Claude candidate — MIS-007 ml-sync (r01, 2026-07-31)

Insumo: design ratificado (`docs/design/MIS-007-ML-SYNC-DESIGN.md`) + P2
(`research/codebase-ingest-side.md`, `research/codebase-read-side.md`,
`research/external-ml-api-facts.md`) + gate P1 (mission.md Interview).

## Architecture spine (ADR-lite)

### ADR-01: Núcleo nativo × adapter (herda design §3)
- Decisão: núcleo agnóstico de provider (Publicação, Custo de canal, Divergência, Envio,
  Inbox, sync_state); adapter ML é o único que conhece endpoints/headers/DTOs ML.
- Prevents: telas/serviços importando tipos ML (repetição da implementação antiga).
- Must preserve: teste da fronteira — provider novo = 1 adapter, zero mudança em núcleo.
- Trade-off: mapeamento adapter→conceito custa structs a mais; aceito.
- Validation impact: revisão de gate verifica ausência de import `connectors` em núcleo/read.

### ADR-02: Read = Postgres, ponto; enriquecimento persiste no INGEST
- Decisão: nenhum handler interativo chama ML. Os 4 sítios morrem. Dados que hoje o read
  busca vivo (shipment, buyer fiscal, tarifa degrau-3) passam a ser persistidos no caminho
  de ingest (import/backfill/scheduler/webhook — classe batch) e lidos do banco.
- Prevents: deadline 15s estourado, ~10.8s de lista, N×3 GETs por render.
- Must preserve: NULL = honesto-desconhecido para campo ainda não ingerido; nunca zero fake.
- Trade-off ACEITO E NOMEADO: entre a Onda 0 e as ondas 1/2 há degradação transitória —
  SLA/frete/rastreio/comprador-fiscal aparecem como desconhecido até o ingest persistir;
  /precos cai para degrau-4 (config, com proveniência visível) até channel_fees camada 2.
  Ordem ratificada pelo operador (design §6). Mitigação: M-02 já move a coleta de shipment
  para o caminho do import (batch), então /pedidos re-ganha os campos na mesma milestone.
- Validation impact: critério "zero call ML no read" medido por teste que instrumenta o
  transport + live-drive <2s.

### ADR-03: Ingest único idempotente; webhook e scheduler = 2 portas do mesmo caminho
- Decisão: um serviço de ingest por entidade (listings, orders) com upsert por chave natural
  (padrão `writer.go:74-95` ON CONFLICT); scheduler, backfill, refresh manual e worker de
  webhook chamam O MESMO serviço. Cursor autoritativo via sync_state (`JobFunc` existente).
- Prevents: fontes divergentes (doença da implementação antiga); duplicata por replay.
- Must preserve: reprocessar notificação/página = zero duplicata; cursor retomável.
- Trade-off: none.
- Validation impact: must-fail de idempotência (reingerir mesmo resource, contar rows).

### ADR-04: 3 camadas de fee com proveniência; degrau-4 config sobrevive como fallback nomeado
- Decisão: `channel_fees` nova com (camada, origem, coletado_em); camada mais forte ganha;
  auditoria 3→2 gera divergência (nunca sobrescreve). Correção P2 sobre o design: o seed
  16/22 já morreu (0081); o que esta missão mata é o degrau-3 VIVO (sítios C/D) e a falta de
  proveniência — `pricing_tariff_defaults` (13/16, editável pelo operador) permanece como
  última camada com proveniência `config`. Frete NÃO tem camada 1 (fato verificado).
- Prevents: número sem origem; sobrescrita silenciosa de tarifa.
- Must preserve: schema das 3 camadas nasce completo; missão popula só 2 e 3.
- Trade-off: /precos pode exibir valor de config defasado quando camadas 2/3 faltam — mas
  com proveniência visível, que é o contrato.
- Validation impact: valor de fee exibido carrega camada+data; divergência 3→2 provada nas
  2 direções.

### ADR-05: Divergência = tabela dedicada única, detectada no ingest, auto-resolve
- Decisão: `divergences(tipo estoque|tarifa, entidade, esperado×observado com timestamps dos
  DOIS lados, detectado_em, resolvido_em, tenant)` — gate P1. Cálculo no INGEST; badge lê flag
  persistida; convergência observada em ingest posterior grava `resolvido_em`.
- Prevents: dois chips inventando shapes; cálculo no read.
- Must preserve: extensível por tipo sem migração destrutiva; timestamps de ambos os lados
  (mitiga falso-positivo R-5).
- Trade-off: join no read de /anuncios (ou view); aceito.
- Validation impact: criar → badge aparece; resolver → some (2 direções).

### ADR-06: Resiliência do client ML centralizada ANTES de qualquer backfill
- Decisão: no adapter (camada `doRawWithHeaders`): backoff exponencial + jitter + honrar
  `Retry-After` no 429; budget de concorrência para hidratação; multiget `/items?ids=` (20)
  substitui N+1 de listings. Shipments seguem GETs individuais paralelos com budget (multiget
  não existe — fato T5).
- Prevents: backfill 12 meses morrendo no primeiro 429 (hoje erro seco).
- Must preserve: uma implementação, todos os caminhos ML passam por ela.
- Trade-off: none.
- Validation impact: teste com 429 injetado (Retry-After respeitado, jitter presente).

### ADR-07: DTO tipado + raw persistido nas entidades de ingest
- Decisão: campos consumidos tipados; payload cru `json.RawMessage` persistido em coluna
  `raw jsonb` da entidade (listings/orders/order_shipments) no upsert do ingest. Raw fica no
  adapter→writer; núcleo não lê raw.
- Prevents: re-fetch para depurar; perda de campo que vira requisito depois.
- Must preserve: cap de tamanho (1MB atual); raw nunca vaza para transport/tela.
- Trade-off: storage extra; aceito (volume: milhares de rows).
- Validation impact: row ingerida tem raw não-nulo; nenhuma resposta HTTP contém raw.

### ADR-08: MASS-CLOSURE morre; ausência = marcação, nunca fechamento em massa
- Decisão: `ApplyCompletedPull` (fecha tudo + re-upserta) é substituído por upsert
  incremental; fechamento de listing só quando o scan COMPLETO terminou e o id não veio
  (padrão keep-absent ADR-04 do mirror: marcar, nunca deletar), ou quando o próprio item
  reporta status closed.
- Prevents: pull parcial fechando o catálogo inteiro (flaw MASS-CLOSURE do audit D-120).
- Must preserve: backfill retomável não pode fechar nada até completar o ciclo.
- Trade-off: listings fechadas no ML podem demorar 1 ciclo p/ refletir; aceito.
- Validation impact: must-fail — pull parcial simulado não fecha listing fora da página.

### ADR-09: Webhook transport fino + inbox; hint não-confiável (gate P1)
- Decisão: `POST /webhooks/{provider}` (classe interativa, 200 imediato) grava
  `notifications_inbox` e retorna; worker in-process consome inbox e chama o ingest (ADR-03).
  Notificação = dica não-confiável: dado real sempre re-buscado na API autenticada; dedupe
  por resource+topic; IP origem gravado, comparado à allowlist oficial só em log. Topic:
  `orders_v2` somente.
- Prevents: processamento no request (timeout ML 500ms-ish); injeção de dado por forja.
- Must preserve: scheduler 5min continua como reconciliação (nunca um sem o outro).
- Trade-off: forja causa fetch extra idempotente; aceito e limitado por dedupe/bound.
- Validation impact: notificação real dispara ingest ≤segundos; forjada não injeta dado;
  replay não duplica.

### ADR-10: Migrações 0086+ aditivas; superfícies novas em tabela nova
- Decisão: numeração a partir de 0086 (0021 duplicado documentado; blocos pré-alocados por
  milestone no P5). Tabelas novas: `order_shipments`, `channel_fees`, `divergences`,
  `notifications_inbox`, `listing_variations`. `listings`/`orders` só ganham colunas
  aditivas. Mirror intocado.
- Prevents: colisão de numeração entre chips; ALTER destrutivo.
- Must preserve: NULL default em coluna nova (honest-unknown).
- Trade-off: none.
- Validation impact: lane de migração + testes de regex existentes.

### Contratos compartilhados a autorar no P4
- OpenAPI+SDK (mesmo commit, regra profile §): `/webhooks/{provider}`, saúde de sync p/
  /integracoes (reusar `GET /sync/runs` OU endpoint novo — decidir no P4), colunas novas nos
  DTOs de listings/orders.
- Interface contracts (`research/*-interface-contract.md`): shape de `channel_fees`,
  `divergences`, `order_shipments`, decomposição JSONB de orders, inbox.

## Milestone headlines (ordem + dependências)

| ID | Headline | Por quê nesta ordem | Depende de |
| --- | --- | --- | --- |
| M-01 | Fundação de resiliência do adapter ML: backoff+jitter+Retry-After, budget de concorrência, multiget items, raw persistível — zero mudança de tela | pré-requisito ratificado das duas ondas; destrava backfills | — |
| M-02 | Onda 0 — orders-decoupling: 4 sítios mortos; `order_shipments` criada e populada no caminho do import (batch); bucket/SLA derivados no ingest; /precos sem degrau-3 vivo (fallback config com proveniência); /pedidos <2s | menor caminho p/ matar a doença read-vivo; re-ganha campos na mesma milestone | M-01 |
| M-03 | Onda 1 — sync de anúncios: backfill scan retomável + multiget, E3 + `listing_variations` + estoque ML, `channel_fees` camada 2 (comissão+frete por anúncio), divergência de estoque, scheduler diário + refresh em lote, re-vínculo pós-backfill, MASS-CLOSURE morto | traz o catálogo ML fresco que /precos e o vínculo consomem | M-01 (∥ M-02) |
| M-04 | Onda 2 — sync de pedidos: backfill 12m (`date_last_updated` desc), incremental 5min, shipments/sla/costs/billing por pedido, pack_id, decomposição persistida + custo congelado, auditoria 3→2, Fila/SLA indexada | verdade absoluta de margem; precisa de resiliência, shipments-schema e camada 2 p/ auditoria | M-01, M-02, M-03 |
| M-05 | Webhook + saúde: `POST /webhooks/{provider}` + inbox + worker (`orders_v2`), registro callback, /integracoes com saúde por entidade + última notificação | fecha o loop "pedido em segundos"; consome o ingest idempotente de M-04 | M-04 |

Paralelismo previsto: M-02 ∥ M-03 (superfícies orders × listings disjuntas; seam compartilhado
= composition root do pricing resolver e `root.go` — lock aditivo a definir no P5).

## Top risks (delta sobre mission.md `## Risks`)
- Degradação transitória da Onda 0 (ADR-02) mal comunicada → operador percebe como regressão.
  Mitigação: M-02 persiste shipment no import na MESMA milestone; critério de validação nomeia
  o estado interino.
- Seam `root.go`/composition root tocado por todos os milestones → colisão de merge.
  Mitigação: ownership matrix P5 + locks aditivos.
- Volume real da conta (F1 dumps) pode ser pequeno → fixtures multi-página obrigatórias.
