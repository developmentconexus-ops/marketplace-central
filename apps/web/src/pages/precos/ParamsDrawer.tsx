import { useEffect, useState } from "react";
import type { JSX } from "react";
import { UnknownValue } from "@marketplace-central/ui";
import type { PricingCalcProfile } from "@marketplace-central/sdk-runtime";
import { ptBrRateToDot } from "./ptbrDecimal";
import type { TariffComponent } from "./tariffBadge";
import { fonteLabel, isNoData } from "./tariffBadge";

/** The 27 Brazilian UFs (26 states + DF), used for the DIFAL destination selector. */
export const BR_UFS = [
  "AC", "AL", "AP", "AM", "BA", "CE", "DF", "ES", "GO", "MA", "MT", "MS", "MG",
  "PA", "PB", "PR", "PE", "PI", "RJ", "RN", "RS", "RO", "RR", "SC", "SP", "SE", "TO",
] as const;

/** The two tax regimes, with the label + seed alíquota the prototype pins to each. */
const REGIMES = [
  { key: "SIMPLES", label: "Simples Nacional", aliquota: "4" },
  { key: "PRESUMIDO", label: "Lucro Presumido", aliquota: "9,25" },
] as const;

/**
 * ML rules the operator cannot edit (they come from Mercado Livre / the ERP), shown
 * verbatim from the design so the drawer reads as "these are givens, not knobs".
 */
const REGRAS_ML =
  "Taxa fixa R$ 6,50 em vendas abaixo de R$ 79 · frete grátis obrigatório " +
  "(vendedor paga) a partir de R$ 79 · frete calculado por peso e CEP · custo do produto vem do ERP";

export interface ParamsDrawerProps {
  open: boolean;
  profile: PricingCalcProfile;
  onSave: (next: PricingCalcProfile) => void;
  onClose: () => void;
  saving?: boolean;
  /** Defaults for the "Restaurar padrão" client-side reset (the page's seed profile). */
  defaults?: PricingCalcProfile;
  /** Active modalidade key ("classico" | "premium" | "full") — drives which ML tier row is live. */
  modalidade?: string;
  /** Resolved comissão tarifa for the active modalidade (display-only, "vem do ML"). */
  comissaoTarifa?: TariffComponent | null;
  /** Opens the full DIFAL-por-UF table (the richer DifalDrawer) from the in-drawer link. */
  onOpenDifal?: () => void;
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
 * ParamsDrawer is the "Parâmetros de cálculo" overlay (design Simulador.dc.html): a
 * fixed right-hand sheet that edits the operator's CalcProfile — tax regime + alíquota,
 * the verde/amarelo margin thresholds, the DIFAL toggle + destino UF — and surfaces the
 * read-only ML context (commission tiers "vem do ML", the fixed ML rules, the DIFAL
 * destination from the analysis CEP). Editable knobs vs. ML givens are visually split.
 * Alíquota is client-validated to 0–35 before the PUT (server re-checks → 422
 * INVALID_RATE). Money/percentages stay decimal strings. "Restaurar padrão" resets the
 * draft client-side; "Concluído" persists it (the page closes on save success).
 */
export function ParamsDrawer({
  open, profile, onSave, onClose, saving, defaults, modalidade, comissaoTarifa, onOpenDifal,
}: ParamsDrawerProps): JSX.Element | null {
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
    <>
      {/* Scrim — click to dismiss, matching the design's fixed overlay. */}
      <div
        data-testid="params-backdrop"
        onClick={onClose}
        className="fixed inset-0 z-[40] bg-ink/35"
      />
      <aside
        role="dialog"
        aria-label="Parâmetros de cálculo"
        data-testid="params-drawer"
        className="fixed inset-y-0 right-0 z-[41] flex w-[400px] max-w-[92vw] flex-col border-l border-border bg-surface text-[12.5px] shadow-[-16px_0_40px_rgba(0,0,0,0.14)]"
      >
        {/* Header */}
        <div className="flex items-center gap-2 border-b border-border bg-surface-2 px-[18px] py-[13px]">
          <b className="text-[13.5px] text-ink">Parâmetros de cálculo</b>
          <span
            data-testid="params-badge-live"
            className="whitespace-nowrap rounded-pill bg-accent-soft px-2 py-0.5 text-[10.5px] font-bold text-accent-ink"
          >
            recalcula ao vivo
          </span>
          <button
            type="button"
            onClick={onClose}
            aria-label="Fechar"
            className="ml-auto px-1.5 py-0.5 text-faint hover:text-ink"
          >
            ✕
          </button>
        </div>

        {/* Body */}
        <div className="flex flex-1 flex-col gap-5 overflow-y-auto px-[18px] py-4">
          {/* IMPOSTO */}
          <section className="flex flex-col gap-2">
            <SectionLabel>IMPOSTO</SectionLabel>
            <div className="flex gap-1.5">
              {REGIMES.map((rg) => {
                const on = draft.regime === rg.key;
                return (
                  <button
                    key={rg.key}
                    type="button"
                    aria-pressed={on}
                    onClick={() => setDraft((d) => ({ ...d, regime: rg.key, aliquota_pct: rg.aliquota }))}
                    className={`whitespace-nowrap rounded-pill border px-3 py-1 font-semibold ${
                      on
                        ? "border-accent bg-accent-soft text-accent-ink"
                        : "border-border text-muted hover:border-accent hover:text-accent-ink"
                    }`}
                  >
                    {rg.label}
                  </button>
                );
              })}
            </div>
            <div className="flex items-center gap-2.5">
              <span className="flex-1 text-muted">Alíquota sobre o preço de venda</span>
              <PctInput aria="Alíquota" value={draft.aliquota_pct} invalid={aliquotaError}
                onChange={(v) => set("aliquota_pct", v)} />
            </div>
            {aliquotaError ? (
              <p role="alert" className="text-[11px] text-warn">Alíquota deve estar entre 0 e 35.</p>
            ) : null}
          </section>

          {/* LIMIARES DE MARGEM */}
          <section className="flex flex-col gap-2">
            <SectionLabel>LIMIARES DE MARGEM</SectionLabel>
            <div className="flex items-center gap-2.5">
              <span className="h-2 w-2 flex-none rounded-pill bg-accent" aria-hidden />
              <span className="flex-1 text-muted">Saudável a partir de</span>
              <PctInput aria="Limiar verde" value={draft.limiar_verde_pct}
                onChange={(v) => set("limiar_verde_pct", v)} />
            </div>
            <div className="flex items-center gap-2.5">
              <span className="h-2 w-2 flex-none rounded-pill bg-amber" aria-hidden />
              <span className="flex-1 text-muted">Apertada a partir de</span>
              <PctInput aria="Limiar amarela" value={draft.limiar_amarelo_pct}
                onChange={(v) => set("limiar_amarelo_pct", v)} />
            </div>
            <p className="text-[11px] text-faint">
              abaixo disso o veredicto marca "não vale" em vermelho
            </p>
          </section>

          {/* MODALIDADES ML — read-only, comes from ML */}
          <section data-testid="params-modalidades-ml" className="flex flex-col gap-2">
            <SectionLabel>MODALIDADES ML</SectionLabel>
            <ModalidadeRow label="Comissão Clássico" active={modalidade === "classico"} tarifa={comissaoTarifa} />
            <ModalidadeRow label="Comissão Premium" active={modalidade === "premium"} tarifa={comissaoTarifa} />
            <div className="flex items-center gap-2.5">
              <span className="flex-1 text-muted">Tarifa Full por unidade</span>
              <span className="flex w-[78px] flex-none items-center gap-1 rounded-control border border-border px-2.5 py-[5px]">
                <span className="text-[11.5px] text-faint">R$</span>
                <input
                  inputMode="decimal"
                  aria-label="Tarifa Full"
                  value={draft.tarifa_full ?? ""}
                  placeholder="—"
                  onChange={(e) => set("tarifa_full", e.target.value.trim() === "" ? null : e.target.value)}
                  className="w-full border-none bg-transparent text-right font-mono text-[13px] font-semibold text-ink outline-none"
                />
              </span>
            </div>
          </section>

          {/* REGRAS DO ML — NÃO EDITÁVEIS */}
          <div
            data-testid="params-regras-ml"
            className="flex flex-col gap-1.5 rounded-card border border-border-2 bg-surface-2 px-3 py-2.5"
          >
            <SectionLabel>REGRAS DO ML — NÃO EDITÁVEIS</SectionLabel>
            <p className="text-[11.5px] leading-relaxed text-muted">{REGRAS_ML}</p>
          </div>

          {/* DIFAL */}
          <section className="flex flex-col gap-2">
            <div className="flex items-center gap-2">
              <SectionLabel className="flex-1">DIFAL</SectionLabel>
              <button
                type="button"
                aria-label="Habilitar DIFAL"
                aria-pressed={draft.difal_enabled}
                onClick={() => set("difal_enabled", !draft.difal_enabled)}
                className={`whitespace-nowrap rounded-pill border px-3 py-[3px] text-[11.5px] font-semibold ${
                  draft.difal_enabled
                    ? "border-accent bg-accent-soft text-accent-ink"
                    : "border-border text-muted"
                }`}
              >
                {draft.difal_enabled ? "considerando ✓" : "ignorar"}
              </button>
            </div>
            <p data-testid="params-difal-context" className="text-[11.5px] leading-snug text-muted">
              Desconta o diferencial de alíquota da UF de destino na margem simulada.{" "}
              {draft.difal_destino_uf ? (
                <>Destino atual: <b className="text-ink">{draft.difal_destino_uf}</b> (do CEP da análise).</>
              ) : (
                <>Informe a UF de destino para descontar o DIFAL.</>
              )}
            </p>
            {draft.difal_enabled ? (
              <label className="flex items-center gap-2.5 text-muted">
                <span className="flex-1">UF de destino</span>
                <select
                  aria-label="UF de destino"
                  value={draft.difal_destino_uf ?? ""}
                  onChange={(e) => set("difal_destino_uf", e.target.value === "" ? null : e.target.value)}
                  className="w-[78px] flex-none rounded-control border border-border bg-surface px-2 py-[5px] text-[13px] text-ink"
                >
                  <option value="">—</option>
                  {BR_UFS.map((uf) => <option key={uf} value={uf}>{uf}</option>)}
                </select>
              </label>
            ) : null}
            {onOpenDifal ? (
              <button
                type="button"
                onClick={onOpenDifal}
                className="self-start text-[11.5px] font-semibold text-accent-ink hover:opacity-75"
              >
                ver tabela completa por UF →
              </button>
            ) : null}
          </section>
        </div>

        {/* Footer */}
        <div className="flex gap-2 border-t border-border px-[18px] py-3">
          <button
            type="button"
            data-testid="params-reset"
            onClick={() => defaults && setDraft(defaults)}
            disabled={!defaults}
            className="whitespace-nowrap rounded-control border border-border px-3.5 py-[9px] text-muted hover:border-warn hover:text-warn disabled:opacity-40"
          >
            Restaurar padrão
          </button>
          <button
            type="button"
            disabled={aliquotaError || saving}
            onClick={() => onSave(normalizeProfile(draft))}
            className="flex-1 rounded-control bg-accent px-3 py-[9px] text-center font-semibold text-accent-ink hover:opacity-90 disabled:opacity-50"
          >
            {saving ? "Salvando…" : "Concluído"}
          </button>
        </div>
      </aside>
    </>
  );
}

function SectionLabel({ children, className = "" }: { children: string; className?: string }): JSX.Element {
  return (
    <div className={`text-[10.5px] font-semibold tracking-[0.08em] text-faint ${className}`}>{children}</div>
  );
}

/** A boxed % input (64px) with the right-aligned mono value the design uses. */
function PctInput({ aria, value, onChange, invalid }: {
  aria: string; value: string; onChange: (v: string) => void; invalid?: boolean;
}): JSX.Element {
  return (
    <span className="flex w-16 flex-none items-center gap-1 rounded-control border border-border px-2.5 py-[5px]">
      <input
        inputMode="decimal"
        aria-label={aria}
        aria-invalid={invalid}
        value={value}
        onChange={(e) => onChange(e.target.value)}
        className="w-full border-none bg-transparent text-right font-mono text-[13px] font-semibold text-ink outline-none"
      />
      <span className="text-[11.5px] text-faint">%</span>
    </span>
  );
}

/**
 * A "Comissão {tier} — vem do ML" row. The percentage is display-only: we only know the
 * live tarifa for the ACTIVE modalidade, so the matching tier shows it (honest "—" via
 * UnknownValue when NO-DATA) and the other tier shows "—" rather than a fabricated rate.
 */
function ModalidadeRow({ label, active, tarifa }: {
  label: string; active: boolean; tarifa?: TariffComponent | null;
}): JSX.Element {
  const show = active && tarifa != null && !isNoData(tarifa);
  const fonte = show ? fonteLabel(tarifa!.fonte) : null;
  return (
    <div className="flex items-center gap-2.5">
      <span className="flex-1 text-muted">{label}</span>
      {show ? (
        <span className="font-mono font-semibold text-ink">{tarifa!.valor}%</span>
      ) : (
        <UnknownValue />
      )}
      <span className="text-[10.5px] text-faint">{fonte ? `${fonte} · vem do ML` : "vem do ML"}</span>
    </div>
  );
}
