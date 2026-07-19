import type { ListingMarketSignal, ListingReadModel } from "@marketplace-central/sdk-runtime";
import { fireEvent, render, screen, within } from "@testing-library/react";
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
  quality_score: 92,
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

const belowMarketSignal: ListingMarketSignal = { ...okSignal, delta_pct: "-4.20" };

describe("AnunciosTable", () => {
  it("renders the ratified 9-column header, no dropped columns", () => {
    renderTable(<AnunciosTable items={[listing]} {...tableSelectionProps} />);

    const headers = screen.getAllByRole("columnheader").map((th) => th.textContent);
    expect(headers).toEqual(["", "MLB", "TÍTULO", "PRODUTO", "PREÇO", "EST.", "SYNC", "QUAL.", "PENDÊNCIA"]);
    expect(screen.queryByText("Vs. mercado")).not.toBeInTheDocument();
    expect(screen.queryByText("Modalidade")).not.toBeInTheDocument();
    expect(screen.queryByText("Vendas 30d")).not.toBeInTheDocument();
    expect(screen.queryByText("Margem")).not.toBeInTheDocument();
  });

  it("renders the listing facts", () => {
    renderTable(<AnunciosTable items={[listing]} asOf="2026-07-16T12:00:00Z" {...tableSelectionProps} />);

    expect(screen.getByText("Camiseta azul")).toBeInTheDocument();
    expect(screen.getByText("MLB123456789")).toBeInTheDocument();
    expect(screen.getByText("product_1")).toBeInTheDocument();
    expect(screen.getByText("R$ 129.90")).toBeInTheDocument();
    expect(screen.getByText("7")).toBeInTheDocument();
    expect(screen.getByText("sincronizado")).toBeInTheDocument();
    // quality_score is an int 0–100 — 92 must render "92%", never "9200%".
    expect(screen.getByText("92%")).toBeInTheDocument();
    expect(screen.queryByText("9200%")).not.toBeInTheDocument();
    expect(screen.getByTitle("Atualização pendente")).toHaveTextContent("Atualização pendente");
  });

  it("renders honest unknowns for nullable values, never a fabricated 0", () => {
    renderTable(
      <AnunciosTable
        items={[
          {
            ...listing,
            price: null,
            published_quantity: null,
            quality_score: null,
            pending_issue: null,
            link: { state: "unresolved", product_id: null, seller_sku: null },
          },
        ]}
        {...tableSelectionProps}
      />,
    );

    // price (—), EST (—), QUAL (—) = 3 honest unknowns. Pendência null is an
    // honest EMPTY cell (no pending issue), not an UnknownValue "—".
    expect(screen.getAllByText("—")).toHaveLength(3);
    expect(screen.queryByText("0")).not.toBeInTheDocument();
    expect(screen.queryByText("0%")).not.toBeInTheDocument();
    expect(screen.queryByText("R$ 0,00")).not.toBeInTheDocument();
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

  it("selects all rows on the page via the header checkbox", () => {
    const onTogglePage = vi.fn();
    renderTable(
      <AnunciosTable
        items={[listing, { ...listing, listing_id: "listing_2", title: "Bermuda cinza" }]}
        {...tableSelectionProps}
        onTogglePage={onTogglePage}
      />,
    );

    fireEvent.click(screen.getByLabelText("Selecionar todos os anúncios desta página"));
    expect(onTogglePage).toHaveBeenCalledWith(["listing_1", "listing_2"]);
  });

  it("uses the fixed sync pill label and a conflict tag for a divergent link", () => {
    renderTable(
      <AnunciosTable
        items={[{ ...listing, sync_state: "error", link: { ...listing.link, state: "conflict" } }]}
        {...tableSelectionProps}
      />,
    );

    expect(screen.getByText("com erro")).toBeInTheDocument();
    expect(screen.getByText("divergente")).toBeInTheDocument();
  });

  describe("PREÇO cell (value + %chip vs mercado, folded per RULING D-56)", () => {
    it("OK: renders the price value plus a sign-colored %chip, no fabricated fallback", () => {
      renderTable(
        <AnunciosTable items={[{ ...listing, signal_status: "OK", market_signal: okSignal }]} {...tableSelectionProps} />,
      );

      expect(screen.getByText("R$ 129.90")).toBeInTheDocument();
      expect(screen.getByText("+8.34%")).toBeInTheDocument();
    });

    it("OK below market: renders a negative %chip", () => {
      renderTable(
        <AnunciosTable
          items={[{ ...listing, signal_status: "OK", market_signal: belowMarketSignal }]}
          {...tableSelectionProps}
        />,
      );

      expect(screen.getByText("-4.20%")).toBeInTheDocument();
    });

    it("STALE: shows the value, an âmbar %chip, and a freshness age marker — never a double underline", () => {
      renderTable(
        <AnunciosTable
          items={[{ ...listing, signal_status: "STALE", market_signal: { ...okSignal, status: "STALE" } }]}
          {...tableSelectionProps}
        />,
      );

      expect(screen.getByText("R$ 129.90")).toBeInTheDocument();
      expect(screen.getByText("+8.34%")).toBeInTheDocument();
      expect(screen.getByLabelText("Data freshness")).toBeInTheDocument();
    });

    it("NO_PRICE_EVIDENCE: renders the price value only, no %chip, no fabricated number", () => {
      renderTable(
        <AnunciosTable
          items={[{ ...listing, signal_status: "NO_PRICE_EVIDENCE", market_signal: null }]}
          {...tableSelectionProps}
        />,
      );

      expect(screen.getByText("R$ 129.90")).toBeInTheDocument();
      expect(within(screen.getByTestId("preco-cell")).queryByText(/%/)).not.toBeInTheDocument();
    });

    it("SEM_VINCULO: renders the price value plus a link to /vinculos, no market number", () => {
      renderTable(
        <AnunciosTable
          items={[{ ...listing, signal_status: "SEM_VINCULO", market_signal: null }]}
          {...tableSelectionProps}
        />,
      );

      expect(screen.getByText("R$ 129.90")).toBeInTheDocument();
      const link = screen.getByRole("link", { name: "sem vínculo" });
      expect(link).toHaveAttribute("href", "/vinculos");
      expect(within(screen.getByTestId("preco-cell")).queryByText(/%/)).not.toBeInTheDocument();
    });

    it("a missing signal_status is treated as the honest absent state (defensive, never throws)", () => {
      expect(() =>
        renderTable(
          <AnunciosTable items={[{ ...listing, signal_status: undefined, market_signal: undefined }]} {...tableSelectionProps} />,
        ),
      ).not.toThrow();

      expect(screen.getByText("R$ 129.90")).toBeInTheDocument();
      expect(within(screen.getByTestId("preco-cell")).queryByText(/%/)).not.toBeInTheDocument();
    });

    it("renders all 4 signal states in one table without throwing", () => {
      const items = [
        { ...listing, listing_id: "l_ok", signal_status: "OK" as const, market_signal: okSignal },
        { ...listing, listing_id: "l_sem_vinculo", signal_status: "SEM_VINCULO" as const, market_signal: null },
        { ...listing, listing_id: "l_no_evidence", signal_status: "NO_PRICE_EVIDENCE" as const, market_signal: null },
        { ...listing, listing_id: "l_stale", signal_status: "STALE" as const, market_signal: { ...okSignal, status: "STALE" as const } },
      ];

      expect(() => renderTable(<AnunciosTable items={items} {...tableSelectionProps} />)).not.toThrow();

      expect(screen.getAllByText("+8.34%")).toHaveLength(2);
      expect(screen.getByRole("link", { name: "sem vínculo" })).toBeInTheDocument();
    });
  });

  describe("group rendering (agrupar por produto)", () => {
    it("renders a group header per CODPROD with per-listing PREÇO chip preserved inside the group", () => {
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
      expect(screen.getByText("+8.34%")).toBeInTheDocument();
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

    // Golden rows per dispatch EXEMPLO-IO (design-screens §Anúncios group-row
    // spec: chevron + meta "ERP est. N · M anúncios" + pill ok/N erro).
    // ERP est. is genuinely unavailable FE-only (ListingGroup carries no ERP
    // stock field, seam FROZEN T-1) → honest "—" per ADR-17, never fabricated.
    describe("enriched group header (chevron + meta + error pill)", () => {
      const grp = (overrides: Partial<import("@marketplace-central/sdk-runtime").ListingGroup> = {}) => ({
        product_id: "product_1",
        product_title: "Produto X",
        listing_count: 3,
        group_state: "ok" as const,
        listings: [
          { ...listing, listing_id: "g1", sync_state: "synced" as const },
          { ...listing, listing_id: "g2", sync_state: "synced" as const },
          { ...listing, listing_id: "g3", sync_state: "synced" as const },
        ],
        ...overrides,
      });

      it("GOLDEN: 3 listings, 1 sync error → 'ERP est. —', '3 anúncios', '1 erro'", () => {
        renderTable(
          <AnunciosTable
            groups={[grp({ listings: [
              { ...listing, listing_id: "g1", sync_state: "error" },
              { ...listing, listing_id: "g2", sync_state: "synced" },
              { ...listing, listing_id: "g3", sync_state: "synced" },
            ] })]}
            {...tableSelectionProps}
          />,
        );

        const header = screen.getByText("Produto X").closest("tr")!;
        expect(within(header).getByText(/ERP est\.\s*—/)).toBeInTheDocument();
        expect(within(header).getByText(/3\s+anúncios/)).toBeInTheDocument();
        expect(within(header).getByText("1 erro")).toBeInTheDocument();
        // ADR-17: never a fabricated ERP stock 0.
        expect(within(header).queryByText(/ERP est\.\s*0/)).not.toBeInTheDocument();
      });

      it("GOLDEN: 3 listings, all synced → 'ERP est. —', '✓ ok' pill, no error count", () => {
        renderTable(<AnunciosTable groups={[grp()]} {...tableSelectionProps} />);

        const header = screen.getByText("Produto X").closest("tr")!;
        expect(within(header).getByText(/ERP est\.\s*—/)).toBeInTheDocument();
        expect(within(header).getByText(/3\s+anúncios/)).toBeInTheDocument();
        expect(within(header).getByText("ok")).toBeInTheDocument();
        expect(within(header).queryByText(/erro/)).not.toBeInTheDocument();
      });

      it("pluralizes honestly: 1 listing → 'anúncio', 2 errors → 'erros'", () => {
        renderTable(
          <AnunciosTable
            groups={[grp({
              listing_count: 2,
              listings: [
                { ...listing, listing_id: "g1", sync_state: "error" },
                { ...listing, listing_id: "g2", sync_state: "error" },
              ],
            })]}
            {...tableSelectionProps}
          />,
        );

        const header = screen.getByText("Produto X").closest("tr")!;
        expect(within(header).getByText(/2\s+anúncios/)).toBeInTheDocument();
        expect(within(header).getByText("2 erros")).toBeInTheDocument();
      });

      it("chevron collapses the group: child listing rows hidden when collapsed", () => {
        renderTable(<AnunciosTable groups={[grp()]} {...tableSelectionProps} />);

        // Expanded by default → child rows present.
        expect(screen.getAllByRole("button", { name: "Camiseta azul" })).toHaveLength(3);

        const toggle = screen.getByRole("button", { name: /grupo Produto X/i });
        expect(toggle).toHaveAttribute("aria-expanded", "true");
        fireEvent.click(toggle);

        expect(toggle).toHaveAttribute("aria-expanded", "false");
        expect(screen.queryByRole("button", { name: "Camiseta azul" })).not.toBeInTheDocument();
        // Header (title + meta) stays visible when collapsed.
        expect(screen.getByText("Produto X")).toBeInTheDocument();
      });
    });
  });
});
