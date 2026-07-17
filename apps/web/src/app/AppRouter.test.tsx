import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { render, screen, waitFor } from "@testing-library/react";
import { beforeEach, describe, expect, it, vi } from "vitest";
import { AppRouter } from "./AppRouter";

const { listIntegrationInstallations, getMutation, listMutationItems } = vi.hoisted(() => ({
  listIntegrationInstallations: vi.fn(),
  getMutation: vi.fn(),
  listMutationItems: vi.fn(),
}));

vi.mock("./ClientContext", () => ({
  useClient: () => ({ mocked: true, listIntegrationInstallations, getMutation, listMutationItems }),
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
    listIntegrationInstallations.mockResolvedValue({ items: [{ installation_id: "inst_test" }] });
    getMutation.mockResolvedValue({
      protocol_id: "MP-000042", installation_id: "inst_test", type: "listing_pause", state: "applied",
      actor: "operator_supplied_unverified", intent: {}, selection: {}, totals: {}, source_as_of: null,
      retried_from: null, created_at: "2026-07-17T12:00:00Z", previewed_at: null, approved_at: null,
      finished_at: "2026-07-17T12:00:03Z",
    });
    listMutationItems.mockResolvedValue({ items: [], next_cursor: null, page_size: 50 });
    window.history.pushState({}, "", "/");
  });

  it("mounts the em-construção stub at /integracoes with a single app-wide installation fetch", async () => {
    window.history.pushState({}, "", "/integracoes");
    renderAppRouter();
    expect(await screen.findByText("Em construção — disponível em breve.")).toBeInTheDocument();
    expect(listIntegrationInstallations).toHaveBeenCalledTimes(1);
  });

  it("mounts the em-construção stub at /vinculos with a single app-wide installation fetch", async () => {
    window.history.pushState({}, "", "/vinculos");
    renderAppRouter();
    expect(await screen.findByText("Em construção — disponível em breve.")).toBeInTheDocument();
    expect(listIntegrationInstallations).toHaveBeenCalledTimes(1);
  });

  it("renders the stock seguro route at its new path", async () => {
    window.history.pushState({}, "", "/estoque");
    renderAppRouter();
    expect(await screen.findByText("Stock Seguro route: inst_test")).toBeInTheDocument();
    expect(listIntegrationInstallations).toHaveBeenCalledTimes(1);
  });

  it("mounts the em-construção stub at /pedidos with a single app-wide installation fetch", async () => {
    window.history.pushState({}, "", "/pedidos");
    renderAppRouter();
    expect(await screen.findByText("Em construção — disponível em breve.")).toBeInTheDocument();
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

  it("renders the anuncios workspace", async () => {
    window.history.pushState({}, "", "/anuncios");
    renderAppRouter();
    expect(await screen.findByRole("tab", { name: "Todos" })).toBeInTheDocument();
  });

  it("renders the product workspace placeholder", async () => {
    window.history.pushState({}, "", "/catalogo/produtos/PROD1");
    renderAppRouter();
    expect(await screen.findByText("Em construção — disponível em breve.")).toBeInTheDocument();
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
