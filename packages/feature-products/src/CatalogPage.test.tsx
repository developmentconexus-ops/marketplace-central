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
    items: [
      {
        internal_product_id: id,
        reference: `REF-${id}`,
        description: `Product ${id}`,
        ean: null,
        active: true,
        sellable_stock: { quantity: id, quality: [] },
        current_price: { amount: "12.90", currency: "BRL", quality: [] },
        cost: {
          amount: id === 1 ? null : "4.50",
          currency: "BRL",
          quality: id === 1 ? ["missing_cost"] : [],
        },
      },
    ],
    next_cursor,
    page_size: 1,
    as_of,
  };
}

// Unlike page(), the stock quantity here is a test input, not derived from
// the id — S12's Sem-estoque / honest-unknown split needs a 0 and a null in
// the same page, and page()'s id-as-quantity shortcut cannot produce either.
function pageOf(
  items: Array<{ id: number; quantity: number | null; quality?: string[] }>,
  next_cursor: string | null,
  as_of: string,
) {
  return {
    items: items.map(({ id, quantity, quality = [] }) => ({
      internal_product_id: id,
      reference: `REF-${id}`,
      description: `Product ${id}`,
      ean: null,
      active: true,
      sellable_stock: { quantity, quality },
      current_price: { amount: "12.90", currency: "BRL", quality: [] },
      cost: { amount: "4.50", currency: "BRL", quality: [] },
    })),
    next_cursor,
    page_size: items.length,
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
    const fetch = vi
      .fn()
      .mockResolvedValueOnce(
        new Response(JSON.stringify(page(1, "cursor-2", "2026-07-14T10:11:12Z"))),
      )
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

  // Search used to render its first page and a line saying only the first page
  // was shown — every match past the 50th was unreachable. It now walks the
  // cursor like the list read, and the cursor must reach the search endpoint.
  it("pages search results past the first page", async () => {
    const fetch = vi
      .fn()
      // Call 0 is the list read the mount fires before the search debounce lands.
      .mockResolvedValueOnce(new Response(JSON.stringify(page(1, null, "2026-07-14T10:10:11Z"))))
      .mockResolvedValueOnce(
        new Response(JSON.stringify(page(11, "cursor-s2", "2026-07-14T10:11:12Z"))),
      )
      .mockResolvedValueOnce(new Response(JSON.stringify(page(12, null, "2026-07-14T10:12:13Z"))));
    const transport = createRefreshableFetch(fetch as unknown as typeof globalThis.fetch);
    const sdk = createMarketplaceCentralClient({ baseUrl: "", fetchImpl: transport.fetchImpl });
    renderPage({ ...sdk, withNoCache: transport.withNoCache });

    fireEvent.change(screen.getByLabelText("Buscar no catálogo"), { target: { value: "CUBA" } });

    expect(await screen.findByText("Product 11")).toBeInTheDocument();
    expect(requestUrl(fetch.mock.calls[1])).toContain("/catalog/products/search?q=CUBA");
    fireEvent.click(screen.getByRole("button", { name: "Carregar mais" }));
    expect(await screen.findByText("Product 12")).toBeInTheDocument();
    expect(requestUrl(fetch.mock.calls[2])).toContain("cursor=cursor-s2");
    expect(screen.getByText("Fim dos resultados")).toBeInTheDocument();
  });

  it("uses catalog staleTime across remounts and refetches after it expires", async () => {
    const fetch = vi
      .fn()
      .mockResolvedValue(new Response(JSON.stringify(page(1, null, "2026-07-14T10:11:12Z"))));
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
    const fetch = vi
      .fn()
      .mockResolvedValue(new Response(JSON.stringify(page(1, null, "2026-07-14T10:11:12Z"))));
    const transport = createRefreshableFetch(fetch as unknown as typeof globalThis.fetch);
    const sdk = createMarketplaceCentralClient({ baseUrl: "", fetchImpl: transport.fetchImpl });
    renderPage(
      { ...sdk, withNoCache: transport.withNoCache },
      createWebQueryClient(),
      "catalogo_cliente",
    );
    await screen.findByText("Product 1");
    expect(requestUrl(fetch.mock.calls[0])).not.toContain("erp_source");
  });

  it("keeps one source's rows out of the other source's cache entry", async () => {
    // S12 added the assortment-counts chip, which also fetches once per mount
    // whenever erpSource is set (same as the facts read). The two ordered
    // Once-responses this test used to queue assumed a single fetch per
    // mount; branching by URL keeps the counts read from stealing a slot
    // meant for the facts page, so the facts-call assertion below still
    // isolates the thing this test is actually about (cache partitioning).
    const factsPages = [
      page(1, null, "2026-07-14T10:11:12Z"),
      page(2, null, "2026-07-14T10:12:13Z"),
    ];
    const fetch = vi.fn((input: RequestInfo | URL) => {
      const url = typeof input === "string" ? input : String((input as Request).url);
      if (url.includes("/catalog/products/counts")) {
        return Promise.resolve(new Response(JSON.stringify({ sellable_count: 0, total_count: 0 })));
      }
      return Promise.resolve(new Response(JSON.stringify(factsPages.shift())));
    });
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
    expect(fetch.mock.calls.filter((call) => !requestUrl(call).includes("/counts"))).toHaveLength(
      2,
    );
  });

  it("renders local as_of and sends no-cache on refresh", async () => {
    const fetch = vi
      .fn()
      .mockResolvedValueOnce(new Response(JSON.stringify(page(1, null, "2026-07-14T10:11:12Z"))))
      .mockResolvedValueOnce(new Response(JSON.stringify(page(1, null, "2026-07-14T11:12:13Z"))));
    const transport = createRefreshableFetch(fetch as unknown as typeof globalThis.fetch);
    const sdk = createMarketplaceCentralClient({ baseUrl: "", fetchImpl: transport.fetchImpl });
    renderPage({ ...sdk, withNoCache: transport.withNoCache });

    const indicator = await screen.findByLabelText("Atualização dos dados");
    await waitFor(() => expect(indicator.textContent).toMatch(/^(agora|há .+)$/));
    fireEvent.click(screen.getByRole("button", { name: "Atualizar" }));
    await waitFor(() => expect(fetch).toHaveBeenCalledTimes(2));
    const headers = new Headers(fetch.mock.calls[1][1]?.headers);
    expect(headers.get("Cache-Control")).toBe("no-cache");
    await waitFor(() => expect(indicator.textContent).toMatch(/^(agora|há .+)$/));
  });

  it("opens with Vendáveis 2 de 4 and ver todos never mutates tenant config", async () => {
    const fetch = vi
      .fn()
      .mockResolvedValue(new Response(JSON.stringify(page(1, null, "2026-07-14T10:11:12Z"))));
    const transport = createRefreshableFetch(fetch as unknown as typeof globalThis.fetch);
    const sdk = createMarketplaceCentralClient({ baseUrl: "", fetchImpl: transport.fetchImpl });
    // Overriding these two (not going through `fetch`) is what gives the
    // "never mutates" assertion a world where it can fail: a client without
    // setSellableAssortment cannot be called, so the assertion below would
    // pass even if the button wired the mutation in by mistake.
    const setSellableAssortment = vi.fn();
    const getCatalogAssortmentCounts = vi
      .fn()
      .mockResolvedValue({ sellable_count: 2, total_count: 4 });
    const client = {
      ...sdk,
      withNoCache: transport.withNoCache,
      setSellableAssortment,
      getCatalogAssortmentCounts,
    };

    renderPage(client, createWebQueryClient(), "xlsx");

    expect(await screen.findByText("Vendáveis 2 de 4")).toBeInTheDocument();
    fireEvent.click(screen.getByRole("button", { name: "Ver todos" }));
    await screen.findByRole("button", { name: "Ver sortimento filtrado" });
    expect(setSellableAssortment).not.toHaveBeenCalled();
  });

  it("does not send include_all in filtered mode (list and search)", async () => {
    const fetch = vi
      .fn()
      .mockResolvedValueOnce(new Response(JSON.stringify(page(1, null, "2026-07-14T10:10:11Z"))))
      .mockResolvedValueOnce(new Response(JSON.stringify(page(11, null, "2026-07-14T10:11:12Z"))));
    const transport = createRefreshableFetch(fetch as unknown as typeof globalThis.fetch);
    const sdk = createMarketplaceCentralClient({ baseUrl: "", fetchImpl: transport.fetchImpl });
    renderPage({ ...sdk, withNoCache: transport.withNoCache });

    expect(await screen.findByText("Product 1")).toBeInTheDocument();
    expect(requestUrl(fetch.mock.calls[0])).not.toContain("include_all");

    fireEvent.change(screen.getByLabelText("Buscar no catálogo"), { target: { value: "CUBA" } });
    expect(await screen.findByText("Product 11")).toBeInTheDocument();
    expect(requestUrl(fetch.mock.calls[1])).toContain("/catalog/products/search?q=CUBA");
    expect(requestUrl(fetch.mock.calls[1])).not.toContain("include_all");
  });

  it("ver todos refetches the list with include_all=true", async () => {
    const fetch = vi
      .fn()
      .mockResolvedValueOnce(new Response(JSON.stringify(page(1, null, "2026-07-14T10:10:11Z"))))
      .mockResolvedValueOnce(new Response(JSON.stringify(page(2, null, "2026-07-14T10:11:12Z"))));
    const transport = createRefreshableFetch(fetch as unknown as typeof globalThis.fetch);
    const sdk = createMarketplaceCentralClient({ baseUrl: "", fetchImpl: transport.fetchImpl });
    renderPage({ ...sdk, withNoCache: transport.withNoCache });

    await screen.findByText("Product 1");
    fireEvent.click(screen.getByRole("button", { name: "Ver todos" }));
    await screen.findByText("Product 2");
    expect(requestUrl(fetch.mock.calls[1])).toContain("include_all=true");
  });

  it("returns to the filtered key without serving the ver-todos page or firing a new fetch", async () => {
    const fetch = vi
      .fn()
      .mockResolvedValueOnce(new Response(JSON.stringify(page(1, null, "2026-07-14T10:10:11Z")))) // filtered
      .mockResolvedValueOnce(new Response(JSON.stringify(page(2, null, "2026-07-14T10:11:12Z")))); // ver todos
    const transport = createRefreshableFetch(fetch as unknown as typeof globalThis.fetch);
    const sdk = createMarketplaceCentralClient({ baseUrl: "", fetchImpl: transport.fetchImpl });
    renderPage({ ...sdk, withNoCache: transport.withNoCache });

    await screen.findByText("Product 1");
    fireEvent.click(screen.getByRole("button", { name: "Ver todos" }));
    await screen.findByText("Product 2");
    expect(fetch).toHaveBeenCalledTimes(2);

    // Same QueryClient, same debounced search: a shared key would either
    // still show Product 2 (the ver-todos page) or force a third fetch.
    fireEvent.click(screen.getByRole("button", { name: "Ver sortimento filtrado" }));
    await screen.findByText("Product 1");
    expect(screen.queryByText("Product 2")).not.toBeInTheDocument();
    expect(fetch).toHaveBeenCalledTimes(2);
  });

  it("shows Sem estoque for quantity 0 in ver-todos mode, and the honest-unknown dash (not the badge) for quantity null", async () => {
    const fetch = vi
      .fn()
      .mockResolvedValueOnce(new Response(JSON.stringify(page(1, null, "2026-07-14T10:10:11Z")))) // filtered mount
      .mockResolvedValueOnce(
        new Response(
          JSON.stringify(
            pageOf(
              [
                { id: 20, quantity: 0 },
                { id: 21, quantity: null, quality: ["missing_stock"] },
              ],
              null,
              "2026-07-14T10:11:12Z",
            ),
          ),
        ),
      );
    const transport = createRefreshableFetch(fetch as unknown as typeof globalThis.fetch);
    const sdk = createMarketplaceCentralClient({ baseUrl: "", fetchImpl: transport.fetchImpl });
    renderPage({ ...sdk, withNoCache: transport.withNoCache });

    await screen.findByText("Product 1");
    fireEvent.click(screen.getByRole("button", { name: "Ver todos" }));

    expect(await screen.findByText("Product 20")).toBeInTheDocument();
    expect(screen.getByText("Sem estoque")).toBeInTheDocument();
    expect(screen.getByText("Product 21")).toBeInTheDocument();
    expect(screen.getByText("— (missing_stock)")).toBeInTheDocument();
    expect(screen.queryAllByText("Sem estoque")).toHaveLength(1);
  });

  it("hides the counts chip and never calls getCatalogAssortmentCounts when erpSource is absent", async () => {
    const fetch = vi
      .fn()
      .mockResolvedValue(new Response(JSON.stringify(page(1, null, "2026-07-14T10:11:12Z"))));
    const transport = createRefreshableFetch(fetch as unknown as typeof globalThis.fetch);
    const sdk = createMarketplaceCentralClient({ baseUrl: "", fetchImpl: transport.fetchImpl });
    const getCatalogAssortmentCounts = vi
      .fn()
      .mockResolvedValue({ sellable_count: 2, total_count: 4 });
    renderPage(
      { ...sdk, withNoCache: transport.withNoCache, getCatalogAssortmentCounts },
      createWebQueryClient(),
      undefined,
    );

    await screen.findByText("Product 1");
    expect(screen.queryByText(/Vendáveis/)).not.toBeInTheDocument();
    expect(getCatalogAssortmentCounts).not.toHaveBeenCalled();
  });

  it("atualizar in ver-todos mode refetches the live include_all key", async () => {
    const fetch = vi
      .fn()
      .mockResolvedValueOnce(new Response(JSON.stringify(page(1, null, "2026-07-14T10:10:11Z")))) // filtered
      .mockResolvedValueOnce(new Response(JSON.stringify(page(2, null, "2026-07-14T10:11:12Z")))) // ver todos
      .mockResolvedValueOnce(new Response(JSON.stringify(page(3, null, "2026-07-14T10:12:13Z")))); // atualizar
    const transport = createRefreshableFetch(fetch as unknown as typeof globalThis.fetch);
    const sdk = createMarketplaceCentralClient({ baseUrl: "", fetchImpl: transport.fetchImpl });
    renderPage({ ...sdk, withNoCache: transport.withNoCache });

    await screen.findByText("Product 1");
    fireEvent.click(screen.getByRole("button", { name: "Ver todos" }));
    await screen.findByText("Product 2");

    fireEvent.click(screen.getByRole("button", { name: "Atualizar" }));
    await screen.findByText("Product 3");
    expect(fetch).toHaveBeenCalledTimes(3);
  });
});
