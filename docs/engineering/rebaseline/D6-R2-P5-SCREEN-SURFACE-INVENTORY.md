# D6-R2 P5 — Current Screen / Material-Surface Inventory

> **Status:** ACCEPTED / CURRENT CONSOLIDATED INVENTORY  
> **Parent:** `D6-FRONTEND.md`  
> **Method:** Engineering Method v1.0.0 + Frontend Product Experience Planning Method v2.3  
> **Product basis:** **106 operations / 31 ordinary Permissions / H-A-S**  
> **Mutable stage/block sequencing:** `docs/roadmap.md`

## 1. Inventory law

P5 maps accepted human jobs to material frontend homes. It is **not** one page per endpoint and does not own current program status.

Separate a material surface only when one of these materially changes:

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

Cosmetic layout/component/responsive variation does not create another material surface.

## 2. Global shell

### G00 — App Shell

Persistent laws:

- Organization is the only global workspace;
- grouped IA follows the operator-LOCKED D6 shell;
- exact Marketplace Installation is page-local context only where required;
- no hidden/default Organization/Installation;
- Organization switch invalidates incompatible URL/server/ephemeral state;
- Permission-conditioned navigation is usability only;
- Personal Notifications uses a bounded shell utility, not another sidebar business mass.

Accepted Notification utility sub-surface:

```text
G00-E bell
→ unread known-present / known-empty / unavailable
→ U01 bounded recent preview
→ source continuation separate from awareness mutation
→ R128 full Inbox
```

## 3. Current content homes

| ID | Route / home | Material owner-separated meaning |
| --- | --- | --- |
| R01 | `/visao-geral` | Installation posture + bounded Economics/Work/Performance composition |
| R10 | `/preparacao` | source Product search, selected subject, requirements/source evidence, correspondence, ListingIntent handoff |
| R20 | `/publicacoes` | Marketplace Listing collection |
| R21 | `/publicacoes/:nativeListingKey` + exact Installation state | Offering Listing truth + owner-separated Availability/Performance/Market/Economics/Work regions |
| R22 | `/publicacoes/intencoes` | ListingIntent collection/work context |
| R23 | `/publicacoes/intencoes/:listingIntentId` | revision-aware ListingIntent authoring/editor/media/evidence/submit-discard outcome |
| R24 | `/publicacoes/precos` | PriceIntent collection/detail/create + Market/Economics handoff |
| R30 | `/disponibilidade` | sellable availability current population + provenance/knowledge/config continuation |
| R40 | `/performance/resumo` | Installation + period Performance summary evidence |
| R41 | `/performance/publicacoes` | Listing Performance collection/detail evidence |
| R42 | `/performance/midia` | Retail Media evidence only; no Ads control |
| R50 | `/mercado` | competitive position/comparable evidence + Economics/PriceIntent continuation |
| R60 | `/economia/prevista` | Expected Economics + scenario evaluation + PriceIntent handoff |
| R61 | `/economia/realizada` | Sale Economics + economic performance summary |
| R62 | `/economia/reconciliacao` | EconomicAttribution collection/detail/resolve |
| R70 | `/vendas` | source-qualified Marketplace Sale collection |
| R71 | `/vendas/:nativeSaleKey` + exact Installation state | Sale + SellingEntity attribution + Economics + Materialization + Fulfillment + Post-Sale + Work composition |
| R80 | `/expedicao/execucoes` | FulfillmentExecution operational queue |
| R81 | `/expedicao/execucoes/:fulfillmentExecutionId` | physical checkpoint/provider/fiscal/artifact execution truth + actions |
| R82 | `/expedicao/envios` | Shipment collection |
| R83 | `/expedicao/envios/:nativeShipmentKey` + exact Installation state | Shipment evidence + exception/Sale/Work continuation |
| R90 | `/pos-venda` | PostSaleResolution collection + bounded create action |
| R91 | `/pos-venda/:postSaleResolutionId` | resolution truth + related owner evidence |
| R100 | `/trabalho` | Work queue |
| R101 | `/trabalho/:workId` | Work truth + assign/clear/hold/resume/escalate; never source truth |
| R110 | `/aprovacoes` | two local lenses: exact-human actionable Requests (`governance.decide`) and immutable history (`governance.read`) |
| R111 request | `/aprovacoes/solicitacoes/:authorizationRequestId` | actionable typed review basis + request-local decision confirmation |
| R111 history | `/aprovacoes/decisoes/:authorizationDecisionId` | immutable AuthorizationDecision history; no decision controls |
| R120 | `/configuracoes/canais` | Marketplace Installation collection/create + Technical Ingress continuation |
| R121 | `/configuracoes/canais/:marketplaceInstallationId` | Installation update/deactivate/binding status |
| R122 | `/configuracoes/entidades-vendedoras` | Selling Entity collection |
| R123 | `/configuracoes/acesso` | members/AccessRoles/assign-revoke |
| R124 | `/configuracoes/disponibilidade` | InventorySource + allocation policy |
| R125 | `/configuracoes/expedicao` | FulfillmentNode + internal Operational Target |
| R126 | `/configuracoes/politica-comercial` | CommercialPolicy current truth/editor |
| R127 | `/configuracoes/delegacoes` | Governance delegation lifecycle |
| R128 | `/notificacoes` | exact-self Personal Notification Inbox/triage |
| R129 | `/configuracoes/notificacoes` | ten fixed ORG_ROUTED Notification route slots + bounded recipient selection |

Organization prefix is `/org/:organizationId` in production route identity. Marketplace Installation and other required qualifiers remain URL/search/path state exactly where their Product subject requires them.

The number of homes is descriptive, never an implementation target.

## 4. High-risk composition laws

### R21 Listing detail

One human Listing route; separate owner regions. Offering alone owns Offering writes. Availability, Performance, Market, Economics and Work do not gain Listing mutation authority by being displayed together.

### R71 Sale detail

One Sale route composes owner-separated Sales, Economics, Materialization, Fulfillment, Post-Sale and Work meaning. No `Materialização` navigation/business owner is invented.

### R81 Fulfillment execution

Physical checkpoint actions are driven by current owner truth, not optimistic button progression. Fiscal/provider readiness/artifacts remain explicit and unknown/unavailable does not become ready.

### R110/R111 Approvals

`governance.decide` and `governance.read` are independent. Actionable Request detail contains typed immutable review basis; successful AuthorizationDecision does not execute the target action. F13 awareness points to the exact Request but grants no capability.

Current `CreateAuthorizationDecision` carrier is:

```text
Idempotency-Key header
body.etag = current AuthorizationRequest StrongETag
body.outcome = authorize | reject
missing/invalid etag → 422
stale request revision → 409
```

No `If-Match`/412/428 remains on this custom `:decide` operation.

### R128 Personal Inbox

Exact-self H-only awareness. URL state:

```text
archive=active|archived
read=all|unread|read
kind=all|<NotificationKind>
```

`all` is frontend state and maps to omission of the Product filter. Source continuation and Notification read/archive mutation remain separate.

### R129 Notification routing

Ten fixed ORG_ROUTED kinds. Inline editor per row. `ListNotificationRouteRecipientCandidates` exposes only `principal_id + display_name`; candidate presence is not authorization. `SetNotificationRoute` uses current route `If-Match`. Configured requires one-or-more recipients; explicit `UNCONFIGURED` is supported.

## 5. Cross-product state grammar

### Context/security

```text
valid Organization
no accessible/stale Organization
required Installation missing
Installation invalid/inaccessible
Permission-conditioned UI hidden
server 403/404 on stale/direct access
source-qualified identity missing
```

### Read/knowledge

```text
loading/revalidating
known populated
known empty / zero
complete / partial
unknown
unavailable
unsupported
stale/current where owner exposes freshness
```

### Consequential write

```text
idle/current
local unsent draft
validation rejected before effect
stale/conflict/precondition failure
submitting same semantic attempt
accepted/converged
accepted/pending
rejected
ambiguous/potentially accepted
known transport failure before effect
```

No blind generic retry after a potentially accepted effect.

### Optional cross-owner region

```text
not permitted to read owner
permitted but no related resource
permitted but unknown/unavailable/unsupported
related current
related stale/partial
```

Host page never infers omitted owner truth.

## 6. P6/P7 current trigger guidance

Reference-study/competing-hypothesis work remains proportional. Previously identified high-risk structural problems include B23 ListingIntent, B40 Performance, B70 Sale detail and B80 Fulfillment when their block opens. B10's prior P6 completed; current B10 issue is an upstream wire falsifier, not a request to redo its full P6/P8 structure.

Do not create P6/P7 ceremony for conventional collection/detail/settings blocks without material ambiguity.

## 7. Current acceptance evidence

P5 does not own LOCK status. The stable current P8 registry and block-specific ratifications/HTML own accepted structure; `docs/roadmap.md` owns mutable stage sequencing.

Current accepted/LOCKED evidence includes:

```text
B00 App Shell / IA
B01 Overview
B00-R2 Notification utility
B11 Personal Inbox
B12 Notification Routing Settings
B110 Approvals
B10 Preparação main structure
```

The B10 correspondence region remains a bounded future revalidation after the paused human-operable read-projection prerequisite. No other locked block is reopened by that finding unless new evidence proves impact.

## 8. Product coverage / upstream findings

Current Product input is 106/31/H-A-S. P5 does not require one visual home per operation; non-human/technical/history/helper operations may be bound inside existing human surfaces or have no separate screen by design.

P5 currently knows one material upstream frontend→Product falsifier: human-operable read projection for Readiness/Offering/Listing-related reads, captured in paused PR #70. Repository-health work preserves that finding but does not implement it.

New frontend evidence may reopen the smallest upstream Product owner when a material human job cannot be realized honestly. It must not be patched through client-side dictionaries, provider calls, N+1 detail fan-out or screen-shaped API convenience.
