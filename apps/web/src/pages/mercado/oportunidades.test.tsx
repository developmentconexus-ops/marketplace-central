import { render, screen, within } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import type {
  CatalogProductFact,
  MarketPriceIntelAggregate,
  MarketPriceIntelVerdict,
} from "@marketplace-central/sdk-runtime";
import { buildOppRows, chunk } from "./oportunidades";
import { OportunidadesTable } from "./OportunidadesTable";

// EXEMPLO-IO golden: codprod 90008 "PAPELEIRA DECA FLEX" — custo ERP R$118,99;
// market (ml_catalog_offers, collected 2026-07-19) mediana R$229,20 / menor R$179,90;
// evidence observed → the Oportunidades row must render the real cost/median/offer-count,
// and honest "—" for VENDAS LÍDER 30D, MARGEM EST. and VEREDICTO (all M-07/snapshot-owned).
function fact90008(): CatalogProductFact {
  return {
    internal_product_id: 90008,
    reference: null,
    description: "PAPELEIRA DECA FLEX",
    ean: null,
    manufacturer_reference: null,
    brand_name: null,
    ncm: null,
    quality_flags: [],
    active: true,
    sellable_stock: { quantity: null, quality: [] },
    current_price: { amount: null, currency: "BRL", quality: [] },
    cost: { amount: "118.99", currency: "BRL", quality: [] },
  };
}

function agg90008(status: MarketPriceIntelAggregate["status"] = "OK"): MarketPriceIntelAggregate {
  return {
    product_id: "90008",
    median: { amount: "229.20", currency: "BRL" },
    min_valid: { amount: "179.90", currency: "BRL" },
    n_offers: 6,
    n_sellers: 5, // >=5 so the aggregate is OK (contract: n_sellers<5 ⇒ INSUFFICIENT_MARKET)
    source: "ml_catalog_offers",
    fetched_at: "2026-07-19T06:00:00Z",
    computed_at: "2026-07-19T06:05:00Z",
    status,
  };
}

function verdict90008(): MarketPriceIntelVerdict {
  return {
    match_status: "ACCEPT",
    price_evidence_status: "OK",
    verdict_label: null, // M-07-owned — always null pre-M-07
    blocking_state: null,
    inputs_used: {},
    market_range: {
      min_valid: { amount: "179.90", currency: "BRL" },
      median: { amount: "229.20", currency: "BRL" },
      currency: "BRL",
      n_offers: 6,
      n_sellers: 4,
    },
  };
}

describe("buildOppRows", () => {
  it("joins the 90008 fact/aggregate/verdict into an honest opportunity row", () => {
    const rows = buildOppRows([fact90008()], [agg90008()], [verdict90008()]);
    expect(rows).toHaveLength(1);
    expect(rows[0]).toMatchObject({
      sku: "90008",
      name: "PAPELEIRA DECA FLEX",
      costAmount: "118.99",
      median: { amount: "229.20", currency: "BRL" },
      nSellers: 5,
      evidenceState: "OK",
    });
    // No client-computed margin field — operational margin is M-07-owned (ADR-17).
    expect(rows[0]).not.toHaveProperty("marginPct");
    // No raw offer count either — CONCORRENTES + sort use distinct sellers, never n_offers.
    expect(rows[0]).not.toHaveProperty("nOffers");
  });

  it("excludes products the market cannot price (status != OK or no median)", () => {
    expect(
      buildOppRows([fact90008()], [agg90008("INSUFFICIENT_MARKET")], [verdict90008()]),
    ).toHaveLength(0);
    // no aggregate at all for the codprod
    expect(buildOppRows([fact90008()], [], [])).toHaveLength(0);
  });

  it("orders by observed competition (n_sellers desc — the CONCORRENTES value), not by margin, offers, or SKU", () => {
    // The higher-n_sellers product (99999) is deliberately given an alphabetically-LATER SKU and
    // FEWER offers than 90008 — so ["99999","90008"] can ONLY come from an n_sellers-desc sort,
    // never from a SKU sort (that yields ["90008","99999"]) nor an n_offers sort.
    const otherFact: CatalogProductFact = { ...fact90008(), internal_product_id: 99999 };
    const moreSellers: MarketPriceIntelAggregate = {
      ...agg90008(),
      product_id: "99999",
      n_sellers: 9,
      n_offers: 6,
    };
    const fewerSellers: MarketPriceIntelAggregate = { ...agg90008(), n_sellers: 5, n_offers: 20 };
    const rows = buildOppRows(
      [fact90008(), otherFact],
      [fewerSellers, moreSellers],
      [verdict90008(), verdict90008()],
    );
    expect(rows.map((r) => r.sku)).toEqual(["99999", "90008"]); // 9 sellers first, despite later SKU + fewer offers
  });
});

describe("buildOppRows — fetchedAt", () => {
  it("carrega a idade do agregado para a linha", () => {
    const rows = buildOppRows([fact90008()], [agg90008()], [verdict90008()]);
    expect(rows[0].fetchedAt).toBe("2026-07-19T06:00:00Z");
  });

  it("renderiza a idade na linha de Oportunidades", () => {
    const rows = buildOppRows([fact90008()], [agg90008()], [verdict90008()]);
    render(<OportunidadesTable rows={rows} />);
    // O agregado é de 2026-07-19; qualquer "agora" seria idade fabricada.
    expect(screen.getAllByLabelText("Atualização dos dados")[0]).toHaveTextContent(/^há \d+ d$/);
  });
});

describe("chunk", () => {
  it("splits into ordered batches no larger than the /market MaxReadIDs cap, preserving order", () => {
    // 5 ids at size 2 → [[a,b],[c,d],[e]] — order preserved so the concatenated per-chunk
    // aggregates/verdicts still align positionally with the original codprod/fact order.
    expect(chunk(["a", "b", "c", "d", "e"], 2)).toEqual([["a", "b"], ["c", "d"], ["e"]]);
    // Exactly at the cap → a single batch (no spurious extra request).
    expect(chunk(["a", "b"], 2)).toEqual([["a", "b"]]);
    expect(chunk([], 200)).toEqual([]);
    expect(() => chunk(["a"], 0)).toThrow();
  });
});

describe("OportunidadesTable — EXEMPLO-IO 90008 render", () => {
  it("renders real cost/median/sellers and honest dashes for M-07/snapshot columns", () => {
    const rows = buildOppRows([fact90008()], [agg90008()], [verdict90008()]);
    render(<OportunidadesTable rows={rows} />);

    const row = screen.getByText("PAPELEIRA DECA FLEX").closest("div")!;
    expect(within(row).getByText("90008")).toBeInTheDocument();
    expect(within(row).getByText("R$ 118,99")).toBeInTheDocument();
    expect(within(row).getByText("R$ 229,20")).toBeInTheDocument();
    // CONCORRENTES = distinct sellers (n_sellers=5), NOT the raw offer count (n_offers=6):
    // rendering offers under a "competitors" header would overstate competition (ADR-17).
    expect(within(row).getByText("5")).toBeInTheDocument();
    expect(within(row).queryByText("6")).toBeNull();
    expect(within(row).getByText("mercado observado")).toBeInTheDocument(); // honest evidence note
    // MARGEM EST. (M-07-owned — no commission/freight/DIFAL), VENDAS LÍDER 30D (no snapshot)
    // and VEREDICTO recommendation are all unbacked → honest dash, never a fabricated % (ADR-17)
    expect(within(row).queryByText(/%$/)).toBeNull();
    expect(within(row).getAllByText("—").length).toBeGreaterThanOrEqual(3);
    // Inert action — no live ML write (D-57)
    expect(within(row).getByRole("button", { name: "Criar anúncio" })).toBeDisabled();
    // Footer must describe the ACTUAL ordering (by observed competition = n_sellers), not a
    // claim the sort doesn't honor — locks the copy so it can't silently drift from the key.
    expect(screen.getByText(/ordenado por diferença mediana − custo/)).toBeInTheDocument();
  });
});
