# Marketplace Central — Architecture Rebaseline

> **Role:** sole current-program status/router after `AGENTS.md`  
> **Current phase:** **D5 — API — OPEN / ACTIVE; D5-B1 ACCEPTED / CANONICAL; D5-B2 OPEN / ACTIVE**  
> **D5-B2 current state:** **B2-A + Operation Matrix Blocks 1–5 ACCEPTED IN-STAGE; Whole-Matrix Fable Rounds 1–2 + GPT adjudication CONVERGED; operator ratification = NEXT**  
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

The current review candidate `D5-B2-WHOLE-MATRIX-REVIEW-CANDIDATE.md` and `AI-DIALOG.md` are **not authority** and are deliberately excluded from the authority path. They remain review input until operator ratification and canonical consolidation.

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
            ├─ lead review — RESTRUCTURE NOW / B2-local corrections
            ├─ Fable Round 1 — COMPLETE / REVISE B2-LOCAL
            ├─ GPT Round 1 adjudication — COMPLETE
            ├─ focused Fable Round 2 — COMPLETE / CONVERGED ON PRICEINTENT
            ├─ GPT final adjudication — COMPLETE / CONVERGED
            └─ operator ratification — NEXT
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
- **Block 5:** Fulfillment exposes physical checkpoints/nodes/artifacts; Shipment remains external observation; Post-Sale uses canonical scoped Resolution; Work owns responsibility/lifecycle without source truth. The provisional generic Work-resolution command and cross-owner Sale P are now converged for deferral pending operator ratification.

These block decisions remain accepted in-stage until the converged whole-matrix package is operator-ratified and canonically consolidated.

## 4. Converged Whole-Matrix review package — pending operator ratification

`docs/engineering/rebaseline/D5-B2-WHOLE-MATRIX-REVIEW-CANDIDATE.md` and `AI-DIALOG.md` remain non-authoritative review input. Fable Rounds 1–2 and GPT adjudication now have **zero surviving material contradictions** and found **no parent-stage reopen**.

The converged package is:

1. **ADD ListingIntent-scoped authored-media intake** under Offering; no Product/media-master authority. Client-supplied arbitrary external URLs are not trusted authored media; metadata/reference/selection remains ListingIntent semantics while binary/storage mechanics remain D7.
2. **ADD Fulfillment-owned internal operating-target Q/C**, exposing effective value plus material default/override provenance; provider deadline remains external authoritative evidence; no generic SLA/target platform.
3. **DEFER generic `SubmitWorkResolution`**. Before wire closure, every Product 1.0 Work-producing condition class must prove a legitimate source-owner auto-resolution or owner-specific human closure path; only a real missing human-evidence path can re-admit a bounded Work evidence-submission capability.
4. **DEFER `GetSaleOperationalView` / baseline P** until D6 proves repeated cross-owner composition pain. Zero-P baseline is valid.
5. **Every admitted C operation declares a complete safety tuple**: consequence class; idempotency (`mandatory client key | named structural anchor/exemption | N/A with reason`); concurrency/precondition (`required | not material with named reason`). Silence is non-conformant. Consequential C keeps D5-B1's fail-closed key default unless a reviewable structural anchor proves duplicate intake unreachable/harmless.
6. **`ResolveBusinessSystemPartyResolution` requires mandatory client idempotency + current resolution/candidate-set precondition by default.** Staleness is material by construction; dedupe never authorizes blind replay of an ambiguous native effect.
7. **`GetCurrentAccessContext` is platform-scoped self-only D2 discovery**: no Principal parameter; only the authenticated Principal's current memberships are returned. Organization-owned business operations remain explicit Organization-path scoped.
8. **Authority revocation is fail-safe/monotonic.** Authorization Delegation and AccessRole revocation are structurally idempotent and not blocked merely by stale snapshots; successful revocation cannot be implicitly resurrected by a concurrent stale update. Re-grant is a new explicit authority action. Business-resource deactivation does not inherit this exception.
9. **B2-A remains the selected Global Maximum**: OIDC/OAuth external identity boundary, MPC-owned Principal/Membership/Permission, audience-bound tokens, no static MPC Product API keys or IdP-owned business permissions.
10. **Initial and later price meaning always uses `PriceIntent`.** Price is never ListingIntent-owned content. `PriceIntent` may target an existing source-qualified Listing or a pre-creation ListingIntent context; this is target duality on the existing PriceIntent identity, not a new identity class.
11. **Pre-creation price revision is explicit supersession, not mutable PriceDraft.** A newer PriceIntent may explicitly supersede the current pending PriceIntent for the same pre-creation target while preserving attribution/history; automation cannot silently supersede a standing human-authored intent. `withdraw/cancel` remains DEFER.
12. **Active publication dispatch fails closed on required separate inputs.** Before the provider effect, current valid/correlated ListingIntent representation + PriceIntent price + Availability-issued meaning must be established as required by the selected lane. Submission may wait pending required inputs; external dispatch never guesses/defaults them. Each owner/intent evaluates convergence independently.
13. **ListingIntent read shape never embeds price as ListingIntent-owned value.** It may expose a typed PriceIntent correlation/reference; price meaning/history remains on the PriceIntent surface, preserving `listing.manage != price.manage`.
14. **Minor bounded decisions:** Marketplace Installation reactivation remains DEFER until a real workflow proves need; `fulfillment.execute` gating artifact reads is intentional least privilege because artifacts may contain operational/PII-sensitive material.

The detailed per-operation safety sweep is part of consolidation/wire preparation; at minimum it includes the operations named in the Round 1 GPT adjudication and every other admitted C row so no safety declaration remains silent.

## 5. What is prohibited now

Until operator ratification:

- do not consolidate the review candidate into the active matrix;
- do not delete/reset review artifacts yet;
- do not begin resource/path/schema/OpenAPI crystallization;
- do not begin D6–D9 design or implementation;
- do not treat the candidate or `AI-DIALOG.md` as authority;
- do not reopen D0–D4/D4-R1/D5-B1 absent new material evidence;
- do not derive operations from legacy routes/current OpenAPI/provider endpoint shape;
- do not weaken Organization scope, source-qualified identity, Permission/Governance separation, idempotency, concurrency, ambiguity, recovery or convergence laws.

## 6. Exact next action

**Operator ratification of the converged D5-B2 Whole-Matrix package.**

Allowed operator outcomes:

- `Aprovado` — ratify the converged package. GPT then consolidates the accepted corrections/hardenings into the active B2 artifacts, performs the complete admitted-C safety sweep, deletes the disposable Whole-Matrix review candidate, resets `AI-DIALOG.md` to its reusable protocol header, updates/revalidates the router/diff/HEAD, and opens **D5-B2 Wire Contract / Resource-Path-Schema Grammar**.
- reject/amend — keep the package non-authoritative and return only the implicated B2 finding for further adjudication.

There is **no Round 3** unless the operator introduces new material evidence/contradiction.

Implementation remains blocked until D9.

## 7. Fresh-session success test

A fresh session must conclude unambiguously:

- D0→D4 + D4-R1 accepted/canonical;
- Decision Reconciliation accepted/canonical;
- D5-B1 accepted/canonical;
- D5-B2 OPEN / ACTIVE;
- B2-A and Matrix Blocks 1–5 accepted in-stage;
- Whole-Matrix Fable Rounds 1–2 and GPT adjudication are complete review input, not authority;
- the Whole-Matrix package has zero surviving material contradictions and no parent-stage reopen;
- initial publication price remains a separate correlated `PriceIntent`, never ListingIntent content;
- operator ratification is the exact next action;
- Wire Contract remains blocked until ratification + canonical consolidation;
- implementation remains blocked until D9.

If not, the active authority tree is inconsistent.