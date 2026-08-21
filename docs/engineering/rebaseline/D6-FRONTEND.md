# D6 — Frontend

> **Status:** OPEN / ACTIVE — D6-B1 Frontend Interaction & Authority Model is the current bounded decision
> **Program:** Architecture Rebaseline / Technical System Design
> **Opened:** 2026-08-21
> **Parent authorities:** `ARCHITECTURE.md`, accepted D0–D5 semantics, and canonical Product OAD at `contracts/api/product/openapi.yaml`
> **Method:** DevelopmentConexus Engineering Method v1.0.0

## 1. Purpose and boundary

D6 defines the target frontend interaction model and frontend topology for Marketplace Central without creating a second business authority.

D6 begins from the accepted Product API and stable architecture constraints. It does not inherit removed frontend structure, legacy routes, package layout, component hierarchy or state-management choices merely because they existed before the rebaseline.

D6 does **not** choose or implement:

- server/router/runtime mechanics;
- database, transaction, worker, scheduler, queue, outbox or deployment mechanics;
- D8 golden-flow execution choreography;
- Product implementation, which remains blocked until D9 is accepted.

React is already the accepted client technology at the stable architecture level. TanStack Query remains the accepted server-state client unless material evidence explicitly reopens that decision. D6 decides how the frontend consumes accepted Product authority; it does not move business policy into React.

---

## 2. Imported invariants

D6 imports rather than re-decides:

1. **The frontend is a Product API client, not a domain authority.** Business policy remains server-owned.
2. **OpenAPI is the single machine-readable Product API wire authority.** The frontend does not create a second hand-written wire contract.
3. **Organization is the isolation root.** Organization-scoped Product interactions preserve explicit `/organizations/{organization_id}/...` scope and never infer Organization from client-local defaults or provider identity.
4. **Ordinary Permission is not business disposition.** Hidden routes/buttons may improve usability but never become authorization.
5. **Knowledge states remain honest.** Known-empty, unknown, unavailable and partial are not interchangeable; material freshness/provenance is not replaced by browser/request time.
6. **Consequential outcomes remain distinct.** Accepted, rejected, pending and ambiguous never collapse into generic success/failure.
7. **Consequential intake preserves idempotency semantics.** Client retry ergonomics do not authorize duplicate intake or blind replay of an ambiguous external effect.
8. **Source-qualified external identity remains explicit where required by D5.** Client state never invents a global identity from a bare provider/native ID.
9. **Provider/business-system protocol ingress remains outside Product frontend semantics.** Provider callbacks/OAuth/technical ingress are not frontend shortcuts to business capability.
10. **No compatibility tax without a real consumer.** Removed frontend/package conventions are evidence only.
11. **No screen-shaped second API authority.** D6 consumes D5; it does not redefine Product API semantics around presentation convenience.

---

# 3. D6-B1 — Frontend Interaction & Authority Model — OPEN / ACTIVE

## 3.1 Decision question

> What is the smallest frontend interaction contract that lets an MPC user observe and initiate accepted Product capabilities without duplicating server/domain authority, weakening D5 semantics or pre-committing D7 runtime mechanics?

## 3.2 Root cause to prevent

The primary frontend architecture risk is not component count or styling consistency.

> **Client-side convenience can silently become a second authority for server state, Organization/access meaning, business disposition, knowledge/freshness or consequential-work lifecycle.**

If that condition is admitted, it can produce:

- stale duplicated server truth in Context/store state;
- route/button visibility being mistaken for authorization;
- empty/loading/error/unknown/partial states collapsing into one presentation state;
- accepted/pending/ambiguous effects being shown or retried as generic success/failure;
- duplicate consequential intake after client retry;
- screen-specific API assumptions becoming a second Product contract;
- provider/native identifiers escaping their accepted namespace qualification;
- frontend package abstractions becoming accidental business owners.

## 3.3 Target invariant

> **Every material frontend observation or user-initiated Product action derives from an accepted semantic owner and canonical Product API operation/capability; the frontend preserves Organization, Permission, knowledge/freshness, concurrency, idempotency and consequential-outcome semantics without duplicating business truth or turning client state into authority.**

## 3.4 State ownership model

D6-B1 starts with four distinct state classes:

| State class | Target owner | D6 law |
| --- | --- | --- |
| Server state | accepted Product/server authority, consumed through TanStack Query | do not duplicate into Context/global client stores by convenience |
| URL/navigation state | browser/navigation contract | may select/view/filter/navigate but does not become business truth |
| Form draft state | user-editing session until submitted/accepted | remains draft; unsent values do not impersonate accepted server meaning |
| Ephemeral UI state | frontend only | presentation-only state such as open panels, focus or temporary selection |

Any proposed fifth durable/global state class requires a concrete consumer, protected property and evidence that these four classes are insufficient.

## 3.5 Interaction laws

For D6 work unless a later accepted D6 decision narrows them:

1. Every target screen or material interaction maps to an accepted owner plus explicit Product API query/capability/operation.
2. Screen composition may combine read authorities, but the composition never becomes write authority.
3. Route/button hiding is usability only; server authorization remains authoritative.
4. Missing Permission/access evidence is never converted into a client-local business rejection.
5. Known-empty, unknown, unavailable, partial and materially stale states receive distinguishable behavior where the Product contract distinguishes them.
6. `accepted`, `pending`, `rejected` and `ambiguous` consequential outcomes remain distinguishable in user-visible interaction state.
7. An ambiguous consequential outcome never enables an automatic or generic "retry" path that could blind-replay an external effect.
8. Required idempotency carriers are preserved across safe client retries/resolution of the same intake; a new semantic request does not silently reuse the old key.
9. Concurrency/precondition failures remain distinct from business/provider preconditions.
10. Provider/OAuth/technical-ingress endpoints are not normal frontend Product operations.
11. Frontend code consumes the canonical Product wire contract through a mechanically derived/conforming client boundary; no hand-written second DTO/wire authority is introduced.
12. A Product API operation is not admitted or reshaped merely because a screen would be easier to implement that way.

## 3.6 YAGNI exclusions at B1 opening

D6-B1 does not introduce without material evidence and a named consumer:

- a second global server-state store;
- generic workflow/state-machine authority;
- generic frontend command/action/mutation business layer;
- a BFF or screen-shaped Product API;
- offline-first synchronization;
- websocket/event-stream architecture;
- micro-frontends or plugin architecture;
- universal design-system/platform work beyond demonstrated Product need;
- SSR/meta-framework selection merely by convention;
- new router/form/state libraries merely by preference.

A later D6 decision may evaluate concrete frontend dependencies when an accepted interaction/topology property requires them. Dependency choice then follows the repository evidence-grounded engineering rules and current official documentation for the exact version.

## 3.7 Known / inferred / unknown / deferred

**Known**

- React is the accepted frontend client technology.
- TanStack Query is the accepted server-state client unless materially reopened.
- The canonical Product OAD is the Product wire authority.
- There are 95 admitted Product operations and 29 ordinary Permissions.
- the global App Shell / IA is operator-approved as Organization-global + task-oriented navigation + contextual Marketplace Installation + low-frequency Settings + read-only Overview composition;
- Product 1.0 currently exposes Mercado Livre as the only connectable marketplace kind;
- the bounded D2-R1 presentation-identity correction is proved without changing the 95/29/H-A-S surface;
- `SearchSourceProductsForMarketplace` admits an optional operation-local SourceInstance narrowing filter while every returned Product remains source-qualified;
- the derived [D6-B1 Frontend Interaction Map](D6-B1-INTERACTION-MAP.md) gives all 95 admitted operations a user-visible interaction home across 32 screen/route states;
- the low-fidelity HTML proof set exists and has passed internal challenge against the interaction-map negative controls without exposing another Product/API gap.

**Inferred**

- separating server, navigation, form-draft and ephemeral UI state removes the main known route for duplicate client authority without requiring another state framework;
- a small representative wireframe set can falsify the interaction model without drawing every read/detail variant as an independent application.

**Unknown**

- exact frontend feature/package topology;
- exact router/form/component/design-system dependency choices, if any;
- exact browser/runtime realization details that depend on later implementation.

**Deferred**

- server/runtime/persistence mechanics → D7;
- executable end-to-end golden flows → D8;
- adversarial global architecture acceptance → D9;
- Product code realization → after D9 acceptance.

## 3.8 Proof strategy before B1 acceptance

D6-B1 is not accepted merely because this document exists.

Before B1 ratification, the smallest proof package must be capable of falsifying the target invariant through at least:

1. a target screen/interaction → semantic owner → Product operation/query/capability → ordinary Permission mapping for the bounded Product 1.0 frontend surface under review;
2. explicit state treatment for known-empty / unknown / unavailable / partial / stale where reachable;
3. explicit consequential interaction treatment for accepted / pending / rejected / ambiguous and required idempotency/concurrency carriers where reachable;
4. negative counterexamples proving that route visibility, client-local state, projections or provider/native vocabulary cannot become business/access authority;
5. a check that no frontend requirement silently demands a screen-shaped API, provider protocol shortcut or D7 mechanic;
6. independent challenge when the resulting B1 decision materially creates or moves an authority/trust boundary, following repository review rules.

The mapping/proof candidate is maintained in [D6-B1 Frontend Interaction Map](D6-B1-INTERACTION-MAP.md).

Browser/component implementation evidence is not required to open or reason about D6 while Product implementation remains blocked. Claims that later depend on actual browser behavior require browser-level evidence at the stage where execution is authorized.

## 3.9 First frontend coverage findings and bounded parent corrections

The first D0–D5 frontend-coverage pass deliberately attacked channel onboarding and global shell identity before screen-by-screen design.

### Available channel versus connected Installation

The intended UX architecture distinguishes:

```text
available channel kind
  != connected Marketplace Installation/account
  != current channel operational context
  != low-frequency configuration/control
```

For Product 1.0 the currently admitted connectable marketplace kind remains **Mercado Livre only**. D5 already rejected a generic provider/integration catalog, so D6 does not reopen D4/D5 merely to imitate marketplace-hub catalogs or to advertise hypothetical Amazon/Shopee support.

The frontend may expose a stable **Add channel** interaction architecture, but the set of connectable channel kinds must follow admitted Product/integration authority. It must not hardcode future providers as currently connectable. A future provider becomes visible as connectable only after its required D4/D5 support is explicitly admitted; this should extend the existing UX model rather than require a different information architecture.

### Authorization ceremony

Provider authorization is not a missing Product operation. The accepted D5 Technical Ingress authorization ceremony remains separate from the Product API: a current Product-authorized human may begin authorization for an exact existing Marketplace Installation, then provider callback/current-authority revalidation establishes or restores the bound credential generation. The frontend must not invent `ConnectMarketplace` Product semantics around that technical ceremony.

### Presentation identity falsifier

D6 found one real bounded read-completeness gap: already-admitted human interactions cannot reliably render Organization/access context when only opaque `organization_id`, `principal_id` and `role_key` values are guaranteed.

The operator-approved [D2-R1 Presentation Identity](D2-R1-PRESENTATION-IDENTITY.md) plus D5 wire correction adds required human-readable `display_name` metadata for the current Principal, accessible Organizations, Organization members and AccessRoles while preserving canonical ID/key authority.

This repair creates no new Product operation, Permission or Principal kind and does not make OIDC/profile/provider names canonical identity.

### SourceInstance discovery falsifier

The first Readiness user flow exposed a second bounded D5 wire gap: `SearchSourceProductsForMarketplace` originally required a `source_instance_id`, while the Product frontend intentionally has no generic SourceInstance registry/discovery operation.

The operator approved the smallest correction:

- keep the Product inventory at 95 operations and 29 ordinary Permissions;
- make `source_instance_id` an **optional narrowing filter only on `SearchSourceProductsForMarketplace`**;
- omission means bounded search across current Organization-scoped SourceInstances admitted/configured for the Readiness search context, never ambient/default-source selection;
- every result remains explicitly `SourceInstance + native_product_key` qualified;
- `GetProductChannelReadiness` and `GetPublicationRequirements` continue to require the exact SourceInstance selected from the result.

The correction is operation-local and does not create a Source Registry screen/domain or alter D2/D4 source authority.

## 3.10 App Shell / IA decision — OPERATOR APPROVED

The approved target is:

> **Organization-global + task-oriented primary navigation + explicit contextual Marketplace Installation where the Product contract admits that dimension + low-frequency Settings separated from routine operation + read-only Overview composition.**

Binding consequences:

- primary navigation groups user work, not internal D1 package names;
- Marketplace Installation context may filter/select an exact marketplace namespace but never becomes ambient business authority;
- `All accounts` is shown only where the admitted read honestly supports optional Installation scope;
- independently paginated source-qualified Marketplace Listing/Sales/Shipment collections are not client-merged into a fake complete cross-account list;
- Settings groups low-frequency configuration without becoming a new business owner;
- Overview uses existing owner Qs only and never becomes `/dashboard` write/read authority;
- permission-conditioned menu/button visibility is usability only;
- contextual drawers may show evidence/history/related Work without becoming write authority;
- responsive behavior changes layout, not business semantics.

Detailed flows, routes, 32 screen states, 95-operation coverage and negative controls live in [D6-B1 Frontend Interaction Map](D6-B1-INTERACTION-MAP.md).

## 3.11 Low-fidelity HTML proof — CANDIDATE / INTERNAL CHALLENGE PASS

The self-contained wireframe prototype is [`qualification/d6-wireframes/index.html`](../../../qualification/d6-wireframes/index.html).

It intentionally contains only low-fidelity evidence:

- grayscale/system-font layout rather than final branding, typography, component library or design system;
- ten representative navigable states covering the approved shell, Readiness, ListingIntent authoring, Availability, Sale composition, Fulfillment execution, Work, Channels settings, Economics/reconciliation and Access/Governance settings;
- illustrative/static values only; no claim of Product runtime, browser integration, provider behavior or persistence;
- responsive structural behavior sufficient to test sidebar/top-context/content stacking without choosing a router/framework.

Internal adversarial review attacked the interaction-map negative controls and found no new Product/API prerequisite. The prototype explicitly demonstrates or blocks the material hazards:

- no hidden/default SourceInstance or Product master;
- no false cross-account merge for source-qualified collections;
- no generic `ConnectMarketplace`, sync/refresh, SetPrice, SetAvailableQuantity, Work close, PostSale close or deferred reactivation action;
- ListingIntent/Price/Availability/Economics remain distinct authorities inside one user flow;
- Sale detail is read-only composition with owner-local actions;
- unknown/unavailable/partial and pending/ambiguous states are visibly distinct;
- stale revision/idempotency/retry semantics are surfaced instead of hidden;
- physical Fulfillment evidence does not trust caller-declared qualification;
- Work does not acquire source-truth closure;
- Governance does not acquire target-domain Permission/execution authority;
- presentation labels remain visibly secondary to canonical IDs/keys.

This is **not browser/runtime proof** and does not ratify D6-B1 by itself. It is an interaction-level falsifier produced while Product implementation remains blocked.

---

## 4. Exact next D6 work

Continue only inside D6-B1:

1. operator-review/adjudicate the [D6-B1 Frontend Interaction Map](D6-B1-INTERACTION-MAP.md) and the [low-fidelity HTML wireframe proof](../../../qualification/d6-wireframes/index.html);
2. if the interaction model is approved, evaluate the smallest frontend feature/package topology and exact dependency needs required by those accepted properties, using current official evidence rather than convention;
3. topology research may select frontend navigation/form/client realization details but must not choose D7 server/runtime/router/database/deployment mechanics;
4. once one coherent D6-B1 interaction + topology candidate exists, submit it to the repository's independent milestone/final review path rather than using Fable/Claude Design as iterative co-authors;
5. do not ratify D6-B1 or merge PR #54 without explicit operator authorization.

Do not begin D7–D9 or Product implementation.