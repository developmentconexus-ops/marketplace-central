# F-02-market-contract-module — spec

```yaml
id: F-02
type: feature-spec
status: in-progress
owner: CHIP-SAT
parent: M-06
created: 2026-07-16
branch: chip/sat-m05f01-m06f02
governance_anchor: a49168e641ffd6f61932ca57c29b1d1bdcde2fb0
```

## Binding inputs

- IC-04 `../../research/market-data-interface-contract.md` — verbatim shapes, CollectorPort
  signature, endpoints, error matrix, canonical honest-empty example.
- Slice cards: `../../CHIP-SAT-P2-PLAN.md` §2 (F02-M1..M7, C1, C2, R1, I1).
- Feature brief: `feature.md` (EARS + negative scenarios).

## Scope

Contract-only `market` Go module:

- Migrations `0043_market_observations.sql` + `0044_market_references.sql`
  (0045 reserved-unused, never an empty file; beyond 0045 = hub REQUEST).
- Domain entities with 6-signal separation; `CollectorPort` exactly per IC-04.
- Append-only postgres repositories (latest-per-key reads, input order preserved).
- `GET /market/observations` + `GET /market/references` (max 200 ids; unknown/malformed
  ids → item-level `no_price_evidence`, request stays 200; 201 ids → 422 `too_many_ids`;
  observations require `installation_id` → 400 `installation_required`).
- OpenAPI + sdk-runtime same commit (C1); category-attribute contract reservation (C2,
  `GET /listings/categories/{category_id}/attributes`, no handler).
- Composition-root registration of reads only (R1); e2e + grep absence proof (I1).

## Hard gates

- NO production CollectorPort implementation, NO seed, NO scheduler, NO UI/frontend files.
- Test double confined to `_test` packages.
- Stored `source` + `captured_at` NOT NULL; synthetic no-evidence rows never persisted.
- Signals never merged; `evidence_state` is the only derived field.
- Every tenant query carries `tenant_id`; INSERT-only writes.

## Slice map (write-DAG)

M1 → M2 → {M3 complex, M4 complex (after M3's shared constructor file), M5} → M6 → M7 →
R1 → I1; M7 → C1 → C2 (same OpenAPI/SDK files, serialized).

Per slice: failing test first → implement → sonnet delta review (REVIEW-STANDARD,
anchor-or-abstain, AI-slop REJECT-on-hit) → commit. Ledger rows in
`../../CHIP-SAT-DISPATCH-LEDGER.md`.
