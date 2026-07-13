# F-08 Assisted Sankhya Linkage Implementation Contract

## Canonical Ownership

`orders` owns one `AssistedSankhyaLinkageService`. Domain/application code owns
identity, state transitions, candidate/proof separation, idempotency, and audit.
Postgres is the only canonical MPC store. `internal_read` owns typed Oracle
reads and all Oracle SQL. Transport only validates requests/auth and maps
responses. Profitability consumes only verified exact Oracle source identity.

No provider payload, Oracle driver type, dynamic customer column, or raw SQL
crosses an adapter boundary.

## Domain States And Proof

- Link state: `unlinked | confirmed | conflict | disabled`.
- Lineage state: `unknown | none | partial | complete | conflict`.
- Candidate: bounded header/line projection with match reasons; never persisted
  as proof and never supplies tax identity.
- Proof: exact header key plus explicit operator-confirmed 313
  (`NUNOTA`,`SEQUENCIA`) per immutable `mpc_line_id`, stored append-only.
- Invoice descendant: exact `TGFVAR` origin match and TOP 306 destination;
  append each (`NUNOTA`,`SEQUENCIA`,`QTDATENDIDA`) observation. Missing or
  conflicting facts remain unknown and profitability tax remains missing.

## Ports

- Orders-owned `SankhyaLinkageRepository`: load current mapping, append
  confirmation/conflict/disable/lineage events, and enforce idempotency and
  uniqueness under one transaction scoped by tenant and installation.
- Orders-owned `SankhyaLinkageReader`: bounded candidate headers/lines, validate
  configured header key, and list exact invoice descendants as typed facts.
- `internal_read` adapter implements the reader contract; customer field
  identifiers come only from validated startup configuration/metadata and are
  safely quoted/allowlisted. Values always use binds.
- Profitability receives `TaxSourceIdentity{DocumentID, LineNumber}` only from
  a verified realized 306 descendant. A 313 origin alone is not sale-tax
  identity.

## Ledger And Invariants

Next migration: `apps/server_core/migrations/0033_orders_sankhya_linkage.sql`.
Use append-only event rows plus a transactionally maintained current projection
if query performance requires it. Every business table carries `tenant_id`.

Required identity columns: `tenant_id`, `installation_id`,
`provider_order_id`, immutable `mpc_line_id`, origin `nunota313`, origin
`sequencia313`; descendant rows additionally carry `nunota306`,
`sequencia306`, nullable attended quantity and observation/source time.
Required audit columns: event ID/type, actor ID/type, reason, idempotency key,
source timestamp, recorded timestamp, previous-event reference, configuration
revision, and safe evidence state.

Uniqueness:

- one current order mapping per (`tenant_id`,`installation_id`,`provider_order_id`);
- one current origin per (`tenant_id`,`installation_id`,`provider_order_id`,`mpc_line_id`);
- an origin (`nunota313`,`sequencia313`) cannot be actively assigned to two
  provider lines in the same tenant/installation;
- descendant identity is unique per origin and (`nunota306`,`sequencia306`);
- idempotency key unique per tenant/installation/operation.

Retries return the existing semantically identical result. Any same-key or
same-origin mismatch appends/returns a conflict and never overwrites history.

## Application Flow

1. Authorize tenant/installation and require validated enabled admin config.
2. Resolve the persisted order and immutable MPC line identities. Do not use
   mutable `line_no` as proof.
3. Read bounded candidates using operator-supplied exact 313 identity plus
   optional non-authoritative filters.
4. On confirmation, validate the configured header key, TOP 313, every exact
   line, actor, reason, source time, and current ledger conflicts.
5. Append the mapping idempotently in Postgres. No Oracle write occurs.
6. Read `TGFVAR` descendants by exact origin, validate TOP 306, append lineage
   observations, and expose verified destination identity to profitability.

## Exact Next-Worker Seams

- Orders domain/application/ports:
  `internal/modules/orders/domain/sankhya_linkage.go`,
  `application/assisted_sankhya_linkage_service.go`,
  `ports/sankhya_linkage_repository.go`, and
  `ports/sankhya_linkage_reader.go`.
- Postgres adapter and migration:
  `orders/adapters/postgres/sankhya_linkage_repo.go` and migration `0033`.
- Oracle read boundary: typed models/port in `internal_read/domain` and
  `internal_read/ports/reader.go`; SQL only in
  `internal_read/adapters/oracle/reader.go`; an orders adapter bridges the
  internal-read port without importing Oracle types.
- Transport/composition: orders HTTP handler/routes and composition wiring.
- Contract pair: update `contracts/api/marketplace-central.openapi.yaml` and
  `packages/sdk-runtime/src/index.ts` together for:
  `GET /orders/{provider_order_id}/sankhya-linkage`,
  `GET /orders/{provider_order_id}/sankhya-linkage/candidates`, and
  `POST /orders/{provider_order_id}/sankhya-linkage/confirm`.
  Tenant comes from authenticated context; `installation_id` remains explicit.
  Confirm requires exact header/lines, actor, reason, source time, and
  idempotency key. Responses expose quality/state, never buyer PII.

The next worker must first introduce immutable `mpc_line_id` without deleting
or regenerating identity during order refresh. Existing positional rows require
an explicit migration/reconciliation state; ambiguous duplicates remain
unlinked.

## Fail-Closed Errors

Stable error classes: configuration missing/invalid, field metadata mismatch,
uniqueness unproved, order/line missing, header-key mismatch, TOP mismatch,
candidate ambiguous, mapping conflict, lineage unknown/conflict, and Oracle
unavailable. None maps to zero/default tax or a resolved link. No automatic
backfill uses partner, buyer, date, product, quantity, price, or value.

## Required Proof Before Enablement

- Unit tests for state transitions, tenant/installation scope, duplicate and
  concurrent retries, conflict append, candidate-not-proof, and unknown tax.
- Postgres integration tests for all unique constraints and append-only audit.
- Oracle adapter tests for exact binds/predicates, bounded result limits,
  configured identifier allowlisting, partial/one-to-many descendants, and
  explicit missing/conflict states.
- OpenAPI/SDK parity tests and HTTP authorization/validation tests.
- Separately authorized bounded live-oracle discovery using a registered
  command; generic smoke evidence cannot activate the workflow.
