# HARNESS-DEBTS — dívidas de melhoria da harness (coletadas em campo)

`status: ABERTO — insumo para a sessão de melhoria da harness ordenada pelo operador 2026-07-29`
`ANÁLISE GLOBAL TRIPLA FECHADA 2026-07-29: .mnfs/HARNESS-DEBTS-GLOBAL/ (SYNTHESIS.md = reconciliação hub+MNOS+Sol-xhigh; vereditos verbatim ao lado) — a sessão de melhoria começa POR LÁ; este arquivo é o inventário de casos que a alimenta`
`fonte da harness: Documents\mnfs-harness (fonte da verdade) · profile local: docs/HARNESS-PROFILE.md`
`proveniência: MIS-006 CHIP-VENDAVEL (principal), M-06, retro 0.4.0 — cada item cita o caso que pagou o custo`

Formato: **facto medido → custo pago → conserto candidato**. Nada aqui é opinião sem caso.

## A. Despacho e briefs

**A-1. Fatia complexa despachada com executor econômico.** S9 reclassificada `complex` e mesmo
assim despachada sol/**low** → reprovada no P4 (F-1 função renomeada em vez de morta, F-2 valor
mágico) → rodada inteira de conserto. Candidato: matriz classe-da-fatia → tier de esforço
(complex = medium/high obrigatório); proposta já feita ao operador, pendente ratificação.

**A-2. Brief/card é ALEGAÇÃO sobre o repo e apodrece.** 3 alegações falsas no brief do M-06;
A-17 "default na seam de load" (falso — sistema é fail-closed, A-22); card do S8 com comando de
lane quebrado (`No test files found`); card do S10 herdou a falsidade do A-17. Custo: cada uma
= rodada perdida ou escalation. Candidato: passo pré-despacho OBRIGATÓRIO do orquestrador —
medir toda alegação executável do card (comandos de lane, símbolos assumidos, claims de
comportamento) antes do worker partir; alegação sem proveniência de medição não entra em card.
Caso adicional (A-27→A-28): âncora de LINHA em ruling cross-árvore apodrece — hub citou
:440-442 medido na MAIN, chip mediu :467-469 na árvore dele; mesmo facto, âncoras
incompatíveis. Ruling ancora por CONTEÚDO (rota/operationId/schema name); linha só com árvore
nomeada.

**A-3. Decisão fora do plano não existe para o worker — mas o registro é manual.** O hub tem
que lembrar de amendar pack/plano ANTES do brief a cada ruling (A-17..A-22 foram 6 amendments
manuais numa semana). Candidato: mecânica de ruling que só emite o ACK ao chip depois do commit
do amendment (hoje é disciplina, não mecanismo).

**A-4. Efforts/modelos dos workers invisíveis no ledger até dar errado.** Custo S9. Candidato:
DISPATCH-LEDGER com coluna modelo/esforço obrigatória + o gate de aceitação confere contra a
matriz A-1.

**A-5. Rodada de follow-up despachada sem snapshot dos bytes da rodada aceita.** S10: sem
snapshot, o orquestrador não conseguia PROVAR "só o arquivo de teste mudou" (base do diff era o
HEAD anterior às duas rodadas); re-revisão inteira no lugar = mais fraca que comparação de
bytes, e nomeada assim. Candidato: regra — snapshot dos arquivos aceitos ANTES de despachar
rodada de follow-up; o diff da rodada N é contra o snapshot da N-1, não contra o HEAD.

## B. Lanes e evidência

**B-1. Lane vacuosa imprime `ok` byte-idêntico ao verde real.** Três formas já pagas: env
ausente → todos skip (S5B RUN 27/PASS 1/SKIP 26); `return` seco em teste com env ausente →
`ok` SOB MUTAÇÃO (S9 M1); sweep com regex que não casa nada → zero falso (mojibake, 0 de 934).
Regra local existe (contagem RUN/PASS/SKIP/FAIL por linha, sweep com DOIS controles, `t.Skip`
nomeado). Candidato: subir ao CORE — evidência de lane sem contagem por linha e sem prova de
execução (failure_token/controle positivo) é inválida por definição.

**B-2. Worker conta o próprio resultado errado NOS DOIS sentidos.** S6: reportou 130 PASS onde
havia 190. Número que não veio de contar erra pra cima e pra baixo. Regra local existe;
candidato a core junto com B-1.

**B-3. Sandbox do worker cega lanes de `packages/*`** (esbuild sobe a árvore p/ resolver
config; sandbox nega travessia; vitest aborta ANTES da coleta). Lanes Go rodam no mesmo sandbox
— cegueira específica de FE workspace. `FINDING-sandbox-blind-fe-lane.md`. Regra chip-local
(brief declara cegueira + orquestrador mede); conserto estrutural: sandbox com leitura da raiz
do repo para resolução de config.

**B-4. Migração nova é invisível até rodar `cmd/migrate`** (13× 42703 na S5B). Já no profile
(§11+§3 @f1cba2a); candidato a core.

**B-5. "Variável setada" ≠ "banco alcançável".** S10-COND round 2: brief exigiu SKIP=0 e levou a
linha de dot-source do A-15; o env estava SET (LEN=144) apontando para porta efêmera de contêiner
de sessão morto → 3 FAIL por conexão recusada. A regra @f1cba2a prova o ENV; nada prova o
ENDPOINT. Custo: rodada do worker provando o óbvio. Candidato: exigência de SKIP=0 vem acoplada a
mecanismo de endpoint vivo — lane que boota a própria base (pg-session por checkout já existe no
repo) ou probe TCP pré-lane com token distinto (`HPG_ENDPOINT_DEAD`), nunca "confie no env".

**B-6. Descoberta de pacote da lane de integração é cega a metade da árvore.** A lane varre só
as 5 primeiras linhas de `internal/modules/**/*_test.go` por `//go:build integration`
(`scripts/harness/Postgres.psm1:42-61`): pacote sem a tag (`tenant_config`) e pacote fora de
`internal/modules/` (`internal/composition`) caem FORA de toda lane. Custo pago (CHIP-VENDAVEL):
o teste que prova que a política armazenada move página E contagem pelo reader composto pulou o
chip inteiro sem ninguém notar; rodado explícito contra banco vivo: RUN=72 PASS=72 SKIP=0.
Candidato: descoberta por população declarada (lista de pacotes de integração versionada e
conferida contra a árvore — pacote novo com teste de DB fora da lista = lane vermelha), não por
convenção de tag+caminho que falha em silêncio.

**B-7. Include de teste por NOME DE ARQUIVO EXATO = zero silencioso.** `apps/web/vitest.config.ts`
inclui `feature-products` por nome exato, não glob: segundo arquivo de teste no pacote roda em
lane NENHUMA e reporta zero. Mesma família do B-1 (verde não-discriminante), no eixo de coleta.
Candidato: include por glob de pacote + contagem de arquivos coletados conferida contra
`git ls-files '*test*'` do escopo.

**B-8. Shell do harness morre para a SESSÃO INTEIRA (subagentes incluídos) sem degradação
graciosa.** CHIP-VENDAVEL 2026-07-31: toda invocação de Bash falhou com
`syntax error: unexpected end of file from `{' command on line 70` ANTES de rodar o comando —
inclusive `true`, inclusive dentro de subagente (processo novo, erro byte-idêntico). O script
de 71 linhas é COMPOSTO em runtime (snapshot 44 linhas + wrapper gerado); o `{` órfão está na
parte gerada, sem arquivo em disco para corrigir. Custo: escada inteira caiu; 6 fatias de
trabalho (F4–F6) ficaram não-verificáveis; o chip "não consegue nem provar que não pode
provar". Wedge é POR SESSÃO (shell do hub vivo no mesmo minuto); conserto = restart da sessão.
Candidato: probe de shell (`true`) no bootstrap de toda sessão despachada — wedge aparece no
minuto 0, não depois de 6 fatias; converge com a intervenção 2 da análise global (atestação
no boot) e com o BOOT-ACK do MNOS.

**B-9. Lane de governança VERMELHA no main — verde absoluto inalcançável para qualquer chip.**
Medido 2026-07-31 em worktree LIMPO destacado no main tip `4ad36272`: status=failed, 51
violações pré-existentes (19 RCFG_UNAPPROVED_READER, 13 GOV_MODULE_DEPENDENCY, 10
RCFG_READER_MISSING, 6 GOV_MODULE_LAYER, 5 RCFG_UNDECLARED_READ, 2 GOV_MODULE_COVERAGE).
Todo chip que rode a lane falha por herança; o critério real vira "zero violação NOVA por
diff de conjunto (code/id/path) vs baseline do main tip" — que nenhum doc prescreve.
Conserto: hub salda (ou baseline-a formalmente) as 51; até lá, o critério de diff de
conjunto entra no profile como regra escrita. (CHIP-VENDAVEL A-30.)

**B-10. Policy scan varre untracked pesado do checkout — lane trava >20min.** A enumeração
de arquivos da `Policy.psm1` só exclui `.git|.mnfs|node_modules|.gomodcache|scripts/.runs|
scripts/tests|contracts/governance`; dumps untracked (`docs/design/evidence/ml-api/`) entram
no `Get-Content -Raw` e a lane no checkout do hub não termina. Caso novo da classe
"governance lane clean worktree" — a causa raiz agora tem nome: o filtro é por diretório
fixo, não por `git ls-files`. Conserto candidato: enumerar por índice do git + untracked
não-ignorado pequeno, ou excluir por tamanho. (CHIP-VENDAVEL A-30.)

## C. Gates e revisão

**C-1. Stop-the-line de CLASSE** — ratificado no profile @1889d0dd (2ª ocorrência do mesmo
padrão de defeito → parar, root-cause, conserto geral ou dívida registrada; gate avalia por
classe, não instância). Nasceu aqui; candidato a subir ao CORE. Caso de origem: superfície de
erro HTTP (2 famílias + 4 writeError) descoberta DEPOIS de 5 instâncias da mesma classe no
predicado vendável.

**C-2. "Chegada tem duas metades"** — gate que verifica que o valor CHEGA ao consumidor mas não
pergunta se vem do PRODUTOR certo aprova fiação falsa (S6 aceito com `IncludeAll` hardcoded;
custou escalation + S9 estendida). Regra local nos briefs de gate do chip; candidato a core.

**C-3. Guards de teste posicionais crescem em silêncio** (janela `slice(âncora, components:)`
engole todo caminho apendado — 2 guards quebraram na S8; um leria SEIS 500 onde exige 4).
Conserto de referência: janelas posição-independentes por lookahead, guard novo nasce por
valor. Candidato: doutrina de teste de contrato no core.

**C-5. Teste de resposta que faz Unmarshal na MESMA struct do Marshal é cego a defeito de tag
POR CONSTRUÇÃO** (as duas pontas usam as mesmas tags; o erro se cancela). S10 rodada 1: fixture
simétrica + round-trip por struct teria mergeado transposição de tag JSON viva com lane verde;
mutações M-A/M-B provaram (teste da rodada 1 PASSA nas duas). Candidato: doutrina de teste de
transporte — compare o FIO (`map[string]bool`/JSON cru), fixtures fixam UM campo por vez nas
duas direções; fixture simétrica é asserção morta (mesma classe ANCHORS-3, 2ª ocorrência).

**C-4. Gate read-only não descarrega critério de execução** (finding antigo, ANCHORS-3):
reviewer sem Bash deixa toda execução certificada pelo próprio chip. Ainda sem conserto
estrutural — hoje mitigado por evidência contada + spot-check do hub.

## D. Protocolo e infraestrutura

**D-1. Mensagens cross-session cruzam** — chip re-perguntou REQUEST que o hub já tinha decidido
(A-22) porque o ruling estava na fila de entrega. Custo baixo, confusão real. Candidato:
convenção "chip drena a fila de entrada antes de compor REQUEST" ou ids de correlação
REQUEST→RULING.

**D-2. `removal_owner` de exceção de governance exige round-trip ao hub** (S10). Candidato:
registro de donos válidos (milestones + chips enfileirados) publicado no pack para o chip
consultar sem parar.

**D-3. Hook dispatch-lint bloqueia spawn_task não-chip** (semeadura de backlog barrada por
exigir BASE-SHA/CONTRATO/etc. de chip). Candidato: distinguir "chip dispatch" de "task seed".

**D-4. Runtime da harness DRIFTED** (retro 0.4.0 deployed; runtime roda hooks antigos;
harness-sync --write gated pelo operador — memória retro-harness-mandate). Sincronizar.

**D-5. Candidatos upstream antigos ainda abertos** (memória harness-upstream-amendment-candidates):
sole-committer, pré-provisionar worktrees, `commit -F` para mensagens multilinha (PowerShell
heredoc é armadilha — o hub mesmo pagou 2× nesta sessão com heredoc/quoting).

**D-7. Stop hook: TERCEIRA acusação falsa da mesma classe (2026-07-29, durante a própria
sessão de análise de debts).** `CLOSED claimed but no evidence pack exists in this worktree` —
nenhum CLOSED foi enviado (token casado por substring dentro de prosa) e o glob `.mnfs/**`
resolveu contra o worktree stale `epic-lehmann-4ffbad`, não contra o checkout do hub. As duas
metades falsas, DE NOVO. O profile já ratificou a regra ("veredito automatizado nomeia caminho
ABS + SHA ou degrada para unknown") após a 2ª ocorrência — e o RUNTIME do hook não mudou: regra
de prosa não conserta hook. Prova viva de que ratificação sem enforcement é inerte; a profecia
do próprio profile ("an alarm wrong twice trains its reader to skip the third") cumpriu-se na
terceira. Candidato: hooks com known-answer test + degradação a `unknown` implementada no
runtime, não na doutrina.

**D-6. Quota do codex como parede invisível** (memória codex-quota-exhausted): despacho contra
quota esgotada falha tarde. Candidato: probe de quota no boot do hub e antes de cada despacho
codex.

## E. Legado do produto que multiplica erro da harness (contexto, não dívida da harness)

- 2 famílias de erro HTTP + 4 `writeError` (→ `CHIP-ERROR-UNIFY`, decisão do operador).
- Script `test` da raiz roda só workspace web (conserto = hub pós-merge, VC-7 @ed1b4183).
- Guards posicionais remanescentes em outras suítes (residual da S8).
