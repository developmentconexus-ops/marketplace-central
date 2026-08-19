# Marketplace Central — Architecture Rebaseline

> **Role:** sole current-program status/router after `AGENTS.md`  
> **Current phase:** **D5 — API — OPEN / ACTIVE; D5-B1 ACCEPTED / CANONICAL; D5-B2 OPEN / ACTIVE**  
> **D5-B2 current state:** **Operation Matrix + Whole-Matrix RATIFIED; Wire W1 + W2 + W3 ACCEPTED / CANONICAL; W4 Permission → Operation / Client-Class Enforcement = NEXT**  
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
17. `docs/engineering/rebaseline/D5-B2-W2-SCHEMA-GRAMMAR.md` — canonical W2
18. `docs/engineering/rebaseline/D5-B2-W3-COLLECTION-GRAMMAR.md` — canonical W3
19. `docs/engineering/rebaseline/EVIDENCE-REGISTER.md`
20. code/OpenAPI/schemas/tests/runtime only as current-state evidence when needed

This file alone owns **program status, allowed/blocked work and exact next action**. Detailed semantics remain in the accepted artifacts above.

Former W3-C staging and Whole-W3 review candidate are absent from the active tree; Git history is the archive. `AI-DIALOG.md` is protocol-only and is not architecture authority.

`docs/engineering/rebaseline/cockpit.html` is a **non-authoritative visual projection**. It never participates in this authority path and must be synchronized only after canonical status changes.

Legacy/current code, routes, OpenAPI, SDK and frontend tables remain evidence only until later target work replaces them.

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
    B2-A Client/Auth                                     ACCEPTED IN-STAGE
    Operation Admission Matrix                           ACCEPTED / RATIFIED
    Whole-Matrix Global Coherence                        ACCEPTED / RATIFIED
    Wire Contract
      W1 Resource / Path / HTTP Grammar                  ACCEPTED / CANONICAL
      W2 Request / Response Schema Grammar               ACCEPTED / CANONICAL
      W3 Collections / Pagination / Filter / Search /
         Cursor Grammar                                  ACCEPTED / CANONICAL
      W4 Permission → Operation / Client-Class
         Enforcement                                     NEXT
      Technical non-Product ingress classification       BLOCKED BY W4
      Final Problem/media consistency                    BLOCKED BY SEQUENCE
      Single OpenAPI wire authority/tooling decision     BLOCKED BY SEQUENCE
D6 — Frontend                                             BLOCKED BY D5
D7 — Runtime / Jobs / Transactions                        BLOCKED
D8 — Golden Flows                                         BLOCKED
D9 — Adversarial Architecture Review                      BLOCKED
Implementation                                            BLOCKED UNTIL D9
```

Whole-W2 and Whole-W3 adversarial review cycles are complete, operator-ratified and incorporated into canonical Wire artifacts. Review dialogue is archived in Git history and carries no active authority.

---

## 3. Load-bearing D5-B2 authority

### 3.1 Semantic/client authority

- Product API is semantic-owner driven, not CRUD-, screen- or provider-shaped;
- humans use OIDC Authorization Code + PKCE; confidential machines use Client Credentials/service-account semantics;
- MPC owns Principal, Membership, AccessRole/Permission and every business authority;
- `GET /access-context` is the only platform-scoped self-only Product Q baseline;
- Organization business API remains `/organizations/{organization_id}/...`;
- ordinary Permission, business disposition/Governance authorization and epistemic authority to establish physical facts are distinct;
- no Product/PIM, generic Integration/Action/Operation/Workflow/Rules/Task/Finance/AI owner;
- zero-P baseline remains until D6 consumer evidence proves a bounded composition need.

### 3.2 Operation authority

- the ratified admission matrix is the sole Product 1.0 operation inventory;
- ListingIntent, PriceIntent and Availability remain distinct through initial publication;
- BusinessOrderIntent/InvoicingIntent remain Materialization owner reactions, not direct Product commands;
- Governance authorizes but does not execute;
- Fulfillment physical facts require admitted human/proven physical-system authority and cannot be fabricated by ordinary automation;
- Post-Sale coordinates scoped consequences and does not gain provider-action vocabulary;
- Work owns responsibility/lifecycle, never source truth;
- every admitted C operation has explicit consequence/idempotency/precondition semantics.

### 3.3 Canonical W1

- no `/v1` compatibility axis without a real stable consumer;
- paths express Organization/canonical identity or source namespace, not D1 package names/workflow order;
- external Listing/Sale/Shipment keep source-qualified native identity; no mirror IDs;
- standard methods only when honest; owner capabilities use `POST {resource-or-keyed-subject-uri}:verb`;
- one opaque ETag revision authority per protected meaning;
- same-resource standard mutation uses HTTP `If-Match`; custom/reference revision proofs use typed ETag request data;
- ProductChannelCorrespondence remains keyed Readiness meaning with correspondence-scoped ETag; no synthetic ID or forced PUT/DELETE;
- Idempotency-Key remains duplicate-intake safety, independent of revision proof;
- provider/business-system protocol ingress remains outside Product API roots.

### 3.4 Canonical W2

- opaque IDs/typed source refs; exact decimal Money; explicit temporal meanings;
- request/read schemas are authority-separated and closed;
- `null` never carries unknown/unavailable/partial/not-applicable;
- knowledge uses smallest owner-specific unions, not universal `Fact<T>`;
- no generic Result/Evidence/ExternalRef/Subject/Scope/Policy/Workflow wrappers;
- ListingIntent historical dispatch/effect basis is append-only without a PublicationAttempt Product resource;
- PriceIntent remains separate, Availability separates desired/observation/convergence, and Market/Economics remain evidence-honest;
- FulfillmentExecution is the one durable Fulfillment lifecycle identity; Work closure never fabricates source truth;
- canonical revision/idempotency grammar is complete;
- W2 owns the single Product Problem Details catalog, including Whole-W3 `invalid-cursor` and `cursor-expired` additions;
- D4/D8 still must prove selected-lane N/A reread and User-Product conditional-requirement behavior before live convergence claims.

### 3.5 Canonical W3

- exactly **26/26** admitted List/Search Q operations have one collection home; zero List/Search operations were added by symmetry;
- responses are owner/operation-specific; no generic `Page<T>`, `data/metadata` or projection wrapper;
- ListItems are semantic subsets of the same owner meaning and may not invent list-only business conclusions;
- forward-only `limit?` + opaque `cursor?` + optional `next_cursor`; no page/offset/skip/previous cursor baseline;
- every continuation repeats the explicit semantic subject/search/filters; cursor carries continuation only; only `limit` may vary;
- `limit` is a positive finite requested maximum, no silent clamp; exact numeric default/max remains bounded `DEFER SAFELY` to final OpenAPI closure after concrete payload measurement;
- cursor is opaque, ephemeral, query/operation/Organization-bound, never authorization and never raw provider/database state;
- `invalid-cursor` and `cursor-expired` are deliberate HTTP 400 cases from the W2 catalog; no silent restart;
- traversal exhaustion, at-most-once deduplication, source enumeration completeness and owner knowledge completeness are independent claims;
- stable MPC/source-qualified member identity appears at most once per traversal; no universal snapshot/no-omission guarantee;
- no universal total count or caller-selectable sort baseline;
- typed operation-specific filters only; no generic filter/query language;
- Source Product Search is source-capability-backed and does not freeze universal tokenizer/case/diacritics/locale/fuzzy/vector semantics or justify a Product mirror/index;
- CompetitivePosition/ExpectedEconomics Lists enumerate currently known existing Listing subjects only; pre-listing reasoning stays explicit point/search/evaluation flow;
- SellableAvailability population is the current Offering target universe: genuine pre-creation `new_listing` draft/submitted intents plus currently known existing Listings, with unknown/unavailable preserved honestly;
- ComparableOffer cursor chains stay on one owner-local Market evaluation/acquisition basis where identity-less mutable price ordering requires it; no public Snapshot/Traversal resource or fabricated ComparableOffer ID;
- provider paging remains D4 protocol; cursor persistence/signing/index/cache/seen-set remains D7 mechanism.

---

## 4. Prohibited now

While W4 is next:

- do not begin technical ingress classification, final OpenAPI/tooling closure, D6–D9 target design or implementation;
- do not reopen accepted D0→D5-B1/B2/W1/W2/W3 for naming/style or implementation convenience;
- do not reconstruct deleted W2/W3 staging/review files as parallel authority;
- do not derive access enforcement from IdP roles, frontend visibility, provider roles or current middleware by inheritance;
- do not collapse ordinary Permission into Governance authorization/business disposition;
- do not let an ordinary automation Principal establish physical facts merely because it has a broad Permission;
- do not create a generic ACL/ReBAC/policy expression engine, permission DSL or provider-scope ontology;
- do not add Product operations/permissions by API symmetry;
- do not choose D7 Keycloak realm/client/deployment mechanics during W4;
- do not choose exact OpenAPI minor/generator or numeric pagination defaults yet.

---

## 5. Exact next action

**Derive W4 — exact Permission → Product operation / allowed-client-class enforcement matrix from the ratified operation inventory and B2-A client/auth authority.**

W4 must decide and adversarially stress-test, for every admitted Product operation:

1. exact ordinary Permission required at the wire boundary, including `authenticated` for the bounded self-only `/access-context` exception;
2. allowed client class: human, automation, system/physical-system or the exact bounded combination already admitted;
3. where client-class admission is stricter than possession of the ordinary Permission;
4. where physical/epistemic authority is additionally required and cannot be inferred from an automation/service credential;
5. where Governance/business disposition remains a separate consequential prerequisite rather than an ordinary access Permission;
6. same-Organization Membership/reference privacy and fail-closed behavior without cross-tenant existence leakage;
7. monotonic/fail-safe revocation cases already ratified versus version-protected ordinary mutations;
8. least-privilege splits such as `fulfillment.read != fulfillment.execute != fulfillment.manage` and `listing.manage != price.manage`;
9. whether any currently named Permission is too broad, redundant, missing or merely mirrors a D1 package rather than a real client capability;
10. whether human-only/machine-capable distinctions remain coherent with B2-A OIDC/Client-Credentials semantics and D2 Principal kinds;
11. negative controls preventing IdP role claims, provider roles/scopes, frontend routes or legacy middleware from becoming MPC Permission authority.

W4 must **not** choose Keycloak realm/client configuration, token claims mapping, database ACL/RLS, middleware package topology or deployment. Those are D7 realization after the Product permission contract is closed.

If W4 exposes a real operation/client-class contradiction, reopen only the smallest implicated B2 operation admission or parent authority. Do not proceed to technical ingress classification until W4 is coherent.

After W4, continue remaining Wire obligations in router order:

1. technical non-Product ingress classification;
2. final Problem/media consistency as still needed;
3. one machine-readable OpenAPI authority and tooling/minor-version decision.

Implementation remains blocked until D9.

---

## 6. Fresh-session success test

A fresh session must conclude:

- D0→D4/D4-R1 + D5-B1 are accepted/canonical;
- D5-B2 Operation Matrix + Whole-Matrix are ratified;
- W1, W2 and **W3 are canonical**;
- W3-C staging and Whole-W3 review candidate are absent from the active tree and preserved only in Git history;
- `AI-DIALOG.md` has no active Whole-W3 review cycle;
- the cockpit is non-authoritative and may be used only as orientation;
- **W4 Permission → operation / client-class enforcement is the exact next action**;
- technical ingress/OpenAPI/D6–D9 remain blocked by sequence;
- implementation remains blocked until D9.

If not, the active authority tree is inconsistent.
