# F-02 Slice 2 — ESCALATION (contract conflict, milestone-owner analysis)

Worker (Sol-low) STOPPED honestly per the slice stop-rule; `slice2-notes.md` has the raw stop. This is the milestone-owner classification + recommendation for the hub.

## Conflict (over-constrained set — pick which constraint gives)
`GET /listings` filters `exception=below_margin` and `has_exception=true|false` require evaluating **below_margin** in the list predicate. below_margin (D-22) = f(listing price [PG], CUSSEMICM cost [Oracle, per product_id], ICMS ceiling [Oracle, global], min_margin [PG marketplaces policy]). The per-row Oracle **cost** input is not in PostgreSQL, and:
- **D-20** forbids a local persisted cost/below-margin projection (live Oracle batch only).
- **D-18** mandates keyset `limit+1` pagination; the F-02 plan requires "filter before truncation" and an EXPLAIN keyset-index proof at 2,000 rows.
- **IC-02** lists below_margin as a first-class list exception filter, and **D-17** puts below_margin in the active-issue set that `has_exception` must reflect **exactly**.

A single keyset `limit+1` SQL statement cannot filter on below_margin (input absent from PG); post-fetch filtering of one `limit+1` page returns short/empty pages and emits a cursor that skips later qualifiers → breaks pagination. The four ratified constraints cannot all hold for this predicate.

## Precise scope of the block (small)
Implementable NOW with ZERO contract change (no below_margin dependency):
- keyset list (`limit+1`, cursor tuple `(title, provider_listing_id, variation_id)`, tenant+installation scoped);
- all direct filters (`status`, `sync_state`, `listing_type_code`, `product_id`, `link_state`);
- 3 of 4 exception filters — `sync_error`, `stale`, `unlinked` — all SQL-derivable from PG (`sync_state`/`sync_error`/`fetched_at`/joined `link_state`);
- `q` bound escaped ILIKE (title / provider_listing_id / snapshot seller_sku);
- **below_margin worst-case DISPLAY** on returned rows (per-page enrichment: one Oracle cost batch + one global ceiling + one policy per emitted page — bounded, no projection);
- next-cursor, injected UTC as-of, IC-02 field-for-field response.

BLOCKED (needs a ruling): only the below_margin-**dependent filter predicates** — `exception=below_margin`, and exact `has_exception=true|false` (false requires below_margin evaluated).

## Options
1. **Defer the below_margin list filter** to a later milestone. M-01 ships below_margin as: per-row DISPLAY flag (list), summary counter (Slice 5), per-UF matrix (Slice 3, GET /listings/{id}) — but NOT as a list filter predicate. Smallest change; documents `exception=below_margin`/`has_exception` as unsupported-in-list for M-01 (IC-02 note).
2. **Iterative bounded keyset scan** for below_margin-filtered queries only (single-page fast path otherwise). Repository walks PG keyset pages, enriches each with its per-page Oracle cost batch + the one global ceiling + one policy, filters, accumulates until `limit` qualifiers or a scan cap; cursor encodes the last **scanned** row so resumption never re-scans or skips (short non-final pages allowed). Preserves D-18 keyset correctness, D-20 no-projection, IC-02 exact filter. Cost: a sparse below_margin filter may issue K PG fetches + K Oracle batches (K bounded by catalog size / scan cap). **Milestone-owner recommendation** — least contract disruption; standard technique for keyset + externally-computed predicate. Needs hub to bless (a) short-non-final pages under D-18 and (b) the multi-batch read as within D-20's "bounded" intent (D-20's single-batch language targeted the C09 summary, not the list).
3. **Authorize a tenant-scoped PG projection** of cost/rate/policy-derived below_margin, reversing D-20's no-projection ruling. Makes below_margin fully SQL-filterable + fast, but adds a sync/staleness surface and contradicts a standing hub decision.

## Requested
Ruling on 1/2/3 (recommend 2). On any choice I proceed: implement the unblocked Slice-2a scope + the chosen below_margin-filter path in one redispatched Slice 2. No source files changed yet; D-21 (min_margin seam) already landed + committed (1b644ed7).
