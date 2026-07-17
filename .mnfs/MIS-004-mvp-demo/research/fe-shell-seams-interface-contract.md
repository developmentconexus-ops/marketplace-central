# Interface Contract

```yaml
id: IC-05
type: interface-contract
status: planned
owner: Mission Strategist
parent: MIS-004
created: 2026-07-17
updated: 2026-07-17
validation_level: QA-0
lifecycle_scope: support
```

## Boundary

Shell M-03 (produtor) ↔ telas M-04..M-09 (consumidoras) em `apps/web`.

## Why This Contract Exists

4 milestones FE paralelos na wave B tocariam Layout/AppRouter/tema simultaneamente sem este contrato. Padrões vêm do W1 (R-03) + design handoff.

## Theme & Tokens (M-03 entrega; ninguém mais edita)

- Tokens papel+verde do design (`docs/design/handoff-2026-07/README.md` §Design Tokens) como CSS variables; `data-theme="light"|"dark"` no root; default light; toggle no header persiste em localStorage.
- Fontes: Instrument Sans (UI) + IBM Plex Mono (valores numéricos/breakdown), self-hosted.
- Chips de margem (`MarginChip` novo em packages/ui): verde ≥18% · âmbar 10–18% · vermelho <10% (limiares vêm do CalcProfile IC-04 quando em contexto de cálculo).

## Nav canônica (HANDOFF — ratificada P1)

Pills no header: `Visão geral(/)` · `Anúncios(/anuncios)` · `Mercado(desabilitada "em breve")` · `Simulador(/precos)` · `Pedidos(/pedidos)` · `Repasses(desabilitada "em breve")` · `⚙`(menu: Configurações → seção DIFAL read/drawer; Integrações; Catálogo; Estoque). Vínculos FORA da nav global — nem pill, nem item de ⚙ (entrada: tela Anúncios → "importação"; rota /vinculos registrada e alcançável por deep-link — ruling P5-R3-01/r04). Rotas existentes (/catalogo, /estoque, /vinculos, /integracoes, /classifications, /marketplaces) permanecem REGISTRADAS e alcançáveis; LegacyRedirects intactos. Sidebar atual morre.

## Route Namespace (indireção — M-03 cria, cada área troca só o seu)

- `apps/web/src/routes/<area>.tsx` exporta o elemento da rota da área (`vinculos.tsx`, `anuncios.tsx`, `produto.tsx`, `precos.tsx`, `pedidos.tsx`, `dashboard.tsx`). `AppRouter.tsx` importa desses arquivos e NÃO é editado por milestones de tela.
- Enquanto a tela não existe, o arquivo da área exporta `WorkspacePlaceholder`.
- Dono por arquivo: vinculos=M-04 · anuncios=M-05 · produto=M-06 · precos=M-07 · pedidos=M-08 · dashboard=M-09.

## Page & State Patterns (padrão W1 — obrigatório)

- Telas moram em `apps/web/src/pages/<área>/**` (padrão M-02 W1). Packages `feature-*` legacy: só feature-simulator/inventory/products/classifications sobrevivem; NENHUM package feature-* novo.
- Dados: TanStack Query via `packages/web-query` (staleTimes por domínio, invalidation helpers); SDK via `ClientContext`; installation via `InstallationContext` (NUNCA recriar seletor).
- Estados: `LoadingState`/`ErrorState`/`EmptyState`/`UnknownValue`/`FreshnessIndicator` de packages/ui — obrigatórios em toda tela nova (ADR-17 + P1c).
- Sem edição inline em tabela; drawer 300–380px padrão design; kanban read-only.

## Ownership

- M-03 exclusivo: `apps/web/src/app/**`, `apps/web/src/routes/` (criação + index), `packages/ui/src/**` (tokens/tema/primitivas base), `index.css`, fontes.
- Pós-M-03: adição de primitiva nova em packages/ui = REQUEST ao hub (lock aditivo nomeado).
- Tela: só `apps/web/src/pages/<área>/**` + `routes/<area>.tsx` próprios.

## Must Not Decide In Feature Execution

Paleta/tokens, mecânica do tema, layout do header, conteúdo das pills, padrão de rota, localização de páginas.

## Validation Impact

QA visual light/dark por tela; nav canônica com pills desabilitadas visíveis; deep-link + F5 em cada rota; placeholder swap sem diff em AppRouter.tsx.
