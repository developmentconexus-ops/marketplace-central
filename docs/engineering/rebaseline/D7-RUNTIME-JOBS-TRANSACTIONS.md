# D7 — Runtime / Jobs / Transactions

> **Status:** OPEN / ACTIVE — D7-A→D7-E OPERATOR-RATIFIED / WHOLE-D7 REVIEW CONVERGED / OPERATOR CLOSEOUT PENDING  
> **Program:** Architecture Rebaseline / Technical System Design  
> **Opened:** 2026-08-21  
> **Parent authorities:** accepted D0–D6 semantics, `ARCHITECTURE.md`, canonical Product OAD, and bounded owner authority routed by `docs/index.md`  
> **Method:** DevelopmentConexus Engineering Method v1.0.0

## 1. Purpose and boundary

D7 defines the smallest target runtime realization capable of executing the already-accepted Marketplace Central Product and frontend contracts without creating new business authority.

D7 owns server/process topology, PostgreSQL isolation/transactions, durable work/effects, session/CSRF/OIDC and machine-token realization, Product HTTP runtime/validation, private byte custody, secrets/configuration, observability, migrations, deployment/health/backup and real-dependency proof seams.

D7 does **not** reopen Product operations, ordinary Permissions, Principal kinds, semantic ownership, frontend interaction meaning or provider/business-system semantics by implementation convenience. D8 golden-flow choreography, D9 adversarial review and Product implementation remain blocked until the roadmap explicitly advances them.

## 2. Imported invariants

1. Go is canonical backend business execution.
2. PostgreSQL is canonical MPC-owned state storage.
3. Organization isolation is structural, never remembered-predicate-only.
4. Product OpenAPI is the one machine-readable Product wire authority.
5. Q/C/E/P meaning remains semantic-owner authority.
6. Consequential writes preserve idempotency, ambiguity, auditability and reconciliation; possible acceptance is never blindly replayed.
7. Unknown/partial/unavailable remain honest.
8. Provider protocol stays inside adapters; Sankhya target integration remains sanctioned API Gateway only.
9. H browser uses server-side session + CSRF and never browser-held OIDC tokens.
10. A/S remain audience-bound Client Credentials bearers and never become authorization shortcuts.
11. Product implementation remains blocked until accepted D9.

## 3. D7 target invariant

> **Every admitted Product or Technical-Ingress interaction can be realized through explicit serving, persistence, transaction, durable-work, authentication and operability boundaries that preserve Organization isolation, owner authority, wire-contract conformity, request-trust separation, idempotency/concurrency, recoverable external-effect semantics and secret/PII safety while platform mechanisms remain business-policy-free.**

## 4. Decision surface and owners

| Slice | Owner | Status |
| --- | --- | --- |
| D7-A — Runtime Envelope & Transaction Ownership | this document | **OPERATOR-RATIFIED** |
| D7-B — PostgreSQL Isolation & Transactions | [D7-B](D7-B-POSTGRESQL-ISOLATION-TRANSACTIONS.md) | **OPERATOR-RATIFIED** |
| D7-C — Durable Work & External Effects | [D7-C](D7-C-DURABLE-WORK-EXTERNAL-EFFECTS.md) | **OPERATOR-RATIFIED** |
| D7-D — Authentication / Session / CSRF / Machine Tokens | [D7-D](D7-D-AUTHENTICATION-SESSION-CSRF.md) | **OPERATOR-RATIFIED** |
| D7-E — Operability / Secrets / Migrations / Deployment & Proof | [D7-E](D7-E-OPERABILITY-DEPLOYMENT-PROOF.md) | **OPERATOR-RATIFIED** |
| D7-R1 — Whole-Stage Coherence Corrections | [D7-R1](D7-R1-WHOLE-STAGE-COHERENCE.md) | **FABLE REVIEWED / GPT ADJUDICATED — BOUNDED FIX APPLIED** |

All five realization slices are ratified. Whole-D7 review found no D0–D6 contradiction and no need to reconstruct D7. D7-R1 contains only the bounded composition repairs that survive review.

## 5. Accepted D7-A — Runtime Envelope & Transaction Ownership

One Go application process per replica serves one public application origin and contains:

```text
same-origin frontend/static delivery
Product API boundary
Technical Non-Product Ingress boundary
H session / CSRF / OIDC mediation
in-process durable worker runner
scheduler seam
PostgreSQL pool
```

One process does not merge semantic owners. Internal Q/C/P enter explicit owner application boundaries rather than self-HTTP. E is durable only where D3 admits an independent reaction.

Owner consequential intake uses one owner-local PostgreSQL transaction for canonical owner state plus required technical correctness records. External network writes happen only after commit. No distributed business transaction spans semantic owners.

API/worker process split reopens only for measured resource, scaling, failure, security or deployment need.

## 6. Accepted D7-B — PostgreSQL Isolation & Transactions

[D7-B](D7-B-POSTGRESQL-ISOLATION-TRANSACTIONS.md) selects:

```text
organization_id on organization-owned state/evidence
+ composite Organization foreign keys
+ transaction-local scope
+ PostgreSQL ENABLE + FORCE RLS
+ runtime role = non-owner / NOSUPERUSER / NOBYPASSRLS
+ READ COMMITTED baseline + explicit row locking
+ opaque random owner revision tokens
+ organization/operation-scoped idempotency
+ pgx/v5 + pgxpool
```

D7-R1 completes the persistence taxonomy with a narrow `authentication_bootstrap` mode for one-time login transaction, ApplicationSession lookup and machine-client binding before Principal resolution. It cannot read Organization-owned business/evidence state. River engine tables are explicit platform technical/library state and do not grant business authority.

No generic ORM/repository abstraction, database/schema/role-per-Organization or global `SERIALIZABLE` baseline is admitted.

## 7. Accepted D7-C — Durable Work & External Effects

[D7-C](D7-C-DURABLE-WORK-EXTERNAL-EFFECTS.md) selects River over the accepted pgx/PostgreSQL stack:

```text
InsertTx atomic owner-state -> durable-work handoff
no second generic MPC outbox for River work
repeat-safe delivery / no exactly-once claim
owner semantic idempotency = correctness
persisted pre-dispatch marker before consequential external writes
possible acceptance => no redispatch
source-authoritative reconciliation
scheduler/job state = wake-up/technical state only
```

River completion, retry, uniqueness, rescue and schedule state never become business truth/history. Crash/timeout after possible dispatch moves to reconciliation, not generic retry.

D7-R1 qualifies continuous-timeline crash semantics for database rollback: after a restore, absence of a marker is not proof that no external effect occurred. Consequential dispatch remains behind the recovery fence until timeline continuity or reconciliation proves safety.

## 8. Accepted D7-D — Authentication / Session / CSRF / Machine Tokens

[D7-D](D7-D-AUTHENTICATION-SESSION-CSRF.md) selects:

```text
Keycloak first OIDC/OAuth provider
H: confidential Authorization Code + PKCE S256
   go-oidc + x/oauth2
   verified (issuer, sub) -> fresh opaque MPC session
   human OIDC token set discarded after callback
   PostgreSQL stores session-handle digest only
   30m idle / 8h absolute baseline
   X-CSRF-Token synchronizer + net/http.CrossOriginProtection
A/S: Keycloak Client Credentials
     audience includes https://conexus.fun
     jwx/v3 trusted-JWKS verification
     explicit machine-client -> A/S Principal binding
```

No browser bearer, persistent human refresh-token cache, Redis session store, IdP role→Permission mapping, realm-per-Organization, wildcard CORS or generic IAM engine is admitted.

D7-R1 clarifies pre-Principal persistence scope and fail-closed multiple-carrier/OAD-validator composition.

## 9. Accepted D7-E — Operability / Deployment / Proof Baseline

[D7-E](D7-E-OPERABILITY-DEPLOYMENT-PROOF.md) selects:

```text
HTTP       Chi v5 + D5-pinned oapi-codegen v2.8 strict server
validation oapi-codegen/nethttp-middleware
policy     generated OAD operation-policy metadata; no handwritten duplicate map
bytes      private S3 API-compatible object storage; authenticated Go delivery; no CDN baseline
config     typed startup config + deployment-injected env/file secrets
migrate    tern/v2 for MPC schema + version-matched River migration tool for River schema
logs       JSON log/slog
telemetry  OpenTelemetry traces + metrics over OTLP/HTTP; logs remain slog
artifact   one immutable OCI application image with Go + compiled frontend
edge       https://conexus.fun through explicit trusted TLS proxy/ingress boundary
backup     PostgreSQL base backup + WAL/PITR or managed equivalent; restore proof
proof      real PostgreSQL/River/Keycloak/browser/router/object-store seams
```

No Kubernetes/service mesh, Redis, external broker, mandatory Collector/Prometheus/ELK stack, Vault dependency, hot configuration, ORM auto-schema, startup auto-migration, CDN/public bucket, multi-region or generic IaC platform is introduced.

## 10. Accepted D7-R1 bounded composition repairs

[D7-R1](D7-R1-WHOLE-STAGE-COHERENCE.md) preserves D7-A→D7-E while closing only five composition seams:

1. `authentication_bootstrap` persistence scope before Principal resolution;
2. separate migration ownership: `tern` for MPC schema, exact-version River tooling for River schema;
3. PITR/database-time rollback recovery fence with **affirmative external continuity witness** and automatic fail-closed arming when continuity is absent/unverifiable;
4. scheme-aware OAD validator composition over already-established D7-D carrier context; no production no-op authentication function and no implicit dual-carrier priority;
5. proof-timing distinction between D7 architecture closeout and later implemented real-dependency conformance.

Independent Fable challenge on exact candidate `c08a4d025cfd89269cc071f4b307695e79f6f8cb` returned **ACCEPT WITH BOUNDED FIXES**. GPT accepted only the Important continuity/arming completion. Fable's `chi-server` generation suggestion and bootstrap-budget note remain non-blocking and require no D7 authority change.

Round 2 is not justified because the blocking amendment adds no new technology, Product meaning or runtime topology.

## 11. Whole-D7 proof contract

D7 closeout leaves executable falsifiers for at least:

- cross-Organization DB access/reference bypass and pooled-scope leakage;
- authentication-bootstrap escape into business/evidence state;
- stale revision/idempotency divergence;
- owner state/handoff atomicity and duplicate/rescued/out-of-order work;
- ambiguous external effect redispatch, including an ordinary boot after PITR with no manual fence-arming step;
- automatic fail-closed continuity detection when restored database lineage cannot be positively established;
- Sankhya Direct Oracle fallback;
- browser token exposure, OIDC replay/fixation, CSRF/cross-origin bypass and disabled-Principal session use;
- wrong machine issuer/audience/JWKS/client binding and A/S→H confusion;
- no-op/mismatched OAD security validation and dual-carrier ambiguity;
- `{id}:verb` runtime dispatch, OAD-invalid request rejection and technical-surface exclusion;
- unauthorized private byte delivery and object/DB recovery integrity;
- secret/PII leakage through logs/traces/metrics/jobs;
- MPC/River migration skew, runtime migration-owner privilege and incompatible-schema boot;
- trusted-proxy spoofing, provider false-death and telemetry coupling;
- restore loss of RLS/history/effect safety, IdP subject continuity or committed binary integrity.

D7 architecture closeout does **not** claim an implemented Product runtime PASS. The current repository intentionally has active runtime population `NONE`. Real PostgreSQL/River/Keycloak/browser/router-validator/object-store execution becomes mandatory for implementation acceptance when the implementation gate opens after D9; mock-only tests cannot substitute for those claims.

## 12. Whole-D7 review result

```text
internal coherence review        COMPLETE
independent Fable challenge      ACCEPT WITH BOUNDED FIXES
GPT adjudication                 CONVERGED
surviving Important finding      PITR continuity/automatic fence arming — APPLIED
fresh post-fix repository gate   PASS @ 42c208abdc320538f9222ccb1ccc4d09705f6577
D0–D6 reopen                     NONE
D7 reconstruction               NONE
Product                          99 operations / 30 Permissions / H-A-S unchanged
D8 / D9                          BLOCKED
Product implementation           BLOCKED
```

## 13. Exact next action

**Operator adjudicate/ratify D7 closeout.** Until that explicit decision, D7 remains OPEN / ACTIVE, PR #58 remains unmerged, D8/D9 remain blocked and Product implementation remains blocked.