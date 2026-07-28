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
