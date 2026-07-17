# F-02-simulador-ui

```yaml
id: F-02
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

Tela `/precos` reconstruída per design (substitui SimulatorPage antiga via `routes/precos.tsx`): picker de anúncio/produto, tabela de simulação com MarginChip, painel de decomposição bidirecional, drawer de parâmetros com labels de origem, toggle DIFAL, cenários salvos, "aplicar preço" ⇒ mutação local SEM execução.

## Inputs

- Design handoff (tela Simulador) + leitura R-02 `research/design-screens-2026-07-17.md` (DIFAL destino real, não SP-hardcoded).
- `sdk-runtime/src/pricing.ts` (F-01); `market.ts` (referência de mercado p/ comparação); envelope `/mutations*` (aplicar).
- IC-04 (decomposição/labels), IC-05 (tokens/MarginChip/seams).

## Expected Output

- Picker: busca anúncio (ou produto vinculado) ⇒ carrega contexto (preço atual, categoria/comissão, custo, UF default de exibição).
- Painel bidirecional: editar PREÇO ⇒ margem recalcula; editar MARGEM-ALVO ⇒ preço recalcula (solver); campo editado por último é a fonte; decomposição completa sempre visível (cada componente + origem).
- Toggle DIFAL on/off com aviso ao desligar ("margem sem DIFAL — não use p/ decisão"); seletor UF p/ simular destino.
- Drawer parâmetros: regime, limiares, alíquota UF com label `seed padrão 2026 — não é orientação fiscal` + origem (seed|override).
- Comparação mercado: se agregado IC-03 disponível, mostrar faixa vs preço simulado COM evidência citada completa IC-03 (`source`, `fetched_at`/freshness, `n_offers`/`n_sellers`, `match_status` do vínculo — estados desconhecidos honestos); sem evidência ⇒ chip NO_PRICE_EVIDENCE.
- Cenários: salvar/recarregar/excluir (nome + inputs).
- "Aplicar preço": cria protocolo tipo existente `price_update` e o leva NO MÁXIMO até estado `previewed` (nunca approve — o executor de `price_update` escreve no ML; approve na demo é PROIBIDO, runbook reforça) + toast com link /protocolos/:id deixando claro "aguardando aprovação — nada foi enviado ao Mercado Livre".
- EARS: While usuário edita preço, when valor muda (debounce), the painel shall re-simular e atualizar decomposição + MarginChip. While componente desconhecido, when decomposição renderiza, the componente shall aparecer como "—" com motivo e a margem como desconhecida (chip neutro). While solver retorna UNREACHABLE_TARGET, when margem-alvo é editada, the painel shall mostrar o teto atingível citado pela API (sem preço fabricado).

## Negative Scenarios

- Anúncio sem custo ERP ⇒ painel funciona com margem desconhecida + call-to-action importar planilha.
- API pricing fora ⇒ ErrorState no painel com retry; picker continua.
- Cenário salvo referenciando anúncio removido ⇒ carrega inputs com aviso "anúncio indisponível".

## State / Interaction Model

- Fonte da verdade da simulação = resposta da API (nunca recalcular no front — motor único IC-04).
- Estado do formulário local; `?item=` na URL p/ deep-link; `?params=1` abre o drawer de parâmetros (deep-link do menu ⚙ Configurações→DIFAL do M-03 — IC-05); cenário carregado substitui formulário inteiro (confirmação se sujo).
- Debounce 300ms nas edições; requests cancelados por edição nova (abort) — última resposta vence, nunca race.
- Keys: `['pricing','simulation',inputsHash]`, `['pricing','scenarios']`, `['pricing','profile']`.

## Constraints

- ZERO cálculo de margem no frontend — exibição apenas (anti-drift do motor).
- Toggle DIFAL é de SIMULAÇÃO (não muda config global).
- Aplicar nunca executa: protocolo para em `previewed` p/ demonstrar governança (pitch); a UI desta tela NÃO oferece approve. Guard-rail de missão: zero writes ML.

## Ownership

- Owned paths: `apps/web/src/pages/precos/**` (rebuild), `apps/web/src/routes/precos.tsx`, `packages/feature-simulator/**` (PricingSimulatorPage legado: retematizar nos tokens ou absorver em pages/precos — package permanece vivo, IC-05 §Page Patterns; superfície da matriz da missão).
- Forbidden paths: `apps/web/src/app/**`, `packages/ui/**`, outros routes/pages, `sdk-runtime/**`, backend.
- Parallel-safe with: none — depends on F-01 (`pricing.ts`) + M-03 seams.

## Validation Expectations

- Screenshot painel com decomposição completa + MarginChip por estado (verde/âmbar/vermelho/neutro).
- Transcript bidirecional: editar preço X ⇒ margem Y da API exibida idêntica; editar margem Y ⇒ preço X' cuja simulação bate.
- Toggle DIFAL off ⇒ aviso visível + componente zerado NA SIMULAÇÃO exibida (API recebe flag).
- "Aplicar" ⇒ protocolo `price_update` em `previewed` visível em /protocolos; NENHUM write executado (preço do anúncio no ML inalterado — conferido no transcript).
- Deep-link `?item=` + F5 restaura simulação.
- Deep-link `/precos?params=1` ⇒ drawer de parâmetros aberto (rota do menu ⚙ Configurações→DIFAL).

## Execution Artifact Rules

`spec.md`, `plan.md`, and `validation.md` are created during feature execution, not mission planning.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer (chip M-07, após F-01 OpenAPI).
- Next action: criar `spec.md`.
- Required files/evidence: `validation.md` + screenshots + transcripts.
- Blockers or open decisions: none.
