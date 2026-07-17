import { BrowserRouter, Route, Routes } from "react-router-dom";
import { ClassificationsPage } from "@marketplace-central/feature-classifications";
import { MarketplaceSettingsPage } from "@marketplace-central/feature-marketplaces";
import { PricingSimulatorPage } from "@marketplace-central/feature-simulator";
import { StockSeguroPage } from "@marketplace-central/feature-inventory";
import { CatalogPage } from "@marketplace-central/feature-products";
import { IntegrationsHubPage } from "@marketplace-central/feature-integrations";
import { ProductLinksPage } from "@marketplace-central/feature-product-links";
import { OrdersPage } from "@marketplace-central/feature-orders";
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

function MarketplaceSettingsPageWrapper() {
  const client = useClient();
  return <MarketplaceSettingsPage client={client} />;
}

function PricingSimulatorPageWrapper() {
  const client = useClient();
  return <PricingSimulatorPage client={client} />;
}

function IntegrationsHubPageWrapper() {
  const client = useClient();
  return <IntegrationsHubPage client={client} />;
}

function ProductLinksPageWrapper() {
  const client = useClient();
  return <ProductLinksPage client={client} />;
}

function StockSeguroPageWrapper() {
  const client = useClient();
  const { installations } = useInstallation();
  return <StockSeguroPage client={client} installations={installations} />;
}

function OrdersPageWrapper() {
  const client = useClient();
  return <OrdersPage client={client} />;
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
            <Route path="/vinculos" element={<ProductLinksPageWrapper />} />
            <Route path="/estoque" element={<StockSeguroPageWrapper />} />
            <Route path="/precos" element={<PricingSimulatorPageWrapper />} />
            <Route path="/pedidos" element={<OrdersPageWrapper />} />
            <Route path="/integracoes" element={<IntegrationsHubPageWrapper />} />
            <Route path="/protocolos/:protocolId" element={<WorkspacePlaceholder />} />
            <Route path="/classifications" element={<ClassificationsPageWrapper />} />
            <Route path="/marketplaces" element={<MarketplaceSettingsPageWrapper />} />
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
