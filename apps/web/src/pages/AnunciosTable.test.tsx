import type { ListingMarketSignal, ListingReadModel } from "@marketplace-central/sdk-runtime";
import { fireEvent, render, screen } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import { describe, expect, it, vi } from "vitest";
import { AnunciosTable } from "./AnunciosTable";

function renderTable(ui: Parameters<typeof render>[0]) {
  return render(<MemoryRouter>{ui}</MemoryRouter>);
}

const tableSelectionProps = {
  selectedIds: new Set<string>(),
  onToggle: () => undefined,
  onTogglePage: () => undefined,
};

const listing: ListingReadModel = {
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
  sync_state: "synced",
  sync_error: null,
  quality_score: 0.9,
  pending_issue: { kind: "stale", message_pt: "Atualização pendente" },
  sales_30d: 12,
  cost: { amount: "70.00", currency: "BRL" },
  below_margin_worst_case: true,
  icms_worst_case_by_uf: null,
  fetched_at: "2026-07-16T12:00:00Z",
};

const okSignal: ListingMarketSignal = {
  status: "OK",
  position: { rank: 2, total: 8 },
  price_to_win: { amount: "119.90", currency: "BRL" },
  delta_pct: "8.34",
  match_status: "ACCEPT",
  n_offers: 6,
  n_sellers: 5,
  evidence: { source: "ml_price_to_win", fetched_at: "2026-07-16T11:00:00Z", freshness: "fresh" },
};

describe("AnunciosTable", () => {
  it("renders the listing facts and pending issue", () => {
    renderTable(<AnunciosTable items={[listing]} asOf="2026-07-16T12:00:00Z" {...tableSelectionProps} />);

    expect(screen.getByText("Camiseta azul")).toBeInTheDocument();
    expect(screen.getByText("MLB123456789")).toBeInTheDocument();
    expect(screen.getByText("Clássico")).toBeInTheDocument();
    expect(screen.getByText("vinculado")).toBeInTheDocument();
    expect(screen.getByText("R$ 129.90")).toBeInTheDocument();
    expect(screen.getByText("7")).toBeInTheDocument();
    expect(screen.getByText("12")).toBeInTheDocument();
    expect(screen.getByText("abaixo da margem")).toBeInTheDocument();
    expect(screen.getByText("sincronizado")).toBeInTheDocument();
    expect(screen.getByLabelText("Atualização pendente")).toHaveAttribute("title", "Atualização pendente");
  });

  it("renders honest unknowns for nullable values", () => {
    renderTable(
      <AnunciosTable
        items={[{ ...listing, price: null, published_quantity: null, sales_30d: null, below_margin_worst_case: null }]}
        {...tableSelectionProps}
      />,
    );

    // 4 pre-existing unknowns (price/estoque/vendas/margem) + 1 honest absent
    // VS MERCADO state for a listing with no market_signal/signal_status.
    expect(screen.getAllByText("—")).toHaveLength(5);
    expect(screen.queryByText("0")).not.toBeInTheDocument();
    expect(screen.queryByText("R$ 0,00")).not.toBeInTheDocument();
  });

  it("explains an unsimulated margin when ERP cost is unknown", () => {
    renderTable(<AnunciosTable items={[{ ...listing, cost: null, below_margin_worst_case: null }]} {...tableSelectionProps} />);

    expect(screen.getByTitle("sem custo no ERP → não simulado")).toHaveTextContent("—");
  });

  it("opens the listing from the title button without reacting to checkbox clicks", () => {
    const onOpen = vi.fn();
    renderTable(<AnunciosTable items={[listing]} {...tableSelectionProps} onOpen={onOpen} />);

    fireEvent.click(screen.getByLabelText("Selecionar anúncio Camiseta azul"));
    expect(onOpen).not.toHaveBeenCalled();

    fireEvent.click(screen.getByRole("button", { name: "Camiseta azul" }));
    expect(onOpen).toHaveBeenCalledOnce();
    expect(onOpen).toHaveBeenCalledWith("listing_1");
  });

  it("uses the fixed sync and conflict labels", () => {
    renderTable(<AnunciosTable items={[{ ...listing, sync_state: "error", link: { ...listing.link, state: "conflict" }}]} {...tableSelectionProps} />);

    expect(screen.getByText("com erro")).toBeInTheDocument();
    expect(screen.getByText("divergente")).toBeInTheDocument();
  });

  describe("VS MERCADO column (signal_status)", () => {
    it("OK: renders position and delta vs price_to_win in mono, no fabricated fallback", () => {
      renderTable(
        <AnunciosTable
          items={[{ ...listing, signal_status: "OK", market_signal: okSignal }]}
          {...tableSelectionProps}
        />,
      );

      expect(screen.getByText("2/8")).toBeInTheDocument();
      expect(screen.getByText("+8.34%")).toBeInTheDocument();
      expect(screen.getByText("2/8").closest(".font-mono")).not.toBeNull();
    });

    it("SEM_VINCULO: renders a grey chip linking to /vinculos, no numeric price", () => {
      renderTable(
        <AnunciosTable
          items={[{ ...listing, signal_status: "SEM_VINCULO", market_signal: null }]}
          {...tableSelectionProps}
        />,
      );

      const link = screen.getByRole("link", { name: /vínculo/i });
      expect(link).toHaveAttribute("href", "/vinculos");
      expect(screen.queryByText("2/8")).not.toBeInTheDocument();
      expect(screen.queryByText(/%/)).not.toBeInTheDocument();
    });

    it("NO_PRICE_EVIDENCE: renders UnknownValue with a tooltip, no number", () => {
      renderTable(
        <AnunciosTable
          items={[{ ...listing, signal_status: "NO_PRICE_EVIDENCE", market_signal: null }]}
          {...tableSelectionProps}
        />,
      );

      expect(screen.getByTitle("sem evidência de preço de mercado")).toHaveTextContent("—");
      expect(screen.queryByText(/%/)).not.toBeInTheDocument();
    });

    it("STALE: shows the value plus a freshness (amber age) indicator", () => {
      renderTable(
        <AnunciosTable
          items={[
            {
              ...listing,
              signal_status: "STALE",
              market_signal: { ...okSignal, status: "STALE" },
            },
          ]}
          {...tableSelectionProps}
        />,
      );

      expect(screen.getByText("2/8")).toBeInTheDocument();
      expect(screen.getByLabelText("Data freshness")).toBeInTheDocument();
    });

    it("renders all 4 states in one table without throwing, and a missing signal_status is treated as the honest absent state", () => {
      const items = [
        { ...listing, listing_id: "l_ok", signal_status: "OK" as const, market_signal: okSignal },
        { ...listing, listing_id: "l_sem_vinculo", signal_status: "SEM_VINCULO" as const, market_signal: null },
        { ...listing, listing_id: "l_no_evidence", signal_status: "NO_PRICE_EVIDENCE" as const, market_signal: null },
        { ...listing, listing_id: "l_stale", signal_status: "STALE" as const, market_signal: { ...okSignal, status: "STALE" as const } },
        { ...listing, listing_id: "l_missing", signal_status: undefined, market_signal: undefined },
      ];

      expect(() => renderTable(<AnunciosTable items={items} {...tableSelectionProps} />)).not.toThrow();

      expect(screen.getAllByText("2/8")).toHaveLength(2);
      expect(screen.getByRole("link", { name: /vínculo/i })).toBeInTheDocument();
      expect(screen.getAllByText("—").length).toBeGreaterThanOrEqual(1);
    });
  });

  describe("group rendering (agrupar por produto)", () => {
    it("renders a group header per CODPROD with per-listing signal columns preserved inside the group", () => {
      renderTable(
        <AnunciosTable
          groups={[
            {
              product_id: "product_1",
              product_title: "Grupo Camisetas",
              listing_count: 1,
              group_state: "ok",
              listings: [{ ...listing, signal_status: "OK", market_signal: okSignal }],
            },
          ]}
          {...tableSelectionProps}
        />,
      );

      expect(screen.getByText("Grupo Camisetas")).toBeInTheDocument();
      expect(screen.getByText("2/8")).toBeInTheDocument();
    });

    it("renders a 1-listing group as a normal group, no special collapse", () => {
      renderTable(
        <AnunciosTable
          groups={[
            {
              product_id: "product_1",
              product_title: "Meia branca",
              listing_count: 1,
              group_state: "ok",
              listings: [{ ...listing, listing_id: "listing_solo", title: "Meia branca" }],
            },
          ]}
          {...tableSelectionProps}
        />,
      );

      expect(screen.getAllByText("Meia branca").length).toBeGreaterThanOrEqual(1);
      expect(screen.getByRole("button", { name: "Meia branca" })).toBeInTheDocument();
    });

    it("renders the synthetic sem-produto group with its honest label", () => {
      renderTable(
        <AnunciosTable
          groups={[
            {
              product_id: null,
              product_title: null,
              listing_count: 1,
              group_state: "attention",
              listings: [listing],
            },
          ]}
          {...tableSelectionProps}
        />,
      );

      expect(screen.getByText("Sem produto")).toBeInTheDocument();
    });
  });
});
