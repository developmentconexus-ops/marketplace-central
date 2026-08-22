# D7-R1 — Whole-Stage Coherence Corrections & Proof-Timing Contract

> **Status:** CANDIDATE / INTERNAL WHOLE-D7 REVIEW — BOUNDED FIXES / INDEPENDENT CHALLENGE PENDING  
> **Parent:** `D7-RUNTIME-JOBS-TRANSACTIONS.md`  
> **Scope:** exact D7-B/C/D/E realization seams only; no D0–D6 semantic reopen  
> **Method:** DevelopmentConexus Engineering Method v1.0.0  
> **Derived:** 2026-08-22

## 1. Purpose

Whole-D7 cross-check found a small set of realization seams that become contradictory only when D7-A→D7-E are composed as one runtime. This amendment is the smallest bounded repair candidate.

It does **not** change Product operations, Permissions, Principal kinds, owner semantics, D3 Q/C/E/P meaning, D4 external-effect meaning, D5 wire/auth contract or D6 frontend authority.

If accepted, this document supersedes only the exact D7 clauses named below.

## 2. Finding R1-F1 — authentication bootstrap is missing from D7-B scope taxonomy

### Problem

D7-B says its persistence scope modes are exhaustive, but D7-D requires PostgreSQL lookups that occur **before** one current Principal or Organization is known:

- one-time OIDC login transaction by `state`;
- ApplicationSession by opaque-handle digest;
- machine client binding by trusted issuer + client identity.

These lookups cannot honestly be `organization`, `principal_self` or `technical_routing` without inventing identity/scope that has not yet been established.

### Bounded correction

D7-B persistence access is refined to:

```text
scoped MPC modes
  organization
  principal_self
  authentication_bootstrap
  technical_routing

explicit platform-global technical/definition state
  product Permission/role definitions where truly global
  River engine tables/library technical state
```

`authentication_bootstrap` is narrowly limited to authentication mechanism state needed to resolve one current Principal:

```text
one-time login transaction
ApplicationSession lookup
machine-client -> Principal binding
```

Binding laws:

- no default Organization exists;
- authentication-bootstrap access cannot read organization-owned business/evidence tables;
- login/session lookup keys do not become business identity;
- after a Principal is resolved, authentication-bootstrap scope ends; eligibility/self-discovery uses `principal_self`, and organization-owned work uses a fresh `organization` unit of work;
- River engine tables are platform technical state, not organization-owned business/history. Their arguments remain bounded routing metadata under D7-C and contain no credentials/raw provider payload/arbitrary PII. Claiming a River job never grants business-table access or Organization authority.

The exact PostgreSQL policy/table spelling remains implementation detail; the structural proof must show auth-bootstrap and technical/library state cannot bypass Organization RLS on business/evidence tables.

## 3. Finding R1-F2 — River schema migrations are not MPC `tern` migrations

### Problem

D7-C selects River, whose queue schema has its own versioned migrations. D7-E selects `tern/v2` for MPC-owned PostgreSQL schema. Treating `tern` as if it also owned River's internal schema would create a copied/parallel vendor-schema authority or leave River schema compatibility undefined.

### Bounded correction

Migration authority is split by ownership, in one controlled release step:

```text
MPC-owned schema
  -> tern/v2

River-owned technical schema
  -> exact-version-pinned River migration CLI/API
     using the same River release line as the application dependency
```

Baseline release sequence:

```text
migration credential / schema owner
  -> tern migrate MPC schema
  -> river migrate-up (pinned; configured River schema)
  -> verify expected MPC + River schema compatibility
  -> deploy/start application with ordinary runtime credential
```

Laws:

- neither MPC nor River migrations run automatically in normal application startup;
- application runtime never receives the migration-owner credential;
- a River library upgrade that carries a schema migration is one reviewed dependency+schema change;
- raw River SQL obtained through River tooling may be inspected as derived vendor evidence but is not hand-edited into an independent MPC schema authority;
- if either migration lane fails, the new application version does not roll forward to serving traffic.

No generic multi-tool migration framework is introduced.

## 4. Finding R1-F3 — normal crash safety does not survive database time rollback automatically

### Problem

D7-C correctly says, on a continuous durable database timeline:

```text
no dispatch marker -> pre-dispatch may proceed
marker + no definitive outcome -> possible acceptance -> reconcile, never redispatch
```

D7-E also admits backup/PITR restore. A restore to an earlier database point can remove a dispatch marker/idempotency/session/access change while an already-performed external effect or later security revocation survived elsewhere. After such a rollback, **absence of a restored marker is no longer evidence that no dispatch occurred**.

Blindly resuming restored River jobs would therefore violate D4/D7-C no-blind-retry semantics.

### Bounded correction — recovery fence

The D7-C `no marker -> may dispatch` rule is valid only on the same durable database timeline with no acknowledged-state rollback.

Any disaster/PITR restore that may have lost acknowledged MPC commits enters a fail-closed **recovery fence** before normal Product readiness:

1. record/identify the restored recovery point and possible state-loss window;
2. invalidate restored human ApplicationSessions and one-time login transactions so a lost post-backup logout/revocation cannot resurrect a browser session;
3. keep normal consequential external-write dispatch disabled, including restored River external-effect jobs;
4. treat restored dispatchable/ambiguous work as reconciliation-only until its owner/source scope is cleared;
5. reacquire authoritative provider/business-system state for affected Installations/SourceInstances and reconcile restored owner intents/effect anchors over the possible loss window;
6. revalidate current identity/access and machine-client configuration before normal authenticated automation resumes when the restore could have rolled those changes back;
7. verify committed binary references at the restored database point against object-store digest/size/type evidence;
8. release the fence only for scopes whose reconciliation/access/integrity checks are sufficient. If coverage cannot establish safety, automated external writes remain disabled and the condition is surfaced for operator resolution rather than guessed away.

New durable intake may be recorded during a bounded recovery posture only if its external dispatch remains behind the same recovery fence.

A zero-data-loss failover whose durable database timeline is proved continuous need not enter the rollback recovery path.

### Required recovery falsifier

Restore a database snapshot from **before** a simulated externally accepted write while the external system retains the effect. Prove that restored River/owner state does not redispatch before reconciliation/fence release.

## 5. Finding R1-F4 — OAD security validation must compose with, not bypass, D7-D authentication

### Problem

Current `nethttp-middleware`/kin-openapi request validation evaluates OpenAPI security requirements through an `AuthenticationFunc`; a missing function fails validation, while a no-op function could make the schema validator appear to satisfy security without enforcing D7-D carrier semantics.

D7-E must therefore define the exact authority seam rather than leave `AuthenticationFunc` as incidental middleware configuration.

### Bounded correction

Product request composition is:

```text
Chi route/operation resolution
-> D7-D carrier authentication
     H: MPC HttpOnly session
     A/S: machine bearer
-> unsafe H CSRF + cross-origin trust check
-> OAD request validator
     scheme-aware AuthenticationFunc checks the already-established carrier context
     MpcHumanSessionAuth succeeds only for H-session context
     MpcMachineBearerAuth succeeds only for A/S-machine context
-> current Principal eligibility / Membership / allowed kind / Permission
   using generated operation-policy metadata
-> owner/resource/idempotency/revision/business/Governance gates
```

Laws:

- the OAD validator does not independently reimplement token/session cryptography;
- `AuthenticationFunc` is never `NoopAuthenticationFunc` in production composition;
- absence of the pre-established carrier context fails security validation;
- a request presenting both human-session and machine-bearer carriers fails closed as ambiguous; no carrier priority/fallback is invented;
- generated OAD operation-policy metadata remains derived; IdP roles/scopes never become MPC Permissions;
- Product validator/security errors are translated to accepted Product failure grammar, not raw kin-openapi errors.

Required negative controls include removal of carrier context, carrier/scheme mismatch and simultaneous H+machine carriers.

## 6. Finding R1-F5 — D7 proof timing must not contradict the implementation gate

### Problem

Some slice prose says real PostgreSQL/River/Keycloak/browser proof is mandatory "before D7 closeout", while repository authority still says Product implementation is blocked until D9 and active runtime baseline is NONE.

D7 may select mechanisms and define executable falsifiers, but it cannot claim a Product runtime PASS that does not yet exist.

### Bounded correction

Use three distinct proof statuses:

```text
D7 architecture closeout
  -> coherent, falsifiable executable proof contract
  -> current primary evidence for selected mechanism capabilities
  -> optional isolated mechanism spikes are evidence only
  -> NO Product runtime conformance claim

D8 / D9
  -> consume the D7 proof seams during golden-flow/adversarial architecture review
  -> still no Product implementation unless roadmap later permits it

post-D9 implementation acceptance
  -> execute required real-dependency proofs against the implemented runtime
  -> PostgreSQL/River/Keycloak/browser/router-validator/object-store claims cannot close on mocks alone
```

Accordingly, prior D7-B/C/D wording that real-dependency execution is required to **close D7** is read as: D7 must define the real-dependency executable falsifier and may not replace it with a mock-only proof claim. Actual Product runtime execution belongs to implementation acceptance after the implementation gate opens.

Existing current gate output that reports no active runtime/schema enforcement remains honest: D7 selects target realization; it does not populate Product runtime code.

## 7. Internal review disposition

Current internal whole-D7 result:

```text
D0–D6 semantic reopen       NONE
D7-A process topology       PRESERVED
D7-B core RLS/transaction   PRESERVED + auth-bootstrap taxonomy repair
D7-C River/effect model     PRESERVED + post-restore fence qualification
D7-D auth profile           PRESERVED + ambiguous-carrier fail-closed clarification
D7-E operability profile    PRESERVED + River migration/recovery/auth-validator seams
Product surface             99 operations / 30 Permissions / H-A-S unchanged
D8 / D9                     NOT OPEN
Product implementation      BLOCKED
```

No microservices, Redis, broker, workflow engine, new business authority, Product operation or Permission is introduced.

## 8. Required next review

This bounded repair remains **candidate** until the independent whole-D7 Fable challenge and GPT adjudication complete.

The independent reviewer must specifically try to falsify:

- whether `authentication_bootstrap` leaks into a generic cross-tenant bypass;
- whether River migration ownership is still duplicated or incomplete;
- whether post-restore external-effect ambiguity can still cause redispatch;
- whether OAD validation can accidentally bypass D7-D carrier separation;
- whether proof-timing language still claims an unimplemented runtime PASS;
- and whether these repairs introduce overengineering relative to the accepted Product 1.0 consumer set.

Do not begin D8, D9 or Product implementation.