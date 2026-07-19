import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { render, screen } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";
import { ProdutoHeader } from "./ProdutoHeader";

const getCatalogProduct = vi.fn();

vi.mock("../../app/ClientContext", () => ({
  useClient: () => ({
    getCatalogProduct: (...args: unknown[]) => getCatalogProduct(...args),
  }),
}));

function renderHeader(productId = "90001") {
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  return render(
    <QueryClientProvider client={queryClient}>
      <ProdutoHeader productId={productId} />
    </QueryClientProvider>,
  );
}

describe("ProdutoHeader", () => {
  it("renders honest unknowns and the sem-EAN chip when identity/numeric facts are null", async () => {
    getCatalogProduct.mockResolvedValue({
      internal_product_id: 90001,
      name: "X",
      ean: null,
      manufacturer_reference: "REF9",
      seller_sku: null,
      brand_name: null,
      product_group_name: null,
      ncm: null,
      quality_flags: ["missing_cost", "invalid_ean"],
      cost_amount: { source: "erp", value: null, quality: "unknown", observed_at: null, quality_reason: null },
      price_amount: { source: "erp", value: null, quality: "unknown", observed_at: null, quality_reason: null },
      stock_quantity: { source: "erp", value: null, quality: "unknown", observed_at: null, quality_reason: null },
    });

    renderHeader();

    expect(await screen.findByText("sem EAN — matching limitado a REVIEW")).toBeInTheDocument();
    expect(screen.getByText("REF9")).toBeInTheDocument();
    expect(screen.getByText("missing_cost")).toBeInTheDocument();
    expect(screen.getByText("invalid_ean")).toBeInTheDocument();

    const unknowns = screen.getAllByText("—");
    expect(unknowns.length).toBeGreaterThanOrEqual(2);

    expect(screen.queryByText("0")).not.toBeInTheDocument();
    expect(getCatalogProduct).toHaveBeenCalledWith(90001);
  });

  it("renders full facts without the sem-EAN chip when values are present", async () => {
    getCatalogProduct.mockResolvedValue({
      internal_product_id: 90002,
      name: "Produto completo",
      ean: "7891234567890",
      manufacturer_reference: "REF1",
      seller_sku: "SKU1",
      brand_name: "Marca X",
      product_group_name: "Grupo",
      ncm: "12345678",
      quality_flags: ["complete"],
      cost_amount: {
        source: "erp",
        value: 12.5,
        quality: "current",
        observed_at: "2026-07-15T10:00:00Z",
        quality_reason: null,
      },
      price_amount: { source: "erp", value: 20, quality: "current", observed_at: "2026-07-15T10:00:00Z", quality_reason: null },
      stock_quantity: { source: "erp", value: 5, quality: "current", observed_at: "2026-07-15T10:00:00Z", quality_reason: null },
    });

    renderHeader("90002");

    expect(await screen.findByText("7891234567890")).toBeInTheDocument();
    expect(screen.queryByText("sem EAN — matching limitado a REVIEW")).not.toBeInTheDocument();
    expect(screen.getByText("12.5")).toBeInTheDocument();
    expect(screen.getByText("complete")).toBeInTheDocument();
    expect(screen.getByText("REF1")).toBeInTheDocument();
  });

  it("keeps REFFORN visible when EAN is present (C01 identity field never dropped)", async () => {
    getCatalogProduct.mockResolvedValue({
      internal_product_id: 90003,
      name: "Produto com EAN e REFFORN",
      ean: "7891234567890",
      manufacturer_reference: "REF-EAN-1",
      seller_sku: null,
      brand_name: null,
      product_group_name: null,
      ncm: null,
      quality_flags: ["complete"],
      cost_amount: { source: "erp", value: null, quality: "unknown", observed_at: null, quality_reason: null },
      price_amount: { source: "erp", value: null, quality: "unknown", observed_at: null, quality_reason: null },
      stock_quantity: { source: "erp", value: null, quality: "unknown", observed_at: null, quality_reason: null },
    });

    renderHeader("90003");

    expect(await screen.findByText("7891234567890")).toBeInTheDocument();
    expect(screen.getByText("REF-EAN-1")).toBeInTheDocument();
  });

  it("renders honest UnknownValue for a null manufacturer_reference even with an EAN present", async () => {
    getCatalogProduct.mockResolvedValue({
      internal_product_id: 90004,
      name: "Produto sem REFFORN",
      ean: "7891234567890",
      manufacturer_reference: null,
      seller_sku: null,
      brand_name: null,
      product_group_name: null,
      ncm: null,
      quality_flags: ["complete"],
      cost_amount: { source: "erp", value: null, quality: "unknown", observed_at: null, quality_reason: null },
      price_amount: { source: "erp", value: null, quality: "unknown", observed_at: null, quality_reason: null },
      stock_quantity: { source: "erp", value: null, quality: "unknown", observed_at: null, quality_reason: null },
    });

    renderHeader("90004");

    expect(await screen.findByText("7891234567890")).toBeInTheDocument();
    expect(screen.getByText("REF")).toBeInTheDocument();
    const unknowns = screen.getAllByText("—");
    expect(unknowns.length).toBeGreaterThanOrEqual(1);
  });
});
