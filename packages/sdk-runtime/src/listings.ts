// F01-S1: additive competitive-signal typed surface for /listings*.
//
// Re-types List/ByProduct/Summary responses with `market_signal` /
// `signal_status` (per the additive OpenAPI delta landed in the same
// commit) without editing the inline client in ./index.ts. No duplicate
// fetch: the wrapper calls the base client and re-types the already
// enriched JSON payload.

import type {
  ListingGroup,
  ListingGroupPage,
  ListingListOptions,
  ListingPage,
  ListingReadModel,
  ListingSummary,
  ListingSummaryExceptions,
  createMarketplaceCentralClient,
} from "./index";

export type MarketplaceCentralClient = ReturnType<typeof createMarketplaceCentralClient>;

// Listings-owned composite over IC-03 states (not a market enum) — see
// domain invariant in the M-05 slice card: SEM_VINCULO (no codprod link) !=
// NO_PRICE_EVIDENCE (linked, no market evidence) != STALE (evidence older
// than TTL), never collapsed, never fabricated for an unlinked listing.
export type SignalStatus = "OK" | "SEM_VINCULO" | "NO_PRICE_EVIDENCE" | "STALE";

export interface ListingSignalEvidence {
  source: string;
  fetched_at: string;
  freshness: string;
}

export interface ListingMarketSignal {
  status: SignalStatus;
  position: { rank: number; total: number } | null;
  price_to_win: { amount: string; currency: string } | null;
  delta_pct: string | null;
  match_status: "ACCEPT" | "REVIEW" | "REJECT" | "NO_CANDIDATE" | null;
  n_offers: number | null;
  n_sellers: number | null;
  evidence: ListingSignalEvidence | null;
}

export interface ListingWithSignal extends ListingReadModel {
  market_signal: ListingMarketSignal | null;
  signal_status: SignalStatus;
}

export interface ListingPageWithSignals extends Omit<ListingPage, "items"> {
  items: ListingWithSignal[];
}

export interface ListingGroupWithSignal extends Omit<ListingGroup, "listings"> {
  listings: ListingWithSignal[];
}

export interface ListingGroupPageWithSignals extends Omit<ListingGroupPage, "groups"> {
  groups: ListingGroupWithSignal[];
}

export interface ListingSummaryExceptionsWithSignals extends ListingSummaryExceptions {
  sem_vinculo: number;
  abaixo_custo: number;
  sem_evidencia: number;
}

export interface ListingSummaryWithSignals extends Omit<ListingSummary, "exceptions"> {
  exceptions: ListingSummaryExceptionsWithSignals;
}

// Only the three /listings* read methods the wrapper re-types — depending
// on the narrow slice of the client it actually uses (not the full ~60
// method surface) keeps the wrapper and its tests decoupled from unrelated
// endpoints.
export type ListingsSignalClient = Pick<
  MarketplaceCentralClient,
  "listListings" | "listListingsByProduct" | "getListingsSummary"
>;

export function withListingSignals<C extends ListingsSignalClient>(client: C) {
  return {
    ...client,
    listListings: (options: ListingListOptions): Promise<ListingPageWithSignals> =>
      client.listListings(options) as unknown as Promise<ListingPageWithSignals>,
    listListingsByProduct: (options: ListingListOptions): Promise<ListingGroupPageWithSignals> =>
      client.listListingsByProduct(options) as unknown as Promise<ListingGroupPageWithSignals>,
    getListingsSummary: (installationId: string): Promise<ListingSummaryWithSignals> =>
      client.getListingsSummary(installationId) as unknown as Promise<ListingSummaryWithSignals>,
  };
}
