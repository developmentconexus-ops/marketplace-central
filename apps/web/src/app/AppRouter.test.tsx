import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { render, screen, waitFor } from "@testing-library/react";
import { beforeEach, describe, expect, it, vi } from "vitest";
import { AppRouter } from "./AppRouter";

vi.mock("./ClientContext", () => ({
  useClient: () => ({
    mocked: true,
    listIntegrationInstallations: () =>
      Promise.resolve({ items: [{ installation_id: "inst_test" }] }),
  }),
}));

vi.mock("@marketplace-central/feature-products", () => ({
  CatalogPage: () => <div>Catalog route</div>,
}));

vi.mock("@marketplace-central/feature-simulator", () => ({
  PricingSimulatorPage: () => <div>Pricing route</div>,
}));

vi.mock("@marketplace-central/feature-classifications", () => ({
  ClassificationsPage: () => <div>Classifications route</div>,
}));

vi.mock("@marketplace-central/feature-marketplaces", () => ({
  MarketplaceSettingsPage: () => <div>Marketplaces route</div>,
}));

function renderAppRouter() {
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  return render(
    <QueryClientProvider client={queryClient}>
      <AppRouter />
    </QueryClientProvider>,
  );
}

vi.mock("@marketplace-central/feature-integrations", () => ({
  IntegrationsHubPage: () => <div>Integrations hub route</div>,
}));

vi.mock("@marketplace-central/feature-product-links", () => ({
  ProductLinksPage: () => <div>Product links route</div>,
}));

vi.mock("@marketplace-central/feature-inventory", () => ({
  StockSeguroPage: () => <div>Stock Seguro route</div>,
}));

vi.mock("@marketplace-central/feature-orders", () => ({
  OrdersPage: () => <div>Orders route</div>,
}));

describe("AppRouter", () => {
  beforeEach(() => {
    window.history.pushState({}, "", "/");
  });

  it("renders the integrations route at its new path", async () => {
    window.history.pushState({}, "", "/integracoes");
    renderAppRouter();
    expect(await screen.findByText("Integrations hub route")).toBeInTheDocument();
  });

  it("renders the product links route at its new path", async () => {
    window.history.pushState({}, "", "/vinculos");
    renderAppRouter();
    expect(await screen.findByText("Product links route")).toBeInTheDocument();
  });

  it("renders the stock seguro route at its new path", async () => {
    window.history.pushState({}, "", "/estoque");
    renderAppRouter();
    expect(await screen.findByText("Stock Seguro route")).toBeInTheDocument();
  });

  it("renders the orders route at its new path", async () => {
    window.history.pushState({}, "", "/pedidos");
    renderAppRouter();
    expect(await screen.findByText("Orders route")).toBeInTheDocument();
  });

  it("renders the catalog route at its new path", async () => {
    window.history.pushState({}, "", "/catalogo");
    renderAppRouter();
    expect(await screen.findByText("Catalog route")).toBeInTheDocument();
  });

  it("renders the pricing simulator route at its new path", async () => {
    window.history.pushState({}, "", "/precos");
    renderAppRouter();
    expect(await screen.findByText("Pricing route")).toBeInTheDocument();
  });

  it.each([
    ["/products", "/catalogo"],
    ["/product-links", "/vinculos"],
    ["/inventory/stock-seguro", "/estoque"],
    ["/orders", "/pedidos"],
    ["/integrations", "/integracoes"],
    ["/simulator", "/precos"],
  ])("redirects %s to %s with the full query string and replaces history", async (legacyPath, targetPath) => {
    const historyLengthBeforePush = window.history.length;
    window.history.pushState({}, "", `${legacyPath}?installation=inst_test&tab=x`);

    renderAppRouter();

    await waitFor(() => {
      expect(window.location.pathname).toBe(targetPath);
      expect(window.location.search).toBe("?installation=inst_test&tab=x");
    });
    expect(window.history.length).toBe(historyLengthBeforePush + 1);
  });

  it("renders the anuncios placeholder", async () => {
    window.history.pushState({}, "", "/anuncios");
    renderAppRouter();
    expect(await screen.findByText("Marcador da rota de anúncios.")).toBeInTheDocument();
  });

  it("renders the product workspace placeholder", async () => {
    window.history.pushState({}, "", "/catalogo/produtos/PROD1");
    renderAppRouter();
    expect(await screen.findByText("Em construção — disponível em breve.")).toBeInTheDocument();
  });

  it("renders the protocol workspace placeholder", async () => {
    window.history.pushState({}, "", "/protocolos/MP-000042");
    renderAppRouter();
    expect(await screen.findByText("Em construção — disponível em breve.")).toBeInTheDocument();
  });

  it("keeps the classifications route mounted", async () => {
    window.history.pushState({}, "", "/classifications");
    renderAppRouter();
    expect(await screen.findByText("Classifications route")).toBeInTheDocument();
  });

  it("keeps the marketplaces route mounted", async () => {
    window.history.pushState({}, "", "/marketplaces");
    renderAppRouter();
    expect(await screen.findByText("Marketplaces route")).toBeInTheDocument();
  });

  it("does not render workspace content for an unknown path", () => {
    window.history.pushState({}, "", "/unknown");
    renderAppRouter();
    expect(document.body.firstElementChild).toBeEmptyDOMElement();
  });
});
