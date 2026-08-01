# F-01-backfill-incremental

```yaml
id: F-01
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-06
created: 2026-07-31
updated: 2026-07-31
validation_level: QA-0
lifecycle_scope: feature
```

## Mission

MIS-007 ml-sync.

## Milestone

M-06 orders-backfill-decomposition.

## Brief

Enumerador + composição: job de orders (JobFunc `scheduler.go:46` — cursor autoritativo,
NUNCA retornar nil: ciclo sem progresso re-retorna o cursor recebido, doc `:36-45`) que
enumera /orders/search (`sort=date_desc`, cursor `date_last_updated` — fato D-120) em fase
`backfill` (janela 12m particionada: window_until + offset, canonical IC-06 §orders) e, ao
exaurir, TRANSICIONA p/ fase `incremental` (watermark = maior date_last_updated − 10min
overlap). Cada id enumerado → `IngestOrder` (M-03 — enumerador NUNCA persiste, ADR-04).
Segunda instância de Scheduler (padrão `synccomposition.NewProductsScheduler`,
root.go:668-678; ADR-08) com intervalo 5min, entity `orders` (enum 0075 JÁ tem — SEM
migração), composição em arquivo novo do módulo orders + 1 linha ancorada root.go.
Run-complete rule (IC-06): watermark só avança quando enumeração exauriu E todos ids da
janela foram ingeridos (falha parcial → cursor fica na janela, re-run retoma).

EARS:
- While cursor nil (1ª run), when job roda, the fase shall ser backfill janela mais recente
  e o cursor retornado shall ser não-nil SEMPRE.
- While backfill exaure, when job completa, the cursor shall virar
  {"phase":"incremental","watermark":...} e RecordSuccess shall gravar incremental=true
  (fix M-02).
- While ingest de 1 id falha transitoriamente, when janela fecha, the watermark shall NÃO
  avançar além do pedido falho (retomável).

## Inputs

IC-06 (cursor canonical + run-complete, binding); M-03 (IngestOrder port); M-02 F-03 (fix
incremental flag + nil-cursor guard); `scheduler.go:85` RegisterJob (duplicate
registration = boot alto); fatos D-120 (search cursor, 403 terceiros, ~poucos pedidos na
conta real).

## Expected Output

Job + enumerador (application do orders) + composição nova + linha root.go + transição de
fase testada com clock fake.

## Constraints

- Enumerador NÃO acumula ids da janela inteira em memória (lição `Ingestion.Pull` NÃO
  herdada — IC-06): processa página a página.
- 403/404 de id → skip contado no resultado do run, nunca falha o run.
- Overlap −10min: replay de id já ingerido é barato (IngestOrder idempotente M-03) — nunca
  filtrar client-side por "já vi".
- scheduler.go INTOCADO.

## Inputs/Outputs

Cursors verbatim IC-06 §orders (backfill e incremental — spec não re-inventa chaves).

## Negative Scenarios

- Kill no meio da janela → re-boot retoma da MESMA janela (cursor persistido intacto).
- Relógio ML vs local: watermark deriva de date_last_updated do PROVIDER, nunca now() local.

## Ownership

- Owned paths: enumerador + job em `orders/application/`, composição nova
  (`orders/composition` ou synccomposition — spec segue padrão products), 1 linha ancorada
  root.go (região schedulers `:661-678`).
- Forbidden paths: `sync/application/scheduler.go`; IngestOrder internals (consome port);
  transport.
- Parallel-safe with: none — primeira do M-06.

## Validation Expectations

- Teste de transição: backfill 2 janelas fake → incremental; cursores exatos por golden.
- Must-fail nil-cursor: job que retornaria nil → guard M-02 pega (prova que o guard
  sustenta — lição must-fail).
- sync_state row real após RunOnce com incremental=true na fase incremental.

## Execution Artifact Rules

Execução cria spec/plan/validation.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer.
- Next action: `spec.md`.
- Required files/evidence: `validation.md`.
- Blockers or open decisions: none.
