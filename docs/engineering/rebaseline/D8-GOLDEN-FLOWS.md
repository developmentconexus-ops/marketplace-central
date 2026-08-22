# D8 — Golden Flows

> **Status:** OPEN / ACTIVE — DERIVED CANDIDATE / OPERATOR-APPROVED FLOW SET  
> **Program:** Architecture Rebaseline / Technical System Design  
> **Opened:** 2026-08-22  
> **Accepted prerequisites:** D0–D7 — ACCEPTED / CLOSED  
> **Method:** DevelopmentConexus Engineering Method v1.0.0

## 1. Purpose

D8 defines the smallest set of representative **golden flows** needed to prove that the accepted D0–D7 Product, ownership, API, frontend and runtime authorities compose coherently end to end before the adversarial D9 review.

D8 owns flow selection, choreography and proof expectations only. It does **not** reopen accepted Product operations, Permissions, Principal kinds, semantic owners, provider contracts, frontend authority or D7 runtime mechanisms by convenience, and it does not begin Product implementation.

## 2. Accepted baseline

```text
D0–D7                  ACCEPTED / CLOSED
Product operations     99
ordinary Permissions   30
Principal kinds        H / A / S only
stable origin          https://conexus.fun
active runtime         NONE
D9                      BLOCKED
Product implementation BLOCKED UNTIL D9
```

Canonical Product wire authority remains [`contracts/api/product/openapi.yaml`](../../../contracts/api/product/openapi.yaml).

## 3. Selection rule / target invariant

> **A D8 flow earns a place only when removing it would leave a material accepted invariant without a representative cross-boundary falsifier.**

D8 therefore selects by **defect class**, not by domain count, screen count or the 99-operation inventory.

The accepted candidate is deliberately irreducible:

```text
3 business golden flows
+ 1 systemic recovery falsifier
```

A flow may carry multiple positive, negative, ambiguity and recovery variants. Authentication, Organization isolation, Permission/Principal admission, idempotency, concurrency, durable-effect safety and frontend/wire composition are cross-cutting controls applied to those flows rather than separate business journeys unless their failure mode cannot be represented honestly inside one.

## 4. Proof horizons

D8 distinguishes two proof horizons so architecture evidence does not become a false implementation claim.

### 4.1 D8 architecture / external-contract proof

D8 may use:

- accepted authority trace and canonical OAD inspection;
- counterexample/falsifier analysis;
- deterministic repository proof;
- bounded real external-system probes explicitly deferred to D8 by accepted D4 authority.

D4 deferred the first controlled real Mercado Livre publication/write proof and selected Sankhya fiscal/destination progression proof to D8. Those probes require separate explicit operator authorization before any live consequential write. They are external-contract evidence, not a Product-runtime implementation.

### 4.2 Post-D9 implementation conformance

The eventual implementation must execute the same protected properties through the real Product runtime. PostgreSQL/RLS, River, Keycloak, browser/CSRF, router/OAD validation, object storage and restart/recovery claims cannot close on mocks alone. D8 defines the acceptance target; it does not populate that runtime while implementation remains blocked.

---

# 5. GF-01 — Publication & Marketplace Convergence

## 5.1 Purpose

GF-01 proves the highest-leverage consequential marketplace-write path while preserving separate Readiness, Offering, Availability and provider authorities.

Owning references are bounded to the canonical OAD plus [D4-R1 Publication Input](D4-R1-PUBLICATION-INPUT.md), [D6 Frontend](D6-FRONTEND.md) and, for effect mechanics only, [D7 Runtime](D7-RUNTIME-JOBS-TRANSACTIONS.md) / [D7-C](D7-C-DURABLE-WORK-EXTERNAL-EFFECTS.md).

## 5.2 Choreography

```text
explicit Organization
→ exact Marketplace Installation
→ explicit SourceInstance-qualified Product
→ Readiness search / requirements / correspondence
→ Offering ListingIntent draft
→ source-following or explicit authored values
→ authored media when needed
→ separate Offering PriceIntent when needed
→ Availability-owned sellable-availability meaning for the target
→ SubmitListingIntent
→ owner-local durable intake / external-effect handoff
→ Mercado Livre provider realization
→ authoritative provider reread
→ Offering evaluates Listing convergence
→ Availability evaluates Availability convergence
```

No step creates a Publication/PIM business owner. One physical provider request may serialize meanings issued by several owners without moving authority into the adapter/runtime mechanism.

## 5.3 Material acceptance properties

GF-01 must preserve all of the following:

1. `FOLLOW_SOURCE` and `EXPLICIT_OVERRIDE` remain distinct; source truth is never rewritten by an override.
2. ListingIntent never owns or fabricates Availability quantity.
3. PriceIntent remains separate Offering meaning; Market/Economics/Performance may inform but never acquire price-write authority.
4. Active provider creation requiring Availability fails closed **before dispatch** when the Availability-owned input/proof is absent.
5. shared User Product / provider blast radius is revalidated before a widened effect can leave MPC.
6. provider `2xx/201/202` is no stronger than accepted/pending evidence; authoritative reread controls convergence.
7. an ambiguous possibly accepted external write is reconciled and is never blindly redispatched.
8. stale revision/precondition remains distinct from business/provider rejection.
9. exact prior idempotent intake is replay-safe; reuse of the key for different semantic input fails closed.
10. recurring automation cannot silently replace a standing human-authored override/resolution.
11. frontend composition consumes canonical Product operations and never adds a screen-shaped Product endpoint or raw provider write path.

## 5.4 Required variants / falsifiers

At minimum the proof matrix must be capable of falsifying:

| Variant | Failure that must remain impossible / explicit |
| --- | --- |
| happy controlled publication | transport success being mistaken for convergence |
| missing Availability-issued input | adapter defaults/infers quantity and dispatches anyway |
| shared provider resource | one narrow intent silently widening provider blast radius |
| stale ListingIntent | silent overwrite instead of accepted precondition semantics |
| duplicate exact intake | duplicate business/effect creation |
| same key, changed meaning | idempotency key becoming business identity |
| timeout/crash after possible dispatch | automatic redispatch instead of reconciliation |
| human override + recurring A Principal | automation silently superseding standing human meaning |
| Organization/reference mismatch | cross-Organization source/intention/reference resolution |

## 5.5 D8 real external gate

Accepted D4-R1 specifically leaves the **first controlled real Mercado Livre creation/write** to D8. When separately authorized, the probe must include authoritative reread plus shared-User-Product blast-radius verification and must not be used to infer a Product-runtime PASS.

---

# 6. GF-02 — Sale → Business System → Fiscal → Fulfillment → Outcome

## 6.1 Purpose

GF-02 proves the core post-sale operating loop across Marketplace Sales, Business-System Materialization, Fulfillment, Sankhya fiscal progression, Shipment observation and bounded exception consequences without creating one cross-owner workflow authority.

Owning references are bounded to the canonical OAD plus [D4 External Integrations](D4-EXTERNAL-INTEGRATIONS.md), [D6 Frontend](D6-FRONTEND.md) and the exact D7 slice required for a runtime falsifier.

## 6.2 Choreography

```text
Mercado Livre Sale evidence
→ authoritative Sale reread
→ transaction Selling Entity attribution when required

→ Business-System Materialization
   → Business Order Intent already derived by accepted owner semantics
   → Party Resolution
   → Destination Realization or explicit external-required / Work
   → sanctioned Sankhya TOP-313 materialization/progression
   → authoritative native reread / convergence

→ Fulfillment
   → Separation
   → Physical Conference

→ Invoicing Intent already derived by accepted owner semantics
→ sanctioned Sankhya 313 → 306 fiscal progression
→ authoritative native reread + origin/result correlation

→ provider-required artifacts/readiness
→ Packing
→ Dispatch Handoff
→ source-qualified Shipment observation
→ realized Economics / reconciliation
```

There is no Product `CreateBusinessOrderIntent`, `CreateInvoicingIntent`, generic Sankhya command, generic workflow progression command or direct-provider shortcut invented for the UI.

## 6.3 Exception branch — same golden flow, not a fourth business flow

```text
Claim / Return / Refund / reverse Shipment evidence
→ PostSaleResolution coordination
→ physical and economic consequences remain independently evidenced
→ Operational Work when ambiguity/manual action remains
```

Post-Sale does not absorb Sales, Fulfillment or Commercial Economics. Work coordinates responsibility; closing Work never manufactures source-domain closure.

## 6.4 Material acceptance properties

GF-02 must preserve all of the following:

1. multiple compatible/contradictory native party candidates remain `AMBIGUOUS`; never first-result-wins.
2. unsafe/unproven destination realization becomes explicit `external-required` / Work; it never silently overwrites customer master data or creates another Party merely to hold an address.
3. Sankhya target transport remains sanctioned API Gateway only; Direct Oracle is never fallback/recovery.
4. native order/progression `2xx` does not prove materialization convergence without authoritative reread.
5. physical readiness/conference remains Fulfillment authority and gates normal invoicing where accepted.
6. ordinary `A` cannot establish physical checkpoints; `S` can establish only the exact admitted physically-qualified checkpoints with current server-owned qualification.
7. Invoicing result remains source-qualified and distinct from its MPC-owned InvoicingIntent.
8. refund/payment evidence does not fabricate physical post-sale closure; physical return evidence does not fabricate financial/economic closure.
9. ambiguous consequential Sankhya/provider effect is reconciled rather than blindly replayed.
10. Sale-detail/read composition does not acquire write/workflow authority from its component reads.

## 6.5 Required variants / falsifiers

| Variant | Failure that must remain impossible / explicit |
| --- | --- |
| one exact native party | valid resolution and correlation are preserved |
| multiple active native parties | guessed/first candidate selection |
| unsafe destination | destructive Partner-master mutation or silent destination loss |
| TOP-313 possible-acceptance timeout | duplicate native order from blind retry |
| pre-physical invoice attempt | invoicing bypasses Fulfillment-owned readiness |
| physical checkpoint by A | client class/Permission treated as physical authority |
| physical checkpoint by unqualified S | token claim/body self-asserts physical qualification |
| 313→306 fiscal progression | `2xx` treated as converged invoice without correlation/reread |
| post-sale refund + reverse shipment | one consequence falsely closes the other |
| Work close | source-domain condition declared resolved only because Work closed |

## 6.6 D8 real external gates

Accepted D4 leaves two related first-lane proofs to D8 when separately authorized:

- controlled Destination Realization showing Party/master correctness, progression/invoicing survival, fiscal/XML correctness and no unrelated-state corruption;
- first selected-lane sanctioned `TOP 313 → TOP 306` fiscal effect with authoritative reread/correlation.

These are irreversible/legal external effects and require explicit operator authorization separate from approval of this architecture document.

---

# 7. GF-03 — Performance Evidence Honesty

## 7.1 Purpose

GF-03 proves the **Commercial Intelligence** half of Product 1.0 and the bounded D6-R1 99/30 amendment. The other two flows are effect-heavy; without GF-03 there is no representative falsifier for performance knowledge, measurement-basis and historical-evidence honesty.

Owning references are bounded to the canonical OAD, [D6-R1 Marketplace Performance Intelligence](D6-R1-MARKETPLACE-PERFORMANCE-INTELLIGENCE.md) and [D6 Frontend](D6-FRONTEND.md).

## 7.2 Choreography

```text
explicit Organization
→ exact Marketplace Installation
→ explicit primary reporting period
→ optional valid comparison period
→ GetMarketplacePerformanceSummary
→ ListMarketplaceListingPerformance
→ GetMarketplaceListingPerformance for selected Listing when needed
→ ListRetailMediaPerformance
→ frontend composes only admitted owner meanings
```

The proof may use a minimal subset per scenario; the four operations are one semantic proof family, not four separate golden flows.

## 7.3 Material acceptance properties

GF-03 must preserve all of the following:

1. known zero/empty remains distinct from `unknown`, `unavailable` and `unsupported`.
2. `partial` evidence is never presented as a complete-period aggregate.
3. all currently known Listings remain representable even when performance evidence for a Listing is unknown/unavailable.
4. comparison emits numeric change only when Performance declares evidence sufficiently complete and measurement basis comparable.
5. preserved historical provider evidence remains source-qualified; MPC custody never becomes MPC-authored provider fact.
6. frontend does not recalculate/redefine provider CVR/ROAS or other provider-defined measures as canonical truth.
7. campaign/catalog/family Retail Media evidence is not attributed to one Listing without sufficient evidence.
8. `performance.read` grants only Performance reads and never implies Market/Economics/Offering/Sales/Availability access.
9. no all-marketplace KPI aggregation appears without proven measurement equivalence.
10. no Ads mutation, optimize/sync/collect operation, opportunity-score DSL or AI authority appears.

## 7.4 Required variants / falsifiers

| Variant | Failure that must remain impossible / explicit |
| --- | --- |
| known zero | rendered as unknown/missing |
| unknown/unavailable | rendered as numeric zero or empty-known |
| partial coverage | displayed as full-period KPI |
| incompatible comparison basis | numeric delta emitted anyway |
| Listing with unavailable performance | omitted from Listing population / survivorship bias |
| campaign/catalog/family evidence | forced into Listing identity |
| `performance.read` only Principal | other strategy-owner data becomes implicitly readable |
| historical preserved evidence | presented as current/live or MPC-authored source fact |

---

# 8. SR-01 — PITR / Timeline Continuity Recovery

## 8.1 Why this is systemic rather than a fourth business flow

Database time rollback can erase the very dispatch/idempotency/authentication evidence used by ordinary continuous-timeline crash safety. That failure cannot be represented honestly as a normal user journey, so D8 keeps one explicit systemic recovery falsifier.

Owning authority is [D7-R1 Whole-Stage Coherence](D7-R1-WHOLE-STAGE-COHERENCE.md), with D7-C/D7-E only as needed for exact mechanism claims.

## 8.2 Falsifier

```text
external consequential write is possibly/actually accepted
→ restore database from a point before the acknowledged dispatch marker/state
→ ordinary application boot
→ NO manual recovery/fence flag is supplied
→ out-of-rollback-domain continuity witness cannot affirm safe continuous lineage
→ recovery fence arms automatically
→ consequential external dispatch remains disabled
→ restored work is reconciliation-only
→ authoritative provider/business-system state is reacquired
→ access/session/integrity state is revalidated as required
→ fence releases only for scopes whose safety is positively established
```

Complementary control:

```text
affirmatively continuous durable lineage
→ continuity witness verifies safe descent
→ recovery fence need not be falsely armed
```

A proof beginning with the fence manually pre-engaged is insufficient.

---

# 9. Cross-cutting acceptance controls

These controls apply to the selected flows and do not create additional business golden flows.

| Control | Binding D8 expectation |
| --- | --- |
| **Product surface conservation** | remains exactly **99 Product operations / 30 ordinary Permissions / H-A-S only**; a required operation 100 or Permission 31 is a material finding, not an implementation convenience |
| **Organization isolation** | path Organization is explicit; secondary references fail closed; eventual implementation must falsify omitted predicates, composite-FK crossing and pooled-scope leakage with real PostgreSQL/RLS |
| **Authentication / ordinary access** | H session+CSRF and A/S bearer remain distinct; Permission, Principal kind, physical qualification, business disposition and Governance are non-equivalent gates |
| **Canonical wire / frontend** | every externally invokable step maps to canonical OAD authority; frontend may compose reads but never invent DTO, route, Product operation or cross-owner write authority |
| **Idempotency / concurrency** | exact replay, changed-fingerprint rejection, opaque revision semantics and stale-precondition behavior survive UI/network retry paths |
| **Durable external effects** | owner state/handoff atomicity, dispatch marker, possible acceptance→reconciliation and no blind redispatch remain binding; River/job state is mechanism only |
| **Governance / Work** | Governance authorizes but never executes the target; Work coordinates responsibility but never becomes source truth or generic command bus |
| **Knowledge honesty** | unknown/partial/unavailable/unsupported/known-zero stay distinct wherever material; provider/source provenance and freshness are not reconstructed by UI convenience |
| **Technical exclusion** | no Direct Oracle fallback, generic provider/plugin/workflow authority, Product Ads management, screen-shaped API or generic retry/recovery Product surface appears |

# 10. Falsifiable acceptance matrix

| ID | Representative claim | D8 passes only if the candidate can expose this defect class |
| --- | --- | --- |
| **GF-01** | owner-preserving marketplace publication + convergence | hidden owner transfer, missing-input dispatch, stale/duplicate intake, widened blast radius or ambiguous write redispatch causes explicit failure/finding |
| **GF-02** | Sale→Sankhya→physical→fiscal→dispatch composition | guessed Party/destination, Oracle fallback, premature invoice, invalid physical actor or cross-owner false closure causes explicit failure/finding |
| **GF-03** | source-qualified Commercial Intelligence | unknown→zero, partial→complete, incompatible comparison, fabricated metric/scope or Permission leakage causes explicit failure/finding |
| **SR-01** | safety across acknowledged-state rollback | restored work cannot redispatch an already-possible effect and stale authority/session cannot be silently resurrected without affirmative continuity proof |

Removing any one row leaves a materially different accepted invariant class without representative falsification; adding a fourth business journey currently adds coverage ceremony rather than a distinct defect class.

# 11. Explicitly rejected D8 expansion

D8 does **not** create independent golden flows merely for:

- marketplace onboarding/channel settings;
- Availability configuration;
- Market/Economics analysis apart from the selected representative consumption;
- Governance administration;
- identity/access administration;
- Operational Work lifecycle;
- Post-Sale as a standalone journey;
- every D6 F1–F12 interaction flow;
- every one of the 99 Product operations;
- every D7 real-dependency falsifier as a separate business journey.

Those properties are either represented by stronger selected variants/cross-cutting controls or belong to later implementation conformance. Conversely, compressing GF-01/GF-02/GF-03 into one or two mega-flows is also rejected: it would reduce failure localization while not removing essential complexity.

---

# 12. D8-F1 — Architecture closure is not implementation-readiness closure

## 12.1 Finding

D8 composition review surfaced a material pre-implementation distinction:

> **Accepted architecture can be coherent while still leaving material realization choices to the coding agent. Implementation must not open while the agent still has to decide which user-visible states exist, which Product operation feeds an element/action, which owner receives a write, which dependency edge is allowed, or what observable result proves a slice complete.**

D6 correctly proved frontend architecture using the 99-operation capability mapping, 39 representative route/screen states, interaction laws, topology and representative low-fidelity wireframes. It did not claim to freeze every material screen/state needed for implementation. That was sufficient for D6 architecture closeout but is not sufficient for the newly explicit implementation-readiness requirement.

This finding does **not** invalidate React, TanStack Router/Query, the generated OAD transport, D6 information architecture or accepted backend authority. No D6 reopen occurs inside D8.

## 12.2 Operator-approved post-D8 sequence

After D8 closes, and before D9 opens, the program will perform the smallest targeted completion:

```text
D8 close
→ D6-R2 — Complete Frontend Realization Closure
→ Pre-D9 Implementation Readiness Contract
→ D9 — Adversarial Architecture Review
→ Product implementation only after accepted D9
```

D6-R2 is a targeted realization-completeness reopen, not a redesign by preference.

## 12.3 D6-R2 minimum obligations

D6-R2 must derive from current 99/30 + D7 authority rather than from a visual-first exercise. For every material user surface it must close, proportionately:

```text
accepted Product capability / user job
→ route + screen/state inventory
→ semantic owner(s)
→ exact Product operations
→ exact Permissions / Principal constraints
→ server / URL / form / ephemeral state ownership
→ actions and write owner
→ knowledge / freshness / effect / recovery states
→ navigation and Organization / Installation / Source qualification
→ wireframe/state variants
→ backend/OAD coherence check
```

All material screens must be defined, including loading/empty/unknown/unavailable/partial, permission/access differences, stale/concurrency, business rejection, pending/ambiguous/converged outcomes and responsive behavior where meaning could otherwise be lost. Multiple material states may share one annotated/base wireframe; the requirement is unambiguous realization, not one image per combinatorial state.

A screen that needs meaning absent accepted backend/OAD authority stops and returns to the **smallest owning authority**; frontend never fabricates the missing contract.

## 12.4 Realization Contract

Each material screen and each later implementation slice must have a concise Realization Contract answering at least:

- expected outcome / user-visible surface;
- owning accepted authorities;
- exact Product operations and external/non-Product seams;
- exact Permissions / Principal classes / physical qualification when relevant;
- data/state ownership and persistence classes required;
- runtime/auth/transaction/job/recovery properties materially exercised;
- allowed and forbidden dependency/import directions;
- negative controls / things that must remain impossible;
- executable acceptance criteria;
- real-dependency proof required where mocks are insufficient;
- **exit state:** what must demonstrably exist when the slice is complete.

The implementation may choose local/private mechanics that do not alter accepted meaning or structure; it may not choose new material UX behavior, authority, contract, dependency direction or acceptance outcome.

## 12.5 Import/dependency graph becomes an implementation acceptance property

Accepted D6 direction remains conceptually:

```text
app/routes
    ↓
features
    ↓
api/<owner-family>
    ↓
api/transport
    ↓
api/generated

features → ui
```

The pre-D9 readiness contract must convert material dependency directions into enforceable acceptance rules rather than rely only on prose. Representative forbidden edges include feature-owned raw Product fetches, feature→generated bypass, generic UI→API authority and reverse dependencies from lower transport/generated layers into feature/application semantics.

## 12.6 D9 entry challenge

D9 must challenge the **implementation-ready package**, not architecture prose alone. One explicit adversarial test is:

> **Can two materially different implementations both satisfy the same specification because a behavior, authority, API/owner connection, state transition, dependency direction or completion criterion was left undecided?**

If yes, the specification still contains a material realization ambiguity and returns only to the smallest implicated owner/gate before implementation opens.

---

# 13. Candidate disposition

Current D8 derivation result:

```text
D0–D7 reopen                 NONE
Product surface              99 operations / 30 Permissions / H-A-S unchanged
business golden flows        3
systemic recovery falsifier  1
new business authority       NONE
new Product operation        NONE
new Permission               NONE
Product implementation       BLOCKED
D9                            BLOCKED
post-D8 realization finding  D6-R2 + Pre-D9 Implementation Readiness Contract
```

**Candidate outcome:** `CURRENT STRUCTURE CONFIRMED` with an implementation-readiness completion prerequisite before D9.

## 14. Exact next action

**Independently challenge this derived D8 candidate, then adjudicate only material findings against current authority. Do not begin D6-R2, D9 or Product implementation. Do not execute a live Mercado Livre or irreversible Sankhya write without separate explicit operator authorization.**
