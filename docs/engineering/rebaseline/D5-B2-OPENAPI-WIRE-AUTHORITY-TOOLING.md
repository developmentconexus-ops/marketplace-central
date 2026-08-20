# D5-B2 — Single OpenAPI Wire Authority / Tooling

> **Status:** ACCEPTED / CANONICAL  
> **Parent stage:** D5-B2 Product Operation / Resource Surface  
> **Canonical inputs:** accepted D0→D4 + D4-R1 + D5-B1 + Operation Admission Matrix + W1 + W2 + W3 + W4 + Technical Ingress + Final Problem/media consistency  
> **Method:** DevelopmentConexus Engineering Method v1.0.0  
> **Lead research opened:** 2026-08-19  
> **Independent Fable review:** COMPLETE  
> **GPT adjudication:** CONVERGED  
> **Round 2:** NOT REQUIRED  
> **Operator final ratification incorporated:** 2026-08-20  
> **Stable Product Problem origin proof:** FILED — `https://conexus.fun`  
> **Canonical Product OAD authoring/proof:** NEXT  
> **Implementation:** BLOCKED until D9 is accepted

## 1. Purpose and authority

This document is the single canonical D5-B2 authority for the Product API's **machine-readable wire authority and bounded tooling**.

It does not replace the semantic homes:

- W1 owns Product resource/path/HTTP grammar;
- W2 owns request/response schemas, Problem Details, idempotency, revision and authored-media grammar;
- W3 owns collection/query/cursor grammar;
- W4 owns Permission, Principal-class and current-access enforcement;
- Technical Ingress owns technical non-Product acquisition/OAuth ingress.

This document decides how those accepted meanings become one machine-readable Product API description and how derived projections are validated and generated without creating a second authority.

It deliberately does **not** choose:

- Go router/server framework;
- Go handler/package topology;
- frontend fetching/cache integration;
- runtime client retry/auth/error middleware;
- production API host/base URL;
- request-validation implementation;
- blob store, proxy, CDN, scanner or transformer;
- technical ingress or authored-media delivery route spelling;
- D6 screen/BFF composition;
- D7 transaction, storage, worker or deployment realization.

---

## 2. Preserved semantic surface

```text
Product operations                              95
ordinary Permissions                            29
Principal kinds                                 H / A / S only
standalone Product media GET/Update/Delete      not admitted
Technical Ingress A/B                           outside Product OAD/SDK
Product path / schema / collection / access     owned by W1 / W2 / W3 / W4
D0→D5-B1 semantic reopen                        none
```

The superseded OpenAPI, handlers, tests and `packages/sdk-runtime` were removed in the operator-ratified runtime reset. They remain available only as Git-history evidence and are not target authority.

### 2.1 2026-08-20 operator sequencing amendment

The checkpointed OA-C13 sequencing originally required the canonical Product OAD replacement to land before retiring the legacy OpenAPI/manual-SDK seam. During the repository cleanup, the operator explicitly ratified a hard cutover first because no production compatibility consumer exists and the superseded runtime/CI topology was creating active-tree and agent-context cost without protecting target behavior.

This amendment changes **retirement sequencing only**:

```text
original: author replacement → retire legacy seam
amended:  retire legacy seam → author canonical replacement
```

It does **not** reopen W1–W4, change any of the 95 Product operations, change the 29 ordinary Permissions, add a Principal kind, authorize D6–D9 or authorize Product implementation. Git history remains the evidence archive for the retired seam.

The current-state failure class is:

> **Two writable wire representations drift while retaining target meaning already rejected by the rebaseline.**

The target design replaces rather than extends that seam.

---

# 3. Canonical decision package

## OA-C1 — One logical OAD is the Product wire authority

The canonical machine-readable Product API wire authority is one OpenAPI Description rooted at:

```text
contracts/api/product/openapi.yaml
```

Every source document reachable from that entry document through repository-local relative `$ref` belongs to the same authority.

Bounded layout:

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

Required properties:

- exactly one Product OAD entry document;
- no independently complete OAD per D1 owner;
- no remote HTTP/Git refs;
- no template/macro/code-first preprocessor;
- no `join` of independent APIs;
- no hand-maintained bundle;
- all component keys globally unique across the source closure;
- bundle inspection rejects collision-generated names such as `Money-2`.

A bounded overlay may exist only as a clearly derived generator/distribution projection for a proven consumer. It never changes source semantics or becomes a second target authority.

## OA-C2 — OAS 3.1.2 is the smallest sufficient feature level

Select exactly:

```yaml
openapi: 3.1.2
jsonSchemaDialect: https://spec.openapis.org/oas/3.1/dialect/base
```

OAS 3.1.2 expresses all accepted W1–W4 properties. No accepted Product meaning requires 3.2, while the selected generator/tooling compatibility is proved at 3.1.2.

The exact value is mechanically pinned; generic lint success is not enough because the selected Redocly release does not independently prove the exact patch value.

Reopen only when a concrete accepted Product contract requires a later-only feature and all selected tools prove it without semantic downgrade or parallel authority.

## OA-C3 — Source, bundle and generated artifacts have distinct status

```text
source OAD                         AUTHORITY
resolved bundle                    DERIVED / TEMPORARY
generated TypeScript               DERIVED
generated Go                       DERIVED
static docs / Problem pages        DERIVED
runtime handlers and traffic       EVIDENCE
```

Pipeline invariant:

```text
lint source
→ verify source closure and global component names
→ deterministic temporary bundle
→ generate Go from bundle
→ generate TypeScript directly from source
→ compile/prove both projections
→ regenerate and prove the intended tree is clean
```

The bundle is mandatory for the selected Go lane but remains uncommitted by default. Commit it only for a named real offline/external consumer; it remains derived.

The Product OAD source tree is excluded from Prettier. Redocly and the OAD-specific authoring/gate rules own its form. This preserves the repository's already-measured failure class where a formatter rewrote the contract and invalidated byte-literal parity checks.

## OA-C4 — `operationId` is exact stable Product operation identity

Every admitted Product operation has one unique PascalCase `operationId`. W4 §8 is the canonical 95-row spelling/mapping source.

Final name-only crystallizations:

```text
GetEffectiveAvailabilityAllocationScopePolicy
UpdateAvailabilityAllocationScopePolicy
ListFulfillmentExecutions
GetFulfillmentExecution
```

The Availability asymmetry is material:

```text
GET
→ configured mode/value
→ effective value
→ effective-source provenance

UPDATE
→ configured owner mode/value only
```

GET and UPDATE do not share one round-trip schema.

The Fulfillment pair only names the already-canonical `FulfillmentExecution` wire home.

A gate parses all 95 W4 rows and diffs the complete `(operationId, class, permission, principal-kinds)` mapping against the OAD. No exception manifest is permitted.

Rules:

- exactly 95 `operationId` values;
- all unique and PascalCase;
- no generator-driven renaming;
- no framework `Controller`/`Handler`, provider name or HTTP-verb suffix by habit;
- a semantic operation change requires explicit architecture adjudication rather than silent `operationId` reuse.

## OA-C5 — Product path/server/version profile

The Product OAD contains only:

```text
/access-context
/organizations/{organization_id}/...
```

plus exact W1-approved paths beneath those roots.

It contains none of:

```text
/v1
/providers
/integrations
/webhooks
oauth callback/begin/refresh
external acquisition ingress
technical authored-media delivery route
provider callback/notification paths
```

The source OAD omits environment-specific `servers`; D7/runtime configuration supplies the base URL. Redocly's `no-empty-servers` rule is explicitly disabled so a tool opinion cannot force a contradiction into the authority.

`info.version` identifies a publication and never implies a `/v1` URI axis.

## OA-C6 — Authentication is standard; MPC Permissions are not OAuth scopes

Define one standard bearer scheme:

```yaml
components:
  securitySchemes:
    MpcBearerAuth:
      type: http
      scheme: bearer
```

Every Product operation, including `GetCurrentAccessContext`, uses that scheme.

The 29 MPC Permissions remain MPC-owned current authorization facts, never OAuth scopes or IdP roles.

Every operation carries bounded projection-only extensions:

```text
x-mpc-operation-class
x-mpc-required-permission
x-mpc-principal-kinds
x-mpc-semantic-owner
```

The two qualified physical checkpoints may carry one bounded qualification marker. These extensions are linted against W4 and cannot widen runtime authority.

## OA-C7 — Closed extension and schema profile

The source OAD follows W2:

- semantic objects declare `additionalProperties: false`;
- material unions use `oneOf` plus required fixed `const` discriminants;
- exact decimals remain strings;
- knowledge states remain explicit;
- provider DTOs do not escape into free-form bags;
- `allOf` is used only for genuine intersection semantics proved by both generators;
- write schemas exclude server-owned fields;
- no universal Product ontology is invented for reuse convenience.

Extension policy is allowlist-based:

```text
allowed
→ the bounded x-mpc-* projection set
→ any later explicitly accepted name-only derived projection

forbidden in source OAD
→ x-go-type
→ x-go-type-import
→ x-go-name
→ x-go-type-skip-optional-pointer
→ any extension that changes generated type or shape
```

`x-enum-varnames` is not admitted now. Generated Go enum-name churn is loud/compiler-visible and remains derived. D7 may reopen only a generator-local naming projection when real cost is proved.

D5 proves declaration of closure and patterns. Runtime rejection of undeclared fields and pattern violations is a named D7 obligation; generated types alone do not discharge it.

## OA-C8 — Authored-media multipart uses the OAS 3.1 form model

`CreateListingIntentMedia` remains:

```text
POST .../listing-intents/{listing_intent_id}:create-media
Idempotency-Key header
multipart/form-data
required raw file part
required text/plain etag part
success body: media + listing_intent_etag
```

Source shape:

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

This is more honest than representing browser `Blob`/`File` bytes as a JSON string merely to obtain prettier generated types.

"Usable projection" means:

- the request is expressible;
- both parts are reachable;
- streaming is possible;
- the parent validator is typed in the success body.

It does **not** mean automatic part binding or runtime validation. Go receives a multipart reader; TypeScript requires FormData assembly; both lanes need bounded runtime code later.

## OA-C9 — Product Problem type URI policy and stable origin

The stable HTTPS origin is:

```text
https://conexus.fun
```

Canonical Product Problem identifiers are:

```text
https://conexus.fun/marketplace-central/problems/product/{slug}
```

Technical protocol/presentation identifiers, when independently contracted, use the disjoint namespace:

```text
https://conexus.fun/marketplace-central/problems/technical/{surface}/{slug}
```

The domain is intentionally shared across DevelopmentConexus projects. Project identity therefore lives in the first path segment rather than consuming the apex for Marketplace Central alone. Future projects may use their own non-overlapping path namespaces without changing Marketplace Central identifiers.

Evidence accepted for the origin proof:

- operator declaration that `conexus.fun` was newly purchased and is controlled for DevelopmentConexus use;
- operator-supplied current browser evidence showing `https://conexus.fun` served over HTTPS with a secure connection and a Hostinger registration/management page;
- Hostinger's documented ability for a registered-domain owner to manage DNS records through hPanel;
- Hostinger's documented HTTPS/SSL support for hosted domains and subdomains.

The previous web content historically associated with the hostname is not authority after the ownership change. What matters for the identifier is current durable control of the registered domain, not the current placeholder content or hosting vendor.

Binding rules:

- Product slugs are canonical W2 kebab-case names;
- `about:blank` remains `about:blank`;
- no version/date segment enters the identifier;
- Product and technical namespaces remain disjoint;
- every custom Product Problem schema constrains `type` to the exact `conexus.fun` URI with `const`;
- generated human-readable pages derive from the same canonical definitions;
- documentation hosting may move between Hostinger, static hosting, a reverse proxy, CDN or another service without changing an already-published identifier;
- the apex domain may host a company site or other projects; Marketplace Central depends only on control/resolution of its own `/marketplace-central/problems/...` namespace;
- no portal, developer platform or central documentation application is required before Product OAD authoring; a minimal static page or permanent redirect is sufficient for first resolution proof;
- loss/transfer of `conexus.fun` control before public Product Problem publication reopens OA-C9 immediately;
- after publication, changing the Problem identifier host is a breaking identity change and is not a routine hosting migration.

The previous ngrok endpoint remains optional preview-only infrastructure:

```text
https://multiradial-unironically-nieves.ngrok-free.dev/
```

It may be used for temporary preview or local tunnel exposure, but never appears in canonical Product `Problem.type` values.

The stable-origin proof gate is **CLOSED / FILED**. Product OAD authoring may now use the exact `conexus.fun` constants above.

## OA-C10 — Exact pinned tooling and routing constraint

Pinned baseline:

```text
@redocly/cli                    2.45.0
openapi-typescript              7.13.0
oapi-codegen                    v2.8.0
github.com/oapi-codegen/runtime v1.7.0
```

Redocly owns source lint, local-ref resolution, deterministic bundle and static-doc compatibility proof.

`redocly.yaml` is **rules-only** and contains no `apis:` block. Every command receives the entry document and output path explicitly so openapi-typescript cannot silently replace its CLI input with config authority.

Required Redocly dispositions include:

- exact `openapi: 3.1.2` assertion;
- `no-empty-servers: off`;
- unresolved-ref and operationId rules;
- W4 extension presence/vocabulary rules;
- extension allowlist/type-override prohibition;
- global component-collision controls.

openapi-typescript generates the mandatory TypeScript contract projection directly from source.

oapi-codegen consumes the resolved temporary bundle. It proves models/strict interfaces but does not select the D7 router or runtime validator.

The standard Go `net/http.ServeMux` is **not** a valid neutral proof vehicle for canonical `{id}:verb` paths. D5 has proved a D7 routing constraint: the runtime must provide partial-segment pattern support or a compatible implementation of the generated mux interface. The D5 fixture uses a bounded custom mux to prove dispatch without selecting Chi, Echo, Gin, Gorilla or another D7 framework.

Compatibility proof must execute both:

```text
exact minimum: GOTOOLCHAIN=go1.25.1 or equivalent CI image
current/newer toolchain: forward-compatibility evidence
```

A newer toolchain honoring the `go 1.25.1` module directive does not by itself prove execution under the exact minimum toolchain.

## OA-C11 — Generated contract projections are mandatory; runtime SDK is D6

```text
TypeScript contract projection from OAD     REQUIRED
Go compatibility projection from OAD        REQUIRED PROOF
hand-authored Product DTO/operation SDK      PROHIBITED
runtime TypeScript client implementation     D6
Go router/server/validation realization      D7
```

A D5 contract SDK exports generated paths, operations and schemas only. It does not manually duplicate endpoint/DTO meaning.

D6 may choose a small typed Fetch transport, a stable generated runtime client or bounded TanStack Query adapters. Paths, params, bodies, responses and Product Problems still derive from the OAD.

The historical `packages/sdk-runtime` remains replacement evidence only; it is neither active nor target authority.

## OA-C12 — Technical surfaces are mechanically excluded

The Product OAD and Product SDK contain none of:

- External Acquisition Ingress;
- OAuth begin/callback/refresh;
- provider notifications/webhooks;
- internal acquisition custody;
- authored-media byte-delivery route;
- provider-protocol Problems;
- object-store/CDN/scanner vocabulary.

Separate technical contracts may later exist for their actual protocols. They are never joined into the Product OAD for documentation convenience.

## OA-C13 — Legacy artifacts remain retired as one measured seam

Target disposition:

```text
contracts/api/marketplace-central.openapi.yaml
→ absent before canonical Product OAD authoring
→ preserved only in Git history; no active compatibility consumer is admitted

packages/sdk-runtime manual Product DTO/operation catalog
→ absent from the active tree
→ generated projections start from the canonical OAD without a migration bridge
→ never present as target Product SDK authority
```

The runtime reset retired the old control set together with its subject population. Canonical OAD authoring owns this forward replacement set:

1. a zero-population guard for the retired OpenAPI/manual-SDK paths;
2. exact 95-operation and 29-Permission checks;
3. deterministic source lint and bundle proof;
4. deterministic TypeScript and Go projection generation;
5. a strict source/generated authority seam;
6. generated Product routing/transport contracts rather than handwritten duplication;
7. contract-drift ratchets tied to accepted W1–W4 meaning;
8. a generated transport boundary whose concrete runtime role is selected by D6/D7.

Useful controls are replaced, not merely deleted:

```text
source OAD
→ deterministic generation
→ generated projections
→ regeneration leaves the intended tree clean
```

D6 inherits the forward constraint that HTTP transport permission follows the selected runtime-client role, not the historical package name `sdk-runtime`.

## OA-C14 — Executable verification and drift controls

Canonical authoring must prove at least:

### Authority and structure

1. one Product OAD entry document;
2. local relative refs only;
3. exact OAS 3.1.2 pin;
4. deterministic temporary bundle;
5. globally unique component keys before bundle;
6. no collision-generated component names in bundle;
7. exactly 95 operations;
8. 95 unique PascalCase `operationId`s;
9. full equality with W4's 95-row mapping;
10. no second complete Product OAD;
11. generated files carry read-only/non-authoritative banners;
12. source OAD extensions are allowlisted.

### W1/W2/W3/W4

13. only approved Product roots and custom `:verb` paths;
14. no `/v1` or technical/provider roots;
15. exact ETag/If-Match and idempotency carrier classes;
16. semantic objects declare closure;
17. `oneOf + const` unions survive both generators;
18. exact decimals remain strings;
19. knowledge states do not collapse to null/default;
20. every custom Product Problem schema uses the exact `https://conexus.fun/marketplace-central/problems/product/{slug}` URI constant;
21. `about:blank`, `405 + Allow` and the media error map remain expressible;
22. multipart file and `etag` are reachable in both lanes;
23. parent ListingIntent validator remains typed result data;
24. all 26 W3 List/Search homes use canonical cursor grammar;
25. no raw provider cursor, offset/page/total/sort DSL;
26. every operation's W4 class, Permission and Principal projection equals W4 authority;
27. exactly 29 ordinary Permission tokens plus the `authenticated` special condition;
28. no OAuth scope becomes Product Permission.

### Generated compatibility

29. Redocly 2.45.0 lint succeeds with checked-in rules-only config;
30. source bundling is deterministic;
31. TypeScript generated from source is deterministic and compiles;
32. Go generated from bundle is deterministic and compiles;
33. custom `:verb` dispatch is proved with a compatible custom mux without selecting D7 framework;
34. exact Go 1.25.1 toolchain proof is executed (`GOTOOLCHAIN=go1.25.1` or equivalent CI image);
35. current/newer Go toolchain proof is also executed as forward-compatibility evidence;
36. code generation is inspected for unsafe source-derived injection;
37. multipart part validation and closed-object runtime rejection remain named D7 obligations, never false D5 PASS claims;
38. regeneration leaves the intended committed projection tree clean.

### Exclusion and retirement

39. Technical Ingress A/B absent;
40. authored-media delivery route absent;
41. technical Problem namespace absent from Product OAD;
42. legacy `/mutations` absent;
43. retired manual SDK types remain absent from the active tree;
44. retired OpenAPI/manual-SDK paths remain at zero population;
45. ngrok or any other temporary host cannot leak into canonical Product Problem constants.

A green command that did not traverse the real source closure, all operations and generated outputs is vacuous.

---

## 4. Review evidence and adjudication

The independent Fable spike proved:

- deterministic Redocly YAML/JSON bundle;
- TypeScript generated from source equals TypeScript generated from bundle byte-for-byte;
- strict TypeScript compilation;
- `const`, `oneOf`, `allOf`, exact decimals, headers, statuses and Problem media type preserved;
- oapi-codegen consumes the resolved bundle and generates compiling Go;
- live `httptest` custom-verb dispatch with a compatible mux;
- multipart file/etag reachability;
- `If-Match`, `Idempotency-Key`, Problem `type` and parent validator round-trip.

The review also proved the limitations incorporated in OA-C1–OA-C14: standard-mux incompatibility, generator-specific type-override divergence, formatter ownership, Redocly config/version gaps, mandatory Go bundle input, component collision renaming, declaration-versus-runtime validation, multipart precision and the exact-toolchain proof gap.

The stable-origin proof was closed with current operator/domain evidence for `conexus.fun`; Hostinger remains a current DNS/hosting mechanism rather than Problem identity authority.

No review finding required D0→D5-B1 or W1→W4 semantic reopen. Round 2 was not required.

---

## 5. Canonical closure and reopen triggers

```text
W1 / W2 / W3 / Technical Ingress               CONFIRMED
W4 semantics                                    CONFIRMED
W4 final operation names                        FILED NAME-ONLY
Product operations                              95 unchanged
ordinary Permissions                            29 unchanged
source model                                    one multi-document design-first OAD
OAS feature level                               3.1.2
lint/bundle                                     Redocly CLI 2.45.0
TypeScript contract generation                  openapi-typescript 7.13.0
Go compatibility generation                     oapi-codegen v2.8.0 + runtime v1.7.0
runtime TypeScript client                       D6
Go router/server/runtime validator              D7
legacy OpenAPI/manual SDK                       RETIRED; canonical replacement authoring next
canonical Product Problem origin                https://conexus.fun
ngrok endpoint                                  PREVIEW ONLY
Product OAD authoring/proof                     NEXT
D6–D9 / implementation                          BLOCKED
```

Reopen only the smallest implicated decision when evidence shows:

- OAS 3.1.2 cannot express an accepted W1–W4 property without loss;
- all credible pinned generators fail a required idiom;
- a real Product consumer requires a runtime SDK property types-only generation cannot support safely;
- control of `conexus.fun` is lost or deliberately transferred before public Product Problem publication;
- a real offline/external consumer requires a committed bundle;
- a real public compatibility consumer requires URI versioning;
- a Product operation cannot be expressed without new meaning;
- a technical surface is mistakenly required by a Product client.

Do not reopen for hosting-vendor changes, generator naming preference, file-count aesthetics, current handler layout, legacy route convenience, desire for `/v1`, or one tool's inability when another credible tool preserves the source semantics.

Current status and exact next action are owned only by `docs/README.md`.