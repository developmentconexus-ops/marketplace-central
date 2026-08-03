# P2.b — Imposto ex-ante no motor de preço (plano re-escopado)

> **For agentic workers:** REQUIRED SUB-SKILL: use `superpowers:subagent-driven-development` ou
> `superpowers:executing-plans` para executar task a task. Passos usam checkbox (`- [ ]`).

**Goal:** trocar a única fonte inventada do imposto — `pricing_calc_profiles` vazio caindo em
"SIMPLES 4%" — pela matriz de ICMS do ERP, dentro do motor de preço que **já existe**.

**Architecture:** o imposto mora onde a fórmula do preço mora (`pricing/domain`). `/pedidos`
continua lendo pela ponte que já existe (`orders/adapters/pricingtax`). Não há módulo novo.

**Spec:** [`docs/superpowers/specs/2026-08-02-p2b-modulo-fiscal-design.md`](../specs/2026-08-02-p2b-modulo-fiscal-design.md)
— a **medição** daquele documento continua valendo integralmente (34/34 notas ao centavo,
1.437/1.437 linhas de PIS/COFINS, lista branca de resolução de célula, armadilhas do DDL do
`TGFICM`). O que esta revisão troca é **a arquitetura**, não a evidência.

**Supersede:** [`2026-08-02-p2b-modulo-fiscal.md`](2026-08-02-p2b-modulo-fiscal.md) (17 tasks).
Aquele plano fica como registro; nenhuma medição dele foi descartada.

---

## Por que este plano substitui o de 17 tasks

O plano anterior (2874 linhas, módulo `fiscal` próprio, leitura *as-of* pela data do pedido,
remoção do `pricingtax`) desenhava um **motor fiscal de registro** para `/pedidos`. Decisão do
operador em 2026-08-02: P2.b entrega **estimativa ex-ante**; o número consolidado vem da nota do
ERP, em fatia própria (P2.c). Com isso:

| construção do plano velho | destino |
|---|---|
| módulo `fiscal` separado (domínio + porta + adapter) | **morre** — a fórmula já é `pricing/domain/decompose.go` |
| remoção do `orders/adapters/pricingtax` | **morre** — é justamente a ponte "uma verdade só"; ela fica |
| leitura *as-of* pela data do pedido, threaded pela porta | **adiada p/ P2.c** — quem reconcilia é a nota |
| card de saúde do espelho em `/integracoes` | **adiada** — não é o que destrava a conta |
| escrita versionada (SCD-2) da matriz | **FICA** — ver abaixo |
| porta de pedido → **item** | **FICA** — é requisito real, não inchaço |

**Sobre o versionamento.** O que separa barato de caro aqui é **escrita vs. leitura**:

- **Escrita versionada: fica.** É uma comparação com a versão aberta e um `UPDATE ... SET
  vigente_ate` dentro do job de sync. Custo pequeno, e preserva o que não volta — spec §4.2 (D-28):
  o histórico começa no primeiro sync, e começar depois é impossível.
- **Leitura as-of: sai.** Todo consumidor desta fatia lê a **célula vigente**
  (`vigente_ate IS NULL`). A assinatura da porta não carrega data. Quem precisa de as-of é o P2.c.

**Sobre a porta por item.** O ICMS depende do `GRUPOICMS` do produto. Um pedido com itens de grupos
diferentes não fecha com uma alíquota única aplicada ao total. Isso é correção de fato, não escopo
extra — mas é uma **extensão da porta existente**, não um módulo novo.

## Escopo — o que esta fatia entrega

Cinco buracos, medidos contra 38 pedidos reais e 21 notas emitidas, e nada além:

1. **`imposto` = `AliquotaPct` × preço** com `pricing_calc_profiles` em zero linhas → default
   SIMPLES 4% entra como fato. Único número fabricado do jeito antigo.
2. **ICMS de saída não é alíquota de regime.** Depende do grupo do produto **e do destino** — a
   mesma peça sai ST em MG e tributada em SP. Produto ST não tem ICMS de saída: o próprio já ficou
   no custo (`CUSSEMICM`, ramo ST).
3. **PIS/COFINS débito não existe na decomposição**, e a base **não é a receita**. É
   `MAX(0, P×(1−a) − S)`: tira o ICMS (STF Tema 69) e tira o ICMS-ST retido na entrada
   (STJ Tema 1125 + SC COSIT 100/2025). Medido em 99,76% ao centavo.
4. **Restituição de ICMS-ST não existe em lugar nenhum do sistema.** Em venda interestadual ela
   abate o custo, e vale 20–25% dele. Sem ela o piso sai ~25% alto.
5. **DIFAL É buraco** — o oposto do que este plano dizia antes.

### Por que o DIFAL virou buraco

`pricing/domain/difal.go` está pronto, mas está pronto **errado em dois eixos**:

- **método:** faz `efetivo = interna − interestadual` ("por fora", como o ERP). O operador decidiu
  seguir a **lei** (gross-up por dentro, LC 87/96 art. 13 §1º I + §7º, Convênio ICMS 236/2021
  cláusula 2ª §1º). Diferença de +2% (SC/RS) a +10% (MA) no piso.
- **tabela:** as internas por UF estão defasadas. Medido contra a lei: **RJ 18 → 22**
  (Lei 10.253/2023 + LC 210/2023, desde 20/03/2024), **PR 18 → 19,5** (Lei 21.850/2023, desde
  18/03/2024), **DF 18 → 20** (Lei 7.326/2023, desde 21/01/2024).

O `TGFICM` do ERP carrega a mesma defasagem — **dois anos** em RJ, PR e DF. Por isso a alíquota
interna **deixa de vir do ERP**.

### Separação de fontes — a decisão estrutural desta fatia

| fato | fonte | por quê |
|---|---|---|
| é ST neste destino? | `TGFICM.CODTRIB` espelhado | cadastro da empresa, mede 99,92% |
| `a_int` — interna do destino | **tabela legal nossa, 27 linhas, versionada** | o ERP está errado em 3 UFs |
| `a_inter` — interestadual | **derivado, sem tabela** | `ORIGPROD ∈ {1,2,3,8}` → 4%; UF ∈ {SP,RJ,PR,SC,RS} → 12%; senão 7% |

O `a_inter` foi confirmado duas vezes: empiricamente (26/26 combinações) e contra a matriz
origem×destino publicada. Saímos só de MG — é uma linha, não uma matriz.

## Fora de escopo — nomeado para ninguém improvisar

| fatia | conteúdo |
|---|---|
| P2.c | ligação pedido↔nota, imposto **realizado**, leitura as-of, detecção de divergência |
| P3 | `/anuncios` e `/mercado` consomem o mesmo cálculo (pior caso ex-ante) |
| P4 | mata os sítios restantes do 4% no simulador |

Corrigir cadastro do ERP não é de nenhuma fatia. Vai ao contador.

## Decisões do operador que governam este plano

- **D-36 — reconciliação é fatia própria (P2.c).** P2.b não depende do vínculo pedido↔nota, que
  hoje tem 11% de identidade (`AD_NUMPEDIDO_ECOM`: 8 linhas em 711.992, gravação parou em 21/07).
- **D-37 — UF sem linha na matriz vira pendência explícita.** Componente nulo + motivo nomeado
  (ADR-17). Nunca cai na alíquota interna como aproximação. Adivinhar já fabricou veredito errado
  uma vez nesta missão.
- **D-40 — custo é `CUSSEMICM`.** Confirmado pelo operador em 2026-08-03. `TGFCUS` tem sete colunas
  de custo e `CUSMEDSEMICMCS` **não é uma delas**; "custo médio sem ICMS" da tela do Sankhya é
  `CUSSEMICM`, que o espelho já lê. **Nenhuma mudança.** Isto derruba explicitamente a §8.3 do
  `sankhya-simulador-preco.md` do MNOS, que mandava trocar por `CUSVARIAVEL` — a diferença chega a
  7% (`22467`: 691,13 × 692,41; `39587`: 409,68 × 435,90 se fosse `CUSMEDICM`).
- **D-41 — DIFAL pela LEI, não pelo ERP.** Gross-up por dentro. Consequência aceita: nossa tela
  discorda da tela de Análise de Rentabilidade do Sankhya por desenho, em 2% a 10%.
- **D-42 — alíquota interna vem de tabela nossa, versionada, semeada da legislação.** Não do
  `TGFICM`. Cada linha carrega `fonte` e `lei`. Reconciliação contra nota emitida é o instrumento de
  verificação, e é rotina — foi assim que a defasagem do RJ apareceu.
- **D-43 — FCP entra embutido na alíquota interna, sem coluna própria.** Só RJ (2%), PB (2%),
  AL (1%) e SE (1%) aplicam FCP de forma geral; nos demais 23 é lista fechada de supérfluo (bebida,
  fumo, cosmético). Nossos produtos — metal sanitário e papeleiro — **não estão em lista nenhuma**.
  PB fica sem FCP porque lá ele não é geral (Dec. 25.618/2004). Medição de apoio: nota do RJ trouxe
  FCP 2% e reconciliou 21/21 ao centavo.

## Regras não negociáveis

1. **Nunca escreva no Oracle.** `METALPRD` é somente leitura.
2. **Desconhecido nunca vira zero nem valor plausível** (ADR-17). Componente sem resposta carrega
   motivo nomeado.
3. **Nenhum teste verde sem o vermelho provado antes.** Todo teste aqui tem controle negativo
   nomeado; rodar e ver falhar **é um passo**, não é opcional.
4. **`tenant_id` em toda query multi-tenant.**
5. **OpenAPI e `sdk-runtime` no mesmo commit.** Nunca um sem o outro.
6. **Nunca `git push`** sem autorização explícita do operador.
7. **Teste Go sempre com `GOCACHE` absoluto**, de dentro de `apps/server_core`.
8. **Chip nunca sobe servidor, porta ou `.env`.** Precisa disso = `REQUEST` ao hub.

## Comunicação com o hub

| situação | evento |
|---|---|
| terminou task e commitou | `COMMITTED` com o SHA |
| terminou o plano | `CLOSED` com SHA final e o que foi medido |
| bateu em algo não previsto | `ESCALATION` — **não improvise** |
| precisa de dependência, `.env`, servidor ou porta | `REQUEST` |
| task grande demais | `SPLIT-REQUEST` |
| travado | `BLOCKED` com o que já tentou |

**Dúvida é evento, não adivinhação.** Plano e código discordando: o código é o fato, o plano é a
alegação — `ESCALATION` com `arquivo:linha`.

---

## Estrutura de arquivos

**Criados:**

| arquivo | responsabilidade |
|---|---|
| `apps/server_core/migrations/0093_icms_matrix.sql` | 5 colunas fiscais em `products_mirror` + `icms_aliquota_interna` + `icms_matrix_mirror` |
| `apps/server_core/migrations/0094_icms_aliquota_interna_seed.sql` | as 27 linhas legais, com `fonte` e `lei` (D-42/D-43) |
| `internal/modules/pricing/domain/icms.go` | `ICMSCell`, `TaxComponents`, `TaxesForItem` — puro |
| `internal/modules/pricing/ports/tax_matrix.go` | `TaxMatrixReader` |
| `internal/modules/pricing/adapters/postgres/matrix_reader.go` | célula **vigente** + alíquota interna legal vigente |
| `internal/modules/internal_read/domain/icms_matrix.go` | `MatrixCell` + `ResolveCell` (lista branca), puro |
| `internal/modules/internal_read/adapters/oracle/icms_matrix.go` | extração do `TGFICM` |
| `internal/modules/internal_read/adapters/mirror/icms_matrix_writer.go` | escrita versionada |

**Modificados:**

| arquivo | mudança |
|---|---|
| `internal/modules/internal_read/adapters/oracle/sync.go` (`sankhyaBaseSQL`) | Q1 lê `GRUPOICMS` e `ORIGPROD`; query nova do bloco de ST (`TGFEFDVMRSTDIA`) |
| `internal/modules/internal_read/adapters/mirror/writer.go` (`upsertSQL`) | as 5 colunas fiscais no insert e no `DO UPDATE SET` |
| `internal/modules/pricing/domain/difal.go` | passa a **derivar** do `a` único (gross-up da lei); a tabela de internas hardcoded **sai** — vem de `icms_aliquota_interna` |
| `internal/modules/pricing/domain/decompose.go` | `Decomposition` ganha `ICMSSaida`, `PisCofins` e `RestituicaoST` |
| `internal/modules/orders/ports/tax_reader.go` | `TaxesForOrder` → `TaxesForItems`; `OrderTaxes` ganha os componentes novos |
| `internal/modules/orders/adapters/pricingtax/reader.go` | consulta a matriz; **fica vivo** |
| `internal/modules/orders/application/enrich_service.go:390,405` | `resolveTaxes` passa itens, não só total |
| `internal/composition/root.go:618-621` | fiação do `TaxMatrixReader` no reader existente |
| `contracts/api/marketplace-central.openapi.yaml` | componentes novos em `decomposicao` |
| `packages/sdk-runtime/src/index.ts` | tipos correspondentes |
| `apps/web/src/pages/pedidos/PedidoDrawer.tsx:155-160` | três linhas novas na `DecomposicaoSection`, uma delas **positiva** |

**Removidos:** nenhum.

---

## Task 1 — Migração: `grupo_icms` e `icms_matrix_mirror`

**Files:**
- Create: `apps/server_core/migrations/0093_icms_matrix.sql`
- Test: `apps/server_core/migrations/icms_matrix_test.go`

A última migração é `0092_listing_variations_tenant_provider_listing_index.sql`. Confirme antes de
numerar.

```sql
ALTER TABLE products_mirror ADD COLUMN IF NOT EXISTS grupo_icms          integer;
ALTER TABLE products_mirror ADD COLUMN IF NOT EXISTS origprod            smallint;
ALTER TABLE products_mirror ADD COLUMN IF NOT EXISTS st_retido_entrada   numeric(14,4);
ALTER TABLE products_mirror ADD COLUMN IF NOT EXISTS restituicao_unit    numeric(14,4);
ALTER TABLE products_mirror ADD COLUMN IF NOT EXISTS fiscal_dt_ref       date;

-- Tabela LEGAL, nossa, semeada da legislação. NÃO vem do TGFICM (D-42).
-- FCP já embutido onde é geral (D-43). Uma linha por UF por vigência.
CREATE TABLE icms_aliquota_interna (
    uf             text        NOT NULL,
    aliquota       numeric(6,3) NOT NULL,
    fcp_embutido   numeric(6,3) NOT NULL DEFAULT 0,
    fonte          text        NOT NULL,
    lei            text        NOT NULL,
    vigente_desde  date        NOT NULL,
    vigente_ate    date,
    PRIMARY KEY (uf, vigente_desde)
);

CREATE UNIQUE INDEX icms_aliquota_interna_vigente
    ON icms_aliquota_interna (uf)
    WHERE vigente_ate IS NULL;

CREATE TABLE icms_matrix_mirror (
    tenant_id          uuid        NOT NULL,
    uf_origem          text        NOT NULL,
    uf_destino         text        NOT NULL,
    grupo_icms         integer     NOT NULL,
    codtrib            smallint,
    aliquota           numeric(6,3),
    red_base           numeric(6,3),
    perc_fcp           numeric(6,3),
    linhas_candidatas  smallint    NOT NULL,
    ambiguo            boolean     NOT NULL,
    vigente_desde      timestamptz NOT NULL,
    vigente_ate        timestamptz,
    synced_at          timestamptz NOT NULL,
    PRIMARY KEY (tenant_id, uf_origem, uf_destino, grupo_icms, vigente_desde)
);

CREATE UNIQUE INDEX icms_matrix_mirror_vigente
    ON icms_matrix_mirror (tenant_id, uf_origem, uf_destino, grupo_icms)
    WHERE vigente_ate IS NULL;
```

`codtrib` é **nullable** — medido em `f1097db2`. `red_base` entra porque é ≠ 0 em 63% das linhas e
compõe a **base**, não a alíquota.

⚠️ Do `icms_matrix_mirror` esta fatia usa **só o `codtrib`** (é ST neste destino?). `aliquota` e
`aliqintdest` continuam sendo gravadas para P2.c poder reconciliar nota emitida, mas **o cálculo não
as lê** (D-42).

O índice parcial é o que impede duas versões abertas da mesma célula. Ele é a garantia; a disciplina
do escritor não é.

**Seed da `icms_aliquota_interna`** — 27 linhas, `vigente_desde` = data da lei, FCP embutido onde é
geral. Verificadas contra a legislação em 2026-08-03:

| UF | alíq. | FCP | UF | alíq. | FCP | UF | alíq. | FCP |
|---|---|---|---|---|---|---|---|---|
| AC | 19,0 | — | MA | 23,0 | — | RJ | **22,0** | 2,0 |
| AL | **21,5** | 1,0 | MT | 17,0 | — | RN | 20,0 | — |
| AM | 20,0 | — | MS | 17,0 | — | RO | 19,5 | — |
| AP | 18,0 | — | MG | 18,0 | — | RR | 20,0 | — |
| BA | 20,5 | — | PA | 19,0 | — | RS | 17,0 | — |
| CE | 20,0 | — | PB | 20,0 | — | SC | 17,0 | — |
| DF | **20,0** | — | PR | **19,5** | — | SP | 18,0 | — |
| ES | 17,0 | — | PE | 20,5 | — | SE | **20,0** | 1,0 |
| GO | 19,0 | — | PI | 22,5 | — | TO | 20,0 | — |

Datas de vigência que importam: **AL 01/04/2026** (Lei 9.776/2025, 20,5 + 1 FECOEP) · **RJ
20/03/2024** (Lei 10.253/2023 + LC 210/2023, 20 + 2 FECP) · **PR 18/03/2024** (Lei 21.850/2023) ·
**DF 21/01/2024** (Lei 7.326/2023).

As quatro em negrito são exatamente onde o `TGFICM` está errado ou onde a tabela comercial mais
citada (brasilnfe) está errada. **Nenhuma fonte isolada acertou as quatro** — por isso `fonte` e
`lei` são colunas, não comentário.

**Steps:**
- [ ] Teste primeiro: aplica a migração e tenta inserir **duas** linhas com `vigente_ate IS NULL`
      para a mesma célula; espera violação de unicidade.
- [ ] Rodar e **ver falhar** (tabela não existe). Registrar a saída.
- [ ] Escrever a migração.
- [ ] Rodar e ver passar.
- [ ] Segundo teste: `grupo_icms` existe em `products_mirror` e aceita `NULL`.

**Controle negativo:** remova o `WHERE vigente_ate IS NULL` do índice e prove que o teste de
duas-versões-abertas passa a aceitar a segunda linha.

---

## Task 2 — `grupo_icms` atravessa o sync de produtos

**Files:**
- Modify: `internal/modules/internal_read/adapters/oracle/sync.go` (`sankhyaBaseSQL`, `readBase`)
- Modify: `internal/modules/internal_read/adapters/mirror/writer.go` (`upsertSQL`) e a `Row`
- Test: integração do writer

`sankhyaBaseSQL` hoje seleciona `p.CODPROD, p.DESCRPROD, p.NCM, p.REFERENCIA, p.REFFORN,
p.CODGRUPOPROD, g.DESCRGRUPOPROD, m.DESCRICAO, p.USOPROD, p.AD_ECOMMERCE`. Acrescente `p.GRUPOICMS`.

**Atenção:** `GRUPOICMS` é o grupo **fiscal**, distinto de `CODGRUPOPROD` (grupo comercial, já
lido). Não reaproveite um pelo outro — nomes parecidos, colunas diferentes, erro barato de cometer e
caro de achar.

**Steps:**
- [ ] Teste de writer: `Row` com `GrupoICMS = 7` → `grupo_icms = 7`; `Row` com `GrupoICMS` nulo →
      coluna `NULL`, **não** `0`.
- [ ] Rodar e ver falhar.
- [ ] `GRUPOICMS` no `SELECT`, no scan de `readBase`, no struct `Row`, no `INSERT` e no
      `DO UPDATE SET`.
- [ ] Rodar e ver passar.

**Controle negativo:** faça o scan gravar `0` quando o Oracle devolve `NULL` e prove que o teste
reprova. Grupo desconhecido não é grupo zero.

---

## Task 3 — Extração e escrita versionada da matriz

**Files:**
- Create: `internal/modules/internal_read/domain/icms_matrix.go` (`ResolveCell`, puro)
- Create: `internal/modules/internal_read/adapters/oracle/icms_matrix.go`
- Create: `internal/modules/internal_read/adapters/mirror/icms_matrix_writer.go`
- Modify: `internal/modules/internal_read/adapters/oracle/sync.go` — `ORIGPROD` e o bloco de ST
- Test: unitário de `ResolveCell` + integração do writer

### Três extrações, não uma

| campo | fonte | destino |
|---|---|---|
| `codtrib` (é ST?) | `TGFICM` por `UF × grupo` | `icms_matrix_mirror.codtrib` |
| `origprod` | `TGFPRO.ORIGPROD` | `products_mirror.origprod` |
| `S` (ST retido) | `TGFEFDVMRSTDIA`, `TIPIMPOSTO = 'S'` | `products_mirror.st_retido_entrada` |
| `R` (restituição) | `TGFEFDVMRSTDIA`, `TIPIMPOSTO IN ('S','I')` | `products_mirror.restituicao_unit` |
| data do bloco de ST | `MAX(DTMOV)` do mesmo bloco | `products_mirror.fiscal_dt_ref` |

Query do bloco de ST, já medida contra o Oracle:

```sql
WITH v AS (
  SELECT d.CODPROD, d.TIPIMPOSTO, d.VLRUNITMED, d.DTMOV,
         ROW_NUMBER() OVER (PARTITION BY d.CODPROD, d.TIPIMPOSTO
                            ORDER BY d.DTMOV DESC) AS RN
  FROM METALPRD.TGFEFDVMRSTDIA d
  WHERE d.CODEMP = 1
    AND d.TIPIMPOSTO IN ('S', 'I')
    AND d.DTMOV <= SYSDATE
    AND NVL(d.VLRUNITMED, 0) <> 0
)
SELECT CODPROD,
       MAX(CASE WHEN TIPIMPOSTO = 'S' THEN VLRUNITMED END)  AS S_ST,
       SUM(VLRUNITMED)                                       AS R_TOTAL,
       MAX(DTMOV)                                            AS DT_REF
FROM v WHERE RN = 1 GROUP BY CODPROD
```

`NVL(VLRUNITMED,0) <> 0` **antes** do `ROW_NUMBER` — sem isso a linha mais recente pode ser um zero
que enterra o valor real. `CODEMP = 1` fixo é **D-17**, não decisão nova.

**Resolução por lista branca** (spec §4.3 — a precedência formal entre linhas concorrentes do
`TGFICM` **não foi estabelecida**: houve uma única observação, insuficiente para generalizar).

Aceite apenas `N/0` (sem restrição), `I/<grupo do produto>` e `S/−1` (curinga). Descarte `O`, `P`,
`H`, `K`, `G`, `L`, `X`, `T`.

- uma linha sobrevive → grava a célula, `ambiguo = false`
- mais de uma → grava `ambiguo = true`, `linhas_candidatas = N`, **e não escolhe**
- nenhuma → não grava célula

**Armadilhas do DDL, já medidas** (`f1097db2`, memória `tgficm-ddl-traps`):
- `ALIQUOTA` é a alíquota **da operação**. `ALIQUFDEST` é do bloco de ST e é ~2,5× maior — a coluna
  errada não estoura teste nenhum, só produz número plausível e errado.
- `REDBASE` compõe a **base**, não a alíquota.
- UF em `TGFICM` **não é código IBGE**; MG = 13. Traduzir exige `TSIUFS` com `CODPAIS = 55`.

Sob **D-42** a escolha entre `ALIQUOTA` e `ALIQUFDEST` deixou de carregar o cálculo — grava as duas,
o cálculo não lê nenhuma. O que carrega peso aqui é **só o `codtrib`**, e ele mede 99,92% contra as
notas emitidas. Não relaxe a resolução por lista branca por causa disso: célula ambígua ainda
significa `codtrib` incerto, e `codtrib` incerto é a diferença entre 0% e 24% de ICMS.

**Escrita versionada:** compare com a versão aberta. Igual → só `synced_at`. Diferente →
`vigente_ate = now()` na aberta e `INSERT` da nova. Sumiu do ERP → fecha e não abre.

**Steps:**
- [ ] Teste de `ResolveCell`: três candidatas (`N/0`, `I/7`, `O/…`) para o mesmo destino → `O`
      descartada, duas sobrevivem, `ambiguo = true` sem escolha.
- [ ] Rodar e ver falhar.
- [ ] Implementar `ResolveCell`.
- [ ] Teste do writer: grava célula; roda de novo com `aliquota` diferente → **duas** linhas, a
      primeira com `vigente_ate` preenchido, a segunda aberta; terceira rodada com o mesmo valor →
      continua duas linhas, só `synced_at` avança.
- [ ] Rodar e ver falhar.
- [ ] Implementar writer e adapter Oracle.

**Controle negativo:** faça o writer sempre inserir versão nova e prove que a terceira rodada (valor
idêntico) passa a criar uma terceira linha — o teste tem que pegar.

---

## Task 4 — `pricing/domain`: ICMS de saída e PIS/COFINS

**Files:**
- Create: `internal/modules/pricing/domain/icms.go`
- Modify: `internal/modules/pricing/domain/decompose.go` (`Decomposition`, `Decompose`)
- Create: `internal/modules/pricing/ports/tax_matrix.go`
- Create: `internal/modules/pricing/adapters/postgres/matrix_reader.go`
- Test: unitários puros + golden

Domínio puro, sem I/O — a célula chega já resolvida pela porta, igual ao que `DecomposeInput` já faz
com comissão, frete e custo.

`Decomposition` ganha `ICMSSaida *string` e `PisCofins *string`. Ambos ponteiro porque ambos podem
ser desconhecidos; entram em `ComponentesDesconhecidos` quando nulos, como `frete`, `difal`,
`tarifa_full` e `custo` já entram.

`Imposto` (alíquota de regime) **permanece no struct** — matá-lo é P4 e removê-lo agora quebra sete
sítios fora do escopo. Mas quando a célula responde, ele **deixa de ser somado**. Registre isso no
comentário do campo: componente somado em silêncio duas vezes é como o 4% sobreviveu até aqui.

### A fórmula — uma alíquota efetiva, não dois componentes independentes

ICMS de operação e DIFAL **telescopam**: a soma é sempre a base × interna do destino. Colapsa em um
parâmetro `a`.

```
a_inter = 0,04  se ORIGPROD ∈ {1,2,3,8}
          0,12  se UF destino ∈ {SP, RJ, PR, SC, RS}
          0,07  demais UFs

a = 0                                      produto ST neste destino (codtrib = 60)
    0,18                                   MG, não-ST
    a_int × (1 − a_inter) / (1 − a_int)    fora de MG, não-ST          [D-41, gross-up da lei]

ICMS_total  = P × a
BASE_PC     = MAX(0, P × (1 − a) − S)
PIS/COFINS  = 0,0925 × BASE_PC
R_aplicável = R  se UF ≠ MG,  senão 0
```

**Exibição.** A tela continua com `ICMS saída` e `DIFAL` separados — deriva os dois do mesmo `a`:
`ICMS_oper = P × a_inter` e `DIFAL = P × a − ICMS_oper`. **Uma fonte, duas linhas.** Nunca dois
cálculos.

Regras:
- produto **ST no destino**: `a = 0`; ICMS de saída e DIFAL = **`0` explícito** com motivo
  `"ST — próprio já no custo"`. Zero legítimo, não desconhecido.
- **restituição só sai de MG.** Gatilho literal `UF <> 13` lido na fonte de
  `PAN_GET_CUSVAR_MNOBRE`. Intra-MG `R = 0` — e o único benefício do ST embutido lá é a redução da
  base de PIS/COFINS, que vale `0,0925 × S`, **não** `S`.
- **`S` e `R` são grandezas diferentes.** `S` = `TIPIMPOSTO='S'` (só ST, entra na base do
  PIS/COFINS). `R` = `'S'` + `'I'` (ST **e** ICMS próprio, abate o custo). Usar um no papel do outro
  é o erro que a §9 do doc do MNOS cometeu e que inflou o piso em ~25%.
- `MAX(0, …)` **não é cosmético**: há itens em que o ST retido excede o líquido e o ERP trava a base
  em zero. Sem o `MAX` sai base negativa e crédito fantasma.
- célula ausente, `ambiguo = true`, ou UF sem linha (**D-37**): `nil` + motivo nomeado.
- **`S` ou `R` com `fiscal_dt_ref` velho não vira zero.** Medido: o produto `15956` — o mais
  vendido, 10 dos 38 pedidos — tem `S` de **2022-02-04**, quatro anos parado, enquanto os outros
  cinco estão em 2026-07-31. Componente sai com aviso de idade, nunca silencioso.

**Steps:**
- [ ] Golden: um item interno, um interestadual, um ST, um com célula ambígua. Quatro saídas
      esperadas escritas **antes** do código.
- [ ] Rodar e ver falhar.
- [ ] Implementar `TaxesForItem` e estender `Decompose`.
- [ ] Rodar e ver passar.
- [ ] Invariante já existente: com **todos** os componentes conhecidos,
      `preço = Σ(componentes) + margem_valor` **exatamente**. Os dois novos entram na soma; se
      quebrar, é porque `Imposto` está sendo somado junto com `ICMSSaida`.

**Controle negativo:** faça o ramo ST devolver `nil` em vez de `0` e prove que o golden reprova. Os
dois são diferentes e a tela mostra coisas diferentes.

---

## Task 5 — Porta por item, fiação e `/pedidos`

**Files:**
- Modify: `internal/modules/orders/ports/tax_reader.go`
- Modify: `internal/modules/orders/adapters/pricingtax/reader.go`
- Modify: `internal/modules/orders/application/enrich_service.go:390,405`
- Modify: `internal/composition/root.go:618-621`

`TaxesForOrder(ctx, total float64, destinoUF string)` vira
`TaxesForItems(ctx, items []TaxItem, destinoUF string)`, com `TaxItem` carregando `codigo_produto`,
valor da linha e quantidade. O adapter resolve `grupo_icms`, `origprod`, `st_retido_entrada` e
`restituicao_unit` por produto, consulta a célula vigente e soma por pedido.

**A soma é por item, sempre.** `MAX(0, …)` na base do PIS/COFINS e o gross-up do DIFAL são ambos
não-lineares — calcular no total do pedido e ratear depois dá número diferente de calcular por linha
e somar. Um pedido com dois grupos é a prova; está no golden da Task 4.

`restituicao_unit` × quantidade entra como componente **positivo** na decomposição, e só quando
`destinoUF ≠ 'MG'`.

`root.go:618` hoje é
`orderspricingtax.NewReader(pricingpostgres.NewCalcRepository(pool), cfg.DefaultTenantID)`.
Acrescente o `TaxMatrixReader` como segundo colaborador. O reader **não** é substituído.

**Steps:**
- [ ] Teste do adapter: pedido com dois itens de grupos diferentes (um ST, um normal) → ICMS de
      saída = só o do item normal.
- [ ] Rodar e ver falhar.
- [ ] Reescrever porta e adapter; ajustar `resolveTaxes`.
- [ ] Rodar e ver passar.
- [ ] `go build ./...` limpo.

**Controle negativo:** aplique a alíquota do primeiro item ao total do pedido e prove que o teste de
dois-grupos reprova. É exatamente o erro que a porta por pedido cometia.

---

## Task 6 — Contrato e tela

**Files:**
- Modify: `contracts/api/marketplace-central.openapi.yaml`
- Modify: `packages/sdk-runtime/src/index.ts`
- Modify: `apps/web/src/pages/pedidos/PedidoDrawer.tsx:155-160`

**No mesmo commit.** OpenAPI sem SDK, ou SDK sem OpenAPI, é violação.

A tela **já** tem a `DecomposicaoSection` renderizando `decomposicao.*` com `DecompRow`,
`UnknownValue` e `componentes_desconhecidos`. São **três linhas novas** — `ICMS saída`, `PIS/COFINS`
e `Restituição ST` — na lista que já existe. Não crie componente novo.

`Restituição ST` é a única linha **positiva** da seção. Se a `DecomposicaoSection` assume que todo
componente subtrai (sinal, cor, ou a soma de verificação), isso quebra — trate no teste, não na
revisão. Intra-MG a linha sai `0` explícito com motivo `"restituição só em operação interestadual"`,
nunca ausente.

Mais duas coisas:
- o `hint` de cada linha nova nomeia **qual lacuna preencher no Sankhya** quando o valor é nulo.
  "Desconhecido" sem endereço não é honesto, é mudo.
- os `hint` hardcoded `"(hub C2)"` nas linhas que agora recebem valor de verdade estão mentindo.
  Corrija os que esta fatia passa a alimentar; deixe os outros.

**Steps:**
- [ ] Campos novos no OpenAPI.
- [ ] Tipos no `sdk-runtime`.
- [ ] `tsc` limpo em `apps/web`.
- [ ] Três `DecompRow` novas + hints corrigidos.
- [ ] Teste de front: componente desconhecido renderiza `UnknownValue` com o hint que nomeia a
      lacuna, **não** `"—"` hardcoded.

**Controle negativo:** hardcode `"—"` numa das linhas novas e prove que o teste reprova.

---

## Task 7 — Verificação ao vivo e fechamento

**Chip não sobe servidor.** `REQUEST` ao hub para o live drive.

### Alvos numéricos — o dry-run já rodou contra os 38 pedidos reais

A lógica desta fatia foi executada em SQL contra o dev stack antes de existir código
(`.mnfs/MIS-007-ml-sync/M-06-orders-backfill-decomposition/_evidence/p2b-dryrun/calc2.sql`). O
resultado é o **alvo de aceite**, não uma expectativa:

```
pedidos | nao_calculavel | positivos_hoje | positivos_novo | viram_negativo | soma_hoje | soma_nova
     38 |              7 |             33 |             28 |              3 |   4560.40 |   1927.01
```

Controle positivo — pedido `2000017515486360` (BA, produto `41912`, `origprod = 0`):

| componente | valor |
|---|---|
| `a_ef` | 23,98% |
| ICMS total | 71,92 |
| PIS/COFINS | 19,00 |
| restituição | +37,69 |
| custo | 154,53 |
| comissão | 40,49 |
| frete | 23,65 |
| **margem nova** | **+28,00** |

Os **7 não-calculáveis** são: 5 pedidos sem vínculo de produto, 2 porque o grupo `311` não tem
célula de `TGFICM` para RJ. Os dois últimos **têm que aparecer como pendência nomeada** — se
saírem com número, a fatia falhou o ADR-17.

**Steps:**
- [ ] Sync da matriz rodado; contar células gravadas, quantas `ambiguo = true`, quantas UFs sem
      linha. Números no pack.
- [ ] Rodar a agregação acima **contra o código**, não contra o SQL do dry-run, e bater os sete
      números. Divergência é defeito ou é achado — nos dois casos, escrito no pack antes do `CLOSED`.
- [ ] Drawer do pedido `2000017515486360`: `+28,00`, com as três linhas novas visíveis. Screenshot.
- [ ] Drawer de um pedido **interno** (MG): ICMS de saída `0` com motivo ST **ou** 18%, DIFAL `0`
      explícito, restituição `0` com motivo. Screenshot.
- [ ] Um dos 2 pedidos de grupo `311`/RJ: pendência nomeada (**D-37**), não alíquota interna.
- [ ] Um item do produto `15956`: aviso de idade do bloco de ST (`fiscal_dt_ref = 2022-02-04`)
      visível. Se sair silencioso, é **D-44** aberta, não `CLOSED`.
- [ ] `CLOSED` ao hub com SHA final e os números medidos.

---

## Dívidas registradas por esta fatia

| id | dívida |
|---|---|
| D-17 | `CODEMP = 1` fixo no caminho de custo |
| D-28 | histórico da matriz começa no primeiro sync; a defasagem da BA (17,0% → 20,5% em 20/07/2026) **não é recuperável** do nosso espelho — está em `TGFHICM`, que esta fatia não lê |
| D-38 | `Imposto` (alíquota de regime) continua no struct e nos sítios do simulador; morte é P4 |
| D-39 | crédito PIS/COFINS de 9,25% é hardcoded no `CUSSEMICM` do ERP; espelhamos o número, não a regra |
| D-44 | bloco de ST envelhece sem sinal na origem. Produto `15956` — o mais vendido, 10 dos 38 pedidos — tem `DTMOV = 2022-02-04`, quatro anos parado, contra 2026-07-31 dos outros cinco. Esta fatia **expõe** a idade; não define política de expiração |
| D-45 | `icms_aliquota_interna` é semeada à mão e **não tem fonte automática** — nenhum órgão publica tabela consultável de alíquota interna/FCP (CONFAZ, SEFAZs e Portal NF-e não publicam; IBPT é carga total da lei da transparência, não serve p/ DIFAL). Mudança de alíquota estadual chega por migração, e o atraso é a janela de erro |
| D-46 | `a_inter` é derivado de `ORIGPROD ∈ {1,2,3,8}` → 4%, sem consultar FCI nem conteúdo de importação por operação. Correto para o catálogo medido; falso se entrar produto com conteúdo de importação entre 40% e 70% resolvido caso a caso |
