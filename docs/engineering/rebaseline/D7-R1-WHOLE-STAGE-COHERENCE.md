# D7-R1 — Whole-Stage Coherence Corrections & Proof-Timing Contract

> **Status:** OPERATOR-RATIFIED / ACCEPTED — WHOLE-D7 CLOSEOUT COMPLETE  
> **Parent:** `D7-RUNTIME-JOBS-TRANSACTIONS.md`  
> **Scope:** exact D7-B/C/D/E realization seams only; no D0–D6 semantic reopen  
> **Method:** DevelopmentConexus Engineering Method v1.0.0  
> **Derived:** 2026-08-22  
> **Ratified:** 2026-08-22  
> **Independent Fable verdict:** ACCEPT WITH BOUNDED FIXES  
> **GPT adjudication:** CONVERGED — Fable F-1 accepted; F-2/F-3 non-blocking

## 1. Purpose

Whole-D7 cross-check found a small set of realization seams that become contradictory only when D7-A→D7-E are composed as one runtime. This amendment is the smallest bounded repair.

It does **not** change Product operations, Permissions, Principal kinds, owner semantics, D3 Q/C/E/P meaning, D4 external-effect meaning, D5 wire/auth contract or D6 frontend authority.

This document supersedes only the exact D7 clauses named below.

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

Any disaster/PITR restore that may have lost acknowledged MPC commits enters a fail-closed **recovery fence** before normal external-dispatch authority resumes:

1. identify the restored recovery point and possible state-loss window;
2. invalidate restored human ApplicationSessions and one-time login transactions so a lost post-backup logout/revocation cannot resurrect a browser session;
3. keep consequential external-write dispatch disabled, including restored River external-effect jobs;
4. treat restored dispatchable/ambiguous work as reconciliation-only until its owner/source scope is cleared;
5. reacquire authoritative provider/business-system state for affected Installations/SourceInstances and reconcile restored owner intents/effect anchors over the possible loss window;
6. revalidate current identity/access and machine-client configuration before normal authenticated automation resumes when the restore could have rolled those changes back;
7. verify committed binary references at the restored database point against object-store digest/size/type evidence;
8. release the fence only for scopes whose reconciliation/access/integrity checks are sufficient. If coverage cannot establish safety, automated external writes remain disabled and the condition is surfaced for operator resolution rather than guessed away.

New durable intake may be recorded during a bounded recovery posture only if its external dispatch remains behind the same recovery fence. Safe reads may remain available where their accepted semantics are still honest; the fence is not permission to create a false whole-application outage.

A zero-data-loss failover whose durable database lineage is affirmatively proved continuous need not enter the rollback recovery path.

### Recovery-fence arming / continuity law — Fable F-1 accepted

**Timeline continuity is established affirmatively; it is never assumed from a syntactically healthy restored database.**

Before ordinary boot may enable consequential external dispatch or trust restored authentication state as continuous, the runtime/deployment boundary must verify a **continuity witness outside the PostgreSQL rollback domain**. The witness must carry enough deployment/database lineage and durable-position/recovery-epoch evidence to distinguish a continuous current database from a database restored to an earlier acknowledged state. Examples include PostgreSQL system/timeline/durable-position evidence anchored in already-selected object/deployment state, or an equivalent deployment-provider continuity control. Exact representation is implementation detail.

Binding laws:

- absence, mismatch, stale/unverifiable lineage, or inability to prove that the current database safely descends from the last externally witnessed state **arms the recovery fence by default**;
- ordinary application boot after a PITR/rollback must fail closed for external dispatch **without requiring an operator to remember a manual `recovery=true` flag or runbook step**;
- the continuity protocol must cover every commit class whose rollback could re-enable an already-possible external effect or resurrect authentication/authority that was subsequently revoked. A casual periodic heartbeat that can lag those safety-sensitive transitions is insufficient;
- a sufficient realization may advance the out-of-rollback-domain witness around safety-sensitive checkpoints or use equivalent provider/deployment lineage guarantees, but D7 does not freeze that implementation here;
- failure to update/verify the witness cannot silently preserve normal dispatch authority;
- continuity evidence is technical safety state only. It does not become business history, Product identity, provider truth or a second authorization authority;
- false-positive fencing is acceptable and recoverable; false-negative continuity is not.

This law closes the unarmed-restore path identified by independent Fable review without introducing a new business mechanism or reopening D4/D7-C effect semantics.

### Required recovery falsifiers

The eventual real-dependency implementation proof must include both:

1. restore a database snapshot from **before** a simulated externally accepted write while the external system retains the effect; prove restored River/owner state does not redispatch before reconciliation/fence release;
2. perform an ordinary application boot after that restore with **no manual fence-arming action**; prove the external continuity check detects/assumes unsafe rollback and automatically keeps dispatch fenced. A proof that starts with the fence pre-engaged is insufficient.

A complementary continuous-lineage failover probe should prove the fence can remain open only when continuity is positively established.

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

## 7. Independent Fable review and GPT adjudication

The isolated `review/d7-fable` challenge reviewed exact candidate:

```text
stage/d7-runtime @ c08a4d025cfd89269cc071f4b307695e79f6f8cb
```

with correct single-file review isolation and green exact-candidate CI.

Fable verdict:

```text
ACCEPT WITH BOUNDED FIXES
```

Material adjudication:

| Fable finding | GPT disposition | D7 action |
| --- | --- | --- |
| F-1 — recovery fence lacks affirmative detection/arming law | **ACCEPT** | incorporated in §4 continuity law + unarmed-restore falsifier |
| F-2 — selected `chi-server` generation mode not yet executed against canonical OAD | **VALID / NON-BLOCKING** | no authority change; optional cheap future gate hardening, not required for D7 closeout |
| F-3 — bootstrap budget at 93% | **VALID / REPOSITORY HYGIENE / NON-BLOCKING** | no cap increase; compress CLOSED D7 roadmap detail when D8 opens if needed |

Fable independently found no reason to reopen D0–D6 or reconstruct D7-A→D7-E, and confirmed 99 Product operations / 30 ordinary Permissions / H-A-S only, provider isolation, Sankhya API-Gateway-only and Product implementation blocked until D9.

Round 2 is not justified: the sole Important finding is a bounded completion of the already-owned R1-F3 recovery law and introduces no new architecture selection.

## 8. Whole-D7 closeout disposition

Operator ratification closes the whole-stage review with the following accepted result:

```text
D0–D6 semantic reopen       NONE
D7-A process topology       OPERATOR-RATIFIED
D7-B RLS/transactions       OPERATOR-RATIFIED + R1 auth-bootstrap completion
D7-C River/effect model     OPERATOR-RATIFIED + R1 recovery-fence qualification
D7-D auth profile           OPERATOR-RATIFIED + R1 dual-carrier/OAD-validator clarification
D7-E operability profile    OPERATOR-RATIFIED + R1 migration/recovery seams
D7-R1                       OPERATOR-RATIFIED / ACCEPTED
Product surface             99 operations / 30 Permissions / H-A-S unchanged
D8 / D9                     NOT OPEN ON THIS UNMERGED BRANCH
Product implementation      BLOCKED
```

No microservices, Redis, broker, workflow engine, new business authority, Product operation or Permission is introduced.

D7-R1 is now accepted authority. The remaining repository action is integration of the D7 closeout into `main`; this document does not itself authorize that merge.
