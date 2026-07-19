# M-09-dashboard-demo — Validation Result

```yaml
id: M-09
type: milestone-validation-result
status: in-progress
branch: claude/adoring-euclid-3ccec3
tip: 4d62353a
base: 89de2fef
code_range: 89de2fef..4d62353a
gate: dual (D08 cold Opus + D09 adversarial sonnet) = PASS/agreement
```

Sections: [kpis](#kpis) · [fila](#fila) · [ausente](#ausente) · [api](#api) · [seams](#seams)

---

## seams

**Changed files (19 = 18 code + 1 evidence), each mapped to ownership/grant:**

| File | Justification |
|---|---|
| `.mnfs/MIS-004-mvp-demo/M-09-dashboard-demo/DISPATCH-LEDGER.md` | evidence path |
| `apps/server_core/internal/modules/dashboard/application/service.go` | owned `modules/dashboard/**` |
| `apps/server_core/internal/modules/dashboard/application/service_test.go` | owned |
| `apps/server_core/internal/modules/dashboard/domain/summary.go` | owned |
| `apps/server_core/internal/modules/dashboard/ports/sources.go` | owned |
| `apps/server_core/internal/composition/root.go` | **grant #3** — wire existing erpQuerySvc into dashboard NewService |
| `contracts/api/marketplace-central.openapi.yaml` | owned `/dashboard*` (additive) |
| `contracts/governance/modules.json` | **grant #2** — add `erp_import` to dashboard deps |
| `packages/sdk-runtime/src/dashboard.ts` | owned `sdk-runtime/src/dashboard.ts` |
| `packages/sdk-runtime/src/dashboard.test.ts` | owned (test sibling) |
| `packages/sdk-runtime/src/index.ts` | **grant #1** — one barrel export line |
| `apps/web/src/pages/dashboard/DashboardPage.tsx` | owned `pages/dashboard/**` (rebuild) |
| `apps/web/src/pages/dashboard/DashboardPage.test.tsx` | owned |
| `apps/web/src/pages/dashboard/FilaDeAtencao.tsx` | owned |
| `apps/web/src/pages/dashboard/FilaDeAtencao.test.tsx` | owned |
| `apps/web/src/pages/dashboard/PedidosRecentes.tsx` | owned |
| `apps/web/src/pages/dashboard/Atalhos.tsx` | owned |
| `apps/web/src/routes/dashboard.tsx` | owned |
| `apps/web/src/app/AppRouter.test.tsx` | **grant #4** — root-route mounts rebuilt dashboard |

- **Migrations:** `git diff --name-only 89de2fef..HEAD | grep -iE 'migration|\.sql$'` → **NONE**. C05 ✓ (no self-assigned migration).
- **Light/dark:** pending live-drive (see §ausente/§kpis screenshots).
- **Verdict C05:** PASS — every write in-bounds or a named grant; zero migration.

## api

- **Additive extension of `GET /dashboard/summary`:** new fields `anuncios_ativos` (nullable), `last_import` (object|null: `{at, age_seconds}`), and enum member `erp_import` added to `degraded[]`. No existing field removed / renamed / retyped.
- **Additivity proven:** sdk-runtime `index.test.ts` (61 old-contract round-trips) stays green under the new `dashboard.ts` types; `dashboard.test.ts` (new) covers the added fields. Full suite 71/71.
- **No cross-schema SQL:** `grep -rniE 'SELECT|FROM|JOIN|pgxpool|db.Query|sql.' modules/dashboard --include=*.go` (excl. tests) → **ZERO**. All cross-module reads via public ports: `OrdersSource`, `ListingsSource`, `LinkageSource`, `SyncSource`, `InstallationSource`, and new `ErpImportSource.ListImports` (→ erp `QueryService`).
- **Live before/after responses:** pending dev-stack (hub REQUEST) — will capture `GET /dashboard/summary` payload against a real fixture.
- **Verdict C04:** PASS (static) — additive + port-only reads confirmed; live payload capture pending.

## kpis

_Pending P7 live-drive (awaiting hub dev-stack). Plan: open `/`, record each KPI, open the owning screen, compare same instant. Vendas parity = API-level per hub ruling (/pedidos renders buckets, not a vendas card); anuncios/exceção vs /anuncios summary; sem-vínculo vs /vinculos pendentes; último import vs M-01 protocolo. Paired evidence (card ↔ source)._

## fila

_Pending P7 live-drive. Plan: click each Fila item, assert it lands on the owning screen with the exception filter pre-applied (sync_error → /anuncios?filter.exception=sync_error; below_margin → …=below_margin; sem-vínculo → /vinculos pendentes). Navigation transcript._

## ausente

_Pending P7 live-drive. Plan: drive a source-empty state (e.g. no ERP import) and confirm honest-absent render (EmptyState/ErrorState + reason) distinct from a real 0 — never a fabricated zero (ADR-17). Light + dark screenshots (note: rasterizer F-ENV-10 broken → visual evidence via computed-style + a11y tree if screenshots unavailable)._
