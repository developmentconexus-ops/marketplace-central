import { useQuery, useQueryClient, type QueryClient } from "@tanstack/react-query";
import type { MutationState } from "@marketplace-central/sdk-runtime";
import {
  invalidateAfterMutation,
  mutationsQueryKeys,
  QUERY_STALE_TIME,
  type MutationInvalidationType,
} from "@marketplace-central/web-query";
import { useEffect } from "react";
import { useClient } from "../../app/ClientContext";

const terminalStates = new Set<MutationState>(["applied", "partially_failed", "failed_preserved", "cancelled"]);
const invalidatedProtocols = new WeakMap<QueryClient, Set<string>>();

export function isMutationTerminal(state: MutationState): boolean {
  return terminalStates.has(state);
}

export function useMutationProtocol(protocolId: string) {
  const client = useClient();
  const queryClient = useQueryClient();
  const query = useQuery({
    queryKey: mutationsQueryKeys.detail(protocolId),
    queryFn: () => client.getMutation(protocolId),
    enabled: protocolId.length > 0,
    staleTime: QUERY_STALE_TIME.mutations,
    refetchInterval: (current) => {
      const protocol = current.state.data;
      return protocol && isMutationTerminal(protocol.state) ? false : 2000;
    },
  });

  useEffect(() => {
    const protocol = query.data;
    if (!protocol || !isMutationTerminal(protocol.state)) return;
    let invalidated = invalidatedProtocols.get(queryClient);
    if (!invalidated) {
      invalidated = new Set<string>();
      invalidatedProtocols.set(queryClient, invalidated);
    }
    if (invalidated.has(protocolId)) return;
    invalidated.add(protocolId);
    void invalidateAfterMutation(queryClient, protocol.type as MutationInvalidationType);
  }, [protocolId, query.data, queryClient]);

  return query;
}
