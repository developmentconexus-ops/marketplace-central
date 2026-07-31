import { useMutation, useQuery, useQueryClient, type UseQueryResult } from "@tanstack/react-query";
import type {
  ActiveSourceConfig,
  ActiveSourceName,
  CatalogAssortmentCounts,
  SellableAssortmentConfig,
  SetActiveSourceRequest,
  SetSellableAssortmentRequest,
} from "@marketplace-central/sdk-runtime";
import { catalogQueryKeys, queryKeyNamespaces } from "./index";

// The active source is per-tenant state in the database, not a browser
// preference. Every mirror-backed read resolves it server-side (routing.Reader
// looks it up and pins it on the request context), so a selection kept only in
// localStorage changed the labels in the UI while the API kept serving the
// other source — the toggle looked live and did nothing. This module is the
// single client-side seam over GET/PUT /config/active-source.

export interface ActiveSourceClient {
  getActiveSource: () => Promise<ActiveSourceConfig>;
  setActiveSource: (req: SetActiveSourceRequest) => Promise<ActiveSourceConfig>;
}

export const activeSourceQueryKeys = {
  config: () => ["config", "active-source"] as const,
};

/**
 * useActiveSourceQuery reads the tenant's configured source. There is no
 * client-side default: a tenant with no row fails closed server-side (400
 * unknown_erp_source), and inventing "xlsx" here would show the operator a
 * source the platform is not actually reading.
 */
export function useActiveSourceQuery(client: ActiveSourceClient): UseQueryResult<ActiveSourceConfig> {
  return useQuery({
    queryKey: activeSourceQueryKeys.config(),
    queryFn: () => client.getActiveSource(),
    staleTime: 60_000,
  });
}

/**
 * useSetActiveSourceMutation flips the tenant's source and then invalidates
 * EVERY query: the source decides what products, costs, stock and prices the
 * whole app reads, so any cached page from the previous source is now wrong
 * data under a correct-looking label.
 */
export function useSetActiveSourceMutation(client: ActiveSourceClient) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (source: ActiveSourceName) => client.setActiveSource({ active_source: source }),
    onSuccess: async (config) => {
      queryClient.setQueryData(activeSourceQueryKeys.config(), config);
      await queryClient.invalidateQueries();
    },
  });
}

export interface SellableAssortmentClient {
  getSellableAssortment: () => Promise<SellableAssortmentConfig>;
  setSellableAssortment: (req: SetSellableAssortmentRequest) => Promise<SellableAssortmentConfig>;
  getCatalogAssortmentCounts: () => Promise<CatalogAssortmentCounts>;
}

export const sellableAssortmentQueryKeys = {
  config: () => ["config", "sellable-assortment"] as const,
};

/**
 * useSellableAssortmentQuery reads the tenant's stored assortment rule. There
 * is no browser default: an absent or failed read must not make an unchecked
 * box look like the server's persisted choice.
 */
export function useSellableAssortmentQuery(
  client: SellableAssortmentClient,
): UseQueryResult<SellableAssortmentConfig> {
  return useQuery({
    queryKey: sellableAssortmentQueryKeys.config(),
    queryFn: () => client.getSellableAssortment(),
    staleTime: 60_000,
  });
}

/**
 * The counts endpoint resolves the active source from request context rather
 * than accepting it as a parameter. Keeping that source in the key prevents
 * a source flip from reusing the previous source's count under a new label.
 */
export function useCatalogAssortmentCountsQuery(
  client: SellableAssortmentClient,
  erpSource: ActiveSourceName | undefined,
): UseQueryResult<CatalogAssortmentCounts> {
  return useQuery({
    queryKey: catalogQueryKeys.counts(erpSource),
    queryFn: () => client.getCatalogAssortmentCounts(),
    staleTime: 300_000,
    enabled: Boolean(erpSource),
  });
}

/**
 * useSetSellableAssortmentMutation sends the complete stored rule because the
 * server rejects partial bodies. The response is authoritative; only the
 * catalog namespace is invalidated because this rule does not affect other
 * tenant screens.
 */
export function useSetSellableAssortmentMutation(client: SellableAssortmentClient) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (config: SetSellableAssortmentRequest) => client.setSellableAssortment(config),
    onSuccess: async (config) => {
      queryClient.setQueryData(sellableAssortmentQueryKeys.config(), config);
      await queryClient.invalidateQueries({ queryKey: queryKeyNamespaces.catalog });
    },
  });
}
