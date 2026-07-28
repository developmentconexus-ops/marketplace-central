# Round 8 — the arms of the signal mechanism, run red one rule at a time

The mechanism (`wireCandidate`'s signal-set check and its signal-vs-absence check)
was already proved to turn REAL fixtures red: `V-signals-round8-RED.log`, 3 tests
failing across 3 files. That proves the mechanism finds defects. It does NOT prove
each individual arm in `wireFixtures.guard.test.ts` is attributable to the rule it
names — an arm that would go red for any reason is not evidence about that rule.

So each rule was neutered ALONE and the arms re-run. Baseline of the unmutated
tree: 8 passed (8).

`wireFixtures.ts` md5 before any mutation and after the restore:
`d1665a51355ba93f87d578aaf4d33756` — identical, verified by `md5sum` both times.
The backup lived outside the worktree (scratchpad), so the restore is a copy, not
a `git` operation.

## Mutation 1 — signal-set check disabled (`if (false && !tupleMatches.some(...))`)

```
   ✓ MUST-FAIL 1 — ACCEPT exists only as (exact_sku, seller_sku) 1ms
   ✓ MUST-FAIL 2 — a provider with no capability declaration is not producible 0ms
   × MUST-FAIL 3 — a producible tuple with a signal set no site emits 11ms
   ✓ MUST-FAIL 4 — one anchor cannot carry a signal and an absence at once 0ms
   ✓ MUST-PASS — `title` FOR and `title` AGAINST on one candidate is real, and stays accepted 0ms
 Test Files  1 failed (1)
      Tests  1 failed | 7 passed (8)
```

Exactly one arm fell, and it is the one that names this rule. MUST-FAIL 1 staying
green is the useful half: the tuple table alone does NOT catch a wrong signal set,
which is why arm 3 had to exist.

## Mutation 2 — signal-vs-absence check disabled (`if (false && signalAnchors.has(...))`)

```
   ✓ MUST-FAIL 1 ... 2ms
   ✓ MUST-FAIL 2 ... 1ms
   ✓ MUST-FAIL 3 ... 1ms
   × MUST-FAIL 4 — one anchor cannot carry a signal and an absence at once 17ms
   ✓ MUST-PASS ... 1ms
 Test Files  1 failed (1)
      Tests  1 failed | 7 passed (8)
```

Again exactly one, and it is the matching one.

## Mutation 3 — the historical WRONG form of the dup-anchor rule

This is the arm that no must-fail can produce, and the reason MUST-PASS is in the
file at all. A guard that is too STRICT makes every must-fail greener; the only
thing that can catch it is a producible candidate asserted not to throw.

The mutation is not an invention — it is the form this rule actually had when it
was first written this round: the dup-absence check hoisted ABOVE the FOR/AGAINST
split, so signals were counted as absences.

```
   ✓ MUST-FAIL 1 ... 1ms
   ✓ MUST-FAIL 2 ... 0ms
   ✓ MUST-FAIL 3 ... 0ms
   × MUST-FAIL 4 — one anchor cannot carry a signal and an absence at once 11ms
   × MUST-PASS — `title` FOR and `title` AGAINST on one candidate is real, and stays accepted 1ms
 Test Files  1 failed (1)
      Tests  2 failed | 6 passed (8)
```

MUST-PASS red is the result being claimed: the wider rule rejects
`applySingleAnchorScore`'s hard-negative candidate, which emits `title` FOR
(`generation_service.go:550-554`) and `title` AGAINST (`:560-562`) together, both
passed through by `:657-661`.

MUST-FAIL 4 also fell here, and NOT for its own reason: under the hoisted check the
`seller_sku` FOR + `seller_sku` UNAVAILABLE fixture throws `carries two absence
reasons` instead of `carries both a signal and a UNAVAILABLE absence`, so its
`toThrow` regex misses. Said out loud because a red count of 2 would otherwise read
as two independent findings; it is one finding and one message collision.

## A first attempt at mutation 3, discarded, and why it is recorded

The first version deleted the split line entirely (`if (false) continue;`). It
turned FOUR arms red — 1, 3, 4 and MUST-PASS — because with the split gone the
FOR reasons also reach the signal-vs-absence check and throw the wrong message
there too. Red, but not attributable: a mutation that breaks several rules at once
cannot tell you which rule an arm is measuring. Replaced with the narrower
reorder above rather than reported as a stronger result.

## Lane on the restored tree

- `npx vitest run` — 67 files / 550 tests, exit 0, 0 `FAIL` lines
  (`V-signals-round8-arms-GREEN.log`)
- `npx --no-install tsc --noEmit -p tsconfig.json` — 12 errors total, **0** under
  `pages/vinculos` (`V-tsc-round8.log`); same 12 as every round since the merge
- `npx vitest run -t "MUST-"` — 6 passed | 544 skipped (`V-mustfail-round8.log`);
  6 = the 5 arms in this describe + the one pre-existing `MUST-` test elsewhere
