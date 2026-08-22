# D7 — Final Independent Fable Challenge

> Review branch only: `review/d7-fable`  
> Candidate branch: `stage/d7-runtime`  
> Candidate HEAD expected: `c08a4d025cfd89269cc071f4b307695e79f6f8cb`  
> Candidate Draft PR: #58 — `docs(d7): open runtime authority stage`  
> Base `main` expected: `8609d6ebdfe62c7c2e9e5ba249e0761adf8243ef`  
> Product candidate: 99 Product operations / 30 ordinary Permissions / Principal kinds H-A-S only  
> D7-A→D7-E: OPERATOR-RATIFIED  
> D7-R1: internal whole-stage bounded repair candidate  
> D8–D9: BLOCKED  
> Product implementation: BLOCKED UNTIL D9

## Purpose

Run the isolated **final whole-D7 adversarial challenge** against the exact current candidate before any D7 closeout or D8 opening.

This is not a preference review and not a technology-fashion review. Attempt to falsify D7-A→D7-E plus the D7-R1 bounded coherence repairs as **one composed target runtime architecture**. Search for material contradictions, impossible seams, hidden business authority, security gaps, failure/recovery gaps, proof blind spots, YAGNI violations, duplicated mechanisms or an implementation claim that the repository has not actually earned.

Do not optimize for agreement with GPT, the operator, prior review history or the current candidate. Reconstruct the reasoning independently from current repository authority and proportional current primary evidence.

Reviewer output is **evidence, not authority**. Do not edit the candidate branch or PR #58. Write only below `## Fable response` in this file on `review/d7-fable`.

## Mandatory revalidation before analysis

Independently record:

1. remote `main` HEAD;
2. `stage/d7-runtime` HEAD;
3. PR #58 base/head/state/draft/mergeability;
4. candidate changed files and count relative to `main`;
5. GitHub Actions / required CI on the exact candidate HEAD;
6. this review branch ancestry/tree relation to the candidate;
7. that `candidate...review` differs by **only** `docs/work/current/ai-dialog.md`.

If `stage/d7-runtime` differs from:

```text
c08a4d025cfd89269cc071f4b307695e79f6f8cb
```

**STOP and report `STALE_REVIEW_CANDIDATE`.** Do not review a moved target.

Expected exact candidate proof includes:

```text
ci / required                 SUCCESS
pr-title                      SUCCESS
accepted D5 baseline          95/95 operations · 29/29 Permissions · 12/12 controls
current Product               99/99 operations · 30/30 Permissions · 28/28 List/Search
Principal kinds               H/A/S
Performance controls          7/7
Auth controls                 5/5
TS + Go projections           PASS
Performance knowledge         2/2
legacy runtime population     0
bootstrap                     19051 / 20480
repository negative controls  1/1
gate                          PASS
```

Do not interpret current `legacy/runtime population = 0`, `runtime schema enforcement = NOT_CLAIMED_D7`, or historical `router selection = NONE_D7` proof output as a contradiction by itself. The repository has no active Product implementation and D7 selects **target realization authority + falsifiable proof seams**, not an implemented runtime. Challenge whether that distinction is coherent and consistently stated.

## Strict reading discipline

Start exactly:

1. `AGENTS.md`
2. `docs/index.md`
3. `docs/roadmap.md`
4. `docs/engineering/rebaseline/D7-RUNTIME-JOBS-TRANSACTIONS.md`
5. `docs/engineering/rebaseline/D7-R1-WHOLE-STAGE-COHERENCE.md`

Then switch to one bounded D7 owner only when a concrete falsifier requires it:

- `D7-B-POSTGRESQL-ISOLATION-TRANSACTIONS.md`
- `D7-C-DURABLE-WORK-EXTERNAL-EFFECTS.md`
- `D7-D-AUTHENTICATION-SESSION-CSRF.md`
- `D7-E-OPERABILITY-DEPLOYMENT-PROOF.md`

Use exact accepted prior owner authority only when required by the counterexample, most likely:

- D2 for Organization/Principal/data ownership;
- D3 for Q/C/E/P and recoverable propagation;
- D4 for provider/business-system effect ambiguity/reconciliation;
- D5-R1 for H session vs A/S bearer;
- D5 W1/W2/OAD for idempotency, ETag, Product paths and auth projection;
- `ARCHITECTURE.md` for stable constraints.

Do not recursively ingest D0–D6 history. Use current official upstream evidence only where a technology-dependent claim materially depends on it.

## Current D7 architecture to attack

### D7-A — runtime envelope

Accepted baseline:

```text
one Go application process per replica
  same-origin frontend/static delivery
  Product API
  Technical Non-Product Ingress
  H session/CSRF/OIDC mediation
  in-process River workers
  scheduler seam
  PostgreSQL pool
```

Owner consequential transactions are local to one semantic owner. External network writes occur only after commit. No distributed business transaction spans owners. No internal self-HTTP exists merely for symmetry. API/worker split is deferred until a measured need exists.

Attempt to prove that one process creates an unavoidable correctness/security/lifecycle contradiction, not merely that separate services are conventional.

### D7-B — PostgreSQL isolation/transactions

Accepted core:

```text
organization_id on organization-owned state/evidence
composite Organization foreign keys
transaction-local PostgreSQL scope
ENABLE + FORCE RLS
runtime role non-owner / NOSUPERUSER / NOBYPASSRLS
READ COMMITTED + explicit locking
opaque random owner revision tokens
organization/operation-scoped idempotency
pgx/v5 + pgxpool
```

D7-R1 adds a narrow `authentication_bootstrap` persistence scope and distinguishes explicitly platform-global technical/library state such as River engine tables.

Attack whether:

- authentication bootstrap can become a generic cross-tenant bypass;
- Principal-self or technical-routing modes leak into business tables;
- River platform rows can accidentally become business authority;
- composite FK/RLS rules are impossible for any accepted state shape;
- idempotency and owner transaction ownership disagree;
- pool reuse can retain tenant context;
- current role topology permits RLS bypass;
- any needed cross-owner operation secretly requires a shared transaction.

### D7-C — durable work/external effects

Accepted core:

```text
River over PostgreSQL/pgx
InsertTx = owner-state -> durable-work atomic handoff
no second generic MPC outbox
repeat-safe; no exactly-once claim
semantic owner idempotency = correctness
pre-dispatch marker before consequential external write
possible acceptance => never redispatch
source-authoritative reconciliation
scheduler/job state = technical wake-up only
```

Attack crash windows precisely. Include at least:

- crash before marker commit;
- crash after marker commit but before network send;
- crash after network send but before outcome persist;
- stuck-job rescue overlapping live work;
- duplicate consumer reaction;
- out-of-order progression;
- multiple replicas;
- exhausted safe retries;
- restored River rows after database rollback.

A correction must preserve D4's no-blind-retry law and may not turn River state into business truth.

### D7-D — authentication/session/CSRF/machines

Accepted core:

```text
Keycloak first OIDC/OAuth provider
H: confidential Authorization Code + PKCE S256
   go-oidc + x/oauth2
   verified issuer/sub -> fresh opaque MPC session
   discard human OIDC token set after callback
   PostgreSQL stores session-handle digest only
   30m idle / 8h absolute baseline
   X-CSRF-Token synchronizer + net/http CrossOriginProtection
A/S: Client Credentials
     audience includes https://conexus.fun
     jwx/v3 trusted-JWKS verification
     explicit machine-client -> A/S Principal binding
```

Attack:

- state/nonce/PKCE replay/fixation;
- human token exposure/retention;
- restored old sessions after PITR;
- disabled Principal still using a session;
- CSRF and cross-origin bypass;
- session/CSRF storage leakage;
- machine wrong issuer/aud/kid/time/client binding;
- A/S resolving to H;
- IdP role/scope becoming MPC Permission;
- simultaneous human cookie + machine bearer;
- OIDC/provider outage being confused with current MPC authorization.

### D7-E — HTTP, storage, operability, recovery

Accepted core:

```text
Chi v5
D5-pinned oapi-codegen v2.8 strict Chi server
nethttp-middleware OAD request validation
generated operation-policy metadata
private S3 API-compatible byte custody + authenticated Go delivery
typed startup config + injected env/file secrets
tern/v2 for MPC-owned schema
JSON slog
OTel traces + metrics over OTLP/HTTP
one immutable OCI image with Go + compiled frontend
trusted TLS proxy boundary at https://conexus.fun
PostgreSQL base backup + WAL/PITR or managed equivalent
Keycloak subject continuity + binary integrity restore proof
```

D7-R1 adds:

- River's own migration tool owns River schema alongside `tern` for MPC schema;
- post-PITR/database rollback enters a recovery/write fence;
- OAD validator security must be scheme-aware over already-established D7-D carrier context; production cannot use `NoopAuthenticationFunc`.

Attack:

- Chi `{id}:verb` compatibility and 405 semantics;
- generated strict server vs runtime validator authority split;
- hidden handwritten Permission/route map;
- OAD auth alternatives and middleware ordering;
- object-first/DB-reference-second binary consistency, retries and orphan cleanup;
- unauthorized media delivery;
- migration ordering/version skew between MPC and River;
- runtime accidentally having migration-owner power;
- secrets in env/files/logs/traces/metrics/jobs;
- trusted proxy spoofing;
- readiness false positives/false negatives;
- telemetry dependency coupling;
- backup/restore coherence among PostgreSQL, Keycloak, River and object storage.

## D7-R1 bounded corrections to attack specifically

### R1-F1 — authentication bootstrap scope

Try to prove that adding `authentication_bootstrap` is too broad, unnecessary, or insufficient. It must only read authentication mechanism state needed to resolve exactly one Principal and must never read organization-owned business/evidence state.

### R1-F2 — River migration ownership

Try to prove migration ownership is duplicated or incomplete. `tern` may own only MPC schema; the exact-version River migration tool owns River internal schema. Neither runs from ordinary app startup.

### R1-F3 — PITR recovery fence

Try to find a path where database rollback loses a dispatch marker/session/access revocation but the restored system still:

- redispatches a possibly accepted effect;
- resurrects an old browser session;
- resumes old machine automation authority blindly;
- silently accepts missing committed bytes;
- or cannot ever leave recovery mode proportionately.

Challenge whether the proposed scope-by-scope reconciliation/write fence is sufficient and minimal.

### R1-F4 — OAD validator ↔ carrier auth composition

Try to prove the ordering or `AuthenticationFunc` contract permits:

- no-op auth bypass;
- H session accepted as machine scheme;
- bearer accepted as H scheme;
- both carriers with implicit priority;
- duplicate cryptographic auth logic in the validator;
- Permission mapping drift from OAD.

### R1-F5 — proof timing

Challenge whether the candidate now cleanly distinguishes:

```text
D7 architecture/proof-contract closeout
!= D8 golden-flow validation
!= D9 architecture review
!= post-D9 Product implementation
!= implemented real-dependency conformance PASS
```

Do not demand Product implementation during D7 merely to satisfy a proof that repository authority explicitly blocks. Conversely, flag any candidate text that claims runtime behavior as already executed when only target authority/evidence exists.

## Cross-slice scenarios you must attempt

1. Human logs in on replica A, callback lands on replica B, then accesses Organization-scoped Product state.
2. Principal is disabled while an MPC session remains unexpired.
3. Machine token is valid cryptographically but has no current MPC A/S binding or current Permission.
4. Request presents valid H cookie and valid machine bearer simultaneously.
5. Product request with an unknown field passes generated Go typing but should fail OAD runtime validation.
6. Product `/{id}:verb` request has wrong method and must preserve correct method/405 behavior.
7. Owner transaction commits business state plus River job, then process dies before worker start.
8. Worker crashes after dispatch marker but before sending bytes externally.
9. Worker sends external write, external system accepts, process dies before MPC outcome persist.
10. Database PITR restores to before scenario 9's dispatch marker while external acceptance survives.
11. Database restore resurrects a session that was revoked after the restored point.
12. River library is upgraded but its schema migration is omitted or mismatched.
13. Object upload succeeds but owner DB transaction rejects/rolls back.
14. DB commits binary reference, then object is missing/corrupt after recovery.
15. OTLP backend fails under load while owner transactions continue.
16. Sankhya is unavailable while non-Sankhya Product reads should remain available.
17. Untrusted client injects forwarded host/proto headers affecting OIDC redirect/security origin.
18. A restored system has incomplete external reconciliation coverage; test whether it fails safe without requiring a platform rewrite.

## YAGNI / dependency necessity challenge

Require a present property/consumer for each selected mechanism:

- one Go modular-monolith process;
- PostgreSQL + pgx/pgxpool;
- RLS + composite FKs;
- River;
- Keycloak;
- go-oidc + x/oauth2;
- jwx/v3;
- Chi v5;
- oapi-codegen strict server;
- nethttp-middleware;
- private S3-compatible storage;
- tern/v2;
- slog;
- OpenTelemetry traces/metrics;
- one OCI image.

Try to prove any selected dependency is redundant or materially insufficient. Also try to prove a deferred mechanism is **already required**.

Do not recommend Kubernetes, microservices, Redis, Kafka/NATS/RabbitMQ, workflow engines, Vault, service mesh, event sourcing, generic IAM, generic repository/ORM, CDN, multi-region, generic IaC or a mandatory observability backend without a concrete current falsifier.

## Product/non-regression fence

Independently verify D7 has not changed:

```text
99 Product operations
30 ordinary Permissions
Principal kinds H / A / S only
stable origin https://conexus.fun
H session + CSRF
A/S Client Credentials bearer
```

Flag any hidden Product operation, Permission, client class, business owner, generic workflow/status entity, provider-shaped Product vocabulary or Direct Oracle fallback.

## Proof-quality challenge

Separate three categories explicitly:

### A. Current executable repository proof

What does current green CI genuinely prove now?

### B. D7 architecture proof contract

What has D7 made falsifiable/implementable but not yet executed as Product runtime behavior?

### C. Implementation acceptance proof

Which claims require a real implemented runtime and real PostgreSQL/River/Keycloak/browser/router-validator/object-store dependencies after the implementation gate opens?

Flag tautological controls, missing subject population, mocks incorrectly promoted to integration proof, or a proof obligation that cannot realistically be executed later.

## Output contract

Append below `## Fable response` using this exact structure.

### 1. Verdict

Choose exactly one:

- `ACCEPT`
- `ACCEPT WITH BOUNDED FIXES`
- `REOPEN SMALLEST AUTHORITY`
- `REJECT / RECONSTRUCT`

### 2. Revalidation record

Record exact main SHA, candidate SHA, PR state/base/head, changed-file count, exact candidate CI status and review-isolation result.

### 3. Executive coherence assessment

Assess whole-D7 coherence, DevelopmentConexus Method alignment, YAGNI/proportionality and whether the candidate is ready for GPT closeout adjudication.

### 4. Material findings

Number findings highest severity first. For each include:

- **classification:** `D7_FIX`, `D7_R1_FIX`, `D0_D6_REOPEN`, `D8_OBLIGATION`, `IMPLEMENTATION_PROOF`, `REPOSITORY_FIX`, `LATER_NON_BLOCKING`, or `REVIEW_FALSE_POSITIVE`;
- **severity:** Critical / Important / Minor;
- exact candidate location;
- governing repository authority;
- current primary external evidence when technology-dependent;
- concrete counterexample/failure;
- smallest correction;
- why it belongs in that exact owner/stage.

If no material finding exists, say so explicitly. Do not manufacture style findings.

### 5. Persistence / isolation / transaction assessment

Adjudicate RLS, composite FKs, transaction-local scope, authentication bootstrap, River technical state, idempotency/revisions and owner transaction boundaries.

### 6. Durable work / external-effect assessment

Adjudicate River atomicity, duplicates/rescue/order, dispatch markers, retry classification, reconciliation and PITR recovery fence.

### 7. Authentication / request-trust assessment

Adjudicate H code+PKCE/session/CSRF, A/S machine bearer, current access checks, dual-carrier behavior and auth bootstrap.

### 8. HTTP / OAD runtime assessment

Adjudicate Chi, strict server, request validation, `AuthenticationFunc`, generated policy metadata, `:verb` routes and technical-surface exclusion.

### 9. Storage / migrations / operability / recovery assessment

Adjudicate S3 byte custody, object↔DB consistency, tern + River migrations, config/secrets, slog/OTel, OCI/edge/readiness and restore safety.

### 10. Product / authority non-regression assessment

Explicitly state whether 99/30/H-A-S, semantic ownership, provider isolation, Sankhya API Gateway-only and no-implementation-before-D9 remain intact.

### 11. YAGNI / dependency assessment

For selected and rejected/deferred mechanisms, identify any unnecessary or missing dependency based on an actual current property.

### 12. Proof assessment

Separate current repository proof, D7 executable proof contract and future implementation real-dependency proof. State any material blind spots.

### 13. Reconstruction decision

Answer explicitly: is there any material reason to reconstruct/reopen broader accepted D0–D6 architecture or D7-A→E before D7 closeout?

### 14. Continuation recommendation

State the smallest exact action after GPT adjudication. No review output authorizes merge, D8 opening, D9 or Product implementation.

---

## Interaction rule

Fable writes **only** to this file on `review/d7-fable`. Do not edit PR #58, `stage/d7-runtime`, `main` or any other review branch.

GPT will independently adjudicate every material finding against current repository authority and executable evidence. Round 2 is justified only for a surviving material contradiction after adjudication/bounded fixes.

---

## Fable response

<!-- Fable: append independent whole-D7 review here. -->