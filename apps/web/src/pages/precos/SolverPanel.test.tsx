import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import { beforeEach, describe, expect, it, vi } from "vitest";
import type { PricingSolveResponse } from "@marketplace-central/sdk-runtime";
import { SolverPanel } from "./SolverPanel";

const pricingSolveTarget = vi.fn<(input: unknown) => Promise<PricingSolveResponse>>();

vi.mock("../../app/ClientContext", () => ({
  useClient: () => ({ pricingSolveTarget }),
}));

function renderPanel(props: Partial<React.ComponentProps<typeof SolverPanel>> = {}) {
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  return render(
    <QueryClientProvider client={queryClient}>
      <SolverPanel productId={90001} comissaoPct="17" modalidade="premium" {...props} />
    </QueryClientProvider>,
  );
}

describe("SolverPanel — margem-alvo → preço", () => {
  beforeEach(() => {
    pricingSolveTarget.mockReset();
  });

  it("disables the solve action until a target margin is entered", () => {
    renderPanel();
    expect(screen.getByTestId("solver-submit")).toBeDisabled();
    fireEvent.click(screen.getByTestId("solver-submit"));
    expect(pricingSolveTarget).not.toHaveBeenCalled();
  });

  it("disables the solve action when no product is selected", () => {
    renderPanel({ productId: null });
    fireEvent.change(screen.getByLabelText("Margem alvo"), { target: { value: "20" } });
    expect(screen.getByTestId("solver-submit")).toBeDisabled();
  });

  it("solves the target margin into a suggested price", async () => {
    pricingSolveTarget.mockResolvedValue({
      reached: true,
      preco: "104.50",
      ceiling_pct: "31.20",
      desconhecidos: [],
      blocking_state: null,
    });
    renderPanel();

    fireEvent.change(screen.getByLabelText("Margem alvo"), { target: { value: "20" } });
    fireEvent.click(screen.getByTestId("solver-submit"));

    await waitFor(() => expect(pricingSolveTarget).toHaveBeenCalledTimes(1));
    expect(pricingSolveTarget.mock.calls[0][0]).toMatchObject({
      margem_alvo_pct: "20",
      comissao_pct: "17",
      modalidade: "premium",
      product_id: 90001,
    });
    expect(await screen.findByTestId("solver-price")).toHaveTextContent("104.50");
    expect(screen.queryByTestId("solver-unreachable")).toBeNull();
  });

  it("surfaces the achievable ceiling when the target is unreachable, without fabricating a price", async () => {
    pricingSolveTarget.mockResolvedValue({
      reached: false,
      preco: null,
      ceiling_pct: "12.80",
      desconhecidos: [],
      blocking_state: null,
      code: "UNREACHABLE_TARGET",
    });
    renderPanel();

    fireEvent.change(screen.getByLabelText("Margem alvo"), { target: { value: "40" } });
    fireEvent.click(screen.getByTestId("solver-submit"));

    const unreachable = await screen.findByTestId("solver-unreachable");
    expect(unreachable).toHaveTextContent("12.80");
    expect(screen.queryByTestId("solver-price")).toBeNull();
  });

  it("shows the SEM_CUSTO blocking banner and no price when cost is unresolved", async () => {
    pricingSolveTarget.mockResolvedValue({
      reached: false,
      preco: null,
      ceiling_pct: "0",
      desconhecidos: ["custo"],
      blocking_state: "SEM_CUSTO",
    });
    renderPanel();

    fireEvent.change(screen.getByLabelText("Margem alvo"), { target: { value: "20" } });
    fireEvent.click(screen.getByTestId("solver-submit"));

    expect(await screen.findByTestId("solver-blocking")).toHaveTextContent("SEM_CUSTO");
    expect(screen.queryByTestId("solver-price")).toBeNull();
  });

  it("normalizes a pt-BR comma target to dot-decimal before the solve SDK call (F-P7-2)", async () => {
    pricingSolveTarget.mockResolvedValue({
      reached: true,
      preco: "104.50",
      ceiling_pct: "31.20",
      desconhecidos: [],
      blocking_state: null,
    });
    renderPanel();

    fireEvent.change(screen.getByLabelText("Margem alvo"), { target: { value: "9,5" } });
    fireEvent.click(screen.getByTestId("solver-submit"));

    await waitFor(() => expect(pricingSolveTarget).toHaveBeenCalledTimes(1));
    expect(pricingSolveTarget.mock.calls[0][0]).toMatchObject({ margem_alvo_pct: "9.5" });
  });

  it("renders an error state when the solve call fails", async () => {
    pricingSolveTarget.mockRejectedValue(new Error("boom"));
    renderPanel();

    fireEvent.change(screen.getByLabelText("Margem alvo"), { target: { value: "20" } });
    fireEvent.click(screen.getByTestId("solver-submit"));

    expect(await screen.findByRole("button", { name: /tentar novamente/i })).toBeInTheDocument();
  });
});
