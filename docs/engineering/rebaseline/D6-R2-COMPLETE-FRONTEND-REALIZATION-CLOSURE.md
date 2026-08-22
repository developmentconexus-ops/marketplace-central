# D6-R2 — Complete Frontend Realization Closure

> **Status:** OPEN / ACTIVE — P0–P3 DERIVED; no operator UX `LOCKED` decision has been made in this stage
> **Program:** Architecture Rebaseline / Technical System Design
> **Parent authority:** accepted D0–D8, accepted D6 frontend authority, canonical Product OAD
> **Execution method:** [`Frontend Product Experience Planning Method v2.1`](../../development/frontend-product-experience-planning-method.md), tailored as recorded here
> **Product implementation:** BLOCKED UNTIL accepted D9

## 1. Purpose and boundary

D6-R2 closes the gap between accepted frontend architecture and an implementation-ready human Product experience.

It does **not** reopen accepted D6 information architecture, route meanings, React/TypeScript, TanStack Query, TanStack Router, `openapi-typescript`, `openapi-fetch`, source topology or dependency direction by preference. It does not implement Product code and does not begin the Pre-D9 Implementation Readiness Contract or D9.

The stage must make a future implementation agent able to realize the frontend without inventing material product, interaction, authority, state or wire decisions while coding.

The required sequence is:

```text
P0–P3  global foundation first
P4     revalidate the accepted IA; redesign only on a material falsifier
P5     derive the complete screen/surface inventory from flows + accepted IA
P6/P7  only when the method's ambiguity triggers are real
P8     render structural wireframes block by block; only the operator may mark a block LOCKED
P9     bind each material screen/action to exact owner + operationId + Permission + H/A/S + identity + state + effect semantics
P10–P13 complete interaction/prototype/adversarial/conformance obligations as triggered by the method
P14    close frontend realization readiness
```

No screen is drawn in P0–P3.

## 2. D6-R2 tailoring laws

1. Accepted D6 IA/routes/topology/technologies are **input authority**, not fresh design hypotheses.
2. P4 is a falsification/revalidation pass over accepted IA. Preference, visual fashion or backend package shape is not a reopen reason.
3. Existing D6-B1 screen states and prototype are evidence of accepted interaction coverage, not automatic P8 `LOCKED` wireframes.
4. P6 reference study and P7 competing hypotheses activate only for real material ambiguity.
5. P8 progresses material blocks one by one. The assistant/reviewer may propose `CANDIDATE`; only the operator may set `LOCKED`.
6. P9 uses the current canonical OAD. No hand-written frontend DTO or screen-shaped Product operation may appear.
7. Product knowledge/outcome semantics remain explicit: unknown is not empty, ambiguous is not failure, stale write is not rejection, external mechanism success is not Product success.
8. Technical Ingress remains outside the Product operation census.
9. Product implementation remains blocked throughout D6-R2.

---

# P0 — Recover accepted authority

## 3. P0 authority pack and recovered truths

The bounded recovery uses accepted composition authorities instead of recursively reopening D0–D8 history.

| Authority | Recovered binding truth for D6-R2 |
| --- | --- |
| D0 Product/System Definition | Product 1.0 is Marketplace Operations + Commercial Intelligence; four accepted human actor classes; MPC is control plane, not marketplace dashboard/ERP/WMS/TMS/CRM replacement. |
| D6 Frontend | Accepted shell/IA, route grammar, Organization/Installation context laws, frontend state ownership, React/TypeScript + TanStack Query/Router + generated OAD transport profile, topology/import laws. |
| D6-B1 Interaction Map | 99/99 Product operations already have coherent interaction homes across 39 derived route/screen states; these homes are coverage evidence, not P8 locks. |
| Canonical Product OAD | OpenAPI 3.1.2; 99 Product operations; 30 ordinary Permissions; Principal kinds H/A/S only; human browser uses server-side session + CSRF; Product implementation remains blocked until D9. |
| D8 Golden Flows / proof closure | Representative system behavior is 3 business golden flows + SR-01; P1/P3 converged, P2 is operator-ratified redefer, P5 narrows full alternate destination to external-required/unsupported, P4/P6 were not triggered; no D0–D7 reopen. |
| Frontend Product Experience Planning Method v2.1 | Human needs before screens, coverage before layout, operator-only `LOCKED`, no screen-shaped API, P0–P14 readiness discipline. |

### 3.1 Accepted frontend architecture that D6-R2 must preserve

```text
Marketplace Central
Organization: <display_name>

VISÃO GERAL

OPERAÇÕES
  Preparação
  Publicações
  Disponibilidade
  Vendas
  Expedição
  Pós-venda

ESTRATÉGIA E INTELIGÊNCIA
  Performance
    Resumo
    Publicações
    Mídia
  Mercado
  Economia

CONTROLE
  Trabalho
  Aprovações

CONFIGURAÇÕES
```

Accepted primary route identities remain:

```text
/org/:organizationId/visao-geral
/org/:organizationId/preparacao
/org/:organizationId/publicacoes/*
/org/:organizationId/disponibilidade
/org/:organizationId/performance/*
/org/:organizationId/mercado
/org/:organizationId/economia/*
/org/:organizationId/vendas/*
/org/:organizationId/expedicao/*
/org/:organizationId/pos-venda/*
/org/:organizationId/trabalho/*
/org/:organizationId/aprovacoes/*
/org/:organizationId/configuracoes/*
```

Accepted realization direction remains:

```text
app/routes
    ↓
features/<lens-or-flow>
    ↓
api/<owner-family>
    ↓
api/transport
    ↓
api/generated

features ──→ ui
app/routes ─→ ui
```

### 3.2 P0 invariants carried into every later phase

- `organization_id` is canonical scope; `display_name` is presentation only.
- Marketplace Installation is explicit navigation/request context where required; never an ambient/default authority.
- Source-qualified external identities remain qualified by their accepted source/Installation/native key.
- H browser access uses server-side session + CSRF; ordinary browser code never owns OIDC access/refresh tokens.
- Permission-conditioned visibility is usability only; server authorization is authoritative.
- TanStack Query owns server state; URL owns navigation state; form draft and ephemeral UI do not become Product truth.
- no universal normalized entity store, generic action bus, generic Strategy/Analytics owner, BFF or screen-shaped API.
- no automatic replay of consequential mutations; idempotency and precondition semantics survive retries.
- `complete / partial / unknown / unavailable / unsupported`, known-zero/known-empty and exposed freshness remain distinguishable.
- Performance preserves exact Installation, explicit periods, comparison admissibility, provider measurement basis and evidence provenance.
- Strategy/Overview/Sale detail are composition only and acquire no write authority.
- D8 P2 remains a future real-proof obligation and may not disappear from implementation-readiness proof planning.
- D8 P5 forbids presenting a full alternate street/fiscal destination override as supported Product capability; the sanctioned contact reference is not that capability.

### 3.3 P0 exit

**DERIVED.** No known frontend requirement in the accepted Product surface is currently unowned. No material contradiction requires a D0–D8 reopen. Unknown user-behavior evidence is registered below rather than converted into design authority.

---

# P1 — Actors, jobs and user needs

## 4. Accepted actor contexts

D6-R2 uses D0 actor classes as role contexts, **not invented personas**.

| Actor context | Trigger / situation | User need / job | Desired outcome | Information needed | Handoffs / constraints |
| --- | --- | --- | --- | --- | --- |
| Marketplace Operations Operator | A marketplace account/product/listing/sale requires normal operation or exception handling | Prepare channel participation, control listings/price intent within policy, understand performance/market/economics, inspect sales and resolve operational exceptions | Marketplace participation is convergent, explainable and operable without bypassing policy | Readiness, requirements, Listing/Intent state, availability, market/economic evidence, performance, sale/work state | Hands physical execution to Fulfillment; hands exceptional authority to Governance/manager; cannot redefine policy to make an action permissible |
| Fulfillment / Dispatch Operator | A marketplace sale reaches internally operated physical execution | Know what is ready, separate/confer correctly, trigger/observe the accepted fiscal/provider prerequisites, pack and hand off dispatch safely | Correct physical goods reach verified dispatch without premature invoicing or fabricated readiness | Sale/execution identity, separation/conference state, fiscal/provider readiness, artifacts, shipment state, exception work | Physical inconsistency blocks the normal path and becomes explicit Work; UI cannot self-authorize physical checkpoints |
| Commercial / Marketplace Manager | Commercial evidence, policy thresholds or exceptions require a decision | Compare market/performance/economics, evaluate price trade-offs, manage legitimate MPC-owned policy and resolve bounded commercial/economic ambiguity | Commercial actions remain profitable/explainable and inside governing policy | Performance, market evidence, expected/realized economics, effective policy/provenance, approvals/exceptions | Does not edit externally governed rules; does not become integration/security administrator by default |
| Owner / Administrator / Policy Approver | Access, integrations, exceptional authority, delegation or containment requires governance | Establish valid access/channel configuration and approve exceptional actions above delegated authority | Organization remains safely operable with explicit accountability and least necessary authority | Access context, members/roles, Installations, authorization decisions/delegations, configuration state | Not a routine approval bottleneck; cannot configure away audit/reconciliation/safety invariants |

## 5. Need statements independent of screens

### N01 — Work in the correct Organization

When a human enters MPC or changes workspace, they need to know which Organization and access context are current, so that every subsequent observation/action is scoped correctly.

### N02 — Establish a usable marketplace channel

When the organization needs to operate on a supported marketplace, an authorized human needs to create/configure the Marketplace Installation and complete any required external authorization ceremony, so that later Product operations have an explicit valid channel context.

### N03 — Determine whether a source product is publishable

When a product should participate in a marketplace, the operator needs to discover the exact source product and understand readiness/requirements/correspondence, so that missing or conflicting conditions are resolved before publication intent is submitted.

### N04 — Express controlled listing and price intent

When a listing or price decision is ready to be proposed, the operator/manager needs to author and submit accepted intents with current revision/evidence, so that the Product can execute/reconcile them without the frontend inventing provider operations.

### N05 — Understand how marketplace participation is performing

When reviewing channel results, the operator/manager needs trustworthy account-, listing- and media-scope performance evidence with explicit periods/coverage/comparability, so that action is based on what is actually known.

### N06 — Understand competitive and economic position

When considering commercial action, the operator/manager needs market evidence, expected/realized economics and stateless scenario evaluation, so that a PriceIntent or policy decision is explainable rather than guessed from a dashboard number.

### N07 — Know and govern sellable availability

When marketplace availability is operating or its policy/source configuration changes, authorized humans need to inspect sellable availability and manage only accepted sources/policy, so that automatic synchronization remains policy-owned and uncertainty stays explicit.

### N08 — Understand a sale across system boundaries

When a marketplace sale arrives or diverges, operations need to see sale identity plus business-order, Party/destination, invoicing, fulfillment, post-sale and economics evidence without collapsing owners, so that the next safe action/handoff is clear.

### N09 — Execute physical fulfillment safely

When an internal fulfillment execution is actionable, the fulfillment operator needs to record accepted physical checkpoints and obtain required artifacts/readiness, so that dispatch occurs only from verified physical/fiscal/provider state.

### N10 — Observe shipment/delivery outcome

After dispatch, operations need to observe source-qualified shipment state and material exceptions, so that delivery divergence becomes explicit work rather than assumed success.

### N11 — Coordinate essential post-sale consequences

When cancellation/return/refund consequences require MPC coordination, operations need to create/inspect the accepted PostSaleResolution without inventing a generic close/refund/cancel command, so that source owners remain authoritative.

### N12 — Resolve operational work without turning Work into source truth

When ambiguity, divergence or failure creates actionable Work, an operator needs to assign/clear/hold/resume/escalate it, so that coordination progresses while the underlying owner remains the truth source.

### N13 — Approve exceptional action and delegate authority safely

When an action requires governance or delegated authority changes, authorized humans need to inspect/create decisions and manage delegations, so that permission, governance decision and target action remain separate.

### N14 — Administer access safely

When Organization membership/access changes, an administrator needs to inspect members/roles and assign/revoke accepted AccessRoles, so that human access is explainable without turning IdP/provider roles into Product roles.

### N15 — Configure fulfillment operating context

When an Organization changes fulfillment nodes or internal operating targets, an authorized human needs to manage the accepted Fulfillment configuration, so that execution has explicit eligible context without a generic WMS model.

### N16 — Govern commercial policy and reconcile economic attribution

When MPC-owned commercial policy or economic attribution requires intervention, an authorized manager needs to inspect/update the accepted policy or resolve bounded attribution ambiguity, so that realized economics remain explainable and current.

## 6. Assumption register

These are **not** accepted Product facts. They may influence later structural candidates only while visibly registered and must be probed before P14 closes.

| ID | Assumption | Evidence level | Phases/blocks influenced | Required P12 probe | Status |
| --- | --- | --- | --- | --- | --- |
| A01 | Relative frequency and information density of the 16 needs are not quantitatively evidenced yet. | accepted product/domain authority proves the jobs, not frequency | P4/P5/P8 density, default landing and prioritization | Observe/operator-validate typical daily/weekly task mix and high-frequency queues | OPEN |
| A02 | Fulfillment users may require stronger tablet/mobile ergonomics than strategy/configuration users, but the real device/work-floor distribution is not evidenced. | domain inference only | P8 responsive structure for fulfillment/work | Validate real devices, gloves/scanning/standing constraints and viewport use | OPEN |
| A03 | Existing Portuguese IA labels are accepted, but detailed field/action terminology has not been user-tested. | operator-ratified IA, no direct terminology test | P4 glossary, P8 labels, P12 comprehension | Run terminology/comprehension walkthrough on material objects/actions | OPEN |
| A04 | The small Overview composition is accepted, but which permitted signals deserve first visual priority is not directly evidenced. | D6 accepts the surface, not final hierarchy | P5/P8 Overview | Task-based operator walkthrough: what must be noticed first and why | OPEN |
| A05 | No bulk-selection/bulk-action UX should be assumed beyond Product operations actually admitted by the OAD. | no accepted bulk capability evidence | P5/P8 collection actions | Probe real repeated-task pain before proposing any bulk interaction | OPEN |

### 6.1 P1 exit

**DERIVED.** Four accepted actor contexts resolve to 16 outcome-oriented needs without inventing personas or pages. The missing behavioral evidence is explicit and does not authorize structural invention.

---

# P2 — End-to-end user-flow inventory

## 7. Cross-flow entry and safety grammar

Every human flow starts from authenticated H-browser context and must preserve:

```text
current access context
→ explicit Organization
→ exact/optional Installation only when the Product operation admits it
→ source-qualified identity where applicable
→ understand authoritative/knowledge state
→ decide
→ execute only admitted Product action or Technical Ingress ceremony
→ preserve consequential outcome semantics
→ hand off to the next semantic owner / Work when required
```

Deep-link and stale-access failures re-resolve server authority; hidden navigation never substitutes for authorization.

## 8. Complete flows

### UF01 — Establish current Organization/access context

**Actors:** all human actors; administration branch for Owner/Administrator.

```text
enter MPC
→ GetCurrentAccessContext
→ select/confirm accessible Organization by canonical id
→ re-scope server/navigation state
→ if administering access: inspect members + AccessRoles
→ assign/revoke AccessRole only with admitted access authority
→ outcome: valid explicit workspace or safe blocked/no-access state
```

Failure/alternate: invalid/deep-linked Organization never becomes fallback/default tenant; Organization switch clears incompatible Installation context.

### UF02 — Establish/manage Marketplace Installation

**Actors:** Owner/Administrator; Marketplace Operations Operator when delegated.

```text
inspect supported/connected channel context
→ list existing Marketplace Installations
→ create Installation when admitted
→ complete exact external OAuth Technical Ingress ceremony when required
→ re-read Installation state
→ update/deactivate through Product operations when authorized
→ inspect Selling Entity context when needed
→ outcome: explicit usable/degraded/inactive Installation state
```

Failure/alternate: no generic provider catalog, no invented `ConnectMarketplace`, no automatic advertiser binding by name/first result.

### UF03 — Prepare a source product for marketplace participation

**Actor:** Marketplace Operations Operator.

```text
enter with exact Organization + marketplace context
→ search admitted source products; SourceInstance may narrow but omission never chooses hidden default
→ select exact SourceInstance + native product key
→ inspect channel readiness + publication requirements
→ resolve/clear correspondence only when needed and authorized
→ re-read readiness
→ outcome: ready state or explicit missing/conflicting/unsupported conditions
```

### UF04 — Author and submit listing/price intent

**Actors:** Marketplace Operations Operator; Commercial Manager for bounded price decision.

```text
start from prepared subject
→ create/read/update ListingIntent draft
→ attach admitted media when needed
→ inspect contextual availability/economics
→ optionally evaluate price scenario
→ create PriceIntent separately when a price decision exists
→ submit ListingIntent with current revision/precondition
→ observe accepted/pending/rejected/ambiguous/convergence state
→ outcome: controlled Product intent with explainable current state
```

Failure/alternate: stale precondition is not business rejection; ambiguous potentially accepted effect is not blindly retried; UI never calls a provider write directly.

### UF05 — Understand Installation-level Performance

**Actors:** Marketplace Operations Operator; Commercial Manager.

```text
choose exact Marketplace Installation
→ choose explicit current period
→ optionally choose explicit comparison period
→ request Performance summary
→ inspect traffic/sales/media evidence with coverage/provenance
→ show delta only when Product says comparable
→ outcome: trustworthy performance understanding or explicit insufficient/unavailable evidence
```

### UF06 — Diagnose Listing Performance

**Actors:** Marketplace Operations Operator; Commercial Manager.

```text
choose exact Installation
→ list listing performance without dropping Listings that lack evidence
→ select source-qualified Listing
→ get Listing performance
→ optionally compose owner-local Offering/Availability/Market/Economics reads
→ outcome: diagnosis without moving write authority into Performance
```

### UF07 — Inspect Retail Media evidence without Ads management

**Actors:** Commercial Manager; Marketplace Operations Operator.

```text
choose exact Installation + period
→ list Retail Media Performance
→ preserve campaign/listing/catalog-group/family-group scope
→ preserve provider measurement basis/attribution/coverage
→ outcome: read-only media evidence or explicit unavailable/unsupported state
```

No campaign/budget/bid/targeting/creative write exists.

### UF08 — Analyze market/economics and express price decision

**Actors:** Commercial Manager; Marketplace Operations Operator within delegated policy.

```text
select explicit subject/context
→ inspect competitive position + comparable offers
→ inspect expected economics
→ evaluate stateless price scenario when useful
→ decide inside effective policy/governance
→ create PriceIntent under Offering when admitted
→ observe PriceIntent/current convergence state
→ outcome: explainable price intent or no-action/approval-needed outcome
```

Market/Economics remain evidence owners; Offering remains PriceIntent owner.

### UF09 — Observe and govern sellable availability

**Actors:** Marketplace Operations Operator for observation; Commercial Manager/authorized administrator for policy/source configuration.

```text
inspect sellable availability list/current state
→ if configuration is needed: inspect Inventory Sources + effective allocation policy
→ create/update/deactivate source or update policy only through admitted operations
→ re-read effective/sellable state
→ outcome: explicit availability and provenance/policy state
```

No manual `SetAvailableQuantity` or generic `SyncAvailability` action is invented.

### UF10 — Understand and reconcile a marketplace Sale through materialization

**Actors:** Marketplace Operations Operator; Commercial Manager for economics; Fulfillment Operator as handoff consumer.

```text
choose exact Installation
→ list/select source-qualified Sale
→ inspect Sale
→ resolve Selling Entity attribution only when admitted/needed
→ compose Sale economics
→ inspect BusinessOrderIntent + Party resolution + destination realization
→ inspect InvoicingIntent
→ inspect FulfillmentExecution + PostSaleResolution where present
→ outcome: known cross-owner state + next safe owner/handoff/work item
```

D8 P5 law: a contact reference is not presented as a supported full alternate street/fiscal destination override.

### UF11 — Execute physical fulfillment and observe Shipment

**Actor:** Fulfillment / Dispatch Operator.

```text
inspect actionable Fulfillment executions
→ select exact execution
→ record separation
→ record physical conference
→ block normal invoicing/dispatch path on physical inconsistency
→ observe/verify fiscal + provider prerequisite readiness owned by accepted semantics
→ inspect required artifacts
→ record packing
→ record dispatch handoff
→ observe source-qualified Shipment through relevant outcome
→ outcome: verified dispatch/delivery state or explicit exception Work
```

UI control visibility never self-authorizes a physical checkpoint. D8 P2 remains a real future fiscal/invoice/label proof obligation.

### UF12 — Coordinate essential Post-Sale resolution

**Actor:** Marketplace Operations Operator.

```text
inspect PostSale resolutions
→ create a resolution when the accepted need exists
→ inspect current resolution + owner evidence
→ coordinate through underlying owners/Work
→ outcome: explicit current resolution state
```

No generic close/refund/cancel command is invented; closure depends on authoritative evidence.

### UF13 — Coordinate operational Work

**Actors:** Marketplace Operations Operator; Fulfillment Operator; managers when escalated.

```text
list Work
→ inspect exact Work item + related source truth
→ assign/clear assignment OR hold/resume OR escalate as admitted
→ return to underlying owner evidence when action is required there
→ outcome: coordination state progresses without Work impersonating source truth
```

### UF14 — Govern decisions, delegations and access

**Actor:** Owner/Administrator/Policy Approver; delegated Commercial Manager for admitted decisions.

```text
inspect pending/history Authorization Decisions
→ inspect exact target/revision evidence
→ create decision when governance authority exists
→ inspect/manage authorization delegations as separately admitted
→ inspect/manage Organization AccessRoles when access administration is required
→ outcome: explicit governance/access state, separate from target execution
```

A Governance decision never grants the target Permission by itself.

### UF15 — Configure fulfillment operating context

**Actors:** Owner/Administrator or delegated operational manager.

```text
inspect Fulfillment Nodes + effective operational targets
→ create/update/deactivate accepted Fulfillment Node configuration as admitted
→ update MPC-owned internal target when authorized
→ re-read effective configuration
→ outcome: explicit execution context and internal target provenance
```

This remains Marketplace-fulfillment configuration, not a company-wide WMS/TMS model.

### UF16 — Govern commercial policy and reconcile economic attribution

**Actor:** Commercial / Marketplace Manager; exceptional authority may hand off to Policy Approver.

```text
inspect current Commercial Policy + effective evidence
→ update only MPC-owned policy with current precondition when authorized
→ separately inspect realized/economic attribution evidence
→ resolve bounded economic attribution ambiguity when authorized
→ re-read realized/effective state
→ outcome: explainable current policy/economics or explicit unresolved/approval-needed state
```

Externally governed rules are never silently converted into editable MPC policy.

### 8.1 P2 exit

**DERIVED.** Every accepted human need N01–N16 has a complete end-to-end flow, including the safety branch that changes safe behavior. No flow requires a new Product operation.

---

# P3 — Frontend Coverage Matrix

## 9. Coverage law

D6-B1 already proved 99/99 Product operations have coherent interaction homes. D6-R2 P3 re-expresses that proof from user needs/flows rather than treating the prior 39 screen states as a design target.

Exact per-operation owner/`operationId`/Permission/Principal-kind/identity/state/effect binding is a **P9 obligation**. P3 must establish complete capability/user-flow coverage and expose any orphan before screen design.

## 10. Coverage matrix

| Product owner / family | Ops | User needs / flows | Accepted candidate context | Reads | Writes / controlled effects | Access/security baseline | UX obligations | P3 status |
| --- | ---: | --- | --- | --- | --- | --- | --- | --- |
| Identity / Access | 5 | N01, N14 / UF01, UF14 | shell + Configurações/Acesso | current access, members, roles | assign/revoke AccessRole | H browser session; `access.read/manage` where applicable; server auth wins | canonical Organization, display labels never identity, stale/deep-link access fails closed | COVERED |
| Marketplace Portfolio | 6 | N02, N08 / UF02, UF10 | Configurações/Canais + Selling Entity context | Installations, Selling Entities | create/update/deactivate Installation | `portfolio.read/manage`; OAuth is separate Technical Ingress | available kind != connected Installation; no fake provider catalog; exact Installation identity | COVERED |
| Product & Channel Readiness | 5 | N03 / UF03 | Preparação | source-product search, readiness, requirements | resolve/clear correspondence | `readiness.read/manage` | SourceInstance explicit; omission is bounded search, never hidden default; unsupported/missing stays visible | COVERED |
| Offering | 12 | N04, N06 / UF04, UF08 | Publicações / ListingIntent / PriceIntent | listings, intents, current state | ListingIntent create/update/discard/media/submit; PriceIntent create | `offering.read`, `listing.manage`, `price.manage` as exact OAD binds | ETag/If-Match, idempotency, no provider direct write, accepted/pending/rejected/ambiguous distinct | COVERED |
| Availability | 9 | N07 / UF09 | Disponibilidade + Configurações/Disponibilidade | sellable state, Inventory Sources, effective policy | source/policy configuration only | `availability.read/manage` | derived availability not manual quantity truth; unknown != zero; no generic sync | COVERED |
| Market Intelligence | 3 | N06 / UF08 | Mercado | competitive position/offers | none | `market.read` | evidence/read composition only; no Price write authority | COVERED |
| Marketplace Performance Intelligence | 4 | N05 / UF05–UF07 | Performance Resumo/Publicações/Mídia + Listing context | summary, listing performance, Retail Media evidence | none | `performance.read`; exact Installation where OAD requires | explicit periods, comparable gate, source coverage/provenance/basis, no Ads controls | COVERED |
| Commercial Economics | 11 | N06, N16 / UF08, UF16 | Economia + Configurações/Política comercial | expected/realized economics, summary, attribution, policy; stateless scenario evaluation | update commercial policy; resolve economic attribution | `economics.read`, `economics.policy.manage`, `economics.reconcile` as exact OAD binds | Cost Basis/provenance honest; scenario evaluation not mutation; external policy not editable local truth; preconditions visible | COVERED |
| Controlled Action Governance | 7 | N13 / UF14 | Aprovações + Configurações/Delegações | decisions/delegations | create decision; establish/update/revoke delegation | `governance.read/decide/manage` as exact OAD binds | decision != target Permission/execution; target/revision explicit | COVERED |
| Marketplace Sales | 3 | N08 / UF10 | Vendas | list/get Sale | bounded Selling Entity attribution resolution | exact Sales OAD permissions; H browser only where admitted | source-qualified Sale identity; no synthetic cross-Installation merge | COVERED |
| Business-System Materialization | 7 | N08 / UF10 | Sale/materialization context | BusinessOrderIntent, Party resolution, destination, InvoicingIntent | only admitted bounded resolution effects; no screen-invented ERP command | exact Materialization OAD permissions | external business/fiscal authority remains external; D8 P5 capability narrowing visible | COVERED |
| Fulfillment + artifacts + Shipment | 17 | N09, N10, N15 / UF11, UF15 | Expedição + Configurações/Expedição | executions, artifacts, shipments, nodes/targets | physical checkpoints + admitted node/target configuration | `fulfillment.read/execute/manage` as exact OAD binds; physical qualification server-owned | conference before normal invoicing; provider prerequisites honest; no WMS/TMS expansion; D8 P2 proof debt preserved | COVERED |
| Post-Sale | 3 | N11 / UF12 | Pós-venda | list/get resolution | create accepted resolution | `post_sale.read/manage` | no generic close/refund/cancel; source owners determine consequence/closure | COVERED |
| Operational Work | 7 | N12 / UF13 | Trabalho | list/get Work | assign/clear/hold/resume/escalate | `work.read/manage` | Work coordinates but never becomes source truth/command bus | COVERED |
| **Total** | **99** | **N01–N16 / UF01–UF16** | accepted D6 IA |  |  | **30 ordinary Permissions; H/A/S canonical principal universe** | **100% accepted Product operation coverage** | **COVERED** |

## 11. Orphan-operation disposition

Current accepted D6-B1 coverage is **99/99**, and the current OAD executable proof remains **99 operations / 30 ordinary Permissions / H-A-S only**.

Therefore P3 finds:

```text
unmapped Product operation families: 0
screen-shaped operations required: 0
new Product capabilities invented by D6-R2: 0
```

This does **not** mean 99 buttons or 99 screens. It also does not let P3 guess per-operation Principal-kind or Permission details: P9 must bind those exact OAD facts for every material screen/action and explicitly surface any operation that is not human-consumable.

If P9 discovers an operation whose current OAD trust class contradicts the accepted D6-B1 human home, that is a material finding and the smallest owning authority must be reopened before P14.

## 12. Cross-cutting UX obligations entering P4+

Every later candidate must preserve at least:

```text
unknown != known-empty
partial != complete
projection != mutation authority
hidden control != authorization
ambiguous outcome != known failure
stale write != business rejection
external mechanism success != Product success
URL/context != business authority
source-qualified identity != globally merged identity
Governance decision != target Permission
Work coordination != source truth
Performance evidence != reconstructed provider metric
contact reference != full alternate fiscal/street destination
```

## 13. P0–P3 findings

### No material architecture contradiction

P0–P3 currently produce **no finding that justifies reopening D0–D8**. The accepted IA and frontend realization topology can represent all recovered human needs and all 14 Product operation families.

### Open evidence findings, not architecture defects

- **E01 — task frequency/density evidence missing:** A01 remains OPEN and may affect P4/P5/P8 prioritization.
- **E02 — real device/work-floor evidence missing:** A02 remains OPEN and may affect responsive interaction structure.
- **E03 — terminology comprehension not directly tested:** A03 remains OPEN.
- **E04 — Overview visual priority not directly evidenced:** A04 remains OPEN.
- **E05 — no accepted bulk-interaction need:** A05 remains OPEN; no bulk UI is admitted by assumption.

These may inform candidates but unresolved material assumptions must be probed by P12 before P14 closes.

## 14. P3 exit

**DERIVED.** Accepted human capabilities are mapped, the accepted 99-operation Product surface remains fully covered without invented capability, and current material unknowns are evidence assumptions rather than hidden design decisions.

No P8 block is `LOCKED`; P8 has not started.

---

# 15. Exact next action — P4 bounded IA revalidation

Run P4 against N01–N16 + UF01–UF16 + the P3 matrix.

P4 must:

1. preserve the accepted D6 shell/route grammar as the baseline;
2. create the durable user-language terminology/object-task/findability trace required by Method v2.1;
3. test whether each need is findable from the accepted IA without exposing backend package taxonomy;
4. record a material falsifier if — and only if — a need cannot be coherently found/contained;
5. avoid layout composition and wireframe drawing;
6. avoid P6/P7 unless a real ambiguity trigger appears.

If P4 converges, P5 derives the complete screen/surface inventory. P8 remains blocked until the required upstream structural decisions are ready for operator-reviewed rendered blocks.
