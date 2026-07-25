// F01-S1 rework (SDK-LISTINGS-M05): the market_signal / signal_status surface
// now lives inline in ./index.ts (no standalone listings.ts wrapper — see
// grant amendment superseding 5318b29's approach). These are behavior tests
// over the inline types: typed field access plus the ADR-17 negative cases
// that must stay distinct (SEM_VINCULO vs NO_PRICE_EVIDENCE vs STALE, and
// below-cost vs unknown-cost).

import { describe, expect, it } from "vitest";
import type {
  ListingMarketSignal,
  ListingReadModel,
  ListingSummaryExceptions,
} from "./index";

function baseListing(overrides: Partial<ListingReadModel>): ListingReadModel {
  return {
    listing_id: "l1",
    installation_id: "inst1",
    provider: "mercado_livre",
    provider_listing_id: "MLB1",
    variation_id: null,
    title: "Produto Teste",
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
    cost: { amount: "80.00", currency: "BRL" },
    below_margin_worst_case: false,
    icms_worst_case_by_uf: null,
    fetched_at: "2026-07-18T00:00:00Z",
    ...overrides,
  };
}

describe("inline market_signal / signal_status typed access", () => {
  it("exposes a fully-populated market_signal with typed evidence and money reuse", () => {
    const signal: ListingMarketSignal = {
      status: "OK",
      position: { rank: 2, total: 8 },
      price_to_win: { amount: "95.50", currency: "BRL" },
      delta_pct: "-4.5",
      match_status: "ACCEPT",
      n_offers: 8,
      n_sellers: 6,
      evidence: {
        source: "mercado_livre_search",
        fetched_at: "2026-07-18T00:00:00Z",
        freshness: "fresh",
      },
    };
    const listing = baseListing({ market_signal: signal, signal_status: "OK" });

    expect(listing.signal_status).toBe("OK");
    expect(listing.market_signal?.price_to_win?.currency).toBe("BRL");
    expect(listing.market_signal?.evidence?.source).toBe("mercado_livre_search");
  });

  it("SEM_VINCULO: market_signal is null while signal_status carries the honest state (never fabricated)", () => {
    const listing = baseListing({
      link: { state: "unresolved", product_id: null, seller_sku: null },
      market_signal: null,
      signal_status: "SEM_VINCULO",
    });

    expect(listing.signal_status).toBe("SEM_VINCULO");
    expect(listing.market_signal).toBeNull();
  });

  it("NO_PRICE_EVIDENCE and STALE stay distinct signal_status values for a linked listing", () => {
    const noEvidence = baseListing({ market_signal: null, signal_status: "NO_PRICE_EVIDENCE" });
    const stale = baseListing({
      market_signal: {
        status: "STALE",
        position: null,
        price_to_win: null,
        delta_pct: null,
        match_status: null,
        n_offers: null,
        n_sellers: null,
        evidence: {
          source: "mercado_livre_search",
          fetched_at: "2026-06-01T00:00:00Z",
          freshness: "expired",
        },
      },
      signal_status: "STALE",
    });

    expect(noEvidence.signal_status).not.toBe(stale.signal_status);
    expect(noEvidence.market_signal).toBeNull();
    expect(stale.market_signal?.status).toBe("STALE");
  });

  it("summary counters keep abaixo_custo (below-cost, known cost) distinct from sem_evidencia (unknown/no evidence)", () => {
    const exceptions: ListingSummaryExceptions = {
      sync_error: 0,
      stale: 1,
      unlinked: 0,
      below_margin_worst_case: 2,
      margin_unknown: 3,
      sem_vinculo: 4,
      abaixo_custo: 5,
      sem_evidencia: 6,
    };

    // Distinct honest counters: below-cost (known cost, price under it) never
    // collapses into unknown-evidence/unknown-cost buckets.
    expect(exceptions.abaixo_custo).toBe(5);
    expect(exceptions.sem_evidencia).toBe(6);
    expect(exceptions.abaixo_custo).not.toBe(exceptions.sem_evidencia);
  });

  it("optional signal fields are omittable so pre-existing W1 consumers remain valid", () => {
    const listing = baseListing({});

    expect(listing.market_signal).toBeUndefined();
    expect(listing.signal_status).toBeUndefined();
  });
});
