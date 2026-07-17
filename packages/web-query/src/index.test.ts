import {
  QUERY_STALE_TIME,
  catalogQueryKeys,
  installationsQueryKeys,
  listingsQueryKeys,
  mutationsQueryKeys,
  ordersQueryKeys,
  queryKeyNamespaces,
  syncQueryKeys,
} from "./index";

describe("web-query registry", () => {
  it("keeps the existing stale times and adds the IC-05 classes", () => {
    expect(QUERY_STALE_TIME).toEqual({
      catalog: 300_000,
      stock: 45_000,
      pricecost: 120_000,
      listings: 45_000,
      mutations: 5_000,
      orders: 120_000,
      sync: 30_000,
      market: 300_000,
    });
  });

  it("preserves existing namespaces and adds the IC-05 namespaces", () => {
    expect(queryKeyNamespaces.catalog).toEqual(["catalog"]);
    expect(queryKeyNamespaces.inventory).toEqual(["inventory"]);
    expect(queryKeyNamespaces.linkage).toEqual(["linkage"]);
    expect(queryKeyNamespaces.profitability).toEqual(["profitability"]);
    expect(queryKeyNamespaces.listings).toEqual(["listings"]);
    expect(queryKeyNamespaces.mutations).toEqual(["mutations"]);
    expect(queryKeyNamespaces.orders).toEqual(["orders"]);
    expect(queryKeyNamespaces.sync).toEqual(["sync"]);
    expect(queryKeyNamespaces.market).toEqual(["market"]);
    expect(queryKeyNamespaces.installations).toEqual(["installations"]);
  });

  it("retains the existing catalog key shapes", () => {
    expect(catalogQueryKeys.facts()).toEqual(["catalog", "facts", { params: {} }]);
    expect(catalogQueryKeys.search("phone")).toEqual(["catalog", "search", "phone"]);
  });

  it("builds exact namespace-prefixed keys for IC-05 queries", () => {
    const filters = { status: "active", q: "phone" };
    const keys = [
      ["listings", listingsQueryKeys.page("inst_1", filters)],
      ["listings", listingsQueryKeys.byProduct("inst_1", filters)],
      ["listings", listingsQueryKeys.detail("listing_1")],
      ["listings", listingsQueryKeys.summary("inst_1")],
      ["mutations", mutationsQueryKeys.list("inst_1", filters)],
      ["mutations", mutationsQueryKeys.detail("MP-000001")],
      ["mutations", mutationsQueryKeys.items("MP-000001")],
      ["orders", ordersQueryKeys.list("inst_1", filters)],
      ["sync", syncQueryKeys.runs("inst_1", filters)],
      ["installations", installationsQueryKeys.list()],
    ] as const;

    for (const [namespace, key] of keys) {
      expect(key[0]).toBe(namespace);
    }

    expect(listingsQueryKeys.page("inst_1", filters)).toEqual([
      "listings",
      "page",
      { installation_id: "inst_1", filters },
    ]);
    expect(listingsQueryKeys.byProduct("inst_1", filters)).toEqual([
      "listings",
      "by-product",
      { installation_id: "inst_1", filters },
    ]);
    expect(listingsQueryKeys.detail("listing_1")).toEqual(["listings", "detail", "listing_1"]);
    expect(listingsQueryKeys.summary("inst_1")).toEqual([
      "listings",
      "summary",
      { installation_id: "inst_1" },
    ]);
    expect(mutationsQueryKeys.list("inst_1", filters)).toEqual([
      "mutations",
      "list",
      { installation_id: "inst_1", filters },
    ]);
    expect(mutationsQueryKeys.detail("MP-000001")).toEqual([
      "mutations",
      "detail",
      "MP-000001",
    ]);
    expect(mutationsQueryKeys.items("MP-000001")).toEqual([
      "mutations",
      "items",
      "MP-000001",
    ]);
    expect(ordersQueryKeys.list("inst_1", filters)).toEqual([
      "orders",
      "list",
      { installation_id: "inst_1", filters },
    ]);
    expect(syncQueryKeys.runs("inst_1", filters)).toEqual([
      "sync",
      "runs",
      { installation_id: "inst_1", filters },
    ]);
    expect(installationsQueryKeys.list()).toEqual(["installations", "list"]);
  });
});
