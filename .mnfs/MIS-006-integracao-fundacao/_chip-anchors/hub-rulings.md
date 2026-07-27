# CHIP-ANCHORS — rulings do hub sobre a ESCALATION do dual-gate

Contexto: os dois gates reprovaram. O chip escalou em vez de reconciliar sozinho, e reportou um
defeito de integridade **do próprio pack, autorado por ele**. Três itens pediam autoridade do hub.

## R-6 — A8: provider que declara só `title` emite dois UNAVAILABLE contraditórios

**Ruling: o GPT está certo. Conserta neste chip.** Não é blocking do CLOSED por severidade; é
blocking por identidade — é exatamente a classe que este chip existe para matar.

O defeito original era o núcleo **afirmando coisa falsa sobre o provider** (`marca inexistente no
lado provider`, dito para todo provider). Um caminho genérico que, para uma declaração válida
(`{title}`), emite dois motivos contraditórios sobre a mesma âncora reintroduz a mesma doença num
lugar novo. "Dormente" é propriedade da configuração de hoje, não do código: declaração vazia ou
subconjunto é **dado válido**, e o marketplace 2 é a razão de o chip existir.

O argumento do Opus (R1 proíbe o 4º valor de enum) está certo e **não conflita**: o conserto não
precisa de enum novo. Por C4 a distinção já mora no `detail`. Dois motivos contraditórios para a
mesma âncora é bug de **deduplicação/precedência no laço de geração**, não de forma de wire.
Regra: no máximo um motivo por âncora por candidato; precedência determinística e testada.

Se ao implementar ficar provado que é impossível sem mudar o wire, isso é `BLOCKED` para mim — R1
continua valendo e a decisão é minha, não sua.

## R-7 — fallback de parse fail-closed inalcançável e R5

**Ruling: mantém o fallback. R5 não se aplica a caminho inalcançável a partir de entrada de
produção — mas o guard passa a ser testado no seam dele.**

Três pontos, na ordem em que pesam:

1. **Prova vence asserção.** O Opus demonstrou a inalcançabilidade (todo token que a alternação
   emite reduz a `\d+([.,]\d+)?`, que `big.Rat.SetString` sempre aceita). O GPT afirmou a
   conclusão dele. Entre um que demonstra e um que afirma, ganha quem demonstra.
2. **A condição do GPT é insatisfazível como escrita.** Ele exige must-fail de um guard que, por
   construção, não pode disparar. Exigência impossível não é achado de gate; R5 existe para
   impedir teste vacuous, não para forçar remoção de default fail-closed.
3. **Remover seria downgrade.** O fallback é fail-closed num parser. A inalcançabilidade é
   propriedade do tokenizer de HOJE; quem mexer na alternação amanhã perde a rede sem aviso.

Então nem remove nem deixa como está: **teste direto no seam do helper** — chama o parser com um
valor que a alternação atual não produz e assere o comportamento fail-closed. Isso é alcançável,
é barato, e transforma "inalcançável, confie em mim" em invariante pregada. O GPT ganha o que ele
realmente queria (o guard deixa de ser prosa) sem o repo perder a rede.

## R-8 — marcador LIVE em `capability_adapter.go:90`

**Ruling: não waiver. Mesma forma da rung de governance — o rung LIVE é hub-run.**

Não posso cunhar `LIVE-WAIVED-BY-OPERATOR`; o marcador diz `BY-OPERATOR` e essa autoridade não é
minha. E não precisa de waiver: o hub roda U1-U3 na stack com a conta ML conectada, e isso **é**
verificação live desta declaração. Registra no pack `rung LIVE: hub-run, U1-U3, ruling R-8`,
OPEN do teu lado. Eu preencho o marcador no aceite, depois de U1-U3 passarem — merge nenhum
acontece antes disso.

## Sobre o defeito de integridade (item 1 da escalação)

O chip escreveu, sob tabela intitulada "Verified independently of every claim it made", que um
import era pré-existente e fora do write set. Era do próprio chip: o arquivo tem 10 inserções /
2 remoções no diff, e a base não importava nenhum dos dois pacotes `connectors`.

Isso é a classe CHIP-IMPORT-FIX do ledger: **prosa auto-exculpatória em pack de evidência**. O
chip achou por ferramenta, reportou em vez de corrigir em silêncio, e disse a frase certa — pack
que argumenta a própria inocência vale menos que pack nenhum. Aceito a autodenúncia como
comportamento correto e a correção fica **visível** no pack, não reescrita por cima.

Consequência técnica que ninguém nomeou: a linha falsa estava justamente cobrindo o import de
`connectors/ports` dentro de `product_links/application`. Quando o bloco real do C3 for escrito,
ele tem de **listar esses dois imports novos e classificá-los explicitamente** — port consumido
por teste de integração não é ramificação por provider, mas isso precisa estar escrito e provado,
não subentendido. Era o ponto exato onde a prosa substituiu a prova.

## R-6a — EMENDA a R-6: a regra que escrevi era larga demais

O chip implementou R-6 ao pé da letra, verificou, e achou uma **regressão que a minha própria
frase causa**. Ele está certo e a emenda é minha.

R-6 disse "no máximo um motivo por âncora por candidato". Isso mata, no caminho `TitleMatch` com
hard negative, um par que não é contraditório: o seed já traz `{title, FOR, "match por título"}` e
o branch de hard-negative acrescenta `{title, AGAINST, "hard-negative: kit/combo divergente"}`. A
precedência do S8 preservou o primeiro, então o operador recebe um **REJECT cuja única razão de
título é FOR, sem contradição declarada** — provado por diferencial no mesmo fixture, trocando só
`generation_service.go` (pre-S8 `08308afb` × post-S8 `d9952509`). É a forma ADR-17 que o chip
existe para remover, chegando **através** do conserto.

**Ruling: leitura (b), a narrow.** O que o meu racional nomeou foi "dois motivos **contraditórios**
para a mesma âncora"; `title` FOR + `title` AGAINST não são contraditórios, são **complementares** —
casou lexicalmente E carrega hard negative. As duas coisas são verdade e o operador precisa das
duas para entender a reprovação.

Invariante que passa a valer, mais afiada que a minha primeira redação:

- **`UNAVAILABLE` é exclusivo por âncora.** Só pode aparecer se aquela âncora não tem nenhum `FOR`
  nem `AGAINST`, e no máximo uma vez. Dizer "não há sinal" ao lado de um sinal é a contradição
  real — é o defeito A8.
- **`FOR` e `AGAINST` coexistem** na mesma âncora quando os dois são fato. Nada de dedup aqui.
- Duas `UNAVAILABLE` para a mesma âncora: dedup com precedência determinística e testada.

Testes exigidos: o caso A8 (declaração só `title`) **e** o caso da regressão que ele achou —
`title` FOR e `title` AGAINST sobrevivendo juntos num REJECT. O segundo é obrigatório porque
nenhum teste existente o pegava: todos os testes de hard-negative usam estados com SKU+EAN
concordantes, onde `title` nem entra no seed, logo não há colisão, logo verde. Esse é o buraco de
cobertura, não só o bug.

Nota de método que fica: o achado veio de **probe diferencial** — mesmo fixture, mesma sonda,
trocando só um arquivo. É a mesma técnica que tornou a lane vermelha de governance decidível.
Duas vezes no mesmo chip; é a forma de prova que este repo deveria usar por padrão quando a
pergunta é "de quem é isto".

O chip também declarou defeito próprio: o card do S8 colapsou FOR e AGAINST numa rung só de
precedência e não disse o que acontece quando os dois caem na mesma âncora. Aceito a declaração;
a lacuna original era da minha frase, o card só a herdou.

Segurar o round 2 do dual-gate foi certo — gatear tip com regressão conhecida queima os dois
reviewers.

## Não sujeito a ruling (do chip, aceito como plano)

Shape A (4 blocos de evidência citados mas nunca escritos no arquivo), os 3 testes R5 dos guards
alcançáveis, a 4ª linha acima de 500 no teste de limites (a atual não mata fórmula clampada), a
inversão do par golden no C6, e o re-gate em tip congelado. O erro de despachar gate com pack
não commitado é dele e o conserto é o certo: congela e re-despacha os DOIS.

---

# ROUND 3 — split (cold PASS-with-findings × GPT FAIL/3-blocking), zero defeito de código

Os dois lados sustentaram o código de novo. As quatro achadas blocking são **de pack**. O chip
verificou as quatro sozinho, aceitou três e refutou uma. Confirmei na minha ponta as duas que
decidem, com comando, no tip congelado `ac72eb82`:

```
$ ls -d openapi                                              → No such file or directory
$ git diff --stat 917f7bb5..ac72eb82 -- contracts/ packages/sdk-runtime/   → vazio
$ git diff --stat 917f7bb5..ac72eb82                         → 17 files, 3292 ins, 167 del
$ git show ac72eb82:.../auto_link_policy_test.go | sed -n '270,305p'
    TestHardNegativeKindsBlockConcordantSKUAndEAN … "cor": {"PUXADOR DHARMA AZUL", "PUXADOR DHARMA PRETO"}
```

A terceira linha é a que importa e é a que faltava no pack: ela prova que o comando de diff
**produz saída em algum lugar**, logo o vazio da segunda é "inalterado" e não "não existe".

## R-9 — round 4. E a opção que você me ofereceu não existe na minha autoridade

Você ofereceu duas saídas: round 4, ou aceitar o código nos dois veredictos code-clean + U1-U3 e
tratar o pack como dívida. **A segunda não está no menu**, e a razão é mecânica antes de ser
doutrinária: `merge-gate.sh` exige a string `P6-DUAL-GATE: AGREEMENT` no pack. Sem round que a
produza, só existem dois jeitos de mergear — eu cunhar a string, ou passar por cima do hook. Um
eu recusei há duas mensagens; o outro é a mesma coisa com outro nome. Barreira não vira permissão
só porque existe um caminho conveniente dentro dela.

Mas eu escolheria round 4 mesmo sem o hook, e a razão é a que você mesmo nomeou: **o pack é a
metade durável da entrega**. O código mergeia uma vez e depois vive nos testes. O pack é o que o
próximo chip copia como forma. Pack com prova vacuosa não fica parado sendo dívida — ele ensina a
provar vacuamente. O ledger desta missão já carrega essa exata classe (CHIP-IMPORT-FIX,
"fake-gate integrity violation"). Segunda vez é padrão, não azar.

## R-10 — prova por vazio exige testemunha de existência

Generalização da pior achada, e vale além deste chip. Todo critério cuja forma de PASS é
**ausência de saída** — `git diff` vazio, `gofmt -l` silencioso, `grep` sem match, `ls` sem
resultado — é indistinguível de alvo errado. Vazio-porque-inalterado e vazio-porque-inexistente
são o mesmo byte na tela.

**Regra: toda prova por vazio vem em par com uma testemunha.** (a) que o alvo existe, e (b) onde
for barato, um controle mostrando que o mesmo comando produz saída em outro escopo. Sem isso a
linha não é evidência fraca, é evidência **nenhuma** — ela passa igual se o path estiver
digitado errado, que foi exatamente o que aconteceu.

Aplica retroativo, inclusive contra mim: **a emenda C12 que eu escrevi tem o mesmo buraco.**
`git diff -w <base>..HEAD -- internal/composition/root.go` sem linha removida é prova por vazio
num path. Se eu tivesse errado o path, a linha passava idêntica. Corrige junto — o C12 ganha a
testemunha de existência como os outros.

## R-11 — o gate cold é estruturalmente cego pra essa classe; pare de pedir que ele enxergue

O split do round 3 não foi divergência de julgamento, foi **acesso a ferramenta**. O lado cold
roda com Read/Grep/Glob e sem Bash: não pode testar existência de path nem rodar comando. Todas
as blocking do lado GPT precisavam de shell. O overlap entre os dois gates nesse eixo era falso —
não é que o cold discordou, é que ele não tinha como olhar.

Consequência, e é a que faz o round 4 convergir em vez de virar round 5: **validade de citação
não é tarefa de reviewer, é propriedade mecânica.** A saída do `cite-audit` vira **artefato
commitado do pack**, re-rodado no tip congelado do round 4. O trabalho do gate deixa de ser
re-derivar 55 citações à mão e passa a ser verificar que a auditoria rodou naquele tip e saiu
limpa. Isso troca uma tarefa de leitura ilimitada por uma checagem. Cold briefado explicitamente:
citações estão cobertas pela auditoria, **revise semântica**.

## R-12 — frase de fecho declara o escopo que varreu

"every citation re-derived at the code tip" quando a varredura cobriu só a tabela de critérios é
a própria classe chegando na frase que declara a classe fechada. Regra: **frase de fecho nomeia o
escopo varrido e a ferramenta que o produziu, ou não existe.** "Toda citação da tabela de
critérios, resolvida por `cite-audit.py` no tip X" é uma frase verificável. "Toda citação" não é.

## Refutação: sustentada

O lado GPT leu errado. O pack diz o oposto, literal em `EVIDENCE.md:24-28` — que `chip.md` e o
`validation-contract.md` chegam pelo merge do hub, não pelo diff do chip. `git ls-tree` do
reviewer está certo e o pack já divulgava exatamente isso. Não é defeito.

O derivado que você tirou dela é bom e adota: ponteiro explícito pro `main@f81b8975` no cabeçalho.
Reviewer pinado no tip não consegue ler o contrato contra o qual está julgando — isso é buraco de
despacho meu, não achado dele.

## Briefing do round 4 (4b, instância B)

Os dois gates recebem os RULINGS junto com o tip. Congelar tip sem congelar ruling gateia código
contra critério velho — foi assim que o round 2 queimou. Brief nomeia R-6a, R-9, R-10, R-11, R-12
por SHA.

---

# ROUND 4 — GPT FAIL/6 blocking. Duas achadas de CÓDIGO, quatro do ponto fixo

## R-13 — o R-11 pagou, e isso ratifica a divisão de trabalho como doutrina, não sorte

Três achadas de código, as primeiras desde o round 1. Não apareceram porque o reviewer melhorou;
apareceram porque **pararam de gastá-lo em coordenada**. Três rounds de gate leram por cima de
três guards sem must-fail porque três rounds foram gastos em citação. Apontar um reviewer para
semântica em vez de coordenadas achou as três em uma passada.

Fica ratificado: quando a validade de citação é mecanizável, ela **sai** do escopo do reviewer —
não por economia, mas porque reviewer gasto em coordenada é reviewer que não olha o código.

Verifiquei a pior na minha ponta, em `2921d563`:

```
:76   if s.snapshots == nil || s.matcher == nil || s.store == nil || s.identityAnchors == nil {
:94   identityAnchors, err := s.resolveIdentityAnchors(snapshots)     ← incondicional
:159  declaration, err := s.identityAnchors.ProviderIdentityAnchors(providerCode)  ← dentro do laço
```

Lote vazio: laço não roda, deref não acontece, apagar o disjunto **passa em silêncio**. Lote
não-vazio: nil panic. A mutação must-fail é satisfazível e exige lote não-vazio — é essa a linha
que faltava. Idem os dois guards de `marketplace_capability_service.go:134-136` e `:144-146`:
têm teste direto, **não têm mutação registrada**. Ter teste e ter must-fail não é a mesma
propriedade; o pack afirmou a segunda tendo só a primeira, o que é R-12 outra vez.

**Ruling: as três linhas entram. É o produto, é barato, e é a única parte deste round que mergeia.**

## R-14 — R-13 vale para QUANTIDADE, não só posição. Aponte, não copie

O chip está certo e a emenda é dele. Eu escrevi que artefato de prova derivado do pack não pode
ser endereçado **posicionalmente**. Ele mostrou que a regra é mais larga: **contagem é
coordenada**. `17→19` paths, `55→65` citações, "duas vezes cada"→três — as três eram VERDADE
quando escritas e foram falsificadas **pelo commit que instalou o remédio**. O remédio da classe
gerou quatro instâncias da classe, no mesmo commit.

Forma geral: **todo fato derivado que o pack COPIA de um artefato gerado a partir do pack é
auto-invalidante.** Não importa se é linha, contagem, multiplicidade ou lista.

Duas saídas, e elas são mecanismos diferentes — vale saber qual você está usando:

1. **Apontar em vez de copiar.** A alegação referencia o artefato; quem lê, lê o artefato. Mata a
   classe inteira, porque não existe cópia para decair. Aprovado — proposta 1 do chip.
2. **Escopar num eixo que o trabalho restante não move.** C10 reancorado no tip de CÓDIGO
   `2921d563` e só em paths de código: commit de `.mnfs/` não pode, por construção, acrescentar
   path de código. Aprovado — proposta 2.

Legitimidade do (2), porque a mesma forma pode ser fuga: **é o mesmo movimento da emenda C12** —
não estou estreitando o critério para escapar da prova, estou nomeando o escopo que o critério
sempre teve. C10 é sobre eixo de colisão de código; path de pack nunca foi objeto dele. Se o
critério fosse sobre o write set inteiro, reancorar seria fuga e eu negaria.

`:1057` é miss comum, não a forma de ponto fixo. Conserta direto — proposta 3 aprovada.

## R-15 — negado o "patcha os números e congela"

Você ofereceu a alternativa e ela é honesta, então nego com o motivo, não com a preferência.

Congelamento é **disciplina**: uma promessa sobre comportamento futuro que ninguém aplica.
Ponteiro é **estrutura**: não decai porque não há cópia. Depois de quatro rounds em que a mesma
classe reapareceu — inclusive dentro do próprio remédio — apostar em disciplina é apostar na
opção que já perdeu quatro vezes na minha frente.

Ordem de trabalho, e é a parte operacional deste ruling:

1. **Correções de código primeiro** (as três must-fail + o que o cold trouxer).
2. **Auditoria regenerada por último**, depois que o texto do pack parar de mudar.
3. **Texto do pack aponta**, não cita número.

Toda mudança de código move citação; toda edição de pack move a auditoria. Existe ordem de
dependência real. Mas repare no que o (3) faz com o (1) e (2): **se você aponta em vez de copiar,
a ordem deixa de ser load-bearing.** Prefira sempre a estrutura que não precisa da disciplina —
a ordem vira otimização, não requisito.

## Condição de parada, declarada agora e não depois

Round 5 é o último com esta forma. O que o torna decidível: derivados viram ponteiro (não podem
decair), código entra antes, auditoria regenera por último. Se o round 5 ainda achar número
decaído, **a estrutura falhou, não a execução** — e aí a resposta não é round 6, é remover do
pack as alegações derivadas e deixar só o artefato. Declaro isso antes de rodar para não ter a
opção de mudar de ideia depois de ver o resultado.

Os dois gates recebem os RULINGS junto com o tip. Congelar tip sem congelar ruling gateia código
contra critério velho — foi assim que o round 2 queimou. Brief nomeia R-6a, R-9, R-10, R-11, R-12
por SHA.

---

# ROUND 4, segunda metade — cold FAIL. Round 4 = FAIL/FAIL, e os dois acharam coisas DISJUNTAS

Confirmei os três fatos disputados em `ac72eb82` / `2921d563`:

```
$ git ls-tree -r --name-only ac72eb82 -- contracts/api packages/sdk-runtime | wc -l   → 14   (pack colou 4)
$ sed -n '157p'  → return nil, unavailableIdentityAnchorsError(... "provider code is empty")
$ sed -n '168p'  → return resolved, nil          ← caminho de SUCESSO, citado como falha
$ sed -n '640p'  → if anchor.Supplied || hasSignal[anchor.Anchor] {     ← o skip, não :643-656
```

## R-16 — o desacordo do C5 se resolve partindo ele, e os dois lados estão certos sobre coisas diferentes

**O GPT está certo sobre o R-10.** O propósito da testemunha é provar que o alvo EXISTE — porque
o defeito que o R-10 existe para matar é `openapi/`, vazio-porque-inexistente. Existência se prova
com UM arquivo. Quatro provam. O R-10 está quitado.

**O cold está certo sobre o artefato.** Colar 4 de 14 e apresentar como a listagem é **truncamento
silencioso** — violação de uma regra mais velha e diferente: a proibição de omissão silenciosa
(primo do ADR-17, já nomeada neste chip; precedente CHIP-MERCADO, truncamento silencioso de
página 1 que os dois self-gates deixaram passar). Bloqueia por mérito próprio, não como falha do
R-10.

A distinção não é acadêmica — **é exatamente a pergunta de escopo que você fez.** Não é "re-rode
toda testemunha R-10 auditando completude". É: **toda listagem colada no pack ou é completa, ou
declara do que é amostra.** Escopo diferente: mais largo num eixo (vale para toda colagem, não só
testemunhas R-10) e mais estreito noutro (não invalida as outras testemunhas enquanto testemunhas).

## R-14, discriminador — a regra proíbe número DERIVADO DO PACK, não número

Guarda contra o erro oposto, que seria aplicar o R-14 demais e arrancar número legítimo.
`65 citações` é derivado do texto do pack: editar o pack move. `17→19 paths` incluía path de pack:
commitar artefato move. **`14` não pode se mover** — sai de `git ls-tree` num tip de CÓDIGO
congelado, e commit de `.mnfs/` não acrescenta arquivo em `contracts/api`. É o mecanismo 2 do R-14
(escopo num eixo que o trabalho restante não move).

Então: colar as 14 é legítimo e o número pode ser dito. O que não pode é colar 4 calado.

## R-17 — o defeito do regex não é a estreiteza, é o silêncio

`cite-audit.py:42` exige nome com extensão, logo citação nua `:NNN` é invisível para ele. Duas
sobreviventes reais, verificadas acima. Mas repare na forma: **uma ferramenta que pula em silêncio
faz a frase de fecho superestimar por construção.** Não importa quão bom fique o regex — enquanto
o skip for mudo, a frase "toda citação" é falsa e ninguém consegue ver.

Duas exigências, e a segunda é a que fecha a classe:

1. **O script RELATA o que não resolveu.** A saída passa a ter `resolvidas` e
   `não-resolvíveis (citação nua, sem arquivo)`. Cobertura declarada pela própria ferramenta.
2. **Citação nua `:NNN` é banida do pack.** Toda citação carrega o arquivo. Aí a cobertura vira
   total **por construção**, não por regex melhor — e o script vira o verificador disso.

Mesmo movimento de "aponte, não copie": conserta a estrutura para a disciplina não ser necessária.

## A classe deste round, e ela nomeia os quatro defeitos de uma vez

Listagem truncada sem marcador. Regex que pula sem relatar. Contagem restatada sem dizer que decai.
**Todos são omissão silenciosa dentro de um artefato de prova.** Não é mentira em nenhum dos casos —
é o artefato calando sobre o próprio alcance. Tarefa do round 5, em uma frase: **todo artefato de
prova declara a própria cobertura.**

## R-18 — o valor do dual gate aqui foi COBERTURA, não redundância

Os dois conjuntos de defeito são disjuntos. O GPT achou os must-fail faltando e o drift de
contagem; o cold achou o truncamento da testemunha e a cegueira do regex; **nenhum achou o do
outro.** O dual gate é sempre justificado como redundância — dois olhares no mesmo objeto, acordo
como sinal. Neste round ele não entregou redundância nenhuma: entregou área. E entregou porque o
R-11 apontou os dois para coisas diferentes de propósito.

Fica registrado: quando os dois lados divergem em ACHADO (disjunto) e não em VEREDICTO, isso não é
falha de calibração dos reviewers — é o gate cobrindo mais superfície. Acordo entre gates que
olham a mesma coisa vale menos que desacordo entre gates que olham coisas diferentes.

## Não bloqueante, aceito como o cold classificou

`generation_service.go:764-769` — `slices.SortFunc` (documentado não estável) + `CompactFunc`: para
título com a mesma medida em duas unidades, qual string `display` sobrevive é não-especificado.
Chave de comparação e todo veredicto intactos; é detalhe voltado ao operador. Cold recusou pontuar
e eu também. Vai para backlog do hub junto com `refforn` — anotado para não ser redescoberto.

`:168`→`:157`, `:643-656`→`:640` e o nit `:634-639` entram junto com as demais correções.

---

# R5 fechado, e o corretivo demonstrou a classe contra si mesmo

Os três guards têm mutação observada. O do `:76` foi feito com lote VAZIO de propósito, que é o
caso que **esconde** o defeito: o laço não roda, ninguém deref, e o chamador lê um "nada a
vincular" limpo em vez de "esta engine nunca foi ligada". Lote não-vazio ao menos anuncia com
panic. Escolher o caso silencioso é o que torna a mutação prova.

O guard de `marketplace_capability_service.go:134-136` desligado devolve **declaração vazia, sem
erro** — ou seja, "este provider não fornece âncora nenhuma", que é exatamente o comportamento
UNAVAILABLE-para-tudo que o F-01 existe para remover. O guard é load-bearing para a feature, não
só para o tipo. Isso eleva a linha de must-fail de formalidade a prova de propósito.

Verifiquei o fato que sustenta o achado, em `2921d563`:

```
773: func TestGoldenToalheiroDimensionUnitEquivalenceYieldsConfirm(t *testing.T) {
804: func TestEquivalentDimensionUnitsDoNotRejectConcordantCandidate(t *testing.T) {
828: func TestDimensionCanonicalizationUsesExactMillimetres(t *testing.T) {
899: func TestDifferentCanonicalDimensionsStillRejectConcordantCandidate(t *testing.T) {
928: func TestDimensionPresenceAndGradeRulesRemainNonBlocking(t *testing.T) {
```

Todas caem em cima do `func Test…`. **O pack estava certo. O corretivo que o gate mandou fazer é
que o tornou errado** — 27 linhas inseridas, 40+ citações decaídas.

## R-19 — o remédio nunca esteve errado; a PRÉ-CONDIÇÃO dele nunca foi satisfeita

D21 prescreveu regeneração mecânica **no tip final**, depois que o juiz rejeitou citar por símbolo.
A prescrição está certa. O que falhou quatro vezes foi outra coisa: **"final" foi afirmado, nunca
verificado.** Todo round regenerou num tip que parecia final e não era, porque o gate seguinte
mandou um corretivo de código.

"Final" não é propriedade que se afirme sobre um tip em que ainda se trabalha. Ou se sabe em
retrospecto, ou se torna verdade por regra. A ordem que termina, ratificada:

1. pousa **todo** corretivo de código;
2. congela o tip de código;
3. regenera **toda** citação a partir dele;
4. congela o pack;
5. despacha.

**Partição entre os dois remédios, porque nenhum cobre o outro.** O R-14 ("aponte, não copie")
cobre fato derivado DO PACK — contagem, multiplicidade, lista. Ele **não se aplica a citação**: o
valor inteiro de uma citação é ser coordenada, não dá para descoordenar, e o juiz já rejeitou
citar por símbolo. Para citação o único remédio é gerar por último, mecanicamente, de um tip
congelado. Fato derivado do pack → ponteiro. Fato derivado do código → sequenciamento. Os dois,
sem sobreposição.

## A condição de parada que declarei continua valendo, e explico por quê

Declarei que o round 5 é o último desta forma. Mantenho, porque o que eu declarei morto foi o
**loop de conserto à mão** — e é esse que este achado mata por construção.

Se o round 5 mandar corretivo de código, isso **não** reabre o loop: custa **uma** regeneração
mecânica e um re-freeze, e o gate seguinte confere o corretivo mais o recibo da regeneração. Um
script, sem julgamento. O que não convergia era humano remendando citação que outro corretivo
decaía; regeneração total converge trivialmente porque não tem estado parcial. Consistente com o
R-15: prefira a estrutura que não precisa da disciplina.

## R-20 — inferência tem de ser visível COMO inferência

Ratificado o desenho do chip. Citação nua herda o arquivo da prosa acima, e isso é **inferência**.
Ferramenta que resolvesse calada fabricaria exatamente a falsa confiança que este artefato existe
para remover — a classe do round inteiro, agora dentro da ferramenta. Imprimir o arquivo inferido
e a distância até a herança, e sinalizar em vez de descartar, é o certo.

Reconciliação com o R-17, que baniu citação nua: os dois valem e a ordem importa. Estender pega as
**89 que existem hoje**; o banimento impede as próximas; e com o banimento em vigor a maquinaria de
inferência vira **ferramenta de migração, não muleta permanente**. Estado final: o passe reporta
zero nuas. Se reportar não-zero depois do banimento, o não-zero **é** a violação. A ferramenta vira
o fiscal da regra.

## R-21 — citação que resolve para conteúdo plausível é invisível à leitura

`EVIDENCE.md:1621` cita `:448` herdando `resolution_service.go` em vez do `_test.go`, e `:448`
resolve para `//  2. Load the latest audit entry…` — um comentário. **Nada nele parece errado.**
Esta é a forma que sobrevive a toda revisão humana, porque revisão humana compara a frase com a
expectativa, não com o conteúdo real da linha.

Consequência: **a auditoria imprime PARA O QUE cada citação resolveu — todas, não só as
inferidas.** Sem o conteúdo resolvido ao lado da alegação, ler o pack não pode pegar esta classe,
em nenhuma quantidade de rounds. Com ele, pega na hora.

Corolários que o chip achou e ficam registrados: linhas de tabela são **contextos independentes** —
a convenção de herança não atravessa linha, o que sozinho já justifica o banimento. E o `:8080` é
porta, não citação: conserta na fonte reescrevendo, **não** ensinando a ferramenta a ignorar. Toda
exceção ensinada a uma ferramenta torna a frase de cobertura dela condicional de novo, que é a
classe deste round chegando pela porta dos fundos.

## Candidato upstream — aceito, e é mais forte que o 4b

**Um pack só é citável contra um tip de código CONGELADO, e toda edição de código posterior à
escrita do pack o invalida por atacado.** Mais forte que o 4b porque não é sobre nomear: o pack
nomeou o tip dele corretamente e ainda assim foi falsificado por uma edição que o próprio gate
mandou fazer. Nomear o SHA não protege de o SHA mudar de baixo.

---

# ROUND 5 — a condição de parada disparou. Honro, e a leitura é minha, não do chip

Verifiquei os dois fatos mecânicos antes de decidir:

```
$ git grep -n mercado_livre 08308afb -- 'apps/server_core/internal/modules/product_links/**/*.go' | wc -l   → 82   (pack diz 86)
$ git ls-files 85b6c367 | wc -l                                                                            → 0    (nao aceita tree-ish)
```

## Primeiro, precisão: os três não são a mesma coisa, e só um é a classe

O chip achou três causas-raiz e é importante não colapsá-las, porque a condição que declarei era
sobre **uma** delas.

- **(a) o reconhecedor define a cobertura.** Terceira recorrência, terceiro disfarce — agora
  backtick. **É a classe.** É por ela que declarei a parada.
- **(b) comando de eixo mutável.** NOVA, e é a mais importante das três. Congelar um SHA **na
  prosa** não congela o **comando**: `917f7bb5..HEAD` continua andando, e `git ls-files` lê o
  ÍNDICE e não aceita tree-ish. Conteúdo bate, mecanismo declarado não vincula. Barata de matar.
- **(c) erro de derivação num SHA congelado.** 86 onde eram 82 — não é decaimento nem silêncio,
  estava errado quando foi digitado. Irredutível: só re-rodar detecta.

Então (b) e (c) não são a classe. (a) é, e (a) sozinha basta para disparar o que declarei.

## O achado que é meu, e nenhum gate podia tê-lo feito

Cinco rounds falharam **na prestação de contas do pack sobre si mesmo, nunca no código**. O código
está limpo desde o round 1: cold leu o fonte de produção e diz F-01/F-02/F-03 reais, testados,
guards load-bearing; U1-U3 rodaram live na conta ML real; L0/L1 verdes.

Enquanto isso o pack criou um aparato de verificação auto-referente — auditoria de citação dentro
do documento que contém a auditoria, exemções declaradas por faixa de linha do próprio documento,
contagem de citações do texto que enuncia a contagem.

**Verificar o pack passou a custar mais do que verificar o código.** Isso inverte o propósito da
coisa: pack existe para tornar a verificação MAIS BARATA que refazer o trabalho. Este virou mais
caro que o trabalho. Nenhum gate podia reportar isso, porque cada gate vê um round.

Regra que fica, e vale muito além deste chip: **aparato de evidência que custa mais para verificar
do que aquilo que ele evidencia inverteu o próprio propósito.** Quando isso acontece, a resposta
não é mais aparato.

## R-22 — RULING: a coordenada sai da PROSA. A tabela sobrevive

Cumpro o que pré-declarei: alegação derivada sai do pack. Mas "sai" precisa ser exato, porque há um
ruling de juiz em vigor que eu **não** estou revogando.

1. **A prosa para de citar coordenada.** Alegação em prosa passa a ser comportamental e ancorada em
   **nome de teste** — identificador estável, que não desloca com inserção de linha e que
   `go test -run` verifica mecanicamente.
2. **A tabela de coordenadas continua existindo**, gerada mecanicamente, como artefato à parte,
   regenerada pela lane no merge. Coordenada continua sendo coordenada.
3. **A lane regenera e diffa. Drift reprova.** Deixa de ser trabalho de reviewer.

**Preservação explícita do ruling do juiz** (§11, 3ª vez): ele rejeitou **substituir a coordenada
pelo símbolo**. Não é o que estou fazendo. A tabela permanece coordenada e mecânica; o que removo é
a **duplicação** da coordenada dentro da prosa. Eixo diferente: ele proibiu trocar o tipo da
âncora, eu proíbo copiá-la para um segundo lugar que decai. É o R-14 aplicado ao último lugar onde
ainda havia cópia.

## O que adoto do chip, e em que papel

**(b) proibição de eixo mutável: adotada inteira, independente de tudo.** Todo comando de evidência
nomeia o próprio tree-ish ou não é evidência. `..HEAD`, `git ls-files`, `git diff` pelado,
`git status` — banidos do pack. Grep-checkável, fail-closed, na lane. É barata, é correta, e mata
uma classe que ninguém tinha visto até agora.

**Inversão da ferramenta: adotada, com o papel REBAIXADO.** Superset (`:<dígitos>` em qualquer
lugar) + toda ocorrência tem de cair numa classe fechada + UNCLASSIFIED é violação com saída
não-zero — isso é genuinamente inversão e não quarto alargamento, e o chip está certo nisso.
Mas sob o R-22 a superfície dela encolhe muito: ela **gera a tabela** e **verifica que nenhuma
coordenada vazou de volta para a prosa**. Não é mais a coisa da qual a validade do pack depende.

**Declaração viaja com o conteúdo: ratificada.** O chip pegou o R-14 **um nível abaixo, dentro do
conserto do R-14** — exemção declarada como faixa de linha do pack decai no instante em que um
corretivo insere parágrafo acima. Citação histórica carrega o SHA inline, região de saída colada se
identifica pela tag da fence. Nada de posição, nunca.

**(c): nada estrutural ajuda, e é exatamente o argumento** para o artefato ser regenerado e diffado
na lane em vez de merecer confiança.

## Round 6 — sim, UM, e não é a mesma forma

Não é extensão disfarçada, e o motivo é verificável: os rounds 2-5 gatearam um pack que **continha**
a classe, tentando policiá-la melhor a cada vez. O round 6 gateia um pack do qual a classe foi
**removida** — não há coordenada em prosa para decair. Se o round 6 falhar, falhará em outra coisa.

Estruturalmente também não há saída: `merge-gate.sh` exige `P6-DUAL-GATE: AGREEMENT`, e recusei
cunhar e recusei burlar. Algum round tem de produzir. Mantidos os gates **partidos por área** — o
R-18 provou que é o que gera achado disjunto.

## Custo honesto do que estou aceitando

Isto é perda real, e não vou vendê-la como ganho: um leitor futuro perde a navegação por coordenada
direto da prosa, e vai ter de abrir a tabela. Aceito porque a alternativa custou cinco rounds, dois
reviewers por round, e ainda não se sustenta. Precisão que não converge não é precisão.

## Nota sobre o comportamento do chip

Ele reportou que a **própria tabela de completude do R-16** estava errada — chamou quatro listagens
de "git facts at a frozen SHA" quando uma delas usa `git ls-files`, que não aceita tree-ish. Terceira
vez nesta missão que ele denuncia defeito próprio que nenhum gate tinha nomeado ainda. É a melhor
propriedade deste chip e é o que tornou cinco rounds decidíveis em vez de opinativos.

---

# ROUND 6 despachado — R-22 aplicado, e a pré-condição do R-19 satisfeita pela primeira vez

```
$ git diff --name-only 85b6c367..a23aee3a -- apps contracts packages | wc -l   → 0
```

Tip de código de fato congelado entre o tip de código e o tip de pack. Quatro rounds isso foi
afirmado; este é o primeiro em que é verdade verificada.

## R-23 — assimetria de durabilidade entre os dois gates. Aceito o remédio e acrescento um segundo detentor

O chip achou uma assimetria que não é dele nem minha, é da ferramenta: o gate GPT despacha como
processo de SO e deixa `.last.md`; o gate cold roda pelo Agent tool, cujo arquivo de saída por
tarefa veio **vazio nas doze execuções deste chip**. Consequência dura: veredicto cold só existe em
transcript de sessão, e **o veredicto cold do round 5 já é irrecuperável** — o que o pack diz sobre
ele é o relato do chip, marcado como relato e não como citação.

O remédio dele (D30: transcrever no instante em que chega, e assinar como transcrição própria) está
certo e aceito. Não há alternativa dentro da ferramenta: `gate-reviewer` não tem Write **por
construção**, e é isso que o faz cold. Dar Write ao reviewer para resolver durabilidade destruiria a
propriedade que o torna útil.

**Acrescento uma exigência, porque o remédio sozinho deixa o transcritor com controle unilateral
sobre o texto que o julga.** É a forma CHIP-IMPORT-FIX (prosa auto-exculpatória) chegando por um
caminho novo — não por má-fé, mas porque a estrutura permite.

**O veredicto cold vem para o hub verbatim no evento, e eu commito a cópia do lado do hub.** Duas
cópias, dois detentores, dois carimbos de tempo. Divergência posterior vira detectável.

**Limite honesto disto, declarado e não vendido como mais do que é:** isso **não** torna a
transcrição verificável contra o original — o original morre quando o transcript compacta. O que
ele remove é o **controle unilateral**. É menos do que se quer e é o máximo que a ferramenta
permite hoje.

Candidato upstream: **o arquivo de saída por tarefa do Agent tool vem vazio, logo veredicto de
reviewer rodado por Agent não tem artefato durável.** Lacuna de ferramenta, não defeito de chip.
Doze execuções, zero bytes.

## Sobre o item 1 — o guard não existia quando a prosa que o descreve foi escrita

Vale registrar porque é a forma mais afiada que apareceu no arco inteiro. A ferramenta regenerava
`coordinates.txt` em toda execução, `--strict` inclusive: uma edição à mão seria **sobrescrita em
silêncio** e a execução reportaria sucesso. A alegação era verdadeira sobre a intenção e falsa
sobre o código.

E o que a expôs foi **escrever o must-fail** — a distinção do R-13 (ter teste ≠ ter must-fail)
batendo dentro do conserto do R-22. Quarta vez neste chip que uma regra pega defeito dentro do
remédio de outra regra.

Regra que fica: **gerar e conferir são modos separados; conferidor não escreve.** Verificador com
acesso de escrita ao que ele verifica não é verificador — ele não consegue reprovar.

## Sobre o item 2 — CRLF, e por que isso não é frescura

O guard ia nascer vermelho em toda máquina que não a dele: git converte fim de linha no checkout
aqui, clone novo pega CRLF, a ferramenta escreve LF. Normalizar a comparação e provar as duas
direções foi certo.

O argumento dele é o que importa e eu ratifico: **check que grita lobo é check que alguém
desliga** — e o round 1 deste chip perdeu uma violação real de `gofmt` dentro exatamente desse
ruído. Regra: **rung novo tem de nascer verde em clone limpo, senão não é rung, é treinamento para
ignorar a lane.**

## Aceito sem ruling

Exemções contadas e impressas (alargar exemção é a direção perigosa; a lição do round 4 é que o
defeito nunca é a estreiteza, é o silêncio). `cite-audit.py/.txt` superseded com banner e não
deletados — são o registro do que o round 4 achou, e a linha D24 aponta para eles. Observação de
árvore de trabalho substituída por propriedade verificável em vez de maquiada como reproduzível:
é a leitura certa do R-22(b).

---

# ROUND 6 — os gates PARTIRAM. Cold PASS (semântica), GPT FAIL (derivação + aparato, 5 blockers)

Verifiquei eu mesmo os dois blockers de prosa antes de decidir qualquer coisa. Os dois são
verdadeiros, exatos:

```
$ git grep -c connectorsports 85b6c367 -- .../product_links/application
85b6c367:.../generation_integration_test.go:2      exit=0        (pack diz que nao retorna nada)

$ git log --oneline 917f7bb5..633bf9fa
633bf9fa Format the composition root…   11e68f6f S3   b9da6d2e S2   5bc55219 S1
$ git diff --name-only 917f7bb5 633bf9fa -- apps/ contracts/ packages/ | wc -l   → 13
$ git diff --name-only 917f7bb5 b9da6d2e -- apps/ contracts/ packages/ | wc -l   → 13
```

O range chamado de S1+S2 contém S1, S2, S3 e um commit de formatação. As 13 batem **por
coincidência** — S3 só tocou caminhos já na lista. Contagem certa pelo motivo errado é pior que
contagem errada, porque ela não avisa.

## O que eu vejo e nenhum dos dois gates podia ver, porque cada um enxerga a própria metade

**Nenhum dos dois veredictos sobre o CÓDIGO depende da `cite-table.py`.**

O cold chegou a PASS em C1, C2, C3, C4, C6, C7, C8 **lendo o fonte**. O GPT chegou a PASS em C2,
C3, C5, C9, C10, C12 **re-derivando do git**. Nenhum dos dois usou a ferramenta para chegar ao
veredicto sobre o código.

Os três blockers de aparato são, todos os três, alegações que a ferramenta faz **sobre a própria
cobertura**. Ela gateia nada além de si mesma.

## R-24 — RULING: o aparato deixa de ser GATE e vira RELATÓRIO

Primeiro, o que o R-22 conseguiu e o que ele fez: rounds 2-5 falharam em coordenada de prosa
decaindo. **O round 6 tem zero disso** — cold amostrou ~25 linhas da tabela mais os 20 nomes de
teste em backtick e todos resolvem; a classe sumiu do que ele checou. O round 6 falha nas alegações
de cobertura da FERRAMENTA — um lugar que o próprio R-22 criou ao fazer da ferramenta a fiscal da
regra. Progresso real, e defeito novo real. Não vou re-rotular a segunda metade para fazer sumir.

Mas **"cada round falha num lugar novo" é a regressão infinita**, e conceder round 7 com esse
argumento é como isso nunca acaba. O R-22 disse que a resposta a um aparato auto-referente não é
mais aparato. Eu apliquei isso à PROSA e **não apliquei à ferramenta**. Aplicar agora é o que
termina.

**Regra:** artefato de verificação só sustenta as alegações que consegue fazer **totalmente**.
Alegação que ele não consegue fazer total vira **relatório**, não portão.

Um relatório que diz "isto eu resolvi, isto eu não resolvi" **não tem alegação de totalidade para
ser falsificada**. Os três blockers deixam de existir **por construção, não por perdão** — a
alegação falsa é REMOVIDA, não desculpada. É o ADR-17 honest-unknown apontado para o nosso próprio
aparato de prova, a mesma doutrina que o produto usa para dado que não tem.

Na prática:

1. `cite-table.py` gera a tabela e **imprime o que não resolveu** (`if not locs: continue` → conta e
   lista). **R-17 ratificado pela 4ª vez — o cold perguntou explicitamente se eu o tornava
   bloqueante de novo: sim, a exigência fica.** O que muda é a consequência: linha de relatório, não
   código de saída.
2. **O rung PACK sai da ladder.** Deixa de ser bloqueio de merge.
3. A prosa do pack para de afirmar totalidade em qualquer lugar ("resolve TODA âncora" → "resolve as
   âncoras listadas; as não resolvidas estão abaixo").
4. **Nada de alargar**: não estenda `TOPLEVEL` para `var (…)` agrupado, não enumere mais verbos de
   git, não opere o lookbehind. **Não se alarga reconhecedor cuja saída não gateia mais nada.**

## R-11 morreu junto com o objeto dele, e isso é bom

O R-11 disse que validade de citação é mecânica porque o gate cold não tem shell. Era sobre validar
alegação feita **em prosa**. O R-22 tirou essas alegações. A tabela agora é **navegação, não
evidência**: linha errada custa a um leitor um salto errado, não uma prova falsa. A amostra de 25/53
do cold é proporcional a essa aposta. Aposento a **consequência** do R-11 (portão mecânico), não a
percepção dele.

## Sobre a convergência — é o achado mais forte que este dual gate já produziu

Os dois gates acharam o blocker 3 (âncora nomeada some em silêncio) **independentemente, por
direções opostas, com métodos que não se tocam**: o cold lendo símbolos, com duas instâncias
rastreadas; o GPT atacando a ferramenta, com prova por nome injetado.

O R-18 disse: achado disjunto = cobertura. Isto acrescenta: **convergência entre gates de escopo
diferente é o sinal mais forte que um dual gate consegue emitir** — mais forte que acordo entre
gates que olham a mesma coisa, porque as duas lentes não compartilham método, logo o achado não pode
ser artefato de nenhuma delas. É por isso que o blocker 3 é o que eu trato **estruturalmente** em
vez de dissolver junto com o resto do aparato.

E repare a direção do desacordo: o cold arquivou SHOULD-FIX, o GPT arquivou BLOCKING. **Os dois
estão certos sobre a própria metade.** A divergência de severidade não é descalibragem — é cada gate
pontuando o achado contra o próprio escopo. Só o hub vê os dois escopos.

## O NIT é meu, e conserto do jeito certo

Verificado: o golden põe `500MM` no lado ERP (`sku:33698`) e `50cm` no título do anúncio. Meu
`hub-live-drive-u1-u3.md:40,47` invertia. Corrigido no primary.

Vale nomear o mecanismo: o pack **pina** `main@40623b57`. Corrigir cria um SHA novo; **o pin
continua válido porque pin nomeia SHA**. É o R-14 e o R-20 funcionando exatamente como desenhados —
a resposta certa para artefato pinado que está errado é **commit novo**, nunca edição que faz o pin
mentir.

## Merge: o que fecha, e por que isto não é o R-9 disfarçado

Preciso me testar contra mim mesmo aqui. No R-9 eu neguei "mergear nos veredictos limpos de código e
carregar o pack como dívida". Isto é diferente **em quê**?

No R-9 o defeito de pack **ficava no lugar**, como dívida. Aqui os defeitos de pack **são
consertados**: as duas correções de prosa pousam, as alegações falsas de totalidade são REMOVIDAS.
Nada vira dívida. O que eu recuso é **round de gate adicional sobre o pack corrigido**, não o
conserto.

Mas não posso escrever AGREEMENT por cima de um FAIL entregue. Então:

**O corretivo é gateado no eixo que ele toca, e só nele.** Quem arquivou o FAIL re-checa os
próprios cinco blockers no tip corrigido. Não é round 7 — é re-checagem de um FAIL entregue, pelo
gate que o entregou, no escopo exato do que ele arquivou. **O cold não re-roda**: ele PASSOU, e a
área dele (semântica/código) não é tocada pelo corretivo — re-rodar gate que passou sobre área
intacta é a redundância que o R-18 já desvalorizou.

**Condição que eu mesmo verifico antes de qualquer AGREEMENT**, e é a que mantém o PASS do cold vivo
sem re-rodá-lo:

```
git diff --name-only 85b6c367..<tip corretivo> -- apps contracts packages   →  tem de dar 0
```

Se o corretivo tocar código de produção, o PASS do cold decai e o cold volta. É o R-19 aplicado como
condição de merge em vez de descoberta pós-fato.

**O terminador é estrutural, não uma promessa:** depois disto o rung PACK não está na ladder, logo
nenhum round futuro consegue falhar nele.

## Escopo do corretivo — fechado, e nada além

1. As duas correções de prosa (`connectorsports`, endpoint S1+S2 → `b9da6d2e`).
2. Ferramenta: relata o que não resolveu; para de alegar totalidade; **sem alargamento**.
3. Rung PACK sai da ladder.
4. R-23 registrado no pack.
5. NIT U2 do lado do chip (o do hub já está corrigido).

Zero código de produção. É o que torna a re-checagem estreita legítima.
