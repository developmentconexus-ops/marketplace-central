import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import { describe, expect, it, vi } from "vitest";
import { AnunciosPage } from "./AnunciosPage";

const listListings = vi.fn();
const getListingsSummary = vi.fn(() =>
  Promise.resolve({
    total: 1,
    active: 1,
    paused: 0,
    exceptions: { below_margin_worst_case: null, margin_unknown: null },
  }),
);

const listingPage = {
  items: [{
    listing_id: "listing_1",
    installation_id: "inst_1",
    provider: "mercado_livre",
    provider_listing_id: "MLB123456789",
    title: "Camiseta azul",
    listing_type: { code: "gold_special", label: "Clássico" },
    status: "active",
    link: { state: "resolved", product_id: "product_1", seller_sku: "CAM-AZ" },
    price: { amount: "129.90", currency: "BRL" },
    published_quantity: 7,
    sync_state: "synced",
    sync_error: null,
    quality_score: 0.9,
    pending_issue: null,
    sales_30d: 12,
    cost: { amount: "70.00", currency: "BRL" },
    below_margin_worst_case: false,
    icms_worst_case_by_uf: null,
    fetched_at: "2026-07-16T12:00:00Z",
  }],
  next_cursor: null,
  page_size: 1,
  as_of: "2026-07-16T12:00:00Z",
};

vi.mock("../app/ClientContext", () => ({
  useClient: () => ({ listListings: (...args: unknown[]) => listListings(...args), getListingsSummary }),
}));

vi.mock("../app/InstallationContext", () => ({
  useInstallation: () => ({ installationId: "inst_1" }),
}));

function renderPage(search: string) {
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  return render(
    <QueryClientProvider client={queryClient}>
      <MemoryRouter initialEntries={[`/anuncios${search}`]}>
        <AnunciosPage />
      </MemoryRouter>
    </QueryClientProvider>,
  );
}

describe("AnunciosPage error recovery", () => {
  it("clears filter params on retry only for invalid_filter errors", async () => {
    listListings.mockRejectedValueOnce({
      status: 400,
      error: { code: "invalid_filter", message: "invalid" },
    });
    listListings.mockResolvedValue({ items: [] });
    renderPage("?installation=inst_1&filter.exception=sync_error&q=abc");

    const [retry] = await screen.findAllByRole("button", { name: "Tentar novamente" });
    fireEvent.click(retry);

    expect(await screen.findByText("Nenhum registro encontrado.")).toBeInTheDocument();
    const lastOptions = listListings.mock.calls.at(-1)?.[0] as Record<string, unknown>;
    expect(lastOptions).toEqual({ installation_id: "inst_1", q: "abc" });
  });

  it("retries with filters intact when the error is not invalid_filter", async () => {
    listListings.mockRejectedValueOnce({
      status: 500,
      error: { code: "internal", message: "boom" },
    });
    listListings.mockResolvedValue({ items: [] });
    renderPage("?installation=inst_1&filter.exception=sync_error");

    const [retry] = await screen.findAllByRole("button", { name: "Tentar novamente" });
    fireEvent.click(retry);

    expect(await screen.findByText("Nenhum registro encontrado.")).toBeInTheDocument();
    const lastOptions = listListings.mock.calls.at(-1)?.[0] as Record<string, unknown>;
    expect(lastOptions).toEqual({ installation_id: "inst_1", exception: "sync_error" });
  });
});

describe("AnunciosPage independent states", () => {
  it("keeps table rows visible when the summary query fails", async () => {
    listListings.mockResolvedValueOnce(listingPage);
    getListingsSummary.mockRejectedValueOnce({ status: 500, error: { code: "internal" } });
    renderPage("?installation=inst_1");

    expect(await screen.findByText("Camiseta azul")).toBeInTheDocument();
    expect(screen.getByRole("alert")).toHaveTextContent("Erro ao carregar.");
  });

  it("keeps summary counters visible when the table query fails", async () => {
    listListings.mockRejectedValueOnce({ status: 500, error: { code: "internal" } });
    getListingsSummary.mockResolvedValueOnce({
      total: 10,
      active: 6,
      paused: 4,
      exceptions: { sync_error: 1, stale: 2, unlinked: 3, below_margin_worst_case: null, margin_unknown: null },
      as_of: "2026-07-16T12:00:00Z",
    });
    renderPage("?installation=inst_1");

    expect(await screen.findByText("10")).toBeInTheDocument();
    expect(screen.getByRole("alert")).toHaveTextContent("Erro ao carregar.");
  });

  it("offers to clear filters after a successful empty response", async () => {
    listListings.mockResolvedValueOnce({ items: [], next_cursor: null, page_size: 0, as_of: "2026-07-16T12:00:00Z" });
    renderPage("?installation=inst_1&filter.exception=sync_error");

    expect(await screen.findByText("Nenhum registro encontrado.")).toBeInTheDocument();
    fireEvent.click(screen.getByRole("button", { name: "Limpar filtros" }));

    await waitFor(() => {
      expect(listListings.mock.calls.at(-1)?.[0]).toEqual({ installation_id: "inst_1" });
    });
  });
});
