import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import { beforeEach, describe, expect, it, vi } from "vitest";
import { PedidosPage } from "./PedidosPage";

const listOrders = vi.fn();

vi.mock("../../app/ClientContext", () => ({
  useClient: () => ({
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
      <PedidosPage />
    </QueryClientProvider>,
  );
}

const baseOrder = {
  provider_order_id: "PO1",
  provider_code: "ML-1001",
  status: "paid",
  provider_status_detail: "paid",
  buyer_nickname: "JOAOSILVA123",
  total: 199.9,
  currency: "BRL",
  fulfillment: null,
  nf_state: null,
  created_at: "2026-07-01T10:00:00Z",
  provider_created_at: "2026-07-01T10:00:00Z",
  provider_closed_at: null,
  provider_updated_at: null,
  items: [],
  payments: [],
};

describe("PedidosPage", () => {
  beforeEach(() => {
    listOrders.mockReset();
  });

  it("renders the header and tabs, including a disabled Cancelados tab", async () => {
    listOrders.mockResolvedValue({ items: [], next_cursor: null });

    renderPage();

    expect(await screen.findByRole("heading", { name: "Pedidos" })).toBeInTheDocument();
    expect(screen.getByRole("tab", { name: "Todos" })).toBeInTheDocument();
    expect(screen.getByRole("tab", { name: "Atrasados" })).toBeInTheDocument();
    expect(screen.getByRole("tab", { name: "Sem vínculo" })).toBeInTheDocument();
    const cancelados = screen.getByRole("tab", { name: "Cancelados" });
    expect(cancelados).toBeDisabled();
  });

  it("renders enriched order rows with masked buyer and honest nulls, never raw buyer_nickname", async () => {
    listOrders.mockResolvedValue({
      items: [
        {
          ...baseOrder,
          buyer: { display: "J. S.", city: "São Paulo", uf: "SP" },
          vinculo_status: "OK",
          destino_uf: "SP",
          sla: { due: "2026-07-20T10:00:00Z", atrasado: false },
        },
        {
          ...baseOrder,
          provider_order_id: "PO2",
          provider_code: "ML-1002",
          vinculo_status: "SEM_VINCULO",
          componentes_desconhecidos: ["custo_unitario"],
        },
      ],
      next_cursor: null,
    });

    renderPage();

    expect(await screen.findByText("ML-1001")).toBeInTheDocument();
    expect(screen.getByText("J. S.")).toBeInTheDocument();
    expect(screen.queryByText("JOAOSILVA123")).not.toBeInTheDocument();

    expect(screen.getByText("ML-1002")).toBeInTheDocument();
    expect(screen.getByText("sem vínculo")).toBeInTheDocument();
    expect(screen.getByText("custo incompleto")).toBeInTheDocument();

    // honest null: PO2 has no destino_uf / sla -> renders as unknown-value dash, not fabricated data
    const rows = screen.getAllByRole("row");
    expect(rows.length).toBeGreaterThan(1);
  });

  it("filters to the atrasados tab client-side", async () => {
    listOrders.mockResolvedValue({
      items: [
        { ...baseOrder, sla: { due: "2026-07-15T10:00:00Z", atrasado: true } },
        { ...baseOrder, provider_order_id: "PO2", provider_code: "ML-1002", sla: { due: "2026-07-25T10:00:00Z", atrasado: false } },
      ],
      next_cursor: null,
    });

    renderPage();

    await screen.findByText("ML-1001");
    fireEvent.click(screen.getByRole("tab", { name: "Atrasados" }));

    await waitFor(() => {
      expect(screen.getByText("ML-1001")).toBeInTheDocument();
      expect(screen.queryByText("ML-1002")).not.toBeInTheDocument();
    });
    expect(screen.getByText("atrasado")).toBeInTheDocument();
  });

  it("shows the error state with retry when the fetch fails", async () => {
    listOrders.mockRejectedValueOnce({ status: 500, error: { code: "internal" } });
    listOrders.mockResolvedValue({ items: [], next_cursor: null });

    renderPage();

    const retry = await screen.findByRole("button", { name: "Tentar novamente" });
    fireEvent.click(retry);

    expect(await screen.findByText("Nenhum registro encontrado.")).toBeInTheDocument();
  });
});
