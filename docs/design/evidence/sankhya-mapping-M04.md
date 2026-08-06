# Sankhya → E2.1 — Mapa ratificado (MIS-006 M-04 SankhyaAdapter)

**Fonte:** db-consult ao especialista (sessão MNOS "Marketplace Central database queries",
`local_ec787804`), verificado AO VIVO em METALPRD 2026-07-22 (`all_tab_columns` + amostras +
distribuição real de NUTAB em itens de venda TGFITE).
**Status:** **6/6 itens ratificados** (item 5 resolvido empiricamente — ver abaixo). Regra
NULL-nunca-0 respeitada. Este doc é a base LITERAL do SQL do sync entrypoint
(`internal_read/adapters/oracle/reader.go`) — não hipotetizar além do que está aqui.

## Campos E2.1 → coluna Sankhya

| Campo E2.1 | Origem Sankhya | Notas |
|---|---|---|
| `codigo_produto` | `TGFPRO.CODPROD` | NUMBER, PK real |
| `descricao` | `TGFPRO.DESCRPROD` | |
| `ncm` | `TGFPRO.NCM` | |
| `marca` | `TGFPRO.CODMARCA → TGFMAR.DESCRICAO` | join `TGFMAR.CODIGO = TGFPRO.CODMARCA`; CODMARCA nulo → marca NULL. `TGFPRO.MARCA` (texto) é cache fiel mas o join por código é canônico |
| **`ean`** | **`TGFPRO.REFERENCIA`** | ⚠️ COLISÃO DE NOME: o EAN mora em REFERENCIA (19.656/40.401 produtos têm formato EAN 12–14 díg). ~51% sem EAN → ean NULL |
| **`referencia`** | **`TGFPRO.REFFORN`** | ⚠️ referência de FORNECEDOR (Docol/Deca), NÃO é o EAN |
| `custo` | `TGFCUS.CUSSEMICM` | custo médio SEM ICMS (base de precificação, golden ao centavo). TEMPORAL — ver query abaixo. Sem linha → NULL |
| `preco_venda` | `TGFEXC.VLRVENDA`, `CODTAB=0`, as-of por produto | ✅ RESOLVIDO (item 5). CODTAB=0 = tabela de venda real (prova TGFITE). As-of: última `TGFTAB.DTVIGOR ≤ data` que contém o CODPROD. Ver query abaixo. Sem linha → NULL |
| `estoque_total` | `SUM(TGFEST.ESTOQUE)` | bruto; disponível = `SUM(ESTOQUE − RESERVADO)`. Filtrar CODPARC=0. Ver query |
| `grupo_codigo` | `TGFPRO.CODGRUPOPROD → TGFGRU.CODGRUPOPROD` | grupo FOLHA (ANALITICO='S') |
| `grupo_descricao` | `TGFGRU.DESCRGRUPOPROD` | NÃO único → sempre join por CODGRUPOPROD (código), nunca por texto |
| multi-local | `TGFEST.CODLOCAL` | → `products_mirror_stock_locations` |

## Regras de query (ratificadas)

**Custo (TGFCUS) — temporal, não escalar corrente:**
- Vigente = `DTATUAL ≤ data_ref` → `MAX(DTATUAL)`, desempate `CODLOCAL DESC`, filtro `CODEMP=1`.
- Sem linha → custo NULL.
- Colunas de custo existentes (todas verificadas): CUSSEMICM (escolhida), CUSMEDICM, CUSMED,
  CUSMEDCALC, CUSREP, CUSGER, CUSVARIAVEL. Alternativa se negócio pedir: CUSREP (reposição).

**Estoque (TGFEST):**
- `ESTOQUE` = saldo BRUTO. Disponível honesto = `ESTOQUE − RESERVADO` (opcional: − quando
  `WMSBLOQUEADO='S'`).
- Filtrar `CODPARC=0` (estoque PRÓPRIO; ≠0 = terceiro/consignação).
- `CONTROLE` = lote/grade (somar por cima).
- `estoque_total = SUM(ESTOQUE) por CODPROD, CODPARC=0` (ou `SUM(ESTOQUE−RESERVADO)` se o mirror
  quiser disponível — decisão de F-01, default recomendado = disponível).
- Produto ausente de TGFEST → NULL; saldo 0 real é legítimo (não confundir).

**EAN / auto-vínculo:**
- `TGFBAR` está VAZIA nesta instância (rowcount=1) — NÃO usar CODBARRA.
- EAN = `TGFPRO.REFERENCIA`, match exato. Sem multi-barcode a resolver.

**Grupo (TGFGRU):**
- Hierarquia: `CODGRUPAI` (pai), `ANALITICO` ('S'=folha), `GRAU` (profundidade).
- TGFPRO aponta pro grupo FOLHA. Use a folha para código/descrição; raiz/caminho = subir por
  CODGRUPAI até GRAU raiz (não necessário para E2.1).

## Item 5 — RESOLVIDO empiricamente (CODTAB=0, as-of por produto)

Resolvido pela prova de dados (não schema-guess), não exigiu decisão do operador:

- **FATO 1 (qual tabela):** `TGFITE` grava o `NUTAB` usado em cada venda. Distribuição de NUTAB em
  TODOS os itens de venda 2026 (TOPs 101,21,116,303,313,701,203) → **todos mapeiam `CODTAB=0`**.
  CODTAB=0 tem versão nova ~diária (`DTVIGOR` diário; NUTAB 11809=27/06 … 11887=21/07). Logo
  CODTAB=0 = a tabela de preço de venda; NUTAB = versão datada. CODTAB 200/501/502 = outras listas
  (não usadas em notas de venda).
- **FATO 2 (VLRVENDA é o preço):** onde `(NUTAB,CODPROD)` existe em TGFEXC,
  `TGFITE.VLRUNIT_vendido == VLRVENDA` EXATO (diff=0) e == `PRECOBASE`. Sem exceção.
- **FATO 3 (as-of por produto — crítico):** cada NUTAB diário é ESPARSO (~121 preços/versão, só
  produtos que mudaram naquele dia). O preço vigente HERDA da última versão CODTAB=0
  (`DTVIGOR ≤ data`) que contém aquele CODPROD. Prova: as-of por produto p/ vendas de jul → 643
  itens, 0 sem preço, 619 == VLRUNIT exato (96,3%); os 24 restantes = preço alterado na mão na
  venda (`TGFITE.ALTPRECO`/negociação), esperado, não é falha de mapeamento.

**SQL literal (preço de venda corrente, as-of):**
```sql
select e.VLRVENDA
from METALPRD.TGFEXC e
join METALPRD.TGFTAB t on t.NUTAB = e.NUTAB
where e.CODPROD = :cod
  and t.CODTAB = 0
  and t.DTVIGOR <= :data_ref        -- hoje p/ preço corrente
order by t.DTVIGOR desc
fetch first 1 rows only;
```
Sem linha → `preco_venda = NULL` (nunca 0). NÃO usar `AD_ECOM` (vazio) nem outro CODTAB.

**Nota de override (operador):** CODTAB=0 = preço de VENDA de nota (varejo/operacional), provado
por transação real. Se o preço-alvo de MARKETPLACE precisar ser diferente do preço de nota, é
ajuste de política do operador — landa como override aditivo depois; o mirror grava CODTAB=0 como
referência base honesta enquanto isso.

## Adendo — preço de LISTA vs PROMOÇÃO são separados (METALPRD 2026-07-22)

Tabelas de preço = só especificação, SEM desconto embutido. Desconto/campanha vive em TGFDES.
NÃO afeta o mapeamento de `preco_venda` do mirror (= CODTAB=0 as-of, inalterado). Registrado para
M-06 (se houver display de promoção/preço-de-canal na tela) e para o motor de preço:

- **CODTAB=0** = preço de venda base (100% das notas). `preco_venda` do mirror = este.
- **CODTAB=3** = lista E-COMMERCE, DERIVADA: `PERCENTUAL=10`, `FORMULA=null`, `CODTABORIG=0`,
  0 linhas próprias em TGFEXC → preço = base×0,9. É lista de canal, NÃO campanha. **Candidato a
  "preço de marketplace"** se o operador quiser expor o preço de canal e-commerce em vez do base.
- CODTAB 200/501/502 = outras listas de canal; nenhuma usa `TGFEXC.PERCDESC` (=0); `TIPO='V'`.

**Promoção/campanha (TGFDES, 3245 linhas, TIPPROMOCAO='P'):** `NUPROMOCAO` (id), `DESCRPROMOCAO`
(nome, ex "17JUNHO26" — número embutido = % desconto), `PERCENTUAL` (o %), `DTINICIAL`/`DTFINAL`
(janela), `LIQUIDACAO='S'`. Escopo por `CODPROD`/`GRUPODESCPROD`/`CODPARC`/`CODVEND`/`CODTAB`/
`CODEMP`/`CODLOCAL` (campanhas de junho miram por `GRUPODESCPROD` = tag grupo-de-desconto no
produto). Ligação na nota: `TGFITE.NUPROMOCAO → TGFDES.NUPROMOCAO`; desconto materializa em
`TGFITE.VLRDESC` (VLRUNIT fica no preço de lista). Hoje 0 promos ativas (junho fechou 20/07); 50
agendadas.

- Promoção ATIVA de um produto hoje:
```sql
select NUPROMOCAO, DESCRPROMOCAO, PERCENTUAL, DTINICIAL, DTFINAL
from METALPRD.TGFDES
where DTINICIAL <= :hoje and DTFINAL >= :hoje
  and (CODPROD = :cod or GRUPODESCPROD = <grupo-desc-do-produto> or (CODPROD=0 and GRUPODESCPROD is null));
```
- Desconto REALIZADO por item vendido = `TGFITE.VLRDESC` + nome via join TGFDES por NUPROMOCAO
  (histórico já resolvido, não recomputar).

**Amarração promoção↔produto RESOLVIDA (para M-06):** `TGFPRO.GRUPODESCPROD = TGFDES.GRUPODESCPROD`
(1.441 produtos taggeados). A tag no produto é a ATUAL (sobrescreve quando entra campanha nova) →
serve p/ promo ATIVA; histórico usa `TGFITE.NUPROMOCAO`. TGFDES também pode mirar `CODPROD` direto.
Validado: prod 39069 tag '20JUNHO26' → casa promo #3488 20% na janela 03/06–20/07. NÃO é coluna do
`products_mirror` (migration 0076 não tem promo) → superfície de M-06, não M-04.

## Pack de queries validado ao vivo (METALPRD 2026-07-22) + shape do mirror ratificado

Bind vars `:cod :emp :data_ref`. Todas rodadas READ-ONLY em METALPRD pelo especialista. Estas são
a base LITERAL do SQL do sync entrypoint — o chip escreve o Go equivalente, TENANT-SCOPED, sobre o
writer de `products_mirror` de M-02 (nunca copia a query crua sem o filtro de tenant, AC-01).

**Shape do `products_mirror` (ratificado pelo schema M-02 migration 0076 + ADRs da missão):**
1. **Grão** = 1 linha por `CODPROD` (as-of hoje) em `products_mirror`; estoque multi-local vai na
   tabela filha `products_mirror_stock_locations` (M-02). `em_promocao`/promo NÃO é coluna do 0076
   → superfície de M-06, fora de M-04.
2. **data_ref** = HOJE (snapshot corrente). Sync é cadence-agnostic (D6); o app não passa data
   histórica no sync. Preço/custo/estoque gravados são o vigente as-of hoje.
3. **Empresa** = `CODEMP=1` fixo (cliente atual single-empresa METALPRD). Revisitar só se surgir
   tenant multi-empresa (não é o caso da fundação).
4. **RAW facts, não pré-calculado**: mirror grava `preco_venda`=CODTAB=0 base e `custo` SEPARADOS;
   NÃO grava um `preco_final` com promo/desconto aplicado. Motor de preço (/precos) + app calculam
   o final downstream (princípio: mirror guarda fato, decisão de preço é por margem). Promo e
   e-commerce (CODTAB=3, base×0,9) derivam do base, não o contrário.
5. **Sem superfície extra em M-04**: margem = `preco−custo` já é do motor de preço existente; lista
   de promo ativa (Q5) = M-06. M-04 usa só Q1–Q4.

```sql
-- Q1 products_mirror BASE (identidade+grupo+marca+ncm+ean). ATIVO='S' = 10.526 produtos.
select p.CODPROD, p.DESCRPROD, p.NCM,
       p.REFERENCIA as EAN,            -- EAN-13 (TGFBAR vazia)
       p.REFFORN    as REF_FORNEC,     -- cód fornecedor (≠ EAN)
       p.CODGRUPOPROD, g.DESCRGRUPOPROD,
       p.CODMARCA,     m.DESCRICAO as MARCA,
       p.GRUPODESCPROD as TAG_DESCONTO, p.ATIVO, p.USOPROD
from METALPRD.TGFPRO p
left join METALPRD.TGFGRU g on g.CODGRUPOPROD = p.CODGRUPOPROD
left join METALPRD.TGFMAR m on m.CODIGO       = p.CODMARCA
where p.ATIVO = 'S';

-- Q2 CUSTO as-of (CUSSEMICM, líquido de ICMS). :emp=1
select CUSSEMICM
from METALPRD.TGFCUS
where CODPROD=:cod and CODEMP=:emp and DTATUAL<=:data_ref
order by DTATUAL desc, CODLOCAL desc
fetch first 1 rows only;

-- Q3 PREÇO DE LISTA as-of (tabela de venda CODTAB=0)
select e.VLRVENDA
from METALPRD.TGFEXC e
join METALPRD.TGFTAB t on t.NUTAB = e.NUTAB
where e.CODPROD=:cod and t.CODTAB=0 and t.DTVIGOR<=:data_ref
order by t.DTVIGOR desc
fetch first 1 rows only;

-- Q4 ESTOQUE por local (próprio = CODPARC=0). Total = soma; disponível = ESTOQUE-RESERVADO.
select CODLOCAL,
       sum(ESTOQUE)           as estoque_bruto,
       sum(ESTOQUE-RESERVADO) as disponivel
from METALPRD.TGFEST
where CODPROD=:cod and CODPARC=0
group by CODLOCAL;

-- Q5 (M-06, NÃO M-04) PROMOÇÃO ATIVA do produto (melhor % vigente); escopo por-produto e por-grupo
select d.NUPROMOCAO, d.DESCRPROMOCAO, d.TIPPROMOCAO,
       d.PERCENTUAL, d.VLRDESC, d.VLRVENDA, d.DTINICIAL, d.DTFINAL
from METALPRD.TGFDES d
join METALPRD.TGFPRO p on p.CODPROD=:cod
where trunc(:data_ref) between trunc(d.DTINICIAL) and trunc(d.DTFINAL)
  and ( d.CODPROD = :cod
     or (nvl(d.CODPROD,0)=0 and d.GRUPODESCPROD = p.GRUPODESCPROD and p.GRUPODESCPROD is not null) )
order by d.PERCENTUAL desc
fetch first 1 rows only;
-- preco_promocional: PERCENTUAL>0 -> preco_lista*(1-PERCENTUAL/100); VLRVENDA>0 -> VLRVENDA fixo; VLRDESC>0 -> preco_lista-VLRDESC.

-- Q6 (histórico, M-06/pedidos) DESCONTO REALIZADO por item vendido — já resolvido, não recomputar
select i.NUNOTA, i.CODPROD, i.QTDNEG, i.VLRUNIT, i.VLRDESC,
       i.NUPROMOCAO, d.DESCRPROMOCAO, d.PERCENTUAL
from METALPRD.TGFITE i
left join METALPRD.TGFDES d on d.NUPROMOCAO = i.NUPROMOCAO
where i.NUNOTA = :nunota;
```

## Estratégia de sync — RATIFICADA (full-snapshot v1, incremental v2 com chaves provadas)

Consulta de change-tracking ao especialista (METALPRD 2026-07-22, tudo `ran`). Decisão de design,
não improviso — registrada para M-04 (v1) e para a evolução (v2).

### Matriz de change-tracking por tabela (provada ao vivo)

| Tabela | Chave incremental | Verdict |
|---|---|---|
| TGFPRO | `DTALTER` (100% preenchida, 40.401/40.401, COM hora intraday) | ✅ confiável; cruzado c/ TSILGT (13.824 eventos/30d) |
| TGFCUS | `DHALTER` / `DTATUAL` — série temporal APPEND-ONLY (produto top = 383 linhas, sem update in-place) | ✅ perfeito p/ incremental |
| TGFEXC | `DHALTREG` (timestamp da linha) | ✅ MAS obrigatório: `DTVIGOR` sozinho PERDE ~4% (23/571 preços desde jun editados IN-PLACE depois da versão criada, `DHALTREG > DTVIGOR`) |
| TGFEST | **NENHUMA** (DTVAL/DTFABRICACAO/DTENTRADA ≠ last-change) | ❌ incremental só via ledger TGFESE (~1,01M linhas, DTMOV) + recompute — não vale; full = 28.010 linhas, trivial |
| TGFMAR | nenhuma (759 linhas) | full trivial |
| TGFGRU | `DHALTER` (baixa freq) | ✅ ou full |
| TSILGT (log nativo) | field-level, 24,9M linhas, MAS cobertura seletiva: só TGFPRO/TGFFIN/TGFPAR/TGFCTB | ❌ NÃO cobre TGFEST/TGFEXC/TGFTAB/TGFCUS → não serve de feed universal |

Cadência real do cliente: preço ~business-daily (24 versões CODTAB=0/mês), custo intraday
(append), cadastro ~1.000 produtos/30d intraday.

### v1 (M-04) = FULL-SNAPSHOT por rodada — decisão ratificada

Não é atalho; é o desenho correto para este dataset:
1. **keep-absent (ADR-031, era ADR-04) exige o conjunto completo** — delta não sabe marcar ausente. Qualquer
   incremental ainda precisaria de sweep completo periódico só para ausência.
2. **TGFEST não tem change-ts** — estoque força full de qualquer jeito; incremental não elimina o
   full, só o complica.
3. **Imune a update in-place** — snapshot relê o valor corrente (Q2/Q3 as-of) a cada rodada; os
   ~4% de edições in-place em TGFEXC são capturados naturalmente. Incremental ingênuo (DTVIGOR)
   os perderia em silêncio.
4. **Volumetria minúscula**: 10.526 ativos × (2 point-lookups indexados) + 28.010 linhas TGFEST.
   Rodada completa = segundos.

Cadência v1: diária via scheduler M-01 (cadence-agnostic, D6) — cursor de sync_state registra a
rodada, não filtra fonte.

### v2 (futuro, se intraday/escala exigir) — desenho pré-aprovado

Incremental liga APENAS nas 3 chaves provadas, mantendo estoque full:
- TGFPRO: `WHERE DTALTER > :cursor` (pega inclusive flip `ATIVO='N'` — desativação toca DTALTER,
  logo ausência de cadastro é detectável incrementalmente).
- TGFCUS: `WHERE DHALTER > :cursor` (append-only).
- TGFEXC: `WHERE DHALTREG > :cursor` — **NUNCA DTVIGOR** (perde in-place).
- TGFEST: continua full (28k linhas).
- Cursor por tabela em `sync_state` (M-01 já suporta), com overlap de segurança contra clock-skew
  + upsert idempotente. Sweep full periódico (ex. 1×/dia) continua obrigatório p/ keep-absent de
  produto deletado fisicamente (raro, não coberto por DTALTER).
