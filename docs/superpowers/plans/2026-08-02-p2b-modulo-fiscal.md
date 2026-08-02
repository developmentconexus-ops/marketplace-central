# P2.b — Módulo `fiscal` Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Trocar o imposto inventado de `/pedidos` (perfil com zero linhas caindo em "SIMPLES 4%") pelo imposto calculado a partir da matriz de ICMS do ERP, espelhada e versionada.

**Architecture:** Um módulo `fiscal` novo — domínio puro sem I/O, porta, adapter Postgres — consumido por `orders` agora e por `listings`/`pricing` nas fatias P3/P4. A matriz do ERP é espelhada em `icms_matrix_mirror` com versão (SCD-2) e consultada as-of pela data do pedido, para que corrigir o cadastro do ERP não apague a divergência histórica.

**Tech Stack:** Go 1.x (`apps/server_core`), Postgres via `pgxpool`, Oracle via `godror` (somente leitura), React + TypeScript (`apps/web`), OpenAPI + `packages/sdk-runtime` escrito à mão.

**Spec:** [`docs/superpowers/specs/2026-08-02-p2b-modulo-fiscal-design.md`](../specs/2026-08-02-p2b-modulo-fiscal-design.md)

---

## Comunicação com o hub — leia antes de começar

Você **pode e deve** mandar mensagem para a sessão do hub. A sessão anterior não fez isso e custou retrabalho.

Mande evento ao hub quando:

| situação | evento |
|---|---|
| terminou uma task e commitou | `COMMITTED` com o SHA |
| terminou o plano inteiro | `CLOSED` com o SHA final e o resumo do que foi medido |
| bateu em algo que o plano não previu | `ESCALATION` — **não improvise** |
| precisa de dependência nova, `.env`, servidor ou porta | `REQUEST` — chip nunca sobe servidor nem cria `.env` |
| descobriu que uma task é grande demais | `SPLIT-REQUEST` |
| está travado | `BLOCKED` com o que já tentou |

**Dúvida é evento, não adivinhação.** Se o plano e o código discordarem, o código é o fato e o plano é a alegação — mande `ESCALATION` com `arquivo:linha`.

## Regras não negociáveis desta fatia

1. **Nunca escreva no Oracle.** Acesso a `METALPRD` é somente leitura.
2. **Desconhecido nunca vira zero nem valor plausível** (ADR-17). Componente sem resposta carrega `Motivo`.
3. **Nenhum teste verde sem o vermelho provado antes.** Todo teste deste plano tem um controle negativo nomeado. Rodar o controle negativo e ver falhar **é um passo**, não é opcional.
4. **`tenant_id` em toda query multi-tenant.**
5. **OpenAPI e `sdk-runtime` no mesmo commit.** Nunca um sem o outro.
6. **Nunca `git push`** sem autorização explícita do operador.
7. **Comando de teste Go sempre com `GOCACHE` absoluto**, de dentro de `apps/server_core`.

---

## Estrutura de arquivos

**Criados:**

| arquivo | responsabilidade |
|---|---|
| `apps/server_core/migrations/0093_fiscal_icms_matrix.sql` | `grupo_icms` no espelho de produtos + tabela `icms_matrix_mirror` |
| `internal/modules/fiscal/domain/rules.go` | tipos: `Componente`, `ICMSRule`, `Resultado` |
| `internal/modules/fiscal/domain/calc.go` | `TaxesForValue` — a fórmula, pura, sem I/O |
| `internal/modules/fiscal/domain/money.go` | arredondamento em centavos, meia-para-cima, exato |
| `internal/modules/fiscal/ports/ports.go` | `TaxMatrixReader` |
| `internal/modules/fiscal/adapters/postgres/matrix_reader.go` | leitura as-of dos dois espelhos |
| `internal/modules/internal_read/domain/icms_matrix.go` | `MatrixCell` + `ResolveCell` (lista branca), puro |
| `internal/modules/internal_read/adapters/oracle/icms_matrix.go` | extração da matriz do Oracle |
| `internal/modules/internal_read/adapters/mirror/icms_matrix_writer.go` | escrita versionada SCD-2 |
| `internal/modules/orders/adapters/fiscal/reader.go` | adapta `fiscal` à porta de `orders` |
| `apps/web/src/pages/pedidos/FiscalSection.tsx` | seção Fiscal do drawer |
| `apps/web/src/pages/integracoes/IcmsMatrixHealthCard.tsx` | card de saúde do espelho |

**Modificados:**

| arquivo | mudança |
|---|---|
| `internal/modules/internal_read/adapters/oracle/sync.go:103-111` | Q1 passa a ler `GRUPOICMS` |
| `internal/modules/internal_read/adapters/mirror/writer.go:74-95` | `grupo_icms` no upsert |
| `internal/modules/orders/ports/tax_reader.go` | porta redesenhada de pedido para **item** |
| `internal/modules/orders/application/enrich_service.go:405-423` | `resolveTaxes` vira `resolveItemTaxes` |
| `internal/composition/root.go:618-621` | fiação do reader novo |
| `contracts/api/marketplace-central.openapi.yaml` | `ValorFiscal` + `OrderFiscal` |
| `packages/sdk-runtime/src/index.ts` | tipos correspondentes |
| `apps/web/src/pages/pedidos/PedidoDrawer.tsx:133-180` | seção DIFAL legada sai, Fiscal entra |

**Removidos:**

| arquivo | motivo |
|---|---|
| `internal/modules/orders/adapters/pricingtax/reader.go` | substituído; é a fonte do 4% fabricado |
| `internal/modules/orders/adapters/pricingtax/reader_test.go` | idem |

---

## Task 1: Migração — coluna de grupo e tabela da matriz

**Files:**
- Create: `apps/server_core/migrations/0093_fiscal_icms_matrix.sql`
- Test: `apps/server_core/migrations/fiscal_icms_matrix_test.go`

- [ ] **Step 1: Escreva o teste que falha**

O diretório `apps/server_core/migrations` já tem testes Go que verificam schema (`sync_state_test.go`, `sellable_assortment_test.go`). Siga o mesmo estilo — abra um deles antes de escrever, para copiar o helper de conexão e a tag de build.

```go
//go:build integration

package migrations_test

import (
	"context"
	"testing"
)

func TestICMSMatrixMirrorSchema(t *testing.T) {
	ctx := context.Background()
	pool := openTestPool(t) // helper existente nos testes vizinhos
	defer pool.Close()

	var grupoICMS string
	err := pool.QueryRow(ctx, `
		SELECT data_type FROM information_schema.columns
		WHERE table_name = 'products_mirror' AND column_name = 'grupo_icms'`).Scan(&grupoICMS)
	if err != nil {
		t.Fatalf("products_mirror.grupo_icms ausente: %v", err)
	}

	wanted := []string{
		"tenant_id", "uf_origem", "uf_destino", "grupo_icms", "codtrib", "zerar",
		"aliquota", "redbase", "aliqintdest", "aliq_uf_dest", "perc_red_base_dest",
		"perc_fcp", "linhas_candidatas", "ambiguo",
		"vigente_desde", "vigente_ate", "synced_at",
	}
	for _, col := range wanted {
		var found string
		if err := pool.QueryRow(ctx, `
			SELECT column_name FROM information_schema.columns
			WHERE table_name = 'icms_matrix_mirror' AND column_name = $1`, col).Scan(&found); err != nil {
			t.Fatalf("icms_matrix_mirror.%s ausente: %v", col, err)
		}
	}

	// Uma célula só pode ter UMA versão aberta. Sem isso, o as-of devolve duas
	// linhas e a leitura vira loteria.
	var idx int
	if err := pool.QueryRow(ctx, `
		SELECT count(*) FROM pg_indexes
		WHERE tablename = 'icms_matrix_mirror'
		  AND indexdef ILIKE '%vigente_ate IS NULL%'`).Scan(&idx); err != nil {
		t.Fatalf("consulta de índice falhou: %v", err)
	}
	if idx == 0 {
		t.Fatal("falta índice único parcial de versão aberta em icms_matrix_mirror")
	}
}
```

- [ ] **Step 2: Rode e confirme que falha**

```bash
cd apps/server_core && GOCACHE=$(pwd)/.gocache go test -tags integration ./migrations/ -run TestICMSMatrixMirrorSchema -v
```

Esperado: `FAIL` com `products_mirror.grupo_icms ausente`.

- [ ] **Step 3: Escreva a migração**

```sql
-- 0093_fiscal_icms_matrix.sql
--
-- P2.b — imposto ex-ante pela matriz do ERP.
--
-- Duas coisas:
--   1. products_mirror ganha grupo_icms (TGFPRO.GRUPOICMS), a chave que liga o
--      produto à célula da matriz. Nullable: a fonte xlsx não tem esse campo.
--   2. icms_matrix_mirror — o espelho VERSIONADO da matriz de ICMS.
--
-- Por que versionado: a alíquota interna da Bahia ficou congelada em 17,0% por
-- 29 meses e foi corrigida para 20,5% em 20/07/2026. Se o cálculo usasse sempre
-- a matriz de hoje, um pedido de março passaria a exibir 20,5% e a divergência
-- desapareceria da tela no exato momento em que o ERP fosse consertado.
--
-- A célula gravada aqui já vem RESOLVIDA pelo sync (lista branca de restrição).
-- ambiguo = true significa que mais de uma linha sobreviveu à lista branca e o
-- sync se RECUSOU a escolher — a precedência entre linhas concorrentes não foi
-- estabelecida por medição. Leitor devolve desconhecido.

ALTER TABLE products_mirror
    ADD COLUMN IF NOT EXISTS grupo_icms INTEGER;

CREATE TABLE IF NOT EXISTS icms_matrix_mirror (
    tenant_id          UUID        NOT NULL,
    uf_origem          TEXT        NOT NULL,
    uf_destino         TEXT        NOT NULL,
    grupo_icms         INTEGER     NOT NULL,
    -- codtrib é NULLABLE: 22 linhas com origem MG têm CODTRIB NULL na TGFICM.
    -- 0 significa "tributado" e gravar 0 no lugar de NULL seria uma afirmação.
    codtrib            SMALLINT,
    -- zerar vem de TGFICM.ZERAR: 'S' zera o imposto da célula.
    zerar              BOOLEAN     NOT NULL DEFAULT false,
    -- aliquota é TGFICM.ALIQUOTA, a alíquota da OPERAÇÃO (7 / 12 / 4 / 18).
    aliquota           NUMERIC(6,3),
    -- redbase é TGFICM.REDBASE, redução da base da operação. Não-zero em 2.669
    -- de 4.209 linhas MG — ignorá-la erra o ICMS na maioria das células.
    redbase            NUMERIC(6,3),
    -- aliqintdest é o que o motor do ERP lê e o que reconciliou o DIFAL das
    -- notas reais. NULL em 82,6% das linhas MG.
    aliqintdest        NUMERIC(6,3),
    -- aliq_uf_dest é TGFICM.ALIQUFDEST. Pertence ao bloco de ST e NUNCA entra
    -- no cálculo. Espelhada porque tem 87,7% de cobertura e fica mais perto da
    -- alíquota legal: é o sinal de divergência da fatia P5.
    aliq_uf_dest       NUMERIC(6,3),
    -- perc_red_base_dest é TGFICM.PERCREDBASEDEST. Papel NÃO VERIFICADO —
    -- espelhada, nunca multiplicada.
    perc_red_base_dest NUMERIC(6,3),
    perc_fcp           NUMERIC(6,3),
    linhas_candidatas  SMALLINT    NOT NULL,
    ambiguo            BOOLEAN     NOT NULL,
    vigente_desde      TIMESTAMPTZ NOT NULL,
    vigente_ate        TIMESTAMPTZ,
    synced_at          TIMESTAMPTZ NOT NULL,
    PRIMARY KEY (tenant_id, uf_origem, uf_destino, grupo_icms, vigente_desde)
);

-- Exatamente uma versão aberta por célula. É o invariante que faz o as-of ser
-- determinístico.
CREATE UNIQUE INDEX IF NOT EXISTS icms_matrix_mirror_open_version
    ON icms_matrix_mirror (tenant_id, uf_origem, uf_destino, grupo_icms)
    WHERE vigente_ate IS NULL;

-- Caminho da consulta as-of.
CREATE INDEX IF NOT EXISTS icms_matrix_mirror_asof
    ON icms_matrix_mirror (tenant_id, uf_origem, uf_destino, grupo_icms, vigente_desde DESC);
```

- [ ] **Step 4: Rode e confirme que passa**

```bash
cd apps/server_core && GOCACHE=$(pwd)/.gocache go test -tags integration ./migrations/ -run TestICMSMatrixMirrorSchema -v
```

Esperado: `PASS`.

- [ ] **Step 5: Controle negativo — prove que o índice único morde**

Insira duas versões abertas da mesma célula direto no banco de teste e confirme que o Postgres recusa:

```sql
INSERT INTO icms_matrix_mirror
 (tenant_id, uf_origem, uf_destino, grupo_icms, codtrib, linhas_candidatas, ambiguo, vigente_desde, synced_at)
VALUES ('00000000-0000-0000-0000-000000000001','MG','BA',122,0,1,false, now(), now()),
       ('00000000-0000-0000-0000-000000000001','MG','BA',122,0,1,false, now() + interval '1 day', now());
```

Esperado: erro `duplicate key value violates unique constraint "icms_matrix_mirror_open_version"`. Se **não** der erro, o índice está errado — conserte antes de seguir.

- [ ] **Step 6: Commit**

```bash
git add apps/server_core/migrations/0093_fiscal_icms_matrix.sql apps/server_core/migrations/fiscal_icms_matrix_test.go
git commit -m "feat(fiscal): grupo_icms no espelho de produtos e icms_matrix_mirror versionado"
```

---

## Task 2: `grupo_icms` atravessa o sync de produtos

**Files:**
- Modify: `apps/server_core/internal/modules/internal_read/adapters/mirror/writer.go:74-95`
- Modify: `apps/server_core/internal/modules/internal_read/adapters/oracle/sync.go:103-111`
- Test: `apps/server_core/internal/modules/internal_read/adapters/mirror/writer_integration_test.go`

O campo existe em `TGFPRO` e a Q1 já lê essa tabela. É só carregar mais uma coluna até o Postgres.

- [ ] **Step 1: Escreva o teste que falha**

Acrescente ao arquivo de teste de integração do writer:

```go
func TestApplySnapshotPersistsGrupoICMS(t *testing.T) {
	ctx := context.Background()
	pool := openTestPool(t)
	defer pool.Close()

	tenantID := newTestTenant(t, pool) // helper existente
	writer := mirror.NewPgWriter(pool)

	grupo := 122
	err := writer.ApplySnapshot(ctx, tenantID, []mirror.Row{{
		CodigoProduto: "15956",
		Descricao:     ptr("PRODUTO TESTE"),
		GrupoICMS:     &grupo,
	}})
	if err != nil {
		t.Fatalf("ApplySnapshot: %v", err)
	}

	var got *int
	if err := pool.QueryRow(ctx,
		`SELECT grupo_icms FROM products_mirror
		  WHERE tenant_id = $1 AND source = 'sankhya' AND codigo_produto = '15956'`,
		tenantID).Scan(&got); err != nil {
		t.Fatalf("leitura: %v", err)
	}
	if got == nil || *got != 122 {
		t.Fatalf("grupo_icms = %v, queria 122", got)
	}
}
```

- [ ] **Step 2: Rode e confirme que falha**

```bash
cd apps/server_core && GOCACHE=$(pwd)/.gocache go test -tags integration ./internal/modules/internal_read/adapters/mirror/ -run TestApplySnapshotPersistsGrupoICMS -v
```

Esperado: `FAIL` — o campo `GrupoICMS` não existe em `mirror.Row`, então nem compila.

- [ ] **Step 3: Acrescente o campo à `Row`**

Em `writer.go`, no struct `Row`, junto dos outros campos opcionais:

```go
	// GrupoICMS é TGFPRO.GRUPOICMS: a chave que liga o produto à célula da
	// matriz de ICMS. Nil quando a fonte não tem o campo (xlsx) ou quando o
	// produto não tem grupo fiscal cadastrado — nos dois casos o imposto sai
	// desconhecido, nunca zero.
	GrupoICMS *int
```

- [ ] **Step 4: Acrescente ao upsert**

Em `writer.go:74-95`, o `INSERT` vai de 14 para 15 parâmetros de valor:

```go
const upsertSQL = `
INSERT INTO products_mirror
	(tenant_id, source, codigo_produto, descricao, referencia, ean, marca,
	 grupo_codigo, grupo_descricao, ncm, custo, preco_venda, usoprod, ad_ecommerce, estoque_total,
	 grupo_icms,
	 absent_in_last_snapshot, stale_since, updated_at)
VALUES ($1, 'sankhya', $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, false, NULL, now())
ON CONFLICT (tenant_id, source, codigo_produto) DO UPDATE SET
	descricao = EXCLUDED.descricao,
	referencia = EXCLUDED.referencia,
	ean = EXCLUDED.ean,
	marca = EXCLUDED.marca,
	grupo_codigo = EXCLUDED.grupo_codigo,
	grupo_descricao = EXCLUDED.grupo_descricao,
	ncm = EXCLUDED.ncm,
	custo = EXCLUDED.custo,
	preco_venda = EXCLUDED.preco_venda,
	usoprod = EXCLUDED.usoprod,
	ad_ecommerce = EXCLUDED.ad_ecommerce,
	estoque_total = EXCLUDED.estoque_total,
	grupo_icms = EXCLUDED.grupo_icms,
	absent_in_last_snapshot = false,
	stale_since = NULL,
	updated_at = now()`
```

Acrescente `row.GrupoICMS` na lista de argumentos, **na posição 15**, na chamada que executa esse SQL. Ordem errada aqui grava grupo no lugar de estoque em silêncio — confira contando os parâmetros.

- [ ] **Step 5: Rode e confirme que passa**

```bash
cd apps/server_core && GOCACHE=$(pwd)/.gocache go test -tags integration ./internal/modules/internal_read/adapters/mirror/ -run TestApplySnapshotPersistsGrupoICMS -v
```

Esperado: `PASS`.

- [ ] **Step 6: Q1 do Oracle passa a ler `GRUPOICMS`**

Em `sync.go:103-111`:

```go
const sankhyaBaseSQL = `
SELECT p.CODPROD, p.DESCRPROD, p.NCM, p.REFERENCIA, p.REFFORN,
       p.CODGRUPOPROD, g.DESCRGRUPOPROD, m.DESCRICAO,
       p.USOPROD, p.AD_ECOMMERCE, p.GRUPOICMS
FROM METALPRD.TGFPRO p
LEFT JOIN METALPRD.TGFGRU g ON g.CODGRUPOPROD = p.CODGRUPOPROD
LEFT JOIN METALPRD.TGFMAR m ON m.CODIGO       = p.CODMARCA
WHERE p.ATIVO = 'S'`
```

Na função `readBase`, acrescente `grupoICMS sql.NullInt64` às variáveis de scan, na **última posição** do `Scan` (a ordem do `Scan` tem que bater com a ordem do `SELECT`), e preencha `r.row.GrupoICMS` quando `grupoICMS.Valid`.

- [ ] **Step 7: Rode a suíte do adapter Oracle**

```bash
cd apps/server_core && GOCACHE=$(pwd)/.gocache go test ./internal/modules/internal_read/adapters/oracle/ -v
```

Esperado: `PASS`. Estes testes usam mock de query — não precisam de Oracle.

- [ ] **Step 8: Commit**

```bash
git add apps/server_core/internal/modules/internal_read/adapters/mirror apps/server_core/internal/modules/internal_read/adapters/oracle/sync.go
git commit -m "feat(fiscal): grupo_icms atravessa Q1 do sankhya ate products_mirror"
```

---

## Task 3: Resolução da célula por lista branca — domínio puro

Esta é a regra que decide **qual linha da matriz vale**. Fica em domínio puro para ser testável sem Oracle e sem Postgres.

**Files:**
- Create: `apps/server_core/internal/modules/internal_read/domain/icms_matrix.go`
- Test: `apps/server_core/internal/modules/internal_read/domain/icms_matrix_test.go`

**O que foi medido, e que esta task implementa:** a matriz do ERP tem linhas com restrição. O domínio observado é `N` (sem restrição, código 0), `S` (curinga, código −1), `I` (grupo ICMS), `O` (TOP), `P` (produto), `H` (NCM), mais cinco letras não decodificadas (`K`, `G`, `L`, `X`, `T`). **A precedência formal entre linhas concorrentes NÃO foi estabelecida** — houve uma única observação, insuficiente para generalizar. Por isso: lista branca, e ambiguidade vira recusa.

- [ ] **Step 1: Escreva os testes que falham**

```go
package domain_test

import (
	"testing"

	"marketplace-central/apps/server_core/internal/modules/internal_read/domain"
)

func tributado() *int { z := 0; return &z }

func linha(tipo1 string, cod1 int, tipo2 string, cod2 int) domain.MatrixLine {
	return domain.MatrixLine{
		Restricao1Tipo: tipo1, Restricao1Codigo: cod1,
		Restricao2Tipo: tipo2, Restricao2Codigo: cod2,
		Sequencia: 1, CodTrib: tributado(),
		AliquotaPct: "7.0", RedBasePct: "0", AliqIntDestPct: "20.5",
	}
}

func TestResolveCellAceitaSemRestricao(t *testing.T) {
	got := domain.ResolveCell([]domain.MatrixLine{linha("N", 0, "I", 122)}, 122)
	if got.Ambiguo {
		t.Fatal("N/0 + I/<grupo> e a forma canonica; nao pode dar ambiguo")
	}
	if got.LinhasCandidatas != 1 {
		t.Fatalf("LinhasCandidatas = %d, queria 1", got.LinhasCandidatas)
	}
	if got.Line == nil || got.Line.AliqIntDestPct != "20.5" {
		t.Fatalf("linha resolvida errada: %+v", got.Line)
	}
}

func TestResolveCellDescartaGrupoDeOutroProduto(t *testing.T) {
	// I/284 nao serve para o produto do grupo 122.
	got := domain.ResolveCell([]domain.MatrixLine{linha("I", 284, "S", -1)}, 122)
	if got.Line != nil {
		t.Fatal("linha de outro grupo nao pode ser escolhida")
	}
	if got.LinhasCandidatas != 0 {
		t.Fatalf("LinhasCandidatas = %d, queria 0", got.LinhasCandidatas)
	}
}

func TestResolveCellDescartaTipoNaoDecodificado(t *testing.T) {
	// K, G, L, X, T nao foram decodificados. Lista BRANCA: o que nao esta na
	// lista e descartado, nunca "provavelmente serve".
	for _, tipo := range []string{"K", "G", "L", "X", "T", "O", "P", "H"} {
		got := domain.ResolveCell([]domain.MatrixLine{linha(tipo, 1000, "N", 0)}, 122)
		if got.Line != nil {
			t.Fatalf("tipo %q foi aceito; a lista branca e N, I e S", tipo)
		}
	}
}

func TestResolveCellRecusaEscolherEntreDuas(t *testing.T) {
	duas := []domain.MatrixLine{
		linha("N", 0, "I", 122),
		linha("S", -1, "I", 122),
	}
	got := domain.ResolveCell(duas, 122)
	if !got.Ambiguo {
		t.Fatal("duas linhas sobreviveram: tem que marcar ambiguo, nao escolher")
	}
	if got.Line != nil {
		t.Fatal("celula ambigua nao pode devolver linha; a precedencia nao foi medida")
	}
	if got.LinhasCandidatas != 2 {
		t.Fatalf("LinhasCandidatas = %d, queria 2", got.LinhasCandidatas)
	}
}

func TestResolveCellAceitaIntraUFComSlotsInvertidos(t *testing.T) {
	// Forma intra-UF medida: grupo no slot 1, sentinela "S" com codigo 1 (nao
	// -1) no slot 2. Sao 56 linhas do recorte real. Exigir codigo do sentinela
	// as rejeitaria todas.
	got := domain.ResolveCell([]domain.MatrixLine{linha("I", 122, "S", 1)}, 122)
	if got.Line == nil {
		t.Fatal("forma intra-UF (grupo no slot 1, sentinela S/1) foi rejeitada")
	}
}
```

- [ ] **Step 2: Rode e confirme que falha**

```bash
cd apps/server_core && GOCACHE=$(pwd)/.gocache go test ./internal/modules/internal_read/domain/ -run TestResolveCell -v
```

Esperado: `FAIL` — não compila, `ResolveCell` não existe.

- [ ] **Step 3: Implemente**

```go
package domain

// MatrixLine é uma linha crua da matriz de ICMS do ERP, já lida do Oracle e
// ainda não resolvida. Percentuais ficam como string: eles vêm de FLOAT e
// viram big.Rat mais tarde, sem passar por float binário.
//
// Sequencia faz parte da chave primária da TGFICM (sete colunas). Hoje vale 1
// em 4.601 de 4.601 linhas, mas a PK permite mais de uma — duas linhas que
// difiram só por Sequencia sobrevivem juntas à lista branca e a célula sai
// ambígua, que é o comportamento correto: a precedência entre elas não foi
// medida.
type MatrixLine struct {
	Restricao1Tipo   string // TIPRESTRICAO
	Restricao1Codigo int    // CODRESTRICAO
	Restricao2Tipo   string // TIPRESTRICAO2
	Restricao2Codigo int    // CODRESTRICAO2
	Sequencia        int    // SEQUENCIA
	// CodTrib é nullable na origem: 22 linhas com origem MG têm CODTRIB NULL.
	// Nil não é 0 — 0 significa "tributado" e seria uma afirmação.
	CodTrib *int
	// Zerar vem de ZERAR VARCHAR2(1). 'S' zera o imposto. Hoje 'S' em 0 linhas
	// MG, mas é NOT NULL na origem e barato de honrar.
	Zerar bool
	// AliquotaPct é ALIQUOTA — a alíquota da OPERAÇÃO origem→destino.
	// NÃO é ALIQUFDEST, que pertence ao bloco de cálculo de ST (17,0 a 22,0) e
	// daria um ICMS cerca de 2,5x maior.
	AliquotaPct string
	// RedBasePct é REDBASE, a redução da base da operação. Não-zero em 2.669 de
	// 4.209 linhas MG. Ignorá-la erra o ICMS na maioria das células.
	RedBasePct string
	// AliqIntDestPct é ALIQINTDEST — a que o motor do ERP lê e a que reconciliou
	// o DIFAL das notas reais. NULL em 82,6% das linhas MG.
	AliqIntDestPct string
	// AliqUFDestPct é ALIQUFDEST. Espelhada mas NÃO usada no cálculo: pertence ao
	// bloco de ST. Fica aqui porque tem 87,7% de cobertura e está mais perto da
	// alíquota legal, o que a torna o sinal de divergência da fatia P5.
	AliqUFDestPct string
	// PercRedBaseDestPct é PERCREDBASEDEST. Papel NÃO VERIFICADO — espelhada,
	// nunca multiplicada. Ver a regra de REDBASE no cálculo.
	PercRedBaseDestPct string
	PercFCPPct         string // PERCICMSFCP
}

// ResolvedCell é o veredito para uma célula (origem, destino, grupo).
type ResolvedCell struct {
	Line             *MatrixLine
	LinhasCandidatas int
	Ambiguo          bool
}

// ResolveCell aplica a LISTA BRANCA de restrição e devolve a linha vencedora.
//
// Aceita exatamente três tipos, que foram os medidos:
//   - "N" → sem restrição (sentinela; o código não significa nada)
//   - "S" → curinga (sentinela; o código não significa nada)
//   - "I" → vale para o grupo de ICMS informado no código
//
// O código dos sentinelas NÃO é constante e não pode ser testado: a forma
// interestadual é "N"/0 + "I"/<grupo>, e a forma intra-UF inverte os papéis dos
// slots para "I"/<grupo> + "S"/1. Exigir -1 no "S" rejeitaria as 56 linhas
// intra-MG do recorte real.
//
// Descarta todo o resto — "O" (TOP), "P" (produto), "H" (NCM) e as cinco letras
// não decodificadas ("K","G","L","X","T"). Lista BRANCA e não lista negra: tipo
// desconhecido vira descarte, nunca escolha silenciosa.
//
// Se mais de uma linha sobreviver, RECUSA escolher: a precedência formal entre
// linhas concorrentes não foi estabelecida por medição, e chutar precedência
// produziria um imposto errado com cara de certo. Ambiguo = true, Line = nil,
// e o leitor devolve desconhecido com pendência.
func ResolveCell(lines []MatrixLine, grupoICMS int) ResolvedCell {
	var sobreviventes []MatrixLine
	for _, line := range lines {
		if slotAceito(line.Restricao1Tipo, line.Restricao1Codigo, grupoICMS) &&
			slotAceito(line.Restricao2Tipo, line.Restricao2Codigo, grupoICMS) {
			sobreviventes = append(sobreviventes, line)
		}
	}
	out := ResolvedCell{LinhasCandidatas: len(sobreviventes)}
	switch len(sobreviventes) {
	case 0:
		return out
	case 1:
		vencedora := sobreviventes[0]
		out.Line = &vencedora
		return out
	default:
		out.Ambiguo = true
		return out
	}
}

func slotAceito(tipo string, codigo, grupoICMS int) bool {
	switch tipo {
	case "N", "S":
		// Sentinelas. O código varia com a forma da linha ("N"/0 interestadual,
		// "S"/1 intra-UF) e não carrega significado — testá-lo rejeita linhas
		// válidas. A guarda contra excesso de casamento é a ambiguidade.
		return true
	case "I":
		return codigo == grupoICMS
	default:
		return false
	}
}
```

- [ ] **Step 4: Rode e confirme que passa**

```bash
cd apps/server_core && GOCACHE=$(pwd)/.gocache go test ./internal/modules/internal_read/domain/ -run TestResolveCell -v
```

Esperado: `PASS` nos quatro.

- [ ] **Step 5: Controle negativo — prove que a recusa de ambiguidade é real**

Troque temporariamente o `default:` de `ResolveCell` por `out.Line = &sobreviventes[0]; return out` e rode de novo.

Esperado: `TestResolveCellRecusaEscolherEntreDuas` **falha**. Se ele passar, o teste não está provando nada — conserte o teste, não o código. Desfaça a alteração depois.

- [ ] **Step 6: Commit**

```bash
git add apps/server_core/internal/modules/internal_read/domain/icms_matrix.go apps/server_core/internal/modules/internal_read/domain/icms_matrix_test.go
git commit -m "feat(fiscal): resolucao da celula da matriz por lista branca, ambiguidade recusa escolher"
```

---

## Task 4: Extração da matriz do Oracle

**Files:**
- Create: `apps/server_core/internal/modules/internal_read/adapters/oracle/icms_matrix.go`
- Test: `apps/server_core/internal/modules/internal_read/adapters/oracle/icms_matrix_test.go`

**⚠ Duas armadilhas medidas, leia antes de escrever:**

1. **A alíquota da operação é `ALIQUOTA`.** `ALIQUFDEST` **não é** — ela pertence ao bloco de cálculo de ST (anda com `MVASTUFDEST`, `BASESTUFDEST`, `CODTABSTUFDEST`) e seu domínio é 17,0 a 22,0, nunca 7 nem 12. Multiplicar a base por ela dá um ICMS ~2,5× maior. O código existente em `icms_ceiling.go:56` usa `ALIQUFDEST` legitimamente, mas como **teto**, não como alíquota de operação — não copie de lá.
2. **`UFORIG`/`UFDEST` não são código IBGE nem ordem alfabética.** MG = **13** (IBGE seria 31). A de-para é `TSIUFS.CODUF → TSIUFS.UF`, e o join **exige `CODPAIS = 55`**: `CODUF = 0` é a sentinela `'SF'` e os códigos 29 e 30 são **ambos `'EX'`**. Sem o filtro o join deixa de ser single-valued e duplica linha.

- [ ] **Step 1: Escreva o teste que falha**

Use o mock de query já usado em `batch_reader_test.go` — o teste alimenta linhas cruas e afirma o agrupamento, sem tocar em Oracle.

```go
func TestReadICMSMatrixAgrupaPorDestinoEGrupo(t *testing.T) {
	db := newFakeOracle(t, [][]any{
		// UFDEST_SIGLA, TIPRESTRICAO, CODRESTRICAO, TIPRESTRICAO2, CODRESTRICAO2,
		// SEQUENCIA, CODTRIB, ZERAR, ALIQUOTA, REDBASE, ALIQINTDEST, ALIQUFDEST,
		// PERCREDBASEDEST, PERCICMSFCP
		{"BA", "N", 0, "I", 122, 1, 0, "N", 7.0, 0.0, 20.5, 20.5, nil, nil},
		{"BA", "N", 0, "I", 284, 1, 60, "N", 7.0, 0.0, nil, 20.5, nil, nil},
		{"SP", "N", 0, "I", 122, 1, 0, "N", 12.0, 0.0, nil, 18.0, nil, nil},
	})
	r := oracle.NewSankhyaAdapter(db)

	got, err := r.ReadICMSMatrix(context.Background(), 13)
	if err != nil {
		t.Fatalf("ReadICMSMatrix: %v", err)
	}
	if len(got) != 3 {
		t.Fatalf("chaves = %d, queria 3", len(got))
	}
	ba122 := got[oracle.MatrixKey{UFDestino: "BA", GrupoICMS: 122}]
	if len(ba122) != 1 {
		t.Fatalf("BA/122 tem %d linhas, queria 1", len(ba122))
	}
	if ba122[0].AliquotaPct != "7" {
		t.Fatalf("AliquotaPct = %q, queria 7 — ALIQUOTA, nao ALIQUFDEST", ba122[0].AliquotaPct)
	}
	// ALIQUFDEST e espelhada mas NUNCA vira aliquota de operacao.
	if ba122[0].AliquotaPct == ba122[0].AliqUFDestPct && ba122[0].AliqUFDestPct != "20.5" {
		t.Fatal("confusao entre ALIQUOTA e ALIQUFDEST")
	}
}

func TestReadICMSMatrixCodTribNuloNaoViraZero(t *testing.T) {
	db := newFakeOracle(t, [][]any{
		{"BA", "N", 0, "I", 122, 1, nil, "N", 7.0, 0.0, 20.5, 20.5, nil, nil},
	})
	r := oracle.NewSankhyaAdapter(db)

	got, err := r.ReadICMSMatrix(context.Background(), 13)
	if err != nil {
		t.Fatalf("ReadICMSMatrix: %v", err)
	}
	linha := got[oracle.MatrixKey{UFDestino: "BA", GrupoICMS: 122}][0]
	if linha.CodTrib != nil {
		t.Fatalf("CodTrib = %v; NULL na origem (22 linhas MG) nao pode virar 0, "+
			"porque 0 significa 'tributado' e seria uma afirmacao", *linha.CodTrib)
	}
}

func TestReadICMSMatrixIgnoraLinhaSemSlotDeGrupo(t *testing.T) {
	// Slot de grupo ausente nos dois lados: a linha nao entra em nenhuma chave.
	db := newFakeOracle(t, [][]any{
		{"BA", "N", 0, "S", -1, 1, 0, "N", 7.0, 0.0, 20.5, 20.5, nil, nil},
	})
	r := oracle.NewSankhyaAdapter(db)

	got, err := r.ReadICMSMatrix(context.Background(), 13)
	if err != nil {
		t.Fatalf("ReadICMSMatrix: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("chaves = %d, queria 0 — sem slot 'I' nao ha grupo a que ancorar", len(got))
	}
}
```

- [ ] **Step 2: Rode e confirme que falha**

```bash
cd apps/server_core && GOCACHE=$(pwd)/.gocache go test ./internal/modules/internal_read/adapters/oracle/ -run TestReadICMSMatrix -v
```

Esperado: `FAIL` — `ReadICMSMatrix` não existe.

- [ ] **Step 3: Implemente**

```go
// MatrixKey é a célula (destino, grupo de ICMS) da matriz.
type MatrixKey struct {
	UFDestino string
	GrupoICMS int
}

// icmsMatrixSQL lê a matriz de ICMS inteira com uma origem.
//
// São 4.209 linhas com origem MG. Tabela minúscula: copiamos inteira todo dia,
// sem delta. Detecção de mudança é MAX(TGFHICM.DHALTER).
//
// TSIUFS traduz o código numérico de UF para sigla. CODPAIS = 55 é
// OBRIGATÓRIO: CODUF = 0 é a sentinela 'SF' e os códigos 29 e 30 são ambos
// 'EX'. Sem o filtro o join duplica linha.
//
// ALIQUOTA é a alíquota da OPERAÇÃO. ALIQUFDEST vem junto porque é o sinal de
// divergência da fatia P5, e nunca entra no cálculo.
//
// Somente leitura. Nada é escrito em METALPRD.
const icmsMatrixSQL = `
SELECT ud.UF AS UFDEST_SIGLA,
       i.TIPRESTRICAO, i.CODRESTRICAO, i.TIPRESTRICAO2, i.CODRESTRICAO2, i.SEQUENCIA,
       i.CODTRIB, i.ZERAR,
       i.ALIQUOTA, i.REDBASE,
       i.ALIQINTDEST, i.ALIQUFDEST, i.PERCREDBASEDEST, i.PERCICMSFCP
FROM METALPRD.TGFICM i
JOIN METALPRD.TSIUFS ud ON ud.CODUF = i.UFDEST AND ud.CODPAIS = 55
WHERE i.UFORIG = :uforig
ORDER BY ud.UF, i.TIPRESTRICAO2, i.CODRESTRICAO2, i.TIPRESTRICAO, i.CODRESTRICAO, i.SEQUENCIA`
```

`ReadICMSMatrix(ctx, originUF int64)` escaneia, monta a `domain.MatrixLine` e agrupa. O grupo sai do slot que tiver tipo `"I"` — **slot 2 na forma normal (4.054 de 4.209 linhas), slot 1 na forma intra-UF**. Linha sem nenhum slot `"I"` não entra em chave nenhuma.

Percentuais chegam como `FLOAT` do Oracle: escaneie para `sql.NullFloat64` e converta para string com `strconv.FormatFloat(v, 'f', -1, 64)`. `-1` de precisão preserva o valor sem inventar casas.

> **Nome do schema:** o especialista mediu em `SANKHYA.TGFICM`; o código existente (`icms_ceiling.go:56`) usa `METALPRD.TGFICM` e funciona. São sinônimos. Use `METALPRD`, igual ao resto do módulo. Se der `ORA-00942`, mande `ESCALATION` — não troque de schema por conta própria.

- [ ] **Step 4: Rode e confirme que passa**

```bash
cd apps/server_core && GOCACHE=$(pwd)/.gocache go test ./internal/modules/internal_read/adapters/oracle/ -run TestReadICMSMatrix -v
```

Esperado: `PASS` nos três.

- [ ] **Step 5: Controle negativo — prove que o filtro de país morde**

Remova `AND ud.CODPAIS = 55` e conte as linhas devolvidas por um destino qualquer.

Esperado: contagem **maior**, porque `CODUF` 29 e 30 são ambos `'EX'`. Se não mudar nada, seu fake de Oracle não tem as linhas duplicadas — acrescente-as, senão o teste não cobre o motivo do filtro existir. Desfaça depois.

- [ ] **Step 6: Commit**

```bash
git add apps/server_core/internal/modules/internal_read/adapters/oracle/
git commit -m "feat(fiscal): leitura da matriz TGFICM com origem MG e de-para de UF por TSIUFS"
```

**Dívida a registrar nesta task:** a origem é fixa em `13` (MG, empresa 1). Mesma classe da **D-17** (`CODEMP=1` fixo na leitura de custo). Deixe o comentário nomeando **D-29** onde o `13` é constante.

---

## Task 5: Escrita versionada da matriz (SCD-2)

**Files:**
- Create: `apps/server_core/internal/modules/internal_read/adapters/mirror/icms_matrix_writer.go`
- Test: `apps/server_core/internal/modules/internal_read/adapters/mirror/icms_matrix_writer_integration_test.go`

- [ ] **Step 1: Escreva os testes que falham**

```go
//go:build integration

package mirror_test

func TestApplyMatrixAbreVersaoNova(t *testing.T) {
	ctx := context.Background()
	pool := openTestPool(t)
	defer pool.Close()
	tenantID := newTestTenant(t, pool)
	w := mirror.NewICMSMatrixWriter(pool)

	celula := func(aliqIntDest string) mirror.MatrixCell {
		return mirror.MatrixCell{
			UFOrigem: "MG", UFDestino: "BA", GrupoICMS: 122, CodTrib: ptrInt(0),
			AliquotaPct: ptr("7.000"), RedBasePct: ptr("0.000"), AliqIntDestPct: ptr(aliqIntDest),
			LinhasCandidatas: 1, Ambiguo: false,
		}
	}

	if err := w.ApplyMatrix(ctx, tenantID, []mirror.MatrixCell{celula("17.000")}); err != nil {
		t.Fatalf("primeira aplicacao: %v", err)
	}
	// O ERP corrigiu a Bahia de 17,0 para 20,5.
	if err := w.ApplyMatrix(ctx, tenantID, []mirror.MatrixCell{celula("20.500")}); err != nil {
		t.Fatalf("segunda aplicacao: %v", err)
	}

	var versoes, abertas int
	pool.QueryRow(ctx, `SELECT count(*), count(*) FILTER (WHERE vigente_ate IS NULL)
		FROM icms_matrix_mirror WHERE tenant_id = $1`, tenantID).Scan(&versoes, &abertas)
	if versoes != 2 {
		t.Fatalf("versoes = %d, queria 2 — a correcao do ERP tem que virar versao nova, nao UPDATE", versoes)
	}
	if abertas != 1 {
		t.Fatalf("versoes abertas = %d, queria 1", abertas)
	}

	// A versao antiga tem que continuar legivel: e ela que prova a defasagem.
	var antiga string
	pool.QueryRow(ctx, `SELECT aliqintdest::text FROM icms_matrix_mirror
		WHERE tenant_id = $1 AND vigente_ate IS NOT NULL`, tenantID).Scan(&antiga)
	if antiga != "17.000" {
		t.Fatalf("versao fechada = %q, queria 17.000 — apagar o passado destroi a evidencia", antiga)
	}
}

func TestApplyMatrixNaoVersionaSemMudanca(t *testing.T) {
	ctx := context.Background()
	pool := openTestPool(t)
	defer pool.Close()
	tenantID := newTestTenant(t, pool)
	w := mirror.NewICMSMatrixWriter(pool)

	c := mirror.MatrixCell{
		UFOrigem: "MG", UFDestino: "SP", GrupoICMS: 122, CodTrib: ptrInt(0),
		AliquotaPct: ptr("7.000"), RedBasePct: ptr("0.000"), AliqIntDestPct: ptr("18.000"),
		LinhasCandidatas: 1, Ambiguo: false,
	}
	for i := 0; i < 3; i++ {
		if err := w.ApplyMatrix(ctx, tenantID, []mirror.MatrixCell{c}); err != nil {
			t.Fatalf("aplicacao %d: %v", i, err)
		}
	}

	var versoes int
	pool.QueryRow(ctx, `SELECT count(*) FROM icms_matrix_mirror WHERE tenant_id = $1`, tenantID).Scan(&versoes)
	if versoes != 1 {
		t.Fatalf("versoes = %d, queria 1 — sync diario sem mudanca nao pode inflar historico", versoes)
	}
}
```

- [ ] **Step 2: Rode e confirme que falha**

```bash
cd apps/server_core && GOCACHE=$(pwd)/.gocache go test -tags integration ./internal/modules/internal_read/adapters/mirror/ -run TestApplyMatrix -v
```

Esperado: `FAIL` — `NewICMSMatrixWriter` não existe.

- [ ] **Step 3: Implemente**

`ApplyMatrix` roda numa transação e, para cada célula:

1. Lê a versão aberta (`vigente_ate IS NULL`) daquela chave.
2. Se não existe → `INSERT` com `vigente_desde = now()`.
3. Se existe e **todos** os campos de valor batem (`codtrib`, `zerar`, `aliquota`, `redbase`, `aliqintdest`, `aliq_uf_dest`, `perc_red_base_dest`, `perc_fcp`, `ambiguo`, `linhas_candidatas`) → só atualiza `synced_at`. **Nenhuma versão nova.**
4. Se existe e algum mudou → `UPDATE ... SET vigente_ate = now()` na antiga, `INSERT` da nova com `vigente_desde = now()`.

Comparação de `NUMERIC` nullable: compare em **texto** (`aliqintdest::text`) ou use `IS NOT DISTINCT FROM`, que trata `NULL` como igual a `NULL`. Comparar com `=` faz toda célula muda parecer mudada em todo sync e infla o histórico sem fim.

Células que sumiram do ERP: fecha a versão aberta (`vigente_ate = now()`) e não abre nova. Nunca `DELETE` — a evidência do passado é o produto desta tabela.

- [ ] **Step 4: Rode e confirme que passa**

```bash
cd apps/server_core && GOCACHE=$(pwd)/.gocache go test -tags integration ./internal/modules/internal_read/adapters/mirror/ -run TestApplyMatrix -v
```

Esperado: `PASS` nos dois.

- [ ] **Step 5: Controle negativo — prove que o teste de não-versionar morde**

Troque a comparação por `=` no lugar de `IS NOT DISTINCT FROM` e rode `TestApplyMatrixNaoVersionaSemMudanca` com uma célula de `aliqintdest` **nulo**.

Esperado: **falha** com `versoes = 3`. Se passar, o teste não cobre a célula muda, que é justamente o caso de PR, RS e SC. Desfaça depois.

- [ ] **Step 6: Commit**

```bash
git add apps/server_core/internal/modules/internal_read/adapters/mirror/
git commit -m "feat(fiscal): escrita versionada da matriz de ICMS, correcao do ERP vira versao nova"
```

---

## Task 6: Domínio `fiscal` — arredondamento em centavos

O cálculo produz dinheiro. Dinheiro em float binário produz `5,1959999`. Esta task existe para isso não acontecer.

**Files:**
- Create: `apps/server_core/internal/modules/fiscal/domain/money.go`
- Test: `apps/server_core/internal/modules/fiscal/domain/money_test.go`

- [ ] **Step 1: Escreva o teste que falha**

```go
package domain

import "testing"

func TestPctOfArredondaMeiaParaCima(t *testing.T) {
	casos := []struct {
		valor float64
		pct   string
		want  float64
		nota  string
	}{
		{299.90, "7.0", 20.99, "20,993 desce"},
		{299.90, "20.5", 61.48, "61,4795 sobe"},
		{129.90, "4.0", 5.20, "5,196 sobe; em float binario daria 5,1959999"},
		{0.10, "50.0", 0.05, "exatamente meio centavo sobe"},
	}
	for _, c := range casos {
		got, err := pctOf(c.valor, c.pct)
		if err != nil {
			t.Fatalf("%s: %v", c.nota, err)
		}
		if got != c.want {
			t.Fatalf("%s: pctOf(%v, %q) = %v, queria %v", c.nota, c.valor, c.pct, got, c.want)
		}
	}
}
```

- [ ] **Step 2: Rode e confirme que falha**

```bash
cd apps/server_core && GOCACHE=$(pwd)/.gocache go test ./internal/modules/fiscal/domain/ -run TestPctOf -v
```

Esperado: `FAIL` — pacote não existe.

- [ ] **Step 3: Implemente com `big.Rat`**

```go
// Package domain contém as regras fiscais. É puro: não abre conexão, não lê
// configuração, não sabe que existe tela.
package domain

import (
	"math/big"
	"strconv"
)

var cem = big.NewRat(100, 1)

// pctOf devolve pct% de valor, arredondado a centavos, meia-para-cima.
//
// Tudo passa por big.Rat porque o resultado é dinheiro: 4% de R$ 129,90 é
// 5,196 → 5,20, nunca o 5,1959999 que o float binário produz.
//
// D-30: o módulo pricing tem uma rotina equivalente (pricingdomain.ParseRat +
// FormatRatHalfUp). Não importamos de lá porque fiscal é mais baixo que pricing
// e importar inverteria a dependência que esta fatia existe para desfazer. Duas
// implementações da mesma regra podem divergir — TestPctOfCasaComPricing existe
// para pegar isso.
func pctOf(valor float64, pct string) (float64, error) {
	amount := new(big.Rat).SetFloat64(valor)
	if amount == nil { // NaN ou ±Inf
		return 0, errValorNaoNumerico
	}
	parsed, ok := new(big.Rat).SetString(pct)
	if !ok {
		return 0, errPercentualInvalido
	}
	v := new(big.Rat).Quo(parsed, cem)
	v.Mul(v, amount)
	return strconv.ParseFloat(v.FloatString(2), 64)
}
```

`FloatString(2)` do `big.Rat` já arredonda meia-para-cima em magnitude. Declare `errValorNaoNumerico` e `errPercentualInvalido` como `errors.New` no mesmo arquivo.

- [ ] **Step 4: Rode e confirme que passa**

```bash
cd apps/server_core && GOCACHE=$(pwd)/.gocache go test ./internal/modules/fiscal/domain/ -run TestPctOf -v
```

Esperado: `PASS`.

- [ ] **Step 5: Escreva o teste de não-divergência com `pricing`**

```go
// TestPctOfCasaComPricing existe por causa da divida D-30: duas rotinas de
// arredondamento na mesma base podem divergir em silencio. Este teste e o
// alarme. Se ele falhar, NAO ajuste o numero esperado — descubra qual das duas
// esta errada.
func TestPctOfCasaComPricing(t *testing.T) {
	casos := []struct {
		valor float64
		pct   string
	}{{299.90, "7.0"}, {299.90, "20.5"}, {129.90, "4.0"}, {1.00, "33.333"}, {0.10, "50.0"}}
	for _, c := range casos {
		meu, err := pctOf(c.valor, c.pct)
		if err != nil {
			t.Fatalf("pctOf: %v", err)
		}
		parsed, err := pricingdomain.ParseRat(c.pct)
		if err != nil {
			t.Fatalf("ParseRat: %v", err)
		}
		v := new(big.Rat).Quo(parsed, cem)
		v.Mul(v, new(big.Rat).SetFloat64(c.valor))
		deles, err := strconv.ParseFloat(pricingdomain.FormatRatHalfUp(v, 2), 64)
		if err != nil {
			t.Fatalf("FormatRatHalfUp: %v", err)
		}
		if meu != deles {
			t.Fatalf("divergencia em %v @ %s: fiscal=%v pricing=%v", c.valor, c.pct, meu, deles)
		}
	}
}
```

Este teste fica em `money_test.go`, que é `package domain` (teste interno). Importar `pricing/domain` **num teste** não cria dependência de produção — o `go list -deps` do pacote de produção continua limpo. A Task 16 verifica isso.

- [ ] **Step 6: Rode os dois e commit**

```bash
cd apps/server_core && GOCACHE=$(pwd)/.gocache go test ./internal/modules/fiscal/domain/ -v
git add apps/server_core/internal/modules/fiscal/domain/
git commit -m "feat(fiscal): arredondamento em centavos com big.Rat e alarme de divergencia com pricing"
```

---

## Task 7: Domínio `fiscal` — tipos

**Files:**
- Create: `apps/server_core/internal/modules/fiscal/domain/rules.go`

Sem teste próprio: são tipos, e os testes da Task 8 os exercitam. Escrever teste de struct vazio é teste vazio.

- [ ] **Step 1: Escreva os tipos**

```go
package domain

import "time"

// Componente é um valor fiscal que pode não existir. Carrega o valor OU o
// motivo, nunca os dois, nunca nenhum dos dois.
//
// ADR-17: desconhecido nunca vira zero. Um DIFAL que não conseguimos calcular
// não é um DIFAL de R$ 0,00 — é uma lacuna do cadastro do ERP, e o Motivo diz
// exatamente qual, em português, para ser levado ao Sankhya.
type Componente struct {
	Valor  *float64
	Motivo string
}

// Conhecido monta um componente com valor.
func Conhecido(v float64) Componente { return Componente{Valor: &v} }

// Desconhecido monta um componente com a pendência. motivo nunca pode ser "".
func Desconhecido(motivo string) Componente { return Componente{Motivo: motivo} }

// ICMSRule é a célula resolvida da matriz do ERP, do jeito que o cálculo
// precisa dela. Percentuais como string: vêm de NUMERIC e viram big.Rat, sem
// float binário no meio.
type ICMSRule struct {
	// CodTrib: 0 = tributado, 60 = substituição tributária, 10 e 40 também
	// ocorrem. Nil quando a origem tem NULL (22 linhas MG) — nil não é 0.
	CodTrib *int
	// Zerar vem de TGFICM.ZERAR = 'S': a célula zera o imposto.
	Zerar bool
	// AliquotaPct é a alíquota da OPERAÇÃO origem→destino (7% ou 12% saindo de
	// MG). É TGFICM.ALIQUOTA. Nunca ALIQUFDEST, que é do bloco de ST.
	AliquotaPct *string
	// RedBasePct é a redução da base da OPERAÇÃO (TGFICM.REDBASE). Não-zero em
	// 63% das linhas MG. Nil é tratado como zero: a célula existe e não
	// declarou redução.
	RedBasePct *string
	// AliqIntDestPct é a alíquota INTERNA do estado de destino, a que fecha o
	// DIFAL, e é a que o motor do ERP lê. Nil = célula muda: a matriz diz que o
	// DIFAL incide mas não diz quanto. NULL em 82,6% das linhas MG — AC, AL, AM,
	// AP, PA, PB, PI, RO, RR, SE, RS e SC são 100% nulos.
	AliqIntDestPct *string
	// PercRedBaseDestPct é a redução de base do lado do destino
	// (TGFICM.PERCREDBASEDEST). Papel NÃO VERIFICADO na origem. Este campo
	// existe para o cálculo saber que ele é desconhecido, nunca para multiplicar.
	PercRedBaseDestPct *string
	// PercFCPPct nil com a célula presente significa FCP zero — a matriz
	// respondeu. É diferente de célula ausente, em que nada é conhecido.
	PercFCPPct *string
	// Ambiguo: mais de uma linha sobreviveu à lista branca e o sync se recusou
	// a escolher.
	Ambiguo bool
	// VigenteDesde é quando esta versão da célula passou a valer no espelho.
	VigenteDesde time.Time
}

// Origem diz de onde o número veio. Nesta fatia é sempre OrigemMatriz; o P6
// acrescenta a leitura da nota emitida, reusando este mesmo bloco.
type Origem string

const (
	OrigemMatriz Origem = "matriz_erp"
	OrigemNota   Origem = "nota_erp"
)

// Resultado é o imposto de UMA linha de pedido.
type Resultado struct {
	ICMS         Componente
	DIFAL        Componente
	FCP          Componente
	PisCofins    Componente
	Total        Componente
	CargaICMSPct Componente
	Origem       Origem
	MatrizDesde  *time.Time
}
```

- [ ] **Step 2: Compile**

```bash
cd apps/server_core && GOCACHE=$(pwd)/.gocache go build ./internal/modules/fiscal/...
```

Esperado: sem saída.

- [ ] **Step 3: Commit**

```bash
git add apps/server_core/internal/modules/fiscal/domain/rules.go
git commit -m "feat(fiscal): tipos do dominio fiscal com componente conhecido-ou-motivo"
```

---

## Task 8: `TaxesForValue` — o coração

**Files:**
- Create: `apps/server_core/internal/modules/fiscal/domain/calc.go`
- Test: `apps/server_core/internal/modules/fiscal/domain/calc_test.go`

**A autoridade aqui é a matriz vigente, não a nota emitida.** A nota 895507 gravou DIFAL 29,99 porque foi emitida com alíquota defasada de 17,0%; a nota-irmã do dia anterior, mesmo grupo e destino, usou 20,5% e reconcilia ao centavo. Se algum dia este teste "quebrar" e alguém for tentado a ajustá-lo para 29,99, isso petrifica a nota defasada como verdade — que é exatamente a doença que esta fatia existe para curar.

- [ ] **Step 1: Escreva os testes que falham**

```go
package domain

import (
	"strings"
	"testing"
	"time"
)

func s(v string) *string { return &v }

func i(v int) *int { return &v }

func regraBahia() ICMSRule {
	return ICMSRule{
		CodTrib:        i(0),
		AliquotaPct:    s("7.0"),
		RedBasePct:     s("0"),
		AliqIntDestPct: s("20.5"),
		VigenteDesde:   time.Date(2026, 7, 20, 0, 0, 0, 0, time.UTC),
	}
}

func valor(t *testing.T, c Componente, nome string) float64 {
	t.Helper()
	if c.Valor == nil {
		t.Fatalf("%s veio desconhecido: %q", nome, c.Motivo)
	}
	return *c.Valor
}

// A AUTORIDADE E A MATRIZ VIGENTE, NAO A NOTA EMITIDA. Ver cabecalho da task.
func TestTaxesForValueBahiaSegueMatrizVigente(t *testing.T) {
	got := TaxesForValue(299.90, regraBahia())

	if v := valor(t, got.ICMS, "ICMS"); v != 20.99 {
		t.Fatalf("ICMS = %v, queria 20.99 (299,90 x 7%%)", v)
	}
	if v := valor(t, got.DIFAL, "DIFAL"); v != 40.49 {
		t.Fatalf("DIFAL = %v, queria 40.49 (299,90 x 20,5%% - 20,99)", v)
	}
	if v := valor(t, got.FCP, "FCP"); v != 0 {
		t.Fatalf("FCP = %v, queria 0 — a matriz respondeu, e a resposta e zero", v)
	}
	// Base = 299,90 x (1 - 0,205) = 238,4205 ; x 9,25% = 22,0538 -> 22,05.
	// A base NAO e o valor de venda: o ICMS e repasse ao estado, nao receita
	// (STF, Tema 69). Medido no proprio ERP em 1.434/1.434 linhas.
	if v := valor(t, got.PisCofins, "PIS/COFINS"); v != 22.05 {
		t.Fatalf("PIS/COFINS = %v, queria 22.05", v)
	}
	if v := valor(t, got.Total, "Total"); v != 83.53 {
		t.Fatalf("Total = %v, queria 83.53", v)
	}

	// Contra-controle: 29,99 e o DIFAL que a nota 895507 gravou, com
	// TGFDIN.ALIQINTDEST = 17,0. A celula da TGFICM (grupo 122, destino BA)
	// vale 20,5 nas duas colunas, e ja valia no dia anterior — a nota-irma
	// 895436 usou 20,5 e reconcilia ao centavo. O 17,0 NAO sai da matriz;
	// vem de uma fonte ainda nao localizada (divida D-35). Espelhamos o
	// cadastro, e o cadastro diz 20,5.
	if got.DIFAL.Valor != nil && *got.DIFAL.Valor == 29.99 {
		t.Fatal("29,99 e o que a nota 895507 gravou — a celula da matriz diz 20,5%")
	}
	if got.Origem != OrigemMatriz {
		t.Fatalf("Origem = %q, queria %q", got.Origem, OrigemMatriz)
	}
}

func TestTaxesForValueCelulaMudaSelaDifalEPisCofins(t *testing.T) {
	// 20 celulas reais assim no nosso recorte, em PR, RS e SC. Nao e hipotese.
	regra := regraBahia()
	regra.AliqIntDestPct = nil

	got := TaxesForValue(299.90, regra)

	if v := valor(t, got.ICMS, "ICMS"); v != 20.99 {
		t.Fatalf("ICMS = %v — falta de aliquota interna nao apaga o ICMS proprio", v)
	}
	if got.DIFAL.Valor != nil {
		t.Fatalf("DIFAL = %v, queria desconhecido", *got.DIFAL.Valor)
	}
	if got.DIFAL.Motivo == "" {
		t.Fatal("DIFAL desconhecido sem motivo — a pendencia tem que dizer o que falta no Sankhya")
	}
	if !strings.Contains(got.DIFAL.Motivo, "ALIQINTDEST") {
		t.Fatalf("motivo = %q; tem que nomear a coluna que falta no ERP", got.DIFAL.Motivo)
	}
	// Cascata: sem a carga de ICMS nao ha base de PIS/COFINS, e sem ela nao ha total.
	if got.PisCofins.Valor != nil {
		t.Fatalf("PIS/COFINS = %v — a base depende da carga de ICMS, que e desconhecida", *got.PisCofins.Valor)
	}
	if got.Total.Valor != nil {
		t.Fatalf("Total = %v — total parcial e mentira com cara de numero", *got.Total.Valor)
	}
}

func TestTaxesForValueReduzBaseDoICMS(t *testing.T) {
	// REDBASE nao-zero em 2.669 de 4.209 linhas MG. Ignorar erra o ICMS na
	// maioria das celulas.
	regra := regraBahia()
	regra.RedBasePct = s("33.33")

	got := TaxesForValue(299.90, regra)

	// 299,90 x (1 - 0,3333) x 7% = 299,90 x 0,6667 x 0,07 = 13,9958 -> 14,00
	if v := valor(t, got.ICMS, "ICMS"); v != 14.00 {
		t.Fatalf("ICMS = %v, queria 14.00 — base reduzida em 33,33%%", v)
	}
	// REDBASE != 0 junto de CODTRIB = 0 tem ZERO observacoes no recorte real: a
	// unica ocorrencia de REDBASE != 0 e REDBASE = 100 com CODTRIB = 60, que e
	// marcacao de ST. Sem observacao, o DIFAL nao se afirma.
	if got.DIFAL.Valor != nil {
		t.Fatalf("DIFAL = %v; REDBASE != 0 sem ST nunca foi observado", *got.DIFAL.Valor)
	}
	if got.DIFAL.Motivo == "" {
		t.Fatal("faltou motivo")
	}
}

func TestTaxesForValueRedBase100ComSTEhST(t *testing.T) {
	// REDBASE = 100 nao e "base reduzida a 100%": e como o ERP marca ST, e vem
	// sempre com CODTRIB = 60. Tratar como reducao zeraria a base de uma venda.
	regra := regraBahia()
	regra.RedBasePct = s("100")
	regra.CodTrib = i(60)

	got := TaxesForValue(299.90, regra)

	for nome, c := range map[string]Componente{"ICMS": got.ICMS, "DIFAL": got.DIFAL, "FCP": got.FCP} {
		if v := valor(t, c, nome); v != 0 {
			t.Fatalf("%s = %v, queria 0", nome, v)
		}
	}
	// Base cheia: 299,90 x 9,25% = 27,74. Se REDBASE tivesse zerado a base, o
	// PIS/COFINS sairia 0 e o pedido pareceria isento.
	if v := valor(t, got.PisCofins, "PIS/COFINS"); v != 27.74 {
		t.Fatalf("PIS/COFINS = %v, queria 27.74 — REDBASE nao pode zerar a base de PIS/COFINS", v)
	}
}

func TestTaxesForValueCodTribNuloNaoCalculaNada(t *testing.T) {
	// 22 linhas MG com CODTRIB NULL. Tratar como 0 afirmaria "tributado".
	regra := regraBahia()
	regra.CodTrib = nil

	got := TaxesForValue(299.90, regra)

	if got.ICMS.Valor != nil {
		t.Fatalf("ICMS = %v; sem CODTRIB nao se sabe nem se incide", *got.ICMS.Valor)
	}
	if got.ICMS.Motivo == "" {
		t.Fatal("faltou motivo")
	}
}

func TestTaxesForValueCodTribNaoMedidoNaoChuta(t *testing.T) {
	// 10 (41 linhas) e 40 (2 linhas) ocorrem em MG e o comportamento deles NAO
	// foi medido. Assumir "igual a 0" seria inventar.
	for _, codigo := range []int{10, 40} {
		regra := regraBahia()
		regra.CodTrib = i(codigo)

		got := TaxesForValue(299.90, regra)

		if got.ICMS.Valor != nil {
			t.Fatalf("CODTRIB %d calculou ICMS = %v; comportamento nao medido", codigo, *got.ICMS.Valor)
		}
	}
}

func TestTaxesForValueZerarAnulaFamiliaICMS(t *testing.T) {
	regra := regraBahia()
	regra.Zerar = true

	got := TaxesForValue(299.90, regra)

	for nome, c := range map[string]Componente{"ICMS": got.ICMS, "DIFAL": got.DIFAL, "FCP": got.FCP} {
		if v := valor(t, c, nome); v != 0 {
			t.Fatalf("%s = %v com ZERAR='S', queria 0", nome, v)
		}
	}
	// Carga zero, base cheia: 299,90 x 9,25% = 27,74.
	if v := valor(t, got.PisCofins, "PIS/COFINS"); v != 27.74 {
		t.Fatalf("PIS/COFINS = %v, queria 27.74", v)
	}
}

func TestTaxesForValueSTZeraFamiliaICMS(t *testing.T) {
	regra := regraBahia()
	regra.CodTrib = i(60)

	got := TaxesForValue(299.90, regra)

	for nome, c := range map[string]Componente{"ICMS": got.ICMS, "DIFAL": got.DIFAL, "FCP": got.FCP} {
		if v := valor(t, c, nome); v != 0 {
			t.Fatalf("%s = %v, queria 0 — no ST o imposto foi pago na entrada e ja esta no CUSSEMICM", nome, v)
		}
	}
	// Carga zero, entao a base e o valor cheio: 299,90 x 9,25% = 27,74.
	if v := valor(t, got.PisCofins, "PIS/COFINS"); v != 27.74 {
		t.Fatalf("PIS/COFINS = %v, queria 27.74", v)
	}
}

func TestTaxesForValueDifalNegativoViraPendencia(t *testing.T) {
	// Aliquota interna do destino MENOR que a da operacao e impossivel na
	// pratica: e defeito de cadastro. Nunca devolver negativo, nunca clampar
	// em zero — clampar transforma defeito do ERP em "nada a pagar".
	regra := regraBahia()
	regra.AliqIntDestPct = s("4.0")

	got := TaxesForValue(299.90, regra)

	if got.DIFAL.Valor != nil {
		t.Fatalf("DIFAL = %v, queria pendencia", *got.DIFAL.Valor)
	}
	if got.DIFAL.Motivo == "" {
		t.Fatal("sem motivo")
	}
}

func TestTaxesForValueCelulaAmbiguaNaoCalculaNada(t *testing.T) {
	regra := regraBahia()
	regra.Ambiguo = true

	got := TaxesForValue(299.90, regra)

	for nome, c := range map[string]Componente{
		"ICMS": got.ICMS, "DIFAL": got.DIFAL, "FCP": got.FCP,
		"PIS/COFINS": got.PisCofins, "Total": got.Total,
	} {
		if c.Valor != nil {
			t.Fatalf("%s = %v — celula ambigua nao autoriza nenhum numero", nome, *c.Valor)
		}
	}
}
```

- [ ] **Step 2: Rode e confirme que falha**

```bash
cd apps/server_core && GOCACHE=$(pwd)/.gocache go test ./internal/modules/fiscal/domain/ -run TestTaxesForValue -v
```

Esperado: `FAIL` — `TaxesForValue` não existe.

- [ ] **Step 3: Implemente**

```go
package domain

import (
	"errors"
	"fmt"
	"math/big"
	"strconv"
)

// pisCofinsPct é a soma de PIS 1,65% e COFINS 7,6% no regime não-cumulativo.
// Medido: 15.7 mil linhas do ERP, e o crédito de 9,25% na fórmula de custo do
// próprio ERP confirma Lucro Real.
const pisCofinsPct = "9.25"

var (
	errValorNaoNumerico   = errors.New("fiscal: valor nao numerico")
	errPercentualInvalido = errors.New("fiscal: percentual invalido")
)

// TaxesForValue calcula o imposto de UMA linha de pedido a partir da célula da
// matriz do ERP.
//
//	base ICMS  = V x (1 - redbase)
//	ICMS       = base ICMS x aliquota
//	DIFAL      = V x aliqintdest - ICMS         (só quando redbase = 0)
//	FCP        = V x perc_fcp
//	PIS/COFINS = (V - ICMS - DIFAL - FCP) x 9,25%
//	Total      = ICMS + DIFAL + FCP + PIS/COFINS
//
// A base de PIS/COFINS desconta a família do ICMS porque o ICMS é repasse ao
// estado, não receita — STF, Tema 69 (RE 574.706). Não é interpretação nossa: o
// Sankhya grava exatamente essa base em 1.434 de 1.434 linhas não-ST.
// Consequência contraintuitiva: quanto maior o ICMS, menor o PIS/COFINS. Somar
// imposto por cima de imposto superestima a carga.
//
// O desconto usa o DINHEIRO apurado, não a soma das alíquotas. As duas formas
// coincidem quando redbase = 0, e só a primeira continua certa quando não.
//
// Nunca devolve zero para significar "não sei". Componente sem resposta carrega
// Motivo em português, nomeando a lacuna do cadastro do ERP.
func TaxesForValue(valorLinha float64, regra ICMSRule) Resultado {
	out := Resultado{Origem: OrigemMatriz}
	if !regra.VigenteDesde.IsZero() {
		desde := regra.VigenteDesde
		out.MatrizDesde = &desde
	}

	if regra.Ambiguo {
		return selaTudo(out, "Matriz de ICMS do ERP tem mais de uma linha aplicavel para este produto e destino. "+
			"A precedencia entre elas nao esta definida, entao nenhum valor pode ser afirmado.")
	}

	// ZERAR = 'S' anula a família do ICMS por decisão da própria matriz.
	if regra.Zerar {
		return semICMS(out, valorLinha)
	}

	// CODTRIB é nullable na origem — 22 linhas com origem MG. Nil não é 0:
	// 0 significa "tributado", e assumi-lo seria afirmar incidência.
	if regra.CodTrib == nil {
		return selaTudo(out, "Matriz de ICMS do ERP nao informa o codigo de tributacao (CODTRIB) "+
			"para este produto e destino.")
	}

	switch *regra.CodTrib {
	case 60:
		// Substituição tributária: o ICMS já foi recolhido na entrada pelo
		// substituto e está dentro do CUSSEMICM. Zero aqui é fato medido
		// (11.233/11.233 linhas), não ausência de resposta.
		return semICMS(out, valorLinha)
	case 0:
		// Tributado — segue.
	default:
		// 10 (41 linhas MG) e 40 (2 linhas). Comportamento NÃO medido.
		return selaTudo(out, fmt.Sprintf(
			"Matriz de ICMS do ERP usa CODTRIB %d para este produto e destino, "+
				"cujo tratamento ainda nao foi apurado.", *regra.CodTrib))
	}

	if regra.AliquotaPct == nil {
		return selaTudo(out, "Matriz de ICMS do ERP nao informa a aliquota da operacao para este produto e destino.")
	}

	// Redução da base da operação. Nil é zero: a célula existe e não declarou
	// redução. Não-zero em 63% das linhas MG.
	redBase := big.NewRat(0, 1)
	if regra.RedBasePct != nil {
		r, ok := new(big.Rat).SetString(*regra.RedBasePct)
		if !ok {
			return selaTudo(out, "Reducao de base invalida na matriz de ICMS do ERP: "+*regra.RedBasePct)
		}
		redBase = r
	}
	baseICMS, err := aplicaReducao(valorLinha, redBase)
	if err != nil {
		return selaTudo(out, "Falha ao aplicar a reducao de base da matriz de ICMS do ERP.")
	}
	icms, err := pctOf(baseICMS, *regra.AliquotaPct)
	if err != nil {
		return selaTudo(out, "Aliquota da operacao invalida na matriz de ICMS do ERP: "+*regra.AliquotaPct)
	}
	out.ICMS = Conhecido(icms)

	// FCP nulo com a célula presente é zero declarado: a matriz respondeu.
	out.FCP = Conhecido(0)
	if regra.PercFCPPct != nil {
		fcp, err := pctOf(valorLinha, *regra.PercFCPPct)
		if err != nil {
			return selaRestante(out, "Percentual de FCP invalido na matriz de ICMS do ERP: "+*regra.PercFCPPct)
		}
		out.FCP = Conhecido(fcp)
	}

	if regra.AliqIntDestPct == nil {
		return selaRestante(out, "Falta ALIQINTDEST na matriz de ICMS do ERP para este grupo com este destino. "+
			"Sem ela o DIFAL nao pode ser calculado.")
	}

	// A base do DIFAL é CHEIA dos dois lados. Medido: BASEDIFAL = BASE em
	// 347/347 linhas de venda não-ST com DIFAL, e PERCREDBASE = 0 em 347/347.
	// PERCREDBASEDEST é coluna morta — 0 em 4.601 de 4.601 linhas da TGFICM,
	// todas as origens.
	//
	// Fórmula do ERP, reconciliada contra o documento:
	//   VLRDIFALDEST = ROUND(BASEDIFAL x (ALIQINTDEST - ALIQUOTA) / 100, 2)
	// 312/347 ao centavo, 347/347 dentro de 5 centavos, erro agregado 0,00003%
	// sobre R$ 93.120 de DIFAL.
	//
	// Note que o crédito descontado é V x ALIQUOTA sobre a base CHEIA, não o
	// ICMS efetivamente apurado — os dois só coincidem quando REDBASE = 0. A
	// combinação REDBASE != 0 com CODTRIB = 0 tem ZERO observações no recorte
	// real (REDBASE != 0 aparece só como REDBASE = 100 junto de CODTRIB = 60,
	// que é a marcação de ST e já saiu no ramo acima). Sem observação, não se
	// afirma: vira pendência.
	if redBase.Sign() != 0 {
		return selaRestante(out, "Matriz de ICMS do ERP reduz a base desta operacao sem marca-la como "+
			"substituicao tributaria. Essa combinacao nao ocorreu em nenhuma venda apurada, "+
			"entao o DIFAL nao pode ser calculado.")
	}

	aid, ok := new(big.Rat).SetString(*regra.AliqIntDestPct)
	if !ok {
		return selaRestante(out, "Aliquota interna do destino invalida na matriz de ICMS do ERP: "+*regra.AliqIntDestPct)
	}
	aliq, ok := new(big.Rat).SetString(*regra.AliquotaPct)
	if !ok {
		return selaRestante(out, "Aliquota da operacao invalida na matriz de ICMS do ERP: "+*regra.AliquotaPct)
	}
	difal, err := pctOf(valorLinha, new(big.Rat).Sub(aid, aliq).FloatString(4))
	if err != nil {
		return selaRestante(out, "Falha ao compor o DIFAL.")
	}
	if difal < 0 {
		return selaRestante(out, fmt.Sprintf(
			"Matriz de ICMS do ERP tem aliquota interna do destino (%s%%) menor que a da operacao (%s%%). "+
				"O DIFAL ficaria negativo, o que indica defeito de cadastro.",
			*regra.AliqIntDestPct, *regra.AliquotaPct))
	}
	out.DIFAL = Conhecido(difal)

	return fechaComPisCofins(out, valorLinha)
}

// semICMS trata os dois casos em que a matriz declara que a família do ICMS não
// incide: ST (CODTRIB 60) e ZERAR = 'S'. Zero aqui é resposta, não lacuna.
func semICMS(out Resultado, valorLinha float64) Resultado {
	out.ICMS, out.DIFAL, out.FCP = Conhecido(0), Conhecido(0), Conhecido(0)
	return fechaComPisCofins(out, valorLinha)
}

// aplicaReducao devolve valorLinha x (1 - reducao), com reducao em pontos
// percentuais.
func aplicaReducao(valorLinha float64, reducaoPct *big.Rat) (float64, error) {
	v := new(big.Rat).SetFloat64(valorLinha)
	if v == nil {
		return 0, errValorNaoNumerico
	}
	fator := new(big.Rat).Quo(reducaoPct, cem)
	fator.Sub(big.NewRat(1, 1), fator)
	return strconv.ParseFloat(new(big.Rat).Mul(v, fator).FloatString(4), 64)
}

// fechaComPisCofins desconta do valor da linha a família do ICMS já apurada,
// aplica 9,25% sobre o que sobrou, e soma o total. Também deriva a carga de
// ICMS em pontos percentuais a partir do dinheiro.
func fechaComPisCofins(out Resultado, valorLinha float64) Resultado {
	carga := 0.0
	for _, c := range []Componente{out.ICMS, out.DIFAL, out.FCP} {
		if c.Valor == nil {
			return out
		}
		carga += *c.Valor
	}
	base, err := somaCentavos(valorLinha, -carga)
	if err != nil {
		return selaRestante(out, "Falha ao compor a base de PIS/COFINS.")
	}
	pc, err := pctOf(base, pisCofinsPct)
	if err != nil {
		return selaRestante(out, "Falha ao calcular PIS/COFINS.")
	}
	out.PisCofins = Conhecido(pc)

	if valorLinha != 0 {
		pct := new(big.Rat).Quo(new(big.Rat).SetFloat64(carga), new(big.Rat).SetFloat64(valorLinha))
		pct.Mul(pct, cem)
		if v, err := strconv.ParseFloat(pct.FloatString(3), 64); err == nil {
			out.CargaICMSPct = Conhecido(v)
		}
	}

	total, err := somaCentavos(carga, pc)
	if err != nil {
		return out
	}
	out.Total = Conhecido(total)
	return out
}

// somaCentavos soma dois valores em reais e devolve o resultado a centavos, sem
// acumular erro de float ao longo da cadeia.
func somaCentavos(a, b float64) (float64, error) {
	ra := new(big.Rat).SetFloat64(a)
	rb := new(big.Rat).SetFloat64(b)
	if ra == nil || rb == nil {
		return 0, errValorNaoNumerico
	}
	return strconv.ParseFloat(new(big.Rat).Add(ra, rb).FloatString(2), 64)
}

// selaTudo marca todos os componentes como desconhecidos com o mesmo motivo.
func selaTudo(out Resultado, motivo string) Resultado {
	d := Desconhecido(motivo)
	out.ICMS, out.DIFAL, out.FCP = d, d, d
	out.PisCofins, out.Total, out.CargaICMSPct = d, d, d
	return out
}

// selaRestante preserva o que já foi apurado e sela o que depende do que falta.
// O ICMS próprio continua sendo um fato mesmo quando o DIFAL não fecha.
func selaRestante(out Resultado, motivo string) Resultado {
	d := Desconhecido(motivo)
	out.DIFAL, out.CargaICMSPct = d, d
	out.PisCofins, out.Total = d, d
	return out
}
```

- [ ] **Step 4: Rode e confirme que passa**

```bash
cd apps/server_core && GOCACHE=$(pwd)/.gocache go test ./internal/modules/fiscal/domain/ -v
```

Esperado: `PASS` em todos.

- [ ] **Step 5: Controle negativo — prove que a base reduzida é o que está sendo testado**

Troque, em `fechaComPisCofins`, a base por `valorLinha` puro (sem descontar `carga`) e rode.

Esperado: `TestTaxesForValueBahiaSegueMatrizVigente` **falha** com `PIS/COFINS = 27.74, queria 22.05`. Se ele passar, o teste não está provando o Tema 69 e o número de 22,05 veio de outro lugar. Desfaça depois.

- [ ] **Step 6: Segundo controle negativo — prove que a cascata da célula muda é real**

Em `selaRestante`, pare de selar `PisCofins` (deixe-o ser calculado com carga zero) e rode.

Esperado: `TestTaxesForValueCelulaMudaSelaDifalEPisCofins` **falha**. Desfaça depois.

- [ ] **Step 7: Commit**

```bash
git add apps/server_core/internal/modules/fiscal/domain/
git commit -m "feat(fiscal): TaxesForValue com base de PIS/COFINS reduzida pela carga de ICMS"
```

---

## Task 9: Porta e adapter Postgres com leitura as-of

**Files:**
- Create: `apps/server_core/internal/modules/fiscal/ports/ports.go`
- Create: `apps/server_core/internal/modules/fiscal/adapters/postgres/matrix_reader.go`
- Test: `apps/server_core/internal/modules/fiscal/adapters/postgres/matrix_reader_integration_test.go`

- [ ] **Step 1: Escreva a porta**

```go
package ports

import (
	"context"
	"time"

	"marketplace-central/apps/server_core/internal/modules/fiscal/domain"
)

// TaxMatrixReader entrega a célula da matriz de ICMS vigente NA DATA pedida.
//
// asOf é a data do fato (data do pedido), não "agora". É o que faz corrigir o
// cadastro do ERP não apagar a divergência histórica.
//
// Célula ausente devolve ok = false, sem erro: não é falha, é o ERP não ter
// resposta, e quem chama transforma isso em pendência.
type TaxMatrixReader interface {
	RuleFor(ctx context.Context, tenantID string, produtoID string, ufDestino string, asOf time.Time) (rule domain.ICMSRule, ok bool, err error)
}
```

- [ ] **Step 2: Escreva os testes que falham**

```go
//go:build integration

package postgres_test

func TestRuleForLeVersaoVigenteNaData(t *testing.T) {
	ctx := context.Background()
	pool := openTestPool(t)
	defer pool.Close()
	tenantID := newTestTenant(t, pool)

	// Produto 15956, grupo 122.
	inserirProduto(t, pool, tenantID, "15956", 122)

	antes := time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)
	corte := time.Date(2026, 7, 20, 0, 0, 0, 0, time.UTC)
	depois := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)

	// Versao antiga: BA a 17,0 (defasada). Versao nova: 20,5 (corrigida).
	inserirCelula(t, pool, tenantID, "MG", "BA", 122, "7.000", "17.000",
		time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), &corte)
	inserirCelula(t, pool, tenantID, "MG", "BA", 122, "7.000", "20.500", corte, nil)

	r := postgres.NewMatrixReader(pool, "MG")

	velha, ok, err := r.RuleFor(ctx, tenantID, "15956", "BA", antes)
	if err != nil || !ok {
		t.Fatalf("pedido antigo: ok=%v err=%v", ok, err)
	}
	if *velha.AliqIntDestPct != "17.000" {
		t.Fatalf("pedido de marco leu %q; tem que ler a matriz que valia NA EPOCA (17,000)", *velha.AliqIntDestPct)
	}

	nova, ok, err := r.RuleFor(ctx, tenantID, "15956", "BA", depois)
	if err != nil || !ok {
		t.Fatalf("pedido novo: ok=%v err=%v", ok, err)
	}
	if *nova.AliqIntDestPct != "20.500" {
		t.Fatalf("pedido de agosto leu %q, queria 20,500", *nova.AliqIntDestPct)
	}
}

func TestRuleForProdutoSemGrupoNaoTemRegra(t *testing.T) {
	ctx := context.Background()
	pool := openTestPool(t)
	defer pool.Close()
	tenantID := newTestTenant(t, pool)

	inserirProdutoSemGrupo(t, pool, tenantID, "99999")
	r := postgres.NewMatrixReader(pool, "MG")

	_, ok, err := r.RuleFor(ctx, tenantID, "99999", "BA", time.Now())
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if ok {
		t.Fatal("produto sem grupo_icms nao pode ter regra; sem grupo nao ha celula")
	}
}

func TestRuleForOutroTenantNaoVaza(t *testing.T) {
	ctx := context.Background()
	pool := openTestPool(t)
	defer pool.Close()
	meu := newTestTenant(t, pool)
	alheio := newTestTenant(t, pool)

	inserirProduto(t, pool, meu, "15956", 122)
	inserirCelula(t, pool, alheio, "MG", "BA", 122, "7.000", "20.500",
		time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), nil)

	r := postgres.NewMatrixReader(pool, "MG")
	_, ok, err := r.RuleFor(ctx, meu, "15956", "BA", time.Now())
	if err != nil {
		t.Fatalf("erro: %v", err)
	}
	if ok {
		t.Fatal("celula de outro tenant foi lida — falta predicado de tenant_id")
	}
}
```

Escreva os helpers `inserirProduto`, `inserirProdutoSemGrupo` e `inserirCelula` no mesmo arquivo, como `INSERT` diretos. Nada de fixture compartilhada: fixture simétrica esconde defeito de assimetria.

- [ ] **Step 3: Rode e confirme que falha**

```bash
cd apps/server_core && GOCACHE=$(pwd)/.gocache go test -tags integration ./internal/modules/fiscal/adapters/postgres/ -v
```

Esperado: `FAIL` — `NewMatrixReader` não existe.

- [ ] **Step 4: Implemente**

`NewMatrixReader(pool *pgxpool.Pool, ufOrigem string)`. Uma consulta só, com join:

```sql
SELECT m.codtrib, m.zerar,
       m.aliquota::text, m.redbase::text,
       m.aliqintdest::text, m.perc_red_base_dest::text, m.perc_fcp::text,
       m.ambiguo, m.vigente_desde
FROM products_mirror p
JOIN icms_matrix_mirror m
  ON  m.tenant_id  = p.tenant_id
  AND m.grupo_icms = p.grupo_icms
WHERE p.tenant_id       = $1
  AND p.codigo_produto  = $2
  AND p.grupo_icms IS NOT NULL
  AND m.uf_origem       = $3
  AND m.uf_destino      = $4
  AND m.vigente_desde  <= $5
  AND (m.vigente_ate IS NULL OR m.vigente_ate > $5)
LIMIT 1
```

`::text` nos numéricos: os percentuais viajam como string até o `big.Rat`, sem passar por float. `pgx` escaneia `NUMERIC` para `float64` sem reclamar, e é exatamente aí que 20,5 vira 20,499999.

Zero linhas → `ok = false`, sem erro.

- [ ] **Step 5: Rode e confirme que passa**

```bash
cd apps/server_core && GOCACHE=$(pwd)/.gocache go test -tags integration ./internal/modules/fiscal/adapters/postgres/ -v
```

Esperado: `PASS` nos três.

- [ ] **Step 6: Controle negativo — prove que o as-of morde**

Remova as duas linhas `AND m.vigente_desde <= $5` e `AND (m.vigente_ate IS NULL OR m.vigente_ate > $5)` e rode.

Esperado: `TestRuleForLeVersaoVigenteNaData` **falha** — vai ler qualquer uma das duas versões. Se passar, o `LIMIT 1` está escondendo o problema e o teste não prova o as-of. Desfaça depois.

- [ ] **Step 7: Segundo controle negativo — prove o isolamento de tenant**

Remova `AND p.tenant_id = $1` e rode.

Esperado: `TestRuleForOutroTenantNaoVaza` **falha**. Desfaça depois.

- [ ] **Step 8: Commit**

```bash
git add apps/server_core/internal/modules/fiscal/ports apps/server_core/internal/modules/fiscal/adapters
git commit -m "feat(fiscal): leitura as-of da matriz espelhada com isolamento de tenant"
```

---

## Task 10: Job de sync da matriz e fiação

**Files:**
- Create: `apps/server_core/internal/modules/internal_read/application/icms_matrix_sync.go`
- Modify: `apps/server_core/internal/composition/root.go` (junto do bloco do `NewProductsScheduler`, hoje em `root.go:695`)
- Test: `apps/server_core/internal/modules/internal_read/application/icms_matrix_sync_test.go`

- [ ] **Step 1: Escreva o teste que falha**

```go
func TestSyncMatrixResolveCadaCelulaAntesDeGravar(t *testing.T) {
	origem := &fakeMatrixSource{lines: map[oracle.MatrixKey][]domain.MatrixLine{
		{UFDestino: "BA", GrupoICMS: 122}: {
			{Restricao1Tipo: "N", Restricao1Codigo: 0, Restricao2Tipo: "I", Restricao2Codigo: 122,
				Sequencia: 1, CodTrib: ptrInt(0), AliquotaPct: "7.0", RedBasePct: "0", AliqIntDestPct: "20.5"},
		},
		{UFDestino: "SP", GrupoICMS: 122}: {
			{Restricao1Tipo: "N", Restricao1Codigo: 0, Restricao2Tipo: "I", Restricao2Codigo: 122,
				Sequencia: 1, CodTrib: ptrInt(0), AliquotaPct: "7.0", RedBasePct: "0", AliqIntDestPct: "18.0"},
			{Restricao1Tipo: "S", Restricao1Codigo: -1, Restricao2Tipo: "I", Restricao2Codigo: 122,
				Sequencia: 1, CodTrib: ptrInt(0), AliquotaPct: "7.0", RedBasePct: "0", AliqIntDestPct: "19.0"},
		},
	}}
	destino := &fakeMatrixWriter{}

	if err := application.SyncICMSMatrix(context.Background(), origem, destino, "tenant-1", "MG"); err != nil {
		t.Fatalf("sync: %v", err)
	}

	ba := destino.cellFor("BA", 122)
	if ba.Ambiguo {
		t.Fatal("BA tem uma linha so; nao pode ser ambigua")
	}
	if *ba.AliqIntDestPct != "20.5" {
		t.Fatalf("BA gravou %q, queria 20.5", *ba.AliqIntDestPct)
	}

	sp := destino.cellFor("SP", 122)
	if !sp.Ambiguo {
		t.Fatal("SP tem duas linhas sobreviventes; tem que gravar ambigua")
	}
	if sp.AliqIntDestPct != nil {
		t.Fatalf("celula ambigua gravou aliquota %q — nao pode escolher", *sp.AliqIntDestPct)
	}
	if sp.LinhasCandidatas != 2 {
		t.Fatalf("LinhasCandidatas = %d, queria 2", sp.LinhasCandidatas)
	}
}
```

Escreva `fakeMatrixSource` e `fakeMatrixWriter` no mesmo arquivo.

- [ ] **Step 2: Rode e confirme que falha**

```bash
cd apps/server_core && GOCACHE=$(pwd)/.gocache go test ./internal/modules/internal_read/application/ -run TestSyncMatrix -v
```

Esperado: `FAIL` — `SyncICMSMatrix` não existe.

- [ ] **Step 3: Implemente**

`SyncICMSMatrix` lê tudo da origem, chama `domain.ResolveCell` para cada chave, monta as `mirror.MatrixCell` (com `nil` nos percentuais quando ambígua) e chama `ApplyMatrix` uma vez, com o lote inteiro.

- [ ] **Step 4: Rode e confirme que passa**

```bash
cd apps/server_core && GOCACHE=$(pwd)/.gocache go test ./internal/modules/internal_read/application/ -run TestSyncMatrix -v
```

Esperado: `PASS`.

- [ ] **Step 5: Fie no `root.go`**

Ao lado do `NewProductsScheduler` (`root.go:695`), registre um scheduler de **24 horas** para `SyncICMSMatrix`, ligado ao mesmo `oracleDB`. Siga exatamente o padrão do scheduler de produtos, inclusive o comportamento com `oracleDB == nil`: o app **tem que continuar subindo sem Oracle** (`root.go:446-452` já trata isso com WARN). Sem Oracle, o espelho fica como está e o cálculo lê o que já foi sincronizado.

- [ ] **Step 6: Compile e rode a suíte**

```bash
cd apps/server_core && GOCACHE=$(pwd)/.gocache go build ./... && GOCACHE=$(pwd)/.gocache go test ./internal/composition/... -v
```

Esperado: build sem saída, testes `PASS`.

- [ ] **Step 7: Commit**

```bash
git add apps/server_core/internal/modules/internal_read/application apps/server_core/internal/composition/root.go
git commit -m "feat(fiscal): sync diario da matriz de ICMS resolvendo a celula antes de gravar"
```

---

## Task 11: Porta de `orders` redesenhada para item

A porta de hoje (`orders/ports/tax_reader.go`) pergunta o imposto do **pedido inteiro** a partir do total. Imposto é por linha, porque o grupo fiscal é do produto. Um pedido com dois produtos de grupos diferentes não tem uma alíquota só.

**Files:**
- Modify: `apps/server_core/internal/modules/orders/ports/tax_reader.go`
- Create: `apps/server_core/internal/modules/orders/adapters/fiscal/reader.go`
- Test: `apps/server_core/internal/modules/orders/adapters/fiscal/reader_test.go`

- [ ] **Step 1: Reescreva a porta**

```go
package ports

import (
	"context"
	"time"
)

// TaxComponent é um valor fiscal que pode não existir. Espelha o Componente do
// módulo fiscal sem que orders precise importar o domínio dele — mesmo padrão
// que ItemCost usa para o custo.
type TaxComponent struct {
	Valor  *float64
	Motivo string
}

// ItemTaxInput identifica uma linha do pedido para o cálculo fiscal.
type ItemTaxInput struct {
	ItemIdentifier string
	ProdutoID      string
	Valor          float64
}

// ItemTaxes é o imposto de uma linha.
type ItemTaxes struct {
	ItemIdentifier string
	ICMS           TaxComponent
	Difal          TaxComponent
	FCP            TaxComponent
	PisCofins      TaxComponent
	Total          TaxComponent
	CargaICMSPct   TaxComponent
	Origem         string
	MatrizDesde    *time.Time
}

// TaxReader calcula o imposto POR LINHA do pedido.
//
// asOf é a data do pedido, não "agora": a matriz do ERP é versionada e um
// pedido antigo tem que ser lido com a matriz que valia na época.
//
// destinoUF vazio é honest-unknown, não erro — sem destino não há DIFAL a
// afirmar. Um reader nil deixa tudo desconhecido em vez de entrar em pânico,
// no mesmo idioma de CostReader e ShipmentReader.
type TaxReader interface {
	TaxesForItems(ctx context.Context, itens []ItemTaxInput, destinoUF string, asOf time.Time) ([]ItemTaxes, error)
}
```

Apague `OrderTaxes` deste arquivo.

- [ ] **Step 2: Escreva o teste do adapter que falha**

```go
func TestTaxesForItemsUsaMatrizEDataDoPedido(t *testing.T) {
	matriz := &fakeMatrix{rule: domain.ICMSRule{
		CodTrib: i(0), AliquotaPct: s("7.0"), RedBasePct: s("0"), AliqIntDestPct: s("20.5"),
		VigenteDesde: time.Date(2026, 7, 20, 0, 0, 0, 0, time.UTC),
	}}
	r := fiscaladapter.NewReader(matriz, "tenant-1")

	got, err := r.TaxesForItems(context.Background(),
		[]ports.ItemTaxInput{{ItemIdentifier: "L1", ProdutoID: "15956", Valor: 299.90}},
		"BA", time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("TaxesForItems: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("linhas = %d, queria 1", len(got))
	}
	if got[0].Total.Valor == nil || *got[0].Total.Valor != 83.53 {
		t.Fatalf("Total = %v, queria 83.53", got[0].Total.Valor)
	}
	if matriz.asOfRecebido != time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC) {
		t.Fatalf("asOf repassado = %v; a data do PEDIDO tem que chegar na matriz, nao time.Now()", matriz.asOfRecebido)
	}
}

func TestTaxesForItemsSemDestinoNaoInventaImposto(t *testing.T) {
	r := fiscaladapter.NewReader(&fakeMatrix{}, "tenant-1")
	got, err := r.TaxesForItems(context.Background(),
		[]ports.ItemTaxInput{{ItemIdentifier: "L1", ProdutoID: "15956", Valor: 299.90}},
		"", time.Now())
	if err != nil {
		t.Fatalf("erro: %v", err)
	}
	if got[0].Total.Valor != nil {
		t.Fatalf("Total = %v sem destino — sem UF de destino nao ha imposto a afirmar", *got[0].Total.Valor)
	}
	if got[0].Total.Motivo == "" {
		t.Fatal("faltou motivo")
	}
}

func TestTaxesForItemsCelulaAusenteVirapendencia(t *testing.T) {
	r := fiscaladapter.NewReader(&fakeMatrix{ausente: true}, "tenant-1")
	got, err := r.TaxesForItems(context.Background(),
		[]ports.ItemTaxInput{{ItemIdentifier: "L1", ProdutoID: "15956", Valor: 299.90}},
		"BA", time.Now())
	if err != nil {
		t.Fatalf("erro: %v", err)
	}
	if got[0].ICMS.Valor != nil {
		t.Fatal("celula ausente nao autoriza nenhum numero")
	}
	if !strings.Contains(got[0].ICMS.Motivo, "Sankhya") && !strings.Contains(got[0].ICMS.Motivo, "ERP") {
		t.Fatalf("motivo = %q; a pendencia tem que apontar o cadastro do ERP", got[0].ICMS.Motivo)
	}
}
```

- [ ] **Step 3: Rode e confirme que falha**

```bash
cd apps/server_core && GOCACHE=$(pwd)/.gocache go test ./internal/modules/orders/adapters/fiscal/ -v
```

Esperado: `FAIL` — pacote não existe.

- [ ] **Step 4: Implemente o adapter**

Para cada item: se `destinoUF == ""` ou `ProdutoID == ""`, sela tudo com motivo. Senão chama `RuleFor`; `ok == false` sela tudo com *"Matriz de ICMS do ERP nao tem linha para este produto com destino <UF>. Cadastre no Sankhya."*; com regra, chama `domain.TaxesForValue` e mapeia `domain.Componente` → `ports.TaxComponent`.

- [ ] **Step 5: Rode e confirme que passa**

```bash
cd apps/server_core && GOCACHE=$(pwd)/.gocache go test ./internal/modules/orders/adapters/fiscal/ -v
```

Esperado: `PASS` nos três.

- [ ] **Step 6: Controle negativo — prove que o `asOf` é repassado**

Troque, no adapter, o `asOf` recebido por `time.Now()` na chamada a `RuleFor` e rode.

Esperado: `TestTaxesForItemsUsaMatrizEDataDoPedido` **falha** na asserção de `asOfRecebido`. É o controle que impede o bug silencioso de ler a matriz de hoje para um pedido de março. Desfaça depois.

- [ ] **Step 7: Commit**

```bash
git add apps/server_core/internal/modules/orders/ports/tax_reader.go apps/server_core/internal/modules/orders/adapters/fiscal/
git commit -m "feat(fiscal): porta de imposto de orders passa a ser por item e le a matriz as-of"
```

---

## Task 12: `enrich_service` usa o imposto por item

**Files:**
- Modify: `apps/server_core/internal/modules/orders/application/enrich_service.go:405-423`
- Test: `apps/server_core/internal/modules/orders/application/enrich_service_test.go`

O padrão a seguir já existe no mesmo arquivo: `resolveItemCosts` (`:487-521`) faz exatamente isto para custo — itera itens, coleta desconhecidos, nunca aborta o pedido por causa de uma linha. Leia essa função antes de escrever.

- [ ] **Step 1: Escreva o teste que falha**

```go
func TestEnrichSomaImpostoDosItens(t *testing.T) {
	svc := newEnrichServiceComTaxes(t, &fakeTaxReader{
		porItem: map[string]ports.ItemTaxes{
			"L1": {ItemIdentifier: "L1", Total: comp(83.53), ICMS: comp(20.99),
				Difal: comp(40.49), FCP: comp(0), PisCofins: comp(22.05)},
		},
	})

	got := svc.Enrich(context.Background(), pedidoBahia(t))

	if got.Decomposicao.Imposto == nil || *got.Decomposicao.Imposto != 83.53 {
		t.Fatalf("imposto = %v, queria 83.53", got.Decomposicao.Imposto)
	}
	// margem = 299,90 - 40,49 - 23,65 - 83,53 - 154,53
	if got.Decomposicao.Margem == nil || *got.Decomposicao.Margem != -2.30 {
		t.Fatalf("margem = %v, queria -2.30", got.Decomposicao.Margem)
	}
}

func TestEnrichNaoClampaMargemNegativaEmZero(t *testing.T) {
	svc := newEnrichServiceComTaxes(t, &fakeTaxReader{
		porItem: map[string]ports.ItemTaxes{
			"L1": {ItemIdentifier: "L1", Total: comp(83.53)},
		},
	})

	got := svc.Enrich(context.Background(), pedidoBahia(t))

	if got.Decomposicao.Margem == nil {
		t.Fatal("margem nula")
	}
	if *got.Decomposicao.Margem >= 0 {
		t.Fatalf("margem = %v; este pedido da PREJUIZO e a tela tem que dizer isso. "+
			"Clamp em zero transforma prejuizo em empate", *got.Decomposicao.Margem)
	}
}

func TestEnrichItemSemImpostoNaoZeraOPedido(t *testing.T) {
	svc := newEnrichServiceComTaxes(t, &fakeTaxReader{
		porItem: map[string]ports.ItemTaxes{
			"L1": {ItemIdentifier: "L1", Total: ports.TaxComponent{Motivo: "Falta ALIQINTDEST na matriz de ICMS do ERP."}},
		},
	})

	got := svc.Enrich(context.Background(), pedidoBahia(t))

	if got.Decomposicao.Imposto != nil {
		t.Fatalf("imposto = %v; item sem imposto conhecido nao pode virar zero", *got.Decomposicao.Imposto)
	}
	if got.Decomposicao.Margem != nil {
		t.Fatalf("margem = %v; sem imposto nao ha margem a afirmar", *got.Decomposicao.Margem)
	}
	if !contains(got.Incompleto, "imposto") {
		t.Fatalf("incompleto = %v; tem que nomear 'imposto'", got.Incompleto)
	}
}
```

- [ ] **Step 2: Rode e confirme que falha**

```bash
cd apps/server_core && GOCACHE=$(pwd)/.gocache go test ./internal/modules/orders/application/ -run TestEnrich -v
```

Esperado: `FAIL` — não compila, a porta mudou de forma.

- [ ] **Step 3: Substitua `resolveTaxes` por `resolveItemTaxes`**

```go
// resolveItemTaxes pede o imposto de cada linha do pedido à porta fiscal,
// usando a data do PEDIDO (não "agora"), porque a matriz do ERP é versionada.
//
// Sem reader, sem itens ou erro do reader degradam para tudo desconhecido
// (ADR-17). Um imposto que não conseguimos apurar nunca vira zero.
func (s EnrichService) resolveItemTaxes(ctx context.Context, order domain.OrderReadModel, shipment *ShipmentEnrichment) []ports.ItemTaxes {
	if s.taxes == nil || len(order.Items) == 0 {
		return nil
	}
	effectiveAt, ok := orderEffectiveAt(order)
	if !ok {
		return nil
	}
	var destinoUF string
	if shipment != nil && shipment.DestinationUF != nil {
		destinoUF = *shipment.DestinationUF
	}
	inputs := make([]ports.ItemTaxInput, 0, len(order.Items))
	for _, item := range order.Items {
		in := ports.ItemTaxInput{ItemIdentifier: itemIdentifier(item)}
		if item.InternalProductID != nil {
			in.ProdutoID = strconv.Itoa(*item.InternalProductID)
		}
		if item.LineTotal != nil {
			in.Valor = *item.LineTotal
		}
		inputs = append(inputs, in)
	}
	out, err := s.taxes.TaxesForItems(ctx, inputs, destinoUF, effectiveAt)
	if err != nil {
		s.logger.Warn("orders: tax lookup failed",
			"provider_order_id", order.ProviderOrderID,
			"destino_uf", destinoUF,
			"error", err,
		)
		return nil
	}
	return out
}

// sumItemTaxes soma o imposto das linhas. Uma linha sem imposto conhecido
// torna o total do pedido desconhecido: somar só o que se sabe produziria um
// imposto menor que o real, com cara de exato.
func sumItemTaxes(itens []ports.ItemTaxes) (*float64, bool) {
	if len(itens) == 0 {
		return nil, false
	}
	total := 0.0
	for _, it := range itens {
		if it.Total.Valor == nil {
			return nil, false
		}
		total += *it.Total.Valor
	}
	arredondado := math.Round(total*100) / 100
	return &arredondado, true
}
```

Ajuste o ponto de montagem da decomposição (`enrich_service.go:395-400`, onde hoje entram `Imposto` e `Difal`) para usar `sumItemTaxes`, e acrescente `"imposto"` a `Incompleto` quando o total não for conhecido. **Nenhum `math.Max(0, ...)` no caminho da margem.** Se você encontrar um clamp já existente, mande `ESCALATION` — é achado de legacy, não conserto seu.

Se `item.LineTotal` não existir com esse nome no modelo, use o campo real e mande `ESCALATION` avisando qual é. Não invente.

- [ ] **Step 4: Rode e confirme que passa**

```bash
cd apps/server_core && GOCACHE=$(pwd)/.gocache go test ./internal/modules/orders/application/ -v
```

Esperado: `PASS`.

- [ ] **Step 5: Controle negativo — prove que a soma parcial é recusada**

Em `sumItemTaxes`, troque o `return nil, false` do item desconhecido por `continue` e rode.

Esperado: `TestEnrichItemSemImpostoNaoZeraOPedido` **falha**. Desfaça depois.

- [ ] **Step 6: Commit**

```bash
git add apps/server_core/internal/modules/orders/application/
git commit -m "feat(fiscal): decomposicao de pedido soma imposto por item e recusa soma parcial"
```

---

## Task 13: Contrato — OpenAPI e SDK no mesmo commit

**Files:**
- Modify: `contracts/api/marketplace-central.openapi.yaml`
- Modify: `packages/sdk-runtime/src/index.ts`
- Modify: `apps/server_core/internal/modules/orders/transport/http_handler.go:428-439`

- [ ] **Step 1: Acrescente os schemas ao OpenAPI**

Em `components.schemas`, ao lado de `OrderDecomposicao`:

```yaml
    ValorFiscal:
      type: object
      description: >-
        Valor fiscal que pode não existir. valor null com motivo preenchido é
        honest-unknown (ADR-17) — nunca zero para significar "não sei".
      required: [valor, motivo]
      properties:
        valor:
          type: number
          format: double
          nullable: true
        motivo:
          type: string
          nullable: true
          description: Pendência em português, nomeando a lacuna do cadastro do ERP.

    OrderFiscal:
      type: object
      description: >-
        Imposto da venda, decomposto. Nesta fatia origem é sempre matriz_erp
        (estimativa ex-ante pela matriz de ICMS do ERP). nota_erp já consta do
        enum porque a fatia de imposto realizado reusa este mesmo bloco.
      required: [origem, matriz_desde, icms, difal, fcp, pis_cofins, total, carga_icms_pct]
      properties:
        origem:
          type: string
          enum: [matriz_erp, nota_erp]
        matriz_desde:
          type: string
          format: date-time
          nullable: true
        icms:          { $ref: '#/components/schemas/ValorFiscal' }
        difal:         { $ref: '#/components/schemas/ValorFiscal' }
        fcp:           { $ref: '#/components/schemas/ValorFiscal' }
        pis_cofins:    { $ref: '#/components/schemas/ValorFiscal' }
        total:         { $ref: '#/components/schemas/ValorFiscal' }
        carga_icms_pct:{ $ref: '#/components/schemas/ValorFiscal' }
```

Acrescente `fiscal: { $ref: '#/components/schemas/OrderFiscal' }` ao schema do pedido em `/orders/{provider_order_id}`, e marque o objeto `difal` existente com `deprecated: true` e a descrição *"Substituído por `fiscal`. Sai na fatia de imposto realizado."*

- [ ] **Step 2: Acrescente os tipos ao SDK**

Em `packages/sdk-runtime/src/index.ts`, junto de `OrderDecomposicao` (perto da linha 695):

```ts
export interface ValorFiscal {
  valor: number | null;
  motivo: string | null;
}

export interface OrderFiscal {
  origem: "matriz_erp" | "nota_erp";
  matriz_desde: string | null;
  icms: ValorFiscal;
  difal: ValorFiscal;
  fcp: ValorFiscal;
  pis_cofins: ValorFiscal;
  total: ValorFiscal;
  carga_icms_pct: ValorFiscal;
}
```

E `fiscal: OrderFiscal;` na interface do pedido de detalhe.

- [ ] **Step 3: Emita o bloco no handler**

Em `http_handler.go:428-439`, onde a decomposição é montada, monte também `fiscal` a partir das `ItemTaxes`. Todo membro sai preenchido — `valor` ou `motivo`, nunca os dois nulos.

- [ ] **Step 4: Verifique os dois lados**

```bash
cd apps/web && npx --no-install tsc --noEmit -p tsconfig.json
cd apps/server_core && GOCACHE=$(pwd)/.gocache go test ./internal/modules/orders/... -v
```

Esperado: `tsc` sem erro, testes `PASS`.

> `npx --no-install tsc` passa vacuosamente se `node_modules` não estiver instalado na raiz do worktree. Se a saída vier vazia rápido demais, confirme que `apps/web/node_modules` existe. Instalação de dependência é `REQUEST` ao hub, não iniciativa sua.

- [ ] **Step 5: Commit — OpenAPI e SDK JUNTOS**

```bash
git add contracts/api/marketplace-central.openapi.yaml packages/sdk-runtime/src/index.ts apps/server_core/internal/modules/orders/transport/http_handler.go
git commit -m "feat(fiscal): contrato do bloco fiscal do pedido, OpenAPI e SDK no mesmo commit"
```

---

## Task 14: Frente — seção Fiscal no drawer

**Files:**
- Create: `apps/web/src/pages/pedidos/FiscalSection.tsx`
- Create: `apps/web/src/pages/pedidos/FiscalSection.test.tsx`
- Modify: `apps/web/src/pages/pedidos/PedidoDrawer.tsx:133-180`

- [ ] **Step 1: Escreva o teste que falha**

```tsx
import { render, screen } from "@testing-library/react";
import { FiscalSection } from "./FiscalSection";

const conhecido = (valor: number) => ({ valor, motivo: null });
const pendente = (motivo: string) => ({ valor: null, motivo });

it("mostra os quatro componentes e a vigencia da matriz", () => {
  render(<FiscalSection fiscal={{
    origem: "matriz_erp",
    matriz_desde: "2026-07-20T14:42:15Z",
    icms: conhecido(20.99),
    difal: conhecido(40.49),
    fcp: conhecido(0),
    pis_cofins: conhecido(22.05),
    total: conhecido(83.53),
    carga_icms_pct: conhecido(20.5),
  }} />);

  expect(screen.getByText("R$ 20,99")).toBeInTheDocument();
  expect(screen.getByText("R$ 40,49")).toBeInTheDocument();
  expect(screen.getByText("R$ 22,05")).toBeInTheDocument();
  expect(screen.getByText(/estimativa pela matriz do ERP/i)).toBeInTheDocument();
  expect(screen.getByText(/20\/07\/2026/)).toBeInTheDocument();
});

it("mostra a pendencia em vez de zero quando o ERP nao responde", () => {
  const motivo = "Falta ALIQINTDEST na matriz de ICMS do ERP para este grupo com este destino.";
  render(<FiscalSection fiscal={{
    origem: "matriz_erp",
    matriz_desde: null,
    icms: conhecido(20.99),
    difal: pendente(motivo),
    fcp: conhecido(0),
    pis_cofins: pendente(motivo),
    total: pendente(motivo),
    carga_icms_pct: pendente(motivo),
  }} />);

  // R$ 0,00 no lugar de desconhecido e o defeito que esta fatia existe para curar.
  expect(screen.queryByText("R$ 0,00")).not.toBeInTheDocument();
  expect(screen.getAllByTitle(motivo).length).toBeGreaterThan(0);
});
```

`R$ 0,00` aparece legitimamente no FCP do primeiro teste — por isso a asserção de ausência está **só no segundo**, onde nenhum componente vale zero de verdade.

- [ ] **Step 2: Rode e confirme que falha**

```bash
cd apps/web && npx vitest run src/pages/pedidos/FiscalSection.test.tsx
```

Esperado: `FAIL` — o módulo não existe.

- [ ] **Step 3: Implemente**

Uma linha por componente, usando o `UnknownValue` que já existe no projeto (veja como `PedidoDrawer.tsx:158` o usa) com `hint = motivo`. Rodapé com *"estimativa pela matriz do ERP"* e, quando `matriz_desde` não for nulo, *"vigência DD/MM/AAAA"*.

- [ ] **Step 4: Rode e confirme que passa**

```bash
cd apps/web && npx vitest run src/pages/pedidos/FiscalSection.test.tsx
```

Esperado: `PASS`.

- [ ] **Step 5: Troque a seção DIFAL legada no drawer**

Em `PedidoDrawer.tsx:133-180`, remova o bloco `Decomposição + DIFAL` que lê `order.difal` e coloque `<FiscalSection fiscal={order.fiscal} />`. As constantes `difalHint` e as `DecompRow` de DIFAL saem junto.

- [ ] **Step 6: Rode a suíte de pedidos e o tsc**

```bash
cd apps/web && npx vitest run src/pages/pedidos/ && npx --no-install tsc --noEmit -p tsconfig.json
```

Esperado: `PASS` e `tsc` limpo. Testes que ainda montam o `difal` legado precisam ser atualizados — atualize-os, não os apague.

- [ ] **Step 7: Commit**

```bash
git add apps/web/src/pages/pedidos/
git commit -m "feat(fiscal): secao fiscal do drawer de pedido com pendencia em vez de zero"
```

---

## Task 15: Card de saúde do espelho em `/integracoes`

**Files:**
- Create: `apps/web/src/pages/integracoes/IcmsMatrixHealthCard.tsx`
- Create: `apps/web/src/pages/integracoes/IcmsMatrixHealthCard.test.tsx`
- Modify: `apps/web/src/pages/integracoes/IntegracoesPage.tsx`
- Modify: `contracts/api/marketplace-central.openapi.yaml` (endpoint `GET /config/icms-matrix/health`)
- Modify: `packages/sdk-runtime/src/index.ts`

- [ ] **Step 1: Escreva o teste que falha**

```tsx
it("mostra contagem de celulas, ambiguas e mudas", () => {
  render(<IcmsMatrixHealthCard health={{
    synced_at: "2026-08-02T03:00:00Z",
    celulas: 210,
    ambiguas: 0,
    mudas: 20,
  }} />);

  expect(screen.getByText("210")).toBeInTheDocument();
  expect(screen.getByText(/20 sem al[íi]quota interna/i)).toBeInTheDocument();
});

it("diz quando o espelho nunca sincronizou, em vez de mostrar zero", () => {
  render(<IcmsMatrixHealthCard health={{ synced_at: null, celulas: 0, ambiguas: 0, mudas: 0 }} />);
  expect(screen.getByText(/nunca sincronizou/i)).toBeInTheDocument();
});
```

O segundo teste é o que separa "espelho vazio" de "espelho zerado" — sem ele, um sync que nunca rodou parece um ERP sem células.

- [ ] **Step 2: Rode e confirme que falha**

```bash
cd apps/web && npx vitest run src/pages/integracoes/IcmsMatrixHealthCard.test.tsx
```

Esperado: `FAIL`.

- [ ] **Step 3: Implemente o endpoint, o tipo do SDK e o card**

`GET /config/icms-matrix/health` devolve `{ synced_at, celulas, ambiguas, mudas }`, contando só as versões abertas (`vigente_ate IS NULL`) do tenant. Um `POST /config/icms-matrix/sync` dispara o sync manual. Siga o `SyncHealthCard.tsx` existente para a forma visual, e registre o card em `IntegracoesPage.tsx`.

- [ ] **Step 4: Rode e confirme que passa**

```bash
cd apps/web && npx vitest run src/pages/integracoes/ && npx --no-install tsc --noEmit -p tsconfig.json
```

Esperado: `PASS` e `tsc` limpo.

- [ ] **Step 5: Commit**

```bash
git add apps/web/src/pages/integracoes/ contracts/api/marketplace-central.openapi.yaml packages/sdk-runtime/src/index.ts apps/server_core/
git commit -m "feat(fiscal): saude do espelho da matriz de ICMS em /integracoes com sync manual"
```

---

## Task 16: Remoção do `pricingtax` e prova de que a dependência morreu

**Files:**
- Delete: `apps/server_core/internal/modules/orders/adapters/pricingtax/reader.go`
- Delete: `apps/server_core/internal/modules/orders/adapters/pricingtax/reader_test.go`
- Modify: `apps/server_core/internal/composition/root.go:618-621`
- Test: `apps/server_core/tests/unit/orders_no_pricing_dependency_test.go`

- [ ] **Step 1: Escreva o teste que falha**

```go
package unit

import (
	"os/exec"
	"strings"
	"testing"
)

// TestOrdersNaoDependeDePricing prova por CONSTRUCAO que a dependencia
// invertida morreu. Grep em import nao serve: uma cadeia indireta some do
// grep e continua no binario. go list -deps ve a arvore inteira.
//
// Escopo: pacotes de PRODUCAO. Testes podem importar pricing (o alarme de
// divergencia de arredondamento, D-30, faz isso de proposito).
func TestOrdersNaoDependeDePricing(t *testing.T) {
	out, err := exec.Command("go", "list", "-deps",
		"marketplace-central/apps/server_core/internal/modules/orders/...").Output()
	if err != nil {
		t.Fatalf("go list: %v", err)
	}
	for _, linha := range strings.Split(string(out), "\n") {
		if strings.Contains(linha, "/modules/pricing/") {
			t.Fatalf("orders ainda depende de %s — a inversao que esta fatia existe para desfazer continua de pe", linha)
		}
	}
}
```

- [ ] **Step 2: Rode e confirme que falha**

```bash
cd apps/server_core && GOCACHE=$(pwd)/.gocache go test ./tests/unit/ -run TestOrdersNaoDependeDePricing -v
```

Esperado: `FAIL`, nomeando `.../modules/pricing/domain`. Se **passar** agora, o teste está errado ou já foi tarde demais — investigue antes de seguir.

- [ ] **Step 3: Apague o `pricingtax` e refaça a fiação**

```bash
git rm apps/server_core/internal/modules/orders/adapters/pricingtax/reader.go apps/server_core/internal/modules/orders/adapters/pricingtax/reader_test.go
```

Em `root.go:618-621`, troque

```go
orderspricingtax.NewReader(pricingpostgres.NewCalcRepository(pool), tenantID)
```

por

```go
ordersfiscal.NewReader(fiscalpostgres.NewMatrixReader(pool, "MG"), tenantID)
```

e remova os imports que ficaram órfãos. A origem `"MG"` é a dívida **D-29** — deixe o comentário nomeando-a.

- [ ] **Step 4: Rode e confirme que passa**

```bash
cd apps/server_core && GOCACHE=$(pwd)/.gocache go build ./... && GOCACHE=$(pwd)/.gocache go test ./tests/unit/ -run TestOrdersNaoDependeDePricing -v
```

Esperado: build limpo e `PASS`.

- [ ] **Step 5: Rode a suíte inteira**

```bash
cd apps/server_core && GOCACHE=$(pwd)/.gocache go test ./... 2>&1 | tail -40
```

Esperado: nenhum `FAIL`. Qualquer pacote quebrado aqui é consequência da porta ter mudado de forma — conserte, e se algum conserto sair do escopo desta fatia, mande `ESCALATION`.

- [ ] **Step 6: Commit**

```bash
git add -A apps/server_core
git commit -m "refactor(fiscal): remove pricingtax e prova por go list que orders nao depende mais de pricing"
```

---

## Task 17: Verificação ao vivo e fechamento

Chip **não sobe servidor**. Esta task pede o drive ao hub.

- [ ] **Step 1: Mande `REQUEST` ao hub**

Peça: subir o dev stack no commit desta fatia, rodar o sync da matriz uma vez, e abrir `/pedidos`.

- [ ] **Step 2: Escreva o que o hub tem que ver**

| ponto | esperado |
|---|---|
| P-1 | Pedido da Bahia de R$ 299,90: ICMS 20,99 · DIFAL 40,49 · PIS/COFINS 22,05 · total 83,53 |
| P-2 | Margem **−2,30 (−0,77%)**, negativa na tela, sem clamp em zero |
| P-3 | Pedido com destino PR: ICMS preenchido, DIFAL e PIS/COFINS com pendência nomeando `ALIQINTDEST`, margem em branco com motivo |
| P-3b | Pedido **intra-MG**: mesma pendência. A matriz é muda para MG→MG (`ALIQINTDEST` NULL, `ALIQUOTA` 0) e o 19 que a nota grava não sai dela — dívida **D-35**. São 15 das 74 notas de e-commerce; se a tela mostrar número aqui, ele foi inventado |
| P-4 | Rodapé do drawer com *"estimativa pela matriz do ERP"* e a vigência |
| P-5 | `/integracoes` mostra o card da matriz com contagem e data de sync |
| P-6 | Nenhum `R$ 0,00` onde a resposta é desconhecida |

- [ ] **Step 3: Mande `CLOSED` ao hub**

Com o SHA final, o resultado de P-1 a P-6, e a lista de dívidas confirmadas: **D-17** (`CODEMP=1`), **D-21** (7 sítios do 4% vivos em `pricing` até o P4), **D-24** (DIFAL não precificável em PR/RS/SC), **D-25**, **D-26**, **D-27**, **D-28** (histórico começa no 1º sync — a `TGFICM` não é versionada e a `TGFHICM` não guarda antes/depois), **D-29** a **D-34** (ver autorrevisão).

Inclua a contagem medida ao vivo: quantas das 38 ordens saíram com imposto completo, quantas com pendência, e o motivo de cada classe. É o número que dimensiona o D-31.

---

## Autorrevisão do plano

**Cobertura da spec:** §4.1 → T1, T2 · §4.2 → T1, T5 · §4.3 → T3, T10 · §4.4 → T10, T15 · §5 → T6, T7, T9 · §6 → T8 · §7 → T13 · §8 → T14, T15 · §9 → T16 · §10 → controles negativos em T1, T3, T5, T6, T8, T9, T11, T12 · §11 → T17.

**M19 respondida — o DDL do `TGFICM` está na Task 4, medido, não inventado.** Três coisas que a resposta mudou no plano:

1. A alíquota da operação é **`ALIQUOTA`**. `ALIQUFDEST` pertence ao bloco de ST (domínio 17,0–22,0) e daria um ICMS ~2,5× maior. Espelhamos `ALIQUFDEST` mesmo assim, sem usá-la: com 87,7% de cobertura e mais perto da alíquota legal, ela é o sinal de divergência da fatia **P5** — e sai de graça agora.
2. **`REDBASE` é não-zero em 2.669 de 4.209 linhas MG.** A fórmula anterior ignorava redução de base. Agora o ICMS usa `V × (1 − redbase)`, e o DIFAL vira **pendência** quando `redbase ≠ 0`, porque o telescópio `ICMS + DIFAL = V × ALIQINTDEST` pressupõe a mesma base dos dois lados e a base do destino (`PERCREDBASEDEST`) não foi apurada.
3. **`CODTRIB` é nullable** (22 linhas MG) e tem `10` e `40` além de `0`/`60`. Nil e os não medidos viram pendência, nunca `0`.

A base de PIS/COFINS passou a descontar o **dinheiro** apurado em vez da soma das alíquotas. As duas formas dão o mesmo número quando `redbase = 0` — o pedido da Bahia continua em 22,05 — e só a primeira continua certa quando não.

**M20 respondida. Reconciliação contra 347 linhas de nota real, 147 notas, 11 UFs, R$ 93.120 de DIFAL.**

| pergunta | veredito |
|---|---|
| qual coluna alimenta o DIFAL | **`ALIQINTDEST`.** 274 acertos contra 18 de `ALIQUFDEST` — e 16 dos 18 são empate em SP, onde as duas valem 18 |
| fórmula do DIFAL | `ROUND(BASE × (ALIQINTDEST − ALIQUOTA)/100, 2)` — 312/347 ao centavo, 347/347 a 5 centavos, agregado 0,00003% |
| base do DIFAL | **cheia dos dois lados**: `BASEDIFAL = BASE` em 347/347 |
| `PERCREDBASEDEST` | **coluna morta** — 0 em 4.601 de 4.601 linhas, todas as origens. **D-31 morre** |
| `REDBASE` no recorte real | pega **2 de 74** pedidos, e a única ocorrência é `REDBASE=100 ∧ CODTRIB=60`, que é marcação de **ST** |
| forma intra-UF | slots invertem: `I/<grupo>` + **`S/1`** — não `S/−1` |

**Uma correção do especialista que ele mesmo emitiu:** o `ALIQINTDEST(BA) = 18` da rodada 11 era **mediana da UF sobre população mista**, não célula. A célula real do grupo 122 vale **20,5**, batendo com o que já tínhamos medido. A tabela de medianas por UF daquela rodada não descreve célula nenhuma e foi descartada.

**Uma recomendação dele que eu recuso, com a tabela dele:** ele sugeriu usar `ALIQUFDEST` como fallback quando `ALIQINTDEST` for NULL, alegando ganho de 2 linhas (1 PE, 1 RS). Mas no **PR**, onde `ALIQINTDEST` também é NULL, `ALIQUFDEST` vale 19 e a nota gravou **18** — o fallback erraria as **11 linhas** do PR. Ganha 2, perde 11. **Célula muda continua pendência.** `ALIQUFDEST` fica espelhada e nunca entra no cálculo.

**⛔ Lacuna nova, aberta, e é a única que importa (D-35):** existe uma fonte de alíquota interna de destino **fora da `TGFICM`**. Prova em dois lugares — a nota 895507 gravou 17,0 **no dia seguinte** às correções, na mesma célula já em 20,5 (não é defasagem de cadastro, é a nota ignorando o cadastro); e as **56 linhas intra-MG** gravaram 19 enquanto a matriz tem `ALIQINTDEST` NULL e `ALIQUOTA` 0. Hipótese mais barata, não testada: herança do pedido 313 para a nota 306.
**Consequência aceita:** o espelho reproduz o **cadastro**, não o documento. Para *simular preço* é o que se quer — o devido, não o errado. Para *conferir nota emitida* não serve, e isso é a fatia P6, que lê `TGFDIN` direto. Escrito no plano para ninguém confundir os dois usos.
**Consequência prática:** pedidos intra-MG saem com DIFAL em pendência. São 15 das 74 notas de e-commerce.

**Achado que barateia o P5/histórico:** `TGFHICM` guarda **before-image de linha inteira**, uma por edição — não é só carimbo de auditoria. 99 linhas, 61 com origem MG, desde 2024-07-31. Dá para reconstruir parte do histórico em vez de partir do zero absoluto no primeiro sync. **Não entra nesta fatia** (a D-28 continua valendo para o que ela não cobre), mas deixa de ser impossível.

**Consistência de tipos:** `domain.MatrixLine` (T3) → `mirror.MatrixCell` (T5) → `domain.ICMSRule` (T7) → `ports.ItemTaxes` (T11) → `OrderFiscal` (T13) → `FiscalSection` (T14). `Componente` só existe em `fiscal/domain`; `orders` usa `TaxComponent` própria e o adapter (T11) faz a tradução. `pctOf` é definida em T6 e usada em T8.

**Dívidas novas criadas por este plano:**

| id | o quê |
|---|---|
| D-29 | origem fixa em MG (`UFORIG = 13`), mesma classe da D-17 |
| D-30 | arredondamento duplicado entre `fiscal` e `pricing`, com teste de alarme em T6 |
| ~~D-31~~ | **MORTA.** `PERCREDBASEDEST` é 0 em 4.601/4.601 e `BASEDIFAL = BASE` em 347/347 — não há redução de base do lado do destino a apurar |
| ~~D-32~~ | **FECHADA.** `ALIQINTDEST` alimenta o cálculo (274 acertos contra 18); `ALIQUFDEST` fica espelhada e nunca é usada, nem como fallback |
| D-33 | `CODTRIB` 10 e 40 (43 linhas MG) com tratamento não apurado |
| D-34 | `TIPCALCDIFAL` (=0 em 100% das linhas vistas), `BASESTUFDEST` e `CODTABSTUFDEST` não espelhados nem apurados |
| **D-35** | ⛔ **existe fonte de alíquota interna de destino fora da `TGFICM`** — a nota 895507 gravou 17,0 com a célula em 20,5, e as 56 linhas intra-MG gravaram 19 com a célula muda. Hipótese não testada: herança do pedido 313 para a nota 306. Enquanto aberta, pedidos intra-MG saem com DIFAL em pendência (15 de 74 notas) |
| D-36 | `REDBASE ≠ 0` com `CODTRIB ≠ 60` nunca foi observado ⇒ DIFAL vira pendência nessa combinação. Custo medido: 2 pedidos em 74, e os 2 são ST |
