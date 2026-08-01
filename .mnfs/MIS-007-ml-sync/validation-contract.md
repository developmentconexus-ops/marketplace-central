# Mission Validation Contract — MIS-007-ml-sync

```yaml
id: MIS-007-VC
type: mission-validation-contract
status: draft
owner: Mission Strategist
parent: MIS-007
created: 2026-08-01
updated: 2026-08-01
validation_level: QA-0
lifecycle_scope: mission
base_sha: dd89d4b3
```

Verdicts binários. Evidência = caminho inspecionável concreto (core §5). Tipos: `ran`
(executado, output salvo), `assumed` (design, não executado), `could-not-run` (bloqueado —
nomear). Nenhum seam contra dependência real provado por stub sem autorização explícita do
operador. Critérios de milestone NÃO são duplicados aqui — cada `M-*/validation-contract.md`
é vinculante por si; este contrato prova o ESTADO FINAL integrado da missão.

## Mission ID

MIS-007

## QA Level

QA-0

## Required Final State

Operação diária ML confiável, mirror-first: zero chamada ML no caminho de read das telas
(4 sítios mortos — A/B orders enrich, C/D pricing degrau-3); `listings` completa (E3 +
estoque + camada 2) com backfill retomável + scheduler diário; `orders` 12 meses + incremental
5min + webhook `orders_v2` com decomposição de margem persistida (custo congelado);
divergências (estoque, tarifa) detectadas no ingest e visíveis; fee sem proveniência morto
(cascata 2→1→config com origem na tela); saúde do sync visível em /integracoes sem SQL.

## Criteria

## Criterion: Zero ML no caminho de read (os 4 sítios mortos)
ID: MIS07-C1
Level: Mission
Type: Architecture
Required: Yes
Status: Pending
Evidence:
- Command: `grep -rn "tarifflive" apps/server_core/` + execução do guard allowlist
  (M-02 F-04) na main pós-merge de M-07 + contador de calls do adapter ML durante
  `GET /orders/{id}` e simulação /precos (log de adapter no dev stack)
- Expected: 0 hits de tarifflive; allowlist com 0 entradas dos sítios A/B/C/D; contador ML = 0
  durante read de pedido e simulação com ledger populado
- Actual:
- Artifact:
Blocking failure: qualquer GET vivo ML disparado por rota de read de tela (orders detail,
listings list, pricing simulate com camada 2 presente)
Blocking failure observed: No
Owner: QA Validator

## Criterion: Q1 — telas <2s com dados reais
ID: MIS07-C2
Level: Mission
Type: Performance
Required: Yes
Status: Pending
Evidence:
- Command: browser QA — /pedidos e /anuncios com a conta real sincronizada; medir load
  (DevTools network, DOMContentLoaded→dados renderizados) 3 amostras cada
- Expected: as 3 amostras <2s por tela (estatística ÚNICA — P7 r01 B-5); nenhuma request
  ML no waterfall
- Actual:
- Artifact:
Blocking failure: qualquer request ML no waterfall, ou qualquer das 3 amostras ≥2s
Blocking failure observed: No
Owner: QA Validator

## Criterion: Q3 — pipeline de orders ponta-a-ponta no live
ID: MIS07-C3
Level: Mission
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: live-drive final (conta ML real): backfill 12m medido; pedido real com
  decomposição conferida contra o pedido no painel ML; webhook disparado por evento REAL
  (design §9)
- Expected: backfill completo sem duplicata (SELECT count por provider_order_id distinto);
  margem do pedido real consistente (liquido = receita_bruta − comissao_total −
  frete_seller − custo_produto; comissao = sale_fee POR UNIDADE × qty — R-4); evento real
  aparece na tela em segundos via inbox→worker→IngestOrder
- Actual:
- Artifact:
Blocking failure: duplicata de domínio, margem divergente do painel ML sem explicação
nomeada em `incompleto[]`, ou webhook real que não produz ingest
Blocking failure observed: No
Owner: QA Validator

## Criterion: Divergência provada nas 2 direções (as 2 espécies)
ID: MIS07-C4
Level: Mission
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: drive de divergência: plantar estoque ERP≠ML (M-05) e tarifa camada2≠observada
  (M-06); depois convergir os dois
- Expected: row aberta com timestamps dos 2 lados NOT NULL (ADR-10) em cada espécie; badge ⚠
  na tela; após convergência, `resolved_at` preenchido e badge some — nas DUAS espécies
- Actual:
- Artifact:
Blocking failure: divergência sem auto-resolve, resolve sem divergência real, ou badge
persistindo após convergência
Blocking failure observed: No
Owner: QA Validator

## Criterion: Q2 — webhook não injeta dado; PII contida
ID: MIS07-C5
Level: Mission
Type: Security
Required: Yes
Status: Pending
Evidence:
- Command: notificação forjada plausível no endpoint público (M-08 drive); grep de schema
  (comando RODÁVEL — sintaxe corrigida P7 r02 A-14):
  `grep -rn "raw" apps/server_core/migrations/00{86..95}*.sql` (brace expansion bash) +
  inspeção de colunas de buyer fiscal (0089)
- Expected: forja → 200 + row inbox + ZERO escrita em tabela de domínio + IP off-allowlist
  logado; nenhuma migração da missão adiciona `raw jsonb` a tabela com dado
  fiscal/comprador (ADR-03/R-6); colunas fiscais tipadas SÓ
- Actual:
- Artifact:
Blocking failure: forja que alcança tabela de domínio, ou raw de billing_info persistido
Blocking failure observed: No
Owner: QA Validator

## Criterion: Q6 — fronteira núcleo×adapter intacta
ID: MIS07-C6
Level: Mission
Type: Architecture
Required: Yes
Status: Pending
Evidence:
- Command: revisão de gate sobre o diff integrado da missão: imports de tipos ML
  (`grep -rn "mercadolivre\|mercado_livre"` em `internal/modules/*/domain` e
  `*/application` dos módulos de núcleo tocados)
- Expected: payloads de provider ficam nos adapters; núcleo (domain/application de
  channelfees/divergences/orders/listings) não importa tipo ML
- Actual:
- Artifact:
Blocking failure: tipo de provider importado em domain/application de módulo de núcleo
Blocking failure observed: No
Owner: QA Validator

## Criterion: Fee com proveniência em toda superfície
ID: MIS07-C7
Level: Mission
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: browser QA /precos (origem + coletado_em + ⚠ config), /anuncios (coluna TARIFA
  com honest-unknown `—`) e /pedidos (drawer de decomposição — termos de fee com
  origem/coletado_em, P7 r01 B-3); SELECT de channel_fees por camada com origem
- Expected: nenhum valor de fee na tela sem origem; origem ∈ vocabulário IC-01
  (api_listing_prices | api_order | api_shipment | config); anúncio sem fee observada = `—`,
  nunca 0
- Actual:
- Artifact:
Blocking failure: valor de fee sem origem, origem fora do vocabulário, ou 0 fabricado
Blocking failure observed: No
Owner: QA Validator

## Criterion: Q4 — saúde do sync visível (dirigido em browser)
ID: MIS07-C8
Level: Mission
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: browser QA /integracoes seção "Saúde do sync" com o estado REAL pós-missão;
  conferência contra SELECT de sync_state. Descarrega TAMBÉM o critério diferido M09-C6
  (re-drive pós-lanes: listings/orders acendem na seção SEM commit no M-09 — endpoint
  entidade-agnóstico provado)
- Expected: entities (products/listings/orders) com timestamps verdadeiros (== SELECT);
  bloco webhook com última notificação real do live-drive
- Actual:
- Artifact:
Blocking failure: tela divergente do banco, ou campo sem observação renderizado como
zero/ok
Blocking failure observed: No
Owner: QA Validator

## Criterion: Contratos de milestone todos PASS
ID: MIS07-C9
Level: Mission
Type: QA
Required: Yes
Status: Pending
Evidence:
- Command: leitura dos 9 `M-*/validation-result.md`
- Expected: 9/9 verdict PASS por QA Validator (só QA passa milestone — HARNESS-CORE)
- Actual:
- Artifact:
Blocking failure: qualquer milestone sem validation-result PASS
Blocking failure observed: No
Owner: QA Validator

## Evidence Requirements

- Evidência de UI exige screenshot/transcript de browser QA salvo em
  `docs/design/evidence/` (PII scrub antes de commit — pendência conhecida em
  `docs/design/evidence/ml-api/`).
- Evidência de contagem/consistência exige output de SELECT salvo, não afirmação.
- Live-drive final requer conta ML real do operador; NENHUM write ML (registro de callback
  incluso) sem autorização explícita do operador registrada.

## Blocking Failures

- Chamada ML em caminho de read de tela = blocking (MIS07-C1/C2).
- Dado fabricado (zero/default onde há desconhecido) em qualquer superfície = blocking
  (MIS07-C7/C8; ADR-17/AGENTS).
- PII raw persistido = blocking (MIS07-C5).
- Duplicata de domínio em replay/retomada = blocking (MIS07-C3).

## Retry Policy

- correction_attempts: 0
- max_correction_attempts: 2
- last_validation_result: n/a

## Handoff

- Current status: draft (P6; vira ativo quando missão `planned` pós-P7).
- Next owner: QA Validator (verdito final de missão, após 9/9 milestones PASS).
- Next action: P7 dual gate; execução por lanes A→B→C→D.
- Required files/evidence: `validation-result.md` (missão); `M-*/validation-result.md` (9);
  evidência viva em `docs/design/evidence/`.
- Blockers or open decisions: none.
