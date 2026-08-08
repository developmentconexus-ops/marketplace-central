import type { createMarketplaceCentralClient } from "@marketplace-central/sdk-runtime";
import { listingsQueryKeys, QUERY_STALE_TIME } from "@marketplace-central/web-query";
import { toListingListOptions, type AnunciosQueryState } from "./anunciosQueryState";

export type AnunciosClient = Pick<
  ReturnType<typeof createMarketplaceCentralClient>,
  "listListings" | "getListingsSummary" | "refreshListings" | "listIntegrationOperationRuns"
>;

// Each factory asks for exactly the methods it calls, not for the whole
// AnunciosClient. A wider parameter forces every caller -- tests included -- to
// supply capabilities the function must never reach for, which makes the type
// unable to say what the function actually does.
export function anunciosPageQuery(
  client: Pick<AnunciosClient, "listListings">,
  installationId: string,
  state: AnunciosQueryState,
) {
  const options = toListingListOptions(state, installationId);
  return {
    queryKey: listingsQueryKeys.page(installationId, options),
    queryFn: () => client.listListings(options),
    staleTime: QUERY_STALE_TIME.listings,
  };
}

export function anunciosSummaryQuery(client: Pick<AnunciosClient, "getListingsSummary">, installationId: string) {
  return {
    queryKey: listingsQueryKeys.summary(installationId),
    queryFn: () => client.getListingsSummary(installationId),
    staleTime: QUERY_STALE_TIME.listings,
  };
}
