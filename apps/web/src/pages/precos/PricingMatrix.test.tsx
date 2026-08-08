import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { render, screen, waitFor, within } from "@testing-library/react";
import { beforeEach, describe, expect, it, vi } from "vitest";
import { PricingMatrix } from "./PricingMatrix";
import type {
  CatalogProductFact,
  PricingCalcProfile,
  PricingDecomposeResponse,
} from "@marketplace-central/sdk-runtime";
import type { MarketPriceIntelAggregate } from "@marketplace-central/sdk-runtime";
import type { MarketAggregatesClient } from "./MarketComparison";

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

function fact(
  over: Partial<CatalogProductFact> & { internal_product_id: number },
): CatalogProductFact {
  return {
    reference: null,
    description: null,
    ean: null,
    manufacturer_reference: null,
    brand_name: null,
    ncm: null,
    quality_flags: [],
    active: true,
    sellable_stock: { quantity: 1, quality: [] },
    current_price: { amount: "189.00", currency: "BRL", quality: [] },
    cost: { amount: "120.00", currency: "BRL", quality: [] },
    ...over,
  };
}

// EXEMPLO-IO three products:
// (a) 90001 custo 120 preço 189, mercado mediana 169 OK n_sellers 8, margem non-null
// (b) 90002 same custo/preço, mercado NO_PRICE_EVIDENCE, margem still computed
// (c) 90003 cost null, decompose margem null
const products: CatalogProductFact[] = [
  fact({ internal_product_id: 90001, manufacturer_reference: "SKU-A", description: "Produto A" }),
  fact({ internal_product_id: 90002, manufacturer_reference: "SKU-B", description: "Produto B" }),
  fact({
    internal_product_id: 90003,
    manufacturer_reference: "SKU-C",
    description: "Produto C",
    cost: { amount: null, currency: "BRL", quality: [] },
  }),
];

const aggregates: MarketPriceIntelAggregate[] = [
  {
    product_id: "90001",
    median: { amount: "169.00", currency: "BRL" },
    min_valid: { amount: "150.00", currency: "BRL" },
    n_offers: 16,
    n_sellers: 8,
    source: "ml_catalog_offers",
    fetched_at: "2026-07-18T12:00:00Z",
    computed_at: "2026-07-18T12:00:00Z",
    status: "OK",
  },
  {
    product_id: "90002",
    median: null,
    min_valid: null,
    n_offers: 0,
    n_sellers: 0,
    source: "ml_catalog_offers",
    fetched_at: "0001-01-01T00:00:00Z",
    computed_at: "0001-01-01T00:00:00Z",
    status: "NO_PRICE_EVIDENCE",
  },
  // 90003 intentionally absent from aggregates → honest "SEM_EVIDENCIA"
];

function decomposeFor(productId: number): PricingDecomposeResponse {
  if (productId === 90003) {
    return {
      decomposition: {
        preco: "189.00",
        comissao: "32.13",
        taxa_fixa: "0",
        frete: null,
        imposto: "7.56",
        difal: null,
        tarifa_full: null,
        custo: null,
        margem_valor: null,
        margem_pct: null,
        componentes_desconhecidos: ["custo_erp"],
      },
      blocking_state: "SEM_CUSTO",
      tarifa: null,
    };
  }
  return {
    decomposition: {
      preco: "189.00",
      comissao: "32.13",
      taxa_fixa: "0",
      frete: "18.30",
      imposto: "7.56",
      difal: null,
      tarifa_full: null,
      custo: "120.00",
      margem_valor: "11.01",
      margem_pct: "5.83",
      componentes_desconhecidos: [],
    },
    blocking_state: null,
    tarifa: null,
  };
}

const pricingDecompose = vi.fn((req: { product_id?: number | null }) =>
  Promise.resolve(decomposeFor(req.product_id as number)),
);

// 90001 has a listing; 90003 has none → "novo".
const listListingsByProduct = vi.fn((req: { product_id: string }) =>
  Promise.resolve({
    groups: [
      {
        product_id: req.product_id,
        product_title: null,
        listing_count: req.product_id === "90001" ? 1 : 0,
        group_state: "ok",
        listings: req.product_id === "90001" ? [{ listing_id: "MLB1" }] : [],
      },
    ],
    next_cursor: null,
    page_size: 1,
  }),
);

vi.mock("../../app/ClientContext", () => ({
  useClient: () => ({ pricingDecompose, listListingsByProduct }),
}));

function renderMatrix(marketClient: MarketAggregatesClient, installationId = "") {
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  return render(
    <QueryClientProvider client={queryClient}>
      <PricingMatrix
        products={products}
        selectedId={null}
        onSelect={() => {}}
        modalidade="premium"
        profile={profile}
        installationId={installationId}
        marketClient={marketClient}
      />
    </QueryClientProvider>,
  );
}

describe("PricingMatrix (EXEMPLO-IO golden)", () => {
  beforeEach(() => {
    pricingDecompose.mockClear();
    listListingsByProduct.mockClear();
    // Restore the default resolving impl so a per-test error override never leaks.
    pricingDecompose.mockImplementation((req: { product_id?: number | null }) =>
      Promise.resolve(decomposeFor(req.product_id as number)),
    );
  });

  it("renders the design columns for a multi-product list", async () => {
    const listMarketAggregates = vi.fn(() => Promise.resolve(aggregates));
    renderMatrix({ listMarketAggregates });

    // 7 design columns present.
    for (const col of [
      "SKU",
      "DESCRIÇÃO",
      "CUSTO",
      "NOSSO PREÇO",
      "PREÇO MERCADO",
      "MARGEM",
      "VEREDICTO",
    ]) {
      expect(screen.getByRole("columnheader", { name: col })).toBeInTheDocument();
    }
    // All three product rows rendered.
    expect(await screen.findByTestId("matrix-row-90001")).toBeInTheDocument();
    expect(screen.getByTestId("matrix-row-90002")).toBeInTheDocument();
    expect(screen.getByTestId("matrix-row-90003")).toBeInTheDocument();
  });

  it("(a) priced row: market median + margem chip + VEREDICTO OK", async () => {
    const listMarketAggregates = vi.fn(() => Promise.resolve(aggregates));
    renderMatrix({ listMarketAggregates });

    const row = await screen.findByTestId("matrix-row-90001");
    // SKU = CODPROD (internal_product_id), never REFFORN — domain-model §7 FIX-1.
    expect(within(row).getByText("90001")).toBeInTheDocument();
    expect(within(row).getByText("Produto A")).toBeInTheDocument();
    // CUSTO + NOSSO PREÇO from the fact.
    expect(within(row).getByText(/120[.,]00/)).toBeInTheDocument();
    expect(within(row).getByText(/189[.,]00/)).toBeInTheDocument();
    // PREÇO MERCADO = median (169) once the single market call resolves — NO rank string.
    await waitFor(() =>
      expect(within(row).getByTestId("matrix-mercado-90001")).toHaveTextContent(/169[.,]00/),
    );
    // VEREDICTO price-evidence (resolves with the market call).
    expect(within(row).getByTestId("matrix-veredicto-90001")).toHaveTextContent("OK");
    // MARGEM: retorno/un + chip %.
    await waitFor(() => expect(within(row).getByText(/11[.,]01/)).toBeInTheDocument());
    expect(within(row).getByText("5,83%")).toBeInTheDocument();
  });

  it("(b) no market evidence: PREÇO MERCADO — , VEREDICTO SEM_EVIDENCIA, margem still computed", async () => {
    const listMarketAggregates = vi.fn(() => Promise.resolve(aggregates));
    renderMatrix({ listMarketAggregates });

    const row = await screen.findByTestId("matrix-row-90002");
    // Market price is the honest dash, never a fabricated number.
    expect(within(row).getByTestId("matrix-mercado-90002")).toHaveTextContent("—");
    // SEM_EVIDENCIA only once the market query RESOLVES to no-evidence (not while loading).
    await waitFor(() =>
      expect(within(row).getByTestId("matrix-veredicto-90002")).toHaveTextContent("SEM_EVIDENCIA"),
    );
    // Margem is still computed from custo/preço (decompose does not need market).
    await waitFor(() => expect(within(row).getByText(/11[.,]01/)).toBeInTheDocument());
  });

  it("(c) no ERP cost: CUSTO — and MARGEM UnknownValue —", async () => {
    const listMarketAggregates = vi.fn(() => Promise.resolve(aggregates));
    renderMatrix({ listMarketAggregates });

    const row = await screen.findByTestId("matrix-row-90003");
    expect(within(row).getByTestId("matrix-custo-90003")).toHaveTextContent("—");
    // Resolved-but-absent aggregate → honest SEM_EVIDENCIA (after the query resolves).
    await waitFor(() =>
      expect(within(row).getByTestId("matrix-veredicto-90003")).toHaveTextContent("SEM_EVIDENCIA"),
    );
    // Margem UnknownValue with the M-07 hint (categorical margin verdict not invented here).
    await waitFor(() => {
      const cell = within(row).getByTestId("matrix-margem-90003");
      expect(cell).toHaveTextContent("—");
      expect(within(cell).getByTitle("margem: M-07")).toBeInTheDocument();
    });
  });

  it('tags a product with no listing "novo" once the listing query resolves', async () => {
    const listMarketAggregates = vi.fn(() => Promise.resolve(aggregates));
    renderMatrix({ listMarketAggregates }, "inst_test");

    // 90003 has no listing → "novo"; 90001 has a listing → no tag.
    await waitFor(() => expect(screen.getByTestId("matrix-novo-90003")).toBeInTheDocument());
    expect(screen.queryByTestId("matrix-novo-90001")).toBeNull();
  });

  it("never tags novo when no installation is known (unknown ≠ novo)", async () => {
    const listMarketAggregates = vi.fn(() => Promise.resolve(aggregates));
    renderMatrix({ listMarketAggregates }); // installationId "" → listing lane disabled

    await screen.findByTestId("matrix-row-90003");
    expect(listListingsByProduct).not.toHaveBeenCalled();
    expect(screen.queryByTestId("matrix-novo-90003")).toBeNull();
  });

  it("queries the market ONCE with all codprods, and omits rank + comissao_pct", async () => {
    const listMarketAggregates = vi.fn(() => Promise.resolve(aggregates));
    renderMatrix({ listMarketAggregates });

    await screen.findByTestId("matrix-row-90001");
    // Single fan-in market call with every codprod.
    await waitFor(() => expect(listMarketAggregates).toHaveBeenCalledTimes(1));
    expect(listMarketAggregates).toHaveBeenCalledWith(["90001", "90002", "90003"]);

    // Rank "19º/23" is absent from aggregate data → never rendered.
    expect(screen.queryByText(/º\/\d/)).toBeNull();

    // Per-row decompose omits comissao_pct so the backend resolver chain runs.
    await waitFor(() => expect(pricingDecompose).toHaveBeenCalled());
    for (const [arg] of pricingDecompose.mock.calls) {
      expect(arg).not.toHaveProperty("comissao_pct");
    }
  });

  it("(insufficient market) VEREDICTO MERCADO_INSUFICIENTE, price stays —", async () => {
    const insufficient: MarketPriceIntelAggregate[] = [
      {
        product_id: "90001",
        median: null,
        min_valid: null,
        n_offers: 2,
        n_sellers: 2,
        source: "ml_catalog_offers",
        fetched_at: "2026-07-18T12:00:00Z",
        computed_at: "2026-07-18T12:00:00Z",
        status: "INSUFFICIENT_MARKET",
      },
    ];
    const listMarketAggregates = vi.fn(() => Promise.resolve(insufficient));
    renderMatrix({ listMarketAggregates });

    const row = await screen.findByTestId("matrix-row-90001");
    await waitFor(() =>
      expect(within(row).getByTestId("matrix-veredicto-90001")).toHaveTextContent(
        "MERCADO_INSUFICIENTE",
      ),
    );
    // No usable median → honest dash, never a fabricated price.
    expect(within(row).getByTestId("matrix-mercado-90001")).toHaveTextContent("—");
  });

  it("market fetch error stays unknown — never a confirmed SEM_EVIDENCIA (ADR-17)", async () => {
    const listMarketAggregates = vi.fn(() => Promise.reject(new Error("boom")));
    renderMatrix({ listMarketAggregates });

    const row = await screen.findByTestId("matrix-row-90001");
    await waitFor(() => expect(listMarketAggregates).toHaveBeenCalled());
    const veredicto = within(row).getByTestId("matrix-veredicto-90001");
    // A failed fetch is unknown, not a verdict — "we couldn't check" ≠ "no evidence".
    await waitFor(() => expect(veredicto).not.toHaveTextContent("SEM_EVIDENCIA"));
    expect(veredicto).toHaveTextContent("—");
    // Positively assert the ERROR hint (not the loading hint) so a regression that
    // sticks in the loading branch — same "—" DOM — cannot pass this test.
    expect(within(veredicto).getByTitle("mercado: falha ao carregar")).toBeInTheDocument();
    const mercado = within(row).getByTestId("matrix-mercado-90001");
    expect(mercado).toHaveTextContent("—");
    expect(within(mercado).getByTitle("mercado: falha ao carregar")).toBeInTheDocument();
  });

  it("decompose failure is honest, not misattributed to the M-07 placeholder (ADR-17)", async () => {
    const listMarketAggregates = vi.fn(() => Promise.resolve(aggregates));
    pricingDecompose.mockImplementation(() => Promise.reject(new Error("boom")));
    renderMatrix({ listMarketAggregates });

    const row = await screen.findByTestId("matrix-row-90001");
    const cell = within(row).getByTestId("matrix-margem-90001");
    await waitFor(() =>
      // Positively assert the ERROR hint so a regression stuck in the plain-loading
      // branch (also "—", no hint) cannot pass; and NOT the engine-pending M-07 hint.
      expect(within(cell).getByTitle("margem: falha ao calcular")).toBeInTheDocument(),
    );
    expect(within(cell).queryByTitle("margem: M-07")).toBeNull();
  });
});
