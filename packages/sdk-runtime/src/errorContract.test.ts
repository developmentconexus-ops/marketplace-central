import { describe, expect, it, vi } from "vitest";
import {
  createMarketplaceCentralClient,
  MarketplaceCentralClientError,
  isApiError,
  hasCode,
} from "./index";

// Covers what errorContract.golden.test.ts (the acceptance test for this
// slice, not to be edited here) does not: the fallback path for a body that
// is NOT a well-formed envelope, and the runtime/type-level shape of the
// narrowing helpers themselves.
describe("SDK error contract — fallback and narrowing (VC-5)", () => {
  it("turns a malformed body (a bare HTML string) into internal_error, preserving the original text at details.raw", async () => {
    const fetchImpl = vi.fn().mockResolvedValue({
      ok: false,
      status: 502,
      json: async () => "<html>502 Bad Gateway</html>",
    });
    const client = createMarketplaceCentralClient({ baseUrl: "http://localhost:8080", fetchImpl });

    try {
      await client.getCatalogAssortmentCounts();
      throw new Error("expected getCatalogAssortmentCounts to reject");
    } catch (e) {
      expect(e).toBeInstanceOf(MarketplaceCentralClientError);
      const err = e as MarketplaceCentralClientError;
      expect(err.code).toBe("internal_error");
      expect(err.details.raw).toBe("<html>502 Bad Gateway</html>");
    }
  });

  it("does NOT reinterpret a legacy flat body as its own code — it lands as internal_error with the flat body at details.raw", async () => {
    const flat = { error: "invalid_q" };
    const fetchImpl = vi.fn().mockResolvedValue({
      ok: false,
      status: 400,
      json: async () => flat,
    });
    const client = createMarketplaceCentralClient({ baseUrl: "http://localhost:8080", fetchImpl });

    try {
      await client.getCatalogAssortmentCounts();
      throw new Error("expected getCatalogAssortmentCounts to reject");
    } catch (e) {
      expect(e).toBeInstanceOf(MarketplaceCentralClientError);
      const err = e as MarketplaceCentralClientError;
      // Must NOT be resurrected as code "invalid_q" — a flat body is not a
      // well-formed envelope (error must be an object, not a string).
      expect(err.code).toBe("internal_error");
      expect(err.details.raw).toEqual(flat);
    }
  });

  it("defaults details to an object ({}), never undefined, when the envelope omits it", async () => {
    const fetchImpl = vi.fn().mockResolvedValue({
      ok: false,
      status: 400,
      json: async () => ({ error: { code: "invalid_q", message: "q inválido" } }),
    });
    const client = createMarketplaceCentralClient({ baseUrl: "http://localhost:8080", fetchImpl });

    try {
      await client.getCatalogAssortmentCounts();
      throw new Error("expected getCatalogAssortmentCounts to reject");
    } catch (e) {
      const err = e as MarketplaceCentralClientError;
      expect(err.details).toEqual({});
      expect(err.details).not.toBeUndefined();
    }
  });

  it("isApiError is a real instanceof check, not a duck-type on the presence of .code", () => {
    const lookalike = { status: 400, code: "x" };
    expect(isApiError(lookalike)).toBe(false);
  });

  it("rejects a bogus code at the TYPE level — discharged by tsc, not by this vitest run", () => {
    const someError: unknown = new MarketplaceCentralClientError(
      400,
      "invalid_erp_source",
      "x",
      {},
    );
    // esbuild (vitest's transform) erases types at runtime, so this line can
    // only prove itself via tsc: hasCode<C extends ApiErrorCode> rejects any
    // C not in the union, so a typo'd code ("invalid_erp_soruce", note the
    // transposed letters) fails to compile. If hasCode's generic constraint
    // is ever loosened, or the typo is ever added to ApiErrorCode for real,
    // the directive below stops being satisfied and `tsc --noEmit` fails on
    // this line. This assertion is NOT proven by running vitest — only by
    // the tsc lane in this slice's validation.
    // @ts-expect-error "invalid_erp_soruce" is not a member of ApiErrorCode (typo of invalid_erp_source)
    const result = hasCode(someError, "invalid_erp_soruce");
    expect(typeof result).toBe("boolean");
  });

  it("hasCode returns false when the code matches but the spec-guaranteed field is missing, while isApiError and .code stay unaffected", () => {
    // invalid_erp_source guarantees details.allowed_range (per
    // ApiErrorDetailsByCode) — an envelope that carries the code but an
    // empty details object is a server bug, not a value hasCode should
    // narrow into. A consumer must fall through to its generic branch
    // instead of reading `undefined` while the type said `string`.
    const err = new MarketplaceCentralClientError(
      400,
      "invalid_erp_source",
      "erp_source inválido",
      {},
    );

    expect(isApiError(err)).toBe(true);
    expect(err.code).toBe("invalid_erp_source");
    expect(hasCode(err, "invalid_erp_source")).toBe(false);
  });
});
