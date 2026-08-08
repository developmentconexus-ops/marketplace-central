import { fireEvent, render as testingRender, screen, waitFor } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { describe, expect, it, vi } from "vitest";
import { queryKeyNamespaces } from "@marketplace-central/web-query";
import type { IntegrationInstallation } from "@marketplace-central/sdk-runtime";
import { StockSeguroPage, type StockSeguroClient } from "./StockSeguroPage";

const baseRender = testingRender;
function render(
  ui: React.ReactNode,
  queryClient = new QueryClient({ defaultOptions: { queries: { retry: 0 } } }),
) {
  return baseRender(<QueryClientProvider client={queryClient}>{ui}</QueryClientProvider>);
}

function makeClient(items: any[], overrides?: Partial<StockSeguroClient>): StockSeguroClient {
  return {
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

const installations: IntegrationInstallation[] = [
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
];

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
        <StockSeguroPage client={makeClient([oversellItem])} installations={installations} />
      </MemoryRouter>,
    );

    expect(screen.getByText("Carregando Stock Seguro...")).toBeInTheDocument();
    expect((await screen.findAllByText("Oversell")).length).toBeGreaterThan(0);
    expect((await screen.findAllByText("Produto testado")).length).toBeGreaterThan(0);
  });

  it("renders empty state", async () => {
    render(
      <MemoryRouter initialEntries={["/inventory/stock-seguro?installation=inst-1"]}>
        <StockSeguroPage client={makeClient([])} installations={installations} />
      </MemoryRouter>,
    );

    expect(
      await screen.findByText("Nenhuma linha de Stock Seguro para os filtros selecionados."),
    ).toBeInTheDocument();
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
          installations={installations}
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
            {
              ...oversellItem,
              state: "conflict",
              link_state: "conflict",
              actionability: "blocked",
              actionable: false,
              blocking_reason: { code: "conflict_link", message: "product link is in conflict" },
            },
            {
              ...oversellItem,
              state: "unresolved",
              link_state: "unresolved",
              actionability: "blocked",
              actionable: false,
              blocking_reason: { code: "unresolved_link", message: "product link is not resolved" },
              identity: { installation_id: "inst-1", provider_item_id: "MLB124" },
            },
            {
              ...oversellItem,
              state: "stale",
              actionability: "blocked",
              actionable: false,
              blocking_reason: {
                code: "stale_provider_source",
                message: "source_older_than_policy",
              },
              identity: { installation_id: "inst-1", provider_item_id: "MLB125" },
            },
          ])}
          installations={installations}
        />
      </MemoryRouter>,
    );

    expect((await screen.findAllByText("Conflito")).length).toBeGreaterThan(0);
    expect(screen.getAllByText("Sem vínculo").length).toBeGreaterThan(0);
    expect(screen.getAllByText("Desatualizado").length).toBeGreaterThan(0);
  });

  it("names the observation a stale-source block was judged on", async () => {
    render(
      <MemoryRouter initialEntries={["/inventory/stock-seguro?installation=inst-1"]}>
        <StockSeguroPage
          client={makeClient([
            {
              ...oversellItem,
              state: "stale",
              actionability: "blocked",
              actionable: false,
              internal_observed_at: "2026-07-18T20:23:27Z",
              blocking_reason: {
                code: "stale_internal_source",
                message: "source_older_than_policy",
              },
            },
          ])}
          installations={installations}
        />
      </MemoryRouter>,
    );

    // Without the date the operator cannot tell whether the snapshot is an hour
    // or a week old, nor which side to refresh.
    expect(await screen.findByText(/Observado em 18\/07\/2026/)).toBeInTheDocument();
  });

  it("renders healthy, undersell, and ineligible states and counts all blocked rows", async () => {
    render(
      <MemoryRouter initialEntries={["/inventory/stock-seguro?installation=inst-1"]}>
        <StockSeguroPage
          client={makeClient([undersellItem, healthyItem, ineligibleItem])}
          installations={installations}
        />
      </MemoryRouter>,
    );

    expect((await screen.findAllByText("Undersell")).length).toBeGreaterThan(0);
    expect(screen.getAllByText("Saudável").length).toBeGreaterThan(0);
    expect(screen.getAllByText("Inelegível").length).toBeGreaterThan(0);
    expect(await screen.findByTestId("stock-summary-total")).toHaveTextContent("3");
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
      risk: {
        ...oversellItem,
        state: "healthy",
        actionable: false,
        actionability: "blocked",
        provider_quantity: 7,
      },
    }));
    const queryClient = new QueryClient({ defaultOptions: { queries: { retry: 0 } } });
    const invalidateQueries = vi.spyOn(queryClient, "invalidateQueries");
    render(
      <MemoryRouter initialEntries={["/inventory/stock-seguro?installation=inst-1"]}>
        <StockSeguroPage
          client={makeClient([oversellItem], { applyInventoryManualStockAction })}
          installations={installations}
        />
      </MemoryRouter>,
      queryClient,
    );

    expect((await screen.findAllByText("Oversell")).length).toBeGreaterThan(0);
    fireEvent.click(await screen.findByRole("button", { name: "Aplicar quantidade recomendada" }));
    expect(await screen.findByRole("button", { name: "Confirmar aplicação" })).toBeInTheDocument();
    fireEvent.change(screen.getByPlaceholderText("Seu nome"), { target: { value: "Leandro" } });
    fireEvent.click(screen.getByRole("button", { name: "Confirmar aplicação" }));

    await waitFor(() => expect(applyInventoryManualStockAction).toHaveBeenCalled());
    expect(await screen.findByText("Ação aplicada para MLB123.")).toBeInTheDocument();
    expect(invalidateQueries).toHaveBeenCalledTimes(1);
    expect(invalidateQueries).toHaveBeenCalledWith({ queryKey: queryKeyNamespaces.inventory });
  });

  it("does not invalidate inventory after a failed stock action", async () => {
    const queryClient = new QueryClient({ defaultOptions: { queries: { retry: 0 } } });
    const invalidateQueries = vi.spyOn(queryClient, "invalidateQueries");
    render(
      <MemoryRouter initialEntries={["/inventory/stock-seguro?installation=inst-1"]}>
        <StockSeguroPage
          client={makeClient([oversellItem], {
            applyInventoryManualStockAction: vi
              .fn()
              .mockRejectedValue(new Error("stock action failed")),
          })}
          installations={installations}
        />
      </MemoryRouter>,
      queryClient,
    );

    fireEvent.click(await screen.findByRole("button", { name: "Aplicar quantidade recomendada" }));
    fireEvent.change(await screen.findByPlaceholderText("Seu nome"), {
      target: { value: "Leandro" },
    });
    fireEvent.click(await screen.findByRole("button", { name: "Confirmar aplicação" }));
    await waitFor(() => expect(screen.getByText("stock action failed")).toBeInTheDocument());
    expect(invalidateQueries).not.toHaveBeenCalled();
  });

  it("blocks the apply until the approving operator is named", async () => {
    const applyInventoryManualStockAction = vi.fn();
    render(
      <MemoryRouter initialEntries={["/inventory/stock-seguro?installation=inst-1"]}>
        <StockSeguroPage
          client={makeClient([oversellItem], {
            applyInventoryManualStockAction: applyInventoryManualStockAction as never,
          })}
          installations={installations}
        />
      </MemoryRouter>,
    );

    fireEvent.click(await screen.findByRole("button", { name: "Aplicar quantidade recomendada" }));
    const confirm = await screen.findByRole("button", { name: "Confirmar aplicação" });
    expect(confirm).toBeDisabled();
    fireEvent.click(confirm);
    expect(applyInventoryManualStockAction).not.toHaveBeenCalled();

    fireEvent.change(screen.getByPlaceholderText("Seu nome"), { target: { value: "Leandro" } });
    expect(screen.getByRole("button", { name: "Confirmar aplicação" })).toBeEnabled();
  });

  it("translates the blocking reason by code and falls back to the API message", async () => {
    render(
      <MemoryRouter initialEntries={["/inventory/stock-seguro?installation=inst-1"]}>
        <StockSeguroPage
          client={makeClient([
            oversellItem,
            {
              ...oversellItem,
              identity: { installation_id: "inst-1", provider_item_id: "MLB129" },
              blocking_reason: { code: "unresolved_link", message: "product link is not resolved" },
            },
            {
              ...oversellItem,
              identity: { installation_id: "inst-1", provider_item_id: "MLB130" },
              blocking_reason: {
                code: "code_the_screen_does_not_know",
                message: "raw engineering detail",
              },
            },
          ])}
          installations={installations}
        />
      </MemoryRouter>,
    );

    expect(
      await screen.findByText("Anúncio sem vínculo com um produto do ERP."),
    ).toBeInTheDocument();
    expect(screen.queryByText("product link is not resolved")).not.toBeInTheDocument();
    expect(screen.getByText("raw engineering detail")).toBeInTheDocument();
  });
});
