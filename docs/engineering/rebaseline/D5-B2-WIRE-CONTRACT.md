# D5-B2 — Wire Contract / Resource-Path-Schema Grammar

> **Status:** OPEN / ACTIVE — W1 Resource / Path / HTTP Grammar ACCEPTED IN-STAGE; W2 Request/Response Schema Grammar next  
> **Parent B2:** `D5-B2-PRODUCT-OPERATION-SURFACE.md`  
> **Operation inventory:** `D5-B2-OPERATION-ADMISSION-MATRIX.md`  
> **Parent authorities:** accepted D0→D4 + D4-R1 + D5-B1 + Decision Reconciliation Baseline + ratified D5-B2 Whole-Matrix  
> **Method:** DevelopmentConexus Engineering Method v1.0.0  
> **Opened / W1 accepted:** 2026-08-18

## 1. Purpose and wire boundary

This artifact converts the ratified D5-B2 operation inventory into one coherent HTTP/OpenAPI wire contract without moving or duplicating accepted business authority.

The wire is derived from accepted semantics. It is not derived from legacy routes, current OpenAPI, controller/package layout, frontend screens or provider endpoint topology.

OpenAPI is the single machine-readable Product API wire authority once the contract is spelled. Accepted D-stage artifacts remain semantic architecture authority; an OpenAPI shape that conflicts with them is defective.

The Wire Contract does not choose D6 screen/BFF composition, D7 server/router/generator/blob/worker/transaction/deployment realization, D8 controlled-effect proof choreography or implementation.

---

# 2. W1 — Resource / Path / HTTP Grammar — ACCEPTED IN-STAGE

## 2.1 Governing invariant

> **Product API paths identify stable MPC scope/resources or explicitly source-qualified external resources; HTTP methods express honest resource-state semantics or explicit owner capabilities. URI hierarchy never defines business ownership, provider topology or workflow state.**

Corollaries:

1. URI may qualify identity; it does not manufacture identity.
2. Resource nesting means identity/lifecycle containment or external namespace qualification, never merely “comes later in the process”.
3. Standard HTTP methods are used where their resource semantics are honest.
4. Owner-specific capabilities are explicit when CRUD or `PATCH status` would lie.
5. No wire convenience may hide Organization scope, source qualification, ordinary Permission, business disposition, Governance or effect safety.

## 2.2 Product API base and versioning

Product paths are defined relative to the OpenAPI `servers` base URL.

No `/v1` path prefix is admitted for Product 1.0 baseline.

Reason:

- there is no entitled production client requiring compatibility with a prior Product API;
- D5-B1/ADR-035 allow hard cutover;
- adding a public versioning axis before a real compatibility consumer exists is accidental complexity.

A real external/public client with an incompatible-evolution contract is the D5 reopen trigger for API versioning strategy.

Exact host/deployment base path remains D7/runtime realization and does not become semantic path grammar.

## 2.3 Organization path law

All Organization-owned Product API operations use:

```text
/organizations/{organization_id}/...
```

The path Organization is an explicit scope claim, never self-authorization. The authenticated Principal still requires current Membership and operation Permission.

All Organization-owned secondary references in body/query/nested structures must resolve inside the path Organization unless a later accepted cross-Organization relationship explicitly says otherwise.

No ambient/default/current-tenant header or process-global Organization axis is admitted.

## 2.4 Platform-scoped self discovery

The only baseline platform-scoped Product Q is:

```http
GET /access-context
```

It is self-only:

- Principal comes exclusively from trusted authenticated token context;
- no `principal_id` parameter exists;
- it returns only current Organizations/Memberships/effective ordinary MPC Permissions for that authenticated Principal;
- it is not `/me`, `/session` or a generic Principal resource surface.

Cross-Principal access enumeration remains Organization-scoped access administration.

## 2.5 MPC-owned canonical resource grammar

An MPC-owned canonical resource uses the Organization root plus a resource collection and opaque MPC-owned identifier, proportionately:

```text
/organizations/{organization_id}/marketplace-installations/{marketplace_installation_id}
/organizations/{organization_id}/listing-intents/{listing_intent_id}
/organizations/{organization_id}/price-intents/{price_intent_id}
/organizations/{organization_id}/inventory-sources/{inventory_source_id}
/organizations/{organization_id}/fulfillment-nodes/{fulfillment_node_id}
/organizations/{organization_id}/post-sale-resolutions/{post_sale_resolution_id}
/organizations/{organization_id}/work/{work_id}
```

Exact remaining collection nouns are W2/later wire spelling against the ratified matrix; this rule freezes the identity grammar, not every noun in advance.

MPC-owned resources are not nested under another resource merely because they refer to it. For example, a ListingIntent may reference a Marketplace Installation without becoming an Installation-owned child identity.

## 2.6 Source-qualified external resource grammar

An externally authoritative resource has no synthetic MPC mirror identity merely to make the URL prettier.

Marketplace-native resources whose namespace is the Marketplace Installation use explicit qualification, for example:

```text
/organizations/{organization_id}/marketplace-installations/{marketplace_installation_id}/listings/{native_listing_key}
/organizations/{organization_id}/marketplace-installations/{marketplace_installation_id}/sales/{native_sale_key}
/organizations/{organization_id}/marketplace-installations/{marketplace_installation_id}/shipments/{native_shipment_key}
```

The nesting above is namespace qualification, not Portfolio ownership of Listing/Sale/Shipment meaning.

Business-system external references likewise remain SourceInstance-qualified in schema/path where the operation scope does not already make the SourceInstance unambiguous.

Two equal native keys under different Marketplace Installations/SourceInstances remain distinguishable without hidden server context.

## 2.7 Keyed Q/view without fake identity

An owner-owned current meaning or derived view does not receive a synthetic canonical ID merely so every GET can look like `/resources/{id}`.

`ProductChannelReadiness`, `CompetitivePosition`, `ExpectedEconomics` and similar keyed Q meanings may use semantically keyed query/path shapes in W2 without minting `ReadinessID`, `CompetitivePositionID` or `ExpectedEconomicsID` unless accepted identity semantics independently require one.

Resource-oriented does not mean database-row-oriented.

## 2.8 No domain/provider namespace as Product topology

Rejected as baseline Product API structure:

```text
/organizations/{org}/offering/...
/organizations/{org}/availability/...
/organizations/{org}/materialization/...

/mercado-livre/...
/sankhya/...
/providers/...
/integrations/...
```

D1 owner names remain semantic architecture, not mandatory public URI taxonomy. Provider/business-system names remain D4 protocol/enrichment vocabulary where materially needed, never the Product API root ontology.

## 2.9 Standard HTTP methods when honest

Use ordinary methods for real resource semantics:

- `GET` — retrieve the current Product representation/query result;
- `POST` to a collection — create a new server-identified MPC-owned resource/occurrence where the admitted operation is resource creation;
- `PATCH` — mutate part of an existing mutable MPC-owned resource only when partial-resource-update semantics are honest;
- `PUT` or `DELETE` are not prohibited, but are admitted only where later concrete semantics genuinely match full replacement/deletion rather than by REST symmetry.

A successful server-created canonical resource should use ordinary HTTP creation semantics proportionately, including `201 Created` and `Location` once the exact response/status grammar is frozen.

W1 does not yet choose JSON Patch, JSON Merge Patch or a typed update schema. W2 must choose request shapes based on domain meaning, especially where `null`, arrays or explicit knowledge states make generic merge-patch semantics unsafe.

## 2.10 Owner-specific capability grammar

When a ratified C operation is not honest resource CRUD, use an explicit capability bound to the owning resource:

```text
POST {resource-uri}:verb
```

Examples of the grammar, not a complete final path inventory:

```text
POST /organizations/{org}/listing-intents/{id}:submit
POST /organizations/{org}/listing-intents/{id}:discard
POST /organizations/{org}/economic-attributions/{id}:resolve
POST /organizations/{org}/work/{id}:hold
POST /organizations/{org}/work/{id}:resume
POST /organizations/{org}/marketplace-installations/{id}:deactivate
```

This convention means “owner capability on this resource”, not a child resource and not a generic Action/Command system.

Rejected:

```text
/actions/...
/commands/...
/operations/...
/{resource}/actions/{provider_action}
```

## 2.11 No `PATCH status` workflow escape hatch

A transition/capability with its own business meaning, safety tuple or owner decision is never encoded merely as a writable status enum.

Rejected examples:

```text
PATCH ListingIntent { status: submitted }
PATCH Work { status: closed }
PATCH BusinessOrderIntent { status: retry }
PATCH AuthorizationDecision { status: approved }
```

Mutable draft/configuration state may use PATCH where it truly represents resource state. Submission, resolution, physical evidence and similar capabilities remain explicit operations.

## 2.12 Nesting fence

Deep lifecycle nesting is prohibited when it merely mirrors process order.

Rejected structural implication:

```text
Sale
  /business-order-intents
    /fulfillment
      /invoicing
        /shipments
```

Sales, Materialization, Fulfillment, Post-Sale and Economics retain independent authorities/identities. References/correlation express lifecycle relationships; URI containment does not fabricate parent-child ownership.

## 2.13 Concurrency grammar — strong ETag + If-Match

For mutable MPC-owned representations/capabilities whose ratified safety tuple requires current-state protection:

1. the authoritative read supplies a strong opaque MPC `ETag` validator;
2. the consequential update/capability supplies `If-Match` with that validator;
3. validator content is opaque to clients and does not expose database sequence, timestamp, Postgres internals or provider-native version as Product semantics.

Conceptual examples:

```http
GET /organizations/O/listing-intents/LI
ETag: "opaque-mpc-validator"

PATCH /organizations/O/listing-intents/LI
If-Match: "opaque-mpc-validator"

POST /organizations/O/listing-intents/LI:submit
If-Match: "opaque-mpc-validator"
```

Projection/read-model validators never authorize writes to underlying owners.

## 2.14 Missing versus stale precondition

When an operation's accepted safety tuple requires `If-Match`:

- missing required precondition → **428 Precondition Required** + RFC 9457 Problem Details;
- supplied `If-Match` does not match current MPC representation → **412 Precondition Failed** + RFC 9457 Problem Details.

Do not overload `409 Conflict` for ordinary optimistic-concurrency mechanics. Business/provider preconditions and valid owner rejection remain separate semantic outcomes.

## 2.15 Idempotency is independent of concurrency

`Idempotency-Key` and `If-Match` solve different failure classes:

```text
Idempotency-Key → duplicate intake
If-Match        → stale client state / lost update
```

W1 preserves D5-B1's ratified `Idempotency-Key` contract. Its exact per-operation placement and Problem Details codes are finalized in later Wire Contract work from the safety matrix.

The MPC contract defines its semantics explicitly; it does not depend on the header having a published IETF RFC.

## 2.16 OpenAPI version defer

OpenAPI remains the single machine-readable Product API wire authority.

W1 does not select an exact OpenAPI minor version merely because one version is newest. The later tooling/machine-readable sub-batch selects the smallest OAS version that expresses the accepted schema/contract and has sufficient real generator/router/validator support.

Toolchain maturity is mechanism evidence, not permission to weaken semantics.

---

# 3. W1 proof / negative controls

Later executable contract proof must be able to falsify at least:

1. an Organization-owned route without explicit Organization scope;
2. a cross-Organization secondary reference accepted under another Organization path;
3. `/access-context` accepting a caller-supplied Principal ID;
4. a fake MPC mirror ID replacing required source qualification for Listing/Sale/Shipment;
5. equal native keys under two Installations becoming indistinguishable;
6. a domain/provider namespace becoming required Product ontology;
7. a capability implemented only through writable `status`;
8. a workflow relationship turned into URI ownership/nesting;
9. an operation requiring current-state protection accepting a write with no `If-Match`;
10. stale `If-Match` being accepted or mapped to an unrelated business rejection;
11. `Idempotency-Key` being treated as a concurrency validator or vice versa;
12. a P/projection validator being accepted to mutate an underlying owner;
13. a fake resource ID being minted solely to make a keyed Q look CRUD-like.

---

# 4. W1 method outcome

**Parent architecture:** `CURRENT STRUCTURE CONFIRMED`.

**Wire result:**

> **Use identity-oriented Organization-scoped resource paths, explicit source namespace qualification for external resources, honest standard HTTP methods and resource-bound owner capabilities. Reject domain-namespaced/provider-shaped paths, process-order nesting, fake CRUD/status transitions and synthetic IDs created only for URI aesthetics.**

No parent-stage reopen is required.

---

# 5. Exact next Wire Contract work

**W2 — Request / Response Schema Grammar + Knowledge / Outcome Semantics.**

W2 must derive the smallest reusable wire primitives and operation-specific shapes without creating universal business wrappers. It must decide proportionately:

1. exact external-reference schema grammar for Marketplace Installation / SourceInstance qualified native identities;
2. opaque MPC ID serialization and same-Organization reference rules;
3. exact `Money` wire representation and other exact domain values required by D2;
4. ListingIntent draft/update shapes, `FOLLOW_SOURCE | EXPLICIT_OVERRIDE`, listing-context media references and PriceIntent correlation without embedding price in ListingIntent;
5. PriceIntent target union (`existing Listing | pre-creation ListingIntent context`) and explicit supersession lineage;
6. knowledge-state representation for known / known-empty / unknown / unavailable / partial without a universal `Fact<T>` envelope;
7. freshness/provenance representation where material without one global Evidence schema;
8. owner capability outcome shapes preserving accepted / rejected / pending / ambiguous separately from applied/converged;
9. intent/tracking resource state without a generic `Operation`/workflow envelope;
10. request update grammar: typed PATCH/update bodies versus generic JSON Patch/Merge Patch where null/arrays/knowledge semantics matter;
11. RFC 9457 Problem Details base + smallest stable MPC problem-code extensions;
12. exact response-body/HTTP-status relationship for resource creation, reads, resource updates and owner-specific capabilities;
13. discriminated provider-enriched evidence only where named consumer/correctness need justifies it;
14. no raw provider DTO passthrough or universal provider/evidence/property bag;
15. schema-level constraints needed to make cross-Organization, bare external identity, invalid union and knowledge-state collapse mechanically rejectable.

After W2, later Wire Contract work must still cover collection pagination/filter/search, exact Permission→operation mapping, complete `Idempotency-Key`/precondition table, technical non-Product ingress classification and OpenAPI/tooling spelling.

Implementation remains blocked until D9.
