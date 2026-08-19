# D5-B2 — W2 Request / Response Schema Grammar

> **Status:** OPEN / ACTIVE — W2-A Core Schema Grammar ACCEPTED IN-STAGE; W2-B ListingIntent + PriceIntent + Availability schemas next  
> **Parent Wire Contract:** `D5-B2-WIRE-CONTRACT.md`  
> **Operation inventory:** `D5-B2-OPERATION-ADMISSION-MATRIX.md`  
> **Parent authorities:** accepted D0→D4 + D4-R1 + D5-B1 + ratified D5-B2 Whole-Matrix + Wire W1  
> **Method:** DevelopmentConexus Engineering Method v1.0.0  
> **W2-A accepted:** 2026-08-18

## 1. Purpose

W2 converts the ratified operation inventory and W1 path/HTTP grammar into precise request/response schema laws without introducing a second cross-domain ontology for serialization convenience.

> **Reusable wire primitives may normalize representation; business meaning, knowledge, lifecycle, evidence and outcomes remain owner-specific.**

W2 is not derived from legacy OpenAPI, provider DTOs, database rows or frontend forms. It does not choose D6 UI composition, D7 server/generator/blob/persistence/runtime realization, D8 live-effect proof or implementation.

The complete W2 package will receive one coherent independent Fable review after its internal sections converge. Do not fragment W2 into micro-reviews unless a material contradiction requires it.

---

# 2. W2-A — Core Schema Grammar — ACCEPTED IN-STAGE

## 2.1 Governing invariant

> **A Product API schema preserves every material distinction the accepted owner can make, while adding no independent business meaning merely for uniform serialization.**

Corollaries:

1. reusable primitives are value/mechanism primitives, not horizontal business wrappers;
2. request schemas never allow clients to author server-owned authority/history fields;
3. unknown/unavailable/partial/not-applicable semantics are explicit where material and never hidden in `null`, zero, false or empty collections;
4. provider/source qualification remains explicit and mechanically distinguishable;
5. schema reuse must not create a universal entity/evidence/result/operation graph.

## 2.2 MPC-owned identifier grammar

MPC-owned canonical identifiers serialize as opaque non-empty strings.

Named semantic schemas such as `OrganizationId`, `ListingIntentId`, `PriceIntentId`, `FulfillmentNodeId` may exist for documentation/tooling, but their wire value does not promise UUID, ULID, database sequence or any other physical encoding.

Do not expose implementation-derived regex/version/timestamp meaning in the ID contract.

## 2.3 Native external identifier grammar

Provider/business-system native keys serialize as opaque strings even when a current provider happens to use numeric identifiers.

Examples include:

- `native_listing_key`;
- `native_sale_key`;
- `native_shipment_key`;
- `native_product_key`;
- native business-document keys when materially exposed.

The Product API does not invite arithmetic over identifiers, narrow them to JavaScript-safe integers or assume a provider can never evolve an identifier to alphanumeric form.

## 2.4 Typed source-qualified references

No universal `ExternalRef`, `{entity_type, entity_id}`, provider graph or relationship graph is admitted.

Use semantic typed references, for example conceptually:

```json
{
  "marketplace_installation_id": "...",
  "native_listing_key": "MLB123"
}
```

```json
{
  "source_instance_id": "...",
  "native_product_key": "42664"
}
```

A request already scoped by `/organizations/{organization_id}` does not repeat an effective Organization field inside each reference. The server resolves every Organization-owned qualifier/reference inside the path Organization and fails closed on cross-Organization resolution.

## 2.5 Exact decimal primitive and Money

D2 exactness is preserved on the wire with an `ExactDecimalString` primitive.

Baseline lexical properties:

- JSON string, not JSON number;
- ordinary base-10 decimal notation;
- optional leading minus when the owning semantic allows negative values;
- one or more integer digits;
- optional fractional part;
- no exponent notation;
- no binary-floating-point conversion as the authoritative contract representation.

`Money` is:

```json
{
  "amount": "189.90",
  "currency": "BRL"
}
```

where `amount` is `ExactDecimalString` and currency is an explicit ISO-4217-style alpha-3 code in the Product API contract.

No universal minor-unit representation or global `round(2)` rule is introduced. Domain/provider rounding remains owned by the applicable accepted rule.

Rates and exact material quantities may reuse `ExactDecimalString` when exactness matters, while their semantic constraints remain owner-specific. W2 does not create a generic Unit-of-Measure/conversion framework.

## 2.6 Temporal grammar

Material instants use unambiguous date-time strings. Field names preserve the accepted temporal meaning, for example `observed_at`, `recorded_at`, `decided_at`, `effective_at` or an externally authoritative deadline/window field.

A generic `timestamp` or response-generation time never substitutes for source/effective/observation/decision time.

Timezone/local-calendar structures are introduced only where an accepted business/external rule actually depends on them.

## 2.7 Request and response schemas are not one CRUD blob

Create/update request schemas and resource/read response schemas are separate contracts when their authorities differ.

Examples of the pattern:

```text
CreateListingIntentRequest
UpdateListingIntentRequest
ListingIntent
```

Server-owned fields such as canonical ID, actor attribution, decision/history metadata, convergence evidence and current owner state are absent from client-write schemas rather than merely relying on clients to ignore them.

## 2.8 Closed semantic objects

Product semantic request objects are closed by default: undeclared properties are contract errors.

This is a wire-level guard against provider DTO/property-bag leakage and accidental client-authored authority fields.

Response objects are likewise closed where this materially protects the semantic contract; provider enrichment uses explicit bounded variants rather than `additionalProperties` property bags.

The exact JSON Schema/OpenAPI keywords used for closure depend on the later selected OAS/schema version, but the semantic closure requirement is accepted here.

## 2.9 `null` is not a knowledge-state transport

`null` never means "unknown", "unavailable", "partial", "not applicable" or "not observed" by convention.

A nullable value is admitted only when semantic null itself is an accepted value for that specific field. Omission and null are not silently overloaded to carry multiple knowledge meanings.

## 2.10 Knowledge-state grammar

No universal Product API `Fact<T>` envelope is admitted.

Where knowledge state is material, use the smallest owner/field-specific discriminated union needed by that meaning. Potential states are drawn only from those materially applicable, such as:

- known value;
- known empty/absent;
- unknown / insufficiently known;
- unavailable;
- partial/incomplete coverage;
- not applicable / unsupported where semantically distinct.

Different meanings need not expose identical state sets.

A known zero/false/empty value remains distinguishable from unknown/absent when material.

## 2.11 Union validation grammar

Material unions are mechanically exclusive.

Baseline schema technique is conceptually:

```text
oneOf
+ per-variant const discriminator value
```

Use `state` for knowledge/lifecycle-state discriminants and `kind` for semantic type/target unions where practical.

An OpenAPI `discriminator` may later be added as a tooling hint if useful, but correctness does not depend on it; the underlying union schemas must validate independently and exclusively.

## 2.12 Freshness and provenance locality

Freshness/provenance metadata attaches at the smallest semantic unit for which it is uniformly true.

Do not create a universal top-level `metadata`, `Evidence`, `Observation` or generated-at envelope.

Where one owner resource shares one observation context, resource-level provenance may be sufficient. Where components have different sources/times, each component preserves the materially correct provenance/time instead of inheriting a fabricated global timestamp.

## 2.13 No universal result/operation envelope

Do not wrap all Product operations in a generic shape such as:

```json
{ "status": "accepted", "data": {}, "error": null }
```

Resource creation returns the created owner resource; resource mutation/capability returns the current relevant owner resource when that honestly represents the outcome; stateless evaluation returns an operation-specific evaluation schema.

Create `oneOf` outcome variants only when one operation has materially distinct valid outcomes that cannot be represented honestly in the returned owner resource.

No generic `Operation`, `CapabilityResult`, `CommandResult` or universal async tracking identity is introduced.

## 2.14 Business outcomes versus HTTP problems

Valid domain `pending`, `rejected`, `ambiguous`, approval-required or unsupported/external-required outcomes remain semantic response meaning where the owner contract admits them.

They are not automatically HTTP 4xx/5xx and do not become `403` merely because the caller dislikes the outcome.

`202 Accepted` is not a baseline synonym for "external effect has not converged". Owner-local durable Intents already provide the tracking identity. Use `202` only if a later concrete HTTP interaction genuinely has HTTP-level deferred processing that cannot be represented better by the current owner resource/result.

## 2.15 Typed PATCH/update bodies

Baseline Product updates use operation-specific typed JSON request bodies, not JSON Patch and not a generic mutation DSL.

For a typed update request:

- omitted property = unchanged;
- present property carries the complete semantic replacement for that field/relationship according to owner rules;
- `null` is not a generic clear instruction;
- clear/inherit/reset operations are explicit semantic variants where the domain needs them;
- array/list fields, when present, represent the intended complete semantic value unless the operation explicitly defines another owner-specific meaning.

Generic JSON Merge Patch is not adopted by default because `null`, arrays and discriminated knowledge/value semantics can carry material domain meaning.

## 2.16 Problem Details baseline

RFC 9457 Problem Details remains the API-level failure shape.

The problem `type` is the primary stable machine identifier for the problem class. Baseline fields are proportionately `type`, `title`, `status`, `detail`, `instance`.

Do not introduce an independent duplicate top-level `code` taxonomy by default. Add a stable MPC extension only when a real programmatic consumer needs information not already represented by the problem type/extensions.

Validation-detail extensions may expose bounded machine-readable locations/details without provider payloads, secrets or PII.

## 2.17 Provider-enriched evidence schema

Provider-enriched evidence is an owner-local closed variant, not a provider property bag.

Use bounded discriminated variants only for named Product 1.0 consumer/correctness needs. Another provider may honestly expose unsupported/not-applicable/unavailable semantics rather than fabricated equivalence.

Raw provider DTOs, arbitrary provider error text and free-form provider-field maps do not cross the Product API boundary.

## 2.18 Server-authority field fence

Client request schemas never contain effective-authority claims such as:

- effective Principal / `created_by` / `approved_by` as the caller identity;
- effective Organization duplicating the path scope;
- `authorized=true`, `approved=true`, `converged=true` or provider-success assertions;
- server-owned history/convergence/evidence fields.

A Principal/Organization ID may appear only when it is the legitimate business/admin **subject** of that operation, e.g. assigning an AccessRole to another Principal, never as the effective actor/scope authority.

---

# 3. W2-A proof / negative controls

Later executable contract proof must be able to falsify at least:

1. an MPC ID schema constraining clients to a physical UUID/ULID representation without accepted need;
2. a native external key becoming numeric/arithmetic wire state;
3. a bare external key accepted where its Marketplace Installation/SourceInstance qualifier is not otherwise unambiguous;
4. a universal `{entity_type, entity_id}` / `ExternalRef` replacing typed references;
5. Money represented as authoritative JSON floating-point number;
6. a global minor-unit/`round(2)` rule silently imposed on domain/provider economics;
7. unknown/unavailable/partial collapsing into `null`, zero, false or empty collection;
8. an invalid knowledge/type union matching multiple or no legal variants while passing validation;
9. client request containing an undeclared provider field/property bag;
10. one request/response CRUD schema making server-owned authority/history writable;
11. a generic Result/Operation wrapper becoming required Product ontology;
12. generic JSON Patch/merge semantics using `null` or array behavior to bypass owner meaning;
13. response generation time impersonating source/owner observation time;
14. arbitrary provider enrichment/property bag escaping its bounded owner-local variant;
15. request attempting to author effective Principal, Organization scope, authorization or convergence state.

---

# 4. W2-A method outcome

**Parent architecture/W1:** `CURRENT STRUCTURE CONFIRMED`.

**W2-A result:**

> **Use strong reusable scalar and union mechanics with owner-specific semantic schemas. Preserve exact decimals, source qualification, knowledge state, provenance and server authority without exporting a universal Fact/Evidence/Result/Operation/resource wrapper ontology.**

No parent-stage reopen is required.

---

# 5. Exact next W2 work

**W2-B — ListingIntent + PriceIntent + Availability concrete schema grammar.**

W2-B must stress-test W2-A against the hardest Product authoring/control shapes and decide proportionately:

1. `CreateListingIntentRequest`, `UpdateListingIntentRequest` and `ListingIntent`;
2. create-versus-edit target semantics without separate architectures;
3. `FOLLOW_SOURCE | EXPLICIT_OVERRIDE` field/value unions and source-candidate reference grammar;
4. requirement/category/product-type resolution references without provider DTO/property-bag leakage;
5. listing-context authored-media intake/reference/selection/order/role without ProductAsset/media-master authority;
6. PriceIntent correlation from ListingIntent without embedding price value;
7. `PriceIntent` target union: existing source-qualified Listing vs pre-creation ListingIntent context;
8. exact Money target and explicit supersession lineage while mutable PriceDraft stays rejected;
9. `SubmitListingIntent` dispatchability/current resource representation and fail-closed required PriceIntent/Availability correlations for active creation;
10. Sellable Availability read shape, including desired owner value, provider actual evidence, knowledge/freshness and convergence without making provider quantity the owner value;
11. Inventory Source/allocation policy request/update shapes where needed to prove the grammar;
12. negative controls proving ListingIntent cannot accept price/availability/provider payload fields or fake source authority.

W2-B does not choose D7 media/blob/upload transport mechanics. Binary/content acquisition mechanics remain D7; Product semantics own only listing-context authored-media identity/reference/selection.

After W2-B, continue W2 through the remaining admitted owner/schema families necessary for global schema coherence. **Only after W2 as a coherent package converges, run one independent Fable review following the canonical Standard Fable review workflow.**

Implementation remains blocked until D9.
