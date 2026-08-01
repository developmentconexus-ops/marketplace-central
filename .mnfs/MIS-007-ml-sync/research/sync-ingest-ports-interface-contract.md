# Interface Contract — ingest ports, cursor e schedulers

```yaml
id: IC-06
type: interface-contract
status: planned
owner: Mission Strategist
parent: MIS-007
created: 2026-07-31
updated: 2026-07-31
validation_level: QA-0
lifecycle_scope: support
```

## Boundary

Seam Go interno entre M-02 (define ports), M-03/M-06 (implementam/estendem IngestOrder),
M-04 (IngestListing), M-08 (worker consome IngestOrder), M-09 (lê phase do cursor).
Infra existente preservada: `sync/application/scheduler.go` + `sync_state` 0075.

## Why This Contract Exists

ADR-04/07/08: 4 milestones produzem chamadas ao mesmo caminho de ingest e 2 escrevem
cursores no MESMO scheduler — os traps conhecidos (nil-cursor APAGA; `incremental` sempre
false; um-job-por-entidade fail-closed) viram defeito cross-worker sem pino.

## Resources Or Entities

- Ports (provider-agnóstico, definidos no MÓDULO DONO — decisão P5: `OrderIngestor` no
  M-03 (`orders/ports/`), `ListingIngestor` no M-04 (`listings/ports/`); não há package
  central de ingest):
  - `IngestOrder(ctx, tenant, installation, providerOrderID) error`
  - `IngestListing(ctx, tenant, installation, providerListingID) error`
  (nomes de package/receiver finais = milestone dono; tenant PODE vir por scoping do repo
  do módulo — idioma vigente, ex. `NewRepository(pool, tenantID)` — em vez de parâmetro;
  ASSINATURA semântica é binding: resource-addressed, idempotente, erro por item nunca
  aborta o run inteiro.)
- Enumeradores (adapter): produzem SÓ ids — scan/scroll_id (listings, 1000/batch);
  `orders/search` `date_last_updated` + `sort=date_desc` (bug `date_asc` provado T7).
- Schedulers: DUAS instâncias (orders 5min; listings diário), padrão
  `synccomposition.NewProductsScheduler` (root.go:672-677).

## Operations

| Operation | Trigger | Input | Output | Notes |
| --- | --- | --- | --- | --- |
| JobFunc listings | scheduler diário / disparo manual | cursor atual | cursor novo NÃO-NIL | fases: `backfill` → `sweep` |
| JobFunc orders | scheduler 5min | cursor atual | cursor novo NÃO-NIL | fases: `backfill` → `incremental` |
| Ingest{Order,Listing} | todo produtor (backfill, scheduler, webhook worker, refresh manual) | provider resource id | upsert completo da entidade | caminho ÚNICO (ADR-04); hidratação multiget 20 (listings) / GETs paralelos com budget (shipments) |

## Fields

### Required Inputs / Outputs — cursor JSONB (por entidade, em `sync_state.cursor`)

Listings:

```json
{"phase":"backfill","scroll_id":"eyJ...","ids_collected":2400,"run_started_at":"2026-07-31T03:00:00Z"}
```

terminal →

```json
{"phase":"sweep","last_full_sweep_at":"2026-07-31T03:41:10Z"}
```

Orders:

```json
{"phase":"backfill","window_until":"2025-07-31T00:00:00Z","offset":1200}
```

terminal →

```json
{"phase":"incremental","watermark":"2026-07-31T12:05:00Z"}
```

- **`nil` PROIBIDO como retorno de job** (scheduler.go:42-45 apaga cursor nil → re-backfill
  infinito silencioso). Fase terminal é um cursor VÁLIDO.
- `watermark` = max `date_last_updated` ingerido; cada run incremental consulta de
  `watermark − 10min` (overlap; upsert idempotente absorve o overlap).
- Round-trip JSONB NUNCA comparado byte-exact (lição 0075).

## Enums And Statuses

`phase ∈ {backfill, incremental, sweep}` — vocabulário fechado; IC-05 exibe verbatim.

## Error Cases

- Erro por item no ingest: registra, segue o batch (run não aborta); run com erros ainda
  pode ser COMPLETO p/ fins de enumeração (completude = enumerador exaurido).
- Erro de enumerador/quota: run INCOMPLETO — cursor preservado no ponto, retomada no
  próximo tick; NUNCA marca-absent (ADR-06).

## Error Matrix

| Case | Status | Code | Notes |
| --- | --- | --- | --- |
| — | — | — | seam interno; sem superfície HTTP própria |

## Persistence Expectations

- Registro de jobs: enum 0075 JÁ contém `listings`/`orders` — **NENHUMA migração para
  registrar** (milestone que propuser está errado); `RegisterJob` é um-job-por-entidade
  fail-closed — preservar.
- Fix único do M-02: `incremental` hardcoded false (scheduler.go:160) passa a refletir o
  tipo do run (backfill=false; incremental/sweep=true) — pré-condição IC-05.
- Completude de run (ADR-06): listings ganham `last_seen_at` no upsert; run COMPLETO
  (enumerador exaurido + hidratação drenada) executa marcação
  `absent_since = now() WHERE last_seen_at < run_started_at` — nunca em run incompleto.
- `Read`/`RecordSuccess`/`RecordFailure` + `consecutive_failures` atômico: preservar.

## Canonical Examples

Acima (cursores). Rejeição canônica: job retorna nil → teste de contrato do M-02 FALHA
nomeando "terminal cursor must be non-nil".

## Database Shape

Nenhuma mudança em `sync_state` (0075 intocada; `schedule jsonb` continua NÃO-lido —
cadência dinâmica fora de escopo).

## Seed Data

Nenhum. Fixtures: >1 página de enumeração (R-3); kill no meio da fase backfill.

## Timestamp And ID Semantics

- `run_started_at` no cursor = início do run corrente (base do marca-absent).
- Todos timestamptz UTC.

## Compatibility Rules

- Fase nova = valor aditivo no vocabulário + row nesta tabela de contrato.
- Campos extras no cursor são livres POR entidade (opaco). `phase` é obrigatório e estável
  nos cursores dos jobs NOVOS (M-04/M-06 — escopo ratificado de ADR-07); o cursor legado de
  products (`products_job.go:22-25` — sem campo phase) fica INTOCADO. O scheduler faz parse
  TOLERANTE: `phase` ausente/desconhecida ⇒ `incremental=false` (comportamento de hoje),
  nunca erro — mesmo pino de M-09 F-01 (auditoria P5 r03 P-1).

## Route Namespace

N/A (seam interno). Wiring: cada milestone entra no root.go via SEU constructor de
composição (ADR-14).

## Transport And Integration

Ingest roda fora de request interativo (scheduler/worker); refresh manual em lote = rota
batch existente (202 async), nunca síncrono no request.

## Must Preserve

- `JobFunc` signature (scheduler.go:46); cadence-agnostic `Start`; isolamento de falha por
  entidade; "cursor read error pula o ciclo, nunca fabrica nil".
- Enumerador NUNCA hidrata; hidratação NUNCA enumera (coletar ids primeiro — design §6).
- `Ingestion.Pull` atual acumula catálogo em memória (10k-page cap) — o backfill novo NÃO
  herda esse acúmulo: processa por batch de hidratação.

## Must Not Decide In Feature Execution

- Shape/vocabulário do cursor; proibição do nil; overlap do watermark; regra de completude;
  duas instâncias de scheduler; não-leitura de `schedule`.

## Validation Impact

- Q3: kill-and-resume em CADA entidade — cursor não-nulo pós-terminal, zero duplicata,
  retomada do ponto.
- Must-fail nil-cursor (acima) + must-fail marca-absent em run incompleto (ADR-06/IC-07).
- Webhook worker e scheduler provados contra o MESMO Ingest* (Q3 só descarrega nesse seam).
