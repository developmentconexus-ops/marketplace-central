SLICE CARD — F02-S1 · feature F-02-header-nav-routes · milestone M-03

depends_on: [F01-S2] (theme layer already committed on this branch).

goal: Introduce the six IC-05 area route modules and make AppRouter consume them, WITHOUT changing any URL,
route path, wrapper behavior, redirect, unknown-route behavior, or the surrounding Layout/InstallationProvider.
This is a pure indirection seam: later milestones (M-04..M-09) will replace ONLY the corresponding route module
body. No header/sidebar/theme work here (that is F02-S2).

complexity: standard.

write_set (create the 6 route modules + edit AppRouter):
- apps/web/src/routes/dashboard.tsx      (new — export `DashboardRoute`)
- apps/web/src/routes/anuncios.tsx       (new — export `AnunciosRoute`)
- apps/web/src/routes/vinculos.tsx       (new — export `VinculosRoute`)
- apps/web/src/routes/produto.tsx        (new — export `ProdutoRoute`)
- apps/web/src/routes/precos.tsx         (new — export `PrecosRoute`, OWNS the useClient wrapper)
- apps/web/src/routes/pedidos.tsx        (new — export `PedidosRoute`)
- apps/web/src/app/AppRouter.tsx         (edit — import the 6 from ../routes/*, swap the 6 elements)

EXACT CURRENT AppRouter.tsx (apps/web/src/app/AppRouter.tsx) — the routes you MUST preserve unchanged are
marked KEEP; the six you MUST re-point are marked MOVE:
```
<Route element={<Layout />}>
  <Route index element={<DashboardPage />} />                                    // MOVE → <DashboardRoute/>
  <Route path="/anuncios" element={<AnunciosPage />} />                          // MOVE → <AnunciosRoute/>
  <Route path="/catalogo" element={<CatalogPageWrapper />} />                    // KEEP exactly
  <Route path="/catalogo/produtos/:productId" element={<WorkspacePlaceholder/>}/> // MOVE → <ProdutoRoute/>
  <Route path="/vinculos" element={<WorkspacePlaceholder />} />                  // MOVE → <VinculosRoute/>
  <Route path="/estoque" element={<StockSeguroPageWrapper />} />                 // KEEP exactly
  <Route path="/precos" element={<PricingSimulatorPageWrapper />} />            // MOVE → <PrecosRoute/>
  <Route path="/pedidos" element={<WorkspacePlaceholder />} />                   // MOVE → <PedidosRoute/>
  <Route path="/integracoes" element={<WorkspacePlaceholder />} />              // KEEP exactly
  <Route path="/protocolos/:protocolId" element={<ProtocoloPage />} />          // KEEP exactly
  <Route path="/classifications" element={<ClassificationsPageWrapper />} />    // KEEP exactly
  <Route path="/marketplaces" element={<WorkspacePlaceholder />} />             // KEEP exactly
  <Route path="/products" element={<LegacyRedirect to="/catalogo" />} />        // KEEP all 6 redirects exactly
  <Route path="/product-links" element={<LegacyRedirect to="/vinculos" />} />
  <Route path="/inventory/stock-seguro" element={<LegacyRedirect to="/estoque" />} />
  <Route path="/orders" element={<LegacyRedirect to="/pedidos" />} />
  <Route path="/integrations" element={<LegacyRedirect to="/integracoes" />} />
  <Route path="/simulator" element={<LegacyRedirect to="/precos" />} />
</Route>
```

REQUIRED ROUTE MODULES (each a named page-level component; keep imports minimal & correct relative paths
from apps/web/src/routes/ → app/ is ../app, pages/ is ../pages, packages are bare specifiers):
- dashboard.tsx: `export function DashboardRoute() { return <DashboardPage />; }` (import DashboardPage from
  "../pages/DashboardPage").
- anuncios.tsx: `export function AnunciosRoute() { return <AnunciosPage />; }` (from "../pages/AnunciosPage").
- vinculos.tsx: `export function VinculosRoute() { return <WorkspacePlaceholder />; }` (from
  "../pages/WorkspacePlaceholder").
- produto.tsx: `export function ProdutoRoute() { return <WorkspacePlaceholder />; }` (same placeholder — the
  productId param screen is M-06's future work).
- precos.tsx: OWNS the client wrapper (moved verbatim from AppRouter's PricingSimulatorPageWrapper):
  ```
  import { PricingSimulatorPage } from "@marketplace-central/feature-simulator";
  import { useClient } from "../app/ClientContext";
  export function PrecosRoute() {
    const client = useClient();
    return <PricingSimulatorPage client={client} />;
  }
  ```
- pedidos.tsx: `export function PedidosRoute() { return <WorkspacePlaceholder />; }` (from
  "../pages/WorkspacePlaceholder").

AppRouter.tsx edits:
- Add `import { DashboardRoute } from "../routes/dashboard";` (and the other five) — path is ../routes/<area>.
- Replace ONLY the six MOVE elements with `<DashboardRoute/>`, `<AnunciosRoute/>`, `<ProdutoRoute/>`,
  `<VinculosRoute/>`, `<PrecosRoute/>`, `<PedidosRoute/>` at their SAME paths (index, /anuncios,
  /catalogo/produtos/:productId, /vinculos, /precos, /pedidos). Do NOT change any path string.
- Remove the now-unused `PricingSimulatorPageWrapper` function and the now-unused imports it/those elements
  needed EXCLUSIVELY: `DashboardPage`, `AnunciosPage`, `PricingSimulatorPage`. KEEP `WorkspacePlaceholder`
  (still used inline by /integracoes and /marketplaces), KEEP `useClient` (still used by CatalogPageWrapper,
  ClassificationsPageWrapper, StockSeguroPageWrapper), KEEP all other imports/wrappers/redirects/Layout/
  InstallationProvider untouched.
- Do NOT add a `path="*"` catch-all. Do NOT alter LegacyRedirect rows. Unknown-route behavior stays as-is.

HARD CONSTRAINTS: do NOT touch apps/web/src/pages/** (forbidden), do NOT touch Layout.tsx or Header (F02-S2),
do NOT change ClientContext/InstallationContext. No new deps.

validation_kind: typecheck / grep-assertion / build (unit test N/A — pure wiring; do not fabricate a test).

commands (run from worktree root; capture output — SANDBOX NOTE in bindings applies to build):
- `npx tsc --noEmit -p apps/web/tsconfig.json`   (expect ONLY the known pre-existing TS2688-'node' env error,
   nothing else; any other TS error is yours)
- `npm run build -w @marketplace-central/web`   (chip re-runs if sandbox blocks)
- `rg -n 'from "\.\./routes/(dashboard|anuncios|vinculos|produto|precos|pedidos)"' apps/web/src/app/AppRouter.tsx`  (6 imports)
- `rg -n 'path="/anuncios"|path="/catalogo/produtos/:productId"|path="/vinculos"|path="/precos"|path="/pedidos"|<Route index' apps/web/src/app/AppRouter.tsx`  (paths intact)
- `rg -n 'path="\*"|Navigate' apps/web/src/app/AppRouter.tsx`   → EXPECT NO MATCHES (no new catch-all)
- `rg -n 'LegacyRedirect to="/(catalogo|vinculos|estoque|pedidos|integracoes|precos)"' apps/web/src/app/AppRouter.tsx`  (all 6 redirects intact)

expected_artifacts: six route modules with the stable named exports above; precos.tsx owns the useClient
wrapper; AppRouter imports & renders them at unchanged paths; catalogo/estoque/integracoes/protocolos/
classifications/marketplaces + all six LegacyRedirect rows preserved verbatim; no catch-all; clean tsc except
the known env TS2688; green build.

open_questions: [] (dispatch-ready)

Report back: status; changed paths vs write_set; each command with evidence type (ran/assumed/could-not-run)
and captured output; confirm no pages/** touched and no path strings changed; anything you did NOT verify.
