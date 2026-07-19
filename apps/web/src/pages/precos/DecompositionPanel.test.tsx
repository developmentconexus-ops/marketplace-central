import { render, screen, within } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import type {
  PricingCalcProfile,
  PricingDecomposition,
} from "@marketplace-central/sdk-runtime";
import { DecompositionPanel } from "./DecompositionPanel";

const profile: PricingCalcProfile = {
  regime: "simples",
  aliquota_pct: "4.00",
  limiar_verde_pct: "18",
  limiar_amarelo_pct: "10",
  tarifa_full: "12",
  difal_enabled: false,
  difal_destino_uf: null,
  origem: "seed",
};

function decomposition(overrides: Partial<PricingDecomposition> = {}): PricingDecomposition {
  return {
    preco: "100.00",
    comissao: "17.00",
    taxa_fixa: "0.00",
    frete: "12.00",
    imposto: "4.00",
    difal: null,
    tarifa_full: null,
    custo: "40.00",
    margem_valor: "27.00",
    margem_pct: "27.00",
    componentes_desconhecidos: [],
    ...overrides,
  };
}

describe("DecompositionPanel", () => {
  it("renders each IC-04 component and the resulting margin band", () => {
    render(<DecompositionPanel decomposition={decomposition()} profile={profile} blockingState={null} />);

    const panel = screen.getByTestId("decomposition-panel");
    expect(within(panel).getByText("100.00")).toBeInTheDocument();
    expect(within(panel).getByText("17.00")).toBeInTheDocument();
    expect(within(panel).getByText("40.00")).toBeInTheDocument();
    // margem 27% >= verde 18 ⇒ healthy band
    expect(within(panel).getByText("27%")).toHaveClass("bg-accent-soft");
  });

  it("renders unknown components as — and NEVER as 0 (ADR-17)", () => {
    const d = decomposition({
      difal: null,
      tarifa_full: null,
      componentes_desconhecidos: ["difal", "tarifa_full"],
    });
    render(<DecompositionPanel decomposition={d} profile={profile} blockingState={null} />);

    const panel = screen.getByTestId("decomposition-panel");
    // Unknown lines (difal, tarifa_full) render the em-dash placeholder, never a
    // fabricated 0. The em-dash spans carry text-faint (UnknownValue), unlike the
    // known "0.00" taxa_fixa which is a real value in text-ink.
    const dashes = within(panel).getAllByText("—", { selector: "span.text-faint" });
    expect(dashes.length).toBeGreaterThanOrEqual(2);
  });

  it("shows the SEM_CUSTO blocking banner and never fabricates a margin when cost is unknown", () => {
    const d = decomposition({
      custo: null,
      margem_valor: null,
      margem_pct: null,
      componentes_desconhecidos: ["custo"],
    });
    render(<DecompositionPanel decomposition={d} profile={profile} blockingState="SEM_CUSTO" />);

    expect(screen.getByTestId("blocking-sem-custo")).toBeInTheDocument();
    expect(screen.getByRole("alert")).toHaveTextContent("SEM_CUSTO");
    // Margin chip falls back to the unknown band, not a 0% claim.
    const panel = screen.getByTestId("decomposition-panel");
    expect(within(panel).queryByText("0%")).toBeNull();
  });

  it("classifies the margin by the operator's CalcProfile thresholds, not hard-coded bands", () => {
    // margem 12% is below the operator's verde 18 but at/above amarelo 10 ⇒ tight.
    const tightProfile: PricingCalcProfile = { ...profile, limiar_verde_pct: "18", limiar_amarelo_pct: "10" };
    const d = decomposition({ margem_pct: "12.00", margem_valor: "12.00" });
    render(<DecompositionPanel decomposition={d} profile={tightProfile} blockingState={null} />);

    expect(screen.getByText("12%")).toHaveClass("bg-amber-soft");
  });

  it("warns that the margin excludes DIFAL when DIFAL is disabled", () => {
    render(<DecompositionPanel decomposition={decomposition()} profile={profile} blockingState={null} />);
    const warning = screen.getByTestId("difal-off-warning");
    expect(warning).toHaveTextContent(/sem DIFAL/i);
    expect(warning).toHaveTextContent(/não use/i);
  });

  it("does not show the DIFAL-off warning when DIFAL is enabled", () => {
    const difalOn: PricingCalcProfile = { ...profile, difal_enabled: true, difal_destino_uf: "SP" };
    render(
      <DecompositionPanel
        decomposition={decomposition({ difal: "3.50" })}
        profile={difalOn}
        blockingState={null}
        difalUf="SP"
      />,
    );
    expect(screen.queryByTestId("difal-off-warning")).toBeNull();
  });

  it("labels the DIFAL line with the destination UF when provided", () => {
    render(
      <DecompositionPanel
        decomposition={decomposition({ difal: "3.50" })}
        profile={profile}
        blockingState={null}
        difalUf="SP"
      />,
    );
    expect(screen.getByText("(−) DIFAL SP")).toBeInTheDocument();
  });
});
