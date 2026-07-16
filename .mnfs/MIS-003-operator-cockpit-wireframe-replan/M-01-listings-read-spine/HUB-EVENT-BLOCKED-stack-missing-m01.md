# HUB EVENT — BLOCKED — C10 stack missing M-01 endpoints

```yaml
event: BLOCKED
from_chip: M-01-listings-read-spine
to: hub (local_efa46c30-1c0c-4075-9671-c2d7ae9efabe)
date: 2026-07-16
depends_on: hub ENV-READY signal (stack docker up, API base http://127.0.0.1:8080)
```

## Finding

The hub `docker:dev` stack is **live and healthy** but was **built from a checkout WITHOUT the
M-01 code** — the M-01 listings endpoints are **not registered in the running backend binary**.
Probed against the hub stack (`http://127.0.0.1:8080`, installation
`inst-mercado_livre-d373dc64-577f-4950-b11b-b90244f30cb2`):

| Call | Hub stack result |
|------|------------------|
| `GET /healthz` | `200 {"status":"ok"}` (stack live) |
| `GET /integrations/installations/{inst}/operations` | `200` (pre-existing integrations route works — ML installation present, prior `pricing_fee_sync` runs `succeeded`) |
| `GET /listings?installation_id={inst}` | **`404 page not found`** |
| `GET /listings/summary?installation_id={inst}` | **`404 page not found`** |
| `POST /listings/refresh {installation_id}` | **`404 page not found`** |

All three M-01 routes 404 → the binary is a **pre-M-01-merge build** (main). M-01 is unmerged
(it is Blocked, not closed), so a `docker:dev` built from the hub's main checkout cannot contain
`/listings/*`. The route exists in the M-01 branch source (`http_handler.go` `POST /listings/refresh`,
confirmed).

## Why the chip cannot proceed

C10 must exercise the M-01 code. That code is only on branch
**`mis-003/m-01-listings-read-spine`** (tip `0468c21119e9b15125324c8edf3cc87084d0fda5`; endpoints
landed by milestone HEAD `e2cde36`). The chip **must not** boot its own server (doctrine
HARNESS §5 L2) — so it cannot supply the missing endpoints itself. Driving C10 requires the
**hub** to rebuild the dev stack backend from the M-01 branch.

## Ask (hub seam)

Rebuild the `docker:dev` backend image from branch `mis-003/m-01-listings-read-spine`
(tip `0468c211`) — e.g. point the compose build context / checkout at the worktree branch — so
`/listings/refresh` + the read endpoints are registered, keep the same migrated Postgres
(38 tables) + the connected ML installation, then **re-signal ENV-READY**. On the new signal the
chip re-drives C10 against the hub stack only (run succeeded, tenant-scoped `count(*)>0`, real
`MLB…` ids, `<20%` unknown) → re-gate → P8 CLOSED.

## Status

M-01 remains **Blocked (pending C10)**. No self-booted server. No `.env` in session. No push.
Token-expiry ressalva noted for the eventual real drive (single attempt; on 401/requires_reauth →
BLOCKED-C10-reauth, no retry loop).
