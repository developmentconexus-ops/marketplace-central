# D5-B2 — Whole-Matrix Global Coherence Review Candidate

> **Status:** NON-AUTHORITATIVE REVIEW CANDIDATE  
> **Program:** Architecture Rebaseline / Technical System Design  
> **Review target:** D5-B2 Client/Auth + Operation Admission Matrix Blocks 1–5 as one coherent Product API surface  
> **Method:** DevelopmentConexus Engineering Method v1.0.0  
> **Prepared:** 2026-08-18  
> **Operator direction approved for independent challenge:** 2026-08-18

## 1. Review role

This file is a bounded adversarial-review package only. It is **not architecture authority** and does not replace the router, D5-B1, the active B2 artifacts, accepted D0–D4/D4-R1, or the Decision Reconciliation Baseline.

Fable must reconstruct repository authority independently from `AGENTS.md` and the router before evaluating this candidate. Current code/OpenAPI/routes are evidence only.

The candidate records the lead review of B2-A + Blocks 1–5 and the operator-approved corrections that should be challenged before canonical consolidation.

## 2. Lead whole-matrix disposition

Lead outcome:

> **RESTRUCTURE NOW — B2-local only. Parent D0/D1/D2/D3/D4/D4-R1/D5-B1 structure remains confirmed.**

The matrix is broadly coherent but the whole-system challenge found four material local corrections plus three hardenings before wire-contract crystallization.

No parent-stage reopen is currently justified.

## 3. Proposed correction F-B2-G1 — ListingIntent-scoped authored-media intake

### Finding

D4-R1 explicitly admits MPC-authored/uploaded media for listing context while rejecting a Product-media master. The current Block 2 matrix lets a client create/update a ListingIntent and select media but does not yet admit the Product operation by which listing-context authored media enters MPC.

### Proposed correction

Admit the semantic capability provisionally named `CreateListingIntentMedia`:

- consumer: human or automation authoring a ListingIntent;
- owner: Marketplace Offering Operations;
- class: C;
- Organization + exact ListingIntent draft scope;
- ordinary Permission: `listing.manage`;
- client idempotency key: mandatory;
- association/selection/order remains controlled by ListingIntent draft revision/concurrency;
- media is listing-context state only, never Product/PIM/media-master authority.

D7 still owns upload/blob/storage/hash/resizing/CDN/presigned-URL mechanics.

## 4. Proposed correction F-B2-G2 — Fulfillment internal operating-target configuration

### Finding

D0 requires organization-operable MPC-owned internal operating-time targets without code edits. Block 5 currently exposes Fulfillment state/checkpoints/nodes but no Product operation for Fulfillment-owned target configuration.

### Proposed correction

Admit bounded Fulfillment operating-target Q/C semantics, provisionally:

- `GetFulfillmentOperatingTargets` — Q, both, `fulfillment.read`;
- `UpdateFulfillmentOperatingTargets` — C/resource update, human baseline, `fulfillment.manage`;
- current-state concurrency/precondition required where lost update is material;
- desired-state update may be structurally idempotent;
- provider deadline/window remains external authoritative evidence and is never rewritten by the MPC internal target.

Do not create an SLA/rules/timer/cron DSL. Scheduling/escalation mechanics remain D7.

## 5. Proposed correction F-B2-G3 — generic `SubmitWorkResolution` over-admission

### Finding

The current Block 5 matrix admits `SubmitWorkResolution`. Across the whole matrix this creates duplicate-authority pressure because concrete source-owner resolution operations already exist, for example:

- `ResolveEconomicAttribution`;
- `ResolveBusinessSystemPartyResolution`;
- `ResolveSaleSellingEntityAttribution`;
- `ResolveProductChannelCorrespondence`.

A generic Work resolution payload capable of deciding those meanings would turn Work into a command bus and violate D1/D3 ownership.

### Proposed correction

Change generic `SubmitWorkResolution` from ADMIT to **DEFER**.

Baseline resolution path:

```text
Work identifies actionable source-owned condition
  -> client invokes the source owner's specific resolution capability
  -> source owner commits semantic result
  -> Work reconciles/updates from source-owned result
```

A bounded Work-to-source evidence-submission capability may be admitted later only if a concrete case proves that Work genuinely holds resolution evidence, no existing owner-specific operation is appropriate, and a generic command bus is still avoided.

## 6. Proposed correction F-B2-G4 — premature cross-owner P

### Finding

The current Block 5 matrix admits `GetSaleOperationalView` as the only baseline P composition. It is semantically legal, but before D6 there is not enough consumer evidence to justify the extra cross-owner permission, property-security, partiality/freshness, caching and evolution surface.

### Proposed correction

Change `GetSaleOperationalView` from ADMIT to **DEFER — D6 consumer proof**.

D5-B1 continues to permit legitimate read-only P compositions. D6 may reopen this operation if a repeated real client need proves that composing Sales + Materialization + Fulfillment + Post-Sale + Work in every consumer is materially worse than one bounded projection.

No P write/concurrency/authorization/retry authority is admitted.

## 7. Hardening H1 — Party Resolution idempotency

`ResolveBusinessSystemPartyResolution` may lead to consequential native party creation when zero safe matches exist and owner rules allow creation. Therefore the baseline should not depend on a future structural-idempotency exemption.

Proposed hardening:

> **`ResolveBusinessSystemPartyResolution` requires a client idempotency key by default, plus current-state precondition/concurrency where material.**

Intake idempotency never authorizes blind replay of an ambiguous native effect.

## 8. Hardening H2 — Current Access Context scope exception

`GetCurrentAccessContext` must discover the authenticated Principal's accessible Organization memberships, so it cannot itself require a previously selected Organization path.

Proposed clarification:

> **`GetCurrentAccessContext` is a bounded platform-scoped D2 discovery Q. Every Organization-owned business Product API operation remains explicitly `/organizations/{organization_id}/...` scoped.**

No ambient/default/current Organization header/session authority is introduced.

## 9. Hardening H3 — Authorization Delegation stale-state protection

Bounded Authorization Delegation/Grant administration was admitted without forcing a premature universal Grant identity model.

Proposed hardening:

- delegation update/revoke requires current-state concurrency/precondition where stale overwrite or stale revocation would materially change authority;
- ordinary Permission to administer Governance remains distinct from the Grant/Delegation semantics themselves;
- a later physical Grant identity/cardinality choice must preserve historical decision authority context and revocation without history rewrite.

## 10. Lead whole-system checks that currently PASS

Subject to the corrections above, the lead review found:

- no duplicate or missing D1 business authority;
- no new D1 semantic edge required;
- Product 1.0 lifecycle remains reachable through legitimate Product operations or accepted owner-triggered reactions;
- human vs machine/system client classes remain coherent, including physical-evidence restrictions;
- Permission catalog remains capability-oriented rather than one-per-endpoint and preserves useful least privilege;
- Q/C classifications remain honest; zero baseline P after G4 is acceptable;
- idempotency and concurrency remain distinct failure controls;
- BusinessOrderIntent/InvoicingIntent/AvailabilityIntent normal paths remain owner-triggered rather than client-commanded;
- Organization and source-qualified identity laws remain intact;
- provider richness remains bounded and source-qualified;
- no baseline generic bulk endpoint is currently justified;
- no Product/PIM, generic Integration, Mutation/Action/Operation, Workflow, Rules, Finance ledger, Task/Case or provider graph is required;
- Structural Inversion against legacy routes/OpenAPI passes;
- second-marketplace/business-system/automation/fulfillment futures have seams without implementation of speculative capabilities.

## 11. Explicit review questions for Fable

Fable should attack this package rather than validate it ceremonially. In particular:

1. Is any admitted operation still a disguised second authority or implementation/runtime command?
2. Is any Product 1.0 capability still unreachable after G1/G2?
3. Does deferring `SubmitWorkResolution` create a real unresolved Work workflow gap, or correctly preserve source-owner authority?
4. Does deferring the P composition create unacceptable client coupling before D6, or correctly apply YAGNI?
5. Are any client-class assignments unsafe, especially automation authoring, Governance decisions and physical Fulfillment facts?
6. Is the Permission floor too coarse or too fragmented for real blast radius?
7. Is any claimed structural-idempotency exemption unsafe? Challenge create/update/resolve/checkpoint operations individually where material.
8. Are concurrency/precondition requirements missing on any standing human decision, configuration, delegation, Work or draft lifecycle?
9. Does B2-A authentication remain the Global Maximum: external OIDC/OAuth, Keycloak as first proof candidate, MPC-owned Principal/Membership/Permission/business authority, no global MPC API key?
10. Is ListingIntent-scoped media intake the smallest correct seam, or does it accidentally create hidden asset/PIM authority?
11. Are Fulfillment operating targets correctly Fulfillment-owned, or does another accepted owner own part of that meaning?
12. Do any deferred operations need to be admitted now because future retrofit would otherwise be disproportionately costly?
13. Does the whole matrix remain correct if legacy OpenAPI/routes/controllers are inverted or deleted?
14. Can a materially better Global Maximum reduce authorities/operations further without losing Product 1.0 reachability or future seams?

Allowed review disposition: `APPROVE`, `REVISE`, or `REJECT`, with material findings only. If a finding requires a parent-stage reopen, identify the exact parent invariant/evidence rather than broadening scope by preference.

## 12. Review protocol / next step

Follow the canonical **Standard Fable review workflow** in `developmentconexus-ops/conexus-methodology/README.md`.

Fable review is non-authoritative evidence. It should be written to the active `AI-DIALOG.md` cycle only; no other repository file should be modified unless the operator separately authorizes that scope.

After Fable:

1. GPT independently adjudicates each material finding;
2. round 2 occurs only if a material contradiction survives;
3. the candidate is revised/converged as needed;
4. operator ratifies the converged package;
5. only then are corrections consolidated into the active B2 matrix and this candidate deleted;
6. the next B2 sub-batch becomes Wire Contract / Resource-Path-Schema Grammar.

Implementation remains blocked until D9.