# D6-R1 — Marketplace Performance Intelligence

> **Status:** ACCEPTED / CURRENT CONSOLIDATED AUTHORITY  
> **Trigger:** D6 frontend falsifier proved first-party marketplace-performance questions could not be answered without Product authority  
> **Current Product:** 106 operations / 31 ordinary Permissions / H-A-S

## 1. Governing decision

Marketplace Performance Intelligence is a bounded Product 1.0 read/derive authority for the Organization's own marketplace-participation performance.

It answers:

```text
what changed?
where?
with what source/period/coverage/measurement basis?
where should a human investigate?
```

while preserving:

```text
correlation != causation
metric availability != metric equivalence
missing != zero
provider access != permission for every derivative use
```

It does not own Offering/Availability/Sales/Market/Economics/Governance/Work, provider protocol, Ads execution or AI authority.

## 2. Owned meaning

Owns:

- exposure/traffic/engagement evidence interpretation;
- conversion/performance only when source measurement semantics support it;
- retail-media performance interpretation;
- explicit period/comparison meaning;
- evidence sufficiency/coverage/freshness for its claims;
- bounded performance-change interpretation without unsupported causality.

Retail-media **analysis** is in Product 1.0. Campaign/bid/budget/targeting/creative management remains deferred.

The human Strategy/Performance Analyst is an H Principal context, not a new Principal kind or automatic mutation authority.

## 3. Authority edges

Performance may consume current/public meaning from:

```text
Marketplace Portfolio
Offering
Sales
Availability
Market Intelligence
Commercial Economics
```

Primarily Q/current owner meaning; D6 may compose Performance with other owners as read-only presentation. Performance does not mutate another owner or automatically create Work because a metric changed.

## 4. Evidence ownership / history

Performance may preserve the **smallest source-qualified provider evidence required for historical comparison/explanation** when provider retention/recoverability is insufficient.

Persisting provider evidence in MPC does not make it MPC-authored source truth. Keep source, scope, period, observation time, coverage, freshness and measurement basis sufficient for the claim.

No Data Lake/Warehouse/generic Metric store/raw provider archive/event store is admitted. D7 chooses physical retention/indexing only as required by this evidence meaning.

## 5. Mercado Livre proving evidence

Initial admitted evidence remains deliberately narrow:

1. **Visits** — provider listing-traffic evidence with provider-defined period/retention semantics.
2. **Product Ads performance** — provider advertiser/campaign and currently proved scoped performance surfaces.

Retail-media subject scope is closed to proven Product meanings:

```text
campaign
listing
marketplace_catalog_group
marketplace_family_group
```

A grouping identity does not prove measure availability at that scope. D4/D7 acquisition must prove each requested measure for the exact source scope; unsupported/unavailable stays explicit.

`advertiser_id` is provider technical namespace, not Organization/MarketplaceInstallation/SellingEntity/Product identity. Ambiguous advertiser candidates require explicit authorized human selection + current revalidation; never first-result/display-name/ID-order selection.

Advertiser binding is Technical Non-Product configuration under current human + exact Organization Membership + `portfolio.manage` + exact Installation. It adds no Product operation or Permission.

Provider contract/retention/attribution restrictions remain evidence semantics and are preserved honestly.

## 6. Product surface

Exactly four Performance Product Q operations remain admitted:

```text
GetMarketplacePerformanceSummary
ListMarketplaceListingPerformance
GetMarketplaceListingPerformance
ListRetailMediaPerformance
```

All are:

```text
owner       MarketplacePerformanceIntelligence
Permission  performance.read
Principal   H / A / S
scope       Organization + exact MarketplaceInstallation
```

No `performance.manage`, `ads.manage`, generic Analytics/Metric operation or AI permission exists.

Current exact paths are owned by W1/OAD:

```text
.../marketplace-installations/{marketplace_installation_id}/performance
.../marketplace-installations/{marketplace_installation_id}/listing-performances
.../marketplace-installations/{marketplace_installation_id}/listings/{native_listing_key}/performance
.../marketplace-installations/{marketplace_installation_id}/retail-media-performance
```

## 7. Period / comparison semantics

Every Performance read uses explicit primary reporting period `[period_from, period_before)` under the Installation reporting calendar. Optional comparison requires both bounds and owner comparability rules. Frontend presets are URL/navigation convenience, not Product vocabulary.

No all-marketplace aggregate baseline exists; multi-marketplace UX composes distinct source-qualified results without assuming metric equivalence/additivity.

## 8. Schema / knowledge laws

Performance uses owner-local closed schemas, not generic Metric/metadata/result envelopes.

Evidence distinguishes proportionately:

```text
complete
partial
unknown
unavailable
unsupported
```

Known zero is distinct from unknown/unavailable. Available evidence structurally contains the value it claims to know; unavailable variants cannot silently carry plausible values.

Money is exact `Money`; exact rates/multiples are decimal strings; non-negative counts remain precision-safe decimal strings where current schema requires it. Percent/multiple units remain explicit.

Provider measurement basis is explicit enough to avoid assuming two same-named metrics are mathematically/attribution equivalent. Comparison emits change only when periods/evidence/basis are sufficiently comparable; otherwise result is explicit insufficient/not-comparable meaning.

## 9. Collection laws

Performance collections use W3 `limit? + cursor?` plus explicit period/comparison. No baseline:

```text
sort/order_by
metrics/group_by/dimensions/granularity
top
total count
opportunity score
```

Listing Performance population includes currently known Listings even when per-item performance is unknown/unavailable, avoiding survivorship bias; universe completeness remains bounded by Offering Listing coverage.

Retail-media collection distinguishes known-empty population from unavailable/unsupported and preserves the closed subject scopes above.

## 10. Frontend realization

D6 routes:

```text
/performance/resumo
/performance/publicacoes
/performance/midia
```

and R21 may show an owner-separated Performance region. Performance UI must keep explicit Installation/period/coverage/basis and cannot turn provider measurements into causal recommendations or Offering writes.

## 11. Reopen triggers

Reopen only when a real source/consumer requires a new performance meaning, measure/scope, Permission/operation or evidence-retention distinction. Do not generalize from one provider metric into a universal analytics platform, and do not admit Ads mutation merely because Ads performance is readable.
