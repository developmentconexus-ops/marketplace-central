# Architecture Map

```yaml
id: AM-001
type: architecture-map
status: planned
owner: Mission Strategist
parent: MIS-001
created: 2026-07-06
updated: 2026-07-06
validation_level: QA-0
lifecycle_scope: support
```

## Topology

```mermaid
graph TD
  web["apps/web React"] -->|SDK runtime only| api["apps/server_core Go API"]
  api --> modules["MPC business modules"]
  modules --> pg[("Postgres MPC-owned state")]
  modules --> caps["Marketplace capability ports"]
  caps --> ml["Mercado Livre adapter"]
  ml --> mlapi["Mercado Livre REST APIs"]
  modules --> sankhyaPort["Sankhya read ports"]
  sankhyaPort --> mnos["MNOS-derived read edge"]
  mnos --> oracle[("Sankhya Oracle read-only")]
  api --> integrations["integrations auth and capability health"]
  integrations --> pg
```

## Stock Seguro Flow

```mermaid
sequenceDiagram
  actor Operator
  Operator->>Web: open Stock Seguro
  Web->>API: GET /inventory/stock-risks
  API->>Inventory: list risk read model
  Inventory->>Postgres: read links, snapshots, policies
  Inventory-->>API: risk rows
  API-->>Web: linked listing risk view
  Operator->>Web: approve manual stock action
  Web->>API: POST /inventory/stock-actions/{id}/apply
  API->>Inventory: apply assisted action
  Inventory->>ProductLinks: verify resolved link
  Inventory->>Connectors: StockWriter.UpdateAvailableQuantity
  Connectors->>MercadoLivre: PUT /items/{item_id}
  MercadoLivre-->>Connectors: provider response
  Inventory->>Postgres: persist audit before/after/policy/response
  API-->>Web: applied or failed result
```

## Link Lifecycle

```mermaid
stateDiagram-v2
  [*] --> candidate
  candidate --> resolved: exact EAN/seller_sku or operator approve
  candidate --> conflict: multiple plausible matches
  candidate --> unresolved: no plausible match
  conflict --> resolved: operator selects one product
  unresolved --> resolved: operator links manually
  resolved --> rejected: operator rejects incorrect link
  rejected --> candidate: new listing snapshot changes evidence
```

## Stock Action Lifecycle

```mermaid
stateDiagram-v2
  [*] --> proposed
  proposed --> blocked: unresolved link, stale source, unsupported provider shape, ineligible product
  proposed --> approved: operator confirms
  approved --> applied: provider accepted update and audit persisted
  approved --> failed: provider rejected or request failed
  failed --> proposed: refreshed snapshot creates new recommendation
  applied --> [*]
  blocked --> [*]
```

## Build Order

```mermaid
graph LR
  M01["M-01 VTEX removal"] --> M02["M-02 capability framework"]
  M02 --> M03["M-03 MNOS/Sankhya read contract"]
  M02 --> M04["M-04 Product Links ML"]
  M03 --> M05["M-05 Stock Seguro ML"]
  M04 --> M05
  M05 --> M06["M-06 Orders + Margin ML"]
  M06 --> M07["M-07 Commercial Intelligence"]
```

## Truth Notes

- Interface contracts own field names, route namespaces, error codes, and capability semantics.
- This map visualizes mission contracts only; feature execution must preserve the referenced contracts.
