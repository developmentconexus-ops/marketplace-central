# D1 — Domains / Boundaries

> **Status:** CLOSED / ACCEPTED AS A WHOLE  
> **Program:** Architecture Rebaseline / Technical System Design  
> **Authority:** operator-approved D1 decisions, adjudicated independent-review amendments and final global-review correction consolidated here  
> **Parent authority:** `D0-PRODUCT-SYSTEM-DEFINITION.md`

## 1. Purpose and scope

D1 answers:

> **Which real business domains exist in Marketplace Central, which responsibility/state does each own, and which semantic boundaries must exist between them?**

D1 does **not** choose canonical IDs/schema (D2), communication mechanisms/events (D3), provider/ERP transport (D4), HTTP contracts (D5), frontend topology (D6), or workers/queues/transactions/deployment (D7).

Current code, schemas, APIs, tests and historical modules are evidence only. A current package does not earn target authority by existing.

## 2. Governing invariant

> **Every material Product 1.0 business responsibility has one semantic authority. Shared correctness/runtime mechanisms do not acquire business authority. Boundaries follow independent meaning, lifecycle and decisions—not legacy modules, provider nouns or deployment topology.**

> **Mechanism ≠ Authority.** Reusable mechanics may be centralized when this removes accidental complexity, but business meaning remains with the domain that understands and owns it.

Twelve business boundaries do **not** imply twelve services, databases, processes or deployments. D7 decides runtime materialization.

## 3. Business boundary catalog

| Boundary | Owns | Explicitly does not own |
|---|---|---|
| **Marketplace Portfolio** | Organization participation in marketplaces through Marketplace Installations; installation membership/lifecycle; organization-facing operational configuration; installation↔eligible Selling Entity participation/configuration; installation-level posture/attention derived from external evidence | Provider credentials/auth/protocol; technical integration registry; offer/sale state; transaction-specific Selling Entity attribution; other domains' policies |
| **Product & Channel Readiness** | Product↔channel correspondence; requirements; supported/missing/conflicting/readiness conclusion | Product master; listing lifecycle; price; sellable availability; provider-native product/catalog ontology |
| **Marketplace Offering Operations** | Marketplace offer/listing representation and lifecycle; listing intent; **price intent/mutation/convergence** | Sellable Availability; native inventory; economic calculation; provider protocol |
| **Availability Control** | Inventory Source/Scope semantics; allocation policy; Sellable Availability; **availability intent/synchronization/convergence** | Native stock/reservation truth; listing/price lifecycle; provider protocol |
| **Market Intelligence** | External comparable-market observations; comparability; competitive position/change; market-evidence sufficiency | Pricing/profitability authority; listing mutation; provider protocol |
| **Commercial Economics** | Cost Basis; economic interpretation; pricing analysis; L0 Expected Economics; L1 Order Economics; L2 settlement/Realized Economics; variance/calibration | Listing/provider state; price actuation; bank/treasury operations; provider protocol |
| **Controlled Action Governance** | Authorization delegation/grant semantics (who/what may authorize which action class/scope); authorization decision/context; authorized target-scope snapshot; authority context; decision state and provenance/correlation | Business thresholds that make an action permitted/approval-required/prohibited; other domains' policy semantics; readiness/economic/operational validity; business intent; intended target selection; execution; provider protocol |
| **Marketplace Sales** | MPC interpretation/context/correlation of marketplace-originated sales; transaction-specific Selling Entity attribution; canonical sale meaning sufficient to feed downstream lifecycle | Business Order Intent/materialization; fulfillment; economics; post-sale resolution; universal provider `Order/Pack/Shipment` model |
| **Business-System Materialization** | Business Order Intent and native-order materialization/correlation; invoicing readiness, Invoicing Intent and fiscal materialization/correlation | ERP-native TOP/document taxonomy as MPC semantics; physical fulfillment; general fiscal-policy domain |
| **Fulfillment Lifecycle** | Physical fulfillment readiness/execution; Fulfillment Node eligibility/selection; separation/conference/packing/dispatch; provider-requirement closure for fulfillment paths MPC claims to control; relevant shipment/delivery observation through terminal outcome | Company-wide WMS/TMS; fiscal authority; post-sale consequence resolution; provider-native requirement/artifact truth or provider protocol |
| **Post-Sale Resolution** | Coordination/correlation/closure of material cancellation/return/refund consequences across owning domains | Native refund/provider protocol itself; physical fulfillment semantics; economic calculation; ERP/fiscal semantics |
| **Operational Work** | Lifecycle of material actionable work: role responsibility, optional assignment, escalation, work state and evidence-backed resolution/closure | Definition of originating business truth, policy, authorization, deadline breach, economic correctness or reconciliation semantics |

## 4. Important non-domains

The following are real concepts/properties but **not independent D1 business domains**.

### 4.1 Organization / Selling Entity / Product master

- `Organization` and `Selling Entity` remain distinct real identities; D2 defines their exact identity representations and relationships.
- Marketplace Portfolio owns the operational configuration linking a Marketplace Installation to eligible Selling Entity participation where material.
- Marketplace Sales owns transaction-specific Selling Entity attribution. Materialization, Economics and other downstream domains consume that attribution rather than independently inferring `Organization = Installation = Selling Entity`.
- Product master remains authoritative in the business system/source. MPC references source-qualified product identity; no PIM/MDM/Product Master domain is justified now.

### 4.2 Policy and action disposition

Business-policy meaning belongs to the domain governed by that policy:

- availability allocation / bounded availability automation → Availability Control;
- margin floors / price boundaries / economic thresholds → Commercial Economics;
- offer/listing constraints → Marketplace Offering Operations;
- internal operational targets → responsible operational domain.

The **action-owning domain** owns the effective action disposition for its action (`permitted`, `approval-required`, `prohibited`, `automation-eligible`) using all governing policies/evidence/validity relevant to that action.

Controlled Action Governance owns only authorization-specific semantics: who/what may authorize which action class/scope and the resulting authorization decision/context. It is not a generic business-policy engine.

Shared default/inheritance/override/provenance mechanics may be decided later without acquiring policy meaning.

### 4.3 Evidence / provenance / freshness / coverage / audit

No `Evidence`, `Freshness`, `Coverage`, `Lineage` or `Audit` business domain exists.

Shared knowledge/provenance/time primitives and audit mechanics may exist later. The consuming business domain remains authority for whether evidence is sufficiently known, current, complete and applicable for its decision. Audit/history composes domain-owned facts; it is not a second business truth.

### 4.4 Time-bound obligations / SLA

No generic SLA/Obligation domain exists.

The domain responsible for satisfying an obligation owns completion semantics, internal target policy, urgency and breach interpretation. External deadlines/windows remain externally authoritative evidence. Operational Work owns escalation/work lifecycle only after actionable work exists. Shared clocks/timers/schedulers are D7 mechanics.

### 4.5 External-action execution safety and target scopes

No `Mutation`/`Action` business domain exists.

The action-owning domain retains business intent, intended target selection, validity and convergence/reconciliation meaning. Governance owns authorization when required. Adapters own external protocol. Reusable attempt/correlation/idempotency/ambiguity/retry-safety mechanics may be centralized in D3/D7 without acquiring business authority.

For material multi-target actions, three scopes remain distinguishable:

1. **Intended target scope** — defined by the action-owning domain as part of business intent.
2. **Authorized target scope** — preserved by Governance as authorization context, derived/constrained from the intended scope rather than becoming a second intent authority.
3. **Attempted/outcome scope** — preserved from actual execution evidence by execution-safety/runtime mechanics, including confirmed/rejected/ambiguous/not-executed distinctions when material.

Blind replay of a potentially accepted external write remains prohibited unless safe replay is established.

### 4.6 Integration capability and provider-requirement closure

Use three levels:

1. **Integration Support / Descriptor** — what an MPC adapter implements; technical D4 concern.
2. **Provider Effective Capability / Requirement** — what provider/installation/resource/operation actually allows/requires; provider-authoritative evidence translated by D4.
3. **Effective Business Capability** — whether MPC can/should perform this action now; owned by the consuming business domain using support + provider evidence + domain state/policy/validity.

No universal marketplace business interface owns all three.

For fulfillment paths MPC claims to control, **Fulfillment Lifecycle owns provider-requirement closure semantics**: from provider-effective requirement evidence, it determines which prerequisites/data/artifacts/acknowledgements are material to that path, orchestrates/reconciles their satisfaction and owns the MPC conclusion that the path is sufficiently provider-ready/closed. Provider-native requirements, states and artifacts remain provider-authoritative; D4 owns protocol translation. If an accepted integration cannot satisfy a required operation, the path is explicit `unsupported` / `external-required` rather than presented as fully MPC-controlled.

## 5. Semantic authority edges

An allowed edge means semantic dependency is legitimate. It does **not** decide sync call, event, projection, queue or orchestration; D3 chooses the mechanism.

> **A consumer may use another context's owned meaning only through an explicit public semantic boundary owned by the producer. Consumption never transfers authority.**

Accepted baseline:

- **Marketplace Portfolio → marketplace-facing domains:** installation participation/configuration/posture and eligible Selling Entity participation as required; Organization, Selling Entity and provider account never collapse by convenience.
- **Readiness → Offering / Availability:** channel correspondence/readiness may govern whether a product can be operated; consumers do not silently recompute readiness.
- **Offering → Availability:** Availability may consume the marketplace representation/target needed to synchronize availability; Offering does not compute Sellable Availability.
- **Offering → Market Intelligence:** Market Intelligence may consume the organization's own offer state for competitive comparison while retaining comparability authority.
- **Offering → Commercial Economics:** Economics may consume offer/listing representation needed for offer-specific economic interpretation, including marketplace listing category/type/current commercial state when material. Provider fee/rule evidence remains externally authoritative and enters through D4. ERP taxonomy must never substitute for marketplace listing category merely because it is locally available.
- **Market Intelligence → Commercial Economics:** comparable-market evidence/interpretation feeds economic reasoning; Economics does not independently reinterpret provider competitor payloads when Market Intelligence owns that meaning.
- **Commercial Economics → Offering:** economic conclusions/candidate price implications may inform price intent; Economics does not write marketplace price.
- **Marketplace Sales → Materialization / Fulfillment / Economics / Post-Sale:** one canonical sale interpretation, including transaction-specific Selling Entity attribution when material, feeds downstream responsibilities; downstream domains do not each invent provider transaction semantics or entity attribution.
- **Materialization ⇄ Fulfillment:** materialized business/fiscal outcomes may gate fulfillment; physical readiness/conference evidence may gate Invoicing Intent. This business workflow cycle must not become a private-code dependency cycle.
- **Materialization → Commercial Economics:** attributable order/fiscal results may contribute economic evidence without transferring fiscal authority.
- **Post-Sale Resolution ⇄ Sales / Materialization / Fulfillment / Economics:** Post-Sale coordinates consequences and observes closure; each owning domain executes/interprets its own semantics.
- **Controlled Action Governance ⇄ action-owning domains:** the domain supplies intended action/scope and effective disposition; Governance applies authorization-specific delegation/grant semantics and returns the authorization decision/context. Governance preserves authorized scope as a snapshot/constrained form of the domain-owned intent; execution-time domain validity remains with the action owner.
- **Operational Work ⇄ domains with actionable work:** source domain defines the issue and semantic closure requirement; Work owns responsibility/assignment/escalation/work state; source domain determines whether resolution evidence actually closes the originating condition.

External-system evidence enters through D4 adapters/consumer-owned ports without transferring semantic authority. In particular, provider requirement/artifact evidence feeds Fulfillment Lifecycle, which owns the business closure conclusion for its claimed path.

A provider API that combines multiple fields/actions does not merge their business authorities.

## 6. Forbidden boundary violations

Unless D1 is explicitly reopened, the target must not introduce:

- cross-context SQL/private-table access as unnamed communication;
- import of another context's private implementation;
- shared mutable business entities with ambiguous ownership;
- reimplementation of another context's owned business rule/meaning;
- consumer mutation of producer-owned state;
- provider DTO/protocol types leaking across business contexts;
- read projections/views becoming write authorities;
- Governance becoming a generic rules engine, business-threshold owner or domain-validity authority;
- Operational Work becoming a reconciliation/business-truth authority;
- Marketplace Portfolio becoming the technical integration registry;
- downstream domains independently inferring Selling Entity transaction attribution;
- Commercial Economics performing marketplace price writes;
- Offering owning Sellable Availability semantics;
- Marketplace Sales absorbing ERP materialization, fulfillment, economics or post-sale semantics;
- D4/integration code becoming authority for fulfillment provider-requirement closure.

D3 may choose communication mechanisms **only inside the D1 semantic edge set**. If D3 discovers a genuinely necessary new semantic dependency, D1 must be reopened rather than hiding that dependency in an event, API, queue, database or projection.

## 7. Legacy semantic disposition

Legacy modules have no right to survive as target authorities by name/history. Valid responsibilities move to the D1 authority that owns their meaning; technical mechanics move to later technical layers; obsolete/authority-less abstractions are retired.

| Current module/context | D1 disposition |
|---|---|
| `catalog` / `contexts/catalog` | **RETIRE / SPLIT** — no MPC Product Master authority; useful primitives/contracts may be reassigned where justified |
| `channelfees` | **ABSORB / SPLIT** — economic meaning → Commercial Economics; provider acquisition → D4 |
| `classifications` | **RETIRE** — no independent Product 1.0 authority proven |
| `connectors` | **TECHNICAL / SPLIT** — D4 adapters + consumer-owned ports |
| `dashboard` | **TECHNICAL** — D6/views projection, not authority |
| `divergences` | **SPLIT** — divergence/reconciliation meaning → originating domain; actionable lifecycle → Operational Work |
| `erp_import` | **TECHNICAL** — D4 acquisition mechanism if still needed |
| `integrations` | **SPLIT** — Marketplace Installation business meaning → Portfolio; credentials/auth/descriptor → D4; operation/runtime mechanics → D7 |
| `internal_read` | **SPLIT / TECHNICAL** — business-system adapters → D4; consumer contracts distributed to Readiness/Availability/Economics/etc. |
| `inventory` | **NUCLEUS** — target semantic authority is Availability Control |
| `listings` / `contexts/listings` | **NUCLEUS + SPLIT** — target semantic authority is Offering; availability/provider-native/projection concerns leave it |
| `market` | **NUCLEUS** — target semantic authority is Market Intelligence |
| `marketplaces` | **SPLIT** — Portfolio + D4 technical concerns + domain-local policy/economic meaning where applicable |
| `mutations` | **TECHNICAL / SPLIT** — intent returns to action owners; authorization → Governance; reusable execution safety → D3/D7; provider protocol → D4 |
| `orders` | **SPLIT** — Marketplace Sales + Business-System Materialization + Fulfillment + other owning domains where responsibility actually belongs |
| `pricing` | **ABSORB / SPLIT** — economic calculation → Commercial Economics; price actuation → Offering |
| `product_links` | **NUCLEUS** — target semantic authority is Product & Channel Readiness; exact identity model deferred to D2 |
| `profitability` | **ABSORB** — Commercial Economics, preserving L0/L1/L2 distinction |
| `sourcekind` | **TECHNICAL PRIMITIVE** — D2/kernel or D4 only if genuinely universal |
| `sync` | **TECHNICAL / SPLIT** — domain-specific completeness/freshness/convergence stay local; polling/cursor/scheduling/runtime → D4/D7 |
| `tenant_config` | **SPLIT / TECHNICAL** — identity/isolation → D2; business config → owning domains; source/integration config → D4 |

If a proven Product 1.0 responsibility cannot fit an accepted boundary without distorting its authority, **reopen D1**; never force-fit it into the nearest context.

## 8. Explicit defers

D1 intentionally does not decide:

- **Composite offers/composition:** while the accepted first Mercado Livre flow does not require component-dependent composition, no composition domain/model is created. If the selected flow requires composition that materially changes availability, offer representation or business-order materialization, D1 reopens to assign semantic ownership before D2/D4 model or integrate it. Real composition may never be silently flattened.
- **D2:** canonical IDs, Organization/Selling Entity/source-qualified product identity, persistence ownership, tenant/isolation model, exact money/evidence/time/value primitives;
- **D3:** synchronous/event/projection matrix, event contracts, outbox/communication implementation, workflow-cycle realization;
- **D4:** exact Mercado Livre/Sankhya/payment contracts, adapter DTOs, provider capability/requirement mappings, auth/credentials, pagination/webhook/polling/source completeness evidence;
- **D5:** HTTP/OpenAPI/bulk/idempotency/precondition contracts;
- **D6:** screens, portfolio views, work inbox, approvals/history/attention UX and projection topology;
- **D7:** process topology, schedulers, workers, queues, retries, locks/versioning, transactions, execution-safety implementation, observability/deployment.

No later stage may silently overturn D1 authority. New material evidence may explicitly reopen it.

## 9. Reopen triggers

Reopen a D1 decision only when new evidence materially changes an assumption, for example:

- a proven Product 1.0 responsibility has no honest owner;
- two accepted contexts demonstrably need to own the same business meaning;
- two accepted boundaries are shown by real accepted flows to have one business meaning/lifecycle with no independently owned decision, or one boundary proves to exercise no independent decision authority;
- a currently absorbed concept gains an independent business lifecycle/decision authority;
- a second real source/master requires identity/mastering behavior not supported by the current boundary;
- MPC becomes actual TMS/WMS/fiscal/settlement/obligation-management authority rather than orchestration/observation;
- a new provider/business-system requirement cannot be represented without dismantling an accepted authority boundary;
- the selected accepted flow requires material composite semantics whose owner is not already unambiguous;
- D2–D7 discover a genuinely necessary semantic dependency not permitted here.

Do not reopen for naming preference, framework fashion or hypothetical future capability alone.

## 10. Decision lineage

The operator-approved D1 decisions consolidated here are:

- **D1.1:** Commercial Intelligence/Economics separate from operational actuation.
- **D1.2:** Readiness, Offering and Availability are distinct responsibilities.
- **D1.3:** policy meaning is domain-local; no central Policy domain.
- **D1.4:** narrow Controlled Action Governance boundary.
- **D1.5:** Business-System Materialization separate from Fulfillment; no separate Fiscal domain now.
- **D1.6:** Fulfillment includes relevant shipment/delivery observation; Post-Sale Resolution separate; no Shipment/Delivery domain now.
- **D1.7:** Operational Work independent; reconciliation semantics domain-owned.
- **D1.8:** Market Intelligence separate from Commercial Economics.
- **D1.9:** Marketplace Sales independent.
- **D1.10:** Marketplace Portfolio independent; Organization/Selling Entity are identities, not separate domains.
- **D1.11:** no MPC Product/Catalog master domain now.
- **D1.12:** three-level capability model: integration support, provider-effective capability, effective business capability.
- **D1.13:** evidence/provenance/freshness/coverage/audit are cross-cutting semantics/mechanics, not domains.
- **D1.14:** time-bound obligation meaning is domain-owned; timers/schedulers are mechanics.
- **D1.15:** external-action execution safety is shared mechanism, not business authority.
- **D1.16:** final catalog ratified at 12 business boundaries.
- **D1.17:** legacy semantic disposition approved.
- **D1.18:** semantic authority edges and forbidden boundary violations approved.
- **D1.G1:** final global review assigns provider-requirement closure for claimed fulfillment paths to Fulfillment Lifecycle while provider-native truth/protocol remain external/D4 concerns.

## 11. Independent review adjudication

The independent adversarial D1 review returned **AMEND**, not `REOPEN`: no boundary was added, removed, merged or split.

The operator approved adjudication of six amendments:

1. add **Offering → Commercial Economics** for offer-specific economic inputs and explicitly block ERP taxonomy from substituting for marketplace listing category;
2. distinguish business-policy thresholds/effective action disposition from Governance's authorization-specific delegation/grant semantics;
3. assign installation↔Selling Entity participation configuration to Portfolio and transaction-specific Selling Entity attribution to Marketplace Sales;
4. distinguish intended, authorized and attempted/outcome target scopes;
5. record composite-offer semantics as an explicit defer/reopen condition rather than a silent gap;
6. add symmetric reopen triggers for excessive fragmentation/vestigial boundaries.

These amendments clarify ownership and legal semantic edges while preserving the 12-boundary structure.

## 12. Final global review

The final internal **Global Coherence + YAGNI / Overengineering / Future-Cost** review re-tested capability coverage, duplicate/missing authority, God Context risk, excessive fragmentation, cross-boundary cycles, legacy disposition, D2–D7 leakage and foreseeable evolution.

It found one final material gap: provider-requirement closure for fulfillment paths was required by D0 but only implicit in the D1 contract. The operator approved the minimal D1.G1 correction recorded above.

After D1.G1, the review found no remaining material contradiction:

- all Product 1.0 responsibilities have an honest semantic owner or explicit external authority/defer;
- no duplicate business authority remains unresolved;
- no additional business domain/framework is justified;
- the 12 boundaries remain independently meaningful without implying runtime fragmentation;
- legacy modules remain evidence rather than target authority;
- D2–D7 retain their intended mechanics/design responsibilities;
- multiple marketplaces, another participating business system and future multi-Organization operation remain additive seams rather than reasons to prebuild speculative frameworks.

## 13. Closure

D1 is **CLOSED / ACCEPTED AS A WHOLE**.

The accepted target contains exactly these 12 business boundaries at D1 altitude:

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

No Product 1.0 semantic remains for D2–D7 to invent at domain/boundary altitude. New material evidence may reopen D1 only through the triggers above.

**Exact next stage:** D2 — Identity / Tenant / Data Ownership.

Product implementation remains blocked until D9 is accepted.

---

## 14. Post-closure accepted amendment — Personal Notifications + AuthorizationRequest boundary

**Accepted / operator-ratified in D6-R2 and consolidated here on 2026-08-25.** The twelve D1 **business boundaries remain exactly twelve**. This amendment adds one bounded supporting semantic owner and clarifies Governance identity/edges; it does not create a 13th business domain or a generic workflow platform.

### 14.1 Personal Notifications is a supporting semantic owner, not a business domain

`Personal Notifications` owns only durable personal-awareness meaning:

- one human recipient's Notification lifecycle (`unread/read`, `active/archived`, bounded supersession);
- bounded Organization routing configuration for the ten accepted `ORG_ROUTED` families;
- personal-awareness deduplication/suppression and immutable presentation/result snapshots needed to explain the occurrence.

It explicitly does **not** own source business truth, action disposition, Work, authorization, sales/materialization/fulfillment/post-sale semantics, source mutation, general subscriptions/preferences, delivery channels or provider protocol.

The accepted source-owner → Personal Notifications awareness edges are:

```text
Marketplace Portfolio
Marketplace Offering Operations
Availability Control
Commercial Economics
Marketplace Sales
Business-System Materialization
Fulfillment Lifecycle
Post-Sale Resolution
Operational Work
Controlled Action Governance
  → Personal Notifications
```

These edges communicate source-owner committed meaning for awareness only. They never transfer the source owner's business authority.

### 14.2 Accepted family/audience boundary

The fourteen accepted families remain the D0 set. Their audience strategy is closed at this altitude:

```text
DIRECT_SOURCE
  OFFERING_ASYNC_ACTION_RESULT
  WORK_ASSIGNMENT
  AUTHORIZATION_DECISION_RESULT

OWNER_DERIVED
  AUTHORIZATION_ACTION_REQUIRED

ORG_ROUTED
  MARKETPLACE_INSTALLATION_ATTENTION
  AVAILABILITY_ATTENTION
  ECONOMIC_RECONCILIATION_ATTENTION
  NEW_MARKETPLACE_SALE
  SALE_ATTENTION
  MATERIALIZATION_ATTENTION
  FULFILLMENT_ACTIONABLE
  FULFILLMENT_ATTENTION
  SHIPMENT_EXCEPTION
  POST_SALE_ATTENTION
```

Exact family identity/state belongs D2/D5. D1 freezes only ownership/audience semantics. Permission/role membership never silently selects routing recipients.

Two overlap rules are deliberately bounded rather than generalized:

1. an exact source attention occurrence may yield to the richer `WORK_ASSIGNMENT` awareness for the same human when that Work/assignment causally replaces the source alert;
2. a `SALE_ATTENTION` occurrence may yield to the richer `POST_SALE_ATTENTION` continuation for the same exact Post-Sale Resolution/human.

No generic notification relation graph, routing DSL or cross-kind precedence engine is admitted.

### 14.3 AuthorizationRequest remains Controlled Action Governance meaning

A materially distinct pre-decision authorization episode is canonical `AuthorizationRequest` state owned by **Controlled Action Governance**. It is separate from:

```text
action-owning Business Intent
AuthorizationDecision
Operational Work
Notification
execution outcome
```

The closed governed target families remain ListingIntent, PriceIntent, BusinessOrderIntent and InvoicingIntent. `AuthorizationRequest` provides stable identity/history for one authorization episode; it does not become a generic approval/workflow/case owner.

`AUTHORIZATION_ACTION_REQUIRED` awareness references the exact pending AuthorizationRequest. `AUTHORIZATION_DECISION_RESULT` remains requester-oriented toward the governed target rather than forcing Governance-history read authority.

A zero-current-decider condition remains explicit Governance/Work coordination; it does not convert Work into an approver or authorization authority.

These accepted amendments extend the legal D1 edge set only as stated above. D3 owns communication/recovery; D5 owns Product wire; D6 owns human realization; D7 owns runtime mechanics.
