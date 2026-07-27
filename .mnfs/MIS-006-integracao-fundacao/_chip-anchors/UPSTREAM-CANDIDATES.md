# Upstream amendment candidates — CHIP-ANCHORS, consolidated for the CLOSED payload

Field findings from this chip that are about the METHOD, not about this repo's code. Raised as
candidates for the hub to accept, amend or drop — none is claimed as doctrine.

## 1. The differential probe — candidate for HARNESS-CORE §5 (verification technique)

**Shape.** Hold the fixture, the probe and the command FIXED. Swap exactly ONE file between two
commits. Run twice. The delta is attributable to that file alone.

**Why it earned a place.** It decided two questions on this chip that argument could not:
- whether the governance lane was already red at base (it was — 53 findings, identical line for line);
- whether S8 regressed the `title FOR`+`AGAINST` path (it did, while S8's own report read clean).

**Two preconditions, both hit here, both worth stating in the doctrine:**
- swapping a production file alone can break the build when the SIBLING TEST file references a symbol
  introduced in the same commit — swap the test file to the same commit, or the probe reports a
  compile error instead of a behavioural delta;
- the probe must assert its own non-vacuity FIRST. My first run PASSED with no output because the
  fixture produced zero candidates and the loop asserted nothing. `if len(x) == 0 { t.Fatal }` before
  reading anything.

## 2. Worker reports are evidence of a claim, not of the behaviour — candidate for §4

Two slices on this chip reported clean while the code was defective (S8's regression; S3 earlier).
The rule that actually caught both: the orchestrator re-runs the single highest-value must-fail
mutation BY HAND rather than accepting the report, and the pack labels each proof
`chip-re-run` vs `worker-observed`. Cheap — one mutation, one revert — and it has a 2/2 hit rate here.

## 3. `-count` on ordering guards must be justified, not conventional

S9's map-iteration mutation reddened on **run 7 of 10**. A `-count=5` lane — the convention used on
the previous slice — would have passed a genuinely non-deterministic ordering on persisted rows.
Candidate rule: any guard pinning ORDER of persisted data runs at `-count=10` minimum, and a lane
that passes at a lower count is not evidence.

## 4. A ruling's WORDING can be the defect — candidate for §4 escalation guidance

R-6 said "at most one reason per anchor"; its rationale said "two CONTRADICTORY reasons". The slice
implemented the sentence, correctly, and destroyed information. The chain the hub itself traced was
**ruling → card → slice**, not worker error. Two things worked and are worth encoding:
- the chip escalated instead of deciding, per "disagreement = BLOCKED with evidence";
- the corrective card carried an EXHAUSTIVE collision matrix (11 rows) precisely because the two
  previous cards had each left one case undefined. S9 reported the matrix total and found no
  uncovered case — the first slice on this function to come back without a defect.

Candidate wording: when a ruling is implemented literally and the result is worse, the defect is
presumed to be in the wording until shown otherwise, and the corrective card must enumerate the case
space rather than restate the principle.

## 4b. A pin that names one axis and silently drops the other — the general shape

Candidates 3 (below) and 5 are the same defect from two angles, and the hub asked for the general
formulation to be kept rather than the two instances:

**A conclusion surviving a dead premise is the failure mode that goes unnoticed *because the result
was right*.**

Two instances, both real on this chip:

- **A measurement names a SHA, or it names nothing.** The governance differential was measured at
  `8e37958a` and then defended with a `git diff` proving synchrony between two *pack* commits — the
  wrong axis. The hub's rule: a hub-run lane rung re-measures at the final tip BY DEFAULT, and the
  only valid window is `<sha measured>..<sha to merge>`. When the re-run finally happened it
  *confirmed* the withdrawn claim, which is exactly why this needs writing down: the claim was still
  never evidence, and a confirming result is the worst possible moment to notice that.
- **A gate brief names a ruling, or it names nothing.** Round-2's briefs were re-pinned to the new
  SHA and left carrying the superseded R-6 text. Dispatched unpatched they would have failed the chip
  for the behaviour R-6a mandates.

Same shape both times: the pin names one axis (the SHA) and drops the other (the ruling, the range)
in silence. Candidate wording: **any re-pin states every axis it is valid along, and a measurement
carries the range it was taken over.**

## 5. Gate briefs must be re-pinned to the ruling in force, not just the SHA

Both round-2 briefs still stated R-6 after R-6a superseded it. Dispatched unpatched, both reviewers
would have failed the chip for keeping same-anchor `FOR`+`AGAINST` — the exact behaviour now
mandated. Candidate checklist item: re-dispatching a gate after a corrective re-pins BOTH the tip
AND every ruling the corrective changed.

## 6. `codex exec` trust check binds to shell cwd — profile §3 addendum

`--sandbox workspace-write` dispatched from the SCRATCHPAD dies instantly with
`Not inside a trusted directory and --skip-git-repo-check was not specified` and writes a zero-byte
`.last.md`. The `.done` sentinel still fires, so a poller sees a completed dispatch with an empty
verdict. Already recorded in memory as "codex dispatch cwd worktree"; worth promoting into profile
§3 next to the other sandbox false alarms, with the detail that the FAILURE IS SILENT to sentinel
polling — check the log head, not just the sentinel.

## 8. Gate output is a claim, not a fact — candidate #2 extended to REVIEWERS

Candidate 2 says a worker's report is evidence of a claim, not of the behaviour. Round 4 proved the
same holds for a GATE. The round-3 GPT reviewer's finding named
`auto_link_policy_test.go:273-303`; the test starts at **274**. I copied the reviewer's line number
into the pack unverified — inside the very corrective whose subject was unverified line numbers — and
the mechanical audit caught it. The reviewer's SUBSTANCE was right (the colour branch is tested, my
sweep had said otherwise); only its coordinate was wrong, which is exactly the shape that survives
inspection.

Candidate wording: **a citation adopted from a reviewer is re-derived like any other before it enters
the pack.** Authority over the verdict is not authority over the line number.

## 9. A document must not cite line numbers of a report generated from itself

The R-11 audit reads `EVIDENCE.md` and prints `pack:<line>` rows. The paragraph explaining the
audit's flagged rows first identified them by those `pack:` numbers — and writing that paragraph
shifted them. **There is no fixed point:** any edit that references the report's coordinates
invalidates the coordinates. Only content-addressed references (the citation text itself, the symbol,
the section heading) survive.

Narrow but sharp, and it generalises to any generated-from-source artifact a pack quotes: coverage
reports, lint output, diff line numbers. Candidate wording: **a pack references generated output by
CONTENT, never by that output's own line numbers.**

## 10. What a mechanical citation audit does NOT buy — state it or it will be over-trusted

`cite-audit.py` proves a citation RESOLVES and shows what the line contains. It cannot tell whether
the prose around it is TRUE. The round-3 defect — "this branch has NO test at all" — would pass the
audit cleanly, because the neighbouring citation resolved perfectly and the false claim was about
absence elsewhere. If R-11's tooling is promoted upstream, the boundary ships with it, and the cold
reviewer's brief is re-pointed at SEMANTICS rather than left to sample line numbers the script has
already covered. Otherwise the tool's green output becomes the new unearned trust.

## 11. A pack is only citable against a FROZEN code tip — ACCEPTED by the hub (R-19)

Every citation in an evidence pack is a coordinate into code. A corrective that inserts lines above
one invalidates it silently, so the ONLY safe order is: land every code corrective, freeze the code
tip, then generate the citations. Round 5 of this chip exists because the remedy was executed in the
wrong order — "final tip" was ASSERTED rather than verified, a 27-line test insertion landed after
it, and 40+ citations decayed in one commit.

The hub's wording: *the remedy was never wrong. Its precondition never held.* And the partition that
follows from it — **pack-derived facts get a pointer, code-derived facts get sequencing.** They are
not competing options; they cover different halves. A count copied out of a report generated FROM the
pack cannot be fixed by sequencing, because every edit moves it; a line number into code cannot be
fixed by a pointer, because the pack has to name a location.

Candidate wording for §4: **a pack is citable only against a frozen code tip. Freeze, then cite; and
never copy a derived number out of an artifact the pack itself generates.**

## 12. Qualifying an under-specified citation is an INFERENCE, and inferences must print

Removing bare `:NNN` citations means giving each one a file, inferred from context. On this chip that
inference was WRONG in six places and every wrong result resolved to a real, plausible line — four
`EVIDENCE.md` line numbers became `link_candidate_repo.go` citations pointing at genuine SQL, and
three `generation_service.go` lines were attributed to `marketplace_capability_service.go`, where
`:76` sits inside a different function in a different module and reads perfectly.

None was findable by reading. All six fell out of one pass of an audit that prints what EVERY
citation resolved to, alongside the file it inferred and the distance it inherited over. That is the
generalisable part: **a migration that fills in missing information must emit the inference it made,
not only the result** — otherwise the migration silently manufactures the confidence the original
gap at least advertised.

Ratified upstream already as R-20 (an inference must be visible as an inference; the permissive pass
is a migration tool, not a permanent crutch) and R-21 (a citation resolving to plausible-but-wrong
content is invisible to reading, so print every resolution). Recorded here because the *technique*
generalises past citations: any bulk qualification, normalisation or backfill of under-specified
references has this shape.

## 7. Deferred findings this chip is carrying out, not fixing

- **F-01 / the `pol` prefix** (A10): reported, deliberately not fixed in-chip.
- **C2 divergence**: the criterion's prose says no anchor name is hardcoded in the generator; its
  *prova mínima* names only `marca`/`refforn`, while `seller_sku`/`ean`/`title` remain hardcoded
  reason seeds. Both gates passed C2 in round 1; the chip found the divergence afterwards, declared
  it, and referred it rather than fixing it under a criterion that does not clearly ask for it.
- **LIVE rung (U1–U3)**: recorded OPEN per R-8, hub-run.
- **Integration duplicate-anchor assertion**: ADDED but NOT RUN — no database, dev stack is a hub
  seam this chip may not boot.
