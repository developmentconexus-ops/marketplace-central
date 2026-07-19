import { useMemo } from "react";
import type { JSX } from "react";
import { useQuery } from "@tanstack/react-query";
import { ErrorState, LoadingState } from "@marketplace-central/ui";
import {
  createMarketPriceIntelClient,
  type MarketPriceIntelAggregate,
  type MarketPriceIntelMoney,
} from "@marketplace-central/sdk-runtime";

/** Minimal market seam this panel consumes — IC-03 is HTTP-only (no Go port). */
export interface MarketAggregatesClient {
  listMarketAggregates: (codprods: string[]) => Promise<MarketPriceIntelAggregate[]>;
}

/**
 * Resolve the IC-03 base URL from the build env, mirroring the app client's rule.
 * Kept local (not imported from ClientContext) so this panel stays decoupled from
 * the central seam — the app's useClient() carries no market client.
 */
function apiBaseUrl(): string {
  const env = import.meta.env.VITE_API_BASE_URL;
  if (env && env.trim()) return env;
  return import.meta.env.DEV ? "http://localhost:8080" : "";
}

export interface MarketComparisonProps {
  /** ML codprod (string) for the selected product, or null when nothing is picked. */
  productId: string | null;
  /** Injectable for tests; defaults to a standalone IC-03 client (useClient has no market seam). */
  client?: MarketAggregatesClient;
}

/** ADR-17: a missing market price is unknown, never a fabricated 0. */
function Money({ value }: { value: MarketPriceIntelMoney | null }): JSX.Element {
  if (value === null) return <span className="text-faint">—</span>;
  return <span className="font-mono text-ink">R$ {value.amount}</span>;
}

/**
 * MarketComparison shows the ML market evidence for the selected product: the
 * median and lowest valid offer, offer/seller counts, the evidence source, and
 * when it was fetched. INSUFFICIENT_MARKET / NO_PRICE_EVIDENCE surface as-is —
 * the panel never fabricates a price when the market can't support one (ADR-17).
 * It talks to IC-03 over HTTP (there is no Go market port), via a standalone
 * client since the app's useClient() exposes only the central seam.
 */
export function MarketComparison({ productId, client }: MarketComparisonProps): JSX.Element {
  const marketClient = useMemo<MarketAggregatesClient>(
    () =>
      client ?? createMarketPriceIntelClient({ baseUrl: apiBaseUrl() }),
    [client],
  );

  const query = useQuery({
    queryKey: ["market", "aggregates", productId],
    queryFn: () => marketClient.listMarketAggregates([productId as string]),
    enabled: productId !== null,
  });

  if (productId === null) {
    return <p className="text-sm text-muted">Selecione um produto para comparar com o mercado.</p>;
  }
  if (query.isLoading) return <LoadingState />;
  if (query.isError) {
    return <ErrorState onRetry={() => void query.refetch()} detail="Falha ao carregar os dados de mercado." />;
  }

  const aggregate = query.data?.find((a) => a.product_id === productId) ?? query.data?.[0] ?? null;
  if (aggregate === null) {
    return <p className="text-sm text-muted">Sem dados de mercado para este produto.</p>;
  }

  const priced = aggregate.status === "OK";

  return (
    <div data-testid="market-comparison" className="flex flex-col gap-2">
      <div className="flex items-center justify-between">
        <h3 className="text-sm font-semibold text-ink">Comparação de mercado</h3>
        <span
          data-testid="market-status"
          className={`rounded px-1.5 py-0.5 text-xs font-medium ${
            priced ? "bg-accent-soft text-accent-ink" : "bg-warn-soft text-warn"
          }`}
        >
          {aggregate.status}
        </span>
      </div>

      {!priced ? (
        <p role="note" className="text-xs text-muted">
          Mercado sem preço confiável — nenhum valor é inferido.
        </p>
      ) : null}

      <dl className="grid grid-cols-2 gap-x-4 gap-y-1 text-sm">
        <dt className="text-muted">Mediana</dt>
        <dd className="text-right"><Money value={aggregate.median} /></dd>
        <dt className="text-muted">Menor válido</dt>
        <dd className="text-right"><Money value={aggregate.min_valid} /></dd>
        <dt className="text-muted">Ofertas</dt>
        <dd className="text-right font-mono text-ink">{aggregate.n_offers}</dd>
        <dt className="text-muted">Vendedores</dt>
        <dd className="text-right font-mono text-ink">{aggregate.n_sellers}</dd>
        <dt className="text-muted">Fonte</dt>
        <dd className="text-right font-mono text-ink">{aggregate.source}</dd>
        <dt className="text-muted">Coletado em</dt>
        <dd className="text-right font-mono text-muted">{aggregate.fetched_at}</dd>
      </dl>
    </div>
  );
}
