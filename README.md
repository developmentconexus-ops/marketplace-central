# Marketplace Central

Marketplace Central is a MetalShopping-style monorepo for pricing simulation and marketplace configuration.

## Apps

- `apps/server_core`: canonical Go backend
- `apps/web`: thin React client

## Packages

- `packages/sdk-runtime`: runtime client for the web app
- `packages/ui`: shared UI primitives
- `packages/web-query`: shared web data-query layer
- `packages/feature-*`: per-workspace React screens — `feature-classifications`, `feature-connectors`, `feature-inventory`, `feature-orders`, `feature-products`, `feature-simulator`

## Docker Dev Environment

The repo includes a Docker Compose dev setup for backend, frontend, Postgres, and optional ngrok OAuth callback tunneling.

Daily local dev:

```powershell
npm run docker:dev
```

OAuth dev with the Mercado Livre callback URL already registered in the app:

```powershell
$env:NGROK_AUTHTOKEN = "your-ngrok-token"
npm run docker:oauth
```

Open the web app at `http://localhost:5174`. The backend is exposed at `http://localhost:8080`, and the optional ngrok inspector is exposed at `http://localhost:4040`.

Details live in [docker/dev/README.md](docker/dev/README.md).
