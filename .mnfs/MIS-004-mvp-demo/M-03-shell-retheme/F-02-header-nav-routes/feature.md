# F-02-header-nav-routes

```yaml
id: F-02
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-03
created: 2026-07-17
updated: 2026-07-17
validation_level: QA-0
lifecycle_scope: feature
```

## Mission

MIS-004 mvp-demo.

## Milestone

M-03-shell-retheme.

## Brief

Substituir sidebar dark por header horizontal do design: logo, pills de navegação canônicas, seletor de installation (InstallationContext existente — reusar, nunca recriar), toggle de tema, menu ⚙. Criar a indireção de rotas por área (IC-05): `apps/web/src/routes/<area>.tsx` — seam que permite às waves B/C trocar conteúdo de tela sem tocar AppRouter.

## Inputs

- IC-05 §Nav + §Route indirection (fonte única: pills, áreas, ownership).
- Design handoff (header/nav specs) + leitura R-02 `research/design-screens-2026-07-17.md`.
- `apps/web/src/app/**` atual: AppRouter.tsx (rotas PT-BR + placeholders W1), Layout, LegacyRedirects, InstallationContext.
- F-01 tokens/`useTheme`.

## Expected Output

- Header: pills `Visão geral` (/), `Anúncios` (/anuncios), `Mercado` (disabled "em breve"), `Simulador` (/precos), `Pedidos` (/pedidos), `Repasses` (disabled "em breve"); ⚙ menu EXATO IC-05: `Configurações` (com item `DIFAL` = deep-link `/precos?params=1` — abre drawer de parâmetros do M-07; até M-07 existir, navega p/ /precos placeholder), `Integrações`, `Catálogo`, `Estoque` — SOMENTE estes 4 itens (IC-05); `Vínculos` fica FORA da navegação global (nem pill, nem ⚙): entrada primária = tela Anúncios, rota /vinculos permanece registrada e alcançável por deep-link; seletor installation; toggle tema. Pill ativa reflete rota atual (incl. deep-link). Rotas /classifications, /marketplaces, /vinculos etc. permanecem registradas e alcançáveis mesmo sem item de menu.
- `apps/web/src/routes/{dashboard,anuncios,vinculos,produto,precos,pedidos}.tsx`: cada um exporta o componente da área — HOJE re-exportando a página/placeholder existente. AppRouter passa a importar SÓ destes arquivos p/ as áreas listadas; screen-milestones (M-04..M-09) editam apenas seu `routes/<area>.tsx` + `pages/<área>/**`.
- Rotas existentes preservadas: paths idênticos aos do W1 (nenhuma URL muda), LegacyRedirects intactos.
- EARS: While usuário em /precos, when página carrega via deep-link + F5, the sistema shall renderizar header com pill Simulador ativa. While pill disabled é clicada, when clique ocorre, the sistema shall não navegar e manter affordance "em breve" (tooltip/label). While installation troca no seletor, when seleção ocorre, the sistema shall propagar via InstallationContext existente (mesmo comportamento atual).

## Negative Scenarios

- Rota desconhecida ⇒ comportamento 404/redirect atual preservado.
- Viewport estreito (janela demo redimensionada) ⇒ header não colapsa em estado inutilizável (overflow com scroll horizontal das pills é aceitável; sumir item não é).

## Interaction Model

Nav é stateless (deriva da URL). Tema via `useTheme` (F-01). Installation via contexto existente — ownership NÃO muda. Nenhum fetch novo no shell.

## Constraints

- IC-05: AppRouter NÃO será editado pelas waves B/C — este F entrega o seam que garante isso.
- Nomes/ordem das pills = IC-05, não improvisar.
- Catálogo/Estoque/Integrações/Protocolos continuam acessíveis (⚙), telas atuais sem retheme profundo aqui.

## Ownership

- Owned paths: `apps/web/src/app/**` (Layout/header/AppRouter/LegacyRedirects), `apps/web/src/routes/**` (novo).
- Forbidden paths: `packages/ui/**` (F-03), `apps/web/src/pages/**` (conteúdo das telas), theme module (F-01).
- Parallel-safe with: F-03 (disjoint: app/routes vs packages/ui).

## Validation Expectations

- Screenshot header light+dark com pill ativa correta em 3 rotas distintas.
- Transcript navegação: cada pill habilitada navega p/ path exato; pills disabled não navegam (URL imutável no clique).
- Deep-link + F5 em /vinculos, /precos, /pedidos ⇒ tela da área renderiza via `routes/<area>.tsx` (placeholder W1 ainda OK nesta fase).
- Grep prova: AppRouter importa áreas SOMENTE de `routes/` (zero import direto de pages p/ as 6 áreas).

## Execution Artifact Rules

`spec.md`, `plan.md`, and `validation.md` are created during feature execution, not mission planning.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer (chip M-03, após F-01).
- Next action: criar `spec.md`.
- Required files/evidence: `validation.md` + screenshots + grep.
- Blockers or open decisions: none.
