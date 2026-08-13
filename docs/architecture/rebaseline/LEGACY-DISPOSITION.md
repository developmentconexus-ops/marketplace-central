# Marketplace Central — Legacy Decision Disposition

> **Status:** CANONICAL D0 ledger  
> **Historical baseline:** `main@de1dc88bcef5a6ed5515378e7c646682c0bc15d2`

This ledger is the bridge between historical architecture material and D1–D9. Historical source files remain retrievable from Git at the baseline above; they are not active target authority.

## Vocabulary

- **CARRY** — the underlying property remains a constraint unless a material finding changes it.
- **REOPEN D#** — evidence survives, but the target boundary/mechanism is decided in the named gate.
- **RETIRED** — a later direction already replaced the decision.
- **HISTORICAL** — explains prior sequencing/evidence only.

## ADR disposition

| ADR | Historical subject | Disposition | Surviving meaning |
|---|---|---|---|
| 001 | direct MetalShopping Postgres reads | RETIRED | external/internal-source access is rederived in D4 |
| 002 | MPC schema on MetalShopping cluster | RETIRED | physical persistence target is D2/D7 |
| 003 | integration spec sequencing | HISTORICAL | D0–D9 owns sequencing |
| 004 | integration catalog/plugin framework | REOPEN D1/D4 | capability-specific consumer-owned ports are the current candidate |
| 005 | Mercado Livre first | CARRY | first operational marketplace remains Mercado Livre |
| 006 | MPC-owned Oracle read boundary | CARRY | business code depends on MPC-owned ports; exact context owner is D1/D4 |
| 007 | godror + OCI runtime | CARRY | current supported Oracle runtime unless D4 finds material reason to change |
| 008 | production deployment topology | REOPEN D7 | no production compatibility forces the old shape |
| 009 | fee provenance | CARRY | external economic facts remain attributable; exact model D2/D4 |
| 010 | polling/visible refresh | REOPEN D4/D7 | freshness mechanism is capability-specific |
| 011 | divergence row model | REOPEN D1/D2 | reconciliation ownership/data model is not presumed |
| 012 | DIFAL owned by pricing | REOPEN D1/D2 | one-authority property survives; owner does not |
| 013 | webhook payload is not automatically domain truth | CARRY | D4 defines authoritative retrieval per topic |
| 014 | market collection on demand | REOPEN D1/D4 | intelligence/market lifecycle follows product need |
| 015 | listings module read-only | REOPEN D1 | observed vs desired/write responsibilities are redesigned |
| 016 | hand-written SDK + same-commit OpenAPI | REOPEN D5 | current guard protects current state; D5 chooses one target contract authority |
| 017 | unknown is never zero | RETIRED by 034 historically; PROPERTY CARRIES | D2 preserves honest absence without assuming a universal wrapper |
| 018 | mutation envelope + in-process poller | REOPEN D1/D3/D7 | safe external change lifecycle remains required; mechanism is open |
| 019 | listings snapshot observer | REOPEN D3/D4 | observation/freshness/event semantics are rederived |
| 020 | market data through CollectorPort | REOPEN D1/D4 | provider access behind consumer-owned ports survives; interface/name does not |
| 021 | frontend/TanStack server-state seam | CARRY property / REOPEN D6 | thin frontend + one server-state cache authority; package topology open |
| 022 | SELLER_SKU equals CODPROD before writes | REOPEN D2/D4 | preserve as historical safety evidence, not universal identity law |
| 023 | legacy module protocol | REOPEN D1 | compiler-backed private boundaries are strong evidence; `internal/modules` is not target authority |
| 024 | one order-ingest writer | CARRY property / REOPEN D1–D3 | one write authority per state survives; exact API does not |
| 025 | selective raw payload, no raw PII | CARRY | D2/D4 define precise retention |
| 026 | scheduler phase vocabulary | REOPEN D7 | platform mechanics vs context freshness/cursor semantics are redesigned |
| 027 | absence from partial pull is not closure | CARRY | incomplete observation cannot fabricate terminal state |
| 028 | automatic link on concordant anchors | REOPEN D1/D2/D4 | linking requires evidence; exact anchors/confidence need current validation |
| 029 | no blind retry on writes | CARRY | D4/D7 model idempotency and outcome-unknown explicitly |
| 030 | second scheduler per installation | REOPEN D7 | old scheduler topology is not target authority |
| 031 | keep-absent merge | REOPEN D2/D4 | only valid if source completeness semantics prove it |
| 032 | catalog-offers feature flag | HISTORICAL | past rollout mechanism only |
| 033 | vendor adapters implement consumer-owned ports | CARRY | target dependency direction |
| 034 | kernel Fact replaces ADR-017 | REOPEN D2 | honest-unknown property carries; exact `Fact[T]`/kernel scope must earn its place |

## Historical document families

| Family | Disposition | What survives |
|---|---|---|
| old root implementation roadmap | HISTORICAL | no active sequencing; D0–D9 replaces it |
| legacy wiki | RETIRED | product north star is absorbed into the current constitution |
| MNFS mission/chip/debt trees | HISTORICAL | Git history only |
| old Superpowers specs/plans/handoffs | HISTORICAL | useful measured facts enter the evidence register; old sequences do not |
| 2026-08-07 repository audit | HISTORICAL EVIDENCE | its lessons inform the current method/gate; current claims require fresh measurement |
| old design/domain/identity proposals | REOPEN D1/D2/D4 | domain/source evidence survives; unratified structural conclusions do not |
| dated marketplace API references/research | HISTORICAL EVIDENCE | D4 re-verifies material provider behavior from current official sources |
| old deploy/runtime runbooks | REOPEN D7 | active scripts are current-state evidence; target operational topology waits for D7 |
| old harness/onda review doctrine | RETIRED | current AGENTS + canonical engineering method + D8 golden-flow proof replace it |
| retired agent workflow/plugin/planning skills | RETIRED | no authority in the D0–D9 program |

## Current mechanisms deliberately not re-adjudicated in D0

`contracts/governance/*`, `scripts/gate.ps1`, `scripts/harness.ps1`, the current OpenAPI/SDK, migrations, application code and tests remain **current-state mechanisms/evidence** until the D-stage that owns each surface changes them.

D0 removes conflicting doctrine; it does not silently redesign runtime controls.

## Future-document rule

Every new document must be visibly one of:

- **CANONICAL** — owns a current decision/status;
- **SUPPORTING EVIDENCE** — can inform but cannot override canonical decisions;
- **IMPLEMENTATION PLAN** — only after D9 and scoped to accepted target design.

Do not keep superseded plans in the active tree merely “for context”. Git provides context without making history look executable.