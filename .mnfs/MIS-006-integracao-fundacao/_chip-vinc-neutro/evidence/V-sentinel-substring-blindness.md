# The sentinel had the shape it names — found by the hub's executor seat

Raised by the hub's independent executor seat, run in its own detached worktree at `2e5331b6`.
Verified here against primary source before touching anything, then fixed red-first.

## The defect, by string

`wireFixtures.guard.test.ts:127` asserted `.toContain(symbol)` with `symbol = "var knownIdentityAnchors"`.

```
'var knownIdentityAnchorsXX'.find('var knownIdentityAnchors') != -1   →   True
```

Substring containment is not symbol presence. A rename that **appends** leaves the substring
intact, so the sentinel passed while the extraction below it matched nothing. The guard was
wider than the fact it claimed to check — the same class this chip has been fixing one layer
down all week, sitting inside the mechanism built to end it.

The test's own comment named the mode it did not fully cover:
*"A rename leaves the file readable and the guard blind, which is the failure mode worth naming."*

## The count was true of a mutation the pack did not declare

The pack reported the arm as `RED, 3/3`. That is the number for mutating **both** seams — the
sentinel loops over both `GO_SEAM` entries. Mutating **one** seam gives a different number, and
under a suffix rename it gave a number with the sentinel *green*. Per the §11 reconciliation
clause, count and mutation are declared together:

| mutation | before the fix | after the fix | sentinel |
|---|---|---|---|
| MUT-1 suffix rename, port only (`knownIdentityAnchorsXX`) | **1 failed / 2 passed** | **2 failed / 1 passed** | was ✓ **blind** → now ✗ |
| MUT-2 total rename, port only (`identityAnchorVocabulary`) | 2 failed / 1 passed | 2 failed / 1 passed | ✗ |
| MUT-3 mercado_livre seam only (`IdentityAnchors:` → `SuppliedAnchors:`) | — | 2 failed / 1 passed | ✗ |
| MUT-4 both seams | 3 failed | **3 failed** | ✗ — this is the pack's `3/3` |
| baseline, tree intact | 3 passed | 3 passed | ✓ |

Restored after every arm: `git diff --stat apps/server_core` = 0 lines. No Go file was left mutated.

MUT-1 and MUT-2 now produce the same result. That equality IS the fix: a suffix rename and a
total rename are one class, and the guard used to answer differently for them.

## The fix, and why it is this fix

The extraction pattern moved INTO `GO_SEAM` as `extract`, and the sentinel now asserts
`GO_SEAM[key].extract.test(text)` — the same object the tests below run. Sentinel and extraction
cannot disagree by construction. A tighter second pattern (`\bsymbol\b\s*=`) would also have gone
red on MUT-1, but it would have left two patterns free to drift apart again, which is what
produced this defect in the first place.

## Reconciliation of the fix itself

Population = every regex literal in the guard that reads Go. Extraction = those the sentinel
covers. Both printed:

```
POPULACAO = 7 literais casados  →  2 sao FALSO POSITIVO
  /../  e  /../s   vieram da STRING DE CAMINHO "../../../../server_core/..." na linha 37
POPULACAO real = 5
  2  localizadores de costura   → GO_SEAM.extract, 3 usos (sentinela + 2 testes)
  3  extracoes secundarias      → dentro de bloco ja localizado, cada uma guardada
  0  padroes sem guarda
assercoes de extracao-zero no arquivo = 5 (text.length, constants.size x2, anchors.length, declared.length)
```

The two false positives are the class again, in the tool measuring the class: the pattern matched
a path string. Printed rather than quietly subtracted, because a residual removed silently is
indistinguishable from a residual that was never there.

Also of record, same category, same hour: the first run of these four arms used
`grep -E "^ *(Test Files|Tests)"` and returned **empty with exit 0** on all four. Vitest writes an
ANSI escape before the leading space, so `^ *` never matches. The hub's executor seat reported
hitting the identical trap independently. Silence read as clean, one hour after the amendment that
names exactly this. Fixed by stripping ANSI *before* the grep, not after.

## Merge-base correction

`git merge-base main HEAD` = **`bcab8269`**, not `5441fe18` (which is an ancestor of it). The
comment at `VinculosPage.test.tsx:38` cited the wrong point while explaining why the
`listErpImports` mock stays. Corrected in place, with the measurement named. No behavioural
consequence — the mock is load-bearing in this tree either way.

## Lanes after the fix

```
tsc     12 errors, 0 under pages/vinculos
vitest  64 files / 531 tests, all green
```
