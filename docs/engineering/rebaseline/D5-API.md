# D5 — API

> **Status:** OPEN / ACTIVE — D5-B1 ACCEPTED / CANONICAL; next coherent D5 work not yet accepted  
> **Program:** Architecture Rebaseline / Technical System Design  
> **Parent authorities:** `D0-PRODUCT-SYSTEM-DEFINITION.md`, `D1-DOMAINS-BOUNDARIES.md`, `D2-IDENTITY-TENANT-DATA-OWNERSHIP.md`, `D3-COMMUNICATION-EVENTS.md`, `D4-EXTERNAL-INTEGRATIONS.md`  
> **Method:** DevelopmentConexus Engineering Method v1.0.0  
> **Opened / B1 accepted:** 2026-08-18

## 1. Purpose and boundary

D5 defines the external API contract by which MPC clients interact with the business authorities already accepted in D0–D4.

D5 does **not** preserve the current OpenAPI, routes, SDK, controller/package shape or provider-facing nouns by inheritance. Those remain current-state evidence only.

D5 does **not** choose:

- frontend screen/component/package/query topology — **D6**;
- worker, queue, outbox, retry, transaction, deployment, RLS or process topology — **D7**;
- end-to-end golden-flow execution/proof choreography — **D8**;
- product implementation — blocked until **D9** is accepted.

D5 also does not silently generalize the product. MPC remains the accepted **Marketplace Operations Control Plane** for **Marketplace Operations + Commercial Intelligence**. Mercado Livre, Sankhya, payment APIs and future external systems participate in that operating loop but do not define MPC ontology. Making marketplace commerce merely one vertical of a broader enterprise operations platform requires a targeted D0 product-boundary reopen followed by D1 coherence/reopen as needed.

---

## 2. Imported parent invariants

D5 imports rather than re-decides:

1. **One semantic authority per material business meaning.** HTTP shape never creates or transfers authority.
2. **No generic Mutation/Action/Command business owner.** Material Business Intents remain domain-local under their D1 owners.
3. **Organization is the canonical isolation root.** Organization scope is explicit and is not inferred from Marketplace Installation, Selling Entity, external account, IdP organization, source key or process-global defaults.
4. **Cross-Organization references are denied by default.** An Organization-scoped request cannot smuggle another Organization's business reference through body/query.
5. **Principal ordinary access is distinct from business disposition and consequential authorization.** Permission to invoke an API capability never proves that a business action is permitted/approved/executable.
6. **External identity is source-qualified.** Marketplace Installation and SourceInstance qualify external namespaces; bare provider/native identifiers are never canonical MPC-global identities.
7. **Q/C/E/P semantics remain binding.** D5 may encode them over HTTP but may not weaken or merge their meanings.
8. **Known, known-empty, unknown and unavailable remain distinct where material.** Partial/coverage and freshness remain honest; source or projection timing never masquerades as stronger knowledge.
9. **Accepted != completed != externally applied != converged.** Pending and ambiguous outcomes remain explicit where reachable.
10. **No blind replay of ambiguous external effects.** Client retry ergonomics cannot weaken D3/D4 effect safety.
11. **Consumer owns meaning; adapter owns protocol.** Provider/business-system DTOs, status taxonomies, auth and native choreography remain behind D4 boundaries.
12. **Semantic Core + Provider-Enriched Evidence.** MPC does not flatten providers to a lowest common denominator and does not mirror provider payload topology into Product API ontology.
13. **No compatibility tax without a consumer.** No production consumer is entitled to the current API; hard cutover is allowed.
14. **Mechanism != Authority.** OpenAPI, generators, routers, validators and SDK tooling are contract/runtime mechanisms, not business authority.

---

# 3. D5-B1 — Semantic API Model & Contract Laws — ACCEPTED / CANONICAL

**Outcome:** `RESTRUCTURE NOW` relative to the current API contract shape.

The operator ratified the converged B1 package after independent Fable challenge and GPT adjudication. Five review findings were accepted as D5-local corrections: freshness/provenance on Q results, source-qualified external identity at the wire boundary, explicit Organization binding in provider ingress, fail-closed idempotency for consequential intake, and cross-Organization secondary-reference protection. No D0–D4 reopen is required.

## 3.1 Root cause

The structural API defect is not route count or naming quality.

> **The external client contract is not currently derived from the accepted semantic authorities, while product operations, provider/business-system protocol concerns and multiple manually maintained wire representations coexist without one explicit wire authority.**

That condition can repeatedly produce:

- provider vocabulary becoming product ontology;
- legacy module names becoming apparent business authorities;
- generic Mutation semantics stealing domain-owned intent;
- bare provider/native IDs collapsing namespaces;
- business rejection being confused with access or transport failure;
- unknown/unavailable/stale/partial data becoming plausible defaults;
- provider 2xx becoming apparent completion/convergence;
- unsafe duplicate intake or replay after ambiguity;
- SDK/OpenAPI/server drift;
- provider richness either leaking wholesale or being discarded for lowest-common-denominator uniformity.

Renaming current routes, generating a client from the current OpenAPI, or adding stronger same-commit parity checks does not remove this defect class.

## 3.2 Governing D5-B1 invariant

> **Every externally invokable MPC Product API operation belongs to exactly one accepted semantic owner or accepted non-domain identity/access authority; its contract preserves explicit Organization scope, source-qualified identity, honest knowledge/freshness/effect semantics and ownership boundaries; provider/business-system protocol remains outside Product API semantics; and the Product API has one machine-readable wire authority from which supported client contracts derive and against which server behavior conforms.**

Corollaries:

1. HTTP shape does not create business authority.
2. Generic technical machinery does not become a generic business owner.
3. Product API vocabulary is MPC business language.
4. Provider-native evidence may be exposed only when semantically qualified and materially useful; raw provider DTO ontology does not cross the boundary.
5. A valid response never makes a stronger knowledge, freshness, authorization or effect claim than the owning authority can justify.
6. Client retry convenience never weakens D3/D4 ambiguity or no-blind-retry safety.
7. No second manually authoritative Product API wire representation is admitted.

## 3.3 Global Maximum selected

Accepted target:

> **Semantic MPC Product API + separate provider/business-system protocol ingress + one machine-readable Product API wire authority.**

Rejected local maxima / overfit alternatives:

- clean current OpenAPI while retaining a hand-written SDK authority;
- generate a client directly from the current mixed OpenAPI;
- generic `/resources`, `/commands`, `/mutations`, workflow/evidence/provider platform;
- one externally independent API per D1 domain;
- GraphQL/gRPC/gateway/event-stream architecture without a real consumer/failure class;
- screen/BFF-shaped API that prematurely makes D6 topology an API authority.

The selected structure removes the known root cause without constructing a generic enterprise platform or deciding D6/D7 realization.

---

## 4. Product API versus protocol ingress

### 4.1 MPC Semantic Product API

The Product API is consumed by MPC applications and future supported SDK clients.

It exposes, only through accepted owners:

- owner current meaning;
- legitimate read projections/compositions;
- MPC-owned configuration/resources;
- owner-specific capabilities and durable domain-local Business Intents where already required by D2;
- Governance and Operational Work interactions under their accepted meanings;
- semantically qualified provider-enriched evidence where a named Product 1.0 consumer/correctness property requires it.

### 4.2 Provider / business-system protocol ingress

Provider callbacks, webhook receivers, OAuth callbacks/handshakes and similar protocol endpoints belong to the D4 integration boundary.

They:

- may be externally reachable HTTP endpoints;
- are not Product API business operations;
- do not enter the normal Product SDK;
- may accept provider protocol vocabulary only inside the protocol boundary;
- do not turn notification payload fields into MPC/domain truth;
- identify the bound Marketplace Installation/SourceInstance fail-closed from authenticated/authoritative protocol markers;
- **never compute Organization from provider data**;
- read Organization from that namespace identity's explicit MPC-owned binding;
- fail closed on mismatch or ambiguity before attribution;
- record Organization explicitly in durable acquisition/recovery state that can outlive the ingress execution context.

No generic third `Internal API` / `Admin API` surface is created merely for taxonomy. A technical surface may be added later only for a concrete runtime/operator consumer.

---

## 5. HTTP / REST semantic model

The accepted target is an HTTP/REST semantic API described by OpenAPI.

Use ordinary resource semantics where the client genuinely reads/creates/updates an MPC-owned resource/configuration whose meaning fits those semantics.

Do **not** force consequential cross-system behavior into fake direct-state CRUD.

Material consequential actions remain represented by their D1 owner and, where D2 already requires durable domain-local identity, by owner-specific Business Intents such as:

- Offering-owned Listing/Price Intents;
- Availability-owned Availability Intents;
- Materialization-owned Business Order / Invoicing Intents;
- durable Fulfillment routing/dispatch intents where materially required.

There is no target generic Product API resource named `Mutation`, `Action`, `Command`, `Operation` or universal `BusinessIntent` merely to normalize those lifecycles.

Explicit owner operations are allowed when standard CRUD semantics would lie. Their names remain in the owning domain vocabulary; they do not become a generic action escape hatch.

---

## 6. Q / C / P at the Product API boundary

| D3 form | Product API meaning | Binding wire law |
|---|---|---|
| **Q** | client asks owner for current owner-owned meaning | preserve known/known-empty/unknown/unavailable/partial where material; when freshness-for-use matters, expose/reference owner-controlled observation/acquisition/provenance time; projection/cache state cannot impersonate current owner truth |
| **C** | client asks owner to accept/perform owner-owned work | target the owning semantic capability/resource; valid business outcome remains owner semantics, not provider transport status |
| **P** | client reads a composition of multiple authorities | explicitly read-only; projection never becomes consequential write/concurrency authority; component freshness/partiality stays honest |
| **E** | committed owner fact for independent reaction | not automatically exposed as a public event-stream API; a real external stream consumer would require a later D5 decision |

D5 exposes product interaction. It does not export the internal event topology by default.

---

## 7. Organization and Principal/access boundary

### 7.1 Organization is path-scoped for Product API operations

For Organization-owned Product API operations, the target contract is:

```text
/organizations/{organization_id}/...
```

This is **decided**, not merely preferred.

Rationale:

- D2 requires explicit Organization scope;
- an unscoped Product route does not exist by default;
- scope is visible in request identity, logs, links, problem `instance` values and URL-keyed client/cache artifacts;
- path scope avoids a separate ambient/default Organization header axis;
- one-Organization Product 1.0 proof does not hard-code a singleton.

The path value is a scope claim, never self-authorizing. The authenticated Principal must have ordinary access to that Organization.

### 7.2 Secondary references stay inside the path Organization

Any Organization-owned resource identity supplied in body, query or nested request structure must resolve inside the path Organization.

A request scoped to Organization A that references Organization B state fails closed; it never silently resolves across tenant boundaries.

### 7.3 Principal ordinary access remains separate from business authority

At the Product API boundary:

- authentication binds/proves the Principal;
- current Membership/Permission proves ordinary access to the operation;
- caller-supplied Principal IDs, role names or `approved=true` do not substitute for those proofs;
- `401/403` represent authentication / ordinary-access failure;
- action-owner `prohibited`, `approval-required`, invalid business state or Governance rejection/pending are **not** ordinary access failures.

Successful API access never implies business-action disposition, Governance authorization or execution-time validity.

---

## 8. Knowledge, freshness and provenance wire law

Where knowledge state is material, the Product API preserves enough distinction to represent the owner's actual claim, including as applicable:

- known value;
- known empty/absent;
- unknown / insufficiently known;
- unavailable;
- partial / incomplete coverage.

No universal `Fact<T>` JSON envelope is required.

> **No null, zero, `false`, empty object, empty list or HTTP success status may silently collapse a materially different knowledge state.**

Freshness is orthogonal to those states.

Where freshness-for-use is material to a client decision:

- response exposes or references owner-controlled observation/acquisition/provenance time sufficient for the consumer to judge freshness;
- HTTP response time never substitutes for source/owner provenance;
- `known` never implies fresh enough for every use;
- `projection.updated_at` never impersonates component/source observation time.

Unsupported/not-applicable provider enrichment remains honestly unsupported/not-applicable/unknown as appropriate; another provider's richer field is never fabricated or suppressed merely for uniformity.

---

## 9. Source-qualified external identity on the wire

Externally authoritative identifiers remain references inside the correct namespace.

> **A provider/native identifier appears in the Product API only with its qualifying Marketplace Installation / SourceInstance identity explicit in the schema, or when that qualifier is unambiguous from the operation's declared scope. A bare external identifier is never a Product API correlation key.**

This applies to provider Listing/Order/Shipment/payment identifiers, business-system Product/document/native keys and equivalent external references.

Two equal native values under different Marketplace Installations or SourceInstances must remain client-distinguishable without relying on hidden server context.

No synthetic MPC mirror ID is created merely to hide source qualification where D2 already says the external resource itself remains externally authoritative.

---

## 10. Consequential capability outcomes

For a contract-valid owner capability request, valid owner outcomes may include where applicable:

- **accepted** — owner accepted/created/continued owner-owned work;
- **rejected** — owner definitively refused under its semantics/current preconditions;
- **pending** — no final owner decision/effect yet;
- **ambiguous** — possible external acceptance cannot yet be classified safely.

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

Provider timeout after possible dispatch cannot be converted to failed/rejected merely for HTTP convenience.

API conditional-request failure and business/provider execution preconditions remain distinct:

- stale MPC representation/concurrency token may be an HTTP conditional-request problem;
- business/readiness/provider precondition outcome belongs to the owning domain; raw provider status/409 vocabulary does not become Product API truth.

---

## 11. Consequential intake idempotency

For every Product API operation classified as **consequential** — capable of creating a durable Business Intent or initiating consequential work/effect — a stable client idempotency key is **mandatory and fail-closed by default**.

Missing key means explicit API problem and **no intake**.

A per-operation exemption is allowed only when the operation contract explicitly proves that duplicate intake is structurally unreachable or harmless through owner-anchor semantics alone. The exemption is declared and reviewable in the D5 operation inventory.

The key is a request-deduplication mechanism, not:

- Business Intent identity;
- provider idempotency proof;
- ordinary access or Governance authorization;
- permission to redispatch an ambiguous external effect.

Baseline law:

1. same Organization + semantic operation + key + semantically equivalent request resolves to the same MPC intake result / durable Intent where one exists;
2. same key + materially different semantic request fails explicitly;
3. retry after locally uncertain first response resolves existing intake before creating another intent;
4. a durable domain-local Intent remains the business object; no generic Mutation row is invented;
5. external ambiguity still reconciles through D3/D4; intake idempotency never authorizes blind external replay.

Exact persistence, retention, lock and cleanup mechanisms belong to D7.

---

## 12. Optimistic concurrency and preconditions

D5 does not impose versions on every resource.

Where stale client state can materially invalidate an MPC-owned update/capability, use an opaque MPC-level concurrency/precondition token with standard HTTP conditional semantics where suitable.

Rules:

- stale client state cannot silently overwrite newer authoritative MPC meaning;
- projections/read models are not consequential concurrency authority;
- provider-native tokens such as Mercado Livre `x-version` remain adapter-local unless a D1-owned semantic precondition genuinely requires source-qualified exposure;
- provider stale-version outcomes translate into owner semantics plus authoritative reread/redecision.

Concrete operations requiring concurrency tokens are later D5 work.

---

## 13. HTTP/API problems versus valid business outcomes

API-level failures use **RFC 9457 Problem Details** (`application/problem+json`) as the baseline problem shape, with stable MPC extensions only where a real programmatic consumer needs them.

API problem classes include proportionately:

- malformed/unparseable request;
- schema/contract validation failure;
- authentication failure;
- ordinary Permission/access denial;
- unsupported HTTP/operation contract;
- API concurrency/conditional-request failure;
- missing/misused/conflicting idempotency key;
- unexpected MPC server failure.

A stable machine-readable MPC `code` may specialize the problem type/extensions where needed.

The following remain valid semantic outcomes rather than Problem Details merely because they are undesirable:

- business `rejected`;
- approval-required / pending Governance decision;
- ambiguous possible external acceptance;
- valid query unknown/unavailable/partial;
- provider capability `unsupported` / `external-required` where the owner contract admits it.

### Provider diagnostic containment

Raw provider/business-system errors, arbitrary text, payloads, PII and secrets do not become Product API problem truth.

Where a named support/operations consumer materially needs cause detail, the owner contract may expose a sanitized, source-qualified diagnostic that:

- remains diagnostic/evidence rather than business disposition;
- preserves enough source correlation for support/reconciliation;
- excludes raw secrets/PII/payload mirroring;
- never replaces the MPC semantic outcome.

---

## 14. Provider-rich evidence without provider ontology

A D1 owner may expose provider-specific enriched evidence when it serves a named Product 1.0 consumer/correctness property.

The enrichment:

- remains source-qualified;
- stays bounded inside the owning semantic contract;
- may use bounded discriminated unions where appropriate;
- remains optional/unsupported/not-applicable when another provider lacks the same meaning;
- never defines top-level MPC resource ownership;
- never becomes raw provider payload passthrough;
- never forces another provider into fabricated equivalence.

Concrete evidence fields are later D5 operation/schema decisions against named consumers.

---

## 15. One machine-readable Product API wire authority

Accepted D-stage artifacts and `ARCHITECTURE.md` remain semantic architecture authority. OpenAPI never outranks them.

Within the Product API HTTP boundary:

> **One OpenAPI document/set is the single machine-readable authority for paths, operations, parameters, headers, status codes and serialized Product API schemas.**

Binding rules:

1. no independently authoritative hand-written SDK wire model;
2. supported SDK/client types are mechanically derived from OpenAPI or mechanically proven to conform to it;
3. server request/response/status/header behavior is mechanically validated against the same admitted OpenAPI contract;
4. implementation/client disagreement with OpenAPI fails verification;
5. OpenAPI conflict with accepted D-stage semantics is an API-spec defect, not authority transfer;
6. provider protocol-ingress schemas may have their own executable contracts but do not contaminate Product SDK semantics;
7. conformance controls count only when a negative fixture proves they actually fire.

Exact generator/router/server technology is deferred to D7/implementation planning after the architecture stages.

---

## 16. ADR-016 disposition — historical

Legacy ADR-016 used same-commit manual OpenAPI+SDK edits because no generator existed. Its own consequences acknowledge that this proves atomicity, not agreement, and leaves manual transcription drift reachable.

D5-B1 supersedes that target shape.

**ADR-016 becomes historical.**

Two durable lessons are rehomed here:

1. **No second manually authoritative Product API wire representation.** Converge authorities instead of synchronizing duplicates by procedure.
2. **Contract-conformance controls must be shown to fire.** A parity/generation/check artifact that never detects deliberate drift is not proof.

No compatibility/transition window is needed for the old manual SDK because implementation remains blocked until D9 and no production consumer is entitled to that legacy surface.

---

## 17. Hard cutover / no speculative versioning

Because there is no current production compatibility obligation, target D5 may:

- delete routes;
- rename resources/operations;
- replace request/response schemas;
- remove generic Mutation/provider-oriented Product API surfaces;
- replace the manual SDK contract;
- remove obsolete compatibility aliases.

D5 does not create parallel legacy versions, deprecation windows, compatibility adapters or dual API contracts merely because mature public APIs often have them.

A literal `/v1` path segment, if later chosen, is only spelling unless a real compatibility policy is explicitly justified.

A future real consumer entitled to an existing contract is material evidence that can reopen this decision.

---

## 18. Bulk is operation-local, not a platform primitive

No generic batch/mutation envelope exists.

A bulk endpoint is admitted only when a named Product 1.0 workflow/consumer proves bulk semantics are needed and cannot be adequately composed through individual operations.

When admitted:

- intended target scope remains action-owner meaning;
- authorized scope remains Governance context;
- attempted/outcome scope remains execution evidence;
- member-level confirmed/rejected/ambiguous/not-executed or equivalent states survive where material;
- one confirmed member never authorizes blind replay of the whole batch;
- provider import/bulk IDs remain external reconciliation evidence rather than a generic MPC Batch identity by default.

---

## 19. Enforcement / proof obligations

D5 owns the protected properties; later stages choose exact realization. The following must become mechanically falsifiable.

### 19.1 Operation-owner map

Every Product API operation must declare:

```text
operation
  -> accepted semantic owner / accepted D2 substrate authority
  -> Q | C | P
  -> ordinary Permission requirement
  -> knowledge/outcome class
  -> Organization scope
  -> source-identity qualification when applicable
  -> idempotency/concurrency requirement when consequential
```

An operation that cannot be classified without inventing a new owner returns to the implicated parent stage.

### 19.2 Boundary classification

Every externally reachable route is classified as:

- Product API;
- provider/business-system protocol ingress; or
- separately justified technical surface.

Unclassified reachable routes fail architecture/conformance review.

### 19.3 Negative fixtures / counterexamples

At minimum prove red on:

- unavailable/partial/unknown represented as plausible known empty/default;
- materially stale known value with no owner-controlled provenance;
- projection `updated_at` impersonating source observation time;
- two equal native IDs under different Installation/SourceInstance qualifiers collapsing in the SDK;
- provider ingress persisting durable acquisition state without explicit Organization;
- Product request scoped to Organization A carrying a secondary Organization B resource reference;
- consequential intake without required idempotency key;
- same key reused for a different semantic request;
- ambiguous possible external acceptance represented as failed/rejected;
- accepted represented as converged;
- ordinary access allowed while business/Governance rejects or pends;
- confirmed + ambiguous bulk members represented as one safe-to-retry failure;
- raw provider PII/secret/error payload crossing the Product API;
- deliberate OpenAPI↔SDK drift;
- deliberate OpenAPI↔server drift.

A green contract/generator/conformance artifact that did not execute the protected subject is no proof.

### 19.4 Structural inversion

The B1 laws must remain true if the current OpenAPI/routes/SDK implementation were opposite in every relevant respect. They derive from D0–D4 authority, not legacy structure.

---

## 20. Explicit Unknowns / later D5 work

B1 intentionally does not decide:

- exact Product 1.0 operation inventory and final path nouns;
- exact owner-specific request/response schemas;
- exact Permission → operation mapping;
- pagination/filter/sort/cursor requirements by real consumer;
- concrete provider-rich fields exposed to clients;
- concrete bulk endpoints;
- concrete operations requiring optimistic concurrency and exact status mappings;
- concrete OpenAPI generator/server-conformance technology;
- concrete technical/admin routes if later justified;
- total D5 batch count.

Unknown stays Unknown. Current OpenAPI/code may inform evidence but does not resolve these by inheritance.

---

## 21. Exact next D5 action

**Derive the Product 1.0 operation/resource surface from accepted D0–D4 owners plus D5-B1 laws, not from legacy routes.**

The next coherent D5 batch must classify the actual external Product API operations by semantic owner and Q/C/P meaning before final path/schema spelling.

For each candidate operation, establish proportionately:

- real Product 1.0 consumer/use;
- semantic owner;
- Q/C/P class;
- Organization scope;
- ordinary Permission requirement;
- canonical/source-qualified subject identity;
- knowledge/freshness/output semantics for reads;
- intent/outcome/idempotency/precondition semantics for consequential capabilities;
- projection/read-only status when composing owners;
- provider-enrichment need if any;
- pagination/filter/sort/bulk only when a real consumer requires them.

Do not design D6 screens or D7 runtime topology as a side effect of this mapping.

---

## 22. Reopen / stop triggers

Revisit B1 or its parent authority only on material evidence such as:

1. marketplace commerce becomes only one vertical of a broader enterprise platform → targeted D0 then D1 review before API generalization;
2. a real external/public client introduces materially different compatibility/security/tenancy obligations;
3. a required Product 1.0 operation cannot fit accepted owner semantics → targeted parent-stage review;
4. a required client interaction cannot preserve D3 Q/C/P/outcome semantics → targeted D3 review;
5. a real provider requirement makes D4 semantic/protocol separation impossible → targeted D4 review;
6. OpenAPI/derived-client/server-conformance cannot express a materially required contract → revisit wire technology, not duplicate manual authorities;
7. production consumers become entitled to an existing contract → revisit compatibility/versioning;
8. multiple real providers/consumers prove a repeated semantic enrichment deserving a smaller shared primitive;
9. a real workflow proves a bulk endpoint necessary → admit operation-local bulk with member semantics;
10. path Organization scope causes a concrete material defect and another explicit mechanism satisfies D2 more sustainably;
11. a real case proves source qualification cannot be represented without a new identity class → use the existing D2 identity reopen trigger.

Framework preference, current-code convenience and hypothetical future providers are not reopen evidence.
