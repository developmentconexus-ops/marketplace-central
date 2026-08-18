# D4-B3 — Sankhya Business-System Contract — REVIEW CANDIDATE

> **Status:** REVIEW CANDIDATE / NON-AUTHORITATIVE / DISPOSABLE  
> **Program:** Architecture Rebaseline / Technical System Design  
> **Parent authority:** accepted D0–D4-B2 only  
> **Method:** DevelopmentConexus Engineering Method v1.0.0  
> **Base evidence HEAD:** `eaab7127518002949ebdfa00aead90151a85ec56`  
> **Independent review HEAD:** `f7ec08d91108ed905133874bb5bcc26f1b729b2b`  
> **Party / Destination review evidence:** `72658432d4aa123298bff96faa91e0c49656fb26`  
> **Expected-Tax closure evidence:** current `AI-DIALOG.md` blob `12d8c8bcd044590b5f4cdb603f21261b7e2848d1`, 2026-08-18  
> **Purpose:** final coherent disposable surface for operator adjudication of D4-B3. This file is not target authority and MUST be deleted before canonical consolidation.

---

## 1. Review question and scope

D4-B3 answers:

> **Through which concrete sanctioned business-system contracts can MPC obtain the internal facts and cause/reconcile the native business-order and fiscal effects required by Product 1.0, while preserving D0–D3 authority, explicit SourceInstance qualification, honest coverage and a future business-system seam without building a generic ERP/workflow/customer framework?**

The target is:

- provider-independent at the MPC semantic boundary;
- concrete about Sankhya inside D4 because Sankhya is the first business system proving the boundary;
- SourceInstance-aware because customer configuration/customization materially changes provider behavior;
- fail-honest where the sanctioned surface or current SourceInstance cannot satisfy a claim;
- bounded by YAGNI: no universal ERP ontology, generic Customer/Party/Address master, plugin framework or workflow/BPM engine.

A future TOTVS/Bling/SAP-like system is a structural replacement test only. It is not modelled now.

Implementation remains blocked until D9.

---

## 2. Imported authority — not reopened

B3 imports rather than re-decides:

1. Business-System Materialization owns `Business Order Intent`, `Invoicing Intent`, native-result correlation and convergence meaning.
2. ERP-native TOP/document taxonomy is not MPC semantics.
3. Product master remains external; MPC references source-qualified Product identity.
4. Availability Control owns Inventory Source/Scope, allocation meaning and Sellable Availability; native stock/reservation truth remains external.
5. Commercial Economics owns Cost Basis and economic interpretation; provider cost/tax variants are evidence only.
6. Marketplace Sales owns canonical marketplace-sale interpretation and transaction-specific Selling Entity attribution.
7. `SourceInstance` identifies one logical externally authoritative business-system/source namespace; credentials are not identity.
8. Consumer owns meaning; adapter owns protocol.
9. Known / absent / unknown / unavailable remain distinct; partial is not closure.
10. Acceptance is not convergence; ambiguous writes are not blindly retried.
11. Provider PII is minimized.
12. Direct Oracle/database access is outside target architecture and is not fallback.
13. Integration Support, Provider Effective Capability/Requirement and Effective Business Capability remain separate authorities.
14. Correspondence does not collapse external and MPC authorities into one identity; material contradiction fails closed.

No D0/D1/D2/D3 or D4-B1/B2 reopen is proposed.

---

## 3. Governing invariant

> **For each Product 1.0 business-system fact or effect, MPC depends on a consumer-owned semantic contract whose external evidence/effect is Organization + SourceInstance qualified. The concrete Sankhya adapter uses only sanctioned provider operations and bounded SourceInstance-specific bindings, preserves provider-native granularity/provenance/partiality, and reconciles consequential effects through authoritative reread/correlation. Provider-declared configuration may be validated where observable but never guarantees execution success; hidden/custom rules remain execution-time uncertainty. Sankhya-native topology, party vocabulary and customer configuration never become MPC business ontology, and future substitutability does not authorize a generic ERP/workflow/party framework before a second real consumer proves the abstraction.**

Corollaries:

- Sankhya first does not mean Sankhya model.
- SourceInstance binding does not mean workflow engine.
- Provider stock decomposition does not imply an MPC `Lot`/Batch entity or proven interchangeability.
- TOP/code identity alone does not prove current configured meaning.
- A validated binding is a precondition, never a prediction that an external effect will succeed.
- Marketplace buyer/account, fiscal/billing party and delivery recipient/destination remain distinct unless transaction evidence establishes their relationship.
- Resolving a fiscal/native party does not prove that the current delivery destination is safely representable.
- Transaction delivery evidence does not become a customer-master update command merely because the business system needs a native representation.

---

## 4. Alternatives / Global Maximum

### A — Sankhya-shaped MPC core

Examples: `Business Order = TOP 313`, `Invoice = TOP 306`, `Marketplace = CODTIPVENDA 27`, `Customer = CODPARC`, `Lot = CONTROLE`.

**REJECT.** One SourceInstance already contains multiple native commercial topologies, and fiscal-party/delivery meaning does not collapse safely into one provider master shape.

### B — generic ERP/workflow/customer framework now

Examples: `GenericERP`, `UniversalParty`, generic resource/capability graph, configurable materialization choreography, universal matching engine.

**REJECT.** There is one concrete business system and no second real consumer proving a stable shared ERP/party ontology or workflow language.

### C — consumer-owned semantics + concrete Sankhya adapter + bounded SourceInstance binding

**GLOBAL MAXIMUM.**

```text
D1 owner meaning
    ↓ consumer-owned semantic port
D4 concrete business-system adapter
    - sanctioned operations
    - bounded SourceInstance binding
    - provider capability/requirements
    - authoritative reread/correlation
    ↓
current Sankhya SourceInstance
```

**Falsifier:** if MPC/domain configuration must express arbitrary provider-native choreography or provider-native customer/master concepts merely so another ERP can fit later, this boundary has failed.

---

## 5. SourceInstance, environment and authentication

Production and sandbox are distinct external environments/namespaces.

Measured auth evidence includes `ambiente=hml` in sandbox and `ambiente=prd` in production. A sandbox `environment` UUID was not universally exposed in production and therefore is not a universal namespace law.

Current target auth:

```text
OAuth 2.0 client_credentials
+ X-Token
+ POST /authenticate
→ Bearer token
```

Measured token TTL is 300 seconds. Refresh scheduling/locking/caching belongs to D7.

MGECOM Bearer compatibility is established in sandbox and production.

---

## 6. Bounded sanctioned entity-read contract

Dedicated REST resources remain preferred when semantically sufficient. `CRUDServiceProvider.loadRecords` is admitted only when a real consumer requires a fact unavailable or materially lossy on a dedicated resource.

Admitted shape:

1. named sanctioned `rootEntity`;
2. explicit minimum result fieldset;
3. `criteria` only over fields of the named root entity using ordinary comparison/logical operators and bound parameters;
4. related data only through provider-declared relation/path mechanisms when genuinely required;
5. no subqueries, cross-table criteria references, Oracle-specific pseudo-tables/functions, arbitrary SQL expressions or query-language escape hatches;
6. needing such an expression is a capability finding, not authorization for Oracle-via-HTTP;
7. no entity read becomes an ERP mirror by convenience.

D7/implementation later chooses mechanical enforcement.

---

## 7. Product / Readiness evidence

Current external Product reference:

```text
SourceInstance + CODPROD
```

`CODPROD` remains provider-native, not an MPC Product-master identity.

Available provider observations include active state, reference, supplier reference, brand, NCM/fiscal facts and alternate-volume barcode where populated.

No universal first-class GTIN/EAN Product field was established. `REFERENCIA == EAN` is forbidden as an identity/correspondence law.

D4 supplies identifier evidence; Product & Channel Readiness owns sufficiency of correspondence.

---

## 8. Company, location and inventory contract

### 8.1 External references stay external

`CODEMP` and `CODLOCAL` do not automatically equal Selling Entity, Inventory Source or Fulfillment Node.

Transaction-specific Selling Entity attribution comes from Marketplace Sales; D4 does not infer it from company codes.

### 8.2 Native stock observations

Measured in sandbox and production:

```text
REST estoque = native ESTOQUE - RESERVADO
```

The REST net surface is useful but lossy because it cannot explain reservation decomposition. When a consumer needs that distinction, bounded sanctioned `Estoque` reads preserve the required provider dimensions/fields.

A negative net value is a real external observation and is not clamped to zero. Resource absence is not automatically `stock=0`.

### 8.3 Provider control dimension

`CONTROLE` is preserved as an opaque provider inventory-partition value when present. It is not canonically named `Lot`, `Batch`, `Tonality`, `Serial` or another MPC identity.

Current evidence shows control is Product-specific and no sanctioned structured tonalidade/calibre/grade source was established.

### 8.4 Sellability invariant

> **Sellable Availability may treat an aggregate quantity as sellable only when available evidence and applicable rules establish that the requested quantity is actually satisfiable by eligible source inventory. Where an external source decomposes inventory in a way material to satisfiability, D4 preserves that decomposition rather than erasing it into a pre-aggregated total.**

Current observed Mercado Livre sold population is control-free. Automatic controlled-product marketplace Availability/Materialization remains unclaimed until satisfiability/interchangeability and selection timing are proven.

### 8.5 Materialization-created inventory commitment

> **When native Business Order materialization creates or changes an external inventory commitment, the inventory observation contract preserves that effect so Availability Control can account for current sellability. Materialization does not acquire allocation authority, and MPC does not model its own materializations as inventory-neutral.**

---

## 9. Cost and Expected Tax

### 9.1 Cost

Bounded sanctioned `rootEntity=Custo` reads expose provider Cost Observations with company/Product/time/provider-local qualifiers.

Native variants such as `CUSGER`, `CUSREP`, `CUSSEMICM` and `CUSMED` remain observations, never Cost Basis by inheritance. Commercial Economics owns Cost Basis selection/interpretation.

### 9.2 Expected Tax — G1 PASS

`POST /v1/fiscal/impostos/calculo` is the selected sanctioned Expected-Tax calculation surface. MPC does not copy Sankhya fiscal rules into its own engine.

The stable SourceInstance binding proven on 2026-08-18 is:

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

The model-attached partner exists to satisfy Sankhya authorization/configuration for negotiation type `27`; it is **not** the fiscal customer of each simulation. The API requires `codigoCliente` explicitly, fails closed when it is missing, and a requested SP client remained SP for fiscal treatment even though the model partner is MG. Therefore the model partner does not silently substitute for or override transaction customer meaning.

Proof pair on the designated model, without company/TOP overrides:

- MG client: ICMS CST 60 / 0; PIS 1.65%; COFINS 7.6%; no DIFAL;
- SP client: ICMS CST 00 / 12% / 20.40; DIFAL destination 10.20; PIS 1.65%; COFINS 7.6%;
- results were byte-for-byte equal to the earlier target-lane pair using a real type-27 TOP-313 order and matched realized ICMS documents to the cent;
- zero persisted negative-NUNOTA headers/items/financial rows, stock/reservation remained unchanged, and the model itself was unchanged by the probes.

Therefore:

> **G1 = PASS. The sanctioned Sankhya surface is materially sufficient for the currently claimed ex-ante Expected-Tax evidence on the designated target binding, with zero persisted business effect.**

#### Expected-Tax adapter obligations

These are contract facts, not capability gaps:

1. **F-G1-1 — returned `valorBase` on PIS/COFINS/CSLL is not reliably the arithmetic base used by the engine.** MPC preserves provider-returned component/value/provenance and does not reconstruct attributable tax as `base × rate` where the provider evidence disproves that identity.
2. **F-G1-2 — unknown request fields may be silently ignored.** The adapter uses a typed/pinned request contract and must not treat HTTP 200 as proof that intended monetary inputs were consumed. Material echoed inputs such as quantity/unit price/discount are validated against what was sent where the surface exposes them.
3. **F-G1-3 — top-level `despesasAcessorias` is an object of unestablished request shape.** No arbitrary guessed shape is admitted; the capability stays Unknown until a real consumer requires and proves it.
4. **F-G1-4 — IPI and ICMS-ST components were not established as separately observable through this response.** Absence is not proof of non-applicability. This was not material to the measured marketplace proof lane, whose realized IPI was zero, but a future claim requiring those components must prove them or remain Unknown/unsupported.
5. **F-G1-6 — `CODTIPVENDA` is fiscally determinant.** Binding must preserve type `27`; ICMS agreement alone is not equivalence evidence for a substitute negotiation type because a wrong type was measured to silently alter PIS while still returning HTTP 200.

PF and PJ customers both exist in the measured e-commerce/marketplace evidence. The specific branch **out-of-state PJ contribuinte** remains unexercised. That is an explicit Unknown, not a new B3 closure gate: MPC delegates fiscal-rule authority to Sankhya and may claim only the component/context coverage actually established. A golden flow encountering an unproven materially different fiscal branch must validate it or fail honestly rather than infer equivalence.

Expected-Tax binding drift must be detectable no stronger than sanctioned evidence permits. Revalidation includes material model/type/partner/source configuration needed by the binding; cadence/cache mechanics remain D7.

Historical realized tax remains realized evidence/B4 input, not a replacement L0 engine.

---

## 10. Coverage, pagination and change semantics

Full/scoped acquisition is the correctness baseline. Point reads prove only their exact source-qualified scope; enumeration proves only the completed traversed scope.

Delta/change-log use is optional and prerequisite-bound. No delta result is accepted as completeness evidence until source change-log coverage is independently established. D7 owns cadence, checkpoints and recovery mechanics.

Provider-specific page origin, `hasMore`, short-page, 404-as-end, ordering and date-format quirks remain adapter-local. A provider option being accepted syntactically does not prove it was honored semantically.

---

## 11. Native order observation and precedence

Native business-order identity remains:

```text
SourceInstance + native NUNOTA
```

REST order enumeration is useful for bounded discovery/observation but documented point filters were empirically unreliable. Bounded sanctioned `CabecalhoNota` read by NUNOTA is the current authoritative point/reread surface for consequential native-order state.

When enumeration and authoritative point reread disagree, point reread governs current consequential state and divergence remains explicit evidence.

Raw `STATUSNOTA`/`PENDENTE` remain adapter-local. MPC does not acquire a universal `Confirmation` lifecycle stage.

Provider-independent meaning:

> **the native business order has reached the externally required state sufficient for the next claimed materialization progression.**

For the current Sankhya binding, `CACSP.confirmarNota` and `A → L` are adapter-local realization.

---

## 12. Business-System Party Resolution and Destination Realization

### 12.1 Party Resolution — provider-independent contract

```text
Marketplace Sale
  ├─ marketplace buyer/account evidence
  └─ fiscal/billing-party evidence
             ↓
Business-System Party Resolution
             ↓
source-native business-party reference
             ↓
Business Order Materialization
```

Rules:

1. Marketplace Sales owns sale/buyer interpretation and supplies only facts legitimately required for materialization.
2. Marketplace buyer/account identity is provenance/context, not automatically fiscal party.
3. A prior explicit native resolution may be reused while materially compatible with current fiscal-party evidence and native party state; a different delivery destination does not by itself invalidate the fiscal-party correspondence.
4. Exactly one sufficiently established compatible native match may be used.
5. Zero native matches may authorize native creation only when every required identity-bearing fact is known from legitimate transaction evidence; creation does not silently decide how a different shipping destination becomes customer master data.
6. Multiple matches or material fiscal/identity contradictions are `AMBIGUOUS`: no guessed selection, first-result-wins or new duplicate creation.
7. Native party create/update is a consequential external effect under §18.
8. MPC may preserve bounded resolution/correspondence lineage when needed to avoid repeating a human-adjudicated ambiguity, without owning customer attributes/lifecycle.
9. Concurrent/repeated materializations for the same unresolved fiscal identity must not independently create duplicate native parties; D7 chooses enforcement mechanism.

Durable Party Resolution state is **KEEP — PRESENT CORRECTNESS NEED** because rereading a source that still contains the same seven-way ambiguity cannot reconstruct the human adjudication.

### 12.2 Destination Realization — provider-independent contract

```text
Marketplace Sale
  └─ delivery-recipient/destination evidence
             ↓
Business-System Destination Realization
  + resolved native business-party reference when required
             ↓
source-supported destination representation
        OR explicit Work / external-required
             ↓
Business Order Materialization
```

Rules:

1. Fiscal/billing party is not automatically delivery recipient.
2. Delivery destination is transaction evidence and is not customer-master address by inheritance.
3. Party Resolution success does not imply destination representability.
4. If the business system can represent the destination without destructive master mutation, Materialization uses the least-destructive sufficiently proven native realization.
5. If no safe native realization is proven, destination becomes explicit `external-required` / Work; it is not dropped or fabricated as master equivalence.
6. MPC never creates another native customer merely to hold another delivery address.
7. MPC never performs an unattended overwrite of registered/master customer address merely because a marketplace transaction carries a different shipping address.
8. Any native contact/address/destination create/update carrying identity/address meaning is a consequential external effect.
9. A prior destination realization may be reused only while it still corresponds to the current sale's destination evidence; reuse does not create an MPC Address master.
10. A future system may realize destination directly on a transaction, through a contact/address object, through another native mechanism, or not at all. D4 standardizes meaning, not one ERP topology.

Party Resolution and Destination Realization are bounded prerequisites of Business-System Materialization, not new D1 domains.

### 12.3 Current Sankhya Party Resolution realization

- source-native reference = `CODPARC`;
- bounded sanctioned `Parceiro` entity reads provide lookup evidence when dedicated customer resources are insufficient;
- legal document `CGC_CPF` is a conditioned lookup signal, not guaranteed unique identity;
- lookup uses ordinary equality with a string-bound parameter;
- production evidence contains a real seven-way duplicate for one legal document; ERP uniqueness therefore cannot be delegated the D7 correctness property;
- current e-commerce origin custom fields are not reliable correspondence authority;
- PF/PJ material differences exist and remain provider-local evidence (`TIPPESSOA`, `CLASSIFICMS`, applicable IE state);
- `Parceiro`, `CODPARC`, `CGC_CPF` and protocol details are Sankhya realization only.

### 12.4 Current Sankhya Destination Realization direction

Corrected measurement established `0/231` explicit transaction destination fields across the observed historical/manual Mercado Livre-tagged document population. This is **incumbent-process evidence only**: the operator established that those sales were manually created and the sale address was placed on Partner registration. It proves neither the target realization nor a provider limitation.

Separate provider evidence proves:

- Partner has registered/master address state;
- TOP 313 can carry `CODCONTATOENTREGA` + delivery city/UF in other e-commerce cases;
- observed contact-based destination references propagated from TOP-313 origins into correlated fiscal results without loss;
- `Contato` can carry address distinct from Partner master.

Target decision:

- unattended Partner master-address mutation as variable-per-sale strategy — **REJECTED**;
- another Partner merely for another destination — **REJECTED**;
- contact-based delivery realization — current strongest Sankhya candidate, **CONDITIONED / NOT YET CLAIMED** until controlled proof establishes SourceInstance configuration and fiscal/fulfillment consequences;
- single mutable Partner delivery-address field — possible provider mechanism but not preferred for repeated variable destinations because it remains shared mutable state;
- no safe realization — explicit `external-required` / Work.

The first alternate-destination/contact path remains D8 proof. A failed proof narrows Sankhya capability; it does not rewrite the provider-independent contract or authorize master overwrite.

### 12.5 G2 status

> **G2 = PASS WITH MATERIAL AMENDMENT at architecture-contract level.**

Closed meaning:

- Party Resolution and Destination Realization remain distinct;
- durable ambiguous Party Resolution is kept without Customer Master authority;
- unattended master overwrite and duplicate-customer-per-address are rejected;
- unsupported destination realization fails honestly;
- concurrent duplicate native-party creation is a D7 correctness obligation.

---

## 13. Business Order Intent materialization

### 13.1 Provider-independent contract

```text
Business Order Intent
  + Organization
  + SourceInstance
  + canonical Marketplace Sale context
  + transaction-specific Selling Entity attribution
  + required Product / quantity / Party Resolution / Destination Realization / materialization facts
→ attempt native materialization
→ source-qualified native business-order result(s)
→ authoritative convergence / rejection / pending / ambiguity
```

The domain does not own TOPs, series, TIPMOV, CACSP, NUNOTA letters, TGFVAR or provider-native customer/address master.

Provider-native intermediate artifacts may exist internally. Their number/sequence are adapter concerns and do not become configurable MPC workflow steps.

### 13.2 Current Metal Nobre target binding

```text
create native e-commerce order
  target binding: TOP 313
  series: PA
  negotiation type 27 as required SourceInstance binding fact
  business-party reference: resolved CODPARC
  destination realization: selected only when safely supported
→ source-required progression
  Sankhya realization: CACSP.confirmarNota / A→L
→ authoritative reread
→ native business-order converged
```

`CODTIPVENDA=27` is **not** a sufficient workflow selector. It was observed across multiple TOP/TIPMOV combinations, including historical/manual `14→303→305`. MPC explicitly selects its SourceInstance target binding; it does not infer native choreography from one negotiation code.

For Expected Tax specifically, type `27` is also a material fiscal binding fact (§9.2). “Not a workflow selector” and “fiscally determinant” are distinct statements.

TOP 313 is e-commerce generally, not Mercado Livre identity. Its provider-native reservation/financial effects are binding facts, not MPC policy.

### 13.3 Creation surface

For the selected binding, `CACSP.incluirNota` is the proven MGECOM creation surface. REST `POST /v1/vendas/pedidos` remains provider-conditioned and is not preferred merely because it is REST.

Measured partial-update behavior of `CACSP.incluirNota` is not generalized into an arbitrary patch API.

### 13.4 Required input binding — not a knob bag

Only values actually required for correct execution are bound, sourced explicitly from stable SourceInstance configuration, current domain-owned intent/context or externally governed/provider-derived prerequisites.

Observed variation in vendor, nature, carrier, partner, delivery or other provider fields does not justify one global setting per field.

---

## 14. Binding validation and hidden SourceInstance rules

`TipoOperacao` is version-qualified by `DHALTER`; `CODTIPOPER` alone is not eternal meaning.

Provider-declared configuration is not a complete description of SourceInstance behavior. Custom triggers, liberação/approval rules and procedural customizations may impose requirements absent from sanctioned configuration surfaces.

Binding validation therefore:

1. detects observable configuration drift;
2. is necessary but never sufficient for a consequential effect;
3. does not claim to pre-validate all hidden provider/custom prerequisites;
4. leaves unexposed rules as execution-time uncertainty;
5. never turns successful validation into a prediction of write success.

Expected-Tax binding validation additionally preserves the proven model/type relationship and must detect material drift such as model `898307` losing required model/type/partner configuration, negotiation type `27` becoming unusable, or the bound partner no longer satisfying the provider constraint. Exact checks/cadence remain D7/implementation mechanics.

---

## 15. Invoicing Intent materialization

```text
Invoicing Intent
  + native business-order result in required source state
  + Fulfillment-owned physical readiness/conference meaning
  + required provider/business-system prerequisites
→ sanctioned native fiscal progression
→ source-qualified fiscal result
→ authoritative reread + origin/result correlation
→ converged / rejected / pending / ambiguous
```

Physical readiness remains Fulfillment authority.

Production history proves `TOP 313 → TOP 306` with distinct native identities and line/quantity correlation.

`SelecaoDocumentoSP.faturar` is the selected sanctioned progression surface. The first controlled real 313→306 fiscal write remains D8 proof because it is an irreversible/legal effect whose architecture contract is already grounded.

---

## 16. Native result correlation

D4 preserves source-native correlation without creating a generic provider graph.

Requirements:

- origin and result identities remain distinct source-qualified references;
- 0..N/partial results remain representable when exposed;
- line/quantity granularity is preserved when material;
- provider relation resources remain adapter evidence, not MPC entities;
- transform response/2xx is not convergence.

Current Sankhya realization may use TGFVAR/`CompraVendavariosPedido`; business domains do not read provider relation resources directly.

---

## 17. Pre-invoice commercial reversal vs post-invoice fiscal return

History proves a branch in which uninvoiced TOP-313 orders correlate to TOP-307 results.

Known: observed 307s originate from uninvoiced 313 orders, reverse commercial/financial posture, have no observed stock write-down and are not observed as results of TOP-306 invoices.

> **Pre-invoice native commercial reversal is not the same business-system consequence as post-invoice fiscal return/reversal.**

Unknown/deferred: sanctioned write command for the 307-class consequence, reservation fate, and post-invoice fiscal-return path.

These do not block B3 while automated reversal remains explicit `external-required`; they become gates before that actuation is claimed.

---

## 18. External-effect semantics

Every consequential business-system write admitted by the selected contract — including party/contact/address create/update, order create/update, source-required order progression, invoicing and later reversal effects — obeys:

1. explicit Organization + SourceInstance target;
2. explicit owning-domain intent/correlation anchor;
3. prerequisites established no stronger than evidence permits;
4. request scoped only to intended effect;
5. response classified no stronger than provider evidence;
6. accepted/rejected/pending/ambiguous preserved where reachable;
7. authoritative reread/correlation after possible acceptance;
8. no blind retry after timeout/connection loss when acceptance is possible;
9. duplicate/ambiguity conditions become explicit work;
10. provider/custom-trigger/liberação failures are translated rather than leaked as business semantics;
11. protocol support never bypasses Readiness, Availability, Fulfillment, Governance or another owner;
12. provider PII is minimized.

---

## 19. Current Product 1.0 proof lane

### Read / evidence lane

- production Sankhya SourceInstance through sanctioned Gateway;
- Product identity/identifier evidence;
- qualified company/location inventory evidence;
- dedicated net stock plus bounded decomposition when material;
- Cost Observations;
- Expected Tax through the proven stable model binding in §9.2.

### Marketplace business-order lane

```text
Marketplace Sale
→ Party Resolution
   resolve CODPARC or explicit Work
→ Destination Realization
   safe supported realization or explicit Work
→ explicit Metal Nobre target e-commerce binding
→ native TOP-313 order
→ source-required progression
→ authoritative reread
→ native business-order convergence
```

`CODTIPVENDA=27` is supporting SourceInstance evidence and a required target/fiscal binding fact; it is not the complete workflow selector. Historical/manual ML-tagged documents on other TOPs do not become target choreography by inheritance.

### Fiscal lane

```text
native business order in required source state
+ Fulfillment readiness
→ Invoicing Intent
→ sanctioned faturamento
→ native TOP-306 fiscal result
→ authoritative reread + correlation
```

First actual selected-lane fiscal effect remains D8.

---

## 20. Closure gates and safe defers

### G1 — Expected Tax — **PASS / CLOSED FOR B3**

The original seller-trigger blocker was remediated by the ERP owner. A second negotiation-type prerequisite was resolved without creating a substitute fiscal type: the operator attached an authorized partner to model `898307`, allowing the real negotiation type `27` to be carried by the stable model.

The final model pair proved destination-sensitive, componentized tax evidence, realized-result agreement, explicit transaction-client precedence and zero residue.

No G1 closure gate remains.

### G2 — Party Resolution + Destination Realization — **PASS WITH MATERIAL AMENDMENT / CLOSED FOR B3**

Provider-independent contract is closed at architecture level. First consequential native-party/contact write and alternate-destination realization remain D8 proofs before those concrete effect paths are claimed.

### G3 — First selected-lane fiscal effect — **DEFER SAFELY → D8**

Controlled actual 313→306 effect plus authoritative fiscal reread/correlation.

### G4 — Controlled-product marketplace lane — **DEFER SAFELY**

Before controlled-product automation is claimed, establish satisfiability/interchangeability and selection timing. No adapter-chosen `CONTROLE`, no aggregate guess.

### G5 — Post-invoice fiscal return — **DEFER SAFELY / EXTERNAL-REQUIRED**

Current history does not establish path/command.

### G6 — Pre-invoice reversal inventory-commitment fate — **DEFER WITH REVERSAL CLAIM**

Reservation fate remains Unknown and must close before automated reversal/convergence is claimed.

### Fiscal-coverage Unknowns — **FAIL-HONEST, NOT B3 CLOSURE GATES**

- out-of-state PJ contribuinte branch not specifically exercised;
- IPI / retained ICMS-ST component visibility not established on the measured marketplace lane;
- accessory-expense request shape unestablished.

These do not authorize plausible zeros or extrapolation. If a claimed flow materially requires one, prove it before claiming that branch or return explicit Unknown/unsupported/Work as appropriate.

### D7 defers

Token refresh; acquisition cadence; checkpoints/cursors; rate/concurrency/backoff; binding-validation cache/cadence; and mechanism preventing concurrent duplicate native-party creation for one unresolved fiscal identity.

> **B3 now has no known `CLOSURE GATE / OPEN` item.**

---

## 21. YAGNI / explicit non-goals

B3 MUST NOT introduce:

- generic ERP business entity or universal ERP ontology;
- universal `ERPAdapter` containing every possible operation;
- generic provider/resource/capability graph;
- generic `Customer`, `Party`, `Address` or CRM master lifecycle merely for integration symmetry;
- universal party/address matching engine;
- plugin registry/factory for speculative providers;
- arbitrary workflow/materialization DSL;
- MPC TOP/NUNOTA/TGFVAR/CONTROLE/Sankhya-status entities;
- universal Lot/Batch/Serial model derived from `CONTROLE`;
- Sankhya product/stock/cost/tax/customer mirror;
- arbitrary SQL/DbExplorer behavior through `loadRecords.criteria` or another Gateway route;
- duplicated Sankhya tax engine;
- support for every historical Metal Nobre process;
- speculative TOTVS/Bling adapters;
- family-level inference of control semantics;
- blanket marketplace→ERP customer-master synchronization;
- automatic Partner-address mutation as a substitute for explicit destination semantics.

Observed historical/manual `14→303→305` usage — including Mercado Livre-tagged documents — remains current-state variability evidence, not a Product 1.0 target workflow requirement.

---

## 22. Provider-independent replacement test

These meanings remain valid if a future accepted business system supplies the same semantics:

- source-qualified Product facts;
- inventory evidence sufficient for Sellable Availability derivation;
- Cost Observations;
- expected/realized fiscal evidence requirements;
- Party Resolution as bounded Materialization prerequisite;
- Destination Realization as distinct bounded Materialization prerequisite;
- buyer/account, fiscal/billing party and delivery evidence remaining distinct;
- safe native destination representation or explicit unsupported/Work without implicit master mutation;
- Business Order Intent;
- native business-order convergence on externally required state;
- Invoicing Intent;
- native fiscal convergence;
- origin/result line/quantity correlation;
- pre-invoice commercial reversal distinct from post-invoice fiscal consequence;
- honest unknown/unavailable;
- authoritative reread;
- no blind retry after ambiguous effect.

Replaceable Sankhya-local pieces include `Parceiro`, `CODPARC`, `CGC_CPF`, `Contato`, address/contact fields, TOP/versioning, NUNOTA, series, statuses, negotiation types, notaModelo, CACSP/SelecaoDocumentoSP, TGFVAR, `CONTROLE`, auth/Gateway, triggers/liberações and all Metal Nobre binding values.

If a future real system reveals genuinely different business meaning, reopen the responsible semantic decision rather than contorting it into a false common model.

Marketplace Installation and business-system SourceInstance remain distinct; no generic `IntegrationInstance` is justified.

---

## 23. Final proposed B3 outcome

### `CURRENT STRUCTURE CONFIRMED`

Provider-independent structure:

- consumer-owned semantic ports;
- SourceInstance-qualified external facts/results;
- concrete provider/business-system adapter;
- bounded Party Resolution under Materialization, not Customer Master authority;
- bounded Destination Realization under Materialization, not Address/CRM master authority;
- explicit external-required/Work for unsupported branches;
- no Integration business domain;
- no generic ERP/workflow/customer framework.

### `CURRENT STRUCTURE CONFIRMED — Sankhya target`

The sanctioned Gateway/API surface is sufficient to define the current Product 1.0 business-system contract without Direct Oracle fallback under the proven capability fences:

- Product / inventory / cost evidence can be acquired through admitted sanctioned reads;
- Expected Tax is proven through the stable model binding with real negotiation type 27;
- native business-order create/progression/reread/correlation is grounded;
- native invoicing progression/correlation is grounded for D8 controlled effect proof;
- party/destination semantics preserve master-data authority and fail honestly when a safe destination realization is not yet proven.

### `DEFER SAFELY`

- D7 runtime/concurrency mechanics;
- first irreversible selected-lane fiscal write → D8;
- first consequential native-party write → D8 when needed;
- first controlled alternate-destination/contact realization → D8 before claiming that capability;
- controlled-product marketplace automation;
- automated pre-invoice reversal until command/reservation-fate proof;
- post-invoice fiscal return unless selected golden flow requires automated actuation;
- unexercised fiscal branches/components until a real claimed flow requires them.

### Closure readiness

> **No known B3 closure gate remains. D4-B3 is READY FOR OPERATOR RATIFICATION AND CANONICAL CONSOLIDATION.**

This statement is candidate disposition only. It does not accept B3, update the router, authorize B4 or implementation.

---

## 24. Proof state

Already evidenced proportionately:

- OAuth/environment binding and short token TTL;
- Product native key/identifier limits;
- company/location/control stock granularity;
- REST net formula, negative-net behavior and reservation decomposition;
- controlled-product/source decomposition evidence;
- Cost Observations;
- change-log prerequisite/failure honesty;
- REST-order point-filter failure and bounded point-reread fallback;
- `CACSP.incluirNota` creation;
- `CACSP.confirmarNota` success/rejection + reread;
- MGECOM OAuth compatibility;
- `SelecaoDocumentoSP.faturar` capability;
- source-native origin/result line/quantity correlation;
- production e-commerce topology `313→306` and observed `313→307` branch;
- binding configuration/version reads plus hidden-rule counterevidence;
- current control-free marketplace sold lane;
- Party Resolution duplicate ambiguity and durable-resolution need;
- historical/manual ML `0/231` destination absence correctly classified as incumbent-process evidence, not target/provider limit;
- TOP 313 contact-based destination reference exists in other e-commerce cases and observed reference propagated into fiscal result;
- `CODTIPVENDA=27` spans multiple native topologies and cannot select complete workflow by itself;
- unattended customer-master address overwrite rejected as transaction delivery strategy;
- Expected-Tax stable model `898307` + type `27` + model partner authorization binding;
- explicit transaction-client precedence over the model partner;
- in-state/out-of-state destination-sensitive tax pair matching realized evidence;
- PIS distortion under wrong negotiation type proving type 27 is fiscally material;
- zero-residue Expected-Tax reread across headers/items/financial/stock/model;
- production-vs-sandbox materialization divergence.

No additional B3 proof is required before operator ratification.

D8/deferred proofs before specific concrete effect capabilities are claimed:

- first consequential native-party create/update;
- controlled Sankhya alternate-destination proof, including contact/configuration/fiscal/fulfillment consequences;
- first irreversible selected-lane fiscal effect;
- any unexercised fiscal branch/component that becomes material to a selected golden flow.

A failed later proof narrows the implicated concrete capability or triggers targeted reopen; it does not authorize Oracle fallback, master-data corruption or fabricated known values.

---

## 25. Reopen triggers

Reopen only the implicated decision when material evidence shows, for example:

- a second real business system cannot implement accepted consumer meaning without Sankhya semantics leaking into the domain;
- Party Resolution grows into independent customer/party business lifecycle authority;
- Destination Realization grows into independent customer/address master lifecycle authority;
- a Product 1.0 normal path requires destination behavior that cannot be represented safely and explicit external-required handling is no longer acceptable;
- selected Sankhya operation disappears/changes materially and no sanctioned replacement is sufficient;
- Expected-Tax model `898307` loses material binding prerequisites, including usable negotiation type `27` or the model partner relationship required by current provider authorization;
- negotiation type `27` becomes inactive/effectively unusable or changes materially enough to invalidate the proven fiscal binding;
- another hidden SourceInstance guard makes the transient tax movement unsatisfiable through the stable sanctioned binding;
- a claimed fiscal branch requires materially necessary components that the current calculation surface cannot represent honestly;
- a controlled Product enters marketplace scope and sellability/selection semantics cannot fit accepted Availability/Fulfillment ownership;
- automated post-invoice fiscal return becomes a claimed normal path and external-required handling is insufficient;
- Product 1.0 requires a fact/effect for which only arbitrary SQL/Direct Oracle could satisfy the claim;
- second-provider repetition proves a smaller shared technical mechanism materially reduces total complexity.

Do not reopen for naming preference, abstract symmetry or hypothetical providers.

---

## 26. Candidate disposition

This file is a disposable design/review surface only.

Its final adjudicated candidate position is:

> **D4-B3 — READY FOR OPERATOR RATIFICATION / CANONICAL CONSOLIDATION; no known closure gate remains; implementation remains blocked until D9.**

It MUST be deleted after operator ratification and before/with canonical D4-B3 consolidation.

Canonical stage/status remains solely in `docs/engineering/rebaseline/README.md`. This file does not itself open/accept B3, authorize B4/D5+, authorize implementation or merge.