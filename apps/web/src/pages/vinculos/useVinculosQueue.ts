import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import type {
  MarketplaceCentralClientError,
  ProductLinkActor,
  ProductLinkCandidateItem,
} from "@marketplace-central/sdk-runtime";
import { invalidateAfterMutation, QUERY_STALE_TIME } from "@marketplace-central/web-query";
import { useClient } from "../../app/ClientContext";
import { vinculosQueryKeys } from "./vinculosQueryKeys";

/**
 * PAGE-LOCAL query-key root (Finding C). The Vínculos queue is keyed under
 * `["product-links", ...]` and does NOT live in the web-query `linkage`
 * namespace. Every approve/reject/manual mutation MUST invalidate this root
 * directly — invalidateAfterMutation("link_apply") alone will NOT refresh it.
 */
export const PRODUCT_LINKS_ROOT = ["product-links"] as const;

/**
 * How many candidates the queue loads at once. The SDK default (20) is a sample,
 * not a work list: a generation run over a whole account produces thousands of
 * rows, so 20 showed one anúncio's variations and hid every real suggestion. The
 * backend ranks by confidence, so this page carries the actionable ones — and the
 * UI says so when the list is capped instead of implying the account has exactly
 * this many pending.
 */
export const VINCULOS_QUEUE_PAGE_SIZE = 200;

/**
 * UI-supplied, unverified operator actor (no auth context in the demo shell).
 * Exported so batch flows (useVinculosBatch) mirror the same actor instead of
 * redefining it.
 */
export const WEB_OPERATOR_ACTOR: ProductLinkActor = {
  actor_type: "operator",
  actor_id: "web-operator",
  actor_name: "Operador",
};

/**
 * A 409 from approve means the listing was already resolved by someone else.
 * The contract code is ALREADY_RESOLVED — surface it (banner + refetch), never
 * swallow it.
 */
export function isAlreadyResolvedConflict(error: unknown): boolean {
  if (typeof error !== "object" || error === null) return false;
  const candidate = error as Partial<MarketplaceCentralClientError>;
  return candidate.status === 409 && candidate.error?.code === "ALREADY_RESOLVED";
}

export interface ApproveVinculoInput {
  candidate_id: string;
  reason?: string;
}

export interface RejectVinculoInput {
  installation_id: string;
  provider_code: string;
  provider_item_id: string;
  provider_variation_id?: string;
  reason?: string;
}

export interface ManualVinculoInput {
  installation_id: string;
  provider_code: string;
  provider_item_id: string;
  provider_variation_id?: string;
  internal_product_id: number;
  internal_product_name?: string;
  internal_reference_code?: string;
  reason?: string;
}

export function rejectInputForCandidate(candidate: ProductLinkCandidateItem): RejectVinculoInput {
  return {
    installation_id: candidate.installation_id,
    provider_code: candidate.provider_code,
    provider_item_id: candidate.provider_item_id,
    provider_variation_id: candidate.provider_variation_id,
  };
}

export function useVinculosQueue(installationId: string) {
  const client = useClient();
  const queryClient = useQueryClient();

  const queueQuery = useQuery({
    queryKey: vinculosQueryKeys.queue(installationId),
    queryFn: () => client.listProductLinkCandidates(installationId, VINCULOS_QUEUE_PAGE_SIZE),
    staleTime: QUERY_STALE_TIME.listings,
  });

  // MANDATORY: refresh the page-local queue/resolved/KPI caches. This covers the
  // success path (item leaves the queue) AND the 409 path (queue re-syncs to the
  // already-resolved reality).
  const invalidateQueue = () => queryClient.invalidateQueries({ queryKey: PRODUCT_LINKS_ROOT });

  // ALSO refresh the cross-domain listings/catalog/mutations caches on success.
  // This does NOT touch the page-local queue — that is why invalidateQueue above
  // is not optional.
  const invalidateCrossDomain = () => void invalidateAfterMutation(queryClient, "link_apply");

  const approve = useMutation({
    mutationFn: (input: ApproveVinculoInput) =>
      client.approveProductLinkCandidate({ ...input, actor: WEB_OPERATOR_ACTOR }),
    onSuccess: invalidateCrossDomain,
    onSettled: invalidateQueue,
  });

  const reject = useMutation({
    mutationFn: (input: RejectVinculoInput) =>
      client.rejectProductLinkListing({ ...input, actor: WEB_OPERATOR_ACTOR }),
    onSuccess: invalidateCrossDomain,
    onSettled: invalidateQueue,
  });

  const manual = useMutation({
    mutationFn: (input: ManualVinculoInput) =>
      client.manualResolveProductLink({ ...input, actor: WEB_OPERATOR_ACTOR }),
    onSuccess: invalidateCrossDomain,
    onSettled: invalidateQueue,
  });

  const conflict =
    isAlreadyResolvedConflict(approve.error) ||
    isAlreadyResolvedConflict(reject.error) ||
    isAlreadyResolvedConflict(manual.error);

  return { queueQuery, approve, reject, manual, conflict };
}
