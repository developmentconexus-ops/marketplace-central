# Round-6 re-check — the gate that filed the FAIL, on its own five

This is not a seventh round. The hub scoped it: *"Whoever filed the FAIL re-checks its own five
blockers at the corrected tip ... by the gate that delivered it, scoped to exactly what it filed."*
The cold gate did not re-run — it passed round 6, and the corrective touches no code, which the chip
verified rather than asserted (`git diff --name-only 85b6c367 56598bea -- apps contracts packages`
returns nothing).

- Dispatch: ledger row **D32**, `gpt-5.6-sol` / `medium`, OS process, `--sandbox read-only`.
- Corrected pack tip re-checked: `56598bea`. Frozen code tip, unchanged: `85b6c367`.
- Prompt: `scratchpad/prompt-p6-sol-recheck.md`. Artifacts: `agent__p6-sol-recheck.last.md` + `.log`.

**Outcome at `56598bea`: FAIL.** Three of the five closed — two RESOLVED on the facts, one
DISSOLVED-BY-R-24 and accepted as such. Two stood, and both for the same reason: the corrective removed
the totality claims from what the tool *prints* and left them standing in what the pack and the
implementation *say*. Three holders of the same claim, one of them reached.

The chip verified all three lines itself before relaying, rather than forwarding the verdict as a
postbox. They were as reported.

**Both were then closed by DELETION under R-25, and the hub verified their absence itself.** The exact
strings are quoted inside the fence below and are deliberately not restated in this preamble: a false
sentence repeated as a live description is the thing R-25 removed. R-25 also fixed the instrument —
*honest-unknown is for gaps, not for falsehoods; you disclose what you do not know, you DELETE what is
wrong* — and ruled that §11 did not fire, because the gate's own observation (*"the two blocking
survivor lines already existed in `f1397cf7^`; they were omitted by the corrective rather than
introduced by it"*) makes this incomplete application rather than recurrence. The two want opposite
responses: recurrence wants a better remedy, incompleteness wants the remedy finished.

Transcribed verbatim below. The chip wrote this file; the gate wrote the text inside the fence, and
none of it was edited to read better. Its own stated limit is the last section: it re-checked five
findings and nothing else — no semantics, no tests, no build, no vet, no Go tests, no governance, no
live drive, no UI, and not the hub-side copies R-23 requires.

Note the FOREIGN-fence SHOULD-FIX below survives on purpose: it was not in the corrective's scope, and
the gate says so plainly instead of counting it against the corrective.

`````quoted-verdict
# VERDICT: FAIL

Two of the five original blockers still stand. Both are surviving totality claims that R-24 required removing.

| # | Status | Re-check evidence |
|---|---|---|
| 1. `connectorsports` claim | **RESOLVED** | At `85b6c367`, unexcluded grep reports `generation_integration_test.go:2`, exit 0; with `':!*_test.go'`, no output, exit 1. At control `f92ca9c7^`, the excluded form reports `generation_service.go:3`, exit 0. The corrected sentence is true and agrees with the import table. |
| 2. S1+S2 range | **RESOLVED** | `git log --reverse 917f7bb5..b9da6d2e` returns only S1 `5bc55219` and S2 `b9da6d2e`. `git diff --name-only … b9da6d2e` returns **13** paths. The old range through `633bf9fa` also returns 13; `SETS_EQUAL=True`, and all intervening paths are a subset of the corrected set. The coincidence is explicitly disclosed. |
| 3. Silently discarded anchors | **DISSOLVED-BY-R-24 — accepted** | Injected name raised UNRESOLVED from 74→75 and printed `TestThisAnchorDefinitelyDoesNotExistR24 named on 1 line`. Injected call span raised 4→5 and printed `TestThisCallDefinitelyDoesNotExistR24() (no declaration under the code roots)`. After virtual regeneration, `--strict` returned **0**. The drops are visible but non-gating. |
| 4. Mutable-axis narrowness | **STANDS** | The current output correctly discloses all five holes. The attack containing bare `HEAD`, `git rev-parse HEAD`, `917f7bb5..main`, plain `grep`, and an unrelated-SHA line returned `mutable-axis commands found 0`, exit 0. However, the pack still says **“Banned outright, grep-checkable in the lane”** and **“the mutable-axis ban is a lane check now”** in [EVIDENCE.md](C:/Users/leandro.theodoro/Documents/marketplace-central/.claude/worktrees/chip-anchors/.mnfs/MIS-006-integracao-fundacao/_chip-anchors/EVIDENCE.md:1028). Those remain coverage claims the partial recognizer cannot deliver. |
| 5. Coordinate SUPERSET | **STANDS** | Injected `123:456` returned `coordinates found in prose 0`, exit 0, and the user-visible report accurately describes the lookbehind and ANY `quoted-*` exemption. But the active implementation comment still says **“THE SUPERSET. Any colon followed by digits, anywhere. The preceding character class is NOT restricted”** immediately above `(?<![0-9])` in [cite-table.py](C:/Users/leandro.theodoro/Documents/marketplace-central/.claude/worktrees/chip-anchors/.mnfs/MIS-006-integracao-fundacao/_chip-anchors/cite-table.py:240). That is the original false claim, still present and directly contradicted by the regex. |

## Merge conditions

- `git diff --name-only 85b6c367 56598bea -- apps contracts packages` returned nothing. `CODE_DIFF_COUNT=0`.
- Baseline `python cite-table.py --strict` returned exit 0 and identified itself as `NOT A LADDER RUNG`.
- The ladder table has no PACK row.
- Sixth-row replay passed exactly as claimed: name and call-span surfaced; regenerated `--strict` exited 0.
- Probes ran against in-memory byte copies because the worktree is read-only. Actual files were never mutated. Before/after SHA-256 values matched; `BYTE_IDENTICAL_PACK=True`, `BYTE_IDENTICAL_TABLE=True`.
- Final status remains only the pre-existing `M package-lock.json`.

## SHOULD-FIX — FOREIGN fence tags

**Still present; corrective did not fix it.**

Control injection:

```text
OUTSIDE_FOREIGN EXIT=1 | mutable-axis commands found 1
INSIDE_FOREIGN  EXIT=0 | mutable-axis commands found 0 | quoted-verdict 10→11 lines
```

A pack author can still hide authored content by labeling it `quoted-verdict`; provenance is not verified. This is not a new defect.

## New defects introduced by the corrective

None found. The two blocking survivor lines already existed in `f1397cf7^`; they were omitted by the corrective rather than introduced by it.

## Honest limit

I did not re-review F-01/F-02/F-03 semantics, tests, production behavior, builds, vet, Go tests, governance, live integration, UI, the external hub copy required by R-23, or the hub-side U2 artifact. I did not report the two mission-backlog items.
`````
