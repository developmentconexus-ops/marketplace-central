# Marketplace Central

Marketplace Central is an internal **Marketplace Operations Control Plane + Commercial Intelligence** product, initially proving its operating loop with Mercado Livre and Sankhya.

## Start here

Read only:

1. [`AGENTS.md`](AGENTS.md)
2. [`docs/README.md`](docs/README.md)

`docs/README.md` is the sole authority for current stage, allowed/blocked work and exact next action. Do not use historical plans, handoffs, review dialogue, Git branches or current code shape as the roadmap.

## Current posture

Architecture Rebaseline D0→D9 is in progress. Product implementation remains blocked until D9 is accepted. Detailed target context, API, persistence, frontend and runtime topology is governed by the accepted D-stage authority indexed in `docs/README.md`.

## Applications

- `apps/server_core` — canonical Go backend/current runtime
- `apps/web` — React client

The current repository contains legacy runtime structures as current-state evidence. Existing code does not become target architecture by inheritance.

## External systems

- Mercado Livre — first marketplace proof
- Sankhya — external business system; target integration uses its sanctioned API Gateway
- Oracle/godror — legacy/current-state evidence, not target fallback
- PostgreSQL — MPC-owned application state

## Verification

```powershell
npm run gate
npm run gate:full
```

The same gate implementation runs locally and in CI. Operational commands describe the current runtime only; target runtime/deployment is a D7 decision.
