# Anti-patterns — defects this repo actually shipped

Each entry is a real case. Read this last, as an adversarial pass over a plan that feels
finished: for every entry, ask *"is my plan doing this?"*

The pattern behind most of them: **an observable that is stable, cheap, and wrong about which
world you are in.** It reproduces, so it feels like evidence; it passes in both the fixed and
the broken world, so it proves nothing.

---

## Planning

**The plan built a second engine beside the working one.** A 17-task fiscal plan constructed a
calculation layer next to `pricing/domain`, which already produced the right number. The
column on screen was empty, and empty column read as "feature absent" instead of "feature not
wired". Re-scoped to 7 tasks. → *An empty screen is not evidence of missing code. Trace the
value from the renderer backwards before planning to produce it.*

**The plan named a column that does not exist.** A query was written against
`sync_state.last_success_at`; the table has `last_full_sync_at` and `last_incremental_at`.
Caught by running the query, not by reading it. → *Measure schema with `\d`, paste the output
into the plan.*

**The plan hand-rolled a scheduler.** A market-collection ticker was specified beside
`syncapp.Scheduler`, discarding `sync_state` reconciliation, per-entity failure isolation,
cursor persistence and — the decisive one — `RecordFailure` feeding the operator's sync-health
card. A broken collection would have shown a green screen. → *Check whether the seam already
exists; check what reuse buys you that a parallel build does not.*

**The brief was an allegation and it rotted.** Three claims in one milestone brief were false
by the time a worker read them. Line anchors measured in one tree pointed at different code in
another. → *Anchor by content — route, `operationId`, schema name, function name. When a line
is unavoidable, name the tree.*

**A declared risk was treated as acceptable debt.** Three declared risks were rejected by the
operator; seven measurements later none was wasted and two *shrank* the work. → *Declared risk
is usually unmeasured work. Measure it; it often gets smaller.*

## Honest values

**An unknown became a plausible default.** `pricing_calc_profiles` had zero rows, and a 4%
SIMPLES rate entered the calculation as if measured. A zero-row configuration table is an
"unconfigured" state, not a default. → *Empty table ≠ default value.*

**A hardcoded tariff.** A literal `R$79` freight constant shipped and was later banned by
name. → *Business constants come from config, DB, or provider.*

**A company code was assumed.** `CODEMP=1` was hardcoded where the ERP has more than one
company. → *Every implicit scope is a hardcode with a longer fuse.*

## Failure visibility

**Token refresh failure was invisible at every layer, and the screen stayed green.** The
operator's account could be disconnected while the cockpit reported health. → *For anything
that can fail unattended, trace failure → pixel in the plan.*

**A silent page-1 truncation survived a live drive.** The dev account had less than one page
of data, so the bug could not appear. → *Prove pagination with a fixture larger than one page,
not with the live account you happen to have.*

**A stale binary made the live drive lie.** The running container predated the slice's commit,
so absence-of-observable read as absence-of-defect. → *Compare container build time to the
slice's commit before believing a live drive.*

## Test design

**Presence assertion over a wrong value.** `expect(screen.getByLabelText("Data freshness"))
.toBeInTheDocument()` passed with any string, which is exactly why a formatter that could not
distinguish 15 minutes from 15 days survived the full suite. → *Assert the value.*

**A criterion unfalsifiable against the type.** "Never calls X" cannot fail when the type does
not declare X. → *Read the acceptance criterion against the signature.*

**A synchronous non-call assertion over a post-`await` target passed even with an injected
click.** The assertion ran before the awaited code could call anything. → *The control must be
able to fail; prove it by injecting the defect.*

**A symmetric fixture asserted nothing.** Swapping the inputs would still pass, so it did not
test order. → *If the fixture is symmetric in the property under test, it is dead.*

**Green integration was not evidence until red named the test.** `status=passed` does not prove
execution; a `failure_token=test=` does. Skipped and green look identical. → *Assert on an
attributable token.*

**A deleted foreign test.** An assertion in someone else's file was removed rather than
restored or restated. → *Restore, or restate the guarantee where it now lives.*

**A partial guard under a total claim.** A docstring promised behaviour for every input while
the guard covered some. That is worse than no guard, because readers stop checking. → *A total
guarantee in prose is a claim about every input.*

## Boundaries and contracts

**A decorator silently dropped an optional port.** The catalog surface 503'd because an
optional dependency was lost in wiring, with nothing at compile time to catch it. → *Assert
optional-port wiring at compile time.*

**Scope loaded into context did not enter the downstream cache key.** Results from one scope
were served for another. → *If it changes the answer, it belongs in the key.*

**A `::text` cast without an alias hijacked a bare `ORDER BY`** in Postgres, silently changing
the ordering. → *Alias casted output columns.*

**A chip's merge reverted a feature from `main`** because the diff was computed against the
dispatch base rather than the target's tip. → *Merge gates diff against the tip of the target.*

## Redundancy

**A relative-time ladder copy-pasted into three files**, with drifting thresholds, while the
fix needed exactly one shared function. → *Grep the concept before writing the fourth copy.*

**Two `FreshnessIndicator` components** with different `aria-label`s, so every label-based test
bound to whichever copy its file imported. → *Grep component names across `ui` and `web-query`.*

**Eleven contract operations with no consumer**, one of them the entire feature a later plan
proposed to build from scratch. → *Search the spec for zero-call-site operations first.*

**Thirty-two call sites collapsed to one table; one hundred and twenty fields collapsed to
three actual sources.** The apparent surface was mostly restatement. → *Count the distinct
sources, not the fields.*

## Process

**A read-only gate certified criteria the chip had certified itself.** A reviewer without the
ability to run anything cannot discharge an execution criterion. → *Match the gate's tools to
what it must verify.*

**A repeated defect got patched three times.** The rule now: a third defect of the same shape
stops the patching and forces a class fix or a registered debt. → *Fix the class or name the
debt; do not patch a third time.*

**A verdict named no tree.** An automated gate reported on "the repo" while resolving paths
against a stale worktree; both halves of its accusation were false. The cost was not the wasted
check — it was that an alarm wrong twice trains its reader to skip the third. → *A verdict
prints the absolute path and tip SHA it measured.*
