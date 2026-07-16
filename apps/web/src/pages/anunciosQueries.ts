import type { createMarketplaceCentralClient } from "@marketplace-central/sdk-runtime";
import { listingsQueryKeys, QUERY_STALE_TIME } from "@marketplace-central/web-query";
import { toListingListOptions, type AnunciosQueryState } from "./anunciosQueryState";

export type AnunciosClient = Pick<
  ReturnType<typeof createMarketplaceCentralClient>,
  "listListings" | "getListingsSummary"
>;

export function anunciosPageQuery(
  client: AnunciosClient,
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

export function anunciosSummaryQuery(client: AnunciosClient, installationId: string) {
  return {
    queryKey: listingsQueryKeys.summary(installationId),
    queryFn: () => client.getListingsSummary(installationId),
    staleTime: QUERY_STALE_TIME.listings,
  };
}
