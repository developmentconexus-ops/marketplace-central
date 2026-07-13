# M-13 Validation Contract

```yaml
id: M-13
type: milestone-validation-contract
status: ready
owner: Mission Strategist
parent: MIS-001
created: 2026-07-13
updated: 2026-07-13
validation_level: QA-3
lifecycle_scope: milestone
```

## Required Outcome

The browser provides one coherent operator cockpit with Visão geral, Produtos /
Product 360, Anúncios, Vendas/Margem, and Operações, connected by stable deep links,
shared installation context, honest quality states, and simulation-only stock review.

QA records the frozen SHA once. Review is fixed-SHA and read-only; QA then runs the
proportional tests and browser checks below.

## Criteria

### M-13-C01 — Navigation and attention closure

- Required: Yes.
- Proof: router/layout tests plus browser interaction.
- Expected: five Portuguese workspaces; attention opens the exact filtered object;
  legacy routes redirect without losing installation/entity context.
- Blocks: duplicate top-level journeys or attention that lands on an unfiltered list.

### M-13-C02 — Product/listing/sale continuity

- Required: Yes.
- Proof: SDK/component tests plus browser deep-link and reload checks.
- Expected: Product 360, listing detail, and sale detail cross-link without mixing
  CODPROD with provider identifiers; reload retains context.
- Blocks: manual identity reconstruction or client-side stock/margin calculation.

### M-13-C03 — Operations convergence

- Required: Yes.
- Proof: component tests and browser check.
- Expected: connection/auth state, probes, runs, source, observed time, and actionable
  errors are available in Operações without duplicate technical navigation.
- Blocks: operational failure is hidden or shown as healthy/default zero.

### M-13-C04 — Safe stock simulation

- Required: Yes.
- Proof: inventory/UI tests and browser inspection of reachable controls.
- Expected: `Simulação`, current/proposed values, reason, source timestamps, preview
  payload, and `executed=false`; no execute/provider-write control is reachable.
- Blocks: copy or behavior implies a real provider update.

### M-13-C05 — Honest operational states

- Required: Yes.
- Proof: deterministic UI tests for applicable loading, error, empty, current,
  stale/unknown, and blocked/conflict states; desktop and 390x844 browser pass.
- Expected: Portuguese operator copy, visible reasons, keyboard-usable primary
  navigation, and no horizontal page overflow on the bounded mobile viewport.
- Blocks: unknown becomes zero or a load-bearing error/conflict is invisible.

### M-13-C06 — SDK-only thin client

- Required: Yes.
- Proof: targeted source scan plus web/SDK tests and build.
- Expected: production web features use `packages/sdk-runtime`; API changes update
  OpenAPI and SDK together; React does not call providers or Oracle.
- Blocks: direct backend/provider fetch or domain calculation in the client.

## Browser Check

1. Confirm the five workspaces and selected installation.
2. Open an attention item to its exact listing.
3. Navigate Listing → Product → Listing → Sale → Listing.
4. Review the simulation and confirm no execution control/result exists.
5. Reload each object route and verify context persists.
6. Exercise representative loading, error/empty, and unknown/conflict states.
7. Repeat the primary journey at 390x844.

Deterministic fixtures may prove state rendering. They must be labeled as fixtures
and are not evidence of a real provider/Oracle integration.

## Evidence Requirements

- Feature `validation.md` files with commands and outcomes.
- Fixed-SHA review findings.
- QA `validation-result.md`, concise command outputs, browser interaction note, and
  desktop/mobile screenshots. A per-state screenshot matrix is optional.

## Retry Policy

Maximum two scoped correction attempts. One writer owns router/layout/shared context
at a time; API/OpenAPI/SDK is one serialized seam.

## Handoff

- Current status: Ready after M-09 passes.
- Next owner: M-13 Milestone Orchestrator, then `mpc-verifier`.
- Open decisions: none; stop if an aggregate would require client-side business logic.
