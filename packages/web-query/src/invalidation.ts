import type { QueryClient } from "@tanstack/react-query";

import { queryKeyNamespaces } from "./index";

export type MutationInvalidationType =
  | "price_update"
  | "stock_correct"
  | "link_apply"
  | "listing_pause"
  | "listing_resync"
  | "listing_edit"
  | "product_enrichment";

export class UnknownMutationInvalidationTypeError extends Error {
  readonly code = "UNKNOWN_MUTATION_INVALIDATION_TYPE";

  constructor(readonly type: unknown) {
    super(`Unknown mutation invalidation type: ${String(type)}`);
    this.name = "UnknownMutationInvalidationTypeError";
  }
}

const mutationInvalidationNamespaces = {
  price_update: ["listings", "mutations"],
  stock_correct: ["listings", "inventory", "mutations"],
  link_apply: ["listings", "linkage", "catalog", "mutations"],
  listing_pause: ["listings", "mutations"],
  listing_resync: ["listings", "mutations"],
  listing_edit: ["listings", "mutations"],
  product_enrichment: ["catalog"],
} as const satisfies Record<MutationInvalidationType, readonly (keyof typeof queryKeyNamespaces)[]>;

export async function invalidateAfterMutation(
  queryClient: QueryClient,
  type: MutationInvalidationType,
): Promise<void> {
  if (!Object.prototype.hasOwnProperty.call(mutationInvalidationNamespaces, type)) {
    throw new UnknownMutationInvalidationTypeError(type);
  }

  const namespaces = mutationInvalidationNamespaces[type as MutationInvalidationType];
  await Promise.all(
    namespaces.map((namespace) =>
      queryClient.invalidateQueries({ queryKey: queryKeyNamespaces[namespace] }),
    ),
  );
}
