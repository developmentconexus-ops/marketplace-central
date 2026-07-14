# Research: Oracle Runtime Dimensioning & Surprise Register

```yaml
id: RES-01
type: research-note
status: planned
owner: Mission Strategist
parent: MIS-002
created: 2026-07-13
updated: 2026-07-13
lifecycle_scope: support
```

Consolidated evidence from 6 read-only investigation agents (2026-07-13 session) plus runtime walkthrough. Facts verified against working tree at planning time.

## Evidence base (verified, session 2026-07-13)

- Adapter internals: godror v0.51.0 (`apps/server_core/go.mod:6`); pool default 4 (`internal_read/adapters/oracle/config.go:13`); uncommitted hardening refactor in working tree (ConnectionParams, CallTimeout 30s per call, `Database` interface, validated envs); no retry, no Close on shutdown, no fetch tuning; `.gomodcache` drift (v0.49.4 on disk).
- Consumption path: catalog 1+3N sequential (`catalog/adapters/internalread/reader.go:67-111`; unbounded `FindProductsForLinking` at `oracle/reader.go:386-428`); profitability per item×tax-source (`profitability/application/service.go:239-287`); inventory per snapshot (`inventory/application/stock_risk_service.go:29-84`); Sankhya `ValidateConfiguration` (ping+2 metadata queries) per call + per-candidate line query (`sankhya_linkage_reader.go:104,150-160`); `FreshnessPolicy` plumbed, never read; zero cache, zero concurrency control; no server/transport timeouts (`cmd/server/main.go:35`).
- Frontend: React 19 + Vite; `packages/sdk-runtime` hand-written fetch wrapper; NO TanStack/SWR; refetch per mount.
- Reference (MetalShopping_Final): 1 JOIN-based snapshot query per entity, zero N+1, batch worker outside HTTP path — proves same ERP answers set-based queries fine; but no pagination, unbounded memory, dead batch-size config there.
- Scale anchor: MetalShopping acceptance recorded 430,794 promotable price rows / 488,420 promoted from the SAME Sankhya ERP → catalog assumed 30k–100k active products (MEASURE in M-01).

## Dimensioning parameters

| Parameter | Assumed | Status |
| --- | --- | --- |
| Active catalog (TGFPRO) | 30k–100k | measure M-01 (COUNT via live lane) |
| Concurrent operators | 2–10 | internal cockpit |
| Network RTT to Oracle | 1–5ms LAN | measure M-01; WAN would 10× round-trip cost |
| Page size | 50 | IC-01 |
| Import size | 20–1000 orders × 1–10 items | ceiling 200/request (ADR-05) |

Key property: keyset pagination makes per-screen cost O(page), not O(catalog); architecture holds 1k→500k products. Danger concentrates in batch flows and Oracle view plans.

## User journeys (simulated against target architecture)

- J1 catalog list cold: 1 query/page ≈ 100–400ms; 5 simultaneous operators collapse to 1 Oracle query via singleflight. OK.
- J2 back-navigation: TanStack memory, zero Oracle within staleTime. OK.
- J3 text search: LIKE '%x%' full-scan risk → debounce + FETCH FIRST 50 + DBA index review (R5).
- J4 stock risks: 1 IN-list query for 50 snapshots; TTL 45s acceptable with force-refresh. OK.
- J5 linkage: cached config validation + 1 JOIN candidates+lines = 1–2 round trips, always fresh. OK.
- J6 margin import: batch cost/tax by IN-list; needs chunk 500 (R2), batch deadline 120s (R3), semaphore 4 (R4).
- J7 link-candidate generation: O(listings) filtered queries; acceptable ≤2k listings (R7).
- J8 staleness UX: `as_of` + force-refresh (R6).
- J9 Oracle down: existing typed unavailability preserved; cache serves within TTL during blips; no retry (declined).
- J10 pricing batch: bottleneck is external freight API, out of scope.

## Surprise register → mission risks R1–R10

Mapped 1:1 into `mission.md ## Risks`. Load-bearing preconditions: R1 (EXPLAIN PLAN gate, M-01-C04), R8 (RTT measurement), R10 (real COUNT).

## Version-sensitive claims

- `@tanstack/react-query` v5 (React 19 compatible) — verify-at-install.
- `golang.org/x/sync` already an indirect dep (`go.mod:19`) — verified 2026-07-13; promote to direct when singleflight imported.
- godror fetch/prefetch tuning API (`godror.FetchArraySize`) — verify-at-install (M-02, optional tuning only).
- Oracle IN-list limit 1000 (ORA-01795) — Oracle-documented, stable behavior; chunk at 500 regardless.
