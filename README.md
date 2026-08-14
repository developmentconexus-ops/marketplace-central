# Marketplace Central

Marketplace Central is an internal marketplace operations and intelligence system, initially Mercado Livre, backed by real Sankhya/Oracle operational facts.

## Current status

**Architecture Rebaseline / Technical System Design is in progress. Product implementation is intentionally paused until the D0–D9 design program is accepted.**

Start here:

1. [`AGENTS.md`](AGENTS.md)
2. [`docs/engineering/rebaseline/README.md`](docs/engineering/rebaseline/README.md)
3. [`docs/README.md`](docs/README.md)

Do not use historical Git plans/handoffs as the current roadmap.

## Applications

- `apps/server_core` — canonical Go backend/current runtime
- `apps/web` — React client

The detailed target context, API, database, frontend and runtime topology is **not yet frozen**; the current repository contains both legacy modules and emerging contexts and is being mapped before implementation continues.

## Development verification

```powershell
npm run gate
npm run gate:full
```

The same gate implementation is used by CI. Operational/local commands in the repository describe the current runtime only; target runtime/deployment is a D7 decision.

## External systems

- Mercado Livre — first marketplace control plane
- Sankhya/Oracle — internal operational source behind MPC-owned adapters
- PostgreSQL — MPC-owned application state

See `ARCHITECTURE.md` and the ADR registry for stable constraints; see the rebaseline README for current design status.