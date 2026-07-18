import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { render, screen } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import { describe, expect, it, vi } from "vitest";
import { PricingPage } from "./PricingPage";
import type {
  PricingCalcProfile,
  PricingDecomposeResponse,
} from "@marketplace-central/sdk-runtime";

const profile: PricingCalcProfile = {
  regime: "SIMPLES",
  aliquota_pct: "4",
  limiar_verde_pct: "18",
  limiar_amarelo_pct: "10",
  tarifa_full: null,
  difal_enabled: false,
  difal_destino_uf: null,
  origem: "operator",
};

const decompose: PricingDecomposeResponse = {
  decomposition: {
    preco: "89.00",
    comissao: "15.13",
    taxa_fixa: "0",
    frete: "18.30",
    imposto: "3.56",
    difal: null,
    tarifa_full: null,
    custo: "61.00",
    margem_valor: "-8.99",
    margem_pct: "-10.10",
    componentes_desconhecidos: [],
  },
  blocking_state: null,
};

const productFactsPage = {
  items: [
    {
      internal_product_id: 90001,
      reference: "ELE-2210",
      description: "Eletrodo 6013",
      ean: "7890000000001",
      manufacturer_reference: "ELE-2210",
      brand_name: null,
      ncm: null,
      quality_flags: [],
      active: true,
      sellable_stock: { quantity: 10, quality: [] },
      current_price: { amount: "89.00", currency: "BRL", quality: [] },
      cost: { amount: "61.00", currency: "BRL", quality: [] },
    },
  ],
  next_cursor: null,
  page_size: 1,
  as_of: "2026-07-18T12:00:00Z",
};

const getPricingProfile = vi.fn(() => Promise.resolve(profile));
const listCatalogProductFacts = vi.fn(() => Promise.resolve(productFactsPage));
const pricingDecompose = vi.fn(() => Promise.resolve(decompose));
const pricingSolveTarget = vi.fn();

vi.mock("../../app/ClientContext", () => ({
  useClient: () => ({
    getPricingProfile,
    listCatalogProductFacts,
    pricingDecompose,
    pricingSolveTarget,
  }),
}));

function renderPage(search = "") {
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  return render(
    <QueryClientProvider client={queryClient}>
      <MemoryRouter initialEntries={[`/precos${search}`]}>
        <PricingPage />
      </MemoryRouter>
    </QueryClientProvider>,
  );
}

describe("PricingPage scaffold", () => {
  it("mounts, loads the calc profile, and renders the shell regions", async () => {
    renderPage();

    // Page marker (also the AppRouter /precos smoke assertion).
    expect(await screen.findByText("Preços & Simulador")).toBeInTheDocument();

    // Shell regions the later slices flesh out.
    expect(screen.getByTestId("region-decomposicao")).toBeInTheDocument();
    expect(screen.getByTestId("params-trigger")).toBeInTheDocument();
    expect(screen.getByTestId("region-comparacao")).toBeInTheDocument();
    expect(screen.getByTestId("region-aplicar")).toBeInTheDocument();

    // Profile drives the calc surface — it must be fetched on mount.
    expect(getPricingProfile).toHaveBeenCalled();
  });
});
