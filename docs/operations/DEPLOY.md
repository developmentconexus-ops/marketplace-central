# Production Deploy Runbook

Binding architecture: [ADR-008](../architecture/decisions/008-production-deploy-topology.md).
Artifacts: `docker/prod/` (images), `deploy/` (compose + Caddyfile + env template),
`scripts/release.sh` (publisher).

Flow in one line:

```
build (scripts/release.sh) → GHCR → host: docker compose pull && up -d
```

One publisher, tag scheme `sha-<commit>`:

- **`bash scripts/release.sh`** — builds both images on the dev machine and
  pushes to GHCR. Prereq once: `docker login ghcr.io` with a `write:packages`
  PAT.

There is no CI publisher. `.github/workflows/release-images.yml` existed until
2026-08-10 and failed all 18 of its runs with `403 Forbidden` on the GHCR push,
never publishing an image; it was deleted rather than left red. See the
[ADR-008 amendment](../architecture/decisions/008-production-deploy-topology.md#amendment-2026-08-10--the-ci-publisher-is-retired-scriptsreleasesh-is-the-publisher)
for the credential precondition that must be settled before a CI publisher
returns.

---

## 1. One-time: provision the host (VPS variant)

Target: Ubuntu Server 24.04 LTS, 2 vCPU / 4 GB (Magalu Cloud BV2-4-40 or Vultr São Paulo).

```bash
# as root on a fresh host
adduser mpc && usermod -aG sudo mpc            # named sudo user, no root logins
curl -fsSL https://get.docker.com | sh          # Docker Engine + compose plugin
usermod -aG docker mpc
curl -fsSL https://tailscale.com/install.sh | sh
tailscale up                                    # join the tailnet
```

Harden SSH (`/etc/ssh/sshd_config`), then `systemctl restart ssh`:

```
PasswordAuthentication no
PermitRootLogin no
ListenAddress <tailscale-100.x.x.x-ip>   # SSH reachable only inside the tailnet
```

Firewall: allow 80/tcp and 443/tcp+udp only (`ufw allow 80,443/tcp && ufw allow 443/udp && ufw enable`).
Port 22 stays closed publicly — operator access is via Tailscale.

DNS: point `MPC_DOMAIN` (A record) at the VPS IP. Caddy provisions TLS automatically.

### Oracle path (Sankhya)

On any always-on machine inside the client network:
install Tailscale, enable subnet routing for the Oracle host's subnet
(`tailscale up --advertise-routes=<oracle-subnet>/24`), approve the route in
the admin console. `SANKHYA_ORACLE_HOST` in `.env` then uses the Oracle
host's LAN IP, reachable from the VPS through the tailnet. The database is
never exposed to the internet.

## 2. One-time: install the stack

```bash
mkdir -p ~/mpc && cd ~/mpc
# copy from the repo: deploy/docker-compose.prod.yml (as docker-compose.prod.yml)
#                     deploy/Caddyfile
# create .env from deploy/env.production.example — fill every CHANGEME
chmod 600 .env
docker login ghcr.io -u <github-user>    # PAT with read:packages only
docker compose -f docker-compose.prod.yml pull
docker compose -f docker-compose.prod.yml up -d
curl -fsS https://<MPC_DOMAIN>/healthz
```

Register the production OAuth callback in the Mercado Livre app:
`https://<MPC_DOMAIN>/integrations/auth/callback` (replaces the ngrok URL for prod).

### The `/orders` edge credential — read this before first boot

`/orders` returns buyer name, CPF/CNPJ and billing address, and **the
application does not check who is asking** — there is no identity middleware in
the composed chain (`internal/composition/root.go:994`). The HTTP basic
credential Caddy enforces on that route is the only control there is.

It is not a placeholder. The boot-time assertion intended to replace it
(`docs/engineering/repo-audit-2026-08-07/GATE-TOPOLOGY.md`, L2-c) is deferred
along with authentication, so this credential is load-bearing indefinitely.

Generate the hash on the host and put both values in `.env`:

```bash
docker run --rm caddy:2-alpine caddy hash-password --plaintext '<a long random password>'
```

`MPC_ORDERS_BASIC_AUTH_USER` and `MPC_ORDERS_BASIC_AUTH_HASH` are both declared
`:?` in `docker-compose.prod.yml` — **the stack refuses to start if either is
missing**, rather than starting with an open door.

What you will see in the browser: opening `/orders` loads the page normally (a
hard nav sends `Accept: text/html` and gets the SPA), then the first data
request prompts once for the credential. The browser caches it for the origin.
`/healthz` and the rest of the API are unaffected — the rule is scoped to the
`/orders` JSON surface.

Verify the rule after any change to `deploy/Caddyfile`:

```bash
npm run harness:edge
```

## 3. Recurring: deploy an update

`bash scripts/release.sh` publishes a `sha-<commit>` tag from the dev machine.
To ship one:

```bash
ssh mpc@<tailscale-ip>                 # via tailnet
cd ~/mpc
# 1. point .env MPC_IMAGE_TAG at the new sha-<commit>
# 2. migrations run automatically on backend start (RUN_MIGRATIONS=1)
docker compose -f docker-compose.prod.yml pull
docker compose -f docker-compose.prod.yml up -d --remove-orphans
# 3. smoke check
curl -fsS https://<MPC_DOMAIN>/healthz
docker compose -f docker-compose.prod.yml logs --tail=50 backend
```

**Rollback:** set `MPC_IMAGE_TAG` back to the previous sha tag, repeat pull/up.
Seconds, no rebuild. (Caveat: only safe across migrations that are
backward-compatible — prefer additive migrations.)

## 4. Backups

Daily `pg_dump` + offsite copy. On the host (`crontab -e` for the mpc user):

```cron
0 3 * * * docker compose -f /home/mpc/mpc/docker-compose.prod.yml exec -T postgres pg_dump -U marketplace marketplace_central | gzip > /home/mpc/backups/mpc-$(date +\%F).sql.gz
30 3 * * * rclone copy /home/mpc/backups remote:mpc-backups --max-age 48h
0 4 * * 0 find /home/mpc/backups -name '*.sql.gz' -mtime +30 -delete
```

`rclone` remote: Backblaze B2 or any S3-compatible bucket. **Test a restore
once** (`gunzip -c dump.sql.gz | docker compose exec -T postgres psql -U
marketplace marketplace_central` into a scratch database) before trusting it.

## 5. On-prem variant (client-hosted)

Same images, same compose, same runbook — differences only:

- Host is a client machine: Linux + Docker Engine required. **Never Docker
  Desktop on a server** (unsupported on Windows Server; paid license for
  companies >250 employees / >US$10M revenue). Windows-only shop → provision a
  Linux VM (Hyper-V) dedicated to the stack.
- No public DNS needed if access is LAN-only: Caddy can serve the LAN
  hostname with internal TLS (`tls internal` in the Caddyfile site block).
- Oracle is on the same LAN — no subnet router needed; Tailscale stays for
  operator SSH.

## 6. Standing rules

- Production hosts run **images only** — never source, never `--build`.
- `.env` on hosts: `chmod 600`, never in git, never in chat/screenshots.
- Pin `MPC_IMAGE_TAG` to sha tags; `latest` is for dev convenience only.
- New backend route prefix ⇒ update BOTH `apps/web/vite.config.ts` proxy table
  and `deploy/Caddyfile` (same PR).
- Deploys are deliberate and attended; no auto-update agents.
