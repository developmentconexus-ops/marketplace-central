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

## Não sujeito a ruling (do chip, aceito como plano)

Shape A (4 blocos de evidência citados mas nunca escritos no arquivo), os 3 testes R5 dos guards
alcançáveis, a 4ª linha acima de 500 no teste de limites (a atual não mata fórmula clampada), a
inversão do par golden no C6, e o re-gate em tip congelado. O erro de despachar gate com pack
não commitado é dele e o conserto é o certo: congela e re-despacha os DOIS.
