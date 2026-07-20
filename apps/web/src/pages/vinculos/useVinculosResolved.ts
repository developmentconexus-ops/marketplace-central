import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import type { ProductLinkWorkflowItem } from "@marketplace-central/sdk-runtime";
import { invalidateAfterMutation, QUERY_STALE_TIME } from "@marketplace-central/web-query";
import { useClient } from "../../app/ClientContext";
import { PRODUCT_LINKS_ROOT, WEB_OPERATOR_ACTOR } from "./useVinculosQueue";
import { vinculosQueryKeys } from "./vinculosQueryKeys";

function isResolved(item: ProductLinkWorkflowItem): boolean {
  return item.current_link?.state === "resolved";
}

/**
 * The audit entry to UNDO for a resolved workflow: the most recent entry that
 * moved the link INTO the resolved state (approve_candidate / manual_resolve).
 * Undo is a genuine reversal of that specific resolution — never a fabricated
 * id. Returns undefined when no such entry exists (desfazer then unavailable).
 */
export function resolutionAuditId(item: ProductLinkWorkflowItem): string | undefined {
  const resolving = item.audit
    .filter((entry) => entry.next_state === "resolved")
    .sort((a, b) => Date.parse(b.created_at) - Date.parse(a.created_at));
  return resolving[0]?.audit_id;
}

export function useVinculosResolved(installationId: string) {
  const client = useClient();
  const queryClient = useQueryClient();

  const resolvedQuery = useQuery({
    queryKey: vinculosQueryKeys.resolved(installationId),
    queryFn: () => client.listProductLinkWorkflows(installationId),
    staleTime: QUERY_STALE_TIME.listings,
  });

  // Undo re-opens the queue AND drops the row from Resolvidos — refresh the
  // page-local root (queue + resolved + KPIs) plus the cross-domain caches.
  const invalidateQueue = () => queryClient.invalidateQueries({ queryKey: PRODUCT_LINKS_ROOT });
  const invalidateCrossDomain = () => void invalidateAfterMutation(queryClient, "link_apply");

  const undo = useMutation({
    mutationFn: (auditId: string) =>
      client.undoProductLinkResolution(auditId, { actor: WEB_OPERATOR_ACTOR }),
    onSuccess: invalidateCrossDomain,
    onSettled: invalidateQueue,
  });

  const items = (resolvedQuery.data?.items ?? []).filter(isResolved);

  return { resolvedQuery, items, undo };
}
