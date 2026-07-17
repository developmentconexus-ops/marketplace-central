# F-01-theme-tokens-fonts

```yaml
id: F-01
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

Sistema de tema papel+verde: CSS variables (tokens do design handoff), atributo `data-theme` light/dark no root com toggle persistido, fontes Instrument Sans (UI) + IBM Plex Mono (números/código) self-hosted, mapeamento Tailwind → tokens. Fundação de TODO o visual das waves B/C.

## Inputs

- `docs/design/handoff-2026-07/README.md` — paleta papel+verde, tokens, tipografia (fonte única de valores).
- IC-05 §Theme — nomes de vars, mecânica data-theme.
- `apps/web/src/index.css`, tailwind config atuais; sidebar `#0F172A` hardcoded a morrer (R-03).

## Expected Output

- `index.css` (ou módulo de theme importado por ele) com todos os tokens em `:root` + overrides `[data-theme="dark"]`; nenhuma cor hardcoded nova.
- Toggle: estado em `localStorage`, aplicado no boot ANTES do primeiro paint (sem flash), default light.
- Fontes self-hosted em assets (sem CDN — app roda offline na demo), `@font-face` + fallbacks de sistema.
- Tailwind config mapeando classes semânticas → vars (ex.: `bg-surface`, `text-primary`, `text-mono`).
- EARS: While tema dark salvo, when app carrega, the sistema shall renderizar dark do primeiro paint. While operador alterna toggle, when clique ocorre, the sistema shall trocar `data-theme` sem reload e persistir.

## Negative Scenarios

- `localStorage` indisponível/corrompido ⇒ default light, sem crash.
- Fonte falha ao carregar ⇒ fallback de sistema legível (font-display: swap).

## Interaction Model

Toggle vive no header (F-02 posiciona; F-01 entrega hook/provider `useTheme` + atributo). Ownership do estado: DOM attr + localStorage — nenhum estado React duplicado fora do hook.

## Constraints

- Valores de cor/spacing vêm do handoff — NÃO inventar tons.
- Nenhuma mudança de layout/nav neste feature (F-02).
- Telas existentes devem continuar renderizando (tokens substituem valores por baixo; quebras visuais pontuais são aceitas até F-03/retheme, quebras FUNCIONAIS não).

## Ownership

- Owned paths: `apps/web/src/index.css`, theme module novo em `apps/web/src/app/theme/**`, assets de fonte, tailwind config.
- Forbidden paths: `apps/web/src/app/Layout*`/AppRouter (F-02), `packages/ui/**` (F-03), `apps/web/src/pages/**`.
- Parallel-safe with: none — F-02 e F-03 dependem dos tokens deste F (edge F-01→F-02, F-01→F-03).

## Validation Expectations

- Screenshot light + dark da mesma rota mostrando tokens aplicados (papel+verde, sem sidebar dark antiga na parte tematizada).
- Transcript de computed style: `getComputedStyle` de um elemento mostrando var resolvida nos dois temas (valores hex exatos do handoff).
- Reload com dark salvo ⇒ primeiro paint dark (sem flash — verificável por gravação/screenshot sequencial ou assert no boot script).

## Execution Artifact Rules

`spec.md`, `plan.md`, and `validation.md` are created during feature execution, not mission planning.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer (chip M-03).
- Next action: criar `spec.md`.
- Required files/evidence: `validation.md` + screenshots.
- Blockers or open decisions: none.
