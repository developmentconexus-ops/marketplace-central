# Handoff — Issue #1 (I1-edge)

Target session: fresh. Read this, then `GATE-TOPOLOGY.md` §1 and issue
[#1](https://github.com/developmentconexus-ops/marketplace-central/issues/1). Nothing else required.

---

## Task

`/orders` returns buyer PII. No identity check anywhere in the chain. Close it at the edge.

## The five facts

| Fact | Where |
|---|---|
| Handler serialises buyer name, CPF/CNPJ, address | `orders/transport/http_handler.go:608-618` |
| Chain is `CORSMiddleware(apierror.Recover(mux))` — no identity | `composition/root.go:994` |
| Caddy routes `/orders` + non-HTML `Accept` → `backend:8080` | `deploy/Caddyfile`, matcher `@orders_api` |
| ngrok publishes the frontend, profile `oauth` | `docker/dev/ngrok-entrypoint.sh` → `ngrok http --url="$callback_host" frontend:5174` |
| Proxy table, third copy of the routes | `apps/web/vite.config.ts` |

## Deliverable

1. Deny rule on the `@orders_api` matcher in `deploy/Caddyfile`.
2. ngrok scope narrowed so it does not publish the API surface.
3. **The fixture** — see below. Not optional.

## The fixture, and why it is the actual deliverable

A request through the composed stack asserting **non-200 without credentials**.

`#1` was written as a stopgap that `GATE-TOPOLOGY.md` L2-c (boot assertion: no PII route composed
without identity middleware) would retire. **Authentication is deferred by operator decision.** L2-c
is not coming. So this edge rule is the only control on the route, indefinitely — a Caddy config and
an ngrok flag, with no type and no boot condition behind it.

Without a fixture, "the door is closed" is a claim about a config file nobody re-reads.

**Sequencing wrinkle, handle it explicitly:** `verify-full` does not exist yet — it ships in #2. So
the fixture lands now as a runnable test with a named command, and #2 wires it. **Write down where it
lives and what invokes it**, or #2 will not find it and the census check in #2 will be the thing that
catches you.

## Out of scope — hard

- **Authentication.** Principals, sessions, tokens. #1 closes a door; it does not decide who has keys.
- **`composition/root.go`**, beyond what the deny rule strictly needs. That is #5.
- **A `127.0.0.1` bind.** Arm A proposed it, then **withdrew it itself** after determining it would
  read as done and be inert. Do not re-propose.

## Acceptance

Level 3. A check that fetches the route through the composed stack and asserts non-200 without
credentials. Not higher — a Caddy config is not representable in the Go type system.

## Unknown that changes urgency, not scope

**Is a production host currently running?** `MPC_DOMAIN` is external and neither synthesis arm could
settle it. If yes, this is immediate rather than merely urgent. Ask before assuming either way.

## Constraints inherited from the audit

- **No push.** Requires explicit operator permission at the time. Local commits only.
- **No `git reset` / `revert` / `stash` / `clean` / history rewrite.**
- **Never dump an environment.** No `docker inspect`, no bare `printenv`, no `. .env`. Diagnose one
  variable by name. Never print a secret value.
- **Oracle/Sankhya is read-only.** Never write to the ERP.
- **Dependency change = operator ACK.** Only `oapi-codegen` and `openapi-typescript` are approved.
- **`gh` active account is global across sessions.** Two same-named private repos exist. Any `gh`
  command that writes: switch, verify identity, mutate — one shell invocation, hard abort between
  verify and mutate, and always pin `-R developmentconexus-ops/marketplace-central`.

## Blocks

**P6, the public flip.** The flip publishes `deploy/Caddyfile`, `apps/web/vite.config.ts` and
`orders/transport/http_handler.go:608-618` — the route, its predicate, and the fields it returns.
The exposure is already live; the flip removes the only thing currently limiting who knows the route.

## Size

Hours.
