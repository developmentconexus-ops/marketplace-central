# ADR-008: Image-based delivery with a VPS-primary production topology

**Date:** 2026-07-19
**Status:** accepted

## Context

Until now every runtime of MPC has been a developer machine: the Docker Compose
stack bind-mounts source and runs `go run`/Vite dev servers, and the Mercado
Livre OAuth callback rides an ngrok tunnel. A client demo and a real pilot
require a production story: a stable HTTPS address, remote updates without
touching the host by hand, secrets discipline, backups, and a defined path to
the client's Sankhya Oracle database (which must never be exposed to the
internet).

Options evaluated (research 2026-07-19, sources in the deploy runbook):

- **Serverless/PaaS (Vercel, Railway, Render, Fly.io).** Vercel cannot host the
  long-running Go server, a Postgres instance, or a persistent VPN tunnel to
  the on-prem Oracle; it only fits the SPA in a hybrid split that adds an infra
  seam without MVP benefit. Container PaaS options force refactoring the
  Compose stack into their per-service model and cost more than a VPS.
- **Build-on-server.** Shipping source to the target and running
  `compose up --build` requires the full toolchain on the host, produces
  non-reproducible builds, and leaks source to client machines.
- **Auto-updaters (Watchtower).** Watchtower was archived in December 2025;
  unattended silent updates are the wrong risk profile for a system carrying
  fiscal data regardless.

## Decision

1. **Delivery is image-based.** CI (GitHub Actions,
   `.github/workflows/release-images.yml`) builds two production images from
   `docker/prod/` multi-stage Dockerfiles — compiled Go binaries plus Oracle
   Instant Client for the server, static Vite build behind nginx for the web —
   and pushes them to GHCR tagged `sha-<commit>` (plus `latest` on main and
   semver on `v*` tags). Targets only ever `docker compose pull && up -d`;
   production hosts never see source or toolchains. Rollback is re-pointing
   `MPC_IMAGE_TAG` at the previous sha tag.
2. **Primary production topology is a Brazil-region VPS** running
   `deploy/docker-compose.prod.yml`: Caddy is the only service publishing host
   ports and terminates TLS for `MPC_DOMAIN`, routing API prefixes (mirroring
   the dev proxy table in `apps/web/vite.config.ts`, including the
   Accept-header split on `/catalog` and `/orders`) to the backend and
   everything else to the SPA. Postgres and the app services stay on the
   internal network.
3. **The Sankhya Oracle link is a private overlay network** (Tailscale): a
   subnet-router node inside the client's network exposes the Oracle host to
   the VPS; the database is never internet-reachable. Operator SSH to
   production hosts also rides the tailnet, not a public port 22.
4. **On-prem is a supported variant, not a fork.** The same images and compose
   file run on a client-hosted Linux + Docker Engine machine (never Docker
   Desktop on servers); remote operations use the same tailnet path. The choice
   between VPS and on-prem is commercial (data residency), not technical.
5. **Deploys are deliberate.** A human runs the update sequence (migrate, pull,
   up, smoke-check) from the runbook (`docs/operations/DEPLOY.md`); no
   unattended auto-update agent is installed.
6. **Production replaces ngrok.** The Mercado Livre OAuth callback registers
   `https://<MPC_DOMAIN>/integrations/auth/callback`; ngrok remains a dev-only
   convenience.

## Consequences

- The dev compose file stays as-is for development; production artifacts live
  in `docker/prod/` and `deploy/` and evolve under the deploy seam (hub-owned).
- A new backend route prefix now has two registration points: the Vite proxy
  table and `deploy/Caddyfile`. The Caddyfile header comment pins this.
- GHCR private-image gratuity is current policy, not contract; if it changes,
  the registry reference in compose/CI is a one-line swap.
- Postgres backup automation (`pg_dump` cron + offsite copy) is part of the
  production definition of done; the runbook carries the procedure.
