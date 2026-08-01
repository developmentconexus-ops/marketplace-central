# M-04-listings-backfill-ingest

```yaml
id: M-04
type: milestone
status: passed
owner: Mission Strategist
parent: MIS-007
created: 2026-07-31
updated: 2026-07-31
validation_level: QA-0
lifecycle_scope: milestone
```

## Mission

MIS-007 ml-sync — [mission.md](../mission.md); IC-06 (ingest/cursor), IC-07 (listings E3 +
variations).

## Outcome

Catálogo ML completo e fresco no Postgres: backfill scan→multiget retomável por cursor,
MASS-CLOSURE substituído por absent≠closed, colunas E3 + `listing_variations` +
`available_quantity`, scheduler diário registrado, refresh manual em lote re-apontado pro
mesmo caminho. Writer único preservado.

## Why This Milestone Exists

Dono ponta-a-ponta do writer de `listings`. Onda 1 do design, itens 1-2 + estoque; fees e
divergência (itens 3-4) são M-05 — separados p/ este milestone fechar QA-ável sem contrato
FE (ADR-14: lane B não carrega contrato FE).

## Features

| ID | Name | Brief |
| --- | --- | --- |
| F-01 | listings-ddl | [F-01-listings-ddl/feature.md](F-01-listings-ddl/feature.md) |
| F-02 | mass-closure-replacement | [F-02-mass-closure-replacement/feature.md](F-02-mass-closure-replacement/feature.md) |
| F-03 | backfill-cursor-ingest | [F-03-backfill-cursor-ingest/feature.md](F-03-backfill-cursor-ingest/feature.md) |
| F-04 | scheduler-refresh-wiring | [F-04-scheduler-refresh-wiring/feature.md](F-04-scheduler-refresh-wiring/feature.md) |

## Dependencies

- M-01 (multiget + resiliência — backfill é o consumidor; hidratação hoje é N+1 e
  `Ingestion.Pull` acumula em memória).
- M-02 (contrato de cursor F-03; fix incremental; `IngestListing` semântica ADR-04).

## Ownership & Concurrency

- Exclusive surfaces: `apps/server_core/internal/modules/listings/**` (application +
  adapters postgres + connectors), migrações 0090-0092; wiring próprio no root.go
  (constructor local, região ancorada — ADR-14).
- Migration block: **0090-0092** (0090 E3 aditivas + lifecycle, 0091 listing_variations,
  0092 índices).
- Predicted seam locks: sync_state/scheduler = consome contrato do M-02 (sem editar
  scheduler.go); snapshots observer (`AbsorbProviderSnapshots`) = superfície de M-04 mas
  com must-fail de não-regressão (ADR-13).
- Runs in parallel with: M-03 (lane B — módulos orders × listings disjuntos; root.go
  serializado pelo hub).
- Internal feature DAG: F-01 → F-02 → F-03 → F-04 (R-B: semântica de closure ANTES do
  cursor — edge ordenado obrigatório).

## Risks

- **R-B (missão)**: cursor antes da semântica = catalog-wiper no primeiro 429. Edge
  F-02→F-03 é a mitigação; must-fail abort-pós-página-1 obrigatório.
- ADR-13: hidratação nova esfomeia snapshots observer → re-vínculo degrada silencioso.
- Conta real pequena esconde defeito de paginação (R-3) → fixtures >1 página.

## Done Means

- Backfill completo na conta real medido (live-drive do hub); kill no meio → retomada sem
  duplicata; cursor terminal `{"phase":"sweep",...}` persistido não-nil.
- Abort pós-página-1 → ZERO rows flipped closed (must-fail nomeado).
- E3 + variations populadas; `available_quantity` no grão certo (IC-07).
- Âncoras de snapshots não-regressivas vs pull pré-mudança (contagem + conteúdo).
- Scheduler diário registrado (sync_state com run real); refresh em lote 202 async pelo
  MESMO IngestListing.

## Handoff

- Current status: planned.
- Next owner: Milestone Orchestrator (hub) — despacho após M-01+M-02 merged.
- Next action: F-01 primeiro; cadeia serial.
- Required files/evidence: `validation-contract.md` (P6), `validation-result.md`.
- Blockers or open decisions: none.

## Correction Handoff

N/A.
