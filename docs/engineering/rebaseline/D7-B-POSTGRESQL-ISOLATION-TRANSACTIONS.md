# D7-B — PostgreSQL Isolation & Transaction Realization

> **Status:** CANDIDATE / OPERATOR RATIFICATION PENDING  
> **Parent:** `D7-RUNTIME-JOBS-TRANSACTIONS.md`  
> **Accepted prerequisite:** D7-A Runtime Envelope & Transaction Ownership Boundary — OPERATOR-RATIFIED  
> **Parent authorities:** accepted D2 Organization/data ownership, accepted D5 wire/idempotency/revision grammar, accepted D7-A transaction ownership  
> **Method:** DevelopmentConexus Engineering Method v1.0.0  
> **Derived:** 2026-08-21

## 1. Purpose

D7-B defines the smallest PostgreSQL realization that makes accepted Organization isolation, owner-local transaction boundaries, Product idempotency and opaque revision/precondition semantics structurally falsifiable.

D7-B does **not** define the full target table census, domain schema layout, D7-C durable-worker mechanism, D7-D session persistence, D7-E migrations/deployment mechanics or Product implementation.

## 2. Target invariant

> **An ordinary query/command bug cannot cross Organization merely because a developer omitted a predicate; an organization-owned row cannot reference a foreign-Organization row; pooled connections cannot retain prior Organization scope after a unit of work; and duplicate/stale consequential intake is resolved inside the owner transaction without exposing database implementation identity as Product revision authority.**

## 3. Accepted inputs

From D2:

- Organization is the canonical isolation root;
- every organization-owned persistent business state or persisted external observation/evidence belongs to exactly one Organization;
- Organization scope is explicit and is not inferred from Installation, Selling Entity, external account, IdP Organization or process-global default;
- a realization based only on remembered `WHERE organization_id = ?` is invalid;
- platform-owned definitions with no tenant-specific business meaning may remain platform-scoped;
- one meaning has one write authority.

From D5/W1/W2:

- Organization path scope is a claim, never self-authorization;
- cross-Organization secondary references fail closed;
- `Idempotency-Key` and revision proof are independent;
- exact prior idempotent intake is resolved before stale-revision re-evaluation;
- strong ETags are opaque server-issued owner validators, never Postgres `xmin`, database sequence, timestamp, provider version or client-authored version;
- standard same-resource conditionals use `If-Match`; custom/reference revision proofs use the same opaque validator as typed request data.

From D7-A:

- one owner-local transaction may atomically persist owner state plus required technical idempotency/audit/durable-handoff records;
- another semantic owner's private business state is never mutated in that transaction;
- external network writes occur only after commit.

## 4. Scope classes

D7-B admits exactly three persistence access modes; none is ambient/default.

### 4.1 Organization scope

Normal Product business state/evidence access runs under one exact Organization.

Examples include:

- organization-owned domain state;
- external observations/evidence persisted for that Organization;
- Membership/RoleAssignment when accessed in an Organization administration flow;
- organization-scoped idempotency/audit/intake records.

### 4.2 Principal-self platform scope

The platform-scoped `GET /access-context` needs bounded current identity/access self-discovery before one exact Organization is selected.

This mode may read only the identity/access state required to establish:

- current Principal access eligibility;
- that Principal's own visible Organization Memberships.

Once an Organization is selected, effective Organization-specific role/Permission resolution returns to Organization scope. Principal-self scope does not become a generic cross-tenant read mode.

### 4.3 Technical-routing platform scope

Technical ingress and the durable worker runner may need to resolve a trusted technical correlation into an exact Organization before touching business state.

Only bounded technical routing/queue records may be accessed in this mode. After exact Organization resolution, business/evidence access starts a new Organization-scoped unit of work.

A technical routing record may carry `organization_id` for routing but gains no business authority by persistence.

## 5. Transaction-local scope contract

Every scoped database unit of work uses an explicit PostgreSQL transaction and establishes scope with transaction-local settings before scoped SQL executes.

Conceptual envelope:

```text
BEGIN
  set_config('mpc.scope_mode',       <organization|principal_self|technical_routing>, true)
  set_config('mpc.organization_id',  <exact organization when required>, true)
  set_config('mpc.principal_id',     <exact principal when required>, true)
  ... scoped SQL ...
COMMIT / ROLLBACK
```

Binding laws:

- `is_local=true` is required so scope ends with the transaction and cannot leak through pooled connection reuse;
- missing, malformed or incompatible scope resolves to no permitted business rows / explicit failure, never a default Organization;
- application code does not set persistent/session-wide tenant state;
- Organization-scoped Qs use a read-only transaction when no write is needed; the transaction exists to bind structural scope, not to manufacture write authority;
- Product commands use a read-write owner transaction;
- platform-scoped configuration/definition reads that genuinely carry no tenant meaning need no fake Organization context;
- no external provider/business-system call occurs while this transaction is open.

## 6. PostgreSQL role boundary

Use separate database identities for schema ownership/migration and application runtime.

### Runtime role

The application runtime role must be:

```text
LOGIN
NOSUPERUSER
NOCREATEROLE
NOBYPASSRLS
not owner of organization-scoped tables
not a member of a role that can assume table-owner/BYPASSRLS power
```

It receives only required DML/sequence/schema privileges.

### Schema/migration owner

A separate non-application role owns schemas/tables and applies migrations/policies. Its credential is never loaded by the running application.

Superuser/`BYPASSRLS` credentials are administrative break-glass infrastructure only and are not application/runtime paths.

## 7. Structural Organization isolation

### 7.1 Organization key on scoped rows

Every organization-owned persistent business row and persisted external observation/evidence row carries:

```text
organization_id NOT NULL
```

The field is data-isolation scope, not business ownership inference.

### 7.2 Composite referential integrity

Organization-owned references must structurally preserve the Organization dimension.

For a scoped parent resource with opaque ID:

```text
UNIQUE (organization_id, parent_id)
```

A scoped child references it through:

```text
FOREIGN KEY (organization_id, parent_id)
  REFERENCES parent (organization_id, parent_id)
```

A globally unique opaque ID does not justify dropping `organization_id` from the FK: the composite relationship is what makes a cross-Organization reference physically invalid.

Do not introduce composite Product identity. Product IDs remain their accepted opaque identities; the composite key is persistence isolation machinery.

### 7.3 RLS baseline

Every organization-scoped table must use both:

```sql
ALTER TABLE ... ENABLE ROW LEVEL SECURITY;
ALTER TABLE ... FORCE ROW LEVEL SECURITY;
```

Policies are explicit for the runtime role and enforce both visibility (`USING`) and new-row/update scope (`WITH CHECK`).

Organization-mode policy meaning is equivalent to:

```text
row.organization_id == transaction-local exact organization_id
```

No active scope/policy means default deny.

`FORCE ROW LEVEL SECURITY` is defense in depth for table-owner access; it does not make superuser or `BYPASSRLS` safe runtime roles, so those remain forbidden to the application.

### 7.4 Bounded non-Organization policies

Principal-self policy exists only on the minimum identity/access relations required by `/access-context`; it admits rows belonging to the exact transaction-local Principal and does not grant general cross-Organization business reads.

Technical-routing mode policies exist only on bounded technical routing/queue relations. Organization-owned business/evidence tables do not admit technical-routing mode.

No universal policy DSL or dynamic ACL engine is introduced.

## 8. Database primitive candidate

### 8.1 `pgx/v5` + `pgxpool` — SELECT CANDIDATE

D7-B selects `pgx/v5` + `pgxpool` as the candidate PostgreSQL connection/transaction primitive because:

- PostgreSQL is already canonical rather than an interchangeable database target;
- `pgxpool` is the concurrency-safe pool for a Go server;
- pgx exposes explicit `BeginTx` / transaction modes and transaction helpers without requiring an ORM/repository abstraction;
- current pgx remains the stable v5 line and tracks supported Go/PostgreSQL releases.

Important implementation obligation from current pgx behavior: context cancellation does not auto-rollback an already-started pgx transaction. The transaction wrapper must deterministically `Commit` or `Rollback`; helper functions such as `BeginTxFunc` may realize that invariant.

### 8.2 Rejected/deferred

- generic ORM / generic repository abstraction — **REJECT**; hides SQL/RLS/constraint semantics that are part of correctness;
- `database/sql` portability as architectural goal — **REJECT**; no multi-database consumer exists;
- `sqlc` — **DEFER** until concrete schema/query surfaces exist; it may later generate typed query glue without becoming persistence authority;
- `tern` — **DEFER to D7-E migration mechanism**;
- schema-per-Organization / database-per-Organization — **REJECT baseline**; no independent isolation/deployment consumer justifies migration/connection/operability multiplication;
- database role per Organization — **REJECT baseline**; unnecessary credential/role explosion when transaction-local RLS can satisfy current isolation properties.

## 9. Transaction isolation and locking

### 9.1 Baseline isolation

Use PostgreSQL `READ COMMITTED` as the default transaction isolation level.

It is sufficient only together with explicit current-row locking/atomic constraints where owner semantics require them. D7-B does not rely on a multi-statement read snapshot accidentally remaining stable.

Do **not** make `SERIALIZABLE` the global default merely for theoretical strength; it introduces serialization retries and still does not replace explicit Product revision/idempotency semantics.

A bounded operation/read may select `REPEATABLE READ`, `SERIALIZABLE` or a stronger lock only when a concrete accepted invariant proves the default insufficient.

### 9.2 Mutable protected owner meanings

For a same-owner mutation/capability protected by the current opaque ETag:

```text
SELECT protected row FOR UPDATE under exact Organization scope
  -> not visible/resolvable => ordinary scoped 404
  -> compare server revision token with supplied validator
       standard If-Match stale => 412
       typed custom/reference stale => 409 resource-revision-conflict
  -> evaluate current owner state/disposition
  -> mutate
  -> rotate revision token in same transaction
```

The row lock prevents two accepted writers from both consuming the same current revision.

Cross-owner referenced revision proof is validated through the referenced owner's admitted read/validation boundary; D7-B does not acquire private cross-owner locks or create a distributed transaction merely to freeze another owner's state.

## 10. Opaque revision realization

Every protected current owner meaning persists one server-owned opaque revision token.

Binding laws:

- token is generated from cryptographically secure randomness;
- token has no business meaning and is not the resource ID;
- the public strong ETag is an opaque encoding of that token;
- never expose PostgreSQL `xmin`, sequence/counter, row timestamp, provider version or client-authored number as the Product validator;
- every material accepted mutation rotates the token in the same owner transaction;
- replay of an already accepted exact idempotent intake returns the previously established result/current validator according to accepted W2 semantics rather than re-running stale revision evaluation.

Exact binary length/encoding helper is implementation detail as long as collision resistance/non-predictability is materially sufficient and the Product representation stays a strong opaque quoted ETag.

## 11. Idempotency persistence boundary

For every operation whose OAD requires `Idempotency-Key`, the owner transaction uses one organization-scoped technical idempotency record.

Candidate uniqueness scope:

```text
organization_id
+ stable Product operation identity
+ digest(Idempotency-Key)
```

The record preserves:

- semantic request fingerprint digest;
- intake state (`processing` only when durably meaningful, accepted/pending/terminal outcome as applicable);
- stable owner intake/result reference or bounded response snapshot required for exact replay;
- created/completed timestamps needed for operation/reconciliation, without making wall-clock time part of request equivalence.

Semantic fingerprint is generated from the accepted W2 material inputs: Organization, operation/path target, material query/body data, revision proofs, and binary content/metadata when applicable. JSON formatting/order, credentials, request time and the Idempotency-Key itself are excluded.

A cryptographic digest (baseline SHA-256 class) may persist canonical fingerprint/key identity; raw bodies, binary bytes, bearer/session credentials and arbitrary PII are not retained merely for dedupe.

### 11.1 New intake order

After AuthN/contract/access gates:

```text
BEGIN owner transaction + exact Organization scope
  -> claim/read idempotency key record
     existing + different fingerprint => reused conflict
     existing + durable same result    => replay result
     existing + genuinely processing   => in-progress result
  -> new intake: evaluate revision + owner/Governance prerequisites
  -> persist owner intake/state + idempotency result + required audit/handoff
COMMIT
```

A validation/business failure that occurs before durable intake rolls back a newly-created idempotency claim. Once durable intake exists, its pending/ambiguous/rejected/accepted meaning remains correlated to that key.

No idempotency record authorizes blind external-effect redispatch.

Exact retention duration is not frozen without a proved retry/history horizon; D7-E/D8 may bound compact record retention without breaking accepted replay/reconciliation semantics.

## 12. Pool and worker safety

- every Organization-scoped pool checkout is used through the transaction-scope wrapper; no session-wide tenant variable survives release;
- transaction-local settings are established before scoped SQL and disappear on commit/rollback;
- workers claim only platform technical durable-work records without business reads, obtain the exact persisted Organization scope from the accepted handoff, then enter a new Organization-scoped owner transaction;
- a worker never uses `BYPASSRLS` to process multiple Organizations;
- provider callbacks/technical ingress may resolve an Organization through the bounded technical-routing surface, then must leave technical-routing mode before business/evidence access.

## 13. Falsifiable proof contract

D7-B cannot be ratified until a real PostgreSQL proof design demonstrates at least:

1. runtime role is not table owner, superuser or `BYPASSRLS`;
2. every organization-scoped table has `organization_id NOT NULL`, RLS enabled, RLS forced and an admitted policy;
3. missing scope cannot read/write organization-owned rows;
4. Organization A cannot read/update/delete Organization B rows even when SQL omits an explicit `organization_id` predicate;
5. Organization A cannot insert/update a row whose `organization_id` is B (`WITH CHECK` failure);
6. a composite FK rejects an Organization A row referencing an Organization B resource even when the target opaque ID is known;
7. a pooled connection reused after an Organization A transaction cannot observe A scope in the next transaction;
8. principal-self scope can see only the exact Principal's identity/access self-discovery population and cannot read business tables;
9. technical-routing scope cannot read organization-owned business/evidence tables;
10. worker processing uses exact job Organization with the ordinary runtime role and RLS still active;
11. two writers using the same current owner revision cannot both commit a material mutation;
12. stale standard `If-Match` and stale typed ETag remain distinguishable as accepted 412 vs 409 semantics;
13. same idempotency key + different semantic fingerprint fails closed;
14. rollback before durable intake leaves no false successful idempotency result;
15. exact replay of a durably accepted intake returns the established result without rerunning the external effect;
16. migration/schema-owner credentials are absent from runtime configuration.

Mocks are insufficient for RLS, transaction-context, locking and FK claims; these proofs require a real PostgreSQL instance.

## 14. Current primary evidence

Revalidated on 2026-08-21 from current primary documentation:

- PostgreSQL Row Security Policies: <https://www.postgresql.org/docs/current/ddl-rowsecurity.html>
- PostgreSQL `CREATE POLICY`: <https://www.postgresql.org/docs/current/sql-createpolicy.html>
- PostgreSQL role attributes / `BYPASSRLS`: <https://www.postgresql.org/docs/current/role-attributes.html>
- PostgreSQL configuration `set_config(..., is_local=true)`: <https://www.postgresql.org/docs/current/functions-admin.html>
- PostgreSQL constraints / composite foreign keys: <https://www.postgresql.org/docs/current/ddl-constraints.html>
- PostgreSQL transaction isolation / row-update behavior: <https://www.postgresql.org/docs/current/transaction-iso.html>
- pgx v5 package/transaction documentation: <https://pkg.go.dev/github.com/jackc/pgx/v5>
- pgx repository/current support policy: <https://github.com/jackc/pgx>

Relevant current PostgreSQL facts:

- enabling RLS without an applicable policy is default deny;
- table owners normally bypass RLS unless `FORCE ROW LEVEL SECURITY` is used;
- superusers and `BYPASSRLS` roles always bypass RLS;
- policy `USING` governs existing-row visibility/modification and `WITH CHECK` governs proposed inserted/updated rows;
- `set_config(..., true)` limits the setting to the current transaction;
- composite foreign keys are natively enforced constraints.

## 15. Adjudication

**Candidate:** accept transaction-local scope + RLS/`FORCE RLS` + non-owner/no-`BYPASSRLS` runtime role + composite Organization FKs + `READ COMMITTED`/explicit row locking + opaque random owner revision token + organization/operation-scoped idempotency records + `pgx/v5`/`pgxpool` as the direct PostgreSQL primitive.

No Product operation, Permission, Principal kind, semantic owner or frontend decision changes.

If ratified, next is **D7-C — Durable Work & External Effects**, where the accepted D7-A handoff property is used to decide River/outbox/scheduler/retry/reconciliation mechanics.

Do not begin D8, D9 or Product implementation.
