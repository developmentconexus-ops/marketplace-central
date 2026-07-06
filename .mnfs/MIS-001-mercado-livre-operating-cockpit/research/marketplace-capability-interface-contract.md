# Interface Contract

```yaml
id: IC-001
type: interface-contract
status: planned
owner: Mission Strategist
parent: MIS-001
created: 2026-07-06
updated: 2026-07-06
validation_level: QA-0
lifecycle_scope: support
```

## Boundary

MPC business modules to marketplace connector adapters.

## Why This Contract Exists

Two feature workers must not invent incompatible provider abstractions. MPC needs native business operations that work for Mercado Livre now and future marketplaces later, while still allowing each provider adapter to map provider-specific endpoints and constraints.

## Resources Or Entities

- `ProviderCode`: stable marketplace code, first value `mercado_livre`.
- `ProviderAccountRef`: tenant-scoped installation/account identity.
- `ProviderListingRef`: provider listing id plus optional variation id.
- `ListingSnapshot`: provider listing state normalized for business modules.
- `StockRead`: provider-reported stock for a listing/variation with timestamp.
- `StockWriteRequest`: requested stock write with idempotency key and audit context.
- `OrderSnapshot`: provider order with item, fee, payment, shipment, and cancellation fields.

## Operations

| Operation | Trigger | Input | Output | Notes |
| --- | --- | --- | --- | --- |
| ListListings | Scheduler or operator refresh | `provider_code`, `account_ref`, optional cursor/filter | listing snapshots newest provider update first | Must include provider item id, title, status, SKU/EAN when available, variations |
| ReadListing | Link resolution or action validation | `provider_listing_ref` | one listing snapshot | Must preserve provider source timestamp when available |
| ReadStock | Stock refresh | `provider_listing_ref` | stock read | Must state whether stock is item-level or variation-level |
| UpdateAvailableQuantity | Manual approved stock action | stock write request | stock write result | Must be idempotent by action id; no auto-write caller in this mission |
| ListOrders | Scheduler refresh | account, date/status cursor | order snapshots newest update first | Must be idempotent by provider order id and provider last update |
| ReadOrder | Order refresh or notification resource | provider order id | order snapshot | Must include items, sale fee when available, payments, shipping id, cancellation detail |

## Fields

### Required Inputs

- `tenant_id`: string, non-empty.
- `provider_code`: string, first implemented value `mercado_livre`.
- `installation_id`: string, identifies credential/capability source.
- `provider_account_id`: string, provider seller/user id when known.
- `provider_item_id`: string for listing operations.
- `provider_variation_id`: nullable string for variation-level operations.
- `idempotency_key`: string for writes, equal to MPC stock action id.

### Required Outputs

- `provider_code`: string.
- `provider_item_id`: string.
- `provider_variation_id`: nullable string.
- `provider_status`: string copied or mapped from provider.
- `seller_sku`: nullable string.
- `ean`: nullable string.
- `title`: string.
- `available_quantity`: integer or null when unsupported/unknown.
- `source_updated_at`: nullable timestamp.
- `fetched_at`: timestamp.
- `raw_provider_ref`: string or object reference suitable for audit, not a raw secret payload.

## Enums And Statuses

- Capability status: `supported`, `unsupported`, `degraded`, `blocked`.
- Stock write result: `applied`, `rejected`, `transient_failure`, `unsupported_shape`.
- Provider operation state: `pending`, `running`, `succeeded`, `failed`, `skipped`.

## Error Cases

- Missing installation.
- Credential unavailable or expired.
- Provider rate limited.
- Provider validation rejected request.
- Unsupported listing shape.
- Provider transient failure.
- Provider payload missing required identity.

## Persistence Expectations

Adapters do not own business state. They may persist operation-run technical evidence only through existing integrations/operation facilities. Product links, stock snapshots, actions, orders, and profit snapshots belong to their business modules.

## Canonical Examples

Success stock read:

```json
{
  "provider_code": "mercado_livre",
  "provider_item_id": "MLB123",
  "provider_variation_id": null,
  "available_quantity": 4,
  "provider_status": "active",
  "seller_sku": "12345",
  "ean": "7890000000000",
  "title": "Produto exemplo",
  "source_updated_at": null,
  "fetched_at": "2026-07-06T15:00:00Z"
}
```

Rejected write result:

```json
{
  "idempotency_key": "stockact_001",
  "result": "unsupported_shape",
  "provider_code": "mercado_livre",
  "provider_item_id": "MLB123",
  "provider_variation_id": "987",
  "message": "variation stock update shape is not supported by this adapter version"
}
```

## Error Matrix

| Case | Status | Code | Notes |
| --- | --- | --- | --- |
| Installation not found | 404 | `INTEGRATIONS_INSTALLATION_NOT_FOUND` | Business service blocks before adapter call when possible |
| Credential unavailable | 409 | `INTEGRATIONS_CREDENTIAL_UNAVAILABLE` | No provider call is attempted |
| Provider rate limited | 503 | `CONNECTORS_PROVIDER_RATE_LIMITED` | Retryable with backoff |
| Provider validation rejection | 422 | `CONNECTORS_PROVIDER_VALIDATION` | Persist provider response in audit |
| Unsupported listing shape | 409 | `CONNECTORS_PROVIDER_UNSUPPORTED_SHAPE` | Stock action remains blocked |
| Provider transient failure | 503 | `CONNECTORS_PROVIDER_TRANSIENT` | Retryable, no duplicate write without same idempotency key |
| Missing provider identity | 502 | `CONNECTORS_PROVIDER_PAYLOAD_INVALID` | Adapter maps malformed response to structured error |

## Database Shape

- Business modules define tables; connector adapters do not create listing/order/stock business tables.
- Provider operation evidence may reference existing operation-run tables if needed.
- All business tables must include `tenant_id`.

## Seed Data

- `provider_code = mercado_livre`.
- One connected Mercado Livre installation fixture.
- One listing without variation and one listing with variation for adapter contract tests.

## Timestamp And ID Semantics

- Provider ids are strings, even when numeric-looking.
- Timestamps are RFC3339 UTC in MPC APIs.
- `fetched_at` is MPC collection time.
- `source_updated_at` is provider update time when documented/available.

## Compatibility Rules

- New providers may add provider-specific metadata, but cannot remove required normalized fields.
- Unsupported capabilities must return `unsupported`/structured errors, not nil behavior.
- Business modules must use capability interfaces, not provider adapter packages directly.

## Route Namespace

- Connector operational routes are internal to server composition.
- User-facing routes for this mission are owned by business modules:
  - `/product-links/*`
  - `/inventory/*`
  - `/orders/*`
  - `/profitability/*`
- Existing integration routes remain under `/integrations/*`.

## Transport And Integration

- Provider secrets remain in integrations credential storage.
- Adapters receive credential material only from application services/ports and never return secrets.
- Browser never calls provider APIs directly.

## Must Preserve

- Capability ports are small and operation-specific.
- Provider-specific payloads stay at adapter boundary.
- Business rules stay out of adapters.

## Must Not Decide In Feature Execution

- Do not add a provider-specific business service for Mercado Livre stock or orders.
- Do not place stock safety policy inside connector adapter code.
- Do not introduce direct frontend/provider calls.

## Validation Impact

- Boundary tests must prove business modules compile against ports and do not import provider HTTP packages.
- Adapter tests must prove Mercado Livre maps documented item/order shapes into normalized snapshots.
