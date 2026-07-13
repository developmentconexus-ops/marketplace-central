# Marketplace Hub UX Benchmark

```yaml
id: R-003
type: research
status: verified
owner: External Researcher
parent: MIS-001
created: 2026-07-13
updated: 2026-07-13
validation_level: QA-0
lifecycle_scope: support
```

## Topic

Public UX and information-architecture patterns for Brazilian marketplace hubs and ERPs.

## Sources Checked

- Magis5 Hub public site and help center: dashboard, product, listing, stock, order,
  pricing, and Sankhya integration categories. Verified 2026-07-13.
  - https://ajuda.magis5.com.br/magis5-hub
- Bling help center: Product versus Mercado Livre Listing, sales dashboard, and
  Margin by Order. Verified 2026-07-13.
  - https://ajuda.bling.com.br/hc/pt-br/articles/34677719147415-Gest%C3%A3o-de-An%C3%BAncios-do-Mercado-Livre-no-Bling-o-que-mudou
  - https://ajuda.bling.com.br/hc/pt-br/articles/26725573512343-Dashboard-Pedidos-de-Venda
  - https://ajuda.bling.com.br/hc/pt-br/articles/35658275718551-Dashboard-Margem-por-Pedido
- ANYMARKET help center: centralized attribute links and product/listing workflows.
  Verified 2026-07-13.
  - https://suporte.anymarket.com.br/hc/pt-br/articles/51175297537555-Como-funciona-a-Central-de-v%C3%ADnculos
- Hub2b public capabilities: catalog, price, stock, orders, and integration health.
  Verified 2026-07-13.
  - https://hub2b.com.br/funcionalidades

Public sources expose documentation and selected screenshots, not authenticated product
areas. Layout conclusions below are evidence-based patterns plus explicit inference.

## Findings

- **Verified:** Bling treats internal Product and channel Listing as separate objects;
  a product owns stock/cost/description while a listing owns offer-specific title,
  price, shipping, and strategy.
- **Verified:** Bling sales views join channel, order, product, cost, freight, tax,
  fee, markup, and contribution instead of making margin a detached utility.
- **Verified:** Magis5 groups operational attention, products/listings, stock, orders,
  and financial reports inside one hub taxonomy.
- **Verified:** ANYMARKET provides centralized cross-marketplace linking so the
  operator does not need to traverse one module per marketplace.
- **Inference:** A useful cockpit needs both entity context (one product/listing/sale)
  and aggregate workspaces for batch filtering; either scale alone is insufficient.
- **Inference:** Dashboards should route to the exact pending objects, not only show
  charts or totals.

## Recommendation

Use Overview, Products/Product 360, Listings, Sales/Margin, and Operations. Keep
product and listing identities distinct, expose the link in both directions, and
show stock/margin simulations in context. Do not copy fiscal, fulfillment, customer
service, multi-company, or multi-role ERP breadth into the MVP.

## Impact On Mission

The evidence replaces module-shaped navigation with object-centered workspaces,
adds a global attention queue, and makes a vertical real-data browser journey the
mission acceptance target.

## Handoff

- Current status: Verified from public primary documentation with access limitations noted.
- Next owner: Mission Strategist.
- Next action: Keep route and object semantics aligned with IC-003.
- Required files/evidence: this note, `mission.md`, and `research/mvp-operator-workspace-interface-contract.md`.
- Blockers or open decisions: None - private authenticated screens were not used as evidence.
