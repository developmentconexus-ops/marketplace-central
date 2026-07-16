# HUB EVENT — COMMITTED — slice 9 (grow canonical listing-status enum)

```yaml
event: COMMITTED
from_chip: M-01-listings-read-spine
to: hub (local_efa46c30-1c0c-4075-9671-c2d7ae9efabe)
date: 2026-07-16
branch: mis-003/m-01-listings-read-spine
sha: c4e8ab913c132d4929c3bc60156a908379a78043
prev_tip: a595f36c (slice 8 code) / d50636e7 (slice 8 evidence)
scope: C10 corrective — grow status enum (hub R2/R3 ruling)
pushed: NO (no operator permission)
```

## What landed

Grew canonical listing-status enum: added `under_review`, `inactive`, `payment_required`,
`not_yet_active`; kept `unknown` as ADR-17 honest fallback. Mapper maps all 7 documented ML
statuses explicitly + slog.Warn(raw provider status) on default→unknown ONLY. One commit across
domain+IsValid, mapper, NEW migration 0037 (additive CHECK widen), migration-count guards 36→37,
transport parse test, OpenAPI (3 enum sites) + SDK union (8-value parity). F-01 not reopened.

## Gates (all green)

- L0: build 0, vet 0 (whole repo).
- L1 unit: listings + composition + platform/migrate + migrations all ok; SDK tsc 0, vitest 43/43.
- L1 integration lane over 0037: migrate 37→0 idempotent; live `listings_status_check` verified
  = 8-value set; `TestListingsReadContractEndToEnd` (8 subtests) + `TestListingsReadPerformance2000`
  PASS (exit 0). Evidence: `_gate-evidence/round-2/slice9-L0-report.md`.
- §14 independent cold sonnet review: no blocking; sole `important` (lane-evidence gap) CLOSED by
  the green lane run. Effective APPROVE. Evidence: `_gate-evidence/round-2/slice9-review.md`.
- Test-first RED→GREEN proof in `F-02-listings-read-api/validation.md`.

## Next (restart pre-armed per hub AJUSTE DE FLUXO — no formal REQUEST)

Requesting hub re-sync the installation on the backend at c4e8ab91, then I re-drive C10:
expect unknown ≪ 20% (likely 0%); I record the per-status counts = R1-closing evidence.
Then single dual-gate DELTA from `e2cde36` covering slices 8+9 (cold Opus subagent + Sol medium,
simultaneous, §8 merge + reconciliation table) → P8 CLOSED.
