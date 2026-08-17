# D4-B3 — Sankhya Business-System Contract — REVIEW CANDIDATE

> **Status:** REVIEW CANDIDATE / NON-AUTHORITATIVE / DISPOSABLE  
> **Base evidence HEAD:** `eaab7127518002949ebdfa00aead90151a85ec56`  
> **Independent review HEAD:** `f7ec08d91108ed905133874bb5bcc26f1b729b2b`  
> **Review evidence commit:** `039e81082c8e5ab687b1537063b210caeb322c3f`  
> **Program:** Architecture Rebaseline / Technical System Design  
> **Parent authority:** accepted D0–D4-B2 only  
> **Method:** DevelopmentConexus Engineering Method v1.0.0  
> **Purpose:** amended coherent surface for closing the remaining D4-B3 evidence gates. This file is not target authority and MUST be deleted before canonical consolidation.

---

## 1. Review question and scope

D4-B3 must answer:

> **Through which concrete sanctioned business-system contracts can MPC obtain the internal facts and cause/reconcile the native business-order and fiscal effects required by Product 1.0, while preserving D0–D3 semantic authority, explicit SourceInstance qualification, honest coverage and a future business-system seam without building a generic ERP/workflow framework?**

The target is intentionally:

- provider-independent at the MPC semantic boundary;
- concrete about Sankhya inside D4 because Sankhya is the first business system proving that boundary;
- SourceInstance-aware because real customer configuration changes how the provider contract behaves;
- bounded by YAGNI: no universal ERP ontology, plugin framework or workflow/BPM engine.

A future TOTVS/Bling/SAP-like system is only a replacement test. It is not modelled now.

Implementation remains blocked until D9.

---

## 2. Imported authority — not reopened

B3 imports rather than re-decides:

1. Business-System Materialization owns `Business Order Intent`, `Invoicing Intent`, native-result correlation and convergence meaning.
2. ERP-native TOP/document taxonomy is not MPC semantics.
3. Product master remains external; MPC references source-qualified Product identity.
4. Availability Control owns Inventory Source/Scope, allocation meaning and Sellable Availability; native stock/reservation truth remains external.
5. Commercial Economics owns Cost Basis and economic interpretation; provider cost variants are evidence only.
6. Marketplace Sales owns canonical marketplace-sale interpretation and transaction-specific Selling Entity attribution.
7. Post-Sale Resolution coordinates consequences without importing provider-native fiscal taxonomy.
8. `SourceInstance` identifies one logical externally authoritative business-system/source namespace; credentials are not identity.
9. Consumer owns meaning; adapter owns protocol.
10. Known / absent / unknown / unavailable remain distinct; partial is not closure.
11. Acceptance is not convergence; ambiguous writes are not blindly retried.
12. Provider PII is minimized.
13. Direct Oracle/database access is outside target architecture and is not fallback.
14. Integration Support, Provider Effective Capability/Requirement and Effective Business Capability remain separate authorities.

No D0/D1/D2/D3 or D4-B1/B2 reopen is proposed.

---

## 3. Governing invariant

> **For each Product 1.0 business-system fact or effect, MPC depends on a consumer-owned semantic contract whose external evidence/effect is Organization + SourceInstance qualified. The concrete Sankhya adapter uses only sanctioned provider operations and bounded SourceInstance-specific bindings, preserves provider-native granularity/provenance/partiality, and reconciles consequential effects through authoritative reread/correlation. Provider-declared configuration may be validated where observable but never guarantees execution success; hidden/custom rules remain execution-time uncertainty. Sankhya-native topology and customer configuration never become MPC business ontology, and future substitutability does not authorize a generic ERP/workflow framework before a second real consumer proves the abstraction.**

Corollaries:

- Sankhya first does not mean Sankhya model.
- SourceInstance binding does not mean workflow engine.
- Provider stock decomposition does not imply an MPC `Lot`/Batch entity or proven interchangeability.
- TOP/code identity alone does not prove current configured meaning.
- Provider configuration evidence is not MPC policy.
- A validated binding is a necessary precondition, never a prediction that a consequential effect will succeed.

---

## 4. Alternatives / Global Maximum

### A — Sankhya-shaped core

Examples: `Business Order = TOP 313`, `Invoice = TOP 306`, `Marketplace = CODTIPVENDA 27`, `Lot = CONTROLE`.

**REJECT.** One current SourceInstance already contains materially different store and e-commerce native document topologies. A Sankhya-shaped core fails the present before any second ERP exists.

### B — generic ERP/workflow framework now

Examples: `GenericERP`, generic resource/capability graph, configurable materialization step sequence, arbitrary provider-operation DSL.

**REJECT.** There is one concrete business system and no second real consumer proving a stable shared provider ontology or workflow language.

### C — provider-independent consumer semantics + concrete Sankhya adapter + bounded SourceInstance binding

**PROPOSED GLOBAL MAXIMUM.**

```text
D1 owner meaning
    ↓ consumer-owned semantic port
D4 Sankhya adapter
    - sanctioned operations
    - bounded SourceInstance binding
    - provider capability/requirements
    - authoritative reread/correlation
    ↓
Sankhya SourceInstance
```

**Falsifier:** if MPC/domain configuration must express arbitrary provider-native document count, sequence or conditional choreography, the boundary has failed and must return to decision rather than evolve into BPM/workflow infrastructure.

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

Dedicated REST resources remain preferred when they are semantically sufficient. `CRUDServiceProvider.loadRecords` is admitted only when a real consumer requires a fact unavailable or materially lossy on the dedicated resource.

The admission is intentionally narrower than what the Gateway parser may technically accept:

1. a named sanctioned `rootEntity` is required;
2. result fieldsets are explicit and minimum for the consumer claim;
3. `criteria` is restricted to predicates over fields of the named root entity using ordinary comparison/logical operators and bound parameters;
4. sanctioned related-entity data, when genuinely required, uses provider-declared relation/path mechanisms with explicit fieldsets rather than ad-hoc cross-table expressions;
5. subqueries, cross-table references inside criteria, Oracle-specific pseudo-tables/functions, arbitrary SQL expressions and query-language escape hatches are **outside the D4 sanctioned entity-read contract**, even if the Gateway happens to execute them;
6. a future need for such an expression is a capability finding requiring adjudication, not authorization for Oracle-via-HTTP;
7. no entity read becomes a replicated ERP mirror by convenience.

This clause protects the accepted Direct-Oracle exclusion at the semantic contract level; D7/implementation later chooses mechanical enforcement.

---

## 7. Product / Readiness evidence

### 7.1 Native Product identity

Current external Product reference:

```text
SourceInstance + CODPROD
```

`CODPROD` remains provider-native, not an MPC Product-master identity.

### 7.2 Identifier evidence

Available provider observations include active state, reference, supplier reference, brand, NCM/fiscal facts and alternate-volume barcode where populated.

No universal first-class GTIN/EAN Product field was established. Therefore `REFERENCIA == EAN` is forbidden as an identity/correspondence law.

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

at the provider dimensions represented by the resource.

The REST net surface is useful but lossy because it cannot explain the reservation decomposition. When a consumer needs that distinction, the sanctioned `Estoque` entity read preserves the minimum required dimensions/fields such as:

```text
CODEMP
CODLOCAL
CODPROD
CONTROLE
TIPO
CODPARC
ESTOQUE
RESERVADO
```

The net surface may legitimately be **negative** when commitments exceed physical stock. Negative is a real external observation; it is neither zero nor unknown and MUST NOT be clamped away at the adapter boundary.

Dedicated REST stock also has provider-defined coverage limits; resource absence is not automatically `stock=0`.

### 8.3 Provider control dimension

`CONTROLE` is preserved as an opaque provider inventory-partition value when present. It is not canonically named `Lot`, `Batch`, `Tonality`, `Serial` or another MPC identity.

Current SourceInstance evidence:

- control is per Product, not reliably inferred from family;
- 3,038 active controlled Products were observed;
- controlled Products concentrate in flooring/covering families including porcelain, but those same families contain uncontrolled counterexamples;
- one raw control-type code (`TIPCONTEST='I'`) is currently observed, so no multi-type MPC taxonomy is justified;
- `CONTROLE` carries physical-looking lot strings and also `ENCOMENDA`;
- no sanctioned structured source for tonalidade/calibre/grade attributes was established.

### 8.4 Provider-independent sellability invariant

Measured evidence proves an aggregate free quantity may exceed every individual provider partition. The D1-facing rule is therefore about **satisfiability**, not provider topology:

> **Sellable Availability may treat an aggregate quantity as sellable only when the available evidence and applicable rules establish that the requested quantity is actually satisfiable by the eligible source inventory. Where an external source decomposes inventory in a way material to satisfiability, D4 preserves that decomposition as evidence rather than erasing it into a pre-aggregated total.**

This does not assert partitions are always non-interchangeable.

Current Mercado Livre **sold** population evidence is control-free: all observed e-commerce items were uncontrolled (`TIPCONTEST='N'`). The current **listed** ML population was not cross-matched for control sensitivity and remains unknown.

First Product 1.0 capability fence:

- uncontrolled current marketplace lane: supported subject to ordinary Availability rules;
- controlled Product automatic marketplace availability/materialization: not claimed until interchangeability/satisfiability and selection timing are proven for that lane;
- no adapter-chosen partition and no aggregate guess.

No D1 reopen is required by current evidence.

### 8.5 Materialization-created inventory commitment

The selected native order binding reserves inventory in this SourceInstance. That is a provider effect whose business meaning remains Availability-owned.

Therefore:

> **When native Business Order materialization creates or changes an external inventory commitment, the inventory observation contract must preserve that effect so Availability Control can account for current sellability. Materialization does not acquire allocation authority, and MPC must not model its own materializations as inventory-neutral.**

This is not a new Materialization→Availability business-authority edge: Availability observes its authoritative external inventory source as already accepted.

---

## 9. Cost and expected-tax evidence

### 9.1 Cost

No dedicated current cost REST surface was established. Bounded sanctioned `rootEntity=Custo` reads expose provider observations with company/Product/time/provider-local qualifiers.

`CUSGER`, `CUSREP`, `CUSSEMICM`, `CUSMED` and other native variants remain **Cost Observations**, never Cost Basis by inheritance. Commercial Economics owns Cost Basis selection/interpretation.

Sentinel/placeholder rows are not silently promoted to valid cost facts.

### 9.2 Expected tax — current classification

`POST /v1/fiscal/impostos/calculo` is Integration-Supported and documented as a calculation surface that applies Sankhya configuration rather than requiring MPC to copy a tax engine.

The observed failure is currently classified more precisely as:

> **SourceInstance configuration prerequisite not yet satisfied for the selected path.**

The exercised calculation encountered a Metal Nobre customization requiring seller data during internal movement preparation. The current evidence does not show that Sankhya lacks the capability; it shows that the selected SourceInstance still needs a native calculation/model configuration that satisfies its own customizations.

B3 decision:

- keep the sanctioned calculation surface as the preferred Expected-Tax path;
- do not copy TGFICM/provider fiscal rules into MPC;
- G1 remains a **B3 closure gate**;
- close G1 by identifying/configuring the correct native model/path and performing a non-persisting re-probe;
- only if the correctly configured sanctioned calculation remains semantically insufficient for L0 Expected Economics does `STOP / SPLIT PREREQUISITE` apply.

Realized fiscal/tax evidence comes from authoritative fiscal results, not the Expected-Tax simulation.

---

## 10. Coverage, pagination and change semantics

Full/scoped acquisition is the correctness baseline. Point reads prove only their exact source-qualified scope; enumeration proves only the completed traversed scope.

`modifiedSince`/change-log use is optional and prerequisite-bound. Current evidence shows change logging disabled and failure honesty differs by resource; Product delta may look like an empty/no-record result when the prerequisite is absent.

Therefore no delta result is accepted as completeness evidence until source change-log coverage is independently established. Even when enabled, the documented short retention window prevents delta from becoming indefinite recovery/history authority.

B3 does not require enabling `LOGTABOPER`; D7 owns cadence, checkpoints and recovery mechanics.

Provider-specific page origin, `hasMore`, short-page, 404-as-end and date-format quirks remain adapter-local.

---

## 11. Native order observation and precedence

The native business-order identity remains:

```text
SourceInstance + native NUNOTA
```

No synthetic `MPCOrder` identity is introduced.

Current REST order enumeration is useful for bounded discovery/observation but its documented point filters were empirically unreliable in both sandbox and production. The sanctioned `CabecalhoNota` read by NUNOTA is therefore the current authoritative point/reread surface for consequential native-order state.

**Precedence rule:** when a bounded enumeration observation and the authoritative point reread disagree for a consequential decision, the point reread governs the current native state and the divergence remains explicit evidence. The two surfaces are not treated as competing authorities.

Raw `STATUSNOTA`/`PENDENTE` vocabulary remains adapter-local.

Measured Sankhya facts include created/unconfirmed (`A`), confirmed (`L`) and orthogonal pendency combinations. MPC does **not** acquire a universal `Confirmation` lifecycle stage.

The provider-independent Materialization meaning is:

> **the native business order has reached the externally required state sufficient for the next claimed materialization progression.**

For the current Sankhya binding, `CACSP.confirmarNota` and `A → L` are the adapter-local realization of that requirement.

`CACSP.confirmarNota` remains officially undocumented but empirically established through MGECOM in sandbox and production, including clean rejection on a custom prerequisite and authoritative reread after success.

---

## 12. Native customer/partner prerequisite

Current e-commerce evidence proves that native orders use real distinct customer/partner records, not one generic Mercado Livre partner.

B3 therefore requires a bounded Materialization prerequisite, not a new Customer Master domain:

1. Marketplace Sales supplies only buyer facts legitimately required for the sale/materialization purpose.
2. D4 resolves the source-native partner using sanctioned customer surfaces.
3. Matching must fail honest on none/multiple/ambiguous candidates; mutable names/emails are not silently treated as canonical identity.
4. Only minimum PII needed for business/fiscal processing crosses the boundary and is retained proportionately.
5. A source-native partner reference must be established before consequential order creation.
6. **Any native customer create or update is itself a consequential external effect** and is governed by the full external-effect contract in §18: explicit correlation anchor, duplicate protection, no blind retry after ambiguous possible acceptance, authoritative reread, auditable outcome and minimum PII.
7. Duplicate or ambiguous partner resolution becomes explicit exception work, never guessed matching.

G2 closes the architecture contract for matching/create/update semantics and the minimum safe source-native path. The first real consequential customer write may be proven in D8 when needed; B3 must not defer discovering the matching/duplicate semantics until after architecture acceptance.

---

## 13. Business Order Intent materialization

### 13.1 Provider-independent semantic contract

```text
Business Order Intent
  + Organization
  + SourceInstance
  + canonical Marketplace Sale context
  + transaction-specific Selling Entity attribution
  + required product / quantity / customer / materialization facts
→ attempt native materialization
→ source-qualified native business-order result(s)
→ authoritative convergence / rejection / pending / ambiguity
```

The domain does not own TOPs, series, TIPMOV, CACSP, NUNOTA letters or TGFVAR.

Provider-native intermediate artifacts may exist internally. Their number and sequence are adapter concerns and do not become configurable MPC workflow steps.

### 13.2 Current Metal Nobre e-commerce binding

Production evidence establishes:

```text
create native e-commerce order
  current provider binding: TOP 313
  order series: PA
  Mercado Livre discriminator in this SourceInstance: CODTIPVENDA 27
→ satisfy source-required progression state
  current Sankhya realization: CACSP.confirmarNota / A→L
→ authoritative reread
→ native business order converged
```

TOP 313 is e-commerce generally, not Mercado Livre identity. Multiple e-commerce negotiation types exist on the same TOP.

The current TOP has provider-native reservation/financial effects. Those are binding facts, not MPC policy.

### 13.3 Creation surface

For the selected current binding, `CACSP.incluirNota` is the proven MGECOM creation surface. The newer REST `POST /v1/vendas/pedidos` remains provider-conditioned because its required native model setup is unresolved; B3 does not prefer it merely because it is REST.

Measured partial-update behavior of `CACSP.incluirNota` is not generalized into an arbitrary patch API.

### 13.4 Required input binding — not a knob bag

For each selected operation, only values actually required for correct execution are bound, sourced explicitly from:

- stable SourceInstance configuration;
- current domain-owned intent/context;
- externally governed/provider-derived prerequisite/default.

Observed variation in vendor, nature, carrier, partner or other fields does not justify one global setting per field.

---

## 14. Binding validation and hidden SourceInstance rules

`TipoOperacao` is version-qualified by `DHALTER`; `CODTIPOPER` alone is not eternal meaning.

Provider-declared properties available through sanctioned configuration reads include, where material:

- movement role;
- stock effect;
- financial posture;
- pendency posture;
- fiscal-model posture;
- active/effective version.

A consequential binding therefore cannot mean only `orderTop=313` / `invoiceTop=306`; provider-declared assumptions material to the contract must still be established.

However, independent review established an essential limit: provider-declared configuration is **not a complete description of SourceInstance behavior**. Custom database triggers, liberação/approval rules and procedural customizations may impose requirements that are not exposed by the sanctioned configuration surface and may even contradict a provider-declared flag.

Therefore:

1. binding validation detects observable provider-configuration drift;
2. it is **necessary but never sufficient** for a consequential effect;
3. it MUST NOT claim to pre-validate all confirmation/provider/customization prerequisites;
4. hidden/unexposed SourceInstance rules remain explicit execution-time uncertainty;
5. consequential execution remains fail-closed under provider rejection/pending/ambiguity even after a binding validates;
6. successful validation is never interpreted as a prediction of write success.

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

Physical readiness remains Fulfillment authority. Materialization does not invoice merely because a native order exists.

Current SourceInstance production history proves the e-commerce result topology:

```text
TOP 313 native order → TOP 306 native fiscal result
```

with distinct native identities and line/quantity correlation.

`SelecaoDocumentoSP.faturar` is the selected sanctioned progression surface. It is documented, OAuth/Bearer-compatible and empirically exercised in production on a non-fiscal transformation; existing native history proves 313→306 fiscal results/correlation.

The first controlled real 313→306 fiscal write remains D8 proof because it is an irreversible/legal effect whose architectural contract is already grounded.

---

## 16. Native result correlation

D4 preserves source-native correlation without creating a generic provider graph.

Current Sankhya evidence includes native result correlation and `CompraVendavariosPedido`/TGFVAR line + `QTDATENDIDA` relation evidence.

Requirements:

- origin and result identities remain distinct source-qualified references;
- 0..N/partial results remain representable when the source exposes them;
- line/quantity granularity is preserved where material;
- provider relation resources remain adapter evidence, not MPC entities;
- 2xx transform response is not convergence.

Commercial Economics/Post-Sale consume Materialization-owned interpretation through semantic boundaries; they do not read TGFVAR directly.

---

## 17. Pre-invoice commercial reversal vs post-invoice fiscal return

Current SourceInstance history proves an alternative branch in which uninvoiced TOP-313 orders correlate to TOP-307 results.

Known:

- observed 307s originate from uninvoiced 313 orders;
- 307 reverses commercial/financial posture and has no stock write-down;
- no observed 307 originates from a 306 fiscal result;
- 313→306 and 313→307 are alternative observed result branches.

Therefore:

> **Pre-invoice native commercial reversal is not the same business-system consequence as post-invoice fiscal return/reversal.**

Unknown/deferred:

- the sanctioned write command that produces the observed 307-class consequence;
- whether that reversal releases the inventory commitment/reservation created by the originating order;
- the post-invoice fiscal-return/reversal path.

Availability and Post-Sale MUST NOT assume either reservation-release outcome from the current evidence.

These unknowns do **not** block B3 whole acceptance while reversal actuation remains explicit `external-required`. They become proof gates before MPC claims automated pre-invoice reversal or automated post-invoice fiscal return.

---

## 18. External-effect semantics

Every consequential business-system write admitted by the selected contract — including native customer create/update, order create/update, source-required order progression, invoicing and any later reversal effect — obeys:

1. target Organization + SourceInstance explicit;
2. owning domain intent/correlation anchor explicit;
3. current known provider/binding prerequisites established no stronger than evidence permits;
4. request scoped only to intended effect;
5. response classified no stronger than provider evidence;
6. accepted/rejected/pending/ambiguous preserved where reachable;
7. authoritative reread/correlation after possible acceptance;
8. timeout/connection loss after possible acceptance is never blindly retried;
9. duplicate/ambiguity conditions become explicit exception work;
10. native business-rule/custom-trigger/liberação failures are translated as provider requirement/rejection evidence rather than raw provider implementation leakage to business domains;
11. protocol support never bypasses Readiness, Availability, Fulfillment, Governance or other owner validity;
12. provider PII is minimized.

Measured confirmation already proves that creation success does not establish full materialization capability and that hidden SourceInstance rules can reject later progression.

---

## 19. Current Product 1.0 proof lane

### Read/fact lane

- Sankhya production SourceInstance through sanctioned Gateway;
- Product identity/identifier evidence as measured;
- qualified company/location inventory evidence;
- dedicated net stock plus bounded entity decomposition when required;
- cost as `Custo` observations;
- expected tax subject to G1.

### Marketplace business-order lane

```text
Marketplace Sale
→ bounded Metal Nobre e-commerce binding
→ native TOP-313 order
→ current Sankhya source-required progression (confirmation)
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

### Inventory-control fence

Current observed ML sold lane is uncontrolled. Automatic controlled-product marketplace Availability/Materialization remains unclaimed until satisfiability/interchangeability and selection timing are proven.

---

## 20. Residual gates and defers

### G1 — Expected Tax — **B3 CLOSURE GATE**

Close the SourceInstance configuration prerequisite and perform a non-persisting calculation re-probe proving the selected sanctioned path is usable and semantically sufficient for L0 Expected Economics.

If, after correct SourceInstance configuration, the sanctioned calculation remains materially insufficient, return `STOP / SPLIT PREREQUISITE`. Do not copy provider tax rules or use Oracle as fallback.

### G2 — Native customer/partner — **B3 CLOSURE GATE**

Establish the concrete minimal safe path for:

- buyer evidence required from Marketplace Sales;
- source-native partner lookup/matching;
- uniqueness/duplicate/ambiguity behavior;
- minimum required PII;
- when create/update is necessary;
- authoritative reread/correlation for customer writes;
- no-blind-retry behavior after ambiguous possible acceptance.

Customer API existence alone is insufficient. B3 closes the contract; first consequential customer write may be controlled in D8 when required.

### G3 — First selected-lane fiscal effect — **DEFER SAFELY → D8**

Controlled actual 313→306 effect plus authoritative fiscal reread/correlation.

### G4 — Controlled-product marketplace lane — **DEFER SAFELY**

Before automated controlled-product operation is claimed, establish satisfiability/interchangeability and selection timing. No adapter-chosen `CONTROLE`, no aggregate guess.

### G5 — Post-invoice fiscal return — **DEFER SAFELY / EXTERNAL-REQUIRED**

Current e-commerce history does not establish the path/command. Post-Sale can coordinate explicit work without pretending automated support.

### G6 — Pre-invoice reversal inventory commitment fate — **DEFER WITH REVERSAL CLAIM**

Whether the observed 307 branch releases the originating reservation remains Unknown. This does not block B3 while reversal actuation is external-required; it must close before automated reversal/convergence is claimed.

### D7 defers

Token refresh mechanics; full/delta cadence; checkpoints/cursors; rate/concurrency/backoff; configuration-validation cache/cadence; worker/process/deployment topology.

---

## 21. YAGNI / explicit non-goals

B3 MUST NOT introduce:

- generic ERP business entity or universal ERP ontology;
- universal `ERPAdapter` containing every possible operation;
- generic provider/resource/capability graph;
- plugin registry/factory framework for speculative providers;
- arbitrary workflow/materialization DSL;
- MPC TOP/NUNOTA/TGFVAR/CONTROLE/Sankhya-status entities;
- universal Lot/Batch/Serial model derived from `CONTROLE`;
- Sankhya product/stock/cost/tax database mirror;
- arbitrary SQL/DbExplorer behavior through `loadRecords.criteria` or another Gateway route;
- duplicated Sankhya tax engine;
- support for every Metal Nobre non-marketplace process;
- speculative TOTVS/Bling adapters;
- family-level inference of control semantics;
- one config knob for every historical provider field.

The observed store lane `14→303→305` remains variability/counterexample evidence, not a Product 1.0 MPC workflow requirement.

---

## 22. Provider-independent replacement test

These MPC meanings should remain valid if a future accepted business system genuinely supplies the same semantics:

- source-qualified Product facts;
- inventory evidence sufficient for Sellable Availability derivation;
- Cost Observations;
- expected/realized fiscal evidence requirements;
- Business Order Intent;
- native business-order convergence on the externally required state;
- Invoicing Intent;
- native fiscal convergence;
- origin/result line/quantity correlation;
- pre-invoice commercial reversal distinct from post-invoice fiscal consequence;
- honest unknown/unavailable;
- authoritative reread;
- no blind retry after ambiguous effect.

Replaceable Sankhya-local pieces include TOP/versioning, NUNOTA, series, `STATUSNOTA`/`PENDENTE`, CACSP/SelecaoDocumentoSP, TGFVAR, `TIPCONTEST`/`CONTROLE`, auth/Gateway, triggers/liberações and all Metal Nobre binding values.

If a future real system reveals genuinely different business meaning, reopen the responsible semantic decision rather than contort it into a false common ERP model.

Marketplace Installation and business-system SourceInstance remain distinct; no generic `IntegrationInstance` is justified.

---

## 23. Proposed B3 outcome after residual gates

Subject to G1/G2 evidence and operator ratification:

### `CURRENT STRUCTURE CONFIRMED`

- consumer-owned semantic ports;
- SourceInstance-qualified external facts/results;
- concrete provider adapter;
- no Integration business domain;
- no generic ERP/workflow framework.

### `CURRENT STRUCTURE CONFIRMED` — Sankhya target

The sanctioned Gateway/API surface is sufficient to define the current Product 1.0 business-system contract without Direct Oracle fallback, subject to honest capability gates.

### `DEFER SAFELY`

- D7 runtime mechanics;
- first irreversible selected-lane fiscal write → D8;
- controlled-product marketplace automation;
- automated pre-invoice reversal until reservation-fate/command proof;
- post-invoice fiscal return unless a selected golden flow requires automated actuation.

### `STOP / SPLIT PREREQUISITE`

If G1, G2 or another materially required Product 1.0 claim cannot be satisfied correctly through sanctioned SourceInstance operations, stop and re-adjudicate that capability. Direct Oracle/database is never admitted implicitly.

No D0/D1/D2/D3/B1/B2 reopen is proposed now.

---

## 24. Proof state

Already evidenced proportionately:

- OAuth/environment binding and short token TTL;
- Product native key/identifier limits;
- company/location/control stock granularity;
- REST net formula, including production re-measurement;
- negative net-stock possibility;
- reservation decomposition;
- controlled-product/source decomposition evidence;
- cost observations;
- change-log prerequisite/failure-honesty behavior;
- REST-order point-filter failure and bounded entity reread fallback, including production re-measurement;
- `CACSP.incluirNota` creation;
- `CACSP.confirmarNota` success/rejection + reread;
- MGECOM OAuth compatibility;
- `SelecaoDocumentoSP.faturar` capability;
- source-native origin/result line/quantity correlation;
- production e-commerce topology `313→306`;
- observed `313→307` pre-invoice branch;
- binding configuration/version reads;
- current control-free marketplace sold lane;
- production-vs-sandbox materialization divergence.

Still required for B3 whole acceptance:

- G1 Expected Tax;
- G2 native customer/partner contract.

D8/deferred proofs are not B3 closure blockers unless their capability becomes a claimed B3 normal path.

A mock/test cannot substitute for a real external-dependency claim.

---

## 25. Reopen triggers

Reopen only the implicated decision when evidence shows, for example:

- a second real business system cannot implement accepted consumer meaning without Sankhya semantics leaking into the domain;
- a selected Sankhya operation disappears/changes materially and no sanctioned replacement is sufficient;
- provider-declared configuration validation cannot protect even the observable assumptions required by consequential execution;
- a controlled Product enters marketplace scope and its sellability/selection semantics cannot fit existing Availability/Fulfillment ownership;
- customer correspondence requires a genuinely independent MPC business lifecycle rather than a bounded Materialization prerequisite;
- automated post-invoice fiscal return becomes a claimed normal path and external-required handling is insufficient;
- Product 1.0 requires a fact/effect for which only arbitrary SQL/Direct Oracle could satisfy the claim;
- second-provider repetition proves a smaller shared technical mechanism materially reduces total complexity.

Do not reopen for naming preference, abstract symmetry or hypothetical providers.

---

## 26. Candidate disposition

This file is a disposable design/review surface only.

It MUST be deleted after G1/G2 evidence, final adjudication/operator ratification and before canonical D4-B3 consolidation.

Canonical stage/status remains solely in `docs/engineering/rebaseline/README.md`. This file does not open/accept B3, authorize B4/D5+, authorize implementation, merge or product writes.
