import { fireEvent, render, screen, within } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";
import type { PricingDifalListResponse, PricingDifalRate } from "@marketplace-central/sdk-runtime";
import { BR_UFS } from "./ParamsDrawer";
import { DIFAL_DISCLAIMER, DifalDrawer } from "./DifalDrawer";

function rate(uf: string, overrides: Partial<PricingDifalRate> = {}): PricingDifalRate {
  const interestadual = ["MG", "PR", "RJ", "RS", "SC", "SP"].includes(uf) ? "12" : "7";
  return {
    uf,
    interna_pct: "18",
    interestadual_pct: interestadual,
    efetivo_pct: "6",
    origem_versao: "padrao-2026",
    override: null,
    ...overrides,
  };
}

const data: PricingDifalListResponse = {
  items: BR_UFS.map((uf) => rate(uf)),
  disclaimer: DIFAL_DISCLAIMER,
};

describe("DifalDrawer", () => {
  it("does not render when closed", () => {
    render(<DifalDrawer open={false} data={data} onOverride={vi.fn()} onClose={vi.fn()} />);
    expect(screen.queryByTestId("difal-drawer")).toBeNull();
  });

  it("lists all 27 UFs with the seed disclaimer on the surface", () => {
    render(<DifalDrawer open data={data} onOverride={vi.fn()} onClose={vi.fn()} />);
    expect(screen.getByTestId("difal-disclaimer")).toHaveTextContent("não é orientação fiscal");
    // One row per UF (27), plus the header row.
    const rows = screen.getAllByRole("row");
    expect(rows).toHaveLength(BR_UFS.length + 1);
  });

  it("shows the active override value in place of the seed interna", () => {
    const withOverride: PricingDifalListResponse = {
      disclaimer: DIFAL_DISCLAIMER,
      items: [rate("SP", { override: { interna_pct: "19.05", updated_at: "2026-07-18T00:00:00Z", actor: "op" } })],
    };
    render(<DifalDrawer open data={withOverride} onOverride={vi.fn()} onClose={vi.fn()} />);
    const row = screen.getByRole("row", { name: /SP/ });
    expect(within(row).getByText("19.05")).toBeInTheDocument();
  });

  it("submits an override for a single UF", () => {
    const onOverride = vi.fn();
    render(<DifalDrawer open data={data} onOverride={onOverride} onClose={vi.fn()} />);
    fireEvent.change(screen.getByLabelText("Override interna SP"), { target: { value: "20.5" } });
    fireEvent.click(within(screen.getByRole("row", { name: /^SP/ })).getByRole("button", { name: "OK" }));
    expect(onOverride).toHaveBeenCalledWith("SP", "20.5");
  });

  it("normalizes a pt-BR comma override to dot-decimal before onOverride (F-P7-2)", () => {
    const onOverride = vi.fn();
    render(<DifalDrawer open data={data} onOverride={onOverride} onClose={vi.fn()} />);
    fireEvent.change(screen.getByLabelText("Override interna SP"), { target: { value: "20,5" } });
    fireEvent.click(within(screen.getByRole("row", { name: /^SP/ })).getByRole("button", { name: "OK" }));
    expect(onOverride).toHaveBeenCalledWith("SP", "20.5");
  });

  it("renders the loading and error states", () => {
    const { rerender } = render(<DifalDrawer open data={undefined} isLoading onOverride={vi.fn()} onClose={vi.fn()} />);
    expect(screen.getByRole("status")).toBeInTheDocument();
    rerender(<DifalDrawer open data={undefined} isError onOverride={vi.fn()} onClose={vi.fn()} />);
    expect(screen.getByText(/Falha ao carregar a tabela DIFAL/)).toBeInTheDocument();
  });
});
