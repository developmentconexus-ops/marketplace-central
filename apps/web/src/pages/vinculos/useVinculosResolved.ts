import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import type { ProductLinkAuditEntry, ProductLinkWorkflowItem } from "@marketplace-central/sdk-runtime";
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
export function resolutionAuditEntry(item: ProductLinkWorkflowItem): ProductLinkAuditEntry | undefined {
  const resolving = (item.audit ?? [])
    .filter((entry) => entry.next_state === "resolved")
    .sort((a, b) => Date.parse(b.created_at) - Date.parse(a.created_at));
  return resolving[0];
}

export function resolutionAuditId(item: ProductLinkWorkflowItem): string | undefined {
  return resolutionAuditEntry(item)?.audit_id;
}

/**
 * Whether the vínculo was resolved BY THE SYSTEM rather than by an operator.
 *
 * The predicate is the resolving audit entry's `actor.actor_type === "system"`.
 * The auto-linker writes that entry with `ActorType: "system", ActorID:
 * "auto_linker"` (resolution_service.go:280) and `audit[].actor.actor_type` is
 * already on the wire, so the badge costs no contract change.
 *
 * It is deliberately NOT the `rule_matched = exact_ean_unique AND actor =
 * system` predicate the M-06 brief asks for: `rule_matched` is not on the wire
 * at all (zero hits in `contracts/` and `packages/sdk-runtime/src/` — it exists
 * only in the DB and in a per-link repo read no route exposes), and that exact
 * pair is forbidden by the CHECK at 0082_product_link_decisions.sql:54
 * (`actor <> 'system' OR rule_matched = 'concordant_codprod_ean'`) — the brief
 * predates D-121, which narrowed auto-approval to concordant CODPROD+EAN.
 *
 * A pre-M-05 link with no resolving audit entry returns false: no badge, which
 * says "not automatic" — true for every manual vínculo — instead of fabricating
 * one (ADR-17).
 */
export function isAutoResolved(item: ProductLinkWorkflowItem): boolean {
  return resolutionAuditEntry(item)?.actor?.actor_type === "system";
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
