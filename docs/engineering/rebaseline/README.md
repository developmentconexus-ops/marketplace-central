# Marketplace Central — Architecture Rebaseline

> **Role:** sole current-program status/router after `AGENTS.md`  
> **Current phase:** **D5 — API — OPEN / ACTIVE; D5-B1 ACCEPTED / CANONICAL; D5-B2 OPEN / ACTIVE**  
> **D5-B2 current state:** **Operation Matrix + Whole-Matrix RATIFIED; Wire W1 + W2 + W3 + W4 ACCEPTED / CANONICAL; Technical non-Product ingress classification = NEXT**  
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
19. `docs/engineering/rebaseline/D5-B2-W4-PERMISSION-CLIENT-CLASS-ENFORCEMENT.md` — canonical W4
20. `docs/engineering/rebaseline/EVIDENCE-REGISTER.md`
21. code/OpenAPI/schemas/tests/runtime only as evidence when needed

This file alone owns **program status, allowed/blocked work and exact next action**. Detailed semantics remain in the accepted artifacts above.

Former Whole-W2/W3/W4 review candidates/staging are absent from the active authority tree; Git history is the archive. `AI-DIALOG.md` is protocol-only and is not architecture authority.

`docs/engineering/rebaseline/cockpit.html` is a **non-authoritative visual projection**. It never participates in this authority path and is synchronized only after canonical status changes.

Legacy/current code, OpenAPI, IdP roles/scopes, middleware, provider roles/scopes and frontend guards remain evidence only until later target work replaces them.

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
      W3 Collections / Query / Cursor Grammar            ACCEPTED / CANONICAL
      W4 Permission / Client-Class Enforcement           ACCEPTED / CANONICAL
      Technical non-Product ingress classification       NEXT
      Final Problem/media consistency                    BLOCKED BY SEQUENCE
      Single OpenAPI wire authority/tooling decision     BLOCKED BY SEQUENCE
D6 — Frontend                                             BLOCKED BY D5
D7 — Runtime / Jobs / Transactions                        BLOCKED
D8 — Golden Flows                                         BLOCKED
D9 — Adversarial Architecture Review                      BLOCKED
Implementation                                            BLOCKED UNTIL D9
```

Whole-W2, Whole-W3 and Whole-W4 adversarial review cycles are complete, operator-ratified and incorporated into canonical authority. Review dialogue is archived in Git history and carries no active authority.

---

## 3. Load-bearing current authority

### 3.1 D2 identity/access confirmation carried by Whole-W4

D2 now explicitly confirms **current Principal access eligibility** as Principal-scoped revocable identity/access state.

Disabling/revoking that eligibility blocks future Product access, including `/access-context`, without deleting Membership/RoleAssignment or rewriting historical actor attribution. Exact lifecycle/storage/cache/revocation mechanism remains D7.

D2 Principal kinds remain exactly:

```text
human | automation | system
```

No `physical_system`, service-account, agent or integration Principal kind is introduced by convenience.

### 3.2 Canonical W1

- Organization business paths use `/organizations/{organization_id}/...`; `/access-context` is the one bounded platform-scoped self-only Q.
- paths express canonical identity/source namespace, not D1 package names/workflow order.
- external Listing/Sale/Shipment keep source-qualified native identity; no mirror IDs.
- honest standard methods plus owner-specific `POST ...:verb` capabilities; no generic actions/commands escape hatch.
- one opaque ETag revision authority per protected meaning.
- same-resource standard mutation uses HTTP `If-Match`; custom/reference revision proofs use typed ETag request data.
- Idempotency-Key remains duplicate-intake safety, independent from revision proof.
- provider/business-system protocol ingress remains outside Product API roots.

### 3.3 Canonical W2

- opaque IDs/typed source refs; exact decimal Money; explicit temporal and knowledge meanings.
- request/read schemas are authority-separated and closed; no generic Result/Evidence/Subject/Scope/Workflow ontology.
- ListingIntent / PriceIntent / Availability remain distinct; historical dispatch/effect basis remains explainable.
- W2 owns the single Product Problem Details catalog, including `invalid-cursor` and `cursor-expired`.
- business rejection/unknown/unavailable/external-required stays Product semantics rather than generic HTTP access failure.
- mutation/capability operations return the current relevant owner representation under the operation's own contract; W4 defines access scope for that response disclosure.

### 3.4 Canonical W3

- exactly **26/26** admitted List/Search Q operations have one collection home; zero additions by symmetry.
- owner-specific collection responses; no generic `Page<T>`/query/filter/sort/projection framework.
- forward `limit?` + opaque `cursor?` + optional `next_cursor`; explicit semantic query repeats on continuation.
- no total-count/caller-sort/universal snapshot baseline.
- stable identity at most once per traversal; exhaustion/coverage/completeness remain separate claims.
- Source Product Search is source-capability-backed without forcing Product mirror/index or universal tokenizer/case/locale/vector semantics.
- CompetitivePosition/ExpectedEconomics/SellableAvailability populations are explicitly bounded and honest.
- ComparableOffer continuation remains on one owner-local evaluation/acquisition basis when identity-less mutable ordering requires it.

### 3.5 Canonical W4

W4 maps **95/95** admitted Product operations using **29 flat exact stored Permissions** with zero hierarchy/wildcard implications.

Canonical access sequence for Organization-owned operations:

```text
external AuthN
→ Principal binding (zero=401; duplicate=500)
→ current Principal access eligibility
→ current Organization Membership
→ allowed Principal kind
→ exact ordinary Permission
→ current server-owned physical qualification, only where required
→ same-Organization resource/reference privacy
→ W1/W2 safety
→ owner business disposition
→ Governance when applicable
→ owner intake/effect
```

Key W4 laws:

- ordinary `both` Q/read → H/A/S; consequential business C admitted to `both` → H/A; `EvaluatePriceScenario` → H/A/S.
- `system` does not substitute for business `automation`.
- Permissions are flat/exact: no `manage ⇒ read`, `execute ⇒ read`, prefix or wildcard implication; AccessRoles may bundle exact Permissions.
- a mutation Permission covers the normal W1/W2 response for the **same operation subject** but does not grant corresponding Get/List/Search operations or Permission inheritance.
- `/access-context` requires valid AuthN + exactly one Principal binding + current Principal eligibility, but no Organization Membership/Product Permission; eligible zero-membership Principal receives an empty Organization set.
- no-Membership path Organization → privacy-preserving 404; current Membership but missing Permission/client-class/required qualification → 403.
- physical checkpoint baseline: Separation/Packing = H only; PhysicalConference/DispatchHandoff = H or currently qualified S. Qualification is server-owned, operation-specific epistemic authority, not Permission/Principal kind/Governance.
- every mutable/revocable access fact must be current; stale token/cache/provisioning state never survives revocation as authority.
- sensitive read choices are explicit: FulfillmentArtifacts use `fulfillment.execute`; AuthorizationDelegations use H + `governance.manage`; Materialization party/destination reads make `materialization.read` PII-bearing. No speculative split Permission is added.
- Governance cannot widen Permission/client-class/physical authority; IdP/OAuth/provider roles/scopes/frontend guards never independently grant Product access.

---

## 4. Prohibited now

While technical non-Product ingress classification is next:

- do not begin final Problem/media consistency, OpenAPI/tooling closure, D6–D9 target design or implementation;
- do not reopen accepted D0→D5-B1/B2/W1/W2/W3/W4 for naming/style or implementation convenience;
- do not reconstruct deleted review/staging files as parallel authority;
- do not derive ingress/Product access from current controllers/routes/middleware/provider endpoints by inheritance;
- do not create a generic Integration/Ingress business domain, provider-shaped Product API root or generic webhook/event resource;
- do not let provider OAuth credentials/webhook signatures authenticate Product Principals;
- do not let Product Permission/Governance become provider-protocol authentication/verification authority;
- do not choose D7 queue/web-server/middleware/secret storage/retry/deployment mechanics yet;
- do not choose exact OpenAPI minor/generator or numeric pagination defaults yet.

---

## 5. Exact next action

**Derive the technical non-Product ingress classification required by D5-B1/D4 without creating Product business operations or provider-shaped Product API authority.**

The next decision must classify every Product 1.0 inbound technical surface that is **not a Product Principal invoking a Product API operation**, including proportionately:

1. marketplace OAuth begin/callback/authorization-result ingress where applicable;
2. marketplace/provider webhook/notification ingress;
3. business-system/provider callback/event ingress where selected flows require it;
4. whether any internal owner-trigger/runtime signal belongs here or remains entirely D3/D7 internal mechanism;
5. which authority authenticates/verifies each ingress: provider protocol credential/signature/state, never MPC Product Permission by convenience;
6. Organization/MarketplaceInstallation/SourceInstance correlation rules without trusting caller-supplied tenant/business authority;
7. pointer/event semantics versus authoritative reread — notification never becomes provider business truth when D4 requires reread;
8. duplicate delivery/replay/out-of-order/unknown-resource handling at the semantic ingress boundary without prematurely choosing queue/storage mechanics;
9. response/acknowledgement meaning so transport acknowledgement is never mistaken for business acceptance/convergence;
10. strict separation from Product API paths, Product Principal authentication and owner Product operations;
11. whether any ingress requires a bounded public technical route and how it avoids becoming `/providers`/`/integrations` Product vocabulary;
12. negative controls preventing provider DTO/status/auth topology from leaking into D1 business ontology or Product API authority.

This work must consume accepted D3/D4/D5 authority and classify only the smallest technical ingress surfaces required by selected Product 1.0 flows.

If a real ingress cannot be represented without changing owner meaning, authentication authority or D3/D4 communication assumptions, reopen only the smallest implicated parent decision.

After ingress classification, continue remaining Wire obligations in router order:

1. final Problem/media consistency as still needed;
2. one machine-readable OpenAPI authority and tooling/minor-version decision.

Implementation remains blocked until D9.

---

## 6. Fresh-session success test

A fresh session must conclude:

- D0→D4/D4-R1 + D5-B1 are accepted/canonical;
- D5-B2 Operation Matrix + Whole-Matrix are ratified;
- W1, W2, W3 and **W4 are canonical**;
- D2 includes current Principal access eligibility as revocable identity/access authority;
- W4 maps 95/95 operations with 29 flat exact Permissions and Whole-W4 review is archived;
- Whole-W4 candidate is absent from the active tree and `AI-DIALOG.md` has no active review cycle;
- cockpit is non-authoritative orientation only;
- **technical non-Product ingress classification is the exact next action**;
- final Problem/media/OpenAPI/D6–D9 remain blocked by sequence;
- implementation remains blocked until D9.

If not, the active authority tree is inconsistent.
