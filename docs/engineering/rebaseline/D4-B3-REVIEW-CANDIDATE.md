# D4-B3 — Sankhya Business-System Contract — REVIEW CANDIDATE

> **Status:** REVIEW CANDIDATE / NON-AUTHORITATIVE / DISPOSABLE  
> **Base HEAD:** `eaab7127518002949ebdfa00aead90151a85ec56`  
> **Program:** Architecture Rebaseline / Technical System Design  
> **Parent authority:** accepted D0–D4-B2 only  
> **Method:** DevelopmentConexus Engineering Method v1.0.0  
> **Purpose:** coherent independent-review surface for D4-B3; this file MUST NOT be used as target authority unless its decisions are later adjudicated, ratified by the operator and consolidated into canonical D4.

---

## 1. Review question and scope

D4-B3 must answer:

> **Through which concrete sanctioned business-system contracts can MPC obtain the internal facts and cause/reconcile the native business-order and fiscal effects required by Product 1.0, while preserving D0–D3 semantic authority, explicit SourceInstance qualification, honest coverage and a future business-system seam without building a generic ERP/workflow framework?**

This candidate is intentionally both:

- **provider-independent at the MPC semantic boundary**; and
- **concrete about Sankhya inside D4**, because Sankhya is the first business system used to prove that boundary.

It does **not** claim that Sankhya is the MPC business-system model.

It does **not** model TOTVS, Bling, SAP or any speculative future system. A future concrete system is only a structural replacement test: if the same accepted MPC business meanings can be supplied/materialized by that system, its protocol belongs in another adapter; if genuinely new business meaning appears, the responsible D-stage reopens rather than hiding it in a generic integration abstraction.

This candidate does not choose:

- D5 MPC HTTP/OpenAPI;
- D6 UI;
- D7 scheduler/worker/cursor/token-refresh/process topology;
- D8 golden-flow execution;
- persistence schema/package layout;
- a generic integration/plugin/provider registry;
- a workflow/BPM engine;
- a universal ERP ontology.

Implementation remains blocked until D9.

---

## 2. Imported authority — not reopened here

B3 imports these accepted rules rather than redefining them:

1. **Business-System Materialization owns `Business Order Intent`, `Invoicing Intent`, native-result correlation and convergence meaning.** ERP-native TOP/document taxonomy is not MPC business semantics.
2. **Product master remains external.** MPC references the source-qualified Product; it does not create a Product master mirror.
3. **Availability Control owns Sellable Availability, Inventory Source/Scope and allocation semantics.** Native stock/reservation truth remains external.
4. **Commercial Economics owns Cost Basis and economic interpretation.** A native cost field/variant is evidence, not Cost Basis by inheritance.
5. **Marketplace Sales owns the canonical marketplace sale interpretation and transaction-specific Selling Entity attribution.** Materialization consumes that meaning rather than independently inferring marketplace/provider semantics.
6. **Post-Sale Resolution coordinates consequences; native business-system reversal/fiscal-result semantics remain external.** A pre-invoice commercial reversal is not automatically a post-invoice fiscal return.
7. **`SourceInstance` identifies one logical externally authoritative business-system/source namespace.** Credential rotation does not itself change SourceInstance; production and test/sandbox are distinct namespaces/environments.
8. **Consumer owns meaning; adapter owns protocol.** Provider DTO/status/operation vocabulary does not cross into business contexts.
9. **Known / absent / unknown / unavailable remain distinct.** Partial acquisition is not closure.
10. **Acceptance is not convergence.** Potentially accepted external writes are not blindly retried.
11. **Provider PII is minimized.**
12. **Direct Oracle/database is outside the target architecture and is not a fallback.**
13. **Integration Support, Provider Effective Capability/Requirement and Effective Business Capability remain separate authorities.**

No D0/D1/D2/D3 or D4-B1/B2 reopen is proposed by this candidate.

---

## 3. Proposed B3 governing invariant

> **For every Product 1.0 business-system fact or effect, MPC depends on a consumer-owned semantic contract whose external evidence/effect is Organization + SourceInstance qualified; the concrete Sankhya adapter may use only sanctioned provider operations and SourceInstance-specific bindings whose material assumptions are verifiable, preserves provider-native granularity/provenance/partiality, and reconciles consequential effects through authoritative reread/correlation. Sankhya-native topology and customer configuration never become MPC business ontology, and future substitutability does not authorize a generic ERP/workflow framework before a second real consumer proves the abstraction.**

Corollaries:

- Sankhya first **does not mean** Sankhya model.
- A SourceInstance binding **does not mean** a workflow engine.
- A provider inventory partition **does not mean** MPC `Lot`/batch identity or proven interchangeability.
- A TOP/code reference **does not prove** that its configured meaning remains valid after provider configuration drift.
- Provider configuration can be a prerequisite/capability fact; it does not become MPC policy merely because it influences execution.

---

## 4. Alternatives and proposed Global Maximum

### Alternative A — Sankhya-shaped core

Example failure shape:

```text
Business Order = TOP 313
Invoice = TOP 306
Marketplace = CODTIPVENDA 27
Lot = CONTROLE
```

**REJECT.** Current production already disproves this shape without needing another ERP: one SourceInstance contains materially different store and e-commerce document topologies, and TOP 313 represents e-commerce generally rather than Mercado Livre specifically.

### Alternative B — generic ERP/workflow framework now

Example failure shape:

```text
GenericERP
GenericResource
GenericCapabilityGraph
MaterializationWorkflow
  steps[]
  arbitrary operations/conditions
```

**REJECT.** There is one concrete business system today and no second consumer proving a stable shared provider ontology or workflow DSL. Such machinery would guess universals from Sankhya and create accidental complexity.

### Alternative C — provider-independent consumer semantics + concrete Sankhya adapter + bounded SourceInstance binding

**PROPOSED GLOBAL MAXIMUM.**

```text
D1 owner meaning
    ↓ consumer-owned semantic port
D4 Sankhya adapter
    - sanctioned operations
    - SourceInstance binding
    - provider requirements/capability
    - authoritative reread/correlation
    ↓
Sankhya
```

A future accepted business system adds another adapter behind the same consumer semantics where those meanings are genuinely equivalent. Shared technical machinery is extracted only after real repetition proves it; no generic provider/business graph is created for symmetry.

### Falsification test — workflow-engine drift

If MPC/domain-level configuration ever needs to express arbitrary **“how many native documents, in what sequence, under which provider-specific steps”**, this boundary has failed. B3 must return to decision rather than letting a binding evolve into a workflow/BPM engine.

Business domains care about their intents, provider requirements relevant to those intents, and convergence on authoritative external results — not an arbitrary provider document choreography.

---

## 5. Sankhya SourceInstance, environment and authentication

### 5.1 Environment qualification

Production and sandbox are distinct SourceInstance environments and cannot be collapsed because native keys/data appear similar.

Measured current auth token evidence establishes:

- `ambiente=hml` for sandbox;
- `ambiente=prd` for production.

A sandbox-only `environment` UUID was observed, while production exposed that claim as null. Therefore the UUID is **not** a cross-environment identity law.

Proposed contract:

1. configured SourceInstance binding identifies the intended environment/host;
2. provider-authoritative environment markers such as `ambiente`, when exposed, must match the binding or attribution fails closed;
3. data-plane responses need not carry a namespace marker on every resource; where none exists, the authenticated/configured SourceInstance binding is the control;
4. credentials are runtime secrets, never SourceInstance identity.

### 5.2 Authentication

Current target auth is:

```text
OAuth 2.0 client_credentials
+ X-Token
+ POST /authenticate
→ Bearer access token
```

The old login/appkey path is not target fallback.

Measured token TTL is 300 seconds. Refresh scheduling/locking/caching belongs to D7; B3 only requires that auth expiration/unavailability never become business absence.

MGECOM Bearer compatibility is empirically established across sanctioned services in both sandbox and production.

---

## 6. Product and Readiness evidence

### 6.1 Native Product identity

For the current Sankhya SourceInstance:

```text
SourceInstance + CODPROD
```

is the external Product reference required by D2.

`CODPROD` remains provider-native; it is not an MPC Product master identity.

### 6.2 Read surfaces

Prefer the dedicated Product REST surface when it supplies the requested semantic facts. Use sanctioned `CRUDServiceProvider.loadRecords` with `rootEntity=Produto` and a minimum explicit fieldset only when a real consumer requires native/custom facts absent from the dedicated REST resource.

This is a **per-claim** selection rule, not “REST always” or “loadRecords always”.

### 6.3 Identifier evidence

Measured Product evidence includes, where populated:

- native product key;
- active state;
- reference;
- supplier reference;
- brand;
- NCM/other fiscal product facts;
- alternate-volume barcode where applicable.

The current Product REST surface exposes **no universal first-class GTIN/EAN field**. An EAN-shaped `referencia` value is convention evidence only.

Therefore:

- `REFERENCIA == EAN` is forbidden as an identity/correspondence law;
- one identifier is insufficient by architecture fiat;
- D4 supplies available identifier evidence;
- Product & Channel Readiness owns whether the evidence is sufficient to establish correspondence.

Wire quirks such as provider field spelling/encoding differences remain entirely inside the adapter.

---

## 7. Company, location and inventory contract

### 7.1 External references remain external

Sankhya `CODEMP` and `CODLOCAL` are external/source-native references.

They do not automatically equal:

- Selling Entity;
- Inventory Source;
- Fulfillment Node.

Current business owners map/qualify those references under their accepted semantics. Transaction-specific Selling Entity attribution continues to come from Marketplace Sales, not from D4 guessing `CODEMP`.

### 7.2 Inventory observation granularity

The dedicated stock REST surface is useful but lossy. Real measurement established that, for the observed SourceInstance:

```text
REST estoque = native ESTOQUE - RESERVADO
```

at provider dimensions including product/company/location/control where present.

It does not expose the decomposition needed to distinguish, for example, no stock from fully reserved stock.

For claims requiring reservation/control decomposition, the minimum sanctioned entity read is `Estoque` with only materially required fields such as:

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

No ERP mirror is introduced.

### 7.3 Dedicated-stock limitations remain part of coverage

The dedicated REST stock resource is not universal inventory authority merely because it is convenient. Current documented/measured limitations include provider-defined exclusions/absence semantics such as WMS/third-party scope and never-moved products not appearing as zero.

Therefore resource absence cannot automatically become `stock=0`.

### 7.4 Provider control partition

`CONTROLE` is preserved as an **opaque provider inventory partition value** when present.

It is **not** canonically named `Lot`, `Batch`, `Tonality`, `Serial` or another MPC entity/identity.

Current production evidence materially shows:

- control is configured **per Product**, not reliably by product family;
- 3,038 active controlled products were observed in the SourceInstance;
- controlled products are concentrated in floor/covering families, including porcelain, but the same families contain uncontrolled counterexamples;
- this SourceInstance currently uses one observed raw control-type code (`TIPCONTEST='I'`), so no multi-type MPC control taxonomy is justified;
- `CONTROLE` can carry values such as physical-looking lot strings **and** `ENCOMENDA`, proving the field itself does not encode one universal business meaning;
- no sanctioned provider entity was found exposing structured tonalidade/calibre/grade attributes for a `CONTROLE` value.

### 7.5 Partition-aware availability fence

Measured evidence proves aggregate free quantity can exceed every individual provider partition. Therefore:

> **Availability Control must not collapse provider-native inventory partitions before interchangeability/satisfaction semantics are sufficiently established for the consequential decision.**

This rule does not decide that partitions are always non-interchangeable. It preserves uncertainty instead of fabricating aggregate sellability.

Current Mercado Livre **sold** population evidence is control-free: all observed e-commerce items had empty `CONTROLE` because their Products were actually uncontrolled (`TIPCONTEST='N'`). The currently **listed** Mercado Livre population was not cross-matched in the final control sweep and remains unknown for this property.

Proposed first Product 1.0 capability fence:

- uncontrolled Product lane: supported subject to ordinary Availability rules;
- controlled Product automatic marketplace availability/materialization lane: **not claimed** until business interchangeability and required selection timing are proven for that operating path;
- controlled-path unknown is explicit `unsupported` / `external-required` / operator work as appropriate, never silent aggregate availability or adapter-chosen partition.

This does not reopen D1. Availability already owns Sellable Availability and evidence sufficiency. A future real controlled-marketplace path is the reopen/adjudication trigger if existing ownership proves insufficient.

---

## 8. Cost and tax evidence

### 8.1 Cost observations

No dedicated current cost REST contract was established. Sanctioned entity read `rootEntity=Custo` is proven live and exposes externally governed observations qualified by Product, company, effective/as-of time and provider-local dimensions/provenance where populated.

Multiple native cost variants are real external observations. B3 does **not** select one as `Cost Basis`.

Commercial Economics receives typed/provenanced Cost Observations and alone owns the Cost Basis used for L0/L1/L2 economic reasoning.

Provider sentinel rows such as zero-cost historical sentinel dates are not silently promoted to economic facts.

### 8.2 Expected-tax calculation

`POST /v1/fiscal/impostos/calculo` is the sanctioned candidate for expected-tax evidence because it applies the business system's configured tax rules rather than requiring MPC to copy a tax-rule engine.

Integration Support is established, but Provider Effective Capability is currently **conditioned** by SourceInstance configuration/customization. A real sandbox call was intercepted by instance customization requiring data the exercised calculation setup did not satisfy; no persistence residue remained.

Proposed decision:

- prefer sanctioned Sankhya tax calculation over copying fiscal-rule tables/logic into MPC;
- do not claim Expected-Tax capability merely because the endpoint exists;
- B3 acceptance must retain an explicit **Tax Evidence Gate**: the selected production SourceInstance calculation path/model/configuration must be proven semantically sufficient and operationally usable for Product 1.0, or B3 returns `STOP / SPLIT PREREQUISITE` for that capability;
- Direct Oracle/TGFICM-style rule copying is not fallback.

Realized tax/fiscal evidence comes from authoritative native fiscal results, not from the expected-tax simulation.

---

## 9. Coverage, pagination and change semantics

### 9.1 Full/scoped observation is the correctness baseline

Current sanctioned enumeration is operationally plausible at the measured SourceInstance scale. B3 does not require delta acquisition to make facts correct.

Per-operation coverage remains explicit:

- point read proves only that exact source-qualified object/scope;
- enumeration proves only the traversed source-defined scope and only after all required pages complete;
- no source-defined snapshot isolation is invented;
- partial traversal is partial.

### 9.2 Delta is optional and prerequisite-bound

`modifiedSince`/change-log behavior depends on Sankhya source configuration. Current measurement observed the change log disabled.

Materially, failure honesty differs by resource:

- some operations return an attributable error when logging is unavailable;
- Product delta can return a no-record/404 shape that is indistinguishable from a legitimate empty delta if the consumer ignores the prerequisite.

Therefore:

> **No delta result is accepted as completeness evidence until the required source change-log coverage is independently established.**

Even when enabled, the documented short retention window means delta is not historical authority or indefinite recovery.

B3 does **not** require enabling `LOGTABOPER` now. Full/scoped acquisition remains the correctness path; D7 later decides whether delta is operationally justified, its cadence and recovery strategy.

### 9.3 Pagination/provider quirks

Provider-specific page origins, page sizes, `hasMore`, 404-as-end behavior, date formats and malformed/provider-specific field names stay in the Sankhya adapter.

A short page is not automatically end-of-data when the provider indicates more results.

### 9.4 Operational viability

Measured behavior is sufficient to treat the sanctioned Gateway as operationally plausible for Product 1.0 design:

- ordinary reads/loadRecords are bounded and usable;
- order enumerations are materially heavier and PII-rich, so they should be scoped/minimized rather than indiscriminately mirrored;
- no rate-limit headers were observed, so concurrency/rate ceilings remain Unknown;
- no aggressive stress test is justified in D4.

D7 owns token refresh, request scheduling, concurrency, backoff, cursors and runtime observability.

---

## 10. Native order observation

### 10.1 Native order identity and point authority

The native business-order result remains:

```text
SourceInstance + native NUNOTA
```

No synthetic `MPCOrder` alias is introduced merely to normalize Sankhya.

The current REST order enumeration is usable for bounded company/window observation but its documented point filters were empirically unreliable for existing orders.

Therefore the current authoritative point/reread surface for the consequential native-order state is the sanctioned `CabecalhoNota` entity read by NUNOTA, with minimum required fields.

This is an adapter decision, not permission for business contexts to read `CabecalhoNota` or other provider entity vocabulary directly.

### 10.2 Confirmation and pendency are distinct facts

Measured provider semantics establish:

```text
STATUSNOTA='A' → created / unconfirmed
STATUSNOTA='L' → confirmed
```

while `PENDENTE` is orthogonal fulfillment/billing pendency. Confirmed + pending combinations are real.

The adapter therefore must not collapse:

```text
created
confirmed
pending/fulfilled
fiscalized
```

into a single boolean lifecycle.

Raw status letters remain provider-local; Materialization consumes semantic evidence such as native order exists / confirmation established / remaining materialization pendency.

### 10.3 `CACSP.confirmarNota`

The current official public reference does not document a confirmation operation, but the provider-native service `CACSP.confirmarNota` was empirically established through the sanctioned MGECOM Gateway in sandbox and production:

- requests targeted the caller-known NUNOTA;
- a missing SourceInstance-required carrier produced a clean provider/business-rule rejection without state change;
- after satisfying that prerequisite, reread proved `A → L` convergence;
- OAuth/Bearer MGECOM transport succeeded;
- no direct database/status mutation was used.

Therefore current Integration Support for confirmation is **SUPPORTED in the observed SourceInstance/runtime**, while the documentation gap remains provider evidence/reopen risk rather than reason to pretend the operation is officially documented.

Confirmation errors can represent business/provider prerequisites, rejection, pending/release behavior or ambiguity; transport success is not the confirmation claim. Authoritative reread owns convergence.

No claim is made that every future Sankhya environment/version supports this undocumented service identically; SourceInstance capability must remain re-verifiable.

---

## 11. Customer/partner prerequisite

Current e-commerce orders prove that native materialization uses real distinct customer/partner references rather than one generic marketplace partner.

B3 therefore requires a bounded native-customer prerequisite contract:

1. Marketplace Sales supplies the marketplace sale/buyer facts legitimately required for materialization;
2. D4 may use sanctioned Sankhya customer lookup/create/update capability where required;
3. Materialization must establish a source-native partner reference sufficient for the native order before consequential order creation;
4. raw buyer PII is minimized to the fields genuinely required by Sankhya/business/fiscal processing;
5. native customer/partner identity remains external/source-qualified; MPC does not create a Customer Master domain or infer canonical customer identity from mutable names/emails/documents by convenience.

Exact create-vs-match policy and required buyer fields must be proven against the selected marketplace/business-system lane before D8. Duplicate/ambiguous native customer correspondence is explicit exception work rather than silent guessed matching.

---

## 12. Business Order Intent materialization

### 12.1 Semantic contract

Materialization owns the MPC intent and convergence meaning:

```text
Business Order Intent
  + Organization
  + SourceInstance
  + canonical marketplace Sale context
  + transaction-specific Selling Entity attribution
  + required product/quantity/customer/materialization facts
→ attempt native materialization
→ source-qualified native business-order result(s)
→ authoritative convergence / rejection / pending / ambiguity
```

The domain does **not** own Sankhya TOPs, series, TIPMOV, CACSP operations, NUNOTA status letters or TGFVAR.

A Business Order Intent may require provider-native intermediate artifacts internally. Their number and sequence are adapter concerns and do not become configurable MPC workflow steps.

### 12.2 Selected current Metal Nobre e-commerce binding

Production evidence establishes the current e-commerce lane:

```text
create native e-commerce order
  current binding: TOP 313
  current order series: PA
  current marketplace discriminator for Mercado Livre: CODTIPVENDA 27
→ confirm native order
→ authoritative reread
→ native business order converged
```

Material facts of current TOP 313 include order movement + reservation/financial behavior under the currently effective, version-qualified TOP configuration.

**These values are a time-bound SourceInstance binding, not MPC semantics.**

TOP 313 is e-commerce generally, not Mercado Livre. Other e-commerce negotiation types exist on the same TOP. The Marketplace Sales/Installation context selects the correct marketplace-specific binding evidence; the adapter does not infer provider identity from TOP alone.

### 12.3 Creation surface

For the selected current SourceInstance binding, `CACSP.incluirNota` is the proven sanctioned MGECOM creation surface and can express the real native operation without requiring the currently unresolved REST `notaModelo` setup.

The v1 `POST /v1/vendas/pedidos` capability remains provider-conditioned rather than globally rejected: it may be usable with a correctly configured native model, but B3 does not prefer it merely because it is newer REST if it cannot currently satisfy the selected binding correctly.

`CACSP.incluirNota` partial update behavior was also measured, but B3 does not generalize this into an arbitrary provider patch API. Only consumer-required sanctioned update use is admitted.

### 12.4 Required input binding, not a bag of knobs

Current documents expose varying company/vendor/nature/carrier/customer values. The target MUST NOT turn every observed field into a global config setting.

For each selected materialization operation, the adapter binding includes only values/derivations required to perform that operation correctly, sourced explicitly from one of:

- stable SourceInstance configuration;
- current domain-owned intent/context;
- externally governed/provider-derived prerequisite/default.

A binding value never becomes business policy merely because Sankhya requires it.

---

## 13. SourceInstance binding validation and configuration drift

Current `TipoOperacao` evidence is version-qualified by `DHALTER`; `CODTIPOPER` alone is not eternal meaning.

Sanctioned reads expose material native properties including movement class, stock effect, financial effect, pendency/lot behavior, fiscal-model posture, activity and version/effective timestamp.

Therefore a consequential binding MUST NOT mean only:

```text
orderTop = 313
invoiceTop = 306
```

It must be able to establish that the referenced current provider configuration still has the properties on which the selected integration contract depends.

For the current e-commerce lane this includes, where material to the consumer claim, evidence such as:

- order vs fiscal movement role;
- reservation vs write-down effect;
- financial generation/reversal posture;
- confirmation/provider prerequisites;
- fiscal-model posture;
- active/current version.

If the referenced provider operation/configuration drifts so the expected properties are no longer established, the binding becomes invalid/unknown and consequential execution fails closed until re-adjudicated/reconfigured.

B3 freezes this **semantic validation obligation**, not its runtime polling/cache mechanism. D7 decides when/how configuration is revalidated and cached.

---

## 14. Invoicing Intent materialization

### 14.1 Semantic contract

```text
Invoicing Intent
  + readiness-gated native business-order result
  + current Fulfillment-owned physical-readiness/conference meaning
  + required provider/business-system prerequisites
→ sanctioned native fiscal progression
→ source-qualified fiscal/document result
→ authoritative reread + origin/result correlation
→ converged / rejected / pending / ambiguous
```

Physical readiness remains Fulfillment authority. Materialization never invoices merely because a provider order exists.

### 14.2 Current e-commerce binding evidence

Production history proves the current SourceInstance e-commerce progression:

```text
TOP 313 native order
  → TOP 306 native fiscal result
```

with distinct native NUNOTAs and line/quantity correlation.

The current TOP 306 configuration is fiscal/write-down shaped and its current document examples carry fiscal numbering/status evidence. Exact TOP/series/config values remain SourceInstance-local.

### 14.3 Effect surface

`SelecaoDocumentoSP.faturar` is the selected sanctioned transform/invoicing surface:

- documented to require an eligible confirmed source document and provider configuration;
- empirically proven with OAuth/Bearer in production on a non-fiscal order→order transformation;
- existing real SourceInstance evidence proves 313→306 order→fiscal results and TGFVAR correlation;
- partial/quantity-level correlation is representable by provider-native relation evidence.

The final selected-lane **controlled fiscal write** remains a D8 effect proof because executing a production fiscal document has materially higher irreversible/legal cost than the architectural signal needed in B3. B3 does not claim D8 has passed.

Whether the generated 306 requires a separately visible post-create confirmation transition is not inferred from steady-state history; fiscal convergence is proven by the authoritative native fiscal reread/result state rather than by guessing internal provider steps.

---

## 15. Native result correlation

The adapter preserves source-native correlation without creating an MPC universal provider graph.

Current supported surfaces include:

- native fiscal result carrying source-order correlation where exposed;
- sanctioned provider relation entity `CompraVendavariosPedido` / TGFVAR semantics for origin/result line and `QTDATENDIDA` quantity correlation.

Requirements:

1. original order identity and result document identity remain distinct source-qualified references;
2. 0..N/partial result relationships remain representable when the provider exposes them;
3. item/quantity scope is preserved when material;
4. provider relation resources remain protocol/evidence, not MPC entities;
5. a 2xx transform response alone is not convergence.

Commercial Economics and Post-Sale may consume the resulting Materialization-owned interpretation/evidence through accepted semantic boundaries; they do not read TGFVAR directly.

---

## 16. Pre-invoice commercial reversal vs post-invoice fiscal return

Production evidence establishes a current SourceInstance e-commerce branch in which uninvoiced TOP-313 orders correlate to TOP-307 results.

Observed semantics:

- the 307 results originate from uninvoiced 313 orders;
- they reverse commercial/financial pendency without a stock write-down;
- no observed 307 originated from a 306 fiscal result.

Therefore:

> **Pre-invoice native commercial reversal is not the same business-system consequence as post-invoice fiscal return/reversal.**

B3 may expose/translate the observed native reversal result to the consuming Post-Sale/Materialization semantics, but the command that produces the 307-class result is **not yet proven** and MUST NOT be invented by equating it with another cancel endpoint.

Current actuation classification:

- observed pre-invoice reversal result: supported as read/correlation evidence;
- write/command to produce the same consequence: Unknown / external-required until proven;
- post-invoice fiscal return/reversal path: Unknown in the observed e-commerce history.

This does not block Product 1.0 from coordinating post-sale work: unsupported business-system consequences remain explicit external-required Work rather than silent success. A later real requirement/controlled proof may extend support without creating a generic reversal model.

---

## 17. External-effect semantics for Sankhya writes

Every consequential write admitted by the selected binding obeys D3/D4-B1:

1. target Organization + SourceInstance is explicit;
2. Materialization/Post-Sale domain intent/correlation anchor is explicit;
3. current required binding/provider prerequisites are established;
4. the request is scoped only to the intended native effect;
5. response is classified no stronger than provider evidence permits;
6. accepted/rejected/pending/ambiguous remain distinguishable where reachable;
7. authoritative reread/correlation follows consequential acceptance;
8. timeout/connection loss after possible acceptance is **not blindly retried**;
9. native business-rule/custom-trigger/liberação failures become translated provider requirement/rejection evidence, not raw ORA/provider implementation leakage to business domains;
10. a provider write never bypasses current Readiness, Fulfillment, Governance or other domain validity merely because D4 supports the protocol.

Measured confirmation proved that SourceInstance custom rules can reject a write after creation succeeded. Therefore validating only the creation payload is insufficient evidence for full materialization capability.

---

## 18. Current Product 1.0 Sankhya proof lane

The current first business-system proof context is deliberately narrow and time-bound.

### 18.1 Read/fact lane

- Sankhya production SourceInstance via sanctioned Gateway;
- Product native key = CODPROD;
- current Product/identifier facts as measured;
- company/location external refs explicitly qualified;
- inventory via dedicated net surface and minimum entity decomposition where the consumer needs reservation/control detail;
- cost via sanctioned `Custo` observations;
- expected tax subject to the explicit Tax Evidence Gate.

### 18.2 Marketplace business-order lane

Current Metal Nobre production evidence:

```text
Marketplace Sale
  → e-commerce materialization binding
  → native TOP-313 order
  → confirmation
  → authoritative reread
```

For Mercado Livre, current SourceInstance evidence distinguishes that lane with negotiation-type binding `27`; TOP 313 alone is not the marketplace identity.

### 18.3 Fiscal lane

```text
confirmed native e-commerce order
  + Fulfillment readiness
  → Invoicing Intent
  → sanctioned faturamento
  → native TOP-306 fiscal result
  → authoritative reread + correlation
```

The production fiscal write itself remains D8 controlled proof; existing production history proves the native result topology/correlation.

### 18.4 Inventory-control support fence

Current observed Mercado Livre sold population is uncontrolled. Controlled marketplace products are **not** included in the currently claimed automatic Availability/Materialization proof lane until partition business semantics are established.

---

## 19. Explicit residual gates and defers

### Gate G1 — Expected Tax

**B3 closure gate.** Prove the sanctioned expected-tax calculation is usable and semantically sufficient for the selected production SourceInstance configuration, or return `STOP / SPLIT PREREQUISITE` for L0 Expected Economics. Do not copy tax rules through Oracle as fallback.

### Gate G2 — Native customer/partner prerequisite

**B3/D8 proof obligation.** Establish the concrete minimal source-native customer resolution/materialization path for a real marketplace sale, including duplicate/ambiguity and PII-minimization behavior. Customer API existence alone is not convergence proof.

### Gate G3 — First selected-lane fiscal effect

**DEFER SAFELY → D8.** Controlled actual 313→306 external effect plus authoritative fiscal reread/correlation. The command/topology is sufficiently grounded for architecture; D8 proves the irreversible effect.

### Gate G4 — Controlled-product marketplace lane

**DEFER SAFELY / not a current Product 1.0 proof lane.** Before automatic marketplace availability/materialization for a controlled Product, establish business interchangeability/partition-feasibility and when a concrete partition must be selected. No adapter-chosen lot and no aggregate guess.

### Gate G5 — Post-invoice business-system return/reversal

**DEFER SAFELY as explicit external-required consequence unless selected Product 1.0 golden flow demands automated actuation.** Current observed e-commerce history does not establish the fiscal-return path or command. Post-Sale remains capable of explicit Work/closure coordination without pretending automated business-system reversal support.

### D7 defers

- token refresh scheduling/locking;
- polling/full/delta cadence;
- cursor/checkpoint storage;
- rate/concurrency/backoff;
- cache of provider configuration/binding validation;
- worker/process/deployment topology.

None of these defers moves business authority.

---

## 20. YAGNI / explicit non-goals

B3 MUST NOT introduce merely for abstraction or anticipated providers:

- generic `ERP` business entity;
- universal `ERPAdapter` interface containing every possible ERP operation;
- generic provider/resource/operation/capability graph;
- plugin registry/factory/self-registration framework;
- arbitrary workflow/materialization DSL;
- MPC `TOP`, `NUNOTA`, `TipoNegociacao`, `TGFVAR`, `CONTROLE`, Sankhya status or Sankhya document entity;
- universal `Lot`/Batch/Serial entity from `CONTROLE`;
- replicated Sankhya product/stock/cost/tax database mirror;
- arbitrary SQL/DbExplorer access disguised as Gateway integration;
- duplicated tax engine from Sankhya rule tables;
- support for every Metal Nobre non-marketplace commercial process;
- speculative TOTVS/Bling adapters;
- claim that every Sankhya customer uses the same TOP/document workflow;
- family-level inference of product control semantics;
- blanket configuration knob for every provider field seen in historical documents.

The store-process 14→303→305 lane is useful as variability/counterexample evidence; it is not a Product 1.0 MPC workflow requirement merely because it exists in the SourceInstance.

---

## 21. Provider-independent replacement test

A future accepted business system may replace Sankhya for the same organization or be used by another Organization.

The following MPC meanings should survive unchanged where the future system genuinely supplies the same semantics:

- source-qualified Product reference;
- Product/identifier facts for Readiness;
- qualified inventory observations and honest partition/granularity semantics;
- Cost Observations;
- expected/realized fiscal evidence requirements;
- Business Order Intent;
- native business-order convergence;
- Invoicing Intent;
- native fiscal convergence;
- native origin/result line/quantity correlation;
- explicit commercial reversal vs fiscal-return consequence distinction;
- honest unknown/unavailable;
- authoritative reread;
- no blind retry of ambiguous effects.

The following Sankhya pieces should be replaceable without changing those business meanings:

- TOPs and their versioning;
- NUNOTA;
- series;
- `STATUSNOTA` / `PENDENTE` provider vocabulary;
- `CACSP` / `SelecaoDocumentoSP` services;
- TGFVAR/provider relation entities;
- `TIPCONTEST` / `CONTROLE` representation;
- Sankhya auth/Gateway;
- customer-specific triggers/liberações;
- all Metal Nobre binding values.

If adding another concrete system requires changing an accepted MPC semantic because the business meaning is genuinely different, reopen the responsible semantic decision. Do not contort the new system into a false common ERP model.

Marketplace Installation and business-system SourceInstance remain distinct concepts; symmetry between adapter families does not justify a generic `IntegrationInstance` identity.

---

## 22. Proposed B3 outcome

Subject to independent review and operator adjudication, the proposed B3 decision is:

### `CURRENT STRUCTURE CONFIRMED` — business-system boundary

The accepted D1/D2/D3/D4-B1 structure remains sound:

- consumer-owned semantic ports;
- SourceInstance-qualified external evidence/results;
- concrete provider adapter;
- no Integration business domain;
- no generic ERP/workflow framework.

### `CURRENT STRUCTURE CONFIRMED` — sanctioned Sankhya target

The sanctioned Gateway/API surface is sufficient to define the current Product 1.0 fact/materialization contract without Direct Oracle fallback, subject to explicit capability gates rather than invented support.

### `DEFER SAFELY`

- runtime/delta/cursor/concurrency mechanics → D7;
- first irreversible selected-lane fiscal write → D8;
- controlled-product marketplace automation until real evidence exists;
- unproven post-invoice fiscal-return actuation unless a selected golden flow makes it mandatory.

### `STOP / SPLIT PREREQUISITE` trigger

If the expected-tax gate, selected native customer prerequisite, or another materially required Product 1.0 claim cannot be satisfied correctly through sanctioned SourceInstance operations, stop/re-adjudicate that capability. Direct Oracle/database is never admitted implicitly.

No D0/D1/D2/D3 reopen is proposed now.

---

## 23. Proof strategy before implementation

B3 architecture claims are falsified/proven proportionally:

### Already evidenced

- current OAuth/SourceInstance environment binding behavior;
- Product native key and identifier limitations;
- company/location stock granularity;
- REST-net vs reservation decomposition;
- controlled-product partition evidence;
- cost entity observations;
- change-log prerequisite/failure-honesty behavior;
- native order point-read fallback;
- `CACSP.incluirNota` creation;
- `CACSP.confirmarNota` confirmation + rejection + reread;
- MGECOM OAuth compatibility;
- `SelecaoDocumentoSP.faturar` effect capability;
- origin/result line/quantity correlation;
- production e-commerce topology 313→306;
- observed 313→307 pre-invoice reversal branch;
- binding-property/version reads;
- current control-free marketplace sold lane;
- production-vs-sandbox materialization divergence.

### Still required for the named claim

- G1 expected-tax selected SourceInstance proof;
- G2 concrete customer/partner prerequisite proof;
- D8 controlled selected-lane fiscal effect;
- any later claim of automated controlled-product marketplace support;
- any later claim of automated post-invoice business-system fiscal return.

A green mock/test cannot substitute for the real external dependency claim.

---

## 24. Reopen triggers

Reopen only the implicated decision when evidence shows, for example:

- a second real business system cannot implement the accepted consumer meaning without forcing Sankhya-specific semantics into the domain;
- a selected current Sankhya provider operation disappears/changes materially and sanctioned replacement is not sufficient;
- SourceInstance binding properties cannot be validated strongly enough to protect consequential execution;
- a controlled Product enters the marketplace operating scope and its interchangeability/selection semantics cannot fit existing Availability/Fulfillment ownership cleanly;
- business-system customer correspondence requires a genuinely independent MPC business lifecycle rather than a bounded materialization prerequisite;
- post-invoice return/reversal becomes a claimed automated normal path and current external-required treatment is insufficient;
- Product 1.0 requires a business-system fact/effect for which only Direct Oracle/arbitrary SQL could satisfy the current claim;
- repeated second-provider evidence proves a shared adapter mechanism/abstraction materially reduces total complexity — at that point extract the smallest common mechanism rather than forecasting it now.

Do not reopen for naming preference, abstract symmetry or hypothetical providers.

---

## 25. Independent-review attack surface

Reviewer should reconstruct authority independently and challenge at least:

1. **Sankhya leakage:** Does any proposed domain contract still encode TOP/NUNOTA/CACSP/TGFVAR/CONTROLE semantics by disguise?
2. **Workflow-engine drift:** Does the SourceInstance binding secretly require a configurable provider step graph? Can the current e-commerce binding remain bounded without supporting the store lane as a generic workflow?
3. **Future business-system seam:** Does the replacement test actually preserve the domain if a real TOTVS/Bling-like system arrives, without already building a universal ERP abstraction?
4. **Binding vs policy:** Are provider flags/configuration being misread as MPC business policy?
5. **Binding drift:** Is version/property validation sufficient as an architecture obligation, or is the proposed binding too weak/too complex?
6. **Inventory partitions:** Does B3 preserve enough granularity for Availability without inventing a `Lot` entity or incorrectly aggregating partitions?
7. **Controlled-product defer:** Is it legitimate to keep the first marketplace proof lane control-free, or does D0 require controlled Product support now?
8. **Tax gate:** Is expected-tax capability correctly classified, or is B3 attempting to close while a material Product 1.0 economics prerequisite is unresolved?
9. **Customer prerequisite:** Is native partner resolution truly a bounded Materialization prerequisite, or does it expose a missing D1 responsibility?
10. **Reversal:** Is observed pre-invoice 307 evidence correctly separated from post-invoice fiscal return, and is external-required actuation sufficient for current Product 1.0?
11. **Point/read authority:** Is sanctioned `loadRecords` use bounded entity access, or has it become an ERP mirror/database abstraction by another name?
12. **Operational viability:** Do the measured latency/coverage/token facts expose a B3 correctness issue rather than a D7 runtime defer?
13. **Sandbox divergence:** Does the observed sandbox featurelock vs production success invalidate any current proof claim?
14. **PII:** Does order/customer acquisition minimize real buyer PII sufficiently for the required materialization/fiscal purpose?
15. **No hidden Oracle:** Does any gap tempt an arbitrary SQL/DB fallback contrary to accepted B1?

Reviewer findings are evidence, not authority. Any finding that creates a new business requirement or moves authority must return to adjudication rather than entering B3 as a disguised correction.

---

## 26. Candidate disposition

This file exists only to enable fresh independent challenge.

It MUST be deleted after review/adjudication and before canonical B3 consolidation.

Canonical status remains whatever `docs/engineering/rebaseline/README.md` says. This candidate does not open/accept B3 by existence and does not authorize B4, D5+, implementation, merge or product writes.
