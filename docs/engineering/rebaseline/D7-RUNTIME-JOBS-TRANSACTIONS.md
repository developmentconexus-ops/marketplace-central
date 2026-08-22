# D7 — Runtime / Jobs / Transactions

> **Status:** OPEN / ACTIVE — D7-A→D7-D OPERATOR-RATIFIED / D7-E CANDIDATE PENDING RATIFICATION  
> **Program:** Architecture Rebaseline / Technical System Design  
> **Opened:** 2026-08-21  
> **Parent authorities:** accepted D0–D6 semantics, `ARCHITECTURE.md`, canonical Product OAD, and bounded owner authority routed by `docs/index.md`  
> **Method:** DevelopmentConexus Engineering Method v1.0.0

## 1. Purpose and boundary

D7 defines the smallest target runtime realization capable of executing the already-accepted Marketplace Central Product and frontend contracts without creating new business authority.

D7 owns server/process topology, PostgreSQL isolation/transactions, durable work/effects, session/CSRF/OIDC and machine-token realization, Product HTTP runtime/validation, private byte custody, secrets/configuration, observability, migrations, deployment/health/backup and real-dependency proof seams.

D7 does **not** reopen Product operations, ordinary Permissions, Principal kinds, semantic ownership, frontend interaction meaning or provider/business-system semantics by implementation convenience. D8 golden-flow choreography, D9 adversarial review and Product implementation remain blocked.

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
| D7-E — Operability / Secrets / Migrations / Deployment & Proof | [D7-E](D7-E-OPERABILITY-DEPLOYMENT-PROOF.md) | **CANDIDATE / OPERATOR RATIFICATION PENDING** |

After D7-E ratification, D7 still requires one whole-stage coherence + proof/adversarial review before closeout. D8 does not open merely because the five slices exist.

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

Principal-self and technical-routing scope are bounded bootstrap modes, never generic cross-tenant business access. No generic ORM/repository abstraction, database/schema/role-per-Organization or global `SERIALIZABLE` baseline is admitted.

Real PostgreSQL negative proof remains mandatory before D7 closeout.

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

Real PostgreSQL + River crash/restart, transactional enqueue, duplicate/rescue and cross-Organization proof remains mandatory.

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
   X-CSRF-Token synchronizer + net/http CrossOriginProtection
A/S: Keycloak Client Credentials
     audience includes https://conexus.fun
     jwx/v3 trusted-JWKS verification
     explicit machine-client -> A/S Principal binding
```

No browser bearer, persistent human refresh-token cache, Redis session store, IdP role→Permission mapping, realm-per-Organization, wildcard CORS or generic IAM engine is admitted.

Real Keycloak/OIDC plus browser-capable cookie/CSRF proof remains mandatory.

## 9. Current D7-E candidate

[D7-E](D7-E-OPERABILITY-DEPLOYMENT-PROOF.md) proposes the remaining smallest runtime/operability profile:

```text
HTTP       Chi v5 + D5-pinned oapi-codegen v2.8 strict server
validation oapi-codegen/nethttp-middleware
policy     generated OAD operation-policy metadata; no handwritten duplicate map
bytes      private S3 API-compatible object storage; authenticated Go delivery; no CDN baseline
config     typed startup config + deployment-injected env/file secrets
migrate    tern/v2 via separate migration owner / release step
logs       JSON log/slog
telemetry  OpenTelemetry traces + metrics over OTLP/HTTP; logs remain slog
artifact   one immutable OCI application image with Go + compiled frontend
edge       https://conexus.fun through explicit trusted TLS proxy/ingress boundary
backup     PostgreSQL base backup + WAL/PITR or managed equivalent; restore proof
proof      real PostgreSQL/River/Keycloak/browser/router/object-store seams
```

The candidate also closes the D5 `{id}:verb` routing constraint, runtime OAD schema validation, private authored-media byte custody, schema-version boot refusal, dependency-aware readiness and Keycloak subject-continuity/binary-integrity recovery.

No Kubernetes/service mesh, Redis, external broker, mandatory Collector/Prometheus/ELK stack, Vault dependency, hot configuration, ORM auto-schema, startup auto-migration, CDN/public bucket, multi-region or generic IaC platform is introduced.

## 10. Whole-D7 proof obligations

Before D7 closes, the combined authority/proof plan must make at least these properties falsifiable:

- cross-Organization DB access/reference bypass;
- transaction-scope leakage through pooled connections;
- stale revision/idempotency divergence;
- owner commit without required durable job or rolled-back owner state with runnable job;
- duplicate/rescued/out-of-order work corrupting business truth;
- ambiguous external effect being blindly redispatched;
- Sankhya recovery falling back to Direct Oracle;
- browser access/refresh token exposure;
- invalid/replayed OIDC state/nonce/issuer/audience creating a session;
- unsafe H request escaping CSRF/cross-origin controls;
- A/S bearer resolving to H or gaining Permission from IdP roles;
- Product `{id}:verb` paths failing selected runtime dispatch;
- runtime request validation accepting OAD-invalid shape;
- technical routes leaking into Product OAD/SDK;
- unauthorized private byte delivery;
- secret/PII leakage through logs/traces/metrics/jobs;
- incompatible schema booting successfully or runtime receiving migration-owner power;
- untrusted proxy headers influencing auth/public-origin security;
- provider/telemetry degradation falsely killing application readiness;
- backup/restore losing RLS/history/effect safety, stable IdP subject continuity or committed binary integrity.

A mock-only green result cannot close a claim whose subject is PostgreSQL, River, Keycloak/OIDC, browser cookie/CSRF, HTTP router/OAD validator or object storage.

## 11. Exact next action

**Operator adjudicate/ratify D7-E. Do not begin whole-D7 closeout review before that decision.**

After D7-E ratification, run one combined D7 coherence + executable-proof/adversarial review over D7-A→D7-E. D8, D9 and Product implementation remain blocked.