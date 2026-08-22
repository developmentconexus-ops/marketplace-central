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

### 1. Verdict

**ACCEPT WITH BOUNDED FIXES**

One Important bounded D7-R1 fix (recovery-fence detection/arming law, F-1). No D0–D6 reopen. No reconstruction. Everything else adjudicated coherent or explicitly deferred to its correct later proof stage.

### 2. Revalidation record

| Item | Independently measured value |
| --- | --- |
| remote `main` HEAD | `8609d6ebdfe62c7c2e9e5ba249e0761adf8243ef` — matches expected |
| `stage/d7-runtime` HEAD | `c08a4d025cfd89269cc071f4b307695e79f6f8cb` — matches expected; **NOT STALE** |
| PR #58 | OPEN · Draft · base `main` · head `stage/d7-runtime` @ `c08a4d02` · MERGEABLE · title `docs(d7): open runtime authority stage` |
| candidate changed files vs `main` | 9 files · +1964/−59 · D7 main + D7-B/C/D/E + D7-R1 + `docs/index.md` + `docs/roadmap.md` + `scripts/gate.ps1` — docs/authority + gate only; no implementation surface |
| CI on exact candidate HEAD | `required` SUCCESS · `conventional-title` SUCCESS (check-runs API on `c08a4d02`) |
| review isolation | `merge-base(stage/d7-runtime, review/d7-fable) = c08a4d02`; `candidate..review` = exactly `docs/work/current/ai-dialog.md` |
| local review-mode gate | `GATE_BASE_REF=stage/d7-runtime npm run gate` → **PASS**; reproduced full expected proof block: baseline 95/95 · 29/29 · 12/12; current 99/99 · 30/30 · 28/28; Performance 7/7; Auth 5/5; H/A/S; TS+Go projections PASS; Performance knowledge 2/2; legacy runtime population 0; bootstrap 19051/20480; repository negative controls 1/1 |
| gate wiring on candidate | review-mode lane verifies exact-candidate tree is `docs/work/**`-free and review diff is single-file; roadmap markers flipped to D7 OPEN/ACTIVE with matching forbidden-duplication guards |

`main` already carried D6 ACCEPTED/CLOSED and D7 NEXT/NOT STARTED before this candidate; PR #58 flips only D7 to OPEN/ACTIVE. Stage sequencing is clean — no stage stacking.

### 3. Executive coherence assessment

D7-A→D7-E compose into one coherent, closed target runtime: one Go process per replica; PostgreSQL as the single durable substrate for business state, sessions, idempotency, dispatch markers and River work — which is what makes the atomicity story (InsertTx) and the crash-window story honest without a second infrastructure service; carrier-split authentication that never hands the browser a token; OAD-generated HTTP surface with no handwritten duplicate; and an operability profile that consistently refuses platform sprawl. D7-R1's five repairs are genuine composition seams, not scope creep, and each names the smallest owner. Method alignment is good: decisions carry present consumers, rejections carry reasons, and proof obligations are named per slice. Proportionality is defensible throughout — nothing selected lacks a current consumer, and nothing deferred is provably already required.

The candidate is ready for GPT closeout adjudication after one bounded R1-F3 amendment (F-1 below).

Cross-slice scenario disposition (all 18 attempted): 1 login-across-replicas — covered (login transaction persisted in PostgreSQL, D7-D §4.1). 2 disabled-Principal live session — covered (eligibility recheck per request, D7-D §5.3). 3 valid crypto/no binding — fails closed (D7-D §8.3). 4 dual carrier — fails closed as ambiguous (R1-F4). 5 unknown field past Go typing — exactly why nethttp-middleware is selected; validator keyword support verified (see §8). 6 `:verb` wrong method / 405+Allow — named proof obligation, Chi capability current. 7 commit+job then death — job durable, runs on restart. 8 crash after marker before send — marker + no outcome ⇒ reconcile (D7-C §8). 9 send accepted, death before persist — same rule, no redispatch. 10 PITR before scenario-9 marker — fence covers **when armed**; unarmed path is F-1. 11 restored revoked session — fence step 2; same F-1 arming caveat. 12 River upgrade without migration — R1-F2 release sequence + readiness on incompatible schema. 13 object success/DB rollback — orphan + bounded sweep (D7-E §4.2). 14 committed ref, missing bytes — integrity probe, surfaced as material recovery failure (D7-E §10.3). 15 OTLP failure under load — export never blocks owner transactions or readiness (D7-E §7.2). 16 Sankhya down, other reads live — readiness law forbids provider false-death (D7-E §9). 17 forwarded-header injection — trusted-proxy-boundary-only law (D7-E §8.2). 18 incomplete reconciliation coverage — fence stays closed per scope, surfaced to operator, no rewrite (R1-F3 step 8).

### 4. Material findings

#### F-1 — recovery fence has no detection/arming law; the unarmed-restore path defeats every fence step

- **classification:** `D7_R1_FIX`
- **severity:** Important
- **location:** D7-R1 §4 (“Any disaster/PITR restore … enters a fail-closed recovery fence”) and its Required recovery falsifier; interacts with D7-E §13 falsifier 17.
- **governing authority:** D7-C §8 crash/ambiguity rule; D4 no-blind-retry law; repository verification doctrine — presence is not execution; a material guard requires a deterministic falsifier.
- **counterexample:** an operator (or managed-DB console automation) performs a PITR restore without executing the MPC restore runbook. The restored database contains no fence marker — fence state either never existed or was itself rolled back, so the restored timeline is self-consistently “normal.” Application replicas reconnect, readiness passes (all D7-E §9 conditions are satisfiable against a rolled-back-but-valid schema), River resumes, and a job whose dispatch marker was lost in the rollback sees `no marker -> may dispatch` and re-sends a possibly accepted external write. Every one of R1-F3’s eight fence steps is written in the passive voice (“enters,” “is recorded,” “are invalidated”) with no actor and no trigger: nothing in the composed architecture *detects* that a rollback occurred. The stated falsifier (“restore a snapshot … prove no redispatch before fence release”) presumes the fence is engaged, so an implementation can pass it while remaining wide open on the most probable real-world failure path — human/procedural error during a disaster.
- **smallest correction:** add one binding law to R1-F3: *timeline continuity is established affirmatively, never assumed.* The implementation must include a durable continuity check whose failure arms the fence by default — e.g., comparing the PostgreSQL system identifier / timeline / last-known-position anchor against a copy held outside the database (the already-selected object store or deployment state qualify), or an equivalent deployment-verified control — and the required recovery falsifier must include the unarmed path: a restore followed by an ordinary application boot with **no** manual arming step must still fail closed before external dispatch. Exact anchor mechanism remains implementation detail.
- **why this owner/stage:** this is precisely the seam R1-F3 already owns; it corrects the fence *law*, not its implementation, so it cannot be deferred to D8/implementation without leaving the law unfalsifiable against its dominant threat. It requires no new mechanism (object store and deployment state already exist in D7-E).

#### F-2 — selected `chi-server` generation mode has never been executed against the canonical OAD

- **classification:** `LATER_NON_BLOCKING`
- **severity:** Minor
- **location:** D7-E §3.2 (`chi-server: true`, `strict-server: true`) vs `scripts/verify-product-oad-baseline.mjs:440` (executed gate proof uses `std-http-server: true` + `strict-server: true`).
- **governing authority:** D7-E §3.2/§14; R1-F5 proof-timing contract.
- **counterexample/failure:** none today — under R1-F5 this is correctly a category-B obligation. But the marginal cost of converting it to category-A evidence is near zero: the gate already runs `oapi-codegen v2.8.0` against the resolved bundle in Go; adding a second config with `chi-server: true` would prove the exact selected generator×mode×document combination (including `{id}:verb` route emission) deterministically generates and compiles, before implementation opens.
- **smallest correction:** optional gate/candidate-proof extension only; no authority text change required. Do not treat as a closeout blocker.
- **why:** cheap falsifier-hardening inside the existing executable proof lane; belongs to gate scope, not D7 authority.

#### F-3 — bootstrap byte budget at 93%

- **classification:** `REPOSITORY_FIX`
- **severity:** Minor
- **location:** gate output `bootstrap_bytes: 19051 / 20480`.
- **failure:** D8 opening will add roadmap/index route rows; the fresh-actor bootstrap cap has ~1.4 KB headroom. First D8 gate may go red on a pure routing edit.
- **smallest correction:** none now; when D8 opens, trim frozen D7 detail from `docs/roadmap.md` (the “Accepted D7 baseline” block can compress once D7 is CLOSED) rather than raising the cap.

No other material finding exists. I did not manufacture style findings.

**Adversarial checks that came back clean (recorded so adjudication does not re-litigate them):**

- **Suspected OAD-validator keyword blind spot — refuted with primary evidence.** The canonical OAD is `openapi: 3.1.2` with pervasive `const` discriminators (241 `additionalProperties: false`, `const` unions in `components.yaml`). `oapi-codegen v2.8.0` depends on `kin-openapi v0.142.0`, and that exact version implements `const` parsing **and** validation (`openapi3/schema.go:136` `Const` field; `visitConstOperation` at `schema.go:2091`). The same loader already parsed the 3.1.2 bundle successfully in the executed gate. The selected `nethttp-middleware` validator therefore can enforce the document as written; residual behavioral proof remains a correctly-named implementation falsifier.
- **`x-mpc-*` operation-policy extensions exist as claimed:** 104 instances each of `x-mpc-operation-class`/`x-mpc-required-permission`/`x-mpc-principal-kinds`/`x-mpc-semantic-owner` plus 2 `x-mpc-required-physical-qualification` across the multi-file OAD — D7-E §3.3’s generated-projection source is real, not aspirational.
- **`GET /_auth/csrf` is not a hidden surface expansion:** D5-R1 explicitly delegates “CSRF secret generation, carriage/bootstrap and rotation” to D7 (D5-R1 lines 77, 210), and D5-B2 Technical Ingress already names the technical presentation surface as outside Lanes A/B. No Product operation or Permission is added.
- **FK-vs-RLS existence-oracle check:** composite FKs `(organization_id, parent_id)` mean referential probes cannot confirm foreign-Organization rows — the Organization dimension is inside the constraint itself. No cross-tenant existence leak through RI internals.

### 5. Persistence / isolation / transaction assessment

Sound. The four-mode scope taxonomy after R1-F1 (`organization` / `principal_self` / `authentication_bootstrap` / `technical_routing` + explicit platform-global technical state) is exhaustive for every access the composed D7-D/C flows actually need — I could not construct a lookup that lacks an honest mode. `authentication_bootstrap` is tightly bounded (three enumerated lookups, no business/evidence tables, ends at Principal resolution) and does not become a cross-tenant bypass: it touches only authentication mechanism state that is not Organization-owned. Transaction-local `set_config(..., is_local=true)` + ENABLE/FORCE RLS + non-owner/NOSUPERUSER/NOBYPASSRLS runtime role + composite FKs is a mutually reinforcing stack with default-deny on missing scope; pooled-connection scope leakage is structurally excluded and named as falsifier 7. `READ COMMITTED` + `SELECT … FOR UPDATE` + same-transaction revision rotation correctly prevents double-spend of one revision without a global `SERIALIZABLE` tax. Idempotency uniqueness `(organization_id, operation, digest(key))` agrees with owner transaction ownership; the intake-order algorithm resolves replay before stale-revision evaluation per accepted W2. River engine tables as platform technical state with no business authority is correct and consistent with D7-C. No accepted state shape requires a cross-owner shared transaction.

### 6. Durable work / external-effect assessment

Sound, with F-1 as the single qualification. `InsertTx` gives real atomic owner-state→work handoff with no second outbox — correct, because everything lives in one PostgreSQL. All nine required crash windows resolve correctly on a continuous timeline: the pre-dispatch marker partitions the world into safely-retryable (no marker), replayable (marker + definitive outcome) and reconcile-only (marker + no outcome), and rescue/duplicate/out-of-order delivery is absorbed by owner semantic idempotency with River uniqueness demoted to optimization. Retry classes R1–R5 preserve D4’s no-blind-retry law, with R5 correctly failing toward ambiguity. Reconciliation never re-issues the original write without a new admitted attempt. A notable strength: because River state, markers and owner state share one database, PITR rolls them back *consistently* — the only inconsistency is with the external world, which is exactly what the R1-F3 fence addresses. The fence’s step logic is sufficient and minimal **once engaged**; its engagement is the F-1 gap.

### 7. Authentication / request-trust assessment

Sound. Confidential code + PKCE S256 with one-time persisted login transaction handles replica-split callbacks; fresh post-login handles kill fixation; `(issuer, sub)`-only binding avoids claim-merge attacks; discarding the human OIDC token set after callback removes an entire secret class with no current consumer — good YAGNI. Digest-only session storage neutralizes the DB-leak fixture. Per-request eligibility recheck makes disabled-Principal sessions dead immediately without back-channel logout infrastructure. CSRF = session-bound synchronizer token (memory-only carriage via authenticated `GET /_auth/csrf`) plus `net/http.CrossOriginProtection` defense in depth; no wildcard CORS. Machine path: jwx/v3 against pinned issuer/JWKS with audience `https://conexus.fun`, unknown-`kid` bounded refresh (no token-supplied JWKS), and the explicit machine-client→A/S binding means cryptographic validity alone never grants access; no-binding/multi-binding/H-resolution all fail closed. Dual-carrier requests fail closed as ambiguous (R1-F4) — this also closes the cookie+bearer CSRF-confusion corner. IdP outage does not invalidate existing sessions (no retained IdP tokens) and does not fake MPC authorization. Auth bootstrap is assessed under §5.

### 8. HTTP / OAD runtime assessment

Sound. Chi v5’s parameter-node-with-static-suffix routing models `{id}:verb` without touching W1 grammar (D5 proved `ServeMux` panics on this class; the executed gate’s custom-mux wrapper proof was D5 expressibility evidence, not a router selection — no contradiction with selecting Chi now). Reusing the D5-pinned `oapi-codegen v2.8.0` with `chi-server`+`strict-server` avoids a second generator; generated operation-policy metadata from real `x-mpc-*` extensions eliminates the handwritten Permission/route map drift class. The R1-F4 composition order (route → carrier auth → CSRF → scheme-aware validator over established context → policy metadata → owner gates) is realizable as ordinary middleware nesting, forbids `NoopAuthenticationFunc` in production, reimplements no cryptography in the validator, and translates validator errors into accepted Problem grammar. Validator keyword capability against the 3.1.2 document is evidenced (see §4 clean checks). 405+Allow and negative controls (carrier removal, scheme mismatch, dual carrier) are correctly named as executable falsifiers. Technical ingress and media delivery stay out of the Product OAD. F-2 notes the one cheap category-B→A conversion available.

### 9. Storage / migrations / operability / recovery assessment

Sound, with F-1 qualifying recovery. Object-first/DB-reference-second ordering, forbidden overwrite of committed keys, digest/size/type retention, orphan sweep with age bound, and authenticated delivery that rechecks full Product authorization before streaming — coherent byte custody with no Asset domain. R1-F2’s migration split is neither duplicated nor incomplete: `tern` owns MPC schema, the pinned River CLI owns River schema, one release step orders both before deploy, runtime never holds the migration credential, and failed lanes stop roll-forward; readiness on incompatible schema closes the skew window. Config/secrets: typed startup config, no CLI-flag secrets, file-mount preference, explicit redaction fence covering every secret class D7-D creates. slog + OTel traces/metrics over OTLP/HTTP matches current OTel Go signal maturity (logs correctly kept out); telemetry failure cannot block owner transactions or readiness. One immutable digest-pinned OCI image, trusted-proxy-only forwarded headers, liveness/readiness split with the provider-false-death prohibition — all proportionate. Backup/restore names the right four coherence domains (PostgreSQL PITR, Keycloak subject continuity, River state, binary integrity) and correctly requires restore drills, not backup existence.

### 10. Product / authority non-regression assessment

**Intact.** Independently measured on the exact candidate: 99/99 Product operations, 30/30 ordinary Permissions, H/A/S only, stable origin `https://conexus.fun`, H session+CSRF and A/S Client Credentials bearer exactly as accepted (gate auth profile `OIDC_SERVER_SESSION_COOKIE_CSRF` / `CLIENT_CREDENTIALS_BEARER_A_S`, 5/5 auth negative controls). No hidden Product operation, Permission, client class, business owner, workflow/status entity or provider-shaped vocabulary appears in any D7 document. Sankhya remains API-Gateway-only, with Direct Oracle fallback named as a falsifier (D7-C §15.16). Semantic ownership and D3 Q/C/E/P meaning are consumed, not amended. Implementation remains blocked until D9 and the candidate touches no implementation surface (gate forbidden-diff lanes confirm).

### 11. YAGNI / dependency assessment

Every selected mechanism has a present consumer: Chi (`{id}:verb` + surface separation), pgx/pgxpool (canonical PostgreSQL, no ORM), RLS+composite FKs (D2 structural isolation), River (D3/D4 durable reactions + effects, no broker), Keycloak + go-oidc/x-oauth2 (D5-R1 human flow), jwx/v3 (machine bearer), oapi-codegen strict server + nethttp-middleware (D5 OAD as sole wire authority), S3-compatible custody (D5-admitted media/artifacts), tern (SQL-first migrations where RLS/FK SQL is the correctness surface), slog/OTel (bounded operability), one OCI image (D7-A topology). Two JOSE stacks coexist (go-oidc for ID tokens, jwx for machine JWTs) — deliberate and defensible since go-oidc is ID-token-centric and abusing it for resource-server validation would be worse than a second well-scoped library; not a finding. No deferred mechanism is already required: I attempted to prove Redis (session scale), a broker (fan-out), an outbox (non-PG transport), Vault (secret rotation) and Kubernetes (replica orchestration) necessary and each attempt died against an actual current property. The rejection list is honest. R1’s own repairs introduce no overengineering — the fence (with F-1) is driven by a real falsifier, not symmetry.

### 12. Proof assessment

**A. Current executable repository proof:** green CI on `c08a4d02` proves document/link/bootstrap governance, docs/work isolation (with a working negative control), roadmap-truth markers, OAD structure (99/30/H-A-S, 28 List/Search, 12+7+5 negative controls, idempotency carriers), deterministic TS+Go projection via real `oapi-codegen v2.8.0` execution (std-http-server + strict), and honest non-claims (`legacy runtime population 0`, `runtime schema enforcement NOT_CLAIMED_D7`). It proves nothing about runtime behavior — and says so.

**B. D7 architecture proof contract:** the composed falsifier set (D7 §10, D7-B §13, D7-C §15, D7-E §13, R1 falsifiers) is genuinely executable later with real dependencies — every claim names its real seam (PostgreSQL, River, Keycloak, browser, router/validator, object store) and none requires infrastructure the architecture excludes. With the F-1 amendment, the recovery-fence falsifier also covers its dominant threat path. I found no tautological control and no obligation that cannot realistically be executed.

**C. Implementation acceptance proof:** correctly gated post-D9 by R1-F5’s three-status contract; the reinterpretation of ratified slice “cannot close without proof” prose is explicit, bounded and honest. Mock-promotion is expressly forbidden for RLS/River/Keycloak/browser/router/object-store claims.

Material blind spots: only F-1 (now named). F-2 is an optional cheap A-category upgrade.

### 13. Reconstruction decision

**No.** There is no material reason to reconstruct or reopen D0–D6 or D7-A→E. All five R1 seams are real and correctly bounded; the one surviving gap (F-1) is a sentence-level completion of R1-F3 itself, at exactly the authority that already owns it.

### 14. Continuation recommendation

Smallest exact action after GPT adjudication: apply the F-1 amendment to D7-R1 §4 (continuity-affirmative fence-arming law + unarmed-restore falsifier) on `stage/d7-runtime`, run the fresh exact-head gate, then present D7 for operator closeout decision. F-2/F-3 need no authority change and must not block. This review authorizes no merge, no D8, no D9 and no Product implementation.

*— Fable, independent whole-D7 challenge, 2026-08-22, on `review/d7-fable` against exact candidate `c08a4d025cfd89269cc071f4b307695e79f6f8cb`.*