# F-02-anuncios-ui-sinais

```yaml
id: F-02
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-05
created: 2026-07-17
updated: 2026-07-17
validation_level: QA-0
lifecycle_scope: feature
```

## Mission

MIS-004 mvp-demo.

## Milestone

M-05-anuncios-sinais.

## Brief

Estender AnunciosPage (existente, código vivo do W1) com os sinais do design: coluna VS MERCADO, chips de exceção clicáveis no header, toggle agrupar-por-produto, seção de evidência no drawer. Retheme da página aos tokens M-03. Regressão zero no comportamento atual.

## Inputs

- Design handoff (tela Anúncios) + leitura R-02 `research/design-screens-2026-07-17.md`.
- `sdk-runtime/src/listings.ts` (F-01) + `signal_status`/`market_signal` shapes.
- `apps/web/src/pages/anuncios/**` atual (filtros/tabela/drawer existentes — mapear antes de editar).
- IC-05 — tokens, MarginChip, FreshnessIndicator; TanStack Query.

## Expected Output

- Coluna VS MERCADO: posição + delta vs price_to_win (mono font); estados: `SEM_VINCULO` chip cinza com link p/ /vinculos; `NO_PRICE_EVIDENCE` "—" + tooltip; `STALE` valor + FreshnessIndicator âmbar.
- Header: chips de exceção com contadores (summary) — clique aplica filtro `?exception=` na URL.
- Toggle "agrupar por produto" (usa endpoint by-product existente) — linhas de grupo com agregado por produto.
- Drawer: seção Evidência (fonte, fetched_at, freshness, `match_status` do vínculo, tamanho da amostra `n_offers`/`n_sellers`, link produto) — amostra pequena visível, nunca escondida (IC-03).
- EARS: While chip exceção clicado, when filtro aplica, the sistema shall refletir na URL e a tabela mostrar só os anúncios da exceção (deep-linkável). While toggle agrupar ativo, when dados carregam, the sistema shall agrupar por CODPROD mantendo colunas de sinal por anúncio dentro do grupo. While anúncio SEM_VINCULO, when linha renderiza, the coluna shall linkar direto p/ /vinculos (nunca número fabricado).

## Negative Scenarios

- Summary falha mas lista ok ⇒ chips ocultos + tabela normal (erro isolado).
- Filtro de exceção na URL com valor inválido ⇒ ignorado + URL limpa (sem crash).
- Grupo com 1 anúncio ⇒ renderiza como grupo normal (sem colapso especial).

## State / Interaction Model

- Filtros/agrupamento/exceção = URL search params (deep-link + F5 restauram tudo); padrão já existente na página — ESTENDER, não substituir.
- Query keys existentes preservados; novos: `['listings','summary']`. Invalidação no refresh manual existente cobre ambos.
- Drawer: padrão atual da página mantido (item na URL).

## Constraints

- Comportamento atual (filtros, ordenação, paginação, drawer, refresh) intocado — mudanças ADITIVAS.
- Retheme = tokens/classes M-03; sem cor hardcoded.
- Nenhum fetch fora do SDK.

## Ownership

- Owned paths: `apps/web/src/pages/anuncios/**`, `apps/web/src/routes/anuncios.tsx`.
- Forbidden paths: `apps/web/src/app/**`, `packages/ui/**`, outros routes/pages, `sdk-runtime/**`.
- Parallel-safe with: none — depends on F-01 (`listings.ts`) + M-03 seams.

## Validation Expectations

- Screenshot tabela com os 4 estados de sinal visíveis simultaneamente (seed QA) + header com contadores.
- Transcript: clique chip exceção ⇒ URL com `?exception=` + tabela filtrada; F5 mantém.
- Regressão C08/C09 (deep-links W1) passa; filtros/paginação atuais funcionam idênticos (vídeo curto ou sequência de screenshots).
- Toggle agrupar ⇒ grupos por CODPROD visíveis.

## Execution Artifact Rules

`spec.md`, `plan.md`, and `validation.md` are created during feature execution, not mission planning.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer (chip M-05, após F-01 OpenAPI).
- Next action: criar `spec.md` (mapear componentes atuais ANTES de editar).
- Required files/evidence: `validation.md` + screenshots + regressão.
- Blockers or open decisions: none.
