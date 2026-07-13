# M-14-real-vertical-mvp-validation

```yaml
id: M-14
type: milestone
status: planned
owner: Mission Strategist
parent: MIS-001
created: 2026-07-13
updated: 2026-07-13
validation_level: QA-4
lifecycle_scope: milestone
```

## Mission

MIS-001 Mercado Livre Operating Cockpit MVP replan.

## Outcome

At one frozen SHA, a bounded values-minimized sample proves real Mercado Livre reads,
real Sankhya Oracle reads, PostgreSQL persistence/idempotency, and the browser journey
Overview → Listing → Product → Listing → Sale/Margin → Listing → Stock Simulation,
with no provider write.

## Why This Milestone Exists

Deterministic module tests prove rules but not that the selected integrations and UX
form a usable product. Previous QA also failed because registered Go commands were not
PowerShell-safe or rooted at the Go module.

## Features

| ID | Name | Brief |
| --- | --- | --- |
| F-01 | Bounded real scenario contract | Define a local values-minimized selector, reset boundary, concise provenance note, and no-write guard |
| F-02 | PowerShell-safe real-read lanes | Reuse or minimally close the PostgreSQL/ML/Oracle validation lanes and record the actual Windows commands |
| F-03 | Vertical browser evidence readiness | Make the accepted real sample addressable by IC-003 and the exact IC-004 browser drive |

## Dependencies

- M-09 and M-13 passed at accepted SHAs.
- Real Mercado Livre and Oracle read credentials available locally without entering artifacts.
- IC-003, IC-004, and M-14 validation contract.

## Risks

- Real data may lack a single sample joining product/listing/order/margin.
- Values or screenshots may expose buyer PII or credentials.
- Refresh/import could accidentally call a provider write path.
- Host-specific command syntax could invalidate evidence again.

## Done Means

- Registered commands are preflighted in Windows PowerShell with exact cwd and absolute GOCACHE.
- Concise evidence identifies source/type/time/command/read-only status without secret/PII values.
- Repeating bounded listing/order imports creates no duplicate durable identities.
- Browser completes the IC-003 journey and retains state after reload.
- Unknown/stale/conflict states are evidenced alongside a current path.
- Stock preview reports `executed=false`; network evidence contains no provider mutation.
- QA alone issues the M-14 verdict and mission rollup may then evaluate MVP completion.

## Handoff

- Current status: Planned.
- Next owner: Milestone Orchestrator after M-13 passes.
- Next action: Establish the no-write preflight and exact command registry before live reads.
- Required files/evidence: feature validations, sanitized real-read/PostgreSQL/browser evidence, fixed-SHA review, M-14 result.
- Blockers or open decisions: A missing real sample is reported as terminal status
  `externally_blocked` under the IC-004/harness checkpoint schema; it is not an API
  error code and must not be invented as application behavior.

## Correction Handoff

- QA failure summary: Not applicable during planning.
- Correction scope: None until QA reports a named failed criterion.
- Attempts used/remaining: 0/2.
- Next artifact: `F-01-bounded-real-scenario/feature.md`.
- Revalidation evidence required: all M-14 criteria at one frozen SHA.
