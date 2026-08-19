# D5-B2 — W2-C Readiness + Market Intelligence + Commercial Economics Schema Grammar

> **Status:** ACCEPTED IN-STAGE / OPERATOR-RATIFIED  
> **Parent W2:** `D5-B2-W2-SCHEMA-GRAMMAR.md`  
> **Parent Wire Contract:** `D5-B2-WIRE-CONTRACT.md`  
> **Operation inventory:** `D5-B2-OPERATION-ADMISSION-MATRIX.md`  
> **Parent authorities:** accepted D0→D4 + D4-R1 + D5-B1 + ratified D5-B2 Whole-Matrix + Wire W1 + W2-A/W2-B  
> **Method:** DevelopmentConexus Engineering Method v1.0.0  
> **Accepted:** 2026-08-18

## 1. Role and governing invariant

This subartifact stress-tests the accepted W2-A schema grammar against source readiness evidence, competitive-market interpretation and economic lineage. It owns only the W2-C schema decisions below; it does not create a second W2 authority for W2-A/B.

> **Readiness, Market Intelligence and Commercial Economics expose interpreted owner meaning together with exactly the evidence sufficiency necessary to understand that meaning. Incomplete source/provider evidence never produces a stronger readiness, market or profitability conclusion than the owner can justify.**

No Product mirror, universal Evidence/MarketObservation, mutable Profitability row, finance ledger, generic Reconciliation resource or rules DSL is introduced.

---

# 2. Product & Channel Readiness

## 2.1 Source-product search is evidence, not Product mastery

`SearchSourceProductsForMarketplace` returns source-qualified search hits, not MPC Product resources.

A search hit is anchored by:

```text
SourceProductRef
  = source_instance_id
  + native_product_key
```

and may expose only source facts required by real readiness/correspondence consumers, for example bounded name/SKU/GTIN evidence with honest knowledge/provenance.

No MPC `product_id`, Product lifecycle or generic `attributes` property bag is admitted merely for search convenience.

## 2.2 ProductChannelReadiness is a keyed Q, not a canonical entity

`ProductChannelReadiness` has no synthetic Readiness ID. Its subject is the current combination of:

```text
source-qualified Product
+ Marketplace Installation
```

The response preserves proportionately:

- subject;
- current correspondence meaning;
- current publication readiness conclusion;
- blockers/insufficiency reasons;
- applicable `requirements_revision` where material;
- owner evaluation time/provenance.

A readiness conclusion uses the smallest owner-specific state set needed to distinguish at least:

```text
ready
blocked
unknown
unavailable
```

`blocked` means sufficient knowledge establishes an unmet condition. `unknown` means evidence is insufficient to conclude. They never collapse.

## 2.3 Correspondence state

Correspondence is an owner-specific union, proportionately including meanings such as:

```text
resolved
unresolved
conflicting
unknown
unavailable
```

Resolved correspondence references a Readiness-owned current candidate. Unresolved/conflicting states may expose bounded candidate summaries required by a human/automation resolution consumer.

Candidate keys are opaque and context-bound. They are not provider Product IDs, business-system columns, MPC canonical entities or universal external references.

`ResolveProductChannelCorrespondence` selects a Readiness-owned candidate/current meaning. Product clients never submit source table/column/JSON-path truth as a replacement for Readiness authority.

## 2.4 PublicationRequirements

`PublicationRequirements` is a keyed Q for the current source Product + Marketplace Installation + publication context. It preserves proportionately:

```text
subject
publication_context
requirements_revision
requirements[]
source_media_candidates[]
evaluated_at / material provenance
```

`requirements_revision` is opaque Readiness/provider-requirement meaning used by ListingIntent resolution and is not an MPC resource ETag.

A `PublicationRequirement` uses a Readiness-owned opaque `requirement_key` and a bounded current effective meaning such as:

- effective obligation/applicability (`required`, `recommended`, `optional` or another proven state);
- value kind/cardinality/constraints needed by W2-B `PublicationValue`;
- current source candidates;
- current option set when the requirement is option-based;
- bounded provider/source enrichment only when explanation/correctness requires it.

Raw provider condition/expression language does not become a Product API rules DSL. D4/provider semantics determine applicability; Readiness returns the effective requirement meaning for the current publication context.

## 2.5 Requirement value specs and candidates

Requirement value specs align with the accepted W2-B value family (`text`, `exact_decimal`, `boolean`, `option`, `text_list`, `option_list`) and expose only constraints materially required for safe authoring.

An option uses an opaque requirement-scoped `option_key`; it is not required to equal a raw provider option identifier.

A source candidate uses an opaque `candidate_key` and may expose its current resolved value using a requirement-specific knowledge union. Client-visible source-candidate meaning never leaks business-system table/column/provider JSON-path topology.

Source media candidates follow the same bounded pattern: an opaque Readiness candidate/reference with enough source/provenance meaning to support ListingIntent selection, not a ProductAsset/media master.

---

# 3. Market Intelligence

## 3.1 Market subject

Market Intelligence supports pre-listing and existing-listing reasoning without depending on ListingIntent identity as universal market subject.

A bounded Market subject union is:

```text
existing_listing
  -> source-qualified Listing

source_product_marketplace_context
  -> source-qualified Product
  + Marketplace Installation
  + publication context when material
```

This is Market Intelligence-local schema meaning, not a universal cross-domain SubjectRef.

## 3.2 CompetitivePosition is a contextual Q

`CompetitivePosition` receives no synthetic canonical ID. It preserves proportionately:

```text
subject
coverage
evidence_sufficiency
own-offer meaning where available
best/relevant comparable meaning
competitive relation / gap
change where supported
provider-enriched evidence where justified
evaluated_at / component provenance
```

No generic mutable MarketObservation resource is introduced.

## 3.3 Coverage and evidence sufficiency are distinct

Market `coverage` answers what observed provider/source scope was actually traversed/available. It never claims universal market completeness merely because a provider enumeration completed.

A bounded coverage meaning may distinguish states such as:

```text
complete_for_declared_source_scope
partial
unknown
unavailable
```

The exact declared source/scope remains explicit enough to prevent `complete` from implying the entire market.

`evidence_sufficiency` independently answers whether Market Intelligence can support the specific competitive conclusion from the evidence it has, for example:

```text
sufficient
insufficient
unavailable
```

Partial coverage may still be sufficient for a bounded conclusion, and complete traversal of a weak scope may still be insufficient. These dimensions do not collapse.

## 3.4 Competitive relation is factual/derived, not policy scoring

Baseline competitive relation is a bounded factual comparison such as relative delivered-price position/gap where evidence supports it. W2-C does not invent generic `good/bad/excellent` scoring or commercial thresholds.

A monetary competitive gap uses `Money`. If the evidence does not justify the comparison, the relation is indeterminate/unknown rather than fabricated.

Change (`improved`, `unchanged`, `worsened`, `unknown` or the smallest proven set) requires a materially comparable prior/current basis; absence of such a basis remains unknown.

## 3.5 Comparable offers are evidence values, not MPC entities

`ListComparableOffers` returns bounded comparable-offer evidence values with the source qualification that actually exists, material price/shipping/delivered-price facts and observation time.

No `ComparableOfferId` is minted merely for list semantics. If a provider does not expose a stable native offer identity, MPC does not fabricate one.

Seller/competitor identity or PII is exposed only when a named correctness/operations consumer requires it; provider payload presence alone is not a reason to retain/expose it.

## 3.6 Provider-rich evidence remains bounded

Provider-specific competitive evidence may use a closed owner-local variant, for example a Mercado Livre enrichment carrying material `price_to_win`, competition state, shipping/free-shipping or boost evidence when the accepted consumer needs it.

No `provider_fields` map, raw DTO or cross-provider fabricated equivalence is admitted. Another provider may honestly expose unsupported/not-applicable/unavailable for a distinct enrichment.

---

# 4. Commercial Economics — Expected / Scenario

## 4.1 Economics subject

Expected/scenario economics uses a bounded Economics-local subject union sufficient for pre-listing and existing-listing analysis, for example:

```text
existing_listing
source_product_marketplace_context
```

The shape may resemble Market subject references but does not create a universal shared `CommerceSubject` authority.

## 4.2 ExpectedEconomics is components-first

`ExpectedEconomics` exposes the current economic interpretation as:

```text
subject
components
conclusion
evaluated_at / material provenance
```

Components remain semantically distinct and may include, where applicable:

- sale/candidate price;
- selected Cost Basis;
- expected tax;
- expected marketplace fee;
- expected seller-borne shipping;
- material promotion/discount effects;
- other accepted Economics-owned components only when required by the concrete lineage.

Each component uses an Economics-local knowledge/evidence shape. There is no universal `Fact<Money>` wire type.

A known monetary component preserves `Money` plus bounded provenance such as `observed`, `modeled`, `configured` or `derived` where that distinction is material. Equal numeric values with different evidence kinds remain semantically distinguishable.

## 4.3 Cost Basis remains Economics meaning

Cost Basis output is explainable under Economics-selected policy/basis meaning. Native Sankhya cost columns/variant names do not become Product API economics ontology by convenience.

When material, Cost Basis response preserves the selected basis key/policy revision and source-qualified evidence/provenance sufficient to explain the amount.

## 4.4 Economic conclusion is fail-honest

A complete expected contribution/margin/profitability conclusion exists only when the components required for that conclusion are sufficiently known under the current Economics rules.

A bounded conclusion union distinguishes at least:

```text
known
insufficient_evidence
unavailable
```

where material.

`insufficient_evidence` identifies the missing/insufficient Economics component classes required for the claim. `unavailable` preserves inability to obtain required authoritative evidence.

Missing/unknown tax, fee, shipping or cost never becomes zero/default merely to produce a numeric margin. There is no baseline "best-effort margin" current-truth field.

## 4.5 EvaluatePriceScenario input cannot spoof authority

`EvaluatePriceScenario` is side-effect-free and accepts only legitimate hypothetical scenario variables, for example:

- Economics subject;
- candidate price as exact `Money`;
- quantity or another bounded scenario variable only when materially required;
- selection of an existing accepted Economics basis/policy where the owner contract supports it.

The caller cannot submit authoritative cost/tax/fee/competitor facts as evidence replacements. Those inputs remain source/Market/Economics authority.

A scenario may be contract-valid yet semantically insufficient for a particular calculation. Missing context-dependent scenario input produces an explicit operation-specific insufficient-input result rather than a fabricated default or a universal schema requirement.

## 4.6 PriceScenarioEvaluation

The result reuses Economics meaning:

```text
subject
scenario
components
conclusion
evaluated_at
```

No `SimulationId`, `RecommendationId`, generic Result wrapper or persisted simulation resource is created by default.

Scenario evaluation never actuates price or emits a generic action authority. A later consequential `PriceIntent` remains Offering-owned and Economics preserves/re-establishes the material decision-time L0 basis needed for history/reconciliation.

---

# 5. SaleEconomics / L0-L1-L2 / R1-R2

## 5.1 SaleEconomics is one lineage Q

`SaleEconomics` has no synthetic Economics ID. Its subject is the source-qualified Marketplace Sale and it preserves:

```text
sale
expected_basis      (L0 historical basis)
order_economics     (L1)
realized_economics  (L2)
reconciliation      (R1/R2)
evaluated_at / coverage
```

L0/L1/L2 are distinct semantic rungs, not three mutable resources or one overwritten profitability row.

## 5.2 L0 historical basis

The historical expected basis is distinct from current `ExpectedEconomics`. It preserves the material Economics snapshot/basis relevant to the sale/decision when available.

Its knowledge shape distinguishes a proven absence of a material historical basis from an unknown/unavailable basis where that distinction matters. Later current-policy/source changes do not rewrite L0 history.

## 5.3 L1 Order Economics

Order Economics preserves only the material order-time components/evidence required by the accepted lineage, including sale amount/quantity/discounts/order fee/seller shipping/fiscal evidence as applicable.

Order transaction fee remains Order evidence and is not reclassified into Payment/settlement meaning merely because a later source reports related charges.

Component provenance/source/time remains local to the component where necessary.

## 5.4 L2 realized economics is occurrence-based

Realized Economics preserves material realized/settlement occurrences rather than overwriting prior evidence into one terminal row.

Conceptually:

```text
material_occurrences[]
current_aggregate
evidence_coverage
```

A refund/reversal appends a distinct material occurrence/interpretation and never erases a previously observed release/settlement/payment occurrence.

A `RealizedEconomicOccurrence` is Economics-local lineage evidence, not permission to expose a standalone generic financial-movement/ledger API.

Each occurrence carries only the Economics meaning required for explanation/reconciliation plus the source-qualified native movement reference and material monetary/effective/observation evidence. Provider-native decomposition remains bounded/source-qualified.

## 5.5 Realized aggregate is honest about completeness

The current realized aggregate may itself be partial/insufficient/unknown when the authoritative occurrence universe/coverage is not sufficiently established.

Observed fees/payments do not make realized profitability "final" merely because the values can be summed.

## 5.6 R1 / R2 are assessments, not resources

R1 (`Expected ↔ Order`) and R2 (`Order ↔ Realized`) remain Economics-owned assessment meaning inside the lineage. They receive no generic Reconciliation ID merely for uniformity.

A bounded assessment may distinguish states such as:

```text
reconciled
divergent
pending
insufficient_evidence
```

`pending` means expected later evidence/progression is not yet complete. `insufficient_evidence` means current evidence quality/coverage is not enough to reconcile. These do not collapse.

When divergent, structured variances use bounded Economics component names and exact values. No generic JSON-path/component map becomes the reconciliation ontology.

---

# 6. Economic Attribution

## 6.1 Domain-local subject union

Economic Attribution may use the D2-accepted Economics-local polymorphic subject, but only with explicit variants actually required by Product 1.0, for example:

```text
marketplace_sale
post_sale_resolution
marketplace_installation_scope
```

No universal `{entity_type, entity_id}` graph is introduced.

## 6.2 Attribution states

Attribution preserves the semantic distinction among at least:

```text
exact
partial
ambiguous
unresolved
```

Exact attribution references one bounded Economics subject.

Partial attribution carries explicit allocations plus an explicit unresolved remainder when material. Monetary allocations use `Money`; same-currency/sum invariants are validated against the authoritative attributable movement/evidence.

Ambiguous attribution exposes only bounded candidate subjects sufficient for human resolution.

Unresolved attribution preserves a bounded reason such as lack of sufficient correlation rather than fabricating a match.

## 6.3 ResolveEconomicAttribution

The human-baseline resolve request selects/interprets the current attribution state and subject allocation; it does not rewrite provider movement amount/currency/status/source facts.

The operation remains protected by the current EconomicAttribution `If-Match` safety contract. Client-supplied external evidence cannot impersonate Economics/source authority.

---

# 7. Commercial Policy seam

W2-C freezes only the schema discipline needed for Commercial Policy, not a speculative rules language.

`CommercialPolicy` is a closed typed Economics-owned configuration object containing only admitted policy meanings such as Cost Basis selection/policy, margin floors, price boundaries and economic approval-trigger thresholds when those fields are concretely defined.

`UpdateCommercialPolicyRequest` is a typed partial desired-state update protected by `If-Match`.

No generic `rules[]`, condition/action expression DSL, provider policy map or Governance-owned threshold model is admitted.

The default/inheritance/override/provenance grammar shared by legitimate owner policies remains for the later W2 cross-owner configuration section so Availability/Economics/Fulfillment can be challenged together rather than duplicating a mechanism here.

---

# 8. W2-C negative controls

Later executable contract/conformance proof must make at least these defects invalid or falsifiable:

1. source-product search result gains an MPC Product ID/master lifecycle;
2. source candidates expose table/column/provider JSON paths as Product business vocabulary;
3. Readiness gains a synthetic canonical Readiness ID merely for CRUD symmetry;
4. `blocked` and `unknown` readiness collapse;
5. provider/raw conditional requirement DSL leaks through Product schemas;
6. completed provider enumeration becomes unqualified universal `coverage=complete`;
7. Market Intelligence returns a known competitive relation from insufficient/unavailable evidence;
8. ComparableOffer receives a fabricated MPC identity where no stable source identity exists;
9. competitor/provider PII crosses the boundary without a named correctness/operations need;
10. `provider_fields`/raw provider payload escapes bounded market enrichment;
11. Expected Economics reports known margin/profit after a required tax/fee/shipping/cost component is unknown/unavailable;
12. caller supplies authoritative cost/tax/fee/competitor evidence to `EvaluatePriceScenario`;
13. PriceScenarioEvaluation creates Simulation/Recommendation durable identity or price actuation authority;
14. L0 current policy/source values rewrite historical decision-time basis;
15. refund/reversal overwrites a prior realized occurrence instead of appending lineage;
16. partial settlement evidence is exposed as final realized profitability;
17. R1/R2 receive a generic Reconciliation resource/ID;
18. Economic Attribution uses a universal entity graph or allows human resolve to rewrite authoritative movement facts;
19. partial attribution percentage/rounding silently loses monetary remainder correctness;
20. Commercial Policy accepts a generic rules/condition/action DSL.

---

# 9. W2-C method outcome

**Parent architecture / W1 / W2-A/B:** `CURRENT STRUCTURE CONFIRMED`.

> **Readiness returns bounded source/requirement candidates without becoming Product master; Market Intelligence returns sufficiency-qualified competitive interpretation without pretending universal market coverage; Commercial Economics exposes component-level evidence plus only conclusions justified by those components, preserving L0/L1/L2 and occurrence history without becoming a generic finance/reconciliation platform.**

No parent-stage reopen is required.

---

# 10. Exact next W2 work

**W2-D — Governance + Marketplace Sales + Business-System Materialization + Fulfillment + Post-Sale + Operational Work schema grammar.**

W2-D must stress-test the same schema laws against consequential operational lifecycles and decide proportionately:

1. Authorization Decision and Delegation schemas without `approved=true`, rules-engine or execution authority;
2. source-qualified Marketplace Sale read/attribution resolution without client-created Sale or provider Order mirror;
3. BusinessOrderIntent/InvoicingIntent read/tracking schemas whose creation remains owner-triggered;
4. Party Resolution candidates/resolution state and Destination Realization without Customer/Address master authority;
5. Fulfillment state/checkpoint/Fulfillment Node/artifact/shipment schemas without generic workflow/WMS/TMS/provider DTO mirror;
6. physical-evidence establishment semantics that cannot be forged by an ordinary automation Principal;
7. Post-Sale Resolution scoped consequences without collapsing cancellation/return/refund into one status/provider action vocabulary;
8. Work responsibility/assignment/hold/escalation schema without becoming source truth or a generic Task/Case platform;
9. owner-specific closure/effect/ambiguity states preserving accepted/pending/ambiguous/converged distinctions;
10. negative controls proving client/provider/workflow fields cannot bypass owner/Governance/source truth.

After W2-D, run the remaining cross-owner W2 configuration/problem/outcome consistency work needed for global coherence. **Do not run Fable until W2 converges as one package.**

Implementation remains blocked until D9.
