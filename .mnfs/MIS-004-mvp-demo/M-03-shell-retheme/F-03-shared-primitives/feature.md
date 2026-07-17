# F-03-shared-primitives

```yaml
id: F-03
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

Primitivas compartilhadas em `packages/ui` nos tokens novos: criar `MarginChip` (estado de margem canônico usado por Anúncios/Simulador/Pedidos/Produto) e retematizar as primitivas existentes (`UnknownValue`, `FreshnessIndicator`, `EmptyState`, `ErrorState`, `LoadingState`, `ConflictTag`) + padrões de tabela/drawer do design.

## Inputs

- IC-05 §Primitives (MarginChip API), IC-04 (limiares default 18/10 — verde ≥18%, âmbar 10–18%, vermelho <10%, neutro desconhecido).
- Design handoff (chips/tabela/drawer) + leitura R-02 `research/design-screens-2026-07-17.md`.
- `packages/ui/src/**` atual (R-03: primitivas existem, estilo antigo).
- F-01 tokens.

## Expected Output

- `MarginChip`: props `{marginPct: number|null, thresholds?: {healthy: number, tight: number}}` — null ⇒ variante neutra "—" com label "desconhecido" (NUNCA 0%); defaults 18/10 sobrescrevíveis por prop (M-07 injeta CalcProfile).
- Primitivas existentes retematizadas: mesma API pública (zero breaking change p/ telas atuais), visual nos tokens.
- Padrões `DataTable` shell + `DetailDrawer` (estrutura/estilo do design; comportamento de dados fica nas telas).
- EARS: While marginPct null, when MarginChip renderiza, the sistema shall mostrar variante neutra sem número. While marginPct=18.0 com defaults, when renderiza, the sistema shall mostrar verde (≥ é inclusivo). While tema alterna, when data-theme muda, the sistema shall re-renderizar primitivas nas cores do tema sem reload.

## Negative Scenarios

- marginPct NaN/Infinity ⇒ tratado como null (neutro), console.error em dev.
- thresholds invertidos (tight > healthy) ⇒ defaults aplicados + warn (nunca chip mentiroso).

## Interaction Model

- Todas as primitivas são STATELESS/presentacionais: zero fetch, zero TanStack, zero contexto de dados — 100% via props.
- `MarginChip`: função pura de `{marginPct, thresholds}` → variante; sem estado interno.
- `DataTable` shell: colunas/linhas/sorting/selection são CONTROLADOS pela página consumidora (props + callbacks); a primitiva não guarda estado de dados nem dispara refetch.
- `DetailDrawer`: aberto/fechado controlado pela página (prop `open` + `onClose`); conteúdo por children; largura 300–380px do design.
- Tema: primitivas leem apenas CSS variables (`data-theme` no root) — re-render automático na troca, sem JS de tema próprio.

## Constraints

- API pública das primitivas existentes NÃO muda (consumidores atuais compilam sem edição).
- Limiar default vem de IC-04 — não inventar terceiro limiar.
- Nenhuma primitiva faz fetch — puramente presentacional.

## Ownership

- Owned paths: `packages/ui/src/**`.
- Forbidden paths: `apps/web/src/app/**` e `routes/` (F-02), `apps/web/src/pages/**`, theme module (F-01 — consome vars, não edita).
- Parallel-safe with: F-02 (disjoint: packages/ui vs app/routes).

## Validation Expectations

- Testes de render MarginChip: 5 casos (null, 25, 18.0, 12, 3) ⇒ variante exata por caso (verde/verde/âmbar-borda.. conforme limiares: 18.0 ⇒ verde inclusivo, 12 ⇒ âmbar, 3 ⇒ vermelho, null ⇒ neutro).
- Build do monorepo verde (prova de API pública preservada).
- Screenshot das primitivas light+dark (storybook/página de dev ou composição de teste).

## Execution Artifact Rules

`spec.md`, `plan.md`, and `validation.md` are created during feature execution, not mission planning.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer (chip M-03, após F-01).
- Next action: criar `spec.md`.
- Required files/evidence: `validation.md` + testes + screenshots.
- Blockers or open decisions: none.
