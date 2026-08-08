import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { fireEvent, render, screen, waitFor, within } from "@testing-library/react";
import { beforeEach, describe, expect, it, vi } from "vitest";
import type {
  PricingScenario,
  PricingScenarioListResponse,
} from "@marketplace-central/sdk-runtime";
import { ScenariosPanel } from "./ScenariosPanel";

const listPricingScenarios = vi.fn<() => Promise<PricingScenarioListResponse>>();
const createPricingScenario = vi.fn<(input: unknown) => Promise<PricingScenario>>();
const deletePricingScenario = vi.fn<(id: string) => Promise<void>>();

vi.mock("../../app/ClientContext", () => ({
  useClient: () => ({ listPricingScenarios, createPricingScenario, deletePricingScenario }),
}));

const scenarioA: PricingScenario = {
  id: "cenario-a",
  name: "Cenário A",
  payload: { preco: "89.00", modalidade: "premium", product_id: 90001 },
  created_at: "2026-07-18T12:00:00Z",
};

function renderPanel(props: Partial<React.ComponentProps<typeof ScenariosPanel>> = {}) {
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  const onReload = props.onReload ?? vi.fn();
  const utils = render(
    <QueryClientProvider client={queryClient}>
      <ScenariosPanel
        payload={{ preco: "104.50", modalidade: "classico", product_id: 90002 }}
        onReload={onReload}
        {...props}
      />
    </QueryClientProvider>,
  );
  return { ...utils, onReload };
}

describe("ScenariosPanel — save / reload / delete", () => {
  beforeEach(() => {
    listPricingScenarios.mockReset();
    createPricingScenario.mockReset();
    deletePricingScenario.mockReset();
  });

  it("lists saved scenarios", async () => {
    listPricingScenarios.mockResolvedValue({ items: [scenarioA] });
    renderPanel();
    expect(await screen.findByText("Cenário A")).toBeInTheDocument();
  });

  it("shows an empty hint when there are no saved scenarios", async () => {
    listPricingScenarios.mockResolvedValue({ items: [] });
    renderPanel();
    expect(await screen.findByTestId("scenarios-empty")).toBeInTheDocument();
  });

  it("saves the current simulation under a name and refreshes the list", async () => {
    listPricingScenarios
      .mockResolvedValueOnce({ items: [] })
      .mockResolvedValue({ items: [scenarioA] });
    createPricingScenario.mockResolvedValue(scenarioA);
    renderPanel();

    await screen.findByTestId("scenarios-empty");
    fireEvent.change(screen.getByLabelText("Nome do cenário"), { target: { value: "Cenário A" } });
    fireEvent.click(screen.getByTestId("scenario-save"));

    await waitFor(() => expect(createPricingScenario).toHaveBeenCalledTimes(1));
    const req = createPricingScenario.mock.calls[0][0] as {
      id: string;
      name: string;
      payload: Record<string, unknown>;
    };
    expect(req.name).toBe("Cenário A");
    expect(req.id).toBeTruthy();
    expect(req.payload).toMatchObject({
      preco: "104.50",
      modalidade: "classico",
      product_id: 90002,
    });
    // List re-fetched → the saved scenario appears.
    expect(await screen.findByText("Cenário A")).toBeInTheDocument();
  });

  it("reloads a saved scenario's payload via the onReload callback", async () => {
    listPricingScenarios.mockResolvedValue({ items: [scenarioA] });
    const { onReload } = renderPanel();

    fireEvent.click(await screen.findByTestId("scenario-reload-cenario-a"));
    expect(onReload).toHaveBeenCalledWith(scenarioA.payload);
  });

  it("deletes a saved scenario and refreshes the list", async () => {
    listPricingScenarios
      .mockResolvedValueOnce({ items: [scenarioA] })
      .mockResolvedValue({ items: [] });
    deletePricingScenario.mockResolvedValue(undefined);
    renderPanel();

    const row = await screen.findByTestId("scenario-row-cenario-a");
    fireEvent.click(within(row).getByTestId("scenario-delete-cenario-a"));

    await waitFor(() => expect(deletePricingScenario).toHaveBeenCalledWith("cenario-a"));
    await waitFor(() => expect(screen.getByTestId("scenarios-empty")).toBeInTheDocument());
  });

  it("renders an error state when the list fails to load", async () => {
    listPricingScenarios.mockRejectedValue(new Error("boom"));
    renderPanel();
    expect(await screen.findByRole("button", { name: /tentar novamente/i })).toBeInTheDocument();
  });
});
