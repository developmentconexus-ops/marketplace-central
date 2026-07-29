# Análise global do HUB (rascunho — vira seção do documento de síntese)

Perspectiva: sessão hub MIS-006, autor de ~65 emendas do profile e do HARNESS-DEBTS.md.
Advertência de auto-referência: esta análise vem da mesma faculdade que produziu as falhas.
Sol xhigh e MNOS entram como contrapeso.

## P-1. A doutrina cresce por acreção, nunca por revisão — a harness não aplica a si mesma a própria regra
- Evidência: profile 1033 linhas, ~65 emendas datadas, ZERO consolidações. Só o §11 ganhou ~15
  subseções novas em 2 dias (2026-07-28). A third-round rule ("terceiro defeito da mesma forma
  para o point-fixing") existe desde D-121 para chips — e nunca disparou para o profile, que
  point-fixou a MESMA classe (vacuous green) SEIS vezes documentadas.
- Causa-raiz: não há dono nem cadência de refatoração da doutrina. Emendar é barato no momento
  (append), consolidar não tem gatilho. O custo aparece depois: cada chip/worker paga a leitura
  de 1033 linhas, e regra enterrada em emenda ~60 não é lida (evidência: chip copiou o
  `GOCACHE=.gocache` que o próprio arquivo proibia 3 linhas acima — a doutrina era longa demais
  para ser coerente consigo mesma).
- Conserto estrutural: (a) third-round rule vale para doutrina — 3ª emenda da mesma classe
  obriga consolidação em MECANISMO antes de nova prosa; (b) sessão periódica de compressão do
  profile com kill-list (que regras morrem porque viraram tooling).

## P-2. Instrumentação ausente vira prosa normativa — a família "vacuous green" inteira é UM defeito de runner
- Evidência: 6+ instâncias documentadas (skip silencioso, `-run` sem match, build tag não
  compilado, summary.txt byte-idêntico, regex ANSI-cega, `wc -m` retornando bytes). Cada uma
  gerou emenda de prosa: "count never tail", "prove que pode ficar vermelho", "print população
  antes do filtro", "strip escapes", "skips são resultado".
- Causa-raiz: as lanes devolvem TEXTO de ferramenta (go test, vitest) e cada leitor re-parseia
  na mão. Não existe camada de resultado estruturado. Toda regra da família é um parser humano
  codificado em doutrina.
- Conserto estrutural: runner único de lane que emite JSON assinado {cmd, cwd ABS, tip SHA,
  env consumido POR NOME, RUN/PASS/SKIP/FAIL contados, duração, exit}; verde sem contagem > 0 é
  inválido POR SCHEMA. Mata: B-1, B-2, metade do §11 vacuous-green, "count never tail",
  "name the tree it measured". O executing seat consome o JSON, não o texto.

## P-3. Assentos desenhados pela NEGAÇÃO de capability, remendados por protocolo
- Evidência: reviewer read-only não executa → terceiro assento (emenda 2026-07-28); worker
  sandbox não roda esbuild/vitest → "chip re-run é a verificação de record" (2026-07-17);
  sandbox nega travessia de árvore → cegueira packages/* (B-3); chip não sobe servidor → seam
  de empréstimo de arquivo do hub (F-5); Opus seat sem Write → verdito só na notificação →
  regras de custódia de 3 parágrafos.
- Causa-raiz: os papéis foram definidos por restrição de segurança/isolamento PRIMEIRO e os
  critérios que eles precisam descarregar vieram DEPOIS. Cada gap virou um protocolo de
  contorno com custo permanente de coordenação.
- Conserto estrutural: matriz papel → critérios que descarrega → capabilities exigidas,
  resolvida ANTES do despacho. Se o critério exige execução, o assento nasce com executor
  (ou o critério é roteado por desenho, não por descoberta em rodada 1). Sandbox de worker FE
  provisionado com leitura da raiz (config resolution) por padrão.

## P-4. Alegações viajam como prosa sem proveniência; a prosa apodrece e ninguém a corrige por mecanismo
- Evidência: A-2 (3 alegações falsas no brief M-06; A-17 do hub falso — "default na seam de
  load" num sistema fail-closed; card S8 com comando de lane quebrado); D-1 (rulings cruzando
  no correio 2×); A-3 (6 amendments manuais numa semana, cada um dependendo de disciplina).
- Causa-raiz: o protocolo transporta CONVICÇÃO, não MEDIÇÃO. Um brief não distingue "medi
  isto" de "acho isto". E não há correlação mecânica REQUEST→RULING→ACK, então ordem de
  entrega vira semântica.
- Conserto estrutural: claim de brief carrega proveniência obrigatória (comando + hash do
  output + SHA) ou é marcado UNMEASURED e não pode fundar critério; correlação por ID nas
  mensagens cross-session; ruling só ACKa depois do commit do amendment (mecanizar A-3).

## P-5. A harness não tem self-test — instrumentos nascem sem known-answer
- Evidência: pack-measure.sh shipou o defeito que existia para matar (`wc -m` = bytes);
  Stop hook acusou falso 2× (mediu worktree errado, e a metade "CLOSED claimed" também era
  falsa); dispatch-lint bloqueou spawn_task legítimo (D-3); sweep regex `[a-z_]` cega a
  acento/maiúscula reportou 1 violação onde havia 5; runtime roda hooks DRIFTED vs 0.4.0
  deployed (D-4).
- Causa-raiz: o código de produto tem gate (must-fail, dual gate, contagem); o código da
  HARNESS não passa por gate nenhum. Hooks e scripts entram em produção por edição direta.
- Conserto estrutural: todo instrumento (hook, sweep, runner, script de despacho) carrega
  known-answer test executado no deploy E no boot da sessão (o hook que não passa no próprio
  teste não dispara acusação — degrada para `unknown`). Versionamento runtime==fonte com
  verificação no boot do hub.

## P-6. Economia de rodadas invertida: o barato no despacho vira caro no gate
- Evidência: A-1 (S9 complex com worker low → reprova P4 → rodada extra); ~20 rodadas de gate
  gastas em prosa antes do BLOCKING=observável (CORTE-YAGNI: "17 rodadas produziram redação");
  packs de evidência 9.7k/11.5k/21.8k linhas para 605/1.8k/1.3k de código; 7 rodadas sobre
  tela INALCANÇÁVEL que um live-drive matou em minutos.
- Causa-raiz dupla: (i) alocação de effort por default econômico sem matriz classe→tier;
  (ii) loop de feedback perverso — instrumento não confiável → mais prosa de evidência →
  gate revisa papelada → mais rodadas → mais prosa. A prosa é o juro composto da
  instrumentação ausente (liga P-6 ao P-2).
- Conserto estrutural: matriz classe-da-fatia → tier obrigatória conferida pelo tooling de
  despacho (A-1+A-4); critério de gate nasce OBSERVÁVEL (o discriminador "o que um
  usuário/caller/linha faz diferente?" aplicado na AUTORIA do contrato, não na triagem do
  finding); orçamento de rodadas por unidade com stop-the-line automático no estouro.

## P-7. Conhecimento em 4 camadas com drift real e sem sincronização mecânica
- Evidência: memória do hub + profile + HARNESS-DEBTS + packs registram a mesma lição em
  redações diferentes; runtime DRIFTED (D-4); candidatos "upstream ao core" acumulam sem
  fluxo (B-1/B-2/B-4/C-1/C-2 todos marcados "candidato a core" e parados); regra de lane FE
  correta existia no profile e o brief S8 saiu com o comando quebrado mesmo assim.
- Causa-raiz: cada camada tem escritor diferente e nenhuma é derivada das outras; a promoção
  profile→core é manual e gated, então nunca acontece no fluxo.
- Conserto estrutural: uma fonte (repo mnfs-harness) + caches DERIVADOS por build; fila de
  promoção a core como artefato versionado com cadência; brief GERADO a partir do profile
  (comandos de lane vêm de um registry, não digitados por card).

## P-8. Hooks disparam por pattern-match sem nomear o que mediram — alarme falso treina surdez
- Evidência: Stop hook 2× falso (árvore errada + token em ORDENS lido como evento);
  dispatch-lint sem distinção chip/task-seed; "an alarm wrong twice trains its reader to skip
  the third" (já ratificado como texto, não como mecanismo).
- Causa-raiz: hooks fazem asserção sobre "o repo" resolvendo caminhos contra cwd volátil
  (worktree topology) e casam tokens sem contexto de gramática (CLOSED dentro de ordem).
- Conserto estrutural: contrato de hook — todo veredito automatizado imprime caminho ABS +
  SHA medido ou degrada para `unknown` (JÁ ratificado como regra; falta ENFORCEMENT no
  runtime do hook); eventos de protocolo (CLOSED/BLOCKED) parseados de campo estruturado da
  mensagem, nunca por substring.

## Meta-observação (a mais global)
Quase todos os padrões acima são instâncias de UMA lei: **o sistema converte falha de
mecanismo em obrigação de disciplina**. Cada falha vira uma regra que um humano/modelo deve
LEMBRAR de aplicar, e regras de disciplina falham do mesmo jeito que a falha original — por
atenção. A sessão de melhoria deveria medir cada regra do profile com a pergunta: "isto pode
virar schema/tooling/capability?" — se sim, a regra é dívida de implementação; se não, é
doutrina legítima. Estimativa honesta: >60% do profile é dívida de implementação disfarçada
de doutrina.
