import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { fireEvent, render, screen, waitFor, within } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import { beforeEach, describe, expect, it, vi } from "vitest";
import type {
  CatalogProductFactPage,
  ListingPage,
  ListingReadModel,
} from "@marketplace-central/sdk-runtime";
import { MercadoPage } from "./MercadoPage";
import type { MercadoMarketClient } from "./marketClient";

const listListings = vi.fn();
const listCatalogProductFacts = vi.fn();
const getListingsSummary = vi.fn();

vi.mock("../../app/ClientContext", () => ({
  useClient: () => ({ listListings, listCatalogProductFacts, getListingsSummary }),
}));
vi.mock("../../app/InstallationContext", () => ({
  useInstallation: () => ({ installationId: "inst_1" }),
}));

const listing9001: ListingReadModel = {
  listing_id: "L1",
  installation_id: "inst_1",
  provider: "mercado_livre",
  provider_listing_id: "MLB-9001",
  variation_id: null,
  title: "Kit Parafuso M8 Premium",
  listing_type: null,
  status: "active",
  link: { state: "resolved", product_id: "90001", seller_sku: null },
  price: { amount: "54.90", currency: "BRL" },
  published_quantity: 10,
  sync_state: "synced",
  sync_error: null,
  quality_score: null,
  pending_issue: null,
  sales_30d: 63,
  cost: { amount: "40.00", currency: "BRL" },
  below_margin_worst_case: null,
  icms_worst_case_by_uf: null,
  fetched_at: "2026-07-19T06:00:00Z",
  market_signal: {
    status: "OK",
    position: { rank: 9, total: 14 },
    price_to_win: null,
    delta_pct: null,
    match_status: "ACCEPT",
    n_offers: 14,
    n_sellers: 10,
    median: { amount: "50.20", currency: "BRL" },
    min_valid: { amount: "47.50", currency: "BRL" },
    max_valid: null,
    evidence: null,
  },
};

const listingsPage: ListingPage = {
  items: [listing9001],
  next_cursor: null,
  page_size: 50,
  as_of: "2026-07-19T06:00:00Z",
};

const catalogPage: CatalogProductFactPage = {
  items: [
    {
      internal_product_id: 90008,
      reference: null,
      description: "PAPELEIRA DECA FLEX",
      ean: null,
      manufacturer_reference: null,
      brand_name: null,
      ncm: null,
      quality_flags: [],
      active: true,
      sellable_stock: { quantity: null, quality: [] },
      current_price: { amount: null, currency: "BRL", quality: [] },
      cost: { amount: "118.99", currency: "BRL", quality: [] },
    },
  ],
  next_cursor: null,
  page_size: 50,
  as_of: "2026-07-19T06:00:00Z",
};

function makeMarketClient(): MercadoMarketClient & { collectMarketPriceIntel: ReturnType<typeof vi.fn> } {
  return {
    collectMarketPriceIntel: vi.fn(),
    listMarketSignals: vi.fn(),
    listMarketAggregates: vi.fn().mockResolvedValue([
      {
        product_id: "90008",
        median: { amount: "229.20", currency: "BRL" },
        min_valid: { amount: "179.90", currency: "BRL" },
        n_offers: 6,
        n_sellers: 5,
        source: "ml_catalog_offers",
        fetched_at: "2026-07-19T06:00:00Z",
        computed_at: "2026-07-19T06:05:00Z",
        status: "OK",
      },
    ]),
    listMarketVerdicts: vi.fn().mockResolvedValue([
      {
        match_status: "ACCEPT",
        price_evidence_status: "OK",
        verdict_label: null,
        blocking_state: null,
        inputs_used: {},
        market_range: null,
      },
    ]),
  } as unknown as MercadoMarketClient & { collectMarketPriceIntel: ReturnType<typeof vi.fn> };
}

function renderPage(marketClient = makeMarketClient()) {
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  const utils = render(
    <QueryClientProvider client={queryClient}>
      {/* The reprice table links a listing's product into /precos, so the page
          needs a router context to render. */}
      <MemoryRouter>
        <MercadoPage marketClient={marketClient} />
      </MemoryRouter>
    </QueryClientProvider>,
  );
  return { ...utils, marketClient };
}

describe("MercadoPage", () => {
  beforeEach(() => {
    listListings.mockReset().mockResolvedValue(listingsPage);
    listCatalogProductFacts.mockReset().mockResolvedValue(catalogPage);
    getListingsSummary
      .mockReset()
      .mockResolvedValue({ total: 1, active: 1, paused: 0, exceptions: {}, as_of: "2026-07-19T06:00:00Z" });
  });

  it("renders the header: title, category filter, margin-min pill, source + refresh", async () => {
    renderPage();
    expect(screen.getByRole("heading", { name: "Mercado" })).toBeInTheDocument();
    expect(screen.getByText("Margem mín: 18%")).toBeInTheDocument();
    expect(screen.getByRole("button", { name: "Categoria ▾" })).toBeDisabled();
    expect(screen.getByRole("button", { name: "Atualizar agora" })).toBeDisabled();
    expect(screen.getByText(/fonte: busca pública ML/)).toBeInTheDocument();
    await screen.findByText("Kit Parafuso M8 Premium");
  });

  it("Reprecificação: renders real listing/signal fields and honest dashes for M-07 columns", async () => {
    renderPage();
    const cell = await screen.findByText("Kit Parafuso M8 Premium");
    const row = cell.closest("div")!;
    expect(within(row).getByText("R$ 54,90")).toBeInTheDocument(); // MEU PREÇO
    expect(within(row).getByText("R$ 47,50")).toBeInTheDocument(); // MENOR CONC.
    expect(within(row).getByText("R$ 50,20")).toBeInTheDocument(); // MEDIANA
    expect(within(row).getByText("9º/14")).toBeInTheDocument(); // POSIÇÃO
    expect(within(row).getByText("63")).toBeInTheDocument(); // VENDAS 30D
    // MARGEM ATUAL / SE IGUALAR MENOR / SUGESTÃO have no backing field → honest dash (ADR-17)
    expect(within(row).getAllByText("—").length).toBeGreaterThanOrEqual(3);
    // "Aplicar" writes to ML and stays inert (D-57); "Simular" only navigates,
    // so it is a real link carrying the listing's linked ERP product.
    expect(within(row).getByRole("button", { name: "Aplicar" })).toBeDisabled();
    expect(within(row).getByRole("link", { name: "Simular" })).toHaveAttribute(
      "href",
      "/precos?produto=90001",
    );
  });

  it("keeps Simular inert for a listing with no resolved product link", async () => {
    listListings.mockReset().mockResolvedValue({
      ...listingsPage,
      items: listingsPage.items.map((item) => ({
        ...item,
        link: { state: "unresolved", product_id: null, seller_sku: null },
      })),
    });
    renderPage();
    const cell = await screen.findByText("Kit Parafuso M8 Premium");
    const row = cell.closest("div")!;
    expect(within(row).queryByRole("link", { name: "Simular" })).not.toBeInTheDocument();
    expect(within(row).getByRole("button", { name: "Simular" })).toBeDisabled();
  });

  it("switches to Oportunidades and renders the 90008 opportunity from market evidence", async () => {
    const { marketClient } = renderPage();
    await screen.findByText("Kit Parafuso M8 Premium");
    fireEvent.click(screen.getByRole("tab", { name: /Oportunidades/ }));
    const cell = await screen.findByText("PAPELEIRA DECA FLEX");
    const row = cell.closest("div")!;
    expect(within(row).getByText("R$ 118,99")).toBeInTheDocument();
    expect(within(row).getByText("R$ 229,20")).toBeInTheDocument();
    // CONCORRENTES = distinct sellers (n_sellers=5), not offers (n_offers=6) — no overstatement
    expect(within(row).getByText("5")).toBeInTheDocument();
    // MARGEM EST. is M-07-owned → honest dash, never a client-fabricated % (ADR-17)
    expect(within(row).queryByText(/%$/)).toBeNull();
    expect(within(row).getAllByText("—").length).toBeGreaterThanOrEqual(3);
    expect(listCatalogProductFacts).toHaveBeenCalled();
    expect(marketClient.listMarketAggregates).toHaveBeenCalledWith(["90008"]);
  });

  it("Monitorados: honest empty state, inert filter chips and +Monitorar (no fabricated rows)", async () => {
    renderPage();
    fireEvent.click(screen.getByRole("tab", { name: "Monitorados" }));
    expect(await screen.findByText("Nenhum monitoramento ativo.")).toBeInTheDocument();
    expect(screen.getByRole("button", { name: "+ Monitorar…" })).toBeDisabled();
    expect(screen.getByRole("button", { name: "Todos" })).toBeDisabled();
  });

  it("collects market evidence for the codprods on the active tab and reports the causes honestly", async () => {
    const { marketClient } = renderPage();
    await screen.findByText("Kit Parafuso M8 Premium");
    marketClient.collectMarketPriceIntel.mockResolvedValue({
      status: "COMPLETED",
      "decisões": [],
      contagens: { ok: 1, no_price_evidence: 0, insufficient_market: 0, no_candidate: 0, sem_custo: 0 },
      causas: [],
    });

    fireEvent.click(screen.getByRole("button", { name: "Atualizar agora" }));

    await waitFor(() => expect(marketClient.collectMarketPriceIntel).toHaveBeenCalledWith("90001"));
    // The summary names what came back per cause — never a bare "pronto".
    expect(await screen.findByText(/1 com preço de mercado/)).toBeInTheDocument();
    expect(screen.getByText(/0 sem anúncio equivalente/)).toBeInTheDocument();
  });

  it("counts a failed collection instead of aborting the sweep", async () => {
    const { marketClient } = renderPage();
    await screen.findByText("Kit Parafuso M8 Premium");
    marketClient.collectMarketPriceIntel.mockRejectedValue({ status: 409, error: { code: "collection_in_progress" } });

    fireEvent.click(screen.getByRole("button", { name: "Atualizar agora" }));

    expect(await screen.findByText(/1 com falha na coleta/)).toBeInTheDocument();
  });

  it("performs no live ML write across tab interaction (D-57)", async () => {
    renderPage();
    const cell = await screen.findByText("Kit Parafuso M8 Premium");
    // Collection is a READ of the public ML search; "Aplicar" is the only
    // control that would write a price, and it stays inert.
    expect(within(cell.closest("div")!).getByRole("button", { name: "Aplicar" })).toBeDisabled();
    fireEvent.click(screen.getByRole("tab", { name: /Oportunidades/ }));
    await screen.findByText("PAPELEIRA DECA FLEX");
  });

  it("uses theme tokens (not hardcoded hex) so light/dark both resolve", async () => {
    renderPage();
    // The margin-min pill uses accent-soft/accent-ink tokens which flip per data-theme.
    expect(screen.getByText("Margem mín: 18%").className).toMatch(/bg-accent-soft/);
    expect(screen.getByText("Margem mín: 18%").className).toMatch(/text-accent-ink/);
    await screen.findByText("Kit Parafuso M8 Premium");
  });

  it("cursor-walks ALL listing pages (no silent page-1 truncation) and counts the honest total", async () => {
    const mk = (id: string, title: string): ListingReadModel => ({
      ...listing9001,
      listing_id: id,
      provider_listing_id: `MLB-${id}`,
      title,
    });
    const page1: ListingPage = {
      items: [mk("L1", "Anúncio Um"), mk("L2", "Anúncio Dois")],
      next_cursor: "cursor-2",
      page_size: 2,
      as_of: "2026-07-19T06:00:00Z",
    };
    const page2: ListingPage = {
      items: [mk("L3", "Anúncio Três"), mk("L4", "Anúncio Quatro")],
      next_cursor: null,
      page_size: 2,
      as_of: "2026-07-19T06:00:00Z",
    };
    listListings
      .mockReset()
      .mockImplementation((opts: { cursor?: string }) =>
        Promise.resolve(opts.cursor === "cursor-2" ? page2 : page1),
      );
    getListingsSummary.mockResolvedValue({ total: 4, active: 4, paused: 0, exceptions: {}, as_of: "2026-07-19T06:00:00Z" });

    renderPage();
    // Rows from BOTH pages must render — a single-page fetch would silently drop page 2.
    await screen.findByText("Anúncio Um");
    expect(screen.getByText("Anúncio Três")).toBeInTheDocument();
    expect(screen.getByText("Anúncio Quatro")).toBeInTheDocument();
    // The walk followed the cursor to page 2 (one call per page).
    expect(listListings).toHaveBeenCalledTimes(2);
    expect(listListings).toHaveBeenLastCalledWith(expect.objectContaining({ cursor: "cursor-2" }));
    // Tab count reflects the honest backed total (4), not the first-page size.
    expect(screen.getByRole("tab", { name: /Reprecificação/ }).textContent).toMatch(/\b4\b/);
  });

  it("surfaces an honest truncation banner when the listings walk hits the scan ceiling", async () => {
    // Server never terminates the cursor → the walk must stop at the ceiling and flag truncation,
    // never loop unbounded and never silently drop the uncovered tail (ADR-17).
    let page = 0;
    listListings.mockReset().mockImplementation(() => {
      page += 1;
      const item: ListingReadModel = { ...listing9001, listing_id: `L${page}`, provider_listing_id: `MLB-${page}` };
      return Promise.resolve({ items: [item], next_cursor: "more", page_size: 1, as_of: "2026-07-19T06:00:00Z" } satisfies ListingPage);
    });
    getListingsSummary.mockResolvedValue({ total: 5000, active: 5000, paused: 0, exceptions: {}, as_of: "2026-07-19T06:00:00Z" });

    renderPage();
    expect(await screen.findByText(/Varredura limitada aos primeiros/)).toBeInTheDocument();
    // Rows still render (populated branch) and the ceiling is honored (bounded, not unbounded).
    expect(screen.getAllByText("Kit Parafuso M8 Premium").length).toBeGreaterThan(0);
    expect(listListings).toHaveBeenCalledTimes(20); // MAX_PAGES ceiling
  });

  it("surfaces truncation even when every walked listing page is empty (empty state must not swallow the cap)", async () => {
    // Guards the asymmetry the refuter caught: listingsTruncated can be true while zero rows
    // rendered — the reprice empty-state branch must still show the honest banner, not a bare
    // "nothing here" that hides a capped scan.
    const emptyEndless: ListingPage = {
      items: [],
      next_cursor: "more",
      page_size: 1,
      as_of: "2026-07-19T06:00:00Z",
    };
    listListings.mockReset().mockResolvedValue(emptyEndless);
    getListingsSummary.mockResolvedValue({ total: 5000, active: 0, paused: 0, exceptions: {}, as_of: "2026-07-19T06:00:00Z" });

    renderPage();
    expect(await screen.findByText(/Varredura limitada aos primeiros/)).toBeInTheDocument();
    expect(screen.getByText("Nenhum anúncio com sinal de mercado ainda.")).toBeInTheDocument();
  });

  it("falls back to the walked count when the listings summary read fails (no crash, honest count)", async () => {
    // The honest total comes from getListingsSummary; if it fails the tab must still render
    // with the walked length — never error out, never present a partial page as complete.
    getListingsSummary.mockReset().mockRejectedValue(new Error("summary down"));

    renderPage();
    await screen.findByText("Kit Parafuso M8 Premium");
    expect(screen.getByRole("tab", { name: /Reprecificação/ }).textContent).toMatch(/\b1\b/);
  });

  it("cursor-walks ALL catalog pages for Oportunidades and batches every walked codprod to /market", async () => {
    const fact = (id: number, description: string) => ({
      ...catalogPage.items[0],
      internal_product_id: id,
      description,
    });
    const cat1: CatalogProductFactPage = {
      items: [fact(90008, "PAPELEIRA DECA FLEX")],
      next_cursor: "cat-2",
      page_size: 1,
      as_of: "2026-07-19T06:00:00Z",
    };
    const cat2: CatalogProductFactPage = {
      items: [fact(90009, "CUBA GRANITO")],
      next_cursor: null,
      page_size: 1,
      as_of: "2026-07-19T06:00:00Z",
    };
    listCatalogProductFacts
      .mockReset()
      .mockImplementation((opts: { cursor?: string }) =>
        Promise.resolve(opts.cursor === "cat-2" ? cat2 : cat1),
      );

    const { marketClient } = renderPage();
    fireEvent.click(screen.getByRole("tab", { name: /Oportunidades/ }));
    await screen.findByText("PAPELEIRA DECA FLEX");
    // Both catalog pages were walked...
    expect(listCatalogProductFacts).toHaveBeenCalledTimes(2);
    // ...and EVERY walked codprod (both pages, in order) went to the batch read — page-2
    // codprods would be missing if the walk stopped at page 1.
    expect(marketClient.listMarketAggregates).toHaveBeenCalledWith(["90008", "90009"]);
    expect(marketClient.listMarketVerdicts).toHaveBeenCalledWith(["90008", "90009"]);
  });

  it("surfaces the truncation banner on the POPULATED Oportunidades tab when the catalog walk hits the ceiling", async () => {
    // Endless catalog cursor → the walk stops at MAX_PAGES with truncated=true; every walked
    // codprod has an OK aggregate so rows render — the banner must accompany them (ADR-17).
    let page = 0;
    listCatalogProductFacts.mockReset().mockImplementation(() => {
      page += 1;
      return Promise.resolve({
        items: [{ ...catalogPage.items[0], internal_product_id: 90000 + page }],
        next_cursor: "more",
        page_size: 1,
        as_of: "2026-07-19T06:00:00Z",
      } satisfies CatalogProductFactPage);
    });
    const okForAll = {
      collectMarketPriceIntel: vi.fn(),
      listMarketSignals: vi.fn(),
      listMarketAggregates: vi.fn().mockImplementation((ids: string[]) =>
        Promise.resolve(
          ids.map((id) => ({
            product_id: id,
            median: { amount: "229.20", currency: "BRL" },
            min_valid: { amount: "179.90", currency: "BRL" },
            n_offers: 6,
            n_sellers: 5,
            source: "ml_catalog_offers",
            fetched_at: "2026-07-19T06:00:00Z",
            computed_at: "2026-07-19T06:05:00Z",
            status: "OK",
          })),
        ),
      ),
      listMarketVerdicts: vi.fn().mockImplementation((ids: string[]) =>
        Promise.resolve(
          ids.map(() => ({
            match_status: "ACCEPT",
            price_evidence_status: "OK",
            verdict_label: null,
            blocking_state: null,
            inputs_used: {},
            market_range: null,
          })),
        ),
      ),
    } as unknown as MercadoMarketClient & { collectMarketPriceIntel: ReturnType<typeof vi.fn> };

    renderPage(okForAll);
    fireEvent.click(screen.getByRole("tab", { name: /Oportunidades/ }));
    expect(await screen.findByText(/Varredura limitada aos primeiros/)).toBeInTheDocument();
    expect(listCatalogProductFacts).toHaveBeenCalledTimes(20); // MAX_PAGES ceiling
  });

  it("surfaces the truncation banner on the EMPTY Oportunidades state when the catalog walk hits the ceiling", async () => {
    // Endless empty catalog pages → truncated=true with zero rows; the opp empty-state branch
    // must still show the honest banner, never a bare "nothing here" that hides a capped scan.
    listCatalogProductFacts.mockReset().mockResolvedValue({
      items: [],
      next_cursor: "more",
      page_size: 1,
      as_of: "2026-07-19T06:00:00Z",
    } satisfies CatalogProductFactPage);

    renderPage();
    fireEvent.click(screen.getByRole("tab", { name: /Oportunidades/ }));
    expect(await screen.findByText(/Varredura limitada aos primeiros/)).toBeInTheDocument();
    expect(screen.getByText("Nenhum produto do catálogo com demanda de mercado observada.")).toBeInTheDocument();
  });
});
