# Marketplace Central — Rebaseline Evidence Register

> **Status:** SUPPORTING EVIDENCE  
> **Baseline:** `main@de1dc88bcef5a6ed5515378e7c646682c0bc15d2`  
> **Snapshot rule:** repository-tree facts below describe that historical baseline, not the current active tree.  
> Evidence informs D1–D9; it does not become target architecture by itself.

## Historical repository facts at baseline `de1dc88`

- `internal/modules/` has 21 legacy module directories.
- `internal/contexts/` has 2 newer contexts: `catalog` and `listings`.
- `internal/adapters/` already has `erp` and `marketplace` families.
- `cmd/` has 7 entrypoint directories.
- Frontend ownership is split between app-local routes/pages and feature packages.
- OpenAPI, a hand-written SDK and routing knowledge coexist at that baseline.
- The migration chain contains both legacy and newer-context state at that baseline.
- `scripts/gate.ps1` and `scripts/harness.ps1` are baseline-era verification/runtime-test mechanisms.

## Operator constraints

- No production users require backward compatibility with the current application.
- Hard cutover is allowed when target design requires it.
- Git history is the archive; no active `old/` source tree is desired.
- Technical design must answer context, identity, database, communication, events, external integrations, API, frontend and runtime questions before implementation planning.
- YAGNI removes accidental complexity, not correctness.

## Product/source direction

Historical material consistently describes MPC as an internal marketplace operations/intelligence control plane:

- Mercado Livre first;
- Sankhya/Oracle supplies internal business facts;
- MPC adds linkage, availability, pricing/profitability, order/fiscal/fulfillment/reconciliation workflows and safe operational controls;
- external protocol details belong behind provider/business-system boundaries;
- unknown operational facts must not silently become plausible defaults.

D1–D4 decide exact technical boundaries and reverify source semantics.

## Domain evidence to carry, not blindly ratify

### Product/listing relationship

Historical work indicates the business may require many-to-many product↔listing relationships, variation-level identity and separate kit/BOM semantics. D1/D2 must validate the actual target model; a convenient current 1:1 table is not proof.

### Product identity conflict

Old drafts treated `CODPROD` as internal ProductID. Later design work proposed opaque MPC identity plus source keys. D2 must explicitly decide internal identity vs `CODPROD`, EAN, REFFORN and source-instance identifiers.

### Provider identity

Historical code/research distinguishes listing item ID, variation ID, catalog-product relation and seller SKU/custom-field values. D2/D4 decide which are provider references versus canonical identities and reverify current provider behavior.

### Pricing source mismatch

Historical review found a real defect class where ERP taxonomy/category could be confused with marketplace listing category when calculating commission/margin. D1/D4 must assign ownership/inputs; D8 must prove the golden flow uses the right evidence.

## Sankhya evidence to reverify in D4

A prior read-only pass in the Metal Nobre environment recorded candidates including:

- product code: `TGFPRO.CODPROD`;
- EAN observed at `TGFPRO.REFERENCIA` in that environment;
- supplier reference: `TGFPRO.REFFORN`;
- cost candidate: `TGFCUS.CUSSEMICM` with product/company/as-of semantics;
- tax rule inputs associated with `TGFICM` and realized item tax facts in `TGFITE`.

These are historical measurements, not permanent contracts. D4 rechecks current schema/query/API semantics before ratifying target ports/contracts.

## Current provider / competitor landscape evidence — 2026-08-14

This section records external research requested during D0.7. Public documentation changes over time; D4 must reverify exact capabilities, scopes, authentication and operating-mode behavior.

### Cross-provider pattern

The researched channels do not expose one uniform marketplace workflow. Responsibilities shift by provider **and by operating/fulfillment mode within the same provider**. Depending on context, seller, marketplace or fulfillment provider may own/control stock, invoicing/fiscal handoff, labels, shipping readiness, tracking, returns or settlement evidence.

Evidence to carry:

- marketplace/provider provenance must remain explicit;
- effective operating/fulfillment context can change capability and authority;
- prerequisites/data handoffs/artifacts can be required before progression;
- capabilities may be supported, unsupported or externally/manual-required;
- MPC intent/workflow must correlate to the actual provider result.

This evidence does not mandate a canonical `OperatingMode` entity or a specific capability-interface design.

### Amazon — SP-API

Official SP-API documentation exposes distinct fulfillment paths rather than one universal flow: seller fulfillment/Buy Shipping, Easy Ship, FBA and Multi-Channel Fulfillment. Shipping/document responsibilities vary by mode; Brazil has shipment-invoicing support for specific FBA Onsite cases. Order management and financial events are separate API concerns.

Evidence: provider brand alone is insufficient to determine responsibility; fulfillment mode can alter inventory/shipping/fiscal behavior; financial events are distinct from order facts.

Source family: Amazon Selling Partner API official documentation (`developer-docs.amazon.com`).

### Magalu

The current developer platform exposes separate domains for portfolio/SKU, stock, price, orders, deliveries, invoices, labels, fulfillment, fiscal documents, tracking, webhooks and financial analysis. Current guidance indicates XML/fiscal submission can gate order progression and fiscal validation can be asynchronous; tracking includes richer readiness states than generic shipped/not-shipped.

Evidence: fiscal handoff can be a true dispatch-readiness prerequisite; submission/2xx is not accepted-state proof; event notification should not replace authoritative state reconciliation.

Source: Magalu official developer documentation (`developers.magalu.com`).

### Grupo Casas Bahia

The marketplace API covers products/offers, orders, freight, labels and tracking. `CB Entrega` and `CB Full` have materially different flows. Full can shift stock/fiscal responsibility and use separate contexts; CB Entrega can require fiscal information before label generation. Callbacks do not remove the need for authoritative order re-read/reconciliation.

Evidence: one provider account can expose multiple operational responsibility sets; label readiness may depend on fiscal state; notification is not full truth.

Source: Grupo Casas Bahia official developer portal (`developers.grupocasasbahia.com.br`).

### Leroy Merlin / Mirakl

Leroy's seller operation uses Mirakl. Mirakl APIs expose offers, orders, tracking, shipment validation, multiple shipments, returns, documents, invoices/accounting and transaction logs. Generic Mirakl capability does **not** prove every Leroy configuration enables every feature.

Evidence: provider-family protocol reuse may be useful later, but business capability/authority remains marketplace-installation specific.

Sources: Leroy seller portal and Mirakl official documentation (`portalseller.leroymerlin.com.br`, `developer.mirakl.com`).

### Shopee

Shopee operates an Open Platform and Brazil logistics include seller-prepared and marketplace-fulfillment flows. Publicly crawlable official technical detail was limited in this research environment, so exact Brazil API operations for label/fiscal/settlement must be reverified with partner access in D4.

Evidence: Shopee must be treated capability-first and mode-sensitive; third-party SDK behavior is not authority.

Sources: Shopee official Open Platform/help/blog domains.

### MadeiraMadeira

Public marketplace guidance describes direct API integration for product portfolio, order capture/status and freight quotation; detailed API/sandbox access is provided through onboarding.

Evidence: some providers expose a narrower/gated direct surface. Third-party hub support does not prove the direct provider capability MPC can rely on. D4 must classify each required capability as directly supported, externally/manual-required or unavailable for the accepted flow.

Source: MadeiraMadeira official marketplace help content.

### ANYMARKET / Magis5 — competitor-hub evidence only

ANYMARKET and Magis5 demonstrate how mature hubs centralize common catalog/order/stock/fiscal/logistics work while still preserving marketplace- and fulfillment-specific exceptions. Their public guidance shows that fulfillment, XML/NF, labels and ERP routing still vary by marketplace/mode.

**Accepted D0 implication:** these systems are benchmark/competitive evidence, not Product 1.0 runtime dependencies. Their value is proving that a lowest-common-denominator universal marketplace model is insufficient — not justifying `Marketplace → hub → MPC` as target architecture.

### Bling — ERP + marketplace/logistics evidence

Bling exposes sales orders, NF-e, stock/deposits/business units, logistics and webhooks alongside marketplace integrations. Current webhook guidance is idempotency-sensitive and ordering is not guaranteed; marketplace logistics can impose provider-specific prerequisites.

Evidence: ERP-native business-unit/deposit constructs must not dictate MPC Selling Entity / Inventory Source semantics; provider-specific prerequisites survive unified ERP/hub interfaces; event delivery needs later idempotency + authoritative re-read/reconciliation design.

Sources: Bling official developer/help documentation.

### Cross-provider conclusion carried forward

The strongest common pattern is a **capability- and authority-sensitive direct provider contract**:

```text
Marketplace Installation / provider
  + effective offer/order fulfillment context
  → effective capabilities / authorities
  → required provider prerequisites/data/artifacts
  → MPC workflow/readiness/reconciliation
```

D4 later decides the technical capability/adapter contracts and verifies each provider/mode. D3/D7 later decide notification, polling, idempotency, retry and convergence mechanics.

## Time-bound obligation evidence — 2026-08-14

### External provider deadline

Current Mercado Livre shipping documentation exposes a maximum dispatch date/time for a shipment. Dispatch after that authoritative deadline can affect delay metrics, reputation and listing exposure. This establishes that provider time can be a material business obligation rather than a passive timestamp.

Research source: Mercado Livre official developer documentation, shipping-management SLA resource (`developers.mercadolivre.com.br`).

### Internal operating safety target

Mature reliability guidance distinguishes external agreements/obligations from internal objectives and recommends a tighter internal objective/safety margin where appropriate, so teams can react before an external commitment is breached.

D0 uses this only as supporting operating-pattern evidence: organization-owned MPC policy may establish an earlier/faster Internal Operational Target, but it must remain distinct from and cannot relax the provider/external obligation.

Research source: Google SRE official Service Level Objectives / Embracing Risk material (`sre.google`).

### Relative time anchor

The operator additionally identified a business-policy pattern such as “act within N time after event receipt.” Architectural implication: relative targets are meaningless unless the triggering event/time and provenance are explicit. Exact canonical timestamps/event semantics belong to D2/D4; timers/schedulers/notifications belong to D7/D6.

This evidence does not reopen buyer Q&A/chat, which remains outside Product 1.0 unless separately re-adjudicated.

## Action-safety / bulk / execution-time evidence — 2026-08-15

### Amazon asynchronous listing/feed outcomes

Current Amazon SP-API guidance distinguishes direct Listings Items operations from `JSON_LISTINGS_FEED` bulk submissions. Bulk submissions are queued/asynchronous and processing reports expose item issues. Amazon also documents that a feed can reach `FATAL` while some, none or all operations may already have completed, and that issues can arise after an initially accepted listings submission during downstream processing.

Evidence carried into D0/D7:

- submission/acceptance is not final convergence;
- bulk work can have partial or later-discovered outcomes;
- already-accepted work and unresolved work must not collapse into one boolean or blind whole-batch retry;
- the concrete queue/report/retry mechanism belongs to D4/D7, not D0.

Research sources: Amazon SP-API official Feeds/Listings workflow documentation (`developer-docs.amazon.com`).

### Mirakl import tracking and price approval

Current Mirakl offer APIs return an import identifier for offer updates so import status/errors can be tracked. Price imports with `Price Approval` enabled create/update **pending prices** while ongoing prices remain active until the approval process advances them.

Evidence carried into D0/D4:

- a proposed/approved commercial change is distinguishable from currently effective provider state;
- import submission is not equivalent to final effective state;
- approval is tied to the proposed change/context rather than a generic eternal permission to mutate any later price;
- exact import/approval mechanics remain provider-specific D4/D7 design.

Research source: Mirakl official seller/operator offer API (`developer.mirakl.com`).

### ANYMARKET safety-limit blocking

Current ANYMARKET monitoring documentation gives an explicit price-update example where the current price is `100`, proposed price is `150`, the configured safety limit is exceeded, and the listing is **not updated**; the operator must review the new price or safety limit.

Evidence carried into D0:

- mature marketplace hubs re-evaluate current safety constraints at action time rather than treating a previous intent as unconditional execution permission;
- blocked external action becomes explicit operational work/monitoring rather than silent success;
- the exact safety-limit/configuration implementation is competitor evidence, not an MPC architecture to copy.

Research source: ANYMARKET official callback/monitoring documentation (`developers.anymarket.com.br`).

### Bling webhook delivery semantics

Current Bling documentation states the same webhook may be sent more than once and events are not guaranteed to arrive in generation order; the integration is expected to be idempotency-aware. Retrying delivery can continue for days and the webhook configuration may eventually be disabled after persistent failures.

Evidence carried into D3/D7:

- duplicate/out-of-order events are reachable real behavior, not hypothetical hardening;
- event receipt does not by itself prove ordering/completeness/current authoritative state;
- idempotency, ordering/re-read/recovery mechanics belong to D3/D7 rather than becoming another D0 business capability.

Research source: Bling official webhook documentation (`developer.bling.com.br`).

## D0 closure benchmark evidence — 2026-08-15

ANYMARKET's current Backoffice API guide publicly groups its integration surface across catalog prerequisites, products/SKUs/images, stock/listings/prices, orders/returns/fiscal documents, callbacks/campaigns/monitoring/users/roles, plus pagination/rate-limit/error mechanics.

Comparison outcome for D0 closure:

- MPC Product 1.0 covers the material marketplace-operating lifecycle families: readiness/catalog linkage, listings, availability, pricing/intelligence, orders/fiscal materialization, fulfillment/shipment, post-sale, reconciliation/monitoring and economic outcome;
- user/role authorization is represented at D0 actor/authority level while exact identities/permissions belong to D2/D5;
- callbacks/pagination/rate-limit/retry are technical mechanics for D3/D4/D7, not missing Product 1.0 capabilities;
- marketplace campaign authoring/discount-campaign management is not required as a Product 1.0 control surface; however observed promotion/discount economic effects remain attributable evidence where they materially affect price/order/realized economics;
- competitor domain breadth does not justify copying every domain as an MPC module/context.

Research source: ANYMARKET official Backoffice API Guide (`developers.anymarket.com.br`).

## Durable uncertainty/safety properties

Evidence repeatedly supports:

- unknown is not zero/default;
- absence from a partial source pull is not automatically terminal state;
- provider payload is not domain truth merely because fields are useful;
- external identifier is not automatically canonical identity;
- blind retry of a potentially accepted write is unsafe;
- provider success response is not automatically convergence;
- unknown deadline/anchor is not equivalent to no obligation;
- raw observation/reprocessing can be useful where mappings evolve, but must be justified per capability rather than universalized.

Exact mechanisms belong to D2–D4/D7.

## Verification lessons already absorbed into the canonical method

- presence is not execution;
- zero executed checks is not proof;
- a measurement needs a stated universe;
- current directory/schema/ADR shape is not sufficient justification for target structure;
- negative fixtures/counterexamples are valuable proof that a control can fail;
- stale binaries/worktrees can invalidate evidence;
- converging authorities is preferable to hand-syncing duplicate authorities with more guards.

Full historical reports remain in Git history.

## Open evidence questions

- **D1:** final contexts and dissolution of legacy modules.
- **D2:** identity model, table ownership, value/knowledge/time-anchor semantics, recoverability/reset.
- **D3:** sync/event/projection map and event/outbox semantics.
- **D4:** current Mercado Livre/Sankhya/provider-mode capability contracts, provider-required fulfillment/readiness artifacts, settlement sources and authoritative deadlines.
- **D5:** one API contract authority and generation/runtime validation.
- **D6:** frontend feature/package, attention/work-queue and data-consumption topology.
- **D7:** process, scheduler/timer, transaction, outbox and deployment topology.
- **D8:** end-to-end golden-flow proof by explicitly supported provider/mode, including obligation/readiness timing where material.

These are design gates, not tasks for an implementation worker to decide locally.