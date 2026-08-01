# F-03-backfill-cursor-ingest

```yaml
id: F-03
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-04
created: 2026-07-31
updated: 2026-07-31
validation_level: QA-0
lifecycle_scope: feature
```

## Mission

MIS-007 ml-sync.

## Milestone

M-04 listings-backfill-ingest.

## Brief

`IngestListing` (ADR-04/IC-06) + backfill retomável: enumerador scan/scroll_id coleta ids
(1000/batch) SEM hidratar; hidratação por multiget 20 (M-01 F-02) processa POR BATCH (nunca
acumula o catálogo em memória — `Ingestion.Pull` atual acumula com cap 10k páginas, o
caminho novo NÃO herda isso); upsert listing+variations via F-02; cursor com fases
`{"phase":"backfill","scroll_id":...,"run_started_at":...}` →
`{"phase":"sweep","last_full_sweep_at":...}` — nil PROIBIDO. Hidratação nova CONTINUA
alimentando `AbsorbProviderSnapshots` (âncoras EAN/SKU do re-vínculo — ADR-13).

EARS:
- While cursor phase=backfill com scroll_id, when job roda, the backfill shall retomar do
  scroll_id sem re-enumerar do zero.
- While processo morto no meio da hidratação, when job re-roda, the ingest shall retomar
  sem duplicar rows (upsert idempotente) — contagem final idêntica à do run ininterrupto.
- While enumeração exaure e hidratação drena, when run fecha, the job shall declarar run
  COMPLETO (gatilho da marcação F-02) e retornar cursor terminal sweep.
- While hidratação processa itens, when snapshots observer é alimentado, the âncoras shall
  ser não-regressivas vs pull pré-mudança.

## Inputs

IC-06 (cursores canônicos, binding); IC-07 (grão de estoque); M-01 F-02 (multiget); F-02
deste milestone (writer re-semantizado); `connectors/source.go:17-19,54-89`
(AbsorbProviderSnapshots — caminho a preservar); contrato de cursor M-02 F-03.

## Expected Output

Job de listings registrável no scheduler + serviço IngestListing chamável por id (refresh
usa) + testes de retomada.

## Constraints

- Enumerador NUNCA hidrata; hidratação NUNCA enumera (design §6).
- Batch em memória bounded (tamanho de batch de hidratação, não catálogo).
- Erro por item não aborta run (IC-06); erro de enumerador → run INCOMPLETO, cursor
  preservado no ponto.
- Zero mudança em scheduler.go (contrato M-02).

## Inputs/Outputs

In: cursor JSONB; out: cursor JSONB não-nil (exemplos IC-06 binding). IngestListing: in
provider_listing_id, out erro-ou-nada, efeito = upsert completo.

## Negative Scenarios

- scroll_id expirado no retorno → run INCOMPLETO com erro nomeado; próximo run re-inicia
  backfill do zero (novo run_started_at) — nunca marca absent, nunca perde rows.
- 429 em série além do budget (M-01) → run INCOMPLETO, cursor no ponto.

## Ownership

- Owned paths: `apps/server_core/internal/modules/listings/application/**` (serviço
  novo/estendido), `internal/modules/listings/adapters/connectors/**` (enumerador/
  hidratação novos).
- Forbidden paths: repository.go além do necessário p/ chamar F-02 (writer é de F-02);
  scheduler.go; root.go (wiring é F-04).
- Parallel-safe with: none — depends on F-02.

## Validation Expectations

- Kill-and-resume: matar após batch k → re-rodar → contagem final == run ininterrupto;
  zero duplicata (PK).
- Cursor pós-terminal não-nil assertado no sync_state (valor phase=sweep).
- Fixture >1 página de scan (R-3).
- Must-fail ADR-13 "snapshot observer starved": contagem+conteúdo de âncoras não-regressivos.

## Execution Artifact Rules

Execução cria spec/plan/validation.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer (complexo — state machine + cursor: candidato a Sol low
  per matriz de dispatch).
- Next action: `spec.md` após F-02.
- Required files/evidence: `validation.md`.
- Blockers or open decisions: none.
