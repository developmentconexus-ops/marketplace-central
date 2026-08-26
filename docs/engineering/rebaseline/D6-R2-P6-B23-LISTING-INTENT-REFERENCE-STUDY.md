# D6-R2 P6 — B23 ListingIntent Authoring Reference Study

> **Status:** COMPLETE — dispositions adjudicated by the operator 2026-08-26
> **Block:** B23 — ListingIntent authoring / R22–R23
> **Method:** [Frontend Product Experience Planning Method v2.3](../../development/frontend-product-experience-planning-method.md) §12 (P6) + §3.10A
> **Trigger:** the first B23 P8 candidate rendered the authoring surface materially simpler than real marketplace/hub authoring; the operator required the conditional reference study that had been skipped

## 1. Sources

- **SOURCE OBSERVATION — Mercado Livre seller flow** (public seller guidance): category is chosen first and drives the attribute census (wrong category hides attributes); title has a hard 60-character limit that sellers are advised to use fully; the technical sheet (ficha técnica) groups attributes with mandatory ones highlighted and a "show all parameters" disclosure (blank attributes remove the listing from buyer filters); photos are the dominant conversion factor with a primary photo; **variations** (color/size/voltage) live inside one listing with per-option quantities; description is a distinct long-form field.
- **SOURCE OBSERVATION — hub platforms** (AnyMarket public guidance, consistent with the ANYMARKET/Magis5/Hub2b evidence already in the Evidence Register): the technical sheet is a dedicated, marketplace-adapted, grouped editor section with visual mandatory/recommended distinction; title templates compose registered product fields; bulk operations exist at hub scale.
- **REPOSITORY AUTHORITY:** D4-R1 publication-input census; W2 value grammar — `PublicationRequirement` carries `requirement_class` (required/recommended/optional/conditional), `applicability`, `not_applicable_allowed` and one of **seven typed value specs** (text, exact_decimal, boolean, option, text_list, option_list, number_unit); `ListingIntentDesired` carries an operator-selectable `publication_context` (category/product-type refs with human descriptors).

## 2. Dispositions

| Capability surfaced | Disposition |
| --- | --- |
| Category/product-type selection driving the requirement census | **PRESENT-IN-AUTHORITY** — `publication_context` + descriptors; must be rendered in R23 |
| Grouped technical sheet (required first, pendency count, progressive "show all") | **PRESENT-IN-AUTHORITY** — projection of `requirement_class`/`applicability`; structural P8 duty |
| Typed inputs per value kind (option select, number+unit, boolean, text list, …) | **PRESENT-IN-AUTHORITY** — the seven value specs; structural P8 duty |
| Title as constrained text with counter | **PRESENT-IN-AUTHORITY** — text value spec constraints |
| Primary photo emphasis + source/authored media | **PRESENT-IN-AUTHORITY** — media selection model |
| "Não se aplica" where allowed | **PRESENT-IN-AUTHORITY** — `not_applicable_allowed` |
| Title templates composed from registered fields | **REJECTED** — hub-scale authoring convenience; MPC authors per-intent with source-following; no template engine in Product 1.0 |
| Bulk authoring/publishing operations | **REJECTED** — already excluded by accepted D6-R2 laws (no bulk framework) |
| **Listing variations (color/size/voltage, per-option quantities)** | **UPSTREAM FINDING — RATIFIED by the operator 2026-08-26.** Material Mercado Livre capability; current D4-R1/W2/OAD model no variation concept anywhere. Per §3.10A + the P7 blocking law, **B23 P8 is PAUSED** until the bounded upstream repair (smallest owning authority: D4-R1 + D5 W2 + OAD) is designed, adjudicated and integrated. |

## 3. Consequence

- The first B23 P8 candidate (`b23-listing-intents.html`) remains preserved evidence; its post-submit honesty structure (server dispatch gate, external-effect states, ambiguous-verify law) was operator-appreciated and is expected to survive the replan.
- P8 resumes only after the variations upstream finding is ratified into authority; the resumed candidate must also apply every PRESENT-IN-AUTHORITY disposition above.
- Per §5.3 bounded rebaseline: no other locked block (B00…B20) reopens; the repair scope is the ListingIntent authoring seam and the publication-input model it consumes.
