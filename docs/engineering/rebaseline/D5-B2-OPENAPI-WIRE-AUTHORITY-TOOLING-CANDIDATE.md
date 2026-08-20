# D5-B2 — Single OpenAPI Wire Authority / Tooling Review Candidate

> **Status:** NON-AUTHORITATIVE CONVERGED ADJUDICATION EVIDENCE — OPERATOR FINAL RATIFICATION NEXT  
> **Parent stage:** D5-B2 Product Operation / Resource Surface  
> **Canonical inputs:** accepted D0→D4 + D4-R1 + D5-B1 + Operation Admission Matrix + W1 + W2 + W3 + W4 + Technical Ingress + Final Problem/media consistency  
> **Method:** DevelopmentConexus Engineering Method v1.0.0  
> **Opened:** 2026-08-19  
> **Independent Fable review:** COMPLETE  
> **GPT adjudication:** CONVERGED  
> **Round 2:** NOT REQUIRED  
> **Operator origin-policy decision:** RECEIVED; bounded as described in OA-C9  
> **Operator final package ratification:** NEXT

## 1. Purpose and authority fence

This candidate consolidates the lead research, Fable independent adversarial review, executable compatibility spike and GPT adjudication for the one machine-readable Product API wire authority and its bounded tooling.

It remains **non-authoritative evidence**. It does not amend W1, W2, W3, W4, Technical Ingress or any semantic parent. It does not make the legacy OpenAPI, current Go handlers or current TypeScript SDK target authority. It does not authorize Product OpenAPI authoring before final operator ratification and canonical filing.

The converged package settles:

1. one logical Product OpenAPI Description;
2. OpenAPI feature level;
3. source/bundle/generated-artifact authority status;
4. exact operation-name law;
5. Product Problem URI namespace policy and the remaining origin proof gate;
6. authored-media multipart representation;
7. pinned lint, bundle and generation tooling;
8. D5 contract SDK versus D6/D7 runtime realization;
9. executable negative controls against drift and duplicate authority;
10. atomic retirement scope for the legacy OpenAPI/manual SDK seam.

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

Implementation remains blocked until D9.

---

## 2. Evidence state

### 2.1 Repository authority preserved

- Product API remains exactly **95 admitted operations**.
- Ordinary Permission vocabulary remains exactly **29**.
- W1 owns Product path/HTTP grammar.
- W2 owns schemas, Problem Details, idempotency and revision grammar.
- W3 owns collection/query/cursor grammar.
- W4 owns Permission, Principal class and current-access enforcement.
- Technical Ingress A/B and authored-media byte delivery remain outside Product OpenAPI and Product SDK.
- Product paths remain without `/v1`.
- Product custom methods remain `POST {resource-or-keyed-subject}:verb`.
- `type` remains the sole global Product Problem discriminator.
- current OpenAPI, code, tests and `sdk-runtime` remain current-state evidence only.

No D0→D5-B1 or W1→W4 semantic reopen was proved.

### 2.2 Current-state failure class

The current tree contains a large hand-maintained OpenAPI and a separately hand-maintained TypeScript SDK. They preserve legacy `/mutations`, provider/current-code vocabulary and old error families.

The failure class is:

> **Two writable wire representations drift while retaining target meaning already rejected by the rebaseline.**

The target design therefore replaces, rather than extends, that seam.

### 2.3 Fable executable spike

Fable exercised a temporary OAS 3.1.2 multi-document fixture with exactly pinned tools:

```text
@redocly/cli                     2.45.0
openapi-typescript               7.13.0
typescript                       5.9.3
oapi-codegen                     v2.8.0
github.com/oapi-codegen/runtime  v1.7.0
repository Go directive          1.25.1
executed Go toolchain            1.26.4
repository Node                  26.3.0
```

The fixture jointly exercised:

- local multi-document refs;
- ordinary Organization GET;
- custom `POST ...:verb`;
- `If-Match` and `Idempotency-Key`;
- closed-object declarations;
- `oneOf + const` unions;
- exact decimal strings;
- RFC 9457 base and custom Problems with exact HTTPS `type` constants;
- `about:blank`;
- `405 + Allow`;
- multipart raw file + typed `etag` part;
- media success carrying the parent ListingIntent validator;
- W3-shaped cursor fields;
- bounded `x-mpc-*` extensions.

### 2.4 Proved compatible

- one multi-document OAD remains one closure;
- Redocly lint and deterministic YAML/JSON bundling;
- source and bundled TypeScript generation produce byte-identical output;
- TypeScript generation is deterministic and compiles under strict mode;
- `const`, `oneOf`, `allOf`, exact decimal strings, headers, statuses and Problem media type survive generation;
- oapi-codegen consumes the resolved bundle and generates compiling Go;
- custom `:verb` dispatch, multipart read, `If-Match`, `Idempotency-Key`, Problem `type` and parent validator were exercised through `httptest` with a compatible mux;
- no parent semantic reopen is necessary.

### 2.5 Proved limitations

- the standard Go `net/http.ServeMux` registration generated by `std-http-server` cannot represent a wildcard followed by `:verb` in the same path segment;
- `x-go-type`/`x-go-name` can silently produce mutually contradictory Go and TypeScript contracts;
- the repository Prettier scope would reach the proposed nested OAD tree unless explicitly excluded;
- Redocly defaults conflict with the deliberate omission of `servers`;
- Redocly does not by default prove the exact OAS patch version;
- `openapi-typescript` can silently use a discovered `redocly.yaml` `apis:` block instead of the CLI input;
- oapi-codegen v2.8.0 requires the resolved bundle for this multi-document shape;
- component-key collisions are silently renamed by the bundler;
- `additionalProperties: false`, string patterns and multipart part typing are declarations, not runtime enforcement;
- the executed Go proof used Go 1.26.4 with module directive 1.25.1, not the exact 1.25.1 toolchain.

---

# 3. Converged decision package

## OA-C1 — One logical OAD is the Product wire authority

The canonical machine-readable Product API wire authority will be one OpenAPI Description rooted at:

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
- all component keys globally unique across the entire source closure;
- bundle inspection rejects collision-generated names such as `Money-2`.

A bounded overlay may exist only as a clearly derived generator/distribution projection for a proven consumer. It never changes source semantics or becomes a second target authority.

## OA-C2 — OAS 3.1.2 is the smallest sufficient feature level

Select exactly:

```yaml
openapi: 3.1.2
jsonSchemaDialect: https://spec.openapis.org/oas/3.1/dialect/base
```

OAS 3.1.2 expresses all currently accepted W1–W4 properties. No accepted Product meaning requires 3.2, while the selected generator/tooling compatibility is proved at 3.1.2.

The exact value is mechanically pinned; lint success alone is not proof because Redocly accepts unsupported patch spellings unless configured otherwise.

Reopen only when a concrete accepted Product contract requires a 3.2-only feature and all selected tools prove it without semantic downgrade or parallel authority.

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

The Product OAD source tree is excluded from Prettier. Redocly and the authoring/gate rules own its shape. This preserves the repository's already-measured failure class where a formatter rewrote the contract and invalidated byte-literal parity tests.

## OA-C4 — `operationId` is exact stable Product operation identity

Every admitted Product operation has one unique PascalCase `operationId`. W4 §8 is the canonical 95-row spelling/mapping source after final filing.

Final name-only crystallizations:

```text
get effective allocation/scope policy
→ GetEffectiveAvailabilityAllocationScopePolicy

update allocation/scope policy
→ UpdateAvailabilityAllocationScopePolicy

ListFulfillmentStates
→ ListFulfillmentExecutions

GetFulfillmentState
→ GetFulfillmentExecution
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

After ratification, these four names are filed into W4 §8 as a name-only correction. Then one gate parses all 95 W4 rows and diffs the complete `(operationId, class, permission, principal-kinds)` mapping against the OAD. No four-row exception manifest is permitted.

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

## OA-C9 — Product Problem type URI policy and origin proof gate

Canonical custom Product Problem identifiers require one immutable HTTPS origin:

```text
<stable-origin>/marketplace-central/problems/product/{slug}
```

Technical protocol/presentation identifiers, when independently contracted, use a disjoint namespace:

```text
<stable-origin>/marketplace-central/problems/technical/{surface}/{slug}
```

Rules:

- Product slugs are canonical W2 kebab-case names;
- `about:blank` remains `about:blank`;
- no version/date segment enters the identifier;
- Product and technical namespaces remain disjoint;
- every Product Problem schema constrains `type` to the exact URI with `const`;
- generated human-readable pages derive from the same definitions;
- the chosen origin must remain controlled after hosting or repository changes;
- changing documentation hosting never changes an already-published Problem identity.

### Operator-selected ngrok endpoint

The operator nominated this zero-cost endpoint:

```text
https://multiradial-unironically-nieves.ngrok-free.dev/
```

Official ngrok plan documentation establishes:

- the free plan provides one **automatically assigned development domain** specific to the account;
- free users cannot choose/customize the hostname;
- free users cannot bring their own domain;
- development endpoints do not have an endpoint timeout, but remain dependent on the ngrok account/service and plan.

Therefore the nominated endpoint is accepted only as:

```text
TEMPORARY PRE-PRODUCTION DOCUMENTATION / RESOLUTION PROOF ENDPOINT
```

It is **not** classified as a DevelopmentConexus-owned verified domain and is not embedded into canonical `Problem.type` values merely "for now". Problem identifiers are permanent; a temporary hostname cannot be swapped later without changing identity.

Before Product OAD authoring may embed custom Product Problem `type` constants, a proof must establish one of:

1. a DevelopmentConexus-owned and verified HTTPS domain; or
2. an explicitly operator-accepted permanent third-party-domain residual, including the obligation never to relinquish/rename the controlling account or hostname.

Until that proof is filed:

```text
ngrok may host generated preview docs
custom Product Problem type constants remain BLOCKED
about:blank remains usable where already canonical
Product OpenAPI authoring remains blocked at the Problem-origin gate
```

Primary evidence:

- https://ngrok.com/docs/pricing-limits/free-plan-limits
- https://ngrok.com/docs/pricing-limits/how-ngrok-charges
- RFC 9457: https://www.rfc-editor.org/rfc/rfc9457.html

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

The current `packages/sdk-runtime` is a replacement target, not target authority.

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

## OA-C13 — Legacy artifacts retire as one measured seam

Target disposition:

```text
contracts/api/marketplace-central.openapi.yaml
→ remove from target active tree when canonical Product OAD replaces it
→ preserve in Git history unless a real legacy production consumer is named

packages/sdk-runtime manual Product DTO/operation catalog
→ stop extending
→ replace target Product imports with generated contract types
→ retain only bounded genuinely legacy/non-Product modules during migration
→ never present as target Product SDK authority
```

The replacement must own the complete mechanical retirement set, not discover it as red gates after deletion:

1. `GOV_API_SDK_SPLIT` implementation;
2. `api-sdk-atomicity` invariant;
3. `openapi-without-sdk` eval;
4. legacy OpenAPI knowledge route;
5. legacy `modules.json` `openapi_prefixes`;
6. first Go byte-literal OpenAPI parity test;
7. second Go byte-literal OpenAPI parity test;
8. obsolete single-level Prettier exclusion;
9. `api-sdk` shared seam;
10. `sdk-runtime` knowledge route;
11. ratchet/baseline entries tied to retiring SDK tests;
12. `GOV_FRONTEND_FETCH` path-name exemption and its synthetic governance fixtures.

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
20. Product Problem schemas use exact URI constants only after OA-C9 proof;
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
43. manual SDK types are not imported as target wire authority;
44. all twelve retirement consumers are removed or deliberately rehomed atomically;
45. ngrok preview hosting cannot leak into canonical Product Problem constants.

A green command that did not traverse the real source closure, all operations and generated outputs is vacuous.

---

## 4. GPT adjudication of Fable findings

| Finding | Adjudication |
|---|---|
| F-1 standard mux cannot serve `:verb` grammar | ACCEPT. Preserve W1; record D7 router constraint; custom-mux compatibility fixture. |
| F-2 generator-specific type overrides create two contracts | ACCEPT. Closed extension allowlist; type/shape overrides forbidden. |
| F-3 Prettier reaches nested OAD | ACCEPT. Exclude `contracts/api/product/**` from Prettier; Redocly owns form. |
| F-4 proposed GitHub Pages origin not established | ACCEPT root cause. Reject ngrok as temporary canonical identity; retain stable-origin proof gate. |
| F-5 legacy retirement under-scoped | ACCEPT corrected twelve-consumer set. |
| F-6 four names not filed in W4 | ACCEPT. Name-only W4 filing after final ratification. |
| F-7 Availability GET lost `effective` | ACCEPT. Use `GetEffectiveAvailabilityAllocationScopePolicy`. |
| F-8 version pin/no-empty-servers gaps | ACCEPT. Exact custom pin; explicitly disable contradictory rule. |
| F-9 Redocly config can hijack TypeScript input | ACCEPT. Rules-only config; explicit CLI entry/output. |
| F-10 Go requires bundle | ACCEPT. Bundle mandatory in Go pipeline, never authority. |
| F-11 component collisions silently rename | ACCEPT. Pre-bundle uniqueness plus bundle collision assertion. |
| F-12 closure is declaration-only | ACCEPT. D5 declaration proof; D7 runtime validation obligation. |
| F-13 Go enum identifiers can churn | ACCEPT evidence; do not admit `x-enum-varnames` now. |
| F-14 multipart "usable" needed precision | ACCEPT. Expressible/reachable, not automatic binding/validation. |
| Additional GPT finding: exact Go toolchain unproved | ACCEPT. Add exact Go 1.25.1 proof alongside current-toolchain proof. |

Round 2 is not required. Every correction is D5-B2-local or a name-only W4 filing; no semantic parent contradiction survives.

---

## 5. Final disposition pending ratification

```text
W1 / W2 / W3 / Technical Ingress               CONFIRMED
W4 semantics                                    CONFIRMED
W4 operation spelling                           NAME-ONLY FILING REQUIRED
Product operations                              95 unchanged
ordinary Permissions                            29 unchanged
source model                                    one multi-document design-first OAD
OAS feature level                               3.1.2
lint/bundle                                     Redocly CLI 2.45.0
TypeScript contract generation                  openapi-typescript 7.13.0
Go compatibility generation                     oapi-codegen v2.8.0 + runtime v1.7.0
runtime TypeScript client                       D6
Go router/server/runtime validator               D7
legacy OpenAPI/manual SDK                       REPLACE / RETIRE AS ONE MEASURED SEAM
ngrok endpoint                                  PREVIEW/PROOF HOST ONLY
canonical Product Problem origin                 BLOCKED BY STABLE-ORIGIN PROOF
independent Fable review                         COMPLETE
GPT adjudication                                 CONVERGED
Round 2                                          NOT REQUIRED
operator final package ratification              NEXT
canonical filing                                 BLOCKED BY RATIFICATION
Product OAD authoring                            BLOCKED BY RATIFICATION + OA-C9 PROOF
D6–D9 / implementation                          BLOCKED
```

## 6. Final ratification effect

After explicit operator ratification of this converged package:

1. file the four operation names into W4 §8 without changing class, Permission or Principal admission;
2. make this decision canonical in the appropriate D5-B2 home;
3. remove this candidate from the active tree;
4. reset `AI-DIALOG.md` to protocol-only;
5. update the router and cockpit projection;
6. open the bounded stable-origin proof if still unresolved;
7. author the Product OAD only after the OA-C9 origin gate is discharged;
8. preserve implementation blocked until D9.

## 7. Reopen triggers

Reopen only the smallest implicated decision when evidence shows:

- OAS 3.1.2 cannot express an accepted W1–W4 property without loss;
- all credible pinned generators fail a required idiom;
- a real Product consumer requires a runtime SDK property types-only generation cannot support safely;
- a stable controlled Problem URI origin cannot be established before custom type publication;
- a real offline/external consumer requires a committed bundle;
- a real public compatibility consumer requires URI versioning;
- a Product operation cannot be expressed without new meaning;
- a technical surface is mistakenly required by a Product client.

Do not reopen for generator naming preference, file-count aesthetics, current handler layout, legacy route convenience, desire for `/v1`, or one tool's inability when another credible tool preserves the source semantics.
