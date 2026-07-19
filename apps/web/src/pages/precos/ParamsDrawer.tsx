import { useEffect, useState } from "react";
import type { JSX } from "react";
import type { PricingCalcProfile } from "@marketplace-central/sdk-runtime";
import { ptBrRateToDot } from "./ptbrDecimal";

/** The 27 Brazilian UFs (26 states + DF), used for the DIFAL destination selector. */
export const BR_UFS = [
  "AC", "AL", "AP", "AM", "BA", "CE", "DF", "ES", "GO", "MA", "MT", "MS", "MG",
  "PA", "PB", "PR", "PE", "PI", "RJ", "RN", "RS", "RO", "RR", "SC", "SP", "SE", "TO",
] as const;

const REGIMES = ["SIMPLES", "PRESUMIDO"] as const;

export interface ParamsDrawerProps {
  open: boolean;
  profile: PricingCalcProfile;
  onSave: (next: PricingCalcProfile) => void;
  onClose: () => void;
  saving?: boolean;
}

/** aliquota is IC-04-bounded to [0,35]; the server re-validates (422 INVALID_RATE). */
function aliquotaValid(value: string): boolean {
  // Validate the normalized (dot-decimal) form so pt-BR "9,25" isn't a false NaN.
  const n = Number(ptBrRateToDot(value));
  return Number.isFinite(n) && n >= 0 && n <= 35;
}

/** Promote the operator's raw pt-BR rate fields to dot-decimal for the PUT. */
function normalizeProfile(p: PricingCalcProfile): PricingCalcProfile {
  return {
    ...p,
    aliquota_pct: ptBrRateToDot(p.aliquota_pct),
    limiar_verde_pct: ptBrRateToDot(p.limiar_verde_pct),
    limiar_amarelo_pct: ptBrRateToDot(p.limiar_amarelo_pct),
    tarifa_full: p.tarifa_full === null ? null : ptBrRateToDot(p.tarifa_full),
  };
}

/**
 * ParamsDrawer edits the operator's CalcProfile — tax regime + alíquota, the
 * verde/amarelo margin thresholds, tarifa Full, and the DIFAL toggle + destino
 * UF. Alíquota is client-validated to 0–35 before the PUT (the server re-checks
 * and answers 422 INVALID_RATE). Money/percentages stay decimal strings.
 */
export function ParamsDrawer({ open, profile, onSave, onClose, saving }: ParamsDrawerProps): JSX.Element | null {
  const [draft, setDraft] = useState<PricingCalcProfile>(profile);

  // Re-seed the draft from the persisted profile whenever the drawer (re)opens.
  useEffect(() => {
    if (open) setDraft(profile);
  }, [open, profile]);

  if (!open) return null;

  const aliquotaError = !aliquotaValid(draft.aliquota_pct);

  function set<K extends keyof PricingCalcProfile>(key: K, value: PricingCalcProfile[K]) {
    setDraft((d) => ({ ...d, [key]: value }));
  }

  return (
    <aside role="dialog" aria-label="Parâmetros de cálculo" data-testid="params-drawer"
      className="flex w-80 flex-col gap-4 border-l border-border bg-surface p-4">
      <div className="flex items-center justify-between">
        <h2 className="text-sm font-semibold text-ink">Parâmetros de cálculo</h2>
        <button type="button" onClick={onClose} aria-label="Fechar" className="text-muted hover:text-ink">✕</button>
      </div>

      <fieldset className="flex flex-col gap-1">
        <legend className="text-xs font-medium text-muted">Imposto</legend>
        <label className="flex flex-col text-xs text-muted">
          Regime
          <select
            aria-label="Regime"
            value={draft.regime}
            onChange={(e) => set("regime", e.target.value)}
            className="mt-1 rounded-md border border-border bg-surface px-2 py-1 text-sm text-ink"
          >
            {REGIMES.map((r) => <option key={r} value={r}>{r}</option>)}
          </select>
        </label>
        <label className="flex flex-col text-xs text-muted">
          Alíquota (%)
          <input
            inputMode="decimal"
            aria-label="Alíquota"
            aria-invalid={aliquotaError}
            value={draft.aliquota_pct}
            onChange={(e) => set("aliquota_pct", e.target.value)}
            className="mt-1 w-24 rounded-md border border-border bg-surface px-2 py-1 font-mono text-sm text-ink"
          />
        </label>
        {aliquotaError ? (
          <p role="alert" className="text-xs text-warn">Alíquota deve estar entre 0 e 35.</p>
        ) : null}
      </fieldset>

      <fieldset className="flex flex-col gap-1">
        <legend className="text-xs font-medium text-muted">Limiares de margem</legend>
        <label className="flex flex-col text-xs text-muted">
          Verde ≥ (%)
          <input inputMode="decimal" aria-label="Limiar verde" value={draft.limiar_verde_pct}
            onChange={(e) => set("limiar_verde_pct", e.target.value)}
            className="mt-1 w-24 rounded-md border border-border bg-surface px-2 py-1 font-mono text-sm text-ink" />
        </label>
        <label className="flex flex-col text-xs text-muted">
          Amarela ≥ (%)
          <input inputMode="decimal" aria-label="Limiar amarela" value={draft.limiar_amarelo_pct}
            onChange={(e) => set("limiar_amarelo_pct", e.target.value)}
            className="mt-1 w-24 rounded-md border border-border bg-surface px-2 py-1 font-mono text-sm text-ink" />
        </label>
      </fieldset>

      <label className="flex flex-col text-xs text-muted">
        Tarifa Full (%)
        <input inputMode="decimal" aria-label="Tarifa Full" value={draft.tarifa_full ?? ""}
          placeholder="—"
          onChange={(e) => set("tarifa_full", e.target.value.trim() === "" ? null : e.target.value)}
          className="mt-1 w-24 rounded-md border border-border bg-surface px-2 py-1 font-mono text-sm text-ink" />
      </label>

      <fieldset className="flex flex-col gap-1">
        <legend className="text-xs font-medium text-muted">DIFAL</legend>
        <label className="flex items-center gap-2 text-sm text-ink">
          <input type="checkbox" aria-label="Habilitar DIFAL" checked={draft.difal_enabled}
            onChange={(e) => set("difal_enabled", e.target.checked)} />
          Aplicar DIFAL na decomposição
        </label>
        {draft.difal_enabled ? (
          <label className="flex flex-col text-xs text-muted">
            UF de destino
            <select aria-label="UF de destino" value={draft.difal_destino_uf ?? ""}
              onChange={(e) => set("difal_destino_uf", e.target.value === "" ? null : e.target.value)}
              className="mt-1 w-24 rounded-md border border-border bg-surface px-2 py-1 text-sm text-ink">
              <option value="">—</option>
              {BR_UFS.map((uf) => <option key={uf} value={uf}>{uf}</option>)}
            </select>
          </label>
        ) : null}
      </fieldset>

      <div className="mt-2 flex justify-end gap-2">
        <button type="button" onClick={onClose} className="rounded-md px-3 py-1.5 text-sm text-muted hover:text-ink">
          Cancelar
        </button>
        <button
          type="button"
          disabled={aliquotaError || saving}
          onClick={() => onSave(normalizeProfile(draft))}
          className="rounded-md bg-accent px-3 py-1.5 text-sm text-accent-ink disabled:opacity-50"
        >
          {saving ? "Salvando…" : "Salvar parâmetros"}
        </button>
      </div>
    </aside>
  );
}
