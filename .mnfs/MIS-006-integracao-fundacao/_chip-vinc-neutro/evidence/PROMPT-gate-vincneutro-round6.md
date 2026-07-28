# Gate review — CHIP-VINC-NEUTRO, round 6 (frozen)

You are a cold, independent gate seat. You have not seen this chip's conversation and you will
not get it. Another seat is reviewing the same frozen input, blind to you and you to it. Do not
try to guess what it found.

**Code under review: `4c000a04..7b5c18eb`** — four commits
(`git log --oneline 4c000a04..7b5c18eb | wc -l` = 4). That range is the CODE delta, and it is the
range the hub's executor seat measured.

**Read the pack at the current `HEAD`, not at `7b5c18eb`.** Commits after `7b5c18eb` touch
`.mnfs/` only — they carry corrections the executor seat forced into the evidence, including one
that restates a count the pack originally got wrong. The invariant that makes this safe is
checkable in one command, and you should run it rather than trust this sentence:

```
git diff --stat 7b5c18eb..HEAD -- apps/     →  empty
```

If that is not empty, stop and say so: the brief is lying to you about what you are reviewing.

```
6e0e9ea4  round 5 — both gate sides REFUTED; make the impossible fixture unwriteable
2e5331b6  the guard had the defect it was built to catch; + the 29th fixture
ea856c32  run §11 count reconciliation against my own pack, before freezing
7b5c18eb  the sentinel had the shape it names — hub executor seat finding
```

(`4c000a04` is the base, so the diff is the four commits after it.)

## Scope of the chip

`/vinculos` front end only. Write-set: `apps/web/src/pages/vinculos/`. Zero Go, zero migrations,
zero `contracts/`, zero `packages/sdk-runtime/`. Go and SDK files are READ as authority, never
edited. `useErpImports.ts` belongs to the hub.

## What is NOT yours, and why

**Do not spend findings on execution criteria.** Lanes, must-fails, the merge check and the `git`
facts are discharged by the hub's executor seat — a third seat with a shell, independent of the
implementer, which ran them in its own detached worktree. Reported there, not here:

- `tsc` 12 errors, 0 under `pages/vinculos/`; `vitest` 64 files / 531 tests green.
- Merge onto `main` @ `0cb6d7e`: 0 conflicts; merged tree 67 files / 544 tests green.
- Guard must-fail, four mutation arms, each with its own count (see the pack).

A finding whose content is *"I could not run X"* has found nothing here. If a claim needs
execution to settle, **name it and route it to the executor seat** — that is a complete and
useful answer, and it costs you nothing.

What IS yours: whether the code and the pack say true things, and whether the checks in the delta
verify what they claim to verify.

## MANDATORY SECTION 1 — SWEEP

A verdict without a `## SWEEP` section is returned as **incomplete** and is not read as PASS.

Sweep the **class**, never re-examine the site. Each class below is given as searchable tokens,
not as a location, deliberately: a brief that names a site teaches the seat to stop at it. Search
the whole of `apps/web/src/pages/vinculos/` for each, and report what you found INCLUDING zero.

| # | class | tokens to search |
|---|---|---|
| C1 | a union enumerated by string literal — type-correct, therefore silent when the union grows | `=== "`, `!== "`, `.filter(`, `? (` on a field whose type is a union, `switch (` |
| C2 | a check WIDER than the fact it claims to verify (this round's defect) | `.toContain(`, `.includes(`, `.startsWith(`, `.match(`, `expect(` on a string where the subject is a symbol |
| C3 | an assertion that PASSES on empty extraction | `matchAll(`, `.match(`, `Array.from(`, `.map(`, any `.length` NOT paired with a positive lower bound |
| C4 | a fixture not built through the throwing constructor | `candidate_id:`, `reasons: [`, `confidence_band:` |
| C5 | a comment or pack sentence stating a repo fact that was never measured | 7-to-40 hex SHAs, file paths, `main`, `deleted`, `no longer` |
| C6 | a `Record<…>` or map indexed by a value that can be outside the key type at runtime | `Record<`, `[` indexing on a wire-typed value |

For each class, state: how many sites the token search returned (population), how many you judged
(extraction), and the verdict per site. **Zero findings in a class is a result — write it.**

## MANDATORY SECTION 2 — COUNT RECONCILIATION

Binding on you exactly as it binds the chip (`docs/HARNESS-PROFILE.md` §11, amendment `0cb6d7e`).

Any sweep you offer must print **two** counts: the **population** (the loose anchor) and the
**extraction** (what your pattern actually yielded). Different without a declared reason means the
sweep is reporting its own blind spot, and it will be read that way.

Also: **a pattern you never show matching is not known to match anything.** If a pattern returned
zero, demonstrate it matching a known member — or say plainly that you did not, and mark the
result unverified rather than clean.

This clause has already caught three real defects in this chip's own sweeps, one of them in the
tool built to enforce the clause. Assume it will catch yours.

## MANDATORY SECTION 3 — pack custody

**Only if you have a shell.** One of the two seats does; the other is physically read-only
(`Read`/`Grep`/`Glob`, no shell at all). If you are the read-only seat, write exactly *"no shell —
custody routed to the executor seat"* and move on. That is a complete answer and costs you
nothing. Do not infer custody from what you can read: a file you can open through `Read` looks
identical whether or not `git` has ever heard of it, which is the whole failure mode below.

If you do have a shell, report the output of both, verbatim:

```
git ls-tree -r HEAD -- .mnfs/MIS-006-integracao-fundacao/_chip-vinc-neutro/ | wc -l
git status --porcelain .mnfs/MIS-006-integracao-fundacao/_chip-vinc-neutro/
```

The `-r` is load-bearing and it is not a style preference: without it `ls-tree` lists tree
ENTRIES, so a pack whose files sit in subdirectories reports a single-digit number that looks
like a file count and is not one. This brief shipped the non-recursive form for one commit. Same
class as C2 in the sweep table — a measurement wider or narrower than the fact it names — which
is why you are being told rather than quietly handed the corrected command.

The second must be empty. An evidence pack that lives only in a working tree is one worktree
teardown away from never having existed — a prior chip in this mission lost six commits' worth of
artifacts to exactly that, and it is why this clause exists. If the pack is dirty, that is a
finding on its own, regardless of what the code says.

This clause is ALSO held by the hub's executor seat, so it is not discharged by your silence and
it is not discharged by your saying you cannot run it. Both of those are correct answers here
precisely because a third seat with a shell owns the same clause.

## MANDATORY SECTION 4 — the two adversarial questions

Answer both explicitly. "No finding" is an acceptable answer; silence is not.

**Q1.** The 29th fixture (`VinculosPage.test.tsx`) had its `match_status` changed from `"REVIEW"`
to `"ACCEPT"` in order to become producible. The test hosting it is named *"renders tabs and shows
the queue KPIs once loaded"* and asserts the two tabs, the `Vincular` button, and the text `MLB1`
— no KPI. Does the button's label or presence depend on `match_status`? And does a test whose
fixture moved status buckets still prove what it proved before? (The name/coverage mismatch
predates this chip; the fixture now travels a different code path, which does not.)

**Q2.** `wireFixtures.ts` closes each candidate's `reasons` array by mirroring what the Go
finalizer emits. That mirror is FE code. **Which test fails if the mirror diverges from the
finalizer while the anchor vocabulary stays unchanged?** If the answer is none, that is a named
gap, not necessarily a defect — but the answer must be written down either way.

## Authority order

`ARCHITECTURE.md`/ADRs → OpenAPI + `packages/sdk-runtime/` → `contracts/governance/` → wiki →
`.mnfs/` → tests/builds/commits. Go sources under `apps/server_core/internal/modules/` are the
authority for generator behaviour; the pack's claims about them are claims, not facts, until you
read the file.

## Verdict format

Open with **APPROVED** or **REFUTED** on its own line. Then:

- every finding with `file:line` and the string you read, quoted;
- severity, and what breaks in the running product if it ships;
- the `## SWEEP` section (section 1) and the reconciliation counts (section 2);
- section 3 — the custody output verbatim, or the one-line routing sentence;
- Q1 and Q2 answered;
- anything you could not settle without execution, named and routed.

Do not soften a finding because the pack sounds thorough. Four of the five previous rounds found
real defects in artifacts whose prose was confident, and the fifth found one inside the mechanism
built to prevent them.
