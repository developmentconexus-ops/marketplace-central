# D6-R2 — Complete Frontend Realization Closure

> **Status:** OPEN / ACTIVE — P0–P4 DERIVED; accepted D6 IA survives bounded revalidation as the D6-R2 `CANDIDATE`; no D6-R2 UX block is operator-`LOCKED`
> **Program:** Architecture Rebaseline / Technical System Design
> **Parent authority:** accepted D0–D8, accepted D6 frontend authority, canonical Product OAD
> **Execution methods:** [DevelopmentConexus Engineering Method v1.0.0](../../development/engineering-method.md) + [Frontend Product Experience Planning Method v2.3](../../development/frontend-product-experience-planning-method.md)
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

No screen is drawn in P0–P4.

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
| Local Engineering Method v1.0.0 + Frontend Product Experience Planning Method v2.3 | Global Maximum/root-cause/proof discipline plus human needs before screens, coverage before layout, operator-only `LOCKED`, no screen-shaped API, P0–P14 readiness discipline. |

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
| Marketplace Operations Operator | A marketplace account/product/listing/sale requires normal operation or exception handling | Prepare channel participation, control listings/price intent within policy, understand performance/market/economics, inspect sales and resolve operational exceptions | Marketplace participation is convergent, explainable and operable without bypassing policy | Marketplace requirements, source values/correspondence, Listing/Intent state, availability, market/economic evidence, performance, sale/work state | Hands physical execution to Fulfillment; hands exceptional authority to Governance/manager; cannot redefine policy to make an action permissible |
| Fulfillment / Dispatch Operator | A marketplace sale reaches internally operated physical execution | Know what is ready, separate/confer correctly, trigger/observe the accepted fiscal/provider prerequisites, pack and hand off dispatch safely | Correct physical goods reach verified dispatch without premature invoicing or fabricated readiness | Sale/execution identity, separation/conference state, fiscal/provider readiness, artifacts, shipment state, exception work | Physical inconsistency blocks the normal path and becomes explicit Work; UI cannot self-authorize physical checkpoints |
| Commercial / Marketplace Manager | Commercial evidence, policy thresholds or exceptions require a decision | Compare market/performance/economics, evaluate price trade-offs, manage legitimate MPC-owned policy and resolve bounded commercial/economic ambiguity | Commercial actions remain profitable/explainable and inside governing policy | Performance, market evidence, expected/realized economics, effective policy/provenance, approvals/exceptions | Does not edit externally governed rules; does not become integration/security administrator by default |
| Owner / Administrator / Policy Approver | Access, integrations, exceptional authority, delegation or containment requires governance | Establish valid access/channel configuration and approve exceptional actions above delegated authority | Organization remains safely operable with explicit accountability and least necessary authority | Access context, members/roles, Installations, authorization decisions/delegations, configuration state | Not a routine approval bottleneck; cannot configure away audit/reconciliation/safety invariants |

## 5. Need statements independent of screens

### N01 — Work in the correct Organization

When a human enters MPC or changes workspace, they need to know which Organization and access context are current, so that every subsequent observation/action is scoped correctly.

### N02 — Establish a usable marketplace channel

When the organization needs to operate on a supported marketplace, an authorized human needs to create/configure the Marketplace Installation and complete any required external authorization ceremony, so that later Product operations have an explicit valid channel context.

### N03 — Prepare a source product for marketplace authoring

When a product should participate in a marketplace, the operator needs to discover the exact source product, understand which fields the marketplace asks for, inspect which source values are available, and resolve correspondence when necessary, so that listing configuration starts from explicit source truth without treating missing source data as publication impossibility.

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
→ inspect marketplace publication requirements + available source values/evidence
→ preserve missing/conflicting/unknown/unavailable/unsupported evidence explicitly
→ resolve/clear correspondence only when needed and authorized
→ authoritative reread after a correspondence effect
→ outcome: exact prepared subject + explicit values/gaps ready for ListingIntent authoring
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
| Product & Channel Readiness | 5 | N03 / UF03 | Preparação | source-product search, publication requirements/source evidence, correspondence/current context | resolve/clear correspondence | `readiness.read/manage` | SourceInstance explicit; omission is bounded search, never hidden default; missing/conflicting/unknown/unavailable/unsupported stays visible; no per-requirement satisfaction inference | COVERED |
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

P0–P3 produce **no finding that justifies reopening D0–D8**. The accepted IA and frontend realization topology can represent all recovered human needs and all 14 Product operation families.

### Open evidence findings, not architecture defects

- **E01 — task frequency/density evidence missing:** A01 remains OPEN and may affect P5/P8 prioritization.
- **E02 — real device/work-floor evidence missing:** A02 remains OPEN and may affect responsive interaction structure.
- **E03 — terminology comprehension not directly tested:** A03 remains OPEN.
- **E04 — Overview visual priority not directly evidenced:** A04 remains OPEN.
- **E05 — no accepted bulk-interaction need:** A05 remains OPEN; no bulk UI is admitted by assumption.

These may inform candidates but unresolved material assumptions must be probed by P12 before P14 closes.

## 14. P3 exit

**DERIVED.** Accepted human capabilities are mapped, the accepted 99-operation Product surface remains fully covered without invented capability, and current material unknowns are evidence assumptions rather than hidden design decisions.

No P8 block is `LOCKED`; P8 has not started.

---

# P4 — Bounded Information-Architecture Revalidation

## 15. Revalidation law

P4 does **not** ask what IA would be preferred if Marketplace Central were being redesigned today. It asks whether the already accepted D6 IA is falsified by the user needs, complete flows or Product coverage recovered in P0–P3.

A material IA falsifier exists only if an accepted need cannot be coherently found/contained without one of these failures:

```text
new top-level navigation created only to mirror a backend owner/package
hidden/default Organization, Installation or source identity required to reach the task
one semantic owner collapsed into another merely for navigation convenience
material flow step has no findable human context
cross-owner composition requires a new business/Strategy/Workflow authority
user must understand an internal API/domain term merely to find the task
```

No layout, visual hierarchy or component selection is part of this pass.

## 16. Durable user-language glossary

Accepted primary IA labels remain accepted D6 authority. Detailed object/action wording below is a **D6-R2 `CANDIDATE` terminology baseline**, because A03 direct comprehension evidence remains open. Canonical IDs and exact technical names remain available for correlation/support where needed but do not become primary labels by convenience.

| Authority concept | Candidate user-facing term | Avoid as primary UX label | Status / note |
| --- | --- | --- | --- |
| Organization | Organização | tenant, org id | `Organização` is presentation; canonical id remains scope authority |
| Marketplace Installation | Conta do marketplace, inside **Canais** | Installation, integração genérica | Exact account identity; `Canal` is grouping language, not a replacement identity |
| Selling Entity | Entidade vendedora | legal-entity DTO/package names | accepted Settings grouping |
| Source Product | Produto de origem | Product master MPC | source remains external/source-qualified |
| SourceInstance | Fonte de origem | source instance id as user noun | exact technical key may be secondary evidence |
| Product Channel Readiness | Preparação | readiness/prontidão como status por campo | primary nav label **Preparação** is accepted; owner truth supplies requirements/correspondence without the UI inventing satisfaction |
| Publication Requirements | Requisitos de publicação | provider requirement DTO | provenance/unsupported state stays explicit |
| Product Channel Correspondence | Correspondência do produto | linkage entity/package | detail term remains candidate pending A03 |
| Marketplace Listing | Publicação | listing resource as unexplained English noun | primary nav label **Publicações** is accepted |
| ListingIntent | Intenção de publicação | provider mutation | MPC intent is distinct from provider state |
| PriceIntent | Intenção de preço | SetPrice, alteração direta no marketplace | Offering intent remains distinct from Economics evidence |
| Sellable Availability | Disponibilidade para venda | estoque = disponibilidade | derived availability is not native stock truth |
| Inventory Source | Fonte de estoque | adapter/source table | configuration object under Disponibilidade |
| Availability Allocation Policy | Política de disponibilidade | regra genérica/sync rule | only MPC-owned policy is editable |
| Marketplace Performance | Performance | analytics genérico | accepted strategic lens |
| Retail Media Performance | Performance de mídia / Mídia | Ads Manager | read-only evidence; no Ads management |
| Competitive Position | Posição no mercado | MarketIntelligence owner name | accepted **Mercado** lens |
| Comparable Offers | Ofertas comparáveis | competitor DTOs | evidence, not Product price authority |
| Expected Economics | Economia prevista | cost engine | accepted **Economia** lens |
| Sale / realized economics | Economia realizada da venda | settlement DTO as navigation | contextual economic evidence |
| Economic Attribution | Reconciliação econômica / atribuição econômica | internal attribution engine | exact final object label remains A03 candidate |
| Marketplace Sale | Venda | order/provider package | source-qualified marketplace sale |
| BusinessOrderIntent / native materialization | Pedido no sistema de negócio / materialização do pedido | Sankhya command/TOP as user navigation | remains contextual to Sale, not a top-level owner area |
| Party Resolution | Cliente/parceiro no sistema de negócio | Party resolver | ambiguity remains explicit |
| Destination Realization | Destino de entrega | `CODCONTATO` as destination capability | D8 P5 narrowing must be visible |
| InvoicingIntent / fiscal progression | Faturamento | generic ERP command | intent/result distinction remains owner-qualified |
| FulfillmentExecution | Execução de expedição | workflow engine | accepted **Expedição** lens |
| Shipment | Envio | carrier/TMS object as new owner | source-qualified observation |
| PostSaleResolution | Resolução pós-venda | generic close/refund/cancel | accepted **Pós-venda** lens |
| Operational Work | Trabalho | task command bus | coordination only; source truth stays elsewhere |
| AuthorizationDecision | Decisão de aprovação | policy engine command | target action remains separate |
| AuthorizationDelegation | Delegação | IAM role | accepted Configurações grouping |
| AccessRole | Perfil de acesso | provider/IdP role | Product access presentation only |
| FulfillmentNode | Nó de expedição | WMS node model | Marketplace-fulfillment configuration only |
| Operational Target | Meta operacional | provider deadline | internal target never rewrites external obligation |

Terminology falsifier result: no accepted user need requires a backend package/API-family name as the way to find the task. A03 remains OPEN for comprehension testing before P14.

## 17. Object/task inventory against accepted navigation

| Human context | Why users care | Primary task family | Important relationships | Findability role | Independent-nav verdict |
| --- | --- | --- | --- | --- | --- |
| Organization shell | establishes scope for everything | confirm/switch current workspace | all Product contexts | persistent global context | **KEEP** global shell context; not a content route |
| Visão geral | orientation/attention entry | understand bounded current posture | Portfolio, Economics, Work, Performance when permitted | home/entry composition | **KEEP** accepted landing; content priority remains A04 |
| Preparação | prepare marketplace-required data before authoring | search source product, inspect marketplace fields/source values, resolve correspondence | Publicações / ListingIntent | search-first task context | **KEEP** independent operational destination |
| Publicações | operate provider-facing offering intent/state | browse listings/intents, author/submit listing/price intent | Preparação, Disponibilidade, Performance, Mercado, Economia | primary operational destination + contextual detail | **KEEP** independent destination |
| Disponibilidade | know sellable availability | inspect sellable state; reach source/policy configuration | Publicações, Configurações | operational read destination | **KEEP** independent destination; configuration remains Settings |
| Vendas | understand source-qualified sale and next safe handoff | browse sale, inspect composed cross-owner state | Economia, materialization, Expedição, Pós-venda, Trabalho | primary sale entry + contextual composition | **KEEP** independent destination |
| Expedição | execute physical fulfillment and observe shipment | work executions/checkpoints/artifacts; observe shipments | Vendas, Trabalho, Configurações | operational queue/browse context | **KEEP** independent destination |
| Pós-venda | coordinate essential cancellation/return/refund consequences | inspect/create resolution | Vendas, Trabalho, Economics consequences | operational exception destination | **KEEP** independent destination |
| Performance | understand our marketplace participation | summary, listing performance, media evidence | Publicações, Mercado, Economia | strategic sub-navigation | **KEEP** separate from Mercado/Economia |
| Mercado | understand competitive environment | competitive position/offers | Publicações, Economia | strategic evidence destination | **KEEP** separate owner lens; no write authority |
| Economia | understand expected/realized economics | inspect, simulate, reconcile attribution | Mercado, Performance, Publicações/Vendas, policy Settings | strategic/economic destination | **KEEP** independent destination |
| Trabalho | coordinate ambiguity/divergence | inspect/assign/hold/resume/escalate work | every underlying owner resource | cross-cutting work queue with contextual return links | **KEEP** control destination; never source truth |
| Aprovações | make required governance decisions | inspect/create decisions for exact target/revision | target action/resource, Delegações | governance queue/context | **KEEP** control destination |
| Configurações | low-frequency organization-operable setup | channels, selling entities, access, availability, fulfillment, commercial policy, delegations | corresponding operational lenses | grouped low-frequency settings | **KEEP** grouping only; never semantic owner |

No row justifies exposing `IdentityAccess`, `BusinessSystemMaterialization`, `ControlledActionGovernance`, `MarketplacePerformanceIntelligence` or another backend owner name as primary navigation.

## 18. Navigation model revalidation

### 18.1 Global frame

The accepted frame survives:

```text
Organization selector
→ grouped primary navigation
→ page-local/context navigation
→ contextual cross-links for cross-owner evidence/handoffs
```

`Organization` remains the only global workspace. Marketplace Installation does **not** become a second global tenant/workspace; exact/all-or-exact Installation selection remains typed page/URL context only where admitted.

### 18.2 Primary grouping

The accepted groups remain coherent with N01–N16:

```text
VISÃO GERAL
OPERAÇÕES
ESTRATÉGIA E INTELIGÊNCIA
CONTROLE
CONFIGURAÇÕES
```

No need requires a new top-level `Strategy`, `Analytics`, `Integrações`, `Materialização`, `Workflow`, `ERP` or `Admin` product area.

### 18.3 Local subcontexts that remain valid

P4 confirms these are user-task subdivisions rather than backend taxonomy exposure:

- **Performance:** Resumo / Publicações / Mídia;
- **Publicações:** provider Listing observation, ListingIntent work and PriceIntent work remain one Offering-oriented human context;
- **Economia:** expected/realized/reconciliation meanings may remain local subcontexts under the accepted route family;
- **Expedição:** execution and Shipment observation remain related but distinct local contexts;
- **Configurações:** Canais, Entidades vendedoras, Acesso, Disponibilidade, Expedição, Política comercial and Delegações remain low-frequency grouped settings.

Exact route/surface splitting is P5, not P4.

### 18.4 Cross-owner navigation is contextual, not a new owner

The following relationships are deliberately satisfied by contextual links/composition rather than new primary navigation:

```text
Preparação → author ListingIntent
Publicação → Performance / Disponibilidade / Mercado / Economia evidence
Mercado / Economia → PriceIntent under Offering when a price decision is made
Venda → Economics / business-order materialization / destination / invoicing / fulfillment / post-sale
Execução/Shipment → related Sale and Work
Pós-venda → related Sale / Work / economic consequences
Trabalho → underlying source-owner resource
Aprovação → exact governed target/revision
```

A cross-link never moves mutation authority into the source page that displays it.

### 18.5 Home/default landing

`Visão geral` remains the accepted default landing because D6 already ratified it as a small read-only composition. P4 does not infer which cards/signals are visually primary. That is A04 and belongs to P5/P8 plus P12 evidence.

### 18.6 Search and breadcrumbs

- no global Product search is introduced; the only search-first context currently proved is **Preparação/source-product search** through admitted Product authority;
- no universal breadcrumb system is introduced merely because nested URLs exist;
- detail surfaces may expose parent/context return navigation when P5 proves it useful, without implying a deeper domain hierarchy than the user model has.

## 19. Findability decisions for material collections

Mechanisms below are the smallest justified baseline. Exact filters/sorts are admitted only when the Product contract and P5/P9 screen contract support them.

| Collection/task | Baseline findability | Explicitly not assumed |
| --- | --- | --- |
| source products for marketplace preparation | **search-first** in Preparação | global search, hidden SourceInstance default |
| Marketplace Listings | browse/list under Publicações with exact Installation | synthetic all-provider merged list |
| ListingIntents / PriceIntents | local browse/work contexts under Publicações | separate top-level Intent areas, generic mutation inbox |
| Listing performance | Performance/Publicações + contextual Listing link | survivorship filtering that hides unknown/unavailable Performance |
| Retail Media evidence | Performance/Mídia | Ads campaign manager |
| sellable availability | browse/list under Disponibilidade | manual quantity editor or generic Sync action |
| Marketplace Sales | browse/list under Vendas with exact Installation | cross-Installation synthetic identity merge |
| fulfillment executions | operational queue/list under Expedição | company-wide WMS queue platform |
| Shipments | local browse under Expedição | TMS/carrier navigation owner |
| Post-Sale resolutions | browse/work context under Pós-venda | generic case-management/CRM platform |
| Work | coordination queue under Trabalho | source-resource replacement |
| Authorization decisions | decision queue/history under Aprovações | target execution inside Governance |
| configuration collections | grouped browse/manage under Configurações | provider/plugin catalog or generic settings object model |

No saved-view platform, alternate-view framework or bulk-selection model is justified at P4. A05 remains OPEN rather than becoming a bulk UX feature.

## 20. N01–N16 findability falsification matrix

| Need | Primary entry from accepted IA | Required contextual continuation | Result |
| --- | --- | --- | --- |
| N01 correct Organization | persistent Organization shell | Configurações/Acesso only for administration | **PASS** |
| N02 usable marketplace channel | Configurações → Canais | external OAuth ceremony for exact account when required | **PASS** |
| N03 prepare marketplace data | Preparação | Publicações/ListingIntent after source-value inspection and correspondence when needed | **PASS** |
| N04 listing/price intent | Publicações | contextual availability/economics and approval when required | **PASS** |
| N05 marketplace Performance | Performance → Resumo/Publicações/Mídia | Listing contextual link where useful | **PASS** |
| N06 competitive/economic position | Mercado + Economia | PriceIntent continuation returns to Publicações/Offering | **PASS** |
| N07 sellable availability | Disponibilidade | Configurações/Disponibilidade for source/policy management | **PASS** |
| N08 cross-system Sale understanding | Vendas | contextual economics/materialization/fulfillment/post-sale/work | **PASS** |
| N09 physical fulfillment | Expedição → Execuções | Sale/Work context as needed | **PASS** |
| N10 shipment/delivery outcome | Expedição → Envios | related Sale/Work | **PASS** |
| N11 essential post-sale | Pós-venda | related Sale/Work/economic evidence | **PASS** |
| N12 operational Work | Trabalho | contextual return to underlying owner | **PASS** |
| N13 exceptional approval/delegation | Aprovações + Configurações/Delegações | exact governed target/revision | **PASS** |
| N14 access administration | Configurações/Acesso | current access context in shell | **PASS** |
| N15 fulfillment operating configuration | Configurações/Expedição | operational Expedição lens | **PASS** |
| N16 policy/economic attribution | Configurações/Política comercial + Economia/Reconciliation context | Aprovações only when governance requires | **PASS** |

Result:

```text
needs tested                     16
findable/containable             16
material IA falsifiers            0
new top-level destinations needed 0
backend-owner navigation leaks     0
```

The sale/materialization path is the strongest negative control: P4 deliberately **rejects** a top-level `Materialização` area. Business-order/Party/destination/invoicing meaning remains findable from the human Sale context and retains its own semantic owner underneath.

## 21. Mental-model validation and P6/P7 trigger evaluation

P4 uses the smallest proportionate validation because the global IA is already operator-ratified D6 authority and P0–P3 produced no new navigation need. The structural walkthrough asks whether a user goal can be located from the accepted tree without knowing backend taxonomy; all 16 pass.

This is **not** direct user-comprehension proof. A03 remains open for P12 terminology/first-click comprehension validation.

Current trigger disposition:

| Trigger | P4 evidence | Disposition |
| --- | --- | --- |
| high-impact IA uncertainty requiring card sorting/tree test now | no uncontained need or competing category placement survives the N01–N16 walk | **NOT TRIGGERED in P4**; A03 remains P12 evidence obligation |
| P6 unfamiliar/high-impact UX problem or meaningful pattern uncertainty | P4 decides navigation/findability only; no new unfamiliar layout problem introduced | **NOT TRIGGERED by P4** |
| P7 genuine structural ambiguity with >1 credible IA hypothesis | accepted D6 tree contains every need without contradiction; alternatives would be preference-driven | **NOT TRIGGERED by P4** |

P6/P7 may still trigger later for an individual P5/P8 material block if a real layout/interaction ambiguity appears.

## 22. Future seams kept out of live IA

The following may shape extensibility but do not earn live routes/navigation now:

- marketplaces beyond the proved Mercado Livre operating loop;
- all-marketplace KPI aggregation without measurement equivalence;
- Ads campaign/budget/bid/creative management;
- global Product search or saved-view platform;
- bulk Product actions without an admitted Product capability and user evidence;
- multi-Organization SaaS provisioning/billing administration;
- generic provider/plugin/integration catalog;
- provider-native composite-offer management unless a selected real flow makes it material;
- AI/agent/autonomous commercial navigation authority.

## 23. P4 findings and exit

### Material finding disposition

**No material IA falsifier found.** P4 does not reopen D6, D0–D5 or any D7/D8 authority.

The accepted D6 IA therefore survives as the D6-R2 **`CANDIDATE` global IA** for downstream screen/surface derivation:

```text
accepted D6 IA
+ N01–N16 / UF01–UF16 falsification
+ 99/99 P3 coverage
→ P4 CANDIDATE REVALIDATED
```

This does **not** create a new D6-R2 `LOCKED` decision. Under Frontend Product Experience Planning Method v2.3, the global frame becomes eligible for an operator `LOCKED` decision only after the first required rendered global-frame block cycle. Existing D6 operator ratification remains architectural input authority; D6-R2 does not pretend that an unrendered P4 trace is the P8 visual lock.

### Open evidence carried forward

- A01 frequency/density remains OPEN for P5/P8 prioritization and P12 validation.
- A02 device/work-floor evidence remains OPEN for responsive P8 structure.
- A03 detailed terminology/comprehension remains OPEN; glossary detail is `CANDIDATE`.
- A04 Overview signal priority remains OPEN.
- A05 bulk interaction remains unsupported by evidence and therefore absent.

### P4 exit

**DERIVED / REVALIDATED.** All 16 accepted needs are findable within the accepted D6 IA, no backend taxonomy must become navigation, no new Product capability is required, and no real P6/P7 ambiguity trigger appears at the IA level.

P8 remains **NOT STARTED** and no D6-R2 UX block is `LOCKED`.

---

# 24. Exact next action — P5 complete screen/material-surface inventory

Derive P5 from N01–N16 + UF01–UF16 + P3 coverage + the P4 candidate IA.

P5 must:

1. derive screens/surfaces from user-flow decisions, not from endpoints;
2. distinguish route/page, material surface, drawer/modal, inline composition, alternate view and material state variant;
3. use the prior D6-B1 39-state inventory as coverage evidence, **not** as a target count to preserve;
4. split a material surface only when semantic truth, safe action, write owner, identity, concurrency/idempotency, disclosure/security, recovery or editor/viewer mode materially changes;
5. preserve cross-owner composition without creating new write authority;
6. name any block whose layout/interaction ambiguity legitimately triggers P6/P7;
7. keep P8 wireframing blocked until P5 gives a coherent block inventory and sequencing basis.
