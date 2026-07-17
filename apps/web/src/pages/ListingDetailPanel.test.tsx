import type { ListingDetail } from "@marketplace-central/sdk-runtime";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";
import { ListingDetailPanel } from "./ListingDetailPanel";

const getListing = vi.fn();
const listListings = vi.fn();

const detail: ListingDetail = {
  listing_id: "listing_1",
  installation_id: "installation_1",
  provider: "mercado_livre",
  provider_listing_id: "MLB123456789",
  title: "Camiseta azul",
  listing_type: { code: "gold_special", label: "Clássico" },
  status: "active",
  link: { state: "resolved", product_id: "product_1", seller_sku: "CAM-AZ" },
  price: { amount: "129.90", currency: "BRL" },
  published_quantity: 7,
  sync_state: "error",
  sync_error: {
    code: "provider_rejected",
    message_pt: "O anúncio foi rejeitado pelo provedor.",
    message_provider: "The provider rejected the listing.",
  },
  quality_score: null,
  pending_issue: null,
  sales_30d: null,
  cost: null,
  below_margin_worst_case: null,
  icms_worst_case_by_uf: null,
  fetched_at: "2026-07-16T12:00:00Z",
  timeline: [
    { at: "2026-07-15T12:00:00Z", kind: "synced", message_pt: "Evento A" },
    { at: "2026-07-16T12:00:00Z", kind: "sync_error", message_pt: "Evento B" },
    { at: "2026-07-14T12:00:00Z", kind: "created", message_pt: "Evento C" },
  ],
};

function renderPanel(listingId = "listing_1", onClose = vi.fn()) {
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  return {
    onClose,
    ...render(
      <QueryClientProvider client={queryClient}>
        <ListingDetailPanel listingId={listingId} onClose={onClose} />
      </QueryClientProvider>,
    ),
  };
}

vi.mock("../app/ClientContext", () => ({
  useClient: () => ({ getListing, listListings }),
}));

describe("ListingDetailPanel", () => {
  it("loads one detail without refetching the listing page", async () => {
    getListing.mockResolvedValueOnce(detail);

    renderPanel();

    await screen.findByText("Camiseta azul");
    expect(getListing).toHaveBeenCalledOnce();
    expect(getListing).toHaveBeenCalledWith("listing_1");
    expect(listListings).not.toHaveBeenCalled();
  });

  it("keeps timeline events in the API returned order without re-sorting", async () => {
    getListing.mockResolvedValueOnce(detail);

    renderPanel();

    await screen.findByText("Evento A");
    const rendered = screen
      .getAllByText(/^Evento [ABC]$/)
      .map((node) => node.textContent);
    expect(rendered).toEqual(["Evento A", "Evento B", "Evento C"]);
  });

  it("shows translated sync errors and discloses provider text only on demand", async () => {
    getListing.mockResolvedValueOnce(detail);

    renderPanel();

    expect(await screen.findByText(detail.sync_error!.message_pt)).toBeInTheDocument();
    expect(screen.queryByText(detail.sync_error!.message_provider!)).not.toBeInTheDocument();
    expect(screen.getByText("▸ técnico")).toBeInTheDocument();

    fireEvent.click(screen.getByText("▸ técnico"));
    expect(screen.getByText(detail.sync_error!.message_provider!)).toBeInTheDocument();
  });

  it("renders nullable price, quality, and sales as unknown", async () => {
    getListing.mockResolvedValueOnce(detail);

    renderPanel();

    await screen.findByText("Camiseta azul");
    expect(screen.getAllByText("—").length).toBeGreaterThanOrEqual(3);
  });

  it("keeps future actions disabled and side-effect free", async () => {
    getListing.mockResolvedValueOnce(detail);

    renderPanel();
    await screen.findByText("Camiseta azul");
    const callsBefore = getListing.mock.calls.length;

    for (const name of ["Corrigir", "Simular", "Pausar"]) {
      const button = screen.getByRole("button", { name });
      expect(button).toBeDisabled();
      expect(button).toHaveAttribute("title", "disponível em breve");
      fireEvent.click(button);
    }

    expect(getListing).toHaveBeenCalledTimes(callsBefore);
  });

  it("closes with the Portuguese accessible label", async () => {
    getListing.mockResolvedValueOnce(detail);
    const onClose = vi.fn();

    renderPanel("listing_1", onClose);
    await screen.findByText("Camiseta azul");

    fireEvent.click(screen.getByRole("button", { name: "Fechar painel" }));
    expect(onClose).toHaveBeenCalledOnce();
  });
});
