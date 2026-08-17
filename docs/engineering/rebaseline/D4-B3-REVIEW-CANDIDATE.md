# D4-B3 — Sankhya Business-System Contract — REVIEW CANDIDATE

> **Status:** REVIEW CANDIDATE / NON-AUTHORITATIVE / DISPOSABLE  
> **Base evidence HEAD:** `eaab7127518002949ebdfa00aead90151a85ec56`  
> **Independent review HEAD:** `f7ec08d91108ed905133874bb5bcc26f1b729b2b`  
> **Latest evidence HEAD before this amendment:** `24ae547ae980ae95eb4c2b85ff0e90774fa2c52c`  
> **Program:** Architecture Rebaseline / Technical System Design  
> **Parent authority:** accepted D0–D4-B2 only  
> **Method:** DevelopmentConexus Engineering Method v1.0.0  
> **Purpose:** coherent disposable surface for closing D4-B3 evidence and review. This file is not target authority and MUST be deleted before canonical consolidation.

---

## 1. Review question and scope

D4-B3 must answer:

> **Through which concrete sanctioned business-system contracts can MPC obtain the internal facts and cause/reconcile the native business-order and fiscal effects required by Product 1.0, while preserving D0–D3 semantic authority, explicit SourceInstance qualification, honest coverage and a future business-system seam without building a generic ERP/workflow/customer framework?**

The target is intentionally:

- provider-independent at the MPC semantic boundary;
- concrete about Sankhya inside D4 because Sankhya is the first business system proving that boundary;
- SourceInstance-aware because real customer configuration/customization changes provider behavior;
- bounded by YAGNI: no universal ERP ontology, generic party/customer master, plugin framework or workflow/BPM engine.

A future TOTVS/Bling/SAP-like system is only a structural replacement test. It is not modelled now.

Implementation remains blocked until D9.

---

## 2. Imported authority — not reopened

B3 imports rather than re-decides:

1. Business-System Materialization owns `Business Order Intent`, `Invoicing Intent`, native-result correlation and convergence meaning.
2. ERP-native TOP/document taxonomy is not MPC semantics.
3. Product master remains external; MPC references source-qualified Product identity.
4. Availability Control owns Inventory Source/Scope, allocation meaning and Sellable Availability; native stock/reservation truth remains external.
5. Commercial Economics owns Cost Basis and economic interpretation; provider cost variants are evidence only.
6. Marketplace Sales owns canonical marketplace-sale interpretation, buyer/sale context and transaction-specific Selling Entity attribution.
7. Post-Sale Resolution coordinates consequences without importing provider-native fiscal taxonomy.
8. `SourceInstance` identifies one logical externally authoritative business-system/source namespace; credentials are not identity.
9. Consumer owns meaning; adapter owns protocol.
10. Known / absent / unknown / unavailable remain distinct; partial is not closure.
11. Acceptance is not convergence; ambiguous writes are not blindly retried.
12. Provider PII is minimized.
13. Direct Oracle/database access is outside target architecture and is not fallback.
14. Integration Support, Provider Effective Capability/Requirement and Effective Business Capability remain separate authorities.
15. Correspondence does not collapse external and MPC authorities into one identity; material contradictory evidence fails closed.

No D0/D1/D2/D3 or D4-B1/B2 reopen is proposed.

---

## 3. Governing invariant

> **For each Product 1.0 business-system fact or effect, MPC depends on a consumer-owned semantic contract whose external evidence/effect is Organization + SourceInstance qualified. The concrete Sankhya adapter uses only sanctioned provider operations and bounded SourceInstance-specific bindings, preserves provider-native granularity/provenance/partiality, and reconciles consequential effects through authoritative reread/correlation. Provider-declared configuration may be validated where observable but never guarantees execution success; hidden/custom rules remain execution-time uncertainty. Sankhya-native topology, party vocabulary and customer configuration never become MPC business ontology, and future substitutability does not authorize a generic ERP/workflow/party framework before a second real consumer proves the abstraction.**

Corollaries:

- Sankhya first does not mean Sankhya model.
- SourceInstance binding does not mean workflow engine.
- Provider stock decomposition does not imply an MPC `Lot`/Batch entity or proven interchangeability.
- TOP/code identity alone does not prove current configured meaning.
- Provider configuration evidence is not MPC policy.
- A validated binding is a necessary precondition, never a prediction that a consequential effect will succeed.
- A marketplace buyer account, fiscal/billing party, delivery recipient/destination and business-system-native party are distinct evidence/identity scopes unless the transaction evidence establishes their relationship.

---

## 4. Alternatives / Global Maximum

### A — Sankhya-shaped core

Examples: `Business Order = TOP 313`, `Invoice = TOP 306`, `Marketplace = CODTIPVENDA 27`, `Customer = CODPARC`, `Lot = CONTROLE`.

**REJECT.** One current SourceInstance already contains materially different commercial topologies, and business-party/delivery evidence does not collapse cleanly into one provider master record.

### B — generic ERP/workflow/customer framework now

Examples: `GenericERP`, `UniversalParty`, generic resource/capability graph, configurable materialization step sequence, universal matching engine.

**REJECT.** There is one concrete business system and no second real consumer proving a stable shared ERP/party ontology or workflow language.

### C — provider-independent consumer semantics + concrete Sankhya adapter + bounded SourceInstance binding

**PROPOSED GLOBAL MAXIMUM.**

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

**Falsifier:** if MPC/domain configuration must express arbitrary provider-native document choreography or provider-native customer/master concepts merely so a different ERP can fit later, the boundary has failed and must return to decision.

---

## 5. SourceInstance, environment and authentication

Production and sandbox are distinct environments/namespaces.

Measured auth evidence:

- sandbox token: `ambiente=hml`;
- production token: `ambiente=prd`;
- sandbox exposed an `environment` UUID while production exposed it as null, so that UUID is not a universal namespace law.

Current target auth:

```text
OAuth 2.0 client_credentials
+ X-Token
+ POST /authenticate
→ Bearer token
```

Measured token TTL is 300 seconds. Refresh scheduling/locking/caching belongs to D7.

MGECOM Bearer compatibility is empirically established in sandbox and production.

---

## 6. Bounded sanctioned entity-read contract

Dedicated REST resources remain preferred when semantically sufficient. `CRUDServiceProvider.loadRecords` is admitted only when a real consumer requires a fact unavailable or materially lossy on the dedicated resource.

The admission is intentionally narrower than what the Gateway parser may technically accept:

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

Measured in sandbox and independently re-measured in production:

```text
REST estoque = native ESTOQUE - RESERVADO
```

The REST net surface is useful but lossy because it cannot explain reservation decomposition. When a consumer needs that distinction, bounded sanctioned `Estoque` reads preserve the materially required provider dimensions/fields.

The net surface may legitimately be negative when commitments exceed physical stock. Negative is a real external observation; it is neither zero nor unknown and is not clamped away at the adapter boundary.

Resource absence is not automatically `stock=0`.

### 8.3 Provider control dimension

`CONTROLE` is preserved as an opaque provider inventory-partition value when present. It is not canonically named `Lot`, `Batch`, `Tonality`, `Serial` or another MPC identity.

Current evidence shows control is Product-specific, the same families contain controlled and uncontrolled products, one raw control-type code is observed in this SourceInstance, `CONTROLE` may contain physical-looking lot strings or `ENCOMENDA`, and no sanctioned structured tonalidade/calibre/grade source was established.

### 8.4 Provider-independent sellability invariant

> **Sellable Availability may treat an aggregate quantity as sellable only when the available evidence and applicable rules establish that the requested quantity is actually satisfiable by eligible source inventory. Where an external source decomposes inventory in a way material to satisfiability, D4 preserves that decomposition as evidence rather than erasing it into a pre-aggregated total.**

Current observed Mercado Livre sold population is control-free. Automatic controlled-product marketplace Availability/Materialization remains unclaimed until satisfiability/interchangeability and selection timing are proven.

### 8.5 Materialization-created inventory commitment

> **When native Business Order materialization creates or changes an external inventory commitment, the inventory observation contract preserves that effect so Availability Control can account for current sellability. Materialization does not acquire allocation authority, and MPC does not model its own materializations as inventory-neutral.**

This is not a new Materialization→Availability business-authority edge; Availability observes its authoritative external inventory source.

---

## 9. Cost and expected-tax evidence

### 9.1 Cost

Bounded sanctioned `rootEntity=Custo` reads expose provider Cost Observations with company/Product/time/provider-local qualifiers.

Native variants such as `CUSGER`, `CUSREP`, `CUSSEMICM` and `CUSMED` remain observations, never Cost Basis by inheritance. Commercial Economics owns Cost Basis selection/interpretation.

### 9.2 Expected Tax — G1 current classification

`POST /v1/fiscal/impostos/calculo` remains the preferred Integration-Supported calculation surface rather than copying Sankhya fiscal rules into MPC.

Current corrected evidence establishes:

- Nota/Pedido models are discoverable as native headers with `TIPMOV='Z'` in this SourceInstance;
- a correct e-commerce model was created by the operator and verified read-only;
- the sanctioned tax-calculation API builds a transient/virtual calculation movement and does not expose a seller/vendedor request field;
- Metal Nobre customization requires seller on confirmed sale/order movements;
- the API does not propagate the model/customer seller into the transient movement;
- therefore the current sanctioned Expected-Tax surface is **incompatible with this SourceInstance customization until that customization is adjusted by the ERP owner**;
- MPC performs no trigger/configuration change and Direct Oracle is not an integration fallback.

G1 remains **CONDITIONED / OPEN**. After operator-owned ERP remediation, close G1 with bounded non-persisting in-state/out-of-state calculation probes plus zero-residue reread. Only if the sanctioned path remains semantically insufficient after its SourceInstance prerequisite is satisfied does `STOP / SPLIT PREREQUISITE` become live.

Historical realized tax evidence remains realized evidence/B4 input, not an invented replacement tax engine for L0.

---

## 10. Coverage, pagination and change semantics

Full/scoped acquisition is the correctness baseline. Point reads prove only their exact source-qualified scope; enumeration proves only the completed traversed scope.

Delta/change-log use is optional and prerequisite-bound. No delta result is accepted as completeness evidence until source change-log coverage is independently established. D7 owns cadence, checkpoints and recovery mechanics.

Provider-specific page origin, `hasMore`, short-page, 404-as-end and date-format quirks remain adapter-local.

---

## 11. Native order observation and precedence

Native business-order identity remains:

```text
SourceInstance + native NUNOTA
```

Current REST order enumeration is useful for bounded discovery/observation but its documented point filters were empirically unreliable in sandbox and production. Bounded sanctioned `CabecalhoNota` read by NUNOTA is the current authoritative point/reread surface for consequential native-order state.

When enumeration and authoritative point reread disagree, the point reread governs current consequential state and divergence remains explicit evidence.

Raw `STATUSNOTA`/`PENDENTE` remain adapter-local. MPC does not acquire a universal `Confirmation` lifecycle stage.

Provider-independent Materialization meaning:

> **the native business order has reached the externally required state sufficient for the next claimed materialization progression.**

For the current Sankhya binding, `CACSP.confirmarNota` and `A → L` are the adapter-local realization.

---

## 12. Business-System Party Resolution

### 12.1 Provider-independent semantic contract

A `Business Order Intent` may require resolving the source-native business-party reference required by the selected business system. This is a bounded Materialization responsibility, not an MPC Customer Master/CRM lifecycle.

The semantic contract keeps distinct:

```text
Marketplace Sale
  ├─ marketplace buyer/account evidence
  ├─ fiscal/billing-party evidence
  └─ delivery-recipient/destination evidence
             ↓
Business-System Party Resolution
             ↓
source-native business-party reference
             ↓
Business Order Materialization
```

Rules:

1. Marketplace Sales owns sale/buyer interpretation and supplies only facts legitimately required for materialization.
2. Marketplace buyer/account identity is provenance/context, not automatically the fiscal party.
3. Fiscal/billing party is not automatically the delivery recipient.
4. Delivery destination is not native master address by inheritance.
5. A prior explicit source-native resolution may be reused only while materially compatible with current fiscal-party evidence and native state.
6. Exactly one sufficiently established compatible native match may be used.
7. Zero native matches may require native creation only when all provider/source-required identity-bearing facts are known from legitimate transaction evidence; otherwise the path is explicit exception/external-required work.
8. Multiple matches or material contradictions are `AMBIGUOUS`: no guessed selection, no first-result-wins and no new duplicate creation.
9. Transaction-specific marketplace billing/delivery data does not overwrite native master data by default.
10. Any native party create/update is a consequential external effect subject to duplicate protection, no-blind-retry after possible acceptance, authoritative reread/correlation, auditability and minimum PII.
11. MPC may preserve bounded resolution/correspondence lineage needed to avoid repeating a resolved ambiguity, but does not thereby own the customer/party master lifecycle.
12. Concurrent/repeated materializations for the same unresolved fiscal identity must not independently create duplicate native parties; D7 chooses the mechanism.

### 12.2 Current Sankhya realization

For the current Sankhya SourceInstance:

- source-native business-party reference is `CODPARC`;
- dedicated customer REST resources provide list/create/update/contact surfaces but are not a sufficient point-lookup surface for safe correspondence;
- bounded sanctioned `Parceiro` entity reads provide the required lookup evidence;
- legal document (`CGC_CPF`) is a conditioned lookup signal, not guaranteed unique identity;
- lookup uses ordinary equality with a string-bound parameter;
- production evidence contains a real seven-way duplicate for one legal document, proving ambiguity is a present failure mode;
- current e-commerce origin custom fields are sparsely populated and are not correspondence authority;
- existing transaction evidence shows delivery data can be represented on the document without requiring blanket Partner-master overwrite;
- any Sankhya Partner create/update remains an external effect under §18.

`Parceiro`, `CODPARC`, `CGC_CPF`, Sankhya contact/address fields and their lookup protocol are adapter/SourceInstance realization only. They are not MPC canonical customer semantics.

### 12.3 G2 status

G2 is **PASS WITH EXPLICIT EXCEPTION PATH at architecture-contract level**. The first real consequential native-party write may be a controlled D8 proof when needed.

A narrow read-only realization probe may still refine which current Sankhya document/contact fields carry transaction delivery data and which minimum Partner fields the selected TOP-313 lane requires; that evidence may amend the Sankhya realization without reopening the provider-independent contract.

---

## 13. Business Order Intent materialization

### 13.1 Provider-independent semantic contract

```text
Business Order Intent
  + Organization
  + SourceInstance
  + canonical Marketplace Sale context
  + transaction-specific Selling Entity attribution
  + required Product / quantity / Business-System Party Resolution / materialization facts
→ attempt native materialization
→ source-qualified native business-order result(s)
→ authoritative convergence / rejection / pending / ambiguity
```

The domain does not own TOPs, series, TIPMOV, CACSP, NUNOTA letters, TGFVAR or a provider-native customer master.

Provider-native intermediate artifacts may exist internally. Their number and sequence are adapter concerns and do not become configurable MPC workflow steps.

### 13.2 Current Metal Nobre e-commerce binding

Production evidence establishes:

```text
create native e-commerce order
  current provider binding: TOP 313
  order series: PA
  Mercado Livre discriminator in this SourceInstance: CODTIPVENDA 27
  business-party reference: resolved CODPARC
→ satisfy source-required progression state
  current Sankhya realization: CACSP.confirmarNota / A→L
→ authoritative reread
→ native business order converged
```

TOP 313 is e-commerce generally, not Mercado Livre identity. The current TOP has provider-native reservation/financial effects; those are binding facts, not MPC policy.

### 13.3 Creation surface

For the selected current binding, `CACSP.incluirNota` is the proven MGECOM creation surface. REST `POST /v1/vendas/pedidos` remains provider-conditioned and is not preferred merely because it is REST.

Measured partial-update behavior of `CACSP.incluirNota` is not generalized into an arbitrary patch API.

### 13.4 Required input binding — not a knob bag

For each selected operation, only values actually required for correct execution are bound, sourced explicitly from stable SourceInstance configuration, current domain-owned intent/context or externally governed/provider-derived prerequisite/default.

Observed variation in vendor, nature, carrier, partner or other fields does not justify one global setting per field.

---

## 14. Binding validation and hidden SourceInstance rules

`TipoOperacao` is version-qualified by `DHALTER`; `CODTIPOPER` alone is not eternal meaning.

Observable provider-declared properties may include movement role, stock/financial/pendency/fiscal posture and active/effective version. Material assumptions must still be established.

Provider-declared configuration is not a complete description of SourceInstance behavior. Custom triggers, liberação/approval rules and procedural customizations may impose requirements absent from the sanctioned configuration surface.

Therefore binding validation:

1. detects observable provider-configuration drift;
2. is necessary but never sufficient for a consequential effect;
3. does not claim to pre-validate all provider/customization prerequisites;
4. leaves hidden/unexposed SourceInstance rules as execution-time uncertainty;
5. never turns successful validation into a prediction of write success.

D7 later decides validation cadence/cache mechanics.

---

## 15. Invoicing Intent materialization

Provider-independent contract:

```text
Invoicing Intent
  + native business-order result in the source-required progression state
  + Fulfillment-owned physical readiness/conference meaning
  + required provider/business-system prerequisites
→ sanctioned native fiscal progression
→ source-qualified fiscal result
→ authoritative reread + origin/result correlation
→ converged / rejected / pending / ambiguous
```

Physical readiness remains Fulfillment authority.

Current production history proves:

```text
TOP 313 native order → TOP 306 native fiscal result
```

with distinct native identities and line/quantity correlation.

`SelecaoDocumentoSP.faturar` is the selected sanctioned progression surface. The first controlled real 313→306 fiscal write remains D8 proof because it is an irreversible/legal effect whose architectural contract is already grounded.

---

## 16. Native result correlation

D4 preserves source-native correlation without creating a generic provider graph.

Requirements:

- origin and result identities remain distinct source-qualified references;
- 0..N/partial results remain representable when exposed;
- line/quantity granularity is preserved when material;
- provider relation resources remain adapter evidence, not MPC entities;
- 2xx transform response is not convergence.

Current Sankhya realization may use TGFVAR/`CompraVendavariosPedido`; business domains do not read that provider resource directly.

---

## 17. Pre-invoice commercial reversal vs post-invoice fiscal return

Current SourceInstance history proves an alternative branch in which uninvoiced TOP-313 orders correlate to TOP-307 results.

Known: the observed 307s originate from uninvoiced 313 orders, reverse commercial/financial posture, have no stock write-down, and are not observed as results of TOP-306 invoices.

> **Pre-invoice native commercial reversal is not the same business-system consequence as post-invoice fiscal return/reversal.**

Unknown/deferred: the sanctioned write command for the observed 307-class consequence, whether it releases the originating inventory commitment, and the post-invoice fiscal-return path.

These unknowns do not block B3 while reversal actuation remains explicit `external-required`; they become proof gates before automated reversal is claimed.

---

## 18. External-effect semantics

Every consequential business-system write admitted by the selected contract — including native party create/update, order create/update, source-required order progression, invoicing and any later reversal effect — obeys:

1. explicit Organization + SourceInstance target;
2. explicit owning-domain intent/correlation anchor;
3. known provider/binding prerequisites established no stronger than evidence permits;
4. request scoped only to intended effect;
5. response classified no stronger than provider evidence;
6. accepted/rejected/pending/ambiguous preserved where reachable;
7. authoritative reread/correlation after possible acceptance;
8. no blind retry after timeout/connection loss when acceptance is possible;
9. duplicate/ambiguity conditions become explicit exception work;
10. provider rule/custom-trigger/liberação failures are translated rather than leaked as business semantics;
11. protocol support never bypasses Readiness, Availability, Fulfillment, Governance or other owner validity;
12. provider PII is minimized.

---

## 19. Current Product 1.0 proof lane

### Read/fact lane

- Sankhya production SourceInstance through sanctioned Gateway;
- Product identity/identifier evidence;
- qualified company/location inventory evidence;
- dedicated net stock plus bounded entity decomposition when required;
- cost as `Custo` observations;
- expected tax subject to G1.

### Marketplace business-order lane

```text
Marketplace Sale
→ Business-System Party Resolution
   current Sankhya realization: resolve CODPARC
→ bounded Metal Nobre e-commerce binding
→ native TOP-313 order
→ current Sankhya source-required progression
→ authoritative reread
→ native business-order convergence
```

`CODTIPVENDA=27` is a current SourceInstance discriminator for Mercado Livre; neither it nor TOP 313 becomes marketplace semantics.

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

## 20. Residual gates and defers

### G1 — Expected Tax — **B3 CLOSURE GATE / OPEN**

Current root cause is a structural incompatibility between the sanctioned calculation API and a Metal Nobre SourceInstance customization requiring seller on the transient calculation movement. ERP-owner remediation is external to MPC and pending.

After remediation, bounded non-persisting in-state/out-of-state probes plus zero-residue reread decide:

- `PASS` → G1 closes;
- still semantically insufficient → `STOP / SPLIT PREREQUISITE` for Expected-Tax capability.

Do not copy provider tax rules or use Oracle as fallback.

### G2 — Business-System Party Resolution — **PASS WITH EXPLICIT EXCEPTION PATH**

The provider-independent contract is closed. A narrow read-only Sankhya realization probe may refine delivery/address/contact realization and minimum create prerequisites without reopening the semantic contract.

### G3 — First selected-lane fiscal effect — **DEFER SAFELY → D8**

Controlled actual 313→306 effect plus authoritative fiscal reread/correlation.

### G4 — Controlled-product marketplace lane — **DEFER SAFELY**

Before automated controlled-product operation is claimed, establish satisfiability/interchangeability and selection timing. No adapter-chosen `CONTROLE`, no aggregate guess.

### G5 — Post-invoice fiscal return — **DEFER SAFELY / EXTERNAL-REQUIRED**

Current e-commerce history does not establish path/command.

### G6 — Pre-invoice reversal inventory commitment fate — **DEFER WITH REVERSAL CLAIM**

Whether the observed 307 branch releases the originating reservation remains Unknown and must close before automated reversal/convergence is claimed.

### D7 defers

Token refresh; full/delta cadence; checkpoints/cursors; rate/concurrency/backoff; binding-validation cache/cadence; and the mechanism that prevents concurrent duplicate native-party creation for the same unresolved fiscal identity.

---

## 21. YAGNI / explicit non-goals

B3 MUST NOT introduce:

- generic ERP business entity or universal ERP ontology;
- universal `ERPAdapter` containing every possible operation;
- generic provider/resource/capability graph;
- generic `Customer`, `Party` or CRM master lifecycle merely for integration symmetry;
- universal party-matching engine;
- plugin registry/factory framework for speculative providers;
- arbitrary workflow/materialization DSL;
- MPC TOP/NUNOTA/TGFVAR/CONTROLE/Sankhya-status entities;
- universal Lot/Batch/Serial model derived from `CONTROLE`;
- Sankhya product/stock/cost/tax/customer database mirror;
- arbitrary SQL/DbExplorer behavior through `loadRecords.criteria` or another Gateway route;
- duplicated Sankhya tax engine;
- support for every Metal Nobre non-marketplace process;
- speculative TOTVS/Bling adapters;
- family-level inference of control semantics;
- blanket marketplace→ERP customer-master synchronization.

The observed store lane `14→303→305` remains variability/counterexample evidence, not a Product 1.0 MPC workflow requirement.

---

## 22. Provider-independent replacement test

These MPC meanings should remain valid if a future accepted business system genuinely supplies the same semantics:

- source-qualified Product facts;
- inventory evidence sufficient for Sellable Availability derivation;
- Cost Observations;
- expected/realized fiscal evidence requirements;
- Business-System Party Resolution as a bounded Materialization prerequisite;
- marketplace buyer/account, fiscal/billing party and delivery evidence remaining distinct unless transaction evidence relates them;
- Business Order Intent;
- native business-order convergence on the externally required state;
- Invoicing Intent;
- native fiscal convergence;
- origin/result line/quantity correlation;
- pre-invoice commercial reversal distinct from post-invoice fiscal consequence;
- honest unknown/unavailable;
- authoritative reread;
- no blind retry after ambiguous effects.

Replaceable Sankhya-local pieces include `Parceiro`, `CODPARC`, `CGC_CPF`, contact/address fields, TOP/versioning, NUNOTA, series, statuses, CACSP/SelecaoDocumentoSP, TGFVAR, `TIPCONTEST`/`CONTROLE`, auth/Gateway, triggers/liberações and all Metal Nobre binding values.

If a future real system reveals genuinely different business meaning, reopen the responsible semantic decision rather than contort it into a false common model.

Marketplace Installation and business-system SourceInstance remain distinct; no generic `IntegrationInstance` is justified.

---

## 23. Proposed B3 outcome after G1

Subject to G1 evidence and operator ratification:

### `CURRENT STRUCTURE CONFIRMED`

- consumer-owned semantic ports;
- SourceInstance-qualified external facts/results;
- concrete provider adapter;
- bounded Business-System Party Resolution under Materialization, not Customer Master authority;
- no Integration business domain;
- no generic ERP/workflow/customer framework.

### `CURRENT STRUCTURE CONFIRMED` — Sankhya target

The sanctioned Gateway/API surface is sufficient to define the current Product 1.0 business-system contract without Direct Oracle fallback, subject to G1 remediation/proof and honest capability fences.

### `DEFER SAFELY`

- D7 runtime/concurrency mechanics;
- first irreversible selected-lane fiscal write → D8;
- first consequential native-party write → D8 when needed;
- controlled-product marketplace automation;
- automated pre-invoice reversal until reservation-fate/command proof;
- post-invoice fiscal return unless selected golden flow requires automated actuation.

### `STOP / SPLIT PREREQUISITE`

If G1 or another materially required Product 1.0 claim cannot be satisfied correctly through sanctioned SourceInstance operations, stop and re-adjudicate that capability. Direct Oracle/database is never admitted implicitly.

No D0/D1/D2/D3/B1/B2 reopen is proposed now.

---

## 24. Proof state

Already evidenced proportionately:

- OAuth/environment binding and short token TTL;
- Product native key/identifier limits;
- company/location/control stock granularity;
- REST net formula and negative-net behavior;
- reservation decomposition;
- controlled-product/source decomposition evidence;
- cost observations;
- change-log prerequisite/failure honesty;
- REST-order point-filter failure and bounded point-reread fallback;
- `CACSP.incluirNota` creation;
- `CACSP.confirmarNota` success/rejection + reread;
- MGECOM OAuth compatibility;
- `SelecaoDocumentoSP.faturar` capability;
- source-native origin/result line/quantity correlation;
- production e-commerce topology `313→306`;
- observed `313→307` branch;
- binding configuration/version reads plus hidden-rule counterevidence;
- current control-free marketplace sold lane;
- G2 safe resolution contract, including measured duplicate ambiguity and no blanket master-update requirement;
- production-vs-sandbox materialization divergence.

Still required for B3 whole acceptance:

- G1 Expected Tax after SourceInstance remediation.

Optional remaining realization evidence before final consolidation:

- narrow read-only Sankhya Party Resolution probe for transaction delivery/contact representation and minimum Partner-create prerequisites.

D8/deferred proofs are not B3 closure blockers unless their capability becomes a claimed B3 normal path.

---

## 25. Reopen triggers

Reopen only the implicated decision when evidence shows, for example:

- a second real business system cannot implement accepted consumer meaning without Sankhya semantics leaking into the domain;
- Business-System Party Resolution grows into an independent customer/party business lifecycle with authority beyond bounded materialization correspondence;
- a selected Sankhya operation disappears/changes materially and no sanctioned replacement is sufficient;
- a controlled Product enters marketplace scope and its sellability/selection semantics cannot fit existing Availability/Fulfillment ownership;
- automated post-invoice fiscal return becomes a claimed normal path and external-required handling is insufficient;
- Product 1.0 requires a fact/effect for which only arbitrary SQL/Direct Oracle could satisfy the claim;
- second-provider repetition proves a smaller shared technical mechanism materially reduces total complexity.

Do not reopen for naming preference, abstract symmetry or hypothetical providers.

---

## 26. Candidate disposition

This file is a disposable design/review surface only.

It MUST be deleted after the remaining evidence, final adjudication/operator ratification and before canonical D4-B3 consolidation.

Canonical stage/status remains solely in `docs/engineering/rebaseline/README.md`. This file does not open/accept B3, authorize B4/D5+, authorize implementation, merge or product writes.
