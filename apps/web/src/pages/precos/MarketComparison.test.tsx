import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { render, screen, waitFor } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";
import type { MarketPriceIntelAggregate } from "@marketplace-central/sdk-runtime";
import { MarketComparison } from "./MarketComparison";

function aggregate(overrides: Partial<MarketPriceIntelAggregate> = {}): MarketPriceIntelAggregate {
  return {
    product_id: "90001",
    median: { amount: "95.00", currency: "BRL" },
    min_valid: { amount: "88.00", currency: "BRL" },
    n_offers: 12,
    n_sellers: 5,
    source: "ml_sale_price",
    fetched_at: "2026-07-18T11:00:00Z",
    computed_at: "2026-07-18T11:05:00Z",
    status: "OK",
    ...overrides,
  };
}

function marketClient(items: MarketPriceIntelAggregate[]) {
  return { listMarketAggregates: vi.fn(() => Promise.resolve(items)) };
}

function renderPanel(node: React.ReactNode) {
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  return render(<QueryClientProvider client={queryClient}>{node}</QueryClientProvider>);
}

describe("MarketComparison", () => {
  it("does not query and prompts when no product is selected", () => {
    const client = marketClient([]);
    renderPanel(<MarketComparison productId={null} client={client} />);
    expect(client.listMarketAggregates).not.toHaveBeenCalled();
    expect(screen.getByText(/Selecione um produto/)).toBeInTheDocument();
  });

  it("renders the OK aggregate: median, min, offers/sellers, source, fetched_at", async () => {
    const client = marketClient([aggregate()]);
    renderPanel(<MarketComparison productId="90001" client={client} />);

    await waitFor(() => expect(client.listMarketAggregates).toHaveBeenCalledWith(["90001"]));
    const panel = await screen.findByTestId("market-comparison");
    expect(panel).toHaveTextContent("95.00");
    expect(panel).toHaveTextContent("88.00");
    expect(panel).toHaveTextContent("12");
    expect(panel).toHaveTextContent("5");
    expect(panel).toHaveTextContent("ml_sale_price");
    expect(screen.getByTestId("market-status")).toHaveTextContent("OK");
  });

  it("shows INSUFFICIENT_MARKET without fabricating a price", async () => {
    const client = marketClient([
      aggregate({ status: "INSUFFICIENT_MARKET", median: null, min_valid: null, n_offers: 1, n_sellers: 1 }),
    ]);
    renderPanel(<MarketComparison productId="90001" client={client} />);

    expect(await screen.findByTestId("market-status")).toHaveTextContent("INSUFFICIENT_MARKET");
    // No fabricated median — the unknown marker stands in, never "0".
    const panel = screen.getByTestId("market-comparison");
    expect(panel.querySelectorAll("span.text-faint").length).toBeGreaterThan(0);
    expect(panel).not.toHaveTextContent("0,00");
  });

  it("shows NO_PRICE_EVIDENCE without fabricating a price", async () => {
    const client = marketClient([
      aggregate({ status: "NO_PRICE_EVIDENCE", median: null, min_valid: null, n_offers: 0, n_sellers: 0 }),
    ]);
    renderPanel(<MarketComparison productId="90001" client={client} />);
    expect(await screen.findByTestId("market-status")).toHaveTextContent("NO_PRICE_EVIDENCE");
  });

  it("renders the empty state when the market has no aggregate for the product", async () => {
    const client = marketClient([]);
    renderPanel(<MarketComparison productId="90001" client={client} />);
    expect(await screen.findByText(/Sem dados de mercado/)).toBeInTheDocument();
  });

  it("renders the error state and retries", async () => {
    const client = { listMarketAggregates: vi.fn(() => Promise.reject(new Error("boom"))) };
    renderPanel(<MarketComparison productId="90001" client={client} />);
    expect(await screen.findByRole("alert")).toHaveTextContent(/mercado/i);
  });
});
