import { useInfiniteQuery, useQuery } from "@tanstack/react-query";
import type { CatalogProductFactPage } from "@marketplace-central/sdk-runtime";
import {
  catalogQueryKeys,
  QUERY_STALE_TIME,
  type RefreshableClient,
} from "@marketplace-central/web-query";

export interface CatalogQueriesClient extends RefreshableClient {
  listCatalogProductFacts: (options?: { cursor?: string; limit?: number }) => Promise<CatalogProductFactPage>;
  searchCatalogProductFacts: (options: { q: string; limit?: number }) => Promise<CatalogProductFactPage>;
}

export function useCatalogFactsQuery(client: CatalogQueriesClient, enabled: boolean) {
  return useInfiniteQuery({
    queryKey: catalogQueryKeys.facts({ limit: 50 }),
    queryFn: ({ pageParam }) => client.listCatalogProductFacts({ cursor: pageParam, limit: 50 }),
    initialPageParam: undefined as string | undefined,
    getNextPageParam: (lastPage) => lastPage.next_cursor ?? undefined,
    staleTime: QUERY_STALE_TIME.catalog,
    enabled,
  });
}

export function useCatalogSearchQuery(client: CatalogQueriesClient, query: string) {
  return useQuery({
    queryKey: catalogQueryKeys.search(query),
    queryFn: () => client.searchCatalogProductFacts({ q: query, limit: 50 }),
    staleTime: QUERY_STALE_TIME.catalog,
    enabled: Boolean(query),
  });
}
