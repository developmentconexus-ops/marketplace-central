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
