# M-09 Dashboard — Hub Close Evidence Pack

```yaml
id: M-09-dashboard-demo
type: hub-close-evidence
author: Dispatch Hub (Wave C)
mission: MIS-004-mvp-demo
created: 2026-07-19
```

- **Branch:** chip/adoring-euclid-3ccec3
- **Base:** 783cbc0d55ad91fe9cb73a97e8329d589587a0be (hub tip)
- **Detailed chip result:** `.mnfs/MIS-004-mvp-demo/M-09-dashboard-demo/validation-result.md` (brought in by this merge)

## Merge-gate markers (harness 0.4.0)

P6-DUAL-GATE: AGREEMENT

Dual gate ran during the chip (DISPATCH-LEDGER D08 cold Opus reviewer + D09 adversarial sonnet refuter) → agreement, recorded in the chip's validation-result.md §gate. Hub verified the marker + diff scope (20 files: dashboard module + contracts + sdk-runtime, all parity-aligned) at close.

LIVE-VERIFIED: 2026-07-19 hub P7 live-drive on the clean docker dev-stack (backend :8080 bound to the M-09 worktree mount, frontend :5174, postgres fixture 30 orders intact). Journey `/` dashboard rendered honestly:
- `anuncios_ativos = 10` (real, mapped from listings.Active in service.go:75)
- `below_margin = null` → rendered `—` (honest unknown, ADR-17)
- `orders_today / orders_7d = 0` (fixture predates the window — honest, not fabricated)
- `last_import = 18/07/2026 17:23 · "há 1 d"` (real ERP import timestamp)
- `pending_links = 1`, `missing_gtin = 3`; "Pedidos recentes" = 5 real fixture orders; "Fila de atenção" functional
- Zero console errors; theme paper (`rgb(251,250,247)`) + green in both light & dark; Instrument Sans + IBM Plex Mono numerals

NOTE: a transient `undefined` on `anuncios_ativos` seen mid-drive was traced to a **hub split-mount** (the backend container had been restarted from the main repo root instead of the worktree, so M-09's FE called main's field-less backend). Corrected via `docker compose up --force-recreate --no-deps backend` from the worktree cwd. **Not an M-09 code defect** — the field is present and parity-aligned across `domain/summary.go:24`, `application/service.go:75`, `contracts/api/*.openapi.yaml`, and `packages/sdk-runtime/src/dashboard.ts`.

## Pre-merge lane check (L0–L2)

- Go dashboard module (`GOCACHE` abs, scratchpad): `go test ./internal/modules/dashboard/...` → **ok** (application 1.531s, transport 2.361s), exit 0; no gomodcache/gocache pollution.
- Web vitest (dashboard + AppRouter) and sdk-runtime parity: green at chip close (DISPATCH-LEDGER D04/D05: "71 sdk tests pass", build+test green). Re-confirmed on main post-merge as the C08 integration gate (main has node_modules + warm caches; chip worktree does not — vitest-node_modules junction trap).

## Verdict

**M-09 Dashboard: PASS.** P6 dual-gate AGREEMENT + P7 live-drive PASS + Go lane green. Cleared for `--no-ff` merge into main.
