# Milestone Validation Contract

```yaml
id: M-08
type: milestone-validation-contract
status: planned
owner: Mission Strategist
parent: MIS-004
created: 2026-07-17
updated: 2026-07-17
validation_level: QA-0
lifecycle_scope: milestone
```

## Milestone ID

M-08-pedidos.

## QA Level

Dual gate frio (Opus + Sol medium) + QA live-drive browser fresh (/pedidos).

## Required Outcome

Projeção de orders com enriquecimento de shipment (IC-06) + custo por pedido (`GetCostAsOf`, M-01) + decomposição de retorno via ports M-07 (IC-04) + tela Pedidos completa (KPIs, Fila, Lista, Kanban read-only, drawer).

## Criteria

## Criterion: Projeção orders com dados vivos
ID: M-08-C01
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: GET /orders (estendido, aditivo) pós-sync contra installation real
- Expected: pedidos com shipment enriquecido via ports IC-06 (UF destino, modalidade, custo de frete) — DTO ML morto no adapter; custo do item via `GetCostAsOf` na data do pedido (source time visível); shape antigo preservado (aditivo); campo indisponível ⇒ null/UNKNOWN, nunca 0
- Actual:
- Artifact: `M-08-pedidos/validation-result.md` §projecao (response bodies)
Blocking failure: DTO provider vazando, custo com data errada, ou indisponível como 0
Blocking failure observed: No
Owner: QA Validator

## Criterion: Decomposição de retorno via ports M-07
ID: M-08-C02
Level: Milestone
Type: Architecture
Required: Yes
Status: Pending
Evidence:
- Command: comparar decomposição de um pedido (drawer/API) com POST /pricing/simulations dos mesmos inputs; inspecionar imports do modules/orders
- Expected: valores IDÊNTICOS componente a componente (mesma fonte: ports `Decompose`/`DifalForUF` — fórmula única IC-04); chip DIFAL do pedido usa `DifalForUF(uf_destino_real)` com versão; grep confirma orders NÃO reimplementa fórmula (consumo read-only dos ports)
- Actual:
- Artifact: `M-08-pedidos/validation-result.md` §decomposicao (diff dos valores + grep)
Blocking failure: divergência de componente entre Pedidos e Simulador, ou fórmula duplicada em orders
Blocking failure observed: No
Owner: QA Validator

## Criterion: Tela Pedidos — KPIs, Fila, Lista, Kanban
ID: M-08-C03
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: live-drive em /pedidos
- Expected: KPIs (contagens/somas do período) batem com a Lista filtrada pelos mesmos parâmetros; Fila de atenção lista exceções com deep-link; Kanban por status ESTRITAMENTE read-only (sem drag/drop mutando estado — nenhuma chamada de escrita ao interagir); aba/coluna Cancelados presente como stub DESABILITADO (Non-Scope); light+dark OK
- Actual:
- Artifact: `M-08-pedidos/validation-result.md` §tela (screenshots light+dark + network transcript)
Blocking failure: KPI divergente da Lista, kanban gerando write, ou Cancelados funcional fora de escopo
Blocking failure observed: No
Drive (UI — agent-browser):
- Fixture: stack local com M-02 fechado (orders sync real) + M-07 ports publicados; sem auth
- Steps:
  - open http://localhost:5174/pedidos
  - assert text "Pedidos"
  - assert text <KPI vendas do período>
  - click <primeiro pedido da Lista>
  - assert text "Retorno"
  - assert text "DIFAL"
- Expected: drawer abre com decomposição completa (componentes + margem) e chip DIFAL da UF real do shipment
Owner: QA Validator

## Criterion: Drawer com decomposição honesta
ID: M-08-C04
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: live-drive: drawer de pedido com custo conhecido e de pedido cujo item não tem custo no snapshot
- Expected: custo conhecido ⇒ decomposição completa com margem e chip (verde/âmbar/vermelho pelos limiares); custo ausente ⇒ UnknownValue + `componentes_desconhecidos` visível — retorno NUNCA exibido como se completo; rótulo DIFAL "seed padrão 2026 — não é orientação fiscal" no chip/tooltip
- Actual:
- Artifact: `M-08-pedidos/validation-result.md` §drawer (screenshots dos 2 casos)
Blocking failure: margem fabricada com custo ausente, ou chip DIFAL sem rótulo
Blocking failure observed: No
Owner: QA Validator

## Criterion: Migrações e seams
ID: M-08-C05
Level: Milestone
Type: Engineering
Required: Yes
Status: Pending
Evidence:
- Command: `ls apps/server_core/migrations/ | grep -E '^006[0-4]'` + `runner_test.go`; diff vs ownership
- Expected: ALTERs/projeções `orders_*` SÓ no bloco 0060–0064; fixture = contagem real; diff só em `modules/orders/**`, `sdk-runtime/src/orders.ts` (aditivo), OpenAPI `/orders*` aditivo, `apps/web/src/pages/pedidos/**`, `routes/pedidos.tsx`; trecho decomposição só commitado APÓS ports M-07 publicados (gate intra-wave B respeitado — evidência no ledger do chip)
- Actual:
- Artifact: `M-08-pedidos/validation-result.md` §seams
Blocking failure: migração fora do bloco, write fora do ownership, ou trecho decomposição antes do gate
Blocking failure observed: No
Owner: QA Validator

## Criterion: Máscara PII do comprador (LGPD)
ID: M-08-C06
Level: Milestone
Type: Security
Required: Yes
Status: Pending
Evidence:
- Command: inspecionar response bodies de GET /orders e GET /orders/{id} (mesmos transcripts de C01) + grep no payload por campos de comprador
- Expected: NENHUM campo de comprador não-mascarado em nenhuma resposta — identificação limitada a "primeiro nome + inicial, cidade/UF"; sem CPF/CNPJ, e-mail, telefone, endereço completo ou nome completo em payload, log ou UI
- Actual:
- Artifact: `M-08-pedidos/validation-result.md` §pii (trechos de payload inspecionados)
Blocking failure: qualquer campo de comprador não-mascarado em payload, log ou tela
Blocking failure observed: No
Owner: QA Validator

## Evidence Requirements

- `M-08-pedidos/validation-result.md` com seções projecao, decomposicao, tela, drawer, seams, pii.
- Dual gate antes do live-drive.

## Blocking Failures

Seams declarados: orders sync + shipments = leitura live ML via ports M-02 (IC-06), live-driven; decomposição = ports M-07 (IC-04). Mock não satisfaz C01/C02 sem autorização do operador registrada.

## Retry Policy

- correction_attempts:
- max_correction_attempts: 2
- last_validation_result:

## Handoff

- Current status: planned (P6).
- Next owner: QA Validator (pós-close F-01/F-02 do chip M-08).
- Next action: rodar critérios pós-gate M-07 F-01.
- Required files/evidence: `M-08-pedidos/validation-result.md`.
- Blockers or open decisions: none.
