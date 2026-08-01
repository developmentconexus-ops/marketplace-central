# Milestone Validation Result — M-09-sync-observability

```yaml
id: M-09-VR
type: milestone-validation-result
status: passed
owner: hub
parent: MIS-007
created: 2026-08-01
validation_level: QA-0
base_sha: dd89d4b3
merge_commit: 7f5fdbe2
chip_tip: 64e25974
correction_attempts: 2
```

## Verdict

**PASS.** 5/5 critérios de milestone não-diferidos + 3/3 de user-drive.
M09-C6 permanece **DIFERIDO** por desenho do próprio contrato (re-drive após M-04/M-06;
bloqueia o close da MISSÃO via MIS07-C8, não o deste milestone).
**M09-U3 reprovou DUAS vezes no live-drive do hub** antes de passar — ver Adjudicações.

## Evidência por critério

| ID | Veredito | Observável medido |
|----|----------|-------------------|
| M09-C1 | PASS | Teste de evidência com Postgres real imprime lado a lado o `SELECT * FROM sync_state` e o JSON de `GET /sync/health`. A row `installation_id=erp` (sentinela do ERP) APARECE no payload — o scan é all-rows, não filtrado por installation ML. `listings` nunca-rodada sai com `last_success_at: null`, `phase: null` (nenhum timestamp fabricado); `products` com `phase:"incremental"` extraído do cursor jsonb. `TestHealthReaderEmptyTenant` prova `entities: []` (lista vazia honesta, não `null`). |
| M09-C2 | PASS | `TestHealthServiceLastSuccessAtIsGreatestNotFullSyncAlone` com a fixture exata do F-r04-1 (full VELHO `2025-12-31`, incremental RECENTE `2026-07-31`): o payload devolve `last_success_at = 2026-07-31T11:00:00Z` — o incremental. O defeito que o F-r04-1 nomeia (entidade saudável mostrada como velha) não reproduz. Badge verde decidido por ESTADO (`last_success_at` não-nulo ∧ `consecutive_failures==0`), nunca por corte temporal. |
| M09-C3 | PASS | `TestHealthHandlerWithWebhookStatsReaderAfterRegisterIsObservedLive` — o fake é injetado DEPOIS do `Register(mux)` e seus valores aparecem no `GET /sync/health` servido pela rota montada: prova de injeção por referência (F-r04-2), não por cópia. `TestHealthServiceDefaultWebhookIsCanonicalZero` fixa o default `{null,0,0}` do IC-05, comparado por estrutura (não string-eq — lição de round-trip JSONB). |
| M09-C4 | PASS | Estados verde/vermelho/cinza-"nunca"/webhook-inicial cobertos em teste de componente **e** confirmados no drive do hub (ver M09-U1/U2/U3). O rótulo de veredito "não configurado" está PROIBIDO e não existe no código — a cópia é o fato "Nenhuma notificação recebida.". |
| M09-C5 | PASS | `/sync/health` no OpenAPI + `getSyncHealth` no `sdk-runtime` no MESMO commit. Diff da `IntegracoesPage.tsx` = exatamente **2 linhas** (import + `<SyncHealthCard />` após `<ProviderConnectCard />`); o card é arquivo novo. `IntegracoesPage.test.tsx` existente verde sem edição. |
| M09-C6 | **DIFERIDO** | Re-drive após M-04/M-06 fecharem. Hoje `listings`/`orders` ainda não têm row em `sync_state`, então o critério é inobservável por construção. Não bloqueia este close (o próprio contrato diz isso); rastreado por MIS07-C8. |
| M09-U1 | PASS | Drive na `/integracoes` com backend de pé: seção mostra `Products / há 49 min / ok` e `Market / nunca`. Confere com o `SELECT`: `products.last_full_sync_at = 2026-08-01 13:39:32.324475+00` contra `now() = 14:28` → 49 min. O absoluto está no `title` (hover). Payload vivo salvo em `_chip-m09/live-health-post-merge.json`. |
| M09-U2 | PASS | Mesmo drive, dois estados honestos simultâneos: `market` = "nunca" em cinza (nulls, NUNCA "0 min atrás") e webhook = "Nenhuma notificação recebida." — FATO observado, não veredito de configuração. |
| M09-U3 | PASS (na 3ª rodada) | Backend PARADO (`docker stop marketplace-central-backend-1`), aba NOVA (`tab-3`), zero injeção de JS. DOM medido da seção: `<div data-testid="sync-health-unknown">…<p role="alert">Erro ao carregar. Não foi possível carregar o status da sincronização. Estado desconhecido.</p><button …>Tentar novamente</button>`. Resto da página intacto e operável na MESMA carga: "Fonte ativa", "Sortimento vendável", "Importar catálogo", "Conectar marketplace", "Importação" — todos renderizados. |

## Adjudicações do hub

**M09-U3 reaberto duas vezes.** O contrato pedia "health 500 simulado → card mostra erro
nomeado". Com o backend inalcançável, o card renderizava **só o `<h2>`**:

```html
<section aria-labelledby="sync-health-title" class="rounded-card border border-border bg-surface p-4"><h2 id="sync-health-title" class="text-sm font-semibold text-ink">Saúde do sync</h2></section>
```

- **Rodada 1** — o chip tinha teste VERDE do ramo `isError` (estado injetado no mock). O
  hub reproduziu o card mudo ao vivo. Diagnóstico do chip: `networkMode`. Fix aplicado em
  `packages/web-query/src/syncHealth.ts`.
- **Rodada 2** — o hub verificou que o fix estava genuinamente NO FIO (o container do FE
  monta o repo em `/workspace` e servia o módulo com `networkMode`), e mesmo assim
  reproduziu o MESMO HTML mudo em aba nova. O teste novo, em jsdom contra recusa de TCP,
  estava verde — jsdom + fetch do Node não reproduz o `onlineManager` do navegador.
  Ordem do hub: **pare de perseguir a causa-raiz; o defeito é de DESENHO** — três
  ternários `? … : null` independentes e nenhum ramo total. Exigido vermelho antes do verde.
- **Rodada 3 (aceita)** — cadeia de ternários trocada por `if/else` sobre `query.status`,
  união fechada de 3 membros (`'pending'|'error'|'success'`) do `@tanstack/query-core`:
  descartados `error` e `success`, o TS estreita o resto para `pending` — o `else` final
  não é chute sobre sobra, é o único valor que o tipo permite. Dentro de `pending`,
  `fetchStatus==='fetching'` → loading; `paused`/`idle` sem dado → `ErrorState` com texto
  honesto de estado desconhecido (ADR-17), deliberadamente distinto da cópia de fato do
  webhook. `SyncHealthCard.stateMatrix.test.tsx` enumera as formas reais de
  `UseQueryResult` incluindo as DUAS indeterminadas; contra o código anterior 2 dos 5
  casos ficaram vermelhos — exatamente os indeterminados. Verificado pelo hub na main
  integrada: 5/5 verde.

`networkMode:"always"` foi MANTIDO como hardening, não como a alegação de causa.

## Ressalvas registradas

- **Classe nomeada e registrada como dívida D-11**: *teste de estado injetado não prova
  transição de estado*. Reincidiu na mesma milestone, e é o motivo de U3 ter reprovado
  duas vezes com o teste verde. Corolário: um observável que só existe no navegador não
  pode ser certificado fora dele.
- **tsc do FE vermelho** — inventário pré-existente (11-12 erros conforme o corte);
  hub-ops confirmou que nenhum arquivo novo (`SyncHealthCard.tsx`, `syncHealth.ts`,
  `sdk-runtime`) aparece na lista.
- Screenshots light+dark em `docs/design/evidence/` não foram produzidos; o observável
  registrado aqui é o **DOM verbatim**, que é mais forte que a captura para o que U3 pede
  (screenshot morto não prova qual ramo renderizou).

## Handoff

- Current status: **passed / merged**.
- Merges: entrega original `4d837303`, fix U3 r1 `6ef6c410`, fix U3 r2 `7f5fdbe2` (todos `--no-ff`).
- Next: worktree `musing-taussig-0b8ae9` a remover após o fechamento da onda A.
- Aberto: **M09-C6** (re-drive pós M-04/M-06) — entidades novas devem acender sem NENHUM
  commit no código do M-09.
