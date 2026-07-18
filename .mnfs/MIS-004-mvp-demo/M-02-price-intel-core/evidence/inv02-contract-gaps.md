## Q1 — IC-06 concrete Go shapes

IC-06 fixes the seven operation names and top-level normalized members:

- `GetOwnItemPricing`: `ItemID`, `Amount`, `Currency`, `RegularAmount`, `FetchedAt` — [IC-06:27](.mnfs/MIS-004-mvp-demo/research/ml-read-ports-interface-contract.md:27)
- `GetPriceToWin`: `ItemID`, `Status`, `CurrentPrice`, `TargetPrice`, `Position{Rank,Total}`, `FetchedAt` — [IC-06:28](.mnfs/MIS-004-mvp-demo/research/ml-read-ports-interface-contract.md:28)
- `SearchCatalogByEAN`: `[]CatalogProduct{CatalogProductID, Title, Attrs}` — [IC-06:29](.mnfs/MIS-004-mvp-demo/research/ml-read-ports-interface-contract.md:29)
- `GetCatalogProduct`: `ID`, `Title`, `BuyBoxWinner{Price,SellerID}`, `Attrs`, `FetchedAt` — [IC-06:30](.mnfs/MIS-004-mvp-demo/research/ml-read-ports-interface-contract.md:30)
- `ListCatalogOffers`: paginated `[]Offer{SellerID, Price, Condition, ShippingMode}` — [IC-06:31](.mnfs/MIS-004-mvp-demo/research/ml-read-ports-interface-contract.md:31)
- `GetShipmentInfo`: `ID`, `Status`, `SLADue`, `Delayed`, `Costs`, `DestinationUF`, `FetchedAt` — [IC-06:32](.mnfs/MIS-004-mvp-demo/research/ml-read-ports-interface-contract.md:32)
- `GetFreeShippingCost`: `Cost`, `FetchedAt` — [IC-06:33](.mnfs/MIS-004-mvp-demo/research/ml-read-ports-interface-contract.md:33)

Known repository convention: monetary values use `Money{Amount string, Currency string}`, not cents or a decimal dependency — [seam map:334-342](.mnfs/MIS-004-mvp-demo/M-02-price-intel-core/evidence/inv01-seam-map.md:334).

The blocking details are not latent:

- `Attrs` has no member schema. The research only proves that catalog search returns products and that identity requires normalized attributes; it does not capture attribute `name`, `value`, `id`, or nullability — [handoff:43-51](docs/research/2026-07-17-pricing-intelligence-implementation-handoff.md:43).
- `Costs` is explicitly `{...}`; no cost member names, monetary types, or nullability are defined — [IC-06:32](.mnfs/MIS-004-mvp-demo/research/ml-read-ports-interface-contract.md:32).
- The `item` request passed to `GetFreeShippingCost(userID, item)` has no field shape — [IC-06:33](.mnfs/MIS-004-mvp-demo/research/ml-read-ports-interface-contract.md:33).
- `SearchCatalogByEAN` has `FetchedAt` omitted from the list result while all seven operations are said to return it — [IC-06:23](.mnfs/MIS-004-mvp-demo/research/ml-read-ports-interface-contract.md:23), [IC-06:29](.mnfs/MIS-004-mvp-demo/research/ml-read-ports-interface-contract.md:29). Neither a list envelope nor per-item timestamp is specified.

**VERDICT: UNDERSPECIFIED — HUB RULING NEEDED.** Escalate: “Freeze `Attrs`, shipment `Costs`, free-shipping `item`, and the `SearchCatalogByEAN` `FetchedAt` placement and Go types.” The captures establish endpoint behavior, not these shapes; choosing them would invent the ML-facing contract.

## Q2 — IC-01 resolver decision table

The resolver is explicitly deterministic:

- Two independent agreeing anchors are required for `ACCEPT`; title/fuzzy matching only ranks candidates — [IC-01:34-36](.mnfs/MIS-004-mvp-demo/research/identity-matching-interface-contract.md:34).
- Contradictory kit/combo, color, measure/dimension, or voltage rejects even with equal EAN — [IC-01:37](.mnfs/MIS-004-mvp-demo/research/identity-matching-interface-contract.md:37).
- Missing EAN caps the result at `REVIEW` — [IC-01:38](.mnfs/MIS-004-mvp-demo/research/identity-matching-interface-contract.md:38).
- States are `ACCEPT`, `REVIEW`, `REJECT`, `NO_CANDIDATE` — [IC-01:41-44](.mnfs/MIS-004-mvp-demo/research/identity-matching-interface-contract.md:41).
- Required normalization vocabulary includes GTIN/checksum, exact GTIN equality, domain/family/category, brand, reference/model, line, attributes, measure, capacity, color, finish, voltage, kit, unit, condition, duplicate GTIN, and conflicting identity — [handoff:93-101](docs/research/2026-07-17-pricing-intelligence-implementation-handoff.md:93).
- Confidence is an API band, not a required exact score: `ALTA|MEDIA|BAIXA`, corresponding to UI thresholds ≥85, 50–84, and <50 — [IC-01:39](.mnfs/MIS-004-mvp-demo/research/identity-matching-interface-contract.md:39).

Therefore fixtures should assert the deterministic `decision` plus `confidence_band`; they should not assert a learned numeric confidence value or require anchor weights. The batch plan’s required fixtures cover missing EAN, zero candidates, equal candidates, EAN collisions, hard negatives, empty attributes, and title-only matches — [batch plan:479-489](.mnfs/MIS-004-mvp-demo/M-02-price-intel-core/evidence/p2-batch-plan.md:479).

**VERDICT: RESOLVED.** Implement a rule gate with `ACCEPT` only for two independent agreeing anchors; `REJECT` for determinative hard negatives; `REVIEW` for missing/colliding EAN, insufficient evidence, ties, or non-hard contradictions; `NO_CANDIDATE` for zero candidates. Fixtures assert the decision and `ALTA|MEDIA|BAIXA` band, not an invented numeric score.

## Q3 — verdict threshold policy ownership

What is frozen:

- IC-03 fixes the four labels and says `SEM_CUSTO` serves market evidence without `verdict_label` — [IC-03:28](.mnfs/MIS-004-mvp-demo/research/market-evidence-read-interface-contract.md:28).
- IC-03 requires known cost/fees/freight/tax for a verdict; unknown inputs remain unknown — [IC-03:28](.mnfs/MIS-004-mvp-demo/research/market-evidence-read-interface-contract.md:28).
- M-07 owns `CalcProfile`, including `limiar_verde_pct` and `limiar_amarelo_pct`, and its decomposition formula — [IC-04:25-30](.mnfs/MIS-004-mvp-demo/research/pricing-difal-interface-contract.md:25).
- M-07 validation explicitly treats ≥18/≥10 as CalcProfile thresholds — [M-07 validation:95-100](.mnfs/MIS-004-mvp-demo/M-07-simulador/validation-contract.md:95).
- M-07’s formula requires commission and freight from live readings and propagates unknown cost, freight, commission, or UF — [IC-04:27-30](.mnfs/MIS-004-mvp-demo/research/pricing-difal-interface-contract.md:27).
- M-02 F-04 nevertheless describes “margin … vs ERP cost ⇒ label + range” without assigning ownership or defining the table — [F-04:33-42](.mnfs/MIS-004-mvp-demo/M-02-price-intel-core/F-04-collect-verdict-api/feature.md:33).

What is not frozen: no artifact assigns the market-range/margin threshold table to M-02 or M-07, and no M-02 contract defines the complete fee, freight, tax, source-time, and profitability-input policy. The batch plan records this exact gap — [batch plan:184-192](.mnfs/MIS-004-mvp-demo/M-02-price-intel-core/evidence/p2-batch-plan.md:184).

**VERDICT: UNDERSPECIFIED — HUB RULING NEEDED.** Escalate: “Does M-02 emit only market range plus `verdict_label: null`/`SEM_CUSTO` unless an existing M-07 `CalcProfile` is supplied, or does M-07 own and publish the threshold/policy table for M-02 to consume? Which exact fee, freight, tax, cost, and source-time inputs are authoritative?” M-02 cannot answer without duplicating or re-deciding M-07 pricing semantics.

## Q4 — concrete public JSON shapes

IC-03 fixes these concrete entity fields:

- `MarketPriceSnapshot`: `scope`, `ref_id`, `source`, nullable `price`, `currency`, `observed_at`, `fetched_at`, `expires_at`, `status`, nullable `failure_reason`, `request_id` — [IC-03:25](.mnfs/MIS-004-mvp-demo/research/market-evidence-read-interface-contract.md:25).
- `CompetitiveSignal`: `listing_id`, `our_price`, nullable `winner_price`, nullable `target_price`, nullable `position{rank,total}`, `fetched_at` — [IC-03:26](.mnfs/MIS-004-mvp-demo/research/market-evidence-read-interface-contract.md:26).
- `MarketAggregate`: `median`, `min_valid`, `n_offers`, `n_sellers`, `computed_at`, `status` — [IC-03:27](.mnfs/MIS-004-mvp-demo/research/market-evidence-read-interface-contract.md:27).
- `Verdict`: `match_status`, `price_evidence_status`, nullable `verdict_label`, nullable `blocking_state`, `inputs_used` — [IC-03:28](.mnfs/MIS-004-mvp-demo/research/market-evidence-read-interface-contract.md:28).
- `POST /market/collections` is synchronous `200` with envelope members `status`, `decisões`, `contagens`, `causas` — [IC-03:35](.mnfs/MIS-004-mvp-demo/research/market-evidence-read-interface-contract.md:35).

Known additive JSON sketch:

```json
{
  "status": "COMPLETED",
  "decisões": [],
  "contagens": {},
  "causas": []
}
```

```json
{
  "match_status": "ACCEPT",
  "price_evidence_status": "OK",
  "verdict_label": null,
  "blocking_state": "SEM_CUSTO",
  "inputs_used": {},
  "market_range": {
    "min_valid": "string|null",
    "median": "string|null",
    "currency": "BRL",
    "n_offers": "integer",
    "n_sellers": "integer"
  }
}
```

The following semantics remain unfrozen:

- Member types and meanings for `decisões`, `contagens`, and `causas`.
- What `inputs_used` enumerates or whether it is an object, array, or named-source map.
- The exact JSON name and shape of the SEM_CUSTO market-price range. IC-03 names aggregate fields but does not define a `market_range` member or whether the range includes min/median/max.
- The `MarketAggregate` entity does not define a `max` field — [IC-03:27](.mnfs/MIS-004-mvp-demo/research/market-evidence-read-interface-contract.md:27).

The batch plan explicitly records these omissions — [batch plan:252-257](.mnfs/MIS-004-mvp-demo/M-02-price-intel-core/evidence/p2-batch-plan.md:252).

**VERDICT: UNDERSPECIFIED — HUB RULING NEEDED.** M-02 may own the additive OpenAPI/SDK serialization, but the hub must first freeze the meanings and types of `decisões`, `contagens`, `causas`, `inputs_used`, and the SEM_CUSTO range member. Otherwise the JSON would invent contract semantics and bind M-05/M-06 consumers.

## Summary

| Q | Verdict | One-line |
|---|---|---|
| Q1 | UNDERSPECIFIED | IC-06 names top-level outputs but omits four required concrete shapes. |
| Q2 | RESOLVED | Deterministic decisions plus confidence bands; no numeric score formula is required. |
| Q3 | UNDERSPECIFIED | Threshold ownership and profitability-input policy are not assigned. |
| Q4 | UNDERSPECIFIED | Named entities are fixed, but several public nested JSON meanings/types are not. |