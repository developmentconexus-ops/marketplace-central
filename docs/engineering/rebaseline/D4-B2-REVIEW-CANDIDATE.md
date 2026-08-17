# D4 Batch B2 — Mercado Livre Operational Contract — Independent Review Candidate

> **STATUS: REVIEW CANDIDATE — NOT ARCHITECTURE AUTHORITY**  
> **Stage:** D4 — External Integrations  
> **Router state remains:** D4 `OPEN / ACTIVE`; D4-B1 `ACCEPTED / CANONICAL`; D4-B2 `NEXT / NOT YET OPENED`; implementation blocked until D9  
> **Parent authorities:** accepted D0–D3 + canonical D4-B1 + stable `ARCHITECTURE.md` constraints  
> **Operator posture:** B2 direction received a pre-review coherence pass against D0–D3, D4-B1, current Mercado Livre official evidence, reference platforms and DevelopmentConexus Engineering Method v1.0.0  
> **Disposable:** delete after adjudication; durable meaning belongs only in an operator-approved canonical D4 consolidation  
> **Date:** 2026-08-17

## Reviewer bootstrap

Reconstruct target authority independently before reading this candidate:

1. `AGENTS.md`
2. `docs/engineering/rebaseline/README.md`
3. `docs/engineering/standards/root-cause-global-maximum-method.md`
4. `ARCHITECTURE.md`
5. `docs/architecture/decisions/README.md`
6. `docs/engineering/rebaseline/D0-PRODUCT-SYSTEM-DEFINITION.md`
7. `docs/engineering/rebaseline/D1-DOMAINS-BOUNDARIES.md`
8. `docs/engineering/rebaseline/D2-IDENTITY-TENANT-DATA-OWNERSHIP.md`
9. `docs/engineering/rebaseline/D3-COMMUNICATION-EVENTS.md`
10. `docs/engineering/rebaseline/D4-EXTERNAL-INTEGRATIONS.md`
11. `docs/engineering/rebaseline/EVIDENCE-REGISTER.md`
12. only claim-specific current provider/reference documentation and repository/runtime evidence needed to challenge B2.

Do **not** treat this file, chat history, reviewer statements, `AI-DIALOG.md`, legacy ADR status labels, current code shape or old live probes as target authority.

Reviewer findings are evidence. Canonical D4 changes only after GPT adjudication and explicit operator approval.

---

## 1. Batch question and boundary

D0–D4-B1 already decide product scope, business authorities, identity/isolation, Q/C/E/P communication/failure semantics and the generic external-contract grounding.

B2 answers one concrete question:

> **For the Product 1.0 capabilities that touch Mercado Livre, which current provider resources, scopes, preconditions, observation surfaces, effect surfaces and reconciliation surfaces must the Mercado Livre adapter understand so MPC can operate Listings, Availability, Sales, Fulfillment and essential Post-Sale without making Mercado Livre protocol or ontology into MPC business authority?**

B2 is organized by four material failure/authority classes, not endpoint-by-endpoint microdecisions:

- **B2-A — Installation & Channel Surface**
- **B2-B — Offering & Availability Effects**
- **B2-C — Sale & Fulfillment Operational Contract**
- **B2-D — Essential Post-Sale Provider Contract**

B2 deliberately excludes:

- competitor/price-to-win/fee/economic/settlement acquisition — **D4-B4**;
- Sankhya facts/commands — **D4-B3**;
- MPC HTTP/OpenAPI/error encoding — **D5**;
- frontend/UX — **D6**;
- workers, queueing, scheduler cadence, webhook receiver topology, retry/backoff, cursor persistence, secret storage, transaction/process/deployment topology — **D7**;
- end-to-end live golden-flow proof — **D8**;
- implementation — blocked until D9.

---

## 2. Evidence classification

### 2.1 Already authoritative from D0–D4-B1

B2 imports rather than re-decides:

1. Mercado Livre is external and remains authoritative for native account/listing/order/shipment/claim/return/capability state.
2. Marketplace Installation is the MPC namespace qualifier for one Organization marketplace participation; provider seller identity remains external.
3. External Listing/Variation, Sale/Order and Shipment identities are source-qualified; provider-native resources do not gain MPC-global IDs merely for normalization.
4. Provider Product/Catalog ontology does not become MPC Product master.
5. Product & Channel Readiness owns Product↔channel correspondence/readiness; D4 supplies provider identifier evidence only.
6. Marketplace Offering Operations owns listing representation/lifecycle plus Listing/Price Intent and convergence.
7. Availability Control owns Inventory Source/Scope, allocation policy, Sellable Availability and Availability Intent/convergence; native provider stock/location truth remains external.
8. Marketplace Sales owns MPC sale interpretation/context/correlation and transaction-specific Selling Entity attribution.
9. Fulfillment Lifecycle owns provider-requirement closure for flows MPC claims to control; D4 translates native requirement/artifact/capability evidence only.
10. Post-Sale Resolution owns resolution/correlation/closure; provider Claim/Return/refund protocol remains external.
11. Provider notifications are triggers/pointers, not MPC domain truth/events; current provider meaning comes from authoritative reread where material.
12. Known value / known empty / unknown / unavailable remain distinct.
13. Partial acquisition never proves stronger absence/closure than its completed scope.
14. Integration Support != Provider Effective Capability/Requirement != Effective Business Capability.
15. Accepted/rejected/pending/ambiguous and `accepted != completed != externally applied != converged` remain binding where applicable.
16. Ambiguous potentially accepted writes are not blindly retried.
17. Material multi-target effects preserve intended, authorized and attempted/outcome scopes with member distinctions where material.
18. Identity-bearing outbound fields must agree with accepted Readiness correspondence before dispatch; mismatch fails closed.
19. One matching provider identifier is not enough for unattended correspondence; Readiness owns corroboration policy.
20. Provider PII is minimized; no universal raw-payload archive is required merely to guard DTO drift.
21. D4-B1 does not create a generic Provider/Resource/Operation/Capability business graph.

### 2.2 Current Mercado Livre official evidence — reverified 2026-08-17

The following external facts are current provider evidence, not architecture authority by themselves.

#### User Products / listing model

Current official Mercado Livre documentation establishes that:

- an **Item** represents a publication/sales condition;
- a **User Product (UP)** represents a seller-specific physical product at a specific variation level;
- one User Product may be associated with one or more Item IDs;
- sellers are progressively activated into the User Products model and expose the `user_product_seller` tag through the user resource;
- legacy and User Product listings can coexist during migration;
- in the User Product model, the old `variations[]` creation shape no longer applies and `user_product_id` is provider-assigned;
- shared User Product attributes changed through `/items` can propagate asynchronously to other Items associated with the same UP;
- `family_id`, `family_name`, User Product and Catalog Product are distinct provider concepts.

Official sources:

- https://developers.mercadolivre.com.br/pt_br/publicacao-de-produtos/user-products
- https://developers.mercadolivre.com.br/pt_br/guia-para-produtos/preco-variacao

#### Prices / price automation

Current official Mercado Livre documentation establishes that:

- current product prices should be observed through the current price resources rather than relying on `price`, `base_price` and `original_price` from `/items` as durable read authority;
- item creation/edit still uses the applicable item publication surface today;
- since 2026-03-18, Price Automation can block direct API price editing for automated items;
- provider-side pricing automation is therefore a provider-effective restriction that can make direct MPC Price Intent execution unavailable even when the adapter technically supports the write operation.

Official sources:

- https://developers.mercadolivre.com.br/devcenter/api-de-precos
- https://developers.mercadolivre.com.br/pt_br/guia-para-produtos/automatizacoes-de-precos

#### User Product stock / distributed and multi-origin stock

Current official Mercado Livre documentation establishes that:

- seller stock can be represented differently depending on seller configuration and User Product migration;
- sellers with `warehouse_management` / `multiwarehouse` tags can have seller warehouses under the multi-origin model;
- User Product stock can expose seller-managed locations and provider-managed `meli_facility` locations;
- a User Product may coexist with seller-managed and Mercado Livre-managed stock typologies;
- when multi-origin applies, stock is managed by User Product/location rather than by treating `/items.available_quantity` as the universal write target;
- User Product stock reads expose `x-version`; writes require the current version and a stale version returns HTTP 409, after which current stock/version must be reread;
- provider-managed Full stock is not silently seller-writable merely because MPC can read it.

Official sources:

- https://developers.mercadolivre.com.br/pt_br/produto-consulta-de-usuarios/estoque-multi-origem
- https://developers.mercadolivre.com.br/pt_br/lojas-oficiais/estoque-distribuido

#### Orders / Sales acquisition

Current official Mercado Livre documentation establishes that:

- `GET /orders/{id}` is the point resource for order detail;
- an order supplies `shipping` as a provider shipment reference rather than making Shipment the same resource as Order;
- `cancel_detail` has distinct group/code/description/requested-by/date meaning;
- order search supports provider filters but documented retention is bounded to up to 12 months;
- provider order search therefore cannot be the permanent historical authority for MPC Sales lineage.

Official source:

- https://developers.mercadolivre.com.br/pt_br/imovel-consulta-de-usuarios/gerenciamento-de-vendas

#### Shipments / fulfillment prerequisites / time

Current official Mercado Livre documentation establishes that:

- Shipment is queried through its own provider resource;
- `/shipments/{id}/sla` exposes a provider-authoritative maximum dispatch date/time where that surface applies;
- provider shipment state can expose fiscal prerequisites such as `invoice_pending` for applicable logistics before progression to printable/dispatch-ready state;
- provider logistics and fulfillment context change which prerequisites, documents, labels, seller actions and deadlines are applicable.

Official sources:

- https://developers.mercadolivre.com.br/pt_br/gerenciamento-de-envios
- https://developers.mercadolivre.com.br/pt_br/produto-autenticacao-autorizacao/importar-nota-fiscal

#### Reputation / restrictions / moderation

Current official Mercado Livre documentation establishes that:

- seller reputation is exposed on the user resource;
- listing moderation can move listings into paused/under-review/forbidden/held/pending-documentation or penalized states with provider `reason`/`remedy` evidence;
- provider moderation/restriction can affect an Installation or Listing's effective ability to operate;
- these surfaces justify operational attention/capability evidence, not an MPC reputation-management or complaint-management product.

Official sources:

- https://developers.mercadolivre.com.br/pt_br/reputacao-de-vendedores
- https://developers.mercadolivre.com.br/pt_br/autenticacao-e-autorizacao/gerenciar-moderacoes

#### Essential returns / claims

Current official Mercado Livre documentation establishes that:

- Claim/Return are separate provider-native resources;
- a Return can have its own `return_id` and reverse Shipment(s);
- provider claim/player state exposes `available_actions`, so the seller's currently allowed return/review actions are provider-effective capability evidence;
- return scope can be material to partial post-sale consequences and must not be collapsed into one global Sale status.

Official sources:

- https://developers.mercadolivre.com.br/pt_br/relatorios-de-faturamento/gerenciar-devolucoes
- https://developers.mercadolivre.com.br/pt_br/produto-consulta-de-usuarios/gerenciar-reclamacoes

### 2.3 Reference-platform evidence — failure-class evidence only

Reference platforms independently show that marketplace/fulfillment responsibility changes by operating mode; they are benchmarks, not target module templates.

- **ANYMARKET / Mercado Livre Full:** Full shifts storage/picking/shipping and stock/order-status responsibility to Mercado Livre; the hub does not treat stock/status control as universally seller-owned.  
  https://developers.anymarket.com.br/integration_guides/order/fulfillment
- **Amazon SP-API:** seller fulfillment, Easy Ship and FBA use different inventory/fulfillment ownership surfaces; provider brand alone does not determine responsibility.  
  https://developer-docs.amazon.com/sp-api/lang-en_US/docs/sp-api-seller-use-cases
- **Casas Bahia / CB Full vs seller logistics:** fulfillment changes product/stock/order/fiscal integration responsibilities.  
  https://developers.grupocasasbahia.com.br/marketplace/docs/etapas-para-integracao
- **Mirakl:** offers, orders, returns, multiple shipments and invoicing are separate provider resource/operation families rather than one universal Order object.  
  https://developer.mirakl.com/content/product/mmp/rest/seller/openapi3

The supported inference is narrow:

> **Mode-sensitive provider capability/responsibility is essential complexity. A lowest-common-denominator provider model would hide real authority differences. This does not justify copying any reference platform's module/service topology into MPC.**

### 2.4 Repository/runtime evidence only

Historical repository work on 2026-08-01 (`18b6479e5c04db7222926891f50bfe8853f4ca54`) performed read-only Mercado Livre probes against the then-real account and found, among other things:

- provider Order scan behavior beyond ordinary offset limits was reachable in that measurement even though current public Order documentation does not clearly make that behavior a durable contract;
- real canceled Orders exposed `cancel_detail` while `status_detail` was not the cancellation source;
- real listing payloads already exposed `user_product_id` as a possible field;
- real moderated/under-review listings and moderation evidence existed;
- adapter DTO omissions had silently discarded provider fields, proving schema-drift/translation defects are reachable.

This is useful evidence, not current Installation state and not target authority. A 2026-08-01 probe cannot prove which seller tags, stock model, price automation or logistics modes are active on 2026-08-17.

### 2.5 Current Installation evidence — UNKNOWN / gate

This review session can inspect repository files and current official documentation, but no safe connected read-only Mercado Livre runner/credential surface is exposed here. The repository exposes governed historical/live Oracle tooling but no directly executable read-only Mercado Livre probe that can be invoked from this session without external credential access.

Therefore these facts remain explicitly **Unknown** at candidate time:

- whether the real seller currently has `user_product_seller`;
- whether the real seller currently has `warehouse_management` / `multiwarehouse`;
- whether selected real listings are legacy, User Product, catalog-linked, composite, or mixed;
- which stock typologies/locations (`selling_address`, `seller_warehouse`, `meli_facility`) are actually active for the selected proof products;
- whether selected price-controlled listings currently have Price Automation active;
- which logistics/fulfillment modes occur in the recent real selected proof sales;
- which provider fiscal/label/SLA obligations apply to those real modes;
- current seller/listing restrictions/moderations relevant to the selected proof surface;
- whether a real Claim/Return case exists that can later provide D8 live proof.

**Unknown is not converted to a default.** Independent B2 review may proceed with this gate open. Canonical B2 ratification must not claim a selected Product 1.0 Mercado Livre proof lane until a fresh read-only real-dependency probe establishes the relevant current Installation/resource facts or the operator explicitly adjudicates a narrower product-proof assumption.

---

## 3. Root cause and target invariant

### Root cause

A marketplace integration becomes unsafe when provider topology is flattened by convenience.

For Mercado Livre, the same MPC business intent can reach different provider resources, scopes and legal effects depending on seller migration state, User Product relationships, Price Automation, stock ownership/location, listing moderation, sale/shipment context and fulfillment mode.

Without an explicit B2 contract:

- Item can be mistaken for Product or User Product;
- legacy Variation semantics can be applied to a User Product seller;
- one Listing intent can silently mutate multiple Items through a shared UP effect;
- `/items.available_quantity` can be treated as universal stock authority even when stock is location/UP/provider-managed;
- provider-managed Full stock can be presented as seller-writable;
- a stale stock write can be retried as if it were an ambiguous transport failure rather than a rejected precondition;
- Price Automation can make a technically supported price write provider-effectively unavailable;
- Order can absorb Shipment and fulfillment meaning;
- provider search retention can be mistaken for historical Sales authority;
- moderation/reputation restrictions can be invisible while operators think the Integration is healthy;
- Claims/Returns can be flattened into one Sale status or omitted from the controlled lifecycle;
- current adapter/code shape can become target authority by existence.

### Target invariant

> **For every Mercado Livre operation required by Product 1.0, the adapter maps a consumer-owned MPC query/intent to the currently applicable source-qualified provider resource and operating context, preserves provider authority/coverage/preconditions/effect scope, and returns only the semantic evidence needed by the owning D1 domain. Mercado Livre resource topology never becomes MPC business ontology, and provider context never silently widens business authority or intended scope.**

---

## 4. Alternatives / Global Maximum / YAGNI

### A — Harden the legacy Item/Variation + `/items` CRUD model

Treat Item as the central provider object, add more fields, preserve direct item price/quantity writes, and patch Orders/Shipments around it.

**Rejected — Local Maximum.** Current User Products, distributed/multi-origin stock, provider-managed Full stock and Price Automation already invalidate the assumption that Item is the universal product/stock/write surface.

### B — Mirror the Mercado Livre ontology inside MPC

Create canonical MPC models/entities for UserProduct, Family, CatalogProduct, ProviderWarehouse, Claim, Return, OperatingMode and a generic Provider Resource graph.

**Rejected — accidental complexity / provider overfit.** D0 rejects a speculative integration framework; D1 already owns business meanings; D2 rejects synthetic aliases over external resources merely for normalization. This would create duplicate authority and make Mercado Livre ontology a business architecture.

### C — Consumer-owned semantic ports with provider-local concrete topology

Keep Item/UP/Family/Catalog/stock locations/Order/Shipment/Claim/Return provider-native inside the Mercado Livre adapter. Expose only the source-qualified references, provider-effective capability/requirement evidence, current observations, coverage and reconciliation data each D1 owner needs.

**Recommended — Global Maximum.** It preserves essential provider differences without creating a generic integration domain or implementing every provider mode in advance.

### YAGNI rule

B2 prepares the seams required by real provider evolution but does **not** implement every Mercado Livre mode by design.

- A provider mode/context that is present in the selected real Product 1.0 proof receives a concrete supported contract.
- A known provider mode not needed by the selected proof may remain explicit `unsupported` / `external-required` without weakening source semantics.
- A second real consumer/failure class may later justify shared mechanism; hypothetical providers/modes do not.

---

## 5. Candidate B2-A — Installation & Channel Surface

### B2-A1 — Installation operational posture / reputation

Marketplace Portfolio may consume provider-authoritative Installation posture evidence needed for operational attention.

1. Current seller identity remains the B1-bound external seller namespace.
2. Seller reputation/health evidence is external observation; it may feed Portfolio attention when materially relevant.
3. Seller restrictions/suspension or provider conditions that make the integration/listing unable to operate are Provider Effective Capability/Requirement evidence.
4. D4 does not turn reputation metrics into MPC-owned reputation state/policy.
5. Reputation optimization, complaint management, buyer messaging and automated reputation-management workflows remain outside Product 1.0.
6. Authentication/application failure remains B1 integration availability; seller business posture/reputation is a separate provider observation and must not be collapsed into auth state.

### B2-A2 — Item / User Product / Family / Catalog topology stays provider-local

1. `item_id` is the provider Listing/offer condition identity under Marketplace Installation.
2. A legacy provider Variation remains an external child identity only for listings that actually use the legacy variation model.
3. `user_product_id` is a provider-native physical-product resource reference; it is **not** MPC Product and gains no MPC canonical entity by normalization.
4. `family_id` / `family_name` are provider grouping concepts, not MPC Product-family authority.
5. `catalog_product_id` remains Mercado Livre catalog ontology and does not replace authoritative business-system Product identity.
6. Legacy and User Product listings may coexist; the adapter detects the actual provider model from current authoritative seller/resource evidence rather than a global static assumption.
7. Provider topology may be referenced internally by source-qualified resource kind + native key where material, but B2 creates no universal `MLResource` entity graph.
8. A later provider migration that changes Item↔UP relationships does not rewrite MPC Product identity or historical Readiness decisions.

### B2-A3 — Listing creation is mode-aware and provider-authoritative

D0 requires a real Product 1.0 creation/publication path, so creation is first-class B2 scope.

1. Offering owns Listing Intent and desired listing meaning.
2. Before dispatch, D4 establishes the provider publication model and mandatory provider requirements for the current seller/category/resource context.
3. If the seller is in the User Product model, B2 does not fabricate legacy `variations[]` creation semantics.
4. If a required provider category/domain/attribute/catalog condition is unknown or unsupported, the Listing Intent cannot be presented as provider-ready merely because a generic POST endpoint exists.
5. Provider-created `item_id`, `user_product_id`, family/catalog relations and resulting current state are reread/observed as external results after creation.
6. Provider creation success is not MPC convergence until Offering reconciles intended listing meaning with authoritative current provider state.
7. B2 does not choose the MPC HTTP create-listing route or frontend form.

### B2-A4 — Listing observation / lifecycle / moderation

1. Point listing observation uses the current authoritative provider listing resource(s) appropriate to the seller/item model.
2. Provider status/substatus/tags/moderation remain provider-native evidence translated to the smallest semantic meaning Offering/Portfolio needs.
3. Moderation `reason`/`remedy` may establish a concrete provider condition or required operator attention; D4 does not own the business closure decision.
4. Listing status is never derived solely from absence in an enumeration.
5. A partial/failed enumeration cannot close or delete a listing by absence.
6. A listing becomes known closed/inactive only when the authoritative source contract affirmatively establishes that meaning for the source-qualified resource.

### B2-A5 — Seller population enumeration / coverage

1. Seller listing population acquisition is scoped to the authenticated Installation seller namespace.
2. B1 scan/scroll/completeness rules apply; successful first-page/offset traversal is not automatically complete for large populations.
3. Provider search restrictions/caps are part of coverage evidence.
4. The result must preserve whether the population observation is completed or partial for the exact provider-defined scope.
5. Enumeration cursors/scroll IDs stay provider protocol and D7 persistence mechanics; they do not become business identities.

### B2-A6 — Product↔channel provider identifier evidence

D4 supplies current provider evidence; Product & Channel Readiness owns correspondence sufficiency.

Where present and materially reliable, B2 may expose/reference provider evidence such as:

- seller SKU / seller custom field;
- GTIN/EAN and relevant provider attributes;
- Item ↔ legacy Variation;
- Item ↔ User Product;
- User Product / Family evidence;
- Catalog Product relation;
- provider category/domain/attribute evidence needed to understand correspondence constraints.

Rules:

1. No one provider field becomes an MPC identity law by B2.
2. `SELLER_SKU == CODPROD` is not resurrected.
3. One matching field does not establish unattended correspondence.
4. Material contradictory provider evidence is preserved for Readiness rather than normalized away.
5. Before a provider write carries an identity-bearing SKU/attribute, the accepted Readiness correspondence is revalidated under D2 §10.1.
6. Standing human Readiness decisions are not silently overwritten by a later recurring acquisition run.

---

## 6. Candidate B2-B — Offering & Availability Effects

### B2-B1 — Price observation authority is separate from price write mechanism

1. Current price observation uses the current authoritative Mercado Livre price surface rather than treating deprecated `/items` price fields as durable read authority merely because they remain present.
2. Offering owns Price Intent; D4 never owns price policy or economic validity.
3. The current provider write mechanism may still be item-oriented even when the authoritative read surface is the price resource. Mechanism and observation authority are recorded separately.
4. Provider price context/channel distinctions stay provider-local unless a consuming domain needs a concrete semantic distinction.
5. A successful item write cannot be called Price convergence until the authoritative current price is reread and matches the intended provider-effective result.

### B2-B2 — Price Automation is a provider-effective restriction

1. Adapter implementation of price mutation proves only Integration Support.
2. When current provider evidence says Price Automation blocks direct API price editing, Provider Effective Capability for direct MPC price write is unavailable/rejected for that listing/context.
3. B2 does **not** automatically disable or reconfigure Mercado Livre Price Automation merely to make the MPC write path succeed.
4. Offering decides how an unavailable/external-required price action affects its business intent/workflow; D4 does not fabricate business permission.
5. A provider response that accepts other item fields while ignoring/rejecting price is not price convergence; authoritative price reread controls the conclusion.

### B2-B3 — Provider-shared User Product effect scope cannot silently widen MPC intent

Mercado Livre can propagate shared User Product attribute changes asynchronously to multiple Items associated with the same UP.

Therefore:

1. Before dispatching a provider-shared field change, D4 establishes the current provider resource/effect scope sufficiently to know whether the operation can affect more than the nominal Item.
2. D4 never silently rewrites a single-Listing domain intent into a multi-Listing intent.
3. If provider effect scope is materially broader than the intended/authorized scope, the action fails closed before dispatch or returns an explicit `requires-scope-redecision`-equivalent semantic condition to the owner; exact encoding is D5/D7.
4. Offering/Governance remain owners of intended/authorized scope; D4 only supplies provider-effective scope evidence.
5. Where a legitimate domain intent already spans all affected resources, D3 multi-target attempted/outcome distinctions remain available where the provider can partially/later apply effects.
6. A User Product relation changing after authorization is execution-time drift; current effect scope is revalidated when materially consequential.

### B2-B4 — Provider stock topology is mode-sensitive

`available_quantity` is not one universal marketplace stock contract.

B2 distinguishes at least these provider cases **only as observed protocol contexts**, not new MPC business entities:

- seller-writable simple/legacy item availability where currently supported;
- User Product/location stock where User Product stock applies;
- seller-managed `selling_address` / `seller_warehouse` locations where applicable;
- provider-managed `meli_facility`/Full stock where the seller/MPC cannot legitimately write the quantity.

Rules:

1. Availability Control owns Inventory Source/Scope and Sellable Availability; provider stock location is external evidence/reference.
2. A provider `store_id`, `network_node_id`, `selling_address` or `meli_facility` never becomes Inventory Source by convenience.
3. D4 supplies the provider stock resources/ownership/location evidence needed for Availability to map eligible external stock participation explicitly.
4. A seller-writable stock mechanism is available only where current provider context supports it.
5. Provider-managed Full stock is read/observed but not presented as a seller-writable Availability capability.
6. Unsupported/provider-managed paths are explicit rather than mapped to a successful no-op.
7. A listing may have multiple provider stock/logistic contexts without creating a universal MPC `OperatingMode` entity.

### B2-B5 — Stock version conflict is a rejected stale precondition

For current User Product stock surfaces that require `x-version`:

1. The version is a provider protocol precondition, not a business identity or global sequence.
2. A stale version returning 409 is **rejected due to stale provider precondition**, not an ambiguous possibly accepted write.
3. The correct next step is authoritative stock/version reread and a fresh owner decision/revalidation.
4. D4/D7 must not blindly repeat the same stale request.
5. A new request after reread is a new execution against current provider state, still subject to current Availability intent, policy, authorization and correspondence validity.

### B2-B6 — Automatic availability synchronization is bounded by provider authority

D0 requires routine sufficiently-known policy-valid availability synchronization to be automatic where authorized.

B2 therefore requires at least one real provider-writable availability lane in the selected Product 1.0 proof set.

- If the selected proof resource is seller-writable, B2 supplies the concrete provider write/reconcile contract.
- If a resource is provider-managed Full stock, Availability observes it but does not pretend to control it.
- A Full-only proof set cannot by itself prove D0's seller/MPC-controlled availability-convergence capability.
- If the real Installation offers no provider-writable availability lane for the selected Product 1.0 scope, surface a targeted product/proof conflict; do not mark the capability complete by observing Full stock only.

This is a Product 1.0 proof-selection constraint, not a requirement to implement every stock mode.

---

## 7. Candidate B2-C — Sale & Fulfillment Operational Contract

### B2-C1 — Provider Order notification triggers authoritative Order reread

1. Order notification/callback evidence is a trigger/pointer only.
2. Marketplace Sales commits MPC sale meaning only after the adapter obtains the provider order evidence required by the Sales contract.
3. Duplicate notifications may cause repeated point rereads without creating duplicate Sales meaning.
4. Provider order arrival order is not MPC business ordering authority.
5. Seller/order namespace mismatch fails closed under B1 acquisition-attribution rules.

### B2-C2 — Order point read and population/history coverage are distinct

1. `GET /orders/{id}` is the point authority for the source-qualified Order details B2 needs.
2. Order search is an enumeration surface with provider filters/retention, not permanent Sales historical authority.
3. Current public documentation says orders are retained/searchable for up to 12 months; MPC historical claims required beyond that bound must survive in accepted MPC/external durable lineage per D2/D3.
4. Historical 2026-08-01 live evidence that `/orders/search?search_type=scan` worked beyond ordinary offset is **evidence only** where current official documentation does not clearly guarantee that behavior.
5. If B2 needs >ordinary-limit complete Order enumeration for the selected real Installation, current real-dependency proof of the chosen traversal is required before claiming complete coverage.
6. Enumeration failure/retention boundary remains partial/unavailable for the uncovered range; it does not fabricate “no older orders”.

### B2-C3 — Order and Shipment are separate provider resources

1. Provider Order remains Sale/checkout evidence; provider Shipment remains its own source-qualified external identity.
2. `shipping.id` on Order is a reference/correlation input, not permission to treat embedded shipping data as Shipment authority forever.
3. Fulfillment/shipment current state is obtained from the provider Shipment surface when materially needed.
4. Provider Pack remains provider-native when relevant; B2 does not introduce a universal MPC Pack.
5. Marketplace Sales does not become Fulfillment authority merely because the Order JSON includes shipping-related fields.

### B2-C4 — Sale cancellation/fraud provider evidence is translated, not normalized into one status

1. `cancel_detail` and equivalent provider evidence remain source facts; B2 translates only the semantic evidence Marketplace Sales/Post-Sale require.
2. Cancellation group/code/requester are not replaced by a guessed generic reason from `status_detail`.
3. Fraud-risk/provider stop-shipment evidence is a provider-effective fulfillment/sale condition and may block progression; it is not silently ignored because payment was previously approved.
4. D4 does not own cancellation Resolution semantics; Marketplace Sales/Post-Sale remain the relevant owners.

### B2-C5 — Effective fulfillment context controls provider requirements/capabilities

1. Marketplace brand alone does not determine fulfillment responsibility.
2. D4 determines current provider-effective logistics/fulfillment context from the concrete Order/Shipment/provider evidence needed by the selected path.
3. Provider strings/codes such as logistic types/substatuses remain inside the adapter; consumers receive semantic requirement/capability evidence.
4. Fulfillment Lifecycle owns the conclusion that a claimed path's provider requirements are sufficiently closed.
5. A provider-managed path may be `provider-operated`/`external-required` rather than pretending MPC controls physical steps the provider owns.
6. B2 does not create a canonical global `OperatingMode` entity solely to normalize Mercado Livre.

### B2-C6 — Provider fiscal prerequisites feed Fulfillment closure without moving fiscal authority

Where the selected Mercado Livre path exposes fiscal prerequisites before label/dispatch readiness:

1. D4 translates the current provider requirement/artifact/state.
2. Fulfillment owns whether that provider requirement is still open/closed for its path.
3. Business-System Materialization remains owner of Business Order/Invoicing Intent and authoritative Sankhya fiscal materialization under B3.
4. Provider-required NF/XML/DC-e/other artifact submission is a provider integration effect only where the selected path actually requires/supports it.
5. Submission/transport success is not proof that the provider requirement closed; authoritative Shipment/provider reread establishes current provider state.
6. If a required artifact operation is unsupported through the accepted integration, the path is explicit unsupported/external-required rather than silently depending on routine seller-portal work while MPC claims full control.

### B2-C7 — Labels / dispatch readiness are provider-effective capabilities

1. Label/document availability is mode/state sensitive.
2. Adapter support for a label endpoint does not prove a label is currently available for the Shipment.
3. Provider-managed modes that do not allow seller label generation remain explicit provider-managed capability, not integration failure.
4. Fulfillment may proceed only when its own business readiness plus provider-effective dispatch prerequisites are satisfied.
5. Provider document payload/format is protocol; it does not become domain state merely because the operator later downloads/prints it.

### B2-C8 — External SLA/deadline is reread provider authority; internal target stays domain policy

1. `/shipments/{id}/sla` or the applicable current provider deadline surface is external-authoritative evidence where available for the selected path.
2. Provider deadline and MPC internal operating target remain different authorities per D0/D1.
3. A consequential action close to dispatch uses sufficiently current provider deadline evidence; stale cached deadline cannot be treated as immutable provider contract.
4. Absence/unavailability of a provider deadline is not “no deadline” unless the source semantics prove non-applicability.
5. D7 owns timers/schedulers/attention mechanics; D4 owns only acquisition/translation of the external deadline evidence.

### B2-C9 — Selected first-flow proof is evidence-driven, not “support every mode”

B2 must not build every Mercado Livre fulfillment/logistics mode merely because documentation lists them.

Before canonical B2 ratification, the real Installation probe selects the minimum real lane set needed to prove accepted D0 outcomes.

At minimum:

1. the proof set includes a provider-writable availability lane if D0 Availability Control is claimed complete;
2. the selected Sale/Shipment lane has explicit fulfillment responsibility and closes every provider prerequisite the MPC claims to control;
3. if the operator selects the D0 internally-operated Fulfillment Node normal path as part of the first proof, the proof set includes a seller-operated physical fulfillment lane — a Full-only provider-operated flow cannot stand in for internal separation/conference/invoicing/packing/dispatch responsibility;
4. if the selected first proof deliberately uses a provider-operated fulfillment lane, B2 may support that lane honestly but must not claim it proves internal physical execution the provider owns;
5. unsupported modes remain explicit and may be added later from real evidence without changing D1 authority.

---

## 8. Candidate B2-D — Essential Post-Sale Provider Contract

### B2-D1 — Claim / Return / reverse Shipment remain provider-native references

1. Provider Claim/Case/Return resources remain source-qualified under Marketplace Installation; B2 does not create MPC Claim/Return aliases merely for normalization.
2. Post-Sale Resolution remains the MPC canonical obligation/correlation/closure owner.
3. One Sale may have multiple/partial post-sale consequences; provider Return scope must remain representable at line/item/quantity scope where material.
4. Reverse Shipment remains provider-native shipment evidence and does not collapse into the original Shipment.

### B2-D2 — `available_actions` is provider-effective capability, not business authority

1. Provider claim/player `available_actions` or equivalent current provider state tells D4 what the seller/provider currently allows.
2. That is Level 2 Provider Effective Capability only.
3. Post-Sale and the consequence-owning domains decide whether MPC should request a refund/review/return-related consequence under business validity/policy/authorization.
4. D4 never turns provider “available” into automatic MPC execution permission.
5. Before a consequential post-sale provider action, current provider capability/state is reread when material.

### B2-D3 — Post-sale observation is bounded to Product 1.0, not CRM/SAC

B2 supports only the provider surfaces necessary for essential cancellation/return/refund consequences of MPC-controlled marketplace sales.

Out of scope remain:

- general complaint management;
- buyer Q&A/chat/messaging automation;
- reputation-management workflows;
- generalized customer-service case management;
- company-wide reverse-logistics replacement.

A provider Claim that is merely customer-service conversation is not automatically a new Product 1.0 workflow. A Claim/Return that creates a material cancellation/return/refund consequence is relevant evidence for Post-Sale Resolution.

### B2-D4 — Financial refund/settlement movement belongs B4

1. B2 may correlate provider Claim/Return consequences to the provider-native financial movement references needed downstream.
2. The authoritative payment/refund/fee/adjustment/settlement movement contract and economic interpretation remain **D4-B4 + Commercial Economics**.
3. Post-Sale closure cannot fabricate financial completion because the provider Return reached a terminal physical state.
4. Commercial Economics cannot infer physical/post-sale closure solely from a financial movement.

---

## 9. Cross-cutting operation matrix

This matrix describes owner/evidence boundaries, not HTTP/API/D7 realization.

| MPC need | Business owner | Mercado Livre provider concern | External effect? | Authoritative reconciliation surface |
|---|---|---|---|---|
| Observe Installation posture | Marketplace Portfolio | seller identity/reputation/restriction | no | current user/moderation/restriction evidence |
| Create listing | Marketplace Offering Operations | current Item/User Product publication model + provider requirements | yes | resulting Item/UP current read |
| Observe listing lifecycle | Marketplace Offering Operations | Item + moderation/current state | no | current Item/moderation read |
| Establish provider identifier evidence | Product & Channel Readiness consumes | Item/Variation/UP/Catalog/SKU/GTIN/attributes | no | current source evidence; Readiness owns correspondence |
| Change listing shared fields | Marketplace Offering Operations | Item/UP effective scope | yes | all materially affected provider resources/current state |
| Observe price | Marketplace Offering Operations consumes | current Price resource/context | no | current Price read |
| Apply Price Intent | Marketplace Offering Operations | item price effect + Price Automation restriction | yes | current authoritative Price read |
| Observe provider stock | Availability Control consumes | applicable Item or UP/location stock | no | current provider stock read |
| Apply Availability Intent | Availability Control | provider-writable Item or UP/location + version/preconditions | yes | current provider stock/quantity read |
| Acquire sale | Marketplace Sales | Order | no | `GET /orders/{id}` |
| Observe shipment/requirements | Fulfillment Lifecycle consumes | Shipment + provider requirements/deadline/artifacts | no | current Shipment/SLA/requirement read |
| Submit provider-required fulfillment artifact | Fulfillment Lifecycle orchestration with Materialization evidence where applicable | selected Shipment provider operation | yes | current Shipment/requirement reread |
| Execute essential post-sale provider action | Post-Sale coordinates; consequence owner remains D1 owner | Claim/Return/provider action | yes where applicable | current Claim/Return/reverse-shipment/provider result reread |

No row creates a generic MarketplaceOperation owner.

---

## 10. Installation Evidence Gate before canonical B2 ratification

### Why the gate exists

The Method requires Unknown to remain Unknown and asks for the smallest real seam, not every hypothetical mode.

Current public documentation proves multiple possible Mercado Livre operating modes, but it does not prove which modes the Metal Nobre Installation currently uses. Historical August 1 measurements are too old to act as current mode authority.

### Required probe properties

Before B2 becomes canonical, execute a **read-only real-dependency probe** against the real Installation with these constraints:

- no Mercado Livre write;
- no secret/token value recorded in artifact/log;
- no buyer PII retained;
- record only the minimal classification/evidence needed for the architecture claim;
- point/resource rereads use current authoritative provider APIs;
- the probe universe/sample is stated so absence is not overclaimed.

### Minimum facts to establish

1. seller tags relevant to publication/stock model: `user_product_seller`, `warehouse_management`, `multiwarehouse` where currently applicable;
2. current selected listing topology: legacy vs User Product, Item↔UP relationships, Catalog relation, real Variation/composite presence where relevant;
3. current provider price model for candidate listings and whether Price Automation is active on the intended price-write proof listing;
4. current stock typologies/locations/ownership for candidate availability-control products, including whether seller-writable stock exists;
5. recent selected real Orders/Shipments and their actual fulfillment/logistics/fiscal/label/SLA contexts;
6. current Installation/listing moderation/restriction evidence materially relevant to the proof set;
7. Claim/Return presence if available for later D8 proof; absence from the current sample does not mean the provider contract is nonexistent.

### Gate outcomes

- **If a seller-writable availability lane exists:** choose the smallest real lane needed to prove D0 Availability Control and keep other modes unsupported/external-required until needed.
- **If no seller-writable availability lane exists for the accepted Product 1.0 scope:** surface a targeted D0/product-proof conflict before calling Availability Control complete; Full/provider-managed observation alone is insufficient proof of MPC-controlled convergence.
- **If Price Automation blocks every candidate Price Intent path:** direct price write remains provider-effectively unavailable; do not disable provider automation by architecture assumption.
- **If the selected first fulfillment path is provider-operated:** support it honestly; do not claim internal Fulfillment Node execution from provider-owned steps.
- **If the operator requires the internal Fulfillment Node path in the first proof:** select/prove a real seller-operated lane or surface the gap.
- **If a provider feature required by D0 cannot be reached on the real Installation:** mark explicit unsupported/external-required or trigger only the actually implicated upstream decision; never fabricate support.

This evidence gate is a prerequisite to **canonical B2 ratification**, not a reason to block independent review of the architecture candidate.

---

## 11. Legacy/current-state disposition proposed by B2

### ADR-015 — legacy read-only Listings module

Proposed final D4 disposition after B2 acceptance: **historical / target structure superseded**.

Rehome only durable meanings under current authority:

- provider listing state is external observation rather than locally editable provider truth;
- Listing Intent/lifecycle/convergence belongs Marketplace Offering Operations;
- source-qualified listing/variation identity belongs D2 + B2 concrete provider mapping;
- unknown/freshness/partial coverage remains D2/D3/D4-B1/B2;
- provider acquisition enters through consumer-owned ports/adapters.

Do **not** retain as target merely because ADR-015 once specified them:

- one canonical read-only `listings` module/table;
- ingestion as the only writer of a shared read model;
- composite `installation~provider_listing_id~variation` MPC ID;
- manual refresh runtime;
- one full pull as permanent sync model;
- “absent from completed pull = closed”.

The last rule directly conflicts with current D0/B1 honest partial-observation authority.

### ADR-022 / ADR-028 provider-evidence obligation

Their old identity formulas remain superseded. B2 fulfills only the D4 evidence obligation:

- establish the current Mercado Livre identifier/correspondence evidence surface;
- preserve contradictory provider evidence;
- enforce pre-dispatch correspondence consistency;
- leave unattended-corroboration policy and human-decision lifecycle to Product & Channel Readiness.

### Historical MIS-007 provider DTO/raw-payload proposals

Current-state evidence showed real DTO omissions and silent provider-field loss. B2 accepts the defect class but **does not** adopt universal raw provider payload persistence as architecture.

- Provider schema/translation drift needs proportionate proof and detection.
- Targeted sanitized fixtures, selective PII-minimized evidence snapshots, contract tests or later runtime observability may prove the mapping.
- Exact mechanism belongs implementation/D7/D9 proof planning.
- Raw buyer/billing/address payload is never persisted merely because it helps adapter debugging.

---

## 12. Explicit deferrals

### D4-B3 — Sankhya Business-System Contract

Owns authoritative internal Product/native keys, stock/cost/tax/fiscal/business-order/invoicing facts/commands through the sanctioned Sankhya API Gateway. B2 may name the semantic evidence it needs from Materialization/Readiness/Availability but does not design Sankhya protocol.

### D4-B4 — Market / Economics / Settlement Contract

Owns:

- competitor/market offer evidence;
- `price_to_win` / catalog competition evidence where justified;
- fee/selling-cost evidence;
- provider Payment/Refund/Fee/Adjustment/Settlement/Payout movement contracts;
- realized-economic correlation/provenance/completeness.

B2 does not use payment movement fields as shortcuts for Sales/Post-Sale truth.

### D5

Owns MPC HTTP/OpenAPI operation/permission/error/idempotency/precondition representation.

### D6

Owns screens, tables, portfolio attention, work/approval UX and projection consumption.

### D7

Owns:

- webhook receiver/ack topology;
- worker/scheduler/process placement;
- polling/reconciliation cadence;
- cursor/scroll persistence;
- retry/backoff/rate-limit machinery;
- token refresh coordination/secrets;
- transaction/outbox/locks/concurrency;
- runtime enforcement implementation.

B2 may require that a property be enforceable/reconcilable but does not choose these mechanisms.

### D8

Owns end-to-end real proof of the selected B2/B3/B4 operating lanes. B2's Installation Evidence Gate selects the real provider context needed before canonical contract ratification; D8 later proves the complete lifecycle through the accepted contracts.

---

## 13. Proof strategy / strongest counterexamples

B2 should not be ratified unless it survives at least these falsification cases:

1. Seller is not yet `user_product_seller` -> legacy listing creation/observation remains supported without pretending UP is mandatory.
2. Seller is `user_product_seller` -> legacy `variations[]` creation is not sent and provider model requirements are honored.
3. Same Installation contains both legacy and User Product listings -> adapter handles actual resource model per listing without two MPC Product identities.
4. One User Product is linked to several Items -> shared-field effect scope is discovered rather than silently treated as one Item.
5. Domain intent/authorization targets Item A but provider-shared change would affect A+B -> fail closed / require owner scope redecision before dispatch.
6. Item↔UP relationship changes after approval -> execution-time scope drift is revalidated.
7. Current authoritative Price differs from deprecated `/items.price` -> price observation uses current Price surface.
8. Price Automation is active -> adapter technical support does not become effective write capability; no false Price convergence.
9. Provider accepts unrelated item fields while price is ignored/rejected -> 2xx does not prove price effect.
10. Seller uses simple seller-writable availability -> selected contract can write/reconcile the correct provider quantity.
11. Seller uses User Product multi-origin stock -> B2 does not write `/items.available_quantity` as a universal shortcut.
12. Provider stock returns `seller_warehouse` and `meli_facility` -> only legitimate seller-managed scope is writable; Full stock remains provider-managed.
13. Stale `x-version` returns 409 -> classified as rejected stale precondition, reread before any new owner decision; no blind retry.
14. Provider stock location ID resembles internal warehouse code -> does not become Inventory Source/Fulfillment Node by coincidence.
15. Listing disappears from one failed/partial enumeration -> it is not marked closed by absence.
16. Listing point read says active while an older pull omitted it -> authoritative point state wins over the incomplete absence claim.
17. Moderation makes an Item under-review/forbidden -> Listing capability/attention changes explicitly; no fabricated normal operation.
18. Seller reputation degrades -> Portfolio may surface attention without creating reputation-optimization authority.
19. Duplicate Order notification -> point reread is safe and Sales does not duplicate a Sale occurrence.
20. Order search no longer returns a >12-month Sale -> MPC historical lineage required for accepted claims remains explainable without provider search history.
21. Selected Installation has >ordinary Order offset limit and current public docs do not guarantee scan -> complete traversal is not claimed from stale historical measurement without a real-dependency proof.
22. Order includes stale/partial shipping data but Shipment point read differs -> Shipment owns current provider shipping observation.
23. Order is canceled with real `cancel_detail` -> adapter does not infer reason from another status field.
24. Fraud-risk provider condition appears after payment -> fulfillment does not ship merely because earlier payment state looked ready.
25. Selected shipment requires invoice before label -> provider requirement remains open until authoritative provider reread shows closure.
26. Selected mode is provider-managed Full and no seller label action exists -> capability is provider-managed, not integration failure or fabricated seller action.
27. Provider SLA/deadline changes after earlier read -> current consequential fulfillment decision can detect/re-read rather than trusting stale time.
28. Full-only real products are available but no seller-writable stock lane exists -> do not falsely claim D0 Availability Control automatic convergence proof.
29. Operator selects internal Fulfillment Node path but only provider-operated shipments are present -> surface proof gap instead of calling provider execution internal fulfillment.
30. Return affects only part of an order -> Post-Sale Resolution can preserve partial scope; one Return does not close the entire Sale automatically.
31. Provider `available_actions` changes between decision and execution -> current provider capability is revalidated; no stale provider permission.
32. Return delivered physically but refund movement is unknown -> Post-Sale does not fabricate economic closure; B4 remains responsible for financial movement evidence.
33. Refund movement exists but physical return/consequence is unresolved -> Economics evidence does not close Post-Sale Resolution by itself.
34. Provider-native composite offer appears in selected proof -> either its material composition is explicitly supported under D0/D1 reopen rules or the path is unsupported; composition is never flattened silently.
35. A second marketplace arrives later -> B2's consumer-owned semantic ports remain usable without exporting `UserProduct`, `meli_facility`, ML moderation or ML Claim vocabulary into domains.
36. Adapter DTO misses a newly material provider field -> proof/observability must detect the contract gap without requiring universal unsanitized raw-payload retention.

---

## 14. Reopen / stop triggers

Reopen only the implicated accepted decision when real evidence requires it.

1. A Mercado Livre resource required for Product 1.0 cannot be represented as an Installation-qualified external reference without inventing a new canonical MPC identity -> targeted D2 review.
2. A necessary consumer dependency is absent from D1 -> STOP / targeted D1 review; do not hide it inside the adapter.
3. Provider effect semantics cannot fit D3 accepted/rejected/pending/ambiguous plus scope/reconciliation rules -> targeted D3 review.
4. Real Mercado Livre capability makes an accepted Product 1.0 outcome impossible or materially different -> targeted D0 review.
5. Selected first-flow proof has no seller-writable availability lane while D0 Availability Control is claimed complete -> STOP / targeted product-proof adjudication.
6. Selected internal Fulfillment Node proof cannot be exercised because the real Installation only exposes provider-operated fulfillment -> operator decides whether another seller-operated lane is required for first proof; do not mislabel provider work as internal work.
7. A real provider-native composite dependency changes availability or business-order materialization in the selected first flow -> use D0/D1's existing composition reopen trigger before modeling it in D4.
8. Price Automation or another provider control makes a required Price Intent path unavailable across the real proof set -> surface provider-effective unsupported/external-required; do not disable the provider control by convenience.
9. Current official/provider real behavior changes materially from the contracts above -> reopen only the affected B2 operation.
10. A second real provider exposes a repeated technical failure class that concrete provider-local adapters cannot handle without duplication -> consider the smallest shared mechanism then; never create a universal provider framework preemptively.

Current-code preference, package layout, old adapter shape and hypothetical future provider features are not reopen evidence.

---

## 15. Proposed batch outcome

Subject to independent review and the Installation Evidence Gate:

- **D0–D3: CURRENT STRUCTURE CONFIRMED; no reopen presently required.**
- **D4-B1: CURRENT STRUCTURE CONFIRMED.** B2 specializes its grounding without weakening it.
- **D4-B2: bounded concrete Mercado Livre target contracts, not a provider ontology.**
- **ADR-015:** proposed historical after B2 acceptance; its read-only-module/composite-ID/manual-refresh/absence=>closed target structure does not survive.
- **ADR-022/028:** current provider identifier evidence obligation satisfied by B2 while Readiness keeps correspondence authority.
- **No new Integration/Provider/OperatingMode/UserProduct/Warehouse/Claim/Return business domain/entity is created.**
- **No D5/D6/D7 topology or implementation is selected.**

Target shape:

```text
D1 owner + D2 identities + D3 intent/query semantics
                     ↓ consumer-owned semantic port
Mercado Livre adapter
  - bound Installation/seller
  - actual Item / UP / Catalog / stock / Order / Shipment / Claim topology
  - provider-effective capability / requirement
  - operation-specific coverage / preconditions / effect scope
  - authoritative reread / reconciliation
                     ↓
Mercado Livre API
```

The provider topology is real and preserved where correctness depends on it; it does not become MPC business ontology.

---

## 16. Independent reviewer focus

Please challenge independently, in particular:

1. Does B2-A…D cover every Mercado Livre responsibility D0 requires without expanding into ads/campaigns/SAC/reputation management/general reverse logistics?
2. Does any B2 rule move business authority from a D1 owner into D4?
3. Does Item/UP/Family/Catalog handling preserve D2 identity law without creating a hidden provider entity graph?
4. Is the User Product shared-effect-scope rule necessary D3 correctness or provider overfit?
5. Does price read-vs-write separation correctly reflect current ML behavior without prematurely designing B4 economics?
6. Is treating Price Automation as Level-2 provider-effective restriction correct, and is refusing automatic deactivation the smallest sustainable Product 1.0 choice?
7. Does stock mode sensitivity preserve Availability authority while correctly refusing seller writes to provider-managed stock?
8. Is 409 `x-version` correctly classified as stale-precondition rejection rather than ambiguity/retry machinery?
9. Does B2 accidentally imply a universal OperatingMode entity or implementation support for every ML logistics mode?
10. Is the first-flow proof rule derived correctly from D0, especially the statement that Full-only cannot prove seller/MPC-controlled availability convergence?
11. Does Sale/Order acquisition preserve D3 evidence/history semantics given provider search retention?
12. Is the historical Order `search_type=scan` measurement being kept at the right evidence strength rather than promoted to a permanent target contract?
13. Are Order and Shipment separated at the right semantic boundary, including provider deadline/fiscal/label readiness?
14. Does essential Post-Sale remain narrow enough to avoid CRM/SAC while still satisfying D0 cancellation/return/refund lifecycle?
15. Are Refund/payment financial movement semantics correctly deferred to B4?
16. Does Installation health/reputation/moderation evidence belong in B2 and feed Portfolio/Offering attention without creating reputation authority?
17. Is the Installation Evidence Gate genuinely necessary before canonical B2 ratification, or can a smaller evidence obligation prove the same Product 1.0 lane selection?
18. Do reference-platform comparisons support only the failure class, or has any benchmark accidentally become target authority?
19. Has B2 preserved YAGNI — seam now, modes only when real — or is there still unnecessary provider-specific detail?
20. What is the strongest counterexample that would force D0/D1/D2/D3 reopen?
21. Does any legacy ADR-015 meaning remain necessary that the proposed disposition would wrongly discard?
22. Does the candidate need a new target ADR, or is accepted D4 artifact authority sufficient for these contracts?

Reviewer findings are evidence, never authority. Do not modify canonical D4 directly. Append independent findings only to `AI-DIALOG.md` under the repository's append-only review protocol and hand back to GPT for adjudication.