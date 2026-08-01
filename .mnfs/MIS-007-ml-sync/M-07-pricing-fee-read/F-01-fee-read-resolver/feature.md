# F-01-fee-read-resolver

```yaml
id: F-01
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-07
created: 2026-07-31
updated: 2026-07-31
validation_level: QA-0
lifecycle_scope: feature
```

## Mission

MIS-007 ml-sync.

## Milestone

M-07 pricing-fee-read.

## Brief

Resolver novo `pricing/adapters/feeledger/resolver.go` implementando `ports.TariffResolver`
(`tariff_resolver.go:30-32` — assinatura existente, request
`{Modalidade, ProductID, PriceBasis}` `:14-24`): resolve comissão via ChannelFeeReader
(IC-01) em cascata camada 2 (subject=listing do anúncio vinculado ao ProductID) → camada 1
(subject=category) → fallback `pricing_tariff_defaults` (resolver existente
`pricingtariffdefaults.NewResolver` root.go:837 REUSADO como base da cadeia — degrau
`config`). Cada resolução carrega proveniência: origem IC-01 + coletado_em
(`domain.ComponentResolution` estendida ADITIVA — campo Fonte já existe, ganha
OrigemLedger/ColetadoEm; spec mapeia contra o domain atual). MORTE no mesmo diff:
`adapters/tarifflive/` deletado, `adapters/tariffcomposite/` deletado ou reduzido (a
cascata nova substitui), fiação root.go:845-851 (o bloco `if feeReader, ferr := ...`;
a declaração `var tariffResolver` em `:844` SOBREVIVE — range medido, P5 r06 F-r06-5) vira
`tariffResolver := feeledger.NewResolver(channelFeeReader, calcTariffResolver)` (sem o
`if feeReader` guard —
resolver novo não depende de ML vivo), allowlist -2 (C/D) com must-fail do guard.
(`baseline_commission_percent: 0.16` em `auth_adapter.go:42-48` é METADATA de catálogo do
provider — contrato publicado wiki/OpenAPI/SDK, SEM call site em pricing; fica INTOCADA —
auditoria P5 r02 N-2.)

EARS:
- While anúncio vinculado tem camada 2, when Resolve roda, the comissão shall vir do ledger
  (detail.percentage_fee/fixed_fee da camada 2) com origem api_listing_prices + coletado_em.
- While ledger vazio p/ o produto, when Resolve roda, the comissão shall vir de
  tariff_defaults com origem `config` — mesmo comportamento de HOJE no caminho defaults
  (13.00/16.00 materialize-on-read `calc_repository.go:239-269` INTOCADO).
- While o store de defaults FALHA (única lacuna real do degrau-4 — `Resolve` do
  tariffdefaults é total: materialize-on-read garante a row; auditoria P5 r02 N-3), when
  Resolve roda, the resultado shall ser erro TIPADO — NUNCA uma constante embutida.

## Inputs

IC-01 §resolution (ordem 2→1→config binding); fato #6 de `research/p5-prerequisites.md`
(cadeia atual verbatim — tarifflive
não implementa o port, composite sim, fiação :837-856); vínculo produto→anúncio
(LinkedListings, leitura por porta); M-02 ChannelFeeReader.

## Expected Output

Resolver novo + deleções + re-fiação root.go (região pricing, merge arbitrado pelo hub) +
extensão aditiva de ComponentResolution.

## Constraints

- `ports.TariffResolver` INTOCADO (os DOIS consumidores — calcSvc `:853-855` e batchOrch
  `:852` — recebem o novo sem mudança).
- pricing_tariff_defaults write path (`PUT /pricing/tariff-defaults`,
  `calc_handler.go:461-478`) INTOCADO — operador segue editando defaults.
- Deleção de teste alheio (testes do tarifflive): restaurar OU provar por observável a
  cobertura equivalente no resolver novo (memória `deleted-test-restore-or-observable`).
- Camada 1: nesta missão ninguém ESCREVE camada 1 (nenhum milestone a popula) — braço
  implementado + testado por fixture, documentado como adormecido até fonte futura.

## Inputs/Outputs

Cascata + origens verbatim IC-01. Degrau numérico: ledger=camada respectiva, defaults=
`config` (vocabulário IC-01; mapeamento p/ campo Degrau legado pinado na spec).

## Negative Scenarios

- Produto sem vínculo → pula direto p/ config (sem erro).
- Fee observada stale (coletado_em velho) → resolve mesmo assim, proveniência carrega o
  timestamp (quem julga é a tela/F-02).
- ProductID nil no request → config (comportamento atual do defaults preservado).

## Ownership

- Owned paths: `pricing/adapters/feeledger/` (novo), deleções tarifflive/tariffcomposite,
  região pricing root.go `:828-858` + remoção dos imports `root.go:99,101`
  (tarifflive/tariffcomposite — sem eles o arquivo não compila; P5 r04 F-r04-5),
  allowlist (C/D), porta de vínculo (nova, read-only).
- Forbidden paths: channel_fees schema/write; calc_repository defaults; orders; listings.
- Parallel-safe with: none — primeira do M-07.

## Validation Expectations

- Tabela-verdade da cascata: 4 fixtures (camada2 hit, camada1 hit, config, store de
  defaults falhando → erro TIPADO) → origem exata em cada / erro nomeado no 4º.
- Allowlist (M-02 F-04) -2 entradas C/D: teste do guard atualizado; must-fail =
  reintroduzir chamada ML read-time em pricing → allowlist reprova nomeando o sítio.
- Before/after: N simulações reais (single + batch) com delta explicado por origem.

## Execution Artifact Rules

Execução cria spec/plan/validation.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer.
- Next action: `spec.md`.
- Required files/evidence: `validation.md`; before/after.
- Blockers or open decisions: none.
