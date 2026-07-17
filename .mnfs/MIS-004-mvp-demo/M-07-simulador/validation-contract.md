# Milestone Validation Contract

```yaml
id: M-07
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

M-07-simulador.

## QA Level

Dual gate frio (Opus + Sol medium) + QA live-drive browser fresh (/precos).

## Required Outcome

Motor de decomposição ÚNICO (fórmula IC-04) + DIFAL seed 27 UFs + simulação bidirecional + cenários + ports `Decompose`/`DifalForUF` p/ M-08 + tela Simulador completa com "aplicar preço" via fila M-03 (teto `previewed`).

## Criteria

## Criterion: Golden tests da fórmula única
ID: M-07-C01
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: rodar golden tests IC-04 (≥10 casos: SIMPLES/PRESUMIDO × preço </≥79 × UF 12%/7% × custo conhecido/desconhecido × modalidade classico/premium/full)
- Expected: cada componente decimal EXATO por caso; soma fecha (preço = Σ componentes + margem); taxa fixa 6,50 se preço<79 senão 0; frete produto se ≥79 senão 0; `tarifa_full` debitada só em `full` (0 explícito nas demais; null em `full` ⇒ componente desconhecido propagado); componente desconhecido ⇒ `margem: null` + `componentes_desconhecidos` nomeando faltantes; precisão decimal (não float64)
- Actual:
- Artifact: `M-07-simulador/validation-result.md` §golden (transcript)
Blocking failure: qualquer componente divergente do valor golden, ou unknown virando 0
Blocking failure observed: No
Owner: QA Validator

## Criterion: DIFAL seed + override
ID: M-07-C02
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: GET /pricing/difal; PUT /pricing/difal/{uf} com Δ>0,049pp e depois Δ≤0,049pp; PUT com alíquota 40; GET UF `XX`
- Expected: 27 UFs ordenadas com `interna_pct` = valores R-04, `interestadual_pct` 12% p/ MG/PR/RJ/RS/SC/SP e 7% demais, efetivo = max(interna−inter, 0), `origem_versao: "padrao-2026"`, campo `disclaimer` = "seed padrão 2026 — não é orientação fiscal"; Δ>0,049pp ⇒ override persistido + auditado; Δ≤0,049pp ⇒ 200 SEM persistir; alíquota 40 ⇒ **422 `INVALID_RATE`**; UF inválida ⇒ **404 `UF_NOT_FOUND`**
- Actual:
- Artifact: `M-07-simulador/validation-result.md` §difal (transcripts + SELECT seed 27 linhas)
Blocking failure: seed divergente do R-04 sem registro, override abaixo do limiar persistido, ou disclaimer ausente
Blocking failure observed: No
Owner: QA Validator

## Criterion: Solver bidirecional
ID: M-07-C03
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: POST /pricing/simulations margem-alvo 15%; re-simular com o preço retornado; POST margem-alvo 60% com comissão 16%
- Expected: preço do solver re-simulado retorna margem 15,00% exata (convergência nos degraus de 79); alvo inatingível ⇒ **200 `UNREACHABLE_TARGET`** citando o teto atingível; shape antigo de /pricing/simulations (W1) preservado (regressão)
- Actual:
- Artifact: `M-07-simulador/validation-result.md` §solver (transcripts)
Blocking failure: round-trip divergente, UNREACHABLE como erro 4xx/5xx, ou regressão no shape antigo
Blocking failure observed: No
Owner: QA Validator

## Criterion: Ports p/ M-08 congelados
ID: M-07-C04
Level: Milestone
Type: Architecture
Required: Yes
Status: Pending
Evidence:
- Command: teste de contrato dos ports Go `Decompose(input)` e `DifalForUF(uf) → {efetivo_pct, versao}`
- Expected: assinaturas EXATAS IC-04 (congeladas — M-08 F-01 consome); `DifalForUF` reflete override ativo; motor é implementação ÚNICA da fórmula (grep por segunda implementação de decomposição em modules/orders = zero)
- Actual:
- Artifact: `M-07-simulador/validation-result.md` §ports (transcript + grep)
Blocking failure: assinatura divergente do IC-04, ou fórmula duplicada fora do pricing
Blocking failure observed: No
Owner: QA Validator

## Criterion: Tela Simulador completa
ID: M-07-C05
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: live-drive em /precos: simulação preço→margem, margem→preço, destino ≠ SP, toggle DIFAL off, produto sem custo
- Expected: decomposição por componente visível com chip margem (verde ≥18 · âmbar ≥10 · vermelho abaixo, limiares do CalcProfile); comparação de mercado exibe source, fetched_at/freshness, n_offers/n_sellers, match_status (IC-03); destino real muda DIFAL; toggle off recalcula sem DIFAL; sem custo ⇒ UnknownValue + veredicto bloqueante SEM_CUSTO (nunca 0); deep-link `/precos?params=1` abre drawer Parâmetros direto; drawer DIFAL lista 27 UFs com disclaimer
- Actual:
- Artifact: `M-07-simulador/validation-result.md` §tela (screenshots light+dark)
Blocking failure: margem desconhecida como 0/verde, comparação sem evidência, ou deep-link não abrindo drawer
Blocking failure observed: No
Drive (UI — agent-browser):
- Fixture: stack local com M-01 (custo real) + M-02 (comissão/frete vivos) fechados; sem auth
- Steps:
  - open http://localhost:5174/precos?params=1
  - assert text "Parâmetros"
  - fill <campo produto/codprod> "<codprod-com-custo>"
  - fill <campo preço> "129,90"
  - fill <campo destino UF> "BA"
  - assert text "DIFAL"
  - assert text "seed padrão 2026"
- Expected: decomposição com DIFAL da BA (efetivo interna−7%); disclaimer visível
Owner: QA Validator

## Criterion: Aplicar preço — fila M-03, teto previewed
ID: M-07-C06
Level: Milestone
Type: Security
Required: Yes
Status: Pending
Evidence:
- Command: live-drive: acionar "aplicar preço" numa simulação; SELECT do protocolo criado; inspecionar log adapter
- Expected: intent `price_update` criado via /mutations com preview+protocolo; status final `previewed` — UI NÃO oferece approve (teto); dispatcher provider OFF confirmado; ZERO requests de escrita ML no log na janela
- Actual:
- Artifact: `M-07-simulador/validation-result.md` §mutacao (SELECT + log excerpt + screenshot protocolo)
Blocking failure: protocolo passando de `previewed`, UI expondo approve, ou write provider
Blocking failure observed: No
Owner: QA Validator

## Criterion: Migrações e seams
ID: M-07-C07
Level: Milestone
Type: Engineering
Required: Yes
Status: Pending
Evidence:
- Command: `ls apps/server_core/migrations/ | grep -E '^005[5-9]'` + `runner_test.go`; diff vs ownership
- Expected: `pricing_calc_profiles`, `pricing_difal_rates`, `pricing_scenarios` no bloco 0055–0059 (brief usa 0055*–0058*; 0059 reserva do bloco); fixture = contagem real; diff só em `modules/pricing/**`, `sdk-runtime/src/pricing.ts`, OpenAPI `/pricing/*`, `apps/web/src/pages/precos/**`, `packages/feature-simulator/**`, `routes/precos.tsx`
- Actual:
- Artifact: `M-07-simulador/validation-result.md` §seams
Blocking failure: migração fora do bloco ou write fora do ownership
Blocking failure observed: No
Owner: QA Validator

## Evidence Requirements

- `M-07-simulador/validation-result.md` com seções golden, difal, solver, ports, tela, mutacao, seams.
- Dual gate antes do live-drive.

## Blocking Failures

Seam declarado: comissão/frete vêm de leituras ML vivas via ports M-02 (IC-06) — live-driven no C05; hardcode de comissão = blocking. DIFAL MS = verify-at-execution (R-04): valor divergente confirmado ⇒ override do operador registrado, não silêncio.

## Retry Policy

- correction_attempts:
- max_correction_attempts: 2
- last_validation_result:

## Handoff

- Current status: planned (P6).
- Next owner: QA Validator (pós-close F-01/F-02 do chip M-07).
- Next action: rodar critérios; publicar ports libera trecho decomposição do M-08 F-01 (gate intra-wave B).
- Required files/evidence: `M-07-simulador/validation-result.md`.
- Blockers or open decisions: none.
