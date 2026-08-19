# Marketplace Central — Architecture Rebaseline

> **Role:** sole current-program status/router after `AGENTS.md`  
> **Current phase:** **D5 — API — OPEN / ACTIVE; D5-B1 ACCEPTED / CANONICAL; D5-B2 OPEN / ACTIVE**  
> **D5-B2 current state:** **Operation Matrix + Whole-Matrix RATIFIED; Wire W1 + W2-A/B/C ACCEPTED IN-STAGE; W2-D operational lifecycle schemas = NEXT**  
> **Implementation:** **BLOCKED until D9 is accepted**  
> **Last updated:** 2026-08-18

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
16. `docs/engineering/rebaseline/D5-B2-WIRE-CONTRACT.md`
17. `docs/engineering/rebaseline/D5-B2-W2-SCHEMA-GRAMMAR.md`
18. `docs/engineering/rebaseline/D5-B2-W2-C-READINESS-MARKET-ECONOMICS.md`
19. `docs/engineering/rebaseline/EVIDENCE-REGISTER.md`
20. code/OpenAPI/schemas/tests/runtime only as current-state evidence when needed

This file alone owns **program status, allowed/blocked work and exact next action**. Detailed semantics remain in the accepted artifacts above. `AI-DIALOG.md`, Git history, deleted review candidates, legacy routes/OpenAPI and current code are never target/status authority by inheritance.

`D5-B2-W2-SCHEMA-GRAMMAR.md` owns W2-A/B. `D5-B2-W2-C-READINESS-MARKET-ECONOMICS.md` owns W2-C only. Its parent file's older “W2-C next” wording is a pre-split snapshot; this router is current status authority.

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
    B2-A Client/Auth                                     ACCEPTED IN-STAGE
    Operation Admission Matrix                           ACCEPTED / RATIFIED
    Whole-Matrix Global Coherence                        ACCEPTED / RATIFIED
    Wire Contract
      W1 Resource / Path / HTTP Grammar                  ACCEPTED IN-STAGE
      W2 Request / Response Schema Grammar               OPEN / ACTIVE
        W2-A Core Schema Grammar                         ACCEPTED IN-STAGE
        W2-B ListingIntent / PriceIntent / Availability  ACCEPTED IN-STAGE
        W2-C Readiness / Market / Economics              ACCEPTED IN-STAGE
        W2-D Governance / Sales / Materialization /
             Fulfillment / Post-Sale / Work              NEXT
      Fable coherent W2 review                           AFTER W2 CONVERGES
D6 — Frontend                                             BLOCKED BY D5
D7 — Runtime / Jobs / Transactions                        BLOCKED
D8 — Golden Flows                                         BLOCKED
D9 — Adversarial Architecture Review                      BLOCKED
Implementation                                            BLOCKED UNTIL D9
```

## 3. Load-bearing D5-B2 decisions

### Client / access

- humans: OIDC Authorization Code + PKCE;
- confidential machines: Client Credentials/service-account semantics;
- MPC owns Principal, Organization Membership, AccessRole/Permission/RoleAssignment and every business authority;
- `GET /access-context` is platform-scoped self-only discovery;
- Organization-owned Product API remains `/organizations/{organization_id}/...`;
- generic automation cannot establish a physical fact merely by holding ordinary Permission.

### Operation inventory / authority

- Product API is semantic-owner driven, not CRUD-, screen- or provider-shaped;
- no Product/PIM, generic Integration/Action/Operation/Workflow/Rules/Task/Finance/AI authority;
- ListingIntent, PriceIntent and Availability remain distinct through initial publication;
- BusinessOrderIntent/InvoicingIntent are Materialization owner reactions, not direct client commands;
- Governance authorizes but never executes;
- Work owns responsibility/lifecycle, never source truth;
- generic Work resolution and cross-owner P remain deferred;
- every admitted C operation has explicit consequence/idempotency/precondition disposition.

### W1 — wire structure

- no `/v1` baseline without a real compatibility consumer;
- MPC IDs are opaque; external identities remain Marketplace Installation/SourceInstance-qualified;
- URI nesting means identity containment/namespace qualification, never workflow order;
- standard HTTP methods only when honest; non-CRUD owner capabilities use `POST {resource-uri}:verb`;
- writable `status` never substitutes for submit/resolve/physical-evidence semantics;
- strong opaque MPC `ETag` + `If-Match`; missing required precondition = 428, stale = 412;
- `Idempotency-Key` is separate from concurrency.

### W2-A — core schema laws

- exact money/decision decimals use decimal strings; `Money = ExactDecimalString + currency`;
- external refs are typed/source-qualified; no universal `ExternalRef`/entity graph;
- request objects are closed and separate from server-owned response/history schemas;
- `null` never carries unknown/unavailable/partial;
- knowledge uses smallest owner-specific exclusive unions, not universal `Fact<T>`;
- no generic Result/Operation/Evidence/property-bag wrappers;
- typed PATCH, omitted=unchanged; no generic JSON Patch/Merge Patch baseline;
- RFC 9457 `type` is primary API problem identifier;
- provider enrichment is bounded/owner-local, never raw DTO passthrough.

### W2-B — authoring / price / availability

- one sparse/declarative ListingIntent covers create/edit via target union;
- dynamic publication uses Readiness requirement keys/revisions/candidates and only `FOLLOW_SOURCE | EXPLICIT_OVERRIDE`;
- bounded PublicationValue; no arbitrary JSON/provider field bag;
- listing media = source candidate or ListingIntent-scoped authored media; no ProductAsset/media master;
- PriceIntent always owns price, including pre-creation price; supersession is explicit;
- client does not PATCH PriceIntent/Availability correlation into ListingIntent;
- `SubmitListingIntent` has no business payload and revalidates current dependencies;
- SellableAvailability separates control, desired value, provider observation and convergence;
- known zero is known zero; provider unavailability never erases known desired state.

### W2-C — readiness / market / economics

- source Product search is source-qualified evidence, not an MPC Product resource;
- ProductChannelReadiness and CompetitivePosition are contextual keyed Qs with no synthetic canonical IDs;
- PublicationRequirements returns effective Readiness-owned requirements/candidates/options, not provider rules DSL;
- market coverage and evidence sufficiency are separate; provider enumeration never means universal market completeness;
- ComparableOffer values gain no fake identity; provider-rich competitive evidence remains bounded and source-qualified;
- ExpectedEconomics is components-first and only exposes a known conclusion when required cost/tax/fee/shipping evidence is sufficiently known;
- EvaluatePriceScenario accepts hypothetical variables, never authoritative evidence overrides;
- SaleEconomics preserves L0/L1/L2 and R1/R2; realized economics is occurrence-based and honest about partial coverage;
- Economic Attribution remains exact/partial/ambiguous/unresolved under an Economics-local bounded subject union; no universal entity/reconciliation graph;
- Commercial Policy stays typed Economics configuration; no generic rules DSL.

## 4. Prohibited now

While W2 is open:

- do not begin D6–D9 target design or implementation;
- do not weaken accepted D0→D5-B1/B2/W1/W2-A/B/C for schema convenience;
- do not derive schemas from legacy OpenAPI, provider DTOs, database rows or frontend forms;
- do not introduce Product mirror/PIM, generic Result/Fact/Evidence/Operation/ExternalRef/property-bag/rules/workflow abstractions;
- do not collapse unknown/unavailable/partial/not-applicable into null/zero/false/empty;
- do not expose bare native IDs, raw provider payloads/errors or client-authored effective authority fields;
- do not fabricate market completeness or profitability from incomplete evidence;
- do not choose D7 server/generator/blob/persistence/queue/transaction/Keycloak realization;
- **do not run Fable yet**: W2 receives one coherent independent review after its internal sections converge unless a material contradiction requires a focused round.

## 5. Exact next action

**Derive W2-D — Governance + Marketplace Sales + Business-System Materialization + Fulfillment + Post-Sale + Operational Work schema grammar.**

W2-D must challenge consequential operational lifecycles and decide proportionately:

1. Authorization Decision + Delegation schemas without `approved=true`, rules-engine or execution authority;
2. Marketplace Sale read/attribution-resolution schema without client-created Sale/provider Order mirror;
3. BusinessOrderIntent/InvoicingIntent tracking schemas whose creation remains owner-triggered;
4. Party Resolution candidates/resolution state + Destination Realization without Customer/Address master authority;
5. Fulfillment state/checkpoint/Node/Artifact/Shipment schemas without generic workflow/WMS/TMS/provider DTO mirror;
6. physical-evidence establishment schema/client constraints that cannot be forged by ordinary automation;
7. Post-Sale Resolution scoped consequences without one cancellation/return/refund status or provider action vocabulary;
8. Work responsibility/assignment/hold/escalation without becoming source truth or Task/Case platform;
9. owner-specific pending/ambiguous/effect/convergence semantics without generic Operation result;
10. negative controls proving client/provider/workflow payloads cannot bypass owner/Governance/source truth.

After W2-D, complete the remaining cross-owner W2 configuration/problem/outcome consistency work. When W2 converges, prepare one non-authoritative W2 review package and run Fable through the canonical Standard Fable workflow before operator ratification/consolidation.

Implementation remains blocked until D9.

## 6. Fresh-session success test

A fresh session must conclude:

- D0→D4/D4-R1 and D5-B1 accepted/canonical;
- D5-B2 Operation Matrix + Whole-Matrix ratified;
- W1 and W2-A/B/C accepted in-stage;
- W2-C is owned by its dedicated subartifact and does not overlap W2-A/B authority;
- **W2-D operational lifecycle schema grammar is exact next action**;
- Fable waits until W2 coherently converges;
- implementation remains blocked until D9.

If not, the active authority tree is inconsistent.