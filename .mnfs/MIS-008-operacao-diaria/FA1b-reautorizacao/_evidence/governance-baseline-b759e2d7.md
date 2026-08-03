# Baseline de governança — pré-merge F-A1b

SHA medido: `b759e2d7` (tip da `main` antes do merge da fatia F-A1b)
SHA de comparação (pós-merge, mesmo método, mesmo resultado): `6147c0d1`

## Método

Motivo: medir a lane de governança dentro do checkout principal contamina o
resultado — cache `.gocache` aquecido por rodadas de teste não relacionadas e
arquivos não commitados de outra sessão concorrente no working tree produzem
violações "novas" que são artefato de ambiente, não do diff real (ver B-9b em
`.mnfs/HARNESS-DEBTS.md`, commit `c3995765`).

1. `git worktree add <scratch>/gov-base-<sha> <sha>` — worktree isolado, fora
   da árvore do repo (scratch da sessão), HEAD destacado.
2. `.gomodcache` espelhado do checkout principal via `robocopy /MIR` (Windows;
   precisa `MSYS_NO_PATHCONV=1` + `cygpath -w` nos dois caminhos, senão o MSYS
   reescreve `/E` como `E:/` e o robocopy falha). Pré-condição verificada:
   `go.sum` idêntico byte-a-byte entre `b759e2d7` e `6147c0d1` (nenhuma mudança
   de dependência no diff), então espelhar o cache populado é seguro — não
   precisou `go mod download all`.
3. `scripts/harness.ps1 -Command governance -BaseSha <sha>` rodado dentro do
   worktree isolado, saída completa redirecionada a arquivo (nunca truncada
   por `tail`).
4. Repetido em um SEGUNDO worktree isolado, no commit de merge (`6147c0d1`),
   mesmo espelhamento de `.gomodcache`.
5. Blocos com `path=` sob `.worktrees/` ou `.claude/worktrees/` são ruído de
   varredura (bug conhecido — a lane não exclui worktrees-irmãos do scan; ver
   B-10b, ainda não corrigido). Filtrado dos dois lados antes do diff.
6. Diff feito por chave `(error_code, id)`, não por `path=` — path é cego a
   classes reportadas em nível de módulo/registro (`GOV_MODULE_COVERAGE`,
   `RCFG_READER_MISSING`) que podem ser disparadas por um diff sem tocar o
   arquivo onde a violação é reportada.

## Resultado

- Blocos reais (excluindo ruído de worktree e o cabeçalho `status=`): **61**
  em cada lado.
- Chaves únicas `(error_code, id)`: **45** em cada lado.
- `base - merge` (chaves que sumiram): vazio.
- `merge - base` (chaves novas): vazio.
- **Conjuntos idênticos.** O merge da F-A1b não introduziu nem removeu
  nenhuma violação de governança de nenhuma classe.

## Lista canônica — 45 chaves `(error_code, id)` no tip `b759e2d7`

```
GOV_MODULE_COVERAGE	sourcekind
GOV_MODULE_DEPENDENCY	catalog-erp_import
GOV_MODULE_DEPENDENCY	internal_read-erp_import
GOV_MODULE_DEPENDENCY	listings-sync
GOV_MODULE_DEPENDENCY	market-erp_import
GOV_MODULE_DEPENDENCY	orders-pricing
GOV_MODULE_DEPENDENCY	product_links-integrations
GOV_MODULE_DEPENDENCY	sync-internal_read
GOV_MODULE_LAYER	catalog-erp_import-adapters
GOV_MODULE_LAYER	internal_read-erp_import-adapters
GOV_MODULE_LAYER	listings-connectors-adapters
GOV_MODULE_LAYER	listings-sync-adapters
GOV_MODULE_LAYER	market-erp_import-adapters
RCFG_READER_MISSING	MC_ERP_SOURCE
RCFG_READER_MISSING	MPC_PROVIDER_AMAZON_APPLICATION_ID
RCFG_READER_MISSING	MPC_PROVIDER_AMAZON_AUTH_VERSION
RCFG_READER_MISSING	MPC_PROVIDER_AMAZON_CLIENT_ID
RCFG_READER_MISSING	MPC_PROVIDER_AMAZON_CLIENT_SECRET
RCFG_READER_MISSING	MPC_PROVIDER_MAGALU_CLIENT_ID
RCFG_READER_MISSING	MPC_PROVIDER_MAGALU_CLIENT_SECRET
RCFG_READER_MISSING	MPC_PROVIDER_SHOPEE_BASE_URL
RCFG_READER_MISSING	MPC_PROVIDER_SHOPEE_PARTNER_ID
RCFG_READER_MISSING	MPC_PROVIDER_SHOPEE_PARTNER_KEY
RCFG_UNAPPROVED_READER	API_PORT
RCFG_UNAPPROVED_READER	MC_DATABASE_URL
RCFG_UNAPPROVED_READER	MC_ERP_SOURCE
RCFG_UNAPPROVED_READER	MPC_ENCRYPTION_KEY
RCFG_UNAPPROVED_READER	MPC_ORACLE_CONNECT_STRING
RCFG_UNAPPROVED_READER	MPC_ORACLE_LIB_DIR
RCFG_UNAPPROVED_READER	MPC_ORACLE_PASSWORD
RCFG_UNAPPROVED_READER	MPC_ORACLE_USERNAME
RCFG_UNAPPROVED_READER	MPC_TEST_DATABASE_URL
RCFG_UNAPPROVED_READER	RUN_MIGRATIONS
RCFG_UNAPPROVED_READER	SANKHYA_ORACLE_HOST
RCFG_UNAPPROVED_READER	SANKHYA_ORACLE_PASSWORD
RCFG_UNAPPROVED_READER	SANKHYA_ORACLE_PORT
RCFG_UNAPPROVED_READER	SANKHYA_ORACLE_SERVICE_NAME
RCFG_UNAPPROVED_READER	SANKHYA_ORACLE_USER
RCFG_UNAPPROVED_READER	SERVER_ADDR
RCFG_UNAPPROVED_READER	VITE_API_BASE_URL
RCFG_UNDECLARED_READ	MPC_DOMAIN
RCFG_UNDECLARED_READ	MPC_GHCR_OWNER
RCFG_UNDECLARED_READ	MPC_IMAGE_SERVER
RCFG_UNDECLARED_READ	MPC_IMAGE_TAG
RCFG_UNDECLARED_READ	MPC_IMAGE_WEB
```

## Nota de correção

Uma medição anterior desta rodada (comunicada em mensagem para a sessão
"Análise fiscal e simulação P2B") citou 46 chaves. Número errado — bug no
parser de extração: a linha de cabeçalho `status=failed` virava um bloco
fantasma sem `error_code`, contado como uma 46ª chave `(None, '')`. Corrigido
aqui: são **45**. O diff vazio (zero chaves novas, zero removidas) não muda —
o artefato do parser era simétrico nos dois lados e não afetava a comparação,
só a contagem total reportada.
