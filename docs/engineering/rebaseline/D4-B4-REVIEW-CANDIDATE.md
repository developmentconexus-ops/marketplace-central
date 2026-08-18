# D4-B4 — Market / Economics / Settlement Contract — REVIEW CANDIDATE

> **Status:** REVIEW CANDIDATE / NON-AUTHORITATIVE / DISPOSABLE — READY FOR OPERATOR RATIFICATION  
> **Program:** Architecture Rebaseline / Technical System Design  
> **Parent authority:** accepted D0–D4-B3 only  
> **Method:** DevelopmentConexus Engineering Method v1.0.0  
> **Initial B4 candidate:** `9c212afc58bbce037cc85582329a3652fd2d665d`  
> **Independent challenge:** `94e49fdf846814b9ff2d488528d1f2c1251aa737` (`AI-DIALOG.md`, non-authoritative)  
> **Post-review consolidation:** `c7e4d770e578f7f5dfaa0db62d88ea69bc35a6a9`  
> **Live M1/E1/S1 evidence:** `d9c4fa558e709e057744f66b133455d84ccdb080` (`AI-DIALOG.md`, non-authoritative)  
> **Purpose:** final coherent B4 candidate after independent challenge, real Installation evidence and operator-approved provider-richness amendment. This file is not target authority and MUST be deleted during canonical consolidation after ratification.

---

## 1. Review question and scope

D4-B4 must answer:

> **Which concrete external market, expected-cost, order-economic, billed-charge, payment/release and financial-movement evidence contracts are required for Product 1.0, and how can MPC exploit materially useful provider capabilities without either collapsing to the lowest common denominator or turning provider vocabulary into MPC business ontology?**

B4 remains an **external-evidence contract**, not a finance domain, marketplace ledger, pricing authority, scraping subsystem or generic provider framework.

The target remains three coherent evidence families:

1. **Market Evidence** — externally authoritative competitive observations consumed by Market Intelligence and, through its interpretation, Commercial Economics.
2. **Expected / Order Economic Evidence** — provider fee/shipping/order evidence required for L0 Expected Economics and L1 Order Economics.
3. **Financial Movement / Release / Billed-Charge Evidence** — provider/payment-native evidence required for Economic Attribution, L2/R2 and applicable Post-Sale refund-consequence closure.

Implementation remains blocked until D9.

---

## 2. Imported authority — do not re-decide

B4 imports the following accepted meaning:

1. **Market Intelligence** owns external comparability, competitor relevance, competitive position/change and market-evidence sufficiency.
2. **Commercial Economics** owns Cost Basis, pricing analysis, L0 Expected Economics, L1 Order Economics, L2 Realized/Settlement Economics, Economic Attribution and Economic Reconciliation.
3. **Post-Sale Resolution** owns coordination/correlation/closure of material cancellation/return/refund consequences.
4. **Marketplace Offering Operations** owns Listing/Price Intent and marketplace price mutation/convergence; Market Evidence never becomes Price Intent by convenience.
5. D4 owns concrete external acquisition/protocol/capability/evidence only.
6. Provider-native Payment, Refund, Fee, Adjustment, Release, Withdrawal/Payout and equivalent resources/occurrences remain external/source-qualified evidence.
7. R1 Expected↔Order and R2 Order↔authoritative realized evidence remain Commercial Economics semantics.
8. R3 provider payout/withdrawal↔Bank Cash Receipt remains unclaimed until an accepted bank source exists.
9. Marketplace Order, Shipment, catalog competition, Payment, Billing and account/report resources remain distinct provider evidence surfaces.
10. Known value, known empty/not-applicable, unknown, unavailable, unsupported and partial remain distinct.
11. Money/cost/pricing values preserve exactness and D2 provenance sufficient for the claim.
12. Consumer owns meaning; adapter owns provider protocol.
13. Provider PII is minimized.
14. Current code/modules/tables/legacy ADR shapes are evidence only.
15. No B4 external write is introduced by symmetry.

No D0/D1/D2/D3 or D4-B1/B2/B3 reopen is required by the evidence obtained in B4.

---

## 3. Global Maximum

### 3.1 Rejected — lowest-common-denominator marketplace contract

Do not reduce every provider to only the fields/capabilities shared by every future marketplace.

Example of the rejected shape:

```text
Mercado Livre can expose:
  price
  buyer shipping
  free-shipping state
  catalog competition
  price_to_win
  winner evidence
  boosts

MadeiraMadeira may expose only:
  price
  availability

therefore MPC keeps only:
  price
  availability
```

**REJECT.** This destroys useful evidence already available from a supported provider and materially weakens Product 1.0 Competitive Intelligence.

### 3.2 Rejected — provider mirror / universal provider ontology

Do not solve provider richness by copying every provider field into MPC entities or by building a generic Provider/Resource/Capability/FinancialTransaction graph.

**REJECT.** That promotes volatile provider vocabulary into business architecture, creates accidental complexity and produces a second authority beside D1 domains.

### 3.3 Selected — Semantic Core + Provider-Enriched Evidence

**CURRENT STRUCTURE CONFIRMED / GLOBAL MAXIMUM.**

```text
provider-specific surfaces
        ↓
D4 concrete adapter
  - namespace
  - capability/coverage
  - source meaning
  - provider-native evidence
  - time/provenance/granularity
        ↓
semantic core where meanings genuinely align
        +
provider-enriched evidence where the provider exposes more
        ↓
D1 consumer
  Market Intelligence
  Commercial Economics
  Post-Sale Resolution
  etc.
```

Shared semantics are normalized **only where the meanings actually align**. Material provider-distinct evidence remains source-qualified and optional.

A provider that lacks an equivalent capability does not force another provider to discard it.

---

## 4. Governing invariants

### 4.1 External evidence invariant

> **Every external market/economic value entering MPC is source-qualified, scope-qualified, granularity-qualified and time/provenance-qualified strongly enough for the exact claim being made. D4 preserves only what the source proves; the D1 consumer retains interpretation and business authority.**

### 4.2 Provider Richness Invariant

> **MPC MUST NOT discard materially useful provider evidence merely because another marketplace lacks an equivalent capability. A supported provider's richer evidence may be acquired and consumed when it serves a named Product 1.0 consumer or correctness property. Provider-specific evidence remains source-qualified/optional and MUST NOT become universal MPC business ontology merely because one provider exposes it.**

Consequences:

- Mercado Livre `price_to_win`, competition boosts, catalog offer shipping evidence and provider competition reasons may be retained/used because they materially improve Market Intelligence.
- A future MadeiraMadeira adapter is not required to fabricate equivalents.
- Absence of a capability is represented honestly as unsupported/not-applicable/unavailable/unknown as appropriate, never zero/default.
- A richer provider surface may support richer UI/analysis later without changing which D1 domain owns the meaning.

### 4.3 Capability-rich does not mean payload-rich

> **Maximize materially useful provider capability, not raw payload retention.**

A provider field/capability enters the target only when at least one is true:

1. a named Product 1.0 consumer requires it;
2. it protects a known correctness property or explains a material outcome;
3. it is materially non-reobservable evidence required by an accepted claim.

This does not authorize indiscriminate mirroring of provider payloads, PII, deprecated fields, UI copy, debug fields or speculative future data.

### 4.4 Monetary/economic rung separation

The following never collapse merely because amounts happen to match:

```text
expected provider fee
≠ expected seller shipping
≠ Order transaction fee
≠ realized seller Shipment cost
≠ billed charge/rebate
≠ Payment approval
≠ money release/account impact
≠ refund/reversal
≠ withdrawal/payout
≠ Bank Cash Receipt
```

`price_to_win` is provider market evidence, not MPC Price Intent or automatic price recommendation.

---

# 5. B4-A — Market Evidence

## 5.1 Provider-independent contract

D4 supplies market observations only through admitted external sources and preserves proportionately:

- Organization + Marketplace Installation/source;
- provider listing/catalog/product scope used by the operation;
- observed offer/price/shipping/competitive-status evidence;
- currency/source dimensions required by the observation;
- provider occurrence/update time where exposed;
- acquisition time;
- operation coverage/pagination/completeness state;
- provider identifiers as external references, not MPC identities;
- provider-enriched evidence whose material consumer is known.

Market Intelligence owns:

- comparability;
- competitor-set relevance;
- competitive position/change;
- whether evidence is sufficient for a market conclusion;
- any derived delivered-price/competitive explanation.

D4 never converts raw competition evidence directly into a Price Intent.

## 5.2 Mercado Livre selected realization

Current real Installation evidence proved useful market surfaces including:

1. `price_to_win` / catalog competition;
2. catalog product offer population;
3. provider offer shipping/free-shipping evidence;
4. provider competition boosts/reasons/status where exposed.

These are admitted as **provider-enriched Market Evidence**, not universal marketplace fields.

## 5.3 Why the enriched evidence is material — measured case

Real Installation measurement:

```text
our offer
  price               69.90
  buyer shipping      44.94
  free_shipping       false

winner
  price               79.90
  buyer shipping       0.00
  free_shipping       true
  shipping tag        mandatory_free_shipping

provider price_to_win 26.75
```

The organization's product price was lower, yet the buyer-facing delivered amount was materially worse:

```text
our apparent delivered amount     = 69.90 + 44.94 = 114.84
winner apparent delivered amount  = 79.90 + 0.00  = 79.90
```

The catalog population also exposed shipping/boost differences. Therefore a price-only normalization would hide the evidence needed to explain the competitive result.

Target consequence:

- D4 may preserve own/winner price, buyer shipping, free-shipping/provider shipping tags, competition status, `price_to_win`, boosts and relevant provider correlation.
- Market Intelligence may derive a delivered-price/competitive explanation from those facts.
- Commercial Economics may later consume the interpretation plus seller-borne cost evidence to evaluate trade-offs.
- Offering remains the only owner of Price Intent/action.

## 5.4 Market knowledge states measured

The bound seller population proved multiple distinct outcomes:

- catalog competition observed;
- `not_listed` / `item_not_opted_in` returned as HTTP 200 with null price fields;
- a `catalog_listing=true` Item can still return competition `not_listed` with `catalog_product_id=null`;
- competitor-owned point access can return 403 unavailable/forbidden;
- catalog offer population has bounded provider paging and is not general-market completeness.

Therefore:

1. HTTP 200 does not imply a positive market observation.
2. Catalog membership does not imply active catalog competition evidence.
3. The competition payload does not become catalog-membership authority.
4. Provider catalog population is not the total market.

## 5.5 M1 verdict

**M1 = PASS / CLOSED AS INSTALLATION LANE-SELECTION EVIDENCE.**

M1 proved a usable real competition lane plus negative/not-applicable/unavailable/coverage-bounded controls. It is not a standing B4 closure gate.

---

# 6. B4-B — Expected / Order Economic Evidence

## 6.1 L0 Expected Economics

Commercial Economics owns L0. D4 supplies provider-dependent inputs only:

```text
candidate sale context
  ├─ Sankhya Expected Tax            ← accepted B3
  ├─ marketplace expected sale fee
  ├─ marketplace expected seller shipping
  └─ provider promotion/discount evidence when materially applicable
              ↓
Commercial Economics
              ↓
L0 Expected Economics
```

Expected selling fee and expected seller shipping stay distinct even when provider rules interact.

## 6.2 Mercado Livre expected sale-fee surface

Selected surface: current `listing_prices` operation.

Live falsification proved that, for the measured MLB context, the meaningful inputs/effects include:

- price;
- listing type;
- category;
- **shipping mode** for logistics-sensitive fixed fee;
- provider-returned currency;
- provider-returned fee components/rate/fixed fee.

The same live probe proved that the measured `listing_prices` operation silently ignored:

- `quantity`;
- `logistic_type` including an intentionally invalid value;
- `billable_weight`;
- an arbitrary unknown parameter.

Those facts are **surface-specific**, not universal provider laws. Weight/dimensions remain material to the separate seller-shipping surface below.

### Adapter obligations from fail-open behavior

1. HTTP 200 never proves a request qualifier was consumed.
2. Requested listing type must be validated against the returned selected listing type/response shape; do not accept an all-types response and take `response[0]`.
3. Category must be explicitly qualified when the claim requires category-specific pricing; omission must not silently become a site-default quote.
4. Provider-returned currency is preserved; do not fabricate `BRL` when the provider already returned the authoritative currency field.
5. The adapter preserves returned component scope and does not invent arithmetic identities.
6. Material request semantics are proven by decorrelation/falsification when the provider can silently ignore fields.

Current `FeeQuote` code remains current-state evidence and is known insufficient for the target because it does not carry the complete proven qualification/fail-open fences.

## 6.3 Expected seller shipping

Selected first-lane surface: Mercado Livre seller-shipping estimation (`/users/{USER_ID}/shipping_options/free` family).

Live evidence proved this is distinct from `listing_prices` and is sensitive to item/dimensions/weight/price context.

D4 preserves **seller-borne expected shipping** for the requested context. It never substitutes buyer-paid freight or a generic `shipping` number.

The same anti-silent-success rule applies to this surface: a 200 response does not prove every submitted input was honored; the adapter must validate the operation's real contract proportionately.

## 6.4 L1 Order Economics

Once a Sale exists, L1 uses transaction-specific evidence rather than re-running L0 and calling it actual.

Current first-flow evidence includes:

- Order price/quantity/discount context;
- provider Order transaction `sale_fee`;
- related Payment reference;
- related Shipment seller-side cost;
- attributable B3 business-system/fiscal evidence.

### Source-specific granularity/decomposition

The live Installation re-proved `order_items[].sale_fee` as **per unit**, including a multi-unit Order.

Payment evidence then showed the same Order-level sale fee represented as multiple charge rows with different destinations (`mp` / `ml`) plus separate shipping charge evidence.

Therefore:

> **Granularity/decomposition is preserved per source. Order fee evidence is not assumed to be the Payment fee evidence re-expressed at another scale.**

Commercial Economics owns aggregation/classification/attribution.

Buyer shipping charge, seller shipping cost and provider fee remain distinct.

## 6.5 Billing / charged-fee evidence

Mercado Livre Billing remains **billed-charge / rebate / bonus / fiscal-reconciliation evidence**.

It may explain divergence such as:

```text
Order transaction fee
≠ billed charge after rebate/bonus/adjustment
```

Billing is not Payment release/cash authority and is not required merely to prove a Payment release.

## 6.6 ADR-009

ADR-009's durable fee/value provenance rule is already homed in D2. B4 cites that authority and adds only source-specific evidence requirements.

B4 does not inherit:

- a universal `channel_fees` table;
- legacy fee layers as economic ontology;
- one global resolution ladder;
- config fallback as provider evidence;
- a universal fee ledger.

## 6.7 E1 verdict

**E1 = PASS / B4 CLOSURE GATE CLOSED.**

Live proof established materially sufficient sanctioned evidence for the selected first-flow L0/L1 marketplace components while falsifying silent-success cases rather than trusting one plausible quote.

The selected lane preserves:

```text
expected selling fee
≠ expected seller shipping
≠ Order transaction fee
≠ realized seller Shipment cost
```

No E1 `STOP / SPLIT PREREQUISITE` remains.

---

# 7. B4-C — Financial Movement / Release / Billed-Charge Evidence

## 7.1 Provider-independent contract

D4 exposes source-qualified external financial evidence without creating synthetic MPC Payment/Refund/Settlement business entities.

A material external financial occurrence preserves proportionately:

- Organization + external account/source namespace;
- provider-native occurrence/resource identity;
- provider-native kind/status;
- currency;
- gross/net/component values supplied by the source;
- source occurrence/approval/release/reversal times separately;
- acquisition time;
- native Order/Payment/Shipment/external-reference anchors where exposed;
- component direction/source/destination where needed for correct interpretation;
- refunds/chargebacks/adjustments/withdrawals as distinct occurrences when the source models them distinctly;
- population/report coverage only when a population/report surface is actually used.

D4 does not decide realized profit, Economic Attribution closure or R2 closure.

Refund/chargeback/adjustment evidence may also feed **Post-Sale Resolution** when its explicit consequence scope needs financial closure evidence. The same external evidence can feed multiple legitimate consumers without transferring their authority.

## 7.2 Same Marketplace Installation credential can reach Payment

For the selected bound Mercado Livre Installation, live evidence proved:

```text
Mercado Livre Order
      ↓ payments[].id
GET https://api.mercadopago.com/v1/payments/{id}
Authorization: Bearer <same bound ML Installation access token>
      ↓ HTTP 200
```

A nonexistent Payment control returned honest 404.

Therefore:

> **A separate Mercado Pago application/credential is NOT required for the selected Product 1.0 per-sale Payment path on this Installation.**

This is an Installation/capability fact, not an eternal provider law. A future authorization/scope change is a reopen/revalidation trigger, not grounds to build a second credential path today.

## 7.3 Approval, release and cash remain distinct

Live Payments proved:

- `approved` can precede `released` by days;
- `money_release_date` may be populated while `money_release_status` is still `pending`;
- therefore timestamp presence does not prove release;
- withdrawal/payout and bank receipt are absent from the selected Payment surface and remain outside S1-A/R3.

D4 preserves provider status and release state/time separately.

## 7.4 Payment component evidence

The selected Payment surface exposed materially useful evidence through `charges_details` including:

- native charge identity/type;
- original/refunded amounts;
- `accounts.from` / `accounts.to` direction;
- Shipment correlation metadata where present;
- Order/external-reference correlation.

Live measurement also proved:

- `fee_details` can be incomplete or empty while real charges exist;
- `net_received_amount` is not refund-adjusted and therefore is not the final realized-economic state after later reversal;
- component direction matters; a set of unsigned amounts is not safely interpretable as one seller cost.

Target obligation:

> **For the selected per-sale L2/R2 evidence path, use the provider's directionally qualified charge/reversal evidence needed to support the claim; do not promote `fee_details`, `net_received_amount` or `money_release_date` alone into realized authority.**

Commercial Economics still owns the economic interpretation.

## 7.5 Refund after release — measured consequence model

A real existing refund case proved:

```text
earlier release fact remains historically true
        +
later refund/reversal evidence is appended
```

The provider did not rewrite the earlier release into non-existence.

The same refund occurrence legitimately feeds:

- **Commercial Economics** for attribution/R2 consequences;
- **Post-Sale Resolution** for refund-consequence closure.

Neither consumer acquires the other's authority.

## 7.6 S1-B account-movement universe

S1-A covered every material movement class reached in the selected real samples: commission components, financing-related charges, seller shipping, release and full refund after release.

No concrete Product 1.0 correctness gap demonstrated a need for a broader account/report population during B4.

Therefore:

**S1-B = DEFER SAFELY / NOT CURRENTLY REQUIRED.**

Residual bounded Unknown/reopen examples:

- material unanchored account adjustment;
- late chargeback fee not reachable from the per-sale Payment relation;
- period-completeness claim that cannot be established from anchored point reads;
- materially required withdrawal/payout population.

If such a real consumer/failure appears, open the smallest necessary read-only population source then. Do not build/report-sync by symmetry.

## 7.7 S1 verdict

**S1 = PASS / B4 CLOSURE GATE CLOSED.**

The bound Installation exposes sanctioned per-sale Payment/release/refund evidence sufficient for the selected Product 1.0 L2/R2 claim without a separate Mercado Pago credential and without report generation.

---

## 8. No report-generation/write surface by convenience

B4's current target is read/evidence-only for Market/Economics/Settlement acquisition.

A provider POST that creates a report artifact is an external effect, not "read support".

Report generation is not admitted by B4 merely because a later GET could consume it.

If a future required evidence class can only be obtained by generating a report, return to explicit D3/D4 external-effect adjudication and operator authorization before adding that effect.

---

## 9. Source admissibility and legacy ADR disposition

### ADR-014 — on-demand/local runtime

**Candidate disposition after B4 ratification: HISTORICAL.**

Its old on-demand/local-Docker target shape has no surviving D4 authority. Honest evidence/absence is already carried by D0–D4; runtime/cadence belongs D7.

### ADR-020 — generic `CollectorPort`

**Candidate disposition: target shape SUPERSEDED.**

One real market source does not justify a generic collector framework.

The surviving platform-level source-admissibility rule should be rehomed into `ARCHITECTURE.md` during canonical B4 consolidation:

> **Absence of an admitted external market-data source never authorizes fabricated evidence or an unadjudicated scraping path by convenience. A materially new market-data source requires explicit source, legality/trust, coverage and provenance adjudication before its evidence can support MPC claims.**

This is not a permanent ban on every future lawful collector. It prevents an unreviewed source from becoming truth because another provider lacks data.

### ADR-032 — catalog-offers default-off flag

**Candidate disposition: target meaning SUPERSEDED / historical after B4.**

A runtime flag is not provider capability authority. D7 later decides any necessary toggle/runtime mechanism.

### ADR-009 — provenance

Remains a **carried constraint with active home in D2**. B4 does not create a second home.

---

## 10. YAGNI / overengineering fence

B4 MUST NOT introduce:

- lowest-common-denominator suppression of useful provider evidence;
- indiscriminate provider payload mirroring;
- universal Provider/Capability/MarketObservation graph;
- generic Financial Transaction / Payment / Refund / Settlement MPC business entity merely for normalization;
- universal D4 financial ledger;
- universal fee ledger / `channel_fees` target table;
- generic `CollectorPort` plugin/framework;
- unadjudicated scraping infrastructure;
- provider PII/raw payload retention by convenience;
- one global `Fee` shape that collapses expected/order/billed/payment decompositions;
- one generic correlation key inferred from amount/time;
- `price_to_win` as automatic recommended price;
- Billing as release/cash authority;
- Payment approval as release;
- `money_release_date` presence as release proof;
- `net_received_amount` as post-refund realized authority;
- `fee_details` as complete fee authority;
- provider release as bank receipt;
- ERP receivable/baixa as marketplace/payment realized authority;
- margin/profit calculation inside D4;
- report generation by symmetry;
- support for every possible provider field/report/movement class before a named consumer/correctness need exists.

---

## 11. Proof status

| Gate / material claim | Final B4 evidence status |
|---|---|
| M1 Market Evidence lane | **PASS / CLOSED** |
| catalog positive + not-applicable + unavailable distinctions | **PROVEN** |
| provider-rich market evidence materially useful | **PROVEN** by price/shipping/winner case |
| E1 expected selling fee | **PASS / CLOSED** |
| silent request-field/fail-open risk | **PROVEN** |
| E1 expected seller shipping | **PASS / CLOSED** |
| L1 `sale_fee` granularity | **PROVEN per-unit on real multi-unit Order** |
| source-specific Order vs Payment fee decomposition | **PROVEN** |
| seller vs buyer shipping scope distinction | **PROVEN** |
| same ML Installation token → Payment API | **PROVEN / HTTP 200** |
| separate Mercado Pago credential needed for selected lane | **PROVEN UNNECESSARY** |
| approval ≠ release | **PROVEN** |
| `money_release_date` presence ≠ release | **PROVEN** |
| `fee_details` incomplete | **PROVEN** |
| `net_received_amount` not refund-adjusted | **PROVEN** |
| refund appended after prior release | **PROVEN** |
| Post-Sale + Economics dual consumption without authority transfer | **PROVEN** |
| S1 per-sale L2/R2 evidence | **PASS / CLOSED** |
| broader account movement universe | **DEFER SAFELY / bounded Unknown** |
| R3 bank side | **DEFERRED / unclaimed** |
| report generation | **NOT ADMITTED / not required** |

No known B4 live-evidence closure gate remains.

---

## 12. Reopen / stop triggers

Reopen only the implicated decision when material evidence shows:

1. a required provider-rich capability cannot fit source-qualified evidence without leaking provider ontology into a D1 domain;
2. a supported provider lacks a legitimate source for a Product 1.0 claim and honest insufficiency is no longer acceptable to the product requirement;
3. a new external market-data source requires scraping/vendor/manual ingestion whose legality/trust/coverage materially differs from current admitted sources;
4. provider expected-cost/shipping semantics change such that the current falsified request contract no longer proves the selected claim;
5. the bound ML token can no longer access required Payment evidence and no sanctioned equivalent exists;
6. a material movement appears without a usable per-sale anchor and Product 1.0 correctness requires a population/recovery source;
7. provider monetary evidence cannot fit D2 source-qualified identity/provenance semantics;
8. financial occurrence/duplicate/partial/recovery semantics cannot fit D3;
9. Product 1.0 genuinely requires report generation or another B4 external write;
10. a bank source becomes accepted and R3 bank-side reconciliation becomes real scope;
11. a second provider proves repeated technical mechanics whose duplication is materially worse than a small shared non-authority mechanism.

Naming preference, desire for provider symmetry, current module convenience and hypothetical provider futures are not reopen evidence.

---

## 13. Final candidate decision

### `CURRENT STRUCTURE CONFIRMED`

The accepted D1/D2 boundary survives:

- D4 acquires/translates concrete external evidence;
- Market Intelligence owns comparability and competitive interpretation;
- Commercial Economics owns L0/L1/L2, attribution and reconciliation;
- Post-Sale Resolution owns refund-consequence closure;
- Offering owns Price Intent/action;
- provider-native vocabulary remains outside MPC business ontology.

### `PROVIDER-RICH / SEMANTICS-FIRST`

B4 additionally freezes that MPC does **not** collapse providers to the lowest common denominator. Material provider-specific evidence is preserved when it serves a named Product 1.0 consumer/correctness property, while unsupported equivalents on another provider remain honestly absent/unsupported rather than fabricated.

### Mercado Livre selected target evidence

The sanctioned Mercado Livre/Mercado Pago surfaces are materially sufficient for the currently claimed B4 Product 1.0 evidence lane:

- catalog competition / offer population provides enriched Market Evidence;
- expected selling fee and expected seller shipping are separately observable and falsifiable;
- Order/Shipment provide L1 transaction evidence with source-specific granularity;
- the same bound ML Installation credential can read the selected Payment resource;
- Payment charge/release/refund evidence is sufficient for selected per-sale L2/R2;
- no separate Mercado Pago credential, account report or generated report is required for this first lane.

### B4 status proposed for operator ratification

```text
Architecture core               CONFIRMED
Provider Richness Invariant     APPROVED BY OPERATOR / candidate only until B4 ratification
M1                              PASS
E1                              PASS
S1                              PASS
B4 closure gate                 NONE KNOWN
D0–D4-B3 reopen                 NONE
B4 canonical status             NOT YET — requires operator ratification of B4 as a whole
Implementation                  BLOCKED until D9
```

> **D4-B4 is READY FOR OPERATOR RATIFICATION / CANONICAL CONSOLIDATION. No known B4 closure gate remains.**

---

## 14. Canonical-consolidation plan after ratification

Only after explicit operator ratification of B4 as a whole:

1. consolidate B4 into `D4-EXTERNAL-INTEGRATIONS.md`;
2. update the router: B4 `ACCEPTED / CANONICAL`, D4 Final Global Coherence `NEXT`;
3. minimally update `ARCHITECTURE.md` with the stable Provider Richness + source-admissibility principles;
4. update ADR registry dispositions for ADR-014/020/032 while keeping ADR-009 homed in D2;
5. delete this disposable candidate;
6. run D4 Final Global Coherence + YAGNI / Overengineering / Future-Cost review;
7. do not begin D5 until D4 final coherence closes;
8. do not implement product code; implementation remains blocked until D9.
