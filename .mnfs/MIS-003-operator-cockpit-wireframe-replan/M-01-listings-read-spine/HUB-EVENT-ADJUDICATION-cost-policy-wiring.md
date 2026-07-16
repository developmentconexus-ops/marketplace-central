# HUB EVENT — ADJUDICATION — cost/policy reader wiring vs M-01 scope

```yaml
event: ADJUDICATION (response to hub composition finding on C10)
from_chip: M-01-listings-read-spine
to: hub (local_efa46c30-1c0c-4075-9671-c2d7ae9efabe)
date: 2026-07-16
trigger: hub found root.go:494 nil-DB BatchReader + :495 permanent policy stub → /listings still 503 even with Oracle creds set
ruling: WIRING OUT-OF-SCOPE (path 2). Fix (a) degrade-only stays mandatory. No slice change.
```

## Finding under adjudication (hub)

- `root.go:494` `listingCostReader := NewCostReader(internalreadoracle.NewBatchReader(nil, sem))` — **nil DB hardcoded** → listings cost reader always errors regardless of env.
- `root.go:495` `listingPolicyReader` = `unavailableListingPolicyReader{...unavailablePolicyService{err: ReadErrorSourceUnavailable}}` — **permanent stub**.
- Consequence: even with Oracle creds now SET + pinging OK in the container, `/listings` still 503.

## Adjudication against contract (repo truth > inference)

**Ruling: live wiring of listings cost/policy readers is NOT M-01 scope — deferred by design. Path 2.**

Evidence (validation-contract.md, the binding M-01 contract):

| Source | Says | Implication |
|--------|------|-------------|
| Required Outcome (:25) | listings table populated by refresh + 5 endpoints canonical shapes + OpenAPI/SDK same commit | no live cost/policy wiring named |
| **M01-C10 (:172, :175)** | "no `unknown` flood: **<20% unknown status**"; blocking = ">20% unmapped **statuses** (adapter mapping gap)" | C10's unknown metric = listing **status** (Postgres-backed, from ML adapter at refresh), **NOT** margin |
| **M01-C07 (:119-131)** | margin criterion = null-**honesty** (`below_margin_worst_case: null`, additive `margin_unknown`), tested on **seeded NULL cost** | M-01 owns the margin FIELD shape + honesty; does **not** require live Oracle values |

`root.go:494-495` are **intentional deferral placeholders**: the policy reader is a full `unavailablePolicyService` stub (an unimplemented adapter, not a missing DB handle), and the cost reader is deliberately nil-DB. Live cost/policy sourcing belongs to a downstream pricing/margin milestone, not the read spine. Naming confirms intent: M-01 is the **read spine** (serves listings + status), margin enrichment is a downstream concern that must **degrade** when its source is absent.

## Consequences

1. **Fix (a) degrade — CONTINUES MANDATORY.** Independent of wiring: Oracle can drop in prod, and per ADR-17 an optional margin fact must never block the read. The in-flight slice (Codex Luna high, test-first) implements fix (a) only and needs **no change** from this adjudication. Scope guard holds: `read_service.go` + tests ONLY; **no `root.go` edit**, no adapter wiring.
2. **C10 secondary is satisfiable once fix (a) lands** — no wiring needed. list/by-product degrade → serve the 34 tenant-scoped rows → observe real `MLB…` ids + compute **status**-unknown% (status is Postgres, unaffected by nil cost reader).
3. **Margin will read 100% `margin_unknown` (null)** with the nil cost/policy readers. This is **C07-consistent** (honest null, never false/zero) and is **NOT** C10's "unknown status" metric → **no flood violation**, no conflict. C10 blocking failure is `>20% unmapped statuses`, not margin.
4. **No complementary SELECT required for the metric** — list will expose ids + status once fix (a) lands. Hub's offered read-only `SELECT provider_listing_id, status …` is welcome as *corroborating* evidence but not load-bearing for the pass.

## Recommended follow-up (not M-01)

Flag a tracked downstream item (pricing/margin milestone or M-02+ chip): **wire listings cost + policy readers to live Oracle** (`root.go:494-495` → real DB handle + implemented policy adapter). Keep it out of M-01 so the read-spine close is not gated on margin sourcing. Fix (a) makes the read spine correct regardless.

## FINAL resolution (2026-07-16) — hub ACK: WIRE (deferral below is VOID)

Hub verified the wiring diagnosis in-worktree (oracleDB :352/:474, marketSvc :491 Postgres-real) and
APPROVED full wiring. **The ACCEPT-WITH-CONDITION deferral path below is SUPERSEDED/VOID** — hub
deleted the board deferral-tracker task ("sem stub alcançável = melhor que deferral"). NO-STUB is
satisfied by REMOVING the stubs in this slice, not deferring. Corrective slice = fix (a) degrade +
`root.go` wiring (cost `nil`→`oracleDB`, policy stub→`marketSvc`, drop dead `unavailablePolicyService`;
keep degrade wrapper). Zero-value `oracleDB` degrades cleanly via `ensureBatchAvailable`
(`batch_reader.go:204-212` `r.db==nil`→SourceUnavailable, no panic) — reviewer confirms. See
`F-02.../slice8-corrective-brief.md` FINAL ADDENDUM for full loci.

---
## (SUPERSEDED) Hub resolution (2026-07-16) — ACCEPT-WITH-CONDITION

Hub accepted the scope analysis (M-01 does not require live cost/policy; wiring would gate the
read-spine close on margin sourcing). A transient "wire now" ratification note crossed in flight and
is **superseded** by this ACCEPT-WITH-CONDITION. Condition (no-stub doctrine): the placeholder is
legal only with a **dated deferral**. Fix (a) stays degrade-only; `root.go` is NOT touched in this
milestone.

## DEFERRAL (formal, no-stub-legal) — REQUIRED by dual gate

- **Date:** 2026-07-16
- **Stubs:** `root.go:494` `NewBatchReader(nil, …)` (listings cost reader, nil DB) + `root.go:495`
  `unavailablePolicyService` (permanent policy stub).
- **Behavior covered:** reads degrade to honest null (ADR-17) via fix (a) — never zero/default.
- **Target:** pricing/margin milestone **M-02+** (hub board task created) — wire cost→live Oracle
  (`nil`→`oracleDB`, real at root.go:352/:474) + real policy adapter
  (`unavailablePolicyService`→`marketSvc`, Postgres service satisfies `listingsmarketplaces.Service`
  field-for-field: service.go:109→repository.go:178). Both chip-verified trivial + env-real.
- **Authority:** hub adjudication 2026-07-16 under no-stub doctrine (docs/HARNESS.md §4/§5).
- **Dual-gate:** a reviewer flagging the `root.go:494-495` stub is pointed HERE as evidence — the
  deferral is dated, targeted, and hub-authorized, so it is not an unregistered violation.

## Next (fix a only)

worker slice (fix a, `read_service.go` + tests ONLY, in flight) → independent review → L0/L1 →
commit → REQUEST restart-backend → re-drive C10 (ids + status-unknown%) → dual-gate DELTA from
e2cde36 (deferral above is gate evidence) → P8 CLOSED. No `root.go` touch. No push.
