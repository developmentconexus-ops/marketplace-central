import { fireEvent, render, screen } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";
import type { PricingCalcProfile } from "@marketplace-central/sdk-runtime";
import { BR_UFS, ParamsDrawer } from "./ParamsDrawer";
import type { TariffComponent } from "./tariffBadge";

const profile: PricingCalcProfile = {
  regime: "SIMPLES",
  aliquota_pct: "4",
  limiar_verde_pct: "18",
  limiar_amarelo_pct: "10",
  tarifa_full: null,
  difal_enabled: false,
  difal_destino_uf: null,
  origem: "operator",
};

function renderDrawer(
  overrides: Partial<PricingCalcProfile> = {},
  props: Partial<Parameters<typeof ParamsDrawer>[0]> = {},
) {
  const onSave = vi.fn();
  const onClose = vi.fn();
  render(
    <ParamsDrawer
      open
      profile={{ ...profile, ...overrides }}
      onSave={onSave}
      onClose={onClose}
      {...props}
    />,
  );
  return { onSave, onClose };
}

describe("ParamsDrawer", () => {
  it("does not render when closed", () => {
    render(<ParamsDrawer open={false} profile={profile} onSave={vi.fn()} onClose={vi.fn()} />);
    expect(screen.queryByTestId("params-drawer")).toBeNull();
  });

  it("edits the profile and saves the draft on Concluído", () => {
    const { onSave } = renderDrawer();
    fireEvent.change(screen.getByLabelText("Alíquota"), { target: { value: "9.25" } });
    fireEvent.change(screen.getByLabelText("Limiar verde"), { target: { value: "20" } });
    fireEvent.click(screen.getByRole("button", { name: "Concluído" }));

    expect(onSave).toHaveBeenCalledTimes(1);
    expect(onSave.mock.calls[0][0]).toMatchObject({ aliquota_pct: "9.25", limiar_verde_pct: "20" });
  });

  it("client-validates alíquota to 0–35 and blocks the save when out of range", () => {
    const { onSave } = renderDrawer();
    fireEvent.change(screen.getByLabelText("Alíquota"), { target: { value: "40" } });

    expect(screen.getByRole("alert")).toHaveTextContent("entre 0 e 35");
    const save = screen.getByRole("button", { name: "Concluído" });
    expect(save).toBeDisabled();
    fireEvent.click(save);
    expect(onSave).not.toHaveBeenCalled();
  });

  it("normalizes pt-BR comma rate fields to dot-decimal in the saved profile (F-P7-2)", () => {
    const { onSave } = renderDrawer();
    fireEvent.change(screen.getByLabelText("Alíquota"), { target: { value: "9,25" } });
    fireEvent.change(screen.getByLabelText("Limiar verde"), { target: { value: "18,5" } });
    fireEvent.change(screen.getByLabelText("Tarifa Full"), { target: { value: "2,5" } });
    const save = screen.getByRole("button", { name: "Concluído" });
    expect(save).not.toBeDisabled();
    fireEvent.click(save);

    expect(onSave).toHaveBeenCalledTimes(1);
    expect(onSave.mock.calls[0][0]).toMatchObject({
      aliquota_pct: "9.25",
      limiar_verde_pct: "18.5",
      tarifa_full: "2.5",
    });
  });

  it("reveals the 27-UF destino selector only when DIFAL is enabled", () => {
    renderDrawer();
    expect(screen.queryByLabelText("UF de destino")).toBeNull();

    fireEvent.click(screen.getByLabelText("Habilitar DIFAL"));
    const select = screen.getByLabelText("UF de destino");
    expect(select).toBeInTheDocument();
    // 27 UFs + the "—" empty option.
    expect(select.querySelectorAll("option")).toHaveLength(BR_UFS.length + 1);
  });

  it("saves the chosen destino UF for a re-decompose", () => {
    const { onSave } = renderDrawer({ difal_enabled: true });
    fireEvent.change(screen.getByLabelText("UF de destino"), { target: { value: "SP" } });
    fireEvent.click(screen.getByRole("button", { name: "Concluído" }));
    expect(onSave.mock.calls[0][0]).toMatchObject({ difal_enabled: true, difal_destino_uf: "SP" });
  });

  // ── Design-parity surface (Simulador.dc.html) ─────────────────────────────

  it("shows the 'recalcula ao vivo' header badge", () => {
    renderDrawer();
    expect(screen.getByTestId("params-badge-live")).toHaveTextContent("recalcula ao vivo");
  });

  it("renders the read-only MODALIDADES ML box with 'vem do ML' tiers + the resolved active rate", () => {
    // Active modalidade = premium; the resolved comissão tarifa is display-only.
    const comissao: TariffComponent = {
      valor: "17",
      fonte: "COTACAO",
      degrau: 3,
      data: "2026-07-18T12:00:00Z",
    };
    renderDrawer({}, { modalidade: "premium", comissaoTarifa: comissao });
    const box = screen.getByTestId("params-modalidades-ml");
    expect(box).toHaveTextContent("Comissão Clássico");
    expect(box).toHaveTextContent("Comissão Premium");
    expect(box).toHaveTextContent("vem do ML");
    // Premium is active ⇒ its resolved % shows; Clássico stays honest "—" (never fabricated).
    expect(box).toHaveTextContent("17%");
  });

  it("renders the static 'REGRAS DO ML — NÃO EDITÁVEIS' box verbatim", () => {
    renderDrawer();
    const box = screen.getByTestId("params-regras-ml");
    expect(box).toHaveTextContent("REGRAS DO ML — NÃO EDITÁVEIS");
    expect(box).toHaveTextContent("Taxa fixa R$ 6,50 em vendas abaixo de R$ 79");
    expect(box).toHaveTextContent("custo do produto vem do ERP");
  });

  it("shows the in-drawer DIFAL destino context from the analysis CEP", () => {
    renderDrawer({ difal_enabled: true, difal_destino_uf: "SP" });
    expect(screen.getByTestId("params-difal-context")).toHaveTextContent(
      "Destino atual: SP (do CEP da análise).",
    );
  });

  it("opens the full DIFAL-por-UF table via the in-drawer link", () => {
    const onOpenDifal = vi.fn();
    renderDrawer({}, { onOpenDifal });
    fireEvent.click(screen.getByRole("button", { name: /ver tabela completa por UF/ }));
    expect(onOpenDifal).toHaveBeenCalledTimes(1);
  });

  it("restores drawer fields to defaults on 'Restaurar padrão' (client-side, no save)", () => {
    const defaults: PricingCalcProfile = {
      ...profile,
      aliquota_pct: "4",
      limiar_verde_pct: "18",
      limiar_amarelo_pct: "10",
    };
    const { onSave } = renderDrawer({ aliquota_pct: "9.25", limiar_verde_pct: "25" }, { defaults });
    // Edit away from defaults, then restore.
    fireEvent.change(screen.getByLabelText("Alíquota"), { target: { value: "12" } });
    fireEvent.click(screen.getByRole("button", { name: "Restaurar padrão" }));

    expect((screen.getByLabelText("Alíquota") as HTMLInputElement).value).toBe("4");
    expect((screen.getByLabelText("Limiar verde") as HTMLInputElement).value).toBe("18");
    expect((screen.getByLabelText("Limiar amarela") as HTMLInputElement).value).toBe("10");
    // Restore is client-side only — never persists.
    expect(onSave).not.toHaveBeenCalled();
  });
});
