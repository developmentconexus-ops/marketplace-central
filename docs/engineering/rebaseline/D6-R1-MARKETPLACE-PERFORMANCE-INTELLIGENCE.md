# D6-R1 — Marketplace Performance Intelligence

> **Status:** OPERATOR-APPROVED SEMANTICS / CANONICALIZATION IN PROGRESS
> **Trigger:** D6 frontend falsifier — the accepted 95-operation surface could operate listings, price, availability, sales and economics but could not answer the strategy team's first-party marketplace-performance questions without frontend-created analytics authority.
> **Scope:** bounded amendment to D0/D1/D2/D3/D4/D5 discovered by D6; D7–D9 and Product implementation remain blocked.
> **Method:** DevelopmentConexus Engineering Method v1.0.0
> **Approved through:** D5 W2/W3 plus historical-evidence custody clarification on 2026-08-21.

## 1. Governing decision

Marketplace Central remains **Marketplace Operations Control Plane + Commercial Intelligence** and adds one bounded Product 1.0 read/derive capability:

> **Marketplace Performance Intelligence** observes and preserves sufficiently qualified first-party marketplace performance evidence, interprets own-participation performance over explicit periods, and supports human strategic investigation without becoming authority for Offering, Availability, Sales, Market Intelligence, Commercial Economics, Governance, Work, provider protocol, retail-media execution or AI.

The Product must be able to answer `what changed?`, `where?`, `with what evidence?` and `where should a human investigate?` while preserving the law:

```text
correlation != causation
metric availability != metric equivalence
missing != zero
provider access != contractual permission for every derivative use
```

Retail-media **analysis** is admitted. Campaign/bid/budget/targeting/creative management remains deferred. AI/MCP remains future mechanism and never business authority.

## 2. D0 bounded amendment — capability and actor

Product 1.0 adds:

**Marketplace Performance Intelligence** — observe exposure, traffic, engagement, conversion/performance and retail-media evidence when the source supports those meanings; preserve source definition, period, attribution basis, coverage, freshness and provenance; derive period comparisons and bounded performance interpretation sufficient for human decision support.

Product 1.0 also recognizes the human actor **Marketplace Strategy & Performance Analyst**. The actor investigates performance, compares periods, identifies material changes and supports decisions. The actor gains no mutation/approval authority by role name and remains Principal kind `H`.

Refinement of the existing non-goal:

```text
IN Product 1.0                         OUT / deferred
observe retail-media performance      create campaign
analyze traffic/conversion evidence   change bid or budget
relate performance signals            change targeting/creative
support human investigation           autonomous optimization
```

## 3. D1 bounded amendment — 13th semantic boundary

Add **Marketplace Performance Intelligence**.

**Owns**

- interpretation of the Organization's own marketplace-participation performance;
- exposure/traffic/engagement meaning from admitted evidence;
- conversion/performance meaning only where measurement semantics support it;
- retail-media performance interpretation;
- explicit-period trend/comparison semantics;
- performance evidence sufficiency/coverage/freshness for its claims;
- bounded performance changes/signals that do not claim unsupported causality.

**Does not own**

- Listing representation/lifecycle or PriceIntent — Marketplace Offering Operations;
- Sellable Availability — Availability Control;
- canonical marketplace Sale — Marketplace Sales;
- competitor/comparability meaning — Market Intelligence;
- margin/profitability/economic attribution — Commercial Economics;
- Marketplace Installation lifecycle — Marketplace Portfolio;
- authorization — Controlled Action Governance;
- actionable-work lifecycle — Operational Work;
- provider protocol/DTO/auth — D4;
- campaigns/bids/budgets/targeting/creative — deferred;
- AI/MCP — future mechanism, not a business boundary.

Accepted feed-forward semantic edges into Performance are Portfolio, Offering, Sales, Availability, Market Intelligence and Commercial Economics. Consumption never transfers authority.

`Strategy Workspace` is **not** a 14th domain. It is D6 read composition.

## 4. D2-R2 — Performance Evidence Ownership

Marketplace Performance Intelligence may durably preserve the external evidence required to sustain its own historical interpretation when provider retention/recoverability is insufficient.

Binding law:

> **Persist the smallest source-qualified evidence needed for history, comparison, explanation, audit or a future Product claim; do not persist arbitrary provider payloads merely because they were fetched.**

Persisted evidence remains attributed to its original authority. Storage in MPC does not turn provider-reported CTR/ROAS/visits into MPC-authored facts.

Performance reads may therefore be satisfied from **preserved historical evidence**, not only live provider passthrough. The Product API must expose period/coverage/provenance honestly enough that preserved history and currently re-observable evidence can coexist without changing authority.

No Data Lake, Warehouse, generic Metric store, event store or raw-provider archive is admitted by this amendment. Physical persistence/granularity/retention/indexing remains D7.

## 5. D3 bounded amendment — Q/P baseline

New Performance dependencies use the smallest existing D3 forms:

| Edge | Baseline |
| --- | --- |
| Portfolio → Performance | Q |
| Offering → Performance | Q |
| Sales → Performance | Q |
| Availability → Performance | Q |
| Market Intelligence → Performance | Q |
| Commercial Economics → Performance | Q |
| Strategy Workspace composition | P in D6 |

No baseline Performance `C` or `E` exists.

Reject KPI/event vocabulary such as `CtrDropped`, `RoasChanged`, `LowConversionDetected` merely because a measurement changed. Performance does not automatically create Work or mutate another owner.

## 6. D4 bounded amendment — Mercado Livre evidence contract

Initial admitted Mercado Livre evidence is deliberately small:

1. **Visits** — provider-authoritative listing-traffic evidence under the documented visits surface and its source-defined period/retention semantics.
2. **Product Ads** — current advertiser/campaign/Ad Group performance surfaces, preserving current provider scope and measurement semantics.

Current Product Ads evidence is not forced to Listing scope. Retail-media subject scope remains one of the actually proven provider meanings:

```text
campaign
listing/item
marketplace catalog grouping
marketplace family grouping
```

`advertiser_id` is a D4 technical retail-media namespace and is **not** MarketplaceInstallation, seller identity, Organization identity or a Product resource. A bound advertiser namespace must be established from sufficient provider evidence; ambiguous candidate sets require explicit current-authorized human selection followed by current-authority/candidate revalidation. Never select by first result, display name, browser memory or ID ordering.

The binding ceremony is **Technical Non-Product Ingress/configuration** under a current human with exact Organization Membership + `portfolio.manage` + exact Marketplace Installation. It adds no Product operation or ordinary Permission.

Provider retention/freshness/attribution restrictions are evidence semantics, not Product historical-retention limits. Current Mercado Livre Product Ads documentation establishes a finite historical query window, provider reporting refresh behavior and provider-defined attribution semantics; D2-R2 exists so material evidence can survive provider expiry.

Technical API access does not imply contractual permission to expose/derive every statistic. D4 must fail honest where the intended use is contractually restricted or requires provider authorization.

Listing-quality evidence is not added to the baseline until a concrete applicable source contract is proven for the selected operating categories.

## 7. D5 operation/permission amendment

Add one ordinary Permission:

```text
performance.read
```

No `performance.manage`, `ads.manage` or AI Permission is admitted.

All four Performance operations are Organization + exact MarketplaceInstallation-scoped Qs, allowed for H/A/S and owned by `MarketplacePerformanceIntelligence`.

<!-- performance-operation-matrix:start -->
| Operation | Class | Permission | Allowed Principal |
| --- | --- | --- | --- |
| `GetMarketplacePerformanceSummary` | Q | `performance.read` | H / A / S |
| `ListMarketplaceListingPerformance` | Q | `performance.read` | H / A / S |
| `GetMarketplaceListingPerformance` | Q | `performance.read` | H / A / S |
| `ListRetailMediaPerformance` | Q | `performance.read` | H / A / S |
<!-- performance-operation-matrix:end -->

This bounded amendment changes the canonical Product inventory:

```text
Product operations       95 → 99
ordinary Permissions     29 → 30
Principal kinds          H / A / S unchanged
List/Search Q operations 26 → 28
new C operations          0
new P operations          0
```

## 8. D5 W1 paths

Exact Product paths:

```text
GET /organizations/{organization_id}/marketplace-installations/{marketplace_installation_id}/performance
GET /organizations/{organization_id}/marketplace-installations/{marketplace_installation_id}/listing-performances
GET /organizations/{organization_id}/marketplace-installations/{marketplace_installation_id}/listings/{native_listing_key}/performance
GET /organizations/{organization_id}/marketplace-installations/{marketplace_installation_id}/retail-media-performance
```

Installation nesting is external/business-participation namespace qualification, not Portfolio ownership of Performance.

Every request requires explicit half-open reporting period `[period_from, period_before)` using calendar dates in the Installation reporting calendar. Optional comparison requires both `comparison_from` and `comparison_before`, equal day-count to the primary period and no overlap. Presets such as `last_30_days` remain D6 URL/navigation convenience and are not Product vocabulary.

No all-marketplace aggregate baseline is admitted. A future multi-marketplace workspace composes distinct Installation results rather than treating provider measurements as automatically additive/equivalent.

## 9. D5 W2 schema laws

Performance uses owner-local closed schemas; there is no Product `Metric {name,value}`, `Fact<T>`, `metadata` envelope, analytics DSL or raw provider DTO passthrough.

### Knowledge / coverage

Owner-local performance evidence distinguishes:

```text
complete
partial
unknown
unavailable
unsupported
```

`partial` carries explicit covered periods. `unknown` and `unavailable` use closed reasons. Known zero/empty remains distinct from unknown/unavailable.

Retail-media unavailable reasons include only bounded Product meanings such as configuration required, provider access unavailable, source unavailable or contract restricted. Those are valid 200 response semantics, not ordinary authorization failures.

### Numeric representation

- money uses existing exact `Money`;
- exact rates/multiples use decimal strings;
- counts use a non-negative integer decimal **string** to preserve client precision;
- percentages and multiples carry explicit units so `0.93%` is never confused with `93%`.

### Measurement basis

Provider-reported retail-media evidence carries an explicit provider-reported measurement basis with optional provider-defined attribution window and opaque basis revision where material. Same field naming across providers never implies mathematical/attribution equivalence.

### Comparisons

Comparison interpretation is server-owned Performance meaning. A change is emitted only when current/comparison evidence is sufficiently complete and the measurement basis is comparable. Otherwise the change is explicitly `insufficient_evidence` or `not_comparable`.

### Retail-media scope

Closed scope variants are:

```text
campaign
listing
marketplace_catalog_group
marketplace_family_group
```

No generic `{entity_type, entity_id}`, dimension map or provider AdGroup entity becomes Product ontology.

### Historical evidence

Performance responses are not required to be live passthroughs. Historical values may come from D2-R2-preserved source-qualified evidence; source authority, observation time, coverage and measurement basis remain honest.

## 10. D5 W3 collection laws

`ListMarketplaceListingPerformance` and `ListRetailMediaPerformance` use only the canonical `limit?` + `cursor?` continuation grammar plus the explicit required reporting period and optional comparison period.

No baseline `sort`, `order_by`, `metrics`, `group_by`, `dimensions`, `granularity`, `top`, opportunity score or total count exists.

`ListMarketplaceListingPerformance` enumerates all currently known Listings in the exact Installation, not only Listings that happen to have known performance; per-item evidence can be unknown/unavailable. This avoids survivorship bias. Population completeness remains bounded by Offering Listing acquisition coverage.

Default listing order is source-qualified native Listing key ascending.

`ListRetailMediaPerformance` preserves closed scope grouping/order (`campaign`, `listing`, `marketplace_catalog_group`, `marketplace_family_group`) then native key ascending. Population empty-known and population unavailable/unsupported remain distinguishable.

No time-series/granularity Product contract is admitted yet. Current-period + optional comparable-period is the baseline; a real D6 consumer may later justify a bounded series read.

No `signals[]`, recommendation score or AI explanation schema is admitted yet.

## 11. Explicit rejects / negative controls

The repaired Product contract must reject or make unreachable at least:

1. generic `/analytics`, `/metrics`, `/strategy` Product roots;
2. generic `Metric`, metric selector, `group_by`, dimension or sort DSL;
3. Performance mutation/refresh/collect/sync/optimize Product operation;
4. `performance.read` implying another ordinary Permission;
5. retailer/provider technical access being treated as MPC ordinary Permission;
6. `advertiser_id` becoming MarketplaceInstallation identity;
7. silent advertiser selection among ambiguous candidates;
8. campaign/family/catalog evidence being attributed to one Listing without sufficient evidence;
9. provider CVR/ROAS being reconstructed or redefined by frontend convenience;
10. known zero/empty collapsing into unknown/unavailable;
11. partial period being presented as a complete-period aggregate;
12. comparison across incompatible measurement basis yielding a numeric delta;
13. provider retention becoming MPC historical-retention limit;
14. arbitrary raw provider payload retention becoming baseline architecture;
15. Strategy Workspace becoming Product API/business authority;
16. AI/MCP becoming current Product authority or operation vocabulary.

## 12. D6 consequence

The previously approved shell remains valid but the information architecture must be revised and rendered in Portuguese.

Target strategic group:

```text
ESTRATÉGIA E INTELIGÊNCIA
  Performance
    Resumo
    Publicações
    Mídia
  Mercado
  Economia
```

A Listing gains an operational context and a Performance context without moving Listing ownership. Strategy Workspace remains read composition of Performance + Market + Economics + Offering + Sales + Availability under their respective Permissions.

Performance requires an exact Marketplace Installation. Future `Todas as contas` experiences show per-Installation results rather than dishonest cross-provider metric aggregation.

## 13. Canonicalization / proof plan

1. Preserve the accepted 95/29 verifier byte-for-byte as a baseline proof.
2. Add a wrapper verifier that executes the baseline proof against a baseline projection, then proves the 99/30 repair and new negative controls on the real OAD.
3. Add bounded access views capable of returning the effective 30-Permission vocabulary without altering identity semantics.
4. Add the four Performance paths and closed schemas; add no mutation/P operation.
5. Prove Redocly bundle/lint plus TypeScript and Go generation against the repaired full OAD.
6. Route this repair through `docs/index.md`, `D6-FRONTEND.md` and `roadmap.md`.
7. Re-derive the D6 interaction map and low-fidelity wireframes in Portuguese after OAD proof.
8. Keep PR #54 Draft; do not merge without explicit operator authorization.

## 14. Evidence used for the bounded reopen

Current external evidence was revalidated on 2026-08-21 from official/current sources, including Mercado Livre Visits, current Product Ads Campaign/Ad Group performance documentation and Developer Program Terms. Amazon Sales & Traffic / Brand Analytics and current commerce-intelligence benchmarks were used only as portability/UX evidence, never to fabricate current connectability or provider equivalence.

Repository current authority remains stronger than research/chat. A provider contract change reopens only the smallest affected D4/D5 meaning.
