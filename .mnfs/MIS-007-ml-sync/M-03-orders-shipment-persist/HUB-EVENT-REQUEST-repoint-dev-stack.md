# HUB EVENT — REQUEST — repoint-dev-stack (serve M-03 worktree binary for browser QA)

```yaml
event: REQUEST
from_chip: M-03-orders-shipment-persist
to: hub
date: 2026-08-01
branch: claude/sleepy-perlman-d0d325
tip: d22d3d20
depends_on: hub docker:dev stack (bind-mount /workspace = this worktree)
```

## Ask

Re-point (or restart) the `docker:dev` backend/frontend containers so the bind-mounted
`/workspace` serves **this worktree** (`C:\Users\leandro.theodoro\Documents\marketplace-central\.claude\worktrees\sleepy-perlman-d0d325`)
at tip **d22d3d20**, not the hub's main checkout.

Checked via `docker inspect marketplace-central-backend-1 --format '{{json .Mounts}}'`
(mounts-only, per the "no full docker inspect" finding — avoids dumping secrets): the
running `marketplace-central-backend-1`/`marketplace-central-frontend-1` containers currently
bind-mount `C:\Users\leandro.theodoro\Documents\marketplace-central` (the hub's own checkout),
confirmed healthy and up ~4h — but that is NOT this worktree, so it does not reflect F-01/F-02/F-03
(commits `cf2c09e1`..`d22d3d20`). Chip does not boot servers or rebuild images/re-point mounts
(hub seam) — per the repo's own documented gotcha (`dev-stack mount cwd-relative`), this needed
a hub REQUEST rather than a self-fix.

## Why

M-03's validation contract (`validation-contract.md`) mandates browser-driven QA that only a real
running stack can satisfy — this is NOT optional under the LEAN code-review model this milestone
otherwise used (that model governs adversarial-review headcount, not the operator's standing
user-drive validation mandate):

- **M03-C6** (Q1 perf): open a real ingested order's drawer, DevTools-measure `GET /orders/{id}` <2s,
  no live ML request in the waterfall.
- **M03-U1**: drawer shows comprador fiscal + shipment from Postgres — verified by killing ML network
  reachability in the stack and confirming the drawer stays complete.
- **M03-U2**: `/pedidos` list + summary before/after, 4 KPI bucket counts, tenant/window held constant,
  zero unexplained delta (any bucket move must be attributed order-by-order to the real shipment
  status the two-tier `GetOrderBucketCounts` read now surfaces, per `order_repo.go`'s F-03 change).
- **M03-U3**: browser console/network capture on both driven screens — zero requests to
  `api.mercadolibre.com`, no new console errors.

All four require the actual F-03 binary (Postgres-backed `ShipmentReader`/`BuyerFiscalReader`,
zero live-ML enrichment) to be what's actually running — the currently-served hub-checkout binary
predates this milestone entirely, so driving against it now would produce false/meaningless evidence.

## What I need from you

1. Re-point `docker:dev`'s bind mount at this worktree (or restart with `--build` against it),
   keeping the same Postgres (order_shipments/orders_marketplace_orders already migrated, 0088/0089
   applied) and the connected Mercado Livre installation (for the "kill ML reachability" step of
   M03-U1, I'll need to know how to simulate that against your existing setup — e.g. a firewall
   rule/env toggle you already use for this, rather than me improvising one).
2. Signal back (ENV-READY or similar) once the container is serving tip `d22d3d20`.
3. Confirm at least one order in the connected tenant has been ingested via F-02's `IngestOrder`
   (not just the old batch importer) so `order_shipments`/the 0089 buyer_* columns are actually
   populated for the row I drive against — if none exists yet, let me know and I'll request a
   live-drive ingest as a separate step (not a self-boot, just calling the already-wired
   `POST /orders/import` endpoint once the stack is up).

## Status while waiting

All non-browser-driven contract criteria (M03-C1 through M03-C5) already have code/test-level
evidence gathered from this milestone's adversarial reviews (per-feature) and the cold whole-diff
review — I'm writing that portion of the evidence pack now (`_chip-m03/EVIDENCE.md`) and will fold
in M03-C6/U1/U2/U3 once the stack is re-pointed. Not blocked on anything else in the meantime.

## Guardrails

No self-booted server. No `.env` touched in this session. No push. No dev-stack mutation attempted
by this chip beyond this request.
