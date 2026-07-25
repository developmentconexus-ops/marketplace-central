import { useInfiniteQuery, useQuery } from "@tanstack/react-query";
import type { ActiveSourceName, CatalogProductFactPage } from "@marketplace-central/sdk-runtime";
import {
  catalogQueryKeys,
  QUERY_STALE_TIME,
  type RefreshableClient,
} from "@marketplace-central/web-query";

export interface CatalogQueriesClient extends RefreshableClient {
  listCatalogProductFacts: (options?: { cursor?: string; limit?: number }) => Promise<CatalogProductFactPage>;
  searchCatalogProductFacts: (options: { q: string; limit?: number }) => Promise<CatalogProductFactPage>;
}

// The source is NOT sent as a request parameter. The reader resolves the
// tenant's configured active source from the database and pins it on the
// request context (routing.Reader), which overwrites any client-supplied
// erp_source — so the screen could label one source while reading another. It
// stays in the query KEY so a flip cannot serve the previous source's cached
// page under the new label.

export function useCatalogFactsQuery(
  client: CatalogQueriesClient,
  enabled: boolean,
  erpSource: ActiveSourceName | undefined,
) {
  return useInfiniteQuery({
    queryKey: catalogQueryKeys.facts({ limit: 50, erp_source: erpSource ?? null }),
    queryFn: ({ pageParam }) => client.listCatalogProductFacts({ cursor: pageParam, limit: 50 }),
    initialPageParam: undefined as string | undefined,
    getNextPageParam: (lastPage) => lastPage.next_cursor ?? undefined,
    staleTime: QUERY_STALE_TIME.catalog,
    enabled,
  });
}

export function useCatalogSearchQuery(
  client: CatalogQueriesClient,
  query: string,
  erpSource: ActiveSourceName | undefined,
) {
  return useQuery({
    queryKey: catalogQueryKeys.search(query, { erp_source: erpSource ?? null }),
    queryFn: () => client.searchCatalogProductFacts({ q: query, limit: 50 }),
    staleTime: QUERY_STALE_TIME.catalog,
    enabled: Boolean(query),
  });
}
