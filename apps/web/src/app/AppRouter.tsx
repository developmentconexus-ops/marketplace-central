import { BrowserRouter, Route, Routes } from "react-router-dom";
import { ClassificationsPage } from "@marketplace-central/feature-classifications";
import { StockSeguroPage } from "@marketplace-central/feature-inventory";
import { CatalogPage } from "@marketplace-central/feature-products";
import { Layout } from "./Layout";
import { ProtocoloPage } from "../pages/mutations/ProtocoloPage";
import { WorkspacePlaceholder } from "../pages/WorkspacePlaceholder";
import { useClient } from "./ClientContext";
import { InstallationProvider, useInstallation } from "./InstallationContext";
import { LegacyRedirect } from "./LegacyRedirect";
import { DashboardRoute } from "../routes/dashboard";
import { AnunciosRoute } from "../routes/anuncios";
import { VinculosRoute } from "../routes/vinculos";
import { ProdutoRoute } from "../routes/produto";
import { PrecosRoute } from "../routes/precos";
import { PedidosRoute } from "../routes/pedidos";

function CatalogPageWrapper() {
  const client = useClient();
  return <CatalogPage client={client} />;
}

function ClassificationsPageWrapper() {
  const client = useClient();
  return <ClassificationsPage client={client} />;
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
            <Route index element={<DashboardRoute />} />
            <Route path="/anuncios" element={<AnunciosRoute />} />
            <Route path="/catalogo" element={<CatalogPageWrapper />} />
            <Route path="/catalogo/produtos/:productId" element={<ProdutoRoute />} />
            <Route path="/vinculos" element={<VinculosRoute />} />
            <Route path="/estoque" element={<StockSeguroPageWrapper />} />
            <Route path="/precos" element={<PrecosRoute />} />
            <Route path="/pedidos" element={<PedidosRoute />} />
            <Route path="/integracoes" element={<WorkspacePlaceholder />} />
            <Route path="/protocolos/:protocolId" element={<ProtocoloPage />} />
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
