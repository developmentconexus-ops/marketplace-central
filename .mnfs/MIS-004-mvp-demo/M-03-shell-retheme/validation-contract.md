# Milestone Validation Contract

```yaml
id: M-03
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

M-03-shell-retheme.

## QA Level

Dual gate frio (Opus + Sol medium) + QA live-drive browser fresh (UI milestone — drives obrigatórios, light E dark).

## Required Outcome

Shell retematizado papel+verde: tokens/fontes/data-theme, header com pills canônicas + ⚙ menu IC-05, indireção de rotas por área, primitivas compartilhadas — sem mudar nenhuma URL.

## Criteria

## Criterion: Tokens e tema light/dark
ID: M-03-C01
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: live-drive: carregar /anuncios, alternar toggle de tema, F5
- Expected: `data-theme` alterna light↔dark no root; paleta papel+verde dos tokens (nenhum resíduo do dark antigo hardcoded); preferência sobrevive a F5; fontes do design handoff carregadas (network sem 404 de font)
- Actual:
- Artifact: `M-03-shell-retheme/validation-result.md` §tema (screenshots light+dark)
Blocking failure: cor hardcoded fora de token em superfície do shell, ou tema resetando no reload
Blocking failure observed: No
Drive (UI — agent-browser):
- Fixture: docker dev stack local seedado (`npm run docker:dev`, web :5174; sem auth — tenant fixo)
- Steps:
  - open http://localhost:5174/anuncios
  - assert text "Anúncios"
  - click <toggle tema>
  - assert url ~ /anuncios
  - key F5
- Expected: página recarrega mantendo o tema escolhido; header papel+verde em ambos os temas
Owner: QA Validator

## Criterion: Pills canônicas e estados
ID: M-03-C02
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: live-drive do header em /
- Expected: pills exatas: `Visão geral` (/), `Anúncios` (/anuncios), `Mercado` (disabled "em breve"), `Simulador` (/precos), `Pedidos` (/pedidos), `Repasses` (disabled "em breve"); clique em pill disabled NÃO navega e mantém affordance "em breve"; pill ativa reflete rota atual inclusive via deep-link+F5 em /precos
- Actual:
- Artifact: `M-03-shell-retheme/validation-result.md` §pills (screenshot + transcript drive)
Blocking failure: pill faltando/sobrando, disabled navegável, ou pill ativa errada em deep-link
Blocking failure observed: No
Drive (UI — agent-browser):
- Fixture: mesma do C01
- Steps:
  - open http://localhost:5174/precos
  - assert text "Simulador"
  - click <pill Mercado>
  - assert url ~ /precos
  - assert text "em breve"
- Expected: URL inalterada após clique em disabled; pill Simulador ativa
Owner: QA Validator

## Criterion: ⚙ menu exato IC-05
ID: M-03-C03
Level: Milestone
Type: Functional
Required: Yes
Status: Pending
Evidence:
- Command: live-drive: abrir ⚙ menu
- Expected: SOMENTE 4 itens: `Configurações` (com item `DIFAL` → deep-link `/precos?params=1`), `Integrações`, `Catálogo`, `Estoque`; SEM `Vínculos` no menu nem nas pills; rota /vinculos segue registrada (open direto renderiza a área)
- Actual:
- Artifact: `M-03-shell-retheme/validation-result.md` §gear (screenshot menu aberto)
Blocking failure: item extra/faltando no menu, ou /vinculos desregistrada
Blocking failure observed: No
Drive (UI — agent-browser):
- Fixture: mesma do C01
- Steps:
  - open http://localhost:5174/
  - click <botão ⚙>
  - assert text "Configurações"
  - click <item DIFAL>
  - assert url ~ /precos\?params=1
  - open http://localhost:5174/vinculos
  - assert url ~ /vinculos
- Expected: deep-link DIFAL leva a /precos?params=1 (placeholder até M-07); /vinculos alcançável por URL direta
Owner: QA Validator

## Criterion: Indireção de rotas sem mudar URL
ID: M-03-C04
Level: Milestone
Type: Architecture
Required: Yes
Status: Pending
Evidence:
- Command: inspecionar `apps/web/src/routes/{dashboard,anuncios,vinculos,produto,precos,pedidos}.tsx` + AppRouter; live-drive nas 6 áreas + 1 LegacyRedirect
- Expected: cada arquivo de rota exporta o componente da área (hoje re-export do existente); AppRouter importa SÓ destes arquivos p/ as áreas listadas; paths idênticos ao W1 (nenhuma URL muda); LegacyRedirects respondem como antes
- Actual:
- Artifact: `M-03-shell-retheme/validation-result.md` §rotas (diff + transcript)
Blocking failure: URL mudada, redirect quebrado, ou área importada fora do arquivo de rota
Blocking failure observed: No
Owner: QA Validator

## Criterion: Primitivas compartilhadas prontas p/ waves B/C
ID: M-03-C05
Level: Milestone
Type: Engineering
Required: Yes
Status: Pending
Evidence:
- Command: build + testes de `packages/ui` (`UnknownValue`, `FreshnessIndicator` e primitivas novas do F-03); storybook/preview ou render de teste light+dark
- Expected: primitivas exportadas de `packages/ui/src/**` conforme Interaction Model do F-03; `UnknownValue` renderiza estado desconhecido visualmente distinto de zero; `FreshnessIndicator` exibe idade a partir de fetched_at
- Actual:
- Artifact: `M-03-shell-retheme/validation-result.md` §primitivas (transcript + screenshot render)
Blocking failure: primitiva do contrato F-03 ausente do export, ou UnknownValue indistinguível de 0
Blocking failure observed: No
Owner: QA Validator

## Criterion: Ownership limpo
ID: M-03-C06
Level: Milestone
Type: Engineering
Required: Yes
Status: Pending
Evidence:
- Command: diff do chip vs matriz de ownership
- Expected: writes só em `apps/web/src/app/**`, `apps/web/src/routes/**`, `packages/ui/src/**`, tokens/fontes/index.css; ZERO migração; ZERO OpenAPI; lanes L0–L2 FE verdes
- Actual:
- Artifact: `M-03-shell-retheme/validation-result.md` §seams
Blocking failure: write fora do ownership ou lane vermelha
Blocking failure observed: No
Owner: QA Validator

## Evidence Requirements

- `M-03-shell-retheme/validation-result.md` com seções tema, pills, gear, rotas, primitivas, seams.
- Screenshots light E dark de header + 1 área.
- Dual gate antes do live-drive.

## Blocking Failures

Sem seam de integração externa (FE puro). Stack local é a única dependência de ambiente (hub seam — chip nunca sobe servidor).

## Retry Policy

- correction_attempts:
- max_correction_attempts: 2
- last_validation_result:

## Handoff

- Current status: planned (P6).
- Next owner: QA Validator (pós-close F-01…F-03 do chip M-03).
- Next action: rodar drives no stack local.
- Required files/evidence: `M-03-shell-retheme/validation-result.md`.
- Blockers or open decisions: none.
