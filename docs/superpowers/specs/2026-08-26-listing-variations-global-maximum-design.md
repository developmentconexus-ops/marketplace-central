# Listing Variations — Global Maximum Design

> **Status:** DESIGN FOR OPERATOR ADJUDICATION — no authority changed yet
> **Trigger:** ratified B23 blocking upstream finding ([P6 study](../../engineering/rebaseline/D6-R2-P6-B23-LISTING-INTENT-REFERENCE-STUDY.md)): listing variations are a material Mercado Livre capability with no model in current authority
> **Smallest owning authority reopened:** D4-R1 (publication input) + D5 W2 (schema grammar application) + canonical OAD
> **Not reopened:** D0–D3, D5 W1/W3/W4 operation/path/permission surface, D6 locked blocks B00…B20, D7, D8

## 1. The human job (evidence)

A Mercado Livre listing frequently sells one product family in one listing with **variation axes** (Cor, Voltagem, Tamanho) and **options** the buyer picks on the listing page; each option can carry its own photos, some own attributes (e.g. GTIN) and its own quantity. Sellers author this inside the listing editor. Hubs expose the same shape. Source ERPs (Sankhya pattern) usually model each option as its **own native product/SKU**.

## 2. Ownership (no new owner, no synthetic identity)

The governing D4-R1 invariant already decides where each piece lives:

```text
Provider vocabulary (which axes/options a category admits)
  → D4 evidence, served through the existing Readiness-owned
    publication-requirements port (census extension)

Which fields are answered per listing vs per option
  → Readiness census (requirement scope)

Desired variation structure (axes chosen, options authored,
per-option resolutions/media, per-option source product)
  → Offering ListingIntent desired state

Per-option quantity                → Sellable Availability (per source product/SKU);
                                     composed by D4 at dispatch under the accepted
                                     R1-G1 composition law — never authored in B23
Per-option price                   → PriceIntent, same composition law
Observed variations on the live listing → Offering MarketplaceListing observation (read-only)
```

Option identity is **provider axis/option canonical keys + optional source-qualified product ref**. MPC invents no variant master and no synthetic variant id.

## 3. Wire shape (W2 grammar, all typed / additionalProperties false / canonical keys)

### 3.1 Census extension (Readiness / GetPublicationRequirements)

- `PublicationRequirement` gains `scope: listing | per_variation` — which requirements are answered once per listing and which per option (provider evidence decides).
- New `VariationAxisSpec[]` on the requirements read: `axis_key`, `display_name`, `axis_class: required | optional`, `options[] {option_key, display_name}` (provider vocabulary; open text axes carry a text value spec instead of options).

### 3.2 ListingIntent desired extension (Offering)

`ListingIntentDesired` gains an optional `variations`:

```text
variations
  axes[]     axis_key                       (chosen axes, canonical)
  options[]  option_coordinates[] {axis_key, option_key | text_value}
             source_product?      SourceProductRef        (SKU-level origin when the ERP models options as SKUs)
             requirement_resolutions[]                    (only per_variation-scoped requirements; follow_source / explicit_override, unchanged kinds)
             media_selection[]                            (source_media / authored_media, unchanged kinds)
```

Absent `variations` = non-variation listing (today's shape, unchanged). Dispatch blockers may now reference option coordinates, so "falta foto na opção Cor: Inox" is server truth like any other blocker.

### 3.3 Observed listing extension (Offering read, feeds B20/B23)

`MarketplaceListing` gains optional `observed_variations[]`: option coordinates + presentation (display label) + per-option observed fields/price/media states, same honest known/unknown/unavailable grammar.

## 4. Explicitly rejected / deferred

- **REJECTED:** per-option price/stock authoring inside B23 (owners already exist; composition at dispatch is already law); generic variant-matrix engine; synthetic MPC variant identity; auto-generation of options from source without explicit human authoring.
- **DEFERRED (Product reason):** kits/bundles (distinct provider concept, separate future finding); variation-aware Performance/Economics splits (evidence not yet proven material — reopen trigger recorded).

## 5. Surface impact

Zero new operations, paths, Permissions or Principal kinds — **106/31/H-A-S preserved**. The repair is schema-shaped only, protected by extending `verify-human-operable-read-projection.mjs` (or a focused sibling verifier) with the new typed shapes and negative controls, wired into the diff-aware gate.

## 6. Frontend consequence (B23 resume scope)

The resumed B23 editor adds one region: **Variações** — choose axes (provider vocabulary), author options, per-option photos and per-variation fields, with per-option blockers surfacing in the dispatch box. R22 list rows show a variation count when present. B20 detail may later render observed variations from the same read (bounded revalidation, only if the operator asks).

## 7. Proof obligations

- OAD lint + operation census unchanged (106/31), auth profile proof unchanged;
- projection verifier: new negative controls (per_variation scope cannot collapse into listing scope; option identity is keys-only; desired variations never carry price/quantity; observed variations honest-state grammar);
- B23 wireframe verifier extended when P8 resumes.
