# AuthorizationRequest D5 Wire Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: use TDD and execute this plan task-by-task. This repository forbids `docs/superpowers/`; this approved local convention uses `docs/plans/` for transient execution plans.

**Goal:** Turn the operator-ratified D5-R5 AuthorizationRequest Product surface into one executable canonical OpenAPI wire and prove the resulting 106-operation / 31-Permission surface.

**Architecture:** Add only two purpose-bounded Governance reads. Reanchor the existing human decision capability on the canonical AuthorizationRequest and preserve request-local concurrency separately from material target validity. Feed the new identity into F13 and Work without adding public workflow operations.

**Tech Stack:** OpenAPI 3.1.2 YAML, Node.js verifier, Redocly 2.45.0, existing repository full gate/generated TypeScript/Go projection proofs.

**Spec:** `docs/engineering/rebaseline/D6-R2-NOTIF-01-D5-R5-AUTHORIZATION-REQUEST-PRODUCT-SURFACE.md`

## Global Constraints

- Product implementation remains BLOCKED UNTIL accepted D9.
- Canonical wire authority is `contracts/api/product/openapi.yaml` and its referenced Product path/component files.
- D1-R2/D2-R6/D3-R3 ownership and lifecycle are binding.
- Exactly two new Product operations; no new ordinary Permission.
- No public AuthorizationRequest create/invalidate/reauthorize API.
- No target ETag as authorization-validity oracle.
- Fable review remains mandatory only after D5 proof + final P9.

---

### Task 1: RED — executable D5-R6 verifier

**Files:**
- Create: `scripts/verify-authorization-request-oad.mjs`
- Modify: `package.json`

**Produces:** A verifier that bundles the canonical OAD and rejects the pre-repair 104/31 contract.

- [ ] Add checks for exact operations `ListMyActionableAuthorizationRequests`, `GetMyActionableAuthorizationRequest`, and reanchored `CreateAuthorizationDecision`.
- [ ] Require 106 operation IDs / 31 ordinary Permissions.
- [ ] Require both actionable reads to be H-only, `governance.decide`, ControlledActionGovernance.
- [ ] Require list query controls exactly `limit,cursor`.
- [ ] Require detail ETag and a closed four-variant target/review-basis pairing.
- [ ] Require decision path to use `authorization_request_id`, `If-Match`, `Idempotency-Key`, and body-only `outcome`.
- [ ] Require F13 `AuthorizationRequestRef`, F14 target-oriented continuation, and Work `authorization_request` origin.
- [ ] Add negative controls for permission leakage, machine actionability, target-ETag regression, generic payload, stale-write removal, F13 target regression and public request-create leakage.
- [ ] Put this verifier before the legacy Notification verifier in both npm gate lanes.
- [ ] Run full gate and observe RED because the OAD is still 104/31 and the new operations are absent.

### Task 2: GREEN — canonical Governance wire

**Files:**
- Modify: `contracts/api/product/openapi.yaml`
- Modify: `contracts/api/product/paths-economics-governance-sales-materialization.yaml`
- Modify: `contracts/api/product/components.yaml`

**Produces:** Two actionable reads plus request-anchored decision command.

- [ ] Add `AuthorizationRequestId` path parameter and canonical target/request reference schemas.
- [ ] Add `ListMyActionableAuthorizationRequests` GET collection with only limit/cursor.
- [ ] Add `GetMyActionableAuthorizationRequest` GET detail with request-local ETag.
- [ ] Move `CreateAuthorizationDecision` from collection POST to `authorization-requests/{id}:decide`.
- [ ] Require both `If-Match` and `Idempotency-Key`; body is exactly `{ outcome }`.
- [ ] Add 412/428 to decision command and preserve 409/422 ambiguity/idempotency/state failures.

### Task 3: GREEN — typed review basis and historical Decision

**Files:**
- Modify: `contracts/api/product/components.yaml`

**Produces:** Structurally paired actionable request detail and truthful immutable Governance history.

- [ ] Add ListingIntent review basis: intent ID, source Product, target, desired, optional requirements revision.
- [ ] Add PriceIntent review basis: intent ID, target, desired price, optional bounded current-price observation and economic-policy evidence.
- [ ] Add BusinessOrderIntent review basis: intent ID, Sale snapshot, target SourceInstance, PartyResolution and DestinationRealization.
- [ ] Add InvoicingIntent review basis: intent ID, Sale snapshot, BusinessOrderIntent snapshot, optional FulfillmentExecution correlation.
- [ ] Add four actionable request variants that pair each target ref with its matching review basis.
- [ ] Add `authorization_request_id` and immutable review basis to `AuthorizationDecision`; remove target ETag from decision semantics.

### Task 4: GREEN — Notification and Work feed-forward

**Files:**
- Modify: `contracts/api/product/paths-notifications.yaml`
- Modify: `contracts/api/product/paths-fulfillment-postsale-work.yaml`

**Produces:** F13 exact request continuation and zero-decider Work origin without new operations.

- [ ] Change F13 source to canonical `AuthorizationRequestRef`.
- [ ] Keep F14 source as canonical no-ETag `AuthorizationTargetRef`.
- [ ] Add `authorization_request` to Work origin union and ListWork `origin_kind` filter.
- [ ] Do not add Work commands or Notification kinds.

### Task 5: GREEN — reconcile existing proofs

**Files:**
- Modify only proof scripts whose current-state assertions are intentionally superseded by D5-R6.

**Produces:** Historical non-regression still proves historical generations while current proofs recognize 106/31.

- [ ] Run the unchanged new verifier first; fix OAD, not verifier, for contract failures.
- [ ] Run full gate and inspect every legacy verifier failure.
- [ ] Preserve historical 95/29 and 99/30 fixtures rather than rewriting history to 106.
- [ ] Update only current-generation Notification/current-OAD census expectations from 104 to 106 and F13 source semantics where required.
- [ ] Re-run until the repository full gate, Redocly and generated projections pass.

### Task 6: D5-R6 durable proof and routing

**Files:**
- Create: `docs/engineering/rebaseline/D6-R2-NOTIF-01-D5-R6-AUTHORIZATION-REQUEST-OAD-WIRE-PROOF.md`
- Modify: `docs/roadmap.md`
- PR metadata only after final proof.

**Produces:** Canonical D5 wire proof and next gate P9.

- [ ] Record exact census, wire laws, negative controls and final CI evidence.
- [ ] Route roadmap to D6 P9 final Screen Contracts; keep B10/D7-R/D8-R blocked.
- [ ] Keep independent Fable review explicitly after final P9 and before Global-Maximum closure/D7-R.
- [ ] Verify exact final HEAD with full gate before any completion claim.
