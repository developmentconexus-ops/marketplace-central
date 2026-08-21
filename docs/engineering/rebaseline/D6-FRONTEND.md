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
- Product implementation is blocked until D9.

**Inferred**

- A screen/capability map is the smallest useful artifact for proving that presentation does not create unnamed business authority.
- Separating server, navigation, form-draft and ephemeral UI state removes the main known route for duplicate client authority without requiring another state framework.

**Unknown**

- final target screen inventory and information architecture;
- exact frontend feature/package topology;
- exact router/form/component/design-system dependency choices, if any;
- which concrete interactions require optimistic concurrency UX or long-running consequential-work UX;
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

Browser/component implementation evidence is not required to open or reason about D6 while Product implementation remains blocked. Claims that later depend on actual browser behavior require browser-level evidence at the stage where execution is authorized.

---

## 4. Exact next D6 work

Continue only inside D6-B1:

1. derive the smallest Product 1.0 screen/interaction inventory from accepted Product objectives and the canonical OAD;
2. map each material interaction to its semantic owner, Product operation/query/capability and ordinary Permission;
3. attack the mapping for duplicate authority, hidden screen-shaped API assumptions, dishonest knowledge states and unsafe consequential retries;
4. only after the interaction model is coherent, evaluate the minimum frontend feature/package topology needed to realize it.

Do not begin D7–D9 or Product implementation.
