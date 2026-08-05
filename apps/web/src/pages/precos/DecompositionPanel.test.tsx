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

  // Task A7: imposto is now nullable (null ⇒ the D-41 per-cell fiscal path
  // is active for this item; the legacy flat rate is never shown alongside
  // icms_saida/difal/pis_cofins). This asserts the RENDERED STATE of the
  // imposto row specifically (its own placeholder span, not merely "a dash
  // exists somewhere on the panel" — difal/tarifa_full are also null by
  // default in the `decomposition()` fixture, so a panel-wide dash count
  // would pass even if imposto itself fell back to a fabricated value).
  it("renders imposto as the unknown-value placeholder (not a fabricated 0) when null", () => {
    const d = decomposition({ imposto: null, difal: "3.50", tarifa_full: "0.00" });
    render(<DecompositionPanel decomposition={d} profile={profile} blockingState={null} />);

    const label = screen.getByText("(−) Imposto");
    const row = label.parentElement as HTMLElement;
    const placeholder = within(row).getByText("—");
    expect(placeholder).toHaveClass("text-faint");
    // No numeric value (fabricated or otherwise) renders in the imposto row.
    expect(row).not.toHaveTextContent("4.00");
    expect(row).not.toHaveTextContent("0.00");
  });

  it("renders imposto as the legacy known value (unchanged) when non-null — negative control for the null case above", () => {
    const d = decomposition({ imposto: "4.00" });
    render(<DecompositionPanel decomposition={d} profile={profile} blockingState={null} />);

    const label = screen.getByText("(−) Imposto");
    const row = label.parentElement as HTMLElement;
    expect(within(row).getByText("4.00")).toBeInTheDocument();
    expect(within(row).queryByText("—")).toBeNull();
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

  it("renders the comissão fonte/degrau carimbo (COTACAO degrau 3, no ESTIMATIVA) beside the value", () => {
    render(
      <DecompositionPanel
        decomposition={decomposition()}
        profile={profile}
        blockingState={null}
        tarifa={{
          comissao: { valor: "14.5", fonte: "COTACAO", degrau: 3, data: "2026-07-19T10:00:00Z", estimativa: false, sem_dados: false },
          frete: { valor: "12.00", fonte: "CATEGORIA", degrau: 2, data: null, estimativa: false, sem_dados: false },
        }}
      />,
    );

    const carimbo = screen.getByTestId("decomp-tarifa-comissao");
    expect(carimbo).toHaveTextContent("Cotação");
    expect(carimbo).toHaveTextContent("degrau 3");
    expect(carimbo).not.toHaveTextContent("ESTIMATIVA");
    expect(within(carimbo).getByLabelText("Atualização dos dados")).toBeInTheDocument();
  });

  it("renders the ESTIMATIVA pill on a degrau-4 PADRAO comissão estimate", () => {
    render(
      <DecompositionPanel
        decomposition={decomposition()}
        profile={profile}
        blockingState={null}
        tarifa={{
          comissao: { valor: "13.00", fonte: "PADRAO", degrau: 4, data: null, estimativa: true, sem_dados: false },
          frete: null,
        }}
      />,
    );

    const carimbo = screen.getByTestId("decomp-tarifa-comissao");
    expect(carimbo).toHaveTextContent("Padrão");
    expect(carimbo).toHaveTextContent("degrau 4");
    expect(carimbo).toHaveTextContent("ESTIMATIVA");
  });

  it("suppresses the frete carimbo and keeps '—' when frete is NO-DATA (valor null / sem_dados)", () => {
    const d = decomposition({ frete: null, componentes_desconhecidos: ["frete"] });
    render(
      <DecompositionPanel
        decomposition={d}
        profile={profile}
        blockingState={null}
        tarifa={{
          comissao: null,
          frete: { valor: null, fonte: "PADRAO", degrau: 4, data: null, estimativa: false, sem_dados: true },
        }}
      />,
    );

    expect(screen.queryByTestId("decomp-tarifa-frete")).toBeNull();
    const panel = screen.getByTestId("decomposition-panel");
    expect(within(panel).getAllByText("—", { selector: "span.text-faint" }).length).toBeGreaterThanOrEqual(1);
    expect(panel).not.toHaveTextContent("ESTIMATIVA");
  });

  it("never stamps a carimbo beside '—' even if tarifa.frete disagrees with a null decomposition.frete (ADR-17)", () => {
    // Pathological backend disagreement: the frete value is unknown (null ⇒ "—") yet the
    // tarifa block carries a resolved frete provenance. The FE must NOT render a source/
    // degrau stamp next to the em-dash — honest-absence is enforced on the shown value.
    const d = decomposition({ frete: null, componentes_desconhecidos: ["frete"] });
    render(
      <DecompositionPanel
        decomposition={d}
        profile={profile}
        blockingState={null}
        tarifa={{
          comissao: null,
          frete: { valor: "12.00", fonte: "COTACAO", degrau: 3, data: "2026-07-19T10:00:00Z", estimativa: false, sem_dados: false },
        }}
      />,
    );

    expect(screen.queryByTestId("decomp-tarifa-frete")).toBeNull();
    const panel = screen.getByTestId("decomposition-panel");
    expect(panel).not.toHaveTextContent("Cotação");
  });

  it("renders no carimbo when no tarifa block is passed (backward compatible)", () => {
    render(<DecompositionPanel decomposition={decomposition()} profile={profile} blockingState={null} />);

    expect(screen.queryByTestId("decomp-tarifa-comissao")).toBeNull();
    expect(screen.queryByTestId("decomp-tarifa-frete")).toBeNull();
    const panel = screen.getByTestId("decomposition-panel");
    expect(within(panel).getByText("17.00")).toBeInTheDocument();
  });

  it("does not add a carimbo to custo/imposto rows even when tarifa is present", () => {
    render(
      <DecompositionPanel
        decomposition={decomposition()}
        profile={profile}
        blockingState={null}
        tarifa={{
          comissao: { valor: "14.5", fonte: "COTACAO", degrau: 3, data: "2026-07-19T10:00:00Z", estimativa: false, sem_dados: false },
          frete: { valor: "12.00", fonte: "CATEGORIA", degrau: 2, data: null, estimativa: false, sem_dados: false },
        }}
      />,
    );

    const panel = screen.getByTestId("decomposition-panel");
    expect(within(panel).getByTestId("decomp-tarifa-comissao")).toBeInTheDocument();
    expect(within(panel).getByTestId("decomp-tarifa-frete")).toBeInTheDocument();
    expect(within(panel).getByText("40.00")).toBeInTheDocument();
    expect(within(panel).getByText("4.00")).toBeInTheDocument();
  });
});
