import { fireEvent, render, screen } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";
import type { PricingCalcProfile } from "@marketplace-central/sdk-runtime";
import { BR_UFS, ParamsDrawer } from "./ParamsDrawer";

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

function renderDrawer(overrides: Partial<PricingCalcProfile> = {}, props: Partial<Parameters<typeof ParamsDrawer>[0]> = {}) {
  const onSave = vi.fn();
  const onClose = vi.fn();
  render(
    <ParamsDrawer open profile={{ ...profile, ...overrides }} onSave={onSave} onClose={onClose} {...props} />,
  );
  return { onSave, onClose };
}

describe("ParamsDrawer", () => {
  it("does not render when closed", () => {
    render(<ParamsDrawer open={false} profile={profile} onSave={vi.fn()} onClose={vi.fn()} />);
    expect(screen.queryByTestId("params-drawer")).toBeNull();
  });

  it("edits the profile and saves the draft", () => {
    const { onSave } = renderDrawer();
    fireEvent.change(screen.getByLabelText("Alíquota"), { target: { value: "9.25" } });
    fireEvent.change(screen.getByLabelText("Limiar verde"), { target: { value: "20" } });
    fireEvent.click(screen.getByRole("button", { name: "Salvar parâmetros" }));

    expect(onSave).toHaveBeenCalledTimes(1);
    expect(onSave.mock.calls[0][0]).toMatchObject({ aliquota_pct: "9.25", limiar_verde_pct: "20" });
  });

  it("client-validates alíquota to 0–35 and blocks the save when out of range", () => {
    const { onSave } = renderDrawer();
    fireEvent.change(screen.getByLabelText("Alíquota"), { target: { value: "40" } });

    expect(screen.getByRole("alert")).toHaveTextContent("entre 0 e 35");
    const save = screen.getByRole("button", { name: "Salvar parâmetros" });
    expect(save).toBeDisabled();
    fireEvent.click(save);
    expect(onSave).not.toHaveBeenCalled();
  });

  it("normalizes pt-BR comma rate fields to dot-decimal in the saved profile (F-P7-2)", () => {
    const { onSave } = renderDrawer();
    // pt-BR "9,25" alíquota must not be a false NaN-block, and must be dot-decimal on save.
    fireEvent.change(screen.getByLabelText("Alíquota"), { target: { value: "9,25" } });
    fireEvent.change(screen.getByLabelText("Limiar verde"), { target: { value: "18,5" } });
    fireEvent.change(screen.getByLabelText("Tarifa Full"), { target: { value: "2,5" } });
    const save = screen.getByRole("button", { name: "Salvar parâmetros" });
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
    fireEvent.click(screen.getByRole("button", { name: "Salvar parâmetros" }));
    expect(onSave.mock.calls[0][0]).toMatchObject({ difal_enabled: true, difal_destino_uf: "SP" });
  });
});
