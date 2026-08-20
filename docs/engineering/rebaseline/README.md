# Marketplace Central — Architecture Rebaseline

> **Role:** sole current-program status/router after `AGENTS.md`  
> **Current phase:** **D5 — API — OPEN / ACTIVE; D5-B1 ACCEPTED / CANONICAL; D5-B2 OPEN / ACTIVE**  
> **D5-B2 current state:** **W1 + W2 + W3 + W4 + Technical Ingress + Final Problem/media consistency ACCEPTED / CANONICAL; single OpenAPI wire authority / tooling decision OPEN / ACTIVE; independent Fable review = NEXT**  
> **Implementation:** **BLOCKED until D9 is accepted**  
> **Last updated:** 2026-08-19

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

`docs/engineering/rebaseline/D5-B2-OPENAPI-WIRE-AUTHORITY-TOOLING-CANDIDATE.md` is **NON-AUTHORITATIVE lead decision evidence**. It does not amend W1/W2/W3/W4, does not make a tooling selection canonical and does not authorize Product OpenAPI authoring.

`AI-DIALOG.md` is the working review channel for the active independent OpenAPI/tooling review. It is never authority.

`docs/engineering/standards/evidence-grounded-production-engineering-for-llm-agents.md` is a derived portable production-research guide. It is not part of this architecture authority path and does not redefine the organizational Method.

`docs/engineering/rebaseline/cockpit.html` is a non-authoritative visual projection. No canonical status changed when this review candidate opened, so the cockpit is not updated by this sub-batch.

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
        independent Fable compatibility/adversarial      NEXT
        GPT final adjudication                           BLOCKED BY REVIEW
        Round 2                                           ONLY IF MATERIAL CONFLICT
        operator final ratification                      BLOCKED BY ADJUDICATION
        canonical Product OpenAPI authoring              BLOCKED BY RATIFICATION
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

## 4. OpenAPI authority/tooling lead package — non-authoritative

The lead candidate derives from canonical W1–W4 and current primary tooling evidence. Its proposed direction is:

```text
one multi-document design-first Product OAD
root: contracts/api/product/openapi.yaml
OAS: 3.1.2

source OAD                            authority
bundle / generated Go / generated TS derived projections

lint + bundle                        @redocly/cli 2.45.0
TypeScript contract projection       openapi-typescript 7.13.0
Go compatibility projection          oapi-codegen v2.8.0

runtime TypeScript transport          D6
Go server framework                   D7
```

It also proposes:

- exact 95 unique PascalCase `operationId` values;
- final D5-local names `Get/UpdateAvailabilityAllocationScopePolicy` and `List/GetFulfillmentExecution(s)`;
- no `/v1` path;
- no environment-specific server URL in the source OAD;
- bearer authentication without turning MPC Permissions into OAuth scopes;
- bounded `x-mpc-*` projections for W4 reviewability;
- OAS 3.1-native multipart file + typed ETag representation;
- Product Problem type namespace under the repository-controlled GitHub Pages origin;
- mandatory generated TypeScript contract types;
- no hand-authored target Product SDK DTO/operation catalog;
- runtime client selection deferred to D6 because currently evaluated candidates are experimental, pre-1.0 or entering maintenance mode;
- current 247 KB legacy OpenAPI and manually authored `sdk-runtime` treated as replacement targets, not compatibility authorities;
- 53 named negative controls/proof obligations before canonical authoring can close.

This direction is not accepted until independent compatibility/adversarial review, GPT adjudication and operator ratification.

---

## 5. Allowed now

Only the following coherent work is open:

> **Independent Fable compatibility and adversarial review of the single OpenAPI wire authority/tooling candidate.**

Fable must:

1. reconstruct authority independently;
2. review the candidate as one package;
3. challenge the Global Maximum and all hidden second-authority risks;
4. run bounded executable compatibility fixtures for the pinned Redocly, openapi-typescript and oapi-codegen versions;
5. test OAS 3.1.2 multi-document refs, `oneOf+const`, closed objects, multipart file+ETag, Problem type constants, headers, custom methods and generated compilation;
6. challenge the GitHub Pages Problem URI lifecycle/control;
7. challenge D5 types-only SDK versus D6 runtime-client defer;
8. write only in `AI-DIALOG.md`;
9. end `HANDOFF → GPT`.

No canonical artifact or Product OpenAPI file may be modified during this review.

---

## 6. Prohibited now

- do not treat the OpenAPI/tooling candidate as authority;
- do not author the new Product OpenAPI before ratification;
- do not modify W1/W2/W3/W4/Technical Ingress semantics as a tooling convenience;
- do not keep the legacy OpenAPI or manual SDK as a parallel target authority;
- do not select OAS 3.2 by recency alone;
- do not downgrade to OAS 3.0 merely because one generator is behind;
- do not create one independently complete OAD per D1 owner;
- do not use code-first Go as target wire authority;
- do not use remote `$ref`, templates, macros, standing overlays or experimental `join` as authority composition;
- do not put Technical Ingress A/B or authored-media byte delivery in Product OpenAPI/SDK;
- do not add a 96th Product operation, 30th Permission, Principal kind, media CRUD or `/v1`;
- do not turn MPC Permissions into OAuth scopes;
- do not select Go router/server framework, runtime request validator, frontend cache/query integration, deployment host or D7 storage;
- do not update the cockpit from an unratified candidate;
- do not begin D6–D9 or implementation.

Implementation remains blocked until D9.

---

## 7. Exact next action

**Fable independent compatibility/adversarial review of `D5-B2-OPENAPI-WIRE-AUTHORITY-TOOLING-CANDIDATE.md` through the active `AI-DIALOG.md` cycle.**

After Fable returns:

```text
GPT revalidates remote HEAD
→ adjudicates every material finding
→ requests Round 2 only if a real contradiction survives
→ asks operator to ratify only a genuinely converged package
→ canonical filing
→ Product OpenAPI authoring/proof
```

Implementation remains blocked until D9.

---

## 8. Fresh-session success test

A fresh session must conclude:

- W1/W2/W3/W4, Technical Ingress and Final Problem/media consistency remain accepted/canonical;
- Product surface remains 95 operations / 29 Permissions;
- the OpenAPI/tooling candidate is non-authoritative;
- the lead proposes one OAS 3.1.2 multi-document Product OAD and generated projections;
- current OpenAPI and manual SDK remain evidence/replacement targets only;
- independent Fable compatibility/adversarial review is the exact next action;
- Product OpenAPI authoring, D6–D9 and implementation remain blocked.

If not, the active authority tree is inconsistent.
