# Marketplace Central

Marketplace Central is an internal **Marketplace Operations Control Plane + Commercial Intelligence** product, initially proving its operating loop with Mercado Livre and Sankhya.

> **Landing page only.** This file carries no independent program status, next action or architecture authority.

## Start here

Read only:

1. [`AGENTS.md`](AGENTS.md)
2. [`docs/README.md`](docs/README.md)

`docs/README.md` is the sole current-program status and selective documentation index. Historical plans, handoffs, review dialogue, branches and current code shape are not the roadmap.

The superseded application, deployment, migration and manually maintained API/SDK trees were removed during the architecture rebaseline. Git history preserves them as evidence; they are not an active compatibility surface.

## Verification

```powershell
npm run gate
npm run gate:full
```

The same compact authority/zero-legacy gate runs locally and in CI while implementation remains blocked through D9.
