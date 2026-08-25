# D7 — Runtime / Jobs / Transactions

> **Status:** ACCEPTED / CLOSED / CURRENT CONSOLIDATED AUTHORITY  
> **Parent authorities:** accepted D0–D6 + canonical Product OAD  
> **Current Product:** 106 operations / 31 ordinary Permissions / H-A-S  
> **Active runtime:** NONE; implementation blocked until D9

## 1. Target invariant

> **Every admitted Product/Technical-Ingress interaction can be realized through explicit serving, persistence, transaction, durable-work, authentication and operability boundaries that preserve Organization isolation, owner authority, wire conformity, request trust, idempotency/concurrency, recoverable external-effect semantics and secret/PII safety while platform mechanisms remain business-policy-free.**

D7 chooses runtime mechanism only. It does not reopen Product operations, Permissions, owners or frontend meaning by convenience.

## 2. Accepted runtime envelope

One Go application process per replica contains proportionately:

```text
same-origin frontend/static delivery
Product API
Technical Non-Product Ingress
H session / CSRF / OIDC mediation
in-process River worker runner
scheduler/recovery seam
PostgreSQL pool
```

One process does not merge owners. Internal Q/C/P cross explicit owner application boundaries rather than self-HTTP. E becomes durable only where D3 requires independent recoverable reaction.

API/worker process split is deferred until measured scaling/resource/security/deployment evidence requires it.

## 3. D7-B — PostgreSQL / isolation / transactions

Accepted realization:

```text
organization_id on Organization-owned state/evidence
composite Organization FKs
transaction-local Organization scope
ENABLE + FORCE RLS
runtime role non-owner / NOSUPERUSER / NOBYPASSRLS
READ COMMITTED baseline + explicit row locks
opaque random owner revision tokens
owner/operation scoped idempotency
pgx/v5 + pgxpool
```

A narrow `authentication_bootstrap` persistence mode exists before Principal resolution for one-time login/session/machine-client binding only; it cannot read Organization business/evidence state.

Owner consequential intake commits owner state + required correctness records in one owner transaction. External network effects occur only after commit. No cross-owner distributed transaction exists.

## 4. D7-C — durable work / external effects

River over PostgreSQL is the accepted durable-reaction mechanism:

```text
InsertTx atomic owner-state → durable reaction
repeat-safe delivery; no exactly-once claim
owner semantic idempotency = correctness
persisted pre-dispatch/attempt evidence for consequential external effect
possible acceptance → reconciliation, not blind redispatch
source-authoritative reread controls convergence
scheduler/job state = technical wake-up only
```

No second generic MPC outbox/broker/Redis is added without evidence. River completion/retry/uniqueness never becomes business truth.

After PITR/acknowledged-state rollback, database absence is not proof no later external effect occurred. Consequential dispatch remains recovery-fenced until timeline continuity/current external truth is positively re-established.

## 5. D7-D — authentication / session / CSRF

Accepted baseline:

```text
Keycloak first provider

H:
  confidential Authorization Code + PKCE S256
  verified (issuer, sub) → opaque MPC server session
  human OIDC token set discarded after callback
  session-handle digest persisted
  30m idle / 8h absolute baseline
  X-CSRF-Token synchronizer + CrossOriginProtection

A/S:
  Client Credentials
  expected audience includes https://conexus.fun
  trusted-JWKS verification
  explicit machine-client → A/S Principal binding
```

No browser bearer, persistent human refresh-token cache, IdP role→MPC Permission mapping, Redis session baseline or A/S→H impersonation.

## 6. D7-E — HTTP / operability / deployment

Accepted direction:

```text
HTTP       Chi v5 + oapi-codegen strict server
validation canonical OAD middleware
bytes      private S3-compatible object storage + authenticated Go delivery
config     typed startup config + injected env/file secrets
migrate    tern/v2 MPC schema + exact-version River migration ownership
logs       structured slog JSON
telemetry  OpenTelemetry traces/metrics over OTLP/HTTP
artifact   one immutable OCI app image (Go + compiled frontend)
edge       https://conexus.fun behind explicit trusted TLS proxy/ingress
backup     PostgreSQL base backup + WAL/PITR or managed equivalent + restore proof
```

No Kubernetes/service mesh/Redis/external broker/Vault/CDN/public bucket/multi-region/generic IaC platform is required by baseline.

## 7. AuthorizationRequest runtime — current accepted composition

The later AuthorizationRequest repair is now part of D7 current authority; no separate workflow engine/runtime stack exists.

### 7.1 Decision intake / current carrier

`CreateAuthorizationDecision` current wire:

```text
POST /organizations/{organization_id}/authorization-requests/{authorization_request_id}:decide
H only / governance.decide
Idempotency-Key header
body.etag = current AuthorizationRequest StrongETag
body.outcome = authorize | reject
```

The custom capability uses **typed body ETag**, not `If-Match`.

Current idempotency namespace:

```text
organization_id
+ effective PrincipalID
+ CreateAuthorizationDecision operation identity
+ digest(Idempotency-Key)
```

Semantic fingerprint includes:

```text
authorization_request_id
+ typed request etag
+ outcome
```

Same raw key under another Principal is a different namespace and never replays/discloses the first Principal's Decision.

### 7.2 Replay before revision precondition

Inside Governance processing:

```text
current AuthN/access gate
→ exact idempotency lookup/claim
   ├─ committed same fingerprint → replay original result
   ├─ different fingerprint      → reused-key validation failure
   ├─ in-progress                → current intake conflict
   └─ new                       → evaluate current Request
```

Exact committed replay resolves **before current Request revision-precondition evaluation**. This lets a lost 201 be recovered even though the Request's supplied revision became historical because that same Decision already committed.

For a genuinely new attempt:

```text
lock exact AuthorizationRequest
→ still PENDING
→ compare typed body.etag
→ current exact-human eligibility
→ action-owner material-validity Q
→ commit terminal meaning or known no-effect response
```

Missing/invalid ETag = 422; stale ETag = 409 revision conflict.

### 7.3 Decision validity Q

Material-validity Q is an in-process action-owner query over current MPC-owned/evidenced truth. Governance does not hold its owner transaction open across provider/business-system network I/O.

Outcomes:

```text
VALID
→ Decision + Request DECIDED + terminal idempotent replay mapping + durable reactions

INVALID
→ Request INVALIDATED + replayable terminal outcome + reconciliation reactions
→ no Decision / no F14

UNKNOWN_OR_UNAVAILABLE
→ no Decision / no Request lifecycle mutation
→ exact typed authorization-validity-unavailable 503
```

Only the exact Product Problem is known-no-effect. Bodyless/proxy/unparseable/non-matching 503 remains ambiguous potentially accepted.

### 7.4 Current eligibility

Actionable Request list/detail and new Decision attempts derive current decision eligibility from Governance authority/delegation + current IdentityAccess truth. No Notification/event/cache is eligibility authority.

A revocation committed before the authoritative eligibility check must be observed. Later drift does not rewrite historical Decision; execution-time source-owner revalidation remains binding.

### 7.5 F13 / F14 durable materialization

F13 materializes only after revalidating:

```text
Request still PENDING
+ exact human still eligible now
```

Same eligibility occurrence replay is idempotent; later legitimate re-eligibility may be a new occurrence. Historical Notifications are not rewritten.

F14 is anchored to immutable Decision occurrence + exact requester recipient; duplicate/rescued work cannot duplicate awareness. It remains target-oriented and grants no `governance.read`.

### 7.6 Zero-decider Work

```text
PENDING Request + known-empty eligible humans
→ ensure one explicit zero-decider Work obligation
```

When a decider exists again or Request becomes terminal, Governance/Work reconciliation closes the no-longer-applicable obligation. Assignment/escalation never grants `governance.decide`.

### 7.7 Recovery sweep

Correctness does not depend on every eligibility/invalidation wake-up. Fast local wakeups are paired with bounded durable recovery over PENDING Governance truth:

```text
scan durable pending Requests
→ Q current eligibility
→ Q material validity where required
→ reconcile F13 / zero-decider Work / invalidation obligations
```

Scheduler cursor/tick is technical state only. A missed event/tick cannot permanently strand required current reconciliation.

## 8. Notification temporal materialization

For ORG_ROUTED occurrences, durable jobs carry only bounded correlation/Organization/source-commit facts; Personal Notifications resolves the route revision that applied at source commit and current eligibility continuity before new materialization.

Late/replayed jobs cannot apply a newer route to an older occurrence. Source occurrence duplicate identity remains owner-semantic; River job uniqueness is only an optimization.

## 9. Whole-D7 proof contract

Eventual implementation must prove with real seams, not mocks alone:

- cross-Organization DB/RLS/reference isolation and pooled-scope leakage prevention;
- authentication-bootstrap containment;
- owner revision/idempotency semantics, including Request typed-ETag replay order;
- owner-state→River atomic handoff and duplicate/rescued/out-of-order reaction safety;
- ambiguous external effect non-redispatch and PITR recovery fencing;
- browser token/session/CSRF/OIDC safety;
- machine issuer/audience/JWKS/client binding and A/S/H separation;
- OAD validation and colon-suffix custom route behavior;
- private byte authorization;
- secret/PII log/job/telemetry safety;
- migration/restore/continuity integrity;
- current Notification temporal routing and AuthorizationRequest recovery obligations.

D7 architecture closeout does **not** claim runtime implementation PASS while active runtime is NONE.

## 10. Reopen triggers

Reopen D7 only when real implementation/proof demonstrates the accepted mechanisms cannot preserve an invariant without materially different runtime/persistence/deployment structure. Framework preference, a desire for another broker/cache/service split or historic proof-file shape is not evidence.
