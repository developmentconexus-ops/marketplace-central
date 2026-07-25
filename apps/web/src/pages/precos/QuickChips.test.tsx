import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";
import type { MarketPriceIntelAggregate } from "@marketplace-central/sdk-runtime";
import { QuickChips } from "./QuickChips";
import type { MarketAggregatesClient } from "./MarketComparison";

const pricedAggregate: MarketPriceIntelAggregate = {
  product_id: "90001",
  status: "OK",
  median: { amount: "75.40", currency: "BRL" },
  min_valid: { amount: "69.90", currency: "BRL" },
  n_offers: 23,
  n_sellers: 19,
  source: "ml_catalog_offers",
  fetched_at: "2026-07-18T12:00:00Z",
  computed_at: "2026-07-18T12:00:00Z",
};

function makeClient(items: MarketPriceIntelAggregate[]): MarketAggregatesClient {
  return { listMarketAggregates: vi.fn(() => Promise.resolve(items)) };
}

function renderChips(
  props: Partial<React.ComponentProps<typeof QuickChips>> = {},
  items: MarketPriceIntelAggregate[] = [pricedAggregate],
) {
  const onSeedPreco = vi.fn();
  const onSeedTarget = vi.fn();
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  render(
    <QueryClientProvider client={queryClient}>
      <QuickChips
        productId="90001"
        coverPct="18"
        onSeedPreco={onSeedPreco}
        onSeedTarget={onSeedTarget}
        client={makeClient(items)}
        {...props}
      />
    </QueryClientProvider>,
  );
  return { onSeedPreco, onSeedTarget };
}

describe("QuickChips", () => {
  it("renders mediana + menor conc. + cobrir chips and seeds the price from the median", async () => {
    const { onSeedPreco } = renderChips();
    const mediana = await screen.findByTestId("quick-chip-mediana");
    expect(mediana).toHaveTextContent("mediana R$ 75,40");
    fireEvent.click(mediana);
    expect(onSeedPreco).toHaveBeenCalledWith("75.40");
  });

  it("seeds the price from the lowest valid offer", async () => {
    const { onSeedPreco } = renderChips();
    fireEvent.click(await screen.findByTestId("quick-chip-menor"));
    expect(onSeedPreco).toHaveBeenCalledWith("69.90");
  });

  it("seeds the solver target from the cover chip (client-side, never a fabricated price)", () => {
    const { onSeedTarget } = renderChips();
    fireEvent.click(screen.getByTestId("quick-chip-cobrir"));
    expect(onSeedTarget).toHaveBeenCalledWith("18");
  });

  it("omits mediana/menor chips when the market has no confident price (ADR-17) — cobrir still shows", async () => {
    renderChips({}, [
      { ...pricedAggregate, status: "INSUFFICIENT_MARKET", median: null, min_valid: null },
    ]);
    // cobrir always available (target seed, not a price).
    expect(await screen.findByTestId("quick-chip-cobrir")).toBeInTheDocument();
    await waitFor(() => expect(screen.queryByTestId("quick-chip-mediana")).toBeNull());
    expect(screen.queryByTestId("quick-chip-menor")).toBeNull();
  });
});
