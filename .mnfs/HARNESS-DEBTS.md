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
Reincidência 2026-08-09 (sessão issue #2, PR #24): variante nova — exit code 107 em toda
invocação, sem mensagem alguma (a de 2026-07-31 tinha syntax error visível), `true` incluso,
sandbox on e off. Wedge apareceu MID-SESSÃO depois de horas de shell saudável (watchers de
background em loop `git fetch` + sleep rodando antes da morte). Trabalho não-shell (Read/
Edit/Write/Grep/WebFetch/gh via API pública) continuou; lanes locais e git ficaram
inexecutáveis enquanto durou. DIFERENÇA da ocorrência original: esta LIMPOU SOZINHA sem
restart da sessão (`true` voltou a passar turnos depois) — wedge 107 é transitório, o de
syntax-error 2026-07-31 era permanente. Re-probar antes de declarar a sessão perdida.

**B-9. Lane de governança VERMELHA no main — verde absoluto inalcançável para qualquer chip.**
Medido 2026-07-31 em worktree LIMPO destacado no main tip `4ad36272`: status=failed, 51
violações pré-existentes (19 RCFG_UNAPPROVED_READER, 13 GOV_MODULE_DEPENDENCY, 10
RCFG_READER_MISSING, 6 GOV_MODULE_LAYER, 5 RCFG_UNDECLARED_READ, 2 GOV_MODULE_COVERAGE).
Todo chip que rode a lane falha por herança; o critério real vira "zero violação NOVA por
diff de conjunto (code/id/path) vs baseline do main tip" — que nenhum doc prescreve.
Conserto: hub salda (ou baseline-a formalmente) as 51; até lá, o critério de diff de
conjunto entra no profile como regra escrita. (CHIP-VENDAVEL A-30.)

**B-9b. O critério de diff de conjunto foi executado pela primeira vez em 2026-08-03 — e o
método é a parte que importa.** A sessão F-A1b mediu o base `b759e2d7` e o merge `6147c0d1`
**cada um em seu próprio worktree destacado fora da árvore do repo**, com o gomodcache
espelhado nos dois (go.sum idêntico conferido antes). Diff de `(error_code, id)`: **vazio nos
dois sentidos, 45 chaves únicas idênticas**. Baseline canônico atual, portanto: **45 chaves no
tip pré-merge `b759e2d7`**, e é este número — não os 17/18/51/68 que circularam — que serve de
referência, porque é o único medido em ambiente equivalente dos dois lados.

**A lista está versionada, e é isso que impede o quinto número:**
`.mnfs/MIS-008-operacao-diaria/FA1b-reautorizacao/_evidence/governance-baseline-b759e2d7.md`
(commit `b7047dd1`) traz as 45 chaves `(error_code, id)` com o SHA e o método no cabeçalho.
Quem precisar do baseline lê esse arquivo em vez de remedir. Nota de método: o número foi
reportado primeiro como 46 e corrigido para 45 pela própria sessão — o parser contava o
`status=failed` como bloco sem `error_code`, gerando uma chave fantasma `(None,'')`. O diff
vazio não muda (o artefato fantasma era simétrico nos dois lados), só o total inflava em 1.
Vale como lembrete de que **contagem agregada tem bug próprio**: a alegação que sobrevive é o
diff de conjuntos, não o número.

O achado do caminho vale mais que o veredito: a primeira tentativa, medindo o merge **direto no
checkout main**, acusou 2 violações NOVAS que não existem — `GOV_MODULE_DEPENDENCY market-sync`
e `RCFG_DYNAMIC_READER_UNBOUNDED dynamic-reader`. Não era o merge e não era o gomodcache
(descartado copiando-o para o worktree base: resultado não mudou). Era o ambiente: o checkout
main tinha `.gocache` aquecido de builds anteriores **e trabalho não-commitado de outra sessão
na árvore de trabalho** no instante da medição.

Isso é uma terceira consequência da mesma raiz do B-10/B-10b, e a pior das três: B-10 custa
tempo, B-10b enterra o sinal em ruído, **B-9b fabrica violação que não existe e a atribui ao
diff em julgamento**. Um chip que confie nesse número é reprovado por trabalho alheio, ou —
pior no sentido inverso — um chip cuja violação real seja mascarada por ruído passa. Regra que
sai daqui: **medição de governança comparativa nunca acontece no checkout main.** Os dois lados
vão para worktrees destacados fora da árvore do repo, com cache equivalente, e o veredito só
vale se disser de qual árvore e de qual SHA cada lado veio. (Medição F-A1b, Task 6.)

**B-10. Policy scan varre untracked pesado do checkout — lane trava >20min.** A enumeração
de arquivos da `Policy.psm1` só exclui `.git|.mnfs|node_modules|.gomodcache|scripts/.runs|
scripts/tests|contracts/governance`; dumps untracked (`docs/design/evidence/ml-api/`) entram
no `Get-Content -Raw` e a lane no checkout do hub não termina. Caso novo da classe
"governance lane clean worktree" — a causa raiz agora tem nome: o filtro é por diretório
fixo, não por `git ls-files`. Conserto candidato: enumerar por índice do git + untracked
não-ignorado pequeno, ou excluir por tamanho. (CHIP-VENDAVEL A-30.)

**B-10b. Mesma causa raiz do B-10, sintoma pior: a lane varre `.claude/worktrees/` e emite
veredito sobre a ÁRVORE ERRADA.** Medido 2026-08-03 na main tip `b759e2d7`, com o worktree
`.claude/worktrees/f00-scheduler-pedidos` vivo e com trabalho não-commitado de outra sessão:

```
status=blocked
Cannot find path '...\.claude\worktrees\f00-scheduler-pedidos\...\sync\composition\installation_scheduler.go'
```

`Get-SourceFiles` não exclui `.claude/worktrees/`, então a varredura entrou no checkout de
outro branch e tropeçou num arquivo que a sessão dona estava criando naquele instante. Dois
custos distintos, e o segundo é o caro:

1. **Corrida** — arquivo enumerado e lido em momentos diferentes; a lane aborta por motivo que
   não é o código de ninguém.
2. **Contaminação silenciosa** — quando NÃO há corrida, a lane termina normalmente e a contagem
   de violações mistura arquivos de N branches. O número parece medido e não é atribuível a
   árvore nenhuma. Isso envenena exatamente o critério que o B-9 estabelece ("zero violação NOVA
   vs baseline do main tip"): baseline sujo faz a comparação mentir nas DUAS direções — esconde
   violação nova do chip e inventa violação que não é dele. Caso concreto no mesmo dia: o
   baseline `17 GOV_MODULE_DEPENDENCY + 9 GOV_MODULE_LAYER` anotado pela sessão do F-00 foi
   medido da main com o worktree dela já populado, e precisou ser remedido.

**Magnitude medida depois, e é pior do que a estimativa.** A sessão F-A1b rodou a lane da main
com BaseSha válido no mesmo dia, com 5 worktrees vivos: **589 blocos de erro, dos quais 520
(88%) vinham de `.claude/worktrees/*` e `.worktrees/*`**. Sobraram 68 violações reais do repo.
Ou seja, no caso silencioso a lane não degrada — ela inverte: a maior parte do que ela reporta
é de árvore que ninguém pediu para medir, e as 68 que importam ficam enterradas. Ninguém lê 589
blocos; lê-se o `status=failed` e conclui-se o que já se esperava. Corolário medido: os
baselines de `GOV_MODULE_DEPENDENCY` anotados no mesmo dia por duas sessões (17 e 18) são o
mesmo repo com contaminação diferente — nenhum dos dois serve como baseline.

Paliativo que funciona e resolve o ruído para todo mundo: checar o base num worktree **fora da
árvore do repo** (`git worktree add "$TMPDIR/gov-base-<sha>" <sha>`), nunca em
`.claude/worktrees/`, que só acrescenta uma árvore ao ruído da própria medição.

É a mesma família do "veredito que não nomeia a árvore": alarme errado duas vezes treina o
leitor a pular a terceira. Conserto candidato: o mesmo do B-10 (enumerar por índice do git em
vez de filtro de diretório fixo — `git ls-files` já é naturalmente escopado ao checkout) **e**
o veredito imprimir sempre caminho absoluto + SHA da árvore que mediu. Paliativo até lá: rodar
a lane de dentro do próprio worktree, nunca da main enquanto houver worktree com trabalho não
commitado. (Plano F-00, Task 0.)

**B-11. `npm run harness:governance` SEMPRE reprova — o atalho não repassa `-BaseSha`.**
Medido 2026-08-03:

```
status=failed
error_code=GOV_SEMANTIC_DRIFT
id=base-sha-invalid
```

`scripts/harness.ps1:113` exige `-BaseSha` casando `^[0-9a-f]{40}$` e sai 1 sem ele; o script
de npm invoca `harness.ps1 -Command governance` sem argumento e não há forma de passá-lo pelo
atalho. Custo pago: duas rodadas perdidas nesta sessão, e três comandos errados escritos num
plano de implementação (F-00 Tasks 0/6/7) que teriam reprovado todo worker no primeiro passo,
contra um erro que não fala nada sobre o trabalho dele. Agravante: o erro se chama
`GOV_SEMANTIC_DRIFT`, que **soa como achado de governança** — quem não abre o `.ps1` conclui
que o repo está quebrado.

Forma que de fato roda:

```bash
cd "$(git rev-parse --show-toplevel)" && pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/harness.ps1 -Command governance -BaseSha "$(git rev-parse main)"
```

Conserto candidato: ou o script de npm resolve o BaseSha sozinho (`git rev-parse main`, com o
valor impresso na saída para o veredito ser atribuível), ou o atalho sai do `package.json` para
não existir caminho que reprova por construção. Um comando publicado que nunca passa é pior que
comando nenhum — ele é copiado para cards e planos. (Plano F-00, Task 0.)

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

**D-8. Lane Oracle sancionada bloqueada pelo próprio `.env`, e o fallback do host é verde-vazio
(2026-08-04).** `scripts/run-live-oracle-docker.ps1` só aceita chaves `MPC_SANKHYA_ORACLE_*` e
rejeita por design tudo com prefixo `MPC_ORACLE_`; o `.env` do repo usa `MPC_ORACLE_*`, então a
lane responde `status=blocked reason=live Oracle .env unsupported_key=MPC_ORACLE_CONNECT_STRING`
e nunca chega a rodar. As duas grafias são legítimas — o script traduz uma na outra ao montar o
`ContainerEnvironment` — mas o `.env` vivo ficou do lado que o preflight recusa.

A metade cara é o fallback. Ao ser recusado pela lane, o caminho "natural" é rodar `go test` no
host, e no host `CGO_ENABLED=0`: o arquivo com `//go:build cgo` sai da compilação e o comando
imprime `testing: warning: no tests to run` + `PASS`, **exit 0**. Um teste live contra Oracle que
não executou nada é indistinguível de um que passou — mesma classe do
`integration-lane-failure-token` (pulado e verde são byte-idênticos). Forçar `CGO_ENABLED=1`
falha honestamente (`cgo: C compiler "gcc" not found`), mas só se alguém pensar em forçar.

Custo medido: conclui em voz alta que o ambiente não conseguia falar com Oracle — depois de ter
rodado dezenas de queries live na mesma sessão. O operador refutou (*"Como não tem Gcc se o tempo
inteiro voce compilou chamou banco e fez tudo?"*) e a medição achou o runner do repo em
`scripts/`. Capacidade ausente no host ≠ capacidade ausente.

Conserto: (a) preflight aceita as duas grafias ou o `.env` ganha as chaves `MPC_SANKHYA_*`;
(b) todo alvo de teste live com tag de build falha fechado quando a tag não é satisfeita —
`-run` que casa zero testes numa lane live é erro, não sucesso.

**D-6. Quota do codex como parede invisível** (memória codex-quota-exhausted): despacho contra
quota esgotada falha tarde. Candidato: probe de quota no boot do hub e antes de cada despacho
codex.

## E. Legado do produto que multiplica erro da harness (contexto, não dívida da harness)

- 2 famílias de erro HTTP + 4 `writeError` (→ `CHIP-ERROR-UNIFY`, decisão do operador).
- Script `test` da raiz roda só workspace web (conserto = hub pós-merge, VC-7 @ed1b4183).
- Guards posicionais remanescentes em outras suítes (residual da S8).

**D-7. Censo de pack escopado por diretório, não por fato** (CHIP-ERROR-UNIFY 2026-07-31):
censo do hub varreu `internal/modules/*/transport` e três produtores ficaram estruturalmente
fora do quadro — inline sem helper (catalog http_handler.go:440), fallback plain-text DENTRO
do `httpx.WriteJSON` (json.go:12) e `writeDeadlineExceeded` no ROUTER de plataforma
(route_deadline.go:129, o flat de maior alcance da árvore: dispara em toda rota com deadline,
módulo nenhum). Generalização a ratificar: **população se define pelo FATO (o que escreve
corpo de erro HTTP), nunca pelo diretório onde se espera que a edição caia**. Candidato:
regra de censo no profile — todo censo TOTAL declara o predicado do fato e varre a árvore
inteira; escopo por diretório só com justificativa escrita.

**D-8. Must-fail com injeção não confirmada = falso-vazio** (CHIP-ERROR-UNIFY 2026-07-31):
regex perl multi-linha `\n` não casou em arquivo CRLF (2545 linhas) → injeção nunca aconteceu
→ saída tsc VAZIA quase aprovada como "guard sustenta". Não-resultado é indistinguível de
sucesso quando a mutação não é verificada. Regra a ratificar: **must-fail só vale se a
INJEÇÃO for confirmada por observável (grep do artefato mutado) ANTES de se ler o veredito**;
em repo CRLF, todo regex `\n` multi-linha é falso-negativo silencioso. Mesma família do D-7
(instrumento que não mede o que afirma).

**D-9. Checker production-panic não reconhece re-panic obrigatório de http.ErrAbortHandler**
(CHIP-ERROR-UNIFY 2026-07-31): recover central correto TEM que re-lançar
`http.ErrAbortHandler` (contrato do net/http p/ abortar conexão em silêncio; convertê-lo em
envelope 500 seria o defeito) e o checker acusa GOV_PRODUCTION_PANIC. Exceção
`production-panic-apierror-recover-abort-handler` registrada com removal_owner=HARNESS-D-9: o
conserto de classe é o CHECKER migrar de regex textual para **go/ast** e aí sancionar o idioma
`if rec == http.ErrAbortHandler { panic(rec) }` por escopo real; quando isso pousar, a exceção
sai. REESCRITA 2026-07-31 (achado do chip, ratificado pelo hub): "ensinar o idioma" ao checker
atual é máximo local — Policy.psm1:363 é `[regex]::Matches($content, 'panic\s*\((?s:.*?)\)')`,
texto cru sem AST/escopo; aproximar o guard por texto passa a desculpar panics reais parecidos,
e afrouxamento de guard é PIOR que lista de exceção (a lista é visível, o afrouxamento não).
A migração p/ AST de quebra mata defeito latente já presente: `.*?` preguiçoso trunca
`panic(fmt.Sprintf("x %d", n))` no primeiro `)` e o fingerprint sai cortado (não mordeu ainda
porque exceções atuais usam `panic(err)`/literal simples). Exceção "temporária até o
instrumento consertar" é honesta; removal_owner de milestone inventado não seria.

**D-10. Lane hermética NÃO é superconjunto da lane de unidade** (CHIP-ERROR-UNIFY
2026-07-31, medido): Postgres.psm1:276 monta o conjunto de testes como ./tests/integration +
só módulos com `//go:build integration` nas 5 primeiras linhas — testes de unidade de
transport (ex. catalog/transport) NUNCA são compilados pela lane. Prova: defeito de shape
injetado → lane `status=passed`/EXIT=0, depois de provada capaz de vermelho
(failure_token=test=). Consequência operante: `status=passed` da integração NÃO cobre pins
de unidade; escada completa exige as DUAS lanes sempre. Candidato: lane hermética rodar
também `go test ./...` do módulo tocado, ou o profile declarar explicitamente a
não-cobertura.

**D-11. Teste de estado injetado não prova transição de estado** (MIS-007/M-09
2026-08-01, medido duas vezes na mesma milestone): `SyncHealthCard.test.tsx` injetava
`isError: true` direto no mock e passava verde; no dev stack real, com o endpoint
inalcançável, a query parava em `status=pending`/`fetchStatus=paused` e o card renderizava
SÓ o `<h2>` — `isLoading`, `isError` e `data` simultaneamente falsos, caindo em todos os
`? … : null`. O teste exercitava o RAMO de erro sem nunca atravessar a máquina de
retry/onlineManager que decide se aquele ramo é alcançado. Segunda rodada: o fix
(`networkMode: "always"`) veio com teste novo em jsdom contra recusa TCP, ficou verde, e o
live-drive do hub REFUTOU de novo com o mesmo HTML — jsdom + fetch de Node não reproduz o
`onlineManager` do navegador. Classe: **um observável que só existe no browser não pode ser
certificado fora dele**; e render com N ramos condicionais independentes e nenhum ramo total
tem combinação de estado não enumerada por construção. Candidato de conserto de classe:
(a) profile exigir, p/ critério de UI, que a prova seja live-drive do hub e não teste de
componente; (b) lint/regra proibindo cadeia de ternários sem ramo final em card que
consome query.

**D-12. Worktrees pré-provisionados nascem stale** (MIS-007/M-02 2026-08-01, relato do
chip): 3 dos 4 worktrees de feature do pool desta sessão bifurcaram de 20 a 480+ commits
atrás da base pretendida — um deles anterior a TODO o planning da MIS-007, e esse
corretamente respondeu BLOCKED em vez de chutar. Cada um se curou com `git merge --ff-only`
verificado, mas custou tempo real em várias features. Já existe candidato de emenda upstream
"pré-provisionar worktrees"; esta é a segunda evidência de campo. Conserto de classe: o
provisionamento carimbar a base e o bootstrap do chip FALHAR ALTO quando
`git merge-base` ≠ BASE-SHA do grant, em vez de deixar o chip descobrir por acidente.

**D-13. Prompt de chip manda "reporte por evento" sem dar o ENDEREÇO nem o MECANISMO**
(MIS-007 onda B, 2026-08-01, achado do operador): os dois chips da onda B terminaram o
trabalho e o hub não recebeu NADA — nem CLOSED, nem BLOCKED. O prompt de despacho dizia
"só eventos ao hub: CLOSED / BLOCKED / ESCALATION / REQUEST" e carregava o marcador
`HUB-SESSION: <local_id>` exigido pelo dispatch-lint, mas **nunca dizia por qual FERRAMENTA
mandar** (`mcp__ccd_session_mgmt__send_message` com aquele `session_id`). O chip não tem
como adivinhar, e o hub descobriu o estado só porque o operador avisou na mão — e porque
foi ler os branches. O `HUB-SESSION` no pack vira decoração: satisfaz o lint e não entrega
mensagem nenhuma. Classe: **marcador de contexto não é canal**; o lint verifica presença de
string, não existência de caminho de retorno. Conserto de classe: (a) o dispatch-lint exigir,
junto do `HUB-SESSION:`, a linha de mecanismo (ferramenta + como chamar); (b) o
`harness-worker` trazer o protocolo de retorno no próprio skill, para não depender de cada
prompt repetir; (c) o hub nunca assumir silêncio = trabalho não terminado — silêncio é
indistinguível de canal quebrado, exatamente como pulado-vs-verde na lane hermética (D-10).

**D-14. Guard de arquitetura por NOME de símbolo tem porta dos fundos silenciosa**
(MIS-007 M-03, 2026-08-01, achado do hub no gate de merge): o `archguard` do M-02 F-04
prova que nenhum sítio interativo novo alcança o ML, e é honesto — `mlExcludedSymbols` é
map de nomes EXATOS (sem wildcard), cada entrada carrega razão, e
`TestRealRepoInteractiveMLSites_MatchesAllowlist` amarra `raw == len(mlAllowlist) +
len(mlExcludedSymbols)` E exige que cada excluído ainda seja achado no scan cru, então nem
entrada morta nem sítio novo escondido passam. O resíduo é outro: a chave da exclusão é o
NOME do construtor, não a CLASSE DA ROTA. `newOrdersBuyerFiscalReaderAdapter` está excluído
porque hoje só alimenta `ordersIngestSvc` (batch, `POST /orders/import`); se alguém amanhã
religar o resultado desse MESMO construtor num caminho interativo, o guard cala — o símbolo
continua na lista. Não é regressão (antes ele estava na allowlist, igualmente permitido),
mas a promessa do guard ("nenhum sítio interativo não-allowlistado alcança o ML") é mais
larga que o que ele mede. Classe: **guard que mede o PRODUTOR quando a alegação é sobre o
CONSUMIDOR** — mesma família de `guard parcial sob frase total`. Conserto de classe: chavear
a exclusão pelo par (símbolo, consumidor) — o teste segue o identificador até o campo da
struct de serviço que o recebe e reprova se ele chegar num serviço registrado fora de
`registerBatchRoutes`. Enquanto não existir, a doc do guard tem que dizer que a garantia é
"nenhum sítio interativo NOVO", não "nenhum caminho interativo".

**D-15. Porta implementada sem chamador de produção e sem guard estrutural**
(MIS-007 M-03, 2026-08-01, C4): `OrderRepository.UpsertOrders` satisfaz `ports.OrderStore`
e não tem nenhum chamador de produção depois que F-02 fez `IngestOrder` virar o writer
único (ADR-04). Nada no repo impede um segundo caminho de escrita de orders renascer por
esse método — a garantia "um writer só" hoje é convenção, não guard. Classe: **superfície
morta que ainda satisfaz uma porta é um writer latente**. Conserto: ou apagar o método, ou
um teste que amarre a lista de implementações de `ports.OrderStore` alcançáveis a partir da
composição.

**D-16. Ticker longo sem tick inicial e sem vencimento persistido = starvation silenciosa**
(MIS-007 M-04, 2026-08-01, achado do hub no drive): `syncapp.Scheduler.Start`
(`scheduler.go:105-118`) é `time.NewTicker(interval)` puro — nenhuma execução no boot e
nenhum registro de "próximo vencimento" em lugar nenhum. Com `interval = 15min` (products)
isso é inofensivo. Com `interval = 24h` (o scheduler de listings que o M-04 registra) vira
starvation: qualquer reinício do processo em janela menor que 24h zera o ticker e a
varredura **nunca roda**. Em dev, onde o stack sobe várias vezes por dia, `sync_state` de
`listings` fica permanentemente vazio, e o card de saúde do M-09 vai reportar "nunca
sincronizado" — dizendo a verdade sobre um agendamento que de fato nunca acontece. O
mecanismo está provado por teste de integração; o defeito é de agendamento, não de job.
Classe: **cadência guardada só na memória do processo**. Conserto: derivar o vencimento de
`sync_state.last_full_sync_at` (rodar no boot se `now - last > interval`), em vez de confiar
num ticker que morre com o processo.

*Addendum (Plano F-00, Task 10, 2026-08-03): o scheduler de pedidos (`NewOrdersSchedulers`)
usa o mesmo `syncapp.Scheduler.Start` com `interval = 15min` — o caso já avaliado aqui como
"inofensivo". Sem debt nova; mesma classe, mesmo veredito.*

**D-17. Emit do `tsc` ao lado do fonte sequestra a resolução do Vite**
(limpeza de repo, 2026-08-01, achado na varredura de untracked): 191 arquivos `.js` estavam
soltos em `apps/web/src/`, `packages/*/src/` e — o pior — `apps/web/vite.config.js`. Todos
com cabeçalho `import { jsx as _jsx } from "react/jsx-runtime"`, ou seja, saída de `tsc`
rodado sem `--noEmit`. Nenhum era rastreado e o `.gitignore` não cobria nenhum, então
sobreviveram a várias missões aparecendo como ruído no `git status`.

Não é sujeira inerte. A ordem padrão de `resolve.extensions` do Vite é
`['.mjs','.js','.mts','.ts','.jsx','.tsx','.json']` — `.js` **antes** de `.tsx` — e o config
do repo não sobrescreve. Um `import "./AnunciosTable"` resolve para o emit velho, não para o
fonte. E o Vite procura o próprio arquivo de configuração em `.js` antes de `.ts`, então o
`vite.config.js` emitido ganhava do `vite.config.ts`: o build inteiro podia estar rodando
sobre configuração compilada desatualizada.

Classe: **artefato de build no diretório do fonte, com precedência sobre o fonte**. Sintoma
que produz: "editei o arquivo e a tela não mudou" — indistinguível de cache, de HMR quebrado
ou de container servindo checkout errado, que foi exatamente o tipo de tempo perdido que
esta base já pagou mais de uma vez.

Conserto aplicado: emit apagado e `.gitignore` fechado sobre os três padrões (`apps/web/src`,
`packages/*/src`, `vite.config.js`). Conserto de classe que **falta**: nada impede o emit de
voltar. Ou o `tsconfig` ganha `"noEmit": true` no lugar em que falta, ou um passo de lane
reprova quando existe `.js` sob um diretório de fonte que é 100% `.ts`/`.tsx`.

**D-18. Diretório de worktree sobrevive ao worktree, e o bind mount do Docker o tranca**
(limpeza de repo, 2026-08-01): `.claude/worktrees/` tinha 46 diretórios enquanto o git
conhecia 7. Os 39 restantes não tinham nem arquivo `.git` nem metadado em `.git/worktrees/`
— resíduo puro de chips já fechados. `git worktree prune` não os vê, porque prune limpa
metadado órfão, não diretório órfão.

Agravante medido: `marketplace-central-frontend-1` e `-backend-1` fazem bind mount da **raiz
do repo**, e a raiz contém `.claude/worktrees/`. Enquanto os containers estão de pé, todo
`rm -rf` ali dentro devolve `Permission denied`. Limpar worktree exige parar o dev stack —
acoplamento que ninguém declarou em lugar nenhum.

Classe: **remoção parcial que deixa o caro para trás**. `git worktree remove` apaga o arquivo
`.git` primeiro; se a remoção do diretório falha depois, o worktree some do `git worktree
list` e o `node_modules` fica. O sucesso aparente do comando esconde o custo real.

Conserto de classe: o hub verifica o diretório depois do `remove`, não o código de saída — e
a criação de worktree passa a ficar fora da árvore montada pelo compose, para que limpeza não
dependa de derrubar a stack.

**D-19. O dev stack não tem dono declarado, e tomar a janela alheia é indistinguível de
defeito da vítima** (2026-08-03, colisão F-00 × F-A1b, medida em container ids):

```
20:53:29Z  backend recriado no worktree do F-00 (janela anunciada pelo hub às 3 sessões)
21:04:35Z  backend recriado montando a main   <- F-A1b, sem avisar
21:28:26Z  hub detecta; sync_state ainda sem linha `orders`
21:28:40Z  backend devolvido ao worktree do F-00
```

Janela real do F-00: **11 minutos, contra um ticker de 15**. O primeiro ciclo do scheduler de
pedidos nunca disparou. O sintoma do lado da vítima — log sem ciclo, `sync_state` sem linha,
tela sem novidade — é **byte-idêntico ao defeito que ela está tentando provar que não existe**.
Se o hub não tivesse comparado o `Created` do container, a conclusão natural do F-00 seria
"meu scheduler não dispara", e ela seria falsa.

O detalhe que faz isso ser classe e não acidente: a sessão que tomou a janela **fez a checagem
certa** (detectou que o stack servia outra árvore, exatamente a verificação anti-binário-velho
que a doutrina exige) e mesmo assim errou o passo seguinte. Observar "o stack serve código de
outra fatia" é a prova de que existe janela em curso de alguém — é motivo para PEDIR a janela,
nunca para tomá-la. Não havia nada no ambiente que dissesse de quem era.

Conserto candidato, na ordem do mais barato: (1) o próprio compose grava o dono da janela num
lugar legível — label no container ou arquivo `.mnfs/STACK-OWNER` com sessão, branch e SHA — e
o `docker inspect` responde "de quem é isto" sem perguntar a ninguém; (2) recriar o backend
passa por comando único do hub que recusa se houver dono diferente registrado; (3) o veredito
de qualquer live drive cita o container id e o `Created`, para que a janela seja reconstituível
depois. Enquanto (1) não existir, o custo cai sempre na vítima, que é quem menos pode detectá-lo.

**D-20. Instalação conectada após o boot não ganha scheduler** (Plano F-00, Task 10,
2026-08-03): `resolveOrdersSchedulers` (mesmo padrão de `resolveListingsSchedulers`) tira um
retrato ÚNICO das instalações no boot do processo e monta um scheduler por instalação
encontrada naquele instante. Uma instalação ML conectada depois — via fluxo de reautorização,
ou primeira conexão de um tenant novo — não ganha scheduler de pedidos até o próximo restart
do processo; em dev, onde o stack sobe várias vezes por dia, o sintoma se mascara sozinho,
mas em produção (restart raro) o atraso é real e silencioso — nenhuma tela avisa "esta
instalação ainda não tem scheduler". Classe: **fan-out por retrato único, não por população
viva**. Conserto: re-listar periodicamente as instalações (mesmo intervalo do scheduler mais
curto) ou emitir um sinal na conclusão do fluxo de conexão/reautorização que registre o
scheduler faltante sem esperar o boot.

**D-21. `listings` e `market` ainda não usam janela incremental** (Plano F-00, Task 10,
2026-08-03): as Tasks 1–3 do plano F-00 alargaram `ports.ListOrdersInput`/`OrderSource` e o
`ProviderOperationService` para aceitar `UpdatedAfter`/`Limit`/`Offset` de ponta a ponta, mas
só o job de pedidos (Task 5/7) usa a janela de fato — `listings` e `market` continuam
enumerando do zero a cada corrida. O caminho está alargado; a migração de cada consumidor
para o cursor incremental é fatia separada, ainda não despachada.

**D-22. `CODEMP` fixo em 1 no leitor de custo de pedidos** (P2/M-06, 2026-08-02; plano
propunha "D-17", já ocupado pelo emit do `tsc` — renumerado):

`internal/composition/orders_adapters.go:53` fixa `CompanyID: 1` na consulta `TGFCUS`. O
predicado de vendável ratificado cobre `CODEMP(1,2)`. Produto cujo custo só exista na
empresa 2 devolve linha nenhuma → custo desconhecido → margem some, sem nada na tela
explicando por quê.

Nenhum dos 38 pedidos medidos em 2026-08-02 caiu nisso (evidência em
`.mnfs/MIS-007-ml-sync/M-06-orders-backfill-decomposition/evidence/p2-premise-check.md`).
Bomba armada, não disparada. Dono: M-06. Registrada por decisão do operador em 2026-08-02.

**D-23. `.dockerignore` deixava 940 MB de cache Go entrar no contexto de build** (revisão da
Onda 0, 2026-08-03) — RESOLVIDA no mesmo dia, registrada pela classe:

`.dockerignore:2` era `.gocache`, padrão sem `**/`, que casa só na raiz do contexto. O
`AGENTS.md` manda usar `GOCACHE` de dentro de `apps/server_core`, então a própria doutrina
garante um segundo diretório de cache que o ignore não pegava. Medido:

```
apps/server_core/.gocache      15511 arquivos   939,9 MB   <- ia no contexto
apps/server_core/.gomodcache   11129 arquivos   220,6 MB   <- excluido por **/.gomodcache
.gocache                       10773 arquivos   560,4 MB   <- excluido por .gocache
```

Sintoma: `docker compose build backend` transferiu 185 MB em 3,5 min (320–900 KB/s no
filesystem do Windows) e uma corrida anterior passou de uma hora sem terminar, com saída
zerada porque estava atrás de `| tail`. Diagnóstico errado óbvio: "docker travado".

Duas lições, ambas de classe:

1. **Padrão de ignore sem `**/` mente por vizinhança.** Na mesma lista, `**/.gomodcache` e
   `**/node_modules` estavam certos e `.gocache` estava errado — a forma correta ao lado da
   errada faz a errada parecer intencional. Toda entrada nova de `.dockerignore` que nomeie
   diretório de cache ou artefato deve nascer com `**/`.

2. **A stack de dev NÃO precisa de rebuild para pegar código novo.**
   `docker/dev/backend.Dockerfile` não compila o servidor: instala toolchain e Oracle Instant
   Client. O binário nasce em runtime — `backend-entrypoint.sh:39,44` roda
   `go run ./apps/server_core/cmd/migrate` e `exec go run ./apps/server_core/cmd/server`
   sobre o bind mount `.:/workspace`. Para tirar binário velho basta
   `docker compose up -d --force-recreate --no-build backend`. Rodar `--build` para "pegar o
   código novo" é uma hora jogada fora contra um layer que não carrega código nenhum.

**D-24. O gate de governança lia checkouts de OUTROS branches** (revisão da Onda 0,
2026-08-03) — RESOLVIDA no mesmo dia, registrada pela classe:

`Get-SourceFiles` (`scripts/harness/Policy.psm1:196`) excluía por regex **ancorada na raiz**:

```
'^(?:\.git|\.mnfs|node_modules|apps/server_core/\.gomodcache|scripts/\.runs|scripts/tests|contracts/governance)/'
```

Âncora `^` significa "só no topo do contexto". Worktree nenhum fica no topo, então
`.worktrees/<branch>/` e `.claude/worktrees/<branch>/` — checkouts inteiros de outros
branches, com `node_modules` junto — entravam na varredura e eram lidos com
`Get-Content -Raw` e casados contra os regexes de leitura de env.

Medido no diff de conjunto da Onda 0, com o mesmo instrumento nos dois lados:

```
antes da correção   HEAD: 319 issues únicos   BASE: 56    "novos": 266
depois              HEAD:  55 issues únicos   BASE: 56    novos:     2
                    264 dos 266 vinham de .claude/worktrees/f00-scheduler-pedidos/
```

O gate reprovava — mas por 264 achados de um branch que não estava sob teste, e as duas
violações verdadeiras ficavam enterradas no meio. Custo secundário: 8102 arquivos-fonte a
mais lidos por corrida, o que levou o `governance-drift` de poucos minutos para mais de 25.

Duas lições, ambas de classe:

1. **Exclusão de instrumento não se ancora na raiz.** O que se exclui de uma varredura é
   *tipo de diretório* (derivado, vendorizado, checkout alheio), e tipo de diretório aparece
   em qualquer profundidade. `^dir/` só está certo para caminho único e fixo do repo
   (`scripts/.runs`, `contracts/governance`); para o resto a forma é `(?:^|/)dir/`.
2. **Verde e vermelho mentem igual quando a população está errada.** Este gate estava
   VERMELHO, o que parece o lado seguro do erro — e ainda assim era inútil, porque ninguém
   ia caçar 2 achados reais dentro de 266. Instrumento com população errada não erra só de
   um lado.

Correção: `(?:^|/)(?:\.git|\.claude|node_modules|\.gocache|\.gomodcache|\.worktrees)/`, mais
a lista ancorada preservada para os caminhos que são de fato únicos no repo.

**D-25. `governance-drift` não calcula diff de conjunto — o `-BaseSha` só serve para uma
checagem** (revisão da Onda 0, 2026-08-03) — ABERTA:

O nome sugere que o comando compare o estado de governança do HEAD contra a base. Não
compara. `Test-GovernanceDrift` usa `$BaseSha` num único ponto (`Policy.psm1:446-461`), e
só para o `GOV_API_SDK_SPLIT` (o OpenAPI e o `sdk-runtime` mudaram no mesmo diff?). Todo o
resto é a lista **absoluta** de violações do HEAD.

Consequência prática: o comando saía `status=failed` com 55 achados no HEAD e 56 na base —
ou seja, a onda MELHOROU o número, e mesmo assim o gate dizia "failed" com a mesma cara de
sempre. Quem lê o exit code aprende a ignorá-lo. O diff de conjunto teve que ser feito à
mão nesta revisão: worktree na base, mesmo instrumento copiado nos dois lados, e diff de
`(error_code, id, path)` — foi só assim que os 2 achados reais apareceram.

Conserto: ou o comando passa a rodar os dois lados e emitir só `HEAD \ BASE`, ou é
renomeado para `governance-snapshot` e o diff vira comando próprio. O nome atual é uma
promessa que o código não cumpre.

---

**D-34. Consertar D-30 destapa `TestModuleBoundaryADR023` (234 violações pré-existentes em
`internal/modules/`), que o build quebrado escondia** (Fecho da Fundação, Tarefa 2,
`internal/composition/module_boundary_arch_test.go`, 2026-08-06) — ABERTA, fora do escopo desta
tarefa.

Depois de fechar D-29/D-30 (`catalog.New(pool)`, ver acima), `go build ./...` e `go vet ./...`
saem limpos (EXIT 0, sem output), mas `go test ./internal/...` continua vermelho — não pela
mesma linha de D-30 (essa já não existe), e sim porque dois testes que vivem nos MESMOS pacotes
do build quebrado só agora conseguem correr:

```
--- FAIL: TestNoVendorTokenInKernel (0.00s)
    repo_test.go:62: ../kernel/channel/channel_test.go:16 adapters/vendor-token-outside-adapters: mercadolivre in string literal
    repo_test.go:62: ../kernel/channel/channel_test.go:20 adapters/vendor-token-outside-adapters: mercadolivre in string literal
    repo_test.go:62: ../kernel/channel/channel_test.go:21 adapters/vendor-token-outside-adapters: mercadolivre in string literal
    repo_test.go:62: ../kernel/channel/channel_test.go:33 adapters/vendor-token-outside-adapters: mercadolivre in string literal
    repo_test.go:62: ../kernel/channel/channel_test.go:43 adapters/vendor-token-outside-adapters: mercadolivre in string literal
    repo_test.go:62: ../kernel/channel/channel_test.go:51 adapters/vendor-token-outside-adapters: mercadolivre in string literal
    repo_test.go:62: 6 vendor tokens in the kernel
FAIL
FAIL	marketplace-central/apps/server_core/internal/arch	3.947s

--- FAIL: TestModuleBoundaryADR023 (0.06s)
    module_boundary_arch_test.go:216: 234 violation(s)
        ADR-023 §2 module boundary violated
        by origin layer: 146 adapters, 42 application, 20 composition, 12 transport, 9 ports,
        2 domain, 2 integration, 1 (module root, no layer)
        by target: 40 connectors/domain, 38 internal_read/domain, 19 connectors/application,
        15 erp_import/adapters, 13 erp_import/domain, 13 integrations/domain, 13 tenant_config,
        11 listings/domain, 9 catalog/domain, ... (61 target buckets total)
FAIL
FAIL	marketplace-central/apps/server_core/internal/composition	10.862s
```

The first is D-28 verbatim, same 6 hits, unrelated to catalog — already registered. The second
is NEW to this registry: `TestModuleBoundaryADR023` lives in package `internal/composition`
(the same package D-30 broke the build of), so before this task it never ran — `go test` reported
`[setup failed]` on the D-30 compile error before it could even attempt this test. Fixing D-30
did not introduce these 234 violations; it removed the compile error that was masking a test that
was already red against the already-committed `internal/modules/` tree. `git log` on
`internal/composition/module_boundary_arch_test.go` shows it predates this session
(`484f40db`, `9555a96c`), and every violation site printed lives under `internal/modules/**`,
the legacy tree that global constraint 1 of this plan (`fecho-global-constraints.md`) forbids
touching: *"Nada em `apps/server_core/internal/modules/` é tocado. A árvore legada é
inventariada, nunca crescida."*

Consequência: `go test ./internal/...` cannot be green while both (a) `internal/composition`
builds (which this task correctly restored) and (b) `internal/modules/` carries its
pre-existing ADR-023 debt un-migrated. These two facts are in tension only because the arch test
that measures (b) happens to live in the same Go package whose build (a) unblocks — an
accident of package layout, not a defect this task created. This plan's own task list already
assigns the fix: T3 ("detector cross-context vê fora dos contextos") and T4 ("detector vendor
com escopo + arch-gate raízes certas") are the tasks that scope these detectors correctly
(kernel/contexts-only, per D-28's same finding) so they stop reproving against the frozen
legacy tree. Conserto candidato: `TestModuleBoundaryADR023` and `TestNoVendorTokenInKernel`
restrict their default scan root to `internal/kernel/` + `internal/contexts/` (mirroring the
scope constraint 1 already states in prose), with `internal/modules/` measured separately as
inventory, never as a pass/fail gate, until a migration task explicitly takes it on.

**D-26. As duas lanes de teste não produzem contagem por linha** (revisão da Onda 0,
2026-08-03) — ABERTA. Dois sintomas, uma causa: a lane reporta veredito e joga fora a
evidência que sustenta o veredito.

*Lane de integração — `status=passed` é byte-idêntico a "tudo pulado".* A lane roda
`go test -tags=integration ... -count=1` sem `-v` (`Postgres.psm1:276`), guarda a saída
**só quando o teste falha** (`Postgres.psm1:418-423`), e o objeto de retorno não expõe
`Stdout`. O `summary.txt` do run tem três linhas: `target`, `status`, `run_id`. Não há
RUN/PASS/SKIP/FAIL, e `failure_token=test=` só aparece na falha. Uma corrida em que todos os
testes fossem pulados imprimiria exatamente o mesmo texto de uma corrida verde.

Para obter a contagem desta onda foi preciso replicar o ciclo à mão: sessão do harness de
pé, `CREATE DATABASE mpc_test_<32 hex>` (qualquer outro nome devolve `HPG_TARGET_INVALID` —
`testsupport/postgres/target.go:19`), `go run ./cmd/testdb migrate`, e então
`go test -tags=integration -v`. Resultado real: **55 RUN / 55 PASS / 0 SKIP / 0 FAIL**, 83
migrações na primeira passada e 0 na segunda. O número é bom; o problema é que a lane não o
conta.

Conserto: acrescentar `-v` aos `$TestArguments` padrão e emitir
`run=/pass=/skip=/fail=` no `summary.txt`, junto com `failure_token=` em qualquer resultado,
não só na falha.

*Lane de front — reprova por contenção de máquina.* A suíte do `apps/web` roda em 24,6 s
sozinha e em 54,4 s concorrendo com a lane de integração. Sob carga, **3 de 601 testes
reprovaram por timeout de 5000 ms**: `ListingDetailPanel > loads one detail without
refetching the listing page`, `ImportChainPanel > renders a known protocol verbatim` e
`SyncHealthCard.realNetwork > reaches the named ErrorState ...` (esta com
`expected 0 to be greater than or equal to 1` — o fetch ainda não havia falhado dentro do
`waitFor`). Os mesmos três arquivos, sozinhos: **26/26 em 5,07 s**, com o `SyncHealthCard`
em 1114 ms — margem de 4× contra o timeout.

A lição não é "são flakes". É que **um vermelho de lane não é atribuível até você separar
carga de defeito**, e o caminho para separar é medir isolado, nunca supor. Supor flake é o
mesmo erro de supor defeito, só na direção confortável.

Conserto: subir `testTimeout` para 15000 na `vitest.config.ts` do `apps/web`, ou serializar
as lanes no runner do harness. A primeira é mais barata e não esconde nada — os testes que
importam levam ~1 s; 15 s continua reprovando pendura de verdade.

---

**D-27. A checagem de dependência de módulo é cega a import da RAIZ de um módulo**
(auditoria de arquitetura do P2.b, 2026-08-04) — ABERTA.

`Policy.psm1:322` procura as arestas com

```
["']marketplace-central/apps/server_core/internal/modules/(?<target>[a-z_]+)/(?<layer>[a-z_]+)[^"']*["']
```

O regex **exige um segundo segmento** depois do nome do módulo. Um módulo que expõe pacote na
própria raiz — `tenant_config` é o caso vivo — é importado como
`".../internal/modules/tenant_config"`, sem camada. O regex não casa, o `foreach` não roda, e
a aresta nunca é confrontada com `dependencies` em `modules.json`.

Medido: `pricing/adapters/postgres/product_fiscal_reader.go:10` importa `tenant_config`, e
`tenant_config` **não** está em `dependencies` de `pricing` (`modules.json:20`). A lane passa.
Outros 15 arquivos fazem o mesmo (`internal_read/adapters/routing/*`,
`internal_read/adapters/cache/cache.go`, `sync/composition/products_job.go`), então o padrão é
da casa, não do P2.b — o que torna a cegueira mais cara, não menos.

Consequência de classe: **qualquer módulo que ganhe um pacote na raiz vira porta aberta**. A
checagem de camada (`GOV_MODULE_LAYER`, mesma linha) herda o mesmo ponto cego.

Conserto: tornar o segundo segmento opcional — `/(?<target>[a-z_]+)(?:/(?<layer>[a-z_]+))?` —
tratando camada vazia como `root`, e decidir explicitamente se `root` entra na lista
`adapters|transport|registry` da checagem de camada. Antes de ligar, medir quantas arestas
novas aparecem: são 16 arquivos, e algumas provavelmente viram exceção declarada em vez de
conserto.

---

**D-28. `ScanVendorTokens` dispara em fixture de teste legítima do kernel** (Onda 1,
Tarefa 6, `internal/arch`, 2026-08-06) — ABERTA.

`internal/arch/repo_test.go:TestNoVendorTokenInKernel` esperava PASS puro (Step 7 do brief
da Tarefa 6 diz "TestNoVendorTokenInKernel PASS"). Medido: FAIL, 6 achados, todos em
`internal/kernel/channel/channel_test.go` — linhas 16, 20, 21, 33, 43, 51 — literais de
string `"MercadoLivre"` / `"mercadolivre"` usadas como dado de teste genérico para
`channel.ParseCode`/`channel.Code`. `channel.Code` é um tipo de valor deliberadamente aberto
(qualquer string de canal é um `Code` válido); o teste usa um nome de marketplace real só
como exemplo plausível, não como acoplamento de vendor.

O detetor está a fazer exatamente o que foi desenhado para fazer — achar o token em
qualquer literal de string, sem distinguir "dado de exemplo num teste de tipo genérico" de
"acoplamento de vendor num caminho de produção". Isto NÃO foi corrigido (nem o teste do
kernel, nem o detetor, nem a fixture) — nenhuma das duas coisas está no escopo aditivo da
Tarefa 6, e a instrução do plano é nunca enfraquecer detector nem fixture para ficar verde.

Consequência de classe: qualquer teste futuro no kernel (ou fora de `contexts/`/`adapters/`)
que use um nome de marketplace real como dado de exemplo vai reprovar `TestNoVendorTokenInKernel`.
Decisão pendente do operador: (a) trocar os literais do teste do kernel por um nome de canal
fictício, (b) restringir `ScanVendorTokens` a produção (excluir `_test.go`) quando chamado
sem `Suffix`, ou (c) aceitar `TestNoVendorTokenInKernel` como conhecido-vermelho até (a)/(b).
Nenhuma das três foi escolhida por esta tarefa.

**FECHADA 2026-08-07 (Fecho da Fundação, Tarefa 13).** Nota de correção sobre a proveniência
desta dívida: o brief da Tarefa 13 (`fecho-task-13-brief.md`) descreve D-28 como "o detector de
vendor acusava a própria lista. Resolvido por escopo no detector (Tarefa 4)" — essa frase
descreve um achado DIFERENTE (o `scan.go:34-35` auto-acusando-se, fechado em 0 pelo
`vendorRuleApplies` da Tarefa 4, ver D-35 item 1) e não o que este D-28 mede. Medido antes de
escrever esta nota: `TestNoVendorTokenInKernel` continuava FAIL em HEAD (`e46cbd4b`), 6 achados,
os MESMOS de `channel_test.go:16,20,21,33,43,51` — a Tarefa 4 nunca tocou `channel_test.go`
(D-35 item 5 confirma: "6 achados — é D-28 verbatim... inalterado"). A alegação do brief não se
sustenta na árvore; escolhida a opção (a) do parágrafo anterior, a única que não altera o
detector nem o alcance da regra.

Aplicado: `internal/kernel/channel/channel_test.go` trocou os cinco literais `"MercadoLivre"`/
`"mercadolivre"` por `"AcmeBazaar"`/`"acmebazaar"` (nome de canal fictício, sem correspondência
com nenhum vendor real). Confirmado por grep antes da troca que nenhum outro ficheiro do repo
depende deste literal a atravessar `channel.ParseCode` — `grep -rn "channel.ParseCode\|channel\.Code("
internal --include=*.go` fora de `_test.go` devolve zero resultados (o pacote não tem nenhum
chamador de produção ainda) e `grep -rln "mercadolivre\|MercadoLivre\|mercado_livre"
internal/kernel internal/contexts` só encontrava `channel_test.go` (o ficheiro corrigido) e um
comentário de exemplo em `internal/kernel/provenance/evidence.go:58` (fora de string literal,
não varrido pelo detector AST). O teste continua a exercitar exatamente o que testava —
normalização de maiúsculas/minúsculas via `ParseCode`/`String()` — só com um dado de exemplo que
não é nome de vendor.

Medido, verbatim:

```
$ cd apps/server_core && GOCACHE="$(pwd)/.gocache" go test ./internal/kernel/channel/... -v
=== RUN   TestParseCodeRejectsEmpty
--- PASS: TestParseCodeRejectsEmpty (0.00s)
=== RUN   TestParseCodeNormalisesCase
--- PASS: TestParseCodeNormalisesCase (0.00s)
=== RUN   TestNewAccountRefRejectsZeroCode
--- PASS: TestNewAccountRefRejectsZeroCode (0.00s)
=== RUN   TestNewAccountRefRejectsEmptyExternal
--- PASS: TestNewAccountRefRejectsEmptyExternal (0.00s)
=== RUN   TestAccountRefRoundTrips
--- PASS: TestAccountRefRoundTrips (0.00s)
PASS
ok  	marketplace-central/apps/server_core/internal/kernel/channel	1.143s

$ cd apps/server_core && GOCACHE="$(pwd)/.gocache" go test ./internal/arch/... -run TestNoVendorTokenInKernel -v
=== RUN   TestNoVendorTokenInKernel
--- PASS: TestNoVendorTokenInKernel (0.00s)
PASS
ok  	marketplace-central/apps/server_core/internal/arch	1.546s
```

`internal/kernel/channel/channel_test.go:16,20,21,33,43,51` no longer round-trip a vendor
literal; `TestNoVendorTokenInKernel` passes clean. This closes D-28 as originally measured. The
self-accusation finding the brief actually meant (`scan.go` matching its own token list) was
already closed by Task 4 — see D-35 item 1 — and stays closed; unaffected by this change.

---

**D-29. O teste de integração da Tarefa 10 (brief verbatim) não compila — importa um pacote
`internal/` de FORA da árvore que o protege** (Onda 1, Tarefa 10, `catalog/internal/postgres`,
2026-08-06) — ABERTA.

O brief manda escrever `apps/server_core/tests/integration/catalog_ingest_test.go` importando
`catalogpostgres "marketplace-central/apps/server_core/internal/contexts/catalog/internal/postgres"`
para construir `catalogpostgres.NewRepository`/`NewULIDFactory`/`NewSummaryReader` diretamente.
Medido, verbatim:

```
$ cd apps/server_core && GOCACHE="$(pwd)/.gocache" go test -tags integration ./tests/integration/ -run TestCatalog -v
# marketplace-central/apps/server_core/tests/integration
package marketplace-central/apps/server_core/tests/integration (test)
	tests\integration\catalog_ingest_test.go:13:2: use of internal package marketplace-central/apps/server_core/internal/contexts/catalog/internal/postgres not allowed
FAIL	marketplace-central/apps/server_core/tests/integration [setup failed]
FAIL
```

A regra de visibilidade `internal/` do Go é por CADA ocorrência do segmento no caminho, não só
a mais externa. `internal/contexts/catalog/internal/postgres` tem DOIS "internal": o de
`apps/server_core/internal/...` (importável por qualquer coisa sob `apps/server_core/`) e o de
`catalog/internal/postgres` (importável só por código sob `internal/contexts/catalog/`).
`tests/integration` satisfaz o primeiro e viola o segundo — a mesma fronteira que a Tarefa 7/8/9
constrói para proteger `domain`/`application`/adapters de outros CONTEXTOS bloqueia, como efeito
colateral não previsto pelo brief, o teste de integração do PRÓPRIO contexto, porque
`tests/integration/` fica fora da subárvore `contexts/catalog/`.

Não corrigido (nem o teste, nem a estrutura de pacotes, nem a fronteira `internal`) — a
instrução da Tarefa 10 foi "não mude nenhum dos dois lados para forçar, reporte o mismatch com
output verbatim". `go build ./internal/contexts/...`, `go vet ./internal/contexts/...` e
`go test ./internal/contexts/catalog/...` (sem a tag `integration`) passam limpos; a migração
`0097_catalog_context.sql` foi aplicada e confirmada (`\dt catalog.*` = 4 tabelas) e a PROVA
CRÍTICA de RLS cross-tenant foi obtida por SQL direto contra o Postgres do dev stack, sob uma
role de teste sem `BYPASSRLS` criada e removida na mesma sessão (a role de aplicação
`marketplace` é `Superuser`+`Bypass RLS`, então o `FORCE ROW LEVEL SECURITY` da migração nunca
seria exercitado através dela — achado relacionado, ver corpo do relatório da Tarefa 10).

**Reconfirmada 2026-08-07 (Fecho da Fundação, Tarefa 13).** `apps/server_core/tests/integration/
catalog_ingest_test.go` continua a chamar `catalog.New(pool)` (não nomeia `catalog/internal/postgres`
diretamente). Medido em HEAD (`e46cbd4b`), sem alteração de ficheiro nesta tarefa:

```
$ cd apps/server_core && GOCACHE="$(pwd)/.gocache" go vet -tags=integration ./tests/...
(EXIT 0, sem output)
```

Fechada continua fechada.

Consequência de classe: **qualquer contexto futuro que siga o mesmo padrão
`internal/contexts/<x>/internal/<adapter>` não pode ter um teste de integração fora de
`internal/contexts/<x>/` que construa o adapter diretamente** — ou o teste entra dentro da
subárvore do contexto (ex.: `internal/contexts/catalog/internal/postgres/*_test.go`, como os
outros pacotes `internal/modules/**/*_test.go` com tag `integration` já fazem — a descoberta de
pacotes do `Postgres.psm1` já cobre esse caso, ver `Get-HarnessIntegrationTestPackages`), ou o
contexto ganha uma função de composição EXPORTADA (ex. `catalog.NewPostgresModule(pool)` em
`module.go`, fora de `internal/`) que o teste de fora chama sem tocar em `internal/postgres`
diretamente. Nenhuma das duas foi escolhida por esta tarefa — decisão do operador/plano.

**FECHADA 2026-08-06 (Fecho da Fundação, Tarefa 2).** A segunda opção do parágrafo acima foi a
escolhida: `catalog.New` mudou de `New(store, ids, reader)` para `New(pool *pgxpool.Pool)` e
passou a montar `postgres.NewRepository`/`NewULIDFactory`/`NewSummaryReader` internamente
(`internal/contexts/catalog/module.go`), então nada fora de `internal/contexts/catalog/`
precisa nomear `catalog/internal/postgres`. `apps/server_core/tests/integration/catalog_ingest_test.go`
foi atualizado para os dois sítios de chamada (`catalog.New(pool)`) e compila limpo sob
`go vet -tags=integration ./tests/...` (EXIT 0). Ver D-30 para o lado do composition root.

---

**D-30. `internal/composition/catalog_wiring.go` (Tarefa 11, brief verbatim) é a MESMA classe
de D-29, agora no composition root** (Onda 1, Tarefa 11, `catalog/internal/postgres`,
2026-08-06) — ABERTA.

O brief da Tarefa 11, Step 8, manda escrever `internal/composition/catalog_wiring.go` importando
`catalogpostgres "marketplace-central/apps/server_core/internal/contexts/catalog/internal/postgres"`
para chamar `catalogpostgres.NewRepository`/`NewULIDFactory`/`NewSummaryReader` diretamente a
partir de `internal/composition`, que fica FORA da subárvore `internal/contexts/catalog/`.
Medido, verbatim:

```
$ cd apps/server_core && GOCACHE="$(pwd)/.gocache" go build ./...
package marketplace-central/apps/server_core/internal/composition
	internal\composition\catalog_wiring.go:9:2: use of internal package marketplace-central/apps/server_core/internal/contexts/catalog/internal/postgres not allowed
```

```
$ cd apps/server_core && GOCACHE="$(pwd)/.gocache" go build ./internal/composition/...
package marketplace-central/apps/server_core/internal/composition
	internal\composition\catalog_wiring.go:9:2: use of internal package marketplace-central/apps/server_core/internal/contexts/catalog/internal/postgres not allowed
```

`go build ./internal/adapters/...` (Steps 1–6, o resto da Tarefa 11) passa limpo — o defeito é
isolado ao ficheiro do Step 8. A instrução da Tarefa 11 é a mesma de D-29: "se o código do
brief não compilar, não mude nenhum dos dois lados para forçar, reporte o mismatch com output
verbatim." Não corrigido — nem `catalog_wiring.go`, nem `module.go`, nem a fronteira `internal`.
O ficheiro foi escrito e cometido tal como o brief manda (é um dos ficheiros nomeados no Step 10
do commit), porque reescrevê-lo para compilar seria escolher, sem mandato do operador, uma das
duas saídas que D-29 já apontou como pendentes (mover o teste/composição para dentro da árvore
de `catalog`, ou dar a `catalog` uma função de composição exportada tipo
`catalog.NewPostgresModule(pool)`).

Consequência: `go build ./...` e `go vet ./...` NÃO estão verdes na árvore inteira depois da
Tarefa 11. O raio de explosão é maior do que só `internal/composition` — `go test ./...` mostra
`cmd/server` e `tests/unit` também falham em `[setup failed]` pela MESMA linha, porque os dois
importam `internal/composition` transitivamente:

```
# marketplace-central/apps/server_core/cmd/server
package marketplace-central/apps/server_core/internal/composition
	internal\composition\catalog_wiring.go:9:2: use of internal package marketplace-central/apps/server_core/internal/contexts/catalog/internal/postgres not allowed
FAIL	marketplace-central/apps/server_core/cmd/server [setup failed]
# marketplace-central/apps/server_core/internal/composition
package marketplace-central/apps/server_core/internal/composition
	internal\composition\catalog_wiring.go:9:2: use of internal package marketplace-central/apps/server_core/internal/contexts/catalog/internal/postgres not allowed
FAIL	marketplace-central/apps/server_core/internal/composition [setup failed]
# marketplace-central/apps/server_core/tests/unit
package marketplace-central/apps/server_core/internal/composition
	internal\composition\catalog_wiring.go:9:2: use of internal package marketplace-central/apps/server_core/internal/contexts/catalog/internal/postgres not allowed
FAIL	marketplace-central/apps/server_core/tests/unit [setup failed]
```

Todo o resto de `go test ./...` (sem a tag `integration`) passa, incluindo os 5 testes novos de
`internal/adapters/erp/sankhyaoracle/catalogfeed` e a prova crítica da fronteira do Step 7. Isto
bloqueia a Tarefa 12 se ela depender de `go build ./...`/`go test ./...` limpo em `cmd/server` ou
`tests/unit`; a Tarefa 12 deve tratar `internal/composition` (e tudo o que o importa) como
conhecido-vermelho por esta dívida, ou o operador precisa escolher uma das duas saídas de
D-29/D-30 antes de continuar. Mesma consequência de classe de D-29: **qualquer composition root
ou teste fora de `internal/contexts/<x>/` que construa um adapter `internal/<x>/internal/<adapter>`
diretamente bate nesta parede.**

**FECHADA 2026-08-06 (Fecho da Fundação, Tarefa 2).** `catalog.New` mudou de
`New(store application.Store, ids application.IDFactory, reader port.Reader)` para
`New(pool *pgxpool.Pool)`, assemblando `postgres.NewRepository`/`NewULIDFactory`/
`NewSummaryReader` dentro de `internal/contexts/catalog/module.go`. `internal/composition/catalog_wiring.go`
perdeu o import `catalogpostgres "…/catalog/internal/postgres"` e passou a chamar só
`catalog.New(pool)`. Medido depois da mudança:

```
$ cd apps/server_core && GOCACHE="$(pwd)/.gocache" go build ./...
EXIT=0
$ cd apps/server_core && GOCACHE="$(pwd)/.gocache" go vet ./...
EXIT=0
```

`cmd/server` e `tests/unit` (que importavam `internal/composition` transitivamente e falhavam
por esta mesma linha) voltam a compilar. Ver D-34 abaixo: `go test ./internal/...` ainda não é
verde na árvore inteira, mas por duas causas SEM relação com D-29/D-30 — surgiram porque o
build deixou de mascará-las, não porque esta correção as introduziu.

**Reconfirmada 2026-08-07 (Fecho da Fundação, Tarefa 13).** Medido em HEAD (`e46cbd4b`), sem
alteração de ficheiro nesta tarefa:

```
$ cd apps/server_core && GOCACHE="$(pwd)/.gocache" go build ./...
BUILD_EXIT=0
$ cd apps/server_core && GOCACHE="$(pwd)/.gocache" go vet ./...
VET_EXIT=0
```

`internal/composition/catalog_wiring.go:10` continua a importar só `marketplace-central/apps/server_core/internal/contexts/catalog`
(a fachada `catalog.New(pool)`), nunca `catalog/internal/postgres`. Fechada continua fechada.

---

**D-31. O `main.go` verbatim do brief da Tarefa 12 não compila — assinaturas de `ScanVendorTokens`
divergem do que o brief assume** (Onda 1, Tarefa 12, `internal/arch/cmd/archscan/main.go`,
2026-08-06) — ABERTA (adaptação aplicada, sem mandato do operador).

O brief manda um slice `[]func(string) ([]arch.Finding, error){ ScanCrossContextInternal,
ScanFloatInContracts, ScanVendorTokens }`. Medido, verbatim:

```
$ cd apps/server_core && GOCACHE="$(pwd)/.gocache" go build ./internal/arch/cmd/archscan/...
internal\arch\cmd\archscan\main.go:19:3: cannot use arch.ScanCrossContextInternal (value of type func(root string) (arch.Findings, error)) as func(string) ([]arch.Finding, error) value in array or slice literal
internal\arch\cmd\archscan\main.go:20:3: cannot use arch.ScanFloatInContracts (value of type func(root string) (arch.Findings, error)) as func(string) ([]arch.Finding, error) value in array or slice literal
internal\arch\cmd\archscan\main.go:21:3: cannot use arch.ScanVendorTokens (value of type func(root string, tokens []string) (arch.Findings, error)) as func(string) ([]arch.Finding, error) value in array or slice literal
```

Duas divergências reais entre o brief e o `internal/arch/scan.go` já commitado na Tarefa 6:
(1) as três funções devolvem `arch.Findings` (tipo nomeado), não `[]arch.Finding`; a igualdade de
tipo de função exige tipos de retorno IDÊNTICOS, não apenas atribuíveis, então a lista não
compila mesmo corrigindo só o nome; (2) `ScanVendorTokens(root string, tokens []string)` tem dois
parâmetros, não um — a Tarefa 6 desenhou-o para aceitar a lista de tokens do chamador, o brief da
Tarefa 12 assume uma assinatura de um parâmetro só.

Corrigido (sem instrução explícita, porque "escreve exactamente o que lá está" produz um binário
que não compila): o slice usa `func(string) (arch.Findings, error)`, e `ScanVendorTokens` entra
como closure `func(root string) (arch.Findings, error) { return arch.ScanVendorTokens(root,
arch.VendorTokens) }`. O comportamento observável — três detetores correm, achados concatenados,
`file:line: regra: detalhe` por linha, exit 1 se houver achado — é idêntico ao que o brief pede;
só a forma da assinatura mudou para bater com o T6 real. Decisão pendente do operador: aceitar a
adaptação, ou ratificar uma mudança de assinatura em `scan.go` que faça o brief verbatim compilar
(fora do escopo desta tarefa — mexeria em Tarefa 6 já commitada).

---

**D-32. `grep -rn 'float64\|float32'` do portão dispara em COMENTÁRIO, não só em código —
"float in kernel" nunca mede 0 enquanto o kernel tiver a explicação do próprio banimento em
prosa** (Onda 1, Tarefa 12, `scripts/arch-gate.sh` Step "no float in the kernel" +
`.mnfs/MEASUREMENTS/2026-08-06-fundacao-kernel.md`, 2026-08-06) — **FECHADA 2026-08-06 (Fecho da
Fundação, Tarefa 4), reconfirmada 2026-08-07 (Tarefa 13).**

A medição de fecho (Tarefa 12, Step 8) esperava `float in kernel: 0`. Medido, verbatim:

```
$ cd apps/server_core && grep -rn 'float64\|float32' internal/kernel
internal/kernel/exact/decimal.go:2:// from float64 anywhere in this package, and that is the point: a binary float
internal/kernel/exact/money.go:44:// constructor from float64.
```

Os dois hits são comentários que EXPLICAM a proibição (`decimal.go:2`: "... this package never
constructs a Decimal from float64 anywhere..."; `money.go:44`: aviso de que não há construtor a
partir de float64) — não há nenhum `float64`/`float32` em código executável no kernel. O mesmo
`grep` cru corre dentro de `scripts/arch-gate.sh` (Step "no float in the kernel", brief verbatim),
então o portão, tal como especificado, **nunca reporta `kernel: no float` enquanto esses dois
comentários existirem** — falso positivo estrutural do instrumento, não uma violação real. Isto
foi medido e provado no Step 6 (prova de que o portão reprova): depois de remover o ficheiro de
sonda `probe_float.go`, o bloco "no float in the kernel" continuou a listar os DOIS comentários
acima e a marcar `fail=1`, mesmo sem nenhum `float64`/`float32` em código. Não corrigido (nem o
`grep` do portão, nem os comentários do kernel) — corrigir qualquer um dos dois lados sem mandato
seria escolher, sem instrução, entre reescrever os comentários para não conterem a palavra
literal ou trocar o `grep` cru por algo AST-consciente (o mesmo problema que `internal/arch`
resolve para os outros dois detetores, e que este Step do portão deliberadamente NÃO usa).
Consequência de classe: qualquer comentário futuro em `internal/kernel` que mencione
`float64`/`float32` em prosa (mesmo para dizer "isto é proibido") reprova o portão pela mesma
razão.

**FECHADA 2026-08-06 (Fecho da Fundação, Tarefa 4), reconfirmada 2026-08-07 (Tarefa 13).** A
Tarefa 4 tirou o passo `grep` cru do portão; a deteção de float no kernel passou a correr por
dentro do `archscan` (detector AST `ScanFloatInContracts`, o mesmo instrumento que já resolvia
D-28/D-30 para os outros dois detetores). Medido em HEAD (`e46cbd4b`):

```
$ grep -n "float in the kernel" scripts/arch-gate.sh
(sem resultado — o Step "no float in the kernel" e o grep cru não existem mais no script)

$ cd apps/server_core && GOCACHE="$(pwd)/.gocache" go run ./internal/arch/cmd/archscan -root internal/kernel
archscan internal/kernel: zero findings (exit 0)
```

Os dois comentários que disparavam o falso positivo (`internal/kernel/exact/decimal.go:2`,
`internal/kernel/exact/money.go:44`) continuam no ficheiro — não precisaram de reescrita, porque
o detector AST varre literais e identificadores de código, não comentários. `bash
scripts/arch-gate.sh` step "architecture detectors" para `internal/kernel` reporta `zero
findings`; o grep cru que disparava em prosa não existe mais no script (confirmado por grep
acima). Fecha exatamente como a Tarefa 4 já tinha registado em D-35 item 2.

**D-33. O portão (`scripts/arch-gate.sh`), corrido do zero na árvore como está hoje, reprova —
e a razão NÃO é nada que a Tarefa 12 tenha introduzido** (Onda 1, Tarefa 12, medição de fecho,
2026-08-06) — ABERTA, soma de causas já conhecidas (D-28, D-30) mais três achados novos.

`./scripts/arch-gate.sh; echo "EXIT=$?"` dá `ARCH GATE: FAIL`, `EXIT=1`, todas as vezes que foi
corrido nesta sessão — antes de injetar a sonda do Step 6, com a sonda, e depois de a remover.
Cinco causas independentes, cada uma isolada por bloco do portão:

1. **`gofmt`** lista ~637 ficheiros como não formatados, cobrindo praticamente toda a árvore
   `internal/modules` e `internal/composition`, não só código tocado por este plano. `git add`
   nesta sessão emitiu o aviso `LF will be replaced by CRLF the next time Git touches it`, ou
   seja o repo está configurado com `core.autocrlf` a normalizar para CRLF no checkout; o `gofmt`
   deste ambiente Windows/Git-Bash formata para LF, então QUALQUER ficheiro do repo (não só os
   desta tarefa) aparece como "não formatado" por causa do fim-de-linha, não do conteúdo. Não é
   um achado do protocolo de módulo — é um artefacto de ambiente anterior a este plano inteiro.
2. **`go vet ./...`** falha exatamente como D-30 documenta: `internal/composition/catalog_wiring.go:9`
   importa `catalog/internal/postgres` de fora da subárvore do contexto. Conhecido, não corrigido
   por instrução explícita do brief da Tarefa 11.
3. **`archscan -root internal`** (o portão corre-o sobre TODA a `internal/`, não só
   `kernel`/`contexts`, porque o brief da Tarefa 12 manda `-root internal` verbatim) encontra
   dezenas de tokens de vendor (`mercadolivre`, `mercado_livre`, `magalu`, ...) em
   `internal/composition/*.go` e `internal/modules/**` — código legado da árvore antiga que
   `internal/modules/` ainda usa livremente porque NENHUM plano até agora migrou esses caminhos
   para o protocolo de contexto/adapter. Isto é esperado: o protocolo só é mandatório dentro de
   `internal/kernel/` e `internal/contexts/`; a Tarefa 12 escreveu o portão exatamente como o
   brief manda (`-root internal`), então ele mede also a árvore legada, ainda não migrada, e
   reprova por ela. Mais D-28 (`TestNoVendorTokenInKernel`/`channel_test.go`, 6 hits) já
   documentado.
4. **`go test ./internal/...`** falha em `internal/composition` (mesma linha de D-30) e também em
   `internal/platform/migrate`: `TestCanonicalSourceListsEveryMigrationByFullFilename` e
   `TestCanonicalSourceDoesNotDependOnCallerWorkingDirectory` — "got 84 canonical migrations, want
   83", um drift de inventário de migrações pré-existente, não relacionado a este plano nem a
   nenhum ficheiro que ele tocou.
5. **`no float in the kernel`** — ver D-32 acima, falso positivo por comentário.

Nenhuma destas cinco causas foi introduzida pela Tarefa 12, e nenhuma foi corrigida pela Tarefa
12 — o mandato era medir o portão contra a árvore real, não fazer a árvore passar. **O veredito
real do portão na árvore como está hoje é FAIL, por razões genuínas e pré-existentes**, e ficar
verde exigiria trabalho fora do escopo deste plano (normalizar terminação de linha do repo,
resolver D-29/D-30, migrar `internal/modules`/`internal/composition` para o protocolo de
contexto, ou investigar o drift de `internal/platform/migrate`). A prova crítica do Step 6 (o
portão reprova quando `float64` é injetado no kernel, e o bloco específico aponta a linha certa)
continua válida isoladamente — só o veredito AGREGADO do portão sobre a árvore inteira é FAIL.

---

**D-35. Tarefa 4 escopa `ScanVendorTokensSuffix` e restringe as raízes do portão — D-32 fecha,
D-28/D-34 persistem, e a restrição de raiz DESTAPA 42 achados novos em `internal/composition`**
(Fecho da Fundação, Tarefa 4, `internal/arch/scan.go` + `scripts/arch-gate.sh`, 2026-08-06) —
ABERTA, fora do escopo aditivo desta tarefa.

Medido depois do fix (`vendorRuleApplies` exclui `/adapters/` e `/internal/arch/`; o portão
troca `-root internal` por um loop sobre `internal/kernel internal/contexts internal/adapters
internal/composition`; o passo `grep` de float sai, o detector AST já corre dentro do
`archscan`):

1. **`archscan -root internal/arch` (RED do brief) fecha**: de 8 achados (`scan.go:34-35`
   acusando a própria lista) para 0 — `vendorRuleApplies` funciona.
2. **O falso positivo de comentário (D-32) fecha**: sem o passo `grep`, o portão não acusa mais
   `internal/kernel/exact/decimal.go:2` / `money.go:44`.
3. **`internal/modules` sai do relatório do portão**: de 487 achados sob `-root internal` para
   os números abaixo — a raiz legada não é mais varrida (constraint 1 respeitada).
4. **`internal/contexts` e `internal/adapters`: zero achados.**
5. **`internal/kernel`: 6 achados** — é D-28 verbatim (`channel_test.go:16,20,21,33,43,51`,
   `mercadolivre` em literal de string), inalterado porque `vendorRuleApplies` não exime
   `kernel/` e o brief da Tarefa 4 não pediu isso.
6. **`internal/composition`: 42 achados, NOVO nesta contagem isolada** (o D-33 já tinha citado
   `internal/composition/*.go` qualitativamente, junto com `internal/modules/**`, sem separar a
   contagem):
   ```
   internal/composition/market_adapters.go: 8 achados (linhas 15,15,229,239,239,242,263,547)
   internal/composition/market_adapters_test.go: 5 achados (314,327,339,353,356)
   internal/composition/orders_adapters.go: 4 achados (8,8,89,97)
   internal/composition/orders_ingest_adapters.go: 6 achados (7,7,26,34,60,68)
   internal/composition/pricing_adapters.go: 1 achado (35)
   internal/composition/pricing_adapters_test.go: 2 achados (44,53)
   internal/composition/root.go: 16 achados (23,23,39,380×3,396,415,590,593,594,753,861,863,965,967)
   ```
   Todos `mercadolivre`/`mercado_livre`/`magalu` em identificador ou literal — a fiação de
   composition ainda referencia vendors diretamente, não através de `adapters/`. A causa raiz da
   descrição da Tarefa 4 diz explicitamente que a Regra 2.3 "governa `contexts/`, `kernel/` e
   `composition/`", então isto não é um falso positivo do detector: é uma violação real, ainda
   não migrada, que só ficou visível porque a raiz do portão deixou de ser `internal` (que
   afogava `composition` em ruído de `internal/modules`) e passou a incluir `internal/composition`
   explicitamente — o próprio brief manda escanear essa raiz.
7. **`go test ./internal/...`** continua FAIL: `TestNoVendorTokenInKernel` (D-28, mesmo teste,
   raiz hard-coded `../kernel` dentro de `repo_test.go`, não passa pelo portão) e
   `TestModuleBoundaryADR023` (D-34, 234 violações em `internal/modules/`, detector diferente
   — `internal/composition/module_boundary_arch_test.go`, ADR-023, não Regra 2.3).
8. **`gofmt -l`** continua a listar dezenas de arquivos pré-existentes fora do escopo desta
   tarefa (causa 1 de D-33, CRLF/LF); confirmado que `internal/arch/scan.go` e
   `internal/arch/scan_test.go` (os dois arquivos que esta tarefa editou) NÃO aparecem na lista
   (`gofmt -l internal/arch/scan.go internal/arch/scan_test.go` → saída vazia).

**`bash scripts/arch-gate.sh` → `ARCH GATE: FAIL`, `EXIT=1`**, mesmo depois do fix — não pelas
duas causas que a Tarefa 4 tinha mandato de fechar (self-accusation do detector, grep em
comentário — ambas fecham, medido acima), mas por quatro causas pré-existentes e já registadas
(D-28, D-34, causa-1-do-D-33) mais uma quinta (item 6, 42 achados de composition) que a própria
restrição de raiz do brief tornou visível pela primeira vez isolada.

Consequência de classe: **"restringir a raiz do portão" é uma mudança que aumenta granularidade
de sinal, não que garante verde** — tirar `internal/modules` do denominador não zera o
numerador de `internal/composition`, só o torna legível. Conserto candidato (fora desta tarefa):
(a) migrar `internal/composition/*.go` para o padrão adapter (mover a fiação de vendor para trás
de `adapters/`, mesmo molde que T3/T4 já aplicaram ao detector), (b) decidir explicitamente se
`channel_test.go` deve trocar `"MercadoLivre"`/`"mercadolivre"` por um nome de canal fictício
(opção (a) do D-28) ou se `TestNoVendorTokenInKernel` deve restringir-se a não-`_test.go`, e
(c) só depois disso medir `bash scripts/arch-gate.sh` outra vez à espera de PASS real.

---

**D-36. `ScanFactValueDiscard` (Regra 4.2, Tarefa 5) acusa 2 sítios em ficheiros `_test.go` que
o brief não previu — nenhum é o defeito que a regra existe para apanhar** (Fecho da Fundação,
Tarefa 5, `internal/arch/scan.go`, 2026-08-06) — ABERTA, fora do escopo aditivo desta tarefa.

O brief nomeou UM sítio real, `repository.go:360` (`desc, _ := p.Description().Value()`, em
produção, sem verificação de estado antes) — esse é o defeito de D-g e está fechado por esta
tarefa (`summarise` agora lê os dois valores e propaga `DescriptionState`). O detector é
sintático por desenho ("não pode saber o tipo do recetor sem um type checker... reporta todo
`.Value()` de dois resultados cujo segundo é descartado. Um sítio que legitimamente faz isso
renomeia o método" — comentário do próprio `ScanFactValueDiscardSuffix`) e por isso, corrido
sobre as 4 raízes reais, também acusa dois sítios de TESTE que descartam o bool depois de já
terem confirmado o estado por outro caminho:

```
$ go run ./internal/arch/cmd/archscan -root internal/contexts
internal/contexts/catalog/internal/domain/product_test.go:158: facts/value-discarded: the bool from .Value() is discarded: unknown would read as the zero value
archscan: 1 finding(s)

$ go run ./internal/arch/cmd/archscan -root internal/adapters
internal/adapters/erp/sankhyaoracle/catalogfeed/mapper_test.go:37: facts/value-discarded: the bool from .Value() is discarded: unknown would read as the zero value
archscan: 1 finding(s)

$ go run ./internal/arch/cmd/archscan -root internal/kernel      # facts/value-discarded: 0 (só D-28, vendor token)
$ go run ./internal/arch/cmd/archscan -root internal/composition # facts/value-discarded: 0 (só D-35, vendor token)
```

Ambos os dois são leituras de facto onde o `Known` já foi estabelecido pela construção do
próprio teste, não pela leitura ao vivo de um facto potencialmente `Unknown`:

1. `product_test.go:158` (`TestApplyDoesNotMutateTheReceiver`) — o produto é construído em
   `obs(t, "original", "sha256:ab91")` com descrição `Known` por construção; `desc, _ :=
   p.Description().Value()` na linha 158 lê essa descrição de volta só para comparar a string
   `"original"`. Não há caminho por onde `Unknown` chegue aqui.
2. `mapper_test.go:37` (mapper do adaptador Sankhya) — a linha 34, três linhas acima, já falha o
   teste explicitamente se `obs.Description.State() != fact.Known`; a linha 37 discard o bool
   depois de o estado já ter sido verificado à parte.

Nenhum dos dois é "desconhecido vira zero" (constraint 9): em ambos o valor só é lido depois de
o estado `Known` já estar garantido por outro caminho no mesmo teste. São exactamente o caso que
o comentário do detector já previa ("um sítio que legitimamente faz isso") — mas o brief da
Tarefa 5 não deu instrução sobre COMO um sítio legítimo deve ficar silencioso perante um
detector puramente sintático (a única saída sem enfraquecer o detector é reescrever a chamada
para nomear o segundo valor, ex. `desc, known := ...; if !known { t.Fatal(...) }`, e essa
reescrita não estava no escopo desta tarefa nem foi pedida).

Consequência de classe: como D-35, um detector sintático correto por desenho destapa sítios reais
que a spec que o encomendou não tinha medido. `archscan` sobre as 4 raízes não fecha em 0 achados
totais para `facts/value-discarded` — fecha em 0 no sítio de PRODUÇÃO que a tarefa tinha mandato
de corrigir, mais 2 em teste, pré-existentes, fora do escopo aditivo. Conserto candidato (fora
desta tarefa): reescrever as duas chamadas para capturar o segundo valor com nome e afirmá-lo
explicitamente, em vez de o descartar — nenhuma mudança de produção necessária.

**D-37. Brief da Tarefa 7 assume um teste existente que não existe.** O brief
(`fecho-task-7-brief.md`) manda "o teste do mapper existente continua a passar (ajusta-o para
`NextPage`, **não o apagues**)" e instrui procurar por chamadas `.Page(` contra `Feed` em
`catalogfeed/mapper_test.go` ou similar para adaptar. Medido: `grep -rn "\.Page(" --include=*.go`
na árvore inteira do módulo devolve zero ocorrências fora da própria definição em
`mapper.go`; `catalogfeed/mapper_test.go` (124 linhas, 5 testes) exercita só `MapProduct`,
nunca `Feed.Page`. Não havia nenhum call-site para adaptar. A alegação do brief é falsa por
ausência, não por erro de linha — mesma classe do A-2 (brief/card é alegação sobre o repo e
apodrece), aqui aplicada a um brief de tarefa em vez de um card de chip. Sem custo pago (a
tarefa prosseguiu sem bloqueio: renomeei `Page`→`NextPage` per spec, `go build`/`go vet`/
`go test ./internal/adapters/erp/sankhyaoracle/...` saem verdes, e como não havia teste a
adaptar, nada foi apagado). Registrado porque o padrão A-2 já tem duas famílias diferentes de
prova (card de milestone, agora brief de tarefa de plano) e o conserto candidato do A-2 (medir
toda alegação executável do brief antes do worker partir) cobriria este caso também.

**D-38. FECHADA por remedição (revisão do fecho da Tarefa 8) — a prova 2 não estava
"estruturalmente" bloqueada; o revisor achou o caminho que a medição original nunca avaliou.**
Redação original (preservada por rastreabilidade, ver histórico git): concluía que "Oracle está
estruturalmente inacessível" neste executor, com base em (a) `go env CGO_ENABLED CC` → `0`/`gcc`
no host bare-metal e (b) `scripts/run-live-oracle-docker.ps1` só correr um binário de smoke
pré-compilado. Ambas as observações (a) e (b) continuam **verdadeiras** — mas a conclusão delas
("estruturalmente inacessível") ia longe demais: nunca avaliou o serviço `backend` do
`docker-compose.yml` na raiz do repo, que já é um container cgo de uso geral (não um runner de
smoke fixo).

Medido nesta remedição:

1. `docker-compose.yml` (`backend` service) + `docker/dev/backend.Dockerfile`: a imagem parte de
   `golang:1.25-bookworm`, define `CGO_ENABLED=1`, instala `build-essential` (gcc) via `apt-get`,
   baixa e instala o Oracle Instant Client (`ORACLE_HOME=/opt/oracle/instantclient`), e o serviço
   faz bind-mount do repo inteiro em `/workspace` com `go_mod_cache`/`go_build_cache` persistentes
   — ou seja, um C compiler não é "presumível", está confirmado por leitura do Dockerfile.
2. `docker ps` mostrava só `marketplace-central-postgres-1` de pé (backend parado); a imagem
   `marketplace-central-backend:latest` já estava construída (2 dias). `.env` tinha
   `MPC_ORACLE_USERNAME`/`MPC_ORACLE_PASSWORD`/`MPC_ORACLE_CONNECT_STRING`/`MC_DEFAULT_TENANT_ID`
   presentes (só checado `SET`/`UNSET`, nunca o valor — nenhum `docker inspect`, só
   `printenv VAR >/dev/null`). `MPC_SANKHYA_INSTANCE` estava ausente em todo o repo fora do
   próprio `cmd/catalogingest/main.go` que o exige — não há convenção prévia porque é um
   discriminador de proveniência puro (`catalogfeed.NewFeed` só valida não-vazio), então foi
   passado ad hoc via `-e` no `docker compose run`, sem alterar `docker-compose.yml`.
3. `docker compose run --rm -e MPC_SANKHYA_INSTANCE=dev-oracle-proof backend bash
   ./docker/dev/backend-entrypoint.sh bash -c "cd apps/server_core && CGO_ENABLED=1 go build -o
   /tmp/catalogingest ./cmd/catalogingest && /tmp/catalogingest"` — builda e liga. Baseline
   `catalog.products` era 0 linhas (`docker exec marketplace-central-postgres-1 psql -U
   marketplace -d marketplace_central -c "SELECT count(*) FROM catalog.products;"` → 0). Após o
   CLI correr contra o Oracle real e o Postgres real do dev stack, a mesma consulta `count(*)`
   subiu monotonicamente em checkpoints repetidos: 0 → 2390 → 5504 → 6971 → 8251 → 9731 → **9789**
   `catalog.products` / **9790** `catalog.source_observations` (última medição,
   2026-08-07T05:55:27Z, ~30 min de runtime real, container `catalogingest-proof2` ainda ativo
   completando a paginação da carteira Sankhya real — deixado a correr, não morto, para não
   truncar dado real).
4. Uma primeira tentativa (container anterior, morto por engano aos ~10 min por eu ter lido "0%
   CPU + stdout parado" como travamento) já tinha deixado 2390 linhas reais gravadas antes de eu o
   matar — o processo não estava preso, estava à espera de round-trips de rede reais para um
   Oracle real (sem `CallTimeout` na página, sem linha de progresso no stdout: só a mensagem final
   `catalog ingest report: ...` seria impressa, e eu nunca a esperei da primeira vez). Registrado
   aqui porque é a mesma classe de erro que este debt corrige: medir demasiado cedo e concluir
   "inacessível" onde só havia "lento".
5. Desta vez deixei o container `catalogingest-proof2` correr até ao fim sem interromper. Terminou
   sozinho com `EXITCODE=0` e imprimiu a linha final: `catalog ingest report: pages=53
   observed=10586 created=8196 changed=0 idempotent=2390 conflicts=94`. `count(*)` final em
   `catalog.products` = **10586**, batendo exatamente com `observed` do relatório — a paginação
   completa da carteira Sankhya real terminou com sucesso, não apenas "ainda a crescer".

Conclusão: Oracle é acessível neste executor pelo caminho `backend` do `docker-compose.yml` —
zero infraestrutura nova, zero binário fingerprinted adicional, só passar `MPC_SANKHYA_INSTANCE`
ad hoc (candidato de conserto, fora desta remedição: documentar essa env var no `.env.example`
e/ou no compose, já que hoje só existe no código do próprio `cmd/catalogingest`). A palavra
"estruturalmente" foi removida porque a evidência agora mostra o oposto: o único obstáculo real
era metodológico (não avaliar o serviço certo), não estrutural.

**D-39. `cmd/catalogingest` herda a substituição silenciosa de tenant de `pgdb.LoadConfig()` — um
operador que esquece `MC_DEFAULT_TENANT_ID` grava linhas reais sob um tenant fabricado, sem
nenhum erro.** Achado pelo revisor do fecho da Tarefa 8. `apps/server_core/cmd/catalogingest/
main.go:44-48` resolve o tenant operante via `dbCfg, err := pgdb.LoadConfig()` seguido de
`tenantID, err := tenant.Parse(dbCfg.DefaultTenantID)`; `apps/server_core/internal/platform/pgdb/
config.go:23-25` substitui silenciosamente `cfg.DefaultTenantID = "tenant_default"` quando
`MC_DEFAULT_TENANT_ID` está vazio, em vez de falhar fechado.

Isto está em tensão direta com a constraint 9 ("Desconhecido nunca vira zero, `""`, `false` ou
default plausível") especificamente para `cmd/catalogingest`: ao contrário da maior parte do que
`pgdb.LoadConfig()` alimenta, este é um comando de operador com efeito de ESCRITA real
(`RunCatalogIngest` grava em `catalog.products`/`source_observations`/`product_identifiers`/
`source_product_keys`), não um leitor nem um serviço servindo tráfego. Um operador que esqueça de
exportar `MC_DEFAULT_TENANT_ID` antes de correr o binário não recebe erro nenhum — o comando
arranca normalmente, liga a um Oracle real, e grava um lote inteiro de linhas de catálogo sob
`tenant_default`, um tenant fabricado, sem nenhum sinal de que o valor pretendido nunca foi
fornecido. Confirmado ao vivo nesta mesma sessão (prova 2 de D-38): a corrida real usou
`dbCfg.DefaultTenantID` tal como veio de `.env` sem eu precisar validar nada — se a variável
estivesse ausente, as ~9.789 linhas gravadas teriam ido para `tenant_default` silenciosamente.

Pré-existente e NÃO desta tarefa: `pgdb.LoadConfig()`/`config.go` é infraestrutura partilhada que
`cmd/server` já usa da mesma forma — `internal/composition/root.go` passa `cfg.DefaultTenantID`
para dezenas de repositórios/serviços na wiring do `NewRootRuntime` (ex. linhas 297-991,
`erpRepo`/`classRepo`/`installationRepo`/`ordersRepo`/`pricingRepo` etc.). Mudar o comportamento
de `pgdb/config.go` para falhar fechado teria raio de explosão sobre `cmd/server` inteiro — fora
do escopo desta tarefa de scrub de erro Oracle + remedição de D-38, e `pgdb/config.go` não foi
tocado.

Conserto candidato (fora desta tarefa): dar a `cmd/catalogingest` (e a outros comandos de
operador com efeito de escrita real, se existirem) uma validação PRÓPRIA logo após
`pgdb.LoadConfig()` que falha fechado se `os.Getenv("MC_DEFAULT_TENANT_ID")` estiver vazio, em vez
de aceitar o fallback silencioso de `pgdb.LoadConfig()` — sem tocar no comportamento de
`cmd/server`, que pode ter motivos distintos (multi-tenant já roteado a outro nível, por
`tenant_config`) para tolerar o default hoje.

**Reconfirmada 2026-08-07 (Fecho da Fundação, Tarefa 13) — ainda ABERTA, código não mudou de
sítio.** Medido em HEAD (`e46cbd4b`):

```
$ grep -n "DefaultTenantID\|MC_DEFAULT_TENANT_ID\|tenant_default" apps/server_core/internal/platform/pgdb/config.go
10:	DefaultTenantID string
17:		DefaultTenantID: os.Getenv("MC_DEFAULT_TENANT_ID"),
23:	if cfg.DefaultTenantID == "" {
24:		cfg.DefaultTenantID = "tenant_default"

$ grep -n "DefaultTenantID\|MC_DEFAULT_TENANT_ID\|tenant.Parse" apps/server_core/cmd/catalogingest/main.go
44:	dbCfg, err := pgdb.LoadConfig()
48:	tenantID, err := tenant.Parse(dbCfg.DefaultTenantID)
```

Mesmo `file:line` do achado original (`config.go:23-24`, `main.go:44-48`) — nenhuma das 12
tarefas deste plano tocou `internal/platform/pgdb/config.go` nem `cmd/catalogingest/main.go`
depois do fecho da Tarefa 8/9 (confirmado por `git diff --name-only 1b2ef2da..e46cbd4b`, ver
lista completa no relatório da Tarefa 13). Continua aberta; não duplicada.

**FECHADA (fecho global, onda de correção final, Item de Trabalho 1).** Guarda própria de
`cmd/catalogingest` adicionada em `apps/server_core/cmd/catalogingest/main.go:48`
(`requireTenantConfigured(os.Getenv)`, chamada logo após `pgdb.LoadConfig()` em `main.go:44-51`,
função definida em `main.go:98-123`). `pgdb.LoadConfig()`/`config.go:23-24` NÃO foi tocado — o
fallback silencioso para `cmd/server` e os demais consumidores partilhados continua exatamente
como estava; a guarda vive só em `cmd/catalogingest`, lendo `os.Getenv("MC_DEFAULT_TENANT_ID")`
diretamente (mesma leitura, sem trim, que `pgdb.LoadConfig` usa) — nunca compara contra o valor
`"tenant_default"`, só contra string vazia, então um tenant legitimamente chamado
`tenant_default` continua a funcionar e um nome de variável com erro de digitação continua a
falhar.

RED antes do código (`apps/server_core/cmd/catalogingest/main_test.go`, ainda sem
`requireTenantConfigured` definida):

```
$ cd apps/server_core && go test ./cmd/catalogingest/... -run TestRequireTenantConfigured -v
# marketplace-central/apps/server_core/cmd/catalogingest [marketplace-central/apps/server_core/cmd/catalogingest.test]
cmd\catalogingest\main_test.go:15:9: undefined: requireTenantConfigured
cmd\catalogingest\main_test.go:31:12: undefined: requireTenantConfigured
cmd\catalogingest\main_test.go:44:12: undefined: requireTenantConfigured
FAIL	marketplace-central/apps/server_core/cmd/catalogingest [build failed]
FAIL
```

GREEN depois da guarda (3 testes: variável ausente falha, variável com nome errado falha,
`tenant_default` legítimo passa):

```
$ cd apps/server_core && go test ./cmd/catalogingest/... -run TestRequireTenantConfigured -v
=== RUN   TestRequireTenantConfigured_MissingEnv
--- PASS: TestRequireTenantConfigured_MissingEnv (0.00s)
=== RUN   TestRequireTenantConfigured_TypoedVarNameStillFails
--- PASS: TestRequireTenantConfigured_TypoedVarNameStillFails (0.00s)
=== RUN   TestRequireTenantConfigured_LegitimateTenantDefaultValuePasses
--- PASS: TestRequireTenantConfigured_LegitimateTenantDefaultValuePasses (0.00s)
PASS
ok  	marketplace-central/apps/server_core/cmd/catalogingest	1.986s
```

Prova ao vivo do binário real (comando de operador, `MC_DEFAULT_TENANT_ID` propositadamente
ausente, `MC_DATABASE_URL`/`MPC_ENCRYPTION_KEY` presentes só para deixar `pgdb.LoadConfig()`
passar e a guarda ser o motivo real da falha — nunca chega a abrir pool Postgres nem Oracle,
porque a guarda corre antes de `pgdb.NewPool`/`oracleconfig.LoadConfigFromEnv`):

```
$ cd apps/server_core && unset MC_DEFAULT_TENANT_ID && MC_DATABASE_URL="postgres://example-not-dialed/db" MPC_ENCRYPTION_KEY="test-key-not-real" go run ./cmd/catalogingest; echo "EXIT_CODE=$?"
catalogingest: MC_DEFAULT_TENANT_ID is required (refusing to silently default to "tenant_default" for a live write path)
exit status 1
EXIT_CODE=1
```

Condição de remoção: já removida — a guarda é permanente enquanto `cmd/catalogingest` existir
como comando de escrita real. Reabre apenas se `main.go:48`/`requireTenantConfigured` forem
removidos ou contornados.

**D-40. `MustParseDecimal` (`apps/server_core/internal/kernel/exact/decimal.go:84`) panica sem
exceção registada — Tarefa 11 do plano `fecho`.** Exceção
`production-panic-kernel-decimal-mustparse` registrada com removal_owner=HARNESS-D-40, seguindo o
precedente de D-9 (schema `invariants.schema.json` exige `removal_owner` no padrão
`M-NN|F-NN|M-NN/F-NN|HARNESS-D-N`; inventar um milestone seria falso — este panic não pertence a
nenhum milestone aberto). `MustParseDecimal` é o idioma de constante literal em teste e
inicialização de pacote: recebe só literais de código-fonte, nunca dado de runtime. Conserto de
classe: o construtor deixar de existir (substituído por `ParseDecimal` explícito em cada
call-site, absorvendo o `error`) — quando isso pousar, a exceção sai.

**D-41. `Decimal.StringFixed` (`apps/server_core/internal/kernel/exact/decimal.go:141`) panica em
escala negativa sem exceção registada — Tarefa 11 do plano `fecho`.** Exceção
`production-panic-kernel-decimal-stringfixed-negative-scale` registrada com
removal_owner=HARNESS-D-41, mesmo motivo de D-40/D-9 para o `removal_owner`. Escala negativa é
erro de programação do chamador (constante de código), não de dado. Conserto de classe: a
assinatura passar a devolver `(string, error)` em vez de panicar — quando isso pousar, a exceção
sai.

**D-42. `Fact.MustValue` (`apps/server_core/internal/kernel/fact/knowledge.go:134`) panica quando
chamado sem `IsUsable` — Tarefa 11 do plano `fecho`.** Exceção
`production-panic-kernel-fact-mustvalue` registrada com removal_owner=HARNESS-D-42, mesmo motivo
de D-40/D-9 para o `removal_owner`. `MustValue` é o contrato do pacote: só é chamável depois de
`IsUsable` ter sido checado na mesma função — isso não é validável estaticamente hoje. Conserto de
classe: o detector `ScanFactValueDiscard` da Tarefa 5 (Regra 4.2) aprender a exigir `IsUsable` no
mesmo bloco antes de qualquer `MustValue`, tornando o panic estruturalmente inalcançável e
sancionado por AST em vez de por exceção de lista — quando isso pousar, a exceção sai.

Achado colateral medido ao registar esta exceção: o `fingerprint` de `exact-occurrence` para um
`panic(...)` multi-linha é sensível a fim-de-linha. `git worktree add` novo com
`core.autocrlf=true` (medido nesta árvore) materializa `\r\n` no `knowledge.go` checked-out,
mas a cópia de trabalho do repo principal (não re-checked-out desde antes do autocrlf, ou tocada
por ferramenta que escreve LF) tinha `\n` puro no mesmo commit — o MESMO blob git produz DOIS
textos de fingerprint diferentes consoante o caminho de checkout. `Policy.psm1:369` usa
`[regex]::Matches($content, 'panic\s*\((?s:.*?)\)')` sobre `Get-Content -Raw`, que preserva o
fim-de-linha do disco; o fingerprint registado aqui usa `\r\n` porque é isso que a lane vê no
worktree limpo (o caminho sancionado pelo brief da Tarefa 11 para correr a lane). Um panic de
UMA linha (`panic(err)`, `panic(ErrNegativeScale)`) não sofre disto — só multi-linha. Conserto
candidato (fora desta tarefa): o checker normalizar `\r\n`→`\n` antes de casar fingerprint, ou o
schema exigir fingerprints só de panics de uma linha.

---

## F. Fecho da Fundação, Tarefa 13 (2026-08-07) — dívidas registadas no fecho do plano

Sete itens vêm nomeados no brief da Tarefa 13 (`fecho-task-13-brief.md`), sem número — cada um
recebe D-número aqui pela primeira vez, para ficar rastreável como o resto deste ficheiro. Mais
dois itens que o brief predatava e o operador mandou acrescentar (D-44 e D-50).

**D-43. Agendamento do ingest de catálogo não existe — é comando de operador, e registar um
scheduler colide na chave de `sync_state`.** `apps/server_core/internal/modules/sync/domain/
sync_state.go:15-35` declara `Entity` como lista fechada e validada em aplicação (não em
constraint de DB): `EntityProducts`/`EntityListings`/`EntityOrders`/`EntityMarket`/
`EntityMarketQueue`/`EntityTariffs`/`EntityICMSMatrix` — `products` já é o job legado do módulo
`internal/modules/sync/composition/products_job.go`. O ingest de catálogo desta plano
(`cmd/catalogingest`, contexto `internal/contexts/catalog`) não tem entidade própria: registar um
segundo job sob `EntityProducts` colidiria na chave `(tenant, installation, entity)` de
`sync_state` com o job legado (a mesma proteção que a Tarefa 5 do MIS-006 desenhou, ver
`sync_state.go:22-26`), e criar uma entidade nova + scheduler exigiria tocar
`internal/modules/sync/` — fora do alcance da constraint 1 deste plano ("nada em
`internal/modules/` é tocado"). Hoje o ingest só corre por invocação manual do operador
(`go run ./cmd/catalogingest` ou o caminho Docker provado em D-38). Decisão de agendamento
pertence ao plano de migração de contextos (`internal/modules/` → `internal/contexts/`), que
ainda não existe.

**D-44. DSN de produção liga como `marketplace` (superuser, bypassa RLS) — o papel `mpc_app` da
Tarefa 10 existe só no ficheiro de migração, ainda não em uso.** Medido contra o dev stack vivo
(`marketplace-central-postgres-1`, o único Postgres deste repo acessível nesta sessão — não há
ambiente de produção alcançável daqui, então esta é a mesma medição que a produção herdaria sem
mudança de configuração):

```
$ docker exec marketplace-central-postgres-1 psql -U marketplace -d marketplace_central -c \
  "SELECT rolname, rolsuper, rolbypassrls FROM pg_roles WHERE rolname IN ('marketplace','mpc_app');"
   rolname   | rolsuper | rolbypassrls
-------------+----------+--------------
 marketplace | t        | t
(1 row)
```

`mpc_app` nem sequer aparece na lista de roles deste banco: a migração
`apps/server_core/migrations/0098_catalog_app_role.sql` (commitada pela Tarefa 10) ainda não foi
aplicada a este dev stack — `schema_migrations` marca `0097_catalog_context.sql` como a última
corrida (`2026-08-06 06:58:25+00`), nenhuma linha para `0098`. Ou seja: o papel de menor
privilégio existe como código migrável e o isolamento foi provado numa lane isolada (Tarefa 10,
role temporária sem `BYPASSRLS`, ver D-29), mas neste banco real, hoje, nem a migração rodou —
muito menos a aplicação liga como `mpc_app`. `internal/composition/root.go` e `cmd/catalogingest`
continuam a construir a `*pgxpool.Pool` a partir de `MC_DATABASE_URL`, que aponta para o utilizador
`marketplace` (`docker-compose.yml:6-7`, `POSTGRES_USER: marketplace`) — nenhum código lê uma
segunda env var para um DSN de aplicação. Consequência: `FORCE ROW LEVEL SECURITY` das tabelas
`catalog.*` nunca é exercitado pela aplicação real, porque o utilizador de ligação é superuser e
bypassa RLS por definição — o isolamento de tenant nessas tabelas depende inteiramente dos
predicados `WHERE tenant_id = $1` escritos à mão em cada query, exatamente como constraint 10
já exige, sem rede de segurança do banco por trás. Task 10's brief previa este debt ser
registrado na Tarefa 12; a Tarefa 12 foi só ADRs, por isso aterra aqui. Conserto candidato (fora
desta tarefa): segunda env var de aplicação (`MC_APP_DATABASE_URL` ou similar) apontando para
`mpc_app`, com `MC_DATABASE_URL` reservado a migrações/admin; e rodar `cmd/migrate` neste dev
stack para aplicar `0098` antes de qualquer prova de RLS em produção.

**D-45. Visibilidade de falha em ecrã — uma falha do ingest hoje é só um código de saída de
terminal.** `cmd/catalogingest/main.go:34-40` imprime o erro em `os.Stderr` e sai `1`; não existe
`sync_health` para o contexto `catalog` (a tabela `sync_state` é do módulo legado `sync`, D-43
acima, e o contexto `catalog` não escreve nela), não existe rota nem card de operador que mostre
"último ingest falhou/há quanto tempo". Um operador que não esteja a olhar para o terminal no
momento exato da corrida não tem outra forma de saber que uma corrida falhou. Escopo do próximo
plano — nenhuma tarefa deste plano tinha mandato de desenhar telemetria de operador.

**D-46. `internal/modules/` continua fora dos quatro detectores do portão.** `scripts/arch-gate.sh`
varre só `internal/kernel internal/contexts internal/adapters internal/composition` (linha 30) e
imprime em voz alta `"internal/modules is the legacy tree and is deliberately NOT scanned"`
(linha 37) — omissão declarada, não escondida. A árvore legada (61 módulos em
`contracts/governance/modules.json`) só entra no protocolo de contexto/kernel quando for migrada;
até lá o portão mede zero sobre ela por desenho, e `TestModuleBoundaryADR023` (D-34, 234
violações, ver acima) é o único instrumento que a vê — e falha o `go test ./internal/...` do
próprio portão por isso mesmo.

**D-47. A árvore de contextos continua ingovernada pelo registo de módulos — com dono e prazo.**
Três factos medidos, cada um `file:line`:

1. `scripts/harness/Policy.psm1:303` fixa `$moduleRoot = Resolve-RepositoryPath $RepositoryRoot
   'apps/server_core/internal/modules'` — o registo de governança só sabe percorrer essa raiz;
   `internal/contexts/` não está no caminho, então uma pasta órfã ali (sem entrada em
   `modules.json`) não reprova nada — fica invisível em vez de bloquear.
2. `contracts/governance/modules.json:5` já declara `{ "id": "catalog", "root":
   "apps/server_core/internal/modules/catalog", ... }` — a chave `id: "catalog"` está ocupada
   pelo módulo legado. O novo contexto `internal/contexts/catalog` não pode entrar no registo
   sob o mesmo `id` sem colidir; a chave real precisaria ser o par `(kind, id)` (módulo legado vs.
   contexto), que o schema atual não modela.
3. `scripts/harness/Policy.psm1:305-308` (mostrado acima) monta `$actualModules` só a partir de
   `Get-ChildItem` sobre `$moduleRoot` — mesmo se o schema ganhasse `kind`, o `foreach` que anda
   pelo disco continuaria sem visitar `internal/contexts/`.

Já ratificado, com dono e prazo: fecha em **T13 da adenda `1c47d906`** (registo passa a ver
contextos), não nesta tarefa — este item é a medição que essa tarefa herdará, não uma decisão
nova.

**D-48. Classe B — a mesma forma copiada à mão, quatro vezes, sem checagem de acordo.** `domain`
→ DTO → OpenAPI → `sdk-runtime` são quatro cópias independentes da mesma forma de dado, e nenhum
compilador liga uma às outras (nenhuma importa as outras). `GOV_API_SDK_SPLIT`
(`scripts/harness/Policy.psm1`, citado em D-25 acima) só exige que OpenAPI e `sdk-runtime` mudem
no MESMO commit — nunca verifica se dizem a mesma coisa. Este plano (Fecho da Fundação) é imune
por construção: constraint 12 proíbe qualquer rota HTTP nesta fatia, então não nasce DTO, não
nasce entrada de `paths` no OpenAPI, não nasce método de SDK — a quarta cópia nunca é escrita
aqui. A imunidade acaba na primeira rota que o contexto `catalog` publicar. Fecha em **T14**
(gerador + generate-and-diff), ainda não despachada.

**D-49. `docs/architecture/decisions/023-module-protocol.md` desalinhado do
`docs/superpowers/specs/2026-08-06-protocolo-de-codigo-design.md` §14.** Duas fontes descrevendo
a mesma regra (o protocolo de módulo) sem um mecanismo que as mantenha em acordo — o spec §14
regista o estado das emendas à ADR-023 (linha 616: "§14, que regista o estado das emendas à
ADR-023"), mas nada compara o TEXTO da ADR-023 contra o texto do spec depois de cada emenda.
Fecha em **T15**; os ADRs deste plano (033, 034) são disjuntos dessa regra — citam-na, não a
substituem.

**D-50. Obrigação de OPERADOR ainda não descarregada: `removal_owner` de exceção de governança
exige round-trip ao hub (D-2), e a Tarefa 11 resolveu unilateralmente usando o precedente de
mecanismo do D-9 — o mecanismo não descarrega o processo.** `contracts/governance/invariants.json`
tem hoje três exceções novas desta tarefa 11 com `removal_owner=HARNESS-D-40`,
`HARNESS-D-41`, `HARNESS-D-42` (linhas 63, 71, 82 — confirmado por `grep -n "removal_owner"
contracts/governance/invariants.json` nesta sessão), seguindo exatamente o padrão que D-9 abriu
(`removal_owner=HARNESS-D-9`, linha 31) para panics que não pertencem a nenhum milestone aberto.

D-2 (seção D acima) diz: "`removal_owner` de exceção de governance exige round-trip ao hub (S10).
Candidato: registro de donos válidos... publicado no pack para o chip consultar sem parar." Isso
nunca foi ratificado como mecanismo formal — D-9 foi um precedente de FATO, criado por uma sessão
sem esperar pelo hub, porque inventar um milestone teria sido pior (falso). A Tarefa 11 deste
plano repetiu o mesmo padrão três vezes (D-40/41/42) sem pedir ACK a ninguém, pela mesma lógica:
o schema `invariants.schema.json` exige `removal_owner` casando `M-NN|F-NN|M-NN/F-NN|HARNESS-D-N`,
e não havia milestone aberto para nomear.

O que falta não é o mecanismo (`HARNESS-D-N` já funciona, tecnicamente, três vezes) — é o
operador confirmar que usar o próprio número da dívida como dono de remoção é uma prática
ratificada, e não uma exceção tolerada caso a caso. Sem essa confirmação, D-2 continua aberta
como obrigação de PROCESSO mesmo com o mecanismo já em uso repetido. Pedido: ACK do operador
sobre o padrão `removal_owner=HARNESS-D-N` como prática permanente para exceções sem milestone —
cross-referências D-2, D-9, D-40, D-41, D-42.

**D-51. `npm run harness:unit` e `npm run harness:integration` — as duas únicas lanes invocáveis
por `npm` — não cobrem nenhum teste que este ramo (fecho global) escreveu sob
`apps/server_core/internal/...`.** `scripts/harness.ps1:61` (`Invoke-Unit`) corre só `go test
./tests/unit/... -count=1`; `scripts/harness/Postgres.psm1:276` (`Invoke-Integration`) corre só
`go test -tags=integration ./tests/integration -count=1`. Nenhuma das duas alcança
`internal/composition`, `internal/contexts/catalog/port`, `internal/kernel/channel`,
`internal/kernel/fact`, `internal/platform/migrate`, ou
`internal/adapters/erp/sankhyaoracle/catalogfeed` — os pacotes onde este ramo escreveu
`scan_test.go`, `feed_test.go`, `channel_test.go`, `combine_test.go`, `runner_test.go` e
`mapper_test.go`.

Medido nesta sessão:

```
$ git diff --name-only 1b2ef2da..a9f92809 -- apps/server_core | grep -E "_test\.go$"
apps/server_core/internal/adapters/erp/sankhyaoracle/catalogfeed/mapper_test.go
apps/server_core/internal/arch/scan_test.go
apps/server_core/internal/contexts/catalog/port/feed_test.go
apps/server_core/internal/kernel/channel/channel_test.go
apps/server_core/internal/kernel/fact/combine_test.go
apps/server_core/internal/platform/migrate/runner_test.go
apps/server_core/tests/integration/catalog_ingest_composition_test.go
apps/server_core/tests/integration/catalog_ingest_test.go
apps/server_core/tests/integration/catalog_provenance_test.go
apps/server_core/tests/integration/catalog_rls_role_test.go
```

Zero desses 6 ficheiros `internal/...` estão sob `apps/server_core/tests/unit/`; os 4 sob
`tests/integration/` SÃO cobertos por `harness:integration`, mas os 6 sob `internal/...` não são
cobertos por lane nenhuma acionável via `npm`. Confirmado correndo `npm run harness:unit`: o passo
Go só toca o pacote `tests/unit`, o resto do sinal (605 testes) é Vitest do `web`:

```
$ npm run harness:unit 2>&1 | grep -E "^ok|^---|^FAIL|target=fake|tests/unit"
target=fake
ok  	marketplace-central/apps/server_core/tests/unit	3.489s
target=fake
```

Correção à premissa original do achado (o revisor do fecho de ramo escreveu "a única lane que os
corre é `scripts/arch-gate.sh` passo 4 — que é inalcançável porque o passo 1 falha"; verificado
aqui e é FALSO): `scripts/arch-gate.sh` não usa `exit` cedo — cada passo só marca uma flag `fail`
local, e a decisão PASS/FAIL só é tomada no fim do script (linhas finais `if [ "$fail" -ne 0 ];
then exit 1`). O passo 4 (`unit tests`, `go test ./internal/...`) CORRE de facto, mesmo com o
passo 1 (`gofmt`) já tendo marcado falha:

```
$ bash scripts/arch-gate.sh > archgate_full.txt 2>&1; echo "EXIT=$?"
EXIT=1
$ grep -n "^== " archgate_full.txt
2:== gofmt ==
644:== go vet ==
646:== architecture detectors ==
700:== unit tests ==
1213:== working tree ==
```

O passo 4 de facto executa e cobre os 6 pacotes `internal/...` deste ramo (todos `ok`); a única
falha dentro dele é `TestModuleBoundaryADR023` (`internal/composition`, dívida herdada, ver
D-49/Item de Trabalho 4 — 234 violações, não deste ramo). `scripts/arch-gate.sh` NÃO é invocável
por nenhum script `npm` (confirmado por `grep -n "arch-gate" package.json` → 0 resultados) — corre
só manualmente via `bash`, então nenhum operador que só usa `npm run harness:*` alguma vez o vê.

Net: `npm run harness:unit` PASS (605 testes) é uma afirmação verdadeira que não carrega nenhuma
evidência sobre o código Go deste ramo — nem porque a lane certa esteja inalcançável (não está),
mas porque `harness:unit`/`harness:integration` nunca apontam para `internal/...` e
`arch-gate.sh`, que aponta, não está ligado a `npm`.

Condição de remoção: `scripts/harness.ps1` ganhar um passo (ou lane nova) que corra `go test
./internal/...` sem a tag `integration` a partir de `apps/server_core`, ligado a um script `npm`
— ou `arch-gate.sh` passar a ser invocável via `npm run` e o seu passo 4 reportado separado do
gofmt/vet/archscan para não ficar mascarado pela dívida ADR-023 herdada.

**FECHADA (issue #3, Tarefa 6).** A segunda das duas condições de remoção acima aconteceu, mas
não do jeito previsto: `scripts/arch-gate.sh` foi apagado no commit `0d1cc3f0` ("gate: retire
arch-gate.sh -- zero executable references, superseded by gate lanes (#3)"), depois de uma
remedição em `60f4d457` medir zero referências executáveis ao script em todo o repo (só docs,
histórico `.mnfs` e planos superados o citavam). A dívida morre com o ficheiro que carregava o
padrão — não foi refatorada nem ligada a `npm`; o passo 4 que este achado media deixou de existir.
A primeira condição (lane nova em `harness.ps1` cobrindo `internal/...`) continua não implementada
e é, portanto, fora do escopo desta baixa.

**D-52. `scripts/arch-gate.sh` passo 1 (`gofmt`) nunca pode passar num checkout Windows deste
repo — `core.autocrlf=true`, blobs git em LF, `.gitattributes` só fixa `*.sh`.** Isso produz ruído
permanente de ~635 ficheiros marcados só por causa de CRLF, escondendo os poucos ficheiros
genuinamente mal formatados dentro do mesmo `gofmt -l`.

Medido nesta sessão — total bruto:

```
$ cd apps/server_core && gofmt -l internal | wc -l
639
```

```
$ git config --get core.autocrlf
true
$ cat .gitattributes
*.sh text eol=lf
```

Separando os dois problemas (Item de Trabalho 2 desta onda): para cada um dos 639 ficheiros,
removi `\r` do conteúdo e comparei contra o mesmo conteúdo passado por `gofmt` — se as duas formas
LF batem, a única diferença era CRLF; se não batem, o ficheiro está genuinamente mal formatado
independente de fim-de-linha:

```
$ while IFS= read -r f; do
    orig_lf=$(tr -d '\r' < "$f"); fmt_lf=$(tr -d '\r' < "$f" | gofmt)
    [ "$orig_lf" != "$fmt_lf" ] && echo "$f"
  done < gofmt_list.txt | wc -l
22
```

**635 ficheiros são ruído puro de CRLF; 22 são genuinamente mal formatados** (lista completa no
relatório do Item de Trabalho 2 desta onda). Interseção desses 22 com os ficheiros que este ramo
tocou (`git diff --name-only 1b2ef2da..a9f92809`) é **vazia** — nenhum dos 22 pertence a este
ramo, então não havia nada para consertar aqui (só registar). Todos os 22 vivem sob
`internal/modules/` (proibido tocar por esta onda, constraint global 1) ou sob `internal/kernel/`
(não tocado por este ramo).

Condição de remoção: uma das (não implementadas aqui, fora do escopo desta onda) — (1) uma regra
`.gitattributes` que force `apps/server_core/**/*.go` a `eol=lf` no checkout, (2) `arch-gate.sh`
normalizar CRLF antes de invocar `gofmt -l` (ex. `gofmt -l <(tr -d '\r' < "$f")` por ficheiro), ou
(3) uma lane que corra onde blobs e worktree já concordam (ex. dentro do container Linux do dev
stack, que não sofre `core.autocrlf`).

**D-53. O lado de leitura do catálogo (`catalog.Reader()`) está morto ao nascer — entregue sem
consumidor e sem nenhum ponto do plano ou do design que o anuncie como capacidade pendente.**
`apps/server_core/internal/contexts/catalog/module.go:46` (`func (m *Module) Reader() port.Reader`),
`apps/server_core/internal/contexts/catalog/port/reader.go` (interface `Reader`, tipo `Summary`) e
`port/reader.go:24` (`Summary.DescriptionState`) só têm chamadores em
`apps/server_core/tests/integration/catalog_ingest_test.go:85,119` e
`apps/server_core/tests/integration/catalog_provenance_test.go:47` — nenhum ecrã, job, adapter de
outro contexto ou comando de operador chama `module.Reader()`.

Medido nesta sessão:

```
$ grep -rn "\.Reader()\|port\.Reader\b\|SummaryReader\b" --include="*.go" internal/contexts/catalog internal/composition cmd tests | grep -v "_test.go"
internal/contexts/catalog/internal/postgres/repository.go:27:// SummaryReader at the bottom of this file, which wraps it.
[... só comentários e a própria definição/wiring interna ...]
internal/contexts/catalog/module.go:20:	reader  port.Reader
internal/contexts/catalog/module.go:36:		reader:  postgres.NewSummaryReader(repo),
internal/contexts/catalog/module.go:46:func (m *Module) Reader() port.Reader { return m.reader }

$ grep -rln "module\.Reader()" --include="*.go" .
./tests/integration/catalog_ingest_test.go
./tests/integration/catalog_provenance_test.go
```

Distinto das duas dívidas já registadas sobre `catalog` (DSN de `mpc_app` que liga como
superuser, e o combinador de kernel) — esta lê no plano como capacidade ENTREGUE (Tarefa 6/7 do
plano `fecho` provisionaram proveniência honesta e um leitor de resumo), mas nada no repositório
fora dos próprios testes de integração jamais pergunta "o que é este produto?" pela porta que a
Tarefa 6 publicou. Não é um bug — o `Reader()` funciona exatamente como especificado — é entrega
sem consumo: o mesmo padrão que a ADR-023 nomeia (§ "três módulos nunca publicaram um
consumível"), aqui invertido (um módulo publicou um consumível que ninguém ainda consome).

Condição de remoção: a tela ou job que hoje pergunta "qual é este produto?" (hoje inexistente,
ou respondida por outro caminho fora de `catalog`) passar a chamar `catalog.New(pool).Reader()` —
ou, alternativamente, o plano que fecha essa lacuna nomeando explicitamente o consumidor previsto
e a data.

**D-54. `apps/server_core/tests/integration/catalog_ingest_composition_test.go` — o teste de
integração bandeira deste plano — depende de estado prévio no banco: assere `count(*) = 3`
absoluto sobre um tenant FIXO enquanto insere chaves únicas por execução, sem `t.Cleanup`.**
Linha 46 abre o pool com um tenant fixo (`testpostgres.OpenPool(t, "tenant_catalog_composition")`);
linha 54 carimba as chaves com `time.Now().UTC().UnixNano()` (únicas a cada corrida); linhas 68,
77, 86, 96, 105, 114 comparam contra números absolutos (`3`, `0`) em vez de "delta desde o início
do teste". Correr o teste duas vezes SEGUIDAS contra o MESMO banco (sem recriar o schema entre
as duas) faria a segunda corrida falhar: a terceira chamada a `RunCatalogIngest` veria as 3 linhas
da primeira corrida MAIS as 3 novas (chaves diferentes por causa do `UnixNano`), e `productCount
!= 3` dispararia em `productCount == 6`.

Isto NÃO é um falso-verde hoje: `scripts/harness/Postgres.psm1` provisiona um banco `mpc_test_<hex>`
efémero por corrida (confirmado — cada invocação de `harness:integration` cria um container/banco
novo), então o CI nunca acumula estado entre corridas e o teste passa legitimamente todas as vezes
que foi medido nesta sessão. O que falta é REPETIBILIDADE da prova, não correção do resultado: um
operador que corra este teste manualmente duas vezes contra um banco de desenvolvimento persistente
(`go test -tags=integration ./tests/integration -run TestCatalogIngestWritesRows -count=2` contra
`MPC_TEST_DATABASE_URL` apontando para um banco já usado) reproduz a quebra.

O padrão correto já existe no próprio pacote: `catalog_rls_role_test.go:38-45` regista
`t.Cleanup` que apaga as linhas que inseriu, tornando esse teste seguro para correr repetidamente
contra o mesmo banco. `catalog_ingest_composition_test.go` não segue o mesmo padrão apesar de ser
escrito depois (Tarefa 8, depois de RLS ter sido fechado na Tarefa 10 — ordem cronológica inversa
à ordem das tarefas, então não é "o padrão ainda não existia quando este teste foi escrito").

Por instrução explícita desta onda de correção: **este debt regista o problema e NÃO conserta o
teste** — mudar uma asserção que o plano já ratificou (`count(*) = 3` foi o critério de aceite
literal da Tarefa 8) não é decisão desta onda de correção de revisão.

Condição de remoção: `catalog_ingest_composition_test.go` ganhar um `t.Cleanup` simétrico ao de
`catalog_rls_role_test.go` (apagando por `product_id = ANY(...)` as chaves carimbadas que criou),
OU as asserções migrarem de `count(*)` absoluto para delta (`count(*) - baselineCount`) medido
antes do primeiro `RunCatalogIngest` — qualquer uma das duas torna o teste seguro para correr
repetidamente contra um banco persistente.

---

**D-55. `docs/architecture/decisions/023-module-protocol.md` — três afirmações em prosa ainda
declaram "35 violações" como se fosse a medição corrente, três medições depois.**

O cabeçalho do ADR-023 foi corrigido de `128` para `234` por medição (commit `82bd18ef`, onda de
correção do review final de ramo). Três frases no corpo do documento não foram alcançadas por essa
correção e continuam a afirmar `35` em presente do indicativo:

- `docs/architecture/decisions/023-module-protocol.md:82` — "The 35 measured violations are exactly
  the violations of this line"
- `docs/architecture/decisions/023-module-protocol.md:300` — "the 35 violations cannot be
  detected..."
- `docs/architecture/decisions/023-module-protocol.md:306` — "it treats 35 symptoms of 3 causes"

Medição corrente, reproduzida por dois revisores independentes nesta sessão:

```
module_boundary_arch_test.go:216: 234 violation(s)
```

`git blame` situa as três frases em `d9f46585` (2026-08-05), portanto são ANTERIORES à própria
correção `35 -> 128` e sobreviveram intactas a DUAS passagens de correção de contagem. Não foram
introduzidas pela onda `82bd18ef`, e o desdobramento por módulo ("60 of 128", "9 of those 35") está
explicitamente marcado como sendo da era 128 — estas três, não.

Por que é dívida e não erro de escrita: o cabeçalho do mesmo ficheiro afirma agora `234`. Um leitor
que chegue ao §82 lê `35` sem qualquer marca de datação e não tem como saber qual dos dois números
é o corrente. Um número de medição sem carimbo de quando foi medido é indistinguível de um número
corrente errado — é exatamente a classe que a correção `128 -> 234` existiu para fechar, deixada
por fechar em três sítios.

Por que NÃO foi corrigida nesta onda: as três frases vivem nas secções Decision e Alternatives
Considered, não no cabeçalho. Reescrever prosa de uma secção de decisão ratificada é emenda de ADR,
não correção de medição, e `ARCHITECTURE.md:5` exige tratamento explícito para isso. A onda de
correção tinha mandato para corrigir a MEDIÇÃO do cabeçalho, não para reescrever o argumento.

Condição de remoção: cada uma das três frases passar a carimbar a sua medição (`"the 35 violations
measured at <commit> ..."`) OU ser reescrita para citar a contagem corrente com a data e o comando
que a produziu, seguindo a convenção `**Amended DATE — título.**` que o próprio ficheiro já usa nas
linhas 62 e 84. Verificação: `grep -n "35" docs/architecture/decisions/023-module-protocol.md`
não devolver nenhuma ocorrência de contagem de violações sem carimbo de medição.

---

**D-56. `AGENTS.md:31-32` continua a imprimir `GOCACHE=.gocache` — o literal relativo que
HARNESS-PROFILE §2 proibiu em 2026-07-28, depois de ter produzido um verde vazio.**

`AGENTS.md` diz, em forma copiável e sem qualificação:

```text
Use `GOCACHE=.gocache` for Go tests.
```

`docs/HARNESS-PROFILE.md:36-39` obriga ao contrário:

```text
GOCACHE must resolve to an ABSOLUTE path on Windows/pwsh (D-14, M-01): relative
  export GOCACHE="$(pwd)/.gocache"   (bash)
  $env:GOCACHE = (Join-Path (Get-Location) '.gocache')   (pwsh)
```

E o log de emendas do próprio perfil (`docs/HARNESS-PROFILE.md:1030`) regista que a forma relativa
já foi retirada de três sítios do perfil precisamente por isto: *"um chip copiou-a e obteve um
EXIT 1 de 83 bytes com zero `=== RUN` — a sexta instância de verde vazio, e a única que a doutrina
entregou ela própria."*

**Como foi medida.** Encontrada em campo, não por varredura: o revisor automático da PR #14 citou
`AGENTS.md:L29-L33` como fundamento para exigir `GOCACHE=.gocache` num workflow de CI. O achado
está correto quanto ao binding existir; o remédio que propôs é o literal banido. Um revisor a ler
o `AGENTS.md` chega à conclusão errada com razão — a fonte que lhe foi dada di-lo.

**Por que é dívida e não erro de escrita.** `AGENTS.md` é carregado em TODO arranque de sessão via
`CLAUDE.md`, enquanto o `HARNESS-PROFILE.md` só é lido por quem segue o ponteiro. O ficheiro de
maior alcance imprime a forma proibida; o de menor alcance carrega a correção. A assimetria garante
que a versão errada é a mais lida.

**Por que NÃO foi corrigida aqui.** `AGENTS.md` é doutrina de hub (`AGENTS.md:4-13` declara a
ordem vinculativa) e um chip não emenda doutrina de hub por iniciativa própria. A correção pertence
ao hub e é de uma linha.

**Condição de remoção.** `AGENTS.md` deixar de conter uma forma relativa copiável — ou apontar para
`HARNESS-PROFILE.md §2` em vez de reimprimir a regra, que é a disposição que o perfil já tomou para
si próprio. Verificação: `grep -n 'GOCACHE=\.gocache' AGENTS.md` não devolver nada.

---

**D-57. `contracts/gate/guards.json` (o inventário de guards do issue #3) não é provado COMPLETO
— só provado não-vazio.** Achado do CodeRabbit na PR #25 (`ci(gate): guard inventory -- evidence
stops certifying itself (#3)`), aceito pelo operador como dívida em vez de trabalho dentro da PR.

`Test-GateGuardInventory` (`scripts/harness/Gate.psm1:454-491`) soma `go_tests` + `pwsh_files` +
`presence_only` do JSON e só falha quando essa soma é zero:

```
scripts/harness/Gate.psm1:490:  $total = $goEntries.Count + $pwshEntries.Count + $presenceEntries.Count
scripts/harness/Gate.psm1:491:  if ($total -eq 0) {
```

Isso prova que o inventário não está vazio; não prova que ele lista toda família de guard que a
harness de fato emite. Uma família nunca registrada, ou uma entrada apagada do JSON, desaparece
sem nenhum sinal — o total só cai, nunca zera.

Omissões medidas nesta sessão (verificadas contra o código, não copiadas do achado):

```
$ grep -n "throw 'HPG_\|Set-Primary 'HPG_\|Add-CleanupCode 'HPG_" scripts/harness/Postgres.psm1 | grep -o "HPG_[A-Z_]*" | sort -u
HPG_CONTAINER_REMOVE_FAILED
HPG_CONTAINER_START_FAILED
HPG_DATABASE_CREATE_FAILED
HPG_DATABASE_DROP_FAILED
HPG_DOCKER_MISSING
HPG_DOCKER_UNAVAILABLE
HPG_IMAGE_MISSING
HPG_MIGRATION_FAILED
HPG_MIGRATION_INVENTORY_INVALID
HPG_MIGRATION_NOT_IDEMPOTENT
HPG_PORT_UNAVAILABLE
HPG_READY_TIMEOUT
HPG_RESOURCE_CONFLICT
HPG_RESOURCE_LEAK
HPG_RUN_ID_INVALID
HPG_TEST_FAILED
HPG_TEST_VACUOUS
```

`contracts/gate/guards.json` só regista `integration-vacuous-run` (token
`guard_ran=integration-vacuous-run`, correspondendo a `HPG_TEST_VACUOUS`). `HPG_READY_TIMEOUT`,
`HPG_RESOURCE_CONFLICT` e `HPG_TEST_FAILED` — três códigos que o mesmo módulo emite — não têm
entrada nenhuma. (Nota de correção: o brief original citava só estes três; a varredura completa
acima mostra que a lista real de códigos não registados é maior — 16 dos 17 `HPG_*` do módulo
ficam fora do inventário, só `HPG_TEST_VACUOUS` está registado, não só os três citados.)

`scripts/harness/Policy.psm1` emite a família `RCFG_*`:

```
$ grep -n "New-PolicyIssue 'RCFG_" scripts/harness/Policy.psm1 | grep -o "RCFG_[A-Z_]*" | sort -u
RCFG_ALIAS_COLLISION
RCFG_ALIAS_UNDECLARED
RCFG_DYNAMIC_READER_UNBOUNDED
RCFG_LANE_VIOLATION
RCFG_READER_MISSING
RCFG_REGISTRY_READER_COVERAGE
RCFG_SECRET_CLASS_MISMATCH
RCFG_UNAPPROVED_READER
RCFG_UNDECLARED_READ
```

Nenhum `RCFG_*` aparece em `contracts/gate/guards.json` (confirmado — o arquivo só cobre `GOV_*`
via `gov-*` em `governance-drift.tests.ps1`, `HPG_TEST_VACUOUS`, `gate-counters` e
`guards-lane-counters`). Nove códigos de uma família inteira de política de runtime-config —
`RCFG_UNAPPROVED_READER`/`RCFG_READER_MISSING` são justamente os que fundamentam o baseline de
51 violações registado em B-9 — vivem fora do instrumento que devia prová-los guardados.

Contexto que corrobora a classe: durante esta mesma PR (#25/#3) descobriu-se que
`GOV_COMPOSITION_MISSING` estava listado no inventário com uma âncora satisfeita só por um
COMENTÁRIO solto — SEM fixture nenhuma por trás. Ou seja, uma entrada presente na lista também
mentia, não só as ausentes; a população precisa de medição tanto quanto cada entrada individual.
O checador por âncora foi substituído por tokens de execução por entrada no commit `4790a358`, e
a lógica de veredito da própria lane ganhou cobertura must-fail em `fbd34d29` — as duas correções
fecham "cada entrada mente", nenhuma fecha "a lista está incompleta".

Classe: **inventário hand-maintained que só falha vazio é população não medida** — mesma família
de B-1/B-6/B-7 (verde não-discriminante por falta de contagem contra uma fonte independente).

Conserto de classe — deliberadamente NÃO registrar entradas à mão para as omissões medidas acima:
registrar por mão é exatamente o padrão que dá o próximo drift. O conserto real é uma checagem de
completude que DERIVA a população esperada de guards de uma fonte autoritativa — as constantes de
código emitido nos módulos da harness (`HPG_*`, `RCFG_*`, etc.), do mesmo jeito que
`contracts/governance/invariants.json` já enumera os códigos `GOV_*` — e reprova a lane `guards`
quando existe um código sem entrada no inventário.

Condição de remoção: `Test-GateGuardInventory` (ou uma checagem irmã na mesma lane) passar a
extrair os códigos emitidos por cada módulo da harness (por regex/AST sobre `throw`/`Set-Primary`/
`New-PolicyIssue`/etc.) e comparar por conjunto contra `contracts/gate/guards.json`, falhando a
lane `guards` para qualquer código sem entrada correspondente. Verificação: injetar um código novo
num módulo da harness sem tocar `guards.json` e confirmar que a lane `guards` vira vermelha —
mesmo padrão de prova must-fail que `fbd34d29` já estabeleceu para o veredito da lane.
