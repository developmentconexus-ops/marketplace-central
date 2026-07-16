# C10 re-drive (post slice 9) — PASS + R1 closed

```yaml
lane: C10 live-provider-read
backend: c4e8ab913c132d4929c3bc60156a908379a78043 (hub restart, 0037 applied, listings_status_check = 8 values live)
installation: inst-mercado_livre-d373dc64-577f-4950-b11b-b90244f30cb2 (tenant_default)
base: http://127.0.0.1:8080
re_sync: POST /listings/refresh {installation_id} -> 202 accepted (async), driven by chip per hub ENV-READY
as_of: 2026-07-16T11:06:26Z
verdict: PASS
```

## The gate that was failing

`validation-contract.md:172,175` — expected `<20%` unknown status; blocking failure = ">20% unmapped
statuses (adapter mapping gap)". Pre-slice-9 re-drive: **unknown 7/34 = 20.6% → FAIL**.

## Post-slice-9 result

`GET /listings?installation_id=…&limit=100` → **200**, 34 rows, **34/34 real `MLB…` provider ids**, next_cursor null.

| status | count | share |
|---|---:|---:|
| paused | 17 | 50.0% |
| active | 10 | 29.4% |
| **under_review** | **7** | **20.6%** |
| unknown | **0** | **0.0%** |

**unknown = 0/34 = 0.0% < 20% → C10 PASS.**

`GET /listings/summary?installation_id=…` → **200**
`{"total":34,"active":10,"paused":17,"exceptions":{"sync_error":0,"stale":0,"unlinked":33,"below_margin_worst_case":0,"margin_unknown":1},"as_of":"2026-07-16T11:06:26Z"}`
(unchanged vs pre-slice-9 — the summary counts active/paused only; the 7 rows moved out of `unknown` into
`under_review`, not into active/paused. Confirms the grow did NOT distort existing buckets.)

## R1 CLOSED — the raw ML status behind the 7 unknown rows

The escalation asked the hub for the raw provider status distribution of the 7 unknown rows (chip was
data-blind: canonical-only reads, no DB/.env/self-boot). The mapped re-drive answers it empirically
**without ever needing the DB query**: all **7 = `under_review`**. Every other documented ML status
(`inactive`, `payment_required`, `not_yet_active`) has 0 rows in this tenant's inventory today — mapped
and ready, but currently unused. R1 requires no further hub action.

## Why R2 (grow, don't collapse) was the right ruling — proven, not asserted

Had `under_review` collapsed into `paused`, the cockpit would report **paused 24 (70.6%)** and an operator
would see 7 in-review items as reactivatable. They are not. The grow keeps the distinction honest and the
`unknown` bucket free to do its ADR-17 job (it is now genuinely empty = no unrecognized status in this
inventory, which is the correct signal, not a masked one).

## Remaining C10 facts (unchanged from the pre-slice-9 drive, still green)

- `GET /listings` 200 (fix (a) degrade holds — was 503 pre-slice-8).
- `GET /listings/by-product` 200, groups served.
- Oracle cost wiring live-proven: linked row `MLB4735328201` (pid 15956) → `cost=91.57` real Oracle.
- 33 unlinked → cost null (legit ADR-17); 1 linked → `margin_unknown` (cost present, ceiling/policy absent
  → honest null, C07-consistent).

## Gate status

C10 secondary (status mapping) was the sole blocker. **Now PASS.** Unblocks: dual-gate DELTA (slices 8+9,
base e2cde36) → P8 CLOSED.
