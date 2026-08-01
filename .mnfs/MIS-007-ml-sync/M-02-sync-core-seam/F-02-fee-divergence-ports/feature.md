# F-02-fee-divergence-ports

```yaml
id: F-02
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

Ports de núcleo (provider-agnóstico — ADR-01) + implementações Postgres:
- `ChannelFeeWriter.UpsertFee` / `ChannelFeeReader.ResolveListingFees` (IC-01): upsert
  current-value na chave natural; resolução comissão LEDGER-ONLY camada 2 → 1 →
  ausente TIPADO (o degrau `config`/`pricing_tariff_defaults` NÃO mora aqui — é composição
  do consumidor pricing, M-07 F-01; auditoria P5 F-10: dois donos tornavam o braço
  honest-absent do M-07 inalcançável); frete camada 2 → honesto-desconhecido;
  camada 3 NUNCA na resolução; recusa de camada 3 fee_kind=`commission` sem
  `detail.sale_fee_unit`/`quantity` — o mandato de detail é SÓ da comissão; camada 3
  fee_kind=`freight` do shipment (M-06 F-02) não tem essa decomposição e detail NULL é
  aceito (IC-01; auditoria P5 r06 F-r06-1).
- `DivergenceRecorder.Evaluate` / `DivergenceReader.ListOpenByEntity` (IC-02): one-open-row
  upsert, detected_at imutável, auto-resolve, tolerâncias (estoque 0; tarifa R$0.01),
  recusa sem observed_at dos 2 lados, sem-vínculo → não avalia.

EARS (amostra):
- While row camada 2 existe p/ o listing, when ResolveListingFees, the reader shall
  retornar {value, value_type, currency, layer:2, detail, origem, coletado_em} — detail =
  jsonb VERBATIM da row (IC-01 §Required Outputs; F-r06-2) — nunca o config.
- While nenhuma row de ledger existe p/ o listing, when ResolveListingFees, the reader
  shall retornar ausente TIPADO (nunca zero, nunca default) — o consumidor decide o
  fallback.
- While row aberta existe e valores convergem (≤ tolerância), when Evaluate, the recorder
  shall gravar resolved_at e NÃO abrir row nova.

## Inputs

IC-01/IC-02 (semântica binding completa); 0086/0087 do F-01; idioma upsert
`internal_read/adapters/mirror/writer.go:74-95` (`upsertSQL`) + keep-absent `:104-112`
(`keepAbsentSQL` — ranges medidos, P5 r06 F-r06-5).

## Expected Output

Packages novos de núcleo (fees/divergences) com interfaces + impl postgres + testes de
contrato. Consumidores (M-03..M-08) importam interfaces, nunca impl.

## Constraints

- Núcleo NÃO importa tipo ML (teste da fronteira Q6).
- Tenant SEMPRE no WHERE (tudo tenant-scoped).
- Sem HTTP; sem consumidor ligado ainda (M-05/M-06/M-07 ligam).
- Nomes finais de package/interface = spec; SEMÂNTICA = IC, não renegociável.

## Inputs/Outputs

In/Out por operação: IC-01/IC-02 `## Operations` + `## Canonical Examples` (não re-declarar
aqui — referência é binding).

## Negative Scenarios

- Evaluate sem expected_observed_at → erro nomeado, nunca default now() (IC-02).
- UpsertFee camada 3 fee_kind=`commission` sem detail obrigatório → erro nomeado (IC-01);
  camada 3 fee_kind=`freight` sem detail → ACEITO (F-r06-1 — o teste cobre os DOIS lados).
- Resolve p/ listing sem NENHUMA row de ledger → resultado "desconhecido" tipado,
  nunca zero (fallback config é do consumidor M-07, fora deste port).

## Ownership

- Owned paths: packages novos (`apps/server_core/internal/modules/channelfees/` e
  `apps/server_core/internal/modules/divergences/` — nomes finais no spec, layering AGENTS
  obriga `modules/`) + seus testes.
- Forbidden paths: migrações (F-01), scheduler (F-03), pricing existente (M-07), adapters.
- Parallel-safe with: none — depends on F-01 (tabelas 0086/0087).

## Validation Expectations

- Testes de contrato cobrindo TODOS os EARS acima + exemplos canônicos dos ICs (incl.
  rejeições) — asserções por valor de campo, não por "não deu erro".
- Round-trip: UpsertFee 2× mesmo subject → 1 row, coletado_em avançou.
- Divergência 2 direções em teste de integração hermética (criar → row aberta; convergir →
  resolved_at).

## Execution Artifact Rules

Execução cria spec/plan/validation.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer.
- Next action: `spec.md` após F-01 merged.
- Required files/evidence: `validation.md`.
- Blockers or open decisions: none.
