# Human-Operable Read Projection & Wire Conformance Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Rehome the approved Global Maximum into canonical D4/D5 authority and repair the Product OAD so human-facing Readiness/Offering reads are operable without changing canonical identity, operation count, Permission count, Principal kinds, or runtime topology.

**Architecture:** Keep canonical refs and all write/decision carriers key/ID based. Add only owner-correct current read projections and purpose/historical snapshots where the approved human job requires recognition, selection, or explanation; repair directly implicated W2→OAD drift; reuse the same MarketplaceListing presentation meaning across Offering and Performance; keep source/authored/observed media presentation trust types distinct. The prerequisite ends when the wire repair is proved and ready for integration; B10 P8/P9 revalidation happens in a later bounded frontend increment after this prerequisite lands.

**Tech Stack:** OpenAPI 3.1.2 YAML; Node.js `>=26.3.0 <27`; PowerShell aggregate gate; `@redocly/cli@2.45.0`; `openapi-typescript@7.13.0`; TypeScript `5.9.3`; Go `1.25.1`; `oapi-codegen v2.8.0`; `github.com/oapi-codegen/runtime v1.7.0`.

**Spec:** `docs/superpowers/specs/2026-08-25-human-operable-read-projection-design.md`

## Global Constraints

- Preserve **106 Product operations / 31 ordinary Permissions / Principal kinds H / A / S**. Any count change is a new material finding and must stop this plan.
- Preserve `MarketplaceListingRef`, `SourceProductRef`, ListingIntent/PriceIntent IDs, correspondence/source candidate keys, option keys, unit keys, and other canonical decision carriers as identity/decision authority.
- Presentation metadata must never authenticate, authorize, scope, correlate, match, converge, or replace canonical identity.
- Client-authored request schemas must not accept presentation labels, presentation locators, server attribution, or source/provider provenance as write authority.
- Preserve `known != missing != unknown != unavailable != unsupported`; known empty remains distinct from unknown/unavailable.
- Do not create `PresentationService`, generic `EntityRef`/`EntityPresentation`, arbitrary metadata bags, provider field bags, Product/PIM master, generic media/asset service, transformation/rule engine, or new Product search/list operation.
- Do not use N+1 point GETs as the production baseline for a collection whose admitted human job requires scan/select/navigation.
- Keep authored/source/observed media presentation references as distinct schema families with distinct trust meaning.
- Runtime remains **NONE** and product implementation paths under `apps|cmd|internal|server|backend|frontend|src|migrations` remain blocked until D9.
- B20 PR #69 remains **PAUSED / NO P8** during this plan.
- B10's main structure remains protected. This plan does not edit B10 HTML or rerun P9; the later bounded frontend increment reopens only the correspondence region after this prerequisite is integrated.
- CI remains one aggregate `npm run gate`. Add only one focused semantic verifier under the existing Product OAD proof; do not add another workflow/check or prose-string ratification test.
- Merge remains separately operator-authorized.

---

## File Structure

### Canonical authority

- `docs/engineering/rebaseline/D4-EXTERNAL-INTEGRATIONS.md` — external evidence/translation rule for human-operable presentation.
- `docs/engineering/rebaseline/D4-R1-PUBLICATION-INPUT-LISTING-AUTHORING-CONTRACT.md` — publication vocabulary, source candidates, correspondence candidates, and presentation/identity fence.
- `docs/engineering/rebaseline/D5-B2-W2-SCHEMA-GRAMMAR.md` — canonical read/write/presentation/snapshot schema grammar.
- `docs/engineering/rebaseline/D5-B2-W3-COLLECTION-GRAMMAR.md` — collection semantic-subset law applied to human scan/select/navigation.

### Canonical machine-readable Product wire

- `contracts/api/product/components.yaml` — shared Product schemas used by Readiness, Offering, Availability, Market, Economics, and selected cross-owner reads.
- `contracts/api/product/paths-identity-portfolio-readiness.yaml` — operation descriptions only when needed to reflect the repaired Readiness semantics; operation IDs/permissions/paths remain unchanged.
- `contracts/api/product/paths-offering-availability-market.yaml` — ListingIntent/MarketplaceListing/PriceIntent/Availability/Market operation descriptions and response bindings as needed; no new path.
- `contracts/api/product/paths-performance.yaml` — remove Performance's parallel Listing label spelling and reuse canonical Listing presentation meaning.
- `contracts/api/product/paths-economics-governance-sales-materialization.yaml` — only if an exact response binding needs adjustment; expected-economics operation surface remains unchanged.
- `contracts/api/product/openapi.yaml` — no new path; touch only if an existing local `$ref` must be redirected after schema factoring.

### Mechanical proof

- `scripts/verify-human-operable-read-projection.mjs` — focused semantic/schema negative proof for this prerequisite.
- `scripts/verify-product-oad.mjs` — invoke the focused verifier against the deterministic resolved bundle; keep the aggregate Product proof as the one gate home.
- `scripts/gate.ps1` — **do not modify** unless execution proves the existing one-line `verify-product-oad.mjs` call cannot carry the targeted proof. Expected plan path is no change.

### Planning/status

- `docs/superpowers/specs/2026-08-25-human-operable-read-projection-design.md` — mark approved/implemented-candidate status as execution advances; never turn it into a second canonical semantic authority.
- `docs/roadmap.md` — sole mutable status/allowed-work/next-action authority.
- `docs/engineering/rebaseline/D6-R2-P9-B10-PREPARATION-SCREEN-CONTRACT.md` — do not rewrite historical P9 during the wire repair; after integration, the later B10 increment records the rerun. The roadmap already owns the current reopen status.

---

### Task 1: Rehome the approved Global Maximum into canonical D4/W2/W3 authority

**Files:**
- Modify: `docs/engineering/rebaseline/D4-EXTERNAL-INTEGRATIONS.md`
- Modify: `docs/engineering/rebaseline/D4-R1-PUBLICATION-INPUT-LISTING-AUTHORING-CONTRACT.md`
- Modify: `docs/engineering/rebaseline/D5-B2-W2-SCHEMA-GRAMMAR.md`
- Modify: `docs/engineering/rebaseline/D5-B2-W3-COLLECTION-GRAMMAR.md`
- Modify: `docs/superpowers/specs/2026-08-25-human-operable-read-projection-design.md`
- Modify: `docs/roadmap.md`

**Interfaces:**
- Consumes: approved design invariant `Canonical Ref != Current Read Projection != Purpose/Historical Snapshot`.
- Produces: canonical authority text that Tasks 2–5 implement mechanically; no operation/path/runtime change.

- [ ] **Step 1: Add the external-evidence presentation fence to D4**

Add a bounded rule under Mercado Livre operational evidence and the general translation fence with this normative meaning:

```markdown
### Human-operable external presentation evidence

When a consumer-owned Product read has a proven human recognition/selection job, D4 preserves the smallest current source/provider presentation evidence needed by that consumer in addition to the canonical external key. Provider/source presentation remains mutable, non-unique evidence and never becomes MPC identity, correspondence authority, authorization, or a generic metadata bag.

For Mercado Livre publication/listing evidence this includes, when applicable, human names/titles associated with the already-used category/product-type/attribute/allowed-value/unit/Listing identities. Adapter-local DTO/field topology remains private.

If presentation cannot currently be acquired, D4 returns honest unknown/unavailable presentation through the consumer-owned semantic port rather than fabricating a label or converting the external key into a name.
```

Do not add a shared provider presentation service or a provider-agnostic free-form field map.

- [ ] **Step 2: Extend D4-R1's publication seam with typed human-operability obligations**

In the Readiness publication-requirements section, add the exact ownership rules:

```text
canonical decision identity
  requirement_key / option_key / unit_key / source_candidate_key / correspondence candidate_key

current human read projection
  provider/source display label + only the bounded context needed to recognize the choice

write/effect
  canonical key only, revalidated against current owner truth
```

State explicitly that Readiness must expose:

```text
SourceProduct/SourceInstance presentation for human source selection
requirement display name
allowed option key + display name
allowed unit key + display name
source candidate key + display label + current value view
correspondence candidate key + display label when a human resolution is admissible
```

and must not expose provider expressions/raw paths/arbitrary maps.

- [ ] **Step 3: Canonicalize W2 read/write/snapshot grammar**

Add a W2 subsection immediately after the request-vs-read law:

```markdown
### Human-operable read projection grammar

A canonical Ref/request carrier remains minimal. A current read may carry an adjacent owner-correct presentation projection when a proven human job requires recognition, selection or explanation. A purpose/historical presentation snapshot is a third meaning and is not refreshed as current source truth.

Presentation is never accepted as identity or write authority. Equal labels never collapse distinct keys. Unknown/unavailable presentation never erases the known canonical subject.
```

Then update the existing W2 Readiness/Offering sections so they normatively name the approved schema families used below:

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

Also state that the current OAD drift for MarketplaceListing/ListingIntent must be repaired without changing W1/W4.

- [ ] **Step 4: Tighten W3 ListItem law without adding generic View machinery**

Add this bounded consequence to W3 §2.2 and the affected matrix rows:

```markdown
When the admitted human consumer must scan/select/navigate a member and the owner can supply current presentation without a second business conclusion, the ListItem carries that owner-semantic presentation subset directly. A point-GET fan-out is not the baseline repair for a deficient collection item. This does not admit a generic `View<T>`, projection DSL, total count, alternate sort or metadata envelope.
```

Call out the currently proven collections only:

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

- [ ] **Step 5: Mark the design approved and roadmap at implementation-plan gate**

Change the spec status to:

```text
DESIGN APPROVED — operator-approved on 2026-08-25; implementation governed by the approved implementation plan and roadmap
```

Change roadmap current prerequisite to `DESIGN APPROVED / IMPLEMENTATION PLAN APPROVAL REQUIRED`, keep B20 paused, and set exact next action to operator approval of this implementation plan. Do not authorize OAD writes yet in the roadmap until the operator approves the plan.

- [ ] **Step 6: Verify the authority-only diff**

Run:

```powershell
git diff --check main...HEAD
git diff --name-only main...HEAD
```

Expected: no whitespace/conflict-marker errors; only the approved design/plan/status and canonical D4/W2/W3 docs are added/modified at this task boundary.

- [ ] **Step 7: Commit the canonical authority rehome**

```bash
git add docs/engineering/rebaseline/D4-EXTERNAL-INTEGRATIONS.md \
        docs/engineering/rebaseline/D4-R1-PUBLICATION-INPUT-LISTING-AUTHORING-CONTRACT.md \
        docs/engineering/rebaseline/D5-B2-W2-SCHEMA-GRAMMAR.md \
        docs/engineering/rebaseline/D5-B2-W3-COLLECTION-GRAMMAR.md \
        docs/superpowers/specs/2026-08-25-human-operable-read-projection-design.md \
        docs/roadmap.md
git commit -m "docs(d6-r2): rehome human-operable read projection authority"
```

---

### Task 2: Repair Readiness publication/source/correspondence read projections with a focused semantic verifier

**Files:**
- Create: `scripts/verify-human-operable-read-projection.mjs`
- Modify: `scripts/verify-product-oad.mjs`
- Modify: `contracts/api/product/components.yaml`
- Modify only if descriptions need semantic clarification: `contracts/api/product/paths-identity-portfolio-readiness.yaml`

**Interfaces:**
- Consumes: Task 1 canonical W2/D4 names and the existing Readiness operations.
- Produces: `SourceProductPresentation`, typed publication read vocabulary/value views, source candidates, correspondence candidate population, source media presentation, and a reusable focused verifier invoked by the existing Product proof.

- [ ] **Step 1: Create the failing Readiness semantic verifier**

Create `scripts/verify-human-operable-read-projection.mjs` with this executable base:

```js
import { readFileSync } from 'node:fs';

const bundlePath = process.argv[2];
if (!bundlePath) throw new Error('usage: node scripts/verify-human-operable-read-projection.mjs <resolved-bundle.json>');
const document = JSON.parse(readFileSync(bundlePath, 'utf8'));
const schemas = document.components?.schemas ?? {};
let negativeControls = 0;

function fail(message) { throw new Error(message); }
function assert(condition, message) { if (!condition) fail(message); }
function schema(name) { const value = schemas[name]; assert(value, `missing schema ${name}`); return value; }
function requireFields(name, fields) {
  const required = new Set(schema(name).required ?? []);
  for (const field of fields) assert(required.has(field), `${name} must require ${field}`);
}
function propertyRef(name, property, ref) {
  assert(schema(name).properties?.[property]?.$ref === ref, `${name}.${property} must reference ${ref}`);
}
function noPropertyDeep(name, forbidden) {
  const text = JSON.stringify(schema(name));
  for (const value of forbidden) assert(!text.includes(`\"${value}\"`), `${name} must not expose write-authority field ${value}`);
}
function expectMutationFailure(label, mutate) {
  const candidate = structuredClone(document);
  let failed = false;
  try { mutate(candidate); validateReadiness(candidate); } catch { failed = true; }
  assert(failed, `negative control unexpectedly passed: ${label}`);
  negativeControls++;
}
function validateReadiness(doc) {
  const s = doc.components.schemas;
  for (const name of [
    'SourceProductPresentationKnown','SourceProductPresentationUnknown','SourceProductPresentationUnavailable','SourceProductPresentation',
    'PublicationCategoryDescriptor','PublicationProductTypeDescriptor','PublicationContextRef','PublicationContextView',
    'PublicationOptionDescriptor','PublicationUnitDescriptor','PublicationValueView','PublicationSourceCandidateView',
    'ProductChannelCorrespondenceCandidate','CorrespondenceCandidatePopulationKnown','CorrespondenceCandidatePopulationUnknown','CorrespondenceCandidatePopulationUnavailable','CorrespondenceCandidatePopulation',
    'SourceMediaPresentationKnown','SourceMediaPresentationUnavailable','SourceMediaPresentation'
  ]) assert(s[name], `missing schema ${name}`);
  for (const field of ['source_product','presentation']) assert((s.SourceProductSearchHit.required ?? []).includes(field), `SourceProductSearchHit must require ${field}`);
  for (const field of ['subject','subject_presentation','correspondence','correspondence_candidate_population','correspondence_etag','readiness','blockers','evaluated_at']) assert((s.ProductChannelReadiness.required ?? []).includes(field), `ProductChannelReadiness must require ${field}`);
  for (const field of ['subject','subject_presentation','publication_context','requirements_revision','requirements','source_media_candidates','evaluated_at']) assert((s.PublicationRequirements.required ?? []).includes(field), `PublicationRequirements must require ${field}`);
  assert((s.PublicationRequirement.required ?? []).includes('display_name'), 'PublicationRequirement must require display_name');
  assert(JSON.stringify(s.PublicationOptionRequirementSpec).includes('PublicationOptionDescriptor'), 'option requirement must expose typed descriptors');
  assert(JSON.stringify(s.PublicationOptionListRequirementSpec).includes('PublicationOptionDescriptor'), 'option-list requirement must expose typed descriptors');
  assert(JSON.stringify(s.PublicationNumberUnitRequirementSpec).includes('PublicationUnitDescriptor'), 'number-unit requirement must expose typed unit descriptors');
  assert(JSON.stringify(s.PublicationSourceEvidenceKnown).includes('PublicationSourceCandidateView'), 'known source evidence must expose candidate views');
  assert(JSON.stringify(s.PublicationSourceEvidenceConflicting).includes('PublicationSourceCandidateView'), 'conflicting source evidence must expose candidate views');
  assert(JSON.stringify(s.SourceMediaCandidate).includes('SourceMediaPresentation'), 'source media candidate must expose its own presentation trust type');
  const writes = ['ResolveCorrespondenceRequest','ClearCorrespondenceRequest','CreateListingIntentDraftRequest','UpdateListingIntentRequest','CreatePriceIntentRequest'];
  for (const name of writes) {
    const text = JSON.stringify(s[name] ?? {});
    for (const forbidden of ['display_name','display_label','access_ref','subject_presentation']) assert(!text.includes(`\"${forbidden}\"`), `${name} must not author ${forbidden}`);
  }
}

validateReadiness(document);
expectMutationFailure('requirement label removed', (d) => { d.components.schemas.PublicationRequirement.required = d.components.schemas.PublicationRequirement.required.filter((x) => x !== 'display_name'); });
expectMutationFailure('correspondence candidate population removed', (d) => { d.components.schemas.ProductChannelReadiness.required = d.components.schemas.ProductChannelReadiness.required.filter((x) => x !== 'correspondence_candidate_population'); });
expectMutationFailure('write accepts presentation label', (d) => { d.components.schemas.ResolveCorrespondenceRequest.properties.display_label = { type: 'string' }; });
assert(negativeControls === 3, `Readiness negative-control count must be 3, found ${negativeControls}`);
console.log('human_operable_read_projection_readiness=PASS');
console.log(`human_operable_read_projection_negative_controls=${negativeControls}/3`);
```

- [ ] **Step 2: Bundle the current baseline and prove the verifier fails for the expected reason**

Run:

```powershell
npx --yes @redocly/cli@2.45.0 bundle contracts/api/product/openapi.yaml --config contracts/api/product/redocly.yaml -o .human-projection-test.json
node scripts/verify-human-operable-read-projection.mjs .human-projection-test.json
Remove-Item .human-projection-test.json
```

Expected: FAIL at the first missing approved schema such as `SourceProductPresentationKnown`. A pass against the pre-repair OAD is a test defect; stop and fix the verifier.

- [ ] **Step 3: Add exact Source Product presentation schemas**

Replace the current flat search-hit presentation with:

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
  SourceProductPresentationUnknown:
    type: object
    additionalProperties: false
    required: [state]
    properties: {state: {const: unknown}}
  SourceProductPresentationUnavailable:
    type: object
    additionalProperties: false
    required: [state]
    properties: {state: {const: unavailable}}
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

Add `subject_presentation: SourceProductPresentation` to both `ProductChannelReadiness` and `PublicationRequirements` and make it required.

- [ ] **Step 4: Split key-based publication context from human read context**

Add:

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
    properties:
      category_key: {$ref: '#/schemas/OpaqueKey'}
      display_name: {type: string, minLength: 1}
  PublicationProductTypeDescriptor:
    type: object
    additionalProperties: false
    required: [product_type_key, display_name]
    properties:
      product_type_key: {$ref: '#/schemas/OpaqueKey'}
      display_name: {type: string, minLength: 1}
  PublicationContextView:
    type: object
    additionalProperties: false
    properties:
      category: {$ref: '#/schemas/PublicationCategoryDescriptor'}
      product_type: {$ref: '#/schemas/PublicationProductTypeDescriptor'}
```

`PublicationRequirements.publication_context` must reference `PublicationContextView`. Keep query parameters key-based. Task 5 uses `PublicationContextRef` inside ListingIntent desired state.

- [ ] **Step 5: Add typed option/unit read descriptors and keep write values key-only**

Add:

```yaml
  PublicationOptionDescriptor:
    type: object
    additionalProperties: false
    required: [option_key, display_name]
    properties:
      option_key: {$ref: '#/schemas/OpaqueKey'}
      display_name: {type: string, minLength: 1}
  PublicationUnitDescriptor:
    type: object
    additionalProperties: false
    required: [unit_key, display_name]
    properties:
      unit_key: {$ref: '#/schemas/OpaqueKey'}
      display_name: {type: string, minLength: 1}
```

Change `PublicationOptionRequirementSpec` and `PublicationOptionListRequirementSpec` from `option_keys` to `options: PublicationOptionDescriptor[]`. Change `PublicationNumberUnitRequirementSpec` from `unit_keys` to `units: PublicationUnitDescriptor[]`; keep `default_unit_key` as the canonical selected key.

Do **not** change client-authored `PublicationOptionValue`, `PublicationOptionListValue`, or `PublicationNumberUnitValue` key-only shapes.

- [ ] **Step 6: Add `PublicationValueView` and typed source candidates**

Add read-only value variants:

```yaml
  PublicationOptionValueView:
    type: object
    additionalProperties: false
    required: [kind, option_key, display_name]
    properties:
      kind: {const: option}
      option_key: {$ref: '#/schemas/OpaqueKey'}
      display_name: {type: string, minLength: 1}
  PublicationOptionListValueView:
    type: object
    additionalProperties: false
    required: [kind, options]
    properties:
      kind: {const: option_list}
      options: {type: array, items: {$ref: '#/schemas/PublicationOptionDescriptor'}}
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

Change `PublicationSourceEvidenceKnown.candidates` to `array`, `minItems: 1`, items `PublicationSourceCandidateView`; change conflicting to `array`, `minItems: 2` with the same item schema. Preserve missing/unknown/unavailable/unsupported as explicit no-candidate variants.

Make `PublicationRequirement.display_name` non-empty and required.

- [ ] **Step 7: Add correspondence candidate-population semantics without a new operation**

Add:

```yaml
  ProductChannelCorrespondenceCandidate:
    type: object
    additionalProperties: false
    required: [candidate_key, display_label]
    properties:
      candidate_key: {$ref: '#/schemas/OpaqueKey'}
      display_label: {type: string, minLength: 1}
  CorrespondenceCandidatePopulationKnown:
    type: object
    additionalProperties: false
    required: [state, candidates]
    properties:
      state: {const: known}
      candidates: {type: array, items: {$ref: '#/schemas/ProductChannelCorrespondenceCandidate'}}
  CorrespondenceCandidatePopulationUnknown:
    type: object
    additionalProperties: false
    required: [state]
    properties: {state: {const: unknown}}
  CorrespondenceCandidatePopulationUnavailable:
    type: object
    additionalProperties: false
    required: [state]
    properties: {state: {const: unavailable}}
  CorrespondenceCandidatePopulation:
    oneOf:
      - {$ref: '#/schemas/CorrespondenceCandidatePopulationKnown'}
      - {$ref: '#/schemas/CorrespondenceCandidatePopulationUnknown'}
      - {$ref: '#/schemas/CorrespondenceCandidatePopulationUnavailable'}
```

Add required `correspondence_candidate_population` to `ProductChannelReadiness`. Keep `ResolveCorrespondenceRequest` exactly key-based: `subject + correspondence_etag + candidate_key`.

- [ ] **Step 8: Add source-media presentation as a distinct trust type**

Add:

```yaml
  SourceMediaPresentationKnown:
    type: object
    additionalProperties: false
    required: [state, access_ref]
    properties:
      state: {const: known}
      access_ref: {type: string, format: uri-reference, minLength: 1}
  SourceMediaPresentationUnavailable:
    type: object
    additionalProperties: false
    required: [state]
    properties: {state: {const: unavailable}}
  SourceMediaPresentation:
    oneOf:
      - {$ref: '#/schemas/SourceMediaPresentationKnown'}
      - {$ref: '#/schemas/SourceMediaPresentationUnavailable'}
  SourceMediaCandidate:
    type: object
    additionalProperties: false
    required: [source_media_candidate_key, content_type, presentation]
    properties:
      source_media_candidate_key: {$ref: '#/schemas/OpaqueKey'}
      content_type: {type: string, minLength: 1}
      presentation: {$ref: '#/schemas/SourceMediaPresentation'}
```

The access reference is response-only. Never add it to `MediaSelection` or client upload/write schemas.

- [ ] **Step 9: Run the focused proof, then hook it into the existing Product verifier**

Run the bundle + verifier command from Step 2. Expected: PASS.

Then in `scripts/verify-product-oad.mjs` add near the other proof-path constants:

```js
const humanOperableReadProjectionVerifier = join(root, 'scripts/verify-human-operable-read-projection.mjs');
```

and inside `currentAuthProof()` after `currentProjectionProof(bundleA);` add:

```js
run(process.execPath, [humanOperableReadProjectionVerifier, bundleA]);
```

Run:

```powershell
node scripts/verify-product-oad.mjs
```

Expected: existing historical/current/auth/generated proofs PASS plus `human_operable_read_projection_readiness=PASS`; operation count remains `106/106`.

- [ ] **Step 10: Commit the Readiness repair**

```bash
git add contracts/api/product/components.yaml \
        contracts/api/product/paths-identity-portfolio-readiness.yaml \
        scripts/verify-human-operable-read-projection.mjs \
        scripts/verify-product-oad.mjs
git commit -m "feat(d5): make readiness reads human-operable"
```

---

### Task 3: Repair MarketplaceListing actual-state reads and remove Performance's parallel Listing-label spelling

**Files:**
- Modify: `contracts/api/product/components.yaml`
- Modify: `contracts/api/product/paths-offering-availability-market.yaml` only if descriptions/bindings need clarification
- Modify: `contracts/api/product/paths-performance.yaml`
- Modify: `scripts/verify-human-operable-read-projection.mjs`

**Interfaces:**
- Consumes: `PublicationContextView`, `PublicationValueView`, canonical `MarketplaceListingRef`.
- Produces: `MarketplaceListingPresentation`, complete-enough Listing actual-state read, observed-media trust type, and one shared Listing presentation meaning reused by Performance.

- [ ] **Step 1: Extend the verifier with failing MarketplaceListing assertions**

Add to the verifier:

```js
function validateMarketplaceListing(doc) {
  const s = doc.components.schemas;
  for (const name of [
    'MarketplaceListingPresentationKnown','MarketplaceListingPresentationUnknown','MarketplaceListingPresentationUnavailable','MarketplaceListingPresentation',
    'MarketplaceListingObservedMedia','MarketplaceListingMediaPresentationKnown','MarketplaceListingMediaPresentationUnavailable','MarketplaceListingMediaPresentation',
    'MarketplaceListingObservationProvenance'
  ]) assert(s[name], `missing schema ${name}`);
  for (const field of ['listing','presentation','lifecycle','observed_at']) assert((s.MarketplaceListingListItem.required ?? []).includes(field), `MarketplaceListingListItem must require ${field}`);
  for (const field of ['listing','presentation','lifecycle','publication_context','observed_fields','observed_media','observed_at','provenance']) assert((s.MarketplaceListing.required ?? []).includes(field), `MarketplaceListing must require ${field}`);
  assert(JSON.stringify(s.ListingObservedField).includes('PublicationValueView'), 'ListingObservedField must use PublicationValueView');
  const perfText = JSON.stringify(s.MarketplaceListingPerformanceListItem ?? {});
  assert(perfText.includes('MarketplaceListingPresentation'), 'Performance Listing item must reuse MarketplaceListingPresentation');
  assert(!Object.hasOwn(s.MarketplaceListingPerformanceListItem?.properties ?? {}, 'display_name'), 'Performance Listing item must not keep a parallel display_name spelling');
}
```

Call `validateMarketplaceListing(document)` from the main verifier and add two negative controls: remove Listing presentation from the collection item; restore a raw Performance `display_name`. Update the expected negative-control count accordingly.

- [ ] **Step 2: Prove the extended verifier fails on the current branch state**

Bundle and run the focused verifier. Expected: FAIL with missing `MarketplaceListingPresentationKnown`.

- [ ] **Step 3: Add MarketplaceListing presentation knowledge states**

Add:

```yaml
  MarketplaceListingPresentationKnown:
    type: object
    additionalProperties: false
    required: [state, display_name]
    properties:
      state: {const: known}
      display_name: {type: string, minLength: 1}
  MarketplaceListingPresentationUnknown:
    type: object
    additionalProperties: false
    required: [state]
    properties: {state: {const: unknown}}
  MarketplaceListingPresentationUnavailable:
    type: object
    additionalProperties: false
    required: [state]
    properties: {state: {const: unavailable}}
  MarketplaceListingPresentation:
    oneOf:
      - {$ref: '#/schemas/MarketplaceListingPresentationKnown'}
      - {$ref: '#/schemas/MarketplaceListingPresentationUnknown'}
      - {$ref: '#/schemas/MarketplaceListingPresentationUnavailable'}
```

Change `MarketplaceListingListItem` to require `presentation` adjacent to the unchanged `listing: MarketplaceListingRef`.

- [ ] **Step 4: Make observed fields human-readable without weakening knowledge state**

Replace the flat `ListingObservedField` with a union whose variants all require `requirement_key` and non-empty `display_name`:

```yaml
  ListingObservedFieldKnown:
    type: object
    additionalProperties: false
    required: [requirement_key, display_name, state, value]
    properties:
      requirement_key: {$ref: '#/schemas/OpaqueKey'}
      display_name: {type: string, minLength: 1}
      state: {const: known}
      value: {$ref: '#/schemas/PublicationValueView'}
  ListingObservedFieldUnknown:
    type: object
    additionalProperties: false
    required: [requirement_key, display_name, state]
    properties:
      requirement_key: {$ref: '#/schemas/OpaqueKey'}
      display_name: {type: string, minLength: 1}
      state: {const: unknown}
  ListingObservedFieldUnavailable:
    type: object
    additionalProperties: false
    required: [requirement_key, display_name, state]
    properties:
      requirement_key: {$ref: '#/schemas/OpaqueKey'}
      display_name: {type: string, minLength: 1}
      state: {const: unavailable}
  ListingObservedFieldNotApplicable:
    type: object
    additionalProperties: false
    required: [requirement_key, display_name, state]
    properties:
      requirement_key: {$ref: '#/schemas/OpaqueKey'}
      display_name: {type: string, minLength: 1}
      state: {const: not_applicable}
  ListingObservedField:
    oneOf:
      - {$ref: '#/schemas/ListingObservedFieldKnown'}
      - {$ref: '#/schemas/ListingObservedFieldUnknown'}
      - {$ref: '#/schemas/ListingObservedFieldUnavailable'}
      - {$ref: '#/schemas/ListingObservedFieldNotApplicable'}
```

- [ ] **Step 5: Restore MarketplaceListing media/provenance axes already required by W2**

Add distinct observed-media presentation:

```yaml
  MarketplaceListingMediaPresentationKnown:
    type: object
    additionalProperties: false
    required: [state, access_ref]
    properties:
      state: {const: known}
      access_ref: {type: string, format: uri-reference, minLength: 1}
  MarketplaceListingMediaPresentationUnavailable:
    type: object
    additionalProperties: false
    required: [state]
    properties: {state: {const: unavailable}}
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

Repair `MarketplaceListing` to require:

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

Keep `observed_price` optional evidence. Do not add ListingIntent/PriceIntent convergence to MarketplaceListing.

- [ ] **Step 6: Normalize Performance to the same Listing presentation meaning**

In `paths-performance.yaml`, change `MarketplaceListingPerformanceListItem` and `MarketplaceListingPerformance` so the existing label field becomes:

```yaml
presentation: {$ref: './components.yaml#/schemas/MarketplaceListingPresentation'}
```

and is required. Remove the local `display_name` property. Preserve Performance owner-specific period/coverage/measurement semantics unchanged.

For `RetailMediaListingScope`, replace the local Listing `display_name` with `presentation: MarketplaceListingPresentation`; keep its exact Listing key identity semantics and path/Installation qualification unchanged.

- [ ] **Step 7: Run targeted and generator proofs**

Run:

```powershell
npx --yes @redocly/cli@2.45.0 lint contracts/api/product/openapi.yaml --config contracts/api/product/redocly.yaml
npx --yes @redocly/cli@2.45.0 bundle contracts/api/product/openapi.yaml --config contracts/api/product/redocly.yaml -o .human-projection-test.json
node scripts/verify-human-operable-read-projection.mjs .human-projection-test.json
Remove-Item .human-projection-test.json
node scripts/verify-product-oad.mjs
```

Expected: lint PASS; focused verifier PASS; Product OAD generated TypeScript/Go projections PASS; `106/106` unchanged.

- [ ] **Step 8: Commit the Listing actual-state repair**

```bash
git add contracts/api/product/components.yaml \
        contracts/api/product/paths-offering-availability-market.yaml \
        contracts/api/product/paths-performance.yaml \
        scripts/verify-human-operable-read-projection.mjs
git commit -m "feat(d5): repair marketplace listing read projection"
```

---

### Task 4: Make the already-proven cross-owner Listing/Product consumers human-operable without enriching every opaque ID

**Files:**
- Modify: `contracts/api/product/components.yaml`
- Modify: `scripts/verify-human-operable-read-projection.mjs`
- Modify only if response description needs clarification: `contracts/api/product/paths-offering-availability-market.yaml`
- Modify only if response description needs clarification: `contracts/api/product/paths-economics-governance-sales-materialization.yaml`

**Interfaces:**
- Consumes: `SourceProductPresentation`, `MarketplaceListingPresentation`, unchanged canonical target/subject refs.
- Produces: adjacent presentation axes for ListingIntent, PriceIntent, Availability, Market, Economics; no new operation or cross-owner write authority.

- [ ] **Step 1: Add failing verifier assertions for only the approved consumer set**

Extend the verifier:

```js
function validateConsumers(doc) {
  const s = doc.components.schemas;
  requireFieldsFrom(s, 'ListingIntentListItem', ['source_product_presentation']);
  requireFieldsFrom(s, 'PriceIntentListItem', ['target_presentation']);
  requireFieldsFrom(s, 'PriceIntent', ['target_presentation']);
  requireFieldsFrom(s, 'SellableAvailability', ['target_presentation']);
  requireFieldsFrom(s, 'CompetitivePosition', ['subject_presentation']);
  requireFieldsFrom(s, 'CompetitivePositionListItem', ['subject_presentation']);
  requireFieldsFrom(s, 'ExpectedEconomics', ['subject_presentation']);
  requireFieldsFrom(s, 'ExpectedEconomicsListItem', ['subject_presentation']);
  requireFieldsFrom(s, 'PriceScenarioEvaluation', ['subject_presentation']);
  for (const name of ['PriceIntentTargetPresentation','AvailabilityTargetPresentation','MarketSubjectPresentation','EconomicsSubjectPresentation']) assert(s[name], `missing schema ${name}`);
}
```

Add helper `requireFieldsFrom(s, name, fields)` locally. Add a negative control that removes `SellableAvailability.target_presentation` and must fail. Do **not** assert presentation on Sales, Shipment, Work, Fulfillment, PostSale or historical Governance lists.

- [ ] **Step 2: Prove the new consumer assertions fail before the repair**

Bundle and run the focused verifier. Expected: FAIL because `ListingIntentListItem` does not yet require `source_product_presentation`.

- [ ] **Step 3: Add ListingIntent list/source presentation only where needed now**

Add `source_product_presentation: SourceProductPresentation` to `ListingIntentListItem` and make it required. Task 5 adds the same field to point `ListingIntent`.

Do not create a composite `ListingIntent.display_name`; the human subject remains the already-known source Product plus explicit target/lifecycle fields.

- [ ] **Step 4: Add PriceIntent target presentation as a typed adjacent axis**

Add:

```yaml
  PriceIntentExistingTargetPresentation:
    type: object
    additionalProperties: false
    required: [kind, listing_presentation]
    properties:
      kind: {const: existing_listing}
      listing_presentation: {$ref: '#/schemas/MarketplaceListingPresentation'}
  PriceIntentPreCreationTargetPresentation:
    type: object
    additionalProperties: false
    required: [kind, source_product_presentation]
    properties:
      kind: {const: pre_creation_listing_intent}
      source_product_presentation: {$ref: '#/schemas/SourceProductPresentation'}
  PriceIntentTargetPresentation:
    oneOf:
      - {$ref: '#/schemas/PriceIntentExistingTargetPresentation'}
      - {$ref: '#/schemas/PriceIntentPreCreationTargetPresentation'}
```

Require `target_presentation` on `PriceIntent` and `PriceIntentListItem`. Keep `target: PriceIntentTarget` unchanged and key-based.

- [ ] **Step 5: Add Availability target presentation without changing its target identity**

Add:

```yaml
  AvailabilityExistingTargetPresentation:
    type: object
    additionalProperties: false
    required: [kind, listing_presentation]
    properties:
      kind: {const: existing_listing}
      listing_presentation: {$ref: '#/schemas/MarketplaceListingPresentation'}
  AvailabilityPreCreationTargetPresentation:
    type: object
    additionalProperties: false
    required: [kind, source_product_presentation]
    properties:
      kind: {const: pre_creation_listing_intent}
      source_product_presentation: {$ref: '#/schemas/SourceProductPresentation'}
  AvailabilityTargetPresentation:
    oneOf:
      - {$ref: '#/schemas/AvailabilityExistingTargetPresentation'}
      - {$ref: '#/schemas/AvailabilityPreCreationTargetPresentation'}
```

Require `target_presentation` on read-only `SellableAvailability`. Leave `AvailabilityTarget` and query target params unchanged.

- [ ] **Step 6: Add Market/Economics subject presentation while preserving canonical subject refs**

Add owner-local adjacent presentation unions:

```yaml
  MarketExistingListingSubjectPresentation:
    type: object
    additionalProperties: false
    required: [kind, listing_presentation]
    properties:
      kind: {const: existing_listing}
      listing_presentation: {$ref: '#/schemas/MarketplaceListingPresentation'}
  MarketSourceProductSubjectPresentation:
    type: object
    additionalProperties: false
    required: [kind, source_product_presentation]
    properties:
      kind: {const: source_product_marketplace_context}
      source_product_presentation: {$ref: '#/schemas/SourceProductPresentation'}
  MarketSubjectPresentation:
    oneOf:
      - {$ref: '#/schemas/MarketExistingListingSubjectPresentation'}
      - {$ref: '#/schemas/MarketSourceProductSubjectPresentation'}
```

Create the same two-variant shape under `EconomicsSubjectPresentation` using Economics names. Require:

```text
CompetitivePosition.subject_presentation
CompetitivePositionListItem.subject_presentation
ExpectedEconomics.subject_presentation
ExpectedEconomicsListItem.subject_presentation
PriceScenarioEvaluation.subject_presentation
```

Keep `MarketSubject`, `EconomicsSubject`, and `EvaluatePriceScenarioRequest.subject` canonical/key-based.

- [ ] **Step 7: Run focused proof and verify scope did not spread**

Run the focused verifier and then inspect:

```powershell
git diff -- contracts/api/product/components.yaml scripts/verify-human-operable-read-projection.mjs
```

Confirm there is no bulk enrichment of `MarketplaceSaleListItem`, `ShipmentListItem`, `WorkListItem`, `FulfillmentExecutionListItem`, `PostSaleResolutionListItem`, or historical `AuthorizationDecisionListItem` merely because they contain IDs.

Then run:

```powershell
node scripts/verify-product-oad.mjs
```

Expected: PASS and `106/106`.

- [ ] **Step 8: Commit the bounded consumer repair**

```bash
git add contracts/api/product/components.yaml scripts/verify-human-operable-read-projection.mjs
git commit -m "feat(d5): project listing identity for human consumers"
```

---

### Task 5: Restore ListingIntent current read/media/history conformance without leaking presentation into writes

**Files:**
- Modify: `contracts/api/product/components.yaml`
- Modify: `contracts/api/product/paths-offering-availability-market.yaml` only if response descriptions need clarification
- Modify: `scripts/verify-human-operable-read-projection.mjs`

**Interfaces:**
- Consumes: key-only `ListingIntentDesired`/`RequirementResolution`/`MediaSelection`, `PublicationContextRef`, `PublicationValueView`, source/authored media presentation families, `SourceProductPresentation`.
- Produces: current human-operable ListingIntent read axes and the proportional historical attempt basis already required by W2, while write schemas stay label-free.

- [ ] **Step 1: Extend the verifier with failing ListingIntent conformance assertions**

Add:

```js
function validateListingIntent(doc) {
  const s = doc.components.schemas;
  for (const name of [
    'ListingIntentResolvedValueKnown','ListingIntentResolvedValueMissing','ListingIntentResolvedValueUnknown','ListingIntentResolvedValueUnavailable','ListingIntentResolvedValueUnsupported','ListingIntentResolvedValue',
    'ListingIntentFollowSourceResolutionView','ListingIntentExplicitOverrideResolutionView','ListingIntentRequirementResolutionView',
    'ListingIntentMediaPresentationKnown','ListingIntentMediaPresentationUnavailable','ListingIntentMediaPresentation','ListingIntentMediaPresentationDescriptor',
    'ListingIntentAttemptRequirementResolution','ListingIntentAttemptMediaBasis','ListingIntentAttemptAvailabilityInput','ListingIntentEffectAttempt'
  ]) assert(s[name], `missing schema ${name}`);
  for (const field of ['source_product_presentation','resolved_requirements','authored_media_presentations','dispatch_blockers','created_by_principal_id','updated_by_principal_id','effect_history']) assert((s.ListingIntent.required ?? []).includes(field), `ListingIntent must require ${field}`);
  assert(JSON.stringify(s.ListingIntentDesired).includes('PublicationContextRef'), 'ListingIntent desired state must carry key-based publication context');
  const writes = ['CreateListingIntentDraftRequest','UpdateListingIntentRequest','RequirementResolution','MediaSelection','CreateListingIntentMediaMultipart'];
  for (const name of writes) {
    const text = JSON.stringify(s[name] ?? {});
    for (const forbidden of ['display_name','display_label','access_ref','authored_by_principal_id','subject_presentation']) assert(!text.includes(`\"${forbidden}\"`), `${name} must remain presentation-free`);
  }
}
```

Add negative controls that: inject `display_label` into `ExplicitOverrideResolution`; remove `effect_history` from ListingIntent; replace authored-media presentation with `SourceMediaPresentation`. Each must fail.

- [ ] **Step 2: Prove the ListingIntent verifier section fails before implementation**

Bundle and run the focused verifier. Expected: FAIL at a missing ListingIntent read schema such as `ListingIntentResolvedValueKnown`.

- [ ] **Step 3: Put publication context in desired state as keys, not labels**

Add optional:

```yaml
publication_context: {$ref: '#/schemas/PublicationContextRef'}
```

to `ListingIntentDesired`. Because `CreateListingIntentDraftRequest` and `UpdateListingIntentRequest` reuse desired state, this remains key-only authored meaning.

- [ ] **Step 4: Add current requirement-resolution read views**

Add value-knowledge union:

```yaml
  ListingIntentResolvedValueKnown:
    type: object
    additionalProperties: false
    required: [state, value]
    properties: {state: {const: known}, value: {$ref: '#/schemas/PublicationValueView'}}
  ListingIntentResolvedValueMissing: {type: object, additionalProperties: false, required: [state], properties: {state: {const: missing}}}
  ListingIntentResolvedValueUnknown: {type: object, additionalProperties: false, required: [state], properties: {state: {const: unknown}}}
  ListingIntentResolvedValueUnavailable: {type: object, additionalProperties: false, required: [state], properties: {state: {const: unavailable}}}
  ListingIntentResolvedValueUnsupported: {type: object, additionalProperties: false, required: [state], properties: {state: {const: unsupported}}}
  ListingIntentResolvedValue:
    oneOf:
      - {$ref: '#/schemas/ListingIntentResolvedValueKnown'}
      - {$ref: '#/schemas/ListingIntentResolvedValueMissing'}
      - {$ref: '#/schemas/ListingIntentResolvedValueUnknown'}
      - {$ref: '#/schemas/ListingIntentResolvedValueUnavailable'}
      - {$ref: '#/schemas/ListingIntentResolvedValueUnsupported'}
```

Add read-only resolution views:

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

Add required `resolved_requirements: ListingIntentRequirementResolutionView[]` to `ListingIntent`. Keep client-authored `RequirementResolution` unchanged.

- [ ] **Step 5: Restore authored-media read presentation while preserving the accepted technical-delivery boundary**

Add:

```yaml
  ListingIntentMediaPresentationKnown:
    type: object
    additionalProperties: false
    required: [state, access_ref]
    properties:
      state: {const: known}
      access_ref: {type: string, format: uri-reference, minLength: 1}
  ListingIntentMediaPresentationUnavailable:
    type: object
    additionalProperties: false
    required: [state]
    properties: {state: {const: unavailable}}
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

Add required `authored_media_presentations: ListingIntentMediaPresentationDescriptor[]` to `ListingIntent`. Keep `CreateListingIntentMediaResult` as stable descriptor + parent validator only; do not add `access_ref` to the create result, upload body, desired media selection, or history.

- [ ] **Step 6: Add current source presentation, blockers and actor attribution to ListingIntent read**

Make these required on `ListingIntent`:

```yaml
source_product_presentation: {$ref: '#/schemas/SourceProductPresentation'}
dispatch_blockers: {type: array, items: {type: string}}
created_by_principal_id: {$ref: '#/schemas/OpaqueId'}
updated_by_principal_id: {$ref: '#/schemas/OpaqueId'}
```

Existing `created_at`/`updated_at`, lifecycle, dispatchability, external-effect state and convergence remain unchanged.

- [ ] **Step 7: Crystallize the proportional W2 historical attempt basis using only already-accepted identities/value families**

Before writing YAML, verify every field below is expressible with existing accepted types. **STOP / split a prerequisite** if implementing this step requires a new Product operation, ordinary Permission, business owner, or new durable cross-owner identity. A new nested historical occurrence key inside ListingIntent is allowed; a new standalone Product resource is not.

Add:

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

Add required `effect_history: ListingIntentEffectAttempt[]` to `ListingIntent`.

This is a nested historical explanation axis, not a new `PublicationAttempt` resource. No presentation locator appears in history. Current labels in effect-history value views are historical snapshots within the ListingIntent attempt basis; canonical keys remain present.

- [ ] **Step 8: Run the full focused verifier and write-authority negative controls**

Run:

```powershell
npx --yes @redocly/cli@2.45.0 bundle contracts/api/product/openapi.yaml --config contracts/api/product/redocly.yaml -o .human-projection-test.json
node scripts/verify-human-operable-read-projection.mjs .human-projection-test.json
Remove-Item .human-projection-test.json
```

Expected: PASS. Verify negative controls fail when presentation fields are injected into write schemas and when authored/source media trust schemas are swapped.

- [ ] **Step 9: Run canonical Product OAD proof**

```powershell
node scripts/verify-product-oad.mjs
```

Expected: Redocly/source bundle deterministic; generated TS/Go compile; auth/security proof unchanged; `product_oad_operations=106/106`; human-operable verifier PASS.

- [ ] **Step 10: Commit ListingIntent conformance**

```bash
git add contracts/api/product/components.yaml \
        contracts/api/product/paths-offering-availability-market.yaml \
        scripts/verify-human-operable-read-projection.mjs
git commit -m "feat(d5): restore listing intent read conformance"
```

---

### Task 6: Global conformance review, independent challenge, aggregate gate, and integration-ready closure

**Files:**
- Modify: `docs/superpowers/specs/2026-08-25-human-operable-read-projection-design.md`
- Modify: `docs/roadmap.md`
- Modify: PR #70 body through GitHub metadata only
- Modify technical files only if a Critical/Important review finding proves a defect against the approved spec

**Interfaces:**
- Consumes: Tasks 1–5 complete wire repair and mechanical proof.
- Produces: one reviewable prerequisite candidate ready for explicit operator merge authorization; B20 remains paused and B10 revalidation remains the next post-integration increment.

- [ ] **Step 1: Run the spec-coverage backcheck**

Read the approved design beside the diff and explicitly verify each required area has a concrete implementation/proof home:

```text
SourceProduct/SourceInstance presentation
requirement/option/unit/context read/write split
FOLLOW_SOURCE source-candidate presentation
correspondence candidate population
MarketplaceListing current presentation + actual-state axes
Performance de-duplication
ListingIntent/PriceIntent/Availability/Market/Economics human consumer projections
ListingIntent current resolved requirements
source/authored/observed media trust separation
ListingIntent historical attempt basis
no write-label authority
no new operation/Permission/Principal kind
```

If any item has no implementation or deliberate stop outcome, do not proceed.

- [ ] **Step 2: Scan for accidental platform expansion**

Run:

```powershell
git diff main...HEAD -- contracts/api/product docs/engineering/rebaseline scripts
```

Reject any accidental introduction of:

```text
PresentationService
EntityPresentation / generic EntityRef
metadata bag / arbitrary map
provider_fields
Product/PIM master
new /presentation or /candidates Product path
new operationId
new ordinary Permission
new Principal kind
new generic media/asset Product operation
```

The focused verifier should mechanically cover operation counts and key write-authority fences; architecture-quality judgment remains review, not CI string matching.

- [ ] **Step 3: Run the one aggregate repository gate**

```powershell
npm run gate
```

Expected terminal evidence includes:

```text
gate: PASS
product_oad_proof: PASS
product_oad_operations=106/106
product_oad_auth_profile=PASS
human_operable_read_projection_readiness=PASS
```

Do not add a second workflow/check for this prerequisite.

- [ ] **Step 4: Request a fresh independent code/contract review**

Use the `superpowers:requesting-code-review` skill against:

```text
BASE_SHA = main at the prerequisite branch point
HEAD_SHA = current prerequisite HEAD
DESCRIPTION = Human-operable Readiness/Offering read projection + W2/OAD conformance repair
PLAN_OR_REQUIREMENTS = this plan + docs/superpowers/specs/2026-08-25-human-operable-read-projection-design.md
```

Reviewer focus must include:

```text
duplicate/shifted authority
presentation becoming identity/write authority
provider DTO leakage
unknown/empty collapse
cross-owner dependency invention
historical/current presentation confusion
media trust collapse
new operation/Permission disguised as repair
W2/OAD drift left unresolved in directly implicated areas
```

Fix every valid Critical/Important finding before continuing. Minor findings are fixed only when they improve correctness/clarity without scope expansion; otherwise record why they are non-blocking.

- [ ] **Step 5: Rerun targeted proof and aggregate gate after review fixes**

```powershell
node scripts/verify-product-oad.mjs
npm run gate
```

Expected: both PASS on the final HEAD. Do not claim completion from an earlier run.

- [ ] **Step 6: Mark the spec and roadmap as implementation-complete candidate, not operator-accepted integration**

Update the spec header to:

```text
DESIGN APPROVED / IMPLEMENTATION CANDIDATE COMPLETE — proof + independent challenge complete; integration still requires operator authorization
```

Update roadmap current prerequisite to:

```text
IMPLEMENTATION CANDIDATE COMPLETE / PROOF PASS / INTEGRATION AUTHORIZATION REQUIRED
```

Set exact next action:

```text
Obtain explicit operator authorization to integrate PR #70. Do not resume B20. After prerequisite integration, open a fresh bounded B10 correspondence-region P8 revalidation increment, rerun P9 after operator re-LOCK, then resume B20.
```

Do not mark D6-R2 accepted/closed and do not authorize Product runtime implementation.

- [ ] **Step 7: Commit closure metadata/docs**

```bash
git add docs/superpowers/specs/2026-08-25-human-operable-read-projection-design.md docs/roadmap.md
git commit -m "docs(d6-r2): close read projection prerequisite candidate"
```

- [ ] **Step 8: Re-run final verification after the closure commit**

```powershell
npm run gate
```

Expected: PASS on the closure commit itself.

- [ ] **Step 9: Update PR #70 body without changing files**

Record:

```text
approved Global Maximum
canonical authority rehome complete
Readiness/Offering OAD repair complete
focused semantic negative proof complete
full gate PASS
independent challenge disposition
Product remains 106/31/H-A-S
B20 remains paused
next gate = explicit operator integration authorization
```

Do not mark the PR ready/merge it until explicit operator authorization.

---

## Plan Completion Boundary

This plan is complete only when PR #70 is an **integration-ready prerequisite candidate** with canonical authority rehomed, OAD repaired, focused semantic proof and aggregate gate green, and independent challenge adjudicated.

This plan explicitly does **not**:

```text
merge PR #70
edit B10 HTML
re-LOCK B10
rerun B10 P9
render B20
start B23/B24/B30/B40/B50 P8
begin Pre-D9/D9
implement runtime/Product code
```

After operator-authorized integration of the prerequisite, create a fresh bounded frontend increment for the B10 correspondence region. That increment must supply real candidate selection states in executable low-fidelity HTML, obtain operator re-LOCK for the affected region, rerun P9 against the repaired canonical OAD, and only then allow PR #69 / B20 to resume.
