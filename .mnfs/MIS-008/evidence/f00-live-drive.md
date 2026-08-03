# F-00 — Task 9 live drive: BLOCKED na precondição de binário

Status: **BLOQUEADO no Step 2** (precondição inegociável do próprio plano). Nenhum dos
Steps 3-8 foi executado porque rodá-los contra o binário errado produziria evidência
inválida — exatamente o risco que o plano nomeia: "ausência de observável lê como
ausência de defeito".

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
