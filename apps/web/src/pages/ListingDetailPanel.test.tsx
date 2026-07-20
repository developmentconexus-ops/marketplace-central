import type { ListingDetail, ListingMarketSignal } from "@marketplace-central/sdk-runtime";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";
import { MemoryRouter } from "react-router-dom";
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

const okSignal: ListingMarketSignal = {
  status: "OK",
  position: { rank: 2, total: 5 },
  price_to_win: { amount: "119.90", currency: "BRL" },
  delta_pct: "-7.5",
  match_status: "ACCEPT",
  n_offers: 5,
  n_sellers: 3,
  median: { amount: "135.50", currency: "BRL" },
  min_valid: { amount: "110.00", currency: "BRL" },
  max_valid: { amount: "149.90", currency: "BRL" },
  evidence: { source: "buybox_api", fetched_at: "2026-07-18T10:00:00Z", freshness: "fresh" },
};

function renderPanel(listingId = "listing_1", onClose = vi.fn()) {
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  return {
    onClose,
    ...render(
      <MemoryRouter>
        <QueryClientProvider client={queryClient}>
          <ListingDetailPanel listingId={listingId} onClose={onClose} />
        </QueryClientProvider>
      </MemoryRouter>,
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

    await screen.findAllByText("Camiseta azul");
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

    await screen.findAllByText("Camiseta azul");
    expect(screen.getAllByText("—").length).toBeGreaterThanOrEqual(3);
  });

  it("keeps future actions disabled and side-effect free", async () => {
    getListing.mockResolvedValueOnce(detail);

    renderPanel();
    await screen.findAllByText("Camiseta azul");
    const callsBefore = getListing.mock.calls.length;

    for (const name of ["Corrigir anúncio", "Simular preço", "Pausar"]) {
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
    await screen.findAllByText("Camiseta azul");

    fireEvent.click(screen.getByRole("button", { name: "Fechar painel" }));
    expect(onClose).toHaveBeenCalledOnce();
  });

  it("renders full market evidence for OK signal_status", async () => {
    getListing.mockResolvedValueOnce({ ...detail, signal_status: "OK", market_signal: okSignal });

    renderPanel();
    await screen.findAllByText("Camiseta azul");

    expect(screen.getByText("2/5")).toBeInTheDocument();
    expect(screen.getByText("-7.5%")).toBeInTheDocument();
    expect(screen.getByText("ACCEPT")).toBeInTheDocument();
    expect(screen.getByText("5")).toBeInTheDocument();
    expect(screen.getByText("3")).toBeInTheDocument();
    expect(screen.getByText("buybox_api")).toBeInTheDocument();
    expect(screen.getByText(/R\$ 119\.90/)).toBeInTheDocument();
    // Honest label rename: the ML target is "Alvo buybox (ML)", not the
    // operator-misleading "Preço p/ vencer".
    expect(screen.getByText("Alvo buybox (ML)")).toBeInTheDocument();
    expect(screen.queryByText("Preço p/ vencer")).not.toBeInTheDocument();
    // Faixa de mercado (min — mediana — máx), competitor-only.
    expect(screen.getByText("Faixa de mercado (concorrentes)")).toBeInTheDocument();
    expect(screen.getByText(/R\$ 110\.00/)).toBeInTheDocument();
    expect(screen.getByText(/R\$ 135\.50/)).toBeInTheDocument();
    expect(screen.getByText(/R\$ 149\.90/)).toBeInTheDocument();
    // Own-seller exclusion means the seller count is competitor-only.
    expect(screen.getByText("Concorrentes")).toBeInTheDocument();
    expect(
      screen.getByRole("link", { name: "Ver produto vinculado" }),
    ).toHaveAttribute("href", "/catalogo/produtos/product_1");
  });

  it("renders honest '—' for absent faixa bounds, never a fabricated 0 (ADR-17)", async () => {
    getListing.mockResolvedValueOnce({
      ...detail,
      signal_status: "OK",
      market_signal: { ...okSignal, min_valid: null, median: null, max_valid: null },
    });

    renderPanel();
    await screen.findAllByText("Camiseta azul");

    // The faixa card is still present and labeled, but every bound reads "—".
    expect(screen.getByText("Faixa de mercado (concorrentes)")).toBeInTheDocument();
    expect(screen.queryByText(/R\$ 110\.00/)).not.toBeInTheDocument();
    expect(screen.queryByText(/R\$ 135\.50/)).not.toBeInTheDocument();
    expect(screen.queryByText(/R\$ 149\.90/)).not.toBeInTheDocument();
  });

  it("shows STALE evidence numbers marked with a freshness indicator, not silently hidden", async () => {
    getListing.mockResolvedValueOnce({
      ...detail,
      signal_status: "STALE",
      market_signal: { ...okSignal, status: "STALE" },
    });

    renderPanel();
    await screen.findAllByText("Camiseta azul");

    expect(screen.getByText("2/5")).toBeInTheDocument();
    expect(screen.getByText(/desatualizada/i)).toBeInTheDocument();
  });

  it("shows an honest sem-evidência state for NO_PRICE_EVIDENCE, never a fabricated number", async () => {
    getListing.mockResolvedValueOnce({ ...detail, signal_status: "NO_PRICE_EVIDENCE", market_signal: null });

    renderPanel();
    await screen.findAllByText("Camiseta azul");

    expect(screen.getByText("Sem evidência de preço de mercado")).toBeInTheDocument();
    expect(screen.queryByText("ACCEPT")).not.toBeInTheDocument();
    expect(screen.queryByText(/^\d+\/\d+$/)).not.toBeInTheDocument();
  });

  it("shows the sem-vínculo honest state with a link to /vinculos, no evidence section", async () => {
    getListing.mockResolvedValueOnce({
      ...detail,
      link: { state: "unresolved", product_id: null, seller_sku: null },
      signal_status: "SEM_VINCULO",
      market_signal: null,
    });

    renderPanel();
    await screen.findAllByText("Camiseta azul");

    expect(screen.getByRole("link", { name: "Vincular produto" })).toHaveAttribute("href", "/vinculos");
    expect(screen.queryByText("ACCEPT")).not.toBeInTheDocument();
    expect(screen.queryByText("Sem evidência de preço de mercado")).not.toBeInTheDocument();
  });

  it("defensively renders the honest-absent evidence state when signal_status is undefined, never throws", async () => {
    getListing.mockResolvedValueOnce({ ...detail });

    renderPanel();
    await screen.findAllByText("Camiseta azul");

    expect(screen.getByText("Sem evidência de preço de mercado")).toBeInTheDocument();
  });

  it("renders the 2x2 facts grid card labels from the ratified drawer (Preço/Est. publicado/Margem est./Qualidade)", async () => {
    getListing.mockResolvedValueOnce(detail);

    renderPanel();
    await screen.findAllByText("Camiseta azul");

    expect(screen.getByText("Preço")).toBeInTheDocument();
    expect(screen.getByText("Est. publicado")).toBeInTheDocument();
    expect(screen.getByText("Margem est.")).toBeInTheDocument();
    expect(screen.getByText("Qualidade")).toBeInTheDocument();
  });

  it("has no inline-edit affordance in the drawer body", async () => {
    getListing.mockResolvedValueOnce(detail);

    renderPanel();
    await screen.findAllByText("Camiseta azul");

    expect(screen.queryByRole("textbox")).not.toBeInTheDocument();
    expect(document.querySelector('input[type="text"]')).not.toBeInTheDocument();
    expect(document.querySelector("textarea")).not.toBeInTheDocument();
  });
});
