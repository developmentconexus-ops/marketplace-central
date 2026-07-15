# F-01-protocolo-core

```yaml
id: F-01
type: feature-brief
status: briefed
owner: Mission Strategist
parent: M-03
created: 2026-07-14
updated: 2026-07-14
validation_level: QA-0
lifecycle_scope: feature
```

## Mission

MIS-003. Binding contract: IC-03 `../../research/mutation-envelope-interface-contract.md` (entities, lifecycle enum, DB shape, id/idempotency formats, caps). ADR-13 (table + in-process poller, NO outbox/bus).

## Milestone

M-03 mutation-envelope-writes.

## Brief

Create the `mutations` module core: `mutation_protocols` + `mutation_items` migrations per IC-03 DB shape, protocol id sequence rendering `MP-000042`, lifecycle state machine (draft→previewed→approved→applying→applied|partially_failed|failed_preserved|cancelled) enforced in domain layer, item idempotency key `{protocol_id}:{listing_id}`, and the in-process poller: claims approved protocols with `FOR UPDATE SKIP LOCKED`, applies items in chunks of 20 through a `WriterPort` (stub in this feature), records per-item before/after audit + failure `{code, message_pt, message_provider, retryable}` from the IC-03 taxonomy, computes terminal protocol state from item outcomes.

EARS:
- While a protocol is approved, when the poller claims it, no second poller instance shall claim it (SKIP LOCKED proof: two concurrent pollers, one claim).
- While items are applying, when the process crashes mid-chunk, a restarted poller shall resume unprocessed items only (idempotency: applied items not re-sent).
- While all items fail, when the run ends, protocol shall be `failed_preserved` with items + failures intact (M-06/MIS-001 preserved-failed semantics).
- While an invalid transition is requested (e.g. approve a draft never previewed), the domain shall reject with 409 lifecycle error, state unchanged.

## Inputs

- IC-03 (verbatim shapes/enums/caps), migration 0026 `inventory_stock_actions` (envelope precedent), R-03 poller notes, composition-root wiring, existing tx helper patterns in server_core.

## Expected Output

- Migrations (next numbers after M-01's) with enum checks + indexes for poller claim query.
- Domain: `MutationProtocol`, `MutationItem`, transition function (pure, table-driven, unit-tested per edge — every IC-03 arrow + every illegal pair).
- Poller runnable from composition root (interval config), chunk 20, per-item timing.
- Stub WriterPort with programmable outcomes for tests.
- Tests: transition table, SKIP LOCKED concurrency, crash-resume, terminal-state computation (all-ok / mixed / all-fail).

## Constraints

- No provider adapter code here (F-02); no HTTP endpoints (F-03).
- No outbox, no message bus, no background framework — plain goroutine + ticker per ADR-13.
- Audit rows immutable: no UPDATE on applied/failed items ever.
- Failure codes restricted to IC-03 enum; unknown provider errors map to `internal` retryable=false, `message_provider` preserved.

## Negative Scenarios

- Approve without preview → 409 lifecycle error.
- Cancel while applying → rejected (only draft/previewed cancellable per IC-03).
- Duplicate poller claim → exactly one winner.
- Item write attempted twice (restart) → second attempt skipped via idempotency key check.

## Validation Expectations

- `go test` output: transition-table test enumerating legal+illegal transitions; concurrency test transcript.
- SQL proof: crash-resume test showing pre-crash applied items untouched, remaining items completed.
- Terminal-state matrix test output: {all ok→applied, mixed→partially_failed, all fail→failed_preserved}.

## Execution Artifact Rules

`spec.md`, `plan.md`, `validation.md` created during feature execution.

## Handoff

- Current status: briefed.
- Next owner: Feature Implementer.
- Next action: compile context pack; read IC-03 + migration 0026 + composition root only.
- Required files/evidence: `validation.md` in this folder.
- Blockers or open decisions: none.
