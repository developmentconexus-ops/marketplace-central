# HUB EVENT — ESCALATION — C10 unknown-status 20.6% > 20% (adapter mapping gap)

```yaml
event: ESCALATION
from_chip: M-01-listings-read-spine
to: hub (local_efa46c30-1c0c-4075-9671-c2d7ae9efabe)
date: 2026-07-16
branch: mis-003/m-01-listings-read-spine
tip: a595f36c
lane: C10 live re-drive (post slice8 fix+wiring)
blocks: C10 pass → dual-gate DELTA → P8 CLOSED
```

## What passed on the C10 re-drive (ENV-READY @ a595f36c)

Read-only probes vs `http://127.0.0.1:8080`, installation
`inst-mercado_livre-d373dc64-577f-4950-b11b-b90244f30cb2`, tenant_default:

- `GET /listings/summary` → **200** `{"total":34,"active":10,"paused":17,"exceptions":{"sync_error":0,"stale":0,"unlinked":33,"below_margin_worst_case":0,"margin_unknown":1}}`
- `GET /listings?...&limit=100` → **200** (was 503 pre-slice → **fix (a) degrade confirmed live**). 34 rows, **34/34 real `MLB…` ids**.
- `GET /listings/by-product` → **200**, groups served.
- **Oracle cost wiring PROVEN LIVE:** linked row `MLB4735328201` (pid 15956) → `cost=91.57` populated from real Oracle (NO-STUB wiring works). 33 unlinked rows → cost null (legit, ADR-17). 1 linked row → `margin_unknown` (cost present, ceiling/policy absent → honest null, C07-consistent).

## The blocker — C10 secondary FAILS

- Status distribution: **active 10 / paused 17 / unknown 7** → **unknown = 7/34 = 20.6%**.
- Contract M01-C10 (`validation-contract.md:172,175`): expected `<20%` unknown status; **blocking failure = ">20% unmapped statuses (adapter mapping gap)"**.
- **20.6% > 20% → contract blocking failure.** Not the read defect (that path now serves 200). Not margin (that's C07 null-honesty, satisfied). This is a genuine **status-mapping** gap.

## Root cause (repo truth, exact)

`apps/server_core/internal/modules/listings/adapters/connectors/mapper.go:102-113`:

```go
func canonicalListingStatus(providerStatus string) listingsdomain.ListingStatus {
    switch strings.ToLower(strings.TrimSpace(providerStatus)) {
    case "active":  return listingsdomain.ListingStatusActive
    case "paused":  return listingsdomain.ListingStatusPaused
    case "closed":  return listingsdomain.ListingStatusClosed
    default:        return listingsdomain.ListingStatusUnknown
    }
}
```

Canonical enum (`domain/listing.go:13-18`) = **only** `active | paused | closed | unknown`.
Mercado Livre `item.status` has more real values than 3 — documented set includes
`under_review`, `inactive`, `payment_required`, `not_yet_active` (plus `active/paused/closed`).
Every ML status outside the 3 mapped cases falls to `unknown`. With 7/34 landing there,
the tenant's real inventory carries ≥1 ML status the mapper doesn't recognize → 20.6%.

## Why this ESCALATES (not a unilateral corrective slice)

1. **Data-blind.** The list endpoint exposes canonical only. I need the **raw provider status
   distribution of the 7 unknown rows** to fix data-driven — and chip has NO DB/.env/self-boot
   access (harness guardrail). Only the hub stack / DB-specialist session can run the read-only
   `SELECT status, count(*) FROM listings WHERE canonical→unknown GROUP BY status`.
2. **Contract/product decision, not mechanical.** Mapping `under_review`/`inactive`/`payment_required`/
   `not_yet_active` requires a ruling: do they collapse into existing canonical values (e.g.
   `inactive`→`closed`? `under_review`/`payment_required`/`not_yet_active`→`paused`?), or does the
   cockpit status enum **grow** new values? Growing the enum ripples: `domain/listing.go` +
   `IsValid()` + `mapper.go` + transport filter allowlist (`transport/query*.go` accepts
   `active|paused|closed|unknown`) + OpenAPI status enum + SDK + tests. That's a design call the
   wireframe/contract owner makes, not the read-spine chip.
3. **Scope boundary.** `mapper.go` is the **F-01 ingestion adapter** (status set at refresh). C10 is
   an M-01 contract but the *mapping* lives in the ingestion seam. Need hub ruling whether the fix
   is an M-01 corrective slice (like slice 8) or an F-01 reopen.

## Requests to hub

- **R1 (data):** run read-only, tenant-scoped, on the live stack (or via DB-specialist session):
  `SELECT status AS raw_ml_status, count(*) FROM listings WHERE tenant_id=<default> AND installation_id='inst-mercado_livre-d373dc64-577f-4950-b11b-b90244f30cb2' GROUP BY status ORDER BY 2 DESC;`
  → returns the exact ML status strings behind the 7 unknown (and confirms active/paused counts).
- **R2 (ruling):** the canonical mapping for each unmapped ML status — collapse into existing
  `active|paused|closed` **or** grow the enum. This is the contract/wireframe decision.
- **R3 (scope):** M-01 corrective slice (mapper+domain+transport+OpenAPI+SDK+tests, worker test-first)
  vs F-01 reopen. Either way I orchestrate; not writing >10 behavioral lines inline.

## Held state (nothing regressed)

- Slice 8 committed a595f36c (code) + d50636e7 (evidence). fix (a) + NO-STUB wiring live-proven.
- Dual-gate DELTA + P8 CLOSED **held** — cannot declare C10 pass at 20.6%. No push.
- Everything else in C10 is green (200s, 34/34 MLB ids, Oracle cost live). Only the status-mapping
  ratio blocks.
```
