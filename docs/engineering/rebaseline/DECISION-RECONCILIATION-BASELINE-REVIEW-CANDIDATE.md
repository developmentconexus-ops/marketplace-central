# Decision Reconciliation Baseline — REVIEW CANDIDATE

> **Status:** NON-AUTHORITATIVE REVIEW CANDIDATE  
> **Purpose:** reconcile the accepted D0→D4 + D4-R1 architecture and already-accepted D5-B1 contract laws before D5-B2 turns the design into concrete API/code-facing contracts.  
> **Prepared against:** `1adc22e8215de26631b19bf9a7a2e0291f7579f1` on `docs/global-methodology-alignment`  
> **Method:** DevelopmentConexus Engineering Method v1.0.0  
> **Important:** this candidate does not replace D0–D5 artifacts. It proposes a canonical **decision-disposition/reconciliation index**, not a second semantic architecture.

---

# 1. Why this baseline exists

D0→D4 were intentionally iterative. Later decisions refined, superseded or split earlier/legacy choices. D4-R1 was discovered only when D5-B2 attempted to derive the concrete Product API surface. D5-B1 is already accepted and D5-B2 will begin binding architecture to routes, schemas, Permissions and implementation-facing contracts.

The repository therefore needs one clean reconciliation point that answers:

> **Which decisions are active now, where is their semantic authority, which older decisions are dead, what remains genuinely open/deferred, and what may the technical stages implement without re-deciding architecture?**

Without such a baseline, a fresh reviewer can accidentally:

- resurrect a legacy ADR because the file still exists;
- follow a stale review candidate;
- treat an old D-stage intermediate formulation as current target meaning;
- implement both an old and a new model to “be safe”;
- add compatibility/framework abstractions for decisions that were deliberately superseded;
- miss a bounded Unknown/Deferred item and silently invent a default;
- reinterpret a D0–D4 authority decision while designing D5/D6/D7 implementation.

This is an authority-routing problem, not a new product-design problem.

---

# 2. Authority rule for this baseline

If accepted, this baseline would own only **decision reconciliation/disposition**:

- which decision generation is current;
- which accepted artifact is the semantic home;
- which legacy/candidate artifacts are superseded/historical/disposable;
- which questions remain intentionally open and who owns their later adjudication.

It would **not** own detailed product/domain semantics already defined by D0–D4/D4-R1/D5-B1.

Conflict rule:

1. router owns current program status/next action;
2. `ARCHITECTURE.md` owns stable cross-stage constraints;
3. accepted D-stage artifacts own detailed semantics in their stage scope;
4. this baseline owns only the reconciled **which-decision-is-current** map;
5. legacy ADRs, review candidates, `AI-DIALOG.md`, code/OpenAPI/schema/history remain evidence only unless an active authority explicitly carries their meaning.

If this baseline ever starts restating enough detail to disagree with the underlying authority, the baseline is defective and must shrink rather than become a competing architecture document.

---

# 3. Reconciled target architecture — current decision generation

## 3.1 Product/system boundary — CURRENT

**Authority:** D0 + stable `ARCHITECTURE.md` constraints.

Current decisions:

- MPC is the internal **Marketplace Operations Control Plane + Commercial Intelligence** product.
- External systems retain authority for facts/processes inherently theirs.
- MPC owns cross-system marketplace operating semantics: observe → understand → reconcile → decide/policy → execute → verify → audit/reconcile.
- Mercado Livre is first marketplace proof; Sankhya is first business-system proof.
- Organization is explicit even while the first proof may use one Organization.
- MPC is **not** an ERP replacement, PIM/MDM, generic marketplace hub, generic integration platform, universal workflow engine, unrestricted autonomous-agent platform or enterprise-wide control-plane framework.
- A future enterprise-generic repositioning is a D0/D1 reopen, not an implementation generalization.

**Superseded interpretations:** dashboard-only product, ERP-shaped product, generic provider/plugin platform, Product/PIM hub.

---

## 3.2 Business authority boundaries — CURRENT

**Authority:** D1.

Twelve semantic business boundaries remain current and justified:

1. Marketplace Portfolio
2. Product & Channel Readiness
3. Marketplace Offering Operations
4. Availability Control
5. Market Intelligence
6. Commercial Economics
7. Controlled Action Governance
8. Marketplace Sales
9. Business-System Materialization
10. Fulfillment Lifecycle
11. Post-Sale Resolution
12. Operational Work

Cross-cutting accepted non-domain authority:

- D2 identity/access substrate.

These are **semantic authorities**, not service/process/database counts.

Current hard boundary laws:

- one business meaning → one semantic authority;
- mechanism ≠ authority;
- Product master remains external;
- no generic Integration, Evidence, Mutation/Action, SLA, Workflow or Policy business domain;
- Governance owns consequential authorization only, not business disposition/policy;
- Work owns responsibility/escalation/work state, not source business truth;
- Offering owns Listing/Price intent, not Sellable Availability;
- Economics owns economic interpretation, not marketplace write;
- Materialization owns Business Order/Invoicing intents, not physical fulfillment;
- provider API resource grouping never merges MPC business authorities.

**Superseded structures:** legacy module names as domain boundaries; `mutations` as business owner; `catalog` as MPC Product master; `integrations` as universal business context; `dashboard` as authority; `sync` as business authority.

---

## 3.3 Identity, tenant and data ownership — CURRENT

**Authority:** D2.

Current decisions:

- Organization is the tenant/isolation root.
- MPC-owned canonical IDs are stable opaque IDs.
- Marketplace Installation is Organization marketplace participation/configuration identity.
- Selling Entity is Organization-scoped operational identity; Portfolio owns registry/lifecycle, Sales owns transaction attribution.
- Inventory Source and Fulfillment Node remain distinct MPC identities.
- SourceInstance minimally qualifies externally authoritative business-system namespaces.
- Principal is the canonical accountable actor; human/automation/system meanings remain distinct.
- Product = `SourceInstance + native Product key`; no MPC Product mirror/master.
- provider Listing/Variation, Sale/Order, Shipment and native financial movements remain external source-qualified identities.
- material domain-local Business Intents have owner-local MPC identity; no universal Intent/Action/Command entity.
- Business Order Intent, Invoicing Intent, Post-Sale Resolution, Authorization Decision, Work and other accepted owner-local identities remain current.
- identity/access substrate owns Membership, product-defined AccessRole/Permission definitions and RoleAssignment; it does not own business disposition or consequential authorization.
- every Organization-owned durable state/evidence is explicitly Organization-scoped; Organization is not inferred from Installation/SourceInstance/provider account/default process state.
- unknown/empty/absent/not-applicable remain distinct when material.
- historical snapshots may preserve decision context without becoming current authority.
- pre-rebaseline persistence has no compatibility/migration claim on the target.
- automation recurrence does not silently reverse standing human decisions in the same semantic scope.

**Superseded structures:** separate canonical Tenant identity; Product mirror; bare native external identity; universal entity/evidence graph; generic shared mutable Product/Listing/Order; legacy CODPROD/EAN identity formulas as universal laws.

---

## 3.4 Communication and failure semantics — CURRENT

**Authority:** D3.

Current grammar:

- **Q** — current owner meaning;
- **C** — request owner capability/work;
- **E** — committed producer occurrence with independent consumer reaction;
- **P** — read-only composition across authorities.

Current laws:

- communication never transfers authority;
- Organization scope remains explicit in communication/recovery state;
- duplicate delivery is allowed; semantic idempotency prevents duplicate business effects;
- no global ordering/exactly-once assumption;
- known/known-empty/unknown/unavailable are distinct;
- freshness/provenance is orthogonal to knowledge state and travels when material;
- `accepted != completed != externally applied != converged`;
- capability outcomes may be accepted/rejected/pending/ambiguous;
- stale provider/domain preconditions may be definitive rejection rather than ambiguity;
- no blind replay after possible external acceptance;
- projections never become write/concurrency/current-truth authority;
- cross-owner workflows are correlated/convergent; no cross-owner atomicity requirement;
- reusable external-effect safety mechanics verify owner-issued proofs without owning business validity/authorization.

**Superseded structures:** generic Command/Event Bus as business architecture; generic Mutation envelope; one global workflow state machine; event sourcing/global event log as mandatory system of record; polling/sync mechanics as business semantics.

---

## 3.5 External integration boundary — CURRENT

**Authority:** D4-B1/B2/B3/B4 + D4-R1.

Governing rule:

> **consumer owns meaning; adapter owns protocol.**

Current decisions:

- D4 is acquisition/protocol/capability/coverage/effect translation, not a business domain or generic evidence store.
- one provider acquisition may feed several consumer-owned semantic ports without one consumer/D4 owning the raw provider resource wholesale.
- credentials/auth are mechanism, not identity/business truth.
- notification/callback is acquisition evidence/pointer; authoritative reread establishes material current provider meaning.
- point/enumeration/delta/notification coverage is operation-scoped.
- Integration Support / Provider Effective Capability / Effective Business Capability remain distinct.
- external effects preserve owner intent/correlation, prerequisites, acceptance/ambiguity and authoritative reread/convergence.
- provider 2xx never proves business convergence by itself.
- provider-specific richness is kept when a named consumer/correctness need exists; no lowest-common-denominator flattening and no raw DTO mirroring.

### Mercado Livre current realization

- Item/User Product/Family/Catalog/stock topology remains provider-local.
- shared User Product effects cannot silently widen intended/authorized scope.
- Listing/Price meaning remains Offering-owned; Sellable Availability remains Availability-owned.
- Sale/Shipment remain distinct external resources.
- fulfillment/provider requirements remain contextual evidence.
- essential post-sale provider resources remain source-qualified evidence; Post-Sale Resolution owns MPC closure semantics.
- first current proof lane is time-bound and revalidated when consequential.

### Sankhya current realization

- sanctioned API Gateway is the **only target transport**; Direct Oracle/database is not fallback.
- Product/company/location/inventory/control/cost/tax/party/document topology remains provider-local evidence.
- bounded sanctioned reads cannot become arbitrary SQL escape hatches.
- Party Resolution and Destination Realization are bounded Materialization prerequisites, not Customer/Address masters.
- Business Order/Invoicing Intent remain MPC meaning; TOP/NUNOTA/CACSP/SelecaoDocumentoSP/etc. remain adapter realization.
- Expected Tax uses the sanctioned fiscal engine under a revalidatable SourceInstance binding; MPC does not duplicate tax rules.

### Market/economics/settlement

- Semantic Core + Provider-Enriched Evidence remains target.
- expected economics, Order economics and realized/settlement evidence remain separate rungs.
- Payment approval, release, refund/reversal, withdrawal/payout and Bank Cash remain distinct.
- no generic Fee/financial ledger or generic market collector framework.

**Superseded structures:** Direct Oracle/godror target integration; generic provider/plugin catalog; universal ERP/Customer/Address model; CollectorPort framework; lowest-common-denominator provider model; generic financial ledger; raw provider payload business entities.

---

## 3.6 Publication/listing authoring seam — CURRENT

**Authority:** D4-R1 + D1/D2/D3 parents.

Current decisions:

- external Product remains external/source-qualified.
- Readiness owns publication requirements, Product↔channel correspondence, source candidates and **source-level readiness**.
- Offering owns `ListingIntent` as the one **create/edit authoring/draft identity** and owns **draft dispatchability** using current Readiness meaning through the existing Readiness→Offering Q edge.
- no separate `PublicationPreparation` aggregate.
- listing values use `FOLLOW_SOURCE | EXPLICIT_OVERRIDE` at baseline only.
- no generic DERIVED/rules/mapping DSL.
- source acquisition remains D4 mechanism/evidence feeding consumer-owned ports; no generic SourceProductObservation business owner.
- embedded Sankhya source adapters need no self-HTTP; future physically external connector ingress is a prepared seam and gets a wire contract only when a real connector exists.
- human and automation authoring use the same Semantic Product API authority; automation cannot impersonate source truth or silently reverse human overrides.
- media may be source-qualified or listing-specific MPC authoring without creating Product-media master.
- provider requirement/schema changes remain D4 evidence + historical intent context, not universal ProductAttribute ontology.
- provider execution can jointly realize multiple owner-issued meanings in one physical request without merging business ownership.
- Mercado Livre initial publication × Availability is **PASS-B**: Offering never owns quantity; Availability issues its own meaning/input; execution machinery may serialize both in one provider request; each owner separately evaluates convergence.
- publication create/edit may be multi-step/partial/asynchronous; no `createListing = success` simplification.

**Superseded/rejected:** Readiness-owned PublicationPreparation; SourceProductObservation owner/service; generic source-ingestion API now; PIM/ProductAsset master; generic mapping/rules engine; AI-specific API/backdoor; separate create vs edit architectures.

---

## 3.7 API laws already accepted before B2 — CURRENT

**Authority:** D5-B1 + `ARCHITECTURE.md`.

These are implementation-facing architecture laws already frozen and must not be redecided by B2:

- semantic/domain-oriented Product API distinct from protocol ingress;
- Organization-owned Product API under `/organizations/{organization_id}/...`;
- secondary Organization-owned references fail closed across Organizations;
- Principal/auth context is not a client-authored business field;
- Permission/invoke access remains distinct from domain disposition/Governance;
- source-qualified external identity on the wire; no bare native correlation key;
- honest knowledge/freshness/provenance on reads;
- accepted/rejected/pending/ambiguous business outcomes distinct from API transport problems;
- fail-closed idempotency key by default for consequential intake unless structural idempotency is explicitly proven;
- concurrency tokens only where stale client state is materially unsafe;
- RFC 9457 Problem Details for API-level failures;
- one machine-readable Product API wire authority: OpenAPI;
- clients derive/conform; server conforms; conformance controls must be proven to fire;
- hard cutover/no compatibility tax absent an entitled consumer;
- bulk only when a real operation/workflow proves it.

**Superseded:** manual OpenAPI+SDK dual authority; legacy routes/SDK as target; provider/integration resource vocabulary as Product API; generic `/mutations`/`/commands`; speculative versioning/compatibility framework.

---

# 4. Decision lineage — old idea → current home

| Historical/earlier idea | Current disposition | Current semantic home |
|---|---|---|
| Marketplace dashboard as product | superseded | D0 control-plane product loop |
| Product mirror / MPC Product master | superseded/rejected | D2 external Product identity + D1 Readiness consumers |
| `Tenant` as separate canonical identity | superseded | D2 Organization root |
| generic `Integration` business domain | rejected | Portfolio business meaning + D4 protocol + D7 mechanics |
| provider plugin/self-registration framework | superseded | concrete D4 adapters; shared mechanism only if later proven |
| Direct Oracle / godror Sankhya target | superseded/historical | D4 sanctioned Gateway only |
| `SELLER_SKU == CODPROD` | superseded as identity law | D2/D4 evidence + Readiness correspondence |
| CODPROD+EAN unattended auto-link formula | superseded | D2 corroboration safety + Readiness policy |
| generic `Mutation` business owner | superseded | domain-local intents + Governance + D3/D7 execution safety |
| generic divergence ledger as truth owner | superseded | source domain correctness + Work lifecycle |
| `sync`/polling phase as product semantics | superseded | domain freshness/coverage + D4/D7 runtime mechanics |
| provider DTO/resource model in core | rejected | consumer-owned semantic ports + adapter-local protocol |
| low-common-denominator marketplace model | rejected | semantic core + provider-enriched evidence |
| generic CollectorPort/market source framework | superseded | explicit source admissibility + D7 mechanics only if needed |
| generic Customer/Address master for Sankhya | rejected | bounded Materialization Party/Destination realization |
| MPC tax engine | rejected | Sankhya sanctioned fiscal engine + Economics interpretation |
| global Fee/Payment/Settlement MPC entity | rejected | source-qualified movements + Economics attribution/reconciliation |
| PublicationPreparation aggregate | rejected by coherence | Offering ListingIntent draft + Readiness Q |
| SourceProductObservation business service | rejected by coherence | D4 acquisition feeding consumer-owned ports |
| generic listing transformation/rules engine | rejected/YAGNI | FOLLOW_SOURCE / EXPLICIT_OVERRIDE; explicit future reopen if repeated need |
| dedicated AI business/API authority | rejected | D2 automation Principal + ordinary Product API/Governance |
| Product API listing quantity owned by Offering | rejected | Availability-owned meaning; joint technical realization only |
| manual OpenAPI + manual SDK authority | superseded | OpenAPI one wire authority + derived/conformant clients/server |
| compatibility/versioning to preserve legacy API | rejected now | hard cutover until a real entitled consumer appears |

---

# 5. What is genuinely still open — do not silently decide in D5-B2

## D5 current/open

- exact Product 1.0 operation/resource inventory;
- exact request/response schemas;
- exact Permission→operation map;
- exact paths/nouns after B1 laws;
- concrete pagination/filter/sort/cursor only where consumers prove need;
- operation-local bulk only where proven;
- exact OpenAPI/server/client generation/conformance tooling.

## D6 deferred

- screen/navigation/editor topology;
- projection/view composition required by UX;
- frontend feature/package structure.

## D7 deferred

- process/server/worker/scheduler topology;
- transaction/outbox/queue/cursor/lease/lock mechanics;
- token refresh/caching/secret realization;
- RLS/runtime Organization isolation enforcement realization;
- idempotency persistence/TTL/locking;
- retry/backoff/rate-control mechanics;
- media blob/cache/CDN realization;
- production deploy topology;
- remaining reopened runtime ADR questions.

## D8 proofs/deferred concrete capabilities

- first controlled Mercado Livre create/Price/Availability writes + authoritative reread/convergence;
- shared-User-Product blast-radius real write proof;
- selected-lane Sankhya irreversible fiscal `313→306` progression;
- controlled alternate destination/contact realization before claiming it;
- first consequential native party create/update if the golden flow reaches it;
- currently unproven fiscal/provider branches only if a selected golden flow materially depends on them.

## D9

- final adversarial architecture/system review;
- legacy transition residues still explicitly routed to D9;
- implementation remains blocked until D9 is accepted.

## Bounded Unknown/Deferred items that remain honest

- unselected Mercado Livre modes/configurations;
- paused/zero-quantity representation-first creation on the selected Installation;
- broader marketplace/payment account movement universe absent a consumer;
- R3 bank-side reconciliation until a bank source is accepted;
- Sankhya controlled-product marketplace path until selection/interchangeability is proven;
- post-invoice fiscal return and some reversal branches until required;
- unproven tax components/branches only when a claimed flow materially needs them;
- future reusable source→listing transformations, Product media library or external connector contract until a real repeated consumer appears.

Unknown/Deferred is **not** permission to invent a default implementation.

---

# 6. Implementation reconciliation guard — what D5+ may not re-decide

Technical stages may choose mechanisms only inside these fixed meanings.

They MUST NOT re-decide by convenience:

1. Product ownership — external source-qualified, not MPC master.
2. Organization isolation/root semantics.
3. the 12 D1 business authorities or their accepted semantic edges.
4. Q/C/E/P meaning and honest failure/knowledge semantics.
5. domain-local Intent ownership versus generic mutation/workflow.
6. Governance versus Permission versus business disposition.
7. Work versus originating business truth.
8. external/provider identities versus MPC canonical identities.
9. consumer-owned port / adapter-owned protocol boundary.
10. sanctioned Sankhya Gateway-only target transport.
11. no blind retry after ambiguous possible acceptance.
12. provider richness without provider DTO mirroring.
13. Offering/Readiness/Availability publication ownership split.
14. ListingIntent as one create/edit authoring identity.
15. FOLLOW_SOURCE / EXPLICIT_OVERRIDE baseline and human-override safety.
16. no hidden PIM/source-observation/rule/connector/AI framework.
17. joint technical realization may compose owner-issued meanings but cannot create ownership transfer/cross-owner hidden semantic edges.
18. owner-specific convergence after provider effects.
19. semantic Product API + explicit Organization scope + honest source-qualified identities.
20. OpenAPI as one wire authority and no manual duplicate SDK authority.
21. hard cutover while no entitled production client exists.

A technical design that seems to require violating one of these is not an implementation workaround; it is evidence for a targeted architecture reopen.

---

# 7. Proposed active-tree cleanup after reconciliation acceptance

The goal is to make a fresh review read current authority, not archaeology.

## 7.1 Disposable review artifacts

After this reconciliation is accepted:

- delete `D5-B2-REVIEW-CANDIDATE.md`; it predates D4-R1 and is explicitly stale/non-authoritative;
- delete `AI-DIALOG.md` after the reconciliation review/consolidation completes; Git history preserves reviewer evidence;
- future review candidates are deleted immediately after their accepted meaning is canonically filed.

## 7.2 Legacy ADR files

The ADR registry already says pre-rebaseline ADRs are history/evidence, not the new system's ADR baseline.

After independent review confirms no unique still-needed meaning is stranded, retire from the **active tree** every legacy ADR whose status is `historical`, `superseded by rebaseline`, or `carried constraint` whose active meaning is fully rehomed in accepted D-stage/`ARCHITECTURE.md` authority.

Proposed legacy files to **retain temporarily** because a future stage or migration gate still explicitly needs them:

- ADR-003 — D9 residue;
- ADR-008 — D7 deploy topology evidence;
- ADR-010 — D7 runtime/polling residue;
- ADR-017 — Fact-domain judgment evidence until replacement Fact ADR;
- ADR-018 — D7 execution-safety/runtime residue;
- ADR-026 — D7 scheduler/runtime residue;
- ADR-030 — D7 residue;
- ADR-034 — current Fact implementation/evidence anchor until new target Fact ADR;
- ADR-035 — rebaseline transition authority until D0–D9 closes.

Everything else is proposed for active-tree retirement **only if independent review confirms its current meaning is already fully rehomed**. Git history remains the archive; ADR numbers are never reused.

## 7.3 Registry after cleanup

The ADR registry should become a short transition index containing only:

- the retained unresolved/carried legacy files above;
- a statement that all other pre-rebaseline ADR files are retired from the active tree and available in Git history;
- any new target ADR-036+ entries when later decisions genuinely benefit from ADR treatment.

Do not maintain a giant prose transition history once the decision has been reconciled; that would recreate the pollution this baseline exists to remove.

---

# 8. Proposed fresh-session read shape after acceptance

```text
AGENTS.md
→ rebaseline README/router
→ DevelopmentConexus Method
→ ARCHITECTURE.md
→ Decision Reconciliation Baseline
→ ADR registry (only unresolved technical/transition residues)
→ current D-stage artifact(s) needed for the work
→ supporting Evidence Register
→ implementation evidence only when necessary
```

The reconciliation baseline gives the reviewer the current map; the underlying D-stage artifact supplies the detailed rule when needed.

No reviewer should need `AI-DIALOG.md`, deleted review candidates or dozens of already-superseded ADR files to discover current target architecture.

---

# 9. Proof / validation plan

Independent challenge should try to prove at least one of these false:

1. a current decision listed here disagrees with D0–D4/D4-R1/D5-B1;
2. an older decision marked superseded actually still owns unique target meaning;
3. the proposed ADR-retirement set would delete the only active evidence/constraint needed by D5/D6/D7/D8/D9;
4. a genuinely open decision is accidentally presented as closed;
5. a closed decision is accidentally presented as open and therefore invites duplicate implementation;
6. the baseline has become a second semantic authority rather than a disposition map;
7. a fresh reviewer following the proposed read shape can still reasonably implement two contradictory architectures;
8. deleting the stale D5-B2 candidate/AI dialog would remove required target authority rather than review history;
9. D5-B2 still lacks an architecture prerequisite after D4-R1 + this reconciliation;
10. implementation guard items exceed accepted authority and accidentally freeze D6/D7 mechanisms.

Structural Inversion test:

> If all current code/OpenAPI/modules/packages were deleted or inverted, the reconciled active decisions above must still follow from accepted D-stage authority.

---

# 10. Reopen triggers

Reopen this reconciliation only when:

- an accepted D-stage decision is amended/reopened;
- a new target ADR changes current disposition;
- a later proof invalidates a current architectural assumption;
- a future Product requirement changes D0/D1 ownership;
- a retained legacy ADR is adjudicated and can leave the active tree;
- D0–D9 closes and ADR-035/transition machinery can be retired.

Do not reopen for implementation naming, package layout, framework preference or because an old Git-history decision is rediscovered.

---

# 11. Candidate outcome

**Proposed outcome:** `CURRENT D0→D4/D4-R1 + D5-B1 DECISION SET RECONCILED / COHERENT`.

If independent review and operator ratification converge:

1. create canonical `DECISION-RECONCILIATION-BASELINE.md`;
2. add it to the authority/read path as the decision-disposition map;
3. simplify the ADR registry to unresolved/transition residues only;
4. retire fully rehomed legacy ADR files from the active tree;
5. delete stale `D5-B2-REVIEW-CANDIDATE.md`;
6. delete `AI-DIALOG.md` after this review evidence has served its purpose;
7. keep Git history as the sole archive of retired decision generations/review dialogue;
8. re-open D5-B2 from the canonical reconciliation + D4-R1 + D5-B1, not from the stale candidate;
9. implementation remains blocked until D9.

**This file is non-authoritative until independent review + GPT adjudication + explicit operator ratification.**