import { useMutation, useQueryClient } from "@tanstack/react-query";
import type {
  ApplyProductLinkBatchResponse,
  PreviewProductLinkBatchResponse,
  ProductLinkBatchApproval,
} from "@marketplace-central/sdk-runtime";
import { invalidateAfterMutation } from "@marketplace-central/web-query";
import { useClient } from "../../app/ClientContext";
import { PRODUCT_LINKS_ROOT, WEB_OPERATOR_ACTOR } from "./useVinculosQueue";

/**
 * Two-phase batch flow: `preview` is a pure dry-run (nothing applied — the
 * `/batch-preview` endpoint only reports predicted OK/FAILED per candidate).
 * `apply` is the only mutation that writes; it is invoked separately, only
 * after the caller confirms the preview.
 *
 * Same page-local invalidation lesson as S7's useVinculosQueue: the
 * `link_apply` cross-domain invalidation does NOT refetch the page-local
 * `["product-links", ...]` queue/resolved/KPI caches, so `apply` always also
 * invalidates PRODUCT_LINKS_ROOT directly (onSettled — covers both the
 * success path and any partial-failure response, since the queue must
 * re-sync to reality either way).
 */
export function useVinculosBatch() {
  const client = useClient();
  const queryClient = useQueryClient();

  const preview = useMutation<PreviewProductLinkBatchResponse, unknown, ProductLinkBatchApproval[]>({
    mutationFn: (approvals) => client.previewProductLinkBatch({ approvals }),
  });

  const apply = useMutation<ApplyProductLinkBatchResponse, unknown, ProductLinkBatchApproval[]>({
    mutationFn: (approvals) =>
      client.applyProductLinkBatch({ approvals, actor: WEB_OPERATOR_ACTOR }),
    onSuccess: () => void invalidateAfterMutation(queryClient, "link_apply"),
    onSettled: () => queryClient.invalidateQueries({ queryKey: PRODUCT_LINKS_ROOT }),
  });

  return { preview, apply, actor: WEB_OPERATOR_ACTOR };
}
