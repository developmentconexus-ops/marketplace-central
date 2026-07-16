# M-06 F-02 plan — market-contract-module

Authoritative slice cards: `../../CHIP-SAT-P2-PLAN.md` §2 (planner: gpt-5.6-sol medium,
OS-process, ledger row 1). Slices F02-M1..M7, C1..C2, R1, I1.

## Write-DAG (from plan)

- M1 → M2 → {M3(complex), M4(complex, after M3's shared constructor file), M5}
- {M3,M4,M5} → M6 → M7 → R1 → I1
- M7 → C1 → C2 (same OpenAPI/SDK files, serialized)
- F-02 starts only after F-01's OpenAPI/SDK/root slices complete (same-file serialization
  inside chip).

## Migrations

0043_market_observations.sql + 0044_market_references.sql; 0045 reserved unused (never an
empty file). Beyond 0045 = REQUEST.

## Hard gates

Contract-only: no production CollectorPort impl, no seed, no scheduler; test double in `_test`
only; grep-proof in I1. 6-signal separation; stored `source`/`captured_at` NOT NULL.
