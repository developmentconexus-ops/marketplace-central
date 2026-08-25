# D6-R2 P5 — Complete Screen / Material-Surface Inventory

> **Status:** DERIVED / CANDIDATE — P5 complete; no P8 block is operator-`LOCKED`
> **Parent:** [D6-R2 Complete Frontend Realization Closure](D6-R2-COMPLETE-FRONTEND-REALIZATION-CLOSURE.md)
> **Methods:** [DevelopmentConexus Engineering Method v1.0.0](../../development/engineering-method.md) + [Frontend Product Experience Planning Method v2.3](../../development/frontend-product-experience-planning-method.md)
> **Inputs:** P0–P3 N01–N16 / UF01–UF16 / 99-operation coverage + P4 revalidated candidate IA
> **Boundary:** screen/surface inventory only; no rendered wireframes, no Product implementation, no P9 wire binding yet

## 1. P5 law

P5 derives human-operable surfaces from accepted user flows and the revalidated IA. It does **not** create one page per endpoint and does not preserve the prior D6-B1 39-state count by symmetry.

A material surface is separated only when one of these materially changes:

```text
primary semantic truth
safe user action / write owner
identity qualification
concurrency / idempotency behavior
security / disclosure context
knowledge / exactness guarantee
recovery path
editor vs viewer mode
```

Cosmetic layout, loading style, component boundaries and responsive rearrangement are not reasons to create a separate material surface.

## 2. P5 structural result

The accepted D6-B1 inventory contained **39 interaction states**. P5 derives the smaller current structural baseline:

```text
persistent global shell                          1
candidate content route/page homes              35
material Product-operation coverage             99 / 99
new Product operations                           0
new top-level navigation destinations            0
mandatory material modal/drawer surfaces         0
mandatory alternate-view frameworks              0
```

The count is descriptive, never an implementation target.

The three-state reduction from the old 39-state evidence is intentional:

- old `S21 Publicação / Operação` + `S22 Publicação / Performance` become one Publication detail route with separate Offering and Performance material regions;
- old `S71 Venda / detalhe composto` + `S72 Pedidos ERP / materialização` + `S73 Faturamento` become one Sale detail route with separate Sales, Economics, Materialization, Fulfillment and Post-Sale owner-specific regions.

This is composition, not ownership collapse. Each region retains its own Product owner, Permission and later P9 wire binding.

## 3. Global shell — persistent structural surface, not a content route

### G00 — App Shell / global frame

**Purpose:** preserve the current human access context and make the accepted IA findable without becoming a second authorization layer.

Material regions:

- `G00-A` current Organization presentation + explicit Organization switch;
- `G00-B` grouped primary navigation (`Visão geral`, `Operações`, `Estratégia e Inteligência`, `Controle`, `Configurações`);
- `G00-C` page-local context host for exact/all-or-exact Marketplace Installation selection where the destination requires it;
- `G00-D` access/deep-link blocked states.

Laws:

- Organization is the only global workspace;
- Marketplace Installation is **not** a second global tenant selector and never supplies ambient authority;
- permission-conditioned navigation reduces noise only; server authorization remains authoritative;
- Organization switch invalidates incompatible server/navigation state;
- exact-required pages render explicit missing/invalid Installation states instead of selecting a first/default account.

No material drawer/modal belongs to the shell baseline.

## 4. Candidate route/page inventory

Primary route-family authority remains accepted D6. The more specific P5 suffixes below are **candidate navigation identities inside those accepted families**. P9 must bind exact route params/search-state carriers to the canonical identity inputs; P5 does not change Product identity semantics.

| ID | Candidate route/page home | Material surfaces inside the home | Need / flow coverage | P5 disposition |
| --- | --- | --- | --- | --- |
| R01 | `/visao-geral` | installation posture; bounded economics summary; Work preview; per-Installation Performance entry point | N01/N05/N06/N12 · UF01/UF05/UF13 | **KEEP** small read-only composition; A04 controls later hierarchy |
| R10 | `/preparacao` | source-product search/results; selected source identity; marketplace publication requirements + source values/evidence; correspondence resolution; continuation to ListingIntent | N03 · UF03/UF04 | **KEEP** one search-first workspace; real structural ambiguity triggers P6 |
| R20 | `/publicacoes` | Marketplace Listing collection; exact Installation context; explicit knowledge/status presentation | N04 · UF04 | **KEEP** collection home |
| R21 | `/publicacoes/:nativeListingKey` + exact Installation URL state | Offering truth; availability context; Performance material region; optional Market/Economics evidence; related Work/intents | N04/N05/N06/N07 · UF04/UF06/UF08/UF09 | **KEEP** one source-qualified Listing detail; old S22 absorbed as separate material region |
| R22 | `/publicacoes/intencoes` | ListingIntent collection/work context | N04 · UF04 | **KEEP** local Publicações subcontext |
| R23 | `/publicacoes/intencoes/:listingIntentId` | identity/revision header; draft editor; media; contextual evidence; submit/discard controls; consequential outcome/recovery | N04/N06 · UF04/UF08 | **KEEP** dedicated material editor; P6 triggered |
| R24 | `/publicacoes/precos` | PriceIntent collection; selected PriceIntent detail; create PriceIntent action; links back to Market/Economics evidence | N04/N06 · UF04/UF08 | **KEEP** one PriceIntent workspace; selected detail need not become a separate page by default |
| R30 | `/disponibilidade` | sellable-availability collection; selected current state/provenance; unknown/zero distinction; configuration continuation | N07 · UF09 | **KEEP** operational read home; no quantity editor/sync action |
| R40 | `/performance/resumo` | exact Installation + period context; traffic/sales/media summary evidence; comparability/coverage/provenance states | N05 · UF05 | **KEEP** strategic summary; P6 triggered for evidence hierarchy |
| R41 | `/performance/publicacoes` | listing-performance collection; selected listing evidence; unknown/unavailable Listings retained; link to R21 | N05 · UF06 | **KEEP** analytical collection |
| R42 | `/performance/midia` | scoped Retail Media evidence; campaign/listing/catalog/family qualification; measurement basis/coverage | N05 · UF07 | **KEEP** read-only media surface; no Ads controls |
| R50 | `/mercado` | explicit subject context; competitive position; comparable offers; continuation to Economics/PriceIntent | N06 · UF08 | **KEEP** evidence-only strategic lens |
| R60 | `/economia/prevista` | expected-economics collection/subject detail; Cost Basis/provenance; stateless price-scenario evaluation; PriceIntent handoff | N06 · UF08 | **KEEP** analysis page; scenario evaluation is not mutation |
| R61 | `/economia/realizada` | sale-economics collection; economic-performance summary; link to source-qualified Sale | N06/N16 · UF08/UF16 | **KEEP** realized-economics page |
| R62 | `/economia/reconciliacao` | attribution collection; selected attribution detail; bounded resolve action; current/recovery state | N16 · UF16 | **KEEP** one reconciliation workspace; selected attribution stays deep-linkable URL state without new top-level route |
| R70 | `/vendas` | source-qualified Sale collection; exact Installation context | N08 · UF10 | **KEEP** primary Sale entry |
| R71 | `/vendas/:nativeSaleKey` + exact Installation URL state | Sale truth; Selling Entity attribution; Sale economics; BusinessOrder/Party/destination materialization; invoicing; fulfillment summary; post-sale summary; Work links | N08/N09/N11/N12/N16 · UF10–UF13/UF16 | **KEEP** one cross-owner Sale detail with owner-separated material regions; P6 triggered |
| R80 | `/expedicao/execucoes` | actionable FulfillmentExecution queue; context/status qualification | N09 · UF11 | **KEEP** operational queue |
| R81 | `/expedicao/execucoes/:fulfillmentExecutionId` | execution truth; physical checkpoint controls; fiscal/provider prerequisite status; artifacts; blocked/exception state; Sale/Work links | N09 · UF11 | **KEEP** dedicated high-risk execution page; P6 triggered |
| R82 | `/expedicao/envios` | source-qualified Shipment collection; exact Installation context | N10 · UF11 | **KEEP** local Expedição collection |
| R83 | `/expedicao/envios/:nativeShipmentKey` + exact Installation URL state | Shipment current evidence; material delivery exception; related Sale/Work | N10/N12 · UF11/UF13 | **KEEP** Shipment detail |
| R90 | `/pos-venda` | resolution collection; create-resolution material action; explicit current/empty state | N11 · UF12 | **KEEP** collection + bounded creation surface |
| R91 | `/pos-venda/:postSaleResolutionId` | resolution truth; related owner evidence; Work/economic consequence links | N11/N12 · UF12/UF13 | **KEEP** resolution detail; no generic close/refund/cancel |
| R100 | `/trabalho` | Work queue; assignment/escalation indicators; source-owner qualification | N12 · UF13 | **KEEP** coordination queue |
| R101 | `/trabalho/:workId` | Work truth; underlying source-resource reference; assign/clear/hold/resume/escalate controls; outcome state | N12 · UF13 | **KEEP** Work detail; never source-resource replacement |
| R110 | `/aprovacoes` | authorization-decision queue/history; target/revision qualification | N13 · UF14 | **KEEP** governance queue |
| R111 | `/aprovacoes/:authorizationDecisionId` | governed target/revision evidence; decision truth; create decision when the exact admitted target requires one; return to target | N13 · UF14 | **KEEP** decision context; target execution remains outside Governance |
| R120 | `/configuracoes/canais` | Installation collection; create Installation; current supported marketplace kind; Technical Ingress continuation | N02 · UF02 | **KEEP** low-frequency settings home |
| R121 | `/configuracoes/canais/:marketplaceInstallationId` | Installation truth; update/deactivate; exact OAuth/binding status where Technical Non-Product applies; Selling Entity context link | N02 · UF02 | **KEEP** account detail; provider ceremony stays non-Product |
| R122 | `/configuracoes/entidades-vendedoras` | Selling Entity collection | N02/N08 · UF02/UF10 | **KEEP** read-only settings collection under current Product surface |
| R123 | `/configuracoes/acesso` | member collection; AccessRole catalog; selected-member assignments; assign/revoke controls | N01/N14 · UF01/UF14 | **KEEP** access administration page |
| R124 | `/configuracoes/disponibilidade` | Inventory Source collection; selected source create/edit/deactivate surface; effective allocation-policy surface; policy update | N07 · UF09 | **KEEP** one settings page with two material Availability sub-surfaces |
| R125 | `/configuracoes/expedicao` | Fulfillment Node collection/editor; internal Operational Target surface | N15 · UF15 | **KEEP** one settings page; node config != company-wide WMS |
| R126 | `/configuracoes/politica-comercial` | current Commercial Policy/provenance; revision-aware editor; save/precondition recovery | N16 · UF16 | **KEEP** dedicated policy page; external rules are not editable local truth |
| R127 | `/configuracoes/delegacoes` | delegation collection; selected delegation detail/editor; establish/update/revoke controls | N13 · UF14 | **KEEP** governance configuration page |

### 4.1 Route/page count

```text
content route/page candidates     35
persistent shell                   1
structural homes total            36
```

The old 39-state count remains coverage evidence only. P5 has not removed a Product capability.

## 5. Material composition rules by high-risk page

### 5.1 R21 Publication detail

R21 is one source-qualified human object but has separate owner regions:

```text
Offering current Listing truth
├─ Availability evidence                [Availability]
├─ Performance evidence                 [MarketplacePerformanceIntelligence]
├─ Market evidence when requested       [MarketIntelligence]
├─ Economics evidence when requested    [CommercialEconomics]
└─ related Work                         [OperationalWork]
```

Only Offering controls Offering writes. A Performance region never acquires Listing/Price write authority merely because it is displayed beside them.

### 5.2 R71 Sale detail

R71 is deliberately the largest cross-owner composition in the baseline:

```text
Marketplace Sale                         [MarketplaceSales]
├─ Selling Entity attribution            [MarketplaceSales]
├─ realized economics                    [CommercialEconomics]
├─ BusinessOrderIntent                   [BusinessSystemMaterialization]
├─ Party resolution                      [BusinessSystemMaterialization]
├─ destination realization               [BusinessSystemMaterialization]
├─ InvoicingIntent                       [BusinessSystemMaterialization]
├─ fulfillment summary                   [Fulfillment]
├─ post-sale summary                     [PostSaleResolution]
└─ related Work                          [OperationalWork]
```

These are separate **material regions inside one Sale route**, not a new `Materialização` navigation owner.

D8 P5 remains visible here: a sanctioned contact reference is not presented as a supported full alternate street/fiscal destination override.

### 5.3 R81 Fulfillment execution

The execution page separates current truth from safe physical action:

```text
execution identity/current state
→ separation checkpoint
→ physical conference checkpoint
→ fiscal/provider prerequisite evidence
→ artifact readiness
→ packing checkpoint
→ dispatch handoff
→ Shipment / exception continuation
```

The frontend never manufactures readiness or advances a checkpoint only because a prior button was clicked. D8 P2 remains an implementation-readiness proof debt for the real invoice/label progression.

## 6. Cross-product state grammar

Every P8/P9 screen contract must use the smallest applicable subset of these material states. They are not generic visual component states; they preserve accepted Product semantics.

### 6.1 Context/security states

```text
authenticated + valid Organization
no accessible Organization
invalid/stale deep-linked Organization
exact Installation required but missing
exact Installation invalid/inaccessible
permission-conditioned navigation hidden
server 403 on direct/stale access
source-qualified identity not found in current context
```

A hidden control never substitutes for server authorization.

### 6.2 Read/knowledge states

```text
loading / revalidating
known populated
known empty / known zero
complete
partial
unknown
unavailable
unsupported
stale / current when authority exposes freshness
```

Unknown, unavailable and unsupported never become empty/zero.

### 6.3 Consequential-write states

```text
idle/current
local unsent draft
validation rejected before effect
precondition required / stale precondition
conflict
submitting same semantic request
accepted / converged
accepted / pending
rejected
ambiguous / potentially accepted
transport failure known before effect
```

An ambiguous potentially accepted effect never receives a blind generic retry. The same Idempotency-Key is reused only for a safe retry of the same semantic intake.

### 6.4 Cross-owner composition states

For an optional owner region on a composed page, P5 requires distinct representation of:

```text
not permitted to read that owner
permitted but no related resource exists
permitted but evidence unknown/unavailable/unsupported
related resource exists and is current
related resource exists but is stale/partial
```

The host page may not infer the omitted owner's truth.

## 7. Drawer/modal and alternate-view disposition

P5 assigns **zero material flows to a mandatory modal/drawer baseline**.

A later P8 block may use a drawer/modal only for bounded presentation/confirmation when all of the following hold:

- closing it cannot lose unrecoverable material work;
- the interaction does not need independent durable identity/navigation;
- it does not hide consequential outcome/recovery state;
- accessibility and responsive realization remain plausible.

Material editors, reconciliation actions, high-risk execution checkpoints and source-qualified details stay route/page or durable inline material surfaces by default.

Likewise, P5 admits no baseline table/grid toggle, saved-view platform, generic alternate-view framework or bulk-selection model. A05 remains OPEN rather than becoming speculative UX.

## 8. P6/P7 trigger evaluation by material block

P4 found no IA-level ambiguity. P5 does find a small number of **block-level structural problems** where reference study is justified before rendering.

| Block | Why ambiguity/risk is real | P6 | P7 |
| --- | --- | --- | --- |
| B00 App Shell + global IA | accepted D6 frame is stable and conventional | **NOT TRIGGERED** | **NOT TRIGGERED** |
| B01 Overview | priority evidence A04 is missing, but external pattern research cannot decide organization-specific priority | **NOT TRIGGERED**; resolve through operator/P12 evidence | **NOT TRIGGERED** |
| B10 Preparation | search-first multi-source selection + marketplace requirements/source values/correspondence can plausibly be split-view, progressive detail or list→detail | **TRIGGERED** | **CONDITIONAL after P6** if >1 credible structure remains |
| B20 Publications collection/detail | conventional collection/detail; authority separation is already explicit | **NOT TRIGGERED** | **NOT TRIGGERED** |
| B23 ListingIntent editor | revision-aware draft + media + contextual evidence + submit/discard + consequential outcomes create a high-impact editor problem | **TRIGGERED** | **CONDITIONAL after P6** |
| B24 PriceIntent | bounded collection/detail/create pattern | **NOT TRIGGERED** | **NOT TRIGGERED** |
| B30 Availability | read collection + separate settings continuation is conventional | **NOT TRIGGERED** | **NOT TRIGGERED** |
| B40 Performance (R40–R42) | analytical hierarchy must preserve coverage/comparability/basis without turning unknown evidence into KPI certainty | **TRIGGERED** | **CONDITIONAL after P6** |
| B50 Market + Economics | accepted separation and conventional evidence/detail patterns | **NOT TRIGGERED** | **NOT TRIGGERED** |
| B70 Sale detail | dense cross-owner lifecycle evidence and next-safe-handoff must remain comprehensible without new owner/navigation | **TRIGGERED** | **CONDITIONAL after P6** |
| B80 Fulfillment execution | sequential physical checkpoints + prerequisite/artifact evidence + work-floor/device uncertainty A02 are safety-critical | **TRIGGERED** | **CONDITIONAL after P6/A02 evidence** |
| B90 Post-Sale | bounded list/detail/create pattern | **NOT TRIGGERED** | **NOT TRIGGERED** |
| B100 Work | conventional coordination queue/detail | **NOT TRIGGERED** | **NOT TRIGGERED** |
| B110 Approvals | conventional queue/detail with explicit target/revision | **NOT TRIGGERED** | **NOT TRIGGERED** |
| B120 Settings | conventional low-frequency configuration grouping; semantic owners remain separated inside | **NOT TRIGGERED** | **NOT TRIGGERED** |

P7 alternatives are not manufactured in advance. Each conditional row is opened only if its P6 evidence leaves two or more materially credible structures.

## 9. P8 block ledger / sequencing basis

P5 creates the following **candidate** rendering sequence. This is sequencing, not a `LOCKED` decision.

| Sequence | P8 block | Homes covered | Upstream condition before render |
| ---: | --- | --- | --- |
| 1 | **B00 App Shell + global IA** | G00 | P5 complete; no P6/P7 trigger |
| 2 | B01 Overview | R01 | B00 operator-LOCKED; A04 may remain explicit candidate assumption until operator walkthrough |
| 3 | B10 Preparation | R10 | P6 B10 first; P7 only if still ambiguous |
| 4 | B20 Publications core | R20–R21 | B10 relationship understood; no block-level P6 required |
| 5 | B23 ListingIntent editor | R22–R23 | P6 B23 first; P7 conditional |
| 6 | B24 PriceIntent | R24 | Publications context established |
| 7 | B30 Availability | R30 + R124 relationship | no P6 required |
| 8 | B40 Performance | R40–R42 + R21 Performance region | P6 B40 first; P7 conditional |
| 9 | B50 Market + Economics | R50, R60–R62, R126 relationship | no P6 required |
| 10 | B70 Sales + materialization composition | R70–R71 | P6 B70 first; P7 conditional |
| 11 | B80 Fulfillment | R80–R83 + R125 relationship | P6 B80 + A02 evidence path; P7 conditional |
| 12 | B90 Post-Sale | R90–R91 | Sale/Work relationships established |
| 13 | B100 Work | R100–R101 | underlying source-owner linking grammar established |
| 14 | B110 Approvals + access/governance | R110–R111, R123, R127 | target/revision and access context grammar established |
| 15 | B120 Remaining Settings | R120–R122, R124–R126 | related operational blocks established |

Method law still applies: only the operator may set a block `LOCKED`, and the next dependent material block does not inherit a candidate as locked authority unless the operator explicitly authorizes parallel progression.

## 10. Coverage backcheck against old D6-B1 states

| Prior evidence | P5 home | Disposition |
| --- | --- | --- |
| S00 | G00 | preserved shell |
| S01 | R01 | preserved Overview |
| S10 | R10 | preserved Preparation |
| S20 | R20 | preserved Listings collection |
| S21 + S22 | R21 | consolidated into one Listing detail with owner-separated Offering/Performance regions |
| S23 + S24 | R22 + R23 | preserved collection + dedicated editor |
| S25 | R24 | preserved PriceIntent workspace |
| S30 | R30 | preserved Availability |
| S40–S42 | R40–R42 | preserved Performance subcontexts |
| S50 | R50 | preserved Market |
| S60–S62 | R60–R62 | preserved Economics contexts |
| S70 | R70 | preserved Sales collection |
| S71 + S72 + S73 | R71 | consolidated into one Sale detail with owner-separated materialization/invoicing regions |
| S80–S83 | R80–R83 | preserved Fulfillment/Shipment surfaces |
| S90–S91 | R90–R91 | preserved Post-Sale |
| S100–S101 | R100–R101 | preserved Work |
| S110–S111 | R110–R111 | preserved Approvals |
| S120–S127 | R120–R127 | preserved Settings meanings |

Result:

```text
prior D6-B1 states dispositioned    39 / 39
accepted Product operations covered 99 / 99
orphan prior interaction homes       0
screen-shaped API need               0
new navigation owner                 0
```

## 11. P5 findings

### No architecture reopen

P5 finds no missing human surface that requires a new Product operation, Permission, semantic owner, top-level IA destination or D0–D8 reopen.

### Open evidence carried forward

- A01 task frequency/density still affects later density/default emphasis, not the existence of a route;
- A02 device/work-floor evidence is material to B80 Fulfillment rendering;
- A03 terminology remains candidate until comprehension testing;
- A04 Overview priority remains organization/operator evidence, not external authority;
- A05 bulk UX remains unsupported and absent.

### Block-level research triggers

P6 is required before P8 for **B10 Preparation, B23 ListingIntent editor, B40 Performance, B70 Sale detail and B80 Fulfillment execution**. P7 is conditional for each and activates only if P6 leaves genuine competing structures.

## 12. P5 exit

**DERIVED / CANDIDATE.** Every UF01–UF16 material step has a screen/surface home; all 39 prior D6-B1 states are dispositioned; 99/99 Product operations remain covered without screen-shaped API or new navigation authority.

P8 remains **NOT STARTED**. No D6-R2 UX block is `LOCKED`.

The first renderable block is **B00 App Shell + global IA**, because it has no P6/P7 trigger and downstream blocks inherit its frame. Per Frontend Product Experience Planning Method v2.3, it must be rendered as a structural candidate and explicitly operator-adjudicated before dependent block rendering becomes baseline.
