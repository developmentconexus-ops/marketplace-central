import { describe, expect, it } from "vitest";
import type { ListingGroupPage, ListingPage, ListingSummary } from "./index";
import { withListingSignals, type ListingsSignalClient } from "./listings";

// Fixture mirrors the additive OpenAPI delta: an OK-shape item with a full
// market_signal, and a SEM_VINCULO item where market_signal is honestly null
// (ADR-17 — unlinked never gets a fabricated signal).
const okItem = {
  listing_id: "l1",
  installation_id: "inst1",
  provider: "mercado_livre",
  provider_listing_id: "ML1",
  title: "Produto A",
  listing_type: null,
  status: "active",
  link: { state: "resolved", product_id: "p1", seller_sku: "sku1" },
  price: { amount: "100.00", currency: "BRL" },
  published_quantity: 10,
  sync_state: "synced",
  sync_error: null,
  quality_score: null,
  pending_issue: null,
  sales_30d: null,
  cost: null,
  below_margin_worst_case: null,
  icms_worst_case_by_uf: null,
  fetched_at: "2026-07-18T12:00:00Z",
  market_signal: {
    status: "OK",
    position: { rank: 1, total: 5 },
    price_to_win: { amount: "95.00", currency: "BRL" },
    delta_pct: "-5.0",
    match_status: "ACCEPT",
    n_offers: 5,
    n_sellers: 4,
    evidence: { source: "ml_price_to_win", fetched_at: "2026-07-18T11:00:00Z", freshness: "fresh" },
  },
  signal_status: "OK",
};

const semVinculoItem = {
  ...okItem,
  listing_id: "l2",
  link: { state: "unresolved", product_id: null, seller_sku: null },
  market_signal: null,
  signal_status: "SEM_VINCULO",
};

const stubClient: ListingsSignalClient = {
  listListings: async () =>
    ({
      items: [okItem, semVinculoItem],
      next_cursor: null,
      page_size: 2,
      as_of: "2026-07-18T12:00:00Z",
    }) as unknown as ListingPage,
  listListingsByProduct: async () =>
    ({
      groups: [],
      next_cursor: null,
      page_size: 0,
      as_of: "2026-07-18T12:00:00Z",
    }) as unknown as ListingGroupPage,
  getListingsSummary: async () =>
    ({
      total: 2,
      active: 2,
      paused: 0,
      exceptions: {
        sync_error: 0,
        stale: 0,
        unlinked: 1,
        below_margin_worst_case: 0,
        margin_unknown: 0,
        sem_vinculo: 1,
        abaixo_custo: 0,
        sem_evidencia: 0,
      },
      as_of: "2026-07-18T12:00:00Z",
    }) as unknown as ListingSummary,
};

describe("withListingSignals", () => {
  it("re-types listListings items so market_signal/signal_status are accessible", async () => {
    const wrapped = withListingSignals(stubClient);
    const page = await wrapped.listListings({ installation_id: "inst1" });

    expect(page.items[0].signal_status).toBe("OK");
    expect(page.items[0].market_signal?.status).toBe("OK");
    expect(page.items[0].market_signal?.n_offers).toBe(5);

    expect(page.items[1].signal_status).toBe("SEM_VINCULO");
    expect(page.items[1].market_signal).toBeNull();
  });

  it("re-types getListingsSummary.exceptions with the extended counters", async () => {
    const wrapped = withListingSignals(stubClient);
    const summary = await wrapped.getListingsSummary("inst1");

    expect(summary.exceptions.sem_vinculo).toBe(1);
    expect(summary.exceptions.abaixo_custo).toBe(0);
    expect(summary.exceptions.sem_evidencia).toBe(0);
  });

  it("does not perform a duplicate fetch — wraps the same stub instance", async () => {
    let calls = 0;
    const countingClient: ListingsSignalClient = {
      ...stubClient,
      listListings: async (...args) => {
        calls += 1;
        return stubClient.listListings(...args);
      },
    };
    const wrapped = withListingSignals(countingClient);
    await wrapped.listListings({ installation_id: "inst1" });
    expect(calls).toBe(1);
  });
});
