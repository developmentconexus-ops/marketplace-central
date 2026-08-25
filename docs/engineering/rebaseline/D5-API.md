# D5 — Product API

> **Status:** CLOSED / ACCEPTED / CURRENT CONSOLIDATED AUTHORITY  
> **Parent authorities:** accepted D0→D4 + D4-R1  
> **Machine-readable wire authority:** `contracts/api/product/openapi.yaml`  
> **Current Product surface:** **106 operations / 31 ordinary Permissions / Principal kinds H-A-S**  
> **Runtime:** NONE; implementation blocked until accepted D9

## 1. Purpose

D5 defines the semantic Product API used by MPC clients. HTTP/OpenAPI exposes accepted owner meaning; it does not create Product authority, provider ontology or frontend-screen authority.

Current Product API structure is:

```text
D5 semantic laws
→ B2 operation/resource surface
→ W1 resource/path/HTTP grammar
→ W2 schema/read-write grammar
→ W3 collection/search/cursor grammar
→ W4 Permission/Principal admission
→ one canonical OpenAPI wire
```

Provider/business-system protocol ingress remains D4/technical ingress and is not counted as Product operations merely because it uses HTTP.

---

## 2. Binding API invariants

1. Every Product operation belongs to exactly one accepted semantic owner or the bounded D2 identity/access substrate.
2. Organization is explicit on Organization-owned Product operations; secondary Organization-owned references fail closed across Organizations.
3. External identities remain source-qualified; a bare provider/native key is never MPC-global identity.
4. Product vocabulary is MPC semantic vocabulary. Provider DTO/status/endpoint topology stays behind D4.
5. Q preserves honest knowledge/freshness/provenance. Query failure never becomes false/zero/empty/ready/permitted.
6. `accepted != completed != externally applied != converged`; pending/ambiguous/divergent states remain explicit when reachable.
7. C operations preserve consequence class, duplicate/idempotency disposition and concurrency/precondition semantics independently.
8. A timeout after possible consequential acceptance is not automatically rejection and never authorizes blind replay.
9. Ordinary Permission, Principal/client-class admission, business disposition, Governance authorization and execution-time validity are separate gates.
10. Read projections/snapshots never become write authority.
11. No generic Product `Mutation`, `Action`, `Workflow`, `Entity`, provider-field bag or screen-shaped BFF operation is admitted by convenience.
12. One machine-readable Product wire authority exists; generated clients/server contracts conform to it rather than forming parallel wire authority.

---

## 3. Authentication / request-trust baseline

Current Product authentication admits exactly:

```text
human H
→ external OIDC Authorization Code flow terminates server-side
→ MPC server-side session cookie
→ CSRF protection on unsafe browser methods
→ browser receives no OIDC access/refresh token

automation/system A|S
→ confidential machine authentication / client-credentials bearer
→ resolved MPC Principal
```

Authentication never supplies Organization/Permission/business authorization by implication. `/access-context` is the bounded platform-scoped self-discovery exception; all other business state is explicitly Organization-scoped.

Exact current security schemes/session/CSRF wire fields live in canonical OpenAPI and are mechanically proved.

---

## 4. Operation admission

A Product operation exists only for a real Product 1.0 client job that needs:

```text
Q — read one owner's current meaning
C — ask one owner to accept/perform owner-owned work
P — consume a justified read-only composition
```

API symmetry, provider endpoint symmetry, current code, debug convenience or hypothetical future need do not admit operations.

The current admitted Product census is **106 operations**. The exact current operationId/path/method/owner/class/Permission/Principal mapping is the canonical OpenAPI projection of the accepted operation surface and is protected by Product OAD proof. `D5-B2-OPERATION-ADMISSION-MATRIX.md` records current admission decisions/families without duplicating another machine wire catalog.

A count change is material and requires explicit operation admission/revalidation; it never enters through incidental schema/frontend work.

---

## 5. Resource and capability model

MPC-owned durable identities use stable Organization-rooted resource grammar when externally addressable. Provider-owned resources retain their source namespace (for marketplace-native resources, usually Marketplace Installation + native key).

Owner-specific capabilities use explicit custom operations only when ordinary resource CRUD semantics would lie. URI nesting expresses stable scope/identity, not internal bounded-context/package topology or workflow stage.

No `/v1` compatibility prefix is required without a real entitled compatibility consumer. Hard cutover remains allowed before production compatibility obligations exist.

Exact path/method semantics are W1 authority.

---

## 6. Request/write safety

### 6.1 Idempotency

Consequential creation/intake defaults to mandatory caller semantic idempotency unless an explicitly admitted structural owner anchor makes duplicate intake unreachable/harmless. An idempotency key deduplicates the same semantic intake; it is not authorization and never makes ambiguous provider replay safe.

### 6.2 Concurrency/preconditions

Owner resource mutations use the smallest typed revision/precondition semantics justified by that owner. Concurrency and idempotency solve different failures.

Current special case: `CreateAuthorizationDecision` decides one `AuthorizationRequest`. The request body carries the exact current **AuthorizationRequest ETag** plus outcome; it does not carry/recompute a target ETag. Request revision protects concurrent terminal decision while D2/D3 material-validity revalidation remains separately authoritative.

### 6.3 External-effect ambiguity

Product command acceptance and provider effect/convergence remain distinct. If caller cannot rule out accepted MPC work, it reconciles by owner semantic anchor/current read rather than blindly creating another intent/effect.

---

## 7. Read / schema semantics

W2 owns exact semantic schema grammar. Current cross-cutting laws include:

- stable IDs/keys are explicit and source-qualified where external;
- exact decimal Money/rates/material quantity where correctness requires it;
- known/known-empty/unknown/unavailable and owner-specific knowledge variants remain distinct;
- request/write schemas carry decision authority only; read schemas may carry current projection/presentation/evidence required by accepted current contracts;
- historical/purpose snapshots are distinct from current source truth;
- provider richness is exposed only as typed, source-qualified evidence needed by a Product consumer;
- arbitrary provider payload/metadata bags are forbidden.

The human-operable read-projection gap discovered later in D6-R2 is a **separate paused prerequisite (#70)** and is not silently implemented by this repository-health consolidation.

---

## 8. Collection/query semantics

W3 owns collection/search/cursor semantics. Shared laws:

- named owner-specific collection envelopes, not universal `Page<T>`;
- ListItem is an owner-semantic subset, not a second business conclusion;
- opaque forward cursor; explicit semantic query repeats on continuation;
- cursor exhaustion != source/owner completeness;
- no universal total count or caller-selectable sort;
- Product cursor never exposes provider/database paging internals;
- filters/search are smallest typed vocabularies justified by real consumers.

---

## 9. Permission / Principal admission

W4 owns the ordinary-access vocabulary and admission laws. Current Product vocabulary is **31 ordinary Permissions**, exactly flat/no wildcard implication. Principal kinds remain only `H`, `A`, `S`.

Special authenticated/self surfaces do not create stored Permissions. Possessing Permission never substitutes for current Membership/access eligibility, business disposition or Governance authorization.

---

## 10. Notifications / AuthorizationRequest post-baseline surface

Later accepted D6-R2 planning added bounded Product surfaces without changing D5 fundamentals:

### Personal Notifications

```text
ListMyNotifications
UpdateMyNotificationAwarenessState
ListNotificationRoutes
ListNotificationRouteRecipientCandidates
SetNotificationRoute
```

Self Inbox awareness is H-only/current-Membership bounded; Organization routing administration uses `notifications.manage`. Notification access never grants source capability.

### AuthorizationRequest

```text
ListMyActionableAuthorizationRequests
GetMyActionableAuthorizationRequest
CreateAuthorizationDecision  // route is request-scoped :decide
```

Actionable Request reads and decision are H-only `governance.decide`. No new `authorization_requests.read` Permission exists. Actionable views expose only what the exact Principal can decide now; Governance history remains separately governed.

Exact paths/schemas are W1/W2/W3/W4 + current OAD.

---

## 11. Problems / response integrity

Product failures use stable Product Problem semantics; provider/native error ontology does not leak into Product API merely because an adapter returned it. Authentication, authorization/privacy not-found, validation, conflict/stale-precondition, unsupported media/content size and internal failures remain distinguishable where applicable.

An error response never fabricates provider/business success/failure after an ambiguous potentially accepted effect.

---

## 12. OpenAPI / generated-contract authority

`contracts/api/product/openapi.yaml` is the sole machine-readable Product wire entrypoint. Local refs compose it; supported TypeScript/Go projections are generated/proved from that authority.

Current proof protects, proportionately:

- OpenAPI validity and resolved refs;
- exact current Product operation/Permission/Principal surface;
- authentication profile;
- source-tree reachability/orphan policy;
- material schema/knowledge/authorization/read invariants;
- deterministic TypeScript/Go projection and generated route expressibility.

Verification mechanics remain tooling, not Product authority.

---

## 13. Explicit non-goals / reopen triggers

D5 does not introduce a generic API gateway/domain bus, provider-normalization platform, BFF/screen API, dynamic query DSL, generic evidence/metadata API, compatibility/version framework without a consumer, or D7 runtime topology.

Reopen only the smallest implicated D5 owner when material evidence proves a new Product operation, Permission, client class, identity/path meaning, knowledge/effect distinction, schema relationship or collection capability is necessary. Frontend convenience alone is not sufficient; a proven human need may falsify the existing contract and reopen the smallest semantic/wire owner under the Engineering + Frontend methods.
