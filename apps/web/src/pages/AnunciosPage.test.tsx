import { MarketplaceCentralClientError, type ListingSummary } from "@marketplace-central/sdk-runtime";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import { MemoryRouter, useLocation } from "react-router-dom";
import { describe, expect, it, vi } from "vitest";
import { AnunciosPage } from "./AnunciosPage";

const listListings = vi.fn();
const listListingsByProduct = vi.fn();
const getListingsSummary = vi.fn(
  (): Promise<ListingSummary> =>
    Promise.resolve({
      total: 1,
      active: 1,
      paused: 0,
      exceptions: {
        sync_error: 0,
        stale: 0,
        unlinked: 0,
        below_margin_worst_case: null,
        margin_unknown: null,
      },
      as_of: "2026-07-16T12:00:00Z",
    }),
);

const listingPage = {
  items: [{
    listing_id: "listing_1",
    installation_id: "inst_1",
    provider: "mercado_livre",
    provider_listing_id: "MLB123456789",
    variation_id: null,
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
  useClient: () => ({
    listListings: (...args: unknown[]) => listListings(...args),
    listListingsByProduct: (...args: unknown[]) => listListingsByProduct(...args),
    getListingsSummary,
  }),
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
        <LocationProbe />
      </MemoryRouter>
    </QueryClientProvider>,
  );
}

function LocationProbe() {
  return <output data-testid="location-search">{useLocation().search}</output>;
}

describe("AnunciosPage active filter chips", () => {
  it("renders and dismisses an exception chip while preserving other query state", async () => {
    listListings.mockReset();
    listListings.mockResolvedValue(listingPage);
    renderPage("?installation=inst_1&tab=ativos&filter.exception=sync_error&q=camiseta");

    expect(await screen.findByText("Exceção: Erro de sync")).toBeInTheDocument();
    fireEvent.click(screen.getByRole("button", { name: "Remover filtro Erro de sync" }));

    await waitFor(() => {
      expect(screen.queryByText("Exceção: Erro de sync")).not.toBeInTheDocument();
    });
    expect(listListings.mock.calls.at(-1)?.[0]).toEqual({
      installation_id: "inst_1",
      q: "camiseta",
      status: "active",
    });
    const search = new URLSearchParams(screen.getByTestId("location-search").textContent ?? "");
    expect(search.get("installation")).toBe("inst_1");
    expect(search.get("tab")).toBe("ativos");
    expect(search.get("q")).toBe("camiseta");
    expect(search.has("filter.exception")).toBe(false);
  });

  it("dismisses a sync_state chip without touching the remaining filters", async () => {
    listListings.mockReset();
    listListings.mockResolvedValue(listingPage);
    renderPage("?installation=inst_1&filter.sync_state=error&filter.link_state=conflict");

    expect(await screen.findByText("Sync: com erro")).toBeInTheDocument();
    expect(screen.getByText("Vínculo: divergente")).toBeInTheDocument();
    fireEvent.click(screen.getByRole("button", { name: "Remover filtro com erro" }));

    await waitFor(() => {
      expect(screen.queryByText("Sync: com erro")).not.toBeInTheDocument();
    });
    expect(screen.getByText("Vínculo: divergente")).toBeInTheDocument();
    expect(listListings.mock.calls.at(-1)?.[0]).toEqual({
      installation_id: "inst_1",
      link_state: "conflict",
    });
    const search = new URLSearchParams(screen.getByTestId("location-search").textContent ?? "");
    expect(search.has("filter.sync_state")).toBe(false);
    expect(search.get("filter.link_state")).toBe("conflict");
  });

  it("renders no filter chips without active filters", async () => {
    listListings.mockReset();
    listListings.mockResolvedValue(listingPage);
    renderPage("?installation=inst_1");

    expect(await screen.findByText("Camiseta azul")).toBeInTheDocument();
    expect(screen.queryByRole("button", { name: /Remover filtro/ })).not.toBeInTheDocument();
  });
});

describe("AnunciosPage exception summary chips", () => {
  it("clicking a summary exception chip deep-links the URL filter and re-issues the list query", async () => {
    listListings.mockReset();
    listListings.mockResolvedValue(listingPage);
    getListingsSummary.mockReset();
    getListingsSummary.mockResolvedValue({
      total: 10,
      active: 6,
      paused: 4,
      exceptions: {
        sync_error: 0,
        stale: 0,
        unlinked: 0,
        below_margin_worst_case: null,
        margin_unknown: null,
        sem_vinculo: 4,
      },
      as_of: "2026-07-16T12:00:00Z",
    });
    renderPage("?installation=inst_1");

    const chip = await screen.findByRole("button", { name: "sem vínculo 4" });
    fireEvent.click(chip);

    await waitFor(() => {
      expect(listListings.mock.calls.at(-1)?.[0]).toEqual({
        installation_id: "inst_1",
        exception: "sem_vinculo",
      });
    });
    const search = new URLSearchParams(screen.getByTestId("location-search").textContent ?? "");
    expect(search.get("filter.exception")).toBe("sem_vinculo");
    expect(await screen.findByText("Exceção: Sem vínculo")).toBeInTheDocument();
  });

  it("renders the header sem-vínculo chip from the real unlinked count and filters the table (EXEMPLO-IO)", async () => {
    listListings.mockReset();
    listListings.mockResolvedValue(listingPage);
    getListingsSummary.mockReset();
    getListingsSummary.mockResolvedValue({
      total: 1284,
      active: 1200,
      paused: 40,
      exceptions: {
        sync_error: 0,
        stale: 0,
        unlinked: 3,
        below_margin_worst_case: null,
        margin_unknown: null,
      },
      as_of: "2026-07-16T12:00:00Z",
    });
    renderPage("?installation=inst_1");

    // Honest count from the real `unlinked` field (no sem_vinculo alias present).
    const chip = await screen.findByRole("button", { name: "sem vínculo 3" });
    fireEvent.click(chip);

    await waitFor(() => {
      expect(listListings.mock.calls.at(-1)?.[0]).toEqual({
        installation_id: "inst_1",
        exception: "sem_vinculo",
      });
    });
    expect(await screen.findByText("Exceção: Sem vínculo")).toBeInTheDocument();
  });

  it("restores the exception filter from the URL on load (F5-restorable)", async () => {
    listListings.mockReset();
    listListings.mockResolvedValue(listingPage);
    getListingsSummary.mockReset();
    getListingsSummary.mockResolvedValue({
      total: 10,
      active: 6,
      paused: 4,
      exceptions: {
        sync_error: 0,
        stale: 0,
        unlinked: 0,
        below_margin_worst_case: null,
        margin_unknown: null,
        sem_vinculo: 4,
      },
      as_of: "2026-07-16T12:00:00Z",
    });
    renderPage("?installation=inst_1&filter.exception=sem_vinculo");

    expect(await screen.findByRole("button", { name: "sem vínculo 4" })).toHaveAttribute(
      "aria-pressed",
      "true",
    );
    expect(listListings.mock.calls.at(-1)?.[0]).toEqual({
      installation_id: "inst_1",
      exception: "sem_vinculo",
    });
  });
});

describe("AnunciosPage error recovery", () => {
  it("clears filter params on retry only for invalid_filter errors", async () => {
    listListings.mockRejectedValueOnce(
      new MarketplaceCentralClientError(400, "invalid_filter", "invalid", {}),
    );
    listListings.mockResolvedValue({ items: [] });
    renderPage("?installation=inst_1&filter.exception=sync_error&q=abc");

    const [retry] = await screen.findAllByRole("button", { name: "Tentar novamente" });
    fireEvent.click(retry);

    expect(await screen.findByText("Nenhum registro encontrado.")).toBeInTheDocument();
    const lastOptions = listListings.mock.calls.at(-1)?.[0] as Record<string, unknown>;
    expect(lastOptions).toEqual({ installation_id: "inst_1", q: "abc" });
  });

  it("retries with filters intact when the error is not invalid_filter", async () => {
    listListings.mockRejectedValueOnce(
      new MarketplaceCentralClientError(500, "internal_error", "boom", {}),
    );
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
    expect(screen.queryByRole("button", { name: /Sem vínculo/ })).not.toBeInTheDocument();
    expect(screen.queryByRole("button", { name: /Abaixo do custo/ })).not.toBeInTheDocument();
    expect(screen.queryByRole("button", { name: /Sem evidência/ })).not.toBeInTheDocument();
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

describe("AnunciosPage agrupar por produto toggle", () => {
  const groupPage = {
    groups: [
      {
        product_id: "product_1",
        product_title: "Grupo Camisetas",
        listing_count: 1,
        group_state: "ok",
        listings: [listingPage.items[0]],
      },
      {
        product_id: "product_2",
        product_title: "Grupo Meias",
        listing_count: 1,
        group_state: "ok",
        listings: [{ ...listingPage.items[0], listing_id: "listing_2", title: "Meia branca" }],
      },
    ],
    next_cursor: null,
    page_size: 2,
    as_of: "2026-07-16T12:00:00Z",
  };

  it("toggling on issues listListingsByProduct, sets the URL flag, and renders per-product groups", async () => {
    listListings.mockReset();
    listListings.mockResolvedValue(listingPage);
    listListingsByProduct.mockReset();
    listListingsByProduct.mockResolvedValue(groupPage);
    renderPage("?installation=inst_1");

    expect(await screen.findByText("Camiseta azul")).toBeInTheDocument();
    expect(listListingsByProduct).not.toHaveBeenCalled();

    fireEvent.click(screen.getByRole("checkbox", { name: "Agrupar por produto" }));

    await waitFor(() => {
      expect(listListingsByProduct).toHaveBeenCalledWith({ installation_id: "inst_1" });
    });
    expect(await screen.findByText("Grupo Meias")).toBeInTheDocument();
    expect(await screen.findByText("Meia branca")).toBeInTheDocument();
    const search = new URLSearchParams(screen.getByTestId("location-search").textContent ?? "");
    expect(search.get("grouped")).toBe("1");
  });

  it("renders a 1-listing group as a normal group (no special collapse)", async () => {
    listListings.mockReset();
    listListingsByProduct.mockReset();
    listListingsByProduct.mockResolvedValue({
      groups: [groupPage.groups[0]],
      next_cursor: null,
      page_size: 1,
      as_of: "2026-07-16T12:00:00Z",
    });
    renderPage("?installation=inst_1&grouped=1");

    expect(await screen.findByText("Camiseta azul")).toBeInTheDocument();
    expect(screen.getByRole("checkbox", { name: "Agrupar por produto" })).toBeChecked();
  });

  it("toggling off restores the flat W1 list unchanged", async () => {
    listListings.mockReset();
    listListings.mockResolvedValue(listingPage);
    listListingsByProduct.mockReset();
    listListingsByProduct.mockResolvedValue(groupPage);
    renderPage("?installation=inst_1&grouped=1");

    expect(await screen.findByRole("checkbox", { name: "Agrupar por produto" })).toBeChecked();
    fireEvent.click(screen.getByRole("checkbox", { name: "Agrupar por produto" }));

    await waitFor(() => {
      expect(listListings).toHaveBeenCalledWith({ installation_id: "inst_1" });
    });
    const search = new URLSearchParams(screen.getByTestId("location-search").textContent ?? "");
    expect(search.has("grouped")).toBe(false);
  });
});
