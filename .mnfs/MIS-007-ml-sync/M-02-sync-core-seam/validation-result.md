# Milestone Validation Result — M-02-sync-core-seam

```yaml
id: M-02-VR
type: milestone-validation-result
status: passed
owner: hub
parent: MIS-007
created: 2026-08-01
validation_level: QA-0
base_sha: 295e293f
merge_commit: 97758970
chip_tip: 7d02bc75
```

## Verdict

**PASS.** 7/7 critérios de milestone + 2/2 de user-drive. Zero blocking failure.
Um HOLD-MERGE emitido pelo hub e sanado antes do merge (ver Adjudicações).

## Evidência por critério

| ID | Veredito | Observável medido |
|----|----------|-------------------|
| M02-C1 | PASS | Lane hermética do hub (`scripts/harness.ps1 -Command integration`, run_id `561a775d34ff48839930638668ce89ed`): `target=ephemeral-postgres`, `migrations_first=76`, `migrations_second=0`, `resource_count=0`, `port=55732`. O `migrations_second=0` é a prova independente de idempotência — medida do lado do hub, não alegada pelo chip. `runner_test.go` fixture 76 == contagem real de `migrations/*.sql`. |
| M02-C2 | PASS | `mis007_f01_core_ddl_must_fail_test.go` contra Postgres real: `23514` em `channel_fees_currency_when_amount_check`; `23505` em `divergences_one_open_row` com aceite após resolve. CHECK provado por INSERT, não por leitura de DDL. |
| M02-C3 | PASS | Suites `channelfees`/`divergences` (domain + postgres) verdes em Postgres real: escada 2→1→ausente-tipada, assimetria comissão/frete da camada 3, tolerância estoque=0 / tarifa R$0,01, one-open-row + auto-resolve, `detected_at` imutável. |
| M02-C4 | PASS | `cursor_contract_test.go` assere a mensagem contendo `"terminal cursor must be non-nil"` — must-fail nomeado, não genérico. |
| M02-C5 | PASS | `products_regression_test.go` roda o `synccomposition.NewProductsJob` REAL através do `Scheduler.RunOnce` REAL e assere `incremental=false` + campos do cursor inalterados. Regressão medida no caminho que produz, não em mock. |
| M02-C6 | PASS | Baseline de 4 sítios reais; fixtures `5th-site` e `aliased_site` falham NOMEANDO o sítio; fixture de allowlist encolhida prova que a lista não está hardcoded. |
| M02-C7 | PASS | As 9 colunas fiscais do 0089 mapeiam 1:1 nos 9 campos que o drawer do FE renderiza; `uf_nome`/`fetched_at` corretamente fora; zero coluna raw/jsonb de `billing_info`. |
| M02-U1 | PASS | Drive do hub pós-merge (tab-2, stack re-apontada): `/catalogo` (10573 produtos, 2940 vendáveis, "dados de 11:22:24", linhas reais com custo/estoque), `/pedidos`, `/anuncios` (34, "dados de 11:22:44"), `/integracoes`. `read_console_messages(onlyErrors)` vazio nas 4. As colunas aditivas do 0089 em `orders_marketplace_orders` não quebraram nenhum read existente. |
| M02-U2 | PASS | Medição do relógio, não de config: backend subiu `2026-08-01T14:16:52Z` (`docker inspect -f '{{.State.StartedAt}}'`). `sync_state.products.last_full_sync_at` estava `13:39:32.324475+00` e avançou para **`14:32:23.975519+00`** — 15min31s após o start, batendo com a cadência de 15 min do `root.go:675`. Polling de 1 em 1 min registrou 3 leituras estáveis antes da virada, então a mudança é atribuível ao tick e não a uma janela ambígua. `last_incremental_at` permanece nulo (products é passe total — correto). |

## Adjudicações do hub

1. **HOLD-MERGE (sanado antes do merge) — fase `sweep` não ratificada.** O chip entregou
   `inferIncremental` com `case "incremental", "sweep": return true`. `"sweep"` não é
   ratificado em lugar nenhum: o ADR-07 (`mission.md:183-187`) define apenas
   `backfill → incremental`, e o próprio doc comment da função prometia "fase não
   reconhecida resolve para false" — contradito pelo próprio código. Risco vivo, e ele é
   CRUZADO: o M-09 já mergeado calcula
   `last_success_at = GREATEST(last_full_sync_at, last_incremental_at)`, e
   `RecordSuccess(..., incremental=true)` só avança `last_incremental_at`. Um passe de
   reparo futuro chamado `"sweep"` congelaria `last_full_sync_at` para sempre e o card de
   saúde mentiria — causa aqui, sintoma em outro milestone meses depois. Corrigido para
   `case "incremental":` com o comentário reescrito nomeando o modo de falha. Fix medido
   pelo hub: 2 arquivos, 9 inserções/6 deleções, nada além. Tip pós-fix `7d02bc75`.
2. **Duas emendas de teste auto-declaradas pelo chip, ambas FORTALECENDO cobertura**
   (F-04 alias bypass no `scanWiringFile` → `resolveCapabilityAliases` + fixture nova;
   F-02 `TestResolveListingFeesIgnoresLayer3` semeava a camada 3 sob `subject_type` que a
   escada nunca consulta — teste vazio, passaria com o filtro deletado). Ambas vieram com
   prova vermelho→verde.

## Ressalvas registradas

- **`status=passed` da lane hermética NÃO prova que os testes de integração executaram.**
  A corrida não emitiu nenhum `failure_token=test=`, e pulado-vs-verde é byte-idêntico
  nessa lane (dívida D-10 de `HARNESS-DEBTS.md`). O que a lane discharge sozinha é
  `migrations_first/second` — o resto de C2/C3 se apoia nas corridas do chip contra
  Postgres real, reportadas com códigos SQLSTATE concretos.
- **tsc do FE vermelho** — inventário pré-existente, M-02 tocou zero arquivo de FE.
- `AssertTerminalCursor` está exportado mas ainda não fiado no `runJob`: diferimento
  intencional e correto de escopo — nenhum critério do M-02 assume enforcement vivo;
  a fiação é do M-04/M-06, quando existirem corpos de job reais.
- **Dívida de harness aberta a partir deste milestone**: D-12 (worktrees
  pré-provisionados stale — 3 de 4 bifurcados 20-480+ commits atrás da base).

## Handoff

- Current status: **passed / merged**.
- Merge: `97758970` (`--no-ff`).
- Next: worktree `gallant-banach-2f909b` a remover após o fechamento da onda A.
- Consumidores de `channelfees`/`divergences` são M-05/M-06/M-07 — hoje zero import no
  repo, confirmado por grep, e isso é o correto.
