import { describe, expect, it, vi } from "vitest";
import { listingsQueryKeys, QUERY_STALE_TIME } from "@marketplace-central/web-query";
import {
  applyAnunciosQueryState,
  clearFilters,
  parseAnunciosQueryState,
  toListingListOptions,
} from "./anunciosQueryState";
import { anunciosPageQuery, anunciosSummaryQuery } from "./anunciosQueries";

describe("anuncios query state", () => {
  it("round-trips a canonical deep link and normalizes its options", () => {
    const searchParams = new URLSearchParams(
      "installation=inst_1&tab=pendencia&filter.exception=sync_error",
    );
    const state = parseAnunciosQueryState(searchParams);

    expect(applyAnunciosQueryState(searchParams, state).toString()).toBe(searchParams.toString());
    expect(toListingListOptions(state, "inst_1")).toEqual({
      installation_id: "inst_1",
      has_exception: true,
      exception: "sync_error",
    });
  });

  it.each([
    ["", {}, "todos"],
    ["tab=ativos", { status: "active" }, "ativos"],
    ["tab=pausados", { status: "paused" }, "pausados"],
    ["tab=pendencia", { has_exception: true }, "pendencia"],
  ])("maps %s to the SDK tab filter", (query, expected, tab) => {
    const state = parseAnunciosQueryState(new URLSearchParams(query));

    expect(state.tab).toBe(tab);
    expect(toListingListOptions(state, "inst_1")).toEqual({
      installation_id: "inst_1",
      ...expected,
    });
  });

  it("drops unknown filter keys and invalid enum values", () => {
    const state = parseAnunciosQueryState(
      new URLSearchParams("filter.bogus=x&filter.exception=nope&filter.sync_state=synced"),
    );

    expect(toListingListOptions(state, "inst_1")).toEqual({
      installation_id: "inst_1",
      sync_state: "synced",
    });
  });

  it.each(["sync_error", "stale", "unlinked", "below_margin", "sem_vinculo", "abaixo_custo", "sem_evidencia"])(
    "round-trips the %s exception filter",
    (exception) => {
      const searchParams = new URLSearchParams(`installation=inst_1&filter.exception=${exception}`);
      const state = parseAnunciosQueryState(searchParams);

      expect(state.filters.exception).toBe(exception);
      expect(applyAnunciosQueryState(searchParams, state).toString()).toBe(searchParams.toString());
      expect(toListingListOptions(state, "inst_1")).toEqual({
        installation_id: "inst_1",
        exception,
      });
    },
  );

  it("ignores an invalid exception value without crashing", () => {
    const searchParams = new URLSearchParams("installation=inst_1&filter.exception=xyz");
    const state = parseAnunciosQueryState(searchParams);

    expect(state.filters.exception).toBeUndefined();
    expect(toListingListOptions(state, "inst_1")).toEqual({ installation_id: "inst_1" });
    expect(clearFilters(searchParams).toString()).toBe("installation=inst_1");
  });

  it("clears every filter while preserving route state", () => {
    const searchParams = new URLSearchParams(
      "q=camiseta&tab=ativos&installation=inst_1&filter.exception=sync_error&filter.bogus=x",
    );

    expect(clearFilters(searchParams).toString()).toBe(
      "q=camiseta&tab=ativos&installation=inst_1",
    );
  });

  it("preserves unrelated params when applying state", () => {
    const searchParams = new URLSearchParams("installation=inst_1&view=compact");
    const state = parseAnunciosQueryState(new URLSearchParams("tab=ativos&q=camiseta"));

    const next = applyAnunciosQueryState(searchParams, state);

    expect(next.toString()).toBe("installation=inst_1&view=compact&tab=ativos&q=camiseta");
  });

  it("defaults grouped to false when the param is absent", () => {
    const state = parseAnunciosQueryState(new URLSearchParams("installation=inst_1"));

    expect(state.grouped).toBe(false);
  });

  it("round-trips the grouped=1 param", () => {
    const searchParams = new URLSearchParams("installation=inst_1&grouped=1");
    const state = parseAnunciosQueryState(searchParams);

    expect(state.grouped).toBe(true);
    expect(applyAnunciosQueryState(searchParams, state).toString()).toBe(searchParams.toString());
  });

  it("omits the grouped param from the URL when false", () => {
    const searchParams = new URLSearchParams("installation=inst_1&grouped=1");
    const state = parseAnunciosQueryState(searchParams);

    const next = applyAnunciosQueryState(searchParams, { ...state, grouped: false });

    expect(next.has("grouped")).toBe(false);
  });
});

describe("anuncios query factories", () => {
  it("uses the normalized options in the listing query key and SDK query function", async () => {
    const client = {
      listListings: vi.fn().mockResolvedValue({ items: [] }),
      getListingsSummary: vi.fn(),
    };
    const state = parseAnunciosQueryState(
      new URLSearchParams("q=camiseta&tab=pendencia&filter.exception=sync_error"),
    );
    const options = anunciosPageQuery(client, "inst_1", state);

    expect(options.queryKey).toEqual(
      listingsQueryKeys.page("inst_1", {
        installation_id: "inst_1",
        q: "camiseta",
        has_exception: true,
        exception: "sync_error",
      }),
    );
    expect(options.staleTime).toBe(QUERY_STALE_TIME.listings);
    await options.queryFn();
    expect(client.listListings).toHaveBeenCalledWith({
      installation_id: "inst_1",
      q: "camiseta",
      has_exception: true,
      exception: "sync_error",
    });
  });

  it("uses the summary key and SDK summary method", async () => {
    const client = {
      listListings: vi.fn(),
      getListingsSummary: vi.fn().mockResolvedValue({ total: 0 }),
    };
    const options = anunciosSummaryQuery(client, "inst_1");

    expect(options.queryKey).toEqual(listingsQueryKeys.summary("inst_1"));
    expect(options.staleTime).toBe(QUERY_STALE_TIME.listings);
    await options.queryFn();
    expect(client.getListingsSummary).toHaveBeenCalledWith("inst_1");
  });
});
