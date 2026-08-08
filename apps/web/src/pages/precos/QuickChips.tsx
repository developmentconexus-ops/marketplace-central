import { useMemo } from "react";
import type { JSX } from "react";
import { useQuery } from "@tanstack/react-query";
import {
  createMarketPriceIntelClient,
  type MarketPriceIntelMoney,
} from "@marketplace-central/sdk-runtime";
import { formatMoneyOr } from "@marketplace-central/ui";
import { apiBaseUrl, type MarketAggregatesClient } from "./MarketComparison";

export interface QuickChipsProps {
  /** ML codprod (string) for the selected product, or null when nothing is picked. */
  productId: string | null;
  /** Green-margin threshold (raw pt-BR/decimal string) the "cobrir X%" chip seeds. */
  coverPct: string;
  /** Seed the sale-price input from a market reference (dot-decimal string). */
  onSeedPreco: (amount: string) => void;
  /** Seed the solver's target margin (raw string). */
  onSeedTarget: (pct: string) => void;
  /** Injectable for tests; defaults to a standalone IC-03 client. */
  client?: MarketAggregatesClient;
}

/**
 * QuickChips are the sim-panel `atalhos` (design §Simulador): one-tap seeds for
 * the working price. "mediana" / "menor conc." pull from the ML market evidence
 * and ONLY render when that value is present — a missing market price is never
 * fabricated into a chip (ADR-17). "cobrir X%" always shows: it seeds the solver
 * target (a client-side pre-fill), never a fabricated price. Market data comes
 * from IC-03 over HTTP via a standalone client (useClient carries no market seam).
 */
export function QuickChips({
  productId,
  coverPct,
  onSeedPreco,
  onSeedTarget,
  client,
}: QuickChipsProps): JSX.Element | null {
  const marketClient = useMemo<MarketAggregatesClient>(
    () => client ?? createMarketPriceIntelClient({ baseUrl: apiBaseUrl() }),
    [client],
  );

  const query = useQuery({
    queryKey: ["market", "aggregates", productId],
    queryFn: () => marketClient.listMarketAggregates([productId as string]),
    enabled: productId !== null,
  });

  const aggregate =
    productId === null
      ? null
      : (query.data?.find((a) => a.product_id === productId) ?? query.data?.[0] ?? null);
  const priced = aggregate?.status === "OK";
  const median = priced ? (aggregate?.median ?? null) : null;
  const minValid = priced ? (aggregate?.min_valid ?? null) : null;

  const chip = (key: string, label: string, onClick: () => void) => (
    <button
      key={key}
      type="button"
      data-testid={`quick-chip-${key}`}
      onClick={onClick}
      className="cursor-pointer whitespace-nowrap rounded-pill border border-border px-2.5 py-0.5 text-[11.5px] text-muted hover:border-accent hover:text-accent-ink"
    >
      {label}
    </button>
  );

  const money = (m: MarketPriceIntelMoney) => formatMoneyOr(m.amount, "—", m.currency);

  return (
    <div data-testid="quick-chips" className="flex flex-wrap gap-1.5">
      {median
        ? chip("mediana", `mediana ${money(median)}`, () => onSeedPreco(median.amount))
        : null}
      {minValid
        ? chip("menor", `menor conc. ${money(minValid)}`, () => onSeedPreco(minValid.amount))
        : null}
      {chip("cobrir", `cobrir ${coverPct}%`, () => onSeedTarget(coverPct))}
    </div>
  );
}
