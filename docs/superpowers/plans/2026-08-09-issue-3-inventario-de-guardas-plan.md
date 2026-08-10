# Issue #3 — inventário de guardas e fim da evidência que se auto-certifica

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development
> (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use
> checkbox (`- [ ]`) syntax for tracking.

**Goal:** Fechar o issue #3 re-escopado: (a) inventário de guardas com meta-check que exige
fixture EXECUTADA, (b) piso de contagem nos 2 lanes sem guarda de zero (`build`, `lint-go`),
(c) fixture `testdata/` própria para o detector ADR-023, (d) retirar `scripts/arch-gate.sh`.

**Architecture:** Tudo dentro do gate existente (`scripts/gate.ps1` + `scripts/harness/Gate.psm1`
+ `contracts/gate/`). Nenhum runtime de produto é tocado; o único Go novo é teste + `testdata/`.
O inventário é um registro JSON que INDEXA fixtures que já existem espalhadas — não as re-implementa.

**Tech Stack:** PowerShell 7 (gate), Go test (fixtures), JSON (`contracts/gate/`), GitHub Actions.

**Branch:** `issue-3-guard-inventory`. Aterrissagem: PR → check `required` → merge (primeiro PR
sob o ruleset `main-landing-path`, o que de quebra exercita o #4 no caminho real).

---

## Medição (mc-planning Fase 1 — tudo aberto nesta sessão, árvore `2854226`)

1. **O que já existe que faz parte disto?**
   - Zero-check por lane: 14 de 16 lanes têm (ex.: `scripts/gate.ps1:144-147` gofmt
     `discovered=0`; `:1085-1088` guarda global `ran=0`). Furos: `build` (`:187-202`, só exit
     code) e `lint-go` (`:486-498`, só `report=missing` + drift de linters).
   - Lane de integração já anti-vácua: `scripts/harness/Postgres.psm1:468-472`
     (`HPG_TEST_VACUOUS` quando `tests_run -eq 0 -or tests_passed -eq 0`), tokens
     `test=`/`package=` em `:252-253`, com fixtures em
     `scripts/tests/postgres-lifecycle.tests.ps1:151,162`.
   - Fixtures de guarda espalhadas: `internal/arch/testdata/` + `scan_test.go` (contagens
     exatas, `TestCountFilesSeesTheFixtureTree` exige n≥10); `governance-drift.tests.ps1`
     (fixtures positivas para 9 códigos GOV, medido por `grep -o 'GOV_[A-Z_]+'`);
     `gate-measure.tests.ps1` (entradas sintéticas vermelhas para os contadores);
     ci.yml `guard self-test` (5 cenários jq no job `required`).
2. **Onde vive o gap?** (i) Não existe registro que ligue guarda→fixture→execução — um guard
   pode morrer e nada nota (universo: `contracts/`, `scripts/`, zero hits para inventário).
   (ii) `Invoke-GateBuild` e `Invoke-GateLintGo` aceitam universo vazio. (iii) O detector
   ADR-023 (`internal/composition/module_boundary_arch_test.go`) só é vermelho porque a árvore
   legada tem 234 violações reais — sem `testdata/`, quando a migração zerar, ele pode virar
   inerte sem sinal. (iv) `GOV_MODULE_LAYER` e `GOV_POSTGRES_DRIVER` são os únicos 2 códigos
   sem fixture positiva (universo: `scripts/tests/governance-drift.tests.ps1`, grep acima).
3. **Quem consome os caminhos alterados?** `ci.yml` chama `./scripts/gate.ps1 -Lane <name>`
   por step (`:95-341`); `npm run gate`/`gate:full` local; lane `boundary` consome a SAÍDA do
   teste ADR-023 (`gate.ps1:339` roda `-run TestModuleBoundaryADR023 -v` e mede o texto) —
   **o formato de saída do teste é contrato do lane e não pode mudar**.
4. **Contrato existente?** Sem OpenAPI/SDK. Contratos do gate: `contracts/gate/baselines.json`
   (ratchet lint) e `contracts/gate/census-exempt.json`. `guards.json` entra ao lado.
5. **Estado vivo?** Runner verde nos PRs #17–#24; ruleset `main-landing-path` (id 20614327)
   ativo com check `required` obrigatório; branch local `issue-3-guard-inventory` a partir de
   `2854226`.
6. **O que prova quebrado hoje / consertado depois?** Hoje: `Test-Path contracts/gate/guards.json`
   = False; grep sem hits para os 2 códigos GOV no arquivo de testes; `gate.ps1:192-201` sem
   contagem. Depois: `./scripts/gate.ps1 -Lane guards` imprime `entries=N ... failures=0` com
   N≥15 E o controle negativo (entrada boga no inventário) produz FAIL nomeando o teste ausente.
7. **Custo:** lane `guards` = 2 execuções `go test -run` direcionadas (segundos) + greps. Zero
   chamadas de provider. CI: +1 step no job `verify-full`.
8. **O que falha silencioso às 3h?** Hoje: um guard morto = verde vácuo (é o issue). Depois:
   guard morto → lane `guards` FAIL → check `required` vermelho → merge bloqueado. Caminho
   falha→pixel: saída do step no PR.

## O que já existe (Fase 2 — por que não serve)

| Artefato próximo | Por que não fecha o gap |
|---|---|
| `census` (`gate.ps1:874+`) | Prova que ARQUIVO de teste está no conjunto executado por algum lane; não liga guarda→fixture NOMEADA, nem falha se o teste for renomeado/apagado dentro do arquivo |
| `selftest` (`gate.ps1:827-872`) | Executa arquivos pwsh inteiros com token PASS por ARQUIVO; não endereça teste individual |
| `gate-measure.tests.ps1` | Prova os CONTADORES com texto sintético; não prova detectores Go nem cobre inventário |
| `scan_test.go` fixtures | Provam archscan; nada garante que continuam EXECUTANDO se renomeados (só o census vê o arquivo) |

**Local vs global (Fase 3):** ≥4 famílias de fixture convergem sem índice — o conserto global é
o REGISTRO + meta-check, não uma 5ª família. Custo extra sobre o conserto local (só ADR-023
fixture): um JSON + ~80 linhas de lane. Tomado. `arch-gate.sh` é legado não-referenciado:
removido em task própria, nunca estendido (Fase 3 pergunta 4).

**Seams (Fase 4):** `scripts/gate.ps1` + `ci.yml` + `contracts/gate/` = harness control files
(seam exclusivo) — este branch é o único escritor (verificado: `git worktree list` = só main;
nenhum branch em voo). Sem migração, sem tenant, sem provider, sem OpenAPI/SDK.

---

### Task 1: Fixture `testdata/` para o detector ADR-023

**Files:**
- Modify: `apps/server_core/internal/composition/module_boundary_arch_test.go`
- Create: `apps/server_core/internal/composition/testdata/adr023/alpha/domain/uses_beta_domain.go`
- Create: `apps/server_core/internal/composition/testdata/adr023/alpha/application/uses_beta_ports.go`
- Create: `apps/server_core/internal/composition/testdata/adr023/gamma/adapters/uses_sourcekind_domain.go`

Restrição dura: a lane `boundary` (`gate.ps1:339`) mede o TEXTO de
`go test ./internal/composition/ -run TestModuleBoundaryADR023 -v`. O refactor não pode mudar
uma linha da saída desse teste (`boundary_files=%d`, `%d violation(s)`, blocos by origin/target).

- [ ] **Step 1: Escrever o teste que falha (compile-red)**

Acrescentar ao fim de `module_boundary_arch_test.go`:

```go
// The detector's red must not depend on the legacy tree still being dirty.
// When the migration empties internal/modules, this fixture is what keeps the
// instrument demonstrably alive (issue #3, mechanism (a)).
func TestModuleBoundaryDetectorFiresOnFixture(t *testing.T) {
	violations, filesParsed, err := scanModuleBoundary(filepath.Join("testdata", "adr023"))
	if err != nil {
		t.Fatalf("scanning fixture tree: %v", err)
	}
	if filesParsed != 3 {
		t.Fatalf("filesParsed = %d, want 3; the fixture tree is not what this test believes", filesParsed)
	}
	if len(violations) != 1 {
		t.Fatalf("violations = %d, want exactly 1 (alpha/domain -> beta/domain)", len(violations))
	}
	v := violations[0]
	if v.fromModule != "alpha" || v.fromLayer != "domain" || v.toModule != "beta" || v.toLayer != "domain" {
		t.Fatalf("violation = %+v, want alpha/domain -> beta/domain", v)
	}
	if v.line <= 0 {
		t.Fatalf("violation line = %d, want > 0", v.line)
	}
}
```

- [ ] **Step 2: Rodar e ver a falha nomeada**

De `apps/server_core`, com `GOCACHE` absoluto
(`$env:GOCACHE = (Join-Path (Get-Location) '.gocache')`):

```bash
go test ./internal/composition/ -run TestModuleBoundaryDetectorFiresOnFixture -count=1 -v
```

Esperado: FAIL de compilação — `undefined: scanModuleBoundary`.

- [ ] **Step 3: Extrair `scanModuleBoundary` sem mudar a saída do teste velho**

Em `module_boundary_arch_test.go`, mover o corpo do walk (linhas do `fset := token.NewFileSet()`
até o fechamento do `WalkDir`, hoje `:91-146`) para:

```go
// scanModuleBoundary walks a modules-shaped tree and returns every cross-module
// import that lands on a layer other than ports, plus the number of Go files
// parsed. Extracted so a fixture tree can prove the detector fires (issue #3);
// TestModuleBoundaryADR023 keeps identical output over ../modules.
func scanModuleBoundary(modulesDir string) ([]boundaryViolation, int, error) {
	fset := token.NewFileSet()
	var violations []boundaryViolation
	filesParsed := 0

	err := filepath.WalkDir(modulesDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}
		filesParsed++

		rel, err := filepath.Rel(modulesDir, path)
		if err != nil {
			return err
		}
		fromModule, fromLayer := moduleAndLayer(rel)

		file, err := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
		if err != nil {
			return err
		}

		for _, spec := range file.Imports {
			importPath, err := strconv.Unquote(spec.Path.Value)
			if err != nil {
				return err
			}
			toModule, toLayer, ok := importedModuleAndLayer(importPath)
			if !ok || toModule == fromModule {
				continue
			}
			if sharedCore[toModule] {
				continue
			}
			if toLayer == "ports" {
				continue
			}
			violations = append(violations, boundaryViolation{
				file:       filepath.ToSlash(filepath.Join("internal/modules", rel)),
				line:       fset.Position(spec.Pos()).Line,
				fromModule: fromModule,
				fromLayer:  fromLayer,
				toModule:   toModule,
				toLayer:    toLayer,
			})
		}
		return nil
	})
	return violations, filesParsed, err
}
```

`TestModuleBoundaryADR023` passa a começar assim (o resto — log, guarda de zero, agregações,
`t.Errorf` — fica byte-idêntico):

```go
func TestModuleBoundaryADR023(t *testing.T) {
	modulesDir := filepath.Join("..", "modules")
	violations, filesParsed, err := scanModuleBoundary(modulesDir)
	if err != nil {
		t.Fatalf("walking %s: %v", modulesDir, err)
	}
	// ... segue idêntico a partir de t.Logf("boundary_files=%d", filesParsed)
```

Nota: o caminho reportado nas violações da fixture carrega o prefixo `internal/modules` — é
cosmético dentro da fixture e o teste novo não o assevera.

- [ ] **Step 4: Criar a árvore de fixture (3 arquivos)**

`testdata/adr023/alpha/domain/uses_beta_domain.go`:

```go
// Package domain is an ADR-023 detector fixture: a domain layer importing a
// sibling module's domain is THE violation. Never compiled (testdata).
package domain

import _ "marketplace-central/apps/server_core/internal/modules/beta/domain"
```

`testdata/adr023/alpha/application/uses_beta_ports.go` (controle: legal, zero findings):

```go
// Package application is an ADR-023 detector fixture: ports is the one public
// surface, so this import is legal and must produce no finding.
package application

import _ "marketplace-central/apps/server_core/internal/modules/beta/ports"
```

`testdata/adr023/gamma/adapters/uses_sourcekind_domain.go` (controle: carve-out shared core):

```go
// Package adapters is an ADR-023 detector fixture: sourcekind is the shared
// core carve-out and must produce no finding at any layer.
package adapters

import _ "marketplace-central/apps/server_core/internal/modules/sourcekind/domain"
```

- [ ] **Step 5: Verde no teste novo + saída inalterada no velho**

```bash
go test ./internal/composition/ -run 'TestModuleBoundaryDetectorFiresOnFixture' -count=1 -v
```

Esperado: `--- PASS: TestModuleBoundaryDetectorFiresOnFixture`.

```bash
go test ./internal/composition/ -run TestModuleBoundaryADR023 -count=1 -v > after.txt; echo exit=$?
```

Comparar com a MESMA invocação capturada antes do refactor (rodar antes do Step 3 e guardar
`before.txt`): `boundary_files=` e a contagem `N violation(s)` idênticos. `diff before.txt
after.txt` limpo fora de timings. Apagar os dois txt (nunca commitá-los).

- [ ] **Step 6: Commit**

```bash
git add apps/server_core/internal/composition
git commit -m "test(arch): ADR-023 detector proves it fires on a committed fixture (#3)"
```

---

### Task 2: Piso de contagem no lane `build`

**Files:**
- Modify: `scripts/gate.ps1` (`Invoke-GateBuild`, hoje `:187-202`)

- [ ] **Step 1: Substituir `Invoke-GateBuild`**

```powershell
function Invoke-GateBuild {
  # HARNESS-PROFILE §11 defines the deterministic L0 lane as build AND vet. The
  # test lane is not a substitute: `go test` without an explicit -vet flag runs a
  # curated subset of vet's checks (`go help testflag`), and vet's checks over
  # non-test code are outside it entirely.
  #
  # `go list ./...` first: build and vet report only exit codes, so without an
  # independent package count this lane cannot tell "everything compiled" from
  # "the pattern matched nothing" (issue #3, mechanism (b)).
  $list = Invoke-GateTool -Name 'go-list' -FilePath 'go' -ArgumentList @('list', './...') -WorkingDirectory $serverCore -Quiet
  if ($list.ExitCode -ne 0) {
    return New-GateVerdict -Lane 'build' -Passed $false -Counts 'packages=unknown' -Reason "go list ./... exited $($list.ExitCode)"
  }
  $packages = @($list.Text -split "`r?`n" | Where-Object { $_ -match '\S' }).Count
  if ($packages -eq 0) {
    return New-GateVerdict -Lane 'build' -Passed $false -Counts 'packages=0' `
      -Reason 'go list ./... resolved zero packages; the universe is empty, not clean'
  }
  $build = Invoke-GateTool -Name 'go-build' -FilePath 'go' -ArgumentList @('build', './...') -WorkingDirectory $serverCore
  if ($build.ExitCode -ne 0) {
    return New-GateVerdict -Lane 'build' -Passed $false -Counts "packages=$packages build=fail" -Reason "go build ./... exited $($build.ExitCode)"
  }
  $vet = Invoke-GateTool -Name 'go-vet' -FilePath 'go' -ArgumentList @('vet', './...') -WorkingDirectory $serverCore
  if ($vet.ExitCode -ne 0) {
    return New-GateVerdict -Lane 'build' -Passed $false -Counts "packages=$packages build=ok vet=fail" -Reason "go vet ./... exited $($vet.ExitCode)"
  }
  Write-Host "packages=$packages build=ok vet=ok"
  return New-GateVerdict -Lane 'build' -Passed $true -Counts "packages=$packages build=ok vet=ok"
}
```

- [ ] **Step 2: Rodar o lane e conferir a contagem**

```bash
pwsh -NoProfile -File ./scripts/gate.ps1 -Lane build
```

Esperado: `packages=<N> build=ok vet=ok` com N na casa das dezenas (universo:
`apps/server_core/./...`), `gate: PASS`.

- [ ] **Step 3: Commit**

```bash
git add scripts/gate.ps1
git commit -m "gate(build): count the package universe; zero packages is a failure, not a pass (#3)"
```

---

### Task 3: Piso de contagem no lane `lint-go`

**Files:**
- Modify: `scripts/gate.ps1` (`Invoke-GateLintGo`, região `:475-499`)

- [ ] **Step 1: Inserir a mesma guarda de universo antes da execução do golangci-lint**

Imediatamente antes da linha `$reportPath = Join-Path $logDirectory 'golangci.json'`
(hoje `:475`), inserir:

```powershell
  # Same universe guard as the build lane: the ratchet compares totals, and a
  # total taken over zero analyzed packages would read as the ratchet's best day
  # ever (issue #3, mechanism (b)).
  $list = Invoke-GateTool -Name 'golangci-go-list' -FilePath 'go' -ArgumentList @('list', './...') -WorkingDirectory $serverCore -Quiet
  if ($list.ExitCode -ne 0) {
    return New-GateVerdict -Lane 'lint-go' -Passed $false -Counts 'packages=unknown' -Reason "go list ./... exited $($list.ExitCode)"
  }
  $packages = @($list.Text -split "`r?`n" | Where-Object { $_ -match '\S' }).Count
  if ($packages -eq 0) {
    return New-GateVerdict -Lane 'lint-go' -Passed $false -Counts 'packages=0' `
      -Reason 'go list ./... resolved zero packages; a lint verdict over the empty set is not a verdict'
  }
```

E acrescentar `packages=$packages ` ao início de cada string `-Counts` do lane que vem DEPOIS
do ponto de inserção — hoje `:487` (`report=missing`), `:496` (`total=`), `:503` (`enabled=`)
e o par ratchet/verde `:527`/`:530` (`$counts`, prefixar na atribuição de `$counts`). A de
`:451` (`baseline=missing`) fica como está: acontece antes da medição do universo.

- [ ] **Step 2: Rodar o lane**

```bash
pwsh -NoProfile -File ./scripts/gate.ps1 -Lane lint-go
```

Esperado: counts começando com `packages=<N>` e o desfecho ratchet habitual (PASS com total ≤
baseline).

- [ ] **Step 3: Commit**

```bash
git add scripts/gate.ps1
git commit -m "gate(lint-go): state the analyzed universe; zero packages fails the lane (#3)"
```

---

### Task 4: Fixtures positivas para `GOV_MODULE_LAYER` e `GOV_POSTGRES_DRIVER`

**Files:**
- Modify: `scripts/tests/governance-drift.tests.ps1` (inserir após o bloco do orphan context,
  hoje `:278-283`)

Estes são controles positivos de instrumentos já implementados
(`Policy.psm1:378-382` e `:412-417`): o mundo vermelho deles é um Policy.psm1 quebrado, não o
código atual — por isso não há arm vermelho executável aqui; a prova de valor é que os 2 únicos
códigos sem fixture passam a ter uma (universo: 11 códigos, 9 fixturados antes desta task).

- [ ] **Step 1: Inserir as duas fixtures**

```powershell
# GOV_MODULE_LAYER: a cross-module import that lands on adapters/ is coupling
# with an extra step, not translation (ADR-023 §2). Policy.psm1:378-383. The file
# sits in domain/, not application/, so GOV_APPLICATION_IMPORT stays out of the
# blast radius and the assertion isolates the code under test. The source module
# must be DECLARED (Policy.psm1:370 skips undeclared sources); orders is, at
# contracts/governance/modules.json:125.
$layerFixture = New-PositiveFixture; $fixtures.Add($layerFixture)
Write-FixtureFile $layerFixture 'apps/server_core/internal/modules/orders/domain/uses_catalog_adapters.go' "package domain`n`nimport _ `"marketplace-central/apps/server_core/internal/modules/catalog/adapters/postgres`"`n"
$layerResult = Test-GovernanceDrift -RepositoryRoot $layerFixture
Assert-True (-not $layerResult.Passed) 'cross-module adapters import produced no finding'
Assert-True ('GOV_MODULE_LAYER' -in @($layerResult.Violations.ErrorCode)) "expected GOV_MODULE_LAYER, got: $(@($layerResult.Violations.ErrorCode) -join ',')"

# GOV_POSTGRES_DRIVER: database/sql under adapters/postgres bypasses pgdb.
# Policy.psm1:414-420.
$driverFixture = New-PositiveFixture; $fixtures.Add($driverFixture)
Write-FixtureFile $driverFixture 'apps/server_core/internal/modules/orders/adapters/postgres/uses_database_sql.go' "package postgres`n`nimport _ `"database/sql`"`n"
$driverResult = Test-GovernanceDrift -RepositoryRoot $driverFixture
Assert-True (-not $driverResult.Passed) 'database/sql under adapters/postgres produced no finding'
Assert-True ('GOV_POSTGRES_DRIVER' -in @($driverResult.Violations.ErrorCode)) "expected GOV_POSTGRES_DRIVER, got: $(@($driverResult.Violations.ErrorCode) -join ',')"
```

Se `orders→catalog` constar como dependência permitida em `modules.json`, a fixture de layer
ainda fecha: `GOV_MODULE_LAYER` dispara independente da aresta ser permitida (o predicado de
`Policy.psm1:378` não consulta `dependencies`). A asserção exige apenas `-in`, então a
presença adicional de `GOV_MODULE_DEPENDENCY` não a quebra.

- [ ] **Step 2: Rodar o arquivo**

```bash
pwsh -NoProfile -ExecutionPolicy Bypass -File ./scripts/tests/governance-drift.tests.ps1
```

Esperado: os tokens PASS habituais do arquivo, sem falha nova, incluindo as 4 asserções novas.

- [ ] **Step 3: Commit**

```bash
git add scripts/tests/governance-drift.tests.ps1
git commit -m "test(governance): positive fixtures for the last two unfixtured codes (#3)"
```

---

### Task 5: Inventário de guardas + lane `guards`

**Files:**
- Create: `contracts/gate/guards.json`
- Modify: `scripts/harness/Gate.psm1` (nova `Measure-GateGuards` + export)
- Create: `scripts/tests/guards-lane.tests.ps1`
- Modify: `scripts/gate.ps1` (nova `Invoke-GateGuards`; `ValidateSet` na linha `:38`; tabela
  `$lanes` `:1034-1051`; `$fullOnly` `:1056`)
- Modify: `.github/workflows/ci.yml` (step no job `verify-full`, após o step `selftest` `:319`)

Semântica de execução, gravada no próprio JSON: entradas `go_tests` são EXECUTADAS pelo lane
(`go test -run` direcionado, prova = `--- PASS: <nome>`); entradas `pwsh_files` delegam a
execução ao lane `selftest` (mesmo run `full`) e o lane `guards` verifica a precondição da
delegação (arquivo dentro do glob de descoberta do selftest) + âncora presente; `presence_only`
é presença documentada com o porquê.

- [ ] **Step 1: `Measure-GateGuards` em `scripts/harness/Gate.psm1`** (junto das demais
  `Measure-*`; o módulo exporta por lista explícita — acrescentar `Measure-GateGuards` ao
  `Export-ModuleMember` de `Gate.psm1:444`):

```powershell
function Measure-GateGuards {
  <#
    Holds a targeted `go test -run -v` output against the inventory's expected
    test names. `\b` after the escaped name: `--- PASS: TestA` must not satisfy
    an expectation of `TestAB`, and vice-versa Go suffixes the line with the
    duration so a bare $ anchor would never match.
  #>
  [CmdletBinding()]
  param(
    [Parameter(Mandatory)][AllowEmptyString()][string]$Text,
    [Parameter(Mandatory)][AllowEmptyCollection()][string[]]$Expected
  )
  $passed = [Collections.Generic.List[string]]::new()
  $missing = [Collections.Generic.List[string]]::new()
  foreach ($name in $Expected) {
    if ([regex]::IsMatch($Text, "(?m)^\s*--- PASS: $([regex]::Escape($name))\b")) { [void]$passed.Add($name) } else { [void]$missing.Add($name) }
  }
  return [pscustomobject]@{
    Ran     = @([regex]::Matches($Text, '(?m)^\s*=== RUN\s')).Count
    Passed  = @($passed)
    Missing = @($missing)
    Failed  = @([regex]::Matches($Text, '(?m)^\s*--- FAIL:')).Count
  }
}
```

- [ ] **Step 2: Self-test vermelho-primeiro — `scripts/tests/guards-lane.tests.ps1`**

Espelhar as convenções de `gate-measure.tests.ps1` (mesmo import do módulo, mesma `Assert-True`;
conferir o cabeçalho daquele arquivo e copiar o boilerplate de import). Corpo:

```powershell
$happy = @"
=== RUN   TestAlpha
--- PASS: TestAlpha (0.00s)
=== RUN   TestBeta
--- PASS: TestBeta (0.01s)
PASS
ok      marketplace-central/apps/server_core/internal/arch     0.117s
"@
$m = Measure-GateGuards -Text $happy -Expected @('TestAlpha', 'TestBeta')
Assert-True ($m.Missing.Count -eq 0) "happy path reported missing: $($m.Missing -join ',')"
Assert-True ($m.Ran -eq 2) "happy path ran=$($m.Ran), want 2"

$m = Measure-GateGuards -Text $happy -Expected @('TestAlpha', 'TestBeta', 'TestGamma')
Assert-True ($m.Missing -contains 'TestGamma') 'a renamed/deleted inventory test must surface as missing'

$vacuous = @"
testing: warning: no tests to run
PASS
ok      marketplace-central/apps/server_core/internal/arch     0.021s [no tests to run]
"@
$m = Measure-GateGuards -Text $vacuous -Expected @('TestAlpha')
Assert-True ($m.Ran -eq 0) 'a no-tests-to-run output must count zero RUN lines'
Assert-True ($m.Missing -contains 'TestAlpha') 'a vacuous run satisfies no expectation'

$prefix = @"
=== RUN   TestAlphaBravo
--- PASS: TestAlphaBravo (0.00s)
PASS
"@
$m = Measure-GateGuards -Text $prefix -Expected @('TestAlpha')
Assert-True ($m.Missing -contains 'TestAlpha') 'TestAlphaBravo must not satisfy TestAlpha (word boundary)'

Write-Output 'PASS guards-lane measurements'
```

Rodar ANTES do Step 1 estar completo é o vermelho (função inexistente →
`Measure-GateGuards: command not found`). Rodar depois:

```bash
pwsh -NoProfile -ExecutionPolicy Bypass -File ./scripts/tests/guards-lane.tests.ps1
```

Esperado: `PASS guards-lane measurements`.

- [ ] **Step 3: `contracts/gate/guards.json`**

Antes de gravar, VERIFICAR cada nome de teste na árvore (`grep -n "^func Test" apps/server_core/internal/arch/scan_test.go`)
— um nome errado aqui nasce como o exato defeito que o lane existe para pegar. Semear:

```json
{
  "$comment": "Issue #3: every guard the delivery path relies on, keyed to the fixture that proves it can fail. go_tests are EXECUTED by the guards lane (proof = '--- PASS: <test>' from a targeted -run). pwsh_files delegate execution to the selftest lane, which runs every scripts/tests/*.tests.ps1 (minus *.integration.*) and requires a pass token per file; the guards lane verifies that delegation precondition plus the anchor. presence_only entries state why execution cannot be asserted here.",
  "go_tests": [
    { "id": "arch-cross-context-internal", "package": "./internal/arch", "test": "TestCrossContextInternalFiresOnFixture" },
    { "id": "arch-float-in-contracts", "package": "./internal/arch", "test": "TestFloatInContractsFiresOnFixture" },
    { "id": "arch-vendor-token", "package": "./internal/arch", "test": "TestVendorTokenFiresOnStringLiteral" },
    { "id": "arch-fact-discard", "package": "./internal/arch", "test": "TestFactValueDiscardIsCaught" },
    { "id": "arch-file-census", "package": "./internal/arch", "test": "TestCountFilesSeesTheFixtureTree" },
    { "id": "adr023-boundary-detector", "package": "./internal/composition", "test": "TestModuleBoundaryDetectorFiresOnFixture" }
  ],
  "pwsh_files": [
    { "id": "gov-module-coverage", "file": "scripts/tests/governance-drift.tests.ps1", "anchor": "GOV_MODULE_COVERAGE" },
    { "id": "gov-module-dependency", "file": "scripts/tests/governance-drift.tests.ps1", "anchor": "GOV_MODULE_DEPENDENCY" },
    { "id": "gov-module-layer", "file": "scripts/tests/governance-drift.tests.ps1", "anchor": "GOV_MODULE_LAYER" },
    { "id": "gov-composition-missing", "file": "scripts/tests/governance-drift.tests.ps1", "anchor": "GOV_COMPOSITION_MISSING" },
    { "id": "gov-application-import", "file": "scripts/tests/governance-drift.tests.ps1", "anchor": "GOV_APPLICATION_IMPORT" },
    { "id": "gov-postgres-driver", "file": "scripts/tests/governance-drift.tests.ps1", "anchor": "GOV_POSTGRES_DRIVER" },
    { "id": "gov-production-panic", "file": "scripts/tests/governance-drift.tests.ps1", "anchor": "GOV_PRODUCTION_PANIC" },
    { "id": "gov-api-sdk-split", "file": "scripts/tests/governance-drift.tests.ps1", "anchor": "GOV_API_SDK_SPLIT" },
    { "id": "gov-frontend-fetch", "file": "scripts/tests/governance-drift.tests.ps1", "anchor": "GOV_FRONTEND_FETCH" },
    { "id": "gov-migration-prefix", "file": "scripts/tests/governance-drift.tests.ps1", "anchor": "GOV_MIGRATION_PREFIX" },
    { "id": "gov-context-unregistered", "file": "scripts/tests/governance-drift.tests.ps1", "anchor": "GOV_CONTEXT_UNREGISTERED" },
    { "id": "integration-vacuous-run", "file": "scripts/tests/postgres-lifecycle.tests.ps1", "anchor": "HPG_TEST_VACUOUS" },
    { "id": "gate-counters", "file": "scripts/tests/gate-measure.tests.ps1", "anchor": "Measure-GateGoTest" },
    { "id": "guards-lane-counters", "file": "scripts/tests/guards-lane.tests.ps1", "anchor": "Measure-GateGuards" }
  ],
  "presence_only": [
    { "id": "required-job-aggregator", "file": ".github/workflows/ci.yml", "anchor": "guard self-test", "why": "executes inside the required job itself on every PR; its five jq scenarios are the execution evidence" },
    { "id": "main-ruleset", "file": "docs/HARNESS-PROFILE.md", "anchor": "main-landing-path", "why": "GitHub-side control; must-fail run 2026-08-09 (GH013), recorded in issue #4's closing comment" }
  ]
}
```

- [ ] **Step 4: `Invoke-GateGuards` em `scripts/gate.ps1`** (inserir após `Invoke-GateSelftest`):

```powershell
function Invoke-GateGuards {
  # Issue #3 mechanism (a): the inventory of guards, each held to the fixture
  # that proves it can fail -- and the fixture must EXECUTE, not merely exist.
  $inventoryPath = Join-Path $repositoryRoot 'contracts/gate/guards.json'
  if (-not (Test-Path -LiteralPath $inventoryPath -PathType Leaf)) {
    return New-GateVerdict -Lane 'guards' -Passed $false -Counts 'inventory=missing' `
      -Reason 'contracts/gate/guards.json does not exist; an unlisted guard population is unmeasurable'
  }
  $inventory = Get-Content -LiteralPath $inventoryPath -Raw | ConvertFrom-Json
  $goEntries = @($inventory.go_tests); $pwshEntries = @($inventory.pwsh_files); $presenceEntries = @($inventory.presence_only)
  $total = $goEntries.Count + $pwshEntries.Count + $presenceEntries.Count
  if ($total -eq 0) {
    return New-GateVerdict -Lane 'guards' -Passed $false -Counts 'entries=0' `
      -Reason 'an empty inventory is not a green inventory'
  }

  $failures = [Collections.Generic.List[string]]::new()
  $goExecuted = 0
  foreach ($group in @($goEntries | Group-Object -Property package)) {
    $names = @($group.Group | ForEach-Object { [string]$_.test })
    $pattern = '^(' + (@($names | ForEach-Object { [regex]::Escape($_) }) -join '|') + ')$'
    $toolName = 'guards-' + (($group.Name -replace '[^A-Za-z0-9]+', '-').Trim('-'))
    $result = Invoke-GateTool -Name $toolName -FilePath 'go' `
      -ArgumentList @('test', $group.Name, '-run', $pattern, '-count=1', '-v') -WorkingDirectory $serverCore
    $measurement = Measure-GateGuards -Text $result.Text -Expected $names
    if ($result.ExitCode -ne 0) { [void]$failures.Add("package $($group.Name) exited $($result.ExitCode)") }
    if ($measurement.Ran -eq 0) { [void]$failures.Add("package $($group.Name): zero tests ran; the inventory names tests that do not execute") }
    foreach ($name in $measurement.Missing) { [void]$failures.Add("no '--- PASS: $name' in $($group.Name)") }
    $goExecuted += @($measurement.Passed).Count
  }

  foreach ($entry in $pwshEntries) {
    $file = Join-Path $repositoryRoot ([string]$entry.file)
    if (-not (Test-Path -LiteralPath $file -PathType Leaf)) { [void]$failures.Add("$($entry.file) is missing"); continue }
    $leaf = [IO.Path]::GetFileName($file)
    $inSelftestGlob = ($leaf -like '*.tests.ps1') -and ($leaf -notlike '*.integration.tests.ps1') -and
      ([IO.Path]::GetDirectoryName($file) -eq (Join-Path $repositoryRoot 'scripts/tests'))
    if (-not $inSelftestGlob) { [void]$failures.Add("$($entry.file) sits outside the selftest lane's discovery glob; its execution is delegated to nothing") }
    if (-not (Select-String -LiteralPath $file -Pattern ([regex]::Escape([string]$entry.anchor)) -Quiet)) {
      [void]$failures.Add("$($entry.file): anchor '$($entry.anchor)' not found")
    }
  }

  foreach ($entry in $presenceEntries) {
    $file = Join-Path $repositoryRoot ([string]$entry.file)
    if (-not (Test-Path -LiteralPath $file -PathType Leaf)) { [void]$failures.Add("$($entry.file) is missing"); continue }
    if (-not (Select-String -LiteralPath $file -Pattern ([regex]::Escape([string]$entry.anchor)) -Quiet)) {
      [void]$failures.Add("$($entry.file): anchor '$($entry.anchor)' not found")
    }
  }

  $counts = "entries=$total go_executed=$goExecuted pwsh=$($pwshEntries.Count) presence=$($presenceEntries.Count) failures=$($failures.Count)"
  Write-Host $counts
  if ($failures.Count -gt 0) {
    foreach ($line in $failures) { Write-Host "  guard: $line" }
    return New-GateVerdict -Lane 'guards' -Passed $false -Counts $counts `
      -Reason 'a registered guard has no executing fixture. A guard that cannot demonstrate a catch is indistinguishable from an absent one.'
  }
  return New-GateVerdict -Lane 'guards' -Passed $true -Counts $counts
}
```

Wiring no mesmo arquivo: `ValidateSet` (linha `:38`) ganha `'guards'`; `$lanes` ganha
`'guards' = ${function:Invoke-GateGuards}` logo após `'selftest'`; `$fullOnly` vira
`@('selftest', 'guards', 'integration', 'edge')`.

- [ ] **Step 5: Vermelho do lane — controle negativo nomeado**

Editar `contracts/gate/guards.json` na árvore de trabalho, adicionando a `go_tests`:
`{ "id": "bogus", "package": "./internal/arch", "test": "TestGuardInventoryBogus" }`. Rodar:

```bash
pwsh -NoProfile -File ./scripts/gate.ps1 -Lane guards
```

Esperado: `lane=guards verdict=FAIL`, com a linha `guard: no '--- PASS: TestGuardInventoryBogus' in ./internal/arch`.
Remover a entrada boga (nunca commitá-la) e capturar a saída vermelha para a evidência do PR.

- [ ] **Step 6: Verde do lane**

```bash
pwsh -NoProfile -File ./scripts/gate.ps1 -Lane guards
```

Esperado: `entries=22 go_executed=6 pwsh=14 presence=2 failures=0`, `gate: PASS`.

- [ ] **Step 7: Step no CI — job `verify-full`, após o step `selftest` (`ci.yml:319`)**

```yaml
      - name: guards
        shell: pwsh
        run: ./scripts/gate.ps1 -Lane guards
```

Pré-condição já satisfeita (medida): o job `verify-full` tem `actions/setup-go@v7` em
`ci.yml:291`. NÃO tocar em `needs`/`EXPECTED` do job `required` (`ci.yml:351,354`) — steps
novos em jobs existentes não alteram a topologia.

- [ ] **Step 8: Selftest + census continuam verdes** (o arquivo novo entra no glob do selftest
  automaticamente; census vê `guards-lane.tests.ps1` como pwsh, fora do censo Go/vitest):

```bash
pwsh -NoProfile -File ./scripts/gate.ps1 -Lane selftest
```

Esperado: `files=<n+1> pass=<n+1> fail=0`, `gate: PASS`.

- [ ] **Step 9: Commit**

```bash
git add contracts/gate/guards.json scripts/harness/Gate.psm1 scripts/tests/guards-lane.tests.ps1 scripts/gate.ps1 .github/workflows/ci.yml
git commit -m "gate(guards): inventory of guards, each held to an executed failing fixture (#3)"
```

---

### Task 6: Retirar `scripts/arch-gate.sh` + linha no amendment log

**Files:**
- Delete: `scripts/arch-gate.sh`
- Modify: `docs/HARNESS-PROFILE.md` (Amendment log, topo do bloco)

- [ ] **Step 1: Re-medir órfão no tip atual** (o grep desta sessão foi em `2854226`):

```bash
grep -rn "arch-gate" package.json .github scripts/gate.ps1 scripts/harness.ps1
```

Esperado: zero hits (as 14 referências existentes são docs/dívidas/planos históricos —
universo: árvore inteira via ripgrep, medido 2026-08-09).

- [ ] **Step 2: Apagar e registrar**

```bash
git rm scripts/arch-gate.sh
```

Nova linha no TOPO do bloco do Amendment log de `docs/HARNESS-PROFILE.md`:

```
2026-08-09 · §2 · ratified · guards lane added to the gate (issue #3: inventory at contracts/gate/guards.json, every registered guard held to an EXECUTED failing fixture; targeted go-test proof for Go detectors, delegation-to-selftest verified for pwsh fixtures) · scripts/arch-gate.sh DELETED, measured retirement: zero executable references (package.json, ci.yml, gate.ps1, harness.ps1 — all clean at 2854226); the D-51 if-defeats-errexit pattern it carried dies with it
```

- [ ] **Step 3: Commit**

```bash
git add docs/HARNESS-PROFILE.md
git commit -m "gate: retire arch-gate.sh -- zero executable references, superseded by gate lanes (#3)"
```

---

### Task 7: Fecho — gate completo, PR, evidência

- [ ] **Step 1: Gate local no conjunto barato** (integração/edge precisam de Docker; rodar se
  o Docker local estiver de pé, senão delegar ao runner):

```bash
pwsh -NoProfile -File ./scripts/gate.ps1 -Lane all
```

Esperado: `gate: PASS`, com `packages=N` visível nos lanes `build` e `lint-go`.

- [ ] **Step 2: Push do branch (autorização em pé para branch de trabalho, HARNESS-PROFILE §9)**

```bash
git push -u origin issue-3-guard-inventory
```

- [ ] **Step 3: PR**

```bash
gh pr create --title "gate: guard inventory -- evidence stops certifying itself (#3)" --body "Closes #3 (re-scoped 2026-08-09 with the operator; measurement showed #2 already paid mechanism (c) and most of (b)).

- (a) contracts/gate/guards.json + guards lane: 22 entries; Go detectors executed by targeted -run with '--- PASS:' proof; pwsh fixtures delegated to selftest with the delegation precondition checked; must-fail reproduced locally (bogus entry -> named FAIL line, output below)
- (b) build and lint-go lanes gain the package-universe floor (the last 2 of 16 lanes without a zero check)
- ADR-023 detector now fires on a committed fixture -- its red no longer depends on the legacy tree staying dirty
- GOV_MODULE_LAYER and GOV_POSTGRES_DRIVER get the last two missing positive fixtures
- scripts/arch-gate.sh retired: zero executable references

<colar aqui a saída vermelha do Step 5 da Task 5 e a verde do Step 6>

🤖 Generated with [Claude Code](https://claude.com/claude-code)"
```

- [ ] **Step 4: Checks verdes no runner** (`gh pr checks --watch`). Merge = permissão do
  operador, por evento (§9). Depois do merge: fechar #3 é automático pelo `Closes #3`;
  comentar no issue as correções de evidência envelhecida (GOV_CONTEXT_UNREGISTERED já
  implementado com fixture em `Policy.psm1:350-358` + `governance-drift.tests.ps1:278-283`;
  HPG_TEST_VACUOUS já coberto) para o histórico do audit.

- [ ] **Step 5: Dívidas** — nenhuma nova; D-51 morre com o arquivo (registrar a baixa em
  `.mnfs/HARNESS-DEBTS.md` se D-51 constar como aberto lá, citando a task 6).

---

## Fora de escopo (inalterado do issue)

Auditar a CORREÇÃO da lógica de cada guarda existente. Este plano prova que cada guarda
consegue falhar, não que falha nas coisas certas.
