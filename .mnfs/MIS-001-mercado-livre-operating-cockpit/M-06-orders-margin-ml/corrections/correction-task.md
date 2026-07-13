# M-06 Correction Ledger

```yaml
id: M-06-COR-01
type: correction-task
status: in_progress
owner: Milestone Orchestrator
parent: M-06
created: 2026-07-12
updated: 2026-07-12
validation_level: QA-0
lifecycle_scope: milestone
```

## Failed Criteria

- M-06-C02: quantity-extended revenue was composed with unextended per-unit
  CUSSEMICM at the round-2 reviewed SHA.
- M-06-C03: the manual-adjustment write accepts caller-supplied actor identity;
  no verified principal or scope authorization is established.
- Gate integrity: exact executable evidence bindings and append-only retry
  accounting were incomplete.

Authoritative result: `../validation-result.md` and `../milestone-review.md`
(round 2, Fail).

## Assigned Scope

- Attempt 1: F-06 only, correcting quantity/cost amount scope without changing
  fee or tax line-amount semantics.
- Attempt 2: F-15 only, propagating non-complete tax input quality into item
  and order snapshot incompleteness while retaining known tax amounts.
- Production authentication/manual adjustments are deliberately deferred by
  owner priority. This ledger does not waive M-06-C03, claim authentication,
  or authorize caller-supplied identity. C03 remains pending/failing for the
  next full gate unless a separately approved trusted-principal boundary is
  implemented. No Candidate A approval may be reissued.
- Next safe product work is truthful, order-specific Oracle tax provenance and
  fresh resolved-order margin evidence. Unknown tax remains unknown; no tax is
  estimated.

## Allowed Paths

- Attempt 1 paths are exactly those recorded in
  `../F-06-quantity-cost-semantics/validation.md`.
- Attempt 2 paths are exactly those recorded in
  `../F-15-partial-tax-snapshot-quality/validation.md`.
- Any later correction work requires new owner authority because the recorded
  two-attempt budget is exhausted.

## Retry Fields

- correction_attempts: 2
- max_correction_attempts: 2
- last_validation_result: Fail (round 3; attempt 2 awaits replacement re-gate)

## Correction Log (append-only)

| Round | Attempt | Scope dispatched | Defect locus | Correction result | New `ran` evidence path | Re-gate verdict |
| --- | --- | --- | --- | --- | --- | --- |
| Historical | — | Earlier unnumbered fixes and live-evidence reconciliation | See round-2 review trail | Historical context retained; owner says these do not consume attempts | `../orchestrator-reconciliation-2026-07-11.md` | Round 2 Fail is the baseline |
| 2 | 1 | F-06 quantity/cost amount-scope correction | `apps/server_core/internal/modules/profitability/application/service.go` (`mapCostInput` composition) | Accepted at `2284c1d3bfcfa359a66777baad6c339083973538` | `../F-06-quantity-cost-semantics/validation.md` | Round 3 Fail at `81b8a4b12c3fe32c011f3d362ede393dd7484381`; C02 integration passed, ★2/★3/★7 remain failing |
| 3 | 2 | F-15 partial-tax snapshot-quality correction | `apps/server_core/internal/modules/profitability/application/service.go` (`applyInput` tax completeness) | Known non-complete tax amounts retained; item/order missing-tax propagation implemented; replacement review pending | `../F-15-partial-tax-snapshot-quality/validation.md` | Pending replacement fixed-SHA review and proportional QA; C03 remains deferred/failing |

## Required Commands Or QA

- Full fixed-SHA independent review after remaining bounded product work.
- Proportional QA must refresh exact `ran` evidence and alone may change the
  milestone verdict.

## Handoff

- Current status: correction attempt 2 implemented; replacement fixed-SHA
  review and proportional QA are pending.
- Next owner: Milestone Orchestrator.
- Next action: freeze the returned F-15 commit, request independent review,
  then proportional QA. Any later correction requires new owner authority.
- Blocker retained: M-06-C03 trusted-principal boundary is deliberately
  deferred, not passed or silently weakened.
