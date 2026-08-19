# Marketplace Central — Architecture Rebaseline

> **Role:** sole current-program status/router after `AGENTS.md`  
> **Current phase:** **D5 — API — OPEN / ACTIVE; D5-B1 ACCEPTED / CANONICAL; D5-B2 OPEN / ACTIVE**  
> **D5-B2 current state:** **B2-A + Operation Matrix + Whole-Matrix ACCEPTED / RATIFIED; Wire W1 + W2-A + W2-B ACCEPTED IN-STAGE; W2-C Readiness + Market Intelligence + Economics schemas = NEXT**  
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

`D5-B2-WIRE-CONTRACT.md` owns accepted Wire W1. `D5-B2-W2-SCHEMA-GRAMMAR.md` owns the active W2 schema batch. W2 receives one coherent independent Fable review only after its sections converge; no micro-review cadence is implied.

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
        W2-B ListingIntent + PriceIntent + Availability ACCEPTED IN-STAGE
        W2-C Readiness + Market Intelligence + Economics = NEXT
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

- MPC/native IDs are opaque strings; external references are typed/source-qualified, never universal `ExternalRef`;
- `Money = ExactDecimalString + explicit currency`; authoritative JSON floating-point money and universal minor units are rejected;
- request and response schemas are separate where authority differs; semantic request objects are closed;
- `null` never transports unknown/unavailable/partial;
- knowledge uses smallest owner-specific exclusive unions, not universal `Fact<T>`;
- provenance/freshness attaches to the smallest semantic unit for which it is true;
- no generic Result/Operation wrapper; valid business pending/rejected/ambiguous remains semantic meaning;
- typed owner-specific PATCH; omitted=unchanged; no generic JSON Patch/Merge Patch baseline;
- RFC 9457 `type` is primary problem identifier; no duplicate problem-code taxonomy by default;
- provider enrichment is bounded owner-local closed schema;
- clients cannot author effective Principal/Organization/approval/convergence/history.

### W2-B — accepted

- one ListingIntent identity covers create/edit through a target union; drafts may be contract-valid while business-incomplete;
- ListingIntent is sparse/declarative Offering meaning, not a full provider mirror;
- dynamic publication uses Readiness-owned requirement keys/revisions/candidates resolved by `FOLLOW_SOURCE | EXPLICIT_OVERRIDE` only;
- `PublicationValue` is a small closed union (`text`, `exact_decimal`, `boolean`, `option`, `text_list`, `option_list`); arbitrary objects/property bags are not baseline;
- requirement/source/media candidate keys are context-bound opaque references, not new global entities/provider IDs;
- media selection is source-media or ListingIntent-scoped authored-media; array order is publication order; no ProductAsset/media master;
- PriceIntent always owns price and targets either an existing source-qualified Listing or pre-creation ListingIntent; price supersession is explicit and attributed;
- clients do not PATCH price/availability correlations into ListingIntent; server establishes/revalidates correlations through accepted owner semantics;
- `SubmitListingIntent` carries no business payload; it submits the current revision under `If-Match` and dispatch later revalidates required PriceIntent/Availability/authorization;
- SellableAvailability target works before/after creation and separates control/effective capability, desired owner quantity, provider observation and convergence;
- known zero availability remains known zero; unavailable provider observation never erases known desired value;
- Inventory Source remains MPC identity mapped to source-qualified inventory scopes, never native stock/location identity.

## 4. Prohibited now

While W2 is open:

- do not begin D6–D9 target design or implementation;
- do not alter accepted D0→D5-B1, B2, W1 or W2-A/B by schema convenience;
- do not derive schemas from legacy OpenAPI, provider DTOs, database rows or frontend forms;
- do not put price/Availability/provider payloads inside ListingIntent;
- do not create universal Fact/Evidence/Result/Operation/Resource/ExternalRef/property-bag schemas;
- do not collapse unknown/unavailable/partial/not-applicable into null/zero/false/empty;
- do not expose bare native IDs, raw provider payloads/errors or client-authored effective authority fields;
- do not create mutable PriceDraft, client-side ListingIntent correlation choreography or provider-mirror desired state;
- do not choose D7 blob/upload/server/generator/persistence mechanics during W2;
- do not run Fable yet: W2 review occurs once after the coherent W2 package converges unless a material contradiction forces a focused round.

## 5. Exact next action

**Derive W2-C — Product & Channel Readiness + Market Intelligence + Commercial Economics schema grammar.**

W2-C must decide and adversarially stress-test:

1. source Product search/reference result shapes without Product mirror identity;
2. ProductChannelReadiness keyed-Q shape, readiness conclusion and blockers without synthetic readiness ID;
3. PublicationRequirements / Requirement / Candidate / Option shapes used by W2-B, including requirement revision and source/media candidate references;
4. Market Intelligence CompetitivePosition and ComparableOffer shapes with source-qualified provider-rich evidence, comparability/evidence sufficiency, coverage and freshness without generic MarketObservation authority;
5. Expected Economics with component-level knowledge/provenance so incomplete tax/fee/shipping/cost evidence cannot masquerade as complete profitability;
6. `EvaluatePriceScenario` request/result where caller supplies legitimate hypothetical variables but never authoritative evidence replacements;
7. `SaleEconomics` L0/L1/L2 plus R1/R2 lineage without mutable Profitability/Reconciliation resource;
8. Economic Attribution exact/partial/ambiguous/unresolved subject/reference grammar without universal entity graph;
9. only enough Commercial Policy schema to test W2-A update semantics; generic Rules DSL remains rejected;
10. negative controls proving clients/providers cannot fabricate source truth, market completeness or authoritative economic evidence.

W2-C does not yet choose collection pagination/filter grammar or full cross-owner policy default/inheritance/override grammar if a later configuration section can adjudicate that pattern more coherently.

After W2-C, continue the remaining W2 owner/schema families needed for global coherence. When W2 converges as a whole, prepare one non-authoritative W2 review package and run Fable using the canonical Standard Fable review workflow before operator ratification/consolidation.

Implementation remains blocked until D9.

## 6. Fresh-session success test

A fresh session must conclude:

- D0→D4/D4-R1 and D5-B1 accepted/canonical;
- D5-B2 OPEN / ACTIVE with B2-A, Operation Matrix and Whole-Matrix accepted;
- Wire W1 and W2-A/W2-B accepted in-stage;
- `D5-B2-W2-SCHEMA-GRAMMAR.md` is active authority for W2;
- W2-A exact values/typed refs/knowledge unions/typed PATCH/no universal wrappers remain current;
- W2-B keeps ListingIntent, PriceIntent and Availability distinct through initial publication and uses server-established correlations;
- **W2-C Readiness + Market Intelligence + Economics is exact next action**;
- Fable review happens after W2 coherently converges, not before;
- implementation remains blocked until D9.

If not, the active authority tree is inconsistent.