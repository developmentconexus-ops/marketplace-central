# Anti-patterns — defects this repo actually shipped

Each entry is a real case. Read this last, as an adversarial pass over a plan that feels
finished: for every entry, ask *"is my plan doing this?"*

The pattern behind most of them: **an observable that is stable, cheap, and wrong about which
world you are in.** It reproduces, so it feels like evidence; it passes in both the fixed and
the broken world, so it proves nothing.

> **Source of truth: `docs/engineering/defect-class-catalog.md`.** This file is the
> *planning-time projection* of it — one line per class, kept short enough to read at the end
> of a plan. The catalog carries the evidence, the prevention rung and the CI gate; this file
> carries only what changes a plan. The `§` tag on each entry is its class there; `§MD` means
> the class lives in the sibling catalog (`MetalDocs/docs/engineering/defect-class-catalog.md`,
> the code-and-architecture axis); `§—` marks a case not yet classified in either, which is a
> TODO, not a category.
>
> **Do not add an untagged entry here.** Two lists that must agree, maintained by hand, is the
> catalog's own §D1. New defect → write the class first, then project it here if it changes
> planning. Classes with no planning consequence (the catalog's Part V, factory machinery)
> deliberately do not appear.

---

## Planning

**The plan built a second engine beside the working one.** §P5 — A 17-task fiscal plan
constructed a calculation layer next to `pricing/domain`, which already produced the right
number. The column on screen was empty, and empty column read as "feature absent" instead of
"feature not wired". Re-scoped to 7 tasks. → *An empty screen is not evidence of missing code.
Trace the value from the renderer backwards before planning to produce it.*

**The plan named a column that does not exist.** §P1 — A query was written against
`sync_state.last_success_at`; the table has `last_full_sync_at` and `last_incremental_at`.
Caught by running the query, not by reading it. → *Measure schema with `\d`, paste the output
into the plan.*

**The plan hand-rolled a scheduler.** §P5 — A market-collection ticker was specified beside
`syncapp.Scheduler`, discarding `sync_state` reconciliation, per-entity failure isolation,
cursor persistence and — the decisive one — `RecordFailure` feeding the operator's sync-health
card. A broken collection would have shown a green screen. → *Check whether the seam already
exists; check what reuse buys you that a parallel build does not.*

**The brief was an allegation and it rotted.** §P1 §P2 — Three claims in one milestone brief
were false by the time a worker read them. Line anchors measured in one tree pointed at
different code in another. → *Anchor by content — route, `operationId`, schema name, function
name. When a line is unavoidable, name the tree.*

**A declared risk was treated as acceptable debt.** §P4 — Three declared risks were rejected by
the operator; seven measurements later none was wasted and two *shrank* the work. → *Declared
risk is usually unmeasured work. Measure it; it often gets smaller.*

**Every task was green and the feature did not exist.** §P3 — P2.b shipped 7 tasks, all lanes
passing, and `icms_matrix_mirror` had **zero rows**: a reader and a writer both landed with the
only caller in a `_test.go`. ADR-17 working correctly *hid* it, because honest-unknown and
never-populated render the same. → *Name the composition site — the production call path ending
in a row or a pixel. Acceptance is a `count(*)`, never a green test.*

## Measurement

**The instrument answered a different question.** §V3 — `n_live_tup` reported `listings = 0`
against a real `count(*)` of 34 (autovacuum estimate). OpenAPI `operationId` was compared to
frontend call sites to find orphan operations, but the SDK is hand-written and renames. The
same sweep over `apps/web/src` alone missed `packages/feature-*` and `web-query`. A `TGFDIN`
join without `CODINC` fanned out and turned 13 exceptions into "342". → *State the universe with
every count. A number that carries a verdict needs a second, independent instrument.*

**A zero was read as proof.** §V7 — A repeated sweep was offered as evidence of ADR-C2
compliance; with nothing to delete it would return 0 in both worlds. Only a positive control
plus `fetched_at` advancing proved anything. → *Every zero-result check ships a positive control
in the same run.*

**The verdict came from one document.** §H4 — A fiscal rule was refuted against a single
invoice carrying a stale cadastro; its sibling from the previous day proved the rule. The ERP is
not self-consistent, so "the ERP disagrees" is not a refutation. → *A refutation needs n > 1 and
a stated distribution. One record is a hypothesis.*

**Two systems disagreed and the plan elected one winner.** §H3 — The ERP's *method* was right
and ours was wrong; our *data* was right and the ERP's was stale. A single-slot verdict was
wrong twice, in opposite directions. → *Decompose before adjudicating: formula, input, scope,
vintage. One verdict per axis.*

## Honest values

**An unknown became a plausible default.** §H1 — `pricing_calc_profiles` had zero rows, and a 4%
SIMPLES rate entered the calculation as if measured. A zero-row configuration table is an
"unconfigured" state, not a default. → *Empty table ≠ default value. And the mirror case is also
a defect: a **known** zero written as unknown destroys information the other way.*

**A hardcoded tariff.** §H2 — A literal `R$79` freight constant shipped and was later banned by
name. A 27-row tax table lived as a Go literal in `difal_seed.go` with zero production
consumers. → *Business constants come from config, DB, or provider.*

**A company code was assumed.** §H2 — `CODEMP=1` was hardcoded where the ERP has more than one
company. → *Every implicit scope is a hardcode with a longer fuse.*

**A sentinel was read as a wildcard.** §H5 — `TGFICM`'s `S` means "this axis is switched off",
never "matches anything", and its two restriction pairs are different axes rather than
interchangeable slots. Reading it as a wildcard admitted rules belonging to other grupos;
correcting the predicate took 881 ambiguous cells to 2. → *Enumerate the value distribution
before writing the predicate. A predicate written from a schema description is a guess.*

## Failure visibility

**Token refresh failure was invisible at every layer, and the screen stayed green.** §V8 — The
operator's account could be disconnected while the cockpit reported health. → *For anything that
can fail unattended, trace failure → pixel in the plan.*

**A silent page-1 truncation survived a live drive.** §V2 — The dev account had less than one
page of data, so the bug could not appear. → *Prove pagination with a fixture larger than one
page, not with the live account you happen to have.*

**A stale binary made the live drive lie.** §V4 — The running container predated the slice's
commit, so absence-of-observable read as absence-of-defect. → *Compare container build time to
the slice's commit before believing a live drive.*

## Test design

**Presence assertion over a wrong value.** §V2 — `expect(screen.getByLabelText("Data freshness"))
.toBeInTheDocument()` passed with any string, which is exactly why a formatter that could not
distinguish 15 minutes from 15 days survived the full suite. → *Assert the value.*

**A criterion unfalsifiable against the type.** §V2 — "Never calls X" cannot fail when the type
does not declare X. → *Read the acceptance criterion against the signature.*

**A synchronous non-call assertion over a post-`await` target passed even with an injected
click.** §V2 — The assertion ran before the awaited code could call anything. → *The control must
be able to fail; prove it by injecting the defect.*

**A symmetric fixture asserted nothing.** §V2 — Swapping the inputs would still pass, so it did
not test order. → *If the fixture is symmetric in the property under test, it is dead.*

**Green integration was not evidence until red named the test.** §V1 — `status=passed` does not
prove execution; a `failure_token=test=` does. Skipped and green look identical. → *Assert on an
attributable token.*

**A deleted foreign test.** §— An assertion in someone else's file was removed rather than
restored or restated. → *Restore, or restate the guarantee where it now lives.*

**A partial guard under a total claim.** §D3 — A docstring promised behaviour for every input
while the guard covered some. That is worse than no guard, because readers stop checking. → *A
total guarantee in prose is a claim about every input.*

## Boundaries and contracts

**A decorator silently dropped an optional port.** §MD — The catalog surface 503'd because an
optional dependency was lost in wiring, with nothing at compile time to catch it. → *Assert
optional-port wiring at compile time.*

**Scope loaded into context did not enter the downstream cache key.** §MD — Results from one
scope were served for another. → *If it changes the answer, it belongs in the key.*

**A `::text` cast without an alias hijacked a bare `ORDER BY`** in Postgres, silently changing
the ordering. §MD → *Alias casted output columns.*

**A chip's merge reverted a feature from `main`** because the diff was computed against the
dispatch base rather than the target's tip. §V6 → *Merge gates diff against the tip of the
target.*

## Redundancy

**A relative-time ladder copy-pasted into three files**, with drifting thresholds, while the fix
needed exactly one shared function. §D1 → *Grep the concept before writing the fourth copy.*

**Two `FreshnessIndicator` components** with different `aria-label`s, so every label-based test
bound to whichever copy its file imported. §D1 → *Grep component names across `ui` and
`web-query`.*

**Three copies of the same 27-state tax table** — two in the database, one a Go literal — whose
values disagreed in exactly three states, by exactly the FCP. §D1 → *Copies that disagree are
more informative than copies that agree: their delta says what the two are really measuring.*

**Eleven contract operations with no consumer**, one of them the entire feature a later plan
proposed to build from scratch; and four override mechanisms in one module, none in use. §D2 →
*Search the spec for zero-call-site operations first.*

**Thirty-two call sites collapsed to one table; one hundred and twenty fields collapsed to three
actual sources.** §D1 — The apparent surface was mostly restatement. → *Count the distinct
sources, not the fields.*

## Process

**A read-only gate certified criteria the chip had certified itself.** §V5 — A reviewer without
the ability to run anything cannot discharge an execution criterion. → *Match the gate's tools to
what it must verify.*

**A repeated defect got patched three times.** §V9 — The rule now: a third defect of the same
shape stops the patching and forces a class fix or a registered debt. → *Fix the class or name
the debt; do not patch a third time.*

**A verdict named no tree.** §V4 — An automated gate reported on "the repo" while resolving paths
against a stale worktree; both halves of its accusation were false. The cost was not the wasted
check — it was that an alarm wrong twice trains its reader to skip the third. → *A verdict prints
the absolute path and tip SHA it measured.*
