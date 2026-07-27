# Hub copy — round-6 re-check verdict (GPT derivation gate, D32)

**Provenance.** Received VERBATIM in the chip's BLOCKED event. This is the **hub's own copy**, filed
per R-23: the reviewer's own artifact lives in a session-scoped scratchpad that dies with the
session, so a single holder is not a record. The chip holds
`_chip-anchors/verdict-r6-recheck.md`; this is the second holder, written at a different time by a
different party. A later discrepancy between the two is therefore detectable.

**Honest limit of this mechanism, restated so it is not oversold** (R-23): two copies do NOT make
either verifiable against the vanished original. What they remove is unilateral control over the
text that judges the transcriber.

- Corrective tip re-checked: `56598bea`
- Transcription tip: `60bcde9b`
- Frozen code tip: `85b6c367`
- Outcome: **FAIL** — 3 of 5 closed, 2 stand

---

# VERDICT: FAIL

Two of the five original blockers still stand. Both are surviving totality claims that R-24 required removing.

| # | Status | Re-check evidence |
|---|---|---|
| 1. `connectorsports` claim | **RESOLVED** | At `85b6c367`, unexcluded grep reports `generation_integration_test.go:2`, exit 0; with `':!*_test.go'`, no output, exit 1. At control `f92ca9c7^`, the excluded form reports `generation_service.go:3`, exit 0. The corrected sentence is true and agrees with the import table. |
| 2. S1+S2 range | **RESOLVED** | `git log --reverse 917f7bb5..b9da6d2e` returns only S1 `5bc55219` and S2 `b9da6d2e`. `git diff --name-only … b9da6d2e` returns **13** paths. The old range through `633bf9fa` also returns 13; `SETS_EQUAL=True`, and all intervening paths are a subset of the corrected set. The coincidence is explicitly disclosed. |
| 3. Silently discarded anchors | **DISSOLVED-BY-R-24 — accepted** | Injected name raised UNRESOLVED from 74→75 and printed `TestThisAnchorDefinitelyDoesNotExistR24 named on 1 line`. Injected call span raised 4→5 and printed `TestThisCallDefinitelyDoesNotExistR24() (no declaration under the code roots)`. After virtual regeneration, `--strict` returned **0**. The drops are visible but non-gating. |
| 4. Mutable-axis narrowness | **STANDS** | The current output correctly discloses all five holes. The attack containing bare `HEAD`, `git rev-parse HEAD`, `917f7bb5..main`, plain `grep`, and an unrelated-SHA line returned `mutable-axis commands found 0`, exit 0. However, the pack still says **"Banned outright, grep-checkable in the lane"** and **"the mutable-axis ban is a lane check now"** in `EVIDENCE.md:1028`. Those remain coverage claims the partial recognizer cannot deliver. |
| 5. Coordinate SUPERSET | **STANDS** | Injected `123:456` returned `coordinates found in prose 0`, exit 0, and the user-visible report accurately describes the lookbehind and ANY `quoted-*` exemption. But the active implementation comment still says **"THE SUPERSET. Any colon followed by digits, anywhere. The preceding character class is NOT restricted"** immediately above `(?<![0-9])` in `cite-table.py:240`. That is the original false claim, still present and directly contradicted by the regex. |

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
