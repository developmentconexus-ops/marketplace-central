# Medição de fecho — fundação do kernel (Tarefa 12, 2026-08-06)

Cada número abaixo foi obtido pelo comando ao lado, correndo de
`apps/server_core`, na árvore como está depois do commit `42b6ed1`
(`feat(arch): portao de arquitetura provado a reprovar e a aprovar`). Nenhum
valor foi estimado.

```
$ cd apps/server_core
$ echo "kernel packages:      $(ls internal/kernel | wc -l)"
kernel packages:      6
$ echo "context packages:     $(ls internal/contexts | wc -l)"
context packages:     1
$ echo "float in kernel:      $(grep -rn 'float64\|float32' internal/kernel | wc -l)"
float in kernel:      2
$ echo "float in contracts:   $(grep -rn 'float64\|float32' internal/contexts/*/contracts | wc -l)"
float in contracts:   0
$ echo "modules touched:      $(git diff --name-only HEAD~12..HEAD -- internal/modules | wc -l)"
modules touched:      0
$ echo "new dependencies:     $(git diff HEAD~12..HEAD -- go.mod | grep -c '^+\s' || true)"
new dependencies:     0
```

| medida | esperado | medido | desvio |
|---|---|---|---|
| kernel packages | 6 | 6 | nenhum |
| context packages | 1 | 1 | nenhum |
| float in kernel | 0 | 2 | ver D-32 |
| float in contracts | 0 | 0 | nenhum |
| modules touched | 0 | 0 | nenhum |
| new dependencies | 0 | 0 | nenhum |

## Desvio: float in kernel = 2, não 0

Os dois hits são comentários, não código executável:

```
internal/kernel/exact/decimal.go:2:// from float64 anywhere in this package, and that is the point: a binary float
internal/kernel/exact/money.go:44:// constructor from float64.
```

Ambos explicam a proibição de `float64` em prosa (doc comment do pacote e de
um construtor). O `grep -rn 'float64\|float32'` — tanto o desta medição como
o do Step "no float in the kernel" de `scripts/arch-gate.sh` — não distingue
comentário de código. Não há nenhum `float64`/`float32` real em código
executável do kernel. Registado como **D-32** em `.mnfs/HARNESS-DEBTS.md`;
não corrigido nesta tarefa (corrigir qualquer lado — reescrever os
comentários ou trocar o `grep` cru por algo AST-consciente — é uma escolha
fora do mandato aditivo desta tarefa).

## Portão corrido do zero contra a árvore de hoje

```
$ ./scripts/arch-gate.sh; echo "EXIT=$?"
...
ARCH GATE: FAIL
EXIT=1
```

**NÃO é `PASS`, ao contrário do "Esperado" do brief.** A árvore como está
hoje reprova o portão por cinco causas independentes, nenhuma introduzida
por esta tarefa: gofmt sinaliza ~637 ficheiros por causa de
`core.autocrlf` (artefacto de ambiente, não de conteúdo), `go vet ./...`
falha por D-30 (`internal/composition/catalog_wiring.go:9`), `archscan
-root internal` encontra tokens de vendor em `internal/modules`/
`internal/composition` (código legado ainda não migrado ao protocolo, mais
D-28 no próprio kernel), `go test ./internal/...` falha em
`internal/composition` (D-30) e em `internal/platform/migrate` (drift de
inventário de migração pré-existente, 84 vs 83), e o bloco "no float in the
kernel" nunca fecha por D-32. Detalhe completo de cada causa em **D-33**
(`.mnfs/HARNESS-DEBTS.md`). A prova crítica do Step 6 — o portão reprova
quando `float64` é injetado no kernel, apontando a linha certa, e a remoção
do ficheiro de sonda não deixa resíduo no `git status` — permanece válida
isoladamente; é o veredito AGREGADO sobre a árvore inteira que é `FAIL`.
