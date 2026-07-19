# Mission Validation Contract

```yaml
id: MIS-004
type: mission-validation-contract
status: draft
owner: Mission Strategist
parent: MIS-004
created: 2026-07-17
updated: 2026-07-17
validation_level: QA-0
lifecycle_scope: mission
```

## Mission ID

MIS-004-mvp-demo.

## QA Level

Doutrina P1b: cada milestone fecha com dual gate frio (Opus + GPT-5.6 Sol medium, concordância) + QA live-drive browser fresh. Fechamento da MISSÃO = rehearsal integrado da demo no docker dev stack local limpo (este contrato). Critérios de milestone vivem nos VCs por milestone (linkados) — não duplicados aqui.

## Required Final State

Jornada completa da demo executável de ponta a ponta no stack local: import .xlsx do cliente → identidade correta → vínculos → sinais competitivos honestos → simulador de margem real com DIFAL → Pedidos funcional — tudo no shell papel+verde, com ZERO writes no Mercado Livre e runbook validado.

## Criteria

## Criterion: Rehearsal integrado da jornada demo
ID: MIS-004-C01
Level: Mission
Type: QA
Required: Yes
Status: Pending
Evidence:
- Command: rehearsal QA live-drive no stack limpo (`npm run docker:dev`, DB re-seedado, installation ML verificada) seguindo o runbook: POST import xlsx exemplo → /vinculos (resolver+batch) → /anuncios (sinais) → /catalogo/produtos/:id (veredicto) → /precos (simulação destino ≠ SP) → /pedidos (drawer decomposição)
- Expected: cada etapa termina no estado que o VC do milestone dono define (import COMPLETED com protocolo `#NNN-E`; vínculo aplicado aparece em Resolvidos; sinal com fonte+fetched_at+n_offers/n_sellers; veredicto um dos estados ADR-06; decomposição soma exata preço = Σ componentes + margem; pedido com chip DIFAL da UF real do shipment)
- Actual:
- Artifact: `.mnfs/MIS-004-mvp-demo/validation-result.md` §rehearsal + screenshots por etapa
Blocking failure: qualquer etapa da jornada exige intervenção fora do runbook, ou exibe zero/default fabricado no lugar de estado UNKNOWN honesto
Blocking failure observed: No
Owner: QA Validator

## Criterion: Prova zero-writes ML
ID: MIS-004-C02
Level: Mission
Type: Security
Required: Yes
Status: Pending
Evidence:
- Command: durante TODO o rehearsal C01, capturar log do adapter connectors (telemetria de rotas) + audit trail de mutations; grep por métodos de escrita contra api.mercadolibre.com
- Expected: zero requests PUT/POST/DELETE ao provider ML no log do adapter (somente GET); protocolo `price_update` criado no rehearsal permanece `previewed` (nunca approved/executed); dispatcher provider confirmado OFF por config
- Actual:
- Artifact: `.mnfs/MIS-004-mvp-demo/validation-result.md` §zero-writes (log excerpt + SELECT status dos protocolos)
Blocking failure: qualquer request de escrita ao ML, ou protocolo price_update fora de `previewed`
Blocking failure observed: No
Owner: QA Validator

## Criterion: Milestones fechados por QA
ID: MIS-004-C03
Level: Mission
Type: QA
Required: Yes
Status: Pending
Evidence:
- Command: inspecionar `M-0X-*/validation-result.md` de cada milestone não-cortado
- Expected: verdict PASS em M-01…M-08; M-09 PASS ou corte registrado no mission close (decisão hub/operador, milestone.md §Risks)
- Actual:
- Artifact: rollup em `.mnfs/MIS-004-mvp-demo/validation-result.md` §milestones (tabela ID → verdict → path)
Blocking failure: milestone não-cortado sem validation-result.md PASS
Blocking failure observed: No
Owner: QA Validator

## Criterion: Observabilidade de coleta em toda UI de preço
ID: MIS-004-C04
Level: Mission
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: no rehearsal, inspecionar TODAS as superfícies que exibem preço de mercado: /anuncios (sinais), /catalogo/produtos/:id (veredicto), /precos (comparação)
- Expected: cada superfície exibe fonte, fetched_at/idade, n_offers/n_sellers, match_status (IC-03 campos obrigatórios); rota flag products/{id}/items com telemetria visível no log quando exercida
- Actual:
- Artifact: `.mnfs/MIS-004-mvp-demo/validation-result.md` §evidencia-preco (screenshot por superfície)
Blocking failure: superfície de preço sem os campos de evidência, ou preço de mercado exibido como número nu
Blocking failure observed: No
Owner: QA Validator

## Criterion: DIFAL fonte única cross-tela
ID: MIS-004-C05
Level: Mission
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: no rehearsal, PUT /pricing/difal/{uf} com override Δ>0,049pp numa UF de teste; re-simular no /precos com essa UF destino; abrir pedido com shipment para essa UF em /pedidos
- Expected: mesmo efetivo_pct refletido no Simulador E no chip Pedidos (fonte única `DifalForUF`); rótulo "seed padrão 2026 — não é orientação fiscal" presente em ambas as superfícies
- Actual:
- Artifact: `.mnfs/MIS-004-mvp-demo/validation-result.md` §difal-cross (screenshots + response bodies)
Blocking failure: valores divergentes entre Simulador e Pedidos para a mesma UF, ou superfície DIFAL sem o rótulo
Blocking failure observed: No
Owner: QA Validator

## Criterion: ADR-17 sem zeros fabricados na jornada
ID: MIS-004-C06
Level: Mission
Type: Architecture
Required: Yes
Status: Pending
Evidence:
- Command: no rehearsal, exercitar ≥3 estados de ausência: produto sem custo (simulação), produto sem evidência de mercado (veredicto), item sem ESTOQUE_RESERVADO (estoque disponível)
- Expected: margem exibe UnknownValue + `componentes_desconhecidos: ["custo_erp"]`; veredicto exibe NO_PRICE_EVIDENCE ou INSUFFICIENT_MARKET (nunca R$0/verde); disponível exibe estado DESCONHECIDO (físico segue consultável)
- Actual:
- Artifact: `.mnfs/MIS-004-mvp-demo/validation-result.md` §adr17 (screenshot por estado)
Blocking failure: qualquer fato operacional desconhecido renderizado como 0, R$0,00 ou default
Blocking failure observed: No
Owner: QA Validator

## Criterion: Security baseline
ID: MIS-004-C07
Level: Mission
Type: Security
Required: Yes
Status: Pending
Evidence:
- Command: `git log --stat` dos merges da missão + grep por `.env` em diffs; inspecionar ao menos 1 query nova POR MILESTONE criador de tabela/ALTER (M-01, M-02, M-04, M-07, M-08) quanto a escopo tenant
- Expected: nenhum arquivo `.env*`/secret/PII em nenhum diff mergeado; queries novas escopam `tenant_id` em TODOS os milestones amostrados (5, não 3); nenhuma superfície nova de auth criada
- Actual:
- Artifact: `.mnfs/MIS-004-mvp-demo/validation-result.md` §security
Blocking failure: secret/PII commitado, ou query nova sem escopo tenant
Blocking failure observed: No
Owner: QA Validator

## Criterion: Maintainability — ladder verde no fechamento
ID: MIS-004-C08
Level: Mission
Type: Engineering
Required: Yes
Status: Pending
Evidence:
- Command: rodar lanes L0–L2 (profile §5) no main pós-merge de todas as waves, `GOCACHE=.gocache`; incluir parity check OpenAPI↔SDK↔handler (smoke L2, ADR-12)
- Expected: L0–L2 exit 0; fixture de migrações `runner_test.go` = contagem real de `apps/server_core/migrations/*.sql`; parity sem drift
- Actual:
- Artifact: `.mnfs/MIS-004-mvp-demo/validation-result.md` §ladder (transcript)
Blocking failure: lane vermelha no main de fechamento
Blocking failure observed: No
Owner: QA Validator

## Criterion: Runbook da demo executado
ID: MIS-004-C09
Level: Mission
Type: Documentation
Required: Yes
Status: Pending
Evidence:
- Command: executar o runbook da demo (preflight: installation ML conectada, coleta na véspera, dry-run planilha real do cliente, proibição de approve em price_update) contra o stack limpo
- Expected: cada passo do runbook conclui com o resultado nele descrito; dry-run da planilha real rejeita ≤10% das linhas (trigger R3) ou desvio registrado com mitigação
- Actual:
- Artifact: `.mnfs/MIS-004-mvp-demo/validation-result.md` §runbook
Blocking failure: passo de preflight sem resultado registrado, ou trigger R3/R6 disparado sem mitigação aplicada
Blocking failure observed: No
Owner: QA Validator

## Criterion: Inventário app-inteiro de telas (design-fidelity sweep)
ID: MIS-004-C10
Level: Mission
Type: QA
Required: Yes
Status: Pending
Evidence:
- Command: no fechamento (pós-merge M-06+M-09), no stack limpo, rodar o live-pass de `SCREEN-INVENTORY.md` §Live-pass protocol — para CADA rota do mapa (jornada + off-journey): (a) computed-style dos elementos de tema (bg/font-family/accent) vs paper+green + Instrument Sans/IBM Plex Mono; (b) a11y tree vs elementos do design R-02 (colunas/chips/seções-drawer/KPIs); registrar present/absent/differs por elemento
- Expected: cada tela da JORNADA da demo com contraparte no design exibe os elementos declarados no design OU a omissão está declarada em feature.md/milestone.md do dono OU tem ruling operador accept-for-demo; nenhuma tela da jornada renderiza off-theme (slate/blue literais) onde o design pede paper+green; placeholders (`/integracoes`,`/marketplaces`) e stubs "em breve" são conhecidos e fora da jornada
- Actual:
- Artifact: `.mnfs/MIS-004-mvp-demo/validation-result.md` §screen-inventory (tabela por tela) + `.mnfs/MIS-004-mvp-demo/SCREEN-INVENTORY.md` (mapa estático base 390d79ab)
Blocking failure: tela DA JORNADA com elemento missing/differs NÃO-declarado e sem ruling accept-for-demo do operador; OU tela da jornada off-theme vs DESIGN-REFERENCE
Blocking failure observed: No
Owner: QA Validator
Note: pixel-diff impossível (rasterizer F-ENV-10) — evidência = computed-style + a11y, limite aceito pelo operador 2026-07-19. Gaps estruturais conhecidos pré-fechamento (Simulador matriz/VEREDICTO, Vínculos tema/tabela, Anúncios grupo-rows) = decisão de escopo do operador ANTES da demo, não auto-fail.

## Evidence Requirements

- `validation-result.md` da missão com seções: rehearsal, zero-writes, milestones, evidencia-preco, difal-cross, adr17, security, ladder, runbook, screen-inventory.
- Screenshots light theme por etapa da jornada (dark coberto nos VCs de milestone).
- Log excerpt do adapter cobrindo a janela completa do rehearsal (C02).
- VCs de milestone: `M-0X-*/validation-contract.md` (linkados, não duplicados).

## Blocking Failures

Qualquer `Blocking failure` acima observado ⇒ verdict FAIL da missão; correção via hub (reserva de migração 0070–0074 só com aprovação), re-rehearsal completo após correção que toque a jornada.

## Retry Policy

Rehearsal integrado: máximo 2 tentativas de correção pós-FAIL antes de escalar ao operador (prazo demo 2026-07-20 — corte de escopo M-09/M-06-abas é a válvula, decisão operador/hub).

## Handoff

- Current status: draft (planning P6).
- Next owner: QA Validator (dispatch pelo hub no fechamento da missão).
- Next action: aguardar milestones fecharem; executar rehearsal.
- Required files/evidence: `.mnfs/MIS-004-mvp-demo/validation-result.md`.
- Blockers or open decisions: corte M-09 (decisão hub/operador na entrada da wave C).
