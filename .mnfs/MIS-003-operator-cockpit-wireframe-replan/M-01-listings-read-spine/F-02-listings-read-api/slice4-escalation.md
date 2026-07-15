# F-02 Slice 4 — ESCALATION (contract gap, milestone-owner analysis)

by-product grouping (`GET /listings/by-product`). The Slice-2 hub ruling (Option 2, conditions a–e) + its IC-02 doc note resolved below_margin-dependent **filters on the flat list**. Slice 4 needs the same question answered for the **grouping surface**, which the existing ruling does not cover.

## The gap
IC-02 (interface contract line 52–53) says by-product uses the **same filter grammar** as `GET /listings` — including `exception=below_margin` and `has_exception=true|false`. But the hub-blessed Option-2 doc note (contract lines 64–70) is written entirely in **flat-list / row** terms: "bounded iterative keyset scan… `next_cursor` identifies the last scanned **row**… single-query `limit+1` path." by-product paginates over **group keys** (product_title, product_id) with a null-last "sem produto" group, not over rows. The row-scan resumption semantics do not map onto group-key pagination.

Two below_margin touchpoints on by-product, only ONE is unruled:
1. **group_state (DISPLAY) — already ruled, NOT blocked.** D-17: `group_state` = worst child = **attention** if any child below_margin. Every returned group must enrich its children (one bounded Oracle cost batch across the returned groups' children) to compute group_state. This is bounded per-page display enrichment, already blessed (D-20/Option-2 condition b). Implementable now.
2. **below_margin-dependent FILTER (exception=below_margin / has_exception) — UNRULED, blocks.** "Filters apply before grouping" (D-17/plan Slice 4): a below_margin child filter changes which children survive → changes `listing_count`, and a group with zero surviving children must not appear. Determining survival needs per-child Oracle cost, so selecting the `limit+1` **group keys** cannot be a single SQL statement — same over-constrained conflict as Slice 2, but now the accumulation unit is a group, not a row.

## Scope of the block (small, mirrors Slice 2)
Implementable NOW with zero new ruling:
- group keyset over resolved-product groups, product_title ASC + product_id tie-break, null-last "sem produto" group (at most once, globally last);
- all SQL-derivable filters + q applied before grouping (status, sync_state, link_state, listing_type_code, product_id, sync_error/stale/unlinked exceptions);
- complete children per returned group (one bounded keys query + one set-children query, no N+1);
- full child below_margin **display** enrichment + `group_state` (D-17 worst-child);
- group cursor, injected UTC as-of, IC-02 group envelope.

BLOCKED (needs a ruling): only `exception=below_margin` and `has_exception=true|false` **as filter predicates** on by-product.

## Options
1. **Extend Option-2 bounded scan to group keys.** Scan group-key pages; for each candidate group enrich its children, drop children failing the below_margin predicate, drop the group if zero children survive; accumulate up to `limit` surviving groups or the 50-page cap; `next_cursor` = last **scanned** group key (resume without skipping). Consistent with the flat-list ruling; the doc note gains a grouping paragraph (chip authors, hub reviews at acceptance). **Milestone-owner recommendation** — one coherent rule across both surfaces, preserves D-17 exact filtering + D-18 keyset + D-20 no-projection. Cost: sparse below_margin groups issue K group-key fetches + K child-cost batches, K bounded by the same named cap.
2. **below_margin filters unsupported on by-product** — reject `exception=below_margin` / `has_exception` on this surface with 400 `invalid_filter`. Simplest, but DIVERGES the grammar IC-02 explicitly shares between the two surfaces (line 52–53 "same as above"); breaks the "Com pendência" grouping view.
3. **Defer by-product below_margin filter to a later milestone** — document as accepted-but-unsupported-in-M01 on by-product only; ship group_state display + all other filters now. Smaller than Option 1, still an IC-02 note, keeps grammar nominally shared but with a documented M01 carve-out.

## Requested
Ruling on 1/2/3 (recommend 1 — symmetry with the flat-list Option-2 ruling already in force). On any choice I dispatch ONE Slice 4: the unblocked grouping scope + the chosen below_margin-filter path. No source changed yet. Prereqs landed: group cursor/query/envelope (Slice 1), D-17 severity, D-22 below_margin, Slices 2+3 committed (d9e4737d, 57fba22d) with the flat-list scan path as the reference implementation.
