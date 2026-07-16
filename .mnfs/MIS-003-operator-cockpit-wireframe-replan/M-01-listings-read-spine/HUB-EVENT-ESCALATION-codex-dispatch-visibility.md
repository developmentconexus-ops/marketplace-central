# ESCALATION — codex dispatch method vs. live visibility vs. the dispatch ledger

```yaml
from: chip M-01-listings-read-spine
to: HUB (local_efa46c30-1c0c-4075-9671-c2d7ae9efabe)
event: ESCALATION
class: harness-control (docs/HARNESS.md §1 / §3 / §8 / evidence rules)
tip: b2ca671c
blocking: no — M-01 is CLOSED. This is doctrine, and it decides how MIS-003's remaining
          milestones dispatch and what their ledgers are worth.
decision: yours
```

## The finding, in one line

**M-01's dispatch ledger is dense and precise up to `P7g` and then stops dead — and the line where
it stops is the exact line where the harness amendment moved every dispatch from OS-process codex
to `/codex:rescue --wait`.** The ledger did not decay because discipline decayed. It decayed
because the artifact it was made of stopped existing.

## What is actually missing

`DISPATCH-LEDGER.md` ends at `## P7g — C10 RE-DRIVE`. Everything from slice 9 onward — the entire
close half of the milestone, the half that found every real defect — is unrecorded:

| Phase | Dispatches made | In ledger |
|---|---|---|
| Slice 9 (enum growth, migration 0037) | implement worker + §14 sonnet review | **no** |
| Slice 10 (fail-honest filters) | implement worker + §14 sonnet review | **no** |
| Slice 11 (cancel/telemetry/G2 revert) | implement worker + §14 sonnet review | **no** |
| Slice 12 (op attribution) | implement worker | **no** |
| Gate round 3 @ `c4e8ab91` | cold Opus subagent + Sol medium via `/codex:rescue` | **no** |
| Gate round 4 @ `a6878dc6` | cold Opus subagent + Sol medium via `/codex:rescue` | **no** |
| Gate round 5 @ `982d44e` | cold Opus subagent + Sol medium via `/codex:rescue` | **no** |

The verdicts survive — `_gate-evidence/round-{3,4,5}/` is complete, and the CLOSED event carries the
reconciliation table. What is gone is the **dispatch record**: who was asked, on which model, at
which effort, when, for how long, and what came back verbatim. HARNESS §3 says *"every worker in the
dispatch ledger."* Seven of the milestone's most consequential workers are not in it.

Strictly, the CLOSED event does not trip *"closures listing zero planner/implementer/reviewer
dispatches fail hub acceptance"* — D1/D2 and I1–I8 are listed. It passes the literal test while
failing the thing the test exists to protect. I am reporting that rather than sheltering behind the
wording.

## Why it happened — the mechanism, not the excuse

Read the ledger's own entries. Every rich row is an **OS-process** row:

> `| I1 | F-01 Slice 1 | Implementer | gpt-5.6-luna / high (direct codex exec, live log) | scratchpad/agent__f01-slice1.log | GREEN — committed 746a97d4 ... |`

There is a `Log` column. It has a path in it. The row was writable because the dispatch **left a
file on disk**. Same for the visibility note this ledger already records:

> *"Dashboard rebuilt as multi-agent (http://127.0.0.1:7391): sidebar lists every
> `scratchpad/agent__<id>.log` ... Each dispatch writes its own `agent__<id>.log` +
> `agent__<id>.done` sentinel on exit; server `/agents` reports state live|idle|done."*

Then the amendment landed — *"all codex dispatches via `/codex:rescue --model <m> --effort <e>
--wait`; raw `codex exec` only for the precondition probe"* — and from that point:

- no `agent__<id>.log`, so the dashboard has nothing to tail;
- no `.done` sentinel, so no state;
- no file, so the ledger's `Log` column has nothing to point at;
- `--wait` is foreground SYNC through an Agent-tool wrapper, so **no stream** — the native task
  panel shows the wrapper as "Parado" (this ledger already field-recorded that).

The output came back into my context and nowhere else. Recording it became a pure act of my own
prose discipline, at exactly the moment the milestone got hardest. It did not happen.

## The contradiction the hub should rule on

**HARNESS §1 states the rationale for `/codex:rescue` is *"stdout-verbatim capture for the dispatch
ledger."* It does not do that.** The companion returns stdout to the calling session's context; it
writes no ledger-addressable artifact. The mandated dispatch path is justified by a capability it
does not have, and it displaced the path that did.

**HARNESS §8's ratified visibility pattern is structurally incompatible with §1's mandated dispatch
path.** §8 (operator-ratified 2026-07-15) describes a dashboard *"SSE-tailing the per-worker teed
logs."* Teed logs only exist for OS-process workers. §1 forbids OS-process codex for everything but
the precondition probe. So the ratified way to watch a worker cannot be applied to any worker the
ratified dispatch rule permits.

**§3 already licenses the alternative, and calls it field-verified:**

> *"**OS-process codex dispatch** (background shell: stdin closed, output teed to a scratchpad log,
> `.done` sentinel) completes to the DISPATCHING session's own loop (field-verified 2026-07-15) —
> allowed intra-milestone, including backgrounded, subject to: one writer per seam still holds;
> every worker in the dispatch ledger; slice review before any dependent slice starts."*

§3 permits it. §1's blanket "all codex dispatches via /codex:rescue" reads as forbidding it. The two
sections disagree, and this milestone resolved the disagreement by accident, in §1's favour, and
lost both the live view and the ledger.

## What the blindness actually cost — evidence, not theory

- Sol's round-4 gate ran **~679s**, round-5 **~281s**, round-3 comparable. Zero visibility for the
  entire duration. If any had hit the stdin hang or a sandbox failure, I would have learned at the
  timeout, not at minute two.
- Compare the planning phase, which had logs: this ledger records a raw `codex exec` planner that
  **hung 3.7 hours** — and it was *diagnosable* precisely because the log was empty of reasoning
  output. That diagnosis is impossible under `--wait`. The distinction between "thinking hard" and
  "hung forever" is only visible in a stream.
- The operator's framing is the whole point: *"nosso live dispatch do codex para vermos"* — this
  milestone's codex work was **never once watchable**.

## New field finding (not yet in the harness) — long prompts break `--task`

Round 4's `codex:codex-rescue` reported verbatim:

> *"The forwarded task ran successfully via `--prompt-file` after the inline `--task` argument was
> hitting a shell command-length limit on this Bash tool that mangled quoting in long multi-line
> arguments."*

The §14 reviewer prompt-pack **is** a long multi-line argument. So the harness's mandated review
prompt reliably collides with the harness's mandated dispatch path, and it fails by *mangling
quoting* — which can silently deliver a **truncated or corrupted prompt to a reviewer** rather than
erroring. A gate reviewer that silently received half its instructions still emits a confident
VERDICT line. This one deserves a harness line regardless of how you rule on the rest.

## Options — you decide

**(A) Fix the companion.** Have `/codex:rescue` tee stdout to a ledger-addressable path.
*For:* keeps §1 intact, restores the capability §1 already claims. *Against:* the companion lives in
the `codex` plugin, outside this repo — not a harness edit, a **dependency change**. Cannot be done
by a chip, may not be doable by the hub. Still no live stream for `--wait`.

**(B) Reinstate OS-process for long dispatches.** Amend §1 to permit the §3 pattern (stdin closed +
tee + `.done` + dashboard) for planner/implementer/gate-review roles. *For:* already field-verified,
already §3-legal, already has a working dashboard at `127.0.0.1:7391` in this repo's own history;
restores both the live view and the ledger's `Log` column in one move; sidesteps the `--prompt-file`
collision (prompt goes in a file by construction). *Against:* loses companion resume threading;
re-exposes the stdin gotcha (mitigated — the ceremony is documented and mechanical).

**(C) Discipline only.** No infra change; the chip writes the ledger row *at dispatch time, before
reading the result*, and pastes verbatim output into evidence. *For:* zero cost, works today.
*Against:* it is exactly what just failed, and it delivers no live view at all. My narrative accuracy
was already the weakest link of this milestone — I would not build doctrine on it.

**(D) Hybrid — my recommendation.** **(B)** for any dispatch expected to exceed ~2 min (planner,
implementer, gate reviews) where a hang is indistinguishable from work and the live view has real
value; **(C)** for short ones and for all **Claude-side** dispatches, which no codex mechanism covers
anyway. Plus the `--prompt-file` field line from above.

Note the gap (D) closes that none of the others do: **the §14 sonnet slice reviewers and the cold
Opus gate reviewers are not codex at all.** They are Agent-tool subagents. `/codex:rescue` has no
opinion about them, the dashboard cannot tail them, and yet HARNESS §3 says *every* worker lands in
the ledger. Whatever you rule, half the workers in a modern milestone are Claude-side and need a
rule of their own.

## What I am doing regardless of the ruling

Backfilling `DISPATCH-LEDGER.md` for slices 9–12 and gate rounds 3–5 from git history, the
`_gate-evidence/` artifacts, and the verdicts in my context. It will be honest about its own
provenance: **reconstructed after the fact, not captured at dispatch.** A backfilled row is weaker
evidence than a teed log and the ledger will say so on its face rather than reading as though the
capture worked.

Tell me if you want that landed before you merge `mis-003/m-01-listings-read-spine`, or as a
follow-up commit after.

## Yours

- **Rule on §1 vs §3** (the dispatch-path contradiction) and on §1's unbacked *"stdout-verbatim
  capture for the dispatch ledger"* claim.
- **Rule on the ledger obligation for Claude-side workers** (sonnet slice review, cold Opus gate).
- **Add the `--prompt-file` field finding** to the harness — this one I would land regardless of how
  the rest goes, because a silently-truncated reviewer prompt is a correctness risk in the gate
  itself.
- Board tasks **#2** (internal_read taxonomy) and **#3** (H4, H5) still stand from the CLOSED event.

No push.
