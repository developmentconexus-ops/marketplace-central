# F-00 — Task 9 live drive: `/integracoes`

Status: **DESBLOQUEADO via hub** (`local_002d928a-48e6-4f41-b5e8-86bc00e01cb6`), coordenado
por REQUEST em fila com a sessão FA3. Ver "Continuação pós-desbloqueio" no fim deste arquivo
para o andamento atual. A seção abaixo é o registro original do bloqueio (preservado como
histórico — o diagnóstico dela é o que motivou o REQUEST ao hub).

## Registro original (BLOQUEADO no Step 2, antes da coordenação com o hub)

## Step 1 — subir a stack

Não precisou subir nada: a stack compartilhada já estava de pé (`docker ps`), usada por
outras sessões concorrentes (`fa1b-reautorizacao`, `fa3-idade-honesta`,
`p2-dinheiro-real-pedidos`, `p2b-imposto-ex-ante` — `git worktree list` na hora da medição).
Por instrução explícita de não derrubar/reconstruir infra compartilhada sem confirmar
segurança, segui para a prova de frescor (Step 2) antes de tocar em qualquer container.

## Step 2 — prova de frescor do binário: FALHOU

### Estado do branch

```
git merge-base HEAD main            -> 6bd22c295796a0adb10aec484474897f06db5ec9
git log --reverse --oneline main..HEAD | head -1   -> 3119f392 plan(F-00): baseline de governanca medido na lane real
git log -1 --oneline HEAD           -> c6df376d feat(orders): scheduler de pedidos ligado no root a cada 15 minutos
git log -1 --format="%H %ci" HEAD   -> c6df376d1f66ad0f3bf6b8e8c0e34ce93a6563ca 2026-08-03 16:57:54 -0300
```

`main` HEAD no momento da medição: `62b5d5d4` ("docs(harness): register B-10b ... B-11").

```
git merge-base --is-ancestor c6df376d main   -> NÃO (c6df376d NÃO é ancestral de main)
```

### Estado do container

```
docker ps --format "table {{.Names}}\t{{.Image}}\t{{.CreatedAt}}\t{{.Status}}"
```

```
NAMES                            IMAGE                          CREATED AT                      STATUS
marketplace-central-backend-1    marketplace-central-backend    2026-08-03 16:51:45 -0300 -03   Up 7 minutes (healthy)
marketplace-central-frontend-1   marketplace-central-frontend   2026-08-03 14:56:08 -0300 -03   Up 2 hours
marketplace-central-postgres-1   postgres:16-bookworm           2026-07-20 19:46:22 -0300 -03   Up 3 days (healthy)
```

```
docker inspect marketplace-central-backend-1 --format '{{.Created}}'
-> 2026-08-03T19:51:45.8846834Z   (= 16:51:45 -0300, bate com docker ps)

docker images marketplace-central-backend --format "{{.ID}}\t{{.CreatedAt}}"
-> 2432c809e3dc   2026-08-03 16:49:11 -0300 -03
```

### A causa raiz não é idade — é o checkout errado

```
docker inspect marketplace-central-backend-1 --format '{{range .Mounts}}{{.Source}} -> {{.Destination}}{{"\n"}}{{end}}'
-> /run/desktop/mnt/host/c/Users/leandro.theodoro/Documents/marketplace-central -> /workspace

docker inspect marketplace-central-backend-1 --format '{{index .Config.Labels "com.docker.compose.project.working_dir"}}'
-> C:\Users\leandro.theodoro\Documents\marketplace-central

docker inspect marketplace-central-backend-1 --format '{{index .Config.Labels "com.docker.compose.project.config_files"}}'
-> C:\Users\leandro.theodoro\Documents\marketplace-central\docker-compose.yml
```

`docker/dev/backend-entrypoint.sh` roda `exec go run ./apps/server_core/cmd/server`
direto contra o `/workspace` montado — ou seja, o processo backend hoje executa o código
do checkout **principal** (`main`, tip `62b5d5d4`), não o código deste worktree
(`.claude/worktrees/f00-scheduler-pedidos`, tip `c6df376d`). Como `c6df376d` não é
ancestral de `main`, **nenhum rebuild de imagem resolveria isso sozinho** — o volume bind
serve o diretório errado independente da idade da imagem. A idade do container (7 min) é
sintoma, não a causa.

## Por que não segui em frente (adjudicação de escopo)

Duas rotas óbvias para corrigir isso foram descartadas, ambas por razão documentada, não
por escolha arbitrária:

1. **Rodar uma stack isolada a partir deste worktree** (`docker compose -p <outro-projeto>
   up --build`, portas alternativas). Inviável sem um `.env` neste worktree — e o backend
   lê credenciais reais de app ML do ambiente:
   `MPC_PROVIDER_MERCADOLIVRE_CLIENT_ID` / `MPC_PROVIDER_MERCADOLIVRE_CLIENT_SECRET`
   (`apps/server_core/internal/modules/integrations/adapters/mercadolivre/auth_adapter.go:70-71`).
   `docs/HARNESS-PROFILE.md` §6 ("Browser-drive seam") nomeia exatamente essa tentação e a
   proíbe: *"creating a `.env` in the chip's worktree — is forbidden: it would place live
   provider credentials inside a chip session's reachable tree."*

2. **Fazer o "borrow" do hub eu mesmo** (`git checkout <branch> -- <arquivos>` dentro do
   checkout principal, dirigir o browser lá, reverter) — é o mecanismo que o próprio §6
   descreve como o caminho correto ("the hub does `git checkout <chip-branch> -- <exact
   files>` into its OWN checkout"). Mas essa tarefa foi despachada com a instrução explícita
   de trabalhar **somente** dentro deste worktree e não tocar nenhum outro
   caminho/worktree/branch — o que essa rota exigiria fazer no checkout principal
   (`C:\Users\leandro.theodoro\Documents\marketplace-central`).

`docs/HARNESS-PROFILE.md` §6 também nomeia a stack dev (`:8080`/`:5174`) como **hub-owned**
("chips never touch containers") e já tem o mecanismo formal para isso — o evento
`stack-sync: <sha>`, que dispara o hub a reconstruir/redeploy a stack naquele SHA sem
negociação por pedido.

Ou seja: as duas correções possíveis pertencem ao hub, não a este worker, e uma delas
recria explicitamente um padrão proibido pela doutrina se eu tentasse fazer sozinho dentro
deste worktree.

## O que É verdade, medido, sem depender do binário

Isto não depende de qual checkout está montado — é git puro:

- Este branch (`f00-scheduler-pedidos`) tem `apps/web` intocado desde que divergiu de
  `main`: `git diff --name-only main...HEAD -- apps/web` → saída vazia.
- `git status --short` no worktree → limpo (nenhuma mudança pendente).
- O código da Task 8 (fiação do scheduler em `internal/composition/root.go`) está commitado
  em `c6df376d`, no branch, pronto para ser servido assim que o binário certo estiver de pé.

## Ação necessária para desbloquear (não executada por este worker)

O hub (ou uma sessão com autoridade sobre o checkout principal / stack compartilhada)
precisa fazer UMA das duas:

- **(A)** Emitir/agir sobre `stack-sync: c6df376d` — rebuild + redeploy do serviço
  `backend` apontado para o tip deste branch (`c6df376d`), preservando o Postgres já em
  pé (mesmos dados, mesma instalação ML conectada) — depois os Steps 3-8 deste plano podem
  rodar contra o binário certo.
- **(B)** Fazer o borrow-and-revert descrito em `docs/HARNESS-PROFILE.md` §6 no seu próprio
  checkout (`git checkout f00-scheduler-pedidos -- apps/server_core/internal/composition/root.go
  <demais arquivos do grant>`), dirigir `/integracoes` e a query SQL lá, devolver o texto
  capturado verbatim, reverter.

Nenhuma escrita, rebuild, restart ou tentativa de re-apontar a stack foi feita por este
worker. Nenhum arquivo fora deste worktree foi tocado.

## Continuação pós-desbloqueio

O hub aceitou o REQUEST (rota A — stack-sync), com uma pré-condição própria: mergear
`afb6b54a` (main, que trouxe a migração 0093 renomeando `entity='market'` → `'market_queue'`
para o enqueuer do erp_import) antes do repoint, porque o código desta branch ainda escrevia
o sentinela velho (`EntityMarket`) e reintroduziria a linha errada no Postgres compartilhado.
Merge feito: `git merge --no-edit afb6b54a`, sem conflito real (`root.go` e
`HARNESS-DEBTS.md` auto-mergearam — vizinhança de linha, não sobreposição semântica). Tip
mergeado: `83c3ee40` (contém `afb6b54a` — superconjunto, nada da FA3 se perde). Verificado
antes de avisar o hub: `go build ./...` limpo; `go test` verde em
`orders/... sync/... erp_import/... composition/...`; governança contra
`-BaseSha afb6b54a` = 14 `GOV_MODULE_DEPENDENCY` + 8 `GOV_MODULE_LAYER` (o `+1` de
dependency é `market/application/collection_job.go`, da FA3, não desta fatia; layer
inalterado 8=8).

Fila: FA3 tinha o slot primeiro (ciclo do market job, 30min, confirmado
`2026/08/03 20:48:31 INFO market collection job: ciclo concluído colhidos=1 falhas=0
teto=50`). Hub aguardando o CLOSED dela (ou timeout) antes de repontar o backend para
`83c3ee40`, já armado e testado a seco (override explícito fora da árvore, bind único de
`/workspace` para este worktree, Postgres e frontend intocados).

### Step 3 — Estado antes (capturado com o backend ainda em `afb6b54a`, pré-repoint)

```sql
SELECT entity, installation_id, last_full_sync_at, last_incremental_at, consecutive_failures, last_error, cursor
FROM sync_state ORDER BY entity, installation_id;
```

```
    entity    |                     installation_id                     |       last_full_sync_at       | last_incremental_at | consecutive_failures | last_error |  cursor
--------------+---------------------------------------------------------+-------------------------------+---------------------+----------------------+------------+-----------
 listings     | inst-mercado_livre-7e0d2125-f525-4174-9ade-8c7dc496a0e0 | 2026-08-03 13:24:05.282956+00 |                     |                    0 |            | {"phase": "backfill", "ids_collected": 34, ...}
 market       | market                                                  | 2026-08-03 20:48:31.961669+00 |                     |                    0 |            |
 market_queue | inst-mercado_livre-7e0d2125-f525-4174-9ade-8c7dc496a0e0 |                               |                     |                    0 |            | {"pending": [...]}
 products     | erp                                                     | 2026-08-03 20:48:51.343395+00 |                     |                    0 |            | {"source": "sankhya", "processed": 10577, ...}
(4 rows)
```

**Confirmado: nenhuma linha `orders`.** Controle negativo mais forte disponível — nenhuma
outra fatia ativa escreve essa entidade, então qualquer linha `orders` que aparecer depois
do primeiro ciclo do scheduler desta fatia é atribuível a este código e a nada mais.

### Interrupção de janela (F-A1b sobrescreveu o backend sem avisar)

Cronologia medida pelo hub e confirmada por mim (`docker inspect --format '{{.Created}}'`):

```
20:53:29Z  backend recriado no meu worktree (janela abre)
21:04:35Z  backend recriado montando a MAIN — sessão F-A1b precisava da main de pé
           para a própria live drive de reautorização, viu o stack servindo minha
           árvore e reconstruiu por cima em vez de pedir a janela
21:28:26Z  hub detecta a troca; sync_state ainda sem linha orders
21:28:40Z  backend recriado no meu worktree de novo (janela reaberta)
```

Minha janela durou 11 minutos contra um ticker de 15 — o primeiro ciclo do scheduler de
pedidos nunca chegou a disparar nela. **Descartado**: qualquer leitura de `sync_state`
feita entre 21:04:35Z e 21:28:40Z não provaria nada sobre este código (binário era o da
main). Nenhuma foi usada como evidência de sucesso ou falha acima — a única coisa medida
nesse intervalo foi a AUSÊNCIA da linha `orders`, que é verdadeira independente de qual
binário rodava (nenhum dos dois escreve essa entidade sem o scheduler certo de pé), então
não contamina o controle negativo, só adia a prova positiva.

Verificado por mim, não só relatado pelo hub:

```
docker inspect marketplace-central-backend-1 --format '{{.Created}}'
-> 2026-08-03T21:28:40.706782975Z   (agora = 21:29:51Z)

docker exec marketplace-central-backend-1 sh -c 'cat /workspace/.git'
-> gitdir: C:/Users/leandro.theodoro/Documents/marketplace-central/.git/worktrees/f00-scheduler-pedidos
```

**Novo baseline do ticker: 21:28:40Z. Primeiro ciclo esperado ~21:43:40Z.** `sync_state`
seguia com as mesmas 4 linhas (listings/market/market_queue/products), sem `orders`, no
momento desta reabertura — o dump do hub bate com o meu do Step 3 original.

Achado incidental, não fabricado por mim: às 21:15:02 o refresh automático de token desta
mesma instalação ML falhou com `invalid_grant` real ("Your authorization code or refresh
token may be expired or it was already used"); a sessão F-A1b então rodou um fluxo de
reautorização ao vivo contra a MESMA instalação (`start_authorize` → `handle_callback`),
que atualizou `integration_installations.updated_at` para 21:19:32Z (status seguiu
`connected`). Uma segunda tentativa de draft às 21:20:01 falhou com
`PROVIDER_ACCOUNT_ALREADY_LINKED` (esperado — a instalação já existia). Registro aqui como
fato observado, não como algo que eu causei ou revertive — nenhuma escrita em provider
partiu deste worker.

### Step 2 (confirmado, segunda vez — container trocou de mão duas vezes)

`docker inspect marketplace-central-backend-1 --format '{{.Created}}'` → `2026-08-03T21:28:40.706782975Z`.
`docker exec ... cat /workspace/.git` → `gitdir: .../worktrees/f00-scheduler-pedidos`.
Processo real (`server starting on :8080`) subiu às `21:28:54Z` — o ticker do scheduler de
pedidos conta a partir daí, não do `Created` do container.

### Step 4/5 — antes/depois (confirmado)

`/integracoes`, antes do primeiro ciclo: card "Saúde do sync" com Listings/Market/Market
Queue/Products, sem linha de pedidos (capturado via `get_page_text`, screenshot indisponível
nesta sessão — ver nota abaixo).

Primeiro ciclo disparou em `run_started_at: 2026-08-03T21:43:54.482207927Z` — exatamente
boot (`21:28:54Z`) + 15min. `sync_state` ganhou a 5ª linha:

```
orders|inst-mercado_livre-7e0d2125-f525-4174-9ade-8c7dc496a0e0
  last_incremental_at=2026-08-03 21:44:30.675159+00
  consecutive_failures=0, last_error=(vazio)
  cursor={"phase":"incremental","offset":0,"run_started_at":"2026-08-03T21:43:54.482207927Z",
          "last_updated_at":"2026-08-03T20:18:45Z","last_run_skipped":0,
          "last_run_imported":39,"last_run_enumerated":39}
```

`orders_marketplace_orders`: `count=39`, `max(fetched_at)=2026-08-03 21:44:29.990968+00`,
`max(provider_updated_at)=2026-08-03 19:16:22+00`. `fetched_at` avançou sem nenhum clique
manual de importar — critério da feature cumprido.

UI (`get_page_text` pós-reload): `Orders(incremental) — há menos de 1 min — ok`, ao lado de
Listings/Market/Market Queue/Products, sem nenhuma mudança em `apps/web` (rótulo genérico do
`entityLabel`, como o plano previa).

**Nota sobre captura visual**: `computer{action:"screenshot"}` retornou erro ("Browser pane
não exibido, não compositando frames") nesta sessão headless. Uso `get_page_text`/DOM como
observável equivalente — verbatim do texto renderizado, suficiente para as asserções deste
passo (presença/ausência de linha, cor não verificável textualmente — ver Step 6 para como
a UI expressa vermelho).

### Interrupção de janela nº2 (F-A1b reconstruiu o backend por cima da minha árvore)

Ver seção acima ("Interrupção de janela") — resumo: minha janela original (20:53:29Z–
21:04:35Z, 11min) não foi suficiente para o primeiro ciclo (ticker de 15min); F-A1b
reconstruiu o backend montando `main` às 21:04:35Z para a live drive de reautorização dela;
hub reabriu minha janela às 21:28:40Z (confirmei via `docker inspect` + `.git` interno) e
o ciclo correto disparou dentro dela, sem mais interrupções.

### Step 6 — provar o vermelho (mutação reversível de account ref)

Coordenado com o hub antes de mutar (nenhuma outra sessão ativa nesta instalação no
momento). Confirmado no código, antes de escolher o valor, que a mutação realmente atinge a
chamada real ao provider: `provider_operation_service.go:246` —
`ProviderAccountID: inst.ExternalAccountID` — e `capability_adapter.go:489` —
`query.Set("seller", accountRef.ProviderAccountID)` no `/orders/search`. Ou seja,
`integration_installations.external_account_id` alimenta literalmente o parâmetro `seller`
da busca de pedidos no Mercado Livre — sem essa checagem, um valor numérico plausível
voltaria 200 com zero resultados (falso-negativo: `consecutive_failures` ficaria em 0 e a
task concluiria "vermelho não acende" por não ter havido erro nenhum).

**Valor original preservado**: `external_account_id = '691607102'`. SQL de restauração
escrito e salvo ANTES da mutação (não improvisado depois):

```sql
UPDATE integration_installations SET external_account_id = '691607102'
WHERE installation_id = 'inst-mercado_livre-7e0d2125-f525-4174-9ade-8c7dc496a0e0';
```

Mutação aplicada às `2026-08-03 21:49:44Z` (autorizada explicitamente pelo operador em
chat, após o classificador de auto mode bloquear o `UPDATE` na primeira tentativa por ser
escrita em DB compartilhado):

```sql
UPDATE integration_installations SET external_account_id = 'INVALID-F00-RED'
WHERE installation_id = 'inst-mercado_livre-7e0d2125-f525-4174-9ade-8c7dc496a0e0';
```

Nenhuma credencial ou token tocado — só este único campo, exatamente o limite que o plano
e o hub estabeleceram. Próximo ciclo esperado ~21:58:54Z. Aguardando confirmação de
`consecutive_failures > 0` + `last_error` preenchido + linha vermelha na UI antes de
restaurar (continuação abaixo).

### Step 6 — vermelho confirmado

Ciclo seguinte rodou e falhou de verdade contra o provider real (não simulado):

```
sync_state.orders: consecutive_failures = 1
last_error = {"at": "2026-08-03T21:58:55.475404763Z",
              "message": "CONNECTORS_PROVIDER_AUTH: provider credentials were rejected"}
```

UI (`/integracoes`, DOM via `get_page_text`/`read_page`, screenshot indisponível na sessão
inteira — ver nota de substituição acima): o card de saúde do sync trocou de rótulo/cor —
`Orders(incremental) — 1 falha`, classe `bg-warn-soft`/`text-warn-ink` (âmbar,
`rgb(208,138,99)` sobre `rgb(58,42,32)`, confirmado via `javascript_tool` `getComputedStyle`)
contra o estado saudável anterior `bg-accent-soft`/`text-accent-ink` (verde,
`rgb(158,207,171)` sobre `rgb(36,51,40)`). Falha real de provider vira sinal visível na
tela — é exatamente o que Step 6 pede provar.

**Restauração**: aplicado o SQL pré-escrito às `2026-08-03 22:00:10Z`. Confirmado por leitura
direta: `external_account_id = '691607102'` (valor original, não um literal parecido).

**Verde de volta (ciclo seguinte, natural, sem intervenção manual)**:

```
sync_state.orders: last_incremental_at = 2026-08-03 22:13:56.350202+00
                    consecutive_failures = 0
                    last_error = (vazio)
```

Container do backend não trocou nesse intervalo (`Created` seguiu
`2026-08-03T21:28:40.706782975Z`, mesmo valor de antes) — o ciclo verde é do MEU binário,
não artefato de rebuild de outra sessão. UI confirma: `Orders(incremental) — há menos de 1
min — ok`, de volta ao rótulo/cor saudável. Loop fechado: vermelho que aparece E limpa,
não estado preso.

### Step 7 — orders avança sozinho (comparação com baseline)

Baseline do plano (§Medição.5): 38 linhas, `fetched_at` 2026-08-02, `provider_updated_at`
2026-07-31.

Observado após o live drive (scheduler rodando sozinho, nenhum import manual):

```
count | max(fetched_at)               | max(provider_updated_at)
   39 | 2026-08-03 21:44:29.990968+00 | 2026-08-03 19:16:22+00
```

`fetched_at` avançou de 2026-08-02 para 2026-08-03 21:44:29Z e a contagem foi de 38→39 —
o scheduler trouxe pedido novo e atualizou o carimbo de tempo sem qualquer clique de
importação manual. Step 7 fechado.

### Pendente

- [x] Step 6 (continuação): vermelho confirmado, restaurado, verde confirmado de volta
- [ ] Step 8: commit final desta evidência
