import { BrowserRouter, Outlet, Route, Routes } from "react-router-dom";
import { ClassificationsPage } from "@marketplace-central/feature-classifications";
import { StockSeguroPage } from "@marketplace-central/feature-inventory";
import { CatalogPage } from "@marketplace-central/feature-products";
import { useActiveSourceQuery } from "@marketplace-central/web-query";
import { Layout } from "./Layout";
import { ProtocoloPage } from "../pages/mutations/ProtocoloPage";
import { WorkspacePlaceholder } from "../pages/WorkspacePlaceholder";
import { useClient } from "./ClientContext";
import { InstallationGate, InstallationProvider, useInstallation } from "./InstallationContext";
import { selectableInstallations } from "./installationSelection";
import { LegacyRedirect } from "./LegacyRedirect";
import { DashboardRoute } from "../routes/dashboard";
import { AnunciosRoute } from "../routes/anuncios";
import { VinculosRoute } from "../routes/vinculos";
import { ProdutoRoute } from "../routes/produto";
import { PrecosRoute } from "../routes/precos";
import { PedidosRoute } from "../routes/pedidos";
import { MercadoRoute } from "../routes/mercado";
import { IntegracoesRoute } from "../routes/integracoes";

function CatalogPageWrapper() {
  const client = useClient();
  // The tenant's active source is server state; the page only needs it to keep
  // one source's rows out of the other's cache entry.
  const activeSource = useActiveSourceQuery(client);
  return <CatalogPage client={client} erpSource={activeSource.data?.active_source} />;
}

function ClassificationsPageWrapper() {
  const client = useClient();
  return <ClassificationsPage client={client} />;
}

function StockSeguroPageWrapper() {
  const client = useClient();
  const { installationId, installations } = useInstallation();
  // Same account list as the header: an abandoned authorization has no listings
  // to compare, so offering it here can only produce an empty screen.
  return (
    <StockSeguroPage
      client={client}
      installations={selectableInstallations(installations, installationId)}
    />
  );
}

function InstallationGatedRoutes() {
  return (
    <InstallationGate>
      <Outlet />
    </InstallationGate>
  );
}

export function AppRouter() {
  return (
    <BrowserRouter>
      <InstallationProvider>
        <Routes>
          <Route element={<Layout />}>
            {/* Setup and ERP-side screens must render with no marketplace account
                connected: /integracoes is where the account is connected, and the
                catalog, stock and import screens read the ERP mirror, which exists
                before any marketplace does. */}
            <Route path="/integracoes" element={<IntegracoesRoute />} />
            <Route path="/catalogo" element={<CatalogPageWrapper />} />
            <Route path="/catalogo/produtos/:productId" element={<ProdutoRoute />} />
            <Route path="/protocolos/:protocolId" element={<ProtocoloPage />} />
            <Route path="/classifications" element={<ClassificationsPageWrapper />} />
            <Route path="/marketplaces" element={<WorkspacePlaceholder />} />
            {/* Everything below reads listings, orders or market data, which only
                exist per connected account — without one they can only render an
                empty screen, so the gate answers for them. */}
            <Route element={<InstallationGatedRoutes />}>
              <Route index element={<DashboardRoute />} />
              <Route path="/anuncios" element={<AnunciosRoute />} />
              <Route path="/mercado" element={<MercadoRoute />} />
              <Route path="/vinculos" element={<VinculosRoute />} />
              <Route path="/estoque" element={<StockSeguroPageWrapper />} />
              <Route path="/precos" element={<PrecosRoute />} />
              <Route path="/pedidos" element={<PedidosRoute />} />
            </Route>
            <Route path="/products" element={<LegacyRedirect to="/catalogo" />} />
            <Route path="/product-links" element={<LegacyRedirect to="/vinculos" />} />
            <Route path="/inventory/stock-seguro" element={<LegacyRedirect to="/estoque" />} />
            <Route path="/orders" element={<LegacyRedirect to="/pedidos" />} />
            <Route path="/integrations" element={<LegacyRedirect to="/integracoes" />} />
            <Route path="/simulator" element={<LegacyRedirect to="/precos" />} />
          </Route>
        </Routes>
      </InstallationProvider>
    </BrowserRouter>
  );
}
