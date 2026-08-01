# Milestone Validation Result — M-03-orders-shipment-persist

```yaml
id: M-03-VR
type: milestone-validation-result
status: passed
owner: hub
parent: MIS-007
created: 2026-08-01
validation_level: QA-0
base_sha: 33750f27
chip_tip: d22d3d2076e93182c69a60cd429cfa1c78483105
merge_sha: c5da0d2d
verdict: PASS (6/6 milestone criteria + U1/U2/U3)
```

Sem PII neste arquivo: o drive tocou dados fiscais reais de compradores; nomes, CPF,
endereço e CEP ficam FORA de qualquer artefato. As provas abaixo usam contagens e
booleanos.

## Baseline vs integrado

Baseline pré-merge (main `33750f27`): build OK, vet OK, `go test ./...` = 117 ok / 0 FAIL.
Integrado (`c5da0d2d`): build OK, vet OK, **117 ok / 0 FAIL** — idêntico.

## Estado do drive

Dev stack rebuildado no merge (backend up 19:22). Migrations 0088/0089 aplicadas
(`order_shipments` existe; 10 colunas `buyer_*` em `orders_marketplace_orders`).
`POST /orders/import` contra a conta real (`inst-mercado_livre-7e0d2125…`,
METALNOBREACABAMENTOS): **38 pedidos ingeridos** (limite 50, a conta só tem 38), 65.1s,
`result=200 count=38`. Distribuição: `cancelado` 7 · `enviado` 30 · `faturar` 1.

## Critérios

### M03-C1 — GET /orders/{id} sem NENHUM GET vivo ML — **PASS**

Instrumento = `integration_operation_runs` (uma linha por operação de provider).
3× `GET /orders/2000016827774084`: contagem **405 → 405, delta = 0**.

**Controle positivo, obrigatório para o instrumento não ser cego:** um
`GET .../probes/orders?limit=2` (que comprovadamente chama o ML) no MESMO banco:
**405 → 406, delta = +1**. Logo `delta=0` no read é ausência medida, não instrumento morto.

Reforço estrutural (do chip): `zero_live_ml_call_test.go` prova por AST que
`shipment_reader.go` e `buyer_fiscal_reader.go` NÃO importam o pacote adapter do ML — não
podem chamar o cliente vivo sem o teste ficar vermelho antes.

Fiação conferida em `root.go:611-613`: `ordersEnrichSvc` recebe
`orderspostgres.NewShipmentReader` + `orderspostgres.NewBuyerFiscalReader`; o reader
ML-vivo (`ordersBuyerFiscalReader`, `:583`) só alimenta `ordersIngestSvc`.

### M03-C2 — Persistência completa em uma passada — **PASS (com lacuna nomeada)**

- `order_shipments`: **38 linhas** para 38 pedidos; `sla_limit_at` não-nulo em 38/38;
  `cost_gross` não-nulo em 38/38.
- Colunas fiscais tipadas preenchidas em 38/38 (`buyer_doc_number`, `buyer_address_zip`,
  `buyer_name` todos não-nulos).
- `bucket` derivado com shipment status REAL, não `""`: p.ex. pedido cancelado persistiu
  `bucket=cancelado` com `order_shipments.status=cancelled`; o único `faturar` tem
  `status=ready_to_ship` persistido.
- **Lacuna honesta:** a conta real não tem NENHUM pedido sem fiscal
  (`count(*) filter (where buyer_doc_number is null)` = **0**), então o braço
  "404 honesto → colunas NULL, ingest não falha" **não foi exercitado ao vivo** — só pelo
  teste de unidade. Não bloqueia (nada regrediu), mas fica registrado como não-medido.

### M03-C3 — Allowlist encolhe -2 com must-fail — **PASS (com dívida D-14)**

`mlAllowlist` = 2 entradas (`pricing.decompose.category_resolver`,
`pricing.solve.commission_quoter`). Símbolo A (`newOrdersShipmentReaderAdapter`) **apagado**.
Símbolo B (`newOrdersBuyerFiscalReaderAdapter`) **reclassificado** para `mlExcludedSymbols`
porque ainda alimenta `ordersIngestSvc.BuyerFiscal` do `POST /orders/import` (batch).

Auditado pelo hub antes do merge: `mlExcludedSymbols` é map de nomes EXATOS (sem wildcard),
cada entrada com razão, e `TestRealRepoInteractiveMLSites_MatchesAllowlist` amarra
`raw == len(mlAllowlist) + len(mlExcludedSymbols)` **e** exige que cada excluído ainda seja
achado no scan cru — nem entrada morta nem sítio novo escondido passam.

Resíduo registrado como **D-14** em `.mnfs/HARNESS-DEBTS.md`: a exclusão é chaveada pelo
NOME do construtor, não pela classe da rota consumidora. Religar o resultado do MESMO
construtor num caminho interativo passa calado. Não é regressão (antes ele estava na
allowlist, igualmente permitido), mas a promessa do guard é mais larga que a medida.

### M03-C4 — Writer único — **PASS (com dívida D-15)**

`ImportService` só enumera `provider_order_id` via `ListOrders` e delega cada um a
`IngestOrder`; `IngestOrder` é o único caminho de escrita (todas as buscas acontecem ANTES
da chamada ao store, logo zero escrita parcial). Landmine registrada como **D-15**:
`OrderRepository.UpsertOrders` ainda satisfaz `ports.OrderStore` sem chamador de produção e
sem guard estrutural — writer latente.

### M03-C5 — Truth table intocada — **PASS**

`git diff 33750f27..c5da0d2d -- '*order_bucket*'` = **vazio**. `DeriveOrderBucket` reusado,
não re-derivado.

### M03-C6 — Q1, detalhe de pedido rápido — **PASS**

API `GET /orders/{id}`, 3 amostras: **1.132s** (frio) · **0.203s** · **0.395s**; depois do
warm-up, 3 amostras a **8ms**. Resource timing do próprio browser para
`/orders/2000016827774084`: **39ms**. Todas <2s. Zero request ML no waterfall (ver U3).

## Critérios de user-drive

### M03-U1 — drawer com comprador fiscal + shipment vindos DO BANCO — **PASS (prova trocada)**

Drive real em `http://localhost:5174/pedidos` → clique na linha do pedido
`2000016827774084` → drawer abre com `Rastreio 47247564687 · delivered`, `Destino`,
`Destinatário`, `Frete real R$ 191,80 (sender R$ 73,95)`, `Comprador`, e o bloco
`COMPRADOR · FISCAL (ERP)` com Nome/Razão, Documento (CPF) e Endereço completos — todos
lidos do Postgres.

**Desvio do critério, declarado:** o "matar rede ML no stack" não foi executado — a edição
de `/etc/hosts` dentro do container foi **negada pela política de permissão** desta sessão.
A prova equivalente e mais direta está em C1: contador server-side de operações de provider
com delta=0 no read E controle positivo +1. Isso mede a MESMA proposição (o read não fala
com o ML) sem depender de sabotar DNS.

### M03-U2 — /pedidos lista e summary sem regressão — **PASS (dataset fraco, nomeado)**

KPIs medidos, mesma janela e tenant, API `GET /orders/summary?by=status`:
`{"novo":0,"faturar":1,"enviar":0,"enviado":30}` — idênticos aos KPIs renderizados na tela
(`NOVOS 0 · A FATURAR 1 · A ENVIAR 0 · ENVIADOS 30`) e ao `group by bucket` do banco.

Atribuição pedido a pedido do contrafactual (o que cada pedido receberia com
`shipmentStatus=""`, o braço de `order_repo.go:378`), calculada sobre os 38 pedidos com as
entradas reais de `DeriveOrderBucket`:

- 7 pedidos `provider_status=cancelled` → `cancelado` pelo primeiro branch, **independe** do
  shipment status. Sem mudança.
- 30 pedidos com tag `delivered` → `enviado` pelo branch de tag, que é
  shipment-lookup-independente. Sem mudança.
- 1 pedido (`2000017379792858`, `paid`, sem tag, shipment `ready_to_ship`) → `ready_to_ship`
  não é `shipped|delivered`, cai no gate de faturado → `faturar` nos dois mundos.

**Diferença total = 0, integralmente atribuída.** Tolerância cumprida.

**Fraqueza declarada:** com esta conta, NENHUM pedido depende do shipment status para o
bucket (a tag `delivered` já carrega todos os despachados). O critério passa, mas este
dataset **não discrimina** o braço que a M-03 conserta. Prova forte do braço só com um
pedido `shipped` sem tag `delivered` — não existe na conta hoje.

### M03-U3 — console/network sem ML e sem erro novo — **PASS**

Duas superfícies dirigidas (Fila com drawer aberto; Lista → aba `Enviados`, 30 IDs
distintos renderizados). `performance.getEntriesByType('resource')` filtrado por
`/mercadolibre|mercadolivre/i` = **0** em ambas. Console: **zero erros** nas duas leituras.

## Achado não-bloqueante (fora de contrato)

Primeira tentativa de `POST /orders/import limit=5` deu **500
`CONNECTORS_PROVIDER_TRANSIENT`** em 16.9s; a repetição imediata passou (10.6s, 5/5). Um
transient de UM pedido aborta a corrida INTEIRA — os pedidos seguintes não são tentados.
Não viola nenhum critério da M-03 (o contrato só exige skip tipado para 403/404 no fetch do
pedido), mas um importador em lote que morre no primeiro 5xx do provider é frágil.
Candidato a chip futuro: per-order transient → skip contado, não abort.

## Handoff

- Status: **passed**, merge `c5da0d2d` na main.
- Dívidas abertas por este milestone: D-14, D-15 em `.mnfs/HARNESS-DEBTS.md`.
- Não-medido: braço de fiscal ausente (conta sem caso); braço de shipment status como único
  levantador de bucket (conta sem caso).
