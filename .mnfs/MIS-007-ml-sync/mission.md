# MIS-007-ml-sync

```yaml
id: MIS-007
type: mission
status: planned
owner: Mission Strategist
parent: none
created: 2026-07-31
updated: 2026-08-01
validation_level: QA-0
lifecycle_scope: mission
planning_phase: readiness
```

Design ratificado de entrada: [docs/design/MIS-007-ML-SYNC-DESIGN.md](../../docs/design/MIS-007-ML-SYNC-DESIGN.md)
(brainstorm 2026-07-31, toda decisão aprovada pelo operador; fatos verificados no §10).

## Objective

Operação diária confiável do lado ML: o vendedor abre /anuncios e /pedidos de manhã e tudo
está lá, atualizado sozinho, rápido (<2s), com margem real por pedido. Zero botão "refresh"
obrigatório. Missão SÓ LEITURA do ML; desenho nasce write-ready.

## Outcome

- Zero chamada ML no caminho de read das telas (os 4 sítios atuais morrem — Onda 0).
- `listings` completa (E3 + estoque ML + comissão camada 2) sincronizada por
  backfill retomável + scheduler diário + refresh manual. (Frete camada 2: fora de escopo —
  ver Non-Scope; P7 r01 B-1.)
- `orders` 12 meses backfill + incremental 5min + webhook `orders_v2` (segundos), com
  decomposição de margem PERSISTIDA (custo congelado) e `order_shipments` (SLA/custos).
- Divergências (estoque ERP×ML; tarifa realizada×estimada) detectadas no ingest, persistidas
  em tabela dedicada, visíveis como badge ⚠, e provadas nas 2 direções.
- Fee sem proveniência morto. (P2 corrigiu o design: o seed 16%/22% JÁ foi removido na main —
  0081 + `registry_test.go:90-103`; `fee_sync.go:29` não existe mais. O que resta: degrau-3
  vivo nos sítios C/D de /precos morre; degrau-4 `pricing_tariff_defaults` 13/16 vira fallback
  explícito com proveniência `config`.)

## Scope

Ver design §2. As 3 ondas do design sobrevivem como LANES do corte por seam (P3 r01
ratificado 2026-07-31): o decoupling da "Onda 0" acontece DENTRO do milestone que persiste o
substituto (regra replacement-before-deletion, ADR-05), nunca como chip separado que deleta
antes de repor. Pré-requisito das lanes: backoff exponencial + jitter + Retry-After no 429;
regra DTO `Raw json.RawMessage` (com exceção de PII — ADR-03).

## Domain Scope

Capability set = escolha explícita do operador no brainstorm ratificado (design §2, §11) +
gate P1 desta sessão:

- (2/3 Entidades+Lifecycle) listings E3 + listing_variations; orders estendida (pack_id,
  decomposição, bucket); order_shipments; channel_fees 3 camadas; divergences;
  notifications_inbox; sync_state entidades `listings`/`orders`.
- (5 Classificação) categoria = ingest do `category_id` real do anúncio (de graça, sem de-para).
- (6 Audit) proveniência em todo fee (camada, origem, coletado_em); auditoria 3→2.
- (7 Notificações) badge/coluna de divergência em /anuncios (sem central de avisos);
  webhook = mecanismo primário, scheduler = reconciliação.
- (8 Busca/relatório) filtro "divergentes" em /anuncios; Fila/SLA bucket em /pedidos.
- (9 Admin) /integracoes: saúde do sync por entidade + status webhook.
- (10 Integração) webhook topic `orders_v2` SOMENTE (gate P1 2026-07-31); ngrok domínio fixo
  como URL de produção.

## Non-Scope

- Mercado/concorrência, catalog offers, sellers, oportunidades, `ml_tariffs` sweep,
  simulação de produto não vinculado → MIS-008 (destino nomeado no design §2).
- De-para categoria + preditor → missão de publicação (write).
- Writes ML de qualquer tipo → missão futura; write vivo exige autorização do operador.
- Onboarding saga/tela de progresso → futura; `sync_state` já grava progresso.
- Webhook topic `items` → backlog nomeado para missão futura (gate P1: scheduler diário já
  garante completude de anúncios; milestone de pedidos fica menor).
- Relaxar cadência do scheduler pós-webhook-saudável → missão futura.
- Frete no nível de anúncio (camada 2): `api_shipping_options` fica como RESERVA aditiva no
  CHECK de `channel_fees`, NENHUM produtor nesta missão — frete camada 2 = honesto-desconhecido
  (IC-01); frete REAL entra só na camada 3 (custo seller do shipment, por pedido — M-06).
  Destino: missão futura de shipping_options (P7 r01 B-1).

## Current State

MIS-006 fechou a fundação: products_mirror + adapters ERP (xlsx/sankhya), vínculo
(auto-aprovação CODPROD+EAN), sync_state 0075 cadence-agnostic + scheduler,
erro unificado (envelope universal + platform/apierror). Fonte ativa ERP = sankhya
(F-LINK-1 refutado; predicado vendável = 2.923). Lado ML: leitura viva no read
(4 sítios), fees estáticos 16%/22%, sem backfill/webhook/decomposição persistida.
Fatos P2 (base main @ dd89d4b3, detalhe em `research/codebase-*.md`):
- Os "4 sítios" confirmados: `GET /orders` (3 GETs ML/pedido, concorrência 8 — a ~10.8s),
  `GET /orders/{id}` (+buyer fiscal 2-step), `POST /pricing/decompose` e `/solve` (categoria +
  `listing_prices` vivos, por linha da matriz).
- Adapter ML: scan/scroll_id existe; hidratação N+1 (sem multiget `/items?ids=`); 429 = erro
  seco, zero retry/backoff/rate-limit; raw descartado (1MB cap).
- `listings`: variação achatada no PK; único writer = `ApplyCompletedPull` com MASS-CLOSURE;
  pull manual (202 async), sem incremental. `orders`: só `POST /orders/import` (limit 20);
  sem pack_id nem shipments persistidos; decomposição calculada no read, nunca gravada.
- Scheduler: só `products` registrado; entities `listings|orders|market|tariffs` órfãs no
  enum da 0075; `incremental` sempre false.
- Migrações até 0085 (0021 duplicado); próximo número = 0086.
- `GET /sync/runs` existe sem consumidor FE; /integracoes sem seção de saúde de sync.
- Upsert canônico: `internal_read/adapters/mirror/writer.go:74-95` (`upsertSQL`, ON
  CONFLICT) + keep-absent `:104-112` (`keepAbsentSQL`) — ranges medidos, F-r06-5.

## Clarified Decisions

- Resolved: ver Interview + design §11 (trilha de 8 decisões ratificadas).
- Accepted assumptions:
  - `channel_fees` nasce com schema das 3 camadas; esta missão popula só camadas 2 e 3
    (camada 1 = MIS-008). Reversível: coluna extra sem consumidor não força retrabalho.
  - Worker do webhook roda in-process (goroutine no server_core); `notifications_inbox`
    desacopla transport de processamento. Reversível: mover worker p/ processo separado
    depois não muda schema nem contrato.
  - Divergência auto-resolve no ingest: quando valores convergem, `resolvido_em` é gravado
    e badge some (design §9 "resolver → some"). Reversível: resolução manual pode ser
    adicionada por cima sem migração destrutiva.
  - Installation/tenant única existente (conta ML do operador). Reversível: tudo já
    tenant-scoped (exceção declarada: agregados de webhook do IC-05 são globais —
    semântica pinada no próprio IC-05, P7 r02 A-10).
  - Rate limit ML ~1500 req/min = fato #11 `assumed` (fonte parcial —
    `research/external-ml-api-facts.md`). Severado do design: ADR-02 torna o limite
    CONFIGURÁVEL (M01-C2 proíbe constante compilada); R-1 carrega mitigação/trigger/owner.
    (Eco registrado P7 r02 A-2.)
- Owner decisions still open: None — gate P1 respondido 2026-07-31.
- Blocked items: None — evidência pendente é P2 (codebase facts), não decisão de dono.

### Clarification Interview

| # | Taxon | Question | Proposed default | Operator answer |
| --- | --- | --- | --- | --- |
| 1 | lifecycle/transitions | Webhook topic `items` entra na Onda 2 ou só `orders_v2`? | Só `orders_v2` | Só `orders_v2` (2026-07-31) |
| 2 | persistence/reset | Divergência: tabela dedicada ou flag por entidade? | Tabela dedicada `divergences` | Tabela dedicada (2026-07-31) |
| 3 | build/runtime | Postura de segurança do POST /webhooks/{provider} público? | Hint não-confiável + refetch autenticado + IP log-only | Hint não-confiável + IP log (2026-07-31) |
| 4 | — | Modo de escrita do planning | Apply | Apply (2026-07-31) |

Demais taxa: sem ambiguidade bloqueante — respondidos pelo design ratificado
(actor/tenant existente; 3 camadas de fee; migrações aditivas + NULL honesto + custo
congelado; nenhuma tela nova; validação padrão MIS-006 + M0X-U*; cadências 5min/diário).

## Architecture Spine

Spine reconciliado P3 r01 (candidato Claude 10 ADRs ⊕ contraproposta Opus A-01..A-14;
detalhe integral: `planning-reviews/p3-reconciliation-r01.md` + os dois insumos). 14 entradas
ADR-lite — cada uma é escolha que 2 workers fariam incompatível:

- **ADR-01 Núcleo nativo × adapter** (design §3). Núcleo agnóstico de provider; adapter ML é
  o único que conhece endpoint/header/DTO ML. Teste da fronteira: provider novo = 1 adapter,
  zero mudança em núcleo. Gate verifica ausência de import `connectors`/tipos ML em
  núcleo/read.
- **ADR-02 Resiliência = decorator no choke point.** Backoff exp + jitter + `Retry-After` +
  token-bucket POR INSTALLATION (compartilhado entre goroutines) UMA vez, envolvendo
  `doRawWithHeaders` (`capability_adapter.go:712-731`). Limite configurável (fato ~1500
  req/min é `assumed`). 429 só vira `ErrCodeProviderRateLimited` após budget esgotado, erro
  nomeia tentativas + último Retry-After. Opt-out no-retry para writes (missão futura).
  `AccessTokenResolver` e backoff de refresh OAuth são mecanismo SEPARADO — não fundir.
  Must-fail nomeia o tempo esperado; "eventually succeeds" = passe vácuo.
- **ADR-03 Raw seletivo, PII nunca.** `Raw json.RawMessage` em todo DTO; persistência de
  `raw jsonb` SÓ em `listings` (EMENDA P7 r01 B-7: `order_shipments` NÃO tem coluna raw —
  payload de shipment carrega PII de entrega (receiver name/endereço/CEP, classe PII já
  registrada em `cmd/mlprobe/main.go:41-43`); `orders` idem, nenhuma coluna raw definida;
  ambos persistem campos TIPADOS somente). **Raw de `billing_info` NUNCA
  persiste** (documento + endereço fiscal — `buyer_fiscal_reader.go:59-94`); campos tipados
  sim. Cap 1MB mantido; truncamento grava marcador explícito. Teste: fixture com documento →
  nenhuma coluna persistida o contém. Raw nunca vaza p/ transport/tela.
- **ADR-04 Ingest resource-addressed, caminho único.** Por entidade, UM writer:
  `Ingest{Order,Listing}(ctx, tenant, installation, providerResourceID)`, idempotente por
  upsert na chave natural (idioma `writer.go:74-95` + keep-absent `:104-112` — F-r06-5:
  ON CONFLICT + keep-absent, nunca
  DELETE). Backfill, scheduler, webhook worker e refresh manual chamam O MESMO. Enumeradores
  (scan/scroll_id; `orders/search` `date_last_updated`+`sort=date_desc`; drain do inbox) só
  produzem IDs. Critérios Q3 descarregam contra ESTE seam.
- **ADR-05 Replacement-before-deletion + guard allowlist encolhente.** Teste arquitetural no
  primeiro milestone agendável: allowlist explícita dos 4 sítios read-time ML; sítio novo
  falha na hora; milestone que persiste o substituto DELETA sua entrada no mesmo commit;
  zero no fechamento da missão. Nenhuma fonte viva morre sem substituto persistido no MESMO
  merge. Assimetria evidenciada: C/D (pricing) podem morrer cedo (fallback `config` honesto
  existe); A/B (orders) só com `order_shipments` populada. Cada deleção prova before/after
  na MESMA tela browser.
- **ADR-06 MASS-CLOSURE morre por absent ≠ closed.** `ApplyCompletedPull` continua writer
  ÚNICO de `listings`, mas para de fechar em massa (`repository.go:390-394`): upsert por row
  + marcação keep-absent escopada por run; `status` = verdade do provider SOMENTE (payload
  `/items`), nunca inferido de ausência; marcação só após run declarado COMPLETO (run
  truncado por 429/deadline/cancel nunca é completo). Must-fail: abort pós-página-1 → zero
  flips p/ `closed`, fixture >1 página.
- **ADR-07 Cursor terminal, nil proibido.** `scheduler.go:42-45` APAGA cursor nil ⇒ jobs
  novos retornam cursor não-nil com phase explícita
  (`{"phase":"backfill",...}` → `{"phase":"incremental","watermark":...}`). Enum 0075 já tem
  `listings`/`orders` — registrar job NÃO precisa de migração. Nunca comparar JSONB
  round-trip byte-exact.
- **ADR-08 Dois Schedulers, cadências fixas.** Instâncias separadas 5min (orders) / diário
  (listings), padrão `synccomposition.NewProductsScheduler` (`root.go:672-677`).
  `schedule jsonb` continua não-lido (fora de escopo). Fix único na fundação:
  `incremental=false` hardcoded (`scheduler.go:160`) — pré-condição do Q4.
- **ADR-09 Ledger de fee enumerado.** JÁ MORTO (nenhum milestone reivindica): seed 16/22
  (0081; `registry_test.go:90-103`; `fee_sync.go:29` NÃO existe). MORRE nesta missão:
  degrau-3 vivo (`root.go:845-851` + `tarifflive/resolver.go:43-69`) → lookup `channel_fees`.
  EMENDA (auditoria P5 r02 N-2): `baseline_commission_percent: 0.16` (`auth_adapter.go:42-48`)
  é METADATA de catálogo do provider — contrato publicado (wiki required, OpenAPI stable key,
  SDK typed), SEM call site em pricing; fica INTOCADA (a disjunção anterior "morre ou vira
  row `origem='config'`" caiu — premissa de fallback silencioso era falsa).
  SOBREVIVE re-rotulado: `pricing_tariff_defaults` 13/16 = fallback mais fraco, proveniência
  `config`, editável. NÃO ADOTADO: `FeeSyncScheduler`/`RegisterFeeSyncerFactory` —
  `channel_fees` é alimentada pelo INGEST de listings/orders, nunca pelo fee-syncer.
  `sale_fee` é POR UNIDADE (fato live #1): camada 3 = `sale_fee × quantity`. `MissingSaleFee`
  continua disparando — nunca satisfeito por row fabricada. Todo fee em tela carrega
  (camada, origem, coletado_em); número sem proveniência = milestone reprovado.
- **ADR-10 Divergences = one-open-row-per-(entity,kind), upsert.** Chave natural
  `(tenant_id, provider, entity_type, entity_id, kind)` com no máximo 1 row aberta;
  `expected_*`/`observed_*` + timestamps de observação dos DOIS lados NOT NULL (mitiga
  falso-positivo R-5); detecção = upsert; convergência grava `resolved_at`.
  `kind ∈ {estoque, tarifa}` com CHECK, extensível aditivamente. Shape decidido ANTES de
  existirem os 2 produtores (cadências 5min × diário). Prova nas 2 direções na MESMA drive.
- **ADR-11 Webhook: ponteiro, nunca dado; always-200.** `POST /webhooks/{provider}` fino:
  corpo bounded, insert em `notifications_inbox`, 200 SEMPRE (incl. topic desconhecido, corpo
  malformado, user_id desconhecido — evita tempestade 8-retries/1h). Só `resource`+`topic`
  usados, só como ponteiro p/ refetch AUTENTICADO via ADR-04. `user_id`→installation pelo
  nosso credential store; desconhecido = gravado e descartado. IP persistido + comparado à
  allowlist oficial LOG-ONLY. Dedupe de transporte (EMENDA P5 r02 N-4, estreitado da tupla
  original): `UNIQUE (provider, notification_id) WHERE notification_id IS NOT NULL` —
  `_id` NÃO está no payload verificado (fato #6, `external-ml-api-facts.md`) e a tupla cheia
  com COALESCE bloquearia p/ SEMPRE re-notificação legítima do mesmo resource; sem `_id`
  não há dedupe de transporte e a idempotência REAL é do IngestOrder (ADR-04) +
  reconciliação 5min. Rota
  classe INTERATIVA — explicitamente fora de `registerBatchRoutes`. 1ª superfície pública
  não-autenticada de escrita do sistema. Scheduler 5min continua reconciliação obrigatória.
- **ADR-12 Migrações 0086+ pré-alocadas pelo hub.** Range disjunto explícito no brief de
  cada milestone; todas aditivas; nunca ALTER destrutivo em `products_mirror`/`listings`;
  migração aplicada NUNCA renomeia (filename = PK de `schema_migrations`); 0021 duplicado
  não se "limpa"; teste por migração no idioma de `listings_test.go` (substrings obrigatórias
  estilo `:25` + regex de corpo estilo `createTableBody` `:101` — F-r06-5).
- **ADR-13 `listing_variations` aditiva; PK de `listings` NÃO muda.** Child table
  `(tenant_id, installation_id, provider, provider_listing_id, variation_id)`, mesmo writer
  único; `listings` mantém PK com sentinela `'-'` (0036). Mudar PK radiaria destrutivo p/
  0022/0025. Acoplamento oculto nomeado: reestrutura da hidratação NÃO pode esfomear
  `AbsorbProviderSnapshots` (`connectors/source.go:54-89`) — must-fail "snapshot observer
  starved" (row count + âncoras não-regressivos vs pull pré-mudança).
- **ADR-14 `root.go` + OpenAPI/SDK = seams serializados pelo hub.** Wiring por milestone =
  1 chamada a constructor de composição local em região ancorada; hub arbitra `root.go`.
  OpenAPI + `sdk-runtime` mesmo commit (profile); SDK é ESCRITO À MÃO (client = objeto
  literal `index.ts:2113-2446`) ⇒ COMMITS de contrato FE serializados: no máximo 1 COMMIT
  de contrato (YAML+SDK) em voo por vez, hub arbitra a ordem; o CÓDIGO dos milestones
  paraleliza. EMENDA (auditoria P5 r03 P-2): a forma original "≤1 milestone com mudança de
  contrato FE em voo" contradizia a lane C ratificada (M-05∥M-06∥M-07, todos com par FE);
  o seam real é o client literal escrito à mão — seam serializado do hub, mesma forma da
  região pricing do root.go. Gate por merge: YAML sem SDK ou SDK sem YAML = blocker.

Contratos compartilhados AUTORADOS (P4, 2026-07-31) em `research/`:
IC-01 `channel-fees-` · IC-02 `divergences-` · IC-03 `orders-persistence-` ·
IC-04 `webhook-inbox-` · IC-05 `sync-health-` · IC-06 `sync-ingest-ports-` ·
IC-07 `listings-sync-interface-contract.md`. Decisões finas tomadas neles (binding):
fee mora SÓ em channel_fees — E3 perde `commission_amount/pct`/`free_shipping_cost`
(IC-07); `/sync/health` novo em vez de reusar `/sync/runs` (IC-05); inbox DDL no M-08,
listings DDL no M-04, core DDL (channel_fees/divergences/order_shipments) no M-02;
webhook documentado no OpenAPI SEM método SDK (IC-04); missed_feeds não-consumido;
tolerância 3→2 = R$0.01; formatos compostos de subject/entity id pinados (IC-01/IC-02).

## Milestone Strategy

Corte por posse de seam, ratificado pelo operador no STOP P3 r01 (2026-07-31) — escopo
idêntico ao design; muda fatiamento/ordem. Headlines (corpos no P5):

| id | headline | depende de |
| --- | --- | --- |
| M-01 `ml-client-hardening` | ADR-02 no choke point + multiget `/items?ids=` + regra raw DTO. Zero schema, zero UI. `capability_adapter.go` congela depois (endpoints novos = arquivos novos). | — |
| M-02 `sync-core-seam` | DDL 0086-0089 (`channel_fees`, `divergences`, `order_shipments`, colunas aditivas de orders) + ports de fee-layer/divergência + fix `incremental` + guard allowlist (ADR-05). Um dono, zero código de provider. | — (∥ M-01) |
| M-03 `orders-shipment-persist` | `order_shipments` populada (shipment+sla+costs+billing tipado) via ingester ADR-04 no caminho do import; sítios A/B morrem; bucket/SLA/frete/rastreio/comprador-fiscal lêem Postgres. | M-01, M-02 |
| M-04 `listings-backfill-ingest` | Backfill scan→multiget retomável, MASS-CLOSURE substituído (ADR-06), E3 + `listing_variations` + `available_quantity`, scheduler diário, refresh em lote. | M-01, M-02 |
| M-05 `listings-fees-divergence` | Camada 2 (`listing_prices` com category_id/price ingeridos; frete = honesto-desconhecido, sem shipping_options nesta missão), divergência de estoque no ingest, /anuncios colunas + badge ⚠ + filtro "divergentes". | M-04, M-02 |
| M-06 `orders-backfill-decomposition` | Backfill 12m + incremental 5min, decomposição persistida (custo congelado), camada 3 (`sale_fee×qty`), auditoria 3→2, bucket indexado. Estende ingester do M-03. | M-03, M-02 |
| M-07 `pricing-fee-read` | Resolver de pricing lê `channel_fees` com proveniência; degrau-3 vivo morre; `pricing_tariff_defaults` re-rotulado `config`; sítios C/D morrem. | M-02 (⤳ M-05: edge de qualidade de dado, não compile — pode rodar cedo com fallback `config`) |
| M-08 `webhook-ingest` | `notifications_inbox` + `POST /webhooks/{provider}` + worker in-process + registro callback (ngrok fixo), `orders_v2`. Worker só chama ingester do M-06. | M-06 |
| M-09 `sync-observability` | /integracoes: saúde por entidade (`sync_state`) + status webhook. `GET /sync/runs` já existe sem consumidor FE; construível desde o dia 1 contra `products`. | — (⤳ M-04/M-06 p/ entidades acenderem) |

Lanes: **A** M-01∥M-02∥M-09 → **B** M-03∥M-04 → **C** M-05∥M-06∥M-07 → **D** M-08.
Edge transversal ADR-14 (emendado P5 r03 P-2): commits de contrato FE serializados pelo
hub — ≤1 COMMIT de contrato em voo por vez; código paraleliza.

## Parallel Execution Plan

### Dependency DAG (edge = artefato que força)

```
M-01 ──┬─→ M-03 ───→ M-06 ──→ M-08
M-02 ──┤                       ↑
       ├─→ M-04 ───→ M-05     M-09 (porta WebhookStatsReader IC-05)
       └─→ M-07  (⤳ M-05: qualidade de dado, não compile)
M-06 (⤳ M-05: auditoria 3→2 precisa de camada 2 populada — fica MUDA sem ela, não quebra)
M-09 (⤳ M-04/M-06: entidades acendem sozinhas; ⤳ M-02: last_incremental_at real — pré-fix
NULL uniforme é honesto e o gate passa; edge dura → M-08 acima)
```

| edge | artefato que força |
| --- | --- |
| M-03→M-01 | ingest de shipment = 3-4 GETs paralelos/pedido contra adapter onde 429 é seco |
| M-04→M-01 | multiget (F-02) é o consumidor; N+1 atual + acúmulo em memória |
| M-03→M-02 | 0088/0089 + ports IC-01/02 + guard (deleta entradas A/B) |
| M-04→M-02 | contrato de cursor + fix incremental |
| M-05→M-04 | camada 2 requer category_id/price ingeridos; divergência requer available_quantity |
| M-05→M-02 | ports ChannelFeeWriter (camada 2) + DivergenceRecorder (estoque) — IC-01/IC-02 |
| M-06→M-03 | estende IngestOrder (writer único ADR-04 — nunca 2º writer) |
| M-06→M-02 | camada 3 + divergência tarifa (ports) |
| M-07→M-02 | ChannelFeeReader (sem ele o resolver não tem o que ler pós-degrau-3) |
| M-07⤳M-05 | fallback config é ratificado-honesto — M-07 pode rodar cedo |
| M-06⤳M-05 | auditoria 3→2 precisa de camada 2 populada — sem ela auditoria fica MUDA, não quebra (já declarado em M-06/milestone; mesma classe do M-07⤳M-05 — F-r06-3) |
| M-08→M-06 | worker chama o IngestOrder ESTENDIDO do M-06 (decomposição+camada 3) — não o v1 do M-03 |
| M-08→M-09 | porta WebhookStatsReader (IC-05) + injeção `WithWebhookStatsReader` da região ancorada do M-08 |
| M-09⤳M-04/M-06 | endpoint entidade-agnóstico; rows aparecem quando jobs registram |
| M-09⤳M-02 | `last_incremental_at` REAL exige fix `incremental` (M-02 F-03, IC-05) — pré-fix, NULL uniforme é honesto e o gate do M-09 PASSA; lane A paraleliza sem ordem de close (F-r07-1) |

### Ownership matrix (6 eixos; célula = DONO exclusivo)

| Milestone | Go packages/files | Migrações | Tabelas DB (write) | OpenAPI/SDK | FE rotas/componentes | root.go wiring |
| --- | --- | --- | --- | --- | --- | --- |
| M-01 | `internal/modules/connectors/adapters/mercado_livre/` — exclusividade sobre os ARQUIVOS EXISTENTES (único a editar capability_adapter.go); arquivos novos no dir permitidos downstream (M-03 F-01 readers) | — | — | — | — | — |
| M-02 | packages novos fees/divergences; scheduler.go (fix pontual); guard test | 0086-0089 | channel_fees, divergences, order_shipments (DDL), orders (DDL aditiva) | — | — | — |
| M-03 | `orders/**` (application+transport), readers ML novos (arquivos novos) | — | order_shipments, orders (rows), channel_fees? NÃO (camada 3 é M-06) | /orders DTOs (par YAML+SDK) | /pedidos (colunas servidas — sem mudança FE além do DTO) | edita região orders existente `:576-601` in-place (troca readers A/B `:591-592` por readers de banco, deleções inclusas — mesma classe de exceção do M-07; F-r07-2), hub arbitra |
| M-04 | `listings/**` | 0090-0092 | listings, listing_variations | — | — | 1 linha ancorada |
| M-05 | ingest ext camada2/divergência + read aditivo (lock aditivo do M-04 PÓS-close: `listings/application/**`, `listings/transport/**`, `listings/composition/**`, `listings/ports/**` e `listings/adapters/postgres/repository.go` — F-r05-3, ampliado P7 r01 B-4), FE /anuncios | — | channel_fees (camada 2), divergences (estoque) | /listings DTOs+param (par) | AnunciosTable/PedidosPage? NÃO — só /anuncios | — |
| M-06 | `orders/**` (pós-close M-03 — herda posse da lane) | 0094-0095 (reserva, pode não usar) | orders (decomposição), channel_fees (camada 3), divergences (tarifa) | /orders DTOs margem (par) | /pedidos margem/Fila | 1 linha ancorada (scheduler 5min) |
| M-07 | `pricing/**` resolvers (tarifflive morre) | — | — (read-only de channel_fees) | /pricing DTOs proveniência (par) | /precos proveniência+⚠ | edita região pricing existente `:828-858` + remoção de imports tarifflive/tariffcomposite (`root.go:99,101` — F-r04-5); hub arbitra (região-edit, uma das DUAS exceções — a outra é M-03 região orders; F-r07-2) |
| M-08 | package webhook/inbox novo | 0093 | notifications_inbox | path /webhooks (SEM método SDK — IC-04) | — | 1 linha ancorada (+ chamada `WithWebhookStatsReader` no seam do M-09, na PRÓPRIA região) |
| M-09 | endpoint /sync/health + FE seção | — | — (read-only) | /sync/health (par) | IntegracoesPage seção nova | 1 linha ancorada |

Regras transversais da matriz:
- **root.go**: hub = resolver of record; cada milestone entra com 1 chamada de constructor
  em região ancorada própria (ADR-14). O BLOCO DE IMPORTS do root.go não pertence a nenhum
  milestone — hub-resolved: cada milestone adiciona/remove SÓ os imports que a própria
  região ancorada exige; conflito mecânico de import = merge do hub (P5 r04 F-r04-5).
- **OpenAPI/SDK** (ADR-14 emendado P5 r03 P-2): ≤1 COMMIT de contrato FE em voo — na
  lane C, M-05/M-06/M-07 têm TODOS par FE ⇒ hub serializa os commits/merges de contrato
  dentro da lane (código pode paralelizar; o COMMIT de contrato não).
- **Posse de `orders/**`**: M-03 na lane B; M-06 herda a posse na lane C (M-03 já fechado).
  Nunca simultâneos.
- **Posse de ingest de listings**: M-05 estende superfícies do M-04 PÓS-close (lock
  aditivo registrado: `listings/application/**`, `listings/transport/**`,
  `listings/composition/**`, `listings/ports/**` e
  `listings/adapters/postgres/repository.go`, additive-only — o write-set REAL do M-05
  inclui transport, o repository de leitura, a fiação em composition e a porta de leitura
  mirror em ports; grant alinhado ao write-DAG na auditoria P5 r05 F-r05-3 e ampliado a
  composition/ports no P7 r01 B-4).
- **Posse do package `sync/application/`** (lane A, M-02 ∥ M-09 — lock aditivo registrado,
  P5 r02 N-8): M-02 F-03 é dono de `scheduler.go` + helper de contrato de cursor (arquivo
  novo); M-09 é ADDITIVE-ONLY — arquivos novos `sync/application/health_*` + `sync/
  transport/**`; `scheduler.go` intocado pelo M-09.
- Migração: ranges pré-alocados acima; range reserva do M-06 pode ficar vazio (gap ok).

### Write-DAG de features (colisões cruzadas já serializadas nos DAGs internos)

Colisões cross-milestone com serial edge: M-05.F-01/F-02 → após M-04 close (lock aditivo);
M-06.* → após M-03 close (herança de lane); M-07 root.go região pricing → merge arbitrado
pelo hub; arquivo da allowlist do guard (dono M-02 F-04) → escrito por M-03 F-03 (remoção
A/B) e M-07 F-01 (remoção C/D), serializados pelo edge lane B → lane C; bloco de imports
do root.go → hub-resolved, sem dono (regra da matriz — F-r04-5). Nenhum write-set overlap
restante sem edge nomeado ou resolução do hub registrada (P5 r04 F-r04-6).

## Quality Attributes

| Attribute | Target (concrete) | Owner (ADR/seam) | Validation criterion |
| --- | --- | --- | --- |
| Q1 Performance | /pedidos e /anuncios <2s no read; zero call ML no read | Onda 0 + arquitetura mirror-first | MIS07-C2: as 3 amostras medidas <2s com dados reais; nenhum sítio de transport importa client ML |
| Q2 Security | webhook = hint não-confiável; dado real sempre refetch na API autenticada; dedupe no inbox; IP origem gravado + comparado à allowlist oficial (log-only) | contrato do webhook (P4) | notificação forjada não injeta dado (só fetch idempotente); IP não-oficial aparece em log |
| Q3 Reliability | backfill retomável por cursor (scan scroll_id / date_last_updated); ingest idempotente (upsert por resource id); buraco >2 dias coberto por reconciliação | sync_state + desenho de ingest único (webhook e scheduler = 2 portas do mesmo caminho) | matar backfill no meio → retomar sem duplicar; reprocessar notificação → zero duplicata de DOMÍNIO (IngestOrder idempotente; dedupe de inbox só com `_id` — ADR-11 emendado; ver M-08 Done Means) |
| Q4 Observability | /integracoes mostra saúde por entidade (sync_state) + última notificação recebida | tela /integracoes + sync_state | critério dirigido em browser real |
| Q6 Maintainability | teste da fronteira (design §3): provider novo = 1 adapter, zero mudança em núcleo | arquitetura núcleo×adapter | revisão de gate: núcleo não importa tipo ML |

## Non-Functional Scope

| Declined attribute | Reason |
| --- | --- |
| Q5 Usability além do padrão existente | nenhuma tela nova; DESIGN-REFERENCE já é gate visual |
| Q7 Compatibility/portability | mesmo stack/runtime existente; nada cruza |
| Verificação criptográfica de webhook (assinatura/token) | decisão P1 do operador: webhook = hint não-confiável + refetch autenticado; `ip_official` é INFORMATIVO/log-only, NUNCA gate de aceitação — derivação do IP pinada no IC-04 (P7 r01 B-9) |

## Validation Strategy

Padrão MIS-006 (design §9): contrato por milestone com critérios dirigidos em browser real
(M0X-U*); fixtures multi-página; live-drive final (conta ML real, backfill medido, pedido
real decomposto, webhook disparado por evento real); must-fail nomeando a falha; divergência
provada nas 2 direções. Evidência: feature → `<feature-root>/validation.md`; milestone →
`<milestone-root>/validation-result.md`; missão → `validation-result.md`.

## Risks

| id | risk | likelihood | impact | mitigation | trigger | owner |
| --- | --- | --- | --- | --- | --- | --- |
| R-1 | Rate limit ML (~1500 req/min, fonte parcial) degrada backfill 12 meses | M | M | backoff+jitter+Retry-After pré-requisito; budget de goroutines; backfill retomável | 429 em série no live-drive | milestone de pedidos |
| R-2 | ngrok domínio fixo cai/muda → webhook surdo | M | L | scheduler 5min é reconciliação garantida; /integracoes expõe última notificação | última notificação > cadência esperada | operador |
| R-3 | Conta ML real com poucos anúncios/pedidos esconde defeito de paginação | H | M | fixtures multi-página obrigatórias (lição CHIP-MERCADO) | fixture <2 páginas em contrato | planning (P6) |
| R-4 | sale_fee POR UNIDADE regride (fato só por medição live, doc não declara) | L | H | decomposição testada contra pedido real no live-drive; must-fail nomeia | margem divergente no pedido real | milestone de pedidos |
| R-5 | Divergência falso-positiva (corte vendável ERP × available_quantity ML defasados) | M | M | divergência calculada no INGEST com timestamps de ambos os lados persistidos (NOT NULL — ADR-10) | badge em anúncio recém-sincronizado | milestone de anúncios |
| R-6 | PII escapa via raw-persistência + 1ª superfície pública de escrita (webhook) | M | H | ADR-03 (raw de billing_info NUNCA; marcador de truncamento) + ADR-11 (corpo de webhook nunca alcança tabela de domínio); teste com fixture contendo documento | migração adicionando `raw jsonb` a tabela com dado fiscal/comprador | M-02 (schema) + M-03/M-08 |
| R-7 | Colisão de merge em seam sem dono: numeração de migração, `root.go`, SDK escrito à mão + regra same-commit | H | M | ADR-12 (ranges no brief), ADR-14 emendado (constructor local + commits de contrato serializados: ≤1 COMMIT FE em voo) | 2 briefs sem ranges disjuntos; 2 COMMITS de contrato simultâneos ou não arbitrados pelo hub | hub (P5 ownership matrix) |

## Handoff

- Current status: **planned** — P7 fechado 2026-08-01: Claude-side Ready (r03) + Sol
  DISPENSADO por decisão do operador (verbatim "Eu aprovo não vai ter Sol" — supersede o
  waiver de retroativo; `sol-unavailable-p7-r03.md` §SUPERSEDIDO). Fold conjunto:
  `readiness-review.md` ⇒ joint Ready.
  P7 r03 (2026-08-01, rodada ENXUTA por ordem do operador — "sonnet subagents conferir e
  somente um revisar", custo): conferência mecânica 17/17 reparos PRESENT
  (`p7-verify-repairs-r03.md`) + 1 assento frio ★2 → PASS (`p7-seat-star2-r03.md`);
  fold `p7-claude-readiness-r03.md` = Ready com 2 desvios registrados (rodada enxuta;
  SEM manifesto congelado — shell wedge B-8 total na sessão, sha256 impossível; input
  prescrito + imutabilidade de janela única substituem). Antes do Sol P7 HIGH: congelar
  `p7-input-r04.sha256` com shell funcional. Dívida A-9 (retention inbox + amplificação
  refetch) declarada sem valor — decisão do operador, não bloqueia.
  Scratch a remover quando shell voltar: `freeze_tmp.py`, `freeze_srv.py` (raiz da missão).
  Histórico r02: crew fria ×5 →
  fold "Needs revision" (`planning-reviews/p7-claude-readiness-r02.md`; ★2 FAIL via seat 5
  double-pass: ★2-A PK `listing_variations` 5-col vs 4-col, ★2-B slug `mercadolivre` vs
  `mercado_livre`); Sol NÃO despachado (procedimento). Reparos ★2-A (tuple ADR-13 5-col
  propagado p/ IC-07 + M-04/F-01, razão registrada no IC-07) e ★2-B (`mercado_livre` nos
  2 loci) + 12 advisories auto-fixáveis aplicados (A-1/2/4/5/6/7/8/10/11/12/13/14).
  Dívida declarada SEM valor decidido (P7 r02 A-9, seats 4+5): `notifications_inbox` sem
  retention/prune + amplificação de refetch (forge flood drena token-bucket M-01) —
  decisão de valor é do operador; rows são inertes (M08-C3), não bloqueia planning.
  r01: fold "Needs revision" B-1..B-9, reparado (`p7-claude-readiness-r01.md`).
  P3 fechado; P4 fechado (7 interface contracts em `research/`);
  P5 FECHADO em PASS (audit loop r01–r08, verdict r08 PASS zero blocking em
  `planning-reviews/p5-claude-decomposition-audit-r08.md`; 5 advisories foldados,
  `p5-reconciliation-r08.md`; baseline pós-fold = `p5-input-r09.sha256`);
  P6 FECHADO (validation contracts autorados 2026-08-01: `validation-contract.md` da
  missão MIS07-C1..C9 + 9 `M-*/validation-contract.md` M0X-C*/M0X-U*, base_sha dd89d4b3);
  planning_phase = readiness (P7).
- Current owner: Mission Strategist (sessão de planning).
- Next owner: hub (sessão harness-hub) — execução por ondas conforme
  `## Parallel Execution Plan` (lane A: M-01 ∥ M-02 ∥ M-09).
- Next action: hub commita os artefatos de planejamento (shell desta sessão morto — B-8;
  nunca push sem permissão do operador), remove scratch `freeze_tmp.py`/`freeze_srv.py`,
  e despacha a onda 1. Decisão A-9 (retention inbox) entra quando o operador quiser.
- Gate Sol: DISPENSADO — decisão do operador 2026-08-01 ("Eu aprovo não vai ter Sol"),
  registrada em `sol-unavailable-p7-r03.md` §SUPERSEDIDO + `readiness-review.md`;
  supersede o waiver de 2026-07-31 (retroativos cancelados).
- Required artifact paths: `research/*-interface-contract.md`, `planning-reviews/p3-*`,
  `planning-reviews/p5-*`, milestones `M-*/`.
- Required evidence paths: `validation-result.md` (missão), `M-*/validation-result.md`.
- Blocked decisions: None — gates P1 e P3 fechados; P5 audit é subgate interno.
