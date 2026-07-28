# V-round9 — B-2 discharged, and the corrected sentence made checkable

Tree: `bfc1d9bb` + round-9 working changes. Measured before commit.

## The finding

`wireFixtures.ts` listed the `side` of an absence among the things `wireCandidate` does NOT
check. `:294` checks it:

```
if (reason.side !== undefined && reason.direction !== "INCOMPARABLE") {
```

The hub sided with Sol against the steelman: *"o docstring contradiz a si mesmo, e isso é
pior que estar incompleto"*. Correct — the "does not check" list is the place a reader
consults precisely so they do not have to read the guard, so a false cell there is worse
than an absent one.

## The correction (R-24, at source)

The gap paragraph now names only the DIRECTION as unchecked, and a following paragraph
states explicitly that `side` is checked, names the rule that checks it, and names the
earlier version as false. `wireFixtures.ts` is comment-only in this round — `git diff
--numstat` is `12 3` and every changed line is inside the block comment, so the guard
MECHANISM is byte-identical to round 8 and round 8's attributability evidence still stands.

## Why a comment fix was not enough

Prose drifts silently; that is the whole finding. So the corrected sentence is now backed by
an arm:

`MUST-FAIL 5 — an UNAVAILABLE absence cannot carry a \`side\``, asserting
`.toThrow(/Only INCOMPARABLE branches set a side/)` on a candidate whose `ean` absence
carries `side: "erp"`.

## Attributability

Per §11 `43bd911c` — an arm proves only what its mutation ISOLATES. Mutation: `:294`
neutered to `if (false && ...)`, nothing else touched.

```
 ✓ MUST-FAIL 1 — ACCEPT exists only as (exact_sku, seller_sku)
 ✓ MUST-FAIL 2 — a provider with no capability declaration is not producible
 ✓ MUST-FAIL 3 — a producible tuple with a signal set no site emits
 ✓ MUST-FAIL 4 — one anchor cannot carry a signal and an absence at once
 ✓ MUST-PASS — `title` FOR and `title` AGAINST on one candidate is real, and stays accepted
 × MUST-FAIL 5 — an UNAVAILABLE absence cannot carry a `side`

 Test Files  1 failed | 5 passed (6)
      Tests  1 failed | 41 passed (42)
```

Exactly one red, and it is the new arm. The mutation isolates the rule the corrected
sentence names.

Restore verified by md5, not by eye:

```
7cdaa910aeef78c0197ee80e9d3c79dc  wireFixtures.ts (pre-mutation)
7cdaa910aeef78c0197ee80e9d3c79dc  wireFixtures.ts (post-restore)
```

`git diff --numstat` for the file is `12 3` before and after, unchanged.

## Lanes on the restored tree

```
vitest src/pages/vinculos   6 files, 42 tests, all pass
vitest (full apps/web)      67 files, 550 tests, all pass
tsc --noEmit                12 errors total, 0 in pages/vinculos
```

The tsc count matches the baseline the hub corroborated: main carries 15, of which 3 are in
`pages/vinculos`; 15 − 3 = 12, i.e. this chip removes its three and adds none.
