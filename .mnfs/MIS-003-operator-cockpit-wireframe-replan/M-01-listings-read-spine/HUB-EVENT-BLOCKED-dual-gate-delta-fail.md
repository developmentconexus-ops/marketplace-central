# HUB EVENT — BLOCKED — dual-gate DELTA FAIL (slice 8 degrade returns confidently-wrong rows)

```yaml
event: BLOCKED
from_chip: M-01-listings-read-spine
to: hub (local_efa46c30-1c0c-4075-9671-c2d7ae9efabe)
date: 2026-07-16
branch: mis-003/m-01-listings-read-spine
tip: c4e8ab913c132d4929c3bc60156a908379a78043 (not pushed)
gate: P6 dual-gate DELTA (slices 8+9, base e2cde36)
verdict: FAIL — BOTH sides, independently, in agreement
blocks: P8 CLOSED
evidence: _gate-evidence/round-2/dual-gate-delta-verdict.md (full reconciliation table)
```

## Outcome

Cold Opus subagent (`model=opus`, clean context) and GPT-5.6 Sol (`--effort medium`, fresh thread) ran
SIMULTANEOUS at fixed SHA on prompt-pack §14, neither seeing the other. **Both returned FAIL.** No
contradictions: Sol raised a strict subset of Opus's blocking set at equal-or-higher severity; Opus added
three verified findings Sol did not reach. Merge (§8) = union, severity = max.

**C10 PASS stands and both reviewers independently certified it honest.** **Slice 9 is CLEAN** — no rework.
The gate fails entirely on **slice 8's degrade path**.

## The defect (agreed blocking, B1)

On an Oracle outage the degraded `passThrough` re-queries with the original `q` and never applies
`matchesDependentFilter`. `needsBelowMarginScan` is true for **`HasException != nil`** as well as
`exception=below_margin`, and the Postgres repository has no SQL case for either. So
`GET /listings?filter.has_exception=false` during an outage returns rows whose `sync_state='error'` /
`sync_error` / unresolved link make them **provable non-matches** — facts that live in Postgres, are
fully available, and do not depend on Oracle at all — as a confident 200 with **every field non-null**.

This is not "unknown rendered as null" (ADR-17 compliant). It is a **known-false fact returned as a
match** — what ADR-17 exists to prevent. The slice-8 plan's justification ("the nullable row fields
visibly communicate that the exception could not be evaluated") is factually inapplicable: nothing in
that response is null. Pre-slice-8 this outage 503'd and served no wrong data. **Slice 8 introduced it:
it traded an honest 503 for a confidently wrong 200.** Both filters are contract-declared and
operator-reachable (`openapi.yaml:225,323`; `transport/query.go:86-93,114`). `passThroughGroups` shares
the defect.

Secondary agreed blocking (B2): **every** cost-reader error — cancellation, timeout, adapter/data defect
— collapses into `factsUnavailable` null, indistinguishable from genuinely-absent data. Only
source-unavailable should degrade; the rest must propagate. And it is swallowed with **zero telemetry**
(`read_service.go` has no `slog` call at all): a total Oracle outage = HTTP 200, null margins
cockpit-wide, not one log line — while the same delta logs a WARN for a single unmapped status string.

## Why I'm escalating rather than driving slice 10 unilaterally

The direction of the fix is settled (degrade is right; its filter semantics are wrong). **What a
fact-dependent filter should DO when its fact cannot be evaluated is a contract/product decision**, and
one option changes the API surface. I'm not choosing that alone. Options, with the reviewers' own
framing:

- **(A) Partial-honest, no contract change.** In the degraded path apply the SQL-computable component
  (`sqlException`) and treat only the margin component as unknown. Kills the `has_exception` false
  positives outright. Residual ambiguity: is a row with no SQL exception but an unevaluable margin a
  match for `has_exception=false`? Needs a rule.
- **(B) Explicit degraded signal.** Serve the unfiltered degraded read but carry a contracted
  `service_degraded` / "filter not applied" signal in the response. Most honest; **requires OpenAPI +
  SDK change** (new seam → hub grant). The slice-8 brief itself listed this as acceptable.
- **(C) Fail the fact-dependent filter honestly.** Only requests that USE a fact-dependent filter error
  during an outage; unfiltered/independent reads still degrade to null gracefully. Narrowest blast
  radius, no contract change, but restores a 503-class response for those specific queries.

**My recommendation: (A) + (C) hybrid** — apply the SQL-computable component always (that data is
available; there is no excuse for getting it wrong), and for the margin-dependent component that
genuinely cannot be evaluated, fail *that* filter honestly rather than guess. (B) is strictly better on
honesty and I'd take it if you'll grant the OpenAPI/SDK seam, but it widens M-01's contract at close time.

## Requests to hub

- **R1 (ruling):** pick A / B / C / A+C for degraded fact-dependent filter semantics. If B, grant the
  OpenAPI + SDK seam.
- **R2 (scope):** confirm slice 10 = slice-8 degrade path only (read_service `passThrough` +
  `passThroughGroups` + error classification + telemetry + the dead `unavailablePolicyService` wrapper +
  re-pin the tests that currently encode the defect). Slice 9 untouched.

## Held state (nothing regressed, nothing pushed)

- c4e8ab91 committed on the worktree branch only. **No push.**
- C10 PASS (unknown 0.0%) stands — independent of this defect, both reviewers certified it.
- Slice 9 clean and certified by both sides.
- **P8 CLOSED withheld.** Only QA passes a milestone, and the gate did not pass.
