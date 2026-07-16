# Live re-drive @ a6878dc6 (post-slice-11) — C10 + non-regression

```yaml
tip: a6878dc6 (slice 11, not pushed)
stack: hub-owned docker compose (chip never booted a server; hub sent ENV-READY)
base_url: http://127.0.0.1:8080
installation_id: inst-mercado_livre-d373dc64-577f-4950-b11b-b90244f30cb2
migrations_applied_on_boot: 0   # correct — slice 11 carries no migration
health: /healthz 200, server up 12:54:48
as_of: 2026-07-16T12:55:47.09492246Z
```

## Why this drive exists

Slice 11 tightened the degrade path (G1: cancel/timeout no longer degrade; G2: the Get
matrix guard reverted to `ceilingErr != nil`). Both changes make the read *stricter*.
The risk they carry is over-correction: a stricter degrade gate that turns a healthy read
into a blanket 503, or a reverted guard that stops emitting the matrix at all. This drive
exists to falsify that, against the live provider, not the mocks.

## HTTP results — all 200

| path | status |
|---|---|
| `/healthz` | 200 |
| `/listings` | 200 |
| `/listings/summary` | 200 |
| `/listings/by-product` | 200 |
| `/listings?has_exception=false` | 200 |
| `/listings?has_exception=true` | 200 |
| `/listings?exception=below_margin` | 200 |

The three fact-dependent filters are the ones slice 11 could have broken: they are the paths
that call `scan`/`scanGroups`, which propagate rather than degrade. All three serve 200 with
the source healthy — the stricter gate did not become a blanket 503.

## C10 — PASS (unchanged from the post-slice-9 drive)

```
rows: 34 | next_cursor: null
status distribution: paused 17, active 10, under_review 7
unknown STATUS: 0 / 34 = 0.0%   (contract bar: <20%; blocking failure: >20% unmapped)
provider ids: 34/34 match ^MLB  (real ids, no synthetic rows)
```

Summary is byte-identical to the slice-9 and slice-10 drives:

```json
{"total":34,"active":10,"paused":17,"exceptions":{"sync_error":0,"stale":0,"unlinked":33,"below_margin_worst_case":0,"margin_unknown":1},"as_of":"2026-07-16T12:55:47.09492246Z"}
```

## Fact-dependent filter counts

```
filter.has_exception=true      -> 33 rows
filter.has_exception=false     ->  0 rows
filter.exception=below_margin  ->  0 rows
```

> **Param-shape correction (made during P7 QA, applied here rather than left to mislead).**
> This table originally wrote the keys bare (`has_exception=true`). The API's filter params are
> namespaced — `filter.<key>=<value>` — and a **bare key is not a lenient alias: it is an
> unrecognized query key, silently ignored.** `status=paused` returns all 34 rows; `filter.status=paused`
> returns 17. The counts above are correct and were re-verified through a fresh browser at
> `c89fae3d` under the namespaced shape, but the labels as first written would have led a reader who
> copied them to get 34 rows back and conclude something had regressed.

33 + 0 = 33 against 34 total. This is **not** an arithmetic bug — it is ADR-17 working.
The one row neither filter claims is `MLB4735328201`, the single linked row (`unlinked` is
33 of 34). Its margin is unknown, so `matchesDependentFilter` declines it in both directions:
`has_exception=true` requires a *provable* exception and it has none provable;
`has_exception=false` requires `BelowMarginWorstCase != nil && !active`, and the value is nil.
A row whose exception status is unknown is honestly claimed by neither answer. Forcing it into
either bucket would be exactly the zero/default substitution ADR-17 forbids.
`matchesDependentFilter` was not touched by slices 10 or 11; this behavior predates both and
was certified at the round-2 gate.

## G2 falsification — the reverted guard, proved live

The regression slice 11 reverted was `Get` nulling the ENTIRE `ICMSWorstCaseByUF` when the
cost fetch failed. `GET /listings/…~MLB4735328201~-` → 200:

```
cost: {"amount":"91.57","currency":"BRL"}
below_margin_worst_case: null
icms_worst_case_by_uf: 25 entries, e.g.
  {"destination_uf":"1","worst_case_icms_pct":"20","price_net_basis":"135.992","below_margin_at_uf":null}
  {"destination_uf":"11","worst_case_icms_pct":"17","price_net_basis":"141.0917","below_margin_at_uf":null}
```

This is the shape the round-2 gate certified and G2 had broken: the matrix is **present**,
with `worst_case_icms_pct` and `price_net_basis` populated (they derive from price and
ceiling alone), and only `below_margin_at_uf` null (it needs the margin, which is unknown).
`below_margin_worst_case: null` at the row level agrees, and the summary's
`margin_unknown: 1` counts exactly this row. Known-true facts are served; the one unknown
fact is null. Under G2 this entire array would have been `null` — 25 known ICMS ceilings
discarded because an unrelated fact was missing.

## Verdict

C10 **PASS** (0.0% unknown, bar <20%). No regression across any read path. Slice 11's
stricter degrade gate did not over-correct, and the G2 revert is confirmed on live data,
not only in the unit matrix.
