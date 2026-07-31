import { describe, expect, it, vi } from "vitest";
import { createMarketplaceCentralClient, MarketplaceCentralClientError, isApiError, hasCode } from "./index";

// This test pins the FUTURE unified SDK error surface (day-1 golden test,
// RED before the chip's SDK refactor lands): a typed client-side error class,
// an `isApiError` guard, and a `hasCode` narrowing helper — none of which
// exist yet. Today `getCatalogAssortmentCounts` also takes no arguments, so
// calling it with `{ erp_source }` is itself part of the pinned target
// contract. Either the assertions below fail, or (more likely at this stage)
// the file fails to COMPILE — both are an acceptable red for this slice; see
// the chip report for which one happened.
describe("unified SDK error contract (golden, RED before refactor)", () => {
  it("throws a typed MarketplaceCentralClientError with a narrowed details.allowed_range", async () => {
    const envelope = {
      error: {
        code: "invalid_erp_source",
        message: "erp_source inválido: use xlsx ou catalogo_cliente",
        details: { allowed_range: "xlsx|catalogo_cliente" },
      },
    };
    const fetchImpl = vi.fn().mockResolvedValue({
      ok: false,
      status: 400,
      json: async () => envelope,
    });
    const client = createMarketplaceCentralClient({ baseUrl: "http://localhost:8080", fetchImpl });

    // "banana" is deliberately not a valid ErpImportSourceInput — it is what
    // triggers the mocked 400. The @ts-expect-error directives below are
    // themselves part of the pinned contract: they prove erp_source is typed
    // as ErpImportSourceInput | undefined, not string. If that option ever
    // loosens to `string`, these directives go unused and the tsc lane fails
    // — which is the behaviour we want.
    // @ts-expect-error erp_source: "banana" is deliberately invalid to trigger the mocked 400
    await expect(client.getCatalogAssortmentCounts({ erp_source: "banana" })).rejects.toBeDefined();

    try {
      // @ts-expect-error erp_source: "banana" is deliberately invalid to trigger the mocked 400
      await client.getCatalogAssortmentCounts({ erp_source: "banana" });
      throw new Error("expected getCatalogAssortmentCounts to reject");
    } catch (e) {
      expect(e instanceof MarketplaceCentralClientError).toBe(true);
      expect(isApiError(e)).toBe(true);
      expect(hasCode(e, "invalid_erp_source")).toBe(true);
      if (hasCode(e, "invalid_erp_source")) {
        const allowed: string = e.details.allowed_range;
        expect(allowed).toBe("xlsx|catalogo_cliente");
        expect(e.status).toBe(400);
        expect(e.message).toBe("erp_source inválido: use xlsx ou catalogo_cliente");
      }
    }
  });
});
