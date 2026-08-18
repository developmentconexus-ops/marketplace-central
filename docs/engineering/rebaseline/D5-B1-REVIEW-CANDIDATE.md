# D5-B1 — Semantic API Model & Contract Laws — CONVERGED REVIEW CANDIDATE

> **Status:** NON-AUTHORITATIVE — Fable reviewed / GPT adjudicated / awaiting explicit operator ratification  
> **Stage:** D5 — API  
> **Batch:** B1 — Semantic API Model & Contract Laws  
> **Parent authority:** accepted D0 → D4 plus `ARCHITECTURE.md` and ADR registry  
> **Method:** DevelopmentConexus Engineering Method v1.0.0  
> **Review evidence:** `AI-DIALOG.md`, Fable D5-B1 Independent Adversarial Review, 2026-08-18  
> **Purpose:** converged review package only. This file does not change architecture authority, router status, ADR status, or implementation permission.

---

## 1. Question and boundary

D5 must not begin by preserving, renaming or mechanically generating the current API.

The first material question is:

> **What external API model and wire laws let MPC clients interact with the business authorities already accepted in D0–D4 without turning legacy modules, provider protocols, transport mechanics or generic mutation abstractions into product semantics?**

The accepted product boundary remains unchanged:

- MPC is the internal **Marketplace Operations Control Plane** for **Marketplace Operations + Commercial Intelligence**;
- Mercado Livre, Sankhya, payment APIs and future external systems participate in that operating loop but do not define MPC ontology;
- MPC is provider-independent and integration-independent;
- MPC is **not** currently a generic enterprise operations platform;
- making marketplace commerce merely one vertical of a broader enterprise control plane requires a targeted D0 product-boundary reopen followed by D1 coherence/reopen as needed; D5 may not introduce that change through generic API naming.

This batch decides **API laws before endpoint inventory**.

### In scope

- client-facing semantic API model;
- Product API versus provider/business-system protocol ingress;
- D1 owner and D3 Q/C/P wire laws;
- Organization and Principal/access scoping;
- knowledge, freshness and provenance representation laws;
- consequential-action outcome, idempotency and precondition laws;
- HTTP/API problem semantics versus valid business outcomes;
- source-qualified external identity and provider-rich evidence containment;
- one machine-readable Product API wire authority and derived/conformant SDK/server contracts;
- hard-cutover / compatibility posture;
- bulk/partial-outcome admission laws;
- architectural enforcement and proof obligations.

### Explicitly not in scope

- complete endpoint/path/operation inventory;
- D6 frontend screens, packages, query/cache topology or UX;
- D7 worker/queue/outbox/retry/transaction/deployment/runtime topology;
- D8 golden-flow execution/proof choreography;
- product implementation, still blocked until D9;
- concrete OpenAPI generator, Go router/framework or TypeScript HTTP runtime;
- exact OIDC provider/session/deployment topology;
- universal pagination/filter/sort DSLs;
- universal event-stream/SSE/WebSocket surface;
- generic workflow, command, mutation, integration, evidence or provider framework.

---

## 2. Evidence classification

### KNOWN

1. D0 defines MPC as a marketplace-operations control plane, not an ERP replacement, marketplace dashboard or generic integration hub.
2. D1 assigns one semantic authority to each material Product 1.0 responsibility and explicitly rejects a generic `Mutation` / `Action` business owner.
3. D2 makes `Organization` the canonical tenant/isolation root, defines `Principal`, source-qualified external identity, Marketplace Installation, SourceInstance and domain-local identities for material Business Intents.
4. D2 forbids Organization scope from being inferred from Marketplace Installation, Selling Entity, external account, IdP organization, source key or process-global defaults.
5. D2 denies cross-Organization references between Organization-owned business state by default.
6. D3 accepts semantic Q/C/E/P communication, explicit Organization scope, semantic idempotency, honest query knowledge state and accepted/rejected/pending/ambiguous capability semantics where applicable.
7. D3 makes freshness orthogonal to known/unknown state and requires owner-controlled provenance/observation time when freshness-for-use is material.
8. D4 requires consumer-owned semantic ports, provider-local protocol/DTO knowledge, explicit namespace/capability/coverage semantics, source-qualified external evidence and acceptance/ambiguity/convergence separation for external effects.
9. D4-B4 accepts **Semantic Core + Provider-Enriched Evidence** rather than lowest-common-denominator flattening or provider-payload mirroring.
10. `ARCHITECTURE.md` and the Evidence Register establish no production compatibility obligation to the current routes/schema/SDK; hard cutover is allowed.
11. Current OpenAPI/routes/SDK are evidence only, not inherited target authority.
12. OpenAPI, a hand-written SDK and handler/routing knowledge currently coexist. Legacy ADR-016 itself admits that same-commit discipline proves atomicity, not agreement, and leaves manual transcription drift reachable.

### INFERRED

- The current external contract is a historical aggregation rather than a coherent D0–D4 semantic API.
- Generating a client directly from the current OpenAPI would automate legacy semantic defects rather than resolve them.
- Endpoint inventory before boundary/wire laws would create local decisions likely to be redone after applying D1–D4 authority consistently.

### UNKNOWN

Still unknown until later D5 work:

- exact Product 1.0 operation set and final path spelling;
- exact D6 read models;
- exact Permission → API-operation catalog;
- per-operation pagination/filter/sort/cursor contracts;
- concrete provider-rich fields exposed to clients;
- concrete bulk operations;
- concrete operations requiring optimistic concurrency;
- exact authentication wire details above the accepted OIDC boundary;
- concrete SDK/generator/server-conformance technology;
- whether a separately named technical/admin surface is needed;
- number of remaining D5 batches.

Unknown stays Unknown. Deferral never authorizes a plausible default.

---

## 3. Root cause and target invariant

### Root cause

The structural defect is not primarily route count or naming.

> **The external client contract is not currently derived from the accepted semantic authorities, while product operations, provider/business-system protocol concerns and multiple manually maintained wire representations coexist without one explicit wire authority.**

That defect class permits provider vocabulary to become ontology, generic Mutation semantics to steal domain intent, unknown/unavailable to collapse into defaults, provider 2xx to masquerade as convergence, unsafe retry after ambiguous acceptance, namespace ambiguity, and OpenAPI/SDK/handler drift.

### Governing invariant

> **Every externally invokable MPC Product API operation belongs to exactly one accepted semantic owner or accepted non-domain identity/access authority; its contract preserves explicit Organization scope, source-qualified identity, honest knowledge/freshness/effect semantics and ownership boundaries; provider/business-system protocol remains outside Product API semantics; and the Product API has one machine-readable wire authority from which supported client contracts derive and against which server behavior conforms.**

Corollaries:

1. HTTP shape does not create business authority.
2. Generic technical machinery does not become a generic business owner.
3. Product API vocabulary is MPC business language.
4. Provider-native evidence may be exposed only when semantically qualified and materially useful; raw provider DTO ontology does not cross the boundary.
5. A valid response never makes a stronger knowledge, freshness, authorization or effect claim than the owning authority can justify.
6. Client retry convenience never weakens D3/D4 ambiguity or no-blind-retry safety.
7. No second manually authoritative wire representation is admitted.

---

## 4. Credible alternatives

### A — Clean current OpenAPI + retain hand-written SDK

Local Maximum. Leaves duplicate manual wire authority and lets legacy topology constrain target semantics. **Rejected.**

### B — Generate SDK directly from current OpenAPI

Removes one drift mechanism but automates the wrong semantic contract. Useful only after semantic redesign. **Rejected as first-order solution.**

### C — Generic platform API (`resources`, `commands`, `mutations`, generic provider/workflow/evidence models)

Erases differentiated D1 meaning and introduces framework authority without a second real product/vertical. **Rejected by authority + YAGNI.**

### D — One external API per D1 domain

Preserves semantic owners but prematurely converts business boundaries into external/runtime fragmentation. **Rejected as over-partitioning.**

### E — Semantic MPC Product API + separate provider protocol ingress + one wire authority

Removes the root cause without genericizing the product or choosing D6/D7 topology. **Selected Global Maximum.**

### F — GraphQL/gateway/gRPC/event-stream as the new architectural center

No accepted consumer/failure class requires this and it does not solve semantic ownership or effect-safety defects by itself. **Rejected for now.**

**Outcome relative to the current API:** `RESTRUCTURE NOW`.

---

# 5. Converged decision set

## B1-D1 — Product API is semantic and domain-oriented

The client-facing API represents the accepted MPC marketplace-operating model.

It is:

- **provider-independent** — Mercado Livre does not define API ontology;
- **business-system-independent** — Sankhya taxonomy/choreography does not define API ontology;
- **integration-independent** — adapters/protocols do not become Product API business resources;
- **domain-oriented** — accepted D1 meanings remain distinct;
- **not enterprise-generic** — marketplace commerce remains the accepted product boundary until explicitly reopened.

`provider-independent` does not mean lowest-common-denominator. D4-B4 remains binding.

---

## B1-D2 — Two external audiences; no speculative third API

### MPC Semantic Product API

Consumed by MPC applications and supported SDK clients. It may expose:

- owner current meaning;
- legitimate read projections/compositions;
- MPC-owned resources/configuration;
- owner capabilities and domain-owned Intents;
- Governance and Operational Work interactions under their accepted semantics;
- qualified provider-enriched evidence when a named consumer/correctness property requires it.

### Provider / business-system protocol ingress

Webhook/callback receivers, OAuth callbacks/handshakes and equivalent protocol endpoints:

- may be externally reachable;
- are not Product API business operations;
- are not part of the normal Product SDK;
- accept provider vocabulary only inside the D4 protocol boundary;
- do not make notification payload domain truth.

Health/readiness/admin/technical surfaces are not pre-created by taxonomy. A real D7/runtime/operator consumer must justify them.

---

## B1-D3 — REST/HTTP resource semantics where truthful; explicit owner operations where CRUD would lie

The target Product API is an HTTP/REST semantic API described by OpenAPI.

Use ordinary resource semantics only when the client is genuinely interacting with MPC-owned state/configuration whose meaning fits those semantics.

Do not pretend a direct-state `PUT` completed an externally authoritative world when MPC only accepted an intent.

Consequential actions remain owned by their D1 domain and, where D2 requires durable identity, by owner-specific Business Intents such as Listing/Price Intent, Availability Intent, Business Order Intent, Invoicing Intent and material Fulfillment intents.

No target Product API business resource named generic `Mutation`, `Action`, `Command`, `Operation` or universal `BusinessIntent` exists merely to normalize these lifecycles.

Explicit owner operations are permitted when CRUD would misrepresent meaning. They remain in the owner's vocabulary rather than becoming a generic action escape hatch.

---

## B1-D4 — D3 Q/C/P semantics survive the wire

| D3 form | Product API meaning | Binding wire law |
|---|---|---|
| **Q** | ask owner for current owner-owned meaning | preserve known/known-empty/unknown/unavailable/partial as applicable; when freshness-for-use is material, expose/reference owner-controlled observation/acquisition/provenance time; cache/projection/HTTP time cannot impersonate owner freshness |
| **C** | ask owner to accept/perform owner-owned work | return owner semantic outcome, never raw provider transport status |
| **P** | read composition of multiple authorities | read-only; projection never becomes write/concurrency authority; component freshness/partiality remain honest and projection `updated_at` never substitutes for source observation time |
| **E** | committed owner fact for independent reaction | not automatically a public stream; external event/stream API requires a named consumer and later D5 decision |

D5 exposes product interaction, not the internal communication topology by default.

---

## B1-D5 — Organization is path-scoped and fail-closed

For every Organization-owned Product API operation, the contract is scoped under:

```text
/organizations/{organization_id}/...
```

This is **decided**, not merely a leading candidate.

Reasons:

1. D2 requires Organization to be explicit.
2. Required path scope removes ambient/token-only/default-tenant dependency.
3. Organization becomes part of operation identity, links, logs and URL-keyed cache/client artifacts.
4. Missing scope fails routing closed.
5. The first one-Organization proof does not hard-code a singleton into the contract.

The path is a scope claim, never self-authorization. The authenticated Principal must have current ordinary access to that Organization.

### Secondary-reference fence

Any Organization-owned identity/reference supplied in path, query or body must resolve inside the path Organization. A request scoped to Organization A cannot smuggle a Sale, Offering, Intent, Work item or other Organization-B-owned reference through a secondary field. Cross-Organization reference fails closed.

### Provider-ingress organization rule

Provider ingress does **not** compute or infer Organization from provider data.

The protocol boundary:

1. identifies the bound Marketplace Installation / SourceInstance fail-closed from authenticated/configured context plus authoritative provider/source markers where available;
2. reads Organization from that identity's explicit MPC-owned binding;
3. fails closed on missing, ambiguous or contradictory binding/markers before attribution;
4. records Organization explicitly in any durable acquisition/recovery state that can outlive the ingress execution context.

Installation/SourceInstance qualify namespace **inside** explicit Organization scope; they never substitute for Organization.

---

## B1-D6 — Principal/access and consequential authorization remain separate

At the Product API boundary:

- AuthN binds the Principal through the accepted OIDC boundary;
- current Membership/Permission determines ordinary invocation/view access;
- caller-supplied `principal_id`, role names or `approved=true` never substitute for those proofs.

`401/403` are AuthN/ordinary-access problems.

They are not:

- domain business rejection/prohibition;
- approval-required disposition;
- Governance approval/rejection;
- proof of consequential authorization.

Action owners retain business disposition/validity. Controlled Action Governance retains consequential grant/delegation/Authorization Decision semantics.

---

## B1-D7 — Knowledge and freshness are semantic output

Where material, Product API contracts preserve enough structure to distinguish:

- known value;
- known empty/absent;
- unknown / insufficiently known;
- unavailable;
- partial / incomplete coverage.

No universal `Fact<T>` JSON envelope is required.

> **No null, zero, `false`, empty object/list or HTTP success code may silently collapse a materially different knowledge state.**

Freshness is a separate axis:

> **Where freshness-for-use is material to a client decision, the response exposes or references owner-controlled source/effective/observation/acquisition/provenance time sufficient for that use. HTTP exchange time, `known` status and projection update time never substitute for source freshness.**

A successful HTTP exchange may legitimately carry semantic `unknown`, `unavailable` or `partial` meaning.

Shared discriminators are reused only where semantics are genuinely identical; no universal Evidence/Knowledge graph is introduced.

---

## B1-D8 — Capability outcomes are not HTTP/API problems

For a contract-valid owner capability request, applicable semantic outcomes include:

- **accepted**;
- **rejected**;
- **pending**;
- **ambiguous** possible acceptance.

Binding distinctions:

```text
accepted != completed
completed != externally applied
externally applied != converged

rejected != access denied
pending != failed
ambiguous != failed
ambiguous != safe-to-retry
```

A provider timeout after possible dispatch cannot become failed/rejected merely for API convenience.

Two precondition classes remain distinct:

1. API representation/concurrency conditional failure may be an HTTP problem;
2. business/provider execution precondition belongs to the owner outcome/intent semantics after provider translation.

---

## B1-D9 — Consequential intake is idempotent by default and fails closed without its key

For every Product API operation that can create a durable Business Intent or initiate a consequential external effect, HTTP intake requires an `Idempotency-Key` by default.

Missing key fails explicitly **before durable intake/effect**.

A per-operation exemption is allowed only when the operation contract proves duplicate requests are structurally unreachable or harmless through owner-anchor/resource semantics alone. The exemption and its proof are recorded in the later operation inventory; optional-by-convenience is not allowed.

The key is request-deduplication mechanism, not:

- canonical Business Intent identity;
- provider idempotency proof;
- authorization;
- permission to replay an ambiguous external write.

Laws:

1. same Organization + semantic operation + key + semantically equivalent request resolves the same MPC intake result / durable Intent;
2. same key + materially different semantic request fails explicitly;
3. retry after local uncertainty resolves existing intake before another Intent can be created;
4. if a domain Intent exists, retry resolves that owner Intent rather than a generic Mutation record;
5. external ambiguity still follows D3/D4 reconciliation and never authorizes blind redispatch.

Exact storage, retention, locks and cleanup are D7 mechanisms.

---

## B1-D10 — Optimistic concurrency only where stale state matters

D5 does not impose version tokens everywhere.

When client action materially depends on current MPC-owned version/state, use an opaque MPC-level concurrency/precondition token with HTTP conditional semantics where suitable.

Provider-native version tokens remain adapter-local unless a D1-owned semantic contract genuinely requires qualified exposure.

Stale provider version responses are translated into owner semantics + authoritative reread/redecision. Projection state is not consequential concurrency authority.

---

## B1-D11 — One standard HTTP problem model; business outcomes stay outside it

HTTP/API-level failures use **RFC 9457 Problem Details** (`application/problem+json`) plus only MPC-owned stable extensions required by real programmatic consumers.

API-problem classes include as applicable:

- malformed/unparseable request;
- schema/contract validation;
- authentication;
- ordinary Permission/access denial;
- unsupported operation/HTTP contract;
- API concurrency/conditional failure;
- missing/misused/conflicting idempotency key;
- unexpected MPC server failure.

A stable MPC machine-readable `code` may specialize the problem type/extensions where real handling requires it.

These are **not** Problem Details merely because they are undesirable:

- business `rejected`;
- approval-required/pending Governance decision;
- `ambiguous` possible external acceptance;
- valid query `unknown`/`unavailable`/`partial`;
- admitted `unsupported` / `external-required` business/provider capability outcome.

Raw provider/business-system error DTOs, arbitrary text, secrets and PII never become Product API problem truth. A named support consumer may receive sanitized source-qualified diagnostics without replacing MPC semantic outcome.

---

## B1-D12 — External identities and provider-rich evidence stay source-qualified

### External identity law

A provider/native identifier crosses the Product API only when its qualifying MPC namespace identity is explicit in the contract:

- Marketplace Installation for marketplace-native identity; or
- SourceInstance where that is the accepted source qualifier.

The qualifier may be explicit in the serialized identity shape or be unambiguous from the operation's declared scope. A bare external identifier is never a Product API correlation key.

Example correctness test: two Installations in one Organization may hold identical-looking native IDs and must remain distinguishable from the client contract alone.

This does not create a universal `ExternalReference`/entity graph. Owner contracts use the smallest typed source-qualified representation their semantics require.

### Provider-rich evidence law

A D1 owner may expose provider-specific enriched evidence only when it serves a named Product 1.0 consumer/correctness property.

The enrichment remains:

- source-qualified;
- bounded inside the owner contract;
- optional/unsupported/not-applicable/unknown when another provider lacks an equivalent;
- prohibited from becoming top-level MPC owner/resource ontology merely because one provider exposes it.

Bounded discriminated enrichment unions are allowed. Arbitrary raw-payload passthrough is rejected.

---

## B1-D13 — OpenAPI is the single machine-readable Product API wire authority

Accepted D-stage artifacts and `ARCHITECTURE.md` remain semantic architecture authority. OpenAPI never outranks them.

Within the HTTP Product API boundary:

> **One OpenAPI document/set is the machine-readable authority for paths, operations, parameters, headers, status codes and serialized schemas.**

Rules:

1. **No second manually authoritative wire representation is admitted.**
2. Supported SDK/client contracts are mechanically derived from OpenAPI or mechanically proven to conform to it.
3. Server request/response/status/header behavior is mechanically validated against the same contract.
4. Implementation/client disagreement with OpenAPI fails verification.
5. OpenAPI disagreement with accepted D-stage semantics is an OpenAPI defect.
6. Provider ingress may have separate executable protocol contracts when useful but does not contaminate Product SDK/authority.
7. Conformance controls count only when a negative/drift fixture demonstrates that they actually fire.

Exact generator/router/runtime tooling remains later realization.

### ADR-016 disposition upon canonical B1 consolidation

ADR-016 becomes **historical** when this batch is ratified and consolidated.

Its old target mechanism — manual OpenAPI + manual SDK same-commit synchronization — is intentionally removed rather than strengthened.

Two lessons survive as active D5 invariants:

1. never maintain a second manually authoritative wire model;
2. conformance/drift controls must be demonstrably capable of failing.

No compatibility transition window needs the old rule because Product implementation is still blocked until D9 and no production client is entitled to the legacy SDK/API.

---

## B1-D14 — No compatibility/versioning machinery without a consumer

There are no production clients entitled to the current API.

Target cutover may delete/rename/replace routes, resources, schemas, generic Mutation surfaces and the manual SDK without compatibility aliases or dual-version infrastructure.

A literal `/v1` path segment is not itself prohibited, but it carries no compatibility policy by implication.

A real future compatibility entitlement is material new evidence and reopens this decision.

---

## B1-D15 — Bulk is operation-local, never a universal primitive

No generic batch/mutation envelope exists.

A bulk endpoint is admitted only when a named Product 1.0 workflow/consumer materially requires bulk semantics beyond adequate composition of individual operations.

If admitted:

- intended scope remains owner-defined;
- authorized scope remains distinct;
- attempted/member outcome scope remains distinct;
- confirmed/rejected/ambiguous/not-executed members remain distinguishable where material;
- one ambiguous member cannot turn the whole batch into safe blind replay;
- provider import/batch IDs remain external reconciliation evidence, not generic MPC Batch identity by default.

---

# 6. Wire-semantics matrix

| Interaction | Authority | Valid semantic result | Intake/precondition law | Must never imply |
|---|---|---|---|---|
| Owner query **Q** | D1 owner / accepted D2 substrate | known / known-empty / unknown / unavailable / partial + freshness/provenance when material | read retry subject to owner freshness/coverage semantics | empty/default fabricated as knowledge; HTTP time == freshness |
| Projection **P** | no new authority | composed read with honest component freshness/partiality | read-only | projection update time == source freshness; projection becomes write authority |
| Ordinary MPC-owned config/resource | owner/substrate | created/updated representation | resource idempotency + conditional update only where material | MPC config write == provider state changed |
| Consequential owner capability **C** | action owner | accepted / rejected / pending / ambiguous as applicable | `Idempotency-Key` mandatory by default; owner/precondition semantics preserved | provider 2xx == converged; timeout == failed; generic Mutation owns intent |
| Governance | Controlled Action Governance | pending / authorized / rejected/invalidated under Governance meaning | decision identity/history preserved | 403 ordinary access == Governance rejection |
| Operational Work | Work lifecycle owner + originating owner for source truth | Work state and source-resolution meaning remain separate | Work-local idempotency | closing Work mutates source truth |
| Provider ingress | D4 protocol boundary | acquisition/protocol result before owner translation/reread | provider-specific auth/duplicate mechanics + explicit Organization binding | callback DTO becomes Product API/domain truth |

---

# 7. Essential versus accidental complexity

## Essential complexity preserved

- Organization isolation and cross-Organization reference safety;
- Principal attribution and ordinary access;
- D1 semantic authority;
- source-qualified external identity;
- known/empty/unknown/unavailable/partial distinctions;
- freshness/provenance/time when material;
- domain-local consequential Intent identity;
- business disposition versus Governance authorization;
- fail-closed consequential intake idempotency;
- precondition/concurrency safety where stale state matters;
- accepted/rejected/pending/ambiguous semantics;
- acceptance/completion/application/convergence separation;
- member-level partiality for real bulk;
- provider-enriched evidence where materially useful;
- provider diagnostic containment.

## Accidental complexity removed/refused

- generic Mutation/Command owner;
- provider/business-system paths as Product API ontology;
- duplicate manual OpenAPI + SDK authorities;
- unused legacy compatibility;
- generic provider/integration/resource graph;
- universal command/capability response envelope;
- universal `Fact<T>`/Evidence JSON model;
- speculative multi-version infrastructure;
- generic bulk framework;
- GraphQL/gRPC/gateway/event-streaming without a named need.

---

# 8. Enforcement candidates

D5 freezes protected properties; D7/implementation later chooses exact mechanisms.

1. **Operation-owner map:** every Product API operation declares exactly one semantic owner/accepted substrate authority plus Q/C/P class.
2. **Route classification:** every externally reachable route is Product API, provider ingress or a separately justified technical surface; unclassified routes fail review.
3. **Provider-vocabulary fence:** provider nouns cannot define Product API ownership/path semantics except qualified enrichment/protocol ingress.
4. **Source-qualified identity fixture:** identical native IDs under two Installations/SourceInstances remain client-distinguishable.
5. **OpenAPI authority:** supported client contract reproducibly derives/conforms; no hand-written second authority.
6. **Server conformance:** request/response/status/header behavior validates against admitted OpenAPI.
7. **Knowledge/freshness fixtures:** unknown/unavailable/partial cannot pass as known-empty/default; freshness-material output without owner-controlled provenance fails contract review; projection `updated_at` cannot substitute for component/source time.
8. **Outcome fixtures:** ambiguous cannot pass as rejected/failed; accepted cannot pass as converged.
9. **Idempotency fixtures:** missing key on consequential intake fails before durable state/effect; same key + different semantic request fails; same key + same request does not create another Intent/effect; any exemption carries explicit owner-anchor proof.
10. **Access/authorization fixtures:** ordinary Permission may allow invocation while domain/Governance still rejects/pends; business rejection cannot become ordinary 403.
11. **Tenant path fixture:** changing path Organization cannot grant cross-Organization access.
12. **Reference-smuggling fixture:** Organization-A request carrying Organization-B secondary reference fails closed before state change/effect.
13. **Ingress binding fixtures:** missing/ambiguous/contradictory provider marker→Installation/SourceInstance binding fails closed; durable acquisition without explicit Organization fails structural validation.
14. **Provider diagnostic redaction:** raw PII/secrets/error payload cannot cross Product API accidentally.
15. **Bulk partiality:** confirmed + ambiguous members cannot become one safe-to-retry batch failure.

---

# 9. Proof strategy before implementation

### P1 — Complete operation ownership

Every later D5 operation must classify:

```text
operation
  -> semantic owner / accepted substrate
  -> Q | C | P
  -> Organization scope
  -> Permission requirement
  -> knowledge/outcome class
  -> source-qualified identity needs
  -> freshness/provenance needs
  -> idempotency/precondition requirement
```

Any operation requiring a new owner returns to the implicated parent stage rather than creating an API exception.

### P2 — Protocol separation

Use concrete current/provider examples to prove protocol nouns terminate at D4 and translate into accepted owner semantics.

### P3 — Knowledge/freshness counterexamples

Falsify at least:

- unavailable owner/source;
- completed known-empty versus incomplete enumeration;
- materially stale but known value;
- projection whose update time is newer than source observation time;
- provider enrichment unsupported on another provider.

### P4 — Effect-safety counterexamples

Falsify at least:

- consequential intake without `Idempotency-Key`;
- same key reused for different semantic request;
- timeout after possible external dispatch;
- stale pre-dispatch condition;
- accepted effect not yet converged;
- one confirmed + one ambiguous bulk member;
- ordinary access allowed while business/Governance rejects.

### P5 — Contract-authority proof

Deliberately introduce OpenAPI↔SDK and OpenAPI↔server drift in a negative fixture. Verification must turn red for the exact protected property. Merely having a generator/checker is not proof.

### P6 — Tenant/namespace proof

Falsify:

- cross-Organization path access;
- cross-Organization secondary-reference smuggling;
- two identical native IDs under different Installation/SourceInstance qualifiers collapsing in the SDK;
- provider ingress persisting durable acquisition state without explicit Organization.

### P7 — Structural inversion

Re-run the Method's Structural Inversion Test: if the current OpenAPI/SDK/routes had the opposite shape, these laws must still follow from D0–D4.

---

# 10. Independent-review adjudication

Fable verdict: **REVISE — direction confirmed; five D5-local corrections; zero D0–D4 reopen.**

GPT adjudication:

| Finding | Verdict | Resolution |
|---|---|---|
| **F-D5-1** freshness/provenance omitted from Q law | **AGREE** | B1-D4/D7, essential complexity, enforcement and proof now carry owner-controlled freshness/provenance when material |
| **F-D5-2** source-qualified external identity missing at wire boundary | **AGREE** | B1-D12 adds qualifier law + negative namespace fixture |
| **F-D5-3** ingress wording looked like forbidden Organization inference | **AGREE** | B1-D5 now separates fail-closed namespace identification from reading explicit MPC-owned Organization binding |
| **F-D5-4** idempotency mandatory class left open | **AGREE** | B1-D9 makes key mandatory/fail-closed by default for consequential intake, with explicit per-operation structural-idempotency exemption only |
| **F-D5-5** tenant tests missed secondary-reference smuggling | **AGREE** | B1-D5 + enforcement/proof add same-Organization resolution for every secondary reference |

Additional review adjudications:

- **Organization path scope:** ACCEPT. Mandatory path scope is the smaller, safer target than a mandatory Organization header; header alternative is rejected for the current product because it preserves an avoidable ambient/cache-key/consistency axis.
- **REST + owner-specific operations:** ACCEPT. No accepted Product 1.0 capability requires CRUD distortion or a generic command surface.
- **RFC 9457 problem split:** ACCEPT.
- **OpenAPI wire authority:** ACCEPT.
- **ADR-016:** historical after canonical consolidation, preserving only the two active lessons stated in B1-D13.
- **Hard cutover:** ACCEPT; no compatibility/versioning tax without a consumer.
- **No second Fable round:** no material contradiction remains after adjudication.

---

# 11. Explicit Unknowns / later D5 work

B1 does not decide:

- exact operation inventory/path nouns;
- concrete owner-specific request/response schemas;
- exact Permission mapping;
- pagination/filter/sort/cursor per real consumer;
- concrete provider-rich fields;
- concrete bulk endpoints;
- concrete concurrency-enabled operations/status mappings;
- concrete OpenAPI tooling/server generation/conformance stack;
- concrete technical/admin routes if later justified;
- total D5 batch count.

The next D5 work, if B1 is ratified, should derive the Product 1.0 operation/resource surface from D0–D4 owners and these laws rather than from legacy routes.

---

# 12. Reopen / stop triggers

Revisit only on material evidence:

1. marketplace commerce becomes only one vertical of a broader enterprise platform → targeted D0 then D1 review before API generalization;
2. a real external/public client introduces materially different compatibility/security/tenancy obligations;
3. required Product 1.0 operation cannot fit accepted owner semantics → targeted parent-stage review;
4. required client interaction cannot preserve D3 Q/C/P/outcome semantics → targeted D3 review;
5. real provider requirement makes D4 semantic/protocol separation impossible → targeted D4 review;
6. OpenAPI/derived client/server conformance cannot express a materially required contract → revisit wire technology, not duplicate manual authorities;
7. production consumers become entitled to an existing contract → revisit compatibility/versioning;
8. multiple real providers/consumers prove a repeated semantic enrichment deserving a smaller shared primitive;
9. real workflow proves bulk endpoint necessary → admit operation-local bulk with member-level semantics;
10. path Organization scope causes a concrete material defect and another explicit mechanism satisfies D2 more sustainably;
11. a real case proves source qualification cannot be represented without a new identity class → existing D2 identity reopen trigger.

Framework preference, current-code convenience and hypothetical future providers are not reopen evidence.

---

# 13. Converged candidate outcome

**Proposed outcome:** `RESTRUCTURE NOW` relative to the current API contract shape.

If explicitly ratified by the operator, canonical D5-B1 consolidation should establish:

- semantic/domain-oriented MPC Product API;
- Product API distinct from provider/business-system protocol ingress;
- REST/HTTP with owner-specific operations where CRUD would lie;
- Organization path scope + cross-Organization reference fail-closed rules;
- honest knowledge + freshness/provenance semantics;
- accepted/rejected/pending/ambiguous effect semantics;
- fail-closed idempotent consequential intake;
- source-qualified external identity + provider-rich/domain-owned enrichment;
- RFC 9457 for API problems, distinct from valid business outcomes;
- OpenAPI as singular machine-readable Product API wire authority;
- derived/conformant SDK/server with demonstrated drift failure;
- ADR-016 historical after consolidation;
- hard cutover with no speculative compatibility/versioning;
- no generic Mutation/Command/Provider/Workflow platform;
- no D0–D4 reopen.

**This file remains NON-AUTHORITATIVE until explicit operator ratification and subsequent canonical consolidation. Implementation remains blocked until D9.**
