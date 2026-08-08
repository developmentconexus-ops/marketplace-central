import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { fireEvent, render, screen } from "@testing-library/react";
import { MemoryRouter, Route, Routes } from "react-router-dom";
import { describe, expect, it, vi } from "vitest";
import type { CanonicalCatalogProduct, MarketPriceIntelVerdict } from "@marketplace-central/sdk-runtime";
import { ProdutoPage } from "./ProdutoPage";

const { getCatalogProduct, listMarketVerdicts } = vi.hoisted(() => ({
  getCatalogProduct: vi.fn(() => new Promise(() => undefined)),
  listMarketVerdicts: vi.fn(() => new Promise(() => undefined)),
}));

vi.mock("../../app/ClientContext", () => ({
  useClient: () => ({
    getCatalogProduct,
  }),
}));

vi.mock("./marketClient", () => ({
  useProdutoMarketClient: () => ({
    listMarketVerdicts,
    collectMarketPriceIntel: vi.fn(() => new Promise(() => undefined)),
  }),
}));

const catalogProductFixture: CanonicalCatalogProduct = {
  internal_product_id: 90008,
  name: "Produto de teste",
  ean: "7890000000001",
  manufacturer_reference: "REF-90008",
  seller_sku: "SKU-90008",
  brand_name: "Marca Teste",
  product_group_name: "Grupo Teste",
  ncm: "12345678",
  quality_flags: ["complete"],
  cost_amount: { source: "erp", value: 42.5, quality: "current", observed_at: "2026-07-19T00:00:00Z", quality_reason: null },
  price_amount: { source: "erp", value: 99.9, quality: "current", observed_at: "2026-07-19T00:00:00Z", quality_reason: null },
  stock_quantity: { source: "erp", value: 10, quality: "current", observed_at: "2026-07-19T00:00:00Z", quality_reason: null },
};

function renderAt(path: string) {
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  return render(
    <QueryClientProvider client={queryClient}>
      <MemoryRouter initialEntries={[path]}>
        <Routes>
          <Route path="/catalogo/produtos/:productId" element={<ProdutoPage />} />
        </Routes>
      </MemoryRouter>
    </QueryClientProvider>,
  );
}

describe("ProdutoPage", () => {
  it("renders the three tabs and defaults to Veredicto", async () => {
    renderAt("/catalogo/produtos/90008");

    expect(await screen.findByRole("tab", { name: "Veredicto" })).toBeInTheDocument();
    expect(screen.getByRole("tab", { name: "Anúncios vinculados" })).toBeInTheDocument();
    expect(screen.getByRole("tab", { name: "Estoque" })).toBeInTheDocument();
    expect(screen.getByRole("tab", { name: "Veredicto" })).toHaveAttribute("aria-selected", "true");
  });

  it("restores the active tab from the ?tab= deep link", async () => {
    renderAt("/catalogo/produtos/90008?tab=estoque");

    expect(await screen.findByRole("tab", { name: "Estoque" })).toHaveAttribute("aria-selected", "true");
    expect(screen.getByRole("tab", { name: "Veredicto" })).toHaveAttribute("aria-selected", "false");
  });

  it("updates the URL tab param when a tab is clicked", async () => {
    renderAt("/catalogo/produtos/90008");

    fireEvent.click(await screen.findByRole("tab", { name: "Estoque" }));

    expect(screen.getByRole("tab", { name: "Estoque" })).toHaveAttribute("aria-selected", "true");
    expect(screen.getByRole("tab", { name: "Veredicto" })).toHaveAttribute("aria-selected", "false");
  });

  it("renders a not-found ErrorState with a link to /catalogo for an invalid productId", async () => {
    renderAt("/catalogo/produtos/abc");

    expect(await screen.findByRole("alert")).toBeInTheDocument();
    expect(screen.getByRole("link", { name: "Voltar para o catálogo" })).toHaveAttribute("href", "/catalogo");
    expect(screen.queryByRole("tab")).not.toBeInTheDocument();
  });

  it("wires the shared catalog product query's cost_amount into the Veredicto box's Custo ERP row", async () => {
    getCatalogProduct.mockResolvedValueOnce(catalogProductFixture);
    const verdict: MarketPriceIntelVerdict = {
      match_status: "ACCEPT",
      price_evidence_status: "OK",
      verdict_label: null,
      blocking_state: null,
      inputs_used: {},
      market_range: { min_valid: null, median: null, currency: "BRL", n_offers: 0, n_sellers: 5 },
    };
    listMarketVerdicts.mockResolvedValueOnce([verdict]);

    renderAt("/catalogo/produtos/90008");

    const custoErpLabel = await screen.findByText("Custo ERP");
    expect(custoErpLabel.nextElementSibling).toHaveTextContent(String(catalogProductFixture.cost_amount.value));
  });
});
