import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { render, screen, waitFor } from "@testing-library/react";
import { beforeEach, describe, expect, it, vi } from "vitest";
import { AppRouter } from "./AppRouter";

const { listIntegrationInstallations, getMutation, listMutationItems, getDashboardSummary, listOrders, listErpImports } =
  vi.hoisted(() => ({
    listIntegrationInstallations: vi.fn(),
    getMutation: vi.fn(),
    listMutationItems: vi.fn(),
    getDashboardSummary: vi.fn(),
    listOrders: vi.fn(),
    listErpImports: vi.fn(),
  }));

vi.mock("./ClientContext", () => ({
  useClient: () => ({
    mocked: true,
    listIntegrationInstallations,
    getMutation,
    listMutationItems,
    getDashboardSummary,
    listOrders,
    listErpImports,
  }),
}));

vi.mock("@marketplace-central/feature-products", () => ({
  CatalogPage: () => <div>Catalog route</div>,
}));

vi.mock("@marketplace-central/feature-classifications", () => ({
  ClassificationsPage: () => <div>Classifications route</div>,
}));

function renderAppRouter() {
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  return render(
    <QueryClientProvider client={queryClient}>
      <AppRouter />
    </QueryClientProvider>,
  );
}

vi.mock("@marketplace-central/feature-inventory", () => ({
  StockSeguroPage: ({ installations }: { installations: Array<{ installation_id: string }> }) => (
    <div>Stock Seguro route: {installations.map((installation) => installation.installation_id).join(", ")}</div>
  ),
}));

describe("AppRouter", () => {
  beforeEach(() => {
    listIntegrationInstallations.mockReset();
    getMutation.mockReset();
    listMutationItems.mockReset();
    getDashboardSummary.mockReset();
    listOrders.mockReset();
    listIntegrationInstallations.mockResolvedValue({ items: [{ installation_id: "inst_test" }] });
    getDashboardSummary.mockResolvedValue({
      sync_errors: 0,
      pending_links: 0,
      below_margin: 0,
      missing_gtin: 0,
      orders_today: 0,
      orders_7d: 0,
      anuncios_ativos: 0,
      last_sync_at: {},
      degraded: [],
      last_import: null,
      as_of: "2026-07-19T12:00:00Z",
    });
    listOrders.mockResolvedValue({ items: [], next_cursor: null });
    listErpImports.mockReset();
    listErpImports.mockResolvedValue({ items: [] });
    getMutation.mockResolvedValue({
      protocol_id: "MP-000042", installation_id: "inst_test", type: "listing_pause", state: "applied",
      actor: "operator_supplied_unverified", intent: {}, selection: {}, totals: {}, source_as_of: null,
      retried_from: null, created_at: "2026-07-17T12:00:00Z", previewed_at: null, approved_at: null,
      finished_at: "2026-07-17T12:00:03Z",
    });
    listMutationItems.mockResolvedValue({ items: [], next_cursor: null, page_size: 50 });
    window.history.pushState({}, "", "/");
  });

  it("mounts the visão-geral dashboard at the / root route", async () => {
    window.history.pushState({}, "", "/");
    renderAppRouter();
    expect(await screen.findByRole("heading", { name: "Visão geral" })).toBeInTheDocument();
    expect(listIntegrationInstallations).toHaveBeenCalledTimes(1);
  });

  it("mounts the plataforma-config screen at /integracoes with a single app-wide installation fetch", async () => {
    window.history.pushState({}, "", "/integracoes");
    renderAppRouter();
    expect(await screen.findByRole("heading", { name: "Configuração da plataforma" })).toBeInTheDocument();
    expect(listIntegrationInstallations).toHaveBeenCalledTimes(1);
  });

  it("mounts the vínculos workspace at /vinculos with a single app-wide installation fetch", async () => {
    window.history.pushState({}, "", "/vinculos");
    renderAppRouter();
    expect(await screen.findByRole("tab", { name: "Fila" })).toBeInTheDocument();
    expect(listIntegrationInstallations).toHaveBeenCalledTimes(1);
  });

  it("renders the stock seguro route at its new path", async () => {
    window.history.pushState({}, "", "/estoque");
    renderAppRouter();
    expect(await screen.findByText("Stock Seguro route: inst_test")).toBeInTheDocument();
    expect(listIntegrationInstallations).toHaveBeenCalledTimes(1);
  });

  it("keeps never-authorized installations out of the stock seguro account selector", async () => {
    // An abandoned authorization has no listings to compare, so offering it in
    // the per-screen selector can only produce an empty screen — the header
    // already hides it and both selectors must show the same account list.
    listIntegrationInstallations.mockResolvedValue({
      items: [
        { installation_id: "inst_test", status: "connected" },
        { installation_id: "inst_abandoned", status: "pending_connection" },
      ],
    });
    window.history.pushState({}, "", "/estoque");
    renderAppRouter();
    expect(await screen.findByText("Stock Seguro route: inst_test")).toBeInTheDocument();
  });

  it("mounts the pedidos workspace at /pedidos with a single app-wide installation fetch", async () => {
    window.history.pushState({}, "", "/pedidos");
    renderAppRouter();
    expect(await screen.findByRole("heading", { name: "Pedidos" })).toBeInTheDocument();
    expect(listIntegrationInstallations).toHaveBeenCalledTimes(1);
  });

  it("renders the catalog route at its new path", async () => {
    window.history.pushState({}, "", "/catalogo");
    renderAppRouter();
    expect(await screen.findByText("Catalog route")).toBeInTheDocument();
  });

  it("renders the pricing simulator route at its new path", async () => {
    window.history.pushState({}, "", "/precos");
    renderAppRouter();
    // New in-app IC-04 /precos page (PricingPage), no longer the feature-simulator stub.
    expect(await screen.findByText("Preços & Simulador")).toBeInTheDocument();
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

  it("renders the anuncios workspace", async () => {
    window.history.pushState({}, "", "/anuncios");
    renderAppRouter();
    expect(await screen.findByRole("tab", { name: "Todos" })).toBeInTheDocument();
  });

  it("renders the product workspace", async () => {
    window.history.pushState({}, "", "/catalogo/produtos/90008");
    renderAppRouter();
    expect(await screen.findByRole("tab", { name: "Veredicto" })).toBeInTheDocument();
  });

  it("mounts the real protocol page on a direct deep link", async () => {
    window.history.pushState({}, "", "/protocolos/MP-000042");
    renderAppRouter();
    expect(await screen.findByRole("heading", { name: "Protocolo MP-000042" })).toBeInTheDocument();
  });

  it("keeps the classifications route mounted", async () => {
    window.history.pushState({}, "", "/classifications");
    renderAppRouter();
    expect(await screen.findByText("Classifications route")).toBeInTheDocument();
  });

  it("mounts the em-construção stub at /marketplaces with a single app-wide installation fetch", async () => {
    window.history.pushState({}, "", "/marketplaces");
    renderAppRouter();
    expect(await screen.findByText("Em construção — disponível em breve.")).toBeInTheDocument();
    expect(listIntegrationInstallations).toHaveBeenCalledTimes(1);
  });

  it("does not render workspace content for an unknown path", () => {
    window.history.pushState({}, "", "/unknown");
    renderAppRouter();
    expect(document.body.firstElementChild).toBeEmptyDOMElement();
  });
});
