import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { render, screen } from "@testing-library/react";
import { beforeEach, describe, expect, it, vi } from "vitest";
import { AppRouter } from "./AppRouter";

vi.mock("./ClientContext", () => ({
  useClient: () => ({
    mocked: true,
    listIntegrationInstallations: () =>
      Promise.resolve({ items: [{ installation_id: "inst_test" }] }),
  }),
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
    window.history.pushState({}, "", "/integrations");
  });

  it("renders the integrations route", async () => {
    renderAppRouter();
    expect(await screen.findByText("Integrations hub route")).toBeInTheDocument();
  });

  it("renders the product links route", async () => {
    window.history.pushState({}, "", "/product-links");
    renderAppRouter();
    expect(await screen.findByText("Product links route")).toBeInTheDocument();
  });

  it("renders the stock seguro route", async () => {
    window.history.pushState({}, "", "/inventory/stock-seguro");
    renderAppRouter();
    expect(await screen.findByText("Stock Seguro route")).toBeInTheDocument();
  });

  it("renders the orders route", async () => {
    window.history.pushState({}, "", "/orders");
    renderAppRouter();
    expect(await screen.findByText("Orders route")).toBeInTheDocument();
  });
});
