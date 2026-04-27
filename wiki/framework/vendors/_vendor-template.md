# <Vendor> Playbook Template

Last verified: YYYY-MM-DD

## 1. Vendor Snapshot

- Auth model:
- API base domains:
- Supported markets:
- Primary capability groups:

## 2. MPC Integration Mapping

- `integrations` responsibilities:
- `connectors` responsibilities:
- Domain module consumers:

## 3. Capability Matrix

| Capability | Shopee/Provider API Area | MPC Module Owner | Status |
|---|---|---|---|
| auth_installation |  | integrations | planned |
| catalog_sync |  | connectors + catalog | planned |
| stock_price_update |  | connectors + catalog/pricing | planned |
| orders_sync |  | connectors + orders | planned |
| returns_refunds |  | connectors + orders | planned |
| logistics_ops |  | connectors + orders | planned |
| push_ingestion |  | connectors + messaging/orders | planned |

## 4. Getting Started Checklist

1. Developer account approved
2. App created
3. Auth callback validated
4. Token refresh path validated
5. Sandbox test evidence captured

## 5. Best Practices Checklist

- Idempotency strategy
- Rate limit strategy
- Retry/backoff strategy
- Signature/credential rotation strategy
- Observability and alerting strategy

## 6. Sources

- List source endpoint:
- Document detail endpoint:
- Primary docs:
