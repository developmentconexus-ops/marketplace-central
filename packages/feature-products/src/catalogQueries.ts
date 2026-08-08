import { useInfiniteQuery } from "@tanstack/react-query";
import type { ActiveSourceName, CatalogProductFactPage } from "@marketplace-central/sdk-runtime";
import {
  catalogQueryKeys,
  QUERY_STALE_TIME,
  type RefreshableClient,
  type SellableAssortmentClient,
} from "@marketplace-central/web-query";

// SellableAssortmentClient is folded in (not just getCatalogAssortmentCounts)
// because useCatalogAssortmentCountsQuery's parameter type is the full
// interface: a client missing setSellableAssortment/getSellableAssortment
// would fail to satisfy it structurally, not just be incomplete.
export interface CatalogQueriesClient extends RefreshableClient, SellableAssortmentClient {
  listCatalogProductFacts: (options?: {
    cursor?: string;
    limit?: number;
    include_all?: boolean;
  }) => Promise<CatalogProductFactPage>;
  searchCatalogProductFacts: (options: {
    q: string;
    cursor?: string;
    limit?: number;
    include_all?: boolean;
  }) => Promise<CatalogProductFactPage>;
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
  includeAll: boolean,
) {
  return useInfiniteQuery({
    queryKey: catalogQueryKeys.facts({
      limit: 50,
      erp_source: erpSource ?? null,
      include_all: includeAll,
    }),
    queryFn: ({ pageParam }) =>
      client.listCatalogProductFacts({
        cursor: pageParam,
        limit: 50,
        // Absent, never `false`: the handler treats both as identical, and an
        // explicit param would land ahead of `q` in the search URL below.
        include_all: includeAll ? true : undefined,
      }),
    initialPageParam: undefined as string | undefined,
    getNextPageParam: (lastPage) => lastPage.next_cursor ?? undefined,
    staleTime: QUERY_STALE_TIME.catalog,
    enabled,
  });
}

// Search pages like the list read. It used to be a single useQuery, which could
// only ever hold the first 50 matches: a query with more matches than that
// showed a full table with no way to reach the rest, and the response gave the
// screen nothing to say so with.
export function useCatalogSearchQuery(
  client: CatalogQueriesClient,
  query: string,
  erpSource: ActiveSourceName | undefined,
  includeAll: boolean,
) {
  return useInfiniteQuery({
    queryKey: catalogQueryKeys.search(query, {
      erp_source: erpSource ?? null,
      include_all: includeAll,
    }),
    queryFn: ({ pageParam }) =>
      client.searchCatalogProductFacts({
        q: query,
        cursor: pageParam,
        limit: 50,
        include_all: includeAll ? true : undefined,
      }),
    initialPageParam: undefined as string | undefined,
    getNextPageParam: (lastPage) => lastPage.next_cursor ?? undefined,
    staleTime: QUERY_STALE_TIME.catalog,
    enabled: Boolean(query),
  });
}
