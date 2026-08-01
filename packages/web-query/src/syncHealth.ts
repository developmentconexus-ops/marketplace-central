import { useQuery, type UseQueryResult } from "@tanstack/react-query";
import type { SyncHealth } from "@marketplace-central/sdk-runtime";
import { QUERY_STALE_TIME, queryKeyNamespaces } from "./index";

// Single client-side seam over GET /sync/health (F-01). Read-only: the card
// this feeds (IntegracoesPage's "Saúde do sync") never mutates, so there is
// no companion mutation here.

export interface SyncHealthClient {
  getSyncHealth: () => Promise<SyncHealth>;
}

export const syncHealthQueryKeys = {
  health: () => [...queryKeyNamespaces.sync, "health"] as const,
};

/**
 * useSyncHealthQuery polls /sync/health every 30s (feature brief: light
 * polling, no websocket) so the card reflects per-entity sync state and the
 * webhook block without a manual refresh action.
 */
export function useSyncHealthQuery(client: SyncHealthClient): UseQueryResult<SyncHealth> {
  return useQuery({
    queryKey: syncHealthQueryKeys.health(),
    queryFn: () => client.getSyncHealth(),
    staleTime: QUERY_STALE_TIME.sync,
    refetchInterval: 30_000,
  });
}
