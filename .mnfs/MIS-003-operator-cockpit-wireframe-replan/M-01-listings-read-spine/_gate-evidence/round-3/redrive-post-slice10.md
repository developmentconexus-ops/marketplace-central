# Live re-drive post slice 10 — NO REGRESSION

```yaml
lane: non-regression re-drive (C10 + fail-honest filter paths)
backend: 7f5a1b8c (hub ENV-READY; "applied 0 migration(s)" = correct, slice 10 has no migration; /healthz 200; server up 11:57:29)
installation: inst-mercado_livre-d373dc64-577f-4950-b11b-b90244f30cb2 (tenant_default)
base: http://127.0.0.1:8080
as_of: 2026-07-16T11:59:19Z
driven_by: chip (canonical reads only — no DB, no .env, no self-boot)
verdict: PASS — no regression
```

## Read paths

| Request | HTTP |
|---|---|
| `GET /healthz` | 200 |
| `GET /listings?installation_id=…&limit=100` | 200 |
| `GET /listings/summary?installation_id=…` | 200 |
| `GET /listings/by-product?installation_id=…&limit=20` | 200 |
| `GET /listings?…&filter.has_exception=false` | 200 |
| `GET /listings?…&filter.has_exception=true` | 200 |
| `GET /listings?…&filter.exception=below_margin` | 200 |

All three fact-dependent filters serve **200** with the fact source healthy — slice 10 did NOT
over-correct into a blanket 503. The fail-honest arm is reachable only under a genuine
`ReadErrorSourceUnavailable`, which cannot be provoked from the canonical surface; it is pinned
deterministically instead (unit tests + `slice10-review.md`).

## C10 — still PASS, byte-identical to the post-slice-9 drive

34 rows, **34/34 real `MLB…` provider ids**, `next_cursor: null`.

| status | count | share |
|---|---:|---:|
| paused | 17 | 50.0% |
| active | 10 | 29.4% |
| under_review | 7 | 20.6% |
| **unknown** | **0** | **0.0%** |

`unknown = 0.0% < 20%` → **C10 PASS holds** (contract `validation-contract.md:172,175`).

Summary unchanged vs post-slice-9:
`{"total":34,"active":10,"paused":17,"exceptions":{"sync_error":0,"stale":0,"unlinked":33,"below_margin_worst_case":0,"margin_unknown":1}}`

## Filter row counts — the one thing worth explaining

| Filter | Rows |
|---|---:|
| `has_exception=true` | 33 |
| `has_exception=false` | 0 |
| `exception=below_margin` | 0 |

33 + 0 = 33, against 34 total. **One row appears in neither filter — and that is correct, not a
regression.** The row is `MLB4735328201` (the single linked one): `link.state=resolved`,
`sync_state=synced`, `cost=91.57` (real Oracle), but `below_margin_worst_case=null` because the
ICMS ceiling / pricing policy is absent for it (`margin_unknown`, the same row the summary counts).

- `has_exception=true` requires a **provable** exception (`sqlException || below_margin==true`) → this
  row has neither → excluded.
- `has_exception=false` requires `BelowMarginWorstCase != nil && !active` → margin is unknown, so the
  service refuses to assert "this row has no exception" → excluded.

A row whose exception status is genuinely **unknown** is claimed by neither side of a boolean filter.
That is ADR-17 exactly: unknown never collapses into a confident `false`. `matchesDependentFilter`
(`read_service.go:488-500`) is untouched by slice 10 — this is pre-existing, intended, and it is the
same principle slice 10 exists to enforce (a known-false or unknown fact must never be returned as a
match). The summary surfacing `margin_unknown: 1` and the filter declining to classify that row are
consistent, not contradictory.

## Other C10 facts (unchanged, still green)

- Oracle cost wiring live-proven: `MLB4735328201` (pid 15956) → `cost=91.57` real Oracle.
- 33 unlinked → `cost` null (legit ADR-17); 1 linked → `margin_unknown` (cost present, ceiling/policy
  absent → honest null, C07-consistent).
