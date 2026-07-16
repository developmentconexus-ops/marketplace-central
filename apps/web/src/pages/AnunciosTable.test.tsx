import type { ListingReadModel } from "@marketplace-central/sdk-runtime";
import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import { AnunciosTable } from "./AnunciosTable";

const tableSelectionProps = {
  selectedIds: new Set<string>(),
  onToggle: () => undefined,
  onTogglePage: () => undefined,
};

const listing: ListingReadModel = {
  listing_id: "listing_1",
  installation_id: "installation_1",
  provider: "mercado_livre",
  provider_listing_id: "MLB123456789",
  title: "Camiseta azul",
  listing_type: { code: "gold_special", label: "Clássico" },
  status: "active",
  link: { state: "resolved", product_id: "product_1", seller_sku: "CAM-AZ" },
  price: { amount: "129.90", currency: "BRL" },
  published_quantity: 7,
  sync_state: "synced",
  sync_error: null,
  quality_score: 0.9,
  pending_issue: { kind: "stale", message_pt: "Atualização pendente" },
  sales_30d: 12,
  cost: { amount: "70.00", currency: "BRL" },
  below_margin_worst_case: true,
  icms_worst_case_by_uf: null,
  fetched_at: "2026-07-16T12:00:00Z",
};

describe("AnunciosTable", () => {
  it("renders the listing facts and pending issue", () => {
    render(<AnunciosTable items={[listing]} asOf="2026-07-16T12:00:00Z" {...tableSelectionProps} />);

    expect(screen.getByText("Camiseta azul")).toBeInTheDocument();
    expect(screen.getByText("MLB123456789")).toBeInTheDocument();
    expect(screen.getByText("Clássico")).toBeInTheDocument();
    expect(screen.getByText("vinculado")).toBeInTheDocument();
    expect(screen.getByText("R$ 129.90")).toBeInTheDocument();
    expect(screen.getByText("7")).toBeInTheDocument();
    expect(screen.getByText("12")).toBeInTheDocument();
    expect(screen.getByText("abaixo da margem")).toBeInTheDocument();
    expect(screen.getByText("sincronizado")).toBeInTheDocument();
    expect(screen.getByLabelText("Atualização pendente")).toHaveAttribute("title", "Atualização pendente");
  });

  it("renders honest unknowns for nullable values", () => {
    render(
      <AnunciosTable
        items={[{ ...listing, price: null, published_quantity: null, sales_30d: null, below_margin_worst_case: null }]}
        {...tableSelectionProps}
      />,
    );

    expect(screen.getAllByText("—")).toHaveLength(4);
    expect(screen.queryByText("0")).not.toBeInTheDocument();
    expect(screen.queryByText("R$ 0,00")).not.toBeInTheDocument();
  });

  it("explains an unsimulated margin when ERP cost is unknown", () => {
    render(<AnunciosTable items={[{ ...listing, cost: null, below_margin_worst_case: null }]} {...tableSelectionProps} />);

    expect(screen.getByTitle("sem custo no ERP → não simulado")).toHaveTextContent("—");
  });

  it("uses the fixed sync and conflict labels", () => {
    render(<AnunciosTable items={[{ ...listing, sync_state: "error", link: { ...listing.link, state: "conflict" }}]} {...tableSelectionProps} />);

    expect(screen.getByText("com erro")).toBeInTheDocument();
    expect(screen.getByText("divergente")).toBeInTheDocument();
  });
});
