# P2.b — Módulo `fiscal`: imposto ex-ante pela matriz do ERP

**Data:** 2026-08-02 · **Missão:** MIS-007-ml-sync / M-06 · **Status:** aprovado pelo operador

**Supersede** `2026-08-02-p2b-imposto-ex-ante-design.md`, que fica como **registro de medição**
(reconciliação contra notas reais, refutações, emendas 1–3). Nenhuma medição daquele documento foi
descartada; o que muda aqui é a arquitetura — módulo próprio, espelho versionado e contrato novo.

---

## 1. Problema

`/pedidos` mostra margem de 69,23 (23%) num pedido cuja margem real é **−2,30 (−0,77%)**.

A causa não é erro de conta: é **fonte inventada**. O caminho de leitura consulta
`pricing_calc_profiles`, que tem **zero linhas**, e cai num default de "SIMPLES 4%". Um imposto
plausível, com cara de fato, que o operador acredita. É pior que campo vazio.

O ERP tem o dado certo. A matriz de ICMS (`TGFICM`) reconcilia **34 de 34 notas reais ao centavo**,
e a fórmula de PIS/COFINS reconcilia **1.437 de 1.437 linhas com erro agregado 0,0000%**.

Esta fatia troca a fonte inventada pela fonte medida.

## 2. Escopo

**Faz:**
- Módulo `fiscal` novo — regras e fórmulas fiscais, puras, sem I/O.
- Espelho **versionado** da matriz de ICMS do ERP.
- `grupo_icms` no espelho de produtos.
- `/pedidos` como primeiro consumidor: imposto decomposto + margem, rotulados **estimativa**.
- Card de saúde do espelho em `/integracoes`.

**Não faz** (fatias seguintes, nomeadas para que ninguém as improvise aqui):

| fatia | conteúdo |
|---|---|
| P3 | `/anuncios` e `/mercado` consomem `fiscal` (pior caso ex-ante) |
| P4 | simulador e `/precos`: piso e sobra; mata os 7 sítios do 4% |
| P5 | tabela de referência legal + detecção de divergência contra o ERP |
| P6 | ligação pedido↔nota emitida; imposto **realizado** lido de `TGFDIN` |

Fora de qualquer fatia: corrigir cadastro do ERP. Isso vai ao contador, não ao código.

## 3. Princípio que governa o desenho

**O ERP é o dono do dado fiscal. Nós lemos, calculamos e mostramos. Nunca digitamos, nunca
completamos lacuna com valor plausível, nunca mantemos tabela legal própria como fonte de cálculo.**

Onde o ERP não responde, a tela diz que não sabe e aponta **qual lacuna preencher no Sankhya**.

## 4. Fontes de dado

Tudo vem de espelho em Postgres. Nenhuma tela consulta Oracle ao vivo — direção mirror-first
ratificada.

### 4.1 `products_mirror` ganha `grupo_icms`

`ALTER TABLE products_mirror ADD COLUMN grupo_icms integer`.

A Q1 do sync já lê `TGFPRO`; acrescenta `GRUPOICMS` ao SELECT e ao upsert do writer.

Decisão registrada: houve no MIS-006 uma regra "nunca ALTER no mirror". Aquilo era regra de colisão
entre chips paralelos no M-07, não arquitetura. Nesta fatia somos escritor único do seam (ADR-024),
então o ALTER é legítimo e evita partir dado de produto em duas tabelas.

### 4.2 `icms_matrix_mirror` — novo, versionado

```sql
CREATE TABLE icms_matrix_mirror (
    tenant_id          uuid        NOT NULL,
    uf_origem          text        NOT NULL,
    uf_destino         text        NOT NULL,
    grupo_icms         integer     NOT NULL,
    codtrib            smallint    NOT NULL,
    aliquota           numeric(6,3),
    aliqintdest        numeric(6,3),
    perc_fcp           numeric(6,3),
    linhas_candidatas  smallint    NOT NULL,
    ambiguo            boolean     NOT NULL,
    vigente_desde      timestamptz NOT NULL,
    vigente_ate        timestamptz,
    synced_at          timestamptz NOT NULL,
    PRIMARY KEY (tenant_id, uf_origem, uf_destino, grupo_icms, vigente_desde)
);
```

**Versionamento.** O sync compara a célula lida com a versão aberta (`vigente_ate IS NULL`). Mudou
qualquer valor → fecha a antiga (`vigente_ate = agora`) e abre uma nova. Célula que sumiu do ERP →
fecha e não abre.

**Consulta as-of:** `vigente_desde <= data_pedido AND (vigente_ate IS NULL OR vigente_ate > data_pedido)`.

Mesmo padrão do custo as-of que o repositório já usa (`GetCostAsOf`).

**Por que versionar.** A alíquota interna da Bahia ficou congelada em 17,0% por 29 meses e foi
corrigida para 20,5% em 20/07/2026. Se o cálculo usasse sempre a matriz de hoje, um pedido de março
passaria a exibir 20,5% e **a divergência desapareceria da tela no exato momento em que o ERP fosse
consertado**. Isso destrói o P5 e o P6 antes de existirem.

**Limitação declarada (D-28).** Nosso histórico começa no primeiro sync. Para os 38 pedidos que já
existem, `as-of` é igual à matriz de hoje — a defasagem da Bahia **não é recuperável** do nosso
espelho. Ela está em `TGFHICM`, que esta fatia não lê. A máquina de versão entra agora porque
começar depois é impossível: o passado não volta.

### 4.3 Resolução da linha acontece no sync, não na leitura

A matriz do ERP tem linhas com restrições. Domínio medido: `N` (sem restrição, código 0), `S`
(curinga, código −1), `I` (grupo ICMS), `O` (TOP), `P` (produto), `H` (NCM), e cinco tipos não
decodificados (`K`, `G`, `L`, `X`, `T`). **A precedência formal entre linhas concorrentes não foi
estabelecida** — houve uma única observação, insuficiente para generalizar.

Regra: **lista branca**. O sync aceita apenas `N/0`, `I/<grupo do produto>` e `S/−1`. Descarta todo o
resto e grava a célula já resolvida.

- Uma linha sobrevive → grava a célula, `ambiguo = false`.
- Mais de uma sobrevive → grava `ambiguo = true`, `linhas_candidatas = N`, **e não escolhe**. O
  leitor devolve desconhecido.
- Nenhuma sobrevive → não grava célula. O leitor devolve desconhecido.

Medido: com a lista branca, as 4 células ambíguas do recorte real resolvem todas para a linha que a
nota confirma. As quatro TOPs de e-commerce (306/307/308/313) não têm linha nenhuma na matriz, então
nada de `O/*` compete hoje — a lista branca é blindagem para o dia em que alguém cadastrar.

Lista branca e não lista negra: tipo desconhecido vira pendência, nunca vira escolha silenciosa.

### 4.4 Cadência

Diária, mais botão manual em `/integracoes`. Alíquota fiscal muda raramente — as 6 edições num único
dia foram excepcionais. O botão resolve o dia em que o próprio operador corrige o cadastro e quer ver
o efeito na hora.

## 5. Módulo `fiscal`

```
internal/modules/fiscal/
├── domain/
│   ├── rules.go        tipos: ICMSRule, Componente, Resultado
│   └── calc.go         TaxesForValue(valor, regra) Resultado   [pura, sem I/O]
├── ports/
│   └── ports.go        TaxMatrixReader
└── adapters/postgres/
    └── matrix_reader.go   lê icms_matrix_mirror + products_mirror, as-of
```

**Por que módulo próprio e não dentro de `pricing`:** hoje `orders` importa repositório de `pricing`
para buscar imposto — dependência invertida, e é de lá que sai o 4% fabricado. Se pedido, anúncio e
simulador consomem a mesma regra, a regra é **mais baixa** que os três; não pode morar dentro de um
deles.

`fiscal` não importa `orders`, `pricing` nem `listings`, e não consulta Oracle. O sync do espelho
continua sendo do `internal_read`, que já é dono dos espelhos.

### 5.1 Tipos

```go
// Componente carrega valor OU motivo. Nunca os dois nulos, nunca os dois preenchidos.
type Componente struct {
	Valor  *float64
	Motivo string // preenchido somente quando Valor == nil
}

type ICMSRule struct {
	CodTrib     int      // 0 = tributado, 60 = ST
	Aliquota    *float64 // alíquota da operação origem→destino
	AliqIntDest *float64 // nil = célula muda
	PercFCP     *float64
	Ambiguo     bool
	VigenteDesde time.Time
}

type Resultado struct {
	ICMS, DIFAL, FCP, PisCofins, Total, CargaICMS Componente
	Origem       string    // "matriz_erp" | "nota_erp" (P6)
	MatrizDesde  time.Time
}
```

`fiscal` é dono da **verdade** e do **texto da pendência**. Cada consumidor é dono da
**apresentação** — a mesma lacuna significa coisas diferentes em `/pedidos`, `/anuncios` e no
simulador.

## 6. Regras de cálculo

`V` = valor do item · `i` = carga total de ICMS sobre `V`

| condição | resultado |
|---|---|
| `ambiguo = true` ou célula ausente | tudo desconhecido + pendência |
| `codtrib = 60` (ST) | `ICMS = DIFAL = FCP = 0`, `i = 0` |
| `codtrib = 0`, `aliquota` nula | ICMS desconhecido, e tudo que depende dele |
| `codtrib = 0`, `aliqintdest` nula | ICMS conhecido · **DIFAL, `i` e PIS/COFINS desconhecidos** |
| `codtrib = 0`, completa | fórmulas abaixo |
| `perc_fcp` nula com linha presente | `FCP = 0` — a matriz respondeu. Diferente de linha ausente |
| `aliqintdest < aliquota` | DIFAL daria negativo → **pendência**. Nunca negativo, nunca zero |

```
ICMS       = V × aliquota
DIFAL      = V × aliqintdest − ICMS
FCP        = V × perc_fcp
i          = aliqintdest + perc_fcp

PIS/COFINS = V × (1 − i) × 0,0925
Total      = ICMS + DIFAL + FCP + PIS/COFINS

Margem     = V × 0,9075 × (1 − i) − comissão − frete − custo
```

**Cascata da célula muda.** Quando `aliqintdest` falta, o `i` fica desconhecido, e como o `i` entra
na base do PIS/COFINS, a margem inteira fica desconhecida. Decisão do operador: **pendência a
resolver** — margem em branco com o motivo apontando a lacuna do Sankhya, e não um valor-teto.
Uma semântica só na tela. Custo aceito: 7 pedidos em 12 meses de 74.

**Por que a base do PIS/COFINS é `V × (1 − i)` e não `V`.** O ICMS é repasse ao estado, não receita.
O STF fixou isso no Tema 69 (RE 574.706). Não estamos aplicando a tese por conta própria: o Sankhya
grava exatamente essa base em **1.434 de 1.434 linhas** não-ST. Consequência contraintuitiva —
quanto maior o ICMS, menor o PIS/COFINS — que é justamente o que a soma ingênua de impostos erra.

**Regime.** Não-cumulativo, PIS 1,65% + COFINS 7,6% = 9,25%, medido em 15.7 mil linhas. O crédito de
9,25% na fórmula de custo do ERP confirma Lucro Real.

**ST.** No produto com ST o ICMS foi pago na entrada e já está dentro do `CUSSEMICM`. Medido:
11.233 de 11.233 linhas intra-MG com ICMS zero. Mas ST é propriedade da **operação**, não do produto
— o mesmo produto sai com `CODTRIB=60` para MG e `CODTRIB=0` pagando 18% para SP. Por isso a
autoridade é o `codtrib` da célula da matriz, nunca o cadastro do produto.

**Custo** é `TGFCUS.CUSSEMICM` as-of, que é o que o repositório já lê em todos os sítios. Nunca
`CUSGER` — essa coluna é preço de venda calculado, e usá-la como custo inverte o sinal da margem.

## 7. Contrato (OpenAPI + `sdk-runtime`, mesmo commit)

Bloco `fiscal` novo no pedido. Todo membro **obrigatório e anulável** — padrão honest-unknown já
usado no repositório.

```yaml
ValorFiscal:
  type: object
  required: [valor, motivo]
  properties:
    valor:  { type: number, format: double, nullable: true }
    motivo: { type: string, nullable: true }   # preenchido só quando valor é null

OrderFiscal:
  type: object
  required: [origem, matriz_desde, icms, difal, fcp, pis_cofins, total, carga_icms_pct]
  properties:
    origem:         { type: string, enum: [matriz_erp, nota_erp] }
    matriz_desde:   { type: string, format: date-time, nullable: true }
    icms:           { $ref: '#/components/schemas/ValorFiscal' }
    difal:          { $ref: '#/components/schemas/ValorFiscal' }
    fcp:            { $ref: '#/components/schemas/ValorFiscal' }
    pis_cofins:     { $ref: '#/components/schemas/ValorFiscal' }
    total:          { $ref: '#/components/schemas/ValorFiscal' }
    carga_icms_pct: { $ref: '#/components/schemas/ValorFiscal' }
```

`origem` já nasce com `nota_erp` no enum porque o P6 usa o **mesmo bloco** — só troca a origem. Nesta
fatia o servidor emite apenas `matriz_erp`.

`decomposicao.imposto` continua existindo, apontando para o total. O objeto `difal` solto de hoje é
marcado **deprecated** e sai no P6.

O bloco aparece em `/orders/{provider_order_id}` (detalhe). Na lista `/orders` não vai o bloco
inteiro — só a margem, que já é coluna existente. Lista carrega o resultado do mesmo cálculo, em
lote, com um único as-of por data de pedido.

## 8. Frente

**Drawer do pedido** ganha seção **Fiscal**: ICMS, DIFAL, FCP, PIS/COFINS e total, com o componente
`UnknownValue` e o motivo como hint onde faltar. Rodapé: *"estimativa pela matriz do ERP, vigência
DD/MM/AAAA"*. A seção DIFAL legada sai.

**Lista** mantém a coluna MARGEM, agora alimentada pelo imposto medido.

**`/integracoes`** ganha card do espelho da matriz no padrão do `SyncHealthCard` existente: última
sincronização, células, ambíguas, mudas, e botão sincronizar.

## 9. O que morre nesta fatia

- `orders/adapters/pricingtax/reader.go` é substituído pelo adapter de `fiscal`.
- `orders` **para de importar** `pricing` — desfaz a dependência invertida.
- `pricing_calc_profiles` sai do caminho de leitura de pedidos.

**Fica vivo até o P4:** os 7 sítios dentro de `pricing` que fabricam 4% (dívida D-21). `/pedidos`
passa a dizer a verdade enquanto `/precos` ainda não. Inconsistência conhecida e aceita pelo
operador, que adiou `/precos`.

## 10. Testes

Todo teste com controle negativo nomeado. Verde sem vermelho provado é byte-idêntico a teste nenhum.

| teste | prova | controle negativo |
|---|---|---|
| Bahia pela matriz vigente | `V=299,90`, alíquota 7,0, aliqintdest 20,5 → ICMS 20,99 · DIFAL 40,49 · PIS/COFINS 22,05 · total 83,53 | trocar 20,5 por 17,0 → tem que dar DIFAL 29,99 e falhar |
| célula muda | `aliqintdest` nula → DIFAL **e** PIS/COFINS nulos, ICMS conhecido | devolver 0 em vez de nulo → falha |
| ST | `codtrib=60` → ICMS, DIFAL e FCP zerados; PIS/COFINS sobre `V` cheio | aplicar alíquota ao ST → falha |
| DIFAL negativo | `aliqintdest < aliquota` → pendência | clampar em 0 → falha |
| as-of | duas versões da célula; pedido anterior à troca pega a antiga | remover o predicado de vigência → falha |
| ambiguidade | duas linhas sobrevivem à lista branca → desconhecido | escolher a primeira → falha |
| margem Bahia | −2,30 e −0,77% | qualquer clamp de margem em zero no caminho → falha |

**Atenção ao ler o teste da Bahia:** neste pedido a comissão e o DIFAL valem o mesmo número (40,49)
por coincidência. Não troque um pelo outro.

**A autoridade é a matriz vigente, não a nota emitida.** A nota 895507 gravou DIFAL 29,99 porque foi
emitida com alíquota defasada; a nota-irmã do dia anterior, mesmo grupo e destino, usou 20,5% e
reconcilia ao centavo. "Consertar" o teste para 29,99 petrificaria a nota defasada como verdade —
exatamente a doença que esta fatia existe para curar.

## 11. Dívidas registradas

| id | dívida |
|---|---|
| D-17 | `CODEMP=1` fixo na leitura de custo |
| D-21 | 7 sítios em `pricing` fabricando 4%, vivos até o P4 |
| D-24 | DIFAL não precificável em PR, RS e SC — 20 células mudas, 7 notas em 12 meses de 74 |
| D-25 | A matriz e os documentos emitidos pelo próprio ERP discordam; o ERP não é auto-consistente |
| D-26 | Fonte da alíquota interna do PR (18,0 em 11 notas) não localizada na base |
| D-27 | DIFAL e FCP dentro da exclusão da base do PIS/COFINS: prática medida do ERP, não jurisprudência fechada. Pendente do contador |
| D-28 | Histórico da matriz começa no primeiro sync; `TGFHICM` guarda o passado real e não é lido |

## 12. Fora do software

`docs/design/evidence/2026-08-02-relatorio-fiscal-difal-operador.md` — divergência de alíquota
interna, 29 meses de atraso na Bahia, lacuna de PR/RS/SC, e as perguntas para o contador. Nenhuma
decisão fiscal é tomada por nós.
