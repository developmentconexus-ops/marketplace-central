# Marketplace Central — Roadmap

<!-- program-status-authority -->

> **Role:** sole mutable current-stage/status/allowed-work/next-action authority. Detailed semantics live in routed owning documents.

## Current checkpoint

| Field | Current value |
| --- | --- |
| Product | **Marketplace Operations Control Plane + Commercial Intelligence** |
| Current stage | **D7 — Runtime / Jobs / Transactions — OPEN / ACTIVE** |
| Accepted baseline | **D0–D6 ACCEPTED / CLOSED** |
| D7 authority | [D7 Runtime / Jobs / Transactions](engineering/rebaseline/D7-RUNTIME-JOBS-TRANSACTIONS.md) |
| D7-A | **Runtime Envelope & Transaction Ownership — OPERATOR-RATIFIED** |
| D7-B | **PostgreSQL Isolation & Transactions — OPERATOR-RATIFIED** |
| D7-C | **Durable Work & External Effects — OPERATOR-RATIFIED** |
| D7-D | **Authentication / Session / CSRF / Machine Tokens — OPERATOR-RATIFIED** |
| D7-E | **Operability / Secrets / Migrations / Deployment & Proof — CANDIDATE / OPERATOR RATIFICATION PENDING** — [candidate](engineering/rebaseline/D7-E-OPERABILITY-DEPLOYMENT-PROOF.md) |
| Canonical Product OAD | `contracts/api/product/openapi.yaml` |
| Product surface | **99 Product operations · 30 ordinary Permissions · Principal kinds H / A / S only** |
| Stable origin | `https://conexus.fun` |
| Active runtime baseline | **NONE** |
| Exact next action | **Operator adjudicate/ratify D7-E. Do not start whole-D7 closeout review before that decision.** |
| Implementation | **BLOCKED UNTIL D9** |

## Stage progression

| Stage | Status |
| --- | --- |
| D0 — Product / System Definition | ACCEPTED / CLOSED |
| D1 — Domains / Boundaries | ACCEPTED / CLOSED |
| D2 — Identity / Tenant / Data Ownership | ACCEPTED / CLOSED |
| D3 — Communication / Events | ACCEPTED / CLOSED |
| D4 — External Integrations | ACCEPTED / CLOSED |
| D4-R1 — Publication Input / Listing Authoring | ACCEPTED / CANONICAL |
| D5 — API | ACCEPTED / CLOSED |
| D6 — Frontend | **ACCEPTED / CLOSED** |
| D7 — Runtime / Jobs / Transactions | **OPEN / ACTIVE — D7-A→D7-D ratified; D7-E candidate** |
| D8 — Golden Flows | BLOCKED |
| D9 — Adversarial Architecture Review | BLOCKED |
| Implementation | BLOCKED UNTIL D9 |

## D7 accepted inputs

```text
Product operations     99
ordinary Permissions   30
Principal kinds        H / A / S only
stable origin          https://conexus.fun
Go backend             canonical business execution
PostgreSQL             canonical MPC-owned state store
active runtime         NONE
```

Human browser authentication remains server-side OIDC → Secure HttpOnly MPC session + CSRF; A/S remain Client Credentials bearer clients. D7 realizes mechanics only and cannot move Product/frontend/business authority.

## D7 accepted decisions

**D7-A:** one Go process per replica hosts same-origin frontend/Product API/Technical Ingress/auth mediation plus in-process durable workers over PostgreSQL. Owner transactions remain local; external writes occur after commit; no distributed business transaction.

**D7-B:** explicit Organization scope + composite Organization FKs + transaction-local PostgreSQL `ENABLE/FORCE RLS` under a non-owner/non-`BYPASSRLS` runtime role; `READ COMMITTED` + explicit locking; opaque revision tokens; organization/operation-scoped idempotency; `pgx/v5 + pgxpool`.

**D7-C:** River is the bounded PostgreSQL durable-work engine. `InsertTx` is the atomic owner-state→work handoff; no second generic outbox. Delivery is repeatable; semantic idempotency is correctness. Possible external acceptance never redispatches and instead reconciles against authoritative source truth.

**D7-D:** Keycloak is the first OIDC/OAuth provider. H uses confidential Authorization Code + PKCE S256, `go-oidc + x/oauth2`, verified `(issuer, sub)`, discard of the human OIDC token set, opaque PostgreSQL-backed MPC session, synchronizer `X-CSRF-Token` + Go cross-origin protection. A/S use Keycloak Client Credentials, stable Product audience, trusted-JWKS `jwx/v3` verification and explicit client→A/S Principal binding. IdP roles never become MPC Permissions.

## D7-E candidate

[D7-E](engineering/rebaseline/D7-E-OPERABILITY-DEPLOYMENT-PROOF.md) proposes:

```text
HTTP       Chi v5 + D5-pinned oapi-codegen v2.8 strict server
validation oapi-codegen/nethttp-middleware + generated OAD policy metadata
bytes      private S3 API-compatible storage + authenticated Go delivery
config     typed startup config + deployment-injected env/file secrets
migrate    tern/v2 via separate migration owner / release step
logs       JSON log/slog
telemetry  OpenTelemetry traces + metrics over OTLP/HTTP
artifact   one immutable OCI image containing Go + compiled frontend
edge       https://conexus.fun behind an explicit trusted TLS proxy boundary
backup     PostgreSQL base backup + WAL/PITR or managed equivalent + restore proof
proof      real PostgreSQL/River/Keycloak/browser/router/object-store seams
```

The candidate closes D5's `{id}:verb` router constraint and runtime OAD validation without another handwritten Product route/DTO/Permission authority. It keeps private binary bytes out of public buckets and uses no CDN baseline. Provider/Sankhya/OTLP degradation does not automatically make the whole app unready.

No Kubernetes/service mesh, Redis, broker, mandatory Collector/Prometheus/ELK stack, Vault dependency, hot config, ORM auto-schema, startup auto-migration, public/persistent presigned object URLs, multi-region or generic IaC platform is admitted.

After D7-E ratification, the next step is **whole-D7 coherence + executable-proof/adversarial review**. D8 does not open until D7 itself closes.

D8–D9 remain blocked. Product implementation remains blocked until D9. Reopen D0–D6 only for a material falsifier at the smallest owner.

One coherent gate lands before the next. For task-specific reading, return to [`index.md`](index.md).