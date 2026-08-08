import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { render, screen } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";
import type { CanonicalCatalogProduct } from "@marketplace-central/sdk-runtime";
import { EstoqueTab } from "./EstoqueTab";
import type { Client } from "../../app/ClientContext";

function makeClient(product: CanonicalCatalogProduct): Client {
  return {
    getCatalogProduct: vi.fn().mockResolvedValue(product),
  } as unknown as Client;
}

function renderTab(client: Client, productId = "90008") {
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  return render(
    <QueryClientProvider client={queryClient}>
      <EstoqueTab productId={productId} client={client} />
    </QueryClientProvider>,
  );
}

const baseProduct: CanonicalCatalogProduct = {
  internal_product_id: 90008,
  name: "Produto 90008",
  ean: "7891234567890",
  manufacturer_reference: null,
  seller_sku: null,
  brand_name: null,
  product_group_name: null,
  ncm: null,
  quality_flags: [],
  cost_amount: {
    source: "erp",
    value: null,
    quality: "unknown",
    observed_at: null,
    quality_reason: null,
  },
  price_amount: {
    source: "erp",
    value: null,
    quality: "unknown",
    observed_at: null,
    quality_reason: null,
  },
  stock_quantity: {
    source: "erp",
    value: null,
    quality: "unknown",
    observed_at: null,
    quality_reason: null,
  },
};

describe("EstoqueTab", () => {
  it("renders FÍSICO with freshness, and honest DESCONHECIDO for RESERVADO/DISPONÍVEL when stock is known", async () => {
    const product: CanonicalCatalogProduct = {
      ...baseProduct,
      quality_flags: [],
      stock_quantity: {
        source: "erp",
        value: 42,
        quality: "current",
        observed_at: "2026-07-15T10:00:00Z",
        quality_reason: null,
      },
    };
    const client = makeClient(product);

    renderTab(client);

    expect(await screen.findByText("42")).toBeInTheDocument();
    const desconhecidos = screen.getAllByText("DESCONHECIDO");
    expect(desconhecidos.length).toBe(2);
    const unknowns = screen.getAllByText("—");
    expect(unknowns.length).toBe(2);
  });

  it("renders an importar-planilha CTA instead of a físico number when stock is unknown", async () => {
    const product: CanonicalCatalogProduct = {
      ...baseProduct,
      quality_flags: ["missing_stock"],
      stock_quantity: {
        source: "erp",
        value: null,
        quality: "unknown",
        observed_at: null,
        quality_reason: null,
      },
    };
    const client = makeClient(product);

    renderTab(client);

    expect(await screen.findByText(/importar planilha/i)).toBeInTheDocument();
    expect(screen.queryByText("0")).not.toBeInTheDocument();
  });

  it("renders the importar-planilha CTA when quality_flags carries missing_stock even if a stale value is present", async () => {
    const product: CanonicalCatalogProduct = {
      ...baseProduct,
      quality_flags: ["missing_stock"],
      stock_quantity: {
        source: "erp",
        value: 0,
        quality: "stale",
        observed_at: "2026-07-10T10:00:00Z",
        quality_reason: "stale erp read",
      },
    };
    const client = makeClient(product);

    renderTab(client);

    expect(await screen.findByText(/importar planilha/i)).toBeInTheDocument();
    expect(screen.queryByText("0")).not.toBeInTheDocument();
  });

  it("never renders the legacy Concorrência / Pedidos / Histórico sections", async () => {
    const product: CanonicalCatalogProduct = {
      ...baseProduct,
      quality_flags: [],
      stock_quantity: {
        source: "erp",
        value: 7,
        quality: "current",
        observed_at: "2026-07-15T10:00:00Z",
        quality_reason: null,
      },
    };
    const client = makeClient(product);

    renderTab(client);

    await screen.findByText("7");
    expect(screen.queryByText("Concorrência")).not.toBeInTheDocument();
    expect(screen.queryByText("Pedidos")).not.toBeInTheDocument();
    expect(screen.queryByText("Histórico")).not.toBeInTheDocument();
  });

  it("shows LoadingState while the query is pending", () => {
    const client = {
      getCatalogProduct: vi.fn(() => new Promise(() => undefined)),
    } as unknown as Client;

    renderTab(client);

    expect(screen.getByRole("status")).toBeInTheDocument();
  });

  it("shows ErrorState when the query fails", async () => {
    const client = {
      getCatalogProduct: vi.fn().mockRejectedValue(new Error("boom")),
    } as unknown as Client;

    renderTab(client);

    expect(await screen.findByRole("alert")).toBeInTheDocument();
  });
});
