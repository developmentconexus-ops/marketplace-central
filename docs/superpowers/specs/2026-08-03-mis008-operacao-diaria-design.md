# MIS-008 — Operação diária: design

**Data:** 2026-08-03
**Estado:** design aprovado, plano não escrito
**Antecessores:** MIS-006 (encerrada), MIS-007 (ondas A e B fechadas), P2.b (em execução em outra sessão)

---

## 1. Problema

Três registros paralelos descrevem o que falta fazer, e só um deles funciona.

- **Fatias `P`** — `P1`..`P6` aparecem em dois planos vivos escritos com um dia de diferença e **significam coisas diferentes** em cada um. No plano do P1, `P3` é motor de DIFAL, `P4` é `quality_score`/`infractions`, `P5` é sub-recursos de envio, `P6` é KPIs de backend. No plano do P2.b, `P2.c` é reconciliação, `P3` é fiscal em `/anuncios`+`/mercado`, `P4` mata o 4%.
- **Dívidas `D`** — `D-16`/`D-17`/`D-18`/`D-38`/`D-39` cada uma tem **dois significados** em documentos diferentes (`.mnfs/HARNESS-DEBTS.md` × specs fiscais). Colisões confirmadas por medição.
- **Milestones `M`** — funcionam. São os únicos que têm endereço (`.mnfs/MIS-*/M-*/`).

Diagnóstico: o `M` funciona porque tem endereço; `P` e `D` não têm. Um id sem endereço não é registro, é apelido.

Consequência medida: a migration `0093` está reivindicada por **dois branches não mergeados** (`worktree-p2-dinheiro-real-pedidos` tip `e338b279`, `NOT_MERGED`; e `p2b-imposto-ex-ante` com `0093_icms_matrix.sql`). `main` está em `0092`.

## 2. Linha de chegada

**Uma pessoa opera uma conta Mercado Livre, todo dia, sozinha.** Um tenant, uma conta, um provider.

Esse critério é a régua de escopo. Tudo que não estiver nesse caminho é podado ou vira dívida registrada — não fica meio construído.

Não está no escopo: multi-tenant, segundo marketplace, escrita automática no provider (a lane `provider-write` exige 7 portões e nenhuma fatia planejada a usa), ML como fonte de verdade de catálogo.

## 3. Registro único

Cria-se `.mnfs/MIS-008-operacao-diaria/` com dois arquivos:

- **`BACKLOG.md`** — namespace novo `F-01`..`F-nn`. Cada `F` tem: id, enunciado, medida (file:line que prova que o defeito existe), critério de aceite falsificável, lane de prova, e onda.
- **`DEBTS.md`** — dívidas com id novo, mais **tabelas de redirect** mapeando os ids antigos (`P*`, `D-*` de cada documento de origem) para o id novo ou para `MORTA`.

Os planos antigos ganham banner de superseding no topo. Nenhum id antigo é reusado.

**Regra:** id sem endereço em `.mnfs` não existe. Se algo não tem linha no `BACKLOG.md` ou no `DEBTS.md`, não está planejado.

### 3.1 Redirects obrigatórios (medidos)

| id antigo | origem | destino |
|---|---|---|
| `D-16` | `.mnfs/HARNESS-DEBTS.md:299-311` | VIVA → bloqueio de `F-00` |
| `D-17` (harness, emit tsc) | `.mnfs/HARNESS-DEBTS.md:313-335` | RESOLVIDA parcial |
| `D-17` (fiscal, `CODEMP=1`) | `docs/superpowers/specs/2026-08-02-p2b-imposto-ex-ante-design.md:348,124` | VIVA → bloqueio de `F-06`, `F-09` |
| `D-18` (harness, worktree órfão) | `.mnfs/HARNESS-DEBTS.md:337-354` | infra, não migra |
| `D-18` (fiscal, ST recuperável) | `...p2b-imposto-ex-ante-design.md:124,356` | dívida de domínio, migra |
| `D-21` | `...p2b-imposto-ex-ante-design.md:301,313,359` | VIVA → bloqueio de `F-06` |
| `D-31`, `D-32`, `D-35` | `docs/superpowers/plans/2026-08-02-p2b-modulo-fiscal.md` | MORTAS (auto-declaradas) |
| `D-38` (imposto-ex-ante) | `docs/superpowers/plans/2026-08-02-p2b-imposto-ex-ante-plan.md:41,599` | VIVA → bloqueio de `F-06` |
| `D-38` (modulo-fiscal) | `...p2b-modulo-fiscal.md:2900` | dado ERP, migra separado |
| `D-39` (dois sentidos) | `...plan.md:600` e `...modulo-fiscal.md:2901` | ambos migram com ids distintos |

## 4. Classes de defeito

Seis classes. **Contagem bruta não é defeito** — cada classe carrega o predicado que separa defeito real de forma legítima. Esse é o aprendizado central da medição: sem predicado, duas classes estavam infladas por ordem de grandeza.

### Classe 1 — escrita que apaga coluna de outro escritor

Predicado: `SET x = EXCLUDED.x` cru só é perigoso com **2+ escritores na mesma tabela** (ADR-C2). Escritor único não corre risco.

- Bruto: 32 sítios de `DO UPDATE SET`
- Real: **1 tabela** — `products_mirror`, escrita por `erp_import/adapters/postgres/mirror_repository.go:117,141` **e** `internal_read/adapters/mirror/writer.go:75,105`
- As outras 8 tabelas auditadas (`listings`, `listing_variations`, `orders_marketplace_orders`, `orders_marketplace_order_items`, `order_shipments`, `channel_fees`, `product_enrichments`, `divergences`) têm escritor único

### Classe 2 — `string` onde `NULL` é legítimo

Predicado: só é defeito quando a fonte pode honestamente omitir o campo.

- `catalog/adapters/internalread/reader.go:284-289` — helper `ptr()` colapsa `nil → ""`. `EAN` e `Reference` chegam vazios em vez de desconhecidos
- `reader.go:277-283` — `Description` e `BrandName` **nem são copiados** do `CanonicalProduct` (`canonical_product.go:69-82`). Perdidos antes de chegar
- `orders/application/ingest_service.go:190,191` — `derefString` sobre `ShippingID` e `BuyerNickname`, ambos com comentário explícito de omissão na origem (`connectors/domain/order_detail.go:27-28,31-32`)
- Indeterminados (não chutar): `SellerSKU`, `Title`, `ProviderStatus` — `string` puro no pipeline inteiro, sem sinal estrutural nem comentário. Resposta vem de leitura de raw, não de adivinhação

### Classe 3 — contrato `required` + `nullable`

Predicado: `nullable` é legítimo quando o valor pode honestamente faltar; é mentira quando sempre existe e o `nullable` só evita recusar payload quebrado.

- Bruto: 120 campos em 42 schemas
- Amostra de 12 medida: **8 legítimos, 3 mentira, 1 indeterminado**
- As 3 mentiras são piores que a classe original:
  - `OrderRead.currency` — **sem produtor**. Não aparece no `SELECT` nem no `INSERT` (`orders/adapters/postgres/order_repo.go:220-278`, `:634-653`). Sempre nulo
  - `OrderRead.fulfillment` — idem
  - `OrderRead.nf_state` — **carrega outra coisa**. A 8ª coluna do `SELECT` (`order_repo.go:222-223`) é o `CASE` de vínculo Sankhya (`'linked'`/NULL, `evidence_state='exact'`). A tela mostra status de vínculo com nome de estado de nota fiscal

**Consequência de método:** a auditoria dos 120 campos **não é fatia**. Vira verificação por schema dentro de cada `F` que toca aquele schema. Varredura cega produziria churn em 8 de cada 12 campos.

### Classe 4 — raw persistido e reconciliação

- `listings.raw` é escrito (`listings/adapters/postgres/repository.go:469`)
- Pedidos: a coluna **existe** — `migrations/0027_orders_marketplace_orders.sql:15`, `raw_provider_ref jsonb NULL` — mas todo escritor grava uma string de path (`"/orders/" + orderID`, `ingest_service.go:194` via `safeOrderProviderReference`). Não é "falta coluna": é coluna certa com conteúdo errado
- `platform/rawkeys/rawkeys.go` `Undeclared` tem **zero call-sites de produção**. Único uso é guard de compile/test em `connectors/adapters/mercado_livre/items_multiget_rawkeys_test.go:43`

### Classe 5 — KPI calculado sobre uma página, não sobre o conjunto

Reformulada pela medição. Não é "cálculo no cliente" — é **agregação sobre resultado truncado apresentada como total**.

- `pedidos/PedidosPage.tsx:62` — ao bater `MAX_ORDER_PAGES`, só `console.warn`; devolve lista parcial como se fosse completa
- `vinculos/useVinculosResolved.ts:59` — "Resolvidos hoje" capado em 20 pelo default do SDK (`packages/sdk-runtime/src/index.ts:2226`), sem aviso de truncamento
- `vinculos/VinculosPage.tsx:37-39,70` — "Alta confiança" é filtro client-side sobre página capada em 200
- Precedente: CHIP-MERCADO já teve o mesmo defeito ("truncamento silencioso de página-1 invisível ao live-drive")

**Terceira ocorrência da mesma classe.** Stop-the-line se aplica: conserto geral, não pontual.

### Classe 6 — duplicação e posse

Classe descoberta por medição; não estava no desenho original e é maior que três das anteriores.

**Margem tem 4 fórmulas divergentes:**

| onde | fórmula | vivo? |
|---|---|---|
| `pricing/domain/decompose.go:180-188` | preço − (comissão+imposto+taxa_fixa+frete+difal+tarifa_full+custo) | sim — `/precos` |
| `orders/domain/order_decomposition.go:93` `BuildProfitability` | Total − Comissão − Frete − Imposto − Difal − Custo; taxa_fixa/tarifa_full sempre `nil` por design ADR-17 | sim — `/pedidos` |
| `pricing/application/service.go:45-47` | base − custo − (base×comissão%) − taxaFixa − frete | **morto** |
| `pricing/application/batch_orchestrator.go:238-240` | idêntica à de cima, reimplementada do zero | **morto** |
| `listings/application/read_service.go:774` `belowMargin` | `preço×(1−tetoICMS) < custo×(1+margemMín)` — booleano de pior caso | sim — dashboard |

**Regra de faixa de margem tem 3 donos, 2 são cópias degradadas:**

- `packages/ui/src/MarginChip.tsx` — limiar configurável, lido do profile por `precos/PricingMatrix.tsx:143-146`
- `pedidos/pedidosFormatters.ts:19-23` `marginBandClass` — reimplementa 18%/10% hardcoded, ignora o profile
- `mercado/mercadoFormatters.ts:8` — `MARGIN_MIN_PCT = 18` hardcoded, renderizado como fato no cabeçalho

**Contrato inflado — 11 operações no OpenAPI e no SDK, zero chamadas do web:**

`POST /pricing/simulations`, `POST /pricing/simulations/batch`, `POST|GET /profitability/manual-adjustments`, `POST /profitability/profit-snapshots/calculate`, `GET /profitability/profit-snapshots`, `GET /marketplaces/fee-schedules`, `GET /market/observations`, `GET /market/references`, `GET /mutations` (métodos SDK em `packages/sdk-runtime/src/index.ts:2329,2342,2346,2348,2352,2361,2363,2392,2437,2443`).

**Costura no cliente** — a API de leitura não entrega a linha pronta, então a tela vira orquestrador:

- `precos/PricingMatrix.tsx:101-141` — 1 + N + N chamadas (`listMarketAggregates`, depois `pricingDecompose` e `listListingsByProduct` por produto)
- `mercado/MercadoPage.tsx:156-173` — `buildOppRows(facts, aggregates, verdicts)`, junção posicional de 3 fontes

**Vazamento de vocabulário ML fora do adapter (ADR-C4)** — 4 sítios, não uma enchente:

- `connectors/domain/capability.go:297` — `ListingTypeID` dentro de `domain`, sem coluna persistida que justifique. O pior dos quatro
- `integrations/application/catalog_match_probe.go:16-17,56` — literais `"gold_special"`/`"gold_pro"`
- `integrations/application/provider_operation_service.go:168`
- `integrations/transport/auth_handler.go:326`

### O que a medição **não** confirmou

Registrado para evitar refactor desnecessário:

- **`transport/` está limpo.** Zero handlers calculando. Só parse de query param e delegação
- **Status de pedido tem fonte única** — `orders/domain/order_bucket.go:48` deriva, `pedidosTabs.ts:30-41` só filtra o campo pronto
- **DIFAL tem fonte única** — `pricing/domain/difal.go:59-76`, reusado por `orders/adapters/pricingtax`
- **Comissão não é duplicada** — separação correta: `sale_fee` verbatim no pedido realizado, `comissao_pct`/policy na simulação
- **ADR-17 na tela passa.** Zero colapsos de desconhecido→zero nas telas roteadas. `?? 0` sempre opera sobre contador ou `length` já resolvido. `degraded` é lido e usado corretamente em `dashboard/DashboardPage.tsx:34,68`

## 5. Continuidade operacional

Eixo que nenhuma classe de defeito enxergava, e que sozinho invalida a linha de chegada.

### 5.1 Falha de token ML é invisível em todas as camadas

- `integrations/background/refresh_ticker.go:64` — `case <-ticker.C: _ = t.RunOnce(ctx)`. Erro descartado, **sem log**
- `refresh_ticker.go:37-51` — `RunOnce` retorna no primeiro erro e não tenta as sessões restantes do lote
- `integrations/application/auth_flow_service.go:440-441` — caminho de erro não grava `RefreshFailureCode` nem `ConsecutiveFailures`. Só o sucesso escreve (`:456-468`), zerando os campos de falha
- `HealthStatus` só é setado para `Healthy` — `auth_flow_service.go:375,398,470`, todos de sucesso. Nenhum call-site produz `Warning`/`Critical`
- `integrations/domain/refresh_policy.go:16-63` `ClassifyRefreshError` existe com **zero chamadores de produção**
- `/integracoes` lê `HealthStatus` via `auth_handler.go:47` → `GetAuthStatus:528-546` — campo congelado no último sucesso

Resultado: token expira, tudo para, a tela continua verde.

### 5.2 Cobertura de sync

Jobs periódicos existentes, todos com cadência hardcoded:

| job | iniciado | cadência | atualiza |
|---|---|---|---|
| RefreshTicker (token) | `root.go:683` | 5min | `integration_auth_sessions` |
| StateCleanup | `root.go:684` | 15min | `integration_oauth_states` |
| FeeSyncScheduler | `root.go:685` | 15min | `fee_schedules`, `marketplace_fee_definitions` |
| ProductsScheduler | `root.go:694-699` (só se `activeSourceLookup != nil`) | 15min | `products_mirror`, `sync_state` |
| ListingsSchedulers | `root.go:775-777` | **24h** | `listings`, `listing_variations` |
| MutationPoller | `cmd/server/main.go:38` | 2s | `mutations` |

Nenhum NO-OP no estado atual (o comentário em `root.go:686-689` documenta o caso histórico do products scheduler, já corrigido).

Tabelas sem job periódico:

| tabela | única escrita |
|---|---|
| `orders_marketplace_orders`, `order_shipments` | `POST /orders/import` manual (`root.go:597-600`) |
| `market_*` (7 tabelas) | clique do operador em `market/transport/collection_handler.go` |

### 5.3 O que já funciona

- Falha de sync de products/listings **é** visível: `sync/application/scheduler.go:153-158` `RecordFailure` → `GET /sync/health` (`health_handler.go:43`, registrado `root.go:913`) → `integracoes/SyncHealthCard.tsx`
- Rate limit tratado: `connectors/adapters/mercado_livre/resilience_decorator.go` — token bucket 900/min (`:48,113-132`), retry com `Retry-After` e backoff exponencial+jitter em 429 (`:211-244,258-289`), `maxAttempts=5`/`maxTotalWait=30s`. Escritas usam `doRawWithHeadersNoRetry` sem retry, decisão intencional (ADR-02, comentário em `:35`)

## 6. Rede de testes

Levantada antes de planejar deleção, porque `F-06` e `F-08` são remoção.

**Predominantemente FORTE.** `TestDecomposeGolden` (`pricing/domain/decompose_golden_test.go:19`) compara struct completa com `reflect.DeepEqual` (`:36-212`) e fecha a soma (`:242`). `BuildProfitability` (`order_decomposition_test.go:5,57,69,96,134,169`), `TestBelowMarginTriState` (`read_service_test.go:776`), `TestScanReadModel_ShippingID` (`order_repo_scan_test.go:37`) e `MarginChip.test.tsx` assertam valor concreto.

Fracos e ausentes:

- `pricing_handler_test.go:85` — FRACO, asserta só presença (`"expected non-empty items array"`)
- `router_registration_test.go:172-174` — FRACO, testa rota-não-404, não o corpo
- `pedidos/pedidosFormatters.ts` — **sem teste próprio**. Cobertura acidental via `PedidosPage.test.tsx` (render, não unidade). É exatamente o arquivo do `marginBandClass` que `F-07` vai matar

**Custo real de apagar `service.go` + `batch_orchestrator.go`:**

- Testes que morrem junto (testam só o morto): `tests/unit/pricing_service_test.go`, `tests/unit/pricing_handler_test.go`, `pricing/application/batch_orchestrator_test.go`
- Teste que precisa **edição** (testa algo que sobrevive): `tests/unit/router_registration_test.go:154-166`
- **Produção que quebra:** `internal/composition/root.go:852-869,893` e `pricing/transport/http_handler.go:17-25,53,80,148`. O `Handler` compartilha struct com `WithCalc`, então tocar nele atinge o caminho vivo de `/precos`

**Contradição não resolvida, registrada como primeira medição de `F-06`:** uma varredura mediu `batch_orchestrator.go:238-240` como fórmula mais simples que `decompose.go` (sem difal, sem limiar de frete); outra encontrou `batch_orchestrator_test.go:160`, teste de paridade que compara byte-a-byte com `decompose`. Não podem ser ambas verdade. `F-06` começa rodando o teste de paridade e lendo o veredito — não assume nenhuma das duas.

## 7. Ondas

Ordem por classe de defeito, atravessando telas. Remoção antes de conserto: consertar `nil→""` em código que vai morrer é trabalho jogado fora.

### Task 0 — pré-requisito de registro

1. Merge de `worktree-p2-dinheiro-real-pedidos` na `main`
2. P2.b renumera `0093_icms_matrix.sql` → `0097`
3. Cria `.mnfs/MIS-008-operacao-diaria/` com `BACKLOG.md`, `DEBTS.md` e as tabelas de redirect da §3.1
4. **Re-mede** `OrderRead.currency`/`fulfillment`/`nf_state` contra a `main` pós-merge

O passo 4 existe porque o branch não mergeado pode já ter consertado dois desses campos. Escrever `F-03` contra o estado de hoje é escrever contra estado velho.

### Onda 0 — sobrevivência

Sem isso a operação diária não existe.

| id | o que | medida |
|---|---|---|
| `F-A1` | Falha de refresh vira fato: grava `RefreshFailureCode`/`ConsecutiveFailures`, `ClassifyRefreshError` ganha chamador, `HealthStatus` sai de `Healthy`, `/integracoes` mostra | §5.1 |
| `F-A2` | `RunOnce` para de abortar o lote no primeiro erro | `refresh_ticker.go:37-51,64` |
| `F-00` | Scheduler periódico de pedidos — **bloqueado por `D-16`** | `root.go:597-600` |
| `F-A3` | **FECHADA 2026-08-03 (`afb6b54a`)**. Coleta de mercado ganha job periódico **e** idade visível em toda linha de `/mercado`. Decisão do operador em 2026-08-03: as duas coisas, não uma ou outra — job resolve o dado velho, idade visível resolve a confiança quando o job falha | §5.2 |

### Onda 1 — remoção

| id | o que | medida |
|---|---|---|
| `F-06` | Consolidar margem. `service.go` e `batch_orchestrator.go` morrem; `belowMargin` fica e é renomeado pro que é (pior caso, não margem); `decompose.go` e `BuildProfitability` ficam separados com a diferença **declarada** — simulação e realizado não são a mesma conta. **Bloqueado por `D-21`, `D-38`, `D-17`(fiscal)**. Começa pela medição de paridade da §6 | §4 classe 6, §6 |
| `F-07` | Faixa de margem com dono único: `MarginChip` com limiar do profile. `marginBandClass` e `MARGIN_MIN_PCT` morrem | `pedidosFormatters.ts:19-23`, `mercadoFormatters.ts:8` |
| `F-08` | Podar as 11 operações sem consumidor do OpenAPI **e** do SDK no mesmo commit. Critério: o que não está no caminho da linha de chegada sai; o que está vira `F` com tela | §4 classe 6 |

### Onda 2 — forma

| id | o que | medida |
|---|---|---|
| `F-01` | Colisão de dois escritores em `products_mirror` + guarda contra um terceiro surgir calado | §4 classe 1 |
| `F-02` | `ptr()` para de colapsar `nil→""`; mapper copia `Description`/`BrandName` | `reader.go:277-289` |
| `F-03` | Campos sem produtor e campo com rótulo errado em `OrderRead` — **re-medir após Task 0** | §4 classe 3 |
| `F-04` | KPI sobre página vira KPI sobre conjunto, **como conserto de classe** (3 ocorrências) | §4 classe 5 |
| `F-09` | API entrega a linha pronta para `/precos` e `/mercado`. **Bloqueado por `D-17`(fiscal)** | §4 classe 6 |
| `F-10` | `ListingTypeID` sai de `domain`; 3 vazamentos em `integrations` idem (ADR-C4) | §4 classe 6 |

### Onda 3 — visibilidade

| id | o que | medida |
|---|---|---|
| `F-A4` | `as_of` em `/pedidos` e `/precos` (as outras 6 telas já mostram) | §8 |
| `F-A5` | `ProtocoloPage.tsx:98-100` — `onError` que engole falha de retry sem feedback. O outro catch mudo (`PedidosPage.tsx:62`) **não** entra aqui: é o mesmo sítio do truncamento e pertence ao `F-04`, para não haver dois donos do mesmo arquivo | §8 |

## 8. Desconhecido na tela

ADR-17 no frontend **passa**. Zero colapsos de desconhecido→zero nas telas roteadas. Todo `?? 0` opera sobre contador ou `length` já resolvido (`PedidosPage.tsx:185` é fallback morto — `difalPending` já filtrou `amount != null` em `:183`). `MarginChip.tsx:44` e `produto/EstoqueTab.tsx:54-64` declaram desconhecido explicitamente.

Os três estados (carregando / vazio legítimo / falhou) são distinguíveis em 8 das 10 páginas roteadas. Parciais: `precos/PricingPage.tsx` (sem `EmptyState`), `integracoes/SyncHealthCard.tsx` (sem estado "vazio" — não se aplica).

`degraded` é lido em **um** lugar: `dashboard/DashboardPage.tsx:18-20,34,68`. Distingue corretamente "KPI nulo por fonte degradada" (erro no card) de "nulo por fonte saudável" (`—` honesto, `formatKpi:14-16`).

Furos:

- `pedidos/PedidosPage.tsx:62` — truncamento só no console
- `mutations/ProtocoloPage.tsx:98-100` — `onError` só reseta o ref; botão volta clicável sem mensagem
- `/pedidos` e `/precos` não mostram data de sincronização; as outras 6 mostram (`AnunciosPage.tsx:456`, `MercadoPage.tsx:141,192`, `ProdutoHeader.tsx:22`, `EstoqueTab.tsx:57`, `SyncHealthCard.tsx:74-82`, `dashboard/DashboardPage.tsx:59-64`)

## 9. Cadeia de prova

Cada `F` prova a corrente inteira ML → DTO → domínio → banco → API → tela. Elo sem prova é elo que já falhou em silêncio nesta casa.

Mapeamento nas lanes ratificadas (`contracts/governance/execution-lanes.json`), sem inventar níveis novos:

| elo | lane | controle negativo nomeado |
|---|---|---|
| payload do provider → DTO | `unit` | fixture sem o campo, prova que vira desconhecido e não zero |
| DTO → domínio | `unit` | valor errado no fixture, prova que a asserção reprova |
| domínio → banco | `integration` (tag `//go:build integration` nas 5 primeiras linhas) | must-fail com `failure_token=test=` citado |
| banco → API | `integration` | schema shape, não só presença |
| API → tela | `browser` | live-drive dirigido pelo operador |
| leitura Oracle | `live-oracle` | somente leitura, `METALPRD` |
| leitura do provider | `live-provider-read` | controle positivo injetado |

`provider-write` não é usada por nenhuma fatia planejada.

**Critério de aceite falsificável.** Todo critério é **um comando com saída esperada** ou **uma string na tela**. Nunca um adjetivo. Isso é o que permite revisão barata (§10) — um revisor pequeno consegue verificar "roda X, espera Y"; não consegue verificar "está correto".

Regra herdada e medida nesta casa: verde de integração só é evidência **depois que o vermelho nomear o teste**. Asserção de presença não pega valor errado.

## 10. Ritual de execução

Por fatia. A metodologia de onda que governa quando cada fatia é planejada está na §12.

1. **Plano profundo** (Opus) — lê o código, os contratos, o SDK, o OpenAPI e os legados antes de decidir. Aproveita o que existe; move para o máximo global mesmo que exija refactor
2. **Implementação**
3. **Spec review** (haiku) — o critério de aceite é verificável mecanicamente por construção
4. **Quality review** (sonnet) — diff contra a fatia
5. **Browser** — live-drive dirigido pelo operador
6. **Evidência em `.mnfs`** — não escrito, não aconteceu

### Checklist anti-atalho

Nomeada porque cada item já mordeu esta casa pelo menos uma vez:

- Nunca hardcodar valor de negócio na tela
- Nunca fabricar default plausível para fato desconhecido (ADR-17)
- Nunca apresentar agregação sobre resultado truncado como total
- Nunca construir motor novo ao lado de motor que já funciona — medir o existente primeiro
- Nunca alterar OpenAPI sem alterar `sdk-runtime` no mesmo commit
- Nunca aceitar verde sem que o vermelho tenha nomeado o teste
- Nunca deixar vocabulário de provider entrar em `domain` (ADR-C4, verificável por grep)

## 11. Modos de falha assumidos

- **Re-medir após Task 0 pode matar `F-03` inteiro.** Custo: uma rodada. É o preço de não escrever contra estado velho
- **Três campos ficaram indeterminados** (`SellerSKU`, `Title`, `ProviderStatus`). Entram como pergunta de payload respondida por leitura de raw, nunca por adivinhação
- **`F-04` como conserto de classe pode virar refactor de paginação no SDK.** Se virar, é refactor aceito — a alternativa é o mesmo defeito voltar pela quarta vez
- **`F-06` toca produção compartilhada** (`http_handler.go` divide struct com `WithCalc`). Não é deleção limpa; a fatia começa medindo, não apagando
- **`F-08` pode revelar que uma das 11 operações é necessária.** Nesse caso ela sai da poda e vira `F` com tela — o critério é a linha de chegada, não a contagem

## 12. Metodologia de execução

### 12.1 Uma onda por vez, porque medição decai

Ondas não são detalhadas de uma vez. **Onda N vira tasks; ondas N+1 em diante ficam no nível da linha do `BACKLOG.md`** — id, enunciado, medida, critério, sem decomposição.

O motivo é medido, não estilístico. A Onda 0 altera o código que a Onda 1 mediria; o `F-06` apaga código que o `F-01` e o `F-02` senão tocariam. Plano escrito contra estado que vai mudar é plano podre. Precedente registrado: o brief do M-06 continha três alegações falsas sobre o repo poucos dias depois de escrito, e a verificação feita para este design encontrou mais uma (a coluna `raw_provider_ref` existe; a suposição era de ausência).

**Re-medição na fronteira de cada onda** é passo obrigatório, não zelo. A Task 0 já é um caso disso.

Dentro da onda, paralelo × serial decide a matriz de colisão do harness. `F-06`, `F-07` e `F-08` tocam `pricing` — dono único, serial.

### 12.2 Alinhamento na fronteira, não no meio do código

Cada onda abre com uma sessão curta de decisão, sobre uma lista conhecida **antes** de começar. Decisões já identificadas:

| onda | decisão | quem decide |
|---|---|---|
| 0 | ~~`F-A3`: job periódico ou manual~~ — **RESOLVIDA 2026-08-03: job periódico e idade visível, as duas** | operador |
| 1 | `F-08`: das 11 operações, quais morrem e quais ganham tela — 11 decisões independentes | operador |
| 1 | `F-06`: `decompose` e `BuildProfitability` permanecem separados (simulação × realizado) ou convergem | operador, com medição de paridade na mesa |
| 2 | `F-03`: escopo real após a re-medição pós-merge | derivado da medição |

### 12.2.1 Sequenciamento contra o P2.b em voo

Decisão do operador em 2026-08-03, medida contra `git diff main...p2b-imposto-ex-ante` (`be8fc56c`, 21 arquivos, +2611 linhas).

O branch do P2.b reescreve exatamente os arquivos da Onda 1:

- `pricing/domain/decompose.go` (+122) — arquivo central do `F-06`
- `internal_read/adapters/oracle/sync.go` (+95) — arquivo do `D-17` fiscal, bloqueio de `F-06` e `F-09`
- `pricing/ports/tax_matrix.go`, `pricing/adapters/postgres/matrix_reader.go` — portas e adapters novos no mesmo módulo
- `migrations/0093_icms_matrix.sql` — a colisão que a Task 0 renumera para `0097`

Além de colidir, o P2.b pode **encerrar sozinho** o `D-21` (SIMPLES 4% fabricado) e o `D-38` (campo `Imposto`), os dois bloqueios do `F-06`.

~~A Onda 0 é disjunta: `F-A1` e `F-A2` vivem em `integrations/`, `F-A3` em `market/`, `F-00` em `sync/`+`orders/`. Zero arquivos em comum com o branch.~~ **CORRIGIDA 2026-08-03**: a medição de `F-A3` mostrou que a causa raiz da idade invisível vive em `packages/web-query` (compartilhado por toda a Onda 0, não confinado a `market/`) — a Fatia A tocou `formatRelativeAge`/`FreshnessIndicator` ali, não só em `market/`. Operador aprovou alargar o escopo declarado em 2026-08-03; a disjunção de arquivos com o branch do P2.b permanece válida (`packages/web-query` não aparece no diff do P2.b), só a alegação "confinado a `market/`" era falsa.

**Regra de sequenciamento:** a Onda 0 é planejada e executada imediatamente. A Onda 1 só é planejada depois de o P2.b aterrissar na `main`, e o primeiro passo do seu plano é **re-medir `decompose.go`** — a medição da §4 classe 6 e da §6 vale para o estado de 2026-08-03 e expira no merge.

Descoberta imprevista durante a implementação continua existindo — é o que os eventos `ESCALATION`/`REQUEST` do harness carregam. A regra é que **decisão previsível não vire interrupção**.

### 12.3 Ficha de medição — gate de entrada no plano

Nenhuma fatia entra no plano da onda sem oito respostas, **cada uma com `file:line`**:

| # | pergunta | defeito que essa pergunta teria pego |
|---|---|---|
| 1 | Quem **já faz isso**? Busca por conceito, não por nome | 4 fórmulas de margem divergentes |
| 2 | Quem **consome**? Zero call-sites de produção = candidato a morte, não a conserto | `batch_orchestrator`, 11 operações do contrato, `ClassifyRefreshError` |
| 3 | Qual **contrato** muda? OpenAPI e `sdk-runtime` no mesmo commit; qual selo ADR-C3 | os 120 campos `required`+`nullable` |
| 4 | Quem é o **produtor real** do valor? | `OrderRead.currency` sem produtor; `nf_state` com produtor errado |
| 5 | Que **teste** pega regressão — e ele é forte ou vacuoso? | `pedidosFormatters.ts` sem teste; dois fracos em `pricing` |
| 6 | Que **dívida registrada** toca isso? | `D-16` bloqueando `F-00`; `D-21`/`D-38` bloqueando `F-06` |
| 7 | O que o **ML realmente entrega**? Payload lido, não documentação | `sale_fee` por unidade; multiget de shipments inexistente |
| 8 | Qual é o **local maximum**? O que já existe se chama como? | `MarginChip` com limiar configurável, ao lado do hardcode |

A ficha é verificável sem julgamento de mérito: basta conferir que as oito respostas existem e carregam `file:line`. Por isso a revisão de spec pode ser barata.

### 12.4 Regra anti-redundância

> Nenhuma fatia cria cálculo, helper ou endpoint novo sem citar o `file:line` do que já existe e **por que não serve**.

Verificável por grep no próprio plano. As quatro fórmulas de margem e os três donos da faixa nasceram exatamente pela ausência dessa linha — cada autor escreveu de boa-fé o que já existia ao lado.

### 12.5 Divisão de trabalho

| papel | quem | por quê |
|---|---|---|
| medir | subagentes sonnet, read-only, em paralelo | barato; preserva o contexto de quem decide |
| sintetizar e decidir | Opus | a decisão precisa segurar as medições juntas — foi o cruzamento entre duas varreduras discordantes sobre `batch_orchestrator` que virou risco declarado em vez de premissa errada |
| conferir a ficha | haiku | critério mecânico (§12.3) |
| revisar o diff | sonnet | escopo da fatia |
| dirigir o browser | operador | §9, lane `browser` |

### 12.6 Custo assumido

Medição é barata porque é subagente. O custo está no plano da onda, e é deliberado: **cada onda começa mais devagar e não volta atrás**. O contra-exemplo está registrado — o plano anterior do P2.b tinha 17 tasks construindo um segundo motor fiscal ao lado do `pricing/domain` que já funcionava; uma rodada de medição derrubou 10 delas.
