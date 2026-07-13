# M-06 read-only architecture discovery: ML order/item → Sankhya 313/306

Date: 2026-07-12/13 UTC  
Scope: architecture recommendation only; no Oracle/provider write, production
change, PII, or secret access.

## Decision summary

Use an immutable two-grain external reference created before/manual-with TOP
313:

- order key: `ml:v1:<provider-account-id>:<provider-order-id>`;
- line key: `ml-line:v1:<provider-account-id>:<provider-order-id>:<mpc-line-id>`.

Persist the order key on the TOP 313 `TGFCAB` row and the line key on each TOP
313 `TGFITE` row. Treat `TGFVAR` (`NUNOTAORIG`,`SEQUENCIAORIG` →
`NUNOTA`,`SEQUENCIA`) as canonical 313 → 306/invoice lineage. Copying the
keys onto 306 is useful for direct lookup but is denormalized convenience, not
the lineage source. MPC must keep an append-only mapping ledger with uniqueness
constraints and never infer links from partner/date/product.

## Facts

### Mercado Livre (current official documentation via Context7)

- An order has stable `id`, optional `pack_id`, `shipping.id`, timestamps,
  status/tags, `payments[].id` with `payments[].order_id`, and `order_items[]`
  containing `item.id`, optional `variation_id`, quantity, unit price, and sale
  fee. The reviewed response does not expose a distinct order-item line ID.
- `GET /packs/{pack_id}` returns a pack with its own ID, shipment ID, and a list
  of order IDs. A pack can therefore group orders and is not an order identity.
- Change/return resources identify the original order (`resource_id`), claim,
  returned item/variation/quantity, return ID, and possible new order IDs.
- Sources:
  - https://developers.mercadolivre.com.br/pt_br/gerenciamento-de-vendas
  - https://developers.mercadolivre.com.br/pt_br/gestao-packs
  - https://developers.mercadolivre.com.br/pt_br/trocas-change

### Sankhya (current primary documentation)

- `TGFCAB.NUNOTA` is the document header key; `TGFITE` is keyed by
  (`NUNOTA`,`SEQUENCIA`).
- `TGFVAR` relates origin and destination document lines and explicitly handles
  one order line being invoiced into multiple destination notes via
  `QTDATENDIDA`.
- The sales-order API accepts additional fields, so a configured custom field
  can be supplied at document creation.
- Sankhya's add-on guidance recommends a unique product prefix, forbids using
  reserved `AD_` in add-on schema definitions, supports extending `TGFCAB` for
  simple fields, and recommends normalization/foreign keys for relationships.
- Sources:
  - https://developer.sankhya.com.br/docs/operacoes-comerciais
  - https://developer.sankhya.com.br/reference/get_ligacaonotapedido
  - https://developer.sankhya.com.br/reference/addpedido
  - https://developer.sankhya.com.br/docs/guia-de-boas-praticas

### Common hub pattern (Bling primary documentation)

- Bling sales-order webhooks distinguish internal order `id`/`numero` from
  store reference `numeroLoja` and store/channel `loja.id`. This supports the
  common pattern “internal ERP identity + immutable external channel identity,”
  rather than buyer/date/product correlation.
- Source: https://developer.bling.com.br/webhooks

### Repository truth

- MPC already keys orders by (`tenant_id`,`installation_id`,`provider_order_id`).
- Items are currently replace-written and keyed by positional `line_no`; this
  is not yet an immutable provider/MPC line identity.
- The ML adapter captures order ID, shipment ID, payments and item/variation,
  but not `pack_id`; it exposes no stable order-item line key.
- F-07 requires exact Oracle (`NUNOTA`,`SEQUENCIA`) for tax and correctly keeps
  tax missing when that identity is absent.
- Paths:
  - `apps/server_core/migrations/0027_orders_marketplace_orders.sql`
  - `apps/server_core/internal/modules/orders/adapters/postgres/order_repo.go`
  - `apps/server_core/internal/modules/connectors/adapters/mercado_livre/capability_adapter.go`
  - `apps/server_core/internal/modules/internal_read/domain/internal_tax.go`

## Inferences and recommended model

### Identifier/grain

1. `provider_order_key` is installation/account scoped. Provider order ID alone
   is not globally sufficient across accounts/providers.
2. Because reviewed ML docs expose no line ID, MPC assigns `mpc_line_id` once
   on first ingestion and never recomputes it from mutable position. Existing
   `line_no` may seed a migration only after identity review; future refreshes
   update a persisted line, not delete/recreate its identity.
3. Item ID, variation ID, seller SKU, quantity and price are validation/snapshot
   attributes, not the line identity. Duplicate identical item rows must either
   retain their previously assigned MPC IDs or block as ambiguous.

### Sankhya placement and lifecycle

- TOP 313 header: one custom text field for `provider_order_key`.
- TOP 313 item: one custom text field for `provider_line_key`.
- At 313 creation, MPC/manual UI reserves the keys in an append-only MPC ledger
  before the Sankhya action and records returned 313 `NUNOTA`/`SEQUENCIA`.
- TOP 306/invoice tax lookup follows `TGFVAR` from each exact 313 line. Partial
  invoicing yields one-to-many destination lines; sum tax only over those exact
  descendants and their realized quantities. Never query product/date totals.
- Optional copy rules may place the same keys on 306 `TGFCAB`/`TGFITE`, but
  reads must still validate them against `TGFVAR`.

### Uniqueness/idempotency

- MPC unique constraints:
  - one order ledger row per (`tenant_id`,`installation_id`,`provider_order_id`);
  - one line per (`tenant_id`,`installation_id`,`provider_order_id`,`mpc_line_id`);
  - one active 313 origin mapping per provider order/line key;
  - destination mappings unique by (`NUNOTA`,`SEQUENCIA`) and append-only.
- Idempotency key: `sankhya-313-create:<provider_order_key>`. A retry reads and
  reconciles the existing ledger/Sankhya key; it never creates a second 313.
- Sankhya-side uniqueness should be a real unique constraint/index or a small
  normalized integration table if the customer's supported customization path
  cannot make the two native extension fields unique. Read-before-create alone
  is insufficient under concurrency.

### Pack, shipment, payment and partner

- Store `pack_id`, `shipment_id`, and `payment_id` as separate child/external
  resource identities related to the order; none may replace the order key.
- A pack may contain multiple order IDs; a payment can be independently updated
  or refunded; a shipment can represent logistics grouping.
- Partner (`CODPARC`) is an operational choice/attribute. The order key remains
  stable if the partner is corrected. Two orders from the same buyer are always
  distinct because partner/date/product are never part of the key.

### Cancellation, return and exchange

- Preserve the original order/line keys and append lifecycle events; do not
  delete or reuse them.
- Before invoicing: cancel/release TOP 313 per approved business policy while
  retaining the ledger.
- After 306: trace reversal/devolution through Sankhya lineage and record
  separate negative/refund tax facts. ML claim/return/change IDs and new order
  IDs are separate references linked to the original order, never mutations of
  its identity.

### Existing manually recreated orders

- No automatic backfill by buyer/partner/date/product/value.
- Build an operator-assisted candidate screen using those fields only to narrow
  candidates; require explicit selection of exact ML order and every 313 line,
  persist actor/reason/time and before/after audit, then attach immutable keys.
- Ambiguous, split, combined or unconfirmed cases stay `unlinked`; tax remains
  missing. Candidate confidence is not linkage evidence.

## Options considered

1. Header-only ML order ID: rejected; cannot prove item tax line, partial
   invoicing, duplicate item rows, or returns at item grain.
2. Partner/date/product matching: rejected; owner explicitly confirms no
   deterministic relationship and same-buyer orders can collide.
3. Custom fields only, read-before-create: insufficient unless Sankhya can
   enforce uniqueness and preserve/copy fields through invoicing.
4. **Recommended:** two custom references on 313 plus an MPC append-only unique
   ledger and `TGFVAR` lineage; add a small Sankhya integration table only if
   supported field/index configuration cannot enforce uniqueness.

## Owner/admin-specific decisions still required

The architecture cannot name these from repository or public documentation:

1. Exact customer-supported custom field names, prefix, types and lengths on
   `TGFCAB` and `TGFITE` (examples above are conceptual, not deployable names).
2. Whether the environment supports unique indexes/constraints on those native
   extensions; otherwise approve a normalized custom integration table and its
   supported deployment mechanism.
3. Whether TOP 313 → 306 copies custom header/item fields automatically or
   needs a configured copy rule/event; `TGFVAR` remains canonical regardless.
4. Which sanctioned manual-entry surface must require the keys at TOP 313
   creation (Sankhya field, action, add-on, or MPC-assisted form) and who owns
   administration/deployment.
5. Cancellation/devolution TOPs and business policy for pre/post-invoice cases.

These are configuration/authority decisions, not facts that code may invent.
