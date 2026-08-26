# D6-R2 P9 — B20 Publicações Screen Contract

> **Status:** DERIVED / PASS — P8 LOCKED 2026-08-26; BACKEND SUFFICIENT; UPSTREAM FINDING NONE
> **Block:** B20 — Publicações core / R20–R21
> **Methods:** [DevelopmentConexus Engineering Method v1.0.0](../../development/engineering-method.md) + [Frontend Product Experience Planning Method v2.3](../../development/frontend-product-experience-planning-method.md)
> **Locked P8 evidence:** `qualification/d6-r2-wireframes/b20-publications.html`
> **Canonical Product OAD:** `contracts/api/product/openapi.yaml`
> **Product implementation:** BLOCKED UNTIL accepted D9

## 1. P9 result

P9 ran only after the operator LOCKed the revision-3 B20 candidate.

The locked human job is:

```text
exact Organization + Marketplace Installation
→ browse/filter observed marketplace listings (name, SKU)
→ open one source-qualified listing detail
→ read the current observation + owner-separated area facts
→ continue by navigation to the responsible owner surface or to ListingIntent authoring
```

The current Product contract is sufficient for this job. The one gap found during P8 (no source-product identity on Listing reads) was repaired inside this increment as the typed `source_product_link` projection. No screen-shaped endpoint, cross-owner aggregate API, generic dashboard service or frontend business authority is required.

**P9 verdict: PASS / BACKEND SUFFICIENT / UPSTREAM FINDING NONE.**

## 2. Route, identity and client-state ownership

Production route family:

`/org/:organizationId/publicacoes` (R20) and `/org/:organizationId/publicacoes/:nativeListingKey` (R21, qualified by exact `marketplace_installation_id`).

| State class | B20 ownership |
| --- | --- |
| `GLOBAL_WORKSPACE_CONTEXT` | `organization_id`; changing Organization invalidates Installation, selection and server state |
| `URL_NAVIGATION_STATE` | `marketplace_installation_id`, list filter, selected `native_listing_key` |
| `SERVER_STATE` | Installation collection, Listing collection/detail, owner-region reads (Availability/Performance/Market/Economics/Work) |
| `LOCAL_EPHEMERAL` | mobile navigation disclosure, technical-details disclosure |

Prototype scenario selectors are Evidence controls, not production state. URL/router state never becomes business authority; TanStack Query remains the server-state owner; no normalized frontend entity mirror.

## 3. Product operation / access binding

| Screen need | Operation | Semantic owner | Permission | Principal kinds | P9 disposition |
| --- | --- | --- | --- | --- | --- |
| populate exact account selector | `ListMarketplaceInstallations` | MarketplacePortfolio | `portfolio.read` | H/A/S | same law as B10; no default account |
| R20 listing collection | `ListMarketplaceListings` | Offering | `offering.read` | H/A/S | cursor collection; coverage state honored |
| R21 listing observation | `GetMarketplaceListing` | Offering | `offering.read` | H/A/S | full observation + `source_product_link` |
| Disponibilidade region | `GetSellableAvailability` | Availability | availability read | H/A/S | read-only region facts |
| Performance region | `GetMarketplaceListingPerformance` | Performance | performance read | H/A/S | read-only region facts |
| Mercado region | `GetCompetitivePosition` | Market | market read | H/A/S | read-only region facts |
| Economia region | `GetExpectedEconomics` | Economics | economics read | H/A/S | read-only region facts |
| downstream authoring | `CreateListingIntentDraft` | Offering | `listing.manage` | H/A | **no B20 mutation home**; navigation-only boundary |

B20 admits **zero write operations**. Human browser realization remains server-side session + CSRF; permission-conditioned visibility is usability only.

## 4. Material screen contract

### 4.1 R20 — listing collection

**OWNER + READ TRUTH:** Offering `ListMarketplaceListings`.
**HUMAN PROJECTION:** `presentation.display_name` (never the canonical key), `source_product_link` (name + SKU + source instance when resolved), observed `lifecycle`, `observed_at`.
**HONESTY LAWS:** known-empty ≠ unknown ≠ unavailable at collection level; per-row presentation and source-link states rendered distinctly; `unresolved` link routes the human to Preparação, never invents a link; collection `coverage` may qualify completeness.
**WIRE MECHANICS:** cursor pagination; filter is client-side narrowing of the consulted page or a bounded server search — never a saved-view platform.
**FORBIDDEN:** bulk-selection framework, synthetic status/score, key-as-label, hidden default Installation.
**BACKEND SUFFICIENCY:** sufficient (after the in-increment `source_product_link` repair).

### 4.2 R21 — listing observation

**OWNER + READ TRUTH:** Offering `GetMarketplaceListing`.
**HUMAN PROJECTION:** display name, lifecycle, publication context (category/product type descriptors), full observed-field census (known/unknown/unavailable/not_applicable), observed media with per-item availability, observed price, provenance with evidence custody, `source_product_link` line.
**WRITE CONTROL:** none.
**FORBIDDEN:** editing controls, per-field sufficiency labels, provider field bags.
**BACKEND SUFFICIENCY:** sufficient.

### 4.3 R21 — owner-separated region facts

**OWNER + READ TRUTH:** each region binds its own owner read (§3); facts render inline, read-only, from already-consulted owner truth.
**HONESTY LAWS:** a degraded region hides its facts and states unknown/unavailable — facts are never faked; unavailable never implies absence.
**AUTHORITY LAW:** composition grants no mutation authority; every `Ver em…` is navigation to the owner surface.
**BACKEND SUFFICIENCY:** sufficient; no cross-owner aggregate endpoint is required or admitted.

### 4.4 Continuation to ListingIntent authoring

**MECHANICS:** navigation-only boundary; no `CreateListingIntentDraft` call from B20; the B23 block owns downstream authoring UX.
**BACKEND SUFFICIENCY:** sufficient for B20.

## 5. Bidirectional trace — PASS

```text
Frontend → Product: account selector → ListMarketplaceInstallations;
list → ListMarketplaceListings; detail → GetMarketplaceListing;
region facts → GetSellableAvailability / GetMarketplaceListingPerformance /
GetCompetitivePosition / GetExpectedEconomics; continuation → navigation only.

Product → Frontend: every admitted read lands on a locked surface;
source_product_link, presentation, lifecycle, observed census, media,
provenance and region facts all have exactly one rendering home.
No admitted B20 operation is orphaned; no locked control lacks an operation.
```

## 6. Adversarial checks

P9 explicitly rejects: canonical key as human label; unavailable collapsing into known-empty; fake facts in degraded regions; region cards gaining writes; B20 creating a ListingIntent by navigation; a screen-shaped aggregate endpoint; frontend-computed listing health/score; reopening Product/OAD merely for rendering convenience. All are excluded by the locked evidence and `verify-d6-r-b20-publications-wireframe.mjs` (11/11 negative controls).

## 7. P9 closure and P10 note

```text
P8 OPERATOR-RATIFIED / LOCKED (2026-08-26, revision 3)
→ exact route/state/identity binding
→ exact owner/operation/Permission binding
→ frontend → backend trace PASS
→ backend → frontend trace PASS
→ adversarial shortcuts rejected
→ BACKEND SUFFICIENT
→ UPSTREAM FINDING NONE
```

**P9: PASS / CLOSED for B20.**

P10: B20 reuses the established laws (B00 exact-context/invalidation, known-empty ≠ unknown ≠ unavailable honesty, navigation-handoff ≠ mutation, typed presentation ≠ canonical ref). The inline owner-region fact card is B20-local until a second locked block proves the same shape; **no new shared component/pattern authority is claimed**. P11, Pre-D9/D9 and Product implementation remain outside this closure.
