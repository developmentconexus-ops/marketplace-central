# HUB EVENT — REQUEST — restart-backend (serve slice8 binary)

```yaml
event: REQUEST
from_chip: M-01-listings-read-spine
to: hub (local_efa46c30-1c0c-4075-9671-c2d7ae9efabe)
date: 2026-07-16
branch: mis-003/m-01-listings-read-spine
tip: a595f36c
depends_on: hub docker:dev stack (bind-mount /workspace = this worktree)
```

## Ask

Restart the `docker:dev` backend container so its entrypoint recompiles the bind-mounted worktree
and serves the slice8 binary (commit **a595f36c**). Keep the same migrated Postgres (38 tables) +
the connected Mercado Livre installation + the Oracle creds you set (cost/policy now wired real).
Re-signal ENV-READY when the new binary is up.

## Why

Slice 8 landed (P7e): read path degrades on optional-fact outage (fix a) AND composition root now
wires the real cost (`oracleDB`) + policy (`marketSvc`) readers (NO-STUB). The running container still
holds the pre-slice binary (List/ByProduct 503). A container restart is needed to pick up a595f36c —
chip does not boot servers or rebuild images (hub seam).

## After ENV-READY — chip re-drives C10 complete (no self-boot, no session .env)

1. `POST /listings/refresh {installation_id}` (or reuse the 34 already ingested) → confirm run succeeded.
2. `GET /listings?installation_id=…` → now **200** (degrade fix + real readers) → observe real `MLB…`
   provider ids + compute `<20%` **status**-unknown ratio (C10 secondary; status is Postgres-backed).
3. `GET /listings/by-product` → 200, groups served.
4. With Oracle wired, cost/margin should now populate (or degrade cleanly per-row if Oracle hiccups).

Then: dual-gate DELTA over commits since `e2cde36` (never-downgrade) → re-gate → **P8 CLOSED**.

## Guardrails

No self-booted server. No `.env` in session. No push. `docker/dev/*.sh` (your CRLF→LF env prep) +
worktree `.env` remain uncommitted in the worktree — commit a595f36c deliberately excluded them.
