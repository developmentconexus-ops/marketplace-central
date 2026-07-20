import type {
  CatalogProductFact,
  MarketPriceIntelAggregate,
  MarketPriceIntelMoney,
  MarketPriceIntelPriceEvidenceStatus,
  MarketPriceIntelVerdict,
} from "@marketplace-central/sdk-runtime";

/** A normalized Oportunidades row: ERP catalog fact joined with ML market evidence. */
export interface OppRow {
  sku: string;
  name: string | null;
  costAmount: string | null;
  median: MarketPriceIntelMoney | null;
  /** Distinct competing sellers (deduped by seller) — the CONCORRENTES column + sort key. */
  nSellers: number | null;
  /** Drives the honest evidence note; the recommendation label itself is M-07-owned. */
  evidenceState: MarketPriceIntelPriceEvidenceStatus | null;
}

function codprod(fact: CatalogProductFact): string {
  return String(fact.internal_product_id);
}

/**
 * Split an array into ordered chunks of at most `size`. Used to keep the /market batch
 * reads under the server's MaxReadIDs=200 cap while preserving codprod order — so the
 * positional fact↔verdict alignment in buildOppRows survives chunking (concatenating the
 * per-chunk results in chunk order reconstructs the original fact ordering).
 */
export function chunk<T>(arr: T[], size: number): T[][] {
  if (size <= 0) throw new Error("chunk size must be > 0");
  const out: T[][] = [];
  for (let i = 0; i < arr.length; i += size) out.push(arr.slice(i, i + size));
  return out;
}

/**
 * Cross the ERP catalog facts with ML aggregates + verdicts into Oportunidades rows.
 * Only products with an OK market aggregate (real observed demand) are surfaced — the
 * radar never fabricates a market for a product that has none (ADR-17). Rows are ordered
 * by observed competition (n_sellers — the same distinct-seller count shown in the
 * CONCORRENTES column), a backed signal; the true opportunity ranking is margin-driven and
 * M-07-owned, so it is NOT computed here. VENDAS LÍDER 30D, MARGEM EST. and the
 * recommendation label are not backed by any endpoint and stay honest "—".
 */
export function buildOppRows(
  facts: CatalogProductFact[],
  aggregates: MarketPriceIntelAggregate[],
  verdicts: MarketPriceIntelVerdict[],
): OppRow[] {
  const aggByCodprod = new Map(aggregates.map((a) => [a.product_id, a]));
  const verdictByCodprod = new Map<string, MarketPriceIntelVerdict>();
  // verdicts carry no product_id field, so align them positionally with the codprods
  // requested — the caller passes verdicts fetched for the same fact ordering.
  facts.forEach((fact, i) => {
    const v = verdicts[i];
    if (v) verdictByCodprod.set(codprod(fact), v);
  });

  const rows: OppRow[] = [];
  for (const fact of facts) {
    const id = codprod(fact);
    const agg = aggByCodprod.get(id) ?? null;
    // Surface only products the market actually supports a price for (status OK).
    if (!agg || agg.status !== "OK" || agg.median === null) continue;
    const verdict = verdictByCodprod.get(id) ?? null;
    const costAmount = fact.cost.amount ?? null;
    rows.push({
      sku: id,
      name: fact.description ?? null,
      costAmount,
      median: agg.median,
      nSellers: agg.n_sellers,
      evidenceState:
        verdict?.price_evidence_status ?? (agg.status as MarketPriceIntelPriceEvidenceStatus),
    });
  }

  // Order by observed competition (n_sellers desc — the CONCORRENTES value shown to the
  // user, so the "ordenado por concorrência observada" footer is literally true); tie-break
  // by SKU for stability. No margin sort — the opportunity ranking is M-07-owned (ADR-17).
  rows.sort((a, b) => (b.nSellers ?? -Infinity) - (a.nSellers ?? -Infinity) || a.sku.localeCompare(b.sku));
  return rows;
}
