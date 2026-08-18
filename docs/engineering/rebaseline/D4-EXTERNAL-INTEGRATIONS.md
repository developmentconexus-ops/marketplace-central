# D4 — External Integrations

> **Status:** OPEN / ACTIVE — D4-B1 accepted/canonical; D4-B2 accepted/canonical; D4-B3 accepted/canonical; D4-B4 next  
> **Program:** Architecture Rebaseline / Technical System Design  
> **Parent authorities:** `D0-PRODUCT-SYSTEM-DEFINITION.md`, `D1-DOMAINS-BOUNDARIES.md`, `D2-IDENTITY-TENANT-DATA-OWNERSHIP.md`, `D3-COMMUNICATION-EVENTS.md`  
> **Method:** DevelopmentConexus Engineering Method v1.0.0  
> **Opened:** 2026-08-17  
> **B1 accepted:** 2026-08-17  
> **B2 accepted:** 2026-08-17  
> **B3 accepted:** 2026-08-18

## 1. Purpose and boundary

D4 defines the concrete external acquisition, translation, capability, requirement, authentication, coverage and effect-reconciliation contracts by which Mercado Livre, Sankhya and later accepted external systems participate in the MPC operating model already fixed by D0–D3.

D4 does **not** create an Integration business domain and does not move business authority into adapters.

D4 decides, proportionately:

- concrete provider/business-system namespace binding under D2 identities;
- authentication/credential protocol required to establish the right external namespace;
- authoritative read/reread surfaces and source-specific observation authority;
- point/enumeration/delta/notification coverage semantics;
- Integration Support and Provider Effective Capability/Requirement evidence;
- concrete external-effect acceptance/ambiguity/reconciliation obligations;
- provider/business-system identifiers, requirements and evidence required by D1 consumers;
- transport admissibility where external/provider constraints materially affect the target.

D4 does **not** choose:

- MPC HTTP/OpenAPI routes, wire errors or client contracts — **D5**;
- frontend screens/forms/projection topology — **D6**;
- workers, schedulers, queues, retries, locks, token refresh scheduling, cursor persistence, transaction/process/deployment topology — **D7**;
- end-to-end golden-flow assembly/proof — **D8**;
- product implementation, which remains blocked until D9 is accepted.

Current code, schemas, OpenAPI, tests, runtime and legacy ADR implementation shapes remain evidence only unless this authority explicitly rehomes their meaning.

---

## 2. Imported parent invariants

D4 imports rather than re-decides:

1. **Consumer owns meaning; adapter owns protocol.**
2. Provider DTO/status/auth/pagination vocabulary never becomes domain authority.
3. Organization scope is explicit; provider/source identifiers and credentials never substitute for it.
4. Marketplace Installation and SourceInstance qualify the correct external namespace without becoming provider credentials/IDs.
5. Externally authoritative resources retain source-qualified identity.
6. Notification/poll/callback payload is acquisition evidence, not automatically MPC domain truth or a D3 domain event.
7. Known value, known empty/absent, unknown and unavailable remain distinct; partial observation is not closure.
8. `accepted != completed != externally applied != converged`; ambiguous possible acceptance is never blindly retried.
9. Provider PII is minimized.
10. Integration Support, Provider Effective Capability/Requirement and Effective Business Capability remain distinct authorities.
11. Mercado Livre is the first marketplace proving the boundary; Sankhya is the first business system proving the business-system boundary. Neither provider becomes MPC ontology.
12. No generic integration/provider/ERP/workflow/customer framework is introduced without a real repeated failure class/consumer.

---

# 3. D4-B1 — External Contract Grounding — ACCEPTED / CANONICAL

**Outcome:** `RESTRUCTURE NOW` for the former direct-Oracle Sankhya target assumption; otherwise `CURRENT STRUCTURE CONFIRMED`. No D0/D1/D2/D3 reopen is required.

## 3.1 Governing external-contract invariant

> **Every external fact/effect entering MPC is qualified by the correct Organization + external namespace, acquired through a contract whose authority/coverage is explicit enough to preserve honest knowledge state, translated through a consumer-owned semantic boundary and, for consequential effects, correlated to an authoritative reconciliation surface. Provider protocol never acquires MPC business authority.**

```text
D1/D2 consumer meaning
        ↓ consumer-owned port
D4 concrete adapter
  - namespace binding
  - auth/protocol
  - operation authority + coverage
  - provider capability/requirements
  - authoritative reread/reconciliation
        ↓
external system
```

D7 later supplies runtime mechanics around this boundary.

## 3.2 Protocol / semantic-boundary fence

1. Consumer-owned ports express business meaning.
2. Concrete adapters own endpoint/wire DTO/auth/pagination/provider-error translation.
3. Provider-local shared HTTP/auth/pagination machinery may exist as mechanism, never business authority.
4. One adapter may implement several consumer ports; this does not imply one universal provider business port.
5. D4 may expose operation-specific Integration Support/Provider Effective evidence but does not create a universal Provider/Resource/Capability graph.
6. Assembly/factory/self-registration/process topology remains composition/D7 mechanism.

Legacy generic plugin/self-registration/catalog/factory shape has no target authority by inheritance.

## 3.3 Mercado Livre Installation ↔ seller binding

1. A Marketplace Installation binds one Organization marketplace participation to one authoritative provider seller/account namespace.
2. External seller identity remains provider-owned and distinct from Installation and Organization.
3. Authorization/re-authorization establishes seller identity from provider-authoritative evidence such as provider user identity surfaces.
4. Token/credential refresh for the same seller does not create a new Installation.
5. A different seller must never silently rebind an existing Installation; mismatch fails closed and requires explicit reconfiguration.
6. Provider user IDs remain opaque and must not be narrowed to Int32.
7. Acquisition-time attribution also fails closed where an acquired resource exposes an authoritative seller/namespace marker contradicting the bound Installation.
8. Where no per-resource marker exists, the established authenticated Installation binding remains the control.

## 3.4 Sankhya SourceInstance binding

1. `SourceInstance` is Organization-qualified and denotes one logical Sankhya authoritative namespace/environment.
2. Production and sandbox/test remain distinct; visually equal native keys do not collapse namespaces.
3. Credential rotation does not re-identify the SourceInstance while the authoritative namespace remains the same.
4. A materially different namespace/environment requires explicit rebinding/new SourceInstance treatment.
5. Source identity is independent of protocol mechanics.
6. Where sanctioned source evidence exposes an authoritative namespace marker, mismatch fails closed; otherwise configured/authorized source binding remains the control.

## 3.5 Sankhya transport — API Gateway is the target

**RESTRUCTURE NOW — transport only.**

For MPC↔Sankhya integration, the sanctioned Sankhya **API Gateway is the target transport**.

**Direct Oracle/database access is not part of target architecture and is not a fallback path.**

The former Oracle/godror path is historical evidence from a period when a usable sanctioned API path was not established. Its transport-independent lessons remain: Sankhya is external, consumers depend on MPC-owned ports/adapters, driver/protocol does not leak into business code, and a wrong legacy source store is not revived for convenience.

A materially required claim that cannot be satisfied correctly/operationally through a sanctioned Gateway/API capability triggers **STOP / SPLIT PREREQUISITE** and explicit operator/architecture adjudication. It never authorizes Oracle implicitly.

## 3.6 Authentication lifecycle is protocol, not identity

### Mercado Livre

- OAuth/access/refresh credentials are adapter/runtime secrets, never business identity/state.
- credential refresh does not re-identify the seller/Installation.
- auth failure makes acquisition unavailable/auth-invalid, not business data absent.

### Sankhya

- target Gateway auth uses OAuth 2.0 `client_credentials` plus provider-required `X-Token` and Bearer token;
- production/sandbox credentials/tokens are not mixed;
- credential changes do not alter SourceInstance while the authoritative namespace is unchanged;
- token cache/refresh locking/scheduling remains D7.

## 3.7 Notification is a trigger; authoritative reread owns current observation

1. Provider notification is acquisition evidence/pointer, not current business truth.
2. Current material provider meaning comes from authoritative resource reread.
3. Duplicate/out-of-order notifications may cause repeated rereads without duplicate business truth or state regression.
4. Bounded missed-feed/recovery features do not become complete durable history by assumption.
5. Notification outage/gaps do not create a completeness claim.
6. Receiver queue/retry/reconciliation scheduling remains D7.

## 3.8 Coverage is operation-scoped

### Point observation

- proves only the exact source-qualified resource/scope established by the endpoint;
- `not found` is known absence only when source semantics prove it;
- timeout/auth/rate-limit/outage/parse failure is unavailable, not absent.

### Enumeration

- proves only the source-defined scope actually traversed;
- every required page/scan segment must complete before a complete-enumeration claim;
- early stop/failure remains partial;
- completed traversal does not invent snapshot isolation the provider did not establish.

### Delta/change

- proves only the source-defined change set under its explicit prerequisites/window;
- does not prove the complete current population;
- empty delta with unproven/disabled change-log prerequisites never becomes “nothing changed”.

### Notification

Carries no completeness claim by itself.

Freshness remains consumer/use-sensitive under D3.

## 3.9 Failure/knowledge semantics survive the boundary

Adapters may normalize external failures, but never map unavailability to plausible business values. Unsupported capability is explicit. Provider-native error DTO/text remains adapter-local. Provenance preserves sufficient source identity and material source/acquisition time without requiring a universal raw-payload archive.

## 3.10 Capability authority fence

1. **Integration Support / Descriptor — D4:** what the adapter implements.
2. **Provider Effective Capability / Requirement — D4 translated provider evidence:** what this installation/source/resource/context currently permits/requires.
3. **Effective Business Capability — D1 owner:** whether MPC may/should perform the action now under business state/policy/readiness/authorization.

Adapter support never bypasses D1 owner validity, and a business domain never fabricates provider support.

## 3.11 Admission rule for consequential external effects

Every admitted external-effect contract identifies:

1. Installation/SourceInstance-qualified target;
2. consumer-owned intent/correlation anchor;
3. material provider/source prerequisites;
4. what the response actually proves — rejection, accepted submission, synchronous effect, pending or ambiguity;
5. when possible acceptance can become ambiguous;
6. authoritative reread/reconciliation surface;
7. member-level outcome where multi-target work can partially succeed;
8. source occurrence/result discriminator only where same-vs-distinct correctness needs it.

Transport 2xx never becomes convergence by itself.

## 3.12 Legacy ADR disposition at B1

- ADR-004 generic plugin framework target meaning: superseded.
- ADR-010 polling-only/no-webhook D4 meaning: superseded; D7 runtime residue remains.
- ADR-006/007 direct-Oracle/godror target meaning: historical/superseded by B1.
- ADR-003 credential/identity prerequisite rehomed; its old strict implementation sequence is not target authority.
- ADR-014/020/032 market/economics residue remains B4 evidence only.
- ADR-015 provider/listing residue is adjudicated by B2 below.

## 3.13 B1 reopen/stop triggers

Reopen only the implicated decision when material evidence shows an accepted identity/authority/communication/product assumption fails, or a required Gateway capability cannot be satisfied by any sanctioned operation. Any proposal to reintroduce Direct Oracle requires explicit operator-requested reopen with new evidence. Framework preference/current-code convenience/hypothetical providers are not reopen evidence.

---

# 4. D4-B2 — Mercado Livre Operational Contract — ACCEPTED / CANONICAL

**Outcome:** `CURRENT STRUCTURE CONFIRMED` with a bounded concrete Mercado Livre contract and one real selected first-flow lane set. No D0–D4-B1 reopen is required.

## 4.1 Governing Mercado Livre invariant

> **For every Mercado Livre operation required by Product 1.0, the adapter maps a consumer-owned MPC query/intent to the currently applicable Installation-qualified provider resource/context, preserves provider authority/coverage/preconditions/effect scope and returns only the semantic evidence needed by the D1 owner. Mercado Livre resource topology never becomes MPC business ontology, and provider context never silently widens intended/authorized scope.**

## 4.2 Installation, Item/User Product and listing surface

1. Provider seller identity remains the B1-bound seller namespace.
2. Seller reputation/restrictions/moderation are provider evidence that may inform Portfolio/Offering attention/effective capability without creating reputation-management authority.
3. `item_id` is provider Listing identity; legacy Variation is an external child identity only where the actual listing uses that model.
4. `user_product_id`, Family and Catalog Product remain provider-native product/catalog topology, not MPC Product entities.
5. Legacy and User Product modes may coexist provider-side; the adapter follows current authoritative resource evidence rather than a static assumption.
6. Provider migration/relationship changes never rewrite MPC Product identity or Readiness authority.
7. Offering owns Listing Intent/convergence; provider create/2xx is not convergence until authoritative current state is reconciled.
8. Listing enumeration completeness depends on actual provider scan/paging rules; absence from partial/unknown coverage never closes/deletes a Listing.

### Product↔channel identifier evidence

D4 supplies provider evidence; Readiness owns correspondence sufficiency. Seller SKU/custom field, GTIN/EAN, Item/Variation/UP/Family/Catalog/category/attribute relationships remain evidence only. No single field becomes identity law, `SELLER_SKU == CODPROD` is not resurrected, contradictory evidence is preserved, and identity-bearing outbound values remain subject to D2 pre-dispatch consistency/corroboration rules.

## 4.3 Price and Availability effects

### Price

1. Current price observation and current write mechanism may be different provider surfaces.
2. Offering owns Price Intent/business validity; D4 owns provider protocol/effective capability.
3. Provider Price Automation or another external control may make direct price write unavailable; D4 does not disable provider controls automatically to make MPC succeed.
4. Mixed/ignored-field 2xx is not convergence; authoritative current price reread controls the conclusion.

### User Product shared-effect scope

Provider shared fields may propagate asynchronously across Items related to one User Product, including availability where applicable. Before dispatching a shared-field effect, D4 establishes current provider blast radius sufficiently to ensure one narrow MPC intent is not silently widened. Relation drift is revalidated at execution time when material.

### Stock topology / writability

1. `available_quantity` is not one universal stock contract.
2. Provider stock locations (`store_id`, `network_node_id`, `selling_address`, `seller_warehouse`, `meli_facility`, etc.) remain external evidence and do not become Inventory Source/Fulfillment Node by convenience.
3. Seller-managed does not automatically mean API-writable; writability depends on concrete resource, site applicability, seller configuration/tags and current resource/listing state.
4. Provider-managed Full/`meli_facility` stock may be observable without being seller/MPC-writable.
5. Unsupported/provider-managed paths remain explicit, never successful no-ops.
6. Where a provider stock surface uses `x-version`/equivalent, stale 409 is a definitive rejected precondition followed by authoritative reread/redecision, not ambiguous retry.
7. D0's controlled Availability-convergence claim requires at least one real writable proof lane; Full-only observation cannot prove MPC-controlled convergence.

## 4.4 Sale, Shipment and fulfillment contract

### Order acquisition / coverage

1. Order notifications are pointers; Marketplace Sales commits sale meaning only after required authoritative Order reread.
2. `GET /orders/{id}` is point authority when the source-qualified ID is known.
3. Seller Order search is an enumeration surface with provider filters/retention, not permanent Sales history authority.
4. Current official documentation and real seller behavior conflict about canceled-Order inclusion in seller search; completed traversal therefore does **not** prove a cancellation-inclusive Sales universe.
5. Cancellation recovery/completeness cannot rely on seller Order-search enumeration alone until the exact recovery universe is proven.
6. Provider retention never replaces accepted D2/D3 durable historical meaning required beyond provider windows.

### Order and Shipment remain separate

- Provider Order remains Sale/checkout evidence.
- Provider Shipment is a distinct source-qualified external identity and owns current shipping/fulfillment evidence when material.
- `shipping.id` is correlation evidence; embedded Order shipping does not become permanent Shipment authority.
- Provider Pack remains provider-native when relevant; no universal MPC Pack entity is created.

### Cancellation/fraud/provider stops

Authoritative `cancel_detail`/equivalent provider evidence remains source fact. Fraud/stop-shipment conditions may block progression even after earlier payment evidence. D4 translates evidence; Post-Sale/Sales/Fulfillment retain their D1 meanings.

### Effective fulfillment context / provider requirements

1. Marketplace brand alone does not determine physical responsibility.
2. Adapter translates current provider logistics/fulfillment context; provider codes stay adapter-local.
3. Fulfillment Lifecycle owns provider-requirement closure for flows MPC claims to control.
4. Provider-operated paths remain explicit rather than being relabeled as internally controlled.
5. Provider fiscal/document/label/readiness requirements remain provider-authoritative evidence; Business-System Materialization owns Sankhya Business Order/Invoicing Intents under B3.
6. Provider submission/2xx is not readiness closure; authoritative Shipment/provider reread controls provider state.
7. Provider deadlines/SLA remain external authoritative evidence and stay distinct from MPC Internal Operational Targets.
8. D7 owns timers/schedulers/attention mechanics.

## 4.5 Essential post-sale contract

1. Provider Claim/Case/Return/reverse Shipment remain Installation-qualified provider resources/evidence, not duplicate MPC aliases.
2. Post-Sale Resolution remains MPC canonical resolution/correlation/closure owner.
3. Consequence scope remains material item/line/quantity scoped; one provider Return does not globally close a Sale.
4. Reverse Shipment remains distinct from original Shipment.
5. Provider `available_actions`/equivalent is Provider Effective Capability, never automatic MPC permission.
6. General complaint management, buyer chat/Q&A and reputation workflows remain outside Product 1.0.
7. Provider Payment/Refund/Fee/Adjustment/Settlement/Payout authority remains B4; physical return closure never fabricates financial closure and vice versa.
8. A measured 403 on one provider-native Return detail surface remains explicit access/scope Unknown rather than provider absence.

## 4.6 External-effect reconciliation specialized for Mercado Livre

Every admitted write:

- is Installation-qualified and correspondence-valid;
- uses current enough provider-effective prerequisites/capability;
- does not widen intended/authorized scope silently;
- classifies provider response no stronger than evidence;
- preserves asynchronous/partial/ignored-field behavior;
- uses authoritative resource reread for current result/convergence;
- never blindly retries ambiguous possible acceptance;
- treats definitive stale-precondition rejection as reread/redecision.

## 4.7 Installation Evidence Gate — CLOSED / PASS

A fresh read-only real-dependency probe on 2026-08-17 selected the smallest current Product 1.0 lane set without implementing every Mercado Livre mode:

- seller active on MLB and currently `user_product_seller`;
- completed current seller scan: 34/34 Items User Product model, zero legacy variations;
- Item↔UP 1:1 across that measurement;
- no `warehouse_management`/`multiwarehouse`, Full/`meli_facility` or multi-origin stock observed;
- selected Availability candidate: non-multi-origin Item-path `available_quantity`, retaining shared-UP blast-radius revalidation;
- Price candidates were not Price Automation-blocked at measurement time;
- selected Sale/Fulfillment lane: seller-operated `me2 / xd_drop_off`; no Full/flex lane observed;
- Shipment SLA surface live, including delayed evidence;
- real Claim/Return + reverse Shipment evidence available for later D8 proof.

These are time-bound Installation facts, not permanent provider promises.

Residual proofs:

- **R1 — D8:** first controlled Price/Availability write + authoritative reread/convergence.
- **R2 — D8 using accepted B3:** live selected-lane fiscal/invoice/label progression in the end-to-end golden flow.

## 4.8 B2 YAGNI / provider-overfit fence

B2 does not introduce a universal Provider/ProviderResource graph, canonical MPC UserProduct/Family/ProviderWarehouse/OperatingMode/Claim/Return aliases, universal marketplace stock surface, universal raw payload archive or support for every documented fulfillment mode. It does not deactivate provider Price Automation merely to obtain control. Unselected modes remain explicit unsupported/external-required until a real consumer requires them.

## 4.9 Legacy disposition / reopen triggers

ADR-015 old read-only listing module/composite-ID/manual-refresh/absence=>closed target shape is historical; D1/D2/D4 now own the surviving meanings. ADR-022/028 legacy identity formulas remain superseded while B2 supplies concrete identifier/correspondence evidence to Readiness.

Reopen only when material provider evidence invalidates an accepted Product/identity/communication assumption or makes a required first-flow capability impossible/different. A second real provider may justify only the smallest proven shared mechanism. Package preference/current adapter shape/hypothetical features are not reopen evidence.

---

# 5. D4-B3 — Sankhya Business-System Contract — ACCEPTED / CANONICAL

**Outcome:** `CURRENT STRUCTURE CONFIRMED` for the provider-independent business-system boundary plus concrete Sankhya realization, with bounded SourceInstance configuration and explicit D7/D8 proofs. No D0/D1/D2/D3/D4-B1/B2 reopen is required.

## 5.1 Governing invariant / Global Maximum

> **For each Product 1.0 business-system fact/effect, MPC depends on a consumer-owned semantic contract whose external evidence/effect is Organization + SourceInstance qualified. The concrete Sankhya adapter uses only sanctioned operations and bounded SourceInstance-specific bindings, preserves provider-native granularity/provenance/partiality and reconciles consequential effects through authoritative reread/correlation. Sankhya-native topology, party vocabulary and customer configuration never become MPC business ontology, and future substitutability does not authorize a generic ERP/workflow/party framework before a second real consumer proves the abstraction.**

Credible alternatives:

- **Sankhya-shaped MPC core** (`Business Order = TOP 313`, `Customer = CODPARC`, `Marketplace = CODTIPVENDA`, `Lot = CONTROLE`) — **REJECTED**.
- **Generic ERP/workflow/customer framework now** — **REJECTED**; no second real consumer proves a stable common ontology/workflow language.
- **Consumer-owned semantics + concrete Sankhya adapter + bounded SourceInstance binding** — **GLOBAL MAXIMUM / ACCEPTED**.

Strong replacement test:

> **Sankhya is the first concrete business system proving the MPC business-system boundary; it is not the MPC business-system model.**

If MPC/domain configuration must express arbitrary provider-native choreography/customer master concepts merely so another ERP can fit later, this boundary has failed.

## 5.2 Environment, auth and bounded sanctioned reads

1. Production and sandbox/test are distinct SourceInstance environments/namespaces.
2. Target auth uses sanctioned Gateway OAuth/client credentials + X-Token/Bearer mechanics; measured short token TTL/caching mechanics remain D7.
3. Dedicated REST resources are preferred when semantically sufficient.
4. `CRUDServiceProvider.loadRecords` is admitted only through a narrow read fence:
   - named sanctioned `rootEntity`;
   - explicit minimum fieldset;
   - criteria only over fields of that root entity using ordinary comparison/logical operators and bound parameters;
   - related data only through provider-declared relation/path when materially required;
   - no subqueries, cross-table criteria references, Oracle pseudo-tables/functions, arbitrary SQL expressions or query-language escape hatches;
   - needing such an expression is a capability finding, never authorization for Oracle-via-HTTP;
   - no entity read becomes an ERP mirror by convenience.
5. Provider paging/ordering/date-format/filter quirks remain adapter-local; syntactic acceptance does not prove semantic honoring.
6. `modifiedSince`/change-log claims remain prerequisite/coverage bound; empty result with unproven change-log coverage is not “nothing changed”.

## 5.3 Product, company/location, inventory and cost

### Product

External Product reference remains:

```text
SourceInstance + native Product key
```

For current Sankhya this key is `CODPROD`; it remains provider-native. Brand/reference/NCM/barcode/fiscal attributes are external evidence; no universal GTIN field or `REFERENCIA == EAN` identity law is introduced. Readiness owns Product↔channel correspondence sufficiency.

### Company/location

`CODEMP`/`CODLOCAL` remain native external references and do not automatically equal Selling Entity, Inventory Source or Fulfillment Node. Marketplace Sales owns transaction-specific Selling Entity attribution.

### Inventory

Measured sanctioned evidence establishes:

```text
REST estoque = native ESTOQUE - RESERVADO
```

The REST surface is useful but lossy when reservation/source decomposition matters. Bounded sanctioned inventory reads preserve necessary provider dimensions. Negative net stock is a real external observation, not clamped zero. Resource absence is not automatically `stock=0`.

`CONTROLE` is preserved as an opaque provider inventory-partition value when present; it is not canonically named Lot/Batch/Tonality/Serial. Current measured Mercado Livre sold population was control-free. Controlled-product marketplace Availability/Materialization remains unclaimed until satisfiability/interchangeability and selection timing are proven.

When native Business Order materialization creates/changes an external inventory commitment, inventory observation must preserve that effect so Availability can derive current sellability. Materialization does not acquire allocation authority.

### Cost

Bounded sanctioned `Custo` observations preserve company/Product/time/provider-local qualifiers. Native variants such as `CUSGER`, `CUSREP`, `CUSSEMICM`, `CUSMED` remain observations; Commercial Economics owns Cost Basis selection/interpretation.

## 5.4 Expected Tax — G1 PASS

`POST /v1/fiscal/impostos/calculo` is the selected sanctioned Expected-Tax surface. MPC does not copy Sankhya fiscal rules into an MPC tax engine.

Stable SourceInstance binding proven on 2026-08-18:

```text
notaModelo = 898307
TIPMOV = Z
CODEMP = 1
CODTIPOPER = 313
CODVEND = 1019
SERIENOTA = PA
CODTIPVENDA = 27
CODPARC(model binding) = 142005
```

The model Partner exists to satisfy current Sankhya authorization/configuration for negotiation type `27`; it is **not** the fiscal customer of each simulation. The API requires `codigoCliente`, fails closed when it is missing, and a requested SP customer retained SP fiscal treatment although the model Partner is MG. The model Partner therefore does not silently substitute for transaction-customer meaning.

Final designated-model proof, no company/TOP override:

- MG customer: ICMS CST 60 / 0; no DIFAL; PIS 1.65%; COFINS 7.6%;
- SP customer: ICMS CST 00 / 12% / 20.40; DIFAL destination 10.20; PIS 1.65%; COFINS 7.6%;
- results matched the earlier exact type-27 target-lane pair and realized ICMS documents to the cent;
- no persisted negative-NUNOTA headers/items/financial rows; stock/reservation and model state remained unchanged by the probes.

Therefore **G1 = PASS / CLOSED FOR B3**.

### Expected-Tax adapter obligations

1. **F-G1-1 — provider `valorBase` on PIS/COFINS/CSLL is not reliably the arithmetic base used.** Preserve provider-returned component/value/provenance; do not reconstruct attributable tax as `base × rate` where evidence disproves that identity.
2. **F-G1-2 — unknown request fields may be silently ignored.** Use a typed/pinned request contract; HTTP 200 does not prove intended monetary inputs were consumed. Validate materially echoed quantity/unit price/discount inputs where exposed.
3. **F-G1-3 — top-level `despesasAcessorias` request shape is unestablished.** Remain Unknown until a real consumer requires/proves it.
4. **F-G1-4 — IPI and retained ICMS-ST components were not established as separately observable in the measured response.** Absence is not proof of non-applicability. A future claimed flow that materially requires them must prove them or remain Unknown/unsupported.
5. **F-G1-6 — `CODTIPVENDA` is fiscally determinant.** Current binding preserves type `27`; ICMS agreement alone is not equivalence evidence for a substitute type because a wrong type was measured to silently change PIS while still returning 200.

PF and PJ occur in measured e-commerce/marketplace evidence. The specific out-of-state PJ-contributor branch remains unexercised and therefore Unknown, not an artificial B3 closure gate. A golden flow that materially depends on an unproven fiscal branch must validate it or fail honestly rather than infer equivalence.

Expected-Tax binding drift must be detectable no stronger than sanctioned evidence permits. Material revalidation includes the model/type/Partner/source configuration required by the proven binding; cadence/cache mechanics remain D7.

Historical realized tax remains realized/B4 evidence, not a replacement Expected-Tax engine.

## 5.5 Native order observation and progression

Native business-order identity remains:

```text
SourceInstance + native NUNOTA
```

REST order enumeration is useful for bounded discovery but measured point filters were unreliable; bounded sanctioned `CabecalhoNota` reread by NUNOTA is the current authoritative point/reread surface for consequential native-order state. When enumeration and point reread disagree, point reread governs current consequential state.

Raw `STATUSNOTA`/`PENDENTE` are adapter-local. MPC does not acquire a universal Confirmation lifecycle stage.

Provider-independent meaning is:

> **the native business order has reached the externally required state sufficient for the next claimed materialization progression.**

Current Sankhya progression realization uses `CACSP.confirmarNota` and observed `A → L`; this remains adapter-local.

## 5.6 Business-System Party Resolution

A Business Order Intent may require resolving the source-native business party required by the selected business system. This is a bounded Materialization prerequisite, not Customer Master/CRM authority.

```text
Marketplace Sale
  ├─ marketplace buyer/account evidence
  └─ fiscal/billing-party evidence
             ↓
Business-System Party Resolution
             ↓
source-native business-party reference
```

Rules:

1. Marketplace buyer/account identity is provenance/context, not automatically fiscal party.
2. A prior explicit native resolution may be reused while materially compatible with current fiscal-party evidence and native party state; another delivery destination does not by itself invalidate fiscal-party correspondence.
3. Exactly one sufficiently established compatible native match may be used.
4. Zero native matches may authorize native creation only when every required identity-bearing fact is known from legitimate transaction evidence; creation does not silently decide that shipping destination becomes customer master data.
5. Multiple matches or material fiscal/identity contradictions are `AMBIGUOUS`: no guessed selection, first-result-wins or new duplicate creation.
6. Native party create/update is a consequential external effect.
7. MPC may preserve bounded resolution/correspondence lineage when needed to avoid repeating a human-adjudicated ambiguity, without owning customer attributes/lifecycle.
8. Concurrent/repeated materializations for the same unresolved fiscal identity must not independently create duplicate native parties; D7 chooses the enforcement mechanism.

Durable Party Resolution state is **KEEP — PRESENT CORRECTNESS NEED** because a measured seven-way active native duplicate means rereading the source cannot reconstruct a prior human adjudication.

### Current Sankhya realization

- source-native party reference: `CODPARC`;
- strongest current legal-document lookup signal: `CGC_CPF`, using ordinary equality with a string-bound parameter on bounded `Parceiro` reads;
- the dedicated customer REST family does not provide the safe point-lookup behavior needed for correspondence;
- legal-document uniqueness is measurably not enforced in this SourceInstance;
- incumbent e-commerce origin custom fields are not reliable correspondence authority;
- `Parceiro`, `CODPARC`, `CGC_CPF` and their lookup protocol remain adapter/SourceInstance realization only.

## 5.7 Business-System Destination Realization

Delivery destination is a separate Materialization prerequisite from Party Resolution:

```text
Marketplace Sale
  └─ delivery-recipient/destination evidence
             ↓
Business-System Destination Realization
  + resolved native party when required
             ↓
safe source-supported destination representation
        OR explicit Work / external-required
```

Rules:

1. Fiscal/billing party is not automatically delivery recipient.
2. Delivery destination is transaction evidence and is not customer-master address by inheritance.
3. Party Resolution success does not imply the current destination is safely representable.
4. Use the least-destructive sufficiently-proven native realization when the business system can represent the destination safely.
5. If no safe realization is proven, return explicit `external-required` / Work; never silently drop the destination or fabricate equivalence with master data.
6. Never create another native customer merely to hold another address.
7. Never perform unattended overwrite of a registered/master customer address merely because a marketplace transaction carries a different shipping address.
8. Native contact/address/destination create/update is itself a consequential external effect.
9. A previous destination realization may be reused only while its native representation still corresponds to current destination evidence; reuse does not create an MPC Address master.
10. Another business system may use a transaction address, contact/address object, other native mechanism or no supported mechanism; D4 standardizes meaning, not ERP topology.

This is intentionally not a new D1 domain.

### Current Sankhya evidence/direction

Historical/manual Mercado Livre documents measured `0/231` explicit transaction-destination representation in inspected native header fields. Operator evidence established these documents were manually created with the sale address placed on Partner registration. Therefore `0/231` is incumbent-process evidence only — not target authority and not proof of a Sankhya provider limit.

TOP 313 and other current e-commerce evidence nevertheless show a contact-based destination reference (`CODCONTATOENTREGA` + delivery city/UF) can exist and observed references propagated to fiscal results; `Contato` can carry an address distinct from the Partner master.

Target decisions:

- unattended Partner master-address mutation as variable-per-sale strategy — **REJECTED**;
- another Partner merely for another destination — **REJECTED**;
- contact-based alternate delivery — current strongest Sankhya candidate, **CONDITIONED / NOT YET CLAIMED** until D8 controlled proof establishes SourceInstance/fiscal/fulfillment consequences;
- a single mutable Partner delivery-address field is possible provider mechanism but not preferred for repeated variable destinations because it preserves shared mutable state;
- no safe realization → explicit Work / `external-required`.

The D8 proof must establish native party/master remain correct, destination survives native order progression/invoicing, fiscal/XML behavior is correct, authoritative reread reconstructs result and unrelated sales/customer state are not silently affected.

## 5.8 Business Order Intent materialization

Provider-independent contract:

```text
Business Order Intent
  + Organization
  + SourceInstance
  + canonical Marketplace Sale context
  + transaction-specific Selling Entity attribution
  + Product / quantity / Party Resolution / Destination Realization / required facts
→ native materialization attempt
→ source-qualified native business-order result(s)
→ authoritative convergence / rejection / pending / ambiguity
```

Provider-native intermediate artifacts/sequences remain adapter-local and never become configurable MPC workflow steps.

### Current Metal Nobre target binding

```text
create native e-commerce order
  target provider binding: TOP 313
  series: PA
  party: resolved CODPARC
  destination: only when safely supported for the sale
→ source-required progression
  current realization: CACSP.confirmarNota / A→L
→ authoritative reread
→ native business-order convergence
```

`CODTIPVENDA=27` is a required current fiscal/SourceInstance binding fact and supporting Mercado Livre negotiation/origin evidence, but **not the complete workflow selector**. It was measured across multiple TOP/TIPMOV topologies including historical/manual `14→303→305`. MPC selects the explicit SourceInstance target binding; it does not infer provider choreography from one negotiation code.

TOP 313 is e-commerce generally, not Mercado Livre identity. Current TOP-native reservation/financial effects are binding facts, not MPC policy.

For current creation, `CACSP.incluirNota` is the proven selected MGECOM surface. REST `POST /v1/vendas/pedidos` remains provider-conditioned and is not preferred merely for being REST. Measured partial-update behavior is not generalized into arbitrary patch semantics.

Only values actually required for correct execution are bound, sourced from stable SourceInstance configuration, domain-owned context/intent or externally governed/provider-derived prerequisite. Observed field variation does not justify a global knob bag.

## 5.9 Binding validation / hidden SourceInstance rules

Provider-declared configuration is not a complete description of SourceInstance behavior. TOP versions/configuration, custom triggers, liberações/approvals and procedural customizations may impose requirements absent from sanctioned configuration surfaces.

Binding validation:

1. detects observable provider/configuration drift;
2. is necessary but never predicts write success;
3. does not claim to prevalidate hidden/unexposed custom rules;
4. leaves execution-time uncertainty explicit;
5. cadence/cache mechanics remain D7.

## 5.10 Invoicing Intent / native correlation

Provider-independent contract:

```text
Invoicing Intent
  + native business-order result in required source state
  + Fulfillment-owned physical readiness/conference
  + required business-system/provider prerequisites
→ sanctioned native fiscal progression
→ source-qualified fiscal result
→ authoritative reread + origin/result correlation
→ converged / rejected / pending / ambiguous
```

Physical readiness remains Fulfillment authority.

Current production evidence establishes distinct native identities with `TOP 313 order → TOP 306 fiscal result` and line/quantity correlation. `SelecaoDocumentoSP.faturar` is the selected sanctioned progression surface. First controlled real `313→306` fiscal effect remains D8 because it is irreversible/legal and the architecture contract is already grounded.

Origin/result identities remain distinct source-qualified references; 0..N/partial results and line/quantity granularity remain representable where material; provider relation resources such as TGFVAR/`CompraVendavariosPedido` remain adapter evidence rather than MPC entities; transform 2xx is not convergence.

## 5.11 Pre-invoice reversal vs post-invoice fiscal consequence

Current SourceInstance history proves an alternative `313→307` branch for uninvoiced orders. This pre-invoice commercial reversal is **not** the same consequence as post-invoice fiscal return/reversal.

The sanctioned command and inventory-reservation fate of the 307 branch remain Unknown, and the post-invoice fiscal-return path remains unproven. Automated actuation stays `external-required`/deferred until the implicated claim requires proof.

## 5.12 Consequential business-system effect contract

Every admitted party/contact/address/order/progression/invoicing/reversal write obeys:

1. explicit Organization + SourceInstance target;
2. explicit owning-domain intent/correlation anchor;
3. prerequisites established no stronger than evidence permits;
4. request scoped only to intended effect;
5. response classified no stronger than provider evidence;
6. accepted/rejected/pending/ambiguous preserved where reachable;
7. authoritative reread/correlation after possible acceptance;
8. no blind retry after timeout/connection loss when acceptance is possible;
9. duplicate/ambiguity conditions become explicit Work;
10. provider/custom-trigger/liberação failures are translated rather than leaked as MPC business semantics;
11. protocol support never bypasses Readiness, Availability, Fulfillment, Governance or another owner;
12. provider PII is minimized.

## 5.13 Current Product 1.0 Sankhya proof lane

### Read/evidence lane

- production SourceInstance through sanctioned Gateway;
- Product identity/identifier evidence;
- qualified company/location inventory evidence;
- dedicated net stock plus bounded decomposition when material;
- Cost Observations;
- Expected Tax through the stable model binding in §5.4.

### Marketplace business-order lane

```text
Marketplace Sale
→ Party Resolution
→ Destination Realization or explicit Work
→ explicit target e-commerce binding
→ native TOP-313 order
→ source-required progression
→ authoritative reread
→ native business-order convergence
```

### Fiscal lane

```text
native order in required state
+ Fulfillment readiness
→ Invoicing Intent
→ sanctioned faturamento
→ native TOP-306 result
→ authoritative reread + correlation
```

## 5.14 Closure gates / safe defers

- **G1 Expected Tax — PASS / CLOSED FOR B3.**
- **G2 Party Resolution + Destination Realization — PASS WITH MATERIAL AMENDMENT / CLOSED FOR B3.**
- **G3 first selected-lane fiscal effect — DEFER SAFELY → D8.**
- **G4 controlled-product marketplace lane — DEFER SAFELY** until satisfiability/interchangeability/selection timing is established.
- **G5 post-invoice fiscal return — DEFER SAFELY / EXTERNAL-REQUIRED.**
- **G6 pre-invoice reversal reservation fate — DEFER WITH REVERSAL CLAIM.**

Fail-honest fiscal coverage Unknowns, not B3 closure gates:

- out-of-state PJ-contributor branch not specifically exercised;
- IPI/retained ICMS-ST component visibility unestablished on measured marketplace lane;
- accessory-expense request shape unestablished.

If a claimed flow materially requires one, prove it before claiming that branch or return explicit Unknown/unsupported/Work.

D7 defers token refresh, acquisition cadence, checkpoints/cursors, rate/concurrency/backoff, binding-validation cadence/cache and the mechanism preventing concurrent duplicate native-party creation.

**No B3 closure gate remains.**

## 5.15 B3 YAGNI / replacement fence

B3 MUST NOT introduce:

- generic ERP business entity/universal ERP ontology;
- universal `ERPAdapter` containing every possible operation;
- generic Provider/Resource/Capability graph;
- generic Customer/Party/Address/CRM master lifecycle for integration symmetry;
- universal party/address matching engine;
- plugin registry/factory for speculative providers;
- arbitrary workflow/materialization DSL;
- MPC TOP/NUNOTA/TGFVAR/CONTROLE/Sankhya-status entities;
- universal Lot/Batch/Serial model derived from `CONTROLE`;
- Sankhya product/stock/cost/tax/customer mirror;
- arbitrary SQL/DbExplorer through Gateway criteria;
- duplicated Sankhya tax engine;
- support for every historical Metal Nobre process;
- speculative second-ERP adapters;
- blanket marketplace→ERP customer-master synchronization;
- automatic Partner-address mutation as a substitute for explicit destination semantics.

Provider-independent replacement test preserves Product facts, inventory evidence, Cost Observations, fiscal-evidence requirements, Party Resolution, Destination Realization, Business Order/Invoicing Intents, native result convergence/correlation, honest unknown/unavailable and no-blind-retry semantics. Sankhya-local `Parceiro`, `CODPARC`, `Contato`, TOP, NUNOTA, negotiation type, notaModelo, CACSP/SelecaoDocumentoSP, TGFVAR, `CONTROLE`, auth/triggers/liberações and binding values are replaceable realization.

If a future real system reveals genuinely different business meaning, reopen the responsible semantic decision rather than contorting it into false commonality.

## 5.16 B3 reopen triggers

Reopen only the implicated decision when material evidence shows, for example:

- a second real business system cannot implement accepted consumer meaning without Sankhya semantics leaking into the domain;
- Party Resolution grows into independent customer/party lifecycle authority;
- Destination Realization grows into independent customer/address master authority;
- a Product 1.0 normal path requires destination behavior that cannot be safely represented and explicit external-required handling is no longer acceptable;
- a selected sanctioned Sankhya operation disappears/changes materially and no sanctioned replacement is sufficient;
- Expected-Tax model `898307` loses material binding prerequisites, negotiation type `27` becomes inactive/materially changed, the model-Partner authorization becomes unusable, or another hidden guard makes the transient calculation movement unsatisfiable;
- a claimed fiscal branch requires a material component the current calculation surface cannot represent honestly;
- a controlled Product enters marketplace scope and sellability/selection cannot fit accepted Availability/Fulfillment ownership;
- automated return/reversal becomes a claimed normal path while external-required handling is insufficient;
- a required fact/effect would only be satisfiable by arbitrary SQL/Direct Oracle;
- second-provider repetition proves a smaller shared technical mechanism materially reduces total complexity.

Do not reopen for naming preference, abstract symmetry or hypothetical providers.

---

# 6. Evidence basis and proof disposition

B1 external facts were reverified against current official provider/source documentation families during review. B2 current provider facts were independently challenged against current Mercado Livre documentation and a real read-only Installation probe. B3 combined official Sankhya/provider evidence with sanctioned sandbox/production measurements and independent Fable challenge; final Expected-Tax closure was measured on the operator-configured stable target model on 2026-08-18.

External sources and reviewer notes are evidence, never second architecture authority. Git history is the archive for detailed probe/review narratives.

D8 still owns the first controlled real effects whose semantic contracts are already accepted:

- first Mercado Livre Price/Availability effect + authoritative convergence reread;
- selected-lane fiscal/invoice/label progression;
- first irreversible Sankhya `313→306` fiscal progression;
- first consequential native party create/update when needed;
- first controlled alternate-destination/contact realization before claiming that concrete capability;
- any currently unexercised fiscal branch/component that becomes material to a selected golden flow.

A failed later proof narrows/reopens only the implicated concrete capability; it never authorizes Oracle fallback, master-data corruption or fabricated known values.

---

# 7. Review and authority protocol

D4 uses the accelerated protocol established by earlier stages:

1. GPT prepares a coherent batch from repository authority and claim-specific evidence.
2. Operator approves the batch direction for independent challenge.
3. A disposable `D4-B<n>-REVIEW-CANDIDATE.md` may be created and is explicitly non-authoritative.
4. Operator invokes Fable separately; Fable reconstructs authority independently and records evidence in `AI-DIALOG.md`.
5. Reviewer findings remain evidence; GPT independently adjudicates them.
6. Only material disagreement receives further reviewer convergence.
7. Operator explicitly ratifies the converged batch before canonical consolidation.
8. Accepted/canonical meaning is filed here and the disposable candidate is deleted.
9. A batch remains OPEN only when an explicit evidence/proof gate is required for whole-batch closure; the router names that gate.
10. At D4 closure, run final Global Coherence + YAGNI / Overengineering / Future-Cost review before whole-stage ratification.

`AI-DIALOG.md`, review candidates, chat summaries and reviewer statements are never target authority.

---

# 8. Current D4 state / exact next action

D4 is **OPEN / ACTIVE**.

- **D4-B1 — External Contract Grounding: ACCEPTED / CANONICAL.**
- **D4-B2 — Mercado Livre Operational Contract: ACCEPTED / CANONICAL; Installation Evidence Gate CLOSED / PASS.**
- **D4-B3 — Sankhya Business-System Contract: ACCEPTED / CANONICAL.**
- **D4-B4 — Market / Economics / Settlement Contract: NEXT / NOT YET OPENED.**
- **Final D4 Global Coherence + YAGNI / Overengineering / Future-Cost review: NOT STARTED.**

No D0/D1/D2/D3 or D4-B1/B2 reopen is currently required.

Exact next action: **open D4-B4 — Market / Economics / Settlement Contract from canonical D4-B1+B2+B3.** B4 owns concrete external market/competitor, fee, provider financial movement, refund/adjustment, settlement/payout and related provenance/coverage/reconciliation contracts required by accepted Product 1.0. Commercial Economics remains authority for economic interpretation, attribution and reconciliation.

B2 R1 and B2 R2/B3 real-effect obligations remain D8 proofs; they do not reopen accepted B2/B3 by themselves.

Do not begin D5 until B4 and the final D4 Global Coherence review are accepted. Product implementation remains blocked until D9.
