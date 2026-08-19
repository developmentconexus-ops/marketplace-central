# Marketplace Central — Architecture Rebaseline

> **Role:** sole current-program status/router after `AGENTS.md`  
> **Current phase:** **D5 — API — OPEN / ACTIVE; D5-B1 ACCEPTED / CANONICAL; D5-B2 OPEN / ACTIVE**  
> **D5-B2 current state:** **Operation Matrix + Whole-Matrix RATIFIED; Wire W1 + W2-A/B/C/D ACCEPTED IN-STAGE; W2-E transversal/final schema consistency = NEXT**  
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
19. `docs/engineering/rebaseline/D5-B2-W2-D-OPERATIONAL-SCHEMAS.md`
20. `docs/engineering/rebaseline/EVIDENCE-REGISTER.md`
21. code/OpenAPI/schemas/tests/runtime only as current-state evidence when needed

This file alone owns **program status, allowed/blocked work and exact next action**. Detailed semantics remain in the accepted artifacts above. `AI-DIALOG.md`, Git history, deleted review candidates, legacy routes/OpenAPI and current code are never target/status authority by inheritance.

`D5-B2-W2-SCHEMA-GRAMMAR.md` owns W2-A/B. `D5-B2-W2-C-READINESS-MARKET-ECONOMICS.md` owns W2-C only. `D5-B2-W2-D-OPERATIONAL-SCHEMAS.md` owns W2-D only. Older “next” wording in child artifacts is a stage snapshot; this router is current status authority.

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
             Fulfillment / Post-Sale / Work              ACCEPTED IN-STAGE
        W2-E Transversal / final consistency             NEXT
      Whole-W2 coherence pass                            AFTER W2-E
      Fable coherent W2 review                           AFTER WHOLE-W2 PASS
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

### W2-A/B/C

- exact money/decision decimals use decimal strings; `Money = ExactDecimalString + currency`;
- external refs are typed/source-qualified; no universal `ExternalRef`/entity graph;
- request objects are closed and separate from server-owned response/history schemas;
- knowledge is owner-specific and explicit; `null`/zero/empty never carry uncertainty by convention;
- no generic Result/Operation/Evidence/property-bag/rules DSL;
- ListingIntent is sparse/declarative and resolves Readiness requirements through `FOLLOW_SOURCE | EXPLICIT_OVERRIDE` only;
- PriceIntent always owns price and Availability owns sellable quantity; server establishes/revalidates correlations;
- Readiness returns bounded source/requirement candidates without Product mirror authority;
- Market Intelligence keeps coverage distinct from evidence sufficiency and never claims universal market completeness;
- Economics is components-first, preserves L0/L1/L2/R1/R2 and never fabricates profitability from missing evidence;
- Economic Attribution remains a bounded Economics-local polymorphic subject, not universal entity graph.

### W2-D — operational schemas

- AuthorizationDecision is immutable occurrence; its target is a closed union and exact reviewed target revision is preserved;
- AuthorizationDelegation now has a justified stable opaque wire ID for update/revoke, without becoming generic Grant/IAM engine;
- Sale and Shipment remain source-qualified external identities; `sale_line_key` is only a Sale-scoped selector/correlation key;
- BusinessOrderIntent/InvoicingIntent remain owner-triggered read/tracking resources; native business/fiscal results remain SourceInstance-qualified;
- PartyResolution/DestinationRealization are BusinessOrderIntent-contained singleton meanings; Party resolve never accepts arbitrary Customer master fields;
- `FulfillmentExecution` is the justified Fulfillment-owned durable identity for admitted physical checkpoint addressing; no new domain/WMS/workflow authority;
- physical checkpoint occurrence is distinct from physical-readiness conclusion; client cannot author effective actor/trusted-evidence claims;
- PostSaleResolution has explicit sale/line/quantity scope and multiple simultaneous consequence tracks; no direct close/provider-action vocabulary;
- Work origin is a Work-local closed union, responsibility role is not AccessRole, and assignment/hold/resume/escalation never become source truth;
- Work closure-path audit passes for currently proven Product 1.0 condition classes; generic `SubmitWorkResolution` remains deferred;
- W2-D identified the need for a typed **referenced-resource precondition** when a POST/create/capability depends on the exact revision of a different MPC resource.

## 4. Prohibited now

While W2 remains open:

- do not begin D6–D9 target design or implementation;
- do not weaken accepted D0→D5-B1/B2/W1/W2-A/B/C/D for schema convenience;
- do not derive schemas from legacy OpenAPI, provider DTOs, database rows or frontend forms;
- do not introduce Product mirror/PIM, generic Result/Fact/Evidence/Operation/ExternalRef/property-bag/rules/workflow abstractions;
- do not collapse unknown/unavailable/partial/not-applicable into null/zero/false/empty;
- do not expose bare native IDs, raw provider payloads/errors or client-authored effective authority fields;
- do not expose TOP/NUNOTA/CODPARC/provider status as canonical Product semantics;
- do not allow generic automation to fabricate physical facts;
- do not create direct Post-Sale/Work close or generic Work resolution;
- do not choose D7 server/generator/blob/persistence/queue/transaction/Keycloak realization;
- **do not run Fable yet**: W2 receives one coherent review after W2-E + Whole-W2 coherence pass unless a material contradiction forces a focused round.

## 5. Exact next action

**Derive W2-E — transversal/final W2 schema consistency.**

W2-E must decide once, across the W2 package:

1. **policy/config grammar** — deterministic default/inherited/effective + explicit override for Availability allocation, Commercial Economics policy and Fulfillment operating targets, while meaning remains owner-local and no generic Rules/SLA platform appears;
2. **capability/business outcome grammar** — owner-local pending/rejected/ambiguous/external-effect/convergence semantics without a universal Result/Operation state;
3. **referenced-resource precondition** — exact wire representation for create/capability operations whose correctness depends on a different MPC resource revision;
4. **safety-axis coherence** — relationship among `ETag`/`If-Match`, `Idempotency-Key`, referenced-resource validators and semantic/provider prerequisites;
5. **Problem Details** — exact RFC 9457 problem types/extensions required by the admitted Product contract, without duplicate global error taxonomy;
6. **response/status/body grammar** — creation/read/PATCH/`:verb` success behavior, valid business outcomes versus transport/access/conditional failures and no ritual `202`;
7. **request closure/authority fence** — final cross-owner negative controls against undeclared/provider/effective-authority fields;
8. **Whole-W2 pre-review audit** — duplicates, missing schema path, inconsistent state names, hidden universal wrappers, operations that still cannot be expressed faithfully.

After W2-E converges, run a Whole-W2 coherence pass. If no material contradiction survives, prepare one disposable non-authoritative W2 review candidate and run Fable through the canonical Standard Fable workflow. Reviewer output remains evidence; operator ratifies before canonical W2 consolidation/next wire sub-batch.

Implementation remains blocked until D9.

## 6. Fresh-session success test

A fresh session must conclude:

- D0→D4/D4-R1 and D5-B1 accepted/canonical;
- D5-B2 Operation Matrix + Whole-Matrix ratified;
- W1 and W2-A/B/C/D accepted in-stage;
- W2-C/D are owned by their dedicated subartifacts and do not create parallel global status authority;
- W2-D crystallized AuthorizationDelegationId and FulfillmentExecutionId only because real wire consumers now require stable identity;
- Work closure-path audit currently passes and generic Work resolution remains deferred;
- **W2-E transversal/final consistency is exact next action**;
- Whole-W2 coherence then one Fable review, not before;
- implementation remains blocked until D9.

If not, the active authority tree is inconsistent.
