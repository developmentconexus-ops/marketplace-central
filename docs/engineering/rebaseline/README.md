# Marketplace Central — Architecture Rebaseline

> **Role:** sole current-program status/router after `AGENTS.md`  
> **Current phase:** **D5 — API — OPEN / ACTIVE; D5-B1 ACCEPTED / CANONICAL; D5-B2 OPEN / ACTIVE**  
> **D5-B2 current state:** **B2-A + Operation Matrix + Whole-Matrix ACCEPTED / RATIFIED; Wire W1 + W2-A ACCEPTED IN-STAGE; W2-B ListingIntent + PriceIntent + Availability schemas = NEXT**  
> **Decision Reconciliation:** **ACCEPTED / CANONICAL**  
> **Implementation:** BLOCKED until D9 is accepted  
> **Last updated:** 2026-08-18

## 1. Authority path

A fresh session reads, in order:

1. `AGENTS.md`
2. this file
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
16. `docs/engineering/rebaseline/D5-B2-WIRE-CONTRACT.md`
17. `docs/engineering/rebaseline/D5-B2-W2-SCHEMA-GRAMMAR.md`
18. `docs/engineering/rebaseline/EVIDENCE-REGISTER.md`
19. code, OpenAPI, schemas, tests and runtime only as current-state evidence when needed

This file alone owns **where the program is and what happens next**. Detailed meaning stays in the accepted artifacts above. `AI-DIALOG.md`, Git history, deleted review candidates, legacy routes/current OpenAPI and current code are never status/target authority by inheritance.

`D5-B2-WIRE-CONTRACT.md` owns accepted Wire W1. `D5-B2-W2-SCHEMA-GRAMMAR.md` owns the active W2 schema batch. W2 will receive one coherent independent Fable review only after the W2 sections converge; no micro-review cadence is implied.

## 2. Program state

```text
D0 — Product / System Definition                     CLOSED / ACCEPTED
D1 — Domains / Boundaries                            CLOSED / ACCEPTED
D2 — Identity / Tenant / Data Ownership              CLOSED / ACCEPTED
D3 — Communication / Events                          CLOSED / ACCEPTED
D4 — External Integrations                           CLOSED / ACCEPTED AS A WHOLE
D4-R1 — Publication Input & Listing Authoring        ACCEPTED / CANONICAL
Decision Reconciliation Baseline                     ACCEPTED / CANONICAL
D5 — API                                              OPEN / ACTIVE
  B1 Semantic API Model & Contract Laws              ACCEPTED / CANONICAL
  B2 Product Operation / Resource Surface             OPEN / ACTIVE
    B2-A Client/Auth                                  ACCEPTED IN-STAGE
    Operation Admission Matrix                       ACCEPTED / RATIFIED
    Whole-Matrix Global Coherence                    ACCEPTED / RATIFIED
    Wire Contract
      W1 Resource / Path / HTTP Grammar              ACCEPTED IN-STAGE
      W2 Request/Response Schema Grammar             OPEN / ACTIVE
        W2-A Core Schema Grammar                     ACCEPTED IN-STAGE
        W2-B ListingIntent + PriceIntent + Availability = NEXT
      Fable W2 coherent review                       AFTER W2 CONVERGES
D6 — Frontend                                         BLOCKED BY D5
D7 — Runtime / Jobs / Transactions                    BLOCKED
D8 — Golden Flows                                     BLOCKED
D9 — Adversarial Architecture Review                  BLOCKED
Implementation                                        BLOCKED UNTIL D9
```

## 3. Accepted D5-B2 load-bearing summary

### Client / authority

- humans use OIDC Authorization Code + PKCE; confidential machines use Client Credentials/service-account semantics;
- MPC owns Principal, Organization Membership, AccessRole/Permission/RoleAssignment and all business authority;
- `GET /access-context` is platform-scoped self-only discovery; Organization business operations remain `/organizations/{organization_id}/...`;
- Keycloak remains first implementation/proof candidate; D7 owns concrete deployment/realm/secrets/token realization.

### Operation inventory

- Product API is semantic-owner driven, not CRUD-, screen- or provider-shaped;
- no Product/PIM, generic Integration/Action/Operation/Workflow/Rules/Task/Finance/AI business authority;
- ListingIntent is the create/edit authoring identity; PriceIntent remains separate even for initial publication; Sellable Availability remains Availability-owned;
- BusinessOrderIntent/InvoicingIntent are owner reactions, not direct client commands;
- Fulfillment physical facts cannot be fabricated by generic machine Permission;
- generic Work resolution and cross-owner P remain deferred;
- every admitted C operation has explicit consequence/idempotency/precondition disposition.

### W1 — accepted

- paths are relative to OpenAPI server URL; no `/v1` baseline without a compatibility consumer;
- Organization-owned resources use opaque MPC IDs; external resources retain explicit Marketplace Installation/SourceInstance qualification and no synthetic mirror ID;
- URI nesting means identity/lifecycle containment or source namespace qualification, never workflow order;
- standard HTTP methods are used when honest; owner-specific non-CRUD capabilities use `POST {resource-uri}:verb`;
- writable `status` never substitutes for submit/resolve/physical-evidence semantics;
- current-state protection uses strong opaque MPC `ETag` + `If-Match`; missing required precondition → 428, stale → 412;
- `Idempotency-Key` and concurrency remain separate;
- exact OpenAPI minor version remains a later tooling decision.

### W2-A — accepted

- MPC and native external IDs serialize as opaque strings; semantic IDs remain typed by schema name/context;
- no universal `ExternalRef`/entity graph; external references are typed and source-qualified;
- exact monetary/decision decimals use decimal strings; `Money = exact decimal amount + explicit currency`; no JSON floating-point money, universal minor-unit model or global `round(2)`;
- material times use semantically named unambiguous instants;
- create/update requests and read resources are separate schemas where authority differs;
- Product semantic request objects are closed; property-bag/provider leakage fails contract validation;
- `null` is never a knowledge-state convention;
- knowledge uses the smallest owner/field-specific exclusive union, not a universal `Fact<T>`;
- unions validate through exclusive variants (`oneOf` + fixed discriminant); discriminator metadata is optional tooling aid, not correctness;
- freshness/provenance attaches to the smallest semantic unit for which it is true; no universal Evidence/metadata envelope;
- no generic Result/Operation wrapper; owner resources/results carry owner outcomes;
- business pending/rejected/ambiguous is not automatically an HTTP problem; `202` is not a synonym for asynchronous external convergence;
- PATCH/update bodies are typed owner-specific JSON; omitted means unchanged; `null` is not generic clear; generic JSON Patch/Merge Patch is not baseline;
- RFC 9457 `type` is the primary problem identifier; no duplicate global problem-code taxonomy by default;
- provider-enriched evidence is a closed owner-local variant, never a raw DTO/property bag;
- clients cannot author effective Principal/Organization/approval/convergence/server-history fields.

## 4. Prohibited now

While W2 is open:

- do not begin D6–D9 target design or implementation;
- do not alter accepted D0→D5-B1, B2, W1 or W2-A by schema convenience;
- do not derive schemas from legacy OpenAPI, provider DTOs, database rows or frontend forms;
- do not put price/Availability inside ListingIntent;
- do not create universal Fact/Evidence/Result/Operation/Resource/ExternalRef/property-bag schemas;
- do not collapse unknown/unavailable/partial/not-applicable into null/zero/false/empty;
- do not expose bare native IDs, raw provider payloads/errors or client-authored effective authority fields;
- do not choose D7 blob/upload/server/generator/persistence mechanics during W2;
- do not run Fable yet: W2 review occurs once after the coherent W2 package converges unless a material contradiction forces a focused round.

## 5. Exact next action

**Derive W2-B — ListingIntent + PriceIntent + Availability concrete schema grammar.**

W2-B must decide and adversarially stress-test:

1. `CreateListingIntentRequest`, `UpdateListingIntentRequest`, `ListingIntent`;
2. one create/edit authoring architecture with an explicit target union, not separate resource models;
3. `FOLLOW_SOURCE | EXPLICIT_OVERRIDE` field/value variants and typed source-candidate references;
4. provider-qualified requirement/category/product-type resolution references without DTO/property-bag leakage;
5. listing-context authored-media intake/reference/selection/order/role without ProductAsset/media-master authority;
6. ListingIntent→PriceIntent typed correlation without embedding a price value;
7. `PriceIntent` target union `existing Listing | pre-creation ListingIntent context`, exact Money target and explicit supersession lineage;
8. fail-closed active-publication requirements for current PriceIntent + Availability owner input before provider dispatch;
9. Sellable Availability read semantics: desired owner value, provider actual evidence, knowledge/freshness and convergence without provider quantity becoming owner truth;
10. enough Inventory Source/allocation-policy schema to prove W2-A grammar under real updates;
11. negative controls proving price, availability, arbitrary provider fields and false source authority cannot enter ListingIntent.

W2-B still does **not** choose media blob upload/storage/hash/CDN mechanics; D7 owns those mechanics.

After W2-B, continue W2 through the remaining schema families needed for global coherence. When W2 converges as a whole, prepare one non-authoritative W2 review package and run Fable using the canonical Standard Fable review workflow before operator ratification/consolidation.

Implementation remains blocked until D9.

## 6. Fresh-session success test

A fresh session must conclude:

- D0→D4/D4-R1 and D5-B1 accepted/canonical;
- D5-B2 OPEN / ACTIVE with B2-A, Operation Matrix and Whole-Matrix accepted;
- Wire W1 and W2-A accepted in-stage;
- `D5-B2-W2-SCHEMA-GRAMMAR.md` is active authority for W2;
- W2-A exact values/typed refs/knowledge unions/typed PATCH/no universal wrappers are current;
- **W2-B ListingIntent + PriceIntent + Availability is exact next action**;
- Fable review happens after W2 coherently converges, not before;
- implementation remains blocked until D9.

If not, the active authority tree is inconsistent.