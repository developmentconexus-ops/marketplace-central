# Marketplace Central — Architecture Rebaseline

> **Role:** sole current-program status/router after `AGENTS.md`  
> **Current phase:** **D5 — API — OPEN / ACTIVE; D5-B1 ACCEPTED / CANONICAL; D5-B2 OPEN / ACTIVE**  
> **D5-B2 current state:** **W1 + W2 + W3 + W4 + Technical Ingress + Final Problem/media consistency ACCEPTED / CANONICAL; single OpenAPI wire authority / tooling decision OPEN / ACTIVE; GPT adjudication CONVERGED; operator final ratification = NEXT**  
> **Implementation:** **BLOCKED until D9 is accepted**  
> **Last updated:** 2026-08-20

## 1. Authority path

A fresh session reads, in order:

1. `AGENTS.md`
2. this router
3. `docs/engineering/standards/root-cause-global-maximum-method.md`
4. `ARCHITECTURE.md`
5. `docs/engineering/rebaseline/DECISION-RECONCILIATION-BASELINE.md`
6. `docs/architecture/decisions/README.md`
7. `docs/engineering/rebaseline/D0-PRODUCT-SYSTEM-DEFINITION.md`
8. `docs/engineering/rebaseline/D1-DOMAINS-BOUNDARIES.md`
9. `docs/engineering/rebaseline/D2-IDENTITY-TENANT-DATA-OWNERSHIP.md`
10. `docs/engineering/rebaseline/D3-COMMUNICATION-EVENTS.md`
11. `docs/engineering/rebaseline/D4-EXTERNAL-INTEGRATIONS.md`
12. `docs/engineering/rebaseline/D4-R1-PUBLICATION-INPUT.md`
13. `docs/engineering/rebaseline/D5-API.md`
14. `docs/engineering/rebaseline/D5-B2-PRODUCT-OPERATION-SURFACE.md`
15. `docs/engineering/rebaseline/D5-B2-OPERATION-ADMISSION-MATRIX.md`
16. `docs/engineering/rebaseline/D5-B2-WIRE-CONTRACT.md` — canonical W1
17. `docs/engineering/rebaseline/D5-B2-W2-SCHEMA-GRAMMAR.md` — canonical W2, including canonical Problem Details and authored-media grammar
18. `docs/engineering/rebaseline/D5-B2-W3-COLLECTION-GRAMMAR.md` — canonical W3
19. `docs/engineering/rebaseline/D5-B2-W4-PERMISSION-CLIENT-CLASS-ENFORCEMENT.md` — canonical W4
20. `docs/engineering/rebaseline/D5-B2-TECHNICAL-INGRESS.md` — canonical technical non-Product ingress
21. `docs/engineering/rebaseline/EVIDENCE-REGISTER.md`
22. code/OpenAPI/schemas/tests/runtime only as evidence when needed

This file alone owns **program status, allowed/blocked work and exact next action**. Detailed semantics remain in accepted artifacts.

`docs/engineering/rebaseline/D5-B2-OPENAPI-WIRE-AUTHORITY-TOOLING-CANDIDATE.md` is **NON-AUTHORITATIVE converged adjudication evidence**. It records the lead direction, Fable compatibility findings, GPT adjudication and the remaining stable-origin proof gate. It does not make the tooling decision canonical and does not authorize Product OpenAPI authoring.

`AI-DIALOG.md` preserves the completed Fable review cycle as working evidence until final ratification/filing. It is never authority.

`docs/engineering/standards/evidence-grounded-production-engineering-for-llm-agents.md` is a derived portable production-research guide. It is not part of this architecture authority path and does not redefine the organizational Method.

`docs/engineering/rebaseline/cockpit.html` is a non-authoritative visual projection. No canonical status changed during review/adjudication, so the cockpit is not updated by this sub-batch.

---

## 2. Program state

```text
D0 — Product / System Definition                         CLOSED / ACCEPTED
D1 — Domains / Boundaries                                CLOSED / ACCEPTED
D2 — Identity / Tenant / Data Ownership                  CLOSED / ACCEPTED
D3 — Communication / Events                              CLOSED / ACCEPTED
D4 — External Integrations                               CLOSED / ACCEPTED AS A WHOLE
D4-R1 — Publication Input & Listing Authoring            ACCEPTED / CANONICAL
Decision Reconciliation                                  ACCEPTED / CANONICAL
D5 — API                                                  OPEN / ACTIVE
  B1 Semantic API Model                                  ACCEPTED / CANONICAL
  B2 Product Operation / Resource Surface                 OPEN / ACTIVE
    Operation Admission Matrix                           ACCEPTED / RATIFIED
    Whole-Matrix Global Coherence                        ACCEPTED / RATIFIED
    Wire Contract
      W1 Resource / Path / HTTP Grammar                  ACCEPTED / CANONICAL
      W2 Request / Response Schema Grammar               ACCEPTED / CANONICAL
      W3 Collections / Query / Cursor Grammar            ACCEPTED / CANONICAL
      W4 Permission / Client-Class Enforcement           ACCEPTED / CANONICAL
      Technical non-Product ingress                      ACCEPTED / CANONICAL
        A External Acquisition Ingress                   ACCEPTED / CANONICAL
        B OAuth / Authorization Ceremony                 ACCEPTED / CANONICAL
        Whole-Ingress review/adjudication                 COMPLETE / FILED
      Final Problem/media consistency                    ACCEPTED / CANONICAL
        lead + Fable + GPT review                        COMPLETE / CONVERGED
        operator final ratification                      COMPLETE / FILED
      Single OpenAPI wire authority/tooling decision     OPEN / ACTIVE
        lead authority/tooling research                  COMPLETE
        lead review candidate                            COMPLETE
        independent Fable compatibility/adversarial      COMPLETE
        GPT final adjudication                           COMPLETE / CONVERGED
        Round 2                                           NOT REQUIRED
        operator origin-policy decision                  RECEIVED / BOUNDED
        operator final package ratification              NEXT
        canonical filing                                 BLOCKED BY RATIFICATION
        stable Product Problem origin proof              BLOCKS TYPE CONSTANTS
        canonical Product OpenAPI authoring              BLOCKED BY RATIFICATION + ORIGIN PROOF
D6 — Frontend                                             BLOCKED BY D5
D7 — Runtime / Jobs / Transactions                        BLOCKED
D8 — Golden Flows                                         BLOCKED
D9 — Adversarial Architecture Review                      BLOCKED
Implementation                                            BLOCKED UNTIL D9
```

---

## 3. Final Problem/media consistency — filed dispositions

The converged lead + Fable + GPT package was operator-ratified and filed substitutively into its existing canonical homes. Detailed semantics live only in those artifacts; this list is a routing index, not a second authority.

| Correction | Canonical home |
|---|---|
| RFC 9457 `about:blank` for status-only failures; `405 + Allow`, `415 Accept`, `413 Retry-After` obligations | W2 §§2.15.1–2.15.2 |
| exact `CreateListingIntentMedia` failure map (400/413/415/422/409/500) | W2 §3.9.2 |
| `400 idempotency-key-required`; `422 idempotency-key-reused` and `409 idempotency-request-in-progress` preserved unchanged | W2 §§3.9.2, 18, 19 |
| bounded `409 resource-state-conflict` family | W2 §§3.9.2, 19 |
| transport size/media-type guards before full-body buffering, without tenant/resource disclosure | W2 §§3.9.3, 18 |
| immutable ListingIntent-scoped authored-media identity | W2 §3.9.5 |
| successful media creation advances and returns the parent ListingIntent validator as typed result data | W2 §3.9.1 |
| exact idempotent retry resolves the prior intake before re-evaluating the now-stale validator | W2 §§3.9.1, 18 |
| stable `ListingIntentMediaDescriptor` vs volatile `ListingIntentMediaPresentationDescriptor` | W2 §3.9.7 |
| access reference never enters history, fingerprint, logs or Problem Details | W2 §3.9.7 |
| authored-media byte delivery = separately justified technical presentation surface reusing current Organization + `offering.read`; no durable anonymous locator; failures do not expand the Product Problem catalog | W2 §3.9.8, W1 §20, W4 §15.2, Technical Ingress §12.1 |
| retention/erasure remains Unknown / deferred to D2 + D7; no `DeleteMedia`, no invented lifetime | W2 §3.9.6 |
| Product and technical-protocol Problem namespaces remain disjoint | W2 §19 |
| source vs authored media stay distinct descriptor families | W2 §3.9.9 |
| converged negative controls / proof obligations | W2 §23 items 43–67 |

Preserved:

```text
Product operations                              95
ordinary Permissions                            29
standalone Product media GET/Update/Delete       not admitted
new Principal kind                               none
Technical Ingress A/B                            canonical
W3 collection/cursor semantics                   canonical
implementation                                   blocked until D9
```

---

## 4. OpenAPI authority/tooling converged package — non-authoritative pending ratification

The lead direction survived independent executable attack with bounded corrections.

```text
one multi-document design-first Product OAD
root: contracts/api/product/openapi.yaml
OAS: 3.1.2

source OAD                            sole wire authority
bundle                               temporary derived input for Go
TypeScript projection                generated directly from source
Go projection                        generated from deterministic bundle
static docs                          derived

lint + bundle                        @redocly/cli 2.45.0
TypeScript contract projection       openapi-typescript 7.13.0
Go compatibility projection          oapi-codegen v2.8.0
Go generated runtime dependency       oapi-codegen/runtime v1.7.0

runtime TypeScript transport          D6
Go router/server/runtime validation   D7
```

### 4.1 Confirmed direction

- one entry document plus local relative refs remains one authority;
- no `/v1`, environment `servers` or technical routes;
- bearer authentication without converting the 29 MPC Permissions into OAuth scopes;
- bounded `x-mpc-*` W4 projections;
- exact OAS 3.1 multipart raw file + `text/plain` ETag part;
- mandatory generated TypeScript contract types;
- runtime TypeScript transport deferred to D6;
- Go server/router/validator realization deferred to D7;
- legacy OpenAPI/manual SDK replaced as one measured seam;
- no 96th Product operation, 30th Permission or new Principal kind.

### 4.2 Name-only operation crystallizations

Pending final filing into W4 §8:

```text
GetEffectiveAvailabilityAllocationScopePolicy
UpdateAvailabilityAllocationScopePolicy
ListFulfillmentExecutions
GetFulfillmentExecution
```

The Availability GET preserves configured value/mode, effective value and effective-source provenance; UPDATE mutates only the configured owner value/mode. They do not share one round-trip schema.

### 4.3 Load-bearing compatibility corrections

- standard `net/http.ServeMux` is excluded for the canonical `{id}:verb` grammar; D7 selects or implements a compatible mux, while D5 proves dispatch through a bounded custom mux;
- source extensions are allowlisted; type/shape-changing `x-go-*` extensions are forbidden;
- Product OAD source files are outside Prettier;
- `redocly.yaml` is rules-only with no `apis:` block;
- exact OAS 3.1.2 is asserted mechanically;
- `no-empty-servers` is disabled deliberately;
- component names are globally unique before bundle and collision suffixes are forbidden after bundle;
- Go generation consumes the temporary bundle; TypeScript generation consumes source directly;
- D5 proves schema declarations; D7 owns runtime unknown-field/pattern/multipart-part validation;
- Go enum naming churn remains derived/compiler-visible; no `x-enum-varnames` enters source now;
- exact Go 1.25.1 proof remains required in addition to the already-run newer-toolchain proof.

### 4.4 Legacy retirement set

The replacement sub-batch must retire or rehome together:

1. `GOV_API_SDK_SPLIT`;
2. `api-sdk-atomicity`;
3. `openapi-without-sdk` eval;
4. legacy OpenAPI knowledge route;
5. legacy `modules.json` path prefixes;
6. two Go literal OpenAPI parity tests;
7. obsolete single-level Prettier exclusion;
8. `api-sdk` shared seam;
9. `sdk-runtime` knowledge route;
10. baselines tied to retiring SDK tests;
11. governance fixtures encoding the old seam;
12. the `GOV_FRONTEND_FETCH` package-name exemption.

Useful atomicity/drift controls are replaced by source-OAD → deterministic-generation → clean-regeneration controls, not simply deleted.

### 4.5 Product Problem origin disposition

The operator nominated:

```text
https://multiradial-unironically-nieves.ngrok-free.dev/
```

Official ngrok documentation classifies the free hostname as one automatically assigned **development domain** specific to the account. The free plan cannot customize that hostname or bring a user-owned domain.

Therefore the endpoint is accepted only for:

```text
temporary preview documentation
resolution/hosting compatibility proof
```

It is not silently reclassified as a DevelopmentConexus-owned verified domain and is not embedded into permanent Product `Problem.type` constants "for now".

Before custom Product Problem types enter the Product OAD, prove either:

1. a DevelopmentConexus-owned verified HTTPS origin; or
2. explicit operator acceptance of a permanent third-party-domain residual, including an obligation never to relinquish/rename the controlling account/hostname.

Until then:

```text
ngrok preview hosting                         ALLOWED
custom Product Problem URI constants          BLOCKED
about:blank where already canonical           UNAFFECTED
Product OAD authoring                         BLOCKED AT ORIGIN GATE
```

This is a bounded proof gate, not a Product semantic reopen.

---

## 5. Allowed now

Only the following coherent action is open:

> **Operator final ratification of the converged OpenAPI wire authority/tooling package recorded in `D5-B2-OPENAPI-WIRE-AUTHORITY-TOOLING-CANDIDATE.md`.**

Ratification must cover the package as adjudicated, including:

- the four final operation names;
- exact pins and source/bundle generation pipeline;
- extension/formatter/config/component controls;
- standard-mux exclusion as a D7 constraint, not a W1 change;
- the twelve-item legacy retirement seam;
- exact Go 1.25.1 proof obligation;
- ngrok as temporary preview/proof hosting only;
- stable Product Problem origin proof before custom `type` constants or Product OAD authoring.

No further Fable round is required unless new material evidence creates a real contradiction.

---

## 6. Prohibited now

- do not treat the candidate, Fable findings or this router summary as canonical detailed semantics;
- do not author the Product OAD before final ratification and the stable-origin proof;
- do not embed the ngrok hostname into permanent Product Problem identifiers merely because it is free today;
- do not modify W1/W2/W3/W4/Technical Ingress as a tooling convenience;
- do not keep legacy OpenAPI or manual SDK as parallel target authority;
- do not select OAS 3.2 by recency or downgrade to 3.0 for one generator;
- do not create independent OADs per D1 owner;
- do not use code-first Go as target authority;
- do not use remote refs, templates, macros or experimental joins as authority composition;
- do not put Technical Ingress A/B or authored-media delivery in Product OpenAPI/SDK;
- do not add a 96th Product operation, 30th Permission, Principal kind, media CRUD or `/v1`;
- do not convert MPC Permissions into OAuth scopes;
- do not select the D7 router, runtime validator, storage, CDN or deployment now;
- do not update the cockpit from an unratified package;
- do not begin D6–D9 or implementation.

Implementation remains blocked until D9.

---

## 7. Exact next action

**Obtain operator final ratification of the converged OpenAPI wire authority/tooling package.**

After ratification:

```text
canonical filing into D5-B2 homes
→ name-only W4 §8 correction
→ candidate removal
→ AI-DIALOG reset to protocol-only
→ router/cockpit synchronization
→ stable Product Problem origin proof if still unresolved
→ Product OAD authoring/proof only after the origin gate
```

Implementation remains blocked until D9.

---

## 8. Fresh-session success test

A fresh session must conclude:

- W1/W2/W3/W4, Technical Ingress and Final Problem/media remain accepted/canonical;
- Product surface remains 95 operations / 29 Permissions;
- the OpenAPI/tooling candidate is non-authoritative converged evidence;
- Fable review and GPT adjudication are complete;
- Round 2 is not required;
- final operator package ratification is the exact next action;
- ngrok is temporary preview/proof hosting, not yet a canonical Product Problem identity origin;
- custom Product Problem URI constants and Product OAD authoring remain blocked by the stable-origin proof;
- D6–D9 and implementation remain blocked.

If not, the active authority tree is inconsistent.
