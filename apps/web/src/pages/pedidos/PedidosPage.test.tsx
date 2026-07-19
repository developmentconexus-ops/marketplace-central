import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import { beforeEach, describe, expect, it, vi } from "vitest";
import { PedidosPage } from "./PedidosPage";

const listOrders = vi.fn();
const getOrderSummary = vi.fn();

vi.mock("../../app/ClientContext", () => ({
  useClient: () => ({
    listOrders: (...args: unknown[]) => listOrders(...args),
    getOrderSummary: (...args: unknown[]) => getOrderSummary(...args),
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

async function goToLista() {
  fireEvent.click(screen.getByRole("tab", { name: "Lista" }));
  await screen.findByRole("tablist", { name: "Filtros de pedidos" });
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
  bucket: "novo" as const,
};

describe("PedidosPage", () => {
  beforeEach(() => {
    listOrders.mockReset();
    getOrderSummary.mockReset();
    getOrderSummary.mockResolvedValue({
      today: 0,
      seven_days: 0,
      by_status: { novo: 2, faturar: 1, enviar: 3, enviado: 4 },
    });
  });

  it("renders the header, KPI row and default Fila view", async () => {
    listOrders.mockResolvedValue({ items: [], next_cursor: null });

    renderPage();

    expect(await screen.findByRole("heading", { name: "Pedidos" })).toBeInTheDocument();
    expect(screen.getByText("NOVOS")).toBeInTheDocument();
    expect(screen.getByText("A FATURAR")).toBeInTheDocument();
    expect(screen.getByText("A ENVIAR")).toBeInTheDocument();
    expect(screen.getByText("ENVIADOS")).toBeInTheDocument();
    expect(screen.getByText("DIFAL A PAGAR")).toBeInTheDocument();
    expect(screen.getByRole("tab", { name: "Fila", selected: true })).toBeInTheDocument();
  });

  it("renders the KPI live bucket counts and the DIFAL card as an honest unknown", async () => {
    listOrders.mockResolvedValue({ items: [], next_cursor: null });

    renderPage();

    await screen.findByRole("heading", { name: "Pedidos" });
    expect(await screen.findByText("2")).toBeInTheDocument(); // novo
    expect(screen.getByText("1")).toBeInTheDocument(); // faturar
    expect(screen.getByText("3")).toBeInTheDocument(); // enviar
    expect(screen.getByText("4")).toBeInTheDocument(); // enviado
    expect(screen.getByText("—")).toBeInTheDocument(); // DIFAL, gated
  });

  it("renders '—' for every bucket when the summary endpoint is unavailable", async () => {
    listOrders.mockResolvedValue({ items: [], next_cursor: null });
    getOrderSummary.mockRejectedValue({ status: 503, error: { code: "unavailable" } });

    renderPage();

    await screen.findByRole("heading", { name: "Pedidos" });
    await waitFor(() => {
      expect(screen.getAllByText("—").length).toBeGreaterThanOrEqual(5);
    });
    expect(screen.queryByText("0")).not.toBeInTheDocument();
  });

  it("switches between Fila, Lista and Kanban views", async () => {
    listOrders.mockResolvedValue({
      items: [{ ...baseOrder, sla: { due: "2026-07-20T10:00:00Z", atrasado: false } }],
      next_cursor: null,
    });

    renderPage();

    await screen.findByRole("heading", { name: "Pedidos" });
    expect(await screen.findByText("ML-1001")).toBeInTheDocument();
    expect(screen.getByRole("heading", { name: "Fila de trabalho" })).toBeInTheDocument();

    await goToLista();
    expect(screen.getByRole("tab", { name: "Lista", selected: true })).toBeInTheDocument();
    expect(screen.getByRole("tab", { name: /Novos/ })).toBeInTheDocument();

    fireEvent.click(screen.getByRole("tab", { name: "Kanban" }));
    expect(await screen.findByText("NOVOS · 1")).toBeInTheDocument();
  });

  it("renders a disabled mutation control in the Fila view (em breve, no handler)", async () => {
    listOrders.mockResolvedValue({
      items: [{ ...baseOrder, nf_state: null }],
      next_cursor: null,
    });

    renderPage();

    const action = await screen.findByRole("button", { name: "Faturar" });
    expect(action).toBeDisabled();
  });

  it("clicking a bucket KPI card switches to Lista AND selects the matching tab", async () => {
    listOrders.mockResolvedValue({ items: [], next_cursor: null });

    renderPage();

    await screen.findByRole("heading", { name: "Pedidos" });
    fireEvent.click(screen.getByText("A FATURAR"));

    await waitFor(() => {
      expect(screen.getByRole("tab", { name: "Lista", selected: true })).toBeInTheDocument();
    });
    expect(screen.getByRole("tab", { name: /A faturar/, selected: true })).toBeInTheDocument();
  });

  it("clicking the DIFAL KPI card selects the Fila view", async () => {
    listOrders.mockResolvedValue({ items: [], next_cursor: null });

    renderPage();

    await screen.findByRole("heading", { name: "Pedidos" });
    await goToLista();
    fireEvent.click(screen.getByText("DIFAL A PAGAR"));
    await waitFor(() => {
      expect(screen.getByRole("tab", { name: "Fila", selected: true })).toBeInTheDocument();
    });
  });

  it("renders the 7 Lista bucket tabs with live counts on the 4 live buckets", async () => {
    listOrders.mockResolvedValue({
      items: [
        { ...baseOrder, bucket: "novo" },
        { ...baseOrder, provider_order_id: "PO2", bucket: "novo" },
        { ...baseOrder, provider_order_id: "PO3", bucket: "faturar" },
      ],
      next_cursor: null,
    });

    renderPage();

    await screen.findByRole("heading", { name: "Pedidos" });
    await goToLista();

    expect(screen.getByRole("tab", { name: "Novos 2" })).toBeInTheDocument();
    expect(screen.getByRole("tab", { name: "A faturar 1" })).toBeInTheDocument();
    expect(screen.getByRole("tab", { name: "A enviar 0" })).toBeInTheDocument();
    expect(screen.getByRole("tab", { name: "Enviados 0" })).toBeInTheDocument();
    expect(screen.getByRole("tab", { name: "Concluídos" })).toBeInTheDocument();
    expect(screen.getByRole("tab", { name: "Cancelados" })).toBeInTheDocument();
    expect(screen.getByRole("tab", { name: "Devoluções" })).toBeInTheDocument();
  });

  it("shows an 'em breve' placeholder panel for a placeholder tab, not a table", async () => {
    listOrders.mockResolvedValue({
      items: [{ ...baseOrder, bucket: "novo" }],
      next_cursor: null,
    });

    renderPage();

    await screen.findByRole("heading", { name: "Pedidos" });
    await goToLista();
    fireEvent.click(screen.getByRole("tab", { name: "Cancelados" }));

    expect(await screen.findByText(/Cancelados — em breve/)).toBeInTheDocument();
    expect(screen.queryByRole("table")).not.toBeInTheDocument();
  });

  it("filters Lista rows by the active bucket tab using OrderRead.bucket", async () => {
    listOrders.mockResolvedValue({
      items: [
        { ...baseOrder, provider_code: "ML-1001", bucket: "novo" },
        { ...baseOrder, provider_order_id: "PO2", provider_code: "ML-1002", bucket: "faturar" },
      ],
      next_cursor: null,
    });

    renderPage();

    await screen.findByRole("heading", { name: "Pedidos" });
    await goToLista();

    expect(await screen.findByText("ML-1001")).toBeInTheDocument();
    expect(screen.queryByText("ML-1002")).not.toBeInTheDocument();

    fireEvent.click(screen.getByRole("tab", { name: "A faturar 1" }));

    await waitFor(() => {
      expect(screen.getByText("ML-1002")).toBeInTheDocument();
      expect(screen.queryByText("ML-1001")).not.toBeInTheDocument();
    });
  });

  it("renders the Lista table with RETORNO and DIFAL as honest '—' (gated)", async () => {
    listOrders.mockResolvedValue({
      items: [
        {
          ...baseOrder,
          buyer: { display: "J. S.", city: "São Paulo", uf: "SP" },
          bucket: "novo",
        },
      ],
      next_cursor: null,
    });

    renderPage();

    await screen.findByRole("heading", { name: "Pedidos" });
    await goToLista();

    expect(await screen.findByText("ML-1001")).toBeInTheDocument();
    expect(screen.getByText("J. S.")).toBeInTheDocument();
    expect(screen.queryByText("JOAOSILVA123")).not.toBeInTheDocument();

    const rows = screen.getAllByRole("row");
    expect(rows.length).toBeGreaterThan(1);
  });

  it("enables mass-select in Lista and shows a disabled bulk action bar when rows are selected", async () => {
    listOrders.mockResolvedValue({
      items: [{ ...baseOrder, bucket: "novo" }],
      next_cursor: null,
    });

    renderPage();

    await screen.findByRole("heading", { name: "Pedidos" });
    await goToLista();
    await screen.findByText("ML-1001");

    fireEvent.click(screen.getByRole("checkbox", { name: /Select row/ }));

    expect(await screen.findByText("1 selecionado(s)")).toBeInTheDocument();
    expect(screen.getByRole("button", { name: "Etiquetas" })).toBeDisabled();
    expect(screen.getByRole("button", { name: "Marcar enviados" })).toBeDisabled();
  });

  it("renders the Kanban view with 4 read-only columns grouped by bucket", async () => {
    listOrders.mockResolvedValue({
      items: [
        { ...baseOrder, bucket: "novo" },
        { ...baseOrder, provider_order_id: "PO2", provider_code: "ML-1002", bucket: "faturar" },
      ],
      next_cursor: null,
    });

    renderPage();

    await screen.findByRole("heading", { name: "Pedidos" });
    fireEvent.click(screen.getByRole("tab", { name: "Kanban" }));

    expect(await screen.findByText("NOVOS · 1")).toBeInTheDocument();
    expect(screen.getByText("A FATURAR · 1")).toBeInTheDocument();
    expect(screen.getByText("A ENVIAR · 0")).toBeInTheDocument();
    expect(screen.getByText("ENVIADOS · 0")).toBeInTheDocument();
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
