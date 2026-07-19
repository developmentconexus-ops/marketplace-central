import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import { MemoryRouter, Route, Routes } from "react-router-dom";
import { describe, expect, it, vi } from "vitest";
import type { CanonicalCatalogProduct, MarketPriceIntelVerdict } from "@marketplace-central/sdk-runtime";
import { InstallationProvider } from "../../app/InstallationContext";
import { ProdutoPage } from "./ProdutoPage";

const getCatalogProduct = vi.fn();
const listIntegrationInstallations = vi.fn();

vi.mock("../../app/ClientContext", () => ({
  useClient: () => ({
    getCatalogProduct: (...args: unknown[]) => getCatalogProduct(...args),
    listIntegrationInstallations: (...args: unknown[]) => listIntegrationInstallations(...args),
  }),
}));

const listMarketVerdicts = vi.fn();
const collectMarketPriceIntel = vi.fn();

vi.mock("./marketClient", () => ({
  useProdutoMarketClient: () => ({
    listMarketVerdicts: (...args: unknown[]) => listMarketVerdicts(...args),
    collectMarketPriceIntel: (...args: unknown[]) => collectMarketPriceIntel(...args),
  }),
}));

const product: CanonicalCatalogProduct = {
  internal_product_id: 90008,
  name: "Parafuso sextavado",
  ean: "7891000000000",
  manufacturer_reference: null,
  seller_sku: null,
  brand_name: "Acme",
  product_group_name: null,
  ncm: "12345678",
  quality_flags: ["complete"],
  cost_amount: { source: "erp", value: 10, quality: "complete", observed_at: "2026-07-01T00:00:00Z", quality_reason: null },
  price_amount: { source: "erp", value: 20, quality: "complete", observed_at: "2026-07-01T00:00:00Z", quality_reason: null },
  stock_quantity: { source: "erp", value: 5, quality: "complete", observed_at: "2026-07-01T00:00:00Z", quality_reason: null },
};

const verdict: MarketPriceIntelVerdict = {
  match_status: "MATCHED",
  price_evidence_status: "OK",
  verdict_label: null,
  blocking_state: null,
  inputs_used: {},
  market_range: { min_valid: null, median: null, currency: "BRL", n_offers: 0, n_sellers: 0 },
};

function renderAt(path: string) {
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  return render(
    <QueryClientProvider client={queryClient}>
      <MemoryRouter initialEntries={[path]}>
        <InstallationProvider>
          <Routes>
            <Route path="/catalogo/produtos/:productId" element={<ProdutoPage />} />
          </Routes>
        </InstallationProvider>
      </MemoryRouter>
    </QueryClientProvider>,
  );
}

describe("ProdutoPage partial-failure isolation", () => {
  it("shows VeredictoBox ErrorState while ProdutoHeader still renders catalog data", async () => {
    getCatalogProduct.mockResolvedValue(product);
    listIntegrationInstallations.mockResolvedValue({ items: [] });
    listMarketVerdicts.mockRejectedValue(new Error("market down"));

    renderAt("/catalogo/produtos/90008");

    expect(await screen.findByText("Parafuso sextavado")).toBeInTheDocument();
    expect(screen.getByText("CODPROD 90008")).toBeInTheDocument();

    expect(await screen.findByText("Erro ao carregar. Não foi possível carregar o veredicto de mercado.")).toBeInTheDocument();
    expect(screen.getByRole("button", { name: "Tentar novamente" })).toBeInTheDocument();
  });

  it("shows ProdutoHeader ErrorState while the tablist survives when catalog rejects", async () => {
    getCatalogProduct.mockRejectedValue(new Error("catalog down"));
    listIntegrationInstallations.mockResolvedValue({ items: [] });
    listMarketVerdicts.mockResolvedValue([verdict]);

    renderAt("/catalogo/produtos/90008");

    expect(await screen.findByText("Erro ao carregar. Não foi possível carregar o produto.")).toBeInTheDocument();

    expect(screen.getByRole("tab", { name: "Veredicto" })).toBeInTheDocument();
    expect(screen.getByRole("tab", { name: "Anúncios vinculados" })).toBeInTheDocument();
    expect(screen.getByRole("tab", { name: "Estoque" })).toBeInTheDocument();
    expect(screen.getByRole("tab", { name: "Veredicto" })).toHaveAttribute("aria-selected", "true");
  });

  it("retries a rejected catalog query and renders resolved data after clicking Tentar novamente", async () => {
    getCatalogProduct.mockRejectedValueOnce(new Error("catalog down")).mockResolvedValue(product);
    listIntegrationInstallations.mockResolvedValue({ items: [] });
    listMarketVerdicts.mockResolvedValue([verdict]);

    renderAt("/catalogo/produtos/90008");

    expect(await screen.findByText("Erro ao carregar. Não foi possível carregar o produto.")).toBeInTheDocument();

    fireEvent.click(screen.getByRole("button", { name: "Tentar novamente" }));

    await waitFor(() => expect(screen.getByText("Parafuso sextavado")).toBeInTheDocument());
    expect(screen.queryByText("Erro ao carregar. Não foi possível carregar o produto.")).not.toBeInTheDocument();
  });
});
