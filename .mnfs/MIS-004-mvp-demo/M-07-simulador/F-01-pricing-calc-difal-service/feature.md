# F-01-pricing-calc-difal-service

```yaml
id: F-01
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-07
created: 2026-07-17
updated: 2026-07-17
validation_level: QA-0
lifecycle_scope: feature
```

## Mission

MIS-004 mvp-demo.

## Milestone

M-07-simulador.

## Brief

Backend do simulador: CalcProfile persistido, tabela DIFAL 27 UFs com seed, motor de decomposição ÚNICO (fórmula IC-04), simulação bidirecional, cenários salvos, port `DifalForUF` p/ M-08. Extensão aditiva de `/pricing/*`.

## Inputs

- IC-04 (`research/pricing-difal-interface-contract.md`) — fórmula, perfis, seed DIFAL, limiares, endpoints, port (fonte única; NÃO redecidir alíquota/fórmula).
- R-04 (`research/difal-interna-rates-2026.md`) — dataset citado das alíquotas internas (modal) das 27 UFs p/ o seed `interna_pct` (MS = verify-at-execution; override do operador esperado).
- `/pricing/simulations` + batch existentes (preservar).
- Custo: Reader `GetCostAsOf` (IC-02); frete: `GetFreeShippingCost`, comissão: dados listing/categoria existentes; componentes indisponíveis ⇒ unknown-propagation.
- Migrations bloco 0055–0059. Toda tabela nova (`pricing_calc_profiles`, `pricing_difal_rates`, `pricing_scenarios`) carrega `tenant_id`; toda query nova escopa `tenant_id` (invariante da missão).

## Expected Output

- Tabelas: `pricing_calc_profiles` (regime `SIMPLES` **default** 4% | `PRESUMIDO` 9,25%; limiares default 18/10; `tarifa_full` nullable = desconhecido; `difal_enabled`; `difal_destino_uf`; overrides), `pricing_difal_rates` (27 UFs seedadas com `interna_pct` do R-04 + `interestadual_pct` 12% MG/PR/RJ/RS/SC/SP, 7% demais; origem SC; efetivo = max(interna−inter, 0); `origem_versao: "padrao-2026"`; override persiste SÓ quando Δ>0,049pp, com auditoria), `pricing_scenarios`.
- Motor: `Decompose(input) → {preco, comissao, taxa_fixa, frete, imposto, difal, tarifa_full, custo, margem_valor, margem_pct, componentes_desconhecidos[]}` — fórmula EXATA IC-04 (taxa fixa 6,50 se preço<79; frete produto se ≥79); modalidades `classico|premium|full`: `tarifa_full` é componente explícito nullable, debitado só em `full` (0 explícito nas demais; null em `full` ⇒ componente desconhecido ⇒ propaga); QUALQUER componente desconhecido ⇒ margem desconhecida + lista de faltantes.
- Bidirecional: `SolveTargetPrice(margem_alvo)` via busca binária (degraus em 79; convergência exata); margem inatingível ⇒ **200 `UNREACHABLE_TARGET`** citando o teto atingível.
- Endpoints: `GET/PUT /pricing/profile`, `GET /pricing/difal`, `PUT /pricing/difal/{uf}` (rota EXATA IC-04; override persiste só Δ>0,049pp, gravado + auditoria), `POST /pricing/simulations` estendido aditivamente (aceita decomposição completa + direção margem→preço), `GET/POST/DELETE /pricing/scenarios`.
- Port Go `DifalForUF(uf) → {efetivo_pct, versao}` — assinatura congelada IC-04 (M-08 consome read-only).
- Seção OpenAPI `/pricing/*` aditiva + `sdk-runtime/src/pricing.ts`.
- EARS: While todos os componentes conhecidos, when simulação roda, the sistema shall retornar decomposição com soma EXATA (preço = Σ componentes + margem). While custo desconhecido, when simulação roda, the sistema shall retornar decomposição parcial + `margem: null` + `componentes_desconhecidos: ["custo_erp"]`. While margem-alvo inatingível (ex.: 60% com comissão 16%), when solver roda, the sistema shall responder 200 `UNREACHABLE_TARGET` citando o teto atingível.

## Inputs/Outputs

Shapes/valores: IC-04 (referência). Codes EXATOS IC-04: UF inválida ⇒ **404 `UF_NOT_FOUND`**; alíquota fora de 0–35% ⇒ **422 `INVALID_RATE`**; margem inatingível ⇒ 200 `UNREACHABLE_TARGET`; simulação com preço ≤ 0 ⇒ 422 `INVALID_PRICE`; simulação de item inexistente ⇒ 404 `ITEM_NOT_FOUND`; cenário inexistente ⇒ 404 `SCENARIO_NOT_FOUND`.

## Negative Scenarios

- UF fora das 27 ⇒ 404 `UF_NOT_FOUND`.
- Override DIFAL com alíquota fora de 0–35% ⇒ 422 `INVALID_RATE`.
- Simulação com preço ≤ 0 ⇒ 422 `INVALID_PRICE` (IC-04).
- Perfil no estado inicial ⇒ regime default `SIMPLES` 4% APLICADO (IC-04), com origem `default` explícita na resposta (usuário vê que é default, não escolha dele).
- Modalidade `full` com `tarifa_full` null ⇒ componente desconhecido propagado (margem desconhecida), nunca 0.

## Constraints

- Label obrigatória em toda resposta DIFAL: "seed padrão 2026 — não é orientação fiscal" (campo `disclaimer`).
- Motor é a ÚNICA implementação da fórmula — M-08 consome via port; duplicação = defeito.
- Precisão: decimal (não float64 p/ dinheiro) conforme padrão do repo.
- `/pricing/simulations` atuais continuam respondendo shape antigo (aditivo).

## Ownership

- Owned paths: `modules/pricing/**`, `apps/server_core/migrations/0055*–0058*` (+ fixture `apps/server_core/internal/platform/migrate/runner_test.go` bump), seção `/pricing/*` OpenAPI, `sdk-runtime/src/pricing.ts`.
- Forbidden paths: `modules/orders/**`, `modules/market/**`, `modules/connectors/**` (consumo via ports), barrel SDK, `apps/web/**` (F-02).
- Parallel-safe with: none — F-02 depende da seção OpenAPI deste F.

## Validation Expectations

- Golden tests IC-04: ≥10 casos com valores decimais EXATOS (SIMPLES/PRESUMIDO × preço </≥79 × UF 12%/7% × custo conhecido/desconhecido) — cada componente da decomposição conferido.
- Transcript solver: margem-alvo 15% ⇒ preço cuja re-simulação retorna 15,00% exato; caso 200 UNREACHABLE_TARGET com teto citado.
- Seed: SELECT das 27 UFs com `interna_pct` (R-04) + `interestadual_pct` exatos e efetivo computado; transcript override com Δ>0,049pp ⇒ persistido + auditado; Δ≤0,049pp ⇒ NÃO persistido.
- Regressão: simulação antiga (shape W1) responde idêntica.

## Execution Artifact Rules

`spec.md`, `plan.md`, and `validation.md` are created during feature execution, not mission planning.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer (chip M-07). Complexidade: candidata a Sol low (solver + motor decimal — hub decide no dispatch).
- Next action: criar `spec.md`.
- Required files/evidence: `validation.md` + golden tests.
- Blockers or open decisions: none.
