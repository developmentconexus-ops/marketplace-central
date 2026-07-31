import { useEffect, useMemo, useState } from "react";
import { useQueryClient } from "@tanstack/react-query";
import { Button, formatMoneyOr, PaginatedTable, SurfaceCard } from "@marketplace-central/ui";
import type { ActiveSourceName, CatalogProductFact } from "@marketplace-central/sdk-runtime";
import {
  catalogQueryKeys,
  FreshnessIndicator,
  useCatalogAssortmentCountsQuery,
} from "@marketplace-central/web-query";
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

// A quantity of 0 is a real read (there is nothing to sell); a null quantity
// is an unresolved read (ADR-17: unknown never becomes zero). Only the first
// gets the badge — the second keeps its honest-unknown dash.
function StockCell({ quantity, quality }: { quantity: number | null; quality: string[] }) {
  if (quantity != null && quantity <= 0) {
    return (
      <span className="inline-flex items-center px-2 py-0.5 rounded-pill text-xs font-medium bg-warn-soft text-warn">
        Sem estoque
      </span>
    );
  }
  return <>{factQuantity(quantity, quality)}</>;
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
  const [includeAll, setIncludeAll] = useState(false);

  useEffect(() => {
    const timer = window.setTimeout(() => setDebouncedSearch(search.trim()), 300);
    return () => window.clearTimeout(timer);
  }, [search]);

  const factsQuery = useCatalogFactsQuery(client, !debouncedSearch, erpSource, includeAll);
  const searchQuery = useCatalogSearchQuery(client, debouncedSearch, erpSource, includeAll);
  // The chip is presentation only: read-only, the sellable rule itself is
  // configured on /integracoes. Nothing on this screen writes it.
  const countsQuery = useCatalogAssortmentCountsQuery(client, erpSource);
  // Both reads page the same way, so the screen drives whichever one is active
  // through one set of controls. Search used to be capped at its first page with
  // no "Carregar mais", which hid every match past the 50th.
  const activeQuery = debouncedSearch ? searchQuery : factsQuery;
  const pages = activeQuery.data?.pages ?? [];
  const facts = useMemo(() => pages.flatMap((page) => page.items), [pages]);
  const isLoading = activeQuery.isPending;
  const error = activeQuery.error;
  const asOf = pages[pages.length - 1]?.as_of;

  async function refresh() {
    setRefreshing(true);
    try {
      const run = async () => {
        if (debouncedSearch) {
          await queryClient.refetchQueries({
            queryKey: catalogQueryKeys.search(debouncedSearch, { erp_source: erpSource ?? null, include_all: includeAll }),
          });
        } else {
          await queryClient.refetchQueries({
            queryKey: catalogQueryKeys.facts({ limit: 50, erp_source: erpSource ?? null, include_all: includeAll }),
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
        <div className="flex items-center gap-3">
          {countsQuery.data ? (
            <span className="text-xs text-faint">
              Vendáveis {countsQuery.data.sellable_count} de {countsQuery.data.total_count}
            </span>
          ) : null}
          <Button
            variant="secondary"
            onClick={() => setIncludeAll((current) => !current)}
          >
            {includeAll ? "Ver sortimento filtrado" : "Ver todos"}
          </Button>
          <Button variant="secondary" disabled={refreshing || isLoading} loading={refreshing} onClick={() => void refresh()}>
            Atualizar
          </Button>
        </div>
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
                <StockCell quantity={fact.sellable_stock.quantity} quality={fact.sellable_stock.quality} />
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

      {activeQuery.hasNextPage ? (
        <Button
          variant="secondary"
          disabled={activeQuery.isFetchingNextPage}
          loading={activeQuery.isFetchingNextPage}
          onClick={() => void activeQuery.fetchNextPage()}
        >
          Carregar mais
        </Button>
      ) : null}
      {!activeQuery.hasNextPage && facts.length > 0 ? (
        <p className="text-sm text-muted">
          {debouncedSearch ? "Fim dos resultados" : "Fim da lista"}
        </p>
      ) : null}
    </div>
  );
}
