# MIS-003 Architecture Map

View of the contracts (IC-02..05, mission.md ADRs). Not a source of truth.

## Topology

```mermaid
graph TD
  web["React SPA apps/web :5174"] -->|"sdk-runtime HTTP"| core["server_core :8080"]
  web -.->|"IC-05 routes /anuncios /catalogo /vinculos /estoque /precos /pedidos /integracoes /protocolos"| web
  core --> pg[("PostgreSQL: listings, mutation_protocols, mutation_items, market_* (empty)")]
  core -->|"internal_read ports (read-only, governed)"| oracle[("Oracle/Sankhya")]
  core -->|"connectors capability adapter"| ml["Mercado Livre API"]
  subgraph core_modules["server_core modules"]
    listings["listings (IC-02) GET /listings*"]
    mutations["mutation envelope (IC-03) /mutations*"]
    market["market (IC-04) /market* contract-only"]
    existing["existing: catalog, product_links, inventory, orders, profitability, pricing, marketplaces, integrations, connectors, internal_read"]
  end
  mutations -->|"poller applies via capability adapters"| ml
  listings -->|"ingestion via connectors ListListings"| ml
```

## Behavior flow — corrigir preço em massa (IC-03)

```mermaid
sequenceDiagram
  actor Op as Operador
  Op->>SPA: seleciona por filtro, "Atualizar preço"
  SPA->>API: POST /mutations {type: price_update, selection, intent}
  API-->>SPA: 201 draft
  SPA->>API: POST /mutations/{id}/preview
  API->>PG: snapshot selection -> mutation_items (before capturado)
  API-->>SPA: previewed + totals
  Op->>SPA: confirma
  SPA->>API: POST /mutations/{id}/approve {execute: true}
  API-->>SPA: approved
  loop poller (in-process, chunks of 20)
    API->>ML: write price (idempotency_key)
    ML-->>API: ok | erro
    API->>PG: item applied | failed {failure.code}
  end
  SPA->>API: GET /mutations/{id} (poll 2s)
  API-->>SPA: applied | partially_failed | failed_preserved
  SPA-->>Op: protocolo MP-nnnnnn com antes/depois por item
```

## Lifecycle — MutationProtocol (IC-03 enum, exact)

```mermaid
stateDiagram-v2
  [*] --> draft
  draft --> previewed: preview
  draft --> cancelled: cancel
  previewed --> approved: approve execute=true
  previewed --> cancelled: cancel
  previewed --> previewed: re-preview
  approved --> applying: poller claim
  applying --> applied: all items ok
  applying --> partially_failed: mixed
  applying --> failed_preserved: all failed
  applied --> [*]
  partially_failed --> [*]
  failed_preserved --> [*]
  cancelled --> [*]
```

Retry never re-enters this machine: it clones retryable failed items into a NEW protocol.

## Build order (mission.md Milestone Strategy, exact)

```mermaid
graph LR
  M01["M-01 listings-read-spine"] --> M02["M-02 frontend-platform-anuncios"]
  M01 --> M03["M-03 mutation-envelope-writes"]
  M02 --> M04["M-04 read-workspaces"]
  M03 --> M04
  M02 --> M05["M-05 visao-geral-pedidos-sync-central"]
  M03 --> M05
  M04 --> M05
  M03 --> M06["M-06 corrigir-atributo-market-contracts"]
  M04 --> M06
```
