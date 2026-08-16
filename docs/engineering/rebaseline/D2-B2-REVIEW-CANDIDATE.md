# D2 Batch B2 — Independent Review Candidate

> **STATUS: REVIEW CANDIDATE — NOT ARCHITECTURE AUTHORITY**  
> **Stage:** D2 — Identity / Tenant / Data Ownership  
> **Operator posture:** B2 direction approved; independent challenge required before canonical filing/closure  
> **Disposable:** delete after adjudication; durable meaning belongs only in the canonical D2 artifact/router  
> **Date:** 2026-08-16

## Reviewer bootstrap

Read the current authority path independently before reviewing this candidate:

1. `AGENTS.md`
2. `docs/engineering/rebaseline/README.md`
3. `docs/engineering/standards/root-cause-global-maximum-method.md`
4. `ARCHITECTURE.md`
5. `docs/architecture/decisions/README.md`
6. `docs/engineering/rebaseline/D0-PRODUCT-SYSTEM-DEFINITION.md`
7. `docs/engineering/rebaseline/D1-DOMAINS-BOUNDARIES.md`
8. `docs/engineering/rebaseline/D2-IDENTITY-TENANT-DATA-OWNERSHIP.md`
9. `docs/engineering/rebaseline/EVIDENCE-REGISTER.md`
10. only the specific legacy ADRs/current code needed as evidence.

Do **not** treat this file, chat history, or `AI-DIALOG.md` as architecture authority.

## Context already authoritative

D2 is OPEN / IN PROGRESS. The canonical D2 artifact already owns the accepted identity/tenant/data-ownership baseline, including:

- identity follows semantic authority;
- Organization is the tenant/isolation root;
- MPC-owned Organization, Marketplace Installation, Selling Entity, Inventory Source, Fulfillment Node and Principal;
- source-qualified external Product, Listing/Variation, Sale/Order, Shipment and native financial identities;
- minimal Source Instance identity; D4 owns source protocols/capabilities/credentials;
- MPC-owned intents, Post-Sale Resolution, Economic Attribution/Reconciliation, Work and Authorization Decision semantics;
- OIDC AuthN with MPC-owned Principal/Membership/Role/Permission ordinary-access state;
- one D1 write authority per canonical business meaning;
- exact Money/rates, bounded `Fact<T>`, provenance/time semantics;
- **clean target database with no legacy-data inventory, archive, migration or compatibility requirement**, per final operator ruling.

This candidate must not reopen that baseline without material contrary evidence.

---

## B2 candidate set

### B2.1 — MPC-owned IDs are opaque, stable and non-reusable

For every MPC-owned canonical identity, the identifier carries no business meaning and does not embed provider/source identifiers, Organization, state, time, CNPJ/CODPROD or other mutable semantics.

Identity survives changes to mutable attributes/correspondences. Once an identity has denoted one canonical entity/occurrence, it is never recycled to denote another.

D2 does **not** select UUIDv4/UUIDv7/ULID/database encoding unless a material later requirement proves that choice architectural. Physical encoding remains implementation/runtime detail as long as the invariant holds.

### B2.2 — External identity binding and ordinary access are explicit durable relations

OIDC external identity binding `(issuer, subject)` maps to at most one MPC Principal. Email, username or other mutable claims never auto-merge Principals.

A future IdP replacement may bind another `(issuer, subject)` to an existing Principal only through explicit identity-administration action; historical Principal attribution is not rewritten.

Organization Membership and RoleAssignment are explicit revocable MPC state. Revocation changes future access but never rewrites historical Work, Authorization Decisions or other actor attribution.

### B2.3 — Source Instance is a namespace identity, not integration configuration

`SourceInstance` is the MPC-owned reference identity for one logical namespace of an externally authoritative business-system/source when Marketplace Installation is not already the correct qualifier.

Changing credentials/token/connection mechanics does not itself create a new Source Instance. Pointing at a materially different external authoritative namespace/environment does.

D4 owns how concrete source identity, connection/configuration and capabilities are verified. Marketplace integrations use Marketplace Installation where it already provides the correct namespace; do not wrap every Installation in a generic SourceInstance merely for uniformity.

### B2.4 — Cross-authority references are typed; no universal entity graph

Known semantic relationships use explicit typed references/contracts owned by their domains, e.g. BusinessOrderIntent→Sale reference.

The target does not introduce a universal `{entity_type, entity_id}` reference system, universal entity table or generic relationship graph.

A domain may own a local polymorphic subject only where its semantics genuinely require it. `Operational Work` is the principal accepted example: its subject may be object, relation, absence, population/coverage scope, etc. That does not generalize into a platform entity graph.

### B2.5 — Material quantities preserve exactness

Quantities used in canonical persisted state or consequential decisions must not silently lose precision through binary floating-point approximation when that could change business correctness.

Domains/sources determine whether a quantity is integer or fractional and what unit/precision is meaningful. D2 does not create a generic Unit-of-Measure/conversion framework.

### B2.6 — Canonical identity lifetime is historical-safe; no universal soft-delete framework

Canonical identities already referenced by material history are not silently deleted/reassigned such that historical meaning changes.

Deactivation/terminal lifecycle is owned by the applicable domain and represented only where semantically needed. D2 does **not** require `deleted_at`, `is_active`, soft-delete or status columns on every table.

Historical occurrences such as Authorization Decisions are not rewritten into a different past decision; later decisions/revocations create later state/history while preserving prior occurrence meaning.

### B2.7 — Role/Permission definitions are product-defined for Product 1.0

Ordinary Roles and Permissions are product-defined definitions/bundles. Organization-specific mutable state is Membership/RoleAssignment.

Product 1.0 does not require per-Organization custom role design, nested groups, generic policy/ACL/ReBAC administration or an IAM designer. Exact permission/operation catalog may be completed in D5 without changing D2 ownership.

### B2.8 — No universal Evidence/Observation identity framework

D1 already rejects Evidence/Provenance/Audit as independent business domains. D2 therefore does not introduce universal `EvidenceID`, `ObservationID`, evidence graph or generic observation business authority.

Domains persist domain-specific observations/evidence where needed (e.g. Market Intelligence observations, Governance decision context, Commercial Economics evidence). Shared provenance/knowledge/time primitives may support them without acquiring authority.

### B2.9 — D2 disposition of legacy ADRs reopened into D2

These old implementation-era ADR shapes do **not** carry into target architecture merely because they existed:

- **ADR-011 — generic shared Divergences ledger:** target structure superseded. Divergence/business correctness stays with the originating D1 authority; actionable lifecycle belongs to Operational Work. Preserve honest evidence/history where material, not a generic divergence business owner.
- **ADR-012 — DIFAL table inside legacy `pricing`:** legacy module/table shape superseded. Retain only the single-authority/no-fabricated-rule principle. Commercial Economics owns economic interpretation; external fiscal facts/rules remain external and D4 re-verifies source contracts.
- **ADR-022 — `SELLER_SKU == CODPROD` provider-write invariant:** superseded as canonical identity invariant. CODPROD and SELLER_SKU are external identifiers/evidence, not universal identity. Readiness correspondence + D4 provider contract must prove the concrete write mapping; Mercado Livre may re-establish a provider-specific rule in D4 if current evidence requires it.
- **ADR-028 — automatic link only from CODPROD+EAN concordance:** superseded as target identity law. Preserve the safety principle that weak evidence cannot fabricate correspondence. Product & Channel Readiness owns matching/correspondence; D4 supplies current source/provider identifiers/evidence.
- **ADR-031 — `products_mirror` keep-absent merge:** legacy table/mechanism superseded. Preserve honest-absence/partial-observation semantics already carried by accepted architecture/ADR-027/Fact scope; no canonical MPC Product mirror is retained merely for that mechanism.

### B2.10 — Legacy ADR set is not the new system's ADR baseline

**Operator direction:** pre-rebaseline ADRs are architecture history of the old system, not a catalog to carry forward into the rebuilt target.

Target policy proposed for review:

1. No pre-rebaseline ADR becomes target architecture merely by being currently marked binding/reopened/superseded.
2. During D0–D9, D-stage artifacts + stable `ARCHITECTURE.md` constraints remain the target authority path.
3. Any still-needed safety/product constraint currently discoverable only through an old ADR must be explicitly rehomed into accepted D-stage/stable architecture or a **new target ADR** before the legacy file is removed.
4. After the rebaseline has enough integrated closure (recommended no later than D9 closure, earlier only when dependency-safe), create a clean **new ADR baseline/series** from the accepted target decisions that genuinely merit ADRs; do not preserve old ADR numbers for compatibility.
5. Then delete legacy ADR files/legacy registry/citation archaeology from the active tree when nothing in the active authority path depends on them. Git history remains the archive.
6. The deletion is documentary cleanup, not a semantic migration: old ADR content stays available in Git history but has no runtime/target authority.

This policy deliberately avoids both extremes: carrying stale ADRs into the new system, and deleting the only current copy of a still-binding safety constraint before it is rehomed.

---

## Requested independent challenge

Return only material findings. Attack especially:

1. Does opaque/stable/non-reusable identity leave any missing identity invariant needed before D3?
2. Is `(issuer, subject) → Principal` uniqueness/migration sufficient without accidentally designing IAM implementation?
3. Does SourceInstance have the smallest correct lifecycle, or does it still leak D4 configuration into D2?
4. Does typed-reference policy cause an impossible polymorphic case outside Work, or correctly prevent a generic entity graph?
5. Is exact quantity semantics justified and correctly bounded, or should any piece defer?
6. Does identity lifetime/deactivation preserve history without smuggling a universal soft-delete model?
7. Are product-defined Roles/Permissions the correct Product 1.0 Global Maximum?
8. Is any universal evidence identity actually required by accepted D0/D1/D2, or is B2.8 correct YAGNI?
9. Re-adjudicate ADR-011/012/022/028/031 against **current target authority**, not their old implementation. Identify any semantic invariant we would accidentally lose.
10. Attack B2.10 hard: can the legacy ADR active tree be retired and replaced by a clean target ADR baseline without breaking the repository authority chain? Identify every class of constraint that must be rehomed first.
11. Identify any B2 item that is actually D3/D4/D5/D7 mechanism and should be narrowed/deferred.
12. If B2 survives, determine whether another material D2 batch is needed or whether D2 can move directly to final Global Coherence + YAGNI / Overengineering / Future-Cost review.

## Expected response shape

- **VERDICT:** APPROVE / REVISE / REJECT
- material findings only, with IDs;
- corrected invariant(s) where needed;
- explicit old-ADR disposition corrections if any;
- reopen triggers;
- final statement: `READY FOR D2 GLOBAL COHERENCE: YES/NO`.
