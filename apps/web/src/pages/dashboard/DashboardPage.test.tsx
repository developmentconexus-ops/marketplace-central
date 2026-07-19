import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { render, screen } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import { beforeEach, describe, expect, it, vi } from "vitest";
import { DashboardPage } from "./DashboardPage";

const getDashboardSummary = vi.fn();
const listOrders = vi.fn();

vi.mock("../../app/ClientContext", () => ({
  useClient: () => ({
    getDashboardSummary: (...args: unknown[]) => getDashboardSummary(...args),
    listOrders: (...args: unknown[]) => listOrders(...args),
  }),
}));

vi.mock("../../app/InstallationContext", () => ({
  useInstallation: () => ({ installationId: "inst_1" }),
}));

function renderPage() {
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  return render(
    <QueryClientProvider client={queryClient}>
      <MemoryRouter initialEntries={["/"]}>
        <DashboardPage />
      </MemoryRouter>
    </QueryClientProvider>,
  );
}

const baseSummary = {
  sync_errors: 0,
  pending_links: 0,
  below_margin: 0,
  missing_gtin: 0,
  orders_today: 0,
  orders_7d: 0,
  anuncios_ativos: 0,
  last_sync_at: {},
  degraded: [] as string[],
  last_import: null,
  as_of: "2026-07-19T12:00:00Z",
};

describe("DashboardPage", () => {
  beforeEach(() => {
    getDashboardSummary.mockReset();
    listOrders.mockReset();
    listOrders.mockResolvedValue({ items: [], next_cursor: null });
  });

  it("renders '0' for a real zero and '—' for a null value from a healthy source, never '0' for the null one", async () => {
    getDashboardSummary.mockResolvedValue({
      ...baseSummary,
      orders_today: 0,
      orders_7d: null,
      degraded: [],
    });

    renderPage();

    const hojeLabel = await screen.findByText("Vendas hoje");
    expect(hojeLabel.nextElementSibling).toHaveTextContent("0");

    const semanaLabel = await screen.findByText("Vendas 7 dias");
    expect(semanaLabel.nextElementSibling).toHaveTextContent("—");
    expect(semanaLabel.nextElementSibling?.textContent?.trim()).not.toBe("0");
  });

  it("shows an isolated ErrorState on the sync-errors card when 'listings' is degraded, while other cards render normally", async () => {
    getDashboardSummary.mockResolvedValue({
      ...baseSummary,
      sync_errors: null,
      below_margin: null,
      anuncios_ativos: null,
      orders_today: 5,
      degraded: ["listings"],
    });

    renderPage();

    await screen.findByText("Vendas hoje");
    const alerts = await screen.findAllByRole("alert");
    expect(alerts.length).toBeGreaterThanOrEqual(1);
  });

  it("shows an ErrorState for último import when null AND erp_import is degraded", async () => {
    getDashboardSummary.mockResolvedValue({
      ...baseSummary,
      last_import: null,
      degraded: ["erp_import"],
    });

    renderPage();

    await screen.findByText("Último import ERP");
    expect(await screen.findAllByRole("alert")).not.toHaveLength(0);
  });

  it("shows an EmptyState CTA for último import when null and NOT degraded (no import yet)", async () => {
    getDashboardSummary.mockResolvedValue({
      ...baseSummary,
      last_import: null,
      degraded: [],
    });

    renderPage();

    await screen.findByText("Nenhuma importação ainda");
  });

  it("shows LoadingState while the query is pending", () => {
    getDashboardSummary.mockReturnValue(new Promise(() => {}));
    renderPage();
    expect(screen.getByRole("status")).toBeInTheDocument();
  });

  it("shows a global ErrorState with retry when the whole summary query fails", async () => {
    getDashboardSummary.mockRejectedValue(new Error("boom"));
    renderPage();
    expect(await screen.findByRole("alert")).toBeInTheDocument();
    expect(screen.getByRole("button", { name: "Tentar novamente" })).toBeInTheDocument();
  });

  it("renders the Visão geral heading", async () => {
    getDashboardSummary.mockResolvedValue(baseSummary);
    renderPage();
    expect(await screen.findByRole("heading", { name: "Visão geral" })).toBeInTheDocument();
  });
});
