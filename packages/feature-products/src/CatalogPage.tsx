import { useEffect, useMemo, useState } from "react";
import { useQueryClient } from "@tanstack/react-query";
import { Button, formatMoneyOr, PaginatedTable, SurfaceCard } from "@marketplace-central/ui";
import type { ActiveSourceName, CatalogProductFact } from "@marketplace-central/sdk-runtime";
import { catalogQueryKeys, FreshnessIndicator } from "@marketplace-central/web-query";
import { useCatalogFactsQuery, useCatalogSearchQuery, type CatalogQueriesClient } from "./catalogQueries";

export interface CatalogPageProps {
  client: CatalogQueriesClient;
  /**
   * The tenant's active source, resolved by the app shell from
   * GET /config/active-source. It never travels in the request (the server
   * resolves it itself); it only partitions the cache so a flip cannot serve
   * the previous source's rows.
   */
  erpSource?: ActiveSourceName;
}

// The reader reports WHY a value is absent (missing_cost, stale, …). Showing
// the reason next to the dash is the honest-unknown contract (ADR-17): never a
// fabricated 0, and never a bare blank the operator has to guess about.
function qualityText(quality: string[]): string {
  return quality.length ? quality.join(", ") : "atual";
}

function factAmount(amount: string | null, quality: string[]): string {
  return amount == null ? `— (${qualityText(quality)})` : formatMoneyOr(amount);
}

function factQuantity(quantity: number | null, quality: string[]): string {
  return quantity == null ? `— (${qualityText(quality)})` : String(quantity);
}

function errorMessage(error: unknown): string {
  if (error && typeof error === "object") {
    const structured = error as { error?: { message?: string }; message?: string };
    return structured.error?.message ?? structured.message ?? "Não foi possível carregar o catálogo.";
  }
  return "Não foi possível carregar o catálogo.";
}

export function CatalogPage({ client, erpSource }: CatalogPageProps) {
  const queryClient = useQueryClient();
  const [search, setSearch] = useState("");
  const [debouncedSearch, setDebouncedSearch] = useState("");
  const [refreshing, setRefreshing] = useState(false);

  useEffect(() => {
    const timer = window.setTimeout(() => setDebouncedSearch(search.trim()), 300);
    return () => window.clearTimeout(timer);
  }, [search]);

  const factsQuery = useCatalogFactsQuery(client, !debouncedSearch, erpSource);
  const searchQuery = useCatalogSearchQuery(client, debouncedSearch, erpSource);
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
          await queryClient.refetchQueries({
            queryKey: catalogQueryKeys.search(debouncedSearch, { erp_source: erpSource ?? null }),
          });
        } else {
          await queryClient.refetchQueries({
            queryKey: catalogQueryKeys.facts({ limit: 50, erp_source: erpSource ?? null }),
          });
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
          <p className="text-xs uppercase tracking-[0.2em] text-faint">Produtos do ERP</p>
          <h1 className="text-2xl font-semibold text-ink">Catálogo</h1>
          <p className="text-sm text-muted">
            <FreshnessIndicator asOf={asOf} />
          </p>
        </div>
        <Button variant="secondary" disabled={refreshing || isLoading} loading={refreshing} onClick={() => void refresh()}>
          Atualizar
        </Button>
      </div>

      <SurfaceCard className="rounded-card">
        <label className="block text-sm text-ink">
          <span className="mb-1 block font-medium">Buscar no catálogo</span>
          <input
            aria-label="Buscar no catálogo"
            className="w-full rounded-control border border-border bg-surface px-3 py-2 text-ink"
            placeholder="Referência, descrição, EAN..."
            value={search}
            onChange={(event) => setSearch(event.target.value)}
          />
        </label>
      </SurfaceCard>

      {isLoading ? (
        <SurfaceCard>
          <p className="text-sm text-muted">Carregando catálogo...</p>
        </SurfaceCard>
      ) : null}
      {error ? (
        <SurfaceCard>
          <div className="space-y-3">
            <p className="text-sm text-warn">{errorMessage(error)}</p>
            <Button variant="primary" onClick={() => void refresh()}>
              Tentar novamente
            </Button>
          </div>
        </SurfaceCard>
      ) : null}
      {!isLoading && !error && facts.length === 0 ? (
        <SurfaceCard>
          <p className="text-sm text-muted">Nenhum produto encontrado.</p>
        </SurfaceCard>
      ) : null}

      {!isLoading && !error && facts.length > 0 ? (
        <PaginatedTable
          items={facts}
          pageSize={50}
          renderHeader={() => (
            <tr>
              <th className="px-4 py-3 font-medium text-muted">Referência</th>
              <th className="px-4 py-3 font-medium text-muted">Descrição</th>
              <th className="px-4 py-3 font-medium text-muted">EAN</th>
              <th className="px-4 py-3 font-medium text-muted">Ativo</th>
              <th className="px-4 py-3 text-right font-medium text-muted">Estoque</th>
              <th className="px-4 py-3 text-right font-medium text-muted">Preço</th>
              <th className="px-4 py-3 text-right font-medium text-muted">Custo</th>
            </tr>
          )}
          renderRow={(fact: CatalogProductFact) => (
            <tr key={fact.internal_product_id} className="border-b border-border">
              <td className="px-4 py-3 font-mono text-xs text-ink">{fact.reference ?? "—"}</td>
              <td className="px-4 py-3 text-ink">{fact.description ?? "—"}</td>
              <td className="px-4 py-3 text-muted">{fact.ean ?? "—"}</td>
              <td className="px-4 py-3 text-muted">{fact.active ? "Sim" : "Não"}</td>
              <td className="px-4 py-3 text-right text-muted">
                {factQuantity(fact.sellable_stock.quantity, fact.sellable_stock.quality)}
              </td>
              <td className="px-4 py-3 text-right font-mono text-muted">
                {factAmount(fact.current_price.amount, fact.current_price.quality)}
              </td>
              <td className="px-4 py-3 text-right font-mono text-muted">
                {factAmount(fact.cost.amount, fact.cost.quality)}
              </td>
            </tr>
          )}
          emptyState={<p>Nenhum produto encontrado.</p>}
        />
      ) : null}

      {!debouncedSearch && factsQuery.hasNextPage ? (
        <Button
          variant="secondary"
          disabled={factsQuery.isFetchingNextPage}
          loading={factsQuery.isFetchingNextPage}
          onClick={() => void factsQuery.fetchNextPage()}
        >
          Carregar mais
        </Button>
      ) : null}
      {!debouncedSearch && !factsQuery.hasNextPage && facts.length > 0 ? (
        <p className="text-sm text-muted">Fim da lista</p>
      ) : null}
      {debouncedSearch && searchQuery.data ? (
        <p className="text-sm text-muted">A busca mostra apenas a primeira página de resultados.</p>
      ) : null}
    </div>
  );
}
