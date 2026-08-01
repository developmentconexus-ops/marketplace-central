# Milestone Validation Result — M-04-listings-backfill-ingest

```yaml
id: M-04-VR
type: milestone-validation-result
status: passed
owner: hub
parent: MIS-007
created: 2026-08-01
validation_level: QA-0
base_sha: 21ca3595
chip_tip: c5e15be47594dbe571430b67767194dd2e03530d
merge_sha: 01f528d2
verdict: PASS (5/6 medidos + U1/U2/U3) · M04-C5 metade viva DEFERRED (tick de 24h)
```

## HOLD-MERGE do hub — como foi quitado

Duas objeções bloqueantes do hub, ambas verificadas por MEDIÇÃO PRÓPRIA (não pelo
auto-relato do chip nem pelos agentes de fix):

**H-1 — amplificação N× no resync (FIXED).** `ResyncWriter.Apply` ignorava
`item.ListingID`/`item.ItemID` e disparava `BackfillPuller.Pull` (catálogo inteiro);
`mutations/application/poller.go:89` chama `Apply` uma vez POR ITEM ⇒ lote de N itens =
N backfills completos. Corrigido em `799b2e89` (delega a `listingsapp.ListingIngestor`,
hidrata+persiste UM anúncio) e endurecido em `c5e15be4` (a instalação embutida no próprio
`ListingID` canônico também é conferida, não só `item.InstallationID`).

Medido aqui: `BackfillPuller` tem **zero referências** no branch inteiro; os 7 testes de
`resync_writer_test.go` rodaram nomeados e passaram, incluindo
`TestResyncWriterBatchCallsIngestorOncePerDistinctListing` (4 ListingIDs distintos → 4
chamadas ao ingestor, zero pulls de catálogo) e
`TestResyncWriterRejectsListingIDEmbeddingForeignInstallation`.

**H-2 — medição do hub estava velha (nada a corrigir).** O teste de integração real de
`Scheduler.RunOnce` já tinha aterrissado em `a09923ac`, depois do meu grep. **Rodei eu
mesmo**: banco efêmero `mpc_test_<32hex>` criado e migrado do zero (**79 migrations
aplicadas**), `go test -tags integration -run TestRealListingsJob` →
`--- PASS: TestRealListingsJobThroughSchedulerPersistsTerminalSweepCursor (0.30s)`, **não
SKIP**. Asserts reais no arquivo: `cursor.Phase == "sweep"`, `LastFullSyncAt == at`, e
`LastIncrementalAt == nil` — "uma fase sweep nunca avança o timestamp incremental",
exatamente a emenda ADR-07 que o hub ratificou na onda A.

**H-3 — ratificado, sem código.** Migration 0090 derruba o CHECK de `status` (correto por
IC-07/ADR-06: vocabulário é do ML, status é verbatim). A defesa contra um status inventado
saiu do banco e passou a ser SÓ o mapper. Registrado no pack.

## Merge e integração

`git merge` automático, **sem conflito**, contra a main já com M-03 (`c5da0d2d`) —
as três regiões do `root.go` do M-04 (import ~66, `installationResyncWriter` ~232-254,
composição ~727-800) não tocam a região de orders 578-621 do M-03.

Integrado (`01f528d2`): build OK, vet OK, `go test -count=1 ./...` = **118 ok / 0 FAIL**.
Migrations 0090/0091/0092 aplicadas no dev stack (`applied 3 migration(s)`);
`listing_variations` existe.

## Critérios

### M04-C1 — Backfill completo + retomada sem duplicata — **PASS (metade viva) / chip (retomada)**

Backfill vivo na conta real, duas corridas: `POST /listings/refresh` → **202** com
`operation_run_id` em 0.96s (async de verdade), run gravado como
`listings_refresh / succeeded / LISTINGS_REFRESH_SUCCEEDED`. As **34 linhas** de `listings`
tiveram `updated_at`/`fetched_at` avançados nas duas corridas (20:46:15 → 20:48:56) e a
contagem **continuou 34** — reentrância sem duplicata na prática.

O braço "kill no meio → retoma sem duplicata" **não foi dirigido por mim** (exigiria matar
o processo no meio de uma corrida de 34 itens); está coberto pelos testes de cursor do chip.

### M04-C2 — Abort pós-página-1 → ZERO rows flipped closed — **PASS (chip) + coerente ao vivo**

Provado pelo chip com must-fail nomeado. Ao vivo, depois de dois sweeps completos:
`absent_since` não-nulo em **0/34** linhas e **zero** linhas em status `closed`. MASS-CLOSURE
morta na prática.

### M04-C3 — E3 + variations no grão certo — **PASS (com lacuna de dataset)**

Preenchimento pós-sweep, 34 linhas: `available_quantity` 34/34 · `sold_quantity` 34/34 ·
`category_id` 34/34 · `permalink` 34/34 · `logistic_type` 34/34 · `last_seen_at` 34/34.

`listing_variations` = **0 linhas**, e isso é honesto, não falha: o payload cru dos 34
anúncios traz `variations` vazio em todos (`sum(jsonb_array_length(raw->'variations')) = 0`)
e `variation_id` é o sentinela `'-'` em 34/34. **A conta não tem nenhum anúncio com
variação**, então o grão de variação não é observável ao vivo aqui — só pelos fixtures do
chip.

### M04-C4 — Âncoras de snapshots não-regressivas (ADR-13) — **PASS (chip)**

Regressão de SellerSKU/EAN foi achada pelo gate frio, corrigida e re-medida pelo chip
(`c251dfae`, `8c9a27eb`). Ao vivo: `/anuncios` mostra `Sem vínculo 4` de 34 — o observador
de snapshots segue alimentando o matcher.

### M04-C5 — Scheduler diário + refresh manual no mesmo caminho — **PASS parcial · metade viva DEFERRED**

**Passa:** refresh manual e scheduler compartilham hidratador e store por construção
(`root.go:757-763`: `listingBackfillHydrator` + `listingRepo` são passados TANTO ao
`BackfillRunner` do refresh QUANTO ao `installationResyncWriter`), e
`NewListingsSchedulers(pool, 24*time.Hour, …).StartAll(ctx)` está registrado no boot
(`root.go:775-777`).

**DEFERRED, com a razão medida:** `sync_state` **não tem linha para `listings`** e não pode
ter hoje. `syncapp.Scheduler.Start` (`scheduler.go:105-118`) é `time.NewTicker(interval)`
puro — **sem tick inicial** — logo a primeira `RunOnce` de listings só acontece 24h após o
boot. O mecanismo está provado pelo teste de integração que EU rodei (ver H-2); o que falta
é exclusivamente o relógio de parede, que é indispensável de provar hoje.

**Achado de desenho que isso expõe (novo, registrado como D-16):** ticker de 24h + zero tick
inicial + zero persistência de "próximo vencimento" ⇒ se o processo reiniciar em intervalo
menor que 24h, a varredura de listings **nunca roda**. O card de saúde do M-09 vai mostrar
`listings` como nunca sincronizado indefinidamente, e estará dizendo a verdade.

### M04-C6 — Paginação provada com fixture >1 página — **PASS (chip)**

Fixture multi-página do chip. Ao vivo não é discriminante: a conta tem 34 anúncios, abaixo
de uma página.

## Critérios de user-drive

### M04-U1 — /anuncios com o catálogo real completo — **PASS**

Drive em `http://localhost:5174/anuncios`. Tela: `Total 34 · Ativos 0 · Pausados 27 ·
Com erro de sync 0 · Desatualizados 0 · Sem vínculo 4`. DB: `select count(*) from listings`
= **34**; `group by status` = `paused 27`, `under_review 7`. Bate exato em total e pausados.
**34 IDs `MLB…` distintos** renderizados na tabela (contados no DOM) — catálogo inteiro na
tela, sem truncamento de página.

Observação sem veredito: os 7 `under_review` não caem em `Ativos` nem em `Pausados` (a soma
das duas abas dá 27, não 34). Coerente com o total, mas o vocabulário de aba não cobre o
status.

### M04-U2 — Refresh manual dirigido na tela — **PASS**

Clique no botão `Atualizar` da tela → 4ª corrida `listings_refresh` gravada,
`max(fetched_at)` avançou de `20:46:15` para `20:48:56`, contagem seguiu 34, distribuição de
status intacta, e o carimbo da tela avançou de `17:48:30` para `17:49:15` **sem F5**. Zero
erro na tela (a única string com "erro" é o rótulo do KPI `Com erro de sync`, valor 0) e
zero erro de console.

### M04-U3 — Status pausado na tela + absent≠closed — **PASS (metade discriminante ausente)**

`Pausados 27` na tela = 27 linhas `paused` no banco, lifecycle de row correto. Depois de
dois sweeps completos, `absent_since` = 0/34 e nenhuma linha virou `closed`.

**Declarado:** o caso que DISCRIMINA absent≠closed (item some do scan e a row NÃO flipa até
o sweep confirmar) exigiria remover um anúncio da conta real do operador — não foi feito.
Ao vivo há só a metade não-regressiva; a discriminante é a do must-fail do chip (C2).

### Rede/console

`performance.getEntriesByType('resource')` filtrado por `/mercadolibre|mercadolivre/i` = **0**
em `/anuncios`. Console: **zero erros**.

## Handoff

- Status: **passed**, merge `01f528d2` na main.
- Dívida aberta por este milestone: **D-16** (starvation do scheduler de 24h sem tick
  inicial) em `.mnfs/HARNESS-DEBTS.md`.
- Não-medido ao vivo, com razão: retomada pós-kill (C1), grão de variação (C3 — conta sem
  anúncio com variação), paginação >1 página (C6 — conta com 34 itens), primeiro tick do
  scheduler (C5 — 24h de relógio), absent→closed discriminante (U3 — exigiria mexer na conta
  do operador).
