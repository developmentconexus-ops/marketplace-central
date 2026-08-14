# Marketplace Central — Rebaseline Evidence Register

> **Status:** SUPPORTING EVIDENCE  
> **Baseline:** `main@de1dc88bcef5a6ed5515378e7c646682c0bc15d2`  
> Evidence informs D1–D9; it does not become target architecture by itself.

## Current repository facts

- `internal/modules/` has 21 legacy module directories.
- `internal/contexts/` has 2 new contexts: `catalog` and `listings`.
- `internal/adapters/` already has `erp` and `marketplace` families.
- `cmd/` has 7 entrypoint directories.
- Frontend ownership is split between app-local routes/pages and feature packages.
- Legacy route redirects remain in `AppRouter.tsx` with no production-user compatibility requirement.
- OpenAPI, a hand-written SDK and routing knowledge currently coexist.
- The migration chain contains both legacy and new-context state.
- `scripts/gate.ps1` and `scripts/harness.ps1` are active verification/runtime-test mechanisms.
- Current governance rules still describe parts of the legacy tree; they are changed with the D-stage that owns each code surface, not during documentation cleanup.

## Operator constraints

- No production users require backward compatibility with the current application.
- Hard cutover is allowed when target design requires it.
- Git history is the archive; no active `old/` source tree is desired.
- Technical design must answer context, identity, database, communication, events, external integrations, API, frontend and runtime questions before implementation planning.
- YAGNI removes accidental complexity, not correctness.

## Product/source direction

Historical material consistently describes MPC as an internal marketplace operations/intelligence cockpit:

- Mercado Livre first;
- Sankhya/Oracle supplies internal business facts;
- MPC adds linkage, stock reconciliation, pricing/profitability, order/reconciliation workflows and safe operational controls;
- external protocol details belong behind adapters/ports;
- unknown operational facts must not silently become plausible defaults.

D1–D4 decide the exact boundaries and reverify external/source semantics.

## Domain evidence to carry, not blindly ratify

### Product/listing relationship

Historical operator/domain work indicates the business may require many-to-many product↔listing relationships, variation-level identity and separate kit/BOM semantics. D1/D2 must validate the actual model; a convenient 1:1 current table is not proof.

### Product identity conflict

Old drafts treated `CODPROD` as internal ProductID. Later design work proposed opaque MPC identity plus source keys. These conclusions conflict. D2 must explicitly decide internal identity vs `CODPROD`, EAN, REFFORN and source-instance identifiers.

### Provider identity

Historical code/research distinguishes listing item ID, variation ID, catalog-product relation and seller SKU/custom-field values. D2/D4 define which are provider references versus internal identities and reverify current provider behavior.

### Pricing source mismatch

Historical review found a real class of defect where ERP taxonomy/category could be confused with Mercado Livre listing category when calculating commission/margin. D1/D4 must assign ownership and inputs; D8 must prove the golden flow uses the right facts.

## Sankhya evidence to reverify in D4

A prior read-only pass in the Metal Nobre environment recorded these candidates:

- product code: `TGFPRO.CODPROD`;
- EAN observed at `TGFPRO.REFERENCIA` in that environment;
- supplier reference: `TGFPRO.REFFORN`;
- cost candidate: `TGFCUS.CUSSEMICM` with product/company/as-of semantics;
- tax rule inputs associated with `TGFICM` and realized item tax facts in `TGFITE`.

These are historical measurements, not permanent contracts. D4 rechecks current schema/query semantics before target ports/contracts are ratified.

## Current provider / hub landscape evidence — 2026-08-14

This section records externally researched integration evidence requested during D0.7. It is **supporting evidence only**. Public documentation changes over time; D4 must reverify the exact capabilities, scopes, authentication and provider-mode behavior before accepting integration contracts.

### Cross-provider pattern

The researched channels do not expose one uniform marketplace workflow. Responsibilities can shift by provider **and by operating/fulfillment mode within the same provider**. Depending on the mode, the seller, marketplace, fulfillment provider or intermediary may own/control stock, invoicing, fiscal documents, labels, shipping readiness, tracking, returns or settlement evidence.

A future integration design therefore needs to preserve, rather than hide:

- underlying marketplace/provider provenance;
- effective operating/fulfillment mode;
- which authority owns each native fact/result;
- which prerequisites/data handoffs/artifacts are required before progression;
- which capabilities are supported, unsupported or externally/manual-required;
- correlation from MPC intent/workflow to the actual provider result.

This evidence does not yet mandate a canonical `Operating Mode` entity or a specific capability-interface design; D0/D4 adjudicate the minimum model necessary.

### Amazon — SP-API

Official SP-API documentation exposes distinct fulfillment paths rather than one universal flow: seller fulfillment/Buy Shipping, Easy Ship, FBA/fulfillment-network flows and Multi-Channel Fulfillment. Shipping APIs can create/retrieve shipping labels; Brazil also has shipment-invoicing support for specific FBA Onsite cases. Amazon exposes order management, notifications and order-linked financial events through separate API families.

Architectural evidence for MPC:

- provider brand alone is insufficient to determine responsibility;
- fulfillment mode can alter inventory location, shipping owner and fiscal prerequisites;
- shipping label/document capability is mode-sensitive;
- financial events are a distinct evidence stream from order facts;
- D4 must describe capabilities/authority for the effective Amazon mode, not simply `provider=amazon`.

Research sources: Amazon Selling Partner API official developer documentation (`developer-docs.amazon.com`).

### Magalu

The current Magalu developer platform exposes separate API domains for portfolio/SKU, stock, price, orders, deliveries, invoices, labels, logistics fulfillment, fiscal documents, tracking, webhooks and financial analysis. Current order/FAQ documentation states that XML submission remains a prerequisite to move many orders to invoiced state and that fiscal validation is asynchronous; logistics APIs expose labels and fiscal-document flows, while tracking webhooks include states such as `WAITING_INVOICE` and `READY_TO_SHIP`. Financial-analysis APIs expose order-linked transactions and chargeback-related categories.

Architectural evidence for MPC:

- provider-required invoice/fiscal handoff can be a true dispatch-readiness prerequisite;
- `API 2xx` / submission is not equivalent to accepted fiscal readiness;
- webhook events are notifications and full provider state should be fetched/reconciled;
- provider tracking/readiness has richer state than a generic `shipped/not shipped` flag;
- economic settlement/chargeback evidence may require a separate provider stream from order economics.

Research sources: Magalu official developer documentation (`developers.magalu.com`).

### Grupo Casas Bahia

The official marketplace API covers products, offers/stock/price/status, orders, freight, labels and tracking. `CB Entrega` and `CB Full` have materially different flows. Current Full integration guidance describes separate tokens/contexts for marketplace and fulfillment, disables ordinary seller stock updates in fulfillment context, and exposes fiscal-document retrieval. `CB Entrega` order flow requires fiscal information before generating marketplace labels; labels can be returned as PNG/PDF/ZPL. Official callbacks notify order lifecycle events, but the documentation also recommends periodic order queries so callback loss does not silently omit orders.

Architectural evidence for MPC:

- one seller account/provider may expose more than one operational integration context;
- fulfillment can shift stock and fiscal authority away from seller normal flow;
- label readiness can depend on prior fiscal state;
- callback delivery is not sufficient proof of complete state and reconciliation/polling may still be required;
- integration routing cannot assume one token/account context equals one business responsibility set.

Research sources: Grupo Casas Bahia official developer portal (`developers.grupocasasbahia.com.br`).

### Leroy Merlin / Mirakl

Leroy Merlin's seller portal directs partners to Mirakl for store, order and product management. Mirakl's current seller APIs expose products/offers, order acceptance, tracking, shipment validation, multiple shipments, returns, document upload/download, invoicing/accounting, seller billing cycles and payment transaction logs.

Important limitation: generic Mirakl capability does **not** prove that every Leroy Merlin tenant/configuration enables every Mirakl API or business feature. D4 must verify Leroy-specific configuration and requirements.

Architectural evidence for MPC:

- provider-family reuse is possible (for example shared Mirakl transport/protocol infrastructure), but business capability/authority must remain marketplace-installation specific;
- generic platform features must not be assumed enabled for a particular marketplace;
- order, shipment, returns and settlement/accounting are independently exposed concerns.

Research sources: Leroy Merlin seller portal and Mirakl official developer documentation (`portalseller.leroymerlin.com.br`, `developer.mirakl.com`).

### Shopee

Shopee operates an Open Platform and Brazil logistics include seller-prepared flows and marketplace fulfillment/Full. Publicly crawlable official developer detail is currently limited from this environment, so exact Brazil seller API operations for labels, fiscal handoff and settlement must be reverified with Open Platform/partner access in D4 rather than inferred from third-party SDKs.

Operational evidence from Shopee's own Brazil help/blog confirms responsibility can shift between seller preparation and Shopee fulfillment/CD, which is enough for D0 to reject a static provider-level fulfillment assumption.

Architectural evidence for MPC:

- Shopee must be modeled capability-first and mode-sensitive;
- exact API capability claims remain intentionally unresolved until D4 verification;
- no third-party SDK behavior should become canonical authority for Shopee semantics.

Research sources: Shopee Open Platform / Shopee Brazil official help and blog (`open.shopee.com`, `help.shopee.com.br`, `shopee.com.br`).

### MadeiraMadeira

Public MadeiraMadeira marketplace guidance describes direct API integration for product portfolio, order capture/status update and freight quotation. Detailed API documentation/sandbox access is provided through onboarding rather than being fully public. Official marketplace guidance also documents support for multiple external integrators/hubs.

Architectural evidence for MPC:

- some providers expose a narrower/gated direct integration surface than Amazon/Magalu/Casas Bahia;
- `provider supported` cannot be binary — individual capabilities may be direct, hub-mediated, externally/manual-required or unavailable;
- transport choice (direct API vs hub) must not alter the underlying marketplace provenance/authority.

Research source: MadeiraMadeira Marketplace official integration/help content (`madeiramadeira.zendesk.com`).

### ANYMARKET — hub / centralizer evidence

ANYMARKET's current Backoffice API centralizes catalog, stock, price, listings, orders, fiscal documents, callbacks and related integration concerns. Its fulfillment documentation explicitly changes behavior by marketplace/mode: for some fulfillment flows the marketplace controls stock and/or order status; seller-issued fiscal flows may require NF XML upload, while marketplace-issued flows expose fiscal documents back to the hub. The platform also contains marketplace-specific flags and flows despite its common API.

Architectural evidence for MPC:

- mature hubs normalize common operations but **do not eliminate provider-specific exceptions**;
- a hub can be transport/aggregation intermediary without becoming authority for the underlying marketplace-native fact;
- MPC must preserve both origin provider and intermediary provenance if a hub is used;
- a generic hub API should not tempt MPC to erase mode-specific authority or readiness requirements.

Research source: ANYMARKET official developer documentation (`developers.anymarket.com.br`).

### Magis5 — hub / ERP-bridge evidence

Magis5 HUB exposes per-marketplace integrations and explicit ERP/faturador/logistics configuration. Current guidance shows materially different setups for Mercado Livre, Amazon, Magalu, Leroy Merlin, MadeiraMadeira, Shopee and Casas Bahia. Casas Bahia fulfillment uses a separate integration/token path; some seller-logistics flows require a separate logistics integration because the marketplace connector does not supply a ready label. Magis5 can route marketplace orders into external ERPs or its own fiscal/ERP products.

Architectural evidence for MPC:

- hub configuration commonly binds marketplace account, ERP routing, invoicing provider and logistics provider, but these should remain separate MPC semantics rather than one opaque connector configuration;
- provider/mode exceptions survive aggregation;
- direct and hub-mediated routes should be interchangeable at the integration layer only if marketplace-native provenance and capability semantics are preserved.

Research source: Magis5 official help center (`ajuda.magis5.com.br`).

### Bling — ERP + marketplace/logistics integration evidence

Bling's current API and help content expose sales orders, NF-e, stock and virtual stock, deposits, business units, logistics and webhooks, alongside many marketplace-specific integrations. Current webhooks are explicitly idempotency-sensitive and not guaranteed to arrive in order. Marketplace logistics can impose provider-specific prerequisites; for example current Magalu Entregas guidance requires an authorized invoice sent to Magalu before label printing. Recent API changes also expose business-unit context on sales orders/channels and distinguish physical from virtual stock.

Architectural evidence for MPC:

- ERP-native business-unit/deposit constructs are useful integration evidence but must not dictate MPC canonical Selling Entity / Inventory Source semantics;
- marketplace-specific fulfillment prerequisites remain even when an ERP/hub offers a unified UI;
- event delivery can be duplicate/out-of-order, reinforcing later D3/D7 need for idempotent event handling plus authoritative re-read/reconciliation.

Research sources: Bling official developer and help documentation (`developer.bling.com.br`, `ajuda.bling.com.br`).

### Cross-provider conclusion carried forward

The strongest common pattern is not a universal list of marketplace fields. It is a **capability- and authority-sensitive operating contract**:

```text
Marketplace Installation / underlying provider
  + effective offer/order fulfillment mode
  + integration route (direct or intermediary)
  → effective capabilities
  → effective native authorities
  → required provider prerequisites/data/artifacts
  → MPC workflow/readiness/reconciliation
```

D0 must decide the product responsibility for provider-required readiness without hardcoding provider artifacts. D4 later decides the technical `Capability Profile` / adapter contracts and verifies each provider/mode. D3/D7 later decide notification, polling, idempotency, retry and convergence mechanics.

## Durable uncertainty/safety properties

Evidence repeatedly supports:

- unknown is not zero/default;
- absence from a partial source pull is not automatically terminal state;
- provider wire payload is not domain truth merely because fields are useful;
- external identifier is not automatically canonical identity;
- blind retry of a potentially accepted write is unsafe;
- provider success response is not automatically convergence;
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
- **D2:** identity model, table ownership, exact value/knowledge semantics, recoverability/reset.
- **D3:** sync/event/projection map and event/outbox semantics.
- **D4:** current Mercado Livre, Sankhya, provider-mode and intermediary capability contracts; exact direct-vs-hub routes; provider-required fulfillment/readiness artifacts; settlement sources.
- **D5:** one API contract authority and generation/runtime validation.
- **D6:** frontend feature/package and data-consumption topology.
- **D7:** process, scheduler, transaction, outbox and deployment topology.
- **D8:** end-to-end golden-flow proof by explicitly supported provider/mode.

These are design gates, not tasks for an implementation worker to decide locally.