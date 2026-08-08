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
  tarifa: {
    comissao: {
      valor: "15.13",
      fonte: "COTACAO",
      degrau: 3,
      data: "2026-07-18T12:00:00Z",
      estimativa: false,
    },
    frete: {
      valor: null,
      fonte: "PADRAO",
      degrau: 4,
      data: null,
      estimativa: false,
      sem_dados: true,
    },
  },
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
    {
      uf: "SP",
      interna_pct: "18",
      interestadual_pct: "12",
      efetivo_pct: "6",
      origem_versao: "padrao-2026",
      override: null,
    },
    {
      uf: "BA",
      interna_pct: "19",
      interestadual_pct: "7",
      efetivo_pct: "12",
      origem_versao: "padrao-2026",
      override: null,
    },
  ],
};

const getPricingProfile = vi.fn(() => Promise.resolve(profile));
// /precos asks for exactly the ids its resolved links name — the page never
// keysets a page of the catalog, so the fake answers the ids ask.
const catalogProductFactsByIds = vi.fn((options: { ids: number[] }) =>
  Promise.resolve({
    ...productFactsPage,
    items: productFactsPage.items.filter((p) => options.ids.includes(p.internal_product_id)),
  }),
);
// One ML listing resolved to product 90001 = the analysis set.
const listListings = vi.fn(() =>
  Promise.resolve({
    items: [
      {
        listing_id: "MLB3758134295",
        link: { state: "resolved", product_id: "90001", seller_sku: null },
      },
    ],
    next_cursor: null,
    page_size: 1,
  }),
);
const pricingDecompose = vi.fn((_req: unknown) => Promise.resolve(decompose));
const pricingSolveTarget = vi.fn();
const putPricingProfile = vi.fn((next: PricingCalcProfile) => Promise.resolve(next));
const listPricingDifal = vi.fn(() => Promise.resolve(difalList));
const putPricingDifalOverride = vi.fn(() =>
  Promise.resolve({ persisted: true, rate: difalList.items[0] }),
);
const listIntegrationInstallations = vi.fn(() =>
  Promise.resolve({ items: [{ installation_id: "inst_test" }] }),
);
const listListingsByProduct = vi.fn(() =>
  Promise.resolve({
    groups: [
      {
        product_id: "90001",
        product_title: null,
        listing_count: 1,
        group_state: "ok",
        listings: [{ listing_id: "MLB3758134295" }],
      },
    ],
    next_cursor: null,
    page_size: 1,
  }),
);

// The market comparison + apply panels own their own client seams/mutations; the
// page test stubs them so /precos page behavior runs without their side effects.
vi.mock("./SolverPanel", () => ({
  SolverPanel: () => <div data-testid="solver-panel-stub" />,
}));
vi.mock("./MarketComparison", () => ({
  MarketComparison: () => <div data-testid="market-comparison-stub" />,
}));
vi.mock("./QuickChips", () => ({
  QuickChips: () => <div data-testid="quick-chips-stub" />,
}));
vi.mock("./ApplyPriceAction", () => ({
  ApplyPriceAction: () => <div data-testid="apply-action-stub" />,
}));
// The stub exposes the onReload seam so the page's applyScenario logic (state
// re-application) is exercised for real — both a present and an absent product.
// The matrix owns its own market + per-row decompose/listing fan-out; stub it here
// (same sanctioned pattern as the child panels below) and expose its onSelect seam
// so the page's row-click → setSelectedId wiring is exercised for real.
vi.mock("./PricingMatrix", () => ({
  PricingMatrix: ({ onSelect }: { onSelect: (id: number) => void }) => (
    <div data-testid="pricing-matrix-stub">
      <button data-testid="stub-row-90001" onClick={() => onSelect(90001)}>
        row-90001
      </button>
      <button data-testid="stub-row-missing" onClick={() => onSelect(99999)}>
        row-missing
      </button>
    </div>
  ),
}));
vi.mock("./ScenariosPanel", () => ({
  ScenariosPanel: ({ onReload }: { onReload: (p: Record<string, unknown>) => void }) => (
    <div data-testid="scenarios-panel-stub">
      <button
        data-testid="stub-reload-valid"
        onClick={() => onReload({ product_id: 90001, preco: "77.00", modalidade: "classico" })}
      >
        reload-valid
      </button>
      <button
        data-testid="stub-reload-missing"
        onClick={() => onReload({ product_id: 99999, preco: "50.00", modalidade: "classico" })}
      >
        reload-missing
      </button>
    </div>
  ),
}));

// The page reads the workspace account from the installation context (which
// prefers a CONNECTED installation) — never from listIntegrationInstallations[0].
vi.mock("../../app/InstallationContext", () => ({
  useInstallation: () => ({
    installationId: "inst_test",
    setInstallationId: () => undefined,
    installations: [],
    status: "ready",
  }),
}));

vi.mock("../../app/ClientContext", () => ({
  useClient: () => ({
    getPricingProfile,
    catalogProductFactsByIds,
    listListings,
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
    pricingDecompose.mockClear();
  });

  it("mounts, loads the calc profile, and renders the shell regions", async () => {
    renderPage();

    // Page marker (also the AppRouter /precos smoke assertion).
    expect(await screen.findByText("Preços & Simulador")).toBeInTheDocument();

    // The simular side panel opens for the auto-selected first product once the
    // catalog resolves — its regions are the preserved single-product surface.
    expect(await screen.findByTestId("region-decomposicao")).toBeInTheDocument();
    expect(screen.getByTestId("region-solver")).toBeInTheDocument();
    expect(screen.getByTestId("params-trigger")).toBeInTheDocument();
    expect(screen.getByTestId("region-comparacao")).toBeInTheDocument();
    expect(screen.getByTestId("region-aplicar")).toBeInTheDocument();
    expect(screen.getByTestId("region-cenarios")).toBeInTheDocument();

    // Profile drives the calc surface — it must be fetched on mount.
    expect(getPricingProfile).toHaveBeenCalled();

    // The decompose response's tarifa block flows through to the panel's carimbo.
    const carimbo = await screen.findByTestId("decomp-tarifa-comissao");
    expect(carimbo).toHaveTextContent("Cotação");
  });

  it("renders the matrix as the main surface and wires row-click to selection", async () => {
    renderPage();
    // Matrix is the main surface.
    expect(await screen.findByTestId("pricing-matrix-stub")).toBeInTheDocument();

    // A row-click for a product NOT in the loaded catalog resolves to the honest
    // "not in list" empty state — proving onSelect drives setSelectedId (never a
    // silent fallback to a different product).
    fireEvent.click(screen.getByTestId("stub-row-missing"));
    expect(await screen.findByTestId("scenario-reload-notice")).toBeInTheDocument();
  });

  it("toggles the Produtos picker and lists the loaded catalog product", async () => {
    renderPage();
    const pill = await screen.findByTestId("produtos-pill");
    // Count reflects the loaded catalog once it resolves (starts at 0).
    await waitFor(() => expect(pill).toHaveTextContent("Produtos: 1"));

    fireEvent.click(pill);
    const dropdown = await screen.findByTestId("produtos-dropdown");
    expect(dropdown).toHaveTextContent("Eletrodo 6013");

    // Selecting the option closes the picker (product stays selected).
    fireEvent.click(screen.getByTestId("produtos-option-90001"));
    await waitFor(() => expect(screen.queryByTestId("produtos-dropdown")).toBeNull());
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
    fireEvent.click(screen.getByRole("button", { name: "Concluído" }));

    await waitFor(() => expect(putPricingProfile).toHaveBeenCalledTimes(1));
    expect(putPricingProfile.mock.calls[0][0]).toMatchObject({ aliquota_pct: "9.25" });
    // onSuccess closes the drawer.
    await waitFor(() => expect(screen.queryByTestId("params-drawer")).toBeNull());
  });

  it("reloads a scenario into page state and re-decomposes for that product", async () => {
    renderPage();
    await screen.findByText("Preços & Simulador");
    await waitFor(() => expect(pricingDecompose).toHaveBeenCalled());
    pricingDecompose.mockClear();

    fireEvent.click(screen.getByTestId("stub-reload-valid"));

    await waitFor(() => expect(pricingDecompose).toHaveBeenCalled());
    const lastArg = pricingDecompose.mock.calls.at(-1)![0] as Record<string, unknown>;
    expect(lastArg).toMatchObject({ preco: "77.00", modalidade: "classico", product_id: 90001 });
    // comissao_pct must be OMITTED so the backend resolver chain runs (COTACAO/PADRAO),
    // never a hardcoded pct that forces the response into a MANUAL override.
    expect(lastArg).not.toHaveProperty("comissao_pct");
  });

  it("does not silently load a different product when a reloaded scenario's product is absent", async () => {
    renderPage();
    await screen.findByText("Preços & Simulador");
    await waitFor(() => expect(pricingDecompose).toHaveBeenCalled());

    fireEvent.click(screen.getByTestId("stub-reload-missing"));

    // Honest empty state + notice — NEVER a decomposition for the wrong product.
    expect(await screen.findByTestId("scenario-reload-notice")).toBeInTheDocument();
    expect(screen.getByText("Selecione um produto e um preço para simular.")).toBeInTheDocument();
    const badCall = pricingDecompose.mock.calls.some(
      ([arg]) =>
        (arg as { product_id?: number }).product_id === 99999 ||
        (arg as { preco?: string }).preco === "50.00",
    );
    expect(badCall).toBe(false);
  });

  it("normalizes a pt-BR comma price to dot-decimal before the decompose SDK call (F-P7-2)", async () => {
    renderPage();
    await screen.findByText("Preços & Simulador");
    await waitFor(() => expect(pricingDecompose).toHaveBeenCalled());
    pricingDecompose.mockClear();

    // Operator types pt-BR "150,00"; the raw string stays in the input, but the
    // value bound to the SDK must be dot-decimal or the Go parser answers 422.
    const priceInput = screen.getByLabelText("Preço de venda") as HTMLInputElement;
    fireEvent.change(priceInput, { target: { value: "150,00" } });

    await waitFor(() => {
      const lastArg = pricingDecompose.mock.calls.at(-1)?.[0] as { preco?: string } | undefined;
      expect(lastArg?.preco).toBe("150.00");
    });
    // The input still shows the raw pt-BR string the operator typed.
    expect(priceInput.value).toBe("150,00");
  });

  it("opens the DIFAL drawer via the Parâmetros drawer link and lists the UF table with the disclaimer", async () => {
    renderPage();
    // DIFAL is reached through the Parâmetros drawer's "ver tabela completa por UF"
    // link (design parity: no standalone page button).
    fireEvent.click(await screen.findByTestId("params-trigger"));
    await screen.findByTestId("params-drawer");
    fireEvent.click(screen.getByRole("button", { name: /ver tabela completa por UF/ }));

    expect(await screen.findByTestId("difal-drawer")).toBeInTheDocument();
    await waitFor(() => expect(listPricingDifal).toHaveBeenCalled());
    expect(await screen.findByTestId("difal-disclaimer")).toHaveTextContent(
      "não é orientação fiscal",
    );
    expect(screen.getByRole("row", { name: /^SP/ })).toBeInTheDocument();
  });

  it("asks the catalog for exactly the linked product ids, deep link included", async () => {
    catalogProductFactsByIds.mockClear();
    renderPage("?produto=90001");
    await screen.findByText("Preços & Simulador");

    await waitFor(() => expect(catalogProductFactsByIds).toHaveBeenCalled());
    // The deep-linked id and the linked listing's id are the SAME product here, so
    // the ask must carry it once — never a keyset page of the whole catalog.
    expect(catalogProductFactsByIds.mock.calls.at(-1)![0]).toEqual({ ids: [90001] });
  });

  it("says why there is nothing to price instead of rendering an empty matrix", async () => {
    listListings.mockResolvedValueOnce({ items: [], next_cursor: null, page_size: 0 });
    renderPage();

    const empty = await screen.findByTestId("pricing-empty");
    expect(empty).toHaveTextContent("Nenhum anúncio vinculado a um produto");
    expect(screen.queryByTestId("pricing-matrix-stub")).toBeNull();
  });
});
