import { fireEvent, render as testingRender, screen, waitFor, within } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import { QueryClientProvider } from "@tanstack/react-query";
import { describe, expect, it, vi } from "vitest";
import { createWebQueryClient, queryKeyNamespaces } from "@marketplace-central/web-query";
import { OrdersPage, type OrdersClient, type ProfitabilityActor } from "./OrdersPage";

const baseRender = testingRender;
function render(ui: React.ReactNode) {
  return baseRender(<QueryClientProvider client={createWebQueryClient()}>{ui}</QueryClientProvider>);
}

function makeInstallation() {
  return {
    installation_id: "inst-1",
    tenant_id: "tenant_default",
    provider_code: "mercado_livre",
    family: "marketplace" as const,
    display_name: "Mercado Livre Primary",
    status: "connected" as const,
    health_status: "healthy" as const,
    external_account_id: "691607102",
    external_account_name: "Metal Nobre",
    connection: {
      state: "connected" as const,
      health: "healthy" as const,
      provider_code: "mercado_livre",
      external_account_id: "691607102",
      external_account_name: "Metal Nobre",
      auth_strategy: "oauth2" as const,
      next_action: "none" as const,
    },
    runtime_capabilities: [],
    created_at: "2026-07-09T12:00:00Z",
    updated_at: "2026-07-09T12:00:00Z",
  };
}

function makeOrders() {
  return [
    {
      installation_id: "inst-1",
      provider_code: "mercado_livre",
      provider_order_id: "2000012659424976",
      provider_status: "paid",
      provider_updated_at: "2026-07-09T12:00:00Z",
      fetched_at: "2026-07-09T12:00:00Z",
      shipping_id: "ship-1",
      items: [
        {
          provider_item_id: "MLB123",
          title: "Pedido #2000012659424976",
          quantity: 1,
          unit_price: 19.9,
          link_quality: "resolved" as const,
        },
      ],
      payments: [],
    },
    {
      installation_id: "inst-1",
      provider_code: "mercado_livre",
      provider_order_id: "2000012659424977",
      provider_status: "paid",
      provider_updated_at: "2026-07-09T13:00:00Z",
      fetched_at: "2026-07-09T13:00:00Z",
      items: [
        {
          provider_item_id: "MLB124",
          title: "Complete order",
          quantity: 1,
          unit_price: 50,
          link_quality: "resolved" as const,
        },
      ],
      payments: [],
    },
    {
      installation_id: "inst-1",
      provider_code: "mercado_livre",
      provider_order_id: "2000012659424978",
      provider_status: "paid",
      provider_updated_at: "2026-07-09T14:00:00Z",
      fetched_at: "2026-07-09T14:00:00Z",
      items: [
        {
          provider_item_id: "MLB125",
          title: "Negative order",
          quantity: 1,
          unit_price: 20,
          link_quality: "resolved" as const,
        },
      ],
      payments: [],
    },
    {
      installation_id: "inst-1",
      provider_code: "mercado_livre",
      provider_order_id: "2000012659424979",
      provider_status: "cancelled",
      provider_updated_at: "2026-07-09T15:00:00Z",
      fetched_at: "2026-07-09T15:00:00Z",
      items: [
        {
          provider_item_id: "MLB126",
          title: "Cancelled order",
          quantity: 1,
          unit_price: 30,
          link_quality: "resolved" as const,
        },
      ],
      payments: [],
    },
    {
      installation_id: "inst-1",
      provider_code: "mercado_livre",
      provider_order_id: "2000012659424980",
      provider_status: "pending",
      provider_updated_at: "2026-07-09T16:00:00Z",
      fetched_at: "2026-07-09T16:00:00Z",
      items: [
        {
          provider_item_id: "MLB127",
          title: "Unknown realization order",
          quantity: 1,
          unit_price: 40,
          link_quality: "resolved" as const,
        },
      ],
      payments: [],
    },
  ];
}

function makeInputs() {
  return [
    {
      installation_id: "inst-1",
      provider_order_id: "2000012659424976",
      scope: "order" as const,
      kind: "freight" as const,
      amount: 12.5,
      currency: "BRL",
      source_system: "manual_adjustment",
      quality: "manual" as const,
      created_at: "2026-07-09T12:00:00Z",
      updated_at: "2026-07-09T12:00:00Z",
    },
    {
      installation_id: "inst-1",
      provider_order_id: "2000012659424977",
      scope: "order" as const,
      kind: "revenue" as const,
      amount: 50,
      currency: "BRL",
      source_system: "orders",
      quality: "complete" as const,
      created_at: "2026-07-09T12:00:00Z",
      updated_at: "2026-07-09T12:00:00Z",
    },
  ];
}

function makeSnapshots() {
  return [
    {
      installation_id: "inst-1",
      provider_order_id: "2000012659424976",
      scope: "order" as const,
      currency: "BRL",
      realization_state: "realized" as const,
      revenue_amount: 19.9,
      sale_fee_amount: 9.63,
      freight_amount: 12.5,
      contribution_amount: -2.23,
      margin_percent: -11.21,
      quality: "incomplete" as const,
      flags: ["missing_link", "missing_commission"] as const,
      created_at: "2026-07-09T12:00:00Z",
      updated_at: "2026-07-09T12:00:00Z",
    },
    {
      installation_id: "inst-1",
      provider_order_id: "2000012659424976",
      provider_item_id: "MLB123",
      scope: "item" as const,
      currency: "BRL",
      realization_state: "realized" as const,
      revenue_amount: 19.9,
      contribution_amount: -2.23,
      margin_percent: -11.21,
      quality: "incomplete" as const,
      flags: ["missing_link", "missing_cost"] as const,
      created_at: "2026-07-09T12:00:00Z",
      updated_at: "2026-07-09T12:00:00Z",
    },
    {
      installation_id: "inst-1",
      provider_order_id: "2000012659424977",
      scope: "order" as const,
      currency: "BRL",
      realization_state: "realized" as const,
      revenue_amount: 50,
      contribution_amount: 10,
      margin_percent: 20,
      quality: "complete" as const,
      created_at: "2026-07-09T12:00:00Z",
      updated_at: "2026-07-09T12:00:00Z",
    },
    {
      installation_id: "inst-1",
      provider_order_id: "2000012659424977",
      provider_item_id: "MLB124",
      scope: "item" as const,
      currency: "BRL",
      realization_state: "realized" as const,
      revenue_amount: 50,
      contribution_amount: 10,
      margin_percent: 20,
      quality: "complete" as const,
      created_at: "2026-07-09T12:00:00Z",
      updated_at: "2026-07-09T12:00:00Z",
    },
    {
      installation_id: "inst-1",
      provider_order_id: "2000012659424978",
      scope: "order" as const,
      currency: "BRL",
      realization_state: "realized" as const,
      revenue_amount: 20,
      contribution_amount: -8,
      margin_percent: -40,
      quality: "negative_margin" as const,
      flags: ["negative_margin"] as const,
      created_at: "2026-07-09T12:00:00Z",
      updated_at: "2026-07-09T12:00:00Z",
    },
    {
      installation_id: "inst-1",
      provider_order_id: "2000012659424979",
      scope: "order" as const,
      currency: "BRL",
      realization_state: "not_realized" as const,
      revenue_amount: 30,
      contribution_amount: null,
      margin_percent: null,
      quality: "not_realized" as const,
      flags: ["order_cancelled"] as const,
      created_at: "2026-07-09T15:00:00Z",
      updated_at: "2026-07-09T15:00:00Z",
    },
    {
      installation_id: "inst-1",
      provider_order_id: "2000012659424980",
      scope: "order" as const,
      currency: "BRL",
      realization_state: "unknown" as const,
      revenue_amount: 40,
      contribution_amount: null,
      margin_percent: null,
      quality: "incomplete" as const,
      flags: ["order_state_unknown"] as const,
      created_at: "2026-07-09T16:00:00Z",
      updated_at: "2026-07-09T16:00:00Z",
    },
  ];
}

function makeAdjustments() {
  return [
    {
      adjustment_id: "adj-1",
		idempotency_key: "manual-adjustment-1",
      installation_id: "inst-1",
      provider_order_id: "2000012659424976",
      scope: "order" as const,
      category: "freight" as const,
      amount: 12.5,
      currency: "BRL",
      reason: "Mercado Envios adjustment",
      actor: { actor_type: "operator", actor_id: "lea", actor_name: "Lea" },
      created_at: "2026-07-09T12:00:00Z",
    },
  ];
}

function makeClient(overrides: Partial<OrdersClient> = {}): OrdersClient {
  const orders = makeOrders();
  const inputs = makeInputs();
  const snapshots = makeSnapshots();
  const adjustments = makeAdjustments();

  return {
    listIntegrationInstallations: vi.fn().mockResolvedValue({ items: [makeInstallation()] }),
    importMarketplaceOrders: vi.fn().mockResolvedValue({}),
    listMarketplaceOrders: vi.fn().mockResolvedValue({ items: orders }),
    importProfitabilityMarginInputs: vi.fn().mockResolvedValue({}),
    listProfitabilityMarginInputs: vi.fn().mockResolvedValue({ items: inputs }),
    createProfitabilityManualAdjustment: vi.fn().mockResolvedValue({
      adjustment_id: "adj-2",
		idempotency_key: "manual-adjustment-2",
      installation_id: "inst-1",
      provider_order_id: "2000012659424976",
      scope: "order",
      category: "commission",
      amount: 4,
      currency: "BRL",
      reason: "Operator verified commission",
      actor: { actor_type: "operator", actor_id: "leandro", actor_name: "Leandro" },
      created_at: "2026-07-09T12:05:00Z",
    }),
    listProfitabilityManualAdjustments: vi.fn().mockResolvedValue({ items: adjustments }),
    calculateProfitabilityProfitSnapshots: vi.fn().mockResolvedValue({}),
    listProfitabilityProfitSnapshots: vi.fn().mockResolvedValue({ items: snapshots }),
    ...overrides,
  };
}

function renderPage(clientOverrides: Partial<OrdersClient> = {}, route = "/orders?installation=inst-1", operator?: ProfitabilityActor, queryClient = createWebQueryClient()) {
  const client = makeClient(clientOverrides);
  render(
    <MemoryRouter initialEntries={[route]}>
      <QueryClientProvider client={queryClient}>
        <OrdersPage client={client} operator={operator} />
      </QueryClientProvider>
    </MemoryRouter>,
  );
  return { client, queryClient };
}

function getSelectedOrderSurface(): HTMLElement {
  const surface = screen.getByText("Selected order").closest("section");
  if (!surface) {
    throw new Error("Selected order surface was not rendered.");
  }
  return surface;
}

function expectFieldValue(container: HTMLElement, label: string, value: string): void {
  const fieldLabel = within(container).getByText(label, { selector: "p" });
  expect(fieldLabel.nextElementSibling).toHaveTextContent(new RegExp(`^${value}$`));
}

describe("OrdersPage", () => {
  it("requires an actor for manual adjustment requests", () => {
    type AdjustmentRequest = Parameters<OrdersClient["createProfitabilityManualAdjustment"]>[0];
    expectTypeOf<AdjustmentRequest>().toMatchTypeOf<{
		idempotency_key: string;
      actor: { actor_type: string; actor_id: string; actor_name?: string };
    }>();
  });

  it("disables adjustment creation when trustworthy operator identity is unavailable", async () => {
    renderPage();

    expect((await screen.findAllByText("Order #2000012659424976")).length).toBeGreaterThan(0);
    expect(screen.getByText(/operator identity is unavailable/i)).toBeInTheDocument();
    expect(screen.queryByLabelText("Adjustment operator")).not.toBeInTheDocument();
    expect(screen.getByRole("button", { name: "Save adjustment" })).toBeDisabled();
  });

  it("disables adjustment creation when required operator identity is blank", async () => {
    renderPage({}, "/orders?installation=inst-1", { actor_type: "operator", actor_id: "   " });

    expect((await screen.findAllByText("Order #2000012659424976")).length).toBeGreaterThan(0);
    expect(screen.getByText(/operator identity is unavailable/i)).toBeInTheDocument();
    expect(screen.getByRole("button", { name: "Save adjustment" })).toBeDisabled();
  });

  it("forwards the immutable operator provenance exactly", async () => {
    const createProfitabilityManualAdjustment = vi.fn().mockResolvedValue(makeAdjustments()[0]);
    renderPage(
      { createProfitabilityManualAdjustment },
      "/orders?installation=inst-1",
      { actor_type: "operator", actor_id: "user-42", actor_name: "Trusted Operator" },
    );

    expect((await screen.findAllByText("Order #2000012659424976")).length).toBeGreaterThan(0);
    fireEvent.change(screen.getByLabelText("Adjustment amount"), { target: { value: "4" } });
    fireEvent.change(screen.getByLabelText("Adjustment reason"), { target: { value: "Verified freight" } });
    fireEvent.click(screen.getByRole("button", { name: "Save adjustment" }));

    await waitFor(() => expect(createProfitabilityManualAdjustment).toHaveBeenCalled());
    expect(createProfitabilityManualAdjustment.mock.calls[0][0].actor).toEqual({
      actor_type: "operator",
      actor_id: "user-42",
      actor_name: "Trusted Operator",
    });
	expect(createProfitabilityManualAdjustment.mock.calls[0][0].idempotency_key).toEqual(expect.any(String));
  });

  it("reuses an idempotency key only for an exact failed-submission retry", async () => {
    const createProfitabilityManualAdjustment = vi.fn().mockRejectedValue(new Error("temporary outage"));
    renderPage(
      { createProfitabilityManualAdjustment },
      "/orders?installation=inst-1",
      { actor_type: "operator", actor_id: "user-42", actor_name: "Trusted Operator" },
    );

    expect((await screen.findAllByText("Order #2000012659424976")).length).toBeGreaterThan(0);
    fireEvent.change(screen.getByLabelText("Adjustment amount"), { target: { value: "4" } });
    fireEvent.change(screen.getByLabelText("Adjustment reason"), { target: { value: "Verified freight" } });

    fireEvent.click(screen.getByRole("button", { name: "Save adjustment" }));
    await waitFor(() => expect(createProfitabilityManualAdjustment).toHaveBeenCalledTimes(1));
    fireEvent.click(screen.getByRole("button", { name: "Save adjustment" }));
    await waitFor(() => expect(createProfitabilityManualAdjustment).toHaveBeenCalledTimes(2));

    fireEvent.click(await screen.findByRole("button", { name: /Order #2000012659424977/ }));
    fireEvent.click(screen.getByRole("button", { name: "Save adjustment" }));
    await waitFor(() => expect(createProfitabilityManualAdjustment).toHaveBeenCalledTimes(3));

    const [first, exactRetry, changedSubmission] = createProfitabilityManualAdjustment.mock.calls.map(([request]) => request);
    expect(exactRetry.idempotency_key).toBe(first.idempotency_key);
    expect(changedSubmission.idempotency_key).not.toBe(first.idempotency_key);
  });

  it("renders loading and then the orders workspace summary", async () => {
    renderPage({
      listMarketplaceOrders: vi.fn().mockImplementation(
        () => new Promise((resolve) => setTimeout(() => resolve({ items: makeOrders() }), 10)),
      ),
    });

    expect(screen.getByText("Loading orders workspace...")).toBeInTheDocument();
    expect((await screen.findAllByText("Order #2000012659424976")).length).toBeGreaterThan(0);
  });

  it("renders empty state", async () => {
    renderPage({
      listMarketplaceOrders: vi.fn().mockResolvedValue({ items: [] }),
      listProfitabilityProfitSnapshots: vi.fn().mockResolvedValue({ items: [] }),
      listProfitabilityManualAdjustments: vi.fn().mockResolvedValue({ items: [] }),
      listProfitabilityMarginInputs: vi.fn().mockResolvedValue({ items: [] }),
    });

    expect(await screen.findByText("No orders available for the selected installation.")).toBeInTheDocument();
  });

  it("renders error state", async () => {
    renderPage({
      listMarketplaceOrders: vi.fn().mockRejectedValue({ error: { message: "orders unavailable" } }),
    });

    expect(await screen.findByText("orders unavailable")).toBeInTheDocument();
  });

  it("renders incomplete state and manual adjustment history", async () => {
    renderPage();

    expect((await screen.findAllByText("Order #2000012659424976")).length).toBeGreaterThan(0);
    expect(screen.getAllByText("Incomplete").length).toBeGreaterThan(0);
    expect(screen.getByText("Data quality")).toBeInTheDocument();
    expect(screen.getAllByText("Manual adjustments").length).toBeGreaterThan(0);
    expect(screen.getByText("Mercado Envios adjustment")).toBeInTheDocument();
    expect(screen.getAllByText("freight").length).toBeGreaterThan(0);
  });

  it("renders complete and negative-margin states and filters incomplete orders", async () => {
    renderPage();

    expect((await screen.findAllByText("Order #2000012659424977")).length).toBeGreaterThan(0);
    expect(screen.getAllByText("Complete").length).toBeGreaterThan(0);
    expect(screen.getAllByText("Negative margin").length).toBeGreaterThan(0);

    fireEvent.change(screen.getByLabelText("Margin quality"), { target: { value: "incomplete" } });

    await waitFor(() => {
      expect(screen.getAllByText("Order #2000012659424976").length).toBeGreaterThan(0);
      expect(screen.queryAllByText("Order #2000012659424977")).toHaveLength(0);
    });
  });

  it("renders cancelled orders as not realized without profitability values", async () => {
    renderPage();

    const cancelledOrder = await screen.findByRole("button", { name: /Order #2000012659424979/ });
    expect(within(screen.getByLabelText("Margin quality")).getByRole("option", { name: "Not realized" })).toBeInTheDocument();
    const notRealizedStatLabel = screen.getByText("Not realized", { selector: "p" });
    expect(notRealizedStatLabel.nextElementSibling).toHaveTextContent(/^1$/);
    expect(within(cancelledOrder).getByText("Not realized")).toBeInTheDocument();
    expect(within(cancelledOrder).getByText("Order cancelled")).toBeInTheDocument();
    expect(within(cancelledOrder).getByText("—")).toBeInTheDocument();
    expect(within(cancelledOrder).queryByText("Negative margin")).not.toBeInTheDocument();

    fireEvent.click(cancelledOrder);
    expect(await screen.findByText("Order not realized")).toBeInTheDocument();
    const selectedCancelledOrder = getSelectedOrderSurface();
    expectFieldValue(selectedCancelledOrder, "Contribution", "—");
    expectFieldValue(selectedCancelledOrder, "Margin", "—");

    fireEvent.change(screen.getByLabelText("Margin quality"), { target: { value: "not_realized" } });
    await waitFor(() => {
      expect(screen.getAllByText("Order #2000012659424979").length).toBeGreaterThan(0);
      expect(screen.queryAllByText("Order #2000012659424978")).toHaveLength(0);
    });
  });

  it("renders unknown realization as incomplete operational state", async () => {
    renderPage();

    const unknownOrder = await screen.findByRole("button", { name: /Order #2000012659424980/ });
    expect(within(unknownOrder).getByText("Incomplete")).toBeInTheDocument();
    expect(within(unknownOrder).getByText("Order state unknown")).toBeInTheDocument();
    expect(within(unknownOrder).getByText("—")).toBeInTheDocument();

    fireEvent.click(unknownOrder);
    await waitFor(() => {
      expect(within(getSelectedOrderSurface()).getByRole("heading", { name: "Order #2000012659424980" })).toBeInTheDocument();
    });
    const selectedUnknownOrder = getSelectedOrderSurface();
    expectFieldValue(selectedUnknownOrder, "Contribution", "—");
    expectFieldValue(selectedUnknownOrder, "Margin", "—");
  });

  it("saves a manual adjustment and refreshes the workspace", async () => {
    const adjustmentList = [...makeAdjustments()];
    const snapshotList = [...makeSnapshots()];
    const createProfitabilityManualAdjustment = vi.fn().mockImplementation(async () => {
      const created = {
        adjustment_id: "adj-2",
		idempotency_key: "manual-adjustment-2",
        installation_id: "inst-1",
        provider_order_id: "2000012659424976",
        scope: "order" as const,
        category: "commission" as const,
        amount: 4,
        currency: "BRL",
        reason: "Operator verified commission",
        actor: { actor_type: "operator", actor_id: "leandro", actor_name: "Leandro" },
        created_at: "2026-07-09T12:05:00Z",
      };
      adjustmentList.unshift(created);
      return created;
    });
    const calculateProfitabilityProfitSnapshots = vi.fn().mockImplementation(async () => {
      snapshotList[0] = {
        ...snapshotList[0],
        commission_amount: 4,
        contribution_amount: -6.23,
        flags: ["missing_link"] as const,
      };
      return {};
    });

    renderPage(
      {
        createProfitabilityManualAdjustment,
        calculateProfitabilityProfitSnapshots,
        listProfitabilityManualAdjustments: vi.fn().mockImplementation(async () => ({ items: adjustmentList })),
        listProfitabilityProfitSnapshots: vi.fn().mockImplementation(async () => ({ items: snapshotList })),
      },
      "/orders?installation=inst-1",
      { actor_type: "operator", actor_id: "leandro", actor_name: "Leandro" },
    );

    expect((await screen.findAllByText("Order #2000012659424976")).length).toBeGreaterThan(0);
    fireEvent.change(screen.getByLabelText("Adjustment category"), { target: { value: "commission" } });
    fireEvent.change(screen.getByLabelText("Adjustment amount"), { target: { value: "4" } });
    fireEvent.change(screen.getByLabelText("Adjustment reason"), { target: { value: "Operator verified commission" } });
    fireEvent.click(screen.getByRole("button", { name: "Save adjustment" }));

    await waitFor(() => expect(createProfitabilityManualAdjustment).toHaveBeenCalled());
    await waitFor(() => expect(calculateProfitabilityProfitSnapshots).toHaveBeenCalled());
    expect(await screen.findByText("Manual adjustment saved for order 2000012659424976.")).toBeInTheDocument();
    expect(screen.getByText("Operator verified commission")).toBeInTheDocument();
  });

  it("invalidates catalog and profitability after importing margin inputs", async () => {
    const queryClient = createWebQueryClient();
    const invalidateQueries = vi.spyOn(queryClient, "invalidateQueries");
    const importProfitabilityMarginInputs = vi.fn().mockResolvedValue({});
    renderPage({ importProfitabilityMarginInputs }, "/orders?installation=inst-1", undefined, queryClient);

    await screen.findAllByText("Order #2000012659424976");
    fireEvent.click(screen.getByRole("button", { name: "Import inputs" }));

    await waitFor(() => expect(importProfitabilityMarginInputs).toHaveBeenCalledWith({ installation_id: "inst-1", limit: 50 }));
    expect(invalidateQueries).toHaveBeenCalledTimes(2);
    expect(invalidateQueries).toHaveBeenCalledWith({ queryKey: queryKeyNamespaces.catalog });
    expect(invalidateQueries).toHaveBeenCalledWith({ queryKey: queryKeyNamespaces.profitability });
  });

  it("does not invalidate after a failed margin-input import", async () => {
    const queryClient = createWebQueryClient();
    const invalidateQueries = vi.spyOn(queryClient, "invalidateQueries");
    renderPage({ importProfitabilityMarginInputs: vi.fn().mockRejectedValue(new Error("import failed")) }, "/orders?installation=inst-1", undefined, queryClient);

    await screen.findAllByText("Order #2000012659424976");
    fireEvent.click(screen.getByRole("button", { name: "Import inputs" }));

    await waitFor(() => expect(screen.getByText("import failed")).toBeInTheDocument());
    expect(invalidateQueries).not.toHaveBeenCalled();
  });
});
