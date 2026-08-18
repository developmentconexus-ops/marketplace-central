# Marketplace Central — Architecture Rebaseline

> **Role:** sole current-program status/router after `AGENTS.md`  
> **Current phase:** **D5 — API — OPEN / ACTIVE; D5-B1 ACCEPTED / CANONICAL; D5-B2 OPEN / ACTIVE**  
> **D5-B2 current state:** **B2-A + Operation Admission Matrix + Whole-Matrix Global Coherence ACCEPTED IN-STAGE / OPERATOR-RATIFIED; Wire Contract = NEXT**  
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
16. `docs/engineering/rebaseline/EVIDENCE-REGISTER.md`
17. code, OpenAPI, schemas, tests and runtime only as current-state evidence when needed

This file alone owns **where the program is and what happens next**. `ARCHITECTURE.md` owns stable cross-stage constraints; the Decision Reconciliation Baseline routes current decision generations; the ADR registry owns ADR status; accepted D-stage/B2 artifacts own detailed semantics in their scope.

`D5-API.md` remains D5-B1 authority. Its old next-action wording is a pre-B2-opening snapshot. Never reconstruct target authority from memory, chat, Git history, retired ADRs, `AI-DIALOG.md`, deleted review candidates or current code/OpenAPI shape.

`AI-DIALOG.md` is a reusable non-authoritative review channel and currently contains no active review cycle. The completed Whole-Matrix Fable/GPT dialogue and deleted review candidate are archived in Git history only.

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
       ├─ Operation Admission Matrix Blocks 1–5 — ACCEPTED IN-STAGE
       ├─ Whole-Matrix Global Coherence — ACCEPTED / OPERATOR-RATIFIED
       │    ├─ independent Fable Rounds 1–2 — COMPLETED
       │    ├─ GPT adjudication — COMPLETED
       │    ├─ B2-local corrections — CONSOLIDATED
       │    ├─ review candidate — DELETED / GIT HISTORY ONLY
       │    └─ AI-DIALOG — RESET / NO ACTIVE REVIEW
       └─ Wire Contract / Resource-Path-Schema Grammar — NEXT
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
- Humans use Authorization Code + PKCE semantics; confidential machine clients use Client Credentials/service-account semantics.
- MPC remains authority for Principal, Organization Membership, AccessRole/Permission/RoleAssignment and all business decisions.
- Tokens are audience-bound to MPC API; no global/static MPC Product API key or IdP-role business authority.
- Keycloak remains first implementation/proof candidate; D7 owns concrete provider/deployment/realm/secrets/token-lifetime realization.
- `GetCurrentAccessContext` is a bounded platform-scoped **self-only** discovery Q; Organization-owned business operations remain Organization-path-scoped.

### Operation Admission Matrix — ratified

The Product API is semantic-owner driven, not CRUD-complete, screen-shaped or provider-shaped.

Load-bearing accepted conclusions:

- no Product/PIM, generic Integration, Mutation/Action/Operation, Workflow, Rule, Finance, Task/Case or AI-specific business authority;
- Listing actual state is Offering Q; `ListingIntent` is the single create/edit authoring identity;
- listing-context authored-media intake exists without a reusable Product/media master;
- price is **always** a distinct `PriceIntent`, including initial publication; pre-creation PriceIntent targets the ListingIntent context and may be physically serialized with Listing + Availability inputs by D4/D7 without ownership merge;
- Sellable Availability remains Availability-owned; no public stock-set or baseline public AvailabilityIntent authoring;
- Market Intelligence owns competitive interpretation; Commercial Economics owns stateless scenario analysis plus durable material L0/L1/L2 lineage; Economics never actuates price;
- Governance owns authorization decisions/delegations only; business thresholds stay with their owners and approval never executes an effect;
- Sales is externally originated/read-centric; Product clients do not create provider sales;
- BusinessOrderIntent and InvoicingIntent are Materialization owner reactions, not direct client commands;
- Party Resolution is a bounded human Materialization resolution; Destination write remains conditioned on D8 proof;
- Fulfillment exposes explicit physical checkpoints, provider-readiness meaning, Fulfillment Nodes and internal operating targets without becoming WMS/TMS or a status workflow;
- generic automation cannot establish physical facts solely from Permission; a machine physical fact requires an explicitly proven system Principal/source;
- Shipment remains source-qualified external observation;
- Post-Sale uses canonical scoped Resolution without copying provider Claim/Return/refund action vocabulary;
- Work owns responsibility/assignment/escalation but not source truth; generic `SubmitWorkResolution` is deferred;
- **zero-P baseline**: cross-owner Sale operational projection is deferred until D6 proves repeated composition need;
- no Product 1.0 bulk operation is admitted merely for symmetry.

### Consequential safety / concurrency

Every admitted C operation now declares a complete safety tuple in `D5-B2-OPERATION-ADMISSION-MATRIX.md`:

```text
consequence class
idempotency disposition
concurrency / precondition disposition
```

Ratified hardenings include:

- no silent safety cell;
- Party Resolution always uses client idempotency plus current candidate-set/resolution precondition;
- access-role and Authorization Delegation revocation is fail-safe/monotonic and not blocked merely by stale snapshots; later re-grant is explicit;
- business-resource deactivation remains current-state protected because stale deactivation can itself be unsafe;
- physical Fulfillment checkpoints are duplicate-sensitive and current-state protected;
- `EvaluatePriceScenario` is explicitly non-consequential/side-effect-free and therefore has no idempotency/concurrency machinery by ritual.

### Whole-Matrix outcome

```text
Parent D0→D4/D4-R1/D5-B1       CURRENT STRUCTURE CONFIRMED
D5-B2 operation inventory       B2-local RESTRUCTURE corrections APPLIED
Whole-Matrix Global Coherence   ACCEPTED / RATIFIED
Parent-stage reopen             NONE
```

## 4. What is prohibited now

While **D5-B2 Wire Contract is NEXT / NOT YET ACCEPTED**:

- do not begin D6–D9 target design or implementation;
- do not alter accepted D0–D4/D4-R1/D5-B1 or ratified B2 operation meaning by wire convenience;
- do not derive paths/schemas from legacy routes/current OpenAPI/provider endpoints;
- do not recreate Product/PIM, generic Integration/Mutation/Action/Operation/Workflow/Rules/AI authority, generic finance ledger, Task/Case engine or market collector platform;
- do not put initial price inside ListingIntent or collapse `listing.manage` with `price.manage` because a provider create payload is combined;
- do not move Availability/Fulfillment meaning into Offering because a provider request combines fields;
- do not expose direct Sankhya TOP/NUNOTA/order/invoice/retry choreography;
- do not expose provider OAuth/webhook protocol as Product business operations;
- do not create direct client commands for owner reactions already defined by D3;
- do not allow a generic machine token to fabricate physical facts;
- do not re-admit generic Work resolution without the concrete closure-path trigger;
- do not invent a P/BFF surface before D6 consumer evidence;
- do not weaken Organization scope, source-qualified identity, honest knowledge/freshness, Permission/Governance separation, idempotency, concurrency, ambiguity, recovery, multi-target scope or convergence laws;
- do not add compatibility/versioning, generic bulk or public event stream without a real entitled consumer;
- do not treat Git-history review artifacts or `AI-DIALOG.md` as current target authority.

## 5. Exact next action

**Open D5-B2 Wire Contract / Resource-Path-Schema Grammar from the ratified Operation Admission Matrix.**

The sub-batch must establish, before any implementation:

1. **Resource/path hierarchy**
   - Organization-owned Product operations remain under `/organizations/{organization_id}/...`;
   - `GetCurrentAccessContext` gets only the bounded platform-scoped self-only discovery shape;
   - source-qualified external identities remain explicit and no synthetic alias is introduced merely for prettier URLs.

2. **HTTP semantic grammar**
   - use ordinary resource methods for real resource state when honest;
   - use owner-specific methods when CRUD/PATCH-status would lie;
   - do not create generic Action/Command/Operation endpoints.

3. **Request / response families**
   - exact schema shapes for admitted operations;
   - ListingIntent authoring + media correlation;
   - PriceIntent target duality (`existing Listing | pre-creation ListingIntent context`) without embedding price in ListingIntent;
   - owner-local Intent/tracking outcomes without generic long-running Operation identity.

4. **Knowledge / outcome semantics**
   - known / known-empty / unknown / unavailable / partial and freshness/provenance where material;
   - accepted / rejected / pending / ambiguous separated from applied/converged where applicable;
   - no plausible defaults from missing evidence.

5. **Problems and access**
   - RFC 9457 Problem Details for malformed/schema/AuthN/ordinary access/idempotency/precondition/unsupported/server problems;
   - valid business rejection/approval-required/pending remains semantic result, not fake 403.

6. **Safety on the wire**
   - exact `Idempotency-Key` placement, same-key/same-request behavior and mismatched-request failure;
   - exact opaque MPC precondition/concurrency mechanism for every matrix safety row;
   - authority revocation remains fail-safe/monotonic;
   - no blind retry of ambiguous external acceptance.

7. **Collections**
   - pagination/filter/search/cursor only for admitted real collection consumers;
   - no bulk until a real member-level workflow proves it.

8. **Permissions / client classes**
   - exact Permission→wire-operation mapping;
   - human vs machine/system restrictions, including physical-fact establishment.

9. **Media seam**
   - wire contract for listing-context media intake/reference without choosing D7 blob/storage/hash/CDN mechanics;
   - arbitrary client external URL is not trusted authored media.

10. **Technical non-Product ingress**
    - provider OAuth start/callback, webhooks and future real external connector ingress remain explicitly separate from Product API and provider vocabulary stays boundary-local.

11. **Work closure-path audit**
    - for every Product 1.0 Work-producing condition class, prove source-owner automatic closure/reconciliation or an admitted owner-specific human capability;
    - only a proven human-evidence gap can re-admit a bounded Work→source evidence-submission operation.

12. **Machine-readable authority**
    - define OpenAPI operation naming/spelling and the path to one machine-readable Product API wire authority;
    - no legacy compatibility surface without a real entitled consumer.

Do not choose D6 UI/BFF composition, D7 transactions/queues/workers/blob storage/Keycloak deployment, D8 controlled-effect proofs or implementation.

If a wire shape cannot preserve the ratified owner/safety/identity/knowledge semantics without distortion, stop and reopen only the implicated B2/parent decision rather than weakening the contract.

Implementation remains blocked until D9.

## 6. Fresh-session success test

A fresh session must conclude unambiguously:

- D0→D4 + D4-R1 accepted/canonical;
- Decision Reconciliation accepted/canonical;
- D5-B1 accepted/canonical;
- D5-B2 OPEN / ACTIVE;
- B2-A accepted in-stage;
- Operation Admission Matrix Blocks 1–5 + Whole-Matrix Global Coherence are operator-ratified and consolidated;
- review candidate is absent from the active tree and historical review dialogue is Git history only;
- every admitted C has explicit consequence/idempotency/precondition disposition;
- current access discovery is self-only platform-scoped, while Organization business calls remain Organization-scoped;
- initial publication price remains PriceIntent, never ListingIntent content;
- generic Work resolution and cross-owner P remain deferred;
- **Wire Contract / Resource-Path-Schema Grammar is the exact next action**;
- implementation remains blocked until D9.

If not, the active authority tree is inconsistent.
