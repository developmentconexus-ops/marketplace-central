import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { fireEvent, render, screen, waitFor, within } from "@testing-library/react";
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

  it("shows actionable frete guidance on SEM_FRETE — never 'inatingível', never a blank ceiling (L3)", async () => {
    pricingSolveTarget.mockResolvedValue({
      reached: false,
      preco: null,
      ceiling_pct: "",
      desconhecidos: ["frete"],
      blocking_state: null,
      code: "SEM_FRETE",
    });
    renderPanel();

    fireEvent.change(screen.getByLabelText("Margem alvo"), { target: { value: "20" } });
    fireEvent.click(screen.getByTestId("solver-submit"));

    const guidance = await screen.findByTestId("solver-sem-frete");
    expect(guidance).toHaveTextContent("Sem dados de frete");
    expect(guidance).toHaveTextContent("Cadastre dimensões");
    expect(guidance).toHaveTextContent("vincule um anúncio ML");
    expect(guidance).not.toHaveTextContent(/inating/i);
    expect(screen.queryByTestId("solver-unreachable")).toBeNull();
    expect(screen.queryByTestId("solver-price")).toBeNull();
  });

  it("routes a bare desconhecidos:['frete'] into the frete guidance branch even without an explicit code", async () => {
    pricingSolveTarget.mockResolvedValue({
      reached: false,
      preco: null,
      ceiling_pct: "",
      desconhecidos: ["frete"],
      blocking_state: null,
    });
    renderPanel();

    fireEvent.change(screen.getByLabelText("Margem alvo"), { target: { value: "20" } });
    fireEvent.click(screen.getByTestId("solver-submit"));

    expect(await screen.findByTestId("solver-sem-frete")).toBeInTheDocument();
    expect(screen.queryByTestId("solver-unreachable")).toBeNull();
  });

  it("never renders a blank ceiling: unreachable with empty ceiling_pct falls to a generic banner, not '%'", async () => {
    pricingSolveTarget.mockResolvedValue({
      reached: false,
      preco: null,
      ceiling_pct: "",
      desconhecidos: [],
      blocking_state: null,
      code: "UNREACHABLE_TARGET",
    });
    renderPanel();

    fireEvent.change(screen.getByLabelText("Margem alvo"), { target: { value: "40" } });
    fireEvent.click(screen.getByTestId("solver-submit"));

    expect(await screen.findByTestId("solver-incomplete")).toBeInTheDocument();
    // The unreachable branch (the only source of a "…é X%" ceiling) must not render,
    // so an empty ceiling can never surface as a blank "%".
    expect(screen.queryByTestId("solver-unreachable")).toBeNull();
    expect(screen.queryByTestId("solver-price")).toBeNull();
  });

  it("renders the price with fonte/degrau/data badges + ESTIMATIVA pill when the tarifa block is present", async () => {
    // The `tarifa` block is owned by backend chip CHIP-T1 and not yet in the SDK type;
    // cast until PricingSolveResponse carries it (the FE reads it tolerantly).
    pricingSolveTarget.mockResolvedValue({
      reached: true,
      preco: "104.50",
      ceiling_pct: "31.20",
      desconhecidos: [],
      blocking_state: null,
      tarifa: {
        comissao: {
          valor: "13.00",
          fonte: "PADRAO",
          degrau: 4,
          data: "2026-07-19T10:00:00Z",
          estimativa: true,
          sem_dados: false,
        },
        frete: {
          valor: null,
          fonte: "PADRAO",
          degrau: 4,
          data: null,
          estimativa: false,
          sem_dados: true,
        },
      },
    } as PricingSolveResponse);
    renderPanel();

    fireEvent.change(screen.getByLabelText("Margem alvo"), { target: { value: "20" } });
    fireEvent.click(screen.getByTestId("solver-submit"));

    expect(await screen.findByTestId("solver-price")).toHaveTextContent("104.50");

    const comissao = screen.getByTestId("tarifa-comissao");
    expect(comissao).toHaveTextContent("13.00%");
    expect(comissao).toHaveTextContent("Padrão");
    expect(comissao).toHaveTextContent("degrau 4");
    expect(comissao).toHaveTextContent("ESTIMATIVA");
    // Resolved value carries a freshness carimbo (data).
    expect(within(comissao).getByLabelText("Atualização dos dados")).toBeInTheDocument();

    // NO-DATA frete renders UnknownValue "—", never R$0, never an ESTIMATIVA pill,
    // and suppresses the source/degrau carimbos beside the unknown.
    const frete = screen.getByTestId("tarifa-frete");
    expect(frete).toHaveTextContent("—");
    expect(frete).not.toHaveTextContent("R$ 0");
    expect(frete).not.toHaveTextContent("ESTIMATIVA");
    expect(frete).not.toHaveTextContent("degrau");
  });

  it("treats an empty-string tarifa valor as NO-DATA — never a blank 'R$ ' with a false ESTIMATIVA pill", async () => {
    pricingSolveTarget.mockResolvedValue({
      reached: true,
      preco: "95.00",
      ceiling_pct: "20.00",
      desconhecidos: [],
      blocking_state: null,
      tarifa: {
        comissao: { valor: "13.00", fonte: "PADRAO", degrau: 4, data: null, estimativa: true, sem_dados: false },
        frete: { valor: "", fonte: "PADRAO", degrau: 4, data: null, estimativa: false, sem_dados: false },
      },
    } as PricingSolveResponse);
    renderPanel();

    fireEvent.change(screen.getByLabelText("Margem alvo"), { target: { value: "20" } });
    fireEvent.click(screen.getByTestId("solver-submit"));

    await screen.findByTestId("solver-price");
    const frete = screen.getByTestId("tarifa-frete");
    expect(frete).toHaveTextContent("—");
    expect(frete).not.toHaveTextContent("R$");
    expect(frete).not.toHaveTextContent("ESTIMATIVA");
  });

  it("does not dress an incomplete-data code (DADOS_INCOMPLETOS) with a ceiling as 'inatingível X%'", async () => {
    pricingSolveTarget.mockResolvedValue({
      reached: false,
      preco: null,
      ceiling_pct: "18.50",
      desconhecidos: [],
      blocking_state: null,
      code: "DADOS_INCOMPLETOS",
    });
    renderPanel();

    fireEvent.change(screen.getByLabelText("Margem alvo"), { target: { value: "40" } });
    fireEvent.click(screen.getByTestId("solver-submit"));

    expect(await screen.findByTestId("solver-incomplete")).toBeInTheDocument();
    expect(screen.queryByTestId("solver-unreachable")).toBeNull();
    // The plausible-looking 18.5% ceiling must not be surfaced for incomplete data.
    expect(screen.queryByText(/18\.50/)).toBeNull();
  });

  it("shows the resolved price when reached=true even if a stale banner code lingers in the response", async () => {
    pricingSolveTarget.mockResolvedValue({
      reached: true,
      preco: "95.00",
      ceiling_pct: "22.00",
      desconhecidos: [],
      blocking_state: null,
      code: "SEM_FRETE",
    });
    renderPanel();

    fireEvent.change(screen.getByLabelText("Margem alvo"), { target: { value: "20" } });
    fireEvent.click(screen.getByTestId("solver-submit"));

    expect(await screen.findByTestId("solver-price")).toHaveTextContent("95.00");
    expect(screen.queryByTestId("solver-sem-frete")).toBeNull();
  });

  it("clears a prior result when the selected product changes, so it never renders under a new product", async () => {
    pricingSolveTarget.mockResolvedValue({
      reached: true,
      preco: "104.50",
      ceiling_pct: "31.20",
      desconhecidos: [],
      blocking_state: null,
    });
    const { rerender } = renderPanel({ productId: 90001 });

    fireEvent.change(screen.getByLabelText("Margem alvo"), { target: { value: "20" } });
    fireEvent.click(screen.getByTestId("solver-submit"));
    expect(await screen.findByTestId("solver-price")).toHaveTextContent("104.50");

    rerender(
      <QueryClientProvider client={new QueryClient({ defaultOptions: { queries: { retry: false } } })}>
        <SolverPanel productId={90002} comissaoPct="17" modalidade="premium" />
      </QueryClientProvider>,
    );

    expect(screen.queryByTestId("solver-price")).toBeNull();
  });

  it("renders the price without badges (no crash) when the tarifa block is missing — tolerant parsing", async () => {
    pricingSolveTarget.mockResolvedValue({
      reached: true,
      preco: "88.00",
      ceiling_pct: "40.00",
      desconhecidos: [],
      blocking_state: null,
    });
    renderPanel();

    fireEvent.change(screen.getByLabelText("Margem alvo"), { target: { value: "20" } });
    fireEvent.click(screen.getByTestId("solver-submit"));

    expect(await screen.findByTestId("solver-price")).toHaveTextContent("88.00");
    expect(screen.queryByTestId("tarifa-comissao")).toBeNull();
    expect(screen.queryByTestId("tarifa-frete")).toBeNull();
  });
});
