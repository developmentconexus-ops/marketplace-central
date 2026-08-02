# P2.b — Imposto ex-ante do produto (escopo A) — Plano de Implementação

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Trocar o imposto fabricado de `/pedidos` (SIMPLES 4% de um perfil com 0 linhas) pelo imposto real lido da matriz fiscal do ERP (`METALPRD.TGFICM` via `TGFPRO.GRUPOICMS`), com selo de desconhecido por componente, e corrigir o mesmo dado lido errado em `/anuncios`.

**Architecture:** Um único leitor de `TGFICM` no módulo `internal_read` passa a ser a fonte de toda alíquota de ICMS do sistema. `orders` consome esse leitor através de uma porta redesenhada **por item** (o CODTRIB depende do `GRUPOICMS` do produto, não do total do pedido). `listings` troca o teto worst-case fabricado pelo pior caso **real** calculado sobre a mesma matriz. O cálculo fiscal em si é função pura no domínio, sem I/O, testável por tabela. Nada de alíquota é digitado por nós: o ERP é dono, nós lemos.

**Tech Stack:** Go 1.x (`apps/server_core`), Oracle via `godror` (read-only), Postgres/pgx, React+TypeScript (`apps/web`), OpenAPI contract-first + `packages/sdk-runtime` escrito à mão.

---

## 0. Contexto obrigatório antes de começar

Leia, nesta ordem:

1. `docs/superpowers/specs/2026-08-02-p2b-imposto-ex-ante-design.md` — o desenho aprovado. **Com as três emendas da §0.1 abaixo.**
2. `Documents/MNOS/docs/product/sankhya-custo-formula.md` — por que `CUSSEMICM` **já** é líquido de ST/IPI/crédito de PIS-COFINS, e portanto o que **não** subtrair de novo.
3. `docs/HARNESS-PROFILE.md` e `AGENTS.md` — regras de costura, ADR-04, ADR-17.

### 0.1 Emendas à spec, decididas pelo operador em 2026-08-02 (posteriores ao commit `f6ddd38`)

| # | spec dizia | agora vale |
|---|---|---|
| E-1 | §2 e §11: `GetICMSCeiling` do `/anuncios` **fora** de escopo, vira dívida D-19 | **Dentro do escopo.** O teto worst-case dirige selo, dois filtros e o contador do dashboard, todos alimentados pela coluna errada. Tarefa 13. **D-19 deixa de existir como dívida** — é resolvida aqui |
| E-2 | §6: cinco detectores de conflito | **Três.** Os detectores "previsto ≠ realizado" e "alíquota mudou sem log" comparam a matriz contra a **nota emitida**, que é o subsistema B (imposto ex-post). Não são calculáveis nesta fatia. Movidos para B, registrados como **D-23** |
| E-3 | §8: "componentes de ICMS, DIFAL, FCP, PIS/COFINS com selo", sem dizer o que acontece com o campo `imposto` já publicado | `imposto` = **ICMS + FCP + PIS/COFINS** (exclui DIFAL); `difal` permanece campo próprio; os quatro componentes são detalhe novo. Isto é forçado pela aritmética que já existe em `BuildProfitability` (`margem = total − comissão − frete − imposto − difal − custo`) — qualquer outra leitura duplicaria valor na tela |
| **E-4** | §1 e §9: a `pricing_difal_rates` "erra +35%" porque daria 40,49 contra a nota de 29,99 | **A refutação era inválida e está retirada.** 40,49 é o resultado da **matriz atual** (`299,90 × 20,5% − 299,90 × 7%`), que é também a alíquota legal da BA desde fev/2024. Quem estava defasada era a **nota** 895507 (usou 17,0); a nota-irmã 895436 do dia anterior usou 20,5. **A fonte é a MATRIZ `TGFICM` vigente** — nunca o documento emitido (não existe ex-ante), nunca uma tabela legal mantida por nós. A `pricing_difal_rates` morre mesmo assim, por ser cópia à mão de dado que o ERP já tem. Detalhe na spec, "EMENDA 1" |
| **E-5** | §5: implícito que `ALIQINTDEST` está lá quando `CODTRIB=0` | **Falso em cerca de metade do recorte.** Medido no recorte real (8 grupos × 10 UFs): **21 células calculáveis** (BA/ES/MA/SP/RJ/PE) · **20 mudas** (`CODTRIB=0` sem `ALIQINTDEST`, concentradas em **PR, RS, SC**) · **6 sem linha** (grupo 311). Decisão do operador: PR/RS/SC saem **desconhecido com motivo legível** e o par (produto, UF) vira **pendência apontando a lacuna da matriz**. Zero dado fiscal digitado por nós. Tarefa 14 detecta a célula muda; Tarefa 6 sela |

### 0.2 Medições que fundamentam este plano (todas feitas em 2026-08-02, não re-medir)

| fato | valor medido | onde |
|---|---|---|
| itens de pedido com `internal_product_id` | **33 de 38** (5 são `link_quality='conflict'` com CODPROD nulo) | `orders_marketplace_order_items` |
| itens por pedido | **1 em 100% dos 38 pedidos** | idem |
| UFs de destino reais | SP 15 · RJ 11 · MG 3 · BA 2 · PR 2 · ES 2 · RS 1 · PE 1 · MA 1 | `order_shipments.dest_state`, **sigla de 2 letras** |
| `TGFICM.UFORIG`/`UFDEST` | **NUMBER** (MG=13, BA=5, RJ=19, SP=25) | `internal_read/domain/icms_ceiling.go:4` + medição do especialista |
| conversão sigla↔número no repo | **não existe** | grep repo-wide |
| `pricing_calc_profiles` | **0 linhas** | Postgres |
| tabela de UF do ERP | `TSIUFS`, **30 linhas**, `'EX'` duplicado, sentinela `CODUF=0/'SF'`; 0 órfãos contra `TGFICM` | especialista Sankhya |
| `GRUPOICMS` nulo, cadastro inteiro | 21.673 de 40.439 (**54%**) — cadastro morto | idem |
| `GRUPOICMS` nulo, **vendável** (`USOPROD='R' AND ATIVO='S'`, n=10.041) | **65 (0,65%)**; 99 com zero-ou-nulo | idem |
| `GRUPOICMS` nulo, **caminho e-commerce** (TOP 306/313, 12 meses, 16 produtos) | **0** — cobertura 100% onde precificamos | idem |
| vendáveis sem linha nenhuma na matriz | **28 de 9.976 (0,28%)**, em 7 grupos: 303, 269, 519, 537, 6, 153, 544 | idem |

> Os 0,28% sem linha vão a desconhecido explícito. Não vale engenharia — vale selo.

**⚠ A medição de cobertura da spec §5 não vale — a correta está aqui.** A spec mediu a linha **genérica** (`TIPRESTRICAO='N', CODRESTRICAO=0`); este plano consulta a linha **do grupo** (`TIPRESTRICAO2='I', CODRESTRICAO2=GRUPOICMS`), outra população. Re-medido no recorte real (8 grupos que os pedidos tocaram × 10 UFs de destino), contando só células com `CODTRIB=0`:

| célula | n | UFs |
|---|---|---|
| calculável (`ALIQINTDEST` presente) | **21** | BA, ES, MA, SP, RJ, PE |
| **muda** (`CODTRIB=0`, sem `ALIQINTDEST`) | **20** | **PR, RS, SC** |
| sem linha na matriz | 6 | grupo 311 |

Célula muda **não** é bug nosso: é lacuna de cadastro no Sankhya. Sai como desconhecido + pendência (E-5). Custo medido de não precificar PR/RS/SC: **7 notas em 12 meses**, de 74. **R7 está fechado — ver §2.1.**

**Separe "se incide" de "quanto incide".** `CODTRIB` responde **se** (60 = não deve, 0 = deve) com cobertura quase total. `ALIQINTDEST` responde **quanto** em ~metade das UFs. São dois fatos com confiabilidade diferente e o código nunca pode tratá-los como uma leitura só.

**A consequência de "1 item em 100% dos pedidos":** a porta por item **não pode ser provada por dado de produção**. O pedido misto ST/não-ST só existe em fixture. A Tarefa 7 tem um teste obrigatório de pedido com 2 itens de CODTRIB diferentes — sem ele, o redesenho da porta passa vazio e ninguém percebe.

---

## 0.3 Comunicação com o hub — LEIA, isto não é opcional

Você é uma sessão de execução. O hub é a sessão que escreveu este plano.

**Você PODE e DEVE mandar mensagem para o hub.** A sessão anterior não fez isso e o resultado foi trabalho jogado fora. Especificamente:

- **Dúvida que muda o que você vai escrever → pergunte ANTES de escrever.** Não escolha sozinho entre duas leituras de um requisito. Uma pergunta custa uma mensagem; uma tarefa refeita custa a tarefa inteira.
- **Se você encontrar legacy ou máximo local que este plano não previu → PARE e mande `ESCALATION`.** O operador foi explícito: "se encontrar algo que seja legacy ou local maximum, pare e vamos buscar solução melhor". Não conserte por conta própria, não contorne, não deixe TODO.
- **Se uma medição deste plano estiver errada contra o repo real → mande `ESCALATION` com o file:line que a contradiz.** Briefs apodrecem. Já aconteceu três vezes nesta missão. O plano é **alegação sobre o repo**, não verdade.
- **Ao terminar → mande `CLOSED`** com: commits produzidos, o que foi medido verde, e o que ficou de fora.
- **Bloqueio que você não resolve em ~15 min → `BLOCKED`**, não fique tentando.
- **Precisa de dependência nova, servidor, `.env` ou `:8080` → `REQUEST` ao hub.** Chip nunca sobe servidor nem cria `.env`.

Eventos válidos: `CLOSED` · `BLOCKED` · `ESCALATION` · `REQUEST` · `SPLIT-REQUEST` · `COMMITTED` · `ACK`.

---

## 0.4 Regras que este plano não repete em toda tarefa

- **Testes Go:** sempre com `GOCACHE` absoluto inline, e sempre a partir de `apps/server_core`:
  ```bash
  cd apps/server_core && GOCACHE=/c/Users/leandro.theodoro/Documents/marketplace-central/apps/server_core/.gocache go test ./internal/modules/orders/... -run TestNome -v
  ```
- **Nunca** `git reset`/`revert`/`stash`/`clean`. Nunca push sem permissão explícita do operador.
- **Nunca** ler, imprimir ou commitar conteúdo de `.env*`. Nunca despejar o ambiente inteiro de um processo ou container.
- **Oracle é read-only.** Nenhum `INSERT`/`UPDATE`/`DDL` em `METALPRD`, em nenhuma circunstância.
- **Toda query com tenant tem predicado `tenant_id`.**
- **ADR-17:** desconhecido nunca vira `0`, nunca vira default plausível. Esta fatia inteira existe porque essa regra foi violada.
- **Mudança de API = OpenAPI + `sdk-runtime` no mesmo commit.**
- Commits frequentes, um por tarefa no mínimo.

---

## 1. File Structure

### Criar

| arquivo | responsabilidade |
|---|---|
| `apps/server_core/internal/modules/internal_read/domain/icms_rule.go` | `ICMSRule` — uma linha da matriz fiscal do ERP: CODTRIB + alíquotas + selo de completude. Sem I/O. |
| `apps/server_core/internal/modules/internal_read/domain/uf_code.go` | `UFCode` — de-para sigla ↔ código numérico do ERP. Sem tabela digitada: só o tipo e o erro de não-encontrado. |
| `apps/server_core/internal/modules/internal_read/adapters/oracle/icms_rule.go` | Query de `TGFICM` por (GRUPOICMS, UFORIG, UFDEST). |
| `apps/server_core/internal/modules/internal_read/adapters/oracle/tax_group.go` | Query de `TGFPRO.GRUPOICMS` por CODPROD. |
| `apps/server_core/internal/modules/internal_read/adapters/oracle/uf_code.go` | Query do de-para de UF. **BLOQUEADA — ver Tarefa 1.** |
| `apps/server_core/internal/modules/internal_read/adapters/oracle/icms_revision.go` | `MAX(TGFHICM.DHALTER)` — carimbo de invalidação de cache. |
| `apps/server_core/internal/modules/orders/domain/tax_calc.go` | **O cálculo fiscal.** Função pura: (valor, ICMSRule) → componentes selados. Zero I/O, zero SQL, zero contexto. |
| `apps/server_core/internal/modules/orders/domain/tax_calc_test.go` | Tabela de casos + controles negativos nomeados. |
| `apps/server_core/internal/modules/orders/adapters/internalread/tax_reader.go` | Adapter orders → internal_read. Substitui o consumo do perfil de `pricing`. |

### Modificar

| arquivo | ação |
|---|---|
| `internal_read/ports/batch_reader.go` | +3 métodos (regras, grupos, revisão) |
| `internal_read/adapters/cache/cache.go` | decorar os métodos novos; invalidação por revisão, não por TTL |
| `internal_read/adapters/oracle/icms_ceiling.go` | **substituído por dentro** — passa a calcular o pior caso REAL sobre `ALIQINTDEST` por grupo |
| `listings/ports/enrichment.go` | `GetICMSCeilingByOrigin` ganha o grupo do produto na assinatura |
| `listings/adapters/internalread/cost_reader.go` | adaptar à assinatura nova |
| `listings/application/read_service.go` | `icmsWorstCaseByUF` + `belowMargin` passam a usar o teto real |
| `orders/ports/tax_reader.go` | **porta redesenhada por item** |
| `orders/adapters/pricingtax/reader.go` | **reescrito**; deixa de consumir `ProfileSource` |
| `orders/application/enrich_service.go` | `resolveTaxes` → `resolveItemTaxes`, espelhando `resolveItemCosts` |
| `orders/domain/order_decomposition.go` | +4 componentes com selo |
| `orders/transport/http_handler.go` | +4 campos no DTO |
| `contracts/api/marketplace-central.openapi.yaml` | `OrderDecomposicao` +4 campos |
| `packages/sdk-runtime/src/index.ts` | espelhar |
| `apps/web/src/pages/pedidos/PedidoDrawer.tsx` | renderizar componentes + motivo legível |
| `internal/composition/root.go` | trocar a fiação; `orderspricingtax` sai |

### Não tocar

- `pricing/adapters/postgres/calc_repository.go` — costura do módulo `pricing` (ADR-04).
- `pricing_difal_rates` e a migração 0057 — a fatia remove o **consumo**, não a tabela.

---

## 2. Ordem das tarefas e por que ela é essa

```
T1 (UF)  ─┐
T2 (domínio ICMSRule) ─┬─ T3 (grupos) ─┬─ T5 (cache) ─┐
                       └─ T4 (matriz) ─┘              │
T6 (cálculo puro, não depende de I/O) ────────────────┤
                                                      ├─ T7 (porta) ─ T8 (adapter) ─ T9 (enrich)
                                                      │                                  │
                                                      │                    T10 (decomposição) ─ T11 (contrato+FE)
                                                      │                                  │
                                                      └────────────── T12 (fiação) ──────┘
                                                                          │
                                                       T13 (listings) ─ T14 (detector) ─ T15 (QA)
```

T6 é puro e não depende de nada — pode ser feito primeiro se você quiser ver valor cedo.

## 2.1 R7 — **FECHADO** em 2026-08-02. O DIFAL é calculável em ~metade das UFs, e isso está decidido

**Não comece T6 nem T8 sem ler isto.** R7 estava aberto quando este plano foi escrito; a medição chegou e o operador decidiu. O que segue é decisão, não risco.

### O que se mediu

Na linha **do grupo** (que é a que este plano consulta), no recorte real de 8 grupos × 10 UFs de destino, contando só células com `CODTRIB=0`:

**21 calculáveis** (BA, ES, MA, SP, RJ, PE) · **20 mudas** (`CODTRIB=0` sem `ALIQINTDEST`, em **PR, RS, SC**) · **6 sem linha** (grupo 311).

O ERP **cobra** DIFAL do PR — 11 notas com `ALIQINTDEST` 18,0 — tirando o número de um lugar que não está em `TGFICM`, `TSIUFS` nem `TGFPAR`, e que **não foi localizado**. Não procure; não é seu trabalho nesta fatia.

Duas medições anteriores continuam valendo e são o motivo de dois testes existirem:

1. **`ALIQUFDEST` ≠ `ALIQINTDEST` onde as duas existem:** MA 23×22, MS 18×17, MT 12×17, PR 17×19, ES 17×19,5, RJ 18×20, TO 17×20. Só SP concorda. Trocar uma pela outra erra o DIFAL em quase toda UF — é por isso que a T4 tem teste de query.
2. **`TSIUFS.AD_ALIQINT` é cadastro morto** (congelado em 17 enquanto GO foi a 19 em 2024). Não leia essa coluna em lugar nenhum.

### A decisão do operador

O especialista recomendou manter tabela de alíquota legal própria, fora do Sankhya. **Recusado**, e o motivo é o mesmo da spec §2: é a `pricing_difal_rates` voltando com boa intenção. Vale:

> **PR, RS, SC → DIFAL desconhecido, selado, com motivo legível na tela, e o par (produto, UF) vira pendência dizendo que falta `ALIQINTDEST` na matriz do ERP.**
> Zero dado fiscal digitado por nós. O conserto acontece no Sankhya, que é o dono.

Custo aceito e medido: 7 notas em 12 meses (de 74) ficam sem preço confiável até alguém preencher a matriz. A Tarefa 14 é o mecanismo pelo qual essa lacuna aparece para quem pode consertá-la.

### A separação que o código tem que refletir

> `TGFICM` decide **se incide** (`CODTRIB` 60 vs 0) e nisso é confiável — só o grupo 311 tem buraco.
> **Quanto incide** no DIFAL, ela não responde em PR, RS e SC.

Hoje as duas perguntas estão coladas na mesma leitura. `ICMSRule` (T2) já as separa por construção — `ChainClosed()` responde a primeira, `DifalKnown()` a segunda. A separação é obrigatória: colapsar as duas transforma "não sei quanto" em "não deve", que é exatamente o defeito que esta fatia existe para matar.

---

## Task 1: De-para de UF (sigla ↔ código do ERP) — **DESBLOQUEADA** (medida 2026-08-02)

**Files:**
- Create: `apps/server_core/internal/modules/internal_read/adapters/oracle/uf_code.go`
- Test: `apps/server_core/internal/modules/internal_read/adapters/oracle/uf_code_test.go`

**O problema:** `TGFICM.UFDEST` é numérico (BA=5, MG=13, RJ=19, SP=25). `order_shipments.dest_state` é sigla (`"BA"`). Nada no repositório converte. Decisão do operador: **JOIN no SQL, o ERP é dono do de-para** — não digitamos a tabela de UFs no nosso código, porque dado do ERP digitado por nós é a classe de erro que matou `pricing_difal_rates`.

**Medição do especialista Sankhya, 2026-08-02 — leia antes de escrever a query:**

Tabela: `METALPRD.TSIUFS`. Colunas: `CODUF` NUMBER NOT NULL · `UF` VARCHAR2(2) NOT NULL · `DESCRICAO` VARCHAR2(40) · `CODPAIS` NUMBER. **Não existe coluna de "ativo".**

Integridade referencial contra `TGFICM`: 27 `UFDEST` distintos e 17 `UFORIG` distintos, **zero sem match**. O JOIN é seguro sem defesa.

**⚠ A tabela tem 30 linhas, não 27, e a sigla NÃO é única:**

| CODUF | UF | o que é |
|---|---|---|
| 0 | `SF` | sentinela `<sem UF>`, `CODPAIS=0` |
| 1–27 | — | as UFs brasileiras (não existe `CODUF=28`) |
| 29 | `EX` | CHINA, `CODPAIS=86` |
| 30 | `EX` | COLON, `CODPAIS=507` |

30 linhas, **29 siglas distintas** — `'EX'` aparece duas vezes. Zero siglas nulas.

**A armadilha:** `WHERE UF = 'EX'` devolve **2 linhas**. Um lookup sigla→código sem filtro faz fanout em silêncio — não estoura, escolhe. Filtre por `CODPAIS = 55`; aí são 27 linhas, 27 siglas, bijeção.

- [ ] **Step 1: Escreva o teste que falha**

```go
func TestUFCodeResolverTranslatesMeasuredSiglas(t *testing.T) {
	r := newUFResolverWithRows(t, tsiufsFixture)
	for sigla, want := range map[string]int64{"BA": 5, "MG": 13, "RJ": 19, "SP": 25} {
		got, err := r.CodeForSigla(context.Background(), sigla)
		if err != nil {
			t.Fatalf("%s: erro inesperado: %v", sigla, err)
		}
		if got != want {
			t.Fatalf("%s deveria ser %d, veio %d", sigla, want, got)
		}
	}
}

func TestUFCodeResolverRejectsUnknownSiglaInsteadOfDefaulting(t *testing.T) {
	r := newUFResolverWithRows(t, tsiufsFixture)
	if _, err := r.CodeForSigla(context.Background(), "XX"); !errors.Is(err, ErrUFUnknown) {
		t.Fatalf("sigla desconhecida tem que dar ErrUFUnknown, nunca um codigo default: %v", err)
	}
}

// TSIUFS tem 'EX' DUAS vezes (29=CHINA, 30=COLON) e a sentinela CODUF=0/'SF'.
// Sem o filtro CODPAIS=55 o lookup por sigla escolhe uma das duas em silencio.
func TestUFCodeResolverFiltersToBrazilSoSiglaIsUnique(t *testing.T) {
	q := buildUFCodeQuery()
	if !strings.Contains(q, "CODPAIS") {
		t.Fatalf("sem filtro de pais, 'EX' casa 2 linhas e o lookup faz fanout:\n%s", q)
	}
	r := newUFResolverWithRows(t, tsiufsFixtureIncludingEXAndSF)
	if _, err := r.CodeForSigla(context.Background(), "EX"); !errors.Is(err, ErrUFUnknown) {
		t.Fatalf("'EX' e pais estrangeiro, nao UF brasileira — nao pode resolver")
	}
	if _, err := r.CodeForSigla(context.Background(), "SF"); !errors.Is(err, ErrUFUnknown) {
		t.Fatalf("'SF' e sentinela <sem UF>, nao pode resolver para 0")
	}
}
```

> **Controle negativo obrigatório:** remova o `WHERE CODPAIS = 55` da query e rode. `TestUFCodeResolverFiltersToBrazilSoSiglaIsUnique` **tem que** ficar vermelho nas três asserções. Se ficar verde sem o filtro, o fixture não tem as linhas `EX`/`SF` e o teste não está medindo nada — conserte o fixture antes de seguir.

- [ ] **Step 2: Rode e confirme o vermelho**

```bash
cd apps/server_core && GOCACHE=/c/Users/leandro.theodoro/Documents/marketplace-central/apps/server_core/.gocache go test ./internal/modules/internal_read/adapters/oracle/ -run TestUFCodeResolver -v
```
Esperado: `undefined: ErrUFUnknown`.

- [ ] **Step 3: Implemente**

```go
package oracle

// ErrUFUnknown e a sigla que nao traduz para codigo do ERP. Uma UF que nao
// sabemos traduzir e imposto DESCONHECIDO, nunca imposto zero e nunca um
// codigo default (ADR-17).
var ErrUFUnknown = errors.New("internal_read: sigla de UF desconhecida")

// buildUFCodeQuery le o de-para do proprio ERP. Nao existe map literal de UFs
// neste repositorio: dado do ERP digitado por nos e a classe de erro que matou
// pricing_difal_rates.
//
// CODPAIS = 55 nao e detalhe: TSIUFS tem 30 linhas, e a sigla 'EX' aparece
// DUAS vezes (CODUF 29 = CHINA, 30 = COLON) alem da sentinela CODUF 0 / 'SF'
// (<sem UF>). Sem o filtro, um lookup por sigla casa 2 linhas e escolhe uma
// sem avisar. Com ele: 27 linhas, 27 siglas, bijecao.
func buildUFCodeQuery() string {
	return `SELECT UF, CODUF FROM METALPRD.TSIUFS WHERE CODPAIS = 55`
}
```

O resolver carrega o mapa inteiro numa leitura (27 linhas) e serve de memória, com a mesma invalidação por revisão da Tarefa 5 — o de-para muda com frequência ainda menor que a matriz.

- [ ] **Step 4: verde** · **Step 5: commit**

```bash
git add apps/server_core/internal/modules/internal_read/adapters/oracle/uf_code.go apps/server_core/internal/modules/internal_read/adapters/oracle/uf_code_test.go
git commit -m "feat(internal_read): de-para de UF lido de TSIUFS, filtrado a CODPAIS=55"
```

---

## Task 2: `ICMSRule` — a linha da matriz fiscal no domínio

**Files:**
- Create: `apps/server_core/internal/modules/internal_read/domain/icms_rule.go`
- Test: `apps/server_core/internal/modules/internal_read/domain/icms_rule_test.go`

- [ ] **Step 1: Escreva o teste que falha**

```go
package domain

import "testing"

func TestICMSRuleCompletenessIsPerColumn(t *testing.T) {
	aliq := 7.0
	// CODTRIB 0 com ALIQINTDEST vazio: ICMS é calculável, DIFAL não é.
	rule := ICMSRule{CodTrib: 0, Aliquota: &aliq, AliqIntDest: nil}
	if !rule.ICMSKnown() {
		t.Fatalf("ICMS deveria ser conhecido: ALIQUOTA presente")
	}
	if rule.DifalKnown() {
		t.Fatalf("DIFAL nao pode ser conhecido com ALIQINTDEST nulo")
	}
	if got := rule.DifalUnknownReason(); got != "ALIQINTDEST vazio na matriz do ERP" {
		t.Fatalf("motivo ilegivel: %q", got)
	}
}

func TestICMSRuleCodTrib60ClosesChain(t *testing.T) {
	rule := ICMSRule{CodTrib: 60}
	if !rule.ChainClosed() {
		t.Fatalf("CODTRIB 60 encerra a cadeia por ST")
	}
	if !rule.ICMSKnown() || !rule.DifalKnown() {
		t.Fatalf("cadeia encerrada e CONHECIDA e vale zero, nao desconhecida")
	}
}

func TestICMSRuleCodTrib10IsUnknownNotZero(t *testing.T) {
	rule := ICMSRule{CodTrib: 10}
	if rule.ICMSKnown() || rule.DifalKnown() {
		t.Fatalf("CODTRIB 10 nao foi caracterizado para PF; e desconhecido, nao zero")
	}
}
```

- [ ] **Step 2: Rode e confirme o vermelho**

```bash
cd apps/server_core && GOCACHE=/c/Users/leandro.theodoro/Documents/marketplace-central/apps/server_core/.gocache go test ./internal/modules/internal_read/domain/ -run TestICMSRule -v
```
Esperado: `undefined: ICMSRule`.

- [ ] **Step 3: Implemente**

```go
package domain

// ICMSRule e uma linha da matriz fiscal do ERP (METALPRD.TGFICM) resolvida
// para um par (GRUPOICMS, UF de destino). Nos LEMOS a decisao fiscal do
// Sankhya; nao a derivamos. Toda alicota aqui veio do ERP — nenhuma constante
// fiscal e digitada neste repositorio.
//
// Completude e POR COLUNA, nunca por linha: uma regra pode ter ALIQUOTA e nao
// ter ALIQINTDEST, e nesse caso o ICMS e calculavel e o DIFAL nao. Colapsar
// isso em "regra incompleta" apagaria um numero que temos (ADR-17 ao contrario:
// esconder conhecido e tao errado quanto inventar desconhecido).
type ICMSRule struct {
	// GroupID e o TGFPRO.GRUPOICMS que enderecou esta linha.
	GroupID int64
	// DestinationUF e o codigo NUMERICO do ERP (BA=5, MG=13, RJ=19, SP=25).
	DestinationUF UF
	// CodTrib e TGFICM.CODTRIB: 0 = tributado, 60 = ST (cadeia encerrada),
	// 10 = nao caracterizado para nossa venda.
	CodTrib int64
	// Aliquota e a interestadual de saida (7% N/NE/CO+ES, 12% S/SE).
	Aliquota *float64
	// AliqIntDest e a interna do destino — a coluna que as notas reais
	// reconciliam. NAO e ALIQUFDEST: as duas divergem (PR 19 x 17, MA 22 x 23)
	// e ALIQUFDEST e a que o leitor de teto usava errado.
	AliqIntDest *float64
	// FCPPercent e TGFICM.PERCICMSFCP. Medido: so o RJ tem, 2%.
	FCPPercent *float64
	Source     SourceMetadata
}

// ChainClosed diz que a substituicao tributaria ja encerrou a cadeia: ICMS,
// DIFAL e FCP sao todos zero, e esse zero e CONHECIDO, nao ausente.
func (r ICMSRule) ChainClosed() bool { return r.CodTrib == 60 }

// characterized separa "CODTRIB que sabemos ler" de "CODTRIB que nao
// caracterizamos para venda a pessoa fisica" (10 e quaisquer outros).
func (r ICMSRule) characterized() bool { return r.CodTrib == 0 || r.CodTrib == 60 }

func (r ICMSRule) ICMSKnown() bool {
	if !r.characterized() {
		return false
	}
	if r.ChainClosed() {
		return true
	}
	return r.Aliquota != nil
}

func (r ICMSRule) DifalKnown() bool {
	if !r.characterized() {
		return false
	}
	if r.ChainClosed() {
		return true
	}
	return r.Aliquota != nil && r.AliqIntDest != nil
}

// FCP ausente na matriz e zero CONHECIDO, nao desconhecido: a medicao mostrou
// que so o RJ tem FCP e que nenhuma nota carregou FCP sem respaldo na matriz.
// Nao existe segunda fonte de FCP para "faltar".
func (r ICMSRule) FCPKnown() bool { return r.characterized() }

func (r ICMSRule) ICMSUnknownReason() string {
	if !r.characterized() {
		return "CODTRIB nao caracterizado para comprador pessoa fisica"
	}
	if r.Aliquota == nil {
		return "ALIQUOTA vazia na matriz do ERP"
	}
	return ""
}

func (r ICMSRule) DifalUnknownReason() string {
	if !r.characterized() {
		return "CODTRIB nao caracterizado para comprador pessoa fisica"
	}
	if r.Aliquota == nil {
		return "ALIQUOTA vazia na matriz do ERP"
	}
	if r.AliqIntDest == nil {
		return "ALIQINTDEST vazio na matriz do ERP"
	}
	return ""
}
```

- [ ] **Step 4: Rode e confirme o verde**

```bash
cd apps/server_core && GOCACHE=/c/Users/leandro.theodoro/Documents/marketplace-central/apps/server_core/.gocache go test ./internal/modules/internal_read/domain/ -run TestICMSRule -v
```
Esperado: PASS, 3 testes.

- [ ] **Step 5: Commit**

```bash
git add apps/server_core/internal/modules/internal_read/domain/icms_rule.go apps/server_core/internal/modules/internal_read/domain/icms_rule_test.go
git commit -m "feat(internal_read): ICMSRule com completude por coluna, nao por linha"
```

---

## Task 3: Ler `TGFPRO.GRUPOICMS` por CODPROD

**Files:**
- Create: `apps/server_core/internal/modules/internal_read/adapters/oracle/tax_group.go`
- Test: `apps/server_core/internal/modules/internal_read/adapters/oracle/tax_group_test.go`
- Modify: `apps/server_core/internal/modules/internal_read/ports/batch_reader.go:11-15`

**Por que método novo e não estender `buildFindProductsQuery`:** `reader.go:414-433` é a costura de leitura de produto e não lê `GRUPOICMS`. Tocá-la arrastaria os consumidores de produto para um problema fiscal que não é deles. Método próprio, pergunta própria — mesmo critério que já separou `GetSalesHistory` (`batch_reader.go:155`).

- [ ] **Step 1: Escreva o teste que falha**

```go
package oracle

import (
	"context"
	"testing"
)

func TestGetTaxGroupsByIDsReturnsNilForProductWithoutGroup(t *testing.T) {
	db := newFakeDB(t, [][]driverValue{
		{int64(39563), int64(122)},
		{int64(15956), nil}, // produto sem GRUPOICMS
	})
	r := NewBatchReader(db, testSemaphore())

	got, err := r.GetTaxGroupsByIDs(context.Background(), []int64{39563, 15956})
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if got[39563] == nil || *got[39563] != 122 {
		t.Fatalf("grupo do 39563 deveria ser 122, veio %v", got[39563])
	}
	if got[15956] != nil {
		t.Fatalf("produto sem GRUPOICMS tem que vir NIL, nunca 0 — 0 e um grupo valido")
	}
}
```

> **Atenção ao controle negativo:** o ponto deste teste é `0` **ser um GRUPOICMS válido**. Se a implementação usar `int64` em vez de `*int64`, "produto sem grupo" e "produto do grupo 0" ficam byte-idênticos e o defeito é invisível. Antes de implementar, troque `*int64` por `int64` de propósito e veja este teste ficar verde — é a prova de que ele mede o que diz medir.

- [ ] **Step 2: Rode e confirme o vermelho**

```bash
cd apps/server_core && GOCACHE=/c/Users/leandro.theodoro/Documents/marketplace-central/apps/server_core/.gocache go test ./internal/modules/internal_read/adapters/oracle/ -run TestGetTaxGroups -v
```
Esperado: `r.GetTaxGroupsByIDs undefined`.

- [ ] **Step 3: Implemente**

```go
package oracle

import (
	"context"
	"database/sql"

	"marketplace-central/apps/server_core/internal/modules/internal_read/adapters/oraclebatch"
)

// GetTaxGroupsByIDs resolve TGFPRO.GRUPOICMS por CODPROD. O grupo e a chave
// que torna TGFICM consultavel ANTES da venda: sem ele a matriz so seria
// legivel depois da nota emitida.
//
// Valor NIL significa produto sem GRUPOICMS cadastrado — o que e diferente de
// grupo 0, que existe e e valido. Confundir os dois faria produto sem cadastro
// herdar a regra do grupo 0 em silencio.
func (r *BatchReader) GetTaxGroupsByIDs(ctx context.Context, ids []int64) (map[int64]*int64, error) {
	ids = uniquePositiveIDs(ids)
	result := make(map[int64]*int64, len(ids))
	for _, id := range ids {
		result[id] = nil
	}
	if len(ids) == 0 {
		return result, nil
	}
	if err := r.ensureBatchAvailable(ctx); err != nil {
		return nil, err
	}
	release, err := r.semaphore.Acquire(ctx)
	if err != nil {
		return nil, wrapOracleError("acquire Oracle batch semaphore", err)
	}
	defer release()

	for _, chunk := range oraclebatch.Chunks(ids, batchChunkSize) {
		query, args := buildTaxGroupsQuery(chunk)
		rows, queryErr := r.db.QueryContext(ctx, query, args...)
		if queryErr != nil {
			return nil, wrapOracleError("read tax groups batch", queryErr)
		}
		for rows.Next() {
			var productID, group sql.NullInt64
			if scanErr := rows.Scan(&productID, &group); scanErr != nil {
				rows.Close()
				return nil, wrapOracleError("scan tax groups batch", scanErr)
			}
			if !productID.Valid || !group.Valid {
				continue
			}
			value := group.Int64
			result[productID.Int64] = &value
		}
		if rowsErr := rows.Err(); rowsErr != nil {
			rows.Close()
			return nil, wrapOracleError("iterate tax groups batch", rowsErr)
		}
		rows.Close()
	}
	return result, nil
}
```

`buildTaxGroupsQuery` segue o mesmo formato de bind por chunk já usado em `buildTaxFactsQuery` (mesmo arquivo, mesma convenção de `sql.Named`):

```go
func buildTaxGroupsQuery(ids []int64) (string, []any) {
	placeholders, args := namedIDBinds(ids)
	return `SELECT CODPROD, GRUPOICMS FROM METALPRD.TGFPRO WHERE CODPROD IN (` + placeholders + `)`, args
}
```

> Se `namedIDBinds` não existir com esse nome, use o helper que `buildTaxFactsQuery` já usa — **não crie um segundo**. Se não houver helper e a construção de binds estiver duplicada inline, isso é máximo local: mande `ESCALATION` ao hub antes de duplicar pela terceira vez.

- [ ] **Step 4: Adicione ao port**

Em `internal_read/ports/batch_reader.go`, dentro da interface existente:

```go
	GetTaxGroupsByIDs(ctx context.Context, ids []int64) (map[int64]*int64, error)
```

- [ ] **Step 5: Rode e confirme o verde**

```bash
cd apps/server_core && GOCACHE=/c/Users/leandro.theodoro/Documents/marketplace-central/apps/server_core/.gocache go test ./internal/modules/internal_read/... -v
```
Esperado: PASS. **Se algum fake de teste em outro pacote quebrar por não implementar o método novo, isso é bom** — é o compilador achando os implementadores da porta. Corrija cada um adicionando o método; não use embedding de interface para calar o compilador (memória `chip-m02-mirror-port-merged`: decorator que apaga porta opcional).

- [ ] **Step 6: Commit**

```bash
git add apps/server_core/internal/modules/internal_read/
git commit -m "feat(internal_read): GRUPOICMS por CODPROD; nil != grupo 0"
```

---

## Task 4: Ler a matriz `TGFICM` por (grupo, UF)

**Files:**
- Create: `apps/server_core/internal/modules/internal_read/adapters/oracle/icms_rule.go`
- Test: `apps/server_core/internal/modules/internal_read/adapters/oracle/icms_rule_test.go`
- Modify: `apps/server_core/internal/modules/internal_read/ports/batch_reader.go`

**A chave de consulta** (spec §3, medida pelo especialista):

```
interestadual:  UFORIG = 13 (MG)  ·  UFDEST = <destino>
                TIPRESTRICAO2 = 'I'  ·  CODRESTRICAO2 = GRUPOICMS
intra-MG:       UFORIG = 13 · UFDEST = 13
                TIPRESTRICAO = 'I'  ·  CODRESTRICAO = GRUPOICMS
```

As posições invertem entre intra e interestadual. `TIPRESTRICAO` é **tag polimórfica de tipo** (`O`=TOP, `I`=GRUPOICMS, `G`=grupo de produto, `P`=produto, `N`=genérica), não uma dimensão. Filtrar por `'I'` é o que endereça o grupo.

- [ ] **Step 1: Escreva o teste que falha**

```go
func TestGetICMSRulesByGroupsReadsAliqIntDestNotAliqUfDest(t *testing.T) {
	// A query tem que citar ALIQINTDEST. ALIQUFDEST diverge dela (PR 19 x 17)
	// e e a coluna que o leitor de teto usava errado.
	q, _ := buildICMSRulesQuery(13, []int64{122})
	if !strings.Contains(q, "ALIQINTDEST") {
		t.Fatalf("query tem que ler ALIQINTDEST:\n%s", q)
	}
	if strings.Contains(q, "ALIQUFDEST") {
		t.Fatalf("ALIQUFDEST e a coluna errada — as notas reais reconciliam com ALIQINTDEST:\n%s", q)
	}
	if strings.Contains(q, "MAX(") {
		t.Fatalf("MAX() colapsa linhas de especificidade diferente — foi o defeito do teto:\n%s", q)
	}
}

func TestGetICMSRulesByGroupsKeepsNullAliqIntDest(t *testing.T) {
	db := newFakeDB(t, [][]driverValue{
		{int64(122), int64(15), int64(0), 12.0, nil, nil},   // PR: ALIQINTDEST nulo
		{int64(122), int64(19), int64(0), 12.0, 20.0, 2.0},  // RJ: completo, com FCP
		{int64(311), int64(5), int64(60), nil, nil, nil},    // BA: ST
	})
	r := NewBatchReader(db, testSemaphore())

	got, err := r.GetICMSRulesByGroups(context.Background(), 13, []int64{122, 311})
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	pr := got[122][15]
	if pr == nil || pr.Aliquota == nil || *pr.Aliquota != 12.0 {
		t.Fatalf("ALIQUOTA do PR deveria sobreviver: %+v", pr)
	}
	if pr.AliqIntDest != nil {
		t.Fatalf("ALIQINTDEST nulo tem que chegar NIL, nunca 0 — 0 viraria DIFAL negativo")
	}
	rj := got[122][19]
	if rj == nil || rj.FCPPercent == nil || *rj.FCPPercent != 2.0 {
		t.Fatalf("FCP do RJ deveria ser 2%%: %+v", rj)
	}
	ba := got[311][5]
	if ba == nil || ba.CodTrib != 60 {
		t.Fatalf("BA do grupo 311 deveria ser CODTRIB 60: %+v", ba)
	}
}
```

> **Controle negativo obrigatório:** antes de implementar, escreva a query com `ALIQUFDEST` de propósito e veja `TestGetICMSRulesByGroupsReadsAliqIntDestNotAliqUfDest` ficar vermelho. Um teste de query que nunca viu a coluna errada não prova que pegaria a coluna errada.

- [ ] **Step 2: Rode e confirme o vermelho**

```bash
cd apps/server_core && GOCACHE=/c/Users/leandro.theodoro/Documents/marketplace-central/apps/server_core/.gocache go test ./internal/modules/internal_read/adapters/oracle/ -run TestGetICMSRules -v
```

- [ ] **Step 3: Implemente**

```go
// GetICMSRulesByGroups le a matriz fiscal do ERP para os grupos pedidos,
// saindo de originUF. Devolve grupo -> UF destino -> regra.
//
// Uma UF ausente do mapa interno NAO e "sem imposto": e "sem linha na matriz",
// que o chamador tem que tratar como desconhecido (spec §5). O leitor nao
// inventa linha faltante.
func (r *BatchReader) GetICMSRulesByGroups(ctx context.Context, originUF domain.UF, groups []int64) (map[int64]map[domain.UF]*domain.ICMSRule, error)
```

Corpo: mesmo esqueleto de `GetTaxGroupsByIDs` (semáforo, chunks, `rows.Close()` em todo caminho de erro). A query:

```go
func buildICMSRulesQuery(originUF domain.UF, groups []int64) (string, []any) {
	placeholders, args := namedIDBinds(groups)
	args = append(args, sql.Named("uforig", int64(originUF)))
	return `
SELECT CODRESTRICAO2 AS GRUPO, UFDEST, CODTRIB, ALIQUOTA, ALIQINTDEST, PERCICMSFCP
  FROM METALPRD.TGFICM
 WHERE UFORIG = :uforig
   AND UFDEST <> :uforig
   AND TIPRESTRICAO2 = 'I'
   AND CODRESTRICAO2 IN (` + placeholders + `)
UNION ALL
SELECT CODRESTRICAO AS GRUPO, UFDEST, CODTRIB, ALIQUOTA, ALIQINTDEST, PERCICMSFCP
  FROM METALPRD.TGFICM
 WHERE UFORIG = :uforig
   AND UFDEST = :uforig
   AND TIPRESTRICAO = 'I'
   AND CODRESTRICAO IN (` + placeholders + `)`, args
}
```

> **Cuidado com bind repetido:** Oracle com `sql.Named` reusa o bind por nome, então `:uforig` aparecer 4 vezes é correto e leva **um** argumento. Mas `placeholders` aparece duas vezes — se `namedIDBinds` gerar nomes posicionais únicos, você precisa dos args **uma vez só** (mesmos nomes reusados), não duplicados. Rode o teste de query e confira a contagem de args antes de seguir. Se o driver reclamar de bind faltando ou sobrando, é aqui.

Se a UNION ficar difícil de amarrar com o helper de binds, **é aceitável fazer duas queries** (uma interestadual, uma intra-MG) e mesclar em Go. Documente a escolha no comentário. Não vale duplicar o helper de binds para forçar a UNION.

- [ ] **Step 4: Adicione ao port + Step 5: verde + Step 6: commit**

```go
	GetICMSRulesByGroups(ctx context.Context, originUF domain.UF, groups []int64) (map[int64]map[domain.UF]*domain.ICMSRule, error)
```

```bash
cd apps/server_core && GOCACHE=/c/Users/leandro.theodoro/Documents/marketplace-central/apps/server_core/.gocache go test ./internal/modules/internal_read/... -v
git add apps/server_core/internal/modules/internal_read/
git commit -m "feat(internal_read): matriz TGFICM por GRUPOICMS, lendo ALIQINTDEST sem MAX()"
```

---

## Task 5: Cache invalidado pelo log do ERP, não por TTL

**Files:**
- Create: `apps/server_core/internal/modules/internal_read/adapters/oracle/icms_revision.go`
- Modify: `apps/server_core/internal/modules/internal_read/adapters/cache/cache.go:410-455`
- Test: `apps/server_core/internal/modules/internal_read/adapters/cache/cache_test.go`

**Por que não TTL** (spec §7): `TGFICM` é sobrescrita in-place, sem versionamento (4.601 linhas com vigência nula em 100%). Um TTL de 24h serve alíquota velha por até 24h **sem saber**. `TGFHICM.DHALTER` muda no instante da alteração. Limite honesto: o log é incompleto (o `19,0` do RJ apareceu sem alteração logada) — é melhor que TTL, não é garantia. **Dívida D-20, já registrada, não re-registrar.**

- [ ] **Step 1: Teste que falha**

```go
func TestICMSCacheInvalidatesWhenRevisionAdvances(t *testing.T) {
	down := &fakeBatchReader{rules: rulesV1, revision: mustTime("2026-07-01T00:00:00Z")}
	c := NewBatchReader(down, newTestCache())

	first, _ := c.GetICMSRulesByGroups(ctx, 13, []int64{122})
	down.rules = rulesV2 // ERP mudou a matriz
	second, _ := c.GetICMSRulesByGroups(ctx, 13, []int64{122})
	if !sameAliquota(first, second) {
		t.Fatalf("sem mudanca no log, o cache tem que servir o valor antigo")
	}

	down.revision = mustTime("2026-08-01T00:00:00Z") // TGFHICM registrou alteracao
	third, _ := c.GetICMSRulesByGroups(ctx, 13, []int64{122})
	if sameAliquota(second, third) {
		t.Fatalf("revisao avancou e o cache continuou servindo aliquota velha")
	}
}
```

> **Controle negativo:** troque a invalidação por um TTL de 24h e veja a terceira asserção ficar vermelha. Se ela passar com TTL, o teste não está medindo a invalidação por log.

- [ ] **Step 2: vermelho** · **Step 3: implemente** · **Step 4: verde** · **Step 5: commit**

A revisão entra na **chave** do cache (`canonicalKey` já existe em `cache.go`), não num campo de expiração: chave nova = miss natural, sem código de despejo. A leitura da revisão é ela própria cacheada com TTL curto (30s) — senão toda leitura de matriz vira duas idas ao Oracle.

```go
func buildICMSRevisionQuery() string {
	return `SELECT MAX(DHALTER) FROM METALPRD.TGFHICM`
}
```

```bash
git commit -m "feat(internal_read): cache de matriz fiscal invalidado por MAX(TGFHICM.DHALTER)"
```

---

## Task 6: O cálculo fiscal — função pura

**Files:**
- Create: `apps/server_core/internal/modules/orders/domain/tax_calc.go`
- Test: `apps/server_core/internal/modules/orders/domain/tax_calc_test.go`

**Esta é a tarefa mais importante do plano.** Tudo antes é encanamento; tudo depois é fiação. Aqui mora o número que o operador vê. Sem I/O, sem `context`, sem SQL — só entrada e saída, para que o teste seja sobre a conta e não sobre o banco.

**A conta** (spec §4):

```
CODTRIB = 60   →  ICMS = 0 · DIFAL = 0 · FCP = 0        (cadeia encerrada por ST)
CODTRIB = 0    →  ICMS  = V × ALIQUOTA
                  DIFAL = V × ALIQINTDEST − ICMS
                  FCP   = V × PERCICMSFCP
CODTRIB = 10   →  DESCONHECIDO
CODTRIB outro  →  DESCONHECIDO

PIS/COFINS     =  V × 8,62%   (9,25% sobre base reduzida ≈ 93,1% do cheio)
IPI            =  0
```

**O que NÃO subtrair de novo:** `CUSSEMICM` já está líquido de ST, IPI e crédito de PIS/COFINS; no ramo com ST o ICMS próprio **fica dentro do custo**. Somar o **débito** de PIS/COFINS da saída não é dupla contagem — o crédito é da entrada e já está no custo, o débito é da receita e hoje não existe em lugar nenhum.

- [ ] **Step 1: Escreva os testes que falham**

```go
package domain

import "testing"

func f(v float64) *float64 { return &v }

// O caso que define a fatia: pedido BA real de R$ 299,90, grupo 122.
//
// A AUTORIDADE AQUI E A MATRIZ VIGENTE, NAO A NOTA EMITIDA. A linha
// UFORIG=13/UFDEST=5/grupo 122 diz CODTRIB=0, ALIQUOTA=7, ALIQINTDEST=20,5,
// FCP=0. Isso da DIFAL 40,49, que e tambem a aliquota legal da BA desde
// fev/2024.
//
// A nota 895507 gravou 29,99 porque foi emitida com 17,0 — aliquota que a
// matriz ja superou (a linha foi corrigida em 2026-07-20). A nota-irma 895436
// do dia anterior usou 20,5. Nao "conserte" este teste para 29,99: fazer isso
// petrifica a nota defasada como verdade, que e exatamente a doenca que esta
// fatia existe para curar.
func TestTaxesForValueBahiaFollowsCurrentMatrix(t *testing.T) {
	rule := ICMSRuleInput{CodTrib: 0, Aliquota: f(7.0), AliqIntDest: f(20.5)}
	got := TaxesForValue(299.90, rule)

	if got.ICMS == nil || *got.ICMS != 20.99 {
		t.Fatalf("ICMS = V x 7%%: esperado 20.99, veio %v", got.ICMS)
	}
	if got.Difal == nil || *got.Difal != 40.49 {
		t.Fatalf("DIFAL = V x 20,5%% - ICMS: esperado 40.49, veio %v", got.Difal)
	}
	// Contra-controle: 29,99 e o valor da nota defasada. Se aparecer, alguem
	// trocou a fonte da matriz pelo documento emitido.
	if got.Difal != nil && *got.Difal == 29.99 {
		t.Fatalf("29.99 e a nota 895507, emitida a 17,0 — a matriz diz 20,5")
	}
}

func TestTaxesForValueChainClosedIsKnownZero(t *testing.T) {
	got := TaxesForValue(299.90, ICMSRuleInput{CodTrib: 60})
	for name, v := range map[string]*float64{"ICMS": got.ICMS, "DIFAL": got.Difal, "FCP": got.FCP} {
		if v == nil {
			t.Fatalf("%s: cadeia encerrada e zero CONHECIDO, nao desconhecido", name)
		}
		if *v != 0 {
			t.Fatalf("%s deveria ser 0 com CODTRIB 60, veio %v", name, *v)
		}
	}
}

// A "celula muda": o ERP diz que o DIFAL INCIDE (CODTRIB=0) mas nao diz
// QUANTO. Medido: 20 celulas assim no recorte real, concentradas em PR, RS e
// SC. Nao e caso hipotetico — e ~metade das UFs que vendemos. Selar o DIFAL e
// manter o ICMS e a propriedade central desta fatia.
func TestTaxesForValueMissingAliqIntDestSealsOnlyDifal(t *testing.T) {
	rule := ICMSRuleInput{CodTrib: 0, Aliquota: f(12.0), AliqIntDest: nil}
	got := TaxesForValue(100.00, rule)

	if got.ICMS == nil || *got.ICMS != 12.00 {
		t.Fatalf("ICMS e calculavel mesmo sem ALIQINTDEST: esperado 12.00, veio %v", got.ICMS)
	}
	if got.Difal != nil {
		t.Fatalf("DIFAL sem ALIQINTDEST tem que ser NIL, nunca 0")
	}
	if got.DifalReason == "" {
		t.Fatalf("desconhecido sem motivo legivel e so um buraco na tela")
	}
}

func TestTaxesForValueFCPOnlyWhenMatrixHasIt(t *testing.T) {
	comFCP := TaxesForValue(100.00, ICMSRuleInput{CodTrib: 0, Aliquota: f(12.0), AliqIntDest: f(20.0), FCPPercent: f(2.0)})
	if comFCP.FCP == nil || *comFCP.FCP != 2.00 {
		t.Fatalf("RJ: FCP 2%% de 100 = 2.00, veio %v", comFCP.FCP)
	}
	semFCP := TaxesForValue(100.00, ICMSRuleInput{CodTrib: 0, Aliquota: f(7.0), AliqIntDest: f(17.0)})
	if semFCP.FCP == nil || *semFCP.FCP != 0 {
		t.Fatalf("FCP ausente na matriz e zero CONHECIDO: nao existe segunda fonte de FCP")
	}
}

func TestTaxesForValueUncharacterizedCodTribIsUnknownNotZero(t *testing.T) {
	got := TaxesForValue(299.90, ICMSRuleInput{CodTrib: 10, Aliquota: f(7.0), AliqIntDest: f(17.0)})
	if got.ICMS != nil || got.Difal != nil {
		t.Fatalf("CODTRIB 10 nao foi caracterizado para PF — calcular com ele e inventar")
	}
}

func TestTaxesForValueAlwaysChargesOutboundPisCofins(t *testing.T) {
	got := TaxesForValue(299.90, ICMSRuleInput{CodTrib: 60})
	// 299.90 x 8.62% = 25.8514 -> 25.85
	if got.PisCofins == nil || *got.PisCofins != 25.85 {
		t.Fatalf("PIS/COFINS de saida e devido mesmo com a cadeia de ICMS encerrada: %v", got.PisCofins)
	}
}
```

> **Controles negativos, execute todos antes de implementar:**
> - Troque `20.5` por `17.0` no teste da Bahia → ele tem que dar 29,99 e falhar. É a prova de que o teste está ancorado na **matriz vigente**, e não na nota defasada que a sessão anterior tomou por autoridade.
> - Faça `CODTRIB=60` cair no ramo do `CODTRIB=0` → `TestTaxesForValueChainClosedIsKnownZero` tem que falhar.
> - Faça `AliqIntDest == nil` selar o ICMS junto com o DIFAL → `TestTaxesForValueMissingAliqIntDestSealsOnlyDifal` tem que falhar. **Este é o controle que prova que o selo é por componente**, que é a propriedade central da fatia.
> - Remova o PIS/COFINS → o último teste falha e a margem sobe ~8,6%.
>
> Se qualquer um desses controles não ficar vermelho, o teste correspondente não está medindo o que declara e **não conta como cobertura**.

- [ ] **Step 2: Rode e confirme o vermelho**

```bash
cd apps/server_core && GOCACHE=/c/Users/leandro.theodoro/Documents/marketplace-central/apps/server_core/.gocache go test ./internal/modules/orders/domain/ -run TestTaxesForValue -v
```
Esperado: `undefined: TaxesForValue`, 6 testes.

- [ ] **Step 3: Implemente**

```go
package domain

import (
	"math/big"
	"strconv"
)

// pisCofinsOutboundPct e o debito de PIS/COFINS da saida.
//
// A aliquota nominal e 9,25% sobre a BASE REDUZIDA (BASERED), medida em
// ~93,1% da base cheia nas notas de 2026. 9,25% x 93,1% = 8,61%; arredondado
// para 8,62% como aproximacao conservadora. Aplicamos sobre o valor CHEIO
// porque a base reduzida por item nao e conhecida ex-ante.
//
// Isto NAO e dupla contagem com o custo: TGFCUS.CUSSEMICM ja carrega o
// CREDITO de PIS/COFINS da entrada (formula METAL_CUSTO); o debito da saida
// nao esta em lugar nenhum hoje. Credito e debito somam — e o nao-cumulativo
// funcionando. Ver Documents/MNOS/docs/product/sankhya-custo-formula.md.
//
// Risco R5 da spec: e aproximacao, declarada aqui e nao escondida.
const pisCofinsOutboundPct = 8.62

// ICMSRuleInput e a regra da matriz do ERP achatada para o calculo. O dominio
// de orders nao importa o dominio de internal_read: o adapter traduz. Isso
// mantem o calculo puro e testavel sem nenhuma dependencia de leitura.
type ICMSRuleInput struct {
	CodTrib     int64
	Aliquota    *float64
	AliqIntDest *float64
	FCPPercent  *float64
	// Missing marca "nao existe linha na matriz para este par". Diferente de
	// linha existente com coluna vazia — e um conflito de ERP diferente, com
	// motivo diferente na tela.
	Missing bool
}

// TaxComponents sao os tributos ex-ante de UM item, ja em reais e arredondados
// a centavos. NIL e desconhecido e NUNCA zero (ADR-17). Cada NIL tem um Reason
// legivel ao lado: desconhecido sem motivo e so um buraco na tela.
type TaxComponents struct {
	ICMS       *float64
	Difal      *float64
	FCP        *float64
	PisCofins  *float64
	ICMSReason string
	DifalReason string
	FCPReason  string
}

// TaxesForValue aplica a regra fiscal do ERP sobre um valor. Funcao PURA.
func TaxesForValue(value float64, rule ICMSRuleInput) TaxComponents {
	var out TaxComponents

	// PIS/COFINS de saida independe do CODTRIB: e devido em toda venda, ate
	// quando a cadeia de ICMS esta encerrada por ST.
	pis := pctOfValue(value, pisCofinsOutboundPct)
	out.PisCofins = &pis

	switch {
	case rule.Missing:
		reason := "sem linha na matriz do ERP para este grupo e UF"
		out.ICMSReason, out.DifalReason, out.FCPReason = reason, reason, reason
		return out

	case rule.CodTrib == 60:
		// Cadeia encerrada por ST: zero CONHECIDO, nao ausente.
		zero := 0.0
		z1, z2, z3 := zero, zero, zero
		out.ICMS, out.Difal, out.FCP = &z1, &z2, &z3
		return out

	case rule.CodTrib != 0:
		reason := "CODTRIB " + strconv.FormatInt(rule.CodTrib, 10) +
			" nao caracterizado para comprador pessoa fisica"
		out.ICMSReason, out.DifalReason, out.FCPReason = reason, reason, reason
		return out
	}

	// CODTRIB 0 — tributado.
	if rule.Aliquota == nil {
		reason := "ALIQUOTA vazia na matriz do ERP"
		out.ICMSReason, out.DifalReason = reason, reason
	} else {
		icms := pctOfValue(value, *rule.Aliquota)
		out.ICMS = &icms

		if rule.AliqIntDest == nil {
			out.DifalReason = "ALIQINTDEST vazio na matriz do ERP"
		} else {
			// Base UNICA: comprador e sempre pessoa fisica nao-contribuinte
			// (LC 87/96 art. 4o §2o). Base dupla so existe para contribuinte.
			difal := round2(pctOfValue(value, *rule.AliqIntDest) - icms)
			out.Difal = &difal
		}
	}

	// FCP ausente da matriz e zero conhecido: a medicao mostrou que so o RJ
	// tem FCP e que nenhuma nota carregou FCP sem respaldo na matriz — nao
	// existe segunda fonte para "faltar".
	fcpPct := 0.0
	if rule.FCPPercent != nil {
		fcpPct = *rule.FCPPercent
	}
	fcp := pctOfValue(value, fcpPct)
	out.FCP = &fcp

	return out
}

// pctOfValue usa big.Rat e nao float para que 7% de 299,90 seja 20,993 -> 20,99
// e nunca um 20,98999999 binario. Mesmo motivo do `cem` que existia no adapter
// antigo (pricingtax/reader.go:20).
func pctOfValue(value, pct float64) float64 {
	amount := new(big.Rat).SetFloat64(value)
	if amount == nil { // NaN ou +-Inf: nao ha tributo honesto sobre nao-numero
		return 0
	}
	rate := new(big.Rat).Quo(new(big.Rat).SetFloat64(pct), big.NewRat(100, 1))
	amount.Mul(amount, rate)
	out, _ := strconv.ParseFloat(amount.FloatString(2), 64)
	return out
}

func round2(v float64) float64 {
	out, _ := strconv.ParseFloat(strconv.FormatFloat(v, 'f', 2, 64), 64)
	return out
}
```

> **Atenção ao arredondamento do DIFAL:** `V × ALIQINTDEST − ICMS` arredonda **o resultado**, não as parcelas separadamente. No caso BA: `299,90 × 20,5% = 61,4795`, menos o ICMS `20,99` = `40,4895 → 40,49`. Se ele falhar por 1 centavo, **não relaxe a tolerância** — mande `ESCALATION` ao hub com o valor obtido, porque um centavo aqui significa que a ordem das operações está diferente da do ERP e isso vira erro grande em outro valor.

- [ ] **Step 4: verde** · **Step 5: commit**

```bash
cd apps/server_core && GOCACHE=/c/Users/leandro.theodoro/Documents/marketplace-central/apps/server_core/.gocache go test ./internal/modules/orders/domain/ -run TestTaxesForValue -v
git add apps/server_core/internal/modules/orders/domain/tax_calc.go apps/server_core/internal/modules/orders/domain/tax_calc_test.go
git commit -m "feat(orders): calculo fiscal ex-ante puro, selado por componente"
```

---

## Task 7: Redesenhar a porta `TaxReader` — por item

**Files:**
- Modify: `apps/server_core/internal/modules/orders/ports/tax_reader.go:10-26`
- Test: `apps/server_core/internal/modules/orders/ports/` (contrato) + os testes de T8

**Por que a porta atual é máximo local:** `TaxesForOrder(ctx, total float64, destinoUF string)` decide imposto sobre o **total do pedido**. Mas o CODTRIB vem do `GRUPOICMS` do **produto**. Um pedido com um item ST e outro não-ST é **inexprimível** nessa assinatura. A porta codifica o modelo antigo (percentual chapado sobre o total) e não comporta o modelo medido. Decisão do operador 2026-08-02: redesenhar.

**Lembre-se:** 100% dos 38 pedidos reais têm 1 item. O defeito é **invisível em produção**. Só o fixture de 2 itens da Step 1 prova o redesenho.

- [ ] **Step 1: Escreva o teste que falha** (em `orders/adapters/internalread/tax_reader_test.go`)

```go
func TestTaxesForOrderMixedCodTribDoesNotCollapse(t *testing.T) {
	// Item A: grupo 122, CODTRIB 0, tributado.  Item B: grupo 311, CODTRIB 60, ST.
	// Nenhum pedido real tem 2 itens hoje — este fixture e a UNICA prova de
	// que a porta por item funciona.
	reader := newTaxReaderWithMatrix(t, matrixFixture)
	items := []ports.TaxableItem{
		{ItemIdentifier: "A", InternalProductID: intp(39563), Value: 100.00},
		{ItemIdentifier: "B", InternalProductID: intp(15956), Value: 100.00},
	}

	got, err := reader.TaxesForItems(context.Background(), items, "BA")
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	// A paga 7%; B nao paga nada (cadeia encerrada). Total 7.00, nao 14.00.
	if got.ICMS == nil || *got.ICMS != 7.00 {
		t.Fatalf("ICMS agregado deveria ser 7.00 (so o item tributado), veio %v", got.ICMS)
	}
	// PIS/COFINS incide nos dois: 8.62% x 200 = 17.24
	if got.PisCofins == nil || *got.PisCofins != 17.24 {
		t.Fatalf("PIS/COFINS incide nos dois itens: esperado 17.24, veio %v", got.PisCofins)
	}
}

func TestTaxesForItemsUnknownItemSealsWholeComponent(t *testing.T) {
	// Item sem vinculo (5 dos 38 pedidos reais estao assim: link_quality=conflict).
	reader := newTaxReaderWithMatrix(t, matrixFixture)
	items := []ports.TaxableItem{
		{ItemIdentifier: "A", InternalProductID: intp(39563), Value: 100.00},
		{ItemIdentifier: "B", InternalProductID: nil, Value: 100.00},
	}

	got, _ := reader.TaxesForItems(context.Background(), items, "BA")
	if got.ICMS != nil {
		t.Fatalf("um item sem CODPROD torna o ICMS do PEDIDO desconhecido: total parcial e mentira")
	}
	if len(got.ComponentesDesconhecidos) == 0 {
		t.Fatalf("o item desconhecido tem que aparecer nomeado")
	}
}
```

> **Controle negativo:** implemente somando os itens conhecidos e ignorando o item sem CODPROD. `TestTaxesForItemsUnknownItemSealsWholeComponent` tem que ficar vermelho. Total parcial apresentado como total é exatamente a doença que este plano trata — e é a mesma regra que `sumItemCosts` (`enrich_service.go:453-465`) já segue para o custo.

- [ ] **Step 2: Rode e confirme o vermelho**

- [ ] **Step 3: A porta nova**

```go
package ports

import "context"

// TaxableItem e um item do pedido no que interessa ao fisco: o produto (que
// determina o GRUPOICMS e portanto o CODTRIB) e o valor da linha.
type TaxableItem struct {
	ItemIdentifier    string
	InternalProductID *int
	Value             float64
}

// OrderTaxes sao os tributos ex-ante do pedido inteiro, agregados a partir dos
// itens. Cada componente NIL e desconhecido, nunca zero (ADR-17), e cada
// desconhecido aparece nomeado em ComponentesDesconhecidos com motivo legivel.
//
// Imposto e a SOMA de ICMS + FCP + PIS/COFINS e EXCLUI o DIFAL, porque
// BuildProfitability ja subtrai os dois separadamente
// (margem = total - comissao - frete - imposto - difal - custo). Somar o DIFAL
// aqui o cobraria duas vezes.
type OrderTaxes struct {
	Imposto   *float64
	Difal     *float64
	ICMS      *float64
	FCP       *float64
	PisCofins *float64

	ComponentesDesconhecidos []string
	Motivos                  map[string]string
}

// TaxReader resolve os tributos ex-ante de um pedido a partir da matriz fiscal
// do ERP. destinoUF e a SIGLA de duas letras como vem do shipment ("BA"); a
// traducao para o codigo numerico do ERP e responsabilidade do adapter.
type TaxReader interface {
	TaxesForItems(ctx context.Context, items []TaxableItem, destinoUF string) (OrderTaxes, error)
}
```

- [ ] **Step 4: verde** · **Step 5: commit**

```bash
git commit -m "refactor(orders): porta de imposto por item; total-level nao expressa CODTRIB por produto"
```

---

## Task 8: O adapter novo — orders lê a matriz do ERP

**Files:**
- Create: `apps/server_core/internal/modules/orders/adapters/internalread/tax_reader.go`
- Delete: `apps/server_core/internal/modules/orders/adapters/pricingtax/` (o pacote inteiro, com seus testes)

**Mudança em relação à spec §8:** a spec dizia "reescrever `pricingtax/reader.go` mantendo o nome do pacote". Com a porta redesenhada, o arquivo é substituído por inteiro e o nome `pricingtax` passa a ser **falso** — o adapter não adapta mais nada de `pricing`. Nome falso em pacote é a mesma classe de defeito que "custo sem ICMS" numa coluna que tem ICMS. Move para `adapters/internalread/`, ao lado do `CostReader` que ele espelha.

**Depende de T1** (tradução de sigla → código). Se T1 ainda estiver bloqueada, pare aqui e mande `BLOCKED`.

- [ ] **Step 1: Teste que falha** — os dois de T7, mais:

```go
func TestTaxReaderUnknownUFSealsInsteadOfDefaulting(t *testing.T) {
	reader := newTaxReaderWithMatrix(t, matrixFixture)
	got, err := reader.TaxesForItems(context.Background(), items, "XX")
	if err != nil {
		t.Fatalf("UF desconhecida nao e erro de leitura, e desconhecido selado: %v", err)
	}
	if got.ICMS != nil || got.Difal != nil {
		t.Fatalf("UF que nao sabemos traduzir tem que selar, nunca cair num default")
	}
}

// Ponta a ponta pela matriz: produto 39563 (grupo 122), destino BA.
// O fixture tem que carregar a linha REAL — CODTRIB=0, ALIQUOTA=7,
// ALIQINTDEST=20,5 — e o resultado tem que ser identico ao da T6. Se este
// teste e o de T6 discordarem, o adapter esta calculando por conta propria em
// vez de delegar para TaxesForValue.
func TestTaxReaderBahiaMatchesPureCalculation(t *testing.T) {
	reader := newTaxReaderWithMatrix(t, matrixFixture)
	items := []ports.TaxableItem{{ItemIdentifier: "A", InternalProductID: intp(39563), Value: 299.90}}
	got, err := reader.TaxesForItems(context.Background(), items, "BA")
	if err != nil {
		t.Fatalf("leitura da matriz: %v", err)
	}
	if got.ICMS == nil || *got.ICMS != 20.99 {
		t.Fatalf("ICMS: esperado 20.99, veio %v", got.ICMS)
	}
	if got.Difal == nil || *got.Difal != 40.49 {
		t.Fatalf("DIFAL: esperado 40.49, veio %v", got.Difal)
	}
	// imposto publicado = ICMS + FCP + PIS/COFINS, sem DIFAL (emenda E-3).
	if got.Imposto == nil || *got.Imposto != 46.84 {
		t.Fatalf("imposto = 20.99 + 0 + 25.85 = 46.84, veio %v", got.Imposto)
	}
}

// Celula muda ponta a ponta: produto do grupo 122 vendido para o PR, onde a
// matriz diz CODTRIB=0 sem ALIQINTDEST. O DIFAL sai desconhecido COM motivo, o
// ICMS sai calculado, e o componente entra em ComponentesDesconhecidos para a
// Tarefa 14 conseguir virar pendencia. Nada de 0, nada de branco mudo.
func TestTaxReaderMuteCellSurfacesPendency(t *testing.T) {
	reader := newTaxReaderWithMatrix(t, matrixFixture)
	items := []ports.TaxableItem{{ItemIdentifier: "A", InternalProductID: intp(39563), Value: 299.90}}
	got, err := reader.TaxesForItems(context.Background(), items, "PR")
	if err != nil {
		t.Fatalf("celula muda nao e erro de leitura: %v", err)
	}
	if got.ICMS == nil {
		t.Fatal("ICMS continua calculavel sem ALIQINTDEST")
	}
	if got.Difal != nil {
		t.Fatalf("PR nao tem ALIQINTDEST na matriz — DIFAL tem que ser desconhecido, veio %v", *got.Difal)
	}
	if !contains(got.ComponentesDesconhecidos, "difal") {
		t.Fatalf("desconhecido invisivel nao vira pendencia: %v", got.ComponentesDesconhecidos)
	}
	if got.Motivos["difal"] == "" {
		t.Fatal("motivo obrigatorio — a tela tem que dizer que falta ALIQINTDEST na matriz do ERP")
	}
	// O imposto publicado NAO carrega DIFAL, entao continua conhecido aqui.
	if got.Imposto == nil {
		t.Fatal("imposto = ICMS+FCP+PIS/COFINS nao depende do DIFAL")
	}
}
```

> **Controle negativo obrigatório:** faça o adapter cair num `0.0` quando `ALIQINTDEST` for nulo → `TestTaxReaderMuteCellSurfacesPendency` tem que falhar em `got.Difal != nil`. Este é o defeito histórico desta base (`pricingtax/reader.go:68-72` escreve `zero := 0.0` quando o DIFAL está desligado); se o controle não reprovar, o selo não existe.

- [ ] **Step 2: vermelho** · **Step 3: implemente**

Estrutura do adapter:

```go
package internalread

// TaxReader traduz a pergunta fiscal de orders para o leitor de matriz do
// internal_read: sigla -> codigo de UF, CODPROD -> GRUPOICMS, grupo+UF -> regra,
// e por fim regra+valor -> componentes (dominio puro de orders).
//
// Substitui o adapter pricingtax, que aplicava a aliquota de um perfil
// (pricing_calc_profiles) que nunca teve linha nenhuma: 4% de SIMPLES sobre um
// tenant em Regime Normal. O imposto exibido era a aliquota de um regime que
// ninguem configurou.
type TaxReader struct {
	matrix   MatrixSource
	ufCodes  UFCodeSource
	originUF internalreaddomain.UF // 13 = MG
}
```

Fluxo de `TaxesForItems`:

1. Traduzir `destinoUF` (sigla) → código. Falhou → **selar tudo**, motivo `"UF de destino desconhecida"`, sem erro.
2. Coletar os `InternalProductID` não-nulos → `GetTaxGroupsByIDs`.
3. Coletar os grupos não-nulos → `GetICMSRulesByGroups(originUF, grupos)`.
4. Por item: montar `ICMSRuleInput` (com `Missing: true` quando não houver linha) e chamar `domain.TaxesForValue(item.Value, rule)`.
5. **Agregar com a regra do `sumItemCosts`:** qualquer item com o componente nil torna o componente do **pedido inteiro** nil, e o `ItemIdentifier` entra em `ComponentesDesconhecidos` com o motivo.
6. `Imposto = ICMS + FCP + PisCofins` — e é nil se qualquer um dos três for nil.

- [ ] **Step 4: verde** · **Step 5: apague o pacote antigo**

```bash
git rm -r apps/server_core/internal/modules/orders/adapters/pricingtax/
```

> Só apague depois que T12 (fiação) compilar sem ele. Se apagar antes, `root.go:618` quebra o build e você perde o sinal dos outros testes.

- [ ] **Step 6: commit**

```bash
git commit -m "feat(orders): imposto lido da matriz do ERP; pricingtax e o perfil fabricado saem"
```

---

## Task 9: `enrich_service` — espelhar `resolveItemCosts`

**Files:**
- Modify: `apps/server_core/internal/modules/orders/application/enrich_service.go:405-423`
- Test: `apps/server_core/internal/modules/orders/application/enrich_service_test.go`

O molde já existe e está a 80 linhas de distância: `resolveItemCosts` (`:487-521`) faz exatamente esta forma para o custo — item sem `InternalProductID` vira desconhecido nomeado, e `sumItemCosts` (`:453-465`) recusa total parcial. **Siga o molde; não invente um segundo formato.**

- [ ] **Step 1: Teste que falha**

```go
func TestResolveTaxesSealsPerComponentAndNamesTheItem(t *testing.T) {
	svc := newEnrichServiceWithTaxes(t, taxReaderSealingDifal)
	enriched := svc.enrichOne(ctx, orderBA299)

	if enriched.Decomposicao.ICMS == nil {
		t.Fatalf("ICMS e conhecido neste par; selar tudo junto apaga numero que temos")
	}
	if enriched.Decomposicao.Difal != nil {
		t.Fatalf("DIFAL sem ALIQINTDEST tem que ficar NIL")
	}
	if !contains(enriched.ComponentesDesconhecidos, "difal") {
		t.Fatalf("componente selado tem que aparecer nomeado")
	}
	if enriched.Decomposicao.MargemValor != nil {
		t.Fatalf("margem com componente desconhecido e margem inventada")
	}
}
```

- [ ] **Step 2: vermelho** · **Step 3: substitua `resolveTaxes`**

```go
// resolveItemTaxes pede ao leitor fiscal os tributos ex-ante do pedido,
// item a item. Espelha resolveItemCosts: item sem InternalProductID vira
// desconhecido NOMEADO, e um componente desconhecido em qualquer item torna o
// componente do pedido inteiro desconhecido — total parcial apresentado como
// total e a mentira que esta fatia existe para matar.
func (s EnrichService) resolveItemTaxes(ctx context.Context, order domain.OrderReadModel, shipment *ShipmentEnrichment) ports.OrderTaxes {
	if s.taxes == nil || len(order.Items) == 0 {
		return ports.OrderTaxes{}
	}
	var destinoUF string
	if shipment != nil && shipment.DestinationUF != nil {
		destinoUF = *shipment.DestinationUF
	}

	items := make([]ports.TaxableItem, 0, len(order.Items))
	for _, item := range order.Items {
		items = append(items, ports.TaxableItem{
			ItemIdentifier:    itemIdentifier(item),
			InternalProductID: item.InternalProductID,
			Value:             lineValue(item),
		})
	}

	taxes, err := s.taxes.TaxesForItems(ctx, items, destinoUF)
	if err != nil {
		s.logger.Warn("orders: tax lookup failed",
			"provider_order_id", order.ProviderOrderID,
			"destino_uf", destinoUF,
			"error", err,
		)
		return ports.OrderTaxes{}
	}
	return taxes
}
```

> **`lineValue(item)`:** a base do tributo é `unit_price × quantity` da linha, **não** `order.Total`. `order.Total` inclui frete, e frete não é base de ICMS aqui. Se `lineValue` não existir, escreva-a ao lado de `sumItemCosts`, que já faz `UnitCost × Quantity`. Se `unit_price` for nil na linha, o item é **desconhecido** — não use o total do pedido como substituto.

- [ ] **Step 4: verde** · **Step 5: commit**

---

## Task 10: Decomposição com os quatro componentes

**Files:**
- Modify: `apps/server_core/internal/modules/orders/domain/order_decomposition.go`
- Test: `apps/server_core/internal/modules/orders/domain/order_decomposition_test.go:169`

**Não crie um segundo produtor.** `OrderDecomposition` já tem `ComponentesDesconhecidos []string` e `BuildProfitability` já anula a margem quando qualquer componente é nil. Estenda; duplicar o produtor foi o erro que derrubou a primeira versão do plano P2.

- [ ] **Step 1: Teste que falha**

```go
// ATENCAO: neste pedido a Comissao e o Difal valem O MESMO NUMERO (40.49) por
// coincidencia. Nao troque um pelo outro ao ler; e o unico lugar do plano onde
// isso acontece e uma asserção que confunda os dois passa por acidente.
func TestBuildProfitabilityKeepsImpostoExcludingDifal(t *testing.T) {
	in := ProfitabilityInput{
		Total: f(299.90), Comissao: f(40.49), Frete: f(23.65), Custo: f(154.53),
		ICMS: f(20.99), FCP: f(0), PisCofins: f(25.85), Difal: f(40.49),
	}
	got := BuildProfitability(in)
	// imposto = 20.99 + 0 + 25.85 = 46.84 ; difal entra separado
	if got.Imposto == nil || *got.Imposto != 46.84 {
		t.Fatalf("imposto = ICMS+FCP+PIS/COFINS, sem DIFAL: %v", got.Imposto)
	}
	// margem = 299.90 - 40.49 - 23.65 - 46.84 - 40.49 - 154.53 = -6.10
	if got.MargemValor == nil || *got.MargemValor != -6.10 {
		t.Fatalf("margem real do pedido BA: esperado -6.10, veio %v", got.MargemValor)
	}
	// MargemPct e FRACAO — o FE multiplica por 100. Nao multiplique aqui.
	if got.MargemPct == nil || round4(*got.MargemPct) != -0.0203 {
		t.Fatalf("margem pct e fracao (~-0.0203), nao -2.03: %v", got.MargemPct)
	}
}
```

> Este teste é o retrato do defeito inteiro: hoje a tela mostra **69,23 (23%)**; depois desta fatia mostra **−6,10 (−2,03%)**. **O pedido dá prejuízo.** Se ele passar com o número velho, nada foi consertado.
>
> **Controle negativo:** some o DIFAL dentro do `imposto` e veja a margem virar `−46,59`. Prova que o teste pega dupla contagem.
>
> **Cuidado com a margem negativa:** confira que nada no caminho (`BuildProfitability`, DTO, FE) clampeia margem em zero ou trata negativo como erro. Um clamp aqui reintroduz a mentira por outra porta. Se encontrar clamp, é achado de legacy → `ESCALATION`, não conserte de lado.

- [ ] **Step 2: vermelho** · **Step 3: estenda o struct e `BuildProfitability`** · **Step 4: verde** · **Step 5: commit**

---

## Task 11: Contrato, SDK e tela

**Files:**
- Modify: `contracts/api/marketplace-central.openapi.yaml:5852-5896` (`OrderDecomposicao`)
- Modify: `packages/sdk-runtime/src/index.ts:695,699`
- Modify: `apps/server_core/internal/modules/orders/transport/http_handler.go:428-439`
- Modify: `apps/web/src/pages/pedidos/PedidoDrawer.tsx:139-172`

**Regra do repo:** OpenAPI e `sdk-runtime` mudam **no mesmo commit**. `sdk-runtime` é escrito à mão — não há gerador para rodar.

- [ ] **Step 1: Teste que falha** (`http_handler_test.go`, padrão de `:330-358`)

```go
func TestDecomposicaoUnknownComponentIsAbsentNotZero(t *testing.T) {
	body := renderOrderJSON(t, orderWithSealedDifal)
	if strings.Contains(body, `"difal":0`) {
		t.Fatalf("desconhecido serializado como 0 e a mentira que esta fatia mata:\n%s", body)
	}
	if !strings.Contains(body, `"icms":20.99`) {
		t.Fatalf("componente conhecido tem que sair no JSON:\n%s", body)
	}
	if !strings.Contains(body, `"ALIQINTDEST vazio na matriz do ERP"`) {
		t.Fatalf("motivo legivel tem que chegar na tela, nao so o selo:\n%s", body)
	}
}
```

> **A asserção é sobre o JSON, não sobre o struct.** `omitempty` num `*float64` omite o nil corretamente, mas num `float64` transformaria `0` em ausente e ausente em `0` — os dois defeitos ao mesmo tempo. Só o corpo serializado prova a forma. **Controle negativo:** troque `*float64` por `float64` num componente e veja este teste ficar vermelho.

- [ ] **Step 2: vermelho** · **Step 3: os quatro campos**

Go (`:428-439`), mantendo o padrão `omitempty` dos irmãos:

```go
	ICMS      *float64 `json:"icms,omitempty"`
	FCP       *float64 `json:"fcp,omitempty"`
	PisCofins *float64 `json:"pis_cofins,omitempty"`
	MotivosDesconhecidos map[string]string `json:"motivos_desconhecidos,omitempty"`
```

OpenAPI, dentro de `OrderDecomposicao`:

```yaml
        icms:
          type: number
          format: double
          description: ICMS proprio da saida, lido da matriz fiscal do ERP. Ausente = desconhecido, nunca zero.
        fcp:
          type: number
          format: double
          description: Fundo de Combate a Pobreza. Zero conhecido quando a matriz nao tem FCP para a UF.
        pis_cofins:
          type: number
          format: double
          description: Debito de PIS/COFINS da saida (8,62% aprox. sobre base reduzida).
        motivos_desconhecidos:
          type: object
          additionalProperties:
            type: string
          description: Motivo legivel por componente selado, ex. "ALIQINTDEST vazio na matriz do ERP".
```

TypeScript (`index.ts`, dentro de `OrderDecomposicao`):

```ts
  icms?: number
  fcp?: number
  pis_cofins?: number
  motivos_desconhecidos?: Record<string, string>
```

- [ ] **Step 4: A tela**

Em `PedidoDrawer.tsx`, `DecompRow` já trata null via `formatMoney` (`pedidosFormatters.ts:71-74`, devolve null para null/undefined). Acrescente as linhas dos componentes e, quando o valor for ausente **e** houver motivo, mostre o motivo no lugar do traço:

```tsx
<DecompRow label="ICMS" value={d.icms} reason={d.motivos_desconhecidos?.icms} />
<DecompRow label="DIFAL" value={d.difal} reason={d.motivos_desconhecidos?.difal} />
<DecompRow label="FCP" value={d.fcp} reason={d.motivos_desconhecidos?.fcp} />
<DecompRow label="PIS/COFINS" value={d.pis_cofins} />
```

`DecompRow` passa a renderizar: valor formatado quando presente; `reason` em texto secundário quando ausente com motivo; o traço atual quando ausente sem motivo.

- [ ] **Step 5: Verde nos dois lados**

```bash
cd apps/server_core && GOCACHE=/c/Users/leandro.theodoro/Documents/marketplace-central/apps/server_core/.gocache go test ./internal/modules/orders/... -v
cd apps/web && npx --no-install tsc --noEmit
```

> **`tsc` verde não é prova sozinho** (memória `chip-import-chain-closed`): sem `node_modules` na raiz do worktree ele resolve `@mc/*` para o branch da main e passa vazio. Confirme que `apps/web/node_modules` existe **neste** checkout antes de acreditar no verde. Se não existir, mande `REQUEST` ao hub — não rode `npm install` por ritual.

- [ ] **Step 6: Commit único**

```bash
git add contracts/api/marketplace-central.openapi.yaml packages/sdk-runtime/src/index.ts apps/server_core/internal/modules/orders/transport/ apps/web/src/pages/pedidos/
git commit -m "contract(orders): componentes fiscais selados com motivo legivel"
```

---

## Task 12: Fiação, e a morte do perfil fabricado

**Files:**
- Modify: `apps/server_core/internal/composition/root.go:618-621`

- [ ] **Step 1: Trocar**

Sai:
```go
	ordersTaxReader := orderspricingtax.NewReader(pricingpostgres.NewCalcRepository(pool), cfg.DefaultTenantID)
```

Entra: construção do `internalread.TaxReader` sobre o `BatchReader` **cacheado** (`root.go:671`), com `originUF = 13`.

> **Note as duas instanciações:** `root.go:671` é a cacheada e `:724` é uma segunda, **sem cache**, para listings. O leitor fiscal tem que pendurar na cacheada. Se você pendurar na de `:724`, cada linha de pedido vai ao Oracle e a lista de pedidos (que já leva ~10,8s) fica inviável.

- [ ] **Step 2: Confirme que o import de `pricingtax` sumiu**

```bash
cd apps/server_core && grep -rn "pricingtax" internal/ && echo "AINDA EXISTE — nao terminou" || echo "limpo"
```

- [ ] **Step 3: Build + suíte inteira**

```bash
cd apps/server_core && GOCACHE=/c/Users/leandro.theodoro/Documents/marketplace-central/apps/server_core/.gocache go build ./... && GOCACHE=/c/Users/leandro.theodoro/Documents/marketplace-central/apps/server_core/.gocache go test ./... 
```

- [ ] **Step 4: Commit**

```bash
git commit -m "wire(orders): leitor fiscal do ERP no lugar do perfil SIMPLES 4% fabricado"
```

---

## Task 13: `/anuncios` — pior caso real, não pior caso fabricado

**Files:**
- Modify: `apps/server_core/internal/modules/internal_read/adapters/oracle/icms_ceiling.go:53-57`
- Modify: `apps/server_core/internal/modules/listings/ports/enrichment.go:9-15`
- Modify: `apps/server_core/internal/modules/listings/adapters/internalread/cost_reader.go:27-39`
- Modify: `apps/server_core/internal/modules/listings/application/read_service.go:539,792-826`
- Test: `read_service_test.go:713,750,776`

**Decisão do operador 2026-08-02:** o **conceito** "pior caso" fica — no anúncio ainda não existe comprador, logo não existe UF de destino, e perguntar "qual o pior destino" é legítimo. O que muda é **o número**.

Hoje: `MAX(ALIQUFDEST) + NVL(MAX(PERCICMSFCP),0)`, agrupado por UFDEST sobre linhas de especificidade diferente, ignorando o produto. Três defeitos empilhados:
1. `ALIQUFDEST` é a coluna errada (as notas reconciliam com `ALIQINTDEST`);
2. `MAX()` mistura a linha genérica com exceções por TOP;
3. **o produto não entra na conta** — um produto ST (cadeia encerrada, ICMS 0) recebe o mesmo teto de 17,5% de um produto tributado.

Depois: máximo sobre as UFs, com `ALIQINTDEST`, **pelo `GRUPOICMS` do próprio produto**. Produto ST → teto `0`.

**Isto muda quantos anúncios ficam marcados.** `below_margin_worst_case` dirige o selo `"Margem abaixo do mínimo no pior caso"` (`read_service.go:765`), os filtros `exception=below_margin` e `has_exception` (`:726,729`) e o contador do dashboard (`dashboard/service.go:74`). A contagem vai **cair** e passar a ser verdadeira. Isso é o resultado esperado, não uma regressão.

- [ ] **Step 1: Teste que falha**

```go
func TestWorstCaseCeilingIsZeroForSTProduct(t *testing.T) {
	// Produto do grupo 311: CODTRIB 60 em todas as UFs (cadeia encerrada).
	item := listingWithProduct(15956)
	got := worstCaseCeiling(item, matrixFixture)
	if got == nil || *got != 0 {
		t.Fatalf("produto ST tem teto 0: a cadeia esta encerrada, nao ha ICMS a pagar. Veio %v", got)
	}
}

func TestWorstCaseCeilingUsesAliqIntDest(t *testing.T) {
	// PR diverge: ALIQUFDEST 19, ALIQINTDEST 17. A nota reconcilia com 17.
	item := listingWithProduct(39563)
	got := worstCaseCeiling(item, matrixFixturePROnly)
	if got == nil || *got != 17.0 {
		t.Fatalf("teto tem que sair de ALIQINTDEST (17), nao de ALIQUFDEST (19). Veio %v", got)
	}
}

func TestBelowMarginCountDropsForSTProducts(t *testing.T) {
	// O selo e os filtros dependem disso. Produto ST marcado abaixo da margem
	// por um ICMS que ele nao paga e alarme falso.
	item := listingWithProduct(15956) // ST
	enriched := enrichListing(t, item, matrixFixture)
	if enriched.BelowMarginWorstCase != nil && *enriched.BelowMarginWorstCase {
		t.Fatalf("produto ST nao pode ser marcado abaixo da margem por ICMS inexistente")
	}
}
```

> **Controle negativo:** volte a query para `MAX(ALIQUFDEST)` e veja `TestWorstCaseCeilingUsesAliqIntDest` **e** `TestWorstCaseCeilingIsZeroForSTProduct` ficarem vermelhos. Se só um ficar, você consertou só metade.

- [ ] **Step 2: vermelho** · **Step 3: implemente**

`icms_ceiling.go` deixa de ter query própria e passa a derivar de `GetICMSRulesByGroups` (T4) — **um leitor de `TGFICM` no sistema, não dois**. A assinatura do port de listings ganha o grupo:

```go
type CostReader interface {
	GetCostFactsByIDs(context.Context, []int64) (map[int64]*CostFact, error)
	// GetICMSCeilingByGroups devolve, por grupo fiscal, o maior ICMS+FCP entre
	// as UFs de destino saindo de originUF. Substitui GetICMSCeilingByOrigin,
	// que ignorava o produto e usava a coluna errada.
	GetICMSCeilingByGroups(ctx context.Context, originUF int64, groups []int64) (map[int64]*ICMSCeiling, error)
}
```

`read_service.go` precisa do grupo do produto antes de pedir o teto — mesmo `GetTaxGroupsByIDs` de T3, sobre os `Link.ProductID` já coletados em `:520-527`.

> **`originUFMG = 13` continua hardcoded** (`read_service.go:29`). Não é escopo desta tarefa e não é o defeito aqui — mas é dívida real, e a **D-17** (`CODEMP=1` fixo) é da mesma família. Se você achar um terceiro hardcode de identidade da empresa, mande `ESCALATION`: três é padrão, não coincidência.

- [ ] **Step 4: Meça o impacto antes e depois**

```bash
cd apps/server_core && GOCACHE=/c/Users/leandro.theodoro/Documents/marketplace-central/apps/server_core/.gocache go test ./internal/modules/listings/... ./internal/modules/dashboard/... -v
```

E registre no pack de evidência **quantos anúncios ficavam marcados antes e quantos ficam depois**. Um número que muda sem ninguém contar não é evidência de nada.

- [ ] **Step 5: Commit**

```bash
git commit -m "fix(listings): pior caso de ICMS pelo grupo do produto e ALIQINTDEST; ST deixa de alarmar"
```

---

## Task 14: Detector de conflitos do ERP — os três calculáveis ex-ante

**Files:**
- Modify: `apps/server_core/internal/modules/orders/adapters/internalread/tax_reader.go`
- Test: idem

Requisito do operador: *"pegar erros e conflitos, ERP não é 100% verdade"*.

**Emenda E-2:** dos cinco detectores da spec §6, só três são calculáveis ex-ante. Os outros dois comparam a matriz contra a **nota emitida** — subsistema B. Movidos, **D-23**.

| detector | como | caso real |
|---|---|---|
| matriz sem linha | `GetICMSRulesByGroups` não devolve o par (grupo, UF) | 15956/RJ — grupo 311 sem RJ. **6 células medidas**, todas do grupo 311 |
| **célula muda** | `ALIQINTDEST` nulo com `CODTRIB=0` | 39563/PR, 39587/PR, 39587/RS. **20 células medidas**, em PR, RS e SC |
| TOP ganhou linha | `TGFICM` passa a ter linha para TOP 306 ou 313 | quebra a premissa da spec §3 |

Os dois primeiros **já são** os motivos de selo de T6 — não escreva um segundo produtor deles. O detector é a **superfície**: o motivo que já existe aparece como pendência, não só como traço na tela.

> **A célula muda é a mais importante das três, e não é caso de borda.** São 20 das 47 células do recorte real — PR, RS e SC inteiros. Decisão do operador (emenda E-5): **não digitamos alíquota nenhuma para tapar isso.** O DIFAL sai desconhecido, o par (produto, UF) vira pendência, e o texto da pendência tem que dizer **o que falta e onde**, porque quem vai consertar abre o Sankhya, não este código. Redação obrigatória, sem eufemismo:
>
> `Falta ALIQINTDEST na matriz de ICMS do ERP para o grupo <G> com destino <UF>. Sem ela o DIFAL não pode ser calculado.`
>
> O ERP **cobra** DIFAL do PR em nota (11 notas a 18,0) por um caminho que não está em `TGFICM`, `TSIUFS` nem `TGFPAR` e que não foi localizado. **Não saia procurando esse caminho** — não é escopo, e adivinhar o número é exatamente o defeito que esta fatia mata.

O terceiro é diferente: é um guard sobre uma premissa. A spec §3 assume que TOPs 306/313 não têm linha em `TGFICM`, e é isso que faz a linha genérica ser unívoca. Se um dia tiverem, a premissa **quebra em silêncio**.

- [ ] **Step 1: Teste que falha**

```go
func TestECommerceTOPWithMatrixRowBreaksTheGenericAssumption(t *testing.T) {
	// A spec §3 so vale porque as TOPs 306/313 nao tem linha em TGFICM.
	// Se ganharem, a precedencia entre linhas concorrentes passa a importar —
	// e ela nao foi decodificada.
	matrix := matrixFixtureWithTOPRow(306)
	_, conflicts := readerWith(matrix).TaxesForItems(ctx, items, "BA")
	if !hasConflict(conflicts, ConflictTOPHasMatrixRow) {
		t.Fatalf("TOP de e-commerce com linha propria tem que virar conflito, nao passar batido")
	}
}
```

> **Controle negativo:** remova a linha da TOP do fixture e veja o teste falhar por não achar o conflito. Um detector que nunca viu o caso que detecta não detecta nada.

- [ ] **Step 2: vermelho** · **Step 3: implemente** · **Step 4: verde** · **Step 5: commit**

```bash
git commit -m "feat(orders): detector dos conflitos de ERP calculaveis ex-ante"
```

---

## Task 15: QA dirigido em browser real

**Não é opcional** e não é auto-certificável: quem implementou não assina o QA (mandato do operador, memória `user-drive-validation-mandate`).

O chip **não sobe servidor**. Peça o dev stack ao hub via `REQUEST`.

- [ ] **P-1** — `/pedidos`, achar o pedido BA de R$ 299,90. Imposto = **46,84** (não 12,00). DIFAL = **40,49** (não 0,00, não 29,99 — 29,99 é a nota defasada). Screenshot.
- [ ] **P-2** — mesmo pedido: margem caiu de ~69,23 (23%) para **−6,10 (−2,03%)**. O pedido dá **prejuízo** e a tela tem que dizer isso, com sinal, sem clamp em zero. Screenshot lado a lado com o valor anterior.
- [ ] **P-3** — um dos pares selados (15956/RJ, sem linha) mostra **motivo legível**, não número, não traço mudo.
- [ ] **P-3b** — **célula muda**: um pedido com destino **PR** mostra ICMS calculado, DIFAL selado com motivo, e a pendência nomeando `ALIQINTDEST` e o grupo. Este é o retrato de 20 das 47 células do recorte — se ele sair branco ou 0,00, a fatia falhou no caso mais comum, não no raro.
- [ ] **P-4** — um pedido com item `link_quality='conflict'` (5 existem): componentes fiscais selados, margem vazia, item nomeado.
- [ ] **P-5** — `/anuncios`: contar anúncios marcados "abaixo da margem" **antes** (main) e **depois**. Confirmar que algum produto ST deixou de ser marcado. Registrar os dois números.
- [ ] **P-6** — dashboard: o contador `below_margin` acompanha a queda de P-5. Se não acompanhar, os dois caminhos divergiram.
- [ ] **P-7** — controle positivo: apagar a linha da matriz do par BA no fixture, recarregar, e ver o ICMS **sumir com motivo**. Sem isso, o verde de P-1 não distingue "leu certo" de "não leu nada".

> P-7 é o que separa este QA de um QA vazio. Repetir a leitura não prova que a leitura funciona — só o controle positivo injetado prova (memória `p1-validated-positive-control`).

---

## 3. Auto-revisão deste plano

**Cobertura da spec:**

| spec | tarefa |
|---|---|
| §3 fonte `TGFICM` por `GRUPOICMS` | T3, T4 |
| §3 `ALIQINTDEST` e não `ALIQUFDEST` | T4 (teste de query), T13 |
| §4 o cálculo | T6 |
| §4 o que não subtrair de novo | T6 (comentário de `pisCofinsOutboundPct`) |
| §5 selo por componente | T2, T6, T7, T10, T11 |
| §5 motivo legível | T2, T6, T11 |
| §6 detectores | T14 — **3 dos 5**, emenda E-2 |
| §5/§9 números da Bahia | **CORRIGIDOS pela emenda E-4** — a spec cita 29,99 como autoridade; vale **40,49** (matriz vigente). T6, T8, T10 e P-1 já estão no número certo |
| cobertura de `ALIQINTDEST` | **CORRIGIDA pela emenda E-5** — 21 calculáveis, 20 mudas, 6 sem linha. D-24 |
| §7 congelar no momento da venda | **LACUNA — ver abaixo** |
| §7 cache por `TGFHICM` | T5 |
| §8 onde encaixa | T8, T9, T12 |
| §8 `pricing_difal_rates` morre por falta de consumidor | T8, T12 (`grep pricingtax`) |
| §9 testes | T2, T4, T6, T7, T10, T11, T13, T14 |
| §9 verificação em browser | T15 |

**Lacuna que eu não fecho e não escondo — spec §7.1, congelar o imposto no momento da venda.**

A spec exige que o imposto **congele**: recalcular pela matriz atual reescreve margem histórica sozinha quando um estado muda alíquota (e as alterações estão acelerando — 3 em 2024, 34 em 2025, 62 até ago/2026). Este plano **não implementa o congelamento**: ele lê a matriz atual a cada leitura de pedido.

Não é esquecimento, é dependência. Congelar exige **persistir** os componentes fiscais no pedido no momento da venda — tabela nova, migração, e a decisão de o que fazer com os 38 pedidos que já existem sem carimbo. Isso é uma fatia própria, com costura própria (`orders/**` escreve, mas é schema novo). Enfiá-la aqui dobraria o tamanho de uma fatia que já toca 15 arquivos.

**Consequência aceita enquanto não for feito:** a margem histórica se move quando o ERP muda alíquota. Hoje ela se move de 12,00 fixo para 12,00 fixo, ou seja, está errada e imóvel; depois desta fatia fica certa e móvel. É melhor, e não é o fim.

**Ação:** registrar como **D-22** e propor a fatia ao operador ao fechar. Não silenciar.

**Placeholders:** varrido. Os únicos "a definir" são a Tarefa 1 (bloqueada, com critério de aceite fechado, faltando só os literais da medição do ERP) e o nome do helper de binds em T3 — que é uma instrução para **usar o existente**, com `ESCALATION` se ele não existir, não uma licença para inventar.

**Consistência de tipos:** `ICMSRule` (T2, internal_read) e `ICMSRuleInput` (T6, orders) são **de propósito** dois tipos — o domínio de orders não importa o domínio de internal_read, e o adapter de T8 traduz. `TaxComponents` (T6, por item) e `OrderTaxes` (T7, agregado do pedido) também são distintos de propósito. `TaxesForOrder` foi renomeada para `TaxesForItems` em **todos** os pontos: T7 (porta), T8 (adapter), T9 (chamador).

---

## 4. Dívidas que esta fatia cria ou move

| id | dívida | estado |
|---|---|---|
| D-18 | ST recuperável dentro de `CUSSEMICM` — custo superestimado, margem pessimista | mantida da spec |
| D-19 | teto de `/anuncios` com `ALIQUFDEST` e `MAX()` | **RESOLVIDA** na T13 (emenda E-1) |
| D-20 | invalidação por `TGFHICM` é melhor que TTL, não é garantia | mantida da spec |
| D-21 | `NewDefaultCalcProfile()` SIMPLES 4% sobrevive em `pricing`; `/precos` continua fabricando em 7 sítios (`calc_service.go:72-73,166,189,256,329,449`) | mantida, **e é a próxima fatia candidata** |
| D-22 | imposto **não congela** no momento da venda — margem histórica se move quando o ERP muda alíquota | **NOVA**, ver auto-revisão |
| D-23 | detectores "previsto ≠ realizado" e "alíquota mudou sem log" movidos para o subsistema B | **NOVA**, emenda E-2 |
| D-24 | **DIFAL não precificável em PR, RS e SC** — 20 células com `CODTRIB=0` e sem `ALIQINTDEST`. Custo medido: 7 notas em 12 meses, de 74. Sai selado + pendência (E-5); o conserto é cadastro no Sankhya, não código | **NOVA**, decisão do operador |
| D-25 | **a matriz do ERP e os documentos do ERP discordam entre si** — pedido 895422 diz `CODTRIB=60, 17,0`, a NF-e do mesmo pedido diz `CODTRIB=0, 20,5`. Não se sabe se é regra (recálculo no faturamento) ou artefato das 6 edições de 2026-07-20. Só o subsistema B (ex-post) pode medir isso | **NOVA** |
| D-26 | **o ERP cobra DIFAL do PR** (11 notas a 18,0) de uma fonte que não está em `TGFICM`, `TSIUFS` nem `TGFPAR` e **não foi localizada**. Enquanto não for, PR fica em D-24 | **NOVA** |
