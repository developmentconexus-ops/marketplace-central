import { useMemo } from "react";
import type { JSX } from "react";
import { useQueries, useQuery } from "@tanstack/react-query";
import { MarginChip, UnknownValue } from "@marketplace-central/ui";
import {
  createMarketPriceIntelClient,
  type CatalogProductFact,
  type MarketPriceIntelAggregate,
  type PricingCalcProfile,
} from "@marketplace-central/sdk-runtime";
import { useClient } from "../../app/ClientContext";
import { apiBaseUrl, type MarketAggregatesClient } from "./MarketComparison";
import type { ModalidadeKey } from "./PricingPage";

export interface PricingMatrixProps {
  products: CatalogProductFact[];
  /** Selected row (page owns selection); drives row highlight / aria-pressed. */
  selectedId: number | null;
  onSelect: (id: number) => void;
  modalidade: ModalidadeKey;
  /** Operator thresholds for the margin chip bands. */
  profile: PricingCalcProfile;
  /** Injectable for tests; defaults to a standalone IC-03 client (useClient has no market seam). */
  marketClient?: MarketAggregatesClient;
}

/** Price-evidence verdict shown in the VEREDICTO column — mapped from the market
 * aggregate status ONLY. The categorical MARGIN verdict (saudável/apertado/…) is
 * M-07-owned and deliberately NOT invented here. A missing aggregate is the honest
 * closest state: SEM_EVIDENCIA. */
function veredictoFor(agg: MarketPriceIntelAggregate | undefined): string {
  switch (agg?.status) {
    case "OK":
      return "OK";
    case "INSUFFICIENT_MARKET":
      return "MERCADO_INSUFICIENTE";
    case "NO_PRICE_EVIDENCE":
    default:
      return "SEM_EVIDENCIA";
  }
}

function productSku(p: CatalogProductFact): string {
  return p.manufacturer_reference ?? `#${p.internal_product_id}`;
}

function productDescription(p: CatalogProductFact): string {
  return p.description ?? p.manufacturer_reference ?? `#${p.internal_product_id}`;
}

/** A decimal-string money value, or the honest "—" when absent (ADR-17: never 0). */
function Money({ amount }: { amount: string | null }): JSX.Element {
  if (amount === null) return <UnknownValue />;
  return <span className="font-mono text-ink">R$ {amount}</span>;
}

/**
 * PricingMatrix is the /precos main surface: the design's product matrix
 * (SKU · DESCRIÇÃO · CUSTO · NOSSO PREÇO · PREÇO MERCADO · MARGEM · VEREDICTO).
 * It reuses the SAME per-product seams the single-product panel uses — CUSTO and
 * NOSSO PREÇO come straight off the CatalogProductFact, market evidence via ONE
 * listMarketAggregates call for every codprod, and per-row margin via the existing
 * pricingDecompose endpoint (no new endpoints; comissao_pct omitted so the backend
 * resolver chain runs). Row-click hands selection back to the page, which opens the
 * existing simular panel. Unknowns stay honest ("—"), never fabricated zeros.
 */
export function PricingMatrix({
  products,
  selectedId,
  onSelect,
  modalidade,
  profile,
  marketClient,
}: PricingMatrixProps): JSX.Element {
  const client = useClient();
  const resolvedMarketClient = useMemo<MarketAggregatesClient>(
    () => marketClient ?? createMarketPriceIntelClient({ baseUrl: apiBaseUrl() }),
    [marketClient],
  );

  const codprods = useMemo(() => products.map((p) => String(p.internal_product_id)), [products]);

  // Lane A — market evidence for the WHOLE page in one call.
  const marketQuery = useQuery({
    queryKey: ["market", "matrix-aggregates", codprods.join(",")],
    queryFn: () => resolvedMarketClient.listMarketAggregates(codprods),
    enabled: codprods.length > 0,
  });
  const marketByCodprod = useMemo(() => {
    const map = new Map<string, MarketPriceIntelAggregate>();
    for (const agg of marketQuery.data ?? []) map.set(agg.product_id, agg);
    return map;
  }, [marketQuery.data]);

  // Lane B — per-row margin via the existing decompose endpoint. comissao_pct is
  // OMITTED so the backend resolver chain (COTACAO → PADRAO) runs, never a MANUAL
  // override. Rows without a current price cannot be decomposed → margin unknown.
  const decomposeResults = useQueries({
    queries: products.map((p) => ({
      queryKey: ["pricing", "matrix-decompose", p.internal_product_id, p.current_price.amount, modalidade],
      queryFn: () =>
        client.pricingDecompose({
          preco: p.current_price.amount as string,
          modalidade,
          product_id: p.internal_product_id,
        }),
      enabled: p.current_price.amount !== null,
    })),
  });

  const thresholds = {
    healthy: Number(profile.limiar_verde_pct),
    tight: Number(profile.limiar_amarelo_pct),
  };

  return (
    <div data-testid="pricing-matrix" className="overflow-x-auto rounded-lg border border-border bg-surface">
      <table className="w-full min-w-[900px] border-collapse text-left">
        <thead>
          <tr className="bg-surface-2 text-[11px] uppercase tracking-[0.04em] text-muted">
            <Th>SKU</Th>
            <Th>DESCRIÇÃO</Th>
            <Th align="right">CUSTO</Th>
            <Th align="right">NOSSO PREÇO</Th>
            <Th align="right">PREÇO MERCADO</Th>
            <Th align="right">MARGEM</Th>
            <Th>VEREDICTO</Th>
          </tr>
        </thead>
        <tbody className="text-[12.5px]">
          {products.map((p, i) => {
            const agg = marketByCodprod.get(String(p.internal_product_id));
            const decomposition = decomposeResults[i]?.data?.decomposition ?? null;
            const margemValor = decomposition?.margem_valor ?? null;
            const margemPct = decomposition?.margem_pct ?? null;
            const pctNum = margemPct !== null && Number.isFinite(Number(margemPct)) ? Number(margemPct) : null;
            const priced = agg?.status === "OK";
            const isSelected = selectedId === p.internal_product_id;

            return (
              <tr
                key={p.internal_product_id}
                data-testid={`matrix-row-${p.internal_product_id}`}
                onClick={() => onSelect(p.internal_product_id)}
                aria-selected={isSelected}
                className={`cursor-pointer border-t border-border-2 ${
                  isSelected ? "bg-accent-soft" : "hover:bg-surface-2"
                }`}
              >
                <Td>
                  <span className="font-mono text-ink">{productSku(p)}</span>
                </Td>
                <Td>
                  <span className="text-ink">{productDescription(p)}</span>
                </Td>
                <Td align="right" testId={`matrix-custo-${p.internal_product_id}`}>
                  <Money amount={p.cost.amount} />
                </Td>
                <Td align="right">
                  <Money amount={p.current_price.amount} />
                </Td>
                <Td align="right" testId={`matrix-mercado-${p.internal_product_id}`}>
                  {/* Rank "19º/23" is absent from aggregate data → price only, no rank. */}
                  {priced ? <Money amount={agg?.median?.amount ?? null} /> : <UnknownValue />}
                </Td>
                <Td align="right" testId={`matrix-margem-${p.internal_product_id}`}>
                  {margemValor !== null ? (
                    <span className="flex items-center justify-end gap-1.5">
                      <span className="font-mono text-ink">{margemValor}</span>
                      <MarginChip marginPct={pctNum} thresholds={thresholds} />
                    </span>
                  ) : (
                    <UnknownValue hint="margem: M-07" />
                  )}
                </Td>
                <Td testId={`matrix-veredicto-${p.internal_product_id}`}>
                  <span
                    className={`inline-flex items-center rounded px-1.5 py-0.5 text-xs font-medium ${
                      priced ? "bg-accent-soft text-accent-ink" : "bg-warn-soft text-warn"
                    }`}
                  >
                    {veredictoFor(agg)}
                  </span>
                </Td>
              </tr>
            );
          })}
        </tbody>
      </table>
    </div>
  );
}

function Th({ children, align }: { children: React.ReactNode; align?: "right" }): JSX.Element {
  return (
    <th scope="col" className={`px-3 py-2 font-medium ${align === "right" ? "text-right" : ""}`}>
      {children}
    </th>
  );
}

function Td({
  children,
  align,
  testId,
}: {
  children: React.ReactNode;
  align?: "right";
  testId?: string;
}): JSX.Element {
  return (
    <td data-testid={testId} className={`px-3 py-2 ${align === "right" ? "text-right" : ""}`}>
      {children}
    </td>
  );
}
