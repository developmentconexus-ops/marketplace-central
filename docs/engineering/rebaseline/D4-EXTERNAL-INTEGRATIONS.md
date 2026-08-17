# D4 — External Integrations

> **Status:** OPEN / ACTIVE — D4-B1 accepted/canonical; D4-B2 contract core canonical with Installation Evidence Gate OPEN  
> **Program:** Architecture Rebaseline / Technical System Design  
> **Parent authorities:** `D0-PRODUCT-SYSTEM-DEFINITION.md`, `D1-DOMAINS-BOUNDARIES.md`, `D2-IDENTITY-TENANT-DATA-OWNERSHIP.md`, `D3-COMMUNICATION-EVENTS.md`  
> **Method:** DevelopmentConexus Engineering Method v1.0.0  
> **Opened:** 2026-08-17  
> **B1 accepted:** 2026-08-17  
> **B2 contract core ratified:** 2026-08-17

## 1. Purpose and boundary

D4 defines the concrete external acquisition, translation, capability, requirement, authentication, coverage and effect-reconciliation contracts by which Mercado Livre, Sankhya and later accepted external systems participate in the MPC operating model already fixed by D0–D3.

D4 does **not** create a new Integration business domain and does not move business authority into adapters.

D4 decides, proportionately:

- concrete provider/business-system namespace binding under D2 identities;
- authentication/credential protocol semantics required to establish the right external namespace;
- authoritative read/reread surfaces and source-specific observation authority;
- point/enumeration/delta/notification coverage semantics;
- Integration Support and provider-effective capability/requirement evidence;
- concrete external-effect acceptance/reconciliation obligations;
- provider/ERP-specific identifiers, requirements and source evidence required by the D1 consumers;
- transport admissibility where external provider guidance materially constrains the target.

D4 does **not** choose:

- MPC HTTP/OpenAPI routes, wire errors or client contract — **D5**;
- frontend screens, forms or projection topology — **D6**;
- workers, schedulers, queues, retries, locks, token-refresh scheduling, cursor persistence, secret-manager technology, transaction/process/deployment topology — **D7**;
- end-to-end golden-flow assembly/proof — **D8**;
- product implementation, which remains blocked until D9 is accepted.

Current code, schemas, OpenAPI, tests, runtime and legacy ADR implementation shapes remain evidence only unless this authority explicitly rehomes their meaning.

---

## 2. Imported parent invariants

D4 imports rather than re-decides these accepted meanings:

1. **Consumer owns meaning; adapter owns protocol.** Business consumers own the semantic ports they need. External adapters own provider/driver protocol translation.
2. **Provider DTOs never become domain authority.** Wire DTOs/status vocabularies/auth headers/pagination tokens stay at the external boundary.
3. **Organization scope is explicit.** Installation, SourceInstance, provider account/resource identifiers and credentials never substitute for Organization scope.
4. **Marketplace Installation is MPC identity; seller/account identity remains external.**
5. **SourceInstance is a transport-independent MPC reference identity for one logical externally authoritative business-system/source namespace.** Credential rotation does not itself re-identify the source.
6. **External resource identities remain source-qualified.** Listing/Variation, Sale/Order, Shipment and native financial movements are not promoted to MPC-global identities merely for normalization.
7. **Notification/poll/callback payload is acquisition evidence, not automatically MPC domain truth or a D3 domain event.**
8. **Knowledge state remains honest.** Known value, known empty/absent, unknown and unavailable remain distinguishable; transport failure cannot become empty/zero/false.
9. **Partial observation is not closure.** Absence from a partial provider/source pull proves no stronger claim than the completed scope actually observed.
10. **Acceptance is not convergence.** `accepted != completed != externally applied != converged`; ambiguous possible acceptance is never blindly retried.
11. **Provider PII is minimized.** Raw provider content is not retained or propagated merely because it was received.
12. **Capability has three authorities.** D1 already distinguishes Integration Support, Provider Effective Capability/Requirement and Effective Business Capability; D4 owns only the first two.
13. **Mercado Livre is first.** D4 does not build a speculative universal provider framework or universal ERP model.

---

## 3. D4-B1 — External Contract Grounding — ACCEPTED

**Outcome:** `RESTRUCTURE NOW` for the former Sankhya direct-Oracle target assumption; otherwise `CURRENT STRUCTURE CONFIRMED` with bounded grounding additions. No D0/D1/D2/D3 reopen is required.

The operator explicitly accepted the converged B1 batch after an independent Fable challenge and GPT adjudication, then clarified the Sankhya transport decision: Direct Oracle existed historically because a sanctioned Sankhya connection path was not yet known/available to the project; now that the API Gateway path is established, Direct Oracle is not wanted in the target architecture. Reviewer findings are evidence; the authority is the operator-approved result recorded here.

### 3.1 Governing external-contract invariant

> **Every external fact or effect entering MPC is qualified by the correct Organization + source namespace, acquired through a contract whose authority and coverage are explicit enough to preserve honest knowledge state, translated through a consumer-owned semantic boundary, and—when it can cause an external effect—correlated to an authoritative reconciliation surface. Provider protocol never acquires MPC business authority.**

The target grounding shape is:

```text
accepted D1/D2 consumer meaning
        ↓ consumer-owned port
D4 provider/business-system adapter
  - source namespace binding
  - auth/protocol
  - operation-specific authority + coverage
  - provider capability/requirement evidence
  - authoritative reread/reconciliation
        ↓
external system

D7 later supplies runtime mechanics around this boundary.
```

### 3.2 External protocol / semantic-boundary fence

1. A business consumer owns the semantic port it needs.
2. The concrete external adapter owns wire protocol, auth headers, endpoint vocabulary, pagination tokens, transport errors and DTO mapping.
3. Provider DTOs/status vocabularies do not cross into business contexts for convenience.
4. Provider-local shared HTTP/auth/pagination machinery is allowed as **mechanism**, but owns no business meaning.
5. One adapter may implement several consumer-owned ports. That does not imply one global generic `Provider` business port.
6. D4 may describe Integration Support for a concrete operation. It does not create a universal Provider/Resource/Operation/Capability entity graph or registry without a proven present consumer/failure class.
7. Assembly/self-registration/factory/process wiring remains composition/D7 mechanism, not D4 business authority.

Legacy generic plugin/self-registration/catalog/factory shape therefore has no target authority by inheritance.

### 3.3 Mercado Livre Installation ↔ external seller binding

A Mercado Livre `Marketplace Installation` binds one Organization marketplace participation/configuration to one authoritative external seller/account namespace.

1. External seller/user identity is provider-owned; it is neither Installation ID nor Organization.
2. Initial authorization/re-authorization establishes authenticated seller identity from provider-authoritative OAuth/user evidence such as provider `user_id` and `/users/me` where applicable.
3. Credential/token refresh for the same seller does not create a new Installation.
4. Authorization for a different seller must never silently rebind an existing Installation. Mismatch fails closed and requires explicit portfolio/integration reconfiguration under the owning semantics.
5. Provider seller/user identifiers are treated as opaque external identifiers and must not be narrowed to Int32; current Mercado Livre documentation explicitly warns that new user IDs can exceed that range.
6. `site_id`, nickname and similar provider attributes are observations/configuration evidence, not Organization or Installation identity roots.
7. Downstream provider resource identity remains Installation + provider-native resource identity per D2.
8. **Acquisition-time attribution is also fail-closed.** Where an acquired provider resource or authoritative identity surface exposes a seller/namespace marker, a mismatch against the bound Installation prevents attribution/storage of that evidence under the binding and surfaces an explicit integration condition. Where the provider exposes no authoritative per-resource marker, the established authorization binding remains the control.

The last rule is an attribution invariant. It does not require calling `/users/me` on every read and does not choose the runtime enforcement mechanism.

### 3.4 Sankhya SourceInstance binding

1. `SourceInstance` is Organization-qualified and identifies one logical Sankhya business-system namespace/environment whose native keys are referenced.
2. Production and sandbox/test are distinct external environments and must not be silently collapsed merely because their entities/keys look alike.
3. Rotating `client_id`, `client_secret`, `X-Token` or other sanctioned API credentials does not by itself create a new SourceInstance when the same authoritative namespace remains bound.
4. Changing to a materially different authoritative namespace/environment requires explicit source rebinding/new SourceInstance treatment per D2.
5. Source identity is independent of concrete protocol mechanics. Credential/auth/API implementation changes do not re-identify the source.
6. B3 owns exact native-key/fact mapping and company/location/as-of qualification.
7. Where a Sankhya response or source identity surface exposes an authoritative namespace marker, mismatch against the SourceInstance binding fails closed for attribution. Where no such marker exists, configured/authorized source binding remains the control.

### 3.5 Sankhya transport — API Gateway is the target

#### Decision

**RESTRUCTURE NOW — transport only.**

For MPC↔Sankhya integration, the provider-sanctioned Sankhya **API Gateway is the target transport**.

**Direct Oracle/database access is not part of the target architecture and is not a fallback path.**

The historical direct-Oracle/godror path existed because the project previously did not have a known/usable sanctioned Sankhya API integration path. That historical constraint no longer justifies carrying database access forward now that the Gateway integration path is available and preferred by the operator.

This correction is deliberately narrow. It does **not** reopen or discard these transport-independent constraints:

- Sankhya/business-system state is external to MPC;
- MPC business consumers depend on MPC-owned semantic ports;
- Sankhya protocol stays behind adapter boundaries;
- PostgreSQL remains the store for MPC-owned canonical state;
- business contexts do not gain ad-hoc ERP transport knowledge;
- no legacy MetalShopping/Postgres source shortcut is resurrected.

Current official Sankhya guidance independently supports Gateway as the standard integration surface, but target selection does not depend on proving that Direct Oracle is contractually forbidden for this specific client. The operator has explicitly chosen not to carry Direct Oracle into the target.

#### B3 evidence gate

B3 must evaluate the concrete Sankhya fact/command surface through the sanctioned Gateway/API path and prove it is sufficient for the required Product 1.0 claims.

“Sufficient” includes:

- semantic correctness of the data/command;
- authoritative source identity and native-key qualification;
- coverage/completeness behavior;
- pagination/window semantics;
- rate limits and blocking risk;
- request/response size and timeout behavior where material;
- operational viability for the required acquisition/effect volume.

B1 does **not** choose polling/job cadence, scheduler topology or concurrency. Those remain D7.

If a materially required Product 1.0 fact/command cannot be obtained correctly and operationally through the sanctioned Gateway/API surface, B3 must **STOP / SPLIT PREREQUISITE** and present the gap for explicit operator/architecture adjudication.

> **Failure of the Gateway surface does not authorize a Direct Oracle fallback.**

The response to an API limitation is to find a sanctioned supported capability, alter the affected Product 1.0 claim explicitly if justified, or reopen the transport decision with new operator-approved evidence. It is never to introduce database access silently.

### 3.6 Authentication and credential lifecycle are protocol, not identity

#### Mercado Livre

1. Private provider access uses the current server-side OAuth contract.
2. Access/refresh credentials are adapter/runtime secrets, never domain business state or domain-event payload.
3. Token/user evidence verifies seller binding; credentials themselves are not seller identity.
4. Refresh-token rotation is provider credential lifecycle; an old refresh token is not assumed reusable after refresh when provider semantics invalidate it.
5. Token lifetime is consumed from current provider protocol/response rather than frozen as an MPC business constant.
6. Expired/revoked/invalid credentials make acquisition auth-invalid/unavailable, not source data absent.

#### Sankhya

1. The target Gateway authentication surface uses OAuth 2.0 `client_credentials` with the provider-required `X-Token` contract.
2. Production and sandbox environment credentials/tokens are not mixed.
3. Legacy appkey/login and direct-database connection behavior are historical evidence, not target auth/integration paths for new work.
4. Credential changes do not alter SourceInstance when the same source namespace remains authorized.
5. Authentication failure is availability/auth state, not business absence.

B1 does not choose credential schema/encryption, secret-manager technology, token-refresh locking/scheduling, callback route, retry/backoff or process placement.

### 3.7 Provider notification is a trigger; authoritative reread owns current observation

For Mercado Livre notification-capable resources:

1. Notification receipt is external acquisition evidence, not MPC domain truth.
2. Resource/topic/provider-user/application metadata may locate and qualify the observation.
3. The authoritative resource read occurs before an owning domain commits MPC meaning when current provider state matters.
4. Repeated notifications may cause repeated rereads without creating duplicate business truth.
5. Notification arrival order is not provider business order and cannot regress a newer authoritative observation merely because it arrived later.
6. Provider `missed_feeds` is a bounded recovery aid, not complete durable provider history.
7. Notification outage/gap does not create a completeness claim. Another authoritative observation/reconciliation path is required where completeness is material.
8. Non-authoritative callback fields stay acquisition/provenance metadata rather than canonical business status.

Webhook exposure/acknowledgement, receiver queueing, retries, reconciliation scheduling and cursor persistence remain D7.

### 3.8 Coverage/completeness is operation-scoped

D4 does not introduce one global sync-phase vocabulary or universal `Coverage` entity. Each consumer-owned acquisition contract exposes or references only the smallest operation-specific provenance/coverage needed for the consumer's correctness claim.

#### Point observation

- concerns one explicitly source-qualified resource/key;
- success proves only what the authoritative endpoint legitimately establishes for that resource at observation time;
- `not found` is known absent only when endpoint/identity semantics actually prove absence for that exact scope;
- auth/transport/rate-limit/timeout/parsing failure is unavailable, not absent.

#### Enumerated observation

- coverage is only for the source-defined scope actually traversed;
- all required pages/scan segments complete before a completed-enumeration claim;
- provider cursors are protocol mechanics, not MPC identities;
- early stop/page failure/depth cap makes the observation partial;
- completed enumeration does not invent stronger snapshot isolation than the provider/source establishes;
- absence outside the enumerated scope proves nothing.

This is material for Mercado Livre because ordinary item-search paging does not establish complete enumeration for large seller populations; current official behavior exposes scan/scroll semantics for larger result sets.

#### Delta/change observation

- proves only the source-defined change set under explicit source preconditions/window;
- does not by itself prove every current object was observed;
- source prerequisites for change tracking are part of the coverage claim;
- empty delta with unknown/disabled prerequisites cannot become “nothing changed”.

This is material for Sankhya because `modifiedSince` depends on the source's change-log behavior; a zero-record response is not proof of zero underlying changes when the required log coverage is not established.

#### Notification trigger

Carries no global completeness claim by itself.

Freshness remains consumer/use-sensitive under D3.

### 3.9 D3 knowledge/failure semantics survive the external boundary

1. Adapters may normalize external failures for consumers, but may not map unavailability to plausible business values.
2. Auth failure, rate limiting, timeout, provider/gateway outage, malformed response and incomplete traversal do not become known empty/zero/false.
3. `known empty/absent` requires affirmative source semantics for the exact queried scope.
4. Unsupported integration operation is explicit rather than masquerading as empty data.
5. If provider behavior establishes only uncertainty, uncertainty remains explicit.
6. Provider-native error DTO/text stays at the external boundary; consumer-facing meaning is consumer-owned.
7. PII/secrets are not retained or propagated merely because a payload/error contains them.
8. Provenance preserves enough source identity and material source/acquisition time for D2/D3 lineage/freshness without requiring universal raw-payload storage.

Exact error class hierarchy and MPC HTTP encoding remain later-stage concerns.

### 3.10 Integration Support != Provider Effective Capability != Effective Business Capability

D4 applies D1's accepted three-level fence:

1. **Integration Support / Descriptor — D4 technical meaning:** can this concrete adapter/protocol attempt/read/write this external operation class at all?
2. **Provider Effective Capability / Requirement — provider-authoritative evidence translated by D4:** for this installation/source/resource/mode/context, what does the external system currently allow/require and what provider prerequisite/artifact/state applies?
3. **Effective Business Capability — consuming D1 authority:** given provider evidence plus MPC readiness/policy/authorization/current business state, may/should MPC perform the action now?

D4 does not own level 3.

No adapter `capability=true` may bypass Readiness, Offering, Availability, Fulfillment, Governance or another D1 owner. Conversely, a business domain cannot fabricate provider support when D4 cannot establish it.

### 3.11 Admission gate for later external-effect contracts

B1 does not choose concrete writes. Every later B2/B3 external-effect contract must identify at least:

1. target Installation/SourceInstance-qualified resource;
2. consumer-owned semantic intent/correlation anchor;
3. material external preconditions/requirements;
4. what the provider/source response actually proves — rejection, accepted submission, synchronous effect, pending work or ambiguity;
5. when outcome may be ambiguous after possible acceptance;
6. authoritative reread/reconciliation surface;
7. member-level outcome semantics where multi-target work can partially succeed;
8. source/provider occurrence/result discriminator only where same-vs-distinct correctness needs it.

> **Transport success is not promoted to `converged` merely because it returned 2xx.**

Current legacy Mercado Livre write code that maps successful HTTP status directly to `applied` remains evidence to re-adjudicate in B2, not target precedent.

Retry/backoff/idempotency-store/attempt-table mechanics remain D7; D3 no-blind-retry/ambiguity semantics remain authority.

### 3.12 Legacy ADR disposition at B1

#### ADR-003 — integration sequencing

The real D4 prerequisite is rehomed: authenticated provider calls require valid credential lifecycle and provider identity binding. The old strict `OAuth -> fee sync -> frontend` implementation sequence is not target architecture. Fee/economics evidence belongs B4; frontend is D6; any D9 residue remains D9.

#### ADR-004 — integration plugin framework

**Superseded by rebaseline for D4 target structure.** Provider-specific protocol remains near provider adapters and business consumers own semantic ports. No generic self-registration/auth-factory/fee-sync/plugin catalog becomes target authority.

#### ADR-010 — polling-only Mercado Livre access

The mission-time `polling-only / no webhooks` rule is **superseded** as a D4 target constraint. Honest freshness/failed-refresh meaning remains carried by accepted D0/D3/D4 semantics. D4-B1 accepts notification→authoritative-reread plus explicit coverage; receiver/scheduler/poll cadence remains D7.

#### ADR-014 / ADR-020 / ADR-032

D4-B1 does not adjudicate their market/economics-specific residue. They remain evidence for **D4-B4**. On-demand/local-Docker/runtime flags/current `CollectorPort` shape have no authority by inheritance.

#### ADR-015

B1 rehomes source-qualified identity, honest coverage and reread principles. The final listing/variation/provider-contract disposition remains **D4-B2**. Legacy composite ID format, one read-model table/module, manual-refresh runtime and “absent from completed pull = closed” do not carry forward automatically.

#### ADR-006 / ADR-007

**Historical after D4-B1 for target architecture.** Their direct-Oracle/default-godror transport meaning is superseded. Direct Oracle is not an admitted target transport and godror/OCI is not a target runtime requirement.

Their durable transport-independent lessons — Sankhya remains external, consumers depend on MPC-owned ports/adapters, business code does not inherit driver/protocol details, and a wrong legacy source store is not revived for convenience — are rehomed in D0–D4/`ARCHITECTURE.md`.

### 3.13 Explicit deferrals and next-batch ownership

#### D4-B2 — Mercado Livre Operational Contract — OPEN BELOW

B2 owns the concrete Mercado Livre operational surface needed by Product 1.0, including:

- listing/item/variation identities and authoritative point reads;
- seller population enumeration and completeness limits;
- catalog/product relations only where the first flow needs them;
- stock/availability read evidence and controlled writes;
- price/listing controlled writes;
- provider sale/order acquisition and authoritative reread;
- fulfillment modes, responsibilities, prerequisites, artifacts and provider-authoritative deadlines;
- concrete write acceptance/outcome/ambiguity/reconciliation semantics;
- composite behavior only if the selected first flow materially requires it;
- **current Product↔channel correspondence identifier evidence** required by D2/Readiness, including the real Mercado Livre `SELLER_SKU` / custom-field / attribute / variation-level semantics where applicable;
- **the current provider identifier-evidence surface available to Product & Channel Readiness for unattended corroboration and pre-dispatch correspondence validation.** Readiness remains authority for whether that evidence is sufficient.

#### D4-B3 — Sankhya Business-System Contract

B3 owns:

- the sanctioned API Gateway target surface under §3.5;
- exact authoritative Product/native key mapping;
- inventory company/location/as-of facts;
- cost/tax evidence and required source qualifiers;
- Business Order Intent materialization;
- Invoicing Intent materialization;
- native result keys and authoritative reread/reconciliation;
- explicit unsupported facts/commands;
- operational viability of the Gateway/API surface for the required correctness/coverage/volume.

If the Gateway/API surface proves insufficient for a required Product 1.0 claim, B3 stops and escalates the gap; **Direct Oracle is not a B3 fallback option** under current authority.

B3 does not choose scheduling cadence/process topology.

#### D4-B4 — Market / Economics / Settlement Contract

B4 owns official market/competitor acquisition where justified, fee evidence, catalog-offer/price-to-win evidence where material, provider financial movements, settlement/adjustments/refunds and their provenance/completeness/economic correlation.

D5–D8 remain untouched.

### 3.14 Proof strategy / strongest counterexamples

B1 is accepted against these falsification cases:

1. Same Mercado Livre seller, new tokens -> Installation remains stable.
2. Existing Installation authorized as a different seller -> fail closed; no silent rebind.
3. Acquired ML resource exposes seller B while bound Installation is seller A -> fail closed for attribution.
4. Provider user ID exceeds Int32 -> external reference remains representable.
5. Sankhya production/sandbox bindings are mixed -> no environment collapse and no fabricated absence.
6. Sankhya credentials rotate for the same authoritative namespace -> SourceInstance remains stable.
7. Duplicate ML notification -> repeated reread remains safe; no duplicate business truth.
8. Notification loss exceeds bounded recovery history -> no false complete-history claim.
9. Seller population exceeds ordinary paging limits -> ordinary successful page traversal cannot be called complete if the provider requires scan/scroll semantics.
10. One enumeration page/segment fails -> observation is partial; unseen resources cannot be closed by absence.
11. Sankhya `modifiedSince` returns zero while change-log coverage is unproven -> zero does not become “no changes”.
12. A provider endpoint exposes a referential quantity -> field naming alone cannot promote it to exact inventory authority.
13. Adapter implements a write but the installation/resource is not provider-eligible -> Integration Support does not become provider/business capability.
14. Authentication expires -> external fact becomes unavailable/stale as appropriate, never zero/empty/false.
15. Endpoint 404 lacks deletion/absence semantics for the exact contract -> no manufactured known absence.
16. Provider receives a write and the connection drops -> ambiguity/reconciliation required; blind retry forbidden.
17. Provider returns 2xx for accepted asynchronous work -> no convergence claim unless the source contract proves it.
18. Multi-target effect partially succeeds -> member outcomes remain representable.
19. Legacy Oracle path is technically easier/faster -> it remains historical and does not enter the target.
20. B3-required Sankhya fact is unavailable or operationally non-viable through Gateway -> STOP / explicit architecture decision; no database fallback.
21. A second marketplace arrives -> consumer-owned ports remain sufficient without a universal provider entity graph unless real repetition proves otherwise.

### 3.15 Reopen / stop triggers

Revisit only the implicated decision when material evidence shows:

1. Mercado Livre authorization cannot establish a stable seller namespace compatible with D2 -> targeted D2 review.
2. A required external resource identity cannot fit Installation/SourceInstance + native-key model -> targeted D2 review.
3. A concrete interaction needs a new semantic business dependency absent from D1 -> targeted D1 review.
4. Provider effect semantics cannot fit D3 accepted/rejected/pending/ambiguous + reread/reconciliation -> targeted D3 review.
5. Provider evidence makes the accepted Product 1.0 operating loop materially impossible/different -> targeted D0 review.
6. B3 proves a required Sankhya fact/command cannot be obtained correctly and operationally through the sanctioned Gateway/API and no sanctioned equivalent exists -> **STOP / SPLIT PREREQUISITE** and return to the operator/architecture decision loop.
7. Any proposal to reintroduce Direct Oracle/database access -> requires an explicit operator-requested reopen of §3.5 with new material evidence; it is never inferred from technical convenience or API limitation.
8. A concrete second provider creates repeated technical failure unsolved by consumer-owned ports + provider-local mechanism -> consider only the smallest proven shared mechanism then.

Framework preference, current-code convenience and speculative future providers are not reopen evidence.

---

## 4. D4-B2 — Mercado Livre Operational Contract — OPEN / CANONICAL CONTRACT CORE

**Outcome:** `CURRENT STRUCTURE CONFIRMED` for D0–D4-B1 with a bounded concrete Mercado Livre contract. No D0/D1/D2/D3/D4-B1 reopen is currently required.

The operator ratified the converged B2 contract core after a pre-review D0–D3/D4-B1 coherence pass, current-provider/reference review, independent Fable challenge and GPT adjudication. The review produced three consolidation findings: canceled-order search coverage must be explicit; stock writability/effect scope must include site/seller/UP context; and the Installation Evidence Gate is necessary but must be described as a lane-selection/completion gate rather than as uncertainty about the mode-conditional contract rules themselves.

B2 is **not closed**. The target contract core below is canonical; the current Metal Nobre Mercado Livre Installation's actually supported Product 1.0 lane set remains Unknown until §4.8 is satisfied. Under the current router, B3/B4 remain unopened while this gate is open.

### 4.1 Governing Mercado Livre invariant

> **For every Mercado Livre operation required by Product 1.0, the adapter maps a consumer-owned MPC query/intent to the currently applicable Installation-qualified provider resource and operating context, preserves provider authority, coverage, preconditions and effect scope, and returns only the semantic evidence needed by the owning D1 domain. Mercado Livre resource topology never becomes MPC business ontology, and provider context never silently widens business authority or intended scope.**

B2 therefore preserves provider-specific differences only where they change correctness or effective capability. It does not create a generic `Provider`, `OperatingMode`, `UserProduct`, `Warehouse`, `Claim` or `Return` MPC business entity merely to normalize Mercado Livre.

### 4.2 Installation & channel surface

#### Installation posture / reputation / restrictions

1. Current seller identity remains the B1-bound external seller namespace.
2. Provider-authoritative seller reputation, restrictions and moderation may feed Marketplace Portfolio/Offering attention and Provider Effective Capability where materially relevant.
3. Authentication/application availability remains separate from provider business posture.
4. D4 does not own reputation policy, complaint management, buyer messaging or reputation optimization; those remain outside Product 1.0 unless explicitly reopened.

#### Item / User Product / Family / Catalog

1. `item_id` is the provider Listing/offer-condition identity under Marketplace Installation.
2. A legacy provider Variation is an external child identity only where the actual listing uses the legacy variation model.
3. `user_product_id` is a provider-native seller physical-product reference; it is not MPC Product and gains no MPC canonical entity merely for normalization.
4. `family_id`/`family_name` are provider grouping concepts; `catalog_product_id` remains Mercado Livre catalog ontology.
5. Legacy and User Product listings may coexist. The adapter determines the actual provider model from current authoritative seller/resource evidence rather than a static installation-wide assumption.
6. Provider topology may be referenced by Installation + provider resource kind/native key where material; B2 does not create a universal provider-resource graph.
7. Provider migration or Item↔UP relationship change never rewrites MPC Product identity or silently rewrites Readiness authority.

#### Listing creation / observation / lifecycle

1. Marketplace Offering Operations owns Listing Intent and desired listing meaning.
2. Before dispatch, D4 establishes the applicable provider publication model and mandatory provider requirements for the current seller/category/resource context.
3. A User Product seller is not forced through a legacy `variations[]` creation model merely because old code supports it.
4. Unknown/unsupported required provider category/domain/attribute/catalog conditions remain explicit; endpoint existence does not prove provider readiness.
5. Provider-created Item/UP/family/catalog relations and current provider state are authoritative observations after creation.
6. Provider creation/2xx is not Listing convergence until Offering reconciles intended meaning with authoritative current provider state.
7. Point lifecycle/moderation evidence wins over unsupported absence inference; listing closure/inactivity requires source semantics that affirm that meaning for the source-qualified resource.

#### Listing population coverage

1. Listing enumeration is scoped to the authenticated Installation seller namespace.
2. Provider search/scan/scroll caps and traversal rules are part of the coverage claim.
3. A completed enumeration claim requires the provider-defined scope to be fully traversed; early stop/failure remains partial.
4. Absence from a partial/unknown-coverage enumeration never closes/deletes a listing.
5. Cursor/scroll persistence is D7 mechanism, not business identity.

#### Product↔channel provider identifier evidence

B2 supplies current provider evidence to Product & Channel Readiness; Readiness owns correspondence sufficiency.

Where present/material, the adapter may expose/reference:

- seller SKU / seller custom field;
- GTIN/EAN and relevant provider attributes;
- Item ↔ legacy Variation;
- Item ↔ User Product;
- User Product / Family evidence;
- Catalog Product relation;
- category/domain/attribute evidence needed to understand correspondence constraints.

No one field becomes MPC identity law. `SELLER_SKU == CODPROD` is not resurrected. Contradictory evidence is preserved rather than normalized away, one match is not enough for unattended correspondence, and identity-bearing outbound fields remain subject to D2 pre-dispatch correspondence revalidation.

### 4.3 Offering & availability effects

#### Price observation vs price write

1. Current price observation uses the current authoritative Mercado Livre price surface for the claim being made; deprecated/convenient Item price fields do not remain read authority by inertia.
2. Offering owns Price Intent and business validity; D4 owns provider protocol/effective capability only.
3. Provider read authority and current write mechanism may be different surfaces and are recorded separately.
4. A successful write is not price convergence until authoritative current price is reread and matches the provider-effective intended result.

#### Price Automation

1. Adapter implementation of price mutation proves only Integration Support.
2. Where current provider evidence says Price Automation blocks direct API price editing, direct MPC price write is Provider Effectively unavailable/rejected for that listing/context.
3. B2 does not automatically disable or reconfigure Mercado Livre Price Automation merely to make MPC's write path succeed.
4. A mixed provider response that applies other fields while ignoring/rejecting price is not price convergence even when transport returns 2xx; authoritative price reread controls the conclusion.
5. Offering decides the business disposition of an unavailable/external-required price action.

#### User Product shared-effect scope

Current Mercado Livre behavior can asynchronously replicate shared User Product fields across multiple Items associated with the same UP. This materially includes `available_quantity` where that provider path applies.

Therefore:

1. Before dispatching a provider-shared field change, D4 establishes the current provider resource/effect scope sufficiently to detect a materially wider blast radius.
2. A single-Listing or single-availability-target intent is never silently widened into a multi-Item intent.
3. If provider effect scope exceeds intended/authorized scope, execution fails closed or returns to the owner for explicit scope redecision before dispatch; exact encoding is D5/D7.
4. Offering/Availability own intended scope; Governance owns authorized scope when applicable; D4 supplies provider-effective scope evidence.
5. A changed Item↔UP relation or other material scope drift is revalidated at execution time under D0/D3 rules.
6. Any availability write using an Item path on a UP-model listing inherits this same shared-effect-scope gate.

#### Provider stock topology / writability

`available_quantity` is not one universal marketplace stock contract.

B2 recognizes provider contexts such as seller-writable simple/legacy Item availability, User Product/location stock, seller-managed locations and provider-managed `meli_facility`/Full stock only as external protocol contexts.

Rules:

1. Availability Control owns Inventory Source/Scope, allocation policy, Sellable Availability and Availability Intent/convergence. Provider stock locations remain external references/evidence.
2. `store_id`, `network_node_id`, `selling_address`, `seller_warehouse` or `meli_facility` never becomes Inventory Source or Fulfillment Node by convenience.
3. **Seller-managed does not automatically mean API-writable.** Provider writability requires the concrete resource typology **plus current site applicability, seller configuration/tags and per-resource/listing state** to establish an enabled write surface.
4. Current official behavior includes site/configuration-specific stock surfaces; no MLA/MLC capability is projected onto MLB merely because the typology name matches.
5. `seller_warehouse`/multi-origin write capability is treated as conditional on the current provider configuration/enablement required by that surface.
6. Provider-managed `meli_facility`/Full quantity is observable but is not presented as seller/MPC-writable availability merely because the adapter can read it.
7. Unsupported/provider-managed paths are explicit, never successful no-ops.
8. The actual Item-vs-UP/location write surface for the real MLB Installation remains probe-gated under §4.8.

#### Stock version conflict

Where the selected provider stock surface uses `x-version` or equivalent optimistic precondition:

1. the provider version is protocol precondition, not business identity/global sequence;
2. stale-version HTTP 409 is a **rejected stale provider precondition**, not ambiguous possible acceptance;
3. the next correctness step is authoritative stock/version reread followed by fresh owner validity/redecision;
4. the same stale request is not blindly replayed;
5. any new request still requires current Availability intent/policy/authorization/correspondence validity.

#### Availability proof constraint

D0 requires automatic normal-path synchronization for sufficiently-known, authorized MPC-controlled availability.

Therefore the Product 1.0 proof set must contain at least one real provider-writable availability lane if MPC claims that controlled convergence capability complete. Provider-managed Full stock can prove observation of provider-owned availability, but a **Full-only** proof set cannot by itself prove MPC-controlled Sellable Availability convergence.

This is a proof-selection rule, not a requirement to implement every Mercado Livre stock mode.

### 4.4 Sale & fulfillment operational contract

#### Order acquisition / point authority

1. Order notifications/callbacks are triggers/pointers; Marketplace Sales commits MPC sale meaning only after the required authoritative Order reread/evidence.
2. Duplicate notification may repeat point reread without duplicating Sale meaning.
3. Provider arrival order is not MPC business-order authority.
4. Seller/order namespace mismatch fails closed under B1 attribution rules.
5. `GET /orders/{id}` is the point authority for the source-qualified Order details needed by B2 when the ID is known.

#### Order enumeration / history coverage

1. Order search is an enumeration surface with provider-defined filters/retention, not permanent Sales historical authority.
2. Current official behavior documents retained/searchable Orders up to a bounded period and, materially, seller-scoped search **filters canceled Orders**.
3. A completed seller Order enumeration therefore proves only its provider-defined scope; it must name the cancellation exclusion and cannot be labeled complete for a Sales population whose claimed universe includes canceled Orders.
4. Cancellation acquisition/recovery cannot rely on seller Order-search enumeration alone.
5. Notification→point-reread remains valid when the source-qualified Order ID is known, but notification/missed-feed behavior is not promoted to complete durable history by assumption.
6. If a correctness claim requires complete recovery/discovery of canceled Orders or another population excluded by provider search, B2/D8 must prove an authoritative discovery/recovery combination for the explicit universe. Until proven, that portion of coverage remains partial/unknown.
7. Provider-search retention cannot serve as permanent Sales historical authority; accepted D2/D3 lineage must retain/recover material historical meaning required beyond the provider window.
8. Historical repository evidence that an undocumented/unclear Order scan mode once worked remains evidence only until current official/real-dependency proof supports the exact coverage claim.

#### Order and Shipment remain separate

1. Provider Order remains Sale/checkout evidence; provider Shipment remains its own source-qualified external identity.
2. `shipping.id` on Order is correlation/reference evidence, not permission to treat embedded Order shipping fields as current Shipment authority indefinitely.
3. Fulfillment/shipment current state is obtained from the provider Shipment surface when materially needed.
4. Provider Pack remains provider-native when relevant; no universal MPC Pack entity is introduced.
5. Marketplace Sales does not become Fulfillment authority because the Order payload contains shipping fields.

#### Cancellation / fraud evidence

1. Provider `cancel_detail`/equivalent authoritative evidence remains source fact; B2 translates only semantic evidence Sales/Post-Sale require.
2. Cancellation group/code/requester are not guessed from unrelated status fields.
3. Provider fraud/stop-shipment conditions may block progression even after earlier payment evidence; they remain provider-effective conditions, not silently ignored.
4. D4 does not own Post-Sale Resolution semantics.

#### Effective fulfillment context

1. Marketplace brand alone does not determine fulfillment responsibility.
2. D4 derives/translates current provider-effective logistics/fulfillment context from the concrete Order/Shipment/provider evidence needed by the selected lane.
3. Provider strings/codes remain inside the adapter; consumers receive semantic requirement/capability evidence.
4. Fulfillment Lifecycle owns provider-requirement closure for flows MPC claims to control.
5. Provider-operated paths remain explicit `provider-operated`/`external-required` where appropriate; MPC never claims physical work owned by the provider.
6. B2 creates no universal global `OperatingMode` entity.

#### Fiscal prerequisites / labels / dispatch readiness

Where the selected lane exposes provider prerequisites before dispatch:

1. D4 translates current provider requirement/artifact/state.
2. Fulfillment owns whether the provider requirement is open/closed for its path.
3. Business-System Materialization remains owner of Business Order/Invoicing Intent and authoritative Sankhya materialization under B3.
4. Provider-required fiscal/document submission is an external effect only where the selected lane actually requires/supports it.
5. Submission/2xx is not closure; authoritative Shipment/provider reread establishes current provider state.
6. Label/document availability is mode/state sensitive; endpoint support does not prove current label capability.
7. Provider-managed modes with no seller label action remain provider-managed rather than integration failure.
8. A required provider operation unsupported by the accepted integration is explicit `unsupported`/`external-required`; MPC does not claim full control while silently depending on routine portal work.

#### External deadline / SLA

1. The applicable current provider Shipment/SLA deadline surface is external-authoritative evidence where available.
2. Provider deadline and MPC Internal Operational Target remain distinct authorities.
3. A consequential decision close to dispatch uses sufficiently current deadline evidence; stale cached provider time is not treated as immutable.
4. Deadline absence/unavailability is not “no deadline” unless source semantics prove non-applicability.
5. D7 owns timers/schedulers/attention mechanics.

#### First-flow proof selection

B2 does not implement every Mercado Livre fulfillment mode merely because documentation lists them.

1. The selected proof lane set is chosen from fresh real Installation evidence.
2. A provider-writable availability lane is required if D0 Availability Control controlled convergence is claimed complete.
3. The selected Sale/Shipment lane must make fulfillment responsibility explicit and close every provider prerequisite MPC claims to control.
4. A Full/provider-operated lane may be supported honestly but cannot stand in for internally operated Fulfillment Node separation/conference/packing/dispatch work if D0 completion claims that internal path was proven.
5. If the first proof requires the internally operated Fulfillment Node normal path, a real seller-operated physical fulfillment lane must be selected or the proof gap must be surfaced.
6. Unsupported modes remain explicit and can be added later without changing D1 authority.

### 4.5 Essential post-sale provider contract

1. Provider Claim/Case/Return/reverse Shipment remain Installation-qualified provider-native references/evidence; B2 does not create duplicate MPC aliases merely for normalization.
2. Post-Sale Resolution remains the MPC canonical resolution/correlation/closure owner.
3. Return/consequence scope remains representable at material item/line/quantity scope; one provider Return does not globally close a Sale.
4. Reverse Shipment remains distinct provider shipment evidence and does not collapse into the original Shipment.
5. Provider `available_actions` or equivalent current state is **Provider Effective Capability**, never automatic MPC permission. Post-Sale/consequence owners retain business decision authority.
6. Current provider capability/state is reread before consequential post-sale action where material.
7. Product 1.0 support remains bounded to essential cancellation/return/refund consequences; general complaint management, buyer Q&A/chat, reputation workflows and general CRM/SAC remain outside scope.
8. B2 may preserve/correlate provider-native financial movement references needed downstream, but authoritative Payment/Refund/Fee/Adjustment/Settlement/Payout contracts and realized-economic interpretation remain **D4-B4 + Commercial Economics**.
9. Physical Return terminal state does not fabricate financial closure; financial movement does not fabricate physical/Post-Sale closure.

### 4.6 External-effect reconciliation rule specialized for Mercado Livre

For every B2 write admitted later by a selected lane:

- target is Installation-qualified and correspondence-valid;
- provider-effective prerequisites/capability are current enough for the action;
- intended and authorized scopes are not silently widened by provider topology;
- provider response is classified only as strongly as it proves;
- possible asynchronous/partial/ignored-field behavior is preserved;
- authoritative point/resource reread establishes external current result/convergence where required;
- ambiguous possible acceptance is never blindly retried;
- definitive stale-precondition rejection returns to reread/redecision rather than ambiguity handling.

A provider 2xx is never sufficient by itself for Price, Listing, Availability, Fulfillment or Post-Sale convergence when current authoritative state can differ from the submitted intent.

### 4.7 YAGNI / provider-overfit fence

B2 prepares seams justified by current provider evolution but does not implement every Mercado Livre capability in advance.

Explicitly rejected as target structure by B2:

- universal Integration/Provider/ProviderResource business graph;
- canonical MPC `UserProduct`, `Family`, `ProviderWarehouse`, `OperatingMode`, `Claim` or `Return` aliases merely for normalization;
- one universal marketplace stock/write surface;
- universal raw Mercado Livre payload archive;
- support for every documented logistics/fulfillment mode before a real selected consumer/proof lane requires it;
- automatic deactivation of provider Price Automation merely to obtain MPC write control;
- campaign authoring, buyer messaging, general SAC/reputation optimization;
- D7 scheduler/queue/retry/cursor/process topology inside D4.

Known provider modes that are not selected for Product 1.0 proof may remain explicit `unsupported`/`external-required` without weakening source semantics. A later real mode/provider may extend concrete adapters through the existing consumer-owned seams.

### 4.8 Installation Evidence Gate — OPEN

#### Status and blocking scope

The mode-conditional contract rules in §§4.2–4.7 are canonical now. The current Metal Nobre Mercado Livre Installation's actual seller/resource modes and Product 1.0 supported lane set are **not** known from public documentation or the historical 2026-08-01 probe.

Unknown remains Unknown.

This gate blocks:

- **D4-B2 closure / whole-batch ACCEPTED status under the current router;**
- declaration of the actual supported Mercado Livre lane/mode set for the real Installation;
- any claim that D0 Availability Control, internal Fulfillment Node execution, Price write or another lane-specific capability has been proven on the real Installation;
- D8 lane selection/proof claims that depend on those facts.

The gate does **not** mean that the already-ratified conditional contracts above are uncertain. Probe outcomes select which of those contracts apply to the real Installation.

Under the current router, B3/B4 remain unopened while B2 is open. If the operator later wants B3/B4 sequencing to proceed while this subgate remains open, that requires an explicit router/sequencing decision; it is not inferred from reviewer convenience.

#### Required probe properties

Execute a read-only real-dependency probe with:

- no Mercado Livre write;
- no secret/token value recorded;
- no buyer PII retained;
- minimal classification/evidence only;
- current provider authoritative point/resource reads;
- stated universe/sample so absence is not overclaimed.

#### Minimum facts to establish

1. seller tags relevant to publication/stock model, including `user_product_seller`, `warehouse_management`, `multiwarehouse` where applicable;
2. selected real listing topology: legacy vs User Product, Item↔UP relations, Catalog relation and real Variation/composite presence where material;
3. current provider price model for candidate listings and whether Price Automation blocks the intended Price-write proof lane;
4. current stock typologies/locations/ownership **and the concrete write surface actually enabled for this seller/site/listing**, including whether `/items.available_quantity` still applies to candidate listings and what UP/shared effect scope it has;
5. recent selected real Orders/Shipments and their actual fulfillment/logistics/fiscal/label/SLA contexts;
6. current Installation/listing moderation/restriction evidence materially relevant to the proof set;
7. Claim/Return presence if available for later D8 proof; sample absence does not imply the provider capability does not exist.

#### Gate outcomes

- If a provider-writable availability lane exists, choose the smallest real lane needed to prove D0 Availability Control and leave other modes unsupported/external-required until needed.
- If no provider-writable availability lane exists for the accepted Product 1.0 scope, surface a targeted D0/product-proof conflict before calling Availability Control complete; Full/provider-managed observation alone is insufficient proof of MPC-controlled convergence.
- If Price Automation blocks every candidate Price Intent lane, direct price write remains Provider Effectively unavailable; do not disable provider automation by architecture assumption.
- If the selected first fulfillment path is provider-operated, support it honestly; do not claim internally operated Fulfillment Node execution from provider-owned work.
- If the operator requires the internal Fulfillment Node path in the first proof, select/prove a real seller-operated lane or surface the gap.
- If a materially required provider feature cannot be reached on the real Installation, mark explicit unsupported/external-required or reopen only the actually implicated upstream decision; never fabricate support.

### 4.9 Legacy/current-state disposition while B2 remains open

#### ADR-015

B2's canonical core supersedes these legacy target assumptions now:

- one canonical read-only `listings` module/table;
- composite MPC listing IDs derived from provider keys;
- manual refresh as target runtime;
- one full pull as permanent synchronization architecture;
- “absent from completed pull = closed” as lifecycle law.

Durable lessons are rehomed in D1/D2/D4: provider state remains external, Listing Intent/convergence belongs Offering, external identity is source-qualified, time/coverage remains honest and acquisition enters through consumer-owned ports/adapters.

However ADR-015 remains in the active legacy registry until the B2 Installation Evidence Gate closes and B2 is accepted as a whole; its final transition to historical is a B2-closure action rather than being inferred from this partial status.

#### ADR-022 / ADR-028

Their legacy identity formulas remain superseded. B2 now rehomes the concrete D4 obligation to supply current Mercado Livre identifier/correspondence evidence and enforce pre-dispatch consistency; Readiness keeps unattended-corroboration policy/human-decision authority.

#### Provider DTO/schema-drift evidence

Historical adapter DTO omissions prove provider translation drift is a real failure class, but B2 does not require universal unsanitized raw-payload persistence. Proportionate sanitized fixtures, contract tests, selective PII-minimized evidence or later runtime observability may prove mappings; exact implementation belongs later stages.

### 4.10 B2 proof strategy / strongest counterexamples

B2 must remain coherent against at least these cases:

1. seller is legacy, User Product, or mixed -> adapter follows actual resource model without duplicate MPC Product identity;
2. one UP links multiple Items -> shared effect scope is known before a narrow write;
3. Item↔UP relation changes after authorization -> current blast radius is revalidated;
4. Price Automation is active or price is ignored inside a mixed 2xx response -> no false Price convergence;
5. seller-managed stock exists but the current MLB site/configuration does not enable that location write surface -> no false Integration Support/Provider Capability;
6. Item availability write on a UP listing would propagate to sibling Items -> Availability intent is not silently widened;
7. `meli_facility` Full stock is observed -> no seller/MPC write capability is fabricated;
8. stale `x-version` returns 409 -> rejected precondition + reread, not ambiguous retry;
9. provider stock location resembles an internal warehouse code -> no Inventory Source/Fulfillment Node identity collapse;
10. listing disappears from partial enumeration but point read is active -> no absence=>closed inference;
11. seller Order search completes while canceled Orders are excluded -> coverage explicitly excludes cancellations rather than claiming complete Sales population;
12. a canceled Order notification was missed -> seller enumeration alone cannot fabricate recovery/completeness; required recovery coverage remains partial/unknown until proved;
13. Order and Shipment disagree -> current Shipment evidence governs shipping/fulfillment observation;
14. fiscal prerequisite/label/SLA changes -> current provider requirement evidence controls Fulfillment closure, not stale transport success;
15. only Full/provider-operated fulfillment is present -> support it honestly without claiming internal physical execution;
16. partial Return/refund consequences occur -> Post-Sale and Economics remain separate authorities and scopes;
17. a second marketplace arrives -> consumer-owned semantic ports remain usable without exporting Mercado Livre vocabulary into domains.

### 4.11 B2 reopen / stop triggers

Reopen only the implicated accepted decision when material evidence shows:

1. required Mercado Livre resource cannot fit Installation-qualified external identity without a new canonical MPC identity -> targeted D2 review;
2. necessary consumer dependency is absent from D1 -> STOP / targeted D1 review;
3. provider effect semantics cannot fit D3 scope/accepted-rejected-pending-ambiguous/reconciliation rules -> targeted D3 review;
4. real Mercado Livre capability makes an accepted Product 1.0 outcome impossible or materially different -> targeted D0 review;
5. no seller-writable availability lane exists while D0 controlled Availability convergence is claimed complete -> targeted product-proof adjudication;
6. required internally operated Fulfillment Node proof has no seller-operated lane -> operator/product-proof adjudication; provider work is not relabeled as internal work;
7. selected flow requires real composite semantics that materially change availability/materialization -> use the existing D0/D1 composition reopen trigger;
8. provider Price Automation or another external control makes a required lane unavailable -> explicit unsupported/external-required or targeted product decision, not silent provider-control removal;
9. current official/provider real behavior materially changes from B2 contracts -> reopen only the affected operation;
10. a second real provider exposes a repeated technical failure class that cannot be handled cleanly by concrete adapters/consumer-owned ports -> consider the smallest proven shared mechanism then.

Package preference, current adapter shape and hypothetical future marketplace features are not reopen evidence.

### 4.12 B2 evidence basis

B2 current-provider facts were independently challenged/reverified during the review cycle against official Mercado Livre documentation families covering:

- User Products / Item publication / shared-field propagation;
- Price APIs and Price Automation;
- distributed/multi-origin/User Product stock and provider-managed Full stock;
- Orders, `cancel_detail`, seller-search retention/scope;
- Shipments, provider requirements, labels and SLA/deadline surfaces;
- seller reputation/moderation/restrictions;
- Claims/Returns and provider-effective actions;
- notifications and bounded missed-feed recovery.

Reference platforms such as ANYMARKET, Amazon, Casas Bahia and Mirakl were used only as failure-class evidence that provider/fulfillment responsibility is mode-sensitive. Their module/service topology is not MPC authority.

Historical repository live probes remain time-bound evidence, never current Installation state by inheritance.

---

## 5. B1 evidence basis

External facts used to ground B1 were reverified against current official sources during the B1 review cycle. The source families included:

### Mercado Livre

- authentication/authorization and OAuth/token lifecycle;
- user identity and `/users/me`;
- notification/resource-reread semantics and bounded missed-notification recovery;
- seller item search/scan behavior.

Official documentation family: `developers.mercadolivre.com.br`.

### Sankhya

- OAuth 2.0 `client_credentials` + `X-Token` Gateway authentication;
- integration guidance / API Gateway posture;
- `loadRecords` pagination;
- `modifiedSince` / change-log dependency;
- authentication migration guidance.

Official documentation family: `developer.sankhya.com.br`.

These sources are external evidence, not a second architecture authority. If unstable provider behavior changes materially, reopen only the decision that depended on it.

---

## 6. Review and authority protocol

D4 uses the accelerated protocol established by the earlier stages:

1. GPT prepares a coherent batch from repository authority and claim-specific evidence.
2. The operator approves the batch direction for independent challenge.
3. A disposable `D4-B<n>-REVIEW-CANDIDATE.md` may be created and is explicitly **not architecture authority**.
4. The operator invokes Fable separately; Fable reconstructs the authority path independently and appends findings to `AI-DIALOG.md`.
5. Reviewer findings remain evidence; GPT independently adjudicates material findings against repository authority/evidence.
6. Material disagreement receives only the minimum further reviewer convergence needed.
7. The operator explicitly ratifies the converged batch before canonical consolidation.
8. Accepted/canonical batch meaning is filed here; the review candidate is deleted once its durable meaning is captured.
9. A batch may remain **OPEN** when operator-ratified canonical contract semantics still have an explicit evidence/proof gate required for whole-batch closure; the router names that gate and blocks any stronger status claim.
10. At D4 closure, run final Global Coherence + YAGNI / Overengineering / Future-Cost review before whole-stage ratification.

`AI-DIALOG.md`, review candidates, chat summaries and reviewer statements are never target authority.

---

## 7. Current D4 state / exact next action

D4 is **OPEN / ACTIVE**.

- **D4-B1 — External Contract Grounding: ACCEPTED / CANONICAL.**
- **D4-B2 — Mercado Livre Operational Contract: OPEN / CANONICAL CONTRACT CORE; INSTALLATION EVIDENCE GATE OPEN.**
- **D4-B3 — Sankhya Business-System Contract: NOT YET OPENED.**
- **D4-B4 — Market / Economics / Settlement Contract: NOT YET OPENED.**
- **Final D4 Global Coherence + YAGNI / Overengineering / Future-Cost review: NOT STARTED.**

No D0/D1/D2/D3 reopen is currently required.

Exact next action: **execute/admit the D4-B2 §4.8 read-only Mercado Livre Installation Evidence Gate and establish the smallest real supported/proof lane set.** Do not infer current seller tags, stock writability, Price Automation or fulfillment modes from historical probes or generic provider documentation.

Until that gate is resolved under the current router, do not open B3/B4 or claim B2 accepted as a whole.

Product implementation remains blocked until D9 is accepted.