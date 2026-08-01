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
 *
 * networkMode: "always" overrides the app QueryClient's default ("online",
 * see createWebQueryClient in ./index.ts). Under the default, TanStack Query
 * gates retries on onlineManager.isOnline(): if the browser's reported
 * connectivity flips to offline while the first attempt is failing (a stale
 * navigator.onLine reading, a flapping adapter, a devtools/CDP network
 * emulation — any of which can coincide with a real backend outage), the
 * query is marked fetchStatus "paused" after its first failure and never
 * retries or surfaces an error again: isLoading, isError, and data are all
 * simultaneously false/undefined, so SyncHealthCard's
 * isLoading ? … : isError ? … : data ? … : null chain falls through every
 * branch and the card renders nothing but its heading, forever — reproduced
 * live by the hub (backend container killed, DevTools showed exactly one
 * request, zero retries) and confirmed here via real instrumentation in
 * SyncHealthCard.realNetwork.test.tsx (status: "pending", fetchStatus:
 * "paused", stuck). "always" makes this card's own polling read-path ignore
 * the browser's online/offline signal and always attempt the fetch, so a
 * real connection failure runs its normal retry/backoff and lands on
 * isError like any other error — the card never goes silently blank.
 */
export function useSyncHealthQuery(client: SyncHealthClient): UseQueryResult<SyncHealth> {
  return useQuery({
    queryKey: syncHealthQueryKeys.health(),
    queryFn: () => client.getSyncHealth(),
    staleTime: QUERY_STALE_TIME.sync,
    refetchInterval: 30_000,
    networkMode: "always",
  });
}
