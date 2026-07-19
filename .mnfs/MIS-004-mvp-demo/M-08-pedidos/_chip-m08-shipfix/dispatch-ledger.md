# CHIP-M08-SHIPFIX — Dispatch Ledger

Branch `chip/m08-shipfix` · base `8b9e0ec01cae7b64f68001054464b3cb7f534160`

| Seq | Step | Result |
|-----|------|--------|
| 1 | Read dispatch spec + HARNESS-CORE §4/§5 + profile ladder bindings | scope A/B/C confirmed |
| 2 | Failing-test-first: costs-legacy-decode + x-format-new header (connector); bucket truth table; tags_json scan; FE KPI==Lista | written |
| 3 | Scope A — x-format-new on /costs only; degrade on ProviderPayloadInvalid; substatus capture | shipping_reader.go, capability_adapter.go, shipping_read.go |
| 4 | Scope B — DeriveOrderBucket(providerStatus, shipmentStatus, tags, hasShipment); thread both call sites | order_bucket.go, read_model.go, order_repo.go, enrich_service.go, http_handler.go |
| 5 | Scope C — FE-derived KPI from list+bucketTabCount; honest "—" | PedidosPage.tsx, pedidosTabs.ts, PedidoDrawer.tsx |
| 6 | Contract additive: OrderRastreio.substatus + SDK | openapi.yaml, sdk-runtime/index.ts |
| 7 | L0/L1 ladder (GOCACHE abs) | Go build/vet/test ./... EXIT=0; web vitest 312✓; SDK tsc✓; vite build✓ |
| 8 | Commit green slice | `aa4bb11` |
| 9 | Governance lane vs BaseSha (clean worktree, 40-hex) | `status=passed` (additive, no drift) |
| 10 | P6 dual gate: cold Opus + adversarial sonnet | PASS / PASS, agreement, no blockers |
| 11 | Actioned 1 valid MINOR (happy-path substatus assertion) | ML pkg re-run green |
| 12 | Commit gate-fix + evidence | `741439d` |
| 13 | CLOSED event → hub (HAND-BACK: hub owns live P7 re-drive) | emitted |

## Operator-surfaced scope addition (reported to hub)

Mid-session the operator relayed ML-AI guidance to incorporate order `tags`
(`delivered`) and `GET /orders/{id}/shipments` status+substatus. Incorporated:
the `delivered` tag is now the robust, shipment-lookup-independent bucket signal
(order-level, survives a failed live /shipments read), requiring
`OrderReadModel.Tags` + read-query threading (no migration — `tags_json` column
pre-exists at migration 0027). Flagged here as a scope addition beyond the
original A/B/C for hub ratification.

## Non-negotiables honored

Zero ML writes · never booted server / bound :8080 / read .env · GOCACHE
absolute · no push · commits via `git commit -F` · evidence written before CLOSED.
