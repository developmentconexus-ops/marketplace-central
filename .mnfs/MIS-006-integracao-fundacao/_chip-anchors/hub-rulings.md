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
