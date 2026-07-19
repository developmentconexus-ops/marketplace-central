import { useState } from "react";
import type { JSX } from "react";
import { ErrorState, LoadingState } from "@marketplace-central/ui";
import type { PricingDifalListResponse } from "@marketplace-central/sdk-runtime";

/** Mandatory on every DIFAL surface (C02/C05) — fallback if the API omits it. */
export const DIFAL_DISCLAIMER = "seed padrão 2026 — não é orientação fiscal";

export interface DifalDrawerProps {
  open: boolean;
  data: PricingDifalListResponse | undefined;
  isLoading?: boolean;
  isError?: boolean;
  onOverride: (uf: string, internaPct: string) => void;
  onRetry?: () => void;
  onClose: () => void;
  savingUf?: string | null;
}

/**
 * DifalDrawer lists the 27-UF DIFAL rate table (interna / interestadual /
 * efetivo, plus any active override) and lets the operator override a UF's
 * internal rate. The seed disclaimer is shown on the surface (C02/C05) — the
 * table is a seed default, not fiscal guidance.
 */
export function DifalDrawer({
  open, data, isLoading, isError, onOverride, onRetry, onClose, savingUf,
}: DifalDrawerProps): JSX.Element | null {
  const [drafts, setDrafts] = useState<Record<string, string>>({});

  if (!open) return null;

  const disclaimer = data?.disclaimer ?? DIFAL_DISCLAIMER;

  return (
    <aside role="dialog" aria-label="Tabela DIFAL" data-testid="difal-drawer"
      className="flex w-96 flex-col gap-3 border-l border-border bg-surface p-4">
      <div className="flex items-center justify-between">
        <h2 className="text-sm font-semibold text-ink">DIFAL por UF de destino</h2>
        <button type="button" onClick={onClose} aria-label="Fechar" className="text-muted hover:text-ink">✕</button>
      </div>

      <p role="note" data-testid="difal-disclaimer" className="rounded-md bg-surface-2 px-2 py-1 text-xs text-muted">
        {disclaimer}
      </p>

      {isLoading ? (
        <LoadingState />
      ) : isError ? (
        <ErrorState onRetry={onRetry ?? (() => {})} detail="Falha ao carregar a tabela DIFAL." />
      ) : (
        <div className="overflow-x-auto">
          <table className="w-full text-xs">
            <thead>
              <tr className="text-left text-muted">
                <th className="py-1">UF</th>
                <th>Interna</th>
                <th>Inter.</th>
                <th>Efetivo</th>
                <th>Override</th>
              </tr>
            </thead>
            <tbody>
              {(data?.items ?? []).map((rate) => (
                <tr key={rate.uf} className="border-t border-border font-mono">
                  <td className="py-1 text-ink">{rate.uf}</td>
                  <td className="text-muted">{rate.override?.interna_pct ?? rate.interna_pct}</td>
                  <td className="text-muted">{rate.interestadual_pct}</td>
                  <td className="text-ink">{rate.efetivo_pct}</td>
                  <td>
                    <span className="flex items-center gap-1">
                      <input
                        inputMode="decimal"
                        aria-label={`Override interna ${rate.uf}`}
                        value={drafts[rate.uf] ?? ""}
                        placeholder={rate.override ? rate.override.interna_pct : "—"}
                        onChange={(e) => setDrafts((d) => ({ ...d, [rate.uf]: e.target.value }))}
                        className="w-16 rounded border border-border bg-surface px-1 py-0.5 text-ink"
                      />
                      <button
                        type="button"
                        disabled={savingUf === rate.uf || (drafts[rate.uf] ?? "").trim() === ""}
                        onClick={() => onOverride(rate.uf, (drafts[rate.uf] ?? "").trim())}
                        className="rounded px-1.5 py-0.5 text-accent hover:bg-accent-soft disabled:opacity-40"
                      >
                        {savingUf === rate.uf ? "…" : "OK"}
                      </button>
                    </span>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </aside>
  );
}
