import { createWebQueryClient, queryKeyNamespaces } from "./index";
import {
  invalidateAfterMutation,
  UnknownMutationInvalidationTypeError,
} from "./invalidation";

describe("invalidateAfterMutation", () => {
  const cases = [
    ["price_update", ["listings", "mutations"]],
    ["stock_correct", ["listings", "inventory", "mutations"]],
    ["link_apply", ["listings", "linkage", "catalog", "mutations"]],
    ["listing_pause", ["listings", "mutations"]],
    ["listing_resync", ["listings", "mutations"]],
    ["listing_edit", ["listings", "mutations"]],
    ["product_enrichment", ["catalog"]],
  ] as const;

  it.each(cases)("invalidates the exact namespaces for %s", async (type, expected) => {
    const queryClient = createWebQueryClient();
    const invalidateQueries = vi.spyOn(queryClient, "invalidateQueries");

    await invalidateAfterMutation(queryClient, type);

    const actual = new Set(
      invalidateQueries.mock.calls.map(([filters]) => filters.queryKey[0]),
    );
    expect(actual).toEqual(new Set(expected));
    expect(invalidateQueries).toHaveBeenCalledTimes(expected.length);
    for (const namespace of expected) {
      expect(invalidateQueries).toHaveBeenCalledWith({
        queryKey: queryKeyNamespaces[namespace],
      });
    }
  });

  it("rejects unknown types before invalidating anything", async () => {
    const queryClient = createWebQueryClient();
    const invalidateQueries = vi.spyOn(queryClient, "invalidateQueries");

    await expect(
      invalidateAfterMutation(queryClient, "bogus" as any),
    ).rejects.toBeInstanceOf(UnknownMutationInvalidationTypeError);
    expect(invalidateQueries).not.toHaveBeenCalled();
  });
});
