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
import { useClient } from "./ClientContext";
import { InstallationProvider } from "./InstallationContext";

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
  return <StockSeguroPage client={client} />;
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
            <Route path="/products" element={<CatalogPageWrapper />} />
            <Route path="/classifications" element={<ClassificationsPageWrapper />} />
            <Route path="/marketplaces" element={<MarketplaceSettingsPageWrapper />} />
            <Route path="/integrations" element={<IntegrationsHubPageWrapper />} />
            <Route path="/product-links" element={<ProductLinksPageWrapper />} />
            <Route path="/inventory/stock-seguro" element={<StockSeguroPageWrapper />} />
            <Route path="/orders" element={<OrdersPageWrapper />} />
            <Route path="/simulator" element={<PricingSimulatorPageWrapper />} />
          </Route>
        </Routes>
      </InstallationProvider>
    </BrowserRouter>
  );
}
