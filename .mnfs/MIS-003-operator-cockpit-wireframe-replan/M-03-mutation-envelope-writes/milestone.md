# M-03-mutation-envelope-writes

```yaml
id: M-03
type: milestone
status: planned
owner: Mission Strategist
parent: MIS-003
created: 2026-07-14
updated: 2026-07-14
validation_level: QA-0
lifecycle_scope: milestone
```

## Mission

MIS-003 — `../mission.md`; contracts: IC-03 (`../research/mutation-envelope-interface-contract.md`), IC-05 (crosswalk + protocolo route), ADR-13/16; governance `contracts/governance/execution-lanes.json` 7 provider-write gates.

## Outcome

The mutation envelope exists end-to-end: `mutation_protocols`/`mutation_items` tables, lifecycle state machine, in-process table-backed poller (FOR UPDATE SKIP LOCKED, chunks of 20), write types price_update / stock_correct / link_apply / listing_pause / listing_resync / listing_edit through provider capability adapters honoring all 7 write gates, preview/approve/cancel/retry endpoints, and the UI surfaces: preview/confirm modal wired to Anúncios bulk actions plus `/protocolos/:protocolId` detail page. `listing_create` returns 422 `type_not_enabled` (contract-only). Observable: bulk price update on stub adapter produces protocolo `MP-nnnnnn` with per-item before/after and honest failure codes; StockActionService envelope precedent folded in.

## Why This Milestone Exists

Every write in the cockpit flows through this one seam. Centralizing lifecycle + gates + audit here means M-04/M-05 surfaces only call `invalidateAfterMutation` and link to protocolos — no per-surface write logic to drift.

## Features

| ID | Name | Brief |
| --- | --- | --- |
| F-01 | protocolo-core | Tables, lifecycle, poller, idempotency, audit trail |
| F-02 | write-types-adapters | Six write intents, PriceWriter, StockWriter live wiring, StockActionService fold |
| F-03 | selection-preview-api | Bulk selection resolution, preview snapshot, approve/cancel/retry endpoints, OpenAPI+SDK |
| F-04 | preview-confirm-ui | Preview/confirm modal in Anúncios, protocolo detail page, crosswalk wiring |

Order F-01 → F-02 → F-03 → F-04. F-02 and F-03 both touch application layer — sequential, one writer.

## Dependencies

M-01 (listings rows are selection targets); M-02 F-02 (`invalidateAfterMutation`, failureCopy) for F-04. Integration lane uses stub provider adapter; live ML write lane requires explicit operator authorization per mission Validation Strategy.

## Ownership & Concurrency

Wave W1 (mission Parallel Execution Plan) — runs concurrent with CHIP-M02 and CHIP-SAT.

- Migration block: **0038–0042** reserved. Do not exceed; more needed → `REQUEST` to hub.
- OpenAPI/SDK: mutation/protocolo paths + schemas only (additive; CHIP-SAT owns
  dashboard/orders/sync-runs and market sections — never touch those).
- Additive contract-locks held: composition root (mutations module registration lines only);
  `connectors` PriceWriter/StockWriter wiring (F-02 only). Released at CLOSED, diffs called
  out in the event.
- F-04 FE gate: preview/confirm modal mounts in M-02 F-03's Anúncios workspace. F-01..F-03
  proceed regardless; start F-04 FE only after the hub confirms M-02 F-03 merged and triggers
  the rebase. If M-02 stalls, `BLOCKED` — do not build a stand-in surface.
- Governance base anchor: pinned in chip prompt at dispatch (profile §2).

## Risks

- RK-03 (provider write hazard): stub adapter default; live lane gated on operator; idempotency keys prevent duplicates on retry/crash.
- Poller starvation/crash mid-chunk (milestone-local risk, no mission RK ID): SKIP LOCKED + item-level status means crash resumes safely; validated by kill-mid-apply test.
- StockActionService fold regression: existing `inventory_stock_actions` flows re-validated after fold.

## Done Means

All IC-03 operations + lifecycle proven per `validation-contract.md` (M03-C01..C12); 7 write gates each proven with negative test; retry clones to new protocol (audit immutability); governance lanes green.

## Handoff

- Current status: planned.
- Next owner: Milestone Orchestrator.
- Next action: dispatch F-01 with IC-03 + migration precedent 0026.
- Required files/evidence: `F-*/validation.md`, `validation-result.md` here.
- Blockers or open decisions: live-write authorization decision deferred to milestone QA time (operator).

## Correction Handoff

Not applicable at planning time.
