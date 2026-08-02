import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { render, screen } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import { describe, expect, it, vi } from "vitest";
import type { ListingGroupPage, ListingReadModel } from "@marketplace-central/sdk-runtime";
import { AnunciosVinculadosTab } from "./AnunciosVinculadosTab";
import type { Client } from "../../app/ClientContext";

function makeClient(overrides: Partial<Client> = {}): Client {
  return {
    listListingsByProduct: vi.fn(),
    ...overrides,
  } as unknown as Client;
}

function renderTab(client: Client, productId = "90008") {
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  return render(
    <QueryClientProvider client={queryClient}>
      <MemoryRouter>
        <AnunciosVinculadosTab productId={productId} client={client} installationId="inst_1" />
      </MemoryRouter>
    </QueryClientProvider>,
  );
}

const fullSignalListing: ListingReadModel = {
  listing_id: "l1",
  installation_id: "inst_1",
  provider: "mercadolivre",
  provider_listing_id: "MLB1",
  variation_id: null,
  title: "Anúncio com sinal",
  listing_type: { code: "gold_special", label: "Clássico" },
  status: "active",
  link: { state: "resolved", product_id: "90008", seller_sku: null },
  price: { amount: "169.90", currency: "BRL" },
  published_quantity: 10,
  sync_state: "synced",
  sync_error: null,
  pending_issue: null,
  cost: null,
  below_margin_worst_case: null,
  icms_worst_case_by_uf: null,
  fetched_at: null,
  market_signal: {
    status: "OK",
    position: { rank: 2, total: 9 },
    price_to_win: { amount: "169.90", currency: "BRL" },
    delta_pct: "-3.2",
    match_status: "ACCEPT",
    n_offers: 9,
    n_sellers: 9,
    median: null,
    min_valid: null,
    max_valid: null,
    evidence: { source: "ml", fetched_at: "2026-07-19T10:00:00Z", freshness: "fresh" },
  },
};

const noSignalListing: ListingReadModel = {
  ...fullSignalListing,
  listing_id: "l2",
  provider_listing_id: "MLB2",
  variation_id: null,
  title: "Anúncio sem vínculo",
  market_signal: { ...fullSignalListing.market_signal!, status: "SEM_VINCULO" },
};

describe("AnunciosVinculadosTab", () => {
  it("renders full market signal for one row and honest dashes for a SEM_VINCULO row", async () => {
    const page: ListingGroupPage = {
      groups: [
        {
          product_id: "90008",
          product_title: "Produto 90008",
          listing_count: 2,
          group_state: "ok",
          listings: [fullSignalListing, noSignalListing],
        },
      ],
      next_cursor: null,
      page_size: 2,
      as_of: "2026-07-19T10:00:00Z",
    };
    const client = makeClient({ listListingsByProduct: vi.fn().mockResolvedValue(page) });

    renderTab(client);

    expect(await screen.findByText("Anúncio com sinal")).toBeInTheDocument();
    expect(screen.getByText("2/9")).toBeInTheDocument();
    // pt-BR through the shared formatter; the regex tolerates the nbsp Intl puts after "R$".
    expect(screen.getByText(/^R\$\s169,90$/)).toBeInTheDocument();
    expect(screen.getByText("-3,2%")).toBeInTheDocument();

    const secondRow = screen.getByText("Anúncio sem vínculo").closest("tr");
    expect(secondRow).not.toBeNull();
    expect(secondRow!.textContent).toContain("—");
    expect(secondRow!.textContent).not.toContain("169.90");
  });

  it("renders EmptyState with a /vinculos link when there are no listings", async () => {
    const page: ListingGroupPage = {
      groups: [],
      next_cursor: null,
      page_size: 0,
      as_of: "2026-07-19T10:00:00Z",
    };
    const client = makeClient({ listListingsByProduct: vi.fn().mockResolvedValue(page) });

    renderTab(client);

    expect(await screen.findByText("Nenhum registro encontrado.")).toBeInTheDocument();
    expect(screen.getByRole("link", { name: "vincular anúncios" })).toHaveAttribute("href", "/vinculos");
  });

  it("does not leak a mismatched group's listings when the lone group belongs to another product", async () => {
    const page: ListingGroupPage = {
      groups: [
        {
          product_id: "12345",
          product_title: "Outro produto",
          listing_count: 1,
          group_state: "ok",
          listings: [{ ...fullSignalListing, listing_id: "l-other" }],
        },
      ],
      next_cursor: null,
      page_size: 1,
      as_of: "2026-07-19T10:00:00Z",
    };
    const client = makeClient({ listListingsByProduct: vi.fn().mockResolvedValue(page) });

    renderTab(client, "90008");

    expect(await screen.findByText("Nenhum registro encontrado.")).toBeInTheDocument();
    expect(screen.queryByText("Anúncio com sinal")).not.toBeInTheDocument();
  });

  it("scopes listListingsByProduct to the given productId", async () => {
    const listListingsByProduct = vi.fn().mockResolvedValue({ groups: [], next_cursor: null, page_size: 0, as_of: "" });
    const client = makeClient({ listListingsByProduct });

    renderTab(client, "90099");

    await screen.findByText("Nenhum registro encontrado.");
    expect(listListingsByProduct).toHaveBeenCalledWith(
      expect.objectContaining({ product_id: "90099", installation_id: "inst_1" }),
    );
  });
});
