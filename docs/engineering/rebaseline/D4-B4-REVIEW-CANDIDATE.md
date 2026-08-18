# D4-B4 — Market / Economics / Settlement Contract — REVIEW CANDIDATE

> **Status:** REVIEW CANDIDATE / NON-AUTHORITATIVE / DISPOSABLE  
> **Program:** Architecture Rebaseline / Technical System Design  
> **Parent authority:** accepted D0–D4-B3 only  
> **Method:** DevelopmentConexus Engineering Method v1.0.0  
> **Base authority HEAD:** `e8a0b60f4478f2e41827ae3174f39e5eca61dfcb`  
> **Initial B4 candidate commit:** `9c212afc58bbce037cc85582329a3652fd2d665d`  
> **Independent review evidence:** current `AI-DIALOG.md`, 2026-08-18 — NON-AUTHORITATIVE  
> **Purpose:** post-review coherent disposable surface for final B4 evidence-gate discharge and operator adjudication. This file is not target authority and MUST be deleted before canonical consolidation.

---

## 1. Review question and scope

D4-B4 must answer:

> **Which concrete external market, fee, billed-charge, payment/release and financial-movement contracts are required for Product 1.0, and how can D4 provide source-qualified evidence with honest coverage/provenance while Market Intelligence, Commercial Economics and Post-Sale Resolution retain their accepted semantic authorities?**

B4 must not create a finance domain, marketplace ledger, universal payment model, scraping subsystem or economic-calculation authority inside D4.

The target remains deliberately divided into only three coherent evidence families:

1. **Market Evidence** — externally authoritative competitive observations for Market Intelligence / Commercial Economics.
2. **Expected / Order Economic Evidence** — provider fee, shipping and transaction-economic evidence for L0 Expected Economics and L1 Order Economics.
3. **Financial Movement / Release / Billed-Charge Evidence** — provider/payment-native evidence for Economic Attribution, L2/R2 and applicable refund-consequence closure.

Implementation remains blocked until D9.

---

## 2. Imported authority — do not re-decide

B4 imports:

1. **Market Intelligence** owns external comparability, competitive position/change and market-evidence sufficiency.
2. **Commercial Economics** owns Cost Basis, economic interpretation, pricing analysis, L0 Expected Economics, L1 Order Economics, L2 Realized/Settlement Economics, Economic Attribution and Economic Reconciliation.
3. **Post-Sale Resolution** owns coordination/correlation/closure of material cancellation/return/refund consequences; provider refund/chargeback/adjustment evidence may legitimately feed that closure without transferring financial interpretation authority.
4. D4 owns concrete external acquisition/protocol/capability/evidence only.
5. Provider-native Payment, Refund, Fee, Adjustment, Settlement, Release, Withdrawal/Payout and equivalent movements remain external/source-qualified identities/evidence.
6. Economic Attribution is MPC-owned Commercial Economics state; D4 does not decide what a provider movement means economically.
7. R1 Expected↔Order and R2 Order↔authoritative realized evidence are Commercial Economics reconciliation semantics.
8. R3 provider payout/withdrawal↔Bank Cash Receipt exists only when an accepted bank source exists; observing the provider side does not open/close the bank side.
9. Marketplace Order, Shipment, Billing/report evidence and provider financial movements remain distinct source resources/evidence classes.
10. Known / absent / unknown / unavailable / partial remain distinct; provider unavailability or partial coverage never becomes zero or complete.
11. Money/tax/cost/pricing values preserve exactness and D2 provenance sufficient for the claim.
12. Consumer owns meaning; adapter owns provider protocol.
13. Provider PII is minimized.
14. Current code/modules/tables/ADRs are evidence only.
15. No B4 external write is introduced merely for symmetry.

No D0/D1/D2/D3 or D4-B1/B2/B3 reopen is proposed by this candidate.

---

## 3. Evidence classification after independent review

### Known

- Current official Mercado Livre competition surfaces include item-level catalog competition / `price_to_win` evidence and catalog-product offer evidence where provider topology supports them.
- Current official Mercado Livre `listing_prices` is a read-only sale-cost quotation surface. Current Brazil fixed-fee behavior is logistics-sensitive; the provider documents logistics qualifiers and warns that omitting material qualifiers can produce a calculated cost different from what the seller is charged.
- Current repository `FeeQuote` evidence calls `listing_prices` with price/listing type/category only and fabricates a default currency in its snapshot; that contract is current-state evidence only and is not target authority.
- Mercado Livre exposes a separate seller-shipping estimation surface for the approximate amount the seller pays for an item/shipping context; expected seller shipping is therefore not the same evidence as sale commission/fixed fee.
- Shipment cost evidence distinguishes seller-side and buyer-side costs; those values cannot collapse into one generic `shipping` amount.
- Current Mercado Livre Order evidence exposes transaction-specific item/fee and Payment references. Current repository evidence measured `sale_fee` as **per unit**; target B4 generalizes only the correctness property that provider fee granularity must be preserved explicitly.
- Mercado Pago Payment evidence distinguishes payment approval from money-release timing; `approved` does not imply `released`.
- Mercado Pago account/release report evidence exposes account-impact/release transaction data and monetary components with native correlation fields where available.
- Mercado Livre Billing is billed-charge/fiscal reconciliation evidence, not current Order authority and not cash/bank authority.
- ADR-014, ADR-020 and ADR-032 are explicitly reopened for D4-B4; their runtime/module/flag shapes have no target authority by inheritance.
- ADR-009's durable provenance principle already has an active home in D2; its legacy `channel_fees` table/layer ladder is not target authority.
- The first independent B4 review could not execute live M1/E1/S1 because no Mercado Livre/Mercado Pago credential was available in that reviewer environment. This is an access/proof prerequisite, not evidence that the B4 boundary failed.

### Inferred

- The smallest coherent external-economic boundary is operation-specific evidence contracts rather than one generic financial ledger/resource graph.
- `price_to_win` / catalog-offer observations are provider market evidence, not an MPC recommended price or comparability conclusion.
- A **Payment-first read-only path** is the smallest candidate source set for per-sale R2; whether it is materially sufficient must be measured on the bound account before B4 closes.
- Billing and payment/release evidence may contain overlapping monetary values while answering different questions; numeric equality does not collapse their authority.

### Unknown / closure-critical until proven

- Exact current access/scope of the bound Mercado Livre/Mercado Pago credentials for the required E1/S1 read operations.
- Whether the selected current `listing_prices` request plus required qualifiers is materially sufficient for the first MLB L0 fee claim; HTTP 200 alone is not proof.
- Whether the concrete expected seller-shipping surface returns materially sufficient seller-borne shipping evidence for the selected first lane.
- Whether Payment-level evidence alone is enough for the selected per-sale R2/L2 claim, including exact fee/net/release/correlation semantics.
- Whether a broader account-movement/recovery surface is additionally required for Product 1.0 correctness beyond the per-sale Payment path.
- Exact real-account pagination/retention/completeness behavior for any broader account/report lane that proves necessary.
- Exact strongest native anchors for every material refund/adjustment/chargeback movement class.

### Unknown but not automatically B4-blocking

- Whether every Product 1.0 listing has applicable catalog competition evidence; lack of catalog evidence can legitimately become Market Intelligence `insufficient evidence`.
- A bank cash source. R3 bank-side reconciliation remains unclaimed.

### Deferred

- collection/report cadence, caches, scheduled polling, retries/backoff, checkpoints and report-worker topology — D7;
- UI representation and operator workflows — D6;
- end-to-end economic golden-flow proof — D8;
- bank-cash R3 proof until a bank source is explicitly accepted.

---

## 4. Root cause

Legacy/current implementation evidence mixes distinct questions under convenient modules/tables:

- market observations vs comparability/economic interpretation;
- expected fee vs expected seller shipping;
- expected fee vs transaction fee;
- Order/Payment evidence vs money release/account impact;
- billed charges/rebates vs realized movement evidence;
- refund movement evidence vs Post-Sale closure meaning;
- value provenance vs one shared `channel_fees` persistence model;
- provider access switches vs provider-effective capability.

That makes silent authority collapse reachable. A plausible number can appear authoritative merely because it came from a provider API or shared table even when its scope, lifecycle, granularity or qualification differs from the question being answered.

The concrete defect class is already reachable: a provider read may return HTTP 200 while a material request qualifier is absent/ignored, producing a plausible but economically wrong value. B4 therefore needs claim-specific qualification and falsification, not endpoint-existence proof.

---

## 5. Target invariant

> **Every external market/economic value entering MPC is source-qualified, scope-qualified, granularity-qualified and time/provenance-qualified strongly enough for the exact claim being made. D4 preserves only what the external source proves; Market Intelligence owns comparability, Commercial Economics owns economic interpretation/attribution/reconciliation, and Post-Sale Resolution owns applicable refund-consequence closure. Expected, order-time, billed, approved, released, withdrawn/payout and bank-cash evidence never collapse merely because they contain similar amounts.**

Corollaries:

- external evidence is not economic meaning;
- `price_to_win` is not MPC Price Intent or automatic price recommendation;
- provider fee quote is not Order fee;
- expected seller shipping is not buyer shipping charge and not realized seller shipping cost;
- transaction fee evidence carries provider granularity; per-unit, per-line and per-order are not interchangeable;
- Payment `approved` is not money `released`;
- Billing charge/rebate is not money release/account impact;
- release/account impact is not withdrawal/payout;
- withdrawal/payout is not Bank Cash Receipt;
- one external refund/chargeback movement may be evidence for both Commercial Economics and Post-Sale Resolution without either acquiring the other's meaning;
- one persistence table is not required because several evidence classes contain Money/provenance;
- provider/runtime feature flags are not capability authority.

---

## 6. Credible alternatives / Global Maximum

### A — Preserve legacy `market` / `channel_fees` / profitability structure as target

**REJECT.** Promotes current module/table/layer boundaries into target authority and preserves evidence-lifecycle conflation.

### B — Create a generic D4 Financial Movement / Economic Evidence ledger

Examples: universal `Payment`, `Settlement`, `Fee`, `MarketObservation`, `FinancialTransaction` or provider-resource graph normalized across future providers.

**REJECT.** External identities remain source-qualified under D2 and business meaning already has D1 owners. A D4 ledger would create duplicate authority and speculative abstraction.

### C — Operation-specific external evidence contracts + domain-owned interpretation

**PROPOSED GLOBAL MAXIMUM / INDEPENDENT REVIEW CORE SURVIVED.**

```text
external provider/payment source
        ↓
D4 concrete adapter
  - source identity
  - operation-specific coverage
  - provider-native components + granularity
  - source occurrence/update/release time
  - acquisition provenance
  - available native correlation anchors
        ↓
Market Intelligence / Commercial Economics / Post-Sale Resolution
  - comparability
  - L0/L1/L2 interpretation
  - Economic Attribution / R1/R2
  - refund-consequence closure where applicable
```

Prepare seams for later providers; do not model their ontology now.

---

# 7. B4-A — Market Evidence

## 7.1 Provider-independent contract

D4 supplies market observations only when an admitted external source establishes them strongly enough for the consuming claim.

A market observation preserves proportionately:

- Organization + Marketplace Installation/provider source;
- provider product/listing/catalog scope;
- offer/price/competitive-status evidence;
- currency/source dimensions required to compare;
- provider occurrence/update time where exposed;
- acquisition time;
- coverage/pagination/completeness state;
- provider identifiers as external references, not MPC identities.

Market Intelligence owns comparability, competitor relevance, competitive position/change and evidence sufficiency.

Commercial Economics may consume Market Intelligence meaning; it does not independently reinterpret raw competitor payload when Market Intelligence owns comparability.

## 7.2 Current Mercado Livre realization candidate

Current official first-provider evidence includes:

1. `GET /items/{ITEM_ID}/price_to_win?...&version=v2` when catalog competition applies.
2. Catalog product offer population through the applicable current catalog product/items surfaces.

Rules:

- provider winning/competing/listed/sharing/boost vocabulary stays adapter-local;
- catalog membership proves a provider relation, not full Market Intelligence comparability;
- non-catalog/unavailable competition evidence becomes explicit insufficient/unavailable evidence, never “no competitors”;
- no general-market completeness claim is derived from one provider catalog population.

## 7.3 Source-admissibility / legacy ADR disposition candidate

### ADR-014 — on-demand/local Docker

**HISTORICAL after B4 if ratified.** Its runtime choice has no surviving D4 target meaning. Honest absence/history is already carried by D0/D3/D4; collection runtime/cadence is D7.

### ADR-020 — generic `CollectorPort`

**SUPERSEDE target shape.** One current market source does not justify a generic collector framework.

The surviving source-admissibility property is broader than B4 and therefore must **not** remain homeless inside B4 after ADR-020 retirement. At canonical consolidation, rehome the stable platform rule into `ARCHITECTURE.md` external-integration principles before marking ADR-020 historical/superseded:

> **Absence of an admitted external market-data source never authorizes fabricated evidence or an unadjudicated scraping path by convenience. A materially new market-data source requires explicit source, legality/trust, coverage and provenance adjudication before its evidence can support MPC claims.**

This does not create a permanent ban on every possible future lawful collector; it prevents an unreviewed source from becoming truth because another source is missing.

### ADR-032 — catalog-offers flag defaults off

**SUPERSEDE target meaning.** Current env-flag behavior is runtime/current-state evidence only. Provider support + current context + consumer need determine capability; feature-toggle mechanics, if retained, belong D7.

## 7.4 M1 — Market Evidence lane-selection proof

**Classification:** `INSTALLATION / LANE-SELECTION EVIDENCE`, not a standing B4 closure gate.

Read-only bounded proof should include, where the current Installation permits:

- one catalog-applicable listing: `price_to_win` plus offer population when applicable;
- one negative/control listing that is non-catalog, moderated/unavailable or otherwise reaches a materially different competition outcome;
- seller/item/catalog attribution;
- observed vs known-empty/not-applicable vs unavailable vs partial distinction;
- pagination/parent-child behavior when reached.

M1 proves which current provider Market Evidence lane exists and that failure states remain honest. If the current first flow has **no legitimate market-evidence lane at all** and D0 Competitive Intelligence would therefore be impossible rather than honestly insufficient for particular items, return to targeted product/B4 adjudication. Otherwise an individual insufficient-market case is not a batch-closure failure.

---

# 8. B4-B — Expected / Order Economic Evidence

## 8.1 L0 Expected Economics evidence

Commercial Economics owns L0. D4 supplies only provider-dependent evidence:

```text
candidate sale context
  ├─ Sankhya Expected Tax          ← accepted B3
  ├─ ML expected selling-cost/fee
  ├─ ML expected seller shipping
  └─ provider promotion/discount evidence when materially applicable
              ↓
Commercial Economics
              ↓
L0 Expected Economics
```

Expected selling cost/fee and expected seller shipping remain distinct evidence components even if provider rules internally interact.

## 8.2 Expected sale-cost / fee qualification

Current Mercado Livre `listing_prices` is the first candidate sale-cost quotation surface.

For current Brazil semantics, the request/evidence contract preserves every provider qualifier materially required by the applicable rule, including where applicable:

- site/currency;
- price;
- listing type;
- category;
- quantity;
- `logistic_type`;
- `shipping_mode`;
- `billable_weight` or another currently documented qualifier where applicable.

A successful provider response does **not** prove these inputs were consumed correctly.

### Decorrelation control

E1 must falsify the assumed mechanism rather than merely obtain one plausible quote. Against a fixed base price/category/listing type, vary one material qualifier at a time where the surface accepts it, for example:

- logistics type;
- shipping mode;
- quantity;
- billable weight when applicable.

Record which returned components actually move. A qualifier claimed to be material but ignored/silently ineffective cannot be treated as proven merely because the endpoint returns 200.

If the selected expected-cost component depends on information the surface cannot represent, use another sanctioned surface that actually represents it or keep that component Unknown. Do not force `listing_prices` to be sufficient by assumption.

### Current-code disposition

The current `FeeQuote` sends price/listing type/category only and returns a reduced fee snapshot. It remains current-state evidence only.

## 8.3 Expected seller-shipping evidence

For the selected MLB first lane, the concrete current candidate is the provider shipping-estimation surface:

```text
GET /users/{USER_ID}/shipping_options/free
```

using `item_id` or the provider-required dimensions/context plus applicable price/listing/logistics/free-shipping qualifiers.

D4 preserves the **seller-borne expected shipping** evidence returned for the requested context. It does not rename buyer shipping charge or a generic freight amount into seller cost.

If the selected normal flow materially needs expected seller shipping and the sanctioned surface cannot represent/return it sufficiently, E1 remains OPEN; Unknown is not zero.

## 8.4 L1 Order Economics evidence

Once a Sale exists, L1 uses transaction-specific evidence instead of re-running L0 and calling it actual.

Potential first-flow evidence includes:

- authoritative Order unit price / quantity / discounts;
- `order_items[].sale_fee` or equivalent provider transaction-fee evidence;
- Payment references needed for later financial correlation;
- authoritative Shipment seller-side cost evidence;
- accepted B3 attributable native/fiscal results.

Rules:

1. L0 quote never substitutes for L1 when transaction-specific evidence exists.
2. Provider transaction fee preserves its **native granularity** as evidence — per-unit/per-line/per-order or equivalent. The adapter never assumes the arithmetic aggregation law from a field name.
3. For the currently measured ML shape, `sale_fee` was observed per unit; quantity must therefore remain available for the consuming economic interpretation rather than silently treating the value as line total.
4. Buyer shipping charge and seller shipping cost remain distinct.
5. Discount/promotion effects preserve who funded/bore them where provider evidence permits; unsupported attribution remains Unknown.
6. Commercial Economics owns component classification, aggregation and attribution.

## 8.5 Billing / charged-fee evidence

Mercado Livre Billing is a **billed-charge / rebate / bonus / fiscal reconciliation** evidence family. It is not current Order authority and it is not the cash/release authority for S1.

Billing may be used later to explain material divergence such as:

```text
Order transaction fee
≠ billed charge after rebate/bonus/adjustment
```

That evidence feeds Commercial Economics Economic Attribution/R2 analysis where material. B4 does not require Billing merely to prove that a Payment was released.

## 8.6 ADR-009 disposition candidate

ADR-009's durable provenance rule is **already homed in D2**. B4 cites that authority and adds only provider-specific evidence requirements; it does not create a second provenance authority.

Legacy B4 does **not** inherit:

- one `channel_fees` table;
- layer 1/2/3 as canonical ontology;
- one global resolution ladder;
- config fallback as equivalent provider evidence;
- a universal fee ledger.

## 8.7 E1 — Expected / Order Economic Evidence — CLOSURE GATE

**Classification:** `B4 CLOSURE GATE / OPEN` until live proof passes.

Read-only bounded proof using a current suitable MLB item plus comparable real Order/Shipment:

### Expected side

1. execute `listing_prices` with the exact current material qualifiers for the selected context;
2. execute the decorrelation controls from §8.2;
3. obtain expected seller-shipping evidence through the concrete shipping-estimation surface in §8.3;
4. preserve returned component scope/granularity/provenance;
5. fail honestly when a required qualifier/component is not representable.

### Order side

6. reread one real Order and preserve price/quantity/discount/transaction-fee evidence;
7. preserve provider fee granularity explicitly;
8. reread the related Shipment seller-side cost evidence when applicable;
9. prove that expected fee, expected seller shipping, Order transaction fee and realized seller shipment cost remain distinct evidence.

E1 does **not** require numerical equality between expected and historical Order evidence when time/rule/promotion context differs. It requires a contract strong enough to explain which components are being compared and why a variance is or is not attributable.

**PASS condition:** the selected first Product 1.0 lane has sanctioned, materially sufficient and falsified provider evidence for the expected provider-cost components required by L0 and the transaction-specific evidence required by L1.

**STOP / SPLIT PREREQUISITE:** a materially required selected-flow expected component cannot be represented honestly through any sanctioned admitted surface.

---

# 9. B4-C — Financial Movement / Release / Billed-Charge Evidence

## 9.1 Provider-independent contract

D4 exposes source-qualified provider/payment evidence without creating synthetic MPC Payment/Refund/Settlement identities.

A material external financial occurrence preserves proportionately:

- Organization + external account/source namespace;
- provider-native occurrence/resource identity;
- provider-native kind/status;
- currency;
- gross/net/component values supplied by the source;
- provider occurrence/approval/release/other source times separately;
- acquisition/report time;
- native Order/Payment/Shipment/external-reference anchors where actually exposed;
- coverage/page/period provenance when the source is population/report based;
- refunds/chargebacks/adjustments/withdrawals as distinct occurrences when the source models them distinctly.

D4 does **not** decide:

- which economic scope an incompletely correlated occurrence belongs to;
- whether two occurrences economically offset each other;
- realized margin/profit;
- R2 closure;
- whether provider release means the same business conclusion as provider payout/withdrawal.

Commercial Economics owns those interpretations.

Refund/chargeback/adjustment evidence also legitimately feeds **Post-Sale Resolution** when its explicit consequence scope requires financial closure evidence. Post-Sale consumes the evidence for its own closure; it does not acquire Economic Attribution or R2 authority.

## 9.2 Evidence rungs remain distinct

### Marketplace Order / embedded Payment

Current transaction/correlation evidence. Not L2 by itself.

### Mercado Pago Payment resource

Candidate per-sale evidence for payment status, release timing, net/refund/fee components and native correlation where actually exposed.

**Approved is not released.** D4 preserves both separately.

### Money release / account-impact evidence

A provider account may expose movement/report evidence showing money released or account impact. D4 preserves the provider's native fields/meaning; it does not relabel all such records as one MPC `Settlement` state.

Commercial Economics decides whether the established evidence is sufficient for the accepted L2/R2 conclusion.

### Mercado Livre Billing

Billed charge/rebate/bonus/fiscal-report evidence only. It may explain charge divergence but is not required merely to prove release/account impact.

### Withdrawal / Payout

Provider-side R3 evidence when exposed. It does not prove bank receipt.

### Bank Cash Receipt

Unclaimed until an accepted bank source exists.

## 9.3 Correlation evidence

D4 preserves only native anchors actually present, for example:

- Order ID;
- Payment/source ID;
- Shipment ID;
- external reference;
- billing detail/document identifier;
- movement/report identifier.

Do not require Pack ID or another anchor merely because it exists in some payload. Add an anchor only when the selected correlation path actually consumes it.

Commercial Economics owns whether available anchors establish exact, partial, ambiguous or unresolved Economic Attribution. Missing exact correlation never becomes amount/date heuristic auto-assignment.

## 9.4 S1 — Realized / Release Evidence — CLOSURE GATE

**Classification:** `B4 CLOSURE GATE / OPEN` until live proof passes.

S1 is intentionally staged from the smallest read-only source set outward.

### S1-A — per-sale R2 candidate minimum — read-only first

Starting from one real bound Marketplace Sale/Order:

1. obtain the source Payment reference from authoritative Order evidence;
2. read the current Payment resource through a sanctioned Mercado Pago surface;
3. establish account/source namespace and Payment identity;
4. preserve payment status separately from release state/time;
5. preserve transaction amount, refunded amount, net/fee components and external references where actually exposed;
6. prove whether this evidence is materially sufficient for Commercial Economics to distinguish L1 Order Economics from the relevant per-sale realized/released evidence without fabricating cash/bank state;
7. use a refunded/adjusted example if one is already readable and materially useful, but do not create one.

Payment-level evidence is a **candidate minimum**, not assumed sufficient before measurement.

### S1-B — account-movement universe / recovery — conditional

Open this layer only if S1-A or accepted Product 1.0 correctness proves that per-sale anchored reads cannot cover a material movement/recovery class, for example late chargeback/refund/account adjustment or period completeness.

Then inspect existing read-only account/report/list surfaces and establish:

- native movement identity;
- period/window/pagination/retention semantics;
- duplicate/overlap behavior;
- available Order/Payment/Shipment anchors;
- partial/unavailable states;
- recovery of material unanchored occurrences.

Do not require both Mercado Livre Billing and Mercado Pago reports by symmetry. Choose the smallest source set that actually proves the claim.

### S1 outcome

**PASS:** the bound account exposes a sanctioned evidence set sufficient for the selected Product 1.0 L2/R2 claim, with honest release/account-impact/correlation/coverage semantics.

**CONDITIONED:** per-sale evidence works but a specifically required movement-universe/recovery property still needs an admitted read-only source/proof.

**STOP / SPLIT PREREQUISITE:** a materially required L2/R2 evidence class cannot be established through a sanctioned surface. Do not substitute Sankhya receivables, spreadsheets or guessed Order net values as marketplace/payment realized authority.

---

## 10. No report-generation/write surface by convenience

B4 currently admits **read/evidence-only** external operations for its normal target.

A provider `POST` that creates/generates a report artifact is an external effect; it is not reclassified as "read support" merely because the resulting artifact is later read.

Therefore report generation is **not part of the current B4 target/gates**.

If S1-B later proves that Product 1.0 correctness actually requires a report that does not already exist/readably recur, stop and return to explicit D3/D4 external-effect adjudication plus operator authorization before any generation call. That later decision must define target, intent/correlation, acceptance/pending/ambiguity, quota/blast radius and authoritative status/download reread.

---

## 11. YAGNI / overengineering fence

B4 MUST NOT introduce:

- generic Financial Transaction / Payment / Refund / Settlement MPC business entity for provider normalization;
- universal financial ledger owned by D4;
- universal fee ledger / `channel_fees` target table;
- universal competitor/provider-resource graph;
- generic `CollectorPort` framework/plugin registry;
- unadjudicated scraping infrastructure;
- fabricated/backfilled market history;
- scheduler/report worker topology;
- one `Fee` type that collapses expected/order/billed/released semantics;
- one generic correlation key inferred from amount/time;
- Billing as release/cash authority;
- Payment approval as release;
- provider release as bank receipt;
- ERP receivable/baixa as marketplace/payment realized authority;
- margin/profit calculation inside D4;
- automatic price recommendation from `price_to_win`;
- support for every Mercado Pago report product/movement class before a Product 1.0 correctness claim requires it.

---

## 12. Independent-review adjudication integrated

| Finding | Adjudication in this candidate |
|---|---|
| F-B4-1 — Post-Sale missing as refund evidence consumer | **ACCEPTED** — §9.1 |
| F-B4-2 — expected-cost mechanism underfalsified | **ACCEPTED WITH AMENDMENT** — keep `listing_prices` as candidate, add decorrelation control and fail-honest alternate/Unknown path (§8.2/E1) |
| F-B4-3 — fee granularity unfenced | **ACCEPTED** — §8.4/E1 |
| F-B4-4 — release collapsed into settlement | **ACCEPTED** — §5/§9.2/S1 |
| F-B4-5 — Billing misfiled as settlement gate | **ACCEPTED** — §8.5/§9.2; no longer S1 closure prerequisite |
| F-B4-6 — Payment-first smaller source set | **ACCEPTED AS CANDIDATE MINIMUM** — S1-A; sufficiency still requires live proof |
| F-B4-7 — expected shipping unnamed | **ACCEPTED** — concrete ML shipping-estimation candidate §8.3/E1 |
| F-B4-8 — report generation mislabeled read support | **ACCEPTED** — excluded from current B4 target §10 |
| F-B4-9 — speculative Pack anchor | **ACCEPTED** — removed from required correlation anchors §9.3 |
| C-1 — ADR-009 already homed in D2 | **ACCEPTED** — §8.6 |
| C-2 — ADR-014 residue empty after runtime defer | **ACCEPTED** — historical candidate §7.3 |
| C-3 — no-scraping/source-admissibility rule needs stable home | **ACCEPTED** — planned `ARCHITECTURE.md` rehome only at canonical consolidation after operator ratification §7.3 |

No accepted finding requires reopening D0–D4-B3.

---

## 13. Gate plan — intentionally small

### M1 — Market Evidence lane selection

Read-only. Prove one usable competition lane plus a negative/not-applicable control where the current Installation supports both. **Not a standing B4 closure gate.**

### E1 — Expected / Order Economic Evidence

Read-only. **B4 CLOSURE GATE.** Prove full selected-flow expected fee + expected seller shipping with decorrelation controls, and distinct L1 Order/Shipment evidence with fee granularity.

### S1 — Realized / Release Evidence

Read-only-first. **B4 CLOSURE GATE.** Start with per-sale Payment evidence; widen to existing account-movement reads only if a concrete correctness gap requires it.

No micro-gates are added unless E1/S1 reveals a genuinely different correctness failure class.

---

## 14. Current candidate disposition

**CURRENT STRUCTURE CONFIRMED at the accepted D1/D2 semantic boundary; bounded D4 evidence-contract restructuring required.**

The independent challenge found no smaller architecture core and no D0–D4-B3 reopen. B4 nevertheless remains **CONDITIONED / NON-AUTHORITATIVE** because the first reviewer environment lacked ML/MP credentials and therefore executed **zero live M1/E1/S1 probes**.

Current status inside this candidate:

```text
Architecture core            SOUND / REVIEWED
M1 market lane               PENDING REAL READ
E1 L0/L1 evidence            CLOSURE GATE / OPEN
S1 L2/R2 evidence            CLOSURE GATE / OPEN
D0–D4-B3 reopen              NONE
B4 canonical acceptance      NOT YET
D4 final coherence           BLOCKED until B4 ratification
Implementation               BLOCKED until D9
```

---

## 15. Next bounded proof / access prerequisite

The next review round should **not** re-review the architecture core unless new evidence contradicts it.

Provide the reviewer, outside the repository and without exposing secrets in chat/docs:

- a valid read credential for the currently bound Mercado Livre Installation used by B2;
- the corresponding Mercado Pago/payment account read credential required for Payment/account reads.

Then execute only:

1. **M1** — catalog competition + negative/not-applicable market control;
2. **E1** — `listing_prices` + decorrelation controls + expected seller shipping + one comparable Order/Shipment;
3. **S1-A** — Order→Payment→release/account-impact proof;
4. **S1-B only if S1-A exposes a concrete movement-universe/recovery gap.**

No report generation/export is authorized or required by this candidate.

---

## 16. Proof strategy / strongest counterexamples

B4 must remain correct when:

1. `price_to_win` is unavailable/not-applicable → no fabricated competitor conclusion.
2. Catalog offer traversal is partial → no complete-market claim.
3. `listing_prices` returns 200 while one allegedly material qualifier is absent/ignored → E1 decorrelation catches or leaves the component Unknown.
4. Expected seller shipping cannot be obtained → not zero; E1 remains open for a flow that materially requires it.
5. Current expected quote and historical Order fee differ → preserve L0/L1 context rather than overwrite one with the other.
6. `sale_fee` is per unit and quantity > 1 → no line-total assumption.
7. Buyer shipping charge differs from seller shipping cost → preserve both scopes.
8. Billing charge differs from Order fee because of rebate/bonus/adjustment → keep both evidence classes.
9. Payment is approved but not released → L1 does not become L2 by naming accident.
10. Refund/chargeback occurs after earlier release → append distinct external evidence; prior history remains true.
11. One refund/chargeback is required for Post-Sale closure and Economics attribution → same evidence can feed both owners without authority transfer.
12. One movement has no exact Order anchor → attribution remains unresolved/ambiguous.
13. Account/report population is partial/retention-limited → no completeness claim.
14. Payout/withdrawal exists with no bank integration → provider side known; bank receipt Unknown/R3 open.
15. Current `channel_fees`, `CollectorPort` or env flag is convenient → convenience does not create target authority.
16. A second marketplace/payment provider arrives → extend concrete evidence adapters through existing domain-owned semantics; add shared mechanism only if real repetition proves it reduces total complexity.

---

## 17. Reopen / stop triggers

Reopen only the implicated accepted decision when material evidence shows:

1. Market Intelligence requires a real source that cannot fit source-qualified observation + comparability ownership → targeted D1/D4 review only as implicated.
2. No legitimate market source can satisfy a required Product 1.0 competitive claim and honest insufficiency is no longer product-acceptable → targeted D0/B4 adjudication; no source fabrication.
3. Provider fee/shipping semantics require a materially new business authority absent from D1 → targeted D1 review.
4. Provider monetary evidence cannot fit D2 source-qualified external-movement identity/provenance → targeted D2 review.
5. Financial occurrence/duplicate/partial/recovery semantics cannot fit D3 → targeted D3 review.
6. E1 proves a materially required selected-flow L0 component has no sanctioned sufficient evidence source → **STOP / SPLIT PREREQUISITE**.
7. S1 proves materially required L2/R2 evidence cannot be established through a sanctioned source → **STOP / SPLIT PREREQUISITE**.
8. Product 1.0 requires provider report generation or another B4 write → explicit D3/D4 external-effect adjudication before adding it.
9. A bank integration becomes an accepted source → open R3 bank side in its responsible stage; do not move bank authority into B4.
10. A real second provider proves a repeated non-authority technical mechanism materially reduces total complexity → consider only that mechanism.

Naming preference, current module convenience, hypothetical providers and desire for one unified ledger are not reopen evidence.

---

## 18. Evidence basis for next independent review

Repository evidence considered:

- accepted D0 Product 1.0 Competitive Intelligence / Pricing & Profitability / Economic Evidence chain;
- accepted D1 Market Intelligence, Commercial Economics and Post-Sale Resolution ownership;
- accepted D2 external financial-movement identity, provenance, Economic Attribution and R1/R2/R3 semantics;
- accepted D4-B1/B2/B3 provider-boundary/coverage/no-fabrication rules;
- ADR-009 provenance history;
- reopened ADR-014 / ADR-020 / ADR-032;
- current Mercado Livre adapter FeeQuote, competition, Order/Shipment and legacy fee code as current-state evidence only;
- independent Fable B4 review recorded in `AI-DIALOG.md` on 2026-08-18;
- current official Mercado Livre/Mercado Pago documentation revalidated after review for sale-cost quotation, seller shipping estimation, Payment/release and account-report evidence.

External/provider documentation remains evidence, not architecture authority. Current real-account behavior must still discharge E1/S1 before B4 can be ratified.
