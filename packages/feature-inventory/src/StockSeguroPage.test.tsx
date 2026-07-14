import { fireEvent, render as testingRender, screen, waitFor } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { describe, expect, it, vi } from "vitest";
import { StockSeguroPage, type StockSeguroClient } from "./StockSeguroPage";

const baseRender = testingRender;
function render(ui: React.ReactNode) {
  return baseRender(<QueryClientProvider client={new QueryClient({ defaultOptions: { queries: { retry: 0 } } })}>{ui}</QueryClientProvider>);
}

function makeClient(items: any[], overrides?: Partial<StockSeguroClient>): StockSeguroClient {
  return {
    listIntegrationInstallations: async () => ({
      items: [
        {
          installation_id: "inst-1",
          tenant_id: "tenant_default",
          provider_code: "mercado_livre",
          family: "marketplace",
          display_name: "Mercado Livre",
          status: "connected",
          health_status: "healthy",
          external_account_id: "acct-1",
          external_account_name: "Acct 1",
          connection: {
            state: "connected",
            health: "healthy",
            provider_code: "mercado_livre",
            external_account_id: "acct-1",
            external_account_name: "Acct 1",
            auth_strategy: "oauth2",
            next_action: "none",
          },
          runtime_capabilities: [],
          created_at: "2026-07-09T12:00:00Z",
          updated_at: "2026-07-09T12:00:00Z",
        },
      ],
    }),
    listInventoryStockRisks: async () => ({ items }),
    applyInventoryManualStockAction: async () => ({
      action: {
        action_id: "act-1",
        state: "applied",
        trigger: "manual",
        provider_ref: {
          tenant_id: "tenant_default",
          installation_id: "inst-1",
          provider_code: "mercado_livre",
          provider_account_id: "acct-1",
          provider_item_id: "MLB123",
        },
        requested_quantity: 7,
        recommended_quantity: 7,
        policy_id: "stock-seguro-default",
        operator: { actor_type: "operator", actor_id: "leandro" },
        idempotency_key: "act-1",
        audit_events: [],
        created_at: "2026-07-09T12:00:00Z",
        updated_at: "2026-07-09T12:00:00Z",
      },
      risk: items[0],
    }),
    ...overrides,
  };
}

const oversellItem = {
  identity: { installation_id: "inst-1", provider_item_id: "MLB123" },
  provider_code: "mercado_livre",
  title: "Produto testado",
  link_state: "resolved",
  state: "oversell",
  actionability: "actionable",
  actionable: true,
  internal_quantity: 8,
  provider_quantity: 9,
  recommended_quantity: 7,
  policy_id: "stock-seguro-default",
};

const undersellItem = {
  ...oversellItem,
  identity: { installation_id: "inst-1", provider_item_id: "MLB126" },
  state: "undersell",
  provider_quantity: 6,
  recommended_quantity: 7,
};

const healthyItem = {
  ...oversellItem,
  identity: { installation_id: "inst-1", provider_item_id: "MLB127" },
  state: "healthy",
  actionability: "blocked",
  actionable: false,
  provider_quantity: 7,
  recommended_quantity: 7,
};

const ineligibleItem = {
  ...oversellItem,
  identity: { installation_id: "inst-1", provider_item_id: "MLB128" },
  state: "ineligible",
  actionability: "blocked",
  actionable: false,
  blocking_reason: { code: "ineligible_product", message: "showroom-only group" },
};

describe("StockSeguroPage", () => {
  it("renders loading and then oversell state", async () => {
    render(
      <MemoryRouter initialEntries={["/inventory/stock-seguro?installation=inst-1"]}>
        <StockSeguroPage client={makeClient([oversellItem])} />
      </MemoryRouter>,
    );

    expect(screen.getByText("Loading Stock Seguro...")).toBeInTheDocument();
    expect((await screen.findAllByText("Oversell")).length).toBeGreaterThan(0);
    expect(screen.getAllByText("Produto testado").length).toBeGreaterThan(0);
  });

  it("renders empty state", async () => {
    render(
      <MemoryRouter initialEntries={["/inventory/stock-seguro?installation=inst-1"]}>
        <StockSeguroPage client={makeClient([])} />
      </MemoryRouter>,
    );

    expect(await screen.findByText("No Stock Seguro rows for the selected filters.")).toBeInTheDocument();
  });

  it("renders error state", async () => {
    render(
      <MemoryRouter initialEntries={["/inventory/stock-seguro?installation=inst-1"]}>
        <StockSeguroPage
          client={makeClient([], {
            listInventoryStockRisks: async () => {
              throw new Error("inventory unavailable");
            },
          })}
        />
      </MemoryRouter>,
    );

    expect(await screen.findByText("inventory unavailable")).toBeInTheDocument();
  });

  it("renders blocked states", async () => {
    render(
      <MemoryRouter initialEntries={["/inventory/stock-seguro?installation=inst-1"]}>
        <StockSeguroPage
          client={makeClient([
            { ...oversellItem, state: "conflict", link_state: "conflict", actionability: "blocked", actionable: false, blocking_reason: { code: "conflict_link", message: "product link is in conflict" } },
            { ...oversellItem, state: "unresolved", link_state: "unresolved", actionability: "blocked", actionable: false, blocking_reason: { code: "unresolved_link", message: "product link is not resolved" }, identity: { installation_id: "inst-1", provider_item_id: "MLB124" } },
            { ...oversellItem, state: "stale", actionability: "blocked", actionable: false, blocking_reason: { code: "stale_provider_source", message: "source_older_than_policy" }, identity: { installation_id: "inst-1", provider_item_id: "MLB125" } },
          ])}
        />
      </MemoryRouter>,
    );

    expect((await screen.findAllByText("Conflito")).length).toBeGreaterThan(0);
    expect(screen.getAllByText("Sem vinculo").length).toBeGreaterThan(0);
    expect(screen.getAllByText("Stale").length).toBeGreaterThan(0);
  });

  it("renders healthy, undersell, and ineligible states and counts all blocked rows", async () => {
    render(
      <MemoryRouter initialEntries={["/inventory/stock-seguro?installation=inst-1"]}>
        <StockSeguroPage client={makeClient([undersellItem, healthyItem, ineligibleItem])} />
      </MemoryRouter>,
    );

    expect((await screen.findAllByText("Undersell")).length).toBeGreaterThan(0);
    expect(screen.getAllByText("Healthy").length).toBeGreaterThan(0);
    expect(screen.getAllByText("Ineligivel").length).toBeGreaterThan(0);
    expect(screen.getByTestId("stock-summary-total")).toHaveTextContent("3");
    expect(screen.getByTestId("stock-summary-actionable")).toHaveTextContent("1");
    expect(screen.getByTestId("stock-summary-blockers")).toHaveTextContent("2");
  });

  it("applies manual action and shows result message", async () => {
    const applyInventoryManualStockAction = vi.fn(async () => ({
      action: {
        action_id: "act-1",
        state: "applied",
        trigger: "manual",
        provider_ref: {
          tenant_id: "tenant_default",
          installation_id: "inst-1",
          provider_code: "mercado_livre",
          provider_account_id: "acct-1",
          provider_item_id: "MLB123",
        },
        requested_quantity: 7,
        recommended_quantity: 7,
        policy_id: "stock-seguro-default",
        operator: { actor_type: "operator", actor_id: "leandro" },
        idempotency_key: "act-1",
        audit_events: [],
        created_at: "2026-07-09T12:00:00Z",
        updated_at: "2026-07-09T12:00:00Z",
      },
      risk: { ...oversellItem, state: "healthy", actionable: false, actionability: "blocked", provider_quantity: 7 },
    }));
    render(
      <MemoryRouter initialEntries={["/inventory/stock-seguro?installation=inst-1"]}>
        <StockSeguroPage client={makeClient([oversellItem], { applyInventoryManualStockAction })} />
      </MemoryRouter>,
    );

    expect((await screen.findAllByText("Oversell")).length).toBeGreaterThan(0);
    fireEvent.click(screen.getByRole("button", { name: "Apply recommended quantity" }));
    expect(await screen.findByRole("button", { name: "Confirm apply" })).toBeInTheDocument();
    fireEvent.click(screen.getByRole("button", { name: "Confirm apply" }));

    await waitFor(() => expect(applyInventoryManualStockAction).toHaveBeenCalled());
    expect(await screen.findByText("Action applied for MLB123.")).toBeInTheDocument();
  });
});
