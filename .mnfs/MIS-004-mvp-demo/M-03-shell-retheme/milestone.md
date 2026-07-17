# M-03-shell-retheme

```yaml
id: M-03
type: milestone
status: planned
owner: Mission Strategist
parent: MIS-004
created: 2026-07-17
updated: 2026-07-17
validation_level: QA-0
lifecycle_scope: milestone
```

## Mission

MIS-004 mvp-demo — shell novo primeiro; todas as telas herdam.

## Outcome

Shell papel+verde (tokens, data-theme light/dark, Instrument Sans/IBM Plex Mono), header com pills canônicas (HANDOFF: Mercado/Repasses "em breve"; Vínculos fora da nav; ⚙ com acesso a Catálogo/Estoque/Integrações/Config), indireção de rotas por área (IC-05), primitivas compartilhadas (MarginChip + retheme das existentes). Sidebar dark morre.

## Why This Milestone Exists

ADR-07 retheme-first (P1b). W1 já entregou rotas/PT-BR/placeholders/InstallationContext (R-03) — este milestone entrega o VISUAL + seams p/ paralelismo da wave B.

## Features

| ID | Name | Brief |
| --- | --- | --- |
| F-01 | theme-tokens-fonts | `F-01-theme-tokens-fonts/feature.md` |
| F-02 | header-nav-routes | `F-02-header-nav-routes/feature.md` |
| F-03 | shared-primitives | `F-03-shared-primitives/feature.md` |

## Dependencies

Nenhuma (wave A). Insumo: tokens do design (`docs/design/handoff-2026-07/README.md`), IC-05.

## Ownership & Concurrency

- Exclusive surfaces: `apps/web/src/app/**` (Layout/AppRouter/theme), `apps/web/src/routes/**` (novo), `apps/web/src/index.css`, fontes/assets, `packages/ui/src/**`, tailwind config.
- Migration block: none.
- Predicted seam locks: none (é o PRODUTOR dos seams IC-05).
- Runs in parallel with: M-01, M-02.
- Internal feature DAG: F-01 → `F-02 ∥ F-03`.

## Risks

Toque em telas existentes (Anúncios/Catálogo/Estoque/Simulador legacy) ao trocar tokens — mitigação: tokens via CSS vars + classes utilitárias, telas existentes seguem funcionais (QA regressão deep-link C08/C09 do W1); pills desabilitadas não podem quebrar rotas registradas.

## Done Means

Toggle light/dark persiste e cobre todas as rotas; pills canônicas + estados "em breve"; `routes/<area>.tsx` com placeholders funcionando (deep-link + F5); MarginChip e primitivas retematizadas; AnunciosPage/telas existentes navegáveis sem regressão; dual gate + QA visual PASS.

## Handoff

- Current status: planned.
- Next owner: hub (dispatch chip wave A).
- Next action: chip implementa F-01 → F-02∥F-03.
- Required files/evidence: `validation-result.md`; screenshots light/dark.
- Blockers or open decisions: none.

## Correction Handoff

n/a (planning).
