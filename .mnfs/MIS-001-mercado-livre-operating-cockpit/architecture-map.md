# Architecture Map — MVP Replan

```yaml
id: AM-001
type: architecture-map
status: planned
owner: Mission Strategist
parent: MIS-001
created: 2026-07-06
updated: 2026-07-13
validation_level: QA-0
lifecycle_scope: support
```

This file visualizes `mission.md` and IC-003. Those artifacts remain authoritative.

## Runtime topology

```mermaid
graph TD
  operator["Trusted local operator"] --> web["apps/web operator workspaces"]
  web -->|"packages/sdk-runtime only"| api["apps/server_core Go API"]
  api --> catalog["catalog / product_links / inventory / orders / profitability / integrations"]
  catalog --> pg[("PostgreSQL MPC state")]
  catalog --> oraclePort["internal_read ports"]
  oraclePort --> oracle[("Sankhya Oracle read-only")]
  catalog --> providerPorts["marketplace capability ports"]
  providerPorts --> ml["Mercado Livre REST read adapters"]
  web -. "stock preview only" .-> preview["No provider mutation"]
```

## Operator journey

```mermaid
sequenceDiagram
  actor Operator
  participant Overview
  participant Product as Product 360
  participant Listing
  participant Sale
  participant API
  participant Sources as PostgreSQL + ML + Oracle
  Operator->>Overview: open attention queue
  Overview->>API: existing SDK list operations
  API->>Sources: read persisted and external-source facts
  Sources-->>API: values + quality + observed_at
  API-->>Overview: attention items
  Operator->>Listing: open stock-attention deep link
  Listing->>Product: inspect canonical product
  Product->>Listing: return to linked listing
  Listing->>Sale: open related sale
  Sale-->>Operator: revenue, inputs, margin, provenance
  Operator->>Listing: review stock simulation
  Listing-->>Operator: current + proposed + rule + preview payload, executed=false
```

## Shared state vocabulary

```mermaid
stateDiagram-v2
  [*] --> current
  current --> stale: freshness threshold exceeded
  current --> unknown: required source absent
  current --> conflict: identities or facts disagree
  stale --> current: successful refresh
  unknown --> current: source becomes known
  conflict --> current: operator resolves evidence
```

## Build order

```mermaid
graph LR
  accepted["Passed M-01..M-05 + M-08"] --> M09["M-09 Canonical Product Foundation"]
  M06["M-06 failed historical gate; reusable order/margin evidence"] -.-> M13
  M09 --> M13["M-13 Integrated Operator Workspaces"]
  M13 --> M14["M-14 Real Vertical MVP Validation"]
  M14 --> M07["M-07 Commercial Intelligence reassessment"]
  M14 -. "post-MVP" .-> M10["M-10 Runtime consolidation"]
  M10 -.-> M11["M-11 Production write execution/auth"]
  M11 -.-> M12["M-12 Listing/SKU mutation"]
```
