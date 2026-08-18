# D4-B4 — Market / Economics / Settlement Contract — REVIEW CANDIDATE

> **Status:** REVIEW CANDIDATE / NON-AUTHORITATIVE / DISPOSABLE  
> **Program:** Architecture Rebaseline / Technical System Design  
> **Parent authority:** accepted D0–D4-B3 only  
> **Method:** DevelopmentConexus Engineering Method v1.0.0  
> **Base authority HEAD:** `e8a0b60f4478f2e41827ae3174f39e5eca61dfcb`  
> **Opened:** 2026-08-18  
> **Purpose:** coherent disposable surface for independent challenge and operator adjudication of D4-B4. This file is not target authority and MUST be deleted before canonical consolidation.

---

## 1. Review question and scope

D4-B4 must answer:

> **Which concrete external market, fee, financial-movement and settlement contracts are required for Product 1.0, and how can D4 provide source-qualified evidence with honest coverage/provenance while Market Intelligence and Commercial Economics retain semantic authority?**

B4 must not create a finance domain, marketplace ledger, universal payment model, scraping subsystem or economic-calculation authority inside D4.

The target is deliberately divided into only three evidence families:

1. **Market Evidence** — externally authoritative competitive observations required by Market Intelligence / Commercial Economics.
2. **Expected / Order Economic Evidence** — provider fee, shipping and transaction-economic evidence required for L0 Expected Economics and L1 Order Economics.
3. **Financial Movement / Settlement Evidence** — provider/payment-native movements required for L2 Realized Economics and R2 reconciliation.

Implementation remains blocked until D9.

---

## 2. Imported authority — do not re-decide

B4 imports:

1. **Market Intelligence** owns external comparability, competitive position/change and market-evidence sufficiency.
2. **Commercial Economics** owns Cost Basis, economic interpretation, pricing analysis, L0 Expected Economics, L1 Order Economics, L2 Settlement/Realized Economics, Economic Attribution and Economic Reconciliation.
3. D4 owns concrete external acquisition/protocol/capability/evidence only.
4. Provider-native Payment, Refund, Fee, Adjustment, Settlement, Payout and equivalent movements remain external/source-qualified identities.
5. Economic Attribution is MPC-owned Commercial Economics state; D4 does not decide what a provider movement means economically.
6. R1 Expected↔Order and R2 Order↔authoritative settlement/realized evidence are Commercial Economics reconciliation semantics.
7. R3 payout/settlement↔Bank Cash Receipt exists only when an accepted bank source exists; B4 does not invent a bank source.
8. Marketplace Order, Shipment and provider financial movements remain distinct source resources.
9. Known / absent / unknown / unavailable / partial remain distinct; provider unavailability or partial coverage never becomes zero or complete.
10. Money/tax/cost/pricing values preserve exactness and provenance sufficient for the claim.
11. Consumer owns meaning; adapter owns provider protocol.
12. Provider PII is minimized.
13. Current code/modules/tables/ADRs are evidence only.
14. No B4 write is introduced merely for symmetry. If Product 1.0 does not require a market/economic/settlement write, B4 remains read/evidence-only.

No D0/D1/D2/D3 or D4-B1/B2/B3 reopen is proposed by this candidate.

---

## 3. Evidence classification

### Known

- Current official Mercado Livre competition surfaces include item-level catalog competition / `price_to_win` evidence.
- Current official Mercado Livre catalog-product offer surfaces expose competing offers where the catalog topology/applicability supports them.
- Current official Mercado Livre `listing_prices` provides sale-cost quotation evidence and, for Brazil under the current cost structure, fixed-fee correctness depends on logistics context; omitting required logistics qualifiers can produce a fee different from the fee actually charged.
- Current repository `FeeQuote` evidence calls `listing_prices` with price/listing type/category only; therefore that current contract is not sufficient target evidence by inheritance.
- Expected seller shipping cost and realized shipment cost are distinct evidence from sale commission/fee evidence.
- Current Mercado Livre Order evidence can expose order/item fees and Payment references, but Order/embedded Payment evidence is not settlement/cash authority by itself.
- Mercado Livre Billing reports expose billing/charge details and are explicitly period/report-oriented rather than the authoritative real-time operational Order surface.
- Mercado Pago account/release reports expose financial transaction/movement evidence, including settlement/refund/chargeback/dispute/withdrawal/payout families and correlation fields such as source/order/shipping references where available.
- ADR-014, ADR-020 and ADR-032 are explicitly reopened for D4-B4; their runtime/module/flag shapes have no target authority by inheritance.
- ADR-009's durable provenance principle is already rehomed in current D2/evidence semantics; its legacy `channel_fees` table/layer ladder is not target authority.

### Inferred

- For Product 1.0, the smallest coherent external-economic boundary is operation-specific evidence contracts rather than one generic financial ledger/resource graph.
- `price_to_win` / catalog-offer observations are provider market evidence, not an MPC recommended price or comparability conclusion.
- Billing and account-settlement reports may overlap in monetary content but answer different evidence questions; they should not be normalized into one authority before real-account evidence proves the precise overlap/correlation.

### Unknown

- Exact current access/scope of the bound Metal Nobre Mercado Livre/Mercado Pago credentials to the required Billing and Mercado Pago report surfaces.
- Exact real-account pagination/retention/completeness behavior needed for the Product 1.0 settlement universe beyond current official documentation.
- Exact strongest source key(s) for deterministic Order/Payment/Shipment ↔ settlement-movement attribution on the bound account in every material movement class.
- Whether every selected Product 1.0 listing has applicable catalog competition evidence; absence of catalog competition is not equivalent to no competitors in a broader market sense.
- Any market evidence outside sanctioned provider/vendor/manual sources. No scraping path is admitted by this candidate.

### Deferred

- polling/report-generation cadence, scheduled collection, cache topology, retry/backoff, statement/report generation scheduling and storage mechanics — D7;
- UI representation and operator workflows — D6;
- end-to-end economic proof across selected real Sale/Shipment/Settlement cases — D8;
- bank-cash R3 proof until a bank source is explicitly accepted.

---

## 4. Root cause

The legacy/current implementation evidence mixes several distinct questions under convenient modules/tables:

- market observations vs economic interpretation;
- expected fee quote vs realized fee;
- order/payment evidence vs account settlement;
- provider billing documents vs actual money movement;
- value provenance vs one shared `channel_fees` persistence model;
- provider access switches vs provider-effective capability.

That structure makes silent authority collapse reachable. A number can look economically authoritative merely because it came from a provider API or shared table even when its scope, lifecycle or qualification is different from the question being answered.

A concrete current example is the legacy `FeeQuote`: it can call a valid provider endpoint successfully while omitting logistics qualifiers that are now materially required for Brazil fixed-fee correctness. HTTP success therefore does not prove economic-input sufficiency.

---

## 5. Target invariant

> **Every external economic or market value entering MPC is source-qualified, scope-qualified and time/provenance-qualified strongly enough for the claim being made; D4 preserves what the external source actually proves, while Market Intelligence owns comparability and Commercial Economics owns interpretation, attribution and reconciliation. Expected, order-time, billed and settled evidence never collapse merely because they contain similar monetary fields.**

Corollaries:

- external evidence is not economic meaning;
- provider fee quote is not realized fee;
- Order fee is not settlement;
- Billing charge is not cash receipt;
- Payment is not Payout;
- Payout is not Bank Cash Receipt;
- `price_to_win` is not MPC Price Intent;
- lack of market evidence is not evidence that no competitor exists;
- a current provider flag is not capability authority;
- one persistence table is not required merely because multiple evidence classes share money/provenance fields.

---

## 6. Credible alternatives / Global Maximum

### A — Preserve legacy `market` / `channel_fees` / profitability structure as target

**REJECT.** It would promote current module/table/layer boundaries into target authority and preserve the root cause: evidence lifecycle/meaning conflated by storage/module convenience.

### B — Create a generic Financial Movement / Economic Evidence ledger owned by D4

Examples: universal `Payment`, `Settlement`, `Fee`, `MarketObservation`, `FinancialTransaction` or provider-resource graph normalized across all future providers.

**REJECT.** External identities already remain source-qualified under D2 and Commercial Economics owns attribution/reconciliation. One generic D4 financial ledger would create duplicate business authority and speculative abstraction.

### C — Operation-specific external evidence contracts + domain-owned interpretation

**PROPOSED GLOBAL MAXIMUM.**

```text
external provider/payment system
        ↓
D4 concrete adapter
  - source identity
  - operation-specific coverage
  - provider-native evidence/components
  - source/occurrence/acquisition time
  - correlation references
        ↓
Market Intelligence / Commercial Economics
  - comparability
  - L0/L1/L2 interpretation
  - Economic Attribution
  - R1/R2 reconciliation
```

Prepare the seam for later providers; do not build their ontology now.

---

# 7. B4-A — Market Evidence

## 7.1 Provider-independent contract

D4 supplies market observations only when a legitimate external source can establish the observation strongly enough for the consuming claim.

A market observation preserves, proportionately:

- Organization + Marketplace Installation/provider source;
- external product/listing/catalog scope used by the source;
- observed offer/price/competitive-status evidence;
- currency and other source dimensions required for comparison;
- source occurrence/update time where exposed;
- acquisition time;
- provider coverage/pagination/completeness semantics;
- raw provider identifiers as external references, not MPC identities.

Market Intelligence owns:

- whether two observations are comparable;
- which competitor set is relevant;
- competitive position/change;
- whether evidence is sufficient for a market conclusion.

Commercial Economics may consume Market Intelligence meaning; it does not independently reinterpret raw competitor payload when Market Intelligence owns comparability.

## 7.2 Current Mercado Livre realization candidate

Current official provider evidence includes:

1. **`GET /items/{ITEM_ID}/price_to_win?...&version=v2`** for catalog-competition position/price-to-win/boost evidence when applicable.
2. **Catalog product offer population** through the current catalog product/items surfaces when a product is catalog-applicable.

Target interpretation:

- `price_to_win` is provider evidence about catalog competition, not MPC recommended price;
- winning/competing/listed/sharing/provider boost vocabulary remains adapter-local evidence translated only as needed;
- catalog-offer seller/price/logistics evidence does not by itself establish Market Intelligence comparability beyond the provider catalog relation;
- a non-catalog listing or unavailable competition surface remains explicit insufficient/unavailable evidence, not “no competitors”.

## 7.3 ADR-014 / ADR-020 / ADR-032 disposition candidate

### ADR-014 — on-demand/local Docker

**SUPERSEDE target-shape portion.** Preserve only the durable honesty rule: no historical market claim before evidence exists. Collection cadence/runtime belongs D7 and is not frozen as on-demand or local Docker.

### ADR-020 — CollectorPort/no scraping

**SUPERSEDE generic `CollectorPort` target shape.** Preserve the durable boundary:

- no fabricated market facts;
- no scraping admitted as a target source merely to fill an evidence gap;
- external market acquisition remains behind D4 adapter/consumer-owned semantics;
- legitimate future official/vendor/manual sources require their own source/coverage/trust contract rather than a universal collector abstraction.

### ADR-032 — catalog-offers flag defaults off

**SUPERSEDE target meaning.** The current environment flag/default-off behavior is runtime/current-state evidence only. Provider-effective support plus a real consumer determines the target capability. Runtime feature-toggle mechanics, if any remain useful, belong D7.

## 7.4 Gate M1 — bounded real Market probe

Read-only probe against one suitable current bound MLB listing:

- obtain current `price_to_win` evidence;
- when the listing/catalog relation legitimately supports it, obtain the related catalog offer population;
- verify seller/item/catalog attribution;
- verify the meaning of no-result/not-applicable/unavailable remains distinguishable;
- verify pagination/parent-child behavior where materially reached;
- do not claim general market completeness beyond the exact provider catalog/competition scope.

**Gate purpose:** prove the first concrete provider Market Evidence lane, not prove all marketplace competition.

---

# 8. B4-B — Expected / Order Economic Evidence

## 8.1 L0 Expected Economics evidence

Commercial Economics owns L0. D4 provides only the external evidence needed to estimate the provider-dependent components.

Current first-flow evidence families:

```text
candidate sale context
  ├─ Sankhya Expected Tax      ← accepted B3
  ├─ Mercado Livre expected sale-fee evidence
  ├─ Mercado Livre expected seller-shipping evidence
  └─ promotion/discount/provider-rule evidence only when materially applicable
           ↓
Commercial Economics
           ↓
L0 Expected Economics
```

Expected sale fee and expected seller shipping cost remain distinct provider evidence components.

## 8.2 Expected sale-fee qualification

Current official Mercado Livre `listing_prices` is the selected first candidate for provider sale-cost quotation.

For current Brazil semantics, target requests must preserve every provider qualifier materially required by the applicable rule, including where applicable:

- site/currency;
- price;
- listing type;
- category;
- quantity;
- logistics type;
- shipping mode;
- billable-weight or other provider-required dimensions when applicable to the current country/rule.

A successful quote missing a material qualifier is not sufficient evidence merely because it returned HTTP 200.

### Current-code disposition

The present `FeeQuote` implementation sends only price/listing type/category and returns commission percentage + fixed fee. Because the current Brazil fixed-fee structure is logistics-sensitive, **the current contract is evidence only and must not become target authority by inheritance.**

B4 does not freeze the target DTO shape or implementation package; it freezes the qualification property.

## 8.3 Expected seller-shipping evidence

Where Commercial Economics needs seller-borne expected shipping cost, D4 obtains it from the applicable provider shipping-quotation/cost surface rather than folding it into sale commission by convention.

The quote must preserve the material item/logistics/free-shipping/price/source context needed by the provider operation. Unknown provider shipping cost remains Unknown, not zero.

## 8.4 L1 Order Economics evidence

Once a Sale exists, L1 uses transaction-specific provider evidence rather than re-running L0 and calling the result realized.

Potential first-flow evidence includes:

- actual Order unit price / quantity / discounts where authoritative for the transaction;
- `order_items[].sale_fee` or equivalent transaction fee evidence;
- Payment references/amount evidence needed to locate external financial movements;
- realized Shipment seller cost from the authoritative Shipment cost surface;
- accepted B3 native/fiscal results where Materialization provides attributable evidence.

Rules:

1. L0 quote is not substituted for L1 when transaction-specific evidence exists.
2. Embedded Order Payment data may be useful correlation/operational evidence but is not L2 settlement authority by itself.
3. Buyer shipping charge and seller shipping cost remain distinct.
4. Discount/promotion effects preserve who funded/bore them where the provider evidence permits; unsupported attribution remains Unknown.
5. Commercial Economics decides economic component classification and attribution.

## 8.5 ADR-009 disposition candidate

**PRESERVE the provenance invariant; SUPERSEDE the legacy storage/layer interpretation.**

Durable property:

> a material economic value is not trustworthy without enough source, scope and time/provenance to understand what the value represents.

B4 does **not** inherit:

- one `channel_fees` table;
- layer 1/2/3 as canonical economic ontology;
- one global resolution ladder;
- config fallback as equivalent provider evidence;
- a universal fee ledger.

## 8.6 Gate E1 — bounded Expected/Order probe

Using one suitable current bound MLB item/order context, read-only:

1. obtain a current `listing_prices` quote with the **full material logistics context**;
2. obtain expected seller-shipping cost when applicable;
3. read one comparable actual Order/Shipment context and preserve transaction-specific fee/shipping evidence;
4. demonstrate that the contract can distinguish expected fee, expected shipping, order fee and realized shipment cost;
5. demonstrate that missing provider qualifiers/components remain Unknown rather than silently defaulted.

**Gate purpose:** prove economically qualified external input surfaces, not calculate margin inside D4.

---

# 9. B4-C — Financial Movement / Settlement Evidence

## 9.1 Provider-independent contract

D4 exposes source-qualified financial movement evidence without creating synthetic MPC Payment/Refund/Settlement identities.

A material external movement preserves, proportionately:

- Organization + external account/source namespace;
- provider-native movement/source identity;
- provider movement kind/status;
- currency;
- gross/net/component values supplied by the source;
- source occurrence/approval/settlement/release time where exposed;
- acquisition/report time;
- Order/Payment/Shipment/Pack/external-reference correlation evidence where exposed;
- report/period/page/coverage provenance needed to assess completeness;
- reversals/refunds/chargebacks/adjustments as distinct occurrences where the source models them distinctly.

D4 does not decide:

- which Sale/Resolution/period the movement economically belongs to when correlation is incomplete;
- whether two provider movements economically offset each other;
- realized profit/margin;
- reconciliation closure.

Those are Commercial Economics Economic Attribution / R2 responsibilities.

## 9.2 Evidence authorities must remain distinct

### Order / Payment operational evidence

Useful for current transaction/correlation state; not by itself proof of account settlement.

### Mercado Livre Billing

Provider billing/charge/fiscal-report evidence. Current official guidance treats Billing integration as report-oriented/periodic and recommends operational Order/Shipment resources for real-time needs.

Therefore Billing may establish billed commission/charge/rebate/discount/shipping-document evidence but does not become current Order state or cash-settlement authority merely because it contains monetary totals.

### Mercado Pago account / released-money reports

Candidate source for actual account-impact/release/settlement movements. Current official reports expose transaction families such as settlement/refund/chargeback/dispute/withdrawal/payout/shipping variants and financial/correlation fields such as net account impact, source/order/shipping identifiers and fee components.

These reports remain provider/payment-native evidence; they do not create an MPC treasury/bank domain.

## 9.3 Settlement completeness / coverage

A successful page/report is never labeled complete beyond the source-defined interval/scope actually traversed.

B4 must establish for the selected settlement lane:

- account namespace qualification;
- date/period boundaries;
- pagination/cursor/from-id semantics where applicable;
- report generation/read lifecycle when the source is asynchronous;
- retention/history limits where material;
- partial response semantics such as HTTP 206 where applicable;
- duplicate/reissued/report-overlap behavior sufficient to avoid double-attribution;
- a reread/recovery strategy at the contract level without choosing D7 schedule/retry machinery.

Provider report cadence/cache recommendations are D7 mechanics, not D4 business meaning.

## 9.4 Economic correlation evidence

D4 must preserve available native anchors; it must not fabricate one universal correlation key.

Potential anchors include:

- provider Order ID;
- Payment/source ID;
- Shipment ID;
- Pack ID;
- external reference;
- provider billing detail/document IDs;
- movement/report identifiers.

Commercial Economics owns whether the available evidence establishes exact, partial, ambiguous or unresolved attribution.

A missing anchor is a real unresolved Economic Attribution condition, not permission to assign by amount/date similarity silently.

## 9.5 R2 / R3 fence

```text
L1 Order Economics
      ↕ R2 — Commercial Economics
provider authoritative realized/settlement evidence
```

B4 supplies the external evidence for R2.

```text
provider payout/settlement
      ↕ R3 — Commercial Economics
Bank Cash Receipt
```

B4 may provide the provider payout side when available. **R3 remains unclaimed until an accepted bank source exists.** ERP receivable/baixa does not silently become bank cash or marketplace-settlement authority.

## 9.6 Gate S1 — Real Settlement Evidence Gate — PRIMARY B4 CLOSURE GATE

Against the currently bound real seller/payment account, use the smallest sanctioned read/report operations needed to establish:

1. which authorized Mercado Livre Billing surface is available for the account and which concrete fee/rebate/charge/shipping evidence it returns;
2. which authorized Mercado Pago report/account surface is available for actual financial movements/release/account impact;
3. one bounded real Sale/Order with enough provider references to test correlation through at least Order/Payment and, where available, Shipment into billed and settled/released evidence;
4. whether refund/adjustment/reversal-like movement classes are structurally representable even if the selected sample does not contain every class;
5. the exact coverage/pagination/report lifecycle of the surfaces used;
6. no assumption that Billing amount = settlement amount or Order payment = account cash merely because totals happen to match in one sample.

If the bound account/credentials cannot expose a materially required Product 1.0 realized/settlement evidence source, B4 returns **STOP / SPLIT PREREQUISITE** for that capability. It does not substitute Sankhya receivables, a manual spreadsheet, or a guessed Order net value as marketplace settlement authority.

**S1 is the primary candidate B4 closure gate.**

---

## 10. No B4 write surface by symmetry

No current Product 1.0 requirement identified in the authority path needs D4-B4 to mutate competitor market data, provider billing, Mercado Pago settlement, payout or financial movements.

Therefore the B4 target is **read/evidence-only unless new material evidence proves a required write**.

Report-generation endpoints may create provider report artifacts as read-support mechanism. If a real gate requires report generation, the reviewer/operator must treat that operation according to the provider's actual external-effect semantics and obtain the bounded authorization required by repository safety rails; generating a report does not create Commercial Economics authority.

---

## 11. YAGNI / overengineering fence

B4 MUST NOT introduce:

- generic Financial Transaction / Payment / Refund / Settlement MPC business entity merely for provider normalization;
- universal financial ledger owned by D4;
- universal fee ledger / `channel_fees` target table by inheritance;
- universal MarketObservation/Competitor entity graph before a real consumer proves it;
- generic `CollectorPort` framework/plugin registry;
- scraping infrastructure;
- full market-history backfill without authoritative historical evidence;
- scheduler/report worker topology;
- one global `Fee` type that collapses expected/order/billed/settled semantics;
- one generic correlation key inferred from amount/time;
- provider Billing documents as settlement/cash authority;
- Mercado Pago account balance as bank-cash authority;
- ERP receivable/baixa as marketplace settlement authority;
- margin/profit calculation inside D4;
- automatic price recommendation from `price_to_win`;
- support for every Mercado Pago report product or every financial movement class before Product 1.0 requires them.

---

## 12. Proof strategy / adversarial challenge

B4 must survive at least:

1. `price_to_win` unavailable/not-applicable for a listing → no fabricated competitor conclusion.
2. Catalog offer query returns only part of a parent/child population → no complete-market claim.
3. Provider quote returns 200 but a material logistics parameter was omitted → quote is not accepted as economically sufficient.
4. Expected shipping unavailable → not zero.
5. Order fee differs from current quote due price/rule/promotion/time change → preserve both L0/L1 evidence rather than overwrite history.
6. Buyer-paid shipping and seller cost differ → no one-field `shipping` collapse.
7. Billing fee differs from Order `sale_fee` due rebate/discount/adjustment → preserve separate evidence and let Economics reconcile.
8. Payment status is approved but money is not yet released/settled → L1 does not become L2.
9. Refund/chargeback arrives after earlier settlement → prior evidence remains history; new movement participates in Economic Attribution/Reconciliation.
10. One movement has no exact Order anchor → attribution remains unresolved/ambiguous.
11. Report page is partial/206 or traversal stops → no settlement completeness claim.
12. Overlapping report intervals expose the same native movement → native identity/provenance prevents double attribution at the semantic boundary; D7 later chooses persistence/runtime mechanism.
13. Payout exists with no bank integration → provider payout is known; bank receipt remains Unknown/unreconciled R3.
14. Current `channel_fees` code is convenient → convenience does not make its table/layers target authority.
15. Current catalog-offers flag is off → runtime flag does not prove provider capability absent.
16. A second marketplace/payment provider arrives → extend concrete evidence adapters through the same domain-owned semantics without creating false cross-provider ontology unless repetition proves a smaller shared mechanism.

---

## 13. Gate plan — intentionally small

### M1 — Market Evidence probe

**Read-only / small.** One suitable real listing; competition + catalog-offer evidence if applicable.

### E1 — Expected/Order Economic Evidence probe

**Read-only / small.** One suitable item/order/shipment context; fully qualified expected fee + expected shipping + transaction-specific fee/shipping evidence.

### S1 — Settlement Evidence Gate

**Primary closure gate.** Bounded real Billing + Mercado Pago/account settlement evidence and correlation/coverage proof.

No micro-gates are added unless one of these three exposes a genuinely distinct correctness failure class.

---

## 14. Preliminary decision

**CURRENT STRUCTURE CONFIRMED at the D1/D2 semantic boundary, with bounded D4 restructuring of legacy evidence contracts.**

Proposed target direction:

```text
Market Evidence ───────────────→ Market Intelligence
                                     ↓
                               Commercial Economics

Expected fee/shipping ─────────→ L0 Expected Economics
Order/Shipment economics ──────→ L1 Order Economics
Billing + Payment movements ───→ Economic Attribution
                                 + L2 Realized Economics
                                 + R2 Reconciliation
Provider payout ────────────────→ R3 provider side only
Bank Cash Receipt ──────────────→ deferred until accepted bank source
```

No generic finance/integration framework is justified.

B4 remains **OPEN / NON-AUTHORITATIVE** until independent challenge, M1/E1/S1 evidence, GPT adjudication and explicit operator ratification.

---

## 15. Reopen / stop triggers

Reopen only the implicated decision when material evidence shows:

1. Market Intelligence requires a real market source that cannot fit source-qualified observation + comparability ownership.
2. Current Mercado Livre market surfaces cannot satisfy a Product 1.0 competitive claim and no legitimate source is available → narrow product/evidence decision; no scraping by convenience.
3. Provider fee/shipping semantics require a materially new business authority absent from D1 → targeted D1 review.
4. Provider monetary evidence cannot fit current D2 source-qualified external-movement identity → targeted D2 review.
5. Settlement effect/occurrence semantics cannot fit D3 current/historical/duplicate/partial/recovery rules → targeted D3 review.
6. The bound Mercado Livre/Mercado Pago account cannot expose materially required L2 settlement evidence through a sanctioned surface → **STOP / SPLIT PREREQUISITE** rather than fabricating realized economics.
7. A future bank integration becomes an accepted Product 1.0 source → open the R3 bank side in its responsible stage without moving bank authority into B4 by assumption.
8. A real second provider exposes repeated technical evidence mechanics whose duplication is materially worse than a small shared non-authority mechanism → consider only that proven mechanism.
9. Product 1.0 introduces a genuine B4 external write requirement → return to D3/D4 external-effect contract before adding the write.

Naming preference, current module convenience, hypothetical providers and desire for one unified ledger are not reopen evidence.

---

## 16. Evidence basis for independent review

Repository evidence considered:

- accepted D1 Market Intelligence / Commercial Economics ownership;
- accepted D2 external financial-movement identity + Economic Attribution/Reconciliation semantics;
- accepted D4-B1/B2/B3 external-evidence/coverage/provider-boundary rules;
- ADR-009 provenance history;
- reopened ADR-014 / ADR-020 / ADR-032;
- current Mercado Livre adapter `FeeQuote`, catalog competition/offers and Order/Shipment evidence as current-state evidence only.

Current official external evidence revalidated on 2026-08-18 includes documentation families for:

- Mercado Livre catalog competition / `price_to_win`;
- Mercado Livre sale costs / `listing_prices` including current logistics-sensitive fixed-fee qualification for Brazil;
- Mercado Livre shipping-cost evidence;
- Mercado Livre Billing integration/details and its report-oriented consumption guidance;
- Mercado Pago account-money and released-money reports, including native movement kinds, net account-impact values and Order/Shipment/source correlation fields.

These sources are evidence, not target authority. If provider behavior changes materially, reopen only the contract that depended on it.
