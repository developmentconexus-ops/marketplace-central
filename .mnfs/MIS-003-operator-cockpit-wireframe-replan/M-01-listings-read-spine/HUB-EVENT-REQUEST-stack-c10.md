# HUB EVENT — REQUEST — M-01 C10 live stack

```yaml
event: REQUEST
from_chip: M-01-listings-read-spine
to: hub (local_efa46c30-1c0c-4075-9671-c2d7ae9efabe)
date: 2026-07-15
depends_on: HUB-EVENT-BLOCKED-C10.md (ruling: caminho (a) approved)
```

## Ask

Per hub ruling (caminho a), the chip needs the dev stack up to drive **C10 (live-provider-read)**
against a real connected Mercado Livre installation. The stack is a **hub seam** — chip does not
boot it. Requesting the hub:

1. Bring up the stack: `npm run docker:dev` (+ `npm run docker:oauth` for the OAuth flow).
2. Operator connects a Mercado Livre installation via the app OAuth flow (credential-entry =
   operator-only; chip never performs it).
3. Signal **env-ready** back to this chip (event) with the reachable API base + the connected
   `installation_id` (no tokens/secrets).

## Chip commitment on env-ready

Re-drive C10 against the **hub stack** only (no self-booted server, no session `.env`):
`POST /listings/refresh {installation_id}` → assert run `succeeded`, tenant-scoped
`SELECT count(*) > 0`, real `MLB…` provider ids, `< 20%` unknown status → capture sanitized
evidence (no tokens) → re-gate (fresh crew + live pass) → **P8 CLOSED**.

## Status

M-01 **held at Blocked (pending C10)**. C01–C09 PASS (P5+P6+P7, fresh evidence). C07 doc
reconciliation committed (`c3aebe12`). Awaiting hub env-ready signal. No push.
