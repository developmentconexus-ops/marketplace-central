# M-05 F-01 plan — aggregate-sync-endpoints

Authoritative slice cards: `../../CHIP-SAT-P2-PLAN.md` §1 (planner: gpt-5.6-sol medium,
OS-process, ledger row 1). Slices F01-O1..O3, S1..S2, D1..D5, C1..C2, R1, I1.

## Write-DAG (from plan)

- O1 → O2(complex) → O3 → C1; O2 → D1
- S1(complex) → S2 → C2; S1 → D3
- D2 independent
- D4 waits {D1, D2, D3, listings summary}; D4 → D5 → R1 → I1
- C1 → C2 serialized (same OpenAPI/SDK files)

## Blocking gate — RESOLVED (hub ruling 2026-07-16, landed main af83547a)

Option 2, honest-null this mission (ADR-17):
- buyer/currency/fulfillment/NF-state: null end-to-end, never defaulted; total = payments sum
  where computable, else null.
- `filter.fulfillment` → 400 `unsupported_filter` (never silently ignored).
- No migration block for F-01; adapter-ingestion canonical columns = successor scope.

Pins granted: (a) canonical order timestamp = `provider_created_at`; null provider_created_at
EXCLUDED from date-filtered/counter results, never bucketed as today. (b) dashboard field names
+ degraded[] {listings,linkage,orders,sync} per plan §1. (c) sync `filter.module` = exact
stored operation_type values.

## Deferred to hub

- Vite proxy rows `/dashboard`, `/sync` (CHIP-SAT is zero-frontend; frontend-seam owner adds).
