import { BrowserRouter, Route, Routes } from "react-router-dom";
import { ClassificationsPage } from "@marketplace-central/feature-classifications";
import { PricingSimulatorPage } from "@marketplace-central/feature-simulator";
import { StockSeguroPage } from "@marketplace-central/feature-inventory";
import { CatalogPage } from "@marketplace-central/feature-products";
import { Layout } from "./Layout";
import { DashboardPage } from "../pages/DashboardPage";
import { AnunciosPage } from "../pages/AnunciosPage";
import { WorkspacePlaceholder } from "../pages/WorkspacePlaceholder";
import { useClient } from "./ClientContext";
import { InstallationProvider, useInstallation } from "./InstallationContext";
import { LegacyRedirect } from "./LegacyRedirect";

function CatalogPageWrapper() {
  const client = useClient();
  return <CatalogPage client={client} />;
}

function ClassificationsPageWrapper() {
  const client = useClient();
  return <ClassificationsPage client={client} />;
}

function PricingSimulatorPageWrapper() {
  const client = useClient();
  return <PricingSimulatorPage client={client} />;
}

function StockSeguroPageWrapper() {
  const client = useClient();
  const { installations } = useInstallation();
  return <StockSeguroPage client={client} installations={installations} />;
}

export function AppRouter() {
  return (
    <BrowserRouter>
      <InstallationProvider>
        <Routes>
          <Route element={<Layout />}>
            <Route index element={<DashboardPage />} />
            <Route path="/anuncios" element={<AnunciosPage />} />
            <Route path="/catalogo" element={<CatalogPageWrapper />} />
            <Route path="/catalogo/produtos/:productId" element={<WorkspacePlaceholder />} />
            <Route path="/vinculos" element={<WorkspacePlaceholder />} />
            <Route path="/estoque" element={<StockSeguroPageWrapper />} />
            <Route path="/precos" element={<PricingSimulatorPageWrapper />} />
            <Route path="/pedidos" element={<WorkspacePlaceholder />} />
            <Route path="/integracoes" element={<WorkspacePlaceholder />} />
            <Route path="/protocolos/:protocolId" element={<WorkspacePlaceholder />} />
            <Route path="/classifications" element={<ClassificationsPageWrapper />} />
            <Route path="/marketplaces" element={<WorkspacePlaceholder />} />
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
