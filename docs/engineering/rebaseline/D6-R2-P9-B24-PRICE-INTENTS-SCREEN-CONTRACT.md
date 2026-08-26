# D6-R2 P9 — B24 Preços Screen Contract

> **Status:** RUN / **UPSTREAM FINDING RAISED — awaiting operator adjudication**; no PASS claimed
> **Block:** B24 — Preços / R24
> **Methods:** [DevelopmentConexus Engineering Method v1.0.0](../../development/engineering-method.md) + [Frontend Product Experience Planning Method v2.3](../../development/frontend-product-experience-planning-method.md)
> **Locked P8 evidence:** `qualification/d6-r2-wireframes/b24-price-intents.html` (revision 5, LOCKED 2026-08-26)
> **Canonical Product OAD:** `contracts/api/product/openapi.yaml` (107/31/H-A-S)
> **Product implementation:** BLOCKED UNTIL accepted D9

## 1. P9 result

P9 ran against the operator-LOCKED revision-5 candidate. The locked human job is:

```text
open the pricing workbench for the organization
→ read, per row, the current price, the Economics margin fact and the Market position fact
→ narrow to what needs attention (below policy / not at the lowest delivered price) — server facts only
→ type a candidate price; the owners re-answer (margin, and position in the comparable range)
→ open the cascade and the reference anchors when the trade-off needs proving
→ create exactly one explicit PriceIntent for that row, superseding any active one
→ read the honest external outcome; verify (never blindly retry) on ambiguity
```

**P9 verdict: BLOCKED on one upstream finding (§6). Everything else traces.**

## 2. Route, identity and client-state ownership

Route family: `/org/:organizationId/publicacoes/precos`, two fixed views (`Decidir` workbench / `Intenções` ledger) — not a saved-view platform.

| State class | B24 ownership |
| --- | --- |
| `GLOBAL_WORKSPACE_CONTEXT` | `organization_id` |
| `URL_NAVIGATION_STATE` | active view, account filter, attention filter |
| `SERVER_STATE` | listings, price intents, expected economics, competitive positions, scenario evaluations, effect state |
| `LOCAL_EPHEMERAL` | the unsent typed price, row disclosure state, mobile navigation |

The typed price is local and inert until it is either sent to an owner for evaluation or committed as an explicit intent. No pricing number displayed on the screen is produced by the screen (`data-client-computation="none"`).

## 3. Product operation / access binding

| Screen need | Operation | Semantic owner | Permission | Principal kinds |
| --- | --- | --- | --- | --- |
| workbench rows (listings) | `ListMarketplaceListings` | Offering | `offering.read` | H/A/S |
| active intent per row + ledger | `ListPriceIntents` | Offering | `offering.read` | H/A/S |
| read one intent / verify outcome | `GetPriceIntent` | Offering | `offering.read` | H/A/S |
| create intent (supersede-only) | `CreatePriceIntent` + Idempotency-Key | Offering | `price.manage` | H/A |
| current margin per row | `ListExpectedEconomics` | CommercialEconomics | `economics.read` | H/A/S |
| margin at the typed price + waterfall + anchor evaluation | `EvaluatePriceScenario` | CommercialEconomics | `economics.read` | H/A/S |
| current market position per row | `ListCompetitivePositions` | MarketIntelligence | `market.read` | H/A/S |
| comparable evidence behind the range | `ListComparableOffers` | MarketIntelligence | `market.read` | H/A/S |

The workbench composes these cursor collections **page-level**: one page of rows, one page of each owner collection, no N+1 and no screen-shaped aggregate endpoint.

## 4. What traces cleanly

- **Write path.** Exactly one consequential write (`CreatePriceIntent`), explicit per row, Idempotency-Key carried, supersede-instead-of-edit, no bulk apply, no repricing engine.
- **Honest outcome.** `pending / applied / rejected / ambiguous / superseded` plus `convergence`; rejection surfaces provider feedback verbatim with a substitute path; ambiguity is verification-only, never blindly retried.
- **Owner separation.** Economics answers margin and policy judgment; Market answers position; Offering owns the write. The screen sets no threshold and computes no money.
- **Honest population.** known-empty ≠ unknown ≠ unavailable, for the workbench, the ledger, economics conclusions and market evidence alike.
- **Attention filter.** Projects server-issued facts (`below_policy`, position not first), never a client score.
- **Anchors.** Owner-evaluated; the market anchor is gated by `evidence_sufficiency` and its absence is explained; clicking only fills the input.

## 5. Backward trace

No admitted operation above is orphaned, and no locked control lacks an operation — with the single exception recorded below.

## 6. UPSTREAM FINDING — the Market owner does not carry the locked position projection

The LOCKED screen presents, per row and live against the typed price:

```text
delivered price range  R$ 237,00 – R$ 289,00   (over N comparable offers)
our rank               3º de 11                (current published price)
our rank at the typed candidate price          (e.g. 6º de 11 at R$ 265,00)
```

The current contract cannot supply this:

| Locked need | Current contract | Gap |
| --- | --- | --- |
| delivered range low/high | `CompetitivePosition` carries `relation` + `delivered_price_gap` only | range absent |
| rank / comparable count | absent everywhere | rank absent |
| evidence sufficiency per workbench row | `CompetitivePositionListItem` omits `evidence_sufficiency` | the column cannot stay honest at collection level |
| position at a **candidate** price | `EvaluatePriceScenario` is CommercialEconomics and answers margin only | no Market scenario operation exists |

`ListComparableOffers` returns the individual delivered prices, so the numbers are *derivable* — but only by the screen taking min/max and counting, which is exactly the client computation this block's LOCK forbids and which the deterministic verifier now protects against. Deriving them in the client would also silently discard `coverage` and `evidence_sufficiency`, turning partial evidence into a confident-looking rank.

Per §3.10A the proven user need is not removed because the API does not yet answer it. The proposed bounded repair, for operator adjudication:

1. extend `CompetitivePosition` with an owner-issued `delivered_price_range` (low/high `Money`) and `market_rank` (`position` + `comparable_count`), both **absent** unless `evidence_sufficiency = sufficient`;
2. add `evidence_sufficiency` (and the same optional range/rank) to `CompetitivePositionListItem`, so the workbench column is honest and page-level composable;
3. admit one stateless Market scenario operation — `EvaluateCompetitivePositionScenario` (class C, `market.read`, H/A/S, owner MarketIntelligence) — answering the position of a **candidate** delivered price, mirroring what `EvaluatePriceScenario` does for Economics. Product surface would move 107 → 108 operations, Permissions unchanged at 31, Principal kinds unchanged.

Rejected alternatives: computing range/rank in the client (breaks the block LOCK and the honesty law); folding the market answer into `EvaluatePriceScenario` (collapses the Economics/Market owner boundary); dropping rank from the locked screen (backend-shaped UX, §3.10A).

**No P9 PASS is claimed for B24 until this finding is RATIFIED, REJECTED or DEFERRED by the operator.**
