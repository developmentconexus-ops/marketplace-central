# REQUEST — restart backend at c89fae3d for P7 browser QA

```yaml
from: chip M-01-listings-read-spine
to: HUB (local_efa46c30-1c0c-4075-9671-c2d7ae9efabe)
event: REQUEST
tip: c89fae3d   # code tip is 982d44e; c89fae3d is evidence-only on top
pushed: no
blocking: yes — P7 cannot start against a stale binary
```

## What happened since ENV-READY at a6878dc6

Gate round 4 ran at `a6878dc6` (cold Opus + Sol medium, simultaneous, fixed SHA) and **merged to
FAIL** on one `important`: **H1**. Sol raised it; Opus did not examine it.

`read_service.go` had 15 `slog` sites. Fourteen named the read path via `op`; one did not — the
cost-fact propagate ERROR inside `enrich`. `enrich` is shared by List/ByProduct/Get and could not
name its caller, so a propagated cost-reader failure reached the operator attributable to no read
path. `Summary`'s inline cost path logs the same message on the same condition *with* `op`, which is
what made this a gap rather than a convention.

**Slice 12** (`982d44e`) threads `op` through `enrich`/`enrichGroups` — both private, single-file, no
external callers. Telemetry attribution only: no behavior change, no added log, no test (a `slog`
field is not behavior a unit test should pin here; asserting the log rather than the read is test
theater — both round-5 reviewers independently endorsed that call).

**Gate round 5** at `982d44e`: **both sides PASS, no contradiction.** H1 resolved; zero blocking,
zero important, one non-gating nit. Full merge table: `_gate-evidence/round-5/dual-gate-round5-verdict.md`.

**Code criteria for M-01 are met. Only QA passes a milestone, so this is not a close.**

## The request

Restart the backend at **`c89fae3d`** (code tip `982d44e`) so P7 browser QA exercises the code that
actually gated. Per the harness the chip never boots a server, binds `:8080`, or loads `.env` — the
dev stack is a hub seam, so I am asking rather than doing.

Expected on boot: **`applied 0 migration(s)`** — slice 12 carries no migration, and neither did
slice 11. If migrations apply, something is wrong; stop and tell me rather than proceeding.

Scope for your independent verification of `a6878dc6..982d44e`: exactly one Go file,
`apps/server_core/internal/modules/listings/application/read_service.go`, 7 hunks. `internal_read/`
untouched. No composition, transport, migration, OpenAPI, or SDK. `docker/dev/*.sh` show as modified
in the worktree — **those are pre-existing hub dev-stack changes, not mine**; I did not touch them
and deliberately excluded them from every commit.

## Self-reported, before you find it

The `slice12-L0-report.md` I wrote claimed 16 `slog` sites. There are 15. The cold Opus reviewer
caught my miscount while checking a claim I had invited it to distrust. H1 turns on the set of sites
*without* `op` being empty — which it is — so the substance stands, but the number was wrong. It is
corrected in place at `c89fae3d`, with the correction left visible rather than silently edited.

Worth one line for the board: round 4 failed on a finding Opus never looked at, and round 5 had Opus
catch an error in my own report. Neither side is reliably the deeper one. The dual gate earns its
cost precisely because which reviewer goes deeper is not predictable in advance.

## Still on your board from this chip

- **Task #2** — `internal_read` taxonomy refactor (the G1 residual, deferred-with-name per your
  ruling (C)): an adapter or data defect is still misclassified as source-unavailable because
  `safeOracleCause` flattens the cause before the application layer sees it. Misclassified but
  logged, no longer silent.
- **H4** (question, from round-4 Opus) — `read_service.go:414`: on a *ceiling* outage the cost fetch
  is skipped and `Cost` nils, so `cost: null` is served while the cost source is healthy. Arguably
  ADR-17 direction (b) one level up. Pre-base, pinned as intended, certified rounds 2-4, outside
  every slice brief. Raised for your queue, not fixed by me — reviews verify, they do not generate
  scope.
- **H5** (question, from round-4 Opus) — `internal_read/adapters/cache/cache.go:179-197`: `group.DoChan`
  shares one loader's result, so caller A's `context.Canceled` can be delivered to an uncancelled
  caller B. Honest either way (no false fact served); the delta did not create it, but slice 11
  changed B's outcome at a distance from degrade to propagate.

Awaiting ENV-READY. No push.
