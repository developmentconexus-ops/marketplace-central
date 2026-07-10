# Docker Dev Environment

This setup runs Marketplace Central with Docker Compose:

- `backend`: Go server on `http://localhost:8080`
- `frontend`: Vite app on `http://localhost:5174`
- `postgres`: MPC-owned dev database on host port `5435`
- `ngrok`: optional OAuth callback tunnel using the Mercado Livre URL already registered in the app

The Compose file reads the local `.env`, but it does not bake secrets into images.

## Start Daily Dev

```powershell
docker compose up --build postgres backend frontend
```

Open:

```text
http://localhost:5174
```

## Start OAuth Dev With Ngrok

Set your ngrok token in the shell or local `.env`:

```powershell
$env:NGROK_AUTHTOKEN = "your-ngrok-token"
```

Then run:

```powershell
docker compose --profile oauth up --build
```

The ngrok service uses the registered Mercado Livre callback:

```text
https://multiradial-unironically-nieves.ngrok-free.dev/integrations/auth/callback
```

The callback URL in `.env` is the source of truth; the container derives the reserved ngrok domain from `MPC_OAUTH_REDIRECT_URI`.

Ngrok's local inspector is exposed at:

```text
http://localhost:4040
```

## Validate Backend Packages

```powershell
docker compose run --rm backend bash ./docker/dev/backend-entrypoint.sh test ./apps/server_core/internal/modules/internal_read/...
```

## Notes

- Compose overrides `MC_DATABASE_URL` inside the backend container to use the `postgres` service.
- The current server still opens `MS_DATABASE_URL` during boot for legacy catalog composition. The container maps it to the dev Postgres URL only so the current binary starts; this is not the Oracle-first internal-read target path.
- Oracle live validation inside Linux containers uses the Instant Client installed in the backend image at `/opt/oracle/instantclient`.
