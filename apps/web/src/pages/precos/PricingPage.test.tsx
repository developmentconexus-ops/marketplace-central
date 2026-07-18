import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import { beforeEach, describe, expect, it, vi } from "vitest";
import { PricingPage } from "./PricingPage";
import type {
  PricingCalcProfile,
  PricingDecomposeResponse,
  PricingDifalListResponse,
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

const difalList: PricingDifalListResponse = {
  disclaimer: "seed padrão 2026 — não é orientação fiscal",
  items: [
    { uf: "SP", interna_pct: "18", interestadual_pct: "12", efetivo_pct: "6", origem_versao: "padrao-2026", override: null },
    { uf: "BA", interna_pct: "19", interestadual_pct: "7", efetivo_pct: "12", origem_versao: "padrao-2026", override: null },
  ],
};

const getPricingProfile = vi.fn(() => Promise.resolve(profile));
const listCatalogProductFacts = vi.fn(() => Promise.resolve(productFactsPage));
const pricingDecompose = vi.fn(() => Promise.resolve(decompose));
const pricingSolveTarget = vi.fn();
const putPricingProfile = vi.fn((next: PricingCalcProfile) => Promise.resolve(next));
const listPricingDifal = vi.fn(() => Promise.resolve(difalList));
const putPricingDifalOverride = vi.fn(() => Promise.resolve({ persisted: true, rate: difalList.items[0] }));
const listIntegrationInstallations = vi.fn(() =>
  Promise.resolve({ items: [{ installation_id: "inst_test" }] }),
);
const listListingsByProduct = vi.fn(() =>
  Promise.resolve({ items: [{ listing_id: "MLB3758134295" }], next_cursor: null, page_size: 1 }),
);

// The market comparison + apply panels own their own client seams/mutations; the
// page test stubs them so /precos page behavior runs without their side effects.
vi.mock("./MarketComparison", () => ({
  MarketComparison: () => <div data-testid="market-comparison-stub" />,
}));
vi.mock("./ApplyPriceAction", () => ({
  ApplyPriceAction: () => <div data-testid="apply-action-stub" />,
}));

vi.mock("../../app/ClientContext", () => ({
  useClient: () => ({
    getPricingProfile,
    listCatalogProductFacts,
    pricingDecompose,
    pricingSolveTarget,
    putPricingProfile,
    listPricingDifal,
    putPricingDifalOverride,
    listIntegrationInstallations,
    listListingsByProduct,
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
  beforeEach(() => {
    putPricingProfile.mockClear();
    listPricingDifal.mockClear();
    putPricingDifalOverride.mockClear();
  });

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

  it("opens the Parâmetros drawer on the ?params=1 deep link", async () => {
    renderPage("?params=1");
    expect(await screen.findByTestId("params-drawer")).toBeInTheDocument();
  });

  it("saves an edited profile and closes the drawer", async () => {
    renderPage();
    fireEvent.click(await screen.findByTestId("params-trigger"));
    const drawer = await screen.findByTestId("params-drawer");
    expect(drawer).toBeInTheDocument();

    fireEvent.change(screen.getByLabelText("Alíquota"), { target: { value: "9.25" } });
    fireEvent.click(screen.getByRole("button", { name: "Salvar parâmetros" }));

    await waitFor(() => expect(putPricingProfile).toHaveBeenCalledTimes(1));
    expect(putPricingProfile.mock.calls[0][0]).toMatchObject({ aliquota_pct: "9.25" });
    // onSuccess closes the drawer.
    await waitFor(() => expect(screen.queryByTestId("params-drawer")).toBeNull());
  });

  it("opens the DIFAL drawer and lists the UF table with the disclaimer", async () => {
    renderPage();
    fireEvent.click(await screen.findByTestId("difal-trigger"));

    expect(await screen.findByTestId("difal-drawer")).toBeInTheDocument();
    await waitFor(() => expect(listPricingDifal).toHaveBeenCalled());
    expect(await screen.findByTestId("difal-disclaimer")).toHaveTextContent("não é orientação fiscal");
    expect(screen.getByRole("row", { name: /^SP/ })).toBeInTheDocument();
  });
});
