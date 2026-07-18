# CHIP-M02 → HUB — contract-amendment DRAFT (IC-06 + IC-03)

Status: DRAFT for `REQUEST contract-amendment`. Chip drafts member tables (D-04: chip holds
research/payload context); hub commits to main. Per D-04 unblock rule, downstream F-01/F-03/F-04
proceed the moment this draft is **sent + hub-ACK'd** — not at main commit.

Frozen D-04 principles applied verbatim: Money = decimal string; `FetchedAt` on list ENVELOPE;
minimal consumed subset only; enum-keyed snake_case count maps; `inputs_used = {input:{source,as_of}}`;
`market_range` mirrors the frozen aggregate exactly (NO invented `max`).

Grounding: IC-06 top-level members stay FROZEN (not re-decided here). Live payload evidence —
`docs/research/2026-07-17-pricing-intelligence-implementation-handoff.md` §2.2/§2.3/§4/§5/§7 and
`docs/research/probes/ml_official_durable_price_probe_test.go:147-152` (buy_box_winner members
`item_id/price/currency_id/condition/listing_type_id`; both `buy_box_winner` and
`buy_box_winner_price_range` null 22/22).

---

## PART A — IC-06 remaining 4 ports (concrete Go shapes)

Package `modules/connectors/domain`. All nullable members are pointers; nil = unknown (ADR-17),
never zero. Money = `*Money{Amount string, Currency string}` (existing per-module VO,
`connectors/domain/money.go`).

### A.1 Shared `CatalogAttribute` (consumed by SearchCatalogByEAN + GetCatalogProduct)

Minimal consumed subset for the F-03 identity gate (handoff §5: gate compares GTIN, brand,
reference/model, line, variant attrs). ML `/products/{id}` and `/products/search` return
attributes as `[]{id, name, value_id, value_name}`.

```go
type CatalogAttribute struct {
    ID        string  `json:"id"`         // ML attribute id, e.g. "GTIN", "BRAND", "MODEL"
    Name      string  `json:"name"`       // human label
    ValueName *string `json:"value_name"` // normalized value; nil when ML returns null
}
```

> **[AMBIGUOUS — recommendation, proceed pending hub veto]** `value_id` OMITTED. The deterministic
> gate matches on attribute `id` + normalized `value_name` (handoff §5 rules 3–5); `value_id` is a
> provider-internal enum key that adds a member the gate does not consume. **Recommend: OMIT.**
> Add later only if F-03 needs structured value equality beyond normalized string compare.

### A.2 `SearchCatalogByEAN(ean) → CatalogSearchResult`

`FetchedAt` on the ENVELOPE (D-04) — resolves the IC-06:29 vs IC-06:23 gap (list result had no
timestamp; the "all ports return FetchedAt" rule is satisfied at envelope level).

```go
type CatalogSearchResult struct {
    Products  []CatalogCandidate `json:"products"`
    FetchedAt time.Time          `json:"fetched_at"`
}

type CatalogCandidate struct {
    CatalogProductID string             `json:"catalog_product_id"`
    Title            string             `json:"title"`
    Attrs            []CatalogAttribute `json:"attrs"`
}
```

Provider relevance order preserved; consumer does NOT assume ranking (IC-06:29).

### A.3 `GetCatalogProduct(id) → CatalogProduct`

```go
type CatalogProduct struct {
    ID           string             `json:"id"`
    Title        string             `json:"title"`
    BuyBoxWinner *BuyBoxWinner      `json:"buy_box_winner"` // nil when null (22/22 in probe)
    Attrs        []CatalogAttribute `json:"attrs"`
    FetchedAt    time.Time          `json:"fetched_at"`
}

type BuyBoxWinner struct {
    Price    *Money `json:"price"`     // decimal string; nil when null
    SellerID string `json:"seller_id"`
}
```

> **[NOTE — no change to frozen member, flag only]** IC-06:30 freezes `BuyBoxWinner{Price, SellerID}`.
> The live `buy_box_winner` payload carries `item_id` (not `seller_id`) plus `price/currency_id/
> condition/listing_type_id` (probe:150-152). Normalizing item→seller requires a second lookup or
> dropping to `item_id`. **Recommend: keep frozen `SellerID`, source it from the winning item's
> `seller_id` when the adapter already holds the item; if unavailable in the buy_box payload alone,
> populate from `ListCatalogOffers` match — implementation detail, does not change the frozen shape.**
> Since buy_box was null 22/22 in every probe, this path is currently unexercised; flagging so the
> hub is aware the frozen `SellerID` may need an item-join at wire-up time.

### A.4 `GetShipmentInfo(shipmentID) → ShipmentInfo`

```go
type ShipmentInfo struct {
    ID            string         `json:"id"`
    Status        string         `json:"status"`
    SLADue        *time.Time     `json:"sla_due"`        // nil = no SLA known
    Delayed       *bool          `json:"delayed"`        // nil = unknown, never false-as-default
    Costs         *ShipmentCosts `json:"costs"`          // nil = costs not fetched/available
    DestinationUF *string        `json:"destination_uf"` // 2-letter UF from receiver state; feeds DIFAL (IC-04)
    FetchedAt     time.Time      `json:"fetched_at"`
}

type ShipmentCosts struct {
    GrossAmount  *Money `json:"gross_amount"`  // total shipping cost, decimal string; nil = unknown
    ReceiverCost *Money `json:"receiver_cost"` // paid by buyer; nil = unknown
    SenderCost   *Money `json:"sender_cost"`   // paid by seller; nil = unknown
}
```

Minimal consumed subset: `/shipments/{id}/costs` returns `gross_amount` + `receiver`/`senders`
breakdown. DIFAL/margin consumers need the amount charged and who pays; all nullable → nil.
`DestinationUF` normalized from `receiver_address.state` to the 2-letter UF code.

### A.5 `GetFreeShippingCost(userID, item) → FreeShippingCost`

```go
type FreeShippingQuery struct {
    ItemID string `json:"item_id"` // own/catalog item the seller-paid free-shipping cost is quoted for
}

type FreeShippingCost struct {
    Cost      *Money    `json:"cost"`       // seller-paid freight for price ≥79; nil = unknown
    FetchedAt time.Time `json:"fetched_at"`
}
```

> **[AMBIGUOUS — most underspecified; recommendation, proceed pending hub veto]** IC-06:33 gives
> `GetFreeShippingCost(userID, item)` with no `item` shape. ML
> `/users/{id}/shipping_options/free` computes from the listing's own dimensions/category, which
> the item already carries server-side. **Recommend the minimal `FreeShippingQuery{ItemID}`** — the
> adapter passes the item id and ML resolves dimensions. If the live route instead requires explicit
> dimensions (weight/length/width/height) or a price threshold in the request, that is a live-lane
> finding to fold in at F-01 probe time (pre-approved live-lane-run, read-only GET). Proceeding on
> `{ItemID}` until the live probe contradicts it.

**Typed errors** (unchanged, IC-06:37): `ErrUnauthorized`, `ErrNotFound`, `ErrRateLimited`,
`ErrCatalogOffersUnavailable`, `ErrProviderUnavailable`. DTOs die at the adapter; GET-only.

---

## PART B — IC-03 public JSON (F-04 collection + verdict)

Additive OpenAPI + sdk-runtime, same commit (M-02 owns serialization; hub freezes meanings).

### B.1 `POST /market/collections` — synchronous 200 envelope

```json
{
  "status": "COMPLETED",
  "decisões": [
    {
      "codprod": "string",
      "match_status": "ACCEPT|REVIEW|REJECT|NO_CANDIDATE",
      "price_evidence_status": "OK|INSUFFICIENT_MARKET|NO_PRICE_EVIDENCE",
      "blocking_state": "NO_CANDIDATE|NO_PRICE_EVIDENCE|INSUFFICIENT_MARKET|SEM_CUSTO|null"
    }
  ],
  "contagens": {
    "ok": 0,
    "no_price_evidence": 0,
    "insufficient_market": 0,
    "no_candidate": 0,
    "sem_custo": 0
  },
  "causas": [
    { "codprod": "string", "reason": "FLAG_DISABLED|PROVIDER_4XX|PROVIDER_5XX|NO_IDENTITY|TIMEOUT", "detail": "string|null" }
  ]
}
```

- `status`: `COMPLETED` | `PARTIAL` (IC-03:35).
- `decisões`: per-codprod outcome array — one row per requested codprod, input order preserved.
- `contagens`: **enum-keyed snake_case count map** (D-04) — outcome enum → count. Keys are the
  price-evidence/blocking outcomes rolled up across `decisões`. Sums the array by category.
- `causas`: only for non-OK rows; `reason` from a closed vocabulary mirroring the IC-03 error
  matrix (FLAG_DISABLED, provider 4xx/5xx, missing identity, timeout).

> **[AMBIGUOUS — recommendation]** `contagens` key case. D-04 says "snake_case count maps"; domain
> enums are UPPER (`NO_PRICE_EVIDENCE`). **Recommend lower snake_case keys** (`no_price_evidence`)
> per D-04 wording, decoupling the wire count-map keys from the Go enum constants. Hub veto flips to
> UPPER if it prefers key==enum.

### B.2 `Verdict` (GET /market/verdicts, per codprod)

```json
{
  "match_status": "ACCEPT|REVIEW|REJECT|NO_CANDIDATE",
  "price_evidence_status": "OK|INSUFFICIENT_MARKET|NO_PRICE_EVIDENCE",
  "verdict_label": null,
  "blocking_state": "SEM_CUSTO",
  "inputs_used": {
    "our_price":    { "source": "ml_sale_price",    "as_of": "2026-07-18T12:00:00Z" },
    "target_price": { "source": "ml_price_to_win",  "as_of": "2026-07-18T12:00:00Z" },
    "market":       { "source": "ml_catalog_offers","as_of": "2026-07-18T12:00:00Z" }
  },
  "market_range": {
    "min_valid": "72.33",
    "median": "88.88",
    "currency": "BRL",
    "n_offers": 6,
    "n_sellers": 5
  }
}
```

- **`verdict_label`: ALWAYS `null` from M-02 in MIS-004 (Q3 CONFIRMED).** M-02 emits market-evidence
  verdict only = `match_status` + `blocking_state` + price targets. `verde`/`âmbar` margin labels are
  M-07-owned via IC-04 `CalcProfile` (≥18/≥10). No margin verdict crosses the M-02 boundary here.
- **`blocking_state`**: `NO_CANDIDATE|NO_PRICE_EVIDENCE|INSUFFICIENT_MARKET|SEM_CUSTO|null`.
  `SEM_CUSTO` = market evidence exists but ERP cost/CalcProfile unavailable ⇒ serve the range WITHOUT
  a label (IC-03:28). Cost comes from a consumer-side minimal cost port mirroring IC-02 GetCostAsOf,
  wired at composition root **post-merge by hub** (Q3).
- **`inputs_used`**: `{ input: { source, as_of } }` (D-04). `source` ∈ the frozen snapshot source
  enum (`ml_sale_price|ml_price_to_win|ml_catalog_offers`) plus `erp_cost` when a cost input is
  present; `as_of` = that input's `fetched_at`/cost as-of (RFC3339).
- **`market_range`**: present when price evidence exists (including SEM_CUSTO). **Mirrors the frozen
  `MarketAggregate` exactly** — `min_valid`, `median` (decimal strings, nil→`null`), `currency`,
  `n_offers`, `n_sellers`. **NO `max`** (the aggregate entity has no `max` field; D-04 "no invented
  max"). Absent value ⇒ `null`, never `0` (IC-03:42, UnknownValue).

### B.3 Aggregates finding (product- vs source-keyed) — NOT in this amendment

Per hub #4: the source-keyed aggregate variant is an **additive internal port method**, decided at
F-04-S1. It does NOT change any exposed JSON — `GET /market/aggregates` stays product-keyed
(latest-per-product), and `market_range` above mirrors that. Therefore excluded from the IC-03
draft; recorded here only for traceability.

---

## Downstream unblock (on send + hub ACK)

- **F-01** remaining 4 ports (SearchCatalogByEAN, GetCatalogProduct, GetShipmentInfo,
  GetFreeShippingCost) dispatch on A.1–A.5. Live-lane-run REQUEST fires when ports are probe-ready
  (read-only GET, dispatcher OFF).
- **F-03** identity resolver (Q2 already RESOLVED) consumes A.1/A.3 `CatalogAttribute` for the gate.
- **F-04** collection + verdict consumes B.1/B.2; carries the n_sellers<5⇒INSUFFICIENT_MARKET
  invariant (F-02-S1 F3 carry-forward) and the F-04-S1 source-keyed-aggregate decision.
