import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";
import type { MarketPriceIntelVerdict } from "@marketplace-central/sdk-runtime";
import { VeredictoBox } from "./VeredictoBox";
import type { ProdutoMarketClient } from "./marketClient";

function makeClient(overrides: Partial<ProdutoMarketClient> = {}): ProdutoMarketClient {
  return {
    listMarketVerdicts: vi.fn(),
    collectMarketPriceIntel: vi.fn(),
    listMarketSignals: vi.fn(),
    listMarketAggregates: vi.fn(),
    ...overrides,
  } as unknown as ProdutoMarketClient;
}

function renderBox(client: ProdutoMarketClient, queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } })) {
  return {
    queryClient,
    ...render(
      <QueryClientProvider client={queryClient}>
        <VeredictoBox productId="90008" client={client} />
      </QueryClientProvider>,
    ),
  };
}

describe("VeredictoBox", () => {
  it("renders the honest INSUFFICIENT_MARKET headline with the seller count", async () => {
    const verdict: MarketPriceIntelVerdict = {
      match_status: "ACCEPT",
      price_evidence_status: "INSUFFICIENT_MARKET",
      verdict_label: null,
      blocking_state: "INSUFFICIENT_MARKET",
      inputs_used: {},
      market_range: {
        min_valid: { amount: "10.00", currency: "BRL" },
        median: { amount: "12.00", currency: "BRL" },
        currency: "BRL",
        n_offers: 2,
        n_sellers: 2,
      },
    };
    const client = makeClient({ listMarketVerdicts: vi.fn().mockResolvedValue([verdict]) });

    renderBox(client);

    expect(await screen.findByText("2 vendedores — mínimo 5")).toBeInTheDocument();
  });

  it("never fabricates a seller count when market_range is null for INSUFFICIENT_MARKET (ADR-17)", async () => {
    const verdict: MarketPriceIntelVerdict = {
      match_status: "ACCEPT",
      price_evidence_status: "INSUFFICIENT_MARKET",
      verdict_label: null,
      blocking_state: "INSUFFICIENT_MARKET",
      inputs_used: {},
      market_range: null,
    };
    const client = makeClient({ listMarketVerdicts: vi.fn().mockResolvedValue([verdict]) });

    renderBox(client);

    const headline = await screen.findByText(/mínimo 5/);
    expect(headline.textContent).not.toMatch(/\b0\b/);
    expect(headline.textContent).not.toMatch(/^\d/);
  });

  it("renders faixa + evidence rows + M-07 unknown margin hint for an OK verdict", async () => {
    const verdict: MarketPriceIntelVerdict = {
      match_status: "ACCEPT",
      price_evidence_status: "OK",
      verdict_label: null,
      blocking_state: null,
      inputs_used: { ml: { source: "ml", as_of: "2026-07-15T10:00:00Z" } },
      market_range: {
        min_valid: { amount: "100.00", currency: "BRL" },
        median: { amount: "120.00", currency: "BRL" },
        currency: "BRL",
        n_offers: 8,
        n_sellers: 6,
      },
    };
    const client = makeClient({ listMarketVerdicts: vi.fn().mockResolvedValue([verdict]) });

    renderBox(client);

    // pt-BR through the shared formatter; the regex tolerates the nbsp Intl puts after "R$".
    expect(await screen.findByText(/^R\$\s100,00$/)).toBeInTheDocument();
    expect(screen.getByText(/^R\$\s120,00$/)).toBeInTheDocument();
    expect(screen.getByText("6")).toBeInTheDocument();
    expect(screen.getByText("ml")).toBeInTheDocument();
    expect(screen.getByText("2026-07-15T10:00:00Z")).toBeInTheDocument();
    expect(screen.getByTitle("veredicto de margem — M-07")).toBeInTheDocument();
  });

  it("never renders a fabricated price/green state for NO_PRICE_EVIDENCE", async () => {
    const verdict: MarketPriceIntelVerdict = {
      match_status: "NO_CANDIDATE",
      price_evidence_status: "NO_PRICE_EVIDENCE",
      verdict_label: null,
      blocking_state: "NO_PRICE_EVIDENCE",
      inputs_used: {},
      market_range: null,
    };
    const client = makeClient({ listMarketVerdicts: vi.fn().mockResolvedValue([verdict]) });

    renderBox(client);

    expect(await screen.findByText(/ainda não tem evidência de preço coletada/i)).toBeInTheDocument();
    expect(screen.queryByText(/R\$\s*0/)).not.toBeInTheDocument();
    expect(screen.queryByText(/R\$ 0,00/)).not.toBeInTheDocument();
  });

  it("runs a single synchronous collect + invalidates both query keys, no polling", async () => {
    const verdict: MarketPriceIntelVerdict = {
      match_status: "NO_CANDIDATE",
      price_evidence_status: "NO_PRICE_EVIDENCE",
      verdict_label: null,
      blocking_state: "NO_PRICE_EVIDENCE",
      inputs_used: {},
      market_range: null,
    };
    let resolveCollect: (value: { status: "COMPLETED" }) => void = () => undefined;
    const collectPromise = new Promise<{ status: "COMPLETED" }>((resolve) => {
      resolveCollect = resolve;
    });
    const collectMarketPriceIntel = vi.fn().mockReturnValue(collectPromise);
    const client = makeClient({
      listMarketVerdicts: vi.fn().mockResolvedValue([verdict]),
      collectMarketPriceIntel,
    });
    const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } });
    const invalidateSpy = vi.spyOn(queryClient, "invalidateQueries");

    renderBox(client, queryClient);

    const button = await screen.findByRole("button", { name: "coletar agora" });
    fireEvent.click(button);

    expect(await screen.findByRole("button", { name: "coletando…" })).toBeDisabled();
    expect(collectMarketPriceIntel).toHaveBeenCalledTimes(1);
    expect(collectMarketPriceIntel).toHaveBeenCalledWith("90008");

    resolveCollect({ status: "COMPLETED" });

    await waitFor(() => {
      expect(screen.getByRole("button", { name: "coletar agora" })).not.toBeDisabled();
    });

    expect(collectMarketPriceIntel).toHaveBeenCalledTimes(1);
    const invalidatedKeys = invalidateSpy.mock.calls.map((call) => (call[0] as { queryKey: unknown[] }).queryKey);
    expect(invalidatedKeys).toContainEqual(["market", "verdict", "90008"]);
    expect(invalidatedKeys).toContainEqual(["listings", "by-product"]);
  });

  it("shows an inline error and re-enables the button when collect rejects", async () => {
    const verdict: MarketPriceIntelVerdict = {
      match_status: "NO_CANDIDATE",
      price_evidence_status: "NO_PRICE_EVIDENCE",
      verdict_label: null,
      blocking_state: "NO_PRICE_EVIDENCE",
      inputs_used: {},
      market_range: null,
    };
    const collectMarketPriceIntel = vi.fn().mockRejectedValue({ status: 504, error: { code: "timeout", message: "boom" } });
    const client = makeClient({
      listMarketVerdicts: vi.fn().mockResolvedValue([verdict]),
      collectMarketPriceIntel,
    });

    renderBox(client);

    const button = await screen.findByRole("button", { name: "coletar agora" });
    fireEvent.click(button);

    await waitFor(() => {
      expect(screen.getByRole("alert")).toBeInTheDocument();
    });
    expect(screen.getByRole("button", { name: "coletar agora" })).not.toBeDisabled();
  });
});
