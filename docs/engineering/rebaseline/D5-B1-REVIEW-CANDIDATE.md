# D5-B1 — Semantic API Model & Contract Laws — REVIEW CANDIDATE

> **Status:** NON-AUTHORITATIVE REVIEW CANDIDATE  
> **Stage:** D5 — API  
> **Batch:** B1 — Semantic API Model & Contract Laws  
> **Parent authority:** accepted D0 → D4 plus `ARCHITECTURE.md` and ADR registry  
> **Method:** DevelopmentConexus Engineering Method v1.0.0  
> **Purpose:** bounded package for independent adversarial review; this file does not change architecture authority, stage status, ADR status, or implementation permission.

---

## 1. Question this batch must answer

D5 must not begin by preserving, renaming, or mechanically generating the current API. The first material question is earlier:

> **What external API model and wire laws let MPC clients interact with the business authorities already accepted in D0–D4 without turning legacy modules, provider protocols, transport mechanics, or generic mutation abstractions into product semantics?**

The answer must preserve the accepted product boundary:

- MPC is the internal **Marketplace Operations Control Plane** for **Marketplace Operations + Commercial Intelligence**;
- Mercado Livre, Sankhya, payment APIs and future external systems are participants in that operating loop, not the product ontology;
- MPC is provider-independent and integration-independent, but the accepted Product 1.0 is not a generic enterprise-control-plane framework;
- a future decision to make marketplace commerce merely one vertical of a broader enterprise platform is a D0/D1 product-boundary reopen, not an API naming choice.

This batch deliberately decides **laws before endpoint inventory**. A later D5 batch may enumerate concrete Product 1.0 operations only after these laws survive challenge.

---

## 2. Scope / non-scope

### 2.1 In scope

This candidate decides, at D5 altitude:

1. the client-facing semantic API model;
2. the boundary between MPC Product API and provider/business-system protocol ingress;
3. how D1 owners and D3 Q/C/P meanings appear at the HTTP boundary;
4. Organization and Principal/access scoping at that boundary;
5. query knowledge-state laws;
6. consequential-action result, idempotency and precondition laws;
7. HTTP problem/error semantics versus valid business outcomes;
8. provider-rich evidence containment in client contracts;
9. one machine-readable wire-contract authority and SDK derivation;
10. hard-cutover / compatibility posture;
11. global bulk/partial-outcome restrictions;
12. proof obligations that a later implementation must satisfy.

### 2.2 Explicitly not in scope

This candidate does **not** decide:

- the complete endpoint/path/operation inventory;
- frontend screens, feature packages, cache/query hooks or UX — D6;
- worker, queue, outbox, retry scheduler, transaction, deployment or process topology — D7;
- end-to-end golden-flow choreography/proof — D8;
- product implementation — blocked until D9;
- a concrete OpenAPI generator, Go router/framework or TypeScript HTTP runtime;
- exact OIDC provider, token format, session topology or identity-provider deployment;
- universal pagination/filter/sort DSLs;
- universal event streaming/SSE/WebSocket surfaces;
- generic workflow, command, mutation, integration, evidence or provider frameworks.

---

## 3. Evidence classification

### 3.1 KNOWN from accepted authority

1. **Product boundary.** D0 defines MPC as a marketplace-operations control plane, not a marketplace dashboard, ERP replacement or generic integration hub.
2. **Semantic ownership.** D1 assigns each material Product 1.0 business responsibility to one semantic authority and explicitly rejects a generic `Mutation`/`Action` business domain.
3. **Identity/isolation.** D2 makes `Organization` the canonical tenant/isolation root, defines `Principal`, source-qualified external identity and domain-local identities for material Business Intents, and forbids inferring Organization from Installation/provider/source identity.
4. **Communication meaning.** D3 accepts semantic Q/C/E/P communication and distinguishes known/known-empty/unknown/unavailable query outcomes plus accepted/rejected/pending/ambiguous capability outcomes where applicable.
5. **External boundary.** D4 requires consumer-owned semantic ports, provider-local protocol/DTO knowledge, explicit namespace/capability/coverage semantics, and acceptance/ambiguity/convergence separation for external effects.
6. **Provider richness.** D4-B4 accepts **Semantic Core + Provider-Enriched Evidence**: no lowest-common-denominator flattening and no raw provider-payload ontology.
7. **No compatibility tax.** `ARCHITECTURE.md` and the Evidence Register state there are no production clients requiring current route/schema/package compatibility; hard cutover is allowed.
8. **Current API is evidence only.** Current OpenAPI/routes/SDK/module names are explicitly not inherited target authority.
9. **Current wire-authority defect exists.** OpenAPI, a hand-written SDK and handler/routing knowledge coexist. Legacy ADR-016 admits that same-commit OpenAPI+SDK discipline proves atomicity, not agreement, and leaves permanent human-transcription drift exposure.

### 3.2 INFERRED

1. The current API is a historical aggregation of several architectural styles rather than one coherent D0–D4 client contract.
2. A client-facing API that exposes `/ml/...`, `/sankhya/...`, generic `/mutations`, provider callback nouns or ERP-native operation vocabulary as product semantics would reintroduce authorities already rejected by D1/D4.
3. Generating an SDK directly from the current OpenAPI would automate the existing semantic defects rather than solve them.
4. Endpoint inventory before API laws would create local decisions likely to be redone once ownership, knowledge state, effect semantics and provider-boundary rules are applied consistently.

### 3.3 UNKNOWN

The following remain Unknown until later D5 work produces evidence:

- exact Product 1.0 operation set;
- exact path naming by semantic owner;
- exact read models needed by D6 consumers;
- exact Permission → operation catalog;
- exact filters/sorts/cursors per operation;
- which provider-rich fields deserve client exposure per named consumer;
- where real bulk workflows exist;
- which operations require optimistic-concurrency preconditions;
- exact authentication wire scheme above the accepted OIDC boundary;
- exact SDK package/generator technology;
- whether a named technical/admin API surface is needed beyond the Product API and provider ingress;
- number of remaining D5 batches.

### 3.4 DEFERRED

- D6: frontend consumption/topology;
- D7: runtime realization, generation framework, retries, queues, transactions and isolation mechanics;
- D8: real golden-flow proof;
- implementation: blocked until D9.

Unknown stays Unknown. Deferral does not authorize a plausible default.

---

## 4. Root cause

The current API problem is not primarily excessive route count or imperfect naming.

The structural defect class is:

> **The external client contract is not currently derived from the accepted semantic authorities, while product operations, provider/business-system protocol concerns and multiple manually maintained wire representations coexist without one explicit contract authority.**

This defect can repeatedly produce:

- provider vocabulary becoming product ontology;
- legacy module names becoming apparent business authorities;
- generic `Mutation` semantics stealing intent from D1 owners;
- business rejection being confused with HTTP/access failure;
- unknown/unavailable data being serialized as empty/default values;
- provider 2xx being interpreted as business completion/convergence;
- retry semantics ignoring ambiguous external acceptance;
- SDK/OpenAPI/handler drift;
- provider-specific enrichment either leaking wholesale or being discarded to fit a lowest common denominator.

Renaming routes or adding more parity checks does not eliminate this class.

---

## 5. Target invariant

> **Every externally invokable MPC Product API operation belongs to exactly one accepted semantic owner or accepted non-domain identity/access authority; its wire contract preserves Organization scope, current knowledge/effect semantics and ownership boundaries; provider/business-system protocol remains outside Product API semantics; and the Product API has one machine-readable wire authority from which client contracts are derived.**

Corollaries:

1. HTTP shape does not create business authority.
2. A generic technical mechanism never becomes a generic business owner.
3. Product API vocabulary is MPC business language.
4. Provider-native identity/evidence may remain visible only when semantically qualified and materially useful; raw provider DTO ontology does not cross the boundary.
5. A valid API response never makes a stronger knowledge/effect claim than the owning domain can justify.
6. Client retry convenience cannot weaken D3/D4 ambiguity or no-blind-retry safety.

---

## 6. Credible alternatives

### Alternative A — Clean the current OpenAPI and keep the hand-written SDK

**Local Maximum.** Improves the existing surface but preserves two manual wire descriptions and allows legacy API topology to constrain target semantics.

Rejected as the target.

### Alternative B — Generate the SDK directly from the current OpenAPI

Removes one drift mechanism but mechanically preserves the current mixed semantic/provider/legacy contract.

Useful only **after** the semantic contract is corrected.

Rejected as first-order solution.

### Alternative C — Generic platform API (`/resources`, `/commands`, `/mutations`, generic workflow/evidence/provider models)

Appears extensible but erases differentiated D1 meaning, contradicts the accepted no-generic-Mutation owner and creates framework authority without a second real product/vertical.

Rejected by YAGNI and authority rules.

### Alternative D — One independent external API per D1 domain

Preserves owners but prematurely converts semantic boundaries into deployment/public-contract fragmentation. D1 explicitly says twelve boundaries do not imply twelve services/processes.

Rejected as over-partitioning.

### Alternative E — Semantic MPC Product API + separate provider protocol ingress + one wire authority

One client-facing API uses MPC business semantics and owner-specific operations. Provider callbacks/OAuth/protocol endpoints remain integration boundaries rather than Product API. OpenAPI is the machine-readable HTTP contract authority; SDK/client types derive from it.

**Recommended Global Maximum.** It removes the current root cause without adding a generic platform or deciding D6/D7 topology.

### Alternative F — Introduce GraphQL/gateway/event-stream API as the new architectural center

No accepted consumer or failure class currently requires it. It changes API technology without solving semantic ownership, uncertainty, external-effect ambiguity or provider ontology leakage by itself.

Rejected for now. Reopen only with a real consumer/constraint.

---

# 7. Candidate decisions

## B1-D1 — Product API is semantic and domain-oriented, not provider- or integration-oriented

The client-facing API represents the accepted MPC marketplace-operating model.

It is:

- **provider-independent:** Mercado Livre does not define the API ontology;
- **business-system-independent:** Sankhya codes/workflows do not define the API ontology;
- **integration-independent:** adapters/protocols do not become client-facing business resources;
- **domain-oriented:** D1 owner meanings remain visible and distinct;
- **not enterprise-generic:** marketplace commerce remains the accepted product boundary unless D0/D1 are explicitly reopened.

A future broader enterprise control plane is legitimate only after product-boundary evidence exists. D5 must not prepare that future by replacing current semantic owners with generic entities/actions now.

### Forbidden implication

`provider-agnostic` does **not** mean `lowest-common-denominator`.

D4-B4 remains binding: a provider can expose richer evidence than another provider without forcing either suppression or fabricated equivalence.

---

## B1-D2 — External surface topology has two fundamentally different audiences

### A. MPC Semantic Product API

Consumed by MPC applications and future supported SDK clients.

It exposes:

- owner current meaning;
- read projections/compositions where legitimate;
- MPC-owned configuration/resources;
- owner capabilities/intents;
- Governance and Work interactions under their accepted meanings;
- semantically qualified provider-enriched evidence where required.

### B. Provider / business-system protocol ingress

Callbacks, webhook receivers, OAuth callbacks/handshakes and similar protocol endpoints belong to the D4 integration boundary.

They:

- may be externally reachable HTTP endpoints;
- are **not** Product API business operations;
- are not included in the normal Product SDK;
- accept provider protocol vocabulary only inside the adapter/protocol boundary;
- do not make notification payload fields domain truth;
- resolve Organization/Installation/SourceInstance through the accepted authenticated namespace binding rather than requiring a provider to know MPC business identifiers it does not own.

### No speculative third API

Health/readiness/administration/operations endpoints may exist later if a concrete runtime/operator consumer requires them. B1 does not create a generic `Internal API` or `Admin API` by taxonomy alone.

---

## B1-D3 — REST/HTTP resource semantics where truthful; explicit owner operations where CRUD would lie

The leading target is an HTTP/REST semantic API described by OpenAPI.

Use ordinary resource semantics where the client is genuinely reading/creating/updating an MPC-owned resource/configuration whose meaning fits those semantics.

Do **not** force consequential cross-system action into a fake direct-state CRUD contract.

Material consequential actions remain represented by their D1 owner and, where D2 already requires durable domain-local identity, by owner-specific Business Intents such as:

- Offering-owned Listing/Price Intent;
- Availability-owned Availability Intent;
- Materialization-owned Business Order Intent / Invoicing Intent;
- durable Fulfillment routing/dispatch intents where material.

There is no target generic Product API resource named `Mutation`, `Action`, `Command`, `Operation` or universal `BusinessIntent` merely to normalize those lifecycles.

Likewise, provider actual state remains external authority. A client request cannot truthfully pretend that a `PUT` completed the external world merely because MPC accepted an intent.

### Resource-oriented does not mean CRUD-only

Explicit owner operations are permitted when standard create/read/update/delete semantics would misrepresent business meaning. The operation name must remain in the owning domain's vocabulary and must not become a generic action escape hatch.

---

## B1-D4 — D3 semantic forms map to the Product API without changing authority

| D3 form | Product API meaning | Baseline wire law |
|---|---|---|
| **Q** | client asks an owner for current owner-owned meaning | response preserves known/empty/unknown/unavailable/partial semantics where material; cached/projection state cannot silently impersonate current owner truth |
| **C** | client asks an owner to accept/perform owner-owned work | request targets the owning semantic capability/resource; valid business outcome is returned as owner semantics, not disguised as provider transport status |
| **P** | client reads a composition of multiple authorities | explicitly read-only; projection/view is never consequential write authority; freshness/partiality remains honest |
| **E** | committed owner fact for independent consumer reaction | **not automatically translated into a public event-stream API**; external event/stream exposure requires a named consumer and later D5 decision |

D3 communication mechanism remains an internal semantic concern. D5 exposes product interaction, not the internal event topology by default.

---

## B1-D5 — Organization scope is explicit in the Product API

For Organization-owned Product API operations, the leading contract is:

```text
/organizations/{organization_id}/...
```

Reasons:

1. D2 requires Organization to be explicit and forbids inferring it from Installation/provider/source identity.
2. The tenant boundary becomes visible in operation identity, logs, authorization checks and SDK call shape.
3. A client cannot accidentally rely on process-global or token-only default tenant context.
4. One-Organization Product 1.0 proof does not hard-code a singleton into the contract.

The backend must verify that the authenticated Principal has ordinary access to the requested Organization; the path value is a scope claim, never self-authorizing.

### Provider ingress exception is not a contradiction

Provider callbacks normally do not know MPC `organization_id`. Protocol ingress resolves the correct Organization from the already authenticated/bound Marketplace Installation or SourceInstance namespace under D4. That resolution is adapter protocol behavior, not Product API tenant inference.

### Alternative retained for adversarial review

An explicit mandatory Organization header is technically credible. Fable should challenge whether it offers a materially better Global Maximum than the path-scoped form. **Token-only inferred Organization is not a credible target** under D2.

---

## B1-D6 — Principal/access and business authorization remain separate

Interactive Principal identity comes from the accepted OIDC boundary; exact token/session technology remains later realization.

At the Product API boundary:

- authentication proves/binds the Principal;
- ordinary identity/access state proves Membership/Permission for the operation;
- caller-supplied `principal_id`, role name or `approved=true` never substitutes for that proof.

D5 later maps stable Permissions to concrete API operations.

Critically:

- `401/403` represent authentication / ordinary-access failure;
- an action-owning domain's `prohibited`/`approval-required`/invalid business disposition is not ordinary access denial;
- Governance approval/rejection is not an API-role check;
- successful ordinary API access never implies consequential business authorization.

Business action disposition remains with the action owner. Controlled Action Governance retains its accepted grant/delegation/Authorization Decision semantics.

---

## B1-D7 — Query knowledge state is semantic output, not an HTTP failure shortcut

Where knowledge state is material, the wire representation must distinguish enough states to preserve the owner's actual claim, including as applicable:

- known value;
- known empty/absent;
- unknown / insufficient evidence;
- unavailable;
- partial / incomplete coverage.

B1 does **not** mandate a universal `Fact<T>` JSON envelope around every field or response.

The law is semantic:

> **No null, zero, `false`, empty object, empty list or HTTP success status may silently collapse a materially different knowledge state.**

Examples:

- `items: []` may represent known-empty only when the owning contract can prove the queried universe/scope is empty;
- source timeout/auth/rate-limit/outage cannot become an empty result;
- unsupported/not-applicable provider enrichment cannot become fabricated parity with another provider;
- a projection with incomplete components cannot present itself as complete current truth.

A successful HTTP exchange may therefore return a valid semantic `unknown`, `unavailable` or `partial` owner outcome. Transport/application availability and domain knowledge availability are different concerns.

Exact shared discriminators/schema reuse are admitted only where repeated semantics are genuinely identical; no universal Evidence/Knowledge graph is introduced.

---

## B1-D8 — Consequential capability outcomes remain distinct from HTTP/API problems

For a contract-valid owner capability request, the domain may legitimately return, where applicable:

- **accepted** — owner accepted the requested work/intent;
- **rejected** — owner definitively refused the requested work under current semantics/preconditions;
- **pending** — no final owner decision/effect yet;
- **ambiguous** — a possibly accepted external effect cannot yet be classified safely.

These are semantic outcomes, not raw provider transport statuses.

Binding distinctions:

```text
accepted != completed
completed != externally applied
externally applied != converged
```

and:

```text
rejected != access denied
pending != failed
ambiguous != failed
ambiguous != safe-to-retry
```

Provider timeout after possible dispatch cannot be converted into a failed/rejected business result merely for API convenience.

### HTTP-level preconditions versus owner/external preconditions

Two classes must not be collapsed:

1. **API concurrency/representation precondition** — e.g. a client supplies a stale MPC ETag/version for an update that explicitly requires current-version matching. This may be an HTTP conditional-request problem.
2. **Business/provider execution precondition** — e.g. the owner/provider detects current business/readiness/provider version state while evaluating/executing a valid intent. The resulting owner intent state remains domain semantics; provider `409`/status vocabulary does not leak through unchanged.

Fable should challenge exact status-code mapping, but the semantic separation is required.

---

## B1-D9 — Consequential API intake has explicit idempotency without replacing domain intent identity

For client requests capable of creating a new material Business Intent or initiating consequential work/effects, the leading contract requires a stable client retry token such as standard `Idempotency-Key` semantics.

The key is a **request-deduplication mechanism**, not:

- the canonical Business Intent identity;
- provider idempotency proof;
- authorization;
- permission to replay an ambiguous external write.

Baseline laws:

1. the same Organization + semantic operation + idempotency key + semantically equivalent request resolves to the same MPC intake result / durable intent where one was created;
2. reuse of a key for a materially different semantic request fails explicitly rather than creating another effect under the same token;
3. a retry after the first call becomes locally uncertain must resolve existing MPC intake before creating another intent;
4. if the owner already created a durable domain-local Intent, idempotent retries resolve to that Intent rather than inventing a generic Mutation record;
5. external-effect ambiguity is still reconciled according to D3/D4; the HTTP idempotency key cannot authorize blind redispatch.

Exact persistence, retention, locking and cleanup of idempotency records are D7 mechanisms. Exact scope may be sharpened during operation design; it must remain at least Organization + semantic operation so identical string keys in unrelated boundaries cannot collide by accident.

Fable should adversarially challenge whether the key should be mandatory for every consequential operation or only operation classes whose intake is retry-reachable.

---

## B1-D10 — Optimistic concurrency/preconditions are explicit only where stale client state can invalidate correctness

D5 does not impose version tokens on every resource.

Where a client update/capability depends materially on the current MPC-owned version, the API uses an opaque MPC-level concurrency/precondition token with standard HTTP conditional semantics where suitable.

Rules:

- stale pre-dispatch client state cannot silently overwrite newer authoritative MPC meaning;
- provider-native version tokens such as external `x-version` remain adapter-local unless a D1-owned semantic precondition genuinely requires qualified exposure;
- a provider stale-version response is translated into the owning domain's semantics and authoritative reread/redecision flow;
- projections/read models are not valid concurrency authorities for consequential writes.

Exact operations requiring preconditions are determined in later D5 operation mapping.

---

## B1-D11 — HTTP/API problems use one standard problem shape; valid business outcomes do not

The leading baseline for HTTP/API-level failures is **RFC 9457 Problem Details** (`application/problem+json`) with MPC-owned stable extensions only where real consumers require them.

API-level problem classes include proportionately:

- malformed/unparseable request;
- schema/contract validation failure;
- authentication failure;
- ordinary Permission/access denial;
- unsupported HTTP/operation contract;
- API concurrency/conditional-request failure;
- idempotency-key misuse/conflict;
- unexpected MPC server failure.

A stable MPC machine-readable `code` may specialize the problem `type`/extensions where programmatic handling needs more precision. Validation details may be structured.

### Provider diagnostic containment

Raw provider/business-system error DTOs, arbitrary text and PII do not become Product API problem truth.

Where a named support/operations consumer materially needs external-cause detail, the owning API may expose a **sanitized, source-qualified diagnostic** that:

- is clearly diagnostic/evidence rather than MPC business disposition;
- preserves enough provider/source correlation for support/reconciliation;
- does not mirror raw payloads or secrets/PII;
- never replaces the MPC semantic outcome.

### Non-errors

The following do not become Problem Details merely because they are undesirable outcomes:

- business `rejected`;
- `approval-required` / pending Governance decision;
- `ambiguous` possible external acceptance;
- valid query `unknown`/`unavailable`/`partial`;
- provider capability honestly `unsupported` / `external-required` where the operation contract admits those outcomes.

---

## B1-D12 — Provider-rich evidence is allowed only as domain-owned qualified enrichment

D4-B4 remains fully visible at the API boundary.

The Product API must not choose between two invalid extremes:

1. flatten every marketplace to a lowest-common-denominator schema; or
2. expose raw provider DTO/resource topology as the product contract.

Target rule:

> **A D1 owner may expose provider-specific enriched evidence when it serves a named Product 1.0 consumer/correctness property, but the enrichment remains source-qualified, bounded inside that owner contract and optional/unsupported when another provider lacks the same meaning.**

Consequences:

- provider richness may appear in Market Intelligence/economic/operational views where materially useful;
- provider-specific fields do not define top-level MPC resource identity or owner boundaries;
- absence of an equivalent field on another provider remains unsupported/not-applicable/unknown as appropriate;
- SDK types may use bounded discriminated enrichment unions without creating a universal `ProviderResource` graph;
- raw arbitrary provider payload passthrough is rejected.

Whether a concrete evidence field belongs in Product API is decided later against a named client/use case, not because D4 can acquire it.

---

## B1-D13 — OpenAPI is the single machine-readable wire authority for the Semantic Product API

Accepted D-stage artifacts and `ARCHITECTURE.md` remain **semantic architecture authority**. OpenAPI does not become authority over business meaning.

Within D5's HTTP boundary, however:

> **One OpenAPI document/set is the machine-readable authority for the Semantic Product API wire contract: paths, operations, parameters, headers, status codes and serialized schemas.**

Rules:

1. there is no independently authoritative hand-written SDK type model;
2. supported SDK/client types are mechanically derived from the OpenAPI contract or otherwise mechanically proven to conform to it;
3. server handlers must be mechanically validated against the same wire contract; exact codegen/router technology is later realization;
4. a change that makes implementation or generated client disagree with OpenAPI must fail verification;
5. a change in OpenAPI that conflicts with accepted D-stage semantics is an API-spec defect; machine-readable wire authority never outranks semantic architecture authority;
6. provider protocol-ingress schemas may have their own executable contracts when needed but do not contaminate the Product SDK or redefine Product API authority.

### ADR-016 disposition candidate

ADR-016 solved the historical no-generator condition with same-commit manual OpenAPI+SDK discipline. Its own Consequences state that this proves **atomicity, not agreement** and leaves permanent transcription risk.

If B1 is accepted and later consolidated, ADR-016's target meaning should become **superseded by D5**. Its useful historical lesson is retained as evidence: duplicated manual wire authorities require procedural synchronization and remain drift-prone.

D5 should converge the authorities rather than add more same-commit guards.

---

## B1-D14 — No compatibility/versioning machinery without a consumer

There are no production clients entitled to the current API.

Therefore target cutover may:

- delete routes;
- rename resources/operations;
- replace request/response schemas;
- delete generic Mutation/provider-facing Product API surfaces;
- replace the manual SDK contract;
- remove obsolete compatibility aliases.

D5 must not introduce parallel legacy versions, compatibility adapters, deprecation windows or dual-write/dual-read API contracts merely because mature public APIs often have them.

Whether the final base path contains a literal `/v1` prefix is a naming decision for later D5 operation topology; a version segment alone must not smuggle in a compatibility policy with no consumer.

If a real external/public client compatibility obligation appears later, that is material new evidence and reopens the versioning/compatibility decision.

---

## B1-D15 — Bulk is not a universal API primitive

No generic batch/mutation envelope is created.

A bulk operation is admitted only when a named Product 1.0 workflow/consumer materially needs bulk semantics that cannot be adequately composed by individual operations.

When admitted:

- intended target scope remains owner-defined;
- authorization scope remains distinct;
- member-level attempted/outcome state remains distinct where partial success/ambiguity can occur;
- one member's confirmed success does not authorize blind replay of the whole batch;
- one `success: false` boolean cannot collapse confirmed/rejected/ambiguous/not-executed members;
- provider bulk/import identifiers remain external reconciliation evidence, not MPC generic Batch identity by default.

Exact bulk endpoints are later D5 work.

---

# 8. Wire-semantics matrix

This table is a contract-law matrix, not an endpoint inventory.

| Interaction class | Semantic authority | Typical HTTP meaning | Valid semantic result | Idempotency / concurrency | Must never mean |
|---|---|---|---|---|---|
| Owner current query (**Q**) | D1 owner / accepted D2 substrate | request was processed | known / known-empty / unknown / unavailable / partial as applicable | safe read retry subject to freshness/coverage semantics | `[]`, `0`, `null` fabricated as knowledge |
| Read projection (**P**) | no new authority; composes owners | projection returned | complete/partial/freshness-qualified view | read-only | projection becomes write/concurrency authority |
| Create/update ordinary MPC-owned configuration | owning domain/substrate | resource/config accepted | created/updated representation | idempotent resource semantics and conditional update where materially required | provider actual state was changed merely because MPC config changed |
| Consequential owner capability (**C**) | action-owning domain | valid invocation reached owner | accepted / rejected / pending / ambiguous where applicable; often references domain-owned Intent | retry key when duplicate intake is reachable; domain/precondition semantics preserved | provider 2xx == converged; timeout == failed; generic Mutation owns intent |
| Governance interaction | Controlled Action Governance | valid authorization request/decision interaction | pending / approved-authorized / rejected/invalidated under Governance semantics | decision identity/history preserved; reapproval does not rewrite past | 403 ordinary access == Governance rejection |
| Operational Work interaction | Operational Work for work lifecycle; source owner for source truth | valid work interaction | work state plus source-domain resolution semantics as separate authorities | Work idempotency/dedupe stays Work-owned | closing Work mutates source truth |
| Provider protocol ingress | D4 adapter/protocol only | provider message/handshake received | acquisition pointer/protocol result; later owner semantics after translation/reread | provider-specific duplicate/auth semantics behind boundary | callback DTO becomes Product API/domain truth |

---

# 9. Naming and boundary examples — illustrative, not endpoint inventory

The following examples show the semantic direction only.

### Semantically plausible Product API vocabulary

```text
/organizations/{organization_id}/marketplace-installations/...
/organizations/{organization_id}/readiness/...
/organizations/{organization_id}/offerings/...
/organizations/{organization_id}/price-intents/...
/organizations/{organization_id}/availability/...
/organizations/{organization_id}/availability-intents/...
/organizations/{organization_id}/market-intelligence/...
/organizations/{organization_id}/economics/...
/organizations/{organization_id}/sales/...
/organizations/{organization_id}/business-order-intents/...
/organizations/{organization_id}/invoicing-intents/...
/organizations/{organization_id}/fulfillment/...
/organizations/{organization_id}/post-sale-resolutions/...
/organizations/{organization_id}/work/...
/organizations/{organization_id}/authorization-decisions/...
```

Not all of these must survive later operation mapping, and plural/path spelling is not decided here.

### Product API vocabulary that requires rejection or explicit proof

```text
/ml/items/...
/ml/catalog-listing/...
/sankhya/orders/...
/sankhya/invoices/...
/integrations/{provider}/...
/mutations/...
/commands/...
/provider-resources/...
```

A provider/business-system name can legitimately appear inside the **protocol ingress boundary**, or inside a **source-qualified enrichment discriminator/evidence field**, without becoming Product API ownership vocabulary.

---

# 10. External benchmark lessons — evidence, not authority

The target direction is consistent with useful patterns seen in mature platforms, but no external platform is copied wholesale.

- **Kubernetes/control-plane pattern:** useful for distinguishing desired intent, observed external state and reconciliation; rejected as a reason to turn all MPC concepts into generic declarative resources/controllers.
- **Crossplane:** useful evidence that a control plane can orchestrate external systems; its provider-resource mirroring is specifically not the MPC ontology because D4 requires consumer-owned semantics.
- **Stripe:** useful evidence for explicit idempotent HTTP intake and stable machine-readable client contracts; its idempotency model cannot erase MPC's ambiguous external-effect state.
- **commercetools:** useful evidence for resource versions/optimistic concurrency where stale state matters; not a reason to treat externally authoritative state as MPC-owned CRUD.
- **Google API design patterns:** useful for resource-oriented APIs plus explicit custom operations when standard CRUD does not express meaning; not authority for MPC operation naming.
- **Unified API/Common Model products such as Merge:** useful counterexample. Their core product value is cross-provider Common Models; D4 deliberately rejects a universal lowest-common-denominator provider model for MPC.

Fable should independently verify or replace any benchmark claim it relies upon. These references create no requirement.

---

# 11. Local maximum vs Global Maximum

### Local maxima rejected

- rename `/ml/*` to `/marketplace/*` while preserving provider resource meaning;
- keep `/mutations` and add more domain metadata;
- generate TypeScript from the current OpenAPI without redesigning semantics;
- add stronger same-commit OpenAPI/SDK/handler parity checks while retaining duplicate authorities;
- create one generic `CommandResult` or `CapabilityResponse<T>` envelope for every domain;
- force every query into a universal `Fact<T>` JSON wrapper;
- hide all provider-specific fields to make providers look uniform;
- expose provider DTOs under an `extensions` bag and call that provider-agnostic.

### Global Maximum candidate

> **A semantic MPC Product API whose operation model follows D1/D2/D3 meaning, whose integration protocol boundary follows D4, whose wire semantics preserve uncertainty/effect safety, and whose machine-readable contract is singular enough that SDK/runtime drift is mechanically detectable.**

This solution is smaller than a generic platform and more sustainable than cleaning the legacy surface.

---

# 12. Essential vs accidental complexity

## Essential complexity preserved

- Organization isolation;
- Principal attribution and ordinary access;
- D1 business authority;
- known/empty/unknown/unavailable/partial distinctions;
- domain-local consequential Intent identity;
- business disposition versus Governance authorization;
- idempotent client intake where duplicates are reachable;
- precondition/concurrency safety where stale state matters;
- accepted/rejected/pending/ambiguous effect semantics;
- external acceptance versus convergence;
- member-level partial/ambiguous outcomes for real bulk;
- provider-enriched evidence when materially useful;
- source/provenance and provider diagnostic containment.

## Accidental complexity removed/refused

- generic Mutation owner;
- provider/business-system paths as Product API ontology;
- duplicate manually authoritative OpenAPI + SDK contracts;
- compatibility with unused legacy routes;
- generic integration/provider/resource graph;
- universal command/capability response envelope;
- universal evidence/Fact JSON model;
- speculative multi-version API infrastructure;
- generic bulk framework;
- GraphQL/gateway/event streaming without a named need.

---

# 13. Enforcement candidates

D5 decides protected properties; D7/implementation later chooses exact mechanisms. The target must make these properties mechanically falsifiable.

1. **Operation-owner map:** every Product API operation must declare its accepted semantic owner or accepted non-domain identity/access authority plus Q/C/P interaction class.
2. **Boundary classification:** every externally reachable route must be classified Product API, provider protocol ingress, or a separately justified technical surface; unclassified routes fail architecture/conformance review.
3. **Provider-vocabulary fence:** provider-native nouns may not define Product API path/operation ownership unless explicitly classified as source-qualified enrichment/protocol ingress.
4. **OpenAPI authority:** SDK/client contract is reproducibly derived or mechanically conformance-checked from OpenAPI; no hand-written second wire authority.
5. **Server conformance:** handler request/response/status/header behavior is validated against the admitted OpenAPI contract.
6. **Knowledge-state negative fixtures:** unavailable/partial/unknown cannot pass as empty/default known values.
7. **Outcome negative fixtures:** ambiguous possible acceptance cannot pass as rejected/failed; accepted cannot pass as converged.
8. **Idempotency negative fixtures:** same retry token + different semantic request must fail; same request must not create another durable Intent/effect.
9. **Access/authorization negative fixtures:** ordinary Permission may allow invocation while domain/Governance still rejects or pends the action; business rejection cannot be encoded as ordinary 403.
10. **Tenant negative fixtures:** an authenticated Principal cannot access another Organization merely by changing path scope; final structural isolation proof belongs to D7/D8.
11. **Provider diagnostic redaction:** raw external PII/secrets/error payload cannot cross Product API accidentally.
12. **Bulk partiality fixtures:** confirmed + ambiguous members cannot become one safe-to-retry batch failure.

---

# 14. Proof strategy before implementation

B1 can be falsified before product implementation through architecture/contract artifacts and counterexamples.

### P1 — Complete operation ownership

For the eventual D5 operation inventory, prove every operation has exactly one:

```text
Product API operation
  -> semantic owner / accepted substrate authority
  -> Q | C | P
  -> ordinary Permission requirement
  -> knowledge/outcome class
  -> idempotency/precondition requirement if consequential
```

Any operation that cannot be classified without inventing a new owner is a D1/D2/D3 issue, not an API convenience exception.

### P2 — Protocol separation

Take current/provider examples such as marketplace webhook callbacks, provider Item operations and Sankhya-native order/invoice identifiers. Demonstrate that protocol nouns terminate at the D4 boundary and translate into accepted owner semantics rather than becoming Product API authorities.

### P3 — Knowledge counterexamples

At minimum falsify:

- search/source unavailable;
- incomplete enumeration;
- known empty;
- unsupported provider enrichment.

The schemas must make each materially different result distinguishable.

### P4 — Consequential-effect counterexamples

At minimum falsify:

- business rejection versus 403 access denial;
- stale API client precondition versus provider/business rejection;
- provider timeout after possible dispatch;
- accepted submission without convergence;
- external state diverging after prior accepted intent;
- retry after ambiguous result;
- one confirmed + one ambiguous bulk member.

### P5 — Contract authority/drift

Implementation-phase proof must demonstrate a deliberate OpenAPI contract change causes derived client/server conformance artifacts to change/fail mechanically. A green gate that never exercises drift detection is no proof.

### P6 — Structural inversion

Re-run D5 design assuming the current OpenAPI, SDK and handlers had the opposite shape or did not exist. B1 laws should remain unchanged except where current evidence reveals a real Product 1.0 consumer requirement.

---

# 15. Adversarial challenge package for Fable

Fable should reconstruct repository authority first, treat this file last/non-authoritatively, apply Method v1.0.0, and attack at least these questions:

1. **Product boundary:** Does this candidate preserve the accepted marketplace-commerce product while correctly avoiding Mercado Livre/Sankhya/integration specificity, or does it accidentally require a D0 reopen?
2. **Global Maximum:** Is semantic REST + explicit owner operations actually the smallest sustainable structure, or is another API model materially superior without adding YAGNI?
3. **Organization scope:** Is `/organizations/{organization_id}/...` superior to an explicit mandatory organization header? Find concrete correctness/future-cost counterexamples. Token-only inference is allowed only if Fable can reconcile it with D2's explicit-scope rule.
4. **Q/C/P mapping:** Does any accepted D1 capability fail to fit the proposed API interaction laws without distorting D3 meaning?
5. **Intent modeling:** Does owner-specific Intent exposure preserve authority, or does the candidate accidentally turn every action into an unnecessary durable resource?
6. **HTTP outcome split:** Challenge the rule that valid business rejected/pending/ambiguous outcomes are not ordinary HTTP problems. Identify where standard HTTP semantics should still carry precondition/conflict meaning without leaking provider state.
7. **Knowledge state:** Can the law be enforced without a universal response envelope? Find any operation class where contextual schemas would become structurally inconsistent or ambiguous.
8. **Idempotency:** Is a request-level `Idempotency-Key` the right seam? Challenge mandatory scope, equivalence definition, conflict behavior and relation to domain-local Intent identity.
9. **Concurrency:** Is MPC-level ETag/conditional semantics sufficient where stale client state matters while provider versions stay local? Identify any unavoidable provider-version exposure.
10. **Provider-rich evidence:** Does bounded domain-owned enrichment avoid both lowest-common-denominator flattening and provider DTO mirroring? Attack future second-provider cases.
11. **Error model:** Is RFC 9457 Problem Details plus owner semantic outcomes sufficient? Identify duplicate/unowned error semantics or missing operational diagnostic needs.
12. **Contract authority:** Should OpenAPI be the single machine-readable HTTP authority? Attack semantic-authority inversion, generated-client limitations and server conformance gaps. Propose a better authority structure only if it removes more defect classes with less total complexity.
13. **ADR-016:** Is superseding same-commit manual SDK discipline justified, and what historical invariant—if any—must survive after derivation/conformance exists?
14. **Surface topology:** Is separating Product API from provider protocol ingress sufficient, or is a third externally reachable surface already evidenced by Product 1.0/runtime constraints?
15. **Compatibility:** Does hard cutover remain correct for D5, and is any versioning seam worth preparing now without a compatibility consumer?
16. **Bulk:** Find a real Product 1.0 workflow that already requires bulk at B1 altitude or confirm operation-local admission is enough.
17. **YAGNI:** Identify any abstraction in this candidate that exists only because other abstractions exist.
18. **Future retrofit:** Find the hardest plausible second marketplace / second business-system change and test whether this API law set preserves the right seam without prebuilding a generic platform.
19. **Proof:** For every recommended correction, state how the property could be falsified before/after implementation.
20. **Reopen discipline:** Distinguish findings that correct D5 from findings that actually require targeted D0–D4 reopen. Reviewer preference does not create a reopen.

Fable should return `APPROVE`, `REVISE` or `REJECT` with material findings only, corrected invariants and reopen triggers. A second review round is justified only if a real material contradiction remains after GPT adjudication.

---

# 16. Reopen / stop triggers

B1 or its parent authority must be revisited only on material evidence such as:

1. **Product boundary change:** a ratified objective makes marketplace commerce only one vertical of a broader enterprise operations platform → targeted D0, then D1 coherence/reopen before D5 generalizes the API.
2. **New API audience:** a real external/public partner/client appears with materially different compatibility, security, tenancy or versioning obligations.
3. **Missing semantic owner:** a required Product 1.0 operation cannot be expressed under accepted D1/D2 authority without distortion → stop and targeted parent-stage reopen.
4. **D3 mismatch:** a required client interaction cannot preserve accepted Q/C/P/outcome semantics → targeted D3 review rather than hiding the mismatch in HTTP.
5. **Provider-contract conflict:** a real provider requirement makes semantic/protocol separation impossible for a required operation → targeted D4 review with evidence.
6. **Contract technology failure:** OpenAPI cannot express a materially required client contract or mechanically derived clients cannot preserve it → revisit wire technology/derivation, not default to duplicate manual authorities.
7. **Real compatibility obligation:** production consumers become entitled to an existing contract → revisit hard-cutover/versioning posture.
8. **Repeated provider enrichment pattern:** multiple real providers/consumers demonstrate the same shared semantic enrichment concept → evaluate a smaller shared semantic primitive; do not jump directly to a universal Provider graph.
9. **Real bulk consumer:** a Product 1.0 workflow proves member-level bulk is necessary → design that operation with explicit partial/ambiguous semantics.
10. **Organization topology evidence:** path scope creates a concrete material defect and another explicit scoping mechanism satisfies D2 more sustainably → revisit B1-D5.

Framework preference, current-code convenience and hypothetical future providers are not reopen evidence.

---

# 17. Candidate outcome

**Proposed outcome:** `RESTRUCTURE NOW` relative to the current API contract shape.

This does **not** authorize product implementation. It means only:

- D5 should design from the accepted semantic target rather than preserve legacy routes/modules;
- Product API and provider protocol ingress are distinct boundaries;
- semantic owner/resource/operation laws come before endpoint inventory;
- unknown/effect/authorization semantics must survive the wire;
- OpenAPI should become the singular machine-readable Product API wire authority;
- supported SDK contracts derive mechanically rather than remain a second manually authoritative model;
- no generic Mutation/Command/Provider platform is introduced;
- no D0–D4 reopen is currently required by this candidate.

**This file remains NON-AUTHORITATIVE until independent review, GPT adjudication and explicit operator ratification are completed, after which only canonical D5/router/ADR consolidation may create target authority.**
