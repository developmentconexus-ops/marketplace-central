import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { fireEvent, render, screen, waitFor, within } from "@testing-library/react";
import { MemoryRouter, useLocation } from "react-router-dom";
import { beforeEach, describe, expect, it, vi } from "vitest";
import { PedidosPage } from "./PedidosPage";

const listOrders = vi.fn();
const getOrderSummary = vi.fn();
const getOrder = vi.fn();
const importMarketplaceOrders = vi.fn();
const markOrderFaturado = vi.fn();

vi.mock("../../app/ClientContext", () => ({
  useClient: () => ({
    listOrders: (...args: unknown[]) => listOrders(...args),
    getOrderSummary: (...args: unknown[]) => getOrderSummary(...args),
    getOrder: (...args: unknown[]) => getOrder(...args),
    importMarketplaceOrders: (...args: unknown[]) => importMarketplaceOrders(...args),
    markOrderFaturado: (...args: unknown[]) => markOrderFaturado(...args),
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
    importMarketplaceOrders.mockReset();
    markOrderFaturado.mockReset();
    getOrderSummary.mockResolvedValue({
      today: 0,
      seven_days: 0,
      by_status: { novo: 2, faturar: 1, enviar: 3, enviado: 4 },
    });
  });

  it("lists every order by default, sending no date filter", async () => {
    listOrders.mockResolvedValue({ items: [], next_cursor: null });

    renderPage();

    await screen.findByRole("heading", { name: "Pedidos" });
    // "todo o período" must not invent a window: date_from stays undefined so the
    // read-list returns the orders that exist.
    await waitFor(() => expect(listOrders).toHaveBeenCalled());
    expect(listOrders.mock.calls[0][0]).toMatchObject({
      installation_id: "inst_1",
      date_from: undefined,
    });
    expect(screen.getByRole("button", { name: /ENVIADOS/ })).toHaveTextContent("todo o período");
  });

  it("applies the chosen period as a date_from filter and relabels the ENVIADOS window", async () => {
    listOrders.mockResolvedValue({ items: [], next_cursor: null });

    renderPage();
    await screen.findByRole("heading", { name: "Pedidos" });
    await waitFor(() => expect(listOrders).toHaveBeenCalled());

    fireEvent.change(screen.getByLabelText("Período dos pedidos"), { target: { value: "7" } });

    await waitFor(() => {
      const last = listOrders.mock.calls[listOrders.mock.calls.length - 1][0];
      expect(typeof last.date_from).toBe("string");
      // Seven days back from now, so the filter matches the label the KPI claims.
      const days = (Date.now() - Date.parse(last.date_from)) / 86_400_000;
      expect(days).toBeGreaterThan(6.9);
      expect(days).toBeLessThan(7.1);
    });
    expect(screen.getByRole("button", { name: /ENVIADOS/ })).toHaveTextContent("últimos 7d");
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

  it("derives the KPI bucket counts from the loaded list (KPI == Lista) and gates DIFAL as unknown", async () => {
    // Counts come from OrderRead.bucket over the full list — the same source the Lista tabs count,
    // so a card and its tab always agree. Bucket split here: novo=2, faturar=1, enviar=3, enviado=4.
    listOrders.mockResolvedValue({
      items: [
        { ...baseOrder, provider_order_id: "A", bucket: "novo" as const },
        { ...baseOrder, provider_order_id: "B", bucket: "novo" as const },
        { ...baseOrder, provider_order_id: "C", bucket: "faturar" as const },
        { ...baseOrder, provider_order_id: "D", bucket: "enviar" as const },
        { ...baseOrder, provider_order_id: "E", bucket: "enviar" as const },
        { ...baseOrder, provider_order_id: "F", bucket: "enviar" as const },
        { ...baseOrder, provider_order_id: "G", bucket: "enviado" as const },
        { ...baseOrder, provider_order_id: "H", bucket: "enviado" as const },
        { ...baseOrder, provider_order_id: "I", bucket: "enviado" as const },
        { ...baseOrder, provider_order_id: "J", bucket: "enviado" as const },
      ],
      next_cursor: null,
    });

    renderPage();

    // KpiCard is a button whose accessible name is "LABEL VALUE" — assert the derived counts.
    expect(await screen.findByRole("button", { name: "NOVOS 2" })).toBeInTheDocument();
    expect(screen.getByRole("button", { name: "A FATURAR 1" })).toBeInTheDocument();
    expect(screen.getByRole("button", { name: "A ENVIAR 3" })).toBeInTheDocument();
    expect(screen.getByRole("button", { name: "ENVIADOS 4" })).toBeInTheDocument();
    // DIFAL stays an honest unknown (decomposição not wired) — never a fabricated 0.
    expect(screen.getByRole("button", { name: /DIFAL A PAGAR —/ })).toBeInTheDocument();
  });

  it("renders '—' for every bucket KPI while the order list is unavailable", async () => {
    listOrders.mockRejectedValue({ status: 503, error: { code: "unavailable" } });

    renderPage();

    await screen.findByRole("heading", { name: "Pedidos" });
    // List error → the bucket counts are unknown, so all four bucket cards + DIFAL render "—",
    // never a fabricated 0 (ADR-17).
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
    expect(await screen.findByText("PO1")).toBeInTheDocument();
    // The Fila view is captioned (design has no "Fila de trabalho" heading; the view toggle sits
    // in the title row), so assert the Fila tab is the selected view instead.
    expect(screen.getByRole("tab", { name: "Fila", selected: true })).toBeInTheDocument();

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

    // Rows are labelled by the order number (provider_order_id), not the channel slug
    // (provider_code), so filtering shows/hides by "PO1"/"PO2".
    expect(await screen.findByText("PO1")).toBeInTheDocument();
    expect(screen.queryByText("PO2")).not.toBeInTheDocument();

    fireEvent.click(screen.getByRole("tab", { name: "A faturar 1" }));

    await waitFor(() => {
      expect(screen.getByText("PO2")).toBeInTheDocument();
      expect(screen.queryByText("PO1")).not.toBeInTheDocument();
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

    expect(await screen.findByText("PO1")).toBeInTheDocument();
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

    expect(await screen.findByText("PO1")).toBeInTheDocument();
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
    await screen.findByText("PO1");

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

    await screen.findByText("PO1");
    expect(screen.getByText("sem ação")).toBeInTheDocument();
    expect(screen.queryByRole("button", { name: "Faturar" })).not.toBeInTheDocument();
  });

  it("renders the golden real Fila row (2000017336572246) from real payload with honest '—' for itens/retorno", async () => {
    // Golden case from the chip contract: a real ML order whose items:[] is empty and whose cost
    // components are all unknown. The Fila row description mirrors the design's descDe —
    // comprador · itens · valor — so with no item titles it reads "Marcia Rocha · R$ 339,98"
    // (destinatário is the recipient; buyer nickname is null by ML privacy). retorno is an honest
    // "—" and there is no fabricated item count or margin. City/UF and status live in the drawer,
    // not the Fila row (design parity).
    const golden = {
      ...baseOrder,
      provider_order_id: "2000017336572246",
      // Real payload: provider_code is the channel slug, not the order number — the num column must
      // show provider_order_id (regression guard: live data showed "mercado_livre" here).
      provider_code: "mercado_livre",
      buyer: { display: "M. R.", city: "São Paulo", uf: "SP" },
      destinatario: "Marcia Rocha",
      destino_uf: "SP",
      total: 339.98,
      bucket: "enviado" as const,
      sla: { due: "2026-07-11T00:00:00Z", atrasado: false },
      items: [],
      retorno_liquido: null,
    };
    listOrders.mockResolvedValue({ items: [golden], next_cursor: null });

    renderPage();

    // Fila is the default view — no tab switch needed.
    await screen.findByRole("heading", { name: "Pedidos" });
    // Composite description (comprador · valor); itens dropped honestly since items:[] is empty.
    expect(await screen.findByText("Marcia Rocha · R$ 339,98")).toBeInTheDocument();
    // retorno_liquido:null → honest "—" for retorno, scoped to THIS Fila row (not the page): the
    // assertion must pin FilaRetorno's UnknownValue, not be satisfied by the unconditional DIFAL
    // KPI dash elsewhere on the page. A regression to a fabricated "R$ 0,00" here must fail this.
    const goldenRow = screen.getByRole("button", { name: "Abrir detalhe do pedido 2000017336572246" });
    expect(within(goldenRow).getByText("—")).toBeInTheDocument();
    // num column shows the order number, never the channel slug.
    expect(within(goldenRow).getByText("2000017336572246")).toBeInTheDocument();
    expect(screen.queryByText("mercado_livre")).not.toBeInTheDocument();
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

    expect(await screen.findByText("PO1")).toBeInTheDocument();
    expect(await screen.findByText("PO2")).toBeInTheDocument();
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
    fireEvent.click(await screen.findByText("PO1"));

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

  it("titles the drawer with the order number, never the channel slug (provider_code)", async () => {
    // Regression guard (hub corrective #1): on real ML data provider_code is the channel slug
    // "mercado_livre" — the drawer title must be the order number (provider_order_id), matching the
    // 4 fixed views. A regression to `order?.provider_code || orderId` must fail this.
    const goldenDetail = {
      ...detailOrder,
      provider_order_id: "2000017336572246",
      provider_code: "mercado_livre",
    };
    listOrders.mockResolvedValue({ items: [goldenDetail], next_cursor: null });
    getOrder.mockResolvedValue(goldenDetail);

    renderPage("?order=2000017336572246");

    const drawer = await screen.findByRole("complementary");
    // Wait for order DETAIL to resolve — until then `title` falls back to orderId for EITHER code
    // path, so the guard is only load-bearing once loaded data is present (buyer comes from getOrder).
    await within(drawer).findByText("J. S.");
    // The drawer body renders neither provider_code nor provider_order_id — only the title does — so
    // the order number inside the drawer IS the title, and the slug can only appear if the title
    // regressed to provider_code. Both assertions fail on a `provider_code || orderId` regression.
    expect(within(drawer).getByText("2000017336572246")).toBeInTheDocument();
    expect(screen.queryByText("mercado_livre")).not.toBeInTheDocument();
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

  it("renders the drawer buyer fiscal identity + carrier/destino/frete_real when present", async () => {
    const withFiscal = {
      ...detailOrder,
      destinatario: "Maria Souza",
      destino_uf: "SP",
      destino_cep: "01310-100",
      frete_real: { bruto: 23.45, receiver: 10.0, sender: 13.45 },
      rastreio: {
        shipment_id: "SHIP-1",
        status: "shipped",
        substatus: "out_for_delivery",
        transportadora: "Total Express",
        url_rastreio: "https://tracking.totalexpress.com.br/x?id=3",
      },
      comprador_fiscal: {
        nome: "Maria Souza LTDA",
        doc_tipo: "CNPJ", // rendered verbatim (opaque, ADR-17)
        doc_numero: "11222333000181",
        endereco: {
          logradouro: "Av. Paulista",
          numero: "1000",
          cidade: "São Paulo",
          uf_codigo: "SP",
          uf_nome: "São Paulo",
          cep: "01310-100",
          pais: "BR",
        },
      },
    };
    listOrders.mockResolvedValue({ items: [withFiscal], next_cursor: null });
    getOrder.mockResolvedValue(withFiscal);

    renderPage("?order=PO1");

    await screen.findByText("Parafuso M8x40 cx100");
    const drawer = within(screen.getByRole("complementary"));
    // Carrier + tracking link (round-4 additive).
    expect(drawer.getByText("Total Express")).toBeInTheDocument();
    const trackLink = drawer.getByRole("link", { name: /rastrear/i });
    expect(trackLink).toHaveAttribute("href", "https://tracking.totalexpress.com.br/x?id=3");
    // Destino + destinatário (round-4 additive).
    expect(drawer.getByText("Maria Souza")).toBeInTheDocument(); // destinatario
    expect(drawer.getByText("SP · 01310-100")).toBeInTheDocument(); // destino_uf · destino_cep
    // Frete real actuals.
    expect(drawer.getByText("R$ 23,45")).toBeInTheDocument(); // frete_real.bruto
    // Buyer fiscal identity for ERP registration.
    expect(drawer.getByText("Maria Souza LTDA")).toBeInTheDocument(); // nome
    expect(drawer.getByText(/CNPJ/)).toBeInTheDocument(); // doc_tipo opaque
    expect(drawer.getByText(/11222333000181/)).toBeInTheDocument(); // doc_numero
    expect(drawer.getByText(/Av\. Paulista/)).toBeInTheDocument(); // endereco logradouro
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
    // Scope to the drawer: the Fila row also renders retorno_liquido, so query within the
    // drawer container (role="complementary") to assert the decomposição/DIFAL block specifically.
    const drawer = within(screen.getByRole("complementary"));
    // Decomposition + margem + retorno + DIFAL render FROM the data, formatted — proving the
    // real-ready path (same components, no UI change, once the hub wires the decomposer).
    expect(drawer.getByText("R$ 145,67")).toBeInTheDocument(); // retorno_liquido
    expect(drawer.getByText("R$ 154,20")).toBeInTheDocument(); // margem_valor
    expect(drawer.getByText("23.4%")).toBeInTheDocument(); // margem_pct
    expect(drawer.getByText("R$ 8,76")).toBeInTheDocument(); // difal.amount
    expect(drawer.getByText("R$ 4,12")).toBeInTheDocument(); // decomposicao.difal (cost component)
    expect(drawer.getByText("SC → SP")).toBeInTheDocument(); // difal.uf_route
    expect(drawer.getByText("não")).toBeInTheDocument(); // difal.paid === false
    // tarifa_full stays null (componentes_desconhecidos) — still honest "—", not fabricated.
    expect(drawer.getAllByText("—").length).toBeGreaterThan(0);
  });

  it("Atualizar re-runs the orders import then refetches the list", async () => {
    listOrders.mockResolvedValue({ items: [], next_cursor: null });
    importMarketplaceOrders.mockResolvedValue({ imported_count: 0, skipped_count: 0 });

    renderPage();

    // Wait for the initial load to settle — while ordersQuery is fetching the button reads
    // "Atualizando…"; it only shows "Atualizar" once idle.
    const button = await screen.findByRole("button", { name: "Atualizar" });
    const callsBefore = listOrders.mock.calls.length;
    fireEvent.click(button);

    await waitFor(() =>
      expect(importMarketplaceOrders).toHaveBeenCalledWith({ installation_id: "inst_1" }),
    );
    // onSuccess refetches the list, so listOrders is invoked again beyond the initial load.
    await waitFor(() => expect(listOrders.mock.calls.length).toBeGreaterThan(callsBefore));
  });

  it("drawer 'Foi faturado' marks a faturar order (our-DB) and invalidates the orders queries", async () => {
    const faturarOrder = { ...detailOrder, bucket: "faturar" as const };
    listOrders.mockResolvedValue({ items: [faturarOrder], next_cursor: null });
    getOrder.mockResolvedValue(faturarOrder);
    markOrderFaturado.mockResolvedValue(undefined);

    renderPage("?order=PO1");

    await screen.findByText("Parafuso M8x40 cx100");
    const drawer = within(screen.getByRole("complementary"));
    fireEvent.click(drawer.getByRole("button", { name: "Foi faturado" }));

    await waitFor(() => expect(markOrderFaturado).toHaveBeenCalledWith("inst_1", "PO1"));
    // Invalidating the "orders" namespace re-derives the drawer bucket + the list/KPI counts.
    await waitFor(() => expect(getOrder.mock.calls.length).toBeGreaterThanOrEqual(2));
  });

  it("drawer omits 'Foi faturado' for a non-faturar order (novo)", async () => {
    listOrders.mockResolvedValue({ items: [detailOrder], next_cursor: null });
    getOrder.mockResolvedValue(detailOrder);

    renderPage("?order=PO1");

    await screen.findByText("Parafuso M8x40 cx100");
    expect(screen.queryByRole("button", { name: "Foi faturado" })).not.toBeInTheDocument();
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
