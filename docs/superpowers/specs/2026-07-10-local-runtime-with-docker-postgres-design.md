# Runtime local com PostgreSQL Docker — design

## Objetivo

Operar Marketplace Central sem Docker para backend e frontend, usando Docker
somente para o PostgreSQL local. O fluxo deve ser reproduzível no Windows,
seguro para o worktree sujo e explícito sobre pré-requisitos, atualização e
diagnóstico.

## Escopo

Criar:

- um orquestrador PowerShell em `scripts/dev-local.ps1`;
- um runbook em `wiki/operations/` e link no índice da wiki;
- uma skill pessoal `local-mpc-runtime` em `$CODEX_HOME/skills`;
- comandos de runtime, status, build, teste e parada local.

Não criar ou remover volumes, não limpar dados, não subir backend/frontend no
Docker, não modificar credenciais, não imprimir valores de `.env`, e não
alterar o worktree além dos arquivos deste recurso.

## Arquitetura

```text
Docker Compose: postgres somente (host port 5435)
                 │
                 ▼
MC_DATABASE_URL ──► backend Go local (:8080) ──► frontend Vite local (:5174)
```

O script executa no host Windows. Ele carrega `.env` apenas no ambiente dos
processos filhos, remove aspas externas de valores e nunca apresenta valores
de configuração. Para o frontend, define `MPC_WEB_PROXY_TARGET` como
`http://localhost:8080`; isso substitui o padrão Docker `http://backend:8080`.

## Interface do script

`scripts/dev-local.ps1` terá os comandos:

| Comando | Comportamento |
| --- | --- |
| `up` | Validar pré-requisitos, iniciar/aguardar apenas `postgres`, aplicar migrations, iniciar backend e frontend locais e registrar PIDs/logs no diretório temporário do projeto. |
| `status` | Mostrar estado do PostgreSQL, health HTTP, portas e processos registrados, sem mostrar segredos. |
| `build` | Rodar build Go e build Vite localmente, com `GOCACHE=.gocache`. |
| `test` | Rodar testes locais impactados; testes de PostgreSQL usam `MC_DATABASE_URL` somente se configurado. |
| `stop` | Parar somente processos backend/frontend lançados por `up`; não parar PostgreSQL e não remover dados ou containers. |

`up` não tenta reparar Docker. Se o PostgreSQL estiver indisponível, o script
falha com estado e próximo passo objetivo; não reinicia Docker Desktop, não
executa `compose down`, não remove volumes e não tenta recuperar banco.

## Pré-requisitos e falhas

O script valida antes de iniciar:

- PowerShell, `go`, `node` e `npm` disponíveis;
- `.env` presente e `MC_DATABASE_URL` definido;
- Docker CLI acessível somente para `docker compose up -d postgres` e health
  do PostgreSQL;
- ferramentas Oracle/CGO apenas quando a configuração Oracle estiver presente
  ou quando o operador pedir validação Oracle ao vivo.

Uma falha de pré-requisito encerra antes de migrations ou processos filhos.
Falhas de migration, backend ou Vite preservam os logs e os PIDs para o
comando `status`; `stop` continua capaz de encerrar os processos lançados.

## Documentação e skill

O runbook documenta instalação, rotina diária, atualização, build/teste,
Oracle live validation, logs e limites do Docker. A skill pessoal orienta
futuras sessões a usar o script, preservar dados, distinguir health local de
evidência de integração e não trocar resultados de mock por validação real.

## Verificação

- validação sintática do script;
- `status` contra ambiente indisponível e disponível, sem expor `.env`;
- `build` local de backend e frontend;
- `up` com PostgreSQL saudável, seguido de `/healthz` e Vite em `:5174`;
- `stop` confirma que somente processos registrados foram encerrados;
- revisão do runbook e validação estrutural da skill.

## Decisões explícitas

- PostgreSQL permanece Dockerizado por escolha do operador.
- Backend/frontend rodam somente no host Windows nesse fluxo.
- Docker Desktop é uma dependência externa para o banco; sua recuperação não
  pertence ao script.
- O procedimento não aprova links de produto nem altera dados de integração.
