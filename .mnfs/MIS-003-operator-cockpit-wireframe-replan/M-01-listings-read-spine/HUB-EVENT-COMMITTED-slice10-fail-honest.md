# HUB EVENT — COMMITTED (slice 10, fail-honest dependent filters)

```yaml
event: COMMITTED
from: chip M-01-listings-read-spine
to: HUB (local_efa46c30-1c0c-4075-9671-c2d7ae9efabe)
sha: 7f5a1b8c
base: 1f6b72d8 (= c4e8ab91 slice 9 + round-2 gate evidence)
scope: exactly 3 source files, hub R2 scope — no migration, no OpenAPI, no SDK, no transport
pushed: NO (never without explicit operator permission)
```

## What landed

Hub ruling **(C) PURE** implemented. Any fact-dependent filter (`has_exception` both
directions AND `exception=below_margin`, uniform via `needsBelowMarginScan`) during a genuine
fact-source outage now returns source-unavailable through the EXISTING error envelope →
transport already 503s `source_unavailable`, so **no contract change**. Reads without such a
filter keep degrading to null (ADR-17), pinned by an over-correction guard test.

Gate covers **both** Oracle facts: ceiling (known pre-scan) and cost (only discovered mid-scan,
threaded out of `enrich` via `(error, error)`). A ceiling-only pre-check would have left the
cost hole open. Strict classification via existing `IsReadErrorCode/ReadErrorSourceUnavailable`
— cancellation/timeout/adapter defects propagate + log ERROR. Degrade logs `slog.Warn` once per
request with `op` + `fact`, never per row. `passThrough`/`passThroughGroups` and the dead
`unavailableListingPolicyReader` wrapper deleted.

Option (A) absent, as ruled. Option (B) `service_degraded` remains deferred with a name.

## Gates (all green, independently run by chip, not relayed)

| Gate | Result |
|---|---|
| L0 build + vet | exit 0 / exit 0 |
| L1 unit | exit 0 |
| §14 cold sonnet review | **APPROVE** — zero blocking, zero important |
| Integration lane (ephemeral pg) | **GREEN** — migrate 37→0, 0037 constraint 8 values intact, `TestListingsReadContractEndToEnd` 8/8 incl `null_cost_honesty` + `tenant_isolation`, `Performance2000` p95 3.93ms |

Evidence: `_gate-evidence/round-3/` (`slice10-L0-report.md`, `slice10-review.md`,
`slice10-candidate.diff`, `delta-round3.diff`).

## Deviation logged

Implementation ran on **sonnet** per explicit operator directive (finish M-01 fast).
Overrides HARNESS §1 (implement = GPT-5.6 Luna high). Operator authority wins; recorded here.
Test-first + §14 cold independent review preserved (implementer ≠ reviewer).

## Ask

Restart the dev stack at `7f5a1b8c` (pre-armed, no formal REQUEST per your last message).
Chip re-drives to confirm no regression, then runs gate round 3 (§9 delta since `c4e8ab91`
+ explicit resolution check of every round-2 finding) → P8 CLOSED.
