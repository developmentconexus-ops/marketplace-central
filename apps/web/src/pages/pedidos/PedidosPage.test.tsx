import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import { MemoryRouter, useLocation } from "react-router-dom";
import { beforeEach, describe, expect, it, vi } from "vitest";
import { PedidosPage } from "./PedidosPage";

const listOrders = vi.fn();
const getOrderSummary = vi.fn();
const getOrder = vi.fn();

vi.mock("../../app/ClientContext", () => ({
  useClient: () => ({
    listOrders: (...args: unknown[]) => listOrders(...args),
    getOrderSummary: (...args: unknown[]) => getOrderSummary(...args),
    getOrder: (...args: unknown[]) => getOrder(...args),
  }),
}));

vi.mock("../../app/InstallationContext", () => ({
  useInstallation: () => ({ installationId: "inst_1" }),
}));

function LocationProbe() {
  return <output data-testid="location-search">{useLocation().search}</output>;
}

function renderPage(search = "") {
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  return render(
    <QueryClientProvider client={queryClient}>
      <MemoryRouter initialEntries={[`/pedidos${search}`]}>
        <PedidosPage />
        <LocationProbe />
      </MemoryRouter>
    </QueryClientProvider>,
  );
}

async function goToLista() {
  fireEvent.click(screen.getByRole("tab", { name: "Lista" }));
  await screen.findByRole("tablist", { name: "Filtros de pedidos" });
}

// F01-C1 honest-empty decomposição/DIFAL — required objects, every member null (F02-S6 real-ready
// path: same shape lights up with real numbers once the hub wires the decomposer, no UI change).
const nullDecomposicao = {
  comissao: null,
  taxa_fixa: null,
  frete: null,
  imposto: null,
  difal: null,
  tarifa_full: null,
  custo: null,
  margem_valor: null,
  margem_pct: null,
  componentes_desconhecidos: ["comissao", "taxa_fixa", "frete", "imposto", "difal", "tarifa_full", "custo"],
};

const nullDifal = {
  amount: null,
  uf_route: null,
  due_date: null,
  paid: null,
};

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
  retorno_liquido: null,
  margem_pct: null,
  decomposicao: nullDecomposicao,
  difal: nullDifal,
};

describe("PedidosPage", () => {
  beforeEach(() => {
    listOrders.mockReset();
    getOrderSummary.mockReset();
    getOrder.mockReset();
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
    // RETORNO/DIFAL render FROM order.retorno_liquido/difal.amount — both null on this fixture,
    // so both columns render honest UnknownValue "—", never a hardcoded string (ADR-17/F02-S6).
    expect(screen.getAllByText("—").length).toBeGreaterThanOrEqual(2);
  });

  it("renders the Lista RETORNO/DIFAL columns with real formatted values when present (F02-S6 real-ready path)", async () => {
    listOrders.mockResolvedValue({
      items: [
        {
          ...baseOrder,
          buyer: { display: "J. S.", city: "São Paulo", uf: "SP" },
          bucket: "novo",
          retorno_liquido: 123.45,
          margem_pct: 0.182,
          difal: { ...nullDifal, amount: 8.76 },
        },
      ],
      next_cursor: null,
    });

    renderPage();

    await screen.findByRole("heading", { name: "Pedidos" });
    await goToLista();

    expect(await screen.findByText("ML-1001")).toBeInTheDocument();
    expect(screen.getByText("R$ 123,45")).toBeInTheDocument();
    expect(screen.getByText("18.2%")).toBeInTheDocument();
    expect(screen.getByText("R$ 8,76")).toBeInTheDocument();
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

  it("Fila shows 'sem ação' (not Faturar) for a cancelado order, agreeing with Lista/Kanban (F02-S5 #1)", async () => {
    listOrders.mockResolvedValue({
      items: [{ ...baseOrder, bucket: "cancelado" as const, nf_state: null }],
      next_cursor: null,
    });

    renderPage();

    await screen.findByText("ML-1001");
    expect(screen.getByText("sem ação")).toBeInTheDocument();
    expect(screen.queryByRole("button", { name: "Faturar" })).not.toBeInTheDocument();
  });

  it("follows next_cursor to load the full order dataset across pages (F02-S5 #2)", async () => {
    listOrders.mockImplementation(({ cursor }: { cursor?: string }) => {
      if (!cursor) {
        return Promise.resolve({
          items: [{ ...baseOrder, bucket: "novo" as const }],
          next_cursor: "page2",
        });
      }
      if (cursor === "page2") {
        return Promise.resolve({
          items: [
            {
              ...baseOrder,
              provider_order_id: "PO2",
              provider_code: "ML-1002",
              bucket: "novo" as const,
            },
          ],
          next_cursor: null,
        });
      }
      return Promise.resolve({ items: [], next_cursor: null });
    });

    renderPage();

    await screen.findByRole("heading", { name: "Pedidos" });
    await goToLista();

    expect(await screen.findByText("ML-1001")).toBeInTheDocument();
    expect(await screen.findByText("ML-1002")).toBeInTheDocument();
    expect(listOrders.mock.calls.length).toBeGreaterThanOrEqual(2);
  });

  it("shows the error state with retry when the fetch fails", async () => {
    listOrders.mockRejectedValueOnce({ status: 500, error: { code: "internal" } });
    listOrders.mockResolvedValue({ items: [], next_cursor: null });

    renderPage();

    const retry = await screen.findByRole("button", { name: "Tentar novamente" });
    fireEvent.click(retry);

    expect(await screen.findByText("Nenhum registro encontrado.")).toBeInTheDocument();
  });

  const detailOrder = {
    ...baseOrder,
    buyer: { display: "J. S.", city: "São Paulo", uf: "SP" },
    bucket: "novo" as const,
    items: [
      {
        provider_item_id: "item_1",
        seller_sku: "PAR-0451",
        title: "Parafuso M8x40 cx100",
        quantity: 2,
        unit_price: 49.9,
        link_quality: "resolved" as const,
        custo_unitario: 21.1,
      },
    ],
  };

  it("clicking a row opens the drawer, calls getOrder and sets the ?order= param", async () => {
    listOrders.mockResolvedValue({ items: [detailOrder], next_cursor: null });
    getOrder.mockResolvedValue(detailOrder);

    renderPage();

    await screen.findByRole("heading", { name: "Pedidos" });
    await goToLista();
    fireEvent.click(await screen.findByText("ML-1001"));

    expect(getOrder).toHaveBeenCalledWith("inst_1", "PO1");
    await waitFor(() => {
      expect(screen.getByTestId("location-search")).toHaveTextContent("order=PO1");
    });
    expect(await screen.findByText("Parafuso M8x40 cx100")).toBeInTheDocument();
  });

  it("opening with a ?order= URL opens the drawer with real detail data", async () => {
    listOrders.mockResolvedValue({ items: [detailOrder], next_cursor: null });
    getOrder.mockResolvedValue(detailOrder);

    renderPage("?order=PO1");

    expect(getOrder).toHaveBeenCalledWith("inst_1", "PO1");
    expect(await screen.findByText("J. S.")).toBeInTheDocument();
  });

  it("renders decomposição/DIFAL and buyer-doc/NF-number/tracking-code as honest '—', and a disabled footer action", async () => {
    listOrders.mockResolvedValue({ items: [detailOrder], next_cursor: null });
    getOrder.mockResolvedValue(detailOrder);

    renderPage("?order=PO1");

    await screen.findByText("Parafuso M8x40 cx100");
    expect(screen.getAllByText("—").length).toBeGreaterThan(0);
    const footerAction = screen.getByRole("button", { name: "Faturar via ERP" });
    expect(footerAction).toBeDisabled();
  });

  it("renders the drawer decomposição/DIFAL block with real formatted values when present (F02-S6 real-ready path)", async () => {
    const withDecomp = {
      ...detailOrder,
      retorno_liquido: 145.67,
      margem_pct: 0.234,
      decomposicao: {
        comissao: 12.5,
        taxa_fixa: 6.5,
        frete: 18.3,
        imposto: 9.4,
        difal: 4.12,
        tarifa_full: null,
        custo: 42.2,
        margem_valor: 154.2,
        margem_pct: 0.234,
        componentes_desconhecidos: ["tarifa_full"],
      },
      difal: {
        amount: 8.76,
        uf_route: "SC → SP",
        due_date: "2026-07-25T00:00:00Z",
        paid: false,
      },
    };
    listOrders.mockResolvedValue({ items: [withDecomp], next_cursor: null });
    getOrder.mockResolvedValue(withDecomp);

    renderPage("?order=PO1");

    await screen.findByText("Parafuso M8x40 cx100");
    // Decomposition + margem + retorno + DIFAL render FROM the data, formatted — proving the
    // real-ready path (same components, no UI change, once the hub wires the decomposer).
    expect(screen.getByText("R$ 145,67")).toBeInTheDocument(); // retorno_liquido
    expect(screen.getByText("R$ 154,20")).toBeInTheDocument(); // margem_valor
    expect(screen.getByText("23.4%")).toBeInTheDocument(); // margem_pct
    expect(screen.getByText("R$ 8,76")).toBeInTheDocument(); // difal.amount
    expect(screen.getByText("R$ 4,12")).toBeInTheDocument(); // decomposicao.difal (cost component)
    expect(screen.getByText("SC → SP")).toBeInTheDocument(); // difal.uf_route
    expect(screen.getByText("não")).toBeInTheDocument(); // difal.paid === false
    // tarifa_full stays null (componentes_desconhecidos) — still honest "—", not fabricated.
    expect(screen.getAllByText("—").length).toBeGreaterThan(0);
  });

  it("closing the drawer clears the ?order= param", async () => {
    listOrders.mockResolvedValue({ items: [detailOrder], next_cursor: null });
    getOrder.mockResolvedValue(detailOrder);

    renderPage("?order=PO1");

    await screen.findByText("Parafuso M8x40 cx100");
    fireEvent.click(screen.getByRole("button", { name: "Fechar detalhe do pedido" }));

    await waitFor(() => {
      expect(screen.getByTestId("location-search")).not.toHaveTextContent("order=PO1");
    });
  });
});
