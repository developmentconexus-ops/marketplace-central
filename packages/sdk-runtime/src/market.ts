// F-04-S5: standalone price-intel collect/verdict SDK surface.
//
// Per D-F4-o this file stays independent of ./index.ts (the hub-owned
// barrel). It defines its own client, error, and money/aggregate types
// rather than importing from index.ts, so a later `export * from './market'`
// in the barrel cannot collide on identifier names. All types here are
// prefixed MarketPriceIntel* for that reason.

export interface MarketPriceIntelMoney {
  amount: string;
  currency: "BRL";
}

export interface MarketPriceIntelPosition {
  rank: number;
  total: number;
}

export interface MarketPriceIntelApiError {
  status: number;
  error: {
    code: string;
    message: string;
    details?: Record<string, unknown>;
  };
}

// --- POST /market/collections (B.1) ---------------------------------------

export interface MarketPriceIntelCollectionRequest {
  codprod: string;
}

// Active ERP snapshot selector for collection identity resolution. Declared
// locally (D-F4-o) rather than imported from ./index.ts.
export type MarketPriceIntelErpSource = "xlsx" | "catalogo_cliente";

export type MarketPriceIntelMatchStatus = "ACCEPT" | "REVIEW" | "REJECT" | "NO_CANDIDATE";
export type MarketPriceIntelPriceEvidenceStatus = "OK" | "INSUFFICIENT_MARKET" | "NO_PRICE_EVIDENCE";
export type MarketPriceIntelBlockingState =
  | "NO_CANDIDATE"
  | "NO_PRICE_EVIDENCE"
  | "INSUFFICIENT_MARKET"
  | "SEM_CUSTO"
  | null;
export type MarketPriceIntelCollectionCauseReason =
  | "NO_IDENTITY"
  | "FLAG_DISABLED"
  | "PROVIDER_4XX"
  | "PROVIDER_5XX"
  | "TIMEOUT"
  | "PRICE_UNAVAILABLE";
export type MarketPriceIntelCollectionStatus = "COMPLETED" | "PARTIAL";

export interface MarketPriceIntelCollectionDecision {
  codprod: string;
  match_status: MarketPriceIntelMatchStatus;
  price_evidence_status: MarketPriceIntelPriceEvidenceStatus;
  blocking_state: MarketPriceIntelBlockingState;
}

export interface MarketPriceIntelCollectionCause {
  codprod: string;
  reason: MarketPriceIntelCollectionCauseReason;
  detail: string | null;
}

export interface MarketPriceIntelCollectionContagens {
  ok: number;
  no_price_evidence: number;
  insufficient_market: number;
  no_candidate: number;
  sem_custo: number;
}

// Field name "decisões" is verbatim per the ratified IC-06/IC-03 amendment.
export interface MarketPriceIntelCollectionResponse {
  status: MarketPriceIntelCollectionStatus;
  "decisões": MarketPriceIntelCollectionDecision[];
  contagens: MarketPriceIntelCollectionContagens;
  causas: MarketPriceIntelCollectionCause[];
}

// --- GET /market/signals ----------------------------------------------------

export interface MarketPriceIntelSignal {
  listing_id: string;
  our_price: MarketPriceIntelMoney | null;
  winner_price: MarketPriceIntelMoney | null;
  target_price: MarketPriceIntelMoney | null;
  position: MarketPriceIntelPosition | null;
  fetched_at: string;
}

// --- GET /market/aggregates --------------------------------------------------

export type MarketPriceIntelSource = "ml_sale_price" | "ml_price_to_win" | "ml_catalog_offers";
export type MarketPriceIntelAggregateStatus = "OK" | "INSUFFICIENT_MARKET" | "NO_PRICE_EVIDENCE";

export interface MarketPriceIntelAggregate {
  product_id: string;
  median: MarketPriceIntelMoney | null;
  min_valid: MarketPriceIntelMoney | null;
  n_offers: number;
  n_sellers: number;
  source: MarketPriceIntelSource;
  fetched_at: string;
  computed_at: string;
  status: MarketPriceIntelAggregateStatus;
}

// --- GET /market/verdicts (B.2) ----------------------------------------------

export interface MarketPriceIntelVerdictInputRef {
  source: string;
  as_of: string;
}

export interface MarketPriceIntelVerdictMarketRange {
  min_valid: MarketPriceIntelMoney | null;
  median: MarketPriceIntelMoney | null;
  currency: "BRL";
  n_offers: number;
  n_sellers: number;
}

export interface MarketPriceIntelVerdict {
  match_status: MarketPriceIntelMatchStatus;
  price_evidence_status: MarketPriceIntelPriceEvidenceStatus;
  // Always null from M-02: margin labels are M-07-owned.
  verdict_label: string | null;
  blocking_state: MarketPriceIntelBlockingState;
  inputs_used: Record<string, MarketPriceIntelVerdictInputRef>;
  market_range: MarketPriceIntelVerdictMarketRange | null;
}

// --- client -------------------------------------------------------------

export function createMarketPriceIntelClient(options: {
  baseUrl: string;
  fetchImpl?: typeof fetch;
}) {
  const fetchImpl = options.fetchImpl ?? fetch;

  async function getJson<T>(path: string): Promise<T> {
    const response = await fetchImpl(`${options.baseUrl}${path}`, { method: "GET" });
    const data = await response.json();
    if (!response.ok) {
      throw { status: response.status, error: data.error } satisfies MarketPriceIntelApiError;
    }
    return data as T;
  }

  async function postJson<T>(path: string, body: unknown): Promise<T> {
    const response = await fetchImpl(`${options.baseUrl}${path}`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(body),
    });
    const data = await response.json();
    if (!response.ok) {
      throw { status: response.status, error: data.error } satisfies MarketPriceIntelApiError;
    }
    return data as T;
  }

  function idsQuery(key: string, ids: string[]): string {
    const query = new URLSearchParams();
    if (ids.length > 0) query.set(key, ids.join(","));
    const encoded = query.toString();
    return encoded ? `?${encoded}` : "";
  }

  return {
    // Single-product synchronous collection (D-F4-p): no batch POST.
    // erpSource selects which imported snapshot resolves the codprod identity
    // (D-F4-o keeps this file independent of ./index.ts, so the source union is
    // declared locally). Absent = backend default (xlsx), URL byte-stable.
    collectMarketPriceIntel: (codprod: string, erpSource?: MarketPriceIntelErpSource) => {
      const body: MarketPriceIntelCollectionRequest = { codprod };
      const query = erpSource ? `?erp_source=${encodeURIComponent(erpSource)}` : "";
      return postJson<MarketPriceIntelCollectionResponse>(`/market/collections${query}`, body);
    },
    listMarketSignals: (listingIds: string[]) => {
      return getJson<MarketPriceIntelSignal[]>(`/market/signals${idsQuery("listing_ids", listingIds)}`);
    },
    listMarketAggregates: (codprods: string[]) => {
      return getJson<MarketPriceIntelAggregate[]>(`/market/aggregates${idsQuery("codprod", codprods)}`);
    },
    // A never-collected codprod still yields a 200 NO_PRICE_EVIDENCE verdict,
    // never a 404 (D-F4-q).
    listMarketVerdicts: (codprods: string[]) => {
      return getJson<MarketPriceIntelVerdict[]>(`/market/verdicts${idsQuery("codprod", codprods)}`);
    },
  };
}
