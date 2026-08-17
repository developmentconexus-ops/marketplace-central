# D4 — External Integrations

> **Status:** OPEN / ACTIVE — D4-B1 accepted and canonical; D4-B2 next  
> **Program:** Architecture Rebaseline / Technical System Design  
> **Parent authorities:** `D0-PRODUCT-SYSTEM-DEFINITION.md`, `D1-DOMAINS-BOUNDARIES.md`, `D2-IDENTITY-TENANT-DATA-OWNERSHIP.md`, `D3-COMMUNICATION-EVENTS.md`  
> **Method:** DevelopmentConexus Engineering Method v1.0.0  
> **Opened:** 2026-08-17  
> **B1 accepted:** 2026-08-17

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

#### D4-B2 — Mercado Livre Operational Contract — NEXT

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

## 4. B1 evidence basis

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

## 5. Review and authority protocol

D4 uses the accelerated protocol established by the earlier stages:

1. GPT prepares a coherent batch from repository authority and claim-specific evidence.
2. The operator approves the batch direction for independent challenge.
3. A disposable `D4-B<n>-REVIEW-CANDIDATE.md` may be created and is explicitly **not architecture authority**.
4. The operator invokes Fable separately; Fable reconstructs the authority path independently and appends findings to `AI-DIALOG.md`.
5. Reviewer findings remain evidence; GPT independently adjudicates material findings against repository authority/evidence.
6. Material disagreement receives only the minimum further reviewer convergence needed.
7. The operator explicitly ratifies the converged batch before canonical consolidation.
8. Accepted batch meaning is filed here; the review candidate is deleted.
9. At D4 closure, run final Global Coherence + YAGNI / Overengineering / Future-Cost review before whole-stage ratification.

`AI-DIALOG.md`, review candidates, chat summaries and reviewer statements are never target authority.

---

## 6. Current D4 state / exact next action

D4 is **OPEN / ACTIVE**.

- **D4-B1 — External Contract Grounding: ACCEPTED / CANONICAL.**
- **D4-B2 — Mercado Livre Operational Contract: NEXT / NOT YET OPENED.**
- **D4-B3 — Sankhya Business-System Contract: NOT YET OPENED.**
- **D4-B4 — Market / Economics / Settlement Contract: NOT YET OPENED.**
- **Final D4 Global Coherence + YAGNI / Overengineering / Future-Cost review: NOT STARTED.**

No D0/D1/D2/D3 reopen is currently required.

Exact next action: **open D4-B2 from this canonical B1 and claim-specific current Mercado Livre official evidence.** B2 must include the identifier-evidence obligations named in §3.13 and must not choose D7 runtime topology.

Product implementation remains blocked until D9 is accepted.
