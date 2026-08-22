# D7-E — Operability / Secrets / Migrations / Deployment & Proof Baseline

> **Status:** CANDIDATE / OPERATOR RATIFICATION PENDING  
> **Parent:** `D7-RUNTIME-JOBS-TRANSACTIONS.md`  
> **Accepted prerequisites:** D7-A + D7-B + D7-C + D7-D — OPERATOR-RATIFIED  
> **Parent authorities:** accepted D5 Product OAD/tooling + accepted D7 runtime/persistence/work/auth realization  
> **Method:** DevelopmentConexus Engineering Method v1.0.0  
> **Derived:** 2026-08-22

## 1. Purpose

D7-E closes the smallest remaining runtime/operability decisions required to make the accepted Marketplace Central architecture implementable and falsifiable after D9 without introducing a platform project.

It owns only:

- concrete Product HTTP routing/server-generation and runtime request-validation realization;
- structured logging plus bounded traces/metrics export;
- configuration/secret injection and rotation boundary;
- database migration execution boundary;
- deployable/container and same-origin edge topology;
- private binary/blob storage and authenticated byte-delivery mechanics already deferred by D5;
- health/readiness/shutdown operational semantics;
- backup/restore expectations for state whose loss would invalidate accepted identity/history;
- the real-dependency proof baseline needed by D7 closeout and later D8/D9.

D7-E does not create Product operations, monitoring Product features, a generic platform/IAM/storage domain, Kubernetes/service-mesh authority, a data warehouse, or Product implementation.

## 2. Target invariant

> **The accepted Product OAD generates the runtime HTTP contract instead of being retyped by hand; one immutable application artifact can start only against compatible configuration/schema/trust state; runtime secrets and PII do not leak into telemetry; durable state and required binary evidence can be restored; external telemetry/marketplace/Sankhya outages do not become false application death; and every D7 safety claim has a named real-dependency falsifier before production implementation is accepted.**

## 3. Product HTTP realization

### 3.1 `go-chi/chi/v5` — SELECT CANDIDATE

Use `github.com/go-chi/chi/v5` as the Product/technical HTTP router on top of `net/http`.

Reason:

- D5 already proved that plain `net/http.ServeMux` is not a neutral realization for canonical partial-segment paths such as `/{id}:submit`;
- Chi stays `net/http`-compatible rather than introducing a larger web framework/runtime model;
- current Chi's radix-tree parser explicitly supports a parameter node with a non-`/` tail delimiter followed by a static suffix, which directly models `{id}:verb` without changing W1 path grammar;
- route groups/middleware stacks are sufficient to keep Product, Technical Ingress, auth/session and technical presentation surfaces distinct inside the one D7-A process.

Exact Chi minor/patch is implementation-manifest authority; the major line is v5. A release upgrade must retain the partial-segment dispatch negative proof.

Echo/Gin/Fiber/Iris/Gorilla are not admitted merely because `oapi-codegen` supports them. They provide no current property absent from Chi for MPC.

### 3.2 `oapi-codegen v2.8.0` strict Chi server — SELECT CANDIDATE

D5 already pins `oapi-codegen v2.8.0` for the canonical Product OAD Go projection. D7-E reuses that same generator rather than adding a second server generator.

Runtime generation uses the canonical resolved Product OAD with:

```text
chi-server:    true
strict-server: true
```

The generated strict boundary owns typed HTTP request/response transport shapes only. Handlers/adapters translate generated wire shapes into owner application calls; generated code never becomes business authority.

No second hand-authored path table, Product DTO catalog or Permission map is allowed.

### 3.3 Generated operation-policy projection

The OAD's accepted per-operation extensions remain the single source for runtime transport enforcement metadata:

```text
operationId
x-mpc-operation-class
x-mpc-required-permission
x-mpc-principal-kinds
x-mpc-semantic-owner
x-mpc-required-physical-qualification when present
```

A deterministic generated projection may emit a small immutable runtime lookup keyed by `operationId`. Runtime carrier/authz middleware consumes that generated metadata; developers do not maintain a parallel handwritten Permission/Principal-kind switch.

This is derived projection only. D2/D5 authority still decides access semantics.

### 3.4 Request validation — `oapi-codegen/nethttp-middleware`

Use `github.com/oapi-codegen/nethttp-middleware` against the canonical resolved/embedded Product specification to enforce the OAD request contract at runtime, including path/query/header/body schema constraints that generated Go types alone do not prove.

Binding laws:

- Product authentication remains the carrier-specific D7-D implementation; the validator cannot merge H cookie and A/S bearer semantics into one generic bearer path;
- W4/current Membership/Permission/business/Governance checks remain MPC authority after authentication;
- validator failures are translated into the accepted D5 Product HTTP/Problem grammar rather than leaking middleware/kin-openapi implementation errors;
- OAD `servers` remain absent, so validator configuration cannot invent environment host authority;
- runtime validation must prove `additionalProperties: false`, required fields, patterns, media types, method routing and canonical `405 + Allow` behavior where applicable;
- strict-server response types are compile-time transport guards; D7-E does not add a response-validation proxy on every production response absent a proved need.

Technical Ingress and authored-media delivery are not added to the Product OAD. They use their own protocol-local routing/validation boundaries.

## 4. Private binary/blob storage

D5 already admits authored ListingIntent media and bounded fulfillment/artifact byte delivery while explicitly deferring physical storage/delivery mechanics to D7.

### 4.1 Private S3-compatible object storage — SELECT CANDIDATE

Use one private S3 API-compatible object-storage bucket/namespace for MPC-owned binary content whose bytes are not appropriate as PostgreSQL relational state.

Examples currently justified:

- ListingIntent-authored media bytes;
- bounded Fulfillment artifacts where accepted owner semantics require retained binary evidence.

PostgreSQL remains authority for semantic metadata, identity, Organization scope and historical references. Object storage is byte custody only and creates no generic Asset/Media domain.

Laws:

- bucket/object namespace is private; public-read ACL is forbidden;
- raw storage keys/URLs never become Product identity or persisted presentation locators;
- no durable anonymous/presigned browser locator is baseline;
- authenticated authored-media delivery remains a Go technical presentation route that rechecks current Product H/A/S authentication, Principal eligibility, exact Organization Membership and `offering.read` for the ListingIntent/media relation before streaming/proxying bytes;
- no CDN is baseline; add one only if measured delivery need proves it while preserving the D5 authorization fence;
- provider/source media remain distinct external evidence and are not copied into this store merely for normalization.

Exact object-storage vendor is deployment configuration. The implementation must pin one tested S3 API behavior and cannot silently assume every "S3-compatible" edge case is identical.

### 4.2 Database/object-store consistency

There is no distributed PostgreSQL↔object-store transaction claim.

For a new authored binary, the safe baseline is:

```text
validate/stream private object first
  -> durable object write succeeds
  -> owner PostgreSQL transaction establishes idempotency/revision/media descriptor/reference
  -> COMMIT
```

A database commit may never create the first authoritative reference to bytes that have not been durably stored.

If object write succeeds and the owner transaction later rolls back/rejects, the object is an unreferenced technical orphan. A bounded recovery sweep deletes old unreferenced objects. Queue/sweep state is not media business truth.

Object overwrite of a committed key with materially different bytes is forbidden. Content digest/size/type needed for integrity is retained in owner metadata proportionately.

## 5. Configuration and secrets

### 5.1 Minimal configuration model

Use one typed startup configuration assembled from deployment-supplied environment values and/or mounted secret files. No generic configuration server or dynamic configuration framework is baseline.

Startup validates all required configuration before serving traffic. Unknown/malformed critical configuration fails startup.

Configuration classes remain explicit, e.g.:

```text
public/runtime configuration
secret-bearing configuration
external endpoint/trust configuration
operational tuning
```

Command-line flags are not used for secrets because process arguments are routinely inspectable. Production `.env` files are not committed or treated as a secret store.

### 5.2 Secret injection

Secrets are injected by the deployment environment as protected environment values or mounted files with least-privilege filesystem/process access. Prefer file-mounted secrets where the deployment platform supports clean rotation; no HashiCorp Vault/cloud-secret-manager dependency is mandatory until a real deployment environment requires it.

Baseline rotation is explicit redeploy/restart with overlap where the external system supports two credentials. Hot secret reload is not required by preference.

At minimum, the application secret-redaction fence covers:

- database credentials;
- Keycloak confidential-client secret;
- provider/business-system credentials;
- MPC session/CSRF material;
- OAuth state/nonce/code/verifier/token material;
- object-storage credentials;
- OTLP authentication headers/tokens if used.

No secret value is a structured log attribute, trace attribute, metric label, Problem detail or River job arg.

## 6. Database migrations — `tern/v2`

Select `github.com/jackc/tern/v2` as the PostgreSQL migration mechanism.

Why:

- PostgreSQL is fixed architecture, so SQL-first migrations are honest;
- tern is a small standalone PostgreSQL migrator from the same pgx ecosystem;
- it avoids an ORM migration model/second schema authority;
- ordered SQL files remain reviewable alongside RLS policies, composite FKs, roles and constraints that are correctness-critical in D7-B.

Migration execution boundary:

```text
release/deploy step
  -> migration credential / schema owner
  -> tern migrate
  -> verify expected schema version
  -> start/roll application runtime using non-owner runtime credential
```

The application process does **not** auto-run migrations at normal startup and never receives the migration-owner credential.

Application startup/readiness fails on an incompatible schema version. It does not silently continue against an older/newer shape.

Production automatic `down` migration is not a rollback strategy. Irreversible/destructive changes require explicit operator review, backup/restore readiness and a forward/restore plan. Initial clean-baseline implementation does not pre-build a universal expand/contract framework; compatibility choreography is added when a real rolling-deployment consumer requires it.

`sqlc` remains an implementation-time query-generation option once real SQL query files exist. D7 does not need it as architecture authority.

## 7. Observability

### 7.1 Logs — standard-library `log/slog`

Use `log/slog` with JSON output to stdout/stderr as the canonical runtime structured-log API.

Logs carry bounded operational facts such as severity, component, operationId, request/trace correlation, outcome class and latency. They do not become audit/business history.

High-cardinality or sensitive identifiers are included only when a named diagnostic need justifies them. Provider payloads, buyer/customer PII, tokens, cookies, CSRF values and secrets are never logged.

No Zap/Zerolog/Logrus dependency is added absent a property `slog` cannot meet.

### 7.2 Traces and metrics — OpenTelemetry stable signals

Use OpenTelemetry Go API/SDK for **traces and metrics**, exported through OTLP to a deployment-configured endpoint.

Current OpenTelemetry Go status marks traces and metrics stable while logs remain less mature. Therefore application logs stay on `slog`; D7-E does not add the OTel logs pipeline merely for symmetry.

Use OTLP/HTTP exporters as the default transport because they require only ordinary HTTPS infrastructure:

```text
go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp
go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp
```

A Collector may sit behind that endpoint, but an OpenTelemetry Collector is not a mandatory MPC service. The app is backend-neutral at the OTLP boundary.

Telemetry export failure does not make Product readiness false. Export is bounded/batched and must not block owner transactions or external-effect safety.

Baseline telemetry dimensions avoid Organization/Principal/native provider IDs as metric labels. Useful bounded measurements include HTTP latency/status/operationId, pool saturation, River queue/execution health, auth/session rejection classes, and external dependency latency/outcome class without raw error payloads.

## 8. Deployable and same-origin topology

### 8.1 One immutable OCI application image

Produce one immutable OCI-compatible application image per release containing:

- the compiled Go application;
- the compiled D6 frontend static assets;
- only runtime CA/timezone/assets actually required.

Use a multi-stage build so Node/Go build toolchains do not ship in the final runtime image.

The image runs the single D7-A Go process. River workers are goroutines inside that service, not a second process/service by default.

Image tag alone is not release identity; deployment pins an immutable digest. Do not deploy floating `latest`.

### 8.2 Public edge

Canonical public origin remains:

```text
https://conexus.fun
```

TLS/certificate termination may be supplied by the deployment ingress/reverse proxy/load balancer. Exact Caddy/Nginx/cloud ingress vendor is not architecture authority.

The application trusts forwarded scheme/host/client-address headers only from an explicitly configured trusted proxy boundary. Untrusted forwarded headers cannot influence OIDC redirect construction, security origin checks or audit attribution.

Frontend/static assets, Product API and H auth/session technical routes remain same-origin. Technical provider ingress may share the host under disjoint technical paths when its protocol permits; it never enters Product OAD vocabulary.

No Kubernetes, service mesh or ingress controller platform is baseline. One VPS/VM/container host or a simple container service can satisfy the topology. Horizontal replicas remain allowed once shared PostgreSQL/object storage/session/River semantics are proved.

## 9. Health, readiness and shutdown

Expose minimal technical health surfaces for deployment automation. Exact path spelling is technical implementation detail but they remain outside Product OAD.

**Liveness** answers only whether the process/event loop is alive enough to respond. It does not synchronously call external systems.

**Readiness** is true only when the local runtime can safely accept new MPC work, including proportionately:

- validated startup configuration;
- PostgreSQL reachable under runtime credential;
- expected schema/migration version;
- required auth trust configuration initialized;
- River worker/runtime initialization healthy for admitted work classes;
- private object-store configuration syntactically/trust-valid where binary features are enabled.

Transient Mercado Livre, Sankhya, Keycloak login endpoint, object-store request, or OTLP backend outage after startup does not automatically make the whole application unready when accepted Product semantics can represent that dependency as unavailable/pending and existing authenticated work can continue safely.

A dependency may drive readiness false only when its absence makes accepting **all** new requests unsafe, not merely because it is degraded.

Graceful shutdown follows accepted D7-A order: mark not-ready; stop new HTTP and job claiming; bounded drain; leave uncompleted River work durable; flush telemetry best-effort; close resources.

## 10. Backup and restore

### 10.1 MPC PostgreSQL

Production PostgreSQL must have a restorable backup path. For a self-managed baseline use PostgreSQL-supported physical base backup plus WAL archiving/Point-In-Time Recovery; a managed service may provide an equivalent PITR mechanism.

`pg_basebackup` is an accepted physical base-backup mechanism; continuous WAL archiving supplies the recovery timeline.

D7 does not invent an RPO/RTO/cadence number without an operator/business requirement. Before production go-live, deployment configuration must set explicit retention/freshness policy and alert when it is violated.

A backup that has never restored is not proof. Restore drills must prove at least:

- schema + RLS/roles/policies restore coherently;
- canonical owner history/idempotency/effect markers survive;
- ApplicationSession rows may restore but expired/revoked/time-invalid sessions still fail current auth checks;
- River technical state may restore without causing blind external redispatch because D7-C dispatch markers/owner state remain authoritative.

### 10.2 Keycloak continuity

Because MPC human identity binding depends on stable Keycloak `(issuer, sub)`, a self-hosted Keycloak deployment must preserve/backup the persistent IdP state needed to keep those subjects stable. A realm export is configuration evidence, not automatically the sole production backup strategy.

Exact Keycloak database/topology backup mechanism belongs to the chosen deployment environment but must have a restore proof before production.

### 10.3 Binary object custody

Committed MPC binary references must not depend on an unbacked single-host filesystem. The selected private object store must provide durability/backup/versioning or equivalent recovery appropriate to retained content.

A restore/integrity probe checks that committed binary references resolve to bytes matching retained digest/size/type metadata. Missing committed bytes are a material recovery failure, not silently converted to `unknown`.

## 11. Real-dependency proof baseline

D7 architecture claims remain design authority until implementation opens, but the accepted implementation must supply executable falsifiers using real dependencies where the property requires them.

### 11.1 Required real seams

**PostgreSQL:** real server, runtime role + migration role, RLS/FORCE RLS, composite FKs, pooled connection reuse, transaction-local scope, revision/idempotency concurrency.

**River:** real PostgreSQL-backed queue, `InsertTx` commit/rollback, duplicate/rescue/restart, ambiguous dispatch-marker recovery, scheduler recovery.

**Keycloak/OIDC:** real Keycloak realm/client fixtures for code+PKCE callback, wrong state/nonce/issuer/audience, client credentials, JWKS rotation/unknown `kid`, Principal binding.

**HTTP/browser:** real HTTP server plus browser-capable test for Secure/HttpOnly/SameSite cookie behavior, CSRF, same-origin/cross-origin rejection and frontend session bootstrap.

**Product OAD runtime:** generated Chi strict server + real request validator; every canonical `:verb` route dispatches; undeclared properties/pattern/media-type violations fail; 405/Allow and auth carrier alternatives remain exact.

**Object storage:** real S3 API-compatible integration proving private upload/read, committed-reference integrity, unauthorized byte-delivery denial and orphan cleanup behavior.

Mocks/fakes remain useful for owner unit tests but cannot close claims about PostgreSQL RLS, River durability, Keycloak protocol behavior, browser cookie/CSRF behavior, router dispatch or object-store custody.

### 11.2 Local/CI harness direction

Use Docker/OCI containers for real dependency fixtures and Docker Compose or an equivalently simple repo-controlled harness during implementation/testing. This is test/development orchestration, not production Kubernetes authority.

A browser proof harness such as Playwright is justified for D7-D/D8 browser semantics; exact version is test-manifest concern.

Live Mercado Livre/Sankhya write tests remain separately operator-authorized. D7 proof uses sanctioned sandboxes/read-only/live-safe probes where available and never fabricates a provider success through mocks as integration proof.

## 12. Explicit non-selections

D7-E does not admit by preference:

- Kubernetes/operators/service mesh;
- a second API gateway/BFF;
- Redis/cache infrastructure;
- Kafka/NATS/RabbitMQ;
- ELK/OpenSearch/Loki stack as required architecture;
- Prometheus/Grafana as required architecture;
- OpenTelemetry Collector as required architecture;
- OpenTelemetry logs pipeline while Go logs are not a required stable signal;
- Vault/cloud secret-manager as a mandatory dependency;
- hot dynamic configuration;
- ORM/schema auto-generation;
- application-startup auto-migration;
- CDN/public object bucket/persistent presigned URLs;
- multi-region active-active;
- standby/replication HA merely by preference;
- generic backup platform;
- generic infrastructure-as-code framework selection before the actual deployment provider exists.

## 13. Falsifiable proof contract

D7-E / whole-D7 closeout must leave a proof plan that can falsify at least:

1. a canonical `{id}:verb` path cannot be registered/dispatched by the selected Chi/oapi-codegen runtime;
2. a handwritten Product route/DTO/Permission map can drift independently from the OAD;
3. request runtime accepts undeclared fields/pattern/media-type violations rejected by the OAD;
4. validator/generator raw errors leak instead of accepted Product Problem semantics;
5. technical ingress/media delivery leaks into Product OAD/SDK;
6. an unauthenticated or unauthorized caller obtains private authored-media bytes;
7. PostgreSQL commits the first media reference before durable object storage exists;
8. object-store orphan becomes business state or committed missing bytes become silent success;
9. secret/PII appears in logs/traces/metrics/River args;
10. telemetry backend outage blocks owner transactions or makes the app unready by itself;
11. runtime starts against incompatible migration/schema version;
12. runtime process can use migration-owner/DDL privileges;
13. app startup races multiple replicas through auto-migration;
14. untrusted forwarded headers alter OIDC/public-origin security decisions;
15. provider/Sankhya outage incorrectly marks the whole app unready;
16. backup exists but cannot restore the required RLS/history/idempotency/effect-marker state;
17. restored River work blindly re-dispatches a possibly accepted external effect;
18. Keycloak restore loses stable `(issuer, sub)` continuity without detection;
19. committed binary references cannot be integrity-verified after restore;
20. a proof claim about PostgreSQL/River/Keycloak/browser/router/object storage passes using only mocks.

These supplement, not replace, the D7-B/C/D slice-specific falsifiers.

## 14. Current evidence basis

Revalidated 2026-08-22 against current primary/current project documentation:

- D5 canonical OpenAPI tooling authority pins `oapi-codegen v2.8.0` and explicitly leaves router/runtime validation to D7;
- current `oapi-codegen` documents Chi + strict-server generation and separate `nethttp-middleware` request validation;
- current Chi v5 routing source models parameter tail delimiters followed by static suffixes, satisfying MPC's `{id}:verb` routing constraint;
- Go standard-library `log/slog` provides structured logging;
- current OpenTelemetry Go status marks traces and metrics stable; OTLP HTTP exporters are supported, while logs remain less mature;
- current tern v2 remains a PostgreSQL SQL-first standalone/library migrator;
- current Docker guidance supports multi-stage Go builds and one application concern/service per container while allowing one application process to own internal goroutines;
- current PostgreSQL 18 documentation provides `pg_basebackup`, continuous WAL archiving and PITR as supported backup/recovery mechanisms.

Exact runtime/test dependency versions other than the already-D5-pinned OAD tooling are implementation-manifest concerns and must be locked rather than floated at build/deploy time.

## 15. Adjudication

**Candidate:** select Chi v5 + D5-pinned `oapi-codegen v2.8.0` strict server + `nethttp-middleware`; generate runtime operation-policy metadata from the OAD; keep private MPC binaries in private S3 API-compatible object storage with authenticated Go delivery; use typed startup config with deployment-injected secrets; use `tern/v2` through a separate migration role/release step; use JSON `slog` plus OpenTelemetry traces/metrics over OTLP/HTTP; package frontend + Go as one immutable OCI application image behind a trusted TLS edge; use dependency-aware liveness/readiness without provider false-death; require PostgreSQL PITR/restore proof, Keycloak identity-continuity backup and binary-integrity recovery; and require real PostgreSQL/River/Keycloak/browser/router/object-store falsifiers before implementation is accepted.

If ratified, D7-A→D7-E are complete. Next is **whole-D7 coherence + executable proof/adversarial review**, not D8 yet.

Do not begin D8, D9 or Product implementation.