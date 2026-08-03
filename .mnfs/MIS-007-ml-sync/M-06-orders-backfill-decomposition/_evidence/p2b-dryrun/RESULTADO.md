# P2.b — dry run da fórmula contra dado real (2026-08-02)

> ⚠️ **RODADA 1. Números de impacto SUPERSEDIDOS pela §10 (rodada 2, 2026-08-03).**
> Duas coisas mudaram depois desta rodada: a base do PIS/COFINS (era 9,25% / 8,62% da receita, é
> `MAX(0, P×(1−a) − S)`) e a existência da restituição de ST, ausente aqui. As §§1–3 (recorte, lista
> branca, reconciliação, defasagem do `TGFICM`) **continuam válidas** — são medição de fato, não de
> fórmula. A §4.1 (impacto), a §5 (controle positivo, **−6,10**) e a linha 6 de "O que muda no
> plano" estão **mortas**; leia a §10.
>
> A §5 é o caso a lembrar: bateu ao centavo com um documento independente e ainda assim estava
> errada — os dois modelos carregavam os mesmos dois defeitos.

Rodado **antes** de executar o plano, a pedido do operador. Nada foi implementado.

**Fontes:** Oracle `METALPRD` ao vivo (sonda read-only `cmd/fiscalprobe`, só `SELECT`) e o
Postgres do dev stack (38 pedidos ML reais, `products_mirror` com 10.628 linhas).

**Artefatos:** `q01..q10*.sql` (consultas Oracle), `calc.sql` (fórmula sobre os 38 pedidos),
`reconcilia.sql` (previsão × nota emitida).

---

## 1. Recorte real

8 produtos distintos nos 38 pedidos. `TGFPRO.GRUPOICMS`:

| grupo | produtos | pedidos |
|---|---|---|
| 122 | 22467, 39563, 39587, 41912, 42194 | 22 |
| 311 | 15956 (FLEX PAPELEIRO — o mais vendido) | 10 |
| 0 | 20303, 20322 | 2 |
| — | sem vínculo (SKU `008720CE`) | 6 |

Destinos: BA, ES, MA, MG, PE, PR, RJ, RS, SP. **Todos os 38 pedidos têm uma linha só** — a
porta por item não é exercitada pelo dado de hoje (continua correta, mas sem cobertura real).

## 2. A lista branca funciona, e a ambiguidade é real

Matriz lida com `UFORIG=13` (MG) para os 3 grupos × 9 UFs. Duas formas coexistem:

- interestadual: `TIPRESTRICAO='N', CODRESTRICAO=0` + `TIPRESTRICAO2='I', CODRESTRICAO2=<grupo>`
- intra-MG: invertido — `TIPRESTRICAO='I', CODRESTRICAO=<grupo>` + `TIPRESTRICAO2='S', CODRESTRICAO2=-1`

**SP grupo 122 tem 3 linhas concorrentes** (`X/101`, `O/101`, `N/0`). A lista branca aceita só a
`N/0` → resolve unívoco. ✔

**MG grupo 0 tem 2 linhas, ambas `O/…`** → lista branca descarta as duas → **sem célula**, e o
leitor devolve desconhecido. ✔ Comportamento correto, e acontece com dado real.

**A linha `X` não é lixo não-decodificado.** As notas com `CLASSIFICMS='X'` (PJ contribuinte)
saem com ICMS zero, coerente com a linha `X/101` (`CODTRIB=10, ALIQUOTA=0`). Descartá-la está
certo **porque venda ML é sempre PF**, não porque seja ruído. Isso precisa estar escrito.

## 3. Reconciliação contra nota emitida — 21 documentos

TOP 306, `CLASSIFICMS='C'`, jun–jul/2026. (A TOP 313 é o mesmo documento pelo `TGFVAR`; incluí-la
dobraria a amostra sem acrescentar evidência.)

| componente | resultado |
|---|---|
| **ICMS de saída** | **20 ok · 1 desconhecido · 0 errado** |
| **DIFAL (variante A)** | 16 ok · 3 desconhecido · **2 divergentes** |
| **FCP** | 21 ok |

### A regra de arredondamento do DIFAL foi decidida por medição

Duas variantes dão resultados diferentes na nota 894296 (SP, R$ 599,80):

```
A: round(V × (ALIQINTDEST − ALIQUOTA) / 100)      = 35,99   ← a nota
B: round(V × ALIQINTDEST) − round(V × ALIQUOTA)   = 35,98
```

**A vence.** É a única nota da amostra que discrimina as duas — as outras 20 são indiferentes.
O plano tinha a variante B escrita.

### As 2 divergências são do ERP contra si mesmo, não da fórmula

**BA, nota 895507 (+10,50):** caso já conhecido. A nota usou 17,0; a matriz diz 20,5. A nota-irmã
895436, mesma UF/grupo, reconcilia **exato** com a matriz. Emenda 1 já documentou.

**RJ (−3,00):** a matriz diz `ALIQINTDEST=18`; a nota de 01/07 implica 19. **Não julguei com uma
nota só** — puxei a série:

| UF | mês | alíq. interestadual | DIFAL implícito | FCP | notas |
|---|---|---|---|---|---|
| RJ | 2026-06 | 12 | **6%** | 2 | 6 |
| RJ | 2026-07 | 12 | **7%** | 2 | 2 |
| PR | 2026-07 | 12 | 6% | 0 | 2 |
| RS | 2026-06 | 12 | 5% | 0 | 1 |

RJ mudou de 6% para 7% na virada de julho, em 2 documentos de dias diferentes, contra 6 de junho
consistentes. **A matriz está defasada no RJ**, ~1 mês.

> Isto inverte o caso da BA. Lá a matriz estava certa e a nota velha. Aqui a matriz está velha.
> **"A fonte é a matriz" continua sendo a única escolha possível ex-ante** — quando precificamos
> o documento não existe — mas ela erra hoje, no RJ, para menos.

### PR e RS: a matriz é muda, e as duas se comportam diferente

Ambas têm `CODTRIB=0` (DIFAL devido) com `ALIQINTDEST` nulo.

- **RS:** DIFAL real 5% → interna 17% = exatamente o `ALIQUFDEST` da matriz. Reconcilia. **1 nota só.**
- **PR:** DIFAL real 6% → interna 18%. O `ALIQUFDEST` do PR é **19**. Não reconcilia com nada da matriz.

O tratamento decidido (D-37: pendência explícita) está certo para o PR. Para o RS há indício de
que `ALIQUFDEST` serve — **mas uma nota não fecha regra**, e foi assim que esta missão já fabricou
um veredito errado. Fica como pendência, não como fallback.

## 4. O achado que muda o plano: PIS/COFINS de saída não tem fonte

```
TGFCAB.VLRPIS / VLRCOFINS / BASEPIS  →  NULL em 100% das notas de 2026,
                                        em TODAS as CODTIPOPER, não só nas nossas.
```

O crédito de PIS/COFINS **da entrada** é medido — está dentro do `CUSSEMICM` via `METAL_CUSTO`
(9,25% hardcoded no ERP). O **débito da saída** não existe em documento nenhum do Sankhya. É
apurado no livro fiscal, fora do alcance desta fatia.

O plano escrevia `PIS/COFINS = 9,25% da receita` como componente calculado. A spec anterior usava
8,62%. **Os dois são chute.** É a mesma classe de defeito que a fatia existe para matar — um
default plausível com cara de fato — só que com melhores modos.

E não é um detalhe: é o segundo maior componente e é ele que decide o sinal da margem.

### Impacto medido nos 38 pedidos

| | pedidos com margem positiva |
|---|---|
| como a tela mostra **hoje** (SIMPLES 4%, DIFAL 0) | **33** |
| fórmula nova **sem** PIS/COFINS | **27** |
| fórmula nova **com** PIS/COFINS a 8,62% | **10** |
| **mudam de sinal só por causa do PIS/COFINS assumido** | **17** |

Soma das margens (só os 28 calculáveis): hoje **4.123,71** · sem PC **1.959,25** · com PC **496,28**.

**10 dos 38 não são calculáveis** — 6 sem vínculo de produto, o resto por célula ausente
(grupo 311 em RJ, grupos mudos em PR/RS, grupo 0 intra-MG).

## 5. Controle positivo

Pedido BA de R$ 299,90 (`2000017515486360`, produto 41912, grupo 122):

| | valor |
|---|---|
| ICMS | 20,99 |
| DIFAL | 40,49 |
| FCP | 0,00 |
| PIS/COFINS (8,62%) | 25,85 |
| **margem calculada aqui** | **−6,10** |
| margem na Emenda 1, por caminho independente | **−6,10** |
| **o que a tela mostra hoje** | **+69,23** |

Bate ao centavo com um documento escrito antes, sem consultar este cálculo.

---

## O que muda no plano

| # | mudança | onde |
|---|---|---|
| 1 | DIFAL = `round(V × (AID − ALIQ)/100)`, **não** diferença de arredondados | Task 4 |
| 2 | PIS/COFINS não pode ser constante no código — decisão do operador pendente | Task 4 |
| 3 | `ALIQUFDEST` não é "coluna de ST": é a interna do destino, e no RS é a única que reconcilia | Task 3 |
| 4 | `TIPRESTRICAO='X'` = ramo contribuinte; descarte é correto **porque ML é PF** | Task 3 |
| 5 | célula ausente é caso comum, não borda: 10 de 38 | Task 4/6 |
| 6 | matriz do RJ está defasada **agora**; a tela vai subestimar o DIFAL em 1pp | dívida nova |

---

# 10. Rodada 2 — corrigida (2026-08-03)

**Artefato:** `calc2.sql`. Substitui `calc.sql` e `impacto.sql`.

## 10.1 O que estava errado na rodada 1

| eixo | rodada 1 | rodada 2 | por quê |
|---|---|---|---|
| base do PIS/COFINS | 9,25% da **receita** | `MAX(0, P×(1−a) − S)` | STF Tema 69 tira o ICMS da base; STJ Tema 1125 + SC COSIT 100/2025 tiram o ICMS-ST |
| restituição de ST | **ausente** | `R = Σ VLRUNITMED` de `TIPIMPOSTO IN ('S','I')`, só se UF ≠ MG | gatilho `UF <> 13` literal em `PAN_GET_CUSVAR_MNOBRE`; vale 20–25% do custo |
| alíquota interna | do `TGFICM` | tabela legal nossa, 27 UFs | `TGFICM` está 2 anos defasado em RJ, PR e DF |
| DIFAL | `interna − interestadual` | gross-up por dentro | decisão do operador: seguir a lei (D-41) |

Os dois primeiros eixos **puxam em direções opostas** ao primeiro erro: a rodada 1 superestimava a
carga. Por isso o veredito mudou de "17 pedidos viram negativo" para "3".

## 10.2 Resultado — 38 pedidos ML reais

```
pedidos | nao_calculavel | positivos_hoje | positivos_novo | viram_negativo | soma_hoje | soma_nova
     38 |              7 |             33 |             28 |              3 |   4560.40 |   1927.01
```

A margem agregada dos pedidos calculáveis cai **58%** (4.560,40 → 1.927,01). Isso é o tamanho do
buraco que a tela esconde hoje com o imposto fixo de 4%.

**7 não-calculáveis:** 5 sem vínculo de produto (SKU `008720CE`), 2 porque o grupo `311` não tem
célula de `TGFICM` para RJ. Os dois últimos são pendência nomeada legítima (D-37), não defeito.

## 10.3 Controle positivo — pedido `2000017515486360` (BA)

Produto `41912`, `grupo_icms = 122`, `origprod = 0`, `codtrib = 0` (não-ST na BA).

| passo | valor |
|---|---|
| `a_inter` | 7,0% (BA não está em {SP,RJ,PR,SC,RS}, origem nacional) |
| `a_int` (legal, BA) | 20,5% |
| `a_ef = 0,205 × (1−0,07) / (1−0,205)` | **23,98%** |
| receita | 300,00 |
| ICMS total | 71,92 |
| base PIS/COFINS = `300 × (1−0,2398) − 22,58` | 205,48 |
| PIS/COFINS | 19,00 |
| restituição (UF ≠ MG) | **+37,69** |
| custo `CUSSEMICM` | 154,53 |
| comissão | 40,49 |
| frete | 23,65 |
| **margem nova** | **+28,00** |
| margem exibida hoje (imposto 4%) | +69,23 |

**Este número supersede o −6,10 da rodada 1.** O −6,10 tinha batido com a Emenda 1 do documento do
MNOS, o que na hora pareceu confirmação — mas os dois modelos compartilhavam os **mesmos dois erros**
(sem restituição, base errada de PIS/COFINS). Concordância entre dois modelos com o mesmo defeito não
é evidência. O que decidiu foi ler o gatilho `UF <> 13` na fonte da procedure e a jurisprudência da
base, não a comparação.

## 10.4 Achado novo — o bloco de ST envelhece sem sinal

`TGFEFDVMRSTDIA`, `MAX(DTMOV)` por produto:

| produto | `DTMOV` | pedidos |
|---|---|---|
| 15956 (FLEX PAPELEIRO) | **2022-02-04** | 10 |
| 22467, 39563, 39587, 41912, 42194 | 2026-07-31 | 22 |

O produto **mais vendido** carrega `S` e `R` de quatro anos atrás, e a origem não sinaliza isso de
forma nenhuma — o valor chega igualzinho ao de ontem. Vira **D-44**: esta fatia expõe a idade
(`fiscal_dt_ref` na tela); política de expiração não é desta fatia.
