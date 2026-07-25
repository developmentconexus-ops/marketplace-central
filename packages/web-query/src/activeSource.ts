import { useMutation, useQuery, useQueryClient, type UseQueryResult } from "@tanstack/react-query";
import type { ActiveSourceConfig, ActiveSourceName, SetActiveSourceRequest } from "@marketplace-central/sdk-runtime";

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
