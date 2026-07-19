import { MemoryRouter } from "react-router-dom";
import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import type { DashboardOverview } from "@marketplace-central/sdk-runtime";
import { FilaDeAtencao } from "./FilaDeAtencao";

const baseSummary: DashboardOverview = {
  sync_errors: 0,
  pending_links: 0,
  below_margin: 0,
  missing_gtin: 0,
  orders_today: 0,
  orders_7d: 0,
  anuncios_ativos: 0,
  last_sync_at: {},
  degraded: [],
  last_import: null,
  as_of: "2026-07-19T12:00:00Z",
};

function renderFila(summary: DashboardOverview) {
  return render(
    <MemoryRouter>
      <FilaDeAtencao summary={summary} />
    </MemoryRouter>,
  );
}

describe("FilaDeAtencao", () => {
  it("deep-links each attention item to the owner screen with the exception filter", () => {
    renderFila({ ...baseSummary, sync_errors: 3, below_margin: 2, pending_links: 5 });

    expect(screen.getByRole("link", { name: /erro de sincroniza/i })).toHaveAttribute(
      "href",
      "/anuncios?filter.exception=sync_error",
    );
    expect(screen.getByRole("link", { name: /abaixo da margem/i })).toHaveAttribute(
      "href",
      "/anuncios?filter.exception=below_margin",
    );
    expect(screen.getByRole("link", { name: /sem v[ií]nculo/i })).toHaveAttribute("href", "/vinculos");
  });

  it("omits an item when its count is null or zero", () => {
    renderFila({ ...baseSummary, sync_errors: 0, below_margin: null, pending_links: 5 });

    expect(screen.queryByRole("link", { name: /erro de sincroniza/i })).not.toBeInTheDocument();
    expect(screen.queryByRole("link", { name: /abaixo da margem/i })).not.toBeInTheDocument();
    expect(screen.getByRole("link", { name: /sem v[ií]nculo/i })).toBeInTheDocument();
  });

  it("never renders a pedidos-atrasados item (no /pedidos filter contract)", () => {
    renderFila({ ...baseSummary, sync_errors: 1, below_margin: 1, pending_links: 1 });
    expect(screen.queryByText(/atrasad/i)).not.toBeInTheDocument();
  });

  it("shows an empty hint when there is nothing to attend to", () => {
    renderFila(baseSummary);
    expect(screen.getByText(/nenhum item na fila/i)).toBeInTheDocument();
  });
});
