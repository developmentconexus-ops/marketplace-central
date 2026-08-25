# Human-Operable Read Projection & Wire Conformance Implementation Plan

> **For agentic workers:** Follow the repository-local Engineering Method v1.0.0 and Frontend Product Experience Planning Method v2.3. Execute this plan task-by-task only after explicit operator approval. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Rehome the approved Global Maximum into canonical D4/D5 authority and repair the Product OAD so the proven human-facing Readiness/Offering flows are operable without changing canonical identity, operation count, Permission count, Principal kinds, or runtime topology.

**Architecture:** Canonical refs and write/decision carriers stay key/ID based. Current human reads gain only owner-correct presentation projections; purpose/historical labels remain snapshots, not current truth. The repair also restores directly implicated W2→OAD semantics and consolidates MarketplaceListing presentation already duplicated by Performance. The prerequisite stops when its wire is proved and integration-ready; B10 correspondence P8/P9 revalidation is a later bounded frontend increment after this prerequisite lands.

**Tech Stack:** OpenAPI 3.1.2 YAML; Node.js `>=26.3.0 <27`; PowerShell aggregate gate; `@redocly/cli@2.45.0`; `openapi-typescript@7.13.0`; TypeScript `5.9.3`; Go `1.25.1`; `oapi-codegen v2.8.0`; `github.com/oapi-codegen/runtime v1.7.0`.

**Spec:** `docs/superpowers/specs/2026-08-25-human-operable-read-projection-design.md`

**Baseline:** `main@181f606ceaf5fadd7b25aab2008d0256ed6ad7de` — PR #71 integrated; current authority owners and diff-aware aggregate gate active.

## Global Constraints

- Preserve **106 Product operations / 31 ordinary Permissions / Principal kinds H / A / S**. A count change is a new material finding and stops this plan.
- Preserve `MarketplaceListingRef`, `SourceProductRef`, ListingIntent/PriceIntent IDs, option/unit/source-candidate/correspondence keys, and other canonical decision carriers as identity/decision authority.
- Presentation never authenticates, authorizes, scopes, correlates, matches, converges, or replaces canonical identity.
- Request/write schemas never accept human labels, presentation locators, server attribution, or source/provider provenance as write authority.
- Preserve `known != missing != unknown != unavailable != unsupported`; known empty remains distinct from unknown/unavailable.
- No `PresentationService`, generic `EntityRef`/`EntityPresentation`, metadata bag, provider field bag, Product/PIM master, generic media/asset service, transform/rule engine, or new Product search/list operation.
- No N+1 point-GET production baseline for a collection whose admitted human job requires scan/select/navigation.
- Source/authored/observed media presentation remain distinct trust schema families.
- Runtime remains **NONE**; implementation paths under `apps|cmd|internal|server|backend|frontend|src|migrations` remain blocked until D9.
- PR #69 / B20 remains **PAUSED / NO P8** throughout this plan.
- B10 main structure remains protected. This plan does not edit B10 HTML or rerun P9.
- CI remains one aggregate `npm run gate`; add one focused semantic verifier under the existing Product OAD proof, not another workflow/check and not prose-string ratification tests.
- Merge remains separately operator-authorized.

---

## File Map

**Canonical authority**
- `docs/engineering/rebaseline/D4-EXTERNAL-INTEGRATIONS.md`
- `docs/engineering/rebaseline/D4-R1-PUBLICATION-INPUT.md`
- `docs/engineering/rebaseline/D5-B2-W2-SCHEMA-GRAMMAR.md`
- `docs/engineering/rebaseline/D5-B2-W3-COLLECTION-GRAMMAR.md`

**Canonical Product wire**
- `contracts/api/product/components.yaml`
- `contracts/api/product/paths-identity-portfolio-readiness.yaml`
- `contracts/api/product/paths-offering-availability-market.yaml`
- `contracts/api/product/paths-performance.yaml`
- `contracts/api/product/paths-economics-governance-sales-materialization.yaml` only when an existing response description/binding must reflect the repaired schema
- `contracts/api/product/openapi.yaml` only if an existing local `$ref` needs redirection; no new path

**Proof**
- Create `scripts/verify-human-operable-read-projection.mjs`
- Modify `scripts/verify-product-oad.mjs`
- Modify `scripts/gate.ps1` only to register `scripts/verify-human-operable-read-projection.mjs` as a Product-proof input in the existing diff-aware affected-path predicate; do not add another workflow/check

**Status/design**
- `docs/superpowers/specs/2026-08-25-human-operable-read-projection-design.md`
- `docs/roadmap.md`

---

### Task 1: Rehome the approved decision into canonical D4/W2/W3 authority

**Files:**
- Modify `docs/engineering/rebaseline/D4-EXTERNAL-INTEGRATIONS.md`
- Modify `docs/engineering/rebaseline/D4-R1-PUBLICATION-INPUT.md`
- Modify `docs/engineering/rebaseline/D5-B2-W2-SCHEMA-GRAMMAR.md`
- Modify `docs/engineering/rebaseline/D5-B2-W3-COLLECTION-GRAMMAR.md`
- Modify `docs/superpowers/specs/2026-08-25-human-operable-read-projection-design.md`
- Modify `docs/roadmap.md`

**Interfaces:**
- Consumes: approved design `Canonical Ref != Current Read Projection != Purpose/Historical Snapshot`.
- Produces: canonical authority for Tasks 2–5; no Product operation/path/runtime change.

- [ ] **Step 1: Add D4's human-operable external-evidence fence**

Add this normative meaning under the general translation fence and Mercado Livre Listing/publication evidence:

```markdown
### Human-operable external presentation evidence

When a consumer-owned Product read has a proven human recognition/selection job, D4 preserves the smallest current source/provider presentation evidence needed by that consumer in addition to the canonical external key. Provider/source presentation is mutable, non-unique evidence and never MPC identity, correspondence authority, authorization, or a generic metadata bag.

For Mercado Livre publication/listing evidence this includes, when applicable, the human names/titles associated with already-used category/product-type/attribute/allowed-value/unit/Listing identities. Adapter DTO/field topology remains private.

When presentation cannot currently be established, D4 propagates honest unknown/unavailable presentation rather than fabricating a label or promoting the native key into a name.
```

- [ ] **Step 2: Extend D4-R1's publication seam**

Record exactly:

```text
canonical decision identity
  requirement_key / option_key / unit_key / source_candidate_key / correspondence candidate_key

current human read projection
  source/provider display presentation needed to recognize the choice

write/effect
  canonical key only + current owner revalidation
```

Require Readiness to preserve human presentation for SourceProduct/SourceInstance, requirements, options, units, FOLLOW_SOURCE candidates, and correspondence candidates when human choice is admitted. Continue rejecting provider expressions/raw paths/arbitrary maps.

- [ ] **Step 3: Add W2's read-projection grammar and exact approved schema-family names**

Add after the existing request-vs-read rule:

```markdown
### Human-operable read projection grammar

A canonical Ref/request carrier remains minimal. A current read may carry an adjacent owner-correct presentation projection when a proven human job requires recognition, selection or explanation. A purpose/historical presentation snapshot is a third meaning and is not refreshed as current source truth.

Presentation is never accepted as identity or write authority. Equal labels never collapse distinct keys. Unknown/unavailable presentation never erases the known canonical subject.
```

Crystallize these schema-family names as the implementation vocabulary:

```text
SourceProductPresentation
PublicationContextRef / PublicationContextView
PublicationOptionDescriptor
PublicationUnitDescriptor
PublicationValue / PublicationValueView
PublicationSourceCandidateView
CorrespondenceCandidatePopulation
MarketplaceListingPresentation
ListingIntent requirement-resolution read views
source/authored/observed media presentation families
```

Also record that MarketplaceListing/ListingIntent OAD drift directly implicated by the approved spec is repaired under W2 without changing W1/W4.

- [ ] **Step 4: Update W3's ListItem law and current SearchSourceProducts semantics**

Add:

```markdown
When the admitted human consumer must scan/select/navigate a member and the owner can supply current presentation without a second business conclusion, the ListItem carries that owner-semantic presentation subset directly. Point-GET fan-out is not the baseline repair for a deficient collection item. This admits no generic `View<T>`, projection DSL, total count, alternate sort, or metadata envelope.
```

Apply only to the currently proven collections:

```text
SearchSourceProductsForMarketplace
ListMarketplaceListings
ListListingIntents
ListPriceIntents
ListSellableAvailability
ListCompetitivePositions
ListExpectedEconomics
ListMarketplaceListingPerformance
```

Also reconcile W3's SearchSourceProducts row with the already-accepted current contract: `source_instance_id` is optional narrowing; omission searches admitted/configured Organization-scoped sources and never selects a hidden/default source. Preserve SourceProductRef as member identity/tie-breaker.

- [ ] **Step 5: Set execution status now that operator plan approval exists**

Execution of this task occurs only after the operator has approved this plan. Therefore set:

```text
spec: DESIGN APPROVED / IMPLEMENTATION PLAN APPROVED
roadmap prerequisite: EXECUTION ACTIVE
roadmap exact next action: execute the approved prerequisite plan task-by-task; B20 remains paused
```

Do not authorize runtime Product implementation or B10/B20 HTML.

- [ ] **Step 6: Verify the authority-only diff**

```powershell
git diff --check main...HEAD
git diff --name-only main...HEAD
```

Expected: no whitespace/conflict-marker error and no OAD/runtime file changed by Task 1.

- [ ] **Step 7: Commit**

```bash
git add docs/engineering/rebaseline/D4-EXTERNAL-INTEGRATIONS.md \
        docs/engineering/rebaseline/D4-R1-PUBLICATION-INPUT.md \
        docs/engineering/rebaseline/D5-B2-W2-SCHEMA-GRAMMAR.md \
        docs/engineering/rebaseline/D5-B2-W3-COLLECTION-GRAMMAR.md \
        docs/superpowers/specs/2026-08-25-human-operable-read-projection-design.md \
        docs/roadmap.md
git commit -m "docs(d6-r2): rehome human-operable read projection authority"
```

---

### Task 2: Repair Readiness/source/publication/correspondence projections and install one focused semantic verifier

**Files:**
- Create `scripts/verify-human-operable-read-projection.mjs`
- Modify `scripts/verify-product-oad.mjs`
- Modify `scripts/gate.ps1` only for the focused-verifier affected-path trigger
- Modify `contracts/api/product/components.yaml`
- Modify `contracts/api/product/paths-identity-portfolio-readiness.yaml` only for descriptions that materially clarify the repaired response

**Interfaces:**
- Produces `SourceProductPresentation`, publication read descriptors/value views, typed source candidates, correspondence candidate population, and source-media presentation.
- Keeps all existing Readiness operation IDs, paths, Permissions, Principal kinds and write request identities.

- [ ] **Step 1: Create a failing semantic verifier that can grow across Tasks 2–5**

Create `scripts/verify-human-operable-read-projection.mjs`:

```js
import { readFileSync } from 'node:fs';

const bundlePath = process.argv[2];
if (!bundlePath) throw new Error('usage: node scripts/verify-human-operable-read-projection.mjs <resolved-bundle.json>');
const document = JSON.parse(readFileSync(bundlePath, 'utf8'));
let negativeControls = 0;

function fail(message) { throw new Error(message); }
function assert(condition, message) { if (!condition) fail(message); }
function schemas(doc) { return doc.components?.schemas ?? {}; }
function requireFieldsFrom(s, name, fields) {
  assert(s[name], `missing schema ${name}`);
  const required = new Set(s[name].required ?? []);
  for (const field of fields) assert(required.has(field), `${name} must require ${field}`);
}

function validateReadiness(doc) {
  const s = schemas(doc);
  for (const name of [
    'SourceProductPresentationKnown','SourceProductPresentationUnknown','SourceProductPresentationUnavailable','SourceProductPresentation',
    'PublicationCategoryDescriptor','PublicationProductTypeDescriptor','PublicationContextRef','PublicationContextView',
    'PublicationOptionDescriptor','PublicationUnitDescriptor','PublicationValueView','PublicationSourceCandidateView',
    'ProductChannelCorrespondenceCandidate','CorrespondenceCandidatePopulationKnown','CorrespondenceCandidatePopulationUnknown','CorrespondenceCandidatePopulationUnavailable','CorrespondenceCandidatePopulation',
    'SourceMediaPresentationKnown','SourceMediaPresentationUnavailable','SourceMediaPresentation'
  ]) assert(s[name], `missing schema ${name}`);
  requireFieldsFrom(s, 'SourceProductSearchHit', ['source_product','presentation']);
  requireFieldsFrom(s, 'ProductChannelReadiness', ['subject','subject_presentation','correspondence','correspondence_candidate_population','correspondence_etag','readiness','blockers','evaluated_at']);
  requireFieldsFrom(s, 'PublicationRequirements', ['subject','subject_presentation','publication_context','requirements_revision','requirements','source_media_candidates','evaluated_at']);
  requireFieldsFrom(s, 'PublicationRequirement', ['requirement_key','display_name','requirement_class','applicability','value_spec','not_applicable_allowed','source_evidence']);
  assert(JSON.stringify(s.PublicationOptionRequirementSpec).includes('PublicationOptionDescriptor'), 'option requirement must expose descriptors');
  assert(JSON.stringify(s.PublicationOptionListRequirementSpec).includes('PublicationOptionDescriptor'), 'option-list requirement must expose descriptors');
  assert(JSON.stringify(s.PublicationNumberUnitRequirementSpec).includes('PublicationUnitDescriptor'), 'number-unit requirement must expose unit descriptors');
  assert(JSON.stringify(s.PublicationSourceEvidenceKnown).includes('PublicationSourceCandidateView'), 'known source evidence must expose candidate views');
  assert(JSON.stringify(s.PublicationSourceEvidenceConflicting).includes('PublicationSourceCandidateView'), 'conflicting source evidence must expose candidate views');
  assert(JSON.stringify(s.SourceMediaCandidate).includes('SourceMediaPresentation'), 'source media must use its own presentation trust type');
  for (const name of ['ResolveCorrespondenceRequest','ClearCorrespondenceRequest','ListingIntentDesired','RequirementResolution','ExplicitOverrideResolution','PublicationValue','MediaSelection','CreatePriceIntentRequest']) {
    const text = JSON.stringify(s[name] ?? {});
    for (const forbidden of ['display_name','display_label','access_ref','subject_presentation']) assert(!text.includes(`\"${forbidden}\"`), `${name} must not author ${forbidden}`);
  }
}

function validateAll(doc) {
  validateReadiness(doc);
  if (typeof validateMarketplaceListing === 'function') validateMarketplaceListing(doc);
  if (typeof validateConsumers === 'function') validateConsumers(doc);
  if (typeof validateListingIntent === 'function') validateListingIntent(doc);
}

function expectMutationFailure(label, mutate) {
  const candidate = structuredClone(document);
  let failed = false;
  try { mutate(candidate); validateAll(candidate); } catch { failed = true; }
  assert(failed, `negative control unexpectedly passed: ${label}`);
  negativeControls++;
}

validateAll(document);
expectMutationFailure('requirement label removed', (d) => { d.components.schemas.PublicationRequirement.required = d.components.schemas.PublicationRequirement.required.filter((x) => x !== 'display_name'); });
expectMutationFailure('correspondence candidate population removed', (d) => { d.components.schemas.ProductChannelReadiness.required = d.components.schemas.ProductChannelReadiness.required.filter((x) => x !== 'correspondence_candidate_population'); });
expectMutationFailure('resolve write accepts label', (d) => { d.components.schemas.ResolveCorrespondenceRequest.properties.display_label = { type: 'string' }; });
assert(negativeControls === 3, `negative-control count must be 3, found ${negativeControls}`);
console.log('human_operable_read_projection=PASS');
console.log(`human_operable_read_projection_negative_controls=${negativeControls}/3`);
```

Later tasks add validators and negative controls, and update the final exact count; `expectMutationFailure` always calls `validateAll`, so a mutation in MarketplaceListing/consumer/ListingIntent proof cannot accidentally be checked only against Readiness.

- [ ] **Step 2: Prove RED against the pre-repair OAD**

```powershell
npx --yes @redocly/cli@2.45.0 bundle contracts/api/product/openapi.yaml --config contracts/api/product/redocly.yaml -o .human-projection-test.json
node scripts/verify-human-operable-read-projection.mjs .human-projection-test.json
Remove-Item .human-projection-test.json
```

Expected: FAIL at the first missing approved schema, e.g. `SourceProductPresentationKnown`. A pass is a test defect.

- [ ] **Step 3: Implement SourceProduct presentation and exact-subject presentation**

Add:

```yaml
  SourceProductPresentationKnown:
    type: object
    additionalProperties: false
    required: [state, display_name, source_instance_display_name]
    properties:
      state: {const: known}
      display_name: {type: string, minLength: 1}
      source_instance_display_name: {type: string, minLength: 1}
      sku: {type: string}
      gtin: {type: string}
      observed_at: {$ref: '#/schemas/Instant'}
  SourceProductPresentationUnknown: {type: object, additionalProperties: false, required: [state], properties: {state: {const: unknown}}}
  SourceProductPresentationUnavailable: {type: object, additionalProperties: false, required: [state], properties: {state: {const: unavailable}}}
  SourceProductPresentation:
    oneOf:
      - {$ref: '#/schemas/SourceProductPresentationKnown'}
      - {$ref: '#/schemas/SourceProductPresentationUnknown'}
      - {$ref: '#/schemas/SourceProductPresentationUnavailable'}
  SourceProductSearchHit:
    type: object
    additionalProperties: false
    required: [source_product, presentation]
    properties:
      source_product: {$ref: '#/schemas/SourceProductRef'}
      presentation: {$ref: '#/schemas/SourceProductPresentationKnown'}
```

Require `subject_presentation: SourceProductPresentation` on both `ProductChannelReadiness` and `PublicationRequirements`.

- [ ] **Step 4: Split publication context and dynamic option/unit read vocabulary from key-only writes**

Add exact read/write schemas:

```yaml
  PublicationContextRef:
    type: object
    additionalProperties: false
    properties:
      category_key: {$ref: '#/schemas/OpaqueKey'}
      product_type_key: {$ref: '#/schemas/OpaqueKey'}
  PublicationCategoryDescriptor:
    type: object
    additionalProperties: false
    required: [category_key, display_name]
    properties: {category_key: {$ref: '#/schemas/OpaqueKey'}, display_name: {type: string, minLength: 1}}
  PublicationProductTypeDescriptor:
    type: object
    additionalProperties: false
    required: [product_type_key, display_name]
    properties: {product_type_key: {$ref: '#/schemas/OpaqueKey'}, display_name: {type: string, minLength: 1}}
  PublicationContextView:
    type: object
    additionalProperties: false
    properties:
      category: {$ref: '#/schemas/PublicationCategoryDescriptor'}
      product_type: {$ref: '#/schemas/PublicationProductTypeDescriptor'}
  PublicationOptionDescriptor:
    type: object
    additionalProperties: false
    required: [option_key, display_name]
    properties: {option_key: {$ref: '#/schemas/OpaqueKey'}, display_name: {type: string, minLength: 1}}
  PublicationUnitDescriptor:
    type: object
    additionalProperties: false
    required: [unit_key, display_name]
    properties: {unit_key: {$ref: '#/schemas/OpaqueKey'}, display_name: {type: string, minLength: 1}}
```

`PublicationRequirements.publication_context` references `PublicationContextView`. Query parameters remain key-only. Option and option-list requirement specs use `options: PublicationOptionDescriptor[]`; number-unit spec uses `units: PublicationUnitDescriptor[]` and keeps `default_unit_key` canonical.

- [ ] **Step 5: Add `PublicationValueView` and typed source candidates while preserving type-family validation**

Add:

```yaml
  PublicationOptionValueView:
    type: object
    additionalProperties: false
    required: [kind, option_key, display_name]
    properties: {kind: {const: option}, option_key: {$ref: '#/schemas/OpaqueKey'}, display_name: {type: string, minLength: 1}}
  PublicationOptionListValueView:
    type: object
    additionalProperties: false
    required: [kind, options]
    properties: {kind: {const: option_list}, options: {type: array, items: {$ref: '#/schemas/PublicationOptionDescriptor'}}
  PublicationNumberUnitValueView:
    type: object
    additionalProperties: false
    required: [kind, value, unit_key, unit_display_name]
    properties:
      kind: {const: number_unit}
      value: {$ref: '#/schemas/ExactDecimalString'}
      unit_key: {$ref: '#/schemas/OpaqueKey'}
      unit_display_name: {type: string, minLength: 1}
  PublicationValueView:
    oneOf:
      - {$ref: '#/schemas/PublicationTextValue'}
      - {$ref: '#/schemas/PublicationExactDecimalValue'}
      - {$ref: '#/schemas/PublicationBooleanValue'}
      - {$ref: '#/schemas/PublicationOptionValueView'}
      - {$ref: '#/schemas/PublicationTextListValue'}
      - {$ref: '#/schemas/PublicationOptionListValueView'}
      - {$ref: '#/schemas/PublicationNumberUnitValueView'}
      - {$ref: '#/schemas/PublicationNotApplicableValue'}
  PublicationSourceCandidateView:
    type: object
    additionalProperties: false
    required: [source_candidate_key, display_label, value]
    properties:
      source_candidate_key: {$ref: '#/schemas/OpaqueKey'}
      display_label: {type: string, minLength: 1}
      value: {$ref: '#/schemas/PublicationValueView'}
```

Change known source evidence to `candidates: array`, `minItems: 1`; conflicting to `array`, `minItems: 2`; items reference `PublicationSourceCandidateView`.

**Also update all existing type-specific source-evidence wrappers** (`PublicationTextSourceEvidence`, `PublicationExactDecimalSourceEvidence`, `PublicationBooleanSourceEvidence`, `PublicationOptionSourceEvidence`, `PublicationTextListSourceEvidence`, `PublicationOptionListSourceEvidence`, `PublicationNumberUnitSourceEvidence`) so their `known|conflicting` candidate arrays narrow each candidate's `value` to the matching value/view family. Do not leave the old object-map constraints in those wrappers.

Make `PublicationRequirement.display_name` required/non-empty.

- [ ] **Step 6: Add correspondence candidate population without a new operation**

```yaml
  ProductChannelCorrespondenceCandidate:
    type: object
    additionalProperties: false
    required: [candidate_key, display_label]
    properties: {candidate_key: {$ref: '#/schemas/OpaqueKey'}, display_label: {type: string, minLength: 1}}
  CorrespondenceCandidatePopulationKnown:
    type: object
    additionalProperties: false
    required: [state, candidates]
    properties:
      state: {const: known}
      candidates: {type: array, items: {$ref: '#/schemas/ProductChannelCorrespondenceCandidate'}}
  CorrespondenceCandidatePopulationUnknown: {type: object, additionalProperties: false, required: [state], properties: {state: {const: unknown}}}
  CorrespondenceCandidatePopulationUnavailable: {type: object, additionalProperties: false, required: [state], properties: {state: {const: unavailable}}}
  CorrespondenceCandidatePopulation:
    oneOf:
      - {$ref: '#/schemas/CorrespondenceCandidatePopulationKnown'}
      - {$ref: '#/schemas/CorrespondenceCandidatePopulationUnknown'}
      - {$ref: '#/schemas/CorrespondenceCandidatePopulationUnavailable'}
```

Require `correspondence_candidate_population` on `ProductChannelReadiness`. Keep `ResolveCorrespondenceRequest = subject + correspondence_etag + candidate_key`; no label is submitted.

- [ ] **Step 7: Add source-media presentation as its own trust type**

```yaml
  SourceMediaPresentationKnown:
    type: object
    additionalProperties: false
    required: [state, access_ref]
    properties: {state: {const: known}, access_ref: {type: string, format: uri-reference, minLength: 1}}
  SourceMediaPresentationUnavailable: {type: object, additionalProperties: false, required: [state], properties: {state: {const: unavailable}}}
  SourceMediaPresentation:
    oneOf:
      - {$ref: '#/schemas/SourceMediaPresentationKnown'}
      - {$ref: '#/schemas/SourceMediaPresentationUnavailable'}
```

Require `presentation: SourceMediaPresentation` on `SourceMediaCandidate`. Do not put `access_ref` in `MediaSelection`, upload/write schemas, or history.

- [ ] **Step 8: Prove GREEN and hook the focused verifier into the existing Product proof**

Bundle and run the focused verifier; expected PASS with `3/3` negative controls.

In `scripts/verify-product-oad.mjs` add:

```js
const humanOperableReadProjectionVerifier = join(root, 'scripts/verify-human-operable-read-projection.mjs');
```

and after `currentProjectionProof(bundleA);`:

```js
run(process.execPath, [humanOperableReadProjectionVerifier, bundleA]);
```

In `scripts/gate.ps1`, add the new verifier to the existing `$productProofPatterns` array:

```powershell
'^scripts/verify-human-operable-read-projection\.mjs$'
```

This is required because PR #71 made Product proof diff-aware. A change to the focused verifier itself is a Product-proof-input change and must not be classified as unrelated documentation/tooling.

Run:

```powershell
node scripts/verify-product-oad.mjs
```

Expected: existing historical/auth/generated proofs PASS, `human_operable_read_projection=PASS`, Product count `106/106`.

- [ ] **Step 9: Commit**

```bash
git add contracts/api/product/components.yaml \
        contracts/api/product/paths-identity-portfolio-readiness.yaml \
        scripts/verify-human-operable-read-projection.mjs \
        scripts/verify-product-oad.mjs \
        scripts/gate.ps1
git commit -m "feat(d5): make readiness reads human-operable"
```

---

### Task 3: Repair MarketplaceListing actual-state and consolidate Performance Listing presentation

**Files:**
- Modify `contracts/api/product/components.yaml`
- Modify `contracts/api/product/paths-performance.yaml`
- Modify `contracts/api/product/paths-offering-availability-market.yaml` only when response description needs clarification
- Modify `scripts/verify-human-operable-read-projection.mjs`

**Interfaces:**
- Consumes `PublicationContextView`, `PublicationValueView`, `MarketplaceListingRef`.
- Produces one `MarketplaceListingPresentation` meaning reused by Offering/Performance plus W2-compliant Listing actual-state presentation/media/provenance axes.

- [ ] **Step 1: Extend verifier RED**

Add:

```js
function validateMarketplaceListing(doc) {
  const s = schemas(doc);
  for (const name of [
    'MarketplaceListingPresentationKnown','MarketplaceListingPresentationUnknown','MarketplaceListingPresentationUnavailable','MarketplaceListingPresentation',
    'ListingObservedFieldKnown','ListingObservedFieldUnknown','ListingObservedFieldUnavailable','ListingObservedFieldNotApplicable','ListingObservedField',
    'MarketplaceListingMediaPresentationKnown','MarketplaceListingMediaPresentationUnavailable','MarketplaceListingMediaPresentation','MarketplaceListingObservedMedia','MarketplaceListingObservationProvenance'
  ]) assert(s[name], `missing schema ${name}`);
  requireFieldsFrom(s, 'MarketplaceListingListItem', ['listing','presentation','lifecycle','observed_at']);
  requireFieldsFrom(s, 'MarketplaceListing', ['listing','presentation','lifecycle','publication_context','observed_fields','observed_media','observed_at','provenance']);
  assert(JSON.stringify(s.ListingObservedFieldKnown).includes('PublicationValueView'), 'known Listing field must use PublicationValueView');
  assert(JSON.stringify(s.MarketplaceListingPerformanceListItem).includes('MarketplaceListingPresentation'), 'Performance Listing item must reuse MarketplaceListingPresentation');
  assert(!Object.hasOwn(s.MarketplaceListingPerformanceListItem?.properties ?? {}, 'display_name'), 'Performance Listing item must not keep parallel display_name');
}
```

Add two negative controls after the existing three: remove collection `presentation`; reintroduce Performance `display_name`. Change the exact final assertion for this task to `negativeControls === 5`.

Run bundle + verifier; expected FAIL on missing MarketplaceListing presentation.

- [ ] **Step 2: Add MarketplaceListing presentation states**

```yaml
  MarketplaceListingPresentationKnown:
    type: object
    additionalProperties: false
    required: [state, display_name]
    properties: {state: {const: known}, display_name: {type: string, minLength: 1}}
  MarketplaceListingPresentationUnknown: {type: object, additionalProperties: false, required: [state], properties: {state: {const: unknown}}}
  MarketplaceListingPresentationUnavailable: {type: object, additionalProperties: false, required: [state], properties: {state: {const: unavailable}}}
  MarketplaceListingPresentation:
    oneOf:
      - {$ref: '#/schemas/MarketplaceListingPresentationKnown'}
      - {$ref: '#/schemas/MarketplaceListingPresentationUnknown'}
      - {$ref: '#/schemas/MarketplaceListingPresentationUnavailable'}
```

Require `presentation` on `MarketplaceListingListItem` next to unchanged `listing: MarketplaceListingRef`.

- [ ] **Step 3: Make observed fields human-readable with exclusive knowledge variants**

Create `ListingObservedFieldKnown|Unknown|Unavailable|NotApplicable`; every variant requires `requirement_key`, non-empty `display_name`, and its fixed state. Only `Known` requires `value: PublicationValueView`. Define `ListingObservedField` as `oneOf` those four variants.

- [ ] **Step 4: Restore observed media and provenance**

Add distinct observed-media presentation:

```yaml
  MarketplaceListingMediaPresentationKnown:
    type: object
    additionalProperties: false
    required: [state, access_ref]
    properties: {state: {const: known}, access_ref: {type: string, format: uri-reference, minLength: 1}}
  MarketplaceListingMediaPresentationUnavailable: {type: object, additionalProperties: false, required: [state], properties: {state: {const: unavailable}}}
  MarketplaceListingMediaPresentation:
    oneOf:
      - {$ref: '#/schemas/MarketplaceListingMediaPresentationKnown'}
      - {$ref: '#/schemas/MarketplaceListingMediaPresentationUnavailable'}
  MarketplaceListingObservedMedia:
    type: object
    additionalProperties: false
    required: [content_type, presentation]
    properties:
      content_type: {type: string, minLength: 1}
      presentation: {$ref: '#/schemas/MarketplaceListingMediaPresentation'}
  MarketplaceListingObservationProvenance:
    type: object
    additionalProperties: false
    required: [source_authority, evidence_custody, acquired_at]
    properties:
      source_authority: {const: marketplace_provider}
      evidence_custody: {type: string, enum: [current_source_observation, preserved_source_evidence]}
      acquired_at: {$ref: '#/schemas/Instant'}
```

Require on `MarketplaceListing`:

```text
listing
presentation
lifecycle
publication_context: PublicationContextView
observed_fields[]
observed_media[]
observed_at
provenance
```

Keep `observed_price` optional; do not add ListingIntent/PriceIntent convergence.

- [ ] **Step 5: Remove Performance's duplicate Listing-label meaning**

In `paths-performance.yaml`, replace local Listing `display_name` fields in `MarketplaceListingPerformanceListItem`, `MarketplaceListingPerformance`, and `RetailMediaListingScope` with required:

```yaml
presentation: {$ref: './components.yaml#/schemas/MarketplaceListingPresentation'}
```

Keep Performance period/coverage/measurement meaning unchanged.

- [ ] **Step 6: Prove and commit**

```powershell
npx --yes @redocly/cli@2.45.0 lint contracts/api/product/openapi.yaml --config contracts/api/product/redocly.yaml
npx --yes @redocly/cli@2.45.0 bundle contracts/api/product/openapi.yaml --config contracts/api/product/redocly.yaml -o .human-projection-test.json
node scripts/verify-human-operable-read-projection.mjs .human-projection-test.json
Remove-Item .human-projection-test.json
node scripts/verify-product-oad.mjs
```

Expected: focused `5/5`; Product `106/106`; generated TS/Go PASS.

```bash
git add contracts/api/product/components.yaml contracts/api/product/paths-performance.yaml contracts/api/product/paths-offering-availability-market.yaml scripts/verify-human-operable-read-projection.mjs
git commit -m "feat(d5): repair marketplace listing read projection"
```

---

### Task 4: Project presentation into the bounded cross-owner consumers proven by P5

**Files:**
- Modify `contracts/api/product/components.yaml`
- Modify `scripts/verify-human-operable-read-projection.mjs`
- Modify operation-description files only when necessary; no path/operation change

**Interfaces:**
- Consumes `SourceProductPresentation`, `MarketplaceListingPresentation`, unchanged canonical target/subject refs.
- Produces adjacent presentation for ListingIntent list, PriceIntent, Availability, Market, Economics only.

- [ ] **Step 1: Verify D3 dependency coherence before changing cross-owner read shapes**

Read accepted D3 owner edges for Offering→Availability and the existing Market/Economics/Performance consumer relationships. If any presentation field below would require a genuinely new semantic owner dependency rather than projecting an already-admitted subject/reference, **STOP / reopen D3** before editing the OAD. Do not create a hidden dependency for UI convenience.

- [ ] **Step 2: Extend verifier RED for exactly the approved consumer set**

Add:

```js
function validateConsumers(doc) {
  const s = schemas(doc);
  requireFieldsFrom(s, 'ListingIntentListItem', ['source_product_presentation']);
  requireFieldsFrom(s, 'PriceIntent', ['target_presentation']);
  requireFieldsFrom(s, 'PriceIntentListItem', ['target_presentation']);
  requireFieldsFrom(s, 'SellableAvailability', ['target_presentation']);
  requireFieldsFrom(s, 'CompetitivePosition', ['subject_presentation']);
  requireFieldsFrom(s, 'CompetitivePositionListItem', ['subject_presentation']);
  requireFieldsFrom(s, 'ExpectedEconomics', ['subject_presentation']);
  requireFieldsFrom(s, 'ExpectedEconomicsListItem', ['subject_presentation']);
  requireFieldsFrom(s, 'PriceScenarioEvaluation', ['subject_presentation']);
  for (const name of ['PriceIntentTargetPresentation','AvailabilityTargetPresentation','MarketSubjectPresentation','EconomicsSubjectPresentation']) assert(s[name], `missing schema ${name}`);
}
```

Add one negative control removing `SellableAvailability.target_presentation`; final exact count becomes `6`. Do not assert presentation on Sale, Shipment, Work, Fulfillment, PostSale, or historical Governance lists.

Run focused verifier; expected FAIL before schema repair.

- [ ] **Step 3: Add only the adjacent presentation unions required by the approved consumers**

Add `source_product_presentation: SourceProductPresentation` to `ListingIntentListItem`.

Define `PriceIntentTargetPresentation`:

```yaml
  PriceIntentExistingTargetPresentation:
    type: object
    additionalProperties: false
    required: [kind, listing_presentation]
    properties: {kind: {const: existing_listing}, listing_presentation: {$ref: '#/schemas/MarketplaceListingPresentation'}}
  PriceIntentPreCreationTargetPresentation:
    type: object
    additionalProperties: false
    required: [kind, source_product_presentation]
    properties: {kind: {const: pre_creation_listing_intent}, source_product_presentation: {$ref: '#/schemas/SourceProductPresentation'}}
  PriceIntentTargetPresentation:
    oneOf:
      - {$ref: '#/schemas/PriceIntentExistingTargetPresentation'}
      - {$ref: '#/schemas/PriceIntentPreCreationTargetPresentation'}
```

Require `target_presentation` on `PriceIntent` and `PriceIntentListItem`; keep `target` key-based.

Define the same two typed meanings under `AvailabilityTargetPresentation`; require it on read-only `SellableAvailability`; keep `AvailabilityTarget`/query params unchanged.

Define `MarketSubjectPresentation` and `EconomicsSubjectPresentation`, each with:

```text
existing_listing -> MarketplaceListingPresentation
source_product_marketplace_context -> SourceProductPresentation
```

Require `subject_presentation` on `CompetitivePosition`, `CompetitivePositionListItem`, `ExpectedEconomics`, `ExpectedEconomicsListItem`, and `PriceScenarioEvaluation`. Keep `MarketSubject`, `EconomicsSubject`, and `EvaluatePriceScenarioRequest.subject` canonical/key-only.

- [ ] **Step 4: Prove the scope did not spread**

Run focused verifier and inspect the diff. Confirm no bulk enrichment of:

```text
MarketplaceSaleListItem
ShipmentListItem
WorkListItem
FulfillmentExecutionListItem
PostSaleResolutionListItem
AuthorizationDecisionListItem
```

Then run `node scripts/verify-product-oad.mjs`; expected `106/106` and focused `6/6`.

- [ ] **Step 5: Commit**

```bash
git add contracts/api/product/components.yaml scripts/verify-human-operable-read-projection.mjs
git commit -m "feat(d5): project listing identity for human consumers"
```

---

### Task 5: Restore ListingIntent current read/media/history conformance without leaking presentation into writes

**Files:**
- Modify `contracts/api/product/components.yaml`
- Modify `contracts/api/product/paths-offering-availability-market.yaml` only if response descriptions need clarification
- Modify `scripts/verify-human-operable-read-projection.mjs`

**Interfaces:**
- Consumes key-only `ListingIntentDesired`, `RequirementResolution`, `MediaSelection`, plus presentation/read families from Tasks 2–4.
- Produces current human-operable ListingIntent read axes and a proportional nested historical attempt basis already required by W2.

- [ ] **Step 1: Extend verifier RED**

Add:

```js
function validateListingIntent(doc) {
  const s = schemas(doc);
  for (const name of [
    'ListingIntentResolvedValueKnown','ListingIntentResolvedValueMissing','ListingIntentResolvedValueUnknown','ListingIntentResolvedValueUnavailable','ListingIntentResolvedValueUnsupported','ListingIntentResolvedValue',
    'ListingIntentFollowSourceResolutionView','ListingIntentExplicitOverrideResolutionView','ListingIntentRequirementResolutionView',
    'ListingIntentMediaPresentationKnown','ListingIntentMediaPresentationUnavailable','ListingIntentMediaPresentation','ListingIntentMediaPresentationDescriptor',
    'ListingIntentAttemptRequirementResolution','ListingIntentAttemptMediaBasis','ListingIntentAttemptAvailabilityInput','ListingIntentEffectAttempt'
  ]) assert(s[name], `missing schema ${name}`);
  requireFieldsFrom(s, 'ListingIntent', ['source_product_presentation','resolved_requirements','authored_media_presentations','dispatch_blockers','created_by_principal_id','updated_by_principal_id','effect_history']);
  assert(JSON.stringify(s.ListingIntentDesired).includes('PublicationContextRef'), 'ListingIntent desired must carry key-based publication context');
  for (const name of ['ListingIntentDesired','RequirementResolution','ExplicitOverrideResolution','PublicationValue','MediaSelection','CreateListingIntentMediaMultipart']) {
    const text = JSON.stringify(s[name] ?? {});
    for (const forbidden of ['display_name','display_label','access_ref','authored_by_principal_id','subject_presentation']) assert(!text.includes(`\"${forbidden}\"`), `${name} must remain presentation-free`);
  }
}
```

Add three negative controls: inject `display_label` into `ExplicitOverrideResolution`; remove `effect_history`; replace authored-media presentation ref with `SourceMediaPresentation`. Final exact count becomes `9`.

Run focused verifier; expected FAIL before implementation.

- [ ] **Step 2: Add key-based publication context to desired state**

Add optional:

```yaml
publication_context: {$ref: '#/schemas/PublicationContextRef'}
```

to `ListingIntentDesired`. This is valid client-authored key meaning; no label enters Create/Update requests.

- [ ] **Step 3: Add current requirement-resolution read views**

Define `ListingIntentResolvedValue` as a `oneOf` of:

```text
known(value: PublicationValueView)
missing
unknown
unavailable
unsupported
```

Define:

```yaml
  ListingIntentFollowSourceResolutionView:
    type: object
    additionalProperties: false
    required: [kind, requirement_key, requirement_display_name, source_candidate_key, source_candidate_display_label, current_value]
    properties:
      kind: {const: follow_source}
      requirement_key: {$ref: '#/schemas/OpaqueKey'}
      requirement_display_name: {type: string, minLength: 1}
      source_candidate_key: {$ref: '#/schemas/OpaqueKey'}
      source_candidate_display_label: {type: string, minLength: 1}
      current_value: {$ref: '#/schemas/ListingIntentResolvedValue'}
  ListingIntentExplicitOverrideResolutionView:
    type: object
    additionalProperties: false
    required: [kind, requirement_key, requirement_display_name, value, authored_by_principal_id, authored_at]
    properties:
      kind: {const: explicit_override}
      requirement_key: {$ref: '#/schemas/OpaqueKey'}
      requirement_display_name: {type: string, minLength: 1}
      value: {$ref: '#/schemas/PublicationValueView'}
      authored_by_principal_id: {$ref: '#/schemas/OpaqueId'}
      authored_at: {$ref: '#/schemas/Instant'}
  ListingIntentRequirementResolutionView:
    oneOf:
      - {$ref: '#/schemas/ListingIntentFollowSourceResolutionView'}
      - {$ref: '#/schemas/ListingIntentExplicitOverrideResolutionView'}
```

Require `resolved_requirements: ListingIntentRequirementResolutionView[]` on `ListingIntent`; keep request `RequirementResolution` unchanged.

- [ ] **Step 4: Restore authored-media read presentation**

```yaml
  ListingIntentMediaPresentationKnown:
    type: object
    additionalProperties: false
    required: [state, access_ref]
    properties: {state: {const: known}, access_ref: {type: string, format: uri-reference, minLength: 1}}
  ListingIntentMediaPresentationUnavailable: {type: object, additionalProperties: false, required: [state], properties: {state: {const: unavailable}}}
  ListingIntentMediaPresentation:
    oneOf:
      - {$ref: '#/schemas/ListingIntentMediaPresentationKnown'}
      - {$ref: '#/schemas/ListingIntentMediaPresentationUnavailable'}
  ListingIntentMediaPresentationDescriptor:
    type: object
    additionalProperties: false
    required: [media, presentation]
    properties:
      media: {$ref: '#/schemas/ListingIntentMediaDescriptor'}
      presentation: {$ref: '#/schemas/ListingIntentMediaPresentation'}
```

Require `authored_media_presentations` on `ListingIntent`. Keep `CreateListingIntentMediaResult` stable descriptor + parent validator only; no access locator in upload result/write/selection/history.

- [ ] **Step 5: Add source presentation, blockers and actor attribution to current ListingIntent read**

Require:

```yaml
source_product_presentation: {$ref: '#/schemas/SourceProductPresentation'}
dispatch_blockers: {type: array, items: {type: string}}
created_by_principal_id: {$ref: '#/schemas/OpaqueId'}
updated_by_principal_id: {$ref: '#/schemas/OpaqueId'}
```

Keep existing lifecycle, dispatchability, effect/convergence and created/updated times.

- [ ] **Step 6: Crystallize the proportional nested W2 historical attempt basis using only accepted identities/types**

Before writing this block, verify every field can be expressed with existing accepted identities/value families. **STOP / split** if it requires a new Product operation, ordinary Permission, business owner, or standalone durable cross-owner resource.

Create:

```yaml
  ListingIntentAttemptFollowSourceResolution:
    type: object
    additionalProperties: false
    required: [kind, requirement_key, source_candidate_key, value, source_observed_at]
    properties:
      kind: {const: follow_source}
      requirement_key: {$ref: '#/schemas/OpaqueKey'}
      source_candidate_key: {$ref: '#/schemas/OpaqueKey'}
      value: {$ref: '#/schemas/PublicationValueView'}
      source_observed_at: {$ref: '#/schemas/Instant'}
  ListingIntentAttemptExplicitOverrideResolution:
    type: object
    additionalProperties: false
    required: [kind, requirement_key, value, authored_by_principal_id, authored_at]
    properties:
      kind: {const: explicit_override}
      requirement_key: {$ref: '#/schemas/OpaqueKey'}
      value: {$ref: '#/schemas/PublicationValueView'}
      authored_by_principal_id: {$ref: '#/schemas/OpaqueId'}
      authored_at: {$ref: '#/schemas/Instant'}
  ListingIntentAttemptRequirementResolution:
    oneOf:
      - {$ref: '#/schemas/ListingIntentAttemptFollowSourceResolution'}
      - {$ref: '#/schemas/ListingIntentAttemptExplicitOverrideResolution'}
  ListingIntentAttemptMediaBasis:
    type: object
    additionalProperties: false
    required: [selection, content_type]
    properties:
      selection: {$ref: '#/schemas/MediaSelection'}
      content_type: {type: string, minLength: 1}
  ListingIntentAttemptAvailabilityInput:
    type: object
    additionalProperties: false
    required: [target, desired_quantity, evaluated_at]
    properties:
      target: {$ref: '#/schemas/AvailabilityTarget'}
      desired_quantity: {$ref: '#/schemas/DesiredQuantity'}
      evaluated_at: {$ref: '#/schemas/Instant'}
  ListingIntentEffectAttempt:
    type: object
    additionalProperties: false
    required: [attempt_key, listing_intent_revision, attempted_at, resolved_requirements, media_basis, external_effect_state, convergence]
    properties:
      attempt_key: {$ref: '#/schemas/OpaqueKey'}
      listing_intent_revision: {$ref: '#/schemas/OpaqueKey'}
      requirements_revision: {$ref: '#/schemas/OpaqueKey'}
      resolved_requirements: {type: array, items: {$ref: '#/schemas/ListingIntentAttemptRequirementResolution'}}
      media_basis: {type: array, items: {$ref: '#/schemas/ListingIntentAttemptMediaBasis'}}
      price_intent_id: {$ref: '#/schemas/OpaqueId'}
      price_intent_etag: {$ref: '#/schemas/StrongETag'}
      availability_input: {$ref: '#/schemas/ListingIntentAttemptAvailabilityInput'}
      authorization_decision_ids: {type: array, uniqueItems: true, items: {$ref: '#/schemas/OpaqueId'}}
      external_effect_state: {type: string, enum: [pending, accepted, rejected, ambiguous]}
      convergence: {type: string, enum: [pending, converged, divergent, unknown, not_applicable]}
      attempted_at: {$ref: '#/schemas/Instant'}
```

Require `effect_history: ListingIntentEffectAttempt[]` on `ListingIntent`. This is a nested historical explanation axis, not a new `PublicationAttempt` Product resource. No presentation access locator enters history.

- [ ] **Step 7: Prove and commit**

Run focused verifier; expected `9/9`. Run `node scripts/verify-product-oad.mjs`; expected `106/106` and generated TS/Go PASS.

```bash
git add contracts/api/product/components.yaml contracts/api/product/paths-offering-availability-market.yaml scripts/verify-human-operable-read-projection.mjs
git commit -m "feat(d5): restore listing intent read conformance"
```

---

### Task 6: Global conformance review, independent challenge, aggregate gate, and integration-ready closure

**Files:**
- Modify `docs/superpowers/specs/2026-08-25-human-operable-read-projection-design.md`
- Modify `docs/roadmap.md`
- Update PR #70 body/comments through GitHub metadata
- Technical files only when a valid review finding proves a defect

**Interfaces:**
- Produces an integration-ready prerequisite candidate; does not merge, re-LOCK B10, rerun B10 P9, or resume B20.

- [ ] **Step 1: Backcheck implementation against the approved spec**

Account explicitly for:

```text
SourceProduct/SourceInstance presentation
requirement/option/unit/context read-write split
FOLLOW_SOURCE source candidates
correspondence candidate population
MarketplaceListing current presentation + actual-state axes
Performance presentation de-duplication
ListingIntent/PriceIntent/Availability/Market/Economics bounded human projections
ListingIntent current resolution/media/actor/blocker axes
source/authored/observed media trust separation
ListingIntent nested historical attempt basis
no write-label authority
no new operation/Permission/Principal kind
```

Any missing item without an explicit STOP outcome blocks closure.

- [ ] **Step 2: Adversarially inspect for accidental platform expansion**

Review `git diff main...HEAD` and reject accidental introduction of:

```text
PresentationService
Generic EntityPresentation/EntityRef
metadata/provider_fields bag
Product/PIM master
new /presentation or /candidates Product path
new operationId/Permission/Principal kind
generic media/asset Product CRUD
```

Do not encode this architecture-quality judgment as prose grep CI; the focused verifier mechanically protects the objective wire fences.

- [ ] **Step 3: Run the one aggregate repository gate**

```powershell
npm run gate
```

Expected evidence includes:

```text
gate: PASS
product_oad_proof: PASS
product_oad_operations=106/106
product_oad_auth_profile=PASS
human_operable_read_projection=PASS
human_operable_read_projection_negative_controls=9/9
```

- [ ] **Step 4: Trigger one fresh independent full PR review**

PR #70 already has CodeRabbit installed. Post this exact PR comment:

```text
@coderabbitai full review
```

This requests a fresh full review of all files rather than another incremental slice. Review findings are evidence, not authority. Adjudicate them against the approved spec and canonical repo authority, with special focus on duplicate/shifted authority, presentation becoming identity/write authority, provider DTO leakage, knowledge-state collapse, hidden cross-owner dependency, current-vs-historical presentation confusion, media trust collapse, operation/Permission expansion, and remaining directly implicated W2/OAD drift.

Fix valid Critical/Important findings before proceeding. Do not implement a reviewer preference that creates a new requirement/authority; record the technical reason when rejecting such a suggestion.

- [ ] **Step 5: Re-run proof after every accepted review fix**

```powershell
node scripts/verify-product-oad.mjs
npm run gate
```

Only the final-HEAD results count as closure evidence.

- [ ] **Step 6: Mark implementation candidate complete, not integrated**

Set:

```text
spec: DESIGN APPROVED / IMPLEMENTATION CANDIDATE COMPLETE — proof and independent challenge complete; integration requires operator authorization
roadmap prerequisite: IMPLEMENTATION CANDIDATE COMPLETE / PROOF PASS / INTEGRATION AUTHORIZATION REQUIRED
roadmap exact next action: obtain explicit operator authorization to integrate PR #70; B20 remains paused; after integration open a fresh bounded B10 correspondence-region revalidation increment
```

Do not mark D6-R2 closed or authorize runtime Product implementation.

- [ ] **Step 7: Commit closure docs**

```bash
git add docs/superpowers/specs/2026-08-25-human-operable-read-projection-design.md docs/roadmap.md
git commit -m "docs(d6-r2): close read projection prerequisite candidate"
```

- [ ] **Step 8: Verify the closure commit itself**

```powershell
npm run gate
```

Expected: PASS on the closure commit.

- [ ] **Step 9: Update PR #70 body**

Record the accepted Global Maximum, canonical authority rehome, Readiness/Offering wire repair, focused negative proof, full gate, independent review disposition, unchanged 106/31/H-A-S inventory, B20 paused state, and exact next gate = explicit operator merge authorization. Do not mark ready/merge without that authorization.

---

## Completion Boundary

This plan is complete only when PR #70 is an **integration-ready prerequisite candidate** with canonical authority rehomed, OAD repaired, focused semantic proof and aggregate gate green, and independent challenge adjudicated.

This plan does **not**:

```text
merge PR #70
edit B10 HTML
re-LOCK B10
rerun B10 P9
render/resume B20
start B23/B24/B30/B40/B50 P8
begin Pre-D9/D9
implement runtime Product code
```

After operator-authorized integration of the prerequisite, open a fresh bounded frontend increment for the B10 correspondence region. That later increment renders real candidate-selection/known-empty/unknown/unavailable states, obtains operator re-LOCK for the affected region, reruns P9 against the repaired OAD, and only then allows PR #69 / B20 to resume.
