# Marketplace Central — Architecture Rebaseline

> **Role:** sole current-program status/router after `AGENTS.md`  
> **Current phase:** **D5 — API — OPEN / ACTIVE; D5-B1 ACCEPTED / CANONICAL; D5-B2 OPEN / ACTIVE**  
> **D5-B2 current state:** **B2-A + Operation Matrix Blocks 1–5 ACCEPTED IN-STAGE; Whole-Matrix Fable Round 1 COMPLETE; GPT adjudication COMPLETE; one focused price-intent contradiction remains; Fable Round 2 = NEXT**  
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

The current review candidate `D5-B2-WHOLE-MATRIX-REVIEW-CANDIDATE.md` and `AI-DIALOG.md` are **not authority** and are deliberately excluded from the authority path. They are review input only.

`D5-API.md` remains D5-B1 authority. Its old next-action wording is a pre-B2-opening snapshot. Never reconstruct target authority from memory, chat, Git history, retired ADRs, `AI-DIALOG.md`, review candidates or current code/OpenAPI shape.

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
       ├─ Matrix Block 1 — Identity/Access + Portfolio + Readiness — ACCEPTED IN-STAGE
       ├─ Matrix Block 2 — Offering + Price + Availability — ACCEPTED IN-STAGE
       ├─ Matrix Block 3 — Market Intelligence + Commercial Economics — ACCEPTED IN-STAGE
       ├─ Matrix Block 4 — Governance + Sales + Materialization — ACCEPTED IN-STAGE
       ├─ Matrix Block 5 — Fulfillment + Post-Sale + Work + P compositions — ACCEPTED IN-STAGE
       └─ Whole-Matrix Global Coherence
            ├─ lead review — RESTRUCTURE NOW / B2-local corrections identified
            ├─ review candidate — PREPARED / NON-AUTHORITATIVE
            ├─ Fable Round 1 — COMPLETE / REVISE B2-LOCAL
            ├─ GPT adjudication — COMPLETE
            └─ focused Fable Round 2 — NEXT / initial-publication price only
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

## 3. Accepted B2 baseline before whole-matrix consolidation

### B2-A — Client/Auth

- Product API authentication uses one standards-based OIDC/OAuth boundary.
- Humans use Authorization Code + PKCE semantics; confidential machine clients use Client Credentials/service-account semantics.
- MPC remains authority for Principal, Organization Membership, AccessRole/Permission/RoleAssignment and all business decisions.
- Tokens are audience-bound to MPC API; no global/static MPC Product API key or IdP-role business authority.
- Keycloak is the first implementation/proof candidate; D7 owns provider/deployment/realm/secrets/token-lifetime realization.

### Matrix Blocks 1–5

- **Block 1:** minimal D2 access context/role assignment, Portfolio Installation lifecycle/configuration, and marketplace-context Readiness Product discovery/requirements/correspondence; no PIM/IAM/integration platform.
- **Block 2:** Listing actual state is Offering Q; `ListingIntent` is create/edit authoring/tracking; `PriceIntent` is separate; Availability owns Sellable Availability; no giant Listing CRUD, direct price/stock set or generic async Operation.
- **Block 3:** Market Intelligence exposes competitive interpretation; Economics owns stateless scenario evaluation plus durable material L0/L1/L2 lineage; no Recommendation/Simulation authority, generic ledger or price actuation in Economics.
- **Block 4:** Governance decisions/delegations remain authorization-only; Sales is externally originated/read-centric; Materialization creates BusinessOrder/Invoicing intents from accepted owner reactions; no direct Sankhya/order/invoice/retry/workflow API.
- **Block 5:** Fulfillment exposes physical checkpoints/nodes/artifacts, Shipment remains external read observation, Post-Sale uses canonical scoped Resolution, Work owns responsibility/lifecycle without source truth, and the provisional Sale operational P remains subject to whole-matrix adjudication.

These block decisions remain accepted in-stage until the whole-matrix package is operator-ratified and canonically consolidated.

## 4. Whole-matrix review state

`docs/engineering/rebaseline/D5-B2-WHOLE-MATRIX-REVIEW-CANDIDATE.md` remains non-authoritative review input.

Fable Round 1 returned **REVISE — B2-local only** and found no parent-stage reopen. GPT independently adjudicated the findings in `AI-DIALOG.md`.

### Converged direction from Round 1 + GPT adjudication

Pending final operator ratification after Round 2, the following are converged review conclusions:

1. add ListingIntent-scoped authored-media intake; never Product/media-master authority;
2. add Fulfillment-owned internal operating-target Q/C with effective-value provenance; never generic SLA/rules authority;
3. defer generic `SubmitWorkResolution`; wire closure must audit every Work-producing condition for a legitimate source-owner closure path;
4. defer `GetSaleOperationalView` until D6 proves repeated P-composition need;
5. every admitted C operation must declare consequence class, idempotency disposition and concurrency/precondition disposition; silence is non-conformant;
6. `ResolveBusinessSystemPartyResolution` requires mandatory client idempotency plus current resolution/candidate-set precondition by default;
7. `GetCurrentAccessContext` is platform-scoped **self-only** discovery for the authenticated Principal; Organization-owned business routes remain Organization-path-scoped;
8. authority revocation is fail-safe/monotonic: Authorization Delegation and AccessRole revocation are structurally idempotent and are not blocked merely by stale snapshots; re-grant is a new explicit authority action;
9. B2-A OIDC/OAuth + MPC-owned Principal/Membership/Permission remains the selected Global Maximum;
10. authored-media URL trust is fail-closed and media reference/selection reads remain on ListingIntent semantics; binary mechanics remain D7;
11. Fulfillment target configuration preserves default/override provenance and does not create a generic target/SLA platform.

None of these items is canonical merely because reviewer/lead currently converge; final consolidation waits for the focused contradiction below and operator ratification.

### One material contradiction remains

Fable F-WM-10 proposed putting **creation-time price inside ListingIntent**, with `PriceIntent` only for later price changes.

GPT rejected that proposal because current accepted authority already gives Listing and Price Intents distinct material identities, D4-R1 explicitly rejects a Publication aggregate absorbing Price, and B2 intentionally separates `listing.manage` from `price.manage`.

GPT's proposed corrected rule is:

> **Price is never ListingIntent-owned content, including initial publication. Initial active publication uses a correlated Offering-owned PriceIntent for the to-be-created listing context; D4/D7 may jointly serialize ListingIntent + PriceIntent + Availability-issued meaning into one provider create request without merging their identities/permissions.**

This is the only material review disagreement that survives Round 1 adjudication.

## 5. What is prohibited now

While focused Round 2 is open:

- do not begin resource/path/schema/OpenAPI crystallization;
- do not consolidate the candidate into the active matrix yet;
- do not begin D6–D9 design or implementation;
- do not treat Fable/GPT review dialogue or the candidate as authority;
- do not reopen already converged Round 1 findings unless the price contradiction logically invalidates one;
- do not mutate accepted D0–D4/D4-R1/D5-B1 semantics by review convenience;
- do not derive operations from legacy routes/current OpenAPI/provider endpoint shape;
- do not weaken Organization scope, source-qualified identity, Permission/Governance separation, idempotency, concurrency, ambiguity, recovery or convergence laws.

## 6. Exact next action

**Run one focused Fable Round 2 on the single surviving material contradiction: initial-publication price semantics.**

Follow the active `AI-DIALOG.md` handoff. Fable must independently decide the smallest Global Maximum between:

- F-WM-10: creation-time price becomes ListingIntent content and PriceIntent begins only after Listing creation; or
- GPT A10: creation-time price remains a separate Offering-owned PriceIntent correlated with the ListingIntent/to-be-created Listing context, preserving `listing.manage != price.manage` and joint physical provider serialization without ownership merge.

The challenge must confront D2 PriceIntent identity, D4-R1 anti-absorption language, B2 least-privilege Permissions and the actual provider requirement that initial creation physically contains a price.

Fable modifies only `AI-DIALOG.md`, appends the Round 2 material resolution and hands back to GPT. **Do not re-review A1–A9/A11 unless the price conclusion logically requires it.**

After Round 2, GPT adjudicates the last contradiction. If no contradiction survives, the converged Whole-Matrix package goes to the operator for ratification. Only after ratification do we consolidate the active matrix, remove the disposable candidate/reset review dialogue as appropriate, and open **D5-B2 Wire Contract / Resource-Path-Schema Grammar**.

Implementation remains blocked until D9.

## 7. Fresh-session success test

A fresh session must conclude unambiguously:

- D0→D4 + D4-R1 accepted/canonical;
- Decision Reconciliation accepted/canonical;
- D5-B1 accepted/canonical;
- D5-B2 OPEN / ACTIVE;
- B2-A and Matrix Blocks 1–5 accepted in-stage;
- Fable Whole-Matrix Round 1 and GPT adjudication are complete review input, not authority;
- no parent-stage reopen was found;
- exactly one material contradiction remains: initial-publication price as ListingIntent content vs correlated PriceIntent;
- focused Fable Round 2 is the exact next action;
- wire-contract design remains blocked until review convergence + operator ratification;
- implementation remains blocked until D9.

If not, the active authority tree is inconsistent.