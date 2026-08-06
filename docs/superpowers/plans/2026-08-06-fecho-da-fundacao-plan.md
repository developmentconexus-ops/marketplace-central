# Plano — Fecho da Fundação

**Objetivo:** deixar a fundação (kernel + analisadores + fatia `catalog`) sem dívida antes do
próximo plano, que migra contexto a contexto. Hoje a árvore **não compila**, a lane de unidade
está **vermelha**, a fatia tem **zero linhas em produção**, e 6 das 15 regras do design caíram
para nível 3 (convenção humana).

**Design que governa:** `docs/superpowers/specs/2026-08-06-protocolo-de-codigo-design.md`.
**Plano anterior:** `docs/superpowers/plans/2026-08-06-fundacao-kernel-e-fatia-catalog.md`
(Tarefas 1–12, fechadas em `401cd842`..`1b2ef2da`).

---

## Global Constraints

Vinculam **toda** tarefa. Um executor que as viole reprova a revisão da tarefa.

1. **Nada em `apps/server_core/internal/modules/` é tocado.** A árvore legada é inventariada,
   nunca crescida (ARCHITECTURE.md decisão congelada 11).
2. **Zero dependências novas.** `go.mod`, `package.json` e `package-lock.json` não mudam.
   Se uma tarefa parecer precisar de uma, é `REQUEST` ao hub, não `go get`.
3. **RED antes do código, em toda tarefa.** Escreve o teste que falha, corre-o, **cola o output
   verbatim** no relatório. Alegar que falhou não é prova. Um passo que diz "verificar que falha"
   sem o texto esperado da falha não é controlo negativo.
4. **Nunca** `git push`, `reset`, `revert`, `stash`, `clean`.
5. `git status --porcelain --untracked-files=all` **vazio** antes de fechar cada tarefa.
6. **Bloqueio:** escreve em `.mnfs/HARNESS-DEBTS.md` com `file:line` e o output verbatim, faz
   commit do que está verde, passa à tarefa seguinte independente. Não inventes contorno.
7. **Nunca despejar um ambiente.** Sem `docker inspect`, sem `docker exec … env`, sem `printenv`
   nu, sem `Get-ChildItem Env:`. Diagnostica **uma** variável pelo nome: `printenv A_VAR`.
   Nunca `. .env` na shell.
8. **Oracle é só leitura.** Nenhuma tarefa escreve no ERP.
9. **Desconhecido nunca vira zero, `""`, `false` ou default plausível** (ADR-017 / kernel `fact`).
10. **`tenant_id` em toda query.** Sem predicado de tenant, a query não entra.
11. Comandos Go correm de `apps/server_core` com `GOCACHE` absoluto:
    `export GOCACHE="$(pwd)/.gocache"`.

---

## Medição

As oito respostas da Fase 1 do `mc-planning`. Medido na árvore `main @1b2ef2da`, 2026-08-06.
O executor herda estas medições — não as re-derives.

### 1. O que já existe que faz parte disto

| Coisa | Onde | Estado medido |
|---|---|---|
| Kernel de 6 membros | `internal/kernel/{tenant,channel,exact,provenance,period,fact}` | existe, compila, testado |
| Detectores nível 2 | `internal/arch/scan.go:127,167,203` | 3 detectores; 1 (vendor) sem filtro de caminho nenhum |
| Fatia `catalog` | `internal/contexts/catalog/` | domínio + aplicação + postgres + porta escritos |
| Adaptador Sankhya | `internal/adapters/erp/sankhyaoracle/` | fachada + mapper + rows escritos |
| Esquema | `migrations/0097_catalog_context.sql` | 4 tabelas, RLS FORCE, políticas |
| Portão | `scripts/arch-gate.sh` | 6 passos |
| Cursor de página legado | `modules/internal_read/ports/catalog_page.go:28` | `Cursor{InternalProductID}` — o antipadrão a **não** repetir |

### 2. Onde o defeito vive de facto

Sete defeitos, todos **causa** e não sintoma:

| # | Causa | `file:line` |
|---|---|---|
| D-a | Fixture de contagem de migrações desactualizada → **lane de unidade vermelha** | `internal/platform/migrate/runner_test.go:26` (`want 83`; há 84 `.sql`) |
| D-b | `catalog.New` tem 2 de 3 parâmetros tipados por `catalog/internal/application` → **não tem chamador legal fora da árvore do contexto** | `internal/contexts/catalog/module.go:21` |
| D-c | Raiz de composição nomeia interno do contexto → **`go build ./...` falha** | `internal/composition/catalog_wiring.go:9` |
| D-d | `ScanCrossContextInternalSuffix` sai cedo em ficheiro fora de `contexts/` → o único sítio que viola a regra é invisível ao instrumento | `internal/arch/scan.go:131-133` |
| D-e | `ScanVendorTokensSuffix` **não filtra caminho nenhum** → acusa a sua própria lista e nunca poderá ficar verde | `internal/arch/scan.go:203-239` vs lista em `:33-36` |
| D-f | Proveniência **fabricada** na rehidratação (`system="catalog"`) e chaves de fonte extra **perdidas** (o `Apply` do loop devolve `Idempotent` porque o hash é o mesmo) | `repository.go:203`, `repository.go:219-226` |
| D-g | `Fact.Value()` com o booleano descartado — exactamente o padrão que o pacote existe para impedir; o comentário promete um detector que não existe | `repository.go:360`, promessa em `fact/knowledge.go:117-120` |

### 3. Quem mais consome o caminho que vai mudar

- `catalog.New` — **um** chamador: `composition/catalog_wiring.go:27`. Não compila.
- `WireCatalog` — **zero** chamadores (grep em todo `apps/server_core`).
- `catalogfeed.Feed.Page` — **zero** chamadores fora de `_test.go`.
- `ScanCrossContextInternal` / `ScanVendorTokens` — chamados só por
  `internal/arch/cmd/archscan/main.go:19-23` e por `scan_test.go`.
- `runner_test.go` — dois testes no pacote, **ambos** falham.

Ou seja: o raio de explosão de todas as mudanças de assinatura deste plano é o próprio plano.
Nada em `internal/modules/` chama nada disto.

### 4. O que o contrato já diz

- **Nada.** `contracts/governance/modules.json` não menciona `contexts`, `kernel` nem
  `adapters/erp` — zero ocorrências. A árvore nova é **ingovernada**: nenhum
  `GOV_MODULE_COVERAGE`, `GOV_MODULE_DEPENDENCY` ou `GOV_MODULE_LAYER` a alcança.
- `contracts/governance/invariants.json:11` — `production-panic` tem escopo
  `apps/server_core` (a árvore toda), portanto **alcança** os 3 `panic` novos:
  `exact/decimal.go:84`, `exact/decimal.go:141`, `fact/knowledge.go:134`. Nenhum tem excepção
  registada. Pânico baselineado não se herda.
- OpenAPI e `sdk-runtime`: **este plano não expõe rota nova**, logo nenhuma tarefa toca no par.
- **Conflito de ordem de verdade (stop-and-classify):** `ARCHITECTURE.md:39` decisão congelada 7
  diz *"External marketplace integrations enter only through the `connectors` module via port
  interfaces"*; o design ratificado põe integrações em `adapters/marketplace/<vendor>`.
  `ARCHITECTURE.md:5` exige ADR explícito para mexer numa decisão congelada. → **Tarefa 11.**

### 5. Estado vivo real

Container `marketplace-central-postgres-1`, Up 24h, healthy:

| Tabela | `count(*)` |
|---|---|
| `catalog.products` | **0** |
| `catalog.product_identifiers` | **0** |
| `catalog.source_product_keys` | **0** |
| `catalog.source_observations` | **0** |

`pg_roles`: `marketplace` tem `rolsuper=t`, `rolbypassrls=t`; `pg_tables` dá
`tableowner = marketplace` nas 4 tabelas com `rowsecurity = t`. **`FORCE ROW LEVEL SECURITY` é
decorativo pela ligação da aplicação** — dono + superuser + bypassrls passam por cima.
`grep -rl "CREATE ROLE\|GRANT " apps/server_core/migrations/` → **zero ficheiros**.

`marketplace-central-backend-1` Exited(255) há 24h; `frontend-1` Exited(1) há 28h. Só o Postgres
está de pé.

### 6. O que prova que está partido hoje, e o que provará consertado

| Defeito | Prova de partido (hoje) | Prova de consertado |
|---|---|---|
| D-a | `go test ./internal/platform/migrate/` → `FAIL` | `ok` |
| D-b/D-c | `go build ./...` → `use of internal package … not allowed` | build limpo |
| D-d | `archscan -root internal/composition` → 0 findings, com a violação lá | ≥1 finding no `catalog_wiring.go` pré-conserto (controlo positivo em testdata) |
| D-e | `archscan -root internal/arch` acusa `scan.go` | 0 findings em `internal/arch`, ≥1 no fixture de controlo positivo |
| D-f | teste que rehidrata produto com 2 chaves de fonte → devolve 1 | devolve 2, e `Evidence().System()=="sankhya"` |
| D-g | detector novo acusa `repository.go:360` | 0 findings + fixture de controlo positivo acusa |
| fatia órfã | `count(*) = 0` nas 4 tabelas | `count(*) > 0` em `catalog.products` **e** `catalog.source_observations` |
| RLS | ligação da app lê linhas de outro tenant | ligação pelo papel novo lê **0** linhas de outro tenant |

**Aceitação é observável, nunca teste verde.** Verde prova a unidade; a linha prova a
funcionalidade (§P3 do `mc-planning`).

### 7. Orçamento de custo

Zero chamadas a provider. Este plano não toca Mercado Livre, portanto não consome o balde de
tokens partilhado (`resilience_decorator.go`, 900/min). A Tarefa 7 lê Oracle: uma varredura
completa de `TGFPRO` activo em páginas de 500, **leitura apenas**, na lane Oracle em Docker.
Postgres: 4 INSERT por produto na primeira passagem, 0 nas seguintes (idempotência por hash).

### 8. O que falha em silêncio às 3 da manhã

Hoje: **tudo**. A fatia não tem chamador de produção, portanto não pode falhar — e não pode
funcionar. Depois deste plano, o ingest é um comando de operador (`cmd/catalogingest`), que sai
com código != 0 e imprime o erro; não há ecrã. **Visibilidade em ecrã é escopo do próximo plano**
e fica registada como dívida nomeada na Tarefa 12, não como comentário.

---

## O que já existe (varredura anti-redundância)

Cada artefacto novo abaixo cita o mais próximo que existe e porque não serve.

| Artefacto novo | Mais próximo que existe | Porque não serve |
|---|---|---|
| `port.ProductFeed` + `port.Cursor` opaco | `modules/internal_read/ports/catalog_page.go:28` `Cursor{InternalProductID int64}` | É precisamente o vazamento a evitar: o identificador do ERP está **na assinatura**, logo um segundo provedor não a pode implementar. É legado; inventaria-se, não se estende (constraint 1). |
| `fact.Combine2` / `fact.Map` | nenhuma aritmética em `kernel/fact/` — **zero métodos** além dos acessores (`knowledge.go:121-151`) | Não existe. A Regra 4.4 do design mandava métodos em `Fact[T]`; Go não permite métodos com parâmetros de tipo próprios, logo a regra é impossível como escrita → emenda na Tarefa 11. |
| `provenance.Derived` | `provenance.NewEvidence` (`evidence.go:30`) exige `externalKey` e `observedAt` de **uma** fonte | Um valor derivado de duas evidências não tem uma origem única; forçar `NewEvidence` obrigaria a escolher uma e mentir sobre a outra. |
| `arch.RuleFactValueDiscard` | os 3 detectores em `scan.go:25-28` | Nenhum olha para descarte de booleano; o próprio `knowledge.go:119-120` promete este detector. |
| Papel de base de dados de menor privilégio | zero `CREATE ROLE`/`GRANT` em 84 migrações | Não existe. Sem ele o `FORCE RLS` da 0097 é decorativo (§Medição 5). |
| `cmd/catalogingest` | `WireCatalog` (`catalog_wiring.go:25`), zero chamadores | Monta a fatia mas nada a corre. Falta o sítio de composição, não a montagem. |

**Máximo local vs global (Fase 3), respondido por escrito:**

1. *Quantas cópias do conceito existem?* Cursor de página: **2** (o legado em `internal_read` e o
   novo em `catalogfeed.Page`). Convergem na porta — por isso o cursor opaco nasce em
   `contexts/catalog/port/`, não no adaptador.
2. *A causa está uma camada abaixo do sintoma?* Sim, em D-b: o sintoma é o build partido na
   composição; a causa é a assinatura da fachada. Consertar a composição (importando o interno)
   seria o máximo local — e é literalmente o que está lá hoje.
3. *Já existe costura para isto?* Para o agendamento, sim: `syncapp.Scheduler`. **Não a usamos
   neste plano** porque `sync/domain/entity.go` tem lista fechada de entidades e já existe job
   `products` legado; registar um segundo colidiria em `sync_state`, e mexer em
   `internal/modules/sync` viola a constraint 1. O agendamento é decisão do próximo plano; aqui o
   sítio de composição é um comando de operador. Registado como dívida na Tarefa 12.
4. *Estou a estender legado para resolver problema actual?* Não. Nenhuma tarefa escreve em
   `internal/modules/`.

---

## Tarefas

### Tarefa 1 — Lane de unidade verde: fixture de contagem de migrações

**Porquê primeiro:** enquanto isto estiver vermelho, nenhum outro sinal de teste é legível.

**RED.** De `apps/server_core`:

```bash
export GOCACHE="$(pwd)/.gocache" && go test ./internal/platform/migrate/
```

Output esperado, verbatim (cola-o no relatório):

```
--- FAIL: TestCanonicalMigrationsMatchFixture (0.00s)
    runner_test.go:26: fixture inventory drift: got 84 canonical migrations, want 83
--- FAIL: TestRunnerResolvesMigrationsFromForeignCWD (0.00s)
    runner_test.go:65: foreign CWD returned 84 migrations, want 83
FAIL
```

**Medir antes de escrever o número.** Não confies em 84:

```bash
ls apps/server_core/migrations/*.sql | wc -l
```

**GREEN.** Em `apps/server_core/internal/platform/migrate/runner_test.go`, substitui os dois
literais `83` pelo número medido (84), nas linhas 26 e 65 e nas mensagens `want %d`.

**Aceitação:** `go test ./internal/platform/migrate/` → `ok`.

Commit: `fix(migrate): fixture de inventario acompanha a 0097`

---

### Tarefa 2 — Build verde: a fachada monta-se a si própria

**Causa (D-b):** `catalog.New` recebe `application.Store` e `application.IDFactory`, ambos
declarados em `catalog/internal/application`. Um construtor exportado cujo parâmetro é um tipo
interno **não é chamável de fora da árvore** — Regra 2.2-a do design. A composição só conseguiu
chamá-lo importando o interno, e é isso que parte o build.

**RED.** De `apps/server_core`:

```bash
export GOCACHE="$(pwd)/.gocache" && go build ./...
```

Output esperado, verbatim:

```
internal/composition/catalog_wiring.go:9:2: use of internal package marketplace-central/apps/server_core/internal/contexts/catalog/internal/postgres not allowed
```

**GREEN — `apps/server_core/internal/contexts/catalog/module.go`**, ficheiro completo:

```go
// Package catalog is the context's façade. The composition root builds a Module
// and hands it to consumers as a port; it never reaches past this file, because
// everything past this file is under internal/.
package catalog

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"marketplace-central/apps/server_core/internal/contexts/catalog/contracts"
	"marketplace-central/apps/server_core/internal/contexts/catalog/internal/application"
	"marketplace-central/apps/server_core/internal/contexts/catalog/internal/postgres"
	"marketplace-central/apps/server_core/internal/contexts/catalog/port"
)

// Module is Catalog, assembled.
type Module struct {
	service *application.Service
	reader  port.Reader
}

// New assembles the context from the ONLY thing an outsider may legitimately
// name: a connection pool.
//
// Every other collaborator — the store, the id factory, the reader — is chosen
// here, inside the context, because their types live under internal/ and a
// parameter typed by one of them would have no legal caller. That was the state
// this constructor was in: the composition root could only satisfy it by
// importing catalog/internal/postgres, which the compiler refused. The refusal
// was correct; the signature was the defect.
func New(pool *pgxpool.Pool) *Module {
	repo := postgres.NewRepository(pool)
	return &Module{
		service: application.NewService(repo, postgres.NewULIDFactory()),
		reader:  postgres.NewSummaryReader(repo),
	}
}

// IngestProduct folds one source observation into the catalogue.
func (m *Module) IngestProduct(ctx context.Context, o contracts.ProductObservation) (contracts.IngestResult, error) {
	return m.service.Ingest(ctx, o)
}

// Reader is Catalog's answer to identity questions from other contexts.
func (m *Module) Reader() port.Reader { return m.reader }
```

**GREEN — `apps/server_core/internal/composition/catalog_wiring.go`**, ficheiro completo:

```go
package composition

import (
	"database/sql"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"marketplace-central/apps/server_core/internal/adapters/erp/sankhyaoracle"
	"marketplace-central/apps/server_core/internal/contexts/catalog"
)

// CatalogWiring is the assembled catalog slice.
//
// Note what this file CANNOT name: nothing under any context's or adapter's
// internal/. That is not a convention — the import above was tried and the
// compiler refused it. This struct is what is left when the refusal is obeyed.
type CatalogWiring struct {
	Module *catalog.Module
	Feed   sankhyaoracle.Bundle
}

// WireCatalog assembles the slice from its two edges.
func WireCatalog(pool *pgxpool.Pool, oracleDB *sql.DB, sankhyaInstance string) (CatalogWiring, error) {
	bundle, err := sankhyaoracle.New(oracleDB, sankhyaInstance, time.Now)
	if err != nil {
		return CatalogWiring{}, err
	}
	return CatalogWiring{Module: catalog.New(pool), Feed: bundle}, nil
}
```

**Nota para o executor:** o teste de integração da fatia (Tarefa 12 do plano anterior) chamava
`catalog.New` com três argumentos. Ajusta-o para `catalog.New(pool)` — **não o apagues**. Se
usava dobras em memória, converte-o para usar o Postgres da lane hermética, ou, se não for
possível na tarefa, **restaura-o como teste de unidade em
`catalog/internal/application/`** onde as dobras são legais. Teste de outrem restaura-se ou
reafirma-se, nunca se apaga.

**Aceitação:** `go build ./...` e `go vet ./...` sem output. `go test ./internal/...` verde.

Commit: `fix(catalog): fachada monta-se a si propria; composicao deixa de nomear interno`

---

### Tarefa 3 — Detector cross-context vê fora dos contextos (Regra 1.3)

**Causa (D-d):** `scan.go:131-133` faz `here, inContext := contextOf(path); if !inContext {
return nil }`. Ficheiro fora de `contexts/` — a raiz de composição, um adaptador, um `cmd` — é
saltado. O **único** sítio que alguma vez violou a regra estava fora, e o instrumento era cego a
ele. Um instrumento que devolveria o mesmo resultado se a resposta fosse a oposta é cego.

**RED.** Cria o fixture `apps/server_core/internal/arch/testdata/outside_context.go.txt`:

```go
package composition

import (
	_ "marketplace-central/apps/server_core/internal/contexts/catalog/internal/postgres"
)
```

E acrescenta a `apps/server_core/internal/arch/scan_test.go`:

```go
// TestCrossContextInternalSeesOutsideContexts is the positive control for the
// blindness that shipped: the composition root imported catalog/internal/postgres
// and the detector reported zero findings, because it skipped every file that was
// not itself inside contexts/.
func TestCrossContextInternalSeesOutsideContexts(t *testing.T) {
	got, err := ScanCrossContextInternalSuffix("testdata", ".go.txt")
	if err != nil {
		t.Fatalf("scan: %v", err)
	}
	var hit *Finding
	for i := range got {
		if strings.HasSuffix(got[i].File, "outside_context.go.txt") {
			hit = &got[i]
			break
		}
	}
	if hit == nil {
		t.Fatalf("test=TestCrossContextInternalSeesOutsideContexts: no finding for a file outside contexts/ importing catalog/internal; got %d finding(s): %+v", len(got), got)
	}
	if hit.Rule != RuleCrossContextInternal {
		t.Fatalf("rule = %q, want %q", hit.Rule, RuleCrossContextInternal)
	}
}
```

Corre e cola o output. Deve falhar com `no finding for a file outside contexts/`.

**GREEN.** Em `apps/server_core/internal/arch/scan.go`, substitui o corpo do visitor de
`ScanCrossContextInternalSuffix` (linhas 129-151) por:

```go
	err := walk(root, suffix, func(path string, fset *token.FileSet, file *ast.File) error {
		// here is "" for a file that lives outside any context — the composition
		// root, an adapter, a cmd. Those are NOT skipped: the one import that
		// ever broke this rule came from exactly there, and a detector that
		// starts by skipping them can only ever report zero.
		here, _ := contextOf(path)
		for _, imp := range file.Imports {
			value, uErr := strconv.Unquote(imp.Path.Value)
			if uErr != nil {
				continue
			}
			target, reaches := importedContextInternal(value)
			if !reaches {
				continue
			}
			// A context reaching into its OWN internal is the design; anything
			// else, including here == "", is a finding.
			if here != "" && target == here {
				continue
			}
			from := here
			if from == "" {
				from = "outside any context"
			}
			out = append(out, Finding{
				File:   filepath.ToSlash(path),
				Line:   fset.Position(imp.Pos()).Line,
				Rule:   RuleCrossContextInternal,
				Detail: from + " imports " + value,
			})
		}
		return nil
	})
```

**Aceitação:** o teste novo passa; `archscan -root internal/contexts` continua a dar 0 findings
(a árvore real está limpa desde a Tarefa 2); o fixture continua a acusar. Zero de um instrumento
que **não** consegue demonstrar uma apanha significa "não medido" — por isso o fixture fica.

Commit: `fix(arch): detector cross-context deixa de saltar ficheiros fora dos contextos`

---

### Tarefa 4 — Detector de vendor com escopo (Regra 2.3) e portão com raízes certas

**Causa (D-e):** `ScanVendorTokensSuffix` (`scan.go:203-239`) não filtra caminho nenhum. Acusa
`scan.go:34-35`, onde a própria lista de tokens é literal. Enquanto a lista estiver no scanner e
o scanner na raiz varrida, o detector **nunca pode ficar verde**, e um detector que não pode ficar
verde é ignorado — que é como um instrumento morre.

A Regra 2.3 do design diz: nome de vendor não aparece **fora de `adapters/`**, e a regra governa
`contexts/`, `kernel/` e `composition/`. O escopo tem de estar no detector, não no comentário.

**RED.** Prova primeiro que o detector se acusa a si próprio, de `apps/server_core`:

```bash
export GOCACHE="$(pwd)/.gocache" && go run ./internal/arch/cmd/archscan -root internal/arch
```

Cola o output (deve listar `internal/arch/scan.go:34` e seguintes). Depois acrescenta a
`scan_test.go`:

```go
// TestVendorTokensIgnoresAdaptersAndOwnList pins the two halves of Regra 2.3:
// a vendor name inside adapters/ is the design, and the detector's own token
// list is not a violation of the rule it implements.
func TestVendorTokensIgnoresAdaptersAndOwnList(t *testing.T) {
	got, err := ScanVendorTokensSuffix("testdata", VendorTokens, ".go.txt")
	if err != nil {
		t.Fatalf("scan: %v", err)
	}
	for _, f := range got {
		if strings.Contains(f.File, "/adapters/") {
			t.Fatalf("test=TestVendorTokensIgnoresAdaptersAndOwnList: adapters/ must be exempt, got %s:%d", f.File, f.Line)
		}
	}
	// Positive control: the fixture outside adapters/ must still be caught.
	var caught bool
	for _, f := range got {
		if strings.HasSuffix(f.File, "vendor_in_context.go.txt") {
			caught = true
		}
	}
	if !caught {
		t.Fatalf("test=TestVendorTokensIgnoresAdaptersAndOwnList: positive control not caught; got %+v", got)
	}
}
```

Fixtures a criar em `apps/server_core/internal/arch/testdata/`:

- `adapters/vendor_ok.go.txt`:

```go
package mercadolivre

const provider = "mercado_livre"
```

- `vendor_in_context.go.txt`:

```go
package pricing

const provider = "mercado_livre"
```

**GREEN.** Em `scan.go`, acrescenta acima de `ScanVendorTokensSuffix`:

```go
// vendorRuleApplies reports whether Regra 2.3 governs this file.
//
// The rule is "a vendor name does not appear OUTSIDE adapters/", so adapters are
// exempt by definition. The scanner package is exempt too: its closed token list
// is the instrument, not a violation, and a detector that permanently accuses
// itself is a detector nobody can ever act on.
func vendorRuleApplies(path string) bool {
	p := filepath.ToSlash(path)
	if strings.Contains(p, "/adapters/") || strings.HasPrefix(p, "adapters/") {
		return false
	}
	if strings.Contains(p, "/internal/arch/") || strings.HasPrefix(p, "internal/arch/") {
		return false
	}
	return true
}
```

E como primeira instrução do visitor de `ScanVendorTokensSuffix`:

```go
		if !vendorRuleApplies(path) {
			return nil
		}
```

**GREEN — `scripts/arch-gate.sh`.** Duas correcções, ambas medidas:

1. O passo `archscan` varre `-root internal`, o que inclui `internal/modules/` — a árvore legada
   que a constraint 1 proíbe tocar. Um portão que reporta sobre código que ninguém pode consertar
   é ruído. Restringe às raízes que a doutrina governa.
2. O passo "no float in the kernel" faz `grep -rn 'float64\|float32'` e dispara em **comentários**
   (dívida D-32). Usa o detector AST que já existe em vez do grep.

Substitui os passos `architecture detectors` e `no float in the kernel` (linhas 29-44) por:

```bash
step "architecture detectors"
for root in internal/kernel internal/contexts internal/adapters internal/composition; do
  if (cd "$SERVER" && go run ./internal/arch/cmd/archscan -root "$root"); then
    echo "archscan $root: zero findings"
  else
    fail=1
  fi
done
echo "note: internal/modules is the legacy tree and is deliberately NOT scanned"
```

O detector de float (`ScanFloatInContracts`, `scan.go:167`) já é AST e já corre dentro do
`archscan`; o `grep` era uma segunda cópia pior do mesmo instrumento e sai.

**Aceitação:** `bash scripts/arch-gate.sh` → `ARCH GATE: PASS`. Os testes de fixture continuam a
acusar (controlo positivo vivo).

Commit: `fix(arch): Regra 2.3 ganha escopo; portao deixa de varrer legado e de grepar comentarios`

---

### Tarefa 5 — Detector de descarte do booleano de `Fact.Value()` (Regra 4.2)

**Causa (D-g):** `fact/knowledge.go:117-120` diz textualmente que *"The vet-level check for that
pattern lives in internal/arch"*. Não vive. E o único sítio que faz exactamente o padrão proibido
é `repository.go:360`: `desc, _ := p.Description().Value()`.

**RED.** Fixture `apps/server_core/internal/arch/testdata/discarded_value.go.txt`:

```go
package edge

func summarise(f Fact) string {
	v, _ := f.Value()
	return v
}
```

Teste em `scan_test.go`:

```go
// TestFactValueDiscardIsCaught is the level-2 instrument for Regra 4.2. The
// package that exists to stop "unknown becomes zero" shipped exactly one call
// site doing `v, _ := f.Value()`, and nothing looked.
func TestFactValueDiscardIsCaught(t *testing.T) {
	got, err := ScanFactValueDiscardSuffix("testdata", ".go.txt")
	if err != nil {
		t.Fatalf("scan: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("test=TestFactValueDiscardIsCaught: want exactly 1 finding, got %d: %+v", len(got), got)
	}
	if !strings.HasSuffix(got[0].File, "discarded_value.go.txt") || got[0].Rule != RuleFactValueDiscard {
		t.Fatalf("finding = %+v", got[0])
	}
}
```

Corre; deve falhar a compilar por `ScanFactValueDiscardSuffix` não existir. Cola o output.

**GREEN.** Em `scan.go`, acrescenta à lista de constantes (`:25-28`):

```go
	RuleFactValueDiscard = "facts/value-discarded"
```

E o detector:

```go
// ScanFactValueDiscardSuffix reports every `v, _ := something.Value()`.
//
// The blank is the whole defect. Fact.Value returns (T, bool) precisely so that
// an Unknown cannot be read as a value; discarding the bool hands back the zero
// value of T, which is the mistake the fact package exists to prevent. This is
// a syntactic detector on purpose: it cannot know the receiver's type without a
// type checker, so it reports every two-result .Value() call whose second result
// is discarded. A call site that legitimately does that renames its method.
func ScanFactValueDiscardSuffix(root, suffix string) (Findings, error) {
	var out Findings
	err := walk(root, suffix, func(path string, fset *token.FileSet, file *ast.File) error {
		ast.Inspect(file, func(n ast.Node) bool {
			assign, ok := n.(*ast.AssignStmt)
			if !ok || len(assign.Lhs) != 2 || len(assign.Rhs) != 1 {
				return true
			}
			blank, ok := assign.Lhs[1].(*ast.Ident)
			if !ok || blank.Name != "_" {
				return true
			}
			call, ok := assign.Rhs[0].(*ast.CallExpr)
			if !ok {
				return true
			}
			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok || sel.Sel.Name != "Value" {
				return true
			}
			out = append(out, Finding{
				File:   filepath.ToSlash(path),
				Line:   fset.Position(assign.Pos()).Line,
				Rule:   RuleFactValueDiscard,
				Detail: "the bool from .Value() is discarded: unknown would read as the zero value",
			})
			return true
		})
		return nil
	})
	out.sortInPlace()
	return out, err
}

// ScanFactValueDiscard scans real Go files.
func ScanFactValueDiscard(root string) (Findings, error) {
	return ScanFactValueDiscardSuffix(root, ".go")
}
```

Regista-o em `internal/arch/cmd/archscan/main.go`, na lista de `:18-24`:

```go
		arch.ScanFactValueDiscard,
```

**GREEN — o sítio real.** `internal/contexts/catalog/port/reader.go`, a struct `Summary` ganha o
estado, porque um consumidor que recebe `Description: ""` não consegue distinguir "sem nome" de
"não sabemos":

```go
// Summary answers "what is this product?", flattened for a consumer that has no
// business knowing Catalog's internal model.
//
// DescriptionState is not decoration. A consumer handed Description == "" cannot
// tell "the source says it has no name" from "we never learned it", and that is
// the difference that decides whether a screen shows a blank or an alert.
type Summary struct {
	ProductID        string
	Description      string
	DescriptionState string
	Identifiers      []contracts.Identifier
	Version          int
}
```

E `internal/contexts/catalog/internal/postgres/repository.go`, `summarise` (`:359-367`):

```go
// summarise flattens an aggregate for a consumer. The knowledge state travels
// with the value: this is the ONLY place the empty string is allowed to stand
// for an unknown description, and it is allowed only because DescriptionState
// says so beside it.
func summarise(p domain.Product) port.Summary {
	desc, known := p.Description().Value()
	if !known {
		desc = ""
	}
	return port.Summary{
		ProductID:        p.ID().String(),
		Description:      desc,
		DescriptionState: p.Description().State().String(),
		Identifiers:      p.Identifiers(),
		Version:          p.Version(),
	}
}
```

**Aceitação:** `archscan` sobre as 4 raízes → 0 findings, **e** o teste do fixture continua a
acusar 1. Zero sem controlo positivo é "não medido".

Commit: `feat(arch): Regra 4.2 vira nivel 2; summary carrega o estado de conhecimento`

---

### Tarefa 6 — Proveniência honesta na rehidratação e observação escrita

Dois defeitos no mesmo caminho (D-f), ambos em `catalog/internal/postgres/repository.go`.

**Defeito 1 — evidência fabricada.** `repository.go:203` monta
`provenance.NewEvidence("catalog", "product", productID, recordedAt, hash)`. Nada disto foi
observado: o sistema não é "catalog", a chave externa não é o `product_id` nosso, e `recordedAt`
é o `updated_at` da nossa linha, não o instante em que a fonte viu o facto. A tabela
`catalog.source_observations` (`0097_catalog_context.sql:65-78`) existe **precisamente** para
guardar isto e nunca recebe uma linha.

**Defeito 2 — chaves de fonte perdidas.** `repository.go:219-226` tenta juntar as chaves extra
chamando `p.Apply(extra)` com uma observação cujo hash é o mesmo. `domain/product.go:89-91`
devolve `DispositionIdempotent` nesse caso e **não junta a chave**. Um produto com duas chaves de
fonte rehidrata com uma.

**RED.** Teste de integração em `apps/server_core/tests/integration/` (tag `//go:build
integration` nas primeiras 5 linhas do ficheiro — a lane só descobre assim):

```go
//go:build integration

package integration

// TestRehydrationKeepsEverySourceKeyAndRealEvidence proves both halves of the
// defect at once: a product reached by two source keys must come back with two,
// and its evidence must name the system that actually observed it.
func TestRehydrationKeepsEverySourceKeyAndRealEvidence(t *testing.T) {
	// ... arrange: ingest the same product under two source keys (sankhya and a
	// second instance), then read it back through the store.
	// assert len(SourceKeys()) == 2
	// assert Description().Evidence().System() == "sankhya"
}
```

Corre a lane de integração da raiz do repo e cola o output:

```bash
npm run harness:integration
```

Espera-se `want 2 source keys, got 1` e `system = "catalog", want "sankhya"`.

**GREEN — `internal/contexts/catalog/internal/domain/product.go`.** O agregado passa a guardar a
evidência, não só o hash, e ganha um reconstituidor de chaves:

```go
type Product struct {
	id           ProductID
	tenant       tenant.ID
	version      int
	description  fact.Fact[string]
	identifiers  []contracts.Identifier
	sourceKeys   []contracts.SourceProductKey
	lastHash     string
	lastEvidence provenance.Evidence
}
```

`NewProduct` passa a preencher `lastEvidence: o.Evidence`, e `Apply` idem no `next`. Acrescenta:

```go
// LastEvidence returns how we last learned about this product. The repository
// writes it to catalog.source_observations, which is what makes rehydration able
// to tell the truth instead of naming itself as the source.
func (p Product) LastEvidence() provenance.Evidence { return p.lastEvidence }

// ReconstituteSourceKeys restores the full set of addresses a persisted product
// answers to. It exists because the alternative — replaying Apply once per extra
// key — is a no-op: Apply sees the same payload hash and reports Idempotent
// without merging anything, so the second key was silently dropped.
func ReconstituteSourceKeys(p Product, keys []contracts.SourceProductKey) Product {
	p.sourceKeys = append([]contracts.SourceProductKey(nil), keys...)
	return p
}
```

**GREEN — `repository.go`, escrita.** No fim de `writeProduct` (depois do loop de
`source_product_keys`, antes do `return nil`):

```go
	e := p.LastEvidence()
	if e.IsZero() {
		return fmt.Errorf("catalog/postgres: product %s has no evidence to record", p.ID())
	}
	if _, err := tx.Exec(ctx, `
		INSERT INTO catalog.source_observations
			(tenant_id, product_id, payload_hash, source_system, object_kind,
			 external_key, observed_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
		ON CONFLICT (tenant_id, product_id, payload_hash) DO NOTHING`,
		p.Tenant().String(), p.ID().String(), e.PayloadHash(), e.System(),
		e.ObjectKind(), e.ExternalKey(), e.ObservedAt().UTC()); err != nil {
		return fmt.Errorf("catalog/postgres: record observation: %w", err)
	}
	return nil
```

**GREEN — `repository.go`, leitura.** Substitui `:203-226` por:

```go
	e, err := loadLatestEvidence(ctx, tx, t, productID, hash)
	if err != nil {
		return domain.Product{}, false, err
	}
	desc, err := rehydrateFact(state, value, reason, e)
	if err != nil {
		return domain.Product{}, false, err
	}

	obs := contracts.ProductObservation{
		Key: keys[0], Description: desc, Identifiers: ident, Evidence: e,
	}
	p, err := domain.NewProduct(id, obs)
	if err != nil {
		return domain.Product{}, false, err
	}
	// Every key, restored directly. Replaying Apply here would report Idempotent
	// on the identical payload hash and drop every key after the first.
	p = domain.ReconstituteSourceKeys(p, keys)
	return domain.ReconstituteVersion(p, version, hash), true, nil
```

E a função nova:

```go
// loadLatestEvidence returns how we actually learned this product's current
// version. It reads the observation that produced the stored payload hash, so
// the rehydrated Fact carries the source system that observed it — not this
// package's own name, which is what it used to fabricate.
func loadLatestEvidence(ctx context.Context, tx pgx.Tx, t tenant.ID, productID, hash string) (provenance.Evidence, error) {
	var system, objectKind, externalKey string
	var observedAt time.Time
	row := tx.QueryRow(ctx, `
		SELECT source_system, object_kind, external_key, observed_at
		  FROM catalog.source_observations
		 WHERE tenant_id = $1 AND product_id = $2 AND payload_hash = $3`,
		t.String(), productID, hash)
	switch err := row.Scan(&system, &objectKind, &externalKey, &observedAt); {
	case errors.Is(err, pgx.ErrNoRows):
		// Not a missing optional. A product whose current version has no recorded
		// observation cannot say where it came from, and inventing one here is the
		// exact fabrication this change removes.
		return provenance.Evidence{}, fmt.Errorf(
			"catalog/postgres: product %s version hash %s has no source observation", productID, hash)
	case err != nil:
		return provenance.Evidence{}, fmt.Errorf("catalog/postgres: load observation: %w", err)
	}
	return provenance.NewEvidence(system, objectKind, externalKey, observedAt.UTC(), hash)
}
```

**Aceitação:** o teste de integração passa; `SELECT count(*) FROM catalog.source_observations`
> 0 na base da lane.

Commit: `fix(catalog): rehidratacao para de fabricar proveniencia e de perder chaves de fonte`

---

### Tarefa 7 — Porta de feed com cursor opaco (o adaptador sobrevive a um segundo provedor)

**Causa:** `catalogfeed.Feed.Page(ctx, t, after int64, limit int)`
(`catalogfeed/mapper.go:59`) leva o CODPROD do Sankhya **na assinatura**. Um segundo ERP cujo
cursor não seja um inteiro crescente não a pode implementar. O precedente já existe e já dói:
`modules/internal_read/ports/catalog_page.go:28` declara `Cursor{InternalProductID int64}` — o
mesmo vazamento, uma árvore acima. Aqui não se repete.

**RED.** `apps/server_core/internal/contexts/catalog/port/feed_test.go`:

```go
// TestFeedPortHidesTheSourceKeyShape is the whole point of the port: the cursor
// is a token, not an ERP row id, so an adapter whose paging key is a string or a
// timestamp can implement it. The legacy readports.Cursor exposes
// InternalProductID and cannot.
func TestFeedPortHidesTheSourceKeyShape(t *testing.T) {
	var _ ProductFeed = stubFeed{}
	c := NewCursor("opaque-token")
	if c.Token() != "opaque-token" || c.IsStart() {
		t.Fatalf("cursor = %+v", c)
	}
	if !(Cursor{}).IsStart() {
		t.Fatalf("the zero cursor must mean start")
	}
}
```

Corre; falha a compilar (`ProductFeed` não existe). Cola o output.

**GREEN — `apps/server_core/internal/contexts/catalog/port/feed.go`**, ficheiro completo:

```go
// Package port carries what other contexts may ASK Catalog, and what Catalog
// asks of a source. Both directions are questions, never tables.
package port

import (
	"context"

	"marketplace-central/apps/server_core/internal/contexts/catalog/contracts"
	"marketplace-central/apps/server_core/internal/kernel/tenant"
)

// Cursor is a source's position in its own feed, and Catalog cannot read it.
//
// The token is opaque on purpose. The legacy read port declared
// Cursor{InternalProductID int64} and thereby wrote one ERP's row-id shape into
// the contract: a source paging by string key, by timestamp, or by a provider's
// scroll id could not implement it without lying. Here the source encodes
// whatever it needs and Catalog hands it back untouched.
type Cursor struct{ token string }

// NewCursor builds a cursor from a source-defined token.
func NewCursor(token string) Cursor { return Cursor{token: token} }

// Token returns the source's own position marker.
func (c Cursor) Token() string { return c.token }

// IsStart reports whether this is the beginning of the feed. The zero cursor is
// the start, so a caller that forgets to seed one reads from the beginning
// rather than from an undefined place.
func (c Cursor) IsStart() bool { return c.token == "" }

// Page is one batch of observations plus where to continue.
type Page struct {
	Observations []contracts.ProductObservation
	Next         Cursor
	// Done means the source has no more rows. It is explicit rather than
	// inferred from an empty page, because an empty page in the middle of a
	// filtered feed is legal and must not stop the walk.
	Done bool
}

// ProductFeed is a source of product observations. Catalog asks; the adapter
// decides how to page.
type ProductFeed interface {
	NextPage(ctx context.Context, t tenant.ID, after Cursor, limit int) (Page, error)
}
```

**GREEN — `internal/adapters/erp/sankhyaoracle/catalogfeed/mapper.go`.** Substitui `Page`
(`:59-80`) por `NextPage`, que traduz o token:

```go
// NextPage implements port.ProductFeed. The Sankhya cursor is CODPROD, and that
// fact stops at this function: the token crossing the port is a string, so the
// day a second ERP pages by something else, nothing in Catalog changes.
func (f Feed) NextPage(ctx context.Context, t tenant.ID, after port.Cursor, limit int) (port.Page, error) {
	var afterCodprod int64
	if !after.IsStart() {
		parsed, err := strconv.ParseInt(after.Token(), 10, 64)
		if err != nil {
			return port.Page{}, fmt.Errorf("catalogfeed: cursor %q is not a CODPROD: %w", after.Token(), err)
		}
		afterCodprod = parsed
	}

	rows, err := oracle.FetchActiveProducts(ctx, f.db, afterCodprod, limit, f.now())
	if err != nil {
		return port.Page{}, err
	}
	out := make([]contracts.ProductObservation, 0, len(rows))
	var last int64
	for _, r := range rows {
		obs, err := MapProduct(t, f.instance, r)
		if err != nil {
			return port.Page{}, err
		}
		out = append(out, obs)
		last = r.Codprod
	}
	page := port.Page{Observations: out, Done: len(rows) < limit}
	if !page.Done {
		page.Next = port.NewCursor(strconv.FormatInt(last, 10))
	}
	return page, nil
}

var _ port.ProductFeed = Feed{}
```

**GREEN — `internal/adapters/erp/sankhyaoracle/sankhyaoracle.go`.** O `Bundle` passa a ser
tipado pela porta do consumidor, que é a Regra 2.2-a aplicada ao adaptador:

```go
// Bundle is everything this adapter offers, typed by the ports its consumers
// own. A field typed by catalogfeed.Feed would publish this adapter's concrete
// type to every caller; a field typed by port.ProductFeed publishes only the
// question Catalog asks.
type Bundle struct {
	Instance    string
	CatalogFeed port.ProductFeed
}
```

**Aceitação:** `go build ./...` verde; `TestFeedPortHidesTheSourceKeyShape` passa; o teste do
mapper existente continua a passar (ajusta-o para `NextPage`, **não o apagues**).

Commit: `feat(catalog): porta de feed com cursor opaco; bundle tipado por porta`

---

### Tarefa 8 — Sítio de composição: o ingest corre e escreve linhas

Sem isto o plano fecha com as lanes verdes e as 4 tabelas a zero — que é exactamente o estado de
hoje. Esta tarefa é a que produz o observável.

**Escopo honesto.** O agendamento **não** entra aqui: `syncapp.Scheduler` já tem job `products`
legado e a lista de entidades em `sync/domain/entity.go` é fechada; registar um segundo colidiria
em `sync_state`, e mexer em `internal/modules/sync` viola a constraint 1. O sítio de composição
deste plano é um comando de operador; o agendamento é decisão do próximo plano e entra como
dívida nomeada na Tarefa 12.

**GREEN — `apps/server_core/internal/composition/catalog_ingest.go`**, ficheiro completo:

```go
package composition

import (
	"context"
	"fmt"

	"marketplace-central/apps/server_core/internal/contexts/catalog"
	"marketplace-central/apps/server_core/internal/contexts/catalog/contracts"
	"marketplace-central/apps/server_core/internal/contexts/catalog/port"
	"marketplace-central/apps/server_core/internal/kernel/tenant"
)

// IngestReport is what one full walk of a feed did. Every disposition is counted
// separately because "nothing changed" and "nothing was read" are different
// outcomes and a single total cannot tell them apart.
type IngestReport struct {
	Pages      int
	Observed   int
	Created    int
	Changed    int
	Idempotent int
	Conflicts  []contracts.Identifier
}

// RunCatalogIngest walks a feed to exhaustion and folds every observation into
// Catalog. It is the production path: the composition root owns it, the context
// owns the decision, the adapter owns the SQL.
func RunCatalogIngest(ctx context.Context, module *catalog.Module, feed port.ProductFeed, t tenant.ID, pageSize int) (IngestReport, error) {
	if pageSize <= 0 {
		return IngestReport{}, fmt.Errorf("composition: page size must be positive, got %d", pageSize)
	}
	var report IngestReport
	cursor := port.Cursor{}
	for {
		page, err := feed.NextPage(ctx, t, cursor, pageSize)
		if err != nil {
			return report, fmt.Errorf("composition: read catalog feed page %d: %w", report.Pages+1, err)
		}
		report.Pages++
		for _, obs := range page.Observations {
			result, err := module.IngestProduct(ctx, obs)
			if err != nil {
				return report, fmt.Errorf("composition: ingest %s: %w", obs.Key, err)
			}
			report.Observed++
			switch result.Disposition {
			case contracts.DispositionCreated:
				report.Created++
			case contracts.DispositionChanged:
				report.Changed++
			case contracts.DispositionIdempotent:
				report.Idempotent++
			}
			report.Conflicts = append(report.Conflicts, result.DuplicateIdentifiers...)
		}
		if page.Done {
			return report, nil
		}
		if page.Next.IsStart() {
			// A source that is not done must advance. Returning the start cursor
			// again would loop forever, and a job that never ends looks exactly
			// like a job that is working.
			return report, fmt.Errorf("composition: feed reported more pages but did not advance the cursor")
		}
		cursor = page.Next
	}
}
```

**GREEN — `apps/server_core/cmd/catalogingest/main.go`.** Comando de operador. Lê a configuração
pelo mesmo caminho que o servidor (`internal/platform/config`), regista o driver Oracle (é a
raiz de composição que o faz — `sankhyaoracle.go:5-7` diz isso explicitamente), monta a fatia com
`WireCatalog`, corre `RunCatalogIngest` e imprime o relatório. Sai com código != 0 em erro, com o
erro no stderr. Nenhuma credencial é impressa.

**RED (integração), antes do comando.** Teste em `apps/server_core/tests/integration/` com
`//go:build integration` nas primeiras 5 linhas, usando uma dobra de `port.ProductFeed` de duas
páginas contra o Postgres real da lane:

```go
// TestCatalogIngestWritesRows is the acceptance this whole foundation lacked: the
// slice had green lanes and four empty tables. A stub feed is legal HERE because
// what is under test is the composition path and the writer, not Oracle.
func TestCatalogIngestWritesRows(t *testing.T) {
	// arrange: pool against the lane database, catalog.New(pool), a two-page stub feed
	// act: RunCatalogIngest(ctx, module, feed, tenantID, 2)
	// assert report.Created == 3 && report.Pages == 2
	// assert SELECT count(*) FROM catalog.products WHERE tenant_id = $1  -> 3
	// assert SELECT count(*) FROM catalog.source_observations ...        -> 3
	// act again with the same feed
	// assert report.Idempotent == 3 && report.Created == 0  (re-polling is free)
	// assert count(*) unchanged                              (idempotence at the row, not at the counter)
}
```

Corre `npm run harness:integration` e cola o output vermelho.

**Aceitação — o observável.** Duas provas, instrumentos independentes:

1. Lane: o teste acima passa, com `count(*)` da base a bater com o relatório.
2. Dev stack: com Oracle acessível, correr o comando e depois, na base:

```bash
docker exec marketplace-central-postgres-1 psql -U marketplace -d marketplace -c "SELECT count(*) FROM catalog.products;"
```

deve devolver > 0. **Nunca** `n_live_tup` nem `reltuples` — foram medidos a mentir neste repo.

**Se Oracle não estiver acessível** no ambiente do executor: a prova 1 fecha a tarefa, a prova 2
vai para `.mnfs/HARNESS-DEBTS.md` com o output verbatim da falha de ligação. Não inventes
contorno, não mockes Oracle na lane.

Commit: `feat(catalog): ingest tem chamador de producao e escreve linhas`

---

### Tarefa 9 — Aritmética de `Fact` (Regra 4.4, na forma que Go permite)

**Causa:** `kernel/fact/knowledge.go` tem **zero** métodos de aritmética. O design (Regra 4.4)
mandava que *"a aritmética vive em métodos de `Fact[T]`"*. **Isso é impossível em Go**: um método
não pode declarar parâmetros de tipo próprios, e somar `Fact[Money]` com `Fact[Money]` para dar
`Fact[Money]` já é possível, mas `Fact[A] × Fact[B] → Fact[C]` não. Isto não foi negligência do
implementador — foi uma promessa que a linguagem não cumpre. A emenda ao design é a Tarefa 11.

Sem isto, o primeiro cálculo que consumir dois `Fact` vai chamar `MustValue()` nos dois e o
kernel volta a ser doutrina.

**RED.** `apps/server_core/internal/kernel/fact/combine_test.go`:

```go
// TestCombine2PropagatesTheWorstState pins the rule that makes the kernel worth
// having: a calculation with one unknown input has an unknown output. Solving a
// target-margin price with one unknown component U gives P = 178.69 + 2.02*U —
// an ignored cost does not err slightly, it errs by twice what was ignored.
func TestCombine2PropagatesTheWorstState(t *testing.T) {
	// known + known   -> known
	// known + estimated -> estimated (with both reasons)
	// known + unknown -> unknown, and fn is NEVER called
	// any + not_applicable -> not_applicable
}

// TestCombine2NeverCallsFnOnAnUnusableInput is the negative control: a fn that
// records it ran must not have run.
func TestCombine2NeverCallsFnOnAnUnusableInput(t *testing.T) { }
```

Corre e cola o vermelho.

**GREEN — `apps/server_core/internal/kernel/provenance/derived.go`:**

```go
package provenance

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
	"time"
)

// Derived builds evidence for a value nobody observed.
//
// NewEvidence cannot serve: it demands one system, one external key and one
// observation time, and a derived value has as many of each as it has inputs.
// Forcing it through NewEvidence would mean picking one input and silently
// discarding the rest — the provenance equivalent of an unknown becoming zero.
//
// ObservedAt is the OLDEST input's time, not the newest and not now: a derived
// number is exactly as fresh as its stalest ingredient.
func Derived(method string, from ...Evidence) (Evidence, error) {
	method = strings.TrimSpace(method)
	if method == "" {
		return Evidence{}, fmt.Errorf("%w: derivation method", ErrIncomplete)
	}
	if len(from) == 0 {
		return Evidence{}, fmt.Errorf("%w: derivation has no inputs", ErrIncomplete)
	}
	refs := make([]string, 0, len(from))
	oldest := time.Time{}
	for _, e := range from {
		if e.IsZero() {
			return Evidence{}, fmt.Errorf("%w: derivation input", ErrIncomplete)
		}
		refs = append(refs, e.Ref())
		if oldest.IsZero() || e.ObservedAt().Before(oldest) {
			oldest = e.ObservedAt()
		}
	}
	sort.Strings(refs)
	h := sha256.New()
	fmt.Fprintf(h, "method=%s\n", method)
	for _, r := range refs {
		fmt.Fprintf(h, "from=%s\n", r)
	}
	return NewEvidence("derived", method, strings.Join(refs, "+"), oldest,
		"sha256:"+hex.EncodeToString(h.Sum(nil)))
}
```

**GREEN — `apps/server_core/internal/kernel/fact/combine.go`:**

```go
package fact

import (
	"fmt"
	"strings"

	"marketplace-central/apps/server_core/internal/kernel/provenance"
)

// Map applies fn to a usable fact and propagates everything else untouched.
//
// This is a package-level function and not a method because a Go method cannot
// declare type parameters of its own: `func (f Fact[A]) Map[B any](...)` does
// not compile. The design said the arithmetic lives in methods of Fact[T]; the
// language says otherwise, and the language wins.
func Map[A, B any](f Fact[A], method string, fn func(A) (B, error)) (Fact[B], error) {
	v, ok := f.Value()
	if !ok {
		return propagate[A, B](f, method)
	}
	out, err := fn(v)
	if err != nil {
		return Fact[B]{}, err
	}
	e, err := provenance.Derived(method, f.Evidence())
	if err != nil {
		return Fact[B]{}, err
	}
	if f.State() == Estimated {
		return NewEstimated(out, f.Reason(), e)
	}
	return NewKnown(out, e)
}

// Combine2 folds two facts into one.
//
// fn is NOT called unless BOTH inputs are usable. That is the whole guarantee:
// there is no code path in which a calculation sees the zero value of an unknown
// input, because the calculation is never reached.
//
// The output state is the worst of the two inputs, in this order:
// NotApplicable < Unknown < Estimated < Known. NotApplicable dominates because
// "this quantity does not exist here" makes the derived quantity not exist
// either, which is a stronger statement than "we failed to learn it".
func Combine2[A, B, C any](a Fact[A], b Fact[B], method string, fn func(A, B) (C, error)) (Fact[C], error) {
	av, aok := a.Value()
	bv, bok := b.Value()
	if !aok || !bok {
		return propagate2[A, B, C](a, b, method)
	}
	out, err := fn(av, bv)
	if err != nil {
		return Fact[C]{}, err
	}
	e, err := provenance.Derived(method, a.Evidence(), b.Evidence())
	if err != nil {
		return Fact[C]{}, err
	}
	if a.State() == Estimated || b.State() == Estimated {
		return NewEstimated(out, joinReasons(method, a, b), e)
	}
	return NewKnown(out, e)
}

func propagate[A, B any](f Fact[A], method string) (Fact[B], error) {
	e, err := provenance.Derived(method, f.Evidence())
	if err != nil {
		return Fact[B]{}, err
	}
	reason := fmt.Sprintf("%s: input is %s (%s)", method, f.State(), f.Reason())
	if f.State() == NotApplicable {
		return NewNotApplicable[B](reason, e)
	}
	return NewUnknown[B](reason, e)
}

func propagate2[A, B, C any](a Fact[A], b Fact[B], method string) (Fact[C], error) {
	e, err := provenance.Derived(method, a.Evidence(), b.Evidence())
	if err != nil {
		return Fact[C]{}, err
	}
	reason := joinReasons(method, a, b)
	if a.State() == NotApplicable || b.State() == NotApplicable {
		return NewNotApplicable[C](reason, e)
	}
	return NewUnknown[C](reason, e)
}

// joinReasons keeps BOTH inputs' reasons. An operator told only the first reason
// fixes one source and sees the value stay unknown with no idea why.
func joinReasons[A, B any](method string, a Fact[A], b Fact[B]) string {
	parts := []string{method + ":"}
	for _, s := range []struct {
		state  Knowledge
		reason string
	}{{a.State(), a.Reason()}, {b.State(), b.Reason()}} {
		if s.state != Known {
			parts = append(parts, fmt.Sprintf("%s (%s)", s.state, s.reason))
		}
	}
	return strings.Join(parts, " ")
}
```

**Aceitação:** os testes acima passam, incluindo o controlo negativo (a `fn` que regista execução
não corre). `go test ./internal/kernel/...` verde.

Commit: `feat(kernel): Fact ganha aritmetica que propaga o estado; proveniencia derivada`

---

### Tarefa 10 — RLS deixa de ser decorativo: papel de menor privilégio

**Causa medida:** o papel `marketplace` tem `rolsuper=t`, `rolbypassrls=t` e **é dono** das quatro
tabelas. `FORCE ROW LEVEL SECURITY` na 0097 não o alcança. Zero `CREATE ROLE`/`GRANT` em 84
migrações. Logo o isolamento de tenant do contexto novo é, hoje, uma alegação sem instrumento.

**RED.** Teste de integração que se liga **como o papel novo** e prova o isolamento:

```go
//go:build integration

// TestCatalogRLSBlocksTheOtherTenant is the only proof that FORCE RLS is not
// decorative. It must connect as the least-privilege role: the same query run as
// the owning superuser returns every row and would pass while proving nothing.
func TestCatalogRLSBlocksTheOtherTenant(t *testing.T) {
	// arrange: as the writer, insert one product for tenant A and one for tenant B
	// act: connect as mpc_app, set_config('app.tenant_id','A',true), select count(*)
	// assert count == 1
	// act: set_config('app.tenant_id','B',true) in a NEW transaction, select count(*)
	// assert count == 1, and the product_id differs from A's
	// negative control: with NO app.tenant_id set, count == 0 (fails closed)
}
```

Corre; falha porque o papel não existe. Cola o output.

**GREEN — `apps/server_core/migrations/0098_catalog_app_role.sql`:**

```sql
-- The catalog schema's RLS was FORCEd in 0097 and was still decorative: the
-- application connects as the table owner, which is also a superuser with
-- rolbypassrls. Every policy was evaluated against a role that skips policies.
--
-- This creates the role the application is meant to use. It is NOSUPERUSER and
-- NOBYPASSRLS explicitly rather than by default, because the defect being fixed
-- is precisely an attribute nobody looked at.
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'mpc_app') THEN
        CREATE ROLE mpc_app NOLOGIN NOSUPERUSER NOBYPASSRLS NOCREATEDB NOCREATEROLE;
    END IF;
END
$$;

ALTER ROLE mpc_app NOSUPERUSER NOBYPASSRLS;

GRANT USAGE ON SCHEMA catalog TO mpc_app;
GRANT SELECT, INSERT, UPDATE, DELETE ON catalog.products            TO mpc_app;
GRANT SELECT, INSERT, UPDATE, DELETE ON catalog.product_identifiers TO mpc_app;
GRANT SELECT, INSERT, UPDATE, DELETE ON catalog.source_product_keys TO mpc_app;
GRANT SELECT, INSERT, UPDATE, DELETE ON catalog.source_observations TO mpc_app;

-- No GRANT on any other schema. A role that can read the legacy tables would
-- make the boundary this context exists to draw invisible again.
```

**Bump do fixture:** `internal/platform/migrate/runner_test.go` passa de 84 para 85 (mede outra
vez com `ls`, não confies neste número).

**Escopo honesto:** apontar a DSN de produção a `mpc_app` é mudança de deploy, fora deste plano.
Entra como dívida nomeada na Tarefa 12. O que esta tarefa entrega é o papel, os GRANTs e a
**prova** de que a política funciona quando a ligação não a contorna.

Commit: `feat(catalog): papel de menor privilegio torna o RLS do contexto provavel`

---

### Tarefa 11 — Governança alcança a árvore nova

**Causa:** `contracts/governance/modules.json` tem **zero** menções a `contexts`, `kernel` ou
`adapters/erp`. Todos os invariantes de módulo têm `scope_paths`
`apps/server_core/internal/modules`. A árvore nova não tem cobertura, dependências declaradas nem
camadas verificadas — nada a impede de regredir depois deste plano.

E `production-panic` tem escopo `apps/server_core` inteiro, portanto **alcança** três `panic`
novos sem excepção registada: `exact/decimal.go:84`, `exact/decimal.go:141`,
`fact/knowledge.go:134`.

**RED.** Da raiz do repo, num worktree limpo, com o SHA de 40 hex completo:

```bash
npm run harness:governance -- -BaseSha <sha-40-hex>
```

Cola o output. Espera-se `GOV_PRODUCTION_PANIC` nos três sítios.

**GREEN.**

1. `contracts/governance/modules.json` ganha entradas para `kernel`, `contexts/catalog` e
   `adapters/erp/sankhyaoracle`, cada uma com o seu `root`, `composition_required` e
   `dependencies` **medidos dos imports reais**, não presumidos. Regra da memória:
   *entry em `modules.json` entra pelo merge, nunca pré-merge* — o executor escreve a entrada na
   mesma tarefa, e a lane de governança só a verá verde depois do merge.
2. `contracts/governance/invariants.json` ganha três `temporary_exceptions` de
   `production-panic`, **cada uma com a sua razão e um `removal_owner`**:
   - `exact/decimal.go:84` — `MustParseDecimal` é o idioma de constante literal em teste e em
     inicialização de pacote; remoção = o construtor deixar de existir.
   - `exact/decimal.go:141` — `StringFixed` com escala negativa é erro de programação, não de
     dados; remoção = a assinatura passar a devolver erro.
   - `fact/knowledge.go:134` — `MustValue` é o contrato do pacote: só é chamável depois de
     `IsUsable` na mesma função. Remoção = o detector da Tarefa 5 aprender a exigir o `IsUsable`
     no mesmo bloco, e então esta entrada morre.

   Se o formato `exception_mode: "exact-occurrence"` exigir a linha exacta, mede-a antes de
   escrever — as linhas rodam.

**Aceitação:** a lane de governança corre sem `GOV_PRODUCTION_PANIC` não declarado.

Commit: `chore(governance): arvore de contextos entra no registo; panicos do kernel declarados`

---

### Tarefa 12 — ADRs: a decisão congelada 7 e as quatro emendas ao design

**Causa (conflito de ordem de verdade, classificado na Medição 4):** `ARCHITECTURE.md:39`
decisão congelada 7 diz que integrações de marketplace entram **só** pelo módulo `connectors`;
o design ratificado põe-nas em `adapters/marketplace/<vendor>`. `ARCHITECTURE.md:5`: decisões
congeladas não se rediscutem sem ADR explícito. A ordem de verdade põe `ARCHITECTURE.md` acima do
design, portanto o design **ainda não venceu** — falta o ADR.

**GREEN — `docs/architecture/decisions/033-integracoes-entram-por-adapters.md`.** ADR que:

- cita a decisão congelada 7 verbatim e diz o que a substitui;
- dá a razão medida: `connectors` é um módulo do tipo que este refactor está a desmontar, e o
  design de 2026-08-06 põe a fronteira no compilador (`internal/` por segmento de caminho) em vez
  de na convenção;
- declara a transição: `internal/modules/connectors/` continua a servir Mercado Livre enquanto
  não for migrado; `adapters/marketplace/<vendor>` é o destino; **nenhum código novo de
  marketplace entra em `connectors`**;
- actualiza `ARCHITECTURE.md`: a decisão 7 passa a citar o ADR-033.

**GREEN — `docs/architecture/decisions/034-fact-substitui-adr-017.md`.** ADR que declara
`017-unknown-is-never-zero` **substituído** por `internal/kernel/fact`: a regra deixa de ser
doutrina citada 1378 vezes e passa a ser um tipo que não deixa construir a combinação inválida.
ADR-017 fica marcado `Superseded by ADR-034` e o texto explica que o predicado sobrevive — muda o
nível de imposição, não a regra.

**GREEN — emendas ao design**, em
`docs/superpowers/specs/2026-08-06-protocolo-de-codigo-design.md`, cada uma com a medição que a
paga:

| Regra | Emenda | Porquê, medido |
|---|---|---|
| 4.1 (valor zero de `Knowledge`) | O spec dizia `Known = iota + 1`; **ratifica-se o código**, `Unknown = iota` (`fact/knowledge.go:28`) | `Fact[string]{}` compila em qualquer pacote — um literal de struct vazio não nomeia campos — portanto "não há literal de struct" é falso nas duas variantes. A guarda real é `Evidence().IsZero()`, e com `Unknown = iota` o zero é o estado seguro. |
| 4.4 (aritmética) | De "métodos de `Fact[T]`" para **funções genéricas de pacote** `Map`/`Combine2` | Go não permite parâmetros de tipo próprios num método. A regra como escrita é inimplementável. |
| 1.3 (instrumento cross-context) | O instrumento é o detector **sem** o salto de ficheiros fora de `contexts/` | O único sítio que violou a regra estava fora, e o detector era cego a ele (`scan.go:131-133`). |
| 2.3 (tokens de vendor) | O escopo (`adapters/` isento, `internal/arch/` isento) vive **no detector** | Sem filtro de caminho o detector acusa a sua própria lista e nunca fica verde (`scan.go:33-36` vs `:203`). |

Cada emenda leva uma linha no log de emendas do spec com data e o `file:line` que a pagou.

Commit: `docs(arch): ADR-033 e ADR-034; quatro emendas ao design com a medicao que as pagou`

---

### Tarefa 13 — Dívidas: fechar o que fechou, registar o que fica

Em `.mnfs/HARNESS-DEBTS.md`, **com a medição de cada uma** — nada entra como opinião.

**Fecham** (com o `file:line` e o comando que o prova):

- **D-28** — o detector de vendor acusava a própria lista. Resolvido por escopo no detector
  (Tarefa 4), **não** por mexer no teste da Tarefa 1 do plano anterior.
- **D-29 / D-30** — o build partido e o teste de integração impossível. Eram o compilador a
  recusar, correctamente, um construtor sem chamador legal. Resolvidos na Tarefa 2.
- **D-32** — o grep de float disparava em comentários. O grep sai; o detector AST fica (Tarefa 4).

**Ficam abertas, nomeadas, com a medição:**

- **Agendamento do ingest de catálogo.** `sync/domain/entity.go` tem lista fechada e já há job
  `products` legado; registar um segundo colide em `sync_state` e mexer em `internal/modules/`
  viola a constraint 1 deste plano. Hoje o ingest é comando de operador
  (`cmd/catalogingest`). Decisão pertence ao plano de migração de contextos.
- **DSN de produção a apontar para `mpc_app`.** O papel e os GRANTs existem (Tarefa 10) e o
  isolamento está provado na lane; a aplicação continua a ligar-se como `marketplace`
  (`rolsuper=t`, `rolbypassrls=t`), portanto em produção o RLS continua contornado. Medição:
  `SELECT rolname, rolsuper, rolbypassrls FROM pg_roles WHERE rolname='marketplace'`.
- **Visibilidade de falha em ecrã.** Uma falha do ingest hoje é um código de saída num terminal.
  Não há ecrã, não há `sync_health`, não há cartão de operador. Escopo do próximo plano.
- **`internal/modules/` ainda fora dos detectores.** O portão varre 4 raízes novas e declara a
  omissão em voz alta (`scripts/arch-gate.sh`). A árvore legada só é varrida quando for migrada.

**Aceitação:** `bash scripts/arch-gate.sh` → `ARCH GATE: PASS`; `go build ./...`, `go vet ./...`,
`npm run harness:unit`, `npm run harness:integration` verdes;
`git status --porcelain --untracked-files=all` vazio.

Commit: `docs(debts): fecho da fundacao — 3 dividas fechadas por medicao, 4 registadas`

---

## Ordem e dependências

```
1 (lane verde)
  └─ 2 (build verde) ──┬─ 3 (detector cross-context)
                       ├─ 4 (escopo vendor + portao)
                       ├─ 5 (detector Value + summary)   [depende de 4: mesmo ficheiro scan.go]
                       ├─ 6 (proveniencia honesta)
                       └─ 7 (porta de feed) ─ 8 (sitio de composicao)  [depende de 6 e 7]
9 (aritmetica de Fact)      — independente de 2..8
10 (papel RLS)              — depende de 6 (escreve observations) e re-bumpa o fixture de 1
11 (governanca)             — depois de 2 e 9 (os panics do kernel ja existem)
12 (ADRs e emendas)         — depois de 3,4,5,9 (as emendas citam o que ficou feito)
13 (dividas)                — ultima
```

Tarefas 3, 4 e 5 tocam todas em `internal/arch/scan.go` — **um dono de cada vez**, em série.
Tarefas 6 e 7 tocam em ficheiros disjuntos e podem correr em paralelo; a 8 espera pelas duas.

## Colisões

`git worktree list` e `git diff --name-only main...<branch>` de cada ramo em voo antes de
despachar. Nenhum ramo em voo à data desta medição (`main` limpo em `1b2ef2da`), portanto a
intersecção é vazia por medição e não por presunção.

Costuras exclusivas que este plano toca, cada uma com um só dono:
`apps/server_core/migrations` (Tarefa 10), `docs/architecture/decisions` (Tarefa 12),
`contracts/governance/` (Tarefa 11). `composition/root.go` **não** é tocado —
`catalog_ingest.go` e `catalog_wiring.go` são ficheiros próprios.
