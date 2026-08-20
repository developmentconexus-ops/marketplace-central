# Marketplace Central — Architecture Rebaseline

> **Role:** sole current-program status/router after `AGENTS.md`  
> **Current phase:** **D5 — API — OPEN / ACTIVE; D5-B1 ACCEPTED / CANONICAL; D5-B2 OPEN / ACTIVE**  
> **D5-B2 current state:** **W1 + W2 + W3 + W4 + Technical Ingress + Final Problem/media consistency + Single OpenAPI wire authority/tooling ACCEPTED / CANONICAL; stable Product Problem origin proof = NEXT**  
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
19. `docs/engineering/rebaseline/D5-B2-W4-PERMISSION-CLIENT-CLASS-ENFORCEMENT.md` — canonical W4 and exact 95-operation names/mapping
20. `docs/engineering/rebaseline/D5-B2-TECHNICAL-INGRESS.md` — canonical technical non-Product ingress
21. `docs/engineering/rebaseline/D5-B2-OPENAPI-WIRE-AUTHORITY-TOOLING.md` — canonical machine-readable wire-authority/tooling decision
22. `docs/engineering/rebaseline/EVIDENCE-REGISTER.md`
23. code/OpenAPI/schemas/tests/runtime only as evidence when needed

This file alone owns **program status, allowed/blocked work and exact next action**. Detailed semantics remain in accepted artifacts.

`AI-DIALOG.md` is protocol-only with no active cycle. It is never authority. Completed review cycles remain in Git history.

`docs/engineering/standards/evidence-grounded-production-engineering-for-llm-agents.md` is a derived portable production-research guide. It is not part of this architecture authority path and does not redefine the organizational Method.

`docs/engineering/rebaseline/cockpit.html` is a non-authoritative visual projection synchronized after canonical status changes.

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
      Single OpenAPI wire authority/tooling decision     ACCEPTED / CANONICAL
        lead authority/tooling research                  COMPLETE
        independent Fable compatibility/adversarial      COMPLETE
        GPT final adjudication                           COMPLETE / CONVERGED
        Round 2                                           NOT REQUIRED
        operator final package ratification              COMPLETE / FILED
        canonical filing                                 COMPLETE
      Stable Product Problem origin proof                OPEN / NEXT
      Canonical Product OpenAPI authoring                 BLOCKED BY ORIGIN PROOF
D6 — Frontend                                             BLOCKED BY D5
D7 — Runtime / Jobs / Transactions                        BLOCKED
D8 — Golden Flows                                         BLOCKED
D9 — Adversarial Architecture Review                      BLOCKED
Implementation                                            BLOCKED UNTIL D9
```

---

## 3. Final Problem/media consistency — filed dispositions

The operator-ratified package is filed substitutively into W1/W2/W4/Technical Ingress. Detailed semantics live only in those homes.

| Correction | Canonical home |
|---|---|
| RFC 9457 `about:blank` for status-only failures; `405 + Allow`, `415 Accept`, `413 Retry-After` obligations | W2 §§2.15.1–2.15.2 |
| exact `CreateListingIntentMedia` failure map | W2 §3.9.2 |
| bounded `409 resource-state-conflict` | W2 §§3.9.2, 19 |
| immutable ListingIntent-scoped authored-media identity | W2 §3.9.5 |
| parent ListingIntent validator returned as typed result | W2 §3.9.1 |
| stable descriptor vs volatile presentation descriptor | W2 §3.9.7 |
| authored-media byte delivery as technical presentation surface | W2 §3.9.8, W1 §20, W4 §15.2, Technical Ingress §12.1 |
| retention/erasure remains Unknown / D2+D7 | W2 §3.9.6 |
| Product and technical Problem namespaces disjoint | W2 §19 |

Preserved:

```text
Product operations                              95
ordinary Permissions                            29
standalone Product media GET/Update/Delete       not admitted
new Principal kind                               none
Technical Ingress A/B                            canonical
W3 collection/cursor semantics                   canonical
```

---

## 4. OpenAPI wire authority/tooling — filed dispositions

The accepted decision lives in `D5-B2-OPENAPI-WIRE-AUTHORITY-TOOLING.md`. This list is routing only, not a second detailed authority.

```text
one logical design-first Product OAD
entry: contracts/api/product/openapi.yaml
OpenAPI: 3.1.2

source OAD                         sole wire authority
resolved bundle                    derived / temporary
generated TypeScript               derived
generated Go                       derived
```

Pinned baseline:

```text
@redocly/cli                    2.45.0
openapi-typescript              7.13.0
oapi-codegen                    v2.8.0
github.com/oapi-codegen/runtime v1.7.0
```

Binding dispositions:

- source refs are local/relative and components are globally unique;
- Redocly config is rules-only with exact OAS pin and `no-empty-servers: off`;
- source OAD extensions are allowlisted; type/shape overrides are forbidden;
- Go generation consumes the temporary bundle; TypeScript generation consumes source;
- the Product OAD source tree is excluded from Prettier;
- standard `net/http.ServeMux` cannot realize canonical `{id}:verb`; D7 must provide compatible partial-segment routing without changing W1;
- D5 proves declarations; D7 owns runtime unknown-field/pattern/multipart validation;
- generated contract types are mandatory; runtime TypeScript transport remains D6;
- the legacy OpenAPI/manual SDK retire as one measured twelve-consumer seam;
- exact Go 1.25.1 and current-toolchain proofs are both required during authoring;
- Technical Ingress A/B and authored-media byte delivery remain outside Product OAD/SDK.

Final W4 operation names:

```text
GetEffectiveAvailabilityAllocationScopePolicy
UpdateAvailabilityAllocationScopePolicy
ListFulfillmentExecutions
GetFulfillmentExecution
```

W4 remains 95 rows / 29 Permissions; the filing is name-only.

---

## 5. Stable Product Problem origin gate

Canonical custom Product Problem identifiers require one immutable HTTPS origin:

```text
<stable-origin>/marketplace-central/problems/product/{slug}
```

Technical protocol/presentation types use a disjoint namespace only in their own future executable contracts.

The operator-nominated endpoint:

```text
https://multiradial-unironically-nieves.ngrok-free.dev/
```

is accepted only for:

```text
temporary preview documentation
resolution/hosting compatibility proof
```

It is not embedded into permanent Product `Problem.type` values merely "for now".

The gate closes only when evidence proves one of:

1. a DevelopmentConexus-owned and verified HTTPS origin; or
2. an explicitly operator-accepted permanent third-party-domain residual, including the obligation never to relinquish/rename the controlling account or hostname.

Until that proof is canonically filed:

```text
ngrok preview hosting                         ALLOWED
about:blank where already canonical           UNAFFECTED
custom Product Problem URI constants          BLOCKED
canonical Product OAD authoring               BLOCKED
```

This is a bounded identity/publication proof gate, not a Product semantic reopen.

---

## 6. Allowed now

Only the following coherent work is open:

> **Prove and file the stable Product Problem HTTPS origin.**

The proof must establish:

- exact immutable origin;
- current control/ownership or explicitly accepted permanent third-party residual;
- HTTPS availability;
- ability to serve or permanently redirect Problem documentation;
- separation of Product and technical namespaces;
- explicit assurance that preview/deployment hosting changes do not change identifier identity.

Do not author the Product OAD while the origin is unresolved, because its custom Product Problem `type` constants would otherwise be temporary identities.

---

## 7. Prohibited now

- do not author the Product OAD before the stable-origin proof is filed;
- do not embed the ngrok hostname into permanent Product Problem identifiers merely because it is currently free/available;
- do not alter W1/W2/W3/W4/Technical Ingress semantics as a tooling convenience;
- do not keep the legacy OpenAPI or manual SDK as parallel target authority;
- do not select OAS 3.2 by recency or downgrade to 3.0 for one generator;
- do not create independent OADs per D1 owner;
- do not use code-first Go as target authority;
- do not use remote refs, templates, macros or experimental joins as authority composition;
- do not put Technical Ingress A/B or authored-media delivery in Product OpenAPI/SDK;
- do not add a 96th Product operation, 30th Permission, Principal kind, media CRUD or `/v1`;
- do not convert MPC Permissions into OAuth scopes;
- do not select the D7 router, runtime validator, storage, CDN or deployment now;
- do not begin D6–D9 or implementation.

Implementation remains blocked until D9.

---

## 8. Exact next action

**Stable Product Problem origin proof.**

After that proof is accepted/filed:

```text
canonical Product OAD authoring + executable proof
→ D5-B2 closure/review as routed
→ D5 close
→ D6
```

Implementation remains blocked until D9.

---

## 9. Fresh-session success test

A fresh session must conclude:

- W1/W2/W3/W4, Technical Ingress, Final Problem/media and OpenAPI wire authority/tooling are accepted/canonical;
- Product surface remains 95 operations / 29 Permissions;
- W4 contains all 95 final operation names with no exception manifest;
- the OpenAPI/tooling review candidate is absent from the active tree;
- `AI-DIALOG.md` is protocol-only with no active cycle;
- ngrok is temporary preview/proof hosting, not a canonical Product Problem identity origin;
- stable Product Problem origin proof is the exact next action;
- Product OAD authoring, D6–D9 and implementation remain blocked.

If not, the active authority tree is inconsistent.
