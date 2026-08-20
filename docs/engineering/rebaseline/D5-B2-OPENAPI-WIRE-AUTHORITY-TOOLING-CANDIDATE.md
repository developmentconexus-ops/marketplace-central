# D5-B2 — Single OpenAPI Wire Authority / Tooling Review Candidate

> **Status:** NON-AUTHORITATIVE LEAD DECISION EVIDENCE — INDEPENDENT REVIEW NEXT  
> **Parent stage:** D5-B2 Product Operation / Resource Surface  
> **Canonical inputs:** accepted D0→D4 + D4-R1 + D5-B1 + Operation Admission Matrix + W1 + W2 + W3 + W4 + Technical Ingress + Final Problem/media consistency  
> **Method:** DevelopmentConexus Engineering Method v1.0.0  
> **Opened:** 2026-08-19  
> **Lead disposition:** `RESTRUCTURE NOW — D5-B2-LOCAL OPENAPI AUTHORITY / TOOLING`  
> **Independent review:** required before operator ratification

## 1. Purpose and authority fence

This candidate derives the one machine-readable Product API wire authority and the bounded tooling needed to prove, publish and consume it.

It does **not** amend W1, W2, W3, W4, Technical Ingress or any semantic parent. It does not make the current OpenAPI, current Go handlers or current TypeScript SDK target authority. It does not begin Product OpenAPI authoring.

The decision must settle:

1. the one OpenAPI Description that will become the Product API wire authority;
2. the OpenAPI feature level;
3. source layout, bundle status and generated-artifact status;
4. `operationId` and final wire-operation spelling;
5. Product Problem `type` URI namespace;
6. multipart authored-media representation;
7. bounded lint, bundle and code-generation tooling;
8. what a generated SDK means at D5 versus what remains D6/D7;
9. executable negative controls against drift and parallel authority.

It deliberately does **not** choose:

- Go router/server framework;
- Go handler/package topology;
- frontend data-fetching/cache framework integration;
- runtime client retry/auth/error middleware;
- deployment host/base URL;
- blob store, proxy, CDN, scanner or image transformer;
- provider/OAuth/notification route spelling;
- authored-media byte-delivery route spelling;
- D6 screen/BFF composition;
- D7 process, transaction, storage or deployment realization.

Implementation remains blocked until D9.

---

## 2. Known / inferred / unknown / deferred

### 2.1 Known from repository authority

- Product API contains exactly **95 admitted operations** and **29 ordinary Permissions**.
- W1 owns Product resource/path/HTTP grammar.
- W2 owns request/response schemas, Problem Details, idempotency and revision grammar.
- W3 owns collection/query/cursor grammar.
- W4 owns Permission, Principal-class and current-access enforcement.
- Technical Ingress A/B and authored-media byte delivery are outside Product OpenAPI and Product SDK.
- Product paths have no `/v1` prefix baseline.
- Product custom methods use `POST {resource-or-keyed-subject}:verb`.
- Product semantic objects are closed by default.
- material unions use `oneOf` plus a fixed `const` discriminant.
- exact economic/decision values use decimal strings rather than authoritative JSON floating point.
- `type` is the sole global Product Problem discriminator.
- current host/deployment base path remains D7.
- the current OpenAPI and current SDK are evidence only.

### 2.2 Known from current repository evidence

The current tree contains:

```text
contracts/api/marketplace-central.openapi.yaml
  openapi: 3.1.0
  approximately 247 KB
  legacy /mutations and other superseded routes
  legacy error families and provider/current-code vocabulary

packages/sdk-runtime/src/index.ts
  approximately 87 KB in one primary file
  manually authored Product/provider/current-code interfaces
  no OpenAPI generation script or generator dependency
```

The current Go module declares Go 1.25.1. The repository Node engine is Node 26. These versions can run the modern tooling evaluated below.

This evidence proves the failure class:

> **A large hand-maintained OpenAPI plus a separately hand-maintained TypeScript SDK creates two drifting wire representations, while preserving legacy API meaning that the rebaseline has rejected.**

It does not prove that either current artifact should survive.

### 2.3 Inferred

- One logical OpenAPI Description can remain one authority while using multiple locally referenced source documents.
- A deterministic bundle is useful for external tools, but a bundle must remain derived output.
- Static generated Go/TypeScript projections can prove compatibility without becoming semantic authority.
- Product runtime clients should consume generated contract types, but D5 need not choose D6 transport/cache/error middleware.
- A resolvable HTTPS Product Problem namespace is globally cheaper than a non-resolvable identifier if a controlled stable origin can be established now.

### 2.4 Unknown / deferred

- actual generated Go package and server-adapter location — D7;
- actual frontend runtime client package and TanStack Query composition — D6;
- production API host and deployment base URL — D7;
- exact numeric W3 `limit` default/maximum values until actual ListItem shapes are serialized and measured;
- concrete authored-media accepted types and byte limits until their product/provider constraints are proved before implementation;
- exact documentation deployment mechanism behind the chosen Problem URI origin;
- server request/response runtime validation library — D7.

Unknown/deferred does not authorize a second contract or manual DTO family.

---

## 3. Primary external evidence

The lead used current primary specifications and official project documentation:

- OpenAPI Specification 3.1.2: https://spec.openapis.org/oas/v3.1.2.html
- OpenAPI Specification 3.2.0: https://spec.openapis.org/oas/v3.2.0.html
- RFC 9457 Problem Details: https://www.rfc-editor.org/rfc/rfc9457.html
- Redocly CLI documentation and releases:
  - https://redocly.com/docs/cli
  - https://redocly.com/docs/cli/commands/lint
  - https://redocly.com/docs/cli/guides/lint-and-bundle
  - https://github.com/Redocly/redocly-cli/releases
- oapi-codegen v2.8.0 release and repository:
  - https://github.com/oapi-codegen/oapi-codegen/releases/tag/v2.8.0
  - https://github.com/oapi-codegen/oapi-codegen
- openapi-typescript documentation and v7.13.0 release:
  - https://openapi-ts.dev/
  - https://github.com/openapi-ts/openapi-typescript/releases/tag/openapi-typescript%407.13.0
- current TypeScript runtime-generator disposition:
  - Redocly `generate-client` is experimental;
  - Hey API is in initial development and requires exact version pinning;
  - openapi-fetch is entering maintenance mode under its maintainers' 2026 roadmap;
  - OpenAPI Generator's OpenAPI 3.1 support remains beta.

External tooling informs expression and compatibility. It never owns MPC semantics.

---

## 4. Credible alternatives

### Alternative A — keep one large YAML and the hand-authored `sdk-runtime`

```text
legacy monolithic YAML
+ manually synchronized TypeScript interfaces/functions
```

**Reject.**

It preserves two writable representations, high merge/conflict cost, legacy target meaning and silent drift. The current tree is direct evidence of the failure class.

### Alternative B — Go code-first generation

```text
Go handlers/types/annotations
→ generated OpenAPI
→ generated TypeScript
```

**Reject for target authority.**

It moves D5 wire meaning into D7/implementation shape, makes framework/code layout the contract source and risks losing W2 distinctions that are awkward in Go. Generated OpenAPI from implementation may be useful as drift evidence later, never as the authority.

### Alternative C — independent domain files joined as several APIs

```text
one OpenAPI root per D1 owner
→ join into a Product API
```

**Reject.**

It creates multiple independently complete descriptions and makes a join step an authority-composition mechanism. D1 owners are not public API roots, and Redocly `join` remains experimental.

### Alternative D — one multi-document, design-first OAS 3.1.2 Description

```text
one entry document
+ locally referenced path/component documents
→ one logical OpenAPI Description
→ deterministic validation/bundle
→ generated Go/TypeScript projections
```

**Select as the Global Maximum.**

It preserves one authority while avoiding a 247 KB editing hotspot; keeps implementation downstream; uses standard `$ref`; supports W1–W4; and permits multiple generators to consume one resolved contract.

### Alternative E — select OAS 3.2 immediately

**Defer safely.**

OAS 3.2 is current, but no W1–W4 property requires a 3.2-only feature. Documentation/generation support remains uneven, including current Redoc tooling that still limits static docs to 3.0/3.1. Selecting 3.2 now adds compatibility risk without product value.

---

# 5. Proposed converged decision package

## OA-C1 — One logical OAD is the Product wire authority

The canonical machine-readable Product API wire authority will be one OpenAPI Description rooted at:

```text
contracts/api/product/openapi.yaml
```

All Product API source documents reachable from that entry document through local relative `$ref` are parts of the same authority.

Recommended bounded layout:

```text
contracts/api/product/
  openapi.yaml
  paths/
  components/
    schemas/
    parameters/
    headers/
    request-bodies/
    responses/
    security/
```

The exact editorial split inside those directories may evolve without changing authority, provided:

- there is exactly one entry document;
- every referenced source belongs to that entry document;
- every reference is repository-local and relative;
- no second independently complete Product OAD exists;
- no source file is also manually maintained in bundled/generated form.

The OAS specification explicitly defines a multi-document OAD as one entry document plus its referenced documents. File count therefore does not create authority count.

### Prohibited source machinery

Do not add:

```text
remote HTTP $ref
Git branch/tag $ref
YAML anchors as cross-file schema mechanism
template/macro language
code-first preprocessor
OpenAPI Overlay as standing target authority
join of independently complete domain APIs
```

A bounded overlay may later be used only to create an explicitly derived distribution projection for a real external consumer. It never changes the canonical source OAD.

---

## OA-C2 — OAS 3.1.2 is the smallest sufficient feature level

Select:

```yaml
openapi: 3.1.2
jsonSchemaDialect: https://spec.openapis.org/oas/3.1/dialect/base
```

Reasons:

1. W2 closed objects, `oneOf`, `const`, exact string patterns and JSON Schema 2020-12 semantics fit OAS 3.1.
2. current Redocly, Redoc, openapi-typescript and oapi-codegen release lines have usable 3.1 support;
3. OAS 3.1.2 is a clarification patch within the 3.1 feature set;
4. OAS 3.2 adds no currently required Product meaning;
5. current docs/generator support for 3.2 is not yet uniformly proven.

Reopen the minor only when:

```text
a concrete accepted Product contract requires a 3.2-only feature
+ selected lint/docs/Go/TypeScript tooling proves that feature
+ no compatibility downgrade/parallel authority is introduced
```

Do not select by recency alone.

---

## OA-C3 — Source, bundle and generated artifacts have distinct status

```text
source OAD (entry + local refs)     AUTHORITY
resolved/bundled OpenAPI            DERIVED DISTRIBUTION
generated Go projection             DERIVED
generated TypeScript projection     DERIVED
static HTML / Problem docs           DERIVED
current runtime traffic/handlers     EVIDENCE
```

Rules:

- source OAD is edited;
- bundle/generated artifacts are never hand-edited;
- deterministic regeneration must produce no diff;
- generated banners identify source path, generator and pinned version;
- no generated artifact may silently add/remove operations or change schemas;
- a tool inability to consume the source OAD is a compatibility finding, not permission to make a second downgraded authority.

The default baseline does **not** commit a bundled OpenAPI copy. Bundle to a temporary/build artifact for tools and distribution. Commit a bundle only if a real offline consumer requires it; it remains derived and regeneration-checked.

---

## OA-C4 — `operationId` is exact stable Product operation identity

Every admitted Product operation has one unique PascalCase `operationId`.

The canonical spelling source entering OpenAPI authoring is the W4 95-row operation mapping, with these final D5-local crystallizations:

```text
get effective allocation/scope policy
→ GetAvailabilityAllocationScopePolicy

update allocation/scope policy
→ UpdateAvailabilityAllocationScopePolicy

ListFulfillmentStates
→ ListFulfillmentExecutions

GetFulfillmentState
→ GetFulfillmentExecution
```

The two Fulfillment changes only name W2's already-canonical wire home, `FulfillmentExecution`; they add no operation or semantic owner.

All other admitted operations preserve their accepted W4 operation name exactly.

Rules:

- exactly 95 `operationId` values;
- all unique;
- PascalCase;
- no generator-driven renaming;
- no `Controller`, `Handler`, provider name or HTTP verb suffix by framework habit;
- operation identity remains stable across path refactoring that preserves semantics;
- a semantic operation change receives explicit architecture adjudication rather than silent `operationId` reuse.

Generated language methods may follow language casing conventions, but the OAD `operationId` remains unchanged.

---

## OA-C5 — Product path/server/version profile

The OAD contains only:

```text
/access-context
/organizations/{organization_id}/...
```

plus the exact W1-approved Product paths beneath those roots.

It contains none of:

```text
/v1
/providers
/integrations
/webhooks
/oauth
/external-events
technical authored-media delivery route
provider callback/notification paths
```

The source OAD omits environment-specific `servers`. Clients receive the base URL through deployment/runtime configuration. This preserves W1's D7-owned host/base path.

`info.version` identifies the OAD publication, not a URI compatibility axis. It never implies a `/v1` prefix. Draft authoring may use a clearly non-production prerelease; the accepted Product 1.0 contract uses a deliberate release value.

---

## OA-C6 — Authentication is standard; MPC Permissions are not OAuth scopes

Define one Product bearer scheme:

```yaml
components:
  securitySchemes:
    MpcBearerAuth:
      type: http
      scheme: bearer
```

Every Product operation, including `GetCurrentAccessContext`, requires that scheme.

Do not use an `oauth2` Security Scheme with the 29 MPC Permissions as scopes. Issuer, authorization URL, token URL, Keycloak realm and machine credential realization remain D7/IdP concerns.

To keep W4 visible and mechanically reviewable, each operation carries a bounded extension projection:

```text
x-mpc-operation-class
x-mpc-required-permission
x-mpc-principal-kinds
x-mpc-semantic-owner
```

For the two system-admitted physical checkpoints, add one bounded qualification marker rather than a generic machine-capability vocabulary.

These extensions:

- project W4 into the wire artifact;
- do not create new Permissions, Principal kinds or Governance semantics;
- are not OAuth scopes;
- cannot widen runtime authority;
- must be changed atomically with the corresponding canonical W4 decision;
- are linted for presence and allowed values on all 95 operations.

`GetCurrentAccessContext` uses the special required-access value `authenticated`, which is not one of the 29 stored Permissions.

---

## OA-C7 — OAS schema profile preserves W2 meaning

The source OAD follows this profile:

- objects with semantic boundaries use `additionalProperties: false`;
- undeclared property bags are rejected;
- material unions use `oneOf`;
- every variant has a required fixed `kind` or `state` using `const`;
- an OpenAPI `discriminator` may be added only for generator ergonomics and never replaces validation;
- `null` appears only where semantic null is legitimate;
- unknown/unavailable/partial/not-applicable remain explicit variants;
- exact decimals and Money amounts are strings;
- MPC and native identifiers are opaque strings;
- temporal names preserve their real meaning;
- write schemas exclude server-owned fields rather than depending only on `readOnly`;
- no universal `Result<T>`, `Fact<T>`, `ExternalRef`, `Operation`, `Evidence`, `Policy`, `Workflow` or metadata envelope;
- no free-form `additionalProperties` escape hatch for provider DTOs;
- `allOf` is used only when intersection semantics are genuinely correct and all selected generators prove it;
- component reuse never creates a business ontology absent from W2.

The OAD may use reusable wire primitives. Business schemas remain owner-specific.

---

## OA-C8 — Authored-media multipart uses the OAS 3.1 form model

`CreateListingIntentMedia` is represented as:

```text
POST .../listing-intents/{listing_intent_id}:create-media
Idempotency-Key header
multipart/form-data request body
one required file property
one required typed etag property
bounded additional metadata only when proven
```

OAS 3.1's native multipart shape is preferred:

```yaml
content:
  multipart/form-data:
    schema:
      type: object
      additionalProperties: false
      required: [file, etag]
      properties:
        file: {}
        etag:
          type: string
    encoding:
      file:
        contentType: <accepted media types when proved>
      etag:
        contentType: text/plain
```

Why:

- OAS 3.1 models raw binary outside JSON `type`;
- multipart binary is not base64;
- `encoding.file.contentType` owns per-part media families;
- `contentMediaType` is not duplicated when an Encoding Object owns the part content type;
- no arbitrary URL, UploadSession or generic Asset appears.

Before canonical OpenAPI authoring, executable compatibility proof must show how each selected generator exposes the `file` part. A generator may use `multipart.Reader`/streaming instead of a typed byte field; that is acceptable. The OAD must not regress to base64 or lie about raw binary merely to obtain a prettier generated type.

The success response carries:

```text
media
listing_intent_etag
```

as typed body data. It does not put the parent validator in the `ETag` header of the distinct custom-method target.

---

## OA-C9 — Product Problem type URIs use one controlled HTTPS namespace

Select the currently repository-evidenced controlled origin:

```text
https://developmentconexus-ops.github.io/marketplace-central/
```

Product custom types:

```text
https://developmentconexus-ops.github.io/marketplace-central/problems/product/{slug}
```

Technical protocol/presentation types, when their own executable contracts are authored:

```text
https://developmentconexus-ops.github.io/marketplace-central/problems/technical/{surface}/{slug}
```

Rules:

- Product slugs are the canonical W2 kebab-case names;
- `about:blank` remains `about:blank`;
- no version/date segment enters the type URI;
- Product and technical namespaces are disjoint;
- technical types never enter the Product OAD;
- each Product type schema constrains `type` to its exact URI using `const`;
- human-readable pages are generated from the same canonical definitions;
- before production exposure, each locator must resolve to documentation or a stable redirect;
- a future custom documentation domain does not replace the established type URI. It may host documentation behind a permanent redirect while the original identity remains stable.

Why this origin:

- control of the GitHub organization/repository is current repository evidence;
- no custom DevelopmentConexus documentation domain is proven in current authority;
- RFC 9457 encourages resolvable absolute URIs because changing from a non-resolvable identifier later changes problem identity.

A reviewer should challenge GitHub Pages lifecycle/control and propose a better already-proven HTTPS origin if one exists. Do not substitute an unverified custom domain.

---

## OA-C10 — Exact pinned validation/generation baseline

Proposed minimum toolchain:

```text
@redocly/cli            2.45.0
openapi-typescript      7.13.0
oapi-codegen            v2.8.0
```

### Redocly CLI

Use for:

- standards validation;
- project-specific lint rules;
- multi-document reference resolution;
- deterministic bundling;
- static documentation compatibility proof.

Pin exactly. Do not use `@latest`.

### openapi-typescript

Use for the mandatory generated TypeScript contract projection:

```text
paths
components
operations
```

Pin exactly. Generated output is read-only and TypeScript-compiled.

### oapi-codegen

Use for bounded Go compatibility proof and, later, generated Go models/strict operation interfaces if D7 confirms the server integration.

Pin exactly. The current v2.8.0 line is the first with initial OAS 3.1 support, so the proof is load-bearing. Generated code must be reviewed; the source OAD must be validated before generator execution.

This tooling choice does **not** select a Go HTTP router. The compatibility fixture may use `std-http-server`/strict generation to avoid deciding Chi/Gin/Echo.

### Not selected as baseline

- Redocly TypeScript client generator — experimental;
- Redocly `join` — experimental and wrong authority shape;
- openapi-fetch/openapi-react-query — entering maintenance mode;
- Hey API runtime SDK — active but pre-1.0/initial development, exact-pin candidate for D6 only;
- OpenAPI Generator — heavier Java toolchain and 3.1 support not yet the selected proven path;
- code-first Go generators — reverse the target authority.

Reopen the tool choice when a pinned release fails a required contract fixture or a materially better stable tool lowers whole-system cost without changing OAD semantics.

---

## OA-C11 — Generated contract projections are mandatory; runtime SDK is D6

The D5 disposition is:

```text
TypeScript contract projection from OAD     REQUIRED
Go compatibility projection from OAD        REQUIRED PROOF
hand-authored Product DTO/operation SDK      PROHIBITED
runtime TypeScript client implementation     D6
Go server framework/adapter realization      D7
```

A generated TypeScript package is the Product API **contract SDK** at this stage. It exports generated path/operation/schema types and no manually authored parallel wire DTOs.

D6 may choose:

- a small generic Fetch transport over generated `paths`;
- a stable generated runtime client;
- bounded TanStack Query adapters.

Whatever D6 chooses:

- endpoint paths, params, bodies, responses and Product Problems derive from this OAD;
- no manually duplicated endpoint/DTO catalog is allowed;
- cache/query keys and UI composition remain D6 meaning, not OpenAPI authority;
- generated runtime output remains derived and regeneration-checked.

The current `packages/sdk-runtime` is legacy current-state evidence. Target Product API work must not extend its manually authored interfaces. During canonical OpenAPI realization it must be retired, replaced or explicitly bounded away from the target Product API so it cannot remain a second authority.

---

## OA-C12 — Technical surfaces are mechanically excluded

The Product OAD contains none of:

- External Acquisition Ingress;
- OAuth begin/callback/refresh;
- provider notifications/webhooks;
- internal typed acquisition custody;
- authored-media byte-delivery route;
- provider-protocol Problem types;
- object-store/CDN/scanner route/error vocabulary.

Separate technical contracts may later exist for their actual protocols. They are not joined, bundled or exposed as Product OAD components merely for documentation convenience.

The Product SDK cannot generate access to those surfaces.

---

## OA-C13 — Legacy artifacts have a replacement, not compatibility, disposition

When the canonical Product OAD is authored and accepted:

```text
contracts/api/marketplace-central.openapi.yaml
→ remove from target active authority path
→ preserve only in Git history unless a real legacy runtime compatibility consumer is explicitly named

packages/sdk-runtime manually authored Product API types
→ stop extending immediately
→ replace target Product API imports with generated contract types
→ retain only genuinely non-Product/legacy modules during a bounded migration
→ never present the package as target Product SDK authority
```

No compatibility consumer currently justifies preserving legacy `/mutations`, generic Integration/provider types or old error schemas in the new Product OAD.

A later migration plan may keep old runtime routes temporarily. That is implementation/current-state migration, not target contract authority and not a reason to contaminate the new OAD.

---

## OA-C14 — Verification and drift controls

The canonical authoring sub-batch must provide executable proof for at least:

### Authority and structure

1. exactly one Product OAD entry document;
2. every `$ref` resolves locally;
3. no remote references;
4. deterministic bundle from a clean checkout;
5. bundle contains exactly 95 operations;
6. 95 unique PascalCase `operationId` values;
7. exact four crystallized names in OA-C4;
8. no unapproved operation/path;
9. no second independently complete Product OAD;
10. generated files contain a non-authoritative/read-only banner.

### W1

11. only `/access-context` and Organization Product roots;
12. no `/v1`;
13. no provider/integration/technical roots;
14. custom methods preserve exact `:verb` spelling;
15. required ETag/If-Match carrier class is represented correctly;
16. `Idempotency-Key` appears only where ratified.

### W2

17. closed semantic objects;
18. oneOf+const unions are mechanically exclusive;
19. exact decimals are strings;
20. unknown/unavailable/not-applicable do not collapse to null/default;
21. request/read schemas preserve server-authority fences;
22. all canonical Product Problem types have exact URI `const`;
23. `about:blank` failures remain status-driven;
24. `405` describes `Allow`;
25. media 400/413/415/422/409/500 map is expressible;
26. multipart file + typed ETag generates usable Go and TypeScript projections;
27. parent ListingIntent ETag is returned as typed result data;
28. no provider/storage/scanner errors leak into Product Problems.

### W3

29. all 26 List/Search homes use the canonical cursor envelope/fields;
30. cursor semantic query fields remain explicit;
31. no raw provider cursor;
32. no offset/page/total/sort DSL by symmetry;
33. numeric limit values are not invented before payload measurement.

### W4

34. every operation has one allowed operation class projection;
35. every Organization operation has one exact required Permission projection;
36. `/access-context` uses the special authenticated condition;
37. projection uses exactly the 29 ordinary Permissions;
38. Principal-kind values are H/A/S meaning only;
39. only the two proven physical checkpoints admit qualified system evidence;
40. no OAuth scope becomes Product Permission.

### Generated compatibility

41. Redocly 2.45.0 lint succeeds under a checked-in strict configuration;
42. Redocly 2.45.0 bundle succeeds and is deterministic;
43. openapi-typescript 7.13.0 output compiles under repository TypeScript;
44. oapi-codegen v2.8.0 consumes the resolved 3.1.2 OAD;
45. generated Go proof compiles under Go 1.25.1;
46. oneOf+const, closed objects, multipart, Problem responses and ETag/header cases survive generation;
47. regenerated committed projections leave the tree clean;
48. generated code is inspected for unsafe source-derived code injection before execution.

### Exclusion

49. Technical Ingress A/B absent;
50. authored-media byte-delivery route absent;
51. technical problem namespace absent from Product OAD;
52. current legacy `/mutations` absent;
53. current hand-authored SDK types are not imported as target wire authority.

A green command that did not traverse the source OAD, all references, all operations and generated outputs is vacuous.

---

## 6. Lead compatibility matrix

| Property | OAS 3.1.2 | Redocly 2.45.0 | openapi-typescript 7.13.0 | oapi-codegen 2.8.0 | Disposition |
|---|---:|---:|---:|---:|---|
| multi-document local `$ref` | normative | supported | resolved/bundle input | bundle input preferred | use one rooted OAD |
| JSON Schema 2020-12 | normative | supported | supported | initial 3.1 support | spike required |
| `oneOf` + `const` | normative | lintable | supported | release claims support | spike required |
| closed objects | normative | lintable | supported | generated behavior must be inspected | spike required |
| multipart file + typed text part | normative | lintable | type output must be inspected | first-class multipart claim | spike required |
| Problem `type` const URI | normative | lintable | supported | output inspection | spike required |
| strict Go operation interfaces | N/A | N/A | N/A | supported, no input validation | D7 realization deferred |
| TypeScript runtime client | N/A | experimental generator available | types only | N/A | D6 deferred |

The lead does not claim compatibility merely from project marketing. The independent review must execute a bounded fixture.

---

## 7. Independent review questions

Fable should challenge this package as one system, especially:

1. Does a multi-document OAD truly remain one authority under our editing/gate model?
2. Is OAS 3.1.2 the smallest sufficient version, or does one W1–W4 property force 3.2/3.0?
3. Can Redocly 2.45.0 lint and bundle the required 3.1.2 idioms deterministically?
4. Can oapi-codegen v2.8.0 generate compileable Go for:
   - `oneOf` + `const`;
   - `additionalProperties: false`;
   - RFC 9457 per-type URI constants;
   - multipart raw file + typed `etag`;
   - multiple success/error responses;
   - custom-method paths;
   - ETag/If-Match and Idempotency-Key headers?
5. Can openapi-typescript 7.13.0 preserve those types without widening to unsafe `any`?
6. Does selecting only generated TS contract types at D5 leave a hidden manual SDK authority, or correctly defer runtime transport to D6?
7. Is a different stable runtime generator already globally cheaper than the defer?
8. Is the GitHub Pages Problem origin genuinely controlled/stable enough, or is another **already-proven** HTTPS origin better?
9. Do the four final `operationId` spellings preserve semantics and generated usability?
10. Are the `x-mpc-*` projections necessary and bounded, or do they duplicate W4 authority?
11. Can the old OpenAPI/SDK be retired without an unnamed compatibility consumer?
12. Do any negative controls choose D6/D7 implementation prematurely?
13. Which hardest foreseeable change — second client, new provider, API split, OAS upgrade, custom domain, or generator replacement — exposes a local rather than Global Maximum?
14. Is any D0→D5-B1/W1–W4 semantic reopen actually unavoidable?

---

## 8. Lead disposition

```text
W1 / W2 / W3 / W4 / Technical Ingress         CURRENT STRUCTURE CONFIRMED
Product operations                              95 unchanged
ordinary Permissions                            29 unchanged
OpenAPI authority/tooling                       RESTRUCTURE NOW — D5-B2 LOCAL
selected source model                           one multi-document design-first OAD
selected OAS feature level                      3.1.2
selected lint/bundle baseline                   Redocly CLI 2.45.0
selected TypeScript contract generator          openapi-typescript 7.13.0
selected Go compatibility generator             oapi-codegen v2.8.0
runtime TypeScript client                        DEFER TO D6
Go server framework                              DEFER TO D7
Problem type namespace                          GitHub Pages Product/technical split
legacy OpenAPI/manual SDK                        REPLACE / RETIRE, not target authority
independent compatibility/adversarial review     NEXT
operator ratification                           BLOCKED BY REVIEW + ADJUDICATION
Product OpenAPI authoring                        BLOCKED BY RATIFICATION
D6–D9 / implementation                          BLOCKED
```

No parent semantic reopen is proven.

---

## 9. Reopen triggers

Reopen only the smallest implicated decision when evidence shows:

- OAS 3.1.2 cannot represent a canonical W1–W4 property without semantic loss;
- all credible pinned Go or TypeScript tools fail a required schema idiom;
- a real Product consumer requires a runtime SDK property that types-only generation cannot support safely;
- the selected Problem URI origin cannot remain controlled/stable or resolvable before production exposure;
- a real external/offline consumer requires a committed bundle;
- a real public compatibility consumer requires URI/API versioning;
- a Product operation cannot be expressed without a new semantic decision;
- a technical surface is mistakenly required by a Product client and must reopen its smallest B2 admission boundary.

Do not reopen for generator naming preference, file-count aesthetics, current handler layout, legacy route convenience, desire for `/v1`, or a tool's inability to preserve semantics when another credible tool can.
