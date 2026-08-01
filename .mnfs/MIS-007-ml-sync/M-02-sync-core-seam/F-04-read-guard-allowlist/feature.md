# F-04-read-guard-allowlist

```yaml
id: F-04
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-02
created: 2026-07-31
updated: 2026-07-31
validation_level: QA-0
lifecycle_scope: feature
```

## Mission

MIS-007 ml-sync.

## Milestone

M-02 sync-core-seam.

## Brief

Guard allowlist encolhente (ADR-05): teste arquitetural que enumera TODO sítio read-time
que pode tocar o client ML, como allowlist explícita de exatamente 4 entradas
(`GET /orders` enrich, `GET /orders/{id}` buyer-fiscal, `POST /pricing/decompose`,
`POST /pricing/solve` — `research/codebase-read-side.md:20`). Sítio NOVO fora da lista =
falha imediata com mensagem que nomeia o sítio. Milestones que matam sítios (M-03: A/B;
M-07: C/D) deletam a própria entrada no MESMO commit; lista vazia no fim da missão = Q1
medido, não esperado.

EARS:
- While allowlist tem as 4 entradas atuais, when guard roda na main de hoje, the teste
  shall PASSAR (baseline verde no merge deste feature).
- While worker adiciona chamada ML num handler interativo novo, when guard roda, the teste
  shall FALHAR nomeando package/símbolo do sítio novo.
- While entrada é deletada junto com o sítio, when guard roda, the teste shall passar com a
  lista menor.

## Inputs

`research/codebase-read-side.md` (4 sítios com file:line); mecanismo: análise de imports /
call-graph nos packages de transport interativo (técnica exata = spec; critério binding =
detectar dependência do client ML a partir de handlers interativos).

## Expected Output

Teste arquitetural novo (package de teste dedicado) + allowlist como dado no próprio teste
(não config externa — mudar a lista = commit revisável).

## Constraints

- Guard NÃO pode passar vácuo: lição CHIP-IMPORT-CHAIN — observável que passa nos dois
  mundos não é evidência. O must-fail (sítio injetado) é parte do DONE deste feature.
- Não bloquear caminhos batch/worker (ingest é ML por design) — o guard mira classe
  interativa.

## Negative Scenarios

- Falso-positivo (import indireto legítimo por package batch) → guard distingue por
  classe de rota/package, documentado no teste.

## Ownership

- Owned paths: arquivo(s) de teste arquitetural novo(s) em server_core.
- Forbidden paths: os 4 sítios em si (morrem em M-03/M-07, não aqui); handlers.
- Parallel-safe with: F-01, F-03 (eixo files).

## Validation Expectations

- Baseline: guard verde na main com 4 entradas.
- Must-fail: injetar chamada ML fake em handler interativo → guard falha NOMEANDO o sítio
  (evidência com output do teste vermelho).

## Execution Artifact Rules

Execução cria spec/plan/validation.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer.
- Next action: `spec.md`.
- Required files/evidence: `validation.md` com o par verde/vermelho.
- Blockers or open decisions: none.
