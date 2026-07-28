> **CORRECTED AT SOURCE (R-24) — the population below was measured with `grep -c`, which counts
> LINES, not matches. It read 51; there are 53. Read `## CORRECTION` at the end before using any
> number in this file. The two counts and the residual are restated there.**

# Count reconciliation, run by this chip against its OWN pack

Binding: `docs/HARNESS-PROFILE.md` §11, amendment `0cb6d7e` — *"a sweep is only as wide as its
pattern; reconcile the extraction against the population."* The hub made it explicit that this
binds the implementer too: *"você roda a mesma varredura contra o próprio pack antes de publicar."*

Run at `2e5331b6`, before freezing. Both counts printed, per the amendment.

## The claim being reconciled

The pack claims every candidate fixture under `apps/web/src/pages/vinculos/` is constructed
through `wireCandidate` / `driftCandidate`, which THROW on an unproducible candidate.

## The two counts

Population anchored loosely — `candidate_id:` in every `*.test.ts*` under `apps/web/src`, NOT
only under `pages/vinculos/`, because the 29th fixture survived four rounds precisely by living
in a file the sweep did not own.

```
POPULACAO=51  EXTRACAO=36  REFS=15  soma_confere=True
```

| | count | what it is |
|---|---|---|
| population | **51** | every `candidate_id:` literal in a test file |
| extraction | **36** | candidate objects built through the throwing constructor |
| residual | **15** | batch payload refs — a DIFFERENT union, different producer |
| | 36 + 15 = 51 | the residual is named, not left as slack |

Per file: `QueueTab.test.tsx` 30 built / 7 refs · `VinculosDesign.golden.test.tsx` 5 built
(1 in `base()`, 4 overrides routed through it) · `VinculosPage.test.tsx` 1 built ·
`BatchPreviewModal.test.tsx` 0 built / 8 refs.

## Must-fail — the amendment's second clause

A pattern never shown matching a known member is not known to match anything. Run on a copy in
the scratchpad, so the tree was never dirtied; the copy was diffed against the tree afterwards
and is byte-identical apart from one trailing newline from the restore.

- **arm A** — raw candidate literal injected into a file the sweep already names →
  `NAO-CONSTRUIDO -> QueueTab.test.tsx:948`, population 51→52, extraction unchanged.
- **arm B** — raw candidate in a NEW file no sweep ever named (the 29th-fixture class) →
  `NAO-CONSTRUIDO -> ArmB.test.tsx:1`, population 52→53, extraction still 36.

Restored: back to `51 / 36 / 15`.

## What the must-fail caught in the sweep itself

The first run of both arms moved the population (51→53) and printed **nothing**. The classifier
routed both injected candidates into the "batch ref" bucket, silently, because its pattern was
`status:` — and `match_status:` contains `status:`.

This is the same defect as round 5's `grep -oh 'anchor: "[a-z_]*"'`, which dropped every anchor
carrying a capital or an accent and reported 1 violation where there were 5. Both times: a loose
pattern, the wrong bucket, no error. Fixed with a word boundary — `(?<![\w])(?:status|cause|…):`
— and both arms then went red as shown above.

A second, smaller one, found the same way: classifying by a 6-line lookback window marked
`QueueTab.test.tsx:633-634` (batch refs) as *built*, because a real `candidate({…})` call sits
four lines above them. Hand count said 30/7 for that file, the window said 32/5; the window was
wrong. Classifying by the line's own content BEFORE consulting the window fixes it. The two
counts disagreeing is what surfaced it — a single number would have been reported as fact.

## What the residual named — REPORTED, not fixed

The 15 residual sites are not noise. They feed a **second union that has no mechanism**:

`ProductLinkBatchPreviewItem.status: "OK" | "FAILED"` (`packages/sdk-runtime/src/index.ts:1160-1164`).
`cause` beside it is a free-form `string`, so no typeset can constrain it at all.

`BatchPreviewModal.tsx` enumerates that union by string literal, in three places that disagree
if it ever grows a member:

- `:44` `filter(status === "OK")` → builds `validApprovals`, i.e. the **apply payload**;
- `:46` `filter(status === "FAILED")` → the counter at `:87`;
- `:94` `status === "OK" ? … : …` → the row rendering.

A third member is therefore: **not applied** (absent from `validApprovals`), **counted in
neither** chip at `:84`/`:87` — so the two chips stop summing to `items.length` — and yet
**rendered in the list as a failure**, via the ternary's else branch at `:97`. Three
inconsistent readings of one item, all silent, all type-correct.

This is not a defect today: the union has exactly two members and the partition is total. It is
the same drift exposure `driftCandidate` exists for, in a union the mechanism never covered and
never claimed to. Filed as a REPORT to the hub, not fixed here — the hub ordered one freeze per
round, and widening the freeze to a component outside the round's delta is the chip deciding
its own scope.

---

## CORRECTION (hub executor seat, round 6 pre-dispatch, per R-24)

The hub's executor seat verified this artifact instead of accepting it, and the population above is
**wrong**. Everything after this heading supersedes the numbers before it.

### The defect: `grep -c` counts LINES

```
grep -rc 'candidate_id:' … | sum   →  51    ← the number reported above
grep -roh 'candidate_id:' | wc -l  →  53    ← occurrences, which is what was CLAIMED
```

The claim was a count of THINGS (fixture sites). The instrument was a count of LINES. The two
diverge exactly when two sites share a line, and two do:

```
pages/vinculos/BatchPreviewModal.test.tsx:46  x2  approvals: [{ candidate_id: "cand_1" }, { candidate_id: "cand_2" }],
pages/vinculos/QueueTab.test.tsx:656          x2  approvals: [{ candidate_id: "cand_1" }, { candidate_id: "cand_2" }],
```

### `soma_confere=True` did not save it, and that is the finding

The original run printed `51 = 36 + 15, soma_confere=True`. The corrected run prints
`53 = 36 + 17, soma_confere=True`. **Both reconcile.** Consistent arithmetic over the wrong
population is still consistent — a reconciliation can close and remain short by exactly as much as
the instrument cannot see. The hub's sentence, kept because it is the lesson: *"uma reconciliação
pode fechar e ainda estar curta exatamente do tanto que o instrumento não enxerga."*

### Corrected counts

| | count | what it is |
|---|---|---|
| population (occurrences) | **53** | every `candidate_id:` MATCH in a test file |
| extraction | **36** | candidate objects built through the throwing constructor — unchanged |
| residual | **17** | batch payload refs; both extra occurrences land here, on `approvals:` lines |
| | 36 + 17 = 53 | |

Per file: `QueueTab.test.tsx` 30 built / 8 refs · `BatchPreviewModal.test.tsx` 0 / 9 ·
`VinculosDesign.golden.test.tsx` 5 / 0 · `VinculosPage.test.tsx` 1 / 0.

### Must-fail of the CORRECTED instrument, three arms

Arm C is new and exists only because of this defect: a raw fixture with **two occurrences on one
line** — the exact shape the old instrument could not count.

```
baseline                     POPULACAO=53  EXTRACAO=36  REFS=17
arm A  raw candidate, known file      -> NAO-CONSTRUIDO, QueueTab.test.tsx:947
arm B  raw candidate, unnamed file    -> NAO-CONSTRUIDO, ArmB.test.tsx:1
arm C  two occurrences on one line    -> NAO-CONSTRUIDO, ArmC.test.tsx:1  (counted x2)
with all three                POPULACAO=57  EXTRACAO=36  REFS=21
```

And the discriminating run — the OLD instrument against the SAME mutated tree:

```
old (per line)        POPULACAO=54  EXTRACAO=36  REFS=18  soma_confere=True
corrected (per match) POPULACAO=57  EXTRACAO=36  REFS=21  soma_confere=True
```

Three short: the two pre-existing double lines plus arm C's second occurrence. The old instrument
reports `soma_confere=True` while missing all three, which is why the arm exists.

### Where this sits in the sequence

Third instance of one class inside this chip's own work, all three an instrument wider or narrower
than the fact: `[a-z_]` dropping capitals and accents (round 5), `status:` matching `match_status:`
(caught by this chip's own must-fail), and now `grep -c` counting lines. The first two were caught
by this chip; this one was not. It was caught by a seat that verified the artifact instead of
reading its conclusion — which is the entire argument for the executor seat existing.

---

## Two more instances of the class, both in the dispatch machinery itself

Recorded because the class has now appeared five times in this chip's own work and the count is
the point: it is not a mistake that gets made once and learned. Neither of these reached a seat —
both were caught by looking at what the command DID instead of what it was for.

### `git ls-tree` without `-r` counts tree entries

The custody clause was committed ordering `git ls-tree HEAD -- <pack>/ | wc -l`. Run:

```
git ls-tree    HEAD -- <pack>/ | wc -l   ->   4
git ls-tree -r HEAD -- <pack>/ | wc -l   ->  40
```

Four is not a wrong file count, it is a correct count of something else: two blobs and one
subtree at the pack root. A seat filling a MANDATORY section would have written `4` truthfully.
The clause exists because a prior chip's pack sat untracked for six commits; the instrument it
shipped with could not have distinguished 40 tracked files from 4.

Also reconciled while here: the hub's executor seat reported **38**; the tree says **39** at
`7b5c18eb` and **40** at the dispatch tip. The +1 is the round-6 brief, added after the freeze.
The 38-vs-39 gap is the executor's measurement and is reported to the hub, not corrected here.

### `"| 9 |"` matched a different table

Re-filing ledger rows 9 and 10 was scripted as `line.startswith("| 9 |")`. `EVIDENCE.md` holds
two tables with those row numbers, and the script rewrote **both** — silently destroying the tsc
error inventory's rows 9 and 10:

```
- | 9  | `src/pages/produto/ProdutoPage.partialFailure.test.tsx(40,45)` | TS2322 same |
- | 10 | `src/pages/produto/ProdutoPage.partialFailure.test.tsx(41,46)` | TS2322 same |
```

The script reported success. `grep -c "DISPATCHED"` returned 4 for two rows, and that discrepancy
— a count that did not match the fact it was checking — is what exposed it; `git diff` then named
the two victims. Restored verbatim, verified by diff: 2 insertions, 2 deletions, both in the
dispatch ledger.

This one is worse than the others in kind. The previous four made a measurement wrong. This one
made an EDIT wrong, and an edit that deletes has no residual to reconcile against — had the
`grep -c` not disagreed, the loss would have been invisible in every subsequent count, because
the deleted rows are not in any population the sweep anchors on. R-25: honest-unknown is for
gaps; a silent deletion is falsity.
