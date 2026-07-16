# C10 Live-Provider-Read — REAL drive vs hub stack (ENV-READY v2)

```yaml
criterion: M01-C10
lane: live-provider-read (POST /listings/refresh vs real connected ML installation)
stack: hub docker:dev (backend built from branch mis-003/m-01-listings-read-spine, bind-mount /workspace)
api_base: http://127.0.0.1:8080
installation_id: inst-mercado_livre-d373dc64-577f-4950-b11b-b90244f30cb2
tenant: tenant_default
external_account: 691607102 (METALNOBREACABAMENTOS)
date: 2026-07-16
primary_outcome: PASS (run succeeded + rows>0 tenant-scoped)
secondary_outcome: could-not-observe-via-API (real MLB ids + <20% unknown) — blocked by a read-path defect (see finding)
self_booted_server: false
env_dot_env_in_session: false
```

## Drive (all against the hub stack only — no self-booted server, no session `.env`)

1. **Baseline** `GET /listings/summary?installation_id=…` → `200 {"total":0,...}` (clean slate).
2. **Trigger** `POST /listings/refresh {installation_id}` → **`202 {"operation_run_id":"op_5d8df1910796888d56c9d97e4ea9bd83"}`**.
3. **Run status** `GET /integrations/installations/{inst}/operations` → run `op_5d8df19…`:
   `operation_type=listings_refresh`, **`status=succeeded`**, `result_code=LISTINGS_REFRESH_SUCCEEDED`,
   `failure_code=""`, `attempt_count=1`, started 03:26:18Z → completed 03:26:25Z (~6s).
   **Token renewed on first use — no 401/requires_reauth.**
4. **Post-refresh** `GET /listings/summary?installation_id=…` → **`200 {"total":34,"active":10,"paused":17,"exceptions":{"sync_error":0,"stale":0,"unlinked":33,"below_margin_worst_case":null,"margin_unknown":null}}`**.

## C10 assertions

| Assertion | Result | Evidence |
|-----------|--------|----------|
| run `succeeded` | **PASS** | op_5d8df19… `status=succeeded` `LISTINGS_REFRESH_SUCCEEDED` |
| rows > 0 tenant-scoped | **PASS** | summary `total=34` (tenant_default, this installation); baseline was 0 → 34 ingested by the real ML pull |
| real `MLB…` provider ids | **COULD-NOT-OBSERVE via API** | `GET /listings` + `/listings/by-product` → `503 source_unavailable` (read-path defect below) |
| `<20%` unknown status | **COULD-NOT-OBSERVE via API** | same 503; summary breaks out active/paused but not unknown/closed (34 total − 10 active − 17 paused = 7 other-status, not itemized) |

Note the summary already exposes the D-22 field `below_margin_worst_case` (+ `margin_unknown`) — confirms the ratified naming is live.

## FINDING (live-QA-caught) — List/ByProduct 503 when Oracle ICMS ceiling reader is down

`ReadService.List` (`read_service.go:122-125`) and `ByProduct` (`:157-159`) call
`facts.GetICMSCeilingByOrigin` and **hard-fail** (`fmt.Errorf("read ICMS ceiling: %w")` → handler
`503 source_unavailable`) when it errors. The concrete reader is **Oracle-backed**
(`internal_read/adapters/oracle/icms_ceiling.go`), and the dev stack has no Oracle wired
(`MPC_ORACLE_USERNAME is required`). By contrast `Summary` (`:59-67`) treats the SAME ceiling
error as `ceilingErr` and **degrades** — sets `BelowMarginWorstCase/MarginUnknown = nil` and
returns the row. Hence: `summary` `200`, `list`/`by-product` `503`.

**Assessment:** this couples the entire listings read spine's availability to Oracle for an
**optional** margin field. Under ADR-17 (unknown operational facts → null, never block) and the
Summary precedent in the same service, List/ByProduct should degrade to `below_margin_worst_case:
null` when the ceiling/cost reader is unavailable, not `503`. Hermetic P4–P7 lanes passed C04/C05
because they inject a working fake `CostReader`; live QA with the real Oracle reader down surfaced
the gap — exactly what P7 live QA is for.

## Disposition (hub ruling requested)

C10 **primary** (live provider read ingestion → run succeeded + 34 real tenant-scoped rows) is
**satisfied**. To complete the **secondary** data-quality assertions (real `MLB…` ids, `<20%`
unknown), one of:

- **(a) FIX (recommended):** correct `List`/`ByProduct` (and per-row `enrich`) to degrade to null
  margin when the ceiling/cost reader is unavailable, matching `Summary` + ADR-17; then re-drive —
  list serves 34 rows → observe ids + unknown%. This is a real M-01 read-spine robustness defect.
- **(b) WIRE Oracle** into the dev stack so the ceiling reader succeeds; then re-drive.
- **(c) DB-side verify:** hub runs one read-only `SELECT provider_listing_id, status FROM listings
  WHERE tenant_id='tenant_default' AND installation_id='…'` to confirm `MLB…` ids + unknown ratio
  (chip does not touch the hub DB/seam).
