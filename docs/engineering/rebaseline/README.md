# Marketplace Central — Architecture Rebaseline

> **Role:** sole current-program status/router after `AGENTS.md`  
> **Current phase:** **D5 — API — OPEN / ACTIVE; D5-B1 ACCEPTED / CANONICAL; D5-B2 OPEN / ACTIVE**  
> **D5-B2 current state:** **B2-A + Operation Admission Matrix + Whole-Matrix ACCEPTED / RATIFIED; Wire W1 Resource/Path/HTTP Grammar ACCEPTED IN-STAGE; W2 Schema Grammar = NEXT**  
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
17. `docs/engineering/rebaseline/EVIDENCE-REGISTER.md`
18. code, OpenAPI, schemas, tests and runtime only as current-state evidence when needed

This file alone owns **where the program is and what happens next**. `ARCHITECTURE.md` owns stable cross-stage constraints; the Decision Reconciliation Baseline routes current decision generations; the ADR registry owns ADR status; accepted D-stage/B2 artifacts own detailed semantics in their scope.

`D5-API.md` remains D5-B1 authority. Never reconstruct target authority from memory, chat, Git history, retired ADRs, `AI-DIALOG.md`, deleted review candidates or current code/OpenAPI shape.

`AI-DIALOG.md` is a reusable non-authoritative review channel and currently contains no active review cycle. Completed review dialogue is Git-history evidence only.

## 2. Program state

```text
Documentary / governance cleanup — DONE
  ↓
D0 — Product / System Definition — CLOSED / ACCEPTED
  ↓
D1 — Domains / Boundaries — CLOSED / ACCEPTED
  ↓
D2 — Identity / Tenant / Data Ownership — CLOSED / ACCEPTED
  ↓
D3 — Communication / Events — CLOSED / ACCEPTED
  ↓
D4 — External Integrations — CLOSED / ACCEPTED AS A WHOLE
  └─ D4-R1 Publication Input & Listing Authoring — ACCEPTED / CANONICAL
       └─ R1-G1 ML initial publication × Availability — PASS-B
  ↓
Decision Reconciliation Baseline — ACCEPTED / CANONICAL
  ↓
D5 — API — OPEN / ACTIVE
  ├─ B1 Semantic API Model & Contract Laws — ACCEPTED / CANONICAL
  └─ B2 Product Operation / Resource Surface — OPEN / ACTIVE
       ├─ B2-A Client & Authentication Admission Model — ACCEPTED IN-STAGE
       ├─ Operation Admission Matrix Blocks 1–5 — ACCEPTED / RATIFIED
       ├─ Whole-Matrix Global Coherence — ACCEPTED / RATIFIED
       └─ Wire Contract / Resource-Path-Schema Grammar — OPEN / ACTIVE
            ├─ W1 Resource / Path / HTTP Grammar — ACCEPTED IN-STAGE
            └─ W2 Request/Response Schema + Knowledge/Outcome Grammar — NEXT
  ↓
D6 — Frontend
  ↓
D7 — Runtime / Jobs / Transactions
  ↓
D8 — Golden Flows
  ↓
D9 — Adversarial Architecture Review
  ↓
Implementation DAG / Plan
  ↓
Implementation
```

Product implementation remains blocked until D9 is accepted.

## 3. Accepted D5-B2 routing summary

### B2-A — Client/Auth

- Product API authentication uses one standards-based OIDC/OAuth boundary.
- Humans use Authorization Code + PKCE; confidential machine clients use Client Credentials/service-account semantics.
- MPC remains authority for Principal, Organization Membership, AccessRole/Permission/RoleAssignment and all business decisions.
- Tokens are audience-bound to MPC API; no global/static MPC Product API key or IdP-role business authority.
- Keycloak remains first implementation/proof candidate; D7 owns concrete provider/deployment/realm/secrets/token-lifetime realization.
- `GetCurrentAccessContext` is platform-scoped **self-only** discovery; Organization-owned business operations remain Organization-path-scoped.

### Operation Admission Matrix / Whole-Matrix — ratified

Load-bearing conclusions:

- Product API is semantic-owner driven, not CRUD-complete, screen-shaped or provider-shaped;
- no Product/PIM, generic Integration, Mutation/Action/Operation, Workflow, Rule, Finance, Task/Case or AI-specific business authority;
- `ListingIntent` is the create/edit listing authoring identity; listing-context authored media exists without Product/media mastery;
- price is **always** a distinct `PriceIntent`, including initial publication; pre-creation PriceIntent targets ListingIntent context;
- Sellable Availability remains Availability-owned; no public stock-set or baseline public AvailabilityIntent authoring;
- Market Intelligence owns competitive interpretation; Economics owns analysis/L0/L1/L2 and never actuates price;
- Governance owns authorization only; approval never executes;
- Sales is externally originated/read-centric;
- BusinessOrderIntent/InvoicingIntent are Materialization owner reactions, not direct client commands;
- Fulfillment exposes physical checkpoints/nodes/provider-readiness/internal targets without WMS/TMS/workflow ownership;
- generic automation cannot fabricate physical facts from Permission alone;
- Shipment remains source-qualified external observation;
- Post-Sale uses scoped canonical Resolution without provider action vocabulary becoming Product API;
- Work owns responsibility/lifecycle, not source truth; generic `SubmitWorkResolution` remains deferred;
- zero-P baseline until D6 proves repeated cross-owner composition need;
- no Product 1.0 bulk merely for symmetry.

Every admitted C operation has an explicit safety tuple:

```text
consequence class
idempotency disposition
concurrency / precondition disposition
```

No silent safety cell is allowed.

### W1 — Resource / Path / HTTP Grammar — accepted in-stage

W1 establishes:

- Product paths are relative to the OpenAPI server URL; no `/v1` baseline without an entitled compatibility consumer;
- Organization-owned Product API stays under `/organizations/{organization_id}/...`;
- `GET /access-context` is the only baseline platform-scoped self-only Product Q;
- MPC-owned canonical resources use opaque MPC IDs under Organization scope;
- externally authoritative resources retain explicit Marketplace Installation / SourceInstance qualification and do not receive synthetic mirror IDs for URI aesthetics;
- URI nesting means identity/lifecycle containment or source namespace qualification, never process order;
- D1 domain names and provider names are not Product API root topology;
- standard HTTP methods are used when resource semantics are honest;
- owner-specific non-CRUD capabilities use `POST {resource-uri}:verb`;
- writable `status` never substitutes for submit/resolve/physical-evidence/workflow capabilities;
- strong opaque MPC `ETag` + `If-Match` is the concurrency grammar where required;
- missing required precondition → `428 Precondition Required`; stale precondition → `412 Precondition Failed`, both as RFC 9457 Problem Details;
- `Idempotency-Key` remains independent from concurrency validators;
- exact OpenAPI minor version remains deferred until schemas/tooling needs are known.

W1 outcome: **CURRENT PARENT STRUCTURE CONFIRMED; identity-oriented resource paths + honest HTTP semantics selected.**

## 4. What is prohibited now

While **W2 is NEXT / NOT YET ACCEPTED**:

- do not begin D6–D9 target design or implementation;
- do not alter accepted D0–D4/D4-R1/D5-B1 or ratified B2/W1 meaning by schema convenience;
- do not derive schemas from current OpenAPI/provider DTOs/database rows/frontend forms;
- do not recreate Product/PIM, generic Integration/Mutation/Action/Operation/Workflow/Rules/AI authority, finance ledger, Task/Case engine or market collector platform;
- do not put price or Availability meaning inside ListingIntent because provider create payloads combine them;
- do not introduce universal `Fact<T>`, Evidence, Result, Operation, Resource or provider property-bag wrappers merely for schema uniformity;
- do not collapse unknown/unavailable/partial into `null`, zero, false, empty collections or HTTP success;
- do not expose bare external IDs, raw provider payloads or provider errors as Product API truth;
- do not turn semantic business rejection/pending/ambiguity into ordinary access/transport errors;
- do not re-admit generic Work resolution or cross-owner P without their ratified triggers;
- do not select D7 blob/generator/router/server technology from schema design convenience;
- do not add versioning/bulk/event stream without a real consumer.

## 5. Exact next action

**Derive D5-B2 Wire Contract W2 — Request / Response Schema Grammar + Knowledge / Outcome Semantics.**

W2 must establish, before collection/tooling work:

1. **Identity/reference schemas**
   - opaque MPC IDs;
   - Marketplace Installation-qualified external references;
   - SourceInstance-qualified business-system references;
   - same-Organization secondary-reference rejection;
   - no fake canonical ID for external/keyed meanings.

2. **Exact values**
   - Money wire representation and other exact D2 values needed by admitted operations;
   - no floating-point money convenience.

3. **Listing authoring schemas**
   - ListingIntent draft/create/update;
   - `FOLLOW_SOURCE | EXPLICIT_OVERRIDE` discriminated meaning;
   - listing-context media reference/intake;
   - PriceIntent correlation without embedding price in ListingIntent.

4. **PriceIntent schema**
   - target union `existing source-qualified Listing | pre-creation ListingIntent context`;
   - exact Money target;
   - explicit supersession lineage; no mutable PriceDraft baseline.

5. **Knowledge / freshness / provenance**
   - known / known-empty / unknown / unavailable / partial where material;
   - freshness orthogonal to knowledge;
   - owner/source provenance sufficient for the consumer;
   - no universal `Fact<T>` or Evidence business wrapper.

6. **Capability outcomes / Intent tracking**
   - accepted / rejected / pending / ambiguous;
   - applied/completed/converged distinctions where material;
   - owner-local resource/Intent state rather than generic `Operation` envelope.

7. **Update body grammar**
   - choose typed PATCH/update request shapes versus JSON Patch/Merge Patch based on real null/array/knowledge semantics;
   - no generic mutation DSL.

8. **Problems**
   - RFC 9457 base shape;
   - smallest stable MPC machine codes/extensions;
   - ordinary API problem vs valid domain outcome separation.

9. **Provider-enriched evidence**
   - bounded discriminated enrichment only for named consumer/correctness need;
   - unsupported/not-applicable honest on providers without the feature;
   - no raw provider DTO passthrough or universal provider field bag.

10. **Schema enforcement / negative controls**
    - invalid union;
    - bare external identity;
    - cross-Organization reference;
    - money precision/format loss;
    - knowledge collapse;
    - client-authored Principal/authority fields;
    - ListingIntent containing price/availability;
    - provider payload/property-bag leakage.

After W2, remaining Wire Contract work still includes collections/pagination/filter/search, complete Permission→operation mapping, exact per-operation safety/header/status table, technical non-Product ingress classification and OpenAPI/tooling spelling.

Do not choose D6 UI topology, D7 runtime/blob/transaction/Keycloak deployment, D8 live-effect proof or implementation.

If a schema cannot preserve ratified identity/owner/safety/knowledge semantics without distortion, reopen only the implicated decision rather than weakening the wire.

Implementation remains blocked until D9.

## 6. Fresh-session success test

A fresh session must conclude unambiguously:

- D0→D4 + D4-R1 accepted/canonical;
- Decision Reconciliation accepted/canonical;
- D5-B1 accepted/canonical;
- D5-B2 OPEN / ACTIVE;
- B2-A + Operation Matrix + Whole-Matrix are accepted/ratified;
- `D5-B2-WIRE-CONTRACT.md` is active authority for Wire Contract decisions;
- W1 Resource/Path/HTTP Grammar is accepted in-stage;
- Organization path scope, `/access-context`, source-qualified external identity, `:verb`, ETag/If-Match, 428/412 and no `/v1` baseline are current W1 decisions;
- **W2 Request/Response Schema Grammar + Knowledge/Outcome Semantics is the exact next action**;
- implementation remains blocked until D9.

If not, the active authority tree is inconsistent.
