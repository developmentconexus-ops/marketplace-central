# M-07-pricing-fee-read

```yaml
id: M-07
type: milestone
status: planned
owner: Mission Strategist
parent: MIS-007
created: 2026-07-31
updated: 2026-07-31
validation_level: QA-0
lifecycle_scope: milestone
```

## Mission

MIS-007 ml-sync — [mission.md](../mission.md); IC-01 (resolução 2→1→config, binding),
ADR-05/ADR-09 (morte dos sites C/D).

## Outcome

Pricing lê tarifa do LEDGER, não do ML vivo: resolver novo lê channel_fees em cascata
camada 2 (anúncio) → 1 (categoria) → `config` (pricing_tariff_defaults — sobrevive como
fallback ratificado-honesto, ADR-09); cadeia viva degrau-3 morre
(`tarifflive/resolver.go` deletado + composite simplificado); DTOs de /pricing carregam
PROVENIÊNCIA (origem + coletado_em) e FE /precos mostra origem com ⚠ quando origem=config.
Allowlist (M-02 F-04) encolhe nas entradas C/D. (`baseline_commission_percent: 0.16` em
`auth_adapter.go:42-48` é metadata de catálogo do provider, contrato publicado, SEM call
site em pricing — fica intocada; auditoria P5 r02 N-2.)

## Why This Milestone Exists

Simulação hoje faz GETs vivos no ML por cotação (root.go:845-851) e um miss cai no
degrau-4 `pricing_tariff_defaults` (13.00/16.00) SEM proveniência visível na tela — número
sem origem em tela de decisão de preço é a doença que ADR-09 mata
(`ml-tariff-design-pending`). Pode rodar CEDO na lane C: fallback `config` é honesto mesmo
antes do M-05 popular camada 2 (edge ⤳ soft — qualidade melhora sozinha).

## Features

| ID | Name | Brief |
| --- | --- | --- |
| F-01 | fee-read-resolver | [F-01-fee-read-resolver/feature.md](F-01-fee-read-resolver/feature.md) |
| F-02 | precos-provenance-fe | [F-02-precos-provenance-fe/feature.md](F-02-precos-provenance-fe/feature.md) |

## Dependencies

- M-02 (ChannelFeeReader IC-01 + allowlist F-04).
- M-05 SOFT (camada 2 populada = proveniência boa; sem ela tudo resolve `config` — correto
  e visível).

## Ownership & Concurrency

- Exclusive surfaces: `pricing/**` (adapters tarifflive/tariffcomposite + resolver novo +
  transport DTOs), par OpenAPI+SDK /pricing, FE /precos, região pricing do root.go
  `:828-858` (uma das DUAS exceções à regra linha-nova — a outra é M-03 região orders
  `:576-601`; edita bloco existente; hub arbitra o merge, R-7; F-r08-1).
- Migration block: nenhum.
- Predicted seam locks: channel_fees READ-ONLY (nunca escreve); pricing_tariff_defaults
  read path intocado (`calc_repository.go:239-269` — materialize-on-read fica); contrato FE
  serializado na lane C (ADR-14).
- Runs in parallel with: M-05, M-06 (código ∥; commits de contrato serializados; root.go
  região pricing = merge arbitrado).
- Internal feature DAG: F-01 → F-02.

## Risks

- `tarifflive.Resolver` NÃO implementa `ports.TariffResolver` (satisfaz seam privado do
  composite — fato #6, `research/p5-prerequisites.md`): a deleção muda a fiação
  root.go:845-851 inteira (o `}` de fechamento em `:851` incluso; `var tariffResolver`
  em `:844` SOBREVIVE — F-r07-5, coerente com F-01); o `if
  feeReader...` guard morre junto (resolver novo não depende de capability ML viva).
- Comportamento de simulação MUDA onde antes o degrau-3 vivo resolvia → before/after de
  simulações reais obrigatório (cada delta explicado: live→camada2 ou live→config).
- Batch orchestrator consome o mesmo resolver (`batchOrch.WithTariffResolver` `:852`) —
  os DOIS caminhos (single + batch) testados.

## Done Means

- Grep: zero referência a tarifflive.
- Allowlist (M-02 F-04) encolhe -2 entradas (C/D); must-fail do guard = reintroduzir
  chamada ML read-time em pricing → allowlist reprova nomeando o sítio.
- Simulação com camada 2 presente usa a tarifa observada do anúncio (detail da camada 2;
  proveniência `api_listing_prices`); sem ledger → `config` com origem visível.
- /precos: origem + coletado_em na tela; ⚠ quando origem=config; nunca valor sem origem.
- Before/after de N simulações reais; par OpenAPI+SDK mesmo commit; tsc verde.

## Handoff

- Current status: planned.
- Next owner: Milestone Orchestrator (hub) — lane C (pode ser o primeiro da lane).
- Next action: F-01 spec.
- Required files/evidence: `validation-contract.md` (P6), `validation-result.md`.
- Blockers or open decisions: none.

## Correction Handoff

N/A.
