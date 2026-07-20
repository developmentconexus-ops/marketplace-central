import { useInfiniteQuery, useQuery } from "@tanstack/react-query";
import type { CatalogProductFactPage, ErpImportSourceInput } from "@marketplace-central/sdk-runtime";
import {
  type ActiveErpSource,
  catalogQueryKeys,
  DEFAULT_ERP_SOURCE,
  QUERY_STALE_TIME,
  type RefreshableClient,
} from "@marketplace-central/web-query";

export interface CatalogQueriesClient extends RefreshableClient {
  listCatalogProductFacts: (options?: { cursor?: string; limit?: number; erp_source?: ErpImportSourceInput }) => Promise<CatalogProductFactPage>;
  searchCatalogProductFacts: (options: { q: string; limit?: number; erp_source?: ErpImportSourceInput }) => Promise<CatalogProductFactPage>;
}

// The default source (xlsx) is omitted from the request so a default-view URL
// stays byte-stable with clients that predate the toggle; only an explicit
// non-default selection travels as erp_source.
function sourceParam(source: ActiveErpSource): ErpImportSourceInput | undefined {
  return source === DEFAULT_ERP_SOURCE ? undefined : source;
}

export function useCatalogFactsQuery(client: CatalogQueriesClient, enabled: boolean, erpSource: ActiveErpSource) {
  return useInfiniteQuery({
    queryKey: catalogQueryKeys.facts({ limit: 50, erp_source: erpSource }),
    queryFn: ({ pageParam }) => client.listCatalogProductFacts({ cursor: pageParam, limit: 50, erp_source: sourceParam(erpSource) }),
    initialPageParam: undefined as string | undefined,
    getNextPageParam: (lastPage) => lastPage.next_cursor ?? undefined,
    staleTime: QUERY_STALE_TIME.catalog,
    enabled,
  });
}

export function useCatalogSearchQuery(client: CatalogQueriesClient, query: string, erpSource: ActiveErpSource) {
  return useQuery({
    queryKey: catalogQueryKeys.search(query, { erp_source: erpSource }),
    queryFn: () => client.searchCatalogProductFacts({ q: query, limit: 50, erp_source: sourceParam(erpSource) }),
    staleTime: QUERY_STALE_TIME.catalog,
    enabled: Boolean(query),
  });
}
