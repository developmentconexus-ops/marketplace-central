import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import { QueryClientProvider } from "@tanstack/react-query";
import { describe, expect, it, vi, afterEach } from "vitest";
import { createMarketplaceCentralClient } from "@marketplace-central/sdk-runtime";
import {
  createRefreshableFetch,
  createWebQueryClient,
  QUERY_STALE_TIME,
} from "@marketplace-central/web-query";
import { CatalogPage } from "./CatalogPage";

function page(id: number, next_cursor: string | null, as_of: string) {
  return {
    items: [{
      internal_product_id: id,
      reference: `REF-${id}`,
      description: `Product ${id}`,
      ean: null,
      active: true,
      sellable_stock: { quantity: id, quality: [] },
      current_price: { amount: "12.90", currency: "BRL", quality: [] },
      cost: { amount: id === 1 ? null : "4.50", currency: "BRL", quality: id === 1 ? ["missing_cost"] : [] },
    }],
    next_cursor,
    page_size: 1,
    as_of,
  };
}

function renderPage(client: any, queryClient = createWebQueryClient(), erpSource?: any) {
  return render(
    <QueryClientProvider client={queryClient}>
      <CatalogPage client={client} erpSource={erpSource} />
    </QueryClientProvider>,
  );
}

afterEach(() => {
  vi.useRealTimers();
  vi.setSystemTime();
});

function requestUrl(call: unknown[]): string {
  const input = call[0];
  return typeof input === "string" ? input : String((input as Request).url);
}

describe("CatalogPage", () => {
  it("appends infinite pages and stops when next_cursor is null", async () => {
    const fetch = vi.fn()
      .mockResolvedValueOnce(new Response(JSON.stringify(page(1, "cursor-2", "2026-07-14T10:11:12Z"))))
      .mockResolvedValueOnce(new Response(JSON.stringify(page(2, null, "2026-07-14T10:12:13Z"))));
    const transport = createRefreshableFetch(fetch as unknown as typeof globalThis.fetch);
    const sdk = createMarketplaceCentralClient({ baseUrl: "", fetchImpl: transport.fetchImpl });
    renderPage({ ...sdk, withNoCache: transport.withNoCache });

    expect(await screen.findByText("Product 1")).toBeInTheDocument();
    expect(screen.getByText("— (missing_cost)")).toBeInTheDocument();
    fireEvent.click(screen.getByRole("button", { name: "Carregar mais" }));
    expect(await screen.findByText("Product 2")).toBeInTheDocument();
    expect(fetch).toHaveBeenCalledTimes(2);
    expect(screen.getByText("Fim da lista")).toBeInTheDocument();
  });

  it("uses catalog staleTime across remounts and refetches after it expires", async () => {
    const fetch = vi.fn().mockResolvedValue(new Response(JSON.stringify(page(1, null, "2026-07-14T10:11:12Z"))));
    const transport = createRefreshableFetch(fetch as unknown as typeof globalThis.fetch);
    const sdk = createMarketplaceCentralClient({ baseUrl: "", fetchImpl: transport.fetchImpl });
    const queryClient = createWebQueryClient();
    const first = renderPage({ ...sdk, withNoCache: transport.withNoCache }, queryClient);
    await screen.findByText("Product 1");
    first.unmount();
    renderPage({ ...sdk, withNoCache: transport.withNoCache }, queryClient);
    await screen.findByText("Product 1");
    expect(fetch).toHaveBeenCalledTimes(1);

    vi.setSystemTime(Date.now() + QUERY_STALE_TIME.catalog + 1);
    first.unmount();
    renderPage({ ...sdk, withNoCache: transport.withNoCache }, queryClient);
    await waitFor(() => expect(fetch).toHaveBeenCalledTimes(2));
  });

  it("never sends erp_source: the server resolves the tenant's active source", async () => {
    // The routing reader looks the active source up per tenant and pins it on
    // the request context, overwriting anything the client sends. A request
    // parameter here would be a no-op that makes the screen look source-aware
    // while reading whatever the database says.
    const fetch = vi.fn().mockResolvedValue(new Response(JSON.stringify(page(1, null, "2026-07-14T10:11:12Z"))));
    const transport = createRefreshableFetch(fetch as unknown as typeof globalThis.fetch);
    const sdk = createMarketplaceCentralClient({ baseUrl: "", fetchImpl: transport.fetchImpl });
    renderPage({ ...sdk, withNoCache: transport.withNoCache }, createWebQueryClient(), "catalogo_cliente");
    await screen.findByText("Product 1");
    expect(requestUrl(fetch.mock.calls[0])).not.toContain("erp_source");
  });

  it("keeps one source's rows out of the other source's cache entry", async () => {
    const fetch = vi.fn()
      .mockResolvedValueOnce(new Response(JSON.stringify(page(1, null, "2026-07-14T10:11:12Z"))))
      .mockResolvedValueOnce(new Response(JSON.stringify(page(2, null, "2026-07-14T10:12:13Z"))));
    const transport = createRefreshableFetch(fetch as unknown as typeof globalThis.fetch);
    const sdk = createMarketplaceCentralClient({ baseUrl: "", fetchImpl: transport.fetchImpl });
    const client = { ...sdk, withNoCache: transport.withNoCache };
    const queryClient = createWebQueryClient();

    const first = renderPage(client, queryClient, "xlsx");
    await screen.findByText("Product 1");
    first.unmount();

    // Same query client, different active source: a shared key would have
    // served Product 1 (the other source's row) with no request at all.
    renderPage(client, queryClient, "catalogo_cliente");
    await screen.findByText("Product 2");
    expect(fetch).toHaveBeenCalledTimes(2);
  });

  it("renders local as_of and sends no-cache on refresh", async () => {
    const fetch = vi.fn()
      .mockResolvedValueOnce(new Response(JSON.stringify(page(1, null, "2026-07-14T10:11:12Z"))))
      .mockResolvedValueOnce(new Response(JSON.stringify(page(1, null, "2026-07-14T11:12:13Z"))));
    const transport = createRefreshableFetch(fetch as unknown as typeof globalThis.fetch);
    const sdk = createMarketplaceCentralClient({ baseUrl: "", fetchImpl: transport.fetchImpl });
    renderPage({ ...sdk, withNoCache: transport.withNoCache });

    const indicator = await screen.findByLabelText("Data freshness");
    await waitFor(() => expect(indicator.textContent).toMatch(/^dados de \d{2}:\d{2}:\d{2}$/));
    fireEvent.click(screen.getByRole("button", { name: "Atualizar" }));
    await waitFor(() => expect(fetch).toHaveBeenCalledTimes(2));
    const headers = new Headers(fetch.mock.calls[1][1]?.headers);
    expect(headers.get("Cache-Control")).toBe("no-cache");
    await waitFor(() => expect(indicator.textContent).toMatch(/dados de \d{2}:\d{2}:\d{2}/));
  });
});
