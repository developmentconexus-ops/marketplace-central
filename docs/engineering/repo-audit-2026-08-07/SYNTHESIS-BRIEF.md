# Phase 2 — Synthesis brief

Method: `C:\Users\leandro.theodoro\Documents\MetalDocs\docs\engineering\repo-audit-playbook.md`, Phase 2.
Two arms run this brief independently with no visibility into each other. Divergence between them is
the point; do not try to guess what the other arm will say.

## Read first

1. `docs/engineering/repo-audit-2026-08-07/PHASE-0.md` — the operator's goal, hard constraints, and
   14 established facts you must not re-derive.
2. All ten lane reports in `docs/engineering/repo-audit-2026-08-07/lanes/`:
   `duplication`, `layering`, `cicd`, `delivery`, `persistence`, `testing`, `observability`,
   `security`, `goidiom`, `frontend`.

The lanes were discovery-only and were forbidden to judge. Judging is your job.

## The goal, in the operator's words

> I have run something really deep into MetalDocs and I am changing the way I code there to move to
> something more professional towards issues, PRs, PR review, CodeRabbit mechanical full validation
> and so much more. For that I had to identify every error in my code, my platform, to improve it and
> to create this full validation. I want to run it here as well so we move on the same path, this way
> it gets so much harder to send bad PRs.

The operative sentence is the last one. **The program succeeds if a mechanism makes a bad change hard
to land — not if the code gets cleaner.** An axis that fixes code and leaves the judging to discipline
has failed the goal even if every finding in it is closed.

Calibration: **solid professional level.** Not Google-tier.

## The inverted evidence standard — binding

For any **structural** conclusion (a boundary, an ownership, a target shape):

**Inadmissible as an argument FOR a structure:**
- that the code is currently organised that way
- that a prior decision document (ADR, `ARCHITECTURE.md`, a design doc) said so
- that migration would be expensive
- that the import graph, route topology, or schema reflects it
- doc-comments describing the current design

**Admissible:**
- the problem domain and any standards or regulation governing it
- how mature systems in the same field solve it
- design principle with a **named failure mode** (not "clean architecture says so")
- the product's observable user-facing behaviour
- measured cost and measured scale

Migration cost is admissible for **sequencing** and never as an argument that a target is right.
"Do A before B because B is cheaper afterwards" is valid. "Keep A because changing it is expensive"
is the trap.

This standard exists because of a real failure: two independent advisors both defended a structure
using only consequences of that structure, and both reversed when that evidence was ruled
inadmissible. The reasoning was correct everywhere else, which is what made the circular section hard
to see.

## The inversion test — mandatory, one line per structural conclusion

> State what would survive if the current implementation were the opposite in every respect.

A conclusion that cannot pass this is not a conclusion. Write the line down. Example of the form:

> *Boot-proven non-bypass DB identity: survives an opposite schema design, because row-level security
> is ineffective for a bypassing connection in every topology.*

## Control versus effect — binding

**A control that exists and does not fire is absent.** A gate script referenced by zero pipelines, a
type-checker wired to no command, a database policy inert against the connecting role, a negative
fixture no runner invokes. Counting these as "we have a control" is the same error class as the
circular argument — reasoning from an artifact instead of from an effect.

Where a control is inert, **size the axis as if the control were absent, because it is.**

## What to produce

Write to your assigned output file. Deliver all seven:

1. **Axes, not a list.** Group the findings into **6–10 axes**, where "connected" means *fixing them
   together is cheaper than separately, because they share a root cause, a mechanism, or a blast
   radius*. Name each axis for its **cause**, not its symptom — "the verifier is not one trusted
   product" beats "CI gaps". A finding belongs to exactly one axis. **If two axes want the same
   finding, they are one axis or the cause is wrong.**
2. **Kill the noise.** Which real findings do not earn a slot, and why — cost of fix versus cost of
   leaving it. Be willing to say "this is fine, leave it forever". A program that treats 150 findings
   as equal is how programs die.
3. **Consolidated do-not-touch list**, merged from every lane's "what is actually fine" section.
4. **The target gate topology.** Climb the firing hierarchy as high as each rule allows:
   1 unrepresentable → 2 boot-fatal → 3 red build → 4 runtime assertion → 5 discipline.
   **Level 5 is not a control.** State for each proposed gate: where it fires, what it blocks versus
   annotates, and its negative fixture. Every guard ships with an input that makes it fail, in the
   same change. Constraints that bind this: **GitHub Actions IS available — see the amended budget
   section of PHASE-0.md, which supersedes the zero-spend constraint**, so level 3 (red build) is
   reachable for any rule expressible as a command that exits non-zero, and proposing level 5 for
   such a rule is a defect in the recommendation; solo operator plus agents, so
   human PR review is **not** an available control while credential custody **is**; and the author of
   nearly all code is a machine whose failure mode is fluent, internally consistent wrongness — the
   wrong noun used correctly everywhere, including inside the guard written to catch it.
5. **The sequence.** Dependency order, rough sizes in days, and explicitly: which axis makes every
   later axis cheaper. Justify the first axis.
6. **Inversion tests.** One line per structural conclusion, per the section above.
7. **Where you disagree with a lane.** A finding that is wrong, mis-sized, or misattributed. Say so
   with evidence. The lanes measured; they did not always weigh correctly.

## Ground rules

- Read-only. Do not edit any file except your own output file. Do not run gates or test suites —
  another session holds write access to parts of this checkout and your runs would collide.
- Cite `file:line` from the lane reports. Where you make a claim no lane made, mark it clearly as
  yours and say what would confirm it.
- Do not re-derive PHASE-0's 14 established facts. Contradicting one with better evidence is a
  finding — say so loudly.
- If a lane's number and another lane's number disagree, do not average them. Name the disagreement.
