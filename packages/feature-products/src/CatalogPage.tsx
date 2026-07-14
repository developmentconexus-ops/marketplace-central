import { useEffect, useMemo, useState } from "react";
import { useQueryClient } from "@tanstack/react-query";
import { Button, PaginatedTable, SurfaceCard } from "@marketplace-central/ui";
import type { CatalogProductFact } from "@marketplace-central/sdk-runtime";
import {
  catalogQueryKeys,
  FreshnessIndicator,
} from "@marketplace-central/web-query";
import { useCatalogFactsQuery, useCatalogSearchQuery, type CatalogQueriesClient } from "./catalogQueries";

export interface CatalogPageProps {
  client: CatalogQueriesClient;
}

function qualityText(quality: string[]): string {
  return quality.length ? quality.join(", ") : "quality: current";
}

function factAmount(amount: string | null, quality: string[]): string {
  return amount == null ? `unknown (${qualityText(quality)})` : amount;
}

function factQuantity(quantity: number | null, quality: string[]): string {
  return quantity == null ? `unknown (${qualityText(quality)})` : String(quantity);
}

function errorMessage(error: unknown): string {
  if (error && typeof error === "object") {
    const structured = error as { error?: { message?: string }; message?: string };
    return structured.error?.message ?? structured.message ?? "Failed to load catalog.";
  }
  return "Failed to load catalog.";
}

export function CatalogPage({ client }: CatalogPageProps) {
  const queryClient = useQueryClient();
  const [search, setSearch] = useState("");
  const [debouncedSearch, setDebouncedSearch] = useState("");
  const [refreshing, setRefreshing] = useState(false);

  useEffect(() => {
    const timer = window.setTimeout(() => setDebouncedSearch(search.trim()), 300);
    return () => window.clearTimeout(timer);
  }, [search]);

  const factsQuery = useCatalogFactsQuery(client, !debouncedSearch);
  const searchQuery = useCatalogSearchQuery(client, debouncedSearch);
  const pages = factsQuery.data?.pages ?? [];
  const facts = useMemo(
    () => (debouncedSearch ? searchQuery.data?.items ?? [] : pages.flatMap((page) => page.items)),
    [debouncedSearch, pages, searchQuery.data?.items],
  );
  const activePage = debouncedSearch ? searchQuery.data : pages[pages.length - 1];
  const isLoading = debouncedSearch ? searchQuery.isPending : factsQuery.isPending;
  const error = debouncedSearch ? searchQuery.error : factsQuery.error;
  const asOf = activePage?.as_of;

  async function refresh() {
    setRefreshing(true);
    try {
      const run = async () => {
        if (debouncedSearch) {
          await queryClient.refetchQueries({ queryKey: catalogQueryKeys.search(debouncedSearch) });
        } else {
          await queryClient.refetchQueries({ queryKey: catalogQueryKeys.facts({ limit: 50 }) });
        }
      };
      await (client.withNoCache ? client.withNoCache(run) : run());
    } finally {
      setRefreshing(false);
    }
  }

  return (
    <div className="space-y-5">
      <div className="flex flex-col gap-3 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <p className="text-xs uppercase tracking-[0.2em] text-slate-500">Oracle catalog</p>
          <h1 className="text-2xl font-semibold text-slate-900">Catalog</h1>
          <p className="text-sm text-slate-500"><FreshnessIndicator asOf={asOf} /></p>
        </div>
        <Button variant="secondary" disabled={refreshing || isLoading} loading={refreshing} onClick={() => void refresh()}>
          Refresh
        </Button>
      </div>

      <SurfaceCard className="rounded-[24px]">
        <label className="block text-sm text-slate-700">
          <span className="mb-1 block font-medium">Search catalog</span>
          <input
            aria-label="Search catalog"
            className="w-full rounded-xl border border-slate-200 px-3 py-2"
            placeholder="Reference, description, EAN..."
            value={search}
            onChange={(event) => setSearch(event.target.value)}
          />
        </label>
      </SurfaceCard>

      {isLoading ? <SurfaceCard><p className="text-sm text-slate-600">Loading catalog...</p></SurfaceCard> : null}
      {error ? <SurfaceCard><div className="space-y-3"><p className="text-sm text-red-700">{errorMessage(error)}</p><Button variant="primary" onClick={() => void refresh()}>Retry</Button></div></SurfaceCard> : null}
      {!isLoading && !error && facts.length === 0 ? <SurfaceCard><p className="text-sm text-slate-600">No catalog facts found.</p></SurfaceCard> : null}

      {!isLoading && !error && facts.length > 0 ? (
        <PaginatedTable
          items={facts}
          pageSize={50}
          renderHeader={() => (
            <tr>
              <th className="px-4 py-3 font-medium text-slate-600">Reference</th>
              <th className="px-4 py-3 font-medium text-slate-600">Description</th>
              <th className="px-4 py-3 font-medium text-slate-600">EAN</th>
              <th className="px-4 py-3 font-medium text-slate-600">Active</th>
              <th className="px-4 py-3 font-medium text-slate-600 text-right">Stock</th>
              <th className="px-4 py-3 font-medium text-slate-600 text-right">Price</th>
              <th className="px-4 py-3 font-medium text-slate-600 text-right">Cost</th>
            </tr>
          )}
          renderRow={(fact: CatalogProductFact) => (
            <tr key={fact.internal_product_id} className="border-b border-slate-50">
              <td className="px-4 py-3 font-mono text-xs text-slate-700">{fact.reference ?? "unknown"}</td>
              <td className="px-4 py-3 text-slate-900">{fact.description ?? "unknown"}</td>
              <td className="px-4 py-3 text-slate-600">{fact.ean ?? "unknown"}</td>
              <td className="px-4 py-3 text-slate-600">{fact.active ? "Yes" : "No"}</td>
              <td className="px-4 py-3 text-right text-slate-600">{factQuantity(fact.sellable_stock.quantity, fact.sellable_stock.quality)}</td>
              <td className="px-4 py-3 text-right font-mono text-slate-600">{factAmount(fact.current_price.amount, fact.current_price.quality)}</td>
              <td className="px-4 py-3 text-right font-mono text-slate-600">{factAmount(fact.cost.amount, fact.cost.quality)}</td>
            </tr>
          )}
          emptyState={<p>No catalog facts found.</p>}
        />
      ) : null}

      {!debouncedSearch && factsQuery.hasNextPage ? (
        <Button variant="secondary" disabled={factsQuery.isFetchingNextPage} loading={factsQuery.isFetchingNextPage} onClick={() => void factsQuery.fetchNextPage()}>
          Load next page
        </Button>
      ) : null}
      {!debouncedSearch && !factsQuery.hasNextPage && facts.length > 0 ? <p className="text-sm text-slate-500">Fim da lista</p> : null}
      {debouncedSearch && searchQuery.data ? <p className="text-sm text-slate-500">Search results are bounded to one page.</p> : null}
    </div>
  );
}
