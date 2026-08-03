# Evidência — F-A1b: Reautorização da conta ML

Fatia: `docs/superpowers/plans/2026-08-03-fa1b-reautorizacao-plan.md`
Merge: `6147c0d1` (`git merge worktree-fa1b-reautorizacao --no-edit`, tip pré-merge `b759e2d7`)
Método de execução: subagent-driven-development (implementador → review de
spec → review de qualidade, por task).

## Tasks → commit → o que prova

| Task | Commit | O que prova |
|---|---|---|
| 1 — botão "Reautorizar" (D-A) | `b6b5e329` | Operador tem ação funcional pra reautorizar sem passar pelo fluxo de "Conectar" (que arriscava criar 2ª instalação). |
| 2 — pré-checagem de vínculo (D-B) | `bf1c29fc` | `ensureProviderAccountUnlinked` checa, antes de qualquer escrita, o mesmo predicado do índice único parcial `(tenant_id, provider_code, external_account_id)` — recusa conta já vinculada a outra instalação ANTES de gravar credencial/sessão, fechando a janela de escrita não-transacional em `applyAuthResult`/`ApplyConnectionSnapshot`. |
| 3 — motivo visível (D-C) | `2780c0b6` | Callback falho redireciona com `reason=<código fechado>` na URL (nunca o texto de erro cru — URL é logada/screenshotada); `IntegracoesPage` lê `?auth=failed&reason=...` e mostra pro operador. |
| 4 — prova em Postgres real | `2b47ca08` | Contra Postgres real (não fake): reautorizar a MESMA instalação não trombra com `uq_integration_installations_active_provider_account`; uma 2ª instalação ativa na MESMA conta É recusada pelo índice parcial real; desconectar a 1ª libera a conta pra 2ª. Must-fail: índice comentado temporariamente em `migrations/0017_oauth_credential_lifecycle.sql`, `failure_token=test=TestReauthKeepsOneRowAndDisconnectedReleasesTheAccount` observado na lane hermética, migration restaurada a diff zero depois. |
| Merge | `6147c0d1` | `go build ./...` e `go test ./internal/modules/integrations/...` verdes no resultado do merge. |

## Zero mudança de contrato

Nenhuma rota nova, nenhum campo novo em OpenAPI/SDK. O botão "Reautorizar"
reusa o fluxo OAuth existente; a mudança é: (a) qual UI dispara o fluxo, (b)
uma checagem de recusa ANTES da escrita, (c) um `reason=` na URL de retorno
que já era passada adiante por `LegacyRedirect.tsx`.

## Controles negativos e must-fail observados

- **Task 2**: sem `ensureProviderAccountUnlinked`, o teste de integração
  (Task 4) trombaria só no índice do Postgres DEPOIS de já ter gravado
  credencial+sessão órfãs (não-transacional) — a checagem em código evita
  esse estado inconsistente antes de chegar no banco.
- **Task 4**: must-fail reproduzido na casa — índice único comentado em
  `migrations/0017_oauth_credential_lifecycle.sql`, teste rodado na lane
  hermética, `failure_token=test=TestReauthKeepsOneRowAndDisconnectedReleasesTheAccount`
  confirmado, migration restaurada a zero diff antes do commit.

## Governança — B-9 / B-9b (medição comparativa)

Preocupação levantada pela sessão "Análise fiscal e simulação P2B" (tratada
como hub, por autorização explícita do operador nesta conversa): o merge
poderia ter introduzido violação de governança não visível por filtro de
`path=` (classes reportadas em nível de módulo/registro, não de arquivo).

Resolvido por medição rigorosa em dois worktrees isolados (fora da árvore do
repo, `.gomodcache` espelhado, `go.sum` confirmado idêntico entre os dois
SHAs) — diff por chave `(error_code, id)`, não por caminho. Resultado: **45
chaves em cada lado, conjuntos idênticos, zero novas, zero removidas**.
Detalhe completo e lista canônica das 45 chaves em
[`governance-baseline-b759e2d7.md`](governance-baseline-b759e2d7.md).

No caminho, uma comparação inicial medida dentro do checkout principal (main)
mostrou 2 chaves "novas" que não sobreviveram à remedição isolada — causa:
`.gocache` aquecido por testes não relacionados + arquivos não commitados de
outra sessão concorrente no working tree do checkout principal. A sessão hub
registrou essa classe de contaminação como **B-9b** em `.mnfs/HARNESS-DEBTS.md`
(commit `c3995765`): medições comparativas de governança nunca devem rodar no
checkout principal; os dois lados vão para worktrees isolados fora da árvore
do repo, com SHA e método nomeados explicitamente no veredito.

Bug de escopo da lane (não desta fatia, pré-existente, informalmente **B-10b**
pela sessão hub): a varredura de governança não exclui `.claude/worktrees/`
nem `.worktrees/` dos worktrees-irmãos, poluindo a saída com ~88% de blocos
duplicados por worktree quando rodada de um checkout com worktrees vivos ao
lado — e a lane sai normal/não-bloqueada mesmo assim. Não corrigido aqui;
registrado como dívida de harness, não de produto.

## Dívida encontrada e não consertada (dívida de produto — não é dívida de harness)

- **D-45** — `HandleCallback` não é transacional: credencial, sessão e
  snapshot de conexão são três escritas separadas
  (`auth_flow_service.go:359,364,375`). A Task 2 fecha a janela de conta
  já-vinculada checando ANTES de escrever, mas outras falhas parciais entre
  as três escritas continuam possíveis (ex.: processo morre entre a escrita
  da credencial e a da sessão).
- **D-46** — `mapIntegrationError` (`transport/http_handler.go:53`) casa erro
  por prefixo de string, não por tipo/código — frágil a qualquer reformulação
  de mensagem de erro upstream.
- **D-47** — `fee_sync_scheduler.go:74` descarta erro do scheduler de fee sync
  com `_, _ =` — falha de sync fica silenciosa, sem log nem métrica.

## Task 5 — live drive: EXECUTADO contra a conta real

Executado em 2026-08-03, com o operador revogando a autorização na própria
conta ML real (ação dele, confirmada explicitamente antes de qualquer SQL
rodar). Pré-condição corrigida antes de começar: a stack de dev estava
servindo o worktree da F-00, não a `main` com F-A1b — reconstruída e
recriada primeiro. Todos os 7 steps do plano executados e conferidos:
revogação → refresh forçado → falha do ticker persistida e visível →
tela mostra badge+motivo+botão → reautorização devolve a conta sem duplicar
instalação → controle positivo confirma recusa de vínculo duplicado →
rascunho de teste apagado. Detalhe completo, com todo log/SQL/dump de tela,
em [`TASK5-LIVE-DRIVE.md`](TASK5-LIVE-DRIVE.md). F-A1b fecha sem pendência
de validação.

## Task 7 do plano F-A1/F-A2 — desbloqueio parcial

`docs/superpowers/plans/2026-08-03-fa1-fa2-token-visivel-plan.md`, Task 7,
editado: bloqueador (2) marcado como desbloqueado por esta fatia (mesmo
evento de revogação real que a Task 5 daqui precisaria produzir cobre a
ambiguidade de "como induzir `invalid_grant`"). Bloqueador (1) (binário do
container desatualizado em relação ao tip da `main`) continua de pé,
independente desta fatia — precisa remedição no dia, comparando `CreatedAt`
do container com o commit mais recente.

## Ledger da missão

Não existe ainda arquivo de ledger em nível de missão para MIS-008 (verificado
por busca em `.mnfs/MIS-008-operacao-diaria/` — ao contrário de MIS-004/
MIS-007, que têm `HUB-LEDGER.md`/`DISPATCH-LEDGER.md`). D-45/D-46/D-47
registradas aqui, no pack de evidência da fatia; não há ledger de missão
onde duplicar o registro. Gap sinalizado, não resolvido — decisão de criar
ou não um ledger de missão fica com o operador.
