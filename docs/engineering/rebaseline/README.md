# Marketplace Central — Architecture Rebaseline

> **Role:** sole current-program status/router after `AGENTS.md`  
> **Current phase:** **D5 — API — OPEN / ACTIVE; D5-B1 ACCEPTED / CANONICAL; D5-B2 OPEN / ACTIVE**  
> **D5-B2 current state:** **W1 + W2 + W3 + W4 + Technical Ingress + Final Problem/media + OpenAPI wire authority/tooling ACCEPTED / CANONICAL; stable Product Problem origin FILED; canonical Product OAD authoring/proof = NEXT**  
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
17. `docs/engineering/rebaseline/D5-B2-W2-SCHEMA-GRAMMAR.md` — canonical W2
18. `docs/engineering/rebaseline/D5-B2-W3-COLLECTION-GRAMMAR.md` — canonical W3
19. `docs/engineering/rebaseline/D5-B2-W4-PERMISSION-CLIENT-CLASS-ENFORCEMENT.md` — canonical W4 and exact 95-operation naming/mapping
20. `docs/engineering/rebaseline/D5-B2-TECHNICAL-INGRESS.md` — canonical technical non-Product ingress
21. `docs/engineering/rebaseline/D5-B2-OPENAPI-WIRE-AUTHORITY-TOOLING.md` — canonical machine-readable wire-authority/tooling decision and Product Problem URI origin
22. `docs/engineering/rebaseline/EVIDENCE-REGISTER.md`
23. code/OpenAPI/schemas/tests/runtime only as evidence when needed

This file alone owns **program status, allowed/blocked work and exact next action**. Detailed semantics remain in accepted artifacts.

`AI-DIALOG.md` is protocol-only with no active cycle. Completed review cycles remain in Git history.

`docs/engineering/standards/evidence-grounded-production-engineering-for-llm-agents.md` is a derived portable production-research guide, not architecture authority.

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
    W1 Resource / Path / HTTP Grammar                    ACCEPTED / CANONICAL
    W2 Request / Response Schema Grammar                 ACCEPTED / CANONICAL
    W3 Collections / Query / Cursor Grammar              ACCEPTED / CANONICAL
    W4 Permission / Client-Class Enforcement             ACCEPTED / CANONICAL
    Technical non-Product ingress                        ACCEPTED / CANONICAL
    Final Problem/media consistency                      ACCEPTED / CANONICAL
    Single OpenAPI wire authority/tooling decision       ACCEPTED / CANONICAL
      Fable compatibility review                         COMPLETE
      GPT adjudication                                   COMPLETE / CONVERGED
      operator ratification                              COMPLETE / FILED
    Stable Product Problem origin proof                  COMPLETE / FILED
      stable origin                                      https://conexus.fun
      ngrok                                              PREVIEW ONLY
    Canonical Product OAD authoring/proof                OPEN / NEXT
D6 — Frontend                                             BLOCKED BY D5
D7 — Runtime / Jobs / Transactions                        BLOCKED
D8 — Golden Flows                                         BLOCKED
D9 — Adversarial Architecture Review                      BLOCKED
Implementation                                            BLOCKED UNTIL D9
```

---

## 3. Filed D5-B2 invariants

Detailed meaning lives in the canonical homes above. This section is routing only.

### Product surface

```text
Product operations                              95
ordinary Permissions                            29
Principal kinds                                 H / A / S only
standalone Product media GET/Update/Delete      not admitted
```

Final W4 operation names include:

```text
GetEffectiveAvailabilityAllocationScopePolicy
UpdateAvailabilityAllocationScopePolicy
ListFulfillmentExecutions
GetFulfillmentExecution
```

### Technical ingress and Problem/media

- provider/OAuth acquisition remains outside Product OpenAPI/SDK;
- provider signal is neither D3 event nor current provider truth;
- authored-media byte delivery is a separate technical presentation surface, not Product operation;
- Product and technical Problem namespaces remain disjoint;
- retention/erasure remains bounded to its accepted D2/D7 reopen conditions.

### Machine-readable Product wire authority

```text
entry document                       contracts/api/product/openapi.yaml
OpenAPI                              3.1.2
source OAD                           sole wire authority
resolved bundle                      derived / temporary
TypeScript projection                generated from source
generated Go                         generated from resolved bundle
```

Pinned baseline:

```text
@redocly/cli                    2.45.0
openapi-typescript              7.13.0
oapi-codegen                    v2.8.0
github.com/oapi-codegen/runtime v1.7.0
```

Binding controls include local refs, globally unique component keys, exact OAS pin, rules-only Redocly config, extension allowlist, exact W4 mapping, exact Go 1.25.1 + current-toolchain proof, and atomic retirement of the measured legacy OpenAPI/manual-SDK seam.

The standard Go `net/http.ServeMux` cannot directly realize canonical `{id}:verb`; D7 must later choose/provide a compatible mux without changing W1.

### Stable Product Problem origin

The stable DevelopmentConexus domain is:

```text
https://conexus.fun
```

Marketplace Central Product Problem identifiers use:

```text
https://conexus.fun/marketplace-central/problems/product/{slug}
```

Future technical protocol/presentation contracts use, when needed:

```text
https://conexus.fun/marketplace-central/problems/technical/{surface}/{slug}
```

`conexus.fun` is intentionally shared across DevelopmentConexus projects; `/marketplace-central/...` owns only this product's namespace. Hosting may change without changing the identifier.

The previous ngrok hostname remains temporary preview/tunnel infrastructure only and is forbidden from Product `Problem.type` constants.

---

## 4. Delivery staging rule

`AGENTS.md` now explicitly distinguishes **target architecture** from **first-release implementation scope**.

When implementation opens after D9, work is classified proportionately as:

```text
BUILD NOW   current consumer/golden-flow or correctness requirement
SEAM NOW    preserve the boundary now, future capability later
PROVE FIRST smallest bounded spike before committing mechanism
DEFER       no current consumer/failure class; retain reopen trigger
```

This does not alter D0→D9 status or allow implementation early. It prevents the accepted target architecture from being misread as an instruction to build every future-scale capability before the first useful internal vertical slice.

---

## 5. Allowed now

Only one coherent architectural sub-batch is open:

> **Author and prove the canonical Product OpenAPI Description from accepted W1–W4 and the canonical OpenAPI authority/tooling profile.**

The sub-batch may proportionately:

1. create `contracts/api/product/openapi.yaml` and its local referenced `paths/` / `components/` source tree;
2. encode exactly the 95 W4 Product operations and 29-Permission projection without adding business meaning;
3. use exact Product Problem constants under `https://conexus.fun/marketplace-central/problems/product/{slug}`;
4. add the bounded rules-only Redocly configuration and deterministic lint/bundle/generation proof required by the canonical tooling authority;
5. generate/prove the derived TypeScript and Go contract projections;
6. apply only the measured mechanical retirement/rehoming work necessary to prevent the legacy OpenAPI/manual SDK from remaining a parallel target authority;
7. run the canonical authoring negative controls and repository gates;
8. stop and reopen only the smallest authority if executable proof discovers a real semantic contradiction.

This is **contract authoring/proof**, not product implementation.

---

## 6. Prohibited now

- do not begin D6–D9 or product implementation;
- do not build all 95 runtime handlers merely because the OAD contains 95 operations;
- do not choose Go router/server framework, runtime validator, deployment topology, queue, storage, CDN or cache in D5;
- do not create self-service tenancy, generic connector/plugin frameworks or scale infrastructure without a current consumer/failure class;
- do not keep the legacy OpenAPI or manual Product SDK as a parallel target authority;
- do not add a 96th Product operation, 30th Permission, new Principal kind, media CRUD or `/v1`;
- do not convert MPC Permissions to OAuth scopes;
- do not place Technical Ingress A/B or authored-media byte delivery in the Product OAD/SDK;
- do not use ngrok or another temporary host in canonical Product Problem identifiers;
- do not weaken accepted security/data/tenant/recovery invariants in the name of an internal first release.

Implementation remains blocked until D9.

---

## 7. Exact next action

**Canonical Product OAD authoring/proof.**

Expected sequence:

```text
source OAD authoring
→ Redocly lint + deterministic bundle
→ TypeScript contract generation/compile
→ Go contract generation/compile on exact minimum + current toolchain
→ 95-operation / W4 / Problem / exclusion negative controls
→ measured legacy authority retirement/rehoming
→ repository gates
→ D5-B2 closure review
```

Do not jump from OpenAPI authoring to runtime implementation.

---

## 8. Fresh-session success test

A fresh session must conclude:

- D0→D4 and D5-B1 are accepted/canonical;
- W1/W2/W3/W4, Technical Ingress, Final Problem/media and OpenAPI authority/tooling are accepted/canonical;
- Product surface remains 95 operations / 29 Permissions;
- `https://conexus.fun` is the filed stable origin and Marketplace Central uses its `/marketplace-central/problems/...` namespace;
- target architecture does not imply first-release scope; implementation staging uses BUILD NOW / SEAM NOW / PROVE FIRST / DEFER after D9;
- canonical Product OAD authoring/proof is the exact next action;
- D6–D9 and implementation remain blocked.

If not, the active authority tree is inconsistent.